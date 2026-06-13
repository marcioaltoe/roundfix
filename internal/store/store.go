package store

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	_ "modernc.org/sqlite"
)

const (
	roundfixHomeDir = ".roundfix"
	databaseName    = "roundfix.db"

	KindFetch   = "fetch"
	KindResolve = "resolve"
	KindWatch   = "watch"

	StateActive           = "Active"
	StateFetched          = "Fetched"
	StateStopped          = "Stopped"
	StateClean            = "Clean"
	StateMaxRoundsReached = "MaxRoundsReached"
	StateBudgetExceeded   = "BudgetExceeded"
	StateTimedOut         = "TimedOut"
	StateFailed           = "Failed"
	// StateUnresolved means the resolve work completed but Unresolved
	// Review Issues remain, so Final Push stayed blocked. It is a deliberate
	// outcome, distinct from StateFailed which means the Run itself broke.
	StateUnresolved = "Unresolved"

	// Intermediate states reflect what an Active Run is doing during a
	// resolve cycle. They are non-terminal: CompleteRun still ends the Run.
	StateResolvingWithAgent = "ResolvingWithAgent"
	StateVerifying          = "Verifying"
	StatePushing            = "Pushing"
)

type Store struct {
	db  *sql.DB
	now func() time.Time
}

type Run struct {
	ID             string
	Kind           string
	State          string
	HeadRepository string
	HeadBranch     string
	BaseRepository string
	PRNumber       string
	GitRoot        string
	LocalBranch    string
	HeadSHA        string
	ArtifactDir    string
	CreatedAt      time.Time
	UpdatedAt      time.Time
	CompletedAt    *time.Time
}

type InteractiveDefaults struct {
	PRNumber string
	Agent    string
}

type CreateRunRequest struct {
	Kind           string
	HeadRepository string
	HeadBranch     string
	BaseRepository string
	PRNumber       string
	GitRoot        string
	LocalBranch    string
	HeadSHA        string
	ArtifactDir    string
}

type ActiveRunError struct {
	Existing Run
}

func (err ActiveRunError) Error() string {
	return fmt.Sprintf("Active Run already exists for Head Repository %q and PR Head Branch %q; existing run_id=%s state=%s", err.Existing.HeadRepository, err.Existing.HeadBranch, err.Existing.ID, err.Existing.State)
}

func DatabasePath(homeDir string) string {
	return filepath.Join(homeDir, roundfixHomeDir, databaseName)
}

// Open opens the writer connection: single connection, immediate
// transactions, WAL enabled at creation before any reader connects, and a
// busy timeout. Single-writer discipline: open at most one writer.
func Open(ctx context.Context, homeDir string) (*Store, error) {
	if strings.TrimSpace(homeDir) == "" {
		return nil, errors.New("open Run Database: home directory is required")
	}
	path := DatabasePath(homeDir)
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return nil, fmt.Errorf("create Roundfix Home %q: %w", filepath.Dir(path), err)
	}

	db, err := sql.Open("sqlite", writerDSN(path))
	if err != nil {
		return nil, fmt.Errorf("open Run Database %q: %w", path, err)
	}
	db.SetMaxOpenConns(1)

	store := &Store{
		db:  db,
		now: func() time.Time { return time.Now().UTC() },
	}
	if err := store.migrate(ctx); err != nil {
		_ = db.Close()
		return nil, err
	}
	return store, nil
}

// OpenReader opens a read-only connection for paging Run Events while the
// writer appends. It never migrates; the Run Database must already exist.
func OpenReader(ctx context.Context, homeDir string) (*Store, error) {
	if strings.TrimSpace(homeDir) == "" {
		return nil, errors.New("open Run Database reader: home directory is required")
	}
	path := DatabasePath(homeDir)
	if _, err := os.Stat(path); err != nil {
		return nil, fmt.Errorf("open Run Database reader %q: %w", path, err)
	}

	db, err := sql.Open("sqlite", readerDSN(path))
	if err != nil {
		return nil, fmt.Errorf("open Run Database reader %q: %w", path, err)
	}
	db.SetMaxOpenConns(1)
	store := &Store{
		db:  db,
		now: func() time.Time { return time.Now().UTC() },
	}
	if err := db.PingContext(ctx); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("open Run Database reader %q: %w", path, err)
	}
	return store, nil
}

