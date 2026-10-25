// Suite: Missing Required Skills in Baseline preview and Doctor.
// Invariant: required external skills are named without changing repository bytes.
// Boundary IN: temporary adopted repositories, real skill inspection, command output and exit codes.
// Boundary OUT: network, host skills and unrelated Doctor health checks.

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

	roundconfig "roundfix/internal/config"
	"roundfix/skills"
)

func installRequiredExternalSkills(t *testing.T, repository string) []string {
	t.Helper()
	required, ok, err := resolveExternalSkillRequirement(repository)
	if err != nil || !ok || len(required) == 0 {
		t.Fatalf("required external skills=%v manifest=%t err=%v", required, ok, err)
	}
	lock := map[string]any{"version": 1, "skills": map[string]any{}}
	lockPath := filepath.Join(repository, "skills-lock.json")
	data, err := os.ReadFile(lockPath)
	if err == nil {
		if err := json.Unmarshal(data, &lock); err != nil {
			t.Fatal(err)
		}
	} else if !os.IsNotExist(err) {
		t.Fatal(err)
	}
	entries := lock["skills"].(map[string]any)
	for _, name := range required {
		folder := filepath.Join(repository, ".agents", "skills", name)
		writeBaselinePlanTestFile(t, repository, filepath.Join(".agents", "skills", name, "SKILL.md"), "---\nname: "+name+"\ndescription: fixture\n---\n# Fixture\n")
		hash, err := skills.SkillFolderHash(t.Context(), folder)
		if err != nil {
			t.Fatal(err)
		}
		entries[name] = map[string]any{"computedHash": hash}
	}
	data, err = json.Marshal(lock)
	if err != nil {
		t.Fatal(err)
	}
	mustWrite(t, lockPath, string(data))
	commitBaselinePlanTestRepository(t, repository)
	return required
}

func removeInstalledSkill(t *testing.T, repository, name string, lockEntry bool) {
	t.Helper()
	if err := os.RemoveAll(filepath.Join(repository, ".agents", "skills", name)); err != nil {
		t.Fatal(err)
	}
	if lockEntry {
		path := filepath.Join(repository, "skills-lock.json")
		var lock map[string]any
		if err := json.Unmarshal(mustReadBytes(t, path), &lock); err != nil {
			t.Fatal(err)
		}
		delete(lock["skills"].(map[string]any), name)
		data, err := json.Marshal(lock)
		if err != nil {
			t.Fatal(err)
		}
		mustWrite(t, path, string(data))
	}
}

func TestBaselineUpdatePreviewNamesAMissingRequiredSkill(t *testing.T) {
	t.Parallel()
	for _, shape := range []string{"folder and lock entry removed", "folder removed with a stale lock entry", "folder no longer matches its lock entry", "every required skill installed"} {
		t.Run(shape, func(t *testing.T) {
			repository := newBaselineUpdateRepository(t)
			name := installRequiredExternalSkills(t, repository)[0]
			reason := "required by the Baseline and not installed; restore with roundfix baseline skills restore --profile go-cli-tui --skill " + name
			switch shape {
			case "folder and lock entry removed":
				removeInstalledSkill(t, repository, name, true)
			case "folder removed with a stale lock entry":
				removeInstalledSkill(t, repository, name, false)
			case "folder no longer matches its lock entry":
				mustWrite(t, filepath.Join(repository, ".agents", "skills", name, "SKILL.md"), "changed fixture\n")
				reason = "required by the Baseline and no longer matches its skills-lock.json entry; restore with roundfix baseline skills restore --profile go-cli-tui --skill " + name
			}
			if shape != "every required skill installed" {
				commitBaselinePlanTestRepository(t, repository)
			}
			before := baselinePlanTestTree(t, repository)
			result, out, stderr, code := runBaselineUpdateTestCommand(t, t.Context(), "baseline", "update", "--repo", repository, "--no-history", "--format=json")
			if shape == "every required skill installed" {
				if code != exitOK || stderr != "" || result.State != "current" || strings.Contains(out, `"missing"`) {
					t.Fatalf("complete preview exit=%d stdout=%s stderr=%s", code, out, stderr)
				}
			} else {
				want := []baselineUpdateSkillDrift{{Skill: name, Reason: reason}}
				message := "guidance matches the current Baseline catalog; 1 required skill(s) are missing or no longer match skills-lock.json: " + name
				action := "rerun with --yes to restore them from the Setup Snapshot, or restore one with roundfix baseline skills restore --profile go-cli-tui --skill " + name
				// Decode the optional list independently so the regression also builds before implementation.
				var payload struct {
					Skills struct{ Missing []baselineUpdateSkillDrift }
				}
				if err := json.Unmarshal([]byte(out), &payload); err != nil {
					t.Fatal(err)
				}
				if code != exitUnverified || stderr != "" || result.OK || result.State != "plan_ready" || result.Category != "approval" || result.Skills.Status != "missing" || !reflect.DeepEqual(payload.Skills.Missing, want) || result.Message != message || result.NextAction != action {
					t.Fatalf("missing preview exit=%d result=%+v stdout=%s stderr=%s", code, result, out, stderr)
				}
				_, text, textErr, textCode := runBaselineUpdateTestCommand(t, t.Context(), "baseline", "update", "--repo", repository, "--no-history", "--format=text")
				if textCode != exitUnverified || textErr != "" || !strings.Contains(text, "Skills: missing\n") || !strings.Contains(text, "Skills missing: 1\n- missing "+name+": "+reason+"\n") {
					t.Fatalf("text exit=%d stdout=%s stderr=%s", textCode, text, textErr)
				}
				skipped, skipOut, skipErr, skipCode := runBaselineUpdateTestCommand(t, t.Context(), "baseline", "update", "--repo", repository, "--no-history", "--no-skills", "--format=json")
				if skipCode != exitOK || skipErr != "" || skipped.State != "current" || strings.Contains(skipOut, `"missing"`) {
					t.Fatalf("skipped exit=%d stdout=%s stderr=%s", skipCode, skipOut, skipErr)
				}
			}
			if before != baselinePlanTestTree(t, repository) {
				t.Fatal("preview changed repository bytes")
			}
		})
	}
}

