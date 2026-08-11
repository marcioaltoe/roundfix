package store

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"strings"
	"sync"
	"time"

	"roundfix/internal/runevent"
)

// Journal batching limits are internal constants, not configuration or CLI
// surface. Durability moves from per-event to per-batch: `synchronous` stays
// at its default, so a crash can lose at most the batch in flight.
const (
	journalBatchSize = 128
	journalMaxLinger = 100 * time.Millisecond
)

// writeCommitError marks a failure to commit a write transaction, as opposed
// to a failure to begin or to run an insert. The message contract is
// unchanged: "commit <operation>: <cause>". The batch writer matches on it so
// an ambiguous commit can be reconciled by reading the assigned cursor range,
// while begin/insert failures are returned as ordinary errors that preserve
// the pending batch.
type writeCommitError struct {
	operation string
	cause     error
}

func (err writeCommitError) Error() string {
	return fmt.Sprintf("commit %s: %v", err.operation, err.cause)
}

func (err writeCommitError) Unwrap() error {
	return err.cause
}

// journalImmediate reports whether an event requires immediate durability.
// Such events flush any pending batch before they are appended and commit on
// their own, so the events a post-mortem needs most never sit in one. Agent
// selection transitions are immediate because an ACP adapter prepares the
// next session only after reading the durable fallback notification back from
// the journal.
func journalImmediate(event runevent.RunEvent) bool {
	switch event.Kind {
	case runevent.KindDaemonOutcome, runevent.KindDaemonVerification,
		runevent.KindDaemonAgentSelectionAttempt,
		runevent.KindDaemonAgentSelectionActive,
		runevent.KindDaemonAgentSelectionFallback,
		runevent.KindDaemonAgentSelectionExhausted,
		runevent.KindDaemonAgentSelectionClosed:
		return true
	default:
		return false
	}
}

// errJournalClosed rejects a Publish after a successful CloseJournal.
var errJournalClosed = errors.New("publish Run Event: journal writer is closed")

// journalWriter is the Store-scoped batched Run Event writer. Every sink the
// Store hands out shares this one writer, so batch boundaries are global to
// the process rather than per sink. It appends events into a pending batch,
// commits the whole batch in one write transaction, and only removes events
// from the pending batch after a successful commit.
type journalWriter struct {
	store *Store

	// batchSize and maxLinger default to the package constants and are the
	// production batching limits. They are unexported fields so tests can
	// observe the boundaries deterministically without changing the constants,
	// which remain internal and are not configuration or CLI surface.
	batchSize int
	maxLinger time.Duration

	mu       sync.Mutex
	pending  []runevent.RunEvent
	inFlight bool
	closed   bool
	timer    *time.Timer

	// commitDone is closed when the in-flight commit completes and is
	// re-allocated on each commit so waiters can observe progress. It lets an
	// immediate-durability publish order its event after a commit that is
	// already draining, without holding the mutex across the database write.
	commitDone chan struct{}
}

func newJournalWriter(store *Store) *journalWriter {
	return &journalWriter{
		store:      store,
		batchSize:  journalBatchSize,
		maxLinger:  journalMaxLinger,
		commitDone: make(chan struct{}),
	}
}

// journalSink adapts the Store-scoped journal writer to runevent.Sink.
type journalSink struct {
	writer *journalWriter
}

func (sink journalSink) Publish(ctx context.Context, event runevent.RunEvent) error {
	return sink.writer.publish(ctx, event)
}

// publish appends one event to the pending batch, closing the batch before an
// immediate-durability event and on count flush. The database commit runs
// outside the mutex (see flush), so a publish never blocks behind another
// publisher's acquireWriteLock wait or SQLite transaction.
func (w *journalWriter) publish(ctx context.Context, event runevent.RunEvent) error {
	if journalImmediate(event) {
		// Flush whatever is pending first (requiring its own commit), then
		// append this event to a fresh batch and flush it immediately so it is
		// durable on its own.
		if err := w.flush(ctx); err != nil {
			return err
		}
		w.mu.Lock()
		if w.closed {
			w.mu.Unlock()
			return errJournalClosed
		}
		w.pending = append(w.pending, event)
		w.mu.Unlock()
		return w.flush(ctx)
	}
	w.mu.Lock()
	if w.closed {
		w.mu.Unlock()
		return errJournalClosed
	}
	w.pending = append(w.pending, event)
	if len(w.pending) >= w.batchSize {
		w.mu.Unlock()
		return w.flush(ctx)
	}
	w.armLingerLocked()
	w.mu.Unlock()
	return nil
}

