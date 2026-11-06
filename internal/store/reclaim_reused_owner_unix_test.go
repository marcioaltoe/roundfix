//go:build unix

// Suite: Active Run reclaim after owner PID reuse.
// Boundary: real parent-process identity reads and a disposable Run Database; no signals.
package store

import (
	"errors"
	"fmt"
	"os"
	"reflect"
	"runtime"
	"testing"
)

func TestReclaimOrphanedRunReclaimsAnOwnerWhosePIDWasReused(t *testing.T) {
	t.Parallel()
	pid := os.Getppid()
	identity, err := OwnerProcessIdentity(t.Context(), pid)
	if errors.Is(err, ErrOwnerProcessUnsupported) {
		t.Skipf("parent process start identity is unsupported: %v", err)
	}
	if err != nil {
		t.Fatalf("read parent process identity: %v", err)
	}
	for _, tt := range []struct {
		name     string
		identity string
		gone     bool
	}{
		{name: "reused", identity: runtime.GOOS + ":0.0", gone: true},
		{name: "matching", identity: identity},
		{name: "foreign identity", identity: "foreign:0.0"},
		{name: "empty identity"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctx := t.Context()
			runStore := openTestStore(t, ctx, t.TempDir())
			defer closeStore(t, runStore)
			req := sampleCreateRunRequest()
			req.OwnerPID = pid
			req.OwnerIdentity = tt.identity
			run, err := runStore.CreateRun(ctx, req)
			if err != nil {
				t.Fatalf("create Active Run: %v", err)
			}
			if err := runStore.ReclaimOrphanedRun(ctx, run, ""); err != nil {
				t.Fatalf("reclaim Active Run: %v", err)
			}
			after, found, err := runStore.Run(ctx, run.ID)
			if err != nil || !found {
				t.Fatalf("read Run: found=%v err=%v", found, err)
			}
			locked, found, err := runStore.ActiveRun(ctx, req.HeadRepository, req.HeadBranch)
			if err != nil {
				t.Fatalf("read Active Run lock: %v", err)
			}
			if tt.gone {
				if after.State != StateFailed || after.CompletedAt == nil || found {
					t.Fatalf("reclaimed Run=%+v, lock present=%v", after, found)
				}
				assertReclamationEvent(t, ctx, runStore, run.ID,
					fmt.Sprintf("owner process %d gone (pid reused); lock reclaimed", pid), pid)
				return
			}
			if !reflect.DeepEqual(after, run) || !found || locked.ID != run.ID {
				t.Fatalf("present owner changed: Run=%+v lock=%+v found=%v", after, locked, found)
			}
			events, err := runStore.RunEventsAfter(ctx, run.ID, 0, 10)
			if err != nil || len(events) != 0 {
				t.Fatalf("present owner journal changed: events=%+v err=%v", events, err)
			}
		})
	}
}
