// Package childexit describes how a child process ended.
package childexit

import (
	"errors"
	"fmt"
	"os/exec"
	"syscall"
)

// Exit is how a child process ended.
type Exit struct {
	Code       int    // The exit status, or -1 when a signal ended the child.
	Signal     string // The signal's name; empty when the child exited.
	Number     int    // The signal's number; zero when the child exited.
	CoreDumped bool
}

// From reads the end of a child from the error Run or Wait returned.
// A nil error is an exit with code 0; ok is false for any other error type.
func From(err error) (exit Exit, ok bool) {
	if err == nil {
		return Exit{}, true
	}
	var exitErr *exec.ExitError
	if !errors.As(err, &exitErr) {
		return Exit{}, false
	}
	if status, ok := exitErr.Sys().(syscall.WaitStatus); ok {
		return fromWaitStatus(status), true
	}
	return Exit{Code: exitErr.ExitCode()}, true
}

func fromWaitStatus(status syscall.WaitStatus) Exit {
	exit := Exit{Code: status.ExitStatus()}
	if status.Signaled() {
		signal := status.Signal()
		exit.Signal = signalName(signal)
		exit.Number = int(signal)
		exit.CoreDumped = status.CoreDump()
	}
	return exit
}

// Signaled reports whether a signal ended the child.
func (exit Exit) Signaled() bool {
	return exit.Number != 0
}

// String describes the exit status or the signal and core dump state.
func (exit Exit) String() string {
	if !exit.Signaled() {
		return fmt.Sprintf("exited %d", exit.Code)
	}
	core := "no core dump"
	if exit.CoreDumped {
		core = "core dumped"
	}
	return fmt.Sprintf("was killed by signal %s (%d), %s", exit.Signal, exit.Number, core)
}
