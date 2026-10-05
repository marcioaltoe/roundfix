package worktree

// Suite: item branch reconciliation against local delivery evidence.
// Invariant: only a non-live, proven branch with a clean derived worktree is removed.
// Boundary IN: real disposable Git repositories and item refs.
// Boundary OUT: Delivery Queue loading, Run cleanup and network.

import (
	"context"
	"path/filepath"
	"strings"
	"testing"
)

func newItemBranchFixture(t *testing.T) (terminalRunFixture, ItemRef) {
	t.Helper()
	f, _ := newMergeEvidenceFixture(t, "item-branch", "completed")
	branch := "roundfix/deliver-" + f.run.SpecSlug + "-0123456789abcdef"
	ref, err := ItemRefFor(f.repoDir, filepath.Join(t.TempDir(), "items"), branch)
	if err != nil {
		t.Fatal(err)
	}
	gitWorktreeTest(t, f.repoDir, "worktree", "add", "-b", branch, ref.Path, mergedHeadTestHead(t, f.ref.Path, "HEAD"))
	return f, ref
}

func inspectItemBranchTest(t *testing.T, f terminalRunFixture, ref ItemRef, live []string) ItemBranchReconciliation {
	t.Helper()
	results, err := InspectItemBranches(context.Background(), f.repoDir, ref.location, live)
	if err != nil {
		t.Fatal(err)
	}
	if len(results) != 1 {
		t.Fatalf("candidates = %+v", results)
	}
	return results[0]
}

func TestItemBranchOfAMergedSpecIsReleasable(t *testing.T) {
	t.Parallel()
	f, ref := newItemBranchFixture(t)
	r := inspectItemBranchTest(t, f, ref, nil)
	if !r.Releasable || r.Branch != ref.Branch || !samePath(r.Worktree, ref.Path) || r.Head == "" || !strings.HasPrefix(r.Proof, "item branch is superseded by the delivery of Spec") {
		t.Fatalf("inspection = %+v", r)
	}
	assertPathExists(t, ref.Path)
	assertRunBranchExists(t, f.repoDir, ref.Branch)
}

func TestItemBranchOfALiveItemIsPreserved(t *testing.T) {
	t.Parallel()
	f, ref := newItemBranchFixture(t)
	r := inspectItemBranchTest(t, f, ref, []string{ref.Branch})
	assertItemBranchPreserved(t, f, ref, r, "live Delivery Queue item")
}

func TestItemBranchOfAnUnmergedSpecIsPreserved(t *testing.T) {
	t.Parallel()
	f, ref := newItemBranchFixture(t)
	gitWorktreeTest(t, f.repoDir, "rm", "--", archivedLeftoverRoot(f)+"/_prd.md")
	gitWorktreeTest(t, f.repoDir, "commit", "-m", "docs: remove archive proof")
	assertItemBranchPreserved(t, f, ref, inspectItemBranchTest(t, f, ref, nil), "no merge evidence")
}

func TestItemBranchWithADirtyWorktreeIsPreserved(t *testing.T) {
	t.Parallel()
	f, ref := newItemBranchFixture(t)
	mustWriteWorktreeTest(t, filepath.Join(ref.Path, "untracked.txt"), "preserve\n")
	assertItemBranchPreserved(t, f, ref, inspectItemBranchTest(t, f, ref, nil), "dirty")
}

func TestApplyItemBranchRemovesTheCleanWorktreeAndBranch(t *testing.T) {
	t.Parallel()
	f, ref := newItemBranchFixture(t)
	r := inspectItemBranchTest(t, f, ref, nil)
	if err := ApplyItemBranch(context.Background(), r); err != nil {
		t.Fatal(err)
	}
	assertPathRemoved(t, ref.Path)
	assertBranchRemoved(t, f.repoDir, ref.Branch)
	assertPathExists(t, f.ref.Path)
	assertRunBranchExists(t, f.repoDir, f.ref.Branch)
}

func TestApplyItemBranchRefusesAMovedHead(t *testing.T) {
	t.Parallel()
	f, ref := newItemBranchFixture(t)
	r := inspectItemBranchTest(t, f, ref, nil)
	commitMergedHeadFiles(t, ref.Path, map[string]string{"later.txt": "later\n"}, "fix: move item head")
	assertItemBranchPreserved(t, f, ref, r, "stale")
}

