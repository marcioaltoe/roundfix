//go:build unix

package store

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"runtime"
	"strings"
	"syscall"
	"testing"
	"time"
)

func TestProcessAliveReportsCurrentProcessAlive(t *testing.T) {
	if !ProcessAlive(os.Getpid()) {
		t.Fatal("expected current process to report alive")
	}
}

func TestProcessAliveReportsReapedChildDead(t *testing.T) {
	cmd := exec.Command("/bin/sh", "-c", "exit 0")
	if err := cmd.Start(); err != nil {
		t.Fatalf("start child process: %v", err)
	}
	pid := cmd.Process.Pid
	if err := cmd.Wait(); err != nil {
		t.Fatalf("wait for child process: %v", err)
	}
	if ProcessAlive(pid) {
		t.Fatalf("expected reaped child pid %d to report dead", pid)
	}
}

func TestOwnerProcessControllerGracefulExitProof(t *testing.T) {
	pid, wait := startOwnerProcessHelper(t, "graceful")
	controller := newOwnerProcessController(250*time.Millisecond, 2*time.Second, 5*time.Millisecond)

	if err := controller.TerminateAndWait(t.Context(), pid, ""); err != nil {
		t.Fatalf("terminate owner process gracefully: %v", err)
	}

	assertOwnerProcessExited(t, pid, wait)
}

func TestOwnerProcessControllerForceKillExitProof(t *testing.T) {
	pid, wait := startOwnerProcessHelper(t, "ignore")
	controller := newOwnerProcessController(20*time.Millisecond, 2*time.Second, 5*time.Millisecond)

	if err := controller.TerminateAndWait(t.Context(), pid, ""); err != nil {
		t.Fatalf("force-kill owner process: %v", err)
	}

	assertOwnerProcessExited(t, pid, wait)
}

func TestOwnerProcessControllerRejectsUnprovenCurrentProcess(t *testing.T) {
	controller := newOwnerProcessController(20*time.Millisecond, 100*time.Millisecond, 5*time.Millisecond)

	err := controller.TerminateAndWait(t.Context(), os.Getpid(), "")

	if !errors.Is(err, ErrOwnerProcessIdentityUnproven) {
		t.Fatalf("current-process error = %v, want ErrOwnerProcessIdentityUnproven", err)
	}
	var controlErr OwnerProcessControlError
	if !errors.As(err, &controlErr) {
		t.Fatalf("current-process error type = %T, want OwnerProcessControlError", err)
	}
	if controlErr.PID != os.Getpid() || controlErr.Step != "prove owner process identity" {
		t.Fatalf("current-process diagnostic = %#v", controlErr)
	}
}

func TestOwnerProcessControllerAcceptsAlreadyAbsentProcess(t *testing.T) {
	cmd := exec.Command("/bin/sh", "-c", "exit 0")
	if err := cmd.Start(); err != nil {
		t.Fatalf("start child process: %v", err)
	}
	pid := cmd.Process.Pid
	if err := cmd.Wait(); err != nil {
		t.Fatalf("wait for child process: %v", err)
	}
	controller := newOwnerProcessController(20*time.Millisecond, 100*time.Millisecond, 5*time.Millisecond)

	if err := controller.TerminateAndWait(t.Context(), pid, ""); err != nil {
		t.Fatalf("accept already absent owner process: %v", err)
	}
}

func TestOwnerProcessControllerRefusesMismatchedOwnerIdentity(t *testing.T) {
	pid, _ := startOwnerProcessHelper(t, "graceful")
	controller := newOwnerProcessController(20*time.Millisecond, 2*time.Second, 5*time.Millisecond)

	err := controller.TerminateAndWait(t.Context(), pid, "identity-token-of-exited-owner-process")

	if !errors.Is(err, ErrOwnerProcessIdentityUnproven) {
		t.Fatalf("mismatched identity error = %v, want ErrOwnerProcessIdentityUnproven", err)
	}
	var controlErr OwnerProcessControlError
	if !errors.As(err, &controlErr) {
		t.Fatalf("mismatched identity error type = %T, want OwnerProcessControlError", err)
	}
	if controlErr.PID != pid || controlErr.Step != "prove owner process identity" {
		t.Fatalf("mismatched identity diagnostic = %#v", controlErr)
	}
	if !ProcessAlive(pid) {
		t.Fatalf("refusal must not signal process %d holding the reused PID", pid)
	}
}

