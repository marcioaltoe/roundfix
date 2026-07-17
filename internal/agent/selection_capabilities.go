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
	maxCapabilityResponseBytes    = 64 * 1024
	maxCapabilityOptions          = 32
	maxCapabilityValues           = 64
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
// for exact assignment and proof.
type SelectCapability struct {
	ID           string
	CurrentValue string
	Values       []string
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
// projection of a documented ACP configOptions response.
func ParseSessionConfigOptions(payload []byte, adapter AdapterEvidence) (SelectionCapabilities, error) {
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
		option, ok := projectSelectCapability(rawOption, issues)
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

	modelOption := selectOptions[selectCapabilityIndex(selectOptions, modelOptionIDs[0])]
	models := make([]ModelCapability, 0, len(modelOption.Values))
	seenModels := make(map[string]struct{}, len(modelOption.Values))
	for _, value := range modelOption.Values {
		model, ok := parseModelCapability(value)
		if !ok {
			issues.add(CapabilityIssueMalformedModelValue)
			continue
		}
		key := model.CanonicalModel + "\x00" + model.ReasoningEffort
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

	var reasoning *SelectCapability
	for index := range selectOptions {
		if selectOptions[index].ID == acpxCodexReasoningEffortKey || selectOptions[index].ID == acpxGenericReasoningEffortKey {
			copyOption := selectOptions[index]
			copyOption.Values = append([]string(nil), copyOption.Values...)
			reasoning = &copyOption
			break
		}
	}

	return SelectionCapabilities{
		CurrentModel:    modelOption.CurrentValue,
		Models:          models,
		Options:         selectOptions,
		ReasoningOption: reasoning,
		Adapter:         boundedAdapter,
	}, nil
}

// AcquireSelectionCapabilities uses only ACPX's public strict-JSON command
// boundary. It does not open ACPX Session files, Codex caches, or other
// runtime-private persistence.
func (runner ACPXRunner) AcquireSelectionCapabilities(ctx context.Context, request CapabilityAcquisitionRequest) (SelectionCapabilities, error) {
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
	output, err := runner.runACPXCommandOutput(ctx, args)
	if err != nil {
		if ctxErr := ctx.Err(); ctxErr != nil {
			return SelectionCapabilities{}, ctxErr
		}
		return SelectionCapabilities{}, &CapabilityAcquisitionError{Err: err}
	}
	return ParseSessionConfigOptions([]byte(output), request.Adapter)
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

func projectSelectCapability(raw acpSelectOptionPayload, issues *capabilityIssueSet) (SelectCapability, bool) {
	if raw.Type != "select" {
		return SelectCapability{}, false
	}
	currentValue, currentIsString := raw.CurrentValue.(string)
	if !currentIsString || !boundedCapabilityID(raw.ID) || !boundedCapabilityValue(currentValue) || raw.Options == nil || len(*raw.Options) == 0 {
		issues.add(CapabilityIssueInvalidOption)
		return SelectCapability{}, false
	}
	if len(*raw.Options) > maxCapabilityValues {
		issues.add(CapabilityIssueTooManyValues)
		return SelectCapability{}, false
	}

	values := make([]string, 0, len(*raw.Options))
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
		values = append(values, value)
		if value == currentValue {
			currentAdvertised = true
		}
	}
	if !currentAdvertised {
		issues.add(CapabilityIssueInvalidCurrentValue)
	}
	return SelectCapability{ID: raw.ID, CurrentValue: currentValue, Values: values}, true
}

func parseModelCapability(adapterValue string) (ModelCapability, bool) {
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
