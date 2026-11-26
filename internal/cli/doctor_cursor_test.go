// Suite: Doctor Cursor login diagnostics.
// Invariant: a missing Cursor login fails Doctor with the classification and maintainer action.
// Boundary IN: Doctor command, health checker, adapter check and compiled agent fixture.
// Boundary OUT: other readiness checks, real Cursor service and credentials.
package cli

import (
	"bytes"
	"context"
	"fmt"
	"strings"
	"testing"

	"roundfix/internal/agent"
	roundconfig "roundfix/internal/config"
	"roundfix/internal/testfixture"
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
	for _, test := range []struct {
		name, output string
		exit         int
	}{
		{name: "message", output: "Not logged in\nprivate account sentinel"},
		{name: "exit", output: "private account sentinel", exit: 1},
	} {
		t.Run(test.name, func(t *testing.T) {
			source := fmt.Sprintf(`package main
import ("fmt"; "os")
func main() { fmt.Print(%q); os.Exit(%d) }
`, test.output, test.exit)
			path := testfixture.FixtureBinary(t, "cursor-agent", source)
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
