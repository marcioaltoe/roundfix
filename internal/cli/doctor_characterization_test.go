package cli

import (
	"bytes"
	"strings"
	"testing"

	"roundfix/internal/agent"
	roundconfig "roundfix/internal/config"
)

// Characterization corpus for Spec 0088-a-third-runtime-that-can-run.
//
// A test named CharacterizationInvariant pins behavior the Spec must not move.
// A test named CharacterizationDeclaredBreak records behavior the Spec changed
// on purpose, keeping the superseded behavior and its provenance in the comment.

// TestCharacterizationDeclaredBreakDoctorProvesConfiguredOptionalCategory
// records the third break this Spec declares. Until Task 04, a configured
// `data` profile whose preferred Agent Selection was failing left the Doctor
// Command printing `profiles: ok (5 distinct tuples; 10 category references)`
// and naming only claude and codex under adapter readiness — measured on
// 2026-08-08, and the reason an OpenCode profile looked like it had never
// registered. Readiness now covers every Agent Work Category the effective
// configuration defines.
func TestCharacterizationDeclaredBreakDoctorProvesConfiguredOptionalCategory(t *testing.T) {
	stdout, runner := runDoctorWithConfiguredDataProfile(t)

	proven := false
	for _, request := range runner.exactRequests {
		if request.Runtime.ID == "opencode" && request.Runtime.Model == "opencode-go/kimi-k3" {
			proven = true
		}
	}
	if !proven {
		t.Fatalf("the configured optional category's Agent Selection was not proven: %#v", runner.exactRequests)
	}
	if !strings.Contains(stdout, "opencode:") {
		t.Fatalf("adapter readiness must name the configured optional category's ACP Runtime: %q", stdout)
	}
	if !strings.Contains(stdout, "profiles: ok (5 distinct tuples; 12 category references)") {
		t.Fatalf("profiles line must count the configured optional category: %q", stdout)
	}
}

// TestCharacterizationInvariantDoctorCountsAreUnchangedWithoutOptionalCategories
// keeps the widened scope from inventing proofs: a configuration defining no
// optional Agent Work Category reports exactly what it reported before.
func TestCharacterizationInvariantDoctorCountsAreUnchangedWithoutOptionalCategories(t *testing.T) {
	stdout, runner := runDoctorWithConfig(t, roundconfig.Builtin())

	if len(runner.exactRequests) != 3 {
		t.Fatalf("expected three distinct profile proofs, got %#v", runner.exactRequests)
	}
	if !strings.Contains(stdout, "profiles: ok (3 distinct tuples; 10 category references)") {
		t.Fatalf("profiles line changed for a configuration with no optional category: %q", stdout)
	}
	if strings.Contains(stdout, "opencode:") {
		t.Fatalf("adapter readiness must not name an ACP Runtime no configured tuple references: %q", stdout)
	}
}

// TestCharacterizationInvariantInheritedCategoryAddsNoTuple keeps an optional
// category that only inherits general out of readiness, because it contributes
// no distinct Agent Selection tuple.
func TestCharacterizationInvariantInheritedCategoryAddsNoTuple(t *testing.T) {
	config := roundconfig.Builtin()
	if _, defined := config.Profiles[roundconfig.CategoryDocs]; defined {
		t.Fatal("the built-in configuration must not define an optional category")
	}

	resolved, err := roundconfig.ResolveProfile(config, roundconfig.CategoryDocs, nil)
	if err != nil {
		t.Fatalf("an optional category must still resolve by inheritance: %v", err)
	}
	if resolved.InheritedFrom != roundconfig.CategoryGeneral {
		t.Fatalf("inherited from = %q, want %q", resolved.InheritedFrom, roundconfig.CategoryGeneral)
	}

	stdout, _ := runDoctorWithConfig(t, config)
	if !strings.Contains(stdout, "profiles: ok (3 distinct tuples; 10 category references)") {
		t.Fatalf("an inherited category must add no reference: %q", stdout)
	}
}

// runDoctorWithConfiguredDataProfile drives the Doctor Command over an
// effective configuration that defines an optional Agent Work Category, and
// returns its stdout with the runner that captured every Exact Agent Selection
// Proof it requested.
func runDoctorWithConfiguredDataProfile(t *testing.T) (string, *profileReadinessExactRunner) {
	t.Helper()

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
	return runDoctorWithConfig(t, config)
}

func runDoctorWithConfig(t *testing.T, config roundconfig.Config) (string, *profileReadinessExactRunner) {
	t.Helper()
	withCLIWorkspace(t)

	runner := &profileReadinessExactRunner{}
	withAgentRunner(t, runner)

	checker := newDoctorFakeHealthChecker(
		CheckResult{Name: HealthCheckNode, Status: CheckStatusOK, Detail: "node accepted"},
		CheckResult{Name: HealthCheckACPX, Status: CheckStatusOK, Detail: "acpx accepted"},
		CheckResult{Name: HealthCheckCodex, Status: CheckStatusOK, Detail: "codex accepted"},
	)
	checker.adapterResults = map[string]CheckResult{
		"claude":   {Name: HealthCheckAdapter, Status: CheckStatusOK, Detail: "claude-agent-acp"},
		"codex":    {Name: HealthCheckAdapter, Status: CheckStatusOK, Detail: "codex-acp"},
		"opencode": {Name: HealthCheckAdapter, Status: CheckStatusOK, Detail: "opencode"},
	}
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

// TestDoctorProfileReadinessFailsOnAConfiguredOptionalCategory is the outcome
// the widened scope exists for: a readiness command never reports ok while a
// configured Agent Selection Profile is failing.
func TestDoctorProfileReadinessFailsOnAConfiguredOptionalCategory(t *testing.T) {
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

	runner := &profileReadinessExactRunner{
		prove: func(req agent.ProbeRequest) (agent.SelectionProof, error) {
			if req.Runtime.ID == "opencode" {
				return agent.SelectionProof{}, &agent.SelectionUnsupportedError{
					Kind:    agent.SelectionModelNotAdvertised,
					Runtime: "opencode",
					Model:   req.Runtime.Model,
				}
			}
			return agent.SelectionProof{}, nil
		},
	}
	withAgentRunner(t, runner)

	checker := newDoctorFakeHealthChecker(
		CheckResult{Name: HealthCheckNode, Status: CheckStatusOK, Detail: "node accepted"},
		CheckResult{Name: HealthCheckACPX, Status: CheckStatusOK, Detail: "acpx accepted"},
		CheckResult{Name: HealthCheckCodex, Status: CheckStatusOK, Detail: "codex accepted"},
	)
	checker.adapterResults = map[string]CheckResult{
		"claude":   {Name: HealthCheckAdapter, Status: CheckStatusOK, Detail: "claude-agent-acp"},
		"codex":    {Name: HealthCheckAdapter, Status: CheckStatusOK, Detail: "codex-acp"},
		"opencode": {Name: HealthCheckAdapter, Status: CheckStatusOK, Detail: "opencode"},
	}
	withDoctorFakeLoaded(t, checker, roundconfig.Loaded{
		Config:  config,
		GitRoot: "/repo/project",
		HomeDir: "/home/roundfix-test",
	})

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	code := runCLI(t, []string{"doctor"}, &stdout, &stderr)

	if code == exitOK {
		t.Fatalf("Doctor must fail while a configured profile is failing: %q", stdout.String())
	}
	if !strings.Contains(stdout.String(), "profiles: failed") {
		t.Fatalf("profiles check must report failed: %q", stdout.String())
	}
	if !strings.Contains(stdout.String(), "data") {
		t.Fatalf("the failing category must be named: %q", stdout.String())
	}
}
