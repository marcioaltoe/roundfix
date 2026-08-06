package store

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"roundfix/internal/runevent"

	_ "modernc.org/sqlite"
)

const (
	roundfixHomeDir = ".roundfix"
	databaseName    = "roundfix.db"

	KindFetch     = "fetch"
	KindResolve   = "resolve"
	KindWatch     = "watch"
	KindImplement = "implement"

	// RunBranchPrefix is the deterministic namespace for spec Run Branches.
	RunBranchPrefix = "roundfix/run-"

	// Active Run locks are keyed by work target (ADR 0016): review Kinds
	// lock the Open Pull Request, the implement Kind locks the Spec.
	targetKindPR   = "pr"
	targetKindSpec = "spec"

	StateActive             = "Active"
	StateFetched            = "Fetched"
	StateStopped            = "Stopped"
	StateClean              = "Clean"
	StateCleanUnverified    = "CleanUnverified"
	StateReviewSkipped      = "ReviewSkipped"
	StateMaxRoundsReached   = "MaxRoundsReached"
	StateBudgetExceeded     = "BudgetExceeded"
	StateTimedOut           = "TimedOut"
	StateFailed             = "Failed"
	StateCheckoutMoved      = "CheckoutMoved"
	StateIntegrationPending = "IntegrationPending"
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

var terminalStates = []string{
	StateFetched,
	StateStopped,
	StateClean,
	StateCleanUnverified,
	StateReviewSkipped,
	StateMaxRoundsReached,
	StateBudgetExceeded,
	StateTimedOut,
	StateFailed,
	StateCheckoutMoved,
	StateIntegrationPending,
	StateUnresolved,
}

type Store struct {
	db  *sql.DB
	now func() time.Time
}

type Run struct {
	ID                    string
	Kind                  string
	State                 string
	HeadRepository        string
	HeadBranch            string
	BaseRepository        string
	PRNumber              string
	GitRoot               string
	LocalBranch           string
	HeadSHA               string
	ArtifactDir           string
	WorkDir               string
	SpecSlug              string
	Agent                 string
	Model                 string
	ReasoningEffort       string
	OwnerPID              *int
	OwnerIdentity         string
	OwnerIdentityUnproven bool
	CreatedAt             time.Time
	UpdatedAt             time.Time
	CompletedAt           *time.Time
}

type TerminalOutcomeConflictError struct {
	RunID     string
	Stored    string
	Requested string
}

func (err TerminalOutcomeConflictError) Error() string {
	return fmt.Sprintf(
		"terminal outcome conflict for Run %q: stored %q, requested %q",
		err.RunID,
		err.Stored,
		err.Requested,
	)
}

type CompleteRunResult struct {
	Run
	Transitioned bool
}

type IntegrationReconciliation struct {
	RunID           string
	PreviousOutcome string
	Classification  string
	RunBranch       string
	RunHead         string
	TargetBranch    string
	TargetHead      string
	Worktree        string
	Reason          string
	Action          string
	Time            time.Time
}

type InteractiveDefaults struct {
	PRNumber string
	Agent    string
}

type CreateRunRequest struct {
	Kind            string
	HeadRepository  string
	HeadBranch      string
	BaseRepository  string
	PRNumber        string
	GitRoot         string
	LocalBranch     string
	HeadSHA         string
	ArtifactDir     string
	WorkDir         string
	SpecSlug        string
	Agent           string
	Model           string
	ReasoningEffort string
	OwnerPID        int
	OwnerIdentity   string
}

// RunStateFilter selects which Run states a listing includes. The zero
// value is StatesActive, so an unset filter lists only non-terminal Runs.
type RunStateFilter int

const (
	// StatesActive keeps only non-terminal Runs.
	StatesActive RunStateFilter = iota
	// StatesTerminal keeps only terminal Runs.
	StatesTerminal
	// StatesAll keeps every Run.
	StatesAll
)

func (filter RunStateFilter) matches(state string) bool {
	switch filter {
	case StatesTerminal:
		return IsTerminalState(state)
	case StatesAll:
		return true
	default:
		return !IsTerminalState(state)
	}
}

// ListRunsQuery scopes a Run listing.
type ListRunsQuery struct {
	GitRoot string
	States  RunStateFilter
	// Limit bounds the listing to the newest matching Runs; 0 is unbounded.
	Limit int
}

type ActiveRunError struct {
	Existing Run
}

// SchemaVersionError reports a Run Database whose schema version differs from
// the one this binary supports. Read-only surfaces never migrate; the guard
// only converts silent SQL errors on unmigrated databases into this
// deterministic diagnostic.
type SchemaVersionError struct {
	Path      string
	Found     int
	Supported int
}

func (err SchemaVersionError) Error() string {
	return fmt.Sprintf(
		"Run Database %q has schema version %d, but this binary supports schema version %d; run one operational roundfix command (resolve, watch, or implement) to migrate the Run Database",
		err.Path,
		err.Found,
		err.Supported,
	)
}

var ErrTerminalRunStopRequest = errors.New("cannot record Stop Request for terminal Run")

func (err ActiveRunError) Error() string {
	if err.Existing.Kind == KindImplement {
		return fmt.Sprintf("Active Run already exists for repository %q and Spec %q; existing run_id=%s state=%s; stop it with: roundfix stop %s", err.Existing.GitRoot, err.Existing.SpecSlug, err.Existing.ID, err.Existing.State, err.Existing.ID)
	}
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
// writer appends. It never migrates; the Run Database must already exist and
// carry the supported schema version, otherwise a typed SchemaVersionError
// names the migration remediation instead of failing on a later query.
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
	version, err := store.MigrationVersion(ctx)
	if err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("open Run Database reader %q: %w", path, err)
	}
	if version != schemaVersion {
		_ = db.Close()
		return nil, SchemaVersionError{Path: path, Found: version, Supported: schemaVersion}
	}
	return store, nil
}

// busyTimeoutMillis bounds how long SQLite waits for a held lock before
// returning SQLITE_BUSY. The Run Database is machine-wide, so concurrent Runs,
// a following `roundfix events` stream, and a `runs list` all contend on it. At
// 5s, an Agent Batch died twice on 2026-08-05 with
// `publish Run Events: begin Run Event append: database is locked (5)` — an
// infrastructure failure with no work defect, which failed the Task and cost a
// recovery Run each time.
const busyTimeoutMillis = 30000

