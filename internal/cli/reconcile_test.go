// Suite: reconcile superseded-branch disposition.
// Invariant: reconcile reports without mutation and discards only after the explicit flag writes durable evidence.
// Boundary IN: public reconcile runner, Run Database metadata, local Git surfaces, and Artifact Root records.
// Boundary OUT: terminal Run classification internals, owned by internal/daemon/reconcile_test.go.
package cli

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"roundfix/internal/daemon"
	"roundfix/internal/store"
)

func TestReconcileDiscardWritesTheRecordBeforeRemoving(t *testing.T) {
	t.Parallel()
	homeDir, repoDir, location := newReconcileWorkspace(t)
	run, ref := createReconcileRun(t, homeDir, repoDir, location, "ma/widget-flow", store.StateStopped)
	first := commitReconcileDispositionChange(t, ref.Path, "first.txt", "first\n", "feat: add first change")
	second := commitReconcileDispositionChange(t, ref.Path, "second.txt", "second\n", "feat: add second change")
	gitImplement(t, repoDir, "merge", "--ff-only", ref.Branch)
	recordPath := filepath.Join(builtinArtifactDirForRepo(t, repoDir), "runs", run.ID, "branch-disposition.json")
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := runCLIContext(
		t,
		context.Background(),
		[]string{"reconcile", run.ID, "--discard-superseded", "--format=json"},
		&stdout,
		&stderr,
	)

	if code != exitOK {
		t.Fatalf("discard exit = %d, want 0 stderr=%q stdout=%q", code, stderr.String(), stdout.String())
	}
	if stderr.Len() != 0 {
		t.Fatalf("discard stderr = %q, want empty", stderr.String())
	}
	var report reconcileReport
	if err := json.Unmarshal(stdout.Bytes(), &report); err != nil {
		t.Fatalf("decode discard report: %v\n%s", err, stdout.String())
	}
	if len(report.Results) != 1 {
		t.Fatalf("discard results = %+v, want one", report.Results)
	}
	result := report.Results[0]
	if result.Classification != "superseded" || result.Action != "discarded" ||
		!strings.Contains(result.Evidence, "branch record: "+recordPath) {
		t.Fatalf("discard result = %+v", result)
	}
	record, err := os.ReadFile(recordPath)
	if err != nil {
		t.Fatalf("read branch record after removal: %v", err)
	}
	for _, want := range []string{daemon.BranchDispositionRecordSchemaVersion, first, second, "first.txt", "second.txt"} {
		if !bytes.Contains(record, []byte(want)) {
			t.Fatalf("branch record does not contain %q: %s", want, record)
		}
	}
	assertReconcilePathState(t, ref.Path, false)
	assertReconcileBranchState(t, repoDir, ref.Branch, false)
}

func TestReconcileDiscardRefusesAnUnreachableCommit(t *testing.T) {
	t.Parallel()
	homeDir, repoDir, location := newReconcileWorkspace(t)
	run, ref := createReconcileRun(t, homeDir, repoDir, location, "ma/widget-flow", store.StateStopped)
	unreachable := commitReconcileDispositionChange(t, ref.Path, "run-only.txt", "run only\n", "feat: retain run-only change")
	commitReconcileDispositionChange(t, repoDir, "target-only.txt", "target only\n", "feat: advance target")
	recordPath := filepath.Join(builtinArtifactDirForRepo(t, repoDir), "runs", run.ID, "branch-disposition.json")
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := runCLIContext(
		t,
		context.Background(),
		[]string{"reconcile", run.ID, "--discard-superseded", "--format=json"},
		&stdout,
		&stderr,
	)

	if code != exitPreflight {
		t.Fatalf("refused discard exit = %d, want %d stderr=%q stdout=%q", code, exitPreflight, stderr.String(), stdout.String())
	}
	var report reconcileReport
	if err := json.Unmarshal(stdout.Bytes(), &report); err != nil {
		t.Fatalf("decode refused discard report: %v\n%s", err, stdout.String())
	}
	if len(report.Results) != 1 ||
		!strings.Contains(report.Results[0].RefusalReason, "unreachable commit "+unreachable) {
		t.Fatalf("refused discard results = %+v", report.Results)
	}
	assertReconcilePathState(t, ref.Path, true)
	assertReconcileBranchState(t, repoDir, ref.Branch, true)
	if _, err := os.Stat(recordPath); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("refused discard record stat error = %v, want not exist", err)
	}
}

