// Suite: owned skill Task version declarations
// Invariant: each Task editing an owned skill declares its own version work.
// Boundary IN: temporary repository records and Task Graphs through CheckStage.
// Boundary OUT: skill version generation and settlement execution.
package speccheck_test

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"roundfix/internal/speccheck"
)

const skillVersionCode = speccheck.CodeSkillVersionUntasked
const skillVersionRecord = "skills/testdata/owned-skill-versions.json"
const skillVersionStep = "## Requirements\n\n1. MUST run `-record-skill-versions`.\n\n"
const skillVersionSource = ".agents/skills/roundfix/references/spec.md"

func skillVersionRepo(t *testing.T, qa string, tasks ...consumerFixtureTask) string {
	t.Helper()
	root := writeConsumerFixture(t, qa, tasks...)
	writeCitationFixtureFile(t, root, skillVersionRecord, `{"skills":{"roundfix":[],"write-tasks":[]}}`)
	return root
}

func skillVersionDeclarations(kind string) string {
	return "## Context\n\n- " + kind + ": " + skillVersionSource + "\n- " + kind + ": skills/roundfix/references/spec.md\n- interface: .agents/skills/roundfix/SKILL.md\n- interface: skills/roundfix/SKILL.md\n- interface: " + skillVersionRecord + "\n"
}

func TestSkillVersionGapNamesATaskWithoutItsRecord(t *testing.T) {
	t.Parallel()
	for _, kind := range []string{"interface", "creates", "deletes"} {
		for _, shape := range []string{"mirror", "SKILL pair", "record", "step", "all", "instruction first", "another Task", "instruction record", "deleted record", "Verification step", "wrong mirror kind", "two skills"} {
			t.Run(kind+"/"+shape, func(t *testing.T) {
				t.Parallel()
				section := skillVersionDeclarations(kind)
				requirement := skillVersionStep
				var missing []string
				var other []consumerFixtureTask
				var commands []string
				mirror := "skills/roundfix/references/spec.md"
				switch shape {
				case "mirror", "wrong mirror kind":
					section = strings.ReplaceAll(section, "- "+kind+": "+mirror+"\n", "")
					if shape == "wrong mirror kind" {
						wrong := "deletes"
						if kind == "deletes" {
							wrong = "interface"
						}
						section += "- " + wrong + ": " + mirror + "\n"
					}
					missing = []string{"interface: " + mirror}
				case "SKILL pair":
					for _, path := range []string{".agents/skills/roundfix/SKILL.md", "skills/roundfix/SKILL.md"} {
						section = strings.ReplaceAll(section, "- interface: "+path+"\n", "")
						missing = append(missing, "interface: "+path)
					}
				case "record", "instruction record", "deleted record":
					section = strings.ReplaceAll(section, "- interface: "+skillVersionRecord+"\n", "")
					if shape == "instruction record" {
						section += "- instruction: " + skillVersionRecord + "\n"
					}
					if shape == "deleted record" {
						section += "- deletes: " + skillVersionRecord + "\n"
					}
					missing = []string{"interface: " + skillVersionRecord}
				case "step", "Verification step":
					requirement = ""
					if shape == "Verification step" {
						commands = []string{"go test ./skills -record-skill-versions"}
					}
					missing = []string{"a requirement naming -record-skill-versions"}
				case "all", "instruction first", "another Task", "two skills":
					section = "## Context\n\n- " + kind + ": " + skillVersionSource + "\n"
					requirement = ""
					if shape == "instruction first" {
						instruction := skillVersionSource
						if kind != "interface" {
							instruction += ".read-only"
						}
						section = strings.Replace(section, "## Context\n\n", "## Context\n\n- instruction: "+instruction+"\n", 1)
					}
					missing = []string{"interface: " + mirror, "interface: .agents/skills/roundfix/SKILL.md", "interface: skills/roundfix/SKILL.md", "interface: " + skillVersionRecord, "a requirement naming -record-skill-versions"}
					if shape == "another Task" {
						other = []consumerFixtureTask{{id: "task_02", needs: []string{"task_01"}, section: skillVersionStep + skillVersionDeclarations(kind)}}
					}
					if shape == "two skills" {
						section += "- interface: .agents/skills/roundfix/references/other.md\n- interface: .agents/skills/write-tasks/references/example.md\n"
						missing = append([]string{"interface: " + mirror, "interface: skills/roundfix/references/other.md"}, missing[1:]...)
					}
				}
				tasks := append([]consumerFixtureTask{{id: "task_01", section: requirement + section, commands: commands}}, other...)
				root := skillVersionRepo(t, consumerDeclinedQA, tasks...)
				result := checkProofFixture(t, root, speccheck.StageAll)
				want := 1
				if shape == "two skills" {
					want = 2
				}
				requireProofCount(t, result, skillVersionCode, want)
				findings := findingsWithCode(result, skillVersionCode)
				finding := findings[0]
				wantSummary := fmt.Sprintf("Task task_01 edits owned skill roundfix without declaring its mirror, both SKILL.md files and %s and naming its -record-skill-versions step", skillVersionRecord)
				if finding.Summary != wantSummary || finding.Severity != speccheck.SeverityGap {
					t.Errorf("finding = %#v, want gap %q", finding, wantSummary)
				}
				wantFix := "In this Task, add " + strings.Join(missing, "; ") + ", following the write-tasks rule \"Owned skill versions\"."
				if finding.Fix != wantFix {
					t.Errorf("fix = %q, want %q", finding.Fix, wantFix)
				}
				content, err := os.ReadFile(filepath.Join(root, "docs/specs/proof/task_01.md"))
				if err != nil {
					t.Fatal(err)
				}
				before, _, _ := strings.Cut(string(content), "- "+kind+": "+skillVersionSource)
				line := strings.Count(before, "\n") + 1
				if len(finding.Where) != 1 || !hasExactLocation(finding, "docs/specs/proof/task_01.md", line) {
					t.Errorf("location = %#v, want Context line %d", finding.Where, line)
				}
				if shape == "two skills" && !strings.Contains(findings[1].Summary, "owned skill write-tasks") {
					t.Errorf("second finding = %#v, want write-tasks", findings[1])
				}
			})
		}
	}
}

