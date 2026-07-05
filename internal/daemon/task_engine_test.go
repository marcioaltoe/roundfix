package daemon

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"

	"roundfix/internal/agent"
	"roundfix/internal/runevent"
	"roundfix/internal/spec"
	"roundfix/internal/store"
)

const taskCycleSlug = "0001-sample-feature"

// taskSpecSeed describes one task file for a test Spec directory. Zero
// values default to status pending, type backend, and one passing
// verification command.
type taskSpecSeed struct {
	id           string
	title        string
	taskType     string
	status       string
	needs        []string
	verification []string
}

type taskCycleFixture struct {
	store    *store.Store
	run      store.Run
	gitRoot  string
	graph    *spec.Graph
	calls    *[]string
	sink     *captureEventSink
	progress *bytes.Buffer
	worktree *engineFakeWorktree
	pusher   *engineFakePusher
	source   *engineFakeSource
}

func newTaskCycleFixture(t *testing.T, seeds []taskSpecSeed) *taskCycleFixture {
	t.Helper()
	ctx := context.Background()
	homeDir := t.TempDir()
	gitRoot := t.TempDir()
	writeSpecDirForTest(t, gitRoot, taskCycleSlug, seeds)

	runStore, err := store.Open(ctx, homeDir)
	if err != nil {
		t.Fatalf("open store: %v", err)
	}
	t.Cleanup(func() { _ = runStore.Close() })

	run, err := runStore.CreateRun(ctx, store.CreateRunRequest{
		Kind:        store.KindImplement,
		GitRoot:     gitRoot,
		LocalBranch: "ma/spec-work",
		SpecSlug:    taskCycleSlug,
	})
	if err != nil {
		t.Fatalf("create run: %v", err)
	}

	graph, err := spec.Load(gitRoot, taskCycleSlug)
	if err != nil {
		t.Fatalf("load spec: %v", err)
	}

	calls := []string{}
	return &taskCycleFixture{
		store:    runStore,
		run:      run,
		gitRoot:  gitRoot,
		graph:    graph,
		calls:    &calls,
		sink:     &captureEventSink{},
		progress: &bytes.Buffer{},
		worktree: &engineFakeWorktree{},
		pusher:   &engineFakePusher{calls: &calls},
		source:   &engineFakeSource{calls: &calls},
	}
}

func (fixture *taskCycleFixture) plan() TaskPlan {
	return TaskPlan{
		RunID:   fixture.run.ID,
		Session: agent.SessionRefForRun(fixture.run.ID, fixture.gitRoot),
		WorkDir: fixture.gitRoot,
		Spec:    fixture.graph.Spec,
		Tasks:   fixture.graph.Tasks,
		Runtime: agent.RuntimeSpec{ID: "codex", DisplayName: "Codex"},
	}
}

func (fixture *taskCycleFixture) qaPlan() TaskPlan {
	plan := fixture.plan()
	plan.QA = true
	return plan
}

func (fixture *taskCycleFixture) engine(t *testing.T, runner agent.Runner, verifier Verifier, committer Committer, worktree WorktreeSnapshotter) *Engine {
	t.Helper()
	engine, err := NewEngine(Dependencies{
		Runner:    runner,
		Verifier:  verifier,
		Committer: committer,
		Pusher:    fixture.pusher,
		Source:    fixture.source,
		Runs:      fixture.store,
		Worktree:  worktree,
		Sink:      fixture.sink,
		Progress:  fixture.progress,
	})
	if err != nil {
		t.Fatalf("new engine: %v", err)
	}
	return engine
}

func writeSpecDirForTest(t *testing.T, gitRoot string, slug string, seeds []taskSpecSeed) {
	t.Helper()
	specDir := filepath.Join(gitRoot, "docs", "specs", slug)
	if err := os.MkdirAll(specDir, 0o755); err != nil {
		t.Fatalf("create spec dir: %v", err)
	}
	mustWriteForTest(t, filepath.Join(specDir, "_prd.md"), "---\nstatus: active\n---\n\n# PRD\n")

	var manifest strings.Builder
	manifest.WriteString("---\nschema: spec-tasks/v1\ngraph:\n  nodes:\n")
	for _, seed := range seeds {
		manifest.WriteString(fmt.Sprintf("    - id: %s\n      file: %s.md\n", seed.id, seed.id))
		if len(seed.needs) > 0 {
			manifest.WriteString(fmt.Sprintf("      needs: [%s]\n", strings.Join(seed.needs, ", ")))
		}
	}
	manifest.WriteString("---\n\n# Task Graph\n")
	mustWriteForTest(t, filepath.Join(specDir, "_tasks.md"), manifest.String())

	for _, seed := range seeds {
		status := seed.status
		if status == "" {
			status = string(spec.StatusPending)
		}
		taskType := seed.taskType
		if taskType == "" {
			taskType = "backend"
		}
		title := seed.title
		if title == "" {
			title = "Do the " + seed.id + " work"
		}
		verification := seed.verification
		if len(verification) == 0 {
			verification = []string{"true"}
		}
		var body strings.Builder
		body.WriteString(fmt.Sprintf("---\ntask: %s\nspec: %s\nstatus: %s\ntype: %s\n---\n\n# %s\n\n## Verification\n\n", seed.id, slug, status, taskType, title))
		for _, command := range verification {
			body.WriteString(fmt.Sprintf("- `%s` — expected: passes.\n", command))
		}
		mustWriteForTest(t, filepath.Join(specDir, seed.id+".md"), body.String())
	}
}

func taskPathFor(gitRoot string, slug string, id string) string {
	return filepath.Join(gitRoot, "docs", "specs", slug, id+".md")
}

func taskFileRel(slug string, id string) string {
	return filepath.Join("docs", "specs", slug, id+".md")
}

