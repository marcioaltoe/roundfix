// Suite: terminal Run consumers under a symlinked Worktree Root.
// Invariant: a linked location preserves recorded spelling and supports proven release.
// Boundary IN: reconcile commands, merged-Spec release, Run Database, and real Git worktrees.
// Boundary OUT: path resolution and unsafe child links, owned by internal/worktree/worktree_root_test.go.
package cli

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"roundfix/internal/store"
	runworktree "roundfix/internal/worktree"
)

func TestReconcileClassifiesARunUnderASymlinkedWorktreeRoot(t *testing.T) {
	t.Parallel()
	fixture := newReconcileSymlinkedRootFixture(t)

	report := runReconcileMergedCommand(t, fixture)

	result := reconcileMergedResult(t, report, fixture.run.ID)
	if result.Classification != "safe" || result.Worktree != fixture.run.WorkDir {
		t.Fatalf("linked-root dry-run result = %+v, want safe with recorded path %q", result, fixture.run.WorkDir)
	}
	for _, result := range report.Results {
		if strings.Contains(result.RefusalReason, "recorded Run Worktree path is unsafe") {
			t.Fatalf("linked-root result refuses the recorded path: %+v", result)
		}
	}
	assertReconcilePathState(t, fixture.ref.Path, true)
	assertReconcileBranchState(t, fixture.repoDir, fixture.ref.Branch, true)
}

func TestReconcileApplyReleasesARunUnderASymlinkedWorktreeRoot(t *testing.T) {
	t.Parallel()
	fixture := newReconcileSymlinkedRootFixture(t)
	var stdout, stderr bytes.Buffer

	code := runCLIContext(t, t.Context(), []string{"reconcile", "--apply", "--format", "json"}, &stdout, &stderr)

	if code != exitOK || stderr.Len() != 0 {
		t.Fatalf("linked-root apply exit = %d; stderr=%q stdout=%q", code, stderr.String(), stdout.String())
	}
	var report reconcileReport
	if err := json.Unmarshal(stdout.Bytes(), &report); err != nil {
		t.Fatalf("decode linked-root apply report: %v\n%s", err, stdout.String())
	}
	result := reconcileMergedResult(t, report, fixture.run.ID)
	if result.Classification != "safe" || result.Action != "released" {
		t.Fatalf("linked-root apply result = %+v, want safe and released", result)
	}
	assertReconcileSymlinkedRootReleased(t, fixture)
}

func TestReleaseMergedSpecRunsCleansUnderASymlinkedWorktreeRoot(t *testing.T) {
	t.Parallel()
	fixture := newReconcileSymlinkedRootFixture(t)
	runStore, err := store.Open(t.Context(), fixture.homeDir)
	if err != nil {
		t.Fatalf("open merged-Spec Run Database: %v", err)
	}
	t.Cleanup(func() {
		if err := runStore.Close(); err != nil {
			t.Errorf("close merged-Spec Run Database: %v", err)
		}
	})

	err = releaseMergedSpecRuns(t.Context(), runStore, fixture.repoDir, runworktree.MergedHead{
		SpecSlug: fixture.run.SpecSlug, TargetBranch: fixture.run.LocalBranch,
		Head: fixture.mergedHead, MergeCommit: fixture.mergedHead, PullRequest: "259",
	})

	if err != nil {
		t.Fatalf("release merged Spec under linked root: %v", err)
	}
	assertReconcileSymlinkedRootReleased(t, fixture)
}

func newReconcileSymlinkedRootFixture(t *testing.T) reconcileMergedFixture {
	t.Helper()
	homeDir, repoDir, realLocation := newReconcileWorkspace(t)
	location := filepath.Join(t.TempDir(), "linked-worktrees")
	if err := os.Symlink(realLocation, location); err != nil {
		t.Fatalf("link worktree location: %v", err)
	}
	ctx := t.Context()
	runStore, err := store.Open(ctx, homeDir)
	if err != nil {
		t.Fatalf("open linked-root Run Database: %v", err)
	}
	defer func() {
		if err := runStore.Close(); err != nil {
			t.Errorf("close linked-root Run Database: %v", err)
		}
	}()
	head := strings.TrimSpace(gitImplementOutput(t, repoDir, "rev-parse", "HEAD"))
	target := "fix/symlinked-root-spec"
	gitImplement(t, repoDir, "branch", target, head)
	run, err := runStore.CreateRun(ctx, store.CreateRunRequest{
		Kind: store.KindImplement, GitRoot: repoDir, LocalBranch: target,
		HeadSHA: head, SpecSlug: reconcileMergedSpecSlug, Agent: "codex",
	})
	if err != nil {
		t.Fatalf("create linked-root Run: %v", err)
	}
	ref, err := runworktree.Create(ctx, runworktree.CreateOptions{
		UserRoot: repoDir, Location: location, RunID: run.ID, HeadSHA: head,
	})
	if err != nil {
		t.Fatalf("create linked-root Run Worktree: %v", err)
	}
	run, err = runStore.SetRunWorkDir(ctx, run.ID, ref.Path)
	if err != nil {
		t.Fatalf("record linked-root Run Worktree: %v", err)
	}
	resolved, err := filepath.EvalSymlinks(ref.Path)
	if err != nil {
		t.Fatalf("resolve linked-root Run Worktree: %v", err)
	}
	registered := gitImplementOutput(t, repoDir, "worktree", "list", "--porcelain")
	if run.WorkDir != ref.Path || !strings.HasPrefix(run.WorkDir, location+string(filepath.Separator)) ||
		resolved == run.WorkDir || !strings.Contains(registered, "worktree "+resolved+"\n") {
		t.Fatalf("linked-root fixture: recorded=%q resolved=%q registrations=\n%s", run.WorkDir, resolved, registered)
	}
	commitReconcileMergedTask(t, ref.Path, run.SpecSlug, reconcileMergedTaskID)
	gitImplement(t, repoDir, "checkout", "main")
	gitImplement(t, repoDir, "merge", "--ff-only", ref.Branch)
	mergedHead := strings.TrimSpace(gitImplementOutput(t, repoDir, "rev-parse", "HEAD"))
	gitImplement(t, repoDir, "branch", "-d", target)
	completed, err := runStore.CompleteRun(ctx, run.ID, store.StateClean)
	if err != nil {
		t.Fatalf("record terminal linked-root Run: %v", err)
	}
	fixture := reconcileMergedFixture{
		homeDir: homeDir, repoDir: repoDir, run: completed.Run, ref: ref, mergedHead: mergedHead,
	}
	seedReconcileDeliveryMergedHead(t, fixture, run.SpecSlug)
	return fixture
}

func assertReconcileSymlinkedRootReleased(t *testing.T, fixture reconcileMergedFixture) {
	t.Helper()
	assertReconcilePathState(t, fixture.ref.Path, false)
	assertReconcileBranchState(t, fixture.repoDir, fixture.ref.Branch, false)
	registered := gitImplementOutput(t, fixture.repoDir, "worktree", "list", "--porcelain")
	if strings.Contains(registered, "branch refs/heads/"+fixture.ref.Branch+"\n") ||
		strings.Contains(registered, "/"+filepath.Base(fixture.ref.Path)+"\n") {
		t.Fatalf("released Run Worktree remains registered:\n%s", registered)
	}
}
