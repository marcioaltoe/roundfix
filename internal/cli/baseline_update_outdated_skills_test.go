package cli

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"testing"

	roundconfig "roundfix/internal/config"
	"roundfix/skills"
)

func installPreviewOwnedSkills(t *testing.T, repository string) {
	t.Helper()
	if _, err := skills.Install(t.Context(), skills.InstallRequest{Target: "project", ProjectDir: repository}); err != nil {
		t.Fatal(err)
	}
}

func lowerPreviewOwnedSkill(t *testing.T, repository, name string) (found, required string) {
	t.Helper()
	readiness, err := skills.CheckRepositoryWithExternal(t.Context(), repository, nil)
	if err != nil {
		t.Fatal(err)
	}
	for _, owned := range readiness.Owned {
		if owned.Skill == name {
			required = owned.Found
		}
	}
	parts := strings.Split(required, ".")
	if len(parts) != 3 {
		t.Fatalf("invalid embedded version %q", required)
	}
	patch, err := strconv.Atoi(parts[2])
	if err != nil || patch < 1 {
		t.Fatalf("cannot lower patch of %q", required)
	}
	found = fmt.Sprintf("%s.%s.%d", parts[0], parts[1], patch-1)
	path := filepath.Join(repository, ".agents", "skills", name, "SKILL.md")
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	old := []byte("\nversion: " + required + "\n")
	if bytes.Count(data, old) != 1 {
		t.Fatal("expected one top-level version")
	}
	if err := os.WriteFile(path, bytes.Replace(data, old, []byte("\nversion: "+found+"\n"), 1), 0o644); err != nil {
		t.Fatal(err)
	}
	return found, required
}

func TestBaselineUpdatePreviewReportsAnOlderOwnedSkill(t *testing.T) {
	t.Parallel()
	for _, changed := range []bool{false, true} {
		t.Run(fmt.Sprintf("guidance_changed=%t", changed), func(t *testing.T) {
			repository := newBaselineUpdateRepository(t)
			if changed {
				repository = staleBaselineUpdateRepository(t)
			}
			installPreviewOwnedSkills(t, repository)
			before, _, _, _ := runBaselineUpdateTestCommand(t, t.Context(), "baseline", "update", "--repo", repository, "--format=json")
			found, required := lowerPreviewOwnedSkill(t, repository, "qa-gate")
			tree := baselinePlanTestTree(t, repository)
			result, out, stderr, code := runBaselineUpdateTestCommand(t, t.Context(), "baseline", "update", "--repo", repository, "--format=json")
			if code != exitUnverified || stderr != "" || result.State != "plan_ready" || result.Category != "approval" || result.Skills.Status != "outdated" {
				t.Fatalf("exit=%d result=%+v stdout=%s stderr=%s", code, result, out, stderr)
			}
			want := []baselineUpdateSkillOutdated{{Skill: "qa-gate", Found: found, Required: required}}
			if !reflect.DeepEqual(result.Skills.Outdated, want) {
				t.Fatalf("outdated=%+v, want %+v", result.Skills.Outdated, want)
			}
			if changed {
				if result.Message != before.Message || result.NextAction != before.NextAction {
					t.Fatal("skill check changed the guidance plan message or action")
				}
			} else {
				if result.Message != "guidance matches the current Baseline catalog; 1 Roundfix-owned skill(s) are older than the ones this binary carries" || result.NextAction != "rerun with --yes to refresh the Repository Skill Set, or run roundfix skills install --target project" {
					t.Fatalf("message/action = %q / %q", result.Message, result.NextAction)
				}
			}
			_, text, textErr, textCode := runBaselineUpdateTestCommand(t, t.Context(), "baseline", "update", "--repo", repository, "--format=text")
			wantLine := fmt.Sprintf("Skills drifted: 0\nSkills outdated: 1\n- outdated qa-gate: found %s, requires %s\n", found, required)
			if textCode != exitUnverified || textErr != "" || !strings.Contains(text, wantLine) {
				t.Fatalf("text exit=%d stdout=%s stderr=%s", textCode, text, textErr)
			}
			if baselinePlanTestTree(t, repository) != tree {
				t.Fatal("preview changed repository bytes")
			}
		})
	}
}

