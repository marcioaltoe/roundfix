//go:build linux

package store

import (
	"bytes"
	"context"
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