func writerDSN(path string) string {
	return "file:" + path + "?_txlock=immediate" +
		fmt.Sprintf("&_pragma=busy_timeout(%d)", busyTimeoutMillis) +
		"&_pragma=journal_mode(WAL)" +
		"&_pragma=foreign_keys(1)"
}

func readerDSN(path string) string {
	return "file:" + path + "?mode=ro" +
		fmt.Sprintf("&_pragma=busy_timeout(%d)", busyTimeoutMillis) +
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
	return store.createRun(ctx, req, true)
}

// CreateRunSkippingActiveLock creates a Run without acquiring the target lock.
func (store *Store) CreateRunSkippingActiveLock(ctx context.Context, req CreateRunRequest) (Run, error) {
	return store.createRun(ctx, req, false)
}

func (store *Store) createRun(ctx context.Context, req CreateRunRequest, acquireActiveLock bool) (Run, error) {
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

	targetKind, targetKey := lockTarget(req)
	if acquireActiveLock {
		existing, found, err := selectActiveRunByTarget(ctx, tx, targetKind, targetKey)
		if err != nil {
			return Run{}, err
		}
		if found {
			return Run{}, ActiveRunError{Existing: existing}
		}
	}

	_, err = tx.ExecContext(ctx, `
INSERT INTO runs (
	id, kind, state, head_repository, head_branch, base_repository,
	pr_number, git_root, local_branch, head_sha, artifact_dir, work_dir,
	spec_slug, agent, model, reasoning_effort, owner_pid, owner_identity,
	owner_identity_unproven, created_at, updated_at
) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
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
		req.WorkDir,
		req.SpecSlug,
		req.Agent,
		req.Model,
		req.ReasoningEffort,
		nullableOwnerPID(req.OwnerPID),
		nullableOwnerIdentity(req.OwnerIdentity),
		req.OwnerPID > 0 && strings.TrimSpace(req.OwnerIdentity) == "",
		formatTime(now),
		formatTime(now),
	)
	if err != nil {
		return Run{}, fmt.Errorf("insert Run record: %w", err)
	}

	if acquireActiveLock {
		_, err = tx.ExecContext(ctx, `
INSERT INTO active_run_locks (target_kind, target_key, run_id, created_at)
VALUES (?, ?, ?, ?)`,
			targetKind,
			targetKey,
			runID,
			formatTime(now),
		)
		if err != nil {
			existing, found, selectErr := selectActiveRunByTarget(ctx, tx, targetKind, targetKey)
			if selectErr == nil && found {
				return Run{}, ActiveRunError{Existing: existing}
			}
			return Run{}, fmt.Errorf("acquire Active Run lock: %w", err)
		}
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

func (store *Store) CompleteRun(ctx context.Context, runID string, terminalState string) (CompleteRunResult, error) {
	runID = strings.TrimSpace(runID)
	if runID == "" {
		return CompleteRunResult{}, errors.New("complete Run: Run ID is required")
	}
	if !IsTerminalState(terminalState) {
		return CompleteRunResult{}, fmt.Errorf("complete Run %q: Run state %q is not terminal", runID, terminalState)
	}
	now := store.now()

	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return CompleteRunResult{}, fmt.Errorf("begin completion for Run %q: %w", runID, err)
	}
	defer rollbackUnlessCommitted(tx)

	terminalClause, terminalArguments := terminalStateExclusion()
	arguments := []any{
		terminalState,
		formatTime(now),
		formatTime(now),
		runID,
	}
	arguments = append(arguments, terminalArguments...)
	result, err := tx.ExecContext(ctx, `
UPDATE runs
SET state = ?, updated_at = ?, completed_at = ?
WHERE id = ?
  AND `+terminalClause,
		arguments...,
	)
	if err != nil {
		return CompleteRunResult{}, fmt.Errorf("compare-and-set terminal outcome for Run %q: %w", runID, err)
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return CompleteRunResult{}, fmt.Errorf("read completion result for Run %q: %w", runID, err)
	}
	if affected == 0 {
		run, err := selectRun(ctx, tx, runID)
		if errors.Is(err, sql.ErrNoRows) {
			return CompleteRunResult{}, fmt.Errorf("complete Run %q: Run does not exist", runID)
		}
		if err != nil {
			return CompleteRunResult{}, fmt.Errorf("read Run %q after completion compare-and-set: %w", runID, err)
		}
		if run.State != terminalState {
			return CompleteRunResult{}, TerminalOutcomeConflictError{
				RunID:     runID,
				Stored:    run.State,
				Requested: terminalState,
			}
		}
		if err := tx.Commit(); err != nil {
			return CompleteRunResult{}, fmt.Errorf("commit completion replay for Run %q: %w", runID, err)
		}
		return CompleteRunResult{Run: run}, nil
	}

	if _, err := tx.ExecContext(ctx, `DELETE FROM active_run_locks WHERE run_id = ?`, runID); err != nil {
		return CompleteRunResult{}, fmt.Errorf("release Active Run lock for Run %q: %w", runID, err)
	}

	run, err := selectRun(ctx, tx, runID)
	if err != nil {
		return CompleteRunResult{}, fmt.Errorf("read completed Run %q: %w", runID, err)
	}
	if err := tx.Commit(); err != nil {
		return CompleteRunResult{}, fmt.Errorf("commit completion for Run %q: %w", runID, err)
	}
	return CompleteRunResult{Run: run, Transitioned: true}, nil
}

func (store *Store) ReconcileIntegration(ctx context.Context, req IntegrationReconciliation) (Run, error) {
	if err := validateIntegrationReconciliation(req); err != nil {
		return Run{}, err
	}
	req.RunID = strings.TrimSpace(req.RunID)
	req.PreviousOutcome = strings.TrimSpace(req.PreviousOutcome)
	req.Classification = strings.TrimSpace(req.Classification)
	req.RunBranch = strings.TrimSpace(req.RunBranch)
	req.RunHead = strings.TrimSpace(req.RunHead)
	req.TargetBranch = strings.TrimSpace(req.TargetBranch)
	req.TargetHead = strings.TrimSpace(req.TargetHead)
	req.Worktree = strings.TrimSpace(req.Worktree)
	req.Reason = strings.TrimSpace(req.Reason)
	req.Action = strings.TrimSpace(req.Action)

	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return Run{}, fmt.Errorf("begin terminal Run reconciliation for Run %q: %w", req.RunID, err)
	}
	defer rollbackUnlessCommitted(tx)

	run, err := selectRun(ctx, tx, req.RunID)
	if errors.Is(err, sql.ErrNoRows) {
		return Run{}, fmt.Errorf("reconcile terminal Run %q: Run does not exist", req.RunID)
	}
	if err != nil {
		return Run{}, fmt.Errorf("read Run %q before terminal reconciliation: %w", req.RunID, err)
	}
	if run.Kind != KindImplement {
		return Run{}, fmt.Errorf(
			"reconcile terminal Run %q: Run kind %q is not %q",
			req.RunID,
			run.Kind,
			KindImplement,
		)
	}
	if strings.TrimSpace(run.LocalBranch) != req.TargetBranch {
		return Run{}, fmt.Errorf(
			"reconcile terminal Run %q: target branch %q does not match recorded target branch %q",
			req.RunID,
			req.TargetBranch,
			run.LocalBranch,
		)
	}
	if strings.TrimSpace(run.WorkDir) != req.Worktree {
		return Run{}, fmt.Errorf(
			"reconcile terminal Run %q: worktree %q does not match recorded worktree %q",
			req.RunID,
			req.Worktree,
			run.WorkDir,
		)
	}
	expectedRunBranch := RunBranchPrefix + req.RunID
	if req.RunBranch != expectedRunBranch {
		return Run{}, fmt.Errorf(
			"reconcile terminal Run %q: Run Branch %q does not match recorded Run Branch %q",
			req.RunID,
			req.RunBranch,
			expectedRunBranch,
		)
	}

	currentOutcome := req.PreviousOutcome
	if req.PreviousOutcome == StateIntegrationPending {
		currentOutcome = StateClean
	}
	payload, err := json.Marshal(map[string]string{
		"event":            "integration_reconciliation",
		"previous_outcome": req.PreviousOutcome,
		"current_outcome":  currentOutcome,
		"classification":   req.Classification,
		"run_branch":       req.RunBranch,
		"run_head":         req.RunHead,
		"target_branch":    req.TargetBranch,
		"target_head":      req.TargetHead,
		"worktree":         req.Worktree,
		"reason":           req.Reason,
		"action":           req.Action,
	})
	if err != nil {
		return Run{}, fmt.Errorf("encode terminal reconciliation for Run %q: %w", req.RunID, err)
	}

	var replay int
	err = tx.QueryRowContext(ctx, `
SELECT 1
FROM run_events
WHERE run_id = ? AND source = ? AND kind = ? AND payload = ?
LIMIT 1`,
		req.RunID,
		string(runevent.SourceDaemon),
		string(runevent.KindDaemonOutcome),
		string(payload),
	).Scan(&replay)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return Run{}, fmt.Errorf("inspect terminal reconciliation replay for Run %q: %w", req.RunID, err)
	}
	if err == nil {
		if run.State != currentOutcome {
			return Run{}, fmt.Errorf(
				"reconcile terminal Run %q: recorded evidence expects outcome %q but stored outcome is %q",
				req.RunID,
				currentOutcome,
				run.State,
			)
		}
		if err := tx.Commit(); err != nil {
			return Run{}, fmt.Errorf("commit terminal reconciliation replay for Run %q: %w", req.RunID, err)
		}
		return run, nil
	}
	if run.State != req.PreviousOutcome {
		if IsTerminalState(run.State) {
			return Run{}, TerminalOutcomeConflictError{
				RunID:     req.RunID,
				Stored:    run.State,
				Requested: currentOutcome,
			}
		}
		return Run{}, fmt.Errorf(
			"reconcile terminal Run %q: stored state %q is not terminal outcome %q",
			req.RunID,
			run.State,
			req.PreviousOutcome,
		)
	}

	if req.PreviousOutcome == StateIntegrationPending {
		result, err := tx.ExecContext(ctx, `
UPDATE runs
SET state = ?, updated_at = ?
WHERE id = ? AND state = ?`,
			StateClean,
			formatTime(req.Time),
			req.RunID,
			StateIntegrationPending,
		)
		if err != nil {
			return Run{}, fmt.Errorf("compare-and-set Integration Pending reconciliation for Run %q: %w", req.RunID, err)
		}
		affected, err := result.RowsAffected()
		if err != nil {
			return Run{}, fmt.Errorf("read Integration Pending reconciliation result for Run %q: %w", req.RunID, err)
		}
		if affected != 1 {
			return Run{}, fmt.Errorf(
				"reconcile Integration Pending Run %q: compare-and-set affected %d rows",
				req.RunID,
				affected,
			)
		}
	}
	if _, err := appendRunEvent(ctx, tx, runevent.RunEvent{
		RunID:   req.RunID,
		Source:  runevent.SourceDaemon,
		Kind:    runevent.KindDaemonOutcome,
		Summary: runevent.BoundSummary(fmt.Sprintf("Run reconciled %s to %s before terminal cleanup.", req.PreviousOutcome, currentOutcome)),
		Time:    req.Time,
		Payload: payload,
	}); err != nil {
		return Run{}, fmt.Errorf("journal terminal reconciliation for Run %q: %w", req.RunID, err)
	}

	reconciled, err := selectRun(ctx, tx, req.RunID)
	if err != nil {
		return Run{}, fmt.Errorf("read reconciled Run %q: %w", req.RunID, err)
	}
	if err := tx.Commit(); err != nil {
		return Run{}, fmt.Errorf("commit terminal reconciliation for Run %q: %w", req.RunID, err)
	}
	return reconciled, nil
}

func validateIntegrationReconciliation(req IntegrationReconciliation) error {
	runID := strings.TrimSpace(req.RunID)
	if runID == "" {
		return errors.New("reconcile terminal Run: Run ID is required")
	}
	if !IsTerminalState(strings.TrimSpace(req.PreviousOutcome)) {
		return fmt.Errorf("reconcile terminal Run %q: previous outcome %q is not terminal", runID, req.PreviousOutcome)
	}
	classification := strings.TrimSpace(req.Classification)
	if classification != "safe" && classification != "superseded" {
		return fmt.Errorf("reconcile terminal Run %q: classification must be safe or superseded", runID)
	}
	if strings.TrimSpace(req.RunBranch) == "" {
		return fmt.Errorf("reconcile terminal Run %q: Run Branch is required", runID)
	}
	if strings.TrimSpace(req.RunHead) == "" {
		return fmt.Errorf("reconcile terminal Run %q: Run Branch head is required", runID)
	}
	if strings.TrimSpace(req.TargetBranch) == "" {
		return fmt.Errorf("reconcile terminal Run %q: target branch is required", runID)
	}
	if strings.TrimSpace(req.TargetHead) == "" {
		return fmt.Errorf("reconcile terminal Run %q: target head is required", runID)
	}
	if strings.TrimSpace(req.Worktree) == "" {
		return fmt.Errorf("reconcile terminal Run %q: worktree is required", runID)
	}
	if strings.TrimSpace(req.Action) != "cleanup" {
		return fmt.Errorf("reconcile terminal Run %q: action must be cleanup", runID)
	}
	if req.Time.IsZero() {
		return fmt.Errorf("reconcile terminal Run %q: timestamp is required", runID)
	}
	return nil
}

func (store *Store) ReclaimOrphanedRun(ctx context.Context, run Run, reason string) error {
	runID := strings.TrimSpace(run.ID)
	if runID == "" {
		return errors.New("reclaim orphaned Run: Run ID is required")
	}
	if run.OwnerPID == nil || *run.OwnerPID <= 0 || ProcessAlive(*run.OwnerPID) {
		return nil
	}
	reason = strings.TrimSpace(reason)
	if reason == "" {
		reason = fmt.Sprintf("owner process %d not running; lock reclaimed", *run.OwnerPID)
	}
	now := store.now()

	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin orphaned Run reclamation: %w", err)
	}
	defer rollbackUnlessCommitted(tx)

	var currentState string
	if err := tx.QueryRowContext(ctx, `SELECT state FROM runs WHERE id = ?`, runID).Scan(&currentState); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("reclaim orphaned Run: Run %q does not exist", runID)
		}
		return fmt.Errorf("read Run %q before orphaned Run reclamation: %w", runID, err)
	}
	if IsTerminalState(currentState) {
		if err := tx.Commit(); err != nil {
			return fmt.Errorf("commit orphaned Run reclamation no-op: %w", err)
		}
		return nil
	}

	result, err := tx.ExecContext(ctx, `
UPDATE runs
SET state = ?, updated_at = ?, completed_at = ?
WHERE id = ? AND state = ?`,
		StateFailed,
		formatTime(now),
		formatTime(now),
		runID,
		currentState,
	)
	if err != nil {
		return fmt.Errorf("mark orphaned Run %q Failed: %w", runID, err)
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("read orphaned Run reclamation result: %w", err)
	}
	if affected == 0 {
		if err := tx.Commit(); err != nil {
			return fmt.Errorf("commit orphaned Run reclamation no-op: %w", err)
		}
		return nil
	}
	if _, err := tx.ExecContext(ctx, `DELETE FROM active_run_locks WHERE run_id = ?`, runID); err != nil {
		return fmt.Errorf("release orphaned Active Run lock: %w", err)
	}
	payload, err := json.Marshal(map[string]any{
		"state":     StateFailed,
		"reason":    reason,
		"owner_pid": *run.OwnerPID,
	})
	if err != nil {
		return fmt.Errorf("encode orphaned Run reclamation event: %w", err)
	}
	if _, err := appendRunEvent(ctx, tx, runevent.RunEvent{
		RunID:   runID,
		Source:  runevent.SourceDaemon,
		Kind:    runevent.KindDaemonOutcome,
		Summary: runevent.BoundSummary(fmt.Sprintf("Run reclaimed orphaned Active Run lock: %s.", reason)),
		Time:    now,
		Payload: payload,
	}); err != nil {
		return err
	}
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit orphaned Run reclamation: %w", err)
	}
	return nil
}

func (store *Store) RequestStop(ctx context.Context, runID string) error {
	runID = strings.TrimSpace(runID)
	if runID == "" {
		return errors.New("Run ID is required")
	}
	now := store.now()

	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin Stop Request: %w", err)
	}
	defer rollbackUnlessCommitted(tx)

	var state string
	if err := tx.QueryRowContext(ctx, `SELECT state FROM runs WHERE id = ?`, runID).Scan(&state); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("request Stop: Run %q does not exist", runID)
		}
		return fmt.Errorf("read Run %q state before Stop Request: %w", runID, err)
	}
	if IsTerminalState(state) {
		return fmt.Errorf("%w %q: state %s", ErrTerminalRunStopRequest, runID, state)
	}
	result, err := tx.ExecContext(ctx, `
UPDATE runs
SET stop_requested_at = COALESCE(stop_requested_at, ?), updated_at = ?
WHERE id = ?`,
		formatTime(now),
		formatTime(now),
		runID,
	)
	if err != nil {
		return fmt.Errorf("record Stop Request for Run %q: %w", runID, err)
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("read Stop Request result: %w", err)
	}
	if affected == 0 {
		return fmt.Errorf("request Stop: Run %q does not exist", runID)
	}
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit Stop Request: %w", err)
	}
	return nil
}

func (store *Store) StopRequested(ctx context.Context, runID string) (bool, error) {
	runID = strings.TrimSpace(runID)
	if runID == "" {
		return false, errors.New("Run ID is required")
	}
	var requestedAt sql.NullString
	if err := store.db.QueryRowContext(ctx, `SELECT stop_requested_at FROM runs WHERE id = ?`, runID).Scan(&requestedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, fmt.Errorf("read Stop Request flag: Run %q does not exist", runID)
		}
		return false, fmt.Errorf("read Stop Request flag for Run %q: %w", runID, err)
	}
	return requestedAt.Valid && strings.TrimSpace(requestedAt.String) != "", nil
}

func (store *Store) SetRunWorkDir(ctx context.Context, runID string, workDir string) (Run, error) {
	runID = strings.TrimSpace(runID)
	if runID == "" {
		return Run{}, errors.New("Run ID is required")
	}
	workDir = strings.TrimSpace(workDir)
	if workDir == "" {
		return Run{}, errors.New("Run Worktree path is required")
	}
	result, err := store.db.ExecContext(ctx, `
UPDATE runs SET work_dir = ?, updated_at = ? WHERE id = ?`,
		workDir,
		formatTime(store.now()),
		runID,
	)
	if err != nil {
		return Run{}, fmt.Errorf("record Run Worktree for Run %q: %w", runID, err)
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return Run{}, fmt.Errorf("read Run Worktree update result: %w", err)
	}
	if affected == 0 {
		return Run{}, fmt.Errorf("record Run Worktree: Run %q does not exist", runID)
	}
	return selectRun(ctx, store.db, runID)
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
	terminalClause, terminalArguments := terminalStateExclusion()
	arguments := []any{
		state,
		formatTime(store.now()),
		runID,
	}
	arguments = append(arguments, terminalArguments...)
	result, err := store.db.ExecContext(ctx, `
UPDATE runs
SET state = ?, updated_at = ?
WHERE id = ?
  AND `+terminalClause,
		arguments...,
	)
	if err != nil {
		return fmt.Errorf("update Run state: %w", err)
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("read Run state update result: %w", err)
	}
	if affected == 0 {
		run, found, selectErr := store.Run(ctx, runID)
		if selectErr != nil {
			return fmt.Errorf("read Run %q after state update compare-and-set: %w", runID, selectErr)
		}
		if !found {
			return fmt.Errorf("update Run state: Run %q does not exist", runID)
		}
		return TerminalOutcomeConflictError{
			RunID:     runID,
			Stored:    run.State,
			Requested: state,
		}
	}
	return nil
}

func (store *Store) ActiveRun(ctx context.Context, headRepository string, headBranch string) (Run, bool, error) {
	return selectActiveRunByTarget(ctx, store.db, targetKindPR, prTargetKey(headRepository, headBranch))
}

// ActiveReviewRunByTarget returns the newest non-terminal Run bound to the
// Head Repository and PR Head Branch, regardless of Active Run lock
// ownership. Runs created with the Branch Integrity bypass hold no lock, so
// target guards must scan the runs table instead of the lock table.
func (store *Store) ActiveReviewRunByTarget(ctx context.Context, headRepository string, headBranch string) (Run, bool, error) {
	sqlQuery := `
SELECT id, kind, state, head_repository, head_branch, base_repository,
       pr_number, git_root, local_branch, head_sha, artifact_dir, work_dir,
       spec_slug, agent, model, reasoning_effort, owner_pid, owner_identity, owner_identity_unproven, created_at, updated_at, completed_at
FROM runs
WHERE head_repository = ? AND head_branch = ?
ORDER BY created_at DESC, id DESC`
	rows, err := store.db.QueryContext(ctx, sqlQuery, strings.TrimSpace(headRepository), strings.TrimSpace(headBranch))
	if err != nil {
		return Run{}, false, fmt.Errorf("scan Active Runs by review target: %w", err)
	}
	defer func() {
		_ = rows.Close()
	}()
	for rows.Next() {
		run, err := scanRun(rows)
		if err != nil {
			return Run{}, false, fmt.Errorf("scan Active Run by review target: %w", err)
		}
		if !IsTerminalState(run.State) {
			return run, true, nil
		}
	}
	if err := rows.Err(); err != nil {
		return Run{}, false, fmt.Errorf("iterate Active Runs by review target: %w", err)
	}
	return Run{}, false, nil
}

// ActiveSpecRun returns the Active Run for one Spec work target, if any.
func (store *Store) ActiveSpecRun(ctx context.Context, gitRoot string, specSlug string) (Run, bool, error) {
	return selectActiveRunByTarget(ctx, store.db, targetKindSpec, specTargetKey(gitRoot, specSlug))
}

// ActiveRunInGitRoot returns the Active Run of any Kind whose Git root
// matches. It backs the ADR 0012 single-working-tree Preflight Validation
// until worktree-per-task lands.
func (store *Store) ActiveRunInGitRoot(ctx context.Context, gitRoot string) (Run, bool, error) {
	row := store.db.QueryRowContext(ctx, `
SELECT r.id, r.kind, r.state, r.head_repository, r.head_branch, r.base_repository,
       r.pr_number, r.git_root, r.local_branch, r.head_sha, r.artifact_dir, r.work_dir,
       r.spec_slug, r.agent, r.model, r.reasoning_effort, r.owner_pid, r.owner_identity, r.owner_identity_unproven, r.created_at, r.updated_at, r.completed_at
FROM active_run_locks l
JOIN runs r ON r.id = l.run_id
WHERE r.git_root = ?
ORDER BY l.created_at, r.id
LIMIT 1`, gitRoot)
	run, err := scanRun(row)
	if errors.Is(err, sql.ErrNoRows) {
		return Run{}, false, nil
	}
	if err != nil {
		return Run{}, false, err
	}
	return run, true, nil
}

// ListRuns returns matching Runs, newest first by creation time.
func (store *Store) ListRuns(ctx context.Context, query ListRunsQuery) ([]Run, error) {
	sqlQuery := `
SELECT id, kind, state, head_repository, head_branch, base_repository,
       pr_number, git_root, local_branch, head_sha, artifact_dir, work_dir,
       spec_slug, agent, model, reasoning_effort, owner_pid, owner_identity, owner_identity_unproven, created_at, updated_at, completed_at
FROM runs`
	args := []any{}
	gitRoot := strings.TrimSpace(query.GitRoot)
	if gitRoot != "" {
		sqlQuery += `
WHERE git_root = ?`
		args = append(args, gitRoot)
	}
	sqlQuery += `
ORDER BY created_at DESC, id DESC`

	rows, err := store.db.QueryContext(ctx, sqlQuery, args...)
	if err != nil {
		return nil, fmt.Errorf("list Runs: %w", err)
	}
	defer func() {
		_ = rows.Close()
	}()

	runs := []Run{}
	for rows.Next() {
		run, err := scanRun(rows)
		if err != nil {
			return nil, fmt.Errorf("scan Run listing: %w", err)
		}
		if !query.States.matches(run.State) {
			continue
		}
		runs = append(runs, run)
		if query.Limit > 0 && len(runs) == query.Limit {
			break
		}
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate Run listing: %w", err)
	}
	return runs, nil
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

func (store *Store) LatestKeptSpecRun(ctx context.Context, gitRoot string, specSlug string) (Run, bool, error) {
	row := store.db.QueryRowContext(ctx, `
SELECT id, kind, state, head_repository, head_branch, base_repository,
       pr_number, git_root, local_branch, head_sha, artifact_dir, work_dir,
       spec_slug, agent, model, reasoning_effort, owner_pid, owner_identity, owner_identity_unproven, created_at, updated_at, completed_at
FROM runs
WHERE kind = ? AND git_root = ? AND spec_slug = ?
  AND work_dir IS NOT NULL AND TRIM(work_dir) <> ''
  AND state IN (?, ?, ?, ?)
ORDER BY updated_at DESC, created_at DESC, id DESC
LIMIT 1`,
		KindImplement,
		gitRoot,
		specSlug,
		StateUnresolved,
		StateFailed,
		StateStopped,
		StateIntegrationPending,
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

func (store *Store) RunCount(ctx context.Context) (int, error) {
	var count int
	if err := store.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM runs`).Scan(&count); err != nil {
		return 0, fmt.Errorf("count Runs: %w", err)
	}
	return count, nil
}

