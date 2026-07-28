package cli

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"

	roundconfig "roundfix/internal/config"
)

const (
	detachTestChildSleepBeforeLiveness = "sleep-before-liveness"
	detachTestChildSlowRunCreation     = "slow-run-creation"
	detachTestChildRunCreated          = "run-created"
	detachTestChildInvalidLiveness     = "invalid-liveness"
	detachTestChildMalformedRunCreated = "malformed-run-created"
	detachTestChildExitTwoWithOutput   = "exit-two-with-output"
	detachTestChildExitOneSilently     = "exit-one-silently"
	detachTestRunID                    = "run-detach-test"
)

func TestRunDetachedCommandTimesOutWaitingForLiveness(t *testing.T) {
	stdout, stderr, code := runDetachParentForTest(t, detachTestChildSleepBeforeLiveness, 50*time.Millisecond, time.Second)

	if code != exitRunFailed {
		t.Fatalf("expected liveness timeout exit %d, got %d stdout=%q stderr=%q", exitRunFailed, code, stdout, stderr)
	}
	if stdout != "" {
		t.Fatalf("expected no stdout, got %q", stdout)
	}
	for _, want := range []string{
		"Detached Run child produced no liveness signal within 50ms; killed (exit:",
		"killed",
	} {
		if !strings.Contains(stderr, want) {
			t.Fatalf("expected liveness timeout diagnostic containing %q, got %q", want, stderr)
		}
	}
}

func TestRunDetachedCommandMonitorReportAllowsRunCreationPastLivenessDeadline(t *testing.T) {
	stdout, stderr, code := runDetachParentForTest(t, detachTestChildRunCreated, 500*time.Millisecond, 2*time.Second)

	if code != exitOK {
		t.Fatalf("expected detach success exit %d, got %d stdout=%q stderr=%q", exitOK, code, stdout, stderr)
	}
	if stderr != "" {
		t.Fatalf("expected no stderr, got %q", stderr)
	}
	runID, consoleLog := parseDetachedReport(t, stdout)
	if runID != detachTestRunID {
		t.Fatalf("expected run id %q, got %q", detachTestRunID, runID)
	}
	wantStdout := fmt.Sprintf(
		"Run ID: %s\nConsole Log: %s\nAttach: roundfix attach %s\nSupervisor monitor: roundfix events %s --follow --filter outcome\nStop: roundfix stop %s\n",
		runID,
		consoleLog,
		runID,
		runID,
		runID,
	)
	if stdout != wantStdout {
		t.Fatalf("detach stdout mismatch\nwant: %q\ngot:  %q", wantStdout, stdout)
	}
}

func TestRunDetachedCommandTimesOutWaitingForRunCreation(t *testing.T) {
	stdout, stderr, code := runDetachParentForTest(t, detachTestChildSlowRunCreation, 500*time.Millisecond, 50*time.Millisecond)

	if code != exitRunFailed {
		t.Fatalf("expected Run-creation timeout exit %d, got %d stdout=%q stderr=%q", exitRunFailed, code, stdout, stderr)
	}
	if stdout != "" {
		t.Fatalf("expected no stdout, got %q", stdout)
	}
	for _, want := range []string{
		"Detached Run child did not create a Run within 50ms; killed (exit:",
		"killed",
	} {
		if !strings.Contains(stderr, want) {
			t.Fatalf("expected Run-creation timeout diagnostic containing %q, got %q", want, stderr)
		}
	}
}

func TestRunDetachedCommandReportsInvalidLivenessHandshake(t *testing.T) {
	stdout, stderr, code := runDetachParentForTest(t, detachTestChildInvalidLiveness, time.Second, time.Second)

	if code != exitRunFailed {
		t.Fatalf("expected invalid liveness exit %d, got %d stdout=%q stderr=%q", exitRunFailed, code, stdout, stderr)
	}
	if stdout != "" {
		t.Fatalf("expected no stdout, got %q", stdout)
	}
	for _, want := range []string{
		"Detached Run child failed liveness handshake",
		"invalid Detached Run liveness marker",
	} {
		if !strings.Contains(stderr, want) {
			t.Fatalf("expected liveness failure diagnostic containing %q, got %q", want, stderr)
		}
	}
}

func TestRunDetachedCommandReportsMalformedRunCreationHandshake(t *testing.T) {
	stdout, stderr, code := runDetachParentForTest(t, detachTestChildMalformedRunCreated, time.Second, time.Second)

	if code != exitRunFailed {
		t.Fatalf("expected malformed run creation exit %d, got %d stdout=%q stderr=%q", exitRunFailed, code, stdout, stderr)
	}
	if stdout != "" {
		t.Fatalf("expected no stdout, got %q", stdout)
	}
	for _, want := range []string{
		"Detached Run child failed Run creation handshake",
		"invalid Detached Run handshake",
	} {
		if !strings.Contains(stderr, want) {
			t.Fatalf("expected Run creation failure diagnostic containing %q, got %q", want, stderr)
		}
	}
}

