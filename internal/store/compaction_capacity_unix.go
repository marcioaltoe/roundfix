//go:build !windows

package store

import (
	"fmt"
	"math"

	"golang.org/x/sys/unix"
)

func availableTemporaryCapacity(path string) (int64, error) {
	var filesystem unix.Statfs_t
	if err := unix.Statfs(path, &filesystem); err != nil {
		return 0, fmt.Errorf("measure filesystem capacity for %q: %w", path, err)
	}
	availableBlocks := uint64(filesystem.Bavail)
	blockSize := uint64(filesystem.Bsize)
	if blockSize != 0 && availableBlocks > uint64(math.MaxInt64)/blockSize {
		return math.MaxInt64, nil
	}
	return int64(availableBlocks * blockSize), nil
}
