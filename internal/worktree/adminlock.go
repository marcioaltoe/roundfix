package worktree

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

const (
	worktreeAdminLockFilename     = "roundfix-worktree-admin.lock"
	worktreeAdminLockPollInterval = 10 * time.Millisecond
)

// worktreeAdminPermits serializes worktree administration for each Git common
// directory inside one Roundfix process. The file lock taken after this permit
// extends the same boundary across Roundfix processes.
var worktreeAdminPermits sync.Map

func runWorktreeCommand(
	ctx context.Context,
	runner gitRunner,
	workDir string,
	args ...string,
) (output string, resultErr error) {
	if len(args) < 2 || args[0] != "worktree" {
		return "", fmt.Errorf("run Git worktree command: invalid arguments %v", args)
	}
	if ctx == nil {
		ctx = context.Background()
	}

	commonDir, err := resolveGitCommonDir(ctx, runner, workDir)
	if err != nil {
		return "", err
	}
	release, err := acquireWorktreeAdminLock(ctx, commonDir)
	if err != nil {
		return "", err
	}
	defer func() {
		if releaseErr := release(); releaseErr != nil {
			releaseErr = fmt.Errorf("release Git worktree administration lock: %w", releaseErr)
			if resultErr == nil {
				resultErr = releaseErr
			} else {
				resultErr = errors.Join(resultErr, releaseErr)
			}
		}
	}()

	return runner.Run(ctx, workDir, args...)
}

func resolveGitCommonDir(ctx context.Context, runner gitRunner, workDir string) (string, error) {
	output, err := runner.Run(ctx, workDir, "rev-parse", "--path-format=absolute", "--git-common-dir")
	if err != nil {
		return "", fmt.Errorf("resolve Git common directory from %q: %w", workDir, err)
	}
	commonDir := strings.TrimSpace(output)
	if commonDir == "" || strings.ContainsAny(commonDir, "\x00\r\n") {
		return "", fmt.Errorf("resolve Git common directory from %q: Git returned invalid path %q", workDir, commonDir)
	}
	if resolved, err := filepath.EvalSymlinks(commonDir); err == nil {
		return filepath.Clean(resolved), nil
	}
	absolute, err := filepath.Abs(commonDir)
	if err != nil {
		return "", fmt.Errorf("resolve Git common directory %q: %w", commonDir, err)
	}
	return filepath.Clean(absolute), nil
}

func acquireWorktreeAdminLock(ctx context.Context, commonDir string) (func() error, error) {
	permit := worktreeAdminPermit(commonDir)
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case <-permit:
	}
	releasePermit := func() {
		permit <- struct{}{}
	}

	lockPath := filepath.Join(commonDir, worktreeAdminLockFilename)
	file, err := os.OpenFile(lockPath, os.O_CREATE|os.O_RDWR, 0o600)
	if err != nil {
		releasePermit()
		return nil, fmt.Errorf("open Git worktree administration lock %q: %w", lockPath, err)
	}
	if err := lockWorktreeAdminFile(ctx, file); err != nil {
		closeErr := file.Close()
		releasePermit()
		if ctxErr := ctx.Err(); ctxErr != nil && errors.Is(err, ctxErr) {
			return nil, ctxErr
		}
		return nil, errors.Join(
			fmt.Errorf("acquire Git worktree administration lock %q: %w", lockPath, err),
			closeErr,
		)
	}

	return func() error {
		unlockErr := unlockWorktreeAdminFile(file)
		closeErr := file.Close()
		releasePermit()
		return errors.Join(unlockErr, closeErr)
	}, nil
}

func worktreeAdminPermit(commonDir string) chan struct{} {
	created := make(chan struct{}, 1)
	created <- struct{}{}
	permit, _ := worktreeAdminPermits.LoadOrStore(commonDir, created)
	return permit.(chan struct{})
}
