// Suite: shared QA Results reading in the mechanical stage.
// Invariant: an orphan table line refuses shape checks before counts, and every Status table is read.
// Boundary IN: RunMechanicalStage over temporary Git repositories and QA Reports.
// Boundary OUT: eligibility and CLI settlement, owned by internal/spec and internal/cli.
package speccheck_test

import (
	"strings"
	"testing"

	"roundfix/internal/speccheck"
)

const mechanicalPantheonResults = "---\nverdict: partial\nrows_blocked_environment: 1\nrows_blocked_finding: 0\nrows_blocked_declared: 1\n---\n\n# QA Report\n\n## Results\n\n| #   | Status | Provenance |\n| --- | ------ | ---------- |\n| Q01 | pass | PRD Goal 1 |\n\n| Q02 | pass | Task 01 Acceptance Criterion 1 |\n| Q03 | blocked (environment: no open Pull Request) | Pull Request row |\n| Q04 | blocked (declared: the probe reads the live service) | PRD Unreachable Acceptance: the probe reads the live service |\n\n## Mechanical skips\n\nNone.\n"

func TestMechanicalReportShapeNamesARowAfterABlankLine(t *testing.T) {
	t.Parallel()
	for _, malformedFrontmatter := range []bool{false, true} {
		name := "valid frontmatter"
		if malformedFrontmatter {
			name = "frontmatter parse finding precedes the break"
		}
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			root := newMechanicalGitRepo(t)
			const path = "docs/specs/mechanical/qa/report.md"
			content := mechanicalPantheonResults
			if malformedFrontmatter {
				content = strings.Replace(content, "rows_blocked_environment: 1", "rows_blocked_environment: []", 1)
			}
			writeMechanicalFile(t, root, path, content)
			writeMechanicalFile(t, root, "docs/specs/mechanical/_prd.md", "## Unreachable Acceptance\n\n- criterion: the probe reads the live service\n  reason: needs a deploy\n  satisfied-by: the operator after release\n")
			result := runMechanical(t, speccheck.MechanicalRequest{RepoRoot: root, ReportPath: path})
			findings := mechanicalFindingsWithCode(result, speccheck.CodeMechanicalReportShape)
			wantCount := 1
			if malformedFrontmatter {
				wantCount++
			}
			if len(findings) != wantCount {
				t.Fatalf("shape findings = %#v, want %d", findings, wantCount)
			}
			var breaks int
			for _, finding := range findings {
				if strings.Contains(finding.Detail, "Results table contains") {
					t.Fatalf("broken report reconciled counts: %#v", finding)
				}
				if strings.HasPrefix(finding.Detail, "report structure cannot be parsed:") && malformedFrontmatter {
					continue
				}
				breaks++
				if finding.File != path || finding.Line != 16 || finding.RowHint != "" ||
					finding.Detail != "Results row at line 16 is outside the Results table that starts at line 12" ||
					finding.Fix != "Remove the blank line at docs/specs/mechanical/qa/report.md:15 so the row continues the table." {
					t.Fatalf("break finding = %#v", finding)
				}
			}
			if breaks != 1 {
				t.Fatalf("break findings = %d, want 1", breaks)
			}
		})
	}
}

func TestMechanicalStageReadsEveryStatusTable(t *testing.T) {
	t.Parallel()
	root := newMechanicalGitRepo(t)
	const path = "docs/specs/mechanical/qa/report.md"
	content := "---\nverdict: pass\nrows_blocked_environment: 0\nrows_blocked_finding: 0\nrows_blocked_declared: 0\n---\n\n## Results\n\n| # | Status | Provenance |\n| - | --- | --- |\n| R01 | pass | Requirement 1 |\n\n| # | Status | Provenance |\n| - | --- | --- |\n| R02 | pending | Requirement 2 |\n"
	writeMechanicalFile(t, root, path, content)
	result := runMechanical(t, speccheck.MechanicalRequest{RepoRoot: root, ReportPath: path})
	findings := mechanicalFindingsWithCode(result, speccheck.CodeMechanicalReportShape)
	if len(findings) != 1 || findings[0].Detail != "row R02 remains pending instead of carrying a terminal status" || findings[0].Line != 16 {
		t.Fatalf("shape findings = %#v, want second table's pending row at line 16", findings)
	}
}
