//go:build !windows

package worktree

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"golang.org/x/sys/unix"
)

func lockWorktreeAdminFile(ctx context.Context, file *os.File) error {
	timer := time.NewTimer(worktreeAdminLockPollInterval)
	defer timer.Stop()
	for {
		err := unix.Flock(int(file.Fd()), unix.LOCK_EX|unix.LOCK_NB)
		if err == nil {
			return nil
		}
		if !errors.Is(err, unix.EWOULDBLOCK) && !errors.Is(err, unix.EAGAIN) {
			return fmt.Errorf("lock Git worktree administration file: %w", err)
		}
		if !timer.Stop() {
			select {
			case <-timer.C:
			default:
			}
		}
		timer.Reset(worktreeAdminLockPollInterval)
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-timer.C:
		}
	}
}

func unlockWorktreeAdminFile(file *os.File) error {
	return unix.Flock(int(file.Fd()), unix.LOCK_UN)
}
