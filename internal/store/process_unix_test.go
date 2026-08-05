//go:build unix

package store

import (
	"bufio"
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"reflect"
	"runtime"
	"strconv"
	"strings"
	"syscall"
	"testing"
	"time"

	"golang.org/x/sys/unix"
)

func TestProcessAliveReportsCurrentProcessAlive(t *testing.T) {
	t.Parallel()
	if !ProcessAlive(os.Getpid()) {
		t.Fatal("expected current process to report alive")
	}
}

func TestProcessAliveReportsReapedChildDead(t *testing.T) {
	t.Parallel()
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
	t.Parallel()
	pid, wait := startOwnerProcessHelper(t, "graceful")
	controller := newOwnerProcessController(250*time.Millisecond, 2*time.Second, 5*time.Millisecond)

	if err := controller.TerminateAndWait(t.Context(), pid, ""); err != nil {
		t.Fatalf("terminate owner process gracefully: %v", err)
	}

	assertOwnerProcessExited(t, pid, wait)
}

func TestOwnerProcessControllerForceKillExitProof(t *testing.T) {
	t.Parallel()
	pid, wait := startOwnerProcessHelper(t, "ignore")
	controller := newOwnerProcessController(20*time.Millisecond, 2*time.Second, 5*time.Millisecond)

	if err := controller.TerminateAndWait(t.Context(), pid, ""); err != nil {
		t.Fatalf("force-kill owner process: %v", err)
	}

	assertOwnerProcessForceKilled(t, pid, wait)
}

func TestOwnerProcessControllerInspectTreeReportsOnlyLiveOwnedProcesses(t *testing.T) {
	t.Parallel()
	const ownerPID = 7101
	const childPID = 7102
	const exitedPID = 7103
	controller := newOwnerProcessController(20*time.Millisecond, 100*time.Millisecond, 5*time.Millisecond)
	controller.processAbsent = func(pid int) (bool, error) {
		return pid == exitedPID, nil
	}
	controller.ownedProcesses = func(pid int) ([]int, error) {
		if pid != ownerPID {
			t.Fatalf("owned process root = %d, want %d", pid, ownerPID)
		}
		return []int{childPID, exitedPID, ownerPID}, nil
	}

	got, err := controller.InspectTree(t.Context(), ownerPID, "")

	if err != nil {
		t.Fatalf("inspect owner process tree: %v", err)
	}
	want := []int{ownerPID, childPID}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("live owned processes = %v, want %v", got, want)
	}
}

func TestOwnerProcessControllerInspectTreeFindsOutlivingGrandchild(t *testing.T) {
	t.Parallel()
	ownerPID, grandchildPID, wait := startOwnerProcessTreeHelper(t)
	controller := newOwnerProcessController(250*time.Millisecond, 2*time.Second, 5*time.Millisecond)

	got, err := controller.InspectTree(t.Context(), ownerPID, "")
	if err != nil {
		t.Fatalf("inspect real owner process tree: %v", err)
	}
	want := []int{ownerPID, grandchildPID}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("live real process tree = %v, want %v", got, want)
	}
	if _, err := controller.TerminateTreeAndWait(t.Context(), ownerPID, ""); err != nil {
		t.Fatalf("terminate inspected real process tree: %v", err)
	}
	assertOwnerProcessExited(t, ownerPID, wait)
}

func TestOwnerProcessControllerTerminateTreeProvesOutlivingGrandchildGone(t *testing.T) {
	t.Parallel()
	ownerPID, grandchildPID, wait := startOwnerProcessTreeHelper(t)
	controller := newOwnerProcessController(250*time.Millisecond, 2*time.Second, 5*time.Millisecond)

	outcomes, err := controller.TerminateTreeAndWait(t.Context(), ownerPID, "")
	if err != nil {
		t.Fatalf("terminate owner process tree: %v", err)
	}

	wantPIDs := map[int]bool{ownerPID: false, grandchildPID: false}
	for _, outcome := range outcomes {
		if _, ok := wantPIDs[outcome.PID]; !ok {
			t.Fatalf("termination outcome for unowned process %d: %#v", outcome.PID, outcome)
		}
		if !outcome.Proven || outcome.Reason != "" {
			t.Fatalf("termination outcome for process %d = %#v, want proven absence", outcome.PID, outcome)
		}
		wantPIDs[outcome.PID] = true
	}
	for pid, observed := range wantPIDs {
		if !observed {
			t.Fatalf("missing termination outcome for process %d: %#v", pid, outcomes)
		}
		if ProcessAlive(pid) {
			t.Fatalf("owned process %d remained alive after tree termination", pid)
		}
	}
	assertOwnerProcessExited(t, ownerPID, wait)
}

