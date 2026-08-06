package store

import (
	"bytes"
	"context"
	"crypto/sha256"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"slices"
	"strings"
	"testing"
	"time"

	"roundfix/internal/runevent"
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
