package store

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"roundfix/internal/runevent"
)

// RunEventsOfKinds returns one Run's selected events in cursor order, with
// payloads. The selection happens in SQL and writes nothing.
func (store *Store) RunEventsOfKinds(ctx context.Context, runID string, kinds ...runevent.Kind) ([]JournalEvent, error) {
	runID = strings.TrimSpace(runID)
	if runID == "" {
		return nil, errors.New("list Run Events of kinds: Run ID is required")
	}
	events := []JournalEvent{}
	if len(kinds) == 0 {
		return events, nil
	}
	args := []any{runID}
	for _, kind := range kinds {
		args = append(args, string(kind))
	}
	placeholders := strings.TrimSuffix(strings.Repeat("?,", len(kinds)), ",")
	rows, err := store.db.QueryContext(ctx, `
SELECT cursor, batch, source, kind, review_issue, tool_id, tool_state, summary, created_at, payload
FROM run_events WHERE run_id = ? AND kind IN (`+placeholders+`) ORDER BY cursor ASC`, args...)
	if err != nil {
		return nil, fmt.Errorf("list selected Run Events for Run %q: %w", runID, err)
	}
	defer func() { _ = rows.Close() }()
	for rows.Next() {
		entry, err := scanJournalEvent(rows, runID)
		if err != nil {
			return nil, err
		}
		events = append(events, entry)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate selected Run Events for Run %q: %w", runID, err)
	}
	return events, nil
}
