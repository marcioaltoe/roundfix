// Suite: current Run disposition characterization.
// Invariant: each later disposition change starts from an executable record of today's outcome.
// Boundary IN: Agent Session ownership, Batch settlement, Task-cycle Git state, and the public CLI preflight.
// Boundary OUT: the replacement behavior owned by tasks 02 through 06.
package daemon

import (
	"bytes"
	"context"
	"errors"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"

	"roundfix/internal/agent"
	"roundfix/internal/gittest"
	"roundfix/internal/rounds"
	"roundfix/internal/runevent"
	"roundfix/internal/spec"
	"roundfix/internal/store"
	runworktree "roundfix/internal/worktree"
)

// Outcome contract test: internal/cli/cli_test.go::TestRunResolveVerificationFailureDoesNotCommit asserts that a failed Verification produces exit 1, an Unresolved Run report, and a failed Review Issue without a commit, source resolution, or Final Push; task_04 rewrites it.
// Outcome contract test: internal/cli/cli_test.go::TestRunResolveAgentFailureMarksBatchFailed asserts that an Agent crash produces exit 1, an Unresolved Run report, and a failed Review Issue; task_04 rewrites it.
// Outcome contract test: internal/cli/cli_test.go::TestRunResolveAgentFailureContinuesWithLaterBatches asserts that every crashed Batch leaves failed Review Issues, later Batches still run, and the Run exits 1 as Unresolved; task_04 rewrites it.
// Outcome contract test: internal/cli/cli_test.go::TestRunResolveClosesAgentSessionForTerminalOutcomes/unresolved asserts that an Agent crash closes the Agent Session and leaves the Run in store.StateUnresolved with exit 1; task_04 rewrites it.
// Outcome contract test: internal/daemon/engine_test.go::TestResolveCycleVerificationFailureFailsBatchAndContinues asserts that final Verification failure leaves the Review Issue failed and unresolved, with a failed Batch and no commit; task_04 rewrites it.
// Outcome contract test: internal/daemon/engine_test.go::TestResolveCycleContinuesToNextBatchAfterFailedBatch asserts that a failed first Batch leaves one failed Review Issue while the second resolves, so one issue remains unresolved; task_04 rewrites it.

func TestRunDispositionCharacterizationWorkStartedFollowsTheFirstAgentOutput(t *testing.T) {
	t.Parallel()
	sink := &captureEventSink{}
	promptErr := errors.New("adapter refused the first prompt without Agent output")
	runner := &dispositionPromptRunner{sink: sink, err: promptErr}
	engine := &Engine{deps: Dependencies{
		Runner:   runner,
		Sink:     sink,
		Progress: &bytes.Buffer{},
		Now: func() time.Time {
			return time.Date(2026, 8, 9, 12, 0, 0, 0, time.UTC)
		},
	}}
	owner := &agentSessionOwner{
		engine:        engine,
		scope:         agentSessionScope{RunID: "run-characterization", Kind: "task", ID: "task_01", Batch: 1},
		activeRuntime: agent.RuntimeSpec{ID: "codex", DisplayName: "Codex"},
		activeSession: agent.SessionRef{Name: "task-session", WorkDir: t.TempDir()},
		active:        true,
		workStarted:   false,
		attemptNumber: 1,
	}

	_, err := owner.Run(context.Background(), agent.ExecuteRequest{
		RunID: "run-characterization",
		Batch: rounds.Batch{Number: 1},
	})

	if !errors.Is(err, promptErr) {
		t.Fatalf("first prompt error = %v, want %v", err, promptErr)
	}
	if runner.workStartedBeforeOutput {
		t.Fatal("Agent work-started status was published before the first Agent output")
	}
	if !runner.workStartedAfterOutput {
		t.Fatal("Agent work-started status was not published with the first Agent output")
	}
	events := sink.snapshot()
	if len(events) != 2 || events[0].Kind != runevent.KindAgentStatus ||
		!strings.Contains(string(events[0].Payload), agent.AgentWorkStartedStatus) ||
		events[1].Kind != runevent.KindAgentMessage {
		t.Fatalf("events = %+v, want work-started immediately before the first Agent message", events)
	}
}

