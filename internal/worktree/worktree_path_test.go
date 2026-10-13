// Suite: worktree path comparison through a symbolic-linked location.
// Invariant: missing paths retain their resolved prefix and match Git registrations.
// Boundary IN: canonicalPath, CreateItem, CleanupItem and UseItem with real Git.
// Boundary OUT: terminal Run containment and CLI/daemon consumer flows.
package worktree

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestCanonicalPathResolvesTheExistingPrefix(t *testing.T) {
	t.Parallel()
	realRoot, err := filepath.EvalSymlinks(t.TempDir())
	if err != nil {
		t.Fatalf("resolve real root: %v", err)
	}
	linkedRoot := filepath.Join(t.TempDir(), "linked-root")
	if err := os.Symlink(realRoot, linkedRoot); err != nil {
		t.Fatalf("create root link: %v", err)
	}
	if err := os.Mkdir(filepath.Join(realRoot, "existing"), 0o755); err != nil {
		t.Fatalf("create existing directory: %v", err)
	}
	// No named component of this relative path exists; its clean spelling stays.
	missingRelative := filepath.Join("missing-"+filepath.Base(t.TempDir()), "parent", "leaf")
	if _, err := os.Stat(missingRelative); !os.IsNotExist(err) {
		t.Fatalf("relative fixture must be absent: %v", err)
	}
	for _, tc := range []struct {
		name, path, want string
	}{
		{"existing", filepath.Join(linkedRoot, "existing"), filepath.Join(realRoot, "existing")},
		{"missing-leaf", filepath.Join(linkedRoot, "existing", "Missing Leaf"), filepath.Join(realRoot, "existing", "Missing Leaf")},
		{"missing-parent", filepath.Join(linkedRoot, "Missing Parent", "Missing Leaf"), filepath.Join(realRoot, "Missing Parent", "Missing Leaf")},
		{"no-existing-prefix", missingRelative, filepath.Clean(missingRelative)},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			if got := canonicalPath(tc.path); got != tc.want {
				t.Fatalf("canonicalPath(%q) = %q, want %q", tc.path, got, tc.want)
			}
		})
	}
}

func TestCleanupItemReleasesAGoneItemWorktreeUnderASymlinkedRoot(t *testing.T) {
	t.Parallel()
	ref, resolvedPath := newSymlinkedItemWorktree(t)
	if err := os.RemoveAll(ref.Path); err != nil {
		t.Fatalf("remove item worktree directory: %v", err)
	}
	assertPathRemoved(t, ref.Path)
	if err := CleanupItem(t.Context(), ref); err != nil {
		t.Fatalf("clean up gone item worktree: %v", err)
	}
	registered := gitWorktreeTest(t, ref.UserRoot, "worktree", "list", "--porcelain")
	if strings.Contains(registered, "worktree "+resolvedPath+"\n") || strings.Contains(registered, "branch refs/heads/"+ref.Branch+"\n") {
		t.Fatalf("item registration remains after cleanup:\n%s", registered)
	}
	assertBranchRemoved(t, ref.UserRoot, ref.Branch)
	assertPathRemoved(t, ref.Path)
}

func TestUseItemRecreatesAGoneItemWorktreeUnderASymlinkedRoot(t *testing.T) {
	t.Parallel()
	ref, resolvedPath := newSymlinkedItemWorktree(t)
	head := commitWorktreeFile(t, ref.Path, "item.txt", "item work\n", "item change")
	if err := os.RemoveAll(ref.Path); err != nil {
		t.Fatalf("remove item worktree directory: %v", err)
	}
	assertPathRemoved(t, ref.Path)
	if err := UseItem(t.Context(), ref); err != nil {
		t.Fatalf("recreate gone item worktree: %v", err)
	}
	assertPathExists(t, ref.Path)
	if got := strings.TrimSpace(gitWorktreeTest(t, ref.Path, "branch", "--show-current")); got != ref.Branch {
		t.Fatalf("recreated branch = %q, want %q", got, ref.Branch)
	}
	if got := strings.TrimSpace(gitWorktreeTest(t, ref.Path, "rev-parse", "HEAD")); got != head {
		t.Fatalf("recreated HEAD = %q, want %q", got, head)
	}
	if got := mustReadWorktreeTest(t, filepath.Join(ref.Path, "item.txt")); got != "item work\n" {
		t.Fatalf("recreated item content = %q", got)
	}
	registered := gitWorktreeTest(t, ref.UserRoot, "worktree", "list", "--porcelain")
	if strings.Count(registered, "worktree "+resolvedPath+"\n") != 1 || strings.Count(registered, "branch refs/heads/"+ref.Branch+"\n") != 1 {
		t.Fatalf("expected one recreated item registration:\n%s", registered)
	}
}

func newSymlinkedItemWorktree(t *testing.T) (ItemRef, string) {
	t.Helper()
	repoDir := initWorktreeRepo(t)
	head := commitWorktreeFile(t, repoDir, "base.txt", "base\n", "initial")
	realLocation := t.TempDir()
	location := filepath.Join(t.TempDir(), "linked-worktrees")
	if err := os.Symlink(realLocation, location); err != nil {
		t.Fatalf("create Worktree Root link: %v", err)
	}
	ref, err := ItemRefFor(repoDir, location, "roundfix/deliver-linked-item-0123456789abcdef")
	if err != nil {
		t.Fatalf("derive item worktree: %v", err)
	}
	if err := CreateItem(t.Context(), ref, ItemCreateOptions{HeadSHA: head}); err != nil {
		t.Fatalf("create item worktree: %v", err)
	}
	resolvedPath, err := filepath.EvalSymlinks(ref.Path)
	if err != nil {
		t.Fatalf("resolve item path before removing it: %v", err)
	}
	if ref.Path == resolvedPath || filepath.Dir(filepath.Dir(ref.Path)) != location {
		t.Fatalf("item path lost linked location spelling: %q", ref.Path)
	}
	registered := gitWorktreeTest(t, repoDir, "worktree", "list", "--porcelain")
	if !strings.Contains(registered, "worktree "+resolvedPath+"\n") || !strings.Contains(registered, "branch refs/heads/"+ref.Branch+"\n") {
		t.Fatalf("expected resolved item registration:\n%s", registered)
	}
	return ref, resolvedPath
}
