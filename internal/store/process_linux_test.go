//go:build linux

package store

import (
	"os"
	"strings"
	"syscall"
	"testing"
)

func TestParseProcStatStartTimeCountsFromLastClosingParenthesis(t *testing.T) {
	t.Parallel()
	stat := []byte("42 (worker ) pool) S 1 2 3 4 5 6 7 8 9 10 11 12 13 14 15 16 17 18 424242 20")

	startTime, err := parseProcStatStartTime(stat)
	if err != nil {
		t.Fatalf("parse process stat start time: %v", err)
	}
	if startTime != "424242" {
		t.Fatalf("process start time = %q, want 424242", startTime)
	}
}

func TestParseProcStatStartTimeRejectsMissingField(t *testing.T) {
	t.Parallel()
	_, err := parseProcStatStartTime([]byte("42 (worker) S 1 2"))
	if err == nil || !strings.Contains(err.Error(), "missing start time field") {
		t.Fatalf("parse error = %v, want missing start time field", err)
	}
}

func TestProcessVanishedAcceptsBothProcfsAnswers(t *testing.T) {
	// A task reaped mid-scan is reported as ENOENT or, while it is still being
	// reaped, as ESRCH. Guarding only the first made a whole process-table
	// enumeration fail on a process that simply exited.
	for _, testCase := range []struct {
		name string
		err  error
		want bool
	}{
		{name: "missing proc entry", err: &os.PathError{Op: "read", Path: "/proc/1/stat", Err: syscall.ENOENT}, want: true},
		{name: "task being reaped", err: &os.PathError{Op: "read", Path: "/proc/1/stat", Err: syscall.ESRCH}, want: true},
		{name: "permission denied is not a vanished process", err: &os.PathError{Op: "read", Path: "/proc/1/stat", Err: syscall.EACCES}, want: false},
		{name: "nil is not a vanished process", err: nil, want: false},
	} {
		t.Run(testCase.name, func(t *testing.T) {
			if got := processVanished(testCase.err); got != testCase.want {
				t.Fatalf("processVanished(%v) = %t, want %t", testCase.err, got, testCase.want)
			}
		})
	}
}
