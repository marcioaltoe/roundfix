package store

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"roundfix/internal/runevent"
)

// JournalEvent is one persisted Run Event with its per-Run replay cursor.
// The cursor is opaque to callers: compare it by ordering and use it as the
// next replay position.
type JournalEvent struct {
	Cursor int64
	Event  runevent.RunEvent
}

// AppendRunEvent persists one Run Event, allocating the next per-Run cursor
// inside the insert transaction. Appending to a missing Run fails clearly.
func (store *Store) AppendRunEvent(ctx context.Context, event runevent.RunEvent) (int64, error) {
	cursors, err := store.AppendRunEvents(ctx, []runevent.RunEvent{event})
	if err != nil {
		return 0, err
	}
	return cursors[0], nil
}

// AppendRunEvents persists a small ordered batch of Run Events in one
// immediate transaction and returns the allocated cursors in input order.
func (store *Store) AppendRunEvents(ctx context.Context, events []runevent.RunEvent) ([]int64, error) {
	if len(events) == 0 {
		return nil, nil
	}
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("begin Run Event append: %w", err)
	}
	defer rollbackUnlessCommitted(tx)

	cursors := make([]int64, 0, len(events))
	for _, event := range events {
		cursor, err := appendRunEvent(ctx, tx, event)
		if err != nil {
			return nil, err
		}
		cursors = append(cursors, cursor)
	}
	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("commit Run Event append: %w", err)
	}
	return cursors, nil
}

func appendRunEvent(ctx context.Context, tx *sql.Tx, event runevent.RunEvent) (int64, error) {
	runID := strings.TrimSpace(event.RunID)
	if runID == "" {
		return 0, errors.New("append Run Event: Run ID is required")
	}
	var exists int
	err := tx.QueryRowContext(ctx, `SELECT 1 FROM runs WHERE id = ?`, runID).Scan(&exists)
	if errors.Is(err, sql.ErrNoRows) {
		return 0, fmt.Errorf("append Run Event: Run %q does not exist", runID)
	}
	if err != nil {
		return 0, fmt.Errorf("append Run Event: check Run %q: %w", runID, err)
	}

	// MAX+1 inside the insert transaction is race-free under the
	// single-writer rule and never reuses a cursor (no pruning yet).
	row := tx.QueryRowContext(ctx, `
INSERT INTO run_events (
	run_id, cursor, batch, source, kind, review_issue,
	tool_id, tool_state, summary, created_at, payload
)
VALUES (?, (SELECT COALESCE(MAX(cursor), 0) + 1 FROM run_events WHERE run_id = ?), ?, ?, ?, ?, ?, ?, ?, ?, ?)
RETURNING cursor`,
		runID,
		runID,
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
	var cursor int64
	if err := row.Scan(&cursor); err != nil {
		return 0, fmt.Errorf("insert Run Event for Run %q: %w", runID, err)
	}
	return cursor, nil
}

// RunEventsAfter lists events for one Run with cursors strictly greater
// than the given cursor, oldest first, bounded by limit.
func (store *Store) RunEventsAfter(ctx context.Context, runID string, cursor int64, limit int) ([]JournalEvent, error) {
	runID = strings.TrimSpace(runID)
	if runID == "" {
		return nil, errors.New("list Run Events: Run ID is required")
	}
	if limit <= 0 {
		return nil, errors.New("list Run Events: a positive limit is required")
	}
	rows, err := store.db.QueryContext(ctx, `
SELECT cursor, batch, source, kind, review_issue, tool_id, tool_state, summary, created_at, payload
FROM run_events
WHERE run_id = ? AND cursor > ?
ORDER BY cursor ASC
LIMIT ?`,
		runID,
		cursor,
		limit,
	)
	if err != nil {
		return nil, fmt.Errorf("list Run Events for Run %q: %w", runID, err)
	}
	defer func() {
		_ = rows.Close()
	}()

	events := []JournalEvent{}
	for rows.Next() {
		entry, err := scanJournalEvent(rows, runID)
		if err != nil {
			return nil, err
		}
		events = append(events, entry)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate Run Events for Run %q: %w", runID, err)
	}
	return events, nil
}

func scanJournalEvent(rows *sql.Rows, runID string) (JournalEvent, error) {
	var entry JournalEvent
	var source string
	var kind string
	var createdAt string
	var payload string
	if err := rows.Scan(
		&entry.Cursor,
		&entry.Event.Batch,
		&source,
		&kind,
		&entry.Event.ReviewIssue,
		&entry.Event.ToolID,
		&entry.Event.ToolState,
		&entry.Event.Summary,
		&createdAt,
		&payload,
	); err != nil {
		return JournalEvent{}, fmt.Errorf("scan Run Event for Run %q: %w", runID, err)
	}
	parsedAt, err := parseTime(createdAt)
	if err != nil {
		return JournalEvent{}, err
	}
	entry.Event.RunID = runID
	entry.Event.Source = runevent.Source(source)
	entry.Event.Kind = runevent.Kind(kind)
	entry.Event.Time = parsedAt
	if payload != "" {
		entry.Event.Payload = []byte(payload)
	}
	return entry, nil
}

// JournalSink adapts the Run Database to the Run Event sink interface. It
// registers as a critical sink: an append failure after Run start must fail
// the Run, never be swallowed.
type JournalSink struct {
	Store *Store
}

func (sink JournalSink) Publish(ctx context.Context, event runevent.RunEvent) error {
	_, err := sink.Store.AppendRunEvent(ctx, event)
	return err
}

// DataVersion exposes SQLite's data_version for this connection: it changes
// when another connection commits, so pollers can detect new events without
// reading rows.
func (store *Store) DataVersion(ctx context.Context) (int64, error) {
	var version int64
	if err := store.db.QueryRowContext(ctx, `PRAGMA data_version`).Scan(&version); err != nil {
		return 0, fmt.Errorf("read Run Database data version: %w", err)
	}
	return version, nil
}