func TestRunDispositionCharacterizationFailedBatchOverwritesSettledIssues(t *testing.T) {
	t.Parallel()
	fixture := newEngineFixture(t)
	const failureReason = "Verification failed after the Agent settled this issue"
	if err := rounds.SetIssueStatus(fixture.issuePaths[0], rounds.StatusResolved, "", ""); err != nil {
		t.Fatalf("settle Review Issue before Batch failure: %v", err)
	}

	err := agent.MarkBatchFailed(rounds.Batch{
		Number: 1,
		Issues: []rounds.Issue{{Path: fixture.issuePaths[0]}},
	}, failureReason)

	if err != nil {
		t.Fatalf("MarkBatchFailed() error = %v", err)
	}
	issue, err := rounds.ParseIssue(fixture.issuePaths[0])
	if err != nil {
		t.Fatalf("parse Review Issue after Batch failure: %v", err)
	}
	if issue.Status != rounds.StatusFailed {
		t.Fatalf("settled Review Issue status after Batch failure = %q, want %q", issue.Status, rounds.StatusFailed)
	}
	if issue.TerminalReason != failureReason {
		t.Fatalf("settled Review Issue reason after Batch failure = %q, want %q", issue.TerminalReason, failureReason)
	}
}

func TestRunDispositionCharacterizationStoppedRunLeavesTasksPending(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	checkoutRoot := t.TempDir()
	gittest.InitRepo(t, checkoutRoot, "-b", "ma/spec-work")
	gittest.PersistIdentity(t, checkoutRoot)
	seeds := []taskSpecSeed{
		{id: "task_01", title: "Build the core"},
		{id: "task_02", title: "Wire the shell"},
	}
	writeSpecDirForTest(t, checkoutRoot, taskCycleSlug, seeds)
	gittest.Run(t, checkoutRoot, "add", ".")
	gittest.Run(t, checkoutRoot, "commit", "-m", "test: seed pending tasks")
	checkoutHead := strings.TrimSpace(gittest.Run(t, checkoutRoot, "rev-parse", "HEAD"))

	homeDir := t.TempDir()
	runStore, err := store.Open(ctx, homeDir)
	if err != nil {
		t.Fatalf("open Run Database: %v", err)
	}
	t.Cleanup(func() { _ = runStore.Close() })
	run, err := runStore.CreateRun(ctx, store.CreateRunRequest{
		Kind:        store.KindImplement,
		GitRoot:     checkoutRoot,
		LocalBranch: "ma/spec-work",
		HeadSHA:     checkoutHead,
		SpecSlug:    taskCycleSlug,
	})
	if err != nil {
		t.Fatalf("create implement Run: %v", err)
	}
	ref, err := runworktree.Create(ctx, runworktree.CreateOptions{
		UserRoot: checkoutRoot,
		Location: t.TempDir(),
		RunID:    run.ID,
		HeadSHA:  checkoutHead,
	})
	if err != nil {
		t.Fatalf("create Run Worktree: %v", err)
	}
	if _, err := runStore.SetRunWorkDir(ctx, run.ID, ref.Path); err != nil {
		t.Fatalf("record Run Worktree: %v", err)
	}

	specsRoot := filepath.Join(ref.Path, "docs", "specs")
	graph, err := spec.Load(specsRoot, taskCycleSlug)
	if err != nil {
		t.Fatalf("load Spec from Run Worktree: %v", err)
	}
	calls := []string{}
	fixture := &taskCycleFixture{
		t:           t,
		seeds:       seeds,
		store:       runStore,
		run:         run,
		gitRoot:     ref.Path,
		specsRoot:   specsRoot,
		artifactDir: filepath.Join(t.TempDir(), "artifacts"),
		graph:       graph,
		calls:       &calls,
		sink:        &captureEventSink{},
		progress:    &bytes.Buffer{},
		worktree:    &engineFakeWorktree{snapshots: [][]string{nil, {"src/core.go"}}},
		pusher:      &engineFakePusher{calls: &calls},
		source:      &engineFakeSource{calls: &calls},
		github:      &taskFakeGHRunner{output: "[]"},
	}
	runner := &taskFakeRunner{
		calls:       fixture.calls,
		gitRoot:     fixture.gitRoot,
		writeByTask: map[string]string{"task_01": "src/core.go"},
		statusByTask: map[string]spec.Status{
			"task_01": spec.StatusCompleted,
			"task_02": spec.StatusCompleted,
		},
	}
	committer := &dispositionStoppingCommitter{
		delegate: GitCommitter{},
		store:    runStore,
		runID:    run.ID,
	}
	engine := fixture.engine(t, runner, &taskFakeVerifier{calls: fixture.calls}, committer, fixture.worktree)
	plan := TaskPlan{
		RunID:                   run.ID,
		Session:                 agent.SessionRefForRun(run.ID, ref.Path),
		WorkDir:                 ref.Path,
		RunWorktree:             ref,
		TargetBranch:            "ma/spec-work",
		Spec:                    graph.Spec,
		SpecsRoot:               specsRoot,
		Tasks:                   graph.Tasks,
		Runtime:                 agent.RuntimeSpec{ID: "codex", DisplayName: "Codex"},
		ArtifactDir:             fixture.artifactDir,
		Concurrency:             1,
		VerificationConcurrency: 1,
	}

	result, cycleErr := engine.TaskCycle(ctx, plan)

	if !errors.Is(cycleErr, ErrStopRequested) {
		t.Fatalf("TaskCycle() error = %v, want ErrStopRequested", cycleErr)
	}
	if result.Completed != 1 || result.Failed != 0 || result.Skipped != 0 {
		t.Fatalf("settlement before stop = %+v, want one completed Task", result)
	}
	if _, err := runStore.CompleteRun(ctx, run.ID, store.StateStopped); err != nil {
		t.Fatalf("complete stopped Run: %v", err)
	}
	if got := taskStatusOnDisk(t, checkoutRoot, "task_01"); got != string(spec.StatusPending) {
		t.Fatalf("checkout task_01 status = %q, want %q", got, spec.StatusPending)
	}
	if got := taskStatusOnDisk(t, ref.Path, "task_01"); got != string(spec.StatusCompleted) {
		t.Fatalf("Run Worktree task_01 status = %q, want %q", got, spec.StatusCompleted)
	}
	runHead := strings.TrimSpace(gittest.Run(t, ref.Path, "rev-parse", "HEAD"))
	if runHead == checkoutHead {
		t.Fatalf("Run Worktree HEAD = checkout HEAD %s, want a settlement commit", runHead)
	}
	committedPaths := gittest.Run(t, ref.Path, "show", "--format=", "--name-only", "HEAD")
	if !strings.Contains(committedPaths, filepath.ToSlash(taskFileRel(taskCycleSlug, "task_01"))) {
		t.Fatalf("Run Worktree settlement commit paths = %q, want task_01", committedPaths)
	}
}