func TestRunDetachedCommandReportsChildExitBeforeHandshakeWithOutput(t *testing.T) {
	stdout, stderr, code := runDetachParentForTest(t, detachTestChildExitTwoWithOutput, time.Second, time.Second)

	if code != exitPreflight {
		t.Fatalf("expected child exit code %d, got %d stdout=%q stderr=%q", exitPreflight, code, stdout, stderr)
	}
	if stdout != "" {
		t.Fatalf("expected no stdout, got %q", stdout)
	}
	for _, want := range []string{
		"Detached Run child failed liveness handshake: EOF; child exited (exit status 2); console output follows",
		"detached child wrote a pre-handshake failure",
	} {
		if !strings.Contains(stderr, want) {
			t.Fatalf("expected child-exit diagnostic containing %q, got %q", want, stderr)
		}
	}
}

func TestRunDetachedCommandReportsSilentChildExitBeforeHandshake(t *testing.T) {
	stdout, stderr, code := runDetachParentForTest(t, detachTestChildExitOneSilently, time.Second, time.Second)

	if code != exitRunFailed {
		t.Fatalf("expected child exit code %d, got %d stdout=%q stderr=%q", exitRunFailed, code, stdout, stderr)
	}
	if stdout != "" {
		t.Fatalf("expected no stdout, got %q", stdout)
	}
	if !strings.Contains(stderr, "Detached Run child failed liveness handshake: EOF; child exited (exit status 1) and produced no output") {
		t.Fatalf("expected silent child diagnostic, got %q", stderr)
	}
}

func runDetachParentForTest(t *testing.T, childMode string, livenessTimeout time.Duration, runCreationTimeout time.Duration) (string, string, int) {
	t.Helper()
	t.Setenv(detachTestChildModeEnv, childMode)
	oldTimeouts := detachTimeouts
	detachTimeouts = detachPhaseTimeouts{
		liveness:    livenessTimeout,
		runCreation: runCreationTimeout,
	}
	t.Cleanup(func() {
		detachTimeouts = oldTimeouts
	})

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	code := runDetachedCommand(
		[]string{"implement", "--detach"},
		commandRequest{name: "implement", detach: true, artifactDir: filepath.Join(t.TempDir(), "artifacts")},
		roundconfig.Loaded{GitRoot: t.TempDir(), HomeDir: t.TempDir()},
		&stdout,
		&stderr,
	)
	return stdout.String(), stderr.String(), code
}

func runDetachTestChild(mode string) int {
	switch mode {
	case detachTestChildSleepBeforeLiveness:
		time.Sleep(2 * time.Second)
		return exitOK
	case detachTestChildSlowRunCreation:
		handshake, ok := writeDetachTestLiveness()
		if !ok {
			return exitRunFailed
		}
		defer func() {
			_ = handshake.Close()
		}()
		time.Sleep(2 * time.Second)
		return exitOK
	case detachTestChildRunCreated:
		handshake, ok := writeDetachTestLiveness()
		if !ok {
			return exitRunFailed
		}
		defer func() {
			_ = handshake.Close()
		}()
		time.Sleep(750 * time.Millisecond)
		consoleLog := filepath.Join(filepath.Dir(os.Getenv(detachConsoleTempEnv)), detachTestRunID, "console.log")
		if _, err := fmt.Fprintf(handshake, "%s\t%s\n", detachTestRunID, consoleLog); err != nil {
			return exitRunFailed
		}
		time.Sleep(75 * time.Millisecond)
		return exitOK
	case detachTestChildInvalidLiveness:
		handshake, ok := openDetachTestHandshake()
		if !ok {
			return exitRunFailed
		}
		defer func() {
			_ = handshake.Close()
		}()
		if _, err := handshake.Write([]byte{'x'}); err != nil {
			return exitRunFailed
		}
		return exitRunFailed
	case detachTestChildMalformedRunCreated:
		handshake, ok := writeDetachTestLiveness()
		if !ok {
			return exitRunFailed
		}
		defer func() {
			_ = handshake.Close()
		}()
		if _, err := fmt.Fprintln(handshake, "malformed"); err != nil {
			return exitRunFailed
		}
		return exitRunFailed
	case detachTestChildExitTwoWithOutput:
		fmt.Fprintln(os.Stderr, "detached child wrote a pre-handshake failure")
		return exitPreflight
	case detachTestChildExitOneSilently:
		return exitRunFailed
	default:
		fmt.Fprintf(os.Stderr, "unknown detach test child mode %q\n", mode)
		return exitRunFailed
	}
}

func writeDetachTestLiveness() (*os.File, bool) {
	handshake, ok := openDetachTestHandshake()
	if !ok {
		return nil, false
	}
	if _, err := handshake.Write([]byte{detachLivenessMarker}); err != nil {
		fmt.Fprintf(os.Stderr, "write detach test liveness marker: %v\n", err)
		_ = handshake.Close()
		return nil, false
	}
	return handshake, true
}

func openDetachTestHandshake() (*os.File, bool) {
	fd, err := strconv.Atoi(strings.TrimSpace(os.Getenv(detachHandshakeFDEnv)))
	if err != nil {
		fmt.Fprintf(os.Stderr, "invalid detach test handshake fd: %v\n", err)
		return nil, false
	}
	handshake := os.NewFile(uintptr(fd), "roundfix-detach-test-handshake")
	if handshake == nil {
		fmt.Fprintf(os.Stderr, "open detach test handshake fd %d failed\n", fd)
		return nil, false
	}
	return handshake, true
}
