package daemon

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"roundfix/internal/runevent"
	"roundfix/internal/spec"
)

func TestTaskCommitRecordsAnUndeclaredTestFile(t *testing.T) {
	t.Parallel()
	fixture := newTaskCycleFixture(t, []taskSpecSeed{{id: "task_01", title: "Update delivery storage"}})
	setRecordedPathsTaskContext(t, fixture,
		"- interface: `internal/store/delivery.go`",
	)
	runner := &taskFakeRunner{
		calls:       fixture.calls,
		gitRoot:     fixture.gitRoot,
		writeByTask: map[string]string{"task_01": "internal/store/delivery_test.go"},
	}
	engine := fixture.engine(t, runner, &taskFakeVerifier{calls: fixture.calls}, GitCommitter{}, GitWorktreeSnapshotter{})

	result, err := engine.TaskCycle(context.Background(), fixture.plan())

	if err != nil {
		t.Fatalf("TaskCycle: %v", err)
	}
	if result.Completed != 1 || result.Failed != 0 {
		t.Fatalf("TaskCycle result = %+v, want one completed Task", result)
	}
	want := []string{"internal/store/delivery_test.go"}
	assertRecordedTaskPaths(t, fixture, "task_01", want)
	assertRecordedCommitEvent(t, fixture, "task_01", want)
}

func TestTaskCommitRecordsTheFilesOfANewUntrackedPackage(t *testing.T) {
	t.Parallel()
	fixture := newTaskCycleFixture(t, []taskSpecSeed{{id: "task_01", title: "Add a package"}})
	var writeErr error
	runner := &taskFakeRunner{
		calls:       fixture.calls,
		gitRoot:     fixture.gitRoot,
		writeByTask: map[string]string{"task_01": "internal/newpkg/new.go"},
		afterTask: func(taskID string) {
			if taskID != "task_01" {
				return
			}
			writeErr = os.WriteFile(filepath.Join(fixture.gitRoot, "internal", "newpkg", "new_test.go"), []byte("package newpkg\n"), 0o644)
		},
	}
	engine := fixture.engine(t, runner, &taskFakeVerifier{calls: fixture.calls}, GitCommitter{}, GitWorktreeSnapshotter{})

	result, err := engine.TaskCycle(context.Background(), fixture.plan())

	if err != nil {
		t.Fatalf("TaskCycle: %v", err)
	}
	if writeErr != nil {
		t.Fatalf("write second package file: %v", writeErr)
	}
	if result.Completed != 1 || result.Failed != 0 {
		t.Fatalf("TaskCycle result = %+v, want one completed Task", result)
	}
	want := []string{"internal/newpkg/new.go", "internal/newpkg/new_test.go"}
	assertRecordedTaskPaths(t, fixture, "task_01", want)
	assertRecordedCommitEvent(t, fixture, "task_01", want)
	if slices.Contains(spec.RecordedTaskPaths(recordedTaskContent(t, fixture, "task_01")), "internal/newpkg") {
		t.Fatal("recorded paths contains the package directory instead of only its files")
	}
}

func TestTaskCommitRecordsNothingWhenEveryPathIsDeclared(t *testing.T) {
	t.Parallel()
	fixture := newTaskCycleFixture(t, []taskSpecSeed{{id: "task_01", title: "Update delivery storage"}})
	setRecordedPathsTaskContext(t, fixture,
		"- interface: `internal/store/delivery.go`",
	)
	taskPath := taskPathFor(fixture.gitRoot, taskCycleSlug, "task_01")
	if err := spec.RecordTaskPaths(taskPath, []string{"old/attempt.go"}); err != nil {
		t.Fatalf("seed old recorded paths: %v", err)
	}
	commitTaskFixtureSource(t, fixture.gitRoot, "seed an earlier recorded-path attempt")
	fixture.reloadGraph()
	runner := &taskFakeRunner{
		calls:       fixture.calls,
		gitRoot:     fixture.gitRoot,
		writeByTask: map[string]string{"task_01": "internal/store/delivery.go"},
	}
	engine := fixture.engine(t, runner, &taskFakeVerifier{calls: fixture.calls}, GitCommitter{}, GitWorktreeSnapshotter{})

	result, err := engine.TaskCycle(context.Background(), fixture.plan())

	if err != nil {
		t.Fatalf("TaskCycle: %v", err)
	}
	if result.Completed != 1 || result.Failed != 0 {
		t.Fatalf("TaskCycle result = %+v, want one completed Task", result)
	}
	assertRecordedTaskPaths(t, fixture, "task_01", nil)
	assertNoRecordedPathsInCommitEvents(t, fixture)
}

