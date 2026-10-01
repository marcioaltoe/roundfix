// Suite: Doctor adapter floors.
// Invariant: official adapters below the floor fail with an install action; the floor passes.
// Boundary IN: real adapter version inspection, health conversion and Doctor aggregation.
// Boundary OUT: real adapters, providers and network; version probes use local shell fixtures.
package cli

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"roundfix/internal/agent"
)

func TestDoctorRefusesAnAdapterBelowTheFloor(t *testing.T) {
	t.Parallel()
	for _, tt := range []struct{ runtime, version, floor, pkg string }{
		{"codex", "1.1.5", agent.PinnedCodexAdapterVersion, agent.CodexAdapterPackage},
		{"claude", "0.63.0", agent.PinnedClaudeAdapterVersion, agent.ClaudeAdapterPackage},
	} {
		t.Run(tt.runtime, func(t *testing.T) {
			result := doctorWithVersionFixture(t, tt.runtime, tt.pkg, tt.version)
			if result.Status != CheckStatusFailed {
				t.Fatalf("Doctor status = %s, want failed: %#v", result.Status, result)
			}
			want := "npm install -g " + tt.pkg + "@" + tt.floor
			if result.NextAction != want {
				t.Fatalf("install action = %q, want %q", result.NextAction, want)
			}
			for _, part := range []string{tt.version, tt.floor, agent.AdapterVersionUnsupported} {
				if !strings.Contains(result.Detail, part) {
					t.Errorf("Doctor detail %q omits %q", result.Detail, part)
				}
			}
		})
	}
}

func TestDoctorAcceptsAnAdapterAtTheFloor(t *testing.T) {
	t.Parallel()
	for _, tt := range []struct{ runtime, floor, pkg string }{
		{"codex", agent.PinnedCodexAdapterVersion, agent.CodexAdapterPackage},
		{"claude", agent.PinnedClaudeAdapterVersion, agent.ClaudeAdapterPackage},
	} {
		t.Run(tt.runtime, func(t *testing.T) {
			result := doctorWithVersionFixture(t, tt.runtime, tt.pkg, tt.floor)
			if result.Status != CheckStatusOK || result.NextAction != "" {
				t.Fatalf("Doctor = %#v, want ok without remediation", result)
			}
			if !strings.Contains(result.Detail, tt.pkg) || !strings.Contains(result.Detail, "version="+tt.floor) {
				t.Fatalf("missing official floor evidence: %#v", result)
			}
		})
	}
}

func doctorWithVersionFixture(t *testing.T, runtime, pkg, version string) CheckResult {
	t.Helper()
	command := filepath.Join(t.TempDir(), "fake-adapter")
	script := "#!/bin/sh\nif [ \"$1\" != \"--version\" ]; then exit 91; fi\nprintf '%s\\n' '" + pkg + " " + version + "'\n"
	if err := os.WriteFile(command, []byte(script), 0o700); err != nil {
		t.Fatal(err)
	}
	checker := newHealthChecker(healthCheckDependencies{checkAdapter: agent.CheckAdapter})
	return doctorAdapterCheck(context.Background(), checker, []agent.RuntimeSpec{{ID: runtime, Protocol: agent.ProtocolStdio, Command: command}}, nil)
}
