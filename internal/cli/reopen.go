package cli

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"roundfix/internal/app"
	roundconfig "roundfix/internal/config"
	"roundfix/internal/spec"
	"roundfix/internal/store"
)

const reopenUsage = `Usage:
  roundfix reopen --spec <slug>

Returns a completed terminal QA Task to pending only when one or more of its
dependencies are no longer completed. Records the invalidated QA Report and
stale dependency ids in the QA Task while preserving the prior Result and QA
Report. reopen creates no Run, writes no Run Event Journal entry, and never
commits or pushes.

Options:
  --spec  Spec slug under the configured Spec Root

Exit codes:
  0  stale QA gate reopened
  1  reopen write failed
  2  Preflight Validation failed
`

type reopenPlan struct {
	specsRoot   string
	qaTaskPath  string
	qaTaskID    string
	reportLabel string
	taskIDs     []string
}

func runReopenCommand(ctx context.Context, args []string, stdout, stderr io.Writer, environment commandEnvironment) int {
	if commandWantsHelp(args) {
		fmt.Fprint(stdout, reopenUsage)
		return exitOK
	}
	slug, err := parseReopenCommand(args)
	if err != nil {
		printPreflightFailure("reopen", err, stderr)
		return exitPreflight
	}
	plan, err := preflightReopen(ctx, slug, stderr, environment)
	if err != nil {
		printPreflightFailure("reopen", err, stderr)
		return exitPreflight
	}
	return reopenFromPlan(ctx, slug, plan, stdout, stderr)
}

func reopenFromPlan(ctx context.Context, slug string, plan reopenPlan, stdout, stderr io.Writer) int {
	if err := ctx.Err(); err != nil {
		printPreflightFailure("reopen", err, stderr)
		return exitPreflight
	}
	rechecked, err := deriveReopenPlan(plan.specsRoot, slug)
	if err != nil {
		printPreflightFailure("reopen", fmt.Errorf("gate changed after preflight: %w", err), stderr)
		return exitPreflight
	}
	if err := compareReopenPlans(plan, rechecked); err != nil {
		printPreflightFailure("reopen", err, stderr)
		return exitPreflight
	}
	if err := spec.ReopenGate(rechecked.qaTaskPath, rechecked.reportLabel, rechecked.taskIDs, time.Now().UTC()); err != nil {
		fmt.Fprintf(stderr, "%s: reopen failed: %v\n", app.Name, err)
		return exitRunFailed
	}
	fmt.Fprintf(stdout, "reopened %s pending — invalidated %s\n", rechecked.qaTaskID, rechecked.reportLabel)
	return exitOK
}

