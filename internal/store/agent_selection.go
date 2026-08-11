package store

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"roundfix/internal/runevent"
)

type AgentSelectionScopeKind string

const (
	AgentSelectionScopeTask   AgentSelectionScopeKind = "task"
	AgentSelectionScopeQA     AgentSelectionScopeKind = "qa"
	AgentSelectionScopeReview AgentSelectionScopeKind = "review"
)

type AgentSelectionRole string

const (
	AgentSelectionRolePreferred AgentSelectionRole = "preferred"
	AgentSelectionRoleFallback  AgentSelectionRole = "fallback"
)

type AgentSelectionStatus string

const (
	AgentSelectionStatusAttempting AgentSelectionStatus = "attempting"
	AgentSelectionStatusActive     AgentSelectionStatus = "active"
	AgentSelectionStatusFailed     AgentSelectionStatus = "failed"
	AgentSelectionStatusClosed     AgentSelectionStatus = "closed"
)

type AgentSelectionAttemptRequest struct {
	RunID           string
	ScopeKind       AgentSelectionScopeKind
	ScopeID         string
	Category        string
	ProfileSource   string
	Attempt         int
	SelectionRole   AgentSelectionRole
	FallbackIndex   int
	Runtime         string
	Model           string
	ReasoningEffort string
	Status          AgentSelectionStatus
	ReasonCode      string
	Reason          string
	Time            time.Time
}

type AgentSelectionAttempt struct {
	ID              int64
	RunID           string
	ScopeKind       AgentSelectionScopeKind
	ScopeID         string
	Category        string
	ProfileSource   string
	Attempt         int
	SelectionRole   AgentSelectionRole
	FallbackIndex   int
	Runtime         string
	Model           string
	ReasoningEffort string
	Status          AgentSelectionStatus
	ReasonCode      string
	Reason          string
	CreatedAt       time.Time
}

type AgentSelectionExhaustedRequest struct {
	RunID          string
	ScopeKind      AgentSelectionScopeKind
	ScopeID        string
	Category       string
	RecoveryAction string
	Time           time.Time
}

// AppendAgentSelectionAttempt appends one durable Agent Selection attempt and
// its matching Run Event in the same transaction. Attempts are strictly
// contiguous per Run scope: the first attempt is 1 and each next attempt is
// exactly previous+1.
func (store *Store) AppendAgentSelectionAttempt(ctx context.Context, req AgentSelectionAttemptRequest) (AgentSelectionAttempt, error) {
	if err := validateAgentSelectionAttemptRequest(req); err != nil {
		return AgentSelectionAttempt{}, err
	}
	createdAt := req.Time
	if createdAt.IsZero() {
		createdAt = store.now()
	}
	attempt := agentSelectionAttemptFromRequest(req, createdAt)

	var persisted AgentSelectionAttempt
	err := store.withWriteTx(ctx, "Agent Selection attempt append", func(tx *sql.Tx) error {
		if err := ensureRunExists(ctx, tx, attempt.RunID, "append Agent Selection attempt"); err != nil {
			return err
		}
		existing, updateExisting, err := ensureNextAgentSelectionAttempt(ctx, tx, attempt)
		if err != nil {
			return err
		}
		if updateExisting {
			attempt.ID = existing.ID
			if err := updateAgentSelectionAttemptStatus(ctx, tx, attempt); err != nil {
				return err
			}
		} else {
			if err := insertAgentSelectionAttempt(ctx, tx, &attempt); err != nil {
				return err
			}
		}
		event, err := agentSelectionAttemptEvent(attempt)
		if err != nil {
			return err
		}
		if _, err := appendRunEvent(ctx, tx, event); err != nil {
			return err
		}
		persisted = attempt
		return nil
	})
	if err != nil {
		return AgentSelectionAttempt{}, err
	}
	return persisted, nil
}

