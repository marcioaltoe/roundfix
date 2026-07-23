//go:build windows

package baseline

import (
	"errors"
	"fmt"
	"os"

	"golang.org/x/sys/windows"
)

func lockTransactionFile(file *os.File) error {
	overlapped := &windows.Overlapped{}
	err := windows.LockFileEx(
		windows.Handle(file.Fd()),
		windows.LOCKFILE_EXCLUSIVE_LOCK|windows.LOCKFILE_FAIL_IMMEDIATELY,
		0,
		1,
		0,
		overlapped,
	)
	if errors.Is(err, windows.ERROR_LOCK_VIOLATION) {
		return ErrTransactionLocked
	}
	if err != nil {
		return fmt.Errorf("lock Baseline transaction: %w", err)
	}
	return nil
}

func unlockTransactionFile(file *os.File) error {
	return windows.UnlockFileEx(windows.Handle(file.Fd()), 0, 1, 0, &windows.Overlapped{})
}
