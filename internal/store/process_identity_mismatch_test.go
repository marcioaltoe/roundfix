package store

import (
	"context"
	"errors"
	"fmt"
	"os"
	"runtime"
	"testing"
)

func TestOwnerProcessControllerTellsAProvenIdentityMismatchFromAnUnreadableIdentity(t *testing.T) {
	t.Parallel()
	for _, tt := range []struct {
		name                           string
		pid                            int
		recorded                       string
		readErr                        error
		absentOnRetry                  bool
		mismatch, unproven, unreadable bool
	}{
		{name: "comparable mismatch", pid: 4242, recorded: runtime.GOOS + ":old", mismatch: true, unproven: true},
		{name: "own PID", pid: os.Getpid(), recorded: runtime.GOOS + ":old", unproven: true},
		{name: "PID zero", recorded: runtime.GOOS + ":old", unproven: true},
		{name: "read failure", pid: 4242, recorded: runtime.GOOS + ":old", readErr: errors.New("read denied"), unreadable: true},
		{name: "unsupported read", pid: 4242, recorded: runtime.GOOS + ":old", readErr: ErrOwnerProcessUnsupported, unreadable: true},
		{name: "foreign token", pid: 4242, recorded: "foreign:old", unreadable: true},
		{name: "absent on second read", pid: 4242, recorded: runtime.GOOS + ":old", readErr: errors.New("exited"), absentOnRetry: true},
	} {
		t.Run(tt.name, func(t *testing.T) {
			controller := NewOwnerProcessController()
			reads := 0
			controller.processAbsent = func(int) (bool, error) { reads++; return tt.absentOnRetry && reads > 1, nil }
			controller.processStartIdentity = func(context.Context, int) (string, error) { return runtime.GOOS + ":live", tt.readErr }
			controller.ownedProcesses = func(int, string) ([]int, error) {
				if !tt.absentOnRetry {
					t.Fatal("refused proof enumerated a process tree")
				}
				return nil, nil
			}
			controller.signalProcess = func(int, bool) error { t.Fatal("refused proof sent a signal"); return nil }
			pids, err := controller.InspectTree(t.Context(), tt.pid, tt.recorded)
			if len(pids) != 0 || errors.Is(err, ErrOwnerProcessIdentityMismatch) != tt.mismatch ||
				errors.Is(err, ErrOwnerProcessIdentityUnproven) != tt.unproven || errors.Is(err, ErrOwnerIdentityUnreadable) != tt.unreadable {
				t.Fatalf("InspectTree = %v, %v", pids, err)
			}
			if tt.absentOnRetry && err != nil {
				t.Fatal(err)
			}
			if tt.mismatch {
				want := fmt.Sprintf("owner process %d failed step %q: owner process identity cannot be proven: live process start identity does not match the recorded owner identity", tt.pid, "prove owner process identity")
				if err.Error() != want {
					t.Fatalf("message = %q, want %q", err, want)
				}
				if err := controller.TerminateAndWait(t.Context(), tt.pid, tt.recorded); !errors.Is(err, ErrOwnerProcessIdentityMismatch) {
					t.Fatalf("termination error = %v", err)
				}
			}
		})
	}
}