func writerDSN(path string) string {
	return "file:" + path + "?_txlock=immediate" +
		"&_pragma=busy_timeout(5000)" +
		"&_pragma=journal_mode(WAL)" +
		"&_pragma=foreign_keys(1)"
}

func readerDSN(path string) string {
	return "file:" + path + "?mode=ro" +
		"&_pragma=busy_timeout(5000)" +
		"&_pragma=foreign_keys(1)"
}

func (store *Store) Close() error {
	return store.db.Close()
}

func (store *Store) CreateFetchRun(ctx context.Context, req CreateRunRequest) (Run, error) {
	req.Kind = KindFetch
	return store.CreateRun(ctx, req)
}

func (store *Store) CreateRun(ctx context.Context, req CreateRunRequest) (Run, error) {
	if err := validateCreateRunRequest(req); err != nil {
		return Run{}, err
	}
	runID, err := newRunID(store.now())
	if err != nil {
		return Run{}, err
	}
	now := store.now()

	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return Run{}, fmt.Errorf("begin Run creation: %w", err)
	}
	defer rollbackUnlessCommitted(tx)

	existing, found, err := selectActiveRun(ctx, tx, req.HeadRepository, req.HeadBranch)
	if err != nil {
		return Run{}, err
	}
	if found {
		return Run{}, ActiveRunError{Existing: existing}
	}

	_, err = tx.ExecContext(ctx, `
INSERT INTO runs (
	id, kind, state, head_repository, head_branch, base_repository,
	pr_number, git_root, local_branch, head_sha, artifact_dir,
	created_at, updated_at
) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		runID,
		req.Kind,
		StateActive,
		req.HeadRepository,
		req.HeadBranch,
		req.BaseRepository,
		req.PRNumber,
		req.GitRoot,
		req.LocalBranch,
		req.HeadSHA,
		req.ArtifactDir,
		formatTime(now),
		formatTime(now),
	)
	if err != nil {
		return Run{}, fmt.Errorf("insert Run record: %w", err)
	}

	_, err = tx.ExecContext(ctx, `
INSERT INTO active_run_locks (head_repository, head_branch, run_id, created_at)
VALUES (?, ?, ?, ?)`,
		req.HeadRepository,
		req.HeadBranch,
		runID,
		formatTime(now),
	)
	if err != nil {
		existing, found, selectErr := selectActiveRun(ctx, tx, req.HeadRepository, req.HeadBranch)
		if selectErr == nil && found {
			return Run{}, ActiveRunError{Existing: existing}
		}
		return Run{}, fmt.Errorf("acquire Active Run lock: %w", err)
	}

	run, err := selectRun(ctx, tx, runID)
	if err != nil {
		return Run{}, err
	}
	if err := tx.Commit(); err != nil {
		return Run{}, fmt.Errorf("commit Run creation: %w", err)
	}
	return run, nil
}

func (store *Store) CompleteRun(ctx context.Context, runID string, terminalState string) (Run, error) {
	if !IsTerminalState(terminalState) {
		return Run{}, fmt.Errorf("Run state %q is not terminal", terminalState)
	}
	now := store.now()

	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return Run{}, fmt.Errorf("begin Run completion: %w", err)
	}
	defer rollbackUnlessCommitted(tx)

	result, err := tx.ExecContext(ctx, `
UPDATE runs
SET state = ?, updated_at = ?, completed_at = ?
WHERE id = ?`,
		terminalState,
		formatTime(now),
		formatTime(now),
		runID,
	)
	if err != nil {
		return Run{}, fmt.Errorf("update Run terminal state: %w", err)
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return Run{}, fmt.Errorf("read Run completion result: %w", err)
	}
	if affected == 0 {
		return Run{}, fmt.Errorf("Run %q does not exist", runID)
	}

	if _, err := tx.ExecContext(ctx, `DELETE FROM active_run_locks WHERE run_id = ?`, runID); err != nil {
		return Run{}, fmt.Errorf("release Active Run lock: %w", err)
	}

	run, err := selectRun(ctx, tx, runID)
	if err != nil {
		return Run{}, err
	}
	if err := tx.Commit(); err != nil {
		return Run{}, fmt.Errorf("commit Run completion: %w", err)
	}
	return run, nil
}

// UpdateRunState records an intermediate, non-terminal state for an Active
// Run. Terminal outcomes must go through CompleteRun.
func (store *Store) UpdateRunState(ctx context.Context, runID string, state string) error {
	if strings.TrimSpace(state) == "" {
		return errors.New("update Run state: state is required")
	}
	if IsTerminalState(state) {
		return fmt.Errorf("update Run state: %q is terminal; use CompleteRun", state)
	}
	result, err := store.db.ExecContext(ctx, `
UPDATE runs SET state = ?, updated_at = ? WHERE id = ?`,
		state,
		formatTime(store.now()),
		runID,
	)
	if err != nil {
		return fmt.Errorf("update Run state: %w", err)
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("read Run state update result: %w", err)
	}
	if affected == 0 {
		return fmt.Errorf("update Run state: Run %q does not exist", runID)
	}
	return nil
}

func (store *Store) ActiveRun(ctx context.Context, headRepository string, headBranch string) (Run, bool, error) {
	return selectActiveRun(ctx, store.db, headRepository, headBranch)
}

func (store *Store) Run(ctx context.Context, runID string) (Run, bool, error) {
	runID = strings.TrimSpace(runID)
	if runID == "" {
		return Run{}, false, errors.New("Run ID is required")
	}
	run, err := selectRun(ctx, store.db, runID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Run{}, false, nil
		}
		return Run{}, false, err
	}
	return run, true, nil
}

func (store *Store) RunCount(ctx context.Context) (int, error) {
	var count int
	if err := store.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM runs`).Scan(&count); err != nil {
		return 0, fmt.Errorf("count Runs: %w", err)
	}
	return count, nil
}

