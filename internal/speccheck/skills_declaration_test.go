// IN: public Spec checks on disposable repositories with committed coverage maps.
// OUT: real repository history, providers, Verification and Daemon settlement.
// Invariant: changed covered sources need a skill-writing Task or a reviewed excuse.
package speccheck

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"roundfix/internal/gittest"
	"roundfix/internal/skillcoverage"
)

const skillsFixtureSlug = "0200-skills"
const skillsFixtureDir = "docs/specs/" + skillsFixtureSlug
const skillsFixtureSurface = "command: example"
const skillsFixtureFile = ".agents/skills/roundfix/references/example.md"
const skillsFixtureMap = `{"schemaVersion":"roundfix/skill-coverage/v1","surfaces":[{"id":"command: example","skills":[".agents/skills/roundfix/references/example.md"],"sources":["internal/cli/example.go","internal/cli/sub/","internal/cli/example_*.go"]},{"id":"command: uncovered","uncovered":"no skill describes it","sources":["internal/cli/uncovered.go"]}]}`

func skillsRepo(t *testing.T, declaration string, mapFirst bool) string {
	t.Helper()
	root := t.TempDir()
	gittest.InitRepo(t, root, "-b", "main")
	if mapFirst {
		writeReceiptFixture(t, root, skillcoverage.MapPath, skillsFixtureMap)
		receiptCommit(t, root)
	}
	writeReceiptFixture(t, root, skillsFixtureDir+"/_prd.md", "---\nstatus: active\n---\n# Skills fixture\n"+declaration)
	skillsTasks(t, root, "interface: internal/cli/example.go")
	receiptCommit(t, root)
	return root
}

func skillsTasks(t *testing.T, root string, contexts ...string) {
	t.Helper()
	manifest := "---\nschema: spec-tasks/v1\nqa: declined\nqa_reason: detector fixture\ngraph:\n  nodes:\n"
	for i, context := range contexts {
		id := fmt.Sprintf("task_%02d", i+1)
		needs := "[]"
		if i > 0 {
			needs = "[task_01]"
		}
		manifest += "    - id: " + id + "\n      file: " + id + ".md\n      needs: " + needs + "\n"
		writeReceiptFixture(t, root, skillsFixtureDir+"/"+id+".md", "---\ntask: "+id+"\nspec: "+skillsFixtureSlug+"\nstatus: pending\ntype: backend\n---\n# Task: fixture\n\n## Context\n\n- "+context+"\n\n## Verification\n\n- `true`\n")
	}
	writeReceiptFixture(t, root, skillsFixtureDir+"/_tasks.md", manifest+"---\n")
}

func skillsCheck(t *testing.T, root string, stage Stage) Result {
	t.Helper()
	result, err := CheckStage(filepath.Join(root, "docs/specs"), root, skillsFixtureSlug, stage)
	if err != nil {
		t.Fatal(err)
	}
	return result
}

func skillsCount(t *testing.T, result Result, code string, want int) []Finding {
	t.Helper()
	findings := receiptCodeFindings(result, code)
	if len(findings) != want {
		t.Fatalf("%s findings = %#v, want %d", code, findings, want)
	}
	for _, f := range findings {
		if f.Severity != SeverityError {
			t.Fatalf("finding severity = %s", f.Severity)
		}
	}
	return findings
}

