package store

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"sync"
	"testing"
	"time"

	"roundfix/internal/runevent"
)

// batchTestWriter builds a test Store whose journal batch size and linger are
// small so tests can observe the boundaries deterministically instead of by
// timing. The production constants remain internal and untouched.
type batchTestWriter struct {
	store *Store
	sink  runevent.Sink
}

func newBatchTestWriter(t *testing.T, batchSize int, linger time.Duration) *batchTestWriter {
	t.Helper()
	ctx := context.Background()
	homeDir := t.TempDir()
	stores := openTestStoreBatch(t, ctx, homeDir, batchSize, linger)
	sink := stores.JournalSink()
	return &batchTestWriter{store: stores, sink: sink}
}

// openTestStoreBatch opens a writer Store with a test-configured journal
// writer. It must not leak the test settings into the Store-scoped writer the
// rest of the suite uses.
func openTestStoreBatch(t *testing.T, ctx context.Context, homeDir string, batchSize int, linger time.Duration) *Store {
	t.Helper()
	store, err := Open(ctx, homeDir)
	if err != nil {
		t.Fatalf("open store: %v", err)
	}
	// Replace the production writer with a test-sized one. The word before
	// newJournalWriter ran inside Open, so the Store-scoped writer is recreated
	// here at the test boundary.
	writer := newJournalWriter(store)
	writer.batchSize = batchSize
	writer.maxLinger = linger
	store.journal = writer
	return store
}

func (w *batchTestWriter) close(t *testing.T) {
	t.Helper()
	if err := w.store.Close(); err != nil {
		t.Fatalf("close store: %v", err)
	}
}

func newRunForBatch(t *testing.T, store *Store) string {
	t.Helper()
	ctx := context.Background()
	run, err := store.CreateRun(ctx, sampleCreateRunRequest())
	if err != nil {
		t.Fatalf("create Run: %v", err)
	}
	return run.ID
}

func batchAgentEvent(runID string, summary string) runevent.RunEvent {
	return runevent.RunEvent{
		RunID:   runID,
		Source:  runevent.SourceAgent,
		Kind:    runevent.KindAgentMessage,
		Summary: summary,
		Time:    time.Date(2026, 6, 10, 12, 0, 0, 0, time.UTC),
		Payload: []byte(`{"sessionId":"sess-1","update":{"text":"` + summary + `"}}`),
	}
}

func batchOutcomeEvent(runID string) runevent.RunEvent {
	return runevent.RunEvent{
		RunID:   runID,
		Source:  runevent.SourceDaemon,
		Kind:    runevent.KindDaemonOutcome,
		Summary: "terminal outcome",
		Time:    time.Date(2026, 6, 10, 12, 1, 0, 0, time.UTC),
		Payload: []byte(`{"state":"clean"}`),
	}
}

func countPublishableEvents(t *testing.T, ctx context.Context, store *Store, runID string) int {
	t.Helper()
	events, err := store.RunEventsAfter(ctx, runID, 0, 1000)
	if err != nil {
		t.Fatalf("read published events: %v", err)
	}
	return len(events)
}

