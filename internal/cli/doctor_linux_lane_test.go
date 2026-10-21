// Boundary: temporary Git repositories and compiled Docker stand-ins only.
package cli

import (
	"bytes"
	"context"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	roundconfig "roundfix/internal/config"
)

func TestDoctorLinuxLaneCheckReportsDockerAndTheImage(t *testing.T) {
	t.Parallel()
	docker, goBinary := linuxLaneTools(t)
	for _, tc := range []struct {
		name, module, image, extra, want string
		noDocker, noRoot                 bool
		args                             [][]string
	}{
		{name: "ok", module: "module fixture\ngo 1.26.1\n", want: "linux-lane: ok (docker reachable; image golang:1.26.1 present)\n", args: [][]string{{"info"}, {"image", "inspect", "golang:1.26.1"}}},
		{name: "image missing", module: "module fixture\ngo 1.26.1\n", extra: "LANE_IMAGE_FAIL=1", want: "linux-lane: warn (DR-LANE-IMAGE-MISSING: image golang:1.26.1 is not present locally; next: docker pull golang:1.26.1)\n", args: [][]string{{"info"}, {"image", "inspect", "golang:1.26.1"}}},
		{name: "docker missing", module: "module fixture\ngo 1.26.1\n", noDocker: true, want: "linux-lane: warn (DR-LANE-DOCKER-MISSING: docker is not on PATH; next: install Docker, or set delivery.linux_lane: false)\n"},
		{name: "daemon unreachable", module: "module fixture\ngo 1.26.1\n", extra: "LANE_INFO_FAIL=1", want: "linux-lane: warn (DR-LANE-DOCKER-UNREACHABLE: docker info failed: daemon unavailable; next: start Docker, then re-run roundfix doctor)\n", args: [][]string{{"info"}}},
		{name: "configured image", module: "module fixture\n", image: "local/lane:ready", want: "linux-lane: ok (docker reachable; image local/lane:ready present)\n", args: [][]string{{"info"}, {"image", "inspect", "local/lane:ready"}}},
		{name: "no go version", module: "module fixture\n", want: "linux-lane: warn (DR-LANE-IMAGE-UNKNOWN: go.mod declares no go version; next: set delivery.linux_lane_image)\n"},
		{name: "no module", want: "linux-lane: skipped (no go.mod at the repository root)\n"},
		{name: "no repository", noRoot: true, want: "linux-lane: skipped (requires a Git repository)\n"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			files := map[string]string{}
			if tc.module != "" {
				files["go.mod"] = tc.module
			}
			f := newLinuxLaneFixture(t, docker, goBinary, files)
			f.workflow.loaded.Config.Delivery.LinuxLaneImage = tc.image
			if tc.noRoot {
				f.workflow.loaded.GitRoot = ""
			}
			if tc.noDocker {
				f.environment = []string{"PATH=" + t.TempDir()}
			}
			if tc.extra != "" {
				f.environment = append(f.environment, tc.extra)
			}
			// A dirty module must not change the image selected from HEAD.
			mustWrite(t, filepath.Join(f.repo, "go.mod"), "module dirty\ngo 1.99\n")
			var output bytes.Buffer
			printDoctorResult(&output, doctorLinuxLaneCheck(f.context(t.Context()), f.workflow.loaded))
			if output.String() != tc.want {
				t.Fatalf("line=%q want=%q", output.String(), tc.want)
			}
			var args [][]string
			for _, start := range f.starts(t) {
				if start.Tool != "docker" {
					t.Fatalf("unexpected tool %q", start.Tool)
				}
				args = append(args, start.Args)
			}
			if !reflect.DeepEqual(args, tc.args) {
				t.Fatalf("arguments=%q want=%q", args, tc.args)
			}
		})
	}
}

func TestDoctorPrintsTheLinuxLaneLineOnlyWhenTheLaneIsOn(t *testing.T) {
	t.Parallel()
	var baseline string
	for _, tc := range []struct {
		name                string
		enabled, dependency bool
	}{
		{"off", false, true}, {"on", true, true}, {"no dependency", true, false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			loaded := roundconfig.Loaded{Config: roundconfig.Builtin(), GitRoot: t.TempDir(), HomeDir: t.TempDir()}
			loaded.Config.Delivery.LinuxLane = tc.enabled
			checker := newDoctorFakeHealthChecker(
				CheckResult{Name: HealthCheckNode, Status: CheckStatusOK},
				CheckResult{Name: HealthCheckACPX, Status: CheckStatusOK},
				CheckResult{Name: HealthCheckCodex, Status: CheckStatusOK},
			)
			withDoctorFakeLoadedAndReadiness(t, checker, loaded, func(context.Context, roundconfig.Config, []roundconfig.WorkCategory, string) profileProofResult {
				return profileProofResult{}
			})
			calls := 0
			if tc.dependency {
				updateCommandDependenciesForTest(t, func(deps *commandDependencies) {
					deps.doctor.linuxLane = func(context.Context, roundconfig.Loaded) CheckResult {
						calls++
						return readinessResult(HealthCheckLinuxLane, "", []readinessFinding{{Code: "DR-LANE-IMAGE-MISSING", Status: CheckStatusWarn, Text: "image golang:1.26.1 is not present locally", Next: "docker pull golang:1.26.1"}})
					}
				})
			}
			var stdout, stderr bytes.Buffer
			if code := runCLI(t, []string{"doctor"}, &stdout, &stderr); code != exitOK || stderr.Len() != 0 {
				t.Fatalf("code=%d stderr=%q", code, stderr.String())
			}
			line := "linux-lane: warn (DR-LANE-IMAGE-MISSING: image golang:1.26.1 is not present locally; next: docker pull golang:1.26.1)\n"
			if tc.enabled && tc.dependency {
				if calls != 1 || !strings.Contains(stdout.String(), "storage: ok (no Run Database)\n"+line+"codex: ok\n") {
					t.Fatalf("calls=%d output=%q", calls, stdout.String())
				}
				if strings.Replace(stdout.String(), line, "", 1) != baseline {
					t.Fatalf("other doctor output changed: %q", stdout.String())
				}
			} else {
				if calls != 0 || strings.Contains(stdout.String(), "linux-lane:") {
					t.Fatalf("calls=%d output=%q", calls, stdout.String())
				}
				if tc.name == "off" {
					baseline = stdout.String()
				} else if stdout.String() != baseline {
					t.Fatalf("missing dependency changed output: %q", stdout.String())
				}
			}
		})
	}
}
