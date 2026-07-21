package agent

import (
	"context"
	"errors"
	"fmt"
	"strings"
)

const (
	SelectionModelNotAdvertised            = "model_not_advertised"
	SelectionReasoningControlNotAdvertised = "reasoning_control_not_advertised"
	SelectionModelVariantNotAdvertised     = "model_variant_not_advertised"
	SelectionRejected                      = "selection_rejected"
	EffectiveSelectionMismatch             = "effective_selection_mismatch"
	SessionCleanupFailed                   = "session_cleanup_failed"

	SelectionEncodingIndependent  = "independent"
	SelectionEncodingModelVariant = "model_variant"
	SelectionEncodingModelManaged = "model_managed"

	SelectionProofStatusProven = "proven"
)

// SelectionAssignment maps one canonical Agent Selection to controls
// advertised by the effective ACP adapter.
type SelectionAssignment struct {
	Runtime         string
	Model           string
	ReasoningEffort string
	AdapterModel    string
	ReasoningKey    string
	ReasoningValue  string
	Encoding        string
}

// SelectionProof records one exact, token-free Agent Selection proof.
type SelectionProof struct {
	Runtime         string
	Model           string
	ReasoningEffort string
	Assignment      SelectionAssignment
	Adapter         AdapterEvidence
	Status          string
}

// SessionSelectionRequest applies one canonical selection to an existing
// Agent Session using the complete capability state most recently returned by
// that Session.
type SessionSelectionRequest struct {
	Runtime      RuntimeSpec
	Session      SessionRef
	Capabilities SelectionCapabilities
}

// ProveExactSelection opens one disposable Agent Session, proves the exact
// requested tuple from advertised effective state, and closes the Session
// before returning.
func (runner ACPXRunner) ProveExactSelection(ctx context.Context, request ProbeRequest) (SelectionProof, error) {
	if ctx == nil {
		return SelectionProof{}, errors.New("selection proof context is required")
	}
	if err := validateRuntimeSelection(request.Runtime); err != nil {
		return SelectionProof{}, err
	}
	if strings.TrimSpace(request.Runtime.ID) == "" {
		return SelectionProof{}, errors.New("ACP Runtime id is required")
	}
	workDir := strings.TrimSpace(request.WorkDir)
	if workDir == "" {
		return SelectionProof{}, errors.New("Agent working directory is required")
	}

	setupCtx, setupCancel := context.WithTimeout(ctx, acpxPreflightSetupTimeout)
	adapter, err := CheckAdapter(setupCtx, request.Runtime)
	if err != nil {
		setupCancel()
		return SelectionProof{}, err
	}
	sessionName, err := disposablePreflightSessionName()
	if err != nil {
		setupCancel()
		return SelectionProof{}, err
	}
	codexEnv, err := runner.codexEnvForSession(setupCtx, request.Runtime, sessionName)
	if err != nil {
		setupCancel()
		return SelectionProof{}, err
	}

	session := SessionRef{Name: sessionName, WorkDir: workDir}
	capabilities, setupErr := runner.startSessionSelection(setupCtx, request.Runtime, session, adapter, codexEnv, true)
	var proof SelectionProof
	if setupErr == nil {
		proof, setupErr = runner.applySessionSelection(setupCtx, SessionSelectionRequest{
			Runtime:      request.Runtime,
			Session:      session,
			Capabilities: capabilities,
		}, codexEnv)
	}
	setupCancel()

	cleanupCtx, cleanupCancel := context.WithTimeout(context.WithoutCancel(ctx), acpxPreflightCleanupTimeout)
	defer cleanupCancel()
	cleanupErr := runner.closeDisposableSession(cleanupCtx, request.Runtime, sessionName, workDir)
	if setupErr != nil && cleanupErr != nil {
		return SelectionProof{}, errors.Join(setupErr, cleanupErr)
	}
	if setupErr != nil {
		return SelectionProof{}, setupErr
	}
	if cleanupErr != nil {
		return SelectionProof{}, cleanupErr
	}
	return proof, nil
}

