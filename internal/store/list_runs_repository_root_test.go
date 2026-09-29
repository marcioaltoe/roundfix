// Suite: Run listing by durable repository identity.
// Invariant: a RepositoryRoot query returns every Run with that exact durable key and no Run from another repository.
// Boundary IN: Run persistence and ListRuns query construction.
// Boundary OUT: repository-root discovery, owned by internal/config/config_test.go.
package store

import (
	"context"
	"path/filepath"
	"testing"

	roundconfig "roundfix/internal/config"
)

func TestListRunsScopesByRepositoryRoot(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	fixtureRoot := t.TempDir()
	mainRoot, linkedRoot := linkedWorktreeFixture(t, fixtureRoot)
	runStore := openTestStore(t, ctx, filepath.Join(fixtureRoot, "home"))
	defer closeStore(t, runStore)

	linkedRequest := sampleCreateRunRequest()
	linkedRequest.GitRoot = linkedRoot
	linked, err := runStore.CreateRun(ctx, linkedRequest)
	if err != nil {
		t.Fatalf("create linked-worktree Run: %v", err)
	}
	otherRequest := sampleCreateRunRequest()
	otherRequest.GitRoot = filepath.Join(fixtureRoot, "other-repository")
	otherRequest.HeadBranch = "feature/other-repository"
	otherRequest.LocalBranch = "feature/other-repository"
	if _, err := runStore.CreateRun(ctx, otherRequest); err != nil {
		t.Fatalf("create other-repository Run: %v", err)
	}
	repositoryRoot, err := roundconfig.RepositoryRoot(mainRoot)
	if err != nil {
		t.Fatalf("resolve durable repository key: %v", err)
	}

	listed, err := runStore.ListRuns(ctx, ListRunsQuery{
		RepositoryRoot: repositoryRoot,
		States:         StatesAll,
	})
	if err != nil {
		t.Fatalf("list Runs by durable repository key: %v", err)
	}
	assertRunIDs(t, listed, []string{linked.ID})
}
