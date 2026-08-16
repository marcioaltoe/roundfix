// Suite: diagnostic signature normalization
// Invariant: only declared volatile diagnostic spans may vary without changing a command-scoped signature.
// Boundary IN: DiagnosticSignature command and diagnostic inputs.
// Boundary OUT: Run Event Journal lookup and repeated-failure reporting.
package daemon

import "testing"

func TestDiagnosticSignature(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		commandA    string
		diagnosticA string
		commandB    string
		diagnosticB string
		wantEqual   bool
	}{
		{
			name:        "normalizes timestamp",
			commandA:    "go test ./internal/daemon",
			diagnosticA: "2026-08-16T12:34:56.123Z assertion failed: got 1, want 2",
			commandB:    "go test ./internal/daemon",
			diagnosticB: "2027-01-02T03:04:05.987-03:00 assertion failed: got 1, want 2",
			wantEqual:   true,
		},
		{
			name:        "normalizes duration",
			commandA:    "go test ./internal/daemon",
			diagnosticA: "assertion failed after 12.4ms: got 1, want 2",
			commandB:    "go test ./internal/daemon",
			diagnosticB: "assertion failed after 9.2s: got 1, want 2",
			wantEqual:   true,
		},
		{
			name:        "normalizes temporary directory path",
			commandA:    "go test ./internal/daemon",
			diagnosticA: "assertion failed; diagnostics: /tmp/roundfix-a123/verification.log",
			commandB:    "go test ./internal/daemon",
			diagnosticB: "assertion failed; diagnostics: /private/tmp/roundfix-b456/verification.log",
			wantEqual:   true,
		},
		{
			name:        "normalizes process identifier",
			commandA:    "go test ./internal/daemon",
			diagnosticA: "worker pid=1234 assertion failed: got 1, want 2",
			commandB:    "go test ./internal/daemon",
			diagnosticB: "worker PID: 98765 assertion failed: got 1, want 2",
			wantEqual:   true,
		},
		{
			name:        "normalizes run identifier",
			commandA:    "go test ./internal/daemon",
			diagnosticA: "run run_20260816T123456Z_a1b2c3d4 assertion failed: got 1, want 2",
			commandB:    "go test ./internal/daemon",
			diagnosticB: "run run_20270102T030405Z_9f8e7d6c assertion failed: got 1, want 2",
			wantEqual:   true,
		},
		{
			name:        "preserves different assertion",
			commandA:    "go test ./internal/daemon",
			diagnosticA: "assertion failed: got 1, want 2",
			commandB:    "go test ./internal/daemon",
			diagnosticB: "assertion failed: got 3, want 2",
			wantEqual:   false,
		},
		{
			name:        "includes failing command",
			commandA:    "go test ./internal/daemon",
			diagnosticA: "assertion failed: got 1, want 2",
			commandB:    "go test ./internal/store",
			diagnosticB: "assertion failed: got 1, want 2",
			wantEqual:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			gotA := DiagnosticSignature(tt.commandA, []byte(tt.diagnosticA))
			gotB := DiagnosticSignature(tt.commandB, []byte(tt.diagnosticB))
			if (gotA == gotB) != tt.wantEqual {
				t.Fatalf("DiagnosticSignature() equality = %t, want %t; first=%q second=%q", gotA == gotB, tt.wantEqual, gotA, gotB)
			}
		})
	}
}