func TestBaselineUpdatePreviewStaysCurrentWhenOwnedSkillsMatch(t *testing.T) {
	t.Parallel()
	repository := newBaselineUpdateRepository(t)
	_, absent, absentErr, absentCode := runBaselineUpdateTestCommand(t, t.Context(), "baseline", "update", "--repo", repository, "--format=json")
	installPreviewOwnedSkills(t, repository)
	result, out, stderr, code := runBaselineUpdateTestCommand(t, t.Context(), "baseline", "update", "--repo", repository, "--format=json")
	if code != exitOK || result.State != "current" || stderr != "" || strings.Contains(out, "\"outdated\"") {
		t.Fatalf("exit=%d result=%+v stdout=%s stderr=%s", code, result, out, stderr)
	}
	if out != absent || code != absentCode || stderr != absentErr {
		t.Fatal("matching and absent skills produce different preview bytes")
	}
	// The shared skills directory must not be a symlink. Keep all profile
	// capabilities reachable so planning still succeeds before the check.
	skillsPath := filepath.Join(repository, ".agents", "skills")
	if err := os.Rename(skillsPath, filepath.Join(repository, "installed-skills")); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink("../installed-skills", skillsPath); err != nil {
		t.Fatal(err)
	}
	if _, err := skills.CheckRepositoryWithExternal(t.Context(), repository, nil); err == nil {
		t.Fatal("fixture did not produce a repository read error")
	}
	_, errorOut, errorErr, errorCode := runBaselineUpdateTestCommand(t, t.Context(), "baseline", "update", "--repo", repository, "--format=json")
	if errorOut != out || errorErr != stderr || errorCode != code {
		t.Fatalf("read error changed preview output: before code=%d stdout=%s stderr=%s; after code=%d stdout=%s stderr=%s", code, out, stderr, errorCode, errorOut, errorErr)
	}
}

func TestBaselineUpdatePreviewWithNoSkillsSkipsTheOwnedSkillCheck(t *testing.T) {
	t.Parallel()
	repository := newBaselineUpdateRepository(t)
	installPreviewOwnedSkills(t, repository)
	_, before, beforeErr, beforeCode := runBaselineUpdateTestCommand(t, t.Context(), "baseline", "update", "--repo", repository, "--no-skills", "--format=json")
	lowerPreviewOwnedSkill(t, repository, "qa-gate")
	result, out, stderr, code := runBaselineUpdateTestCommand(t, t.Context(), "baseline", "update", "--repo", repository, "--no-skills", "--format=json")
	if code != exitOK || result.State != "current" || out != before || stderr != beforeErr || code != beforeCode || strings.Contains(out, "\"outdated\"") {
		t.Fatalf("--no-skills changed: exit=%d stdout=%s stderr=%s", code, out, stderr)
	}
}

func TestDoctorFailsForAnOwnedSkillOlderThanTheBundle(t *testing.T) {
	t.Parallel()
	repository := newBaselineUpdateRepository(t)
	installPreviewOwnedSkills(t, repository)
	found, required := lowerPreviewOwnedSkill(t, repository, "qa-gate")
	checker := newDoctorFakeHealthChecker(CheckResult{Name: HealthCheckNode, Status: CheckStatusOK}, CheckResult{Name: HealthCheckACPX, Status: CheckStatusOK}, CheckResult{Name: HealthCheckCodex, Status: CheckStatusOK})
	withDoctorFakeLoadedAndReadiness(t, checker, roundconfig.Loaded{Config: roundconfig.Builtin(), GitRoot: repository}, func(context.Context, roundconfig.Config, []roundconfig.WorkCategory, string) profileProofResult {
		return profileProofResult{}
	})
	updateCommandDependenciesForTest(t, func(dependencies *commandDependencies) {
		dependencies.doctor.resolveExternal = func(string) ([]string, bool, error) { return nil, true, nil }
		dependencies.doctor.checkSkills = skills.CheckRepositoryWithExternal
	})
	var stdout, stderr bytes.Buffer
	code := runCLI(t, []string{"doctor"}, &stdout, &stderr)
	want := fmt.Sprintf("skills: failed (below minimum: skill \"qa-gate\" requires %s, found %s; next: roundfix skills install --target project)", required, found)
	if code != exitRunFailed || !strings.Contains(stdout.String(), want) {
		t.Fatalf("doctor exit=%d stdout=%s stderr=%s; want %s", code, stdout.String(), stderr.String(), want)
	}
}
