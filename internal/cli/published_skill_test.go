package cli

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	roundskills "roundfix/skills"
)

const publishedCommandFixture = `skills:
  authoring:
    claude:
      prd: {model: opus, reasoning_effort: xhigh}
      techspec: {model: opus, reasoning_effort: xhigh}
`

func publishedCommandRepository(t *testing.T, all bool) (string, string) {
	t.Helper()
	home, root := withCLIWorkspace(t)
	writeAuthoringCommandFile(t, root, ".roundfixrc.yml", publishedCommandFixture)
	files, err := roundskills.Files()
	if err != nil {
		t.Fatal(err)
	}
	for _, file := range files {
		if file.Path == "write-prd/SKILL.md" || all && file.Path == "write-techspec/SKILL.md" {
			writeAuthoringCommandFile(t, root, "skills/"+file.Path, string(file.Data))
		}
	}
	return home, root
}

func linkPublishedCommandRepository(t *testing.T, root string) {
	t.Helper()
	if err := os.Symlink(filepath.Join("..", ".agents", "skills"), filepath.Join(root, ".claude", "skills")); err != nil {
		t.Fatal(err)
	}
}

func TestThePublishedSkillSurfaceTranscripts(t *testing.T) {
	t.Parallel()
	t.Run("transcript_6", func(t *testing.T) {
		_, root := publishedCommandRepository(t, false)
		files, err := roundskills.Files()
		if err != nil {
			t.Fatal(err)
		}
		want := fmt.Sprintf("Installed Roundfix skill for project: %s (%d file(s))\n", filepath.Join(root, ".agents", "skills"), len(files)) +
			"Withheld declared header: write-prd (this repository publishes skills/write-prd/SKILL.md)\n" +
			"Applied declared header: write-techspec (model opus, effort xhigh)\n" +
			"Generated subagent: .claude/agents/roundfix-prd-author.md\n" +
			"Generated subagent: .claude/agents/roundfix-techspec-author.md\n"
		assertAuthoringCommand(t, []string{"skills", "install", "--target", "project"}, want, "", 0)
	})
	t.Run("transcript_7", func(t *testing.T) {
		_, root := publishedCommandRepository(t, false)
		installAuthoringReportFixture(t)
		linkPublishedCommandRepository(t, root)
		path := ".agents/skills/write-prd/SKILL.md"
		data, err := os.ReadFile(filepath.Join(root, filepath.FromSlash(path)))
		if err != nil {
			t.Fatal(err)
		}
		writeAuthoringCommandFile(t, root, path, strings.Replace(string(data), "---\n", "---\nmodel: opus\n", 1))
		wantOut := "Roundfix skill check passed: " + strings.Join(roundskills.Names(), ", ") + "\n"
		wantErr := "roundfix: warning: the installed files do not match the Authoring Model Declaration (skills.authoring.claude):\n" +
			"  write-prd holds a header, and this repository publishes skills/write-prd/SKILL.md\n" +
			"  next: roundfix skills install --target project\n"
		assertAuthoringCommand(t, []string{"skills", "check"}, wantOut, wantErr, 0)
	})
}

func TestTheChecksCountAPublishedSkillAsApplied(t *testing.T) {
	t.Parallel()
	for _, name := range []string{"skills_check", "doctor", "update_check"} {
		t.Run(name, func(t *testing.T) {
			_, root := publishedCommandRepository(t, true)
			installAuthoringReportFixture(t)
			linkPublishedCommandRepository(t, root)
			switch name {
			case "skills_check":
				assertAuthoringCommand(t, []string{"skills", "check"}, "Roundfix skill check passed: "+strings.Join(roundskills.Names(), ", ")+"\n", "", 0)
			case "doctor":
				configureAuthoringDoctorFixture(t)
				var stdout, stderr bytes.Buffer
				code := runCLI(t, []string{"doctor"}, &stdout, &stderr)
				if code != 0 || stderr.Len() != 0 {
					t.Fatalf("exit=%d stdout=%s stderr=%s", code, &stdout, &stderr)
				}
				assertAuthoringLineAfter(t, stdout.String(), "skills:", "authoring-models: ok (claude: prd, techspec applied)", false)
			case "update_check":
				configureAuthoringUpdateFixture(t)
				var report updateCheckReport
				raw := runUpdateCheckTest(t, "update", "check", "--format", "json")
				if err := json.Unmarshal([]byte(raw), &report); err != nil {
					t.Fatal(err)
				}
				want := `{"state":"applied","harness":"claude","declared":["prd","techspec"],"problems":[]}`
				data, err := json.Marshal(report.AuthoringModels)
				if err != nil || string(data) != want {
					t.Fatalf("authoring=%s err=%v", data, err)
				}
				if out := runUpdateCheckTest(t, "update", "check"); strings.Contains(out, "Authoring models:") {
					t.Fatalf("applied text=%s", out)
				}
			}
		})
	}
	// A header is named through each check's existing report without changing exit.
	for _, name := range []string{"doctor_header", "update_header"} {
		t.Run(name, func(t *testing.T) {
			_, root := publishedCommandRepository(t, true)
			installAuthoringReportFixture(t)
			linkPublishedCommandRepository(t, root)
			writeAuthoringCommandFile(t, root, ".agents/skills/write-prd/SKILL.md", "---\nmodel: opus\n---\n")
			problem := "write-prd holds a header, and this repository publishes skills/write-prd/SKILL.md"
			if name == "doctor_header" {
				configureAuthoringDoctorFixture(t)
				var stdout, stderr bytes.Buffer
				code := runCLI(t, []string{"doctor"}, &stdout, &stderr)
				want := "authoring-models: warn (the installed files do not match the Authoring Model Declaration: " + problem + "; next: roundfix skills install --target project)"
				if code != 0 || stderr.Len() != 0 {
					t.Fatalf("exit=%d stderr=%s", code, &stderr)
				}
				assertAuthoringLineAfter(t, stdout.String(), "skills:", want, false)
			} else {
				configureAuthoringUpdateFixture(t)
				out := runUpdateCheckTest(t, "update", "check")
				want := "Authoring models: the installed files do not match the declaration: " + problem + "; apply with roundfix skills install --target project.\n"
				if !strings.Contains(out, want) {
					t.Fatalf("stdout=%s want=%s", out, want)
				}
			}
		})
	}
}

