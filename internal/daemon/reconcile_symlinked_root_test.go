// Suite: superseded Run Branch disposition under a symlinked Worktree Root.
// Invariant: a linked location does not refuse a clean Run whose commits reached its target.
// Boundary IN: persisted terminal Run metadata, linked location, and real Git ancestry.
// Boundary OUT: unsafe child links, owned by internal/worktree/worktree_root_test.go.
package daemon

import (
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"roundfix/internal/gittest"
	"roundfix/internal/store"
	runworktree "roundfix/internal/worktree"
)

func TestClassifySupersededBranchUnderASymlinkedWorktreeRoot(t *testing.T) {
	t.Parallel()
	ctx := t.Context()
	root, err := filepath.EvalSymlinks(t.TempDir())
	if err != nil {
		t.Fatalf("resolve Git root: %v", err)
	}
	gittest.InitRepo(t, root, "-b", "main")
	gittest.PersistIdentity(t, root)
	gittest.Run(t, root, "commit", "--allow-empty", "-m", "test: seed target")
	base := strings.TrimSpace(gittest.Run(t, root, "rev-parse", "HEAD"))
	location := filepath.Join(t.TempDir(), "linked-worktrees")
	if err := os.Symlink(t.TempDir(), location); err != nil {
		t.Fatalf("link worktree location: %v", err)
	}
	runStore, err := store.Open(ctx, t.TempDir())
	if err != nil {
		t.Fatalf("open linked-root Run Database: %v", err)
	}
	t.Cleanup(func() {
		if err := runStore.Close(); err != nil {
			t.Errorf("close linked-root Run Database: %v", err)
		}
	})
	run, err := runStore.CreateRun(ctx, store.CreateRunRequest{
		Kind: store.KindImplement, GitRoot: root, LocalBranch: "main",
		HeadSHA: base, SpecSlug: "disposition-spec", Agent: "codex",
	})
	if err != nil {
		t.Fatalf("create linked-root Run: %v", err)
	}
	ref, err := runworktree.Create(ctx, runworktree.CreateOptions{
		UserRoot: root, Location: location, RunID: run.ID, HeadSHA: base,
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
	registered := gittest.Run(t, root, "worktree", "list", "--porcelain")
	if run.WorkDir != ref.Path || !strings.HasPrefix(run.WorkDir, location+string(filepath.Separator)) ||
		resolved == run.WorkDir || !strings.Contains(registered, "worktree "+resolved+"\n") {
		t.Fatalf("linked-root fixture: recorded=%q resolved=%q registrations=\n%s", run.WorkDir, resolved, registered)
	}
	fixture := supersededBranchFixture{root: root, ref: ref, run: run}
	commit := fixture.commitRunChange(t, "work.txt", "work\n", "feat: add work")
	gittest.Run(t, root, "merge", "--ff-only", ref.Branch)
	completed, err := runStore.CompleteRun(ctx, run.ID, store.StateStopped)
	if err != nil {
		t.Fatalf("record terminal linked-root Run: %v", err)
	}

	got, err := ClassifySupersededBranch(ctx, completed.Run, nil)

	if err != nil {
		t.Fatalf("ClassifySupersededBranch() error = %v", err)
	}
	if !got.Superseded || !got.Reachable || got.RefusalReason != "" ||
		!slices.Equal(got.Commits, []string{commit}) || got.Worktree != ref.Path {
		t.Fatalf("linked-root disposition = %+v, want reachable superseded with recorded path %q", got, ref.Path)
	}
	assertSupersededBranchSurface(t, fixture, true)
}
