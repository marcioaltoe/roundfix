package daemon

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os/exec"
	"sort"
	"strings"
	"sync"
	"time"

	"roundfix/internal/agent"
	"roundfix/internal/reviewsource"
	"roundfix/internal/rounds"
	"roundfix/internal/runevent"
	"roundfix/internal/store"
	runworktree "roundfix/internal/worktree"
)

const publicOutcomeReasonMaxRunes = 240

// ReviewSourceResolver propagates settled Review Issue outcomes to the Review
// Source after verification or terminal Batch failure.
type ReviewSourceResolver interface {
	ResolveIssues(ctx context.Context, req reviewsource.ResolveRequest) error
	ResolveIssue(ctx context.Context, req reviewsource.IssueResolveRequest) error
	ReplyToIssue(ctx context.Context, req reviewsource.IssueCommentRequest) (reviewsource.IssueCommentResult, error)
}

// RunStateStore records intermediate Run states during a cycle. Terminal
// completion stays with the caller.
type RunStateStore interface {
	UpdateRunState(ctx context.Context, runID string, state string) error
	StopRequested(ctx context.Context, runID string) (bool, error)
}

var ErrStopRequested = errors.New("stop requested")

// Dependencies are the engine's explicit collaborators, replacing the CLI
// package globals that previously wired orchestration.
type Dependencies struct {
	Runner        agent.Runner
	Verifier      Verifier
	Committer     Committer
	Pusher        Pusher
	Source        ReviewSourceResolver
	Runs          RunStateStore
	Worktree      WorktreeSnapshotter
	TaskWorktrees TaskWorktreeManager
	PriorChanges  PriorChangedResolver
	Sink          runevent.Sink
	Now           func() time.Time
	Progress      io.Writer
}

// Engine executes one resolve cycle over a validated plan and exposes Final
// Push as a separate explicit operation, so resolve and watch share one
// orchestration implementation. The Daemon owns Batch commits and Final
// Push per ADR 0001.
type Engine struct {
	deps Dependencies
}

type TaskWorktreeManager interface {
	CreateTask(ctx context.Context, run runworktree.Ref, taskID string, opts runworktree.TaskCreateOptions) (runworktree.TaskRef, error)
	IntegrateTask(ctx context.Context, run runworktree.Ref, task runworktree.TaskRef) (runworktree.TaskIntegration, error)
	CleanupTask(ctx context.Context, task runworktree.TaskRef) error
}

type GitTaskWorktreeManager struct{}

func (GitTaskWorktreeManager) CreateTask(ctx context.Context, run runworktree.Ref, taskID string, opts runworktree.TaskCreateOptions) (runworktree.TaskRef, error) {
	return runworktree.CreateTaskWithOptions(ctx, run, taskID, opts)
}

func (GitTaskWorktreeManager) IntegrateTask(ctx context.Context, run runworktree.Ref, task runworktree.TaskRef) (runworktree.TaskIntegration, error) {
	return runworktree.IntegrateTask(ctx, run, task)
}

func (GitTaskWorktreeManager) CleanupTask(ctx context.Context, task runworktree.TaskRef) error {
	return runworktree.CleanupTask(ctx, task)
}

type PriorChangedResolver interface {
	PriorChangedFiles(ctx context.Context, workDir string, initialHead string) ([]string, error)
}

type GitPriorChangedResolver struct{}

func (GitPriorChangedResolver) PriorChangedFiles(ctx context.Context, workDir string, initialHead string) ([]string, error) {
	return runworktree.PriorChangedFiles(ctx, workDir, initialHead)
}

// PullRequestRef identifies the Open Pull Request a cycle works on.
type PullRequestRef struct {
	Number         string
	BaseRepository string
	HeadRepository string
	HeadBranch     string
}

// CyclePlan is the validated input for one resolve cycle: deduplicated
// Review Issues already assembled into Batches for an already-created Run.
type CyclePlan struct {
	RunID        string
	Session      agent.SessionRef
	GitRoot      string
	ArtifactDir  string
	ReviewRoot   string
	AgentLogs    bool
	SourceName   string
	AgentName    string
	Runtime      agent.RuntimeSpec
	Verification string
	AutoCommit   bool
	PullRequest  PullRequestRef
	Batches      []rounds.Batch
	Duplicates   []rounds.DuplicateAssociation
	TotalIssues  int
}

// BatchOutcome reports what one Batch produced. CommitSkipped means
// auto-commit was on but the Agent changed nothing, which is a success.
// Failed means the Batch ended with its assigned Review Issues marked
// failed (Agent error or failed verification); the cycle continues with
// the next Batch and the failed issues are retried in a later Round.
type BatchOutcome struct {
	Batch                 int
	Issues                int
	Failed                bool
	FailureReason         string
	Committed             bool
	CommitSkipped         bool
	ResolvedSourceThreads int
}

func agentLogPath(enabled bool, artifactDir string, runID string, batchNumber int) string {
	if !enabled {
		return ""
	}
	return agent.LogPath(artifactDir, runID, batchNumber)
}

type verificationAttemptRequest struct {
	RunID       string
	WorkDir     string
	ArtifactDir string
	BatchNumber int
	WorkItem    string
	Attempt     int
	Commands    []string
	Publish     func(context.Context, string, map[string]any) error
}

type verificationAttemptOutcome struct {
	Failure        string
	CommandFailure *VerificationCommandError
}

func terminalReasonLine(reason string) string {
	return strings.Join(strings.Fields(reason), " ")
}

func agentTerminalReason(step string, err error) string {
	return agentFailureReason(err, terminalReasonLine(fmt.Sprintf("%s failed: %v", step, err)))
}

func agentFailureReason(err error, fallback string) string {
	if reason, ok := modelNotAdvertisedTerminalReason(err); ok {
		return terminalReasonLine(reason)
	}
	return fallback
}

func modelNotAdvertisedTerminalReason(err error) (string, bool) {
	var modelErr *agent.ModelNotAdvertisedError
	if !errors.As(err, &modelErr) || modelErr == nil {
		return "", false
	}
	advertised := advertisedModelsReason(modelErr.Advertised)
	return fmt.Sprintf("Agent Model %q not advertised by runtime %q; advertised: %s", strings.TrimSpace(modelErr.Model), strings.TrimSpace(modelErr.Runtime), advertised), true
}

func advertisedModelsReason(models []string) string {
	advertised := make([]string, 0, len(models))
	for _, model := range models {
		model = strings.TrimSpace(model)
		if model != "" {
			advertised = append(advertised, model)
		}
	}
	if len(advertised) == 0 {
		return "unavailable"
	}
	return strings.Join(advertised, ", ")
}

func unsettledTerminalReason(step string) string {
	return terminalReasonLine(fmt.Sprintf("%s left issue unsettled after Batch", step))
}

func verificationTerminalReason(commandErr *VerificationCommandError) string {
	if commandErr == nil {
		return ""
	}
	command := strings.TrimSpace(commandErr.Command)
	if command == "" {
		command = "<unknown>"
	}
	diagnostics := strings.TrimSpace(commandErr.OutputPath)
	if diagnostics == "" {
		diagnostics = "unavailable"
	}
	return terminalReasonLine(fmt.Sprintf("Verification failed: command %q exited with %s; diagnostics: %s", command, verificationExitStatus(commandErr), diagnostics))
}

func verificationExitStatus(commandErr *VerificationCommandError) string {
	if commandErr == nil || commandErr.Err == nil {
		return "unknown exit status"
	}
	var exitErr *exec.ExitError
	if errors.As(commandErr.Err, &exitErr) {
		return fmt.Sprintf("exit status %d", exitErr.ExitCode())
	}
	status := strings.TrimSpace(commandErr.Err.Error())
	if status == "" {
		return "unknown exit status"
	}
	return status
}

