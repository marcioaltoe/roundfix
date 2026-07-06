//go:build unix

package cli

import (
	"errors"
	"os/exec"
	"syscall"
)

// configureDetachedSysProcAttr makes the detached child a session leader so it
// survives the caller's process group being reaped.
func configureDetachedSysProcAttr(cmd *exec.Cmd) {
	cmd.SysProcAttr = &syscall.SysProcAttr{Setsid: true}
}

// killDetachedProcessGroup kills the detached child's whole process group and
// reports whether it handled the kill.
func killDetachedProcessGroup(cmd *exec.Cmd) bool {
	err := syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL)
	return err == nil || errors.Is(err, syscall.ESRCH)
}
