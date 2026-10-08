// Suite: whole-Run retention and incremental Run Database compaction.
// Boundary: Store APIs against real SQLite files in temporary homes.
package store

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"testing"
	"time"
)

var retentionFixtureCutoff = time.Date(2026, 10, 8, 12, 0, 0, 123456789, time.UTC)
var retentionDependentTables = []string{"runs", "run_events", "run_agent_selections", "run_token_usage", "active_run_locks"}
var retentionUnchangedTables = []string{"run_windows", "interactive_defaults", "delivery_queues", "delivery_queue_items", "delivery_queue_runs", "delivery_action_intents", "delivery_action_receipts"}

func seedRetainedRun(t *testing.T, s *Store, state string, completed time.Time) (Run, int64) {
	t.Helper()
	ctx := context.Background()
	s.now = func() time.Time { return completed }
	req := sampleCreateRunRequest()
	req.HeadBranch = t.Name() + fmt.Sprint(completed.UnixNano())
	run, err := s.CreateRunSkippingActiveLock(ctx, req)
	if err != nil {
		t.Fatal(err)
	}
	event := sampleRunEvent(run.ID, "retenção 🐾")
	event.Payload = []byte(`{"message":"ação"}`)
	if _, err := s.AppendRunEvent(ctx, event); err != nil {
		t.Fatal(err)
	}
	if _, err := s.AppendAgentSelectionAttempt(ctx, AgentSelectionAttemptRequest{
		RunID: run.ID, ScopeKind: AgentSelectionScopeTask, ScopeID: "task_01", Category: "backend", ProfileSource: "project",
		Attempt: 1, SelectionRole: AgentSelectionRolePreferred, Runtime: "codex", Model: "fixture", Status: AgentSelectionStatusClosed,
	}); err != nil {
		t.Fatal(err)
	}
	if err := s.AppendTokenUsage(ctx, TokenUsageRecord{RunID: run.ID, ScopeKind: "task", ScopeID: "task_01", Session: "fixture", Basis: "unreported"}); err != nil {
		t.Fatal(err)
	}
	if IsTerminalState(state) {
		if _, err := s.CompleteRun(ctx, run.ID, state); err != nil {
			t.Fatal(err)
		}
	}
	// Retention must also cascade a leftover lock on a terminal Run.
	retentionExec(t, s.db, `INSERT INTO active_run_locks (target_kind,target_key,run_id,created_at) VALUES ('pr',?,?,?)`, run.ID, run.ID, formatTime(completed))
	events, err := s.RunEventsAfter(ctx, run.ID, 0, 100)
	if err != nil {
		t.Fatal(err)
	}
	var estimate int64
	for _, entry := range events {
		estimate += int64(len(entry.Event.Payload) + len(entry.Event.Summary))
	}
	return run, estimate
}

func retentionExec(t *testing.T, db *sql.DB, query string, args ...any) {
	t.Helper()
	if _, err := db.ExecContext(context.Background(), query, args...); err != nil {
		t.Fatal(err)
	}
}

func seedRetentionQueue(t *testing.T, s *Store, runID, reference string) {
	t.Helper()
	ctx := context.Background()
	root := filepath.Join(t.TempDir(), "queue")
	if _, err := s.CreateDeliveryQueue(ctx, root, []string{"fixture"}); err != nil {
		t.Fatal(err)
	}
	if reference == "item" {
		retentionExec(t, s.db, `UPDATE delivery_queue_items SET run_id = ? WHERE git_root = ?`, runID, root)
	}
	if reference == "link" {
		retentionExec(t, s.db, `INSERT INTO delivery_queue_runs (git_root,spec_slug,run_id) VALUES (?, 'fixture', ?)`, root, runID)
	}
	retentionExec(t, s.db, `INSERT INTO delivery_action_intents (git_root,spec_slug,action,created_at) VALUES (?, 'fixture','push','fixture')`, root)
	retentionExec(t, s.db, `INSERT INTO delivery_action_receipts (intent_id,result,created_at) VALUES (last_insert_rowid(),'receipt','fixture')`)
	retentionExec(t, s.db, `INSERT OR REPLACE INTO run_windows VALUES (?, 123, 456)`, root)
	retentionExec(t, s.db, `INSERT OR REPLACE INTO interactive_defaults VALUES (?, 'value', 'fixture')`, root)
}