func TestOwnerProcessControllerTerminateTreeLeavesUnrelatedProcessRunning(t *testing.T) {
	t.Parallel()
	ownerPID, _, wait := startOwnerProcessTreeHelper(t)
	unrelatedPID, _ := startOwnerProcessHelper(t, "graceful")
	controller := newOwnerProcessController(250*time.Millisecond, 2*time.Second, 5*time.Millisecond)

	if _, err := controller.TerminateTreeAndWait(t.Context(), ownerPID, ""); err != nil {
		t.Fatalf("terminate owner process tree: %v", err)
	}

	assertOwnerProcessExited(t, ownerPID, wait)
	if !ProcessAlive(unrelatedPID) {
		t.Fatalf("tree termination signalled unrelated process %d", unrelatedPID)
	}
}

func TestOwnerProcessControllerTerminateTreeReportsUnprovenAbsence(t *testing.T) {
	t.Parallel()
	const pid = 424242
	hostErr := errors.New("host process table unavailable")
	controller := newOwnerProcessController(20*time.Millisecond, 100*time.Millisecond, 5*time.Millisecond)
	absenceReads := 0
	controller.processAbsent = func(gotPID int) (bool, error) {
		if gotPID != pid {
			t.Fatalf("absence read pid = %d, want %d", gotPID, pid)
		}
		absenceReads++
		if absenceReads <= 2 {
			return false, nil
		}
		return false, hostErr
	}
	controller.ownedProcesses = func(gotPID int) ([]int, error) {
		if gotPID != pid {
			t.Fatalf("owned process root = %d, want %d", gotPID, pid)
		}
		return []int{pid}, nil
	}
	signalCalls := 0
	controller.signalProcess = func(gotPID int, force bool) error {
		if gotPID != pid || force {
			t.Fatalf("signal = (%d, force=%v), want (%d, force=false)", gotPID, force, pid)
		}
		signalCalls++
		return nil
	}

	outcomes, err := controller.TerminateTreeAndWait(t.Context(), pid, "")
	if err != nil {
		t.Fatalf("terminate process with unobservable absence: %v", err)
	}

	if signalCalls != 1 {
		t.Fatalf("signal calls = %d, want 1", signalCalls)
	}
	if len(outcomes) != 1 {
		t.Fatalf("termination outcomes = %#v, want one", outcomes)
	}
	outcome := outcomes[0]
	if outcome.PID != pid || outcome.Proven || strings.TrimSpace(outcome.Reason) == "" {
		t.Fatalf("termination outcome = %#v, want unproven with non-empty reason", outcome)
	}
	if !strings.Contains(outcome.Reason, hostErr.Error()) {
		t.Fatalf("termination reason = %q, want host diagnostic %q", outcome.Reason, hostErr)
	}
}

