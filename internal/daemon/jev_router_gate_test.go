// Suite: routed prompt ceiling and accounting.
// Invariant: refused prompts never run; executed routed prompts leave one record without credentials.
// Boundary IN: TaskCycle, session ownership, fallback notifications, and Judge Log append.
// Boundary OUT: key HTTP transport and monthly spend aggregation, owned by internal/jevrouter.
package daemon

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"math"
	"net/http"
	"net/http/httptest"
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

func routerCeilingReason(t *testing.T) string {
	t.Helper()
	q, err := judge.Load()
	if err != nil {
		t.Fatal(err)
	}
	return fmt.Sprintf("jev_ceiling_reached: month's Jev spend US$%.4f of US$%.4f", q.MonthlyCeilingUSD, q.MonthlyCeilingUSD)
}

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
	gate := &fakeJevRouterGate{beforeErrors: []error{routerRefusal(routerCeilingReason(t))}}
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
	if payload["reason_code"] != "jev_ceiling_reached" || !strings.Contains(payload["reason"].(string), routerCeilingReason(t)) {
		t.Fatalf("refusal classification: %+v", payload)
	}
	if gate.beforeCalls != 1 || len(gate.records) != 0 {
		t.Fatalf("gate calls: before=%d after=%d", gate.beforeCalls, len(gate.records))
	}
}

func TestJevRouterCeilingFailsTheTaskAfterWork(t *testing.T) {
	gate := &fakeJevRouterGate{beforeErrors: []error{nil, routerRefusal(routerCeilingReason(t))}}
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

func TestJevRouterUnboundedKeyFallsBackBeforeWork(t *testing.T) {
	q, err := judge.Load()
	if err != nil {
		t.Fatal(err)
	}
	reason := fmt.Sprintf("jev_router_key_unbounded: set a monthly credit limit of at most US$%.4f on the key at OpenRouter", q.MonthlyCeilingUSD)
	gate := &fakeJevRouterGate{beforeErrors: []error{routerRefusal(reason)}}
	fixture, runner, result := routerTaskFixture(t, gate, false)
	if result.Completed != 1 || result.Failed != 0 {
		t.Fatalf("fallback result: %+v", result)
	}
	requests := runner.runRequests()
	if len(requests) != 1 || requests[0].Runtime.Model != "good-model" {
		t.Fatalf("refused prompt ran: %+v", requests)
	}
	payload := eventPayloadMap(t, singleEventOfKind(t, fixture.sink, runevent.KindDaemonAgentSelectionFallback))
	if payload["reason_code"] != "jev_router_key_unbounded" {
		t.Fatalf("classification: %+v", payload)
	}
}

func TestJevRouterUnboundedKeyFailsAfterWork(t *testing.T) {
	gate := &fakeJevRouterGate{beforeErrors: []error{nil, routerRefusal("jev_router_key_unbounded: set a monthly credit limit")}}
	fixture, runner, result := routerTaskFixture(t, gate, true)
	if result.Failed != 1 || result.Completed != 0 {
		t.Fatalf("post-work result: %+v", result)
	}
	if len(runner.runRequests()) != 1 || len(runner.prepareRequests()) != 1 {
		t.Fatal("refused prompt ran or fallback prepared")
	}
	if len(eventsOfKind(fixture.sink, runevent.KindDaemonAgentSelectionFallback)) != 0 {
		t.Fatal("fallback after work")
	}
	if len(result.Outcomes) != 1 || !strings.Contains(result.Outcomes[0].Reason, "jev_router_key_unbounded") {
		t.Fatalf("outcomes: %+v", result.Outcomes)
	}
}

func TestJevRouterGateChecksTheReportedKeyLimit(t *testing.T) {
	questions, err := judge.Load()
	if err != nil {
		t.Fatal(err)
	}
	for _, tc := range []struct{ name, fields, reason string }{
		{"unlimited", `"limit":null,"limit_reset":null`, "jev_router_key_unbounded"},
		{"lifetime", `"limit":1,"limit_reset":null`, "jev_router_key_unbounded"},
		{"over ceiling", fmt.Sprintf(`"limit":%v,"limit_reset":"monthly"`, questions.MonthlyCeilingUSD+1), "jev_router_key_unbounded"},
		{"bounded", fmt.Sprintf(`"limit":%v,"limit_reset":"monthly","limit_remaining":1`, questions.MonthlyCeilingUSD), ""},
		{"exhausted", `"limit":1,"limit_reset":"monthly","limit_remaining":0`, "jev_ceiling_reached"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			calls := 0
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path == "/credits" {
					fmt.Fprint(w, `{"data":{"total_credits":200,"total_usage":0}}`)
					return
				}
				calls++
				if r.Header.Get("Authorization") != "Bearer "+routerSentinelKey {
					t.Error("missing bearer key")
				}
				fmt.Fprintf(w, `{"data":{"usage_monthly":0.5,%s}}`, tc.fields)
			}))
			defer server.Close()
			gate := &jevRouterGate{deps: jevrouter.Deps{MinCreditUSD: 0.5, HomeDir: t.TempDir(), Env: []string{agent.JevRouterKeyEnv + "=" + routerSentinelKey}, Client: server.Client(), Endpoint: server.URL}, now: time.Now}
			usage, err := gate.Before(context.Background())
			if calls != 1 {
				t.Fatalf("key calls=%d", calls)
			}
			if tc.reason == "" {
				if err != nil || usage != 0.5 {
					t.Fatalf("usage=%v err=%v", usage, err)
				}
				return
			}
			var refusal *agent.SelectionFailureError
			if usage != 0 || !errors.As(err, &refusal) || !strings.HasPrefix(refusal.Reason, tc.reason+":") || selectionReasonCode(err) != tc.reason {
				t.Fatalf("usage=%v err=%v", usage, err)
			}
		})
	}
}

