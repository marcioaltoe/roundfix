package daemon

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"sort"
	"strings"
	"sync"
	"testing"
	"time"

	"roundfix/internal/agent"
	roundconfig "roundfix/internal/config"
	"roundfix/internal/gittest"
	"roundfix/internal/rounds"
	"roundfix/internal/runevent"
	"roundfix/internal/spec"
	"roundfix/internal/store"
	runworktree "roundfix/internal/worktree"
)

const taskCycleSlug = "0001-sample-feature"

func TestTaskCommitMessageDerivesSubjectAndTrailers(t *testing.T) {
	tests := []struct {
		name     string
		taskType spec.TaskType
		title    string
		want     string
	}{
		{
			name:     "default type lowercases ascii letter",
			taskType: spec.TaskTypeBackend,
			title:    "Build the acpx invocation core",
			want:     "feat: build the acpx invocation core\n\nRoundfix-Spec: 0003-dogfood-polish\nRoundfix-Task: task_01",
		},
		{
			name:     "docs type passes through",
			taskType: spec.TaskTypeDocs,
			title:    "Write the usage docs",
			want:     "docs: write the usage docs\n\nRoundfix-Spec: 0003-dogfood-polish\nRoundfix-Task: task_01",
		},
		{
			name:     "test type passes through",
			taskType: spec.TaskTypeTest,
			title:    "Add commit-message tests",
			want:     "test: add commit-message tests\n\nRoundfix-Spec: 0003-dogfood-polish\nRoundfix-Task: task_01",
		},
		{
			name:     "chore type passes through",
			taskType: spec.TaskTypeChore,
			title:    "Refresh fixtures",
			want:     "chore: refresh fixtures\n\nRoundfix-Spec: 0003-dogfood-polish\nRoundfix-Task: task_01",
		},
		{
			name:     "digit first title passes through",
			taskType: spec.TaskTypeBackend,
			title:    "2FA setup",
			want:     "feat: 2FA setup\n\nRoundfix-Spec: 0003-dogfood-polish\nRoundfix-Task: task_01",
		},
		{
			name:     "unicode first title lowercases",
			taskType: spec.TaskTypeBackend,
			title:    "Über tracing",
			want:     "feat: über tracing\n\nRoundfix-Spec: 0003-dogfood-polish\nRoundfix-Task: task_01",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			task := spec.Task{ID: "task_01", Type: tt.taskType, Title: tt.title}

			got := TaskCommitMessage("0003-dogfood-polish", task)

			if got != tt.want {
				t.Fatalf("TaskCommitMessage() = %q, want %q", got, tt.want)
			}
			if task.Title != tt.title {
				t.Fatalf("TaskCommitMessage mutated task title to %q, want %q", task.Title, tt.title)
			}
		})
	}
}

func TestQACommitMessageDerivesUnscopedSubjectAndTrailer(t *testing.T) {
	got := QACommitMessage("0003-dogfood-polish", "pass")
	want := "docs: qa report for 0003-dogfood-polish (pass)\n\nRoundfix-Spec: 0003-dogfood-polish"

	if got != want {
		t.Fatalf("QACommitMessage() = %q, want %q", got, want)
	}
}

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
	store       *store.Store
	run         store.Run
	gitRoot     string
	specsRoot   string
	artifactDir string
	graph       *spec.Graph
	calls       *[]string
	sink        *captureEventSink
	progress    *bytes.Buffer
	worktree    *engineFakeWorktree
	pusher      *engineFakePusher
	source      *engineFakeSource
}

func newTaskCycleFixture(t *testing.T, seeds []taskSpecSeed) *taskCycleFixture {
	t.Helper()
	ctx := context.Background()
	homeDir := t.TempDir()
	gitRoot := t.TempDir()
	artifactDir := filepath.Join(t.TempDir(), "artifacts")
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

	specsRoot := filepath.Join(gitRoot, "docs", "specs")
	graph, err := spec.Load(specsRoot, taskCycleSlug)
	if err != nil {
		t.Fatalf("load spec: %v", err)
	}

	calls := []string{}
	return &taskCycleFixture{
		store:       runStore,
		run:         run,
		gitRoot:     gitRoot,
		specsRoot:   specsRoot,
		artifactDir: artifactDir,
		graph:       graph,
		calls:       &calls,
		sink:        &captureEventSink{},
		progress:    &bytes.Buffer{},
		worktree:    &engineFakeWorktree{},
		pusher:      &engineFakePusher{calls: &calls},
		source:      &engineFakeSource{calls: &calls},
	}
}

func (fixture *taskCycleFixture) plan() TaskPlan {
	return TaskPlan{
		RunID:                   fixture.run.ID,
		Session:                 agent.SessionRefForRun(fixture.run.ID, fixture.gitRoot),
		WorkDir:                 fixture.gitRoot,
		RunWorktree:             runworktree.Ref{RunID: fixture.run.ID, Path: fixture.gitRoot, Branch: runworktree.BranchName(fixture.run.ID), UserRoot: fixture.gitRoot},
		TargetBranch:            fixture.run.LocalBranch,
		Spec:                    fixture.graph.Spec,
		SpecsRoot:               fixture.specsRoot,
		Tasks:                   fixture.graph.Tasks,
		Runtime:                 agent.RuntimeSpec{ID: "codex", DisplayName: "Codex"},
		ArtifactDir:             fixture.artifactDir,
		AgentLogs:               true,
		Concurrency:             1,
		VerificationConcurrency: 1,
	}
}

func (fixture *taskCycleFixture) qaPlan() TaskPlan {
	plan := fixture.plan()
	plan.QA = true
	return plan
}

func (fixture *taskCycleFixture) useExternalSpecRoot(t *testing.T, seeds []taskSpecSeed) string {
	t.Helper()
	specsRoot := filepath.Join(t.TempDir(), "external-specs")
	writeSpecDirAtRootForTest(t, specsRoot, taskCycleSlug, seeds)
	graph, err := spec.Load(specsRoot, taskCycleSlug)
	if err != nil {
		t.Fatalf("load external spec: %v", err)
	}
	fixture.specsRoot = specsRoot
	fixture.graph = graph
	return specsRoot
}

func (fixture *taskCycleFixture) engine(t *testing.T, runner agent.Runner, verifier Verifier, committer Committer, worktree WorktreeSnapshotter) *Engine {
	t.Helper()
	return fixture.engineWithTaskWorktrees(t, runner, verifier, committer, worktree, nil)
}

func (fixture *taskCycleFixture) engineWithTaskWorktrees(t *testing.T, runner agent.Runner, verifier Verifier, committer Committer, worktree WorktreeSnapshotter, taskWorktrees TaskWorktreeManager) *Engine {
	return fixture.engineWithTaskWorktreesAndPriorChanges(t, runner, verifier, committer, worktree, taskWorktrees, nil)
}

func (fixture *taskCycleFixture) engineWithTaskWorktreesAndPriorChanges(t *testing.T, runner agent.Runner, verifier Verifier, committer Committer, worktree WorktreeSnapshotter, taskWorktrees TaskWorktreeManager, priorChanges PriorChangedResolver) *Engine {
	t.Helper()
	engine, err := NewEngine(Dependencies{
		Runner:        runner,
		Verifier:      verifier,
		Committer:     committer,
		Pusher:        fixture.pusher,
		Source:        fixture.source,
		Runs:          fixture.store,
		Worktree:      worktree,
		TaskWorktrees: taskWorktrees,
		PriorChanges:  priorChanges,
		Sink:          fixture.sink,
		Progress:      fixture.progress,
	})
	if err != nil {
		t.Fatalf("new engine: %v", err)
	}
	return engine
}

func writeSpecDirForTest(t *testing.T, gitRoot string, slug string, seeds []taskSpecSeed) {
	t.Helper()
	writeSpecDirAtRootForTest(t, filepath.Join(gitRoot, "docs", "specs"), slug, seeds)
}

func writeSpecDirAtRootForTest(t *testing.T, specsRoot string, slug string, seeds []taskSpecSeed) {
	t.Helper()
	specDir := filepath.Join(specsRoot, slug)
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
	return taskPathInSpecRootFor(filepath.Join(gitRoot, "docs", "specs"), slug, id)
}

func taskPathInSpecRootFor(specsRoot string, slug string, id string) string {
	return filepath.Join(specsRoot, slug, id+".md")
}

func taskFileRel(slug string, id string) string {
	return filepath.Join("docs", "specs", slug, id+".md")
}

func taskStatusOnDisk(t *testing.T, gitRoot string, id string) string {
	t.Helper()
	return taskStatusInSpecRootOnDisk(t, filepath.Join(gitRoot, "docs", "specs"), id)
}

