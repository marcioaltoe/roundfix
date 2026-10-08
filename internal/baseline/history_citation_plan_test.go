// Suite: Relocation Citations in Baseline Plans
// Invariant: citation impact is warning-only, digest-bound planning evidence that never changes apply mutations.
// Boundary IN: BuildPlan, Plan Digest assembly, and ApplyPlan over real Git repositories.
// Boundary OUT: citation parsing details and CLI rendering.

package baseline

import (
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

const (
	planCitationSource = "docs/adr/0040-retired.md"
	planCitationPath   = "docs/references/layout.md"
)

func TestBaselinePlanReportsRelocationCitations(t *testing.T) {
	t.Parallel()
	repository := newPlanRepository(t)
	writeInspectionFile(t, repository, planCitationSource, "---\nstatus: superseded\n---\n\n# Retired decision\n")
	writeInspectionFile(t, repository, planCitationPath, "See [the retired decision](../adr/0040-retired.md).\n")
	commitInspectionRepository(t, repository, "seed cited retired ADR")

	plan := buildTestPlan(t, repository)
	warning := planRelocationCitationWarning(t, plan)
	if warning.Path != planCitationPath ||
		!strings.Contains(warning.Message, "line 1 cites ../adr/0040-retired.md") {
		t.Fatalf("Relocation Citation warning = %#v, want citing path and line", warning)
	}

	current := newPlanRepository(t)
	writeInspectionFile(t, current, "docs/adr/0041-current.md", "---\nstatus: accepted\n---\n")
	writeInspectionFile(t, current, planCitationPath, "See [the current decision](../adr/0041-current.md).\n")
	commitInspectionRepository(t, current, "seed current ADR citation")
	currentPlan := buildTestPlan(t, current)
	for _, finding := range currentPlan.Warnings {
		if strings.HasPrefix(finding.Code, relocationCitationCode) {
			t.Fatalf("current-layout plan warning = %#v, want no Relocation Citation scan finding", finding)
		}
	}
}

func TestRelocationCitationsBindThePlanDigestOnly(t *testing.T) {
	t.Parallel()
	repository := newPlanRepository(t)
	writeInspectionFile(t, repository, planCitationSource, "---\nstatus: superseded\n---\n")
	writeInspectionFile(t, repository, planCitationPath, "docs/adr/0040-retired.md\n")
	commitInspectionRepository(t, repository, "seed citation-bound plan")

	withCitation := buildTestPlan(t, repository)
	writeInspectionFile(t, repository, planCitationPath, "No retired decision citation.\n")
	withoutCitation := buildTestPlan(t, repository)

	if withCitation.PlanDigest == withoutCitation.PlanDigest {
		t.Fatalf("citation-only edit left Plan Digest unchanged at %s", withCitation.PlanDigest)
	}
	planRelocationCitationWarning(t, withCitation)
	if warning := findPlanRelocationCitationWarning(withoutCitation); warning != nil {
		t.Fatalf("plan after citation removal still reports warning %#v", warning)
	}
	if !reflect.DeepEqual(withCitation.Preimages, withoutCitation.Preimages) ||
		!reflect.DeepEqual(withCitation.Postimages, withoutCitation.Postimages) ||
		!reflect.DeepEqual(withCitation.ManagedEntries, withoutCitation.ManagedEntries) ||
		!reflect.DeepEqual(withCitation.HistoryMoves, withoutCitation.HistoryMoves) ||
		!reflect.DeepEqual(withCitation.FileChanges, withoutCitation.FileChanges) {
		t.Fatalf(
			"citation-only edit changed mutation ledgers:\npreimages equal=%t\npostimages equal=%t\nmanaged entries equal=%t\nhistory moves equal=%t\nfile changes equal=%t",
			reflect.DeepEqual(withCitation.Preimages, withoutCitation.Preimages),
			reflect.DeepEqual(withCitation.Postimages, withoutCitation.Postimages),
			reflect.DeepEqual(withCitation.ManagedEntries, withoutCitation.ManagedEntries),
			reflect.DeepEqual(withCitation.HistoryMoves, withoutCitation.HistoryMoves),
			reflect.DeepEqual(withCitation.FileChanges, withoutCitation.FileChanges),
		)
	}
	assertPlanDoesNotMutateCitingFile(t, withCitation)
}

func TestRelocationCitationsLeaveApplyUnchanged(t *testing.T) {
	t.Parallel()
	withCitationRepository := newPlanRepository(t)
	writeInspectionFile(t, withCitationRepository, planCitationSource, "---\nstatus: superseded\n---\n")
	writeInspectionFile(t, withCitationRepository, planCitationPath, "docs/adr/0040-retired.md\n")
	commitInspectionRepository(t, withCitationRepository, "seed apply comparison")

	cloneParent := t.TempDir()
	withoutCitationRepository := filepath.Join(cloneParent, "repository")
	runInspectionCommand(t, cloneParent, "git", "clone", "--quiet", withCitationRepository, withoutCitationRepository)
	writeInspectionFile(t, withoutCitationRepository, planCitationPath, "No retired decision citation.\n")

	withCitationBytes, err := os.ReadFile(filepath.Join(withCitationRepository, filepath.FromSlash(planCitationPath)))
	if err != nil {
		t.Fatal(err)
	}
	withoutCitationBytes, err := os.ReadFile(filepath.Join(withoutCitationRepository, filepath.FromSlash(planCitationPath)))
	if err != nil {
		t.Fatal(err)
	}

	withCitationPlan := buildTestPlan(t, withCitationRepository)
	withoutCitationPlan := buildTestPlan(t, withoutCitationRepository)
	if _, err := ApplyPlan(t.Context(), withCitationRepository, withCitationPlan, withCitationPlan.PlanDigest); err != nil {
		t.Fatalf("ApplyPlan() with citation: %v", err)
	}
	if _, err := ApplyPlan(t.Context(), withoutCitationRepository, withoutCitationPlan, withoutCitationPlan.PlanDigest); err != nil {
		t.Fatalf("ApplyPlan() without citation: %v", err)
	}

	assertPlanCitationBytes(t, withCitationRepository, withCitationBytes)
	assertPlanCitationBytes(t, withoutCitationRepository, withoutCitationBytes)
	withCitationTree := visibleTreeWithoutPath(snapshotVisibleTree(t, withCitationRepository), planCitationPath)
	withoutCitationTree := visibleTreeWithoutPath(snapshotVisibleTree(t, withoutCitationRepository), planCitationPath)
	if !reflect.DeepEqual(withCitationTree, withoutCitationTree) {
		t.Fatalf("applied trees differ beyond the citing file:\nwith citation: %+v\nwithout citation: %+v", withCitationTree, withoutCitationTree)
	}
}

func findPlanRelocationCitationWarning(plan PlanDocument) *Finding {
	for index := range plan.Warnings {
		if plan.Warnings[index].Code == relocationCitationCode {
			return &plan.Warnings[index]
		}
	}
	return nil
}

func planRelocationCitationWarning(t *testing.T, plan PlanDocument) Finding {
	t.Helper()
	warning := findPlanRelocationCitationWarning(plan)
	if warning == nil {
		t.Fatalf("Plan warnings = %#v, want %s", plan.Warnings, relocationCitationCode)
	}
	return *warning
}

func assertPlanCitationBytes(t *testing.T, repository string, want []byte) {
	t.Helper()
	got, err := os.ReadFile(filepath.Join(repository, filepath.FromSlash(planCitationPath)))
	if err != nil {
		t.Fatalf("read citing file after apply: %v", err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("citing file after apply = %q, want %q", got, want)
	}
}

func assertPlanDoesNotMutateCitingFile(t *testing.T, plan PlanDocument) {
	t.Helper()
	for _, preimage := range plan.Preimages {
		if preimage.Path == planCitationPath {
			t.Fatalf("citing file entered Plan preimages: %#v", preimage)
		}
	}
	for _, postimage := range plan.Postimages {
		if postimage.Path == planCitationPath {
			t.Fatalf("citing file entered Plan postimages: %#v", postimage)
		}
	}
	for _, entry := range plan.ManagedEntries {
		if entry.Path == planCitationPath {
			t.Fatalf("citing file entered managed-entry ledger: %#v", entry)
		}
	}
	for _, change := range plan.FileChanges {
		if change.Path == planCitationPath {
			t.Fatalf("citing file entered file changes: %#v", change)
		}
	}
	for _, move := range plan.HistoryMoves {
		if move.From == planCitationPath || move.To == planCitationPath {
			t.Fatalf("citing file entered History Relocations: %#v", move)
		}
	}
}

func visibleTreeWithoutPath(entries []visibleTreeEntry, omitted string) []visibleTreeEntry {
	filtered := make([]visibleTreeEntry, 0, len(entries)-1)
	for _, entry := range entries {
		if entry.Path != omitted {
			filtered = append(filtered, entry)
		}
	}
	return filtered
}