// Snapshot every column, sorting rows only to avoid query ordering assumptions.
func retentionSnapshot(t *testing.T, db *sql.DB, tables []string, excludeRun string) map[string][]string {
	t.Helper()
	result := map[string][]string{}
	for _, table := range tables {
		query := `SELECT * FROM ` + table
		var args []any
		if excludeRun != "" {
			key := "run_id"
			if table == "runs" {
				key = "id"
			}
			query += ` WHERE ` + key + ` <> ?`
			args = append(args, excludeRun)
		}
		rows, err := db.QueryContext(context.Background(), query, args...)
		if err != nil {
			t.Fatal(err)
		}
		cols, err := rows.Columns()
		if err != nil {
			t.Fatal(err)
		}
		result[table] = []string{}
		for rows.Next() {
			values := make([]any, len(cols))
			pointers := make([]any, len(cols))
			for i := range values {
				pointers[i] = &values[i]
			}
			if err := rows.Scan(pointers...); err != nil {
				t.Fatal(err)
			}
			raw, err := json.Marshal(values)
			if err != nil {
				t.Fatal(err)
			}
			result[table] = append(result[table], string(raw))
		}
		if err := rows.Err(); err != nil {
			t.Fatal(err)
		}
		if err := rows.Close(); err != nil {
			t.Fatal(err)
		}
		sort.Strings(result[table])
	}
	return result
}

func TestRunRetentionCandidatesCountEveryDependentRow(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	home := t.TempDir()
	s := openTestStore(t, ctx, home)
	defer closeStore(t, s)
	old, estimate := seedRetainedRun(t, s, StateClean, retentionFixtureCutoff.Add(-time.Hour))
	seedRetainedRun(t, s, StateActive, retentionFixtureCutoff.Add(-2*time.Hour))
	seedRetainedRun(t, s, StateClean, retentionFixtureCutoff.Add(time.Hour))
	queued, queueEstimate := seedRetainedRun(t, s, StateFailed, retentionFixtureCutoff.Add(-3*time.Hour))
	seedRetentionQueue(t, s, queued.ID, "item")
	reader, err := OpenReader(ctx, home)
	if err != nil {
		t.Fatal(err)
	}
	defer closeStore(t, reader)
	candidates, err := reader.RunRetentionCandidates(ctx, retentionFixtureCutoff)
	if err != nil {
		t.Fatal(err)
	}
	if len(candidates) != 2 {
		t.Fatalf("candidates = %+v", candidates)
	}
	want := map[string]int64{old.ID: estimate, queued.ID: queueEstimate}
	for _, c := range candidates {
		// The fixture appends one event, one selection event and one usage event.
		expected := RunRetentionRows{Runs: 1, RunEvents: 3, AgentSelections: 1, TokenUsage: 1, ActiveRunLocks: 1}
		bytes, ok := want[c.RunID]
		if !ok || c.Rows != expected || c.EstimatedBytes != bytes || c.QueueReferenced != (c.RunID == queued.ID) {
			t.Fatalf("candidate=%+v expected rows=%+v bytes=%d", c, expected, bytes)
		}
		if c.Repository == "" || c.GitRoot == "" || c.ArtifactDir == "" || !c.CompletedAt.Before(retentionFixtureCutoff) {
			t.Fatalf("metadata = %+v", c)
		}
	}
	// A historical retry link alone also protects the Run.
	retentionExec(t, s.db, `UPDATE delivery_queue_items SET run_id = ''`)
	seedRetentionQueue(t, s, queued.ID, "link")
	candidates, err = reader.RunRetentionCandidates(ctx, retentionFixtureCutoff)
	if err != nil {
		t.Fatal(err)
	}
	for _, c := range candidates {
		if c.RunID == queued.ID && !c.QueueReferenced {
			t.Fatal("retry link ignored")
		}
	}
}

