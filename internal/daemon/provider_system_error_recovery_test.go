// Suite: Provider System Error recovery at the Task and QA session boundary.
// Boundary OUT: runtime processes and Verification are package fakes.
package daemon

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"reflect"
	"strings"
	"testing"

	"roundfix/internal/agent"
	roundconfig "roundfix/internal/config"
	"roundfix/internal/runevent"
	"roundfix/internal/spec"
	"roundfix/internal/store"
)

const providerSystemErrorDetailForTest = "Selected model is at capacity. Please try a different model."
const providerSystemErrorNoticeForTest = "The previous Agent Session ended with a provider system error; the working tree holds its changes. Keep them and continue."

type providerSystemErrorRunner struct {
	*taskFakeRunner
	turns      []error
	seen       []agent.ExecuteRequest
	closed     []agent.SessionRef
	qaFallback string
	noOutput   bool
}

func (runner *providerSystemErrorRunner) Run(ctx context.Context, req agent.ExecuteRequest, sink runevent.Sink) (agent.ExecuteResult, error) {
	runner.seen = append(runner.seen, req)
	index := len(runner.seen) - 1
	if index < len(runner.turns) && runner.turns[index] != nil {
		if !runner.noOutput {
			if err := sink.Publish(ctx, runevent.RunEvent{RunID: req.RunID, Source: runevent.SourceAgent, Kind: runevent.KindAgentMessage, Summary: providerSystemErrorDetailForTest}); err != nil {
				return agent.ExecuteResult{}, err
			}
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
	if runErr == nil {
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

func (runner *providerSystemErrorRunner) EndSession(_ context.Context, _ agent.RuntimeSpec, session agent.SessionRef) error {
	runner.closed = append(runner.closed, session)
	return nil
}

func providerSystemErrorForTest() error {
	return &agent.BatchFailureError{Reason: agent.ProviderSystemErrorReason, Protocol: &agent.ProtocolFailure{
		Step: "session/prompt", PromptSent: true, ProviderSystemError: true, Detail: providerSystemErrorDetailForTest,
	}}
}

func assertProviderSystemErrorEvents(t *testing.T, sink *captureEventSink, scope string, recoveries ...string) {
	t.Helper()
	var got []string
	for _, event := range eventsOfKind(sink, runevent.KindDaemonTask) {
		var payload struct {
			Phase         string            `json:"phase"`
			ScopeKind     string            `json:"scope_kind"`
			ScopeID       string            `json:"scope_id"`
			Selection     map[string]string `json:"selection"`
			NextSelection map[string]string `json:"next_selection"`
			Detail        string            `json:"detail"`
			Recovery      string            `json:"recovery"`
			RetrySpent    *bool             `json:"retry_spent"`
		}
		if err := json.Unmarshal(event.Payload, &payload); err != nil {
			t.Fatal(err)
		}
		if payload.Phase != "provider_system_error" {
			continue
		}
		got = append(got, payload.Recovery)
		if payload.ScopeKind != scope || payload.ScopeID == "" || payload.Selection["model"] == "" || payload.Detail != providerSystemErrorDetailForTest || payload.RetrySpent == nil || *payload.RetrySpent {
			t.Fatalf("provider failure payload: %+v", payload)
		}
		if payload.Recovery == "fallback" && payload.NextSelection["model"] == "" {
			t.Fatalf("missing next selection: %+v", payload)
		}
		if payload.Recovery == "exhausted" && payload.NextSelection != nil {
			t.Fatalf("exhaustion names a next selection: %+v", payload)
		}
	}
	if !reflect.DeepEqual(got, recoveries) {
		t.Fatalf("recoveries=%v want %v", got, recoveries)
	}
}

func assertProviderSystemErrorAttempts(t *testing.T, fixture *taskCycleFixture, scope, id string, failures int) {
	t.Helper()
	attempts, err := fixture.store.AgentSelectionAttemptsForScope(t.Context(), fixture.run.ID, store.AgentSelectionScopeKind(scope), id)
	if err != nil {
		t.Fatal(err)
	}
	failed := 0
	for _, attempt := range attempts {
		if attempt.Status == store.AgentSelectionStatusFailed {
			failed++
			if attempt.ReasonCode != "provider_system_error" || !strings.Contains(attempt.Reason, providerSystemErrorDetailForTest) {
				t.Fatalf("failed attempt=%+v", attempt)
			}
		}
	}
	if failed != failures {
		t.Fatalf("failed attempts=%d want=%d history=%+v", failed, failures, attempts)
	}
}

func runProviderSystemErrorTask(t *testing.T, fallbacks bool, turns []error, repair bool) (*providerSystemErrorRunner, TaskCycleResult, *captureEventSink) {
	t.Helper()
	fixture := newTaskCycleFixture(t, []taskSpecSeed{{id: "task_01", taskType: string(spec.TaskTypeBackend)}})
	runner := &providerSystemErrorRunner{taskFakeRunner: &taskFakeRunner{calls: fixture.calls, gitRoot: fixture.gitRoot}, turns: turns}
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
	result, err := engine.TaskCycle(t.Context(), plan)
	if err != nil {
		t.Fatal(err)
	}
	recovery := "exhausted"
	if result.Completed == 1 {
		recovery = "fallback"
	}
	wantProgress := "roundfix: task task_01 hit a provider system error: " + providerSystemErrorDetailForTest + "; recovery: " + recovery + "\n"
	if strings.Count(fixture.progress.String(), wantProgress) != 1 {
		t.Fatalf("provider notice: %s", fixture.progress)
	}
	assertProviderSystemErrorAttempts(t, fixture, "task", "task_01", 1)
	var sessions []agent.SessionRef
	for i, req := range runner.seen {
		if i == 0 || req.Session != runner.seen[i-1].Session {
			sessions = append(sessions, req.Session)
		}
	}
	if !reflect.DeepEqual(runner.closed, sessions) {
		t.Fatalf("closed=%v want=%v", runner.closed, sessions)
	}
	return runner, result, fixture.sink
}

func TestProviderSystemErrorBeforeFirstHandoffActivatesTheFallback(t *testing.T) {
	t.Parallel()
	runner, result, sink := runProviderSystemErrorTask(t, true, []error{providerSystemErrorForTest()}, false)
	if result.Completed != 1 || len(runner.seen) != 2 {
		t.Fatalf("result=%+v requests=%v", result, runner.seen)
	}
	req := runner.seen[1]
	if req.Runtime.Model != "fallback-model" || !strings.HasSuffix(req.Session.Name, "-fallback-01") || !strings.HasPrefix(req.Prompt, runner.seen[0].Prompt) || !strings.Contains(req.Prompt, providerSystemErrorNoticeForTest) || strings.Contains(req.Prompt, "Verification Feedback") {
		t.Fatalf("recovery=%+v", req)
	}
	assertProviderSystemErrorEvents(t, sink, "task", "fallback")
	event := singleEventOfKind(t, sink, runevent.KindDaemonAgentSelectionFallback)
	if !strings.Contains(string(event.Payload), `"reason_code":"provider_system_error"`) {
		t.Fatalf("fallback=%s", event.Payload)
	}
}

func TestProviderSystemErrorAfterFirstHandoffSettlesRuntimeInfrastructure(t *testing.T) {
	t.Parallel()
	runner, result, sink := runProviderSystemErrorTask(t, true, []error{nil, providerSystemErrorForTest()}, true)
	wantReason := "runtime infrastructure: task task_01 hit a provider system error: " + providerSystemErrorDetailForTest
	if result.Failed != 1 || len(runner.seen) != 2 || result.Outcomes[0].Reason != wantReason {
		t.Fatalf("result=%+v turns=%d", result, len(runner.seen))
	}
	if runner.seen[0].Session != runner.seen[1].Session || !strings.Contains(runner.seen[1].Prompt, "Verification Feedback") {
		t.Fatalf("feedback=%+v", runner.seen[1])
	}
	if events := eventsOfKind(sink, runevent.KindDaemonAgentSelectionFallback); len(events) != 0 {
		t.Fatalf("unexpected fallback: %v", events)
	}
	assertProviderSystemErrorEvents(t, sink, "task", "exhausted")
}

func TestProviderSystemErrorWithoutFallbackSettlesRuntimeInfrastructure(t *testing.T) {
	t.Parallel()
	runner, result, sink := runProviderSystemErrorTask(t, false, []error{providerSystemErrorForTest()}, false)
	wantReason := "runtime infrastructure: task task_01 hit a provider system error: " + providerSystemErrorDetailForTest
	if result.Failed != 1 || len(runner.seen) != 1 || result.Outcomes[0].Reason != wantReason {
		t.Fatalf("result=%+v requests=%v", result, runner.seen)
	}
	assertProviderSystemErrorEvents(t, sink, "task", "exhausted")
}

func TestQAProviderSystemErrorBeforeTheReportFallsBackAndRecordsIt(t *testing.T) {
	t.Parallel()
	fixture := newTaskCycleFixture(t, []taskSpecSeed{{id: "task_01", status: string(spec.StatusCompleted)}})
	runner := &providerSystemErrorRunner{taskFakeRunner: &taskFakeRunner{calls: fixture.calls, gitRoot: fixture.gitRoot, qaReport: qaReportForTest(spec.VerdictPass)}, turns: []error{providerSystemErrorForTest()}}
	engine := fixture.engine(t, runner, &taskFakeVerifier{calls: fixture.calls}, &engineFakeCommitter{calls: fixture.calls}, fixture.worktree)
	engine.deps.MechanicalStage = &fakeQAMechanicalStage{}
	plan := fixture.qaPlan()
	plan.AgentSelections = selectionProfilesForTest(map[roundconfig.WorkCategory]roundconfig.AgentSelectionProfile{roundconfig.CategoryQA: selectionProfileForTest(selectionForTest("codex", "preferred-model", "high"), selectionForTest("claude", "fallback-model", "high"))})
	plan.RuntimeFactory = runtimeFactoryForLifecycleTest(nil)
	result, err := engine.TaskCycle(t.Context(), plan)
	if err != nil || !result.QAAccepted {
		t.Fatalf("result=%+v err=%v", result, err)
	}
	if len(runner.seen) != 2 || runner.seen[1].Runtime.Model != "fallback-model" || !strings.Contains(runner.seen[1].Prompt, providerSystemErrorNoticeForTest) || len(runner.closed) != 2 {
		t.Fatalf("requests=%v closed=%v", runner.seen, runner.closed)
	}
	path, _, err := seededQAReportPathFromPromptForTest(runner.seen[1].Prompt, runner.gitRoot)
	if err != nil {
		t.Fatal(err)
	}
	finalReport, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	for _, line := range []string{
		"## Agent runtime fallback",
		"- provider_system_error: codex/preferred-model/high: " + providerSystemErrorDetailForTest,
		"- fallback: claude/fallback-model/high",
		"- retry_spent: false",
	} {
		if strings.Count(runner.qaFallback, line) != 1 || strings.Count(string(finalReport), line) != 1 {
			t.Fatalf("missing or repeated %q: seed=%s final=%s", line, runner.qaFallback, finalReport)
		}
	}
	if !strings.Contains(runner.qaFallback, "verdict: pending") {
		t.Fatalf("report before fallback=%s", runner.qaFallback)
	}
	assertProviderSystemErrorAttempts(t, fixture, "qa", "qa", 1)
	assertProviderSystemErrorEvents(t, fixture.sink, "qa", "fallback")
}

func TestQAProviderSystemErrorExhaustionSettlesInsteadOfHalting(t *testing.T) {
	t.Parallel()
	fixture := newTaskCycleFixture(t, []taskSpecSeed{{id: "task_01", status: string(spec.StatusCompleted)}})
	runner := &providerSystemErrorRunner{taskFakeRunner: &taskFakeRunner{calls: fixture.calls, gitRoot: fixture.gitRoot}, turns: []error{providerSystemErrorForTest()}}
	engine := fixture.engine(t, runner, &taskFakeVerifier{calls: fixture.calls}, &engineFakeCommitter{calls: fixture.calls}, fixture.worktree)
	engine.deps.MechanicalStage = &fakeQAMechanicalStage{}
	plan := fixture.qaPlan()
	plan.AgentSelections = selectionProfilesForTest(map[roundconfig.WorkCategory]roundconfig.AgentSelectionProfile{roundconfig.CategoryQA: selectionProfileForTest(selectionForTest("codex", "preferred-model", "high"))})
	plan.RuntimeFactory = runtimeFactoryForLifecycleTest(nil)
	result, err := engine.TaskCycle(t.Context(), plan)
	if err != nil || result.QAAccepted || result.QAVerdict != spec.VerdictPending || len(runner.seen) != 1 || len(runner.closed) != 1 {
		t.Fatalf("result=%+v err=%v turns=%d closed=%v", result, err, len(runner.seen), runner.closed)
	}
	assertProviderSystemErrorAttempts(t, fixture, "qa", "qa", 1)
	assertProviderSystemErrorEvents(t, fixture.sink, "qa", "exhausted")
	wantReason := "runtime infrastructure: qa qa hit a provider system error: " + providerSystemErrorDetailForTest
	settled := 0
	for _, event := range eventsOfKind(fixture.sink, runevent.KindDaemonTask) {
		var payload map[string]any
		if err := json.Unmarshal(event.Payload, &payload); err != nil {
			t.Fatal(err)
		}
		if payload["phase"] == "settled" && payload["status"] == "failed" && payload["reason"] == wantReason {
			settled++
		}
	}
	if settled != 1 {
		t.Fatalf("failed QA settlements=%d", settled)
	}
	path, _, err := seededQAReportPathFromPromptForTest(runner.seen[0].Prompt, runner.gitRoot)
	if err != nil {
		t.Fatal(err)
	}
	report, err := spec.ReadQAReportFile(path)
	if err != nil || report.Verdict != spec.VerdictPending {
		t.Fatalf("pending report=%+v err=%v", report, err)
	}
}

func TestProviderSystemErrorRecoveryEligibility(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		name     string
		scope    string
		noOutput bool
	}{
		{"task without output falls back", "task", true},
		{"review after output keeps the error", "review", false},
		{"review without output keeps the error", "review", true},
		{"QA with a finished report exhausts", "qa", false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			fixture := newTaskCycleFixture(t, []taskSpecSeed{{id: "task_01"}})
			cause := providerSystemErrorForTest()
			if tc.noOutput {
				cause = &agent.SelectionFailureError{Runtime: "codex", Reason: agent.ProviderSystemErrorReason, Protocol: &agent.ProtocolFailure{
					Step: "session/prompt", PromptSent: true, ProviderSystemError: true, Detail: providerSystemErrorDetailForTest,
				}}
			}
			runner := &providerSystemErrorRunner{taskFakeRunner: &taskFakeRunner{calls: fixture.calls, gitRoot: fixture.gitRoot}, turns: []error{cause}, noOutput: tc.noOutput}
			engine := fixture.engine(t, runner, &taskFakeVerifier{calls: fixture.calls}, &engineFakeCommitter{calls: fixture.calls}, fixture.worktree)
			owner, err := engine.agentSessionOwner(agentSelectionOwnerConfig{
				Profiles: selectionProfilesForTest(map[roundconfig.WorkCategory]roundconfig.AgentSelectionProfile{
					roundconfig.CategoryBackend: selectionProfileForTest(selectionForTest("codex", "preferred-model", "high"), selectionForTest("claude", "fallback-model", "high")),
				}),
				RuntimeFactory: runtimeFactoryForLifecycleTest(nil),
			}, agentSessionScope{RunID: fixture.run.ID, Kind: tc.scope, ID: "fixture-scope", Category: roundconfig.CategoryBackend, Session: agent.SessionRefForTask(fixture.run.ID, "task_01", fixture.gitRoot)})
			if err != nil {
				t.Fatal(err)
			}
			owner.reportPending = func() bool { return false }
			_, runErr := owner.Run(t.Context(), agent.ExecuteRequest{RunID: fixture.run.ID, Prompt: "fixture turn", GitRoot: fixture.gitRoot})
			if err := owner.Close(t.Context()); err != nil {
				t.Fatal(err)
			}
			switch tc.scope {
			case "task":
				if runErr != nil || len(runner.seen) != 2 || runner.seen[1].Runtime.Model != "fallback-model" {
					t.Fatalf("err=%v turns=%v", runErr, runner.seen)
				}
				assertProviderSystemErrorAttempts(t, fixture, tc.scope, "fixture-scope", 1)
				assertProviderSystemErrorEvents(t, fixture.sink, tc.scope, "fallback")
			case "review":
				if runErr != cause || len(runner.seen) != 1 {
					t.Fatalf("err=%v want unchanged %v turns=%v", runErr, cause, runner.seen)
				}
				assertProviderSystemErrorAttempts(t, fixture, tc.scope, "fixture-scope", 0)
				assertProviderSystemErrorEvents(t, fixture.sink, tc.scope)
				if events := eventsOfKind(fixture.sink, runevent.KindDaemonAgentSelectionFallback); len(events) != 0 {
					t.Fatalf("review fell back: %v", events)
				}
			case "qa":
				want := "runtime infrastructure: qa fixture-scope hit a provider system error: " + providerSystemErrorDetailForTest
				if runErr == nil || runErr.Error() != want || len(runner.seen) != 1 {
					t.Fatalf("err=%v want=%s turns=%v", runErr, want, runner.seen)
				}
				assertProviderSystemErrorAttempts(t, fixture, tc.scope, "fixture-scope", 1)
				assertProviderSystemErrorEvents(t, fixture.sink, tc.scope, "exhausted")
			}
		})
	}
}

func TestProviderSystemErrorDoesNotUseLostRolloutBudget(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		name               string
		turns              []error
		providerRecoveries []string
		lostRecoveries     []string
		models             []string
	}{
		{
			name:               "provider recoveries leave the Lost Rollout budget available",
			turns:              []error{providerSystemErrorForTest(), providerSystemErrorForTest(), providerSystemErrorForTest(), lostRolloutErrorForTest()},
			providerRecoveries: []string{"fallback", "fallback", "fallback"},
			lostRecoveries:     []string{"new_session"},
			models:             []string{"preferred-model", "fallback-1", "fallback-2", "fallback-3", "fallback-3"},
		},
		{
			name:               "spent Lost Rollout budget does not prevent provider fallback",
			turns:              []error{lostRolloutErrorForTest(), lostRolloutErrorForTest(), providerSystemErrorForTest()},
			providerRecoveries: []string{"fallback"},
			lostRecoveries:     []string{"fallback", "fallback"},
			models:             []string{"preferred-model", "fallback-1", "fallback-2", "fallback-3"},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			fixture := newTaskCycleFixture(t, []taskSpecSeed{{id: "task_01", taskType: string(spec.TaskTypeBackend)}})
			runner := &providerSystemErrorRunner{taskFakeRunner: &taskFakeRunner{calls: fixture.calls, gitRoot: fixture.gitRoot}, turns: tc.turns}
			engine := fixture.engine(t, runner, &taskFakeVerifier{calls: fixture.calls}, &engineFakeCommitter{calls: fixture.calls}, fixture.worktree)
			plan := fixture.plan()
			plan.AgentSelections = selectionProfilesForTest(map[roundconfig.WorkCategory]roundconfig.AgentSelectionProfile{
				roundconfig.CategoryBackend: selectionProfileForTest(selectionForTest("codex", "preferred-model", "high"),
					selectionForTest("claude", "fallback-1", "high"), selectionForTest("codex", "fallback-2", "high"), selectionForTest("claude", "fallback-3", "high")),
			})
			plan.RuntimeFactory = runtimeFactoryForLifecycleTest(nil)
			result, err := engine.TaskCycle(t.Context(), plan)
			if err != nil || result.Completed != 1 {
				t.Fatalf("result=%+v err=%v", result, err)
			}
			var models []string
			for _, req := range runner.seen {
				models = append(models, req.Runtime.Model)
			}
			if !reflect.DeepEqual(models, tc.models) {
				t.Fatalf("selections=%v want=%v", models, tc.models)
			}
			assertProviderSystemErrorEvents(t, fixture.sink, "task", tc.providerRecoveries...)
			assertLostRolloutEvents(t, fixture.sink, tc.lostRecoveries...)
		})
	}
}
