package store

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"roundfix/internal/runevent"
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
	if version != 9 {
		t.Fatalf("expected migration version 9, got %d", version)
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

func TestListRunsScopesByRepositoryAndOrdersNewestFirst(t *testing.T) {
	ctx := context.Background()
	runStore := openTestStore(t, ctx, t.TempDir())
	defer closeStore(t, runStore)

	base := time.Date(2026, 7, 6, 10, 0, 0, 0, time.UTC)
	repoA := filepath.Join("tmp", "repo-a")
	repoB := filepath.Join("tmp", "repo-b")
	repoAOlder := createListedRun(t, ctx, runStore, listedRunSeed{
		gitRoot:   repoA,
		branch:    "feature/repo-a-older",
		prNumber:  "101",
		createdAt: base,
		state:     StateActive,
	})
	repoANewer := createListedRun(t, ctx, runStore, listedRunSeed{
		gitRoot:   repoA,
		branch:    "feature/repo-a-newer",
		prNumber:  "102",
		createdAt: base.Add(2 * time.Minute),
		state:     StateResolvingWithAgent,
	})
	repoBRun := createListedRun(t, ctx, runStore, listedRunSeed{
		gitRoot:   repoB,
		branch:    "feature/repo-b",
		prNumber:  "201",
		createdAt: base.Add(4 * time.Minute),
		state:     StateActive,
	})

	scoped, err := runStore.ListRuns(ctx, ListRunsQuery{GitRoot: repoA})
	if err != nil {
		t.Fatalf("list repository-scoped Runs: %v", err)
	}
	assertRunIDs(t, scoped, []string{repoANewer.ID, repoAOlder.ID})
	for _, run := range scoped {
		if run.GitRoot != repoA {
			t.Fatalf("expected only repository %q, got Run %s in %q", repoA, run.ID, run.GitRoot)
		}
	}

	allRepos, err := runStore.ListRuns(ctx, ListRunsQuery{})
	if err != nil {
		t.Fatalf("list Runs across repositories: %v", err)
	}
	assertRunIDs(t, allRepos, []string{repoBRun.ID, repoANewer.ID, repoAOlder.ID})
}

func TestListRunsStateFilterAndLimit(t *testing.T) {
	ctx := context.Background()
	runStore := openTestStore(t, ctx, t.TempDir())
	defer closeStore(t, runStore)

	gitRoot := filepath.Join("tmp", "state-filter")
	base := time.Date(2026, 7, 6, 11, 0, 0, 0, time.UTC)
	// Every Run state, interleaving non-terminal and terminal, oldest first.
	seededStates := []string{
		StateActive,
		StateFetched,
		StateResolvingWithAgent,
		StateStopped,
		StateVerifying,
		StateClean,
		StateCleanUnverified,
		StatePushing,
		StateMaxRoundsReached,
		StateBudgetExceeded,
		StateTimedOut,
		StateFailed,
		StateIntegrationPending,
		StateUnresolved,
	}
	newestFirstAll := make([]string, 0, len(seededStates))
	newestFirstActive := []string{}
	newestFirstTerminal := []string{}
	for index, state := range seededStates {
		run := createListedRun(t, ctx, runStore, listedRunSeed{
			gitRoot:   gitRoot,
			branch:    fmt.Sprintf("feature/state-%02d", index),
			prNumber:  fmt.Sprintf("3%02d", index),
			createdAt: base.Add(time.Duration(index*2) * time.Minute),
			state:     state,
		})
		newestFirstAll = append([]string{run.ID}, newestFirstAll...)
		if IsTerminalState(state) {
			newestFirstTerminal = append([]string{run.ID}, newestFirstTerminal...)
		} else {
			newestFirstActive = append([]string{run.ID}, newestFirstActive...)
		}
	}
	activeOnlyRoot := filepath.Join("tmp", "state-filter-active-only")
	createListedRun(t, ctx, runStore, listedRunSeed{
		gitRoot:   activeOnlyRoot,
		branch:    "feature/lone-active",
		prNumber:  "400",
		createdAt: base.Add(time.Hour),
		state:     StateActive,
	})

	cases := []struct {
		name  string
		query ListRunsQuery
		want  []string
	}{
		{
			name:  "unset filter defaults to active",
			query: ListRunsQuery{GitRoot: gitRoot},
			want:  newestFirstActive,
		},
		{
			name:  "active filter keeps only non-terminal Runs",
			query: ListRunsQuery{GitRoot: gitRoot, States: StatesActive},
			want:  newestFirstActive,
		},
		{
			name:  "terminal filter keeps only terminal Runs",
			query: ListRunsQuery{GitRoot: gitRoot, States: StatesTerminal},
			want:  newestFirstTerminal,
		},
		{
			name:  "all filter with zero limit keeps every Run",
			query: ListRunsQuery{GitRoot: gitRoot, States: StatesAll},
			want:  newestFirstAll,
		},
		{
			name:  "limit bounds the newest Runs",
			query: ListRunsQuery{GitRoot: gitRoot, States: StatesAll, Limit: 3},
			want:  newestFirstAll[:3],
		},
		{
			name:  "limit applies after the state filter",
			query: ListRunsQuery{GitRoot: gitRoot, States: StatesTerminal, Limit: 2},
			want:  newestFirstTerminal[:2],
		},
		{
			name:  "limit above the match count keeps every match",
			query: ListRunsQuery{GitRoot: gitRoot, States: StatesActive, Limit: len(newestFirstAll) + 1},
			want:  newestFirstActive,
		},
		{
			name:  "filter matching nothing returns an empty slice",
			query: ListRunsQuery{GitRoot: activeOnlyRoot, States: StatesTerminal},
			want:  []string{},
		},
	}
	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			runs, err := runStore.ListRuns(ctx, testCase.query)
			if err != nil {
				t.Fatalf("list Runs: %v", err)
			}
			if runs == nil {
				t.Fatal("expected empty slice, got nil")
			}
			assertRunIDs(t, runs, testCase.want)
		})
	}
}

func TestListRunsEmptyDatabaseReturnsEmptySlice(t *testing.T) {
	ctx := context.Background()
	runStore := openTestStore(t, ctx, t.TempDir())
	defer closeStore(t, runStore)

	runs, err := runStore.ListRuns(ctx, ListRunsQuery{})
	if err != nil {
		t.Fatalf("list empty Run Database: %v", err)
	}
	if runs == nil {
		t.Fatal("expected empty slice, got nil")
	}
	if len(runs) != 0 {
		t.Fatalf("expected no Runs, got %#v", runs)
	}
}

func TestCreateRunPersistsWorkDir(t *testing.T) {
	ctx := context.Background()
	store := openTestStore(t, ctx, t.TempDir())
	defer closeStore(t, store)

	req := sampleCreateRunRequest()
	req.WorkDir = filepath.Join("tmp", "roundfix-home", "worktrees", "repo-id", "run-id")
	created, err := store.CreateRun(ctx, req)
	if err != nil {
		t.Fatalf("expected Run creation, got %v", err)
	}
	if created.WorkDir != req.WorkDir {
		t.Fatalf("expected created WorkDir %q, got %q", req.WorkDir, created.WorkDir)
	}

	found, ok, err := store.Run(ctx, created.ID)
	if err != nil || !ok {
		t.Fatalf("lookup persisted Run: ok=%v err=%v", ok, err)
	}
	if found.WorkDir != req.WorkDir {
		t.Fatalf("expected looked-up WorkDir %q, got %q", req.WorkDir, found.WorkDir)
	}

	active, ok, err := store.ActiveRun(ctx, req.HeadRepository, req.HeadBranch)
	if err != nil || !ok {
		t.Fatalf("lookup active Run: ok=%v err=%v", ok, err)
	}
	if active.WorkDir != req.WorkDir {
		t.Fatalf("expected active WorkDir %q, got %q", req.WorkDir, active.WorkDir)
	}

	inGitRoot, ok, err := store.ActiveRunInGitRoot(ctx, req.GitRoot)
	if err != nil || !ok {
		t.Fatalf("lookup active Run in Git root: ok=%v err=%v", ok, err)
	}
	if inGitRoot.WorkDir != req.WorkDir {
		t.Fatalf("expected Git-root active WorkDir %q, got %q", req.WorkDir, inGitRoot.WorkDir)
	}
}

func TestCreateRunPersistsAgentSelectionAcrossRunQueries(t *testing.T) {
	ctx := context.Background()
	runStore := openTestStore(t, ctx, t.TempDir())
	defer closeStore(t, runStore)

	req := sampleCreateRunRequest()
	req.Kind = KindResolve
	req.Agent = "codex"
	req.Model = "gpt-5.6-sol"
	req.ReasoningEffort = "experimental-reasoning"
	created, err := runStore.CreateRun(ctx, req)
	if err != nil {
		t.Fatalf("expected Resolve Run creation, got %v", err)
	}
	assertRunSelection(t, created, "gpt-5.6-sol", "experimental-reasoning")

	found, ok, err := runStore.Run(ctx, created.ID)
	if err != nil || !ok {
		t.Fatalf("lookup persisted Run: ok=%v err=%v", ok, err)
	}
	assertRunSelection(t, found, "gpt-5.6-sol", "experimental-reasoning")

	active, ok, err := runStore.ActiveRun(ctx, req.HeadRepository, req.HeadBranch)
	if err != nil || !ok {
		t.Fatalf("lookup active Run: ok=%v err=%v", ok, err)
	}
	assertRunSelection(t, active, "gpt-5.6-sol", "experimental-reasoning")

	inGitRoot, ok, err := runStore.ActiveRunInGitRoot(ctx, req.GitRoot)
	if err != nil || !ok {
		t.Fatalf("lookup active Run in Git root: ok=%v err=%v", ok, err)
	}
	assertRunSelection(t, inGitRoot, "gpt-5.6-sol", "experimental-reasoning")

	listed, err := runStore.ListRuns(ctx, ListRunsQuery{GitRoot: req.GitRoot})
	if err != nil {
		t.Fatalf("list Runs: %v", err)
	}
	if len(listed) != 1 {
		t.Fatalf("expected one listed Run, got %#v", listed)
	}
	assertRunSelection(t, listed[0], "gpt-5.6-sol", "experimental-reasoning")
}

