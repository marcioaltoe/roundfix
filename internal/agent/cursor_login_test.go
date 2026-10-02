// Suite: Cursor login refusal.
// Invariant: an unlogged Cursor selection opens no Agent Session and retains no status output.
// Boundary IN: adapter checks, exact proof, explicit environment and subprocess cancellation.
// Boundary OUT: Cursor service and real credentials; the executable is the compiled fixture.
package agent

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func cursorFixtureRuntime(t *testing.T, fixture fakeAdapterFixture) RuntimeSpec {
	t.Helper()
	path := filepath.Join(t.TempDir(), "cursor-agent")
	if err := provisionFakeAdapter(path, fixture); err != nil {
		t.Fatal(err)
	}
	return RuntimeSpec{ID: "cursor", Protocol: ProtocolStdio, Command: path, Model: "default[]"}
}

func cursorRefusals() map[string]fakeAdapterFixture {
	output := "Not logged in\nprivate account sentinel"
	return map[string]fakeAdapterFixture{
		"message": {StatusOutput: &output},
		"stderr":  {StatusStderr: &output},
		"exit":    {StatusExitCode: 1},
	}
}

func assertCursorLoginRequired(t *testing.T, err error) {
	t.Helper()
	var loginErr CursorLoginRequiredError
	if !errors.As(err, &loginErr) {
		t.Fatalf("error = %T %v, want CursorLoginRequiredError", err, err)
	}
	if err.Error() != "cursor-agent is not logged in" || loginErr.Classification() != CursorLoginRequired || loginErr.NextAction() != "run cursor-agent login in a terminal yourself" {
		t.Fatalf("unexpected refusal: %#v", loginErr)
	}
}

func TestCursorAdapterRefusedWithoutLogin(t *testing.T) {
	t.Parallel()
	for name, fixture := range cursorRefusals() {
		t.Run(name, func(t *testing.T) {
			runtime := cursorFixtureRuntime(t, fixture)
			_, err := checkAdapter(context.Background(), runtime, []string{})
			assertCursorLoginRequired(t, err)
		})
	}
}

func TestCursorAdapterReadyWithLogin(t *testing.T) {
	t.Parallel()
	output := "Logged in as private account sentinel"
	runtime := cursorFixtureRuntime(t, fakeAdapterFixture{StatusOutput: &output})
	evidence, err := checkAdapter(context.Background(), runtime, []string{})
	if err != nil {
		t.Fatal(err)
	}
	if evidence.Command != runtime.Command {
		t.Fatalf("evidence = %#v", evidence)
	}
}

func TestCursorProofRefusesBeforeAnySessionWithoutLogin(t *testing.T) {
	t.Parallel()
	for name, fixture := range cursorRefusals() {
		t.Run(name, func(t *testing.T) {
			harness := newFakeACPXHarness(t)
			runtime := cursorFixtureRuntime(t, fixture)
			_, err := harness.runner.ProveExactSelection(context.Background(), ProbeRequest{Runtime: runtime, WorkDir: harness.gitRoot})
			assertCursorLoginRequired(t, err)
			if _, err := os.Stat(harness.invocationsPath); !errors.Is(err, os.ErrNotExist) {
				t.Fatalf("acpx invoked or invocation evidence unreadable: %v", err)
			}
		})
	}
}

func TestCursorLoginHonorsCallerDeadline(t *testing.T) {
	t.Parallel()
	runtime := cursorFixtureRuntime(t, fakeAdapterFixture{StatusBlock: true})
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()
	assertCursorLoginRequired(t, checkCursorLogin(ctx, runtime.Command, []string{}))
	if !errors.Is(ctx.Err(), context.DeadlineExceeded) {
		t.Fatalf("context = %v", ctx.Err())
	}
}

func TestCursorLoginUsesExplicitEnvironment(t *testing.T) {
	t.Parallel()
	const statusEnv = "ROUNDFIX_TEST_CURSOR_STATUS"
	runtime := cursorFixtureRuntime(t, fakeAdapterFixture{StatusOutputEnv: statusEnv})
	assertCursorLoginRequired(t, checkCursorLogin(context.Background(), runtime.Command, []string{statusEnv + "=Not logged in"}))
	if err := checkCursorLogin(context.Background(), runtime.Command, []string{statusEnv + "=Logged in"}); err != nil {
		t.Fatal(err)
	}
}

func TestCursorCustomAdapterAlsoRequiresLogin(t *testing.T) {
	t.Parallel()
	output := "Not logged in"
	runtime := cursorFixtureRuntime(t, fakeAdapterFixture{StatusOutput: &output})
	runtime.ID = "cursor-custom"
	_, err := checkAdapter(context.Background(), runtime, []string{})
	assertCursorLoginRequired(t, err)
}
