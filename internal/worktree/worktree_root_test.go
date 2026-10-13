// Suite: terminal Run reconciliation under a symbolic-linked Worktree Root.
// Invariant: location links are followed; repository and worktree links are refused.
// Boundary IN: Create, InspectTerminalRun, ApplyTerminalRun and PruneTerminalReport.
// Boundary OUT: missing-prefix path comparison and consumer surfaces belong to later Tasks.
package worktree

import (
	"context"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"roundfix/internal/store"
)

func TestInspectTerminalRunResolvesASymlinkedWorktreeRoot(t *testing.T) {
	t.Parallel()
	fixture := newSymlinkedWorktreeRootFixture(t, "inspect-linked-root")
	fixture.commitRunChange(t, "integrated.txt", "integrated\n")
	gitWorktreeTest(t, fixture.repoDir, "merge", "--ff-only", fixture.ref.Branch)

	result, err := InspectTerminalRun(t.Context(), fixture.run)
	if err != nil {
		t.Fatalf("inspect terminal Run: %v", err)
	}
	assertTerminalRunReconciliation(t, result, fixture.run, ReconciliationSafe)
	if result.Reason == reconciliationReasonWorktreePath {
		t.Fatalf("linked Worktree Root was refused: %q", result.Reason)
	}
	assertPathExists(t, fixture.ref.Path)
	assertRunBranchExists(t, fixture.repoDir, fixture.ref.Branch)
}

func TestApplyTerminalRunReleasesAWorktreeUnderASymlinkedRoot(t *testing.T) {
	t.Parallel()
	fixture := newSymlinkedWorktreeRootFixture(t, "apply-linked-root")
	fixture.commitRunChange(t, "integrated.txt", "integrated\n")
	gitWorktreeTest(t, fixture.repoDir, "merge", "--ff-only", fixture.ref.Branch)

	result, err := InspectTerminalRun(t.Context(), fixture.run)
	if err != nil {
		t.Fatalf("inspect terminal Run: %v", err)
	}
	if err := ApplyTerminalRun(t.Context(), &recordingTerminalRunStore{}, result); err != nil {
		t.Fatalf("apply terminal Run reconciliation: %v", err)
	}
	assertPathRemoved(t, fixture.ref.Path)
	assertBranchRemoved(t, fixture.repoDir, fixture.ref.Branch)
	released, err := InspectTerminalRun(t.Context(), fixture.run)
	if err != nil {
		t.Fatalf("inspect released Run: %v", err)
	}
	assertTerminalRunReconciliation(t, released, fixture.run, ReconciliationReleased)
}

func TestPruneTerminalReportPrunesUnderASymlinkedRoot(t *testing.T) {
	t.Parallel()
	fixture := newSymlinkedWorktreeRootFixture(t, "prune-linked-root")
	fixture.commitRunChange(t, "integrated.txt", "integrated\n")
	gitWorktreeTest(t, fixture.repoDir, "merge", "--ff-only", fixture.ref.Branch)
	task, err := CreateTask(t.Context(), fixture.ref, "task_01", nil)
	if err != nil {
		t.Fatalf("create kept Task Worktree: %v", err)
	}
	taskPath := canonicalPath(task.Path)
	location := filepath.Dir(filepath.Dir(fixture.ref.Path))

	pruned, err := PruneTerminalReport(t.Context(), fixture.repoDir, location, &recordingTerminalRunStore{}, func(_ context.Context, runID string) (store.Run, bool, error) {
		return fixture.run, runID == fixture.run.ID, nil
	})
	if err != nil {
		t.Fatalf("prune terminal Run: %v", err)
	}
	want := []PrunedRef{
		{RunID: fixture.run.ID, Path: fixture.run.WorkDir, Branch: fixture.ref.Branch},
		{RunID: fixture.run.ID, TaskID: task.TaskID, Path: taskPath, Branch: task.Branch},
	}
	if !slices.Equal(pruned, want) {
		t.Fatalf("pruned refs = %#v, want %#v", pruned, want)
	}
	assertPathRemoved(t, fixture.ref.Path)
	assertPathRemoved(t, task.Path)
	assertBranchRemoved(t, fixture.repoDir, fixture.ref.Branch)
	assertBranchRemoved(t, fixture.repoDir, task.Branch)
}