// RunIDs lists persisted Run IDs for artifact cleanup comparisons.
func (store *Store) RunIDs(ctx context.Context) ([]string, error) {
	rows, err := store.db.QueryContext(ctx, `SELECT id FROM runs ORDER BY id`)
	if err != nil {
		return nil, fmt.Errorf("list Run IDs: %w", err)
	}
	defer func() {
		_ = rows.Close()
	}()

	runIDs := []string{}
	for rows.Next() {
		var runID string
		if err := rows.Scan(&runID); err != nil {
			return nil, fmt.Errorf("scan Run ID: %w", err)
		}
		runIDs = append(runIDs, runID)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate Run IDs: %w", err)
	}
	return runIDs, nil
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
	for _, terminalState := range terminalStates {
		if state == terminalState {
			return true
		}
	}
	return false
}

func terminalStateExclusion() (string, []any) {
	placeholders := make([]string, len(terminalStates))
	arguments := make([]any, len(terminalStates))
	for index, state := range terminalStates {
		placeholders[index] = "?"
		arguments[index] = state
	}
	return "state NOT IN (" + strings.Join(placeholders, ", ") + ")", arguments
}

const schemaVersion = 12

// activeRunLocksColumns is the schema v4 lock-table shape (ADR 0016): one
// Active Run per work target, keyed by (target_kind, target_key).
const activeRunLocksColumns = `(
	target_kind TEXT NOT NULL,
	target_key TEXT NOT NULL,
	run_id TEXT NOT NULL UNIQUE,
	created_at TEXT NOT NULL,
	PRIMARY KEY (target_kind, target_key),
	FOREIGN KEY (run_id) REFERENCES runs(id) ON DELETE CASCADE
)`

