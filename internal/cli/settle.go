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

Re-runs one Task's Verification commands in its kept Task Worktree when
available, then its kept Run Worktree, then the current repository. A surface
settles when its task file is failed, or completed while the verified work is
still uncommitted there — the state a refusing commit hook leaves behind. On
pass, settles completed and commits all Run Worktree changes plus the task file;
for a kept Task Worktree, commits all Task Worktree changes plus the task file
and integrates them onto the Run Branch. settle creates no Run, writes no Run
Event Journal entries, and never pushes.

Options:
  --spec  Spec slug under docs/specs/
  --task  Task id from the Spec Task Graph

Exit codes:
  0  settled completed and committed
  1  verification did not pass; the Task status is unchanged and no commit is created
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
	detail string
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

	collaborators := commandDependenciesForContext(ctx).newEngineCollaborators()
	fmt.Fprintf(stderr, "Settle surface: %s\n", plan.workDir)
	if !runVerificationInSettleSurface(ctx, plan, collaborators.verifier, stdout, stderr) {
		return exitRunFailed
	}

	commitResult, err := settleTaskAndCommit(ctx, plan, collaborators)
	if err != nil {
		fmt.Fprintf(stderr, "%s: settle failed after verification: %v\n", app.Name, err)
		return exitRunFailed
	}
	shortSHA := commitResult.shortSHA
	if plan.taskSurface {
		if integratedSHA, err := integrateTaskWorktree(ctx, plan); err != nil {
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
		for _, staged := range commitResult.paths {
			if staged.deleted {
				fmt.Fprintf(stdout, "commit %s — deleted\n", staged.path)
				continue
			}
			fmt.Fprintf(stdout, "commit %s\n", staged.path)
		}
	}
	fmt.Fprintf(stdout, "settled %s completed — %s\n", plan.task.ID, shortSHA)
	return exitOK
}