func (req verificationAttemptRequest) payload(phase runevent.VerificationPhase, command string) map[string]any {
	payload := map[string]any{
		"attempt": req.Attempt,
		"phase":   string(phase),
	}
	if req.BatchNumber > 0 {
		payload["batch"] = req.BatchNumber
	}
	if req.WorkItem != "" {
		payload["work_item"] = req.WorkItem
		payload["task"] = req.WorkItem
	}
	if command != "" {
		payload["command"] = command
	}
	return payload
}

func (req verificationAttemptRequest) summary(phase runevent.VerificationPhase, command string) string {
	target := fmt.Sprintf("Batch %03d", req.BatchNumber)
	if req.WorkItem != "" {
		target = fmt.Sprintf("Task %s", req.WorkItem)
	}
	switch phase {
	case runevent.VerificationPhaseStarted:
		return fmt.Sprintf("Verification attempt %d for %s started: %s", req.Attempt, target, command)
	case runevent.VerificationPhaseCommandPassed:
		return fmt.Sprintf("Verification attempt %d for %s command passed: %s", req.Attempt, target, command)
	case runevent.VerificationPhaseFailed:
		return fmt.Sprintf("Verification attempt %d for %s failed: %s", req.Attempt, target, command)
	case runevent.VerificationPhaseVerdict:
		return fmt.Sprintf("Verification attempt %d for %s verdict recorded.", req.Attempt, target)
	default:
		return fmt.Sprintf("Verification attempt %d for %s phase %s.", req.Attempt, target, phase)
	}
}

func (engine *Engine) runVerificationAttempt(ctx context.Context, req verificationAttemptRequest) (verificationAttemptOutcome, error) {
	if req.Attempt < 1 {
		return verificationAttemptOutcome{}, fmt.Errorf("run verification: attempt is required")
	}
	if req.Publish == nil {
		return verificationAttemptOutcome{}, fmt.Errorf("run verification attempt %d: event publisher is required", req.Attempt)
	}
	for _, command := range req.Commands {
		if err := req.Publish(ctx, req.summary(runevent.VerificationPhaseStarted, command), req.payload(runevent.VerificationPhaseStarted, command)); err != nil {
			return verificationAttemptOutcome{}, err
		}
		outputPath := VerificationOutputPath(req.ArtifactDir, req.RunID, req.BatchNumber, req.Attempt)
		_, err := engine.deps.Verifier.Verify(ctx, VerifyRequest{
			WorkDir:    req.WorkDir,
			Command:    command,
			OutputPath: outputPath,
		})
		if err != nil {
			if isStop(ctx, err) {
				return verificationAttemptOutcome{}, err
			}
			var commandErr *VerificationCommandError
			if errors.As(err, &commandErr) {
				if err := req.publishFailedCommand(ctx, command, commandErr); err != nil {
					return verificationAttemptOutcome{}, err
				}
				fmt.Fprintf(engine.deps.Progress, "Verification failed (attempt %d); diagnostics: %s\n", req.Attempt, commandErr.OutputPath)
				return verificationAttemptOutcome{
					Failure:        fmt.Sprintf("verification failed: %v", err),
					CommandFailure: commandErr,
				}, nil
			}
			if err := req.publishInfrastructureFailure(ctx, command, err); err != nil {
				return verificationAttemptOutcome{}, err
			}
			return verificationAttemptOutcome{}, fmt.Errorf("run verification attempt %d: %w", req.Attempt, err)
		}
		if err := req.Publish(ctx, req.summary(runevent.VerificationPhaseCommandPassed, command), req.payload(runevent.VerificationPhaseCommandPassed, command)); err != nil {
			return verificationAttemptOutcome{}, err
		}
	}
	if err := req.publishVerdict(ctx, runevent.VerificationVerdictPassed, "", ""); err != nil {
		return verificationAttemptOutcome{}, err
	}
	fmt.Fprintf(engine.deps.Progress, "Verification passed (attempt %d).\n", req.Attempt)
	return verificationAttemptOutcome{}, nil
}

func (req verificationAttemptRequest) publishFailedCommand(ctx context.Context, command string, commandErr *VerificationCommandError) error {
	payload := req.payload(runevent.VerificationPhaseFailed, command)
	payload["error"] = commandErr.Error()
	if commandErr.OutputPath != "" {
		payload["diagnostic_path"] = commandErr.OutputPath
	}
	if err := req.Publish(ctx, req.summary(runevent.VerificationPhaseFailed, command), payload); err != nil {
		return err
	}
	return req.publishVerdict(ctx, runevent.VerificationVerdictFailed, commandErr.OutputPath, "")
}

func (req verificationAttemptRequest) publishInfrastructureFailure(ctx context.Context, command string, err error) error {
	payload := req.payload(runevent.VerificationPhaseFailed, command)
	payload["error"] = err.Error()
	if publishErr := req.Publish(ctx, req.summary(runevent.VerificationPhaseFailed, command), payload); publishErr != nil {
		return publishErr
	}
	return req.publishVerdict(ctx, runevent.VerificationVerdictFailed, "", "")
}

func (req verificationAttemptRequest) publishVerdict(ctx context.Context, verdict runevent.VerificationVerdict, diagnosticPath string, failure string) error {
	payload := req.payload(runevent.VerificationPhaseVerdict, "")
	payload["verdict"] = string(verdict)
	if diagnosticPath != "" {
		payload["diagnostic_path"] = diagnosticPath
	}
	if failure != "" {
		payload["error"] = failure
	}
	summary := fmt.Sprintf("Verification attempt %d verdict: %s", req.Attempt, verdict)
	if req.WorkItem != "" {
		summary = fmt.Sprintf("Verification attempt %d for Task %s verdict: %s", req.Attempt, req.WorkItem, verdict)
	} else if req.BatchNumber > 0 {
		summary = fmt.Sprintf("Verification attempt %d for Batch %03d verdict: %s", req.Attempt, req.BatchNumber, verdict)
	}
	return req.Publish(ctx, summary, payload)
}

// CycleResult reports per-Batch outcomes and the remaining Unresolved
// Review Issue count after the cycle.
type CycleResult struct {
	Batches   []BatchOutcome
	Remaining int
}

// FinalPushRequest names the push target for the engine's second operation.
// Push gating policy (no Unresolved Review Issues, auto-push enabled) stays
// with the caller.
type FinalPushRequest struct {
	RunID   string
	WorkDir string
	Remote  string
	Branch  string
}

func NewEngine(deps Dependencies) (*Engine, error) {
	missing := []string{}
	if deps.Runner == nil {
		missing = append(missing, "Runner")
	}
	if deps.Verifier == nil {
		missing = append(missing, "Verifier")
	}
	if deps.Committer == nil {
		missing = append(missing, "Committer")
	}
	if deps.Pusher == nil {
		missing = append(missing, "Pusher")
	}
	if deps.Source == nil {
		missing = append(missing, "Source")
	}
	if deps.Runs == nil {
		missing = append(missing, "Runs")
	}
	if deps.Worktree == nil {
		missing = append(missing, "Worktree")
	}
	if len(missing) > 0 {
		return nil, fmt.Errorf("create Run engine: missing dependencies: %s", strings.Join(missing, ", "))
	}
	if deps.Sink == nil {
		deps.Sink = runevent.Discard
	}
	if deps.TaskWorktrees == nil {
		deps.TaskWorktrees = GitTaskWorktreeManager{}
	}
	if deps.PriorChanges == nil {
		deps.PriorChanges = GitPriorChangedResolver{}
	}
	if deps.Now == nil {
		deps.Now = time.Now
	}
	if deps.Progress == nil {
		deps.Progress = io.Discard
	}
	deps.Progress = &lockedWriter{writer: deps.Progress}
	return &Engine{deps: deps}, nil
}

