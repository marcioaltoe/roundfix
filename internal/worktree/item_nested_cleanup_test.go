// Suite: half-removed Delivery Queue item cleanup.
// Invariant: cleanup never deletes a registered worktree nested under an item path.
// Boundary IN: CleanupItem, real Git worktree registration, and item-path removal.
// Boundary OUT: Delivery Queue warning persistence, owned by internal/cli delivery tests.
package worktree

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestCleanupItemRefusesAHalfRemovedItemHoldingARegisteredWorktree(t *testing.T) {
	repoDir, ref, head := newHalfRemovedItemForNestedCleanupTest(t, "nested")
	nestedPath := filepath.Join(ref.Path, "nested-worktree")
	gitWorktreeTest(t, repoDir, "worktree", "add", "--detach", nestedPath, head)
	uncommittedPath := filepath.Join(nestedPath, "uncommitted.txt")
	mustWriteWorktreeTest(t, uncommittedPath, "preserve nested work\n")

	err := CleanupItem(t.Context(), ref)

	if err == nil || !strings.Contains(err.Error(), nestedPath) {
		t.Fatalf("CleanupItem() error = %v, want refusal naming nested worktree %q", err, nestedPath)
	}
	if got := mustReadWorktreeTest(t, uncommittedPath); got != "preserve nested work\n" {
		t.Fatalf("nested uncommitted file = %q, want preserved content", got)
	}
	registered := gitWorktreeTest(t, repoDir, "worktree", "list", "--porcelain")
	if !strings.Contains(registered, "worktree "+canonicalPath(nestedPath)) {
		t.Fatalf("registered worktrees = %q, want nested path %q", registered, nestedPath)
	}
}

func TestCleanupItemRemovesAHalfRemovedItemWithoutANestedWorktree(t *testing.T) {
	repoDir, ref, _ := newHalfRemovedItemForNestedCleanupTest(t, "plain")

	if err := CleanupItem(t.Context(), ref); err != nil {
		t.Fatalf("CleanupItem() error = %v", err)
	}

	assertPathRemoved(t, ref.Path)
	assertBranchRemoved(t, repoDir, ref.Branch)
}

func newHalfRemovedItemForNestedCleanupTest(t *testing.T, suffix string) (string, ItemRef, string) {
	t.Helper()
	repoDir := initWorktreeRepo(t)
	mustWriteWorktreeTest(t, filepath.Join(repoDir, "tracked.txt"), "base\n")
	gitWorktreeTest(t, repoDir, "add", "tracked.txt")
	gitWorktreeTest(t, repoDir, "commit", "-m", "initial")
	head := strings.TrimSpace(gitWorktreeTest(t, repoDir, "rev-parse", "HEAD"))
	ref, err := ItemRefFor(
		repoDir,
		filepath.Join(t.TempDir(), "worktrees"),
		"roundfix/deliver-half-removed-"+suffix,
	)
	if err != nil {
		t.Fatalf("derive item Worktree ref: %v", err)
	}
	if err := CreateItem(t.Context(), ref, ItemCreateOptions{HeadSHA: head}); err != nil {
		t.Fatalf("create item Worktree: %v", err)
	}

	marker := strings.TrimSpace(mustReadWorktreeTest(t, filepath.Join(ref.Path, ".git")))
	adminDir := strings.TrimSpace(strings.TrimPrefix(marker, "gitdir:"))
	if err := os.RemoveAll(adminDir); err != nil {
		t.Fatalf("remove item Worktree admin directory %q: %v", adminDir, err)
	}
	if _, err := os.Lstat(adminDir); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("item Worktree admin directory %q remains: %v", adminDir, err)
	}
	registered := gitWorktreeTest(t, repoDir, "worktree", "list", "--porcelain")
	if strings.Contains(registered, "worktree "+ref.Path) {
		t.Fatalf("half-removed item Worktree remained registered: %s", registered)
	}
	return repoDir, ref, head
}