func taskStatusInSpecRootOnDisk(t *testing.T, specsRoot string, id string) string {
	t.Helper()
	content, err := os.ReadFile(taskPathInSpecRootFor(specsRoot, taskCycleSlug, id))
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
// parses from the prompt, including attempted Agent status edits that the
// Daemon must normalize. A QA prompt writes qaReport as the Spec's QA Report,
// the way the qa-gate Agent does; an empty qaReport writes none.
type taskFakeRunner struct {
	calls            *[]string
	gitRoot          string
	store            *store.Store
	statusByTask     map[string]spec.Status
	statusByTaskCall map[string][]spec.Status
	errByTask        map[string]error
	errByTaskCall    map[string][]error
	taskCalls        map[string]int
	writeByTask      map[string]string
	resultByTask     map[string]string
	rawStatusByTask  map[string]string
	anomalyByTask    map[string]string
	afterTask        func(string)
	qaReport         string
	qaPrompts        []string
	seenStates       []string
	prompts          []string
	requests         []agent.ExecuteRequest
	writeLogs        bool
}

func (runner *taskFakeRunner) Probe(context.Context, agent.ProbeRequest) error { return nil }

func (runner *taskFakeRunner) Run(ctx context.Context, req agent.ExecuteRequest, sink runevent.Sink) (agent.ExecuteResult, error) {
	*runner.calls = append(*runner.calls, "agent")
	runner.prompts = append(runner.prompts, req.Prompt)
	runner.requests = append(runner.requests, req)
	if runner.writeLogs && strings.TrimSpace(req.LogPath) != "" {
		if err := os.MkdirAll(filepath.Dir(req.LogPath), 0o755); err != nil {
			return agent.ExecuteResult{}, err
		}
		if err := os.WriteFile(req.LogPath, []byte("fake agent output\n"), 0o644); err != nil {
			return agent.ExecuteResult{}, err
		}
	}
	if runner.store != nil {
		runner.seenStates = append(runner.seenStates, runStateForTest(runner.store, req.RunID))
	}
	taskID := taskIDFromPrompt(req.Prompt)
	if runner.taskCalls == nil {
		runner.taskCalls = map[string]int{}
	}
	taskCall := runner.taskCalls[taskID]
	runner.taskCalls[taskID] = taskCall + 1
	if taskID == "" && strings.Contains(req.Prompt, "Spec QA gate") {
		runner.qaPrompts = append(runner.qaPrompts, req.Prompt)
		if runner.qaReport != "" {
			reportPath := filepath.Join(qaSpecDirFromPromptForTest(req.Prompt, runner.gitRoot), "qa", qaReportNameForTest)
			if err := os.MkdirAll(filepath.Dir(reportPath), 0o755); err != nil {
				return agent.ExecuteResult{}, err
			}
			if err := os.WriteFile(reportPath, []byte(runner.qaReport), 0o644); err != nil {
				return agent.ExecuteResult{}, err
			}
		}
		return agent.ExecuteResult{LogPath: req.LogPath}, nil
	}
	if taskErrs := runner.errByTaskCall[taskID]; taskCall < len(taskErrs) && taskErrs[taskCall] != nil {
		err := taskErrs[taskCall]
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
	if rawStatus, ok := runner.rawStatusByTask[taskID]; ok {
		if err := setRawTaskStatusForTest(taskPathFromPromptForTest(req.Prompt, runner.gitRoot, taskCycleSlug, taskID), rawStatus); err != nil {
			return agent.ExecuteResult{}, err
		}
	}
	if status, ok := runner.statusByTask[taskID]; ok {
		if err := spec.SetStatus(taskPathFromPromptForTest(req.Prompt, runner.gitRoot, taskCycleSlug, taskID), status); err != nil {
			return agent.ExecuteResult{}, err
		}
	}
	if statuses := runner.statusByTaskCall[taskID]; taskCall < len(statuses) {
		if err := spec.SetStatus(taskPathFromPromptForTest(req.Prompt, runner.gitRoot, taskCycleSlug, taskID), statuses[taskCall]); err != nil {
			return agent.ExecuteResult{}, err
		}
	}
	if result := runner.resultByTask[taskID]; result != "" {
		path := taskPathFromPromptForTest(req.Prompt, runner.gitRoot, taskCycleSlug, taskID)
		file, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY, 0)
		if err != nil {
			return agent.ExecuteResult{}, err
		}
		if _, err := file.WriteString(result); err != nil {
			_ = file.Close()
			return agent.ExecuteResult{}, err
		}
		if err := file.Close(); err != nil {
			return agent.ExecuteResult{}, err
		}
	}
	if runner.afterTask != nil {
		runner.afterTask(taskID)
	}
	return agent.ExecuteResult{LogPath: req.LogPath, TransportAnomaly: runner.anomalyByTask[taskID]}, nil
}

func (runner *taskFakeRunner) EndSession(context.Context, agent.RuntimeSpec, agent.SessionRef) error {
	return nil
}

func taskIDFromPrompt(prompt string) string {
	for _, line := range strings.Split(prompt, "\n") {
		if strings.HasPrefix(line, "Task: ") {
			return strings.TrimPrefix(line, "Task: ")
		}
		if strings.HasPrefix(line, "Work Item: ") {
			value := strings.TrimSpace(strings.TrimPrefix(line, "Work Item: "))
			if strings.HasPrefix(value, "task_") {
				return value
			}
		}
	}
	return ""
}

func taskPathFromPromptForTest(prompt string, gitRoot string, slug string, taskID string) string {
	for _, line := range strings.Split(prompt, "\n") {
		path, ok := strings.CutPrefix(line, "Task file: ")
		if !ok {
			continue
		}
		path = strings.TrimSpace(path)
		if filepath.IsAbs(path) {
			return path
		}
		return filepath.Join(gitRoot, path)
	}
	return taskPathFor(gitRoot, slug, taskID)
}

func qaSpecDirFromPromptForTest(prompt string, gitRoot string) string {
	for _, line := range strings.Split(prompt, "\n") {
		path, ok := strings.CutPrefix(line, "Spec directory: ")
		if !ok {
			continue
		}
		path = strings.TrimSpace(path)
		if filepath.IsAbs(path) {
			return path
		}
		return filepath.Join(gitRoot, path)
	}
	return filepath.Join(gitRoot, "docs", "specs", taskCycleSlug)
}

func setRawTaskStatusForTest(path string, status string) error {
	content, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	lines := strings.Split(string(content), "\n")
	for index, line := range lines {
		if strings.HasPrefix(line, "status:") {
			lines[index] = "status: " + status
			return os.WriteFile(path, []byte(strings.Join(lines, "\n")), 0o644)
		}
	}
	return fmt.Errorf("task file %q has no status line", path)
}

const qaReportNameForTest = "qa-report-2026-01-01.md"

func qaReportRelPathForTest() string {
	return filepath.Join("docs", "specs", taskCycleSlug, "qa", qaReportNameForTest)
}

func qaReportForTest(verdict string) string {
	return fmt.Sprintf("---\nverdict: %s\n---\n\n# QA Report\n", verdict)
}

func TestPerWorkAgentSessionMixedTaskTypesAndQA(t *testing.T) {
	fixture := newTaskCycleFixture(t, []taskSpecSeed{
		{id: "task_01", title: "Backend work", taskType: string(spec.TaskTypeBackend)},
		{id: "task_02", title: "Frontend work", taskType: string(spec.TaskTypeFrontend)},
	})
	runner := &selectionLifecycleRunner{
		gitRoot: fixture.gitRoot,
		statusByTask: map[string]spec.Status{
			"task_01": spec.StatusCompleted,
			"task_02": spec.StatusCompleted,
		},
		qaReport: qaReportForTest(spec.VerdictPass),
	}
	engine := fixture.engine(t, runner, &taskFakeVerifier{calls: fixture.calls}, &engineFakeCommitter{calls: fixture.calls}, fixture.worktree)
	plan := fixture.qaPlan()
	plan.AgentSelections = selectionProfilesForTest(map[roundconfig.WorkCategory]roundconfig.AgentSelectionProfile{
		roundconfig.CategoryBackend:  selectionProfileForTest(selectionForTest("codex", "backend-model", "high"), selectionForTest("codex", "backend-fallback", "high")),
		roundconfig.CategoryFrontend: selectionProfileForTest(selectionForTest("claude", "frontend-model", "medium"), selectionForTest("codex", "frontend-fallback", "high")),
		roundconfig.CategoryQA:       selectionProfileForTest(selectionForTest("codex", "qa-model", "high"), selectionForTest("codex", "qa-fallback", "high")),
	})
	plan.RuntimeFactory = runtimeFactoryForLifecycleTest(nil)

	result, err := engine.TaskCycle(context.Background(), plan)

	if err != nil {
		t.Fatalf("TaskCycle returned error: %v", err)
	}
	if result.Completed != 2 || result.QAVerdict != spec.VerdictPass {
		t.Fatalf("expected two completed Tasks and passing QA, got %+v", result)
	}
	requests := runner.runRequests()
	if got, want := lifecycleRequestSummary(requests), []string{
		"roundfix-" + fixture.run.ID + "-task_01|codex|backend-model",
		"roundfix-" + fixture.run.ID + "-task_02|claude|frontend-model",
		"roundfix-" + fixture.run.ID + "-qa|codex|qa-model",
	}; strings.Join(got, "\n") != strings.Join(want, "\n") {
		t.Fatalf("expected per-work run selections %v, got %v", want, got)
	}
	if got := runner.closedSessions(); strings.Join(got, "\n") != strings.Join([]string{
		"roundfix-" + fixture.run.ID + "-task_01",
		"roundfix-" + fixture.run.ID + "-task_02",
		"roundfix-" + fixture.run.ID + "-qa",
	}, "\n") {
		t.Fatalf("expected every owned session closed once, got %v", got)
	}
	if got := countAgentStatusEvents(fixture.sink, agent.AgentWorkStartedStatus); got != 3 {
		t.Fatalf("expected one agent_work_started marker per work action, got %d", got)
	}
}

func TestAgentSelectionFallbackPublishesBeforeNextSession(t *testing.T) {
	fixture := newTaskCycleFixture(t, []taskSpecSeed{{id: "task_01", taskType: string(spec.TaskTypeBackend)}})
	runner := &selectionLifecycleRunner{
		gitRoot:           fixture.gitRoot,
		statusByTask:      map[string]spec.Status{"task_01": spec.StatusCompleted},
		prepareErrByModel: map[string]error{"bad-model": selectionStartErrForTest("codex", "bad-model")},
		sink:              fixture.sink,
		progress:          fixture.progress,
	}
	engine := fixture.engine(t, runner, &taskFakeVerifier{calls: fixture.calls}, &engineFakeCommitter{calls: fixture.calls}, fixture.worktree)
	plan := fixture.plan()
	plan.AgentSelections = selectionProfilesForTest(map[roundconfig.WorkCategory]roundconfig.AgentSelectionProfile{
		roundconfig.CategoryBackend: selectionProfileForTest(selectionForTest("codex", "bad-model", "high"), selectionForTest("codex", "good-model", "high")),
	})
	plan.RuntimeFactory = runtimeFactoryForLifecycleTest(nil)

	result, err := engine.TaskCycle(context.Background(), plan)

	if err != nil {
		t.Fatalf("TaskCycle returned error: %v", err)
	}
	if result.Completed != 1 {
		t.Fatalf("expected fallback Task to complete, got %+v", result)
	}
	if !runner.fallbackPreparedAfterNotification() {
		t.Fatal("expected fallback notification event to be committed before fallback session preparation")
	}
	if !runner.fallbackPreparedAfterVisibleMessage() {
		t.Fatal("expected caller-visible fallback message before fallback session preparation")
	}
	event := singleEventOfKind(t, fixture.sink, runevent.KindDaemonAgentSelectionFallback)
	payload := eventPayloadMap(t, event)
	if payload["category"] != "backend" || payload["scope_kind"] != "task" || payload["scope_id"] != "task_01" || payload["automatic"] != true {
		t.Fatalf("unexpected fallback payload: %+v", payload)
	}
	if payload["fallback_index"].(float64) != 1 {
		t.Fatalf("expected fallback index 1, got %+v", payload)
	}
	if got := lifecycleRequestSummary(runner.runRequests()); strings.Join(got, "\n") != "roundfix-"+fixture.run.ID+"-task_01-fallback-01|codex|good-model" {
		t.Fatalf("expected only fallback prompt to run, got %v", got)
	}
}

func TestCrossRuntimeFallbackUsesRuntimeFactory(t *testing.T) {
	fixture := newTaskCycleFixture(t, []taskSpecSeed{{id: "task_01", taskType: string(spec.TaskTypeBackend)}})
	seenSelections := []string{}
	runner := &selectionLifecycleRunner{
		gitRoot:           fixture.gitRoot,
		statusByTask:      map[string]spec.Status{"task_01": spec.StatusCompleted},
		prepareErrByModel: map[string]error{"codex-bad": selectionStartErrForTest("codex", "codex-bad")},
	}
	engine := fixture.engine(t, runner, &taskFakeVerifier{calls: fixture.calls}, &engineFakeCommitter{calls: fixture.calls}, fixture.worktree)
	plan := fixture.plan()
	plan.AgentSelections = selectionProfilesForTest(map[roundconfig.WorkCategory]roundconfig.AgentSelectionProfile{
		roundconfig.CategoryBackend: selectionProfileForTest(selectionForTest("codex", "codex-bad", "high"), selectionForTest("claude", "claude-good", "medium")),
	})
	plan.RuntimeFactory = runtimeFactoryForLifecycleTest(&seenSelections)

	_, err := engine.TaskCycle(context.Background(), plan)

	if err != nil {
		t.Fatalf("TaskCycle returned error: %v", err)
	}
	if got, want := strings.Join(seenSelections, ","), "codex/codex-bad/high,claude/claude-good/medium"; got != want {
		t.Fatalf("expected runtime factory to receive preferred then cross-runtime fallback, got %q", got)
	}
	if got := lifecycleRequestSummary(runner.runRequests()); strings.Join(got, "\n") != "roundfix-"+fixture.run.ID+"-task_01-fallback-01|claude|claude-good" {
		t.Fatalf("expected Claude fallback prompt, got %v", got)
	}
}

func TestNoFallbackAfterAgentWorkStarted(t *testing.T) {
	fixture := newTaskCycleFixture(t, []taskSpecSeed{{id: "task_01", taskType: string(spec.TaskTypeBackend)}})
	runner := &selectionLifecycleRunner{
		gitRoot:       fixture.gitRoot,
		runErrByModel: map[string]error{"preferred-model": selectionStartErrForTest("codex", "preferred-model")},
	}
	engine := fixture.engine(t, runner, &taskFakeVerifier{calls: fixture.calls}, &engineFakeCommitter{calls: fixture.calls}, fixture.worktree)
	plan := fixture.plan()
	plan.AgentSelections = selectionProfilesForTest(map[roundconfig.WorkCategory]roundconfig.AgentSelectionProfile{
		roundconfig.CategoryBackend: selectionProfileForTest(selectionForTest("codex", "preferred-model", "high"), selectionForTest("codex", "fallback-model", "high")),
	})
	plan.RuntimeFactory = runtimeFactoryForLifecycleTest(nil)

	result, err := engine.TaskCycle(context.Background(), plan)

	if err != nil {
		t.Fatalf("TaskCycle returned error: %v", err)
	}
	if result.Failed != 1 {
		t.Fatalf("expected post-start selection-looking error to fail Task, got %+v", result)
	}
	if got := lifecycleRequestSummary(runner.prepareRequests()); strings.Join(got, "\n") != "roundfix-"+fixture.run.ID+"-task_01|codex|preferred-model" {
		t.Fatalf("expected no fallback preparation after agent_work_started, got %v", got)
	}
	if events := eventsOfKind(fixture.sink, runevent.KindDaemonAgentSelectionFallback); len(events) != 0 {
		t.Fatalf("expected no fallback event after work started, got %+v", events)
	}
}

func TestAgentSessionOwnerCleanup(t *testing.T) {
	fixture := newTaskCycleFixture(t, []taskSpecSeed{{id: "task_01", taskType: string(spec.TaskTypeBackend)}})
	runner := &selectionLifecycleRunner{
		gitRoot: fixture.gitRoot,
		prepareErrByModel: map[string]error{
			"bad-preferred": selectionStartErrForTest("codex", "bad-preferred"),
			"bad-fallback":  selectionStartErrForTest("claude", "bad-fallback"),
		},
	}
	engine := fixture.engine(t, runner, &taskFakeVerifier{calls: fixture.calls}, &engineFakeCommitter{calls: fixture.calls}, fixture.worktree)
	plan := fixture.plan()
	plan.AgentSelections = selectionProfilesForTest(map[roundconfig.WorkCategory]roundconfig.AgentSelectionProfile{
		roundconfig.CategoryBackend: selectionProfileForTest(selectionForTest("codex", "bad-preferred", "high"), selectionForTest("claude", "bad-fallback", "medium")),
	})
	plan.RuntimeFactory = runtimeFactoryForLifecycleTest(nil)

	result, err := engine.TaskCycle(context.Background(), plan)

	if err != nil {
		t.Fatalf("TaskCycle returned error: %v", err)
	}
	if result.Failed != 1 {
		t.Fatalf("expected exhausted selection to settle Task failed, got %+v", result)
	}
	if got := strings.Join(runner.closedSessions(), "\n"); got != strings.Join([]string{
		"roundfix-" + fixture.run.ID + "-task_01",
		"roundfix-" + fixture.run.ID + "-task_01-fallback-01",
	}, "\n") {
		t.Fatalf("expected failed preferred and fallback sessions closed, got %q", got)
	}
	event := singleEventOfKind(t, fixture.sink, runevent.KindDaemonAgentSelectionExhausted)
	payload := eventPayloadMap(t, event)
	attempts, ok := payload["attempts"].([]any)
	if !ok || len(attempts) != 2 {
		t.Fatalf("expected exhausted event to list two attempts, got %+v", payload)
	}
	if !strings.Contains(result.Outcomes[0].Reason, "bad-preferred") || !strings.Contains(result.Outcomes[0].Reason, "bad-fallback") || !strings.Contains(result.Outcomes[0].Reason, "roundfix profiles validate --category backend") {
		t.Fatalf("expected failure reason to list attempts and recovery, got %q", result.Outcomes[0].Reason)
	}
}

func TestAgentSessionOwnerCleanupClosesOnCancellation(t *testing.T) {
	fixture := newTaskCycleFixture(t, []taskSpecSeed{{id: "task_01", taskType: string(spec.TaskTypeBackend)}})
	runner := &selectionLifecycleRunner{
		gitRoot:       fixture.gitRoot,
		runErrByModel: map[string]error{"preferred-model": agent.StopError{Err: context.Canceled}},
	}
	engine := fixture.engine(t, runner, &taskFakeVerifier{calls: fixture.calls}, &engineFakeCommitter{calls: fixture.calls}, fixture.worktree)
	plan := fixture.plan()
	plan.AgentSelections = selectionProfilesForTest(map[roundconfig.WorkCategory]roundconfig.AgentSelectionProfile{
		roundconfig.CategoryBackend: selectionProfileForTest(selectionForTest("codex", "preferred-model", "high"), selectionForTest("codex", "fallback-model", "high")),
	})
	plan.RuntimeFactory = runtimeFactoryForLifecycleTest(nil)

	_, err := engine.TaskCycle(context.Background(), plan)

	if !agent.IsStopError(err) {
		t.Fatalf("expected StopError from canceled Agent work, got %v", err)
	}
	if got := strings.Join(runner.closedSessions(), "\n"); got != "roundfix-"+fixture.run.ID+"-task_01" {
		t.Fatalf("expected active session closed on cancellation, got %q", got)
	}
	if got := lifecycleRequestSummary(runner.prepareRequests()); strings.Join(got, "\n") != "roundfix-"+fixture.run.ID+"-task_01|codex|preferred-model" {
		t.Fatalf("expected cancellation not to prepare fallback, got %v", got)
	}
}

type selectionLifecycleRunner struct {
	mu                                   sync.Mutex
	gitRoot                              string
	statusByTask                         map[string]spec.Status
	qaReport                             string
	prepareErrByModel                    map[string]error
	runErrByModel                        map[string]error
	closeErr                             error
	prepared                             []agent.ExecuteRequest
	ran                                  []agent.ExecuteRequest
	closed                               []agent.SessionRef
	sink                                 *captureEventSink
	progress                             *bytes.Buffer
	sawFallbackNotificationBeforePrepare bool
	sawFallbackMessageBeforePrepare      bool
}

func (runner *selectionLifecycleRunner) Probe(context.Context, agent.ProbeRequest) error { return nil }

func (runner *selectionLifecycleRunner) PrepareSession(_ context.Context, req agent.ExecuteRequest, _ runevent.Sink) error {
	runner.mu.Lock()
	runner.prepared = append(runner.prepared, req)
	if strings.Contains(req.Session.Name, "fallback") {
		runner.sawFallbackNotificationBeforePrepare = runner.sink == nil || len(eventsOfKind(runner.sink, runevent.KindDaemonAgentSelectionFallback)) > 0
		runner.sawFallbackMessageBeforePrepare = runner.progress == nil || strings.Contains(runner.progress.String(), "activating fallback")
	}
	err := runner.prepareErrByModel[req.Runtime.Model]
	runner.mu.Unlock()
	return err
}

func (runner *selectionLifecycleRunner) Run(ctx context.Context, req agent.ExecuteRequest, sink runevent.Sink) (agent.ExecuteResult, error) {
	return runner.RunPrepared(ctx, req, sink)
}

func (runner *selectionLifecycleRunner) RunPrepared(_ context.Context, req agent.ExecuteRequest, _ runevent.Sink) (agent.ExecuteResult, error) {
	runner.mu.Lock()
	runner.ran = append(runner.ran, req)
	err := runner.runErrByModel[req.Runtime.Model]
	runner.mu.Unlock()
	if err != nil {
		return agent.ExecuteResult{LogPath: req.LogPath}, err
	}
	taskID := taskIDFromPrompt(req.Prompt)
	if taskID == "" && strings.Contains(req.Prompt, "Spec QA gate") {
		if runner.qaReport != "" {
			reportPath := filepath.Join(qaSpecDirFromPromptForTest(req.Prompt, runner.gitRoot), "qa", qaReportNameForTest)
			if err := os.MkdirAll(filepath.Dir(reportPath), 0o755); err != nil {
				return agent.ExecuteResult{}, err
			}
			if err := os.WriteFile(reportPath, []byte(runner.qaReport), 0o644); err != nil {
				return agent.ExecuteResult{}, err
			}
		}
		return agent.ExecuteResult{LogPath: req.LogPath}, nil
	}
	if taskID == "" && len(req.Batch.Issues) > 0 {
		for _, issue := range req.Batch.Issues {
			if err := rounds.SetIssueStatus(issue.Path, rounds.StatusResolved, "", ""); err != nil {
				return agent.ExecuteResult{}, err
			}
		}
		return agent.ExecuteResult{LogPath: req.LogPath}, nil
	}
	if status, ok := runner.statusByTask[taskID]; ok {
		if err := spec.SetStatus(taskPathFromPromptForTest(req.Prompt, runner.gitRoot, taskCycleSlug, taskID), status); err != nil {
			return agent.ExecuteResult{}, err
		}
	}
	return agent.ExecuteResult{LogPath: req.LogPath}, nil
}

func (runner *selectionLifecycleRunner) EndSession(_ context.Context, _ agent.RuntimeSpec, session agent.SessionRef) error {
	runner.mu.Lock()
	defer runner.mu.Unlock()
	runner.closed = append(runner.closed, session)
	return runner.closeErr
}

func (runner *selectionLifecycleRunner) prepareRequests() []agent.ExecuteRequest {
	runner.mu.Lock()
	defer runner.mu.Unlock()
	return append([]agent.ExecuteRequest(nil), runner.prepared...)
}

func (runner *selectionLifecycleRunner) runRequests() []agent.ExecuteRequest {
	runner.mu.Lock()
	defer runner.mu.Unlock()
	return append([]agent.ExecuteRequest(nil), runner.ran...)
}

func (runner *selectionLifecycleRunner) closedSessions() []string {
	runner.mu.Lock()
	defer runner.mu.Unlock()
	sessions := make([]string, 0, len(runner.closed))
	for _, session := range runner.closed {
		sessions = append(sessions, session.Name)
	}
	return sessions
}

func (runner *selectionLifecycleRunner) fallbackPreparedAfterNotification() bool {
	runner.mu.Lock()
	defer runner.mu.Unlock()
	return runner.sawFallbackNotificationBeforePrepare
}

func (runner *selectionLifecycleRunner) fallbackPreparedAfterVisibleMessage() bool {
	runner.mu.Lock()
	defer runner.mu.Unlock()
	return runner.sawFallbackMessageBeforePrepare
}

func selectionForTest(runtime string, model string, effort string) roundconfig.AgentSelection {
	return roundconfig.AgentSelection{Runtime: runtime, Model: model, ReasoningEffort: effort}
}

func selectionProfileForTest(preferred roundconfig.AgentSelection, fallbacks ...roundconfig.AgentSelection) roundconfig.AgentSelectionProfile {
	return roundconfig.AgentSelectionProfile{Preferred: preferred, Fallbacks: append([]roundconfig.AgentSelection(nil), fallbacks...)}
}

func selectionProfilesForTest(profiles map[roundconfig.WorkCategory]roundconfig.AgentSelectionProfile) AgentSelectionProfiles {
	resolved := AgentSelectionProfiles{}
	for category, profile := range profiles {
		resolved[category] = roundconfig.ResolvedProfile{
			Category: category,
			Source:   roundconfig.ProfileSourceProject,
			Profile:  profile,
		}
	}
	return resolved
}

func runtimeFactoryForLifecycleTest(seen *[]string) AgentRuntimeFactory {
	return func(selection roundconfig.AgentSelection) (agent.RuntimeSpec, error) {
		if seen != nil {
			*seen = append(*seen, selection.Runtime+"/"+selection.Model+"/"+selection.ReasoningEffort)
		}
		return agent.RuntimeSpec{
			ID:              selection.Runtime,
			DisplayName:     selection.Runtime,
			Model:           selection.Model,
			ReasoningEffort: selection.ReasoningEffort,
		}, nil
	}
}

func selectionStartErrForTest(runtime string, model string) error {
	return &agent.SelectionPreflightError{
		Runtime:         runtime,
		Model:           model,
		ReasoningEffort: "high",
		Operation:       "set model",
		Err:             errors.New("adapter rejected " + model),
	}
}

func lifecycleRequestSummary(requests []agent.ExecuteRequest) []string {
	summary := make([]string, 0, len(requests))
	for _, req := range requests {
		summary = append(summary, req.Session.Name+"|"+req.Runtime.ID+"|"+req.Runtime.Model)
	}
	return summary
}

func countAgentStatusEvents(sink *captureEventSink, status string) int {
	count := 0
	for _, event := range eventsOfKind(sink, runevent.KindAgentStatus) {
		var payload struct {
			Status string `json:"status"`
		}
		if json.Unmarshal(event.Payload, &payload) == nil && payload.Status == status {
			count++
		}
	}
	return count
}

func singleEventOfKind(t *testing.T, sink *captureEventSink, kind runevent.Kind) runevent.RunEvent {
	t.Helper()
	events := eventsOfKind(sink, kind)
	if len(events) != 1 {
		t.Fatalf("expected one %s event, got %d: %+v", kind, len(events), events)
	}
	return events[0]
}

func eventsOfKind(sink *captureEventSink, kind runevent.Kind) []runevent.RunEvent {
	if sink == nil {
		return nil
	}
	var events []runevent.RunEvent
	for _, event := range sink.snapshot() {
		if event.Kind == kind {
			events = append(events, event)
		}
	}
	return events
}

// taskFakeVerifier records every verification command verbatim and fails
// the commands scripted in failOn.
type taskFakeVerifier struct {
	calls           *[]string
	store           *store.Store
	runID           string
	failOn          map[string]error
	script          []error
	temporaryOnCall map[int]bool
	commands        []string
	workDirs        []string
	outputPaths     []string
	seenStates      []string
}

func (verifier *taskFakeVerifier) Verify(_ context.Context, req VerifyRequest) (VerifyResult, error) {
	*verifier.calls = append(*verifier.calls, "verify")
	verifier.commands = append(verifier.commands, req.Command)
	verifier.workDirs = append(verifier.workDirs, req.WorkDir)
	verifier.outputPaths = append(verifier.outputPaths, req.OutputPath)
	if verifier.store != nil {
		verifier.seenStates = append(verifier.seenStates, runStateForTest(verifier.store, verifier.runID))
	}
	if verifier.temporaryOnCall[len(verifier.commands)] {
		commandErr := &VerificationCommandError{Command: req.Command, OutputPath: req.OutputPath, Err: errors.New("exit status 75")}
		return VerifyResult{OutputPath: req.OutputPath}, &TemporaryVerificationFailureError{CommandFailure: commandErr}
	}
	if len(verifier.script) > 0 {
		err := verifier.script[0]
		verifier.script = verifier.script[1:]
		if err != nil {
			return VerifyResult{OutputPath: req.OutputPath}, &VerificationCommandError{Command: req.Command, OutputPath: req.OutputPath, Err: err}
		}
		return VerifyResult{OutputPath: req.OutputPath}, nil
	}
	if err := verifier.failOn[req.Command]; err != nil {
		return VerifyResult{OutputPath: req.OutputPath}, &VerificationCommandError{Command: req.Command, OutputPath: req.OutputPath, Err: err}
	}
	return VerifyResult{OutputPath: req.OutputPath}, nil
}

type fakePriorChangedResolver struct {
	mu       sync.Mutex
	byWork   map[string][]string
	calls    []priorChangedCall
	resolve  func(workDir string) []string
	errByDir map[string]error
}

type priorChangedCall struct {
	workDir     string
	initialHead string
}

func (resolver *fakePriorChangedResolver) PriorChangedFiles(_ context.Context, workDir string, initialHead string) ([]string, error) {
	resolver.mu.Lock()
	defer resolver.mu.Unlock()
	resolver.calls = append(resolver.calls, priorChangedCall{workDir: workDir, initialHead: initialHead})
	if err := resolver.errByDir[workDir]; err != nil {
		return nil, err
	}
	if resolver.resolve != nil {
		return append([]string(nil), resolver.resolve(workDir)...), nil
	}
	return append([]string(nil), resolver.byWork[workDir]...), nil
}

func (resolver *fakePriorChangedResolver) seenCalls() []priorChangedCall {
	resolver.mu.Lock()
	defer resolver.mu.Unlock()
	return append([]priorChangedCall(nil), resolver.calls...)
}

type taskSchedulerRunner struct {
	mu           sync.Mutex
	started      chan string
	release      map[string]chan struct{}
	statusByTask map[string]spec.Status
	errByTask    map[string]error
	onStart      func(taskID string, req agent.ExecuteRequest) error
	active       int
	maxActive    int
	starts       []string
	requests     []agent.ExecuteRequest
}

func newTaskSchedulerRunner(taskIDs ...string) *taskSchedulerRunner {
	runner := &taskSchedulerRunner{
		started: make(chan string, len(taskIDs)),
		release: make(map[string]chan struct{}, len(taskIDs)),
	}
	for _, taskID := range taskIDs {
		runner.release[taskID] = make(chan struct{})
	}
	return runner
}

func (runner *taskSchedulerRunner) Probe(context.Context, agent.ProbeRequest) error { return nil }

func (runner *taskSchedulerRunner) Run(ctx context.Context, req agent.ExecuteRequest, _ runevent.Sink) (agent.ExecuteResult, error) {
	taskID := taskIDFromPrompt(req.Prompt)
	runner.mu.Lock()
	runner.active++
	if runner.active > runner.maxActive {
		runner.maxActive = runner.active
	}
	runner.starts = append(runner.starts, taskID)
	runner.requests = append(runner.requests, req)
	runner.mu.Unlock()
	runner.started <- taskID
	defer func() {
		runner.mu.Lock()
		runner.active--
		runner.mu.Unlock()
	}()

	if runner.onStart != nil {
		if err := runner.onStart(taskID, req); err != nil {
			return agent.ExecuteResult{}, err
		}
	}
	if release := runner.release[taskID]; release != nil {
		select {
		case <-release:
		case <-ctx.Done():
			return agent.ExecuteResult{}, agent.StopError{Err: ctx.Err()}
		}
	}
	if err := runner.errByTask[taskID]; err != nil {
		return agent.ExecuteResult{}, err
	}
	if status, ok := runner.statusByTask[taskID]; ok {
		if err := spec.SetStatus(taskPathFor(req.GitRoot, taskCycleSlug, taskID), status); err != nil {
			return agent.ExecuteResult{}, err
		}
	}
	return agent.ExecuteResult{LogPath: req.LogPath}, nil
}

func (runner *taskSchedulerRunner) EndSession(context.Context, agent.RuntimeSpec, agent.SessionRef) error {
	return nil
}

func (runner *taskSchedulerRunner) releaseTask(taskID string) {
	if release := runner.release[taskID]; release != nil {
		close(release)
	}
}

func (runner *taskSchedulerRunner) maxObservedActive() int {
	runner.mu.Lock()
	defer runner.mu.Unlock()
	return runner.maxActive
}

func (runner *taskSchedulerRunner) startedTasks() []string {
	runner.mu.Lock()
	defer runner.mu.Unlock()
	return append([]string(nil), runner.starts...)
}

// qaPrompts reports the QA gate prompts the Agent received, which is empty
// whenever the QA step never started.
func (runner *taskSchedulerRunner) qaPrompts() []string {
	runner.mu.Lock()
	defer runner.mu.Unlock()
	var prompts []string
	for _, req := range runner.requests {
		if strings.Contains(req.Prompt, "Spec QA gate") {
			prompts = append(prompts, req.Prompt)
		}
	}
	return prompts
}

func waitSchedulerStarts(t *testing.T, runner *taskSchedulerRunner, count int) []string {
	t.Helper()
	started := make([]string, 0, count)
	for len(started) < count {
		select {
		case taskID := <-runner.started:
			started = append(started, taskID)
		case <-time.After(2 * time.Second):
			t.Fatalf("timed out waiting for %d Task start(s), got %v", count, started)
		}
	}
	return started
}

func assertNoSchedulerStart(t *testing.T, runner *taskSchedulerRunner) {
	t.Helper()
	select {
	case taskID := <-runner.started:
		t.Fatalf("expected no additional Task to start, got %s", taskID)
	default:
	}
}

func assertTaskSet(t *testing.T, got []string, want ...string) {
	t.Helper()
	sort.Strings(got)
	sort.Strings(want)
	if strings.Join(got, "|") != strings.Join(want, "|") {
		t.Fatalf("expected started Tasks %v, got %v", want, got)
	}
}

type taskSchedulerVerifier struct {
	mu       sync.Mutex
	failOn   map[string]error
	commands []string
	workDirs []string
}

func (verifier *taskSchedulerVerifier) Verify(_ context.Context, req VerifyRequest) (VerifyResult, error) {
	verifier.mu.Lock()
	verifier.commands = append(verifier.commands, req.Command)
	verifier.workDirs = append(verifier.workDirs, req.WorkDir)
	verifier.mu.Unlock()
	if err := verifier.failOn[req.Command]; err != nil {
		return VerifyResult{OutputPath: req.OutputPath}, &VerificationCommandError{Command: req.Command, OutputPath: req.OutputPath, Err: err}
	}
	return VerifyResult{OutputPath: req.OutputPath}, nil
}

type taskVerificationStart struct {
	taskID  string
	attempt int
}

type taskCapacityVerifier struct {
	mu          sync.Mutex
	started     chan taskVerificationStart
	release     map[string]chan struct{}
	fail        map[string]error
	calls       map[string]int
	active      int
	maxActive   int
	totalStarts int
}

func newTaskCapacityVerifier(taskIDs ...string) *taskCapacityVerifier {
	verifier := &taskCapacityVerifier{
		started: make(chan taskVerificationStart, len(taskIDs)*2),
		release: make(map[string]chan struct{}, len(taskIDs)*2),
		fail:    map[string]error{},
		calls:   map[string]int{},
	}
	for _, taskID := range taskIDs {
		verifier.release[taskVerificationKey(taskID, 1)] = make(chan struct{})
	}
	return verifier
}

func (verifier *taskCapacityVerifier) Verify(ctx context.Context, req VerifyRequest) (VerifyResult, error) {
	taskID := strings.TrimPrefix(req.Command, "verify-")
	verifier.mu.Lock()
	verifier.calls[taskID]++
	attempt := verifier.calls[taskID]
	key := taskVerificationKey(taskID, attempt)
	verifier.active++
	verifier.totalStarts++
	if verifier.active > verifier.maxActive {
		verifier.maxActive = verifier.active
	}
	release := verifier.release[key]
	fail := verifier.fail[key]
	verifier.mu.Unlock()

	verifier.started <- taskVerificationStart{taskID: taskID, attempt: attempt}
	defer func() {
		verifier.mu.Lock()
		verifier.active--
		verifier.mu.Unlock()
	}()

	if fail != nil {
		return VerifyResult{OutputPath: req.OutputPath}, &VerificationCommandError{
			Command:    req.Command,
			OutputPath: req.OutputPath,
			Err:        fail,
		}
	}
	if release != nil {
		select {
		case <-release:
		case <-ctx.Done():
			return VerifyResult{OutputPath: req.OutputPath}, ctx.Err()
		}
	}
	return VerifyResult{OutputPath: req.OutputPath}, nil
}

func (verifier *taskCapacityVerifier) addAttempt(taskID string, attempt int) {
	verifier.mu.Lock()
	defer verifier.mu.Unlock()
	verifier.release[taskVerificationKey(taskID, attempt)] = make(chan struct{})
}

func (verifier *taskCapacityVerifier) failAttempt(taskID string, attempt int, err error) {
	verifier.mu.Lock()
	defer verifier.mu.Unlock()
	verifier.fail[taskVerificationKey(taskID, attempt)] = err
}

func (verifier *taskCapacityVerifier) releaseAttempt(taskID string, attempt int) {
	verifier.mu.Lock()
	release := verifier.release[taskVerificationKey(taskID, attempt)]
	verifier.mu.Unlock()
	if release != nil {
		close(release)
	}
}

func (verifier *taskCapacityVerifier) waitStart(t *testing.T) taskVerificationStart {
	t.Helper()
	select {
	case started := <-verifier.started:
		return started
	case <-time.After(2 * time.Second):
		t.Fatal("timed out waiting for Verification start")
		return taskVerificationStart{}
	}
}

func (verifier *taskCapacityVerifier) assertNoStart(t *testing.T) {
	t.Helper()
	select {
	case started := <-verifier.started:
		t.Fatalf("expected no Verification start, got %+v", started)
	default:
	}
}

func (verifier *taskCapacityVerifier) observed() (maxActive int, totalStarts int) {
	verifier.mu.Lock()
	defer verifier.mu.Unlock()
	return verifier.maxActive, verifier.totalStarts
}

func taskVerificationKey(taskID string, attempt int) string {
	return fmt.Sprintf("%s:%d", taskID, attempt)
}

type taskSchedulerCommitter struct {
	mu       sync.Mutex
	messages []string
	paths    [][]string
}

func (committer *taskSchedulerCommitter) Commit(_ context.Context, req CommitRequest) error {
	committer.mu.Lock()
	defer committer.mu.Unlock()
	committer.messages = append(committer.messages, req.Message)
	committer.paths = append(committer.paths, append([]string(nil), req.Paths...))
	return nil
}

func (committer *taskSchedulerCommitter) commitCount() int {
	committer.mu.Lock()
	defer committer.mu.Unlock()
	return len(committer.messages)
}

func (committer *taskSchedulerCommitter) commitMessages() []string {
	committer.mu.Lock()
	defer committer.mu.Unlock()
	return append([]string(nil), committer.messages...)
}

type fakeTaskWorktrees struct {
	mu               sync.Mutex
	fsMu             sync.Mutex
	conflictByTask   map[string]string
	createErrByTask  map[string]error
	onCreate         func(taskID string, ref runworktree.TaskRef, opts runworktree.TaskCreateOptions) error
	onIntegrate      func(taskID string)
	created          map[string]runworktree.TaskRef
	createOptions    map[string]runworktree.TaskCreateOptions
	integrated       []string
	cleaned          []string
	integratedSignal chan string
}

func newFakeTaskWorktrees() *fakeTaskWorktrees {
	return &fakeTaskWorktrees{
		created:          map[string]runworktree.TaskRef{},
		createOptions:    map[string]runworktree.TaskCreateOptions{},
		integratedSignal: make(chan string, 16),
	}
}

func (worktrees *fakeTaskWorktrees) CreateTask(_ context.Context, run runworktree.Ref, taskID string, opts runworktree.TaskCreateOptions) (runworktree.TaskRef, error) {
	path := filepath.Join(filepath.Dir(run.Path), run.RunID+"."+taskID)
	worktrees.fsMu.Lock()
	defer worktrees.fsMu.Unlock()
	if err := os.RemoveAll(path); err != nil {
		return runworktree.TaskRef{}, err
	}
	if err := copyTreeForSchedulerTest(run.Path, path); err != nil {
		return runworktree.TaskRef{}, err
	}
	ref := runworktree.TaskRef{
		RunID:    run.RunID,
		TaskID:   taskID,
		Path:     path,
		Branch:   runworktree.TaskBranchName(run.RunID, taskID),
		UserRoot: run.UserRoot,
		BaseSHA:  "base-" + taskID,
	}
	optsCopy := opts
	optsCopy.CopyList = append([]string(nil), opts.CopyList...)
	worktrees.mu.Lock()
	worktrees.created[taskID] = ref
	worktrees.createOptions[taskID] = optsCopy
	createErr := worktrees.createErrByTask[taskID]
	onCreate := worktrees.onCreate
	worktrees.mu.Unlock()
	if onCreate != nil {
		if err := onCreate(taskID, ref, opts); err != nil {
			return ref, err
		}
	}
	if createErr != nil {
		return ref, createErr
	}
	return ref, nil
}

func (worktrees *fakeTaskWorktrees) IntegrateTask(_ context.Context, run runworktree.Ref, task runworktree.TaskRef) (runworktree.TaskIntegration, error) {
	worktrees.mu.Lock()
	worktrees.integrated = append(worktrees.integrated, task.TaskID)
	conflict := worktrees.conflictByTask[task.TaskID]
	onIntegrate := worktrees.onIntegrate
	worktrees.mu.Unlock()
	worktrees.integratedSignal <- task.TaskID
	if onIntegrate != nil {
		onIntegrate(task.TaskID)
	}
	if conflict != "" {
		return runworktree.TaskIntegration{Mode: runworktree.ModeTaskConflict, Reason: conflict}, nil
	}
	worktrees.fsMu.Lock()
	defer worktrees.fsMu.Unlock()
	if err := copyTreeForSchedulerTest(task.Path, run.Path); err != nil {
		return runworktree.TaskIntegration{}, err
	}
	return runworktree.TaskIntegration{Mode: runworktree.ModeTaskCherryPick}, nil
}

func (worktrees *fakeTaskWorktrees) CleanupTask(_ context.Context, task runworktree.TaskRef) error {
	worktrees.mu.Lock()
	worktrees.cleaned = append(worktrees.cleaned, task.TaskID)
	worktrees.mu.Unlock()
	return os.RemoveAll(task.Path)
}

func (worktrees *fakeTaskWorktrees) integratedTasks() []string {
	worktrees.mu.Lock()
	defer worktrees.mu.Unlock()
	return append([]string(nil), worktrees.integrated...)
}

func (worktrees *fakeTaskWorktrees) cleanedTasks() []string {
	worktrees.mu.Lock()
	defer worktrees.mu.Unlock()
	return append([]string(nil), worktrees.cleaned...)
}

func (worktrees *fakeTaskWorktrees) taskRef(taskID string) runworktree.TaskRef {
	worktrees.mu.Lock()
	defer worktrees.mu.Unlock()
	return worktrees.created[taskID]
}

func (worktrees *fakeTaskWorktrees) taskCreateOptions(taskID string) runworktree.TaskCreateOptions {
	worktrees.mu.Lock()
	defer worktrees.mu.Unlock()
	return worktrees.createOptions[taskID]
}

func (worktrees *fakeTaskWorktrees) integratedTask(taskID string) bool {
	worktrees.mu.Lock()
	defer worktrees.mu.Unlock()
	for _, integrated := range worktrees.integrated {
		if integrated == taskID {
			return true
		}
	}
	return false
}

func waitIntegratedTask(t *testing.T, worktrees *fakeTaskWorktrees) string {
	t.Helper()
	select {
	case taskID := <-worktrees.integratedSignal:
		return taskID
	case <-time.After(2 * time.Second):
		t.Fatal("timed out waiting for Task integration")
		return ""
	}
}

func waitTaskCycleResult(t *testing.T, resultCh <-chan struct {
	result TaskCycleResult
	err    error
}) struct {
	result TaskCycleResult
	err    error
} {
	t.Helper()
	select {
	case outcome := <-resultCh:
		return outcome
	case <-time.After(2 * time.Second):
		t.Fatal("timed out waiting for TaskCycle result")
		return struct {
			result TaskCycleResult
			err    error
		}{}
	}
}

func taskStartedEvents(t *testing.T, sink *captureEventSink) []string {
	t.Helper()
	var started []string
	for _, event := range taskEventsOfKind(sink, runevent.KindDaemonTask) {
		if eventPayloadString(t, event, "phase") == "started" {
			started = append(started, fmt.Sprintf("%s:%d", event.ReviewIssue, event.Batch))
		}
	}
	return started
}

func droppedStageEvents(t *testing.T, sink *captureEventSink) []runevent.RunEvent {
	t.Helper()
	var dropped []runevent.RunEvent
	for _, event := range taskEventsOfKind(sink, runevent.KindDaemonCommit) {
		if eventPayloadString(t, event, "decision") == "dropped" {
			dropped = append(dropped, event)
		}
	}
	return dropped
}

func noOpTaskCommitWarningEvents(t *testing.T, sink *captureEventSink) []runevent.RunEvent {
	t.Helper()
	var warnings []runevent.RunEvent
	for _, event := range taskEventsOfKind(sink, runevent.KindDaemonCommit) {
		if eventPayloadString(t, event, "decision") == "warning" && eventPayloadString(t, event, "warning") == "no_op_task_commit" {
			warnings = append(warnings, event)
		}
	}
	return warnings
}

func TestPublishNoOpTaskCommitWarningReturnsProgressWriteError(t *testing.T) {
	engine := &Engine{deps: Dependencies{Progress: failingProgressWriter{err: errors.New("progress closed")}}}

	err := engine.publishNoOpTaskCommitWarning(context.Background(), TaskPlan{RunID: "run_1"}, "task_01", 1, "spec-only")

	if err == nil {
		t.Fatal("expected progress write error")
	}
	for _, want := range []string{"write no-op Task commit warning", "run_1", "task_01", "progress closed"} {
		if !strings.Contains(err.Error(), want) {
			t.Fatalf("expected error to contain %q, got %v", want, err)
		}
	}
}

type failingProgressWriter struct {
	err error
}

func (writer failingProgressWriter) Write([]byte) (int, error) {
	return 0, writer.err
}

func assertNoOpTaskCommitWarning(t *testing.T, fixture *taskCycleFixture, taskID string, shape string) {
	t.Helper()
	warnings := noOpTaskCommitWarningEvents(t, fixture.sink)
	if len(warnings) != 1 {
		t.Fatalf("expected one no-op Task commit warning event, got %+v", warnings)
	}
	if warnings[0].ReviewIssue != taskID {
		t.Fatalf("expected warning event for %s, got %+v", taskID, warnings[0])
	}
	if got := eventPayloadString(t, warnings[0], "task"); got != taskID {
		t.Fatalf("expected warning payload task %q, got %q", taskID, got)
	}
	if got := eventPayloadString(t, warnings[0], "shape"); got != shape {
		t.Fatalf("expected warning shape %q, got %q", shape, got)
	}
	wantLine := fmt.Sprintf("roundfix: warning: Task %s completed with no changes outside the Spec Root (%s)\n", taskID, shape)
	if count := strings.Count(fixture.progress.String(), wantLine); count != 1 {
		t.Fatalf("expected one warning line %q, count=%d progress=%q", wantLine, count, fixture.progress.String())
	}
}

func assertNoNoOpTaskCommitWarning(t *testing.T, fixture *taskCycleFixture) {
	t.Helper()
	if warnings := noOpTaskCommitWarningEvents(t, fixture.sink); len(warnings) != 0 {
		t.Fatalf("expected no no-op Task commit warning events, got %+v", warnings)
	}
	if strings.Contains(fixture.progress.String(), "no changes outside the Spec Root") {
		t.Fatalf("expected no no-op warning line, got progress=%q", fixture.progress.String())
	}
}

func hasTaskSettlementEvent(t *testing.T, events []runevent.RunEvent, taskID string, status spec.Status) bool {
	t.Helper()
	for _, event := range events {
		if event.ReviewIssue != taskID {
			continue
		}
		if eventPayloadString(t, event, "phase") == "settled" && eventPayloadString(t, event, "status") == string(status) {
			return true
		}
	}
	return false
}

func taskOutcomeByID(outcomes []TaskOutcome, taskID string) (TaskOutcome, bool) {
	for _, outcome := range outcomes {
		if outcome.Task == taskID {
			return outcome, true
		}
	}
	return TaskOutcome{}, false
}

func copyTreeForSchedulerTest(source string, destination string) error {
	if err := os.MkdirAll(destination, 0o755); err != nil {
		return err
	}
	return filepath.WalkDir(source, func(path string, entry os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(source, path)
		if err != nil {
			return err
		}
		if rel == "." {
			return nil
		}
		target := filepath.Join(destination, rel)
		if entry.IsDir() {
			return os.MkdirAll(target, 0o755)
		}
		info, err := entry.Info()
		if err != nil {
			return err
		}
		if !info.Mode().IsRegular() {
			return nil
		}
		content, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
			return err
		}
		return os.WriteFile(target, content, info.Mode().Perm())
	})
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

func assertTaskAnomalyBeforeVerification(t *testing.T, events []runevent.RunEvent, taskID string, anomaly string) {
	t.Helper()
	anomalyIndex := -1
	verificationIndex := -1
	for index, event := range events {
		if event.Kind == runevent.KindDaemonTask && event.ReviewIssue == taskID && eventPayloadString(t, event, "anomaly") == anomaly {
			anomalyIndex = index
			if event.Source != runevent.SourceDaemon {
				t.Fatalf("expected daemon Task anomaly event, got %+v", event)
			}
		}
		if event.Kind == runevent.KindDaemonVerification && event.ReviewIssue == taskID && strings.Contains(string(event.Payload), `"phase":"started"`) {
			verificationIndex = index
		}
	}
	if anomalyIndex < 0 {
		t.Fatalf("expected transport anomaly event for %s, got %+v", taskID, events)
	}
	if verificationIndex < 0 {
		t.Fatalf("expected verification event for %s, got %+v", taskID, events)
	}
	if anomalyIndex > verificationIndex {
		t.Fatalf("expected anomaly before verification, anomaly index %d verification index %d", anomalyIndex, verificationIndex)
	}
}

type observedGateContext struct {
	context.Context
	entered chan struct{}
	once    sync.Once
}

func (ctx *observedGateContext) Err() error {
	ctx.once.Do(func() {
		close(ctx.entered)
	})
	return ctx.Context.Err()
}

type gateAcquireResult struct {
	release func()
	err     error
}

func waitGateAcquireResult(t *testing.T, results <-chan gateAcquireResult) gateAcquireResult {
	t.Helper()
	select {
	case result := <-results:
		return result
	case <-time.After(2 * time.Second):
		t.Fatal("timed out waiting for Verification gate acquisition")
		return gateAcquireResult{}
	}
}

func waitObservedGateEntry(t *testing.T, entered <-chan struct{}) {
	t.Helper()
	select {
	case <-entered:
	case <-time.After(2 * time.Second):
		t.Fatal("timed out waiting for Verification gate entry")
	}
}

func TestTaskCycleExclusiveRetryDrainsSharedAttemptsAndBlocksLaterShared(t *testing.T) {
	gate := newVerificationGate(2)
	releaseFirst, err := gate.Acquire(context.Background(), verificationShared)
	if err != nil {
		t.Fatalf("acquire first active shared capacity: %v", err)
	}
	releaseSecond, err := gate.Acquire(context.Background(), verificationShared)
	if err != nil {
		t.Fatalf("acquire second active shared capacity: %v", err)
	}

	exclusiveEntered := make(chan struct{})
	exclusiveResults := make(chan gateAcquireResult, 1)
	go func() {
		release, err := gate.Acquire(&observedGateContext{Context: context.Background(), entered: exclusiveEntered}, verificationExclusive)
		exclusiveResults <- gateAcquireResult{release: release, err: err}
	}()
	waitObservedGateEntry(t, exclusiveEntered)

	sharedEntered := make(chan struct{})
	sharedResults := make(chan gateAcquireResult, 1)
	go func() {
		release, err := gate.Acquire(&observedGateContext{Context: context.Background(), entered: sharedEntered}, verificationShared)
		sharedResults <- gateAcquireResult{release: release, err: err}
	}()
	waitObservedGateEntry(t, sharedEntered)

	select {
	case result := <-exclusiveResults:
		t.Fatalf("exclusive acquired before active shared drained: %+v", result)
	default:
	}
	select {
	case result := <-sharedResults:
		t.Fatalf("later shared bypassed queued exclusive: %+v", result)
	default:
	}

	releaseFirst()
	select {
	case result := <-exclusiveResults:
		t.Fatalf("exclusive acquired before both active shared attempts drained: %+v", result)
	default:
	}
	releaseSecond()
	exclusive := waitGateAcquireResult(t, exclusiveResults)
	if exclusive.err != nil {
		t.Fatalf("acquire exclusive after shared drain: %v", exclusive.err)
	}
	select {
	case result := <-sharedResults:
		t.Fatalf("later shared overlapped active exclusive: %+v", result)
	default:
	}

	exclusive.release()
	shared := waitGateAcquireResult(t, sharedResults)
	if shared.err != nil {
		t.Fatalf("acquire shared after exclusive release: %v", shared.err)
	}
	shared.release()
}

func TestTaskCycleExclusiveRetryCancellationRestoresFullSharedCapacity(t *testing.T) {
	gate := newVerificationGate(2)
	releaseFirst, err := gate.Acquire(context.Background(), verificationShared)
	if err != nil {
		t.Fatalf("acquire first shared capacity: %v", err)
	}
	releaseSecond, err := gate.Acquire(context.Background(), verificationShared)
	if err != nil {
		t.Fatalf("acquire second shared capacity: %v", err)
	}

	exclusiveCtx, cancelExclusive := context.WithCancel(context.Background())
	exclusiveEntered := make(chan struct{})
	exclusiveResults := make(chan gateAcquireResult, 1)
	go func() {
		release, err := gate.Acquire(&observedGateContext{Context: exclusiveCtx, entered: exclusiveEntered}, verificationExclusive)
		exclusiveResults <- gateAcquireResult{release: release, err: err}
	}()
	waitObservedGateEntry(t, exclusiveEntered)
	cancelExclusive()

	canceled := waitGateAcquireResult(t, exclusiveResults)
	if !errors.Is(canceled.err, context.Canceled) {
		t.Fatalf("expected queued exclusive cancellation, got %v", canceled.err)
	}
	if canceled.release != nil {
		t.Fatal("expected canceled acquisition not to return a release function")
	}

	releaseFirst()
	releaseFirst()
	releaseSecond()
	releaseSecond()
	releaseAfterCancelFirst, err := gate.Acquire(context.Background(), verificationShared)
	if err != nil {
		t.Fatalf("reacquire first shared capacity after cancellation: %v", err)
	}
	releaseAfterCancelSecond, err := gate.Acquire(context.Background(), verificationShared)
	if err != nil {
		t.Fatalf("reacquire second shared capacity after cancellation: %v", err)
	}
	releaseAfterCancelFirst()
	releaseAfterCancelSecond()
}

// TestTaskCycleVerificationGateTryAcquireDistinguishesQueuedAttempts pins
// the fact the cancellation contract depends on: whether an attempt owns
// capacity right away, or is queued and must start honoring Stop Requests.
func TestTaskCycleVerificationGateTryAcquireDistinguishesQueuedAttempts(t *testing.T) {
	gate := newVerificationGate(1)
	release, acquired := gate.TryAcquire(verificationShared)
	if !acquired {
		t.Fatal("expected free capacity granted without waiting")
	}
	if _, queued := gate.TryAcquire(verificationShared); queued {
		t.Fatal("expected a full gate to report the next attempt as queued")
	}
	release()
	regained, acquired := gate.TryAcquire(verificationShared)
	if !acquired {
		t.Fatal("expected a released permit to be reusable")
	}

	exclusiveEntered := make(chan struct{})
	exclusiveResults := make(chan gateAcquireResult, 1)
	go func() {
		release, err := gate.Acquire(&observedGateContext{Context: context.Background(), entered: exclusiveEntered}, verificationExclusive)
		exclusiveResults <- gateAcquireResult{release: release, err: err}
	}()
	waitObservedGateEntry(t, exclusiveEntered)
	if _, acquired := gate.TryAcquire(verificationShared); acquired {
		t.Fatal("expected a queued exclusive retry to keep later shared attempts queued")
	}
	if _, acquired := gate.TryAcquire(verificationExclusive); acquired {
		t.Fatal("expected a queued exclusive retry not to be bypassed")
	}

	regained()
	exclusive := waitGateAcquireResult(t, exclusiveResults)
	if exclusive.err != nil {
		t.Fatalf("acquire exclusive after shared drain: %v", exclusive.err)
	}
	exclusive.release()
	final, acquired := gate.TryAcquire(verificationShared)
	if !acquired {
		t.Fatal("expected full shared capacity back after the exclusive retry released")
	}
	final()
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

func TestTaskCyclePublishesCapacities(t *testing.T) {
	fixture := newTaskCycleFixture(t, []taskSpecSeed{{id: "task_01"}})
	engine := fixture.engine(t, &taskFakeRunner{calls: fixture.calls, gitRoot: fixture.gitRoot}, &taskFakeVerifier{calls: fixture.calls}, &engineFakeCommitter{calls: fixture.calls}, fixture.worktree)
	plan := fixture.plan()
	plan.Concurrency = 1
	plan.VerificationConcurrency = 3

	result, err := engine.TaskCycle(context.Background(), plan)
	if err != nil {
		t.Fatalf("task cycle: %v", err)
	}
	if result.Completed != 1 {
		t.Fatalf("expected Task to complete, got %+v", result)
	}

	statusEvents := taskEventsOfKind(fixture.sink, runevent.KindDaemonStatus)
	if len(statusEvents) == 0 {
		t.Fatal("expected Task-cycle-start event")
	}
	payload := eventPayloadMap(t, statusEvents[0])
	if payload["concurrency"] != float64(1) || payload["task_capacity"] != float64(1) || payload["verification_capacity"] != float64(3) {
		t.Fatalf("unexpected Task-cycle capacities: %+v", payload)
	}
}

func TestTaskCycleRejectsInvalidCapacitiesBeforeSideEffects(t *testing.T) {
	tests := []struct {
		name                 string
		taskCapacity         int
		verificationCapacity int
		want                 string
	}{
		{name: "missing Task Capacity", taskCapacity: 0, verificationCapacity: 1, want: "Task Capacity"},
		{name: "negative Task Capacity", taskCapacity: -1, verificationCapacity: 1, want: "Task Capacity"},
		{name: "missing Verification Capacity", taskCapacity: 1, verificationCapacity: 0, want: "Verification Capacity"},
		{name: "negative Verification Capacity", taskCapacity: 1, verificationCapacity: -1, want: "Verification Capacity"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fixture := newTaskCycleFixture(t, []taskSpecSeed{{id: "task_01"}})
			engine := fixture.engine(t, &taskFakeRunner{calls: fixture.calls, gitRoot: fixture.gitRoot}, &taskFakeVerifier{calls: fixture.calls}, &engineFakeCommitter{calls: fixture.calls}, fixture.worktree)
			plan := fixture.plan()
			plan.Concurrency = tt.taskCapacity
			plan.VerificationConcurrency = tt.verificationCapacity
			stateBefore := runStateForTest(fixture.store, fixture.run.ID)

			_, err := engine.TaskCycle(context.Background(), plan)

			if err == nil || !strings.Contains(err.Error(), tt.want) {
				t.Fatalf("expected %s validation error, got %v", tt.want, err)
			}
			if len(*fixture.calls) != 0 {
				t.Fatalf("expected no Agent, Verification, or settlement action, got %v", *fixture.calls)
			}
			if events := fixture.sink.snapshot(); len(events) != 0 {
				t.Fatalf("expected no Run Event before capacity validation, got %+v", events)
			}
			if fixture.progress.Len() != 0 {
				t.Fatalf("expected no progress output before capacity validation, got %q", fixture.progress.String())
			}
			if got := runStateForTest(fixture.store, fixture.run.ID); got != stateBefore {
				t.Fatalf("Run state changed from %q to %q", stateBefore, got)
			}
		})
	}
}

func TestTaskCycleIntegratedVerificationCapacityOneBoundsConcurrentTaskWorktrees(t *testing.T) {
	fixture := newTaskCycleFixture(t, []taskSpecSeed{
		{id: "task_01", verification: []string{"verify-task_01"}},
		{id: "task_02", verification: []string{"verify-task_02"}},
	})
	taskWorktrees := newFakeTaskWorktrees()
	runner := newTaskSchedulerRunner("task_01", "task_02")
	verifier := newTaskCapacityVerifier("task_01", "task_02")
	engine := fixture.engineWithTaskWorktrees(t, runner, verifier, &taskSchedulerCommitter{}, fixture.worktree, taskWorktrees)
	plan := fixture.plan()
	plan.Concurrency = 2
	plan.VerificationConcurrency = 1

	resultCh := make(chan struct {
		result TaskCycleResult
		err    error
	}, 1)
	go func() {
		result, err := engine.TaskCycle(context.Background(), plan)
		resultCh <- struct {
			result TaskCycleResult
			err    error
		}{result: result, err: err}
	}()

	assertTaskSet(t, waitSchedulerStarts(t, runner, 2), "task_01", "task_02")
	if got := runner.maxObservedActive(); got != 2 {
		t.Fatalf("expected two overlapping Agent turns, got max active %d", got)
	}
	runner.releaseTask("task_01")
	runner.releaseTask("task_02")

	first := verifier.waitStart(t)
	verifier.assertNoStart(t)
	verifier.releaseAttempt(first.taskID, first.attempt)
	second := verifier.waitStart(t)
	if second.taskID == first.taskID {
		t.Fatalf("expected the other Task to acquire after release, got %+v then %+v", first, second)
	}
	verifier.releaseAttempt(second.taskID, second.attempt)

	outcome := waitTaskCycleResult(t, resultCh)
	if outcome.err != nil {
		t.Fatalf("task cycle: %v", outcome.err)
	}
	if outcome.result.Completed != 2 {
		t.Fatalf("expected both Tasks completed, got %+v", outcome.result)
	}
	if maxActive, starts := verifier.observed(); maxActive != 1 || starts != 2 {
		t.Fatalf("expected two Verification starts bounded to one active, got max=%d starts=%d", maxActive, starts)
	}
}

func TestTaskCycleVerificationCapacityTwoOverlapsReadyAttemptsWithoutPermitLoss(t *testing.T) {
	fixture := newTaskCycleFixture(t, []taskSpecSeed{
		{id: "task_01", verification: []string{"verify-task_01"}},
		{id: "task_02", verification: []string{"verify-task_02"}},
	})
	taskWorktrees := newFakeTaskWorktrees()
	runner := newTaskSchedulerRunner("task_01", "task_02")
	verifier := newTaskCapacityVerifier("task_01", "task_02")
	engine := fixture.engineWithTaskWorktrees(t, runner, verifier, &taskSchedulerCommitter{}, fixture.worktree, taskWorktrees)
	plan := fixture.plan()
	plan.Concurrency = 2
	plan.VerificationConcurrency = 2

	resultCh := make(chan struct {
		result TaskCycleResult
		err    error
	}, 1)
	go func() {
		result, err := engine.TaskCycle(context.Background(), plan)
		resultCh <- struct {
			result TaskCycleResult
			err    error
		}{result: result, err: err}
	}()

	assertTaskSet(t, waitSchedulerStarts(t, runner, 2), "task_01", "task_02")
	runner.releaseTask("task_01")
	runner.releaseTask("task_02")
	first := verifier.waitStart(t)
	second := verifier.waitStart(t)
	if first.taskID == second.taskID {
		t.Fatalf("expected distinct ready Tasks to overlap, got %+v and %+v", first, second)
	}
	if maxActive, _ := verifier.observed(); maxActive != 2 {
		t.Fatalf("expected two active Verification attempts, got %d", maxActive)
	}
	verifier.releaseAttempt(first.taskID, first.attempt)
	verifier.releaseAttempt(second.taskID, second.attempt)

	outcome := waitTaskCycleResult(t, resultCh)
	if outcome.err != nil {
		t.Fatalf("task cycle: %v", outcome.err)
	}
	if outcome.result.Completed != 2 {
		t.Fatalf("expected both Tasks completed without permit loss, got %+v", outcome.result)
	}
	if maxActive, starts := verifier.observed(); maxActive != 2 || starts != 2 {
		t.Fatalf("expected exactly two overlapping Verification starts, got max=%d starts=%d", maxActive, starts)
	}
}

func TestTaskCycleWaitingForVerificationPrecedesStartedWhenImmediatelyAvailable(t *testing.T) {
	fixture := newTaskCycleFixture(t, []taskSpecSeed{{id: "task_01", verification: []string{"verify-task_01"}}})
	engine := fixture.engine(t, &taskFakeRunner{calls: fixture.calls, gitRoot: fixture.gitRoot}, &taskFakeVerifier{calls: fixture.calls}, &engineFakeCommitter{calls: fixture.calls}, fixture.worktree)

	result, err := engine.TaskCycle(context.Background(), fixture.plan())
	if err != nil {
		t.Fatalf("task cycle: %v", err)
	}
	if result.Completed != 1 {
		t.Fatalf("expected Task completed, got %+v", result)
	}

	events := taskEventsOfKind(fixture.sink, runevent.KindDaemonVerification)
	waitingIndex := -1
	startedIndex := -1
	for index, event := range events {
		payload := eventPayloadMap(t, event)
		switch payload["phase"] {
		case string(runevent.VerificationPhaseWaiting):
			waitingIndex = index
			if payload["mode"] != "shared" || payload["capacity"] != float64(1) {
				t.Fatalf("unexpected waiting payload: %+v", payload)
			}
		case string(runevent.VerificationPhaseStarted):
			startedIndex = index
		}
	}
	if waitingIndex < 0 || startedIndex < 0 || waitingIndex >= startedIndex {
		t.Fatalf("expected waiting before started for immediately available capacity, got %+v", events)
	}
}

func TestTaskCycleRepairReacquiresVerificationCapacityAfterFeedback(t *testing.T) {
	fixture := newTaskCycleFixture(t, []taskSpecSeed{
		{id: "task_01", verification: []string{"verify-task_01"}},
		{id: "task_02", verification: []string{"verify-task_02"}},
	})
	fixture.sink.published = make(chan runevent.RunEvent, 128)
	taskWorktrees := newFakeTaskWorktrees()
	runner := newTaskSchedulerRunner("task_01", "task_02")
	repairRelease := make(chan struct{})
	var runnerCallsMu sync.Mutex
	runnerCalls := map[string]int{}
	runner.onStart = func(taskID string, _ agent.ExecuteRequest) error {
		runnerCallsMu.Lock()
		runnerCalls[taskID]++
		call := runnerCalls[taskID]
		runnerCallsMu.Unlock()
		if taskID == "task_01" && call == 2 {
			<-repairRelease
		}
		return nil
	}
	verifier := newTaskCapacityVerifier("task_01", "task_02")
	verifier.failAttempt("task_01", 1, errors.New("deterministic failure"))
	verifier.addAttempt("task_01", 2)
	engine := fixture.engineWithTaskWorktrees(t, runner, verifier, &taskSchedulerCommitter{}, fixture.worktree, taskWorktrees)
	plan := fixture.plan()
	plan.Concurrency = 2
	plan.VerificationConcurrency = 1

	resultCh := make(chan struct {
		result TaskCycleResult
		err    error
	}, 1)
	go func() {
		result, err := engine.TaskCycle(context.Background(), plan)
		resultCh <- struct {
			result TaskCycleResult
			err    error
		}{result: result, err: err}
	}()

	assertTaskSet(t, waitSchedulerStarts(t, runner, 2), "task_01", "task_02")
	runner.releaseTask("task_01")
	if started := verifier.waitStart(t); started != (taskVerificationStart{taskID: "task_01", attempt: 1}) {
		t.Fatalf("expected task_01 attempt 1 first, got %+v", started)
	}
	if started := waitSchedulerStarts(t, runner, 1); len(started) != 1 || started[0] != "task_01" {
		t.Fatalf("expected task_01 Verification Feedback Agent, got %v", started)
	}

	runner.releaseTask("task_02")
	if started := verifier.waitStart(t); started != (taskVerificationStart{taskID: "task_02", attempt: 1}) {
		t.Fatalf("expected task_02 to verify while task_01 repair held no capacity, got %+v", started)
	}
	close(repairRelease)
	waitPublishedVerificationPhase(t, fixture.sink.published, "task_01", 2, runevent.VerificationPhaseWaiting)
	verifier.assertNoStart(t)

	verifier.releaseAttempt("task_02", 1)
	if started := verifier.waitStart(t); started != (taskVerificationStart{taskID: "task_01", attempt: 2}) {
		t.Fatalf("expected task_01 attempt 2 to reacquire after task_02, got %+v", started)
	}
	verifier.releaseAttempt("task_01", 2)

	outcome := waitTaskCycleResult(t, resultCh)
	if outcome.err != nil {
		t.Fatalf("task cycle: %v", outcome.err)
	}
	if outcome.result.Completed != 2 {
		t.Fatalf("expected both Tasks completed, got %+v", outcome.result)
	}
	if maxActive, starts := verifier.observed(); maxActive != 1 || starts != 3 {
		t.Fatalf("expected released/reacquired capacity across three attempts, got max=%d starts=%d", maxActive, starts)
	}
	assertVerificationWaitingBeforeStarted(t, fixture.sink.snapshot())
	assertVerificationFeedbackJournaled(t, fixture.sink.snapshot(), "task_01", "task_02")
}

// assertVerificationFeedbackJournaled proves the Verification Feedback turn
// is observable per Task: only the repaired Task journals it, and it lands
// between the attempt that failed and the attempt that follows, so a
// consumer can move that Task — and only that Task — back to Agent work.
func assertVerificationFeedbackJournaled(t *testing.T, events []runevent.RunEvent, repaired string, untouched string) {
	t.Helper()
	feedbackIndex := -1
	failedVerdict := -1
	nextWaiting := -1
	for index, event := range events {
		switch event.Kind {
		case runevent.KindDaemonTask:
			if eventPayloadString(t, event, "phase") != "verification_feedback" {
				continue
			}
			if event.ReviewIssue == untouched {
				t.Fatalf("expected no Verification Feedback event for %s, got %+v", untouched, event)
			}
			if event.ReviewIssue != repaired {
				continue
			}
			if feedbackIndex >= 0 {
				t.Fatalf("expected one Verification Feedback event for %s, got a second at %d", repaired, index)
			}
			feedbackIndex = index
		case runevent.KindDaemonVerification:
			payload := eventPayloadMap(t, event)
			taskID, _ := payload["task"].(string)
			attempt, _ := payload["attempt"].(float64)
			if taskID != repaired {
				continue
			}
			if payload["phase"] == string(runevent.VerificationPhaseVerdict) && int(attempt) == 1 {
				failedVerdict = index
			}
			if payload["phase"] == string(runevent.VerificationPhaseWaiting) && int(attempt) == 2 {
				nextWaiting = index
			}
		}
	}
	if feedbackIndex < 0 {
		t.Fatalf("expected a verification_feedback event for %s, got none", repaired)
	}
	if failedVerdict < 0 || feedbackIndex < failedVerdict {
		t.Fatalf("expected the Verification Feedback event after the failed verdict, got verdict=%d feedback=%d", failedVerdict, feedbackIndex)
	}
	if nextWaiting < 0 || feedbackIndex > nextWaiting {
		t.Fatalf("expected the Verification Feedback event before the next waiting event, got feedback=%d waiting=%d", feedbackIndex, nextWaiting)
	}
}

func TestTaskCycleVerificationCapacityCancellationWhileQueuedStartsNoCommandOrSettlement(t *testing.T) {
	fixture := newTaskCycleFixture(t, []taskSpecSeed{
		{id: "task_01", verification: []string{"verify-task_01"}},
		{id: "task_02", verification: []string{"verify-task_02"}},
	})
	fixture.sink.published = make(chan runevent.RunEvent, 128)
	taskWorktrees := newFakeTaskWorktrees()
	runner := newTaskSchedulerRunner("task_01", "task_02")
	verifier := newTaskCapacityVerifier("task_01", "task_02")
	engine := fixture.engineWithTaskWorktrees(t, runner, verifier, &taskSchedulerCommitter{}, fixture.worktree, taskWorktrees)
	plan := fixture.plan()
	plan.Concurrency = 2
	plan.VerificationConcurrency = 1
	ctx, cancel := context.WithCancel(context.Background())

	resultCh := make(chan struct {
		result TaskCycleResult
		err    error
	}, 1)
	go func() {
		result, err := engine.TaskCycle(ctx, plan)
		resultCh <- struct {
			result TaskCycleResult
			err    error
		}{result: result, err: err}
	}()

	assertTaskSet(t, waitSchedulerStarts(t, runner, 2), "task_01", "task_02")
	runner.releaseTask("task_01")
	runner.releaseTask("task_02")
	active := verifier.waitStart(t)
	queuedTask := "task_01"
	if active.taskID == queuedTask {
		queuedTask = "task_02"
	}
	waitPublishedVerificationPhase(t, fixture.sink.published, queuedTask, 1, runevent.VerificationPhaseWaiting)
	verifier.assertNoStart(t)
	cancel()

	outcome := waitTaskCycleResult(t, resultCh)
	if outcome.err == nil || !errors.Is(outcome.err, context.Canceled) {
		t.Fatalf("expected queued cancellation, got result=%+v err=%v", outcome.result, outcome.err)
	}
	if _, starts := verifier.observed(); starts != 1 {
		t.Fatalf("expected zero Verification commands for queued Task, total starts=%d", starts)
	}
	for _, status := range []spec.Status{spec.StatusCompleted, spec.StatusFailed} {
		if hasTaskSettlementEvent(t, fixture.sink.snapshot(), queuedTask, status) {
			t.Fatalf("expected no false terminal settlement for queued Task %s", queuedTask)
		}
	}
	taskRef := taskWorktrees.taskRef(queuedTask)
	if got := taskStatusInSpecRootOnDisk(t, filepath.Join(taskRef.Path, "docs", "specs"), queuedTask); got != string(spec.StatusInProgress) {
		t.Fatalf("expected queued Task left resumable in_progress, got %q", got)
	}
}

// TestTaskCycleStopRequestWhileQueuedForVerificationStartsNoCommandAndStaysResumable
// covers the public Stop path the direct-cancellation case cannot reach: a
// Stop Request is durable Run-store state, so the queued attempt has to
// read it rather than wait for a context that is never cancelled.
func TestTaskCycleStopRequestWhileQueuedForVerificationStartsNoCommandAndStaysResumable(t *testing.T) {
	fixture := newTaskCycleFixture(t, []taskSpecSeed{
		{id: "task_01", verification: []string{"verify-task_01"}},
		{id: "task_02", verification: []string{"verify-task_02"}},
	})
	fixture.sink.published = make(chan runevent.RunEvent, 128)
	taskWorktrees := newFakeTaskWorktrees()
	runner := newTaskSchedulerRunner("task_01", "task_02")
	verifier := newTaskCapacityVerifier("task_01", "task_02")
	engine := fixture.engineWithTaskWorktrees(t, runner, verifier, &taskSchedulerCommitter{}, fixture.worktree, taskWorktrees)
	plan := fixture.plan()
	plan.Concurrency = 2
	plan.VerificationConcurrency = 1

	resultCh := make(chan struct {
		result TaskCycleResult
		err    error
	}, 1)
	go func() {
		result, err := engine.TaskCycle(context.Background(), plan)
		resultCh <- struct {
			result TaskCycleResult
			err    error
		}{result: result, err: err}
	}()

	assertTaskSet(t, waitSchedulerStarts(t, runner, 2), "task_01", "task_02")
	runner.releaseTask("task_01")
	runner.releaseTask("task_02")
	active := verifier.waitStart(t)
	queuedTask := "task_01"
	if active.taskID == queuedTask {
		queuedTask = "task_02"
	}
	waitPublishedVerificationPhase(t, fixture.sink.published, queuedTask, 1, runevent.VerificationPhaseWaiting)
	verifier.assertNoStart(t)

	if err := fixture.store.RequestStop(context.Background(), fixture.run.ID); err != nil {
		t.Fatalf("record public Stop Request: %v", err)
	}
	// The queued Task leaves the queue on the Stop Request alone, while the
	// active attempt still holds the only Verification permit.
	waitPublishedStopEvent(t, fixture.sink.published)
	verifier.assertNoStart(t)
	verifier.releaseAttempt(active.taskID, 1)

	outcome := waitTaskCycleResult(t, resultCh)
	if !errors.Is(outcome.err, ErrStopRequested) {
		t.Fatalf("expected the queued Task to end the Run stopped, got result=%+v err=%v", outcome.result, outcome.err)
	}
	if _, starts := verifier.observed(); starts != 1 {
		t.Fatalf("expected zero Verification commands for the queued Task, total starts=%d", starts)
	}
	for _, status := range []spec.Status{spec.StatusCompleted, spec.StatusFailed} {
		if hasTaskSettlementEvent(t, fixture.sink.snapshot(), queuedTask, status) {
			t.Fatalf("expected no false terminal settlement for queued Task %s", queuedTask)
		}
	}
	queuedRoot := filepath.Join(taskWorktrees.taskRef(queuedTask).Path, "docs", "specs")
	if got := taskStatusInSpecRootOnDisk(t, queuedRoot, queuedTask); got != string(spec.StatusInProgress) {
		t.Fatalf("expected queued Task left resumable in_progress, got %q", got)
	}
	// The already-running attempt still settles and integrates before the
	// Run stops, so its status reaches the Run Worktree.
	if got := taskStatusOnDisk(t, fixture.gitRoot, active.taskID); got != string(spec.StatusCompleted) {
		t.Fatalf("expected the in-flight Task settled completed, got %q", got)
	}
	if !taskWorktrees.integratedTask(active.taskID) {
		t.Fatalf("expected the in-flight Task integrated before the Run stopped, got %v", taskWorktrees.integratedTasks())
	}
	if outcome.result.Completed != 1 || outcome.result.Failed != 0 || outcome.result.Skipped != 0 {
		t.Fatalf("expected only the in-flight Task counted, got %+v", outcome.result)
	}
}

func waitPublishedStopEvent(t *testing.T, published <-chan runevent.RunEvent) {
	t.Helper()
	for {
		select {
		case event := <-published:
			if event.Kind == runevent.KindDaemonStatus && strings.Contains(event.Summary, "Stop Request") {
				return
			}
		case <-time.After(2 * time.Second):
			t.Fatal("timed out waiting for the Stop Request to reach the Run Event Stream")
		}
	}
}

func waitPublishedVerificationPhase(t *testing.T, published <-chan runevent.RunEvent, taskID string, attempt int, phase runevent.VerificationPhase) {
	t.Helper()
	for {
		select {
		case event := <-published:
			if event.Kind != runevent.KindDaemonVerification || event.ReviewIssue != taskID {
				continue
			}
			payload := eventPayloadMap(t, event)
			if payload["attempt"] == float64(attempt) && payload["phase"] == string(phase) {
				return
			}
		case <-time.After(2 * time.Second):
			t.Fatalf("timed out waiting for Task %s Verification attempt %d phase %s", taskID, attempt, phase)
		}
	}
}

func assertVerificationWaitingBeforeStarted(t *testing.T, events []runevent.RunEvent) {
	t.Helper()
	waiting := map[string]int{}
	for index, event := range events {
		if event.Kind != runevent.KindDaemonVerification {
			continue
		}
		payload := eventPayloadMap(t, event)
		taskID, _ := payload["task"].(string)
		attempt, _ := payload["attempt"].(float64)
		key := taskVerificationKey(taskID, int(attempt))
		switch payload["phase"] {
		case string(runevent.VerificationPhaseWaiting):
			waiting[key] = index
		case string(runevent.VerificationPhaseStarted):
			waitingIndex, ok := waiting[key]
			if !ok || waitingIndex >= index {
				t.Fatalf("Verification %s started without an earlier waiting event", key)
			}
		}
	}
}

func TestTaskCycleSchedulesIndependentWaveWithConcurrencyCap(t *testing.T) {
	fixture := newTaskCycleFixture(t, []taskSpecSeed{
		{id: "task_01", title: "Do first independent work"},
		{id: "task_02", title: "Do second independent work"},
		{id: "task_03", title: "Do third independent work"},
		{id: "task_04", title: "Do fourth independent work"},
	})
	taskWorktrees := newFakeTaskWorktrees()
	runner := newTaskSchedulerRunner("task_01", "task_02", "task_03", "task_04")
	verifier := &taskSchedulerVerifier{}
	committer := &taskSchedulerCommitter{}
	engine := fixture.engineWithTaskWorktrees(t, runner, verifier, committer, fixture.worktree, taskWorktrees)
	plan := fixture.plan()
	plan.Concurrency = 2

	resultCh := make(chan struct {
		result TaskCycleResult
		err    error
	}, 1)
	go func() {
		result, err := engine.TaskCycle(context.Background(), plan)
		resultCh <- struct {
			result TaskCycleResult
			err    error
		}{result: result, err: err}
	}()

	assertTaskSet(t, waitSchedulerStarts(t, runner, 2), "task_01", "task_02")
	assertNoSchedulerStart(t, runner)
	runner.releaseTask("task_01")
	if got := strings.Join(waitSchedulerStarts(t, runner, 1), "|"); got != "task_03" {
		t.Fatalf("expected task_03 to start when one slot opened, got %s", got)
	}
	runner.releaseTask("task_02")
	if got := strings.Join(waitSchedulerStarts(t, runner, 1), "|"); got != "task_04" {
		t.Fatalf("expected task_04 to start when the next slot opened, got %s", got)
	}
	runner.releaseTask("task_03")
	runner.releaseTask("task_04")

	outcome := waitTaskCycleResult(t, resultCh)
	if outcome.err != nil {
		t.Fatalf("task cycle: %v", outcome.err)
	}
	if outcome.result.Completed != 4 || outcome.result.Failed != 0 || outcome.result.Skipped != 0 {
		t.Fatalf("expected all four independent Tasks completed, got %+v", outcome.result)
	}
	if runner.maxObservedActive() != 2 {
		t.Fatalf("expected observed execution overlap capped at 2, got max active %d", runner.maxObservedActive())
	}
	if got := len(taskWorktrees.integratedTasks()); got != 4 {
		t.Fatalf("expected four Task integrations onto the Run Branch, got %d", got)
	}
	if got := committer.commitCount(); got != 4 {
		t.Fatalf("expected four Task commits before integration, got %d", got)
	}
	startEvents := taskStartedEvents(t, fixture.sink)
	sort.Strings(startEvents)
	if got := strings.Join(startEvents, "|"); got != "task_01:1|task_02:2|task_03:3|task_04:4" {
		t.Fatalf("expected deterministic start ordinals in journal replay, got %s", got)
	}
}

func TestTaskCycleCreatesTaskWorktreesWithBootstrapBeforeAgentWork(t *testing.T) {
	fixture := newTaskCycleFixture(t, []taskSpecSeed{
		{id: "task_01", title: "Prepare first Task"},
		{id: "task_02", title: "Prepare second Task"},
	})
	const command = "make task-bootstrap"
	var bootstrapOutput bytes.Buffer
	bootstrapWriter := &lockedWriter{writer: &bootstrapOutput}
	taskWorktrees := newFakeTaskWorktrees()
	taskWorktrees.onCreate = func(taskID string, ref runworktree.TaskRef, opts runworktree.TaskCreateOptions) error {
		if opts.Bootstrap.Command != command {
			return fmt.Errorf("bootstrap command for %s = %q", taskID, opts.Bootstrap.Command)
		}
		if opts.Bootstrap.Timeout != time.Second {
			return fmt.Errorf("bootstrap timeout for %s = %s", taskID, opts.Bootstrap.Timeout)
		}
		if opts.BootstrapOutput == nil {
			return fmt.Errorf("missing bootstrap output writer for %s", taskID)
		}
		if _, err := opts.BootstrapOutput.Write([]byte("bootstrap " + taskID + "\n")); err != nil {
			return err
		}
		return os.WriteFile(filepath.Join(ref.Path, "bootstrap.ready"), []byte(taskID), 0o644)
	}
	runner := &taskSchedulerRunner{
		started: make(chan string, 2),
		statusByTask: map[string]spec.Status{
			"task_01": spec.StatusCompleted,
			"task_02": spec.StatusCompleted,
		},
		onStart: func(taskID string, req agent.ExecuteRequest) error {
			got, err := os.ReadFile(filepath.Join(req.GitRoot, "bootstrap.ready"))
			if err != nil {
				return fmt.Errorf("read bootstrap marker for %s: %w", taskID, err)
			}
			if string(got) != taskID {
				return fmt.Errorf("bootstrap marker for %s = %q", taskID, got)
			}
			return nil
		},
	}
	verifier := &taskSchedulerVerifier{}
	committer := &taskSchedulerCommitter{}
	engine := fixture.engineWithTaskWorktrees(t, runner, verifier, committer, fixture.worktree, taskWorktrees)
	plan := fixture.plan()
	plan.Concurrency = 2
	plan.Bootstrap = runworktree.BootstrapSpec{Command: command, Timeout: time.Second}
	plan.BootstrapOutput = bootstrapWriter

	result, err := engine.TaskCycle(context.Background(), plan)

	if err != nil {
		t.Fatalf("task cycle: %v", err)
	}
	if result.Completed != 2 || result.Failed != 0 || result.Skipped != 0 {
		t.Fatalf("expected both bootstrapped Tasks completed, got %+v", result)
	}
	assertTaskSet(t, runner.startedTasks(), "task_01", "task_02")
	for _, taskID := range []string{"task_01", "task_02"} {
		opts := taskWorktrees.taskCreateOptions(taskID)
		if opts.Bootstrap.Command != command || opts.Bootstrap.Timeout != time.Second {
			t.Fatalf("expected bootstrap options for %s, got %#v", taskID, opts.Bootstrap)
		}
		if !strings.Contains(bootstrapOutput.String(), "bootstrap "+taskID) {
			t.Fatalf("expected bootstrap output for %s, got %q", taskID, bootstrapOutput.String())
		}
	}
}

func TestTaskCycleTaskWorktreeBootstrapFailureIsolatesIndependentTasks(t *testing.T) {
	fixture := newTaskCycleFixture(t, []taskSpecSeed{
		{id: "task_01", title: "Fail bootstrap"},
		{id: "task_02", title: "Keep running"},
		{id: "task_03", title: "Also keep running"},
	})
	const command = "make task-bootstrap"
	taskWorktrees := newFakeTaskWorktrees()
	taskWorktrees.createErrByTask = map[string]error{
		"task_01": &runworktree.BootstrapError{Command: command, Err: errors.New("exit status 7")},
	}
	runner := &taskSchedulerRunner{
		started: make(chan string, 3),
		statusByTask: map[string]spec.Status{
			"task_02": spec.StatusCompleted,
			"task_03": spec.StatusCompleted,
		},
	}
	verifier := &taskSchedulerVerifier{}
	committer := &taskSchedulerCommitter{}
	engine := fixture.engineWithTaskWorktrees(t, runner, verifier, committer, fixture.worktree, taskWorktrees)
	plan := fixture.plan()
	plan.Concurrency = 2
	plan.Bootstrap = runworktree.BootstrapSpec{Command: command, Timeout: time.Second}

	result, err := engine.TaskCycle(context.Background(), plan)

	if err != nil {
		t.Fatalf("expected bootstrap failure to fail only its Task, got %v", err)
	}
	if result.Completed != 2 || result.Failed != 1 || result.Skipped != 0 {
		t.Fatalf("expected one bootstrap-failed Task and two completed independents, got %+v", result)
	}
	assertTaskSet(t, runner.startedTasks(), "task_02", "task_03")
	if got := taskStatusOnDisk(t, taskWorktrees.taskRef("task_01").Path, "task_01"); got != string(spec.StatusFailed) {
		t.Fatalf("expected bootstrap-failed Task settled failed in its Task Worktree, got %q", got)
	}
	for _, integrated := range taskWorktrees.integratedTasks() {
		if integrated == "task_01" {
			t.Fatalf("expected bootstrap-failed Task not integrated, got %v", taskWorktrees.integratedTasks())
		}
	}
	if len(taskWorktrees.integratedTasks()) != 2 {
		t.Fatalf("expected only independent completed Tasks integrated, got %v", taskWorktrees.integratedTasks())
	}
	wantReason := "worktree bootstrap failed: " + command + ": exit status 7"
	var journaled bool
	for _, event := range taskEventsOfKind(fixture.sink, runevent.KindDaemonTask) {
		if event.ReviewIssue == "task_01" && strings.Contains(string(event.Payload), wantReason) {
			journaled = true
		}
	}
	if !journaled {
		t.Fatalf("expected bootstrap failure reason %q journaled", wantReason)
	}
}

func TestTaskCycleGatesDependenciesAndSkipsFailedDependencyChainsUnderConcurrency(t *testing.T) {
	fixture := newTaskCycleFixture(t, []taskSpecSeed{
		{id: "task_01", title: "Build prerequisite"},
		{id: "task_03", title: "Fail independent prerequisite", verification: []string{"fail task_03"}},
		{id: "task_02", title: "Use prerequisite", needs: []string{"task_01"}},
		{id: "task_04", title: "Blocked by failed prerequisite", needs: []string{"task_03"}},
		{id: "task_05", title: "Blocked by skipped prerequisite", needs: []string{"task_04"}},
	})
	taskWorktrees := newFakeTaskWorktrees()
	runner := &taskSchedulerRunner{
		started: make(chan string, 5),
		onStart: func(taskID string, _ agent.ExecuteRequest) error {
			if taskID == "task_02" && !taskWorktrees.integratedTask("task_01") {
				return errors.New("task_02 started before task_01 settled")
			}
			return nil
		},
	}
	verifier := &taskSchedulerVerifier{failOn: map[string]error{"fail task_03": errors.New("gate broke")}}
	committer := &taskSchedulerCommitter{}
	engine := fixture.engineWithTaskWorktrees(t, runner, verifier, committer, fixture.worktree, taskWorktrees)
	plan := fixture.plan()
	plan.Concurrency = 2

	result, err := engine.TaskCycle(context.Background(), plan)

	if err != nil {
		t.Fatalf("task cycle: %v", err)
	}
	if result.Completed != 2 || result.Failed != 1 || result.Skipped != 2 {
		t.Fatalf("expected dependency-gated counts 2 completed, 1 failed, 2 skipped, got %+v", result)
	}
	assertTaskSet(t, runner.startedTasks(), "task_01", "task_02", "task_03", "task_03")
	if got := taskStatusOnDisk(t, fixture.gitRoot, "task_04"); got != string(spec.StatusPending) {
		t.Fatalf("expected skipped task_04 left pending on disk, got %q", got)
	}
	if got := taskStatusOnDisk(t, fixture.gitRoot, "task_05"); got != string(spec.StatusPending) {
		t.Fatalf("expected skipped task_05 left pending on disk, got %q", got)
	}
	if got := len(taskWorktrees.integratedTasks()); got != 2 {
		t.Fatalf("expected only completed Tasks integrated, got %d", got)
	}
}

func TestTaskCycleIntegrationConflictSettlesTaskFailedAndKeepsTaskWorktree(t *testing.T) {
	fixture := newTaskCycleFixture(t, []taskSpecSeed{
		{id: "task_01", title: "Write first independent change"},
		{id: "task_02", title: "Collide during integration"},
		{id: "task_03", title: "Write third independent change"},
	})
	taskWorktrees := newFakeTaskWorktrees()
	taskWorktrees.conflictByTask = map[string]string{"task_02": "shared.txt"}
	runner := &taskSchedulerRunner{started: make(chan string, 3)}
	verifier := &taskSchedulerVerifier{}
	committer := &taskSchedulerCommitter{}
	engine := fixture.engineWithTaskWorktrees(t, runner, verifier, committer, fixture.worktree, taskWorktrees)
	plan := fixture.plan()
	plan.Concurrency = 2

	result, err := engine.TaskCycle(context.Background(), plan)

	if err != nil {
		t.Fatalf("task cycle: %v", err)
	}
	if result.Completed != 2 || result.Failed != 1 || result.Skipped != 0 {
		t.Fatalf("expected one integration-conflicted Task failed and independents completed, got %+v", result)
	}
	cleaned := strings.Join(taskWorktrees.cleanedTasks(), "|")
	if strings.Contains(cleaned, "task_02") {
		t.Fatalf("expected conflicting Task Worktree kept, cleaned %s", cleaned)
	}
	conflicting := taskWorktrees.taskRef("task_02")
	if _, err := os.Stat(conflicting.Path); err != nil {
		t.Fatalf("expected conflicting Task Worktree kept at %q: %v", conflicting.Path, err)
	}
	if got := taskStatusOnDisk(t, conflicting.Path, "task_02"); got != string(spec.StatusFailed) {
		t.Fatalf("expected conflicting Task settled failed in its Task Worktree, got %q", got)
	}
	var journaled bool
	for _, event := range taskEventsOfKind(fixture.sink, runevent.KindDaemonTask) {
		if event.ReviewIssue == "task_02" && strings.Contains(string(event.Payload), "integration conflict: shared.txt") {
			journaled = true
		}
	}
	if !journaled {
		t.Fatal("expected integration conflict reason journaled for task_02")
	}
	if got := len(taskWorktrees.integratedTasks()); got != 3 {
		t.Fatalf("expected all completed Task branches offered to the integration queue, got %d", got)
	}
}

func TestTaskCycleStopRequestMidWaveDrainsRunningTasksAndStartsNothingNew(t *testing.T) {
	fixture := newTaskCycleFixture(t, []taskSpecSeed{
		{id: "task_01", title: "Start first"},
		{id: "task_02", title: "Start second"},
		{id: "task_03", title: "Must not start"},
		{id: "task_04", title: "Must also not start"},
	})
	taskWorktrees := newFakeTaskWorktrees()
	taskWorktrees.onIntegrate = func(taskID string) {
		if taskID == "task_01" {
			_ = fixture.store.RequestStop(context.Background(), fixture.run.ID)
		}
	}
	runner := newTaskSchedulerRunner("task_01", "task_02", "task_03", "task_04")
	verifier := &taskSchedulerVerifier{}
	committer := &taskSchedulerCommitter{}
	engine := fixture.engineWithTaskWorktrees(t, runner, verifier, committer, fixture.worktree, taskWorktrees)
	plan := fixture.plan()
	plan.Concurrency = 2

	resultCh := make(chan struct {
		result TaskCycleResult
		err    error
	}, 1)
	go func() {
		result, err := engine.TaskCycle(context.Background(), plan)
		resultCh <- struct {
			result TaskCycleResult
			err    error
		}{result: result, err: err}
	}()

	assertTaskSet(t, waitSchedulerStarts(t, runner, 2), "task_01", "task_02")
	runner.releaseTask("task_01")
	if got := waitIntegratedTask(t, taskWorktrees); got != "task_01" {
		t.Fatalf("expected task_01 to integrate first, got %s", got)
	}
	assertNoSchedulerStart(t, runner)
	runner.releaseTask("task_02")

	outcome := waitTaskCycleResult(t, resultCh)
	if !errors.Is(outcome.err, ErrStopRequested) {
		t.Fatalf("expected ErrStopRequested after draining running Tasks, got %v", outcome.err)
	}
	if outcome.result.Completed != 2 || outcome.result.Failed != 0 || outcome.result.Skipped != 0 {
		t.Fatalf("expected the two running Tasks settled before stop, got %+v", outcome.result)
	}
	assertTaskSet(t, runner.startedTasks(), "task_01", "task_02")
	if got := len(taskWorktrees.integratedTasks()); got != 2 {
		t.Fatalf("expected both running Tasks integrated before stop, got %d", got)
	}
	if got := taskStatusOnDisk(t, fixture.gitRoot, "task_03"); got != string(spec.StatusPending) {
		t.Fatalf("expected task_03 left pending, got %q", got)
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
		calls:     fixture.calls,
		gitRoot:   fixture.gitRoot,
		store:     fixture.store,
		writeLogs: true,
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
	expectedFirst := "docs: write the usage docs\n\nRoundfix-Spec: " + taskCycleSlug + "\nRoundfix-Task: task_01"
	if committer.messages[0] != expectedFirst {
		t.Fatalf("expected docs commit message with trailers, got %q", committer.messages[0])
	}
	expectedSecond := "feat: add the backend behavior\n\nRoundfix-Spec: " + taskCycleSlug + "\nRoundfix-Task: task_02"
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
	expectedLog := agent.LogPath(fixture.artifactDir, fixture.run.ID, 1)
	if runner.requests[0].LogPath != expectedLog {
		t.Fatalf("expected Artifact Directory log path %q, got %q", expectedLog, runner.requests[0].LogPath)
	}
	if _, err := os.Stat(expectedLog); err != nil {
		t.Fatalf("expected fake Agent log under Artifact Directory: %v", err)
	}
	if _, err := os.Stat(filepath.Join(fixture.gitRoot, ".roundfix")); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("expected no repo .roundfix directory, got err=%v", err)
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
	statusEvents := taskEventsOfKind(fixture.sink, runevent.KindDaemonStatus)
	if len(statusEvents) == 0 || !strings.Contains(string(statusEvents[0].Payload), `"concurrency":1`) {
		t.Fatalf("expected TaskCycle start event to journal concurrency, got %+v", statusEvents)
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

func TestVerifyTaskRejectsMissingRetryStateWithoutMutatingRun(t *testing.T) {
	fixture := newTaskCycleFixture(t, []taskSpecSeed{{id: "task_01"}})
	engine := fixture.engine(
		t,
		&taskFakeRunner{calls: fixture.calls},
		&taskFakeVerifier{calls: fixture.calls},
		&engineFakeCommitter{calls: fixture.calls},
		fixture.worktree,
	)
	before := runStateForTest(fixture.store, fixture.run.ID)

	_, err := engine.verifyTask(context.Background(), fixture.plan(), fixture.graph.Tasks[0], 1, 1, nil)

	if err == nil || !strings.Contains(err.Error(), "temporary retry state is required") {
		t.Fatalf("verifyTask() error = %v, want missing temporary retry state", err)
	}
	if after := runStateForTest(fixture.store, fixture.run.ID); after != before {
		t.Fatalf("verifyTask() changed Run state from %q to %q for invalid input", before, after)
	}
}

func TestTaskCycleRewritesNormalizedStatusAfterAgentReload(t *testing.T) {
	fixture := newTaskCycleFixture(t, []taskSpecSeed{{id: "task_01", title: "Normalize the task status"}})
	runner := &taskFakeRunner{
		calls:           fixture.calls,
		gitRoot:         fixture.gitRoot,
		rawStatusByTask: map[string]string{"task_01": "done"},
	}
	verifier := &taskFakeVerifier{calls: fixture.calls}
	committer := &engineFakeCommitter{calls: fixture.calls}
	engine := fixture.engine(t, runner, verifier, committer, fixture.worktree)

	result, err := engine.TaskCycle(context.Background(), fixture.plan())

	if err != nil {
		t.Fatalf("task cycle: %v", err)
	}
	if result.Completed != 1 || result.Failed != 0 || result.Skipped != 0 {
		t.Fatalf("expected synonym-authored Task completed, got %+v", result)
	}
	if got := taskStatusOnDisk(t, fixture.gitRoot, "task_01"); got != string(spec.StatusCompleted) {
		t.Fatalf("expected status synonym rewritten canonical on disk, got %q", got)
	}
	if got := strings.Join(*fixture.calls, ">"); got != "agent>verify>commit" {
		t.Fatalf("expected normal Agent flow after synonym normalization, got %q", got)
	}
	if len(committer.paths) != 1 || strings.Join(committer.paths[0], "|") != taskFileRel(taskCycleSlug, "task_01") {
		t.Fatalf("expected only the canonicalized task file committed, got %v", committer.paths)
	}
}

func TestTaskCycleDaemonStatusInProgressBeforeAgentAndStartEvent(t *testing.T) {
	fixture := newTaskCycleFixture(t, []taskSpecSeed{{id: "task_01"}})
	runner := &taskFakeRunner{calls: fixture.calls, gitRoot: fixture.gitRoot}
	engine := fixture.engine(t, runner, &taskFakeVerifier{calls: fixture.calls}, &engineFakeCommitter{calls: fixture.calls}, fixture.worktree)

	result, err := engine.TaskCycle(context.Background(), fixture.plan())

	if err != nil {
		t.Fatalf("task cycle: %v", err)
	}
	if result.Completed != 1 {
		t.Fatalf("expected Task to complete, got %+v", result)
	}
	if len(runner.prompts) != 1 || !strings.Contains(runner.prompts[0], "\nstatus: in_progress\n") {
		t.Fatalf("expected first Agent prompt to embed the Daemon-owned in_progress Task, got %q", runner.prompts)
	}
	taskEvents := taskEventsOfKind(fixture.sink, runevent.KindDaemonTask)
	if len(taskEvents) < 1 {
		t.Fatal("expected Task start event")
	}
	payload := eventPayloadMap(t, taskEvents[0])
	if payload["phase"] != "started" || payload["status"] != string(spec.StatusInProgress) {
		t.Fatalf("expected Task start event to agree with in_progress on disk, got %+v", payload)
	}
}

func TestTaskCycleAgentStatusVariantsReachDaemonVerification(t *testing.T) {
	statuses := []spec.Status{
		spec.StatusPending,
		spec.StatusInProgress,
		spec.StatusCompleted,
		spec.StatusFailed,
	}
	for _, status := range statuses {
		t.Run(string(status), func(t *testing.T) {
			fixture := newTaskCycleFixture(t, []taskSpecSeed{{id: "task_01"}})
			runner := &taskFakeRunner{
				calls:        fixture.calls,
				gitRoot:      fixture.gitRoot,
				statusByTask: map[string]spec.Status{"task_01": status},
			}
			verifier := &taskFakeVerifier{calls: fixture.calls}
			engine := fixture.engine(t, runner, verifier, &engineFakeCommitter{calls: fixture.calls}, fixture.worktree)

			result, err := engine.TaskCycle(context.Background(), fixture.plan())

			if err != nil {
				t.Fatalf("task cycle: %v", err)
			}
			if result.Completed != 1 || result.Failed != 0 {
				t.Fatalf("expected status %q to be normalized before passing settlement, got %+v", status, result)
			}
			if got := strings.Join(*fixture.calls, ">"); got != "agent>verify>commit" {
				t.Fatalf("expected status %q to reach Daemon Verification, got %q", status, got)
			}
			if got := taskStatusOnDisk(t, fixture.gitRoot, "task_01"); got != string(spec.StatusCompleted) {
				t.Fatalf("expected Daemon-settled completed status after %q handoff, got %q", status, got)
			}
		})
	}
}

func TestTaskCycleAgentStatusNormalizationPreservesResultThroughVerification(t *testing.T) {
	fixture := newTaskCycleFixture(t, []taskSpecSeed{{id: "task_01"}})
	const resultEvidence = "\n## Result\n\n- Focused check passed and implementation is ready.\n"
	runner := &taskFakeRunner{
		calls:        fixture.calls,
		gitRoot:      fixture.gitRoot,
		statusByTask: map[string]spec.Status{"task_01": spec.StatusFailed},
		resultByTask: map[string]string{"task_01": resultEvidence},
	}
	committer := &engineFakeCommitter{calls: fixture.calls}
	engine := fixture.engine(t, runner, &taskFakeVerifier{calls: fixture.calls}, committer, fixture.worktree)

	result, err := engine.TaskCycle(context.Background(), fixture.plan())

	if err != nil {
		t.Fatalf("task cycle: %v", err)
	}
	if result.Completed != 1 {
		t.Fatalf("expected passing Verification to settle completed, got %+v", result)
	}
	content, err := os.ReadFile(taskPathFor(fixture.gitRoot, taskCycleSlug, "task_01"))
	if err != nil {
		t.Fatalf("read Task after settlement: %v", err)
	}
	if !strings.Contains(string(content), resultEvidence) {
		t.Fatalf("expected Agent-authored Result prose to survive status normalization, got:\n%s", content)
	}
	if len(committer.paths) != 1 || !slices.Contains(committer.paths[0], taskFileRel(taskCycleSlug, "task_01")) {
		t.Fatalf("expected passing Task and Result evidence committed together, got %v", committer.paths)
	}
}

func TestTaskCycleAgentStatusRepairHandoffReachesFinalVerification(t *testing.T) {
	fixture := newTaskCycleFixture(t, []taskSpecSeed{{id: "task_01"}})
	runner := &taskFakeRunner{
		calls:            fixture.calls,
		gitRoot:          fixture.gitRoot,
		statusByTaskCall: map[string][]spec.Status{"task_01": {spec.StatusCompleted, spec.StatusFailed}},
	}
	verifier := &taskFakeVerifier{calls: fixture.calls, script: []error{errors.New("gate broke"), nil}}
	engine := fixture.engine(t, runner, verifier, &engineFakeCommitter{calls: fixture.calls}, fixture.worktree)

	result, err := engine.TaskCycle(context.Background(), fixture.plan())

	if err != nil {
		t.Fatalf("task cycle: %v", err)
	}
	if result.Completed != 1 || result.Failed != 0 {
		t.Fatalf("expected repair handoff status to be normalized before final Verification, got %+v", result)
	}
	if got := strings.Join(*fixture.calls, ">"); got != "agent>verify>agent>verify>commit" {
		t.Fatalf("expected one same-Session repair and final Verification, got %q", got)
	}
	if len(runner.requests) != 2 || runner.requests[0].Session != runner.requests[1].Session {
		t.Fatalf("expected Verification Feedback to reuse the Agent Session, got %+v", runner.requests)
	}
}

func TestTaskCycleCommitsAfterTransportAnomalyAndPassingVerification(t *testing.T) {
	fixture := newTaskCycleFixture(t, []taskSpecSeed{{id: "task_01", title: "Classify by parsed result"}})
	fixture.worktree.snapshots = [][]string{nil, {"internal/agent/fix.go"}}
	const anomaly = "acpx exited with exit code 1 after parsed session/prompt result\n--- acpx stderr tail ---\nMessage buffer exceeded 10485760 bytes"
	runner := &taskFakeRunner{
		calls:         fixture.calls,
		gitRoot:       fixture.gitRoot,
		anomalyByTask: map[string]string{"task_01": anomaly},
	}
	verifier := &taskFakeVerifier{calls: fixture.calls}
	committer := &engineFakeCommitter{calls: fixture.calls}
	engine := fixture.engine(t, runner, verifier, committer, fixture.worktree)

	result, err := engine.TaskCycle(context.Background(), fixture.plan())

	if err != nil {
		t.Fatalf("task cycle: %v", err)
	}
	if result.Completed != 1 || result.Failed != 0 || result.Skipped != 0 {
		t.Fatalf("expected anomaly Task to complete after verification, got %+v", result)
	}
	if got := strings.Join(*fixture.calls, ">"); got != "agent>verify>commit" {
		t.Fatalf("expected anomaly to preserve Task flow, got %q", got)
	}
	if got := taskStatusOnDisk(t, fixture.gitRoot, "task_01"); got != string(spec.StatusCompleted) {
		t.Fatalf("expected Daemon to settle completed after verification, got %q", got)
	}
	if len(committer.messages) != 1 {
		t.Fatalf("expected Task commit after passing verification, got %v", committer.messages)
	}
	if got := strings.Join(committer.paths[0], "|"); got != taskFileRel(taskCycleSlug, "task_01")+"|internal/agent/fix.go" {
		t.Fatalf("expected task file and Agent change committed, got %q", got)
	}
	assertTaskAnomalyBeforeVerification(t, fixture.sink.snapshot(), "task_01", anomaly)
}

func TestTaskCycleTransportAnomalyStillLetsVerificationGateSettleFailure(t *testing.T) {
	fixture := newTaskCycleFixture(t, []taskSpecSeed{{id: "task_01", verification: []string{"make verify"}}})
	const anomaly = "acpx exited with exit code 1 after parsed session/prompt result\n--- acpx stderr tail ---\nMessage buffer exceeded 10485760 bytes"
	runner := &taskFakeRunner{
		calls:         fixture.calls,
		gitRoot:       fixture.gitRoot,
		statusByTask:  map[string]spec.Status{"task_01": spec.StatusCompleted},
		anomalyByTask: map[string]string{"task_01": anomaly},
	}
	verifier := &taskFakeVerifier{
		calls:  fixture.calls,
		failOn: map[string]error{"make verify": errors.New("gate broke")},
	}
	committer := &engineFakeCommitter{calls: fixture.calls}
	engine := fixture.engine(t, runner, verifier, committer, fixture.worktree)

	result, err := engine.TaskCycle(context.Background(), fixture.plan())

	if err != nil {
		t.Fatalf("expected failed verification to fail only the Task, got %v", err)
	}
	if result.Completed != 0 || result.Failed != 1 || result.Skipped != 0 {
		t.Fatalf("expected verification failure after anomaly to settle failed, got %+v", result)
	}
	if got := strings.Join(*fixture.calls, ">"); got != "agent>verify>agent>verify" {
		t.Fatalf("expected one repair attempt and no commit after final failed verification, got %q", got)
	}
	if got := taskStatusOnDisk(t, fixture.gitRoot, "task_01"); got != string(spec.StatusFailed) {
		t.Fatalf("expected Daemon to settle failed after verification failure, got %q", got)
	}
	if len(committer.messages) != 0 {
		t.Fatalf("expected no commit after failed verification, got %v", committer.messages)
	}
	assertTaskAnomalyBeforeVerification(t, fixture.sink.snapshot(), "task_01", anomaly)
}

func TestTaskCycleVerificationFailureRepairsSameSessionAndRerunsFullSequence(t *testing.T) {
	fixture := newTaskCycleFixture(t, []taskSpecSeed{{id: "task_01", verification: []string{"echo first", "echo second"}}})
	runner := &taskFakeRunner{
		calls:        fixture.calls,
		gitRoot:      fixture.gitRoot,
		statusByTask: map[string]spec.Status{"task_01": spec.StatusCompleted},
	}
	verifier := &taskFakeVerifier{
		calls: fixture.calls,
		script: []error{
			nil,
			errors.New("exit status 9"),
			nil,
			nil,
		},
	}
	committer := &engineFakeCommitter{calls: fixture.calls}
	engine := fixture.engine(t, runner, verifier, committer, fixture.worktree)
	plan := fixture.plan()

	result, err := engine.TaskCycle(context.Background(), plan)

	if err != nil {
		t.Fatalf("task cycle: %v", err)
	}
	if result.Completed != 1 || result.Failed != 0 || result.Skipped != 0 {
		t.Fatalf("expected repaired Task completed, got %+v", result)
	}
	if got := strings.Join(*fixture.calls, ">"); got != "agent>verify>verify>agent>verify>verify>commit" {
		t.Fatalf("expected initial Agent, attempt 1, repair, full final attempt, commit; got %q", got)
	}
	if got := strings.Join(verifier.commands, "|"); got != "echo first|echo second|echo first|echo second" {
		t.Fatalf("expected final attempt to rerun every configured command, got %q", got)
	}
	if len(runner.requests) != 2 {
		t.Fatalf("expected initial prompt plus one Verification Feedback prompt, got %d", len(runner.requests))
	}
	expectedSession := agent.SessionRefForRun(fixture.run.ID, fixture.gitRoot)
	if runner.requests[0].Session != expectedSession || runner.requests[1].Session != expectedSession {
		t.Fatalf("expected repair to reuse SessionRef %#v, got %#v then %#v", expectedSession, runner.requests[0].Session, runner.requests[1].Session)
	}
	repairPrompt := runner.requests[1].Prompt
	expectedPath := VerificationOutputPath(fixture.artifactDir, fixture.run.ID, 1, 1)
	for _, expected := range []string{
		"Verification Feedback",
		"Work Item: task_01",
		"Failed command: echo second",
		"Diagnostic artifact: " + expectedPath,
		"verification failed",
	} {
		if !strings.Contains(repairPrompt, expected) {
			t.Fatalf("expected repair prompt to contain %q, got:\n%s", expected, repairPrompt)
		}
	}
	if got := taskStatusOnDisk(t, fixture.gitRoot, "task_01"); got != string(spec.StatusCompleted) {
		t.Fatalf("expected repaired Task settled completed, got %q", got)
	}
}

func TestTaskCycleRepairAgentErrorFailsTaskWithoutFinalVerification(t *testing.T) {
	fixture := newTaskCycleFixture(t, []taskSpecSeed{{id: "task_01", verification: []string{"make verify"}}})
	runner := &taskFakeRunner{
		calls:         fixture.calls,
		gitRoot:       fixture.gitRoot,
		statusByTask:  map[string]spec.Status{"task_01": spec.StatusCompleted},
		errByTaskCall: map[string][]error{"task_01": {nil, errors.New("repair crashed")}},
	}
	verifier := &taskFakeVerifier{calls: fixture.calls, script: []error{errors.New("gate broke")}}
	committer := &engineFakeCommitter{calls: fixture.calls}
	engine := fixture.engine(t, runner, verifier, committer, fixture.worktree)

	result, err := engine.TaskCycle(context.Background(), fixture.plan())

	if err != nil {
		t.Fatalf("expected repair Agent error to fail only the Task, got %v", err)
	}
	if result.Completed != 0 || result.Failed != 1 || result.Skipped != 0 {
		t.Fatalf("expected one failed Task, got %+v", result)
	}
	if got := strings.Join(*fixture.calls, ">"); got != "agent>verify>agent" {
		t.Fatalf("expected no final Verification after repair Agent error, got %q", got)
	}
	if got := taskStatusOnDisk(t, fixture.gitRoot, "task_01"); got != string(spec.StatusFailed) {
		t.Fatalf("expected repair Agent error settled failed, got %q", got)
	}
	if len(committer.messages) != 0 {
		t.Fatalf("expected no commit after repair Agent error, got %v", committer.messages)
	}
}

func TestTaskCycleAgentStatusFailedReachesDaemonVerificationAndUnblocksDependent(t *testing.T) {
	fixture := newTaskCycleFixture(t, []taskSpecSeed{
		{id: "task_01", title: "Agent writes failed"},
		{id: "task_02", title: "Depends on authoritative settlement", needs: []string{"task_01"}},
		{id: "task_03", title: "Independent cleanup", taskType: "chore"},
	})
	runner := &taskFakeRunner{
		calls:   fixture.calls,
		gitRoot: fixture.gitRoot,
		statusByTask: map[string]spec.Status{
			"task_01": spec.StatusFailed,
			"task_03": spec.StatusCompleted,
		},
	}
	verifier := &taskFakeVerifier{calls: fixture.calls}
	committer := &engineFakeCommitter{calls: fixture.calls}
	engine := fixture.engine(t, runner, verifier, committer, fixture.worktree)

	result, err := engine.TaskCycle(context.Background(), fixture.plan())

	if err != nil {
		t.Fatalf("task cycle: %v", err)
	}
	if result.Completed != 3 || result.Failed != 0 || result.Skipped != 0 {
		t.Fatalf("expected Agent-authored failed status to be ignored after passing Verification, got %+v", result)
	}
	if got := strings.Join(*fixture.calls, ">"); got != "agent>verify>commit>agent>verify>commit>agent>verify>commit" {
		t.Fatalf("expected every Task, including the Agent-authored failed Task, to verify and commit, got %q", got)
	}
	if got := strings.Join(verifier.commands, "|"); got != "true|true|true" {
		t.Fatalf("expected all Tasks to reach Daemon Verification, got %q", got)
	}
	if got := taskStatusOnDisk(t, fixture.gitRoot, "task_01"); got != string(spec.StatusCompleted) {
		t.Fatalf("expected task_01 Daemon-settled completed, got %q", got)
	}
	if got := taskStatusOnDisk(t, fixture.gitRoot, "task_02"); got != string(spec.StatusCompleted) {
		t.Fatalf("expected dependent task_02 to run after authoritative settlement, got %q", got)
	}
	if got := taskStatusOnDisk(t, fixture.gitRoot, "task_03"); got != string(spec.StatusCompleted) {
		t.Fatalf("expected independent task_03 completed, got %q", got)
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
	verifier := &taskFakeVerifier{calls: fixture.calls, failOn: map[string]error{"echo t1": errors.New("exit status 7")}}
	committer := &engineFakeCommitter{calls: fixture.calls}
	engine := fixture.engine(t, runner, verifier, committer, fixture.worktree)

	result, err := engine.TaskCycle(context.Background(), fixture.plan())

	if err != nil {
		t.Fatalf("expected failed Task to not halt the cycle, got %v", err)
	}
	if result.Completed != 1 || result.Failed != 1 || result.Skipped != 1 {
		t.Fatalf("expected 1 completed, 1 failed, 1 skipped, got %+v", result)
	}
	if got := strings.Join(*fixture.calls, ">"); got != "agent>verify>agent>verify>agent>verify>commit" {
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
	if len(committer.messages) != 1 || !strings.HasPrefix(committer.messages[0], "chore: tidy the build files") {
		t.Fatalf("expected only the independent Task committed with chore mapping, got %v", committer.messages)
	}
	taskOneOutcome, ok := taskOutcomeByID(result.Outcomes, "task_01")
	if !ok {
		t.Fatalf("expected task_01 outcome, got %+v", result.Outcomes)
	}
	if taskOneOutcome.Status != string(spec.StatusFailed) {
		t.Fatalf("expected task_01 failed outcome, got %+v", taskOneOutcome)
	}
	taskTwoOutcome, ok := taskOutcomeByID(result.Outcomes, "task_02")
	if !ok {
		t.Fatalf("expected task_02 outcome, got %+v", result.Outcomes)
	}
	if taskTwoOutcome.Status != string(taskRunSkipped) || taskTwoOutcome.Reason != "needs not completed: task_01" {
		t.Fatalf("expected task_02 skipped outcome with unmet need reason, got %+v", taskTwoOutcome)
	}
	taskThreeOutcome, ok := taskOutcomeByID(result.Outcomes, "task_03")
	if !ok {
		t.Fatalf("expected task_03 outcome, got %+v", result.Outcomes)
	}
	if taskThreeOutcome.Status != string(spec.StatusCompleted) || taskThreeOutcome.Reason != "" {
		t.Fatalf("expected task_03 completed outcome without reason, got %+v", taskThreeOutcome)
	}
	if runner.requests[2].Batch.Number != 2 {
		t.Fatalf("expected skipped Task to not consume an execution ordinal, got Batch %d", runner.requests[2].Batch.Number)
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
	if got := eventPayloadString(t, *skipEvent, "reason"); got != taskTwoOutcome.Reason {
		t.Fatalf("expected skipped outcome reason to match journal payload %q, got %q", got, taskTwoOutcome.Reason)
	}
	var failedEvent *runevent.RunEvent
	for _, event := range taskEventsOfKind(fixture.sink, runevent.KindDaemonTask) {
		if event.ReviewIssue == "task_01" && eventPayloadString(t, event, "phase") == "settled" {
			matched := event
			failedEvent = &matched
		}
	}
	if failedEvent == nil {
		t.Fatal("expected failed settlement event journaling the verification failure reason")
	}
	if got := eventPayloadString(t, *failedEvent, "reason"); got != taskOneOutcome.Reason {
		t.Fatalf("expected failed outcome reason to match journal payload %q, got %q", got, taskOneOutcome.Reason)
	}
	for _, expected := range []string{`Verification failed: command "echo t1"`, "exit status", "diagnostics:"} {
		if !strings.Contains(taskOneOutcome.Reason, expected) {
			t.Fatalf("expected verification outcome reason to contain %q, got %q", expected, taskOneOutcome.Reason)
		}
	}
	kinds := fixture.sink.kinds()
	last := taskEventsOfKind(fixture.sink, runevent.KindDaemonOutcome)
	if kinds[len(kinds)-1] != runevent.KindDaemonOutcome || !strings.Contains(string(last[0].Payload), `"skipped":1`) {
		t.Fatalf("expected outcome event with counts at cycle end, got %v", kinds)
	}
}

func TestTaskCycleTemporaryVerificationFlowPassesExclusiveRetryWithoutAgentRepair(t *testing.T) {
	fixture := newTaskCycleFixture(t, []taskSpecSeed{{id: "task_01", verification: []string{"verify task"}}})
	runner := &taskFakeRunner{calls: fixture.calls, gitRoot: fixture.gitRoot}
	verifier := &taskFakeVerifier{calls: fixture.calls, temporaryOnCall: map[int]bool{1: true}}
	engine := fixture.engine(t, runner, verifier, &engineFakeCommitter{calls: fixture.calls}, fixture.worktree)

	result, err := engine.TaskCycle(context.Background(), fixture.plan())

	if err != nil {
		t.Fatalf("task cycle: %v", err)
	}
	if result.Completed != 1 || result.Failed != 0 {
		t.Fatalf("expected exclusive retry to complete Task, got %+v", result)
	}
	if got := strings.Join(*fixture.calls, ">"); got != "agent>verify>verify>commit" {
		t.Fatalf("expected zero Agent repair turns around exclusive retry, got %q", got)
	}
	initialPath := VerificationOutputPath(fixture.artifactDir, fixture.run.ID, 1, 1)
	retryPath := VerificationRetryOutputPath(fixture.artifactDir, fixture.run.ID, 1, 1, 1)
	if got := strings.Join(verifier.outputPaths, "|"); got != initialPath+"|"+retryPath {
		t.Fatalf("expected distinct initial and retry diagnostics, got %q", got)
	}

	var sawTemporary, sawExclusiveWaiting, sawRetryPassed bool
	for _, event := range taskEventsOfKind(fixture.sink, runevent.KindDaemonVerification) {
		payload := eventPayloadMap(t, event)
		switch payload["phase"] {
		case string(runevent.VerificationPhaseFailed):
			if payload["classification"] == string(runevent.VerificationClassificationTemporary) &&
				payload["retry_available"] == true &&
				payload["diagnostic_path"] == initialPath {
				sawTemporary = true
			}
		case string(runevent.VerificationPhaseWaiting):
			if payload["retry"] == float64(1) && payload["mode"] == "exclusive" {
				sawExclusiveWaiting = true
			}
		case string(runevent.VerificationPhaseVerdict):
			if payload["retry"] == float64(1) && payload["verdict"] == string(runevent.VerificationVerdictPassed) {
				sawRetryPassed = true
			}
		}
	}
	if !sawTemporary || !sawExclusiveWaiting || !sawRetryPassed {
		t.Fatalf("expected temporary, exclusive-wait, and retry-pass evidence; temporary=%v waiting=%v passed=%v", sawTemporary, sawExclusiveWaiting, sawRetryPassed)
	}
}

func TestTaskCycleDeterministicRetryUsesAgentRepairThenAttemptTwo(t *testing.T) {
	fixture := newTaskCycleFixture(t, []taskSpecSeed{{id: "task_01", verification: []string{"verify task"}}})
	runner := &taskFakeRunner{calls: fixture.calls, gitRoot: fixture.gitRoot}
	verifier := &taskFakeVerifier{
		calls:           fixture.calls,
		temporaryOnCall: map[int]bool{1: true},
		script:          []error{errors.New("exit status 1"), nil},
	}
	engine := fixture.engine(t, runner, verifier, &engineFakeCommitter{calls: fixture.calls}, fixture.worktree)

	result, err := engine.TaskCycle(context.Background(), fixture.plan())

	if err != nil {
		t.Fatalf("task cycle: %v", err)
	}
	if result.Completed != 1 || result.Failed != 0 {
		t.Fatalf("expected deterministic retry failure repaired by attempt 2, got %+v", result)
	}
	if got := strings.Join(*fixture.calls, ">"); got != "agent>verify>verify>agent>verify>commit" {
		t.Fatalf("expected one Agent repair after deterministic exclusive retry, got %q", got)
	}
	if runner.taskCalls["task_01"] != 2 {
		t.Fatalf("expected initial Agent plus one Verification Feedback turn, got %d", runner.taskCalls["task_01"])
	}
	if len(verifier.outputPaths) != 3 ||
		verifier.outputPaths[1] != VerificationRetryOutputPath(fixture.artifactDir, fixture.run.ID, 1, 1, 1) ||
		verifier.outputPaths[2] != VerificationOutputPath(fixture.artifactDir, fixture.run.ID, 1, 2) {
		t.Fatalf("expected retry then numbered attempt 2 paths, got %v", verifier.outputPaths)
	}
}

func TestTaskCycleRetryBudgetExhaustsOnRepeatedTemporaryVerification(t *testing.T) {
	fixture := newTaskCycleFixture(t, []taskSpecSeed{{id: "task_01", verification: []string{"verify task"}}})
	runner := &taskFakeRunner{calls: fixture.calls, gitRoot: fixture.gitRoot}
	verifier := &taskFakeVerifier{calls: fixture.calls, temporaryOnCall: map[int]bool{1: true, 2: true}}
	engine := fixture.engine(t, runner, verifier, &engineFakeCommitter{calls: fixture.calls}, fixture.worktree)

	result, err := engine.TaskCycle(context.Background(), fixture.plan())

	if err != nil {
		t.Fatalf("task cycle: %v", err)
	}
	if result.Completed != 0 || result.Failed != 1 {
		t.Fatalf("expected repeated temporary failure to settle failed, got %+v", result)
	}
	if got := strings.Join(*fixture.calls, ">"); got != "agent>verify>verify" {
		t.Fatalf("expected no Agent repair or second retry after exhaustion, got %q", got)
	}
	var sawExhaustion bool
	for _, event := range taskEventsOfKind(fixture.sink, runevent.KindDaemonVerification) {
		payload := eventPayloadMap(t, event)
		if payload["phase"] == string(runevent.VerificationPhaseFailed) &&
			payload["retry"] == float64(1) &&
			payload["classification"] == string(runevent.VerificationClassificationTemporary) &&
			payload["retry_available"] == false &&
			payload["reason"] == string(runevent.VerificationReasonTemporaryFailure) {
			sawExhaustion = true
		}
	}
	if !sawExhaustion {
		t.Fatal("expected bounded temporary retry exhaustion evidence")
	}
}

func TestTaskCycleTemporaryVerificationAfterDeterministicRetryAndRepairDoesNotRetryAgain(t *testing.T) {
	fixture := newTaskCycleFixture(t, []taskSpecSeed{{id: "task_01", verification: []string{"verify task"}}})
	runner := &taskFakeRunner{calls: fixture.calls, gitRoot: fixture.gitRoot}
	verifier := &taskFakeVerifier{
		calls:           fixture.calls,
		temporaryOnCall: map[int]bool{1: true, 3: true},
		script:          []error{errors.New("exit status 1")},
	}
	engine := fixture.engine(t, runner, verifier, &engineFakeCommitter{calls: fixture.calls}, fixture.worktree)

	result, err := engine.TaskCycle(context.Background(), fixture.plan())

	if err != nil {
		t.Fatalf("task cycle: %v", err)
	}
	if result.Completed != 0 || result.Failed != 1 {
		t.Fatalf("expected attempt-2 temporary failure with spent retry budget to settle failed, got %+v", result)
	}
	if got := strings.Join(*fixture.calls, ">"); got != "agent>verify>verify>agent>verify" {
		t.Fatalf("expected one retry and one Agent repair only, got %q", got)
	}
	if len(verifier.outputPaths) != 3 {
		t.Fatalf("expected no second exclusive retry, got paths %v", verifier.outputPaths)
	}
}

func TestTaskCycleVerificationSequenceStopsAtFirstFailureWithOneVerdictPerAttempt(t *testing.T) {
	fixture := newTaskCycleFixture(t, []taskSpecSeed{{id: "task_01", verification: []string{"echo first", "echo second", "echo third"}}})
	runner := &taskFakeRunner{calls: fixture.calls, gitRoot: fixture.gitRoot}
	verifier := &taskFakeVerifier{calls: fixture.calls, failOn: map[string]error{"echo second": errors.New("exit status 9")}}
	engine := fixture.engine(t, runner, verifier, &engineFakeCommitter{calls: fixture.calls}, fixture.worktree)

	result, err := engine.TaskCycle(context.Background(), fixture.plan())

	if err != nil {
		t.Fatalf("expected command failure to fail only the Task, got %v", err)
	}
	if result.Completed != 0 || result.Failed != 1 || result.Skipped != 0 {
		t.Fatalf("expected one failed Task, got %+v", result)
	}
	if got := strings.Join(verifier.commands, "|"); got != "echo first|echo second|echo first|echo second" {
		t.Fatalf("expected verification to stop on the first failing command, got %q", got)
	}
	expectedPath := VerificationOutputPath(fixture.artifactDir, fixture.run.ID, 1, 1)
	finalPath := VerificationOutputPath(fixture.artifactDir, fixture.run.ID, 1, 2)
	if len(verifier.outputPaths) != 4 ||
		verifier.outputPaths[0] != expectedPath ||
		verifier.outputPaths[1] != expectedPath ||
		verifier.outputPaths[2] != finalPath ||
		verifier.outputPaths[3] != finalPath {
		t.Fatalf("expected attempt diagnostic paths %q then %q, got %v", expectedPath, finalPath, verifier.outputPaths)
	}
	if progress := fixture.progress.String(); !strings.Contains(progress, "Verification failed (attempt 1); diagnostics: "+expectedPath) {
		t.Fatalf("expected bounded verdict summary without raw failure text, got %q", progress)
	}
	if progress := fixture.progress.String(); !strings.Contains(progress, "Verification failed (attempt 2); diagnostics: "+finalPath) {
		t.Fatalf("expected final bounded verdict summary without raw failure text, got %q", progress)
	}
	verificationEvents := taskEventsOfKind(fixture.sink, runevent.KindDaemonVerification)
	phases := []string{}
	verdicts := 0
	for _, event := range verificationEvents {
		if event.ReviewIssue != "task_01" || event.Batch != 1 {
			t.Fatalf("expected task verification event to carry Work Item and Batch, got %+v", event)
		}
		payload := eventPayloadMap(t, event)
		phase, _ := payload["phase"].(string)
		phases = append(phases, phase)
		attempt := int(payload["attempt"].(float64))
		expectedDiagnosticPath := VerificationOutputPath(fixture.artifactDir, fixture.run.ID, 1, attempt)
		if phase == string(runevent.VerificationPhaseFailed) && payload["diagnostic_path"] != expectedDiagnosticPath {
			t.Fatalf("expected failed command diagnostic path %q, got %v", expectedDiagnosticPath, payload["diagnostic_path"])
		}
		if phase == string(runevent.VerificationPhaseVerdict) {
			verdicts++
			if payload["verdict"] != string(runevent.VerificationVerdictFailed) {
				t.Fatalf("expected failed verdict, got %v", payload["verdict"])
			}
		}
	}
	if got := strings.Join(phases, "|"); got != "waiting|started|command-passed|started|failed|verdict|waiting|started|command-passed|started|failed|verdict" {
		t.Fatalf("expected command phases and one aggregate verdict per attempt, got %s", got)
	}
	if verdicts != 2 {
		t.Fatalf("expected exactly one aggregate verdict per attempt, got %d", verdicts)
	}
}

func TestTaskCycleDaemonStatusRemainsInProgressOnVerificationInfrastructureError(t *testing.T) {
	fixture := newTaskCycleFixture(t, []taskSpecSeed{{id: "task_01"}})
	infraErr := errors.New("diagnostic artifact write failed")
	runner := &taskFakeRunner{calls: fixture.calls, gitRoot: fixture.gitRoot, statusByTask: map[string]spec.Status{"task_01": spec.StatusCompleted}}
	engine := fixture.engine(t, runner, &engineInfrastructureVerifier{calls: fixture.calls, err: infraErr}, &engineFakeCommitter{calls: fixture.calls}, fixture.worktree)

	result, err := engine.TaskCycle(context.Background(), fixture.plan())

	if !errors.Is(err, infraErr) {
		t.Fatalf("expected infrastructure error identity preserved, got %v", err)
	}
	if result.Completed != 0 || result.Failed != 0 || result.Skipped != 0 {
		t.Fatalf("expected no Task settlement on infrastructure error, got %+v", result)
	}
	if got := strings.Join(*fixture.calls, ">"); got != "agent>verify" {
		t.Fatalf("expected Task cycle to halt after verification infrastructure error, got %q", got)
	}
	if got := taskStatusOnDisk(t, fixture.gitRoot, "task_01"); got != string(spec.StatusInProgress) {
		t.Fatalf("expected Daemon-owned in_progress status preserved on infrastructure error, got %q", got)
	}
	verificationEvents := taskEventsOfKind(fixture.sink, runevent.KindDaemonVerification)
	verdicts := 0
	for _, event := range verificationEvents {
		payload := eventPayloadMap(t, event)
		if payload["phase"] == string(runevent.VerificationPhaseVerdict) {
			verdicts++
			if payload["verdict"] != string(runevent.VerificationVerdictFailed) {
				t.Fatalf("expected failed infrastructure verdict, got %v", payload["verdict"])
			}
		}
	}
	if verdicts != 1 {
		t.Fatalf("expected one aggregate verdict before halting, got %d", verdicts)
	}
}

func TestTaskCycleDaemonStatusUnreadableAgentArtifactSettlesFailedWithoutVerification(t *testing.T) {
	fixture := newTaskCycleFixture(t, []taskSpecSeed{{id: "task_01"}})
	runner := &taskFakeRunner{
		calls:           fixture.calls,
		gitRoot:         fixture.gitRoot,
		rawStatusByTask: map[string]string{"task_01": "not-a-status"},
	}
	engine := fixture.engine(t, runner, &taskFakeVerifier{calls: fixture.calls}, &engineFakeCommitter{calls: fixture.calls}, fixture.worktree)

	result, err := engine.TaskCycle(context.Background(), fixture.plan())

	if err != nil {
		t.Fatalf("expected unreadable Agent artifact to fail only the Task, got %v", err)
	}
	if result.Completed != 0 || result.Failed != 1 || result.Skipped != 0 {
		t.Fatalf("expected one Daemon-settled failed Task, got %+v", result)
	}
	if got := strings.Join(*fixture.calls, ">"); got != "agent" {
		t.Fatalf("expected unreadable Task artifact to start no Verification, got %q", got)
	}
	if len(result.Outcomes) != 1 || !strings.Contains(result.Outcomes[0].Reason, "reload task file after the Agent") {
		t.Fatalf("expected unreadable artifact reason distinct from Agent and Verification failures, got %+v", result.Outcomes)
	}
	if got := taskStatusOnDisk(t, fixture.gitRoot, "task_01"); got != string(spec.StatusFailed) {
		t.Fatalf("expected only the Daemon to settle the unreadable Task failed, got %q", got)
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
		if got := strings.Join(verifier.commands, "|"); got != "echo gamma|echo gamma" {
			t.Fatalf("expected verification to stop at the first failing command, got %q", got)
		}
		if len(committer.messages) != 0 {
			t.Fatalf("expected no commit for a failed Task, got %v", committer.messages)
		}
	})
}

func TestTaskCycleAgentFailureStartsNoVerificationAndSettlesFailed(t *testing.T) {
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
	if len(result.Outcomes) != 2 || result.Outcomes[0].Reason != "Agent failed: agent crashed" {
		t.Fatalf("expected existing Agent failure reason to stay unchanged, got %+v", result.Outcomes)
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

func TestTaskCycleModelNotAdvertisedFailureSettlesAndReportsReason(t *testing.T) {
	fixture := newTaskCycleFixture(t, []taskSpecSeed{
		{id: "task_01", title: "Use the selected Agent Model"},
	})
	runner := &taskFakeRunner{
		calls:     fixture.calls,
		gitRoot:   fixture.gitRoot,
		errByTask: map[string]error{"task_01": modelNotAdvertisedBatchErrorForTest()},
	}
	engine := fixture.engine(t, runner, &taskFakeVerifier{calls: fixture.calls}, &engineFakeCommitter{calls: fixture.calls}, fixture.worktree)

	result, err := engine.TaskCycle(context.Background(), fixture.plan())

	if err != nil {
		t.Fatalf("expected model rejection to fail only the Task, got %v", err)
	}
	if result.Completed != 0 || result.Failed != 1 || result.Skipped != 0 {
		t.Fatalf("expected one failed Task, got %+v", result)
	}
	if len(result.Outcomes) != 1 {
		t.Fatalf("expected one Task outcome, got %+v", result.Outcomes)
	}
	outcome := result.Outcomes[0]
	if outcome.Task != "task_01" || outcome.Status != string(spec.StatusFailed) || outcome.Reason != modelNotAdvertisedReasonForTest {
		t.Fatalf("expected model rejection report outcome, got %+v", outcome)
	}
	if got := taskStatusOnDisk(t, fixture.gitRoot, "task_01"); got != string(spec.StatusFailed) {
		t.Fatalf("expected model rejection settled failed, got %q", got)
	}
	if !strings.Contains(fixture.progress.String(), "Task task_01 failed: "+modelNotAdvertisedReasonForTest+"\n") {
		t.Fatalf("expected progress report to include model rejection reason, got %q", fixture.progress.String())
	}
	for _, event := range taskEventsOfKind(fixture.sink, runevent.KindDaemonTask) {
		if event.ReviewIssue == "task_01" && eventPayloadString(t, event, "reason") == modelNotAdvertisedReasonForTest {
			return
		}
	}
	t.Fatalf("expected Task settlement event to journal reason %q", modelNotAdvertisedReasonForTest)
}

func TestTaskCycleSpecRootOnlyTaskCommitWarnsAndStillCommits(t *testing.T) {
	fixture := newTaskCycleFixture(t, []taskSpecSeed{{id: "task_01", title: "Only settle the task file"}})
	fixture.worktree.snapshots = [][]string{nil, nil}
	runner := &taskFakeRunner{
		calls:        fixture.calls,
		gitRoot:      fixture.gitRoot,
		statusByTask: map[string]spec.Status{"task_01": spec.StatusCompleted},
	}
	committer := &engineFakeCommitter{calls: fixture.calls}
	engine := fixture.engine(t, runner, &taskFakeVerifier{calls: fixture.calls}, committer, fixture.worktree)

	result, err := engine.TaskCycle(context.Background(), fixture.plan())

	if err != nil {
		t.Fatalf("task cycle: %v", err)
	}
	if result.Completed != 1 || result.Failed != 0 || result.Skipped != 0 {
		t.Fatalf("expected one completed Task, got %+v", result)
	}
	if len(committer.paths) != 1 {
		t.Fatalf("expected one spec-root-only Task commit, got %v", committer.paths)
	}
	if got := strings.Join(committer.paths[0], "|"); got != taskFileRel(taskCycleSlug, "task_01") {
		t.Fatalf("expected only the task file committed, got %q", got)
	}
	if got := taskStatusOnDisk(t, fixture.gitRoot, "task_01"); got != string(spec.StatusCompleted) {
		t.Fatalf("expected Task settled completed, got %q", got)
	}
	assertNoOpTaskCommitWarning(t, fixture, "task_01", taskNoOpShapeSpecRootOnly)
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
	assertNoNoOpTaskCommitWarning(t, fixture)
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

func TestTaskCycleStopRequestAfterTaskSettlementHaltsBeforeNextTask(t *testing.T) {
	fixture := newTaskCycleFixture(t, []taskSpecSeed{
		{id: "task_01", title: "Build the core"},
		{id: "task_02", title: "Wire the shell"},
	})
	fixture.worktree.snapshots = [][]string{nil, {"src/core.go"}, nil, {"src/shell.go"}}
	runner := &taskFakeRunner{
		calls:   fixture.calls,
		gitRoot: fixture.gitRoot,
		statusByTask: map[string]spec.Status{
			"task_01": spec.StatusCompleted,
			"task_02": spec.StatusCompleted,
		},
	}
	committer := &engineFakeCommitter{
		calls: fixture.calls,
		afterCommit: func(context.Context, CommitRequest) error {
			return fixture.store.RequestStop(context.Background(), fixture.run.ID)
		},
	}
	engine := fixture.engine(t, runner, &taskFakeVerifier{calls: fixture.calls}, committer, fixture.worktree)

	result, err := engine.TaskCycle(context.Background(), fixture.plan())

	if !errors.Is(err, ErrStopRequested) {
		t.Fatalf("expected ErrStopRequested, got %v", err)
	}
	if got := strings.Join(*fixture.calls, ">"); got != "agent>verify>commit" {
		t.Fatalf("expected first Task to verify and commit before stop, got %q", got)
	}
	if result.Completed != 1 || result.Failed != 0 || result.Skipped != 0 {
		t.Fatalf("expected one settled completed Task before stop, got %+v", result)
	}
	if got := taskStatusOnDisk(t, fixture.gitRoot, "task_01"); got != string(spec.StatusCompleted) {
		t.Fatalf("expected task_01 completed before stop, got %q", got)
	}
	if got := taskStatusOnDisk(t, fixture.gitRoot, "task_02"); got != string(spec.StatusPending) {
		t.Fatalf("expected task_02 left pending, got %q", got)
	}
	if len(committer.messages) != 1 || !strings.HasPrefix(committer.messages[0], "feat: build the core") {
		t.Fatalf("expected exactly the completed Task committed, got %v", committer.messages)
	}
	kinds := fixture.sink.kinds()
	if len(kinds) == 0 || kinds[len(kinds)-1] != runevent.KindDaemonStatus {
		t.Fatalf("expected daemon stop event after Task settlement, got %v", kinds)
	}
}

func TestTaskCycleStopRequestBeforeQAStepSkipsQA(t *testing.T) {
	fixture := newTaskCycleFixture(t, []taskSpecSeed{{id: "task_01", title: "Build the feature"}})
	fixture.worktree.snapshots = [][]string{nil, {"src/feature.go"}, {"src/feature.go"}, {"src/feature.go", qaReportRelPathForTest()}}
	runner := &taskFakeRunner{
		calls:        fixture.calls,
		gitRoot:      fixture.gitRoot,
		statusByTask: map[string]spec.Status{"task_01": spec.StatusCompleted},
		qaReport:     qaReportForTest(spec.VerdictPass),
	}
	committer := &engineFakeCommitter{
		calls: fixture.calls,
		afterCommit: func(context.Context, CommitRequest) error {
			return fixture.store.RequestStop(context.Background(), fixture.run.ID)
		},
	}
	engine := fixture.engine(t, runner, &taskFakeVerifier{calls: fixture.calls}, committer, fixture.worktree)

	result, err := engine.TaskCycle(context.Background(), fixture.qaPlan())

	if !errors.Is(err, ErrStopRequested) {
		t.Fatalf("expected ErrStopRequested, got %v", err)
	}
	if got := strings.Join(*fixture.calls, ">"); got != "agent>verify>commit" {
		t.Fatalf("expected Task settlement only, with QA skipped, got %q", got)
	}
	if result.Completed != 1 || result.Failed != 0 || result.Skipped != 0 || result.QAVerdict != "" || result.QAReportPath != "" {
		t.Fatalf("expected completed Task and no QA verdict, got %+v", result)
	}
	if len(runner.qaPrompts) != 0 {
		t.Fatalf("expected Stop Request to skip QA, got %d QA prompt(s)", len(runner.qaPrompts))
	}
	if events := taskEventsOfKind(fixture.sink, runevent.KindDaemonQA); len(events) != 0 {
		t.Fatalf("expected no daemon.qa event when Stop Request precedes QA, got %+v", events)
	}
	kinds := fixture.sink.kinds()
	if len(kinds) == 0 || kinds[len(kinds)-1] != runevent.KindDaemonStatus {
		t.Fatalf("expected daemon stop event before QA, got %v", kinds)
	}
}

func TestTaskCycleStopDuringAgentPreservesDaemonStatusAndHalts(t *testing.T) {
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
	// Worktree preserved: the stopped Task keeps the Daemon's resumable status.
	if got := taskStatusOnDisk(t, fixture.gitRoot, "task_01"); got != string(spec.StatusInProgress) {
		t.Fatalf("expected stopped Task left resumable in_progress, got %q", got)
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

func TestTaskCycleStopAfterAgentStatusAuthorshipPreservesResultInProgress(t *testing.T) {
	fixture := newTaskCycleFixture(t, []taskSpecSeed{{id: "task_01"}})
	ctx, cancel := context.WithCancel(context.Background())
	const resultEvidence = "\n## Result\n\n- Agent work is implementation-ready.\n"
	runner := &taskFakeRunner{
		calls:        fixture.calls,
		gitRoot:      fixture.gitRoot,
		statusByTask: map[string]spec.Status{"task_01": spec.StatusCompleted},
		resultByTask: map[string]string{"task_01": resultEvidence},
		afterTask:    func(string) { cancel() },
	}
	engine := fixture.engine(t, runner, &taskFakeVerifier{calls: fixture.calls}, &engineFakeCommitter{calls: fixture.calls}, fixture.worktree)

	result, err := engine.TaskCycle(ctx, fixture.plan())

	if !errors.Is(err, context.Canceled) {
		t.Fatalf("expected Stop Request after Agent work, got %v", err)
	}
	if result.Completed != 0 || result.Failed != 0 || result.Skipped != 0 {
		t.Fatalf("expected no invented terminal Task settlement, got %+v", result)
	}
	if got := strings.Join(*fixture.calls, ">"); got != "agent" {
		t.Fatalf("expected Stop Request to start no Verification or commit, got %q", got)
	}
	if got := taskStatusOnDisk(t, fixture.gitRoot, "task_01"); got != string(spec.StatusInProgress) {
		t.Fatalf("expected resumable Daemon-owned in_progress status, got %q", got)
	}
	content, readErr := os.ReadFile(taskPathFor(fixture.gitRoot, taskCycleSlug, "task_01"))
	if readErr != nil {
		t.Fatalf("read stopped Task: %v", readErr)
	}
	if !strings.Contains(string(content), resultEvidence) {
		t.Fatalf("expected Stop Request to preserve Agent Result evidence, got:\n%s", content)
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
	gittest.InitRepo(t, repoDir, "-b", "main")
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
	if firstMessage != "docs: write the usage docs\n\nRoundfix-Spec: "+taskCycleSlug+"\nRoundfix-Task: task_01" {
		t.Fatalf("expected docs commit message with trailers, got %q", firstMessage)
	}
	secondMessage := strings.TrimSpace(runGitForTest(t, repoDir, "show", "-s", "--format=%B", "HEAD"))
	if secondMessage != "feat: add the backend behavior\n\nRoundfix-Spec: "+taskCycleSlug+"\nRoundfix-Task: task_02" {
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

func TestTaskCommitDropsSymlinkCrossingTaskFileAndCommitsRepositoryPaths(t *testing.T) {
	fixture := newTaskCycleFixture(t, []taskSpecSeed{{id: "task_01", title: "Build the feature"}})
	externalRoot := filepath.Join(t.TempDir(), "knowledge-specs")
	writeSpecDirAtRootForTest(t, externalRoot, taskCycleSlug, []taskSpecSeed{{id: "task_01", title: "Build the feature"}})
	linkPath := filepath.Join(fixture.gitRoot, "docs", "specs")
	if err := os.RemoveAll(linkPath); err != nil {
		t.Fatalf("remove default Spec Root fixture: %v", err)
	}
	if err := os.Symlink(externalRoot, linkPath); err != nil {
		t.Skipf("symlink unavailable: %v", err)
	}
	fixture.worktree.snapshots = [][]string{{"src/agent-change.go"}}
	committer := &engineFakeCommitter{calls: fixture.calls}
	engine := fixture.engine(t, &taskFakeRunner{calls: fixture.calls, gitRoot: fixture.gitRoot}, &taskFakeVerifier{calls: fixture.calls}, committer, fixture.worktree)
	plan := fixture.plan()
	plan.SpecsRoot = linkPath

	err := engine.commitTask(context.Background(), plan, fixture.graph.Tasks[0], 1, nil)

	if err != nil {
		t.Fatalf("commitTask: %v", err)
	}
	if len(committer.paths) != 1 {
		t.Fatalf("expected one Task commit, got %d", len(committer.paths))
	}
	if got := strings.Join(committer.paths[0], "|"); got != "src/agent-change.go" {
		t.Fatalf("expected only repository path staged, got %q", got)
	}
	dropped := droppedStageEvents(t, fixture.sink)
	if len(dropped) != 1 {
		t.Fatalf("expected one dropped-path event, got %+v", dropped)
	}
	if got := eventPayloadString(t, dropped[0], "reason"); got != "crosses a symbolic link" {
		t.Fatalf("expected symlink reason, got %q", got)
	}
	if got := eventPayloadString(t, dropped[0], "path"); got != taskFileRel(taskCycleSlug, "task_01") {
		t.Fatalf("expected dropped task path %q, got %q", taskFileRel(taskCycleSlug, "task_01"), got)
	}
	wantWarning := "roundfix: task file " + taskFileRel(taskCycleSlug, "task_01") + " kept outside the repository; omitted from the commit\n"
	if !strings.Contains(fixture.progress.String(), wantWarning) {
		t.Fatalf("expected progress warning %q, got %q", wantWarning, fixture.progress.String())
	}
}

func TestTaskCycleSettlesCompletedWithoutCommitWhenOnlyExternalTaskFileChanged(t *testing.T) {
	fixture := newTaskCycleFixture(t, []taskSpecSeed{{id: "task_01", title: "External artifact only"}})
	externalRoot := fixture.useExternalSpecRoot(t, []taskSpecSeed{{id: "task_01", title: "External artifact only"}})
	fixture.worktree.snapshots = [][]string{nil, nil}
	runner := &taskFakeRunner{
		calls:        fixture.calls,
		gitRoot:      fixture.gitRoot,
		statusByTask: map[string]spec.Status{"task_01": spec.StatusCompleted},
	}
	committer := &engineFakeCommitter{calls: fixture.calls}
	engine := fixture.engine(t, runner, &taskFakeVerifier{calls: fixture.calls}, committer, fixture.worktree)

	result, err := engine.TaskCycle(context.Background(), fixture.plan())

	if err != nil {
		t.Fatalf("task cycle: %v", err)
	}
	if result.Completed != 1 || result.Failed != 0 || result.Skipped != 0 {
		t.Fatalf("expected the external-only Task completed, got %+v", result)
	}
	if len(committer.messages) != 0 {
		t.Fatalf("expected no commit when only the external task file changed, got %v", committer.messages)
	}
	if got := taskStatusInSpecRootOnDisk(t, externalRoot, "task_01"); got != string(spec.StatusCompleted) {
		t.Fatalf("expected external task settled completed, got %q", got)
	}
	taskEvents := taskEventsOfKind(fixture.sink, runevent.KindDaemonTask)
	if !hasTaskSettlementEvent(t, taskEvents, "task_01", spec.StatusCompleted) {
		t.Fatalf("expected normal settled event, got %+v", taskEvents)
	}
	dropped := droppedStageEvents(t, fixture.sink)
	if len(dropped) != 1 {
		t.Fatalf("expected one dropped-path event, got %+v", dropped)
	}
	if got := eventPayloadString(t, dropped[0], "reason"); got != "external to repository" {
		t.Fatalf("expected external reason, got %q", got)
	}
	wantPath := taskPathInSpecRootFor(externalRoot, taskCycleSlug, "task_01")
	if got := eventPayloadString(t, dropped[0], "path"); got != wantPath {
		t.Fatalf("expected dropped external path %q, got %q", wantPath, got)
	}
	assertNoOpTaskCommitWarning(t, fixture, "task_01", taskNoOpShapeEmptyStageable)
}

func TestTaskCycleQAReportExternalProceedsWithoutStaging(t *testing.T) {
	fixture := newTaskCycleFixture(t, []taskSpecSeed{{id: "task_01", status: string(spec.StatusCompleted)}})
	externalRoot := fixture.useExternalSpecRoot(t, []taskSpecSeed{{id: "task_01", status: string(spec.StatusCompleted)}})
	fixture.worktree.snapshots = [][]string{nil, nil}
	runner := &taskFakeRunner{
		calls:    fixture.calls,
		gitRoot:  fixture.gitRoot,
		qaReport: qaReportForTest(spec.VerdictPass),
	}
	committer := &engineFakeCommitter{calls: fixture.calls}
	engine := fixture.engine(t, runner, &taskFakeVerifier{calls: fixture.calls}, committer, fixture.worktree)

	result, err := engine.TaskCycle(context.Background(), fixture.qaPlan())

	if err != nil {
		t.Fatalf("task cycle: %v", err)
	}
	wantReportPath := filepath.Join(externalRoot, taskCycleSlug, "qa", qaReportNameForTest)
	if result.QAVerdict != spec.VerdictPass || result.QAReportPath != wantReportPath {
		t.Fatalf("expected external QA pass report %q, got %+v", wantReportPath, result)
	}
	if len(committer.messages) != 0 {
		t.Fatalf("expected no commit for an external QA Report, got %v", committer.messages)
	}
	dropped := droppedStageEvents(t, fixture.sink)
	if len(dropped) != 1 {
		t.Fatalf("expected one dropped QA report event, got %+v", dropped)
	}
	if got := eventPayloadString(t, dropped[0], "reason"); got != "external to repository" {
		t.Fatalf("expected external reason, got %q", got)
	}
	if got := eventPayloadString(t, dropped[0], "path"); got != wantReportPath {
		t.Fatalf("expected dropped QA path %q, got %q", wantReportPath, got)
	}
	if !strings.Contains(fixture.progress.String(), "roundfix: QA Report "+wantReportPath+" kept outside the repository; omitted from the commit\n") {
		t.Fatalf("expected external QA progress warning, got %q", fixture.progress.String())
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
				writeLogs:    true,
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
			expectedQALog := agent.LogPath(fixture.artifactDir, fixture.run.ID, 2)
			if last.LogPath != expectedQALog {
				t.Fatalf("expected QA Agent log under Artifact Directory %q, got %q", expectedQALog, last.LogPath)
			}
			if _, err := os.Stat(expectedQALog); err != nil {
				t.Fatalf("expected fake QA Agent log under Artifact Directory: %v", err)
			}
			if _, err := os.Stat(filepath.Join(fixture.gitRoot, ".roundfix")); !errors.Is(err, os.ErrNotExist) {
				t.Fatalf("expected no repo .roundfix directory after QA step, got err=%v", err)
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
				wantMessage := "docs: qa report for " + taskCycleSlug + " (" + tt.wantVerdict + ")\n\nRoundfix-Spec: " + taskCycleSlug
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
		statusByTask: map[string]spec.Status{"task_01": spec.StatusCompleted},
		qaReport:     qaReportForTest(spec.VerdictPass),
	}
	committer := &engineFakeCommitter{calls: fixture.calls}
	verifier := &taskFakeVerifier{calls: fixture.calls, failOn: map[string]error{"true": errors.New("gate broke")}}
	engine := fixture.engine(t, runner, verifier, committer, fixture.worktree)

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
	if len(committer.messages) != 1 || !strings.HasPrefix(committer.messages[0], "docs: qa report for "+taskCycleSlug+" (pass)") {
		t.Fatalf("expected the QA Report commit, got %v", committer.messages)
	}
}

// The QA gate reasons about the user's Pull Request from a Run Worktree
// that is structurally never on the user's branch, so the prompt carries
// both branch names and the user checkout as facts.
func TestTaskCycleQAPromptStatesRunBranchAndSpecTargetBranch(t *testing.T) {
	fixture := newTaskCycleFixture(t, []taskSpecSeed{{id: "task_01", status: string(spec.StatusCompleted)}})
	runner := &taskFakeRunner{
		calls:    fixture.calls,
		gitRoot:  fixture.gitRoot,
		qaReport: qaReportForTest(spec.VerdictPass),
	}
	engine := fixture.engine(t, runner, &taskFakeVerifier{calls: fixture.calls}, &engineFakeCommitter{calls: fixture.calls}, fixture.worktree)

	result, err := engine.TaskCycle(context.Background(), fixture.qaPlan())

	if err != nil {
		t.Fatalf("task cycle: %v", err)
	}
	if result.QAVerdict != spec.VerdictPass {
		t.Fatalf("expected the pass verdict, got %+v", result)
	}
	if len(runner.qaPrompts) != 1 {
		t.Fatalf("expected one QA prompt, got %d", len(runner.qaPrompts))
	}
	prompt := runner.qaPrompts[0]
	for _, expected := range []string{
		"Run Worktree branch: " + runworktree.BranchName(fixture.run.ID) + " (this checkout only — a per-Run branch that is never pushed and has no Pull Request of its own)\n",
		"Spec target branch: ma/spec-work (the user branch this Spec's commits land on; any Pull Request for this Spec is open on this branch, never on the Run Worktree branch)\n",
		"User checkout: " + fixture.gitRoot + " (the user's repository root this Run Worktree was created from)\n",
	} {
		if !strings.Contains(prompt, expected) {
			t.Fatalf("expected the QA prompt to state %q, got:\n%s", expected, prompt)
		}
	}
}

func TestTaskCycleQAPromptStaysUsableWithoutRecordedTargetBranch(t *testing.T) {
	fixture := newTaskCycleFixture(t, []taskSpecSeed{{id: "task_01", status: string(spec.StatusCompleted)}})
	runner := &taskFakeRunner{
		calls:    fixture.calls,
		gitRoot:  fixture.gitRoot,
		qaReport: qaReportForTest(spec.VerdictPass),
	}
	engine := fixture.engine(t, runner, &taskFakeVerifier{calls: fixture.calls}, &engineFakeCommitter{calls: fixture.calls}, fixture.worktree)
	plan := fixture.qaPlan()
	plan.TargetBranch = ""

	result, err := engine.TaskCycle(context.Background(), plan)

	if err != nil {
		t.Fatalf("task cycle: %v", err)
	}
	if result.QAVerdict != spec.VerdictPass {
		t.Fatalf("expected a Run with no recorded target branch to still reach a verdict, got %+v", result)
	}
	if len(runner.qaPrompts) != 1 {
		t.Fatalf("expected one QA prompt, got %d", len(runner.qaPrompts))
	}
	prompt := runner.qaPrompts[0]
	if strings.Contains(prompt, "Spec target branch:") {
		t.Fatalf("expected no Spec target branch line for a Run that recorded none, got:\n%s", prompt)
	}
	if !strings.Contains(prompt, "Spec: "+taskCycleSlug) || !strings.Contains(prompt, "Run the qa-gate process for this Spec") {
		t.Fatalf("expected the Spec identity and the QA contract to survive, got:\n%s", prompt)
	}
}

// The QA gate runs only when every Task in the Task Graph ended completed
// (ADR 0015). Each case leaves exactly one Task non-completed by a
// different route and asserts the gate never started: no QA prompt, no
// daemon.qa event, and no QA Report commit.
func TestTaskCycleQAStepRequiresEveryGraphTaskCompleted(t *testing.T) {
	tests := []struct {
		name          string
		seeds         []taskSpecSeed
		wantCompleted int
		wantFailed    int
		wantSkipped   int
	}{
		{
			name:       "failed Task",
			seeds:      []taskSpecSeed{{id: "task_01", verification: []string{"broken"}}},
			wantFailed: 1,
		},
		{
			name: "skipped Task behind an unmet dependency",
			seeds: []taskSpecSeed{
				{id: "task_01", verification: []string{"broken"}},
				{id: "task_02", needs: []string{"task_01"}},
			},
			wantFailed:  1,
			wantSkipped: 1,
		},
		{
			name: "earlier Run completed one Task while another chain does not finish",
			seeds: []taskSpecSeed{
				{id: "task_01", status: string(spec.StatusCompleted)},
				{id: "task_02", verification: []string{"broken"}},
				{id: "task_03", needs: []string{"task_02"}},
			},
			wantFailed:  1,
			wantSkipped: 1,
		},
		{
			name: "earlier Run completed every Task but one still has to run and fails",
			seeds: []taskSpecSeed{
				{id: "task_01", status: string(spec.StatusCompleted)},
				{id: "task_02", status: string(spec.StatusCompleted)},
				{id: "task_03", verification: []string{"broken"}},
			},
			wantFailed: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fixture := newTaskCycleFixture(t, tt.seeds)
			runner := &taskFakeRunner{
				calls:    fixture.calls,
				gitRoot:  fixture.gitRoot,
				qaReport: qaReportForTest(spec.VerdictPass),
			}
			committer := &engineFakeCommitter{calls: fixture.calls}
			verifier := &taskFakeVerifier{calls: fixture.calls, failOn: map[string]error{"broken": errors.New("gate broke")}}
			engine := fixture.engine(t, runner, verifier, committer, fixture.worktree)

			result, err := engine.TaskCycle(context.Background(), fixture.qaPlan())

			if err != nil {
				t.Fatalf("task cycle: %v", err)
			}
			if result.Completed != tt.wantCompleted || result.Failed != tt.wantFailed || result.Skipped != tt.wantSkipped {
				t.Fatalf("expected %d completed, %d failed, %d skipped, got %+v", tt.wantCompleted, tt.wantFailed, tt.wantSkipped, result)
			}
			assertNoQAStep(t, fixture, runner.qaPrompts, committer.messages, result)
		})
	}
}

// A Task the scheduler can never start — spec.Load rejects a need outside
// the Task Graph, so only a hand-built plan reaches here — settles skipped
// through the scheduler's defensive pass. The gate stays shut instead of
// crediting the Tasks that did run.
func TestTaskCycleQAStepSkippedWhenATaskNeverBecomesReady(t *testing.T) {
	fixture := newTaskCycleFixture(t, []taskSpecSeed{{id: "task_01", title: "Build the feature"}})
	runner := &taskFakeRunner{
		calls:        fixture.calls,
		gitRoot:      fixture.gitRoot,
		statusByTask: map[string]spec.Status{"task_01": spec.StatusCompleted},
		qaReport:     qaReportForTest(spec.VerdictPass),
	}
	committer := &engineFakeCommitter{calls: fixture.calls}
	engine := fixture.engine(t, runner, &taskFakeVerifier{calls: fixture.calls}, committer, fixture.worktree)
	plan := fixture.qaPlan()
	plan.Tasks = append(append([]spec.Task(nil), plan.Tasks...), spec.Task{ID: "task_99", Needs: []string{"task_absent"}})

	result, err := engine.TaskCycle(context.Background(), plan)

	if err != nil {
		t.Fatalf("task cycle: %v", err)
	}
	if result.Completed != 1 || result.Failed != 0 || result.Skipped != 1 {
		t.Fatalf("expected the ready Task completed and the unreachable Task skipped, got %+v", result)
	}
	assertNoQAStep(t, fixture, runner.qaPrompts, committer.messages, result)
}

// Every Task completing is still not enough: a Stop Request that lands
// mid-wave ends the Run before the gate, so the QA step never starts even
// though the Task Graph finished completed.
func TestTaskCycleStopRequestMidWaveSkipsQAWithEveryTaskCompleted(t *testing.T) {
	fixture := newTaskCycleFixture(t, []taskSpecSeed{
		{id: "task_01", title: "Start first"},
		{id: "task_02", title: "Start second"},
	})
	taskWorktrees := newFakeTaskWorktrees()
	taskWorktrees.onIntegrate = func(taskID string) {
		if taskID == "task_01" {
			_ = fixture.store.RequestStop(context.Background(), fixture.run.ID)
		}
	}
	runner := newTaskSchedulerRunner("task_01", "task_02")
	committer := &taskSchedulerCommitter{}
	engine := fixture.engineWithTaskWorktrees(t, runner, &taskSchedulerVerifier{}, committer, fixture.worktree, taskWorktrees)
	plan := fixture.qaPlan()
	plan.Concurrency = 2

	resultCh := make(chan struct {
		result TaskCycleResult
		err    error
	}, 1)
	go func() {
		result, err := engine.TaskCycle(context.Background(), plan)
		resultCh <- struct {
			result TaskCycleResult
			err    error
		}{result: result, err: err}
	}()

	assertTaskSet(t, waitSchedulerStarts(t, runner, 2), "task_01", "task_02")
	runner.releaseTask("task_01")
	if got := waitIntegratedTask(t, taskWorktrees); got != "task_01" {
		t.Fatalf("expected task_01 to integrate first, got %s", got)
	}
	runner.releaseTask("task_02")

	outcome := waitTaskCycleResult(t, resultCh)
	if !errors.Is(outcome.err, ErrStopRequested) {
		t.Fatalf("expected ErrStopRequested, got %v", outcome.err)
	}
	if outcome.result.Completed != 2 || outcome.result.Failed != 0 || outcome.result.Skipped != 0 {
		t.Fatalf("expected every Task completed before the stop, got %+v", outcome.result)
	}
	// Only a completed Task integrates, so both integrations confirm the
	// Task Graph itself finished completed.
	assertTaskSet(t, taskWorktrees.integratedTasks(), "task_01", "task_02")
	assertNoQAStep(t, fixture, runner.qaPrompts(), committer.commitMessages(), outcome.result)
}

// assertNoQAStep asserts the observable evidence a QA step leaves behind is
// absent: the Agent got no QA prompt, no daemon.qa Run Event exists, no QA
// Report commit was created, and the cycle settled no verdict.
func assertNoQAStep(t *testing.T, fixture *taskCycleFixture, qaPrompts []string, commitMessages []string, result TaskCycleResult) {
	t.Helper()
	if len(qaPrompts) != 0 {
		t.Fatalf("expected no QA prompt, got %v", qaPrompts)
	}
	if events := taskEventsOfKind(fixture.sink, runevent.KindDaemonQA); len(events) != 0 {
		t.Fatalf("expected no daemon.qa event, got %+v", events)
	}
	for _, message := range commitMessages {
		if strings.HasPrefix(message, "docs: qa report for ") {
			t.Fatalf("expected no QA Report commit, got %v", commitMessages)
		}
	}
	if result.QAVerdict != "" || result.QAReportPath != "" {
		t.Fatalf("expected no settled QA verdict, got %+v", result)
	}
	reportDir := filepath.Join(fixture.specsRoot, taskCycleSlug, "qa")
	if _, err := os.Stat(reportDir); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("expected no QA Report directory, got err=%v", err)
	}
}

// allTasksRunCompleted is the whole guarantee behind the QA step, so it
// answers for every Task in the Task Graph — never for the subset this Run
// happened to execute.
func TestAllTasksRunCompletedCoversEveryGraphTask(t *testing.T) {
	tasks := []spec.Task{{ID: "task_01"}, {ID: "task_02"}, {ID: "task_03"}}
	tests := []struct {
		name     string
		statuses map[string]taskRunStatus
		want     bool
	}{
		{
			name:     "every graph Task completed",
			statuses: map[string]taskRunStatus{"task_01": taskRunCompleted, "task_02": taskRunCompleted, "task_03": taskRunCompleted},
			want:     true,
		},
		{
			name:     "one Task still pending",
			statuses: map[string]taskRunStatus{"task_01": taskRunCompleted, "task_02": taskRunPending, "task_03": taskRunCompleted},
		},
		{
			name:     "one Task still running",
			statuses: map[string]taskRunStatus{"task_01": taskRunCompleted, "task_02": taskRunRunning, "task_03": taskRunCompleted},
		},
		{
			name:     "one Task failed",
			statuses: map[string]taskRunStatus{"task_01": taskRunCompleted, "task_02": taskRunFailed, "task_03": taskRunCompleted},
		},
		{
			name:     "one Task skipped",
			statuses: map[string]taskRunStatus{"task_01": taskRunCompleted, "task_02": taskRunSkipped, "task_03": taskRunCompleted},
		},
		{
			name:     "one graph Task missing from the status map",
			statuses: map[string]taskRunStatus{"task_01": taskRunCompleted, "task_02": taskRunCompleted},
		},
		{
			name:     "a completed status outside the graph never stands in for a graph Task",
			statuses: map[string]taskRunStatus{"task_01": taskRunCompleted, "task_02": taskRunCompleted, "task_04": taskRunCompleted},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := allTasksRunCompleted(tasks, tt.statuses); got != tt.want {
				t.Fatalf("allTasksRunCompleted = %v, want %v", got, tt.want)
			}
		})
	}
}

// initialTaskRunStatuses seeds Tasks an earlier Run completed, so the QA
// gate can credit them, and re-runs every other status.
func TestInitialTaskRunStatusesSeedsEarlierRunCompletions(t *testing.T) {
	statuses := initialTaskRunStatuses([]spec.Task{
		{ID: "task_01", Status: spec.StatusCompleted},
		{ID: "task_02", Status: spec.StatusPending},
		{ID: "task_03", Status: spec.StatusInProgress},
		{ID: "task_04", Status: spec.StatusFailed},
	})
	want := map[string]taskRunStatus{
		"task_01": taskRunCompleted,
		"task_02": taskRunPending,
		"task_03": taskRunPending,
		"task_04": taskRunPending,
	}
	if len(statuses) != len(want) {
		t.Fatalf("expected one status per graph Task, got %v", statuses)
	}
	for id, expected := range want {
		if statuses[id] != expected {
			t.Fatalf("expected %s seeded %q, got %q", id, expected, statuses[id])
		}
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
