package speccheck

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"roundfix/internal/worktree"
)

// ErrDisposableCheckoutCreate is the named reason returned when Git cannot
// materialize the repository at HEAD.
var ErrDisposableCheckoutCreate = errors.New("create disposable checkout")

// DisposableCheckout materializes repoRoot at HEAD in a temporary detached
// worktree. Callers must defer cleanup immediately after a successful return;
// cleanup uses a fresh context so caller cancellation and panic unwinding do
// not prevent removal.
func DisposableCheckout(ctx context.Context, repoRoot string) (dir string, cleanup func() error, err error) {
	if ctx == nil {
		ctx = context.Background()
	}
	repoRoot = strings.TrimSpace(repoRoot)
	if repoRoot == "" {
		return "", nil, fmt.Errorf("%w: repository root is required", ErrDisposableCheckoutCreate)
	}

	tempRoot, err := os.MkdirTemp("", "roundfix-verification-checkout-")
	if err != nil {
		return "", nil, fmt.Errorf("%w: allocate temporary directory: %w", ErrDisposableCheckoutCreate, err)
	}
	dir = filepath.Join(tempRoot, "worktree")
	cleanup = func() error {
		return removeDisposableCheckout(repoRoot, tempRoot, dir)
	}

	if err := worktree.AddDetachedWorktree(ctx, repoRoot, dir, "HEAD"); err != nil {
		cleanupErr := cleanup()
		return "", nil, fmt.Errorf("%w: %w", ErrDisposableCheckoutCreate, errors.Join(err, cleanupErr))
	}
	return dir, cleanup, nil
}

func removeDisposableCheckout(repoRoot string, tempRoot string, dir string) error {
	var cleanupErr error
	if _, err := os.Lstat(dir); err == nil {
		if err := worktree.RemoveWorktreeForce(context.Background(), repoRoot, dir); err != nil {
			cleanupErr = fmt.Errorf("remove disposable checkout: %w", err)
		}
	} else if !errors.Is(err, os.ErrNotExist) {
		cleanupErr = fmt.Errorf("inspect disposable checkout before removal: %w", err)
	}
	if err := os.RemoveAll(tempRoot); err != nil {
		cleanupErr = errors.Join(cleanupErr, fmt.Errorf("remove disposable checkout directory: %w", err))
	}
	return cleanupErr
}
