package store

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestOpenCreatesRunDatabaseAndAppliesMigrations(t *testing.T) {
	ctx := context.Background()
	homeDir := t.TempDir()

	store := openTestStore(t, ctx, homeDir)
	defer closeStore(t, store)

	if _, err := os.Stat(DatabasePath(homeDir)); err != nil {
		t.Fatalf("expected Run Database file at %s: %v", DatabasePath(homeDir), err)
	}
	version, err := store.MigrationVersion(ctx)
	if err != nil {
		t.Fatalf("expected migration version, got %v", err)
	}
	if version != 3 {
		t.Fatalf("expected migration version 3, got %d", version)
	}
}

func TestInteractiveDefaultsRememberLastPullRequestAndAgent(t *testing.T) {
	ctx := context.Background()
	store := openTestStore(t, ctx, t.TempDir())
	defer closeStore(t, store)

	if defaults, err := store.InteractiveDefaults(ctx); err != nil {
		t.Fatalf("read empty defaults: %v", err)
	} else if defaults.PRNumber != "" || defaults.Agent != "" {
		t.Fatalf("expected empty defaults, got %#v", defaults)
	}

	if err := store.RememberInteractiveDefaults(ctx, InteractiveDefaults{
		PRNumber: "123",
		Agent:    "codex",
	}); err != nil {
		t.Fatalf("remember defaults: %v", err)
	}
	if err := store.RememberInteractiveDefaults(ctx, InteractiveDefaults{
		PRNumber: "456",
	}); err != nil {
		t.Fatalf("update defaults: %v", err)
	}

	defaults, err := store.InteractiveDefaults(ctx)
	if err != nil {
		t.Fatalf("read defaults: %v", err)
	}
	if defaults.PRNumber != "456" {
		t.Fatalf("expected remembered PR 456, got %q", defaults.PRNumber)
	}
	if defaults.Agent != "codex" {
		t.Fatalf("expected remembered Agent codex, got %q", defaults.Agent)
	}
}

func TestCreateFetchRunCompletesFetchedAndReleasesLock(t *testing.T) {
	ctx := context.Background()
	store := openTestStore(t, ctx, t.TempDir())
	defer closeStore(t, store)

	run, err := store.CreateFetchRun(ctx, sampleCreateRunRequest())
	if err != nil {
		t.Fatalf("expected Fetch Run creation, got %v", err)
	}
	if run.ID == "" {
		t.Fatal("expected Run id")
	}
	if run.Kind != KindFetch {
		t.Fatalf("expected kind %q, got %q", KindFetch, run.Kind)
	}
	if run.State != StateActive {
		t.Fatalf("expected active state, got %q", run.State)
	}
	if run.CompletedAt != nil {
		t.Fatalf("expected active run without completion timestamp, got %v", run.CompletedAt)
	}

	active, found, err := store.ActiveRun(ctx, run.HeadRepository, run.HeadBranch)
	if err != nil {
		t.Fatalf("expected active run lookup, got %v", err)
	}
	if !found || active.ID != run.ID {
		t.Fatalf("expected active lock for %s, found=%v active=%#v", run.ID, found, active)
	}

	completed, err := store.CompleteRun(ctx, run.ID, StateFetched)
	if err != nil {
		t.Fatalf("expected Fetched completion, got %v", err)
	}
	if completed.State != StateFetched {
		t.Fatalf("expected Fetched state, got %q", completed.State)
	}
	if completed.CompletedAt == nil {
		t.Fatal("expected completion timestamp")
	}
	_, found, err = store.ActiveRun(ctx, run.HeadRepository, run.HeadBranch)
	if err != nil {
		t.Fatalf("expected active lookup after release, got %v", err)
	}
	if found {
		t.Fatal("expected terminal Fetched run to release Active Run lock")
	}

	second, err := store.CreateFetchRun(ctx, sampleCreateRunRequest())
	if err != nil {
		t.Fatalf("expected new Fetch Run after lock release, got %v", err)
	}
	if second.ID == run.ID {
		t.Fatal("expected second Run to have a distinct id")
	}
}

