package goensemblrest

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestAllDomainEndpoints(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]any{
			"url":    r.URL.String(),
			"method": r.Method,
			"status": "success",
			"ping":   1,
		})
	}))
	defer server.Close()

	client, err := NewClient(
		WithBaseURL(server.URL),
		WithMaxAttempts(1),
		WithRateLimit(10000, time.Second),
	)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}
	defer client.Close()

	ctx := context.Background()
	var res map[string]any

	// Archive
	t.Run("Archive", func(t *testing.T) {
		if err := client.GetArchiveByID(ctx, "ENSG00000157764", &res); err != nil {
			t.Errorf("GetArchiveByID failed: %v", err)
		}
		if err := client.GetArchiveByMultipleIDs(ctx, []string{"ENSG1", "ENSG2"}, &res); err != nil {
			t.Errorf("GetArchiveByMultipleIDs failed: %v", err)
		}
	})

	// Comparative Genomics
	t.Run("Comparative", func(t *testing.T) {
		if err := client.GetCafeGeneTreeByID(ctx, "ENSGT00390000003602", &res); err != nil {
			t.Errorf("GetCafeGeneTreeByID failed: %v", err)
		}
		if err := client.GetCafeGeneTreeMemberBySymbol(ctx, "human", "BRCA2", &res); err != nil {
			t.Errorf("GetCafeGeneTreeMemberBySymbol failed: %v", err)
		}
		if err := client.GetCafeGeneTreeMemberByID(ctx, "human", "ENSG167664", &res); err != nil {
			t.Errorf("GetCafeGeneTreeMemberByID failed: %v", err)
		}
		if err := client.GetGeneTreeByID(ctx, "ENSGT00390000003602", &res); err != nil {
			t.Errorf("GetGeneTreeByID failed: %v", err)
		}
		if err := client.GetGeneTreeMemberBySymbol(ctx, "human", "BRCA2", &res); err != nil {
			t.Errorf("GetGeneTreeMemberBySymbol failed: %v", err)
		}
		if err := client.GetGeneTreeMemberByID(ctx, "human", "ENSG167664", &res); err != nil {
			t.Errorf("GetGeneTreeMemberByID failed: %v", err)
		}
		if err := client.GetAlignmentByRegion(ctx, "human", "X:1000..2000:1", &res); err != nil {
			t.Errorf("GetAlignmentByRegion failed: %v", err)
		}
		if err := client.GetHomologyByID(ctx, "human", "ENSG157764", &res); err != nil {
			t.Errorf("GetHomologyByID failed: %v", err)
		}
		if err := client.GetHomologyBySymbol(ctx, "human", "BRCA2", &res); err != nil {
			t.Errorf("GetHomologyBySymbol failed: %v", err)
		}
	})

	// Cross References
	t.Run("Xrefs", func(t *testing.T) {
		if err := client.GetXrefsBySymbol(ctx, "human", "BRCA2", &res); err != nil {
			t.Errorf("GetXrefsBySymbol failed: %v", err)
		}
		if err := client.GetXrefsByID(ctx, "ENSG157764", &res); err != nil {
			t.Errorf("GetXrefsByID failed: %v", err)
		}
		if err := client.GetXrefsByName(ctx, "human", "BRCA2", &res); err != nil {
			t.Errorf("GetXrefsByName failed: %v", err)
		}
	})

	// Info
	t.Run("Info", func(t *testing.T) {
		if err := client.GetInfoAnalysis(ctx, "human", &res); err != nil {
			t.Errorf("GetInfoAnalysis failed: %v", err)
		}
		if err := client.GetInfoAssembly(ctx, "human", &res, WithQuery("bands", "1")); err != nil {
			t.Errorf("GetInfoAssembly failed: %v", err)
		}
		if err := client.GetInfoAssemblyRegion(ctx, "human", "X", &res); err != nil {
			t.Errorf("GetInfoAssemblyRegion failed: %v", err)
		}
		if err := client.GetInfoBiotypes(ctx, "human", &res); err != nil {
			t.Errorf("GetInfoBiotypes failed: %v", err)
		}
		if err := client.GetInfoBiotypesByGroup(ctx, "coding", "gene", &res); err != nil {
			t.Errorf("GetInfoBiotypesByGroup failed: %v", err)
		}
		if err := client.GetInfoBiotypesByName(ctx, "protein_coding", "gene", &res); err != nil {
			t.Errorf("GetInfoBiotypesByName failed: %v", err)
		}
		if err := client.GetInfoComparaMethods(ctx, &res); err != nil {
			t.Errorf("GetInfoComparaMethods failed: %v", err)
		}
		if err := client.GetInfoComparaSpeciesSets(ctx, "EPO", &res); err != nil {
			t.Errorf("GetInfoComparaSpeciesSets failed: %v", err)
		}
		if err := client.GetInfoComparas(ctx, &res); err != nil {
			t.Errorf("GetInfoComparas failed: %v", err)
		}
		if err := client.GetInfoData(ctx, &res); err != nil {
			t.Errorf("GetInfoData failed: %v", err)
		}
		if err := client.GetInfoEgVersion(ctx, &res); err != nil {
			t.Errorf("GetInfoEgVersion failed: %v", err)
		}
		if err := client.GetInfoExternalDbs(ctx, "human", &res); err != nil {
			t.Errorf("GetInfoExternalDbs failed: %v", err)
		}
		if err := client.GetInfoDivisions(ctx, &res); err != nil {
			t.Errorf("GetInfoDivisions failed: %v", err)
		}
		if err := client.GetInfoGenomesByName(ctx, "arabidopsis_thaliana", &res); err != nil {
			t.Errorf("GetInfoGenomesByName failed: %v", err)
		}
		if err := client.GetInfoGenomesByAccession(ctx, "U00096", &res); err != nil {
			t.Errorf("GetInfoGenomesByAccession failed: %v", err)
		}
		if err := client.GetInfoGenomesByAssembly(ctx, "GCA_902167145.1", &res); err != nil {
			t.Errorf("GetInfoGenomesByAssembly failed: %v", err)
		}
		if err := client.GetInfoGenomesByDivision(ctx, "EnsemblPlants", &res); err != nil {
			t.Errorf("GetInfoGenomesByDivision failed: %v", err)
		}
		if err := client.GetInfoGenomesByTaxonomy(ctx, "Homo sapiens", &res); err != nil {
			t.Errorf("GetInfoGenomesByTaxonomy failed: %v", err)
		}
		if err := client.GetInfoPing(ctx, &res); err != nil {
			t.Errorf("GetInfoPing failed: %v", err)
		}
		if err := client.GetInfoRest(ctx, &res); err != nil {
			t.Errorf("GetInfoRest failed: %v", err)
		}
		if err := client.GetInfoSoftware(ctx, &res); err != nil {
			t.Errorf("GetInfoSoftware failed: %v", err)
		}
		if err := client.GetInfoSpecies(ctx, &res); err != nil {
			t.Errorf("GetInfoSpecies failed: %v", err)
		}
		if err := client.GetInfoVariationBySpecies(ctx, "human", &res); err != nil {
			t.Errorf("GetInfoVariationBySpecies failed: %v", err)
		}
		if err := client.GetInfoVariationConsequenceTypes(ctx, &res); err != nil {
			t.Errorf("GetInfoVariationConsequenceTypes failed: %v", err)
		}
		if err := client.GetInfoVariationPopulationIndividuals(ctx, "human", "1000GENOMES:phase_3:ASW", &res); err != nil {
			t.Errorf("GetInfoVariationPopulationIndividuals failed: %v", err)
		}
		if err := client.GetInfoVariationPopulations(ctx, "human", &res); err != nil {
			t.Errorf("GetInfoVariationPopulations failed: %v", err)
		}
	})

	// Linkage Disequilibrium
	t.Run("LD", func(t *testing.T) {
		if err := client.GetLdID(ctx, "human", "rs56116432", "1000GENOMES:phase_3:KHV", &res); err != nil {
			t.Errorf("GetLdID failed: %v", err)
		}
		if err := client.GetLdPairwise(ctx, "human", "rs6792369", "rs1042779", &res); err != nil {
			t.Errorf("GetLdPairwise failed: %v", err)
		}
		if err := client.GetLdRegion(ctx, "human", "6:25837556..25843455", "1000GENOMES:phase_3:KHV", &res); err != nil {
			t.Errorf("GetLdRegion failed: %v", err)
		}
	})

	// Lookup
	t.Run("Lookup", func(t *testing.T) {
		if err := client.GetLookupByID(ctx, "ENSG00000157764", &res); err != nil {
			t.Errorf("GetLookupByID failed: %v", err)
		}
		if err := client.GetLookupByMultipleIDs(ctx, []string{"ENSG1", "ENSG2"}, &res); err != nil {
			t.Errorf("GetLookupByMultipleIDs failed: %v", err)
		}
		if err := client.GetLookupBySymbol(ctx, "human", "BRCA2", &res); err != nil {
			t.Errorf("GetLookupBySymbol failed: %v", err)
		}
		if err := client.GetLookupByMultipleSymbols(ctx, "human", []string{"BRCA1", "BRCA2"}, &res); err != nil {
			t.Errorf("GetLookupByMultipleSymbols failed: %v", err)
		}
	})

	// Mapping
	t.Run("Mapping", func(t *testing.T) {
		if err := client.GetMapCdnaToRegion(ctx, "ENST00000288602", "100..300", &res); err != nil {
			t.Errorf("GetMapCdnaToRegion failed: %v", err)
		}
		if err := client.GetMapCdsToRegion(ctx, "ENST00000288602", "1..100", &res); err != nil {
			t.Errorf("GetMapCdsToRegion failed: %v", err)
		}
		if err := client.GetMapAssemblyOneToTwo(ctx, "human", "GRCh37", "X:1000000..1000100:1", "GRCh38", &res); err != nil {
			t.Errorf("GetMapAssemblyOneToTwo failed: %v", err)
		}
		if err := client.GetMapTranslationToRegion(ctx, "ENSP00000288602", "100..300", &res); err != nil {
			t.Errorf("GetMapTranslationToRegion failed: %v", err)
		}
	})

	// Ontologies
	t.Run("Ontologies", func(t *testing.T) {
		if err := client.GetAncestorsByID(ctx, "GO:0005667", &res); err != nil {
			t.Errorf("GetAncestorsByID failed: %v", err)
		}
		if err := client.GetAncestorsChartByID(ctx, "GO:0005667", &res); err != nil {
			t.Errorf("GetAncestorsChartByID failed: %v", err)
		}
		if err := client.GetDescendantsByID(ctx, "GO:0005667", &res); err != nil {
			t.Errorf("GetDescendantsByID failed: %v", err)
		}
		if err := client.GetOntologyByID(ctx, "GO:0005667", &res); err != nil {
			t.Errorf("GetOntologyByID failed: %v", err)
		}
		if err := client.GetOntologyByName(ctx, "transcription factor complex", &res); err != nil {
			t.Errorf("GetOntologyByName failed: %v", err)
		}
		if err := client.GetTaxonomyClassificationByID(ctx, "9606", &res); err != nil {
			t.Errorf("GetTaxonomyClassificationByID failed: %v", err)
		}
		if err := client.GetTaxonomyByID(ctx, "9606", &res); err != nil {
			t.Errorf("GetTaxonomyByID failed: %v", err)
		}
		if err := client.GetTaxonomyByName(ctx, "human", &res); err != nil {
			t.Errorf("GetTaxonomyByName failed: %v", err)
		}
	})

	// Overlap
	t.Run("Overlap", func(t *testing.T) {
		if err := client.GetOverlapByID(ctx, "ENSG00000157764", &res); err != nil {
			t.Errorf("GetOverlapByID failed: %v", err)
		}
		if err := client.GetOverlapByRegion(ctx, "human", "7:140424943-140624564", &res); err != nil {
			t.Errorf("GetOverlapByRegion failed: %v", err)
		}
		if err := client.GetOverlapByTranslation(ctx, "ENSP00000288602", &res); err != nil {
			t.Errorf("GetOverlapByTranslation failed: %v", err)
		}
	})

	// Phenotype
	t.Run("Phenotype", func(t *testing.T) {
		if err := client.GetPhenotypeByAccession(ctx, "human", "EFO:0003900", &res); err != nil {
			t.Errorf("GetPhenotypeByAccession failed: %v", err)
		}
		if err := client.GetPhenotypeByGene(ctx, "human", "ENSG00000157764", &res); err != nil {
			t.Errorf("GetPhenotypeByGene failed: %v", err)
		}
		if err := client.GetPhenotypeByRegion(ctx, "human", "7:140424943-140624564", &res); err != nil {
			t.Errorf("GetPhenotypeByRegion failed: %v", err)
		}
		if err := client.GetPhenotypeByTerm(ctx, "human", "diabetes", &res); err != nil {
			t.Errorf("GetPhenotypeByTerm failed: %v", err)
		}
	})

	// Regulation
	t.Run("Regulation", func(t *testing.T) {
		if err := client.GetRegulationBindingMatrix(ctx, "human", "MA0139.1", &res); err != nil {
			t.Errorf("GetRegulationBindingMatrix failed: %v", err)
		}
	})

	// Sequence
	t.Run("Sequence", func(t *testing.T) {
		if err := client.GetSequenceByID(ctx, "ENSG00000157764", &res); err != nil {
			t.Errorf("GetSequenceByID failed: %v", err)
		}
		if err := client.GetSequenceByMultipleIDs(ctx, []string{"ENSG1", "ENSG2"}, &res); err != nil {
			t.Errorf("GetSequenceByMultipleIDs failed: %v", err)
		}
		if err := client.GetSequenceByRegion(ctx, "human", "X:1000..2000:1", &res); err != nil {
			t.Errorf("GetSequenceByRegion failed: %v", err)
		}
		if err := client.GetSequenceByMultipleRegions(ctx, "human", []string{"X:1000..2000:1"}, &res); err != nil {
			t.Errorf("GetSequenceByMultipleRegions failed: %v", err)
		}
	})

	// Transcript
	t.Run("Transcript", func(t *testing.T) {
		if err := client.GetTranscriptHaplotypes(ctx, "human", "ENST00000288602", &res); err != nil {
			t.Errorf("GetTranscriptHaplotypes failed: %v", err)
		}
	})

	// VEP
	t.Run("VEP", func(t *testing.T) {
		if err := client.GetVariantConsequencesByHGVSNotation(ctx, "human", "AGT:c.803T>C", &res); err != nil {
			t.Errorf("GetVariantConsequencesByHGVSNotation failed: %v", err)
		}
		if err := client.GetVariantConsequencesByMultipleHGVSNotations(ctx, "human", []string{"AGT:c.803T>C"}, &res); err != nil {
			t.Errorf("GetVariantConsequencesByMultipleHGVSNotations failed: %v", err)
		}
		if err := client.GetVariantConsequencesByID(ctx, "human", "COSM476", &res); err != nil {
			t.Errorf("GetVariantConsequencesByID failed: %v", err)
		}
		if err := client.GetVariantConsequencesByMultipleIDs(ctx, "human", []string{"COSM476"}, &res); err != nil {
			t.Errorf("GetVariantConsequencesByMultipleIDs failed: %v", err)
		}
		if err := client.GetVariantConsequencesByRegion(ctx, "human", "9:22125503-22125502:1", "C", &res); err != nil {
			t.Errorf("GetVariantConsequencesByRegion failed: %v", err)
		}
		if err := client.GetVariantConsequencesByMultipleRegions(ctx, "human", []string{"21 26978314 rs116645811 G A . . ."}, &res); err != nil {
			t.Errorf("GetVariantConsequencesByMultipleRegions failed: %v", err)
		}
	})

	// Variation
	t.Run("Variation", func(t *testing.T) {
		if err := client.GetVariationRecoderByID(ctx, "human", "rs699", &res); err != nil {
			t.Errorf("GetVariationRecoderByID failed: %v", err)
		}
		if err := client.GetVariationRecoderByMultipleIDs(ctx, "human", []string{"rs699"}, &res); err != nil {
			t.Errorf("GetVariationRecoderByMultipleIDs failed: %v", err)
		}
		if err := client.GetVariationByID(ctx, "human", "rs699", &res); err != nil {
			t.Errorf("GetVariationByID failed: %v", err)
		}
		if err := client.GetVariationByPMCID(ctx, "human", "PMC2176140", &res); err != nil {
			t.Errorf("GetVariationByPMCID failed: %v", err)
		}
		if err := client.GetVariationByPMID(ctx, "human", "10072433", &res); err != nil {
			t.Errorf("GetVariationByPMID failed: %v", err)
		}
		if err := client.GetVariationByMultipleIDs(ctx, "human", []string{"rs699"}, &res); err != nil {
			t.Errorf("GetVariationByMultipleIDs failed: %v", err)
		}
	})

	// GA4GH
	t.Run("GA4GH", func(t *testing.T) {
		if err := client.GetGA4GHBeacon(ctx, &res); err != nil {
			t.Errorf("GetGA4GHBeacon failed: %v", err)
		}
		if err := client.GetGA4GHBeaconQuery(ctx, "C", "GRCh38", "T", "1", 1000, &res); err != nil {
			t.Errorf("GetGA4GHBeaconQuery failed: %v", err)
		}
		if err := client.PostGA4GHBeaconQuery(ctx, map[string]any{"assemblyId": "GRCh38"}, &res); err != nil {
			t.Errorf("PostGA4GHBeaconQuery failed: %v", err)
		}
		if err := client.GetGA4GHFeaturesByID(ctx, "feat1", &res); err != nil {
			t.Errorf("GetGA4GHFeaturesByID failed: %v", err)
		}
		if err := client.SearchGA4GHFeatures(ctx, map[string]any{"referenceName": "1"}, &res); err != nil {
			t.Errorf("SearchGA4GHFeatures failed: %v", err)
		}
		if err := client.SearchGA4GHCallset(ctx, map[string]any{"variantSetId": "vs1"}, &res); err != nil {
			t.Errorf("SearchGA4GHCallset failed: %v", err)
		}
		if err := client.GetGA4GHCallsetByID(ctx, "cs1", &res); err != nil {
			t.Errorf("GetGA4GHCallsetByID failed: %v", err)
		}
		if err := client.SearchGA4GHDatasets(ctx, map[string]any{"pageSize": 10}, &res); err != nil {
			t.Errorf("SearchGA4GHDatasets failed: %v", err)
		}
		if err := client.GetGA4GHDatasetsByID(ctx, "ds1", &res); err != nil {
			t.Errorf("GetGA4GHDatasetsByID failed: %v", err)
		}
		if err := client.SearchGA4GHFeaturesets(ctx, map[string]any{"datasetId": "ds1"}, &res); err != nil {
			t.Errorf("SearchGA4GHFeaturesets failed: %v", err)
		}
		if err := client.GetGA4GHFeaturesetsByID(ctx, "fs1", &res); err != nil {
			t.Errorf("GetGA4GHFeaturesetsByID failed: %v", err)
		}
		if err := client.GetGA4GHVariantsByID(ctx, "v1", &res); err != nil {
			t.Errorf("GetGA4GHVariantsByID failed: %v", err)
		}
		if err := client.SearchGA4GHVariantAnnotations(ctx, map[string]any{"variantAnnotationSetId": "vas1"}, &res); err != nil {
			t.Errorf("SearchGA4GHVariantAnnotations failed: %v", err)
		}
		if err := client.SearchGA4GHVariants(ctx, map[string]any{"variantSetId": "vs1"}, &res); err != nil {
			t.Errorf("SearchGA4GHVariants failed: %v", err)
		}
		if err := client.SearchGA4GHVariantsets(ctx, map[string]any{"datasetId": "ds1"}, &res); err != nil {
			t.Errorf("SearchGA4GHVariantsets failed: %v", err)
		}
		if err := client.GetGA4GHVariantsetsByID(ctx, "vs1", &res); err != nil {
			t.Errorf("GetGA4GHVariantsetsByID failed: %v", err)
		}
		if err := client.SearchGA4GHReferences(ctx, map[string]any{"referenceSetId": "rs1"}, &res); err != nil {
			t.Errorf("SearchGA4GHReferences failed: %v", err)
		}
		if err := client.GetGA4GHReferencesByID(ctx, "ref1", &res); err != nil {
			t.Errorf("GetGA4GHReferencesByID failed: %v", err)
		}
		if err := client.SearchGA4GHReferencesets(ctx, map[string]any{"pageSize": 10}, &res); err != nil {
			t.Errorf("SearchGA4GHReferencesets failed: %v", err)
		}
		if err := client.GetGA4GHReferencesetsByID(ctx, "rs1", &res); err != nil {
			t.Errorf("GetGA4GHReferencesetsByID failed: %v", err)
		}
		if err := client.SearchGA4GHVariantAnnotationsets(ctx, map[string]any{"variantSetId": "vs1"}, &res); err != nil {
			t.Errorf("SearchGA4GHVariantAnnotationsets failed: %v", err)
		}
		if err := client.GetGA4GHVariantAnnotationsetsByID(ctx, "vas1", &res); err != nil {
			t.Errorf("GetGA4GHVariantAnnotationsetsByID failed: %v", err)
		}
	})
}