func TestBatchClosesOnCountLingerAndImmediate(t *testing.T) {
	t.Parallel()

	t.Run("count closes a batch", func(t *testing.T) {
		w := newBatchTestWriter(t, 3, time.Hour)
		defer w.close(t)
		w.store.FlushJournal(context.Background())
		runID := newRunForBatch(t, w.store)

		for i := 0; i < 3; i++ {
			if err := w.sink.Publish(context.Background(), batchAgentEvent(runID, "event")); err != nil {
				t.Fatalf("publish event %d: %v", i, err)
			}
		}
		if got := w.store.journal.pendingCount(); got != 0 {
			t.Fatalf("expected count flush to empty pending, got %d pending", got)
		}
		if got := countPublishableEvents(t, context.Background(), w.store, runID); got != 3 {
			t.Fatalf("expected 3 committed events after count flush, got %d", got)
		}
	})

	t.Run("linger closes a quiet batch", func(t *testing.T) {
		w := newBatchTestWriter(t, 100, 50*time.Millisecond)
		defer w.close(t)
		runID := newRunForBatch(t, w.store)

		if err := w.sink.Publish(context.Background(), batchAgentEvent(runID, "event")); err != nil {
			t.Fatalf("publish event: %v", err)
		}
		deadline := time.Now().Add(2 * time.Second)
		for {
			if w.store.journal.pendingCount() == 0 {
				break
			}
			if time.Now().After(deadline) {
				t.Fatal("linger deadline did not close the batch")
			}
			time.Sleep(10 * time.Millisecond)
		}
		if got := countPublishableEvents(t, context.Background(), w.store, runID); got != 1 {
			t.Fatalf("expected 1 committed event after linger, got %d", got)
		}
	})

	t.Run("immediate flushes prior events", func(t *testing.T) {
		w := newBatchTestWriter(t, 100, time.Hour)
		defer w.close(t)
		runID := newRunForBatch(t, w.store)

		if err := w.sink.Publish(context.Background(), batchAgentEvent(runID, "before")); err != nil {
			t.Fatalf("publish agent event: %v", err)
		}
		if got := w.store.journal.pendingCount(); got != 1 {
			t.Fatalf("expected pending before immediate, got %d", got)
		}
		if err := w.sink.Publish(context.Background(), batchOutcomeEvent(runID)); err != nil {
			t.Fatalf("publish immediate event: %v", err)
		}
		if got := w.store.journal.pendingCount(); got != 0 {
			t.Fatalf("expected immediate flush to empty pending, got %d", got)
		}
		if got := countPublishableEvents(t, context.Background(), w.store, runID); got != 2 {
			t.Fatalf("expected both events committed, got %d", got)
		}
	})

	t.Run("explicit flush closes a batch", func(t *testing.T) {
		w := newBatchTestWriter(t, 100, time.Hour)
		defer w.close(t)
		runID := newRunForBatch(t, w.store)

		if err := w.sink.Publish(context.Background(), batchAgentEvent(runID, "held")); err != nil {
			t.Fatalf("publish event: %v", err)
		}
		if got := countPublishableEvents(t, context.Background(), w.store, runID); got != 0 {
			t.Fatalf("expected held event not yet committed, got %d", got)
		}
		if err := w.store.FlushJournal(context.Background()); err != nil {
			t.Fatalf("flush journal: %v", err)
		}
		if got := countPublishableEvents(t, context.Background(), w.store, runID); got != 1 {
			t.Fatalf("expected event committed after flush, got %d", got)
		}
	})
}

func TestBatchPublishPreservesOrderAndContiguousCursors(t *testing.T) {
	t.Parallel()
	w := newBatchTestWriter(t, 2, time.Hour)
	defer w.close(t)
	runID := newRunForBatch(t, w.store)

	for _, summary := range []string{"a", "b", "c", "d"} {
		if err := w.sink.Publish(context.Background(), batchAgentEvent(runID, summary)); err != nil {
			t.Fatalf("publish %s: %v", summary, err)
		}
	}
	events, err := w.store.RunEventsAfter(context.Background(), runID, 0, 100)
	if err != nil {
		t.Fatalf("read events: %v", err)
	}
	if len(events) != 4 {
		t.Fatalf("expected 4 committed events, got %d", len(events))
	}
	for index, event := range events {
		if event.Cursor != int64(index+1) {
			t.Fatalf("expected contiguous cursor %d, got %d", index+1, event.Cursor)
		}
		if event.Event.Summary != []string{"a", "b", "c", "d"}[index] {
			t.Fatalf("expected publisher order, got %q at %d", event.Event.Summary, index)
		}
	}
}

func TestBatchPreservesPayloadBytes(t *testing.T) {
	t.Parallel()
	w := newBatchTestWriter(t, 2, time.Hour)
	defer w.close(t)
	runID := newRunForBatch(t, w.store)

	raw := []byte(`{"very":"nested","payload":{"with":[1,2,3],"and":"spaces ~!@#$%^&*()"}}`)
	event := batchAgentEvent(runID, "raw")
	event.Payload = raw
	if err := w.sink.Publish(context.Background(), event); err != nil {
		t.Fatalf("publish event: %v", err)
	}
	if err := w.store.FlushJournal(context.Background()); err != nil {
		t.Fatalf("flush journal: %v", err)
	}
	committed, err := w.store.RunEventsAfter(context.Background(), runID, 0, 10)
	if err != nil {
		t.Fatalf("read event: %v", err)
	}
	if len(committed) != 1 {
		t.Fatalf("expected 1 event, got %d", len(committed))
	}
	if string(committed[0].Event.Payload) != string(raw) {
		t.Fatalf("payload bytes changed: got %q want %q", committed[0].Event.Payload, raw)
	}
}

func TestBatchRejectsPublishAfterClose(t *testing.T) {
	t.Parallel()
	w := newBatchTestWriter(t, 100, time.Hour)
	defer w.close(t)
	runID := newRunForBatch(t, w.store)
	if err := w.store.CloseJournal(context.Background()); err != nil {
		t.Fatalf("close journal: %v", err)
	}
	err := w.sink.Publish(context.Background(), batchAgentEvent(runID, "late"))
	if err == nil || !strings.Contains(err.Error(), "closed") {
		t.Fatalf("expected publish after close to be rejected, got %v", err)
	}
}

