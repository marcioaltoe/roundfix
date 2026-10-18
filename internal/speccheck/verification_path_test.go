// Suite: Verification path declarations
// Invariant: pending non-QA Tasks name only paths the plan provides or the tree holds.
// Boundary IN: temporary repositories and Task Graphs through CheckStage.
// Boundary OUT: shell execution and paths that cannot be recognized literally.
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

func TestVerificationPathUndeclaredRaisesTheMissingFile(t *testing.T) {
	t.Parallel()
	const missing = "internal/docscontract/glossary_contract_test.go"
	for _, command := range []string{
		"grep -qF -- Glossary " + missing,
		"rg -e internal/pattern -e internal/other " + missing,
		"grep -f internal/patterns " + missing,
		"grep -einternal/pattern-file " + missing,
		"grep -qFe internal/pattern " + missing,
		"rg --regexp internal/pattern " + missing,
		"rtk proxy grep -qF -- internal/pattern " + missing,
		"printf internal/printed; cat " + missing,
		"echo internal/printed && cat " + missing,
		"true || cat " + missing,
		"printf x | cat " + missing,
		"(cat " + missing + ")",
		"cat < " + missing,
		"printf x > " + missing,
		"grep -qF Glossary ./" + missing + "; touch executed",
	} {
		t.Run(command, func(t *testing.T) {
			t.Parallel()
			root := writeConsumerFixture(t, consumerDeclinedQA, consumerFixtureTask{id: "task_01", commands: []string{"echo setup", command}})
			writeCitationFixtureFile(t, root, "internal/existing.go", "package existing\n")
			for _, stage := range []speccheck.Stage{speccheck.StageAll, speccheck.StageTasks} {
				result := checkProofFixture(t, root, stage)
				if findings := findingsWithCode(result, speccheck.CodeVerifyPathUndeclared); len(findings) != 1 {
					t.Fatalf("findings = %#v, want exactly one path error", findings)
				}
				finding := findingsWithCode(result, "SC-VERIFY-PATH-UNDECLARED")[0]
				wantSummary := fmt.Sprintf("Task task_01 Verification names %s, which the Task does not declare, no Task it depends on creates, and the checked tree does not hold", missing)
				wantFix := fmt.Sprintf("Declare %s in this Task, create it in a Task this Task needs, or correct the command, following the write-tasks rule \"Verification paths\".", missing)
				if finding.Code != speccheck.CodeVerifyPathUndeclared || finding.Severity != speccheck.SeverityError || finding.Summary != wantSummary || finding.Fix != wantFix {
					t.Errorf("finding = %#v, want error %q with fix %q", finding, wantSummary, wantFix)
				}
				content, err := os.ReadFile(filepath.Join(root, "docs/specs/proof/task_01.md"))
				if err != nil {
					t.Fatal(err)
				}
				before, _, _ := strings.Cut(string(content), "- `"+command+"`")
				line := strings.Count(before, "\n") + 1
				if len(finding.Where) != 1 || !hasExactLocation(finding, "docs/specs/proof/task_01.md", line) {
					t.Errorf("location = %#v, want Verification line %d", finding.Where, line)
				}
			}
			for _, stage := range []speccheck.Stage{speccheck.StagePRD, speccheck.StageTechSpec} {
				result := checkProofFixture(t, root, stage)
				requireProofCount(t, result, speccheck.CodeVerifyPathUndeclared, 0)
				if !hasSkip(result, speccheck.CodeVerifyPathUndeclared, "stage "+string(stage)) {
					t.Errorf("missing stage skip: %#v", result.Skipped)
				}
			}
			if _, err := os.Lstat(filepath.Join(root, "executed")); !os.IsNotExist(err) {
				t.Fatalf("command execution marker: %v", err)
			}
		})
	}
}

