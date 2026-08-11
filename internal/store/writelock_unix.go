//go:build !windows

package store

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"golang.org/x/sys/unix"
)

// writeLockPollInterval bounds how long a waiting Roundfix writer sleeps
// between non-blocking attempts at the machine-wide advisory write lock. A
// small interval keeps lock waits responsive to cancellation while bounding
// the polling cost.
const writeLockPollInterval = 10 * time.Millisecond

// acquireWriteLock takes the machine-wide advisory write lock on file,
// waiting until it is free or ctx is cancelled. Roundfix processes serialize
// here before SQLite ever sees a writer, so no Roundfix writer races another
// Roundfix writer inside SQLite.
func acquireWriteLock(file *os.File, ctx context.Context) error {
	for {
		err := unix.Flock(int(file.Fd()), unix.LOCK_EX|unix.LOCK_NB)
		if err == nil {
			return nil
		}
		if !errors.Is(err, unix.EWOULDBLOCK) && !errors.Is(err, unix.EAGAIN) {
			return fmt.Errorf("acquire machine-wide write lock: %w", err)
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(writeLockPollInterval):
		}
	}
}

// releaseWriteLock releases the machine-wide advisory write lock held on
// file. It must run on every exit path after the SQLite transaction that held
// it has committed or rolled back.
func releaseWriteLock(file *os.File) error {
	return unix.Flock(int(file.Fd()), unix.LOCK_UN)
}
