package daemon

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"sort"
	"strings"
	"time"

	"roundfix/internal/agent"
	"roundfix/internal/reviewsource"
	"roundfix/internal/rounds"
	"roundfix/internal/runevent"
	"roundfix/internal/store"
)

// ReviewSourceResolver resolves Review Source threads for terminal Review
// Issues after verification.
type ReviewSourceResolver interface {
	ResolveIssues(ctx context.Context, req reviewsource.ResolveRequest) error
}

// ReviewSourceResolverFunc adapts a function to ReviewSourceResolver.
type ReviewSourceResolverFunc func(ctx context.Context, req reviewsource.ResolveRequest) error

func (fn ReviewSourceResolverFunc) ResolveIssues(ctx context.Context, req reviewsource.ResolveRequest) error {
	return fn(ctx, req)
}

// RunStateStore records intermediate Run states during a cycle. Terminal
// completion stays with the caller.
type RunStateStore interface {
	UpdateRunState(ctx context.Context, runID string, state string) error
}

// Dependencies are the engine's explicit collaborators, replacing the CLI
// package globals that previously wired orchestration.
type Dependencies struct {
	Runner    agent.Runner
	Verifier  Verifier
	Committer Committer
	Pusher    Pusher
	Source    ReviewSourceResolver
	Runs      RunStateStore
	Worktree  WorktreeSnapshotter
	Sink      runevent.Sink
	Now       func() time.Time
	Progress  io.Writer
}