// armLingerLocked schedules a commit after journalMaxLinger since the current
// batch's first event, so a quiet publisher's last line is not held
// indefinitely. A batch's deadline runs from its first event, so a later event
// in the same batch must not extend it: once armed, the timer is left running.
func (w *journalWriter) armLingerLocked() {
	if w.inFlight {
		return
	}
	if w.timer != nil {
		return
	}
	w.timer = time.AfterFunc(w.maxLinger, func() {
		w.mu.Lock()
		pending := len(w.pending)
		w.mu.Unlock()
		if err := w.flush(context.Background()); err != nil {
			slog.Error("journal linger flush failed",
				slog.Int("pending_events", pending),
				slog.Any("error", err))
		}
	})
}

// FlushJournal commits any pending batch now.
func (w *journalWriter) flush(ctx context.Context) error {
	for {
		w.mu.Lock()
		batch, wait, err := w.drainLocked()
		var done <-chan struct{}
		if wait {
			done = w.commitDone
		}
		w.mu.Unlock()
		if err != nil {
			return err
		}
		if batch != nil {
			if err := w.commitDrained(ctx, batch); err != nil {
				return err
			}
			continue
		}
		if wait {
			// A commit is draining concurrently. Wait for it to finish so
			// events appended (or re-queued) behind it are not skipped; then
			// re-run the drain loop.
			if err := waitCommit(ctx, done); err != nil {
				return err
			}
			continue
		}
		return nil
	}
}