func TestCreateRunRejectsDuplicateActiveRunWithoutNewRecord(t *testing.T) {
	ctx := context.Background()
	store := openTestStore(t, ctx, t.TempDir())
	defer closeStore(t, store)

	first, err := store.CreateRun(ctx, sampleCreateRunRequest())
	if err != nil {
		t.Fatalf("expected first Run, got %v", err)
	}
	_, err = store.CreateRun(ctx, sampleCreateRunRequest())
	var activeErr ActiveRunError
	if !errors.As(err, &activeErr) {
		t.Fatalf("expected ActiveRunError, got %T %v", err, err)
	}
	if activeErr.Existing.ID != first.ID {
		t.Fatalf("expected existing run %s, got %s", first.ID, activeErr.Existing.ID)
	}
	if !strings.Contains(err.Error(), "existing run_id="+first.ID) {
		t.Fatalf("expected existing run id in error, got %q", err.Error())
	}
	count, err := store.RunCount(ctx)
	if err != nil {
		t.Fatalf("expected run count, got %v", err)
	}
	if count != 1 {
		t.Fatalf("expected duplicate active rejection to avoid new Run records, got count %d", count)
	}
}

func TestStoppedRunReleasesActiveLock(t *testing.T) {
	ctx := context.Background()
	store := openTestStore(t, ctx, t.TempDir())
	defer closeStore(t, store)

	req := sampleCreateRunRequest()
	req.Kind = KindResolve
	run, err := store.CreateRun(ctx, req)
	if err != nil {
		t.Fatalf("expected active Resolve Run, got %v", err)
	}
	if _, err := store.CompleteRun(ctx, run.ID, StateStopped); err != nil {
		t.Fatalf("expected Stopped completion, got %v", err)
	}
	_, found, err := store.ActiveRun(ctx, run.HeadRepository, run.HeadBranch)
	if err != nil {
		t.Fatalf("expected active lookup, got %v", err)
	}
	if found {
		t.Fatal("expected Stopped terminal outcome to release Active Run lock")
	}

	second, err := store.CreateRun(ctx, req)
	if err != nil {
		t.Fatalf("expected new Run after Stopped lock release, got %v", err)
	}
	if second.ID == run.ID {
		t.Fatal("expected distinct run id")
	}
}

func TestRunLooksUpExistingRunByID(t *testing.T) {
	ctx := context.Background()
	store := openTestStore(t, ctx, t.TempDir())
	defer closeStore(t, store)

	created, err := store.CreateRun(ctx, sampleCreateRunRequest())
	if err != nil {
		t.Fatalf("expected Run creation, got %v", err)
	}

	found, ok, err := store.Run(ctx, created.ID)
	if err != nil {
		t.Fatalf("lookup Run: %v", err)
	}
	if !ok || found.ID != created.ID {
		t.Fatalf("expected Run lookup for %s, ok=%v found=%#v", created.ID, ok, found)
	}

	_, ok, err = store.Run(ctx, "run_missing")
	if err != nil {
		t.Fatalf("lookup missing Run: %v", err)
	}
	if ok {
		t.Fatal("expected missing Run lookup")
	}
}

func TestCreateRunAllowsDifferentHeadBranch(t *testing.T) {
	ctx := context.Background()
	store := openTestStore(t, ctx, t.TempDir())
	defer closeStore(t, store)

	first, err := store.CreateRun(ctx, sampleCreateRunRequest())
	if err != nil {
		t.Fatalf("expected first Run, got %v", err)
	}
	secondReq := sampleCreateRunRequest()
	secondReq.HeadBranch = "feature/other-review"
	second, err := store.CreateRun(ctx, secondReq)
	if err != nil {
		t.Fatalf("expected simultaneous Run on different PR Head Branch, got %v", err)
	}
	if second.ID == first.ID {
		t.Fatal("expected distinct run ids")
	}
	count, err := store.RunCount(ctx)
	if err != nil {
		t.Fatalf("expected run count, got %v", err)
	}
	if count != 2 {
		t.Fatalf("expected two Run records, got %d", count)
	}
}

