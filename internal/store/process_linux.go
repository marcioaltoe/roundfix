//go:build linux

package store

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"strconv"
)

func processStartIdentity(_ context.Context, pid int) (string, error) {
	stat, err := os.ReadFile(fmt.Sprintf("/proc/%d/stat", pid))
	if err != nil {
		return "", fmt.Errorf("read start time for process %d: %w", pid, err)
	}
	startTime, err := parseProcStatStartTime(stat)
	if err != nil {
		return "", fmt.Errorf("read start time for process %d: %w", pid, err)
	}
	return "linux:" + startTime, nil
}

func parseProcStatStartTime(stat []byte) (string, error) {
	commEnd := bytes.LastIndexByte(stat, ')')
	if commEnd < 0 {
		return "", fmt.Errorf("malformed proc stat: missing comm terminator")
	}
	fields := bytes.Fields(stat[commEnd+1:])
	const startTimeIndexAfterComm = 19
	if len(fields) <= startTimeIndexAfterComm {
		return "", fmt.Errorf("malformed proc stat: missing start time field")
	}
	startTime := string(fields[startTimeIndexAfterComm])
	if _, err := strconv.ParseUint(startTime, 10, 64); err != nil {
		return "", fmt.Errorf("malformed proc stat start time %q: %w", startTime, err)
	}
	return startTime, nil
}

func processTreePIDs(ownerPID int) ([]int, error) {
	entries, err := os.ReadDir("/proc")
	if err != nil {
		return nil, fmt.Errorf("read process table: %w", err)
	}
	parents := make([]processParent, 0, len(entries))
	sessionPIDs := make([]int, 0)
	for _, entry := range entries {
		pid, err := strconv.Atoi(entry.Name())
		if err != nil || pid <= 0 {
			continue
		}
		stat, err := os.ReadFile(fmt.Sprintf("/proc/%d/stat", pid))
		if err != nil {
			if errors.Is(err, os.ErrNotExist) {
				continue
			}
			return nil, fmt.Errorf("read process %d ownership: %w", pid, err)
		}
		parentPID, sessionID, err := parseProcStatOwnership(stat)
		if err != nil {
			return nil, fmt.Errorf("read process %d ownership: %w", pid, err)
		}
		parents = append(parents, processParent{pid: pid, parentPID: parentPID})
		if sessionID == ownerPID {
			sessionPIDs = append(sessionPIDs, pid)
		}
	}
	if len(sessionPIDs) > 0 {
		return sessionPIDs, nil
	}
	return descendantProcessPIDs(ownerPID, parents), nil
}

func parseProcStatOwnership(stat []byte) (int, int, error) {
	commEnd := bytes.LastIndexByte(stat, ')')
	if commEnd < 0 {
		return 0, 0, fmt.Errorf("malformed proc stat: missing comm terminator")
	}
	fields := bytes.Fields(stat[commEnd+1:])
	const sessionIndexAfterComm = 3
	if len(fields) <= sessionIndexAfterComm {
		return 0, 0, fmt.Errorf("malformed proc stat: missing ownership fields")
	}
	parentPID, err := strconv.Atoi(string(fields[1]))
	if err != nil {
		return 0, 0, fmt.Errorf("malformed proc stat parent pid %q: %w", fields[1], err)
	}
	sessionID, err := strconv.Atoi(string(fields[sessionIndexAfterComm]))
	if err != nil {
		return 0, 0, fmt.Errorf("malformed proc stat session id %q: %w", fields[sessionIndexAfterComm], err)
	}
	return parentPID, sessionID, nil
}
