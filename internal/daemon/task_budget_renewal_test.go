package daemon

// Suite: renewing Implement Run Budget
// Invariant: each timely Task settlement renews one allowance while a late settlement cannot extend it.
// Boundary IN: TaskCycle scheduling, settlement, QA routing, cancellation, and result reporting.
// Boundary OUT: CLI setup and post-cycle work, covered by internal/cli/implement_budget_renewal_test.go.

import (
	"context"
	"errors"
	"runtime"
	"strings"
	"sync"
	"testing"
	"time"

	"roundfix/internal/agent"
	"roundfix/internal/runevent"
	"roundfix/internal/spec"
	"roundfix/internal/store"
	"roundfix/internal/testwait"
)

type budgetTestClock struct {
	mu  sync.Mutex
	now time.Time
}

func (clock *budgetTestClock) Now() time.Time {
	clock.mu.Lock()
	defer clock.mu.Unlock()
	return clock.now
}

func (clock *budgetTestClock) Set(now time.Time) {
	clock.mu.Lock()
	clock.now = now
	clock.mu.Unlock()
}

func budgetedPlan(plan TaskPlan, startedAt time.Time, maximum time.Duration) TaskPlan {
	plan.RunStartedAt = startedAt
	plan.BudgetEnabled = true
	plan.MaxRunDuration = maximum
	return plan
}

func TestTaskBudgetRenewsAtEachTaskSettlement(t *testing.T) {
	fixture := newTaskCycleFixture(t, []taskSpecSeed{
		{id: "task_01"},
		{id: "task_02", needs: []string{"task_01"}},
	})
	startedAt := time.Now()
	const maximum = time.Hour
	clock := &budgetTestClock{now: startedAt}
	settlements := 0
	runner := &taskFakeRunner{
		calls: fixture.calls, gitRoot: fixture.gitRoot,
		statusByTask: map[string]spec.Status{"task_01": spec.StatusCompleted, "task_02": spec.StatusCompleted},
		afterTask: func(string) {
			settlements++
			clock.Set(startedAt.Add(time.Duration(settlements) * 54 * time.Minute))
		},
	}
	engine := fixture.engine(t, runner, &taskFakeVerifier{calls: fixture.calls}, &engineFakeCommitter{calls: fixture.calls}, fixture.worktree)
	engine.deps.Now = clock.Now

	result, err := engine.TaskCycle(context.Background(), budgetedPlan(fixture.plan(), startedAt, maximum))

	if err != nil {
		t.Fatalf("TaskCycle() error = %v, want two timely settlements", err)
	}
	if result.TerminalOutcome != "" || result.Completed != 2 {
		t.Fatalf("TaskCycle() result = %+v, want two completed Tasks without a terminal outcome", result)
	}
	wantDeadline := startedAt.Add(2*54*time.Minute + maximum)
	if !result.BudgetDeadline.Equal(wantDeadline) {
		t.Fatalf("BudgetDeadline = %s, want %s", result.BudgetDeadline, wantDeadline)
	}
}

func TestTaskBudgetStartsTheQAGateWithARenewedAllowance(t *testing.T) {
	fixture := newTaskCycleFixture(t, []taskSpecSeed{{id: "task_01"}})
	startedAt := time.Now()
	const maximum = time.Hour
	clock := &budgetTestClock{now: startedAt}
	runner := &taskFakeRunner{
		calls: fixture.calls, gitRoot: fixture.gitRoot,
		statusByTask: map[string]spec.Status{"task_01": spec.StatusCompleted},
		qaReport:     qaReportForTest(spec.VerdictPass),
		afterTask:    func(string) { clock.Set(startedAt.Add(54 * time.Minute)) },
		afterQA:      func() { clock.Set(startedAt.Add(108 * time.Minute)) },
	}
	engine := fixture.engine(t, runner, &taskFakeVerifier{calls: fixture.calls}, &engineFakeCommitter{calls: fixture.calls}, fixture.worktree)
	engine.deps.Now = clock.Now
	engine.deps.MechanicalStage = &fakeQAMechanicalStage{}

	result, err := engine.TaskCycle(context.Background(), budgetedPlan(fixture.qaPlan(), startedAt, maximum))

	if err != nil {
		t.Fatalf("TaskCycle() error = %v, want QA to use the renewed allowance", err)
	}
	if result.QAVerdict != spec.VerdictPass || !result.QAAccepted {
		t.Fatalf("QA result = verdict %q accepted=%t, want accepted pass", result.QAVerdict, result.QAAccepted)
	}
	wantDeadline := startedAt.Add(108*time.Minute + maximum)
	if !result.BudgetDeadline.Equal(wantDeadline) {
		t.Fatalf("BudgetDeadline = %s, want QA settlement deadline %s", result.BudgetDeadline, wantDeadline)
	}
}

