//go:build unix

package store

import (
	"strings"
	"testing"
	"time"
)

func TestParsePSOwnedProcess(t *testing.T) {
	tests := []struct {
		name        string
		output      string
		wantCPUTime time.Duration
		wantCommand string
		wantErr     string
	}{
		{
			name:        "reads age source CPU time and command",
			output:      "Fri Aug 14 18:56:56 2026  01:02:03.45 roundfix fixture\n",
			wantCPUTime: time.Hour + 2*time.Minute + 3450*time.Millisecond,
			wantCommand: "roundfix fixture",
		},
		{
			name:        "reads CPU time longer than one day",
			output:      "Fri Aug 14 18:56:56 2026  1-02:03:04 roundfix fixture\n",
			wantCPUTime: 26*time.Hour + 3*time.Minute + 4*time.Second,
			wantCommand: "roundfix fixture",
		},
		{
			name:    "rejects incomplete process table entry",
			output:  "Fri Aug 14 18:56:56 2026\n",
			wantErr: "missing command",
		},
		{
			name:    "rejects overflowing CPU time",
			output:  "Fri Aug 14 18:56:56 2026  999999999999999999-00:00:00 roundfix fixture\n",
			wantErr: "overflows",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			process, err := parsePSOwnedProcess([]byte(test.output))
			if test.wantErr != "" {
				if err == nil || !strings.Contains(err.Error(), test.wantErr) {
					t.Fatalf("parse error = %v, want text %q", err, test.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("parse process table entry: %v", err)
			}
			if process.CPUTime != test.wantCPUTime {
				t.Fatalf("CPU time = %s, want %s", process.CPUTime, test.wantCPUTime)
			}
			if process.Command != test.wantCommand {
				t.Fatalf("command = %q, want %q", process.Command, test.wantCommand)
			}
			wantStarted := time.Date(2026, time.August, 14, 18, 56, 56, 0, time.Local)
			if !process.Started.Equal(wantStarted) {
				t.Fatalf("started = %s, want %s", process.Started, wantStarted)
			}
		})
	}
}
