package cli

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"slices"
	"strings"
	"testing"

	"roundfix/internal/baseline"
	roundconfig "roundfix/internal/config"
	"roundfix/internal/gittest"
	"roundfix/skills"
)

func setDoctorTestManifestProfile(t *testing.T, root, profile string) {
	t.Helper()
	path := filepath.Join(root, filepath.FromSlash(doctorSetupManifestPath))
	var manifest map[string]any
	if err := json.Unmarshal(mustReadBytes(t, path), &manifest); err != nil {
		t.Fatal(err)
	}
	manifest["profile"] = profile
	data, err := json.Marshal(manifest)
	if err != nil {
		t.Fatal(err)
	}
	mustWrite(t, path, string(data))
}

func TestDoctorChecksAnExternalSkillOnlyForPresence(t *testing.T) {
	t.Parallel()
	for _, test := range []struct {
		name    string
		missing bool
	}{
		{name: "differs from its Setup Snapshot"},
		{name: "missing", missing: true},
	} {
		t.Run(test.name, func(t *testing.T) {
			root := t.TempDir()
			gittest.MarkRepository(t, root)
			writeDoctorReadyRepositoryFixture(t, root)
			setDoctorTestManifestProfile(t, root, "go-cli-tui")
			if test.missing {
				if err := os.RemoveAll(filepath.Join(root, ".agents", "skills", "testing-boss")); err != nil {
					t.Fatal(err)
				}
			}
			checker := newDoctorFakeHealthChecker(CheckResult{Name: HealthCheckNode, Status: CheckStatusOK}, CheckResult{Name: HealthCheckACPX, Status: CheckStatusOK}, CheckResult{Name: HealthCheckCodex, Status: CheckStatusOK})
			withDoctorFakeLoadedAndReadiness(t, checker, roundconfig.Loaded{Config: roundconfig.Builtin(), GitRoot: root}, func(context.Context, roundconfig.Config, []roundconfig.WorkCategory, string) profileProofResult {
				return profileProofResult{}
			})
			updateCommandDependenciesForTest(t, func(dependencies *commandDependencies) {
				fixture := dependencies.doctor
				dependencies.doctor = defaultDoctorDependencies()
				dependencies.doctor.loadConfig = fixture.loadConfig
				dependencies.doctor.readiness = fixture.readiness
				dependencies.doctor.healthChecker = fixture.healthChecker
				dependencies.doctor.profileReadiness = fixture.profileReadiness
				dependencies.doctor.residue = fixture.residue
				dependencies.doctor.storage = fixture.storage
				dependencies.doctor.diskSpace = nil
				dependencies.doctor.linuxLane = nil
				dependencies.doctor.resolveExternal = func(string) ([]string, bool, error) { return []string{"testing-boss"}, true, nil }
				dependencies.doctor.checkSkills = skills.CheckRepositoryWithExternal
			})
			before := snapshotDoctorPath(t, root)
			var stdout, stderr bytes.Buffer
			code := runCLI(t, []string{"doctor"}, &stdout, &stderr)
			owned := len(skills.Names())
			want := fmt.Sprintf("skills: ok (%d required: %d Roundfix-owned, 1 external)", owned+1, owned)
			wantCode := exitOK
			if test.missing {
				want = "skills: failed (missing: testing-boss; next: bunx skills add marcioaltoe/skills@testing-boss)"
				wantCode = exitRunFailed
			}
			if code != wantCode || stderr.Len() != 0 || !slices.Contains(strings.Split(stdout.String(), "\n"), want) || strings.Contains(stdout.String(), "DR-SKILL-DIFFERS-FROM-SNAPSHOT") {
				t.Fatalf("exit=%d stdout=%s stderr=%s; want exit=%d line=%s", code, stdout.String(), stderr.String(), wantCode, want)
			}
			if !reflect.DeepEqual(before, snapshotDoctorPath(t, root)) {
				t.Fatal("Doctor changed repository bytes")
			}
		})
	}
}

func TestBaselineUpdateNeverRestoresAPresentExternalSkill(t *testing.T) {
	t.Parallel()
	for _, test := range []struct {
		name    string
		missing bool
	}{
		{name: "present copy that differs from its Setup Snapshot"},
		{name: "missing copy", missing: true},
	} {
		t.Run(test.name, func(t *testing.T) {
			root := t.TempDir()
			gittest.MarkRepository(t, root)
			writeDoctorReadyRepositoryFixture(t, root)
			setDoctorTestManifestProfile(t, root, "go-cli-tui")
			path := filepath.Join(root, ".agents", "skills", "testing-boss", "SKILL.md")
			before := mustReadBytes(t, path)
			if test.missing {
				if err := os.RemoveAll(filepath.Dir(path)); err != nil {
					t.Fatal(err)
				}
			}
			var requests []baseline.SkillsRestoreRequest
			result, err := runBaselineUpdateSkillsStageWith(t.Context(), baselineUpdateSkillsRequest{
				Repository: root, ProfileID: "go-cli-tui",
			}, baselineUpdateSkillsDependencies{
				recordLock:         baseline.RecordInstalledOwnedSkills,
				resolveProjectRoot: defaultResolveSkillsProjectRoot,
				install:            skills.Install,
				ownedNames:         skills.Names,
				resolveExternal:    func(string) ([]string, bool, error) { return []string{"testing-boss"}, true, nil },
				checkRepository:    skills.CheckRepositoryWithExternal,
				restore: func(_ context.Context, request baseline.SkillsRestoreRequest) (baseline.SkillsRestorePayload, error) {
					requests = append(requests, request)
					if request.Repository != root || request.ProfileID != "go-cli-tui" || !slices.Equal(request.Skills, []string{"testing-boss"}) {
						t.Fatalf("restore request=%+v", request)
					}
					digest := strings.Repeat("a", 64)
					payload := baseline.SkillsRestorePayload{PlanDigest: &digest}
					if request.Confirmation == "" {
						payload.Finding = &baseline.RestoreFinding{Code: "plan.confirmation.required"}
						return payload, &baseline.SkillsRestoreError{Category: baseline.SkillsRestoreAction, Finding: *payload.Finding}
					}
					if request.Confirmation != digest {
						t.Fatalf("restore confirmation=%q, want %q", request.Confirmation, digest)
					}
					payload.Applied = true
					return payload, nil
				},
			})
			if err != nil || result.Status != baselineUpdateSkillsVerified || len(result.Drifted) != 0 {
				t.Fatalf("result=%+v err=%v", result, err)
			}
			data, err := json.Marshal(result)
			if err != nil || bytes.Contains(data, []byte(`"differing"`)) {
				t.Fatalf("JSON=%s err=%v", data, err)
			}
			if test.missing {
				if !slices.Equal(result.Restored, []string{"testing-boss"}) || len(requests) != 2 {
					t.Fatalf("result=%+v requests=%+v", result, requests)
				}
			} else {
				if len(result.Restored) != 0 || len(requests) != 0 {
					t.Fatalf("present external skill restored: result=%+v requests=%+v", result, requests)
				}
				if !bytes.Equal(before, mustReadBytes(t, path)) {
					t.Fatal("refresh changed the present skill's bytes")
				}
			}
		})
	}
}
