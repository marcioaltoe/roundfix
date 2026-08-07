package cli

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"strings"

	"roundfix/internal/app"
	"roundfix/internal/baseline"
)

const baselineUpdateResultSchema = "roundfix/baseline-update-result/v1"

type baselineUpdateRequest struct {
	repo           string
	format         string
	confirmation   string
	yes            bool
	adoptSuggested bool
}

type baselineUpdateResult struct {
	SchemaVersion      string                        `json:"schemaVersion"`
	Operation          string                        `json:"operation"`
	State              string                        `json:"state"`
	Category           string                        `json:"category,omitempty"`
	Message            string                        `json:"message,omitempty"`
	NextAction         string                        `json:"nextAction,omitempty"`
	PriorCatalog       baseline.CatalogIdentity      `json:"priorCatalog"`
	CurrentCatalog     baseline.CatalogIdentity      `json:"currentCatalog"`
	FileChanges        []baseline.FileChange         `json:"fileChanges"`
	Retention          []baseline.RetentionEvidence  `json:"retention"`
	ClauseDelta        *baseline.ClauseDelta         `json:"clauseDelta,omitempty"`
	Warnings           []baseline.Finding            `json:"warnings"`
	NewDecisions       []baseline.DecisionSuggestion `json:"newDecisions"`
	AdoptedSuggestions []baseline.DecisionSuggestion `json:"adoptedSuggestions"`
	PlanDigest         string                        `json:"planDigest,omitempty"`
	ApprovedPlanDigest string                        `json:"approvedPlanDigest,omitempty"`
	StatusMatrix       *baseline.ResultStatusMatrix  `json:"statusMatrix,omitempty"`
}

