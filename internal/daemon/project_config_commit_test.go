// Suite: Project Config commit policy
// Invariant: Project Config enters only an authorized Task commit and every new exclusion is observable.
// Boundary IN: Daemon Task, Batch, and QA Report settlement and commit preparation.
// Boundary OUT: Git pathspec exclusion inside GitCommitter, covered by daemon_test.go.
package daemon

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"roundfix/internal/runevent"
	"roundfix/internal/spec"
)

func TestTaskCommitCarriesProjectConfigBoundedByTheAuthorization(t *testing.T) {
	t.Parallel()
	fixture := newTaskCycleFixture(t, []taskSpecSeed{{id: "task_01", title: "Update Project Config"}})
	setProjectConfigAuthorizationForTest(t, fixture)
	fixture.worktree.snapshots = [][]string{nil, {projectConfigPath}}
	runner := &taskFakeRunner{
		calls:       fixture.calls,
		gitRoot:     fixture.gitRoot,
		writeByTask: map[string]string{"task_01": projectConfigPath},
	}
	committer := &engineFakeCommitter{calls: fixture.calls}
	engine := fixture.engine(t, runner, &taskFakeVerifier{calls: fixture.calls}, committer, fixture.worktree)

	result, err := engine.TaskCycle(context.Background(), fixture.plan())

	if err != nil {
		t.Fatalf("TaskCycle: %v", err)
	}
	if result.Completed != 1 || result.Failed != 0 || result.Skipped != 0 {
		t.Fatalf("TaskCycle result = %+v, want one completed Task", result)
	}
	if got := taskStatusOnDisk(t, fixture.gitRoot, "task_01"); got != string(spec.StatusCompleted) {
		t.Fatalf("Task status = %q, want %q", got, spec.StatusCompleted)
	}
	assertCommittedPaths(t, committer, projectConfigPath, taskFileRel(taskCycleSlug, "task_01"))
	if dropped := droppedStageEvents(t, fixture.sink); len(dropped) != 0 {
		t.Fatalf("authorized Project Config produced dropped-path events: %+v", dropped)
	}
}

func TestTaskOutsideTheAuthorizationFailsOnProjectConfig(t *testing.T) {
	t.Parallel()
	fixture := newTaskCycleFixture(t, []taskSpecSeed{{id: "task_01", title: "Update Project Config"}})
	fixture.worktree.snapshots = [][]string{nil, {projectConfigPath}}
	runner := &taskFakeRunner{
		calls:       fixture.calls,
		gitRoot:     fixture.gitRoot,
		writeByTask: map[string]string{"task_01": projectConfigPath},
	}
	committer := &engineFakeCommitter{calls: fixture.calls}
	engine := fixture.engine(t, runner, &taskFakeVerifier{calls: fixture.calls}, committer, fixture.worktree)

	result, err := engine.TaskCycle(context.Background(), fixture.plan())

	if err != nil {
		t.Fatalf("TaskCycle: %v", err)
	}
	wantReason := "Task commit lost output: " + projectConfigPath + " (" + projectConfigOutsideAuthorizationReason + ")"
	if result.Completed != 0 || result.Failed != 1 || result.Skipped != 0 {
		t.Fatalf("TaskCycle result = %+v, want one failed Task", result)
	}
	outcome, found := taskOutcomeByID(result.Outcomes, "task_01")
	if !found || outcome.Status != string(spec.StatusFailed) || outcome.Reason != wantReason {
		t.Fatalf("Task outcome = %+v, found=%t, want failed with %q", outcome, found, wantReason)
	}
	if got := taskStatusOnDisk(t, fixture.gitRoot, "task_01"); got != string(spec.StatusFailed) {
		t.Fatalf("Task status = %q, want %q", got, spec.StatusFailed)
	}
	if len(committer.paths) != 0 {
		t.Fatalf("outside-authority Task created commits: %v", committer.paths)
	}
	assertProjectConfigDropped(t, fixture.sink, fixture.progress.String(), "task_01", projectConfigOutsideAuthorizationReason)
	assertTaskSettlementReason(t, fixture.sink, "task_01", wantReason)
}

func TestBatchCommitReportsProjectConfigAsDropped(t *testing.T) {
	t.Parallel()
	fixture := newEngineFixture(t)
	regularPath := "src/fixed.go"
	writeProjectConfigCommitFixtureFiles(t, fixture.gitRoot, projectConfigPath, regularPath)
	fixture.worktree.snapshots = [][]string{{projectConfigPath, regularPath}}
	committer := &engineFakeCommitter{calls: fixture.calls}
	engine := fixture.engine(
		t,
		&engineFakeRunner{calls: fixture.calls, store: fixture.store},
		&engineFakeVerifier{calls: fixture.calls, store: fixture.store, runID: fixture.run.ID},
		committer,
		&engineFakePusher{calls: fixture.calls},
		&engineFakeSource{calls: fixture.calls},
	)
	plan := fixture.plan()

	committed, skipped, err := engine.commitBatch(context.Background(), plan, plan.Batches[0], nil)

	if err != nil {
		t.Fatalf("commitBatch: %v", err)
	}
	if !committed || skipped {
		t.Fatalf("commitBatch outcome = committed %t skipped %t, want created", committed, skipped)
	}
	assertCommittedPaths(t, committer, regularPath)
	assertProjectConfigDropped(t, fixture.sink, fixture.progress.String(), "", projectConfigNonTaskCommitExclusionReason)
}