func (store *Store) AppendAgentSelectionExhausted(ctx context.Context, req AgentSelectionExhaustedRequest) (int64, error) {
	if err := validateAgentSelectionExhaustedRequest(req); err != nil {
		return 0, err
	}
	createdAt := req.Time
	if createdAt.IsZero() {
		createdAt = store.now()
	}
	runID := strings.TrimSpace(req.RunID)

	var cursor int64
	err := store.withWriteTx(ctx, "Agent Selection exhausted event append", func(tx *sql.Tx) error {
		if err := ensureRunExists(ctx, tx, runID, "append Agent Selection exhausted event"); err != nil {
			return err
		}
		attempts, err := selectAgentSelectionAttemptsForScope(ctx, tx, runID, req.ScopeKind, req.ScopeID)
		if err != nil {
			return err
		}
		if len(attempts) == 0 {
			return fmt.Errorf("append Agent Selection exhausted event for Run %q scope %s:%s: at least one persisted attempt is required", runID, req.ScopeKind, req.ScopeID)
		}
		latest := attempts[len(attempts)-1]
		if latest.Status != AgentSelectionStatusFailed {
			return fmt.Errorf("append Agent Selection exhausted event for Run %q scope %s:%s: latest attempt %d must be failed, got %q", runID, req.ScopeKind, req.ScopeID, latest.Attempt, latest.Status)
		}
		if strings.TrimSpace(req.Category) != latest.Category {
			return fmt.Errorf("append Agent Selection exhausted event for Run %q scope %s:%s: category %q does not match latest failed attempt category %q", runID, req.ScopeKind, req.ScopeID, strings.TrimSpace(req.Category), latest.Category)
		}
		req.Category = latest.Category
		event, err := agentSelectionExhaustedEvent(req, attempts, createdAt.UTC())
		if err != nil {
			return err
		}
		cursor, err = appendRunEvent(ctx, tx, event)
		return err
	})
	if err != nil {
		return 0, err
	}
	return cursor, nil
}

func (store *Store) AgentSelectionAttempts(ctx context.Context, runID string) ([]AgentSelectionAttempt, error) {
	runID = strings.TrimSpace(runID)
	if runID == "" {
		return nil, errors.New("list Agent Selection attempts: Run ID is required")
	}
	rows, err := store.db.QueryContext(ctx, `
SELECT id, run_id, scope_kind, scope_id, category, profile_source, attempt,
       selection_role, fallback_index, runtime, model, reasoning_effort,
       status, reason_code, reason, created_at
FROM run_agent_selections
WHERE run_id = ?
ORDER BY scope_kind, scope_id, attempt`, runID)
	if err != nil {
		return nil, fmt.Errorf("list Agent Selection attempts for Run %q: %w", runID, err)
	}
	defer func() {
		_ = rows.Close()
	}()
	return scanAgentSelectionAttempts(rows, runID)
}

// ActiveAgentSelectionScopes returns the latest persisted lifecycle for each
// Run scope only when that lifecycle is active. Scope kinds and IDs are
// ordered explicitly so Agent Session cleanup is deterministic.
func (store *Store) ActiveAgentSelectionScopes(ctx context.Context, runID string) ([]AgentSelectionAttempt, error) {
	runID = strings.TrimSpace(runID)
	if runID == "" {
		return nil, errors.New("list active Agent Selection scopes: Run ID is required")
	}
	rows, err := store.db.QueryContext(ctx, `
SELECT selection.id, selection.run_id, selection.scope_kind, selection.scope_id,
       selection.category, selection.profile_source, selection.attempt,
       selection.selection_role, selection.fallback_index, selection.runtime,
       selection.model, selection.reasoning_effort, selection.status,
       selection.reason_code, selection.reason, selection.created_at
FROM run_agent_selections AS selection
WHERE selection.run_id = ?
  AND selection.status = ?
  AND selection.attempt = (
      SELECT MAX(latest.attempt)
      FROM run_agent_selections AS latest
      WHERE latest.run_id = selection.run_id
        AND latest.scope_kind = selection.scope_kind
        AND latest.scope_id = selection.scope_id
  )
ORDER BY CASE selection.scope_kind
             WHEN 'task' THEN 1
             WHEN 'qa' THEN 2
             WHEN 'review' THEN 3
             ELSE 4
         END,
         selection.scope_id`,
		runID,
		string(AgentSelectionStatusActive),
	)
	if err != nil {
		return nil, fmt.Errorf("list active Agent Selection scopes for Run %q: %w", runID, err)
	}
	defer func() {
		_ = rows.Close()
	}()
	return scanAgentSelectionAttempts(rows, runID)
}

