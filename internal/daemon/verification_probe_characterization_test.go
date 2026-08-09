package daemon

import (
	"context"
	"errors"
	"strings"
	"testing"

	"roundfix/internal/spec"
)

// TestVerificationProbeCharacterizationVacuousGateSettlesCompleted records
// that a gate which already passes still lets the Daemon settle the Task.
// Declared break: task_03 changes this to refusal before the Agent turn.
func TestVerificationProbeCharacterizationVacuousGateSettlesCompleted(t *testing.T) {
	t.Parallel()
	const command = "gate already passes"
	fixture := newTaskCycleFixture(t, []taskSpecSeed{{
		id:           "task_01",
		verification: []string{command},
	}})
	runner := &taskFakeRunner{calls: fixture.calls, gitRoot: fixture.gitRoot}
	verifier := &taskFakeVerifier{calls: fixture.calls}
	committer := &engineFakeCommitter{calls: fixture.calls}
	engine := fixture.engine(t, runner, verifier, committer, fixture.worktree)

	result, err := engine.TaskCycle(context.Background(), fixture.plan())

	if err != nil {
		t.Fatalf("TaskCycle() error = %v", err)
	}
	if result.Completed != 1 || result.Failed != 0 || result.Skipped != 0 {
		t.Fatalf("vacuous Verification settlement = %+v, want one completed Task", result)
	}
	if got := taskStatusOnDisk(t, fixture.gitRoot, "task_01"); got != string(spec.StatusCompleted) {
		t.Fatalf("vacuous Verification status = %q, want %q", got, spec.StatusCompleted)
	}
	if got := strings.Join(*fixture.calls, ">"); got != "agent>verify>commit" {
		t.Fatalf("vacuous Verification call sequence = %q, want Agent, Verification, and commit", got)
	}
	if got := strings.Join(verifier.commands, "|"); got != command {
		t.Fatalf("fake verifier commands = %q, want %q", got, command)
	}
}

// TestVerificationProbeCharacterizationUnobservedVerdictSettlesFailed records
// that an unobserved verdict has the same failure report as a command rejection.
// Declared break: task_02 changes the unobserved case to an explicit unknown cause.
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
	verifier := &taskFakeVerifier{
		calls: fixture.calls,
		failOn: map[string]error{
			unobservedCommand: unobservedErr,
			failedCommand:     errors.New("exit status 7"),
		},
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
	unobservedShape := currentVerificationFailureShape(
		unobserved.Reason,
		unobservedCommand,
		unobservedErr.Error(),
		VerificationOutputPath(fixture.artifactDir, fixture.run.ID, 1, 2),
	)
	observedShape := currentVerificationFailureShape(
		observed.Reason,
		failedCommand,
		"exit status 7",
		VerificationOutputPath(fixture.artifactDir, fixture.run.ID, 2, 2),
	)
	if unobservedShape != observedShape {
		t.Fatalf("unobserved failure shape = %q, command failure shape = %q", unobservedShape, observedShape)
	}
	const currentShape = `Verification failed: command "<command>" exited with <cause>; diagnostics: <diagnostics>`
	if unobservedShape != currentShape {
		t.Fatalf("normalized current failure shape = %q, want %q", unobservedShape, currentShape)
	}
	if strings.Contains(strings.ToLower(unobserved.Reason), "unknown") {
		t.Fatalf("unobserved verdict unexpectedly distinguished as unknown: %q", unobserved.Reason)
	}
	if len(committer.messages) != 0 {
		t.Fatalf("failed Tasks produced commits: %v", committer.messages)
	}
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

func currentVerificationFailureShape(reason string, command string, cause string, diagnosticPath string) string {
	return strings.NewReplacer(
		command, "<command>",
		cause, "<cause>",
		diagnosticPath, "<diagnostics>",
	).Replace(reason)
}
