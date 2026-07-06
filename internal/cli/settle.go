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
	graph        *spec.Graph
	task         spec.Task
	run          store.Run
	runRef       runworktree.Ref
	taskRef      runworktree.TaskRef
	hasRun       bool
	taskSurface  bool
}

func runSettleCommand(ctx context.Context, args []string, stdout, stderr io.Writer) int {
	if commandWantsHelp(args) {
		fmt.Fprint(stdout, settleUsage)
		return exitOK
	}
	req, err := parseSettleCommand(args)
	if err != nil {
		printPreflightFailure("settle", err, stderr)
		return exitPreflight
	}
	plan, err := preflightSettle(ctx, req, stderr)
	if err != nil {
		printPreflightFailure("settle", err, stderr)
		return exitPreflight
	}

	collaborators := newEngineCollaborators()
	for _, command := range plan.task.Verification {
		if err := collaborators.verifier.Verify(ctx, daemon.VerifyRequest{
			WorkDir: plan.workDir,
			Command: command,
			Stream:  stderr,
		}); err != nil {
			fmt.Fprintf(stdout, "verify %s — failed\n", command)
			fmt.Fprintf(stdout, "%s stays failed — verification failed\n", plan.task.ID)
			return exitRunFailed
		}
		fmt.Fprintf(stdout, "verify %s — ok\n", command)
	}

	shortSHA, err := settleTaskAndCommit(ctx, plan, collaborators)
	if err != nil {
		fmt.Fprintf(stderr, "%s: settle failed after verification: %v\n", app.Name, err)
		return exitRunFailed
	}
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
	fmt.Fprintf(stdout, "settled %s completed — %s\n", plan.task.ID, shortSHA)
	return exitOK
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

func preflightSettle(ctx context.Context, req settleRequest, stderr io.Writer) (settlePlan, error) {
	loadedConfig, err := roundconfig.Load(roundconfig.LoadOptions{Stderr: stderr})
	if err != nil {
		return settlePlan{}, err
	}
	gitState, err := preflight.InspectGit(ctx, loadedConfig.GitRoot, nil)
	if err != nil {
		return settlePlan{}, validationError{message: fmt.Sprintf("repository resolves: %v", err)}
	}
	plan := settlePlan{
		userRoot:     gitState.Root,
		workDir:      gitState.Root,
		targetBranch: gitState.Branch,
		homeDir:      loadedConfig.HomeDir,
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
		if ok, err := settlePathExists(taskRef.Path); err != nil {
			return settlePlan{}, err
		} else if ok {
			plan.workDir = taskRef.Path
			plan.targetBranch = run.LocalBranch
			plan.run = run
			plan.runRef = runRef
			plan.taskRef = taskRef
			plan.hasRun = true
			plan.taskSurface = true
		} else if ok, err := settlePathExists(run.WorkDir); err != nil {
			return settlePlan{}, err
		} else if ok {
			plan.workDir = run.WorkDir
			plan.targetBranch = run.LocalBranch
			plan.run = run
			plan.runRef = runRef
			plan.hasRun = true
		}
	}
	graph, err := spec.Load(plan.workDir, req.specSlug)
	if err != nil {
		return settlePlan{}, validationError{message: fmt.Sprintf("Spec loads valid: %v", err)}
	}
	task, found := findSettleTask(graph.Tasks, req.taskID)
	if !found {
		return settlePlan{}, validationError{message: fmt.Sprintf("Task %q does not exist in Task Graph for Spec %q; choose a task id from docs/specs/%s/_tasks.md", req.taskID, req.specSlug, req.specSlug)}
	}
	if task.Status != spec.StatusFailed {
		return settlePlan{}, settleStatusError(req.specSlug, task)
	}
	if err := ensureNoSettleActiveRun(ctx, loadedConfig.HomeDir, gitState.Root, req.specSlug); err != nil {
		return settlePlan{}, err
	}
	plan.graph = graph
	plan.task = task
	return plan, nil
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
	switch task.Status {
	case spec.StatusPending, spec.StatusInProgress:
		return validationError{message: fmt.Sprintf("Task %s status is %s; run the Implement Command instead: roundfix implement --spec %s --agent <agent>", task.ID, task.Status, slug)}
	case spec.StatusCompleted:
		return validationError{message: fmt.Sprintf("Task %s status is completed; nothing to do", task.ID)}
	default:
		return validationError{message: fmt.Sprintf("Task %s status is %s; settle requires failed", task.ID, task.Status)}
	}
}

func ensureNoSettleActiveRun(ctx context.Context, homeDir string, gitRoot string, specSlug string) error {
	if _, err := os.Stat(store.DatabasePath(homeDir)); errors.Is(err, os.ErrNotExist) {
		return nil
	} else if err != nil {
		return validationError{message: fmt.Sprintf("inspect Run Database before settle: %v", err)}
	}
	runStore, err := store.OpenReader(ctx, homeDir)
	if err != nil {
		return validationError{message: fmt.Sprintf("open Run Database before settle: %v", err)}
	}
	defer func() {
		_ = runStore.Close()
	}()
	if active, found, err := runStore.ActiveSpecRun(ctx, gitRoot, specSlug); err != nil {
		return validationError{message: fmt.Sprintf("check Active Run for Spec target: %v", err)}
	} else if found {
		return validationError{message: fmt.Sprintf("Active Run %s already holds Spec target %q in working tree %q; stop it with: roundfix stop %s", active.ID, specSlug, gitRoot, active.ID)}
	}
	if active, found, err := runStore.ActiveRunInGitRoot(ctx, gitRoot); err != nil {
		return validationError{message: fmt.Sprintf("check Active Run for working tree: %v", err)}
	} else if found {
		return validationError{message: fmt.Sprintf("Active Run %s already holds working tree %q; stop it with: roundfix stop %s", active.ID, gitRoot, active.ID)}
	}
	return nil
}

func settleTaskAndCommit(ctx context.Context, plan settlePlan, collaborators engineCollaborators) (string, error) {
	taskPath := plan.task.File
	changed, err := collaborators.worktree.Snapshot(ctx, plan.workDir)
	if err != nil {
		return "", err
	}
	changed = ensureSettleCommitPath(changed, taskPath)
	if err := spec.SetStatus(filepathInRoot(plan.workDir, taskPath), spec.StatusCompleted); err != nil {
		return "", fmt.Errorf("settle Task %s completed: %w", plan.task.ID, err)
	}
	message := daemon.TaskCommitMessage(plan.graph.Spec.Slug, plan.task)
	if err := collaborators.committer.Commit(ctx, daemon.CommitRequest{
		WorkDir: plan.workDir,
		Message: message,
		Paths:   changed,
	}); err != nil {
		return "", err
	}
	shortSHA, err := settleShortHEAD(ctx, plan.workDir)
	if err != nil {
		return "", err
	}
	return shortSHA, nil
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
	runStore, err := store.Open(ctx, plan.homeDir)
	if err != nil {
		return "", err
	}
	defer func() {
		_ = runStore.Close()
	}()
	if result.Mode == runworktree.ModePending {
		if _, err := runStore.CompleteRun(ctx, plan.run.ID, store.StateIntegrationPending); err != nil {
			return "", err
		}
		return implementIntegrationCommand(ref), nil
	}
	if err := cleanupCleanRunWorktree(ctx, ref); err != nil {
		return "", err
	}
	if _, err := runStore.CompleteRun(ctx, plan.run.ID, store.StateClean); err != nil {
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