func runBaselineUpdateCommand(
	ctx context.Context,
	args []string,
	stdout io.Writer,
	stderr io.Writer,
	environment commandEnvironment,
) int {
	if commandWantsHelp(args) {
		fmt.Fprint(stdout, commandUsage("baseline update"))
		return exitOK
	}
	request, err := parseBaselineUpdateCommand(args)
	jsonOutput := request.format == "json" || baselineUpdateJSONRequested(args)
	if err != nil {
		return writeBaselineUpdateFailure(
			newBaselineUpdateResult(),
			err,
			"invalid",
			"correct the command input and rerun roundfix baseline update",
			exitPreflight,
			jsonOutput,
			stdout,
			stderr,
		)
	}

	catalog, err := baseline.LoadEmbeddedCatalog()
	if err != nil {
		return writeBaselineUpdateFailure(
			newBaselineUpdateResult(), err, "execution",
			"repair the embedded Baseline catalog and rerun roundfix baseline update",
			exitRunFailed, jsonOutput, stdout, stderr,
		)
	}
	result := newBaselineUpdateResult()
	result.CurrentCatalog = baseline.CatalogIdentity{
		SchemaVersion: baseline.CatalogSchemaVersion(),
		Digest:        catalog.Digest(),
	}

	input, resolveErr := baseline.ResolveManifestInput(request.repo, catalog)
	if input.Manifest.CatalogDigest != "" {
		result.PriorCatalog = baseline.CatalogIdentity{
			SchemaVersion: baseline.CatalogSchemaVersion(),
			Digest:        input.Manifest.CatalogDigest,
		}
	}
	if errors.Is(resolveErr, baseline.ErrNoManifest) {
		result.State = "action_required"
		result.Category = "adoption"
		result.Message = "the repository has no Setup Manifest and cannot be updated"
		result.NextAction = "run roundfix baseline for first adoption, then rerun roundfix baseline update"
		return writeBaselineUpdateOutcome(result, exitUnverified, jsonOutput, stdout, stderr)
	}
	if resolveErr != nil {
		category := "manifest"
		nextAction := "repair the Setup Manifest or run roundfix baseline for first adoption"
		return writeBaselineUpdateFailure(
			result, resolveErr, category, nextAction, exitPreflight,
			jsonOutput, stdout, stderr,
		)
	}
	result.NewDecisions = append(result.NewDecisions, input.NewDecisions...)
	if len(input.NewDecisions) != 0 && !request.adoptSuggested {
		ids := baselineUpdateDecisionIDs(input.NewDecisions)
		result.State = "action_required"
		result.Category = "decision"
		result.Message = "the current Baseline catalog requires new decisions: " + strings.Join(ids, ", ")
		result.NextAction = "answer the named decisions through roundfix baseline, or rerun roundfix baseline update with --adopt-suggested"
		return writeBaselineUpdateOutcome(result, exitUnverified, jsonOutput, stdout, stderr)
	}
	if request.adoptSuggested {
		result.AdoptedSuggestions = append(result.AdoptedSuggestions, input.NewDecisions...)
		for _, suggestion := range input.NewDecisions {
			input.Decisions = append(input.Decisions, baseline.DecisionValue{
				ID:    suggestion.ID,
				Value: suggestion.SuggestedValue,
			})
		}
	}

	executableDirectories, err := environment.executableDirectories("resolve Baseline executable search path")
	if err != nil {
		return writeBaselineUpdateFailure(
			result, err, "invalid",
			"correct the executable search path and rerun roundfix baseline update",
			exitPreflight, jsonOutput, stdout, stderr,
		)
	}
	outcome, err := baseline.BuildPlan(ctx, baseline.PlanRequest{
		Repository:            request.repo,
		ProfileID:             input.ProfileID,
		Decisions:             input.Decisions,
		Preservation:          baseline.RootPreservationRequest{Mode: baseline.PreservationModeManagedRefresh},
		ExecutableDirectories: executableDirectories,
	})
	if err != nil {
		return writeBaselineUpdateOperationFailure(result, err, jsonOutput, stdout, stderr)
	}
	if outcome.Plan == nil {
		result.State = outcome.Result.State
		result.Category = outcome.Result.Category
		result.Message = outcome.Result.Message
		result.NextAction = outcome.Result.NextAction
		result.Warnings = append(result.Warnings, outcome.Result.Warnings...)
		result.ClauseDelta = outcome.Result.ClauseDelta
		exit := exitUnverified
		if outcome.Result.Category == "preflight" {
			exit = exitPreflight
		}
		return writeBaselineUpdateOutcome(result, exit, jsonOutput, stdout, stderr)
	}

	plan := *outcome.Plan
	result.PlanDigest = plan.PlanDigest
	result.CurrentCatalog = plan.Catalog
	result.FileChanges = append(result.FileChanges, plan.FileChanges...)
	result.Retention = append(result.Retention, plan.Retention...)
	result.ClauseDelta = plan.ClauseDelta
	result.Warnings = append(result.Warnings, plan.Warnings...)
	if !request.yes && request.confirmation == "" {
		if len(plan.FileChanges) == 0 {
			result.State = "current"
			result.Message = "the repository already matches the current Baseline catalog"
			return writeBaselineUpdateOutcome(result, exitOK, jsonOutput, stdout, stderr)
		}
		result.State = "plan_ready"
		result.Category = "approval"
		result.Message = "managed-refresh plan is ready; repository bytes are unchanged"
		result.NextAction = "review the plan and rerun with --confirm-plan " + plan.PlanDigest + ", or use --yes to approve the digest computed in one invocation"
		return writeBaselineUpdateOutcome(result, exitUnverified, jsonOutput, stdout, stderr)
	}

	confirmation := request.confirmation
	if request.yes {
		confirmation = plan.PlanDigest
	}
	applyResult, err := baseline.ApplyPlan(ctx, request.repo, plan, confirmation)
	if err != nil {
		return writeBaselineUpdateApplyFailure(result, err, jsonOutput, stdout, stderr)
	}
	result.State = applyResult.State
	result.Message = applyResult.Message
	result.NextAction = applyResult.NextAction
	result.ApprovedPlanDigest = applyResult.PlanDigest
	result.StatusMatrix = applyResult.StatusMatrix
	return writeBaselineUpdateOutcome(result, exitOK, jsonOutput, stdout, stderr)
}