func (store *Store) MigrationVersion(ctx context.Context) (int, error) {
	var version int
	if err := store.db.QueryRowContext(ctx, `PRAGMA user_version`).Scan(&version); err != nil {
		return 0, fmt.Errorf("read Run Database migration version: %w", err)
	}
	return version, nil
}

func (store *Store) InteractiveDefaults(ctx context.Context) (InteractiveDefaults, error) {
	rows, err := store.db.QueryContext(ctx, `SELECT key, value FROM interactive_defaults`)
	if err != nil {
		return InteractiveDefaults{}, fmt.Errorf("read Interactive Input defaults: %w", err)
	}
	defer func() {
		_ = rows.Close()
	}()

	var defaults InteractiveDefaults
	for rows.Next() {
		var key string
		var value string
		if err := rows.Scan(&key, &value); err != nil {
			return InteractiveDefaults{}, fmt.Errorf("scan Interactive Input default: %w", err)
		}
		switch key {
		case "pr_number":
			defaults.PRNumber = value
		case "agent":
			defaults.Agent = value
		}
	}
	if err := rows.Err(); err != nil {
		return InteractiveDefaults{}, fmt.Errorf("iterate Interactive Input defaults: %w", err)
	}
	return defaults, nil
}

func (store *Store) RememberInteractiveDefaults(ctx context.Context, defaults InteractiveDefaults) error {
	now := formatTime(store.now())
	updates := []struct {
		key   string
		value string
	}{
		{key: "pr_number", value: strings.TrimSpace(defaults.PRNumber)},
		{key: "agent", value: strings.TrimSpace(defaults.Agent)},
	}
	for _, update := range updates {
		if update.value == "" {
			continue
		}
		if _, err := store.db.ExecContext(ctx, `
INSERT INTO interactive_defaults (key, value, updated_at)
VALUES (?, ?, ?)
ON CONFLICT(key) DO UPDATE SET value = excluded.value, updated_at = excluded.updated_at`,
			update.key,
			update.value,
			now,
		); err != nil {
			return fmt.Errorf("remember Interactive Input default %s: %w", update.key, err)
		}
	}
	return nil
}

