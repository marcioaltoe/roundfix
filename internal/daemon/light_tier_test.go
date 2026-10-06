// Boundary: TaskCycle with a fake runner, real temporary Run store and private spend log.
package daemon

import (
	"context"
	"encoding/json"
	"errors"
	"math"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"roundfix/internal/agent"
	roundconfig "roundfix/internal/config"
	"roundfix/internal/gittest"
	"roundfix/internal/lighttier"
	"roundfix/internal/runevent"
	"roundfix/internal/spec"
	"roundfix/internal/store"
)

func lightPlanForTest(t *testing.T, fixture *taskCycleFixture) TaskPlan {
	t.Helper()
	plan := fixture.plan()
	for index := range plan.Tasks {
		plan.Tasks[index].Complexity = "low"
	}
	plan.AgentSelections = selectionProfilesForTest(map[roundconfig.WorkCategory]roundconfig.AgentSelectionProfile{
		roundconfig.CategoryBackend: selectionProfileForTest(selectionForTest("codex", "category-model", "high"), selectionForTest("claude", "category-fallback", "high")),
		roundconfig.CategoryQA:      selectionProfileForTest(selectionForTest("claude", "qa-model", "high")),
	})
	plan.RuntimeFactory = runtimeFactoryForLifecycleTest(nil)
	variable, _ := roundconfig.OpenRouterImplementKey(func(string) string { return "fixture" })
	plan.LightTier = lighttier.Plan{Models: []string{roundconfig.DefaultLightModel, "x-ai/grok-4.5"}, CeilingUSD: 10, KeyVariable: variable, KeyPresent: true, HomeDir: t.TempDir(), Repository: fixture.gitRoot}
	return plan
}

func lightRunnerForTest(f *taskCycleFixture) *selectionLifecycleRunner {
	return &selectionLifecycleRunner{gitRoot: f.gitRoot, statusByTask: map[string]spec.Status{"task_01": spec.StatusCompleted}, sink: f.sink, progress: f.progress}
}

func lightPhasesForTest(t *testing.T, f *taskCycleFixture, phase string) []map[string]any {
	t.Helper()
	var rows []map[string]any
	for _, event := range eventsOfKind(f.sink, runevent.KindDaemonTask) {
		payload := eventPayloadMap(t, event)
		if payload["phase"] == phase {
			rows = append(rows, payload)
		}
	}
	return rows
}

