// Suite: QA settlement at the successful Agent turn boundary.
// Boundary OUT: Agent processes and Verification are scripted package fakes.
package daemon

import (
	"context"
	"encoding/json"
	"path/filepath"
	"reflect"
	"testing"

	"roundfix/internal/agent"
	"roundfix/internal/runevent"
	"roundfix/internal/spec"
)

type unfinishedTurnRunner struct {
	*taskFakeRunner
	backgroundTasks []agent.AgentBackgroundTask
}

func (runner *unfinishedTurnRunner) Run(ctx context.Context, req agent.ExecuteRequest, sink runevent.Sink) (agent.ExecuteResult, error) {
	result, err := runner.taskFakeRunner.Run(ctx, req, sink)
	result.BackgroundTasks = runner.backgroundTasks
	return result, err
}

func qaBackgroundTasksForTest() []agent.AgentBackgroundTask {
	return []agent.AgentBackgroundTask{
		{ID: "bvlc7cgwx", Tool: "Bash", Description: "Run the load driver"},
		{ID: "bekzicxjj", Tool: "Monitor", Description: "Watch the driver events"},
	}
}

func runQABackgroundTurnForTest(t *testing.T, tasks []agent.AgentBackgroundTask, report string) (*taskCycleFixture, TaskCycleResult) {
	t.Helper()
	fixture := newTaskCycleFixture(t, []taskSpecSeed{{id: "task_01", status: string(spec.StatusCompleted)}})
	runner := &unfinishedTurnRunner{taskFakeRunner: &taskFakeRunner{
		calls: fixture.calls, gitRoot: fixture.gitRoot, qaReport: report, preserveQASeed: report == "",
	}, backgroundTasks: tasks}
	engine := fixture.engine(t, runner, &taskFakeVerifier{calls: fixture.calls}, &engineFakeCommitter{calls: fixture.calls}, fixture.worktree)
	engine.deps.MechanicalStage = &fakeQAMechanicalStage{}
	result, err := engine.TaskCycle(t.Context(), fixture.qaPlan())
	if err != nil {
		t.Fatalf("QA turn halted the Run: %v", err)
	}
	if len(runner.qaPrompts) != 1 {
		t.Fatalf("QA turns=%d want=1", len(runner.qaPrompts))
	}
	if result.QAReportPath != runner.qaReportPath {
		t.Fatalf("report=%s want seed=%s", result.QAReportPath, runner.qaReportPath)
	}
	return fixture, result
}

func assertQABackgroundSettlement(t *testing.T, fixture *taskCycleFixture, status spec.Status, reason string) {
	t.Helper()
	if got := taskStatusOnDisk(t, fixture.gitRoot, fixture.graph.QATaskID); got != string(status) {
		t.Fatalf("QA status=%s want=%s", got, status)
	}
	count := 0
	for _, event := range eventsOfKind(fixture.sink, runevent.KindDaemonTask) {
		payload := eventPayloadMap(t, event)
		if payload["phase"] == "settled" && payload["task"] == fixture.graph.QATaskID {
			count++
			gotReason, _ := payload["reason"].(string)
			if payload["status"] != string(status) || gotReason != reason {
				t.Fatalf("settlement=%v want status=%s reason=%s", payload, status, reason)
			}
		}
	}
	if count != 1 {
		t.Fatalf("QA settlements=%d want=1", count)
	}
}

func unfinishedTurnEventsForTest(t *testing.T, fixture *taskCycleFixture) []runevent.RunEvent {
	t.Helper()
	var events []runevent.RunEvent
	for _, event := range eventsOfKind(fixture.sink, runevent.KindDaemonTask) {
		if eventPayloadMap(t, event)["phase"] == "turn_unfinished" {
			events = append(events, event)
		}
	}
	return events
}

func TestQATurnThatLeavesABackgroundTaskIsUnfinished(t *testing.T) {
	t.Parallel()
	tasks := qaBackgroundTasksForTest()
	fixture, result := runQABackgroundTurnForTest(t, tasks, "")
	if result.QAAccepted || result.QAVerdict != spec.VerdictPending {
		t.Fatalf("result=%+v", result)
	}
	const reason = "unfinished turn: the QA turn ended with its report pending after starting background task bvlc7cgwx (Bash: Run the load driver); bekzicxjj (Monitor: Watch the driver events)"
	assertQABackgroundSettlement(t, fixture, spec.StatusFailed, reason)
	events := unfinishedTurnEventsForTest(t, fixture)
	if len(events) != 1 {
		t.Fatalf("unfinished events=%d want=1", len(events))
	}
	var payload struct {
		Phase           string `json:"phase"`
		ScopeKind       string `json:"scope_kind"`
		ScopeID         string `json:"scope_id"`
		BackgroundTasks []struct {
			ID          string `json:"id"`
			Tool        string `json:"tool"`
			Description string `json:"description"`
		} `json:"background_tasks"`
	}
	if err := json.Unmarshal(events[0].Payload, &payload); err != nil {
		t.Fatal(err)
	}
	if payload.Phase != "turn_unfinished" || payload.ScopeKind != "qa" || payload.ScopeID != fixture.graph.QATaskID || events[0].RunID != fixture.run.ID {
		t.Fatalf("event=%+v payload=%+v", events[0], payload)
	}
	var got []agent.AgentBackgroundTask
	for _, task := range payload.BackgroundTasks {
		got = append(got, agent.AgentBackgroundTask{ID: task.ID, Tool: task.Tool, Description: task.Description})
	}
	if !reflect.DeepEqual(got, tasks) {
		t.Fatalf("event tasks=%+v want=%+v", got, tasks)
	}
	report, err := spec.ReadQAReportFile(filepath.Join(fixture.gitRoot, result.QAReportPath))
	if err != nil || report.Verdict != spec.VerdictPending {
		t.Fatalf("report=%+v err=%v", report, err)
	}
}

func TestQATurnWithBackgroundTaskAndFinishedReportSettlesNormally(t *testing.T) {
	t.Parallel()
	for _, verdict := range []string{spec.VerdictPass, spec.VerdictFail} {
		t.Run(verdict, func(t *testing.T) {
			t.Parallel()
			fixture, result := runQABackgroundTurnForTest(t, qaBackgroundTasksForTest(), qaReportForTest(verdict))
			if result.QAVerdict != verdict || result.QAAccepted != (verdict == spec.VerdictPass) {
				t.Fatalf("result=%+v", result)
			}
			status, reason := spec.StatusFailed, "QA verdict "+verdict
			if verdict == spec.VerdictPass {
				status, reason = spec.StatusCompleted, ""
			}
			assertQABackgroundSettlement(t, fixture, status, reason)
			if events := unfinishedTurnEventsForTest(t, fixture); len(events) != 0 {
				t.Fatalf("unexpected unfinished events=%v", events)
			}
		})
	}
}

func TestQATurnWithoutBackgroundTaskKeepsAPendingReport(t *testing.T) {
	t.Parallel()
	fixture, result := runQABackgroundTurnForTest(t, nil, "")
	if result.QAVerdict != spec.VerdictPending || result.QAAccepted {
		t.Fatalf("result=%+v", result)
	}
	assertQABackgroundSettlement(t, fixture, spec.StatusFailed, "QA verdict pending")
	if events := unfinishedTurnEventsForTest(t, fixture); len(events) != 0 {
		t.Fatalf("unexpected unfinished events=%v", events)
	}
	report, err := spec.ReadQAReportFile(filepath.Join(fixture.gitRoot, result.QAReportPath))
	if err != nil || report.Verdict != spec.VerdictPending {
		t.Fatalf("report=%+v err=%v", report, err)
	}
}
