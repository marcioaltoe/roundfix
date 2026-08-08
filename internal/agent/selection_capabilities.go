package agent

import (
	"context"
	"encoding/json"
	"errors"
	"sort"
	"strings"
)

const (
	CapabilityEvidenceInvalid = "capability_evidence_invalid"

	CapabilityIssueAdapterIdentityTooLong = "adapter_identity_too_long"
	CapabilityIssueAmbiguousModelVariant  = "ambiguous_model_variant"
	CapabilityIssueAmbiguousModelOption   = "ambiguous_model_option"
	CapabilityIssueAmbiguousReasoning     = "ambiguous_reasoning_option"
	CapabilityIssueContradictoryResponse  = "contradictory_response"
	CapabilityIssueDuplicateOptionID      = "duplicate_option_id"
	CapabilityIssueDuplicateOptionValue   = "duplicate_option_value"
	CapabilityIssueInvalidCurrentValue    = "invalid_current_value"
	CapabilityIssueInvalidOption          = "invalid_select_option"
	CapabilityIssueMalformedModelValue    = "malformed_model_value"
	CapabilityIssueMalformedResponse      = "malformed_response"
	CapabilityIssueMissingAdapter         = "missing_adapter_identity"
	CapabilityIssueMissingModel           = "missing_model_state"
	CapabilityIssueMissingOptions         = "missing_config_options"
	CapabilityIssueResponseTooLarge       = "response_too_large"
	CapabilityIssueTooManyOptions         = "too_many_config_options"
	CapabilityIssueTooManyValues          = "too_many_option_values"

	MaxCapabilityDiagnosticIssues = 8
	// maxCapabilityResponseBytes bounds one adapter response. The measured
	// ceiling is the OpenCode session snapshot, 50,590 bytes with 417
	// advertised models on 2026-08-08, so this leaves roughly twenty times
	// today's payload. It bounds reading; retention bounds what is kept.
	maxCapabilityResponseBytes = 1024 * 1024
	maxCapabilityOptions       = 32
	// maxRetainedCapabilityValues bounds how many advertised values one
	// projected option keeps, not how many a runtime may advertise.
	maxRetainedCapabilityValues = 64
	// maxAdvertisedCapabilityValues is the absolute ceiling above which an
	// advertised option is refused as implausible rather than merely large.
	maxAdvertisedCapabilityValues = 4096
	maxCapabilityIDBytes          = 64
	maxCapabilityValueBytes       = 160
	maxAdapterIdentityBytes       = 512
)

// SelectionCapabilities is the bounded projection of one documented ACP
// configOptions response. Empty ModelCapability.ReasoningEffort is explicit
// model-managed state; it is never inferred from a model's display name.
type SelectionCapabilities struct {
	CurrentModel    string
	Models          []ModelCapability
	Options         []SelectCapability
	ReasoningOption *SelectCapability
	Adapter         AdapterEvidence
}

// ModelCapability maps one advertised adapter value to its canonical model
// and, only for an explicitly encoded variant, its reasoning effort.
type ModelCapability struct {
	AdapterValue    string
	CanonicalModel  string
	ReasoningEffort string
	ModelManaged    bool
}

// SelectCapability contains only the documented select-option fields needed
// for exact assignment and proof. Values holds what the projection retained;
// AdvertisedCount holds how many valid values the adapter offered, so a
// diagnostic can report both instead of implying the adapter offers only what
// was kept.
type SelectCapability struct {
	ID              string
	CurrentValue    string
	Values          []string
	AdvertisedCount int
}

// SelectionRetention names the advertised values a capability projection must
// keep when an option advertises more values than Roundfix retains. It is the
// requested Agent Selection, never inferred from the payload.
type SelectionRetention struct {
	Model           string
	ReasoningEffort string
}

// RetentionFor returns the retention one runtime's requested Agent Selection
// implies.
func RetentionFor(runtime RuntimeSpec) SelectionRetention {
	return SelectionRetention{
		Model:           strings.TrimSpace(runtime.Model),
		ReasoningEffort: strings.TrimSpace(runtime.ReasoningEffort),
	}
}

// CapabilityAcquisitionRequest identifies one public ACPX config-option
// update whose response carries the complete current configOptions state.
type CapabilityAcquisitionRequest struct {
	Runtime  RuntimeSpec
	Session  SessionRef
	ConfigID string
	Value    string
	Adapter  AdapterEvidence
}