func TestRunDispositionCharacterizationPreflightRefusesOnAnUnintegratedBranch(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	homeDir := t.TempDir()
	checkoutRoot := t.TempDir()
	gittest.InitRepo(t, checkoutRoot, "-b", "feature/review")
	gittest.PersistIdentity(t, checkoutRoot)
	mustWriteForTest(t, filepath.Join(checkoutRoot, "base.txt"), "base\n")
	mustWriteForTest(t, filepath.Join(checkoutRoot, ".roundfixrc.yml"), "review_source:\n  request_review: true\nwatch:\n  auto_push: false\n")
	gittest.Run(t, checkoutRoot, "add", ".")
	gittest.Run(t, checkoutRoot, "commit", "-m", "test: seed review branch")
	baseHead := strings.TrimSpace(gittest.Run(t, checkoutRoot, "rev-parse", "HEAD"))

	runStore, err := store.Open(ctx, homeDir)
	if err != nil {
		t.Fatalf("open Run Database: %v", err)
	}
	run, err := runStore.CreateRun(ctx, store.CreateRunRequest{
		Kind:           store.KindImplement,
		HeadRepository: "owner/project",
		HeadBranch:     "feature/review",
		GitRoot:        checkoutRoot,
		LocalBranch:    "feature/review",
		HeadSHA:        baseHead,
		SpecSlug:       "characterization-spec",
	})
	if err != nil {
		_ = runStore.Close()
		t.Fatalf("create stopped Run: %v", err)
	}
	ref, err := runworktree.Create(ctx, runworktree.CreateOptions{
		UserRoot: checkoutRoot,
		Location: t.TempDir(),
		RunID:    run.ID,
		HeadSHA:  baseHead,
	})
	if err != nil {
		_ = runStore.Close()
		t.Fatalf("create unintegrated Run Worktree: %v", err)
	}
	if _, err := runStore.SetRunWorkDir(ctx, run.ID, ref.Path); err != nil {
		_ = runStore.Close()
		t.Fatalf("record unintegrated Run Worktree: %v", err)
	}
	mustWriteForTest(t, filepath.Join(ref.Path, "run-only.txt"), "unintegrated work\n")
	gittest.Run(t, ref.Path, "add", "run-only.txt")
	gittest.Run(t, ref.Path, "commit", "-m", "test: retain run-only work")
	mustWriteForTest(t, filepath.Join(checkoutRoot, "target-only.txt"), "newer target work\n")
	gittest.Run(t, checkoutRoot, "add", "target-only.txt")
	gittest.Run(t, checkoutRoot, "commit", "-m", "test: diverge the target branch")
	if _, err := runStore.CompleteRun(ctx, run.ID, store.StateStopped); err != nil {
		_ = runStore.Close()
		t.Fatalf("complete stopped Run: %v", err)
	}
	if err := runStore.Close(); err != nil {
		t.Fatalf("close Run Database before CLI invocation: %v", err)
	}

	binary := filepath.Join(t.TempDir(), "roundfix")
	_, thisFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("resolve characterization test source path")
	}
	moduleRoot := filepath.Clean(filepath.Join(filepath.Dir(thisFile), "..", ".."))
	build := exec.Command("go", "build", "-buildvcs=false", "-o", binary, "./cmd/roundfix")
	build.Dir = moduleRoot
	buildOutput, err := build.CombinedOutput()
	if err != nil {
		t.Fatalf("build Roundfix CLI: %v\n%s", err, buildOutput)
	}
	command := exec.Command(binary,
		"resolve",
		"--pr", "123",
		"--head-repo", "owner/project",
		"--head-branch", "feature/review",
		"--no-input",
	)
	command.Dir = checkoutRoot
	command.Env = dispositionCommandEnvironment(homeDir)
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	command.Stdout = &stdout
	command.Stderr = &stderr
	err = command.Run()
	var exitErr *exec.ExitError
	if !errors.As(err, &exitErr) || exitErr.ExitCode() != 2 {
		t.Fatalf("resolve exit error = %v, stdout=%q stderr=%q; want Branch Integrity exit 2", err, stdout.String(), stderr.String())
	}
	if stdout.Len() != 0 {
		t.Fatalf("Branch Integrity refusal stdout = %q, want empty", stdout.String())
	}
	for _, want := range []string{
		"Branch Integrity Preflight refused pending Run Branch work",
		ref.Branch,
		"git merge --ff-only " + ref.Branch,
		"did not create a Run",
	} {
		if !strings.Contains(stderr.String(), want) {
			t.Fatalf("Branch Integrity refusal stderr does not contain %q: %q", want, stderr.String())
		}
	}

	runStore, err = store.Open(ctx, homeDir)
	if err != nil {
		t.Fatalf("reopen Run Database: %v", err)
	}
	defer runStore.Close()
	runs, err := runStore.ListRuns(ctx, store.ListRunsQuery{GitRoot: checkoutRoot, States: store.StatesAll})
	if err != nil {
		t.Fatalf("list Runs after Branch Integrity refusal: %v", err)
	}
	if len(runs) != 1 || runs[0].ID != run.ID {
		t.Fatalf("Runs after Branch Integrity refusal = %+v, want only stopped Run %s", runs, run.ID)
	}
}

