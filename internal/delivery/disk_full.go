package delivery

import (
	"errors"
	"fmt"
)

const BlockerDiskFull = "disk-full"

// DiskFullError retains the step's error while naming the full filesystem.
type DiskFullError struct {
	Path           string
	AvailableBytes int64
	Err            error
}

func (err *DiskFullError) Error() string { return err.Err.Error() }
func (err *DiskFullError) Unwrap() error { return err.Err }
func (err *DiskFullError) Blocker() string {
	return fmt.Sprintf("%s: %s: %d MiB free", BlockerDiskFull, err.Path, err.AvailableBytes/(1024*1024))
}

func stepErrorBlocker(err error) string {
	var full *DiskFullError
	if errors.As(err, &full) {
		return full.Blocker()
	}
	return BlockerDeliveryError + ": " + err.Error()
}
