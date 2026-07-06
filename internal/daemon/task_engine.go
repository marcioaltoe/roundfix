package daemon

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"unicode"
	"unicode/utf8"

	"roundfix/internal/agent"
	"roundfix/internal/rounds"
	"roundfix/internal/runevent"
	"roundfix/internal/spec"
	"roundfix/internal/store"
	runworktree "roundfix/internal/worktree"
)

// TaskPlan is the validated input for one Task cycle over an
// already-created implement Run: the full Task Graph in the deterministic
// topological order spec.Load produced. WorkDir is the git root and the
// Agent working directory.
type TaskPlan struct {
	RunID       string
	Session     agent.SessionRef
	WorkDir     string
	RunWorktree runworktree.Ref
	ArtifactDir string
	Spec        spec.Spec
	Tasks       []spec.Task
	Runtime     agent.RuntimeSpec
	QA          bool
	Concurrency int
	CopyList    []string
}

// TaskCycleResult reports what one Task cycle settled. Skipped counts
// Tasks left pending because a needed Task did not end completed.
// QAVerdict stays empty when the QA step did not run; when it ran it is
// the report verdict (pass, fail, partial) or a Daemon settlement
// (missing, unreadable). QAReportPath is the newest QA Report relative to
// the working tree, empty when no report exists.
type TaskCycleResult struct {
	Completed, Failed, Skipped int
	QAVerdict                  string
	QAReportPath               string
}

// QA verdict settlements the Daemon adds beyond the report-authored
// pass/fail/partial values: a QA step that finds no report settles
// "missing" and one whose report verdict cannot be read settles
// "unreadable". Every value except spec.VerdictPass ends the Run
// Unresolved (ADR 0015).
const (
	qaVerdictMissing    = "missing"
	qaVerdictUnreadable = "unreadable"
)

// TaskCycle executes the Task Graph for one Spec as a sibling of
// ResolveCycle: each non-completed Task, in topological order, runs
// through one Agent invocation as a Batch of one, then the Task's
// Verification commands run verbatim, the Daemon settles the task status
// (ADR 0014), and a passing Task gets its own commit (ADR 0013). A failed
// Task creates no commit, its worktree changes are preserved, and the
// cycle continues with independent Tasks (generalizing ADR 0010); only
// Stop Requests and infrastructure errors halt the cycle. The Pusher and
// the Review Source resolver are never invoked for spec Runs.
func (engine *Engine) TaskCycle(ctx context.Context, plan TaskPlan) (TaskCycleResult, error) {
	if err := validateTaskPlan(plan); err != nil {
		return TaskCycleResult{}, err
	}
	concurrency := taskConcurrency(plan)
	if err := engine.publishDaemonEvent(ctx, plan.RunID, 0, runevent.KindDaemonStatus,
		fmt.Sprintf("Task cycle started with concurrency %d.", concurrency),
		map[string]any{"spec": plan.Spec.Slug, "tasks": len(plan.Tasks), "concurrency": concurrency},
	); err != nil {
		return TaskCycleResult{}, err
	}
	statuses := initialTaskRunStatuses(plan.Tasks)
	result, ordinal, err := engine.runTaskScheduler(ctx, plan, statuses)
	if err != nil {
		return result, err
	}
	// QA step (ADR 0015): the qa-gate runs only when plan.QA is set and
	// every Task in the Task Graph — including Tasks completed by earlier
	// Runs — ended completed; any failed, skipped, or pending Task leaves
	// the outcome to the Task results alone.
	if plan.QA && allTasksRunCompleted(plan.Tasks, statuses) {
		if err := engine.stopIfRequested(ctx, plan.RunID, ordinal+1); err != nil {
			return result, fmt.Errorf("stop run %q before the QA step: %w", plan.RunID, err)
		}
		ordinal++
		verdict, reportPath, err := engine.runQAGate(ctx, plan, ordinal)
		if err != nil {
			return result, err
		}
		result.QAVerdict = verdict
		result.QAReportPath = reportPath
	}
	if err := engine.publishDaemonEvent(ctx, plan.RunID, 0, runevent.KindDaemonOutcome,
		fmt.Sprintf("Task cycle finished: %d completed, %d failed, %d skipped.", result.Completed, result.Failed, result.Skipped),
		map[string]any{"completed": result.Completed, "failed": result.Failed, "skipped": result.Skipped},
	); err != nil {
		return result, err
	}
	return result, nil
}

type taskRunStatus string

const (
	taskRunPending   taskRunStatus = "pending"
	taskRunRunning   taskRunStatus = "running"
	taskRunCompleted taskRunStatus = "completed"
	taskRunFailed    taskRunStatus = "failed"
	taskRunSkipped   taskRunStatus = "skipped"
)

type taskWorkerResult struct {
	task             spec.Task
	ordinal          int
	status           spec.Status
	taskPlan         TaskPlan
	taskRef          runworktree.TaskRef
	usesTaskWorktree bool
	err              error
}

