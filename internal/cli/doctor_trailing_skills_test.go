package cli

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"roundfix/internal/baseline"
	roundconfig "roundfix/internal/config"
	"roundfix/skills"
)

func TestDoctorWarnsOnASkillThatTrailsItsSnapshot(t *testing.T) {
	t.Parallel()
	root := t.TempDir()
	writeDoctorReadyRepositoryFixture(t, root)
	manifestPath := filepath.Join(root, filepath.FromSlash(doctorSetupManifestPath))
	var manifest map[string]any
	if err := json.Unmarshal(mustReadBytes(t, manifestPath), &manifest); err != nil {
		t.Fatal(err)
	}
	manifest["profile"] = "go-cli-tui"
	data, err := json.Marshal(manifest)
	if err != nil {
		t.Fatal(err)
	}
	mustWrite(t, manifestPath, string(data))
	checker := newDoctorFakeHealthChecker(CheckResult{Name: HealthCheckNode, Status: CheckStatusOK}, CheckResult{Name: HealthCheckACPX, Status: CheckStatusOK}, CheckResult{Name: HealthCheckCodex, Status: CheckStatusOK})
	withDoctorFakeLoadedAndReadiness(t, checker, roundconfig.Loaded{Config: roundconfig.Builtin(), GitRoot: root}, func(context.Context, roundconfig.Config, []roundconfig.WorkCategory, string) profileProofResult {
		return profileProofResult{}
	})
	updateCommandDependenciesForTest(t, func(dependencies *commandDependencies) {
		dependencies.doctor.resolveExternal = func(string) ([]string, bool, error) { return []string{"testing-boss"}, true, nil }
		dependencies.doctor.checkSkills = skills.CheckRepositoryWithExternal
		dependencies.doctor.trailingSkills = defaultDoctorDependencies().trailingSkills
	})
	before := snapshotDoctorPath(t, root)
	var stdout, stderr bytes.Buffer
	code := runCLI(t, []string{"doctor"}, &stdout, &stderr)
	want := fmt.Sprintf("skills: warn (%d required: %d Roundfix-owned, 1 external; DR-SKILL-TRAILS-SNAPSHOT: trails the Setup Snapshot: testing-boss; next: roundfix baseline update)", len(skills.Names())+1, len(skills.Names()))
	if code != exitOK || stderr.Len() != 0 || !strings.Contains(stdout.String(), want) {
		t.Fatalf("exit=%d stdout=%s stderr=%s; want %s", code, stdout.String(), stderr.String(), want)
	}
	if !reflect.DeepEqual(before, snapshotDoctorPath(t, root)) {
		t.Fatal("Doctor changed repository bytes")
	}
}

func TestDoctorSnapshotComparisonPreservesReadinessAndFiltersSkills(t *testing.T) {
	t.Parallel()
	for _, test := range []struct {
		name           string
		readiness      skills.RepositoryReadiness
		comparisonErr  error
		trailing       []string
		wantStatus     string
		wantDetail     string
		wantComparable []string
	}{
		{name: "matching", wantStatus: "ok", wantDetail: "skills: ok (0 required: 0 Roundfix-owned, 0 external)", wantComparable: []string{"testing-boss", "golang-cli", "golang-context"}},
		{name: "mixed failure", readiness: skills.RepositoryReadiness{MissingExternal: []string{"golang-cli"}, OutdatedExternal: []string{"golang-context"}}, trailing: []string{"testing-boss"}, wantStatus: "failed", wantDetail: "DR-SKILL-TRAILS-SNAPSHOT: trails the Setup Snapshot: testing-boss", wantComparable: []string{"testing-boss"}},
		{name: "unknown profile", comparisonErr: errors.New("profile is not embedded"), wantStatus: "ok", wantDetail: "snapshot comparison unavailable", wantComparable: []string{"testing-boss", "golang-cli", "golang-context"}},
	} {
		t.Run(test.name, func(t *testing.T) {
			checker := newDoctorFakeHealthChecker(CheckResult{Name: HealthCheckNode, Status: CheckStatusOK}, CheckResult{Name: HealthCheckACPX, Status: CheckStatusOK}, CheckResult{Name: HealthCheckCodex, Status: CheckStatusOK})
			withDoctorFakeLoadedAndReadiness(t, checker, roundconfig.Loaded{Config: roundconfig.Builtin(), GitRoot: "/repo"}, func(context.Context, roundconfig.Config, []roundconfig.WorkCategory, string) profileProofResult {
				return profileProofResult{}
			})
			updateCommandDependenciesForTest(t, func(dependencies *commandDependencies) {
				dependencies.doctor.resolveExternal = func(string) ([]string, bool, error) {
					return []string{"testing-boss", "golang-cli", "golang-context"}, true, nil
				}
				dependencies.doctor.checkSkills = func(context.Context, string, []string) (skills.RepositoryReadiness, error) {
					return test.readiness, nil
				}
				dependencies.doctor.trailingSkills = func(_ string, comparable []string) ([]string, error) {
					if !reflect.DeepEqual(comparable, test.wantComparable) {
						t.Fatalf("comparable=%v, want %v", comparable, test.wantComparable)
					}
					return test.trailing, test.comparisonErr
				}
			})
			var stdout, stderr bytes.Buffer
			code := runCLI(t, []string{"doctor"}, &stdout, &stderr)
			wantCode := exitOK
			if test.wantStatus == "failed" {
				wantCode = exitRunFailed
			}
			if code != wantCode || !strings.Contains(stdout.String(), "skills: "+test.wantStatus) || !strings.Contains(stdout.String(), test.wantDetail) {
				t.Fatalf("exit=%d stdout=%s stderr=%s", code, stdout.String(), stderr.String())
			}
		})
	}
}

