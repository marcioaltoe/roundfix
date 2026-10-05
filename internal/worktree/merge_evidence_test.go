package worktree

// Suite: delivery archive evidence through terminal Run inspection and apply.
// Invariant: only pre-delivery Runs with completed same-Spec Tasks can be released;
// dirty paths remain bounded by the archived scope, and apply re-proves the head.
// Boundary IN: disposable local Git repositories and terminal Run metadata.
// Boundary OUT: network and Delivery Queue persistence.

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"
	"testing"
)

func newMergeEvidenceFixture(t *testing.T, id string, status string) (terminalRunFixture, string) {
	t.Helper()
	f := newMergedHeadTestFixture(t, id)
	f.run.LocalBranch = "delivery/" + id
	gitWorktreeTest(t, f.repoDir, "branch", f.run.LocalBranch, "main")
	commitMergedHeadFiles(t, f.ref.Path, map[string]string{"operator.txt": "inherited operator work\n", "renamed.txt": "inherited file\n"}, "fix: operator work inherited from item branch")
	commitMergedHeadTask(t, f.ref.Path, f.run.SpecSlug, "task_01", "feature.txt", "delivered\n")
	delivery := commitMergedHeadSpec(t, f.repoDir, f.run.SpecSlug, true, map[string]string{"task_01": status}, "qa-report-2026-10-01.md", map[string]string{
		archivedLeftoverRoot(f) + "/_prd.md": "archived planning\n",
		"operator.txt":                       "inherited operator work\n", "renamed.txt": "inherited file\n", "feature.txt": "delivered\n",
	})
	commitMergedHeadFiles(t, f.repoDir, map[string]string{"operator.txt": "later default edit\n"}, "fix: later Spec edits inherited work")
	gitWorktreeTest(t, f.repoDir, "mv", "renamed.txt", "later-name.txt")
	gitWorktreeTest(t, f.repoDir, "commit", "-m", "refactor: later Spec renames inherited file")
	return f, delivery
}

func TestMergeEvidenceSupersedesARunWhoseFilesDivergedAfterTheMerge(t *testing.T) {
	t.Parallel()
	f, delivery := newMergeEvidenceFixture(t, "delivery-diverged", "completed")
	result := inspectMergedHeadTestRun(t, f, nil)
	want := boundedReconciliationReason(fmt.Sprintf("Run work is superseded by the delivery of Spec %s: delivery commit %s archived it on default branch %q; 1 Task commit(s) completed, 1 other commit(s) superseded", f.run.SpecSlug, delivery[:12], "main"))
	if result.Reason != want {
		t.Fatalf("reason = %q, want %q", result.Reason, want)
	}
	assertMergeEvidenceRelease(t, f, nil)
}

func TestMergeEvidenceSupersedesADivergedRunWithAMergeRecord(t *testing.T) {
	t.Parallel()
	f, _ := newMergeEvidenceFixture(t, "delivery-record", "completed")
	// The recorded candidate also diverges from the earlier Run's inherited work.
	candidate := commitMergedHeadFiles(t, f.repoDir, map[string]string{"operator.txt": "candidate differs\n"}, "fix: candidate work")
	assertMergeEvidenceRelease(t, f, []MergedHead{mergedHeadTestRecord(f, candidate)})
}

func assertMergeEvidenceRelease(t *testing.T, f terminalRunFixture, records []MergedHead) {
	t.Helper()
	result := inspectMergedHeadTestRun(t, f, records)
	assertTerminalRunReconciliation(t, result, f.run, ReconciliationSuperseded)
	if !strings.Contains(result.Reason, "superseded by the delivery of Spec") || len(result.Reason) > reconciliationReasonMaxBytes {
		t.Fatalf("reason = %q", result.Reason)
	}
	if result.evidence.proofHead != mergedHeadTestHead(t, f.repoDir, "main") || result.evidence.proofRef != "main" {
		t.Fatalf("proof head/ref = %s/%s", result.evidence.proofHead, result.evidence.proofRef)
	}
	if err := ApplyTerminalRun(context.Background(), &recordingTerminalRunStore{}, result); err != nil {
		t.Fatal(err)
	}
	assertPathRemoved(t, f.ref.Path)
	assertBranchRemoved(t, f.repoDir, f.ref.Branch)
}