func (store *Store) migrate(ctx context.Context) error {
	if _, err := store.db.ExecContext(ctx, `PRAGMA foreign_keys = ON`); err != nil {
		return fmt.Errorf("apply Run Database migration: %w", err)
	}
	version, err := store.MigrationVersion(ctx)
	if err != nil {
		return err
	}
	switch version {
	case schemaVersion:
		return nil
	case 0:
		return store.applyMigration(ctx, createSchemaStatements())
	case 3:
		statements := append(migrateV3ToV4Statements(), migrateV4ToV5Statements()...)
		statements = append(statements, migrateV5ToV6Statements()...)
		statements = append(statements, migrateV6ToV7Statements()...)
		statements = append(statements, migrateV7ToV8Statements()...)
		statements = append(statements, migrateV8ToV9Statements()...)
		statements = append(statements, migrateV9ToV10Statements()...)
		statements = append(statements, migrateV10ToV11Statements()...)
		statements = append(statements, migrateV11ToV12Statements()...)
		return store.applyMigration(ctx, statements)
	case 4:
		statements := append(migrateV4ToV5Statements(), migrateV5ToV6Statements()...)
		statements = append(statements, migrateV6ToV7Statements()...)
		statements = append(statements, migrateV7ToV8Statements()...)
		statements = append(statements, migrateV8ToV9Statements()...)
		statements = append(statements, migrateV9ToV10Statements()...)
		statements = append(statements, migrateV10ToV11Statements()...)
		statements = append(statements, migrateV11ToV12Statements()...)
		return store.applyMigration(ctx, statements)
	case 5:
		if err := store.ensureAgentColumn(ctx); err != nil {
			return err
		}
		statements := append(migrateV5ToV6Statements(), migrateV6ToV7Statements()...)
		statements = append(statements, migrateV7ToV8Statements()...)
		statements = append(statements, migrateV8ToV9Statements()...)
		statements = append(statements, migrateV9ToV10Statements()...)
		statements = append(statements, migrateV10ToV11Statements()...)
		statements = append(statements, migrateV11ToV12Statements()...)
		return store.applyMigration(ctx, statements)
	case 6:
		statements := append(migrateV6ToV7Statements(), migrateV7ToV8Statements()...)
		statements = append(statements, migrateV8ToV9Statements()...)
		statements = append(statements, migrateV9ToV10Statements()...)
		statements = append(statements, migrateV10ToV11Statements()...)
		statements = append(statements, migrateV11ToV12Statements()...)
		return store.applyMigration(ctx, statements)
	case 7:
		statements := append(migrateV7ToV8Statements(), migrateV8ToV9Statements()...)
		statements = append(statements, migrateV9ToV10Statements()...)
		statements = append(statements, migrateV10ToV11Statements()...)
		statements = append(statements, migrateV11ToV12Statements()...)
		return store.applyMigration(ctx, statements)
	case 8:
		statements := append(migrateV8ToV9Statements(), migrateV9ToV10Statements()...)
		statements = append(statements, migrateV10ToV11Statements()...)
		statements = append(statements, migrateV11ToV12Statements()...)
		return store.applyMigration(ctx, statements)
	case 9:
		statements := append(migrateV9ToV10Statements(), migrateV10ToV11Statements()...)
		statements = append(statements, migrateV11ToV12Statements()...)
		return store.applyMigration(ctx, statements)
	case 10:
		statements := append(migrateV10ToV11Statements(), migrateV11ToV12Statements()...)
		return store.applyMigration(ctx, statements)
	case 11:
		return store.applyMigration(ctx, migrateV11ToV12Statements())
	default:
		return fmt.Errorf("migrate Run Database: schema version %d is not supported", version)
	}
}

