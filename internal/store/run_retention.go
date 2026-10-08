package store

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"
)

const IncrementalCompactionSlicePages = 2560

const runRetentionSweepsSchema = `CREATE TABLE IF NOT EXISTS run_retention_sweeps (
 id INTEGER PRIMARY KEY CHECK (id = 1),
 completed_at TEXT NOT NULL,
 retention_days INTEGER NOT NULL
)`

type RunRetentionRows struct{ Runs, RunEvents, AgentSelections, TokenUsage, ActiveRunLocks int64 }
type RunRetentionCandidate struct {
	RunID, Repository, GitRoot, ArtifactDir, WorkDir string
	CompletedAt                                      time.Time
	QueueReferenced                                  bool
	Rows                                             RunRetentionRows
	EstimatedBytes                                   int64
}
type RunRetentionKeptError struct{ RunID, Reason string }

func (err RunRetentionKeptError) Error() string {
	return fmt.Sprintf("Run Retention kept Run %q: %s", err.RunID, err.Reason)
}

type RunRetentionSweep struct {
	CompletedAt   time.Time
	RetentionDays int
}
type IncrementalCompaction struct {
	Incremental              bool
	PagesReleased, FreePages int64
}

const retentionQueueReferenceSQL = `EXISTS (SELECT 1 FROM delivery_queue_items WHERE run_id = r.id)
 OR EXISTS (SELECT 1 FROM delivery_queue_runs WHERE run_id = r.id)`
const retentionRowCountsSQL = `(SELECT COUNT(*) FROM run_events WHERE run_id = r.id),
 (SELECT COUNT(*) FROM run_agent_selections WHERE run_id = r.id),
 (SELECT COUNT(*) FROM run_token_usage WHERE run_id = r.id),
 (SELECT COUNT(*) FROM active_run_locks WHERE run_id = r.id)`

// RunRetentionCandidates reads only the terminal Runs at the cutoff, including
// queue references so callers can report why a candidate must be kept.
func (store *Store) RunRetentionCandidates(ctx context.Context, cutoff time.Time) ([]RunRetentionCandidate, error) {
	placeholders := strings.TrimRight(strings.Repeat("?,", len(terminalStates)), ",")
	args := make([]any, 0, len(terminalStates)+1)
	for _, state := range terminalStates {
		args = append(args, state)
	}
	args = append(args, formatTime(cutoff))
	rows, err := store.db.QueryContext(ctx, `SELECT r.id,
 COALESCE(NULLIF(r.repository_root, ''), r.git_root), r.git_root, r.artifact_dir,
 COALESCE(r.work_dir, ''), r.completed_at, (`+retentionQueueReferenceSQL+`),
 `+retentionRowCountsSQL+`,
 (SELECT COALESCE(SUM(length(CAST(payload AS BLOB)) + length(CAST(summary AS BLOB))), 0)
 FROM run_events WHERE run_id = r.id)
 FROM runs r WHERE r.state IN (`+placeholders+`)
 AND TRIM(r.completed_at) <> '' AND julianday(r.completed_at) <= julianday(?)
 ORDER BY r.completed_at, r.id`, args...)
	if err != nil {
		return nil, fmt.Errorf("select Run Retention candidates: %w", err)
	}
	defer func() { _ = rows.Close() }()
	candidates := []RunRetentionCandidate{}
	for rows.Next() {
		var c RunRetentionCandidate
		var completed string
		if err := rows.Scan(&c.RunID, &c.Repository, &c.GitRoot, &c.ArtifactDir, &c.WorkDir,
			&completed, &c.QueueReferenced, &c.Rows.RunEvents, &c.Rows.AgentSelections,
			&c.Rows.TokenUsage, &c.Rows.ActiveRunLocks, &c.EstimatedBytes); err != nil {
			return nil, fmt.Errorf("scan Run Retention candidate: %w", err)
		}
		c.CompletedAt, err = parseTime(completed)
		if err != nil {
			return nil, fmt.Errorf("read Run %q retention completion: %w", c.RunID, err)
		}
		// julianday rounds fractional seconds; Go preserves the strict boundary.
		if c.CompletedAt.Before(cutoff) {
			c.Rows.Runs = 1
			candidates = append(candidates, c)
		}
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate Run Retention candidates: %w", err)
	}
	return candidates, nil
}