func TestDoctorNamesAMissingRequiredSkill(t *testing.T) {
	t.Parallel()
	for _, lockEntry := range []bool{true, false} {
		shape := "folder removed with a stale lock entry"
		if lockEntry {
			shape = "folder and lock entry removed"
		}
		t.Run(shape, func(t *testing.T) {
			repository := newBaselineUpdateRepository(t)
			installPreviewOwnedSkills(t, repository)
			name := installRequiredExternalSkills(t, repository)[0]
			removeInstalledSkill(t, repository, name, lockEntry)
			checker := newDoctorFakeHealthChecker(CheckResult{Name: HealthCheckNode, Status: CheckStatusOK}, CheckResult{Name: HealthCheckACPX, Status: CheckStatusOK}, CheckResult{Name: HealthCheckCodex, Status: CheckStatusOK})
			withDoctorFakeLoadedAndReadiness(t, checker, roundconfig.Loaded{Config: roundconfig.Builtin(), GitRoot: repository}, func(context.Context, roundconfig.Config, []roundconfig.WorkCategory, string) profileProofResult {
				return profileProofResult{}
			})
			updateCommandDependenciesForTest(t, func(dependencies *commandDependencies) {
				dependencies.doctor.resolveExternal = resolveExternalSkillRequirement
				dependencies.doctor.checkSkills = skills.CheckRepositoryWithExternal
			})
			var stdout, stderr bytes.Buffer
			code := runCLI(t, []string{"doctor"}, &stdout, &stderr)
			var skillLine string
			for _, line := range strings.Split(stdout.String(), "\n") {
				if strings.Contains(line, "skills: failed (") {
					skillLine = line
				}
			}
			if code != exitRunFailed || !strings.Contains(skillLine, name) || !strings.Contains(skillLine, fmt.Sprintf("bunx skills add marcioaltoe/skills@%s", name)) {
				t.Fatalf("doctor exit=%d stdout=%s stderr=%s", code, stdout.String(), stderr.String())
			}
		})
	}
}

func TestMissingRequiredSkillsCollectsNamedEvidence(t *testing.T) {
	t.Parallel()
	for _, test := range []struct {
		name      string
		readiness skills.RepositoryReadiness
		err       error
		want      []baselineUpdateSkillDrift
	}{
		{name: "unnamed inspection error", err: errors.New("unreadable lock")},
		{name: "sorted deduplicated evidence", readiness: skills.RepositoryReadiness{MissingExternal: []string{"zeta"}, OutdatedExternal: []string{"zeta", "beta", "alpha"}}, err: fmt.Errorf("inspect: %w", &skills.RepositoryReadinessError{MissingExternal: []string{"alpha"}}), want: []baselineUpdateSkillDrift{
			{Skill: "alpha", Reason: "required by the Baseline and not installed; restore with roundfix baseline skills restore --profile fixture --skill alpha"},
			{Skill: "beta", Reason: "required by the Baseline and no longer matches its skills-lock.json entry; restore with roundfix baseline skills restore --profile fixture --skill beta"},
			{Skill: "zeta", Reason: "required by the Baseline and not installed; restore with roundfix baseline skills restore --profile fixture --skill zeta"},
		}},
	} {
		t.Run(test.name, func(t *testing.T) {
			if got := missingRequiredSkills("fixture", test.readiness, test.err); !reflect.DeepEqual(got, test.want) {
				t.Fatalf("missing=%+v want=%+v", got, test.want)
			}
		})
	}
}

func TestRequiredSkillsManifestProfileReadsOnlyProfile(t *testing.T) {
	t.Parallel()
	root := t.TempDir()
	if _, err := requiredSkillsManifestProfile(root); err == nil {
		t.Fatal("missing manifest accepted")
	}
	writeBaselinePlanTestFile(t, root, doctorSetupManifestPath, "{")
	if _, err := requiredSkillsManifestProfile(root); err == nil {
		t.Fatal("malformed manifest accepted")
	}
	mustWrite(t, filepath.Join(root, filepath.FromSlash(doctorSetupManifestPath)), `{"profile":"fixture","ignored":true}`)
	if profile, err := requiredSkillsManifestProfile(root); err != nil || profile != "fixture" {
		t.Fatalf("profile=%q err=%v", profile, err)
	}
}
