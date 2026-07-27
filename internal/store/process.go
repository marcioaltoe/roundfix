package store

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"
)

const (
	ownerProcessGracePeriod  = 2 * time.Second
	ownerProcessStopWindow   = 10 * time.Second
	ownerProcessPollInterval = 25 * time.Millisecond
)

var (
	ErrOwnerProcessUnsupported      = errors.New("owner process control is unsupported")
	ErrOwnerProcessIdentityUnproven = errors.New("owner process identity cannot be proven")
	errOwnerProcessAlreadyAbsent    = errors.New("owner process is already absent")
	errGracefulSignalUnsupported    = errors.New("graceful owner process termination is unsupported")
)

type OwnerProcessControlError struct {
	PID  int
	Step string
	Err  error
}

func (err OwnerProcessControlError) Error() string {
	return fmt.Sprintf("owner process %d failed step %q: %v", err.PID, err.Step, err.Err)
}

func (err OwnerProcessControlError) Unwrap() error {
	return err.Err
}

type OwnerProcessControl struct {
	gracePeriod  time.Duration
	stopWindow   time.Duration
	pollInterval time.Duration
}

func NewOwnerProcessController() *OwnerProcessControl {
	return newOwnerProcessController(
		ownerProcessGracePeriod,
		ownerProcessStopWindow,
		ownerProcessPollInterval,
	)
}

func newOwnerProcessController(gracePeriod, stopWindow, pollInterval time.Duration) *OwnerProcessControl {
	return &OwnerProcessControl{
		gracePeriod:  gracePeriod,
		stopWindow:   stopWindow,
		pollInterval: pollInterval,
	}
}

// TerminateAndWait proves ownership of pid before sending any signal: when
// the Run recorded a start-time identity token, the live process must present
// the same token, otherwise a reused PID would let Force Stop terminate an
// unrelated process. An empty recorded token comes from a Run created before
// identity recording existed and keeps the legacy PID-only proof, mirroring
// the ADR-0044 precedent that PID-less legacy Runs degrade gracefully.
func (controller *OwnerProcessControl) TerminateAndWait(ctx context.Context, pid int, recordedIdentity string) error {
	if pid <= 0 || pid == os.Getpid() {
		return ownerProcessControlError(pid, "prove owner process identity", ErrOwnerProcessIdentityUnproven)
	}
	absent, err := processAbsent(pid)
	if err != nil {
		return ownerProcessControlError(pid, "prove owner process identity", err)
	}
	if absent {
		return nil
	}
	if recorded := strings.TrimSpace(recordedIdentity); recorded != "" {
		liveIdentity, identityErr := processStartIdentity(ctx, pid)
		if identityErr != nil {
			// The owner may have exited between the liveness check and the
			// identity read; proven absence is the proof Force Stop needs.
			if absent, absentErr := processAbsent(pid); absentErr == nil && absent {
				return nil
			}
			return ownerProcessControlError(pid, "prove owner process identity",
				fmt.Errorf("%w: read live owner identity: %v", ErrOwnerProcessIdentityUnproven, identityErr))
		}
		if liveIdentity != recorded {
			return ownerProcessControlError(pid, "prove owner process identity",
				fmt.Errorf("%w: live process start identity does not match the recorded owner identity", ErrOwnerProcessIdentityUnproven))
		}
	}

	stopCtx, cancel := context.WithTimeout(ctx, controller.stopWindow)
	defer cancel()

	gracefulErr := signalOwnerProcess(pid, false)
	switch {
	case errors.Is(gracefulErr, errOwnerProcessAlreadyAbsent):
		return nil
	case gracefulErr == nil:
		absent, err = controller.waitForAbsence(stopCtx, pid, controller.gracePeriod)
		if err != nil {
			return ownerProcessControlError(pid, "prove exit after graceful termination", err)
		}
		if absent {
			return nil
		}
	case errors.Is(gracefulErr, errGracefulSignalUnsupported):
	default:
		return ownerProcessControlError(pid, "send graceful termination", gracefulErr)
	}

	if err := signalOwnerProcess(pid, true); err != nil {
		if errors.Is(err, errOwnerProcessAlreadyAbsent) {
			return nil
		}
		return ownerProcessControlError(pid, "send force kill", err)
	}
	absent, err = controller.waitForAbsence(stopCtx, pid, 0)
	if err != nil {
		return ownerProcessControlError(pid, "prove exit after force kill", err)
	}
	if !absent {
		return ownerProcessControlError(pid, "prove exit after force kill", context.DeadlineExceeded)
	}
	return nil
}

func (controller *OwnerProcessControl) waitForAbsence(ctx context.Context, pid int, window time.Duration) (bool, error) {
	pollInterval := controller.pollInterval
	if pollInterval <= 0 {
		pollInterval = ownerProcessPollInterval
	}
	ticker := time.NewTicker(pollInterval)
	defer ticker.Stop()

	var timer *time.Timer
	var windowElapsed <-chan time.Time
	if window > 0 {
		timer = time.NewTimer(window)
		windowElapsed = timer.C
		defer timer.Stop()
	}

	for {
		absent, err := processAbsent(pid)
		if err != nil {
			return false, err
		}
		if absent {
			return true, nil
		}
		select {
		case <-ctx.Done():
			return false, ctx.Err()
		case <-windowElapsed:
			return false, nil
		case <-ticker.C:
		}
	}
}

func ownerProcessControlError(pid int, step string, err error) error {
	return OwnerProcessControlError{PID: pid, Step: step, Err: err}
}

// ProcessAlive reports whether pid is alive or cannot be proven absent.
func ProcessAlive(pid int) bool {
	if pid <= 0 {
		return true
	}
	absent, err := processAbsent(pid)
	return err != nil || !absent
}

// OwnerProcessIdentity returns the opaque start-time identity token for pid.
// Callers compare tokens verbatim and must not parse them.
func OwnerProcessIdentity(ctx context.Context, pid int) (string, error) {
	if pid <= 0 {
		return "", fmt.Errorf("read owner process identity: pid %d is invalid", pid)
	}
	return processStartIdentity(ctx, pid)
}
