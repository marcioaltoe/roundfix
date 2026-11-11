package cli

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	roundconfig "roundfix/internal/config"
	"roundfix/internal/store"
)

func TestImplementGivesTheEngineTheUserConfigMachineCapacity(t *testing.T) {
	t.Parallel()
	for _, capacity := range []int{1, 0} {
		t.Run(fmt.Sprint(capacity), func(t *testing.T) {
			t.Parallel()
			home, checkout := t.TempDir(), t.TempDir()
			configDir := filepath.Join(home, ".roundfix")
			if err := os.MkdirAll(configDir, 0o700); err != nil {
				t.Fatal(err)
			}
			if err := os.WriteFile(filepath.Join(configDir, "config.yml"), []byte(fmt.Sprintf("runs:\n  max_verifying: %d\n", capacity)), 0o600); err != nil {
				t.Fatal(err)
			}
			loaded, err := roundconfig.Load(roundconfig.LoadOptions{HomeDir: home, WorkDir: checkout})
			if err != nil {
				t.Fatal(err)
			}
			var stderr bytes.Buffer
			gate := newImplementMachineVerificationGate(loaded, checkout, &stderr)
			ctx, cancel := context.WithCancel(t.Context())
			release, err := gate.Acquire(ctx, store.MachineVerificationHolder{
				Kind: store.MachineVerificationHolderRun, RunID: "run-implement", Repository: "/transient/run-worktree", SpecSlug: "spec-implement", Step: "Task task_01 Verification attempt 1",
			}, func(store.MachineVerificationWait) { t.Error("unexpected wait for empty home") })
			if err != nil {
				cancel()
				t.Fatal(err)
			}
			defer release()
			view, err := store.ReadMachineVerification(t.Context(), home)
			if err != nil {
				cancel()
				t.Fatal(err)
			}
			if len(view.Unreadable) != 0 {
				t.Fatalf("unreadable entries=%v", view.Unreadable)
			}
			if len(view.Entries) != capacity {
				t.Fatalf("entries=%+v, want %d", view.Entries, capacity)
			}
			if capacity == 1 {
				entry := view.Entries[0]
				want := store.MachineVerificationHolder{Kind: store.MachineVerificationHolderRun, RunID: "run-implement", Repository: checkout, SpecSlug: "spec-implement", Step: "Task task_01 Verification attempt 1"}
				if entry.State != "holding" || entry.Holder != want {
					t.Fatalf("entry=%+v, want holding %+v", entry, want)
				}
			} else if _, err := os.Stat(filepath.Join(configDir, "verification-capacity")); !errors.Is(err, os.ErrNotExist) {
				t.Fatalf("off capacity created a directory: %v", err)
			}
			// A Run cancellation must not prevent the release from removing its entry.
			cancel()
			release()
			view, err = store.ReadMachineVerification(context.WithoutCancel(ctx), home)
			if err != nil || len(view.Entries) != 0 {
				t.Fatalf("after release: %+v, %v", view, err)
			}
			if stderr.Len() != 0 {
				t.Fatalf("stderr=%q", stderr.String())
			}
		})
	}
}

func TestMachineVerificationReleaseWarnsWithoutReturningATaskError(t *testing.T) {
	t.Parallel()
	home := t.TempDir()
	var stderr bytes.Buffer
	loaded, err := roundconfig.Load(roundconfig.LoadOptions{HomeDir: home, WorkDir: t.TempDir()})
	if err != nil {
		t.Fatal(err)
	}
	gate := newImplementMachineVerificationGate(loaded, loaded.GitRoot, &stderr)
	release, err := gate.Acquire(t.Context(), store.MachineVerificationHolder{Kind: store.MachineVerificationHolderRun, RunID: "run-release", SpecSlug: "spec-release", Step: "repository check"}, nil)
	if err != nil {
		t.Fatal(err)
	}
	view, err := store.ReadMachineVerification(t.Context(), home)
	if err != nil || len(view.Entries) != 1 {
		t.Fatalf("holding view=%+v, %v", view, err)
	}
	path := view.Entries[0].Path
	// A nonempty directory at the lease path makes removal fail on every host,
	// including hosts whose test user can bypass ordinary mode restrictions.
	if err := os.Remove(path); err != nil {
		t.Fatal(err)
	}
	if err := os.Mkdir(path, 0o700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(path, "keep"), []byte("occupied"), 0o600); err != nil {
		t.Fatal(err)
	}
	release()
	if !strings.HasPrefix(stderr.String(), "roundfix: warning: ") || !strings.Contains(stderr.String(), path) {
		t.Fatalf("warning=%q", stderr.String())
	}
	if err := os.Remove(filepath.Join(path, "keep")); err != nil {
		t.Fatal(err)
	}
	if err := os.Remove(path); err != nil {
		t.Fatal(err)
	}
	release()
}