func (store *Store) applyMigration(ctx context.Context, statements []string) error {
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin Run Database migration: %w", err)
	}
	defer rollbackUnlessCommitted(tx)
	for _, statement := range statements {
		if _, err := tx.ExecContext(ctx, statement); err != nil {
			return fmt.Errorf("apply Run Database migration: %w", err)
		}
	}
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit Run Database migration: %w", err)
	}
	return nil
}

// createSchemaStatements creates schema v12 directly on a fresh Run Database.
// spec_slug and the PR-shaped columns use the empty string for "not set";
// which fields a Run must carry is enforced by Kind in CreateRun.
func createSchemaStatements() []string {
	return []string{
		`CREATE TABLE IF NOT EXISTS runs (
			id TEXT PRIMARY KEY,
			kind TEXT NOT NULL,
			state TEXT NOT NULL,
			head_repository TEXT NOT NULL DEFAULT '',
			head_branch TEXT NOT NULL DEFAULT '',
			base_repository TEXT NOT NULL DEFAULT '',
			pr_number TEXT NOT NULL DEFAULT '',
			git_root TEXT NOT NULL,
			local_branch TEXT NOT NULL,
			head_sha TEXT NOT NULL DEFAULT '',
			artifact_dir TEXT NOT NULL DEFAULT '',
			work_dir TEXT,
			spec_slug TEXT NOT NULL DEFAULT '',
			agent TEXT NOT NULL DEFAULT '',
			model TEXT NOT NULL DEFAULT '',
			reasoning_effort TEXT NOT NULL DEFAULT '',
			owner_pid INTEGER,
			owner_identity TEXT,
			owner_identity_unproven INTEGER NOT NULL DEFAULT 0,
			stop_requested_at TEXT,
			created_at TEXT NOT NULL,
			updated_at TEXT NOT NULL,
			completed_at TEXT NOT NULL DEFAULT ''
		)`,
		`CREATE TABLE IF NOT EXISTS active_run_locks ` + activeRunLocksColumns,
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
		`CREATE TABLE IF NOT EXISTS run_agent_selections ` + runAgentSelectionsColumns,
		`CREATE INDEX IF NOT EXISTS idx_run_agent_selections_scope
			ON run_agent_selections (run_id, scope_kind, scope_id, attempt)`,
		`PRAGMA user_version = 12`,
	}
}

