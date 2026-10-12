//go:build unix

package cli

import (
	"context"
	"os/exec"
	"syscall"

	"roundfix/internal/store"
)

// configureDetachedSysProcAttr makes the detached child a session leader so it
// survives the caller's process group being reaped.
func configureDetachedSysProcAttr(cmd *exec.Cmd) {
	cmd.SysProcAttr = &syscall.SysProcAttr{Setsid: true}
}

// killDetachedProcessGroup kills the detached child's whole process group and
// reports whether it handled the kill.
func killDetachedProcessGroup(cmd *exec.Cmd, leaderIdentity string) bool {
	_, err := store.SignalProvenProcessGroup(context.Background(), cmd.Process.Pid, leaderIdentity, syscall.SIGKILL)
	return err == nil
}
