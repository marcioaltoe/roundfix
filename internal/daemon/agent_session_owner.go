package daemon

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"sync/atomic"

	"roundfix/internal/agent"
	roundconfig "roundfix/internal/config"
	"roundfix/internal/runevent"
	"roundfix/internal/spec"
	"roundfix/internal/store"
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

	activeRuntime             agent.RuntimeSpec
	activeSession             agent.SessionRef
	active                    bool
	workStarted               atomic.Bool
	candidateIndex            int
	selectionFailurePublished atomic.Bool
	attemptNumber             int
	attempts                  []agentSelectionAttempt
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
	for {
		if !owner.active {
			if err := owner.activate(ctx, req); err != nil {
				return agent.ExecuteResult{LogPath: req.LogPath}, err
			}
		}
		activeReq := owner.activeRequest(req)
		result, err := owner.runPrepared(ctx, activeReq)
		if err == nil {
			if publishErr := owner.publishWorkStartedOnce(ctx, activeReq); publishErr != nil {
				owner.closeActive(context.WithoutCancel(ctx))
				return result, publishErr
			}
			return result, nil
		}
		var selectionErr *agent.SelectionFailureError
		if !errors.As(err, &selectionErr) || owner.workStarted.Load() {
			return result, err
		}
		if !owner.selectionFailurePublished.Load() {
			if publishErr := owner.publishSelectionFailed(context.WithoutCancel(ctx), activeReq); publishErr != nil {
				return result, publishErr
			}
		}
		if fallbackErr := owner.fallbackAfterSelectionFailure(ctx, req, selectionErr); fallbackErr != nil {
			return result, fallbackErr
		}
	}
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
		if index < owner.candidateIndex {
			continue
		}
		if err := ctx.Err(); err != nil {
			return err
		}
		owner.candidateIndex = index
		owner.selectionFailurePublished.Store(false)
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
			if failure := selectionFailureForStart(runtime, err); failure != nil {
				err = failure
				if publishErr := owner.publishSelectionFailed(context.WithoutCancel(ctx), prepareReq); publishErr != nil {
					return publishErr
				}
			}
			owner.closeAttempt(context.WithoutCancel(ctx), runtime, session)
			if !isSelectionStartFailure(err) || index == len(candidates)-1 {
				owner.attempts = append(owner.attempts, agentSelectionAttempt{Candidate: candidate, Err: err})
				if isSelectionStartFailure(err) {
					if persistErr := owner.persistSelectionAttempt(context.WithoutCancel(ctx), candidate, store.AgentSelectionStatusFailed, err); persistErr != nil {
						return persistErr
					}
				}
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
			if persistErr := owner.persistSelectionAttempt(context.WithoutCancel(ctx), candidate, store.AgentSelectionStatusFailed, err); persistErr != nil {
				return persistErr
			}
			if err := owner.publishFallback(ctx, candidate, candidates[index+1], err); err != nil {
				return err
			}
			continue
		}
		if err := owner.persistSelectionAttempt(ctx, candidate, store.AgentSelectionStatusActive, nil); err != nil {
			owner.closeAttempt(context.WithoutCancel(ctx), runtime, session)
			return err
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
	sink := &agentSessionEventSink{owner: owner, req: req, next: owner.engine.deps.Sink}
	if runner, ok := owner.engine.deps.Runner.(agent.PreparedPromptRunner); ok {
		return runner.RunPrepared(ctx, req, sink)
	}
	return owner.engine.deps.Runner.Run(ctx, req, sink)
}

func (owner *agentSessionOwner) fallbackAfterSelectionFailure(ctx context.Context, req agent.ExecuteRequest, cause error) error {
	candidates := owner.candidates()
	failedIndex := owner.candidateIndex
	if failedIndex < 0 || failedIndex >= len(candidates) {
		return cause
	}
	failed := candidates[failedIndex]
	owner.closeAttempt(context.WithoutCancel(ctx), owner.activeRuntime, owner.activeSession)
	owner.active = false
	owner.attempts = append(owner.attempts, agentSelectionAttempt{Candidate: failed, Err: cause})
	if err := owner.persistSelectionAttempt(context.WithoutCancel(ctx), failed, store.AgentSelectionStatusFailed, cause); err != nil {
		return err
	}
	nextIndex := failedIndex + 1
	if nextIndex >= len(candidates) {
		if err := owner.publishExhausted(context.WithoutCancel(ctx)); err != nil {
			return err
		}
		return &AgentSelectionExhaustedError{
			ScopeKind: owner.scope.Kind,
			ScopeID:   owner.scope.ID,
			Category:  owner.scope.Category,
			Attempts:  append([]agentSelectionAttempt(nil), owner.attempts...),
		}
	}
	if err := owner.publishFallback(ctx, failed, candidates[nextIndex], cause); err != nil {
		return err
	}
	owner.candidateIndex = nextIndex
	owner.selectionFailurePublished.Store(false)
	return nil
}