func taskStatusOnDisk(t *testing.T, gitRoot string, id string) string {
	t.Helper()
	content, err := os.ReadFile(taskPathFor(gitRoot, taskCycleSlug, id))
	if err != nil {
		t.Fatalf("read task file %s: %v", id, err)
	}
	for _, line := range strings.Split(string(content), "\n") {
		if strings.HasPrefix(line, "status:") {
			return strings.TrimSpace(strings.TrimPrefix(line, "status:"))
		}
	}
	t.Fatalf("task file %s has no status line", id)
	return ""
}

// taskFakeRunner scripts per-Task Agent behavior keyed by the Task id it
// parses from the prompt, mirroring how a real Agent learns its assignment
// from the prompt contract. A QA prompt writes qaReport as the Spec's QA
// Report, the way the qa-gate Agent does; an empty qaReport writes none.
type taskFakeRunner struct {
	calls        *[]string
	gitRoot      string
	store        *store.Store
	statusByTask map[string]spec.Status
	errByTask    map[string]error
	writeByTask  map[string]string
	qaReport     string
	qaPrompts    []string
	seenStates   []string
	prompts      []string
	requests     []agent.ExecuteRequest
}

func (runner *taskFakeRunner) Probe(context.Context, agent.RuntimeSpec) error { return nil }

func (runner *taskFakeRunner) Run(ctx context.Context, req agent.ExecuteRequest, sink runevent.Sink) (agent.ExecuteResult, error) {
	*runner.calls = append(*runner.calls, "agent")
	runner.prompts = append(runner.prompts, req.Prompt)
	runner.requests = append(runner.requests, req)
	if runner.store != nil {
		runner.seenStates = append(runner.seenStates, runStateForTest(runner.store, req.RunID))
	}
	taskID := taskIDFromPrompt(req.Prompt)
	if taskID == "" && strings.Contains(req.Prompt, "Spec QA gate") {
		runner.qaPrompts = append(runner.qaPrompts, req.Prompt)
		if runner.qaReport != "" {
			reportPath := filepath.Join(runner.gitRoot, qaReportRelPathForTest())
			if err := os.MkdirAll(filepath.Dir(reportPath), 0o755); err != nil {
				return agent.ExecuteResult{}, err
			}
			if err := os.WriteFile(reportPath, []byte(runner.qaReport), 0o644); err != nil {
				return agent.ExecuteResult{}, err
			}
		}
		return agent.ExecuteResult{LogPath: req.LogPath}, nil
	}
	if err := runner.errByTask[taskID]; err != nil {
		if agent.IsStopError(err) && sink != nil {
			// Mirror the real runner: a stopped Agent publishes its status
			// event before returning.
			payload, _ := json.Marshal(struct {
				Status string `json:"status"`
			}{Status: "stopped"})
			_ = sink.Publish(ctx, runevent.RunEvent{
				RunID:   req.RunID,
				Batch:   req.Batch.Number,
				Source:  runevent.SourceAgent,
				Kind:    runevent.KindAgentStatus,
				Summary: "SESSION STOPPED\n",
				Payload: payload,
			})
		}
		return agent.ExecuteResult{}, err
	}
	if path := runner.writeByTask[taskID]; path != "" {
		full := filepath.Join(runner.gitRoot, path)
		if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
			return agent.ExecuteResult{}, err
		}
		if err := os.WriteFile(full, []byte("agent change for "+taskID+"\n"), 0o644); err != nil {
			return agent.ExecuteResult{}, err
		}
	}
	if status, ok := runner.statusByTask[taskID]; ok {
		if err := spec.SetStatus(taskPathFor(runner.gitRoot, taskCycleSlug, taskID), status); err != nil {
			return agent.ExecuteResult{}, err
		}
	}
	return agent.ExecuteResult{LogPath: req.LogPath}, nil
}

func (runner *taskFakeRunner) EndSession(context.Context, agent.RuntimeSpec, agent.SessionRef) error {
	return nil
}

func taskIDFromPrompt(prompt string) string {
	for _, line := range strings.Split(prompt, "\n") {
		if strings.HasPrefix(line, "Task: ") {
			return strings.TrimPrefix(line, "Task: ")
		}
	}
	return ""
}

const qaReportNameForTest = "qa-report-2026-01-01.md"

func qaReportRelPathForTest() string {
	return filepath.Join("docs", "specs", taskCycleSlug, "qa", qaReportNameForTest)
}

func qaReportForTest(verdict string) string {
	return fmt.Sprintf("---\nverdict: %s\n---\n\n# QA Report\n", verdict)
}

// taskFakeVerifier records every verification command verbatim and fails
// the commands scripted in failOn.
type taskFakeVerifier struct {
	calls      *[]string
	store      *store.Store
	runID      string
	failOn     map[string]error
	commands   []string
	workDirs   []string
	seenStates []string
}

func (verifier *taskFakeVerifier) Verify(_ context.Context, req VerifyRequest) error {
	*verifier.calls = append(*verifier.calls, "verify")
	verifier.commands = append(verifier.commands, req.Command)
	verifier.workDirs = append(verifier.workDirs, req.WorkDir)
	if verifier.store != nil {
		verifier.seenStates = append(verifier.seenStates, runStateForTest(verifier.store, verifier.runID))
	}
	return verifier.failOn[req.Command]
}

func taskEventsOfKind(sink *captureEventSink, kind runevent.Kind) []runevent.RunEvent {
	sink.mu.Lock()
	defer sink.mu.Unlock()
	var matched []runevent.RunEvent
	for _, event := range sink.events {
		if event.Kind == kind {
			matched = append(matched, event)
		}
	}
	return matched
}