func TestSkillsTaskIsRequiredForAChangedSurface(t *testing.T) {
	root := skillsRepo(t, "", true)
	for _, kind := range []string{"interface", "creates", "deletes"} {
		for _, source := range []string{"internal/cli/example.go", "internal/cli/sub/other.go", "internal/cli/example_more.go"} {
			skillsTasks(t, root, kind+": "+source)
			result := skillsCheck(t, root, StageTasks)
			f := skillsCount(t, result, CodeSkillsUntasked, 1)[0]
			for _, text := range []string{skillsFixtureSurface, "task_01", source, skillsFixtureFile} {
				if !strings.Contains(f.Summary, text) {
					t.Fatalf("summary %q omits %q", f.Summary, text)
				}
			}
			if !strings.Contains(RenderText(result, VerificationCoverage{}), "[error] SC-SKILLS-UNTASKED: ") {
				t.Fatal("missing transcript prefix")
			}
		}
	}
	for _, kind := range []string{"interface", "creates"} {
		skillsTasks(t, root, "interface: internal/cli/example.go", kind+": "+skillsFixtureFile)
		skillsCount(t, skillsCheck(t, root, StageAll), CodeSkillsUntasked, 0)
	}
	for _, context := range []string{"interface: internal/cli/uncovered.go", "instruction: internal/cli/example.go", "interface: internal/cli/unknown.go"} {
		skillsTasks(t, root, context)
		skillsCount(t, skillsCheck(t, root, StageAll), CodeSkillsUntasked, 0)
	}
	for _, kind := range []string{"deletes", "instruction"} {
		skillsTasks(t, root, "interface: internal/cli/example.go", kind+": "+skillsFixtureFile)
		skillsCount(t, skillsCheck(t, root, StageAll), CodeSkillsUntasked, 1)
	}
}

func TestAnUnchangedDeclarationExcusesASurface(t *testing.T) {
	root := skillsRepo(t, "## Skills\n\n- unchanged: "+skillsFixtureSurface+" — help is unchanged\n", true)
	skillsCount(t, skillsCheck(t, root, StageAll), CodeSkillsUntasked, 1)
	for _, kind := range []string{"interface", "creates"} {
		skillsTasks(t, root, "interface: internal/cli/example.go", kind+": "+skillcoverage.MapPath)
		skillsCount(t, skillsCheck(t, root, StageAll), CodeSkillsUntasked, 0)
	}
	skillsTasks(t, root, "interface: internal/cli/example.go", "instruction: "+skillcoverage.MapPath)
	skillsCount(t, skillsCheck(t, root, StageAll), CodeSkillsUntasked, 1)
}

func TestAMalformedSkillsDeclarationIsReported(t *testing.T) {
	for _, entry := range []string{"- changed: command: example — reason", "- unchanged: unknown — reason", "- unchanged: command: example — "} {
		t.Run(entry, func(t *testing.T) {
			root := skillsRepo(t, "## Skills\n\n"+entry+"\n", true)
			for _, stage := range []Stage{StagePRD, StageTechSpec, StageTasks, StageAll} {
				f := skillsCount(t, skillsCheck(t, root, stage), CodeSkillsMalformed, 1)[0]
				if f.Where[0] != (Location{Path: skillsFixtureDir + "/_prd.md", Line: 7}) {
					t.Fatalf("location = %#v", f.Where)
				}
				if stage == StagePRD || stage == StageTechSpec {
					skillsCount(t, skillsCheck(t, root, stage), CodeSkillsUntasked, 0)
				}
			}
		})
	}
	root := skillsRepo(t, "", true)
	skillsCount(t, skillsCheck(t, root, StageAll), CodeSkillsMalformed, 0)
}