func TestTaskCommitNeverRecordsAGovernedPath(t *testing.T) {
	t.Parallel()
	fixture := newTaskCycleFixture(t, []taskSpecSeed{{id: "task_01", title: "Update governed guidance"}})
	authorizationPath := filepath.Join(fixture.specsRoot, taskCycleSlug, "_authorization.md")
	authorization, err := os.ReadFile(authorizationPath)
	if err != nil {
		t.Fatalf("read fixture authorization: %v", err)
	}
	authorization = []byte(strings.Replace(string(authorization), "paths:\n", "paths:\n  - docs/agents/agent-instructions.md\n", 1))
	if err := os.WriteFile(authorizationPath, authorization, 0o644); err != nil {
		t.Fatalf("authorize governed fixture path: %v", err)
	}
	commitTaskFixtureSource(t, fixture.gitRoot, "authorize governed fixture path")
	runner := &taskFakeRunner{
		calls:       fixture.calls,
		gitRoot:     fixture.gitRoot,
		writeByTask: map[string]string{"task_01": "docs/agents/agent-instructions.md"},
	}
	engine := fixture.engine(t, runner, &taskFakeVerifier{calls: fixture.calls}, GitCommitter{}, GitWorktreeSnapshotter{})

	result, err := engine.TaskCycle(context.Background(), fixture.plan())

	if err != nil {
		t.Fatalf("TaskCycle: %v", err)
	}
	if result.Completed != 1 || result.Failed != 0 {
		t.Fatalf("TaskCycle result = %+v, want one completed Task", result)
	}
	assertRecordedTaskPaths(t, fixture, "task_01", nil)
	assertNoRecordedPathsInCommitEvents(t, fixture)
}

func TestAFailedTaskRecordsNoPaths(t *testing.T) {
	t.Parallel()
	fixture := newTaskCycleFixture(t, []taskSpecSeed{{id: "task_01", title: "Fail after changing a file"}})
	runner := &taskFakeRunner{
		calls:       fixture.calls,
		gitRoot:     fixture.gitRoot,
		writeByTask: map[string]string{"task_01": "internal/failed.go"},
	}
	verifier := &taskFakeVerifier{
		calls:  fixture.calls,
		script: []error{errors.New("first failure"), errors.New("final failure")},
	}
	engine := fixture.engine(t, runner, verifier, GitCommitter{}, GitWorktreeSnapshotter{})

	result, err := engine.TaskCycle(context.Background(), fixture.plan())

	if err != nil {
		t.Fatalf("TaskCycle: %v", err)
	}
	if result.Completed != 0 || result.Failed != 1 {
		t.Fatalf("TaskCycle result = %+v, want one failed Task", result)
	}
	assertRecordedTaskPaths(t, fixture, "task_01", nil)
	assertNoRecordedPathsInCommitEvents(t, fixture)
}

