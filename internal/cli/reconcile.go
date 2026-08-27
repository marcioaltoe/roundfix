package cli

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"
	"time"

	"roundfix/internal/app"
	roundconfig "roundfix/internal/config"
	"roundfix/internal/daemon"
	"roundfix/internal/runevent"
	"roundfix/internal/spec"
	"roundfix/internal/store"
	runworktree "roundfix/internal/worktree"
)

const reconcileSchemaVersion = "roundfix-reconcile/v1"

type reconcileOptions struct {
	runID             string
	apply             bool
	discardSuperseded bool
	carryForward      bool
	format            string
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

	inspected   runworktree.RunWorktreeReconciliation
	disposition daemon.BranchDisposition
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
	CarryForwards       []spec.CarryForward     `json:"carryForwards"`
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
	selected     []store.Run
	all          []store.Run
	coverage     []daemon.RunTaskCoverage
	taskEvidence map[string]map[string]reconcileTaskEvidence
}

type reconcileTaskEvidence struct {
	verificationPassed bool
	settledCompleted   bool
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
	carryForwardRefusal := ""
	if opts.carryForward && len(runs.selected) != 1 {
		carryForwardRefusal = fmt.Sprintf(
			"Run %q could not be selected; carry-forward accepts one Run with outcome %s",
			opts.runID,
			strings.Join(carryForwardAcceptedStates, " or "),
		)
	} else if opts.carryForward && !slices.Contains(carryForwardAcceptedStates, runs.selected[0].State) {
		carryForwardRefusal = fmt.Sprintf(
			"Run %q has outcome %s; carry-forward accepts outcomes %s",
			opts.runID,
			runs.selected[0].State,
			strings.Join(carryForwardAcceptedStates, " and "),
		)
	}
	resolvedSpecsRoot, resolveSpecsErr := roundconfig.ResolveSpecsRoot(loaded, repository)
	if resolveSpecsErr != nil {
		if opts.carryForward {
			printReconcileOperationalFailure(
				fmt.Errorf("resolve Specs Root for carry-forward: %w", resolveSpecsErr),
				reconcileRetryCommand(opts.runID),
				stderr,
			)
			return exitRunFailed
		}
	} else if !resolvedSpecsRoot.External {
		if carryForwardRefusal == "" {
			carryForwards, inspectErr := inspectCarryForwards(ctx, repository, resolvedSpecsRoot, runs)
			if inspectErr != nil {
				if opts.carryForward {
					printReconcileOperationalFailure(inspectErr, reconcileRetryCommand(opts.runID), stderr)
					return exitRunFailed
				}
			} else {
				report.CarryForwards = carryForwards
			}
		}
	} else if opts.carryForward && carryForwardRefusal == "" {
		carryForwardRefusal = "carry-forward requires a repository-local Specs Root"
	}
	if opts.carryForward && carryForwardRefusal == "" {
		carryForwardRefusal = carryForwardRefusalReason(report.CarryForwards)
		if carryForwardRefusal == "" {
			if err := applyCarryForwards(ctx, repository, report.CarryForwards); err != nil {
				for index := range report.CarryForwards {
					report.CarryForwards[index].Action = "preserve"
					report.CarryForwards[index].RefusalReason = err.Error()
				}
				report.Summary.OperationalFailures++
			} else {
				for index := range report.CarryForwards {
					report.CarryForwards[index].Action = "carried forward"
				}
			}
		}
	}
	discardRefusals := 0
	if opts.discardSuperseded {
		artifactRoot, resolveErr := roundconfig.ResolveArtifactDirectory(
			loaded.Config.Defaults.ArtifactDir,
			repository,
			loaded.HomeDir,
		)
		if resolveErr != nil {
			printReconcileOperationalFailure(
				fmt.Errorf("resolve Artifact Root for branch disposition records: %w", resolveErr),
				reconcileRetryCommand(opts.runID),
				stderr,
			)
			return exitRunFailed
		}
		discardRefusals = discardSupersededReconcileResults(ctx, artifactRoot, &report)
	}
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
	if discardRefusals > 0 {
		printReconcileValidationFailure(
			validationError{message: fmt.Sprintf("%d Run Branch disposition(s) could not be proven superseded", discardRefusals)},
			stderr,
		)
		return exitPreflight
	}
	if carryForwardRefusal != "" {
		printReconcileValidationFailure(validationError{message: carryForwardRefusal}, stderr)
		return exitPreflight
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
		case arg == "--discard-superseded":
			opts.discardSuperseded = true
		case arg == "--carry-forward":
			opts.carryForward = true
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
	mutationModes := 0
	for _, enabled := range []bool{opts.apply, opts.discardSuperseded, opts.carryForward} {
		if enabled {
			mutationModes++
		}
	}
	if mutationModes > 1 {
		return opts, validationError{message: "--apply, --discard-superseded, and --carry-forward are mutually exclusive"}
	}
	if opts.carryForward && opts.runID == "" {
		return opts, validationError{message: "--carry-forward requires one Run ID"}
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
	coverage, taskEvidence, err := loadReconcileTaskCoverage(ctx, reader, all)
	if err != nil {
		return reconcileRunSelection{}, err
	}
	return reconcileRunSelection{selected: selected, all: all, coverage: coverage, taskEvidence: taskEvidence}, nil
}

func loadReconcileTaskCoverage(
	ctx context.Context,
	reader *store.Store,
	runs []store.Run,
) ([]daemon.RunTaskCoverage, map[string]map[string]reconcileTaskEvidence, error) {
	const pageSize = 256
	coverage := make([]daemon.RunTaskCoverage, 0, len(runs))
	evidenceByRun := make(map[string]map[string]reconcileTaskEvidence, len(runs))
	for _, run := range runs {
		evidence := make(map[string]reconcileTaskEvidence)
		var cursor int64
		for {
			entries, err := reader.RunEventsAfter(ctx, run.ID, cursor, pageSize)
			if err != nil {
				return nil, nil, fmt.Errorf("read Task coverage for later Run %q: %w", run.ID, err)
			}
			for _, entry := range entries {
				cursor = entry.Cursor
				if entry.Event.Kind != runevent.KindDaemonTask && entry.Event.Kind != runevent.KindDaemonVerification {
					continue
				}
				var payload struct {
					Task    string `json:"task"`
					Phase   string `json:"phase"`
					Status  string `json:"status"`
					Verdict string `json:"verdict"`
				}
				if err := json.Unmarshal(entry.Event.Payload, &payload); err != nil {
					return nil, nil, fmt.Errorf("decode Task coverage event %d for later Run %q: %w", entry.Cursor, run.ID, err)
				}
				taskID := strings.TrimSpace(payload.Task)
				if taskID == "" {
					taskID = strings.TrimSpace(entry.Event.ReviewIssue)
				}
				if taskID == "" {
					continue
				}
				taskEvidence := evidence[taskID]
				if entry.Event.Kind == runevent.KindDaemonTask && payload.Phase == "settled" && payload.Status == "completed" {
					taskEvidence.settledCompleted = true
				}
				if entry.Event.Kind == runevent.KindDaemonVerification &&
					payload.Phase == string(runevent.VerificationPhaseVerdict) &&
					payload.Verdict == string(runevent.VerificationVerdictPassed) {
					taskEvidence.verificationPassed = true
				}
				evidence[taskID] = taskEvidence
			}
			if len(entries) < pageSize {
				break
			}
		}
		tasks := make([]string, 0, len(evidence))
		for taskID, taskEvidence := range evidence {
			if taskEvidence.settledCompleted {
				tasks = append(tasks, taskID)
			}
		}
		slices.Sort(tasks)
		coverage = append(coverage, daemon.RunTaskCoverage{Run: run, CompletedTasks: tasks})
		evidenceByRun[run.ID] = evidence
	}
	return coverage, evidenceByRun, nil
}

func applyCarryForwards(ctx context.Context, repository string, candidates []spec.CarryForward) (returnErr error) {
	status, err := reconcileGitText(ctx, repository, "status", "--porcelain=v1", "--untracked-files=all")
	if err != nil {
		return fmt.Errorf("inspect checkout before carry-forward: %w", err)
	}
	if strings.TrimSpace(status) != "" {
		return errors.New("checkout changed after carry-forward proof; commit, stash, or discard its changes and retry")
	}
	checkoutHead, err := reconcileGitText(ctx, repository, "rev-parse", "HEAD")
	if err != nil {
		return fmt.Errorf("read checkout HEAD before carry-forward: %w", err)
	}
	stagingWorktree, cleanup, err := createCarryForwardStaging(ctx, repository, checkoutHead)
	if err != nil {
		return err
	}
	defer func() {
		returnErr = errors.Join(returnErr, cleanup())
	}()

	for _, candidate := range candidates {
		if err := stageCarryForwardCandidate(ctx, stagingWorktree, candidate); err != nil {
			return err
		}
	}
	stagedHead, err := reconcileGitText(ctx, stagingWorktree, "rev-parse", "HEAD")
	if err != nil {
		return fmt.Errorf("read staged carry-forward HEAD: %w", err)
	}
	currentHead, err := reconcileGitText(ctx, repository, "rev-parse", "HEAD")
	if err != nil {
		return fmt.Errorf("re-read checkout HEAD before carry-forward: %w", err)
	}
	if currentHead != checkoutHead {
		return fmt.Errorf("checkout HEAD moved from %s to %s during carry-forward proof", checkoutHead, currentHead)
	}
	status, err = reconcileGitText(ctx, repository, "status", "--porcelain=v1", "--untracked-files=all")
	if err != nil {
		return fmt.Errorf("re-inspect checkout before carry-forward: %w", err)
	}
	if strings.TrimSpace(status) != "" {
		return errors.New("checkout changed during carry-forward proof; no Task was carried")
	}
	if _, err := reconcileGitRaw(ctx, repository, "merge", "--ff-only", stagedHead); err != nil {
		return fmt.Errorf("fast-forward checkout with carried Tasks: %w", err)
	}
	return nil
}

func reconcileGitText(ctx context.Context, workDir string, args ...string) (string, error) {
	output, err := reconcileGitRaw(ctx, workDir, args...)
	return strings.TrimSpace(string(output)), err
}

func reconcileGitRaw(ctx context.Context, workDir string, args ...string) ([]byte, error) {
	return reconcileGitRawInput(ctx, workDir, nil, args...)
}

func reconcileGitRawInput(ctx context.Context, workDir string, input []byte, args ...string) ([]byte, error) {
	commandArgs := append([]string{"-c", "core.fsmonitor=false", "-c", "commit.gpgSign=false"}, args...)
	command := exec.CommandContext(ctx, "git", commandArgs...)
	command.Dir = workDir
	command.Env = append(os.Environ(), "GIT_TERMINAL_PROMPT=0")
	if input != nil {
		command.Stdin = strings.NewReader(string(input))
	}
	output, err := command.CombinedOutput()
	if err != nil {
		diagnostic := strings.TrimSpace(string(output))
		if diagnostic == "" {
			return nil, fmt.Errorf("git %s: %w", strings.Join(args, " "), err)
		}
		return nil, fmt.Errorf("git %s: %w: %s", strings.Join(args, " "), err, diagnostic)
	}
	return output, nil
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
	} else if opts.discardSuperseded {
		mode = "discard-superseded"
	} else if opts.carryForward {
		mode = "carry-forward"
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
		} else {
			disposition, dispositionErr := daemon.ClassifySupersededBranch(ctx, run, runs.coverage)
			applyReconcileBranchDisposition(&result, opts, disposition, dispositionErr)
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

func applyReconcileBranchDisposition(
	result *reconcileResult,
	opts reconcileOptions,
	disposition daemon.BranchDisposition,
	dispositionErr error,
) {
	if result == nil {
		return
	}
	result.disposition = disposition
	if dispositionErr != nil {
		if opts.discardSuperseded {
			result.Action = "preserve"
			result.RefusalReason = dispositionErr.Error()
		}
		return
	}
	if !disposition.Superseded {
		if opts.discardSuperseded && strings.TrimSpace(disposition.RefusalReason) != "" {
			result.Action = "preserve"
			result.RefusalReason = disposition.RefusalReason
		}
		return
	}
	// Keep legacy --apply classification and cleanup unchanged. The named
	// disposition reports only branches that actually held Run commits.
	if opts.apply || (!opts.discardSuperseded && len(disposition.Commits) == 0) {
		return
	}
	result.Classification = string(runworktree.ReconciliationSuperseded)
	result.Evidence = disposition.Reason
	result.RefusalReason = ""
	if opts.discardSuperseded {
		result.Action = "discard after writing branch record"
	} else {
		result.Action = "would discard with --discard-superseded"
	}
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
				if branch == classification.Current || slices.Contains(classification.Preserved, branch) {
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
			if slices.Contains(classification.Releasable, branch) {
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

func discardSupersededReconcileResults(
	ctx context.Context,
	artifactRoot string,
	report *reconcileReport,
) int {
	if report == nil {
		return 0
	}
	refusals := 0
	for index := range report.Results {
		result := &report.Results[index]
		disposition := result.disposition
		if !disposition.Superseded {
			result.Action = "preserve"
			if strings.TrimSpace(result.RefusalReason) == "" {
				result.RefusalReason = disposition.RefusalReason
			}
			refusals++
			continue
		}
		recordPath := filepath.Join(artifactRoot, "runs", result.RunID, "branch-disposition.json")
		if err := daemon.RecordAndDiscardSupersededBranch(ctx, recordPath, disposition, time.Now().UTC()); err != nil {
			result.Action = "preserve"
			if _, statErr := os.Stat(recordPath); statErr == nil {
				result.Action += "; branch record written at " + recordPath
			}
			result.RefusalReason = err.Error()
			report.Summary.OperationalFailures++
			continue
		}
		result.Action = "discarded"
		result.RefusalReason = ""
		result.Evidence = disposition.Reason + "; branch record: " + recordPath
		report.Summary.Applied++
	}
	return refusals
}

func applyReconcileReport(
	ctx context.Context,
	homeDir string,
	opts reconcileOptions,
	report *reconcileReport,
) {
	applyReconcileProcesses(ctx, report)
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

}

func applyReconcileProcesses(ctx context.Context, report *reconcileReport) {
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
	for _, candidate := range report.CarryForwards {
		fmt.Fprintf(&output, "Carry-forward candidate: Task %s (Run %s)\n", candidate.TaskID, candidate.RunID)
		fmt.Fprintf(&output, "  commit: %s\n", textReconcileValue(candidate.Commit))
		fmt.Fprintf(&output, "  task-file: %s\n", textReconcileValue(candidate.TaskFile))
		fmt.Fprintf(&output, "  inputs-moved: %t\n", candidate.InputsMoved)
		fmt.Fprintf(&output, "  moved-inputs: %v\n", candidate.MovedInputs)
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
