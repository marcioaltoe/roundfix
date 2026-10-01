package cli

// Boundary: real Doctor aggregation; readiness and system checks are fixtures.
// Recommendations must neither open an Agent Session nor fail Doctor.
import (
	"bytes"
	"context"
	"fmt"
	"strings"
	"testing"

	roundconfig "roundfix/internal/config"
)

func doctorBuiltinRecommendationsLine() string {
	return fmt.Sprintf("recommendations: ok (snapshot %s; %d current, 0 differ, 0 pinned)\n", roundconfig.ModelRecommendationSnapshotVersion, len(roundconfig.RequiredWorkCategories()))
}

func TestDoctorReportsRecommendationsAfterProfiles(t *testing.T) {
	t.Parallel()
	for _, status := range []string{"ok", "found", "pinned", "skipped"} {
		t.Run(status, func(t *testing.T) {
			cfg := roundconfig.Builtin()
			if status != "ok" {
				entry := cfg.Profiles[roundconfig.CategoryBackend]
				entry.Profile.Preferred.Model = "chosen-model"
				if status == "pinned" {
					entry.Deviation = &roundconfig.ProfileDeviation{From: roundconfig.ModelRecommendationSnapshotVersion, Reason: "Keep the validated model"}
				}
				if status == "skipped" {
					cfg.Profiles = roundconfig.Profiles{}
				}
				if status != "skipped" {
					cfg.Profiles[roundconfig.CategoryBackend] = entry
				}
			}
			checker := newDoctorFakeHealthChecker(CheckResult{Name: HealthCheckNode, Status: CheckStatusOK}, CheckResult{Name: HealthCheckACPX, Status: CheckStatusOK}, CheckResult{Name: HealthCheckCodex, Status: CheckStatusOK})
			withDoctorFakeLoadedAndReadiness(t, checker, roundconfig.Loaded{Config: cfg, GitRoot: "/repo/project"}, func(context.Context, roundconfig.Config, []roundconfig.WorkCategory, string) profileProofResult {
				return profileProofResult{}
			})
			var stdout, stderr bytes.Buffer
			code := runCLI(t, []string{"doctor"}, &stdout, &stderr)
			wantCode := exitOK
			// Missing required profiles already fail adapter readiness independently.
			if status == "skipped" {
				wantCode = exitRunFailed
			}
			if code != wantCode || stderr.Len() != 0 {
				t.Fatalf("exit=%d stdout=%q stderr=%q", code, &stdout, &stderr)
			}
			lines := strings.Split(stdout.String(), "\n")
			profileIndex := -1
			for i, line := range lines {
				if strings.HasPrefix(line, "profiles:") {
					profileIndex = i
				}
			}
			if profileIndex < 0 || profileIndex+1 >= len(lines) {
				t.Fatalf("missing profiles: %s", &stdout)
			}
			want := strings.TrimSuffix(doctorBuiltinRecommendationsLine(), "\n")
			n := len(roundconfig.RequiredWorkCategories())
			switch status {
			case "found":
				want = fmt.Sprintf("recommendations: found (snapshot %s; %d current, 1 differ, 0 pinned; run roundfix profiles check)", roundconfig.ModelRecommendationSnapshotVersion, n-1)
			case "pinned":
				want = fmt.Sprintf("recommendations: ok (snapshot %s; %d current, 0 differ, 1 pinned)", roundconfig.ModelRecommendationSnapshotVersion, n-1)
			case "skipped":
				result := doctorRecommendationsResult(cfg)
				if result.Status != CheckStatusSkipped || result.Detail == "" {
					t.Fatalf("result=%+v", result)
				}
				want = "recommendations: skipped (" + result.Detail + ")"
			}
			if lines[profileIndex+1] != want {
				t.Fatalf("line=%q want=%q", lines[profileIndex+1], want)
			}
			if len(checker.agentRequests) != 0 {
				t.Fatalf("recommendations opened an Agent Session: %#v", checker.agentRequests)
			}
		})
	}
}

func TestDoctorRecommendationsNeverFailDoctor(t *testing.T) {
	t.Parallel()
	cfg := roundconfig.Builtin()
	entry := cfg.Profiles[roundconfig.CategoryBackend]
	entry.Profile.Preferred.Model = "deliberate-choice"
	cfg.Profiles[roundconfig.CategoryBackend] = entry
	result := doctorRecommendationsResult(cfg)
	if result.Status != CheckStatusFound {
		t.Fatalf("difference status=%q want found", result.Status)
	}
	cfg.Profiles = roundconfig.Profiles{}
	result = doctorRecommendationsResult(cfg)
	if result.Status != CheckStatusSkipped {
		t.Fatalf("uncheckable status=%q want skipped", result.Status)
	}
}
