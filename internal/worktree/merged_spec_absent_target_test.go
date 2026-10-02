package worktree

// Suite: absent-target reconciliation through real Git and the Run Store.
// Invariant: only revalidated merged-head proof permits recorded cleanup.
// Boundary IN: disposable repositories, linked worktrees, and local Run Stores.
// Boundary OUT: GitHub and the user's Run Database.

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"roundfix/internal/spec"
	"roundfix/internal/store"
)

func newStoredAbsentTargetFixture(t *testing.T) (storedTerminalRunFixture, string) {
	t.Helper()
	f := newStoredTerminalRunFixture(t, store.StateIntegrationPending)
	gitWorktreeTest(t, f.repoDir, "checkout", "-b", "trunk")
	gitWorktreeTest(t, f.repoDir, "symbolic-ref", "refs/remotes/origin/HEAD", "refs/remotes/origin/trunk")
	active := filepath.ToSlash(filepath.Join("docs/specs", f.run.SpecSlug, "_prd.md"))
	commitMergedHeadFiles(t, f.ref.Path, map[string]string{active: "old planning\n"}, "docs: plan Spec")
	gitWorktreeTest(t, f.repoDir, "branch", "-f", f.run.LocalBranch, f.ref.Branch)
	gitWorktreeTest(t, f.repoDir, "merge", "--squash", f.run.LocalBranch)
	gitWorktreeTest(t, f.repoDir, "commit", "-m", "feat: squash merged Spec")
	gitWorktreeTest(t, f.repoDir, "rm", "--", active)
	root := filepath.ToSlash(filepath.Join(spec.ArchiveDir(spec.ArchiveKindSpec), f.run.SpecSlug))
	commitMergedHeadSpec(t, f.repoDir, f.run.SpecSlug, true, map[string]string{"task_01": "completed"}, "qa-report-2026-10-02.md", map[string]string{
		root + "/_prd.md": "archived planning\n",
	})
	commitMergedHeadFiles(t, f.repoDir, map[string]string{
		root + "/task_01.md": "---\ntask: task_01\nspec: " + f.run.SpecSlug + "\nstatus: completed\ntype: backend\ncomplexity: low\n---\n\n# Task\n\n## Context\n\n- creates: `declared.txt`\n\n## Verification\n\n- `true`\n",
	}, "docs: record declared scope")
	gitWorktreeTest(t, f.repoDir, "branch", "-D", f.run.LocalBranch)
	return f, mergedHeadTestHead(t, f.repoDir, "trunk")
}

func inspectStoredAbsentTarget(t *testing.T, f storedTerminalRunFixture, merged []MergedHead) RunWorktreeReconciliation {
	t.Helper()
	result, err := InspectTerminalRunMerged(context.Background(), f.run, merged)
	if err != nil {
		t.Fatalf("inspect: %v", err)
	}
	return result
}

func assertStoredAbsentTargetReleased(t *testing.T, f storedTerminalRunFixture, result RunWorktreeReconciliation, head, ref string) {
	t.Helper()
	ctx := context.Background()
	assertTerminalRunReconciliation(t, result, f.run, ReconciliationSuperseded)
	assertPathExists(t, f.ref.Path)
	assertRunBranchExists(t, f.repoDir, f.ref.Branch)
	if err := ApplyTerminalRun(ctx, f.runStore, result); err != nil {
		t.Fatalf("apply merged Run with absent target: %v", err)
	}
	assertPathRemoved(t, f.ref.Path)
	assertBranchRemoved(t, f.repoDir, f.ref.Branch)
	events, err := f.runStore.RunEventsAfter(ctx, f.run.ID, 0, 10)
	if err != nil || len(events) != 1 {
		t.Fatalf("reconciliation events: count=%d err=%v", len(events), err)
	}
	var payload map[string]string
	if err := json.Unmarshal(events[0].Event.Payload, &payload); err != nil {
		t.Fatal(err)
	}
	for key, want := range map[string]string{
		"event": "integration_reconciliation", "classification": "superseded",
		"target_head": head, "target_branch": ref, "run_head": result.RunHead,
		"run_branch": f.ref.Branch, "action": "cleanup",
	} {
		if payload[key] != want {
			t.Errorf("event %s = %q, want %q", key, payload[key], want)
		}
	}
	stored, found, err := f.runStore.Run(ctx, f.run.ID)
	if err != nil || !found || stored.State != store.StateClean || stored.LocalBranch != f.run.LocalBranch {
		t.Fatalf("stored Run after reconciliation: %+v found=%v err=%v", stored, found, err)
	}
}