func TestReconcileDiscardKeepsSurfaceWhenRecordCannotBeWritten(t *testing.T) {
	t.Parallel()
	homeDir, repoDir, location := newReconcileWorkspace(t)
	run, ref := createReconcileRun(t, homeDir, repoDir, location, "ma/widget-flow", store.StateStopped)
	commitReconcileDispositionChange(t, ref.Path, "work.txt", "work\n", "feat: add work")
	gitImplement(t, repoDir, "merge", "--ff-only", ref.Branch)
	artifactRoot := builtinArtifactDirForRepo(t, repoDir)
	if err := os.MkdirAll(artifactRoot, 0o755); err != nil {
		t.Fatalf("create Artifact Root: %v", err)
	}
	if err := os.WriteFile(filepath.Join(artifactRoot, "runs"), []byte("not a directory\n"), 0o644); err != nil {
		t.Fatalf("block branch record directory: %v", err)
	}
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := runCLIContext(
		t,
		context.Background(),
		[]string{"reconcile", run.ID, "--discard-superseded", "--format=json"},
		&stdout,
		&stderr,
	)

	if code != exitRunFailed {
		t.Fatalf("record failure exit = %d, want %d stderr=%q stdout=%q", code, exitRunFailed, stderr.String(), stdout.String())
	}
	if !strings.Contains(stderr.String(), "reconciliation operation(s) failed") {
		t.Fatalf("record failure stderr = %q", stderr.String())
	}
	assertReconcilePathState(t, ref.Path, true)
	assertReconcileBranchState(t, repoDir, ref.Branch, true)
}

func TestReconcileWithoutTheFlagRemovesNothing(t *testing.T) {
	t.Parallel()
	homeDir, repoDir, location := newReconcileWorkspace(t)
	run, ref := createReconcileRun(t, homeDir, repoDir, location, "ma/widget-flow", store.StateStopped)
	commitReconcileDispositionChange(t, ref.Path, "work.txt", "work\n", "feat: add work")
	gitImplement(t, repoDir, "merge", "--ff-only", ref.Branch)
	recordPath := filepath.Join(builtinArtifactDirForRepo(t, repoDir), "runs", run.ID, "branch-disposition.json")
	beforeDatabase := readReconcileBytes(t, store.DatabasePath(homeDir))
	beforeGit := reconcileGitSurface(t, repoDir)
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := runCLIContext(
		t,
		context.Background(),
		[]string{"reconcile", run.ID, "--format=json"},
		&stdout,
		&stderr,
	)

	if code != exitOK {
		t.Fatalf("report exit = %d, want 0 stderr=%q stdout=%q", code, stderr.String(), stdout.String())
	}
	var report reconcileReport
	if err := json.Unmarshal(stdout.Bytes(), &report); err != nil {
		t.Fatalf("decode report: %v\n%s", err, stdout.String())
	}
	if len(report.Results) != 1 || report.Results[0].Classification != "superseded" ||
		report.Results[0].Action == "discarded" {
		t.Fatalf("read-only disposition report = %+v", report.Results)
	}
	if got := readReconcileBytes(t, store.DatabasePath(homeDir)); !bytes.Equal(got, beforeDatabase) {
		t.Fatal("reconcile without discard flag changed the Run Database")
	}
	if got := reconcileGitSurface(t, repoDir); got != beforeGit {
		t.Fatalf("reconcile without discard flag changed Git state\nbefore:\n%s\nafter:\n%s", beforeGit, got)
	}
	assertReconcilePathState(t, ref.Path, true)
	assertReconcileBranchState(t, repoDir, ref.Branch, true)
	if _, err := os.Stat(recordPath); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("read-only branch record stat error = %v, want not exist", err)
	}
}

func commitReconcileDispositionChange(t *testing.T, workDir, path, content, message string) string {
	t.Helper()
	mustWrite(t, filepath.Join(workDir, path), content)
	gitImplement(t, workDir, "add", path)
	gitImplement(t, workDir, "commit", "-m", message)
	return strings.TrimSpace(gitImplementOutput(t, workDir, "rev-parse", "HEAD"))
}
