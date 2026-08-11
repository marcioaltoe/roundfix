package store

import (
	"bytes"
	"context"
	"crypto/sha256"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"roundfix/internal/runevent"
	"runtime"
	"slices"
	"strings"
	"testing"
	"time"
)

func sampleRunEvent(runID string, summary string) runevent.RunEvent {
	return runevent.RunEvent{
		RunID:   runID,
		Batch:   1,
		Source:  runevent.SourceAgent,
		Kind:    runevent.KindAgentMessage,
		Summary: summary,
		Time:    time.Date(2026, 6, 10, 12, 0, 0, 0, time.UTC),
		Payload: []byte(`{"sessionId":"sess-1","update":{"sessionUpdate":"agent_message_chunk"}}`),
	}
}

func TestMigrationAddsJournalSchemaAndWALMode(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	homeDir := t.TempDir()
	store := openTestStore(t, ctx, homeDir)

	var journalMode string
	if err := store.db.QueryRowContext(ctx, `PRAGMA journal_mode`).Scan(&journalMode); err != nil {
		t.Fatalf("read journal mode: %v", err)
	}
	if !strings.EqualFold(journalMode, "wal") {
		t.Fatalf("expected WAL journal mode, got %q", journalMode)
	}
	var busyTimeout int
	if err := store.db.QueryRowContext(ctx, `PRAGMA busy_timeout`).Scan(&busyTimeout); err != nil {
		t.Fatalf("read busy timeout: %v", err)
	}
	if busyTimeout <= 0 {
		t.Fatalf("expected busy timeout configured, got %d", busyTimeout)
	}
	closeStore(t, store)

	// Reopening an existing Run Database must keep working with the journal.
	reopened := openTestStore(t, ctx, homeDir)
	defer closeStore(t, reopened)
	run, err := reopened.CreateRun(ctx, sampleCreateRunRequest())
	if err != nil {
		t.Fatalf("create run after reopen: %v", err)
	}
	if _, err := reopened.AppendRunEvent(ctx, sampleRunEvent(run.ID, "after reopen")); err != nil {
		t.Fatalf("append after reopen: %v", err)
	}
}

func TestAppendRunEventAllocatesMonotonicCursorsPerRun(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	store := openTestStore(t, ctx, t.TempDir())
	defer closeStore(t, store)

	first, err := store.CreateRun(ctx, sampleCreateRunRequest())
	if err != nil {
		t.Fatalf("create first run: %v", err)
	}
	secondReq := sampleCreateRunRequest()
	secondReq.HeadBranch = "feature/other"
	second, err := store.CreateRun(ctx, secondReq)
	if err != nil {
		t.Fatalf("create second run: %v", err)
	}

	for index, expected := range []int64{1, 2, 3} {
		cursor, err := store.AppendRunEvent(ctx, sampleRunEvent(first.ID, "event"))
		if err != nil {
			t.Fatalf("append event %d: %v", index, err)
		}
		if cursor != expected {
			t.Fatalf("expected cursor %d, got %d", expected, cursor)
		}
	}
	cursor, err := store.AppendRunEvent(ctx, sampleRunEvent(second.ID, "event"))
	if err != nil {
		t.Fatalf("append to second run: %v", err)
	}
	if cursor != 1 {
		t.Fatalf("expected independent per-Run cursor 1, got %d", cursor)
	}

	cursors, err := store.AppendRunEvents(ctx, []runevent.RunEvent{
		sampleRunEvent(first.ID, "batch a"),
		sampleRunEvent(first.ID, "batch b"),
	})
	if err != nil {
		t.Fatalf("append batch: %v", err)
	}
	if len(cursors) != 2 || cursors[0] != 4 || cursors[1] != 5 {
		t.Fatalf("expected batch cursors [4 5], got %v", cursors)
	}
}

func TestAppendRunEventToMissingRunFailsClearly(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	store := openTestStore(t, ctx, t.TempDir())
	defer closeStore(t, store)

	_, err := store.AppendRunEvent(ctx, sampleRunEvent("run_missing", "orphan"))

	if err == nil {
		t.Fatal("expected append to missing Run to fail")
	}
	if !strings.Contains(err.Error(), `Run "run_missing" does not exist`) {
		t.Fatalf("expected clear missing-Run error, got %v", err)
	}
}

func TestPruneTerminalRunsDeletesOnlyEligibleJournalRows(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	runStore := openTestStore(t, ctx, t.TempDir())
	defer closeStore(t, runStore)

	cutoff := time.Date(2026, 7, 1, 12, 0, 0, 0, time.UTC)
	oldCreated := cutoff.Add(-72 * time.Hour)
	oldCompleted := cutoff.Add(-24 * time.Hour)
	recentCompleted := cutoff.Add(24 * time.Hour)

	tests := []struct {
		name                   string
		branch                 string
		terminalState          string
		nonTerminalState       string
		terminalEmptyCompleted bool
		completedAt            time.Time
		eventCount             int
		wantEventCount         int
		wantPruned             bool
		wantLock               bool
	}{
		{
			name:           "terminal clean before cutoff",
			branch:         "old-clean",
			terminalState:  StateClean,
			completedAt:    oldCompleted,
			eventCount:     2,
			wantPruned:     true,
			wantEventCount: 0,
		},
		{
			name:           "terminal clean unverified before cutoff",
			branch:         "old-clean-unverified",
			terminalState:  StateCleanUnverified,
			completedAt:    oldCompleted.Add(30 * time.Second),
			eventCount:     1,
			wantPruned:     true,
			wantEventCount: 0,
		},
		{
			name:           "terminal Review Skipped before cutoff",
			branch:         "old-review-skipped",
			terminalState:  StateReviewSkipped,
			completedAt:    oldCompleted.Add(45 * time.Second),
			eventCount:     1,
			wantPruned:     true,
			wantEventCount: 0,
		},
		{
			name:           "terminal unresolved before cutoff",
			branch:         "old-unresolved",
			terminalState:  StateUnresolved,
			completedAt:    oldCompleted.Add(time.Minute),
			eventCount:     1,
			wantPruned:     true,
			wantEventCount: 0,
		},
		{
			name:           "terminal clean after cutoff",
			branch:         "recent-clean",
			terminalState:  StateClean,
			completedAt:    recentCompleted,
			eventCount:     1,
			wantEventCount: 1,
		},
		{
			name:           "active run with old created at",
			branch:         "active-old",
			eventCount:     1,
			wantEventCount: 1,
			wantLock:       true,
		},
		{
			name:             "non-terminal run with old created at",
			branch:           "resolving-old",
			nonTerminalState: StateResolvingWithAgent,
			eventCount:       1,
			wantEventCount:   1,
			wantLock:         true,
		},
		{
			name:                   "terminal state with empty completed at",
			branch:                 "empty-completed",
			terminalEmptyCompleted: true,
			eventCount:             1,
			wantEventCount:         1,
			wantLock:               true,
		},
	}

	runByName := map[string]Run{}
	wantPrunedRunIDs := []string{}
	wantLockCount := 0
	for _, tt := range tests {
		createdAt := oldCreated.Add(time.Duration(len(runByName)) * time.Minute)
		runStore.now = func() time.Time { return createdAt }
		req := sampleCreateRunRequest()
		req.HeadBranch = "feature/" + tt.branch
		req.PRNumber = tt.branch
		run, err := runStore.CreateRun(ctx, req)
		if err != nil {
			t.Fatalf("%s: create Run: %v", tt.name, err)
		}
		for eventIndex := 0; eventIndex < tt.eventCount; eventIndex++ {
			if _, err := runStore.AppendRunEvent(ctx, sampleRunEvent(run.ID, tt.branch)); err != nil {
				t.Fatalf("%s: append Run Event: %v", tt.name, err)
			}
		}
		switch {
		case tt.terminalState != "":
			completedAt := tt.completedAt
			runStore.now = func() time.Time { return completedAt }
			completed, completeErr := runStore.CompleteRun(ctx, run.ID, tt.terminalState)
			err = completeErr
			if err != nil {
				t.Fatalf("%s: complete Run: %v", tt.name, err)
			}
			run = completed.Run
		case tt.nonTerminalState != "":
			if err := runStore.UpdateRunState(ctx, run.ID, tt.nonTerminalState); err != nil {
				t.Fatalf("%s: update Run state: %v", tt.name, err)
			}
		case tt.terminalEmptyCompleted:
			if _, err := runStore.db.ExecContext(ctx, `
UPDATE runs
SET state = ?, updated_at = ?, completed_at = ''
WHERE id = ?`,
				StateClean,
				formatTime(oldCompleted),
				run.ID,
			); err != nil {
				t.Fatalf("%s: seed empty completed_at terminal Run: %v", tt.name, err)
			}
			run.State = StateClean
		}

		runByName[tt.name] = run
		if tt.wantPruned {
			wantPrunedRunIDs = append(wantPrunedRunIDs, run.ID)
		}
		if tt.wantLock {
			wantLockCount++
		}
	}

	initialRunCount, err := runStore.RunCount(ctx)
	if err != nil {
		t.Fatalf("count seeded Runs: %v", err)
	}
	initialLockCount := countActiveRunLocks(t, ctx, runStore)

	result, err := runStore.PruneTerminalRuns(ctx, time.Time{})
	if err != nil {
		t.Fatalf("prune terminal Runs with zero cutoff: %v", err)
	}
	if len(result.RunIDs) != 0 || result.Events != 0 {
		t.Fatalf("expected empty zero-cutoff prune result, got %#v", result)
	}
	for _, tt := range tests {
		run := runByName[tt.name]
		if got := countRunEvents(t, ctx, runStore, run.ID); got != tt.eventCount {
			t.Fatalf("%s: expected zero cutoff to keep %d Run Events, got %d", tt.name, tt.eventCount, got)
		}
	}

	result, err = runStore.PruneTerminalRuns(ctx, cutoff)

	if err != nil {
		t.Fatalf("prune terminal Runs: %v", err)
	}
	gotRunIDs := append([]string(nil), result.RunIDs...)
	slices.Sort(gotRunIDs)
	slices.Sort(wantPrunedRunIDs)
	if !slices.Equal(gotRunIDs, wantPrunedRunIDs) {
		t.Fatalf("expected pruned Run ids %v, got %v", wantPrunedRunIDs, result.RunIDs)
	}
	if result.Events != 5 {
		t.Fatalf("expected 5 pruned Run Events, got %d", result.Events)
	}
	for _, tt := range tests {
		run := runByName[tt.name]
		if got := countRunEvents(t, ctx, runStore, run.ID); got != tt.wantEventCount {
			t.Fatalf("%s: expected %d remaining Run Events, got %d", tt.name, tt.wantEventCount, got)
		}
		if _, ok, err := runStore.Run(ctx, run.ID); err != nil || !ok {
			t.Fatalf("%s: expected Run row to survive, ok=%v err=%v", tt.name, ok, err)
		}
	}
	if got, err := runStore.RunCount(ctx); err != nil || got != initialRunCount {
		t.Fatalf("expected all Run rows to survive, got count=%d err=%v", got, err)
	}
	if got := countActiveRunLocks(t, ctx, runStore); got != initialLockCount || got != wantLockCount {
		t.Fatalf("expected Active Run locks untouched at %d, got %d", wantLockCount, got)
	}
}