func (runner ACPXRunner) startSessionSelection(ctx context.Context, runtime RuntimeSpec, session SessionRef, adapter AdapterEvidence, codexEnv []string, disposable bool) (SelectionCapabilities, error) {
	args, err := acpxEnsureArgs(runtime, session.Name, session.WorkDir)
	if err != nil {
		return SelectionCapabilities{}, err
	}
	if err := runner.runACPXCommandWithEnv(ctx, args, codexEnv); err != nil {
		if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
			return SelectionCapabilities{}, err
		}
		classified := classifyModelNotAdvertised(runtime, err)
		var modelErr *ModelNotAdvertisedError
		if errors.As(classified, &modelErr) {
			return SelectionCapabilities{}, modelErr
		}
		assignment := canonicalSelection(runtime)
		operation := "ensure live Agent Session"
		if disposable {
			operation = "ensure disposable Agent Session"
		}
		return SelectionCapabilities{}, &SelectionRejectedError{Assignment: assignment, Operation: operation, Err: err}
	}

	capabilities, err := runner.observeSelectionCapabilities(ctx, runtime, session, adapter, codexEnv)
	if err == nil {
		return capabilities, nil
	}
	if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
		return SelectionCapabilities{}, err
	}
	var evidenceErr *CapabilityEvidenceError
	if errors.As(err, &evidenceErr) {
		return SelectionCapabilities{}, err
	}
	return SelectionCapabilities{}, &SelectionRejectedError{Assignment: canonicalSelection(runtime), Operation: "inspect Agent Session capabilities", Err: err}
}

func canonicalSelection(runtime RuntimeSpec) SelectionAssignment {
	return SelectionAssignment{
		Runtime:         strings.TrimSpace(runtime.ID),
		Model:           strings.TrimSpace(runtime.Model),
		ReasoningEffort: strings.TrimSpace(runtime.ReasoningEffort),
	}
}

// SelectionUnsupportedError reports an exact requested tuple that cannot be
// represented by the advertised controls.
type SelectionUnsupportedError struct {
	Kind                string
	Runtime             string
	Model               string
	ReasoningEffort     string
	AdvertisedModels    []string
	AdvertisedReasoning []string
}

func (err *SelectionUnsupportedError) Error() string {
	if err == nil {
		return ""
	}
	message := fmt.Sprintf("Agent Selection %q/%q/%q is unsupported by advertised ACP controls", err.Runtime, err.Model, err.ReasoningEffort)
	if len(err.AdvertisedModels) > 0 {
		message += "; advertised Agent Models: " + strings.Join(err.AdvertisedModels, ", ")
	}
	if len(err.AdvertisedReasoning) > 0 {
		message += "; advertised reasoning efforts: " + strings.Join(err.AdvertisedReasoning, ", ")
	}
	return message + "; recovery: update the ACP Runtime or adapter and retry, or configure a different exact Agent Selection"
}

func (err *SelectionUnsupportedError) Classification() string {
	if err == nil {
		return ""
	}
	return err.Kind
}

// SelectionRejectedError reports an adapter rejection while applying an
// advertised assignment.
type SelectionRejectedError struct {
	Assignment SelectionAssignment
	Operation  string
	Err        error
}

func (err *SelectionRejectedError) Error() string {
	if err == nil {
		return ""
	}
	message := fmt.Sprintf("apply Agent Selection %q/%q/%q during %s: adapter rejected selection", err.Assignment.Runtime, err.Assignment.Model, err.Assignment.ReasoningEffort, err.Operation)
	if err.Err != nil {
		message += ": " + err.Err.Error()
	}
	return message + "; recovery: update the ACP Runtime or adapter and retry the exact Agent Selection"
}

func (err *SelectionRejectedError) Unwrap() error {
	if err == nil {
		return nil
	}
	return err.Err
}

func (err *SelectionRejectedError) Classification() string {
	return SelectionRejected
}

// EffectiveSelectionError reports a successful adapter update whose complete
// returned state does not represent the requested canonical tuple.
type EffectiveSelectionError struct {
	Assignment         SelectionAssignment
	EffectiveModel     string
	EffectiveReasoning string
}

func (err *EffectiveSelectionError) Error() string {
	if err == nil {
		return ""
	}
	return fmt.Sprintf("effective Agent Selection mismatch: requested %q/%q/%q, observed model %q and reasoning %q; recovery: update the ACP Runtime or adapter and retry", err.Assignment.Runtime, err.Assignment.Model, err.Assignment.ReasoningEffort, err.EffectiveModel, err.EffectiveReasoning)
}

func (err *EffectiveSelectionError) Classification() string {
	return EffectiveSelectionMismatch
}

