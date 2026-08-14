// Suite: disposable Verification checkout lifecycle
// Invariant: authored commands run at HEAD without changing the operator checkout, and their temporary tree is always removable.
// Boundary IN: the public disposable-checkout helper and a real local Git repository
// Boundary OUT: Verification command execution and CLI reporting
package speccheck_test

import (
	"bytes"
	"context"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"roundfix/internal/gittest"
	"roundfix/internal/speccheck"
)

func TestDisposableCheckoutIsolatesOperatorTreeAndIndex(t *testing.T) {
	repoRoot := newDisposableCheckoutRepository(t)

	if err := os.WriteFile(filepath.Join(repoRoot, "tracked.txt"), []byte("operator staged\n"), 0o644); err != nil {
		t.Fatalf("write operator's staged file: %v", err)
	}
	gittest.Run(t, repoRoot, "add", "tracked.txt")
	if err := os.WriteFile(filepath.Join(repoRoot, "operator-untracked.txt"), []byte("operator only\n"), 0o644); err != nil {
		t.Fatalf("write operator's untracked file: %v", err)
	}
	statusBefore := gittest.Run(t, repoRoot, "status", "--porcelain=v1", "--untracked-files=all")
	indexBefore, err := os.ReadFile(filepath.Join(repoRoot, ".git", "index"))
	if err != nil {
		t.Fatalf("read operator index before checkout: %v", err)
	}

	dir, cleanup, err := speccheck.DisposableCheckout(context.Background(), repoRoot)
	if err != nil {
		t.Fatalf("DisposableCheckout() error: %v", err)
	}
	if filepath.Clean(dir) == filepath.Clean(repoRoot) {
		t.Fatalf("DisposableCheckout() path = repository root %q", repoRoot)
	}

	wantHead := strings.TrimSpace(gittest.Run(t, repoRoot, "rev-parse", "HEAD"))
	gotHead := strings.TrimSpace(gittest.Run(t, dir, "rev-parse", "HEAD"))
	if gotHead != wantHead {
		t.Fatalf("disposable checkout HEAD = %q, want %q", gotHead, wantHead)
	}
	tracked, err := os.ReadFile(filepath.Join(dir, "tracked.txt"))
	if err != nil {
		t.Fatalf("read tracked file from disposable checkout: %v", err)
	}
	if got, want := string(tracked), "committed\n"; got != want {
		t.Fatalf("disposable tracked file = %q, want HEAD content %q", got, want)
	}
	if _, err := os.Stat(filepath.Join(dir, "operator-untracked.txt")); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("operator untracked file visible in disposable checkout: %v", err)
	}

	if err := os.WriteFile(filepath.Join(dir, "tracked.txt"), []byte("disposable staged\n"), 0o644); err != nil {
		t.Fatalf("write tracked file in disposable checkout: %v", err)
	}
	if err := os.WriteFile(filepath.Join(dir, "disposable-only.txt"), []byte("disposable only\n"), 0o644); err != nil {
		t.Fatalf("write untracked file in disposable checkout: %v", err)
	}
	gittest.Run(t, dir, "add", "tracked.txt", "disposable-only.txt")

	if got := gittest.Run(t, repoRoot, "status", "--porcelain=v1", "--untracked-files=all"); got != statusBefore {
		t.Fatalf("operator status changed\nbefore:\n%s\nafter:\n%s", statusBefore, got)
	}
	indexAfter, err := os.ReadFile(filepath.Join(repoRoot, ".git", "index"))
	if err != nil {
		t.Fatalf("read operator index after disposable writes: %v", err)
	}
	if !bytes.Equal(indexAfter, indexBefore) {
		t.Fatal("operator index bytes changed after staging inside disposable checkout")
	}

	if err := cleanup(); err != nil {
		t.Fatalf("cleanup disposable checkout: %v", err)
	}
	if _, err := os.Stat(dir); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("disposable checkout remains after cleanup: %v", err)
	}
}

func TestDisposableCheckoutCleanupRunsWhenCallerFailsOrPanics(t *testing.T) {
	repoRoot := newDisposableCheckoutRepository(t)

	t.Run("caller returns error", func(t *testing.T) {
		var checkoutDir string
		callerErr := func() (returnErr error) {
			var cleanup func() error
			var err error
			checkoutDir, cleanup, err = speccheck.DisposableCheckout(context.Background(), repoRoot)
			if err != nil {
				return err
			}
			defer func() {
				returnErr = errors.Join(returnErr, cleanup())
			}()
			return errors.New("caller failed")
		}()
		if callerErr == nil || !strings.Contains(callerErr.Error(), "caller failed") {
			t.Fatalf("caller error = %v, want caller failure", callerErr)
		}
		if _, err := os.Stat(checkoutDir); !errors.Is(err, os.ErrNotExist) {
			t.Fatalf("disposable checkout remains after caller error: %v", err)
		}
	})

	t.Run("caller panics", func(t *testing.T) {
		var checkoutDir string
		var cleanupErr error
		var recovered any
		func() {
			defer func() {
				recovered = recover()
			}()
			var cleanup func() error
			var err error
			checkoutDir, cleanup, err = speccheck.DisposableCheckout(context.Background(), repoRoot)
			if err != nil {
				t.Fatalf("DisposableCheckout() error: %v", err)
			}
			defer func() {
				cleanupErr = cleanup()
			}()
			panic("caller panicked")
		}()
		if recovered != "caller panicked" {
			t.Fatalf("recovered panic = %v, want caller panic", recovered)
		}
		if cleanupErr != nil {
			t.Fatalf("cleanup after panic: %v", cleanupErr)
		}
		if _, err := os.Stat(checkoutDir); !errors.Is(err, os.ErrNotExist) {
			t.Fatalf("disposable checkout remains after caller panic: %v", err)
		}
	})
}

func TestDisposableCheckoutCreationRefusalHasNamedReason(t *testing.T) {
	_, cleanup, err := speccheck.DisposableCheckout(context.Background(), t.TempDir())
	if err == nil {
		if cleanup != nil {
			if cleanupErr := cleanup(); cleanupErr != nil {
				t.Fatalf("cleanup unexpected checkout: %v", cleanupErr)
			}
		}
		t.Fatal("DisposableCheckout() error = nil, want creation refusal")
	}
	if !errors.Is(err, speccheck.ErrDisposableCheckoutCreate) {
		t.Fatalf("DisposableCheckout() error = %v, want %v", err, speccheck.ErrDisposableCheckoutCreate)
	}
	if !strings.HasPrefix(err.Error(), speccheck.ErrDisposableCheckoutCreate.Error()+":") {
		t.Fatalf("DisposableCheckout() error = %q, want named creation reason first", err)
	}
}

func newDisposableCheckoutRepository(t *testing.T) string {
	t.Helper()
	repoRoot := filepath.Join(t.TempDir(), "repo")
	gittest.InitRepo(t, repoRoot, "-b", "main")
	if err := os.WriteFile(filepath.Join(repoRoot, "tracked.txt"), []byte("committed\n"), 0o644); err != nil {
		t.Fatalf("write committed file: %v", err)
	}
	gittest.Run(t, repoRoot, "add", "tracked.txt")
	gittest.Run(t, repoRoot, "commit", "-q", "-m", "seed")
	return repoRoot
}