func TestPruneTerminalRunsNoOpsWhenCutoffSelectsNothing(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	runStore := openTestStore(t, ctx, t.TempDir())
	defer closeStore(t, runStore)

	cutoff := time.Date(2026, 7, 1, 12, 0, 0, 0, time.UTC)
	runStore.now = func() time.Time { return cutoff.Add(-time.Hour) }
	run, err := runStore.CreateRun(ctx, sampleCreateRunRequest())
	if err != nil {
		t.Fatalf("create Run: %v", err)
	}
	if _, err := runStore.AppendRunEvent(ctx, sampleRunEvent(run.ID, "recent terminal")); err != nil {
		t.Fatalf("append Run Event: %v", err)
	}
	runStore.now = func() time.Time { return cutoff.Add(time.Hour) }
	if _, err := runStore.CompleteRun(ctx, run.ID, StateClean); err != nil {
		t.Fatalf("complete Run: %v", err)
	}

	result, err := runStore.PruneTerminalRuns(ctx, cutoff)

	if err != nil {
		t.Fatalf("prune terminal Runs: %v", err)
	}
	if len(result.RunIDs) != 0 || result.Events != 0 {
		t.Fatalf("expected empty prune result, got %#v", result)
	}
	if got := countRunEvents(t, ctx, runStore, run.ID); got != 1 {
		t.Fatalf("expected recent terminal journal kept, got %d events", got)
	}
	if count, err := runStore.RunCount(ctx); err != nil || count != 1 {
		t.Fatalf("expected Run row to survive, count=%d err=%v", count, err)
	}
}

// TestRetentionScanIsBoundedByCandidates proves that eligibility work is
// bounded by the candidate set — terminal Runs — rather than by the event
// table. It seeds a batch of identical terminal Runs and a competing large
// body of unrelated run_events, then asserts that TerminalRunPruneCandidates
// yields exactly the eligible candidate IDs with their own event counts,
// unchanged by the unrelated rows. This is behavioral coverage; the structural
// concern is expressed through the query's SQL predicate (completed_at bias),
// which the bounded candidate query observes by construction.
func TestRetentionScanIsBoundedByCandidates(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	runStore := openTestStore(t, ctx, t.TempDir())
	defer closeStore(t, runStore)

	cutoff := time.Date(2026, 7, 1, 12, 0, 0, 0, time.UTC)
	runStore.now = func() time.Time { return cutoff.Add(-2 * time.Hour) }

	// Seed several identical terminal Runs that are eligible for pruning and
	// one recent terminal Run that is not.
	const eligibleCount = 3
	wantIDs := make([]string, 0, eligibleCount)
	for i := 0; i < eligibleCount; i++ {
		run, err := runStore.CreateRun(ctx, sampleImplementCreateRunRequest())
		if err != nil {
			t.Fatalf("create eligible Run %d: %v", i, err)
		}
		for j := 0; j < 5; j++ {
			if _, err := runStore.AppendRunEvent(ctx, sampleRunEvent(run.ID, "eligible event")); err != nil {
				t.Fatalf("append eligible Run %d event %d: %v", i, j, err)
			}
		}
		if _, err := runStore.CompleteRun(ctx, run.ID, StateClean); err != nil {
			t.Fatalf("complete eligible Run %d: %v", i, err)
		}
		wantIDs = append(wantIDs, run.ID)
	}

	runStore.now = func() time.Time { return cutoff.Add(time.Hour) }
	recent, err := runStore.CreateRun(ctx, sampleImplementCreateRunRequest())
	if err != nil {
		t.Fatalf("create recent Run: %v", err)
	}
	for j := 0; j < 5; j++ {
		if _, err := runStore.AppendRunEvent(ctx, sampleRunEvent(recent.ID, "recent event")); err != nil {
			t.Fatalf("append recent Run event %d: %v", j, err)
		}
	}
	if _, err := runStore.CompleteRun(ctx, recent.ID, StateClean); err != nil {
		t.Fatalf("complete recent Run: %v", err)
	}

	// Seed a competing body of unrelated run_events on a Run that never
	// completes, so the candidate query cannot be confused by table volume.
	noisy, err := runStore.CreateRun(ctx, sampleImplementCreateRunRequest())
	if err != nil {
		t.Fatalf("create noisy Run: %v", err)
	}
	for j := 0; j < 200; j++ {
		if _, err := runStore.AppendRunEvent(ctx, sampleRunEvent(noisy.ID, "noise event")); err != nil {
			t.Fatalf("append noisy Run event %d: %v", j, err)
		}
	}

	candidates, err := runStore.TerminalRunPruneCandidates(ctx, cutoff)
	if err != nil {
		t.Fatalf("list terminal Run prune candidates: %v", err)
	}
	if len(candidates) != len(wantIDs) {
		t.Fatalf("expected %d eligible candidates, got %d", len(wantIDs), len(candidates))
	}
	gotIDs := make([]string, 0, len(candidates))
	gotCounts := map[string]int{}
	for _, candidate := range candidates {
		gotIDs = append(gotIDs, candidate.RunID)
		gotCounts[candidate.RunID] = candidate.Events
	}
	slices.Sort(gotIDs)
	slices.Sort(wantIDs)
	if !slices.Equal(gotIDs, wantIDs) {
		t.Fatalf("eligible candidate ids %v, want %v", gotIDs, wantIDs)
	}
	for _, runID := range wantIDs {
		if gotCounts[runID] != 5 {
			t.Fatalf("eligible Run %q event count = %d, want 5 (unrelated rows must not leak)", runID, gotCounts[runID])
		}
	}
	if _, ok := gotCounts[recent.ID]; ok {
		t.Fatalf("recent Run %q must not be eligible at cutoff", recent.ID)
	}
}

