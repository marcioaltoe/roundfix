package daemon

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"roundfix/internal/agent"
	roundconfig "roundfix/internal/config"
	"roundfix/internal/runevent"
	"roundfix/internal/spec"
)

type AgentRuntimeFactory func(roundconfig.AgentSelection) (agent.RuntimeSpec, error)

type AgentSelectionProfiles map[roundconfig.WorkCategory]roundconfig.ResolvedProfile

type agentSessionScope struct {
	RunID    string
	Kind     string
	ID       string
	Category roundconfig.WorkCategory
	Session  agent.SessionRef
	Batch    int
}

type agentSessionOwner struct {
	engine         *Engine
	scope          agentSessionScope
	profile        roundconfig.ResolvedProfile
	runtimeFactory AgentRuntimeFactory

	activeRuntime agent.RuntimeSpec
	activeSession agent.SessionRef
	active        bool
	workStarted   bool
	attempts      []agentSelectionAttempt
}

func (plan TaskPlan) hasAgentSelectionProfiles() bool {
	return len(plan.AgentSelections) > 0
}

func (plan TaskPlan) agentSelectionOwnerConfig() agentSelectionOwnerConfig {
	return agentSelectionOwnerConfig{
		Profiles:       plan.AgentSelections,
		RuntimeFactory: plan.RuntimeFactory,
	}
}

func (plan CyclePlan) hasAgentSelectionProfiles() bool {
	return len(plan.AgentSelections) > 0
}

func (plan CyclePlan) agentSelectionOwnerConfig() agentSelectionOwnerConfig {
	return agentSelectionOwnerConfig{
		Profiles:       plan.AgentSelections,
		RuntimeFactory: plan.RuntimeFactory,
	}
}

type agentSelectionCandidate struct {
	Role          string
	FallbackIndex int
	Selection     roundconfig.AgentSelection
}

type agentSelectionAttempt struct {
	Candidate agentSelectionCandidate
	Err       error
}

type AgentSelectionExhaustedError struct {
	ScopeKind string
	ScopeID   string
	Category  roundconfig.WorkCategory
	Attempts  []agentSelectionAttempt
}

func (err *AgentSelectionExhaustedError) Error() string {
	if err == nil {
		return ""
	}
	parts := make([]string, 0, len(err.Attempts))
	for _, attempt := range err.Attempts {
		parts = append(parts, fmt.Sprintf("%s %s: %s", attempt.Candidate.Role, selectionLabel(attempt.Candidate.Selection), strings.TrimSpace(attempt.Err.Error())))
	}
	return fmt.Sprintf("Agent Selection exhausted for %s %s (%s); attempted: %s; recovery: %s",
		err.ScopeKind,
		err.ScopeID,
		err.Category,
		strings.Join(parts, "; "),
		agentSelectionRecoveryAction(err.Category),
	)
}

func (engine *Engine) taskAgentSessionOwner(plan TaskPlan, task spec.Task, ordinal int) (*agentSessionOwner, error) {
	if !plan.hasAgentSelectionProfiles() {
		return nil, nil
	}
	category := roundconfig.WorkCategory(task.Type)
	return engine.agentSessionOwner(plan.agentSelectionOwnerConfig(), agentSessionScope{
		RunID:    plan.RunID,
		Kind:     "task",
		ID:       task.ID,
		Category: category,
		Session:  agent.SessionRefForTask(plan.RunID, task.ID, plan.WorkDir),
		Batch:    ordinal,
	})
}

func (engine *Engine) qaAgentSessionOwner(plan TaskPlan, ordinal int) (*agentSessionOwner, error) {
	if !plan.hasAgentSelectionProfiles() {
		return nil, nil
	}
	return engine.agentSessionOwner(plan.agentSelectionOwnerConfig(), agentSessionScope{
		RunID:    plan.RunID,
		Kind:     "qa",
		ID:       "qa",
		Category: roundconfig.CategoryQA,
		Session:  agent.SessionRefForQA(plan.RunID, plan.WorkDir),
		Batch:    ordinal,
	})
}

func (engine *Engine) reviewAgentSessionOwner(plan CyclePlan, batchNumber int) (*agentSessionOwner, error) {
	if !plan.hasAgentSelectionProfiles() {
		return nil, nil
	}
	return engine.agentSessionOwner(plan.agentSelectionOwnerConfig(), agentSessionScope{
		RunID:    plan.RunID,
		Kind:     "review",
		ID:       fmt.Sprintf("batch-%03d", batchNumber),
		Category: roundconfig.CategoryReview,
		Session:  agent.SessionRefForReview(plan.RunID, batchNumber, plan.GitRoot),
		Batch:    batchNumber,
	})
}

