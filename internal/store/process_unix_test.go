//go:build unix

package store

import (
	"os"
	"os/exec"
	"testing"
)

func TestProcessAliveReportsCurrentProcessAlive(t *testing.T) {
	if !ProcessAlive(os.Getpid()) {
		t.Fatal("expected current process to report alive")
	}
}

func TestProcessAliveReportsReapedChildDead(t *testing.T) {
	cmd := exec.Command("/bin/sh", "-c", "exit 0")
	if err := cmd.Start(); err != nil {
		t.Fatalf("start child process: %v", err)
	}
	pid := cmd.Process.Pid
	if err := cmd.Wait(); err != nil {
		t.Fatalf("wait for child process: %v", err)
	}
	if ProcessAlive(pid) {
		t.Fatalf("expected reaped child pid %d to report dead", pid)
	}
}