func TestVerificationPathUndeclaredIgnoresWhatThePlanProvides(t *testing.T) {
	t.Parallel()
	const missing = "internal/new.go"
	for _, shape := range []string{"creates", "interface", "deletes", "instruction", "recorded", "dependency", "transitive dependency", "dependency directory", "directory prefix", "existing", "existing directory", "grep pattern", "rg pattern", "grep explicit pattern", "rg pattern file", "printf", "echo", "negated test", "negated bracket", "spec root", "custom spec root", "completed", "in progress", "failed", "QA Task", "no graph", "variable", "glob", "ellipsis", "unknown root", "single segment"} {
		t.Run(shape, func(t *testing.T) {
			t.Parallel()
			task := consumerFixtureTask{id: "task_03", commands: []string{"cat " + missing}}
			qa := consumerDeclinedQA
			var others []consumerFixtureTask
			switch shape {
			case "creates", "interface", "deletes", "instruction":
				task.section = "## Context\n\n- " + shape + ": " + missing + "\n"
			case "dependency", "transitive dependency", "dependency directory":
				others = []consumerFixtureTask{{id: "task_01", section: "## Context\n\n- creates: " + missing + "\n"}}
				task.needs = []string{"task_01"}
				if shape == "transitive dependency" {
					others = append(others, consumerFixtureTask{id: "task_02", needs: []string{"task_01"}})
					task.needs = []string{"task_02"}
				}
				if shape == "dependency directory" {
					task.commands = []string{"cat ./internal/"}
				}
			case "directory prefix":
				task.section = "## Context\n\n- creates: internal/new/file.go\n"
				task.commands = []string{"cat ./internal/new/"}
			case "grep pattern":
				task.commands = []string{"grep -qF -- " + missing + " internal/existing.go"}
			case "rg pattern":
				task.commands = []string{"rg " + missing + " internal/existing.go"}
			case "grep explicit pattern":
				task.commands = []string{"grep -e " + missing + " internal/existing.go"}
			case "rg pattern file":
				task.commands = []string{"rg -f " + missing + " internal/existing.go"}
			case "printf", "echo":
				task.commands = []string{shape + " '%s' " + missing}
			case "negated test":
				task.commands = []string{"test ! -f " + missing}
			case "negated bracket":
				task.commands = []string{"[ ! -e " + missing + " ]"}
			case "spec root":
				task.commands = []string{"cat docs/specs/other/missing.md"}
			case "custom spec root":
				task.commands = []string{"cat plans/other/missing.md"}
			case "completed":
				task.status = "completed"
			case "in progress":
				task.status = "in_progress"
			case "failed":
				task.status = "failed"
			case "QA Task":
				task.kind = "qa"
				task.commands = nil
				qa = "qa: task_03\n"
			case "variable":
				task.commands = []string{"cat internal/$file"}
			case "glob":
				task.commands = []string{"cat internal/*.go internal/a?.go internal/[ab].go"}
			case "ellipsis":
				task.commands = []string{"go test ./internal/..."}
			case "unknown root":
				task.commands = []string{"cat absent/file.go"}
			case "single segment":
				task.commands = []string{"cat absent.go"}
			}
			tasks := append(others, task)
			root := writeConsumerFixture(t, qa, tasks...)
			writeCitationFixtureFile(t, root, "internal/existing.go", "package existing\n")
			taskPath := filepath.Join(root, "docs/specs/proof/task_03.md")
			if shape == "recorded" {
				if err := spec.RecordTaskPaths(taskPath, []string{missing}); err != nil {
					t.Fatal(err)
				}
			}
			if shape == "existing" {
				writeCitationFixtureFile(t, root, missing, "")
			}
			if shape == "existing directory" {
				if err := os.Mkdir(filepath.Join(root, missing), 0o755); err != nil {
					t.Fatal(err)
				}
			}
			if shape == "no graph" {
				if err := os.Remove(filepath.Join(root, "docs/specs/proof/_tasks.md")); err != nil {
					t.Fatal(err)
				}
			}
			specsRoot := filepath.Join(root, "docs/specs")
			if shape == "custom spec root" {
				specsRoot = filepath.Join(root, "plans")
				if err := os.Rename(filepath.Join(root, "docs/specs"), specsRoot); err != nil {
					t.Fatal(err)
				}
			}
			result, err := speccheck.CheckStage(specsRoot, root, "proof", speccheck.StageTasks)
			if err != nil {
				t.Fatal(err)
			}
			requireProofCount(t, result, speccheck.CodeVerifyPathUndeclared, 0)
			if shape == "no graph" && !hasSkip(result, speccheck.CodeVerifyPathUndeclared, "docs/specs/proof/_tasks.md") {
				t.Errorf("missing graph skip: %#v", result.Skipped)
			}
		})
	}
}

func TestVerificationPathUndeclaredDoesNotBorrowOtherDeclarations(t *testing.T) {
	t.Parallel()
	for _, shape := range []string{"unneeded creation", "needed interface", "needed instruction", "needed deletion", "parent declaration", "similar prefix"} {
		t.Run(shape, func(t *testing.T) {
			t.Parallel()
			task := consumerFixtureTask{id: "task_02", commands: []string{"cat internal/new/file.go"}}
			source := consumerFixtureTask{id: "task_01", section: "## Context\n\n- creates: internal/new/file.go\n"}
			switch shape {
			case "needed interface", "needed instruction", "needed deletion":
				kind := map[string]string{"needed interface": "interface", "needed instruction": "instruction", "needed deletion": "deletes"}[shape]
				source.section = strings.Replace(source.section, "creates:", kind+":", 1)
				task.needs = []string{"task_01"}
			case "parent declaration":
				task.section = "## Context\n\n- interface: internal/new\n"
			case "similar prefix":
				task.commands = []string{"cat internal/ne"}
				task.section = "## Context\n\n- creates: internal/new/file.go\n"
			}
			root := writeConsumerFixture(t, consumerDeclinedQA, source, task)
			writeCitationFixtureFile(t, root, "internal/existing.go", "package existing\n")
			findings := findingsWithCode(checkProofFixture(t, root, speccheck.StageTasks), speccheck.CodeVerifyPathUndeclared)
			if len(findings) != 1 || !strings.Contains(findings[0].Summary, "Task task_02 Verification names ") {
				t.Fatalf("findings = %#v, want task_02's unsupplied path", findings)
			}
		})
	}
}