func parseReopenCommand(args []string) (string, error) {
	var slug string
	flags := flag.NewFlagSet("reopen", flag.ContinueOnError)
	flags.SetOutput(io.Discard)
	flags.StringVar(&slug, "spec", "", "Spec slug under the configured Spec Root")
	if err := flags.Parse(args); err != nil {
		return "", validationError{message: err.Error()}
	}
	if remaining := flags.Args(); len(remaining) > 0 {
		return "", validationError{message: fmt.Sprintf("unexpected argument %q", remaining[0])}
	}
	slug = strings.TrimSpace(slug)
	if slug == "" {
		return "", validationError{message: "missing required --spec; pass --spec <slug>"}
	}
	if slug == "." || slug == ".." || filepath.IsAbs(slug) || filepath.Base(slug) != slug || strings.ContainsAny(slug, `/\`) {
		return "", validationError{message: fmt.Sprintf("invalid Spec slug %q; --spec must be one Spec directory name", slug)}
	}
	return slug, nil
}

func preflightReopen(ctx context.Context, slug string, stderr io.Writer, environment commandEnvironment) (reopenPlan, error) {
	loaded, err := loadCommandConfig(environment, stderr)
	if err != nil {
		return reopenPlan{}, err
	}
	if loaded.GitRoot == "" {
		return reopenPlan{}, validationError{message: "reopen requires a git repository working tree"}
	}
	if err := ctx.Err(); err != nil {
		return reopenPlan{}, err
	}
	resolvedSpecsRoot, err := roundconfig.ResolveSpecsRoot(loaded, loaded.GitRoot)
	if err != nil {
		return reopenPlan{}, err
	}
	if err := ensureNoReopenActiveRun(ctx, loaded.HomeDir, loaded.GitRoot, slug); err != nil {
		return reopenPlan{}, err
	}
	return deriveReopenPlan(resolvedSpecsRoot.Path, slug)
}

func deriveReopenPlan(specsRoot string, slug string) (reopenPlan, error) {
	graph, loadErr := spec.LoadForRecovery(specsRoot, slug)
	if loadErr == nil {
		return reopenHealthyGateRefusal(graph)
	}
	var stale spec.StaleGateError
	if !errors.As(loadErr, &stale) || graph == nil {
		return reopenPlan{}, loadErr
	}
	qaTask, found := reopenTaskByID(graph.Tasks, stale.QATaskID)
	if !found {
		return reopenPlan{}, validationError{message: fmt.Sprintf("terminal QA Task %q is absent from the loaded graph", stale.QATaskID)}
	}
	reportPath, err := spec.NewestQAReport(graph.Spec.Dir)
	if err != nil {
		return reopenPlan{}, fmt.Errorf("resolve invalidated QA Report: %w", err)
	}
	reportLabel, err := filepath.Rel(graph.Spec.Dir, reportPath)
	if err != nil {
		return reopenPlan{}, fmt.Errorf("format invalidated QA Report path: %w", err)
	}
	qaTaskPath := filepath.Join(specsRoot, qaTask.File)
	specDir := filepath.Join(specsRoot, slug)
	if err := ensureReopenTaskInsideSpecsRoot(specsRoot, specDir, qaTaskPath); err != nil {
		return reopenPlan{}, err
	}
	return reopenPlan{
		specsRoot:   specsRoot,
		qaTaskPath:  qaTaskPath,
		qaTaskID:    qaTask.ID,
		reportLabel: filepath.ToSlash(reportLabel),
		taskIDs:     append([]string(nil), stale.TaskIDs...),
	}, nil
}

func compareReopenPlans(before reopenPlan, after reopenPlan) error {
	if before.qaTaskID != after.qaTaskID {
		return validationError{message: fmt.Sprintf("gate changed after preflight: terminal QA Task changed from %q to %q", before.qaTaskID, after.qaTaskID)}
	}
	if !sameTaskIDSet(before.taskIDs, after.taskIDs) {
		return validationError{message: fmt.Sprintf("gate changed after preflight: stale dependencies for terminal QA Task %q changed from %q to %q", before.qaTaskID, before.taskIDs, after.taskIDs)}
	}
	return nil
}

func sameTaskIDSet(left []string, right []string) bool {
	if len(left) != len(right) {
		return false
	}
	counts := make(map[string]int, len(left))
	for _, taskID := range left {
		counts[taskID]++
	}
	for _, taskID := range right {
		if counts[taskID] == 0 {
			return false
		}
		counts[taskID]--
	}
	return true
}

func ensureNoReopenActiveRun(ctx context.Context, homeDir string, gitRoot string, specSlug string) error {
	if _, err := os.Stat(store.DatabasePath(homeDir)); errors.Is(err, os.ErrNotExist) {
		return nil
	} else if err != nil {
		return validationError{message: fmt.Sprintf("inspect Run Database before reopen: %v", err)}
	}
	runStore, err := store.OpenReader(ctx, homeDir)
	if err != nil {
		return validationError{message: fmt.Sprintf("open Run Database before reopen: %v", err)}
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

func ensureReopenTaskInsideSpecsRoot(specsRoot string, specDir string, taskPath string) error {
	if _, inside := repositoryRelativePath(specDir, taskPath); !inside {
		return validationError{message: fmt.Sprintf("QA Task path %q is outside Spec directory %q", taskPath, specDir)}
	}
	resolvedRoot, err := filepath.EvalSymlinks(specsRoot)
	if err != nil {
		return fmt.Errorf("resolve configured Spec Root %q: %w", specsRoot, err)
	}
	resolvedTask, err := filepath.EvalSymlinks(taskPath)
	if err != nil {
		return fmt.Errorf("resolve QA Task path %q: %w", taskPath, err)
	}
	if _, inside := repositoryRelativePath(resolvedRoot, resolvedTask); !inside {
		return validationError{message: fmt.Sprintf("resolved QA Task path %q is outside configured Spec Root %q", resolvedTask, resolvedRoot)}
	}
	return nil
}

func reopenHealthyGateRefusal(graph *spec.Graph) (reopenPlan, error) {
	if graph == nil || graph.QATaskID == "" {
		return reopenPlan{}, validationError{message: "Spec has no terminal QA Task"}
	}
	qaTask, found := reopenTaskByID(graph.Tasks, graph.QATaskID)
	if !found {
		return reopenPlan{}, validationError{message: fmt.Sprintf("terminal QA Task %q is absent from the loaded graph", graph.QATaskID)}
	}
	if qaTask.Status != spec.StatusCompleted {
		return reopenPlan{}, validationError{message: fmt.Sprintf("terminal QA Task %q is not completed (status: %s)", qaTask.ID, qaTask.Status)}
	}
	return reopenPlan{}, validationError{message: fmt.Sprintf("terminal QA Task %q is not stale; every dependency is completed", qaTask.ID)}
}

func reopenTaskByID(tasks []spec.Task, taskID string) (spec.Task, bool) {
	for _, task := range tasks {
		if task.ID == taskID {
			return task, true
		}
	}
	return spec.Task{}, false
}
