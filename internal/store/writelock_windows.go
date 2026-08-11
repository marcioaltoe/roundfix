//go:build windows

package store

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"golang.org/x/sys/windows"
)

// writeLockPollInterval bounds how long a waiting Roundfix writer sleeps
// between non-blocking attempts at the machine-wide advisory write lock.
const writeLockPollInterval = 10 * time.Millisecond

// acquireWriteLock takes the machine-wide advisory write lock on file,
// waiting until it is free or ctx is cancelled. Roundfix processes serialize
// here before SQLite ever sees a writer.
func acquireWriteLock(ctx context.Context, file *os.File) error {
	timer := time.NewTimer(writeLockPollInterval)
	defer timer.Stop()
	for {
		overlapped := &windows.Overlapped{}
		err := windows.LockFileEx(
			windows.Handle(file.Fd()),
			windows.LOCKFILE_EXCLUSIVE_LOCK|windows.LOCKFILE_FAIL_IMMEDIATELY,
			0,
			1,
			0,
			overlapped,
		)
		if err == nil {
			return nil
		}
		if !errors.Is(err, windows.ERROR_LOCK_VIOLATION) {
			return fmt.Errorf("acquire machine-wide write lock: %w", err)
		}
		if !timer.Stop() {
			select {
			case <-timer.C:
			default:
			}
		}
		timer.Reset(writeLockPollInterval)
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-timer.C:
		}
	}
}

// releaseWriteLock releases the machine-wide advisory write lock held on
// file. It must run on every exit path after the SQLite transaction that held
// it has committed or rolled back.
func releaseWriteLock(file *os.File) error {
	return windows.UnlockFileEx(windows.Handle(file.Fd()), 0, 1, 0, &windows.Overlapped{})
}