// TestRetentionScanOutsideWriteTransaction proves the eligibility scan left
// the write transaction: it completes while the machine-wide advisory write
// lock is held from an independent descriptor. With the lock held, a
// TerminalRunPruneCandidates scan must complete (the scan never needs the
// writer), while a prune that reaches the write path must block until the
// deadline.
func TestRetentionScanOutsideWriteTransaction(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	homeDir := t.TempDir()
	runStore := openTestStore(t, ctx, homeDir)
	defer closeStore(t, runStore)

	// Seed a terminal Run completed before the cutoff so it is eligible for
	// pruning and the prune reaches the write transaction path.
	cutoff := time.Date(2026, 7, 1, 12, 0, 0, 0, time.UTC)
	runStore.now = func() time.Time { return cutoff.Add(-time.Hour) }
	run, err := runStore.CreateRun(ctx, sampleCreateRunRequest())
	if err != nil {
		t.Fatalf("create Run: %v", err)
	}
	if _, err := runStore.AppendRunEvent(ctx, sampleRunEvent(run.ID, "retention event")); err != nil {
		t.Fatalf("append Run Event: %v", err)
	}
	if _, err := runStore.CompleteRun(ctx, run.ID, StateClean); err != nil {
		t.Fatalf("complete Run: %v", err)
	}

	// Hold the write lock from an independent descriptor. flock is owned by
	// the open file description, so reusing runStore.writeLockFile would grant
	// the lock to the Store as well and prove nothing.
	holder, err := openWriteLockFile(DatabasePath(homeDir))
	if err != nil {
		t.Fatalf("open independent write lock: %v", err)
	}
	defer func() {
		_ = releaseWriteLock(holder)
		_ = holder.Close()
	}()
	if err := acquireWriteLock(holder, ctx); err != nil {
		t.Fatalf("acquire write lock: %v", err)
	}

	// The eligibility scan must complete before the deadline while the lock is
	// held, because it never needs the writer.
	scanCtx, cancelScan := context.WithTimeout(ctx, 500*time.Millisecond)
	defer cancelScan()
	candidates, err := runStore.TerminalRunPruneCandidates(scanCtx, cutoff)
	if err != nil {
		t.Fatalf("eligibility scan blocked on the held write lock: %v", err)
	}
	if len(candidates) != 1 || candidates[0].RunID != run.ID || candidates[0].Events != 1 {
		t.Fatalf("eligible candidate = %+v, want one Run with one event", candidates)
	}

	// A prune with eligible rows reaches the write path and must block on the
	// held lock until the deadline expires.
	pruneCtx, cancelPrune := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancelPrune()
	start := time.Now()
	_, pruneErr := runStore.PruneTerminalRuns(pruneCtx, cutoff)
	if pruneErr == nil {
		t.Fatal("expected eligible prune to block on the held write lock")
	}
	if got := time.Since(start); got < 150*time.Millisecond {
		t.Fatalf("eligible prune returned too fast (%v); did not block on the held write lock", got)
	}
	// The recent terminal journal must be preserved: the prune never ran.
	if got := countRunEvents(t, ctx, runStore, run.ID); got != 1 {
		t.Fatalf("expected preserved terminal journal, got %d events", got)
	}
}

func TestDurableTableLifecyclePolicyCoversEveryTable(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	runStore := openTestStore(t, ctx, t.TempDir())
	defer closeStore(t, runStore)

	policy := readDurableTableLifecyclePolicy(t)

	if err := validateDurableTableLifecyclePolicy(ctx, runStore.db, policy); err != nil {
		t.Fatal(err)
	}
}

func TestDurableTableLifecyclePolicyRejectsUnstatedFixtureTable(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	runStore := openTestStore(t, ctx, t.TempDir())
	defer closeStore(t, runStore)
	if _, err := runStore.db.ExecContext(ctx, `CREATE TABLE lifecycle_fixture (id TEXT PRIMARY KEY)`); err != nil {
		t.Fatalf("create ungoverned durable table fixture: %v", err)
	}

	err := validateDurableTableLifecyclePolicy(ctx, runStore.db, readDurableTableLifecyclePolicy(t))

	if err == nil {
		t.Fatal("expected an ungoverned durable table to fail lifecycle policy validation")
	}
	if !strings.Contains(err.Error(), `durable table "lifecycle_fixture" has no lifecycle policy`) {
		t.Fatalf("expected missing fixture policy diagnostic, got %v", err)
	}
}

func TestRetentionPreservesRunLifecycleRecords(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	runStore := openTestStore(t, ctx, t.TempDir())
	defer closeStore(t, runStore)

	cutoff := time.Date(2026, 7, 1, 12, 0, 0, 0, time.UTC)
	runStore.now = func() time.Time { return cutoff.Add(-2 * time.Hour) }
	terminalRun, err := runStore.CreateRun(ctx, sampleImplementCreateRunRequest())
	if err != nil {
		t.Fatalf("create terminal Run: %v", err)
	}
	if _, err := runStore.AppendAgentSelectionAttempt(ctx, AgentSelectionAttemptRequest{
		RunID: terminalRun.ID, ScopeKind: AgentSelectionScopeTask, ScopeID: "task_04",
		Category: "backend", ProfileSource: "project", Attempt: 1,
		SelectionRole: AgentSelectionRolePreferred, Runtime: "codex", Model: "gpt-5.6-sol",
		ReasoningEffort: "high", Status: AgentSelectionStatusClosed,
		Time: cutoff.Add(-90 * time.Minute),
	}); err != nil {
		t.Fatalf("append terminal Run Agent Selection: %v", err)
	}
	runStore.now = func() time.Time { return cutoff.Add(-time.Hour) }
	if _, err := runStore.CompleteRun(ctx, terminalRun.ID, StateClean); err != nil {
		t.Fatalf("complete terminal Run: %v", err)
	}

	activeRequest := sampleImplementCreateRunRequest()
	activeRequest.SpecSlug = "active-lifecycle-fixture"
	runStore.now = func() time.Time { return cutoff.Add(-30 * time.Minute) }
	activeRun, err := runStore.CreateRun(ctx, activeRequest)
	if err != nil {
		t.Fatalf("create Active Run: %v", err)
	}
	if _, err := runStore.AppendAgentSelectionAttempt(ctx, AgentSelectionAttemptRequest{
		RunID: activeRun.ID, ScopeKind: AgentSelectionScopeTask, ScopeID: "task_04",
		Category: "backend", ProfileSource: "project", Attempt: 1,
		SelectionRole: AgentSelectionRolePreferred, Runtime: "codex", Model: "gpt-5.6-sol",
		ReasoningEffort: "high", Status: AgentSelectionStatusActive,
		Time: cutoff.Add(-15 * time.Minute),
	}); err != nil {
		t.Fatalf("append Active Run Agent Selection: %v", err)
	}

	initialRuns, err := runStore.RunCount(ctx)
	if err != nil {
		t.Fatalf("count Runs before retention: %v", err)
	}
	initialLocks := countActiveRunLocks(t, ctx, runStore)
	initialSelections := countAgentSelections(t, ctx, runStore)

	result, err := runStore.PruneTerminalRuns(ctx, cutoff)

	if err != nil {
		t.Fatalf("prune terminal Run Events: %v", err)
	}
	if !slices.Equal(result.RunIDs, []string{terminalRun.ID}) || result.Events != 1 {
		t.Fatalf("expected only the terminal Run journal pruned, got %#v", result)
	}
	if got := countRunEvents(t, ctx, runStore, terminalRun.ID); got != 0 {
		t.Fatalf("expected terminal Run journal pruned, got %d events", got)
	}
	if got := countRunEvents(t, ctx, runStore, activeRun.ID); got != 1 {
		t.Fatalf("expected Active Run journal preserved, got %d events", got)
	}
	if got, err := runStore.RunCount(ctx); err != nil || got != initialRuns {
		t.Fatalf("expected compact Run index preserved at %d rows, got %d err=%v", initialRuns, got, err)
	}
	if got := countActiveRunLocks(t, ctx, runStore); got != initialLocks {
		t.Fatalf("expected Active Run locks preserved at %d rows, got %d", initialLocks, got)
	}
	if got := countAgentSelections(t, ctx, runStore); got != initialSelections {
		t.Fatalf("expected Agent Selection records preserved at %d rows, got %d", initialSelections, got)
	}
}