func assertStoredAbsentTargetPreserved(t *testing.T, f storedTerminalRunFixture) {
	t.Helper()
	assertPathExists(t, f.ref.Path)
	assertRunBranchExists(t, f.repoDir, f.ref.Branch)
	events, err := f.runStore.RunEventsAfter(context.Background(), f.run.ID, 0, 10)
	if err != nil || len(events) != 0 {
		t.Fatalf("preserved Run events: count=%d err=%v", len(events), err)
	}
	stored, found, err := f.runStore.Run(context.Background(), f.run.ID)
	if err != nil || !found || stored.State != f.run.State {
		t.Fatalf("preserved Run: %+v found=%v err=%v", stored, found, err)
	}
}

func TestApplyReleasesAMergedRunWhoseTargetBranchIsGone(t *testing.T) {
	t.Parallel()
	for _, test := range []struct {
		name          string
		dirty         bool
		targetPresent bool
	}{
		{name: "clean"},
		{name: "declared leftovers", dirty: true},
		{name: "present target keeps its request", dirty: true, targetPresent: true},
	} {
		t.Run(test.name, func(t *testing.T) {
			f, head := newStoredAbsentTargetFixture(t)
			ref := "trunk"
			if test.targetPresent {
				gitWorktreeTest(t, f.repoDir, "branch", f.run.LocalBranch, f.run.HeadSHA)
				head, ref = f.run.HeadSHA, f.run.LocalBranch
			}
			if test.dirty {
				mustWriteWorktreeTest(t, filepath.Join(f.ref.Path, "declared.txt"), "superseded\n")
			}
			result := inspectStoredAbsentTarget(t, f, nil)
			assertStoredAbsentTargetReleased(t, f, result, head, ref)
		})
	}
}

func TestApplyKeepsAnUndeclaredLeftoverWhenTheTargetIsGone(t *testing.T) {
	t.Parallel()
	f, head := newStoredAbsentTargetFixture(t)
	mustWriteWorktreeTest(t, filepath.Join(f.ref.Path, "declared.txt"), "superseded\n")
	undeclared := filepath.Join(f.ref.Path, "unrelated.txt")
	mustWriteWorktreeTest(t, undeclared, "preserve\n")
	result := inspectStoredAbsentTarget(t, f, nil)
	assertTerminalRunReconciliation(t, result, f.run, ReconciliationDirty)
	if err := ApplyTerminalRun(context.Background(), f.runStore, result); err == nil {
		t.Fatal("apply accepted undeclared leftover")
	}
	assertStoredAbsentTargetPreserved(t, f)
	content, err := os.ReadFile(undeclared)
	if err != nil || string(content) != "preserve\n" {
		t.Fatalf("undeclared leftover changed: %q err=%v", content, err)
	}
	// Removing only the undeclared path must allow the declared leftovers to
	// be released. This distinguishes scoped preservation from refusing all Runs.
	if err := os.Remove(undeclared); err != nil {
		t.Fatal(err)
	}
	result = inspectStoredAbsentTarget(t, f, nil)
	assertStoredAbsentTargetReleased(t, f, result, head, "trunk")
}

func TestApplyStillRefusesAnEmptyTargetHeadWithoutProof(t *testing.T) {
	t.Parallel()
	f, head := newStoredAbsentTargetFixture(t)
	runHead := mergedHeadTestHead(t, f.repoDir, f.ref.Branch)
	_, err := f.runStore.ReconcileIntegration(context.Background(), store.IntegrationReconciliation{
		RunID: f.run.ID, PreviousOutcome: f.run.State, Classification: "superseded",
		RunBranch: f.ref.Branch, RunHead: runHead, TargetBranch: f.run.LocalBranch,
		Worktree: f.run.WorkDir, Reason: "no merged-head proof", Action: "cleanup", Time: time.Now().UTC(),
	})
	if err == nil || !strings.Contains(err.Error(), "target head is required") {
		t.Fatalf("empty target head refusal: %v", err)
	}
	assertStoredAbsentTargetPreserved(t, f)
	// The inspected proof is the positive control for the refused request.
	// A delivery record reads the immutable commit ref rather than a live branch.
	mustWriteWorktreeTest(t, filepath.Join(f.ref.Path, "declared.txt"), "superseded\n")
	result := inspectStoredAbsentTarget(t, f, []MergedHead{{SpecSlug: f.run.SpecSlug, Head: head}})
	assertStoredAbsentTargetReleased(t, f, result, head, head)
}
