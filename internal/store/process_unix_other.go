//go:build unix && !linux && !darwin

package store

import (
	"context"
	"fmt"
)

func processStartIdentity(_ context.Context, _ int) (string, error) {
	return "", fmt.Errorf("owner process identity is unreadable: %w", ErrOwnerProcessUnsupported)
}

func processTreePIDs(_ int, _ string) ([]int, error) {
	return nil, ErrOwnerProcessUnsupported
}

func liveProcessGroupMembers(_ context.Context, _ int) ([]int, error) {
	return nil, ErrOwnerProcessUnsupported
}