func TestTaskCycleValidatesPlan(t *testing.T) {
	fixture := newTaskCycleFixture(t, []taskSpecSeed{{id: "task_01"}})
	engine := fixture.engine(t, &taskFakeRunner{calls: fixture.calls, gitRoot: fixture.gitRoot}, &taskFakeVerifier{calls: fixture.calls}, &engineFakeCommitter{calls: fixture.calls}, fixture.worktree)

	_, err := engine.TaskCycle(context.Background(), TaskPlan{})

	if err == nil || !strings.Contains(err.Error(), "required") {
		t.Fatalf("expected plan validation error, got %v", err)
	}
	if len(*fixture.calls) != 0 {
		t.Fatalf("expected no daemon actions for invalid plan, got %v", *fixture.calls)
	}
}

func TestTaskCycleExecutesAgentVerifySettleCommitContract(t *testing.T) {
	fixture := newTaskCycleFixture(t, []taskSpecSeed{
		{id: "task_01", title: "Write the usage docs", taskType: "docs", verification: []string{"echo alpha", "echo beta"}},
		{id: "task_02", title: "Add the backend behavior", needs: []string{"task_01"}, verification: []string{"echo gamma"}},
	})
	// Snapshot script: task_01 before/after, then task_02 before/after.
	fixture.worktree.snapshots = [][]string{nil, {"src/one.go"}, nil, {"src/two.go"}}
	runner := &taskFakeRunner{
		calls:   fixture.calls,
		gitRoot: fixture.gitRoot,
		store:   fixture.store,
		statusByTask: map[string]spec.Status{
			"task_01": spec.StatusCompleted,
			"task_02": spec.StatusCompleted,
		},
	}
	verifier := &taskFakeVerifier{calls: fixture.calls, store: fixture.store, runID: fixture.run.ID}
	committer := &engineFakeCommitter{calls: fixture.calls}
	engine := fixture.engine(t, runner, verifier, committer, fixture.worktree)

	result, err := engine.TaskCycle(context.Background(), fixture.plan())

	if err != nil {
		t.Fatalf("task cycle: %v", err)
	}
	if got := strings.Join(*fixture.calls, ">"); got != "agent>verify>verify>commit>agent>verify>commit" {
		t.Fatalf("expected agent>verify>settle>commit contract per Task, got %q", got)
	}
	if result.Completed != 2 || result.Failed != 0 || result.Skipped != 0 || result.QAVerdict != "" {
		t.Fatalf("unexpected cycle result: %+v", result)
	}
	expectedFirst := "docs: Write the usage docs\n\nRoundfix-Spec: " + taskCycleSlug + "\nRoundfix-Task: task_01"
	if committer.messages[0] != expectedFirst {
		t.Fatalf("expected docs commit message with trailers, got %q", committer.messages[0])
	}
	expectedSecond := "feat: Add the backend behavior\n\nRoundfix-Spec: " + taskCycleSlug + "\nRoundfix-Task: task_02"
	if committer.messages[1] != expectedSecond {
		t.Fatalf("expected default feat commit message with trailers, got %q", committer.messages[1])
	}
	if got := strings.Join(committer.paths[0], "|"); got != taskFileRel(taskCycleSlug, "task_01")+"|src/one.go" {
		t.Fatalf("expected snapshot diff plus task file staged, got %q", got)
	}
	if got := strings.Join(verifier.commands, "|"); got != "echo alpha|echo beta|echo gamma" {
		t.Fatalf("expected every Verification command verbatim in order, got %q", got)
	}
	for _, workDir := range verifier.workDirs {
		if workDir != fixture.gitRoot {
			t.Fatalf("expected verification in WorkDir %q, got %q", fixture.gitRoot, workDir)
		}
	}
	if runner.requests[0].Batch.Number != 1 || runner.requests[1].Batch.Number != 2 {
		t.Fatalf("expected 1-based execution ordinals as Batch numbers, got %d and %d", runner.requests[0].Batch.Number, runner.requests[1].Batch.Number)
	}
	expectedSession := agent.SessionRefForRun(fixture.run.ID, fixture.gitRoot)
	for _, req := range runner.requests {
		if req.Session != expectedSession {
			t.Fatalf("expected shared Agent Session %#v, got %#v", expectedSession, req.Session)
		}
	}
	expectedLog := filepath.Join(fixture.gitRoot, ".roundfix", "runs", fixture.run.ID, "agent", "batch-001.log")
	if runner.requests[0].LogPath != expectedLog {
		t.Fatalf("expected review-path log convention, got %q", runner.requests[0].LogPath)
	}
	if runner.requests[0].GitRoot != fixture.gitRoot {
		t.Fatalf("expected Agent working directory %q, got %q", fixture.gitRoot, runner.requests[0].GitRoot)
	}
	if !strings.Contains(runner.prompts[0], "Task: task_01") || !strings.Contains(runner.prompts[0], "## Verification") {
		t.Fatalf("expected task prompt with identity and fresh file content, got %q", runner.prompts[0])
	}
	for _, state := range runner.seenStates {
		if state != store.StateResolvingWithAgent {
			t.Fatalf("expected ResolvingWithAgent during Agent run, got %q", state)
		}
	}
	for _, state := range verifier.seenStates {
		if state != store.StateVerifying {
			t.Fatalf("expected Verifying during verification, got %q", state)
		}
	}
	if got := taskStatusOnDisk(t, fixture.gitRoot, "task_01"); got != string(spec.StatusCompleted) {
		t.Fatalf("expected task_01 completed on disk, got %q", got)
	}
	taskEvents := taskEventsOfKind(fixture.sink, runevent.KindDaemonTask)
	if len(taskEvents) != 4 {
		t.Fatalf("expected started+settled events per Task, got %d", len(taskEvents))
	}
	if taskEvents[0].ReviewIssue != "task_01" || taskEvents[0].Batch != 1 {
		t.Fatalf("expected Task id in Work Item field with ordinal Batch, got %+v", taskEvents[0])
	}
	if !strings.Contains(string(taskEvents[1].Payload), `"settled"`) || !strings.Contains(string(taskEvents[1].Payload), `"completed"`) {
		t.Fatalf("expected settled completed payload, got %s", taskEvents[1].Payload)
	}
	kinds := fixture.sink.kinds()
	if kinds[len(kinds)-1] != runevent.KindDaemonOutcome {
		t.Fatalf("expected outcome event at cycle end, got %v", kinds)
	}
	for _, kind := range kinds {
		if kind == runevent.KindDaemonPush || kind == runevent.KindDaemonSourceResolution {
			t.Fatalf("expected no push or Review Source events in the spec path, got %v", kinds)
		}
	}
	if len(fixture.pusher.remotes) != 0 || len(fixture.source.requests) != 0 {
		t.Fatal("expected Pusher and Review Source resolver never invoked for spec Runs")
	}
}