func TestOwnerProcessControllerMatchingOwnerIdentityProceeds(t *testing.T) {
	pid, wait := startOwnerProcessHelper(t, "graceful")
	identity, err := OwnerProcessIdentity(t.Context(), pid)
	if err != nil {
		t.Fatalf("read owner process identity: %v", err)
	}
	controller := newOwnerProcessController(250*time.Millisecond, 2*time.Second, 5*time.Millisecond)

	if err := controller.TerminateAndWait(t.Context(), pid, identity); err != nil {
		t.Fatalf("terminate owner process with matching identity: %v", err)
	}

	assertOwnerProcessExited(t, pid, wait)
}

func TestOwnerProcessControllerProveOwnerLeavesProvenOwnerRunning(t *testing.T) {
	pid, _ := startOwnerProcessHelper(t, "graceful")
	identity, err := OwnerProcessIdentity(t.Context(), pid)
	if err != nil {
		t.Fatalf("read owner process identity: %v", err)
	}
	controller := newOwnerProcessController(250*time.Millisecond, 2*time.Second, 5*time.Millisecond)

	if err := controller.ProveOwner(t.Context(), pid, identity); err != nil {
		t.Fatalf("prove matching owner identity: %v", err)
	}

	if !ProcessAlive(pid) {
		t.Fatalf("owner proof must send no signal, but process %d exited", pid)
	}
}

func TestOwnerProcessControllerProveOwnerRefusesMismatchedIdentity(t *testing.T) {
	pid, _ := startOwnerProcessHelper(t, "graceful")
	controller := newOwnerProcessController(20*time.Millisecond, 2*time.Second, 5*time.Millisecond)

	err := controller.ProveOwner(t.Context(), pid, "identity-token-of-exited-owner-process")

	if !errors.Is(err, ErrOwnerProcessIdentityUnproven) {
		t.Fatalf("mismatched identity proof error = %v, want ErrOwnerProcessIdentityUnproven", err)
	}
	var controlErr OwnerProcessControlError
	if !errors.As(err, &controlErr) {
		t.Fatalf("mismatched identity proof error type = %T, want OwnerProcessControlError", err)
	}
	if controlErr.PID != pid || controlErr.Step != "prove owner process identity" {
		t.Fatalf("mismatched identity proof diagnostic = %#v", controlErr)
	}
	if !ProcessAlive(pid) {
		t.Fatalf("refused proof must not signal process %d holding the reused PID", pid)
	}
}

func TestOwnerProcessControllerProveOwnerAcceptsAbsentOwner(t *testing.T) {
	cmd := exec.Command("/bin/sh", "-c", "exit 0")
	if err := cmd.Start(); err != nil {
		t.Fatalf("start child process: %v", err)
	}
	pid := cmd.Process.Pid
	if err := cmd.Wait(); err != nil {
		t.Fatalf("wait for child process: %v", err)
	}
	controller := newOwnerProcessController(20*time.Millisecond, 100*time.Millisecond, 5*time.Millisecond)

	if err := controller.ProveOwner(t.Context(), pid, "identity-token-of-exited-owner-process"); err != nil {
		t.Fatalf("absence is its own proof, got: %v", err)
	}
}

func TestOwnerProcessControllerProveOwnerRejectsCurrentProcess(t *testing.T) {
	controller := newOwnerProcessController(20*time.Millisecond, 100*time.Millisecond, 5*time.Millisecond)

	err := controller.ProveOwner(t.Context(), os.Getpid(), "")

	if !errors.Is(err, ErrOwnerProcessIdentityUnproven) {
		t.Fatalf("current-process proof error = %v, want ErrOwnerProcessIdentityUnproven", err)
	}
}

func TestOwnerProcessIdentityIsStableForOneProcess(t *testing.T) {
	if runtime.GOOS != "linux" && runtime.GOOS != "darwin" {
		t.Skip("owner process identity is unsupported on this Unix platform")
	}
	first, err := OwnerProcessIdentity(t.Context(), os.Getpid())
	if err != nil {
		t.Fatalf("read current process identity: %v", err)
	}
	second, err := OwnerProcessIdentity(t.Context(), os.Getpid())
	if err != nil {
		t.Fatalf("re-read current process identity: %v", err)
	}
	if first != second {
		t.Fatalf("identity token changed for one process: %q then %q", first, second)
	}
	if prefix := runtime.GOOS + ":"; !strings.HasPrefix(first, prefix) {
		t.Fatalf("identity token = %q, want prefix %q", first, prefix)
	}
}