type lockedWriter struct {
	mu     sync.Mutex
	writer io.Writer
}

func (writer *lockedWriter) Write(data []byte) (int, error) {
	writer.mu.Lock()
	defer writer.mu.Unlock()
	return writer.writer.Write(data)
}

// ResolveCycle executes one resolve cycle: for each Batch it runs the
// Agent, settles assigned issue statuses, creates the Batch commit when
// auto-commit is enabled, and propagates settled Review Issue outcomes to the
// Review Source. A failed Batch (Agent error or failed verification) marks
// its assigned issues failed and the cycle continues with the next Batch; only
// Stop Requests and infrastructure errors halt the cycle. A Stop Request halts
// before any new Batch, verification, commit, or Review Source mutation; Agent
// worktree changes are preserved.
func (engine *Engine) ResolveCycle(ctx context.Context, plan CyclePlan) (CycleResult, error) {
	if err := validateCyclePlan(plan); err != nil {
		return CycleResult{}, err
	}
	if err := engine.publishDaemonEvent(ctx, plan.RunID, 0, runevent.KindDaemonSelection,
		fmt.Sprintf("Selected %d Review Issue(s) into %d Batch(es); %d duplicate occurrence(s) associated.", plan.TotalIssues, len(plan.Batches), len(plan.Duplicates)),
		map[string]any{"issues": plan.TotalIssues, "batches": len(plan.Batches), "duplicates": len(plan.Duplicates)},
	); err != nil {
		return CycleResult{}, err
	}
	result := CycleResult{Remaining: plan.TotalIssues}
	for index, batch := range plan.Batches {
		if err := ctx.Err(); err != nil {
			if publishErr := engine.publishStop(ctx, plan.RunID, batch.Number); publishErr != nil {
				return result, fmt.Errorf("publish stop event for run %q before Batch %03d: %w", plan.RunID, batch.Number, errors.Join(err, publishErr))
			}
			return result, fmt.Errorf("stop run %q before Batch %03d: %w", plan.RunID, batch.Number, err)
		}
		if err := engine.stopIfRequested(ctx, plan.RunID, batch.Number); err != nil {
			return result, fmt.Errorf("stop run %q before Batch %03d: %w", plan.RunID, batch.Number, err)
		}
		outcome, remaining, err := engine.resolveBatch(ctx, plan, batch, index+1, len(plan.Batches))
		if err != nil {
			engine.reportPending(plan, index)
			return result, err
		}
		result.Batches = append(result.Batches, outcome)
		result.Remaining = remaining
		if !outcome.Failed && remaining > 0 && index < len(plan.Batches)-1 {
			fmt.Fprintf(engine.deps.Progress, "Batch %03d/%03d completed; %d Unresolved Review Issue(s) remain.\n", batch.Number, len(plan.Batches), remaining)
		}
		if err := engine.stopIfRequested(ctx, plan.RunID, batch.Number); err != nil {
			return result, fmt.Errorf("stop run %q after Batch %03d settlement: %w", plan.RunID, batch.Number, err)
		}
	}
	if result.Remaining > 0 {
		if err := engine.propagateRunEndUnresolved(ctx, plan); err != nil {
			return result, err
		}
	}
	return result, nil
}

// FinalPush sends the PR Head Branch. It is invoked explicitly by the
// caller, never per Batch or Round, preserving ADR 0001 semantics.
func (engine *Engine) FinalPush(ctx context.Context, req FinalPushRequest) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	if strings.TrimSpace(req.RunID) != "" {
		if err := engine.deps.Runs.UpdateRunState(ctx, req.RunID, store.StatePushing); err != nil {
			return err
		}
	}
	if err := engine.deps.Pusher.Push(ctx, PushRequest{
		WorkDir: req.WorkDir,
		Remote:  req.Remote,
		Branch:  req.Branch,
	}); err != nil {
		return err
	}
	return engine.publishDaemonEvent(ctx, req.RunID, 0, runevent.KindDaemonPush,
		fmt.Sprintf("Final Push completed: git push %s HEAD:%s", req.Remote, req.Branch),
		map[string]any{"decision": "pushed", "remote": req.Remote, "branch": req.Branch},
	)
}

func (engine *Engine) resolveBatch(ctx context.Context, plan CyclePlan, batch rounds.Batch, batchIndex int, batchTotal int) (BatchOutcome, int, error) {
	outcome := BatchOutcome{Batch: batch.Number, Issues: len(batch.Issues)}
	// The before-snapshot is taken at Batch start, so anything already
	// dirty — pre-existing user work or edits from earlier in the Run —
	// never reaches a Batch commit.
	var before []string
	if plan.AutoCommit {
		snapshot, err := engine.deps.Worktree.Snapshot(ctx, plan.GitRoot)
		if err != nil {
			return outcome, 0, err
		}
		before = snapshot
	}
	failure, err := engine.runBatchAgent(ctx, plan, batch, batchIndex, batchTotal)
	if err != nil {
		return outcome, 0, err
	}
	if failure == "" {
		verification, verifyErr := engine.verifyBatch(ctx, plan, batch, 1)
		if verifyErr != nil {
			return outcome, 0, verifyErr
		}
		if verification.Failure != "" {
			failure, err = engine.repairBatchVerification(ctx, plan, batch, verification)
			if err != nil {
				return outcome, 0, err
			}
		}
	}
	if failure != "" {
		return engine.completeFailedBatch(ctx, plan, batch, outcome, failure)
	}
	marked, err := rounds.MarkDuplicatedAfterTerminal(ctx, plan.Duplicates)
	if err != nil {
		return outcome, 0, fmt.Errorf("mark duplicate Review Issues after terminal outcomes for run %q batch %03d: %w", plan.RunID, batch.Number, err)
	}
	if marked > 0 {
		fmt.Fprintf(engine.deps.Progress, "Marked %d older duplicate Review Issue occurrence(s) as duplicated.\n", marked)
	}
	duplicates, err := duplicatedIssuesForBatch(ctx, batch, plan.Duplicates)
	if err != nil {
		return outcome, 0, fmt.Errorf("load duplicated Review Issues for run %q Batch %03d: %w", plan.RunID, batch.Number, err)
	}
	committed, skipped, err := engine.commitBatch(ctx, plan, batch, before)
	if err != nil {
		return outcome, 0, err
	}
	outcome.Committed = committed
	outcome.CommitSkipped = skipped
	propagated, err := engine.propagateBatchSources(ctx, plan, batch, duplicates)
	if err != nil {
		return outcome, 0, err
	}
	outcome.ResolvedSourceThreads = propagated.Resolved
	if propagated.Failed > 0 {
		outcome.Failed = true
		outcome.FailureReason = fmt.Sprintf("Review Source propagation failed for %d Review Issue(s)", propagated.Failed)
	}
	remaining, err := remainingUnresolvedIssues(ctx, plan)
	if err != nil {
		return outcome, 0, fmt.Errorf("compute remaining unresolved issues for run %q after Batch %03d: %w", plan.RunID, batch.Number, err)
	}
	phase := "completed"
	summary := fmt.Sprintf("Batch %03d completed; %d Unresolved Review Issue(s) remain.", batch.Number, remaining)
	payload := map[string]any{"phase": phase, "batch": batch.Number, "remaining": remaining}
	if propagated.Failed > 0 {
		phase = "failed"
		summary = fmt.Sprintf("Batch %03d completed with Review Source propagation failure; %d Unresolved Review Issue(s) remain.", batch.Number, remaining)
		payload["phase"] = phase
		payload["source_propagation_failed"] = propagated.Failed
	}
	if err := engine.publishDaemonEvent(ctx, plan.RunID, batch.Number, runevent.KindDaemonBatch, summary, payload); err != nil {
		return outcome, 0, fmt.Errorf("publish completion event for run %q batch %03d: %w", plan.RunID, batch.Number, err)
	}
	return outcome, remaining, nil
}

