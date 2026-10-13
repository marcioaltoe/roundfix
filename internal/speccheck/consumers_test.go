// Suite: Task consumer declarations
// Invariant: pending non-QA Tasks name literal suites run by themselves or by a
// checkpoint between the Task and the QA gate; fenced examples never declare.
// Boundary IN: temporary Task Graphs and artifacts through CheckStage.
// Boundary OUT: semantic judgment and execution of authored commands.
package speccheck_test

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"roundfix/internal/spec"
	"roundfix/internal/speccheck"
)

func TestConsumersGapFollowsTheSection(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name, section string
		count         int
	}{
		{"missing", "", 1},
		{"empty", "## Consumers\n", 1},
		{"prose", "## Consumers\n\nThere are no consumers.\n", 1},
		{"none without reason", "## Consumers\n\nNone.\n", 1},
		{"separate reason", "## Consumers\n\nNone.\n\nNo shared contract changes.\n", 1},
		{"plain suite", "## Consumers\n\n- Report — consumers: ./consumer\n", 1},
		{"empty suite", "## Consumers\n\n- Report — consumers: ` `\n", 1},
		{"empty contract", "## Consumers\n\n- — consumers: `./consumer`\n", 1},
		{"blank contract", "## Consumers\n\n-    — consumers: `./consumer`\n", 1},
		{"trailing prose", "## Consumers\n\n- Report — consumers: `./consumer` and others\n", 1},
		{"bad separator", "## Consumers\n\n- Report — consumers: `./consumer` `./other`\n", 1},
		{"bad checkpoint", "## Consumers\n\n- checkpoint: task_02.md\n", 1},
		{"fenced section", "```md\n## Consumers\n\nNone. No shared contract changes.\n```\n", 1},
		{"fenced entry", "## Consumers\n\n~~~md\n- Report — consumers: `./consumer`\n~~~\n", 1},
		{"none with reason", consumerNoneSection, 0},
		{"wrapped reason", "## Consumers\n\nNone.\nNo shared contract changes.\n", 0},
		{"suite", "## Consumers\n\n- Report — consumers: `./consumer`\n", 0},
		{"wrapped suites", "## Consumers\n\n- Report — consumers: `./consumer`,\n  `./other`\n", 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			root := writeConsumerFixture(t, consumerDeclinedQA, consumerFixtureTask{id: "task_01", section: tt.section, commands: []string{"go test ./consumer ./other"}})
			result := checkProofFixture(t, root, speccheck.StageTasks)
			requireConsumerCount(t, result, tt.count)
			for _, finding := range findingsWithCode(result, speccheck.CodeConsumersUnproven) {
				if !hasLocation(finding, "docs/specs/proof/task_01.md") {
					t.Errorf("locations = %#v, want declaring Task", finding.Where)
				}
				for _, text := range []string{"## Consumers", "write-tasks", "Consumers of a shared contract"} {
					if !strings.Contains(finding.Fix, text) {
						t.Errorf("fix = %q, want %q", finding.Fix, text)
					}
				}
			}
		})
	}
}

func TestConsumersGapNamesASuiteNoVerificationRuns(t *testing.T) {
	t.Parallel()
	const section = "## Consumers\n\n- Report — consumers: `./consumer`, `./other`\n"
	root := writeConsumerFixture(t, consumerDeclinedQA, consumerFixtureTask{id: "task_01", section: section, commands: []string{"go test ./unrelated"}})
	result := checkProofFixture(t, root, speccheck.StageAll)
	requireConsumerCount(t, result, 2)
	for _, suite := range []string{"./consumer", "./other"} {
		found := false
		for _, finding := range findingsWithCode(result, speccheck.CodeConsumersUnproven) {
			if strings.Contains(finding.Summary, suite) && hasExactLocation(finding, "docs/specs/proof/task_01.md", 13) {
				found = true
			}
		}
		if !found {
			t.Errorf("findings = %#v, want suite %s at its declaration", result.Findings, suite)
		}
	}
	marker := filepath.Join(t.TempDir(), "executed")
	writeConsumerTask(t, root, consumerFixtureTask{id: "task_01", section: section, commands: []string{"go test ./consumer", "go test ./other; touch " + marker}})
	requireConsumerCount(t, checkProofFixture(t, root, speccheck.StageAll), 0)
	if _, err := os.Stat(marker); !os.IsNotExist(err) {
		t.Fatalf("authored command marker stat = %v, want absent", err)
	}
	// The contract deliberately matches literal substrings without shell inference.
	writeConsumerTask(t, root, consumerFixtureTask{id: "task_01", section: section, commands: []string{"go test ./consumer-extra ./other"}})
	requireConsumerCount(t, checkProofFixture(t, root, speccheck.StageAll), 0)
}

