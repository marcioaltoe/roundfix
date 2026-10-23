package cli

import (
	"bytes"
	"context"
	"strings"
	"testing"

	"roundfix/internal/agent"
	roundconfig "roundfix/internal/config"
)

const doctorLightTierWarningLine = "light-tier: warn (DR-LIGHT-ACCESS-MODE: light runtime opencode offers no full-access mode; light Tasks run on their category profile; next: set openrouter.light_models: [] in User Config, or defaults.agent_full_access: false)\n"

func TestDoctorLightTierResultNamesTheLightRuntimeThatCannotServeFullAccess(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		name                             string
		fullAccess, emptyModels, serving bool
		want                             string
	}{
		{name: "full access off"},
		{name: "no light model", fullAccess: true, emptyModels: true},
		{name: "static runtime", fullAccess: true, want: doctorLightTierWarningLine},
		{name: "serving runtime", fullAccess: true, serving: true, want: "light-tier: ok (light runtime opencode serves full access)\n"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			config := roundconfig.Builtin()
			config.Defaults.AgentFullAccess = tc.fullAccess
			if tc.emptyModels {
				config.OpenRouter.LightModels = nil
			}
			factory := func(selection roundconfig.AgentSelection) (agent.RuntimeSpec, error) {
				if tc.want == "" {
					t.Fatal("inactive check must not build a runtime")
				}
				runtime, err := runtimeForProfileSelectionWithOptions(selection, profileProofOptions{EnableFullAccess: config.Defaults.AgentFullAccess})
				if tc.serving {
					runtime.SupportedFullAccessMode = "full-access"
				}
				return runtime, err
			}
			result, printed := doctorLightTierResult(config, factory)
			if printed != (tc.want != "") {
				t.Fatalf("printed=%v want=%v", printed, tc.want != "")
			}
			if printed {
				var output bytes.Buffer
				printDoctorResult(&output, result)
				if output.String() != tc.want {
					t.Fatalf("line=%q want=%q", output.String(), tc.want)
				}
			}
		})
	}
}

func TestDoctorPrintsTheLightTierLineOnlyWithFullAccessAndLightModels(t *testing.T) {
	t.Parallel()
	var baseline string
	for _, tc := range []struct {
		name                    string
		fullAccess, emptyModels bool
	}{
		{name: "built-in"},
		{name: "full access", fullAccess: true},
		{name: "no light model", fullAccess: true, emptyModels: true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			loaded := roundconfig.Loaded{Config: roundconfig.Builtin(), GitRoot: t.TempDir(), HomeDir: t.TempDir()}
			loaded.Config.Defaults.AgentFullAccess = tc.fullAccess
			if tc.emptyModels {
				loaded.Config.OpenRouter.LightModels = nil
			}
			checker := newDoctorFakeHealthChecker(
				CheckResult{Name: HealthCheckNode, Status: CheckStatusOK},
				CheckResult{Name: HealthCheckACPX, Status: CheckStatusOK},
				CheckResult{Name: HealthCheckCodex, Status: CheckStatusOK},
			)
			withDoctorFakeLoadedAndReadiness(t, checker, loaded, func(context.Context, roundconfig.Config, []roundconfig.WorkCategory, string) profileProofResult {
				return profileProofResult{}
			})
			var stdout, stderr bytes.Buffer
			if code := runCLI(t, []string{"doctor"}, &stdout, &stderr); code != exitOK || stderr.Len() != 0 {
				t.Fatalf("code=%d stderr=%q", code, stderr.String())
			}
			output := stdout.String()
			if tc.fullAccess && !tc.emptyModels {
				if strings.Count(output, "light-tier:") != 1 || !strings.Contains(output, doctorBuiltinRecommendationsLine()+doctorLightTierWarningLine+"pre-pr-review:") {
					t.Fatalf("light-tier line missing or out of order: %q", output)
				}
				if strings.Replace(output, doctorLightTierWarningLine, "", 1) != baseline {
					t.Fatalf("other doctor output changed: %q", output)
				}
			} else {
				if strings.Contains(output, "light-tier:") {
					t.Fatalf("unexpected light-tier line: %q", output)
				}
				if tc.name == "built-in" {
					baseline = output
				} else if output != baseline {
					t.Fatalf("other doctor output changed: %q", output)
				}
			}
		})
	}
}