func TestTheSkillsRuleStartsAtTheMapHorizon(t *testing.T) {
	assertSkip := func(t *testing.T, root, reason string) {
		t.Helper()
		result := skillsCheck(t, root, StageAll)
		for _, code := range []string{CodeSkillsMalformed, CodeSkillsUntasked} {
			skillsCount(t, result, code, 0)
			found := false
			for _, s := range result.Skipped {
				found = found || s.Code == code && s.Missing == reason
			}
			if !found {
				t.Fatalf("missing %s skip %q: %#v", code, reason, result.Skipped)
			}
		}
	}
	t.Run("no map", func(t *testing.T) { assertSkip(t, skillsRepo(t, "", false), skillcoverage.MapPath) })
	t.Run("older PRD", func(t *testing.T) {
		root := skillsRepo(t, "## Skills\nmalformed\n", false)
		writeReceiptFixture(t, root, skillcoverage.MapPath, skillsFixtureMap)
		receiptCommit(t, root)
		assertSkip(t, root, "a PRD committed at or after "+skillcoverage.MapPath)
	})
	t.Run("uncommitted PRD is held", func(t *testing.T) {
		root := t.TempDir()
		gittest.InitRepo(t, root, "-b", "main")
		writeReceiptFixture(t, root, skillcoverage.MapPath, skillsFixtureMap)
		receiptCommit(t, root)
		writeReceiptFixture(t, root, skillsFixtureDir+"/_prd.md", "---\nstatus: active\n---\n## Skills\nbad entry\n")
		skillsTasks(t, root, "interface: internal/cli/example.go")
		skillsCount(t, skillsCheck(t, root, StageAll), CodeSkillsUntasked, 1)
		skillsCount(t, skillsCheck(t, root, StagePRD), CodeSkillsMalformed, 1)
	})
	t.Run("missing PRD", func(t *testing.T) {
		root := skillsRepo(t, "", true)
		if err := os.Remove(filepath.Join(root, skillsFixtureDir, "_prd.md")); err != nil {
			t.Fatal(err)
		}
		assertSkip(t, root, skillsFixtureDir+"/_prd.md")
	})
}

func TestSkillsDeclarationIgnoresQATasks(t *testing.T) {
	for _, tc := range []struct {
		name, source, covering string
		want                   int
	}{
		{"QA cannot cover a source", "interface: internal/cli/example.go", "interface: " + skillsFixtureFile, 1},
		{"QA cannot record a review", "interface: internal/cli/example.go", "interface: " + skillcoverage.MapPath, 1},
		{"QA source needs no skills Task", "interface: internal/cli/unknown.go", "interface: internal/cli/example.go", 0},
	} {
		t.Run(tc.name, func(t *testing.T) {
			root := skillsRepo(t, "## Skills\n- unchanged: "+skillsFixtureSurface+" — no text change\n", true)
			skillsTasks(t, root, tc.source, tc.covering)
			manifestPath := filepath.Join(root, skillsFixtureDir, "_tasks.md")
			manifest, err := os.ReadFile(manifestPath)
			if err != nil {
				t.Fatal(err)
			}
			writeReceiptFixture(t, root, skillsFixtureDir+"/_tasks.md", strings.Replace(string(manifest), "qa: declined\nqa_reason: detector fixture", "qa: task_02", 1))
			taskPath := filepath.Join(root, skillsFixtureDir, "task_02.md")
			task, err := os.ReadFile(taskPath)
			if err != nil {
				t.Fatal(err)
			}
			writeReceiptFixture(t, root, skillsFixtureDir+"/task_02.md", strings.Replace(string(task), "type: backend", "type: qa", 1))
			skillsCount(t, skillsCheck(t, root, StageAll), CodeSkillsUntasked, tc.want)
		})
	}
}

func TestSkillsDeclarationUsesOldestMapAddition(t *testing.T) {
	root := skillsRepo(t, "", true)
	if err := os.Remove(filepath.Join(root, skillcoverage.MapPath)); err != nil {
		t.Fatal(err)
	}
	receiptCommit(t, root)
	writeReceiptFixture(t, root, skillcoverage.MapPath, skillsFixtureMap)
	receiptCommit(t, root)
	// Re-adding the map does not exempt a PRD introduced after its first addition.
	skillsCount(t, skillsCheck(t, root, StageAll), CodeSkillsUntasked, 1)
}

func TestSkillsDeclarationMissingGraphSkipsBothCodes(t *testing.T) {
	root := skillsRepo(t, "", true)
	if err := os.Remove(filepath.Join(root, skillsFixtureDir, "_tasks.md")); err != nil {
		t.Fatal(err)
	}
	result := skillsCheck(t, root, StageAll)
	for _, code := range []string{CodeSkillsMalformed, CodeSkillsUntasked} {
		found := false
		for _, skip := range result.Skipped {
			found = found || skip.Code == code && skip.Missing == skillsFixtureDir+"/_tasks.md"
		}
		if !found {
			t.Fatalf("missing %s graph skip: %#v", code, result.Skipped)
		}
	}
}
