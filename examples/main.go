package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	ensembl "github.com/gawbul/goensemblrest"
)

func main() {
	fmt.Println("Initializing Ensembl REST client...")

	// 1. Create a client with custom options
	client, err := ensembl.NewClient(
		ensembl.WithTimeout(60*time.Second),
		ensembl.WithMaxAttempts(5),
		ensembl.WithRateLimit(15, time.Second),
	)
	if err != nil {
		log.Fatalf("Failed to initialize client: %v", err)
	}
	defer client.Close()

	ctx := context.Background()

	// Print client user agent & base URL
	fmt.Printf("User-Agent: %s\n", client.UserAgent())
	fmt.Printf("Base URL:   %s\n\n", client.BaseURL())

	// 2. Service Ping
	fmt.Println("--> Checking service health via /info/ping...")
	var ping ensembl.PingResponse
	if err := client.GetInfoPing(ctx, &ping); err != nil {
		log.Printf("Ping failed: %v\n", err)
	} else {
		fmt.Printf("Ping status: %d (Healthy)\n\n", ping.Ping)
	}

	// 3. Archive endpoint
	fmt.Println("--> Fetching Archive for ENSG00000157764 (BRAF)...")
	var archive ensembl.ArchiveRecord
	if err := client.GetArchiveByID(ctx, "ENSG00000157764", &archive); err != nil {
		log.Printf("Archive fetch error: %v\n", err)
	} else {
		fmt.Printf("Archive ID: %s, Latest: %s, Version: %d\n\n", archive.ID, archive.Latest, archive.Version)
	}

	// 4. Lookup endpoint with query options
	fmt.Println("--> Looking up Gene BRCA2 with expansion...")
	var lookup ensembl.LookupRecord
	if err := client.GetLookupBySymbol(ctx, "homo_sapiens", "BRCA2", &lookup, ensembl.WithQuery("expand", "1")); err != nil {
		log.Printf("Lookup error: %v\n", err)
	} else {
		fmt.Printf("Gene: %s (%s), Biotype: %s, Location: %s:%d-%d (Strand: %d)\n",
			lookup.DisplayName, lookup.ID, lookup.Biotype, lookup.SeqRegion, lookup.Start, lookup.End, lookup.Strand)
		fmt.Printf("Number of Transcripts: %d\n\n", len(lookup.Transcripts))
	}

	// 5. Cross references
	fmt.Println("--> Fetching external cross-references for BRCA2...")
	var xrefs []ensembl.XrefRecord
	if err := client.GetXrefsBySymbol(ctx, "homo_sapiens", "BRCA2", &xrefs); err != nil {
		log.Printf("Xrefs error: %v\n", err)
	} else {
		fmt.Printf("Found %d external references:\n", len(xrefs))
		limit := 5
		if len(xrefs) < limit {
			limit = len(xrefs)
		}
		for i := 0; i < limit; i++ {
			fmt.Printf("  - DB: %-15s PrimaryID: %-15s DisplayID: %s\n", xrefs[i].DBName, xrefs[i].PrimaryID, xrefs[i].DisplayID)
		}
		fmt.Println()
	}

	// 6. Sequence retrieval (FASTA format via header option)
	fmt.Println("--> Fetching Sequence in FASTA format...")
	var fastaSeq string
	if err := client.GetSequenceByID(ctx, "ENSG00000157764", &fastaSeq, ensembl.WithRequestContentType("text/x-fasta")); err != nil {
		log.Printf("FASTA sequence fetch error: %v\n", err)
	} else {
		lines := 3
		fmt.Printf("FASTA sequence preview:\n%s\n...\n\n", truncateLines(fastaSeq, lines))
	}

	// 7. Error Handling Demonstration
	fmt.Println("--> Demonstrating structured error handling on invalid query...")
	var errTarget map[string]any
	err = client.GetArchiveByID(ctx, "INVALID_QUERY_ID", &errTarget)
	if err != nil {
		if errors.Is(err, ensembl.ErrBadRequest) {
			fmt.Println("Successfully caught ErrBadRequest using errors.Is!")
		}
		var apiErr *ensembl.APIError
		if errors.As(err, &apiErr) {
			fmt.Printf("Structured API Error details:\n  Status Code: %d\n  Server Message: %s\n",
				apiErr.StatusCode, apiErr.Message)
		}
	}
	fmt.Println()

	// 8. Rate limit inspection
	rateInfo := client.RateLimit()
	fmt.Println("--> Current rate limit telemetry from headers:")
	rawJSON, _ := json.MarshalIndent(rateInfo, "", "  ")
	fmt.Println(string(rawJSON))

	fmt.Println("\nAll examples completed successfully!")
}

func truncateLines(s string, maxLines int) string {
	var out []rune
	lines := 0
	for _, r := range s {
		out = append(out, r)
		if r == '\n' {
			lines++
			if lines >= maxLines {
				break
			}
		}
	}
	return string(out)
}
