//go:build windows

package store

import (
	"context"
	"errors"
	"fmt"
	"os"
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

const (
	processTerminate               = 0x0001
	processQueryLimitedInformation = 0x1000
	errorInvalidParameter          = syscall.Errno(87)
	stillActive                    = 259
)

func processAbsent(pid int) (bool, error) {
	handle, err := syscall.OpenProcess(processQueryLimitedInformation, false, uint32(pid))
	if err != nil {
		if errors.Is(err, errorInvalidParameter) {
			return true, nil
		}
		return false, err
	}
	defer func() {
		_ = syscall.CloseHandle(handle)
	}()
	var exitCode uint32
	if err := syscall.GetExitCodeProcess(handle, &exitCode); err != nil {
		return false, err
	}
	return exitCode != stillActive, nil
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

func processTreePIDs(ownerPID int) ([]int, error) {
	snapshot, err := windows.CreateToolhelp32Snapshot(windows.TH32CS_SNAPPROCESS, 0)
	if err != nil {
		return nil, fmt.Errorf("read process table: %w", err)
	}
	defer func() {
		_ = windows.CloseHandle(snapshot)
	}()

	entry := windows.ProcessEntry32{Size: uint32(unsafe.Sizeof(windows.ProcessEntry32{}))}
	if err := windows.Process32First(snapshot, &entry); err != nil {
		return nil, fmt.Errorf("read first process table entry: %w", err)
	}
	parents := make([]processParent, 0)
	for {
		parents = append(parents, processParent{
			pid:       int(entry.ProcessID),
			parentPID: int(entry.ParentProcessID),
		})
		err := windows.Process32Next(snapshot, &entry)
		if errors.Is(err, windows.ERROR_NO_MORE_FILES) {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("read process table entry: %w", err)
		}
	}
	return descendantProcessPIDs(ownerPID, parents), nil
}
