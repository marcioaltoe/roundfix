package store

import (
	"context"
	"testing"
	"time"

	"roundfix/internal/runevent"
)

// Suite: Run Event header projection
// Invariant: a payload-free read path projects cursor, batch, source, kind,
// summary, and creation time without touching the payload column, leaving the
// full read and every payload-reading consumer unchanged.
// Boundary IN: Store appends and the header projection read.
// Boundary OUT: the full read, the stream schema, and every command output.

func mustHeaderRun(t *testing.T, ctx context.Context, s *Store, summaries []string) string {
	t.Helper()
	run, err := s.CreateRun(ctx, sampleCreateRunRequest())
	if err != nil {
		t.Fatalf("create run: %v", err)
	}
	for index, summary := range summaries {
		event := sampleRunEvent(run.ID, summary)
		event.Kind = runevent.KindDaemonTask
		event.Source = runevent.SourceDaemon
		event.Batch = 1 + index%3
		if _, err := s.AppendRunEvent(ctx, event); err != nil {
			t.Fatalf("append %q: %v", summary, err)
		}
	}
	return run.ID
}

func TestHeaderProjectionProjectsOnlyHeaderColumns(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	s := openTestStore(t, ctx, t.TempDir())
	defer closeStore(t, s)

	runID := mustHeaderRun(t, ctx, s, []string{"one", "two", "three"})

	headers, err := s.RunEventHeadersAfter(ctx, runID, 0)
	if err != nil {
		t.Fatalf("list headers after cursor: %v", err)
	}
	if len(headers) != 3 {
		t.Fatalf("expected three headers, got %d", len(headers))
	}
	for index, header := range headers {
		if header.Cursor != int64(index+1) {
			t.Fatalf("header %d: expected cursor %d, got %d", index, index+1, header.Cursor)
		}
		if header.Source != runevent.SourceDaemon || header.Kind != runevent.KindDaemonTask {
			t.Fatalf("header %d: expected daemon task source/kind, got %s/%s", index, header.Source, header.Kind)
		}
		if header.Time.IsZero() {
			t.Fatalf("header %d: expected a creation time, got zero", index)
		}
		if header.Summary == "" {
			t.Fatalf("header %d: expected a summary", index)
		}
	}
}

func TestHeaderProjectionMatchesFullReadHeaders(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	s := openTestStore(t, ctx, t.TempDir())
	defer closeStore(t, s)

	runID := mustHeaderRun(t, ctx, s, []string{"one", "two", "three"})

	headers, err := s.RunEventHeadersAfter(ctx, runID, 0)
	if err != nil {
		t.Fatalf("list headers: %v", err)
	}
	events, err := s.RunEventsAfter(ctx, runID, 0, 10)
	if err != nil {
		t.Fatalf("list full events: %v", err)
	}
	if len(headers) != len(events) {
		t.Fatalf("header/full row count mismatch: %d vs %d", len(headers), len(events))
	}
	for index := range events {
		if headers[index].Cursor != events[index].Cursor {
			t.Fatalf("row %d: cursor mismatch %d vs %d", index, headers[index].Cursor, events[index].Cursor)
		}
		if headers[index].Batch != events[index].Event.Batch {
			t.Fatalf("row %d: batch mismatch %d vs %d", index, headers[index].Batch, events[index].Event.Batch)
		}
		if headers[index].Source != events[index].Event.Source {
			t.Fatalf("row %d: source mismatch %s vs %s", index, headers[index].Source, events[index].Event.Source)
		}
		if headers[index].Kind != events[index].Event.Kind {
			t.Fatalf("row %d: kind mismatch %s vs %s", index, headers[index].Kind, events[index].Event.Kind)
		}
		if headers[index].Summary != events[index].Event.Summary {
			t.Fatalf("row %d: summary mismatch %q vs %q", index, headers[index].Summary, events[index].Event.Summary)
		}
		if !headers[index].Time.Equal(events[index].Event.Time) {
			t.Fatalf("row %d: time mismatch %v vs %v", index, headers[index].Time, events[index].Event.Time)
		}
	}
}

func TestRunEventHeadersAfterCursorSkipsOlder(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	s := openTestStore(t, ctx, t.TempDir())
	defer closeStore(t, s)

	runID := mustHeaderRun(t, ctx, s, []string{"one", "two", "three"})

	page, err := s.RunEventHeadersAfter(ctx, runID, 2)
	if err != nil {
		t.Fatalf("list headers after cursor: %v", err)
	}
	if len(page) != 1 || page[0].Cursor != 3 {
		t.Fatalf("expected cursor [3], got %+v", page)
	}

	empty, err := s.RunEventHeadersAfter(ctx, runID, 3)
	if err != nil {
		t.Fatalf("list headers at tail: %v", err)
	}
	if len(empty) != 0 {
		t.Fatalf("expected no headers at tail, got %+v", empty)
	}
}

func TestRunEventHeadersAfterRequiresRunAndCursorForward(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	s := openTestStore(t, ctx, t.TempDir())
	defer closeStore(t, s)

	if _, err := s.RunEventHeadersAfter(ctx, "   ", 0); err == nil {
		t.Fatal("expected blank Run ID to be rejected")
	}
	if _, err := s.RunEventHeadersAfter(ctx, "run_missing", 0); err != nil {
		t.Fatalf("expected missing Run to list no headers, got %v", err)
	}

	runID := mustHeaderRun(t, ctx, s, []string{"one", "two"})
	if _, err := s.RunEventHeadersAfter(ctx, runID, -1); err != nil {
		t.Fatalf("expected negative cursor to list all headers, got %v", err)
	}
}

func TestEventHeadersOrderAscendingAcrossBatches(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	s := openTestStore(t, ctx, t.TempDir())
	defer closeStore(t, s)

	run, err := s.CreateRun(ctx, sampleCreateRunRequest())
	if err != nil {
		t.Fatalf("create run: %v", err)
	}
	first := sampleRunEvent(run.ID, "first")
	first.Batch = 2
	first.Time = time.Date(2026, 6, 10, 12, 0, 0, 0, time.UTC)
	second := sampleRunEvent(run.ID, "second")
	second.Batch = 3
	second.Time = time.Date(2026, 6, 10, 12, 0, 1, 0, time.UTC)
	if _, err := s.AppendRunEvent(ctx, first); err != nil {
		t.Fatalf("append first event: %v", err)
	}
	if _, err := s.AppendRunEvent(ctx, second); err != nil {
		t.Fatalf("append second event: %v", err)
	}

	headers, err := s.RunEventHeadersAfter(ctx, run.ID, 0)
	if err != nil {
		t.Fatalf("list headers: %v", err)
	}
	if len(headers) != 2 {
		t.Fatalf("expected two headers, got %d", len(headers))
	}
	if headers[0].Cursor != 1 || headers[1].Cursor != 2 {
		t.Fatalf("expected cursors [1 2], got %+v", headers)
	}
	if !headers[1].Time.After(headers[0].Time) {
		t.Fatalf("expected ascending creation times, got %v then %v", headers[0].Time, headers[1].Time)
	}
}