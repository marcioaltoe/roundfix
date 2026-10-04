//go:build darwin

package cli

import (
	"fmt"

	"golang.org/x/sys/unix"
)

// detachFixtureZombieState is SZOMB in Darwin's <sys/proc.h>.
const detachFixtureZombieState = 5

// detachFixtureGroupLiveMembers excludes exited, unreaped group members.
// An error is a failed reading, never an empty group.
func detachFixtureGroupLiveMembers(pgid int) ([]int, error) {
	group, err := unix.SysctlKinfoProcSlice("kern.proc.pgrp", pgid)
	if err != nil {
		return nil, fmt.Errorf("read fixture process group %d: %w", pgid, err)
	}
	var live []int
	for _, process := range group {
		if process.Proc.P_stat != detachFixtureZombieState {
			live = append(live, int(process.Proc.P_pid))
		}
	}
	return live, nil
}