func TestTaskCycleFailedTaskSkipsDependentsAndContinuesIndependents(t *testing.T) {
	fixture := newTaskCycleFixture(t, []taskSpecSeed{
		{id: "task_01", verification: []string{"echo t1"}},
		{id: "task_02", needs: []string{"task_01"}, verification: []string{"echo t2"}},
		{id: "task_03", title: "Tidy the build files", taskType: "chore", verification: []string{"echo t3"}},
	})
	// task_01's Agent claims completed, but the Daemon settles failed when
	// verification fails (ADR 0014).
	runner := &taskFakeRunner{
		calls:   fixture.calls,
		gitRoot: fixture.gitRoot,
		statusByTask: map[string]spec.Status{
			"task_01": spec.StatusCompleted,
			"task_03": spec.StatusCompleted,
		},
	}
	verifier := &taskFakeVerifier{calls: fixture.calls, failOn: map[string]error{"echo t1": errors.New("gate broke")}}
	committer := &engineFakeCommitter{calls: fixture.calls}
	engine := fixture.engine(t, runner, verifier, committer, fixture.worktree)

	result, err := engine.TaskCycle(context.Background(), fixture.plan())

	if err != nil {
		t.Fatalf("expected failed Task to not halt the cycle, got %v", err)
	}
	if result.Completed != 1 || result.Failed != 1 || result.Skipped != 1 {
		t.Fatalf("expected 1 completed, 1 failed, 1 skipped, got %+v", result)
	}
	if got := strings.Join(*fixture.calls, ">"); got != "agent>verify>agent>verify>commit" {
		t.Fatalf("expected no commit for the failed Task and a full run for the independent Task, got %q", got)
	}
	if got := taskStatusOnDisk(t, fixture.gitRoot, "task_01"); got != string(spec.StatusFailed) {
		t.Fatalf("expected Agent-claimed completed settled failed after failed verification, got %q", got)
	}
	if got := taskStatusOnDisk(t, fixture.gitRoot, "task_02"); got != string(spec.StatusPending) {
		t.Fatalf("expected skipped dependent left pending, got %q", got)
	}
	if got := taskStatusOnDisk(t, fixture.gitRoot, "task_03"); got != string(spec.StatusCompleted) {
		t.Fatalf("expected independent Task completed, got %q", got)
	}
	if len(committer.messages) != 1 || !strings.HasPrefix(committer.messages[0], "chore: Tidy the build files") {
		t.Fatalf("expected only the independent Task committed with chore mapping, got %v", committer.messages)
	}
	if runner.requests[1].Batch.Number != 2 {
		t.Fatalf("expected skipped Task to not consume an execution ordinal, got Batch %d", runner.requests[1].Batch.Number)
	}
	var skipEvent *runevent.RunEvent
	for _, event := range taskEventsOfKind(fixture.sink, runevent.KindDaemonTask) {
		if strings.Contains(string(event.Payload), `"skipped"`) {
			matched := event
			skipEvent = &matched
		}
	}
	if skipEvent == nil {
		t.Fatalf("expected a skipped daemon.task event, got %v", fixture.sink.kinds())
	}
	if skipEvent.ReviewIssue != "task_02" || skipEvent.Batch != 0 || !strings.Contains(string(skipEvent.Payload), "task_01") {
		t.Fatalf("expected skip event naming the unmet need, got %+v payload %s", skipEvent, skipEvent.Payload)
	}
	var settledFailed bool
	for _, event := range taskEventsOfKind(fixture.sink, runevent.KindDaemonTask) {
		if event.ReviewIssue == "task_01" && strings.Contains(string(event.Payload), "verification failed") {
			settledFailed = true
		}
	}
	if !settledFailed {
		t.Fatal("expected failed settlement event journaling the verification failure reason")
	}
	kinds := fixture.sink.kinds()
	last := taskEventsOfKind(fixture.sink, runevent.KindDaemonOutcome)
	if kinds[len(kinds)-1] != runevent.KindDaemonOutcome || !strings.Contains(string(last[0].Payload), `"skipped":1`) {
		t.Fatalf("expected outcome event with counts at cycle end, got %v", kinds)
	}
}

