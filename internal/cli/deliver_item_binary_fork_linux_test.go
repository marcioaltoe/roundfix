//go:build linux

// Boundary: executable copies in a temporary directory while this process forks.
package cli

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
	"syscall"
	"testing"

	"roundfix/internal/testfixture"
)

func TestDeliveryItemBinaryCopyStartsWhileTheOwnerForks(t *testing.T) {
	t.Parallel()
	fixture := testfixture.FixtureBinary(t, "fork-lock-item", "package main\n\nfunc main() {}\n")
	program, err := os.ReadFile(fixture)
	if err != nil {
		t.Fatal(err)
	}
	dir := t.TempDir()
	ctx, cancel := context.WithCancel(t.Context())
	var forks sync.WaitGroup
	for range 4 {
		forks.Add(1)
		go func() {
			defer forks.Done()
			for ctx.Err() == nil {
				_ = exec.Command(fixture).Run()
			}
		}()
	}
	defer func() { cancel(); forks.Wait() }()
	for i := range 100 {
		source := filepath.Join(dir, fmt.Sprintf("build-%d", i), "roundfix")
		if err := os.MkdirAll(filepath.Dir(source), 0o755); err != nil {
			t.Fatal(err)
		}
		// A distinct trailer gives each build its own copy directory.
		if err := os.WriteFile(source, append(append([]byte(nil), program...), fmt.Sprintf("build %d", i)...), 0o600); err != nil {
			t.Fatal(err)
		}
		copied, err := copyDeliveryItemBinary(source, filepath.Join(dir, "item-binary"))
		if err != nil {
			t.Fatal(err)
		}
		if err := exec.Command(copied).Run(); err != nil {
			t.Fatalf("start %d of a fresh copy: %v (text file busy: %t)", i, err, errors.Is(err, syscall.ETXTBSY))
		}
		if err := os.RemoveAll(filepath.Dir(copied)); err != nil {
			t.Fatal(err)
		}
	}
}