func (engine *Engine) agentSessionOwner(config agentSelectionOwnerConfig, scope agentSessionScope) (*agentSessionOwner, error) {
	if config.RuntimeFactory == nil {
		return nil, errors.New("Agent Selection runtime factory is required")
	}
	profile, ok := config.Profiles[scope.Category]
	if !ok {
		return nil, fmt.Errorf("Agent Selection Profile for category %q is required", scope.Category)
	}
	return &agentSessionOwner{
		engine:         engine,
		scope:          scope,
		profile:        profile,
		runtimeFactory: config.RuntimeFactory,
	}, nil
}

func (engine *Engine) runAgentSession(ctx context.Context, owner *agentSessionOwner, req agent.ExecuteRequest) (agent.ExecuteResult, error) {
	if owner == nil {
		return engine.deps.Runner.Run(ctx, req, engine.deps.Sink)
	}
	return owner.Run(ctx, req)
}

func (owner *agentSessionOwner) Run(ctx context.Context, req agent.ExecuteRequest) (agent.ExecuteResult, error) {
	if owner == nil || owner.engine == nil {
		return agent.ExecuteResult{LogPath: req.LogPath}, errors.New("Agent Session owner is required")
	}
	if !owner.active {
		if err := owner.activate(ctx, req); err != nil {
			return agent.ExecuteResult{LogPath: req.LogPath}, err
		}
	}
	activeReq := owner.activeRequest(req)
	if !owner.workStarted {
		if err := owner.publishWorkStarted(ctx, activeReq); err != nil {
			owner.closeActive(context.WithoutCancel(ctx))
			return agent.ExecuteResult{LogPath: req.LogPath}, err
		}
		owner.workStarted = true
	}
	return owner.runPrepared(ctx, activeReq)
}

func (owner *agentSessionOwner) Close(ctx context.Context) error {
	if owner == nil || !owner.active {
		return nil
	}
	return owner.closeActive(ctx)
}

func (owner *agentSessionOwner) activate(ctx context.Context, req agent.ExecuteRequest) error {
	candidates := owner.candidates()
	if len(candidates) == 0 {
		return fmt.Errorf("Agent Selection Profile for category %q has no selections", owner.scope.Category)
	}
	for index, candidate := range candidates {
		if err := ctx.Err(); err != nil {
			return err
		}
		runtime, err := owner.runtimeFactory(candidate.Selection)
		if err != nil {
			return err
		}
		session := owner.sessionForCandidate(candidate)
		prepareReq := req
		prepareReq.Runtime = runtime
		prepareReq.Session = session
		prepareReq.GitRoot = session.WorkDir
		if strings.TrimSpace(prepareReq.GitRoot) == "" {
			prepareReq.GitRoot = req.GitRoot
		}
		if err := owner.prepareSession(ctx, prepareReq); err != nil {
			owner.closeAttempt(context.WithoutCancel(ctx), runtime, session)
			if !isSelectionStartFailure(err) || index == len(candidates)-1 {
				owner.attempts = append(owner.attempts, agentSelectionAttempt{Candidate: candidate, Err: err})
				if isSelectionStartFailure(err) && index == len(candidates)-1 {
					if publishErr := owner.publishExhausted(context.WithoutCancel(ctx)); publishErr != nil {
						return publishErr
					}
					return &AgentSelectionExhaustedError{
						ScopeKind: owner.scope.Kind,
						ScopeID:   owner.scope.ID,
						Category:  owner.scope.Category,
						Attempts:  append([]agentSelectionAttempt(nil), owner.attempts...),
					}
				}
				return err
			}
			owner.attempts = append(owner.attempts, agentSelectionAttempt{Candidate: candidate, Err: err})
			if err := owner.publishFallback(ctx, candidate, candidates[index+1], err); err != nil {
				return err
			}
			continue
		}
		owner.activeRuntime = runtime
		owner.activeSession = session
		owner.active = true
		return nil
	}
	return &AgentSelectionExhaustedError{
		ScopeKind: owner.scope.Kind,
		ScopeID:   owner.scope.ID,
		Category:  owner.scope.Category,
		Attempts:  append([]agentSelectionAttempt(nil), owner.attempts...),
	}
}

func (owner *agentSessionOwner) candidates() []agentSelectionCandidate {
	profile := owner.profile.Profile
	candidates := []agentSelectionCandidate{{
		Role:      "preferred",
		Selection: profile.Preferred,
	}}
	for index, selection := range profile.Fallbacks {
		candidates = append(candidates, agentSelectionCandidate{
			Role:          "fallback",
			FallbackIndex: index + 1,
			Selection:     selection,
		})
	}
	return candidates
}