func TestCreateRunPersistsOwnerPIDAcrossRunQueries(t *testing.T) {
	ctx := context.Background()
	runStore := openTestStore(t, ctx, t.TempDir())
	defer closeStore(t, runStore)

	wantPID := os.Getpid()
	req := sampleCreateRunRequest()
	req.OwnerPID = wantPID
	created, err := runStore.CreateRun(ctx, req)
	if err != nil {
		t.Fatalf("expected Run creation, got %v", err)
	}
	assertRunOwnerPID(t, created, wantPID)

	found, ok, err := runStore.Run(ctx, created.ID)
	if err != nil || !ok {
		t.Fatalf("lookup persisted Run: ok=%v err=%v", ok, err)
	}
	assertRunOwnerPID(t, found, wantPID)

	active, ok, err := runStore.ActiveRun(ctx, req.HeadRepository, req.HeadBranch)
	if err != nil || !ok {
		t.Fatalf("lookup active Run: ok=%v err=%v", ok, err)
	}
	assertRunOwnerPID(t, active, wantPID)

	inGitRoot, ok, err := runStore.ActiveRunInGitRoot(ctx, req.GitRoot)
	if err != nil || !ok {
		t.Fatalf("lookup active Run in Git root: ok=%v err=%v", ok, err)
	}
	assertRunOwnerPID(t, inGitRoot, wantPID)

	listed, err := runStore.ListRuns(ctx, ListRunsQuery{GitRoot: req.GitRoot})
	if err != nil {
		t.Fatalf("list Runs: %v", err)
	}
	if len(listed) != 1 {
		t.Fatalf("expected one listed Run, got %#v", listed)
	}
	assertRunOwnerPID(t, listed[0], wantPID)

	var rawOwnerPID int
	if err := runStore.db.QueryRowContext(ctx, `SELECT owner_pid FROM runs WHERE id = ?`, created.ID).Scan(&rawOwnerPID); err != nil {
		t.Fatalf("read owner_pid: %v", err)
	}
	if rawOwnerPID != wantPID {
		t.Fatalf("expected raw owner_pid %d, got %d", wantPID, rawOwnerPID)
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

func TestCompleteRunWinnerAndIdenticalReplay(t *testing.T) {
	ctx := context.Background()
	runStore := openTestStore(t, ctx, t.TempDir())
	defer closeStore(t, runStore)

	createdAt := time.Date(2026, 7, 27, 9, 0, 0, 0, time.UTC)
	completedAt := createdAt.Add(time.Minute)
	runStore.now = func() time.Time { return createdAt }
	run, err := runStore.CreateRun(ctx, sampleCreateRunRequest())
	if err != nil {
		t.Fatalf("create Run: %v", err)
	}
	if _, err := runStore.AppendRunEvent(ctx, sampleRunEvent(run.ID, "before completion")); err != nil {
		t.Fatalf("append initial Run Event: %v", err)
	}

	runStore.now = func() time.Time { return completedAt }
	first, err := runStore.CompleteRun(ctx, run.ID, StateStopped)
	if err != nil {
		t.Fatalf("complete Run: %v", err)
	}
	if !first.Transitioned {
		t.Fatal("expected first completion to report a transition")
	}
	if first.Run.State != StateStopped || first.Run.CompletedAt == nil {
		t.Fatalf("expected terminal Stopped Run, got %+v", first.Run)
	}
	if got := countActiveRunLocks(t, ctx, runStore); got != 0 {
		t.Fatalf("expected winning completion to release one Active Run lock, got %d", got)
	}

	runStore.now = func() time.Time { return completedAt.Add(time.Hour) }
	replay, err := runStore.CompleteRun(ctx, run.ID, StateStopped)
	if err != nil {
		t.Fatalf("replay identical completion: %v", err)
	}
	if replay.Transitioned {
		t.Fatal("expected identical replay to report no transition")
	}
	if replay.Run.State != first.Run.State ||
		!replay.Run.UpdatedAt.Equal(first.Run.UpdatedAt) ||
		replay.Run.CompletedAt == nil ||
		!replay.Run.CompletedAt.Equal(*first.Run.CompletedAt) {
		t.Fatalf("expected replay to preserve the terminal row, first=%+v replay=%+v", first.Run, replay.Run)
	}
	if got := countActiveRunLocks(t, ctx, runStore); got != 0 {
		t.Fatalf("expected replay to preserve released lock state, got %d locks", got)
	}
	if got := countRunEvents(t, ctx, runStore, run.ID); got != 1 {
		t.Fatalf("expected replay to preserve journal history, got %d events", got)
	}
}

func TestTerminalOutcomeConflictPreservesWinner(t *testing.T) {
	ctx := context.Background()
	runStore := openTestStore(t, ctx, t.TempDir())
	defer closeStore(t, runStore)

	completedAt := time.Date(2026, 7, 27, 10, 0, 0, 0, time.UTC)
	run, err := runStore.CreateRun(ctx, sampleCreateRunRequest())
	if err != nil {
		t.Fatalf("create Run: %v", err)
	}
	if _, err := runStore.AppendRunEvent(ctx, sampleRunEvent(run.ID, "winner evidence")); err != nil {
		t.Fatalf("append initial Run Event: %v", err)
	}
	runStore.now = func() time.Time { return completedAt }
	winner, err := runStore.CompleteRun(ctx, run.ID, StateStopped)
	if err != nil {
		t.Fatalf("complete winning outcome: %v", err)
	}

	runStore.now = func() time.Time { return completedAt.Add(time.Hour) }
	_, err = runStore.CompleteRun(ctx, run.ID, StateFailed)
	var conflict TerminalOutcomeConflictError
	if !errors.As(err, &conflict) {
		t.Fatalf("expected TerminalOutcomeConflictError, got %T %v", err, err)
	}
	if conflict.RunID != run.ID || conflict.Stored != StateStopped || conflict.Requested != StateFailed {
		t.Fatalf("unexpected terminal conflict: %+v", conflict)
	}

	stored, found, err := runStore.Run(ctx, run.ID)
	if err != nil || !found {
		t.Fatalf("read terminal winner: found=%v err=%v", found, err)
	}
	if stored.State != winner.Run.State ||
		!stored.UpdatedAt.Equal(winner.Run.UpdatedAt) ||
		stored.CompletedAt == nil ||
		!stored.CompletedAt.Equal(*winner.Run.CompletedAt) {
		t.Fatalf("expected conflict to preserve winner, winner=%+v stored=%+v", winner.Run, stored)
	}
	if got := countActiveRunLocks(t, ctx, runStore); got != 0 {
		t.Fatalf("expected conflict to preserve released lock state, got %d locks", got)
	}
	if got := countRunEvents(t, ctx, runStore, run.ID); got != 1 {
		t.Fatalf("expected conflict to preserve journal history, got %d events", got)
	}
}

func TestTerminalOutcomeRejectsIntermediateStateUpdate(t *testing.T) {
	ctx := context.Background()
	runStore := openTestStore(t, ctx, t.TempDir())
	defer closeStore(t, runStore)
	run, err := runStore.CreateRun(ctx, sampleCreateRunRequest())
	if err != nil {
		t.Fatalf("create Run: %v", err)
	}
	winner, err := runStore.CompleteRun(ctx, run.ID, StateStopped)
	if err != nil {
		t.Fatalf("complete Run: %v", err)
	}

	err = runStore.UpdateRunState(ctx, run.ID, StateResolvingWithAgent)
	var conflict TerminalOutcomeConflictError
	if !errors.As(err, &conflict) {
		t.Fatalf("expected TerminalOutcomeConflictError, got %T %v", err, err)
	}
	if conflict.Stored != StateStopped || conflict.Requested != StateResolvingWithAgent {
		t.Fatalf("unexpected conflict: %+v", conflict)
	}
	stored, found, err := runStore.Run(ctx, run.ID)
	if err != nil || !found {
		t.Fatalf("read stored Run: found=%v err=%v", found, err)
	}
	if stored.State != StateStopped || !stored.CompletedAt.Equal(*winner.CompletedAt) {
		t.Fatalf("intermediate update changed terminal winner: before=%+v after=%+v", winner.Run, stored)
	}
}

func TestTerminalOutcomeEveryStoredTerminalStateIsImmutable(t *testing.T) {
	ctx := context.Background()
	runStore := openTestStore(t, ctx, t.TempDir())
	defer closeStore(t, runStore)

	sourceOutcomes := []string{
		StateFetched,
		StateStopped,
		StateClean,
		StateCleanUnverified,
		StateMaxRoundsReached,
		StateBudgetExceeded,
		StateTimedOut,
		StateFailed,
		StateIntegrationPending,
		StateUnresolved,
	}
	for index, sourceOutcome := range sourceOutcomes {
		t.Run(sourceOutcome, func(t *testing.T) {
			req := sampleCreateRunRequest()
			req.HeadBranch = fmt.Sprintf("feature/terminal-immutable-%d", index)
			req.PRNumber = fmt.Sprintf("terminal-immutable-%d", index)
			run, err := runStore.CreateRun(ctx, req)
			if err != nil {
				t.Fatalf("create Run: %v", err)
			}
			first, err := runStore.CompleteRun(ctx, run.ID, sourceOutcome)
			if err != nil {
				t.Fatalf("complete Run %s: %v", sourceOutcome, err)
			}
			requested := StateStopped
			if sourceOutcome == StateStopped {
				requested = StateFailed
			}

			_, err = runStore.CompleteRun(ctx, run.ID, requested)
			var conflict TerminalOutcomeConflictError
			if !errors.As(err, &conflict) {
				t.Fatalf("expected terminal conflict for %s to %s, got %T %v", sourceOutcome, requested, err, err)
			}
			stored, found, err := runStore.Run(ctx, run.ID)
			if err != nil || !found {
				t.Fatalf("read immutable terminal Run: found=%v err=%v", found, err)
			}
			if stored.State != sourceOutcome ||
				stored.CompletedAt == nil ||
				first.CompletedAt == nil ||
				!stored.CompletedAt.Equal(*first.CompletedAt) {
				t.Fatalf("expected source outcome %s unchanged, got %+v", sourceOutcome, stored)
			}
		})
	}
}

func TestCompleteRunConcurrentTerminalOutcomesHaveOneWinner(t *testing.T) {
	ctx := context.Background()
	runStore := openTestStore(t, ctx, t.TempDir())
	defer closeStore(t, runStore)

	run, err := runStore.CreateRun(ctx, sampleCreateRunRequest())
	if err != nil {
		t.Fatalf("create Run: %v", err)
	}
	type completion struct {
		requested string
		result    CompleteRunResult
		err       error
	}
	start := make(chan struct{})
	completions := make(chan completion, 2)
	for _, requested := range []string{StateStopped, StateFailed} {
		go func(requested string) {
			<-start
			result, err := runStore.CompleteRun(ctx, run.ID, requested)
			completions <- completion{requested: requested, result: result, err: err}
		}(requested)
	}
	close(start)

	winner := ""
	conflicts := 0
	for range 2 {
		completed := <-completions
		if completed.err == nil {
			if !completed.result.Transitioned {
				t.Fatalf("expected successful competing completion for %s to transition", completed.requested)
			}
			if winner != "" {
				t.Fatalf("expected one winner, got %s and %s", winner, completed.requested)
			}
			winner = completed.requested
			continue
		}
		var conflict TerminalOutcomeConflictError
		if !errors.As(completed.err, &conflict) {
			t.Fatalf("expected competing completion conflict, got %T %v", completed.err, completed.err)
		}
		conflicts++
	}
	if winner == "" || conflicts != 1 {
		t.Fatalf("expected exactly one winner and one conflict, winner=%q conflicts=%d", winner, conflicts)
	}
	stored, found, err := runStore.Run(ctx, run.ID)
	if err != nil || !found {
		t.Fatalf("read competing completion winner: found=%v err=%v", found, err)
	}
	if stored.State != winner || stored.CompletedAt == nil {
		t.Fatalf("expected stable terminal winner %s, got %+v", winner, stored)
	}
	if got := countActiveRunLocks(t, ctx, runStore); got != 0 {
		t.Fatalf("expected one winning lock release, got %d locks", got)
	}
}

func TestCompleteRunDatabaseFailureNamesOperationAndRun(t *testing.T) {
	ctx := context.Background()
	runStore := openTestStore(t, ctx, t.TempDir())
	run, err := runStore.CreateRun(ctx, sampleCreateRunRequest())
	if err != nil {
		t.Fatalf("create Run: %v", err)
	}
	if err := runStore.Close(); err != nil {
		t.Fatalf("close store: %v", err)
	}

	_, err = runStore.CompleteRun(ctx, run.ID, StateStopped)
	if err == nil ||
		!strings.Contains(err.Error(), "begin completion") ||
		!strings.Contains(err.Error(), run.ID) {
		t.Fatalf("expected wrapped completion failure naming operation and Run %s, got %v", run.ID, err)
	}
}

func TestReconcileIntegrationPendingRecordsEvidence(t *testing.T) {
	ctx := context.Background()
	runStore := openTestStore(t, ctx, t.TempDir())
	defer closeStore(t, runStore)

	req := sampleCreateRunRequest()
	req.Kind = KindImplement
	req.SpecSlug = "0037-terminal-outcome-integrity"
	req.LocalBranch = "feature/terminal-outcomes"
	run, err := runStore.CreateRun(ctx, req)
	if err != nil {
		t.Fatalf("create Implement Run: %v", err)
	}
	completed, err := runStore.CompleteRun(ctx, run.ID, StateIntegrationPending)
	if err != nil {
		t.Fatalf("complete Run Integration Pending: %v", err)
	}
	reconciledAt := time.Date(2026, 7, 27, 11, 0, 0, 0, time.UTC)

	reconciled, err := runStore.ReconcileIntegration(ctx, IntegrationReconciliation{
		RunID:        run.ID,
		RunHead:      "run-head",
		TargetBranch: req.LocalBranch,
		TargetHead:   "target-head",
		Time:         reconciledAt,
	})
	if err != nil {
		t.Fatalf("reconcile Integration Pending Run: %v", err)
	}
	if reconciled.State != StateClean {
		t.Fatalf("expected reconciled Clean Run, got %+v", reconciled)
	}
	if reconciled.CompletedAt == nil || completed.Run.CompletedAt == nil ||
		!reconciled.CompletedAt.Equal(*completed.Run.CompletedAt) {
		t.Fatalf("expected reconciliation to preserve original completion time, before=%+v after=%+v", completed.Run, reconciled)
	}
	events, err := runStore.RunEventsAfter(ctx, run.ID, 0, 10)
	if err != nil {
		t.Fatalf("read reconciliation event: %v", err)
	}
	if len(events) != 1 {
		t.Fatalf("expected one reconciliation event, got %d", len(events))
	}
	event := events[0].Event
	if event.Kind != runevent.KindDaemonOutcome || !event.Time.Equal(reconciledAt) {
		t.Fatalf("unexpected reconciliation event: %+v", event)
	}
	var payload map[string]string
	if err := json.Unmarshal(event.Payload, &payload); err != nil {
		t.Fatalf("decode reconciliation event: %v", err)
	}
	wantPayload := map[string]string{
		"event":            "integration_reconciliation",
		"previous_outcome": StateIntegrationPending,
		"current_outcome":  StateClean,
		"run_head":         "run-head",
		"target_branch":    req.LocalBranch,
		"target_head":      "target-head",
	}
	for key, want := range wantPayload {
		if payload[key] != want {
			t.Fatalf("expected reconciliation payload %s=%q, got %q in %#v", key, want, payload[key], payload)
		}
	}
}

func TestReconcileIntegrationDatabaseFailureNamesOperationAndRun(t *testing.T) {
	ctx := context.Background()
	runStore := openTestStore(t, ctx, t.TempDir())
	req := sampleCreateRunRequest()
	req.Kind = KindImplement
	req.SpecSlug = "0037-terminal-outcome-integrity"
	run, err := runStore.CreateRun(ctx, req)
	if err != nil {
		t.Fatalf("create Implement Run: %v", err)
	}
	if _, err := runStore.CompleteRun(ctx, run.ID, StateIntegrationPending); err != nil {
		t.Fatalf("complete Run Integration Pending: %v", err)
	}
	if err := runStore.Close(); err != nil {
		t.Fatalf("close store: %v", err)
	}

	_, err = runStore.ReconcileIntegration(ctx, IntegrationReconciliation{
		RunID:        run.ID,
		RunHead:      "run-head",
		TargetBranch: req.LocalBranch,
		TargetHead:   "target-head",
		Time:         time.Date(2026, 7, 27, 15, 0, 0, 0, time.UTC),
	})
	if err == nil ||
		!strings.Contains(err.Error(), "begin Integration Pending reconciliation") ||
		!strings.Contains(err.Error(), run.ID) {
		t.Fatalf("expected wrapped reconciliation failure naming operation and Run %s, got %v", run.ID, err)
	}
}

func TestReconcileIntegrationRejectsIncompleteEvidence(t *testing.T) {
	ctx := context.Background()
	runStore := openTestStore(t, ctx, t.TempDir())
	defer closeStore(t, runStore)

	req := sampleCreateRunRequest()
	req.Kind = KindImplement
	req.SpecSlug = "0037-terminal-outcome-integrity"
	run, err := runStore.CreateRun(ctx, req)
	if err != nil {
		t.Fatalf("create Implement Run: %v", err)
	}
	completed, err := runStore.CompleteRun(ctx, run.ID, StateIntegrationPending)
	if err != nil {
		t.Fatalf("complete Run Integration Pending: %v", err)
	}
	valid := IntegrationReconciliation{
		RunID:        run.ID,
		RunHead:      "run-head",
		TargetBranch: req.LocalBranch,
		TargetHead:   "target-head",
		Time:         time.Date(2026, 7, 27, 12, 0, 0, 0, time.UTC),
	}
	tests := []struct {
		name   string
		mutate func(*IntegrationReconciliation)
	}{
		{name: "missing Run ID", mutate: func(value *IntegrationReconciliation) { value.RunID = "" }},
		{name: "missing Run head", mutate: func(value *IntegrationReconciliation) { value.RunHead = "" }},
		{name: "missing target branch", mutate: func(value *IntegrationReconciliation) { value.TargetBranch = "" }},
		{name: "missing target head", mutate: func(value *IntegrationReconciliation) { value.TargetHead = "" }},
		{name: "missing timestamp", mutate: func(value *IntegrationReconciliation) { value.Time = time.Time{} }},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input := valid
			tt.mutate(&input)
			if _, err := runStore.ReconcileIntegration(ctx, input); err == nil {
				t.Fatal("expected incomplete reconciliation evidence to fail")
			}
		})
	}
	stored, found, err := runStore.Run(ctx, run.ID)
	if err != nil || !found {
		t.Fatalf("read Integration Pending Run: found=%v err=%v", found, err)
	}
	if stored.State != StateIntegrationPending ||
		stored.CompletedAt == nil ||
		completed.Run.CompletedAt == nil ||
		!stored.CompletedAt.Equal(*completed.Run.CompletedAt) {
		t.Fatalf("expected invalid evidence to preserve Integration Pending Run, got %+v", stored)
	}
	if got := countRunEvents(t, ctx, runStore, run.ID); got != 0 {
		t.Fatalf("expected invalid evidence to append no event, got %d", got)
	}
}

