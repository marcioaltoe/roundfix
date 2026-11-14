package cli

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"roundfix/internal/app"
	roundconfig "roundfix/internal/config"
	roundskills "roundfix/skills"
)

const authoringHeaderProblem = "write-prd header is not the declared one"
const authoringLinkProblem = ".claude/skills/write-prd does not lead to .agents/skills/write-prd, so Claude Code does not load the declared header; link it there"
const authoringDiffersProblem = ".claude/agents/roundfix-prd-author.md differs from the declaration"
const authoringMissingProblem = ".claude/agents/roundfix-researcher.md is missing"

func TestSkillsCheckWarnsOfAnUnappliedDeclarationAndKeepsItsExit(t *testing.T) {
	t.Parallel()
	t.Run("declaration", func(t *testing.T) {
		_, root := authoringReportFixture(t, false)
		driftAuthoringReportFixture(t, root)
		nested := filepath.Join(root, "nested")
		if err := os.Mkdir(nested, 0o755); err != nil {
			t.Fatal(err)
		}
		setCommandWorkDirForTest(t, nested)
		wantOut := "Roundfix skill check passed: " + strings.Join(roundskills.Names(), ", ") + "\n"
		wantErr := "roundfix: warning: the installed files do not match the Authoring Model Declaration (skills.authoring.claude):\n  " + authoringHeaderProblem + "\n  " + authoringLinkProblem + "\n  " + authoringDiffersProblem + "\n  " + authoringMissingProblem + "\n  next: roundfix skills install --target project\n"
		assertAuthoringCommand(t, []string{"skills", "check"}, wantOut, wantErr, 0)
		installAuthoringReportFixture(t)
		wantErr = "roundfix: warning: the installed files do not match the Authoring Model Declaration (skills.authoring.claude):\n  " + authoringLinkProblem + "\n"
		assertAuthoringCommand(t, []string{"skills", "check"}, wantOut, wantErr, 0)
		linkAuthoringReportFixture(t, root)
		assertAuthoringCommand(t, []string{"skills", "check"}, wantOut, "", 0)
	})
	t.Run("not_read", func(t *testing.T) {
		home, root := authoringReportFixture(t, false)
		// The declaration reader suppresses even the loader's User Config warning.
		writeUpdateCheckUserConfig(t, home, authoringCommandFixture)
		writeAuthoringCommandFile(t, root, ".roundfixrc.yml", "skills: {authoring: {unknown: {}}}\n")
		wantErr := fmt.Sprintf("roundfix: warning: Authoring Model Declaration not read: parse config %q: skills.authoring.unknown is not a supported harness; this Roundfix supports claude\n", filepath.Join(root, ".roundfixrc.yml"))
		assertAuthoringCommand(t, []string{"skills", "check"}, "Roundfix skill check passed: "+strings.Join(roundskills.Names(), ", ")+"\n", wantErr, 0)
	})
	t.Run("unversioned", func(t *testing.T) {
		_, root := authoringReportFixture(t, false)
		updateCommandDependenciesForTest(t, func(deps *commandDependencies) {
			deps.checkOwnedSkills = func() roundskills.BundleReadiness {
				return roundskills.BundleReadiness{Owned: []roundskills.OwnedSkillReadiness{{Skill: "roundfix", State: roundskills.ReadinessUnversioned}}}
			}
		})
		assertAuthoringCommand(t, []string{"skills", "check"}, "Roundfix skill check unversioned: roundfix\n", "roundfix: warning: the installed files do not match the Authoring Model Declaration (skills.authoring.claude):\n  "+authoringLinkProblem+"\n", 0)
		// A failed bundle check never adds declaration diagnostics.
		writeAuthoringCommandFile(t, root, ".roundfixrc.yml", "skills: invalid\n")
		updateCommandDependenciesForTest(t, func(deps *commandDependencies) {
			deps.checkOwnedSkills = func() roundskills.BundleReadiness {
				return roundskills.BundleReadiness{Diagnostics: []roundskills.Diagnostic{{Path: "bundle", Message: "fixture failure"}}}
			}
		})
		assertAuthoringCommand(t, []string{"skills", "check"}, "", "roundfix: Roundfix skill check failed:\n  bundle: fixture failure\n", exitRunFailed)
	})
}