const runAgentSelectionsColumns = `(
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	run_id TEXT NOT NULL,
	scope_kind TEXT NOT NULL,
	scope_id TEXT NOT NULL,
	category TEXT NOT NULL,
	profile_source TEXT NOT NULL,
	attempt INTEGER NOT NULL,
	selection_role TEXT NOT NULL,
	fallback_index INTEGER NOT NULL DEFAULT 0,
	runtime TEXT NOT NULL,
	model TEXT NOT NULL,
	reasoning_effort TEXT NOT NULL,
	status TEXT NOT NULL,
	reason_code TEXT NOT NULL DEFAULT '',
	reason TEXT NOT NULL DEFAULT '',
	created_at TEXT NOT NULL,
	FOREIGN KEY (run_id) REFERENCES runs(id) ON DELETE CASCADE,
	CHECK (scope_kind IN ('task', 'qa', 'review')),
	CHECK (selection_role IN ('preferred', 'fallback')),
	CHECK (status IN ('attempting', 'active', 'failed', 'closed')),
	CHECK (attempt > 0),
	CHECK ((selection_role = 'preferred' AND fallback_index = 0) OR (selection_role = 'fallback' AND fallback_index > 0)),
	UNIQUE (run_id, scope_kind, scope_id, attempt)
)`