func TestReconcileIntegrationRejectsStaleTargetBranch(t *testing.T) {
	ctx := context.Background()
	runStore := openTestStore(t, ctx, t.TempDir())
	defer closeStore(t, runStore)

	req := sampleCreateRunRequest()
	req.Kind = KindImplement
	req.SpecSlug = "0037-terminal-outcome-integrity"
	req.LocalBranch = "feature/terminal-outcomes"
	run, err := runStore.CreateRun(ctx, req)
	if err != nil {
		t.Fatalf("create Implement Run: %v", err)
	}
	completed, err := runStore.CompleteRun(ctx, run.ID, StateIntegrationPending)
	if err != nil {
		t.Fatalf("complete Run Integration Pending: %v", err)
	}

	_, err = runStore.ReconcileIntegration(ctx, IntegrationReconciliation{
		RunID:        run.ID,
		RunHead:      "run-head",
		TargetBranch: "feature/stale-target",
		TargetHead:   "target-head",
		Time:         time.Date(2026, 7, 27, 12, 30, 0, 0, time.UTC),
	})
	if err == nil || !strings.Contains(err.Error(), "does not match recorded target branch") {
		t.Fatalf("expected stale target branch rejection, got %v", err)
	}
	stored, found, readErr := runStore.Run(ctx, run.ID)
	if readErr != nil || !found {
		t.Fatalf("read Integration Pending Run: found=%v err=%v", found, readErr)
	}
	if stored.State != StateIntegrationPending ||
		stored.CompletedAt == nil ||
		completed.CompletedAt == nil ||
		!stored.CompletedAt.Equal(*completed.CompletedAt) {
		t.Fatalf("expected stale target evidence to preserve Integration Pending Run, got %+v", stored)
	}
	if got := countRunEvents(t, ctx, runStore, run.ID); got != 0 {
		t.Fatalf("expected stale target evidence to append no event, got %d", got)
	}
}