func (store *Store) AgentSelectionAttemptsForScope(ctx context.Context, runID string, scopeKind AgentSelectionScopeKind, scopeID string) ([]AgentSelectionAttempt, error) {
	runID = strings.TrimSpace(runID)
	if runID == "" {
		return nil, errors.New("list Agent Selection attempts: Run ID is required")
	}
	if !validAgentSelectionScopeKind(scopeKind) {
		return nil, fmt.Errorf("list Agent Selection attempts for Run %q: scope kind %q is invalid", runID, scopeKind)
	}
	if strings.TrimSpace(scopeID) == "" {
		return nil, fmt.Errorf("list Agent Selection attempts for Run %q: scope ID is required", runID)
	}
	return selectAgentSelectionAttemptsForScope(ctx, store.db, runID, scopeKind, scopeID)
}

func validateAgentSelectionExhaustedRequest(req AgentSelectionExhaustedRequest) error {
	if strings.TrimSpace(req.RunID) == "" {
		return errors.New("append Agent Selection exhausted event: Run ID is required")
	}
	if !validAgentSelectionScopeKind(req.ScopeKind) {
		return fmt.Errorf("append Agent Selection exhausted event for Run %q: scope kind %q is invalid; allowed: task, qa, review", strings.TrimSpace(req.RunID), req.ScopeKind)
	}
	if strings.TrimSpace(req.ScopeID) == "" {
		return fmt.Errorf("append Agent Selection exhausted event for Run %q: scope ID is required", strings.TrimSpace(req.RunID))
	}
	if strings.TrimSpace(req.Category) == "" {
		return fmt.Errorf("append Agent Selection exhausted event for Run %q: category is required", strings.TrimSpace(req.RunID))
	}
	return nil
}

type agentSelectionQuerier interface {
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
}

func selectAgentSelectionAttemptsForScope(ctx context.Context, querier agentSelectionQuerier, runID string, scopeKind AgentSelectionScopeKind, scopeID string) ([]AgentSelectionAttempt, error) {
	rows, err := querier.QueryContext(ctx, `
SELECT id, run_id, scope_kind, scope_id, category, profile_source, attempt,
       selection_role, fallback_index, runtime, model, reasoning_effort,
       status, reason_code, reason, created_at
FROM run_agent_selections
WHERE run_id = ? AND scope_kind = ? AND scope_id = ?
ORDER BY attempt`, runID, string(scopeKind), scopeID)
	if err != nil {
		return nil, fmt.Errorf("list Agent Selection attempts for Run %q scope %s:%s: %w", runID, scopeKind, scopeID, err)
	}
	defer func() {
		_ = rows.Close()
	}()
	return scanAgentSelectionAttempts(rows, runID)
}

func validateAgentSelectionAttemptRequest(req AgentSelectionAttemptRequest) error {
	if strings.TrimSpace(req.RunID) == "" {
		return errors.New("append Agent Selection attempt: Run ID is required")
	}
	if !validAgentSelectionScopeKind(req.ScopeKind) {
		return fmt.Errorf("append Agent Selection attempt for Run %q: scope kind %q is invalid; allowed: task, qa, review", strings.TrimSpace(req.RunID), req.ScopeKind)
	}
	if strings.TrimSpace(req.ScopeID) == "" {
		return fmt.Errorf("append Agent Selection attempt for Run %q: scope ID is required", strings.TrimSpace(req.RunID))
	}
	if strings.TrimSpace(req.Category) == "" {
		return fmt.Errorf("append Agent Selection attempt for Run %q: category is required", strings.TrimSpace(req.RunID))
	}
	if strings.TrimSpace(req.ProfileSource) == "" {
		return fmt.Errorf("append Agent Selection attempt for Run %q: profile source is required", strings.TrimSpace(req.RunID))
	}
	if req.Attempt <= 0 {
		return fmt.Errorf("append Agent Selection attempt for Run %q: attempt must be greater than zero", strings.TrimSpace(req.RunID))
	}
	if !validAgentSelectionRole(req.SelectionRole) {
		return fmt.Errorf("append Agent Selection attempt for Run %q: selection role %q is invalid; allowed: preferred, fallback", strings.TrimSpace(req.RunID), req.SelectionRole)
	}
	if req.SelectionRole == AgentSelectionRolePreferred && req.FallbackIndex != 0 {
		return fmt.Errorf("append Agent Selection attempt for Run %q: preferred selection fallback index must be 0", strings.TrimSpace(req.RunID))
	}
	if req.SelectionRole == AgentSelectionRoleFallback && req.FallbackIndex <= 0 {
		return fmt.Errorf("append Agent Selection attempt for Run %q: fallback selection requires a positive fallback index", strings.TrimSpace(req.RunID))
	}
	if strings.TrimSpace(req.Runtime) == "" {
		return fmt.Errorf("append Agent Selection attempt for Run %q: runtime is required", strings.TrimSpace(req.RunID))
	}
	if strings.TrimSpace(req.Model) == "" {
		return fmt.Errorf("append Agent Selection attempt for Run %q: model is required", strings.TrimSpace(req.RunID))
	}
	if !validAgentSelectionStatus(req.Status) {
		return fmt.Errorf("append Agent Selection attempt for Run %q: status %q is invalid; allowed: attempting, active, failed, closed", strings.TrimSpace(req.RunID), req.Status)
	}
	return nil
}