// CapabilityEvidenceError is a bounded typed failure for malformed,
// ambiguous, or contradictory public capability evidence.
type CapabilityEvidenceError struct {
	Issues []string
}

func (err *CapabilityEvidenceError) Error() string {
	if err == nil || len(err.Issues) == 0 {
		return "capability evidence invalid"
	}
	return "capability evidence invalid: " + strings.Join(err.Issues, ", ")
}

func (err *CapabilityEvidenceError) Classification() string {
	return CapabilityEvidenceInvalid
}

// CapabilityAcquisitionError keeps ACPX command diagnostics out of the public
// error string while preserving the underlying error for typed inspection.
type CapabilityAcquisitionError struct {
	Err error
}

func (err *CapabilityAcquisitionError) Error() string {
	return "acquire Agent Selection capabilities through acpx: command failed"
}

func (err *CapabilityAcquisitionError) Unwrap() error {
	if err == nil {
		return nil
	}
	return err.Err
}

type acpxCapabilityResponse struct {
	Action        string                    `json:"action"`
	ConfigID      string                    `json:"configId"`
	Value         string                    `json:"value"`
	ConfigOptions *[]acpSelectOptionPayload `json:"configOptions"`
}

type acpxSessionCapabilityResponse struct {
	Schema string `json:"schema"`
	ACPX   *struct {
		CurrentModelID string                    `json:"current_model_id"`
		ConfigOptions  *[]acpSelectOptionPayload `json:"config_options"`
	} `json:"acpx"`
}

type acpxModelSetResponse struct {
	Action  string `json:"action"`
	ModelID string `json:"modelId"`
}

type acpSelectOptionPayload struct {
	ID           string                   `json:"id"`
	Category     string                   `json:"category"`
	Type         string                   `json:"type"`
	CurrentValue any                      `json:"currentValue"`
	Options      *[]acpSelectValuePayload `json:"options"`
}

type acpSelectValuePayload struct {
	Value string `json:"value"`
}