func (engine *Engine) runTaskScheduler(ctx context.Context, plan TaskPlan, statuses map[string]taskRunStatus) (TaskCycleResult, int, error) {
	result := TaskCycleResult{}
	concurrency := taskConcurrency(plan)
	useTaskWorktrees := concurrency > 1 && hasParallelTaskPotential(plan.Tasks, statuses)
	results := make(chan taskWorkerResult)
	running := 0
	ordinal := 0
	var stopErr error
	var fatalErr error

	for {
		if fatalErr == nil && stopErr == nil {
			skipped, err := engine.skipBlockedTasks(ctx, plan, statuses, &result)
			if err != nil {
				return result, ordinal, err
			}
			_ = skipped
			for running < concurrency {
				task, ok := nextReadyTask(plan.Tasks, statuses)
				if !ok {
					break
				}
				if err := taskStartBoundary(ctx, engine, plan.RunID, ordinal, task.ID); err != nil {
					stopErr = err
					break
				}
				ordinal++
				statuses[task.ID] = taskRunRunning
				running++
				go func(task spec.Task, ordinal int) {
					results <- engine.executeTaskWorker(ctx, plan, task, ordinal, useTaskWorktrees)
				}(task, ordinal)
			}
		}

		if running == 0 {
			switch {
			case fatalErr != nil:
				return result, ordinal, fatalErr
			case stopErr != nil:
				return result, ordinal, stopErr
			case allTaskRunStatusesTerminal(statuses):
				return result, ordinal, nil
			default:
				// A valid graph should not get here. If it does, keep the
				// shipped behavior by reporting remaining Tasks skipped.
				if err := engine.skipPendingTasks(ctx, plan, statuses, &result); err != nil {
					return result, ordinal, err
				}
			}
			continue
		}

		workerResult := <-results
		running--
		if workerResult.err != nil {
			statuses[workerResult.task.ID] = taskRunFailed
			if fatalErr == nil {
				fatalErr = workerResult.err
			}
			continue
		}
		if fatalErr != nil {
			continue
		}
		settled, err := engine.integrateTaskSettlement(ctx, plan, workerResult)
		if err != nil {
			if fatalErr == nil {
				fatalErr = err
			}
			continue
		}
		if settled == spec.StatusCompleted {
			statuses[workerResult.task.ID] = taskRunCompleted
			result.Completed++
		} else {
			statuses[workerResult.task.ID] = taskRunFailed
			result.Failed++
		}
		if stopErr == nil {
			if err := engine.stopIfRequested(ctx, plan.RunID, workerResult.ordinal); err != nil {
				stopErr = fmt.Errorf("stop run %q after Task %s settlement: %w", plan.RunID, workerResult.task.ID, err)
			}
		}
	}
}

func (engine *Engine) executeTaskWorker(ctx context.Context, plan TaskPlan, task spec.Task, ordinal int, usesTaskWorktree bool) taskWorkerResult {
	taskPlan := plan
	var taskRef runworktree.TaskRef
	if usesTaskWorktree {
		var err error
		taskRef, err = engine.deps.TaskWorktrees.CreateTask(ctx, plan.RunWorktree, task.ID, plan.CopyList)
		if err != nil {
			return taskWorkerResult{task: task, ordinal: ordinal, usesTaskWorktree: true, err: err}
		}
		taskPlan = taskPlanForTaskWorktree(plan, task, taskRef)
	}
	settled, err := engine.executeTask(ctx, taskPlan, task, ordinal)
	if usesTaskWorktree {
		if closeErr := engine.deps.Runner.EndSession(context.WithoutCancel(ctx), taskPlan.Runtime, taskPlan.Session); closeErr != nil {
			err = errors.Join(err, fmt.Errorf("close Agent Session for run %q Task %s: %w", plan.RunID, task.ID, closeErr))
		}
	}
	return taskWorkerResult{
		task:             task,
		ordinal:          ordinal,
		status:           settled,
		taskPlan:         taskPlan,
		taskRef:          taskRef,
		usesTaskWorktree: usesTaskWorktree,
		err:              err,
	}
}

func (engine *Engine) integrateTaskSettlement(ctx context.Context, plan TaskPlan, result taskWorkerResult) (spec.Status, error) {
	if result.status != spec.StatusCompleted || !result.usesTaskWorktree {
		return result.status, nil
	}
	integration, err := engine.deps.TaskWorktrees.IntegrateTask(ctx, plan.RunWorktree, result.taskRef)
	if err != nil {
		return "", fmt.Errorf("integrate Task %s into Run Branch: %w", result.task.ID, err)
	}
	if integration.Mode == runworktree.ModeTaskConflict {
		reason := strings.TrimSpace(integration.Reason)
		if reason == "" {
			reason = "unknown conflict"
		}
		reason = "integration conflict: " + reason
		if err := engine.settleTask(ctx, result.taskPlan, result.task, result.ordinal, spec.StatusFailed, reason); err != nil {
			return "", err
		}
		fmt.Fprintf(engine.deps.Progress, "Task %s failed: %s\n", result.task.ID, reason)
		return spec.StatusFailed, nil
	}
	if err := engine.deps.TaskWorktrees.CleanupTask(ctx, result.taskRef); err != nil {
		return "", fmt.Errorf("cleanup Task Worktree for %s: %w", result.task.ID, err)
	}
	return spec.StatusCompleted, nil
}