func validAgentSelectionScopeKind(scopeKind AgentSelectionScopeKind) bool {
	switch scopeKind {
	case AgentSelectionScopeTask, AgentSelectionScopeQA, AgentSelectionScopeReview:
		return true
	default:
		return false
	}
}

func validAgentSelectionRole(role AgentSelectionRole) bool {
	switch role {
	case AgentSelectionRolePreferred, AgentSelectionRoleFallback:
		return true
	default:
		return false
	}
}

func validAgentSelectionStatus(status AgentSelectionStatus) bool {
	switch status {
	case AgentSelectionStatusAttempting, AgentSelectionStatusActive, AgentSelectionStatusFailed, AgentSelectionStatusClosed:
		return true
	default:
		return false
	}
}

func agentSelectionAttemptFromRequest(req AgentSelectionAttemptRequest, createdAt time.Time) AgentSelectionAttempt {
	return AgentSelectionAttempt{
		RunID:           strings.TrimSpace(req.RunID),
		ScopeKind:       req.ScopeKind,
		ScopeID:         req.ScopeID,
		Category:        req.Category,
		ProfileSource:   req.ProfileSource,
		Attempt:         req.Attempt,
		SelectionRole:   req.SelectionRole,
		FallbackIndex:   req.FallbackIndex,
		Runtime:         req.Runtime,
		Model:           req.Model,
		ReasoningEffort: req.ReasoningEffort,
		Status:          req.Status,
		ReasonCode:      req.ReasonCode,
		Reason:          req.Reason,
		CreatedAt:       createdAt.UTC(),
	}
}

func ensureRunExists(ctx context.Context, tx *sql.Tx, runID string, operation string) error {
	var exists int
	err := tx.QueryRowContext(ctx, `SELECT 1 FROM runs WHERE id = ?`, runID).Scan(&exists)
	if errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("%s: Run %q does not exist", operation, runID)
	}
	if err != nil {
		return fmt.Errorf("%s: check Run %q: %w", operation, runID, err)
	}
	return nil
}

