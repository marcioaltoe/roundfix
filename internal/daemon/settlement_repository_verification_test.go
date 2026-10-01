package daemon

import (
	"context"
	"errors"
	"reflect"
	"strings"
	"testing"

	"roundfix/internal/spec"
)

const settlementRepositoryCommand = "repository check"

func settlementRepositoryFixture(t *testing.T, declared []string) (*taskCycleFixture, TaskPlan, *taskFakeRunner, *taskFakeVerifier, *engineFakeCommitter) {
	t.Helper()
	fixture := newTaskCycleFixture(t, []taskSpecSeed{{id: "task_01", verification: declared}})
	plan := fixture.qaPlan()
	// Isolate non-QA settlement from execution of the separate terminal gate.
	for i := range plan.Tasks {
		if plan.Tasks[i].Type == spec.TaskTypeQA {
			plan.Tasks[i].Status = spec.StatusCompleted
		}
	}
	plan.RepositoryVerification = settlementRepositoryCommand
	plan.RepositoryVerificationAtSettlement = true
	fixture.worktree.snapshots = [][]string{nil, {"src/one.go"}}
	runner := &taskFakeRunner{calls: fixture.calls, gitRoot: fixture.gitRoot, statusByTask: map[string]spec.Status{"task_01": spec.StatusCompleted}}
	verifier := &taskFakeVerifier{calls: fixture.calls}
	committer := &engineFakeCommitter{calls: fixture.calls}
	return fixture, plan, runner, verifier, committer
}

func TestGatedGraphRunsRepositoryVerificationAfterDeclaredCommands(t *testing.T) {
	t.Parallel()
	fixture, plan, runner, verifier, committer := settlementRepositoryFixture(t, []string{"declared check"})
	// A repository failure forces the second attempt, exposing its full order.
	verifier.script = []error{nil, errors.New("exit status 1"), nil, nil}
	result, err := fixture.engine(t, runner, verifier, committer, fixture.worktree).TaskCycle(context.Background(), plan)
	if err != nil || result.Completed != 1 {
		t.Fatalf("TaskCycle = %+v, %v", result, err)
	}
	want := []string{"declared check", settlementRepositoryCommand, "declared check", settlementRepositoryCommand}
	if !reflect.DeepEqual(verifier.commands, want) {
		t.Fatalf("commands = %v, want %v", verifier.commands, want)
	}
}

func TestRepositoryVerificationFailureAtSettlementReturnsFeedback(t *testing.T) {
	t.Parallel()
	fixture, plan, runner, verifier, committer := settlementRepositoryFixture(t, []string{"declared check"})
	verifier.script = []error{nil, errors.New("exit status 1"), nil, nil}
	result, err := fixture.engine(t, runner, verifier, committer, fixture.worktree).TaskCycle(context.Background(), plan)
	if err != nil || result.Completed != 1 || result.Failed != 0 {
		t.Fatalf("TaskCycle = %+v, %v", result, err)
	}
	if len(runner.requests) != 2 {
		t.Fatalf("Agent requests = %d, want 2", len(runner.requests))
	}
	if prompt := runner.requests[1].Prompt; !strings.Contains(prompt, "Verification Feedback") || !strings.Contains(prompt, "Failed command: "+settlementRepositoryCommand) {
		t.Fatalf("repair prompt omitted repository failure: %s", prompt)
	}
	if runner.requests[0].Session != runner.requests[1].Session {
		t.Fatal("repair changed Agent Session")
	}
	if got := taskStatusOnDisk(t, fixture.gitRoot, "task_01"); got != "completed" {
		t.Fatalf("status = %q", got)
	}
	if len(committer.messages) != 1 {
		t.Fatalf("commits = %v, want one repaired Task commit", committer.messages)
	}
}

func TestRepositoryVerificationFailureOnFinalAttemptFailsTheTask(t *testing.T) {
	t.Parallel()
	fixture, plan, runner, verifier, committer := settlementRepositoryFixture(t, []string{"declared check"})
	verifier.failOn = map[string]error{settlementRepositoryCommand: errors.New("exit status 7")}
	result, err := fixture.engine(t, runner, verifier, committer, fixture.worktree).TaskCycle(context.Background(), plan)
	if err != nil || result.Failed != 1 || result.Completed != 0 {
		t.Fatalf("TaskCycle = %+v, %v", result, err)
	}
	outcome, ok := taskOutcomeByID(result.Outcomes, "task_01")
	if !ok || !strings.Contains(outcome.Reason, settlementRepositoryCommand) {
		t.Fatalf("failure outcome = %+v", outcome)
	}
	if got := taskStatusOnDisk(t, fixture.gitRoot, "task_01"); got != "failed" {
		t.Fatalf("status = %q", got)
	}
	if len(runner.requests) != 2 || len(committer.messages) != 0 {
		t.Fatalf("requests = %d, commits = %v", len(runner.requests), committer.messages)
	}
}