func TestReconcileIntegrationRejectsEveryOtherSourceOutcome(t *testing.T) {
	ctx := context.Background()
	runStore := openTestStore(t, ctx, t.TempDir())
	defer closeStore(t, runStore)

	sourceOutcomes := []string{
		StateFetched,
		StateStopped,
		StateClean,
		StateCleanUnverified,
		StateMaxRoundsReached,
		StateBudgetExceeded,
		StateTimedOut,
		StateFailed,
		StateUnresolved,
	}
	for index, sourceOutcome := range sourceOutcomes {
		t.Run(sourceOutcome, func(t *testing.T) {
			req := sampleCreateRunRequest()
			req.Kind = KindImplement
			req.SpecSlug = fmt.Sprintf("terminal-source-%d", index)
			req.LocalBranch = fmt.Sprintf("feature/terminal-source-%d", index)
			run, err := runStore.CreateRun(ctx, req)
			if err != nil {
				t.Fatalf("create Implement Run: %v", err)
			}
			completed, err := runStore.CompleteRun(ctx, run.ID, sourceOutcome)
			if err != nil {
				t.Fatalf("complete Run %s: %v", sourceOutcome, err)
			}

			_, err = runStore.ReconcileIntegration(ctx, IntegrationReconciliation{
				RunID:        run.ID,
				RunHead:      "run-head",
				TargetBranch: req.LocalBranch,
				TargetHead:   "target-head",
				Time:         time.Date(2026, 7, 27, 13, 0, 0, 0, time.UTC),
			})
			var conflict TerminalOutcomeConflictError
			if !errors.As(err, &conflict) {
				t.Fatalf("expected terminal conflict for source %s, got %T %v", sourceOutcome, err, err)
			}
			stored, found, err := runStore.Run(ctx, run.ID)
			if err != nil || !found {
				t.Fatalf("read rejected source Run: found=%v err=%v", found, err)
			}
			if stored.State != sourceOutcome ||
				stored.CompletedAt == nil ||
				completed.Run.CompletedAt == nil ||
				!stored.CompletedAt.Equal(*completed.Run.CompletedAt) {
				t.Fatalf("expected source outcome %s unchanged, got %+v", sourceOutcome, stored)
			}
			if got := countRunEvents(t, ctx, runStore, run.ID); got != 0 {
				t.Fatalf("expected source outcome rejection to append no event, got %d", got)
			}
		})
	}
}

func TestReconcileIntegrationRollsBackWhenJournalFails(t *testing.T) {
	ctx := context.Background()
	runStore := openTestStore(t, ctx, t.TempDir())
	defer closeStore(t, runStore)

	req := sampleCreateRunRequest()
	req.Kind = KindImplement
	req.SpecSlug = "0037-terminal-outcome-integrity"
	run, err := runStore.CreateRun(ctx, req)
	if err != nil {
		t.Fatalf("create Implement Run: %v", err)
	}
	completed, err := runStore.CompleteRun(ctx, run.ID, StateIntegrationPending)
	if err != nil {
		t.Fatalf("complete Run Integration Pending: %v", err)
	}
	if _, err := runStore.db.ExecContext(ctx, `
CREATE TRIGGER reject_reconciliation_event
BEFORE INSERT ON run_events
BEGIN
	SELECT RAISE(FAIL, 'reconciliation journal unavailable');
END`); err != nil {
		t.Fatalf("create reconciliation journal failure trigger: %v", err)
	}

	_, err = runStore.ReconcileIntegration(ctx, IntegrationReconciliation{
		RunID:        run.ID,
		RunHead:      "run-head",
		TargetBranch: req.LocalBranch,
		TargetHead:   "target-head",
		Time:         time.Date(2026, 7, 27, 14, 0, 0, 0, time.UTC),
	})
	if err == nil || !strings.Contains(err.Error(), run.ID) {
		t.Fatalf("expected reconciliation journal failure naming Run %s, got %v", run.ID, err)
	}
	stored, found, readErr := runStore.Run(ctx, run.ID)
	if readErr != nil || !found {
		t.Fatalf("read rolled-back Run: found=%v err=%v", found, readErr)
	}
	if stored.State != StateIntegrationPending ||
		stored.CompletedAt == nil ||
		completed.Run.CompletedAt == nil ||
		!stored.CompletedAt.Equal(*completed.Run.CompletedAt) {
		t.Fatalf("expected failed journal append to roll back Run, got %+v", stored)
	}
	if got := countRunEvents(t, ctx, runStore, run.ID); got != 0 {
		t.Fatalf("expected failed reconciliation to append no event, got %d", got)
	}
}

func TestRequestStopRecordsStopRequestForActiveRun(t *testing.T) {
	ctx := context.Background()
	runStore := openTestStore(t, ctx, t.TempDir())
	defer closeStore(t, runStore)
	now := time.Date(2026, 7, 5, 10, 30, 0, 0, time.UTC)
	runStore.now = func() time.Time { return now }

	run, err := runStore.CreateRun(ctx, sampleCreateRunRequest())
	if err != nil {
		t.Fatalf("create Run: %v", err)
	}
	requested, err := runStore.StopRequested(ctx, run.ID)
	if err != nil {
		t.Fatalf("read initial Stop Request flag: %v", err)
	}
	if requested {
		t.Fatal("expected new Active Run without a Stop Request")
	}
	if err := runStore.UpdateRunState(ctx, run.ID, StateResolvingWithAgent); err != nil {
		t.Fatalf("update active Run state: %v", err)
	}

	if err := runStore.RequestStop(ctx, run.ID); err != nil {
		t.Fatalf("request Stop: %v", err)
	}

	requested, err = runStore.StopRequested(ctx, run.ID)
	if err != nil {
		t.Fatalf("read Stop Request flag: %v", err)
	}
	if !requested {
		t.Fatal("expected Stop Request flag recorded")
	}
	active, found, err := runStore.ActiveRun(ctx, run.HeadRepository, run.HeadBranch)
	if err != nil {
		t.Fatalf("active run lookup after Stop Request: %v", err)
	}
	if !found || active.ID != run.ID {
		t.Fatalf("expected Stop Request to keep Active Run lock, found=%v active=%#v", found, active)
	}
	var recorded string
	if err := runStore.db.QueryRowContext(ctx, `SELECT stop_requested_at FROM runs WHERE id = ?`, run.ID).Scan(&recorded); err != nil {
		t.Fatalf("read stop_requested_at: %v", err)
	}
	if recorded != formatTime(now) {
		t.Fatalf("expected stop_requested_at %q, got %q", formatTime(now), recorded)
	}
}

func TestRequestStopRejectsTerminalRunWithNamedError(t *testing.T) {
	ctx := context.Background()
	runStore := openTestStore(t, ctx, t.TempDir())
	defer closeStore(t, runStore)

	run, err := runStore.CreateRun(ctx, sampleCreateRunRequest())
	if err != nil {
		t.Fatalf("create Run: %v", err)
	}
	if _, err := runStore.CompleteRun(ctx, run.ID, StateClean); err != nil {
		t.Fatalf("complete Run: %v", err)
	}

	err = runStore.RequestStop(ctx, run.ID)

	if !errors.Is(err, ErrTerminalRunStopRequest) {
		t.Fatalf("expected ErrTerminalRunStopRequest, got %T %v", err, err)
	}
	requested, flagErr := runStore.StopRequested(ctx, run.ID)
	if flagErr != nil {
		t.Fatalf("read Stop Request flag: %v", flagErr)
	}
	if requested {
		t.Fatal("expected terminal Run to keep Stop Request flag unset")
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

func TestIsTerminalState(t *testing.T) {
	tests := []struct {
		state string
		want  bool
	}{
		{state: StateActive},
		{state: StateResolvingWithAgent},
		{state: StateVerifying},
		{state: StatePushing},
		{state: StateFetched, want: true},
		{state: StateStopped, want: true},
		{state: StateClean, want: true},
		{state: StateCleanUnverified, want: true},
		{state: StateMaxRoundsReached, want: true},
		{state: StateBudgetExceeded, want: true},
		{state: StateTimedOut, want: true},
		{state: StateFailed, want: true},
		{state: StateIntegrationPending, want: true},
		{state: StateUnresolved, want: true},
		{state: "unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.state, func(t *testing.T) {
			if got := IsTerminalState(tt.state); got != tt.want {
				t.Fatalf("expected IsTerminalState(%q) = %v, got %v", tt.state, tt.want, got)
			}
		})
	}
}

type listedRunSeed struct {
	gitRoot   string
	branch    string
	prNumber  string
	createdAt time.Time
	state     string
}

func createListedRun(t *testing.T, ctx context.Context, store *Store, seed listedRunSeed) Run {
	t.Helper()
	store.now = func() time.Time { return seed.createdAt }
	req := sampleCreateRunRequest()
	req.GitRoot = seed.gitRoot
	req.HeadBranch = seed.branch
	req.LocalBranch = seed.branch
	req.PRNumber = seed.prNumber
	req.ArtifactDir = filepath.Join(seed.gitRoot, ".roundfix")
	run, err := store.CreateRun(ctx, req)
	if err != nil {
		t.Fatalf("create listed Run: %v", err)
	}

	switch {
	case seed.state == "" || seed.state == StateActive:
		return run
	case IsTerminalState(seed.state):
		store.now = func() time.Time { return seed.createdAt.Add(time.Minute) }
		completed, err := store.CompleteRun(ctx, run.ID, seed.state)
		if err != nil {
			t.Fatalf("complete listed Run as %s: %v", seed.state, err)
		}
		return completed.Run
	default:
		if err := store.UpdateRunState(ctx, run.ID, seed.state); err != nil {
			t.Fatalf("update listed Run state to %s: %v", seed.state, err)
		}
		updated, ok, err := store.Run(ctx, run.ID)
		if err != nil {
			t.Fatalf("read updated listed Run: %v", err)
		}
		if !ok {
			t.Fatalf("expected updated listed Run %s to exist", run.ID)
		}
		return updated
	}
}

func assertRunIDs(t *testing.T, runs []Run, want []string) {
	t.Helper()
	if len(runs) != len(want) {
		t.Fatalf("expected Run IDs %v, got %v", want, runIDs(runs))
	}
	for index, run := range runs {
		if run.ID != want[index] {
			t.Fatalf("expected Run IDs %v, got %v", want, runIDs(runs))
		}
	}
}

func runIDs(runs []Run) []string {
	ids := make([]string, 0, len(runs))
	for _, run := range runs {
		ids = append(ids, run.ID)
	}
	return ids
}

func assertRunSelection(t *testing.T, run Run, wantModel string, wantReasoning string) {
	t.Helper()
	if run.Model != wantModel || run.ReasoningEffort != wantReasoning {
		t.Fatalf("expected Run %s selection %q/%q, got %q/%q", run.ID, wantModel, wantReasoning, run.Model, run.ReasoningEffort)
	}
}

func assertRunOwnerPID(t *testing.T, run Run, wantPID int) {
	t.Helper()
	if run.OwnerPID == nil {
		t.Fatalf("expected Run %s owner PID %d, got nil", run.ID, wantPID)
	}
	if *run.OwnerPID != wantPID {
		t.Fatalf("expected Run %s owner PID %d, got %d", run.ID, wantPID, *run.OwnerPID)
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
		OwnerPID:       os.Getpid(),
	}
}

// buildV3Fixture creates a populated schema v3 Run Database via raw SQL:
// runs in several states plus one Active Run lock in the v3
// (head_repository, head_branch) shape.
func buildV3Fixture(t *testing.T, homeDir string) {
	t.Helper()
	path := DatabasePath(homeDir)
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("create fixture directory: %v", err)
	}
	db, err := sql.Open("sqlite", writerDSN(path))
	if err != nil {
		t.Fatalf("open fixture database: %v", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			t.Fatalf("close fixture database: %v", err)
		}
	}()

	statements := []string{
		`CREATE TABLE runs (
			id TEXT PRIMARY KEY,
			kind TEXT NOT NULL,
			state TEXT NOT NULL,
			head_repository TEXT NOT NULL,
			head_branch TEXT NOT NULL,
			base_repository TEXT NOT NULL DEFAULT '',
			pr_number TEXT NOT NULL,
			git_root TEXT NOT NULL,
			local_branch TEXT NOT NULL,
			head_sha TEXT NOT NULL,
			artifact_dir TEXT NOT NULL,
			created_at TEXT NOT NULL,
			updated_at TEXT NOT NULL,
			completed_at TEXT NOT NULL DEFAULT ''
		)`,
		`CREATE TABLE active_run_locks (
			head_repository TEXT NOT NULL,
			head_branch TEXT NOT NULL,
			run_id TEXT NOT NULL UNIQUE,
			created_at TEXT NOT NULL,
			PRIMARY KEY (head_repository, head_branch),
			FOREIGN KEY (run_id) REFERENCES runs(id) ON DELETE CASCADE
		)`,
		`CREATE INDEX idx_runs_head ON runs (head_repository, head_branch)`,
		`CREATE TABLE interactive_defaults (
			key TEXT PRIMARY KEY,
			value TEXT NOT NULL,
			updated_at TEXT NOT NULL
		)`,
		`CREATE TABLE run_events (
			run_id TEXT NOT NULL,
			cursor INTEGER NOT NULL,
			batch INTEGER NOT NULL DEFAULT 0,
			source TEXT NOT NULL,
			kind TEXT NOT NULL,
			review_issue TEXT NOT NULL DEFAULT '',
			tool_id TEXT NOT NULL DEFAULT '',
			tool_state TEXT NOT NULL DEFAULT '',
			summary TEXT NOT NULL DEFAULT '',
			created_at TEXT NOT NULL,
			payload TEXT NOT NULL DEFAULT '',
			PRIMARY KEY (run_id, cursor),
			FOREIGN KEY (run_id) REFERENCES runs(id) ON DELETE CASCADE
		)`,
		`INSERT INTO runs (id, kind, state, head_repository, head_branch, base_repository,
			pr_number, git_root, local_branch, head_sha, artifact_dir, created_at, updated_at, completed_at)
		 VALUES ('run_v3_active', 'resolve', 'Active', 'owner/project', 'feature/review', 'owner/project',
			'123', 'tmp/repo', 'feature/review', 'abc123', 'tmp/repo/.roundfix',
			'2026-07-01T10:00:00Z', '2026-07-01T10:00:00Z', '')`,
		`INSERT INTO runs (id, kind, state, head_repository, head_branch, base_repository,
			pr_number, git_root, local_branch, head_sha, artifact_dir, created_at, updated_at, completed_at)
		 VALUES ('run_v3_clean', 'resolve', 'Clean', 'owner/project', 'feature/done', 'owner/project',
			'99', 'tmp/repo', 'feature/done', 'def456', 'tmp/repo/.roundfix',
			'2026-07-01T08:00:00Z', '2026-07-01T09:00:00Z', '2026-07-01T09:00:00Z')`,
		`INSERT INTO runs (id, kind, state, head_repository, head_branch, base_repository,
			pr_number, git_root, local_branch, head_sha, artifact_dir, created_at, updated_at, completed_at)
		 VALUES ('run_v3_fetched', 'fetch', 'Fetched', 'owner/other', 'feature/fetch', 'owner/other',
			'7', 'tmp/other', 'feature/fetch', 'fed789', 'tmp/other/.roundfix',
			'2026-07-01T07:00:00Z', '2026-07-01T07:30:00Z', '2026-07-01T07:30:00Z')`,
		`INSERT INTO active_run_locks (head_repository, head_branch, run_id, created_at)
		 VALUES ('owner/project', 'feature/review', 'run_v3_active', '2026-07-01T10:00:00Z')`,
		`PRAGMA user_version = 3`,
	}
	for _, statement := range statements {
		if _, err := db.Exec(statement); err != nil {
			t.Fatalf("build v3 fixture: %v", err)
		}
	}
}

