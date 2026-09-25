package daemon

import (
	"context"
	"errors"
	"strings"
	"testing"

	"roundfix/internal/spec"
)

func TestTemporaryOnRetryKeepsTheFirstRunDeterministicFailure(t *testing.T) {
	t.Parallel()

	initialDeterministic := &VerificationCommandError{Command: "verify deterministic", OutputPath: "initial-deterministic.log", Err: errors.New("initial failure")}
	initialTemporary := &VerificationCommandError{Command: "verify temporary", OutputPath: "initial-temporary.log", Err: errors.New("exit status 75")}
	initial := verificationAttemptOutcome{
		CommandFailure: initialDeterministic,
		CommandFailures: []verificationAttemptFailure{
			{CommandFailure: initialDeterministic},
			{CommandFailure: initialTemporary},
		},
		TemporaryFailure: &TemporaryVerificationFailureError{CommandFailure: initialTemporary},
	}
	retryCommandFailure := &VerificationCommandError{Command: "verify deterministic", OutputPath: "retry-deterministic.log", Err: errors.New("exit status 75")}
	retryTemporary := &TemporaryVerificationFailureError{CommandFailure: retryCommandFailure}
	retry := verificationAttemptOutcome{
		CommandFailure:   retryCommandFailure,
		CommandFailures:  []verificationAttemptFailure{{CommandFailure: retryCommandFailure}},
		ReachedCommands:  []string{"verify deterministic"},
		TemporaryFailure: retryTemporary,
	}

	got := retainCollectedVerificationFailures(retry, initial)

	if got.TemporaryFailure != retryTemporary {
		t.Fatalf("temporary failure = %+v, want retry temporary failure %+v", got.TemporaryFailure, retryTemporary)
	}
	if len(got.CommandFailures) != 1 || got.CommandFailures[0].CommandFailure != initialDeterministic {
		t.Fatalf("merged failures = %+v, want the first-run deterministic failure", got.CommandFailures)
	}
	if got.CommandFailure != initialDeterministic {
		t.Fatalf("repair target = %+v, want first-run deterministic failure", got.CommandFailure)
	}
	reason := taskVerificationFailureReason(got)
	for _, want := range []string{"verify deterministic", "initial-deterministic.log"} {
		if !strings.Contains(reason, want) {
			t.Fatalf("Task failure reason %q does not contain %q", reason, want)
		}
	}
	if strings.Contains(reason, "retry-deterministic.log") {
		t.Fatalf("Task failure reason names retry temporary diagnostics: %q", reason)
	}
}

func TestPassOnRetryReplacesTheFirstRunFailure(t *testing.T) {
	t.Parallel()

	initialDeterministic := &VerificationCommandError{Command: "verify deterministic", OutputPath: "initial-deterministic.log", Err: errors.New("initial failure")}
	initialTemporary := &VerificationCommandError{Command: "verify temporary", OutputPath: "initial-temporary.log", Err: errors.New("exit status 75")}
	initial := verificationAttemptOutcome{
		CommandFailure: initialDeterministic,
		CommandFailures: []verificationAttemptFailure{
			{CommandFailure: initialDeterministic},
			{CommandFailure: initialTemporary},
		},
		TemporaryFailure: &TemporaryVerificationFailureError{CommandFailure: initialTemporary},
	}
	retryCommandFailure := &VerificationCommandError{Command: "verify temporary", OutputPath: "retry-temporary.log", Err: errors.New("exit status 75")}
	retryTemporary := &TemporaryVerificationFailureError{CommandFailure: retryCommandFailure}
	retry := verificationAttemptOutcome{
		CommandFailure:   retryCommandFailure,
		CommandFailures:  []verificationAttemptFailure{{CommandFailure: retryCommandFailure}},
		ReachedCommands:  []string{"verify deterministic", "verify temporary"},
		TemporaryFailure: retryTemporary,
	}

	got := retainCollectedVerificationFailures(retry, initial)

	if got.TemporaryFailure != retryTemporary {
		t.Fatalf("temporary failure = %+v, want retry temporary failure %+v", got.TemporaryFailure, retryTemporary)
	}
	for _, failure := range got.CommandFailures {
		if failure.CommandFailure != nil && failure.CommandFailure.Command == initialDeterministic.Command {
			t.Fatalf("merged failures retained the first-run verdict for a command that passed on retry: %+v", got.CommandFailures)
		}
	}
	if reason := taskVerificationFailureReason(got); strings.Contains(reason, initialDeterministic.OutputPath) {
		t.Fatalf("Task failure reason retained first-run diagnostics after the command passed: %q", reason)
	}
}

func TestIndependentVerificationTemporaryRetryNamesTheDeterministicFailure(t *testing.T) {
	t.Parallel()

	fixture := newTaskCycleFixture(t, []taskSpecSeed{{
		id:               "task_01",
		verificationMode: spec.VerificationModeIndependent,
		verification:     []string{"verify deterministic", "verify temporary"},
	}})
	runner := &taskFakeRunner{calls: fixture.calls, gitRoot: fixture.gitRoot}
	verifier := &taskFakeVerifier{
		calls:           fixture.calls,
		temporaryOnCall: map[int]bool{2: true, 3: true},
		script:          []error{errors.New("deterministic failure")},
	}
	engine := fixture.engine(t, runner, verifier, &engineFakeCommitter{calls: fixture.calls}, fixture.worktree)

	result, err := engine.TaskCycle(context.Background(), fixture.plan())

	if err != nil {
		t.Fatalf("TaskCycle: %v", err)
	}
	if result.Completed != 0 || result.Failed != 1 || len(result.Outcomes) != 1 {
		t.Fatalf("expected the temporary retry to settle the Task failed, got %+v", result)
	}
	if len(runner.requests) != 1 {
		t.Fatalf("expected no Verification Feedback repair turn, got %d Agent requests", len(runner.requests))
	}
	if got := strings.Join(verifier.commands, "|"); got != "verify deterministic|verify temporary|verify deterministic" {
		t.Fatalf("verification commands = %q, want the first run through temporary failure and one exclusive retry", got)
	}
	if len(verifier.outputPaths) != 3 {
		t.Fatalf("verification diagnostics = %v, want three command diagnostics", verifier.outputPaths)
	}
	reason := result.Outcomes[0].Reason
	for _, want := range []string{"verify deterministic", verifier.outputPaths[0]} {
		if !strings.Contains(reason, want) {
			t.Fatalf("Task failure reason %q does not contain %q", reason, want)
		}
	}
	if strings.Contains(reason, verifier.outputPaths[2]) {
		t.Fatalf("Task failure reason names retry temporary diagnostics: %q", reason)
	}
}