func TestConsumersCheckpointSitsBetweenTheTaskAndTheGate(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name, checkpoint, qa string
		needs                []string
		want                 int
	}{
		{"missing", "task_99", "qa: task_05\n", []string{"task_02"}, 1},
		{"self", "task_01", "qa: task_05\n", []string{"task_02"}, 1},
		{"qa task", "task_05", "qa: task_05\n", []string{"task_02"}, 1},
		{"parallel task", "task_03", "qa: task_05\n", nil, 1},
		{"no gate", "task_03", "", []string{"task_02"}, 1},
		{"direct", "task_02", "qa: task_05\n", []string{"task_02"}, 0},
		{"transitive", "task_03", "qa: task_05\n", []string{"task_02"}, 0},
		{"declined", "task_03", consumerDeclinedQA, []string{"task_02"}, 0},
		{"declined still needs task", "task_03", consumerDeclinedQA, nil, 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			tasks := []consumerFixtureTask{
				{id: "task_01", section: "## Consumers\n\n- checkpoint: " + tt.checkpoint + "\n"},
				{id: "task_02", needs: []string{"task_01"}, section: consumerNoneSection, commands: []string{"go test ./consumer"}},
				{id: "task_03", needs: tt.needs, section: consumerNoneSection, commands: []string{"go test ./consumer"}},
				{id: "task_04", needs: []string{"task_02", "task_03"}, section: consumerNoneSection},
			}
			if strings.Contains(tt.qa, "task_05") {
				tasks = append(tasks, consumerFixtureTask{id: "task_05", kind: "qa", needs: []string{"task_04"}})
			}
			root := writeConsumerFixture(t, tt.qa, tasks...)
			requireConsumerCount(t, checkProofFixture(t, root, speccheck.StageTasks), tt.want)
			if tt.want > 0 {
				writeConsumerTask(t, root, consumerFixtureTask{id: "task_01", section: "## Consumers\n\n- checkpoint: task_02\n"})
				if tt.qa == "" {
					writeConsumerManifest(t, root, consumerDeclinedQA, tasks)
				}
				requireConsumerCount(t, checkProofFixture(t, root, speccheck.StageTasks), 0)
			}
		})
	}
	t.Run("suite must run in a named valid checkpoint", func(t *testing.T) {
		t.Parallel()
		const declaration = "## Consumers\n\n- Report — consumers: `./consumer`\n"
		root := writeConsumerFixture(t, "qa: task_03\n",
			consumerFixtureTask{id: "task_01", section: declaration},
			consumerFixtureTask{id: "task_02", needs: []string{"task_01"}, section: consumerNoneSection, commands: []string{"go test ./consumer"}},
			consumerFixtureTask{id: "task_03", needs: []string{"task_02"}, kind: "qa"},
		)
		requireConsumerCount(t, checkProofFixture(t, root, speccheck.StageTasks), 1)
		writeConsumerTask(t, root, consumerFixtureTask{id: "task_01", section: declaration + "- checkpoint: task_02\n"})
		requireConsumerCount(t, checkProofFixture(t, root, speccheck.StageTasks), 0)
		writeConsumerTask(t, root, consumerFixtureTask{id: "task_02", section: consumerNoneSection, commands: []string{"go test ./unrelated"}})
		requireConsumerCount(t, checkProofFixture(t, root, speccheck.StageTasks), 1)
		writeConsumerTask(t, root, consumerFixtureTask{id: "task_02", section: consumerNoneSection, commands: []string{"go test ./consumer"}})
		requireConsumerCount(t, checkProofFixture(t, root, speccheck.StageTasks), 0)
		writeConsumerTask(t, root, consumerFixtureTask{id: "task_01", section: declaration + "- checkpoint: task_03\n"})
		requireConsumerCount(t, checkProofFixture(t, root, speccheck.StageTasks), 2)
		writeConsumerTask(t, root, consumerFixtureTask{id: "task_01", section: declaration + "- checkpoint: task_02\n"})
		requireConsumerCount(t, checkProofFixture(t, root, speccheck.StageTasks), 0)
	})
}

func TestConsumersSkipCompletedAndQATasks(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name, status string
		want         int
	}{{"pending", "pending", 1}, {"in progress", "in_progress", 1}, {"failed", "failed", 1}, {"completed", "completed", 0}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			root := writeConsumerFixture(t, "qa: task_02\n", consumerFixtureTask{id: "task_01", status: tt.status}, consumerFixtureTask{id: "task_02", kind: "qa", needs: []string{"task_01"}})
			result := checkProofFixture(t, root, speccheck.StageTasks)
			requireConsumerCount(t, result, tt.want)
			for _, finding := range findingsWithCode(result, speccheck.CodeConsumersUnproven) {
				if !hasLocation(finding, "docs/specs/proof/task_01.md") {
					t.Errorf("reported QA Task: %#v", finding)
				}
			}
			writeConsumerTask(t, root, consumerFixtureTask{id: "task_01", status: tt.status, section: consumerNoneSection})
			requireConsumerCount(t, checkProofFixture(t, root, speccheck.StageTasks), 0)
		})
	}
}

