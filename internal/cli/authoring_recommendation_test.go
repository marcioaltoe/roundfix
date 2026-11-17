package cli

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	roundskills "roundfix/skills"
)

const authoringRecommendationBlock = `  skills:
    authoring:
      claude:
        idea:
          model: fable
          reasoning_effort: xhigh
        prd:
          model: fable
          reasoning_effort: high
        techspec:
          model: fable
          reasoning_effort: high
        tasks:
          model: opus
          reasoning_effort: high
        research:
          model: sonnet
          reasoning_effort: medium
`

const authoringRecommendationDoctorLine = "authoring-models: warn (the authoring skills are installed without an Authoring Model Declaration: write-idea, write-prd, write-techspec, write-tasks; next: add the Recommended Authoring Declaration below to .roundfixrc.yml, then run roundfix skills install --target project)"
const authoringRecommendationUpdateLine = "Authoring models: the authoring skills are installed without a declaration: write-idea, write-prd, write-techspec, write-tasks; add the Recommended Authoring Declaration below to .roundfixrc.yml, then apply with roundfix skills install --target project."

func TestDoctorRecommendsAnAuthoringDeclarationWhenNoneIsDeclared(t *testing.T) {
	t.Parallel()
	for _, name := range []string{"undeclared", "partly_published", "published", "no_authoring_skill", "leftover_generated", "declared"} {
		t.Run(name, func(t *testing.T) {
			home, root := authoringRecommendationFixture(t, name)
			configureAuthoringDoctorFixture(t)
			before := snapshotDoctorPath(t, root)
			beforeHome := snapshotDoctorPath(t, home)
			var stdout, stderr bytes.Buffer
			code := runCLI(t, []string{"doctor"}, &stdout, &stderr)
			if code != 0 || stderr.Len() != 0 {
				t.Fatalf("exit=%d stdout=%s stderr=%s", code, &stdout, &stderr)
			}
			want := ""
			switch name {
			case "undeclared", "partly_published":
				want = authoringRecommendationDoctorLine + "\n" + authoringRecommendationBlock
				if name == "partly_published" {
					want = strings.Replace(want, "write-idea, write-prd, write-techspec, write-tasks", "write-idea, write-techspec", 1)
				}
			case "leftover_generated":
				want = "authoring-models: warn (the installed files do not match the Authoring Model Declaration: .claude/agents/roundfix-prd-author.md has no declared entry; next: roundfix skills install --target project)\n"
			case "declared":
				want = "authoring-models: ok (claude: prd, research applied)\n"
			}
			assertAuthoringRecommendationPosition(t, stdout.String(), "skills:", "authoring-models:", want, "baseline-rules:")
			if after := snapshotDoctorPath(t, root); !reflect.DeepEqual(before, after) {
				t.Fatal("Doctor changed the repository")
			}
			if after := snapshotDoctorPath(t, home); !reflect.DeepEqual(beforeHome, after) {
				t.Fatal("Doctor changed the fixture home")
			}
		})
	}
}

func TestUpdateCheckRecommendsAnAuthoringDeclarationWhenNoneIsDeclared(t *testing.T) {
	t.Parallel()
	for _, name := range []string{"undeclared", "partly_published", "published", "no_authoring_skill", "leftover_generated", "declared"} {
		t.Run(name, func(t *testing.T) {
			home, root := authoringRecommendationFixture(t, name)
			configureAuthoringUpdateFixture(t)
			before := snapshotDoctorPath(t, root)
			beforeHome := snapshotDoctorPath(t, home)
			raw := runUpdateCheckTest(t, "update", "check", "--format", "json")
			var fields map[string]json.RawMessage
			if err := json.Unmarshal([]byte(raw), &fields); err != nil {
				t.Fatal(err)
			}
			wantJSON, wantText := "", ""
			switch name {
			case "undeclared", "partly_published":
				installed := `["write-idea","write-prd","write-techspec","write-tasks"]`
				wantText = authoringRecommendationUpdateLine + "\n" + authoringRecommendationBlock
				if name == "partly_published" {
					installed = `["write-idea","write-techspec"]`
					wantText = strings.Replace(wantText, "write-idea, write-prd, write-techspec, write-tasks", "write-idea, write-techspec", 1)
				}
				block := strings.TrimPrefix(strings.ReplaceAll(authoringRecommendationBlock, "\n  ", "\n"), "  ")
				encoded, err := json.Marshal(block)
				if err != nil {
					t.Fatal(err)
				}
				wantJSON = `{"state":"undeclared","harness":"claude","declared":[],"problems":[],"installed":` + installed + `,"recommendation":` + string(encoded) + `}`
			case "leftover_generated":
				wantJSON = `{"state":"unapplied","harness":"claude","declared":[],"problems":[".claude/agents/roundfix-prd-author.md has no declared entry"],"next":"roundfix skills install --target project"}`
				wantText = "Authoring models: the installed files do not match the declaration: .claude/agents/roundfix-prd-author.md has no declared entry; apply with roundfix skills install --target project.\n"
			case "declared":
				wantJSON = `{"state":"applied","harness":"claude","declared":["prd","research"],"problems":[]}`
			}
			if got := string(fields["authoringModels"]); got != wantJSON {
				t.Fatalf("authoringModels=%s want=%s", got, wantJSON)
			}
			if string(fields["next"]) != `"none"` {
				t.Fatalf("changed next: %s", raw)
			}
			output := runUpdateCheckTest(t, "update", "check")
			assertAuthoringRecommendationPosition(t, output, "Baseline:", "Authoring models:", wantText, "")
			if strings.Contains(output, "Next:") {
				t.Fatalf("changed next: %s", output)
			}
			if after := snapshotDoctorPath(t, root); !reflect.DeepEqual(before, after) {
				t.Fatal("Update Check changed the repository")
			}
			if after := snapshotDoctorPath(t, home); !reflect.DeepEqual(beforeHome, after) {
				t.Fatal("development Update Check changed the fixture home")
			}
		})
	}
}