func TestTaskCycleSettlesForgottenAgentStatus(t *testing.T) {
	t.Run("completed on passing verification", func(t *testing.T) {
		fixture := newTaskCycleFixture(t, []taskSpecSeed{{id: "task_01", verification: []string{"echo alpha", "echo beta"}}})
		runner := &taskFakeRunner{calls: fixture.calls, gitRoot: fixture.gitRoot}
		verifier := &taskFakeVerifier{calls: fixture.calls}
		committer := &engineFakeCommitter{calls: fixture.calls}
		engine := fixture.engine(t, runner, verifier, committer, fixture.worktree)

		result, err := engine.TaskCycle(context.Background(), fixture.plan())

		if err != nil {
			t.Fatalf("task cycle: %v", err)
		}
		if result.Completed != 1 || result.Failed != 0 {
			t.Fatalf("expected one completed Task, got %+v", result)
		}
		if got := taskStatusOnDisk(t, fixture.gitRoot, "task_01"); got != string(spec.StatusCompleted) {
			t.Fatalf("expected Daemon to settle the forgotten status completed, got %q", got)
		}
		if got := strings.Join(verifier.commands, "|"); got != "echo alpha|echo beta" {
			t.Fatalf("expected every Verification command to run, got %q", got)
		}
		if len(committer.messages) != 1 {
			t.Fatalf("expected the settled Task committed, got %v", committer.messages)
		}
	})

	t.Run("failed on failing verification", func(t *testing.T) {
		fixture := newTaskCycleFixture(t, []taskSpecSeed{{id: "task_01", verification: []string{"echo gamma", "echo delta"}}})
		runner := &taskFakeRunner{calls: fixture.calls, gitRoot: fixture.gitRoot}
		verifier := &taskFakeVerifier{calls: fixture.calls, failOn: map[string]error{"echo gamma": errors.New("gate broke")}}
		committer := &engineFakeCommitter{calls: fixture.calls}
		engine := fixture.engine(t, runner, verifier, committer, fixture.worktree)

		result, err := engine.TaskCycle(context.Background(), fixture.plan())

		if err != nil {
			t.Fatalf("expected failed verification to fail only the Task, got %v", err)
		}
		if result.Completed != 0 || result.Failed != 1 {
			t.Fatalf("expected one failed Task, got %+v", result)
		}
		if got := taskStatusOnDisk(t, fixture.gitRoot, "task_01"); got != string(spec.StatusFailed) {
			t.Fatalf("expected Daemon to settle the forgotten status failed, got %q", got)
		}
		if got := strings.Join(verifier.commands, "|"); got != "echo gamma" {
			t.Fatalf("expected verification to stop at the first failing command, got %q", got)
		}
		if len(committer.messages) != 0 {
			t.Fatalf("expected no commit for a failed Task, got %v", committer.messages)
		}
	})
}

func TestTaskCycleAgentErrorSettlesTaskFailedAndContinues(t *testing.T) {
	fixture := newTaskCycleFixture(t, []taskSpecSeed{
		{id: "task_01"},
		{id: "task_02"},
	})
	runner := &taskFakeRunner{
		calls:        fixture.calls,
		gitRoot:      fixture.gitRoot,
		errByTask:    map[string]error{"task_01": errors.New("agent crashed")},
		statusByTask: map[string]spec.Status{"task_02": spec.StatusCompleted},
	}
	verifier := &taskFakeVerifier{calls: fixture.calls}
	committer := &engineFakeCommitter{calls: fixture.calls}
	engine := fixture.engine(t, runner, verifier, committer, fixture.worktree)

	result, err := engine.TaskCycle(context.Background(), fixture.plan())

	if err != nil {
		t.Fatalf("expected Agent error to fail only the Task, got %v", err)
	}
	if result.Completed != 1 || result.Failed != 1 || result.Skipped != 0 {
		t.Fatalf("expected the independent Task to continue after the Agent error, got %+v", result)
	}
	if got := strings.Join(*fixture.calls, ">"); got != "agent>agent>verify>commit" {
		t.Fatalf("expected no verification or commit for the errored Task, got %q", got)
	}
	if got := taskStatusOnDisk(t, fixture.gitRoot, "task_01"); got != string(spec.StatusFailed) {
		t.Fatalf("expected Agent error settled failed, got %q", got)
	}
	var journaled bool
	for _, event := range taskEventsOfKind(fixture.sink, runevent.KindDaemonTask) {
		if event.ReviewIssue == "task_01" && strings.Contains(string(event.Payload), "agent crashed") {
			journaled = true
		}
	}
	if !journaled {
		t.Fatal("expected the Agent failure reason journaled in the settlement event")
	}
}

func TestTaskCycleCommitStagesSnapshotDiffPlusTaskFile(t *testing.T) {
	fixture := newTaskCycleFixture(t, []taskSpecSeed{{id: "task_01"}})
	taskFile := taskFileRel(taskCycleSlug, "task_01")
	// user-wip.txt and the task file are dirty before the Task starts: the
	// pre-existing user change stays out of the commit, while the task file
	// is re-added because the settled status must ride in the Task commit.
	fixture.worktree.snapshots = [][]string{
		{"user-wip.txt", taskFile},
		{"user-wip.txt", taskFile, "src/x.go"},
	}
	runner := &taskFakeRunner{calls: fixture.calls, gitRoot: fixture.gitRoot, statusByTask: map[string]spec.Status{"task_01": spec.StatusCompleted}}
	committer := &engineFakeCommitter{calls: fixture.calls}
	engine := fixture.engine(t, runner, &taskFakeVerifier{calls: fixture.calls}, committer, fixture.worktree)

	if _, err := engine.TaskCycle(context.Background(), fixture.plan()); err != nil {
		t.Fatalf("task cycle: %v", err)
	}
	if len(committer.paths) != 1 {
		t.Fatalf("expected one commit, got %v", committer.paths)
	}
	if got := strings.Join(committer.paths[0], "|"); got != taskFile+"|src/x.go" {
		t.Fatalf("expected the task file ensured and the pre-existing change excluded, got %q", got)
	}
}

func TestTaskCycleStopBeforeTaskPublishesStopAndDoesNothing(t *testing.T) {
	fixture := newTaskCycleFixture(t, []taskSpecSeed{{id: "task_01"}})
	engine := fixture.engine(t, &taskFakeRunner{calls: fixture.calls, gitRoot: fixture.gitRoot}, &taskFakeVerifier{calls: fixture.calls}, &engineFakeCommitter{calls: fixture.calls}, fixture.worktree)
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	result, err := engine.TaskCycle(ctx, fixture.plan())

	if !errors.Is(err, context.Canceled) {
		t.Fatalf("expected canceled cycle, got %v", err)
	}
	if len(*fixture.calls) != 0 {
		t.Fatalf("expected no daemon actions after Stop Request, got %v", *fixture.calls)
	}
	if result.Completed != 0 || result.Failed != 0 || result.Skipped != 0 {
		t.Fatalf("expected no settled Tasks, got %+v", result)
	}
	kinds := fixture.sink.kinds()
	if len(kinds) == 0 || kinds[len(kinds)-1] != runevent.KindDaemonStatus {
		t.Fatalf("expected daemon stop event published to the sink, got %v", kinds)
	}
}