type durableTableLifecyclePolicy struct {
	owner string
	rule  string
}

func readDurableTableLifecyclePolicy(t *testing.T) map[string]durableTableLifecyclePolicy {
	t.Helper()
	_, testFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("locate lifecycle policy: runtime.Caller failed")
	}
	policyPath := filepath.Join(filepath.Dir(testFile), "..", "..", "docs", "user-guide", "run-database-lifecycle.md")
	content, err := os.ReadFile(policyPath)
	if err != nil {
		t.Fatalf("read durable table lifecycle policy: %v", err)
	}

	const begin = "<!-- durable-table-lifecycle:begin -->"
	const end = "<!-- durable-table-lifecycle:end -->"
	policyContent := string(content)
	beginIndex := strings.Index(policyContent, begin)
	if beginIndex < 0 {
		t.Fatalf("read durable table lifecycle policy: marker %q is missing", begin)
	}
	policyBlock := policyContent[beginIndex+len(begin):]
	endIndex := strings.Index(policyBlock, end)
	if endIndex < 0 {
		t.Fatalf("read durable table lifecycle policy: marker %q is missing", end)
	}
	policyBlock = policyBlock[:endIndex]

	policy := map[string]durableTableLifecyclePolicy{}
	for _, line := range strings.Split(policyBlock, "\n") {
		line = strings.TrimSpace(line)
		if !strings.HasPrefix(line, "| `") {
			continue
		}
		columns := strings.Split(line, "|")
		if len(columns) != 5 {
			t.Fatalf("read durable table lifecycle policy: malformed table row %q", line)
		}
		table := strings.Trim(strings.TrimSpace(columns[1]), "`")
		entry := durableTableLifecyclePolicy{
			owner: strings.TrimSpace(columns[2]),
			rule:  strings.TrimSpace(columns[3]),
		}
		if table == "" || entry.owner == "" || entry.rule == "" {
			t.Fatalf("read durable table lifecycle policy: table, owner, and retention rule must be non-empty in %q", line)
		}
		if _, exists := policy[table]; exists {
			t.Fatalf("read durable table lifecycle policy: duplicate policy for table %q", table)
		}
		policy[table] = entry
	}
	if len(policy) == 0 {
		t.Fatal("read durable table lifecycle policy: no table policies found")
	}
	return policy
}

func validateDurableTableLifecyclePolicy(ctx context.Context, db *sql.DB, policy map[string]durableTableLifecyclePolicy) error {
	rows, err := db.QueryContext(ctx, `
SELECT name
FROM sqlite_schema
WHERE type = 'table' AND name NOT LIKE 'sqlite_%'
ORDER BY name`)
	if err != nil {
		return fmt.Errorf("list durable Run Database tables: %w", err)
	}
	defer func() {
		_ = rows.Close()
	}()

	tables := map[string]struct{}{}
	for rows.Next() {
		var table string
		if err := rows.Scan(&table); err != nil {
			return fmt.Errorf("scan durable Run Database table: %w", err)
		}
		tables[table] = struct{}{}
		if _, ok := policy[table]; !ok {
			return fmt.Errorf("durable table %q has no lifecycle policy", table)
		}
	}
	if err := rows.Err(); err != nil {
		return fmt.Errorf("iterate durable Run Database tables: %w", err)
	}
	for table := range policy {
		if _, ok := tables[table]; !ok {
			return fmt.Errorf("lifecycle policy names unknown durable table %q", table)
		}
	}
	return nil
}

