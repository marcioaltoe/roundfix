// Suite: partial policy agreement through isolated local command surfaces.
// Boundary IN: QA accept, archive, settle, Run Branch reports and temporary homes.
// Boundary OUT: network and external Agent execution.
package cli

import (
	"bytes"
	"fmt"
	"path/filepath"
	"strings"
	"testing"

	"roundfix/internal/delivery"
	"roundfix/internal/gittest"
	"roundfix/internal/spec"
	"roundfix/internal/store"
)

func TestQAReportAcceptArchiveAndSettleAgreeOnEveryPartialShape(t *testing.T) {
	for _, shape := range qaPartialShapes() {
		t.Run(shape.name, func(t *testing.T) {
			for _, command := range []string{"qa-report", "archive", "settle"} {
				t.Run(command, func(t *testing.T) {
					status := spec.StatusCompleted
					if command == "settle" {
						status = spec.StatusFailed
					}
					home, repo := newImplementWorkspace(t, []implementSeed{{id: "task_01", status: string(spec.StatusCompleted)}, implementQAGateSeed(string(status), "task_01")})
					dir := filepath.Join(repo, "docs", "specs", implementTestSlug)
					prd := filepath.Join(dir, "_prd.md")
					mustWrite(t, prd, mustRead(t, prd)+shape.unreachable())
					rel := filepath.ToSlash(filepath.Join("docs", "specs", implementTestSlug, "qa", "qa-report-2026-10-06.md"))
					writeQAReportTestFile(t, filepath.Join(repo, rel), shape.report())
					// Settle executes only fixture-authored Verification on committed provenance.
					gitImplement(t, repo, "add", "-A")
					gitImplement(t, repo, "commit", "-m", "docs: partial policy fixture")
					before := mustRead(t, implementTaskPath(repo, "task_qa"))
					args := []string{command, implementTestSlug}
					if command == "qa-report" {
						args = []string{command, "accept", filepath.Join(repo, rel)}
					}
					if command == "settle" {
						args = []string{command, "--spec", implementTestSlug, "--task", "task_qa"}
					}
					var stdout, stderr bytes.Buffer
					code := runCLIContext(t, t.Context(), args, &stdout, &stderr)
					accepted := shape.reason == ""
					if (code == exitOK) != accepted {
						t.Fatalf("exit=%d stdout=%q stderr=%q; want accepted=%v", code, stdout.String(), stderr.String(), accepted)
					}
					if !accepted && !strings.Contains(stderr.String(), shape.reason) {
						t.Fatalf("stderr=%q, want reason %q", stderr.String(), shape.reason)
					}
					if command == "qa-report" {
						wantStderr := ""
						wantCode := exitOK
						if !accepted {
							wantCode = exitRunFailed
							wantStderr = "roundfix: QA Report eligibility refused: " + shape.reason + "\n"
						}
						if code != wantCode || stdout.Len() != 0 || stderr.String() != wantStderr {
							t.Fatalf("accept transcript: exit=%d stdout=%q stderr=%q, want exit=%d empty stdout stderr=%q", code, stdout.String(), stderr.String(), wantCode, wantStderr)
						}
					}
					if command == "settle" && !accepted {
						want := "Settle surface: " + repo + "\nroundfix: settle QA Report is ineligible: " + shape.reason + " (report " + rel + ")\n"
						if code != exitRunFailed || stderr.String() != want || mustRead(t, implementTaskPath(repo, "task_qa")) != before {
							t.Fatalf("settle refusal: exit=%d stderr=%q, want %q and unchanged Task", code, stderr.String(), want)
						}
					}
					assertNoRunDatabase(t, home)
				})
			}
		})
	}
}

