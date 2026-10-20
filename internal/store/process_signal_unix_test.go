//go:build unix

// Suite: Proven Signals.
// Invariant: an identity mismatch leaves the live target untouched.
// Boundary IN: test-owned cat processes, recorded start identities and injected owner readers.
// Boundary OUT: unrelated processes and daemon-owned Verification.
package store

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"syscall"
	"testing"

	"roundfix/internal/testwait"
)

type provenSignalHelper struct {
	command  *exec.Cmd
	identity string
	input    io.WriteCloser
	output   *bufio.Scanner
	reaped   bool
	wait     <-chan error
}

func startProvenSignalHelper(t *testing.T, group int) *provenSignalHelper {
	t.Helper()
	command := exec.Command("cat")
	command.SysProcAttr = &syscall.SysProcAttr{Setpgid: true, Pgid: group}
	input, err := command.StdinPipe()
	if err != nil {
		t.Fatal(err)
	}
	output, err := command.StdoutPipe()
	if err != nil {
		t.Fatal(err)
	}
	if err := command.Start(); err != nil {
		t.Fatal(err)
	}
	helper := &provenSignalHelper{command: command, input: input, output: bufio.NewScanner(output)}
	t.Cleanup(func() {
		_ = input.Close()
		if !helper.reaped {
			_ = command.Process.Kill()
			if helper.wait != nil {
				<-helper.wait
			} else {
				_ = command.Wait()
			}
		}
	})
	helper.identity, err = OwnerProcessIdentity(t.Context(), command.Process.Pid)
	if err != nil {
		t.Fatalf("record helper identity at start: %v", err)
	}
	assertProvenSignalHelperAlive(t, helper)
	return helper
}

func assertProvenSignalHelperAlive(t *testing.T, helper *provenSignalHelper) {
	t.Helper()
	if _, err := fmt.Fprintln(helper.input, "alive"); err != nil {
		t.Fatalf("write helper liveness probe: %v", err)
	}
	ready := make(chan bool, 1)
	go func() { ready <- helper.output.Scan() }()
	if !testwait.Until(t, "helper liveness response", ready, (<-chan error)(nil)) || helper.output.Text() != "alive" {
		t.Fatalf("helper liveness response = %q, error = %v", helper.output.Text(), helper.output.Err())
	}
}

func waitProvenSignalHelper(t *testing.T, helper *provenSignalHelper) {
	t.Helper()
	ended := make(chan error, 1)
	helper.wait = ended
	go func() { ended <- helper.command.Wait() }()
	err := testwait.Until(t, "signalled helper exit", ended, (<-chan error)(nil))
	helper.reaped = true
	var exit *exec.ExitError
	if !errors.As(err, &exit) || !exit.Sys().(syscall.WaitStatus).Signaled() {
		t.Fatalf("helper exit = %v, want signal exit", err)
	}
}

func TestSignalProvenProcessSendsToTheRecordedProcess(t *testing.T) {
	t.Parallel()
	helper := startProvenSignalHelper(t, 0)
	sent, err := SignalProvenProcess(t.Context(), helper.command.Process.Pid, helper.identity, syscall.SIGTERM)
	if err != nil || !sent {
		t.Fatalf("signal recorded helper: sent=%v error=%v", sent, err)
	}
	waitProvenSignalHelper(t, helper)
	// A reaped target is an ordinary no-send result.
	sent, err = SignalProvenProcess(t.Context(), helper.command.Process.Pid, helper.identity, syscall.SIGTERM)
	if err != nil || sent {
		t.Fatalf("signal gone helper: sent=%v error=%v", sent, err)
	}
}

func TestSignalProvenProcessSendsNothingToAReusedPID(t *testing.T) {
	t.Parallel()
	helper := startProvenSignalHelper(t, 0)
	sent, err := SignalProvenProcess(t.Context(), helper.command.Process.Pid, helper.identity+"-old", syscall.SIGKILL)
	if err != nil || sent {
		t.Fatalf("signal mismatched helper: sent=%v error=%v", sent, err)
	}
	assertProvenSignalHelperAlive(t, helper)
}

