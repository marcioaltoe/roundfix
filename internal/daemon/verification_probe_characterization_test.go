package daemon

import (
	"context"
	"errors"
	"strings"
	"testing"

	"roundfix/internal/spec"
)

// TestVerificationProbeCharacterizationVacuousGateSettlesFailedBeforeAgent
// records the refusal behavior that replaced the declared break from task_01.
func TestVerificationProbeCharacterizationVacuousGateSettlesFailedBeforeAgent(t *testing.T) {
	t.Parallel()
	const command = "gate already passes"
	fixture := newTaskCycleFixture(t, []taskSpecSeed{{
		id:           "task_01",
		verification: []string{command},
	}})
	runner := &taskFakeRunner{calls: fixture.calls, gitRoot: fixture.gitRoot}
	verifier := &characterizationVerdictVerifier{calls: fixture.calls}
	committer := &engineFakeCommitter{calls: fixture.calls}
	engine := fixture.engine(t, runner, verifier, committer, fixture.worktree)

	result, err := engine.TaskCycle(context.Background(), fixture.plan())

	if err != nil {
		t.Fatalf("TaskCycle() error = %v", err)
	}
	if result.Completed != 0 || result.Failed != 1 || result.Skipped != 0 {
		t.Fatalf("vacuous Verification settlement = %+v, want one failed Task", result)
	}
	if got := taskStatusOnDisk(t, fixture.gitRoot, "task_01"); got != string(spec.StatusFailed) {
		t.Fatalf("vacuous Verification status = %q, want %q", got, spec.StatusFailed)
	}
	if got := strings.Join(*fixture.calls, ">"); got != "verify" {
		t.Fatalf("vacuous Verification call sequence = %q, want only the pre-work probe", got)
	}
	if got := strings.Join(verifier.commands, "|"); got != command {
		t.Fatalf("fake verifier commands = %q, want %q", got, command)
	}
	if runner.taskCalls["task_01"] != 0 || len(committer.messages) != 0 {
		t.Fatalf("vacuous Verification spent Agent or commit work: agent=%d commits=%v", runner.taskCalls["task_01"], committer.messages)
	}
}

// TestVerificationProbeCharacterizationUnobservedVerdictSettlesFailed records
// that an unobserved verdict settles failed with an explicit unknown cause,
// distinct from a command rejection.
func TestVerificationProbeCharacterizationUnobservedVerdictSettlesFailed(t *testing.T) {
	t.Parallel()
	const (
		unobservedCommand = "gate verdict not observed"
		failedCommand     = "gate ran and rejected the work"
	)
	unobservedErr := errors.New("runner did not observe the command verdict")
	fixture := newTaskCycleFixture(t, []taskSpecSeed{
		{id: "task_01", verification: []string{unobservedCommand}},
		{id: "task_02", verification: []string{failedCommand}},
	})
	runner := &taskFakeRunner{calls: fixture.calls, gitRoot: fixture.gitRoot}
	verifier := &characterizationVerdictVerifier{
		calls:     fixture.calls,
		unknownOn: map[string]error{unobservedCommand: unobservedErr},
		failOn:    map[string]error{failedCommand: errors.New("exit status 7")},
	}
	committer := &engineFakeCommitter{calls: fixture.calls}
	engine := fixture.engine(t, runner, verifier, committer, fixture.worktree)

	result, err := engine.TaskCycle(context.Background(), fixture.plan())

	if err != nil {
		t.Fatalf("TaskCycle() error = %v", err)
	}
	if result.Completed != 0 || result.Failed != 2 || result.Skipped != 0 {
		t.Fatalf("Verification settlements = %+v, want both Tasks failed", result)
	}
	unobserved := failedTaskOutcome(t, result.Outcomes, "task_01")
	observed := failedTaskOutcome(t, result.Outcomes, "task_02")
	unobservedPath := verificationProbeOutputPath(fixture.artifactDir, fixture.run.ID, 1, 1)
	for _, want := range []string{"Verification unknown", unobservedCommand, unobservedErr.Error(), unobservedPath} {
		if !strings.Contains(unobserved.Reason, want) {
			t.Fatalf("unobserved terminal reason %q does not contain %q", unobserved.Reason, want)
		}
	}
	observedPath := VerificationOutputPath(fixture.artifactDir, fixture.run.ID, 2, 2)
	for _, want := range []string{"Verification failed", failedCommand, "exit status 7", observedPath} {
		if !strings.Contains(observed.Reason, want) {
			t.Fatalf("command-failure terminal reason %q does not contain %q", observed.Reason, want)
		}
	}
	if unobserved.Reason == observed.Reason {
		t.Fatalf("unobserved and command-failure terminal reasons are identical: %q", unobserved.Reason)
	}
	if runner.taskCalls["task_01"] != 0 {
		t.Fatalf("unobserved probe spent %d Agent turns, want 0", runner.taskCalls["task_01"])
	}
	if len(committer.messages) != 0 {
		t.Fatalf("failed Tasks produced commits: %v", committer.messages)
	}
}

type characterizationVerdictVerifier struct {
	calls     *[]string
	unknownOn map[string]error
	failOn    map[string]error
	commands  []string
}

func (verifier *characterizationVerdictVerifier) Verify(_ context.Context, req VerifyRequest) (VerifyResult, error) {
	*verifier.calls = append(*verifier.calls, "verify")
	verifier.commands = append(verifier.commands, req.Command)
	result := VerifyResult{OutputPath: req.OutputPath}
	if err := verifier.unknownOn[req.Command]; err != nil {
		return result, &VerificationUnknownError{Command: req.Command, DiagnosticPath: req.OutputPath, Err: err}
	}
	if err := verifier.failOn[req.Command]; err != nil {
		return result, &VerificationCommandError{Command: req.Command, OutputPath: req.OutputPath, Err: err}
	}
	return result, nil
}

func failedTaskOutcome(t *testing.T, outcomes []TaskOutcome, taskID string) TaskOutcome {
	t.Helper()
	outcome, ok := taskOutcomeByID(outcomes, taskID)
	if !ok {
		t.Fatalf("missing outcome for %s: %+v", taskID, outcomes)
	}
	if outcome.Status != string(spec.StatusFailed) {
		t.Fatalf("%s status = %q, want %q", taskID, outcome.Status, spec.StatusFailed)
	}
	return outcome
}
