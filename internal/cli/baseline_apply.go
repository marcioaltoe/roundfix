package cli

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"roundfix/internal/app"
	"roundfix/internal/baseline"
)

type baselineApplyCommandRequest struct {
	repo         string
	plan         string
	confirmation string
	format       string
}

func runBaselineApplyCommand(ctx context.Context, args []string, stdout, stderr io.Writer) int {
	if commandWantsHelp(args) {
		fmt.Fprint(stdout, commandUsage("baseline apply"))
		return exitOK
	}
	request, err := parseBaselineApplyCommand(args)
	jsonOutput := request.format == "json" || baselineApplyJSONRequested(args)
	if err != nil {
		return printBaselineApplyFailure(baseline.PlanDocument{}, err, jsonOutput, stdout, stderr)
	}
	data, err := os.ReadFile(request.plan)
	if err != nil {
		return printBaselineApplyFailure(
			baseline.PlanDocument{},
			fmt.Errorf("read Baseline Plan %q: %w", request.plan, err),
			jsonOutput,
			stdout,
			stderr,
		)
	}
	document, err := baseline.ParsePlanDocument(data)
	if err != nil {
		return printBaselineApplyFailure(document, err, jsonOutput, stdout, stderr)
	}
	result, err := baseline.ApplyPlan(ctx, request.repo, document, request.confirmation)
	if err != nil {
		return printBaselineApplyFailure(document, err, jsonOutput, stdout, stderr)
	}
	if err := writeBaselineApplyResult(result, jsonOutput, stdout); err != nil {
		fmt.Fprintf(stderr, "%s: baseline apply output failed: %v\n", app.Name, err)
		return exitRunFailed
	}
	return exitOK
}

func parseBaselineApplyCommand(args []string) (baselineApplyCommandRequest, error) {
	request := baselineApplyCommandRequest{repo: ".", format: "text"}
	flags := flag.NewFlagSet("baseline apply", flag.ContinueOnError)
	flags.SetOutput(io.Discard)
	flags.StringVar(&request.repo, "repo", ".", "Git worktree or a path inside it")
	flags.StringVar(&request.plan, "plan", "", "Portable roundfix/baseline-plan/v1 JSON file")
	flags.StringVar(&request.confirmation, "confirm-plan", "", "Exact approved Plan Digest")
	flags.StringVar(&request.format, "format", "text", "Output format: text or json")
	if err := flags.Parse(args); err != nil {
		return baselineApplyCommandRequest{}, validationError{
			message: fmt.Sprintf("invalid baseline apply arguments: %v; run '%s baseline apply --help' for usage", err, app.Name),
		}
	}
	if remaining := flags.Args(); len(remaining) != 0 {
		return baselineApplyCommandRequest{}, validationError{
			message: fmt.Sprintf("unexpected argument %q; run '%s baseline apply --help' for usage", remaining[0], app.Name),
		}
	}
	request.repo = strings.TrimSpace(request.repo)
	request.plan = strings.TrimSpace(request.plan)
	request.confirmation = strings.TrimSpace(request.confirmation)
	request.format = strings.TrimSpace(request.format)
	if request.repo == "" {
		return request, validationError{message: "--repo cannot be empty"}
	}
	if request.plan == "" {
		return request, validationError{message: "--plan is required"}
	}
	if request.confirmation == "" {
		return request, validationError{message: "--confirm-plan is required"}
	}
	if request.format != "text" && request.format != "json" {
		return request, validationError{
			message: fmt.Sprintf("unsupported --format %q; use text or json", request.format),
		}
	}
	return request, nil
}

func baselineApplyJSONRequested(args []string) bool {
	for index, arg := range args {
		if arg == "--format=json" {
			return true
		}
		if arg == "--format" && index+1 < len(args) && args[index+1] == "json" {
			return true
		}
	}
	return false
}

func printBaselineApplyFailure(
	document baseline.PlanDocument,
	err error,
	jsonOutput bool,
	stdout io.Writer,
	stderr io.Writer,
) int {
	kind := baseline.ApplyErrorInvalid
	nextAction := "correct the command input and rerun roundfix baseline apply"
	var applyErr *baseline.ApplyError
	if errors.As(err, &applyErr) {
		kind = applyErr.Kind
		nextAction = applyErr.NextAction
	}
	state := "failed"
	exit := exitPreflight
	switch {
	case errors.Is(err, context.Canceled), errors.Is(err, context.DeadlineExceeded):
		kind = baseline.ApplyErrorExecution
		nextAction = "rerun roundfix baseline apply when the operation can complete"
		exit = exitSIGINT
	case kind == baseline.ApplyErrorApproval, kind == baseline.ApplyErrorStale:
		state = "action_required"
		exit = exitUnverified
	case kind == baseline.ApplyErrorExecution:
		exit = exitRunFailed
	case kind == baseline.ApplyErrorVerification:
		state = "recovery_required"
		exit = exitRunFailed
	}
	fmt.Fprintf(stderr, "%s: baseline apply failed: %v\n", app.Name, err)
	result := baseline.Result{
		SchemaVersion:      baseline.ResultSchemaVersion,
		Operation:          "apply",
		State:              state,
		Category:           string(kind),
		Message:            err.Error(),
		NextAction:         nextAction,
		PlanDigest:         document.PlanDigest,
		VerifiedPostimages: []baseline.Postimage{},
		Warnings:           []baseline.Finding{},
		Recommendations:    []string{},
	}
	if writeErr := writeBaselineApplyResult(result, jsonOutput, stdout); writeErr != nil {
		fmt.Fprintf(stderr, "%s: baseline apply output failed: %v\n", app.Name, writeErr)
		return exitRunFailed
	}
	return exit
}

func writeBaselineApplyResult(result baseline.Result, jsonOutput bool, stdout io.Writer) error {
	if jsonOutput {
		data, err := baseline.MarshalResult(result)
		if err != nil {
			return err
		}
		_, err = stdout.Write(data)
		return err
	}
	fmt.Fprintf(stdout, "Baseline apply: %s\n", strings.ReplaceAll(result.State, "_", " "))
	if result.PlanDigest != "" {
		fmt.Fprintf(stdout, "Plan Digest: %s\n", result.PlanDigest)
	}
	if result.Message != "" {
		fmt.Fprintf(stdout, "Result: %s\n", result.Message)
	}
	if result.NextAction != "" {
		fmt.Fprintf(stdout, "Next action: %s\n", result.NextAction)
	}
	fmt.Fprintf(stdout, "Verified postimages: %d\n", len(result.VerifiedPostimages))
	if len(result.VerifiedHistoryMoves) != 0 {
		fmt.Fprintf(stdout, "Verified history moves: %d\n", len(result.VerifiedHistoryMoves))
		for _, move := range result.VerifiedHistoryMoves {
			fmt.Fprintf(stdout, "- move %s -> %s (%s)\n", move.From, move.To, move.ContentIdentity)
		}
	}
	for _, recommendation := range result.Recommendations {
		fmt.Fprintf(stdout, "Recommendation (not run): %s\n", recommendation)
	}
	for _, warning := range result.Warnings {
		fmt.Fprintf(stdout, "Warning: %s: %s: %s\n", warning.Code, warning.Path, warning.Message)
	}
	return nil
}