func TestOwnerProcessControllerRejectsUnprovenCurrentProcess(t *testing.T) {
	t.Parallel()
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
	t.Parallel()
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
	t.Parallel()
	pid, _ := startOwnerProcessHelper(t, "graceful")
	controller := newOwnerProcessController(20*time.Millisecond, 2*time.Second, 5*time.Millisecond)

	err := controller.TerminateAndWait(t.Context(), pid, runtime.GOOS+":identity-token-of-exited-owner-process")

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
	t.Parallel()
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
	t.Parallel()
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
	t.Parallel()
	pid, _ := startOwnerProcessHelper(t, "graceful")
	controller := newOwnerProcessController(20*time.Millisecond, 2*time.Second, 5*time.Millisecond)

	err := controller.ProveOwner(t.Context(), pid, runtime.GOOS+":identity-token-of-exited-owner-process")

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

func TestOwnerProcessControllerProveOwnerDoesNotReportLegacyIdentityAsReusedPID(t *testing.T) {
	t.Parallel()
	pid, _ := startOwnerProcessHelper(t, "graceful")
	controller := newOwnerProcessController(20*time.Millisecond, 2*time.Second, 5*time.Millisecond)

	err := controller.ProveOwner(t.Context(), pid, "identity-token-from-legacy-ps-reader")

	if !errors.Is(err, ErrOwnerIdentityUnreadable) {
		t.Fatalf("legacy identity error = %v, want ErrOwnerIdentityUnreadable", err)
	}
	if errors.Is(err, ErrOwnerProcessIdentityUnproven) {
		t.Fatalf("legacy identity error = %v, must not report a reused PID", err)
	}
	if !strings.Contains(err.Error(), "not comparable on "+runtime.GOOS) {
		t.Fatalf("legacy identity error = %v, want platform comparison diagnostic", err)
	}
}

func TestOwnerProcessControllerProveOwnerClassifiesIdentityEvidence(t *testing.T) {
	t.Parallel()
	type absenceResult struct {
		absent bool
		err    error
	}

	hostErr := errors.New("resource temporarily unavailable")
	nativePrefix := runtime.GOOS + ":"
	foreignPrefix := "linux:"
	if runtime.GOOS == "linux" {
		foreignPrefix = "darwin:"
	}
	tests := []struct {
		name              string
		recordedIdentity  string
		liveIdentity      string
		identityErr       error
		absenceResults    []absenceResult
		wantAbsent        bool
		wantErr           error
		wantNotErr        error
		wantWrappedErr    error
		wantMessage       []string
		wantIdentityReads int
	}{
		{
			name:             "live identity read failure is unreadable",
			recordedIdentity: nativePrefix + "recorded",
			identityErr:      hostErr,
			absenceResults:   []absenceResult{{}, {}},
			wantErr:          ErrOwnerIdentityUnreadable,
			wantNotErr:       ErrOwnerProcessIdentityUnproven,
			wantWrappedErr:   hostErr,
			wantMessage: []string{
				"read live owner identity",
				"resource temporarily unavailable",
				"resolve the host resource failure, then retry",
			},
			wantIdentityReads: 1,
		},
		{
			name:              "owner exit after identity read failure is proven absent",
			recordedIdentity:  nativePrefix + "recorded",
			identityErr:       hostErr,
			absenceResults:    []absenceResult{{}, {absent: true}},
			wantAbsent:        true,
			wantIdentityReads: 1,
		},
		{
			name:              "legacy token without prefix is unreadable",
			recordedIdentity:  "legacy-ps-token",
			liveIdentity:      nativePrefix + "live",
			absenceResults:    []absenceResult{{}},
			wantErr:           ErrOwnerIdentityUnreadable,
			wantNotErr:        ErrOwnerProcessIdentityUnproven,
			wantMessage:       []string{"recorded owner identity is not comparable on " + runtime.GOOS},
			wantIdentityReads: 1,
		},
		{
			name:              "token from another platform is unreadable",
			recordedIdentity:  foreignPrefix + "recorded",
			liveIdentity:      nativePrefix + "live",
			absenceResults:    []absenceResult{{}},
			wantErr:           ErrOwnerIdentityUnreadable,
			wantNotErr:        ErrOwnerProcessIdentityUnproven,
			wantMessage:       []string{"recorded owner identity is not comparable on " + runtime.GOOS},
			wantIdentityReads: 1,
		},
		{
			name:              "comparable mismatch remains unproven",
			recordedIdentity:  nativePrefix + "recorded",
			liveIdentity:      nativePrefix + "live",
			absenceResults:    []absenceResult{{}},
			wantErr:           ErrOwnerProcessIdentityUnproven,
			wantNotErr:        ErrOwnerIdentityUnreadable,
			wantMessage:       []string{"does not match the recorded owner identity"},
			wantIdentityReads: 1,
		},
		{
			name:              "matching identity remains proven",
			recordedIdentity:  nativePrefix + "same",
			liveIdentity:      nativePrefix + "same",
			absenceResults:    []absenceResult{{}},
			wantIdentityReads: 1,
		},
		{
			name:             "initial absence remains proof",
			recordedIdentity: nativePrefix + "recorded",
			absenceResults:   []absenceResult{{absent: true}},
			wantAbsent:       true,
		},
		{
			name:             "empty recorded identity keeps PID-only proof",
			recordedIdentity: "  ",
			absenceResults:   []absenceResult{{}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			controller := newOwnerProcessController(20*time.Millisecond, 100*time.Millisecond, 5*time.Millisecond)
			absenceReads := 0
			controller.processAbsent = func(int) (bool, error) {
				if absenceReads >= len(tt.absenceResults) {
					t.Fatalf("unexpected process absence read %d", absenceReads+1)
				}
				result := tt.absenceResults[absenceReads]
				absenceReads++
				return result.absent, result.err
			}
			identityReads := 0
			controller.processStartIdentity = func(context.Context, int) (string, error) {
				identityReads++
				return tt.liveIdentity, tt.identityErr
			}

			absent, err := controller.proveOwner(t.Context(), os.Getpid()+1, tt.recordedIdentity)

			if absent != tt.wantAbsent {
				t.Fatalf("proven absent = %v, want %v", absent, tt.wantAbsent)
			}
			if tt.wantErr == nil {
				if err != nil {
					t.Fatalf("proof error = %v, want nil", err)
				}
			} else if !errors.Is(err, tt.wantErr) {
				t.Fatalf("proof error = %v, want %v", err, tt.wantErr)
			}
			if tt.wantNotErr != nil && errors.Is(err, tt.wantNotErr) {
				t.Fatalf("proof error = %v, must not match %v", err, tt.wantNotErr)
			}
			if tt.wantWrappedErr != nil && !errors.Is(err, tt.wantWrappedErr) {
				t.Fatalf("proof error = %v, want wrapped %v", err, tt.wantWrappedErr)
			}
			for _, want := range tt.wantMessage {
				if !strings.Contains(err.Error(), want) {
					t.Fatalf("proof error = %v, want diagnostic containing %q", err, want)
				}
			}
			if identityReads != tt.wantIdentityReads {
				t.Fatalf("identity reads = %d, want %d", identityReads, tt.wantIdentityReads)
			}
		})
	}
}

func TestOwnerProcessControllerProveOwnerAcceptsAbsentOwner(t *testing.T) {
	t.Parallel()
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
	t.Parallel()
	controller := newOwnerProcessController(20*time.Millisecond, 100*time.Millisecond, 5*time.Millisecond)

	err := controller.ProveOwner(t.Context(), os.Getpid(), "")

	if !errors.Is(err, ErrOwnerProcessIdentityUnproven) {
		t.Fatalf("current-process proof error = %v, want ErrOwnerProcessIdentityUnproven", err)
	}
}

func TestOwnerProcessIdentityIsStableForOneProcess(t *testing.T) {
	t.Parallel()
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
	t.Parallel()
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

func TestOwnerProcessHelperIgnoreModeStaysAlive(t *testing.T) {
	cmd := exec.Command(os.Args[0], "-test.run=^TestOwnerProcessHelper$")
	cmd.Env = append(os.Environ(), "ROUNDFIX_OWNER_PROCESS_HELPER=ignore")
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		t.Fatalf("open owner process helper stdout: %v", err)
	}
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Start(); err != nil {
		t.Fatalf("start owner process helper: %v", err)
	}

	wait := make(chan error, 1)
	waitStarted := false
	waited := false
	t.Cleanup(func() {
		_ = cmd.Process.Kill()
		if waitStarted {
			if !waited {
				<-wait
			}
			return
		}
		_ = cmd.Wait()
	})

	scanner := bufio.NewScanner(stdout)
	if !scanner.Scan() {
		t.Fatalf("owner process helper did not become ready: %v", scanner.Err())
	}
	if scanner.Text() != "ready" {
		t.Fatalf("owner process helper readiness = %q, want ready", scanner.Text())
	}

	waitStarted = true
	go func() {
		wait <- cmd.Wait()
	}()
	select {
	case err := <-wait:
		waited = true
		t.Fatalf("owner process helper exited after readiness: %v", err)
	default:
	}

	if err := cmd.Process.Signal(syscall.SIGTERM); err != nil {
		t.Fatalf("signal owner process helper with SIGTERM: %v", err)
	}
	if err := cmd.Process.Signal(syscall.SIGUSR1); err != nil {
		t.Fatalf("signal owner process helper liveness probe: %v", err)
	}
	if !scanner.Scan() {
		t.Fatalf("owner process helper did not acknowledge liveness: %v", scanner.Err())
	}
	if scanner.Text() != "alive" {
		t.Fatalf("owner process helper liveness = %q, want alive", scanner.Text())
	}
	select {
	case err := <-wait:
		waited = true
		t.Fatalf("owner process helper exited after SIGTERM: %v", err)
	default:
	}

	if err := cmd.Process.Kill(); err != nil {
		t.Fatalf("kill owner process helper: %v", err)
	}
	<-wait
	waited = true
	output := stderr.String()
	if strings.Contains(output, "fatal error") || strings.Contains(output, "all goroutines are asleep") {
		t.Fatalf("owner process helper emitted a runtime fatal error: %s", output)
	}
}

func TestOwnerProcessHelper(t *testing.T) {
	t.Parallel()
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
		signals := make(chan os.Signal, 1)
		signal.Notify(signals, syscall.SIGUSR1)
		defer signal.Stop(signals)
		fmt.Fprintln(os.Stdout, "ready")
		for range signals {
			fmt.Fprintln(os.Stdout, "alive")
		}
	case "tree-root":
		middle := exec.Command(os.Args[0], "-test.run=^TestOwnerProcessHelper$")
		middle.Env = ownerProcessHelperEnv("tree-middle")
		output, err := middle.Output()
		if err != nil {
			t.Fatalf("start outliving grandchild: %v", err)
		}
		outputScanner := bufio.NewScanner(bytes.NewReader(output))
		if !outputScanner.Scan() {
			t.Fatalf("outliving grandchild omitted its pid: %v", outputScanner.Err())
		}
		grandchildPID, err := strconv.Atoi(strings.TrimSpace(outputScanner.Text()))
		if err != nil {
			t.Fatalf("parse outliving grandchild pid %q: %v", output, err)
		}
		signals := make(chan os.Signal, 1)
		signal.Notify(signals, syscall.SIGTERM)
		defer signal.Stop(signals)
		fmt.Fprintf(os.Stdout, "%d\nready\n", grandchildPID)
		<-signals
	case "tree-middle":
		grandchild := exec.Command(os.Args[0], "-test.run=^TestOwnerProcessHelper$")
		grandchild.Env = ownerProcessHelperEnv("tree-grandchild")
		stdout, err := grandchild.StdoutPipe()
		if err != nil {
			t.Fatalf("open grandchild stdout: %v", err)
		}
		if err := grandchild.Start(); err != nil {
			t.Fatalf("start grandchild: %v", err)
		}
		scanner := bufio.NewScanner(stdout)
		if !scanner.Scan() || scanner.Text() != "ready" {
			t.Fatalf("grandchild readiness = %q, error = %v", scanner.Text(), scanner.Err())
		}
		fmt.Fprintln(os.Stdout, grandchild.Process.Pid)
		if err := grandchild.Process.Release(); err != nil {
			t.Fatalf("release grandchild: %v", err)
		}
	case "tree-grandchild":
		signals := make(chan os.Signal, 1)
		signal.Notify(signals, syscall.SIGTERM)
		defer signal.Stop(signals)
		fmt.Fprintln(os.Stdout, "ready")
		<-signals
	default:
		t.Fatalf("unknown owner process helper mode %q", mode)
	}
}