func TestSignalProvenProcessRefusesAnUnprovableTarget(t *testing.T) {
	t.Parallel()
	helper := startProvenSignalHelper(t, 0)
	for _, tc := range []struct {
		name     string
		pid      int
		identity string
		signal   syscall.Signal
	}{
		{"zero PID", 0, helper.identity, syscall.SIGKILL},
		{"negative PID", -helper.command.Process.Pid, helper.identity, syscall.SIGKILL},
		{"self", os.Getpid(), helper.identity, syscall.SIGKILL},
		{"no identity", helper.command.Process.Pid, "", syscall.SIGKILL},
		{"blank identity", helper.command.Process.Pid, "  ", syscall.SIGKILL},
		{"zero signal", helper.command.Process.Pid, helper.identity, 0},
	} {
		t.Run(tc.name, func(t *testing.T) {
			for _, operation := range []func(context.Context, int, string, syscall.Signal) (bool, error){SignalProvenProcess, SignalProvenProcessGroup} {
				sent, err := operation(t.Context(), tc.pid, tc.identity, tc.signal)
				if sent || !errors.Is(err, ErrOwnerProcessIdentityUnproven) {
					t.Fatalf("unprovable target: sent=%v error=%v", sent, err)
				}
			}
		})
	}
	assertProvenSignalHelperAlive(t, helper)
}

func TestSignalProvenProcessGroupProvesItsLeader(t *testing.T) {
	t.Parallel()
	leader := startProvenSignalHelper(t, 0)
	member := startProvenSignalHelper(t, leader.command.Process.Pid)
	sent, err := SignalProvenProcessGroup(t.Context(), leader.command.Process.Pid, leader.identity+"-old", syscall.SIGKILL)
	if err != nil || sent {
		t.Fatalf("signal mismatched group: sent=%v error=%v", sent, err)
	}
	assertProvenSignalHelperAlive(t, leader)
	assertProvenSignalHelperAlive(t, member)
	sent, err = SignalProvenProcessGroup(t.Context(), leader.command.Process.Pid, leader.identity, syscall.SIGKILL)
	if err != nil || !sent {
		t.Fatalf("signal proven group: sent=%v error=%v", sent, err)
	}
	waitProvenSignalHelper(t, leader)
	waitProvenSignalHelper(t, member)
}

func TestSignalProvenProcessGroupSignalsMembersOfAGoneLeader(t *testing.T) {
	t.Parallel()
	leader := startProvenSignalHelper(t, 0)
	member := startProvenSignalHelper(t, leader.command.Process.Pid)
	if err := leader.command.Process.Kill(); err != nil {
		t.Fatal(err)
	}
	waitProvenSignalHelper(t, leader)
	assertProvenSignalHelperAlive(t, member)
	sent, err := SignalProvenProcessGroup(t.Context(), leader.command.Process.Pid, leader.identity, syscall.SIGKILL)
	if err != nil || !sent {
		t.Fatalf("signal members of gone leader: sent=%v error=%v", sent, err)
	}
	waitProvenSignalHelper(t, member)
	sent, err = SignalProvenProcessGroup(t.Context(), leader.command.Process.Pid, leader.identity, syscall.SIGKILL)
	if err != nil || sent {
		t.Fatalf("signal empty group: sent=%v error=%v", sent, err)
	}
}

func TestOwnerProcessControllerReprovesBeforeForceKill(t *testing.T) {
	t.Parallel()
	for _, name := range []string{"changed", "absent", "unreadable", "matching", "legacy"} {
		t.Run(name, func(t *testing.T) {
			helper := startProvenSignalHelper(t, 0)
			// The stop window only bounds the wait: a short one races the grace
			// timer on a loaded runner and fails before re-proof is reached.
			controller := newOwnerProcessController(ownerProcessTestGracePeriod, ownerProcessTestWaitBudget, ownerProcessTestPollInterval)
			recorded := helper.identity
			signalled, forced := false, false
			hostErr := errors.New("identity reader unavailable")
			controller.processAbsent = func(int) (bool, error) { return forced || (signalled && name == "absent"), nil }
			controller.processStartIdentity = func(context.Context, int) (string, error) {
				if signalled {
					if name == "changed" {
						return recorded + "-stranger", nil
					}
					if name == "unreadable" {
						return "", hostErr
					}
				}
				return recorded, nil
			}
			controller.signalProcess = func(_ int, force bool) error {
				signalled = true
				forced = force
				if force {
					return helper.command.Process.Kill()
				}
				return nil
			}
			identity := recorded
			if name == "legacy" {
				identity = ""
			}
			err := controller.TerminateAndWait(t.Context(), helper.command.Process.Pid, identity)
			if name == "unreadable" {
				if !errors.Is(err, hostErr) || !errors.Is(err, ErrOwnerIdentityUnreadable) {
					t.Fatalf("re-proof read failure = %v", err)
				}
			} else if err != nil {
				t.Fatalf("terminate after grace: %v", err)
			}
			if want := name == "matching" || name == "legacy"; forced != want {
				t.Fatalf("forced signal = %v, want %v", forced, want)
			}
			if !signalled {
				t.Fatal("graceful signal missing")
			}
			if forced {
				waitProvenSignalHelper(t, helper)
			} else {
				assertProvenSignalHelperAlive(t, helper)
			}
		})
	}
}
