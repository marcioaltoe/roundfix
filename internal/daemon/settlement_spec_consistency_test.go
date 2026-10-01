package daemon

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"roundfix/internal/speccheck"
)

type specConsistencyChecker struct {
	read  func(int) ([]speccheck.Finding, error)
	calls int
}

func (c *specConsistencyChecker) RefusingFindings(string, string, string) ([]speccheck.Finding, error) {
	c.calls++
	return c.read(c.calls)
}
func (*specConsistencyChecker) AuditCommit(context.Context, speccheck.MechanicalRequest, speccheck.ProspectiveTaskCommit) ([]speccheck.MechanicalFinding, error) {
	return nil, nil
}

func specConsistencyFixture(t *testing.T, checker *specConsistencyChecker) (*taskCycleFixture, TaskPlan, *taskFakeRunner, *engineFakeCommitter, *Engine) {
	t.Helper()
	fixture, plan, runner, verifier, committer := settlementRepositoryFixture(t, []string{"declared check"})
	engine := fixture.engine(t, runner, verifier, committer, fixture.worktree)
	engine.deps.SettlementChecker = checker
	return fixture, plan, runner, committer, engine
}

func taskConsistencyFinding() speccheck.Finding {
	return speccheck.Finding{Code: "SC-TEST", Summary: "A new contradiction.", Where: []speccheck.Location{{Path: "_prd.md", Line: 7}, {Path: "_techspec.md", Line: 9}}, Fix: "Repair the contradiction."}
}

func TestSettlementRefusesASpecConsistencyFindingTheTaskIntroduced(t *testing.T) {
	t.Parallel()
	finding := taskConsistencyFinding()
	checker := &specConsistencyChecker{read: func(call int) ([]speccheck.Finding, error) {
		if call == 1 {
			return nil, nil
		}
		return []speccheck.Finding{finding}, nil
	}}
	fixture, plan, runner, committer, engine := specConsistencyFixture(t, checker)
	result, err := engine.TaskCycle(context.Background(), plan)
	if err != nil || result.Failed != 1 || result.Completed != 0 || checker.calls != 3 || len(committer.messages) != 0 {
		t.Fatalf("result %+v, err %v, reads %d, commits %v", result, err, checker.calls, committer.messages)
	}
	if len(runner.requests) != 2 || runner.requests[0].Session != runner.requests[1].Session {
		t.Fatalf("requests = %+v", runner.requests)
	}
	for _, text := range []string{"Verification Feedback", "settlement check: spec consistency", finding.Code, finding.Summary, "_prd.md:7", "_techspec.md:9", "Fix: " + finding.Fix} {
		if !strings.Contains(runner.requests[1].Prompt, text) {
			t.Errorf("Feedback missing %q: %s", text, runner.requests[1].Prompt)
		}
	}
	if got := taskStatusOnDisk(t, fixture.gitRoot, "task_01"); got != "failed" {
		t.Fatalf("status = %q", got)
	}
	outcome, _ := taskOutcomeByID(result.Outcomes, "task_01")
	if !strings.Contains(outcome.Reason, finding.Code) {
		t.Fatalf("reason = %q", outcome.Reason)
	}
	diagnostic, err := os.ReadFile(filepath.Join(plan.ArtifactDir, "runs", plan.RunID, "verification", "batch-001-attempt-2-settlement-spec-consistency.log"))
	if err != nil || !strings.Contains(string(diagnostic), finding.Code) {
		t.Fatalf("diagnostic %q, err %v", diagnostic, err)
	}
}

func TestSettlementIgnoresASpecConsistencyFindingPresentAtStart(t *testing.T) {
	t.Parallel()
	checker := &specConsistencyChecker{read: func(call int) ([]speccheck.Finding, error) {
		finding := taskConsistencyFinding()
		if call > 1 {
			finding.Where = []speccheck.Location{{Path: "moved.md", Line: 20}}
			finding.Fix = "Updated repair."
		}
		return []speccheck.Finding{finding}, nil
	}}
	_, plan, runner, committer, engine := specConsistencyFixture(t, checker)
	result, err := engine.TaskCycle(context.Background(), plan)
	if err != nil || result.Completed != 1 || result.Failed != 0 || checker.calls != 2 || len(runner.requests) != 1 || len(committer.messages) != 1 {
		t.Fatalf("result %+v, err %v, reads %d, requests %d, commits %v", result, err, checker.calls, len(runner.requests), committer.messages)
	}
}

func TestSpecConsistencyBaselineErrorFailsTheTaskBeforeAgentWork(t *testing.T) {
	t.Parallel()
	checker := &specConsistencyChecker{read: func(int) ([]speccheck.Finding, error) { return nil, errors.New("cannot read baseline") }}
	fixture, plan, runner, committer, engine := specConsistencyFixture(t, checker)
	result, err := engine.TaskCycle(context.Background(), plan)
	if err != nil || result.Failed != 1 || result.Completed != 0 || checker.calls != 1 || len(runner.requests) != 0 || len(committer.messages) != 0 {
		t.Fatalf("result %+v, err %v, reads %d, requests %d, commits %v", result, err, checker.calls, len(runner.requests), committer.messages)
	}
	outcome, _ := taskOutcomeByID(result.Outcomes, "task_01")
	if !strings.Contains(outcome.Reason, "cannot read baseline") {
		t.Fatalf("reason = %q", outcome.Reason)
	}
	if got := taskStatusOnDisk(t, fixture.gitRoot, "task_01"); got != "failed" {
		t.Fatalf("status = %q", got)
	}
}

func TestSpecConsistencyCheckerErrorFailsTheCheck(t *testing.T) {
	t.Parallel()
	checker := &specConsistencyChecker{read: func(call int) ([]speccheck.Finding, error) {
		if call == 1 {
			return nil, nil
		}
		return nil, errors.New("cannot read attempt findings")
	}}
	_, plan, runner, committer, engine := specConsistencyFixture(t, checker)
	result, err := engine.TaskCycle(context.Background(), plan)
	if err != nil || result.Failed != 1 || checker.calls != 3 || len(committer.messages) != 0 {
		t.Fatalf("result %+v, err %v, reads %d", result, err, checker.calls)
	}
	if len(runner.requests) != 2 || !strings.Contains(runner.requests[1].Prompt, "cannot read attempt findings") || !strings.Contains(runner.requests[1].Prompt, "settlement check: spec consistency") {
		t.Fatalf("requests = %+v", runner.requests)
	}
}

func TestSpecConsistencyRepairKeepsTheTaskStartBaseline(t *testing.T) {
	t.Parallel()
	checker := &specConsistencyChecker{read: func(call int) ([]speccheck.Finding, error) {
		if call == 2 {
			return []speccheck.Finding{taskConsistencyFinding()}, nil
		}
		return nil, nil
	}}
	_, plan, runner, committer, engine := specConsistencyFixture(t, checker)
	result, err := engine.TaskCycle(context.Background(), plan)
	if err != nil || result.Completed != 1 || checker.calls != 3 || len(runner.requests) != 2 || len(committer.messages) != 1 {
		t.Fatalf("result %+v, err %v, reads %d", result, err, checker.calls)
	}
}
