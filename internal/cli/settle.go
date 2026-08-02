package cli

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"roundfix/internal/app"
	roundconfig "roundfix/internal/config"
	"roundfix/internal/daemon"
	"roundfix/internal/preflight"
	"roundfix/internal/spec"
	"roundfix/internal/store"
	runworktree "roundfix/internal/worktree"
)

const settleUsage = `Usage:
  roundfix settle --spec <slug> --task <task_id>

Re-runs one failed Task's Verification commands in its kept Task Worktree when
available, then its kept Run Worktree, then the current repository. On pass,
settles completed and commits all Run Worktree changes plus the task file; for
a kept Task Worktree, commits all Task Worktree changes plus the task file and
integrates them onto the Run Branch. settle creates no Run, writes no Run Event
Journal entries, and never pushes.

Options:
  --spec  Spec slug under docs/specs/
  --task  Task id from the Spec Task Graph

Exit codes:
  0  settled completed and committed
  1  verification failed; status stays failed and no commit is created
  2  Preflight Validation failed
`

type settleRequest struct {
	specSlug string
	taskID   string
}

type settlePlan struct {
	userRoot     string
	workDir      string
	targetBranch string
	homeDir      string
	artifactDir  string
	specsRoot    string
	graph        *spec.Graph
	task         spec.Task
	run          store.Run
	runRef       runworktree.Ref
	taskRef      runworktree.TaskRef
	hasRun       bool
	taskSurface  bool
}

type settleSurfaceCandidate struct {
	workDir      string
	targetBranch string
	run          store.Run
	runRef       runworktree.Ref
	taskRef      runworktree.TaskRef
	hasRun       bool
	taskSurface  bool
	current      bool
}

type settleSurfaceChoice struct {
	candidate settleSurfaceCandidate
	specsRoot string
	graph     *spec.Graph
	task      spec.Task
	found     bool
}

type settleSurfaceReport struct {
	path   string
	status string
}

func runSettleCommand(ctx context.Context, args []string, stdout, stderr io.Writer, environment commandEnvironment) int {
	if commandWantsHelp(args) {
		fmt.Fprint(stdout, settleUsage)
		return exitOK
	}
	req, err := parseSettleCommand(args)
	if err != nil {
		printPreflightFailure("settle", err, stderr)
		return exitPreflight
	}
	plan, err := preflightSettle(ctx, req, stderr, environment)
	if err != nil {
		printPreflightFailure("settle", err, stderr)
		return exitPreflight
	}

	collaborators := newEngineCollaborators()
	verificationRunID := settleVerificationRunID(plan)
	fmt.Fprintf(stderr, "Settle surface: %s\n", plan.workDir)
	for _, command := range plan.task.Verification {
		_, err := collaborators.verifier.Verify(ctx, daemon.VerifyRequest{
			WorkDir:    plan.workDir,
			Command:    command,
			OutputPath: daemon.VerificationOutputPath(plan.artifactDir, verificationRunID, 1, 1),
		})
		if err != nil {
			var commandErr *daemon.VerificationCommandError
			if !errors.As(err, &commandErr) {
				fmt.Fprintf(stderr, "%s: settle verification failed: %v\n", app.Name, err)
				return exitRunFailed
			}
			fmt.Fprintf(stdout, "verify %s — failed (diagnostics: %s)\n", command, commandErr.OutputPath)
			fmt.Fprintf(stdout, "%s stays failed — verification failed\n", plan.task.ID)
			return exitRunFailed
		}
		fmt.Fprintf(stdout, "verify %s — ok\n", command)
	}

	commitResult, err := settleTaskAndCommit(ctx, plan, collaborators)
	if err != nil {
		fmt.Fprintf(stderr, "%s: settle failed after verification: %v\n", app.Name, err)
		return exitRunFailed
	}
	shortSHA := commitResult.shortSHA
	if plan.taskSurface {
		if integratedSHA, err := integrateSettledTaskWorktree(ctx, plan); err != nil {
			fmt.Fprintf(stderr, "%s: settle failed after verification: %v\n", app.Name, err)
			return exitRunFailed
		} else if integratedSHA != "" {
			shortSHA = integratedSHA
		}
	}
	if plan.hasRun {
		integrationCommand, err := integrateSettledRun(ctx, plan)
		if err != nil {
			fmt.Fprintf(stderr, "%s: settle failed after verification: %v\n", app.Name, err)
			return exitRunFailed
		}
		if integrationCommand != "" {
			fmt.Fprintf(stdout, "integration pending — %s\n", integrationCommand)
			return exitRunFailed
		}
	}
	if len(commitResult.paths) > 0 {
		emitSettleSharedFailedWarning(stderr, plan)
		for _, path := range commitResult.paths {
			fmt.Fprintf(stdout, "commit %s\n", path)
		}
	}
	fmt.Fprintf(stdout, "settled %s completed — %s\n", plan.task.ID, shortSHA)
	return exitOK
}