func TestAQATaskRecordsNoPaths(t *testing.T) {
	t.Parallel()
	fixture := newTaskCycleFixture(t, []taskSpecSeed{{id: "task_01", status: string(spec.StatusCompleted)}})
	runner := &taskFakeRunner{
		calls:    fixture.calls,
		gitRoot:  fixture.gitRoot,
		qaReport: "---\nverdict: pass\n---\n\n# QA Report\n\n## Results\n\n| # | Status | Evidence |\n| - | --- | --- |\n| R01 | pass | observed fixture |\n",
	}
	engine := fixture.engine(t, runner, &taskFakeVerifier{calls: fixture.calls}, GitCommitter{}, GitWorktreeSnapshotter{})
	engine.deps.MechanicalStage = &fakeQAMechanicalStage{}
	plan := fixture.qaPlan()
	qaTaskID := fixture.graph.QATaskID

	result, err := engine.TaskCycle(context.Background(), plan)

	if err != nil {
		t.Fatalf("TaskCycle: %v", err)
	}
	if result.QAVerdict != spec.VerdictPass {
		t.Fatalf("QA verdict = %q, want pass", result.QAVerdict)
	}
	assertRecordedTaskPaths(t, fixture, qaTaskID, nil)
	assertNoRecordedPathsInCommitEvents(t, fixture)
}

func setRecordedPathsTaskContext(t *testing.T, fixture *taskCycleFixture, lines ...string) {
	t.Helper()
	path := taskPathFor(fixture.gitRoot, taskCycleSlug, "task_01")
	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read Task fixture: %v", err)
	}
	contextSection := "## Context\n\n" + strings.Join(lines, "\n") + "\n\n"
	updated := strings.Replace(string(content), "## Verification\n", contextSection+"## Verification\n", 1)
	if updated == string(content) {
		t.Fatal("Task fixture has no Verification heading")
	}
	if err := os.WriteFile(path, []byte(updated), 0o644); err != nil {
		t.Fatalf("write Task Context: %v", err)
	}
	commitTaskFixtureSource(t, fixture.gitRoot, "declare Task Context")
	fixture.reloadGraph()
}

func recordedTaskContent(t *testing.T, fixture *taskCycleFixture, taskID string) []byte {
	t.Helper()
	content, err := os.ReadFile(taskPathFor(fixture.gitRoot, taskCycleSlug, taskID))
	if err != nil {
		t.Fatalf("read Task %s: %v", taskID, err)
	}
	return content
}

func assertRecordedTaskPaths(t *testing.T, fixture *taskCycleFixture, taskID string, want []string) {
	t.Helper()
	if got := spec.RecordedTaskPaths(recordedTaskContent(t, fixture, taskID)); !slices.Equal(got, want) {
		t.Fatalf("Task %s recorded paths = %v, want %v", taskID, got, want)
	}
}

func assertRecordedCommitEvent(t *testing.T, fixture *taskCycleFixture, taskID string, want []string) {
	t.Helper()
	for _, event := range taskEventsOfKind(fixture.sink, runevent.KindDaemonCommit) {
		payload := eventPayloadMap(t, event)
		if event.ReviewIssue != taskID || payload["decision"] != "created" {
			continue
		}
		raw, ok := payload["recorded_paths"].([]any)
		if !ok {
			t.Fatalf("Task %s commit event payload = %#v, want recorded_paths", taskID, payload)
		}
		got := make([]string, 0, len(raw))
		for _, value := range raw {
			path, ok := value.(string)
			if !ok {
				t.Fatalf("recorded_paths value = %#v, want string", value)
			}
			got = append(got, path)
		}
		if !slices.Equal(got, want) {
			t.Fatalf("Task %s commit event recorded_paths = %v, want %v", taskID, got, want)
		}
		return
	}
	t.Fatalf("Task %s has no created commit event", taskID)
}

func assertNoRecordedPathsInCommitEvents(t *testing.T, fixture *taskCycleFixture) {
	t.Helper()
	for _, event := range taskEventsOfKind(fixture.sink, runevent.KindDaemonCommit) {
		payload := eventPayloadMap(t, event)
		if _, present := payload["recorded_paths"]; present {
			t.Fatalf("commit event unexpectedly carries recorded_paths: %#v", payload)
		}
	}
}
