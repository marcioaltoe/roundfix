// Suite: reconcile carry-forward staging debris.
// Invariant: reconcile reports every staging candidate and mutates one only in an explicit releasing mode.
// Boundary IN: public reconcile runner, JSON report, text report, and real Git worktree registrations.
// Boundary OUT: owner-process proof and locking internals, owned by internal/worktree/staging_test.go.
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

	"roundfix/internal/store"
	runworktree "roundfix/internal/worktree"
)

func TestReconcileDryRunReportsStaleStagingWithoutRemovingIt(t *testing.T) {
	_, repoDir, _ := newReconcileWorkspace(t)
	worktreePath := addLockedInitializingStagingForReconcileTest(t, repoDir)
	t.Cleanup(func() { removeReconcileStagingForTest(t, repoDir, worktreePath) })

	report, textReport := runStagingReconcileForTest(t, repoDir, []string{"reconcile", "--format=json"})
	if len(report.StagingCandidates) != 1 {
		t.Fatalf("staging candidates = %+v, want one", report.StagingCandidates)
	}
	candidate := report.StagingCandidates[0]
	if candidate.Kind != "staging" || candidate.Worktree != worktreePath || candidate.Action != "would reclaim with --apply" || !strings.Contains(candidate.Proof, "locked initializing") {
		t.Fatalf("stale staging candidate = %+v", candidate)
	}
	if report.DebrisSummary.StagingCandidates != 1 || report.DebrisSummary.StagingApplied != 0 {
		t.Fatalf("dry-run debris summary = %+v", report.DebrisSummary)
	}
	if !stagingRegisteredForReconcileTest(t, repoDir, worktreePath) {
		t.Fatal("dry-run removed the stale staging worktree")
	}
	for _, want := range []string{"Staging candidate:", "locked initializing", "staging-candidates=1 staging-applied=0"} {
		if !strings.Contains(textReport, want) {
			t.Fatalf("text report does not contain %q:\n%s", want, textReport)
		}
	}
}

func TestReconcileApplyReleasesStaleStaging(t *testing.T) {
	_, repoDir, _ := newReconcileWorkspace(t)
	worktreePath := addLockedInitializingStagingForReconcileTest(t, repoDir)

	report, _ := runStagingReconcileForTest(t, repoDir, []string{"reconcile", "--apply", "--format=json"})
	if len(report.StagingCandidates) != 1 || report.StagingCandidates[0].Action != "released" {
		t.Fatalf("applied staging candidates = %+v", report.StagingCandidates)
	}
	if report.DebrisSummary.StagingApplied != 1 {
		t.Fatalf("apply debris summary = %+v, want one staging applied", report.DebrisSummary)
	}
	assertReconcileStagingReleased(t, repoDir, worktreePath)
}

func TestReconcileCarryForwardReleasesStaleStagingFirst(t *testing.T) {
	fixture := newCarryForwardFixture(t, store.StateUnresolved, []implementSeed{{id: "task_01", title: "Build the core"}})
	worktreePath := addLockedInitializingStagingForReconcileTest(t, fixture.repoDir)

	report, _ := runStagingReconcileForTest(t, fixture.repoDir, []string{"reconcile", fixture.run.ID, "--carry-forward", "--format=json"})
	if len(report.StagingCandidates) != 1 || report.StagingCandidates[0].Action != "released" {
		t.Fatalf("carry-forward staging candidates = %+v", report.StagingCandidates)
	}
	if report.DebrisSummary.StagingApplied != 1 {
		t.Fatalf("carry-forward debris summary = %+v, want one staging applied", report.DebrisSummary)
	}
	assertReconcileStagingReleased(t, fixture.repoDir, worktreePath)
	if _, err := os.Stat(filepath.Join(fixture.repoDir, "src", "task_01.txt")); err != nil {
		t.Fatalf("carried-forward implementation: %v", err)
	}
}