func TestConsumersRunAtTheTasksStage(t *testing.T) {
	t.Parallel()
	root := writeConsumerFixture(t, consumerDeclinedQA, consumerFixtureTask{id: "task_01"})
	for _, stage := range []speccheck.Stage{speccheck.StagePRD, speccheck.StageTechSpec} {
		result := checkProofFixture(t, root, stage)
		requireConsumerCount(t, result, 0)
		if !hasSkip(result, speccheck.CodeConsumersUnproven, "stage "+string(stage)) {
			t.Errorf("skips = %#v, want stage skip", result.Skipped)
		}
	}
	for _, stage := range []speccheck.Stage{speccheck.StageAll, speccheck.StageTasks} {
		result := checkProofFixture(t, root, stage)
		requireConsumerCount(t, result, 1)
		gap := speccheck.Result{Slug: result.Slug, Findings: findingsWithCode(result, speccheck.CodeConsumersUnproven)}
		speccheck.PromoteGaps(&gap)
		if gap.Findings[0].Severity != speccheck.SeverityError {
			t.Errorf("strict severity = %q, want error", gap.Findings[0].Severity)
		}
	}
	if err := os.Remove(filepath.Join(root, "docs", "specs", "proof", "_techspec.md")); err != nil {
		t.Fatal(err)
	}
	requireConsumerCount(t, checkProofFixture(t, root, speccheck.StageAll), 1)
	writeConsumerTask(t, root, consumerFixtureTask{id: "task_01", section: consumerNoneSection})
	for _, stage := range []speccheck.Stage{speccheck.StageAll, speccheck.StageTasks} {
		requireConsumerCount(t, checkProofFixture(t, root, stage), 0)
	}
	if err := os.Remove(filepath.Join(root, "docs", "specs", "proof", "_tasks.md")); err != nil {
		t.Fatal(err)
	}
	for _, stage := range []speccheck.Stage{speccheck.StageAll, speccheck.StageTasks} {
		result := checkProofFixture(t, root, stage)
		requireConsumerCount(t, result, 0)
		if !hasSkip(result, speccheck.CodeConsumersUnproven, "_tasks.md") {
			t.Errorf("skips = %#v, want absent graph skip", result.Skipped)
		}
	}
	if err := os.Remove(filepath.Join(root, "docs", "specs", "proof", "_prd.md")); err != nil {
		t.Fatal(err)
	}
	result := checkProofFixture(t, root, speccheck.StageAll)
	requireConsumerCount(t, result, 0)
	if !hasSkip(result, speccheck.CodeConsumersUnproven, "_prd.md") {
		t.Errorf("skips = %#v, want absent PRD skip", result.Skipped)
	}
}

const consumerNoneSection = "## Consumers\n\nNone. This fixture changes no shared contract.\n"
const consumerDeclinedQA = "qa: declined\nqa_reason: fixture only\n"

type consumerFixtureTask struct {
	id, status, kind, section string
	needs, commands           []string
}

func writeConsumerFixture(t *testing.T, qa string, tasks ...consumerFixtureTask) string {
	t.Helper()
	root := writeProofFixture(t, "", "", "## Absence Evidence\n\nNone. No component is new.\n\n## Coverage Map\n\nDeviations from the PRD: None.\n")
	writeCitationFixtureFile(t, root, "docs/specs/proof/_prd.md", "---\nspec: proof\nstatus: active\n---\n\n"+proofConstraints("")+"## Success Metrics\n\nNone. No outcome is measured.\n")
	for _, task := range tasks {
		writeConsumerTask(t, root, task)
	}
	writeConsumerManifest(t, root, qa, tasks)
	return root
}
func writeConsumerManifest(t *testing.T, root, qa string, tasks []consumerFixtureTask) {
	t.Helper()
	var manifest strings.Builder
	manifest.WriteString("---\nschema: spec-tasks/v1\nspec: proof\n" + qa + "graph:\n  nodes:\n")
	for _, task := range tasks {
		fmt.Fprintf(&manifest, "    - id: %s\n      file: %s.md\n      needs: [%s]\n", task.id, task.id, strings.Join(task.needs, ", "))
	}
	manifest.WriteString("---\n\n# Consumer fixture Task Graph\n")
	writeCitationFixtureFile(t, root, "docs/specs/proof/_tasks.md", manifest.String())
}
func writeConsumerTask(t *testing.T, root string, task consumerFixtureTask) {
	t.Helper()
	if task.status == "" {
		task.status = "pending"
	}
	if task.kind == "" {
		task.kind = "backend"
	}
	var content strings.Builder
	fmt.Fprintf(&content, "---\ntask: %s\nspec: proof\nstatus: %s\ntype: %s\ncomplexity: low\n---\n\n# Task: Consumer fixture\n\n%s\n## Verification\n\n", task.id, task.status, task.kind, task.section)
	for _, command := range task.commands {
		fmt.Fprintf(&content, "- `%s`\n", command)
	}
	if len(task.commands) == 0 {
		command := "test -f task-output.txt"
		if task.kind == "qa" {
			command = spec.DerivedQAVerification("proof")[0]
		}
		fmt.Fprintf(&content, "- `%s`\n", command)
	}
	writeCitationFixtureFile(t, root, "docs/specs/proof/"+task.id+".md", content.String())
}
func requireConsumerCount(t *testing.T, result speccheck.Result, count int) {
	t.Helper()
	requireProofCount(t, result, speccheck.CodeConsumersUnproven, count)
}
