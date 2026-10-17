// Suite: guide Task surface record declarations
// Invariant: a guide Task records its own fingerprint at its own settlement.
// Boundary IN: disposable Spec artifacts through CheckStage and strict promotion.
// Boundary OUT: command execution, help inference and real repository state.
package speccheck_test

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"roundfix/internal/skillcoverage"
	"roundfix/internal/speccheck"
)

const surfaceRecordCode = speccheck.CodeSurfaceRecordUntasked
const surfaceGuide = "docs/user-guide/commands/example.md"
const surfaceRecordRequirement = "## Requirements\n\n1. MUST re-record the fingerprints with\n   `-record-skill-coverage`.\n"

func surfaceRecordRepo(t *testing.T, qa string, tasks ...consumerFixtureTask) string {
	t.Helper()
	root := writeConsumerFixture(t, qa, tasks...)
	writeCitationFixtureFile(t, root, skillcoverage.RecordPath, "{}\n")
	return root
}

func surfaceRecordContext(kind, record string) string {
	return "## Context\n\n- " + kind + ": " + surfaceGuide + "\n" + record
}

func requireSurfaceRecordFindings(t *testing.T, result speccheck.Result, want int) []speccheck.Finding {
	t.Helper()
	requireProofCount(t, result, surfaceRecordCode, want)
	findings := findingsWithCode(result, surfaceRecordCode)
	for _, finding := range findings {
		if finding.Severity != speccheck.SeverityGap || len(finding.Where) != 1 {
			t.Fatalf("finding = %#v, want one-location gap", finding)
		}
		if !strings.Contains(finding.Fix, `write-tasks rule "Skills Tasks for changed Behavior Surfaces"`) {
			t.Errorf("fix = %q, want authoring rule", finding.Fix)
		}
	}
	return findings
}

func TestSurfaceRecordGapNamesAGuideTaskWithoutItsRecord(t *testing.T) {
	t.Parallel()
	for _, kind := range []string{"interface", "creates", "deletes"} {
		for _, shape := range []string{"0268 plan", "record without step", "step without record", "record and step in another Task", "record only deleted", "record only instruction", "step only in Verification"} {
			t.Run(kind+"/"+shape, func(t *testing.T) {
				t.Parallel()
				record, requirement := "", ""
				var other []consumerFixtureTask
				var commands []string
				switch shape {
				case "record without step", "step only in Verification":
					record = "- interface: " + skillcoverage.RecordPath + "\n"
				case "record only deleted":
					record = "- deletes: " + skillcoverage.RecordPath + "\n"
					requirement = surfaceRecordRequirement
				case "record only instruction":
					record = "- instruction: " + skillcoverage.RecordPath + "\n"
					requirement = surfaceRecordRequirement
				case "step without record":
					requirement = surfaceRecordRequirement
				case "record and step in another Task":
					other = append(other, consumerFixtureTask{id: "task_02", needs: []string{"task_01"}, section: surfaceRecordRequirement + "\n## Context\n\n- interface: " + skillcoverage.RecordPath + "\n"})
				}
				if shape == "step only in Verification" {
					commands = []string{"go test ./example -record-skill-coverage"}
				}
				section := requirement + "\n" + surfaceRecordContext(kind, record)
				if shape == "0268 plan" {
					section += "- interface: docs/user-guide/configuration.md\n- interface: docs/user-guide/usage.md\n"
				}
				tasks := append([]consumerFixtureTask{{id: "task_01", section: section, commands: commands}}, other...)
				root := surfaceRecordRepo(t, consumerDeclinedQA, tasks...)
				if shape == "0268 plan" {
					writeCitationFixtureFile(t, root, "docs/specs/proof/_authorization.md", "## Limits\n\n- Do not edit "+skillcoverage.RecordPath+".\n")
				}
				want := 1
				if shape == "0268 plan" {
					want = 2
				}
				findings := requireSurfaceRecordFindings(t, checkProofFixture(t, root, speccheck.StageAll), want)
				wantSummary := fmt.Sprintf("Task task_01 declares user guide %s without declaring %s and naming its -record-skill-coverage step", surfaceGuide, skillcoverage.RecordPath)
				found := false
				for _, finding := range findings {
					if finding.Summary != wantSummary {
						continue
					}
					found = true
					content, err := os.ReadFile(filepath.Join(root, "docs/specs/proof/task_01.md"))
					if err != nil {
						t.Fatal(err)
					}
					before, _, _ := strings.Cut(string(content), "- "+kind+": "+surfaceGuide)
					line := strings.Count(before, "\n") + 1
					if !hasExactLocation(finding, "docs/specs/proof/task_01.md", line) {
						t.Errorf("location = %#v, want Context line %d", finding.Where, line)
					}
				}
				if !found {
					t.Fatalf("findings = %#v, want summary %q", findings, wantSummary)
				}
			})
		}
	}
}