func TestLightTierDispatchOrder(t *testing.T) {
	t.Parallel()
	for _, name := range []string{"light", "medium", "governed", "off", "first-start-fails", "light-starts-fail"} {
		t.Run(name, func(t *testing.T) {
			f := newTaskCycleFixture(t, []taskSpecSeed{{id: "task_01"}})
			runner := lightRunnerForTest(f)
			plan := lightPlanForTest(t, f)
			first := "openrouter/" + plan.LightTier.Models[0]
			second := "openrouter/" + plan.LightTier.Models[1]
			want := first
			switch name {
			case "medium":
				plan.Tasks[0].Complexity = "medium"
				want = "category-model"
			case "governed":
				plan.Tasks[0].Context = []spec.TaskContextRef{{Kind: spec.ContextKindCreates, Path: ".agents/skills/example/SKILL.md"}}
				want = "category-model"
			case "off":
				plan.LightTier = lighttier.Plan{}
				want = "category-model"
			case "first-start-fails":
				runner.prepareErrByModel = map[string]error{first: selectionStartErrForTest("opencode", first)}
				want = second
			case "light-starts-fail":
				runner.prepareErrByModel = map[string]error{first: selectionStartErrForTest("opencode", first), second: selectionStartErrForTest("opencode", second)}
				want = "category-model"
			}
			engine := f.engine(t, runner, &taskFakeVerifier{calls: f.calls}, &engineFakeCommitter{calls: f.calls}, f.worktree)
			result, err := engine.TaskCycle(t.Context(), plan)
			if err != nil || result.Completed != 1 {
				t.Fatalf("cycle = %+v, %v", result, err)
			}
			requests := runner.runRequests()
			if len(requests) != 1 || requests[0].Runtime.Model != want {
				t.Fatalf("work selections = %+v, want %s", requests, want)
			}
			if requests[0].Runtime.ID == "opencode" && (requests[0].Runtime.ReasoningEffort != "" || requests[0].Runtime.OpenRouterKeyVariable != plan.LightTier.KeyVariable) {
				t.Fatal("light runtime has effort or lacks key name")
			}
			if requests[0].Runtime.ID != "opencode" && requests[0].Runtime.OpenRouterKeyVariable != "" {
				t.Fatal("category selection inherited light credential")
			}
			if name == "light-starts-fail" {
				prepared := runner.prepareRequests()
				if len(prepared) != 3 || prepared[0].Runtime.Model != first || prepared[1].Runtime.Model != second || prepared[2].Runtime.Model != want {
					t.Fatalf("start order = %+v", prepared)
				}
			}
			attempts, err := f.store.AgentSelectionAttemptsForScope(t.Context(), f.run.ID, store.AgentSelectionScopeTask, "task_01")
			if err != nil {
				t.Fatal(err)
			}
			if name == "light" && (len(attempts) == 0 || attempts[0].ProfileSource != "light-tier") {
				t.Fatalf("profile source = %+v", attempts)
			}
			qa, err := engine.qaAgentSessionOwner(plan, 2)
			if err != nil || qa.profile.Profile.Preferred.Model != "qa-model" || qa.lightCandidates != 0 {
				t.Fatalf("QA owner changed: %+v %v", qa, err)
			}
		})
	}
}

func TestLightTierSkipsWithAWarning(t *testing.T) {
	t.Parallel()
	for _, reason := range []string{"key_missing", "spend_unreadable", "ceiling_reached"} {
		t.Run(reason, func(t *testing.T) {
			f := newTaskCycleFixture(t, []taskSpecSeed{{id: "task_01"}})
			runner := lightRunnerForTest(f)
			plan := lightPlanForTest(t, f)
			engine := f.engine(t, runner, &taskFakeVerifier{calls: f.calls}, &engineFakeCommitter{calls: f.calls}, f.worktree)
			now := engine.deps.Now()
			if reason == "key_missing" {
				plan.LightTier.KeyPresent = false
			}
			{
				if err := lighttier.AppendSpend(plan.LightTier.HomeDir, now, lighttier.SpendLine{CostUSD: 10, CostSource: "opencode"}); err != nil {
					t.Fatal(err)
				}
			}
			if reason == "spend_unreadable" || reason == "key_missing" {
				path := filepath.Join(plan.LightTier.HomeDir, ".roundfix", "openrouter", "implement", now.UTC().Format("2006-01")+".jsonl")
				if err := os.WriteFile(path, []byte("bad log\n"), 0600); err != nil {
					t.Fatal(err)
				}
			}
			result, err := engine.TaskCycle(t.Context(), plan)
			if err != nil || result.Completed != 1 {
				t.Fatalf("cycle = %+v %v", result, err)
			}
			if runner.runRequests()[0].Runtime.Model != "category-model" {
				t.Fatal("skip changed profile")
			}
			rows := lightPhasesForTest(t, f, "light_tier_skipped")
			if len(rows) != 1 || rows[0]["reason_code"] != reason {
				t.Fatalf("skip events = %v", rows)
			}
			warning := "roundfix: warning: light tier skipped for Task task_01: "
			if strings.Count(f.progress.String(), warning) != 1 || !strings.Contains(f.progress.String(), "; it runs on its backend profile") {
				t.Fatalf("progress = %s", f.progress)
			}
		})
	}
}

type lightCostRunner struct {
	*selectionLifecycleRunner
	costs       []*agent.ReportedCost
	afterPrompt func()
}