func initialTaskRunStatuses(tasks []spec.Task) map[string]taskRunStatus {
	statuses := make(map[string]taskRunStatus, len(tasks))
	for _, task := range tasks {
		if task.Status == spec.StatusCompleted {
			statuses[task.ID] = taskRunCompleted
			continue
		}
		statuses[task.ID] = taskRunPending
	}
	return statuses
}

func taskConcurrency(plan TaskPlan) int {
	if plan.Concurrency < 1 {
		return 1
	}
	return plan.Concurrency
}

func hasParallelTaskPotential(tasks []spec.Task, statuses map[string]taskRunStatus) bool {
	remaining := make([]spec.Task, 0, len(tasks))
	for _, task := range tasks {
		if statuses[task.ID] != taskRunCompleted {
			remaining = append(remaining, task)
		}
	}
	for left := 0; left < len(remaining); left++ {
		for right := left + 1; right < len(remaining); right++ {
			if !taskDependsOn(remaining[left].ID, remaining[right].ID, tasks) && !taskDependsOn(remaining[right].ID, remaining[left].ID, tasks) {
				return true
			}
		}
	}
	return false
}

func taskDependsOn(taskID string, needID string, tasks []spec.Task) bool {
	byID := make(map[string]spec.Task, len(tasks))
	for _, task := range tasks {
		byID[task.ID] = task
	}
	seen := map[string]bool{}
	var visit func(string) bool
	visit = func(id string) bool {
		if seen[id] {
			return false
		}
		seen[id] = true
		task, found := byID[id]
		if !found {
			return false
		}
		for _, need := range task.Needs {
			if need == needID || visit(need) {
				return true
			}
		}
		return false
	}
	return visit(taskID)
}

func taskStartBoundary(ctx context.Context, engine *Engine, runID string, ordinal int, taskID string) error {
	if err := ctx.Err(); err != nil {
		if publishErr := engine.publishStop(ctx, runID, ordinal); publishErr != nil {
			return fmt.Errorf("publish stop event for run %q before Task %s: %w", runID, taskID, errors.Join(err, publishErr))
		}
		return fmt.Errorf("stop run %q before Task %s: %w", runID, taskID, err)
	}
	if err := engine.stopIfRequested(ctx, runID, ordinal); err != nil {
		return fmt.Errorf("stop run %q before Task %s: %w", runID, taskID, err)
	}
	return nil
}

func nextReadyTask(tasks []spec.Task, statuses map[string]taskRunStatus) (spec.Task, bool) {
	for _, task := range tasks {
		if statuses[task.ID] != taskRunPending {
			continue
		}
		if taskNeedsCompleted(task, statuses) {
			return task, true
		}
	}
	return spec.Task{}, false
}

func taskNeedsCompleted(task spec.Task, statuses map[string]taskRunStatus) bool {
	for _, need := range task.Needs {
		if statuses[need] != taskRunCompleted {
			return false
		}
	}
	return true
}

func (engine *Engine) skipBlockedTasks(ctx context.Context, plan TaskPlan, statuses map[string]taskRunStatus, result *TaskCycleResult) (bool, error) {
	skippedAny := false
	for {
		skippedPass := false
		for _, task := range plan.Tasks {
			if statuses[task.ID] != taskRunPending {
				continue
			}
			if !hasTerminalFailedNeed(task, statuses) {
				continue
			}
			if err := engine.skipTask(ctx, plan.RunID, task, statuses); err != nil {
				return skippedAny, err
			}
			result.Skipped++
			skippedPass = true
			skippedAny = true
		}
		if !skippedPass {
			return skippedAny, nil
		}
	}
}

func (engine *Engine) skipPendingTasks(ctx context.Context, plan TaskPlan, statuses map[string]taskRunStatus, result *TaskCycleResult) error {
	for _, task := range plan.Tasks {
		if statuses[task.ID] != taskRunPending {
			continue
		}
		if err := engine.skipTask(ctx, plan.RunID, task, statuses); err != nil {
			return err
		}
		result.Skipped++
	}
	return nil
}

func (engine *Engine) skipTask(ctx context.Context, runID string, task spec.Task, statuses map[string]taskRunStatus) error {
	unmet := unmetRunNeeds(task, statuses)
	reason := fmt.Sprintf("needs not completed: %s", strings.Join(unmet, ", "))
	statuses[task.ID] = taskRunSkipped
	fmt.Fprintf(engine.deps.Progress, "Task %s skipped: %s\n", task.ID, reason)
	return engine.publishTaskEvent(ctx, runID, 0, task.ID, runevent.KindDaemonTask,
		fmt.Sprintf("Task %s skipped: %s", task.ID, reason),
		map[string]any{"task": task.ID, "phase": "skipped", "reason": reason},
	)
}