func IsTerminalState(state string) bool {
	switch state {
	case StateFetched, StateStopped, StateClean, StateMaxRoundsReached, StateBudgetExceeded, StateTimedOut, StateFailed, StateUnresolved:
		return true
	default:
		return false
	}
}

func (store *Store) migrate(ctx context.Context) error {
	statements := []string{
		`PRAGMA foreign_keys = ON`,
		`CREATE TABLE IF NOT EXISTS runs (
			id TEXT PRIMARY KEY,
			kind TEXT NOT NULL,
			state TEXT NOT NULL,
			head_repository TEXT NOT NULL,
			head_branch TEXT NOT NULL,
			base_repository TEXT NOT NULL DEFAULT '',
			pr_number TEXT NOT NULL,
			git_root TEXT NOT NULL,
			local_branch TEXT NOT NULL,
			head_sha TEXT NOT NULL,
			artifact_dir TEXT NOT NULL,
			created_at TEXT NOT NULL,
			updated_at TEXT NOT NULL,
			completed_at TEXT NOT NULL DEFAULT ''
		)`,
		`CREATE TABLE IF NOT EXISTS active_run_locks (
			head_repository TEXT NOT NULL,
			head_branch TEXT NOT NULL,
			run_id TEXT NOT NULL UNIQUE,
			created_at TEXT NOT NULL,
			PRIMARY KEY (head_repository, head_branch),
			FOREIGN KEY (run_id) REFERENCES runs(id) ON DELETE CASCADE
		)`,
		`CREATE INDEX IF NOT EXISTS idx_runs_head ON runs (head_repository, head_branch)`,
		`CREATE TABLE IF NOT EXISTS interactive_defaults (
			key TEXT PRIMARY KEY,
			value TEXT NOT NULL,
			updated_at TEXT NOT NULL
		)`,
		`CREATE TABLE IF NOT EXISTS run_events (
			run_id TEXT NOT NULL,
			cursor INTEGER NOT NULL,
			batch INTEGER NOT NULL DEFAULT 0,
			source TEXT NOT NULL,
			kind TEXT NOT NULL,
			review_issue TEXT NOT NULL DEFAULT '',
			tool_id TEXT NOT NULL DEFAULT '',
			tool_state TEXT NOT NULL DEFAULT '',
			summary TEXT NOT NULL DEFAULT '',
			created_at TEXT NOT NULL,
			payload TEXT NOT NULL DEFAULT '',
			PRIMARY KEY (run_id, cursor),
			FOREIGN KEY (run_id) REFERENCES runs(id) ON DELETE CASCADE
		)`,
		`PRAGMA user_version = 3`,
	}
	for _, statement := range statements {
		if _, err := store.db.ExecContext(ctx, statement); err != nil {
			return fmt.Errorf("apply Run Database migration: %w", err)
		}
	}
	return nil
}

