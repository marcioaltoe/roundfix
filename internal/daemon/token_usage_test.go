// Suite: prompt usage recording.
// Invariant: every executed prompt records usage without changing its outcome.
// Boundary IN: owner dispatch, fake prompt runner, real temporary Run Database.
// Boundary OUT: adapter parsing, owned by internal/agent/usage_test.go.
package daemon

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"reflect"
	"strings"
	"testing"
	"time"

	"roundfix/internal/agent"
	roundconfig "roundfix/internal/config"
	"roundfix/internal/runevent"
	"roundfix/internal/spec"
	"roundfix/internal/store"
)

type usagePromptRunner struct {
	result   agent.ExecuteResult
	err      error
	prepared int
	plain    int
}

func (r *usagePromptRunner) Run(context.Context, agent.ExecuteRequest, runevent.Sink) (agent.ExecuteResult, error) {
	r.plain++
	return r.result, r.err
}
func (r *usagePromptRunner) RunPrepared(context.Context, agent.ExecuteRequest, runevent.Sink) (agent.ExecuteResult, error) {
	r.prepared++
	return r.result, r.err
}

// Only expose the Runner interface to exercise the non-prepared dispatch.
type usageRunnerOnly struct{ runner *usagePromptRunner }

func (r usageRunnerOnly) Run(ctx context.Context, req agent.ExecuteRequest, sink runevent.Sink) (agent.ExecuteResult, error) {
	return r.runner.Run(ctx, req, sink)
}

func usageEngine(t *testing.T, runner agent.Runner) (*Engine, *store.Store, agent.ExecuteRequest) {
	t.Helper()
	ctx := context.Background()
	root := t.TempDir()
	s, err := store.Open(ctx, t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = s.Close() })
	run, err := s.CreateRun(ctx, store.CreateRunRequest{Kind: store.KindImplement, GitRoot: root, LocalBranch: "feat/usage", SpecSlug: "usage"})
	if err != nil {
		t.Fatal(err)
	}
	engine := &Engine{deps: Dependencies{Runs: s, Runner: runner, Sink: &captureEventSink{}, Progress: &bytes.Buffer{}, Now: func() time.Time { return time.Date(2026, 10, 1, 12, 4, 0, 0, time.UTC) }}}
	req := agent.ExecuteRequest{RunID: run.ID, Session: agent.SessionRef{Name: "usage-session", WorkDir: root}, Runtime: agent.RuntimeSpec{ID: "codex", Model: "gpt-6.1-sol", ReasoningEffort: "high"}, Prompt: "implement"}
	return engine, s, req
}

func TestEachPromptRecordsItsUsageWithTheOwnerScope(t *testing.T) {
	t.Parallel()
	for _, kind := range []string{"task", "qa", "review"} {
		t.Run(kind, func(t *testing.T) {
			for _, prepared := range []bool{false, true} {
				t.Run(map[bool]string{false: "runner", true: "prepared"}[prepared], func(t *testing.T) {
					runner := &usagePromptRunner{result: agent.ExecuteResult{Usage: agent.TurnUsage{Basis: agent.UsageBasisRequestSum, TotalTokens: 123, Readings: 3}}}
					var dispatch agent.Runner = usageRunnerOnly{runner}
					if prepared {
						dispatch = runner
					}
					engine, s, req := usageEngine(t, dispatch)
					owner := &agentSessionOwner{engine: engine, scope: agentSessionScope{Kind: kind, ID: "work"}, attemptNumber: 2}
					for range 2 {
						result, err := owner.runPrepared(context.Background(), req)
						if err != nil || !reflect.DeepEqual(result, runner.result) {
							t.Fatalf("result=%+v err=%v", result, err)
						}
					}
					report, err := s.RunTokenUsage(context.Background(), req.RunID)
					if err != nil || report.Total.Prompts != 2 || *report.Total.Tokens != 246 || report.Scopes[0].ScopeKind != kind || report.Scopes[0].ScopeID != "work" {
						t.Fatalf("report=%+v err=%v", report, err)
					}
					events, err := s.RunEventsAfter(context.Background(), req.RunID, 0, 10)
					if err != nil || len(events) != 2 {
						t.Fatalf("events=%v err=%v", events, err)
					}
					for _, event := range events {
						var payload map[string]any
						if err := json.Unmarshal(event.Event.Payload, &payload); err != nil {
							t.Fatal(err)
						}
						if payload["attempt"] != float64(2) || payload["runtime"] != "codex" || payload["model"] != "gpt-6.1-sol" || payload["reasoning_effort"] != "high" || payload["session"] != req.Session.Name || payload["readings"] != float64(3) {
							t.Fatalf("payload=%v", payload)
						}
					}
				})
			}
		})
	}
}

