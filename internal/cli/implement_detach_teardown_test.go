// Suite: detached test fixture ownership.
// Invariant: fixture processes end when their owning test binary is killed.
// Boundary IN: the compiled test binary, private TMPDIR and recorded PIDs.
// Boundary OUT: production process control and unrelated host processes.
package cli

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"testing"
	"time"

	"roundfix/internal/store"
	"roundfix/internal/testwait"
)

func TestImplementDetachChildEndsWhenItsTestBinaryDies(t *testing.T) {
	t.Parallel()
	runDetachFixtureDeathTest(t, "TestRunImplementDetachSurvivesCallerProcessGroupKill", false)
}

func TestDetachSurvivorEndsWhenItsTestBinaryDies(t *testing.T) {
	t.Parallel()
	runDetachFixtureDeathTest(t, "TestDetachedChildIsTerminatedAtTeardown", true)
}

func runDetachFixtureDeathTest(t *testing.T, innerTest string, survivor bool) {
	t.Helper()
	root := t.TempDir()
	tmpDir := filepath.Join(root, "tmp")
	mustMkdir(t, tmpDir)
	record := filepath.Join(root, "owner-pid")
	logPath := filepath.Join(root, "inner.log")
	log, err := os.Create(logPath)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if err := log.Close(); err != nil {
			t.Errorf("close inner test log: %v", err)
		}
		if t.Failed() {
			output, err := os.ReadFile(logPath)
			t.Logf("inner test output (read error %v):\n%s", err, output)
		}
	})
	// The outer test owns the deadline. The inner binary must still be alive
	// at that deadline so its own timeout cannot rescue a broken owner watch.
	cmd := exec.Command(os.Args[0], "-test.run=^"+innerTest+"$", "-test.v",
		"-test.timeout="+(testwait.Bound(t)+time.Minute).String())
	cmd.Env = withEnvValue(os.Environ(), "TMPDIR", tmpDir)
	cmd.Env = withEnvValue(cmd.Env, detachTestOwnerRecordEnv, record)
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	cmd.Stdout, cmd.Stderr = log, log
	if err := cmd.Start(); err != nil {
		t.Fatalf("start inner test binary: %v", err)
	}
	ended := make(chan error, 1)
	go func() { ended <- cmd.Wait() }()
	ownerPID := 0
	reaped := false
	t.Cleanup(func() {
		// On a failed owner-watch assertion, reclaim only our recorded group.
		killDetachFixtureGroup(t, cmd.Process.Pid)
		if !reaped {
			_ = testwait.Until(t, "inner test binary exit", ended, (<-chan error)(nil))
		}
		if ownerPID > 0 {
			killDetachFixtureGroup(t, ownerPID)
		}
	})
	testwait.Poll(t, "live recorded detached fixture PID", ended, func() (bool, string) {
		pidPath := record
		if survivor {
			paths, err := filepath.Glob(filepath.Join(tmpDir, "roundfix-detach-test-*", "pid"))
			if err != nil {
				t.Fatalf("find survivor PID record: %v", err)
			}
			if len(paths) != 1 {
				return false, fmt.Sprintf("survivor PID records: %v", paths)
			}
			pidPath = paths[0]
		}
		content, err := os.ReadFile(pidPath)
		if errors.Is(err, os.ErrNotExist) {
			return false, "PID record not written yet"
		}
		if err != nil {
			t.Fatalf("read fixture PID record: %v", err)
		}
		pid, err := strconv.Atoi(strings.TrimSpace(string(content)))
		if err != nil {
			return false, fmt.Sprintf("PID record not complete: %q", content)
		}
		if pid <= 0 || pid == os.Getpid() || pid == cmd.Process.Pid {
			t.Fatalf("invalid detached fixture PID %d", pid)
		}
		ownerPID = pid
		if !store.ProcessAlive(ownerPID) {
			t.Fatalf("fixture %d exited before its test binary was killed", ownerPID)
		}
		return true, "fixture alive"
	})
	killDetachFixtureGroup(t, cmd.Process.Pid)
	waitErr := testwait.Until(t, "killed inner test binary", ended, (<-chan error)(nil))
	reaped = true
	var exitErr *exec.ExitError
	if !errors.As(waitErr, &exitErr) {
		t.Fatalf("expected inner binary killed by SIGKILL, got %v", waitErr)
	}
	status, ok := exitErr.Sys().(syscall.WaitStatus)
	if !ok || !status.Signaled() || status.Signal() != syscall.SIGKILL {
		t.Fatalf("inner binary was not killed by SIGKILL: %v", waitErr)
	}
	// No cleanup signal is sent to the fixture before this assertion. Group
	// absence also proves the fake ACPX and its shell children have exited.
	waitForDetachFixtureExit(t, ownerPID)
}

func killDetachFixtureGroup(t *testing.T, pid int) {
	t.Helper()
	if pid <= 0 || pid == os.Getpid() {
		t.Fatalf("refuse to signal invalid fixture process group %d", pid)
	}
	err := syscall.Kill(-pid, syscall.SIGKILL)
	switch {
	case err == nil, errors.Is(err, syscall.ESRCH):
	case errors.Is(err, syscall.EPERM):
		// A sandboxed Verification may deny signals to the group; the probe
		// below skips the test there, so cleanup only records the denial.
		t.Logf("kill fixture process group %d: %v", pid, err)
	default:
		t.Errorf("kill fixture process group %d: %v", pid, err)
	}
}

func waitForDetachFixtureExit(t *testing.T, pid int) {
	t.Helper()
	testwait.Poll(t, fmt.Sprintf("fixture %d and its process group to exit", pid), (<-chan error)(nil), func() (bool, string) {
		alive := store.ProcessAlive(pid)
		err := syscall.Kill(-pid, 0)
		if errors.Is(err, syscall.EPERM) {
			t.Skipf("probe fixture process group %d: %v; this environment denies the signal the test observes, and the CI Verification gate runs it", pid, err)
		}
		if err != nil && !errors.Is(err, syscall.ESRCH) {
			t.Fatalf("probe fixture process group %d: %v", pid, err)
		}
		groupAlive := err == nil
		return !alive && !groupAlive, fmt.Sprintf("PID alive=%v; process group alive=%v", alive, groupAlive)
	})
}

// A record requests an inner death-test rendezvous. Normal survival tests
// take their original path; death tests hold the fixture at its live wait
// until the outer binary kills them. The wait remains deadline-bounded.
func waitForDetachDeathTestKill(t *testing.T, record string) {
	t.Helper()
	waitForFile(t, record+".release", nil)
}
