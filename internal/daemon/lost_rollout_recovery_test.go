// Suite: Lost Rollout recovery at the Task and QA session boundary.
// Boundary OUT: runtime processes and Verification are package fakes.
package daemon

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"strings"
	"testing"

	"roundfix/internal/agent"
	roundconfig "roundfix/internal/config"
	"roundfix/internal/runevent"
	"roundfix/internal/spec"
	"roundfix/internal/store"
)

type lostRolloutRunner struct {
	*taskFakeRunner
	turns      []error
	seen       []agent.ExecuteRequest
	closed     []agent.SessionRef
	qaFallback string
}

func (runner *lostRolloutRunner) Run(ctx context.Context, req agent.ExecuteRequest, sink runevent.Sink) (agent.ExecuteResult, error) {
	runner.seen = append(runner.seen, req)
	index := len(runner.seen) - 1
	if index < len(runner.turns) && runner.turns[index] != nil {
		if err := sink.Publish(ctx, runevent.RunEvent{RunID: req.RunID, Source: runevent.SourceAgent, Kind: runevent.KindAgentMessage, Summary: "work started before the runtime lost the turn"}); err != nil {
			return agent.ExecuteResult{}, err
		}
		return agent.ExecuteResult{}, runner.turns[index]
	}
	if strings.Contains(req.Prompt, "Spec QA gate") {
		path, _, err := seededQAReportPathFromPromptForTest(req.Prompt, runner.gitRoot)
		if err != nil {
			return agent.ExecuteResult{}, err
		}
		content, err := os.ReadFile(path)
		if err != nil {
			return agent.ExecuteResult{}, err
		}
		runner.qaFallback = string(content)
	}
	result, runErr := runner.taskFakeRunner.Run(ctx, req, sink)
	if runErr == nil && runner.qaFallback != "" {
		if start := strings.Index(runner.qaFallback, "\n## Agent runtime fallback"); start >= 0 {
			path, _, err := seededQAReportPathFromPromptForTest(req.Prompt, runner.gitRoot)
			if err != nil {
				return result, err
			}
			content, err := os.ReadFile(path)
			if err != nil {
				return result, err
			}
			if err := os.WriteFile(path, append(content, []byte(runner.qaFallback[start:])...), 0o644); err != nil {
				return result, err
			}
		}
	}
	return result, runErr
}
func (runner *lostRolloutRunner) EndSession(_ context.Context, _ agent.RuntimeSpec, session agent.SessionRef) error {
	runner.closed = append(runner.closed, session)
	return nil
}
func lostRolloutErrorForTest() error {
	return &agent.BatchFailureError{Protocol: &agent.ProtocolFailure{Step: "session/prompt", Code: -32603, Message: "Internal error", PromptSent: true, LostRollout: true, Detail: agent.LostRolloutPhrase + " fixture-thread"}}
}
func assertLostRolloutEvents(t *testing.T, sink *captureEventSink, recoveries ...string) {
	t.Helper()
	var got []string
	for _, event := range eventsOfKind(sink, runevent.KindDaemonTask) {
		var payload map[string]any
		if err := json.Unmarshal(event.Payload, &payload); err != nil {
			t.Fatal(err)
		}
		if payload["phase"] != "rollout_lost" {
			continue
		}
		got = append(got, payload["recovery"].(string))
		if payload["retry_spent"] != false || payload["step"] != "session/prompt" || payload["scope_id"] == "" || payload["selection"] == nil || !strings.Contains(payload["detail"].(string), agent.LostRolloutPhrase) {
			t.Fatalf("loss payload: %+v", payload)
		}
		if payload["recovery"] == "fallback" && payload["next_selection"] == nil {
			t.Fatalf("missing next selection: %+v", payload)
		}
	}
	if strings.Join(got, ",") != strings.Join(recoveries, ",") {
		t.Fatalf("recoveries=%v want %v", got, recoveries)
	}
}
func runLostRolloutTask(t *testing.T, fallbacks bool, turns []error, repair bool) (*lostRolloutRunner, TaskCycleResult, *captureEventSink) {
	t.Helper()
	fixture := newTaskCycleFixture(t, []taskSpecSeed{{id: "task_01", taskType: string(spec.TaskTypeBackend)}})
	runner := &lostRolloutRunner{taskFakeRunner: &taskFakeRunner{calls: fixture.calls, gitRoot: fixture.gitRoot}, turns: turns}
	verifier := &taskFakeVerifier{calls: fixture.calls}
	if repair {
		verifier.script = []error{errors.New("fixture Verification failure"), nil}
	}
	engine := fixture.engine(t, runner, verifier, &engineFakeCommitter{calls: fixture.calls}, fixture.worktree)
	plan := fixture.plan()
	profile := selectionProfileForTest(selectionForTest("codex", "preferred-model", "high"))
	if fallbacks {
		profile.Fallbacks = []roundconfig.AgentSelection{selectionForTest("claude", "fallback-model", "high")}
	}
	plan.AgentSelections = selectionProfilesForTest(map[roundconfig.WorkCategory]roundconfig.AgentSelectionProfile{roundconfig.CategoryBackend: profile})
	plan.RuntimeFactory = runtimeFactoryForLifecycleTest(nil)
	result, err := engine.TaskCycle(context.Background(), plan)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(fixture.progress.String(), "roundfix: task task_01 lost its rollout at session/prompt") {
		t.Fatalf("missing loss notice: %s", fixture.progress)
	}
	for i, req := range runner.seen {
		if i > 0 && strings.Contains(req.Session.Name, "rollout") && !strings.Contains(req.Prompt, lostRolloutNotice) {
			t.Fatalf("missing recovery notice: %s", req.Prompt)
		}
	}
	attempts, readErr := fixture.store.AgentSelectionAttemptsForScope(t.Context(), fixture.run.ID, store.AgentSelectionScopeKind("task"), "task_01")
	if readErr != nil {
		t.Fatal(readErr)
	}
	failed := 0
	for _, attempt := range attempts {
		if attempt.Status == store.AgentSelectionStatusFailed {
			failed++
			if attempt.ReasonCode != "rollout_lost" {
				t.Fatalf("failed attempt=%+v", attempt)
			}
		}
	}
	losses := 0
	for _, turn := range turns {
		if _, ok := agent.DescribeLostRollout(turn); ok {
			losses++
		}
	}
	if failed != losses {
		t.Fatalf("failed attempts=%d losses=%d history=%+v", failed, losses, attempts)
	}
	for i, req := range runner.seen {
		for _, prior := range runner.seen[:i] {
			if strings.Contains(req.Session.Name, "rollout") && req.Session.Name == prior.Session.Name {
				t.Fatalf("reused lost session %s", req.Session.Name)
			}
		}
	}
	return runner, result, fixture.sink
}
func TestLostRolloutBeforeFirstHandoffActivatesTheFallback(t *testing.T) {
	runner, result, sink := runLostRolloutTask(t, true, []error{lostRolloutErrorForTest()}, false)
	if result.Completed != 1 || len(runner.seen) != 2 {
		t.Fatalf("result=%+v requests=%v", result, runner.seen)
	}
	req := runner.seen[1]
	if req.Runtime.Model != "fallback-model" || !strings.HasSuffix(req.Session.Name, "-fallback-01-rollout-01") || strings.Contains(req.Prompt, "Verification Feedback") {
		t.Fatalf("recovery=%+v", req)
	}
	assertLostRolloutEvents(t, sink, "fallback")
	event := singleEventOfKind(t, sink, runevent.KindDaemonAgentSelectionFallback)
	if !strings.Contains(string(event.Payload), `"reason_code":"rollout_lost"`) {
		t.Fatalf("fallback=%s", event.Payload)
	}
}
func TestLostRolloutAfterFirstHandoffResumesInANewSession(t *testing.T) {
	runner, result, sink := runLostRolloutTask(t, true, []error{nil, lostRolloutErrorForTest()}, true)
	if result.Completed != 1 || len(runner.seen) != 3 {
		t.Fatalf("result=%+v turns=%d", result, len(runner.seen))
	}
	req := runner.seen[2]
	if req.Runtime.Model != "preferred-model" || !strings.HasSuffix(req.Session.Name, "-rollout-01") || !strings.HasPrefix(req.Prompt, runner.seen[0].Prompt) || !strings.Contains(req.Prompt, runner.seen[1].Prompt) {
		t.Fatalf("recovery=%+v", req)
	}
	assertLostRolloutEvents(t, sink, "new_session")
}
func TestLostRolloutWithoutFallbackResumesTheSameSelection(t *testing.T) {
	runner, result, sink := runLostRolloutTask(t, false, []error{lostRolloutErrorForTest()}, false)
	if result.Completed != 1 || len(runner.seen) != 2 || runner.seen[1].Runtime.Model != "preferred-model" || !strings.HasSuffix(runner.seen[1].Session.Name, "-rollout-01") {
		t.Fatalf("result=%+v turns=%v", result, runner.seen)
	}
	assertLostRolloutEvents(t, sink, "new_session")
}
func TestThirdLostRolloutSettlesRuntimeInfrastructure(t *testing.T) {
	runner, result, sink := runLostRolloutTask(t, true, []error{lostRolloutErrorForTest(), lostRolloutErrorForTest(), lostRolloutErrorForTest()}, false)
	if result.Failed != 1 || len(runner.seen) != 3 || !strings.HasPrefix(result.Outcomes[0].Reason, "runtime infrastructure: lost rollout") {
		t.Fatalf("result=%+v turns=%d", result, len(runner.seen))
	}
	if len(runner.closed) != 3 {
		t.Fatalf("closed=%v", runner.closed)
	}
	assertLostRolloutEvents(t, sink, "fallback", "new_session", "exhausted")
}
func TestQALostRolloutBeforeTheReportFallsBackAndRecordsIt(t *testing.T) {
	fixture := newTaskCycleFixture(t, []taskSpecSeed{{id: "task_01", status: string(spec.StatusCompleted)}})
	runner := &lostRolloutRunner{taskFakeRunner: &taskFakeRunner{calls: fixture.calls, gitRoot: fixture.gitRoot, qaReport: qaReportForTest(spec.VerdictPass)}, turns: []error{lostRolloutErrorForTest()}}
	engine := fixture.engine(t, runner, &taskFakeVerifier{calls: fixture.calls}, &engineFakeCommitter{calls: fixture.calls}, fixture.worktree)
	engine.deps.MechanicalStage = &fakeQAMechanicalStage{}
	plan := fixture.qaPlan()
	plan.AgentSelections = selectionProfilesForTest(map[roundconfig.WorkCategory]roundconfig.AgentSelectionProfile{roundconfig.CategoryQA: selectionProfileForTest(selectionForTest("codex", "preferred-model", "high"), selectionForTest("claude", "fallback-model", "high"))})
	plan.RuntimeFactory = runtimeFactoryForLifecycleTest(nil)
	_, err := engine.TaskCycle(context.Background(), plan)
	if err != nil {
		t.Fatal(err)
	}
	if len(runner.seen) != 2 || runner.seen[1].Runtime.Model != "fallback-model" {
		t.Fatalf("requests=%v", runner.seen)
	}
	path, _, readErr := seededQAReportPathFromPromptForTest(runner.seen[1].Prompt, runner.gitRoot)
	if readErr != nil {
		t.Fatal(readErr)
	}
	finalReport, readErr := os.ReadFile(path)
	if readErr != nil {
		t.Fatal(readErr)
	}
	if !strings.Contains(string(finalReport), "## Agent runtime fallback") || !strings.Contains(string(finalReport), "- retry_spent: false") {
		t.Fatalf("final report=%s", finalReport)
	}
	if !strings.Contains(runner.qaFallback, "## Agent runtime fallback") || !strings.Contains(runner.qaFallback, "- retry_spent: false") || !strings.Contains(runner.qaFallback, "verdict: pending") {
		t.Fatalf("report=%s", runner.qaFallback)
	}
	assertLostRolloutEvents(t, fixture.sink, "fallback")
}