func TestTaskCycleStopDuringAgentPreservesTaskAndHalts(t *testing.T) {
	fixture := newTaskCycleFixture(t, []taskSpecSeed{
		{id: "task_01"},
		{id: "task_02"},
	})
	runner := &taskFakeRunner{
		calls:     fixture.calls,
		gitRoot:   fixture.gitRoot,
		errByTask: map[string]error{"task_01": agent.StopError{Err: context.Canceled}},
	}
	engine := fixture.engine(t, runner, &taskFakeVerifier{calls: fixture.calls}, &engineFakeCommitter{calls: fixture.calls}, fixture.worktree)

	_, err := engine.TaskCycle(context.Background(), fixture.plan())

	if !agent.IsStopError(err) {
		t.Fatalf("expected StopError surfaced, got %v", err)
	}
	if got := strings.Join(*fixture.calls, ">"); got != "agent" {
		t.Fatalf("expected no further Tasks or daemon actions after stop, got %q", got)
	}
	// Worktree preserved: the stopped Task keeps the status the Agent left.
	if got := taskStatusOnDisk(t, fixture.gitRoot, "task_01"); got != string(spec.StatusPending) {
		t.Fatalf("expected stopped Task left unsettled, got %q", got)
	}
	kinds := fixture.sink.kinds()
	if len(kinds) == 0 || kinds[len(kinds)-1] != runevent.KindAgentStatus {
		t.Fatalf("expected the runner's stop event to reach the sink last, got %v", kinds)
	}
	for _, kind := range kinds {
		if kind == runevent.KindDaemonVerification || kind == runevent.KindDaemonCommit || kind == runevent.KindDaemonOutcome {
			t.Fatalf("expected no unsafe daemon events after stop, got %v", kinds)
		}
	}
}

func TestTaskCycleRerunsStaleTasksAndSkipsCompletedTasks(t *testing.T) {
	fixture := newTaskCycleFixture(t, []taskSpecSeed{
		{id: "task_01", status: string(spec.StatusCompleted)},
		{id: "task_02", status: string(spec.StatusInProgress), needs: []string{"task_01"}},
		{id: "task_03", status: string(spec.StatusFailed)},
	})
	runner := &taskFakeRunner{
		calls:   fixture.calls,
		gitRoot: fixture.gitRoot,
		statusByTask: map[string]spec.Status{
			"task_02": spec.StatusCompleted,
			"task_03": spec.StatusCompleted,
		},
	}
	committer := &engineFakeCommitter{calls: fixture.calls}
	engine := fixture.engine(t, runner, &taskFakeVerifier{calls: fixture.calls}, committer, fixture.worktree)

	result, err := engine.TaskCycle(context.Background(), fixture.plan())

	if err != nil {
		t.Fatalf("task cycle: %v", err)
	}
	if result.Completed != 2 || result.Failed != 0 || result.Skipped != 0 {
		t.Fatalf("expected both stale Tasks re-run fresh, got %+v", result)
	}
	if len(runner.requests) != 2 {
		t.Fatalf("expected the completed Task to not run, got %d Agent runs", len(runner.requests))
	}
	// A Task completed in a prior Run satisfies needs, and execution
	// ordinals count executed Tasks only.
	if taskIDFromPrompt(runner.prompts[0]) != "task_02" || runner.requests[0].Batch.Number != 1 {
		t.Fatalf("expected stale in_progress Task re-run as Batch 001, got %q Batch %d", taskIDFromPrompt(runner.prompts[0]), runner.requests[0].Batch.Number)
	}
	if taskIDFromPrompt(runner.prompts[1]) != "task_03" || runner.requests[1].Batch.Number != 2 {
		t.Fatalf("expected failed Task re-run as Batch 002, got %q Batch %d", taskIDFromPrompt(runner.prompts[1]), runner.requests[1].Batch.Number)
	}
	for _, event := range taskEventsOfKind(fixture.sink, runevent.KindDaemonTask) {
		if event.ReviewIssue == "task_01" {
			t.Fatalf("expected no events for the already-completed Task, got %+v", event)
		}
	}
}

