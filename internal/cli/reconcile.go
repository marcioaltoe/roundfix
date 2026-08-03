package cli

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"roundfix/internal/app"
	"roundfix/internal/store"
	runworktree "roundfix/internal/worktree"
)

const reconcileSchemaVersion = "roundfix-reconcile/v1"

type reconcileOptions struct {
	runID  string
	apply  bool
	format string
}

type reconcileResult struct {
	RunID             string `json:"runId"`
	Outcome           string `json:"outcome"`
	Classification    string `json:"classification"`
	RunBranch         string `json:"runBranch"`
	RunHead           string `json:"runHead"`
	TargetBranch      string `json:"targetBranch"`
	TargetHead        string `json:"targetHead"`
	Worktree          string `json:"worktree"`
	SupersedingReport string `json:"supersedingReport"`
	Evidence          string `json:"evidence"`
	Action            string `json:"action"`
	RefusalReason     string `json:"refusalReason"`

	inspected runworktree.RunWorktreeReconciliation
}

type reconcileSummary struct {
	Total               int `json:"total"`
	Safe                int `json:"safe"`
	Superseded          int `json:"superseded"`
	Unintegrated        int `json:"unintegrated"`
	Dirty               int `json:"dirty"`
	Unknown             int `json:"unknown"`
	Released            int `json:"released"`
	Applied             int `json:"applied"`
	Preserved           int `json:"preserved"`
	OperationalFailures int `json:"operationalFailures"`
}

type reconcileReport struct {
	SchemaVersion string            `json:"schemaVersion"`
	Mode          string            `json:"mode"`
	Repository    string            `json:"repository"`
	ApplyCommand  string            `json:"applyCommand"`
	Results       []reconcileResult `json:"results"`
	Summary       reconcileSummary  `json:"summary"`
}

func runReconcileCommand(ctx context.Context, args []string, stdout, stderr io.Writer, environment commandEnvironment) int {
	if commandWantsHelp(args) {
		fmt.Fprint(stdout, commandUsage("reconcile"))
		return exitOK
	}
	opts, err := parseReconcileOptions(args)
	if err != nil {
		printReconcileValidationFailure(err, stderr)
		return exitPreflight
	}

	loaded, err := loadCommandConfig(environment, stderr)
	if err != nil {
		printReconcileValidationFailure(err, stderr)
		return exitPreflight
	}
	repository := strings.TrimSpace(loaded.GitRoot)
	if repository == "" {
		printReconcileValidationFailure(
			validationError{message: "reconcile requires running inside a Git repository"},
			stderr,
		)
		return exitPreflight
	}
	repository, err = filepath.EvalSymlinks(repository)
	if err != nil {
		printReconcileOperationalFailure(
			fmt.Errorf("resolve current Git repository %q: %w", loaded.GitRoot, err),
			reconcileRetryCommand(opts.runID),
			stderr,
		)
		return exitRunFailed
	}

	runs, err := loadReconcileRuns(ctx, loaded.HomeDir, repository, opts.runID)
	if err != nil {
		var invalid validationError
		if errors.As(err, &invalid) {
			printReconcileValidationFailure(err, stderr)
			return exitPreflight
		}
		printReconcileOperationalFailure(err, reconcileRetryCommand(opts.runID), stderr)
		return exitRunFailed
	}

	report := inspectReconcileRuns(ctx, repository, opts, runs)
	if opts.apply && report.Summary.Safe+report.Summary.Superseded > 0 {
		applyReconcileReport(ctx, loaded.HomeDir, opts, &report)
	}
	if err := printReconcileReport(stdout, opts.format, report); err != nil {
		printReconcileOperationalFailure(
			fmt.Errorf("write reconciliation report: %w", err),
			reconcileRetryCommand(opts.runID),
			stderr,
		)
		return exitRunFailed
	}
	if report.Summary.OperationalFailures > 0 {
		printReconcileOperationalFailure(
			fmt.Errorf("%d reconciliation operation(s) failed", report.Summary.OperationalFailures),
			reconcileRetryCommand(opts.runID),
			stderr,
		)
		return exitRunFailed
	}
	return exitOK
}

