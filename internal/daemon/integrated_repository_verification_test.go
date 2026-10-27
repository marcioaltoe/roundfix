package daemon

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"

	"roundfix/internal/runevent"
	"roundfix/internal/spec"
	"roundfix/internal/store"
)

type integratedRecordingVerifier struct {
	mu               sync.Mutex
	requests         []VerifyRequest
	integrated       []VerifyRequest
	fail             bool
	beforeIntegrated func()
	temporaryCalls   int
	integratedHook   func(context.Context) error
}

func (v *integratedRecordingVerifier) Verify(ctx context.Context, req VerifyRequest) (VerifyResult, error) {
	v.mu.Lock()
	defer v.mu.Unlock()
	result := VerifyResult{OutputPath: req.OutputPath}
	if isPreWorkProbeRequest(req) {
		return result, &VerificationCommandError{Command: req.Command, OutputPath: req.OutputPath, Err: errors.New("exit status 1")}
	}
	v.requests = append(v.requests, req)
	if filepath.Base(req.OutputPath) != "integrated-repository.log" {
		return result, nil
	}
	if v.beforeIntegrated != nil {
		v.beforeIntegrated()
	}
	v.integrated = append(v.integrated, req)
	if err := os.MkdirAll(filepath.Dir(req.OutputPath), 0o755); err != nil {
		return result, err
	}
	if err := os.WriteFile(req.OutputPath, []byte("integrated diagnostics\n"), 0o644); err != nil {
		return result, err
	}
	if v.integratedHook != nil {
		if err := v.integratedHook(ctx); err != nil {
			return result, err
		}
	}
	if len(v.integrated) <= v.temporaryCalls {
		return result, &TemporaryVerificationFailureError{CommandFailure: &VerificationCommandError{Command: req.Command, OutputPath: req.OutputPath, Err: errors.New("exit status 75")}}
	}
	if v.fail {
		return result, &VerificationCommandError{Command: req.Command, OutputPath: req.OutputPath, Err: errors.New("exit status 1")}
	}
	return result, nil
}

func integratedFixture(t *testing.T) (*taskCycleFixture, TaskPlan, *taskSchedulerRunner, *integratedRecordingVerifier, *Engine) {
	t.Helper()
	fixture := newTaskCycleFixture(t, []taskSpecSeed{
		{id: "task_01", verification: []string{"check-one"}},
		{id: "task_02", verification: []string{"check-two"}},
	})
	plan := fixture.plan()
	plan.Concurrency = 2
	plan.RepositoryVerification = "  repository check  "
	plan.RepositoryVerificationAtSettlement = true
	runner := newTaskSchedulerRunner("task_01", "task_02")
	runner.statusByTask = map[string]spec.Status{"task_01": spec.StatusCompleted, "task_02": spec.StatusCompleted}
	runner.releaseTask("task_01")
	runner.releaseTask("task_02")
	verifier := &integratedRecordingVerifier{}
	engine := fixture.engineWithTaskWorktrees(t, runner, verifier, &taskSchedulerCommitter{}, fixture.worktree, newFakeTaskWorktrees())
	return fixture, plan, runner, verifier, engine
}

func TestIntegratedRepositoryVerificationRunsOnceAfterAParallelUngatedGraph(t *testing.T) {
	t.Parallel()
	fixture, plan, _, verifier, engine := integratedFixture(t)
	verifier.beforeIntegrated = func() {
		for _, id := range []string{"task_01", "task_02"} {
			if !hasTaskSettlementEvent(t, eventsOfKind(fixture.sink, runevent.KindDaemonTask), id, spec.StatusCompleted) {
				t.Errorf("%s not settled before integrated verification", id)
			}
		}
	}
	result, err := engine.TaskCycle(context.Background(), plan)
	if err != nil || result.Completed != 2 {
		t.Fatalf("TaskCycle = %+v, %v", result, err)
	}
	if len(verifier.integrated) != 1 || len(verifier.requests) != 5 {
		t.Fatalf("integrated=%v; commands=%v", verifier.integrated, verifier.requests)
	}
	req := verifier.integrated[0]
	wantPath := filepath.Join(plan.ArtifactDir, "runs", plan.RunID, "verification", "integrated-repository.log")
	if req.WorkDir != plan.WorkDir || req.Command != "repository check" || req.OutputPath != wantPath {
		t.Fatalf("integrated request = %+v", req)
	}
	if _, err := os.Stat(wantPath); err != nil {
		t.Fatal(err)
	}
	found := false
	for _, event := range eventsOfKind(fixture.sink, runevent.KindDaemonVerification) {
		if event.Summary == "Integrated Repository Verification passed: repository check." {
			found = true
		}
	}
	if !found {
		t.Fatal("missing integrated success event")
	}
}

func TestIntegratedRepositoryVerificationFailureLeavesTaskStatuses(t *testing.T) {
	t.Parallel()
	fixture, plan, runner, verifier, engine := integratedFixture(t)
	// Keep both settled files in the Run tree; the scheduler copy fake replaces
	// whole parallel trees rather than cherry-picking each Task commit.
	plan.Tasks[1].Needs = []string{"task_01"}
	verifier.fail = true
	result, err := engine.TaskCycle(context.Background(), plan)
	want := "Integrated Repository Verification failed: repository check; diagnostics: " + filepath.Join(plan.ArtifactDir, "runs", plan.RunID, "verification", "integrated-repository.log")
	if err != nil || result.IntegratedVerificationFailure != want || result.Completed != 2 || result.Failed != 0 {
		t.Fatalf("TaskCycle = %+v, %v; want %q", result, err, want)
	}
	for _, outcome := range result.Outcomes {
		if outcome.Status != "completed" {
			t.Fatalf("Task status changed: %+v", outcome)
		}
	}
	for _, id := range []string{"task_01", "task_02"} {
		if got := taskStatusOnDisk(t, fixture.gitRoot, id); got != "completed" {
			t.Fatalf("%s status = %s", id, got)
		}
	}
	if len(runner.startedTasks()) != 2 {
		t.Fatalf("Agent requests = %v", runner.startedTasks())
	}
	found := false
	for _, event := range eventsOfKind(fixture.sink, runevent.KindDaemonVerification) {
		if event.Summary == "Integrated Repository Verification failed: repository check." {
			found = true
		}
	}
	if !found {
		t.Fatal("missing integrated failure event")
	}
}