func TestBaselineUpdatePreviewListsATrailingSkill(t *testing.T) {
	t.Parallel()
	root := newBaselineUpdateRepository(t)
	source, err := filepath.Abs("../..")
	if err != nil {
		t.Fatal(err)
	}
	for _, relative := range []string{".agents/skills", "skills-lock.json"} {
		copyThisRepositorySkillSetPath(t, filepath.Join(source, relative), filepath.Join(root, relative))
	}
	// Make the installed tree intentionally different while keeping its lock intact.
	path := filepath.Join(root, ".agents", "skills", "testing-boss", "SKILL.md")
	if err := os.WriteFile(path, append(mustReadBytes(t, path), []byte("\nOld upstream installation fixture.\n")...), 0o644); err != nil {
		t.Fatal(err)
	}
	lockPath := filepath.Join(root, "skills-lock.json")
	var lock map[string]any
	if err := json.Unmarshal(mustReadBytes(t, lockPath), &lock); err != nil {
		t.Fatal(err)
	}
	hash, err := skills.SkillFolderHash(t.Context(), filepath.Dir(path))
	if err != nil {
		t.Fatal(err)
	}
	lock["skills"].(map[string]any)["testing-boss"].(map[string]any)["computedHash"] = hash
	data, err := json.Marshal(lock)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(lockPath, data, 0o644); err != nil {
		t.Fatal(err)
	}
	before := baselinePlanTestTree(t, root)
	result, out, stderr, code := runBaselineUpdateTestCommand(t, t.Context(), "baseline", "update", "--repo", root, "--format=json")
	if code != exitUnverified || result.State != "plan_ready" || stderr != "" {
		t.Fatalf("exit=%d stdout=%s stderr=%s", code, out, stderr)
	}
	found := false
	for _, drift := range result.Skills.Drifted {
		if drift.Skill == "testing-boss" && strings.Contains(drift.Reason, "trails its Setup Snapshot") {
			found = true
		}
	}
	if !found {
		t.Fatalf("trailing skill absent from preview: %s", out)
	}
	_, out, stderr, code = runBaselineUpdateTestCommand(t, t.Context(), "baseline", "update", "--repo", root, "--format=text")
	if code != exitUnverified || stderr != "" || !strings.Contains(out, "- drifted testing-boss: trails its Setup Snapshot") {
		t.Fatalf("exit=%d stdout=%s stderr=%s", code, out, stderr)
	}
	if baselinePlanTestTree(t, root) != before {
		t.Fatal("preview changed repository bytes")
	}
	_, out, stderr, code = runBaselineUpdateTestCommand(t, t.Context(), "baseline", "update", "--repo", root, "--no-skills", "--format=json")
	if code != exitOK || stderr != "" || strings.Contains(out, "trails its Setup Snapshot") {
		t.Fatalf("--no-skills exit=%d stdout=%s stderr=%s", code, out, stderr)
	}
}

func TestBaselineUpdateRestoresTrailingSkillsThroughExistingRequest(t *testing.T) {
	t.Parallel()
	var requests []baseline.SkillsRestoreRequest
	deps := baselineUpdateSkillsDependencies{
		resolveProjectRoot: func(context.Context, string) (string, error) { return "/repo", nil },
		install: func(context.Context, skills.InstallRequest) (skills.InstallResult, error) {
			return skills.InstallResult{}, nil
		},
		ownedNames:      func() []string { return nil },
		resolveExternal: func(string) ([]string, bool, error) { return []string{"testing-boss"}, true, nil },
		checkRepository: func(context.Context, string, []string) (skills.RepositoryReadiness, error) {
			return skills.RepositoryReadiness{}, nil
		},
		trailingSkills: func(root, profile string, required []string) ([]string, error) {
			if root != "/repo" || profile != "go-cli-tui" || !reflect.DeepEqual(required, []string{"testing-boss"}) {
				t.Fatalf("comparison input=%s %s %v", root, profile, required)
			}
			return []string{"testing-boss"}, nil
		},
		restore: func(_ context.Context, request baseline.SkillsRestoreRequest) (baseline.SkillsRestorePayload, error) {
			requests = append(requests, request)
			payload := baseline.SkillsRestorePayload{PlanDigest: pointerToString(strings.Repeat("a", 64))}
			if request.Confirmation == "" {
				payload.Finding = &baseline.RestoreFinding{Code: "plan.confirmation.required"}
				return payload, &baseline.SkillsRestoreError{Category: baseline.SkillsRestoreAction, Finding: baseline.RestoreFinding{Code: "plan.confirmation.required"}}
			}
			payload.Applied = true
			return payload, nil
		},
	}
	result, err := runBaselineUpdateSkillsStageWith(t.Context(), baselineUpdateSkillsRequest{Repository: "/repo", ProfileID: "go-cli-tui", SourceDir: "/offline"}, deps)
	if err != nil || !reflect.DeepEqual(result.Restored, []string{"testing-boss"}) || len(requests) != 2 {
		t.Fatalf("result=%+v requests=%+v err=%v", result, requests, err)
	}
	for _, request := range requests {
		if request.ProfileID != "go-cli-tui" || request.SourceDir != "/offline" || !reflect.DeepEqual(request.Skills, []string{"testing-boss"}) {
			t.Fatalf("restore request=%+v", request)
		}
	}
	if requests[1].Confirmation != strings.Repeat("a", 64) {
		t.Fatal("restore did not confirm its preview")
	}
}
