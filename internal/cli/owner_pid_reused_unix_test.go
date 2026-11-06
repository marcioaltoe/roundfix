//go:build unix

// Suite: reused owner PIDs in Active Run reclaim, Run slots and queue start.
// Boundary: real parent-process identity reads and disposable Roundfix homes; no signals.
package cli

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"testing"

	"roundfix/internal/store"
)

func TestImplementReclaimsAnActiveRunWhoseOwnerPIDWasReused(t *testing.T) {
	t.Parallel()
	pid, identity := parentOwnerIdentityForReuseTest(t)
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
			db, repo := openOwnerPIDReuseTestStore(t)
			run, err := db.CreateRun(ctx, store.CreateRunRequest{
				Kind: store.KindImplement, GitRoot: repo, LocalBranch: "fix/reused-owner",
				SpecSlug: "0001-reused-owner", OwnerPID: pid, OwnerIdentity: tt.identity,
			})
			if err != nil {
				t.Fatalf("create Active Run: %v", err)
			}
			var stderr bytes.Buffer
			after, reclaimed, err := reclaimOrphanedActiveRun(ctx, db, run, &stderr)
			if err != nil || reclaimed != tt.gone {
				t.Fatalf("reclaim=%v err=%v, want reclaim=%v", reclaimed, err, tt.gone)
			}
			stored, found, err := db.Run(ctx, run.ID)
			if err != nil || !found || !reflect.DeepEqual(stored, after) {
				t.Fatalf("read Run: stored=%+v returned=%+v found=%v err=%v", stored, after, found, err)
			}
			locked, found, err := db.ActiveSpecRun(ctx, repo, run.SpecSlug)
			if err != nil {
				t.Fatalf("read Active Run lock: %v", err)
			}
			events, err := db.RunEventsAfter(ctx, run.ID, 0, 10)
			if err != nil {
				t.Fatalf("read Run Events: %v", err)
			}
			if tt.gone {
				reason := fmt.Sprintf("owner process %d gone (pid reused); lock reclaimed", pid)
				want := fmt.Sprintf("roundfix: reclaimed orphaned Active Run %s: %s\n", run.ID, reason)
				if after.State != store.StateFailed || after.CompletedAt == nil || found || stderr.String() != want {
					t.Fatalf("reclaim Run=%+v lock present=%v stderr=%q, want %q", after, found, stderr.String(), want)
				}
				if len(events) != 1 || !bytes.Contains(events[0].Event.Payload, []byte(reason)) {
					t.Fatalf("reclamation event missing reuse reason: %+v", events)
				}
				return
			}
			if !reflect.DeepEqual(after, run) || !found || locked.ID != run.ID || stderr.Len() != 0 || len(events) != 0 {
				t.Fatalf("present owner changed: Run=%+v lock=%+v found=%v stderr=%q events=%+v", after, locked, found, stderr.String(), events)
			}
		})
	}
}

func TestRunSlotsExcludeAHolderWhoseOwnerPIDWasReused(t *testing.T) {
	t.Parallel()
	pid, identity := parentOwnerIdentityForReuseTest(t)
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
			db, repo := openOwnerPIDReuseTestStore(t)
			run, err := db.CreateRun(ctx, store.CreateRunRequest{
				Kind: store.KindImplement, GitRoot: repo, LocalBranch: "fix/reused-owner",
				SpecSlug: "0001-reused-owner", OwnerPID: pid, OwnerIdentity: tt.identity,
			})
			if err != nil {
				t.Fatalf("create Active Run: %v", err)
			}
			workflow := &commandDeliveryWorkflow{store: db}
			workflow.loaded.Config.Runs.MaxActive = 2
			slots, err := workflow.RunSlots(ctx)
			if err != nil || slots.MaxActive != 2 {
				t.Fatalf("Run slots=%+v err=%v", slots, err)
			}
			if tt.gone {
				if len(slots.Holders) != 0 {
					t.Fatalf("reused PID retained a Run slot: %+v", slots.Holders)
				}
			} else if len(slots.Holders) != 1 || slots.Holders[0].RunID != run.ID || slots.Holders[0].Repository != repo || slots.Holders[0].SpecSlug != run.SpecSlug {
				t.Fatalf("present owner lost its Run slot: %+v", slots.Holders)
			}
			after, found, err := db.Run(ctx, run.ID)
			if err != nil || !found || !reflect.DeepEqual(after, run) {
				t.Fatalf("Run slot read mutated Run: after=%+v found=%v err=%v", after, found, err)
			}
		})
	}
}

func TestDeliverStartReleasesAQueueOwnerWhosePIDWasReused(t *testing.T) {
	t.Parallel()
	pid, identity := parentOwnerIdentityForReuseTest(t)
	for _, tt := range []struct {
		name     string
		identity string
		gone     bool
	}{
		{name: "reused", identity: runtime.GOOS + ":0.0", gone: true},
		{name: "matching", identity: identity},
		{name: "foreign identity", identity: "foreign:0.0"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctx := t.Context()
			db, repo := openOwnerPIDReuseTestStore(t)
			if _, err := db.CreateDeliveryQueue(ctx, repo, []string{"0001-reused-owner"}); err != nil {
				t.Fatalf("create Delivery Queue: %v", err)
			}
			if err := db.ClaimDeliveryQueueOwner(ctx, repo, pid, tt.identity); err != nil {
				t.Fatalf("claim Delivery Queue owner: %v", err)
			}
			before, found, err := db.DeliveryQueue(ctx, repo)
			if err != nil || !found {
				t.Fatalf("read Delivery Queue before release: found=%v err=%v", found, err)
			}
			if err := releaseGoneDeliveryQueueOwner(ctx, db, repo); err != nil {
				t.Fatalf("release gone queue owner: %v", err)
			}
			after, found, err := db.DeliveryQueue(ctx, repo)
			if err != nil || !found {
				t.Fatalf("read Delivery Queue after release: found=%v err=%v", found, err)
			}
			want := before
			if tt.gone {
				want.OwnerPID = 0
				want.OwnerIdentity = ""
			}
			if !reflect.DeepEqual(after, want) {
				t.Fatalf("Delivery Queue after release=%+v, want %+v", after, want)
			}
			if tt.gone {
				if err := db.ClaimDeliveryQueueOwner(ctx, repo, pid, identity); err != nil {
					t.Fatalf("take over released queue owner: %v", err)
				}
			}
		})
	}
}

func parentOwnerIdentityForReuseTest(t *testing.T) (int, string) {
	t.Helper()
	pid := os.Getppid()
	identity, err := store.OwnerProcessIdentity(t.Context(), pid)
	if errors.Is(err, store.ErrOwnerProcessUnsupported) {
		t.Skipf("parent process start identity is unsupported: %v", err)
	}
	if err != nil {
		t.Fatalf("read parent process identity: %v", err)
	}
	return pid, identity
}

func openOwnerPIDReuseTestStore(t *testing.T) (*store.Store, string) {
	t.Helper()
	home := t.TempDir()
	db, err := store.Open(t.Context(), home)
	if err != nil {
		t.Fatalf("open Run Database: %v", err)
	}
	t.Cleanup(func() {
		if err := db.Close(); err != nil {
			t.Errorf("close Run Database: %v", err)
		}
	})
	return db, filepath.Join(home, "repository")
}