func TestJevRouterDefaultGateUsesTheConfiguredCeiling(t *testing.T) {
	t.Setenv("HOME", t.TempDir())
	q, err := judge.Load()
	if err != nil {
		t.Fatal(err)
	}
	fixture := newTaskCycleFixture(t, []taskSpecSeed{{id: "task_01", taskType: string(spec.TaskTypeDocs)}})
	for _, ceiling := range []float64{0, 50} {
		engine, err := NewEngine(Dependencies{Runner: &selectionLifecycleRunner{}, Verifier: &taskFakeVerifier{}, Committer: &engineFakeCommitter{}, Pusher: &engineFakePusher{}, Source: &engineFakeSource{}, Runs: fixture.store, Worktree: fixture.worktree, JevMonthlyCeilingUSD: ceiling})
		if err != nil {
			t.Fatal(err)
		}
		gate, ok := engine.deps.JevRouter.(*jevRouterGate)
		if !ok {
			t.Fatalf("gate type=%T", engine.deps.JevRouter)
		}
		want := q.WithMonthlyCeiling(ceiling).MonthlyCeilingUSD
		if gate.deps.Ceiling != want {
			t.Fatalf("ceiling=%v want=%v", gate.deps.Ceiling, want)
		}
	}
}

func TestJevRouterGateAcceptsAKeyLimitAtTheConfiguredCeiling(t *testing.T) {
	t.Setenv("HOME", t.TempDir())
	t.Setenv(agent.JevRouterKeyEnv, routerSentinelKey)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "Bearer "+routerSentinelKey {
			t.Error("missing bearer key")
		}
		if r.URL.Path == "/credits" {
			fmt.Fprint(w, `{"data":{"total_credits":200,"total_usage":0}}`)
			return
		}
		fmt.Fprint(w, `{"data":{"usage_monthly":0.5,"limit":50,"limit_reset":"monthly","limit_remaining":49.5}}`)
	}))
	defer server.Close()
	fixture := newTaskCycleFixture(t, []taskSpecSeed{{id: "task_01", taskType: string(spec.TaskTypeDocs)}})
	for _, ceiling := range []float64{50, 0} {
		engine, err := NewEngine(Dependencies{Runner: &selectionLifecycleRunner{}, Verifier: &taskFakeVerifier{}, Committer: &engineFakeCommitter{}, Pusher: &engineFakePusher{}, Source: &engineFakeSource{}, Runs: fixture.store, Worktree: fixture.worktree, JevMonthlyCeilingUSD: ceiling})
		if err != nil {
			t.Fatal(err)
		}
		gate := engine.deps.JevRouter.(*jevRouterGate)
		gate.deps.Client, gate.deps.Endpoint = server.Client(), server.URL
		usage, err := gate.Before(context.Background())
		if ceiling == 50 {
			if err != nil || usage != 0.5 {
				t.Fatalf("configured gate: usage=%v err=%v", usage, err)
			}
		} else {
			want := fmt.Sprintf("jev_router_key_unbounded: set a monthly credit limit of at most US$%.4f on the key at OpenRouter", gate.deps.Ceiling)
			var refusal *agent.SelectionFailureError
			if !errors.As(err, &refusal) || refusal.Reason != want || usage != 0 {
				t.Fatalf("default gate: usage=%v err=%v", usage, err)
			}
		}
	}
}