func (owner *agentSessionOwner) sessionForCandidate(candidate agentSelectionCandidate) agent.SessionRef {
	session := owner.scope.Session
	session.WorkDir = strings.TrimSpace(session.WorkDir)
	if candidate.FallbackIndex > 0 {
		session.Name = fmt.Sprintf("%s-fallback-%02d", strings.TrimSpace(session.Name), candidate.FallbackIndex)
	}
	return session
}

func (owner *agentSessionOwner) activeRequest(req agent.ExecuteRequest) agent.ExecuteRequest {
	req.Runtime = owner.activeRuntime
	req.Session = owner.activeSession
	if strings.TrimSpace(req.GitRoot) == "" || strings.TrimSpace(req.GitRoot) != strings.TrimSpace(owner.activeSession.WorkDir) {
		req.GitRoot = owner.activeSession.WorkDir
	}
	return req
}

func (owner *agentSessionOwner) prepareSession(ctx context.Context, req agent.ExecuteRequest) error {
	preparer, ok := owner.engine.deps.Runner.(agent.SessionPreparer)
	if !ok {
		return nil
	}
	return preparer.PrepareSession(ctx, req, owner.engine.deps.Sink)
}

func (owner *agentSessionOwner) runPrepared(ctx context.Context, req agent.ExecuteRequest) (agent.ExecuteResult, error) {
	if runner, ok := owner.engine.deps.Runner.(agent.PreparedPromptRunner); ok {
		return runner.RunPrepared(ctx, req, owner.engine.deps.Sink)
	}
	return owner.engine.deps.Runner.Run(ctx, req, owner.engine.deps.Sink)
}

func (owner *agentSessionOwner) publishWorkStarted(ctx context.Context, req agent.ExecuteRequest) error {
	raw, err := json.Marshal(struct {
		Status string `json:"status"`
	}{Status: agent.AgentWorkStartedStatus})
	if err != nil {
		return err
	}
	if err := owner.engine.deps.Sink.Publish(ctx, runevent.RunEvent{
		RunID:   req.RunID,
		Batch:   req.Batch.Number,
		Source:  runevent.SourceAgent,
		Kind:    runevent.KindAgentStatus,
		Summary: "SESSION AGENT_WORK_STARTED\n",
		Time:    owner.engine.deps.Now(),
		Payload: raw,
	}); err != nil {
		return fmt.Errorf("publish Agent work-started Run Event: %w", err)
	}
	return nil
}

func (owner *agentSessionOwner) publishFallback(ctx context.Context, failed agentSelectionCandidate, next agentSelectionCandidate, cause error) error {
	payload := owner.notificationPayload(failed, next, cause)
	summary := fmt.Sprintf("Agent Selection fallback for %s %s (%s): activating fallback %d.",
		owner.scope.Kind,
		owner.scope.ID,
		owner.scope.Category,
		next.FallbackIndex,
	)
	if err := owner.engine.publishDaemonEvent(ctx, owner.profileRunID(), owner.notificationBatch(), runevent.KindDaemonAgentSelectionFallback, summary, payload); err != nil {
		return err
	}
	if _, err := fmt.Fprintf(owner.engine.deps.Progress, "roundfix: %s %s (%s) Agent Selection failed (%s); activating fallback %d %s.\n",
		owner.scope.Kind,
		owner.scope.ID,
		owner.scope.Category,
		selectionReasonCode(cause),
		next.FallbackIndex,
		selectionLabel(next.Selection),
	); err != nil {
		return fmt.Errorf("write Agent Selection fallback notification: %w", err)
	}
	return nil
}

func (owner *agentSessionOwner) publishExhausted(ctx context.Context) error {
	payload := map[string]any{
		"event":          "agent_selection_exhausted",
		"category":       string(owner.scope.Category),
		"scope_kind":     owner.scope.Kind,
		"scope_id":       owner.scope.ID,
		"scope_identity": owner.scope.Kind + ":" + owner.scope.ID,
		"attempts":       owner.attemptPayloads(),
		"recovery":       agentSelectionRecoveryAction(owner.scope.Category),
	}
	summary := fmt.Sprintf("Agent Selection exhausted for %s %s (%s).", owner.scope.Kind, owner.scope.ID, owner.scope.Category)
	return owner.engine.publishDaemonEvent(ctx, owner.profileRunID(), owner.notificationBatch(), runevent.KindDaemonAgentSelectionExhausted, summary, payload)
}