func TestTaskBudgetSettlementAfterTheDeadlineDoesNotRenew(t *testing.T) {
	fixture := newTaskCycleFixture(t, []taskSpecSeed{
		{id: "task_01"},
		{id: "task_02", needs: []string{"task_01"}},
	})
	startedAt := time.Now()
	const maximum = time.Hour
	clock := &budgetTestClock{now: startedAt}
	runner := &taskFakeRunner{calls: fixture.calls, gitRoot: fixture.gitRoot, statusByTask: map[string]spec.Status{"task_01": spec.StatusCompleted, "task_02": spec.StatusCompleted}}
	committer := &engineFakeCommitter{calls: fixture.calls, afterCommit: func(context.Context, CommitRequest) error {
		clock.Set(startedAt.Add(maximum))
		return nil
	}}
	engine := fixture.engine(t, runner, &taskFakeVerifier{calls: fixture.calls}, committer, fixture.worktree)
	engine.deps.Now = clock.Now

	result, err := engine.TaskCycle(context.Background(), budgetedPlan(fixture.plan(), startedAt, maximum))

	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("TaskCycle() error = %v, want DeadlineExceeded", err)
	}
	if result.Completed != 1 || result.TerminalOutcome != store.StateBudgetExceeded {
		t.Fatalf("TaskCycle() result = %+v, want the late settlement kept and the Run budget-exceeded", result)
	}
	if !result.BudgetDeadline.Equal(startedAt.Add(maximum)) {
		t.Fatalf("late settlement renewed deadline to %s", result.BudgetDeadline)
	}
	if got := taskStatusOnDisk(t, fixture.gitRoot, "task_02"); got != string(spec.StatusPending) {
		t.Fatalf("task_02 status = %q, want pending", got)
	}
}

type settleThenStallRunner struct {
	gitRoot string
	started chan struct{}
	once    sync.Once
}

func (*settleThenStallRunner) Probe(context.Context, agent.ProbeRequest) error { return nil }

func (runner *settleThenStallRunner) Run(ctx context.Context, req agent.ExecuteRequest, _ runevent.Sink) (agent.ExecuteResult, error) {
	taskID := taskIDFromPrompt(req.Prompt)
	if taskID == "task_01" {
		if err := spec.SetStatus(taskPathFromPromptForTest(req.Prompt, runner.gitRoot, taskCycleSlug, taskID), spec.StatusCompleted); err != nil {
			return agent.ExecuteResult{}, err
		}
		return agent.ExecuteResult{}, nil
	}
	runner.once.Do(func() { close(runner.started) })
	<-ctx.Done()
	return agent.ExecuteResult{}, agent.StopError{Err: ctx.Err()}
}

func (*settleThenStallRunner) EndSession(context.Context, agent.RuntimeSpec, agent.SessionRef) error {
	return nil
}