func TestSettleNamesTheReportItRefused(t *testing.T) {
	home, repo := newImplementWorkspace(t, []implementSeed{implementQAGateSeed(string(spec.StatusFailed))})
	location := configureSettleWorktreeLocation(t, repo, filepath.Join(home, "worktrees"))
	_, _, taskRef := createImplementRunWorktreeFixture(t, home, repo, location, implementTestSlug, "task_qa", store.StateUnresolved)
	dir := filepath.Join(taskRef.Path, "docs", "specs", implementTestSlug, "qa")
	writeQAReportTestFile(t, filepath.Join(dir, "qa-report-2026-10-05.md"), qaPartialShapeNamed(t, "Spec 0213 PR only").report())
	shape := qaPartialShapeNamed(t, "network without outside evidence")
	rel := filepath.ToSlash(filepath.Join("docs", "specs", implementTestSlug, "qa", "qa-report-2026-10-06-02.md"))
	writeQAReportTestFile(t, filepath.Join(taskRef.Path, rel), shape.report())
	// An eligible checkout report must not replace the refused Task surface report.
	writeQAReportTestFile(t, filepath.Join(repo, rel), qaPartialShapeNamed(t, "Spec 0213 PR only").report())
	before := mustRead(t, implementTaskPath(taskRef.Path, "task_qa"))
	var stdout, stderr bytes.Buffer
	code := runCLIContext(t, t.Context(), []string{"settle", "--spec", implementTestSlug, "--task", "task_qa"}, &stdout, &stderr)
	want := "Settle surface: " + taskRef.Path + "\nroundfix: settle QA Report is ineligible: " + shape.reason + " (report " + rel + ")\n"
	if code != exitRunFailed || stderr.String() != want {
		t.Fatalf("exit=%d stdout=%q stderr=%q, want exit 1 stderr=%q", code, stdout.String(), stderr.String(), want)
	}
	if got := mustRead(t, implementTaskPath(taskRef.Path, "task_qa")); got != before {
		t.Fatal("refusal changed Task")
	}
}

func TestRunSpecDoesNotParkAQualifyingPartialAsEnvironmentOnly(t *testing.T) {
	const pullRequestRows = `
## Results

| ID | Provenance | Status |
| --- | --- | --- |
| PR | Pull Request row | blocked (environment: no open Pull Request) |
| CI | Pull Request row | blocked (environment: no open Pull Request) |
`
	tests := []struct {
		name, content string
		want          bool
	}{
		{"pull request rows only", "---\nverdict: partial\nrows_blocked_finding: 0\nrows_blocked_environment: 2\n---\n" + pullRequestRows, false},
		{"PR and network outside evidence", qaPartialShapeNamed(t, "PR and network").report(), false},
		{"another environment row", qaPartialShapeNamed(t, "another environment beside network").report(), true},
		{"network without outside evidence", qaPartialShapeNamed(t, "network without outside evidence").report(), true},
		{"declared rows only", "---\nverdict: partial\nrows_blocked_finding: 0\nrows_blocked_environment: 0\nrows_blocked_declared: 1\n---\n\n## Results\n\n| ID | Provenance | Status |\n| --- | --- | --- |\n| D1 | Declared acceptance | blocked (declared: unreachable acceptance) |\n", false},
		{"finding beside pull request rows", "---\nverdict: partial\nrows_blocked_finding: 1\nrows_blocked_environment: 2\n---\n" + pullRequestRows + "| F1 | Finding row | blocked (finding: regression) |\n", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			home, repo := newImplementWorkspace(t, []implementSeed{{id: "task_01"}})
			workflow := newItemRecoveryWorkflowForRepository(t, home, repo)
			head := itemRecoveryHead(t, repo)
			const runID = "qa-policy-partial"
			gittest.Run(t, repo, "checkout", "-b", store.RunBranchPrefix+runID)
			dir := filepath.Join(repo, "docs", "specs", implementTestSlug, "qa")
			writeQAReportTestFile(t, filepath.Join(dir, "qa-report-2026-10-05.md"), qaPartialShapeNamed(t, "network without outside evidence").report())
			writeQAReportTestFile(t, filepath.Join(dir, "qa-report-2026-10-06.md"), tt.content)
			gittest.Run(t, repo, "add", "docs/specs")
			gittest.Run(t, repo, "commit", "-m", "docs: Run QA report")
			gittest.Run(t, repo, "checkout", "--detach", head)
			result, err := workflow.runResult(t.Context(), repo, implementTestSlug, roundfixCommandResult{exitCode: exitRunFailed}, "", nil, &store.Run{ID: runID, State: store.StateUnresolved})
			if err != nil {
				t.Fatal(err)
			}
			if result.Outcome != delivery.RunOutcomeUnresolved || result.QAEnvironmentPartial != tt.want {
				t.Fatalf("result=%+v want environment partial=%v", result, tt.want)
			}
		})
	}
}

