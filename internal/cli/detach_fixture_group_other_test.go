//go:build !darwin && !linux

package cli

import (
	"errors"
	"fmt"
	"syscall"
)

// Other Unix systems retain the signal probe. Only ESRCH proves absence;
// permission and other probe failures must not look like an empty group.
func detachFixtureGroupLiveMembers(pgid int) ([]int, error) {
	err := syscall.Kill(-pgid, 0)
	if errors.Is(err, syscall.ESRCH) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("probe fixture process group %d: %w", pgid, err)
	}
	return []int{pgid}, nil
}
