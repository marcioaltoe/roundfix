package worktree

// Suite: archived Spec representation and scoped leftovers.
// Invariant: cleanup requires represented commits and re-proved dirty paths.
// Boundary IN: disposable Git repositories, archived Task files, linked worktrees.
// Boundary OUT: GitHub, delivery records, and the real Run Database.

import (
	"context"
	"path/filepath"
	"strings"
	"testing"

	"roundfix/internal/spec"
	"roundfix/internal/store"
)

func archivedLeftoverRoot(f terminalRunFixture) string {
	return filepath.ToSlash(filepath.Join(spec.ArchiveDir(spec.ArchiveKindSpec), f.run.SpecSlug))
}

func newArchivedLeftoverFixture(t *testing.T, id string) terminalRunFixture {
	t.Helper()
	f := newMergedHeadTestFixture(t, id)
	commitMergedHeadFiles(t, f.ref.Path, map[string]string{filepath.ToSlash(filepath.Join("docs/specs", f.run.SpecSlug, "_prd.md")): "old planning\n"}, "docs: plan Spec")
	root := archivedLeftoverRoot(f)
	commitMergedHeadSpec(t, f.repoDir, f.run.SpecSlug, true, map[string]string{"task_01": "completed"}, "qa-report-2026-10-01.md", map[string]string{root + "/_prd.md": "archived planning\n"})
	commitMergedHeadFiles(t, f.repoDir, map[string]string{root + "/task_01.md": "---\ntask: task_01\nspec: " + f.run.SpecSlug + "\nstatus: completed\ntype: backend\ncomplexity: low\n---\n\n# Task\n\n## Context\n\n- interface: `declared.txt`\n- creates: `created.txt`\n- deletes: `deleted.txt`\n- instruction: `instruction.txt`\n\n## Verification\n\n- `true`\n\n## Recorded paths\n\n- `recorded.txt`\n"}, "docs: record scope")
	return f
}

func addDeclaredLeftovers(t *testing.T, f terminalRunFixture) {
	t.Helper()
	mustWriteWorktreeTest(t, filepath.Join(f.ref.Path, "declared.txt"), "leftover\n")
	mustWriteWorktreeTest(t, filepath.Join(f.ref.Path, "docs/specs", f.run.SpecSlug, "_prd.md"), "unfinished planning\n")
}

func TestMergedSpecRunWithSpecDirectoryCommitsIsSuperseded(t *testing.T) {
	t.Parallel()
	f := newArchivedLeftoverFixture(t, "archived-spec-commits")
	result := inspectMergedHeadTestRun(t, f, nil)
	assertTerminalRunReconciliation(t, result, f.run, ReconciliationSuperseded)
	recordResult := inspectMergedHeadTestRun(t, f, []MergedHead{mergedHeadTestRecord(f, mergedHeadTestHead(t, f.repoDir, "main"))})
	if !strings.Contains(recordResult.Reason, "1 Spec-directory path(s) archived") || len(recordResult.Reason) > reconciliationReasonMaxBytes {
		t.Fatalf("record reason = %q", recordResult.Reason)
	}
	if !strings.Contains(result.Reason, "1 Spec-directory path(s) archived") {
		t.Fatalf("reason = %q", result.Reason)
	}
}

func TestMergedSpecRunWithDeclaredLeftoversIsSuperseded(t *testing.T) {
	t.Parallel()
	f := newArchivedLeftoverFixture(t, "archived-declared-leftovers")
	addDeclaredLeftovers(t, f)
	result := inspectMergedHeadTestRun(t, f, nil)
	assertTerminalRunReconciliation(t, result, f.run, ReconciliationSuperseded)
	if !strings.Contains(result.Reason, "; 2 uncommitted path(s) superseded by the archived Spec") {
		t.Fatalf("reason = %q", result.Reason)
	}
}