func TestQAReportCommitReportsProjectConfigAsDropped(t *testing.T) {
	t.Parallel()
	fixture := newTaskCycleFixture(t, []taskSpecSeed{{id: "task_01", status: string(spec.StatusCompleted)}})
	regularPath := "qa/evidence/observed.txt"
	reportPath := qaReportRelPathForTest()
	writeProjectConfigCommitFixtureFiles(t, fixture.gitRoot, projectConfigPath, regularPath, reportPath)
	fixture.worktree.snapshots = [][]string{{projectConfigPath, regularPath, reportPath}}
	committer := &engineFakeCommitter{calls: fixture.calls}
	engine := fixture.engine(t, &taskFakeRunner{calls: fixture.calls}, &taskFakeVerifier{calls: fixture.calls}, committer, fixture.worktree)
	plan := fixture.plan()
	plan.Authorization = spec.AuthorizationResolution{}

	err := engine.commitQAReport(context.Background(), plan, 2, nil, nil, spec.VerdictPass, reportPath, spec.Task{})

	if err != nil {
		t.Fatalf("commitQAReport: %v", err)
	}
	assertCommittedPaths(t, committer, reportPath, regularPath)
	assertProjectConfigDropped(t, fixture.sink, fixture.progress.String(), "", projectConfigNonTaskCommitExclusionReason)
}

func TestPreexistingProjectConfigChangeStaysOutOfTheTaskCommit(t *testing.T) {
	t.Parallel()

	t.Run("Task commit", func(t *testing.T) {
		fixture := newTaskCycleFixture(t, []taskSpecSeed{{id: "task_01", title: "Change ordinary source"}})
		regularPath := "internal/ordinary.go"
		writeProjectConfigCommitFixtureFiles(t, fixture.gitRoot, projectConfigPath)
		fixture.worktree.snapshots = [][]string{{projectConfigPath}, {projectConfigPath, regularPath}}
		runner := &taskFakeRunner{
			calls:       fixture.calls,
			gitRoot:     fixture.gitRoot,
			writeByTask: map[string]string{"task_01": regularPath},
		}
		committer := &engineFakeCommitter{calls: fixture.calls}
		engine := fixture.engine(t, runner, &taskFakeVerifier{calls: fixture.calls}, committer, fixture.worktree)

		result, err := engine.TaskCycle(context.Background(), fixture.plan())

		if err != nil {
			t.Fatalf("TaskCycle: %v", err)
		}
		if result.Completed != 1 || result.Failed != 0 || result.Skipped != 0 {
			t.Fatalf("TaskCycle result = %+v, want one completed Task", result)
		}
		assertCommittedPaths(t, committer, taskFileRel(taskCycleSlug, "task_01"), regularPath)
		assertNoProjectConfigDrop(t, fixture.sink, fixture.progress.String())
	})

	t.Run("Batch commit", func(t *testing.T) {
		fixture := newEngineFixture(t)
		regularPath := "src/fixed.go"
		writeProjectConfigCommitFixtureFiles(t, fixture.gitRoot, projectConfigPath, regularPath)
		fixture.worktree.snapshots = [][]string{{projectConfigPath, regularPath}}
		committer := &engineFakeCommitter{calls: fixture.calls}
		engine := fixture.engine(
			t,
			&engineFakeRunner{calls: fixture.calls, store: fixture.store},
			&engineFakeVerifier{calls: fixture.calls, store: fixture.store, runID: fixture.run.ID},
			committer,
			&engineFakePusher{calls: fixture.calls},
			&engineFakeSource{calls: fixture.calls},
		)
		plan := fixture.plan()

		committed, skipped, err := engine.commitBatch(context.Background(), plan, plan.Batches[0], []string{projectConfigPath})

		if err != nil {
			t.Fatalf("commitBatch: %v", err)
		}
		if !committed || skipped {
			t.Fatalf("commitBatch outcome = committed %t skipped %t, want created", committed, skipped)
		}
		assertCommittedPaths(t, committer, regularPath)
		assertNoProjectConfigDrop(t, fixture.sink, fixture.progress.String())
	})

	t.Run("QA Report commit", func(t *testing.T) {
		fixture := newTaskCycleFixture(t, []taskSpecSeed{{id: "task_01", status: string(spec.StatusCompleted)}})
		regularPath := "qa/evidence/observed.txt"
		reportPath := qaReportRelPathForTest()
		writeProjectConfigCommitFixtureFiles(t, fixture.gitRoot, projectConfigPath, regularPath, reportPath)
		fixture.worktree.snapshots = [][]string{{projectConfigPath, regularPath, reportPath}}
		committer := &engineFakeCommitter{calls: fixture.calls}
		engine := fixture.engine(t, &taskFakeRunner{calls: fixture.calls}, &taskFakeVerifier{calls: fixture.calls}, committer, fixture.worktree)
		plan := fixture.plan()
		plan.Authorization = spec.AuthorizationResolution{}

		err := engine.commitQAReport(context.Background(), plan, 2, []string{projectConfigPath}, nil, spec.VerdictPass, reportPath, spec.Task{})

		if err != nil {
			t.Fatalf("commitQAReport: %v", err)
		}
		assertCommittedPaths(t, committer, reportPath, regularPath)
		assertNoProjectConfigDrop(t, fixture.sink, fixture.progress.String())
	})
}