// PlanSelectionAssignment deterministically maps a canonical Agent Selection
// to the controls advertised by one Agent Session.
func PlanSelectionAssignment(runtime RuntimeSpec, capabilities SelectionCapabilities) (SelectionAssignment, error) {
	assignment := SelectionAssignment{
		Runtime:         strings.TrimSpace(runtime.ID),
		Model:           strings.TrimSpace(runtime.Model),
		ReasoningEffort: strings.TrimSpace(runtime.ReasoningEffort),
	}
	if assignment.Runtime == "" {
		return SelectionAssignment{}, errors.New("ACP Runtime id is required")
	}
	if assignment.Model == "" {
		return SelectionAssignment{}, errors.New("agent model is required")
	}

	models := modelsForCanonical(capabilities.Models, assignment.Model)
	if len(models) == 0 {
		return SelectionAssignment{}, unsupportedSelection(SelectionModelNotAdvertised, assignment, capabilities)
	}
	base, hasBase := modelManagedRepresentation(models)
	if assignment.ReasoningEffort == "" {
		if !hasBase {
			return SelectionAssignment{}, unsupportedSelection(SelectionModelVariantNotAdvertised, assignment, capabilities)
		}
		assignment.AdapterModel = base.AdapterValue
		assignment.Encoding = SelectionEncodingModelManaged
		return assignment, nil
	}

	if hasBase && capabilities.ReasoningOption != nil && containsCapabilityValue(capabilities.ReasoningOption.Values, assignment.ReasoningEffort) {
		assignment.AdapterModel = base.AdapterValue
		assignment.ReasoningKey = capabilities.ReasoningOption.ID
		assignment.ReasoningValue = assignment.ReasoningEffort
		assignment.Encoding = SelectionEncodingIndependent
		return assignment, nil
	}
	if variant, ok := exactVariantRepresentation(models, assignment.ReasoningEffort); ok {
		assignment.AdapterModel = variant.AdapterValue
		assignment.Encoding = SelectionEncodingModelVariant
		return assignment, nil
	}
	if hasModelVariants(models) {
		return SelectionAssignment{}, unsupportedSelection(SelectionModelVariantNotAdvertised, assignment, capabilities)
	}
	return SelectionAssignment{}, unsupportedSelection(SelectionReasoningControlNotAdvertised, assignment, capabilities)
}

// ApplySessionSelection applies and proves one assignment on an existing
// Agent Session. Each successful update replaces the prior capability state.
func (runner ACPXRunner) ApplySessionSelection(ctx context.Context, request SessionSelectionRequest) (SelectionProof, error) {
	if ctx == nil {
		return SelectionProof{}, errors.New("selection application context is required")
	}
	codexEnv, err := runner.codexEnvForSession(ctx, request.Runtime, request.Session.Name)
	if err != nil {
		return SelectionProof{}, err
	}
	return runner.applySessionSelection(ctx, request, codexEnv)
}

func (runner ACPXRunner) applySessionSelection(ctx context.Context, request SessionSelectionRequest, codexEnv []string) (SelectionProof, error) {
	assignment, err := PlanSelectionAssignment(request.Runtime, request.Capabilities)
	if err != nil {
		return SelectionProof{}, err
	}
	state := request.Capabilities

	if state.CurrentModel != assignment.AdapterModel {
		state, err = runner.applySelectionOption(ctx, request, assignment, "model", assignment.AdapterModel, codexEnv)
		if err != nil {
			return SelectionProof{}, err
		}
		replanned, planErr := PlanSelectionAssignment(request.Runtime, state)
		if planErr != nil {
			return SelectionProof{}, planErr
		}
		if replanned != assignment {
			return SelectionProof{}, effectiveSelectionError(assignment, state)
		}
	}
	if assignment.Encoding == SelectionEncodingIndependent {
		state, err = runner.applySelectionOption(ctx, request, assignment, assignment.ReasoningKey, assignment.ReasoningValue, codexEnv)
		if err != nil {
			return SelectionProof{}, err
		}
	}
	if !selectionStateMatches(assignment, state) {
		return SelectionProof{}, effectiveSelectionError(assignment, state)
	}
	return SelectionProof{
		Runtime:         assignment.Runtime,
		Model:           assignment.Model,
		ReasoningEffort: assignment.ReasoningEffort,
		Assignment:      assignment,
		Adapter:         state.Adapter,
		Status:          SelectionProofStatusProven,
	}, nil
}