func TestSkillVersionGapAcceptsATaskThatRecordsItsVersion(t *testing.T) {
	t.Parallel()
	for _, shape := range []string{"complete", "creates", "deletion", "upstream", "completed", "QA Task", "instruction", "creates declarations", "two skills"} {
		t.Run(shape, func(t *testing.T) {
			t.Parallel()
			task := consumerFixtureTask{id: "task_01", section: skillVersionStep + skillVersionDeclarations("interface")}
			qa := consumerDeclinedQA
			switch shape {
			case "creates":
				task.section = skillVersionStep + skillVersionDeclarations("creates")
			case "deletion":
				task.section = skillVersionStep + skillVersionDeclarations("deletes")
			case "upstream":
				task.section = "## Context\n\n- interface: .agents/skills/upstream/SKILL.md\n"
			case "completed":
				task.section = "## Context\n\n- interface: " + skillVersionSource + "\n"
				task.status = "completed"
			case "QA Task":
				task.section = "## Context\n\n- interface: " + skillVersionSource + "\n"
				task.kind = "qa"
				qa = "qa: task_01\n"
			case "instruction":
				task.section = "## Context\n\n- instruction: " + skillVersionSource + "\n"
			case "creates declarations":
				task.section = strings.ReplaceAll(task.section, "- interface:", "- creates:")
			case "two skills":
				task.section += "- interface: .agents/skills/write-tasks/SKILL.md\n- interface: skills/write-tasks/SKILL.md\n"
			}
			root := skillVersionRepo(t, qa, task)
			requireProofCount(t, checkProofFixture(t, root, speccheck.StageAll), skillVersionCode, 0)
		})
	}
}

func TestSkillVersionGapSkipsWithoutTheRecord(t *testing.T) {
	t.Parallel()
	for _, shape := range []string{"absent record", "record directory", "absent graph", "malformed record", "unreadable record"} {
		t.Run(shape, func(t *testing.T) {
			t.Parallel()
			root := skillVersionRepo(t, consumerDeclinedQA, consumerFixtureTask{id: "task_01", section: "## Context\n\n- interface: " + skillVersionSource + "\n"})
			missing := skillVersionRecord
			if shape == "absent graph" {
				missing = "docs/specs/proof/_tasks.md"
			}
			path := filepath.Join(root, missing)
			if shape == "malformed record" {
				writeCitationFixtureFile(t, root, missing, "{broken")
			} else if shape == "unreadable record" {
				// A symlink loop gives a deterministic read error, even as root.
				if err := os.Remove(path); err != nil {
					t.Fatal(err)
				}
				if err := os.Symlink(filepath.Base(path), path); err != nil {
					t.Fatal(err)
				}
			} else {
				if err := os.Remove(path); err != nil {
					t.Fatal(err)
				}
				if shape == "record directory" {
					if err := os.Mkdir(path, 0o755); err != nil {
						t.Fatal(err)
					}
				}
			}
			result, err := speccheck.Check(filepath.Join(root, "docs/specs"), root, "proof")
			if shape == "malformed record" || shape == "unreadable record" {
				if err == nil || !strings.Contains(err.Error(), skillVersionRecord) {
					t.Fatalf("error = %v, want record error", err)
				}
				return
			}
			if err != nil {
				t.Fatal(err)
			}
			requireProofCount(t, result, skillVersionCode, 0)
			if !hasSkip(result, skillVersionCode, missing) {
				t.Errorf("skips = %#v, want %s", result.Skipped, missing)
			}
		})
	}
}

func TestSkillVersionGapRunsAtTheTasksStageAndFailsStrict(t *testing.T) {
	t.Parallel()
	root := skillVersionRepo(t, consumerDeclinedQA, consumerFixtureTask{id: "task_01", section: "## Context\n\n- interface: " + skillVersionSource + "\n"})
	for _, stage := range []speccheck.Stage{speccheck.StagePRD, speccheck.StageTechSpec} {
		result := checkProofFixture(t, root, stage)
		requireProofCount(t, result, skillVersionCode, 0)
		if !hasSkip(result, skillVersionCode, "stage "+string(stage)) {
			t.Errorf("skips = %#v, want stage skip", result.Skipped)
		}
	}
	for _, stage := range []speccheck.Stage{speccheck.StageAll, speccheck.StageTasks} {
		result := checkProofFixture(t, root, stage)
		requireProofCount(t, result, skillVersionCode, 1)
		gap := speccheck.Result{Findings: findingsWithCode(result, skillVersionCode)}
		speccheck.PromoteGaps(&gap)
		if gap.Findings[0].Severity != speccheck.SeverityError {
			t.Errorf("strict finding = %#v, want error", gap.Findings[0])
		}
	}
}