func hasTerminalFailedNeed(task spec.Task, statuses map[string]taskRunStatus) bool {
	for _, need := range task.Needs {
		if statuses[need] == taskRunFailed || statuses[need] == taskRunSkipped {
			return true
		}
	}
	return false
}

func unmetRunNeeds(task spec.Task, statuses map[string]taskRunStatus) []string {
	var unmet []string
	for _, need := range task.Needs {
		if statuses[need] != taskRunCompleted {
			unmet = append(unmet, need)
		}
	}
	return unmet
}

func allTaskRunStatusesTerminal(statuses map[string]taskRunStatus) bool {
	for _, status := range statuses {
		switch status {
		case taskRunCompleted, taskRunFailed, taskRunSkipped:
		default:
			return false
		}
	}
	return true
}

func taskPlanForTaskWorktree(plan TaskPlan, task spec.Task, taskRef runworktree.TaskRef) TaskPlan {
	taskPlan := plan
	taskPlan.WorkDir = taskRef.Path
	taskPlan.Session = agent.SessionRefForTask(plan.RunID, task.ID, taskRef.Path)
	taskPlan.Spec = specForWorkDir(plan, taskRef.Path)
	return taskPlan
}

func specForWorkDir(plan TaskPlan, workDir string) spec.Spec {
	specRef := plan.Spec
	if rel, err := filepath.Rel(plan.WorkDir, plan.Spec.Dir); err == nil && rel != "." && !strings.HasPrefix(rel, "..") {
		specRef.Dir = filepath.Join(workDir, rel)
	}
	return specRef
}

// executeTask runs one Task end to end: before-snapshot, Agent, reload,
// Verification, settlement, and the Task commit on success. It returns
// the settled status; the returned error is reserved for Stop Requests
// and infrastructure failures, which halt the cycle.
func (engine *Engine) executeTask(ctx context.Context, plan TaskPlan, task spec.Task, ordinal int) (spec.Status, error) {
	// The before-snapshot is taken before the Agent starts, so anything
	// already dirty — pre-existing user work or a failed Task's preserved
	// changes — never reaches this Task's commit.
	before, err := engine.deps.Worktree.Snapshot(ctx, plan.WorkDir)
	if err != nil {
		return "", err
	}
	failure, err := engine.runTaskAgent(ctx, plan, &task, ordinal)
	if err != nil {
		return "", err
	}
	if failure == "" && task.Status == spec.StatusFailed {
		failure = "Agent settled the Task failed"
	}
	if failure == "" {
		failure, err = engine.verifyTask(ctx, plan, task, ordinal)
		if err != nil {
			return "", err
		}
	}
	settled := spec.StatusCompleted
	if failure != "" {
		settled = spec.StatusFailed
	}
	if err := engine.settleTask(ctx, plan, task, ordinal, settled, failure); err != nil {
		return "", err
	}
	if failure != "" {
		// No commit for a failed Task; its worktree changes are preserved
		// and the cycle continues with independent Tasks.
		fmt.Fprintf(engine.deps.Progress, "Task %s failed: %s\n", task.ID, failure)
		return settled, nil
	}
	if err := engine.commitTask(ctx, plan, task, ordinal, before); err != nil {
		return "", err
	}
	fmt.Fprintf(engine.deps.Progress, "Task %s completed.\n", task.ID)
	return settled, nil
}