func TestBaselineUpdateReportsAWithheldHeader(t *testing.T) {
	t.Parallel()
	_, root := publishedCommandRepository(t, false)
	var stderr bytes.Buffer
	result, err := runBaselineUpdateSkillsStageWith(t.Context(), baselineUpdateSkillsRequest{Repository: root}, authoringStageDependencies(t, root, &stderr))
	if err != nil || result.Status != baselineUpdateSkillsVerified || stderr.Len() != 0 {
		t.Fatalf("result=%+v err=%v stderr=%s", result, err, &stderr)
	}
	report := newBaselineUpdateResult()
	report.Skills = result
	var output bytes.Buffer
	if err := writeBaselineUpdateResult(report, true, &output); err != nil {
		t.Fatal(err)
	}
	assertAuthoringJSON(t, output.Bytes(), `{"headers":["write-techspec"],"withheld":["write-prd"],"subagents":[".claude/agents/roundfix-prd-author.md",".claude/agents/roundfix-techspec-author.md"]}`)
	output.Reset()
	if err := writeBaselineUpdateResult(report, false, &output); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(output.String(), "Authoring models: 1 header(s), 2 generated subagent(s), 0 removed\n") {
		t.Fatalf("text=%s", &output)
	}
}

func TestAUserLevelCopyOfAPublishedSkillIsReportedByItsVersionOnly(t *testing.T) {
	t.Parallel()
	for _, name := range []string{"a_current_copy_is_not_named", "an_older_copy_hides_no_header"} {
		t.Run(name, func(t *testing.T) {
			home, root := publishedCommandRepository(t, true)
			installAuthoringReportFixture(t)
			linkPublishedCommandRepository(t, root)
			version := bundledUserSkillVersion(t, "write-prd")
			if name == "an_older_copy_hides_no_header" {
				version = "0.0.1"
			}
			path := filepath.Join(home, ".claude", "skills", "write-prd")
			writeStaleUserSkillFixture(t, path, "write-prd", version)
			var stdout, stderr bytes.Buffer
			code := runCLI(t, []string{"skills", "check"}, &stdout, &stderr)
			if code != 0 || strings.Contains(stderr.String(), "hides the declared header") {
				t.Fatalf("exit=%d stderr=%s", code, &stderr)
			}
			if name == "a_current_copy_is_not_named" && stderr.Len() != 0 {
				t.Fatalf("current copy named: %s", &stderr)
			}
			if name == "an_older_copy_hides_no_header" && !strings.Contains(stderr.String(), "; older than this binary\n") {
				t.Fatalf("older copy not named: %s", &stderr)
			}
			configureAuthoringDoctorFixture(t)
			stdout.Reset()
			stderr.Reset()
			code = runCLI(t, []string{"doctor"}, &stdout, &stderr)
			if code != 0 || stderr.Len() != 0 || strings.Contains(stdout.String(), "hides the declared header") {
				t.Fatalf("exit=%d stdout=%s stderr=%s", code, &stdout, &stderr)
			}
			if strings.Contains(stdout.String(), "user-skills:") != (name == "an_older_copy_hides_no_header") {
				t.Fatalf("doctor copy report=%s", &stdout)
			}
			configureAuthoringUpdateFixture(t)
			text := runUpdateCheckTest(t, "update", "check")
			if strings.Contains(text, "hides the declared header") || strings.Contains(text, "User skills:") != (name == "an_older_copy_hides_no_header") {
				t.Fatalf("text=%s", text)
			}
			var report updateCheckReport
			raw := runUpdateCheckTest(t, "update", "check", "--format", "json")
			if err := json.Unmarshal([]byte(raw), &report); err != nil {
				t.Fatal(err)
			}
			if name == "a_current_copy_is_not_named" {
				if len(report.UserSkills) != 0 {
					t.Fatalf("current=%+v", report.UserSkills)
				}
			} else {
				if len(report.UserSkills) != 1 {
					t.Fatalf("copies=%+v", report.UserSkills)
				}
				copy := report.UserSkills[0]
				want := "write-prd " + version + " (binary " + bundledUserSkillVersion(t, "write-prd") + ") at " + path + "; older than this binary"
				if copy.HidesDeclaredHeader || !copy.Stale || copy.Describe() != want {
					t.Fatalf("copy=%+v", copy)
				}
			}
			data, err := os.ReadFile(filepath.Join(path, "SKILL.md"))
			if err != nil || !strings.Contains(string(data), "version: "+version) {
				t.Fatalf("copy changed: %q %v", data, err)
			}
		})
	}
}