func TestJevRouterCreditLowFallsBackBeforeWork(t *testing.T) {
	t.Setenv("HOME", t.TempDir())
	gate := &fakeJevRouterGate{beforeErrors: []error{routerRefusal(routerCreditLowReason())}}
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
	if payload["reason_code"] != "openrouter_credit_low" || !strings.Contains(payload["reason"].(string), routerCreditLowReason()) {
		t.Fatalf("refusal classification: %+v", payload)
	}
	if gate.beforeCalls != 1 || len(gate.records) != 0 {
		t.Fatalf("gate calls: before=%d after=%d", gate.beforeCalls, len(gate.records))
	}
}

func TestJevRouterCreditLowFailsTheTaskAfterWork(t *testing.T) {
	t.Setenv("HOME", t.TempDir())
	gate := &fakeJevRouterGate{beforeErrors: []error{nil, routerRefusal(routerCreditLowReason())}}
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
	if graph.Tasks[0].Status != spec.StatusFailed || len(result.Outcomes) != 1 || !strings.Contains(result.Outcomes[0].Reason, "openrouter_credit_low") {
		t.Fatalf("Task outcome: %+v", result.Outcomes)
	}
}

func routerCreditLowReason() string {
	return fmt.Sprintf("openrouter_credit_low: OpenRouter credit US$5.1100 is below the US$%.4f floor", jevrouter.DefaultMinCreditUSD)
}

func routerCreditServer(t *testing.T, balance, remaining float64, creditsStatus int) *httptest.Server {
	t.Helper()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "Bearer "+routerSentinelKey {
			t.Error("missing bearer key")
		}
		switch r.URL.Path {
		case "/key":
			fmt.Fprintf(w, `{"data":{"usage_monthly":0.5,"limit":50,"limit_reset":"monthly","limit_remaining":%v}}`, remaining)
		case "/credits":
			w.WriteHeader(creditsStatus)
			fmt.Fprintf(w, `{"data":{"total_credits":%v,"total_usage":10}}`, balance+10)
		default:
			t.Errorf("unexpected path %s", r.URL.Path)
			http.NotFound(w, r)
		}
	}))
	t.Cleanup(server.Close)
	return server
}

