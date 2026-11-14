package cli

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"roundfix/internal/app"
	roundconfig "roundfix/internal/config"
	roundskills "roundfix/skills"
)

func TestSkillsCheckWarnsOfStaleUserSkillCopiesAndKeepsItsExit(t *testing.T) {
	t.Parallel()
	home, copies, descriptions := staleUserSkillsFixture(t)
	setCommandHomeDirForTest(t, home)
	setCommandWorkDirForTest(t, t.TempDir())
	var stdout, stderr bytes.Buffer
	code := runCLI(t, []string{"skills", "check"}, &stdout, &stderr)
	wantOut := "Roundfix skill check passed: " + strings.Join(roundskills.Names(), ", ") + "\n"
	wantErr := "roundfix: warning: user-level Roundfix skill copies that an Agent harness may load in place of the project copy:\n  " + strings.Join(descriptions, "\n  ") + "\n  remove: roundfix skills remove --target claude --apply\n  remove: roundfix skills remove --target opencode --apply\n  refresh: roundfix skills install --target claude\n  refresh: roundfix skills install --target opencode\n"
	if code != 0 || stdout.String() != wantOut || stderr.String() != wantErr {
		t.Fatalf("exit=%d\nstdout=%q\nstderr=%q\nwant stdout=%q\nwant stderr=%q", code, stdout.String(), stderr.String(), wantOut, wantErr)
	}
	t.Run("unversioned still warns", func(t *testing.T) {
		setCommandHomeDirForTest(t, home)
		setCommandWorkDirForTest(t, t.TempDir())
		updateCommandDependenciesForTest(t, func(deps *commandDependencies) {
			deps.checkOwnedSkills = func() roundskills.BundleReadiness {
				return roundskills.BundleReadiness{Owned: []roundskills.OwnedSkillReadiness{{Skill: "roundfix", State: roundskills.ReadinessUnversioned}}}
			}
		})
		var stdout, stderr bytes.Buffer
		code := runCLI(t, []string{"skills", "check"}, &stdout, &stderr)
		if code != 0 || stdout.String() != "Roundfix skill check unversioned: roundfix\n" || stderr.String() != wantErr {
			t.Fatalf("exit=%d stdout=%q stderr=%q", code, stdout.String(), stderr.String())
		}
	})
	t.Run("failed check has no advisory", func(t *testing.T) {
		setCommandHomeDirForTest(t, home)
		setCommandWorkDirForTest(t, t.TempDir())
		updateCommandDependenciesForTest(t, func(deps *commandDependencies) {
			deps.checkOwnedSkills = func() roundskills.BundleReadiness {
				return roundskills.BundleReadiness{Diagnostics: []roundskills.Diagnostic{{Path: "bundle", Message: "fixture failure"}}}
			}
		})
		var stdout, stderr bytes.Buffer
		code := runCLI(t, []string{"skills", "check"}, &stdout, &stderr)
		if code != exitRunFailed || stdout.Len() != 0 || stderr.String() != "roundfix: Roundfix skill check failed:\n  bundle: fixture failure\n" {
			t.Fatalf("exit=%d stdout=%q stderr=%q", code, stdout.String(), stderr.String())
		}
	})
	// Bring both copies to the embedded versions: the next check is silent.
	for _, copy := range copies {
		writeStaleUserSkillFixture(t, copy.Path, copy.Skill, copy.Embedded)
	}
	stdout.Reset()
	stderr.Reset()
	code = runCLI(t, []string{"skills", "check"}, &stdout, &stderr)
	if code != 0 || stdout.String() != wantOut || stderr.Len() != 0 {
		t.Fatalf("current copies: exit=%d stdout=%q stderr=%q", code, stdout.String(), stderr.String())
	}
}

func TestDoctorReportsStaleUserSkillCopiesWithoutFailing(t *testing.T) {
	t.Parallel()
	home, _, descriptions := staleUserSkillsFixture(t)
	setCommandHomeDirForTest(t, home)
	setCommandWorkDirForTest(t, t.TempDir())
	checker := newDoctorFakeHealthChecker(
		CheckResult{Name: HealthCheckNode, Status: CheckStatusOK},
		CheckResult{Name: HealthCheckACPX, Status: CheckStatusOK},
		CheckResult{Name: HealthCheckCodex, Status: CheckStatusOK},
	)
	withDoctorFakeLoadedAndReadiness(t, checker, roundconfig.Loaded{Config: roundconfig.Builtin(), GitRoot: t.TempDir(), HomeDir: home}, func(context.Context, roundconfig.Config, []roundconfig.WorkCategory, string) profileProofResult {
		return profileProofResult{}
	})
	// Distinct command and loaded homes prove Doctor reads loaded.HomeDir.
	setCommandHomeDirForTest(t, t.TempDir())
	var stdout, stderr bytes.Buffer
	code := runCLI(t, []string{"doctor"}, &stdout, &stderr)
	if code != 0 || stderr.Len() != 0 {
		t.Fatalf("exit=%d stdout=%q stderr=%q", code, stdout.String(), stderr.String())
	}
	want := "user-skills: warn (user-level Roundfix skill copies that an Agent harness may load in place of the project copy: " + strings.Join(descriptions, ", ") + "; or refresh with roundfix skills install --target claude && roundfix skills install --target opencode; next: roundfix skills remove --target claude --apply && roundfix skills remove --target opencode --apply)"
	lines := strings.Split(strings.TrimSuffix(stdout.String(), "\n"), "\n")
	found := false
	for i, line := range lines {
		if strings.HasPrefix(line, "skills:") {
			found = true
			if i+1 >= len(lines) || lines[i+1] != want {
				t.Fatalf("line after skills:\n%s\nwant:\n%s", stdout.String(), want)
			}
		}
	}
	if !found || strings.Count(stdout.String(), "user-skills:") != 1 {
		t.Fatalf("missing or duplicate check: %s", stdout.String())
	}
}

