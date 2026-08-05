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

type reconcileProcessController interface {
	InspectTree(context.Context, int, string) ([]int, error)
	TerminateTreeAndWait(context.Context, int, string) ([]store.TerminationOutcome, error)
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
	SchemaVersion       string                  `json:"schemaVersion"`
	Mode                string                  `json:"mode"`
	Repository          string                  `json:"repository"`
	ApplyCommand        string                  `json:"applyCommand"`
	Results             []reconcileResult       `json:"results"`
	Runs                []reconcileResult       `json:"runs"`
	Summary             reconcileSummary        `json:"summary"`
	ProcessCandidates   []reconcileDebrisResult `json:"processCandidates"`
	RunBranchCandidates []reconcileDebrisResult `json:"runBranchCandidates"`
	PreservedCandidates []reconcileDebrisResult `json:"preservedCandidates"`
	DebrisSummary       reconcileDebrisSummary  `json:"debrisSummary"`
}

type reconcileDebrisResult struct {
	Kind              string `json:"kind"`
	RunID             string `json:"runId"`
	Outcome           string `json:"outcome"`
	OwnerPID          int    `json:"ownerPid"`
	ProcessIDs        []int  `json:"processIds"`
	RunBranch         string `json:"runBranch"`
	TargetBranch      string `json:"targetBranch"`
	Worktree          string `json:"worktree"`
	SpecSlug          string `json:"specSlug"`
	SupersedingReport string `json:"supersedingReport"`
	Proof             string `json:"proof"`
	Action            string `json:"action"`
	RefusalReason     string `json:"refusalReason"`

	run            store.Run
	classification runworktree.BranchSetClassification
}

type reconcileDebrisSummary struct {
	ProcessCandidates   int `json:"processCandidates"`
	RunBranchCandidates int `json:"runBranchCandidates"`
	Preserved           int `json:"preserved"`
	ProcessesApplied    int `json:"processesApplied"`
	RunBranchesApplied  int `json:"runBranchesApplied"`
}

type reconcileRunSelection struct {
	selected []store.Run
	all      []store.Run
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
	if opts.apply && (report.Summary.Safe+report.Summary.Superseded > 0 ||
		len(report.ProcessCandidates) > 0 || len(report.RunBranchCandidates) > 0) {
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

func loadReconcileRuns(ctx context.Context, homeDir, repository, runID string) (reconcileRunSelection, error) {
	reader, err := store.OpenReader(ctx, homeDir)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			if runID == "" {
				return reconcileRunSelection{selected: []store.Run{}, all: []store.Run{}}, nil
			}
			return reconcileRunSelection{}, validationError{message: fmt.Sprintf("Run %q does not exist", runID)}
		}
		return reconcileRunSelection{}, fmt.Errorf("open Run Database for reconciliation: %w", err)
	}
	defer func() {
		_ = reader.Close()
	}()

	if runID != "" {
		run, found, err := reader.Run(ctx, runID)
		if err != nil {
			return reconcileRunSelection{}, fmt.Errorf("read Run %q for reconciliation: %w", runID, err)
		}
		if !found {
			return reconcileRunSelection{}, validationError{message: fmt.Sprintf("Run %q does not exist", runID)}
		}
		if run.Kind != store.KindImplement {
			return reconcileRunSelection{}, validationError{
				message: fmt.Sprintf("Run %q is a review Run; reconcile accepts only terminal spec Runs", runID),
			}
		}
		if !store.IsTerminalState(run.State) {
			return reconcileRunSelection{}, validationError{
				message: fmt.Sprintf("Run %q is Active; stop it before reconciliation", runID),
			}
		}
		if !sameRepository(run.GitRoot, repository) {
			return reconcileRunSelection{}, validationError{
				message: fmt.Sprintf(
					"Run %q belongs to repository %q, not current repository %q",
					runID,
					run.GitRoot,
					repository,
				),
			}
		}
	}

	runs, err := reader.ListRuns(ctx, store.ListRunsQuery{
		GitRoot: repository,
		States:  store.StatesAll,
	})
	if err != nil {
		return reconcileRunSelection{}, fmt.Errorf("list Runs for reconciliation: %w", err)
	}
	all := make([]store.Run, 0, len(runs))
	selected := make([]store.Run, 0, len(runs))
	for _, run := range runs {
		if run.Kind != store.KindImplement {
			continue
		}
		all = append(all, run)
		if store.IsTerminalState(run.State) && (runID == "" || run.ID == runID) {
			selected = append(selected, run)
		}
	}
	return reconcileRunSelection{selected: selected, all: all}, nil
}

