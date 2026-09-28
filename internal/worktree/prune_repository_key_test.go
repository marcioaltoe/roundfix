// Suite: implement-preflight terminal Run pruning by repository key.
// Invariant: pruning inspects only Runs owned by the current repository, even after a checkout disappears.
// Boundary IN: PruneTerminalReport, real Git worktrees, and persisted Run metadata.
// Boundary OUT: Run Database key derivation and CLI reconciliation.
package worktree

import (
	"context"
	"path/filepath"
	"testing"

	"roundfix/internal/store"
)

func TestPruneTerminalReportReleasesARunWhoseCheckoutWasRemoved(t *testing.T) {
	fixture := newTerminalRunFixture(t, "prune-removed-checkout")
	userRoot := canonicalPath(fixture.repoDir)
	fixture.run.GitRoot = filepath.Join(t.TempDir(), "removed-checkout")
	fixture.run.RepositoryRoot = userRoot
	location := filepath.Dir(filepath.Dir(fixture.ref.Path))

	pruned, err := PruneTerminalReport(
		context.Background(),
		userRoot,
		location,
		&recordingTerminalRunStore{},
		func(_ context.Context, runID string) (store.Run, bool, error) {
			return fixture.run, runID == fixture.run.ID, nil
		},
	)
	if err != nil {
		t.Fatalf("prune terminal Run whose checkout was removed: %v", err)
	}
	if len(pruned) != 1 || pruned[0].RunID != fixture.run.ID {
		t.Fatalf("pruned refs = %#v, want Run %q", pruned, fixture.run.ID)
	}
	assertPathRemoved(t, fixture.ref.Path)
	assertBranchRemoved(t, fixture.repoDir, fixture.ref.Branch)
}

func TestPruneTerminalReportSkipsARunKeyedToAnotherRepository(t *testing.T) {
	fixture := newTerminalRunFixture(t, "prune-other-repository")
	fixture.run.RepositoryRoot = canonicalPath(initWorktreeRepo(t))
	location := filepath.Dir(filepath.Dir(fixture.ref.Path))

	pruned, err := PruneTerminalReport(
		context.Background(),
		fixture.repoDir,
		location,
		&recordingTerminalRunStore{},
		func(_ context.Context, runID string) (store.Run, bool, error) {
			return fixture.run, runID == fixture.run.ID, nil
		},
	)
	if err != nil {
		t.Fatalf("prune terminal Run keyed to another repository: %v", err)
	}
	if len(pruned) != 0 {
		t.Fatalf("pruned refs = %#v, want none for another repository", pruned)
	}
	assertPathExists(t, fixture.ref.Path)
	assertRunBranchExists(t, fixture.repoDir, fixture.ref.Branch)
}