// runVerificationInSettleSurface re-runs the Task's Verification in the
// selected settle surface and reports one verdict line per command. The work
// settle recovers already passed this Verification under the Daemon, so settle
// runs the same commands through the same daemon.Verifier collaborator the
// Implement Run engine uses: task file order, verbatim text, no edit, and no
// Agent session to repair a failure. The first command that does not pass ends
// the loop, and the caller leaves the surface exactly as it found it.
//
// Implement's single Verification repair (ADR-0038) has no counterpart here. A
// temporary failure (exit code 75) reaches settle wrapped around its command
// failure and stops the settle like any other failed command; the Supervisor
// re-runs settle once the cause is gone.
func runVerificationInSettleSurface(ctx context.Context, plan settlePlan, verifier daemon.Verifier, stdout io.Writer, stderr io.Writer) bool {
	outputPath := daemon.VerificationOutputPath(plan.artifactDir, settleVerificationRunID(plan), 1, 1)
	for _, verificationCommand := range plan.task.Verification {
		_, err := verifier.Verify(ctx, daemon.VerifyRequest{
			WorkDir:    plan.workDir,
			Command:    verificationCommand,
			OutputPath: outputPath,
		})
		if err == nil {
			fmt.Fprintf(stdout, "verify %s — ok\n", verificationCommand)
			continue
		}
		var commandErr *daemon.VerificationCommandError
		if errors.As(err, &commandErr) {
			fmt.Fprintf(stdout, "verify %s — failed (diagnostics: %s)\n", verificationCommand, commandErr.OutputPath)
			fmt.Fprintf(stdout, "%s stays %s — verification failed\n", plan.task.ID, plan.task.Status)
			return false
		}
		// The runner never observed a verdict, so the recovered work is
		// neither verified nor refuted. Report it apart from a failed command
		// and keep the cause on stderr: it names an environment problem the
		// Supervisor clears before settling again.
		var unknownErr *daemon.VerificationUnknownError
		if errors.As(err, &unknownErr) {
			fmt.Fprintf(stdout, "verify %s — verdict unknown\n", verificationCommand)
			fmt.Fprintf(stdout, "%s stays %s — verification verdict unknown\n", plan.task.ID, plan.task.Status)
			fmt.Fprintf(stderr, "%s: settle verification verdict unknown: %v\n", app.Name, unknownErr)
			return false
		}
		fmt.Fprintf(stderr, "%s: settle verification failed: %v\n", app.Name, err)
		return false
	}
	return true
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
		choice, reports, err := resolveSettleSurface(ctx, req, resolvedSpecsRoot, gitState.Root, candidatesWithCurrent(candidates, gitState.Root, gitState.Branch))
		if err != nil {
			return settlePlan{}, err
		}
		if !choice.found {
			return settlePlan{}, settleNoSurfaceError(req.specSlug, req.taskID, reports)
		}
		applySettleSurfaceChoice(&plan, choice)
	}
	if plan.graph == nil {
		choice, reports, err := resolveSettleSurface(ctx, req, resolvedSpecsRoot, gitState.Root, candidatesWithCurrent(nil, gitState.Root, gitState.Branch))
		if err != nil {
			return settlePlan{}, err
		}
		if !choice.found {
			return settlePlan{}, settleNoSurfaceError(req.specSlug, req.taskID, reports)
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

func resolveSettleSurface(ctx context.Context, req settleRequest, resolvedSpecsRoot roundconfig.SpecsRoot, gitRoot string, candidates []settleSurfaceCandidate) (settleSurfaceChoice, []settleSurfaceReport, error) {
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
		settles, detail, err := settleSurfaceSettles(ctx, candidate.workDir, report.status)
		if err != nil {
			return settleSurfaceChoice{}, reports, err
		}
		report.detail = detail
		reports = append(reports, report)
		if settles {
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

// settleSurfaceSettles reports whether a candidate surface holds work settle
// recovers, and what to say about the surface when it does not. taskStatus is
// the status the surface's own task file carries, in the canonical text
// spec.Load normalizes to — the same text the candidate report enumerates.
//
// failed is the ordinary settle target. completed is the hook-refusal state:
// the Daemon settles the Task completed once its Verification passes and then
// commits, so a commit a repository hook refuses leaves verified work
// uncommitted in the surface that produced it. A completed Task whose surface
// holds nothing uncommitted was committed already and settles nothing.
func settleSurfaceSettles(ctx context.Context, workDir string, taskStatus string) (bool, string, error) {
	if taskStatus == "failed" {
		return true, "", nil
	}
	if taskStatus == "completed" {
		changed, err := commandDependenciesForContext(ctx).newEngineCollaborators().worktree.Snapshot(ctx, workDir)
		if err != nil {
			return false, "", validationError{message: fmt.Sprintf("inspect settle surface %q for uncommitted work: %v", workDir, err)}
		}
		if len(changed) == 0 {
			return false, settleCommittedSurfaceDetail, nil
		}
		return true, "", nil
	}
	return false, "", nil
}

const settleCommittedSurfaceDetail = "no uncommitted work"

func settleNoSurfaceError(slug string, taskID string, reports []settleSurfaceReport) error {
	return validationError{message: fmt.Sprintf("Task %s has no settleable surface; candidates: %s; %s", taskID, formatSettleSurfaceReports(reports), settleStatusGuidance(slug, taskID, reports))}
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
		if detail := strings.TrimSpace(report.detail); detail != "" {
			status += " (" + detail + ")"
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
	return fmt.Sprintf("Task %s status is unavailable; %s", taskID, settleStatusRequirement)
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

// settleStatusRequirement states the contract every refusal closes on, so a
// Supervisor reading one refusal knows which two states settle recovers.
const settleStatusRequirement = "settle requires failed, or completed with uncommitted work"

func settleStatusMessage(slug string, task spec.Task) string {
	switch task.Status {
	case spec.StatusPending, spec.StatusInProgress:
		return fmt.Sprintf("Task %s status is %s; run the Implement Command instead: roundfix implement --spec %s", task.ID, task.Status, slug)
	case spec.StatusCompleted:
		return fmt.Sprintf("Task %s status is completed with no uncommitted work; nothing to settle", task.ID)
	default:
		return fmt.Sprintf("Task %s status is %s; %s", task.ID, task.Status, settleStatusRequirement)
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
	paths    []settleStagedPath
}

// settleStagedPath is one path the settle commit carries and whether the
// commit removes it. A removal is the one entry the path alone cannot report:
// the file is gone from the surface, so the settle report says so rather than
// leaving a Supervisor to read the commit to find out what it deleted.
type settleStagedPath struct {
	path    string
	deleted bool
}

// settleTaskAndCommit settles the Task completed in the surface that holds its
// verified work and commits that surface as one standard Task commit. The
// status flip is written before staging so it rides in the same commit as the
// code changes (ADR 0013), the way the Daemon's own Task commit carries it.
func settleTaskAndCommit(ctx context.Context, plan settlePlan, collaborators engineCollaborators) (settleCommitResult, error) {
	taskPath := filepath.Join(plan.specsRoot, plan.task.File)
	if err := spec.SetStatus(taskPath, spec.StatusCompleted); err != nil {
		return settleCommitResult{}, fmt.Errorf("settle Task %s completed: %w", plan.task.ID, err)
	}
	if err := addAllChanges(ctx, plan.workDir); err != nil {
		return settleCommitResult{}, err
	}
	committedPaths, err := stagedSettlePaths(ctx, plan.workDir)
	if err != nil {
		return settleCommitResult{}, err
	}
	if len(committedPaths) > 0 {
		// The committer stages the surface again before it commits. That
		// repeats the stage-all above rather than narrowing it, so the commit
		// carries exactly the index stagedSettlePaths just reported.
		if err := collaborators.committer.Commit(ctx, daemon.CommitRequest{
			WorkDir: plan.workDir,
			Message: daemon.TaskCommitMessage(plan.graph.Spec.Slug, plan.task),
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

// addAllChanges stages every change in the settle surface with `git add --all`.
//
// The Daemon's Task commit stages an explicit path list, because it diffs a
// before and an after snapshot to keep the commit to what one Agent turn
// touched. A settle surface has no before snapshot: everything uncommitted
// there is the recovered Task's work, and two shapes reach this call that a
// path list cannot carry. A file the Task deleted stops matching its own
// pathspec once the removal is staged, so re-staging it fails on a pathspec
// that matches nothing. A file the Task renamed is one porcelain record naming
// the new path, so staging that path alone leaves the old path's deletion out
// of the commit. Staging the whole surface records both as removals.
func addAllChanges(ctx context.Context, workDir string) error {
	if _, err := (preflight.ExecGitRunner{}).RunGit(ctx, workDir, "add", "--all"); err != nil {
		return fmt.Errorf("stage settle surface %q: %w", workDir, err)
	}
	return nil
}

// stagedSettlePaths reports what addAllChanges staged, so the settle report
// names exactly the paths the Task commit carries and marks the ones it
// removes. Rename detection stays off: a rename is an added path and a removed
// path in the commit, and settle reports both.
//
// `--name-status -z` writes one status field and one path field per entry,
// each NUL-terminated, so the fields read as pairs and the split leaves one
// empty field after the last path. A status without its path is truncated
// output, and settle refuses rather than report a commit it cannot enumerate.
func stagedSettlePaths(ctx context.Context, workDir string) ([]settleStagedPath, error) {
	output, err := (preflight.ExecGitRunner{}).RunGit(ctx, workDir, "diff", "--cached", "--name-status", "--no-renames", "-z")
	if err != nil {
		return nil, fmt.Errorf("read staged paths in settle surface %q: %w", workDir, err)
	}
	fields := strings.Split(output, "\x00")
	var paths []settleStagedPath
	for index := 0; index < len(fields); index += 2 {
		status := fields[index]
		if status == "" {
			continue
		}
		if index+1 >= len(fields) || fields[index+1] == "" {
			return nil, fmt.Errorf("read staged paths in settle surface %q: status %q carries no path", workDir, status)
		}
		paths = append(paths, settleStagedPath{
			path:    fields[index+1],
			deleted: strings.HasPrefix(status, "D"),
		})
	}
	sort.Slice(paths, func(first int, second int) bool {
		return paths[first].path < paths[second].path
	})
	return paths, nil
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

// integrateTaskWorktree replays the settled Task commit onto the Run Branch and
// removes the Task Worktree once it lands. A conflict returns before the
// removal, so the surface that still holds the work stays where the Supervisor
// can reach it.
func integrateTaskWorktree(ctx context.Context, plan settlePlan) (string, error) {
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
	if err := commandDependenciesForContext(ctx).cleanupCleanRunWorktree(ctx, ref); err != nil {
		return "", err
	}
	return "", nil
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
