// Suite: Root Backup reporting.
// Invariant: adoption names exact backup bytes and managed refresh adds none.
// Boundary IN: plan, apply and update commands in temporary Git repositories.
// Boundary OUT: transaction fault injection.
package cli

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"roundfix/internal/baseline"
)

func rootBackupTestPlan(t *testing.T) (string, string, string, string) {
	t.Helper()
	repo := newBaselineApplyTestRepository(t)
	const original = "# Local instructions\nKeep this exact copy.\n"
	writeBaselinePlanTestFile(t, repo, "AGENTS.md", original)
	commitBaselinePlanTestRepository(t, repo)
	plan, planPath := baselineApplyTestPlan(t, repo)
	backup := ""
	for _, entry := range plan.ManagedEntries {
		if entry.ID == "backup:AGENTS.md" {
			backup = entry.Path
		}
	}
	if backup == "" {
		t.Fatal("plan omitted AGENTS.md backup")
	}
	return repo, planPath, plan.PlanDigest, backup
}

func TestBaselineApplyNamesTheRootBackupItWrote(t *testing.T) {
	t.Parallel()
	for _, format := range []string{"text", "json"} {
		t.Run(format, func(t *testing.T) {
			repo, planPath, digest, backup := rootBackupTestPlan(t)
			original, err := os.ReadFile(filepath.Join(repo, "AGENTS.md"))
			if err != nil {
				t.Fatal(err)
			}
			var out, diagnostics bytes.Buffer
			code := RunContext(t.Context(), []string{"baseline", "apply", "--repo", repo, "--plan", planPath, "--confirm-plan", digest, "--format", format}, &out, &diagnostics)
			if code != exitOK {
				t.Fatalf("exit=%d stdout=%s stderr=%s", code, &out, &diagnostics)
			}
			content, err := os.ReadFile(filepath.Join(repo, backup))
			if err != nil || !bytes.Equal(content, original) {
				t.Fatalf("backup=%q err=%v", content, err)
			}
			message := "exact copy of AGENTS.md before this apply; commit it with this change, Roundfix never deletes it"
			if format == "text" {
				if !strings.Contains(out.String(), "Warning: baseline.preservation.backup.written: "+backup+": "+message+"\n") {
					t.Fatal(out.String())
				}
			} else {
				result, err := baseline.ParseResult(out.Bytes())
				if err != nil {
					t.Fatal(err)
				}
				count := 0
				for _, warning := range result.Warnings {
					if warning.Code == "baseline.preservation.backup.written" {
						count++
						if warning.Path != backup || warning.Message != message {
							t.Fatalf("warning=%+v", warning)
						}
					}
				}
				if count != 1 {
					t.Fatalf("backup warnings=%d", count)
				}
			}
		})
	}
}

func TestBaselinePlanTextLabelsTheRootBackup(t *testing.T) {
	t.Parallel()
	_, planPath, _, backup := rootBackupTestPlan(t)
	data, err := os.ReadFile(planPath)
	if err != nil {
		t.Fatal(err)
	}
	plan, err := baseline.ParsePlanDocument(data)
	if err != nil {
		t.Fatal(err)
	}
	var out bytes.Buffer
	printBaselinePlanText(plan, &out)
	for _, change := range plan.FileChanges {
		if change.Path == backup {
			want := fmt.Sprintf("- create %s (%d managed entries; root backup of AGENTS.md)\n", backup, len(change.ManagedEntries))
			if !strings.Contains(out.String(), want) {
				t.Fatalf("missing %q: %s", want, &out)
			}
			return
		}
	}
	t.Fatal("backup file change absent")
}

func TestBaselineUpdateWritesNoRootBackupWarning(t *testing.T) {
	t.Parallel()
	repo := staleBaselineUpdateRepository(t)
	var out, diagnostics bytes.Buffer
	code := RunContext(t.Context(), []string{"baseline", "update", "--repo", repo, "--yes", "--no-skills", "--no-history", "--format=json"}, &out, &diagnostics)
	if code != exitOK {
		t.Fatalf("exit=%d stdout=%s stderr=%s", code, &out, &diagnostics)
	}
	if strings.Contains(out.String(), "baseline.preservation.backup.written") {
		t.Fatal(out.String())
	}
	var result baselineUpdateResult
	if err := json.Unmarshal(out.Bytes(), &result); err != nil {
		t.Fatal(err)
	}
	if result.State != "verified" || result.ApprovedPlanDigest == "" || result.ApprovedPlanDigest != result.PlanDigest {
		t.Fatalf("update did not apply: %+v", result)
	}
	// Inspect the apply result too: update has its own warning projection.
	catalog, err := baseline.LoadEmbeddedCatalog()
	if err != nil {
		t.Fatal(err)
	}
	input, err := baseline.ResolveManifestInput(repo, catalog)
	if err != nil {
		t.Fatal(err)
	}
	outcome, err := baseline.BuildPlan(t.Context(), baseline.PlanRequest{
		Repository: repo, ProfileID: input.ProfileID, Decisions: input.Decisions,
		Preservation: baseline.RootPreservationRequest{Mode: baseline.PreservationModeManagedRefresh},
	})
	if err != nil || outcome.Plan == nil {
		t.Fatalf("managed refresh plan: %v %+v", err, outcome.Result)
	}
	applied, err := baseline.ApplyPlan(t.Context(), repo, *outcome.Plan, outcome.Plan.PlanDigest)
	if err != nil {
		t.Fatal(err)
	}
	for _, warning := range applied.Warnings {
		if warning.Code == "baseline.preservation.backup.written" {
			t.Fatalf("managed refresh backup warning: %+v", warning)
		}
	}
}
