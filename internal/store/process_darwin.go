//go:build darwin

package store

import (
	"context"
	"errors"
	"fmt"
	"unsafe"

	"golang.org/x/sys/unix"
)

func processStartIdentity(_ context.Context, pid int) (string, error) {
	raw, err := unix.SysctlRaw("kern.proc.pid", pid)
	if err != nil {
		return "", fmt.Errorf("read start time for process %d: %w", pid, err)
	}
	if len(raw) == 0 {
		return "", fmt.Errorf("read start time for process %d: %w", pid, unix.ESRCH)
	}
	if len(raw) != unix.SizeofKinfoProc {
		return "", fmt.Errorf("read start time for process %d: sysctl returned %d bytes, want %d", pid, len(raw), unix.SizeofKinfoProc)
	}
	process := *(*unix.KinfoProc)(unsafe.Pointer(&raw[0]))
	started := process.Proc.P_starttime
	return fmt.Sprintf("darwin:%d.%d", started.Sec, started.Usec), nil
}

func processTreePIDs(ownerPID int, _ string) ([]int, error) {
	group, err := unix.SysctlKinfoProcSlice("kern.proc.pgrp", ownerPID)
	if err != nil {
		return nil, fmt.Errorf("read owner process group: %w", err)
	}
	sessionPIDs := make([]int, 0, len(group))
	for _, process := range group {
		pid := int(process.Proc.P_pid)
		if pid <= 0 {
			continue
		}
		sessionID, err := unix.Getsid(pid)
		if err != nil {
			if errors.Is(err, unix.ESRCH) {
				continue
			}
			return nil, fmt.Errorf("read process %d session: %w", pid, err)
		}
		if sessionID == ownerPID {
			sessionPIDs = append(sessionPIDs, pid)
		}
	}
	if len(sessionPIDs) > 0 {
		return sessionPIDs, nil
	}
	if len(group) == 0 {
		absent, absentErr := processAbsent(ownerPID)
		if absentErr == nil && absent {
			return nil, nil
		}
		if absentErr != nil {
			return nil, fmt.Errorf("prove owner absence after empty process group query: %w", absentErr)
		}
	}

	processes, err := unix.SysctlKinfoProcSlice("kern.proc.all")
	if err != nil {
		return nil, fmt.Errorf("read process table for non-session owner: %w", err)
	}
	parents := make([]processParent, 0, len(processes))
	for _, process := range processes {
		pid := int(process.Proc.P_pid)
		if pid > 0 {
			parents = append(parents, processParent{pid: pid, parentPID: int(process.Eproc.Ppid)})
		}
	}
	return descendantProcessPIDs(ownerPID, parents), nil
}
