# goensemblrest

[![Go Reference](https://pkg.go.dev/badge/github.com/gawbul/goensemblrest.svg)](https://pkg.go.dev/github.com/gawbul/goensemblrest)
[![Go Report Card](https://goreportcard.com/badge/github.com/gawbul/goensemblrest)](https://goreportcard.com/report/github.com/gawbul/goensemblrest)
[![CI Tests](https://github.com/gawbul/goensemblrest/actions/workflows/pull_request.yaml/badge.svg)](https://github.com/gawbul/goensemblrest/actions/workflows/pull_request.yaml)
[![Nightly Drift Check](https://github.com/gawbul/goensemblrest/actions/workflows/nightly.yaml/badge.svg)](https://github.com/gawbul/goensemblrest/actions/workflows/nightly.yaml)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](https://opensource.org/licenses/MIT)

`goensemblrest` is an idiomatic, high-performance Go client library for the [Ensembl REST API](https://rest.ensembl.org/).

Designed as the modern Go counterpart to [`pyEnsemblRest`](https://github.com/gawbul/pyEnsemblRest), it brings full endpoint parity, type safety, zero external dependencies, robust rate limiting, automatic transient-error retries, and comprehensive `context.Context` support.

---

## Features

- **Standard Library First**: **Zero external dependencies** — builds exclusively on the Go standard library (`net/http`, `context`, `encoding/json`, `sync`, `time`, `net/url`).
- **Modern Go**: Requires and targets **Go 1.26.x** (with backwards compatibility down to Go 1.24+).
- **100+ Endpoints**: Full coverage across Archive, Comparative Genomics, Cross References, Information, Linkage Disequilibrium, Lookup, Mapping, Ontologies, Overlap, Phenotypes, Regulation, Sequences, Transcript Haplotypes, VEP, Variation, and GA4GH.
- **Sliding-Window Rate Limiter**: Thread-safe client-side rate limiting (defaulting to Ensembl's 15 req/sec limit) that prevents IP throttling and respects `Retry-After` response headers.
- **Automatic Retries**: Exponential backoff for transient server errors (HTTP 500, HTTP 408) and recoverable API conditions.
- **Idiomatic Error Handling**: Deep integration with `errors.Is` (e.g. `ensembl.ErrNotFound`, `ensembl.ErrBadRequest`, `ensembl.ErrRateLimit`) and structured `*ensembl.APIError` unwrapping via `errors.As`.
- **Dynamic & Static Dispatch**: Use strongly-typed domain methods or dynamic dispatch via `client.Call(ctx, "apiName", params, body, target)`.
- **Context-Aware**: Full cancellation and deadline support throughout network calls and rate limiting backoff sleeps.

---

## Installation

```bash
go get github.com/gawbul/goensemblrest
```

---

## Quick Start

```go
package main

import (
	"context"
	"fmt"
	"log"
	"time"

	ensembl "github.com/gawbul/goensemblrest"
)

func main() {
	// Initialize client with defaults (15 req/sec, 60s timeout)
	client, err := ensembl.NewClient()
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	defer client.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Lookup gene information
	var gene ensembl.LookupRecord
	err = client.GetLookupBySymbol(ctx, "homo_sapiens", "BRCA2", &gene, ensembl.WithQuery("expand", "1"))
	if err != nil {
		log.Fatalf("lookup failed: %v", err)
	}

	fmt.Printf("Gene: %s (%s)\n", gene.DisplayName, gene.ID)
	fmt.Printf("Biotype: %s\n", gene.Biotype)
	fmt.Printf("Location: %s:%d-%d (Strand: %d)\n", gene.SeqRegion, gene.Start, gene.End, gene.Strand)
	fmt.Printf("Transcripts: %d\n", len(gene.Transcripts))
}
```

---

## Client Configuration & Options

Configure the client using functional options:

```go
client, err := ensembl.NewClient(
	// Custom base URL (e.g. local mirror or Grch37 server)
	ensembl.WithBaseURL("https://grch37.rest.org"),

	// Custom request timeout
	ensembl.WithTimeout(30 * time.Second),

	// Custom rate limiting (e.g. 10 requests per 1 second window)
	ensembl.WithRateLimit(10, time.Second),

	// Maximum retry attempts for transient server errors
	ensembl.WithMaxAttempts(3),

	// Custom User-Agent
	ensembl.WithUserAgent("MyAnalysisPipeline/1.0"),

	// Default persistent header
	ensembl.WithHeader("X-Custom-Pipeline", "GenomicsV1"),

	// Custom *http.Client (e.g. proxy, custom transport)
	ensembl.WithHTTPClient(&http.Client{
		Timeout: 45 * time.Second,
	}),
)
```

---

## Request Customization

Individual endpoint calls accept `RequestOption` modifiers:

```go
// Add query parameters
client.GetInfoAssembly(ctx, "homo_sapiens", &res, ensembl.WithQuery("bands", "1"))

// Add multiple query parameters
client.GetLookupByID(ctx, "ENSG00000157764", &res, ensembl.WithQueryParams(map[string]string{
	"expand": "1",
	"phenotypes": "1",
}))

// Request alternate content formats (e.g. FASTA, GFF3, plain text)
var fastaSeq string
client.GetSequenceByID(ctx, "ENSG00000157764", &fastaSeq, ensembl.WithRequestContentType("text/x-fasta"))

// Custom request header
client.GetInfoPing(ctx, &ping, ensembl.WithRequestHeader("X-Trace-ID", "abc-123"))
```

---

## Error Handling

`goensemblrest` maps HTTP status codes and Ensembl error responses into structured `*ensembl.APIError` types compatible with standard Go `errors.Is` and `errors.As`:

```go
var record ensembl.ArchiveRecord
err := client.GetArchiveByID(ctx, "UNKNOWN_IDENTIFIER", &record)
if err != nil {
	// 1. Sentinel error checking with errors.Is
	switch {
	case errors.Is(err, ensembl.ErrNotFound):
		fmt.Println("Identifier not found (404)")
	case errors.Is(err, ensembl.ErrBadRequest):
		fmt.Println("Invalid query parameter or ID format (400)")
	case errors.Is(err, ensembl.ErrRateLimit):
		fmt.Println("Rate limit exceeded (429)")
	case errors.Is(err, ensembl.ErrServiceUnavailable):
		fmt.Println("Ensembl service is temporarily unavailable (503)")
	case errors.Is(err, ensembl.ErrMaxRetriesReached):
		fmt.Println("Exhausted maximum retry attempts")
	}

	// 2. Structured error inspection with errors.As
	var apiErr *ensembl.APIError
	if errors.As(err, &apiErr) {
		fmt.Printf("HTTP Status: %d\n", apiErr.StatusCode)
		fmt.Printf("Ensembl Message: %s\n", apiErr.Message)
		if apiErr.RetryAfter != nil {
			fmt.Printf("Retry After: %.1f seconds\n", *apiErr.RetryAfter)
		}
	}
}
```

---

## Endpoint Categories

### 1. Archive
- `GetArchiveByID(ctx, id, target, opts...)`
- `GetArchiveByMultipleIDs(ctx, ids, target, opts...)`

### 2. Comparative Genomics
- `GetCafeGeneTreeByID(ctx, id, target, opts...)`
- `GetCafeGeneTreeMemberBySymbol(ctx, species, symbol, target, opts...)`
- `GetCafeGeneTreeMemberByID(ctx, species, id, target, opts...)`
- `GetGeneTreeByID(ctx, id, target, opts...)`
- `GetGeneTreeMemberBySymbol(ctx, species, symbol, target, opts...)`
- `GetGeneTreeMemberByID(ctx, species, id, target, opts...)`
- `GetAlignmentByRegion(ctx, species, region, target, opts...)`
- `GetHomologyByID(ctx, species, id, target, opts...)`
- `GetHomologyBySymbol(ctx, species, symbol, target, opts...)`

### 3. Cross References (Xrefs)
- `GetXrefsBySymbol(ctx, species, symbol, target, opts...)`
- `GetXrefsByID(ctx, id, target, opts...)`
- `GetXrefsByName(ctx, species, name, target, opts...)`

### 4. Information (Info)
- `GetInfoAnalysis(ctx, species, target, opts...)`
- `GetInfoAssembly(ctx, species, target, opts...)`
- `GetInfoAssemblyRegion(ctx, species, regionName, target, opts...)`
- `GetInfoBiotypes(ctx, species, target, opts...)`
- `GetInfoBiotypesByGroup(ctx, group, objectType, target, opts...)`
- `GetInfoBiotypesByName(ctx, name, objectType, target, opts...)`
- `GetInfoComparaMethods(ctx, target, opts...)`
- `GetInfoComparaSpeciesSets(ctx, methods, target, opts...)`
- `GetInfoComparas(ctx, target, opts...)`
- `GetInfoData(ctx, target, opts...)`
- `GetInfoEgVersion(ctx, target, opts...)`
- `GetInfoExternalDbs(ctx, species, target, opts...)`
- `GetInfoDivisions(ctx, target, opts...)`
- `GetInfoGenomesByName(ctx, name, target, opts...)`
- `GetInfoGenomesByAccession(ctx, accession, target, opts...)`
- `GetInfoGenomesByAssembly(ctx, assemblyID, target, opts...)`
- `GetInfoGenomesByDivision(ctx, division, target, opts...)`
- `GetInfoGenomesByTaxonomy(ctx, taxonName, target, opts...)`
- `GetInfoPing(ctx, target, opts...)`
- `GetInfoRest(ctx, target, opts...)`
- `GetInfoSoftware(ctx, target, opts...)`
- `GetInfoSpecies(ctx, target, opts...)`
- `GetInfoVariationBySpecies(ctx, species, target, opts...)`
- `GetInfoVariationConsequenceTypes(ctx, target, opts...)`
- `GetInfoVariationPopulationIndividuals(ctx, species, populationName, target, opts...)`
- `GetInfoVariationPopulations(ctx, species, target, opts...)`

### 5. Linkage Disequilibrium (LD)
- `GetLdID(ctx, species, id, populationName, target, opts...)`
- `GetLdPairwise(ctx, species, id1, id2, target, opts...)`
- `GetLdRegion(ctx, species, region, populationName, target, opts...)`

### 6. Lookup
- `GetLookupByID(ctx, id, target, opts...)`
- `GetLookupByMultipleIDs(ctx, ids, target, opts...)`
- `GetLookupBySymbol(ctx, species, symbol, target, opts...)`
- `GetLookupByMultipleSymbols(ctx, species, symbols, target, opts...)`

### 7. Mapping
- `GetMapCdnaToRegion(ctx, id, region, target, opts...)`
- `GetMapCdsToRegion(ctx, id, region, target, opts...)`
- `GetMapAssemblyOneToTwo(ctx, species, asmOne, region, asmTwo, target, opts...)`
- `GetMapTranslationToRegion(ctx, id, region, target, opts...)`

### 8. Ontologies & Taxonomy
- `GetAncestorsByID(ctx, id, target, opts...)`
- `GetAncestorsChartByID(ctx, id, target, opts...)`
- `GetDescendantsByID(ctx, id, target, opts...)`
- `GetOntologyByID(ctx, id, target, opts...)`
- `GetOntologyByName(ctx, name, target, opts...)`
- `GetTaxonomyClassificationByID(ctx, id, target, opts...)`
- `GetTaxonomyByID(ctx, id, target, opts...)`
- `GetTaxonomyByName(ctx, name, target, opts...)`

### 9. Overlap
- `GetOverlapByID(ctx, id, target, opts...)`
- `GetOverlapByRegion(ctx, species, region, target, opts...)`
- `GetOverlapByTranslation(ctx, id, target, opts...)`

### 10. Phenotype Annotations
- `GetPhenotypeByAccession(ctx, species, accession, target, opts...)`
- `GetPhenotypeByGene(ctx, species, gene, target, opts...)`
- `GetPhenotypeByRegion(ctx, species, region, target, opts...)`
- `GetPhenotypeByTerm(ctx, species, term, target, opts...)`

### 11. Regulation
- `GetRegulationBindingMatrix(ctx, species, bindingMatrix, target, opts...)`

### 12. Sequences
- `GetSequenceByID(ctx, id, target, opts...)`
- `GetSequenceByMultipleIDs(ctx, ids, target, opts...)`
- `GetSequenceByRegion(ctx, species, region, target, opts...)`
- `GetSequenceByMultipleRegions(ctx, species, regions, target, opts...)`

### 13. Transcript Haplotypes
- `GetTranscriptHaplotypes(ctx, species, id, target, opts...)`

### 14. Variant Effect Predictor (VEP)
- `GetVariantConsequencesByHGVSNotation(ctx, species, hgvsNotation, target, opts...)`
- `GetVariantConsequencesByMultipleHGVSNotations(ctx, species, hgvsNotations, target, opts...)`
- `GetVariantConsequencesByID(ctx, species, id, target, opts...)`
- `GetVariantConsequencesByMultipleIDs(ctx, species, ids, target, opts...)`
- `GetVariantConsequencesByRegion(ctx, species, region, allele, target, opts...)`
- `GetVariantConsequencesByMultipleRegions(ctx, species, variants, target, opts...)`

### 15. Variation & Recoder
- `GetVariationRecoderByID(ctx, species, id, target, opts...)`
- `GetVariationRecoderByMultipleIDs(ctx, species, ids, target, opts...)`
- `GetVariationByID(ctx, species, id, target, opts...)`
- `GetVariationByPMCID(ctx, species, pmcid, target, opts...)`
- `GetVariationByPMID(ctx, species, pmid, target, opts...)`
- `GetVariationByMultipleIDs(ctx, species, ids, target, opts...)`

### 16. GA4GH
- `GetGA4GHBeacon(ctx, target, opts...)`
- `GetGA4GHBeaconQuery(ctx, alt, asm, ref, refName, start, target, opts...)`
- `PostGA4GHBeaconQuery(ctx, body, target, opts...)`
- `GetGA4GHFeaturesByID(ctx, id, target, opts...)`
- `SearchGA4GHFeatures(ctx, body, target, opts...)`
- `SearchGA4GHCallset(ctx, body, target, opts...)`
- `GetGA4GHCallsetByID(ctx, id, target, opts...)`
- `SearchGA4GHDatasets(ctx, body, target, opts...)`
- `GetGA4GHDatasetsByID(ctx, id, target, opts...)`
- `SearchGA4GHFeaturesets(ctx, body, target, opts...)`
- `GetGA4GHFeaturesetsByID(ctx, id, target, opts...)`
- `GetGA4GHVariantsByID(ctx, id, target, opts...)`
- `SearchGA4GHVariantAnnotations(ctx, body, target, opts...)`
- `SearchGA4GHVariants(ctx, body, target, opts...)`
- `SearchGA4GHVariantsets(ctx, body, target, opts...)`
- `GetGA4GHVariantsetsByID(ctx, id, target, opts...)`
- `SearchGA4GHReferences(ctx, body, target, opts...)`
- `GetGA4GHReferencesByID(ctx, id, target, opts...)`
- `SearchGA4GHReferencesets(ctx, body, target, opts...)`
- `GetGA4GHReferencesetsByID(ctx, id, target, opts...)`
- `SearchGA4GHVariantAnnotationsets(ctx, body, target, opts...)`
- `GetGA4GHVariantAnnotationsetsByID(ctx, id, target, opts...)`

---

## Dynamic Invocation & Introspection

You can introspect the endpoint catalog at runtime and call any endpoint dynamically by name:

```go
// Introspect all 100+ endpoints
endpoints := client.Endpoints()
for name, spec := range endpoints {
	fmt.Printf("%-35s %-6s %s\n", name, spec.Method, spec.URL)
}

// Call dynamically by name
var result map[string]any
err := client.Call(
	ctx,
	"getLookupById",
	map[string]string{"id": "ENSG00000157764"},
	nil,
	&result,
	ensembl.WithQuery("expand", "1"),
)
```

---

## Running Tests & CI

### Running Unit Tests

```bash
# Run unit tests
go test -v ./...

# Run unit tests with race detection
go test -v -race ./...

# Run unit tests and generate coverage profile
go test -v -race -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### Running Live Smoke Tests

```bash
# Execute smoke tests against live rest.ensembl.org API
ENSEMBL_LIVE_TESTS=1 go test -v -tags=live -race -run TestLiveEnsemblAPI ./...
```

---

## License

MIT License. Copyright (c) 2020-2026 Steve Moss. See [LICENSE](LICENSE) for details.