func TestRemoveRetainedRunDeletesTheRunWhole(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	s := openTestStore(t, ctx, t.TempDir())
	defer closeStore(t, s)
	old, _ := seedRetainedRun(t, s, StateClean, retentionFixtureCutoff.Add(-time.Hour))
	other, _ := seedRetainedRun(t, s, StateActive, retentionFixtureCutoff)
	seedRetentionQueue(t, s, other.ID, "link")
	beforeOther := retentionSnapshot(t, s.db, retentionDependentTables, old.ID)
	beforeUntouched := retentionSnapshot(t, s.db, retentionUnchangedTables, "")
	removed, err := s.RemoveRetainedRun(ctx, old.ID, retentionFixtureCutoff)
	if err != nil {
		t.Fatal(err)
	}
	want := RunRetentionRows{1, 3, 1, 1, 1}
	if removed != want {
		t.Fatalf("removed=%+v want=%+v", removed, want)
	}
	for _, table := range retentionDependentTables {
		key := "run_id"
		if table == "runs" {
			key = "id"
		}
		var count int
		if err := s.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM `+table+` WHERE `+key+` = ?`, old.ID).Scan(&count); err != nil {
			t.Fatal(err)
		}
		if count != 0 {
			t.Fatalf("%s kept %d rows", table, count)
		}
	}
	if after := retentionSnapshot(t, s.db, retentionDependentTables, old.ID); !reflect.DeepEqual(after, beforeOther) {
		t.Fatal("other Run rows changed")
	}
	if after := retentionSnapshot(t, s.db, retentionUnchangedTables, ""); !reflect.DeepEqual(after, beforeUntouched) {
		t.Fatal("unrelated tables changed")
	}
	if got, err := s.RemoveRetainedRun(ctx, old.ID, retentionFixtureCutoff); err != nil || got != (RunRetentionRows{}) {
		t.Fatalf("repeat removal = %+v, %v", got, err)
	}
}

func TestRemoveRetainedRunRefusesARunThatChanged(t *testing.T) {
	t.Parallel()
	for _, change := range []string{"item", "link", "recent", "active", "boundary", "nanosecond before"} {
		t.Run(change, func(t *testing.T) {
			t.Parallel()
			ctx := context.Background()
			s := openTestStore(t, ctx, t.TempDir())
			defer closeStore(t, s)
			run, _ := seedRetainedRun(t, s, StateClean, retentionFixtureCutoff.Add(-time.Hour))
			if candidates, err := s.RunRetentionCandidates(ctx, retentionFixtureCutoff); err != nil || len(candidates) != 1 {
				t.Fatalf("scan=%+v %v", candidates, err)
			}
			switch change {
			case "item", "link":
				seedRetentionQueue(t, s, run.ID, change)
			case "active":
				retentionExec(t, s.db, `UPDATE runs SET state = ? WHERE id = ?`, StateActive, run.ID)
			case "recent":
				retentionExec(t, s.db, `UPDATE runs SET completed_at = ? WHERE id = ?`, formatTime(retentionFixtureCutoff.Add(time.Hour)), run.ID)
			case "boundary":
				retentionExec(t, s.db, `UPDATE runs SET completed_at = ? WHERE id = ?`, formatTime(retentionFixtureCutoff), run.ID)
			case "nanosecond before":
				retentionExec(t, s.db, `UPDATE runs SET completed_at = ? WHERE id = ?`, formatTime(retentionFixtureCutoff.Add(-time.Nanosecond)), run.ID)
			}
			if change == "boundary" || change == "nanosecond before" {
				candidates, err := s.RunRetentionCandidates(ctx, retentionFixtureCutoff)
				if err != nil {
					t.Fatal(err)
				}
				want := 0
				if change == "nanosecond before" {
					want = 1
				}
				if len(candidates) != want {
					t.Fatalf("strict cutoff candidates = %+v", candidates)
				}
			}
			before := retentionSnapshot(t, s.db, retentionDependentTables, "")
			got, err := s.RemoveRetainedRun(ctx, run.ID, retentionFixtureCutoff)
			if change == "nanosecond before" {
				if err != nil || got.Runs != 1 {
					t.Fatalf("strict before removal=%+v %v", got, err)
				}
				return
			}
			var kept RunRetentionKeptError
			if !errors.As(err, &kept) || kept.RunID != run.ID || kept.Reason == "" || got != (RunRetentionRows{}) {
				t.Fatalf("removal=%+v error=%v", got, err)
			}
			if after := retentionSnapshot(t, s.db, retentionDependentTables, ""); !reflect.DeepEqual(before, after) {
				t.Fatal("refusal changed rows")
			}
		})
	}
}

func TestRunRetentionSweepRecordRoundTrips(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	home := t.TempDir()
	s := openTestStore(t, ctx, home)
	defer closeStore(t, s)
	reader, err := OpenReader(ctx, home)
	if err != nil {
		t.Fatal(err)
	}
	defer closeStore(t, reader)
	if got, found, err := reader.LastRunRetentionSweep(ctx); err != nil || found || got != (RunRetentionSweep{}) {
		t.Fatalf("initial=%+v %v %v", got, found, err)
	}
	for _, days := range []int{30, 7} {
		want := RunRetentionSweep{retentionFixtureCutoff.Add(time.Duration(days) * time.Hour), days}
		if err := s.RecordRunRetentionSweep(ctx, want); err != nil {
			t.Fatal(err)
		}
		got, found, err := reader.LastRunRetentionSweep(ctx)
		if err != nil || !found || got.RetentionDays != want.RetentionDays || !got.CompletedAt.Equal(want.CompletedAt) {
			t.Fatalf("record=%+v %v %v", got, found, err)
		}
	}
	var count int
	if err := s.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM run_retention_sweeps`).Scan(&count); err != nil || count != 1 {
		t.Fatalf("count=%d %v", count, err)
	}
}