// waitCommit blocks until the in-flight commit's commitDone channel closes or
// ctx is done. It releases the mutex before blocking so the committing
// goroutine can finish and close the channel.
func waitCommit(ctx context.Context, done <-chan struct{}) error {
	select {
	case <-done:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

// drainLocked takes the current pending batch out from under the mutex and
// marks the writer in-flight, or reports whether nothing can be drained
// because a commit is already in flight. It must be called with w.mu held.
// Databases writes happen in commitDrained, after the mutex is released, so
// publishers are never serialized behind the advisory write lock or a SQLite
// transaction.
func (w *journalWriter) drainLocked() (batch []runevent.RunEvent, inFlight bool, err error) {
	if len(w.pending) == 0 {
		if !w.inFlight {
			w.stopTimerLocked()
		}
		return nil, false, nil
	}
	if w.inFlight {
		// A linger or count flush is already committing; tell the caller to
		// wait on commitDone instead of skipping these pending events.
		return nil, true, nil
	}
	w.inFlight = true
	w.stopTimerLocked()
	batch = w.pending
	w.pending = nil
	return batch, false, nil
}

// commitDrained commits a drained batch outside the mutex, restoring the whole
// batch (plus anything published meanwhile) on failure. Committing outside the
// mutex keeps the write lock wait and SQLite transaction out of the mutex
// critical section.
func (w *journalWriter) commitDrained(ctx context.Context, batch []runevent.RunEvent) error {
	err := commitJournalBatch(ctx, w.store, batch)
	w.mu.Lock()
	w.inFlight = false
	// Signal any waiter blocked in flush that the in-flight commit (and its
	// commitDone observation) finished, then re-arm the channel for the next
	// commit. A single drain owns this transition: only the goroutine that
	// observed inFlight true and drained the batch reaches commitDrained, so
	// there is exactly one close per commit.
	close(w.commitDone)
	w.commitDone = make(chan struct{})
	if err != nil {
		w.pending = append(batch, w.pending...)
	} else if len(w.pending) > 0 && !w.closed && w.timer == nil {
		// Re-arm the linger timer for events committed into a fresh batch.
		w.armLingerLocked()
	}
	w.mu.Unlock()
	return err
}

func (w *journalWriter) stopTimerLocked() {
	if w.timer != nil {
		w.timer.Stop()
		w.timer = nil
	}
}

// close flushes the pending batch and marks the writer closed only after that
// flush commits. A failed Close preserves the batch and remains retryable;
// every later Publish after a successful Close is rejected.
func (w *journalWriter) close(ctx context.Context) error {
	for {
		w.mu.Lock()
		idle := len(w.pending) == 0 && !w.inFlight && w.timer == nil
		if w.closed || idle {
			w.closed = true
			w.stopTimerLocked()
			w.mu.Unlock()
			return nil
		}
		w.mu.Unlock()
		if err := w.flush(ctx); err != nil {
			return err
		}
	}
}

// batchAssignment records the contiguous cursor range a Run received inside a
// write transaction.
type batchAssignment struct {
	runID  string
	start  int64
	events []runevent.RunEvent
}

// commitJournalBatch commits one ordered batch in a single write transaction,
// allocating a contiguous cursor range per Run, and reconciles an ambiguous
// commit by reading the assigned range back.
func commitJournalBatch(ctx context.Context, store *Store, batch []runevent.RunEvent) error {
	var assignment []batchAssignment
	err := store.withWriteTx(ctx, "Run Event batch append", func(tx *sql.Tx) error {
		var err error
		assignment, err = planAndInsertBatch(ctx, tx, batch)
		return err
	})
	if err != nil {
		if isWriteCommitError(err) {
			return reconcileBatchCommit(ctx, store, batch, assignment)
		}
		return err
	}
	return nil
}

// isWriteCommitError reports whether err is a writeCommitError value, the
// marker withWriteTx returns when the transaction commit itself failed.
func isWriteCommitError(err error) bool {
	var commitErr writeCommitError
	return errors.As(err, &commitErr)
}

// planAndInsertBatch inserts events in publisher order, allocating one
// contiguous cursor range per Run inside the transaction. A missing Run fails
// clearly, as a per-event append does today.
func planAndInsertBatch(ctx context.Context, tx *sql.Tx, batch []runevent.RunEvent) ([]batchAssignment, error) {
	var order []string
	groups := map[string][]runevent.RunEvent{}
	for _, event := range batch {
		runID := strings.TrimSpace(event.RunID)
		if runID == "" {
			return nil, errors.New("append Run Event: Run ID is required")
		}
		if _, ok := groups[runID]; !ok {
			order = append(order, runID)
		}
		groups[runID] = append(groups[runID], event)
	}

	assignments := make([]batchAssignment, 0, len(order))
	for _, runID := range order {
		var exists int
		err := tx.QueryRowContext(ctx, `SELECT 1 FROM runs WHERE id = ?`, runID).Scan(&exists)
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("append Run Event: Run %q does not exist", runID)
		}
		if err != nil {
			return nil, fmt.Errorf("append Run Event: check Run %q: %w", runID, err)
		}
		var max sql.NullInt64
		if err := tx.QueryRowContext(ctx, `SELECT MAX(cursor) FROM run_events WHERE run_id = ?`, runID).Scan(&max); err != nil {
			return nil, fmt.Errorf("append Run Event: read cursor for Run %q: %w", runID, err)
		}
		start := max.Int64 + 1
		events := groups[runID]
		for index, event := range events {
			if err := insertRunEventAtCursor(ctx, tx, event, start+int64(index)); err != nil {
				return nil, err
			}
		}
		assignments = append(assignments, batchAssignment{runID: runID, start: start, events: events})
	}
	return assignments, nil
}

// insertRunEventAtCursor inserts one event at an explicit cursor, preserving
// the raw producer payload byte-for-byte (ADR 0008: nothing is rewritten,
// pruned, or compressed).
func insertRunEventAtCursor(ctx context.Context, tx *sql.Tx, event runevent.RunEvent, cursor int64) error {
	_, err := tx.ExecContext(ctx, `
INSERT INTO run_events (
	run_id, cursor, batch, source, kind, review_issue,
	tool_id, tool_state, summary, created_at, payload
)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		strings.TrimSpace(event.RunID),
		cursor,
		event.Batch,
		string(event.Source),
		string(event.Kind),
		event.ReviewIssue,
		event.ToolID,
		event.ToolState,
		event.Summary,
		formatTime(event.Time),
		string(event.Payload),
	)
	if err != nil {
		return fmt.Errorf("insert Run Event for Run %q: %w", event.RunID, err)
	}
	return nil
}

// reconcileBatchCommit resolves an ambiguous commit by reading the assigned
// cursor range: an exact field-and-payload match settles the batch, no rows
// permits one retry holding the same cursors, and a partial or different match
// fails as corruption.
func reconcileBatchCommit(ctx context.Context, store *Store, batch []runevent.RunEvent, assignment []batchAssignment) error {
	settled, anyRows, err := readBatchAssignment(ctx, store, assignment)
	if err != nil {
		return err
	}
	if settled {
		return nil
	}
	if anyRows {
		return fmt.Errorf("commit Run Event batch append: ambiguous commit left a partial or different cursor range")
	}
	// No rows committed: one retry with the same cursor range.
	err = store.withWriteTx(ctx, "Run Event batch append retry", func(tx *sql.Tx) error {
		for _, a := range assignment {
			for index, event := range a.events {
				if err := insertRunEventAtCursor(ctx, tx, event, a.start+int64(index)); err != nil {
					return err
				}
			}
		}
		return nil
	})
	if err != nil {
		if isWriteCommitError(err) {
			settled, anyRows, rerr := readBatchAssignment(ctx, store, assignment)
			if rerr != nil {
				return rerr
			}
			if settled {
				return nil
			}
			if anyRows {
				return fmt.Errorf("commit Run Event batch append retry: ambiguous commit left a partial or different cursor range")
			}
			return fmt.Errorf("commit Run Event batch append retry: ambiguous commit left no rows after one retry")
		}
		return err
	}
	return nil
}

// readBatchAssignment classifies the stored state of the assigned cursor range
// against the intended events. settled is true only when every assigned event
// matches field-and-payload; anyRows reports whether any row exists at all.
func readBatchAssignment(ctx context.Context, store *Store, assignment []batchAssignment) (settled bool, anyRows bool, err error) {
	allSettled := true
	for _, a := range assignment {
		rows, err := store.db.QueryContext(ctx, `
SELECT cursor, batch, source, kind, review_issue, tool_id, tool_state, summary, created_at, payload
FROM run_events
WHERE run_id = ? AND cursor BETWEEN ? AND ?
ORDER BY cursor ASC`,
			a.runID,
			a.start,
			a.start+int64(len(a.events))-1,
		)
		if err != nil {
			return false, false, fmt.Errorf("read ambiguous Run Event batch for Run %q: %w", a.runID, err)
		}
		stored, rowsErr := collectStoredEvents(rows)
		if rowsErr != nil {
			return false, false, rowsErr
		}
		if len(stored) > 0 {
			anyRows = true
		}
		if len(stored) != len(a.events) || !storedEventsMatch(a.events, stored) {
			allSettled = false
		}
	}
	return allSettled, anyRows, nil
}

// collectStoredEvents reads a cursor range of stored journal events.
func collectStoredEvents(rows *sql.Rows) ([]storedJournalEvent, error) {
	defer func() { _ = rows.Close() }()
	stored := []storedJournalEvent{}
	for rows.Next() {
		var row storedJournalEvent
		var source string
		var kind string
		var createdAt string
		var payload string
		if err := rows.Scan(
			&row.cursor,
			&row.batch,
			&source,
			&kind,
			&row.reviewIssue,
			&row.toolID,
			&row.toolState,
			&row.summary,
			&createdAt,
			&payload,
		); err != nil {
			return nil, fmt.Errorf("scan ambiguous Run Event batch: %w", err)
		}
		row.source = source
		row.kind = kind
		row.createdAt = createdAt
		row.payload = payload
		stored = append(stored, row)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate ambiguous Run Event batch: %w", err)
	}
	return stored, nil
}

// storedJournalEvent is a projection of one stored run_events row used to
// compare against the intended event field-for-field.
type storedJournalEvent struct {
	cursor      int64
	batch       int
	source      string
	kind        string
	reviewIssue string
	toolID      string
	toolState   string
	summary     string
	createdAt   string
	payload     string
}

// storedEventsMatch compares an intended event slice against stored rows,
// including the raw payload bytes, in cursor order.
func storedEventsMatch(intended []runevent.RunEvent, stored []storedJournalEvent) bool {
	if len(intended) != len(stored) {
		return false
	}
	for index, want := range intended {
		got := stored[index]
		if got.batch != want.Batch ||
			got.source != string(want.Source) ||
			got.kind != string(want.Kind) ||
			got.reviewIssue != want.ReviewIssue ||
			got.toolID != want.ToolID ||
			got.toolState != want.ToolState ||
			got.summary != want.Summary ||
			got.createdAt != formatTime(want.Time) ||
			got.payload != string(want.Payload) {
			return false
		}
	}
	return true
}