func parseBaselineUpdateCommand(args []string) (baselineUpdateRequest, error) {
	request := baselineUpdateRequest{repo: ".", format: "text"}
	flags := flag.NewFlagSet("baseline update", flag.ContinueOnError)
	flags.SetOutput(io.Discard)
	flags.StringVar(&request.repo, "repo", ".", "Git worktree or a path inside it")
	flags.StringVar(&request.format, "format", "text", "Output format: text or json")
	flags.StringVar(&request.confirmation, "confirm-plan", "", "Exact Plan Digest reviewed in a previous invocation")
	flags.BoolVar(&request.yes, "yes", false, "Approve the Plan Digest computed in this invocation")
	flags.BoolVar(&request.adoptSuggested, "adopt-suggested", false, "Adopt catalog suggestions for decisions absent from the Setup Manifest")
	if err := flags.Parse(args); err != nil {
		return baselineUpdateRequest{}, validationError{
			message: fmt.Sprintf("invalid baseline update arguments: %v; run '%s baseline update --help' for usage", err, app.Name),
		}
	}
	if remaining := flags.Args(); len(remaining) != 0 {
		return baselineUpdateRequest{}, validationError{
			message: fmt.Sprintf("unexpected argument %q; run '%s baseline update --help' for usage", remaining[0], app.Name),
		}
	}
	request.repo = strings.TrimSpace(request.repo)
	request.format = strings.TrimSpace(request.format)
	request.confirmation = strings.TrimSpace(request.confirmation)
	if request.repo == "" {
		return request, validationError{message: "--repo cannot be empty"}
	}
	if request.format != "text" && request.format != "json" {
		return request, validationError{
			message: fmt.Sprintf("unsupported --format %q; use text or json", request.format),
		}
	}
	if request.yes && request.confirmation != "" {
		return request, validationError{message: "--yes and --confirm-plan are mutually exclusive"}
	}
	return request, nil
}

func newBaselineUpdateResult() baselineUpdateResult {
	return baselineUpdateResult{
		SchemaVersion:      baselineUpdateResultSchema,
		Operation:          "update",
		FileChanges:        []baseline.FileChange{},
		Retention:          []baseline.RetentionEvidence{},
		Warnings:           []baseline.Finding{},
		NewDecisions:       []baseline.DecisionSuggestion{},
		AdoptedSuggestions: []baseline.DecisionSuggestion{},
	}
}

func baselineUpdateDecisionIDs(suggestions []baseline.DecisionSuggestion) []string {
	ids := make([]string, len(suggestions))
	for index, suggestion := range suggestions {
		ids[index] = suggestion.ID
	}
	return ids
}

func writeBaselineUpdateOperationFailure(
	result baselineUpdateResult,
	err error,
	jsonOutput bool,
	stdout io.Writer,
	stderr io.Writer,
) int {
	if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
		return writeBaselineUpdateFailure(
			result, err, "execution",
			"rerun roundfix baseline update when the operation can complete",
			exitSIGINT, jsonOutput, stdout, stderr,
		)
	}
	return writeBaselineUpdateFailure(
		result, err, "invalid",
		"repair the repository or manifest input and rerun roundfix baseline update",
		exitPreflight, jsonOutput, stdout, stderr,
	)
}

func writeBaselineUpdateApplyFailure(
	result baselineUpdateResult,
	err error,
	jsonOutput bool,
	stdout io.Writer,
	stderr io.Writer,
) int {
	category := string(baseline.ApplyErrorInvalid)
	nextAction := "correct the command input and rerun roundfix baseline update"
	exit := exitPreflight
	state := "failed"
	var applyErr *baseline.ApplyError
	if errors.As(err, &applyErr) {
		category = string(applyErr.Kind)
		nextAction = applyErr.NextAction
		switch applyErr.Kind {
		case baseline.ApplyErrorApproval, baseline.ApplyErrorStale:
			state = "action_required"
			exit = exitUnverified
		case baseline.ApplyErrorExecution, baseline.ApplyErrorVerification:
			exit = exitRunFailed
		}
	}
	if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
		category = string(baseline.ApplyErrorExecution)
		nextAction = "rerun roundfix baseline update when the operation can complete"
		exit = exitSIGINT
	}
	result.State = state
	return writeBaselineUpdateFailure(
		result, err, category, nextAction, exit, jsonOutput, stdout, stderr,
	)
}

func writeBaselineUpdateFailure(
	result baselineUpdateResult,
	err error,
	category string,
	nextAction string,
	exit int,
	jsonOutput bool,
	stdout io.Writer,
	stderr io.Writer,
) int {
	if result.State == "" {
		result.State = "failed"
	}
	result.Category = category
	result.Message = err.Error()
	result.NextAction = nextAction
	fmt.Fprintf(stderr, "%s: baseline update failed: %v\n", app.Name, err)
	return writeBaselineUpdateOutcome(result, exit, jsonOutput, stdout, stderr)
}

