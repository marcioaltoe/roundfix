//go:build unix

// Suite: Machine Verification Capacity owner reclamation.
// Boundary: real process identity reads and temporary entry files; no signals.
package store

import (
	"context"
	"errors"
	"os"
	"runtime"
	"testing"

	"roundfix/internal/testwait"
)

func TestMachineVerificationCapacityRemovesOnlyAnOwnerProvenGone(t *testing.T) {
	t.Parallel()
	parent := os.Getppid()
	identity, err := OwnerProcessIdentity(t.Context(), parent)
	if err != nil {
		t.Fatalf("read parent identity: %v", err)
	}
	for _, tc := range []struct {
		name     string
		pid      int
		identity string
		gone     bool
	}{
		{name: "reused", pid: parent, identity: runtime.GOOS + ":0.0", gone: true},
		{name: "absent", pid: reapedOwnerPID(t), gone: true},
		{name: "matching", pid: parent, identity: identity},
		{name: "foreign", pid: parent, identity: "foreign:0.0"},
		{name: "empty", pid: parent},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			home := t.TempDir()
			path := writeMachineTestEntry(t, home, 1, "holding", tc.pid, tc.identity)
			acquisition := startMachineTestAcquisition(t, home, "new", 1)
			if tc.gone {
				lease := acquisition.granted(t)
				if lease.Waited() != 0 {
					t.Error("proven-gone owner delayed acquisition")
				}
				select {
				case wait := <-acquisition.waits:
					t.Fatalf("proven-gone entry caused wait: %+v", wait)
				default:
				}
				if _, err := os.Stat(path); !errors.Is(err, os.ErrNotExist) {
					t.Fatalf("proven-gone entry remains: %v", err)
				}
				releaseMachineTestLease(t, lease)
				return
			}
			wait := testwait.Until(t, "live or unproven owner wait", acquisition.waits, acquisition.done)
			if len(wait.Holding) != 1 || wait.Holding[0].Path != path {
				t.Fatalf("wait = %+v", wait)
			}
			acquisition.cancel()
			testwait.Until[struct{}, struct{}](t, "cancelled acquisition exit", acquisition.done, nil)
			if !errors.Is(acquisition.err, context.Canceled) || acquisition.lease != nil {
				t.Fatalf("live owner displaced: lease=%v err=%v", acquisition.lease, acquisition.err)
			}
			if _, err := os.Stat(path); err != nil {
				t.Fatalf("live or unproven owner removed: %v", err)
			}
		})
	}
}