func TestTaskCycleRealRepoCommitsPerTaskExcludingPreexistingDirt(t *testing.T) {
	fixture := newTaskCycleFixture(t, []taskSpecSeed{
		{id: "task_01", title: "Write the usage docs", taskType: "docs", verification: []string{"true"}},
		{id: "task_02", title: "Add the backend behavior", needs: []string{"task_01"}, verification: []string{"true"}},
	})
	repoDir := fixture.gitRoot
	runGitForTest(t, repoDir, "init", "-q", "-b", "main")
	runGitForTest(t, repoDir, "config", "user.name", "Roundfix Test")
	runGitForTest(t, repoDir, "config", "user.email", "test@example.com")
	runGitForTest(t, repoDir, "config", "commit.gpgsign", "false")
	runGitForTest(t, repoDir, "add", "-A")
	runGitForTest(t, repoDir, "commit", "-q", "-m", "initial")
	// Pre-existing user work that must never enter a Task commit.
	mustWriteForTest(t, filepath.Join(repoDir, "user-wip.txt"), "wip\n")

	runner := &taskFakeRunner{
		calls:   fixture.calls,
		gitRoot: repoDir,
		writeByTask: map[string]string{
			"task_01": "docs/usage.md",
			"task_02": "internal/feature.go",
		},
		// task_01's Agent settles its status; task_02's Agent forgets and
		// the Daemon settles completed after passing verification.
		statusByTask: map[string]spec.Status{"task_01": spec.StatusCompleted},
	}
	engine := fixture.engine(t, runner, ExecVerifier{}, GitCommitter{}, GitWorktreeSnapshotter{})

	result, err := engine.TaskCycle(context.Background(), fixture.plan())

	if err != nil {
		t.Fatalf("task cycle: %v", err)
	}
	if result.Completed != 2 || result.Failed != 0 || result.Skipped != 0 {
		t.Fatalf("expected both Tasks completed, got %+v", result)
	}
	firstMessage := strings.TrimSpace(runGitForTest(t, repoDir, "show", "-s", "--format=%B", "HEAD~1"))
	if firstMessage != "docs: Write the usage docs\n\nRoundfix-Spec: "+taskCycleSlug+"\nRoundfix-Task: task_01" {
		t.Fatalf("expected docs commit message with trailers, got %q", firstMessage)
	}
	secondMessage := strings.TrimSpace(runGitForTest(t, repoDir, "show", "-s", "--format=%B", "HEAD"))
	if secondMessage != "feat: Add the backend behavior\n\nRoundfix-Spec: "+taskCycleSlug+"\nRoundfix-Task: task_02" {
		t.Fatalf("expected feat commit message with trailers, got %q", secondMessage)
	}
	firstFiles := commitFilesForTest(t, repoDir, "HEAD~1")
	if got := strings.Join(firstFiles, "|"); got != taskFileRel(taskCycleSlug, "task_01")+"|docs/usage.md" {
		t.Fatalf("expected only task_01's change and task file committed, got %q", got)
	}
	secondFiles := commitFilesForTest(t, repoDir, "HEAD")
	if got := strings.Join(secondFiles, "|"); got != taskFileRel(taskCycleSlug, "task_02")+"|internal/feature.go" {
		t.Fatalf("expected only task_02's change and task file committed, got %q", got)
	}
	status := runGitForTest(t, repoDir, "status", "--porcelain=v1")
	if !strings.Contains(status, "?? user-wip.txt") {
		t.Fatalf("expected pre-existing user work preserved and uncommitted, got %q", status)
	}
	if got := taskStatusOnDisk(t, repoDir, "task_02"); got != string(spec.StatusCompleted) {
		t.Fatalf("expected forgotten status settled completed, got %q", got)
	}
}

func TestTaskCycleQAVerdictMatrixSettlesRunAndCommitsReport(t *testing.T) {
	reportRel := qaReportRelPathForTest()
	tests := []struct {
		name        string
		report      string
		wantVerdict string
		wantCommit  bool
	}{
		{name: "pass", report: qaReportForTest(spec.VerdictPass), wantVerdict: spec.VerdictPass, wantCommit: true},
		{name: "partial", report: qaReportForTest(spec.VerdictPartial), wantVerdict: spec.VerdictPartial, wantCommit: true},
		{name: "fail", report: qaReportForTest(spec.VerdictFail), wantVerdict: spec.VerdictFail, wantCommit: true},
		{name: "missing report", report: "", wantVerdict: "missing", wantCommit: false},
		{name: "unreadable verdict", report: "---\nsummary: no verdict field\n---\n\n# QA Report\n", wantVerdict: "unreadable", wantCommit: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fixture := newTaskCycleFixture(t, []taskSpecSeed{{id: "task_01", title: "Build the feature"}})
			// Snapshot script: task_01 before/after, then the QA step
			// before/after — only the QA Report is new in the QA diff.
			fixture.worktree.snapshots = [][]string{nil, {"src/one.go"}, {"src/one.go"}, {"src/one.go", reportRel}}
			runner := &taskFakeRunner{
				calls:        fixture.calls,
				gitRoot:      fixture.gitRoot,
				store:        fixture.store,
				statusByTask: map[string]spec.Status{"task_01": spec.StatusCompleted},
				qaReport:     tt.report,
			}
			committer := &engineFakeCommitter{calls: fixture.calls}
			engine := fixture.engine(t, runner, &taskFakeVerifier{calls: fixture.calls}, committer, fixture.worktree)

			result, err := engine.TaskCycle(context.Background(), fixture.qaPlan())

			if err != nil {
				t.Fatalf("task cycle: %v", err)
			}
			if result.Completed != 1 || result.Failed != 0 || result.Skipped != 0 {
				t.Fatalf("unexpected Task counts: %+v", result)
			}
			if result.QAVerdict != tt.wantVerdict {
				t.Fatalf("expected QA verdict %q, got %q", tt.wantVerdict, result.QAVerdict)
			}
			wantReportPath := reportRel
			if tt.report == "" {
				wantReportPath = ""
			}
			if result.QAReportPath != wantReportPath {
				t.Fatalf("expected QA Report path %q, got %q", wantReportPath, result.QAReportPath)
			}
			if len(runner.qaPrompts) != 1 || !strings.Contains(runner.qaPrompts[0], "Spec: "+taskCycleSlug) {
				t.Fatalf("expected one QA prompt naming the Spec, got %v", runner.qaPrompts)
			}
			// The QA step runs as the next Batch ordinal after the Tasks.
			last := runner.requests[len(runner.requests)-1]
			if last.Batch.Number != 2 {
				t.Fatalf("expected the QA step as Batch 2, got %d", last.Batch.Number)
			}
			qaEvents := taskEventsOfKind(fixture.sink, runevent.KindDaemonQA)
			if len(qaEvents) != 1 || qaEvents[0].Batch != 2 {
				t.Fatalf("expected one daemon.qa event on the QA Batch, got %+v", qaEvents)
			}
			payload := string(qaEvents[0].Payload)
			if !strings.Contains(payload, fmt.Sprintf("%q", tt.wantVerdict)) || !strings.Contains(payload, wantReportPath) {
				t.Fatalf("expected daemon.qa payload with verdict and report path, got %s", payload)
			}
			if tt.wantCommit {
				if len(committer.messages) != 2 {
					t.Fatalf("expected the QA Report committed in its own commit, got %v", committer.messages)
				}
				wantMessage := "docs(qa): qa report for " + taskCycleSlug + " (" + tt.wantVerdict + ")\n\nRoundfix-Spec: " + taskCycleSlug
				if committer.messages[1] != wantMessage {
					t.Fatalf("expected QA commit message %q, got %q", wantMessage, committer.messages[1])
				}
				if got := strings.Join(committer.paths[1], "|"); got != reportRel {
					t.Fatalf("expected only the QA step's changes in the QA commit, got %q", got)
				}
			} else if len(committer.messages) != 1 {
				t.Fatalf("expected no QA commit for a missing report, got %v", committer.messages)
			}
			kinds := fixture.sink.kinds()
			if kinds[len(kinds)-1] != runevent.KindDaemonOutcome {
				t.Fatalf("expected the outcome event after the QA step, got %v", kinds)
			}
		})
	}
}

