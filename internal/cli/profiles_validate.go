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
	Selection  roundconfig.AgentSelection `json:"selection"`
	Status     string                     `json:"status"`
	References []profileProofReference    `json:"references"`
	Error      string                     `json:"error,omitempty"`
}

type profileProofReference struct {
	Category      roundconfig.WorkCategory  `json:"category"`
	Source        roundconfig.ProfileSource `json:"source"`
	InheritedFrom roundconfig.WorkCategory  `json:"inherited_from"`
	Role          string                    `json:"role"`
	FallbackIndex int                       `json:"fallback_index,omitempty"`
}

type profileProofResult struct {
	Proofs []profileProofReport
	Err    error
}

type profileProofError struct {
	Selection  roundconfig.AgentSelection
	References []profileProofReference
	Err        error
}

func (err profileProofError) Error() string {
	var builder strings.Builder
	fmt.Fprintf(&builder, "profile proof failed for runtime %q, model %q, reasoning_effort %q", err.Selection.Runtime, err.Selection.Model, err.Selection.ReasoningEffort)
	fmt.Fprintf(&builder, "; affected categories: %s", formatProfileProofReferences(err.References))
	if err.Err != nil {
		fmt.Fprintf(&builder, "; adapter error: %v", err.Err)
	}
	builder.WriteString("; next: update the profile with `roundfix profiles configure --scope user|project` or rerun `roundfix profiles validate` after fixing the ACP Runtime")
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
	printProfilesValidateSuccess(req, result, stdout)
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
	if ctx == nil {
		ctx = context.Background()
	}
	proofs, err := buildProfileProofReports(config, categories)
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
		runtime, err := runtimeForProfileSelection(proofs[index].Selection)
		if err != nil {
			proofs[index].Status = "failed"
			proofs[index].Error = err.Error()
			return profileProofResult{Proofs: proofs, Err: profileProofError{Selection: proofs[index].Selection, References: proofs[index].References, Err: err}}
		}
		if err := runner.Probe(ctx, agent.ProbeRequest{Runtime: runtime, WorkDir: workDir}); err != nil {
			proofs[index].Status = "failed"
			proofs[index].Error = err.Error()
			return profileProofResult{Proofs: proofs, Err: profileProofError{Selection: proofs[index].Selection, References: proofs[index].References, Err: err}}
		}
		proofs[index].Status = "passed"
	}
	return profileProofResult{Proofs: proofs}
}

func buildProfileProofReports(config roundconfig.Config, categories []roundconfig.WorkCategory) ([]profileProofReport, error) {
	proofs := []profileProofReport{}
	bySelection := map[roundconfig.AgentSelection]int{}
	for _, category := range categories {
		resolved, err := roundconfig.ResolveProfile(config, category, nil)
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
	return agent.RuntimeFor(agent.RuntimeOptions{
		Agent:           strings.TrimSpace(selection.Runtime),
		Model:           strings.TrimSpace(selection.Model),
		ReasoningEffort: strings.TrimSpace(selection.ReasoningEffort),
	})
}

func printProfilesValidateSuccess(req profilesValidateRequest, result profileProofResult, stdout io.Writer) {
	if req.json {
		_ = json.NewEncoder(stdout).Encode(profilesValidateResponseForResult(result, ""))
		return
	}
	fmt.Fprintln(stdout, "Profiles validate passed.")
	for index, proof := range result.Proofs {
		fmt.Fprintf(stdout, "%d. %s — %s\n", index+1, formatProfileSelection(proof.Selection), proof.Status)
		for _, reference := range proof.References {
			fmt.Fprintf(stdout, "   - %s\n", formatProfileProofReference(reference))
		}
	}
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
