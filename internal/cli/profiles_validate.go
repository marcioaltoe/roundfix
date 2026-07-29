package cli

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"roundfix/internal/agent"
	roundconfig "roundfix/internal/config"
)

const profilesValidateSchema = "roundfix/profiles-validate/v1"

type profilesValidateRequest struct {
	category string
	json     bool
}

type profilesValidateResponse struct {
	Schema string               `json:"schema"`
	OK     bool                 `json:"ok"`
	Proofs []profileProofReport `json:"proofs"`
	Error  string               `json:"error,omitempty"`
}

type profileProofReport struct {
	Selection           roundconfig.AgentSelection `json:"selection"`
	Status              string                     `json:"status"`
	References          []profileProofReference    `json:"references"`
	Classification      string                     `json:"classification,omitempty"`
	Encoding            string                     `json:"encoding,omitempty"`
	AdapterCommand      string                     `json:"adapter_command,omitempty"`
	AdapterVersion      string                     `json:"adapter_version,omitempty"`
	AdvertisedModels    []string                   `json:"advertised_models,omitempty"`
	AdvertisedReasoning []string                   `json:"advertised_reasoning,omitempty"`
	NextAction          string                     `json:"next_action,omitempty"`
	Error               string                     `json:"error,omitempty"`
}

type profileProofReference struct {
	Category      roundconfig.WorkCategory  `json:"category"`
	Source        roundconfig.ProfileSource `json:"source"`
	InheritedFrom roundconfig.WorkCategory  `json:"inherited_from"`
	Role          string                    `json:"role"`
	FallbackIndex int                       `json:"fallback_index,omitempty"`
}

type profileReadiness struct {
	Proofs []profileProofReport
	Err    error
}

type profileProofResult = profileReadiness

type profileProofOptions struct {
	PreferredOverride *roundconfig.AgentSelection
	CommandOverride   string
	CommandOverrides  map[string]string
	EnableFullAccess  bool
}

type profileProofError struct {
	Selection      roundconfig.AgentSelection
	References     []profileProofReference
	Classification string
	NextAction     string
	Err            error
}

func (err profileProofError) Error() string {
	var builder strings.Builder
	fmt.Fprintf(&builder, "profile proof failed for runtime %q, model %q, reasoning_effort %q", err.Selection.Runtime, err.Selection.Model, err.Selection.ReasoningEffort)
	fmt.Fprintf(&builder, "; affected categories: %s", formatProfileProofReferences(err.References))
	if err.Classification != "" {
		fmt.Fprintf(&builder, "; classification: %s", err.Classification)
	}
	if err.Err != nil {
		fmt.Fprintf(&builder, "; adapter error: %v", err.Err)
	}
	if err.NextAction != "" {
		fmt.Fprintf(&builder, "; next: %s", err.NextAction)
	}
	fmt.Fprint(&builder, "; fallback: Fallback Chains activate only after Run creation (ADR-0050); Preflight proves every configured tuple and substitutes none")
	return builder.String()
}

func (err profileProofError) Unwrap() error {
	return err.Err
}

func runProfilesValidateCommand(ctx context.Context, args []string, stdout, stderr io.Writer) int {
	if commandWantsHelp(args) {
		fmt.Fprint(stdout, commandUsage("profiles validate"))
		return exitOK
	}
	req, err := parseProfilesValidateCommand(args)
	if err != nil {
		return printProfilesValidateError(req, profileProofResult{}, err, stdout, stderr)
	}
	categories, err := profilesValidateCategories(req)
	if err != nil {
		return printProfilesValidateError(req, profileProofResult{}, err, stdout, stderr)
	}
	loaded, err := roundconfig.Load(roundconfig.LoadOptions{Stderr: stderr})
	if err != nil {
		return printProfilesValidateError(req, profileProofResult{}, err, stdout, stderr)
	}
	workDir := loaded.GitRoot
	if strings.TrimSpace(workDir) == "" {
		workDir, err = os.Getwd()
		if err != nil {
			return printProfilesValidateError(req, profileProofResult{}, fmt.Errorf("resolve validation working directory: %w", err), stdout, stderr)
		}
	}
	result := proveProfileSelections(ctx, loaded.Config, categories, workDir, newEngineCollaborators().runner)
	if result.Err != nil {
		return printProfilesValidateError(req, result, result.Err, stdout, stderr)
	}
	if err := printProfilesValidateSuccess(req, result, stdout); err != nil {
		return printProfilesValidateOutputError(err, stderr)
	}
	return exitOK
}