// runTaskAgent runs the Agent over one Task as a Batch of one, reading the
// task file fresh for the prompt and reloading it afterwards. It returns a
// non-empty failure reason when the Task failed but the cycle should
// continue; the returned error is reserved for Stop Requests and
// infrastructure failures, which halt the cycle.
func (engine *Engine) runTaskAgent(ctx context.Context, plan TaskPlan, task *spec.Task, ordinal int) (string, error) {
	if err := engine.deps.Runs.UpdateRunState(ctx, plan.RunID, store.StateResolvingWithAgent); err != nil {
		return "", fmt.Errorf("update run %q to state %q before Task %s: %w", plan.RunID, store.StateResolvingWithAgent, task.ID, err)
	}
	if err := engine.publishTaskEvent(ctx, plan.RunID, ordinal, task.ID, runevent.KindDaemonTask,
		fmt.Sprintf("Task %s started as Batch %03d: %s", task.ID, ordinal, task.Title),
		map[string]any{"task": task.ID, "phase": "started", "batch": ordinal},
	); err != nil {
		return "", fmt.Errorf("publish start event for run %q Task %s: %w", plan.RunID, task.ID, err)
	}
	taskPath := filepath.Join(plan.WorkDir, task.File)
	content, err := os.ReadFile(taskPath)
	if err != nil {
		return "", fmt.Errorf("read Task %q file %q before the Agent: %w", task.ID, taskPath, err)
	}
	prompt, err := agent.BuildTaskPrompt(agent.TaskPromptRequest{
		SpecSlug:    plan.Spec.Slug,
		TaskID:      task.ID,
		TaskPath:    task.File,
		TaskContent: string(content),
	})
	if err != nil {
		return "", fmt.Errorf("build Task prompt for run %q Task %s: %w", plan.RunID, task.ID, err)
	}
	logPath := agent.LogPath(plan.ArtifactDir, plan.RunID, ordinal)
	fmt.Fprintf(engine.deps.Progress, "Task: %s (Batch %03d) %s\n", task.ID, ordinal, task.Title)
	fmt.Fprintf(engine.deps.Progress, "Agent log: %s\n", logPath)

	runResult, runErr := engine.deps.Runner.Run(ctx, agent.ExecuteRequest{
		Runtime:     plan.Runtime,
		Session:     plan.Session,
		RunID:       plan.RunID,
		Batch:       rounds.Batch{Number: ordinal},
		LogPath:     logPath,
		ArtifactDir: plan.ArtifactDir,
		Prompt:      prompt,
		GitRoot:     plan.WorkDir,
	}, engine.deps.Sink)
	if runErr != nil {
		if isStop(ctx, runErr) {
			// The runner already published the stopped status event;
			// Agent-created worktree changes stay untouched.
			return "", fmt.Errorf("run Agent for run %q Task %s: %w", plan.RunID, task.ID, runErr)
		}
		return fmt.Sprintf("Agent failed: %v", runErr), nil
	}
	if err := ctx.Err(); err != nil {
		if publishErr := engine.publishStop(ctx, plan.RunID, ordinal); publishErr != nil {
			return "", fmt.Errorf("publish stop event for run %q after Agent Task %s: %w", plan.RunID, task.ID, errors.Join(err, publishErr))
		}
		return "", fmt.Errorf("stop run %q after Agent Task %s: %w", plan.RunID, task.ID, err)
	}
	if anomaly := strings.TrimSpace(runResult.TransportAnomaly); anomaly != "" {
		if err := engine.publishTaskEvent(ctx, plan.RunID, ordinal, task.ID, runevent.KindDaemonTask,
			fmt.Sprintf("Task %s transport anomaly: %s", task.ID, anomaly),
			map[string]any{"task": task.ID, "phase": "transport_anomaly", "batch": ordinal, "anomaly": anomaly},
		); err != nil {
			return "", fmt.Errorf("publish transport anomaly event for run %q Task %s: %w", plan.RunID, task.ID, err)
		}
	}
	if err := spec.ReloadTask(plan.WorkDir, task); err != nil {
		// The Agent left the task file unreadable; the Task fails and the
		// Daemon settles the status by rewriting the frontmatter value.
		return fmt.Sprintf("reload task file after the Agent: %v", err), nil
	}
	return "", nil
}

// verifyTask runs every Verification command of the Task sequentially and
// verbatim through the Verifier, in WorkDir; the first failing command
// fails the Task. defaults.verification is never appended: the Daemon gate
// runs only the Task's own Verification commands (ADR 0014). It returns a
// non-empty failure reason when verification failed; the returned error is
// reserved for Stop Requests and infrastructure failures.
func (engine *Engine) verifyTask(ctx context.Context, plan TaskPlan, task spec.Task, ordinal int) (string, error) {
	if err := engine.deps.Runs.UpdateRunState(ctx, plan.RunID, store.StateVerifying); err != nil {
		return "", fmt.Errorf("update run %q to state %q before Task %s verification: %w", plan.RunID, store.StateVerifying, task.ID, err)
	}
	for _, command := range task.Verification {
		if err := engine.publishTaskEvent(ctx, plan.RunID, ordinal, task.ID, runevent.KindDaemonVerification,
			fmt.Sprintf("Verification started: %s", command),
			map[string]any{"phase": "started", "command": command, "task": task.ID},
		); err != nil {
			return "", fmt.Errorf("publish verification start event for run %q Task %s: %w", plan.RunID, task.ID, err)
		}
		if err := engine.deps.Verifier.Verify(ctx, VerifyRequest{
			WorkDir: plan.WorkDir,
			Command: command,
			Stream:  engine.deps.Progress,
		}); err != nil {
			if isStop(ctx, err) {
				// A Stop Request during verification keeps the Agent's task
				// status untouched; the run ends Stopped, not failed.
				return "", fmt.Errorf("verify run %q Task %s: %w", plan.RunID, task.ID, err)
			}
			if publishErr := engine.publishTaskEvent(ctx, plan.RunID, ordinal, task.ID, runevent.KindDaemonVerification,
				fmt.Sprintf("Verification failed: %s", command),
				map[string]any{"phase": "failed", "command": command, "task": task.ID, "error": err.Error()},
			); publishErr != nil {
				return "", fmt.Errorf("publish verification failure event for run %q Task %s: %w", plan.RunID, task.ID, publishErr)
			}
			return fmt.Sprintf("verification failed: %v", err), nil
		}
		if err := engine.publishTaskEvent(ctx, plan.RunID, ordinal, task.ID, runevent.KindDaemonVerification,
			fmt.Sprintf("Verification command passed: %s", command),
			map[string]any{"phase": "passed", "command": command, "task": task.ID},
		); err != nil {
			return "", fmt.Errorf("publish verification pass event for run %q Task %s: %w", plan.RunID, task.ID, err)
		}
		fmt.Fprintf(engine.deps.Progress, "Verification command passed: %s\n", command)
	}
	return "", nil
}

