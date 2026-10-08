// Suite: legacy Run repository-key backfill.
// Invariant: a write open records only repository identity proven by surviving Git metadata.
// Boundary IN: the SQLite-backed store and local Git worktree administration files.
// Boundary OUT: repository-scoped CLI and cleanup consumers.
package store

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"roundfix/internal/gittest"
)

func TestOpenKeysALegacyRunFromItsRunWorktree(t *testing.T) {
	t.Parallel()
	fixture := newLegacyRepositoryKeyFixture(t)

	runStore := openTestStore(t, context.Background(), fixture.homeDir)
	defer closeStore(t, runStore)

	stored := legacyRepositoryKeyRun(t, runStore, fixture.run.ID)
	if stored.RepositoryRoot != fixture.mainRoot {
		t.Fatalf("repository root = %q, want %q", stored.RepositoryRoot, fixture.mainRoot)
	}
	runs, err := runStore.ListRuns(context.Background(), ListRunsQuery{
		GitRoot: fixture.mainRoot,
		States:  StatesAll,
	})
	if err != nil {
		t.Fatalf("list Runs from main checkout: %v", err)
	}
	if len(runs) != 1 || runs[0].ID != fixture.run.ID {
		t.Fatalf("Runs from main checkout = %#v, want Run %q", runs, fixture.run.ID)
	}
}

func TestOpenLeavesALegacyRunUnkeyedWhenItsWorktreeIsGone(t *testing.T) {
	t.Parallel()
	fixture := newLegacyRepositoryKeyFixture(t)
	gittest.Run(t, fixture.mainRoot, "worktree", "remove", fixture.workDir)

	runStore := openTestStore(t, context.Background(), fixture.homeDir)
	defer closeStore(t, runStore)

	stored := legacyRepositoryKeyRun(t, runStore, fixture.run.ID)
	if stored.RepositoryRoot != "" {
		t.Fatalf("repository root = %q, want empty for a gone Run Worktree", stored.RepositoryRoot)
	}
}

func TestOpenRejectsAWorktreeKeyWithoutACommonDirectory(t *testing.T) {
	t.Parallel()
	fixture := newLegacyRepositoryKeyFixture(t)
	gittest.Run(t, fixture.mainRoot, "worktree", "remove", fixture.workDir)
	if err := os.MkdirAll(fixture.workDir, 0o755); err != nil {
		t.Fatalf("recreate recorded Run Worktree directory: %v", err)
	}
	administrativeDir := filepath.Join(t.TempDir(), "gitdir-without-common-directory")
	if err := os.MkdirAll(administrativeDir, 0o755); err != nil {
		t.Fatalf("create incomplete worktree administration directory: %v", err)
	}
	if err := os.WriteFile(
		filepath.Join(fixture.workDir, ".git"),
		[]byte("gitdir: "+administrativeDir+"\n"),
		0o644,
	); err != nil {
		t.Fatalf("write incomplete Run Worktree Git pointer: %v", err)
	}

	runStore := openTestStore(t, context.Background(), fixture.homeDir)
	defer closeStore(t, runStore)

	stored := legacyRepositoryKeyRun(t, runStore, fixture.run.ID)
	if stored.RepositoryRoot != "" {
		t.Fatalf("repository root = %q, want empty without a common Git directory", stored.RepositoryRoot)
	}
}

func TestOpenKeysALegacyRunToItsOwnRepository(t *testing.T) {
	t.Parallel()
	fixture := newLegacyRepositoryKeyFixture(t)
	gittest.Run(t, fixture.mainRoot, "worktree", "remove", fixture.workDir)
	otherRoot := filepath.Join(t.TempDir(), "other-main")
	gittest.InitRepo(t, otherRoot, "--initial-branch=main")
	gittest.Run(t, otherRoot, "commit", "--allow-empty", "-m", "seed other repository")
	gittest.Run(t, otherRoot, "worktree", "add", "-b", "roundfix/legacy-run", fixture.workDir)
	otherRoot, err := filepath.EvalSymlinks(otherRoot)
	if err != nil {
		t.Fatalf("resolve other repository root: %v", err)
	}

	runStore := openTestStore(t, context.Background(), fixture.homeDir)
	defer closeStore(t, runStore)

	stored := legacyRepositoryKeyRun(t, runStore, fixture.run.ID)
	if stored.RepositoryRoot != otherRoot {
		t.Fatalf("repository root = %q, want owning repository %q", stored.RepositoryRoot, otherRoot)
	}
	runs, err := runStore.ListRuns(context.Background(), ListRunsQuery{
		GitRoot: fixture.mainRoot,
		States:  StatesAll,
	})
	if err != nil {
		t.Fatalf("list Runs from original repository: %v", err)
	}
	if len(runs) != 0 {
		t.Fatalf("Runs from original repository = %#v, want none", runs)
	}
}