func parseProfilesValidateCommand(args []string) (profilesValidateRequest, error) {
	req := profilesValidateRequest{}
	fs := flagSet("profiles validate")
	fs.StringVar(&req.category, "category", "", "Agent Work Category to validate")
	fs.BoolVar(&req.json, "json", false, "Print roundfix/profiles-validate/v1 JSON")
	if err := fs.Parse(args); err != nil {
		return req, validationError{message: err.Error()}
	}
	if remaining := fs.Args(); len(remaining) > 0 {
		return req, validationError{message: fmt.Sprintf("unexpected argument %q", remaining[0])}
	}
	req.category = strings.TrimSpace(req.category)
	return req, nil
}

func profilesValidateCategories(req profilesValidateRequest) ([]roundconfig.WorkCategory, error) {
	if req.category == "" {
		return roundconfig.RequiredWorkCategories(), nil
	}
	category, ok := roundconfig.ParseWorkCategory(req.category)
	if !ok {
		return nil, validationError{message: fmt.Sprintf("unknown profile category %q; supported values: %s", req.category, profilesCategoryList())}
	}
	return []roundconfig.WorkCategory{category}, nil
}

func proveProfileSelections(ctx context.Context, config roundconfig.Config, categories []roundconfig.WorkCategory, workDir string, runner agent.Runner) profileProofResult {
	return proveProfileSelectionsWithOptions(ctx, config, categories, workDir, runner, profileProofOptions{})
}

func proveProfileSelectionsWithOptions(ctx context.Context, config roundconfig.Config, categories []roundconfig.WorkCategory, workDir string, runner agent.Runner, options profileProofOptions) profileProofResult {
	if ctx == nil {
		ctx = context.Background()
	}
	proofs, err := buildProfileProofReportsWithOptions(config, categories, options)
	if err != nil {
		return profileProofResult{Err: err}
	}
	if runner == nil {
		return profileProofResult{Proofs: proofs, Err: errors.New("agent runner is required for profile validation")}
	}
	for index := range proofs {
		if err := ctx.Err(); err != nil {
			return profileProofResult{Proofs: proofs, Err: err}
		}
		runtime, err := runtimeForProfileSelectionWithOptions(proofs[index].Selection, options)
		if err != nil {
			applyProfileProofFailure(&proofs[index], err)
			return profileProofResult{Proofs: proofs, Err: profileProofError{
				Selection:      proofs[index].Selection,
				References:     proofs[index].References,
				Classification: proofs[index].Classification,
				NextAction:     proofs[index].NextAction,
				Err:            err,
			}}
		}
		proof, err := proveProfileSelection(ctx, runner, agent.ProbeRequest{Runtime: runtime, WorkDir: workDir})
		if err != nil {
			applyProfileProofFailure(&proofs[index], err)
			return profileProofResult{Proofs: proofs, Err: profileProofError{
				Selection:      proofs[index].Selection,
				References:     proofs[index].References,
				Classification: proofs[index].Classification,
				NextAction:     proofs[index].NextAction,
				Err:            err,
			}}
		}
		proofs[index].Status = "passed"
		proofs[index].Encoding = strings.TrimSpace(proof.Assignment.Encoding)
		proofs[index].AdapterCommand = strings.TrimSpace(proof.Adapter.Command)
		proofs[index].AdapterVersion = strings.TrimSpace(proof.Adapter.Version)
	}
	return profileProofResult{Proofs: proofs}
}

func proveProfileSelection(ctx context.Context, runner agent.Runner, request agent.ProbeRequest) (agent.SelectionProof, error) {
	if prover, ok := runner.(agent.SelectionProver); ok {
		return prover.ProveProfileSelection(ctx, request)
	}
	if err := runner.Probe(ctx, request); err != nil {
		return agent.SelectionProof{}, err
	}
	return agent.SelectionProof{}, nil
}

