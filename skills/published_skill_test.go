package skills

import (
	"bytes"
	"io/fs"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func publishedSkillDeclaration() *AuthoringDeclaration {
	return &AuthoringDeclaration{Harness: "claude", Entries: map[string]AuthoringEntry{
		"prd": {Model: "opus", Effort: "xhigh"}, "techspec": {Model: "opus", Effort: "xhigh"},
	}}
}

func publishSkillFixture(t *testing.T, root, name string) {
	t.Helper()
	data, err := fs.ReadFile(embedded, name+"/SKILL.md")
	if err != nil {
		t.Fatal(err)
	}
	writeAuthoringFixture(t, root, "skills/"+name+"/SKILL.md", string(data))
}

func TestInstallationWithholdsTheHeaderOfAPublishedSkill(t *testing.T) {
	t.Parallel()
	for _, name := range []string{"one_published", "second_run_changes_no_byte", "every_authoring_skill_published", "none_published", "a_previous_header_is_removed"} {
		t.Run(name, func(t *testing.T) {
			root := t.TempDir()
			declaration := publishedSkillDeclaration()
			if name == "a_previous_header_is_removed" {
				if _, err := Install(t.Context(), InstallRequest{ProjectDir: root, Authoring: declaration}); err != nil {
					t.Fatal(err)
				}
				data, err := os.ReadFile(filepath.Join(root, ".agents/skills/write-prd/SKILL.md"))
				if err != nil || !hasAuthoringHeader(data) {
					t.Fatalf("previous header absent: %v", err)
				}
			}
			var headers, withheld []string
			if name == "every_authoring_skill_published" {
				for _, phase := range AuthoringPhases() {
					declaration.Entries[phase.Key] = AuthoringEntry{Model: "opus", Effort: "xhigh"}
					publishSkillFixture(t, root, phase.Skill)
					withheld = append(withheld, phase.Skill)
				}
			} else if name == "none_published" {
				headers = []string{"write-prd", "write-techspec"}
			} else {
				publishSkillFixture(t, root, "write-prd")
				headers, withheld = []string{"write-techspec"}, []string{"write-prd"}
			}
			result, err := Install(t.Context(), InstallRequest{ProjectDir: root, Authoring: declaration})
			if err != nil {
				t.Fatal(err)
			}
			var subagents []string
			for _, phase := range AuthoringPhases() {
				entry, declared := declaration.Entries[phase.Key]
				if !declared {
					continue
				}
				want, err := fs.ReadFile(embedded, phase.Skill+"/SKILL.md")
				if err != nil {
					t.Fatal(err)
				}
				if name == "none_published" || phase.Skill == "write-techspec" && name != "every_authoring_skill_published" {
					want, err = ApplyDeclaredHeader(want, entry)
					if err != nil {
						t.Fatal(err)
					}
				}
				got, err := os.ReadFile(filepath.Join(root, ".agents/skills", phase.Skill, "SKILL.md"))
				if err != nil || !bytes.Equal(got, want) {
					t.Fatalf("installed %s differs: %v", phase.Skill, err)
				}
				path, want, err := GeneratedSubagent("claude", phase.Key, entry)
				if err != nil {
					t.Fatal(err)
				}
				subagents = append(subagents, path)
				got, err = os.ReadFile(filepath.Join(root, filepath.FromSlash(path)))
				if err != nil || !bytes.Equal(got, want) {
					t.Fatalf("subagent %s differs: %v", path, err)
				}
			}
			want := AuthoringInstall{Headers: headers, Withheld: withheld, Subagents: subagents}
			if !reflect.DeepEqual(result.Authoring, want) {
				t.Fatalf("result=%+v want=%+v", result.Authoring, want)
			}
			if name == "second_run_changes_no_byte" {
				before := inspectAuthoringTreeSnapshot(t, root)
				second, err := Install(t.Context(), InstallRequest{ProjectDir: root, Authoring: declaration})
				if err != nil || !reflect.DeepEqual(second, result) {
					t.Fatalf("second=%+v err=%v", second, err)
				}
				if !reflect.DeepEqual(before, inspectAuthoringTreeSnapshot(t, root)) {
					t.Fatal("second install changed tree")
				}
			}
		})
	}
}

func TestInspectAuthoringAcceptsAPublishedSkillWithoutAHeader(t *testing.T) {
	t.Parallel()
	link := AuthoringProblem{".claude/skills/write-techspec does not lead to .agents/skills/write-techspec, so Claude Code does not load the declared header; link it there", false}
	header := AuthoringProblem{"write-prd holds a header, and this repository publishes skills/write-prd/SKILL.md", true}
	for _, name := range []string{"no_header_is_no_problem", "a_header_is_the_problem", "the_tree_is_unchanged", "missing_installed_skill", "effort_only", "nested_model_is_no_header"} {
		t.Run(name, func(t *testing.T) {
			root := t.TempDir()
			declaration := publishedSkillDeclaration()
			publishSkillFixture(t, root, "write-prd")
			if _, err := Install(t.Context(), InstallRequest{ProjectDir: root, Authoring: declaration}); err != nil {
				t.Fatal(err)
			}
			want := []AuthoringProblem{link}
			switch name {
			case "a_header_is_the_problem", "effort_only", "the_tree_is_unchanged":
				key := "model: opus"
				if name == "effort_only" {
					key = "effort: high"
				}
				writeAuthoringFixture(t, root, ".agents/skills/write-prd/SKILL.md", "---\nname: write-prd\n"+key+"\n---\n")
				want = []AuthoringProblem{header, link}
			case "missing_installed_skill":
				if err := os.Remove(filepath.Join(root, ".agents/skills/write-prd/SKILL.md")); err != nil {
					t.Fatal(err)
				}
			case "nested_model_is_no_header":
				writeAuthoringFixture(t, root, ".agents/skills/write-prd/SKILL.md", "---\nname: write-prd\nmetadata:\n  model: nested\n  effort: nested\n---\n")
			}
			before := inspectAuthoringTreeSnapshot(t, root)
			got := InspectAuthoring(root, declaration)
			if got.Harness != "claude" || !reflect.DeepEqual(got.Declared, []string{"prd", "techspec"}) || !reflect.DeepEqual(got.Problems, want) {
				t.Fatalf("state=%+v want problems=%+v", got, want)
			}
			if !reflect.DeepEqual(before, inspectAuthoringTreeSnapshot(t, root)) {
				t.Fatal("inspection changed files, directories or links")
			}
			if name == "a_header_is_the_problem" {
				if _, err := Install(t.Context(), InstallRequest{ProjectDir: root, Authoring: declaration}); err != nil {
					t.Fatal(err)
				}
				if got := InspectAuthoring(root, declaration); !reflect.DeepEqual(got.Problems, []AuthoringProblem{link}) {
					t.Fatalf("repair=%+v", got)
				}
			}
		})
	}
}

func TestTheSetupReferenceNamesThePublishedSkillRule(t *testing.T) {
	t.Parallel()
	phrases := []struct{ name, text string }{
		{"a_published_skill_receives_no_header", "A repository that publishes an authoring skill under `skills/<name>/` receives no Declared Skill Header for it"},
		{"the_subagents_are_generated", "still generates every declared subagent"},
		{"the_checks_agree", "The checks count the declaration as applied when that installed skill holds no `model` and no `effort`, and report a header found there."},
		{"an_older_binary_writes_it", "A Roundfix older than the release that introduced this rule writes the header into a published skill, so upgrade it before installing there."},
	}
	for name, source := range roundfixLayoutSources() {
		t.Run(name, func(t *testing.T) {
			data, err := fs.ReadFile(source, "roundfix/references/setup.md")
			if err != nil {
				t.Fatal(err)
			}
			text := strings.Join(strings.Fields(string(data)), " ")
			for _, phrase := range phrases {
				t.Run(phrase.name, func(t *testing.T) {
					if count := strings.Count(text, phrase.text); count != 1 {
						t.Fatalf("phrase count=%d want=1: %s", count, phrase.text)
					}
				})
			}
		})
	}
}