func TestMergeEvidenceKeepsARunWithAPendingTaskCommit(t *testing.T) {
	t.Parallel()
	f, _ := newMergeEvidenceFixture(t, "delivery-pending", "pending")
	result := inspectMergedHeadTestRun(t, f, nil)
	assertTerminalRunReconciliation(t, result, f.run, ReconciliationUnintegrated)
	if !strings.Contains(result.Reason, "Task task_01 is not completed") || !strings.Contains(result.Reason, mergedHeadTestHead(t, f.ref.Path, "HEAD")[:12]) {
		t.Fatalf("reason = %q", result.Reason)
	}
	assertMergeEvidencePreserved(t, f, result)
}

func TestMergeEvidenceRequiresTheArchivedPRD(t *testing.T) {
	t.Parallel()
	f, _ := newMergeEvidenceFixture(t, "delivery-no-prd", "completed")
	gitWorktreeTest(t, f.repoDir, "rm", "--", archivedLeftoverRoot(f)+"/_prd.md")
	gitWorktreeTest(t, f.repoDir, "commit", "-m", "docs: remove archive proof")
	result := inspectMergedHeadTestRun(t, f, nil)
	assertTerminalRunReconciliation(t, result, f.run, ReconciliationUnintegrated)
	assertMergeEvidencePreserved(t, f, result)
}

func TestMergeEvidenceKeepsARunThatHoldsTheDeliveryCommit(t *testing.T) {
	t.Parallel()
	f, delivery := newMergeEvidenceFixture(t, "delivery-held", "completed")
	gitWorktreeTest(t, f.ref.Path, "reset", "--hard", delivery)
	f.commitRunChange(t, "operator.txt", "post-delivery work\n")
	result := inspectMergedHeadTestRun(t, f, nil)
	assertTerminalRunReconciliation(t, result, f.run, ReconciliationUnintegrated)
	assertMergeEvidencePreserved(t, f, result)
}

func TestMergeEvidenceKeepsAnUndeclaredDirtyPath(t *testing.T) {
	t.Parallel()
	f, _ := newMergeEvidenceFixture(t, "delivery-dirty", "completed")
	mustWriteWorktreeTest(t, filepath.Join(f.ref.Path, "outside.txt"), "uncommitted\n")
	result := inspectMergedHeadTestRun(t, f, nil)
	assertTerminalRunReconciliation(t, result, f.run, ReconciliationDirty)
	assertMergeEvidencePreserved(t, f, result)
}

func TestMergeEvidenceApplyRefusesAMovedDefaultBranch(t *testing.T) {
	t.Parallel()
	f, _ := newMergeEvidenceFixture(t, "delivery-stale", "completed")
	result := inspectMergedHeadTestRun(t, f, nil)
	assertTerminalRunReconciliation(t, result, f.run, ReconciliationSuperseded)
	commitMergedHeadFiles(t, f.repoDir, map[string]string{"later.txt": "later\n"}, "fix: move default head")
	err := ApplyTerminalRun(context.Background(), &recordingTerminalRunStore{}, result)
	if err == nil || !strings.Contains(err.Error(), "evidence is stale") {
		t.Fatalf("apply error = %v", err)
	}
	assertPathExists(t, f.ref.Path)
	assertRunBranchExists(t, f.repoDir, f.ref.Branch)
}

func assertMergeEvidencePreserved(t *testing.T, f terminalRunFixture, result RunWorktreeReconciliation) {
	t.Helper()
	if err := ApplyTerminalRun(context.Background(), &recordingTerminalRunStore{}, result); err == nil {
		t.Fatal("apply accepted unproven Run")
	}
	assertPathExists(t, f.ref.Path)
	assertRunBranchExists(t, f.repoDir, f.ref.Branch)
}