func TestTaskCycleQAStepSkippedUnlessEveryTaskCompleted(t *testing.T) {
	fixture := newTaskCycleFixture(t, []taskSpecSeed{
		{id: "task_01"},
		{id: "task_02", needs: []string{"task_01"}},
	})
	runner := &taskFakeRunner{
		calls:        fixture.calls,
		gitRoot:      fixture.gitRoot,
		statusByTask: map[string]spec.Status{"task_01": spec.StatusFailed},
		qaReport:     qaReportForTest(spec.VerdictPass),
	}
	committer := &engineFakeCommitter{calls: fixture.calls}
	engine := fixture.engine(t, runner, &taskFakeVerifier{calls: fixture.calls}, committer, fixture.worktree)

	result, err := engine.TaskCycle(context.Background(), fixture.qaPlan())

	if err != nil {
		t.Fatalf("task cycle: %v", err)
	}
	if result.Failed != 1 || result.Skipped != 1 {
		t.Fatalf("expected the failed and skipped Tasks settled, got %+v", result)
	}
	if len(runner.qaPrompts) != 0 {
		t.Fatalf("expected the QA step never invoked with a failed Task, got %d QA prompt(s)", len(runner.qaPrompts))
	}
	if result.QAVerdict != "" || result.QAReportPath != "" {
		t.Fatalf("expected no QA verdict when the step is skipped, got %+v", result)
	}
	if events := taskEventsOfKind(fixture.sink, runevent.KindDaemonQA); len(events) != 0 {
		t.Fatalf("expected no daemon.qa event when the step is skipped, got %+v", events)
	}
	if len(committer.messages) != 0 {
		t.Fatalf("expected no commits, got %v", committer.messages)
	}
}

func TestTaskCycleQAOnlyRunWhenEveryTaskAlreadyCompleted(t *testing.T) {
	fixture := newTaskCycleFixture(t, []taskSpecSeed{
		{id: "task_01", status: string(spec.StatusCompleted)},
		{id: "task_02", status: string(spec.StatusCompleted), needs: []string{"task_01"}},
	})
	fixture.worktree.snapshots = [][]string{nil, {qaReportRelPathForTest()}}
	runner := &taskFakeRunner{
		calls:    fixture.calls,
		gitRoot:  fixture.gitRoot,
		store:    fixture.store,
		qaReport: qaReportForTest(spec.VerdictPass),
	}
	committer := &engineFakeCommitter{calls: fixture.calls}
	engine := fixture.engine(t, runner, &taskFakeVerifier{calls: fixture.calls}, committer, fixture.worktree)

	result, err := engine.TaskCycle(context.Background(), fixture.qaPlan())

	if err != nil {
		t.Fatalf("task cycle: %v", err)
	}
	if result.Completed != 0 || result.Failed != 0 || result.Skipped != 0 {
		t.Fatalf("expected no Task executions in a QA-only Run, got %+v", result)
	}
	if result.QAVerdict != spec.VerdictPass || result.QAReportPath != qaReportRelPathForTest() {
		t.Fatalf("expected the pass verdict from the fresh report, got %+v", result)
	}
	if got := strings.Join(*fixture.calls, ">"); got != "agent>commit" {
		t.Fatalf("expected only the QA Agent run and the QA Report commit, got %q", got)
	}
	if len(runner.qaPrompts) != 1 || !strings.Contains(runner.qaPrompts[0], "PRD: ") {
		t.Fatalf("expected the QA prompt with the PRD path, got %v", runner.qaPrompts)
	}
	if runner.requests[0].Batch.Number != 1 {
		t.Fatalf("expected the QA step as Batch 1 in a QA-only Run, got %d", runner.requests[0].Batch.Number)
	}
	for _, state := range runner.seenStates {
		if state != store.StateResolvingWithAgent {
			t.Fatalf("expected ResolvingWithAgent during the QA Agent run, got %q", state)
		}
	}
	if len(committer.messages) != 1 || !strings.HasPrefix(committer.messages[0], "docs(qa): qa report for "+taskCycleSlug+" (pass)") {
		t.Fatalf("expected the QA Report commit, got %v", committer.messages)
	}
}

// commitFilesForTest returns the file names of one commit, sorted with the
// spec task file first for deterministic comparison.
func commitFilesForTest(t *testing.T, repoDir string, rev string) []string {
	t.Helper()
	output := runGitForTest(t, repoDir, "show", "--name-only", "--format=", rev)
	var files []string
	for _, line := range strings.Split(strings.TrimSpace(output), "\n") {
		if line != "" {
			files = append(files, line)
		}
	}
	sort.Strings(files)
	return files
}