func (owner *agentSessionOwner) notificationPayload(failed agentSelectionCandidate, next agentSelectionCandidate, cause error) map[string]any {
	return map[string]any{
		"event":            "agent_selection_fallback",
		"category":         string(owner.scope.Category),
		"scope_kind":       owner.scope.Kind,
		"scope_id":         owner.scope.ID,
		"scope_identity":   owner.scope.Kind + ":" + owner.scope.ID,
		"failed_selection": selectionPayload(failed.Selection),
		"next_selection":   selectionPayload(next.Selection),
		"fallback_index":   next.FallbackIndex,
		"reason_code":      selectionReasonCode(cause),
		"reason":           strings.TrimSpace(cause.Error()),
		"automatic":        true,
	}
}

func (owner *agentSessionOwner) attemptPayloads() []map[string]any {
	attempts := make([]map[string]any, 0, len(owner.attempts))
	for _, attempt := range owner.attempts {
		attempts = append(attempts, map[string]any{
			"role":            attempt.Candidate.Role,
			"fallback_index":  attempt.Candidate.FallbackIndex,
			"selection":       selectionPayload(attempt.Candidate.Selection),
			"reason_code":     selectionReasonCode(attempt.Err),
			"reason":          strings.TrimSpace(attempt.Err.Error()),
			"recovery_action": agentSelectionRecoveryAction(owner.scope.Category),
		})
	}
	return attempts
}

func (owner *agentSessionOwner) profileRunID() string {
	return owner.scope.RunID
}

func (owner *agentSessionOwner) notificationBatch() int {
	return owner.scope.Batch
}

func (owner *agentSessionOwner) closeAttempt(ctx context.Context, runtime agent.RuntimeSpec, session agent.SessionRef) {
	if strings.TrimSpace(session.Name) == "" {
		return
	}
	_ = owner.engine.deps.Runner.EndSession(ctx, runtime, session)
}

func (owner *agentSessionOwner) closeActive(ctx context.Context) error {
	if !owner.active {
		return nil
	}
	err := owner.engine.deps.Runner.EndSession(ctx, owner.activeRuntime, owner.activeSession)
	if err == nil {
		_ = owner.publishSessionClosed(ctx)
	}
	owner.active = false
	return nil
}

func (owner *agentSessionOwner) publishSessionClosed(ctx context.Context) error {
	raw, err := json.Marshal(struct {
		Status string `json:"status"`
	}{Status: agent.AgentSessionClosedStatus})
	if err != nil {
		return err
	}
	if err := owner.engine.deps.Sink.Publish(ctx, runevent.RunEvent{
		RunID:   owner.scope.RunID,
		Batch:   owner.scope.Batch,
		Source:  runevent.SourceAgent,
		Kind:    runevent.KindAgentStatus,
		Summary: "SESSION SESSION_CLOSED\n",
		Time:    owner.engine.deps.Now(),
		Payload: raw,
	}); err != nil {
		return fmt.Errorf("publish Agent Session closed Run Event: %w", err)
	}
	return nil
}

func isSelectionStartFailure(err error) bool {
	var selectionErr *agent.SelectionPreflightError
	return errors.As(err, &selectionErr)
}

func selectionReasonCode(err error) string {
	var modelErr *agent.ModelNotAdvertisedError
	if errors.As(err, &modelErr) {
		return "model_not_advertised"
	}
	var selectionErr *agent.SelectionPreflightError
	if errors.As(err, &selectionErr) {
		operation := strings.ToLower(selectionErr.Operation)
		switch {
		case strings.Contains(operation, "model"):
			return "model_unavailable"
		case strings.Contains(operation, "reasoning") || strings.Contains(operation, "effort"):
			return "reasoning_effort_unavailable"
		case strings.Contains(operation, "runtime"):
			return "runtime_unavailable"
		default:
			return "selection_start_failed"
		}
	}
	return "selection_start_failed"
}

func selectionPayload(selection roundconfig.AgentSelection) map[string]string {
	return map[string]string{
		"runtime":          strings.TrimSpace(selection.Runtime),
		"model":            strings.TrimSpace(selection.Model),
		"reasoning_effort": strings.TrimSpace(selection.ReasoningEffort),
	}
}

func selectionLabel(selection roundconfig.AgentSelection) string {
	effort := strings.TrimSpace(selection.ReasoningEffort)
	if effort == "" {
		effort = "<model-managed>"
	}
	return fmt.Sprintf("%s/%s/%s", strings.TrimSpace(selection.Runtime), strings.TrimSpace(selection.Model), effort)
}

func agentSelectionRecoveryAction(category roundconfig.WorkCategory) string {
	return fmt.Sprintf("run `roundfix profiles validate --category %s`, then update the configured profile with `roundfix profiles configure --scope user` or `roundfix profiles configure --scope project`", category)
}

type agentSelectionOwnerConfig struct {
	Profiles       AgentSelectionProfiles
	RuntimeFactory AgentRuntimeFactory
}
