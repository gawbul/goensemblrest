//go:build live

package goensemblrest_test

import (
	"context"
	"os"
	"testing"
	"time"

	ensembl "github.com/gawbul/goensemblrest"
)

func skipUnlessLive(t *testing.T) {
	t.Helper()
	if os.Getenv("ENSEMBL_LIVE_TESTS") == "" {
		t.Skip("Skipping live Ensembl REST test: ENSEMBL_LIVE_TESTS is not set")
	}
}

func TestLiveEnsemblAPI(t *testing.T) {
	skipUnlessLive(t)

	client, err := ensembl.NewClient(
		ensembl.WithTimeout(30*time.Second),
		ensembl.WithMaxAttempts(3),
		ensembl.WithRateLimit(10, time.Second),
		ensembl.WithUserAgent("goensemblrest-test/1.0 (+https://github.com/gawbul/goensemblrest)"),
	)
	if err != nil {
		t.Fatalf("failed to initialize live client: %v", err)
	}
	defer client.Close()

	t.Run("Ping", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
		defer cancel()

		var ping ensembl.PingResponse
		if err := client.GetInfoPing(ctx, &ping); err != nil {
			t.Fatalf("GetInfoPing failed: %v", err)
		}
		if ping.Ping != 1 {
			t.Errorf("expected ping: 1, got %d", ping.Ping)
		}
	})

	t.Run("Species", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
		defer cancel()

		var species ensembl.SpeciesResponse
		if err := client.GetInfoSpecies(ctx, &species); err != nil {
			t.Fatalf("GetInfoSpecies failed: %v", err)
		}
		if len(species.Species) == 0 {
			t.Errorf("expected species list, got 0")
		}
	})

	t.Run("Lookup", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
		defer cancel()

		var gene ensembl.LookupRecord
		if err := client.GetLookupByID(ctx, "ENSG00000157764", &gene); err != nil {
			t.Fatalf("GetLookupByID failed: %v", err)
		}
		if gene.ID != "ENSG00000157764" {
			t.Errorf("expected ID ENSG00000157764, got %q", gene.ID)
		}
		if gene.Species != "homo_sapiens" {
			t.Errorf("expected species homo_sapiens, got %q", gene.Species)
		}
	})

	t.Run("Sequence", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
		defer cancel()

		var seq ensembl.SequenceRecord
		if err := client.GetSequenceByID(ctx, "ENST00000288602", &seq); err != nil {
			t.Fatalf("GetSequenceByID failed: %v", err)
		}
		if len(seq.Seq) == 0 {
			t.Errorf("expected non-empty sequence")
		}
	})

	t.Run("Xrefs", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
		defer cancel()

		var xrefs []ensembl.XrefRecord
		if err := client.GetXrefsByID(ctx, "ENST00000288602", &xrefs); err != nil {
			t.Fatalf("GetXrefsByID failed: %v", err)
		}
		if len(xrefs) == 0 {
			t.Errorf("expected xref records, got 0")
		}
	})
}