// ParseSessionConfigOptions validates and bounds the machine-readable ACPX
// projection of a documented ACP configOptions response. Retention names the
// requested Agent Selection so an advertised option larger than Roundfix
// retains keeps the values that selection needs.
func ParseSessionConfigOptions(payload []byte, adapter AdapterEvidence, retention SelectionRetention) (SelectionCapabilities, error) {
	issues := newCapabilityIssueSet()
	if len(payload) == 0 {
		issues.add(CapabilityIssueMalformedResponse)
		return SelectionCapabilities{}, issues.err()
	}
	if len(payload) > maxCapabilityResponseBytes {
		issues.add(CapabilityIssueResponseTooLarge)
		return SelectionCapabilities{}, issues.err()
	}

	var response acpxCapabilityResponse
	if err := json.Unmarshal(payload, &response); err != nil {
		issues.add(CapabilityIssueMalformedResponse)
		return SelectionCapabilities{}, issues.err()
	}

	boundedAdapter := validateCapabilityAdapter(adapter, issues)
	if response.Action != "config_set" || !boundedCapabilityID(response.ConfigID) || !boundedCapabilityValue(response.Value) {
		issues.add(CapabilityIssueMalformedResponse)
	}
	if response.ConfigOptions == nil {
		issues.add(CapabilityIssueMissingOptions)
		return SelectionCapabilities{}, issues.err()
	}
	if len(*response.ConfigOptions) == 0 {
		issues.add(CapabilityIssueMissingOptions)
	}
	if len(*response.ConfigOptions) > maxCapabilityOptions {
		issues.add(CapabilityIssueTooManyOptions)
	}

	seenIDs := make(map[string]struct{}, len(*response.ConfigOptions))
	selectOptions := make([]SelectCapability, 0, len(*response.ConfigOptions))
	modelOptionIDs := make([]string, 0, 1)
	reasoningOptionIDs := make([]string, 0, 1)
	configuredValueMatched := false

	for _, rawOption := range boundedConfigOptions(*response.ConfigOptions) {
		option, ok := projectSelectCapability(rawOption, retention, issues)
		if rawOption.ID != "" {
			if _, exists := seenIDs[rawOption.ID]; exists {
				issues.add(CapabilityIssueDuplicateOptionID)
			} else {
				seenIDs[rawOption.ID] = struct{}{}
			}
		}
		if !ok {
			continue
		}
		selectOptions = append(selectOptions, option)
		if option.ID == response.ConfigID {
			configuredValueMatched = option.CurrentValue == response.Value
		}
		if rawOption.Category == "model" || option.ID == "model" {
			modelOptionIDs = append(modelOptionIDs, option.ID)
		}
		if option.ID == acpxCodexReasoningEffortKey || option.ID == acpxGenericReasoningEffortKey {
			reasoningOptionIDs = append(reasoningOptionIDs, option.ID)
		}
	}

	if !configuredValueMatched {
		issues.add(CapabilityIssueContradictoryResponse)
	}
	if len(modelOptionIDs) == 0 {
		issues.add(CapabilityIssueMissingModel)
	}
	if len(modelOptionIDs) > 1 {
		issues.add(CapabilityIssueAmbiguousModelOption)
	}
	if len(reasoningOptionIDs) > 1 {
		issues.add(CapabilityIssueAmbiguousReasoning)
	}

	sort.Slice(selectOptions, func(left, right int) bool {
		return selectOptions[left].ID < selectOptions[right].ID
	})
	if err := issues.err(); err != nil {
		return SelectionCapabilities{}, err
	}

	var reasoning *SelectCapability
	if len(reasoningOptionIDs) == 1 {
		copyOption := selectOptions[selectCapabilityIndex(selectOptions, reasoningOptionIDs[0])]
		copyOption.Values = append([]string(nil), copyOption.Values...)
		reasoning = &copyOption
	}

	modelOption := selectOptions[selectCapabilityIndex(selectOptions, modelOptionIDs[0])]
	opaqueModels := reasoning != nil
	models := make([]ModelCapability, 0, len(modelOption.Values))
	seenModels := make(map[string]struct{}, len(modelOption.Values))
	for _, value := range modelOption.Values {
		model, ok := parseModelCapability(value, opaqueModels)
		if !ok {
			issues.add(CapabilityIssueMalformedModelValue)
			continue
		}
		// Under opaque parsing an advertised identifier is a whole model
		// identity, so two entries that share a canonical prefix are an alias
		// group for one model, not an ambiguous variant encoding. Advertised
		// values are already unique here, so keying on the advertised value
		// keeps the dedup honest without inventing an ambiguity. The variant
		// encoding keeps its fail-closed canonical/effort ambiguity check.
		key := model.AdapterValue
		if !opaqueModels {
			key = model.CanonicalModel + "\x00" + model.ReasoningEffort
		}
		if _, exists := seenModels[key]; exists {
			issues.add(CapabilityIssueAmbiguousModelVariant)
			continue
		}
		seenModels[key] = struct{}{}
		models = append(models, model)
	}
	if err := issues.err(); err != nil {
		return SelectionCapabilities{}, err
	}

	return SelectionCapabilities{
		CurrentModel:    modelOption.CurrentValue,
		Models:          models,
		Options:         selectOptions,
		ReasoningOption: reasoning,
		Adapter:         boundedAdapter,
	}, nil
}

// ParseSessionCapabilitySnapshot validates ACPX's public acpx.session.v1
// projection and reuses the ACP configOptions validation contract, including
// its retention rule.
func ParseSessionCapabilitySnapshot(payload []byte, adapter AdapterEvidence, retention SelectionRetention) (SelectionCapabilities, error) {
	issues := newCapabilityIssueSet()
	if len(payload) == 0 {
		issues.add(CapabilityIssueMalformedResponse)
		return SelectionCapabilities{}, issues.err()
	}
	if len(payload) > maxCapabilityResponseBytes {
		issues.add(CapabilityIssueResponseTooLarge)
		return SelectionCapabilities{}, issues.err()
	}

	var snapshot acpxSessionCapabilityResponse
	if err := json.Unmarshal(payload, &snapshot); err != nil || snapshot.Schema != "acpx.session.v1" || snapshot.ACPX == nil || !boundedCapabilityValue(snapshot.ACPX.CurrentModelID) {
		issues.add(CapabilityIssueMalformedResponse)
		return SelectionCapabilities{}, issues.err()
	}
	if snapshot.ACPX.ConfigOptions == nil {
		issues.add(CapabilityIssueMissingOptions)
		return SelectionCapabilities{}, issues.err()
	}

	projection, err := json.Marshal(acpxCapabilityResponse{
		Action:        "config_set",
		ConfigID:      "model",
		Value:         snapshot.ACPX.CurrentModelID,
		ConfigOptions: snapshot.ACPX.ConfigOptions,
	})
	if err != nil {
		issues.add(CapabilityIssueMalformedResponse)
		return SelectionCapabilities{}, issues.err()
	}
	return ParseSessionConfigOptions(projection, adapter, retention)
}

