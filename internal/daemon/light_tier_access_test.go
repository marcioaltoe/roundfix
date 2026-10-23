// Boundary: TaskCycle, fake sessions, temporary repository and Run store.
package daemon

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"testing"

	"roundfix/internal/agent"
	roundconfig "roundfix/internal/config"
	"roundfix/internal/lighttier"
	"roundfix/internal/runevent"
	"roundfix/internal/store"
)

type lightAccessRunner struct {
	*selectionLifecycleRunner
	outputBeforeFailure bool
	prepareCancel       context.CancelFunc
}

func (r *lightAccessRunner) PrepareSession(ctx context.Context, req agent.ExecuteRequest, sink runevent.Sink) error {
	err := r.selectionLifecycleRunner.PrepareSession(ctx, req, sink)
	if r.prepareCancel != nil && req.Runtime.ID == "opencode" {
		r.prepareCancel()
	}
	if err != nil {
		return err
	}
	return req.Runtime.ValidateRequestedAccessPolicy()
}

func (r *lightAccessRunner) RunPrepared(ctx context.Context, req agent.ExecuteRequest, sink runevent.Sink) (agent.ExecuteResult, error) {
	if r.outputBeforeFailure && req.Runtime.ID == "opencode" {
		if err := sink.Publish(ctx, runevent.RunEvent{RunID: req.RunID, Source: runevent.SourceAgent, Kind: runevent.KindAgentMessage, Summary: "work began"}); err != nil {
			return agent.ExecuteResult{}, err
		}
	}
	return r.selectionLifecycleRunner.RunPrepared(ctx, req, sink)
}

func accessRuntimeFactory(full bool, second string) AgentRuntimeFactory {
	return func(selection roundconfig.AgentSelection) (agent.RuntimeSpec, error) {
		runtime, err := agent.RuntimeFor(agent.RuntimeOptions{Agent: selection.Runtime, Model: selection.Model, ReasoningEffort: selection.ReasoningEffort, EnableFullAccess: full})
		if selection.Model == second {
			runtime.SupportedFullAccessMode = "full-access"
		}
		return runtime, err
	}
}

func TestLightTierSkipsWhenTheLightRuntimeCannotServeTheRequestedAccess(t *testing.T) {
	t.Parallel()
	for _, name := range []string{"full-access", "full-access-key-missing", "runtime-default", "second-model-serves"} {
		t.Run(name, func(t *testing.T) {
			f := newTaskCycleFixture(t, []taskSpecSeed{{id: "task_01"}})
			plan := lightPlanForTest(t, f)
			runner := &lightAccessRunner{selectionLifecycleRunner: lightRunnerForTest(f)}
			second := ""
			want := "category-model"
			if name == "runtime-default" {
				want = "openrouter/" + plan.LightTier.Models[0]
			}
			if name == "second-model-serves" {
				second = "openrouter/" + plan.LightTier.Models[1]
				want = second
			}
			plan.RuntimeFactory = accessRuntimeFactory(name != "runtime-default", second)
			if name == "full-access-key-missing" {
				plan.LightTier.KeyPresent = false
			}
			engine := f.engine(t, runner, &taskFakeVerifier{calls: f.calls}, &engineFakeCommitter{calls: f.calls}, f.worktree)
			result, err := engine.TaskCycle(t.Context(), plan)
			if err != nil || result.Completed != 1 {
				t.Fatalf("cycle=%+v %v", result, err)
			}
			requests := runner.runRequests()
			if len(requests) != 1 || requests[0].Runtime.Model != want {
				t.Fatalf("work=%+v", requests)
			}
			if name != "runtime-default" && requests[0].Runtime.RequestedPolicy() != agent.AccessPolicyFullAccess {
				t.Fatal("access policy weakened")
			}
			if len(lightPhasesForTest(t, f, "light_tier_yielded")) != 0 {
				t.Fatal("unexpected yield")
			}
			skips := lightPhasesForTest(t, f, "light_tier_skipped")
			if want == "category-model" {
				if len(skips) != 1 || skips[0]["reason_code"] != "access_mode_unavailable" {
					t.Fatalf("skips=%v", skips)
				}
				for _, req := range runner.prepareRequests() {
					if req.Runtime.ID == "opencode" {
						t.Fatal("prepared light session")
					}
				}
				const warning = `roundfix: warning: light tier skipped for Task task_01: light runtime opencode offers no mode for access policy "full-access"; it runs on its backend profile`
				if strings.Count(f.progress.String(), warning+"\n") != 1 {
					t.Fatalf("progress=%s", f.progress)
				}
			} else if len(skips) != 0 {
				t.Fatalf("skips=%v", skips)
			}
		})
	}
}