func parseReconcileOptions(args []string) (reconcileOptions, error) {
	opts := reconcileOptions{format: "text"}
	for index := 0; index < len(args); index++ {
		arg := args[index]
		switch {
		case arg == "--apply":
			opts.apply = true
		case arg == "--format":
			if index+1 >= len(args) {
				return opts, validationError{message: "flag needs an argument: --format"}
			}
			index++
			opts.format = strings.TrimSpace(args[index])
		case strings.HasPrefix(arg, "--format="):
			opts.format = strings.TrimSpace(strings.TrimPrefix(arg, "--format="))
		case strings.HasPrefix(arg, "-"):
			return opts, validationError{message: fmt.Sprintf("unknown flag %q", arg)}
		default:
			if opts.runID != "" {
				return opts, validationError{message: fmt.Sprintf("unexpected argument %q", arg)}
			}
			opts.runID = strings.TrimSpace(arg)
			if opts.runID == "" {
				return opts, validationError{message: "Run ID cannot be empty"}
			}
		}
	}
	if opts.format != "text" && opts.format != "json" {
		return opts, validationError{
			message: fmt.Sprintf("unknown --format %q; use text or json", opts.format),
		}
	}
	return opts, nil
}

func loadReconcileRuns(ctx context.Context, homeDir, repository, runID string) ([]store.Run, error) {
	reader, err := store.OpenReader(ctx, homeDir)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			if runID == "" {
				return []store.Run{}, nil
			}
			return nil, validationError{message: fmt.Sprintf("Run %q does not exist", runID)}
		}
		return nil, fmt.Errorf("open Run Database for reconciliation: %w", err)
	}
	defer func() {
		_ = reader.Close()
	}()

	if runID != "" {
		run, found, err := reader.Run(ctx, runID)
		if err != nil {
			return nil, fmt.Errorf("read Run %q for reconciliation: %w", runID, err)
		}
		if !found {
			return nil, validationError{message: fmt.Sprintf("Run %q does not exist", runID)}
		}
		if run.Kind != store.KindImplement {
			return nil, validationError{
				message: fmt.Sprintf("Run %q is a review Run; reconcile accepts only terminal spec Runs", runID),
			}
		}
		if !store.IsTerminalState(run.State) {
			return nil, validationError{
				message: fmt.Sprintf("Run %q is Active; stop it before reconciliation", runID),
			}
		}
		if !sameRepository(run.GitRoot, repository) {
			return nil, validationError{
				message: fmt.Sprintf(
					"Run %q belongs to repository %q, not current repository %q",
					runID,
					run.GitRoot,
					repository,
				),
			}
		}
		return []store.Run{run}, nil
	}

	runs, err := reader.ListRuns(ctx, store.ListRunsQuery{
		GitRoot: repository,
		States:  store.StatesTerminal,
	})
	if err != nil {
		return nil, fmt.Errorf("list terminal Runs for reconciliation: %w", err)
	}
	specRuns := make([]store.Run, 0, len(runs))
	for _, run := range runs {
		if run.Kind == store.KindImplement {
			specRuns = append(specRuns, run)
		}
	}
	return specRuns, nil
}