func (r *lightCostRunner) RunPrepared(ctx context.Context, req agent.ExecuteRequest, sink runevent.Sink) (agent.ExecuteResult, error) {
	result, err := r.selectionLifecycleRunner.RunPrepared(ctx, req, sink)
	if len(r.costs) > 0 {
		result.Usage.Cost = r.costs[0]
		r.costs = r.costs[1:]
	}
	if r.afterPrompt != nil {
		r.afterPrompt()
	}
	return result, err
}

func TestLightTierRecordsSpend(t *testing.T) {
	t.Parallel()
	f := newTaskCycleFixture(t, []taskSpecSeed{{id: "task_01"}})
	plan := lightPlanForTest(t, f)
	runner := &lightCostRunner{selectionLifecycleRunner: lightRunnerForTest(f), costs: []*agent.ReportedCost{{Amount: 0.12, Currency: "USD"}, nil, {Amount: 0.32, Currency: "USD"}, {Amount: 9, Currency: "EUR"}}}
	engine := f.engine(t, runner, &taskFakeVerifier{calls: f.calls}, &engineFakeCommitter{calls: f.calls}, f.worktree)
	owner, err := engine.taskAgentSessionOwner(t.Context(), plan, plan.Tasks[0], 1)
	if err != nil {
		t.Fatal(err)
	}
	defer owner.Close(t.Context())
	for range 4 {
		if _, err := owner.Run(t.Context(), agent.ExecuteRequest{RunID: plan.RunID, Prompt: "work", GitRoot: plan.WorkDir}); err != nil {
			t.Fatal(err)
		}
	}
	total, err := lighttier.ReadMonth(plan.LightTier.HomeDir, engine.deps.Now())
	if err != nil || math.Abs(total-0.32) > 1e-9 {
		t.Fatalf("spend = %v %v", total, err)
	}
	path := filepath.Join(plan.LightTier.HomeDir, ".roundfix", "openrouter", "implement", engine.deps.Now().UTC().Format("2006-01")+".jsonl")
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	rows := strings.Split(strings.TrimSpace(string(data)), "\n")
	if len(rows) != 4 {
		t.Fatalf("spend rows = %d", len(rows))
	}
	for index, row := range rows {
		var line lighttier.SpendLine
		if err := json.Unmarshal([]byte(row), &line); err != nil {
			t.Fatal(err)
		}
		if line.RunID != plan.RunID || line.Repository != f.gitRoot || line.Spec != plan.Spec.Slug || line.Task != "task_01" || line.Session != owner.activeSession.Name || line.Model != owner.activeRuntime.Model {
			t.Fatalf("identity = %+v", line)
		}
		if (index == 1 || index == 3) && (line.CostSource != "unreported" || line.CostUSD != 0) {
			t.Fatalf("unreported = %+v", line)
		}
	}
	// A failed append warns, preserves the last recorded reading, and a later
	// successful append includes the cost since that reading.
	runner.costs = []*agent.ReportedCost{{Amount: 0.5, Currency: "USD"}, {Amount: 0.7, Currency: "USD"}}
	dir := filepath.Dir(path)
	backup := dir + ".saved"
	runner.afterPrompt = func() {
		if err := os.Rename(dir, backup); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(dir, []byte("obstruction"), 0600); err != nil {
			t.Fatal(err)
		}
	}
	if _, err := owner.Run(t.Context(), agent.ExecuteRequest{RunID: plan.RunID, Prompt: "work", GitRoot: plan.WorkDir}); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(f.progress.String(), "Light Spend Log not recorded") {
		t.Fatal("append failure has no warning")
	}
	runner.afterPrompt = nil
	if err := os.Remove(dir); err != nil {
		t.Fatal(err)
	}
	if err := os.Rename(backup, dir); err != nil {
		t.Fatal(err)
	}
	if _, err := owner.Run(t.Context(), agent.ExecuteRequest{RunID: plan.RunID, Prompt: "work", GitRoot: plan.WorkDir}); err != nil {
		t.Fatal(err)
	}
	total, err = lighttier.ReadMonth(plan.LightTier.HomeDir, engine.deps.Now())
	if err != nil || math.Abs(total-0.7) > 1e-9 {
		t.Fatalf("spend after append recovery = %v %v", total, err)
	}

}