func TestMergedSpecRunWithAnUndeclaredLeftoverStaysDirty(t *testing.T) {
	t.Parallel()
	for _, file := range []string{"unrelated.txt", "instruction.txt"} {
		t.Run(file, func(t *testing.T) {
			f := newArchivedLeftoverFixture(t, "archived-undeclared-leftover")
			addDeclaredLeftovers(t, f)
			mustWriteWorktreeTest(t, filepath.Join(f.ref.Path, file), "preserve\n")
			result := inspectMergedHeadTestRun(t, f, nil)
			assertTerminalRunReconciliation(t, result, f.run, ReconciliationDirty)
			if result.Reason != reconciliationReasonDirty {
				t.Fatalf("reason = %q", result.Reason)
			}
			if err := ApplyTerminalRun(context.Background(), &recordingTerminalRunStore{}, result); err == nil {
				t.Fatal("apply accepted dirty work")
			}
			assertPathExists(t, f.ref.Path)
			assertRunBranchExists(t, f.repoDir, f.ref.Branch)
		})
	}
}

func TestMergedSpecRunWithAnUnrepresentedTaskCommitStaysUnintegrated(t *testing.T) {
	t.Parallel()
	f := newArchivedLeftoverFixture(t, "archived-unrepresented-task")
	commitMergedHeadTask(t, f.ref.Path, f.run.SpecSlug, "task_02", "unique.txt", "preserve\n")
	result := inspectMergedHeadTestRun(t, f, nil)
	assertTerminalRunReconciliation(t, result, f.run, ReconciliationUnintegrated)
	addDeclaredLeftovers(t, f)
	dirty := inspectMergedHeadTestRun(t, f, nil)
	assertTerminalRunReconciliation(t, dirty, f.run, ReconciliationDirty)
	assertPathExists(t, f.ref.Path)
}

func TestApplyRefusesWhenTheLeftoverSetChanged(t *testing.T) {
	t.Parallel()
	f := newArchivedLeftoverFixture(t, "archived-leftover-stale")
	addDeclaredLeftovers(t, f)
	result := inspectMergedHeadTestRun(t, f, nil)
	mustWriteWorktreeTest(t, filepath.Join(f.ref.Path, "recorded.txt"), "new scoped leftover\n")
	err := ApplyTerminalRun(context.Background(), &recordingTerminalRunStore{}, result)
	if err == nil || !strings.Contains(err.Error(), "evidence is stale") {
		t.Fatalf("apply error = %v", err)
	}
	assertPathExists(t, f.ref.Path)
	assertRunBranchExists(t, f.repoDir, f.ref.Branch)
}

func TestApplyRemovesASupersededLeftoverWorktree(t *testing.T) {
	t.Parallel()
	f := newArchivedLeftoverFixture(t, "archived-leftover-apply")
	addDeclaredLeftovers(t, f)
	result := inspectMergedHeadTestRun(t, f, nil)
	if err := ApplyTerminalRun(context.Background(), &recordingTerminalRunStore{}, result); err != nil {
		t.Fatalf("apply: %v", err)
	}
	assertPathRemoved(t, f.ref.Path)
	assertBranchRemoved(t, f.repoDir, f.ref.Branch)
}

func TestAbsentTargetBranchIsReleasedOnTheMergedHead(t *testing.T) {
	t.Parallel()
	f := newArchivedLeftoverFixture(t, "archived-absent-target")
	f.run.LocalBranch = "delivery/absent-target"
	gitWorktreeTest(t, f.repoDir, "symbolic-ref", "refs/remotes/origin/HEAD", "refs/remotes/origin/main")
	// The working root recorded by a Run can itself be a linked worktree.
	linkedRoot := filepath.Join(t.TempDir(), "linked-checkout")
	gitWorktreeTest(t, f.repoDir, "worktree", "add", "--detach", linkedRoot, "main")
	f.run.GitRoot = canonicalPath(linkedRoot)
	result, err := ClassifyRunBranchSet(context.Background(), canonicalPath(f.repoDir), f.run.LocalBranch, f.run.SpecSlug, []store.Run{f.run})
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Releasable) != 1 || result.Releasable[0] != f.ref.Branch || len(result.Preserved) != 0 {
		t.Fatalf("classification = %#v", result)
	}
	_, targetErr := localBranchExists(context.Background(), execGitRunner{}, f.run.GitRoot, f.run.LocalBranch)
	if targetErr != nil {
		t.Fatalf("target lookup: %v", targetErr)
	}
	terminal := inspectMergedHeadTestRun(t, f, nil)
	assertTerminalRunReconciliation(t, terminal, f.run, ReconciliationSuperseded)
	if err := ApplyRunBranchCandidate(context.Background(), result, f.ref.Branch); err != nil {
		t.Fatalf("apply branch candidate: %v", err)
	}
	assertPathRemoved(t, f.ref.Path)
	assertBranchRemoved(t, f.repoDir, f.ref.Branch)
}

