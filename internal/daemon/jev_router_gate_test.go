// Suite: routed prompt ceiling and accounting.
// Invariant: refused prompts never run; executed routed prompts leave one record without credentials.
// Boundary IN: TaskCycle, session ownership, fallback notifications, and Judge Log append.
// Boundary OUT: key HTTP transport and monthly spend aggregation, owned by internal/jevrouter.
package daemon

import (
	"bytes"
	"context"
	"errors"
	"math"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"

	"roundfix/internal/agent"
	roundconfig "roundfix/internal/config"
	"roundfix/internal/jevrouter"
	"roundfix/internal/judge"
	"roundfix/internal/runevent"
	"roundfix/internal/spec"
)

const routerSentinelKey = "router-gate-test-secret-never-log"
const routerCeilingReason = "jev_ceiling_reached: month's Jev spend US$5.0000 of US$5.0000"

type fakeJevRouterGate struct {
	beforeErrors            []error
	beforeCalls             int
	usageBefore, usageAfter float64
	records                 []jevrouter.PromptRecord
	home                    string
	now                     time.Time
	afterErr                error
}

func (gate *fakeJevRouterGate) Before(context.Context) (float64, error) {
	index := gate.beforeCalls
	gate.beforeCalls++
	if index < len(gate.beforeErrors) && gate.beforeErrors[index] != nil {
		return 0, gate.beforeErrors[index]
	}
	return gate.usageBefore, nil
}

func (gate *fakeJevRouterGate) After(ctx context.Context, record jevrouter.PromptRecord) error {
	gate.records = append(gate.records, record)
	if err := ctx.Err(); err != nil {
		return err
	}
	if gate.afterErr != nil {
		return gate.afterErr
	}
	if gate.home == "" {
		return nil
	}
	record.UsageAfter = gate.usageAfter
	return (jevrouter.Ledger{HomeDir: gate.home}).Append(record, gate.now)
}

func routerRefusal(reason string) error {
	return &agent.SelectionFailureError{Runtime: "opencode", Reason: reason}
}

func routerTaskFixture(t *testing.T, gate *fakeJevRouterGate, feedback bool) (*taskCycleFixture, *selectionLifecycleRunner, TaskCycleResult) {
	t.Helper()
	t.Setenv(agent.JevRouterKeyEnv, routerSentinelKey)
	fixture := newTaskCycleFixture(t, []taskSpecSeed{{id: "task_01", taskType: string(spec.TaskTypeDocs)}})
	runner := &selectionLifecycleRunner{
		gitRoot:      fixture.gitRoot,
		statusByTask: map[string]spec.Status{"task_01": spec.StatusCompleted},
		sink:         fixture.sink, progress: fixture.progress,
	}
	verifier := &taskFakeVerifier{calls: fixture.calls}
	if feedback {
		verifier.script = []error{errors.New("verification needs a repair")}
	}
	engine := fixture.engine(t, runner, verifier, &engineFakeCommitter{calls: fixture.calls}, fixture.worktree)
	engine.deps.JevRouter = gate
	plan := fixture.plan()
	plan.AgentSelections = selectionProfilesForTest(map[roundconfig.WorkCategory]roundconfig.AgentSelectionProfile{
		roundconfig.CategoryDocs: selectionProfileForTest(
			selectionForTest("opencode", agent.JevRouterModel, ""),
			selectionForTest("codex", "good-model", "high"),
		),
	})
	plan.RuntimeFactory = runtimeFactoryForLifecycleTest(nil)
	result, err := engine.TaskCycle(context.Background(), plan)
	if err != nil {
		t.Fatalf("TaskCycle: %v", err)
	}
	assertNoRouterSecret(t, fixture.progress.String(), fixture.sink)
	return fixture, runner, result
}

func assertNoRouterSecret(t *testing.T, text string, sink *captureEventSink) {
	t.Helper()
	if strings.Contains(text, routerSentinelKey) {
		t.Fatal("sentinel key reached a line")
	}
	for _, event := range sink.snapshot() {
		if strings.Contains(event.Summary, routerSentinelKey) || bytes.Contains(event.Payload, []byte(routerSentinelKey)) {
			t.Fatal("sentinel key reached a Run Event")
		}
	}
}

