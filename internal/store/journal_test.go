package store

import (
	"bytes"
	"context"
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

func TestRunEventsAfterCursorReturnsOnlyNewerAndRespectsLimit(t *testing.T) {
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