func (owner *agentSessionOwner) publishWorkStartedOnce(ctx context.Context, req agent.ExecuteRequest) error {
	if !owner.workStarted.CompareAndSwap(false, true) {
		return nil
	}
	if err := owner.publishWorkStarted(ctx, req); err != nil {
		owner.workStarted.Store(false)
		return err
	}
	return nil
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

func (owner *agentSessionOwner) publishSelectionFailed(ctx context.Context, req agent.ExecuteRequest) error {
	raw, err := json.Marshal(struct {
		Status string `json:"status"`
	}{Status: agent.AgentSelectionFailedStatus})
	if err != nil {
		return err
	}
	if err := owner.engine.deps.Sink.Publish(ctx, runevent.RunEvent{
		RunID:   req.RunID,
		Batch:   req.Batch.Number,
		Source:  runevent.SourceAgent,
		Kind:    runevent.KindAgentStatus,
		Summary: "SESSION AGENT_SELECTION_FAILED\n",
		Time:    owner.engine.deps.Now(),
		Payload: raw,
	}); err != nil {
		return fmt.Errorf("publish Agent selection-failed Run Event: %w", err)
	}
	owner.selectionFailurePublished.Store(true)
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
	if err != nil {
		owner.active = false
		return fmt.Errorf("end Agent Session %q: %w", owner.activeSession.Name, err)
	}
	if persistErr := owner.persistSelectionAttempt(ctx, owner.activeCandidate(), store.AgentSelectionStatusClosed, nil); persistErr != nil {
		owner.active = false
		return persistErr
	}
	_ = owner.publishSessionClosed(ctx)
	owner.active = false
	return nil
}

type agentSelectionAttemptAppender interface {
	AppendAgentSelectionAttempt(context.Context, store.AgentSelectionAttemptRequest) (store.AgentSelectionAttempt, error)
}

type agentSelectionAttemptReader interface {
	AgentSelectionAttemptsForScope(context.Context, string, store.AgentSelectionScopeKind, string) ([]store.AgentSelectionAttempt, error)
}

func (owner *agentSessionOwner) persistSelectionAttempt(ctx context.Context, candidate agentSelectionCandidate, status store.AgentSelectionStatus, cause error) error {
	appender, ok := owner.engine.deps.Runs.(agentSelectionAttemptAppender)
	if !ok || appender == nil {
		return nil
	}
	attemptNumber, err := owner.nextSelectionAttemptNumber(ctx)
	if err != nil {
		return err
	}
	req := store.AgentSelectionAttemptRequest{
		RunID:           owner.scope.RunID,
		ScopeKind:       store.AgentSelectionScopeKind(owner.scope.Kind),
		ScopeID:         owner.scope.ID,
		Category:        string(owner.scope.Category),
		ProfileSource:   string(owner.profile.Source),
		Attempt:         attemptNumber,
		SelectionRole:   store.AgentSelectionRole(candidate.Role),
		FallbackIndex:   candidate.FallbackIndex,
		Runtime:         strings.TrimSpace(candidate.Selection.Runtime),
		Model:           strings.TrimSpace(candidate.Selection.Model),
		ReasoningEffort: strings.TrimSpace(candidate.Selection.ReasoningEffort),
		Status:          status,
		Time:            owner.engine.deps.Now(),
	}
	if cause != nil {
		req.ReasonCode = selectionReasonCode(cause)
		req.Reason = strings.TrimSpace(cause.Error())
	}
	if _, err := appender.AppendAgentSelectionAttempt(ctx, req); err != nil {
		return fmt.Errorf("persist Agent Selection %s for %s %s (%s): %w", status, owner.scope.Kind, owner.scope.ID, owner.scope.Category, err)
	}
	owner.attemptNumber = attemptNumber
	return nil
}

func (owner *agentSessionOwner) nextSelectionAttemptNumber(ctx context.Context) (int, error) {
	if owner.attemptNumber > 0 {
		return owner.attemptNumber + 1, nil
	}
	reader, ok := owner.engine.deps.Runs.(agentSelectionAttemptReader)
	if !ok || reader == nil {
		return 1, nil
	}
	attempts, err := reader.AgentSelectionAttemptsForScope(ctx, owner.scope.RunID, store.AgentSelectionScopeKind(owner.scope.Kind), owner.scope.ID)
	if err != nil {
		return 0, fmt.Errorf("read Agent Selection history for %s %s (%s): %w", owner.scope.Kind, owner.scope.ID, owner.scope.Category, err)
	}
	for _, attempt := range attempts {
		if attempt.Attempt > owner.attemptNumber {
			owner.attemptNumber = attempt.Attempt
		}
	}
	return owner.attemptNumber + 1, nil
}

func (owner *agentSessionOwner) activeCandidate() agentSelectionCandidate {
	candidates := owner.candidates()
	if owner.candidateIndex >= 0 && owner.candidateIndex < len(candidates) {
		return candidates[owner.candidateIndex]
	}
	active := selectionPayload(roundconfig.AgentSelection{
		Runtime:         owner.activeRuntime.ID,
		Model:           owner.activeRuntime.Model,
		ReasoningEffort: owner.activeRuntime.ReasoningEffort,
	})
	for _, candidate := range candidates {
		if equalSelectionPayload(selectionPayload(candidate.Selection), active) {
			return candidate
		}
	}
	return agentSelectionCandidate{Role: "preferred", Selection: roundconfig.AgentSelection{
		Runtime:         owner.activeRuntime.ID,
		Model:           owner.activeRuntime.Model,
		ReasoningEffort: owner.activeRuntime.ReasoningEffort,
	}}
}

func equalSelectionPayload(left map[string]string, right map[string]string) bool {
	return left["runtime"] == right["runtime"] &&
		left["model"] == right["model"] &&
		left["reasoning_effort"] == right["reasoning_effort"]
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
	var failure *agent.SelectionFailureError
	if errors.As(err, &failure) {
		return true
	}
	var selectionErr *agent.SelectionPreflightError
	if errors.As(err, &selectionErr) {
		return true
	}
	var adapterErr agent.AdapterProbeError
	return errors.As(err, &adapterErr)
}

func selectionFailureForStart(runtime agent.RuntimeSpec, err error) *agent.SelectionFailureError {
	var failure *agent.SelectionFailureError
	if errors.As(err, &failure) {
		return failure
	}
	if !isSelectionStartFailure(err) {
		return nil
	}
	return &agent.SelectionFailureError{
		Runtime: strings.TrimSpace(runtime.ID),
		Reason:  selectionReasonCode(err),
		Err:     err,
	}
}

func selectionReasonCode(err error) string {
	var failure *agent.SelectionFailureError
	if errors.As(err, &failure) {
		if failure.Err != nil {
			return selectionReasonCode(failure.Err)
		}
		return "selection_failed"
	}
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

type agentSessionEventSink struct {
	owner *agentSessionOwner
	req   agent.ExecuteRequest
	next  runevent.Sink
}

func (sink *agentSessionEventSink) Publish(ctx context.Context, event runevent.RunEvent) error {
	if sink == nil || sink.owner == nil {
		return errors.New("Agent Session event sink owner is required")
	}
	if agentStatusEventIs(event, agent.AgentWorkStartedStatus) {
		sink.owner.workStarted.Store(true)
	}
	if agentStatusEventIs(event, agent.AgentSelectionFailedStatus) {
		sink.owner.selectionFailurePublished.Store(true)
	}
	if agentOutputEvent(event) {
		if err := sink.owner.publishWorkStartedOnce(ctx, sink.req); err != nil {
			return err
		}
	}
	if sink.next == nil {
		return nil
	}
	return sink.next.Publish(ctx, event)
}

func agentOutputEvent(event runevent.RunEvent) bool {
	if event.Source != runevent.SourceAgent {
		return false
	}
	switch event.Kind {
	case runevent.KindAgentMessage,
		runevent.KindAgentThought,
		runevent.KindAgentToolStarted,
		runevent.KindAgentToolUpdated,
		runevent.KindAgentPlan,
		runevent.KindAgentRaw:
		return true
	default:
		return false
	}
}

func agentStatusEventIs(event runevent.RunEvent, want string) bool {
	if event.Source != runevent.SourceAgent || event.Kind != runevent.KindAgentStatus {
		return false
	}
	var payload struct {
		Status string `json:"status"`
	}
	return json.Unmarshal(event.Payload, &payload) == nil && payload.Status == want
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