func retentionPragma(t *testing.T, s *Store, name string) int64 {
	t.Helper()
	value, err := storagePragmaInt64(context.Background(), s.db, name)
	if err != nil {
		t.Fatal(err)
	}
	return value
}
func TestNewRunDatabaseUsesIncrementalAutoVacuum(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	s := openTestStore(t, ctx, t.TempDir())
	defer closeStore(t, s)
	if mode := retentionPragma(t, s, "auto_vacuum"); mode != 2 {
		t.Fatalf("mode=%d", mode)
	}
}

func openDefaultRetentionFixture(t *testing.T, home string) *Store {
	t.Helper()
	path := DatabasePath(home)
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		t.Fatal(err)
	}
	db, err := sql.Open("sqlite", migratorDSN(path))
	if err != nil {
		t.Fatal(err)
	}
	retentionExec(t, db, `CREATE TABLE default_mode_fixture (id INTEGER)`)
	retentionExec(t, db, `DROP TABLE default_mode_fixture`)
	if err := db.Close(); err != nil {
		t.Fatal(err)
	}
	s := openTestStore(t, context.Background(), home)
	if mode := retentionPragma(t, s, "auto_vacuum"); mode != 0 {
		t.Fatalf("default fixture mode=%d", mode)
	}
	return s
}

func TestCompactIncrementallyShrinksTheFile(t *testing.T) {
	t.Parallel()
	for _, mode := range []string{"incremental", "default"} {
		t.Run(mode, func(t *testing.T) {
			t.Parallel()
			ctx := context.Background()
			home := t.TempDir()
			var s *Store
			if mode == "default" {
				s = openDefaultRetentionFixture(t, home)
			} else {
				s = openTestStore(t, ctx, home)
			}
			defer closeStore(t, s)
			run := seedTerminalRunEvents(t, ctx, s)
			checkpointRunDatabase(t, ctx, s)
			pagesBefore := retentionPragma(t, s, "page_count")
			bytesBefore := databaseFileSize(t, home)
			if _, err := s.RemoveRetainedRun(ctx, run.ID, time.Now().Add(time.Hour)); err != nil {
				t.Fatal(err)
			}
			freeBefore := retentionPragma(t, s, "freelist_count")
			if freeBefore == 0 {
				t.Fatal("fixture has no free pages")
			}
			var released int64
			for remaining := freeBefore; remaining > 0; {
				result, err := s.CompactIncrementally(ctx, 17)
				if err != nil {
					t.Fatal(err)
				}
				if mode == "default" {
					if result != (IncrementalCompaction{}) || retentionPragma(t, s, "freelist_count") != freeBefore || retentionPragma(t, s, "page_count") != pagesBefore {
						t.Fatalf("default changed: %+v", result)
					}
					break
				}
				if !result.Incremental || result.PagesReleased <= 0 || result.PagesReleased > 17 || result.FreePages >= remaining {
					t.Fatalf("no bounded progress: %+v remaining=%d", result, remaining)
				}
				released += result.PagesReleased
				remaining = result.FreePages
			}
			checkpointRunDatabase(t, ctx, s)
			if mode == "incremental" {
				if retentionPragma(t, s, "freelist_count") != 0 || retentionPragma(t, s, "page_count") != pagesBefore-released || databaseFileSize(t, home) >= bytesBefore {
					t.Fatal("incremental slices did not shrink file")
				}
			} else if databaseFileSize(t, home) != bytesBefore {
				t.Fatal("default file changed")
			}
		})
	}
}

