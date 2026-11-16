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
	"roundfix/internal/phaserecord"
	"roundfix/internal/spec"
	"roundfix/internal/store"
)

type capturePromptRecorder struct {
	entries []phaserecord.PromptEntry
	err     error
}

func (r *capturePromptRecorder) RecordPrompt(ctx context.Context, entry phaserecord.PromptEntry) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	r.entries = append(r.entries, entry)
	return r.err
}

func TestEachPromptLeavesOnePromptEntryWithItsScope(t *testing.T) {
	t.Parallel()
	for _, kind := range []string{"task", "qa", "review", "session"} {
		t.Run(kind, func(t *testing.T) {
			for _, prepared := range []bool{true, false} {
				t.Run(map[bool]string{true: "prepared", false: "plain"}[prepared], func(t *testing.T) {
					runner := &usagePromptRunner{result: agent.ExecuteResult{ACPSessionID: "adapter-session"}}
					var dispatch agent.Runner = runner
					if !prepared {
						dispatch = usageRunnerOnly{runner}
					}
					engine, s, req := usageEngine(t, dispatch)
					recorder := &capturePromptRecorder{}
					engine.deps.PhaseRecord = recorder
					run, _, err := s.Run(t.Context(), req.RunID)
					if err != nil {
						t.Fatal(err)
					}
					req.Session.WorkDir = t.TempDir()
					req.GitRoot = req.Session.WorkDir
					scopeID, specSlug := "task_02", "0397-fixture"
					if kind == "session" {
						scopeID, specSlug = req.Session.Name, run.SpecSlug
					}
					// Both clocks are fixed and distinct to assert start and end independently.
					started := engine.deps.Now()
					ended := started.Add(time.Second)
					calls := 0
					engine.deps.Now = func() time.Time {
						calls++
						if calls == 1 {
							return started
						}
						return ended
					}
					for range 2 {
						if kind == "session" {
							_, err = engine.runAgentSession(t.Context(), nil, req)
						} else {
							owner := &agentSessionOwner{engine: engine, scope: agentSessionScope{Kind: kind, ID: scopeID, Spec: specSlug}}
							_, err = owner.runPrepared(t.Context(), req)
						}
						if err != nil {
							t.Fatal(err)
						}
					}
					if len(recorder.entries) != 2 {
						t.Fatalf("entries = %+v", recorder.entries)
					}
					for i, entry := range recorder.entries {
						want := phaserecord.PromptEntry{Started: ended, Ended: ended, RunID: req.RunID, Spec: specSlug, ScopeKind: kind, ScopeID: scopeID, Repository: run.GitRoot, WorkDir: req.Session.WorkDir, Harness: req.Runtime.ID, Session: runner.result.ACPSessionID}
						if i == 0 {
							want.Started = started
						}
						if !reflect.DeepEqual(entry, want) || entry.Started.After(entry.Ended) {
							t.Fatalf("entry = %+v, want %+v", entry, want)
						}
					}
				})
			}
		})
	}
}

