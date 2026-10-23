// Suite: Subsumed test selectors
// Invariant: a pending non-QA Task never passes a pytest node id beside an argument of the same invocation that contains it.
// Boundary IN: temporary repositories and Task Graphs through CheckStage.
// Boundary OUT: running pytest or any Verification command.
package speccheck_test

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"roundfix/internal/speccheck"
)

func TestSubsumedSelectorRaisesANodeIdBesideItsContainingPath(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		command, selector, argument string
	}{
		{"pytest tests/test_a.py::test_new tests/test_a.py", "tests/test_a.py::test_new", "tests/test_a.py"},
		{"uv run pytest -q tests/test_a.py tests/test_a.py::test_new", "tests/test_a.py::test_new", "tests/test_a.py"},
		{"uv run --group dev pytest tests tests/test_a.py::test_new", "tests/test_a.py::test_new", "tests"},
		{"python -m pytest tests/ tests/test_a.py::TestNew", "tests/test_a.py::TestNew", "tests/"},
		{"python3.12 -m pytest . tests/test_a.py::test_new", "tests/test_a.py::test_new", "."},
		{"pytest ./tests/test_a.py tests/test_a.py::test_new", "tests/test_a.py::test_new", "./tests/test_a.py"},
		{"pytest tests/test_a.py::TestNew tests/test_a.py::TestNew::test_it", "tests/test_a.py::TestNew::test_it", "tests/test_a.py::TestNew"},
		{"pytest tests/test_a.py::test_it tests/test_a.py::test_it[new]", "tests/test_a.py::test_it[new]", "tests/test_a.py::test_it"},
		{"out=\"$(pytest tests/test_a.py::test_new tests/test_a.py 2>&1)\" || exit 1", "tests/test_a.py::test_new", "tests/test_a.py"},
		{"cd tools && .venv/bin/pytest tests/test_a.py::test_new tests", "tests/test_a.py::test_new", "tests"},
	} {
		t.Run(tc.command, func(t *testing.T) {
			t.Parallel()
			root := writeConsumerFixture(t, consumerDeclinedQA, consumerFixtureTask{id: "task_01", commands: []string{"echo setup", tc.command}})
			content, err := os.ReadFile(filepath.Join(root, "docs/specs/proof/task_01.md"))
			if err != nil {
				t.Fatal(err)
			}
			before, _, _ := strings.Cut(string(content), "- `"+tc.command+"`")
			line := strings.Count(before, "\n") + 1
			wantSummary := fmt.Sprintf("Task task_01 Verification passes the test selector %s together with %s, which contains it; pytest runs %s whole and ignores the selector, so a test that does not exist yet goes unnoticed", tc.selector, tc.argument, tc.argument)
			wantFix := fmt.Sprintf("Name only the selectors in one command and run %s in a separate command, following the write-tasks rule \"Test selectors and their containing path\".", tc.argument)
			for _, stage := range []speccheck.Stage{speccheck.StageAll, speccheck.StageTasks} {
				findings := findingsWithCode(checkProofFixture(t, root, stage), speccheck.CodeVerifySubsumedSelector)
				if len(findings) != 1 {
					t.Fatalf("stage %s findings = %#v, want exactly one selector error", stage, findings)
				}
				finding := findings[0]
				if finding.Code != "SC-VERIFY-SUBSUMED-SELECTOR" || finding.Severity != speccheck.SeverityError || finding.Summary != wantSummary || finding.Fix != wantFix {
					t.Errorf("stage %s finding = %#v, want error %q with fix %q", stage, finding, wantSummary, wantFix)
				}
				if len(finding.Where) != 1 || !hasExactLocation(finding, "docs/specs/proof/task_01.md", line) {
					t.Errorf("stage %s location = %#v, want Verification line %d", stage, finding.Where, line)
				}
			}
		})
	}
}

func TestSubsumedSelectorAcceptsSelectorsNoArgumentContains(t *testing.T) {
	t.Parallel()
	for _, command := range []string{
		"pytest tests/test_a.py::test_new",
		"pytest tests/test_a.py::test_new tests/test_b.py",
		"pytest tests/test_a.py::test_new tests/test_a.py::test_old",
		"pytest tests/test_ab.py::test_new tests/test_a",
		"pytest --keep-duplicates tests/test_a.py::test_new tests/test_a.py",
		"pytest --deselect tests/test_a.py::test_old tests/test_a.py::test_old[new]",
		"pytest -k tests/test_a.py tests/test_a.py::test_new",
		"pytest tests/test_a.py::test_new; pytest tests/test_a.py",
		"grep -qF -- tests/test_a.py::test_new tests/test_a.py",
		"pytest tests/test_a.py::test_new tests/test_a.py::test_new",
	} {
		t.Run(command, func(t *testing.T) {
			t.Parallel()
			root := writeConsumerFixture(t, consumerDeclinedQA, consumerFixtureTask{id: "task_01", commands: []string{command}})
			for _, stage := range []speccheck.Stage{speccheck.StageAll, speccheck.StageTasks} {
				if findings := findingsWithCode(checkProofFixture(t, root, stage), speccheck.CodeVerifySubsumedSelector); len(findings) != 0 {
					t.Errorf("stage %s findings = %#v, want none", stage, findings)
				}
			}
		})
	}
}

func TestSubsumedSelectorJudgesOnlyPendingNonQATasks(t *testing.T) {
	t.Parallel()
	const command = "pytest tests/test_a.py::test_new tests/test_a.py"
	for _, status := range []string{"completed", "in_progress"} {
		t.Run(status, func(t *testing.T) {
			t.Parallel()
			root := writeConsumerFixture(t, consumerDeclinedQA, consumerFixtureTask{id: "task_01", status: status, commands: []string{command}})
			for _, stage := range []speccheck.Stage{speccheck.StageAll, speccheck.StageTasks} {
				if findings := findingsWithCode(checkProofFixture(t, root, stage), speccheck.CodeVerifySubsumedSelector); len(findings) != 0 {
					t.Errorf("stage %s findings = %#v, want none for a %s Task", stage, findings, status)
				}
			}
		})
	}
	t.Run("prd stage", func(t *testing.T) {
		t.Parallel()
		root := writeConsumerFixture(t, consumerDeclinedQA, consumerFixtureTask{id: "task_01", commands: []string{command}})
		result := checkProofFixture(t, root, speccheck.StagePRD)
		if findings := findingsWithCode(result, speccheck.CodeVerifySubsumedSelector); len(findings) != 0 {
			t.Errorf("findings = %#v, want none at the PRD stage", findings)
		}
		if !hasSkip(result, speccheck.CodeVerifySubsumedSelector, "stage "+string(speccheck.StagePRD)) {
			t.Errorf("skips = %#v, want the code skipped at the PRD stage", result.Skipped)
		}
	})
}