func TestOpenMigratesV3RunDatabasePreservingRunsAndRekeyingLocks(t *testing.T) {
	ctx := context.Background()
	homeDir := t.TempDir()
	buildV3Fixture(t, homeDir)

	store := openTestStore(t, ctx, homeDir)
	defer closeStore(t, store)

	version, err := store.MigrationVersion(ctx)
	if err != nil {
		t.Fatalf("read migration version: %v", err)
	}
	if version != 9 {
		t.Fatalf("expected user_version 9 after migration, got %d", version)
	}

	count, err := store.RunCount(ctx)
	if err != nil {
		t.Fatalf("count runs: %v", err)
	}
	if count != 3 {
		t.Fatalf("expected all 3 v3 run rows to survive, got %d", count)
	}

	clean, ok, err := store.Run(ctx, "run_v3_clean")
	if err != nil || !ok {
		t.Fatalf("expected run_v3_clean to survive, ok=%v err=%v", ok, err)
	}
	if clean.Kind != KindResolve || clean.State != StateClean || clean.PRNumber != "99" ||
		clean.HeadRepository != "owner/project" || clean.HeadBranch != "feature/done" ||
		clean.HeadSHA != "def456" || clean.ArtifactDir != "tmp/repo/.roundfix" ||
		clean.WorkDir != "" || clean.SpecSlug != "" || clean.CompletedAt == nil {
		t.Fatalf("expected run_v3_clean fields preserved, got %#v", clean)
	}

	active, found, err := store.ActiveRun(ctx, "owner/project", "feature/review")
	if err != nil {
		t.Fatalf("active run lookup after migration: %v", err)
	}
	if !found || active.ID != "run_v3_active" {
		t.Fatalf("expected re-keyed lock to keep run_v3_active active, found=%v active=%#v", found, active)
	}

	var lockCount int
	if err := store.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM active_run_locks`).Scan(&lockCount); err != nil {
		t.Fatalf("count migrated locks: %v", err)
	}
	if lockCount != 1 {
		t.Fatalf("expected exactly one migrated lock row, got %d", lockCount)
	}
	var targetKind, targetKey, runID string
	if err := store.db.QueryRowContext(ctx,
		`SELECT target_kind, target_key, run_id FROM active_run_locks`).Scan(&targetKind, &targetKey, &runID); err != nil {
		t.Fatalf("read migrated lock row: %v", err)
	}
	if targetKind != "pr" || targetKey != "owner/project#feature/review" || runID != "run_v3_active" {
		t.Fatalf("expected lock re-keyed to (pr, owner/project#feature/review, run_v3_active), got (%s, %s, %s)",
			targetKind, targetKey, runID)
	}

	implementRun, err := store.CreateRun(ctx, sampleImplementCreateRunRequest())
	if err != nil {
		t.Fatalf("expected implement Run creation on migrated database, got %v", err)
	}
	if implementRun.Kind != KindImplement || implementRun.SpecSlug != "0001-implement-command" {
		t.Fatalf("expected persisted implement Run with Spec slug, got %#v", implementRun)
	}
}

// buildV4Fixture creates a populated schema v4 Run Database via raw SQL:
// runs in several states plus one Active Run lock in the v4 work-target
// shape. It intentionally omits stop_requested_at so the v5 migration must
// add it without disturbing existing rows or locks.
func buildV4Fixture(t *testing.T, homeDir string) {
	t.Helper()
	path := DatabasePath(homeDir)
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("create fixture directory: %v", err)
	}
	db, err := sql.Open("sqlite", writerDSN(path))
	if err != nil {
		t.Fatalf("open fixture database: %v", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			t.Fatalf("close fixture database: %v", err)
		}
	}()

	statements := []string{
		`CREATE TABLE runs (
			id TEXT PRIMARY KEY,
			kind TEXT NOT NULL,
			state TEXT NOT NULL,
			head_repository TEXT NOT NULL DEFAULT '',
			head_branch TEXT NOT NULL DEFAULT '',
			base_repository TEXT NOT NULL DEFAULT '',
			pr_number TEXT NOT NULL DEFAULT '',
			git_root TEXT NOT NULL,
			local_branch TEXT NOT NULL,
			head_sha TEXT NOT NULL DEFAULT '',
			artifact_dir TEXT NOT NULL DEFAULT '',
			spec_slug TEXT NOT NULL DEFAULT '',
			created_at TEXT NOT NULL,
			updated_at TEXT NOT NULL,
			completed_at TEXT NOT NULL DEFAULT ''
		)`,
		`CREATE TABLE active_run_locks ` + activeRunLocksColumns,
		`CREATE INDEX idx_runs_head ON runs (head_repository, head_branch)`,
		`CREATE TABLE interactive_defaults (
			key TEXT PRIMARY KEY,
			value TEXT NOT NULL,
			updated_at TEXT NOT NULL
		)`,
		`CREATE TABLE run_events (
			run_id TEXT NOT NULL,
			cursor INTEGER NOT NULL,
			batch INTEGER NOT NULL DEFAULT 0,
			source TEXT NOT NULL,
			kind TEXT NOT NULL,
			review_issue TEXT NOT NULL DEFAULT '',
			tool_id TEXT NOT NULL DEFAULT '',
			tool_state TEXT NOT NULL DEFAULT '',
			summary TEXT NOT NULL DEFAULT '',
			created_at TEXT NOT NULL,
			payload TEXT NOT NULL DEFAULT '',
			PRIMARY KEY (run_id, cursor),
			FOREIGN KEY (run_id) REFERENCES runs(id) ON DELETE CASCADE
		)`,
		`INSERT INTO runs (id, kind, state, head_repository, head_branch, base_repository,
			pr_number, git_root, local_branch, head_sha, artifact_dir, spec_slug, created_at, updated_at, completed_at)
		 VALUES ('run_v4_active', 'resolve', 'ResolvingWithAgent', 'owner/project', 'feature/review', 'owner/project',
			'123', 'tmp/repo', 'feature/review', 'abc123', 'tmp/repo/.roundfix', '',
			'2026-07-01T10:00:00Z', '2026-07-01T10:05:00Z', '')`,
		`INSERT INTO runs (id, kind, state, head_repository, head_branch, base_repository,
			pr_number, git_root, local_branch, head_sha, artifact_dir, spec_slug, created_at, updated_at, completed_at)
		 VALUES ('run_v4_clean', 'watch', 'Clean', 'owner/project', 'feature/done', 'owner/project',
			'99', 'tmp/repo', 'feature/done', 'def456', 'tmp/repo/.roundfix', '',
			'2026-07-01T08:00:00Z', '2026-07-01T09:00:00Z', '2026-07-01T09:00:00Z')`,
		`INSERT INTO runs (id, kind, state, head_repository, head_branch, base_repository,
			pr_number, git_root, local_branch, head_sha, artifact_dir, spec_slug, created_at, updated_at, completed_at)
		 VALUES ('run_v4_implement', 'implement', 'Stopped', '', '', '',
			'', 'tmp/spec-repo', 'ma/spec-work', '', '', '0001-widget-flow',
			'2026-07-01T07:00:00Z', '2026-07-01T07:30:00Z', '2026-07-01T07:30:00Z')`,
		`INSERT INTO active_run_locks (target_kind, target_key, run_id, created_at)
		 VALUES ('pr', 'owner/project#feature/review', 'run_v4_active', '2026-07-01T10:00:00Z')`,
		`PRAGMA user_version = 4`,
	}
	for _, statement := range statements {
		if _, err := db.Exec(statement); err != nil {
			t.Fatalf("build v4 fixture: %v", err)
		}
	}
}

func TestOpenMigratesV4RunDatabasePreservingRunsLocksAndAddingStopRequests(t *testing.T) {
	ctx := context.Background()
	homeDir := t.TempDir()
	buildV4Fixture(t, homeDir)

	runStore := openTestStore(t, ctx, homeDir)
	defer closeStore(t, runStore)

	version, err := runStore.MigrationVersion(ctx)
	if err != nil {
		t.Fatalf("read migration version: %v", err)
	}
	if version != 9 {
		t.Fatalf("expected user_version 9 after migration, got %d", version)
	}
	count, err := runStore.RunCount(ctx)
	if err != nil {
		t.Fatalf("count runs: %v", err)
	}
	if count != 3 {
		t.Fatalf("expected all 3 v4 run rows to survive, got %d", count)
	}
	clean, ok, err := runStore.Run(ctx, "run_v4_clean")
	if err != nil || !ok {
		t.Fatalf("expected run_v4_clean to survive, ok=%v err=%v", ok, err)
	}
	if clean.Kind != KindWatch || clean.State != StateClean || clean.PRNumber != "99" ||
		clean.HeadRepository != "owner/project" || clean.HeadBranch != "feature/done" ||
		clean.HeadSHA != "def456" || clean.WorkDir != "" || clean.SpecSlug != "" || clean.CompletedAt == nil {
		t.Fatalf("expected run_v4_clean fields preserved, got %#v", clean)
	}
	implement, ok, err := runStore.Run(ctx, "run_v4_implement")
	if err != nil || !ok {
		t.Fatalf("expected run_v4_implement to survive, ok=%v err=%v", ok, err)
	}
	if implement.Kind != KindImplement || implement.SpecSlug != "0001-widget-flow" || implement.GitRoot != "tmp/spec-repo" ||
		implement.WorkDir != "" || implement.CompletedAt == nil {
		t.Fatalf("expected implement fields preserved, got %#v", implement)
	}
	active, found, err := runStore.ActiveRun(ctx, "owner/project", "feature/review")
	if err != nil {
		t.Fatalf("active run lookup after migration: %v", err)
	}
	if !found || active.ID != "run_v4_active" || active.State != StateResolvingWithAgent {
		t.Fatalf("expected v4 active lock to survive, found=%v active=%#v", found, active)
	}
	requested, err := runStore.StopRequested(ctx, "run_v4_active")
	if err != nil {
		t.Fatalf("read migrated Stop Request flag: %v", err)
	}
	if requested {
		t.Fatal("expected migrated v4 Run to have no Stop Request")
	}
	var stopRequestedAt any
	if err := runStore.db.QueryRowContext(ctx, `SELECT stop_requested_at FROM runs WHERE id = 'run_v4_active'`).Scan(&stopRequestedAt); err != nil {
		t.Fatalf("read migrated stop_requested_at column: %v", err)
	}
	if stopRequestedAt != nil {
		t.Fatalf("expected migrated stop_requested_at NULL, got %#v", stopRequestedAt)
	}
}

// buildV5Fixture creates a populated schema v5 Run Database via raw SQL:
// persisted rows have agents and Stop Request state but no work_dir column.
func buildV5Fixture(t *testing.T, homeDir string) {
	t.Helper()
	path := DatabasePath(homeDir)
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("create fixture directory: %v", err)
	}
	db, err := sql.Open("sqlite", writerDSN(path))
	if err != nil {
		t.Fatalf("open fixture database: %v", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			t.Fatalf("close fixture database: %v", err)
		}
	}()

	statements := []string{
		`CREATE TABLE runs (
			id TEXT PRIMARY KEY,
			kind TEXT NOT NULL,
			state TEXT NOT NULL,
			head_repository TEXT NOT NULL DEFAULT '',
			head_branch TEXT NOT NULL DEFAULT '',
			base_repository TEXT NOT NULL DEFAULT '',
			pr_number TEXT NOT NULL DEFAULT '',
			git_root TEXT NOT NULL,
			local_branch TEXT NOT NULL,
			head_sha TEXT NOT NULL DEFAULT '',
			artifact_dir TEXT NOT NULL DEFAULT '',
			spec_slug TEXT NOT NULL DEFAULT '',
			agent TEXT NOT NULL DEFAULT '',
			stop_requested_at TEXT,
			created_at TEXT NOT NULL,
			updated_at TEXT NOT NULL,
			completed_at TEXT NOT NULL DEFAULT ''
		)`,
		`CREATE TABLE active_run_locks ` + activeRunLocksColumns,
		`CREATE INDEX idx_runs_head ON runs (head_repository, head_branch)`,
		`CREATE TABLE interactive_defaults (
			key TEXT PRIMARY KEY,
			value TEXT NOT NULL,
			updated_at TEXT NOT NULL
		)`,
		`CREATE TABLE run_events (
			run_id TEXT NOT NULL,
			cursor INTEGER NOT NULL,
			batch INTEGER NOT NULL DEFAULT 0,
			source TEXT NOT NULL,
			kind TEXT NOT NULL,
			review_issue TEXT NOT NULL DEFAULT '',
			tool_id TEXT NOT NULL DEFAULT '',
			tool_state TEXT NOT NULL DEFAULT '',
			summary TEXT NOT NULL DEFAULT '',
			created_at TEXT NOT NULL,
			payload TEXT NOT NULL DEFAULT '',
			PRIMARY KEY (run_id, cursor),
			FOREIGN KEY (run_id) REFERENCES runs(id) ON DELETE CASCADE
		)`,
		`INSERT INTO runs (id, kind, state, head_repository, head_branch, base_repository,
			pr_number, git_root, local_branch, head_sha, artifact_dir, spec_slug, agent,
			stop_requested_at, created_at, updated_at, completed_at)
		 VALUES ('run_v5_active', 'resolve', 'Verifying', 'owner/project', 'feature/review', 'owner/project',
			'123', 'tmp/repo', 'feature/review', 'abc123', 'tmp/repo/.roundfix', '', 'codex',
			'2026-07-01T10:06:00Z', '2026-07-01T10:00:00Z', '2026-07-01T10:06:00Z', '')`,
		`INSERT INTO runs (id, kind, state, head_repository, head_branch, base_repository,
			pr_number, git_root, local_branch, head_sha, artifact_dir, spec_slug, agent,
			stop_requested_at, created_at, updated_at, completed_at)
		 VALUES ('run_v5_clean', 'watch', 'Clean', 'owner/project', 'feature/done', 'owner/project',
			'99', 'tmp/repo', 'feature/done', 'def456', 'tmp/repo/.roundfix', '', 'claude',
			NULL, '2026-07-01T08:00:00Z', '2026-07-01T09:00:00Z', '2026-07-01T09:00:00Z')`,
		`INSERT INTO runs (id, kind, state, head_repository, head_branch, base_repository,
			pr_number, git_root, local_branch, head_sha, artifact_dir, spec_slug, agent,
			stop_requested_at, created_at, updated_at, completed_at)
		 VALUES ('run_v5_implement', 'implement', 'Active', '', '', '',
			'', 'tmp/spec-repo', 'ma/spec-work', '', '', '0001-widget-flow', 'opencode',
			NULL, '2026-07-01T07:00:00Z', '2026-07-01T07:30:00Z', '')`,
		`INSERT INTO active_run_locks (target_kind, target_key, run_id, created_at)
		 VALUES ('pr', 'owner/project#feature/review', 'run_v5_active', '2026-07-01T10:00:00Z')`,
		`INSERT INTO active_run_locks (target_kind, target_key, run_id, created_at)
		 VALUES ('spec', 'tmp/spec-repo#0001-widget-flow', 'run_v5_implement', '2026-07-01T07:00:00Z')`,
		`PRAGMA user_version = 5`,
	}
	for _, statement := range statements {
		if _, err := db.Exec(statement); err != nil {
			t.Fatalf("build v5 fixture: %v", err)
		}
	}
}

func TestOpenMigratesV5RunDatabasePreservingRunsLocksAndAddingWorkDir(t *testing.T) {
	ctx := context.Background()
	homeDir := t.TempDir()
	buildV5Fixture(t, homeDir)

	runStore := openTestStore(t, ctx, homeDir)
	defer closeStore(t, runStore)

	version, err := runStore.MigrationVersion(ctx)
	if err != nil {
		t.Fatalf("read migration version: %v", err)
	}
	if version != 9 {
		t.Fatalf("expected user_version 9 after migration, got %d", version)
	}
	count, err := runStore.RunCount(ctx)
	if err != nil {
		t.Fatalf("count runs: %v", err)
	}
	if count != 3 {
		t.Fatalf("expected all 3 v5 run rows to survive, got %d", count)
	}

	active, found, err := runStore.ActiveRun(ctx, "owner/project", "feature/review")
	if err != nil {
		t.Fatalf("active review Run lookup after migration: %v", err)
	}
	if !found || active.ID != "run_v5_active" || active.State != StateVerifying ||
		active.Agent != "codex" || active.WorkDir != "" {
		t.Fatalf("expected v5 review lock and row to survive with empty WorkDir, found=%v active=%#v", found, active)
	}
	requested, err := runStore.StopRequested(ctx, "run_v5_active")
	if err != nil {
		t.Fatalf("read migrated Stop Request flag: %v", err)
	}
	if !requested {
		t.Fatal("expected populated v5 Stop Request to survive")
	}

	implement, found, err := runStore.ActiveSpecRun(ctx, "tmp/spec-repo", "0001-widget-flow")
	if err != nil {
		t.Fatalf("active spec Run lookup after migration: %v", err)
	}
	if !found || implement.ID != "run_v5_implement" || implement.Agent != "opencode" || implement.WorkDir != "" {
		t.Fatalf("expected v5 spec lock and row to survive with empty WorkDir, found=%v implement=%#v", found, implement)
	}

	clean, ok, err := runStore.Run(ctx, "run_v5_clean")
	if err != nil || !ok {
		t.Fatalf("expected run_v5_clean to survive, ok=%v err=%v", ok, err)
	}
	if clean.Kind != KindWatch || clean.State != StateClean || clean.Agent != "claude" ||
		clean.WorkDir != "" || clean.CompletedAt == nil {
		t.Fatalf("expected run_v5_clean fields preserved with empty WorkDir, got %#v", clean)
	}

	var rawWorkDir any
	if err := runStore.db.QueryRowContext(ctx, `SELECT work_dir FROM runs WHERE id = 'run_v5_active'`).Scan(&rawWorkDir); err != nil {
		t.Fatalf("read migrated work_dir column: %v", err)
	}
	if rawWorkDir != nil {
		t.Fatalf("expected migrated work_dir NULL for legacy row, got %#v", rawWorkDir)
	}
	var lockCount int
	if err := runStore.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM active_run_locks`).Scan(&lockCount); err != nil {
		t.Fatalf("count migrated locks: %v", err)
	}
	if lockCount != 2 {
		t.Fatalf("expected both v5 locks to survive, got %d", lockCount)
	}
}