func TestAFailedAndAStoppedPromptStillRecordUsage(t *testing.T) {
	t.Parallel()
	for _, cause := range []error{errors.New("agent failed"), ErrStopRequested, context.Canceled} {
		t.Run(cause.Error(), func(t *testing.T) {
			runner := &usagePromptRunner{err: cause, result: agent.ExecuteResult{Output: "partial", Usage: agent.TurnUsage{Basis: agent.UsageBasisTurn, TotalTokens: 42}}}
			engine, s, req := usageEngine(t, runner)
			owner := &agentSessionOwner{engine: engine, scope: agentSessionScope{Kind: "task", ID: "task_01"}, attemptNumber: 1}
			ctx, cancel := context.WithCancel(context.Background())
			cancel()
			result, err := owner.runPrepared(ctx, req)
			if !errors.Is(err, cause) || !reflect.DeepEqual(result, runner.result) {
				t.Fatalf("outcome changed: %+v %v", result, err)
			}
			report, err := s.RunTokenUsage(context.Background(), req.RunID)
			if err != nil || report.Total.Prompts != 1 || *report.Total.Tokens != 42 {
				t.Fatalf("report=%+v err=%v", report, err)
			}
			events, err := s.RunEventsAfter(context.Background(), req.RunID, 0, 10)
			if err != nil || len(events) != 1 {
				t.Fatalf("events=%v err=%v", events, err)
			}
		})
	}
}

func TestAPromptWithoutAnOwnerRecordsASessionScope(t *testing.T) {
	t.Parallel()
	runner := &usagePromptRunner{}
	engine, s, req := usageEngine(t, runner)
	_, err := engine.runAgentSession(context.Background(), nil, req)
	if err != nil {
		t.Fatal(err)
	}
	report, err := s.RunTokenUsage(context.Background(), req.RunID)
	if err != nil || report.Total.Prompts != 1 || report.Total.ReportedPrompts != 0 || report.Total.Tokens != nil || report.Scopes[0].ScopeKind != "session" || report.Scopes[0].ScopeID != req.Session.Name {
		t.Fatalf("report=%+v err=%v", report, err)
	}
}

type rejectingUsageStore struct{ *store.Store }

func (s rejectingUsageStore) AppendTokenUsage(context.Context, store.TokenUsageRecord) error {
	return errors.New("usage disk failure")
}

func TestAFailedUsageWriteLeavesTheOutcomeUnchanged(t *testing.T) {
	t.Parallel()
	for _, owned := range []bool{false, true} {
		t.Run(map[bool]string{false: "session", true: "owner"}[owned], func(t *testing.T) {
			for _, cause := range []error{nil, errors.New("prompt failed"), ErrStopRequested} {
				runner := &usagePromptRunner{err: cause, result: agent.ExecuteResult{Output: "answer", Usage: agent.TurnUsage{Basis: agent.UsageBasisTurn, TotalTokens: 4}}}
				engine, s, req := usageEngine(t, runner)
				engine.deps.Runs = rejectingUsageStore{s}
				var result agent.ExecuteResult
				var err error
				if owned {
					owner := &agentSessionOwner{engine: engine, scope: agentSessionScope{Kind: "task", ID: "task_01"}}
					result, err = owner.runPrepared(context.Background(), req)
				} else {
					result, err = engine.runAgentSession(context.Background(), nil, req)
				}
				if !errors.Is(err, cause) || !reflect.DeepEqual(result, runner.result) {
					t.Fatalf("write changed prompt result: %+v err=%v want=%v", result, err, cause)
				}
				warning := engine.deps.Progress.(*bytes.Buffer).String()
				if !strings.Contains(warning, "roundfix: warning: token usage not recorded for ") || !strings.Contains(warning, "usage disk failure\n") {
					t.Fatalf("warning=%q", warning)
				}
			}
		})
	}
	// Cross the settlement boundary too: compare a successful real Task cycle
	// with and without the rejected usage write.
	for _, reject := range []bool{false, true} {
		t.Run(map[bool]string{false: "cycle-recorded", true: "cycle-rejected"}[reject], func(t *testing.T) {
			fixture := newTaskCycleFixture(t, []taskSpecSeed{{id: "task_01"}})
			runner := &taskFakeRunner{calls: fixture.calls, gitRoot: fixture.gitRoot}
			verifier := &taskFakeVerifier{calls: fixture.calls}
			committer := &engineFakeCommitter{calls: fixture.calls}
			engine := fixture.engine(t, runner, verifier, committer, fixture.worktree)
			if reject {
				engine.deps.Runs = rejectingUsageStore{fixture.store}
			}
			result, err := engine.TaskCycle(context.Background(), fixture.plan())
			if err != nil || result.Completed != 1 || result.Failed != 0 || taskStatusOnDisk(t, fixture.gitRoot, "task_01") != string(spec.StatusCompleted) {
				t.Fatalf("usage write changed settlement: %+v err=%v", result, err)
			}
			if reject && !strings.Contains(fixture.progress.String(), "usage disk failure") {
				t.Fatalf("missing cycle warning: %s", fixture.progress)
			}
		})
	}
}