// Engine executes one resolve cycle over a validated plan and exposes Final
// Push as a separate explicit operation, so resolve and watch share one
// orchestration implementation. The Daemon owns Batch commits and Final
// Push per ADR 0001.
type Engine struct {
	deps Dependencies
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
	GitRoot      string
	ArtifactDir  string
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
type BatchOutcome struct {
	Batch                 int
	Issues                int
	Committed             bool
	CommitSkipped         bool
	ResolvedSourceThreads int
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
	if deps.Now == nil {
		deps.Now = time.Now
	}
	if deps.Progress == nil {
		deps.Progress = io.Discard
	}
	return &Engine{deps: deps}, nil
}

// ResolveCycle executes one resolve cycle: for each Batch it runs the
// Agent, verifies assigned issue statuses, creates the Batch commit when
// auto-commit is enabled, and resolves Review Source threads for resolved
// and invalid issues. A Stop Request halts before any new Batch,
// verification, commit, or Review Source mutation; Agent worktree changes
// are preserved.
func (engine *Engine) ResolveCycle(ctx context.Context, plan CyclePlan) (CycleResult, error) {
	if err := validateCyclePlan(plan); err != nil {
		return CycleResult{}, err
	}
	result := CycleResult{Remaining: plan.TotalIssues}
	for index, batch := range plan.Batches {
		if err := ctx.Err(); err != nil {
			engine.publishStop(ctx, plan.RunID, batch.Number)
			return result, err
		}
		outcome, remaining, err := engine.resolveBatch(ctx, plan, batch, index+1, len(plan.Batches))
		if err != nil {
			engine.reportPending(plan, index)
			return result, err
		}
		result.Batches = append(result.Batches, outcome)
		result.Remaining = remaining
		if remaining > 0 && index < len(plan.Batches)-1 {
			fmt.Fprintf(engine.deps.Progress, "Batch %03d/%03d completed; %d Unresolved Review Issue(s) remain.\n", batch.Number, len(plan.Batches), remaining)
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
	return engine.deps.Pusher.Push(ctx, PushRequest{
		WorkDir: req.WorkDir,
		Remote:  req.Remote,
		Branch:  req.Branch,
	})
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
	if err := engine.runBatchAgent(ctx, plan, batch, batchIndex, batchTotal); err != nil {
		return outcome, 0, err
	}
	if err := engine.verifyBatch(ctx, plan, batch); err != nil {
		return outcome, 0, err
	}
	marked, err := rounds.MarkDuplicatedAfterTerminal(ctx, plan.Duplicates)
	if err != nil {
		return outcome, 0, err
	}
	if marked > 0 {
		fmt.Fprintf(engine.deps.Progress, "Marked %d older duplicate Review Issue occurrence(s) as duplicated.\n", marked)
	}
	committed, skipped, err := engine.commitBatch(ctx, plan, batch, before)
	if err != nil {
		return outcome, 0, err
	}
	outcome.Committed = committed
	outcome.CommitSkipped = skipped
	resolved, err := engine.resolveBatchSources(ctx, plan, batch)
	if err != nil {
		return outcome, 0, err
	}
	outcome.ResolvedSourceThreads = resolved
	remaining, err := remainingUnresolvedIssues(ctx, plan)
	if err != nil {
		return outcome, 0, err
	}
	return outcome, remaining, nil
}

func (engine *Engine) runBatchAgent(ctx context.Context, plan CyclePlan, batch rounds.Batch, batchIndex int, batchTotal int) error {
	if err := engine.deps.Runs.UpdateRunState(ctx, plan.RunID, store.StateResolvingWithAgent); err != nil {
		return err
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
	logPath := agent.LogPath(plan.ArtifactDir, plan.RunID, batch.Number)
	fmt.Fprintf(engine.deps.Progress, "Batch: %03d/%03d (%d Review Issue(s))\n", batchIndex, batchTotal, len(batch.Issues))
	fmt.Fprintf(engine.deps.Progress, "Agent log: %s\n", logPath)

	_, runErr := engine.deps.Runner.Run(ctx, agent.ExecuteRequest{
		Runtime:      plan.Runtime,
		RunID:        plan.RunID,
		Batch:        batch,
		LogPath:      logPath,
		Prompt:       prompt,
		ArtifactDir:  plan.ArtifactDir,
		GitRoot:      plan.GitRoot,
		Verification: plan.Verification,
		AllowAddDirs: []string{plan.ArtifactDir},
	}, engine.deps.Sink)
	if runErr != nil {
		if isStop(ctx, runErr) {
			// The runner already published the stopped status event;
			// Agent-created worktree changes stay untouched.
			return runErr
		}
		_ = agent.MarkBatchFailed(batch)
		return runErr
	}
	if err := ctx.Err(); err != nil {
		engine.publishStop(ctx, plan.RunID, batch.Number)
		return err
	}
	if err := agent.ValidateAssignedIssuesTerminal(batch); err != nil {
		_ = agent.MarkBatchFailed(batch)
		return err
	}
	fmt.Fprintln(engine.deps.Progress, "Agent Batch reached terminal local Review Issue statuses.")
	return nil
}

func (engine *Engine) verifyBatch(ctx context.Context, plan CyclePlan, batch rounds.Batch) error {
	if err := engine.deps.Runs.UpdateRunState(ctx, plan.RunID, store.StateVerifying); err != nil {
		return err
	}
	if err := engine.deps.Verifier.Verify(ctx, VerifyRequest{
		WorkDir: plan.GitRoot,
		Command: plan.Verification,
		Stream:  engine.deps.Progress,
	}); err != nil {
		_ = agent.MarkBatchFailed(batch)
		return err
	}
	fmt.Fprintf(engine.deps.Progress, "Verification command passed: %s\n", plan.Verification)
	return nil
}

func (engine *Engine) commitBatch(ctx context.Context, plan CyclePlan, batch rounds.Batch, before []string) (bool, bool, error) {
	if !plan.AutoCommit {
		fmt.Fprintln(engine.deps.Progress, "Auto-commit disabled; no Batch commit created.")
		return false, false, nil
	}
	if err := ctx.Err(); err != nil {
		engine.publishStop(ctx, plan.RunID, batch.Number)
		return false, false, err
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
		return false, true, nil
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

func (engine *Engine) resolveBatchSources(ctx context.Context, plan CyclePlan, batch rounds.Batch) (int, error) {
	issues, err := terminalAssignedSourceIssues(batch)
	if err != nil {
		return 0, err
	}
	if len(issues) == 0 {
		return 0, nil
	}
	if err := ctx.Err(); err != nil {
		engine.publishStop(ctx, plan.RunID, batch.Number)
		return 0, err
	}
	if err := engine.deps.Source.ResolveIssues(ctx, reviewsource.ResolveRequest{
		Source:         plan.SourceName,
		PRNumber:       plan.PullRequest.Number,
		BaseRepository: plan.PullRequest.BaseRepository,
		Issues:         issues,
	}); err != nil {
		return 0, err
	}
	fmt.Fprintf(engine.deps.Progress, "Resolved %d Review Source thread(s).\n", len(issues))
	return len(issues), nil
}

// publishStop records a Stop Request observed at a daemon boundary so the
// stop is visible in the event stream before the engine returns.
func (engine *Engine) publishStop(ctx context.Context, runID string, batchNumber int) {
	payload, err := json.Marshal(struct {
		Status string `json:"status"`
	}{Status: "stopped"})
	if err != nil {
		return
	}
	_ = engine.deps.Sink.Publish(context.WithoutCancel(ctx), runevent.RunEvent{
		RunID:   runID,
		Batch:   batchNumber,
		Source:  runevent.SourceDaemon,
		Kind:    runevent.KindDaemonStatus,
		Summary: "Stop Request: cycle halted",
		Time:    engine.deps.Now(),
		Payload: payload,
	})
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

func terminalAssignedSourceIssues(batch rounds.Batch) ([]reviewsource.ResolvedIssue, error) {
	issues := make([]reviewsource.ResolvedIssue, 0, len(batch.Issues))
	for _, assigned := range batch.Issues {
		issue, err := rounds.ParseIssue(assigned.Path)
		if err != nil {
			return nil, err
		}
		if issue.Status != rounds.StatusResolved && issue.Status != rounds.StatusInvalid {
			continue
		}
		if strings.TrimSpace(issue.SourceRef) == "" {
			continue
		}
		issues = append(issues, reviewsource.ResolvedIssue{
			FilePath:  issue.Path,
			Status:    issue.Status,
			SourceRef: issue.SourceRef,
		})
	}
	return issues, nil
}

func remainingUnresolvedIssues(ctx context.Context, plan CyclePlan) (int, error) {
	selection, err := rounds.SelectCompatibleIssues(ctx, rounds.SelectRequest{
		ArtifactDir:    plan.ArtifactDir,
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

func isStop(ctx context.Context, err error) bool {
	return agent.IsStopError(err) || errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) || (ctx != nil && ctx.Err() != nil)
}
