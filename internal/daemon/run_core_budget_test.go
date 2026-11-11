package daemon

import (
	"context"
	"path/filepath"
	"strings"
	"testing"

	"roundfix/internal/spec"
)

type runCoreBudgetVerifier struct {
	*taskFakeVerifier
	requests []VerifyRequest
}

func (v *runCoreBudgetVerifier) Verify(ctx context.Context, req VerifyRequest) (VerifyResult, error) {
	v.requests = append(v.requests, req)
	return v.taskFakeVerifier.Verify(ctx, req)
}

func TestEveryRunVerificationAndAgentSessionReceivesTheCoreBudget(t *testing.T) {
	t.Parallel()
	for _, row := range []struct {
		name  string
		cores int
	}{{"budget five", 5}, {"no budget", 0}} {
		t.Run(row.name, func(t *testing.T) {
			t.Parallel()
			fixture := newTaskCycleFixture(t, []taskSpecSeed{{id: "task_01"}})
			runner := &taskFakeRunner{calls: fixture.calls, gitRoot: fixture.gitRoot, statusByTask: map[string]spec.Status{"task_01": spec.StatusCompleted}}
			verifier := &runCoreBudgetVerifier{taskFakeVerifier: &taskFakeVerifier{calls: fixture.calls}}
			engine := fixture.engine(t, runner, verifier, &engineFakeCommitter{calls: fixture.calls}, fixture.worktree)
			engine.deps.VerificationCores = row.cores
			plan := fixture.plan()
			plan.RepositoryVerification = "repository check"
			plan.RepositoryVerificationAtSettlement = true
			result, err := engine.TaskCycle(t.Context(), plan)
			if err != nil || result.Completed != 1 || result.Failed != 0 {
				t.Fatalf("Task cycle = %+v, %v", result, err)
			}
			var probe, settlement bool
			for _, req := range verifier.requests {
				if req.Cores != row.cores {
					t.Errorf("Verification %q (%s) cores = %d, want %d", req.Command, req.OutputPath, req.Cores, row.cores)
				}
				probe = probe || isPreWorkProbeRequest(req)
				settlement = settlement || (req.Command == plan.RepositoryVerification && strings.Contains(filepath.Base(req.OutputPath), "batch-001-attempt-"))
			}
			if !probe || !settlement || len(runner.requests) == 0 {
				t.Fatalf("missing Run boundaries: probe=%t settlement=%t Agent requests=%d", probe, settlement, len(runner.requests))
			}
			for _, req := range runner.requests {
				if req.Cores != row.cores {
					t.Errorf("Agent Session %s cores = %d, want %d", req.Session.Name, req.Cores, row.cores)
				}
			}
		})
	}
}
