// Suite: runs list over terminal Runs from removed linked checkouts.
// Invariant: a vanished recorded checkout is a stored fact, not a CLI warning.
// Boundary IN: public CLI dispatch, the real Run Database, and real local Git repositories.
// Boundary OUT: retained-count validation details, owned by internal/worktree/retained_vanished_checkout_test.go.
package cli

import (
	"bytes"
	"context"
	"path/filepath"
	"strings"
	"testing"

	"roundfix/internal/store"
	runworktree "roundfix/internal/worktree"
)

func TestRunsListPrintsNoWarningForARunWhoseCheckoutWasDeleted(t *testing.T) {
	homeDir, repoDir := newImplementWorkspace(t, []implementSeed{{id: "task_01"}})
	linkedRoot := filepath.Join(t.TempDir(), "linked")
	const linkedBranch = "feature/removed-linked-runs-list"
	gitImplement(t, repoDir, "worktree", "add", "-b", linkedBranch, linkedRoot)
	linkedRoot, err := filepath.EvalSymlinks(linkedRoot)
	if err != nil {
		t.Fatalf("resolve linked checkout: %v", err)
	}
	run := createReconcileMetadataRun(t, homeDir, store.CreateRunRequest{
		Kind:        store.KindImplement,
		GitRoot:     linkedRoot,
		LocalBranch: linkedBranch,
		HeadSHA:     strings.TrimSpace(gitImplementOutput(t, linkedRoot, "rev-parse", "HEAD")),
		SpecSlug:    "runs-list-vanished-checkout",
		Agent:       "codex",
	}, store.StateStopped)
	gitImplement(t, repoDir, "branch", runworktree.BranchName(run.ID), run.HeadSHA)
	gitImplement(t, repoDir, "worktree", "remove", linkedRoot)
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := runCLIContext(
		t,
		context.Background(),
		[]string{"runs", "list", "--state", "terminal"},
		&stdout,
		&stderr,
	)

	if code != exitOK {
		t.Fatalf("runs list exit = %d, want %d; stderr=%q stdout=%q", code, exitOK, stderr.String(), stdout.String())
	}
	if !strings.Contains(stdout.String(), run.ID+"  ") {
		t.Fatalf("runs list stdout omitted Run %q: %q", run.ID, stdout.String())
	}
	if strings.Contains(stderr.String(), "warning:") {
		t.Fatalf("runs list reported a vanished-checkout warning: %q", stderr.String())
	}
	wantNote := "(1 terminal Run Worktree retained; run 'roundfix reconcile' to inspect)\n"
	if stderr.String() != wantNote {
		t.Fatalf("runs list stderr = %q, want retained guidance %q", stderr.String(), wantNote)
	}
}