func settleVerificationRunID(plan settlePlan) string {
	replacer := strings.NewReplacer("/", "_", "\\", "_")
	return "settle-" + replacer.Replace(plan.graph.Spec.Slug) + "-" + replacer.Replace(plan.task.ID)
}

func parseSettleCommand(args []string) (settleRequest, error) {
	var req settleRequest
	fs := flag.NewFlagSet("settle", flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	fs.StringVar(&req.specSlug, "spec", "", "Spec slug under docs/specs/")
	fs.StringVar(&req.taskID, "task", "", "Task id from the Spec Task Graph")
	if err := fs.Parse(args); err != nil {
		return req, validationError{message: err.Error()}
	}
	if remaining := fs.Args(); len(remaining) > 0 {
		return req, validationError{message: fmt.Sprintf("unexpected argument %q", remaining[0])}
	}
	req.specSlug = strings.TrimSpace(req.specSlug)
	req.taskID = strings.TrimSpace(req.taskID)
	if req.specSlug == "" {
		return req, validationError{message: "missing required --spec; pass --spec <slug>"}
	}
	if req.taskID == "" {
		return req, validationError{message: "missing required --task; pass --task <task_id>"}
	}
	return req, nil
}

func preflightSettle(ctx context.Context, req settleRequest, stderr io.Writer, environment commandEnvironment) (settlePlan, error) {
	loadedConfig, err := loadCommandConfig(environment, stderr)
	if err != nil {
		return settlePlan{}, err
	}
	gitWorkDir := loadedConfig.GitRoot
	if strings.TrimSpace(gitWorkDir) == "" {
		gitWorkDir = environment.workDir
	}
	gitState, err := preflight.InspectGit(ctx, gitWorkDir, nil)
	if err != nil {
		return settlePlan{}, validationError{message: fmt.Sprintf("repository resolves: %v", err)}
	}
	plan := settlePlan{
		userRoot:     gitState.Root,
		workDir:      gitState.Root,
		targetBranch: gitState.Branch,
		homeDir:      loadedConfig.HomeDir,
	}
	artifactDir, err := roundconfig.ValidateArtifactDirectory(loadedConfig.Config.Defaults.ArtifactDir, gitState.Root, loadedConfig.HomeDir)
	if err != nil {
		return settlePlan{}, validationError{message: fmt.Sprintf("Artifact Directory resolves: %v", err)}
	}
	plan.artifactDir = artifactDir
	resolvedSpecsRoot, err := roundconfig.ResolveSpecsRoot(loadedConfig, gitState.Root)
	if err != nil {
		return settlePlan{}, err
	}
	if run, found, err := latestKeptSettleRun(ctx, loadedConfig.HomeDir, gitState.Root, req.specSlug); err != nil {
		return settlePlan{}, err
	} else if found {
		runRef := runworktree.Ref{
			RunID:    run.ID,
			Path:     run.WorkDir,
			Branch:   runworktree.BranchName(run.ID),
			UserRoot: gitState.Root,
		}
		taskRef, err := runworktree.TaskRefFor(runRef, req.taskID)
		if err != nil {
			return settlePlan{}, err
		}
		plan.run = run
		plan.runRef = runRef
		candidates := []settleSurfaceCandidate{
			{
				workDir:      taskRef.Path,
				targetBranch: run.LocalBranch,
				run:          run,
				runRef:       runRef,
				taskRef:      taskRef,
				hasRun:       true,
				taskSurface:  true,
			},
			{
				workDir:      run.WorkDir,
				targetBranch: run.LocalBranch,
				run:          run,
				runRef:       runRef,
				hasRun:       true,
			},
		}
		choice, reports, err := resolveSettleSurface(req, resolvedSpecsRoot, gitState.Root, candidatesWithCurrent(candidates, gitState.Root, gitState.Branch))
		if err != nil {
			return settlePlan{}, err
		}
		if !choice.found {
			return settlePlan{}, settleNoFailedSurfaceError(req.specSlug, req.taskID, reports)
		}
		applySettleSurfaceChoice(&plan, choice)
	}
	if plan.graph == nil {
		choice, reports, err := resolveSettleSurface(req, resolvedSpecsRoot, gitState.Root, candidatesWithCurrent(nil, gitState.Root, gitState.Branch))
		if err != nil {
			return settlePlan{}, err
		}
		if !choice.found {
			return settlePlan{}, settleNoFailedSurfaceError(req.specSlug, req.taskID, reports)
		}
		applySettleSurfaceChoice(&plan, choice)
	}
	if err := ensureNoSettleActiveRun(ctx, loadedConfig.HomeDir, gitState.Root, req.specSlug, stderr); err != nil {
		return settlePlan{}, err
	}
	return plan, nil
}

func candidatesWithCurrent(candidates []settleSurfaceCandidate, gitRoot string, branch string) []settleSurfaceCandidate {
	result := append([]settleSurfaceCandidate(nil), candidates...)
	result = append(result, settleSurfaceCandidate{
		workDir:      gitRoot,
		targetBranch: branch,
		current:      true,
	})
	return result
}

func applySettleSurfaceChoice(plan *settlePlan, choice settleSurfaceChoice) {
	candidate := choice.candidate
	plan.workDir = candidate.workDir
	plan.targetBranch = candidate.targetBranch
	plan.specsRoot = choice.specsRoot
	plan.graph = choice.graph
	plan.task = choice.task
	plan.run = candidate.run
	plan.runRef = candidate.runRef
	plan.taskRef = candidate.taskRef
	plan.hasRun = candidate.hasRun
	plan.taskSurface = candidate.taskSurface
}

func resolveSettleSurface(req settleRequest, resolvedSpecsRoot roundconfig.SpecsRoot, gitRoot string, candidates []settleSurfaceCandidate) (settleSurfaceChoice, []settleSurfaceReport, error) {
	var reports []settleSurfaceReport
	for _, candidate := range candidates {
		report := settleSurfaceReport{path: candidate.workDir}
		ok, err := settlePathExists(candidate.workDir)
		if err != nil {
			return settleSurfaceChoice{}, reports, err
		}
		if !ok {
			report.status = "path does not exist"
			reports = append(reports, report)
			continue
		}
		specsRoot := specsRootForWorkDir(resolvedSpecsRoot, gitRoot, candidate.workDir)
		graph, err := spec.Load(specsRoot, req.specSlug)
		if err != nil {
			if candidate.current {
				return settleSurfaceChoice{}, reports, validationError{message: fmt.Sprintf("Spec loads valid: %v", err)}
			}
			report.status = "Spec loads invalid: " + err.Error()
			reports = append(reports, report)
			continue
		}
		task, found := findSettleTask(graph.Tasks, req.taskID)
		if !found {
			if candidate.current {
				return settleSurfaceChoice{}, reports, validationError{message: fmt.Sprintf("Task %q does not exist in Task Graph for Spec %q; choose a task id from docs/specs/%s/_tasks.md", req.taskID, req.specSlug, req.specSlug)}
			}
			report.status = "task not found"
			reports = append(reports, report)
			continue
		}
		report.status = string(task.Status)
		reports = append(reports, report)
		if task.Status == spec.StatusFailed {
			return settleSurfaceChoice{
				candidate: candidate,
				specsRoot: specsRoot,
				graph:     graph,
				task:      task,
				found:     true,
			}, reports, nil
		}
	}
	return settleSurfaceChoice{}, reports, nil
}

func settleNoFailedSurfaceError(slug string, taskID string, reports []settleSurfaceReport) error {
	return validationError{message: fmt.Sprintf("Task %s has no failed settle surface; candidates: %s; %s", taskID, formatSettleSurfaceReports(reports), settleStatusGuidance(slug, taskID, reports))}
}

func formatSettleSurfaceReports(reports []settleSurfaceReport) string {
	parts := make([]string, 0, len(reports))
	for _, report := range reports {
		status := strings.TrimSpace(report.status)
		if status == "" {
			status = "status unavailable"
		} else if spec.AllowedStatus(spec.Status(status)) {
			status = "status " + status
		}
		parts = append(parts, fmt.Sprintf("%s: %s", report.path, status))
	}
	return strings.Join(parts, "; ")
}

func settleStatusGuidance(slug string, taskID string, reports []settleSurfaceReport) string {
	for _, report := range reports {
		status := spec.Status(strings.TrimSpace(report.status))
		if spec.AllowedStatus(status) {
			return settleStatusMessage(slug, spec.Task{ID: taskID, Status: status})
		}
	}
	return fmt.Sprintf("Task %s status is unavailable; settle requires failed", taskID)
}

func latestKeptSettleRun(ctx context.Context, homeDir string, gitRoot string, specSlug string) (store.Run, bool, error) {
	if _, err := os.Stat(store.DatabasePath(homeDir)); errors.Is(err, os.ErrNotExist) {
		return store.Run{}, false, nil
	} else if err != nil {
		return store.Run{}, false, validationError{message: fmt.Sprintf("inspect Run Database before settle: %v", err)}
	}
	runStore, err := store.OpenReader(ctx, homeDir)
	if err != nil {
		return store.Run{}, false, validationError{message: fmt.Sprintf("open Run Database before settle: %v", err)}
	}
	defer func() {
		_ = runStore.Close()
	}()
	run, found, err := runStore.LatestKeptSpecRun(ctx, gitRoot, specSlug)
	if err != nil {
		return store.Run{}, false, validationError{message: fmt.Sprintf("find kept Run Worktree for Spec: %v", err)}
	}
	return run, found, nil
}

func settlePathExists(path string) (bool, error) {
	path = strings.TrimSpace(path)
	if path == "" {
		return false, nil
	}
	info, err := os.Stat(path)
	if errors.Is(err, os.ErrNotExist) {
		return false, nil
	}
	if err != nil {
		return false, validationError{message: fmt.Sprintf("settle surface %q is unavailable: %v", path, err)}
	}
	if !info.IsDir() {
		return false, validationError{message: fmt.Sprintf("settle surface %q is not a directory", path)}
	}
	return true, nil
}

func findSettleTask(tasks []spec.Task, taskID string) (spec.Task, bool) {
	for _, task := range tasks {
		if task.ID == taskID {
			return task, true
		}
	}
	return spec.Task{}, false
}

func settleStatusError(slug string, task spec.Task) error {
	return validationError{message: settleStatusMessage(slug, task)}
}

func settleStatusMessage(slug string, task spec.Task) string {
	switch task.Status {
	case spec.StatusPending, spec.StatusInProgress:
		return fmt.Sprintf("Task %s status is %s; run the Implement Command instead: roundfix implement --spec %s", task.ID, task.Status, slug)
	case spec.StatusCompleted:
		return fmt.Sprintf("Task %s status is completed; nothing to do", task.ID)
	default:
		return fmt.Sprintf("Task %s status is %s; settle requires failed", task.ID, task.Status)
	}
}

func ensureNoSettleActiveRun(ctx context.Context, homeDir string, gitRoot string, specSlug string, stderr io.Writer) error {
	if _, err := os.Stat(store.DatabasePath(homeDir)); errors.Is(err, os.ErrNotExist) {
		return nil
	} else if err != nil {
		return validationError{message: fmt.Sprintf("inspect Run Database before settle: %v", err)}
	}
	runStore, err := store.Open(ctx, homeDir)
	if err != nil {
		return validationError{message: fmt.Sprintf("open Run Database before settle: %v", err)}
	}
	defer func() {
		_ = runStore.Close()
	}()
	if active, found, err := runStore.ActiveSpecRun(ctx, gitRoot, specSlug); err != nil {
		return validationError{message: fmt.Sprintf("check Active Run for Spec target: %v", err)}
	} else if found {
		if _, reclaimed, err := reclaimOrphanedActiveRun(ctx, runStore, active, stderr); err != nil {
			return err
		} else if reclaimed {
			return nil
		}
		return validationError{message: fmt.Sprintf("Active Run %s already holds Spec target %q in working tree %q; stop it with: roundfix stop %s", active.ID, specSlug, gitRoot, active.ID)}
	}
	if active, found, err := runStore.ActiveRunInGitRoot(ctx, gitRoot); err != nil {
		return validationError{message: fmt.Sprintf("check Active Run for working tree: %v", err)}
	} else if found {
		if _, reclaimed, err := reclaimOrphanedActiveRun(ctx, runStore, active, stderr); err != nil {
			return err
		} else if reclaimed {
			return nil
		}
		return validationError{message: fmt.Sprintf("Active Run %s already holds working tree %q; stop it with: roundfix stop %s", active.ID, gitRoot, active.ID)}
	}
	return nil
}

type settleCommitResult struct {
	shortSHA string
	paths    []string
}

func settleTaskAndCommit(ctx context.Context, plan settlePlan, collaborators engineCollaborators) (settleCommitResult, error) {
	taskPath := filepath.Join(plan.specsRoot, plan.task.File)
	changed, err := collaborators.worktree.Snapshot(ctx, plan.workDir)
	if err != nil {
		return settleCommitResult{}, err
	}
	changed = ensureSettleCommitPath(changed, settleArtifactCommitPath(plan, taskPath))
	stageable, _ := daemon.FilterStageablePaths(plan.workDir, changed)
	if err := spec.SetStatus(taskPath, spec.StatusCompleted); err != nil {
		return settleCommitResult{}, fmt.Errorf("settle Task %s completed: %w", plan.task.ID, err)
	}
	committedPaths := []string(nil)
	if len(stageable) > 0 {
		committedPaths = append([]string(nil), stageable...)
		message := daemon.TaskCommitMessage(plan.graph.Spec.Slug, plan.task)
		if err := collaborators.committer.Commit(ctx, daemon.CommitRequest{
			WorkDir: plan.workDir,
			Message: message,
			Paths:   stageable,
		}); err != nil {
			return settleCommitResult{}, err
		}
	}
	shortSHA, err := settleShortHEAD(ctx, plan.workDir)
	if err != nil {
		return settleCommitResult{}, err
	}
	return settleCommitResult{shortSHA: shortSHA, paths: committedPaths}, nil
}

func emitSettleSharedFailedWarning(stderr io.Writer, plan settlePlan) {
	failed := otherFailedSettleTasks(plan.graph.Tasks, plan.task.ID)
	if len(failed) == 0 {
		return
	}
	fmt.Fprintf(stderr, "%s: warning: other failed Tasks in Spec %q may have work included in this settle commit: %s\n", app.Name, plan.graph.Spec.Slug, strings.Join(failed, ", "))
}

func otherFailedSettleTasks(tasks []spec.Task, settledTaskID string) []string {
	var failed []string
	for _, task := range tasks {
		if task.ID != settledTaskID && task.Status == spec.StatusFailed {
			failed = append(failed, task.ID)
		}
	}
	sort.Strings(failed)
	return failed
}

func integrateSettledTaskWorktree(ctx context.Context, plan settlePlan) (string, error) {
	result, err := runworktree.IntegrateTask(ctx, plan.runRef, plan.taskRef)
	if err != nil {
		return "", err
	}
	if result.Mode == runworktree.ModeTaskConflict {
		return "", fmt.Errorf("task worktree integration conflict on %s", result.Reason)
	}
	shortSHA, err := settleShortHEAD(ctx, plan.runRef.Path)
	if err != nil {
		return "", fmt.Errorf("read run branch HEAD after task worktree integration: %w", err)
	}
	if err := runworktree.CleanupTask(ctx, plan.taskRef); err != nil {
		return "", err
	}
	return shortSHA, nil
}

func integrateSettledRun(ctx context.Context, plan settlePlan) (string, error) {
	ref := plan.runRef
	result, err := integrateCleanImplementRun(ctx, ref, plan.targetBranch)
	if err != nil {
		return "", err
	}
	if result.Mode == runworktree.ModePending {
		return implementIntegrationCommand(ref), nil
	}
	if err := cleanupCleanRunWorktree(ctx, ref); err != nil {
		return "", err
	}
	return "", nil
}

func ensureSettleCommitPath(paths []string, path string) []string {
	for _, existing := range paths {
		if existing == path {
			result := append([]string(nil), paths...)
			sort.Strings(result)
			return result
		}
	}
	result := append(append([]string(nil), paths...), path)
	sort.Strings(result)
	return result
}

func settleShortHEAD(ctx context.Context, gitRoot string) (string, error) {
	output, err := preflight.ExecGitRunner{}.RunGit(ctx, gitRoot, "rev-parse", "--short", "HEAD")
	if err != nil {
		return "", err
	}
	shortSHA := strings.TrimSpace(output)
	if shortSHA == "" {
		return "", errors.New("git returned an empty short SHA")
	}
	return shortSHA, nil
}

func filepathInRoot(root string, path string) string {
	return filepath.Join(root, path)
}

func settleArtifactCommitPath(plan settlePlan, artifactPath string) string {
	if relative, err := filepath.Rel(plan.workDir, artifactPath); err == nil && relative != ".." && !strings.HasPrefix(relative, ".."+string(filepath.Separator)) && !filepath.IsAbs(relative) {
		return relative
	}
	return artifactPath
}
