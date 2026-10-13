// Suite: exported Git worktree disposal operations.
// Invariant: every exported add or removal uses the repository administration lock and preserves dirty work without force.
// Boundary IN: public disposal helpers, the administration lock, and real local Git repositories.
// Boundary OUT: Daemon disposition recording and disposable checkout cleanup ordering.
package worktree

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"roundfix/internal/testwait"
)

func TestRemoveRegisteredWorktreeWaitsForTheAdminLock(t *testing.T) {
	t.Parallel()
	repoDir, head := newStagingRepository(t)
	worktreePath := filepath.Join(t.TempDir(), "registered")
	gitWorktreeTest(t, repoDir, "worktree", "add", "--detach", worktreePath, head)
	release := holdDisposalAdminLock(t, repoDir)

	runner := newDisposalLockRunner()
	waitCtx, cancelWait := context.WithCancel(t.Context())
	defer cancelWait()
	result := make(chan error, 1)
	go func() {
		result <- removeRegisteredWorktree(waitCtx, runner, repoDir, worktreePath)
		close(result)
	}()
	t.Cleanup(func() {
		cancelWait()
		testwait.Until(t, "remove waiter cleanup", result, (<-chan error)(nil))
	})
	testwait.Until(t, "remove common directory resolution", runner.resolved, result)
	cancelWait()
	err := testwait.Until(t, "cancelled removal", result, (<-chan error)(nil))
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("RemoveRegisteredWorktree() error = %v, want context canceled", err)
	}
	requireNoDisposalWorktreeCommand(t, runner.started)
	if _, err := os.Stat(worktreePath); err != nil {
		t.Fatalf("registered worktree was removed while lock was held: %v", err)
	}
	if got := gitWorktreeTest(t, worktreePath, "rev-parse", "HEAD"); strings.TrimSpace(got) != head {
		t.Fatalf("registered worktree HEAD = %q, want %q", got, head)
	}
	requireDisposalWorktreeRegistration(t, repoDir, worktreePath, true)
	release()
	if err := removeRegisteredWorktree(t.Context(), runner, repoDir, worktreePath); err != nil {
		t.Fatalf("remove after releasing administration lock: %v", err)
	}
	if _, err := os.Stat(worktreePath); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("registered worktree remains after removal: %v", err)
	}
	requireDisposalWorktreeRegistration(t, repoDir, worktreePath, false)
}

func TestAddDetachedWorktreeWaitsForTheAdminLock(t *testing.T) {
	t.Parallel()
	repoDir, head := newStagingRepository(t)
	worktreePath := filepath.Join(t.TempDir(), "detached")
	release := holdDisposalAdminLock(t, repoDir)

	runner := newDisposalLockRunner()
	waitCtx, cancelWait := context.WithCancel(t.Context())
	defer cancelWait()
	result := make(chan error, 1)
	go func() {
		result <- addDetachedWorktree(waitCtx, runner, repoDir, worktreePath, head)
		close(result)
	}()
	t.Cleanup(func() {
		cancelWait()
		testwait.Until(t, "add waiter cleanup", result, (<-chan error)(nil))
	})
	testwait.Until(t, "add common directory resolution", runner.resolved, result)
	cancelWait()
	err := testwait.Until(t, "cancelled addition", result, (<-chan error)(nil))
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("AddDetachedWorktree() error = %v, want context canceled", err)
	}
	requireNoDisposalWorktreeCommand(t, runner.started)
	if _, err := os.Stat(worktreePath); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("detached worktree path exists while lock was held: %v", err)
	}
	requireDisposalWorktreeRegistration(t, repoDir, worktreePath, false)
	release()
	if err := addDetachedWorktree(t.Context(), runner, repoDir, worktreePath, head); err != nil {
		t.Fatalf("add after releasing administration lock: %v", err)
	}
	if got := gitWorktreeTest(t, worktreePath, "rev-parse", "HEAD"); strings.TrimSpace(got) != head {
		t.Fatalf("detached worktree HEAD = %q, want %q", got, head)
	}
	requireDisposalWorktreeRegistration(t, repoDir, worktreePath, true)
}