func TestUpdateCheckNamesStaleUserSkillCopies(t *testing.T) {
	t.Parallel()
	home, copies, descriptions := staleUserSkillsFixture(t)
	setCommandHomeDirForTest(t, home)
	setCommandWorkDirForTest(t, t.TempDir())
	configureUpdateCheckTest(t, "0.0.0-dev", func(context.Context) ([]app.Release, error) {
		t.Fatal("development build reached release lookup")
		return nil, nil
	})
	var stdout, stderr bytes.Buffer
	code := runCLI(t, []string{"update", "check"}, &stdout, &stderr)
	want := "Roundfix 0.0.0-dev is a development build; no release comparison (update.mode notify).\nBaseline: not checked (outside a Git repository).\nUser skills: copies that an Agent harness may load in place of the project copy: " + strings.Join(descriptions, ", ") + "; remove with roundfix skills remove --target claude --apply && roundfix skills remove --target opencode --apply, or refresh with roundfix skills install --target claude && roundfix skills install --target opencode.\n"
	if code != 0 || stderr.Len() != 0 || stdout.String() != want {
		t.Fatalf("exit=%d stdout=%q stderr=%q\nwant=%q", code, stdout.String(), stderr.String(), want)
	}
	stdout.Reset()
	stderr.Reset()
	code = runCLI(t, []string{"update", "check", "--format", "json"}, &stdout, &stderr)
	var report updateCheckReport
	if err := json.Unmarshal(stdout.Bytes(), &report); err != nil {
		t.Fatal(err)
	}
	if code != 0 || stderr.Len() != 0 || !reflect.DeepEqual(report.UserSkills, copies) || report.Next != "none" {
		t.Fatalf("exit=%d stderr=%q report=%+v want copies=%+v", code, stderr.String(), report, copies)
	}
	for _, copy := range copies {
		writeStaleUserSkillFixture(t, copy.Path, copy.Skill, copy.Embedded)
	}
	stdout.Reset()
	stderr.Reset()
	code = runCLI(t, []string{"update", "check", "--format", "json"}, &stdout, &stderr)
	var fields map[string]json.RawMessage
	if err := json.Unmarshal(stdout.Bytes(), &fields); err != nil {
		t.Fatal(err)
	}
	if _, exists := fields["userSkills"]; exists || code != 0 || stderr.Len() != 0 {
		t.Fatalf("current copies: exit=%d stdout=%q stderr=%q", code, stdout.String(), stderr.String())
	}
}

func staleUserSkillsFixture(t *testing.T) (string, []roundskills.StaleUserSkillCopy, []string) {
	t.Helper()
	home := t.TempDir()
	vault := filepath.Join(home, "vault", "roundfix")
	writeStaleUserSkillFixture(t, vault, "roundfix", "0.0.1")
	claude := filepath.Join(home, ".claude", "skills", "roundfix")
	if err := os.MkdirAll(filepath.Dir(claude), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(vault, claude); err != nil {
		t.Fatal(err)
	}
	qa := filepath.Join(home, ".opencode", "skills", "qa-gate")
	writeStaleUserSkillFixture(t, qa, "qa-gate", "0.0.1")
	resolved, err := filepath.EvalSymlinks(vault)
	if err != nil {
		t.Fatal(err)
	}
	versions := make(map[string]string)
	for _, owned := range roundskills.CheckReadiness().Owned {
		versions[owned.Skill] = owned.Minimum
	}
	for _, name := range []string{"roundfix", "qa-gate"} {
		if versions[name] == "" {
			t.Fatalf("missing embedded version for %s", name)
		}
	}
	copies := []roundskills.StaleUserSkillCopy{
		{Skill: "roundfix", Path: claude, Resolved: resolved, Found: "0.0.1", Embedded: versions["roundfix"], Refresh: "roundfix skills install --target claude", Remove: "roundfix skills remove --target claude --apply", Stale: true},
		{Skill: "qa-gate", Path: qa, Found: "0.0.1", Embedded: versions["qa-gate"], Refresh: "roundfix skills install --target opencode", Remove: "roundfix skills remove --target opencode --apply", Stale: true},
	}
	descriptions := []string{
		fmt.Sprintf("roundfix 0.0.1 (binary %s) at %s -> %s", versions["roundfix"], claude, resolved) + "; older than this binary",
		fmt.Sprintf("qa-gate 0.0.1 (binary %s) at %s", versions["qa-gate"], qa) + "; older than this binary",
	}
	return home, copies, descriptions
}

func writeStaleUserSkillFixture(t *testing.T, path, name, version string) {
	t.Helper()
	if err := os.MkdirAll(path, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(path, "SKILL.md"), []byte(fmt.Sprintf("---\nname: %s\nversion: %s\n---\n", name, version)), 0o644); err != nil {
		t.Fatal(err)
	}
}