func ensureNextAgentSelectionAttempt(ctx context.Context, tx *sql.Tx, attempt AgentSelectionAttempt) (AgentSelectionAttempt, bool, error) {
	var previous sql.NullInt64
	if err := tx.QueryRowContext(ctx, `
SELECT MAX(attempt)
FROM run_agent_selections
WHERE run_id = ? AND scope_kind = ? AND scope_id = ?`,
		attempt.RunID,
		string(attempt.ScopeKind),
		attempt.ScopeID,
	).Scan(&previous); err != nil {
		return AgentSelectionAttempt{}, false, fmt.Errorf("read previous Agent Selection attempt for Run %q scope %s:%s: %w", attempt.RunID, attempt.ScopeKind, attempt.ScopeID, err)
	}
	next := 1
	if previous.Valid {
		next = int(previous.Int64) + 1
	}
	if attempt.Attempt == next {
		return AgentSelectionAttempt{}, false, nil
	}
	if previous.Valid && attempt.Attempt == int(previous.Int64) {
		existing, err := selectAgentSelectionAttemptForScopeAttempt(ctx, tx, attempt.RunID, attempt.ScopeKind, attempt.ScopeID, attempt.Attempt)
		if err != nil {
			return AgentSelectionAttempt{}, false, err
		}
		if err := validateAgentSelectionAttemptStatusUpdate(existing, attempt); err != nil {
			return AgentSelectionAttempt{}, false, err
		}
		return existing, true, nil
	}
	return AgentSelectionAttempt{}, false, fmt.Errorf("append Agent Selection attempt for Run %q scope %s:%s: attempt must be %d, got %d", attempt.RunID, attempt.ScopeKind, attempt.ScopeID, next, attempt.Attempt)
}

func selectAgentSelectionAttemptForScopeAttempt(ctx context.Context, tx *sql.Tx, runID string, scopeKind AgentSelectionScopeKind, scopeID string, attempt int) (AgentSelectionAttempt, error) {
	row := tx.QueryRowContext(ctx, `
SELECT id, run_id, scope_kind, scope_id, category, profile_source, attempt,
       selection_role, fallback_index, runtime, model, reasoning_effort,
       status, reason_code, reason, created_at
FROM run_agent_selections
WHERE run_id = ? AND scope_kind = ? AND scope_id = ? AND attempt = ?`,
		runID,
		string(scopeKind),
		scopeID,
		attempt,
	)
	existing, err := scanAgentSelectionAttempt(row)
	if err != nil {
		return AgentSelectionAttempt{}, fmt.Errorf("read Agent Selection attempt %d for Run %q scope %s:%s: %w", attempt, runID, scopeKind, scopeID, err)
	}
	return existing, nil
}

func validateAgentSelectionAttemptStatusUpdate(existing AgentSelectionAttempt, next AgentSelectionAttempt) error {
	if existing.Category != next.Category ||
		existing.ProfileSource != next.ProfileSource ||
		existing.SelectionRole != next.SelectionRole ||
		existing.FallbackIndex != next.FallbackIndex ||
		existing.Runtime != next.Runtime ||
		existing.Model != next.Model ||
		existing.ReasoningEffort != next.ReasoningEffort {
		return fmt.Errorf("append Agent Selection attempt for Run %q scope %s:%s attempt %d: selection fields are immutable for lifecycle updates", next.RunID, next.ScopeKind, next.ScopeID, next.Attempt)
	}
	if !validAgentSelectionStatusTransition(existing.Status, next.Status) {
		return fmt.Errorf("append Agent Selection attempt for Run %q scope %s:%s attempt %d: invalid status transition %q to %q", next.RunID, next.ScopeKind, next.ScopeID, next.Attempt, existing.Status, next.Status)
	}
	return nil
}

func validAgentSelectionStatusTransition(current AgentSelectionStatus, next AgentSelectionStatus) bool {
	switch current {
	case AgentSelectionStatusAttempting:
		return next == AgentSelectionStatusActive || next == AgentSelectionStatusFailed
	case AgentSelectionStatusActive:
		return next == AgentSelectionStatusClosed || next == AgentSelectionStatusFailed
	default:
		return false
	}
}

func insertAgentSelectionAttempt(ctx context.Context, tx *sql.Tx, attempt *AgentSelectionAttempt) error {
	row := tx.QueryRowContext(ctx, `
INSERT INTO run_agent_selections (
	run_id, scope_kind, scope_id, category, profile_source, attempt,
	selection_role, fallback_index, runtime, model, reasoning_effort,
	status, reason_code, reason, created_at
) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
RETURNING id`,
		attempt.RunID,
		string(attempt.ScopeKind),
		attempt.ScopeID,
		attempt.Category,
		attempt.ProfileSource,
		attempt.Attempt,
		string(attempt.SelectionRole),
		attempt.FallbackIndex,
		attempt.Runtime,
		attempt.Model,
		attempt.ReasoningEffort,
		string(attempt.Status),
		attempt.ReasonCode,
		attempt.Reason,
		formatTime(attempt.CreatedAt),
	)
	if err := row.Scan(&attempt.ID); err != nil {
		return fmt.Errorf("insert Agent Selection attempt for Run %q scope %s:%s attempt %d: %w", attempt.RunID, attempt.ScopeKind, attempt.ScopeID, attempt.Attempt, err)
	}
	return nil
}