// runBatchAgent runs the Agent over one Batch. It returns a non-empty
// failure reason when the Batch failed but the cycle should continue;
// the returned error is reserved for Stop Requests and infrastructure
// failures, which halt the cycle.
func (engine *Engine) runBatchAgent(ctx context.Context, plan CyclePlan, batch rounds.Batch, batchIndex int, batchTotal int) (string, error) {
	if err := engine.deps.Runs.UpdateRunState(ctx, plan.RunID, store.StateResolvingWithAgent); err != nil {
		return "", fmt.Errorf("update run %q to state %q before Batch %03d: %w", plan.RunID, store.StateResolvingWithAgent, batch.Number, err)
	}
	if err := engine.publishDaemonEvent(ctx, plan.RunID, batch.Number, runevent.KindDaemonBatch,
		fmt.Sprintf("Batch %03d/%03d started with %d Review Issue(s).", batchIndex, batchTotal, len(batch.Issues)),
		map[string]any{"phase": "started", "batch": batch.Number, "issues": len(batch.Issues)},
	); err != nil {
		return "", fmt.Errorf("publish start event for run %q batch %03d: %w", plan.RunID, batch.Number, err)
	}
	prompt := agent.BuildPrompt(agent.PromptRequest{
		RunID:        plan.RunID,
		Batch:        batch,
		Agent:        plan.AgentName,
		Model:        plan.Runtime.Model,
		ArtifactDir:  plan.ArtifactDir,
		GitRoot:      plan.GitRoot,
		Verification: plan.Verification,
	})
	logPath := agentLogPath(plan.AgentLogs, plan.ArtifactDir, plan.RunID, batch.Number)
	fmt.Fprintf(engine.deps.Progress, "Batch: %03d/%03d (%d Review Issue(s))\n", batchIndex, batchTotal, len(batch.Issues))
	if logPath != "" {
		fmt.Fprintf(engine.deps.Progress, "Agent log: %s\n", logPath)
	}

	runResult, runErr := engine.deps.Runner.Run(ctx, agent.ExecuteRequest{
		Runtime:      plan.Runtime,
		Session:      plan.Session,
		RunID:        plan.RunID,
		Batch:        batch,
		LogPath:      logPath,
		Prompt:       prompt,
		ArtifactDir:  plan.ArtifactDir,
		GitRoot:      plan.GitRoot,
		Verification: plan.Verification,
	}, engine.deps.Sink)
	if runErr != nil {
		if isStop(ctx, runErr) {
			// The runner already published the stopped status event;
			// Agent-created worktree changes stay untouched.
			return "", fmt.Errorf("run Agent for run %q batch %03d: %w", plan.RunID, batch.Number, runErr)
		}
		reason := agentTerminalReason("Agent", runErr)
		if err := agent.MarkBatchFailed(batch, reason); err != nil {
			return "", fmt.Errorf("mark run %q batch %03d failed after Agent error: %w", plan.RunID, batch.Number, err)
		}
		return reason, nil
	}
	if err := ctx.Err(); err != nil {
		if publishErr := engine.publishStop(ctx, plan.RunID, batch.Number); publishErr != nil {
			return "", fmt.Errorf("publish stop event for run %q after Agent batch %03d: %w", plan.RunID, batch.Number, errors.Join(err, publishErr))
		}
		return "", fmt.Errorf("stop run %q after Agent batch %03d: %w", plan.RunID, batch.Number, err)
	}
	if anomaly := strings.TrimSpace(runResult.TransportAnomaly); anomaly != "" {
		if err := engine.publishDaemonEvent(ctx, plan.RunID, batch.Number, runevent.KindDaemonBatch,
			fmt.Sprintf("Batch %03d transport anomaly: %s", batch.Number, anomaly),
			map[string]any{"phase": "transport_anomaly", "batch": batch.Number, "anomaly": anomaly},
		); err != nil {
			return "", fmt.Errorf("publish transport anomaly event for run %q batch %03d: %w", plan.RunID, batch.Number, err)
		}
	}
	settled, err := agent.SettleAssignedIssues(ctx, batch, unsettledTerminalReason("Agent"))
	if err != nil {
		return "", fmt.Errorf("settle assigned Review Issues for run %q batch %03d: %w", plan.RunID, batch.Number, err)
	}
	if len(settled) > 0 {
		fmt.Fprintf(engine.deps.Progress, "Marked %d assigned Review Issue(s) the Agent left unsettled as failed.\n", len(settled))
		if err := engine.publishDaemonEvent(ctx, plan.RunID, batch.Number, runevent.KindDaemonBatch,
			fmt.Sprintf("Marked %d assigned Review Issue(s) the Agent left unsettled as failed.", len(settled)),
			map[string]any{"phase": "settled", "batch": batch.Number, "failed": len(settled)},
		); err != nil {
			return "", fmt.Errorf("publish settlement event for run %q batch %03d: %w", plan.RunID, batch.Number, err)
		}
	}
	if _, err := fmt.Fprintln(engine.deps.Progress, "Agent Batch finished with settled Review Issue statuses."); err != nil {
		return "", fmt.Errorf("write batch completion progress for run %q batch %03d: %w", plan.RunID, batch.Number, err)
	}
	return "", nil
}

// verifyBatch runs one configured verification attempt. Command failures
// return a typed outcome for the repair loop; the returned error is reserved
// for Stop Requests and infrastructure failures, which halt the cycle.
func (engine *Engine) verifyBatch(ctx context.Context, plan CyclePlan, batch rounds.Batch, attempt int) (verificationAttemptOutcome, error) {
	if err := engine.deps.Runs.UpdateRunState(ctx, plan.RunID, store.StateVerifying); err != nil {
		return verificationAttemptOutcome{}, fmt.Errorf("update run %q to state %q before Batch %03d verification: %w", plan.RunID, store.StateVerifying, batch.Number, err)
	}
	verification, err := engine.runVerificationAttempt(ctx, verificationAttemptRequest{
		RunID:       plan.RunID,
		WorkDir:     plan.GitRoot,
		ArtifactDir: plan.ArtifactDir,
		BatchNumber: batch.Number,
		Attempt:     attempt,
		Commands:    []string{plan.Verification},
		Publish: func(ctx context.Context, summary string, payload map[string]any) error {
			if err := engine.publishDaemonEvent(ctx, plan.RunID, batch.Number, runevent.KindDaemonVerification, summary, payload); err != nil {
				return fmt.Errorf("publish verification event for run %q batch %03d: %w", plan.RunID, batch.Number, err)
			}
			return nil
		},
	})
	if err != nil {
		if isStop(ctx, err) {
			// A Stop Request during verification keeps Agent statuses
			// untouched; the run ends Stopped, not failed.
			return verificationAttemptOutcome{}, fmt.Errorf("verify run %q batch %03d: %w", plan.RunID, batch.Number, err)
		}
		return verificationAttemptOutcome{}, fmt.Errorf("verify run %q batch %03d: %w", plan.RunID, batch.Number, err)
	}
	return verification, nil
}