type dispositionPromptRunner struct {
	sink                    *captureEventSink
	err                     error
	workStartedBeforeOutput bool
	workStartedAfterOutput  bool
}

func (*dispositionPromptRunner) Probe(context.Context, agent.ProbeRequest) error { return nil }

func (runner *dispositionPromptRunner) Run(ctx context.Context, req agent.ExecuteRequest, sink runevent.Sink) (agent.ExecuteResult, error) {
	return runner.RunPrepared(ctx, req, sink)
}

func (runner *dispositionPromptRunner) RunPrepared(ctx context.Context, req agent.ExecuteRequest, sink runevent.Sink) (agent.ExecuteResult, error) {
	for _, event := range runner.sink.snapshot() {
		if event.Kind == runevent.KindAgentStatus && strings.Contains(string(event.Payload), agent.AgentWorkStartedStatus) {
			runner.workStartedBeforeOutput = true
		}
	}
	if err := sink.Publish(ctx, runevent.RunEvent{
		RunID:   req.RunID,
		Batch:   req.Batch.Number,
		Source:  runevent.SourceAgent,
		Kind:    runevent.KindAgentMessage,
		Summary: "first Agent output",
	}); err != nil {
		return agent.ExecuteResult{LogPath: req.LogPath}, err
	}
	for _, event := range runner.sink.snapshot() {
		if event.Kind == runevent.KindAgentStatus && strings.Contains(string(event.Payload), agent.AgentWorkStartedStatus) {
			runner.workStartedAfterOutput = true
		}
	}
	return agent.ExecuteResult{LogPath: req.LogPath}, runner.err
}

func (*dispositionPromptRunner) EndSession(context.Context, agent.RuntimeSpec, agent.SessionRef) error {
	return nil
}

type dispositionStoppingCommitter struct {
	delegate Committer
	store    *store.Store
	runID    string
}

func (committer *dispositionStoppingCommitter) Commit(ctx context.Context, req CommitRequest) error {
	if err := committer.delegate.Commit(ctx, req); err != nil {
		return err
	}
	return committer.store.RequestStop(ctx, committer.runID)
}

func dispositionCommandEnvironment(homeDir string) []string {
	environment := gittest.IsolatedEnv()
	filtered := environment[:0]
	for _, entry := range environment {
		key, _, _ := strings.Cut(entry, "=")
		if key == "HOME" {
			continue
		}
		filtered = append(filtered, entry)
	}
	return append(filtered, "HOME="+homeDir)
}