func inspectReconcileRuns(
	ctx context.Context,
	repository string,
	opts reconcileOptions,
	runs reconcileRunSelection,
) reconcileReport {
	mode := "dry-run"
	if opts.apply {
		mode = "apply"
	}
	report := reconcileReport{
		SchemaVersion:       reconcileSchemaVersion,
		Mode:                mode,
		Repository:          repository,
		ApplyCommand:        reconcileApplyCommand(opts.runID),
		Results:             make([]reconcileResult, 0, len(runs.selected)),
		ProcessCandidates:   make([]reconcileDebrisResult, 0),
		RunBranchCandidates: make([]reconcileDebrisResult, 0),
		PreservedCandidates: make([]reconcileDebrisResult, 0),
	}
	classifications := make(map[string]string, len(runs.selected))
	for _, run := range runs.selected {
		inspected, err := runworktree.InspectTerminalRun(ctx, run)
		result := newReconcileResult(inspected, opts.apply)
		if err != nil {
			result.Classification = string(runworktree.ReconciliationUnknown)
			result.RefusalReason = err.Error()
			result.Action = "inspect Git state and rerun: " + reconcileRetryCommand(run.ID)
			report.Summary.OperationalFailures++
		}
		report.Results = append(report.Results, result)
		classifications[run.ID] = result.Classification
		countReconcileClassification(&report.Summary, result.Classification)
	}
	// Keep the legacy results field unchanged while naming the same Run
	// collection explicitly for additive schema consumers. The slice is fully
	// populated before this alias is taken, and apply only updates its elements.
	report.Runs = report.Results
	report.Summary.Total = len(report.Results)
	inspectReconcileProcesses(ctx, opts, runs, &report)
	inspectReconcileRunBranches(ctx, repository, opts, runs, classifications, &report)
	report.DebrisSummary.ProcessCandidates = len(report.ProcessCandidates)
	report.DebrisSummary.RunBranchCandidates = len(report.RunBranchCandidates)
	report.DebrisSummary.Preserved = len(report.PreservedCandidates)
	return report
}

func inspectReconcileProcesses(
	ctx context.Context,
	opts reconcileOptions,
	runs reconcileRunSelection,
	report *reconcileReport,
) {
	selected := reconcileRunIDs(runs.selected)
	controller := commandDependenciesForContext(ctx).reconcileProcesses
	for _, run := range runs.all {
		if run.OwnerPID == nil || *run.OwnerPID <= 0 {
			continue
		}
		pid := *run.OwnerPID
		if !store.IsTerminalState(run.State) {
			if opts.runID == "" {
				report.PreservedCandidates = append(report.PreservedCandidates, reconcileDebrisResult{
					Kind:          "process",
					RunID:         run.ID,
					Outcome:       run.State,
					OwnerPID:      pid,
					Action:        "preserve",
					RefusalReason: fmt.Sprintf("process tree belongs to Active Run %q", run.ID),
				})
			}
			continue
		}
		if !selected[run.ID] {
			continue
		}
		if run.OwnerIdentityUnproven {
			report.PreservedCandidates = append(report.PreservedCandidates, reconcileDebrisResult{
				Kind:          "process",
				RunID:         run.ID,
				Outcome:       run.State,
				OwnerPID:      pid,
				Action:        "preserve",
				RefusalReason: "recorded owner process identity is unproven",
			})
			continue
		}
		processIDs, err := controller.InspectTree(ctx, pid, run.OwnerIdentity)
		if err != nil {
			report.PreservedCandidates = append(report.PreservedCandidates, reconcileDebrisResult{
				Kind:          "process",
				RunID:         run.ID,
				Outcome:       run.State,
				OwnerPID:      pid,
				Action:        "preserve",
				RefusalReason: fmt.Sprintf("owned process tree cannot be proven: %v", err),
			})
			continue
		}
		if len(processIDs) == 0 {
			continue
		}
		report.ProcessCandidates = append(report.ProcessCandidates, reconcileDebrisResult{
			Kind:       "process",
			RunID:      run.ID,
			Outcome:    run.State,
			OwnerPID:   pid,
			ProcessIDs: append([]int(nil), processIDs...),
			Proof: fmt.Sprintf(
				"terminal Run %q outcome %q owns inspected live process tree %v",
				run.ID,
				run.State,
				processIDs,
			),
			Action: debrisCandidateAction(opts.apply),
			run:    run,
		})
	}
}