const (
	maxProfileProofAdvertisedValues = 16
	maxProfileProofValueBytes       = 160
)

func applyProfileProofFailure(report *profileProofReport, err error) {
	report.Status = "failed"
	report.Error = err.Error()
	report.Classification = profileProofClassification(err)
	report.NextAction = profileProofNextAction(err)

	var unsupported *agent.SelectionUnsupportedError
	if errors.As(err, &unsupported) {
		report.AdvertisedModels = boundedProfileProofValues(unsupported.AdvertisedModels)
		report.AdvertisedReasoning = boundedProfileProofValues(unsupported.AdvertisedReasoning)
	}
	var rejected *agent.SelectionRejectedError
	if errors.As(err, &rejected) {
		report.Encoding = strings.TrimSpace(rejected.Assignment.Encoding)
	}
	var mismatch *agent.EffectiveSelectionError
	if errors.As(err, &mismatch) {
		report.Encoding = strings.TrimSpace(mismatch.Assignment.Encoding)
	}
	var lineage *agent.AdapterLineageError
	if errors.As(err, &lineage) {
		report.AdapterCommand = strings.TrimSpace(lineage.Command)
		report.AdapterVersion = strings.TrimSpace(lineage.Version)
	}
	var version *agent.AdapterVersionError
	if errors.As(err, &version) {
		report.AdapterCommand = strings.TrimSpace(version.Command)
		report.AdapterVersion = strings.TrimSpace(version.FoundVersion)
	}
}

func profileProofClassification(err error) string {
	if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
		return ""
	}
	var classified interface{ Classification() string }
	if errors.As(err, &classified) {
		return strings.TrimSpace(classified.Classification())
	}
	return agent.SelectionRejected
}

func profileProofNextAction(err error) string {
	const configureAction = "update the profile with `roundfix profiles configure --scope user|project`"
	var installer interface{ InstallCommand() string }
	if errors.As(err, &installer) {
		if command := strings.TrimSpace(installer.InstallCommand()); command != "" {
			return "run `" + command + "`, then rerun `roundfix profiles validate`; if the tuple remains unavailable, " + configureAction
		}
	}
	var cleanup *agent.AgentSessionCleanupError
	if errors.As(err, &cleanup) {
		return "restore Agent Session cleanup, then rerun `roundfix profiles validate`; if the tuple remains unavailable, " + configureAction
	}
	var unsupported *agent.SelectionUnsupportedError
	if errors.As(err, &unsupported) {
		return "update the ACP Runtime or adapter, or " + configureAction + " with a different exact Agent Selection, then rerun `roundfix profiles validate`"
	}
	return configureAction + " or update the ACP Runtime or adapter, then rerun `roundfix profiles validate`"
}

func boundedProfileProofValues(values []string) []string {
	bounded := make([]string, 0, min(len(values), maxProfileProofAdvertisedValues))
	for _, value := range values {
		value = strings.TrimSpace(value)
		if value == "" || len(value) > maxProfileProofValueBytes {
			continue
		}
		bounded = append(bounded, value)
		if len(bounded) == maxProfileProofAdvertisedValues {
			break
		}
	}
	return bounded
}

func buildProfileProofReports(config roundconfig.Config, categories []roundconfig.WorkCategory) ([]profileProofReport, error) {
	return buildProfileProofReportsWithOptions(config, categories, profileProofOptions{})
}

func buildProfileProofReportsWithOptions(config roundconfig.Config, categories []roundconfig.WorkCategory, options profileProofOptions) ([]profileProofReport, error) {
	proofs := []profileProofReport{}
	bySelection := map[roundconfig.AgentSelection]int{}
	for _, category := range categories {
		resolved, err := roundconfig.ResolveProfile(config, category, options.PreferredOverride)
		if err != nil {
			return nil, err
		}
		addProfileProofReference(&proofs, bySelection, resolved.Profile.Preferred, profileProofReference{
			Category:      category,
			Source:        resolved.Source,
			InheritedFrom: resolved.InheritedFrom,
			Role:          "preferred",
		})
		for index, fallback := range resolved.Profile.Fallbacks {
			addProfileProofReference(&proofs, bySelection, fallback, profileProofReference{
				Category:      category,
				Source:        resolved.Source,
				InheritedFrom: resolved.InheritedFrom,
				Role:          "fallback",
				FallbackIndex: index + 1,
			})
		}
	}
	return proofs, nil
}