func assertItemBranchPreserved(t *testing.T, f terminalRunFixture, ref ItemRef, r ItemBranchReconciliation, reason string) {
	t.Helper()
	err := ApplyItemBranch(context.Background(), r)
	if err == nil || !strings.Contains(err.Error(), reason) {
		t.Fatalf("inspection = %+v; apply = %v, want %q", r, err, reason)
	}
	assertPathExists(t, ref.Path)
	assertRunBranchExists(t, f.repoDir, ref.Branch)
}

func TestItemBranchSafetyBoundaries(t *testing.T) {
	t.Parallel()
	for _, change := range []string{"wrong-path", "default-moved", "dirty-after-proof", "worktree-removed", "post-delivery", "pending-task", "other-spec-task", "bare-branch", "unrelated-refs", "ambiguous-short-ref", "duplicate-worktree"} {
		t.Run(change, func(t *testing.T) {
			f, ref := newItemBranchFixture(t)
			r := inspectItemBranchTest(t, f, ref, nil)
			switch change {
			case "wrong-path":
				moved := filepath.Join(t.TempDir(), "moved")
				gitWorktreeTest(t, f.repoDir, "worktree", "move", ref.Path, moved)
				r = inspectItemBranchTest(t, f, ref, nil)
				if r.Releasable || !strings.Contains(r.RefusalReason, "derived path") {
					t.Fatalf("inspection = %+v", r)
				}
			case "default-moved":
				commitMergedHeadFiles(t, f.repoDir, map[string]string{"later.txt": "later\n"}, "fix: move default")
			case "dirty-after-proof":
				mustWriteWorktreeTest(t, filepath.Join(ref.Path, "dirty.txt"), "dirty\n")
			case "worktree-removed":
				gitWorktreeTest(t, f.repoDir, "worktree", "remove", ref.Path)
			case "post-delivery":
				gitWorktreeTest(t, ref.Path, "reset", "--hard", "main")
				r = inspectItemBranchTest(t, f, ref, nil)
			case "pending-task":
				task := archivedLeftoverRoot(f) + "/task_01.md"
				commitMergedHeadFiles(t, f.repoDir, map[string]string{task: "---\ntask: task_01\nstatus: pending\n---\n"}, "docs: reopen archived Task")
				r = inspectItemBranchTest(t, f, ref, nil)
			case "other-spec-task":
				commitMergedHeadTask(t, ref.Path, "other-spec", "task_01", "other.txt", "other\n")
				r = inspectItemBranchTest(t, f, ref, nil)
			case "bare-branch":
				gitWorktreeTest(t, f.repoDir, "worktree", "remove", ref.Path)
				r = inspectItemBranchTest(t, f, ref, nil)
				if !r.Releasable || r.Worktree != "" {
					t.Fatalf("inspection = %+v", r)
				}
				if err := ApplyItemBranch(context.Background(), r); err != nil {
					t.Fatal(err)
				}
				assertBranchRemoved(t, f.repoDir, ref.Branch)
				return
			case "ambiguous-short-ref":
				gitWorktreeTest(t, f.repoDir, "tag", ref.Branch)
				r = inspectItemBranchTest(t, f, ref, nil)
				if !r.Releasable {
					t.Fatalf("inspection = %+v", r)
				}
				return
			case "duplicate-worktree":
				gitWorktreeTest(t, f.repoDir, "worktree", "add", "--force", filepath.Join(t.TempDir(), "duplicate"), ref.Branch)
				r = inspectItemBranchTest(t, f, ref, nil)
			case "unrelated-refs":
				for _, branch := range []string{"roundfix/deliver-invalid", "roundfix/deliver-slug-0123456789abcdeg", "operator/work"} {
					gitWorktreeTest(t, f.repoDir, "branch", branch)
				}
				r = inspectItemBranchTest(t, f, ref, nil)
				if !r.Releasable {
					t.Fatalf("inspection = %+v", r)
				}
				return
			}
			if err := ApplyItemBranch(context.Background(), r); err == nil {
				t.Fatal("apply accepted changed or unproven item")
			}
			assertRunBranchExists(t, f.repoDir, ref.Branch)
		})
	}
}