func inspectReconcileRuns(
	ctx context.Context,
	repository string,
	opts reconcileOptions,
	runs []store.Run,
) reconcileReport {
	mode := "dry-run"
	if opts.apply {
		mode = "apply"
	}
	report := reconcileReport{
		SchemaVersion: reconcileSchemaVersion,
		Mode:          mode,
		Repository:    repository,
		ApplyCommand:  reconcileApplyCommand(opts.runID),
		Results:       make([]reconcileResult, 0, len(runs)),
	}
	for _, run := range runs {
		inspected, err := runworktree.InspectTerminalRun(ctx, run)
		result := newReconcileResult(inspected, opts.apply)
		if err != nil {
			result.Classification = string(runworktree.ReconciliationUnknown)
			result.RefusalReason = err.Error()
			result.Action = "inspect Git state and rerun: " + reconcileRetryCommand(run.ID)
			report.Summary.OperationalFailures++
		}
		report.Results = append(report.Results, result)
		countReconcileClassification(&report.Summary, result.Classification)
	}
	report.Summary.Total = len(report.Results)
	return report
}

func newReconcileResult(
	inspected runworktree.RunWorktreeReconciliation,
	apply bool,
) reconcileResult {
	result := reconcileResult{
		RunID:             inspected.RunID,
		Outcome:           inspected.Outcome,
		Classification:    string(inspected.State),
		RunBranch:         inspected.Branch,
		RunHead:           inspected.RunHead,
		TargetBranch:      inspected.TargetBranch,
		TargetHead:        inspected.TargetHead,
		Worktree:          inspected.Path,
		SupersedingReport: inspected.SupersedingReport,
		Evidence:          inspected.Reason,
		inspected:         inspected,
	}
	switch inspected.State {
	case runworktree.ReconciliationSafe, runworktree.ReconciliationSuperseded:
		if apply {
			result.Action = "release after fresh safety proof"
		} else {
			result.Action = "would release with --apply"
		}
	case runworktree.ReconciliationReleased:
		result.Action = "none"
	default:
		result.Action = "preserve"
		result.RefusalReason = inspected.Reason
	}
	return result
}

func applyReconcileReport(
	ctx context.Context,
	homeDir string,
	opts reconcileOptions,
	report *reconcileReport,
) {
	runStore, err := store.Open(ctx, homeDir)
	if err != nil {
		for index := range report.Results {
			if !reconcileClassificationReleasable(report.Results[index].Classification) {
				continue
			}
			report.Results[index].Action = "preserved; repair the Run Database and rerun: " +
				reconcileApplyCommand(opts.runID)
			report.Results[index].RefusalReason = fmt.Sprintf("open Run Database for apply: %v", err)
			report.Summary.OperationalFailures++
		}
		return
	}
	defer func() {
		_ = runStore.Close()
	}()

	for index := range report.Results {
		result := &report.Results[index]
		if !reconcileClassificationReleasable(result.Classification) {
			continue
		}
		if err := runworktree.ApplyTerminalRun(ctx, runStore, result.inspected); err != nil {
			result.Action = "preserved; inspect remaining Git state and rerun: " +
				reconcileApplyCommand(result.RunID)
			result.RefusalReason = err.Error()
			report.Summary.OperationalFailures++
			continue
		}
		result.Action = "released"
		result.RefusalReason = ""
		report.Summary.Applied++
	}
}

func reconcileClassificationReleasable(classification string) bool {
	switch runworktree.ReconciliationState(classification) {
	case runworktree.ReconciliationSafe, runworktree.ReconciliationSuperseded:
		return true
	default:
		return false
	}
}

func countReconcileClassification(summary *reconcileSummary, classification string) {
	switch runworktree.ReconciliationState(classification) {
	case runworktree.ReconciliationSafe:
		summary.Safe++
	case runworktree.ReconciliationSuperseded:
		summary.Superseded++
	case runworktree.ReconciliationUnintegrated:
		summary.Unintegrated++
		summary.Preserved++
	case runworktree.ReconciliationDirty:
		summary.Dirty++
		summary.Preserved++
	case runworktree.ReconciliationUnknown:
		summary.Unknown++
		summary.Preserved++
	case runworktree.ReconciliationReleased:
		summary.Released++
	}
}