// migrateV3ToV4Statements rewrites a v3 Run Database in place (ADR 0016):
// every run row survives, and each existing lock row — always a review
// lock in v3 — is re-keyed to ("pr", "<head_repository>#<head_branch>").
func migrateV3ToV4Statements() []string {
	return []string{
		`ALTER TABLE runs ADD COLUMN spec_slug TEXT NOT NULL DEFAULT ''`,
		`CREATE TABLE active_run_locks_v4 ` + activeRunLocksColumns,
		`INSERT INTO active_run_locks_v4 (target_kind, target_key, run_id, created_at)
		 SELECT 'pr', head_repository || '#' || head_branch, run_id, created_at
		 FROM active_run_locks`,
		`DROP TABLE active_run_locks`,
		`ALTER TABLE active_run_locks_v4 RENAME TO active_run_locks`,
		`PRAGMA user_version = 4`,
	}
}

func migrateV4ToV5Statements() []string {
	return []string{
		`ALTER TABLE runs ADD COLUMN agent TEXT NOT NULL DEFAULT ''`,
		`ALTER TABLE runs ADD COLUMN stop_requested_at TEXT`,
		`PRAGMA user_version = 5`,
	}
}

func migrateV5ToV6Statements() []string {
	return []string{
		`ALTER TABLE runs ADD COLUMN work_dir TEXT`,
		`PRAGMA user_version = 6`,
	}
}

func migrateV6ToV7Statements() []string {
	return []string{
		`ALTER TABLE runs ADD COLUMN model TEXT NOT NULL DEFAULT ''`,
		`ALTER TABLE runs ADD COLUMN reasoning_effort TEXT NOT NULL DEFAULT ''`,
		`PRAGMA user_version = 7`,
	}
}

func migrateV7ToV8Statements() []string {
	return []string{
		`ALTER TABLE runs ADD COLUMN owner_pid INTEGER`,
		`PRAGMA user_version = 8`,
	}
}

func migrateV8ToV9Statements() []string {
	return []string{
		`CREATE TABLE IF NOT EXISTS run_agent_selections ` + runAgentSelectionsColumns,
		`CREATE INDEX IF NOT EXISTS idx_run_agent_selections_scope
			ON run_agent_selections (run_id, scope_kind, scope_id, attempt)`,
		`PRAGMA user_version = 9`,
	}
}

