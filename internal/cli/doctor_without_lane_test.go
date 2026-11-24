// Boundary: doctor with legacy checkout configuration and fake readiness checks.
package cli

import (
	"bytes"
	"context"
	"os"
	"path/filepath"
	roundconfig "roundfix/internal/config"
	"roundfix/internal/gittest"
	"strings"
	"testing"
)

func TestDoctorRunsNoLinuxLaneCheck(t *testing.T) {
	t.Parallel()
	home, repo := withCLIWorkspace(t)
	gittest.InitRepo(t, repo, "-b", "main")
	record := recordingRetiredDocker(t)
	mustWrite(t, filepath.Join(repo, ".roundfixrc.yml"), "delivery:\n  linux_lane: true\n")
	var diagnostics bytes.Buffer
	loaded, err := roundconfig.Load(roundconfig.LoadOptions{WorkDir: repo, HomeDir: home, Stderr: &diagnostics})
	if err != nil {
		t.Fatal(err)
	}
	checker := newDoctorFakeHealthChecker(CheckResult{Name: HealthCheckNode, Status: CheckStatusOK}, CheckResult{Name: HealthCheckACPX, Status: CheckStatusOK}, CheckResult{Name: HealthCheckCodex, Status: CheckStatusOK})
	withDoctorFakeLoadedAndReadiness(t, checker, loaded, func(context.Context, roundconfig.Config, []roundconfig.WorkCategory, string) profileProofResult {
		return profileProofResult{}
	})
	updateCommandDependenciesForTest(t, func(d *commandDependencies) {
		d.doctor.verificationCores = func(context.Context, roundconfig.Loaded) CheckResult {
			return CheckResult{Name: HealthCheckVerificationCores, Status: CheckStatusOK}
		}
	})
	var out, stderr bytes.Buffer
	code := runCLI(t, []string{"doctor"}, &out, &stderr)
	if code != 0 || stderr.Len() != 0 || strings.Contains(out.String(), "linux-lane:") || strings.Contains(out.String(), "DR-LANE-") || !strings.Contains(out.String(), "verification-cores: ok\ncodex: ok\n") {
		t.Fatalf("exit=%d stdout=%q stderr=%q", code, out.String(), stderr.String())
	}
	if _, err := os.Stat(record); !os.IsNotExist(err) {
		t.Fatalf("docker started: %v", err)
	}
}
