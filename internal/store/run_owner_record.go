package store

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"roundfix/internal/runevent"
)

// ReleaseStaleRunOwnerRecord writes the maintainer's Run Database, clearing
// only the inspected terminal Run owner. The process now holding that PID is
// unrelated: this operation never enumerates its tree or sends it a signal.
func (store *Store) ReleaseStaleRunOwnerRecord(ctx context.Context, runID string, pid int, identity string, reason string) (bool, error) {
	runID = strings.TrimSpace(runID)
	identity = strings.TrimSpace(identity)
	reason = strings.TrimSpace(reason)
	if runID == "" {
		return false, errors.New("release stale Run owner record: Run ID is required")
	}
	if pid < 1 {
		return false, errors.New("release stale Run owner record: PID is required")
	}
	if identity == "" {
		return false, errors.New("release stale Run owner record: process identity is required")
	}
	if reason == "" {
		return false, errors.New("release stale Run owner record: reason is required")
	}
	released := false
	err := store.withWriteTx(ctx, "stale Run owner record release", func(tx *sql.Tx) error {
		run, err := selectRun(ctx, tx, runID)
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("release stale Run owner record of Run %q: Run does not exist", runID)
		}
		if err != nil {
			return fmt.Errorf("read Run %q before stale Run owner record release: %w", runID, err)
		}
		if !IsTerminalState(run.State) {
			return fmt.Errorf("release stale Run owner record of Run %q: Run is %s, not terminal", runID, run.State)
		}
		result, err := tx.ExecContext(ctx, `UPDATE runs SET owner_pid = NULL, owner_identity = '', owner_identity_unproven = 0 WHERE id = ? AND owner_pid = ? AND owner_identity = ?`, runID, pid, identity)
		if err != nil {
			return fmt.Errorf("clear stale Run owner record of Run %q: %w", runID, err)
		}
		affected, err := result.RowsAffected()
		if err != nil {
			return fmt.Errorf("read stale Run owner record release result: %w", err)
		}
		if affected == 0 {
			return nil
		}
		payload, err := json.Marshal(map[string]any{"event": "stale_run_owner_record_release", "owner_pid": pid, "owner_identity": identity, "reason": reason})
		if err != nil {
			return fmt.Errorf("encode stale Run owner record release: %w", err)
		}
		if _, err := appendRunEvent(ctx, tx, runevent.RunEvent{
			RunID:   runID,
			Source:  runevent.SourceDaemon,
			Kind:    runevent.KindDaemonOutcome,
			Summary: "Run released its Stale Run Owner Record: " + reason + ".",
			Time:    store.now(),
			Payload: payload,
		}); err != nil {
			return fmt.Errorf("journal stale Run owner record release: %w", err)
		}
		released = true
		return nil
	})
	return released && err == nil, err
}