func TestJevRouterCeilingFallsBackBeforeWork(t *testing.T) {
	gate := &fakeJevRouterGate{beforeErrors: []error{routerRefusal(routerCeilingReason)}}
	fixture, runner, result := routerTaskFixture(t, gate, false)
	if result.Completed != 1 || result.Failed != 0 {
		t.Fatalf("fallback result: %+v", result)
	}
	if !runner.fallbackPreparedAfterNotification() || !runner.fallbackPreparedAfterVisibleMessage() {
		t.Fatal("fallback prepared before both notifications")
	}
	requests := runner.runRequests()
	if len(requests) != 1 || requests[0].Runtime.Model != "good-model" {
		t.Fatalf("refused prompt reached runner: %+v", requests)
	}
	payload := eventPayloadMap(t, singleEventOfKind(t, fixture.sink, runevent.KindDaemonAgentSelectionFallback))
	if payload["reason_code"] != "jev_ceiling_reached" || !strings.Contains(payload["reason"].(string), routerCeilingReason) {
		t.Fatalf("refusal classification: %+v", payload)
	}
	if gate.beforeCalls != 1 || len(gate.records) != 0 {
		t.Fatalf("gate calls: before=%d after=%d", gate.beforeCalls, len(gate.records))
	}
}

func TestJevRouterCeilingFailsTheTaskAfterWork(t *testing.T) {
	gate := &fakeJevRouterGate{beforeErrors: []error{nil, routerRefusal(routerCeilingReason)}}
	fixture, runner, result := routerTaskFixture(t, gate, true)
	if result.Failed != 1 || result.Completed != 0 {
		t.Fatalf("post-work result: %+v", result)
	}
	if len(runner.runRequests()) != 1 || len(runner.prepareRequests()) != 1 {
		t.Fatal("refused feedback prompt ran or fallback prepared")
	}
	if len(eventsOfKind(fixture.sink, runevent.KindDaemonAgentSelectionFallback)) != 0 {
		t.Fatal("fallback activated after work")
	}
	if gate.beforeCalls != 2 || len(gate.records) != 1 {
		t.Fatalf("gate calls: before=%d after=%d", gate.beforeCalls, len(gate.records))
	}
	record := gate.records[0]
	if record.RunID != fixture.run.ID || record.Spec != taskCycleSlug || record.ScopeKind != "task" || record.ScopeID != "task_01" || record.Category != "docs" || record.Attempt < 1 {
		t.Fatalf("Task prompt identity: %+v", record)
	}
	graph, err := spec.Load(fixture.specsRoot, taskCycleSlug)
	if err != nil {
		t.Fatal(err)
	}
	if graph.Tasks[0].Status != spec.StatusFailed || len(result.Outcomes) != 1 || !strings.Contains(result.Outcomes[0].Reason, "jev_ceiling_reached") {
		t.Fatalf("Task outcome: %+v", result.Outcomes)
	}
}

func TestJevRouterUnreadableSpendFallsBack(t *testing.T) {
	gate := &fakeJevRouterGate{beforeErrors: []error{routerRefusal("jev_spend_unreadable: read key usage: HTTP 503")}}
	fixture, runner, result := routerTaskFixture(t, gate, false)
	if result.Completed != 1 || len(runner.runRequests()) != 1 || runner.runRequests()[0].Runtime.Model != "good-model" {
		t.Fatalf("fallback result: %+v", result)
	}
	if !runner.fallbackPreparedAfterNotification() || !runner.fallbackPreparedAfterVisibleMessage() {
		t.Fatal("fallback preceded notification")
	}
	payload := eventPayloadMap(t, singleEventOfKind(t, fixture.sink, runevent.KindDaemonAgentSelectionFallback))
	if payload["reason_code"] != "jev_spend_unreadable" {
		t.Fatalf("classification: %+v", payload)
	}
	if gate.beforeCalls != 1 || len(gate.records) != 0 {
		t.Fatal("refused prompt was recorded or fallback gated")
	}
}

type routerPromptRunner struct {
	fallbackBoundaryRunner
	result agent.ExecuteResult
	err    error
	cancel context.CancelFunc
}

func (runner *routerPromptRunner) RunPrepared(ctx context.Context, req agent.ExecuteRequest, sink runevent.Sink) (agent.ExecuteResult, error) {
	runner.ran = append(runner.ran, req.Runtime.Model)
	if runner.cancel != nil {
		runner.cancel()
	}
	return runner.result, runner.err
}
func (runner *routerPromptRunner) Run(ctx context.Context, req agent.ExecuteRequest, sink runevent.Sink) (agent.ExecuteResult, error) {
	return runner.RunPrepared(ctx, req, sink)
}

func routerPromptOwner(t *testing.T, gate *fakeJevRouterGate, runner *routerPromptRunner) (*agentSessionOwner, *captureEventSink) {
	t.Helper()
	t.Setenv(agent.JevRouterKeyEnv, routerSentinelKey)
	sink := &captureEventSink{}
	owner := fallbackBoundaryOwner(t, sink, runner)
	owner.engine.deps.JevRouter = gate
	owner.scope.Spec = "router-spec"
	owner.scope.Category = roundconfig.CategoryDocs
	owner.attemptNumber = 3
	owner.active = true
	owner.activeRuntime = agent.RuntimeSpec{ID: "opencode-custom", Model: agent.JevRouterModel}
	owner.activeSession = owner.scope.Session
	return owner, sink
}