// AcquireSelectionCapabilities uses only ACPX's public strict-JSON command
// boundary. It does not open ACPX Session files, Codex caches, or other
// runtime-private persistence.
func (runner ACPXRunner) AcquireSelectionCapabilities(ctx context.Context, request CapabilityAcquisitionRequest) (SelectionCapabilities, error) {
	return runner.acquireSelectionCapabilities(ctx, request, nil)
}

func (runner ACPXRunner) acquireSelectionCapabilities(ctx context.Context, request CapabilityAcquisitionRequest, env []string) (SelectionCapabilities, error) {
	if ctx == nil {
		return SelectionCapabilities{}, errors.New("capability acquisition context is required")
	}
	if err := ctx.Err(); err != nil {
		return SelectionCapabilities{}, err
	}
	if issues := capabilityAdapterIssues(request.Adapter); len(issues) > 0 {
		return SelectionCapabilities{}, newCapabilityEvidenceError(issues...)
	}
	args, err := acpxCapabilityResponseArgs(request.Runtime, request.ConfigID, request.Value, request.Session)
	if err != nil {
		return SelectionCapabilities{}, err
	}
	output, err := runner.runACPXCommandOutputWithEnv(ctx, args, env)
	if err != nil {
		if ctxErr := ctx.Err(); ctxErr != nil {
			return SelectionCapabilities{}, ctxErr
		}
		return SelectionCapabilities{}, &CapabilityAcquisitionError{Err: err}
	}
	capabilities, parseErr := ParseSessionConfigOptions([]byte(output), request.Adapter, RetentionFor(request.Runtime))
	if parseErr != nil {
		if request.ConfigID == "model" && exactACPXModelSetResponse([]byte(output), request.Value) {
			return runner.observeSelectionCapabilities(ctx, request.Runtime, request.Session, request.Adapter, env)
		}
		return SelectionCapabilities{}, parseErr
	}
	var response acpxCapabilityResponse
	if err := json.Unmarshal([]byte(output), &response); err != nil || response.ConfigID != request.ConfigID || response.Value != request.Value {
		return SelectionCapabilities{}, newCapabilityEvidenceError(CapabilityIssueContradictoryResponse)
	}
	return capabilities, nil
}

func (runner ACPXRunner) observeSelectionCapabilities(ctx context.Context, runtime RuntimeSpec, session SessionRef, adapter AdapterEvidence, env []string) (SelectionCapabilities, error) {
	if ctx == nil {
		return SelectionCapabilities{}, errors.New("capability observation context is required")
	}
	if err := ctx.Err(); err != nil {
		return SelectionCapabilities{}, err
	}
	if issues := capabilityAdapterIssues(adapter); len(issues) > 0 {
		return SelectionCapabilities{}, newCapabilityEvidenceError(issues...)
	}
	args, err := acpxSessionCapabilityArgs(runtime, session)
	if err != nil {
		return SelectionCapabilities{}, err
	}
	output, err := runner.runACPXCommandOutputWithEnv(ctx, args, env)
	if err != nil {
		if ctxErr := ctx.Err(); ctxErr != nil {
			return SelectionCapabilities{}, ctxErr
		}
		return SelectionCapabilities{}, &CapabilityAcquisitionError{Err: err}
	}
	return ParseSessionCapabilitySnapshot([]byte(output), adapter, RetentionFor(runtime))
}

func exactACPXModelSetResponse(payload []byte, expectedModel string) bool {
	var response acpxModelSetResponse
	if err := json.Unmarshal(payload, &response); err != nil {
		return false
	}
	return response.Action == "model_set" && response.ModelID == expectedModel
}

func acpxSessionCapabilityArgs(runtime RuntimeSpec, session SessionRef) ([]string, error) {
	workDir := strings.TrimSpace(session.WorkDir)
	if workDir == "" {
		return nil, errors.New("Agent working directory is required")
	}
	sessionName := strings.TrimSpace(session.Name)
	if sessionName == "" {
		return nil, errors.New("Agent Session name is required")
	}
	agentArgs, err := acpxAgentArgs(runtime)
	if err != nil {
		return nil, err
	}
	args := []string{"--cwd", workDir, "--format", "json", "--json-strict"}
	args = append(args, agentArgs...)
	return append(args, "sessions", "show", sessionName), nil
}