func TestDoctorReportsTheAuthoringModels(t *testing.T) {
	t.Parallel()
	for _, name := range []string{"ok", "warn", "warn_without_next", "after_user_skills"} {
		t.Run(name, func(t *testing.T) {
			home, root := authoringReportFixture(t, name == "ok" || name == "after_user_skills")
			var problems []string
			if name == "warn" {
				driftAuthoringReportFixture(t, root)
				problems = authoringReportFixtureProblems()
			}
			if name == "warn_without_next" {
				problems = []string{authoringLinkProblem}
			}
			if name == "after_user_skills" {
				writeStaleUserSkillFixture(t, filepath.Join(home, ".claude", "skills", "roundfix"), "roundfix", "0.0.1")
			}
			configureAuthoringDoctorFixture(t)
			var stdout, stderr bytes.Buffer
			code := runCLI(t, []string{"doctor"}, &stdout, &stderr)
			want := "authoring-models: ok (claude: prd, research applied)"
			if len(problems) > 0 {
				want = "authoring-models: warn (the installed files do not match the Authoring Model Declaration: " + strings.Join(problems, "; ")
				if name == "warn" {
					want += "; next: roundfix skills install --target project"
				}
				want += ")"
			}
			if code != 0 || stderr.Len() != 0 {
				t.Fatalf("exit=%d stdout=%s stderr=%s", code, &stdout, &stderr)
			}
			assertAuthoringLineAfter(t, stdout.String(), "skills:", want, name == "after_user_skills")
		})
	}
	t.Run("project_config_not_read", func(t *testing.T) {
		home, root := authoringReportFixture(t, true)
		loaded := roundconfig.Loaded{Config: roundconfig.Builtin(), HomeDir: home, GitRoot: root, RepositoryVersionErr: errors.New("newer repository fixture")}
		// Even a supplied declaration must not leak into a recovery report.
		loaded.Config.Skills.Authoring.Claude = &roundconfig.AuthoringHarness{PRD: &roundconfig.AuthoringEntry{Model: "opus"}}
		configureAuthoringDoctorLoaded(t, loaded)
		var stdout, stderr bytes.Buffer
		if code := runCLI(t, []string{"doctor"}, &stdout, &stderr); code != exitRunFailed || strings.Contains(stdout.String(), "authoring-models:") {
			t.Fatalf("exit=%d stdout=%s stderr=%s", code, &stdout, &stderr)
		}
		if report := updateCheckAuthoringModels(loaded); report != nil {
			t.Fatalf("Project Config not read: %+v", report)
		}
	})
	t.Run("outside_repository", func(t *testing.T) {
		home, _ := withCLIWorkspace(t)
		configureAuthoringDoctorLoaded(t, roundconfig.Loaded{Config: roundconfig.Builtin(), HomeDir: home})
		var stdout, stderr bytes.Buffer
		if code := runCLI(t, []string{"doctor"}, &stdout, &stderr); code != exitRunFailed || strings.Contains(stdout.String(), "authoring-models:") {
			t.Fatalf("exit=%d stdout=%s stderr=%s", code, &stdout, &stderr)
		}
	})
}

func TestUpdateCheckReportsTheAuthoringModels(t *testing.T) {
	t.Parallel()
	for _, name := range []string{"applied", "unapplied", "unapplied_without_next"} {
		t.Run(name, func(t *testing.T) {
			home, root := authoringReportFixture(t, name == "applied")
			// Force userSkills to be present, proving JSON and text placement.
			writeStaleUserSkillFixture(t, filepath.Join(home, ".claude", "skills", "roundfix"), "roundfix", "0.0.1")
			problems := []string{}
			next := ""
			if name == "unapplied" {
				driftAuthoringReportFixture(t, root)
				problems = authoringReportFixtureProblems()
				next = "roundfix skills install --target project"
			}
			if name == "unapplied_without_next" {
				problems = []string{authoringLinkProblem}
			}
			configureAuthoringUpdateFixture(t)
			raw := runUpdateCheckTest(t, "update", "check", "--format", "json")
			var report updateCheckReport
			if err := json.Unmarshal([]byte(raw), &report); err != nil {
				t.Fatal(err)
			}
			state := "applied"
			if len(problems) > 0 {
				state = "unapplied"
			}
			want := &updateCheckAuthoring{State: state, Harness: "claude", Declared: []string{"prd", "research"}, Problems: problems, Next: next}
			if !reflect.DeepEqual(report.AuthoringModels, want) || report.Next != "none" || report.SchemaVersion != updateCheckSchema {
				t.Fatalf("report=%+v authoring=%+v want=%+v", report, report.AuthoringModels, want)
			}
			var fields map[string]json.RawMessage
			if err := json.Unmarshal([]byte(raw), &fields); err != nil {
				t.Fatal(err)
			}
			// Assert every member, including [] instead of null and omitted next.
			wantJSON := fmt.Sprintf(`{"state":%q,"harness":"claude","declared":["prd","research"],"problems":%s`, state, mustAuthoringJSON(t, problems))
			if next != "" {
				wantJSON += `,"next":"roundfix skills install --target project"`
			}
			wantJSON += "}"
			if string(fields["authoringModels"]) != wantJSON {
				t.Fatalf("authoringModels=%s want=%s", fields["authoringModels"], wantJSON)
			}
			userIndex := strings.Index(raw, `"userSkills":`)
			authoringIndex := strings.Index(raw, `"authoringModels":`)
			nextIndex := strings.LastIndex(raw, `"next":`)
			if userIndex < 0 || authoringIndex <= userIndex || nextIndex <= authoringIndex {
				t.Fatalf("field order: %s", raw)
			}
			output := runUpdateCheckTest(t, "update", "check")
			if name == "applied" {
				if strings.Contains(output, "Authoring models:") {
					t.Fatalf("applied text=%s", output)
				}
			} else {
				wantLine := "Authoring models: the installed files do not match the declaration: " + strings.Join(problems, "; ")
				if next != "" {
					wantLine += "; apply with roundfix skills install --target project"
				}
				wantLine += "."
				assertAuthoringLineAfter(t, output, "User skills:", wantLine, false)
			}
			if strings.Contains(output, "Next:") {
				t.Fatalf("declaration changed report next: %s", output)
			}
		})
	}
}