func countAgentSelections(t *testing.T, ctx context.Context, store *Store) int {
	t.Helper()
	var count int
	if err := store.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM run_agent_selections`).Scan(&count); err != nil {
		t.Fatalf("count Agent Selection records: %v", err)
	}
	return count
}

func TestStorageReportReconcilesMeasuredTotalsWithoutMutation(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	homeDir := t.TempDir()
	existingRepository := t.TempDir()
	missingRepository := filepath.Join(t.TempDir(), "removed-repository")
	artifactRoot := filepath.Join(t.TempDir(), "artifacts")
	missingArtifactRoot := filepath.Join(t.TempDir(), "removed-artifacts")

	runStore := openTestStore(t, ctx, homeDir)
	firstRequest := sampleCreateRunRequest()
	firstRequest.GitRoot = existingRepository
	firstRequest.ArtifactDir = artifactRoot
	first, err := runStore.CreateRun(ctx, firstRequest)
	if err != nil {
		t.Fatalf("create first Run: %v", err)
	}
	if _, err := runStore.AppendRunEvent(ctx, sampleRunEvent(first.ID, "measured event")); err != nil {
		t.Fatalf("append measured Run Event: %v", err)
	}
	for index := 0; index < 32; index++ {
		event := sampleRunEvent(first.ID, fmt.Sprintf("prunable measured event %d", index))
		event.Payload = bytes.Repeat([]byte("x"), 8192)
		if _, err := runStore.AppendRunEvent(ctx, event); err != nil {
			t.Fatalf("append prunable measured Run Event %d: %v", index, err)
		}
	}
	if _, err := runStore.CompleteRun(ctx, first.ID, StateClean); err != nil {
		t.Fatalf("complete first Run: %v", err)
	}
	if _, err := runStore.PruneTerminalRuns(ctx, time.Now().UTC().Add(time.Hour)); err != nil {
		t.Fatalf("create free pages by pruning measured Run Events: %v", err)
	}

	secondRequest := sampleCreateRunRequest()
	secondRequest.HeadBranch = "feature/missing-root"
	secondRequest.PRNumber = "456"
	secondRequest.GitRoot = missingRepository
	secondRequest.ArtifactDir = missingArtifactRoot
	if _, err := runStore.CreateRun(ctx, secondRequest); err != nil {
		t.Fatalf("create Run with missing recorded roots: %v", err)
	}
	closeStore(t, runStore)

	artifactPath := filepath.Join(artifactRoot, "runs", first.ID, "agent", "batch-001.log")
	if err := os.MkdirAll(filepath.Dir(artifactPath), 0o755); err != nil {
		t.Fatalf("create artifact directory: %v", err)
	}
	if err := os.WriteFile(artifactPath, []byte("measured artifact"), 0o644); err != nil {
		t.Fatalf("write measured artifact: %v", err)
	}
	databaseBefore := storageTreeSnapshot(t, filepath.Join(homeDir, roundfixHomeDir))
	artifactBefore := storageTreeSnapshot(t, artifactRoot)

	reader, err := OpenStorageReader(ctx, homeDir)
	if err != nil {
		t.Fatalf("open Run Database reader: %v", err)
	}
	report, err := reader.StorageReport(ctx)
	if err != nil {
		t.Fatalf("measure storage report: %v", err)
	}
	pageSize, err := storagePragmaInt64(ctx, reader.db, "page_size")
	if err != nil {
		t.Fatalf("measure SQLite page size: %v", err)
	}
	closeStore(t, reader)
	if report.ReconciliationToleranceBytes != pageSize {
		t.Fatalf("expected tolerance to equal measured SQLite page size %d, got %d", pageSize, report.ReconciliationToleranceBytes)
	}
	if report.DatabaseFreeBytes == 0 {
		t.Fatal("expected the reconciliation fixture to contain measured free pages")
	}
	if !strings.Contains(report.ReconciliationToleranceReason, "one SQLite page") {
		t.Fatalf("expected page-sized tolerance reason, got %q", report.ReconciliationToleranceReason)
	}

	var tableBytes int64
	var runRows int64
	for _, group := range report.Tables {
		tableBytes += group.Bytes
		if group.Table == "runs" {
			runRows = group.Rows
		}
	}
	databaseDifference := absInt64(report.DatabaseBytes - (tableBytes + report.DatabaseFreeBytes))
	if databaseDifference > report.ReconciliationToleranceBytes {
		t.Fatalf("database groups differ from measured file by %d bytes, tolerance %d: %#v", databaseDifference, report.ReconciliationToleranceBytes, report)
	}
	var artifactBytes int64
	for _, group := range report.ArtifactRoots {
		artifactBytes += group.Bytes
	}
	if artifactBytes != report.ArtifactBytes {
		t.Fatalf("Artifact Root groups total %d bytes, measured artifacts total %d", artifactBytes, report.ArtifactBytes)
	}
	var repositoryRows int64
	var repositoryBytes int64
	for _, group := range report.Repositories {
		repositoryRows += group.Rows
		repositoryBytes += group.Bytes
	}
	var stateRows int64
	var stateBytes int64
	for _, group := range report.States {
		stateRows += group.Rows
		stateBytes += group.Bytes
	}
	var artifactRootRows int64
	for _, group := range report.ArtifactRoots {
		artifactRootRows += group.Rows
	}
	if repositoryRows != runRows || stateRows != runRows || artifactRootRows != runRows {
		t.Fatalf("Run row groupings do not reconcile: table=%d repositories=%d states=%d Artifact Roots=%d", runRows, repositoryRows, stateRows, artifactRootRows)
	}
	if repositoryBytes != stateBytes {
		t.Fatalf("Run artifact byte groupings do not reconcile: repositories=%d states=%d", repositoryBytes, stateBytes)
	}
	if !storageRepositoryGroup(report.Repositories, missingRepository).Missing {
		t.Fatalf("expected missing repository %q to remain reported: %#v", missingRepository, report.Repositories)
	}
	if !storageArtifactRootGroup(report.ArtifactRoots, missingArtifactRoot).Missing {
		t.Fatalf("expected missing Artifact Root %q to remain reported: %#v", missingArtifactRoot, report.ArtifactRoots)
	}
	if len(report.States) < 2 {
		t.Fatalf("expected Active and Clean state groups, got %#v", report.States)
	}
	if len(report.Tables) == 0 {
		t.Fatal("expected table groups")
	}

	databaseAfter := storageTreeSnapshot(t, filepath.Join(homeDir, roundfixHomeDir))
	artifactAfter := storageTreeSnapshot(t, artifactRoot)
	if !reflect.DeepEqual(databaseAfter, databaseBefore) {
		t.Fatalf("storage report changed Run Database tree\nbefore: %#v\nafter:  %#v", databaseBefore, databaseAfter)
	}
	if !reflect.DeepEqual(artifactAfter, artifactBefore) {
		t.Fatalf("storage report changed artifact tree\nbefore: %#v\nafter:  %#v", artifactBefore, artifactAfter)
	}
}

func storageTreeSnapshot(t *testing.T, root string) map[string]string {
	t.Helper()
	snapshot := map[string]string{}
	err := filepath.WalkDir(root, func(path string, entry os.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		relative, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}
		if entry.IsDir() {
			snapshot[relative] = "directory"
			return nil
		}
		if entry.Type()&os.ModeSymlink != 0 {
			target, err := os.Readlink(path)
			if err != nil {
				return err
			}
			snapshot[relative] = "symlink:" + target
			return nil
		}
		content, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		digest := sha256.Sum256(content)
		snapshot[relative] = fmt.Sprintf("file:%d:%x", len(content), digest)
		return nil
	})
	if err != nil {
		t.Fatalf("snapshot storage tree %q: %v", root, err)
	}
	return snapshot
}

func storageRepositoryGroup(groups []StorageRepositoryGroup, repository string) StorageRepositoryGroup {
	for _, group := range groups {
		if group.Repository == repository {
			return group
		}
	}
	return StorageRepositoryGroup{}
}

func storageArtifactRootGroup(groups []StorageArtifactRootGroup, artifactRoot string) StorageArtifactRootGroup {
	for _, group := range groups {
		if group.ArtifactRoot == artifactRoot {
			return group
		}
	}
	return StorageArtifactRootGroup{}
}

func absInt64(value int64) int64 {
	if value < 0 {
		return -value
	}
	return value
}

func TestRunEventsAfterCursorReturnsOnlyNewerAndRespectsLimit(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	store := openTestStore(t, ctx, t.TempDir())
	defer closeStore(t, store)

	run, err := store.CreateRun(ctx, sampleCreateRunRequest())
	if err != nil {
		t.Fatalf("create run: %v", err)
	}
	for _, summary := range []string{"one", "two", "three", "four", "five"} {
		if _, err := store.AppendRunEvent(ctx, sampleRunEvent(run.ID, summary)); err != nil {
			t.Fatalf("append %q: %v", summary, err)
		}
	}

	page, err := store.RunEventsAfter(ctx, run.ID, 2, 2)
	if err != nil {
		t.Fatalf("list after cursor: %v", err)
	}
	if len(page) != 2 || page[0].Cursor != 3 || page[1].Cursor != 4 {
		t.Fatalf("expected cursors [3 4], got %+v", page)
	}
	if page[0].Event.Summary != "three" || page[1].Event.Summary != "four" {
		t.Fatalf("expected ordered summaries, got %+v", page)
	}

	rest, err := store.RunEventsAfter(ctx, run.ID, 4, 10)
	if err != nil {
		t.Fatalf("list tail: %v", err)
	}
	if len(rest) != 1 || rest[0].Cursor != 5 {
		t.Fatalf("expected final event, got %+v", rest)
	}

	if _, err := store.RunEventsAfter(ctx, run.ID, 0, 0); err == nil {
		t.Fatal("expected positive-limit requirement")
	}
}

func TestRunEventPayloadRoundTripsByteExact(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	store := openTestStore(t, ctx, t.TempDir())
	defer closeStore(t, store)

	run, err := store.CreateRun(ctx, sampleCreateRunRequest())
	if err != nil {
		t.Fatalf("create run: %v", err)
	}
	payload := []byte(`{"sessionId":"sess-1","update":{"sessionUpdate":"tool_call_update","toolCallId":"call_1","status":"completed","content":[{"type":"diff","path":"a.go","oldText":"x","newText":"y"}],"rawOutput":{"aggregated_output":"ok","unicode":"é✓"}}}`)
	event := sampleRunEvent(run.ID, "tool update")
	event.Kind = runevent.KindAgentToolUpdated
	event.ToolID = "call_1"
	event.ToolState = "completed"
	event.ReviewIssue = "issue_001"
	event.Payload = payload

	if _, err := store.AppendRunEvent(ctx, event); err != nil {
		t.Fatalf("append: %v", err)
	}
	stored, err := store.RunEventsAfter(ctx, run.ID, 0, 1)
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if len(stored) != 1 {
		t.Fatalf("expected one event, got %d", len(stored))
	}
	got := stored[0].Event
	if !bytes.Equal(got.Payload, payload) {
		t.Fatalf("expected byte-exact payload round-trip\nwant: %s\ngot:  %s", payload, got.Payload)
	}
	if got.ToolID != "call_1" || got.ToolState != "completed" || got.ReviewIssue != "issue_001" {
		t.Fatalf("expected normalized columns preserved, got %+v", got)
	}
	if got.Kind != runevent.KindAgentToolUpdated || got.Source != runevent.SourceAgent || got.Batch != 1 {
		t.Fatalf("expected kind/source/batch preserved, got %+v", got)
	}
	if !got.Time.Equal(event.Time) {
		t.Fatalf("expected timestamp preserved, got %v", got.Time)
	}
}

func TestReaderPagesEventsWhileWriterAppends(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	homeDir := t.TempDir()
	writer := openTestStore(t, ctx, homeDir)
	defer closeStore(t, writer)

	run, err := writer.CreateRun(ctx, sampleCreateRunRequest())
	if err != nil {
		t.Fatalf("create run: %v", err)
	}
	reader, err := OpenReader(ctx, homeDir)
	if err != nil {
		t.Fatalf("open reader: %v", err)
	}
	defer closeStore(t, reader)

	const total = 40
	appended := make(chan error, 1)
	go func() {
		for index := 0; index < total; index++ {
			if _, err := writer.AppendRunEvent(ctx, sampleRunEvent(run.ID, "concurrent")); err != nil {
				appended <- err
				return
			}
		}
		appended <- nil
	}()

	var cursor int64
	seen := 0
	for seen < total {
		page, err := reader.RunEventsAfter(ctx, run.ID, cursor, 7)
		if err != nil {
			t.Fatalf("reader page: %v", err)
		}
		for _, entry := range page {
			if entry.Cursor != cursor+1 {
				t.Fatalf("expected contiguous cursors, got %d after %d", entry.Cursor, cursor)
			}
			cursor = entry.Cursor
			seen++
		}
		select {
		case err := <-appended:
			if err != nil {
				t.Fatalf("writer append: %v", err)
			}
			appended = nil
		default:
		}
	}
	if seen != total {
		t.Fatalf("expected %d events paged, got %d", total, seen)
	}

	if _, err := reader.AppendRunEvent(ctx, sampleRunEvent(run.ID, "rejected")); err == nil {
		t.Fatal("expected read-only connection to reject writes")
	}
}

func TestDataVersionSignalsWriterCommitsToPollers(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	homeDir := t.TempDir()
	writer := openTestStore(t, ctx, homeDir)
	defer closeStore(t, writer)
	run, err := writer.CreateRun(ctx, sampleCreateRunRequest())
	if err != nil {
		t.Fatalf("create run: %v", err)
	}

	reader, err := OpenReader(ctx, homeDir)
	if err != nil {
		t.Fatalf("open reader: %v", err)
	}
	defer closeStore(t, reader)

	before, err := reader.DataVersion(ctx)
	if err != nil {
		t.Fatalf("read data version: %v", err)
	}
	if _, err := writer.AppendRunEvent(ctx, sampleRunEvent(run.ID, "signal")); err != nil {
		t.Fatalf("append: %v", err)
	}
	after, err := reader.DataVersion(ctx)
	if err != nil {
		t.Fatalf("read data version after append: %v", err)
	}
	if after == before {
		t.Fatalf("expected data version change after writer commit, got %d twice", after)
	}
}

func TestRunEventsBeforeReturnsImmediatelyPrecedingEventsAscending(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	store := openTestStore(t, ctx, t.TempDir())
	defer closeStore(t, store)

	run, err := store.CreateRun(ctx, sampleCreateRunRequest())
	if err != nil {
		t.Fatalf("create run: %v", err)
	}
	for _, summary := range []string{"one", "two", "three", "four", "five"} {
		if _, err := store.AppendRunEvent(ctx, sampleRunEvent(run.ID, summary)); err != nil {
			t.Fatalf("append %q: %v", summary, err)
		}
	}

	page, err := store.RunEventsBefore(ctx, run.ID, 4, 2)
	if err != nil {
		t.Fatalf("list before cursor: %v", err)
	}
	if len(page) != 2 || page[0].Cursor != 2 || page[1].Cursor != 3 {
		t.Fatalf("expected the two events immediately before cursor 4, ascending, got %+v", page)
	}
	if page[0].Event.Summary != "two" || page[1].Event.Summary != "three" {
		t.Fatalf("expected ordered summaries, got %+v", page)
	}

	if _, err := store.RunEventsBefore(ctx, run.ID, 4, 0); err == nil {
		t.Fatal("expected positive-limit requirement")
	}
	if _, err := store.RunEventsBefore(ctx, "", 4, 1); err == nil {
		t.Fatal("expected Run ID requirement")
	}
}

func TestRunEventsBeforeBoundaryCursors(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	store := openTestStore(t, ctx, t.TempDir())
	defer closeStore(t, store)

	run, err := store.CreateRun(ctx, sampleCreateRunRequest())
	if err != nil {
		t.Fatalf("create run: %v", err)
	}
	for _, summary := range []string{"one", "two", "three"} {
		if _, err := store.AppendRunEvent(ctx, sampleRunEvent(run.ID, summary)); err != nil {
			t.Fatalf("append %q: %v", summary, err)
		}
	}

	for _, cursor := range []int64{0, 1} {
		page, err := store.RunEventsBefore(ctx, run.ID, cursor, 10)
		if err != nil {
			t.Fatalf("list before %d: %v", cursor, err)
		}
		if len(page) != 0 {
			t.Fatalf("expected nothing before cursor %d, got %+v", cursor, page)
		}
	}
	page, err := store.RunEventsBefore(ctx, run.ID, 2, 10)
	if err != nil {
		t.Fatalf("list before 2: %v", err)
	}
	if len(page) != 1 || page[0].Cursor != 1 {
		t.Fatalf("expected exactly the first event, got %+v", page)
	}
	missing, err := store.RunEventsBefore(ctx, "run_missing", 10, 5)
	if err != nil {
		t.Fatalf("expected empty page for unknown Run like the forward read, got %v", err)
	}
	if len(missing) != 0 {
		t.Fatalf("expected empty page, got %+v", missing)
	}
}

func TestRunEventsBeforePagesComposeTailToHeadWithoutDuplicates(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	store := openTestStore(t, ctx, t.TempDir())
	defer closeStore(t, store)

	run, err := store.CreateRun(ctx, sampleCreateRunRequest())
	if err != nil {
		t.Fatalf("create run: %v", err)
	}
	const total = 23
	for index := 0; index < total; index++ {
		if _, err := store.AppendRunEvent(ctx, sampleRunEvent(run.ID, "event")); err != nil {
			t.Fatalf("append %d: %v", index, err)
		}
	}

	// Walk tail→head with a page size that does not divide the total, so
	// the head page is partial — an exact page-edge case.
	cursor := int64(total + 1)
	seen := []int64{}
	for {
		page, err := store.RunEventsBefore(ctx, run.ID, cursor, 5)
		if err != nil {
			t.Fatalf("page before %d: %v", cursor, err)
		}
		if len(page) == 0 {
			break
		}
		for index := len(page) - 1; index >= 0; index-- {
			seen = append(seen, page[index].Cursor)
		}
		cursor = page[0].Cursor
	}
	if len(seen) != total {
		t.Fatalf("expected every event exactly once, got %d of %d", len(seen), total)
	}
	for index, cursor := range seen {
		if cursor != int64(total-index) {
			t.Fatalf("expected strictly descending walk without gaps, got %v", seen)
		}
	}
}

func TestRunEventsBeforeOnReaderWhileWriterAppends(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	homeDir := t.TempDir()
	writer := openTestStore(t, ctx, homeDir)
	defer closeStore(t, writer)
	run, err := writer.CreateRun(ctx, sampleCreateRunRequest())
	if err != nil {
		t.Fatalf("create run: %v", err)
	}
	for index := 0; index < 10; index++ {
		if _, err := writer.AppendRunEvent(ctx, sampleRunEvent(run.ID, "seed")); err != nil {
			t.Fatalf("seed append: %v", err)
		}
	}
	reader, err := OpenReader(ctx, homeDir)
	if err != nil {
		t.Fatalf("open reader: %v", err)
	}
	defer closeStore(t, reader)

	appended := make(chan error, 1)
	go func() {
		for index := 0; index < 20; index++ {
			if _, err := writer.AppendRunEvent(ctx, sampleRunEvent(run.ID, "live")); err != nil {
				appended <- err
				return
			}
		}
		appended <- nil
	}()

	// Backward pages over the stable prefix stay correct while the writer
	// appends past the tail: short autocommit reads, no read transaction.
	for round := 0; round < 10; round++ {
		page, err := reader.RunEventsBefore(ctx, run.ID, 6, 5)
		if err != nil {
			t.Fatalf("reader page: %v", err)
		}
		if len(page) != 5 || page[0].Cursor != 1 || page[4].Cursor != 5 {
			t.Fatalf("expected stable backward page [1..5], got %+v", page)
		}
	}
	if err := <-appended; err != nil {
		t.Fatalf("writer append: %v", err)
	}
}

func countRunEvents(t *testing.T, ctx context.Context, store *Store, runID string) int {
	t.Helper()
	var count int
	if err := store.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM run_events WHERE run_id = ?`, runID).Scan(&count); err != nil {
		t.Fatalf("count Run Events for %s: %v", runID, err)
	}
	return count
}