func TestTaskDeclaringRepositoryVerificationRunsItOnce(t *testing.T) {
	t.Parallel()
	fixture, plan, runner, verifier, committer := settlementRepositoryFixture(t, []string{"declared check", settlementRepositoryCommand})
	// Entry precondition, then both attempts; the probe is tracked separately.
	verifier.script = []error{nil, nil, errors.New("exit status 1"), nil, nil}
	result, err := fixture.engine(t, runner, verifier, committer, fixture.worktree).TaskCycle(context.Background(), plan)
	if err != nil || result.Completed != 1 {
		t.Fatalf("TaskCycle = %+v, %v", result, err)
	}
	want := []string{settlementRepositoryCommand, "declared check", settlementRepositoryCommand, "declared check", settlementRepositoryCommand}
	if !reflect.DeepEqual(verifier.commands, want) {
		t.Fatalf("commands = %v, want %v (entry then one per attempt)", verifier.commands, want)
	}
}

func TestRepositoryAtSettlementOffRunsOnlyDeclaredCommands(t *testing.T) {
	t.Parallel()
	fixture, plan, runner, verifier, committer := settlementRepositoryFixture(t, []string{"declared check"})
	plan.RepositoryVerificationAtSettlement = false
	result, err := fixture.engine(t, runner, verifier, committer, fixture.worktree).TaskCycle(context.Background(), plan)
	if err != nil || result.Completed != 1 {
		t.Fatalf("TaskCycle = %+v, %v", result, err)
	}
	if !reflect.DeepEqual(verifier.commands, []string{"declared check"}) {
		t.Fatalf("commands = %v", verifier.commands)
	}
}

// Record probe requests while preserving the existing verifier's semantics.
type settlementProbeVerifier struct {
	*taskFakeVerifier
	probes []string
}

func (v *settlementProbeVerifier) Verify(ctx context.Context, req VerifyRequest) (VerifyResult, error) {
	if isPreWorkProbeRequest(req) {
		v.probes = append(v.probes, req.Command)
	}
	return v.taskFakeVerifier.Verify(ctx, req)
}

func TestPreWorkProbeNeverRunsTheSettlementRepositoryVerification(t *testing.T) {
	t.Parallel()
	fixture, plan, runner, verifier, committer := settlementRepositoryFixture(t, []string{"declared check"})
	recording := &settlementProbeVerifier{taskFakeVerifier: verifier}
	result, err := fixture.engine(t, runner, recording, committer, fixture.worktree).TaskCycle(context.Background(), plan)
	if err != nil || result.Completed != 1 {
		t.Fatalf("TaskCycle = %+v, %v", result, err)
	}
	if !reflect.DeepEqual(recording.probes, []string{"declared check"}) {
		t.Fatalf("probe commands = %v", recording.probes)
	}
	if !reflect.DeepEqual(verifier.commands, []string{"declared check", settlementRepositoryCommand}) {
		t.Fatalf("settlement commands = %v", verifier.commands)
	}
	if !reflect.DeepEqual(plan.Tasks[0].Verification, []string{"declared check"}) {
		t.Fatalf("authored Task was mutated: %v", plan.Tasks[0].Verification)
	}
}

func TestEmptySettlementRepositoryCommandAppendsNothing(t *testing.T) {
	t.Parallel()
	fixture, plan, runner, verifier, committer := settlementRepositoryFixture(t, []string{"declared check"})
	plan.RepositoryVerification = "  "
	result, err := fixture.engine(t, runner, verifier, committer, fixture.worktree).TaskCycle(context.Background(), plan)
	if err != nil || result.Completed != 1 {
		t.Fatalf("TaskCycle = %+v, %v", result, err)
	}
	if !reflect.DeepEqual(verifier.commands, []string{"declared check"}) {
		t.Fatalf("commands = %v", verifier.commands)
	}
}

func TestUngatedGraphIgnoresEnabledRepositorySettlementVerification(t *testing.T) {
	t.Parallel()
	fixture := newTaskCycleFixture(t, []taskSpecSeed{{id: "task_01", verification: []string{"declared check"}}})
	plan := fixture.plan()
	plan.RepositoryVerification = settlementRepositoryCommand
	plan.RepositoryVerificationAtSettlement = true
	fixture.worktree.snapshots = [][]string{nil, {"src/one.go"}}
	runner := &taskFakeRunner{calls: fixture.calls, gitRoot: fixture.gitRoot, statusByTask: map[string]spec.Status{"task_01": spec.StatusCompleted}}
	verifier := &taskFakeVerifier{calls: fixture.calls}
	committer := &engineFakeCommitter{calls: fixture.calls}
	result, err := fixture.engine(t, runner, verifier, committer, fixture.worktree).TaskCycle(context.Background(), plan)
	if err != nil || result.Completed != 1 {
		t.Fatalf("TaskCycle = %+v, %v", result, err)
	}
	if !reflect.DeepEqual(verifier.commands, []string{"declared check"}) {
		t.Fatalf("commands = %v", verifier.commands)
	}
}