// buildV6Fixture creates a populated schema v6 Run Database via raw SQL:
// persisted rows have work_dir but no model or reasoning_effort columns.
func buildV6Fixture(t *testing.T, homeDir string) {
	t.Helper()
	path := DatabasePath(homeDir)
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("create fixture directory: %v", err)
	}
	db, err := sql.Open("sqlite", writerDSN(path))
	if err != nil {
		t.Fatalf("open fixture database: %v", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			t.Fatalf("close fixture database: %v", err)
		}
	}()

	statements := []string{
		`CREATE TABLE runs (
			id TEXT PRIMARY KEY,
			kind TEXT NOT NULL,
			state TEXT NOT NULL,
			head_repository TEXT NOT NULL DEFAULT '',
			head_branch TEXT NOT NULL DEFAULT '',
			base_repository TEXT NOT NULL DEFAULT '',
			pr_number TEXT NOT NULL DEFAULT '',
			git_root TEXT NOT NULL,
			local_branch TEXT NOT NULL,
			head_sha TEXT NOT NULL DEFAULT '',
			artifact_dir TEXT NOT NULL DEFAULT '',
			work_dir TEXT,
			spec_slug TEXT NOT NULL DEFAULT '',
			agent TEXT NOT NULL DEFAULT '',
			stop_requested_at TEXT,
			created_at TEXT NOT NULL,
			updated_at TEXT NOT NULL,
			completed_at TEXT NOT NULL DEFAULT ''
		)`,
		`CREATE TABLE active_run_locks ` + activeRunLocksColumns,
		`CREATE INDEX idx_runs_head ON runs (head_repository, head_branch)`,
		`CREATE TABLE interactive_defaults (
			key TEXT PRIMARY KEY,
			value TEXT NOT NULL,
			updated_at TEXT NOT NULL
		)`,
		`CREATE TABLE run_events (
			run_id TEXT NOT NULL,
			cursor INTEGER NOT NULL,
			batch INTEGER NOT NULL DEFAULT 0,
			source TEXT NOT NULL,
			kind TEXT NOT NULL,
			review_issue TEXT NOT NULL DEFAULT '',
			tool_id TEXT NOT NULL DEFAULT '',
			tool_state TEXT NOT NULL DEFAULT '',
			summary TEXT NOT NULL DEFAULT '',
			created_at TEXT NOT NULL,
			payload TEXT NOT NULL DEFAULT '',
			PRIMARY KEY (run_id, cursor),
			FOREIGN KEY (run_id) REFERENCES runs(id) ON DELETE CASCADE
		)`,
		`INSERT INTO runs (id, kind, state, head_repository, head_branch, base_repository,
			pr_number, git_root, local_branch, head_sha, artifact_dir, work_dir, spec_slug, agent,
			stop_requested_at, created_at, updated_at, completed_at)
		 VALUES ('run_v6_active', 'resolve', 'ResolvingWithAgent', 'owner/project', 'feature/review', 'owner/project',
			'123', 'tmp/repo', 'feature/review', 'abc123', 'tmp/repo/.roundfix', 'tmp/repo/.roundfix/worktrees/run_v6_active', '', 'codex',
			NULL, '2026-07-01T10:00:00Z', '2026-07-01T10:06:00Z', '')`,
		`INSERT INTO runs (id, kind, state, head_repository, head_branch, base_repository,
			pr_number, git_root, local_branch, head_sha, artifact_dir, work_dir, spec_slug, agent,
			stop_requested_at, created_at, updated_at, completed_at)
		 VALUES ('run_v6_implement', 'implement', 'Stopped', '', '', '',
			'', 'tmp/spec-repo', 'ma/spec-work', '', '', 'tmp/spec-worktree', '0001-widget-flow', 'claude',
			NULL, '2026-07-01T07:00:00Z', '2026-07-01T07:30:00Z', '2026-07-01T07:30:00Z')`,
		`INSERT INTO active_run_locks (target_kind, target_key, run_id, created_at)
		 VALUES ('pr', 'owner/project#feature/review', 'run_v6_active', '2026-07-01T10:00:00Z')`,
		`PRAGMA user_version = 6`,
	}
	for _, statement := range statements {
		if _, err := db.Exec(statement); err != nil {
			t.Fatalf("build v6 fixture: %v", err)
		}
	}
}