// pendingCount reports how many events are currently in the pending batch.
func (w *journalWriter) pendingCount() int {
	w.mu.Lock()
	defer w.mu.Unlock()
	return len(w.pending)
}

// seedAssignmentRows inserts rows at the exact assigned cursor range so tests
// can exercise the ambiguous-commit reconciliation without a real SQLite
// commit failure.
func (w *batchTestWriter) seedAssignmentRows(t *testing.T, assignment batchAssignment) {
	t.Helper()
	ctx := context.Background()
	if err := w.store.withWriteTx(ctx, "test assign", func(tx *sql.Tx) error {
		for index, event := range assignment.events {
			if err := insertRunEventAtCursor(ctx, tx, event, assignment.start+int64(index)); err != nil {
				return err
			}
		}
		return nil
	}); err != nil {
		t.Fatalf("seed assigned rows: %v", err)
	}
}

func TestBatchAmbiguousCommit(t *testing.T) {
	t.Parallel()

	newFixture := func(t *testing.T) *batchTestWriter {
		w := newBatchTestWriter(t, 100, time.Hour)
		w.store.FlushJournal(context.Background())
		t.Cleanup(func() { w.close(t) })
		return w
	}

	t.Run("exact match settles the batch without retry", func(t *testing.T) {
		w := newFixture(t)
		runID := newRunForBatch(t, w.store)
		events := []runevent.RunEvent{batchAgentEvent(runID, "a"), batchAgentEvent(runID, "b")}
		assignment := []batchAssignment{{runID: runID, start: 1, events: events}}
		w.seedAssignmentRows(t, assignment[0])

		if err := reconcileBatchCommit(context.Background(), w.store, events, assignment); err != nil {
			t.Fatalf("exact match should settle without retry, got %v", err)
		}
		if got := countPublishableEvents(t, context.Background(), w.store, runID); got != 2 {
			t.Fatalf("expected exact-match settled rows unchanged, got %d events", got)
		}
	})

	t.Run("no rows retries once with the same cursors", func(t *testing.T) {
		w := newFixture(t)
		runID := newRunForBatch(t, w.store)
		events := []runevent.RunEvent{batchAgentEvent(runID, "a")}
		assignment := []batchAssignment{{runID: runID, start: 1, events: events}}

		if err := reconcileBatchCommit(context.Background(), w.store, events, assignment); err != nil {
			t.Fatalf("no-rows retry should succeed, got %v", err)
		}
		// The retry must have inserted at cursor 1 (the original range), not a
		// freshly allocated one.
		committed, err := w.store.RunEventsAfter(context.Background(), runID, 0, 10)
		if err != nil {
			t.Fatalf("read retried events: %v", err)
		}
		if len(committed) != 1 || committed[0].Cursor != 1 {
			t.Fatalf("expected one retried event at cursors 1, got %+v", committed)
		}
	})

	t.Run("partial match fails as corruption", func(t *testing.T) {
		w := newFixture(t)
		runID := newRunForBatch(t, w.store)
		events := []runevent.RunEvent{batchAgentEvent(runID, "a"), batchAgentEvent(runID, "b")}
		assignment := []batchAssignment{{runID: runID, start: 1, events: events}}
		// Only the first of two assigned events is present.
		if err := w.store.withWriteTx(context.Background(), "test partial", func(tx *sql.Tx) error {
			return insertRunEventAtCursor(context.Background(), tx, events[0], 1)
		}); err != nil {
			t.Fatalf("seed partial rows: %v", err)
		}

		err := reconcileBatchCommit(context.Background(), w.store, events, assignment)
		if err == nil || !strings.Contains(err.Error(), "partial or different") {
			t.Fatalf("expected partial match to fail as corruption, got %v", err)
		}
	})

	t.Run("different payload fails as corruption", func(t *testing.T) {
		w := newFixture(t)
		runID := newRunForBatch(t, w.store)
		events := []runevent.RunEvent{batchAgentEvent(runID, "a")}
		assignment := []batchAssignment{{runID: runID, start: 1, events: events}}
		different := events[0]
		different.Payload = []byte(`{"different":"bytes"}`)
		if err := w.store.withWriteTx(context.Background(), "test different", func(tx *sql.Tx) error {
			return insertRunEventAtCursor(context.Background(), tx, different, 1)
		}); err != nil {
			t.Fatalf("seed different rows: %v", err)
		}

		err := reconcileBatchCommit(context.Background(), w.store, events, assignment)
		if err == nil || !strings.Contains(err.Error(), "partial or different") {
			t.Fatalf("expected different-payload match to fail as corruption, got %v", err)
		}
	})
}