func (engine *Engine) repairBatchVerification(ctx context.Context, plan CyclePlan, batch rounds.Batch, first verificationAttemptOutcome) (string, error) {
	if first.CommandFailure == nil {
		return first.Failure, nil
	}
	failure, err := engine.runBatchVerificationRepair(ctx, plan, batch, first)
	if err != nil {
		return "", err
	}
	if failure != "" {
		return failure, nil
	}
	final, err := engine.verifyBatch(ctx, plan, batch, 2)
	if err != nil {
		return "", err
	}
	if final.Failure != "" {
		reason := verificationTerminalReason(final.CommandFailure)
		if reason == "" {
			reason = terminalReasonLine(final.Failure)
		}
		if markErr := agent.MarkBatchFailed(batch, reason); markErr != nil {
			return "", fmt.Errorf("mark run %q batch %03d failed after verification error: %w", plan.RunID, batch.Number, markErr)
		}
		return final.Failure, nil
	}
	return "", nil
}

func (engine *Engine) runBatchVerificationRepair(ctx context.Context, plan CyclePlan, batch rounds.Batch, first verificationAttemptOutcome) (string, error) {
	if err := ctx.Err(); err != nil {
		if publishErr := engine.publishStop(ctx, plan.RunID, batch.Number); publishErr != nil {
			return "", fmt.Errorf("publish stop event for run %q before Batch %03d Verification Feedback: %w", plan.RunID, batch.Number, errors.Join(err, publishErr))
		}
		return "", fmt.Errorf("stop run %q before Batch %03d Verification Feedback: %w", plan.RunID, batch.Number, err)
	}
	if err := engine.stopIfRequested(ctx, plan.RunID, batch.Number); err != nil {
		return "", fmt.Errorf("stop run %q before Batch %03d Verification Feedback: %w", plan.RunID, batch.Number, err)
	}
	if err := engine.deps.Runs.UpdateRunState(ctx, plan.RunID, store.StateResolvingWithAgent); err != nil {
		return "", fmt.Errorf("update run %q to state %q before Batch %03d Verification Feedback: %w", plan.RunID, store.StateResolvingWithAgent, batch.Number, err)
	}
	prompt, err := agent.BuildVerificationRepairPrompt(fmt.Sprintf("Batch %03d", batch.Number), agent.VerificationFeedback{
		Command:        first.CommandFailure.Command,
		DiagnosticPath: first.CommandFailure.OutputPath,
		Failure:        first.Failure,
		Attempt:        1,
	})
	if err != nil {
		return "", fmt.Errorf("build Verification Feedback prompt for run %q batch %03d: %w", plan.RunID, batch.Number, err)
	}
	logPath := agentLogPath(plan.AgentLogs, plan.ArtifactDir, plan.RunID, batch.Number)
	fmt.Fprintf(engine.deps.Progress, "Verification Feedback: Batch %03d\n", batch.Number)
	if logPath != "" {
		fmt.Fprintf(engine.deps.Progress, "Agent log: %s\n", logPath)
	}
	runResult, runErr := engine.deps.Runner.Run(ctx, agent.ExecuteRequest{
		Runtime:      plan.Runtime,
		Session:      plan.Session,
		RunID:        plan.RunID,
		Batch:        batch,
		LogPath:      logPath,
		Prompt:       prompt,
		ArtifactDir:  plan.ArtifactDir,
		GitRoot:      plan.GitRoot,
		Verification: plan.Verification,
	}, engine.deps.Sink)
	if runErr != nil {
		if isStop(ctx, runErr) {
			return "", fmt.Errorf("run Verification Feedback Agent for run %q batch %03d: %w", plan.RunID, batch.Number, runErr)
		}
		reason := agentTerminalReason("Verification Feedback Agent", runErr)
		if err := agent.MarkBatchFailed(batch, reason); err != nil {
			return "", fmt.Errorf("mark run %q batch %03d failed after Verification Feedback Agent error: %w", plan.RunID, batch.Number, err)
		}
		return reason, nil
	}
	if err := ctx.Err(); err != nil {
		if publishErr := engine.publishStop(ctx, plan.RunID, batch.Number); publishErr != nil {
			return "", fmt.Errorf("publish stop event for run %q after Batch %03d Verification Feedback: %w", plan.RunID, batch.Number, errors.Join(err, publishErr))
		}
		return "", fmt.Errorf("stop run %q after Batch %03d Verification Feedback: %w", plan.RunID, batch.Number, err)
	}
	if anomaly := strings.TrimSpace(runResult.TransportAnomaly); anomaly != "" {
		if err := engine.publishDaemonEvent(ctx, plan.RunID, batch.Number, runevent.KindDaemonBatch,
			fmt.Sprintf("Batch %03d Verification Feedback transport anomaly: %s", batch.Number, anomaly),
			map[string]any{"phase": "verification_feedback_transport_anomaly", "batch": batch.Number, "anomaly": anomaly},
		); err != nil {
			return "", fmt.Errorf("publish Verification Feedback transport anomaly event for run %q batch %03d: %w", plan.RunID, batch.Number, err)
		}
	}
	settled, err := agent.SettleAssignedIssues(ctx, batch, unsettledTerminalReason("Verification Feedback Agent"))
	if err != nil {
		return "", fmt.Errorf("settle assigned Review Issues after Verification Feedback for run %q batch %03d: %w", plan.RunID, batch.Number, err)
	}
	if len(settled) > 0 {
		fmt.Fprintf(engine.deps.Progress, "Marked %d assigned Review Issue(s) the Verification Feedback Agent left unsettled as failed.\n", len(settled))
		if err := engine.publishDaemonEvent(ctx, plan.RunID, batch.Number, runevent.KindDaemonBatch,
			fmt.Sprintf("Marked %d assigned Review Issue(s) the Verification Feedback Agent left unsettled as failed.", len(settled)),
			map[string]any{"phase": "settled", "batch": batch.Number, "failed": len(settled)},
		); err != nil {
			return "", fmt.Errorf("publish Verification Feedback settlement event for run %q batch %03d: %w", plan.RunID, batch.Number, err)
		}
	}
	if _, err := fmt.Fprintln(engine.deps.Progress, "Verification Feedback Agent finished with settled Review Issue statuses."); err != nil {
		return "", fmt.Errorf("write Verification Feedback progress for run %q batch %03d: %w", plan.RunID, batch.Number, err)
	}
	return "", nil
}

// completeFailedBatch records a failed Batch outcome and the remaining
// Unresolved Review Issue count so the cycle can continue with the next
// Batch. The Batch issues were already marked failed at the failure site.
func (engine *Engine) completeFailedBatch(ctx context.Context, plan CyclePlan, batch rounds.Batch, outcome BatchOutcome, failure string) (BatchOutcome, int, error) {
	outcome.Failed = true
	outcome.FailureReason = failure
	if _, err := engine.propagateBatchSources(ctx, plan, batch, nil); err != nil {
		return outcome, 0, err
	}
	remaining, err := remainingUnresolvedIssues(ctx, plan)
	if err != nil {
		return outcome, 0, fmt.Errorf("compute remaining unresolved issues for run %q after failed Batch %03d: %w", plan.RunID, batch.Number, err)
	}
	fmt.Fprintf(engine.deps.Progress, "Batch %03d failed: %s\n", batch.Number, failure)
	fmt.Fprintln(engine.deps.Progress, "The tracked checkout was clean at Preflight; inspect current tracked and untracked changes before retrying.")
	if err := engine.publishDaemonEvent(ctx, plan.RunID, batch.Number, runevent.KindDaemonBatch,
		fmt.Sprintf("Batch %03d failed; %d Unresolved Review Issue(s) remain.", batch.Number, remaining),
		map[string]any{"phase": "failed", "batch": batch.Number, "remaining": remaining, "error": failure},
	); err != nil {
		return outcome, 0, fmt.Errorf("publish failure event for run %q batch %03d: %w", plan.RunID, batch.Number, err)
	}
	return outcome, remaining, nil
}