func TestRemoveRegisteredWorktreeKeepsADirtyWorktree(t *testing.T) {
	t.Parallel()
	repoDir, head := newStagingRepository(t)
	worktreePath := filepath.Join(t.TempDir(), "dirty")
	gitWorktreeTest(t, repoDir, "worktree", "add", "--detach", worktreePath, head)
	dirtyPath := filepath.Join(worktreePath, "dirty.txt")
	mustWriteWorktreeTest(t, dirtyPath, "uncommitted\n")

	err := RemoveRegisteredWorktree(t.Context(), repoDir, worktreePath)

	if err == nil {
		t.Fatal("RemoveRegisteredWorktree() error = nil, want dirty-worktree refusal")
	}
	if got := mustReadWorktreeTest(t, dirtyPath); got != "uncommitted\n" {
		t.Fatalf("dirty worktree file = %q, want preserved content", got)
	}
	registered := gitWorktreeTest(t, repoDir, "worktree", "list", "--porcelain")
	if !strings.Contains(registered, "worktree "+canonicalPath(worktreePath)) {
		t.Fatalf("registered worktrees = %q, want dirty path %q", registered, worktreePath)
	}
}

func holdDisposalAdminLock(t *testing.T, repoDir string) func() {
	t.Helper()
	commonDir := strings.TrimSpace(gitWorktreeTest(t, repoDir, "rev-parse", "--path-format=absolute", "--git-common-dir"))
	started := make(chan struct{}, 1)
	release := make(chan struct{})
	holder := &adminLockTestRunner{
		commonDir: commonDir,
		admin: func(ctx context.Context, _ string, _ ...string) (string, error) {
			started <- struct{}{}
			select {
			case <-release:
				return "", nil
			case <-ctx.Done():
				return "", ctx.Err()
			}
		},
	}
	result := make(chan error, 1)
	go func() {
		_, err := runWorktreeCommand(t.Context(), holder, repoDir, "worktree", "prune")
		result <- err
	}()
	waitForAdminCommandStart(t, started)

	released := false
	releaseHolder := func() {
		t.Helper()
		if released {
			return
		}
		released = true
		close(release)
		if err := <-result; err != nil {
			t.Fatalf("release administration lock holder: %v", err)
		}
	}
	t.Cleanup(releaseHolder)
	return releaseHolder
}

// disposalLockRunner observes the real Git boundary without bounding Git's runtime.
type disposalLockRunner struct {
	resolved chan struct{}
	started  chan struct{}
}

func newDisposalLockRunner() *disposalLockRunner {
	return &disposalLockRunner{resolved: make(chan struct{}, 1), started: make(chan struct{}, 1)}
}

func (runner *disposalLockRunner) Run(ctx context.Context, workDir string, args ...string) (string, error) {
	if len(args) > 0 && args[0] == "worktree" {
		runner.started <- struct{}{}
	}
	output, err := (execGitRunner{}).Run(ctx, workDir, args...)
	if err == nil && slices.Equal(args, []string{"rev-parse", "--path-format=absolute", "--git-common-dir"}) {
		runner.resolved <- struct{}{}
	}
	return output, err
}

func requireNoDisposalWorktreeCommand(t *testing.T, started <-chan struct{}) {
	t.Helper()
	select {
	case <-started:
		t.Fatal("worktree command started while the administration lock was held")
	default:
	}
}

func requireDisposalWorktreeRegistration(t *testing.T, repoDir, path string, want bool) {
	t.Helper()
	registered := gitWorktreeTest(t, repoDir, "worktree", "list", "--porcelain")
	if got := strings.Contains(registered, "worktree "+canonicalPath(path)+"\n"); got != want {
		t.Fatalf("worktree %q registered = %v, want %v; list:\n%s", path, got, want, registered)
	}
}
