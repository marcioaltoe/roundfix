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

	"roundfix/internal/agent"
	"roundfix/internal/rounds"
	"roundfix/internal/runevent"
	"roundfix/internal/spec"
	"roundfix/internal/store"
)

// specRunArtifactDirName is the documented Artifact Directory default.
// TaskPlan carries no Artifact Directory, so spec Runs keep Agent logs
// under the default directory inside the working tree, following the
// review path's runs/<run-id>/agent/batch-NNN.log convention.
const specRunArtifactDirName = ".roundfix"

// TaskPlan is the validated input for one Task cycle over an
// already-created implement Run: the full Task Graph in the deterministic
// topological order spec.Load produced. WorkDir is the git root and the
// Agent working directory.
type TaskPlan struct {
	RunID   string
	WorkDir string
	Spec    spec.Spec
	Tasks   []spec.Task
	Runtime agent.RuntimeSpec
	QA      bool
}

// TaskCycleResult reports what one Task cycle settled. Skipped counts
// Tasks left pending because a needed Task did not end completed.
// QAVerdict stays empty when the QA step did not run.
type TaskCycleResult struct {
	Completed, Failed, Skipped int
	QAVerdict                  string
}

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
	// The live status map is seeded from the loaded task files, so Tasks
	// completed in prior Runs satisfy needs, and updated as this Run
	// settles Tasks.
	statuses := make(map[string]spec.Status, len(plan.Tasks))
	for _, task := range plan.Tasks {
		statuses[task.ID] = task.Status
	}
	result := TaskCycleResult{}
	ordinal := 0
	for _, task := range plan.Tasks {
		if statuses[task.ID] == spec.StatusCompleted {
			continue
		}
		if err := ctx.Err(); err != nil {
			if publishErr := engine.publishStop(ctx, plan.RunID, ordinal); publishErr != nil {
				return result, fmt.Errorf("publish stop event for run %q before Task %s: %w", plan.RunID, task.ID, errors.Join(err, publishErr))
			}
			return result, fmt.Errorf("stop run %q before Task %s: %w", plan.RunID, task.ID, err)
		}
		if unmet := unmetNeeds(task, statuses); len(unmet) > 0 {
			result.Skipped++
			reason := fmt.Sprintf("needs not completed: %s", strings.Join(unmet, ", "))
			fmt.Fprintf(engine.deps.Progress, "Task %s skipped: %s\n", task.ID, reason)
			if err := engine.publishTaskEvent(ctx, plan.RunID, 0, task.ID, runevent.KindDaemonTask,
				fmt.Sprintf("Task %s skipped: %s", task.ID, reason),
				map[string]any{"task": task.ID, "phase": "skipped", "reason": reason},
			); err != nil {
				return result, err
			}
			continue
		}
		ordinal++
		settled, err := engine.executeTask(ctx, plan, task, ordinal)
		if err != nil {
			return result, err
		}
		statuses[task.ID] = settled
		if settled == spec.StatusCompleted {
			result.Completed++
		} else {
			result.Failed++
		}
	}
	// QA step seam (task_07): when plan.QA is set and every Task ended
	// completed, the QA step runs here — after the Task walk, before the
	// cycle outcome event — and fills result.QAVerdict.
	if err := engine.publishDaemonEvent(ctx, plan.RunID, 0, runevent.KindDaemonOutcome,
		fmt.Sprintf("Task cycle finished: %d completed, %d failed, %d skipped.", result.Completed, result.Failed, result.Skipped),
		map[string]any{"completed": result.Completed, "failed": result.Failed, "skipped": result.Skipped},
	); err != nil {
		return result, err
	}
	return result, nil
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
	logPath := agent.LogPath(filepath.Join(plan.WorkDir, specRunArtifactDirName), plan.RunID, ordinal)
	fmt.Fprintf(engine.deps.Progress, "Task: %s (Batch %03d) %s\n", task.ID, ordinal, task.Title)
	fmt.Fprintf(engine.deps.Progress, "Agent log: %s\n", logPath)

	_, runErr := engine.deps.Runner.Run(ctx, agent.ExecuteRequest{
		Runtime: plan.Runtime,
		RunID:   plan.RunID,
		Batch:   rounds.Batch{Number: ordinal},
		LogPath: logPath,
		Prompt:  prompt,
		GitRoot: plan.WorkDir,
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
	changed := ensureTaskFile(diffSnapshots(before, after), task.File)
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
// frontmatter (ADR 0013): `<type>: <task title>` where docs, test, and
// chore pass through and every other type maps to feat, plus the
// Roundfix-Spec and Roundfix-Task trailers separated from the subject by a
// blank line per git trailer convention.
func TaskCommitMessage(slug string, task spec.Task) string {
	return fmt.Sprintf("%s: %s\n\nRoundfix-Spec: %s\nRoundfix-Task: %s", taskCommitType(task.Type), task.Title, slug, task.ID)
}

func taskCommitType(taskType string) string {
	switch taskType {
	case "docs", "test", "chore":
		return taskType
	default:
		return "feat"
	}
}

// ensureTaskFile guarantees the settled task file rides in its own Task
// commit even when it was already dirty before the Task started and the
// snapshot diff therefore excluded it.
func ensureTaskFile(changed []string, taskFile string) []string {
	for _, path := range changed {
		if path == taskFile {
			return changed
		}
	}
	changed = append(changed, taskFile)
	sort.Strings(changed)
	return changed
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
		"Run ID":       plan.RunID,
		"working tree": plan.WorkDir,
		"Spec slug":    plan.Spec.Slug,
	}
	for label, value := range required {
		if strings.TrimSpace(value) == "" {
			return fmt.Errorf("task cycle: %s is required", label)
		}
	}
	if len(plan.Tasks) == 0 {
		return errors.New("task cycle: at least one Task is required")
	}
	return nil
}
