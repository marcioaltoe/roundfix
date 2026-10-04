//go:build linux

package cli

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
)

// detachFixtureGroupLiveMembers excludes exited, unreaped group members.
// An error is a failed reading, never an empty group.
func detachFixtureGroupLiveMembers(pgid int) ([]int, error) {
	entries, err := os.ReadDir("/proc")
	if err != nil {
		return nil, fmt.Errorf("read process table for fixture group %d: %w", pgid, err)
	}
	var live []int
	for _, entry := range entries {
		pid, err := strconv.Atoi(entry.Name())
		if err != nil || pid <= 0 || !entry.IsDir() {
			continue
		}
		stat, err := os.ReadFile(filepath.Join("/proc", entry.Name(), "stat"))
		if errors.Is(err, syscall.ENOENT) || errors.Is(err, syscall.ESRCH) {
			continue
		}
		if err != nil {
			return nil, fmt.Errorf("read process %d for fixture group %d: %w", pid, pgid, err)
		}
		// comm can contain spaces and ')'; state, ppid and pgrp follow
		// the last ')', at indices 0, 1 and 2 respectively.
		end := strings.LastIndexByte(string(stat), ')')
		if end < 0 {
			return nil, fmt.Errorf("read process %d for fixture group %d: missing comm delimiter", pid, pgid)
		}
		fields := strings.Fields(string(stat[end+1:]))
		if len(fields) < 3 || len(fields[0]) != 1 {
			return nil, fmt.Errorf("read process %d for fixture group %d: incomplete stat", pid, pgid)
		}
		group, err := strconv.Atoi(fields[2])
		if err != nil {
			return nil, fmt.Errorf("read process %d group for fixture group %d: %w", pid, pgid, err)
		}
		if group == pgid && fields[0] != "Z" {
			live = append(live, pid)
		}
	}
	return live, nil
}