func TestLightTierYieldsToTheCategoryProfileBeforeAgentWork(t *testing.T) {
	t.Parallel()
	for _, name := range []string{"factory-refused", "setup-refused", "prompt-failed-before-work"} {
		t.Run(name, func(t *testing.T) {
			f := newTaskCycleFixture(t, []taskSpecSeed{{id: "task_01"}})
			plan := lightPlanForTest(t, f)
			first, second := "openrouter/"+plan.LightTier.Models[0], "openrouter/"+plan.LightTier.Models[1]
			runner := &lightAccessRunner{selectionLifecycleRunner: lightRunnerForTest(f)}
			cause := error(errors.New("light session refused"))
			code := "light_session_failed"
			plan.RuntimeFactory = accessRuntimeFactory(false, "")
			switch name {
			case "factory-refused":
				factory := plan.RuntimeFactory
				plan.RuntimeFactory = func(s roundconfig.AgentSelection) (agent.RuntimeSpec, error) {
					if s.Runtime == "opencode" {
						return agent.RuntimeSpec{}, cause
					}
					return factory(s)
				}
			case "setup-refused":
				cause = &agent.AccessPolicyError{Kind: agent.AccessPolicyRejected, Runtime: "opencode", Policy: agent.AccessPolicyFullAccess, Predicate: agent.AccessModeAcceptedPredicate}
				code = agent.AccessPolicyRejected
				// The runtime declares support; its adapter still refuses at setup.
				plan.RuntimeFactory = accessRuntimeFactory(true, first)
				runner.prepareErrByModel = map[string]error{first: cause}
			case "prompt-failed-before-work":
				runner.runErrByModel = map[string]error{first: cause}
			}
			engine := f.engine(t, runner, &taskFakeVerifier{calls: f.calls}, &engineFakeCommitter{calls: f.calls}, f.worktree)
			result, err := engine.TaskCycle(t.Context(), plan)
			if err != nil || result.Completed != 1 {
				t.Fatalf("cycle=%+v %v", result, err)
			}
			requests := runner.runRequests()
			if len(requests) == 0 || requests[len(requests)-1].Runtime.Model != "category-model" {
				t.Fatalf("work=%+v", requests)
			}
			if name == "prompt-failed-before-work" && (len(requests) != 2 || requests[0].Prompt != requests[1].Prompt) {
				t.Fatal("category runtime did not receive the same Task prompt")
			}
			if name == "setup-refused" && requests[len(requests)-1].Runtime.RequestedPolicy() != agent.AccessPolicyFullAccess {
				t.Fatal("category access weakened")
			}
			for _, req := range runner.prepareRequests() {
				if req.Runtime.Model == second {
					t.Fatal("opened remaining light model")
				}
			}
			if name != "factory-refused" && len(runner.closedSessions()) < 2 {
				t.Fatal("light and category sessions must close")
			}
			rows := lightPhasesForTest(t, f, "light_tier_yielded")
			if len(rows) != 1 || rows[0]["model"] != first || rows[0]["reason_code"] != code || rows[0]["reason"] != cause.Error() || rows[0]["task"] != "task_01" {
				t.Fatalf("yield=%v", rows)
			}
			fallbacks := eventsOfKind(f.sink, runevent.KindDaemonAgentSelectionFallback)
			if len(fallbacks) != 1 {
				t.Fatalf("fallbacks=%v", fallbacks)
			}
			payload := eventPayloadMap(t, fallbacks[0])
			if payload["next_selection"].(map[string]any)["model"] != "category-model" {
				t.Fatalf("fallback=%v", payload)
			}
			if !runner.fallbackPreparedAfterNotification() || !runner.fallbackPreparedAfterVisibleMessage() {
				t.Fatal("fallback opened before notification")
			}
			warning := fmt.Sprintf("roundfix: warning: light tier yielded for Task task_01 before Agent work: %s; it runs on its backend profile\n", cause)
			if strings.Count(f.progress.String(), warning) != 1 {
				t.Fatalf("progress=%s", f.progress)
			}
			attempts, err := f.store.AgentSelectionAttemptsForScope(t.Context(), f.run.ID, store.AgentSelectionScopeTask, "task_01")
			if err != nil {
				t.Fatal(err)
			}
			found := false
			for _, attempt := range attempts {
				if attempt.Model == first && attempt.Status == store.AgentSelectionStatusFailed && attempt.Reason == cause.Error() {
					found = true
				}
			}
			if !found {
				t.Fatalf("missing failed light attempt: %+v", attempts)
			}
		})
	}
}

func TestLightTierDoesNotYieldAfterWorkOrOutsideTheLightTier(t *testing.T) {
	t.Parallel()
	for _, name := range []string{"after-work-started", "category-refused", "canceled-error", "context-done"} {
		t.Run(name, func(t *testing.T) {
			f := newTaskCycleFixture(t, []taskSpecSeed{{id: "task_01"}})
			plan := lightPlanForTest(t, f)
			runner := &lightAccessRunner{selectionLifecycleRunner: lightRunnerForTest(f)}
			first := "openrouter/" + plan.LightTier.Models[0]
			cause := errors.New("refused")
			ctx := t.Context()
			switch name {
			case "after-work-started":
				runner.outputBeforeFailure = true
				runner.runErrByModel = map[string]error{first: cause}
			case "category-refused":
				plan.LightTier = lighttier.Plan{}
				runner.prepareErrByModel = map[string]error{"category-model": cause}
			case "canceled-error":
				runner.prepareErrByModel = map[string]error{first: fmt.Errorf("setup: %w", context.Canceled)}
			case "context-done":
				canceled, cancel := context.WithCancel(ctx)
				defer cancel()
				ctx = canceled
				runner.prepareCancel = cancel
				runner.prepareErrByModel = map[string]error{first: cause}
			}
			engine := f.engine(t, runner, &taskFakeVerifier{calls: f.calls}, &engineFakeCommitter{calls: f.calls}, f.worktree)
			result, err := engine.TaskCycle(ctx, plan)
			if name == "after-work-started" || name == "category-refused" {
				if err != nil || result.Failed != 1 {
					t.Fatalf("cycle=%+v %v", result, err)
				}
			}
			if len(lightPhasesForTest(t, f, "light_tier_yielded")) != 0 || len(eventsOfKind(f.sink, runevent.KindDaemonAgentSelectionFallback)) != 0 {
				t.Fatal("unexpected yield or fallback")
			}
			if name != "category-refused" {
				for _, req := range runner.prepareRequests() {
					if req.Runtime.Model == "category-model" {
						t.Fatal("opened category session")
					}
				}
			}
		})
	}
}