type reconcileBranchGroupKey struct {
	targetBranch string
	specSlug     string
}

func inspectReconcileRunBranches(
	ctx context.Context,
	repository string,
	opts reconcileOptions,
	runs reconcileRunSelection,
	existingClassifications map[string]string,
	report *reconcileReport,
) {
	selected := reconcileRunIDs(runs.selected)
	groups := make(map[reconcileBranchGroupKey][]store.Run)
	groupOrder := make([]reconcileBranchGroupKey, 0)
	for _, run := range runs.all {
		key := reconcileBranchGroupKey{
			targetBranch: strings.TrimSpace(run.LocalBranch),
			specSlug:     strings.TrimSpace(run.SpecSlug),
		}
		if key.targetBranch == "" || key.specSlug == "" {
			continue
		}
		if _, found := groups[key]; !found {
			groupOrder = append(groupOrder, key)
		}
		groups[key] = append(groups[key], run)
	}

	for _, key := range groupOrder {
		groupRuns := groups[key]
		classification, err := commandDependenciesForContext(ctx).classifyRunBranchSet(
			ctx,
			repository,
			key.targetBranch,
			key.specSlug,
			groupRuns,
		)
		if err != nil {
			for _, run := range groupRuns {
				if !selected[run.ID] {
					continue
				}
				report.PreservedCandidates = append(report.PreservedCandidates, reconcileDebrisResult{
					Kind:          "runBranch",
					RunID:         run.ID,
					Outcome:       run.State,
					RunBranch:     runworktree.BranchName(run.ID),
					TargetBranch:  key.targetBranch,
					Worktree:      run.WorkDir,
					SpecSlug:      key.specSlug,
					Action:        "preserve",
					RefusalReason: fmt.Sprintf("Run Branch set cannot be proven: %v", err),
				})
			}
			continue
		}

		for _, run := range groupRuns {
			active := !store.IsTerminalState(run.State)
			if !selected[run.ID] && !(opts.runID == "" && active) {
				continue
			}
			branch := runworktree.BranchName(run.ID)
			if active {
				if branch == classification.Current || containsReconcileString(classification.Preserved, branch) {
					report.PreservedCandidates = append(report.PreservedCandidates, reconcileDebrisResult{
						Kind:          "runBranch",
						RunID:         run.ID,
						Outcome:       run.State,
						RunBranch:     branch,
						TargetBranch:  key.targetBranch,
						Worktree:      run.WorkDir,
						SpecSlug:      key.specSlug,
						Action:        "preserve",
						RefusalReason: fmt.Sprintf("Run Branch belongs to Active Run %q", run.ID),
					})
				}
				continue
			}
			if reconcileClassificationReleasable(existingClassifications[run.ID]) ||
				existingClassifications[run.ID] == string(runworktree.ReconciliationReleased) {
				continue
			}
			if containsReconcileString(classification.Releasable, branch) {
				if existingClassifications[run.ID] != string(runworktree.ReconciliationUnintegrated) {
					report.PreservedCandidates = append(report.PreservedCandidates, reconcileDebrisResult{
						Kind:          "runBranch",
						RunID:         run.ID,
						Outcome:       run.State,
						RunBranch:     branch,
						TargetBranch:  key.targetBranch,
						Worktree:      run.WorkDir,
						SpecSlug:      key.specSlug,
						Action:        "preserve",
						RefusalReason: fmt.Sprintf("Run Worktree classification %q is not a clean releasable branch surface", existingClassifications[run.ID]),
					})
					continue
				}
				proof := classification.ReleasableProofs[branch]
				if proof == "" {
					report.PreservedCandidates = append(report.PreservedCandidates, reconcileDebrisResult{
						Kind:          "runBranch",
						RunID:         run.ID,
						Outcome:       run.State,
						RunBranch:     branch,
						TargetBranch:  key.targetBranch,
						Worktree:      run.WorkDir,
						SpecSlug:      key.specSlug,
						Action:        "preserve",
						RefusalReason: "Run Branch superseding proof is missing",
					})
					continue
				}
				report.RunBranchCandidates = append(report.RunBranchCandidates, reconcileDebrisResult{
					Kind:              "runBranch",
					RunID:             run.ID,
					Outcome:           run.State,
					RunBranch:         branch,
					TargetBranch:      key.targetBranch,
					Worktree:          run.WorkDir,
					SpecSlug:          key.specSlug,
					SupersedingReport: proof,
					Proof: fmt.Sprintf(
						"Run Branch %q is superseded by QA Report %q; registered Run Worktree %q was inspected clean",
						branch,
						proof,
						run.WorkDir,
					),
					Action:         debrisCandidateAction(opts.apply),
					run:            run,
					classification: classification,
				})
				continue
			}
			if reason := classification.PreservedReasons[branch]; reason != "" {
				report.PreservedCandidates = append(report.PreservedCandidates, reconcileDebrisResult{
					Kind:          "runBranch",
					RunID:         run.ID,
					Outcome:       run.State,
					RunBranch:     branch,
					TargetBranch:  key.targetBranch,
					Worktree:      run.WorkDir,
					SpecSlug:      key.specSlug,
					Action:        "preserve",
					RefusalReason: reason,
				})
			} else if branch == classification.Current {
				report.PreservedCandidates = append(report.PreservedCandidates, reconcileDebrisResult{
					Kind:              "runBranch",
					RunID:             run.ID,
					Outcome:           run.State,
					RunBranch:         branch,
					TargetBranch:      key.targetBranch,
					Worktree:          run.WorkDir,
					SpecSlug:          key.specSlug,
					SupersedingReport: classification.CurrentReport,
					Action:            "preserve",
					RefusalReason:     "Run Branch carries current QA Report evidence and is not superseded",
				})
			}
		}
	}
}

