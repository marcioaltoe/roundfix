package speccheck

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
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

	if err := runDisposableCheckoutGit(ctx, repoRoot, "worktree", "add", "--detach", dir, "HEAD"); err != nil {
		cleanupErr := cleanup()
		return "", nil, fmt.Errorf("%w: %w", ErrDisposableCheckoutCreate, errors.Join(err, cleanupErr))
	}
	return dir, cleanup, nil
}

func removeDisposableCheckout(repoRoot string, tempRoot string, dir string) error {
	var cleanupErr error
	if _, err := os.Lstat(dir); err == nil {
		if err := runDisposableCheckoutGit(context.Background(), repoRoot, "worktree", "remove", "--force", dir); err != nil {
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

func runDisposableCheckoutGit(ctx context.Context, repoRoot string, args ...string) error {
	commandArgs := []string{
		"-C", repoRoot,
		"-c", "user.name=Roundfix Test",
		"-c", "user.email=test@example.com",
		"-c", "commit.gpgsign=false",
		"-c", "core.fsmonitor=false",
	}
	commandArgs = append(commandArgs, args...)
	command := exec.CommandContext(ctx, "git", commandArgs...)
	command.Env = disposableCheckoutGitEnv()
	output, err := command.CombinedOutput()
	if err == nil {
		return nil
	}
	diagnostic := strings.TrimSpace(string(output))
	if diagnostic == "" {
		return fmt.Errorf("git %s: %w", strings.Join(args, " "), err)
	}
	return fmt.Errorf("git %s: %w: %s", strings.Join(args, " "), err, diagnostic)
}

func disposableCheckoutGitEnv() []string {
	env := make([]string, 0, len(os.Environ())+5)
	for _, entry := range os.Environ() {
		key, _, _ := strings.Cut(entry, "=")
		if strings.HasPrefix(key, "GIT_CONFIG_") || key == "GIT_TEST_MAINT_AUTO_DETACH" {
			continue
		}
		env = append(env, entry)
	}
	return append(env,
		"GIT_CONFIG_GLOBAL=/dev/null",
		"GIT_CONFIG_SYSTEM=/dev/null",
		"GIT_TEST_MAINT_AUTO_DETACH=0",
		"GIT_OPTIONAL_LOCKS=0",
		"GIT_TERMINAL_PROMPT=0",
	)
}