func TestJevRouterGateRefusesCreditBelowTheFloor(t *testing.T) {
	t.Setenv("HOME", t.TempDir())
	for _, tc := range []struct {
		name               string
		balance, remaining float64
		refused            bool
	}{
		{"account below floor", 5.11, 40.31, true}, {"account above floor", 20, 40.31, false}, {"key below floor", 200, 3, true}, {"at floor", jevrouter.DefaultMinCreditUSD, 40.31, false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			server := routerCreditServer(t, tc.balance, tc.remaining, http.StatusOK)
			gate := &jevRouterGate{deps: jevrouter.Deps{HomeDir: t.TempDir(), Env: []string{agent.JevRouterKeyEnv + "=" + routerSentinelKey}, Client: server.Client(), Endpoint: server.URL, Ceiling: 50}, now: time.Now}
			runner := &routerPromptRunner{}
			owner, _ := routerPromptOwner(t, &fakeJevRouterGate{}, runner)
			owner.engine.deps.JevRouter = gate
			_, err := owner.runPrepared(context.Background(), owner.activeRequest(agent.ExecuteRequest{RunID: owner.scope.RunID}))
			if !tc.refused {
				if err != nil || len(runner.ran) != 1 {
					t.Fatalf("prompt not sent: %v", err)
				}
				return
			}
			left := min(tc.balance, tc.remaining)
			want := fmt.Sprintf("openrouter_credit_low: OpenRouter credit US$%.4f is below the US$%.4f floor", left, jevrouter.DefaultMinCreditUSD)
			var refusal *agent.SelectionFailureError
			if !errors.As(err, &refusal) || refusal.Runtime != "opencode" || refusal.Reason != want || len(runner.ran) != 0 {
				t.Fatalf("refusal=%v sent=%v", err, runner.ran)
			}
		})
	}
}

func TestJevRouterGateRefusesUnreadableCredits(t *testing.T) {
	t.Setenv("HOME", t.TempDir())
	server := routerCreditServer(t, 200, 40.31, http.StatusForbidden)
	gate := &jevRouterGate{deps: jevrouter.Deps{HomeDir: t.TempDir(), Env: []string{agent.JevRouterKeyEnv + "=" + routerSentinelKey}, Client: server.Client(), Endpoint: server.URL, Ceiling: 50}, now: time.Now}
	runner := &routerPromptRunner{}
	owner, _ := routerPromptOwner(t, &fakeJevRouterGate{}, runner)
	owner.engine.deps.JevRouter = gate
	_, err := owner.runPrepared(context.Background(), owner.activeRequest(agent.ExecuteRequest{RunID: owner.scope.RunID}))
	var refusal *agent.SelectionFailureError
	if !errors.As(err, &refusal) || refusal.Runtime != "opencode" || refusal.Reason != "jev_spend_unreadable: read account credits: HTTP 403" || len(runner.ran) != 0 {
		t.Fatalf("refusal=%v sent=%v", err, runner.ran)
	}
}

func TestJevRouterDefaultGateUsesTheConfiguredCreditFloor(t *testing.T) {
	t.Setenv("HOME", t.TempDir())
	t.Setenv(agent.JevRouterKeyEnv, routerSentinelKey)
	fixture := newTaskCycleFixture(t, []taskSpecSeed{{id: "task_01", taskType: string(spec.TaskTypeDocs)}})
	server := routerCreditServer(t, 5.11, 40.31, http.StatusOK)
	engine, err := NewEngine(Dependencies{Runner: &selectionLifecycleRunner{}, Verifier: &taskFakeVerifier{}, Committer: &engineFakeCommitter{}, Pusher: &engineFakePusher{}, Source: &engineFakeSource{}, Runs: fixture.store, Worktree: fixture.worktree, JevMonthlyCeilingUSD: 50, JevRouterMinCreditUSD: 4})
	if err != nil {
		t.Fatal(err)
	}
	gate := engine.deps.JevRouter.(*jevRouterGate)
	gate.deps.Client, gate.deps.Endpoint = server.Client(), server.URL
	runner := &routerPromptRunner{}
	owner, _ := routerPromptOwner(t, &fakeJevRouterGate{}, runner)
	owner.engine = engine
	owner.engine.deps.Runner = runner
	if _, err := owner.Run(context.Background(), agent.ExecuteRequest{RunID: owner.scope.RunID}); err != nil || len(runner.ran) != 1 {
		t.Fatalf("configured prompt: %v sent=%v", err, runner.ran)
	}
}