func TestAFailedPhaseRecordWriteLeavesTheOutcomeUnchanged(t *testing.T) {
	t.Parallel()
	for _, name := range []string{"answered", "failed", "stopped", "canceled", "usage_write_failed"} {
		t.Run(name, func(t *testing.T) {
			var cause error
			switch name {
			case "failed":
				cause = errors.New("prompt failed")
			case "stopped":
				cause = ErrStopRequested
			case "canceled":
				cause = context.Canceled
			}
			runner := &usagePromptRunner{err: cause, result: agent.ExecuteResult{Output: "answer", ACPSessionID: "reported", Usage: agent.TurnUsage{Basis: agent.UsageBasisTurn, TotalTokens: 42}}}
			engine, s, req := usageEngine(t, runner)
			recorder := &capturePromptRecorder{err: errors.New("record disk failure")}
			engine.deps.PhaseRecord = recorder
			if name == "usage_write_failed" {
				engine.deps.Runs = rejectingUsageStore{s}
			}
			before, _, err := s.Run(t.Context(), req.RunID)
			if err != nil {
				t.Fatal(err)
			}
			ctx, cancel := context.WithCancel(t.Context())
			defer cancel()
			if name == "canceled" {
				cancel()
			}
			owner := &agentSessionOwner{engine: engine, scope: agentSessionScope{Kind: "task", ID: "task_02"}}
			for range 2 {
				result, err := owner.runPrepared(ctx, req)
				if !errors.Is(err, cause) || !reflect.DeepEqual(result, runner.result) {
					t.Fatalf("outcome = %+v, %v; want %+v, %v", result, err, runner.result, cause)
				}
			}
			if len(recorder.entries) != 2 {
				t.Fatalf("entries = %+v", recorder.entries)
			}
			warning := "roundfix: warning: phase record not written for task task_02: record disk failure\n"
			got := engine.deps.Progress.(*bytes.Buffer).String()
			want := strings.Repeat(warning, 2)
			wantPrompts := int64(2)
			if name == "usage_write_failed" {
				want = strings.Repeat(warning+"roundfix: warning: token usage not recorded for task task_02: usage disk failure\n", 2)
				wantPrompts = 0
			}
			if got != want {
				t.Fatalf("warning = %q, want %q", got, want)
			}
			usage, err := s.RunTokenUsage(t.Context(), req.RunID)
			if err != nil || int64(usage.Total.Prompts) != wantPrompts {
				t.Fatalf("usage = %+v, %v", usage, err)
			}
			if wantPrompts > 0 && (usage.Total.Tokens == nil || *usage.Total.Tokens != 84) {
				t.Fatalf("usage = %+v", usage)
			}
			events, err := s.RunEventsAfter(t.Context(), req.RunID, 0, 10)
			if err != nil || int64(len(events)) != wantPrompts {
				t.Fatalf("events = %+v, %v", events, err)
			}
			for _, event := range events {
				var payload map[string]any
				if err := json.Unmarshal(event.Event.Payload, &payload); err != nil {
					t.Fatal(err)
				}
				if payload["scope_kind"] != "task" || payload["scope_id"] != "task_02" || payload["runtime"] != req.Runtime.ID || payload["total_tokens"] != float64(42) {
					t.Fatalf("event changed: %s", event.Event.Payload)
				}
			}
			after, _, err := s.Run(t.Context(), req.RunID)
			if err != nil || after.State != before.State {
				t.Fatalf("Run state changed: %+v, %v", after, err)
			}
		})
	}
	t.Run("task_settlement", func(t *testing.T) {
		f := newTaskCycleFixture(t, []taskSpecSeed{{id: "task_01"}})
		engine := f.engine(t, &taskFakeRunner{calls: f.calls, gitRoot: f.gitRoot}, &taskFakeVerifier{calls: f.calls}, &engineFakeCommitter{calls: f.calls}, f.worktree)
		recorder := &capturePromptRecorder{err: errors.New("record disk failure")}
		engine.deps.PhaseRecord = recorder
		result, err := engine.TaskCycle(t.Context(), f.plan())
		if err != nil || result.Completed != 1 || result.Failed != 0 || taskStatusOnDisk(t, f.gitRoot, "task_01") != string(spec.StatusCompleted) {
			t.Fatalf("settlement changed: %+v, %v", result, err)
		}
		if strings.Count(f.progress.String(), "phase record not written for session "+f.plan().Session.Name+": record disk failure\n") != len(recorder.entries) || len(recorder.entries) == 0 {
			t.Fatalf("warnings: %s", f.progress)
		}
	})
}

func TestANilPromptRecorderRecordsNothing(t *testing.T) {
	t.Parallel()
	runner := &usagePromptRunner{result: agent.ExecuteResult{Output: "answer"}}
	engine, s, req := usageEngine(t, runner)
	if engine.deps.PhaseRecord != nil {
		t.Fatal("fixture has a recorder")
	}
	result, err := engine.runAgentSession(t.Context(), nil, req)
	if err != nil || !reflect.DeepEqual(result, runner.result) {
		t.Fatalf("outcome changed: %+v, %v", result, err)
	}
	usage, err := s.RunTokenUsage(t.Context(), req.RunID)
	if err != nil || usage.Total.Prompts != 1 || engine.deps.Progress.(*bytes.Buffer).Len() != 0 {
		t.Fatalf("usage = %+v, %v", usage, err)
	}
}

var _ PromptRecorder = (*capturePromptRecorder)(nil)
var _ RunStateStore = (*store.Store)(nil)

// Expose the Run reader and state methods, without a token-usage writer.
type phaseRunStateOnly struct {
	RunStateStore
	reader *store.Store
}

func (s phaseRunStateOnly) Run(ctx context.Context, id string) (store.Run, bool, error) {
	return s.reader.Run(ctx, id)
}

func TestAPromptIsRecordedWithoutAUsageAppender(t *testing.T) {
	t.Parallel()
	runner := &usagePromptRunner{result: agent.ExecuteResult{ACPSessionID: "adapter-session"}}
	engine, s, req := usageEngine(t, runner)
	engine.deps.Runs = phaseRunStateOnly{RunStateStore: s, reader: s}
	recorder := &capturePromptRecorder{}
	engine.deps.PhaseRecord = recorder
	result, err := engine.runAgentSession(t.Context(), nil, req)
	if err != nil || !reflect.DeepEqual(result, runner.result) {
		t.Fatalf("result=%+v, %v", result, err)
	}
	if len(recorder.entries) != 1 || recorder.entries[0].Session != "adapter-session" || recorder.entries[0].Repository != req.Session.WorkDir {
		t.Fatalf("entries=%+v", recorder.entries)
	}
	usage, err := s.RunTokenUsage(t.Context(), req.RunID)
	if err != nil || usage.Total.Prompts != 0 {
		t.Fatalf("usage=%+v, %v", usage, err)
	}
}