func (runner ACPXRunner) applySelectionOption(ctx context.Context, request SessionSelectionRequest, assignment SelectionAssignment, key string, value string, codexEnv []string) (SelectionCapabilities, error) {
	state, err := runner.acquireSelectionCapabilities(ctx, CapabilityAcquisitionRequest{
		Runtime:  request.Runtime,
		Session:  request.Session,
		ConfigID: key,
		Value:    value,
		Adapter:  request.Capabilities.Adapter,
	}, codexEnv)
	if err == nil {
		return state, nil
	}
	if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
		return SelectionCapabilities{}, err
	}
	var evidenceErr *CapabilityEvidenceError
	if errors.As(err, &evidenceErr) {
		return SelectionCapabilities{}, err
	}
	return SelectionCapabilities{}, &SelectionRejectedError{Assignment: assignment, Operation: "set " + key, Err: err}
}

func modelsForCanonical(models []ModelCapability, canonical string) []ModelCapability {
	matched := make([]ModelCapability, 0, len(models))
	for _, model := range models {
		if model.CanonicalModel == canonical {
			matched = append(matched, model)
		}
	}
	return matched
}

func modelManagedRepresentation(models []ModelCapability) (ModelCapability, bool) {
	for _, model := range models {
		if model.ModelManaged && model.ReasoningEffort == "" {
			return model, true
		}
	}
	return ModelCapability{}, false
}

func exactVariantRepresentation(models []ModelCapability, effort string) (ModelCapability, bool) {
	for _, model := range models {
		if !model.ModelManaged && model.ReasoningEffort == effort {
			return model, true
		}
	}
	return ModelCapability{}, false
}

func hasModelVariants(models []ModelCapability) bool {
	for _, model := range models {
		if !model.ModelManaged && model.ReasoningEffort != "" {
			return true
		}
	}
	return false
}

func containsCapabilityValue(values []string, want string) bool {
	for _, value := range values {
		if value == want {
			return true
		}
	}
	return false
}

func unsupportedSelection(kind string, assignment SelectionAssignment, capabilities SelectionCapabilities) *SelectionUnsupportedError {
	models := make([]string, 0, len(capabilities.Models))
	for _, model := range capabilities.Models {
		models = append(models, model.AdapterValue)
	}
	reasoning := []string(nil)
	if capabilities.ReasoningOption != nil {
		reasoning = append(reasoning, capabilities.ReasoningOption.Values...)
	}
	return &SelectionUnsupportedError{
		Kind:                kind,
		Runtime:             assignment.Runtime,
		Model:               assignment.Model,
		ReasoningEffort:     assignment.ReasoningEffort,
		AdvertisedModels:    models,
		AdvertisedReasoning: reasoning,
	}
}

func selectionStateMatches(assignment SelectionAssignment, state SelectionCapabilities) bool {
	if state.CurrentModel != assignment.AdapterModel {
		return false
	}
	model, ok := modelByAdapterValue(state.Models, state.CurrentModel)
	if !ok || model.CanonicalModel != assignment.Model {
		return false
	}
	switch assignment.Encoding {
	case SelectionEncodingIndependent:
		return model.ReasoningEffort == "" && state.ReasoningOption != nil && state.ReasoningOption.ID == assignment.ReasoningKey && state.ReasoningOption.CurrentValue == assignment.ReasoningValue
	case SelectionEncodingModelVariant:
		return model.ReasoningEffort == assignment.ReasoningEffort
	case SelectionEncodingModelManaged:
		return model.ModelManaged && model.ReasoningEffort == "" && state.ReasoningOption == nil
	default:
		return false
	}
}

func modelByAdapterValue(models []ModelCapability, adapterValue string) (ModelCapability, bool) {
	for _, model := range models {
		if model.AdapterValue == adapterValue {
			return model, true
		}
	}
	return ModelCapability{}, false
}

func effectiveSelectionError(assignment SelectionAssignment, state SelectionCapabilities) *EffectiveSelectionError {
	effectiveModel := state.CurrentModel
	effectiveReasoning := ""
	if model, ok := modelByAdapterValue(state.Models, state.CurrentModel); ok {
		effectiveModel = model.CanonicalModel
		effectiveReasoning = model.ReasoningEffort
	}
	if effectiveReasoning == "" && state.ReasoningOption != nil {
		effectiveReasoning = state.ReasoningOption.CurrentValue
	}
	return &EffectiveSelectionError{Assignment: assignment, EffectiveModel: effectiveModel, EffectiveReasoning: effectiveReasoning}
}
