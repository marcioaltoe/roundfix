package cli

import (
	"bytes"
	"strings"
	"testing"

	roundconfig "roundfix/internal/config"
)

// Characterization corpus for Spec 0088-a-third-runtime-that-can-run.
//
// A test named CharacterizationInvariant pins behavior the Spec must not move.
// A test named CharacterizationToday pins behavior the Spec intends to break,
// so a later Task editing it is a declared break rather than a silent one.

// TestCharacterizationTodayDoctorIgnoresConfiguredOptionalCategory pins the
// diagnostic hole Spec 0088 closes. Measured on 2026-08-08, a configured `data`
// profile whose preferred Agent Selection was failing left the Doctor Command
// printing `profiles: ok (5 distinct tuples; 10 category references)` and naming
// only claude and codex under adapter readiness.
func TestCharacterizationTodayDoctorIgnoresConfiguredOptionalCategory(t *testing.T) {
	stdout, runner := runDoctorWithConfiguredDataProfile(t)

	for _, request := range runner.exactRequests {
		if request.Runtime.ID == "opencode" {
			t.Fatalf("today the Doctor Command proves no configured optional category, got %#v", runner.exactRequests)
		}
	}
	if strings.Contains(stdout, "opencode") {
		t.Fatalf("today the Doctor Command never names an optional category's ACP Runtime: %q", stdout)
	}
	if !strings.Contains(stdout, doctorReadyAdapterLine) {
		t.Fatalf("adapter line = %q, want %q", stdout, doctorReadyAdapterLine)
	}
	if !strings.Contains(stdout, "profiles: ok (") {
		t.Fatalf("today the profiles check reports ok over a configured optional category: %q", stdout)
	}
}

// runDoctorWithConfiguredDataProfile drives the Doctor Command over an
// effective configuration that defines an optional Agent Work Category, and
// returns its stdout with the runner that captured every Exact Agent Selection
// Proof it requested.
func runDoctorWithConfiguredDataProfile(t *testing.T) (string, *profileReadinessExactRunner) {
	t.Helper()
	withCLIWorkspace(t)

	config := roundconfig.Builtin()
	config.Profiles[roundconfig.CategoryData] = roundconfig.ProfileEntry{
		Profile: roundconfig.AgentSelectionProfile{
			Preferred: roundconfig.AgentSelection{Runtime: "opencode", Model: "opencode-go/kimi-k3"},
			Fallbacks: []roundconfig.AgentSelection{
				{Runtime: "claude", Model: "opus", ReasoningEffort: "high"},
			},
		},
		Source: roundconfig.ProfileSourceProject,
	}

	runner := &profileReadinessExactRunner{}
	withAgentRunner(t, runner)

	checker := newDoctorFakeHealthChecker(
		CheckResult{Name: HealthCheckNode, Status: CheckStatusOK, Detail: "node accepted"},
		CheckResult{Name: HealthCheckACPX, Status: CheckStatusOK, Detail: "acpx accepted"},
		CheckResult{Name: HealthCheckCodex, Status: CheckStatusOK, Detail: "codex accepted"},
	)
	withDoctorFakeLoaded(t, checker, roundconfig.Loaded{
		Config:  config,
		GitRoot: "/repo/project",
		HomeDir: "/home/roundfix-test",
	})

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	if code := runCLI(t, []string{"doctor"}, &stdout, &stderr); code != exitOK {
		t.Fatalf("Doctor exit code = %d, want %d; stdout=%q stderr=%q", code, exitOK, stdout.String(), stderr.String())
	}
	if stderr.Len() != 0 {
		t.Fatalf("Doctor stderr = %q, want empty", stderr.String())
	}
	return stdout.String(), runner
}
