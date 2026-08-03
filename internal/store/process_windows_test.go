//go:build windows

package store

import (
	"os/exec"
	"syscall"
	"testing"
)

const synchronizeProcess = 0x00100000

func TestProcessAliveReportsExitedUnreapedChildDead(t *testing.T) {
	t.Parallel()
	cmd := exec.Command("cmd", "/c", "exit", "0")
	if err := cmd.Start(); err != nil {
		t.Fatalf("start child process: %v", err)
	}
	defer func() {
		if err := cmd.Wait(); err != nil {
			t.Fatalf("wait for child process: %v", err)
		}
	}()

	handle, err := syscall.OpenProcess(
		synchronizeProcess|processQueryLimitedInformation,
		false,
		uint32(cmd.Process.Pid),
	)
	if err != nil {
		t.Fatalf("open child process: %v", err)
	}
	defer func() {
		if err := syscall.CloseHandle(handle); err != nil {
			t.Fatalf("close child process handle: %v", err)
		}
	}()
	event, err := syscall.WaitForSingleObject(handle, 5_000)
	if err != nil {
		t.Fatalf("wait for child process exit: %v", err)
	}
	if event != syscall.WAIT_OBJECT_0 {
		t.Fatalf("child process wait event = %#x, want WAIT_OBJECT_0", event)
	}

	if ProcessAlive(cmd.Process.Pid) {
		t.Fatalf("expected exited unreaped child pid %d to report dead", cmd.Process.Pid)
	}
}