func setProjectConfigAuthorizationForTest(t *testing.T, fixture *taskCycleFixture) {
	t.Helper()
	record := fmt.Sprintf(`---
status: approved
granted: 2026-09-28
action: update Project Config in the fixture Task
consuming: %s
paths:
  - %s
operations:
  - implement
  - commit
  - push
---

# Approved fixture Project Config authority
`, taskCycleSlug, projectConfigPath)
	mustWriteForTest(t, filepath.Join(fixture.gitRoot, "docs", "specs", taskCycleSlug, "_authorization.md"), record)
	commitTaskFixtureSource(t, fixture.gitRoot, "authorize fixture Project Config")
}

func writeProjectConfigCommitFixtureFiles(t *testing.T, root string, paths ...string) {
	t.Helper()
	for _, path := range paths {
		fullPath := filepath.Join(root, path)
		if err := os.MkdirAll(filepath.Dir(fullPath), 0o755); err != nil {
			t.Fatalf("create fixture directory for %s: %v", path, err)
		}
		mustWriteForTest(t, fullPath, path+"\n")
	}
}

func assertCommittedPaths(t *testing.T, committer *engineFakeCommitter, want ...string) {
	t.Helper()
	if len(committer.paths) != 1 {
		t.Fatalf("commit count = %d, want one; paths=%v", len(committer.paths), committer.paths)
	}
	if !reflect.DeepEqual(committer.paths[0], want) {
		t.Fatalf("committed paths = %v, want %v", committer.paths[0], want)
	}
}

func assertProjectConfigDropped(t *testing.T, sink *captureEventSink, progress string, taskID string, reason string) {
	t.Helper()
	events := droppedStageEvents(t, sink)
	if len(events) != 1 {
		t.Fatalf("dropped-path events = %+v, want one", events)
	}
	wantPayload := map[string]any{
		"decision": "dropped",
		"path":     projectConfigPath,
		"reason":   reason,
	}
	if taskID != "" {
		wantPayload["task"] = taskID
	}
	if got := eventPayloadMap(t, events[0]); !reflect.DeepEqual(got, wantPayload) {
		t.Fatalf("dropped-path payload = %#v, want %#v", got, wantPayload)
	}
	wantSummary := "Project Config " + projectConfigPath + " omitted from the commit: " + reason + "."
	if events[0].Summary != wantSummary {
		t.Fatalf("dropped-path summary = %q, want %q", events[0].Summary, wantSummary)
	}
	wantProgress := "roundfix: Project Config " + projectConfigPath + " omitted from the commit: " + reason + "\n"
	if count := strings.Count(progress, wantProgress); count != 1 {
		t.Fatalf("Project Config progress line count = %d, want one; progress=%q", count, progress)
	}
}

func assertNoProjectConfigDrop(t *testing.T, sink *captureEventSink, progress string) {
	t.Helper()
	for _, event := range droppedStageEvents(t, sink) {
		if eventPayloadString(t, event, "path") == projectConfigPath {
			t.Fatalf("pre-existing Project Config produced a dropped-path event: %+v", event)
		}
	}
	if strings.Contains(progress, "Project Config "+projectConfigPath+" omitted from the commit") {
		t.Fatalf("pre-existing Project Config produced an exclusion line: %q", progress)
	}
}

func assertTaskSettlementReason(t *testing.T, sink *captureEventSink, taskID string, reason string) {
	t.Helper()
	want := map[string]any{
		"phase":  "settled",
		"reason": reason,
		"status": string(spec.StatusFailed),
		"task":   taskID,
	}
	for _, event := range taskEventsOfKind(sink, runevent.KindDaemonTask) {
		payload := eventPayloadMap(t, event)
		if event.ReviewIssue == taskID && payload["phase"] == "settled" {
			if !reflect.DeepEqual(payload, want) {
				t.Fatalf("Task settlement payload = %#v, want %#v", payload, want)
			}
			return
		}
	}
	t.Fatalf("no settlement event found for %s", taskID)
}