func authoringReportFixture(t *testing.T, link bool) (string, string) {
	t.Helper()
	home, root := withCLIWorkspace(t)
	writeAuthoringCommandFile(t, root, ".roundfixrc.yml", authoringCommandFixture)
	installAuthoringReportFixture(t)
	if link {
		linkAuthoringReportFixture(t, root)
	}
	return home, root
}

func installAuthoringReportFixture(t *testing.T) {
	t.Helper()
	var stdout, stderr bytes.Buffer
	if code := runCLI(t, []string{"skills", "install", "--target", "project"}, &stdout, &stderr); code != 0 || stderr.Len() != 0 {
		t.Fatalf("install exit=%d stdout=%s stderr=%s", code, &stdout, &stderr)
	}
}

func driftAuthoringReportFixture(t *testing.T, root string) {
	t.Helper()
	writeAuthoringCommandFile(t, root, ".roundfixrc.yml", strings.Replace(authoringCommandFixture, "model: opus", "model: sonnet", 1))
	if err := os.Remove(filepath.Join(root, ".claude", "agents", "roundfix-researcher.md")); err != nil {
		t.Fatal(err)
	}
}

func linkAuthoringReportFixture(t *testing.T, root string) {
	t.Helper()
	path := filepath.Join(root, ".claude", "skills", "write-prd")
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(filepath.Join(root, ".agents", "skills", "write-prd"), path); err != nil {
		t.Fatal(err)
	}
}

func authoringReportFixtureProblems() []string {
	return []string{authoringHeaderProblem, authoringLinkProblem, authoringDiffersProblem, authoringMissingProblem}
}

func configureAuthoringDoctorFixture(t *testing.T) {
	t.Helper()
	loaded, err := loadCommandConfig(commandEnvironmentForTest(t), io.Discard)
	if err != nil {
		t.Fatal(err)
	}
	configureAuthoringDoctorLoaded(t, loaded)
}

func configureAuthoringDoctorLoaded(t *testing.T, loaded roundconfig.Loaded) {
	t.Helper()
	checker := newDoctorFakeHealthChecker(CheckResult{Name: HealthCheckNode, Status: CheckStatusOK}, CheckResult{Name: HealthCheckACPX, Status: CheckStatusOK}, CheckResult{Name: HealthCheckCodex, Status: CheckStatusOK})
	withDoctorFakeLoadedAndReadiness(t, checker, loaded, func(context.Context, roundconfig.Config, []roundconfig.WorkCategory, string) profileProofResult {
		return profileProofResult{}
	})
}

func configureAuthoringUpdateFixture(t *testing.T) {
	t.Helper()
	configureUpdateCheckTest(t, "0.0.0-dev", func(context.Context) ([]app.Release, error) {
		t.Fatal("development build reached release lookup")
		return nil, nil
	})
}

func assertAuthoringCommand(t *testing.T, args []string, wantOut, wantErr string, wantCode int) {
	t.Helper()
	var stdout, stderr bytes.Buffer
	code := runCLI(t, args, &stdout, &stderr)
	if code != wantCode || stdout.String() != wantOut || stderr.String() != wantErr {
		t.Fatalf("args=%v exit=%d want=%d\nstdout=%q want=%q\nstderr=%q want=%q", args, code, wantCode, stdout.String(), wantOut, stderr.String(), wantErr)
	}
}

func assertAuthoringLineAfter(t *testing.T, output, prefix, want string, afterUser bool) {
	t.Helper()
	lines := strings.Split(strings.TrimSuffix(output, "\n"), "\n")
	for i, line := range lines {
		if !strings.HasPrefix(line, prefix) {
			continue
		}
		at := i + 1
		if afterUser {
			if at >= len(lines) || !strings.HasPrefix(lines[at], "user-skills:") {
				t.Fatalf("missing user-skills line after %s: %s", prefix, output)
			}
			at++
		}
		if at >= len(lines) || lines[at] != want {
			t.Fatalf("line after %s: %s\nwant=%s", prefix, output, want)
		}
		return
	}
	t.Fatalf("missing prefix %s in %s", prefix, output)
}

func mustAuthoringJSON(t *testing.T, value []string) string {
	t.Helper()
	data, err := json.Marshal(value)
	if err != nil {
		t.Fatal(err)
	}
	return string(data)
}