// settleTask writes the Daemon-owned final status when the Agent left
// anything else and journals the settlement (ADR 0014). Callers pass
// completed only after every Verification command passed, so completed is
// never settled without passing verification.
func (engine *Engine) settleTask(ctx context.Context, plan TaskPlan, task spec.Task, ordinal int, status spec.Status, reason string) error {
	if task.Status != status {
		taskPath := filepath.Join(plan.WorkDir, task.File)
		if err := spec.SetStatus(taskPath, status); err != nil {
			return fmt.Errorf("settle Task %s status %q for run %q: %w", task.ID, status, plan.RunID, err)
		}
	}
	payload := map[string]any{"task": task.ID, "phase": "settled", "status": string(status)}
	summary := fmt.Sprintf("Task %s settled %s.", task.ID, status)
	if reason != "" {
		payload["reason"] = reason
		summary = fmt.Sprintf("Task %s settled %s: %s", task.ID, status, reason)
	}
	if err := engine.publishTaskEvent(ctx, plan.RunID, ordinal, task.ID, runevent.KindDaemonTask, summary, payload); err != nil {
		return fmt.Errorf("publish settlement event for run %q Task %s: %w", plan.RunID, task.ID, err)
	}
	return nil
}

// commitTask creates the Task commit from the snapshot diff plus the task
// file itself, so the settled status and Result section always ride in the
// same commit as the code changes (ADR 0013).
func (engine *Engine) commitTask(ctx context.Context, plan TaskPlan, task spec.Task, ordinal int, before []string) error {
	if err := ctx.Err(); err != nil {
		if publishErr := engine.publishStop(ctx, plan.RunID, ordinal); publishErr != nil {
			return fmt.Errorf("publish stop event for run %q before Task %s commit: %w", plan.RunID, task.ID, errors.Join(err, publishErr))
		}
		return fmt.Errorf("stop run %q before Task %s commit: %w", plan.RunID, task.ID, err)
	}
	after, err := engine.deps.Worktree.Snapshot(ctx, plan.WorkDir)
	if err != nil {
		return err
	}
	changed := ensureCommitPath(diffSnapshots(before, after), task.File)
	message := TaskCommitMessage(plan.Spec.Slug, task)
	if err := engine.deps.Committer.Commit(ctx, CommitRequest{
		WorkDir: plan.WorkDir,
		Message: message,
		Paths:   changed,
	}); err != nil {
		return err
	}
	subject, _, _ := strings.Cut(message, "\n")
	fmt.Fprintf(engine.deps.Progress, "Task commit created: %s\n", subject)
	if err := engine.publishTaskEvent(ctx, plan.RunID, ordinal, task.ID, runevent.KindDaemonCommit,
		fmt.Sprintf("Task commit created: %s", subject),
		map[string]any{"decision": "created", "task": task.ID, "paths": len(changed)},
	); err != nil {
		return fmt.Errorf("publish commit event for run %q Task %s: %w", plan.RunID, task.ID, err)
	}
	return nil
}

// TaskCommitMessage derives the Task commit message from the task
// frontmatter (ADR 0013): `<type>: <task title>` with the title's first
// rune lowercased, where docs, test, and chore pass through and every
// other type maps to feat, plus the
// Roundfix-Spec and Roundfix-Task trailers separated from the subject by a
// blank line per git trailer convention.
func TaskCommitMessage(slug string, task spec.Task) string {
	return fmt.Sprintf("%s: %s\n\nRoundfix-Spec: %s\nRoundfix-Task: %s", taskCommitType(task.Type), lowerFirstRune(task.Title), slug, task.ID)
}

func taskCommitType(taskType string) string {
	switch taskType {
	case "docs", "test", "chore":
		return taskType
	default:
		return "feat"
	}
}

func lowerFirstRune(value string) string {
	first, size := utf8.DecodeRuneInString(value)
	if first == utf8.RuneError && size == 0 {
		return value
	}
	lowered := unicode.ToLower(first)
	if lowered == first {
		return value
	}
	return string(lowered) + value[size:]
}

