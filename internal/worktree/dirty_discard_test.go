package worktree

// Suite: confirmed Dirty Discard boundaries.
// Invariant: only inspected, unchanged leftovers of an archived ended Run are discarded.
// Boundary IN: disposable Git repositories and linked Run Worktrees.
// Boundary OUT: Branch Disposition persistence and the Run Database.

import (
	"errors"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"roundfix/internal/store"
)

func TestAnEndedRunOfAnArchivedSpecWithAnOutsideLeftoverOffersADirtyDiscard(t *testing.T) {
	t.Parallel()
	f := newArchivedLeftoverFixture(t, "dirty-discard-offer")
	f.run.State = store.StateFailed
	addDeclaredLeftovers(t, f)
	mustWriteWorktreeTest(t, filepath.Join(f.ref.Path, "unrelated.txt"), "leftover\n")
	result := inspectMergedHeadTestRun(t, f, nil)
	want := []string{"declared.txt", "docs/specs/" + f.run.SpecSlug + "/_prd.md", "unrelated.txt"}
	if result.State != ReconciliationDirty || result.Reason != reconciliationReasonDirty+"; unrelated.txt is outside the archived Spec's declared and recorded paths" || !slices.Equal(result.DiscardPaths, want) {
		t.Fatalf("candidate = %+v, want paths %q and original outside-path reason", result, want)
	}
	if result.evidence == nil || result.evidence.proofHead == "" || !slices.Equal(result.evidence.dirtyPaths, want) {
		t.Fatalf("candidate evidence = %+v", result.evidence)
	}
}

func TestADirtyDiscardIsNeverOfferedOutsideItsBounds(t *testing.T) {
	t.Parallel()
	for _, name := range []string{"unrepresented Task commit", "not archived on default", "Clean outcome", "scoped leftovers"} {
		t.Run(name, func(t *testing.T) {
			f := newArchivedLeftoverFixture(t, "dirty-discard-bounds")
			f.run.State = store.StateStopped
			addDeclaredLeftovers(t, f)
			var merged []MergedHead
			switch name {
			case "unrepresented Task commit":
				commitMergedHeadTask(t, f.ref.Path, f.run.SpecSlug, "task_02", "unique.txt", "keep\n")
			case "not archived on default":
				// Keep an archived proof head, but move main back before its archive.
				head := mergedHeadTestHead(t, f.repoDir, "main")
				merged = []MergedHead{mergedHeadTestRecord(f, head)}
				initial := strings.TrimSpace(gitWorktreeTest(t, f.repoDir, "rev-list", "--max-parents=0", "HEAD"))
				gitWorktreeTest(t, f.repoDir, "checkout", "-b", "unarchived-default", initial)
				gitWorktreeTest(t, f.repoDir, "branch", "-f", "main", initial)
				gitWorktreeTest(t, f.repoDir, "checkout", "main")
			case "Clean outcome":
				f.run.State = store.StateClean
			}
			if name != "scoped leftovers" {
				mustWriteWorktreeTest(t, filepath.Join(f.ref.Path, "unrelated.txt"), "keep\n")
			}
			result := inspectMergedHeadTestRun(t, f, merged)
			if len(result.DiscardPaths) != 0 {
				t.Fatalf("out-of-bounds candidate = %+v", result)
			}
			if name == "scoped leftovers" && result.State != ReconciliationSuperseded {
				t.Fatalf("scoped state = %q", result.State)
			}
		})
	}
}

func TestDiscardDirtyTerminalRunRemovesTheRunWorktreeAndRunBranch(t *testing.T) {
	t.Parallel()
	f := newArchivedLeftoverFixture(t, "dirty-discard-remove")
	f.run.State = store.StateUnresolved
	addDeclaredLeftovers(t, f)
	mustWriteWorktreeTest(t, filepath.Join(f.ref.Path, "unrelated.txt"), "discard\n")
	result := inspectMergedHeadTestRun(t, f, nil)
	if err := DiscardDirtyTerminalRun(t.Context(), result); err != nil {
		t.Fatal(err)
	}
	assertPathRemoved(t, f.ref.Path)
	assertBranchRemoved(t, f.repoDir, f.ref.Branch)
}