// The same measured shapes exercise command and Daemon settlement independently.
type qaPartialShape struct {
	name, rows                                   string
	environment, finding, declared, declarations int
	reason                                       string
}

func qaPartialShapes() []qaPartialShape {
	pr := "| " + spec.QANoOpenPullRequestStatus + " | " + spec.QAPullRequestRowSource + " |\n"
	network := "| " + spec.QANetworkDeniedStatusPrefix + "docs.example.com) | Requirement 7; " + spec.QAOutsideEvidenceRowSource + " (published guide) |\n"
	declared := "| blocked (declared) | Requirement 1 |\n| blocked (declared) | Requirement 2 |\n"
	return []qaPartialShape{
		{"Spec 0213 PR only", pr, 1, 0, 0, 0, ""},
		{"Spec 0220 PR only", pr + pr, 2, 0, 0, 0, ""},
		{"Oraculum PR only", pr, 1, 0, 0, 0, ""},
		{"Spec 0227 annotated PR", "| " + spec.QANoOpenPullRequestStatus + " | Pull Request row (Requirement 10 constrains execution) |\n", 1, 0, 0, 0, ""},
		{"Fluxus covered declarations and PR", pr + declared, 1, 0, 2, 2, ""},
		{"network only", network, 1, 0, 0, 0, ""},
		{"PR and network", pr + network, 2, 0, 0, 0, ""},
		{"all three kinds", pr + network + declared, 2, 0, 2, 2, ""},
		{"declared only", declared, 0, 0, 2, 2, ""},
		{"network without outside evidence", "| " + spec.QANetworkDeniedStatusPrefix + "docs.example.com) | Requirement 8 |\n", 1, 0, 0, 0, "rows_blocked_environment is 1; expected 0"},
		{"empty network host", "| " + spec.QANetworkDeniedStatusPrefix + ") | outside-evidence row |\n", 1, 0, 0, 0, "rows_blocked_environment is 1; expected 0"},
		{"another environment beside network", network + "| blocked (environment: credentials unavailable) | Requirement 8 |\n", 2, 0, 0, 0, "rows_blocked_environment is 2, 1 outside the pre-PR Pull Request row and network-denied outside-evidence rows; expected 0 outside them"},
		{"PR and skipped", pr + "| skipped | Requirement 9 |\n", 1, 0, 0, 0, "partial records 1 skipped row(s); a qualifying partial records none"},
		{"finding", pr + "| blocked (finding) | Requirement 10 |\n", 1, 1, 0, 0, "rows_blocked_finding is 1; expected 0"},
		{"uncovered declaration", pr + declared, 1, 0, 2, 1, "rows_blocked_declared is 2, but Spec declares 1 unreachable acceptance; shortfall is 1"},
		{"no unmet row", "| pass | Requirement 1 |\n", 0, 0, 0, 0, `newest QA Report verdict is "partial"; expected "pass"`},
	}
}

func (shape qaPartialShape) report() string {
	return fmt.Sprintf("---\nverdict: partial\nrows_blocked_environment: %d\nrows_blocked_finding: %d\nrows_blocked_declared: %d\n---\n\n## Results\n\n| Status | Provenance |\n| --- | --- |\n%s", shape.environment, shape.finding, shape.declared, shape.rows)
}

func (shape qaPartialShape) unreachable() string {
	if shape.declarations == 0 {
		return ""
	}
	text := "\n## Unreachable Acceptance\n\n"
	for i := 0; i < shape.declarations; i++ {
		text += fmt.Sprintf("- criterion: acceptance criterion %d\n  reason: unavailable environment\n  satisfied-by: task_01\n", i+1)
	}
	return text
}

func qaPartialShapeNamed(t *testing.T, name string) qaPartialShape {
	t.Helper()
	for _, shape := range qaPartialShapes() {
		if shape.name == name {
			return shape
		}
	}
	t.Fatalf("unknown partial shape %q", name)
	return qaPartialShape{}
}