func TestTaskBudgetCancelsAStalledTaskOneAllowanceAfterTheLastSettlement(t *testing.T) {
	fixture := newTaskCycleFixture(t, []taskSpecSeed{
		{id: "task_01"},
		{id: "task_02", needs: []string{"task_01"}},
	})
	runner := &settleThenStallRunner{gitRoot: fixture.gitRoot, started: make(chan struct{})}
	engine := fixture.engine(t, runner, &taskFakeVerifier{calls: fixture.calls}, &engineFakeCommitter{calls: fixture.calls}, fixture.worktree)
	engine.deps.Now = time.Now
	const maximum = 250 * time.Millisecond
	resultCh := make(chan struct {
		result TaskCycleResult
		err    error
	}, 1)
	go func() {
		result, err := engine.TaskCycle(context.Background(), budgetedPlan(fixture.plan(), time.Now(), maximum))
		resultCh <- struct {
			result TaskCycleResult
			err    error
		}{result: result, err: err}
	}()

	testwait.Until(t, "stalled second Task to start", runner.started, resultCh)
	var noEnd <-chan struct{}
	finished := testwait.Until(t, "renewed allowance to cancel the stalled Task", resultCh, noEnd)
	if !errors.Is(finished.err, context.DeadlineExceeded) {
		t.Fatalf("TaskCycle() error = %v, want DeadlineExceeded", finished.err)
	}
	if finished.result.Completed != 1 {
		t.Fatalf("completed Tasks = %d, want the first settlement preserved", finished.result.Completed)
	}
	if got := taskStatusOnDisk(t, fixture.gitRoot, "task_01"); got != string(spec.StatusCompleted) {
		t.Fatalf("task_01 status = %q, want completed", got)
	}
}

type failedSettlementRunner struct {
	clock *budgetTestClock
	at    time.Time
}

func (*failedSettlementRunner) Probe(context.Context, agent.ProbeRequest) error { return nil }

func (runner *failedSettlementRunner) Run(context.Context, agent.ExecuteRequest, runevent.Sink) (agent.ExecuteResult, error) {
	runner.clock.Set(runner.at)
	return agent.ExecuteResult{}, errors.New("scripted Agent failure")
}

func (*failedSettlementRunner) EndSession(context.Context, agent.RuntimeSpec, agent.SessionRef) error {
	return nil
}

func TestTaskBudgetAFailedSettlementRenewsTheAllowance(t *testing.T) {
	fixture := newTaskCycleFixture(t, []taskSpecSeed{{id: "task_01"}})
	startedAt := time.Now()
	const maximum = time.Hour
	clock := &budgetTestClock{now: startedAt}
	engine := fixture.engine(t, &failedSettlementRunner{clock: clock, at: startedAt.Add(54 * time.Minute)}, &taskFakeVerifier{calls: fixture.calls}, &engineFakeCommitter{calls: fixture.calls}, fixture.worktree)
	engine.deps.Now = clock.Now

	result, err := engine.TaskCycle(context.Background(), budgetedPlan(fixture.plan(), startedAt, maximum))

	if err != nil {
		t.Fatalf("TaskCycle() error = %v, want a settled Task failure", err)
	}
	if result.Failed != 1 {
		t.Fatalf("failed Tasks = %d, want 1", result.Failed)
	}
	wantDeadline := startedAt.Add(54*time.Minute + maximum)
	if !result.BudgetDeadline.Equal(wantDeadline) {
		t.Fatalf("BudgetDeadline = %s, want failed settlement renewal %s", result.BudgetDeadline, wantDeadline)
	}
}

func TestTaskCycleReportsTheRenewedBudgetDeadline(t *testing.T) {
	fixture := newTaskCycleFixture(t, []taskSpecSeed{{id: "task_01"}})
	startedAt := time.Now()
	const maximum = time.Hour
	settledAt := startedAt.Add(15 * time.Minute)
	clock := &budgetTestClock{now: startedAt}
	runner := &taskFakeRunner{calls: fixture.calls, gitRoot: fixture.gitRoot, statusByTask: map[string]spec.Status{"task_01": spec.StatusCompleted}, afterTask: func(string) { clock.Set(settledAt) }}
	engine := fixture.engine(t, runner, &taskFakeVerifier{calls: fixture.calls}, &engineFakeCommitter{calls: fixture.calls}, fixture.worktree)
	engine.deps.Now = clock.Now

	result, err := engine.TaskCycle(context.Background(), budgetedPlan(fixture.plan(), startedAt, maximum))

	if err != nil {
		t.Fatalf("TaskCycle() error = %v", err)
	}
	if want := settledAt.Add(maximum); !result.BudgetDeadline.Equal(want) {
		t.Fatalf("BudgetDeadline = %s, want %s", result.BudgetDeadline, want)
	}
}