// runQAGate runs the qa-gate step as the Run's last Batch: before-snapshot,
// Agent with the QA prompt, verdict settling from the newest QA Report
// frontmatter, the daemon.qa Run Event, and the QA Report commit — created
// either way, pass or fail, while a missing report commits nothing
// (ADR 0015). It returns the settled verdict and the report path relative
// to the working tree; the returned error is reserved for Stop Requests
// and infrastructure failures. Unlike a Task, the QA step has no per-item
// settlement to fall back on, so a non-stop Agent failure halts the cycle.
func (engine *Engine) runQAGate(ctx context.Context, plan TaskPlan, ordinal int) (string, string, error) {
	if err := ctx.Err(); err != nil {
		if publishErr := engine.publishStop(ctx, plan.RunID, ordinal); publishErr != nil {
			return "", "", fmt.Errorf("publish stop event for run %q before the QA step: %w", plan.RunID, errors.Join(err, publishErr))
		}
		return "", "", fmt.Errorf("stop run %q before the QA step: %w", plan.RunID, err)
	}
	// The before-snapshot keeps everything already dirty out of the QA
	// Report commit: only the QA step's own report and evidence ride in it.
	before, err := engine.deps.Worktree.Snapshot(ctx, plan.WorkDir)
	if err != nil {
		return "", "", err
	}
	if err := engine.deps.Runs.UpdateRunState(ctx, plan.RunID, store.StateResolvingWithAgent); err != nil {
		return "", "", fmt.Errorf("update run %q to state %q before the QA step: %w", plan.RunID, store.StateResolvingWithAgent, err)
	}
	prompt, err := agent.BuildQAPrompt(plan.Spec.Slug, plan.Spec.Dir, filepath.Join(plan.Spec.Dir, "_prd.md"))
	if err != nil {
		return "", "", fmt.Errorf("build QA prompt for run %q: %w", plan.RunID, err)
	}
	logPath := agent.LogPath(plan.ArtifactDir, plan.RunID, ordinal)
	fmt.Fprintf(engine.deps.Progress, "QA step (Batch %03d) for Spec %s\n", ordinal, plan.Spec.Slug)
	fmt.Fprintf(engine.deps.Progress, "Agent log: %s\n", logPath)
	if _, runErr := engine.deps.Runner.Run(ctx, agent.ExecuteRequest{
		Runtime:     plan.Runtime,
		Session:     plan.Session,
		RunID:       plan.RunID,
		Batch:       rounds.Batch{Number: ordinal},
		LogPath:     logPath,
		ArtifactDir: plan.ArtifactDir,
		Prompt:      prompt,
		GitRoot:     plan.WorkDir,
	}, engine.deps.Sink); runErr != nil {
		return "", "", fmt.Errorf("run Agent for run %q QA step: %w", plan.RunID, runErr)
	}
	if err := ctx.Err(); err != nil {
		if publishErr := engine.publishStop(ctx, plan.RunID, ordinal); publishErr != nil {
			return "", "", fmt.Errorf("publish stop event for run %q after the QA step Agent: %w", plan.RunID, errors.Join(err, publishErr))
		}
		return "", "", fmt.Errorf("stop run %q after the QA step Agent: %w", plan.RunID, err)
	}
	verdict, reportPath := engine.settleQAVerdict(plan)
	if err := engine.publishDaemonEvent(ctx, plan.RunID, ordinal, runevent.KindDaemonQA,
		fmt.Sprintf("QA verdict %s for Spec %s.", verdict, plan.Spec.Slug),
		map[string]any{"verdict": verdict, "report": reportPath},
	); err != nil {
		return "", "", fmt.Errorf("publish QA event for run %q: %w", plan.RunID, err)
	}
	fmt.Fprintf(engine.deps.Progress, "QA verdict: %s\n", verdict)
	if reportPath == "" {
		// A missing report means nothing to commit (ADR 0015).
		return verdict, reportPath, nil
	}
	if err := engine.commitQAReport(ctx, plan, ordinal, before, verdict, reportPath); err != nil {
		return "", "", err
	}
	return verdict, reportPath, nil
}

// settleQAVerdict settles the QA verdict from the newest QA Report: the
// report's own verdict when readable, missing when no report exists, and
// unreadable when the report exists but its verdict cannot be read
// (ADR 0015). The report path comes back relative to the working tree,
// empty when no report exists.
func (engine *Engine) settleQAVerdict(plan TaskPlan) (string, string) {
	verdict := ""
	switch value, err := spec.QAVerdict(plan.Spec.Dir); {
	case err == nil:
		verdict = value
	case errors.Is(err, spec.ErrNoQAReport):
		verdict = qaVerdictMissing
	default:
		fmt.Fprintf(engine.deps.Progress, "QA Report verdict unreadable: %v\n", err)
		verdict = qaVerdictUnreadable
	}
	reportPath := ""
	if newest, err := spec.NewestQAReport(plan.Spec.Dir); err == nil {
		reportPath = newest
		if relative, relErr := filepath.Rel(plan.WorkDir, newest); relErr == nil {
			reportPath = relative
		}
	}
	return verdict, reportPath
}

