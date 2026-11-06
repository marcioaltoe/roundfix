package store

import (
	"context"
	"errors"
	"os"
	"runtime"
	"testing"
)

func TestOwnerGoneTreatsOnlyAProvenMismatchOrAbsenceAsGone(t *testing.T) {
	t.Parallel()
	for _, tt := range []struct {
		name          string
		pid           int
		recorded      string
		absent        bool
		absentErr     error
		readErr       error
		absentOnRetry bool
		gone, reused  bool
	}{
		{name: "absent", pid: 4242, recorded: runtime.GOOS + ":old", absent: true, gone: true},
		{name: "comparable mismatch", pid: 4242, recorded: runtime.GOOS + ":old", gone: true, reused: true},
		{name: "matching", pid: 4242, recorded: runtime.GOOS + ":live"},
		{name: "unreadable identity", pid: 4242, recorded: runtime.GOOS + ":old", readErr: errors.New("read denied")},
		{name: "unsupported identity", pid: 4242, recorded: runtime.GOOS + ":old", readErr: ErrOwnerProcessUnsupported},
		{name: "foreign token", pid: 4242, recorded: "foreign:old"},
		{name: "empty recorded identity", pid: 4242},
		{name: "PID zero"},
		{name: "negative PID", pid: -1},
		{name: "own PID", pid: os.Getpid(), recorded: runtime.GOOS + ":old"},
		{name: "process table failure", pid: 4242, recorded: runtime.GOOS + ":old", absentErr: errors.New("process table denied")},
		{name: "absent after identity read failure", pid: 4242, recorded: runtime.GOOS + ":old", readErr: errors.New("exited"), absentOnRetry: true, gone: true},
	} {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			controller := NewOwnerProcessController()
			reads := 0
			controller.processAbsent = func(int) (bool, error) {
				reads++
				return tt.absent || (tt.absentOnRetry && reads > 1), tt.absentErr
			}
			controller.processStartIdentity = func(context.Context, int) (string, error) {
				return runtime.GOOS + ":live", tt.readErr
			}
			controller.signalProcess = func(int, bool) error {
				t.Fatal("OwnerGone sent a signal")
				return nil
			}
			controller.ownedProcesses = func(int, string) ([]int, error) {
				t.Fatal("OwnerGone enumerated a process tree")
				return nil, nil
			}
			gone, reused := controller.OwnerGone(t.Context(), tt.pid, tt.recorded)
			if gone != tt.gone || reused != tt.reused {
				t.Fatalf("OwnerGone = (%t, %t), want (%t, %t)", gone, reused, tt.gone, tt.reused)
			}
			if tt.pid <= 0 || tt.pid == os.Getpid() {
				gone, reused = OwnerGone(t.Context(), tt.pid, tt.recorded)
				if gone || reused {
					t.Fatalf("default OwnerGone = (%t, %t), want (false, false)", gone, reused)
				}
			}
		})
	}
}