func startOwnerProcessHelper(t *testing.T, mode string) (int, <-chan error) {
	t.Helper()
	cmd := exec.Command(os.Args[0], "-test.run=^TestOwnerProcessHelper$")
	cmd.Env = ownerProcessHelperEnv(mode)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		t.Fatalf("open owner process helper stdout: %v", err)
	}
	cmd.Stderr = os.Stderr
	if err := cmd.Start(); err != nil {
		t.Fatalf("start owner process helper: %v", err)
	}
	wait := make(chan error, 1)
	waitStarted := false
	t.Cleanup(func() {
		if ProcessAlive(cmd.Process.Pid) {
			_ = cmd.Process.Kill()
		}
		if !waitStarted {
			_ = cmd.Wait()
			return
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
	waitStarted = true
	go func() {
		wait <- cmd.Wait()
		close(wait)
	}()
	return cmd.Process.Pid, wait
}

func startOwnerProcessTreeHelper(t *testing.T) (int, int, <-chan error) {
	t.Helper()
	cmd := exec.Command(os.Args[0], "-test.run=^TestOwnerProcessHelper$")
	cmd.Env = ownerProcessHelperEnv("tree-root")
	cmd.SysProcAttr = &syscall.SysProcAttr{Setsid: true}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		t.Fatalf("open owner process tree stdout: %v", err)
	}
	cmd.Stderr = os.Stderr
	if err := cmd.Start(); err != nil {
		t.Fatalf("start owner process tree: %v", err)
	}
	if sessionID, err := unix.Getsid(cmd.Process.Pid); err != nil || sessionID != cmd.Process.Pid {
		t.Fatalf("owner process tree session = %d, error = %v, want %d", sessionID, err, cmd.Process.Pid)
	}
	if groupID, err := unix.Getpgid(cmd.Process.Pid); err != nil || groupID != cmd.Process.Pid {
		t.Fatalf("owner process tree group = %d, error = %v, want %d", groupID, err, cmd.Process.Pid)
	}
	wait := make(chan error, 1)
	waitStarted := false
	t.Cleanup(func() {
		_ = syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL)
		if !waitStarted {
			_ = cmd.Wait()
			return
		}
		select {
		case <-wait:
		case <-time.After(2 * time.Second):
		}
	})
	scanner := bufio.NewScanner(stdout)
	if !scanner.Scan() {
		t.Fatalf("owner process tree omitted grandchild pid: %v", scanner.Err())
	}
	grandchildPID, err := strconv.Atoi(scanner.Text())
	if err != nil {
		t.Fatalf("parse owner process tree grandchild pid %q: %v", scanner.Text(), err)
	}
	if !scanner.Scan() || scanner.Text() != "ready" {
		t.Fatalf("owner process tree readiness = %q, error = %v", scanner.Text(), scanner.Err())
	}
	waitStarted = true
	go func() {
		wait <- cmd.Wait()
		close(wait)
	}()
	return cmd.Process.Pid, grandchildPID, wait
}