func TestCompleteRunAcceptsUnresolvedAsTerminal(t *testing.T) {
	ctx := context.Background()
	store := openTestStore(t, ctx, t.TempDir())
	defer closeStore(t, store)

	run, err := store.CreateRun(ctx, sampleCreateRunRequest())
	if err != nil {
		t.Fatalf("expected Run creation, got %v", err)
	}
	completed, err := store.CompleteRun(ctx, run.ID, StateUnresolved)
	if err != nil {
		t.Fatalf("expected Unresolved completion, got %v", err)
	}
	if completed.State != StateUnresolved || completed.CompletedAt == nil {
		t.Fatalf("expected completed Unresolved Run, got %+v", completed)
	}
	_, found, err := store.ActiveRun(ctx, run.HeadRepository, run.HeadBranch)
	if err != nil {
		t.Fatalf("expected active lookup, got %v", err)
	}
	if found {
		t.Fatal("expected Unresolved Run to release the Active Run lock")
	}
}

func TestCompleteRunRejectsNonTerminalState(t *testing.T) {
	ctx := context.Background()
	store := openTestStore(t, ctx, t.TempDir())
	defer closeStore(t, store)

	run, err := store.CreateRun(ctx, sampleCreateRunRequest())
	if err != nil {
		t.Fatalf("expected Run creation, got %v", err)
	}
	if _, err := store.CompleteRun(ctx, run.ID, StateActive); err == nil {
		t.Fatal("expected non-terminal completion to fail")
	}
	active, found, err := store.ActiveRun(ctx, run.HeadRepository, run.HeadBranch)
	if err != nil {
		t.Fatalf("expected active lookup, got %v", err)
	}
	if !found || active.ID != run.ID {
		t.Fatalf("expected active lock to remain after failed completion, found=%v active=%#v", found, active)
	}
}

func openTestStore(t *testing.T, ctx context.Context, homeDir string) *Store {
	t.Helper()
	store, err := Open(ctx, homeDir)
	if err != nil {
		t.Fatalf("open store: %v", err)
	}
	return store
}

func closeStore(t *testing.T, store *Store) {
	t.Helper()
	if err := store.Close(); err != nil {
		t.Fatalf("close store: %v", err)
	}
}

func sampleCreateRunRequest() CreateRunRequest {
	return CreateRunRequest{
		Kind:           KindFetch,
		HeadRepository: "owner/project",
		HeadBranch:     "feature/review",
		BaseRepository: "owner/project",
		PRNumber:       "123",
		GitRoot:        filepath.Join("tmp", "repo"),
		LocalBranch:    "feature/review",
		HeadSHA:        "abc123",
		ArtifactDir:    filepath.Join("tmp", "repo", ".roundfix"),
	}
}

func TestUpdateRunStateRejectsTerminalStatesAndMissingRuns(t *testing.T) {
	ctx := context.Background()
	store := openTestStore(t, ctx, t.TempDir())
	defer closeStore(t, store)

	run, err := store.CreateRun(ctx, sampleCreateRunRequest())
	if err != nil {
		t.Fatalf("create run: %v", err)
	}

	if err := store.UpdateRunState(ctx, run.ID, StateResolvingWithAgent); err != nil {
		t.Fatalf("expected intermediate state update, got %v", err)
	}
	updated, _, err := store.Run(ctx, run.ID)
	if err != nil {
		t.Fatalf("lookup run: %v", err)
	}
	if updated.State != StateResolvingWithAgent {
		t.Fatalf("expected ResolvingWithAgent, got %q", updated.State)
	}

	if err := store.UpdateRunState(ctx, run.ID, StateClean); err == nil {
		t.Fatal("expected terminal state rejection; terminal outcomes go through CompleteRun")
	}
	if err := store.UpdateRunState(ctx, "run_missing", StateVerifying); err == nil {
		t.Fatal("expected missing Run rejection")
	}
}