func TestReconcileKeepsALiveStagingAsPreserved(t *testing.T) {
	_, repoDir, _ := newReconcileWorkspace(t)
	head := strings.TrimSpace(gitImplementOutput(t, repoDir, "rev-parse", "HEAD"))
	staging, err := runworktree.AddCarryForwardStaging(t.Context(), repoDir, head, t.TempDir())
	if err != nil {
		t.Fatalf("add live staging: %v", err)
	}
	t.Cleanup(func() {
		if stagingRegisteredForReconcileTest(t, repoDir, staging.Worktree) {
			if err := staging.Remove(context.Background()); err != nil {
				t.Errorf("remove live staging: %v", err)
			}
		}
	})

	report, _ := runStagingReconcileForTest(t, repoDir, []string{"reconcile", "--apply", "--format=json"})
	if len(report.StagingCandidates) != 1 || report.StagingCandidates[0].Action != "preserve" {
		t.Fatalf("live staging candidates = %+v, want preserved entry", report.StagingCandidates)
	}
	foundPreserved := false
	for _, candidate := range report.PreservedCandidates {
		if candidate.Kind == "staging" && candidate.Worktree == staging.Worktree &&
			candidate.Action == "preserve" && strings.Contains(candidate.RefusalReason, "live") {
			foundPreserved = true
			break
		}
	}
	if !foundPreserved {
		t.Fatalf("live staging is missing from preserved candidates: %+v", report.PreservedCandidates)
	}
	if !stagingRegisteredForReconcileTest(t, repoDir, staging.Worktree) {
		t.Fatal("apply removed a live owner's staging worktree")
	}
}

func runStagingReconcileForTest(t *testing.T, repoDir string, args []string) (reconcileReport, string) {
	t.Helper()
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	code := runCLIContext(t, context.Background(), args, &stdout, &stderr)
	if code != exitOK {
		t.Fatalf("reconcile exit = %d, want 0; stderr=%q stdout=%q", code, stderr.String(), stdout.String())
	}
	if stderr.Len() != 0 {
		t.Fatalf("reconcile stderr = %q, want empty", stderr.String())
	}
	var report reconcileReport
	if err := json.Unmarshal(stdout.Bytes(), &report); err != nil {
		t.Fatalf("decode reconcile report: %v\n%s", err, stdout.String())
	}

	var textStdout bytes.Buffer
	var textStderr bytes.Buffer
	textArgs := make([]string, 0, len(args))
	for _, arg := range args {
		if arg != "--format=json" {
			textArgs = append(textArgs, arg)
		}
	}
	if slicesApplyMode(args) {
		return report, reconcileText(report)
	}
	textCode := runCLIContext(t, context.Background(), textArgs, &textStdout, &textStderr)
	if textCode != exitOK {
		t.Fatalf("text reconcile exit = %d, want 0; stderr=%q stdout=%q", textCode, textStderr.String(), textStdout.String())
	}
	return report, textStdout.String()
}

func slicesApplyMode(args []string) bool {
	for _, arg := range args {
		if arg == "--apply" || arg == "--carry-forward" {
			return true
		}
	}
	return false
}

func addLockedInitializingStagingForReconcileTest(t *testing.T, repoDir string) string {
	t.Helper()
	root, err := os.MkdirTemp(t.TempDir(), "roundfix-carry-forward-")
	if err != nil {
		t.Fatalf("create staging root: %v", err)
	}
	root, err = filepath.EvalSymlinks(root)
	if err != nil {
		t.Fatalf("resolve staging root: %v", err)
	}
	worktreePath := filepath.Join(root, "worktree")
	head := strings.TrimSpace(gitImplementOutput(t, repoDir, "rev-parse", "HEAD"))
	gitImplement(t, repoDir, "worktree", "add", "--detach", worktreePath, head)
	gitImplement(t, repoDir, "worktree", "lock", "--reason", "initializing", worktreePath)
	return worktreePath
}

func removeReconcileStagingForTest(t *testing.T, repoDir string, worktreePath string) {
	t.Helper()
	if stagingRegisteredForReconcileTest(t, repoDir, worktreePath) {
		gitImplement(t, repoDir, "worktree", "remove", "--force", "--force", worktreePath)
	}
	if err := os.RemoveAll(filepath.Dir(worktreePath)); err != nil {
		t.Errorf("remove staging root: %v", err)
	}
}

func stagingRegisteredForReconcileTest(t *testing.T, repoDir string, worktreePath string) bool {
	t.Helper()
	return strings.Contains(gitImplementOutput(t, repoDir, "worktree", "list", "--porcelain"), "worktree "+worktreePath)
}

func assertReconcileStagingReleased(t *testing.T, repoDir string, worktreePath string) {
	t.Helper()
	if stagingRegisteredForReconcileTest(t, repoDir, worktreePath) {
		t.Fatal("staging worktree remains registered")
	}
	if _, err := os.Stat(filepath.Dir(worktreePath)); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("staging root stat error = %v, want not exist", err)
	}
}
