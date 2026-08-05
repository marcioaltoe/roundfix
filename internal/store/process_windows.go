//go:build windows

package store

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
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

func processStartIdentity(_ context.Context, pid int) (string, error) {
	started, err := processCreationTime(pid)
	if err != nil {
		return "", fmt.Errorf("read start time for process %d: %w", pid, err)
	}
	return fmt.Sprintf("windows:%d", started), nil
}

func processTreePIDs(ownerPID int, recordedIdentity string) ([]int, error) {
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

	ownerStarted, comparable := windowsProcessStart(recordedIdentity)
	if !comparable {
		ownerStarted, err = processCreationTime(ownerPID)
		if err != nil {
			return nil, fmt.Errorf("read owner process creation time: %w", err)
		}
	}
	candidates := descendantProcessPIDs(ownerPID, parents)
	started := make([]processParentWithStart, 0, len(candidates)-1)
	parentsByPID := make(map[int]int, len(parents))
	for _, process := range parents {
		parentsByPID[process.pid] = process.parentPID
	}
	for _, pid := range candidates {
		if pid == ownerPID {
			continue
		}
		created, creationErr := processCreationTime(pid)
		if errors.Is(creationErr, errorInvalidParameter) {
			continue
		}
		if creationErr != nil {
			return nil, fmt.Errorf("read process %d creation time: %w", pid, creationErr)
		}
		started = append(started, processParentWithStart{
			pid:       pid,
			parentPID: parentsByPID[pid],
			started:   created,
		})
	}
	return descendantProcessPIDsAfterStart(ownerPID, ownerStarted, started), nil
}

func processCreationTime(pid int) (uint64, error) {
	handle, err := windows.OpenProcess(processQueryLimitedInformation, false, uint32(pid))
	if err != nil {
		return 0, err
	}
	defer func() {
		_ = windows.CloseHandle(handle)
	}()
	var created, exited, kernel, user windows.Filetime
	if err := windows.GetProcessTimes(handle, &created, &exited, &kernel, &user); err != nil {
		return 0, err
	}
	return uint64(created.HighDateTime)<<32 | uint64(created.LowDateTime), nil
}

func windowsProcessStart(identity string) (uint64, bool) {
	value, ok := strings.CutPrefix(strings.TrimSpace(identity), "windows:")
	if !ok || value == "" {
		return 0, false
	}
	started, err := strconv.ParseUint(value, 10, 64)
	return started, err == nil
}
