// Suite: Doctor Cursor login diagnostics.
// Invariant: a missing Cursor login fails Doctor with the classification and maintainer action.
// Boundary IN: Doctor command, health checker, adapter check and compiled agent fixture.
// Boundary OUT: other readiness checks, real Cursor service and credentials.
package cli

import (
	"bytes"
	"context"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"roundfix/internal/agent"
	roundconfig "roundfix/internal/config"
)

// Keep the other Doctor boundaries offline while exercising the real adapter check.
type cursorDoctorHealthChecker struct {
	*doctorFakeHealthChecker
	adapter HealthChecker
}

func (checker cursorDoctorHealthChecker) Adapter(ctx context.Context, runtime agent.RuntimeSpec) CheckResult {
	return checker.adapter.Adapter(ctx, runtime)
}

func TestDoctorAdapterNamesTheCursorLogin(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	binary := filepath.Join(dir, "agent.test")
	build := exec.CommandContext(t.Context(), "go", "test", "-c", "-o", binary, "../agent")
	if output, err := build.CombinedOutput(); err != nil {
		t.Fatalf("compile agent fixture: %v\n%s", err, output)
	}
	for _, test := range []struct {
		name, output string
		exit         int
	}{
		{name: "message", output: "Not logged in\nprivate account sentinel"},
		{name: "exit", output: "private account sentinel", exit: 1},
	} {
		t.Run(test.name, func(t *testing.T) {
			path := filepath.Join(t.TempDir(), "cursor-agent")
			provision := exec.CommandContext(t.Context(), binary)
			provision.Env = append(os.Environ(), "ROUNDFIX_FAKE_ADAPTER_PROVISION="+path, "ROUNDFIX_FAKE_ADAPTER_ARGV0=true")
			if output, err := provision.CombinedOutput(); err != nil {
				t.Fatalf("provision agent fixture: %v\n%s", err, output)
			}
			payload, err := json.Marshal(map[string]any{"status_output": test.output, "status_exit_code": test.exit})
			if err != nil {
				t.Fatal(err)
			}
			if err := os.WriteFile(path+".fixture.json", payload, 0o600); err != nil {
				t.Fatal(err)
			}
			checker := cursorDoctorHealthChecker{
				doctorFakeHealthChecker: newDoctorFakeHealthChecker(
					CheckResult{Name: HealthCheckNode, Status: CheckStatusOK},
					CheckResult{Name: HealthCheckACPX, Status: CheckStatusOK},
					CheckResult{Name: HealthCheckCodex, Status: CheckStatusOK},
				),
				adapter: newHealthChecker(healthCheckDependencies{checkAdapter: func(ctx context.Context, runtime agent.RuntimeSpec) (agent.AdapterEvidence, error) {
					runtime.Protocol = agent.ProtocolStdio
					runtime.Command = path
					return agent.CheckAdapter(ctx, runtime)
				}}),
			}
			config := roundconfig.Builtin()
			for category, entry := range config.Profiles {
				entry.Profile.Preferred = roundconfig.AgentSelection{Runtime: "cursor", Model: "grok-4-20[thinking=true]"}
				entry.Profile.Fallbacks = []roundconfig.AgentSelection{{Runtime: "cursor", Model: "default[]"}}
				config.Profiles[category] = entry
			}
			withDoctorFakeLoadedAndReadiness(t, checker, roundconfig.Loaded{Config: config, GitRoot: dir}, func(context.Context, roundconfig.Config, []roundconfig.WorkCategory, string) profileProofResult {
				return profileProofResult{}
			})
			var stdout, stderr bytes.Buffer
			code := runCLI(t, []string{"doctor"}, &stdout, &stderr)
			want := "adapter: failed (cursor: cursor-agent is not logged in; classification: " + agent.CursorLoginRequired + "; next: run cursor-agent login in a terminal yourself)"
			if code != exitRunFailed || stderr.Len() != 0 || !strings.Contains(stdout.String(), want+"\n") {
				t.Fatalf("code=%d stdout=%q stderr=%q; want %q", code, stdout.String(), stderr.String(), want)
			}
			if strings.Contains(stdout.String()+stderr.String(), "private account sentinel") {
				t.Fatal("Doctor exposed status output")
			}
		})
	}
}