func acpxCapabilityResponseArgs(runtime RuntimeSpec, configID string, value string, session SessionRef) ([]string, error) {
	workDir := strings.TrimSpace(session.WorkDir)
	if workDir == "" {
		return nil, errors.New("Agent working directory is required")
	}
	sessionName := strings.TrimSpace(session.Name)
	if sessionName == "" {
		return nil, errors.New("Agent Session name is required")
	}
	configID = strings.TrimSpace(configID)
	if !boundedCapabilityID(configID) {
		return nil, errors.New("ACP config option id is required and must be bounded")
	}
	value = strings.TrimSpace(value)
	if !boundedCapabilityValue(value) {
		return nil, errors.New("ACP config option value is required and must be bounded")
	}
	agentArgs, err := acpxAgentArgs(runtime)
	if err != nil {
		return nil, err
	}
	args := []string{"--cwd", workDir, "--format", "json", "--json-strict"}
	args = append(args, agentArgs...)
	return append(args, "set", configID, value, "-s", sessionName), nil
}

func boundedConfigOptions(options []acpSelectOptionPayload) []acpSelectOptionPayload {
	if len(options) > maxCapabilityOptions {
		return options[:maxCapabilityOptions]
	}
	return options
}

func projectSelectCapability(raw acpSelectOptionPayload, retention SelectionRetention, issues *capabilityIssueSet) (SelectCapability, bool) {
	if raw.Type != "select" {
		return SelectCapability{}, false
	}
	currentValue, currentIsString := raw.CurrentValue.(string)
	if !currentIsString || !boundedCapabilityID(raw.ID) || !boundedCapabilityValue(currentValue) || raw.Options == nil || len(*raw.Options) == 0 {
		issues.add(CapabilityIssueInvalidOption)
		return SelectCapability{}, false
	}
	if len(*raw.Options) > maxAdvertisedCapabilityValues {
		issues.add(CapabilityIssueTooManyValues)
		return SelectCapability{}, false
	}

	advertised := make([]string, 0, len(*raw.Options))
	seen := make(map[string]struct{}, len(*raw.Options))
	currentAdvertised := false
	for _, rawValue := range *raw.Options {
		value := strings.TrimSpace(rawValue.Value)
		if !boundedCapabilityValue(value) {
			issues.add(CapabilityIssueInvalidOption)
			continue
		}
		if _, exists := seen[value]; exists {
			issues.add(CapabilityIssueDuplicateOptionValue)
			continue
		}
		seen[value] = struct{}{}
		advertised = append(advertised, value)
		if value == currentValue {
			currentAdvertised = true
		}
	}
	if !currentAdvertised {
		issues.add(CapabilityIssueInvalidCurrentValue)
	}
	return SelectCapability{
		ID:              raw.ID,
		CurrentValue:    currentValue,
		Values:          retainAdvertisedValues(advertised, currentValue, retention),
		AdvertisedCount: len(advertised),
	}, true
}

// retainAdvertisedValues bounds what the projection keeps rather than what the
// adapter may advertise. At or below the bound every value survives in
// advertised order, so a runtime with a small catalog projects exactly as it
// did before retention existed. Above it, the current value and every value the
// requested Agent Selection binds to are kept first, then advertised order
// fills the remainder for diagnostics. Retention narrows only what is kept:
// a value planning would have bound is never dropped, so an unadvertised
// Agent Model still fails as unsupported rather than as invalid evidence.
func retainAdvertisedValues(advertised []string, currentValue string, retention SelectionRetention) []string {
	if len(advertised) <= maxRetainedCapabilityValues {
		return advertised
	}

	keep := make(map[string]struct{}, maxRetainedCapabilityValues)
	for _, value := range advertised {
		if len(keep) >= maxRetainedCapabilityValues {
			break
		}
		if value == currentValue || value == retention.ReasoningEffort || bindsRequestedModel(value, retention.Model) {
			keep[value] = struct{}{}
		}
	}
	for _, value := range advertised {
		if len(keep) >= maxRetainedCapabilityValues {
			break
		}
		keep[value] = struct{}{}
	}

	retained := make([]string, 0, len(keep))
	for _, value := range advertised {
		if _, exists := keep[value]; exists {
			retained = append(retained, value)
		}
	}
	return retained
}