func TestDiscardDirtyTerminalRunRefusesAChangedPathSet(t *testing.T) {
	t.Parallel()
	for _, name := range []string{"added path", "not a candidate", "changed result", "changed proof head", "changed Run head", "changed target head", "changed merged record"} {
		t.Run(name, func(t *testing.T) {
			f := newArchivedLeftoverFixture(t, "dirty-discard-stale")
			f.run.State = store.StateTimedOut
			mustWriteWorktreeTest(t, filepath.Join(f.ref.Path, "unrelated.txt"), "keep\n")
			head := mergedHeadTestHead(t, f.repoDir, "main")
			result := inspectMergedHeadTestRun(t, f, []MergedHead{mergedHeadTestRecord(f, head)})
			switch name {
			case "added path":
				mustWriteWorktreeTest(t, filepath.Join(f.ref.Path, "another.txt"), "keep\n")
			case "not a candidate":
				result = RunWorktreeReconciliation{RunID: f.run.ID}
			case "changed result":
				result.DiscardPaths[0] = "different.txt"
			case "changed proof head":
				// Without a pinned delivery record the proof follows the default head.
				result = inspectMergedHeadTestRun(t, f, nil)
				commitWorktreeFile(t, f.repoDir, "later.txt", "later\n", "feat: advance default")
			case "changed Run head":
				commitMergedHeadTask(t, f.ref.Path, f.run.SpecSlug, "task_02", "new-task.txt", "unique\n")
			case "changed target head":
				commitWorktreeFile(t, f.repoDir, "later-target.txt", "later\n", "feat: advance target")
			case "changed merged record":
				result.evidence.merged[0].Head = strings.Repeat("0", 40)
			}
			if err := DiscardDirtyTerminalRun(t.Context(), result); err == nil {
				t.Fatal("accepted changed or uninspected evidence")
			}
			assertPathExists(t, f.ref.Path)
			assertRunBranchExists(t, f.repoDir, f.ref.Branch)
		})
	}
}

func TestDirtyDiscardKeepsTheDeliveryProofWhenDefaultHasAdvanced(t *testing.T) {
	t.Parallel()
	f := newArchivedLeftoverFixture(t, "dirty-discard-delivery-proof")
	f.run.State = store.StateBudgetExceeded
	addDeclaredLeftovers(t, f)
	mustWriteWorktreeTest(t, filepath.Join(f.ref.Path, "unrelated.txt"), "discard\n")
	proofHead := mergedHeadTestHead(t, f.repoDir, "main")
	record := mergedHeadTestRecord(f, proofHead)
	commitWorktreeFile(t, f.repoDir, "later.txt", "later\n", "feat: advance default")
	result := inspectMergedHeadTestRun(t, f, []MergedHead{record})
	if len(result.DiscardPaths) != 3 || result.evidence == nil || result.evidence.proofHead != proofHead || !slices.Equal(result.evidence.merged, []MergedHead{record}) || !slices.Equal(result.evidence.mergedSnapshot, []MergedHead{record}) {
		t.Fatalf("delivery proof candidate = %+v; evidence = %+v", result, result.evidence)
	}
	if err := DiscardDirtyTerminalRun(t.Context(), result); err != nil {
		t.Fatal(err)
	}
	assertPathRemoved(t, f.ref.Path)
	assertBranchRemoved(t, f.repoDir, f.ref.Branch)
}

func TestDirtyDiscardReportsALockedWorktreeRemovalFailure(t *testing.T) {
	t.Parallel()
	f := newArchivedLeftoverFixture(t, "dirty-discard-locked")
	f.run.State = store.StateStopped
	mustWriteWorktreeTest(t, filepath.Join(f.ref.Path, "unrelated.txt"), "keep\n")
	result := inspectMergedHeadTestRun(t, f, nil)
	gitWorktreeTest(t, f.repoDir, "worktree", "lock", f.ref.Path)
	err := DiscardDirtyTerminalRun(t.Context(), result)
	var failure *TerminalRunApplyError
	if !errors.As(err, &failure) || !failure.WorktreeRemaining || !failure.RunBranchRemaining {
		t.Fatalf("removal failure = %v, want remaining worktree and Run Branch", err)
	}
	assertPathExists(t, f.ref.Path)
	assertRunBranchExists(t, f.repoDir, f.ref.Branch)
}
