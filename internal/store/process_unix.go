//go:build unix

package store

import (
	"context"
	"errors"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
	"syscall"
)

func processAbsent(pid int) (bool, error) {
	err := syscall.Kill(pid, syscall.Signal(0))
	if err == nil {
		return false, nil
	}
	if errors.Is(err, syscall.ESRCH) {
		return true, nil
	}
	return false, err
}

func signalOwnerProcess(pid int, force bool) error {
	signal := syscall.SIGTERM
	if force {
		signal = syscall.SIGKILL
	}
	err := syscall.Kill(pid, signal)
	if errors.Is(err, syscall.ESRCH) {
		return errOwnerProcessAlreadyAbsent
	}
	return err
}

// processStartIdentity returns the process start time exactly as ps prints
// it. The verbatim string is the opaque identity token: two processes reusing
// one PID cannot share it, and equality is the only supported comparison.
func processStartIdentity(ctx context.Context, pid int) (string, error) {
	output, err := exec.CommandContext(ctx, "ps", "-p", strconv.Itoa(pid), "-o", "lstart=").Output()
	if err != nil {
		return "", fmt.Errorf("read start time for process %d: %w", pid, err)
	}
	identity := strings.TrimSpace(string(output))
	if identity == "" {
		return "", fmt.Errorf("read start time for process %d: ps reported no start time", pid)
	}
	return identity, nil
}
