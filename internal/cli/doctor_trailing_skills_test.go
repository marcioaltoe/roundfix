package cli

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"roundfix/internal/baseline"
	roundconfig "roundfix/internal/config"
	"roundfix/skills"
)

func TestDoctorNamesASetupSnapshotDifferenceWithoutWarning(t *testing.T) {
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
		dependencies.doctor.snapshotDifferences = defaultDoctorDependencies().snapshotDifferences
	})
	before := snapshotDoctorPath(t, root)
	var stdout, stderr bytes.Buffer
	code := runCLI(t, []string{"doctor"}, &stdout, &stderr)
	want := fmt.Sprintf("skills: ok (%d required: %d Roundfix-owned, 1 external; DR-SKILL-DIFFERS-FROM-SNAPSHOT: differs from the Setup Snapshot, newer or older unknown: testing-boss)", len(skills.Names())+1, len(skills.Names()))
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
		differing      []string
		wantStatus     string
		wantDetail     string
		wantComparable []string
	}{
		{name: "matching", wantStatus: "ok", wantDetail: "skills: ok (0 required: 0 Roundfix-owned, 0 external)", wantComparable: []string{"testing-boss", "golang-cli", "golang-context"}},
		{name: "mixed failure", readiness: skills.RepositoryReadiness{MissingExternal: []string{"golang-cli"}, OutdatedExternal: []string{"golang-context"}}, differing: []string{"testing-boss"}, wantStatus: "failed", wantDetail: "DR-SKILL-DIFFERS-FROM-SNAPSHOT: differs from the Setup Snapshot, newer or older unknown: testing-boss", wantComparable: []string{"testing-boss"}},
		{name: "unknown profile", comparisonErr: errors.New("profile is not embedded"), wantStatus: "ok", wantDetail: "snapshot comparison unavailable", wantComparable: []string{"testing-boss", "golang-cli", "golang-context"}},
	} {
		t.Run(test.name, func(t *testing.T) {
			checker := newDoctorFakeHealthChecker(CheckResult{Name: HealthCheckNode, Status: CheckStatusOK}, CheckResult{Name: HealthCheckACPX, Status: CheckStatusOK}, CheckResult{Name: HealthCheckCodex, Status: CheckStatusOK})
			withDoctorFakeLoadedAndReadiness(t, checker, roundconfig.Loaded{Config: roundconfig.Builtin(), GitRoot: "/repo"}, func(context.Context, roundconfig.Config, []roundconfig.WorkCategory, string) profileProofResult {
				return profileProofResult{}
			})
			var doctorDeps doctorDependencies
			updateCommandDependenciesForTest(t, func(dependencies *commandDependencies) {
				dependencies.doctor.resolveExternal = func(string) ([]string, bool, error) {
					return []string{"testing-boss", "golang-cli", "golang-context"}, true, nil
				}
				dependencies.doctor.checkSkills = func(context.Context, string, []string) (skills.RepositoryReadiness, error) {
					return test.readiness, nil
				}
				dependencies.doctor.snapshotDifferences = func(_ string, comparable []string) ([]string, error) {
					if !reflect.DeepEqual(comparable, test.wantComparable) {
						t.Fatalf("comparable=%v, want %v", comparable, test.wantComparable)
					}
					return test.differing, test.comparisonErr
				}
				doctorDeps = dependencies.doctor
			})
			readinessResult := doctorSkillReadinessResult(test.readiness, nil)
			comparisonResult := repositorySkillsCheck(t.Context(), doctorDeps, "/repo")
			if comparisonResult.Status != readinessResult.Status || comparisonResult.NextAction != readinessResult.NextAction {
				t.Fatalf("snapshot comparison changed readiness status or action: got %+v, readiness %+v", comparisonResult, readinessResult)
			}
			var stdout, stderr bytes.Buffer
			code := runCLI(t, []string{"doctor"}, &stdout, &stderr)
			wantCode := exitOK
			if test.wantStatus == "failed" {
				wantCode = exitRunFailed
			}
			if strings.Contains(stdout.String(), "roundfix baseline update") {
				t.Fatal("snapshot comparison added a next action")
			}
			if code != wantCode || !strings.Contains(stdout.String(), "skills: "+test.wantStatus) || !strings.Contains(stdout.String(), test.wantDetail) {
				t.Fatalf("exit=%d stdout=%s stderr=%s", code, stdout.String(), stderr.String())
			}
		})
	}
}

func TestBaselineUpdateKeepsASetupSnapshotDifference(t *testing.T) {
	t.Parallel()
	for _, test := range []struct {
		name         string
		readiness    skills.RepositoryReadiness
		wantRestored []string
	}{
		{name: "difference alone", wantRestored: []string{}},
		{name: "missing", readiness: skills.RepositoryReadiness{MissingExternal: []string{"golang-cli"}}, wantRestored: []string{"golang-cli"}},
		{name: "lock mismatch", readiness: skills.RepositoryReadiness{OutdatedExternal: []string{"golang-cli"}}, wantRestored: []string{"golang-cli"}},
	} {
		t.Run(test.name, func(t *testing.T) {
			var requests []baseline.SkillsRestoreRequest
			deps := baselineUpdateSkillsDependencies{
				recordLock:         func(context.Context, string, []string) ([]string, error) { return nil, nil },
				resolveProjectRoot: func(context.Context, string) (string, error) { return "/repo", nil },
				install: func(context.Context, skills.InstallRequest) (skills.InstallResult, error) {
					return skills.InstallResult{}, nil
				},
				ownedNames:      func() []string { return nil },
				resolveExternal: func(string) ([]string, bool, error) { return []string{"testing-boss", "golang-cli"}, true, nil },
				checkRepository: func(context.Context, string, []string) (skills.RepositoryReadiness, error) {
					return test.readiness, nil
				},
				snapshotDifferences: func(root, profile string, required []string) ([]string, error) {
					wantComparable := []string{"testing-boss"}
					if len(test.wantRestored) == 0 {
						wantComparable = append(wantComparable, "golang-cli")
					}
					if root != "/repo" || profile != "go-cli-tui" || !reflect.DeepEqual(required, wantComparable) {
						t.Fatalf("comparison input=%s %s %v", root, profile, required)
					}
					return []string{"testing-boss"}, nil
				},
				restore: func(_ context.Context, request baseline.SkillsRestoreRequest) (baseline.SkillsRestorePayload, error) {
					for _, name := range request.Skills {
						if name == "testing-boss" {
							t.Fatal("restore called for the differing skill")
						}
					}
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
			wantDifferences := []baselineUpdateSkillDrift{{Skill: "testing-boss", Reason: "differs from its Setup Snapshot, newer or older unknown; kept; restore deliberately with roundfix baseline skills restore --profile go-cli-tui --skill testing-boss"}}
			if err != nil || result.Status != baselineUpdateSkillsVerified || !reflect.DeepEqual(result.Restored, test.wantRestored) || !reflect.DeepEqual(result.Differing, wantDifferences) || len(result.Drifted) != 0 || len(requests) != 2*len(test.wantRestored) {
				t.Fatalf("result=%+v requests=%+v err=%v", result, requests, err)
			}
			for _, request := range requests {
				if request.ProfileID != "go-cli-tui" || request.SourceDir != "/offline" || !reflect.DeepEqual(request.Skills, []string{"golang-cli"}) {
					t.Fatalf("restore request=%+v", request)
				}
			}
			if len(requests) != 0 && requests[1].Confirmation != strings.Repeat("a", 64) {
				t.Fatal("restore did not confirm its preview")
			}
		})
	}
}