func printReconcileReport(stdout io.Writer, format string, report reconcileReport) error {
	if format == "json" {
		encoder := json.NewEncoder(stdout)
		encoder.SetEscapeHTML(false)
		return encoder.Encode(report)
	}
	_, err := io.WriteString(stdout, reconcileText(report))
	return err
}

func reconcileText(report reconcileReport) string {
	var output strings.Builder
	fmt.Fprintf(&output, "Repository: %s\n", report.Repository)
	fmt.Fprintf(&output, "Mode: %s\n", report.Mode)
	for _, result := range report.Results {
		fmt.Fprintf(&output, "Run: %s\n", textReconcileValue(result.RunID))
		fmt.Fprintf(&output, "  outcome: %s\n", textReconcileValue(result.Outcome))
		fmt.Fprintf(&output, "  classification: %s\n", textReconcileValue(result.Classification))
		fmt.Fprintf(&output, "  run-branch: %s\n", textReconcileValue(result.RunBranch))
		fmt.Fprintf(&output, "  run-head: %s\n", textReconcileValue(result.RunHead))
		fmt.Fprintf(&output, "  target-branch: %s\n", textReconcileValue(result.TargetBranch))
		fmt.Fprintf(&output, "  target-head: %s\n", textReconcileValue(result.TargetHead))
		fmt.Fprintf(&output, "  worktree: %s\n", textReconcileValue(result.Worktree))
		fmt.Fprintf(&output, "  superseding-report: %s\n", textReconcileValue(result.SupersedingReport))
		fmt.Fprintf(&output, "  evidence: %s\n", textReconcileValue(result.Evidence))
		fmt.Fprintf(&output, "  action: %s\n", textReconcileValue(result.Action))
		fmt.Fprintf(&output, "  refusal-reason: %s\n", textReconcileValue(result.RefusalReason))
	}
	summary := report.Summary
	fmt.Fprintf(
		&output,
		"Summary: total=%d safe=%d superseded=%d unintegrated=%d dirty=%d unknown=%d released=%d applied=%d preserved=%d operational-failures=%d\n",
		summary.Total,
		summary.Safe,
		summary.Superseded,
		summary.Unintegrated,
		summary.Dirty,
		summary.Unknown,
		summary.Released,
		summary.Applied,
		summary.Preserved,
		summary.OperationalFailures,
	)
	if report.Mode == "dry-run" {
		fmt.Fprintf(&output, "Apply with: %s\n", report.ApplyCommand)
	}
	return output.String()
}

func textReconcileValue(value string) string {
	if strings.TrimSpace(value) == "" {
		return "-"
	}
	return value
}

func reconcileApplyCommand(runID string) string {
	if strings.TrimSpace(runID) == "" {
		return "roundfix reconcile --apply"
	}
	return "roundfix reconcile " + strings.TrimSpace(runID) + " --apply"
}

func reconcileRetryCommand(runID string) string {
	if strings.TrimSpace(runID) == "" {
		return "roundfix reconcile"
	}
	return "roundfix reconcile " + strings.TrimSpace(runID)
}

func sameRepository(left, right string) bool {
	left = strings.TrimSpace(left)
	right = strings.TrimSpace(right)
	if left == "" || right == "" {
		return false
	}
	if resolved, err := filepath.EvalSymlinks(left); err == nil {
		left = resolved
	}
	if resolved, err := filepath.EvalSymlinks(right); err == nil {
		right = resolved
	}
	return filepath.Clean(left) == filepath.Clean(right)
}

func printReconcileValidationFailure(err error, stderr io.Writer) {
	fmt.Fprintf(stderr, "%s: reconcile failed: %v\n", app.Name, err)
	fmt.Fprintf(stderr, "Run '%s reconcile --help' for usage.\n", app.Name)
}

func printReconcileOperationalFailure(err error, nextAction string, stderr io.Writer) {
	fmt.Fprintf(stderr, "%s: reconcile failed: %v\n", app.Name, err)
	fmt.Fprintf(stderr, "Next safe action: %s\n", nextAction)
}