func reconcileRunIDs(runs []store.Run) map[string]bool {
	ids := make(map[string]bool, len(runs))
	for _, run := range runs {
		ids[run.ID] = true
	}
	return ids
}

func containsReconcileString(values []string, want string) bool {
	for _, value := range values {
		if value == want {
			return true
		}
	}
	return false
}

func debrisCandidateAction(apply bool) string {
	if apply {
		return "reclaim after fresh proof"
	}
	return "would reclaim with --apply"
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
	applyReconcileWorktrees(ctx, homeDir, opts, report)
	for index := range report.RunBranchCandidates {
		candidate := &report.RunBranchCandidates[index]
		if err := runworktree.ApplyRunBranchCandidate(ctx, candidate.classification, candidate.RunBranch); err != nil {
			candidate.Action = "preserve"
			candidate.RefusalReason = err.Error()
			report.DebrisSummary.Preserved++
			report.Summary.OperationalFailures++
			continue
		}
		candidate.Action = "released"
		candidate.RefusalReason = ""
		report.DebrisSummary.RunBranchesApplied++
	}

	controller := commandDependenciesForContext(ctx).reconcileProcesses
	for index := range report.ProcessCandidates {
		candidate := &report.ProcessCandidates[index]
		outcomes, err := controller.TerminateTreeAndWait(
			ctx,
			candidate.OwnerPID,
			candidate.run.OwnerIdentity,
		)
		if err == nil {
			for _, outcome := range outcomes {
				if outcome.Proven {
					continue
				}
				err = fmt.Errorf("process %d absence was not proven: %s", outcome.PID, outcome.Reason)
				break
			}
		}
		if err != nil {
			candidate.Action = "preserve"
			candidate.RefusalReason = err.Error()
			report.DebrisSummary.Preserved++
			report.Summary.OperationalFailures++
			continue
		}
		if len(outcomes) == 0 {
			candidate.Action = "none; process tree already absent"
			continue
		}
		candidate.Action = "terminated"
		candidate.RefusalReason = ""
		report.DebrisSummary.ProcessesApplied++
	}
}