func TestJevRouterPromptAppendsOneLine(t *testing.T) {
	input, output := int64(31), int64(7)
	gate := &fakeJevRouterGate{home: t.TempDir(), now: time.Date(2026, 10, 2, 12, 0, 0, 0, time.UTC), usageBefore: 1.2, usageAfter: 1.35}
	runner := &routerPromptRunner{result: agent.ExecuteResult{Output: "prompt output must not reach Judge Log", Usage: agent.TurnUsage{InputTokens: &input, OutputTokens: &output}}}
	owner, sink := routerPromptOwner(t, gate, runner)
	ticks := []time.Time{gate.now, gate.now.Add(250 * time.Millisecond)}
	owner.engine.deps.Now = func() time.Time {
		if len(ticks) > 1 {
			now := ticks[0]
			ticks = ticks[1:]
			return now
		}
		return ticks[0]
	}
	result, err := owner.Run(context.Background(), agent.ExecuteRequest{RunID: owner.scope.RunID, Prompt: "private prompt"})
	if err != nil || !reflect.DeepEqual(result, runner.result) {
		t.Fatalf("prompt result changed: %+v, %v", result, err)
	}
	if gate.beforeCalls != 1 || len(gate.records) != 1 {
		t.Fatal("expected one Before and one After")
	}
	rows, err := judge.ReadMonth(context.Background(), gate.home, gate.now)
	if err != nil || len(rows) != 1 {
		t.Fatalf("Judge Log: %d lines, %v", len(rows), err)
	}
	row := rows[0]
	if row.Judgment != "router-prompt" || math.Abs(row.CostUSD-0.15) > 1e-12 || row.InputTokens != 31 || row.OutputTokens != 7 || row.LatencyMS != 250 || row.Attempts != 3 || row.Spec != "router-spec" || row.QuestionID != "docs" || row.Target != owner.scope.RunID+" task "+owner.scope.ID || row.Repository != owner.scope.Session.WorkDir || row.Outcome != "clear" {
		t.Fatalf("router record: %+v", row)
	}
	raw, err := os.ReadFile(filepath.Join(gate.home, ".roundfix/judge/2026-10.jsonl"))
	if err != nil {
		t.Fatal(err)
	}
	assertNoRouterSecret(t, string(raw), sink)
	if strings.Contains(string(raw), "private prompt") || strings.Contains(string(raw), runner.result.Output) {
		t.Fatal("prompt or answer reached Judge Log")
	}
}

func TestNonRoutedPromptCallsNoGate(t *testing.T) {
	for _, tc := range []struct{ name, runtime, model string }{
		{"other OpenCode model", "opencode", "other-model"},
		{"other runtime", "codex", agent.JevRouterModel},
	} {
		t.Run(tc.name, func(t *testing.T) {
			gate := &fakeJevRouterGate{beforeErrors: []error{errors.New("gate must not run")}}
			runner := &routerPromptRunner{result: agent.ExecuteResult{Output: "unchanged output"}}
			owner, sink := routerPromptOwner(t, gate, runner)
			owner.activeRuntime = agent.RuntimeSpec{ID: tc.runtime, Model: tc.model}
			result, err := owner.Run(context.Background(), agent.ExecuteRequest{RunID: owner.scope.RunID})
			if err != nil || !reflect.DeepEqual(result, runner.result) {
				t.Fatalf("non-routed result: %+v %v", result, err)
			}
			if gate.beforeCalls != 0 || len(gate.records) != 0 || len(runner.ran) != 1 {
				t.Fatal("non-routed prompt gated or did not run")
			}
			assertNoRouterSecret(t, owner.engine.deps.Progress.(*bytes.Buffer).String(), sink)
		})
	}
}

func TestJevRouterAfterFailurePreservesPromptResult(t *testing.T) {
	for _, failed := range []bool{false, true} {
		name := "success"
		if failed {
			name = "failure"
		}
		t.Run(name, func(t *testing.T) {
			gate := &fakeJevRouterGate{afterErr: errors.New("append denied")}
			runner := &routerPromptRunner{result: agent.ExecuteResult{Output: "partial output"}}
			if failed {
				runner.err = errors.New("prompt failed")
			}
			owner, sink := routerPromptOwner(t, gate, runner)
			result, err := owner.Run(context.Background(), agent.ExecuteRequest{RunID: owner.scope.RunID})
			if !reflect.DeepEqual(result, runner.result) || !errors.Is(err, runner.err) {
				t.Fatalf("After changed prompt result: %+v %v", result, err)
			}
			if len(gate.records) != 1 || gate.records[0].Failed != failed {
				t.Fatal("missing or incorrect failed prompt record")
			}
			progress := owner.engine.deps.Progress.(*bytes.Buffer).String()
			if !strings.Contains(progress, "append denied") {
				t.Fatal("After error absent from progress")
			}
			assertNoRouterSecret(t, progress, sink)
		})
	}
}

