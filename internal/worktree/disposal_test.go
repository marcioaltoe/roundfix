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
	"strings"
	"testing"
	"time"
)

func TestRemoveRegisteredWorktreeWaitsForTheAdminLock(t *testing.T) {
	t.Parallel()
	repoDir, head := newStagingRepository(t)
	worktreePath := filepath.Join(t.TempDir(), "registered")
	gitWorktreeTest(t, repoDir, "worktree", "add", "--detach", worktreePath, head)
	release := holdDisposalAdminLock(t, repoDir)

	waitCtx, cancelWait := context.WithTimeout(t.Context(), 250*time.Millisecond)
	defer cancelWait()
	err := RemoveRegisteredWorktree(waitCtx, repoDir, worktreePath)
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("RemoveRegisteredWorktree() error = %v, want context deadline exceeded", err)
	}
	if _, err := os.Stat(worktreePath); err != nil {
		t.Fatalf("registered worktree was removed while lock was held: %v", err)
	}
	release()
}

func TestAddDetachedWorktreeWaitsForTheAdminLock(t *testing.T) {
	t.Parallel()
	repoDir, head := newStagingRepository(t)
	worktreePath := filepath.Join(t.TempDir(), "detached")
	release := holdDisposalAdminLock(t, repoDir)

	waitCtx, cancelWait := context.WithTimeout(t.Context(), 250*time.Millisecond)
	defer cancelWait()
	err := AddDetachedWorktree(waitCtx, repoDir, worktreePath, head)
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("AddDetachedWorktree() error = %v, want context deadline exceeded", err)
	}
	if _, err := os.Stat(worktreePath); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("detached worktree path exists while lock was held: %v", err)
	}
	release()
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