func TestCompactConvertsADefaultModeDatabase(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	home := t.TempDir()
	s := openDefaultRetentionFixture(t, home)
	seedPrunableRunEvents(t, ctx, s)
	checkpointRunDatabase(t, ctx, s)
	closeStore(t, s)
	before := readDatabaseBytes(t, DatabasePath(home))
	reader, err := OpenStorageReader(ctx, home)
	if err != nil {
		t.Fatal(err)
	}
	preview, err := reader.PreviewCompaction(ctx)
	if err != nil {
		t.Fatal(err)
	}
	closeStore(t, reader)
	if after := readDatabaseBytes(t, DatabasePath(home)); !bytes.Equal(before, after) {
		t.Fatal("preview changed source bytes")
	}
	s = openTestStore(t, ctx, home)
	defer closeStore(t, s)
	if retentionPragma(t, s, "auto_vacuum") != 0 {
		t.Fatal("preview converted source")
	}
	// Compact's fingerprint belongs to its connection, so take a fresh preview.
	preview, err = s.PreviewCompaction(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if after := readDatabaseBytes(t, DatabasePath(home)); !bytes.Equal(before, after) {
		t.Fatal("writer preview changed source bytes")
	}
	result, err := s.Compact(ctx, preview)
	if err != nil {
		t.Fatal(err)
	}
	if retentionPragma(t, s, "auto_vacuum") != 2 {
		t.Fatal("Compact did not convert source")
	}
	if absInt64(result.BytesAfter-preview.BytesAfter) > preview.ReconciliationToleranceBytes {
		t.Fatalf("preview=%+v result=%+v", preview, result)
	}
}

func TestOpenPreservesExistingAutoVacuumModes(t *testing.T) {
	t.Parallel()
	for _, mode := range []int{0, 1, 2} {
		t.Run(fmt.Sprint(mode), func(t *testing.T) {
			t.Parallel()
			ctx := context.Background()
			home := t.TempDir()
			path := DatabasePath(home)
			if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
				t.Fatal(err)
			}
			db, err := sql.Open("sqlite", migratorDSN(path))
			if err != nil {
				t.Fatal(err)
			}
			retentionExec(t, db, fmt.Sprintf("PRAGMA auto_vacuum = %d", mode))
			retentionExec(t, db, `CREATE TABLE mode_fixture (id INTEGER)`)
			retentionExec(t, db, `DROP TABLE mode_fixture`)
			if err := db.Close(); err != nil {
				t.Fatal(err)
			}
			s := openTestStore(t, ctx, home)
			defer closeStore(t, s)
			if got := retentionPragma(t, s, "auto_vacuum"); got != int64(mode) {
				t.Fatalf("existing mode=%d want=%d", got, mode)
			}
		})
	}
}

func TestPreviewCompactionOnIncrementalStorageReader(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	home := t.TempDir()
	s := openTestStore(t, ctx, home)
	seedPrunableRunEvents(t, ctx, s)
	checkpointRunDatabase(t, ctx, s)
	closeStore(t, s)
	before := readDatabaseBytes(t, DatabasePath(home))
	reader, err := OpenStorageReader(ctx, home)
	if err != nil {
		t.Fatal(err)
	}
	preview, err := reader.PreviewCompaction(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if retentionPragma(t, reader, "auto_vacuum") != 2 {
		t.Fatal("preview changed incremental mode")
	}
	closeStore(t, reader)
	if after := readDatabaseBytes(t, DatabasePath(home)); !bytes.Equal(before, after) {
		t.Fatal("incremental reader preview changed database bytes")
	}
	s = openTestStore(t, ctx, home)
	defer closeStore(t, s)
	// Compact revalidates a fingerprint belonging to its own connection.
	writerPreview, err := s.PreviewCompaction(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if writerPreview.BytesAfter != preview.BytesAfter {
		t.Fatalf("reader projected %d bytes; writer projected %d", preview.BytesAfter, writerPreview.BytesAfter)
	}
	result, err := s.Compact(ctx, writerPreview)
	if err != nil {
		t.Fatal(err)
	}
	if absInt64(result.BytesAfter-preview.BytesAfter) > preview.ReconciliationToleranceBytes {
		t.Fatalf("preview=%+v result=%+v", preview, result)
	}
}
