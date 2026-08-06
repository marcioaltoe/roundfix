//go:build windows

package store

import (
	"fmt"
	"math"

	"golang.org/x/sys/windows"
)

func availableTemporaryCapacity(path string) (int64, error) {
	pathPointer, err := windows.UTF16PtrFromString(path)
	if err != nil {
		return 0, fmt.Errorf("encode filesystem path %q: %w", path, err)
	}
	var availableBytes uint64
	if err := windows.GetDiskFreeSpaceEx(pathPointer, &availableBytes, nil, nil); err != nil {
		return 0, fmt.Errorf("measure filesystem capacity for %q: %w", path, err)
	}
	if availableBytes > math.MaxInt64 {
		return math.MaxInt64, nil
	}
	return int64(availableBytes), nil
}