func (engine *Engine) commitBatch(ctx context.Context, plan CyclePlan, batch rounds.Batch, before []string) (bool, bool, error) {
	if !plan.AutoCommit {
		fmt.Fprintln(engine.deps.Progress, "Auto-commit disabled; no Batch commit created.")
		err := engine.publishDaemonEvent(ctx, plan.RunID, batch.Number, runevent.KindDaemonCommit,
			"Auto-commit disabled; no Batch commit created.",
			map[string]any{"decision": "disabled", "batch": batch.Number},
		)
		return false, false, err
	}
	if err := ctx.Err(); err != nil {
		if publishErr := engine.publishStop(ctx, plan.RunID, batch.Number); publishErr != nil {
			return false, false, fmt.Errorf("publish stop event for run %q before Batch %03d commit: %w", plan.RunID, batch.Number, errors.Join(err, publishErr))
		}
		return false, false, fmt.Errorf("stop run %q before Batch %03d commit: %w", plan.RunID, batch.Number, err)
	}
	after, err := engine.deps.Worktree.Snapshot(ctx, plan.GitRoot)
	if err != nil {
		return false, false, err
	}
	changed := diffSnapshots(before, after)
	if len(changed) == 0 {
		// A triage-only Batch changed nothing: skipping the commit is a
		// success, never a nothing-to-commit failure.
		fmt.Fprintf(engine.deps.Progress, "Batch commit skipped: Batch %03d made no worktree changes.\n", batch.Number)
		err := engine.publishDaemonEvent(ctx, plan.RunID, batch.Number, runevent.KindDaemonCommit,
			fmt.Sprintf("Batch commit skipped: Batch %03d made no worktree changes.", batch.Number),
			map[string]any{"decision": "skipped", "batch": batch.Number},
		)
		return false, true, err
	}
	message := BatchCommitMessage(batch.Number)
	if err := engine.deps.Committer.Commit(ctx, CommitRequest{
		WorkDir: plan.GitRoot,
		Message: message,
		Paths:   changed,
	}); err != nil {
		return false, false, err
	}
	fmt.Fprintf(engine.deps.Progress, "Batch commit created: %s\n", message)
	if err := engine.publishDaemonEvent(ctx, plan.RunID, batch.Number, runevent.KindDaemonCommit,
		fmt.Sprintf("Batch commit created: %s", message),
		map[string]any{"decision": "created", "batch": batch.Number, "paths": len(changed)},
	); err != nil {
		return false, false, err
	}
	return true, false, nil
}

// diffSnapshots returns the paths dirty after the Batch that were not
// already dirty before it, sorted for deterministic staging. Project Config
// stays excluded as defense in depth.
func diffSnapshots(before []string, after []string) []string {
	seen := make(map[string]bool, len(before))
	for _, path := range before {
		seen[path] = true
	}
	changed := []string{}
	for _, path := range after {
		if seen[path] || path == ".roundfixrc.yml" {
			continue
		}
		seen[path] = true
		changed = append(changed, path)
	}
	sort.Strings(changed)
	return changed
}

type sourcePropagationSummary struct {
	Resolved int
	Failed   int
}

func (engine *Engine) propagateBatchSources(ctx context.Context, plan CyclePlan, batch rounds.Batch, duplicated []rounds.Issue) (sourcePropagationSummary, error) {
	issues, err := settledBatchSourceIssues(batch, duplicated)
	if err != nil {
		return sourcePropagationSummary{}, err
	}
	if len(issues) == 0 {
		return sourcePropagationSummary{}, nil
	}
	if err := ctx.Err(); err != nil {
		if publishErr := engine.publishStop(ctx, plan.RunID, batch.Number); publishErr != nil {
			return sourcePropagationSummary{}, fmt.Errorf("publish stop event for run %q before Batch %03d source propagation: %w", plan.RunID, batch.Number, errors.Join(err, publishErr))
		}
		return sourcePropagationSummary{}, fmt.Errorf("stop run %q before Batch %03d source propagation: %w", plan.RunID, batch.Number, err)
	}
	summary := sourcePropagationSummary{}
	resolvedRefs := map[string]bool{}
	for _, issue := range issues {
		action, err := engine.propagateSourceIssue(ctx, plan, batch.Number, issue, "batch", resolvedRefs)
		if err != nil {
			return summary, err
		}
		if action.Resolved {
			summary.Resolved++
		}
		if action.Failed {
			summary.Failed++
		}
	}
	if summary.Resolved > 0 {
		fmt.Fprintf(engine.deps.Progress, "Resolved %d Review Source thread(s).\n", summary.Resolved)
	}
	return summary, nil
}

func (engine *Engine) propagateRunEndUnresolved(ctx context.Context, plan CyclePlan) error {
	issues, err := unresolvedSourceIssues(ctx, plan)
	if err != nil {
		return fmt.Errorf("load unresolved Review Issues for run %q source propagation: %w", plan.RunID, err)
	}
	for _, issue := range issues {
		if _, err := engine.propagateSourceIssue(ctx, plan, 0, issue, "run_end", nil); err != nil {
			return err
		}
	}
	return nil
}

type sourceIssueAction struct {
	Action   string
	Comment  bool
	Resolve  bool
	Resolved bool
	Failed   bool
}

func (engine *Engine) propagateSourceIssue(ctx context.Context, plan CyclePlan, batchNumber int, issue rounds.Issue, phase string, resolvedRefs map[string]bool) (sourceIssueAction, error) {
	if strings.TrimSpace(issue.SourceRef) == "" {
		return sourceIssueAction{Action: "skipped_no_source_ref"}, nil
	}
	action := sourceActionForIssue(issue, phase)
	if action.Action == "" {
		return sourceIssueAction{Action: "skipped_status"}, nil
	}
	if action.Comment {
		marker := outcomeCommentMarker(plan.RunID, issue, action.Action)
		result, err := engine.deps.Source.ReplyToIssue(ctx, reviewsource.IssueCommentRequest{
			Source:         plan.SourceName,
			PRNumber:       plan.PullRequest.Number,
			BaseRepository: plan.PullRequest.BaseRepository,
			SourceRef:      issue.SourceRef,
			Marker:         marker,
			Body:           outcomeCommentBody(plan.RunID, issue, action.Action, marker),
		})
		if err != nil {
			action.Failed = true
			if markErr := markSourcePropagationFailed(issue, action.Action, err); markErr != nil {
				return action, markErr
			}
			engine.reportSourcePropagationFailure(issue, action.Action, err)
			if eventErr := engine.publishSourcePropagationEvent(ctx, plan, batchNumber, issue, phase, action.Action, false, false, err); eventErr != nil {
				return action, eventErr
			}
			return action, nil
		}
		if eventErr := engine.publishSourcePropagationEvent(ctx, plan, batchNumber, issue, phase, action.Action, result.Posted, false, nil); eventErr != nil {
			return action, eventErr
		}
	}
	if action.Resolve {
		sourceRef := strings.TrimSpace(issue.SourceRef)
		if resolvedRefs != nil && resolvedRefs[sourceRef] {
			if eventErr := engine.publishSourcePropagationEvent(ctx, plan, batchNumber, issue, phase, "resolve_skipped_already_resolved", false, false, nil); eventErr != nil {
				return action, eventErr
			}
			return action, nil
		}
		if err := engine.deps.Source.ResolveIssue(ctx, reviewsource.IssueResolveRequest{
			Source:         plan.SourceName,
			PRNumber:       plan.PullRequest.Number,
			BaseRepository: plan.PullRequest.BaseRepository,
			SourceRef:      issue.SourceRef,
		}); err != nil {
			action.Failed = true
			if markErr := markSourcePropagationFailed(issue, "resolve", err); markErr != nil {
				return action, markErr
			}
			engine.reportSourcePropagationFailure(issue, "resolve", err)
			if eventErr := engine.publishSourcePropagationEvent(ctx, plan, batchNumber, issue, phase, "resolve", false, false, err); eventErr != nil {
				return action, eventErr
			}
			return action, nil
		}
		if resolvedRefs != nil {
			resolvedRefs[sourceRef] = true
		}
		action.Resolved = true
		if eventErr := engine.publishSourcePropagationEvent(ctx, plan, batchNumber, issue, phase, "resolve", false, true, nil); eventErr != nil {
			return action, eventErr
		}
	}
	return action, nil
}