func TestArchivedSpecLeftoverScopeKinds(t *testing.T) {
	t.Parallel()
	for _, file := range []string{"created.txt", "deleted.txt", "recorded.txt"} {
		t.Run(file, func(t *testing.T) {
			f := newArchivedLeftoverFixture(t, "archived-scope-kind")
			mustWriteWorktreeTest(t, filepath.Join(f.ref.Path, file), "scoped\n")
			result := inspectMergedHeadTestRun(t, f, nil)
			assertTerminalRunReconciliation(t, result, f.run, ReconciliationSuperseded)
		})
	}
}

func TestArchivedSpecLeftoversRequireTheArchivedPRD(t *testing.T) {
	t.Parallel()
	f := newArchivedLeftoverFixture(t, "archived-prd-required")
	gitWorktreeTest(t, f.repoDir, "rm", "--", archivedLeftoverRoot(f)+"/_prd.md")
	gitWorktreeTest(t, f.repoDir, "commit", "-m", "docs: remove archived PRD")
	addDeclaredLeftovers(t, f)
	result := inspectMergedHeadTestRun(t, f, nil)
	assertTerminalRunReconciliation(t, result, f.run, ReconciliationDirty)
}

func TestArchivedSpecLeftoverRenameKeepsAnOutsideSource(t *testing.T) {
	t.Parallel()
	f := newArchivedLeftoverFixture(t, "archived-rename-source")
	f.commitRunChange(t, "outside.txt", "rename me\n")
	commitWorktreeFile(t, f.repoDir, "outside.txt", "rename me\n", "feat: represent outside content")
	gitWorktreeTest(t, f.ref.Path, "mv", "outside.txt", "created.txt")
	result := inspectMergedHeadTestRun(t, f, nil)
	assertTerminalRunReconciliation(t, result, f.run, ReconciliationDirty)
}

func TestArchivedSpecDirtyPathsKeepLiteralFilenames(t *testing.T) {
	t.Parallel()
	f := newArchivedLeftoverFixture(t, "archived-literal-paths")
	file := filepath.Join(f.ref.Path, "docs/specs", f.run.SpecSlug, "quoted \"name\"\nfile.md")
	mustWriteWorktreeTest(t, file, "literal\n")
	result := inspectMergedHeadTestRun(t, f, nil)
	assertTerminalRunReconciliation(t, result, f.run, ReconciliationSuperseded)
	if len(result.evidence.dirtyPaths) != 1 || result.evidence.dirtyPaths[0] != filepath.ToSlash(filepath.Join("docs/specs", f.run.SpecSlug, "quoted \"name\"\nfile.md")) {
		t.Fatalf("paths = %q", result.evidence.dirtyPaths)
	}
}

func TestArchivedSpecWithMergeEvidenceSupersedesOtherCommitPaths(t *testing.T) {
	t.Parallel()
	f := newArchivedLeftoverFixture(t, "archived-mixed-paths")
	commitMergedHeadFiles(t, f.ref.Path, map[string]string{
		filepath.ToSlash(filepath.Join("docs/specs", f.run.SpecSlug, "_prd.md")): "new plan\n",
		"outside.txt": "unique\n",
	}, "docs: planning with outside work")
	result := inspectMergedHeadTestRun(t, f, nil)
	assertTerminalRunReconciliation(t, result, f.run, ReconciliationSuperseded)
}
