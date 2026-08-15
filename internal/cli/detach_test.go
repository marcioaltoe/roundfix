package cli

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"

	roundconfig "roundfix/internal/config"
	"roundfix/internal/store"
)

const (
	detachTestChildSleepBeforeLiveness = "sleep-before-liveness"
	detachTestChildSlowRunCreation     = "slow-run-creation"
	detachTestChildRunCreated          = "run-created"
	detachTestChildInvalidLiveness     = "invalid-liveness"
	detachTestChildMalformedRunCreated = "malformed-run-created"
	detachTestChildExitTwoWithOutput   = "exit-two-with-output"
	detachTestChildExitOneSilently     = "exit-one-silently"
	detachTestChildIgnoreSentinel      = "ignore-sentinel"
	detachTestRunID                    = "run-detach-test"
	detachTestPIDPathEnv               = "ROUNDFIX_DETACH_TEST_PID_PATH"
	detachTestSentinelPathEnv          = "ROUNDFIX_DETACH_TEST_SENTINEL_PATH"
)

type detachTestSurvivor struct {
	root         string
	pidPath      string
	sentinelPath string
	pid          int
}

func TestRunDetachedCommandTimesOutWaitingForLiveness(t *testing.T) {
	t.Parallel()
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
	t.Parallel()
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
	t.Parallel()
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
	t.Parallel()
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
	t.Parallel()
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
	t.Parallel()
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

func TestRunDetachedReviewRefusalLeavesArtifactDirectoryAbsent(t *testing.T) {
	t.Parallel()
	for _, command := range []string{"resolve", "watch"} {
		t.Run(command, func(t *testing.T) {
			updateCommandDependenciesForTest(t, func(dependencies *commandDependencies) {
				dependencies.detachTimeouts = detachPhaseTimeouts{
					liveness:    time.Second,
					runCreation: time.Second,
				}
			})
			artifactParent := t.TempDir()
			artifactDir := filepath.Join(artifactParent, "artifacts")
			var stdout bytes.Buffer
			var stderr bytes.Buffer
			environment := commandEnvironmentForTest(t)

			code := runDetachedCommand(
				[]string{command, "--detach"},
				commandRequest{name: command, detach: true, artifactDir: artifactDir},
				roundconfig.Loaded{GitRoot: t.TempDir(), HomeDir: t.TempDir()},
				&stdout,
				&stderr,
				withEnvValue(environment.environ, detachTestChildModeEnv, detachTestChildExitTwoWithOutput),
				t.TempDir(),
				environment.dependencies.detachTimeouts,
			)

			if code != exitPreflight {
				t.Fatalf("expected child refusal exit %d, got %d stdout=%q stderr=%q", exitPreflight, code, stdout.String(), stderr.String())
			}
			if _, err := os.Stat(artifactDir); !errors.Is(err, os.ErrNotExist) {
				t.Fatalf("Artifact Directory exists after detached refusal: stat error = %v", err)
			}
			entries, err := os.ReadDir(artifactParent)
			if err != nil {
				t.Fatalf("read Artifact Directory parent after refusal: %v", err)
			}
			if len(entries) != 0 {
				t.Fatalf("detached refusal left filesystem entries: %v", entries)
			}
		})
	}
}

func TestRunDetachedCommandReportsSilentChildExitBeforeHandshake(t *testing.T) {
	t.Parallel()
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

func TestDetachedChildIsTerminatedAtTeardown(t *testing.T) {
	survivor := newDetachTestSurvivor(t)
	t.Cleanup(func() {
		if err := os.RemoveAll(survivor.root); err != nil {
			t.Errorf("remove detached child teardown directory: %v", err)
		}
	})
	command := exec.Command(os.Args[0], "-test.run=^$")
	command.Env = withEnvValue(os.Environ(), detachTestChildModeEnv, detachTestChildIgnoreSentinel)
	command.Env = withEnvValue(command.Env, detachTestSentinelPathEnv, survivor.sentinelPath)
	if err := command.Start(); err != nil {
		t.Fatalf("start detached child that ignores its sentinel: %v", err)
	}
	survivor.pid = command.Process.Pid

	t.Cleanup(func() {
		if err := survivor.terminate(); err != nil {
			t.Error(err)
		}
		_ = command.Process.Release()
		if store.ProcessAlive(survivor.pid) {
			t.Errorf("detached child %d remains alive after teardown", survivor.pid)
		}
	})
	if !store.ProcessAlive(survivor.pid) {
		t.Fatalf("detached child %d exited before teardown", survivor.pid)
	}
}

func runDetachParentForTest(t *testing.T, childMode string, livenessTimeout time.Duration, runCreationTimeout time.Duration) (string, string, int) {
	t.Helper()
	updateCommandDependenciesForTest(t, func(dependencies *commandDependencies) {
		dependencies.detachTimeouts = detachPhaseTimeouts{
			liveness:    livenessTimeout,
			runCreation: runCreationTimeout,
		}
	})

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	environment := commandEnvironmentForTest(t)
	childEnvironment := withEnvValue(environment.environ, detachTestChildModeEnv, childMode)
	var survivor *detachTestSurvivor
	if childMode == detachTestChildRunCreated {
		survivor = newDetachTestSurvivor(t)
		childEnvironment = withEnvValue(childEnvironment, detachTestPIDPathEnv, survivor.pidPath)
		childEnvironment = withEnvValue(childEnvironment, detachTestSentinelPathEnv, survivor.sentinelPath)
		t.Cleanup(func() {
			if err := survivor.terminate(); err != nil {
				t.Error(err)
			}
			if err := os.RemoveAll(survivor.root); err != nil {
				t.Errorf("remove detached child teardown directory: %v", err)
			}
		})
	}
	code := runDetachedCommand(
		[]string{"implement", "--detach"},
		commandRequest{name: "implement", detach: true, artifactDir: filepath.Join(t.TempDir(), "artifacts")},
		roundconfig.Loaded{GitRoot: t.TempDir(), HomeDir: t.TempDir()},
		&stdout,
		&stderr,
		childEnvironment,
		t.TempDir(),
		environment.dependencies.detachTimeouts,
	)
	if survivor != nil && code == exitOK {
		survivor.pid = readDetachTestSurvivorPID(t, survivor.pidPath)
		if !store.ProcessAlive(survivor.pid) {
			t.Fatalf("detached child %d did not survive its caller", survivor.pid)
		}
	}
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
		if !recordDetachTestSurvivor() {
			return exitRunFailed
		}
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
		return waitForDetachTestSentinel()
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
	case detachTestChildIgnoreSentinel:
		for {
			time.Sleep(time.Hour)
		}
	default:
		fmt.Fprintf(os.Stderr, "unknown detach test child mode %q\n", mode)
		return exitRunFailed
	}
}

func newDetachTestSurvivor(t *testing.T) *detachTestSurvivor {
	t.Helper()
	root, err := os.MkdirTemp("", "roundfix-detach-test-")
	if err != nil {
		t.Fatalf("create detached child teardown directory: %v", err)
	}
	return &detachTestSurvivor{
		root:         root,
		pidPath:      filepath.Join(root, "pid"),
		sentinelPath: filepath.Join(root, "release-sentinel"),
	}
}

func recordDetachTestSurvivor() bool {
	pidPath := strings.TrimSpace(os.Getenv(detachTestPIDPathEnv))
	if pidPath == "" {
		fmt.Fprintln(os.Stderr, "detached test child pid path is empty")
		return false
	}
	if err := os.WriteFile(pidPath, []byte(strconv.Itoa(os.Getpid())), 0o600); err != nil {
		fmt.Fprintf(os.Stderr, "record detached test child pid: %v\n", err)
		return false
	}
	return true
}

func readDetachTestSurvivorPID(t *testing.T, pidPath string) int {
	t.Helper()
	content, err := os.ReadFile(pidPath)
	if err != nil {
		t.Fatalf("read detached test child pid: %v", err)
	}
	pid, err := strconv.Atoi(strings.TrimSpace(string(content)))
	if err != nil || pid <= 0 {
		t.Fatalf("parse detached test child pid %q: %v", content, err)
	}
	return pid
}

func waitForDetachTestSentinel() int {
	sentinelPath := strings.TrimSpace(os.Getenv(detachTestSentinelPathEnv))
	if sentinelPath == "" {
		fmt.Fprintln(os.Stderr, "detached test child sentinel path is empty")
		return exitRunFailed
	}
	ticker := time.NewTicker(10 * time.Millisecond)
	defer ticker.Stop()
	for {
		if _, err := os.Stat(sentinelPath); err == nil {
			return exitOK
		} else if !errors.Is(err, os.ErrNotExist) {
			fmt.Fprintf(os.Stderr, "read detached test child sentinel: %v\n", err)
			return exitRunFailed
		}
		<-ticker.C
	}
}

func (survivor *detachTestSurvivor) terminate() error {
	if survivor == nil || survivor.pid <= 0 {
		return nil
	}
	process, err := os.FindProcess(survivor.pid)
	if err != nil {
		return fmt.Errorf("find detached child %d: %w", survivor.pid, err)
	}
	killErr := process.Kill()
	killFailedWhileAlive := killErr != nil && store.ProcessAlive(survivor.pid)
	var sentinelErr error
	if err := os.WriteFile(survivor.sentinelPath, []byte("release\n"), 0o600); err != nil {
		sentinelErr = fmt.Errorf("write detached child release sentinel: %w", err)
	}

	wait := make(chan error, 1)
	go func() {
		_, waitErr := process.Wait()
		wait <- waitErr
	}()
	deadline := time.NewTimer(detachWaitTimeout)
	defer deadline.Stop()
	select {
	case waitErr := <-wait:
		if waitErr != nil && store.ProcessAlive(survivor.pid) {
			return errors.Join(killErr, sentinelErr, fmt.Errorf("wait for detached child %d after kill: %w", survivor.pid, waitErr))
		}
	case <-deadline.C:
		return errors.Join(killErr, sentinelErr, fmt.Errorf("detached child %d remains alive after non-cooperative teardown", survivor.pid))
	}
	if store.ProcessAlive(survivor.pid) {
		return errors.Join(killErr, sentinelErr, fmt.Errorf("detached child %d remains alive after non-cooperative teardown", survivor.pid))
	}
	if killFailedWhileAlive {
		return errors.Join(fmt.Errorf("kill detached child %d: %w", survivor.pid, killErr), sentinelErr)
	}
	return sentinelErr
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
