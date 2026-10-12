//go:build !windows

package daemon

import (
	"context"
	"os"
	"os/exec"
	"sync"
	"syscall"

	"roundfix/internal/store"
)

// Kill the shell and its formatter children before restoring their files.
// The returned function publishes the leader identity immediately after Start;
// cancellation waits for that publication, including when Start fails.
func configureQAFormatProcess(command *exec.Cmd) func() {
	command.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	var lock sync.Mutex
	var identity string
	var identityErr error
	lock.Lock()
	command.Cancel = func() error {
		lock.Lock()
		defer lock.Unlock()
		if identityErr != nil {
			return identityErr
		}
		sent, err := store.SignalProvenProcessGroup(context.Background(), command.Process.Pid, identity, syscall.SIGKILL)
		if err == nil && !sent {
			return os.ErrProcessDone
		}
		return err
	}
	return func() {
		defer lock.Unlock()
		if command.Process != nil {
			identity, identityErr = store.OwnerProcessIdentity(context.Background(), command.Process.Pid)
		}
	}
}
