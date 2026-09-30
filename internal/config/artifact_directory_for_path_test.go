// Suite: path-derived default Artifact Directories
// Invariant: deriving a default from a recorded checkout never re-resolves its repository identity.
// Boundary IN: Artifact Directory path derivation and real Git worktree identity.
// Boundary OUT: GC sanitation classification, owned by internal/cli/gc_sanitize_pre_key_root_test.go.
package config

import (
	"path/filepath"
	"testing"

	"roundfix/internal/gittest"
)

func TestDefaultArtifactDirectoryForPathKeepsTheCheckoutPath(t *testing.T) {
	t.Parallel()
	fixtureRoot, err := filepath.EvalSymlinks(t.TempDir())
	if err != nil {
		t.Fatalf("resolve fixture root: %v", err)
	}
	seedRoot := filepath.Join(fixtureRoot, "seed")
	bareRoot := filepath.Join(fixtureRoot, "repository.git")
	checkoutRoot := filepath.Join(fixtureRoot, "checkout")
	homeDir := filepath.Join(fixtureRoot, "home")
	gittest.InitRepo(t, seedRoot, "--initial-branch=main")
	gittest.Run(t, seedRoot, "commit", "--allow-empty", "-m", "seed repository")
	gittest.Run(t, seedRoot, "clone", "--bare", seedRoot, bareRoot)
	gittest.Harden(t, bareRoot)
	gittest.Run(t, bareRoot, "worktree", "add", "-b", "feature/checkout-default", checkoutRoot, "main")

	checkoutDefault, err := DefaultArtifactDirectoryForPath(checkoutRoot, homeDir)
	if err != nil {
		t.Fatalf("derive checkout default Artifact Directory: %v", err)
	}
	keyDefault, err := ResolveArtifactDirectory("", checkoutRoot, homeDir)
	if err != nil {
		t.Fatalf("resolve repository-key default Artifact Directory: %v", err)
	}

	wantCheckoutDefault := filepath.Join(homeDir, ".roundfix", "artifacts", repoID(checkoutRoot))
	if checkoutDefault != wantCheckoutDefault {
		t.Fatalf("checkout default = %q, want path-derived default %q", checkoutDefault, wantCheckoutDefault)
	}
	wantKeyDefault := filepath.Join(homeDir, ".roundfix", "artifacts", repoID(bareRoot))
	if keyDefault != wantKeyDefault {
		t.Fatalf("repository-key default = %q, want common Git directory default %q", keyDefault, wantKeyDefault)
	}
	if checkoutDefault == keyDefault {
		t.Fatalf("checkout default %q unexpectedly re-resolved to repository-key default", checkoutDefault)
	}

	if _, err := DefaultArtifactDirectoryForPath("", homeDir); err == nil {
		t.Fatal("empty path accepted for path-derived default Artifact Directory")
	}
	if _, err := DefaultArtifactDirectoryForPath(checkoutRoot, ""); err == nil {
		t.Fatal("empty Roundfix Home accepted for path-derived default Artifact Directory")
	}
}
