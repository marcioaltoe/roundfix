// Suite: retained terminal Runs when Git is unavailable.
// Invariant: a Git launch failure at an existing recorded root remains an inspection failure.
// Boundary IN: retained-Run counting and the Git runner boundary.
// Boundary OUT: public runs-list warning rendering, owned by internal/cli tests.
package worktree

import (
	"context"
	"io/fs"
	"os/exec"
	"testing"

	"roundfix/internal/store"
)

func TestCountRetainedTerminalRunsReportsAGitLaunchFailure(t *testing.T) {
	t.Parallel()
	root := canonicalPath(t.TempDir())
	runner := &retainedVanishedCheckoutGitRunner{run: func(_ string, _ []string) (string, error) {
		return "", &exec.Error{Name: "git", Err: fs.ErrNotExist}
	}}
	run := store.Run{
		ID:      "git-unavailable",
		Kind:    store.KindImplement,
		State:   store.StateStopped,
		GitRoot: root,
	}

	retained, failures := countRetainedTerminalRuns(context.Background(), runner, []store.Run{run})

	if retained != 0 {
		t.Fatalf("retained terminal Runs = %d, want 0", retained)
	}
	assertSingleRetainedRootFailure(t, failures, "git")
}