// commitQAReport creates the QA Report commit from the QA step's snapshot
// diff with the report ensured, so the report and its evidence always ride
// in their own commit, separate from every Task commit (ADR 0015).
func (engine *Engine) commitQAReport(ctx context.Context, plan TaskPlan, ordinal int, before []string, verdict string, reportPath string) error {
	if err := ctx.Err(); err != nil {
		if publishErr := engine.publishStop(ctx, plan.RunID, ordinal); publishErr != nil {
			return fmt.Errorf("publish stop event for run %q before the QA Report commit: %w", plan.RunID, errors.Join(err, publishErr))
		}
		return fmt.Errorf("stop run %q before the QA Report commit: %w", plan.RunID, err)
	}
	after, err := engine.deps.Worktree.Snapshot(ctx, plan.WorkDir)
	if err != nil {
		return err
	}
	changed := ensureCommitPath(diffSnapshots(before, after), reportPath)
	message := QACommitMessage(plan.Spec.Slug, verdict)
	if err := engine.deps.Committer.Commit(ctx, CommitRequest{
		WorkDir: plan.WorkDir,
		Message: message,
		Paths:   changed,
	}); err != nil {
		return err
	}
	subject, _, _ := strings.Cut(message, "\n")
	fmt.Fprintf(engine.deps.Progress, "QA Report commit created: %s\n", subject)
	if err := engine.publishDaemonEvent(ctx, plan.RunID, ordinal, runevent.KindDaemonCommit,
		fmt.Sprintf("QA Report commit created: %s", subject),
		map[string]any{"decision": "created", "report": reportPath, "paths": len(changed)},
	); err != nil {
		return fmt.Errorf("publish QA commit event for run %q: %w", plan.RunID, err)
	}
	return nil
}

// QACommitMessage is the QA Report commit contract (ADR 0013, ADR 0015):
// `docs: qa report for <slug> (<verdict>)` plus the Roundfix-Spec trailer
// separated from the subject by a blank line per git trailer convention.
func QACommitMessage(slug string, verdict string) string {
	return fmt.Sprintf("docs: qa report for %s (%s)\n\nRoundfix-Spec: %s", slug, verdict, slug)
}

// ensureCommitPath guarantees the settled task file or QA Report rides in
// its own commit even when it was already dirty before the step started
// and the snapshot diff therefore excluded it.
func ensureCommitPath(changed []string, path string) []string {
	for _, existing := range changed {
		if existing == path {
			return changed
		}
	}
	changed = append(changed, path)
	sort.Strings(changed)
	return changed
}

func allTasksCompleted(tasks []spec.Task, statuses map[string]spec.Status) bool {
	for _, task := range tasks {
		if statuses[task.ID] != spec.StatusCompleted {
			return false
		}
	}
	return true
}

func allTasksRunCompleted(tasks []spec.Task, statuses map[string]taskRunStatus) bool {
	for _, task := range tasks {
		if statuses[task.ID] != taskRunCompleted {
			return false
		}
	}
	return true
}

func unmetNeeds(task spec.Task, statuses map[string]spec.Status) []string {
	var unmet []string
	for _, need := range task.Needs {
		if statuses[need] != spec.StatusCompleted {
			unmet = append(unmet, need)
		}
	}
	return unmet
}

// publishTaskEvent appends one daemon-owned Run Event for a Task, carrying
// the Task id in the existing per-issue Work Item field with no journal
// schema change. The Batch number is the Task's 1-based execution ordinal,
// or 0 for Tasks that never executed.
func (engine *Engine) publishTaskEvent(ctx context.Context, runID string, ordinal int, taskID string, kind runevent.Kind, summary string, payload any) error {
	raw, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("encode daemon event payload: %w", err)
	}
	if err := engine.deps.Sink.Publish(ctx, runevent.RunEvent{
		RunID:       runID,
		Batch:       ordinal,
		Source:      runevent.SourceDaemon,
		Kind:        kind,
		ReviewIssue: taskID,
		Summary:     runevent.BoundSummary(summary),
		Time:        engine.deps.Now(),
		Payload:     raw,
	}); err != nil {
		return fmt.Errorf("publish daemon event %s: %w", kind, err)
	}
	return nil
}

func validateTaskPlan(plan TaskPlan) error {
	required := map[string]string{
		"Run ID":             plan.RunID,
		"Agent Session":      plan.Session.Name,
		"working tree":       plan.WorkDir,
		"Artifact Directory": plan.ArtifactDir,
		"Spec slug":          plan.Spec.Slug,
	}
	for label, value := range required {
		if strings.TrimSpace(value) == "" {
			return fmt.Errorf("task cycle: %s is required", label)
		}
	}
	if len(plan.Tasks) == 0 {
		return errors.New("task cycle: at least one Task is required")
	}
	if taskConcurrency(plan) > 1 {
		concurrentRequired := map[string]string{
			"Run Worktree path":   plan.RunWorktree.Path,
			"Run Branch":          plan.RunWorktree.Branch,
			"user checkout root":  plan.RunWorktree.UserRoot,
			"Run Worktree Run ID": plan.RunWorktree.RunID,
		}
		for label, value := range concurrentRequired {
			if strings.TrimSpace(value) == "" {
				return fmt.Errorf("task cycle: %s is required when concurrency is greater than 1", label)
			}
		}
	}
	return nil
}