func ownerProcessHelperEnv(mode string) []string {
	const key = "ROUNDFIX_OWNER_PROCESS_HELPER"
	env := make([]string, 0, len(os.Environ())+1)
	for _, entry := range os.Environ() {
		name, _, _ := strings.Cut(entry, "=")
		if name != key {
			env = append(env, entry)
		}
	}
	return append(env, key+"="+mode)
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

func assertOwnerProcessForceKilled(t *testing.T, pid int, wait <-chan error) {
	t.Helper()
	select {
	case err := <-wait:
		var exitErr *exec.ExitError
		if !errors.As(err, &exitErr) {
			t.Fatalf("owner process %d exited prematurely before controller force-kill escalation: %v", pid, err)
		}
		status, ok := exitErr.ProcessState.Sys().(syscall.WaitStatus)
		if !ok || !status.Signaled() || status.Signal() != syscall.SIGKILL {
			t.Fatalf("owner process %d exit = %v, want controller force-kill signal %v", pid, err, syscall.SIGKILL)
		}
	case <-time.After(2 * time.Second):
		t.Fatalf("owner process %d did not exit after controller force-kill escalation", pid)
	}
	if ProcessAlive(pid) {
		t.Fatalf("owner process %d remained alive after controller force-kill escalation", pid)
	}
}