func TestMergeEvidenceKeepsATaskCommitOfAnotherSpec(t *testing.T) {
	t.Parallel()
	f, _ := newMergeEvidenceFixture(t, "delivery-other-spec", "completed")
	commit := commitMergedHeadTask(t, f.ref.Path, "other-spec", "task_01", "other.txt", "preserve\n")
	result := inspectMergedHeadTestRun(t, f, nil)
	assertTerminalRunReconciliation(t, result, f.run, ReconciliationUnintegrated)
	if !strings.Contains(result.Reason, commit[:12]) || !strings.Contains(result.Reason, "belongs to Spec other-spec") {
		t.Fatalf("reason = %q", result.Reason)
	}
	assertMergeEvidencePreserved(t, f, result)
}

func TestMergeEvidenceSupersedesScopedDirtyPathsAtTheDefaultHead(t *testing.T) {
	t.Parallel()
	f, delivery := newMergeEvidenceFixture(t, "delivery-scoped", "completed")
	// A record without the archive must fall back to the default head. Dirty
	// scope is declared there, rather than at the selected record's head.
	taskPath := archivedLeftoverRoot(f) + "/task_01.md"
	task, err := gitBlobAtHead(context.Background(), execGitRunner{}, f.repoDir, delivery, taskPath)
	if err != nil {
		t.Fatal(err)
	}
	commitMergedHeadFiles(t, f.repoDir, map[string]string{taskPath: string(task) + "\n## Context\n\n- interface: `declared.txt`\n"}, "docs: declare archived scope")
	mustWriteWorktreeTest(t, filepath.Join(f.ref.Path, "declared.txt"), "scoped leftover\n")
	records := []MergedHead{mergedHeadTestRecord(f, f.run.HeadSHA)}
	result := inspectMergedHeadTestRun(t, f, records)
	assertTerminalRunReconciliation(t, result, f.run, ReconciliationSuperseded)
	if result.evidence.proofRef != "main" || !strings.Contains(result.Reason, "uncommitted path(s)") {
		t.Fatalf("proof = %#v", result)
	}
	if err := ApplyTerminalRun(context.Background(), &recordingTerminalRunStore{}, result); err != nil {
		t.Fatal(err)
	}
	assertPathRemoved(t, f.ref.Path)
	assertBranchRemoved(t, f.repoDir, f.ref.Branch)
}

func TestMergeEvidenceDoesNotRequireAnArchivedQAReport(t *testing.T) {
	t.Parallel()
	for _, absentTarget := range []bool{false, true} {
		t.Run(map[bool]string{false: "ancestry-miss", true: "absent-target"}[absentTarget], func(t *testing.T) {
			f, _ := newMergeEvidenceFixture(t, "delivery-no-qa", "completed")
			gitWorktreeTest(t, f.repoDir, "rm", "--", archivedLeftoverRoot(f)+"/qa/qa-report-2026-10-01.md")
			gitWorktreeTest(t, f.repoDir, "commit", "-m", "docs: remove QA report")
			if absentTarget {
				gitWorktreeTest(t, f.repoDir, "branch", "-D", f.run.LocalBranch)
			}
			assertMergeEvidenceRelease(t, f, nil)
		})
	}
}

func TestMergeEvidenceSupersedesNonTaskCommitsWithoutRepresentation(t *testing.T) {
	t.Parallel()
	t.Run("newer-QA", func(t *testing.T) {
		f, _ := newMergeEvidenceFixture(t, "delivery-newer-qa", "completed")
		commitQAReport(t, f.ref.Path, f.run.SpecSlug, "qa-report-2026-10-02.md", false, "fail")
		assertMergeEvidenceRelease(t, f, nil)
	})
	t.Run("non-Task-Spec-trailers", func(t *testing.T) {
		f, _ := newMergeEvidenceFixture(t, "delivery-non-task", "completed")
		commitMergedHeadFiles(t, f.ref.Path, map[string]string{"other.txt": "committed work\n"}, "fix: operator work\n\nRoundfix-Spec: first-spec\nRoundfix-Spec: second-spec")
		assertMergeEvidenceRelease(t, f, nil)
	})
}
