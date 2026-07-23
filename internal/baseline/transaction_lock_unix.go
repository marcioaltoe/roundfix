//go:build !windows

package baseline

import (
	"errors"
	"fmt"
	"os"

	"golang.org/x/sys/unix"
)

func lockTransactionFile(file *os.File) error {
	err := unix.Flock(int(file.Fd()), unix.LOCK_EX|unix.LOCK_NB)
	if errors.Is(err, unix.EWOULDBLOCK) || errors.Is(err, unix.EAGAIN) {
		return ErrTransactionLocked
	}
	if err != nil {
		return fmt.Errorf("lock Baseline transaction: %w", err)
	}
	return nil
}

func unlockTransactionFile(file *os.File) error {
	return unix.Flock(int(file.Fd()), unix.LOCK_UN)
}