// migrateV9ToV10Statements adds the owner start-time identity token column.
// Legacy rows keep NULL, which Force Stop treats as the pre-identity
// PID-only owner proof.
func migrateV9ToV10Statements() []string {
	return []string{
		`ALTER TABLE runs ADD COLUMN owner_identity TEXT`,
		`PRAGMA user_version = 10`,
	}
}

func migrateV10ToV11Statements() []string {
	return []string{
		`PRAGMA user_version = 11`,
	}
}

func migrateV11ToV12Statements() []string {
	return []string{
		`ALTER TABLE runs ADD COLUMN owner_identity_unproven INTEGER NOT NULL DEFAULT 0`,
		`PRAGMA user_version = 12`,
	}
}

func (store *Store) ensureAgentColumn(ctx context.Context) error {
	exists, err := store.runColumnExists(ctx, "agent")
	if err != nil {
		return err
	}
	if exists {
		return nil
	}
	return store.applyMigration(ctx, []string{
		`ALTER TABLE runs ADD COLUMN agent TEXT NOT NULL DEFAULT ''`,
		`PRAGMA user_version = 5`,
	})
}

func (store *Store) runColumnExists(ctx context.Context, name string) (bool, error) {
	rows, err := store.db.QueryContext(ctx, `PRAGMA table_info(runs)`)
	if err != nil {
		return false, fmt.Errorf("inspect runs columns: %w", err)
	}
	defer func() {
		_ = rows.Close()
	}()
	for rows.Next() {
		var cid int
		var columnName string
		var columnType string
		var notNull int
		var defaultValue any
		var primaryKey int
		if err := rows.Scan(&cid, &columnName, &columnType, &notNull, &defaultValue, &primaryKey); err != nil {
			return false, fmt.Errorf("scan runs column: %w", err)
		}
		if columnName == name {
			return true, nil
		}
	}
	if err := rows.Err(); err != nil {
		return false, fmt.Errorf("inspect runs columns: %w", err)
	}
	return false, nil
}

func validateCreateRunRequest(req CreateRunRequest) error {
	if req.Kind == "" {
		return errors.New("Run kind is required")
	}
	var required map[string]string
	switch req.Kind {
	case KindFetch, KindResolve, KindWatch:
		required = map[string]string{
			"Head Repository":    req.HeadRepository,
			"PR Head Branch":     req.HeadBranch,
			"pull request":       req.PRNumber,
			"Git root":           req.GitRoot,
			"local branch":       req.LocalBranch,
			"HEAD":               req.HeadSHA,
			"Artifact Directory": req.ArtifactDir,
		}
	case KindImplement:
		required = map[string]string{
			"Git root":     req.GitRoot,
			"local branch": req.LocalBranch,
			"Spec slug":    req.SpecSlug,
		}
	default:
		return fmt.Errorf("Run kind %q is invalid", req.Kind)
	}
	for label, value := range required {
		if strings.TrimSpace(value) == "" {
			return fmt.Errorf("%s is required to create a Run", label)
		}
	}
	return nil
}

func nullableOwnerPID(pid int) any {
	if pid <= 0 {
		return nil
	}
	return pid
}

func nullableOwnerIdentity(identity string) any {
	identity = strings.TrimSpace(identity)
	if identity == "" {
		return nil
	}
	return identity
}

// lockTarget derives the Active Run lock key for a create request (ADR
// 0016): review Kinds lock the Open Pull Request, implement locks the Spec.
func lockTarget(req CreateRunRequest) (string, string) {
	if req.Kind == KindImplement {
		return targetKindSpec, specTargetKey(req.GitRoot, req.SpecSlug)
	}
	return targetKindPR, prTargetKey(req.HeadRepository, req.HeadBranch)
}

func prTargetKey(headRepository string, headBranch string) string {
	return headRepository + "#" + headBranch
}

func specTargetKey(gitRoot string, specSlug string) string {
	return gitRoot + "#" + specSlug
}

type runQuerier interface {
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
}

type runScanner interface {
	Scan(dest ...any) error
}

func selectActiveRunByTarget(ctx context.Context, querier runQuerier, targetKind string, targetKey string) (Run, bool, error) {
	row := querier.QueryRowContext(ctx, `
SELECT r.id, r.kind, r.state, r.head_repository, r.head_branch, r.base_repository,
       r.pr_number, r.git_root, r.local_branch, r.head_sha, r.artifact_dir, r.work_dir,
       r.spec_slug, r.agent, r.model, r.reasoning_effort, r.owner_pid, r.owner_identity, r.owner_identity_unproven, r.created_at, r.updated_at, r.completed_at
FROM active_run_locks l
JOIN runs r ON r.id = l.run_id
WHERE l.target_kind = ? AND l.target_key = ?`,
		targetKind,
		targetKey,
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
       pr_number, git_root, local_branch, head_sha, artifact_dir, work_dir,
       spec_slug, agent, model, reasoning_effort, owner_pid, owner_identity, owner_identity_unproven, created_at, updated_at, completed_at
FROM runs
WHERE id = ?`, runID)
	run, err := scanRun(row)
	if err != nil {
		return Run{}, fmt.Errorf("select Run %q: %w", runID, err)
	}
	return run, nil
}

func scanRun(row runScanner) (Run, error) {
	var run Run
	var createdAt string
	var updatedAt string
	var completedAt string
	var workDir sql.NullString
	var ownerPID sql.NullInt64
	var ownerIdentity sql.NullString
	var ownerIdentityUnproven int
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
		&workDir,
		&run.SpecSlug,
		&run.Agent,
		&run.Model,
		&run.ReasoningEffort,
		&ownerPID,
		&ownerIdentity,
		&ownerIdentityUnproven,
		&createdAt,
		&updatedAt,
		&completedAt,
	)
	if err != nil {
		return Run{}, err
	}
	if workDir.Valid {
		run.WorkDir = workDir.String
	}
	if ownerPID.Valid {
		pid := int(ownerPID.Int64)
		run.OwnerPID = &pid
	}
	if ownerIdentity.Valid {
		run.OwnerIdentity = ownerIdentity.String
	}
	run.OwnerIdentityUnproven = ownerIdentityUnproven != 0
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
