//go:build windows

package store

import (
	"context"
	"errors"
	"os"
	"syscall"
)

const (
	processTerminate               = 0x0001
	processQueryLimitedInformation = 0x1000
	errorInvalidParameter          = syscall.Errno(87)
)

func processAbsent(pid int) (bool, error) {
	handle, err := syscall.OpenProcess(processQueryLimitedInformation, false, uint32(pid))
	if err != nil {
		if errors.Is(err, errorInvalidParameter) {
			return true, nil
		}
		return false, err
	}
	_ = syscall.CloseHandle(handle)
	return false, nil
}

func signalOwnerProcess(pid int, force bool) error {
	if !force {
		process, err := os.FindProcess(pid)
		if err != nil {
			return err
		}
		err = process.Signal(os.Interrupt)
		if errors.Is(err, os.ErrProcessDone) {
			return errOwnerProcessAlreadyAbsent
		}
		if errors.Is(err, syscall.EWINDOWS) {
			return errGracefulSignalUnsupported
		}
		return err
	}

	handle, err := syscall.OpenProcess(processTerminate, false, uint32(pid))
	if err != nil {
		if errors.Is(err, errorInvalidParameter) {
			return errOwnerProcessAlreadyAbsent
		}
		return err
	}
	defer func() {
		_ = syscall.CloseHandle(handle)
	}()
	if err := syscall.TerminateProcess(handle, 1); err != nil {
		return err
	}
	return nil
}

// processStartIdentity is unsupported on Windows: Runs created here record no
// identity token and keep the legacy PID-only owner proof.
func processStartIdentity(_ context.Context, _ int) (string, error) {
	return "", ErrOwnerProcessUnsupported
}