func applyReconcileWorktrees(
	ctx context.Context,
	homeDir string,
	opts reconcileOptions,
	report *reconcileReport,
) {
	hasCandidates := false
	for _, result := range report.Results {
		if reconcileClassificationReleasable(result.Classification) {
			hasCandidates = true
			break
		}
	}
	if !hasCandidates {
		return
	}
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
	for _, candidate := range report.ProcessCandidates {
		fmt.Fprintf(&output, "Process candidate: Run %s owner PID %d\n", candidate.RunID, candidate.OwnerPID)
		fmt.Fprintf(&output, "  process-ids: %v\n", candidate.ProcessIDs)
		fmt.Fprintf(&output, "  proof: %s\n", textReconcileValue(candidate.Proof))
		fmt.Fprintf(&output, "  action: %s\n", textReconcileValue(candidate.Action))
		fmt.Fprintf(&output, "  refusal-reason: %s\n", textReconcileValue(candidate.RefusalReason))
	}
	for _, candidate := range report.RunBranchCandidates {
		fmt.Fprintf(&output, "Run Branch candidate: %s (Run %s)\n", candidate.RunBranch, candidate.RunID)
		fmt.Fprintf(&output, "  target-branch: %s\n", textReconcileValue(candidate.TargetBranch))
		fmt.Fprintf(&output, "  worktree: %s\n", textReconcileValue(candidate.Worktree))
		fmt.Fprintf(&output, "  superseding-report: %s\n", textReconcileValue(candidate.SupersedingReport))
		fmt.Fprintf(&output, "  proof: %s\n", textReconcileValue(candidate.Proof))
		fmt.Fprintf(&output, "  action: %s\n", textReconcileValue(candidate.Action))
		fmt.Fprintf(&output, "  refusal-reason: %s\n", textReconcileValue(candidate.RefusalReason))
	}
	for _, candidate := range report.PreservedCandidates {
		fmt.Fprintf(&output, "Preserved candidate: kind=%s Run=%s\n", candidate.Kind, candidate.RunID)
		fmt.Fprintf(&output, "  owner-pid: %d\n", candidate.OwnerPID)
		fmt.Fprintf(&output, "  run-branch: %s\n", textReconcileValue(candidate.RunBranch))
		fmt.Fprintf(&output, "  worktree: %s\n", textReconcileValue(candidate.Worktree))
		fmt.Fprintf(&output, "  action: %s\n", textReconcileValue(candidate.Action))
		fmt.Fprintf(&output, "  refusal-reason: %s\n", textReconcileValue(candidate.RefusalReason))
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
	fmt.Fprintf(
		&output,
		"Debris summary: process-candidates=%d run-branch-candidates=%d preserved=%d processes-applied=%d run-branches-applied=%d\n",
		report.DebrisSummary.ProcessCandidates,
		report.DebrisSummary.RunBranchCandidates,
		report.DebrisSummary.Preserved,
		report.DebrisSummary.ProcessesApplied,
		report.DebrisSummary.RunBranchesApplied,
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