// RemoveRetainedRun rechecks eligibility and counts the cascade in the same
// write transaction. Repeated removals of an absent Run reclaim nothing.
// RemoveRetainedRun deletes one terminal Run past the cutoff in a write
// transaction that first re-reads it. Each artifact step runs inside that
// transaction after the recheck and before the row is deleted, so the Run's
// artifacts go first and a state change or queue link made meanwhile waits for
// the transaction; a failed step rolls the removal back.
func (store *Store) RemoveRetainedRun(ctx context.Context, runID string, cutoff time.Time, artifactSteps ...func() error) (RunRetentionRows, error) {
	var removed RunRetentionRows
	err := store.withWriteTx(ctx, "Run Retention removal", func(tx *sql.Tx) error {
		var state, completed string
		var queueReferenced bool
		err := tx.QueryRowContext(ctx, `SELECT r.state, r.completed_at, (`+retentionQueueReferenceSQL+`)
  FROM runs r WHERE r.id = ?`, runID).Scan(&state, &completed, &queueReferenced)
		if errors.Is(err, sql.ErrNoRows) {
			return nil
		}
		if err != nil {
			return fmt.Errorf("recheck retained Run %q: %w", runID, err)
		}
		if !IsTerminalState(state) {
			return RunRetentionKeptError{runID, "Run is no longer terminal"}
		}
		if strings.TrimSpace(completed) == "" {
			return RunRetentionKeptError{runID, "Run has no completion time"}
		}
		completedAt, err := parseTime(completed)
		if err != nil {
			return fmt.Errorf("read retained Run %q completion: %w", runID, err)
		}
		if !completedAt.Before(cutoff) {
			return RunRetentionKeptError{runID, "Run is no longer before the cutoff"}
		}
		if queueReferenced {
			return RunRetentionKeptError{runID, "Run is queue-referenced"}
		}
		for _, step := range artifactSteps {
			if err := step(); err != nil {
				return err
			}
		}
		if err := tx.QueryRowContext(ctx, `SELECT `+retentionRowCountsSQL+` FROM runs r WHERE r.id = ?`, runID).
			Scan(&removed.RunEvents, &removed.AgentSelections, &removed.TokenUsage, &removed.ActiveRunLocks); err != nil {
			return fmt.Errorf("count retained Run %q rows: %w", runID, err)
		}
		if _, err := tx.ExecContext(ctx, `DELETE FROM runs WHERE id = ?`, runID); err != nil {
			return fmt.Errorf("remove retained Run %q: %w", runID, err)
		}
		removed.Runs = 1
		return nil
	})
	if err != nil {
		return RunRetentionRows{}, err
	}
	return removed, nil
}

func (store *Store) LastRunRetentionSweep(ctx context.Context) (RunRetentionSweep, bool, error) {
	var sweep RunRetentionSweep
	var completed string
	err := store.db.QueryRowContext(ctx, `SELECT completed_at, retention_days FROM run_retention_sweeps WHERE id = 1`).Scan(&completed, &sweep.RetentionDays)
	if errors.Is(err, sql.ErrNoRows) {
		return sweep, false, nil
	}
	if err != nil {
		return sweep, false, fmt.Errorf("read last Run Retention sweep: %w", err)
	}
	sweep.CompletedAt, err = parseTime(completed)
	if err != nil {
		return RunRetentionSweep{}, false, fmt.Errorf("read Run Retention sweep completion: %w", err)
	}
	return sweep, true, nil
}

func (store *Store) RecordRunRetentionSweep(ctx context.Context, sweep RunRetentionSweep) error {
	return store.withWriteTx(ctx, "Run Retention sweep record", func(tx *sql.Tx) error {
		if _, err := tx.ExecContext(ctx, `INSERT INTO run_retention_sweeps (id, completed_at, retention_days)
  VALUES (1, ?, ?) ON CONFLICT(id) DO UPDATE SET completed_at = excluded.completed_at,
  retention_days = excluded.retention_days`, formatTime(sweep.CompletedAt), sweep.RetentionDays); err != nil {
			return fmt.Errorf("record Run Retention sweep: %w", err)
		}
		return nil
	})
}

// CompactIncrementally releases at most maxPages, holding the machine-wide
// write lock for this slice alone. Default-mode databases require Compact.
func (store *Store) CompactIncrementally(ctx context.Context, maxPages int64) (IncrementalCompaction, error) {
	if maxPages <= 0 {
		return IncrementalCompaction{}, errors.New("compact incrementally: positive maxPages is required")
	}
	var result IncrementalCompaction
	err := store.withWriteTx(ctx, "incremental compaction", func(tx *sql.Tx) error {
		mode, err := storagePragmaInt64(ctx, tx, "auto_vacuum")
		if err != nil {
			return err
		}
		if mode != 2 {
			return nil
		}
		result.Incremental = true
		before, err := storagePragmaInt64(ctx, tx, "page_count")
		if err != nil {
			return err
		}
		if _, err := tx.ExecContext(ctx, fmt.Sprintf("PRAGMA incremental_vacuum(%d)", maxPages)); err != nil {
			return fmt.Errorf("compact Run Database incrementally: %w", err)
		}
		after, err := storagePragmaInt64(ctx, tx, "page_count")
		if err != nil {
			return err
		}
		result.PagesReleased = before - after
		result.FreePages, err = storagePragmaInt64(ctx, tx, "freelist_count")
		return err
	})
	if err != nil {
		return IncrementalCompaction{}, err
	}
	return result, nil
}

// RunRetentionStorage reads logical database bytes, including pages still in WAL.
func (store *Store) RunRetentionStorage(ctx context.Context) (bytes, freePages, mode int64, err error) {
	pages, err := storagePragmaInt64(ctx, store.db, "page_count")
	if err != nil {
		return 0, 0, 0, err
	}
	size, err := storagePragmaInt64(ctx, store.db, "page_size")
	if err != nil {
		return 0, 0, 0, err
	}
	freePages, err = storagePragmaInt64(ctx, store.db, "freelist_count")
	if err != nil {
		return 0, 0, 0, err
	}
	mode, err = storagePragmaInt64(ctx, store.db, "auto_vacuum")
	return pages * size, freePages, mode, err
}

// CheckpointRunRetention returns incremental compaction's WAL pages without
// waiting for readers to release their snapshots.
func (store *Store) CheckpointRunRetention(ctx context.Context) error {
	if _, err := store.db.ExecContext(ctx, "PRAGMA wal_checkpoint(PASSIVE)"); err != nil {
		return fmt.Errorf("checkpoint Run Retention: %w", err)
	}
	return nil
}