func TestOpenMigratesV6RunDatabaseAddingSelectionDefaults(t *testing.T) {
	ctx := context.Background()
	homeDir := t.TempDir()
	buildV6Fixture(t, homeDir)

	runStore := openTestStore(t, ctx, homeDir)
	defer closeStore(t, runStore)

	version, err := runStore.MigrationVersion(ctx)
	if err != nil {
		t.Fatalf("read migration version: %v", err)
	}
	if version != 9 {
		t.Fatalf("expected user_version 9 after migration, got %d", version)
	}
	count, err := runStore.RunCount(ctx)
	if err != nil {
		t.Fatalf("count runs: %v", err)
	}
	if count != 2 {
		t.Fatalf("expected both v6 run rows to survive, got %d", count)
	}

	active, found, err := runStore.ActiveRun(ctx, "owner/project", "feature/review")
	if err != nil {
		t.Fatalf("active review Run lookup after migration: %v", err)
	}
	if !found || active.ID != "run_v6_active" || active.Agent != "codex" || active.WorkDir == "" {
		t.Fatalf("expected v6 review lock and row to survive, found=%v active=%#v", found, active)
	}
	assertRunSelection(t, active, "", "")

	implement, ok, err := runStore.Run(ctx, "run_v6_implement")
	if err != nil || !ok {
		t.Fatalf("expected run_v6_implement to survive, ok=%v err=%v", ok, err)
	}
	if implement.Kind != KindImplement || implement.Agent != "claude" || implement.WorkDir != "tmp/spec-worktree" ||
		implement.CompletedAt == nil {
		t.Fatalf("expected v6 implement fields preserved, got %#v", implement)
	}
	assertRunSelection(t, implement, "", "")

	var rawModel string
	var rawReasoning string
	if err := runStore.db.QueryRowContext(ctx, `SELECT model, reasoning_effort FROM runs WHERE id = 'run_v6_active'`).Scan(&rawModel, &rawReasoning); err != nil {
		t.Fatalf("read migrated selection columns: %v", err)
	}
	if rawModel != "" || rawReasoning != "" {
		t.Fatalf("expected migrated legacy selection defaults to be empty strings, got %q/%q", rawModel, rawReasoning)
	}
}

// buildV7Fixture creates a populated schema v7 Run Database via raw SQL:
// persisted rows have model and reasoning_effort but no owner_pid column.
func buildV7Fixture(t *testing.T, homeDir string) {
	t.Helper()
	path := DatabasePath(homeDir)
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("create fixture directory: %v", err)
	}
	db, err := sql.Open("sqlite", writerDSN(path))
	if err != nil {
		t.Fatalf("open fixture database: %v", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			t.Fatalf("close fixture database: %v", err)
		}
	}()

	statements := []string{
		`CREATE TABLE runs (
			id TEXT PRIMARY KEY,
			kind TEXT NOT NULL,
			state TEXT NOT NULL,
			head_repository TEXT NOT NULL DEFAULT '',
			head_branch TEXT NOT NULL DEFAULT '',
			base_repository TEXT NOT NULL DEFAULT '',
			pr_number TEXT NOT NULL DEFAULT '',
			git_root TEXT NOT NULL,
			local_branch TEXT NOT NULL,
			head_sha TEXT NOT NULL DEFAULT '',
			artifact_dir TEXT NOT NULL DEFAULT '',
			work_dir TEXT,
			spec_slug TEXT NOT NULL DEFAULT '',
			agent TEXT NOT NULL DEFAULT '',
			model TEXT NOT NULL DEFAULT '',
			reasoning_effort TEXT NOT NULL DEFAULT '',
			stop_requested_at TEXT,
			created_at TEXT NOT NULL,
			updated_at TEXT NOT NULL,
			completed_at TEXT NOT NULL DEFAULT ''
		)`,
		`CREATE TABLE active_run_locks ` + activeRunLocksColumns,
		`CREATE INDEX idx_runs_head ON runs (head_repository, head_branch)`,
		`CREATE TABLE interactive_defaults (
			key TEXT PRIMARY KEY,
			value TEXT NOT NULL,
			updated_at TEXT NOT NULL
		)`,
		`CREATE TABLE run_events (
			run_id TEXT NOT NULL,
			cursor INTEGER NOT NULL,
			batch INTEGER NOT NULL DEFAULT 0,
			source TEXT NOT NULL,
			kind TEXT NOT NULL,
			review_issue TEXT NOT NULL DEFAULT '',
			tool_id TEXT NOT NULL DEFAULT '',
			tool_state TEXT NOT NULL DEFAULT '',
			summary TEXT NOT NULL DEFAULT '',
			created_at TEXT NOT NULL,
			payload TEXT NOT NULL DEFAULT '',
			PRIMARY KEY (run_id, cursor),
			FOREIGN KEY (run_id) REFERENCES runs(id) ON DELETE CASCADE
		)`,
		`INSERT INTO runs (id, kind, state, head_repository, head_branch, base_repository,
			pr_number, git_root, local_branch, head_sha, artifact_dir, work_dir, spec_slug, agent,
			model, reasoning_effort, stop_requested_at, created_at, updated_at, completed_at)
		 VALUES ('run_v7_active', 'resolve', 'ResolvingWithAgent', 'owner/project', 'feature/review', 'owner/project',
			'123', 'tmp/repo', 'feature/review', 'abc123', 'tmp/repo/.roundfix', 'tmp/repo/.roundfix/worktrees/run_v7_active', '', 'codex',
			'gpt-5.5', 'xhigh', NULL, '2026-07-01T10:00:00Z', '2026-07-01T10:06:00Z', '')`,
		`INSERT INTO active_run_locks (target_kind, target_key, run_id, created_at)
		 VALUES ('pr', 'owner/project#feature/review', 'run_v7_active', '2026-07-01T10:00:00Z')`,
		`PRAGMA user_version = 7`,
	}
	for _, statement := range statements {
		if _, err := db.Exec(statement); err != nil {
			t.Fatalf("build v7 fixture: %v", err)
		}
	}
}