func writeBaselineUpdateOutcome(
	result baselineUpdateResult,
	exit int,
	jsonOutput bool,
	stdout io.Writer,
	stderr io.Writer,
) int {
	if err := writeBaselineUpdateResult(result, jsonOutput, stdout); err != nil {
		fmt.Fprintf(stderr, "%s: baseline update output failed: %v\n", app.Name, err)
		return exitRunFailed
	}
	return exit
}

func writeBaselineUpdateResult(result baselineUpdateResult, jsonOutput bool, stdout io.Writer) error {
	if jsonOutput {
		return json.NewEncoder(stdout).Encode(result)
	}
	fmt.Fprintf(stdout, "Baseline update: %s\n", strings.ReplaceAll(result.State, "_", " "))
	if result.Category != "" {
		fmt.Fprintf(stdout, "Category: %s\n", result.Category)
	}
	if result.Message != "" {
		fmt.Fprintf(stdout, "Result: %s\n", result.Message)
	}
	if result.PriorCatalog.Digest == "" {
		fmt.Fprintln(stdout, "Prior catalog: unavailable")
	} else {
		fmt.Fprintf(stdout, "Prior catalog: %s\n", result.PriorCatalog.Digest)
	}
	if result.CurrentCatalog.Digest != "" {
		fmt.Fprintf(stdout, "Current catalog: %s\n", result.CurrentCatalog.Digest)
	}
	fmt.Fprintf(stdout, "File changes: %d\n", len(result.FileChanges))
	for _, change := range result.FileChanges {
		fmt.Fprintf(stdout, "- %s %s (%d managed entries)\n", change.Action, change.Path, len(change.ManagedEntries))
	}
	fmt.Fprintf(stdout, "Retention evidence: %d\n", len(result.Retention))
	for _, evidence := range result.Retention {
		fmt.Fprintf(stdout, "- %s: %s", evidence.FromClause, evidence.Disposition)
		if len(evidence.Targets) != 0 {
			fmt.Fprintf(stdout, " -> %s", strings.Join(evidence.Targets, ", "))
		}
		fmt.Fprintln(stdout)
	}
	if result.ClauseDelta != nil {
		data, err := json.Marshal(result.ClauseDelta.Counts)
		if err != nil {
			return fmt.Errorf("encode clause delta counts: %w", err)
		}
		fmt.Fprintf(stdout, "Clause delta: %s\n", data)
	}
	for _, decision := range result.NewDecisions {
		value, err := json.Marshal(decision.SuggestedValue)
		if err != nil {
			return fmt.Errorf("encode new decision %q: %w", decision.ID, err)
		}
		fmt.Fprintf(stdout, "New decision: %s suggested=%s (%s)\n", decision.ID, value, decision.Summary)
	}
	for _, decision := range result.AdoptedSuggestions {
		value, err := json.Marshal(decision.SuggestedValue)
		if err != nil {
			return fmt.Errorf("encode adopted decision %q: %w", decision.ID, err)
		}
		fmt.Fprintf(stdout, "Adopted suggestion: %s=%s\n", decision.ID, value)
	}
	for _, warning := range result.Warnings {
		fmt.Fprintf(stdout, "Warning: %s: %s: %s\n", warning.Code, warning.Path, warning.Message)
	}
	if result.PlanDigest != "" {
		fmt.Fprintf(stdout, "Plan Digest: %s\n", result.PlanDigest)
	}
	if result.ApprovedPlanDigest != "" {
		fmt.Fprintf(stdout, "Approved Plan Digest: %s\n", result.ApprovedPlanDigest)
	}
	if result.NextAction != "" {
		fmt.Fprintf(stdout, "Next action: %s\n", result.NextAction)
	}
	return nil
}

func baselineUpdateJSONRequested(args []string) bool {
	for index, arg := range args {
		if arg == "--format=json" || arg == "-format=json" {
			return true
		}
		if (arg == "--format" || arg == "-format") &&
			index+1 < len(args) && strings.TrimSpace(args[index+1]) == "json" {
			return true
		}
	}
	return false
}
