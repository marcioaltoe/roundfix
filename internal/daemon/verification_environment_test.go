package daemon

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"roundfix/internal/judge"
)

const verificationEnvironmentFixture = "ROUNDFIX_VERIFICATION_ENVIRONMENT_FIXTURE"

func TestVerificationReceivesNoJudgeKeyVariable(t *testing.T) {
	t.Parallel()
	// Sequential: the isolated child sets the Judge Key Variables and ROUNDFIX_VERIFICATION_ENVIRONMENT_CONTROL before any tests start.
	command := exec.CommandContext(t.Context(), os.Args[0])
	command.Env = append(judge.WithoutKeyVariables(os.Environ()), verificationEnvironmentFixture+"=1")
	if output, err := command.CombinedOutput(); err != nil {
		t.Fatalf("verification environment fixture: %v\n%s", err, output)
	}
}

// runVerificationEnvironmentFixture runs before the child's test harness,
// so its process environment changes cannot race any other test.
func runVerificationEnvironmentFixture() (bool, error) {
	if os.Getenv(verificationEnvironmentFixture) != "1" {
		return false, nil
	}
	questions, err := judge.Load()
	if err != nil {
		return true, err
	}
	variables := append(questions.KeyVariables(), "TYPESAFE_API_KEY")
	var checks []string
	for index, variable := range variables {
		if err := os.Setenv(variable, fmt.Sprintf("verification-fixture-%d", index)); err != nil {
			return true, fmt.Errorf("set judge fixture variable %s: %w", variable, err)
		}
		checks = append(checks, fmt.Sprintf(`[ "${%s+x}" = "" ]`, variable))
	}
	const runID = "run-verification-fixture"
	if err := os.Setenv(judge.RunMarkerVariable, "inherited-run"); err != nil {
		return true, err
	}
	checks = append(checks, fmt.Sprintf(`[ "$%s" = "%s" ]`, judge.RunMarkerVariable, runID))
	const control = "ROUNDFIX_VERIFICATION_ENVIRONMENT_CONTROL"
	value := fmt.Sprintf("control-%d", len(variables))
	if err := os.Setenv(control, value); err != nil {
		return true, fmt.Errorf("set control variable: %w", err)
	}
	checks = append(checks, fmt.Sprintf(`[ "$%s" = "%s" ]`, control, value))
	dir, err := os.MkdirTemp("", "roundfix-verification-environment-")
	if err != nil {
		return true, err
	}
	_, verifyErr := (ExecVerifier{}).Verify(context.Background(), VerifyRequest{
		RunID:      runID,
		WorkDir:    dir,
		Command:    strings.Join(checks, " && "),
		OutputPath: filepath.Join(dir, "verification.log"),
	})
	cleanupErr := os.RemoveAll(dir)
	if err := errors.Join(verifyErr, cleanupErr); err != nil {
		return true, err
	}
	if os.Getenv(judge.RunMarkerVariable) != "inherited-run" {
		return true, fmt.Errorf("fixture parent Run Marker changed")
	}
	for index, variable := range variables {
		if os.Getenv(variable) != fmt.Sprintf("verification-fixture-%d", index) {
			return true, fmt.Errorf("fixture parent variable %s changed", variable)
		}
	}
	return true, nil
}

func TestEnginePassesRunIDToEveryVerificationCommand(t *testing.T) {
	t.Parallel()
	fixture, plan, runner, verifier, committer := settlementRepositoryFixture(t, []string{"declared check"})
	recording := &qaGateRecordingVerifier{delegate: verifier, onVerify: func(req VerifyRequest) {
		if req.RunID != plan.RunID {
			t.Errorf("Verification %q in %q got RunID %q, want %q", req.Command, req.OutputPath, req.RunID, plan.RunID)
		}
	}}
	result, err := fixture.engine(t, runner, recording, committer, fixture.worktree).TaskCycle(t.Context(), plan)
	if err != nil || result.Completed != 1 {
		t.Fatalf("TaskCycle=%+v err=%v", result, err)
	}
	var probe, task, repository bool
	for _, req := range recording.requests {
		switch {
		case strings.Contains(req.OutputPath, "-probe-"):
			probe = true
		case req.Command == "declared check":
			task = true
		case req.Command == settlementRepositoryCommand:
			repository = true
		}
	}
	if !probe || !task || !repository {
		t.Fatalf("did not exercise every Verification boundary: probe=%t task=%t repository=%t requests=%+v", probe, task, repository, recording.requests)
	}
}