func updateAgentSelectionAttemptStatus(ctx context.Context, tx *sql.Tx, attempt AgentSelectionAttempt) error {
	result, err := tx.ExecContext(ctx, `
UPDATE run_agent_selections
SET status = ?, reason_code = ?, reason = ?, created_at = ?
WHERE id = ?`,
		string(attempt.Status),
		attempt.ReasonCode,
		attempt.Reason,
		formatTime(attempt.CreatedAt),
		attempt.ID,
	)
	if err != nil {
		return fmt.Errorf("update Agent Selection attempt for Run %q scope %s:%s attempt %d: %w", attempt.RunID, attempt.ScopeKind, attempt.ScopeID, attempt.Attempt, err)
	}
	updated, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("count updated Agent Selection attempt for Run %q scope %s:%s attempt %d: %w", attempt.RunID, attempt.ScopeKind, attempt.ScopeID, attempt.Attempt, err)
	}
	if updated != 1 {
		return fmt.Errorf("update Agent Selection attempt for Run %q scope %s:%s attempt %d: updated %d rows", attempt.RunID, attempt.ScopeKind, attempt.ScopeID, attempt.Attempt, updated)
	}
	return nil
}

func agentSelectionAttemptEvent(attempt AgentSelectionAttempt) (runevent.RunEvent, error) {
	payload, err := agentSelectionAttemptPayload(attempt)
	if err != nil {
		return runevent.RunEvent{}, err
	}
	return runevent.RunEvent{
		RunID:   attempt.RunID,
		Source:  runevent.SourceDaemon,
		Kind:    agentSelectionEventKind(attempt.Status),
		Summary: runevent.BoundSummary(agentSelectionSummary(attempt)),
		Time:    attempt.CreatedAt,
		Payload: payload,
	}, nil
}

func agentSelectionEventKind(status AgentSelectionStatus) runevent.Kind {
	switch status {
	case AgentSelectionStatusAttempting:
		return runevent.KindDaemonAgentSelectionAttempt
	case AgentSelectionStatusActive:
		return runevent.KindDaemonAgentSelectionActive
	case AgentSelectionStatusFailed:
		return runevent.KindDaemonAgentSelectionFallback
	case AgentSelectionStatusClosed:
		return runevent.KindDaemonAgentSelectionClosed
	default:
		return runevent.KindDaemonAgentSelectionAttempt
	}
}

func agentSelectionPayloadEventName(status AgentSelectionStatus) string {
	switch status {
	case AgentSelectionStatusAttempting:
		return "agent_selection_attempt"
	case AgentSelectionStatusActive:
		return "agent_selection_active"
	case AgentSelectionStatusFailed:
		return "agent_selection_fallback"
	case AgentSelectionStatusClosed:
		return "agent_selection_closed"
	default:
		return "agent_selection_attempt"
	}
}

func agentSelectionAttemptPayload(attempt AgentSelectionAttempt) (json.RawMessage, error) {
	payload, err := json.Marshal(map[string]any{
		"event":            agentSelectionPayloadEventName(attempt.Status),
		"scope_kind":       string(attempt.ScopeKind),
		"scope_id":         attempt.ScopeID,
		"scope_identity":   string(attempt.ScopeKind) + ":" + attempt.ScopeID,
		"category":         attempt.Category,
		"profile_source":   attempt.ProfileSource,
		"attempt":          attempt.Attempt,
		"selection_role":   string(attempt.SelectionRole),
		"fallback_index":   attempt.FallbackIndex,
		"runtime":          attempt.Runtime,
		"model":            attempt.Model,
		"reasoning_effort": attempt.ReasoningEffort,
		"status":           string(attempt.Status),
		"reason_code":      attempt.ReasonCode,
		"reason":           attempt.Reason,
		"selection": map[string]string{
			"runtime":          attempt.Runtime,
			"model":            attempt.Model,
			"reasoning_effort": attempt.ReasoningEffort,
		},
	})
	if err != nil {
		return nil, fmt.Errorf("encode Agent Selection attempt event: %w", err)
	}
	return payload, nil
}

