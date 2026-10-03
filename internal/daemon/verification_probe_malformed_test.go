package daemon

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestProbeReportsACommandTheShellCannotParse(t *testing.T) {
	t.Parallel()
	commands := []string{`grep -qF '| \`, "false"}
	verdicts, err := ProbeCommands(context.Background(), ExecVerifier{}, t.TempDir(), commands, func(i int) string { return filepath.Join(t.TempDir(), fmt.Sprintf("%d.log", i)) })
	if err != nil {
		t.Fatal(err)
	}
	if len(verdicts) != 2 || !verdicts[0].Unknown || verdicts[0].Vacuous || !errors.Is(verdicts[0].Cause, ErrVerificationMalformed) {
		t.Fatalf("verdicts = %+v", verdicts)
	}
	var unknown *VerificationUnknownError
	if !errors.As(verdicts[0].Cause, &unknown) || !strings.Contains(strings.ToLower(unknown.Err.Error()), "syntax error") {
		t.Fatalf("cause = %v", verdicts[0].Cause)
	}
	if verdicts[1].Unknown || verdicts[1].Vacuous {
		t.Fatalf("parsable failure = %+v", verdicts[1])
	}
}

func TestProbeRunsNoMalformedCommand(t *testing.T) {
	t.Parallel()
	marker := filepath.Join(t.TempDir(), "ran")
	// A streaming shell can execute the first line before rejecting the second.
	command := fmt.Sprintf("touch %q\nif", marker)
	verdicts, err := ProbeCommands(context.Background(), ExecVerifier{}, t.TempDir(), []string{command}, func(int) string { return filepath.Join(t.TempDir(), "probe.log") })
	if err != nil {
		t.Fatal(err)
	}
	if len(verdicts) != 1 || !errors.Is(verdicts[0].Cause, ErrVerificationMalformed) {
		t.Fatalf("verdicts = %+v", verdicts)
	}
	if _, err := os.Stat(marker); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("malformed command ran: %v", err)
	}
}

func TestPreWorkProbeRefusesAMalformedCommand(t *testing.T) {
	t.Parallel()
	fixture := newTaskCycleFixture(t, []taskSpecSeed{{id: "task_01", verification: []string{"if"}}})
	runner := &taskFakeRunner{calls: fixture.calls, gitRoot: fixture.gitRoot}
	verifier := &probeScriptVerifier{}
	committer := &engineFakeCommitter{calls: fixture.calls}
	engine := fixture.engine(t, runner, verifier, committer, fixture.worktree)
	result, err := engine.TaskCycle(context.Background(), fixture.plan())
	if err != nil {
		t.Fatal(err)
	}
	outcome := failedTaskOutcome(t, result.Outcomes, "task_01")
	if !strings.Contains(outcome.Reason, ErrVerificationMalformed.Error()) {
		t.Fatalf("reason = %q", outcome.Reason)
	}
	if runner.taskCalls["task_01"] != 0 || len(verifier.requests) != 0 || len(committer.messages) != 0 {
		t.Fatalf("malformed task spent work: runner=%v verifier=%v commits=%v", runner.taskCalls, verifier.requests, committer.messages)
	}
}

func TestProbeParserUnavailableRunsVerifier(t *testing.T) {
	t.Setenv("PATH", t.TempDir())
	verifier := &probeScriptVerifier{}
	verdicts, err := ProbeCommands(context.Background(), verifier, t.TempDir(), []string{"if"}, func(int) string { return "unused" })
	if err != nil || len(verifier.requests) != 1 || len(verdicts) != 1 || !verdicts[0].Vacuous {
		t.Fatalf("fallback: verdicts=%+v requests=%+v err=%v", verdicts, verifier.requests, err)
	}
}

func TestProbeCancelledParseDoesNotRunVerifier(t *testing.T) {
	t.Parallel()
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	verifier := &probeScriptVerifier{}
	_, err := ProbeCommands(ctx, verifier, t.TempDir(), []string{"false"}, func(int) string { return "unused" })
	if !errors.Is(err, context.Canceled) || len(verifier.requests) != 0 {
		t.Fatalf("cancel: requests=%+v err=%v", verifier.requests, err)
	}
}