func TestLightTierEscalatesOnce(t *testing.T) {
	t.Parallel()
	for _, name := range []string{"repair-passes", "repair-fails", "category-start-fails", "light-start-fails"} {
		t.Run(name, func(t *testing.T) {
			f := newTaskCycleFixture(t, []taskSpecSeed{{id: "task_01"}})
			plan := lightPlanForTest(t, f)
			// Keep the authored Task low after ReloadTask, matching a real Spec.
			path := taskPathFor(f.gitRoot, taskCycleSlug, "task_01")
			data, err := os.ReadFile(path)
			if err != nil {
				t.Fatal(err)
			}
			if err := os.WriteFile(path, []byte(strings.Replace(string(data), "type: backend\n", "type: backend\ncomplexity: low\n", 1)), 0644); err != nil {
				t.Fatal(err)
			}
			gittest.Run(t, f.gitRoot, "add", ".")
			gittest.Run(t, f.gitRoot, "commit", "-m", "seed light task")
			plan.HeadSHA = strings.TrimSpace(gittest.Run(t, f.gitRoot, "rev-parse", "HEAD"))
			runner := lightRunnerForTest(f)
			if name == "category-start-fails" {
				runner.prepareErrByModel = map[string]error{"category-model": selectionStartErrForTest("codex", "category-model")}
			}
			if name == "light-start-fails" {
				runner.prepareErrByModel = map[string]error{"openrouter/" + plan.LightTier.Models[0]: selectionStartErrForTest("opencode", "unavailable"), "openrouter/" + plan.LightTier.Models[1]: selectionStartErrForTest("opencode", "unavailable")}
			}
			second := error(nil)
			if name == "repair-fails" {
				second = errors.New("still red")
			}
			verifier := &taskFakeVerifier{calls: f.calls, script: []error{errors.New("first failure"), second}, outputByCall: map[int]string{1: "widget does not compile"}}
			engine := f.engine(t, runner, verifier, &engineFakeCommitter{calls: f.calls}, f.worktree)
			result, err := engine.TaskCycle(t.Context(), plan)
			if err != nil {
				t.Fatal(err)
			}
			if name == "repair-fails" && result.Failed != 1 {
				t.Fatalf("expected failed after bounded repair: %+v", result)
			}
			if name != "repair-fails" && result.Completed != 1 {
				t.Fatalf("expected repaired: %+v", result)
			}
			requests := runner.runRequests()
			if len(requests) != 2 {
				t.Fatalf("repair count = %d", len(requests))
			}
			rows := lightPhasesForTest(t, f, "light_tier_escalated")
			if name == "light-start-fails" {
				if len(rows) != 0 || requests[0].Session.Name != requests[1].Session.Name {
					t.Fatal("standard repair changed")
				}
				return
			}
			want := "category-model"
			if name == "category-start-fails" {
				want = "category-fallback"
			}
			if requests[1].Runtime.Model != want || requests[0].Session.Name == requests[1].Session.Name || requests[1].Runtime.OpenRouterKeyVariable != "" {
				t.Fatalf("repair selection = %+v", requests[1])
			}
			if len(rows) != 1 || rows[0]["model"] != requests[0].Runtime.Model || len(rows[0]["failed_commands"].([]any)) != 1 {
				t.Fatalf("escalation events = %+v", rows)
			}
			if !strings.Contains(requests[1].Prompt, requests[0].Prompt) || !strings.Contains(requests[1].Prompt, "working tree holds another model's attempt") || !strings.Contains(requests[1].Prompt, "Verification Feedback") {
				t.Fatal("repair lacks handoff context")
			}
			if strings.Count(f.progress.String(), "escalates from the light tier") != 1 {
				t.Fatal("escalation progress missing or repeated")
			}
			if len(runner.closedSessions()) < 2 {
				t.Fatal("both sessions must close")
			}
		})
	}
}