func TestIntegratedRepositoryVerificationDoesNotRun(t *testing.T) {
	t.Parallel()
	for _, name := range []string{"capacity-one", "one-completed", "failed-task", "skipped-task", "switch-off", "empty-command", "gated-graph"} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			fixture, plan, runner, verifier, engine := integratedFixture(t)
			switch name {
			case "capacity-one":
				plan.Concurrency = 1
			case "one-completed":
				plan.Tasks = plan.Tasks[:1]
			case "failed-task":
				runner.errByTask = map[string]error{"task_01": errors.New("agent failed")}
			case "skipped-task":
				plan.Tasks[1].Needs = []string{"task_01"}
				runner.errByTask = map[string]error{"task_01": errors.New("agent failed")}
			case "switch-off":
				plan.RepositoryVerificationAtSettlement = false
			case "empty-command":
				plan.RepositoryVerification = "  "
			case "gated-graph":
				gated := fixture.qaPlan()
				plan.Tasks = gated.Tasks
				plan.Spec = gated.Spec
				for i := range plan.Tasks {
					if plan.Tasks[i].Type == spec.TaskTypeQA {
						plan.Tasks[i].Status = spec.StatusCompleted
					}
				}
			}
			result, err := engine.TaskCycle(context.Background(), plan)
			if err != nil {
				t.Fatal(err)
			}
			if len(verifier.integrated) != 0 || strings.TrimSpace(result.IntegratedVerificationFailure) != "" {
				t.Fatalf("unexpected integrated check: %+v, %+v", verifier.integrated, result)
			}
		})
	}
}

func TestIntegratedRepositoryVerificationRetriesTemporaryFailureExclusively(t *testing.T) {
	t.Parallel()
	for _, temporaryCalls := range []int{1, 2} {
		t.Run(fmt.Sprintf("temporary-calls-%d", temporaryCalls), func(t *testing.T) {
			t.Parallel()
			fixture, plan, _, verifier, engine := integratedFixture(t)
			verifier.temporaryCalls = temporaryCalls
			result, err := engine.TaskCycle(context.Background(), plan)
			if err != nil || len(verifier.integrated) != 2 {
				t.Fatalf("result=%+v err=%v integrated=%v", result, err, verifier.integrated)
			}
			if (result.IntegratedVerificationFailure != "") != (temporaryCalls == 2) {
				t.Fatalf("result=%+v", result)
			}
			found := false
			for _, event := range eventsOfKind(fixture.sink, runevent.KindDaemonVerification) {
				var payload struct {
					Retry int    `json:"retry"`
					Mode  string `json:"mode"`
				}
				if err := json.Unmarshal(event.Payload, &payload); err != nil {
					t.Fatal(err)
				}
				if payload.Retry == 1 && payload.Mode == "exclusive" {
					found = true
				}
			}
			if !found {
				t.Fatal("missing exclusive retry event")
			}
		})
	}
}

func TestIntegratedRepositoryVerificationHonorsStopRequest(t *testing.T) {
	t.Parallel()
	fixture, plan, _, verifier, engine := integratedFixture(t)
	verifier.integratedHook = func(ctx context.Context) error {
		if err := fixture.store.RequestStop(context.Background(), plan.RunID); err != nil {
			return err
		}
		<-ctx.Done()
		return ctx.Err()
	}
	result, err := engine.TaskCycle(context.Background(), plan)
	if !errors.Is(err, ErrStopRequested) || result.Completed != 2 || result.IntegratedVerificationFailure != "" {
		t.Fatalf("result=%+v err=%v", result, err)
	}
}

func TestIntegratedRepositoryVerificationHonorsRunBudget(t *testing.T) {
	t.Parallel()
	_, plan, runner, verifier, engine := integratedFixture(t)
	startedAt := time.Now()
	engine.deps.Now = func() time.Time { return startedAt }
	timer := make(chan time.Time, 1)
	engine.deps.BudgetTimer = func(time.Duration) (<-chan time.Time, func() bool) {
		return timer, func() bool { return true }
	}
	plan.RunStartedAt = startedAt
	plan.BudgetEnabled = true
	plan.MaxRunDuration = time.Hour
	verifier.integratedHook = func(ctx context.Context) error {
		timer <- startedAt.Add(plan.MaxRunDuration)
		<-ctx.Done()
		return ctx.Err()
	}
	result, err := engine.TaskCycle(t.Context(), plan)
	if !errors.Is(err, context.DeadlineExceeded) || result.TerminalOutcome != store.StateBudgetExceeded || result.Completed != 2 || result.IntegratedVerificationFailure != "" {
		t.Fatalf("result=%+v err=%v", result, err)
	}
	if len(runner.startedTasks()) != 2 {
		t.Fatalf("Agent calls=%v", runner.startedTasks())
	}
}