func agentSelectionExhaustedEvent(req AgentSelectionExhaustedRequest, attempts []AgentSelectionAttempt, createdAt time.Time) (runevent.RunEvent, error) {
	runID := strings.TrimSpace(req.RunID)
	payloadAttempts := make([]map[string]any, 0, len(attempts))
	for _, attempt := range attempts {
		payloadAttempts = append(payloadAttempts, map[string]any{
			"attempt":          attempt.Attempt,
			"selection_role":   string(attempt.SelectionRole),
			"fallback_index":   attempt.FallbackIndex,
			"runtime":          attempt.Runtime,
			"model":            attempt.Model,
			"reasoning_effort": attempt.ReasoningEffort,
			"status":           string(attempt.Status),
			"reason_code":      attempt.ReasonCode,
			"reason":           attempt.Reason,
		})
	}
	payload, err := json.Marshal(map[string]any{
		"event":           "agent_selection_exhausted",
		"scope_kind":      string(req.ScopeKind),
		"scope_id":        req.ScopeID,
		"scope_identity":  string(req.ScopeKind) + ":" + req.ScopeID,
		"category":        req.Category,
		"attempts":        payloadAttempts,
		"recovery_action": req.RecoveryAction,
	})
	if err != nil {
		return runevent.RunEvent{}, fmt.Errorf("encode Agent Selection exhausted event: %w", err)
	}
	return runevent.RunEvent{
		RunID:   runID,
		Source:  runevent.SourceDaemon,
		Kind:    runevent.KindDaemonAgentSelectionExhausted,
		Summary: runevent.BoundSummary(fmt.Sprintf("Agent Selection exhausted for %s %s (%s).", req.ScopeKind, req.ScopeID, req.Category)),
		Time:    createdAt,
		Payload: payload,
	}, nil
}

func agentSelectionSummary(attempt AgentSelectionAttempt) string {
	return fmt.Sprintf("Agent Selection attempt %d for %s %s (%s) %s.",
		attempt.Attempt,
		attempt.ScopeKind,
		attempt.ScopeID,
		attempt.Category,
		attempt.Status,
	)
}

func scanAgentSelectionAttempts(rows *sql.Rows, runID string) ([]AgentSelectionAttempt, error) {
	attempts := []AgentSelectionAttempt{}
	for rows.Next() {
		attempt, err := scanAgentSelectionAttempt(rows)
		if err != nil {
			return nil, fmt.Errorf("scan Agent Selection attempt for Run %q: %w", runID, err)
		}
		attempts = append(attempts, attempt)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate Agent Selection attempts for Run %q: %w", runID, err)
	}
	return attempts, nil
}

type agentSelectionScanner interface {
	Scan(dest ...any) error
}

func scanAgentSelectionAttempt(rows agentSelectionScanner) (AgentSelectionAttempt, error) {
	var attempt AgentSelectionAttempt
	var scopeKind string
	var selectionRole string
	var status string
	var createdAt string
	if err := rows.Scan(
		&attempt.ID,
		&attempt.RunID,
		&scopeKind,
		&attempt.ScopeID,
		&attempt.Category,
		&attempt.ProfileSource,
		&attempt.Attempt,
		&selectionRole,
		&attempt.FallbackIndex,
		&attempt.Runtime,
		&attempt.Model,
		&attempt.ReasoningEffort,
		&status,
		&attempt.ReasonCode,
		&attempt.Reason,
		&createdAt,
	); err != nil {
		return AgentSelectionAttempt{}, err
	}
	parsedAt, err := parseTime(createdAt)
	if err != nil {
		return AgentSelectionAttempt{}, err
	}
	attempt.ScopeKind = AgentSelectionScopeKind(scopeKind)
	attempt.SelectionRole = AgentSelectionRole(selectionRole)
	attempt.Status = AgentSelectionStatus(status)
	attempt.CreatedAt = parsedAt
	return attempt, nil
}