func countActiveRunLocks(t *testing.T, ctx context.Context, store *Store) int {
	t.Helper()
	var count int
	if err := store.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM active_run_locks`).Scan(&count); err != nil {
		t.Fatalf("count Active Run locks: %v", err)
	}
	return count
}

func TestCompactionPreviewMatchesResultWithinDeclaredTolerance(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	homeDir := t.TempDir()
	runStore := openTestStore(t, ctx, homeDir)
	defer closeStore(t, runStore)
	seedPrunableRunEvents(t, ctx, runStore)

	preview, err := runStore.PreviewCompaction(ctx)
	if err != nil {
		t.Fatalf("preview Run Database compaction: %v", err)
	}
	if preview.BytesBefore <= 0 {
		t.Fatalf("expected measured bytes before compaction, got %#v", preview)
	}
	if preview.BytesReclaimable <= 0 {
		t.Fatalf("expected measured reclaimable bytes, got %#v", preview)
	}
	if preview.BytesAfter != preview.BytesBefore-preview.BytesReclaimable {
		t.Fatalf("preview relation does not reconcile: %#v", preview)
	}
	if preview.ReconciliationToleranceBytes <= 0 || preview.ReconciliationToleranceReason == "" {
		t.Fatalf("expected declared compaction tolerance, got %#v", preview)
	}

	result, err := runStore.Compact(ctx, preview)
	if err != nil {
		t.Fatalf("compact Run Database: %v", err)
	}
	reclaimDifference := absInt64(result.BytesReclaimed - preview.BytesReclaimable)
	if reclaimDifference > preview.ReconciliationToleranceBytes {
		t.Fatalf(
			"compaction reclaimed %d bytes, preview projected %d, difference %d exceeds tolerance %d",
			result.BytesReclaimed,
			preview.BytesReclaimable,
			reclaimDifference,
			preview.ReconciliationToleranceBytes,
		)
	}
	if result.BytesBefore != preview.BytesBefore || result.BytesAfter != result.BytesBefore-result.BytesReclaimed {
		t.Fatalf("compaction result does not reconcile with preview: preview=%#v result=%#v", preview, result)
	}
	databaseInfo, err := os.Stat(DatabasePath(homeDir))
	if err != nil {
		t.Fatalf("stat compacted Run Database: %v", err)
	}
	if result.BytesAfter != databaseInfo.Size() {
		t.Fatalf("compaction reported %d bytes after, file has %d", result.BytesAfter, databaseInfo.Size())
	}
	var integrity string
	if err := runStore.db.QueryRowContext(ctx, `PRAGMA quick_check`).Scan(&integrity); err != nil {
		t.Fatalf("check compacted Run Database integrity: %v", err)
	}
	if integrity != "ok" {
		t.Fatalf("compacted Run Database integrity check returned %q", integrity)
	}
}

func TestCompactRefusesActiveRunAndPreservesDatabaseBytes(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	homeDir := t.TempDir()
	runStore := openTestStore(t, ctx, homeDir)
	defer closeStore(t, runStore)

	active, err := runStore.CreateRun(ctx, sampleCreateRunRequest())
	if err != nil {
		t.Fatalf("create Active Run: %v", err)
	}
	preview, err := runStore.PreviewCompaction(ctx)
	if err != nil {
		t.Fatalf("preview Run Database compaction: %v", err)
	}
	before := databaseFileSize(t, homeDir)

	_, err = runStore.Compact(ctx, preview)

	var activeErr ActiveRunCompactionError
	if !errors.As(err, &activeErr) {
		t.Fatalf("expected Active Run compaction refusal, got %T %v", err, err)
	}
	if activeErr.RunID != active.ID || !strings.Contains(err.Error(), active.ID) {
		t.Fatalf("expected refusal to name Active Run %q, got %#v: %v", active.ID, activeErr, err)
	}
	after := databaseFileSize(t, homeDir)
	if after != before {
		t.Fatalf("Active Run refusal changed Run Database bytes: before=%d after=%d", before, after)
	}
}

func TestCompactRefusesAnotherWriterAndPreservesDatabaseBytes(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	homeDir := t.TempDir()
	runStore := openTestStore(t, ctx, homeDir)
	defer closeStore(t, runStore)
	seedPrunableRunEvents(t, ctx, runStore)
	preview, err := runStore.PreviewCompaction(ctx)
	if err != nil {
		t.Fatalf("preview Run Database compaction: %v", err)
	}
	otherWriter := openTestStore(t, ctx, homeDir)
	defer closeStore(t, otherWriter)
	before := databaseFileSize(t, homeDir)

	_, err = runStore.Compact(ctx, preview)

	var writerErr WriterPresentCompactionError
	if !errors.As(err, &writerErr) {
		t.Fatalf("expected writer-present compaction refusal, got %T %v", err, err)
	}
	if !strings.Contains(err.Error(), "another Run Database connection") {
		t.Fatalf("expected refusal to name the blocking writer, got %v", err)
	}
	after := databaseFileSize(t, homeDir)
	if after != before {
		t.Fatalf("writer refusal changed Run Database bytes: before=%d after=%d", before, after)
	}
}

func TestCompactRefusesInsufficientTemporaryCapacityBeforeMutation(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	homeDir := t.TempDir()
	runStore := openTestStore(t, ctx, homeDir)
	defer closeStore(t, runStore)
	seedPrunableRunEvents(t, ctx, runStore)
	preview, err := runStore.PreviewCompaction(ctx)
	if err != nil {
		t.Fatalf("preview Run Database compaction: %v", err)
	}
	required := preview.BytesBefore * 2
	runStore.temporaryCapacity = func(string) (int64, error) {
		return required - 1, nil
	}
	before := databaseFileSize(t, homeDir)

	_, err = runStore.Compact(ctx, preview)

	var capacityErr CompactionCapacityError
	if !errors.As(err, &capacityErr) {
		t.Fatalf("expected temporary-capacity refusal, got %T %v", err, err)
	}
	if capacityErr.RequiredBytes != required || capacityErr.AvailableBytes != required-1 || capacityErr.ShortfallBytes != 1 {
		t.Fatalf("capacity refusal did not name the measured shortfall: %#v", capacityErr)
	}
	if !strings.Contains(err.Error(), "shortfall=1") {
		t.Fatalf("expected refusal to name one-byte shortfall, got %v", err)
	}
	after := databaseFileSize(t, homeDir)
	if after != before {
		t.Fatalf("capacity refusal changed Run Database bytes: before=%d after=%d", before, after)
	}
}

func TestCompactRefusesStalePreviewAndPreservesDatabaseBytes(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	homeDir := t.TempDir()
	runStore := openTestStore(t, ctx, homeDir)
	defer closeStore(t, runStore)
	seedPrunableRunEvents(t, ctx, runStore)
	preview, err := runStore.PreviewCompaction(ctx)
	if err != nil {
		t.Fatalf("preview Run Database compaction: %v", err)
	}
	request := sampleCreateRunRequest()
	request.HeadBranch = "feature/after-preview"
	request.PRNumber = "after-preview"
	run, err := runStore.CreateRun(ctx, request)
	if err != nil {
		t.Fatalf("create Run after preview: %v", err)
	}
	if _, err := runStore.CompleteRun(ctx, run.ID, StateClean); err != nil {
		t.Fatalf("complete Run after preview: %v", err)
	}
	before := databaseFileSize(t, homeDir)

	_, err = runStore.Compact(ctx, preview)

	var staleErr CompactionPreviewStaleError
	if !errors.As(err, &staleErr) {
		t.Fatalf("expected stale-preview refusal, got %T %v", err, err)
	}
	after := databaseFileSize(t, homeDir)
	if after != before {
		t.Fatalf("stale-preview refusal changed Run Database bytes: before=%d after=%d", before, after)
	}
}

func TestPruneTerminalRunsNeverVacuumsAllocatedPages(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	homeDir := t.TempDir()
	runStore := openTestStore(t, ctx, homeDir)
	defer closeStore(t, runStore)

	run := seedTerminalRunEvents(t, ctx, runStore)
	checkpointRunDatabase(t, ctx, runStore)
	before := databaseFileSize(t, homeDir)
	result, err := runStore.PruneTerminalRuns(ctx, time.Now().UTC().Add(time.Hour))
	if err != nil {
		t.Fatalf("prune Run Event Journal: %v", err)
	}
	if !slices.Contains(result.RunIDs, run.ID) || result.Events == 0 {
		t.Fatalf("expected fixture Run Events pruned, got %#v", result)
	}
	checkpointRunDatabase(t, ctx, runStore)
	after := databaseFileSize(t, homeDir)
	if after != before {
		t.Fatalf("retention sweep compacted Run Database as a side effect: before=%d after=%d", before, after)
	}
	freePages, err := storagePragmaInt64(ctx, runStore.db, "freelist_count")
	if err != nil {
		t.Fatalf("measure free pages after retention sweep: %v", err)
	}
	if freePages == 0 {
		t.Fatal("expected retention sweep to leave deleted pages allocated for explicit compaction")
	}
}

func seedPrunableRunEvents(t *testing.T, ctx context.Context, runStore *Store) Run {
	t.Helper()
	run := seedTerminalRunEvents(t, ctx, runStore)
	if _, err := runStore.PruneTerminalRuns(ctx, time.Now().UTC().Add(time.Hour)); err != nil {
		t.Fatalf("prune Run Event Journal fixture: %v", err)
	}
	return run
}

func seedTerminalRunEvents(t *testing.T, ctx context.Context, runStore *Store) Run {
	t.Helper()
	run, err := runStore.CreateRun(ctx, sampleCreateRunRequest())
	if err != nil {
		t.Fatalf("create compaction fixture Run: %v", err)
	}
	for index := 0; index < 48; index++ {
		event := sampleRunEvent(run.ID, fmt.Sprintf("compaction fixture %d", index))
		event.Payload = bytes.Repeat([]byte{byte(index + 1)}, 16*1024)
		if _, err := runStore.AppendRunEvent(ctx, event); err != nil {
			t.Fatalf("append compaction fixture Run Event %d: %v", index, err)
		}
	}
	completed, err := runStore.CompleteRun(ctx, run.ID, StateClean)
	if err != nil {
		t.Fatalf("complete compaction fixture Run: %v", err)
	}
	return completed.Run
}

func checkpointRunDatabase(t *testing.T, ctx context.Context, runStore *Store) {
	t.Helper()
	var busy int
	var logFrames int
	var checkpointedFrames int
	if err := runStore.db.QueryRowContext(ctx, `PRAGMA wal_checkpoint(TRUNCATE)`).Scan(&busy, &logFrames, &checkpointedFrames); err != nil {
		t.Fatalf("checkpoint Run Database: %v", err)
	}
	if busy != 0 {
		t.Fatalf("checkpoint Run Database remained busy: log=%d checkpointed=%d", logFrames, checkpointedFrames)
	}
}

func databaseFileSize(t *testing.T, homeDir string) int64 {
	t.Helper()
	info, err := os.Stat(DatabasePath(homeDir))
	if err != nil {
		t.Fatalf("stat Run Database: %v", err)
	}
	return info.Size()
}

func TestDiscoverArtifactRootsReturnsEveryRecordedRepository(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	runStore := openTestStore(t, ctx, t.TempDir())
	defer closeStore(t, runStore)

	firstRequest := sampleCreateRunRequest()
	firstRequest.GitRoot = filepath.Join(t.TempDir(), "repository-a")
	firstRequest.ArtifactDir = filepath.Join(t.TempDir(), "artifacts-a")
	first, err := runStore.CreateRun(ctx, firstRequest)
	if err != nil {
		t.Fatalf("create first Artifact Root Run: %v", err)
	}
	if _, err := runStore.CompleteRun(ctx, first.ID, StateClean); err != nil {
		t.Fatalf("complete first Artifact Root Run: %v", err)
	}

	secondRequest := sampleCreateRunRequest()
	secondRequest.HeadBranch = "feature/other-artifact-root"
	secondRequest.PRNumber = "456"
	secondRequest.GitRoot = filepath.Join(t.TempDir(), "repository-b")
	secondRequest.ArtifactDir = filepath.Join(t.TempDir(), "artifacts-b")
	second, err := runStore.CreateRun(ctx, secondRequest)
	if err != nil {
		t.Fatalf("create second Artifact Root Run: %v", err)
	}

	roots, err := DiscoverArtifactRoots(ctx, runStore)
	if err != nil {
		t.Fatalf("discover Artifact Roots: %v", err)
	}
	if len(roots) != 2 {
		t.Fatalf("expected roots from both repositories, got %#v", roots)
	}
	runsByRoot := map[string]ArtifactRootRun{}
	for _, root := range roots {
		if len(root.Runs) != 1 {
			t.Fatalf("expected one Run for Artifact Root %q, got %#v", root.Path, root.Runs)
		}
		runsByRoot[root.Path] = root.Runs[0]
	}
	if got := runsByRoot[firstRequest.ArtifactDir]; got.ID != first.ID || got.Repository != firstRequest.GitRoot || got.State != StateClean || got.CompletedAt == nil {
		t.Fatalf("first Artifact Root lost durable Run evidence: %#v", got)
	}
	if got := runsByRoot[secondRequest.ArtifactDir]; got.ID != second.ID || got.Repository != secondRequest.GitRoot || got.State != StateActive || got.CompletedAt != nil {
		t.Fatalf("second Artifact Root lost durable Run evidence: %#v", got)
	}
}
