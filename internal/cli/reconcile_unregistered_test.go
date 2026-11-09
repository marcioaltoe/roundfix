package cli

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"roundfix/internal/store"
)

func TestReconcileApplyReleasesAnIgnoredOnlyRunWorktreeAndEmptyRepositoryDirectories(t *testing.T) {
	t.Parallel()
	run, root := newReconcileUnregisteredIgnoredFixture(t)
	other := filepath.Join(root, "other-0123abcd")
	notes := filepath.Join(root, "notes")
	for _, path := range []string{other, notes} {
		if err := os.Mkdir(path, 0o755); err != nil {
			t.Fatal(err)
		}
	}
	text := string(itemReconcileOutput(t))
	for _, want := range []string{
		"  classification: safe\n",
		"  evidence: recorded Run Worktree is not registered in the Git root and holds only files the repository ignores\n",
		"Empty repository directory: " + other + "\n  action: would remove with --apply\n",
	} {
		if !strings.Contains(text, want) {
			t.Fatalf("dry run lacks %q:\n%s", want, text)
		}
	}
	dry := itemReconcileReport(t)
	if len(dry.EmptyDirectoryCandidates) != 1 {
		t.Fatalf("dry candidates = %+v", dry.EmptyDirectoryCandidates)
	}
	candidate := dry.EmptyDirectoryCandidates[0]
	if candidate.Kind != "emptyRepositoryDirectory" || candidate.Worktree != other || candidate.Proof != emptyRepositoryDirectoryProof || candidate.Action != "would remove with --apply" {
		t.Fatalf("dry candidate = %+v", candidate)
	}
	assertReconcilePathState(t, run.WorkDir, true)
	assertReconcilePathState(t, other, true)
	applied := itemReconcileReport(t, "--apply")
	if applied.Summary.Applied != 1 || applied.DebrisSummary.EmptyDirectoriesRemoved != 2 || len(applied.EmptyDirectoryCandidates) != 2 {
		t.Fatalf("apply = %+v", applied)
	}
	for _, candidate := range applied.EmptyDirectoryCandidates {
		if candidate.Action != "removed" || candidate.RefusalReason != "" {
			t.Fatalf("apply candidate = %+v", candidate)
		}
	}
	assertReconcilePathState(t, run.WorkDir, false)
	assertReconcilePathState(t, filepath.Dir(run.WorkDir), false)
	assertReconcilePathState(t, other, false)
	assertReconcilePathState(t, notes, true)
	next := itemReconcileReport(t)
	if len(next.Results) != 1 || next.Results[0].Classification != "released" {
		t.Fatalf("next results = %+v", next.Results)
	}
	// The additive collections and counter stay absent without candidates.
	raw := itemReconcileOutput(t, "--format=json")
	for _, field := range []string{"emptyDirectoryCandidates", "emptyDirectoriesRemoved"} {
		if bytes.Contains(raw, []byte(field)) {
			t.Fatalf("empty report acquired %s: %s", field, raw)
		}
	}
}

func TestReconcileWithARunIDRemovesNoEmptyRepositoryDirectory(t *testing.T) {
	t.Parallel()
	run, root := newReconcileUnregisteredIgnoredFixture(t)
	other := filepath.Join(root, "other-0123abcd")
	if err := os.Mkdir(other, 0o755); err != nil {
		t.Fatal(err)
	}
	for _, args := range [][]string{{run.ID}, {run.ID, "--apply"}} {
		report := itemReconcileReport(t, args...)
		if len(report.EmptyDirectoryCandidates) != 0 || report.DebrisSummary.EmptyDirectoriesRemoved != 0 {
			t.Fatalf("Run-ID report swept directories: %+v", report)
		}
		assertReconcilePathState(t, other, true)
	}
	assertReconcilePathState(t, run.WorkDir, false)
	assertReconcilePathState(t, filepath.Dir(run.WorkDir), true)
	raw := itemReconcileOutput(t, run.ID, "--format=json")
	var envelope map[string]json.RawMessage
	if err := json.Unmarshal(raw, &envelope); err != nil {
		t.Fatal(err)
	}
	if _, found := envelope["emptyDirectoryCandidates"]; found {
		t.Fatalf("Run-ID report lists empty directories: %s", raw)
	}
	if text := string(itemReconcileOutput(t, run.ID)); strings.Contains(text, "Empty repository directory:") {
		t.Fatalf("Run-ID text lists empty directories: %s", text)
	}
}

func newReconcileUnregisteredIgnoredFixture(t *testing.T) (store.Run, string) {
	t.Helper()
	home, repo, _ := newReconcileWorkspace(t)
	root := filepath.Join(home, ".roundfix", "worktrees")
	if err := os.WriteFile(filepath.Join(repo, ".gitignore"), []byte(".storage/\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	gitImplement(t, repo, "add", ".gitignore")
	gitImplement(t, repo, "commit", "-m", "test: ignore runtime storage")
	run, ref := createReconcileRun(t, home, repo, root, "ma/widget-flow", store.StateFailed)
	gitImplement(t, repo, "worktree", "remove", "--force", ref.Path)
	gitImplement(t, repo, "branch", "-D", ref.Branch)
	file := filepath.Join(ref.Path, ".storage", "postgres", "18", "docker", "base", "1")
	if err := os.MkdirAll(filepath.Dir(file), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(file, []byte("test volume\n"), 0o444); err != nil {
		t.Fatal(err)
	}
	return run, root
}

func TestReconcileEmptyDirectorySweepWithoutRunCandidates(t *testing.T) {
	t.Parallel()
	home, _, _ := newReconcileWorkspace(t)
	path := filepath.Join(home, ".roundfix", "worktrees", "other-0123abcd")
	if err := os.MkdirAll(path, 0o755); err != nil {
		t.Fatal(err)
	}
	report := itemReconcileReport(t, "--apply")
	if report.Summary.Total != 0 || report.DebrisSummary.EmptyDirectoriesRemoved != 1 || len(report.EmptyDirectoryCandidates) != 1 || report.EmptyDirectoryCandidates[0].Action != "removed" {
		t.Fatalf("empty-only apply = %+v", report)
	}
	assertReconcilePathState(t, path, false)
}

func TestReconcileEmptyDirectoryInspectionFailureIsReported(t *testing.T) {
	t.Parallel()
	path := filepath.Join(t.TempDir(), "not-a-directory")
	if err := os.WriteFile(path, []byte("keep\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	for _, apply := range []bool{false, true} {
		report := reconcileReport{}
		if apply {
			applyReconcileEmptyDirectories(path, &report)
		} else {
			inspectReconcileEmptyDirectories(path, false, &report)
		}
		if report.Summary.OperationalFailures != 1 || len(report.EmptyDirectoryCandidates) != 1 {
			t.Fatalf("failed inspection = %+v", report)
		}
		candidate := report.EmptyDirectoryCandidates[0]
		if candidate.Action != "preserve" || candidate.RefusalReason == "" || report.DebrisSummary.EmptyDirectoriesRemoved != 0 {
			t.Fatalf("failed candidate = %+v", candidate)
		}
		text := reconcileText(report)
		want := "Empty repository directory: " + path + "\n  action: preserve\n  refusal-reason: " + candidate.RefusalReason + "\nSummary:"
		if !strings.Contains(text, want) {
			t.Fatalf("failure text lacks %q:\n%s", want, text)
		}
	}
	assertReconcilePathState(t, path, true)
}