func TestJevRouterCanceledPromptStillRecords(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	gate := &fakeJevRouterGate{home: t.TempDir(), now: time.Now()}
	runner := &routerPromptRunner{err: context.Canceled, cancel: cancel}
	owner, _ := routerPromptOwner(t, gate, runner)
	_, err := owner.Run(ctx, agent.ExecuteRequest{RunID: owner.scope.RunID})
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("prompt error: %v", err)
	}
	rows, readErr := judge.ReadMonth(context.Background(), gate.home, gate.now)
	if readErr != nil || len(rows) != 1 || rows[0].Outcome != "skipped" {
		t.Fatalf("canceled prompt log: %+v %v", rows, readErr)
	}
}

// Expose only Runner, exercising the compatibility path without RunPrepared.
type routerUnpreparedRunner struct{ agent.Runner }

func TestJevRouterPlainRunnerRecordsEveryPrompt(t *testing.T) {
	gate := &fakeJevRouterGate{home: t.TempDir(), now: time.Now(), usageBefore: 0.1, usageAfter: 0.2}
	runner := &routerPromptRunner{result: agent.ExecuteResult{Output: "work"}}
	owner, sink := routerPromptOwner(t, gate, runner)
	owner.engine.deps.Runner = routerUnpreparedRunner{Runner: runner}
	for range 2 {
		if _, err := owner.Run(context.Background(), agent.ExecuteRequest{RunID: owner.scope.RunID}); err != nil {
			t.Fatal(err)
		}
	}
	rows, err := judge.ReadMonth(context.Background(), gate.home, gate.now)
	if err != nil || len(rows) != 2 || gate.beforeCalls != 2 || len(runner.ran) != 2 {
		t.Fatalf("repeated prompts: lines=%d before=%d runs=%d error=%v", len(rows), gate.beforeCalls, len(runner.ran), err)
	}
	for _, row := range rows {
		if row.Judgment != "router-prompt" || row.InputTokens != 0 || row.OutputTokens != 0 {
			t.Fatalf("unreported tokens: %+v", row)
		}
	}
	assertNoRouterSecret(t, owner.engine.deps.Progress.(*bytes.Buffer).String(), sink)
}

func TestJevRouterKeyRefusalKeepsItsClassification(t *testing.T) {
	owner, sink := routerPromptOwner(t, &fakeJevRouterGate{}, &routerPromptRunner{})
	owner.active = false
	owner.profile.Profile = selectionProfileForTest(selectionForTest("opencode", agent.JevRouterModel, ""), selectionForTest("codex", "good-model", "high"))
	runner := &fallbackBoundaryRunner{prepareErrByModel: map[string]error{
		agent.JevRouterModel: routerRefusal(agent.JevRouterKeyMissing + ": ROUNDFIX_OPENROUTER_API_KEY is not set"),
	}}
	owner.engine.deps.Runner = runner
	if _, err := owner.Run(context.Background(), agent.ExecuteRequest{RunID: owner.scope.RunID}); err != nil {
		t.Fatal(err)
	}
	payload := eventPayloadMap(t, singleEventOfKind(t, sink, runevent.KindDaemonAgentSelectionFallback))
	if payload["reason_code"] != agent.JevRouterKeyMissing {
		t.Fatalf("key classification: %+v", payload)
	}
	assertNoRouterSecret(t, owner.engine.deps.Progress.(*bytes.Buffer).String(), sink)
}

func TestJevRouterDefaultGateUsesProcessHomeAndJudgeCeiling(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv(agent.JevRouterKeyEnv, routerSentinelKey)
	fixture := newTaskCycleFixture(t, []taskSpecSeed{{id: "task_01", taskType: string(spec.TaskTypeDocs)}})
	engine := fixture.engine(t, &selectionLifecycleRunner{}, &taskFakeVerifier{calls: fixture.calls}, &engineFakeCommitter{calls: fixture.calls}, fixture.worktree)
	gate, ok := engine.deps.JevRouter.(*jevRouterGate)
	if !ok {
		t.Fatalf("default gate type = %T", engine.deps.JevRouter)
	}
	questions, err := judge.Load()
	if err != nil {
		t.Fatal(err)
	}
	if gate.deps.HomeDir != home || gate.deps.Ceiling != questions.MonthlyCeilingUSD || gate.deps.Endpoint != "https://openrouter.ai/api/v1" || gate.now == nil {
		t.Fatal("default gate does not use process Home, loaded ceiling, production endpoint, and clock")
	}
	if !reflect.DeepEqual(gate.deps.Env, os.Environ()) {
		t.Fatal("default gate does not use process environment")
	}
}
