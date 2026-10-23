//go:build unix

package store

import (
	"errors"
	"runtime"
	"testing"
)

func TestInspectTreeReportsAProvenMismatchOfALiveProcessWithoutSignallingIt(t *testing.T) {
	t.Parallel()
	pid, _ := startOwnerProcessHelper(t, "graceful")
	controller := NewOwnerProcessController()
	pids, err := controller.InspectTree(t.Context(), pid, runtime.GOOS+":ended-owner")
	if !errors.Is(err, ErrOwnerProcessIdentityMismatch) || len(pids) != 0 {
		t.Fatalf("InspectTree = %v, %v", pids, err)
	}
	if !ProcessAlive(pid) {
		t.Fatal("mismatch inspection stopped the live helper")
	}
}
