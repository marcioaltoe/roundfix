//go:build darwin

package cli

import (
	"errors"
	"fmt"
	"os/exec"
	"syscall"
	"testing"
	"time"

	"golang.org/x/sys/unix"

	"roundfix/internal/testwait"
)

// detachFixtureZombieState is SZOMB in Darwin's <sys/proc.h>.
const detachFixtureZombieState = 5

// detachFixtureExitingFlag is P_WEXIT in Darwin's <sys/proc.h>.
const detachFixtureExitingFlag = 0x2000

// detachFixtureGroupLiveMembers excludes exiting and exited, unreaped members.
// An error is a failed reading, never an empty group.
func detachFixtureGroupLiveMembers(pgid int) ([]int, error) {
	group, err := unix.SysctlKinfoProcSlice("kern.proc.pgrp", pgid)
	if err != nil {
		return nil, fmt.Errorf("read fixture process group %d: %w", pgid, err)
	}
	var live []int
	for _, process := range group {
		if process.Proc.P_stat != detachFixtureZombieState && process.Proc.P_flag&detachFixtureExitingFlag == 0 {
			live = append(live, int(process.Proc.P_pid))
		}
	}
	return live, nil
}

func TestDetachFixtureGroupWithOnlyAnExitingMemberHasEnded(t *testing.T) {
	t.Parallel()
	deadline := time.Now().Add(testwait.Bound(t))
	launches := 0
	testwait.Poll(t, "non-zombie exiting fixture group with EPERM", (<-chan error)(nil), func() (bool, string) {
		if !time.Now().Before(deadline) {
			return false, fmt.Sprintf("deadline reached after %d launches", launches)
		}
		cmd := exec.Command("/usr/bin/true")
		cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
		if err := cmd.Start(); err != nil {
			t.Fatalf("start short-lived fixture: %v", err)
		}
		launches++
		pid := cmd.Process.Pid
		// Keep the child unreaped during the probe. This defer also reaps it
		// when an assertion fails or the polling deadline expires.
		defer func() {
			if err := cmd.Process.Kill(); err != nil && !errors.Is(err, syscall.ESRCH) {
				t.Errorf("stop short-lived fixture %d: %v", pid, err)
			}
			if err := cmd.Wait(); err != nil {
				var exitErr *exec.ExitError
				if !errors.As(err, &exitErr) || !exitErr.ProcessState.Sys().(syscall.WaitStatus).Signaled() {
					t.Errorf("reap short-lived fixture %d: %v", pid, err)
				}
			}
		}()
		// A ticker can miss the brief exiting window. Read continuously for
		// this child until it becomes a zombie, under the testwait deadline.
		for time.Now().Before(deadline) {
			group, err := unix.SysctlKinfoProcSlice("kern.proc.pgrp", pid)
			if err != nil {
				t.Fatalf("read raw fixture group %d: %v", pid, err)
			}
			if len(group) != 1 || int(group[0].Proc.P_pid) != pid {
				t.Fatalf("raw fixture group %d: want only child, got %v", pid, group)
			}
			process := group[0].Proc
			if process.P_stat == detachFixtureZombieState {
				return false, fmt.Sprintf("launch %d: fixture %d became a zombie before the window was observed", launches, pid)
			}
			if process.P_flag&detachFixtureExitingFlag == 0 {
				continue
			}
			probeErr := syscall.Kill(-pid, 0)
			if !errors.Is(probeErr, syscall.EPERM) {
				if probeErr != nil {
					t.Fatalf("probe exiting fixture group %d: %v", pid, probeErr)
				}
				continue
			}
			live, err := detachFixtureGroupLiveMembers(pid)
			if err != nil || len(live) != 0 {
				t.Fatalf("exiting fixture group %d: state=%d flags=%#x probe=%v live members=%v reading error=%v; want no live member", pid, process.P_stat, process.P_flag, probeErr, live, err)
			}
			t.Logf("observed exiting fixture group after %d launches: pid=%d state=%d flags=%#x probe=%v", launches, pid, process.P_stat, process.P_flag, probeErr)
			return true, "exiting group has no live member"
		}
		return false, fmt.Sprintf("deadline reached after %d launches", launches)
	})
}