func addProfileProofReference(proofs *[]profileProofReport, bySelection map[roundconfig.AgentSelection]int, selection roundconfig.AgentSelection, reference profileProofReference) {
	if index, ok := bySelection[selection]; ok {
		(*proofs)[index].References = append((*proofs)[index].References, reference)
		return
	}
	bySelection[selection] = len(*proofs)
	*proofs = append(*proofs, profileProofReport{
		Selection:  selection,
		Status:     "pending",
		References: []profileProofReference{reference},
	})
}

func runtimeForProfileSelection(selection roundconfig.AgentSelection) (agent.RuntimeSpec, error) {
	return runtimeForProfileSelectionWithOptions(selection, profileProofOptions{})
}

func runtimeForProfileSelectionWithOptions(selection roundconfig.AgentSelection, options profileProofOptions) (agent.RuntimeSpec, error) {
	commandOverride := options.CommandOverride
	if command := strings.TrimSpace(options.CommandOverrides[strings.TrimSpace(selection.Runtime)]); command != "" {
		commandOverride = command
	}
	return agent.RuntimeFor(agent.RuntimeOptions{
		Agent:            strings.TrimSpace(selection.Runtime),
		CommandOverride:  commandOverride,
		Model:            strings.TrimSpace(selection.Model),
		ReasoningEffort:  strings.TrimSpace(selection.ReasoningEffort),
		EnableFullAccess: options.EnableFullAccess,
	})
}

func printProfilesValidateSuccess(req profilesValidateRequest, result profileProofResult, stdout io.Writer) error {
	if req.json {
		return json.NewEncoder(stdout).Encode(profilesValidateResponseForResult(result, ""))
	}
	fmt.Fprintln(stdout, "Profiles validate passed.")
	for index, proof := range result.Proofs {
		fmt.Fprintf(stdout, "%d. %s — %s\n", index+1, formatProfileSelection(proof.Selection), proof.Status)
		for _, reference := range proof.References {
			fmt.Fprintf(stdout, "   - %s\n", formatProfileProofReference(reference))
		}
	}
	return nil
}

func printProfilesValidateOutputError(err error, stderr io.Writer) int {
	printProfilesFailure(fmt.Errorf("encode profiles validate JSON: %w", err), stderr)
	return exitRunFailed
}

func printProfilesValidateError(req profilesValidateRequest, result profileProofResult, err error, stdout, stderr io.Writer) int {
	if req.json {
		_ = json.NewEncoder(stdout).Encode(profilesValidateResponseForResult(result, err.Error()))
	}
	printProfilesFailure(err, stderr)
	return exitPreflight
}

func profilesValidateResponseForResult(result profileProofResult, errorMessage string) profilesValidateResponse {
	return profilesValidateResponse{
		Schema: profilesValidateSchema,
		OK:     errorMessage == "",
		Proofs: result.Proofs,
		Error:  errorMessage,
	}
}

func formatProfileProofReferences(references []profileProofReference) string {
	values := make([]string, 0, len(references))
	for _, reference := range references {
		values = append(values, formatProfileProofReference(reference))
	}
	return strings.Join(values, ", ")
}

func formatProfileProofReference(reference profileProofReference) string {
	label := string(reference.Category) + " " + reference.Role
	if reference.Role == "fallback" {
		label = fmt.Sprintf("%s fallback[%d]", reference.Category, reference.FallbackIndex)
	}
	if reference.InheritedFrom != "" {
		label += " inherited_from=" + string(reference.InheritedFrom)
	}
	if reference.Source != "" {
		label += " source=" + string(reference.Source)
	}
	return label
}