func TestOwnerProcessIdentityDoesNotSpawnPS(t *testing.T) {
	if runtime.GOOS != "linux" && runtime.GOOS != "darwin" {
		t.Skip("owner process identity is unsupported on this Unix platform")
	}
	binDir := t.TempDir()
	psPath := filepath.Join(binDir, "ps")
	if err := os.WriteFile(psPath, []byte("#!/bin/sh\nexit 99\n"), 0o755); err != nil {
		t.Fatalf("write failing ps executable: %v", err)
	}
	t.Setenv("PATH", binDir)

	if _, err := OwnerProcessIdentity(t.Context(), os.Getpid()); err != nil {
		t.Fatalf("read current process identity without ps: %v", err)
	}
}

func TestOwnerProcessIdentityIgnoresCallerTimezone(t *testing.T) {
	if runtime.GOOS != "linux" && runtime.GOOS != "darwin" {
		t.Skip("owner process identity is unsupported on this Unix platform")
	}
	t.Setenv("TZ", "Pacific/Honolulu")
	first, err := OwnerProcessIdentity(t.Context(), os.Getpid())
	if err != nil {
		t.Fatalf("read current process identity in first timezone: %v", err)
	}
	t.Setenv("TZ", "Asia/Tokyo")
	second, err := OwnerProcessIdentity(t.Context(), os.Getpid())
	if err != nil {
		t.Fatalf("read current process identity in second timezone: %v", err)
	}
	if first != second {
		t.Fatalf("identity token changed with caller timezone: %q then %q", first, second)
	}
}

func TestOwnerProcessIdentityFailsForAbsentProcess(t *testing.T) {
	if runtime.GOOS != "linux" && runtime.GOOS != "darwin" {
		t.Skip("owner process identity is unsupported on this Unix platform")
	}
	cmd := exec.Command("/bin/sh", "-c", "exit 0")
	if err := cmd.Start(); err != nil {
		t.Fatalf("start child process: %v", err)
	}
	pid := cmd.Process.Pid
	if err := cmd.Wait(); err != nil {
		t.Fatalf("wait for child process: %v", err)
	}

	identity, err := OwnerProcessIdentity(t.Context(), pid)
	if err == nil {
		t.Fatalf("expected identity read failure for reaped pid %d, got %q", pid, identity)
	}
	want := error(syscall.ENOENT)
	if runtime.GOOS == "darwin" {
		want = syscall.ESRCH
	}
	if !errors.Is(err, want) {
		t.Fatalf("identity read error for reaped pid %d = %v, want %v", pid, err, want)
	}
}

func TestOwnerProcessHelper(t *testing.T) {
	mode := os.Getenv("ROUNDFIX_OWNER_PROCESS_HELPER")
	if mode == "" {
		return
	}
	switch mode {
	case "graceful":
		signals := make(chan os.Signal, 1)
		signal.Notify(signals, syscall.SIGTERM)
		defer signal.Stop(signals)
		fmt.Fprintln(os.Stdout, "ready")
		<-signals
	case "ignore":
		signal.Ignore(syscall.SIGTERM)
		defer signal.Reset(syscall.SIGTERM)
		fmt.Fprintln(os.Stdout, "ready")
		select {}
	default:
		t.Fatalf("unknown owner process helper mode %q", mode)
	}
}

func startOwnerProcessHelper(t *testing.T, mode string) (int, <-chan error) {
	t.Helper()
	cmd := exec.Command(os.Args[0], "-test.run=^TestOwnerProcessHelper$")
	cmd.Env = append(os.Environ(), "ROUNDFIX_OWNER_PROCESS_HELPER="+mode)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		t.Fatalf("open owner process helper stdout: %v", err)
	}
	cmd.Stderr = os.Stderr
	if err := cmd.Start(); err != nil {
		t.Fatalf("start owner process helper: %v", err)
	}
	wait := make(chan error, 1)
	go func() {
		wait <- cmd.Wait()
		close(wait)
	}()
	t.Cleanup(func() {
		if ProcessAlive(cmd.Process.Pid) {
			_ = cmd.Process.Kill()
		}
		select {
		case <-wait:
		case <-time.After(2 * time.Second):
		}
	})
	scanner := bufio.NewScanner(stdout)
	if !scanner.Scan() {
		t.Fatalf("owner process helper did not become ready: %v", scanner.Err())
	}
	if scanner.Text() != "ready" {
		t.Fatalf("owner process helper readiness = %q, want ready", scanner.Text())
	}
	return cmd.Process.Pid, wait
}

func assertOwnerProcessExited(t *testing.T, pid int, wait <-chan error) {
	t.Helper()
	select {
	case <-wait:
	case <-time.After(2 * time.Second):
		t.Fatalf("owner process %d did not exit", pid)
	}
	if ProcessAlive(pid) {
		t.Fatalf("owner process %d remained alive after exit proof", pid)
	}
}
