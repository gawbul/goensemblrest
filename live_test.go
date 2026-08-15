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
		ensembl.WithRateLimit(15, time.Second),
	)
	if err != nil {
		t.Fatalf("failed to initialize live client: %v", err)
	}
	defer client.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	t.Run("Ping", func(t *testing.T) {
		var ping ensembl.PingResponse
		if err := client.GetInfoPing(ctx, &ping); err != nil {
			t.Fatalf("GetInfoPing failed: %v", err)
		}
		if ping.Ping != 1 {
			t.Errorf("expected ping: 1, got %d", ping.Ping)
		}
	})

	t.Run("Archive", func(t *testing.T) {
		var archive ensembl.ArchiveRecord
		if err := client.GetArchiveByID(ctx, "ENSG00000157764", &archive); err != nil {
			t.Fatalf("GetArchiveByID failed: %v", err)
		}
		if archive.ID != "ENSG00000157764" {
			t.Errorf("expected ID ENSG00000157764, got %q", archive.ID)
		}
	})

	t.Run("Lookup", func(t *testing.T) {
		var gene ensembl.LookupRecord
		if err := client.GetLookupByID(ctx, "ENSG00000157764", &gene, ensembl.WithQuery("expand", "1")); err != nil {
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
		var seq ensembl.SequenceRecord
		if err := client.GetSequenceByID(ctx, "ENSG00000157764", &seq); err != nil {
			t.Fatalf("GetSequenceByID failed: %v", err)
		}
		if len(seq.Seq) == 0 {
			t.Errorf("expected non-empty sequence")
		}
	})

	t.Run("Xrefs", func(t *testing.T) {
		var xrefs []ensembl.XrefRecord
		if err := client.GetXrefsBySymbol(ctx, "homo_sapiens", "BRCA2", &xrefs); err != nil {
			t.Fatalf("GetXrefsBySymbol failed: %v", err)
		}
		if len(xrefs) == 0 {
			t.Errorf("expected xref records, got 0")
		}
	})
}