func TestSurfaceRecordGapNamesLimitsThatForbidTheRecord(t *testing.T) {
	t.Parallel()
	root := surfaceRecordRepo(t, consumerDeclinedQA, consumerFixtureTask{id: "task_01", section: surfaceRecordRequirement + surfaceRecordContext("interface", "- interface: "+skillcoverage.RecordPath+"\n")})
	const path = "docs/specs/proof/_authorization.md"
	writeCitationFixtureFile(t, root, path, "# Authority\n\n## Regeneration\n\n"+skillcoverage.RecordPath+"\n\n## Limits\n\n- Do not edit `"+skillcoverage.RecordPath+"`.\n\n## Other\n\n"+skillcoverage.RecordPath+"\n")
	findings := requireSurfaceRecordFindings(t, checkProofFixture(t, root, speccheck.StageAll), 1)
	if !hasExactLocation(findings[0], path, 9) {
		t.Errorf("location = %#v, want Limits line 9", findings[0].Where)
	}
	if !strings.Contains(findings[0].Summary, skillcoverage.RecordPath) {
		t.Errorf("summary = %q, want record path", findings[0].Summary)
	}
}

func TestSurfaceRecordAcceptsAGuideTaskThatRecordsItself(t *testing.T) {
	t.Parallel()
	for _, shape := range []string{"interface record", "creates record", "completed", "QA Task", "no guide", "instruction guide", "non-markdown guide", "multiple guides", "no record", "record directory"} {
		t.Run(shape, func(t *testing.T) {
			t.Parallel()
			task := consumerFixtureTask{id: "task_01", section: surfaceRecordContext("interface", "")}
			qa := consumerDeclinedQA
			switch shape {
			case "interface record", "creates record", "multiple guides":
				kind := "interface"
				if shape == "creates record" {
					kind = "creates"
				}
				task.section = surfaceRecordRequirement + surfaceRecordContext("interface", "- "+kind+": "+skillcoverage.RecordPath+"\n")
				if shape == "multiple guides" {
					task.section += "- deletes: docs/user-guide/old.md\n"
				}
			case "completed":
				task.status = "completed"
			case "QA Task":
				task.kind = "qa"
				qa = "qa: task_01\n"
			case "no guide":
				task.section = "## Context\n\n- interface: internal/cli/usage.go\n"
			case "instruction guide":
				task.section = surfaceRecordContext("instruction", "")
			case "non-markdown guide":
				task.section = "## Context\n\n- interface: docs/user-guide/image.png\n"
			}
			root := surfaceRecordRepo(t, qa, task)
			if shape == "no record" || shape == "record directory" {
				if err := os.Remove(filepath.Join(root, skillcoverage.RecordPath)); err != nil {
					t.Fatal(err)
				}
				if shape == "record directory" {
					if err := os.Mkdir(filepath.Join(root, skillcoverage.RecordPath), 0o755); err != nil {
						t.Fatal(err)
					}
				}
			}
			result := checkProofFixture(t, root, speccheck.StageAll)
			requireSurfaceRecordFindings(t, result, 0)
			if (shape == "no record" || shape == "record directory") && !hasSkip(result, surfaceRecordCode, skillcoverage.RecordPath) {
				t.Errorf("skips = %#v, want record skip", result.Skipped)
			}
		})
	}
}

func TestSurfaceRecordRunsAtTheTasksStageAndFailsStrict(t *testing.T) {
	t.Parallel()
	root := surfaceRecordRepo(t, consumerDeclinedQA, consumerFixtureTask{id: "task_01", section: surfaceRecordContext("interface", "")})
	for _, stage := range []speccheck.Stage{speccheck.StagePRD, speccheck.StageTechSpec} {
		result := checkProofFixture(t, root, stage)
		requireSurfaceRecordFindings(t, result, 0)
		if !hasSkip(result, surfaceRecordCode, "stage "+string(stage)) {
			t.Errorf("skips = %#v, want stage skip", result.Skipped)
		}
	}
	for _, stage := range []speccheck.Stage{speccheck.StageAll, speccheck.StageTasks} {
		result := checkProofFixture(t, root, stage)
		findings := requireSurfaceRecordFindings(t, result, 1)
		gap := speccheck.Result{Findings: findings}
		speccheck.PromoteGaps(&gap)
		if gap.Findings[0].Severity != speccheck.SeverityError {
			t.Errorf("strict finding = %#v, want error", gap.Findings[0])
		}
	}
	if err := os.Remove(filepath.Join(root, "docs/specs/proof/_tasks.md")); err != nil {
		t.Fatal(err)
	}
	result := checkProofFixture(t, root, speccheck.StageAll)
	requireSurfaceRecordFindings(t, result, 0)
	if !hasSkip(result, surfaceRecordCode, "_tasks.md") {
		t.Errorf("skips = %#v, want graph skip", result.Skipped)
	}
}