func TestInspectTerminalRunRefusesAPathResolvingOutsideItsWorktreeRoot(t *testing.T) {
	t.Parallel()
	for _, name := range []string{"repository-directory-link", "worktree-directory-link", "dangling-worktree-link"} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			fixture := newSymlinkedWorktreeRootFixture(t, "unsafe-linked-root")
			location := filepath.Dir(filepath.Dir(fixture.ref.Path))
			root := filepath.Join(location, "unsafe-root")
			if err := os.Mkdir(root, 0o755); err != nil {
				t.Fatalf("create alternate Worktree Root: %v", err)
			}
			parent := filepath.Join(root, "repository")
			recordedPath := filepath.Join(parent, filepath.Base(fixture.ref.Path))
			realPath := canonicalPath(fixture.ref.Path)
			if name == "repository-directory-link" {
				if err := os.Symlink(filepath.Dir(realPath), parent); err != nil {
					t.Fatalf("create repository directory link: %v", err)
				}
			} else {
				if err := os.Mkdir(parent, 0o755); err != nil {
					t.Fatalf("create repository directory: %v", err)
				}
				target := realPath
				if name == "dangling-worktree-link" {
					target = filepath.Join(canonicalPath(t.TempDir()), "missing-worktree")
				}
				if err := os.Symlink(target, recordedPath); err != nil {
					t.Fatalf("create worktree directory link: %v", err)
				}
			}
			fixture.run.WorkDir = recordedPath

			result, err := InspectTerminalRun(t.Context(), fixture.run)
			if err != nil {
				t.Fatalf("inspect terminal Run: %v", err)
			}
			assertTerminalRunReconciliation(t, result, fixture.run, ReconciliationUnknown)
			if result.Reason != reconciliationReasonWorktreePath {
				t.Fatalf("unsafe path reason = %q, want %q", result.Reason, reconciliationReasonWorktreePath)
			}
			if err := ApplyTerminalRun(t.Context(), &recordingTerminalRunStore{}, result); err == nil {
				t.Fatal("apply accepted an unsafe recorded path")
			}
			pruned, err := PruneTerminalReport(t.Context(), fixture.repoDir, location, &recordingTerminalRunStore{}, func(_ context.Context, runID string) (store.Run, bool, error) {
				return fixture.run, runID == fixture.run.ID, nil
			})
			if err != nil || len(pruned) != 0 {
				t.Fatalf("prune unsafe Run: refs=%#v err=%v", pruned, err)
			}
			assertPathExists(t, realPath)
			assertRunBranchExists(t, fixture.repoDir, fixture.ref.Branch)
			if info, err := os.Lstat(recordedPath); err != nil || info == nil {
				t.Fatalf("recorded path was removed: %v", err)
			}
		})
	}
}

func TestInspectTerminalRunTreatsAMissingWorktreeUnderASymlinkedRootAsAbsent(t *testing.T) {
	t.Parallel()
	fixture := newSymlinkedWorktreeRootFixture(t, "missing-linked-root")
	gitWorktreeTest(t, fixture.repoDir, "worktree", "remove", fixture.ref.Path)
	gitWorktreeTest(t, fixture.repoDir, "branch", "-D", fixture.ref.Branch)

	result, err := InspectTerminalRun(t.Context(), fixture.run)
	if err != nil {
		t.Fatalf("inspect terminal Run: %v", err)
	}
	assertTerminalRunReconciliation(t, result, fixture.run, ReconciliationReleased)
}

func newSymlinkedWorktreeRootFixture(t *testing.T, runID string) terminalRunFixture {
	t.Helper()
	repoDir := canonicalPath(initWorktreeRepo(t))
	commitWorktreeFile(t, repoDir, "shared.txt", "base\n", "initial")
	headSHA := strings.TrimSpace(gitWorktreeTest(t, repoDir, "rev-parse", "HEAD"))
	realLocation := canonicalPath(t.TempDir())
	location := filepath.Join(canonicalPath(t.TempDir()), "linked-worktrees")
	if err := os.Symlink(realLocation, location); err != nil {
		t.Fatalf("create Worktree Root link: %v", err)
	}
	ref, err := Create(t.Context(), CreateOptions{
		UserRoot: repoDir,
		Location: location,
		RunID:    runID,
		HeadSHA:  headSHA,
	})
	if err != nil {
		t.Fatalf("create Run Worktree under linked root: %v", err)
	}
	if ref.Path == canonicalPath(ref.Path) {
		t.Fatalf("fixture lost linked location spelling: %q", ref.Path)
	}
	return terminalRunFixture{
		integrationFixture: integrationFixture{repoDir: repoDir, ref: ref},
		run: store.Run{
			ID:          runID,
			Kind:        store.KindImplement,
			State:       store.StateClean,
			GitRoot:     repoDir,
			LocalBranch: "main",
			WorkDir:     ref.Path,
			SpecSlug:    "0262-a-run-worktree-under-a-symlinked-root-is-cleaned",
		},
	}
}