func (r *usagePromptRunner) Probe(context.Context, agent.ProbeRequest) error { return nil }
func (r *usagePromptRunner) EndSession(context.Context, agent.RuntimeSpec, agent.SessionRef) error {
	return nil
}
func (r usageRunnerOnly) Probe(context.Context, agent.ProbeRequest) error { return nil }
func (r usageRunnerOnly) EndSession(context.Context, agent.RuntimeSpec, agent.SessionRef) error {
	return nil
}

type usageFallbackRunner struct{ usagePromptRunner }

func (r *usageFallbackRunner) RunPrepared(_ context.Context, req agent.ExecuteRequest, _ runevent.Sink) (agent.ExecuteResult, error) {
	result := agent.ExecuteResult{Usage: agent.TurnUsage{Basis: agent.UsageBasisTurn, TotalTokens: 10}}
	if req.Runtime.Model == "preferred-model" {
		return result, &agent.SelectionFailureError{Runtime: "codex", Reason: "quota exhausted"}
	}
	return result, nil
}

func TestUsageRecordsPreferredAndFallbackPrompts(t *testing.T) {
	t.Parallel()
	runner := &usageFallbackRunner{}
	engine, s, req := usageEngine(t, runner)
	owner, err := engine.agentSessionOwner(agentSelectionOwnerConfig{
		Profiles: selectionProfilesForTest(map[roundconfig.WorkCategory]roundconfig.AgentSelectionProfile{
			roundconfig.CategoryBackend: selectionProfileForTest(selectionForTest("codex", "preferred-model", "high"), selectionForTest("claude", "fallback-model", "low")),
		}), RuntimeFactory: runtimeFactoryForLifecycleTest(nil),
	}, agentSessionScope{RunID: req.RunID, Kind: "task", ID: "task_01", Category: roundconfig.CategoryBackend, Session: req.Session})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := engine.runAgentSession(context.Background(), owner, req); err != nil {
		t.Fatal(err)
	}
	report, err := s.RunTokenUsage(context.Background(), req.RunID)
	if err != nil || report.Total.Prompts != 2 || report.Total.Sessions != 2 || report.Total.Tokens == nil || *report.Total.Tokens != 20 {
		t.Fatalf("fallback report=%+v err=%v", report, err)
	}
	events, err := s.RunEventsAfter(context.Background(), req.RunID, 0, 30)
	if err != nil {
		t.Fatal(err)
	}
	attempts := []int{}
	models := []string{}
	sessions := []string{}
	for _, entry := range events {
		if entry.Event.Kind != runevent.KindDaemonTokenUsage {
			continue
		}
		var p struct {
			Attempt int    `json:"attempt"`
			Model   string `json:"model"`
			Session string `json:"session"`
		}
		if err := json.Unmarshal(entry.Event.Payload, &p); err != nil {
			t.Fatal(err)
		}
		attempts = append(attempts, p.Attempt)
		models = append(models, p.Model)
		sessions = append(sessions, p.Session)
	}
	if !reflect.DeepEqual(attempts, []int{1, owner.attemptNumber}) || !reflect.DeepEqual(models, []string{"preferred-model", "fallback-model"}) || !reflect.DeepEqual(sessions, []string{req.Session.Name, req.Session.Name + "-fallback-01"}) {
		t.Fatalf("attempts=%v models=%v sessions=%v", attempts, models, sessions)
	}
}