// bindsRequestedModel reports whether one advertised value represents the
// requested model, matching the way modelsForCanonical binds: the exact
// advertised value, or a value whose canonical prefix before a trailing
// bracketed effort equals the request.
func bindsRequestedModel(value string, requested string) bool {
	if requested == "" {
		return false
	}
	if value == requested {
		return true
	}
	if !strings.HasSuffix(value, "]") {
		return false
	}
	open := strings.LastIndexByte(value, '[')
	if open <= 0 {
		return false
	}
	return strings.TrimSpace(value[:open]) == requested
}

func parseModelCapability(adapterValue string, opaque bool) (ModelCapability, bool) {
	adapterValue = strings.TrimSpace(adapterValue)
	if !boundedCapabilityValue(adapterValue) {
		return ModelCapability{}, false
	}
	if !strings.Contains(adapterValue, "[") && !strings.Contains(adapterValue, "]") {
		return ModelCapability{AdapterValue: adapterValue, CanonicalModel: adapterValue, ModelManaged: true}, true
	}
	if !strings.HasSuffix(adapterValue, "]") {
		return ModelCapability{}, false
	}
	open := strings.LastIndexByte(adapterValue, '[')
	if open <= 0 {
		return ModelCapability{}, false
	}
	canonical := strings.TrimSpace(adapterValue[:open])
	effort := strings.TrimSpace(adapterValue[open+1 : len(adapterValue)-1])
	if canonical == "" || effort == "" || strings.ContainsAny(canonical, "[]") || strings.ContainsAny(effort, "[]") {
		return ModelCapability{}, false
	}
	if opaque {
		return ModelCapability{AdapterValue: adapterValue, CanonicalModel: canonical, ModelManaged: true}, true
	}
	return ModelCapability{AdapterValue: adapterValue, CanonicalModel: canonical, ReasoningEffort: effort}, true
}

func selectCapabilityIndex(options []SelectCapability, id string) int {
	for index := range options {
		if options[index].ID == id {
			return index
		}
	}
	return 0
}

func validateCapabilityAdapter(adapter AdapterEvidence, issues *capabilityIssueSet) AdapterEvidence {
	bounded := AdapterEvidence{
		Command: strings.TrimSpace(adapter.Command),
		Package: strings.TrimSpace(adapter.Package),
		Version: strings.TrimSpace(adapter.Version),
	}
	for _, issue := range capabilityAdapterIssues(bounded) {
		issues.add(issue)
	}
	return bounded
}

func capabilityAdapterIssues(adapter AdapterEvidence) []string {
	if strings.TrimSpace(adapter.Command) == "" {
		return []string{CapabilityIssueMissingAdapter}
	}
	for _, value := range []string{adapter.Command, adapter.Package, adapter.Version} {
		if len(value) > maxAdapterIdentityBytes {
			return []string{CapabilityIssueAdapterIdentityTooLong}
		}
	}
	return nil
}

func boundedCapabilityID(value string) bool {
	return value != "" && strings.TrimSpace(value) == value && len(value) <= maxCapabilityIDBytes
}

func boundedCapabilityValue(value string) bool {
	return value != "" && strings.TrimSpace(value) == value && len(value) <= maxCapabilityValueBytes
}

type capabilityIssueSet struct {
	values map[string]struct{}
}

func newCapabilityIssueSet() *capabilityIssueSet {
	return &capabilityIssueSet{values: make(map[string]struct{})}
}

func (issues *capabilityIssueSet) add(value string) {
	if value != "" {
		issues.values[value] = struct{}{}
	}
}

func (issues *capabilityIssueSet) err() error {
	if len(issues.values) == 0 {
		return nil
	}
	values := make([]string, 0, len(issues.values))
	for value := range issues.values {
		values = append(values, value)
	}
	return newCapabilityEvidenceError(values...)
}

func newCapabilityEvidenceError(values ...string) *CapabilityEvidenceError {
	unique := make(map[string]struct{}, len(values))
	for _, value := range values {
		if value != "" {
			unique[value] = struct{}{}
		}
	}
	bounded := make([]string, 0, len(unique))
	for value := range unique {
		bounded = append(bounded, value)
	}
	sort.Strings(bounded)
	if len(bounded) > MaxCapabilityDiagnosticIssues {
		bounded = bounded[:MaxCapabilityDiagnosticIssues]
	}
	return &CapabilityEvidenceError{Issues: bounded}
}