func TestBatchBeginInsertCommitFailurePreservesBatch(t *testing.T) {
	t.Parallel()

	t.Run("insert failure preserves the batch and fails the Run", func(t *testing.T) {
		w := newBatchTestWriter(t, 100, time.Hour)
		// No Run exists; appending to a missing Run must fail, preserving the
		// whole pending batch. With batching the failure surfaces at flush, not
		// at publish.
		evict := runevent.RunEvent{
			RunID:   "run_missing",
			Source:  runevent.SourceAgent,
			Kind:    runevent.KindAgentMessage,
			Summary: "orphan",
			Time:    time.Now().UTC(),
		}
		if err := w.sink.Publish(context.Background(), evict); err != nil {
			t.Fatalf("publish defers the failure to flush, got %v", err)
		}
		if got := w.store.journal.pendingCount(); got != 1 {
			t.Fatalf("expected pending before flush, got %d", got)
		}
		if err := w.store.FlushJournal(context.Background()); err == nil {
			t.Fatal("expected flush of a missing-Run batch to fail")
		}
		if got := w.store.journal.pendingCount(); got != 1 {
			t.Fatalf("expected failed flush to preserve the whole pending batch, got %d", got)
		}
		// A second, valid publish must not drop the preserved orphan batch: the
		// flush still fails on the preserved first event.
		runID := newRunForBatch(t, w.store)
		if err := w.sink.Publish(context.Background(), batchAgentEvent(runID, "ok")); err != nil {
			t.Fatalf("publish valid event: %v", err)
		}
		if err := w.store.FlushJournal(context.Background()); err == nil {
			t.Fatalf("expected preserved orphan batch to still fail flush, got nil")
		}
		if got := w.store.journal.pendingCount(); got != 2 {
			t.Fatalf("expected 2 preserved pending events, got %d", got)
		}
		// A failed Close also preserves the batch and remains retryable.
		if err := w.store.CloseJournal(context.Background()); err == nil {
			t.Fatal("expected CloseInt to preserve the unflushable batch with an error")
		}
		w.store.db.Close()
		w.store.writeLockFile.Close()
	})

	t.Run("multiple sinks share one store-scoped writer", func(t *testing.T) {
		w := newBatchTestWriter(t, 100, time.Hour)
		defer w.close(t)
		// Two independent sink handles publish into the SAME writer: the batch
		// boundary is store-scoped, not sink-scoped.
		sinkA := w.store.JournalSink()
		sinkB := w.store.JournalSink()
		runID := newRunForBatch(t, w.store)
		if err := sinkA.Publish(context.Background(), batchAgentEvent(runID, "a")); err != nil {
			t.Fatalf("publish a: %v", err)
		}
		if err := sinkB.Publish(context.Background(), batchAgentEvent(runID, "b")); err != nil {
			t.Fatalf("publish b: %v", err)
		}
		if got := w.store.journal.pendingCount(); got != 2 {
			t.Fatalf("expected both sinks to share one pending batch, got %d", got)
		}
		if err := w.store.FlushJournal(context.Background()); err != nil {
			t.Fatalf("flush: %v", err)
		}
		if got := countPublishableEvents(t, context.Background(), w.store, runID); got != 2 {
			t.Fatalf("expected 2 events from shared writer, got %d", got)
		}
	})

	t.Run("concurrent publishers keep order and contiguous cursors", func(t *testing.T) {
		w := newBatchTestWriter(t, 5, time.Hour)
		defer w.close(t)
		runID := newRunForBatch(t, w.store)

		const publishers = 4
		const each = 25
		var wg sync.WaitGroup
		for p := 0; p < publishers; p++ {
			wg.Add(1)
			go func(p int) {
				defer wg.Done()
				for i := 0; i < each; i++ {
					_ = w.sink.Publish(context.Background(), batchAgentEvent(runID, fmt.Sprintf("p%d-%d", p, i)))
				}
			}(p)
		}
		wg.Wait()
		if err := w.store.FlushJournal(context.Background()); err != nil {
			t.Fatalf("flush after concurrent publish: %v", err)
		}
		events, err := w.store.RunEventsAfter(context.Background(), runID, 0, 1000)
		if err != nil {
			t.Fatalf("read events: %v", err)
		}
		if len(events) != publishers*each {
			t.Fatalf("expected %d events, got %d", publishers*each, len(events))
		}
		for index, event := range events {
			if event.Cursor != int64(index+1) {
				t.Fatalf("expected contiguous cursor %d, got %d", index+1, event.Cursor)
			}
		}
	})
}
