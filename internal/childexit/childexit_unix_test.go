//go:build darwin || linux

// Boundary: constructed wait statuses only; no process is started.
package childexit

import (
	"errors"
	"syscall"
	"testing"
)

func TestExitDescribesHowAChildEnded(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		name        string
		status      syscall.WaitStatus
		useStatus   bool
		err         error
		want        Exit
		ok          bool
		description string
	}{
		{name: "exit 2", status: syscall.WaitStatus(2 << 8), useStatus: true, want: Exit{Code: 2}, ok: true, description: "exited 2"},
		{name: "SIGKILL without core", status: syscall.WaitStatus(syscall.SIGKILL), useStatus: true, want: Exit{Code: -1, Signal: "SIGKILL", Number: int(syscall.SIGKILL)}, ok: true, description: "was killed by signal SIGKILL (9), no core dump"},
		{name: "SIGABRT with core", status: syscall.WaitStatus(syscall.SIGABRT) | 0x80, useStatus: true, want: Exit{Code: -1, Signal: "SIGABRT", Number: int(syscall.SIGABRT), CoreDumped: true}, ok: true, description: "was killed by signal SIGABRT (6), core dumped"},
		{name: "nil error", ok: true, description: "exited 0"},
		{name: "non-exit error", err: errors.New("child did not start")},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			var got Exit
			var ok bool
			if tc.useStatus {
				got, ok = fromWaitStatus(tc.status), true
			} else {
				got, ok = From(tc.err)
			}
			if got != tc.want || ok != tc.ok {
				t.Fatalf("exit=%+v ok=%t want=%+v ok=%t", got, ok, tc.want, tc.ok)
			}
			if got.Signaled() != (tc.want.Number != 0) {
				t.Fatalf("Signaled()=%t want=%t", got.Signaled(), tc.want.Number != 0)
			}
			if tc.ok && got.String() != tc.description {
				t.Fatalf("String()=%q want=%q", got.String(), tc.description)
			}
		})
	}
}