func markSourcePropagationFailed(issue rounds.Issue, action string, cause error) error {
	if strings.TrimSpace(issue.Path) == "" {
		return nil
	}
	if issue.Status == rounds.StatusFailed {
		return nil
	}
	reason := terminalReasonLine(fmt.Sprintf("Review Source propagation failed during %s: %v", action, cause))
	if err := rounds.SetIssueStatus(issue.Path, rounds.StatusFailed, "", reason); err != nil {
		return fmt.Errorf("mark Review Issue %q failed after Review Source propagation failure: %w", issue.Path, err)
	}
	return nil
}

func sourceActionForIssue(issue rounds.Issue, phase string) sourceIssueAction {
	switch issue.Status {
	case rounds.StatusResolved:
		if phase == "run_end" {
			return sourceIssueAction{}
		}
		return sourceIssueAction{Action: "resolved", Resolve: true}
	case rounds.StatusInvalid:
		if phase == "run_end" {
			return sourceIssueAction{}
		}
		return sourceIssueAction{Action: "invalid", Comment: true, Resolve: true}
	case rounds.StatusDuplicated:
		if phase == "run_end" {
			return sourceIssueAction{}
		}
		return sourceIssueAction{Action: "duplicated", Comment: true, Resolve: true}
	case rounds.StatusFailed:
		return sourceIssueAction{Action: "failed", Comment: true}
	default:
		if phase == "run_end" {
			return sourceIssueAction{Action: "unresolved", Comment: true}
		}
		return sourceIssueAction{}
	}
}

func (engine *Engine) reportSourcePropagationFailure(issue rounds.Issue, action string, err error) {
	if engine.deps.Progress == nil {
		return
	}
	display := strings.TrimSpace(issue.SourceRef)
	if display == "" {
		display = issue.Path
	}
	fmt.Fprintf(engine.deps.Progress, "Review Source propagation failed for %s (%s): %v\n", display, action, err)
}

func (engine *Engine) publishSourcePropagationEvent(ctx context.Context, plan CyclePlan, batchNumber int, issue rounds.Issue, phase string, action string, commentPosted bool, resolved bool, err error) error {
	payload := map[string]any{
		"phase":       phase,
		"issue_path":  issue.Path,
		"source_ref":  issue.SourceRef,
		"status":      issue.Status,
		"action":      action,
		"commented":   commentPosted,
		"resolved":    resolved,
		"batch":       batchNumber,
		"review_hash": issue.ReviewHash,
	}
	summary := fmt.Sprintf("Review Source propagation for %s: %s.", sourceIssueDisplay(issue), action)
	if err != nil {
		payload["error"] = err.Error()
		summary = fmt.Sprintf("Review Source propagation failed for %s: %s.", sourceIssueDisplay(issue), action)
	}
	return engine.publishDaemonEvent(ctx, plan.RunID, batchNumber, runevent.KindDaemonSourceResolution, summary, payload)
}

func sourceIssueDisplay(issue rounds.Issue) string {
	if strings.TrimSpace(issue.SourceRef) != "" {
		return issue.SourceRef
	}
	if strings.TrimSpace(issue.ReviewHash) != "" {
		return issue.ReviewHash
	}
	return issue.Path
}

func outcomeCommentMarker(runID string, issue rounds.Issue, action string) string {
	key := issue.SourceRef
	if strings.TrimSpace(key) == "" {
		key = issue.ReviewHash
	}
	if strings.TrimSpace(key) == "" {
		key = issue.Path
	}
	sum := sha256.Sum256([]byte(runID + "\x00" + key + "\x00" + action))
	return fmt.Sprintf("<!-- roundfix:outcome run=%s issue=%s action=%s -->", runID, hex.EncodeToString(sum[:8]), action)
}

func outcomeCommentBody(runID string, issue rounds.Issue, action string, marker string) string {
	lines := []string{
		marker,
		fmt.Sprintf("Roundfix outcome for Run %s: %s.", runID, action),
	}
	switch action {
	case "invalid":
		lines = append(lines,
			"Reason: "+publicOutcomeReason(issue.TerminalReason, "The Review Issue was triaged invalid."),
			"Next action: this thread is resolved after recording the triage outcome.",
		)
	case "duplicated":
		lines = append(lines,
			"Canonical Review Issue: "+publicDuplicateReference(issue),
			"Next action: this duplicate thread is resolved after recording the canonical issue.",
		)
	case "failed":
		lines = append(lines,
			"Reason: "+publicOutcomeReason(issue.TerminalReason, "The Batch failed before this Review Issue could settle."),
			"Next action: this thread stays open; a later Round retries it after the failed step is addressed.",
		)
	case "unresolved":
		lines = append(lines,
			"Reason: this Run ended with this Review Issue still unresolved.",
			"Next action: this thread stays open; a later Round retries it.",
		)
	}
	return strings.Join(lines, "\n\n")
}

func publicDuplicateReference(issue rounds.Issue) string {
	duplicateOf := strings.TrimSpace(issue.DuplicateOf)
	if duplicateOf == "" {
		return "another Review Issue in this Run"
	}
	if canonical, err := rounds.ParseIssue(duplicateOf); err == nil {
		if sourceRef := strings.TrimSpace(canonical.SourceRef); sourceRef != "" {
			return sourceRef
		}
		return "another Review Issue in this Run"
	}
	if isPublicSourceRef(duplicateOf) {
		return duplicateOf
	}
	return "another Review Issue in this Run"
}

func isPublicSourceRef(value string) bool {
	return strings.HasPrefix(value, "thread:") || strings.HasPrefix(value, "comment:")
}

func publicOutcomeReason(value string, fallback string) string {
	reason := strings.Join(strings.Fields(value), " ")
	if reason == "" {
		return fallback
	}
	fields := strings.Fields(reason)
	for index, field := range fields {
		fields[index] = publicOutcomeReasonToken(field)
	}
	reason = strings.Join(fields, " ")
	reason = strings.ReplaceAll(reason, "@", "＠")
	reason = strings.ReplaceAll(reason, "`", "'")
	if reason == "" {
		return fallback
	}
	runes := []rune(reason)
	if len(runes) > publicOutcomeReasonMaxRunes {
		reason = string(runes[:publicOutcomeReasonMaxRunes-1]) + "…"
	}
	return reason
}

func publicOutcomeReasonToken(token string) string {
	trimmed := strings.Trim(token, `"'(),.;:[]{}<>`)
	if looksLikeLocalPath(trimmed) {
		return "<path>"
	}
	return token
}