func TestSurfaceTranscript1UpdateCheckRecommendsTheDeclaration(t *testing.T) {
	t.Parallel()
	authoringRecommendationFixture(t, "undeclared")
	configureAuthoringUpdateFixture(t)
	var stdout, stderr bytes.Buffer
	code := runCLI(t, []string{"update", "check"}, &stdout, &stderr)
	lines := strings.SplitN(stdout.String(), "\n", 3)
	if code != 0 || stderr.Len() != 0 || len(lines) != 3 {
		t.Fatalf("exit=%d stdout=%q stderr=%q", code, &stdout, &stderr)
	}
	// The first transcript ellipsis permits only version and Baseline lines.
	if lines[0] != "Roundfix 0.0.0-dev is a development build; no release comparison (update.mode notify)." || !strings.HasPrefix(lines[1], "Baseline:") {
		t.Fatalf("unexpected transcript prefix: %q", stdout.String())
	}
	want := authoringRecommendationUpdateLine + "\n" + authoringRecommendationBlock
	// The final ellipsis permits only Next when the Baseline has a ready plan.
	if strings.HasPrefix(lines[1], "Baseline: plan ready ") {
		want += "Next: apply it with the Roundfix skill's --update mode; update.mode notify applies nothing on its own.\n"
	}
	if lines[2] != want {
		t.Fatalf("transcript stdout=%q want suffix=%q", stdout.String(), want)
	}
}

func authoringRecommendationFixture(t *testing.T, name string) (string, string) {
	t.Helper()
	if name == "declared" {
		return authoringReportFixture(t, true)
	}
	home, root := withCLIWorkspace(t)
	installAuthoringReportFixture(t)
	switch name {
	case "partly_published", "published":
		names := []string{"write-prd", "write-tasks"}
		if name == "published" {
			names = []string{"write-idea", "write-prd", "write-techspec", "write-tasks"}
		}
		for _, skill := range names {
			writeAuthoringCommandFile(t, root, "skills/"+skill+"/SKILL.md", "# Published authoring skill\n")
		}
	case "no_authoring_skill":
		for _, phase := range roundskills.AuthoringPhases() {
			if err := os.RemoveAll(filepath.Join(root, ".agents", "skills", phase.Skill)); err != nil {
				t.Fatal(err)
			}
		}
	case "leftover_generated":
		writeAuthoringCommandFile(t, root, ".claude/agents/roundfix-prd-author.md", "---\nname: roundfix-prd-author\n---\n\n<!-- roundfix:generated skills.authoring.claude.prd -->\n")
	}
	return home, root
}

func assertAuthoringRecommendationPosition(t *testing.T, output, preceding, prefix, want, following string) {
	t.Helper()
	lines := strings.SplitAfter(output, "\n")
	for i, line := range lines {
		if !strings.HasPrefix(line, preceding) {
			continue
		}
		rest := strings.Join(lines[i+1:], "")
		if want == "" {
			if strings.Contains(output, prefix) || strings.Contains(output, "  skills:\n") {
				t.Fatalf("unexpected recommendation: %s", output)
			}
		} else if !strings.HasPrefix(rest, want) {
			t.Fatalf("output after %s=%q want prefix=%q", preceding, rest, want)
		}
		if following == "" && rest != want {
			t.Fatalf("unexpected trailing output: %q want=%q", rest, want)
		}
		if following != "" && !strings.HasPrefix(strings.TrimPrefix(rest, want), following) {
			t.Fatalf("wrong following check: %s", output)
		}
		return
	}
	t.Fatalf("missing preceding %s in %s", preceding, output)
}