func TestOpenRepositoryKeyBackfillIsIdempotent(t *testing.T) {
	t.Parallel()
	fixture := newLegacyRepositoryKeyFixture(t)

	first := openTestStore(t, context.Background(), fixture.homeDir)
	if got := connectionTotalChanges(t, first); got != 1 {
		closeStore(t, first)
		t.Fatalf("first write open changed %d rows, want 1", got)
	}
	closeStore(t, first)

	second := openTestStore(t, context.Background(), fixture.homeDir)
	defer closeStore(t, second)
	if got := connectionTotalChanges(t, second); got != 0 {
		t.Fatalf("second write open changed %d rows, want 0", got)
	}
	stored := legacyRepositoryKeyRun(t, second, fixture.run.ID)
	if stored.RepositoryRoot != fixture.mainRoot {
		t.Fatalf("repository root after second open = %q, want %q", stored.RepositoryRoot, fixture.mainRoot)
	}
}

func TestOpenReaderDoesNotBackfillRepositoryKeys(t *testing.T) {
	t.Parallel()
	fixture := newLegacyRepositoryKeyFixture(t)

	reader, err := OpenReader(context.Background(), fixture.homeDir)
	if err != nil {
		t.Fatalf("open Run Database reader: %v", err)
	}
	defer closeStore(t, reader)

	stored := legacyRepositoryKeyRun(t, reader, fixture.run.ID)
	if stored.RepositoryRoot != "" {
		t.Fatalf("reader repository root = %q, want empty without a write open", stored.RepositoryRoot)
	}
}

type legacyRepositoryKeyFixture struct {
	homeDir  string
	mainRoot string
	workDir  string
	run      Run
}

func newLegacyRepositoryKeyFixture(t *testing.T) legacyRepositoryKeyFixture {
	t.Helper()
	ctx := context.Background()
	fixtureRoot := t.TempDir()
	mainRoot, sourceRoot := linkedWorktreeFixture(t, fixtureRoot)
	homeDir := filepath.Join(fixtureRoot, "home")
	runStore := openTestStore(t, ctx, homeDir)

	request := sampleImplementCreateRunRequest()
	request.GitRoot = sourceRoot
	request.LocalBranch = "main"
	run, err := runStore.CreateRunSkippingActiveLock(ctx, request)
	if err != nil {
		closeStore(t, runStore)
		t.Fatalf("create legacy Run fixture: %v", err)
	}
	workDir := filepath.Join(fixtureRoot, "run-worktree")
	gittest.Run(t, mainRoot, "worktree", "add", "-b", RunBranchPrefix+run.ID, workDir)
	run, err = runStore.SetRunWorkDir(ctx, run.ID, workDir)
	if err != nil {
		closeStore(t, runStore)
		t.Fatalf("record legacy Run Worktree: %v", err)
	}
	if _, err := runStore.db.ExecContext(ctx,
		`UPDATE runs SET repository_root = '' WHERE id = ?`, run.ID,
	); err != nil {
		closeStore(t, runStore)
		t.Fatalf("clear legacy Run repository root: %v", err)
	}
	closeStore(t, runStore)
	gittest.Run(t, mainRoot, "worktree", "remove", sourceRoot)

	return legacyRepositoryKeyFixture{
		homeDir:  homeDir,
		mainRoot: mainRoot,
		workDir:  workDir,
		run:      run,
	}
}

func legacyRepositoryKeyRun(t *testing.T, runStore *Store, runID string) Run {
	t.Helper()
	run, found, err := runStore.Run(context.Background(), runID)
	if err != nil {
		t.Fatalf("read legacy Run %q: %v", runID, err)
	}
	if !found {
		t.Fatalf("legacy Run %q was not found", runID)
	}
	return run
}

func connectionTotalChanges(t *testing.T, runStore *Store) int {
	t.Helper()
	var changes int
	if err := runStore.db.QueryRowContext(context.Background(), `SELECT total_changes()`).Scan(&changes); err != nil {
		t.Fatalf("read connection total changes: %v", err)
	}
	return changes
}