func looksLikeLocalPath(value string) bool {
	if value == "" || strings.Contains(value, "://") {
		return false
	}
	if strings.HasPrefix(value, "/") ||
		strings.HasPrefix(value, "~/") ||
		strings.HasPrefix(value, "./") ||
		strings.HasPrefix(value, "../") ||
		strings.HasPrefix(value, ".roundfix/") ||
		strings.Contains(value, "\\") {
		return true
	}
	if len(value) >= 3 && value[1] == ':' && (value[2] == '\\' || value[2] == '/') {
		return true
	}
	return strings.Contains(value, "/")
}

// publishDaemonEvent appends one daemon-owned Run Event. Publication is
// part of the Run state contract: a critical sink failure propagates and
// fails the cycle rather than being swallowed.
func (engine *Engine) publishDaemonEvent(ctx context.Context, runID string, batchNumber int, kind runevent.Kind, summary string, payload any) error {
	raw, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("encode daemon event payload: %w", err)
	}
	if err := engine.deps.Sink.Publish(ctx, runevent.RunEvent{
		RunID:   runID,
		Batch:   batchNumber,
		Source:  runevent.SourceDaemon,
		Kind:    kind,
		Summary: runevent.BoundSummary(summary),
		Time:    engine.deps.Now(),
		Payload: raw,
	}); err != nil {
		return fmt.Errorf("publish daemon event %s: %w", kind, err)
	}
	return nil
}

// publishStop records a Stop Request observed at a daemon boundary so the
// stop is visible in the event stream before the engine returns.
func (engine *Engine) publishStop(ctx context.Context, runID string, batchNumber int) error {
	payload, err := json.Marshal(struct {
		Status string `json:"status"`
	}{Status: "stopped"})
	if err != nil {
		return fmt.Errorf("encode stop event payload: %w", err)
	}
	if err := engine.deps.Sink.Publish(context.WithoutCancel(ctx), runevent.RunEvent{
		RunID:   runID,
		Batch:   batchNumber,
		Source:  runevent.SourceDaemon,
		Kind:    runevent.KindDaemonStatus,
		Summary: "Stop Request: cycle halted",
		Time:    engine.deps.Now(),
		Payload: payload,
	}); err != nil {
		return fmt.Errorf("publish stop event: %w", err)
	}
	return nil
}

func (engine *Engine) stopIfRequested(ctx context.Context, runID string, batchNumber int) error {
	requested, err := engine.deps.Runs.StopRequested(ctx, runID)
	if err != nil {
		return fmt.Errorf("read Stop Request flag for run %q: %w", runID, err)
	}
	if !requested {
		return nil
	}
	if err := engine.publishStop(ctx, runID, batchNumber); err != nil {
		return err
	}
	return ErrStopRequested
}

func (engine *Engine) reportPending(plan CyclePlan, failedIndex int) {
	pendingBatches := plan.Batches[failedIndex+1:]
	pending := 0
	for _, batch := range pendingBatches {
		pending += len(batch.Issues)
	}
	if pending > 0 {
		fmt.Fprintf(engine.deps.Progress, "Resolve stopped after Batch %03d/%03d failed; %d planned Review Issue(s) remain pending in %d Batch(es).\n", plan.Batches[failedIndex].Number, len(plan.Batches), pending, len(pendingBatches))
	}
}

func validateCyclePlan(plan CyclePlan) error {
	required := map[string]string{
		"Run ID":             plan.RunID,
		"Agent Session":      plan.Session.Name,
		"git root":           plan.GitRoot,
		"Artifact Directory": plan.ArtifactDir,
	}
	for label, value := range required {
		if strings.TrimSpace(value) == "" {
			return fmt.Errorf("resolve cycle: %s is required", label)
		}
	}
	if len(plan.Batches) == 0 {
		return errors.New("resolve cycle: at least one Batch is required")
	}
	return nil
}

func settledBatchSourceIssues(batch rounds.Batch, duplicated []rounds.Issue) ([]rounds.Issue, error) {
	issues := make([]rounds.Issue, 0, len(batch.Issues)+len(duplicated))
	for _, issue := range duplicated {
		if strings.TrimSpace(issue.SourceRef) == "" {
			continue
		}
		issues = append(issues, issue)
	}
	for _, assigned := range batch.Issues {
		issue, err := rounds.ParseIssue(assigned.Path)
		if err != nil {
			return nil, err
		}
		switch issue.Status {
		case rounds.StatusResolved, rounds.StatusInvalid, rounds.StatusDuplicated, rounds.StatusFailed:
		default:
			continue
		}
		if strings.TrimSpace(issue.SourceRef) == "" {
			continue
		}
		issues = append(issues, issue)
	}
	return issues, nil
}

func duplicatedIssuesForBatch(ctx context.Context, batch rounds.Batch, associations []rounds.DuplicateAssociation) ([]rounds.Issue, error) {
	if len(associations) == 0 {
		return nil, nil
	}
	assigned := make(map[string]bool, len(batch.Issues))
	for _, issue := range batch.Issues {
		assigned[issue.Path] = true
	}
	duplicated := []rounds.Issue{}
	for _, association := range associations {
		if err := ctx.Err(); err != nil {
			return duplicated, err
		}
		if !assigned[association.Newest.Path] {
			continue
		}
		issue, err := rounds.ParseIssue(association.Older.Path)
		if err != nil {
			return duplicated, err
		}
		if issue.Status != rounds.StatusDuplicated {
			continue
		}
		duplicated = append(duplicated, issue)
	}
	return duplicated, nil
}

func remainingUnresolvedIssues(ctx context.Context, plan CyclePlan) (int, error) {
	selection, err := rounds.SelectCompatibleIssues(ctx, rounds.SelectRequest{
		ArtifactDir:    plan.ArtifactDir,
		ReviewRoot:     plan.ReviewRoot,
		PRNumber:       plan.PullRequest.Number,
		HeadRepository: plan.PullRequest.HeadRepository,
		HeadBranch:     plan.PullRequest.HeadBranch,
	})
	if err != nil {
		var noArtifacts rounds.NoCompatibleArtifactsError
		if errors.As(err, &noArtifacts) {
			return 0, nil
		}
		return 0, err
	}
	return len(selection.Issues), nil
}

func unresolvedSourceIssues(ctx context.Context, plan CyclePlan) ([]rounds.Issue, error) {
	selection, err := rounds.SelectCompatibleIssues(ctx, rounds.SelectRequest{
		ArtifactDir:    plan.ArtifactDir,
		ReviewRoot:     plan.ReviewRoot,
		PRNumber:       plan.PullRequest.Number,
		HeadRepository: plan.PullRequest.HeadRepository,
		HeadBranch:     plan.PullRequest.HeadBranch,
	})
	if err != nil {
		var noArtifacts rounds.NoCompatibleArtifactsError
		if errors.As(err, &noArtifacts) {
			return nil, nil
		}
		return nil, err
	}
	issues := make([]rounds.Issue, 0, len(selection.Issues))
	for _, selected := range selection.Issues {
		if err := ctx.Err(); err != nil {
			return issues, err
		}
		issue, err := rounds.ParseIssue(selected.Path)
		if err != nil {
			return nil, err
		}
		if strings.TrimSpace(issue.SourceRef) == "" {
			continue
		}
		issues = append(issues, issue)
	}
	return issues, nil
}

func isStop(ctx context.Context, err error) bool {
	return agent.IsStopError(err) || errors.Is(err, ErrStopRequested) || errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) || (ctx != nil && ctx.Err() != nil)
}