func TestTaskCycleReportsNoBudgetDeadlineWhenTheBudgetIsDisabled(t *testing.T) {
	fixture := newTaskCycleFixture(t, []taskSpecSeed{{id: "task_01", status: string(spec.StatusCompleted)}})
	engine := fixture.engine(t, &taskFakeRunner{calls: fixture.calls, gitRoot: fixture.gitRoot}, &taskFakeVerifier{calls: fixture.calls}, &engineFakeCommitter{calls: fixture.calls}, fixture.worktree)

	result, err := engine.TaskCycle(context.Background(), fixture.plan())

	if err != nil {
		t.Fatalf("TaskCycle() error = %v", err)
	}
	if !result.BudgetDeadline.IsZero() {
		t.Fatalf("BudgetDeadline = %s, want zero", result.BudgetDeadline)
	}
}

func TestTaskBudgetReasonNamesTheSettlementThatRenewedIt(t *testing.T) {
	fixture := newTaskCycleFixture(t, []taskSpecSeed{
		{id: "task_01"},
		{id: "task_02", needs: []string{"task_01"}},
	})
	runner := &settleThenStallRunner{gitRoot: fixture.gitRoot, started: make(chan struct{})}
	engine := fixture.engine(t, runner, &taskFakeVerifier{calls: fixture.calls}, &engineFakeCommitter{calls: fixture.calls}, fixture.worktree)
	engine.deps.Now = time.Now
	const maximum = 200 * time.Millisecond
	resultCh := make(chan struct {
		result TaskCycleResult
		err    error
	}, 1)
	go func() {
		result, err := engine.TaskCycle(context.Background(), budgetedPlan(fixture.plan(), time.Now(), maximum))
		resultCh <- struct {
			result TaskCycleResult
			err    error
		}{result: result, err: err}
	}()

	testwait.Until(t, "stalled Task to start after the renewing settlement", runner.started, resultCh)
	var noEnd <-chan struct{}
	finished := testwait.Until(t, "renewed allowance to produce its reason", resultCh, noEnd)

	if !errors.Is(finished.err, context.DeadlineExceeded) {
		t.Fatalf("TaskCycle() error = %v, want DeadlineExceeded", finished.err)
	}
	if !strings.HasSuffix(finished.result.TerminalReason, " since Task task_01 settled.") {
		t.Fatalf("TerminalReason = %q, want renewing Task suffix", finished.result.TerminalReason)
	}
}

func budgetWatchdogGoroutines() int {
	stack := make([]byte, 1<<20)
	n := runtime.Stack(stack, true)
	return strings.Count(string(stack[:n]), "(*taskCycleBudget).watch")
}

func TestTaskBudgetWatchdogStopsWhenTheCycleReturns(t *testing.T) {
	before := budgetWatchdogGoroutines()
	fixture := newTaskCycleFixture(t, []taskSpecSeed{{id: "task_01", status: string(spec.StatusCompleted)}})
	engine := fixture.engine(t, &taskFakeRunner{calls: fixture.calls, gitRoot: fixture.gitRoot}, &taskFakeVerifier{calls: fixture.calls}, &engineFakeCommitter{calls: fixture.calls}, fixture.worktree)
	engine.deps.Now = time.Now

	if _, err := engine.TaskCycle(context.Background(), budgetedPlan(fixture.plan(), time.Now(), time.Hour)); err != nil {
		t.Fatalf("TaskCycle() error = %v", err)
	}
	var noEnd <-chan struct{}
	testwait.Poll(t, "TaskCycle watchdog to stop", noEnd, func() (bool, string) {
		got := budgetWatchdogGoroutines()
		return got == before, "watchdog goroutines remain"
	})
}
