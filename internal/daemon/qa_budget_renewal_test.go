package daemon

// Suite: QA settlement budget renewal
// Invariant: a timely QA Task settlement renews the Run Budget before the QA Report commit.
// Boundary IN: TaskCycle QA settlement, QA Report commit, and Run Budget result reporting.
// Boundary OUT: CLI post-cycle integration, push, and cleanup.

import (
	"context"
	"strings"
	"testing"
	"time"

	"roundfix/internal/spec"
)

func TestQASettlementRenewsTheBudgetBeforeTheReportCommit(t *testing.T) {
	fixture := newTaskCycleFixture(t, []taskSpecSeed{{id: "task_01"}})
	startedAt := time.Now()
	const maximum = time.Hour
	clock := &budgetTestClock{now: startedAt}
	runner := &taskFakeRunner{
		calls:        fixture.calls,
		gitRoot:      fixture.gitRoot,
		statusByTask: map[string]spec.Status{"task_01": spec.StatusCompleted},
		qaReport:     qaReportForTest(spec.VerdictPass),
		afterTask:    func(string) { clock.Set(startedAt.Add(50 * time.Minute)) },
		afterQA:      func() { clock.Set(startedAt.Add(109 * time.Minute)) },
	}
	qaCommitCompleted := false
	committer := &engineFakeCommitter{calls: fixture.calls, afterCommit: func(_ context.Context, request CommitRequest) error {
		if !strings.HasPrefix(request.Message, "docs: qa report") {
			return nil
		}
		clock.Set(startedAt.Add(120 * time.Minute))
		qaCommitCompleted = true
		return nil
	}}
	engine := fixture.engine(t, runner, &taskFakeVerifier{calls: fixture.calls}, committer, fixture.worktree)
	engine.deps.Now = clock.Now
	engine.deps.MechanicalStage = &fakeQAMechanicalStage{}

	result, err := engine.TaskCycle(context.Background(), budgetedPlan(fixture.qaPlan(), startedAt, maximum))

	if err != nil {
		t.Fatalf("TaskCycle() error = %v, want the QA Report commit to finish under the renewed allowance", err)
	}
	if !qaCommitCompleted {
		t.Fatal("QA Report commit did not complete")
	}
	if result.QAVerdict != spec.VerdictPass || !result.QAAccepted {
		t.Fatalf("QA result = verdict %q accepted=%t, want accepted pass", result.QAVerdict, result.QAAccepted)
	}
	wantDeadline := startedAt.Add(109*time.Minute + maximum)
	if !result.BudgetDeadline.Equal(wantDeadline) || result.BudgetRenewedBy != "task_02" {
		t.Fatalf("QA settlement budget = deadline %s renewed by %q, want deadline %s renewed by task_02", result.BudgetDeadline, result.BudgetRenewedBy, wantDeadline)
	}
}