func TestOpenMigratesV7RunDatabaseAddingOwnerPID(t *testing.T) {
	ctx := context.Background()
	homeDir := t.TempDir()
	buildV7Fixture(t, homeDir)

	runStore := openTestStore(t, ctx, homeDir)
	defer closeStore(t, runStore)

	version, err := runStore.MigrationVersion(ctx)
	if err != nil {
		t.Fatalf("read migration version: %v", err)
	}
	if version != 9 {
		t.Fatalf("expected user_version 9 after migration, got %d", version)
	}

	active, found, err := runStore.ActiveRun(ctx, "owner/project", "feature/review")
	if err != nil {
		t.Fatalf("active review Run lookup after migration: %v", err)
	}
	if !found || active.ID != "run_v7_active" {
		t.Fatalf("expected v7 active Run to survive, found=%v active=%#v", found, active)
	}
	if active.OwnerPID != nil {
		t.Fatalf("expected migrated v7 Run to have no owner PID, got %d", *active.OwnerPID)
	}

	var rawOwnerPID any
	if err := runStore.db.QueryRowContext(ctx, `SELECT owner_pid FROM runs WHERE id = 'run_v7_active'`).Scan(&rawOwnerPID); err != nil {
		t.Fatalf("read migrated owner_pid column: %v", err)
	}
	if rawOwnerPID != nil {
		t.Fatalf("expected migrated owner_pid NULL for legacy row, got %#v", rawOwnerPID)
	}
}

func TestCreateRunRejectsSecondActiveRunForSameSpecTarget(t *testing.T) {
	ctx := context.Background()
	store := openTestStore(t, ctx, t.TempDir())
	defer closeStore(t, store)

	first, err := store.CreateRun(ctx, sampleImplementCreateRunRequest())
	if err != nil {
		t.Fatalf("expected first implement Run, got %v", err)
	}

	_, err = store.CreateRun(ctx, sampleImplementCreateRunRequest())
	var activeErr ActiveRunError
	if !errors.As(err, &activeErr) {
		t.Fatalf("expected ActiveRunError for same Spec target, got %T %v", err, err)
	}
	if activeErr.Existing.ID != first.ID {
		t.Fatalf("expected blocking run %s, got %s", first.ID, activeErr.Existing.ID)
	}
	want := `Active Run already exists for repository "tmp/spec-repo" and Spec "0001-implement-command"; existing run_id=` +
		first.ID + ` state=Active; stop it with: roundfix stop ` + first.ID
	if err.Error() != want {
		t.Fatalf("expected spec-target error naming the work target and blocking run,\nwant %q\ngot  %q", want, err.Error())
	}

	otherSlug := sampleImplementCreateRunRequest()
	otherSlug.SpecSlug = "0002-other-feature"
	if _, err := store.CreateRun(ctx, otherSlug); err != nil {
		t.Fatalf("expected different Spec slug in same repository to pass the lock, got %v", err)
	}

	otherRepo := sampleImplementCreateRunRequest()
	otherRepo.GitRoot = "tmp/other-repo"
	if _, err := store.CreateRun(ctx, otherRepo); err != nil {
		t.Fatalf("expected same Spec slug in different repository to pass the lock, got %v", err)
	}
}

func TestCompletedImplementRunReleasesSpecTargetLock(t *testing.T) {
	ctx := context.Background()
	store := openTestStore(t, ctx, t.TempDir())
	defer closeStore(t, store)

	first, err := store.CreateRun(ctx, sampleImplementCreateRunRequest())
	if err != nil {
		t.Fatalf("expected implement Run, got %v", err)
	}
	if _, err := store.CompleteRun(ctx, first.ID, StateStopped); err != nil {
		t.Fatalf("expected Stopped completion, got %v", err)
	}
	second, err := store.CreateRun(ctx, sampleImplementCreateRunRequest())
	if err != nil {
		t.Fatalf("expected new implement Run after lock release, got %v", err)
	}
	if second.ID == first.ID {
		t.Fatal("expected distinct run id")
	}
}

func TestReviewKindActiveRunErrorTextUnchanged(t *testing.T) {
	ctx := context.Background()
	store := openTestStore(t, ctx, t.TempDir())
	defer closeStore(t, store)

	first, err := store.CreateRun(ctx, sampleCreateRunRequest())
	if err != nil {
		t.Fatalf("expected first Run, got %v", err)
	}
	_, err = store.CreateRun(ctx, sampleCreateRunRequest())
	if err == nil {
		t.Fatal("expected duplicate review Run rejection")
	}
	want := `Active Run already exists for Head Repository "owner/project" and PR Head Branch "feature/review"; existing run_id=` +
		first.ID + ` state=Active`
	if err.Error() != want {
		t.Fatalf("expected review-path error text unchanged,\nwant %q\ngot  %q", want, err.Error())
	}
}

func TestActiveRunInGitRootFindsActiveRunsOfAnyKind(t *testing.T) {
	ctx := context.Background()
	store := openTestStore(t, ctx, t.TempDir())
	defer closeStore(t, store)

	review, err := store.CreateRun(ctx, sampleCreateRunRequest())
	if err != nil {
		t.Fatalf("expected review Run, got %v", err)
	}
	found, ok, err := store.ActiveRunInGitRoot(ctx, review.GitRoot)
	if err != nil {
		t.Fatalf("active run in git root: %v", err)
	}
	if !ok || found.ID != review.ID {
		t.Fatalf("expected review Run active in git root, ok=%v found=%#v", ok, found)
	}

	if _, ok, err := store.ActiveRunInGitRoot(ctx, "tmp/elsewhere"); err != nil || ok {
		t.Fatalf("expected no Active Run in unrelated git root, ok=%v err=%v", ok, err)
	}

	if _, err := store.CompleteRun(ctx, review.ID, StateClean); err != nil {
		t.Fatalf("complete review Run: %v", err)
	}
	if _, ok, err := store.ActiveRunInGitRoot(ctx, review.GitRoot); err != nil || ok {
		t.Fatalf("expected completed Run to leave git root free, ok=%v err=%v", ok, err)
	}

	implementReq := sampleImplementCreateRunRequest()
	implementReq.GitRoot = review.GitRoot
	implementRun, err := store.CreateRun(ctx, implementReq)
	if err != nil {
		t.Fatalf("expected implement Run, got %v", err)
	}
	found, ok, err = store.ActiveRunInGitRoot(ctx, review.GitRoot)
	if err != nil {
		t.Fatalf("active run in git root after implement create: %v", err)
	}
	if !ok || found.ID != implementRun.ID || found.Kind != KindImplement {
		t.Fatalf("expected implement Run active in git root, ok=%v found=%#v", ok, found)
	}
}

func TestCreateRunValidatesRequiredFieldsByKind(t *testing.T) {
	ctx := context.Background()
	store := openTestStore(t, ctx, t.TempDir())
	defer closeStore(t, store)

	tests := []struct {
		name    string
		request func() CreateRunRequest
		wantErr string
	}{
		{
			name: "implement missing Spec slug",
			request: func() CreateRunRequest {
				req := sampleImplementCreateRunRequest()
				req.SpecSlug = ""
				return req
			},
			wantErr: "Spec slug is required to create a Run",
		},
		{
			name: "implement missing Git root",
			request: func() CreateRunRequest {
				req := sampleImplementCreateRunRequest()
				req.GitRoot = ""
				return req
			},
			wantErr: "Git root is required to create a Run",
		},
		{
			name: "implement missing local branch",
			request: func() CreateRunRequest {
				req := sampleImplementCreateRunRequest()
				req.LocalBranch = ""
				return req
			},
			wantErr: "local branch is required to create a Run",
		},
		{
			name: "review missing pull request",
			request: func() CreateRunRequest {
				req := sampleCreateRunRequest()
				req.PRNumber = ""
				return req
			},
			wantErr: "pull request is required to create a Run",
		},
		{
			name: "review missing Head Repository",
			request: func() CreateRunRequest {
				req := sampleCreateRunRequest()
				req.HeadRepository = ""
				return req
			},
			wantErr: "Head Repository is required to create a Run",
		},
		{
			name: "review missing HEAD",
			request: func() CreateRunRequest {
				req := sampleCreateRunRequest()
				req.HeadSHA = ""
				return req
			},
			wantErr: "HEAD is required to create a Run",
		},
		{
			name: "unknown kind",
			request: func() CreateRunRequest {
				req := sampleCreateRunRequest()
				req.Kind = "deploy"
				return req
			},
			wantErr: `Run kind "deploy" is invalid`,
		},
		{
			name: "empty kind",
			request: func() CreateRunRequest {
				req := sampleCreateRunRequest()
				req.Kind = ""
				return req
			},
			wantErr: "Run kind is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := store.CreateRun(ctx, tt.request())
			if err == nil {
				t.Fatalf("expected validation error %q, got nil", tt.wantErr)
			}
			if err.Error() != tt.wantErr {
				t.Fatalf("expected error %q, got %q", tt.wantErr, err.Error())
			}
		})
	}
}

func sampleImplementCreateRunRequest() CreateRunRequest {
	return CreateRunRequest{
		Kind:            KindImplement,
		GitRoot:         "tmp/spec-repo",
		LocalBranch:     "ma/implement-spec",
		SpecSlug:        "0001-implement-command",
		ArtifactDir:     "tmp/spec-repo/.roundfix",
		Agent:           "codex",
		Model:           "gpt-5.5",
		ReasoningEffort: "xhigh",
		OwnerPID:        os.Getpid(),
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