func TestQALostRolloutExhaustionSettlesInsteadOfHalting(t *testing.T) {
	fixture := newTaskCycleFixture(t, []taskSpecSeed{{id: "task_01", status: string(spec.StatusCompleted)}})
	runner := &lostRolloutRunner{taskFakeRunner: &taskFakeRunner{calls: fixture.calls, gitRoot: fixture.gitRoot}, turns: []error{lostRolloutErrorForTest(), lostRolloutErrorForTest(), lostRolloutErrorForTest()}}
	engine := fixture.engine(t, runner, &taskFakeVerifier{calls: fixture.calls}, &engineFakeCommitter{calls: fixture.calls}, fixture.worktree)
	engine.deps.MechanicalStage = &fakeQAMechanicalStage{}
	plan := fixture.qaPlan()
	plan.AgentSelections = selectionProfilesForTest(map[roundconfig.WorkCategory]roundconfig.AgentSelectionProfile{roundconfig.CategoryQA: selectionProfileForTest(selectionForTest("codex", "preferred-model", "high"), selectionForTest("claude", "fallback-model", "high"))})
	plan.RuntimeFactory = runtimeFactoryForLifecycleTest(nil)
	result, err := engine.TaskCycle(context.Background(), plan)
	if err != nil || result.QAAccepted {
		t.Fatalf("result=%+v err=%v", result, err)
	}
	assertLostRolloutEvents(t, fixture.sink, "fallback", "new_session", "exhausted")
	found := false
	for _, event := range eventsOfKind(fixture.sink, runevent.KindDaemonTask) {
		var payload map[string]any
		if err := json.Unmarshal(event.Payload, &payload); err != nil {
			t.Fatal(err)
		}
		if payload["phase"] == "settled" && payload["status"] == "failed" && strings.Contains(string(event.Payload), "runtime infrastructure: lost rollout") {
			found = true
		}
	}
	if !found {
		t.Fatal("missing failed QA settlement with runtime infrastructure reason")
	}
}
