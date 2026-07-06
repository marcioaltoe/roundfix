//go:build windows

package cli

import (
	"os/exec"
	"syscall"
)

// Windows process-creation flags (from processthreadsapi.h) that give the
// detached child its own process group and detach it from the caller's console
// so it survives the caller.
const (
	createNewProcessGroup = 0x00000200
	detachedProcess       = 0x00000008
)

func configureDetachedSysProcAttr(cmd *exec.Cmd) {
	cmd.SysProcAttr = &syscall.SysProcAttr{CreationFlags: createNewProcessGroup | detachedProcess}
}

// killDetachedProcessGroup has no syscall-level process-group kill on Windows;
// return false so the caller falls back to Process.Kill.
func killDetachedProcessGroup(_ *exec.Cmd) bool {
	return false
}