func validateCreateRunRequest(req CreateRunRequest) error {
	if req.Kind == "" {
		return errors.New("Run kind is required")
	}
	switch req.Kind {
	case KindFetch, KindResolve, KindWatch:
	default:
		return fmt.Errorf("Run kind %q is invalid", req.Kind)
	}
	required := map[string]string{
		"Head Repository":    req.HeadRepository,
		"PR Head Branch":     req.HeadBranch,
		"pull request":       req.PRNumber,
		"Git root":           req.GitRoot,
		"local branch":       req.LocalBranch,
		"HEAD":               req.HeadSHA,
		"Artifact Directory": req.ArtifactDir,
	}
	for label, value := range required {
		if strings.TrimSpace(value) == "" {
			return fmt.Errorf("%s is required to create a Run", label)
		}
	}
	return nil
}

type runQuerier interface {
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
}

func selectActiveRun(ctx context.Context, querier runQuerier, headRepository string, headBranch string) (Run, bool, error) {
	row := querier.QueryRowContext(ctx, `
SELECT r.id, r.kind, r.state, r.head_repository, r.head_branch, r.base_repository,
       r.pr_number, r.git_root, r.local_branch, r.head_sha, r.artifact_dir,
       r.created_at, r.updated_at, r.completed_at
FROM active_run_locks l
JOIN runs r ON r.id = l.run_id
WHERE l.head_repository = ? AND l.head_branch = ?`,
		headRepository,
		headBranch,
	)
	run, err := scanRun(row)
	if errors.Is(err, sql.ErrNoRows) {
		return Run{}, false, nil
	}
	if err != nil {
		return Run{}, false, err
	}
	return run, true, nil
}

func selectRun(ctx context.Context, querier runQuerier, runID string) (Run, error) {
	row := querier.QueryRowContext(ctx, `
SELECT id, kind, state, head_repository, head_branch, base_repository,
       pr_number, git_root, local_branch, head_sha, artifact_dir,
       created_at, updated_at, completed_at
FROM runs
WHERE id = ?`, runID)
	run, err := scanRun(row)
	if err != nil {
		return Run{}, fmt.Errorf("select Run %q: %w", runID, err)
	}
	return run, nil
}

func scanRun(row *sql.Row) (Run, error) {
	var run Run
	var createdAt string
	var updatedAt string
	var completedAt string
	err := row.Scan(
		&run.ID,
		&run.Kind,
		&run.State,
		&run.HeadRepository,
		&run.HeadBranch,
		&run.BaseRepository,
		&run.PRNumber,
		&run.GitRoot,
		&run.LocalBranch,
		&run.HeadSHA,
		&run.ArtifactDir,
		&createdAt,
		&updatedAt,
		&completedAt,
	)
	if err != nil {
		return Run{}, err
	}
	parsedCreatedAt, err := parseTime(createdAt)
	if err != nil {
		return Run{}, err
	}
	parsedUpdatedAt, err := parseTime(updatedAt)
	if err != nil {
		return Run{}, err
	}
	run.CreatedAt = parsedCreatedAt
	run.UpdatedAt = parsedUpdatedAt
	if completedAt != "" {
		parsedCompletedAt, err := parseTime(completedAt)
		if err != nil {
			return Run{}, err
		}
		run.CompletedAt = &parsedCompletedAt
	}
	return run, nil
}

func rollbackUnlessCommitted(tx *sql.Tx) {
	_ = tx.Rollback()
}

func newRunID(now time.Time) (string, error) {
	var random [8]byte
	if _, err := rand.Read(random[:]); err != nil {
		return "", fmt.Errorf("generate Run id: %w", err)
	}
	return fmt.Sprintf("run_%s_%s", now.UTC().Format("20060102T150405Z"), hex.EncodeToString(random[:])), nil
}

func formatTime(value time.Time) string {
	return value.UTC().Format(time.RFC3339Nano)
}

func parseTime(value string) (time.Time, error) {
	parsed, err := time.Parse(time.RFC3339Nano, value)
	if err != nil {
		return time.Time{}, fmt.Errorf("parse Run timestamp %q: %w", value, err)
	}
	return parsed, nil
}
