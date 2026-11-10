// Suite: the shared QA Results reader and its settlement consumers.
// Invariant: every Results table line belongs to a header-led table or yields a precise refusal.
// Boundary IN: ReadQAResults, QAReportEligibility, archive and evidence checks on temporary reports.
// Boundary OUT: CLI diagnostics and mechanical findings, owned by their package tests.
package spec

import (
	"errors"
	"fmt"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

const pantheonQAResults = "---\nverdict: partial\nrows_blocked_environment: 1\nrows_blocked_finding: 0\nrows_blocked_declared: 1\n---\n\n# QA Report\n\n## Results\n\n| #   | Status | Provenance |\n| --- | ------ | ---------- |\n| Q01 | pass | PRD Goal 1 |\n\n| Q02 | pass | Task 01 Acceptance Criterion 1 |\n| Q03 | blocked (environment: no open Pull Request) | Pull Request row |\n| Q04 | blocked (declared: the probe reads the live service) | PRD Unreachable Acceptance: the probe reads the live service |\n\n## Mechanical skips\n\nNone.\n"

const pantheonQADeclaration = "## Unreachable Acceptance\n\n- criterion: the probe reads the live service\n  reason: needs a deploy\n  satisfied-by: the operator after release\n"

func TestQAResultsRefusesARowAfterABlankLine(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	path := filepath.Join(dir, "qa/qa-report-2026-11-10.md")
	writeFile(t, filepath.Join(dir, "_prd.md"), pantheonQADeclaration)
	writeFile(t, path, pantheonQAResults)
	report, err := ReadQAReportFile(path)
	if err != nil {
		t.Fatal(err)
	}
	want := fmt.Sprintf("QA Report %s: Results row at line 16 is outside the Results table that starts at line 12; remove the blank line at %s:15 so the row continues the table", path, path)
	for _, content := range []string{pantheonQAResults, strings.ReplaceAll(pantheonQAResults, "\n", "\r\n")} {
		results := ReadQAResults(path, []byte(content))
		wantBreak := QAResultsBreak{Path: path, Line: 16, TableLine: 12, LastRowLine: 14, BlankFrom: 15, BlankTo: 15}
		if results.Break == nil || *results.Break != wantBreak || results.Break.Error() != want {
			t.Fatalf("break = %+v, want %+v with error %q", results.Break, wantBreak, want)
		}
	}
	if report.ResultsBreak == nil || report.ResultsBreak.Error() != want {
		t.Fatalf("report break = %v, want %q", report.ResultsBreak, want)
	}
	if err := QAReportEligibility(dir, report); err == nil || err.Error() != want || strings.Contains(err.Error(), "rows_blocked_environment") {
		t.Fatalf("eligibility = %v, want %q", err, want)
	}
	writeFile(t, path, strings.Replace(pantheonQAResults, "\n\n| Q02", "\n| Q02", 1))
	report, err = ReadQAReportFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if report.RowsBlockedPrePullRequest != 1 || len(report.DeclaredRows) != 1 || report.DeclaredRows[0].ID != "Q04" {
		t.Fatalf("repaired rows = %+v, want one pre-PR row and declared Q04", report)
	}
	if report.ResultsBreak != nil {
		t.Fatalf("repaired break = %v", report.ResultsBreak)
	}
	if err := QAReportEligibility(dir, report); err != nil {
		t.Fatalf("repaired eligibility = %v", err)
	}
}

func TestQAResultsNamesEveryBlankLineOrTheRowsPlace(t *testing.T) {
	t.Parallel()
	const path = "qa/report.md"
	tests := []struct {
		name, content, detail, repair string
		line, table, last, from, to   int
	}{
		{
			name: "two blank lines", content: strings.Replace(pantheonQAResults, "\n\n| Q02", "\n\n\n| Q02", 1),
			line: 17, table: 12, last: 14, from: 15, to: 16,
			detail: "Results row at line 17 is outside the Results table that starts at line 12",
			repair: "remove the blank lines at qa/report.md:15-16 so the row continues the table",
		},
		{
			name: "prose separates the row", content: strings.Replace(pantheonQAResults, "\n\n| Q02", "\nSome prose.\n| Q02", 1),
			line: 16, table: 12, last: 14,
			detail: "Results row at line 16 is outside the Results table that starts at line 12",
			repair: "move the row directly below the last row of the table, at qa/report.md:14",
		},
		{
			name: "header has no separator", content: strings.Replace(pantheonQAResults, "| --- | ------ | ---------- |\n", "", 1),
			line:   12,
			detail: "Results row at line 12 is outside every Results table",
			repair: "add a separator line directly below the table header at qa/report.md:12",
		},
		{
			name: "deeper heading does not reset the table", content: strings.Replace(pantheonQAResults, "\n\n| Q02", "\n### Outside evidence\n| Q02", 1),
			line: 16, table: 12, last: 14,
			detail: "Results row at line 16 is outside the Results table that starts at line 12",
			repair: "move the row directly below the last row of the table, at qa/report.md:14",
		},
		{
			name: "empty table uses its header for the last row", content: "## Results\n| # | Status |\n| - | --- |\nSome prose.\n| Q01 | pass |\n",
			line: 5, table: 2, last: 2,
			detail: "Results row at line 5 is outside the Results table that starts at line 2",
			repair: "move the row directly below the last row of the table, at qa/report.md:2",
		},
		{
			name: "blank after an empty table excludes the separator", content: "## Results\n| # | Status |\n| - | --- |\n\n| Q01 | pass |\n",
			line: 5, table: 2, last: 2, from: 4, to: 4,
			detail: "Results row at line 5 is outside the Results table that starts at line 2",
			repair: "remove the blank line at qa/report.md:4 so the row continues the table",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			results := ReadQAResults(path, []byte(tt.content))
			want := QAResultsBreak{Path: path, Line: tt.line, TableLine: tt.table, LastRowLine: tt.last, BlankFrom: tt.from, BlankTo: tt.to}
			if results.Break == nil || *results.Break != want {
				t.Fatalf("break = %+v, want %+v", results.Break, want)
			}
			if results.Break.Detail() != tt.detail || results.Break.Repair() != tt.repair || results.Break.Error() != "QA Report "+path+": "+tt.detail+"; "+tt.repair {
				t.Fatalf("break text = %q, want detail %q and repair %q", results.Break.Error(), tt.detail, tt.repair)
			}
		})
	}
}

func TestQAResultsReadsEveryStatusTable(t *testing.T) {
	t.Parallel()
	const first = "| # | Status | Provenance |\n| - | --- | --- |\n| Q01 | pass | Requirement 1 |\n"
	content := "---\nverdict: partial\nrows_blocked_environment: 3\nrows_blocked_declared: 1\n---\n\n" +
		"```markdown\n## Results\n| orphan | ignored |\n```\n\n## rEsUlTs ##\n\n" + first +
		"\n| # | Status | Provenance |\n| - | --- | --- |\n| Q02 | " + QANoOpenPullRequestStatus + " | " + QAPullRequestRowSource + " |\n" +
		"\n### Outside evidence\n\n| # | Provenance | STATUS | Evidence |\n| - | --- | --- | --- |\n" +
		"| Q03 | outside-evidence row | " + QANetworkDeniedStatusPrefix + "docs.example.com) | denied |\n" +
		"| Q04 | outside-evidence row | " + QANetworkProhibitedStatus + " | prohibited |\n" +
		"| Q05 | Requirement 5 | blocked (declared: the probe reads the live service) | declaration |\n" +
		"| Q06 | Requirement 6 | SkIpPeD | skipped |\n" +
		"\n| Case | Observation |\n| --- | --- |\n| comparison | skipped |\n" +
		"\n```markdown\n| orphan | skipped |\n```\n\n~~~markdown\n" + first + "~~~\n" +
		"\n## Notes\n\n| orphan | skipped |\n\n## Results\n\n" + first
	results := ReadQAResults("qa/report.md", []byte(content))
	if results.Break != nil || len(results.Tables) != 4 || len(results.StatusTables()) != 3 {
		t.Fatalf("results = %+v, want four tables and three Status tables", results)
	}
	var ids []string
	lines := strings.Split(content, "\n")
	for _, table := range results.StatusTables() {
		if !reflect.DeepEqual(table.Header, markdownTableCells(lines[table.HeaderLine-1])) {
			t.Fatalf("header at line %d = %q", table.HeaderLine, table.Header)
		}
		for _, row := range table.Rows {
			ids = append(ids, row.Cells[0])
			if !strings.HasPrefix(lines[row.Line-1], "| "+row.Cells[0]+" |") {
				t.Fatalf("row %+v does not name its source line", row)
			}
		}
	}
	if !reflect.DeepEqual(ids, []string{"Q01", "Q02", "Q03", "Q04", "Q05", "Q06"}) {
		t.Fatalf("row IDs = %q", ids)
	}
	dir := t.TempDir()
	path := filepath.Join(dir, "qa/report.md")
	writeFile(t, path, content)
	report, err := ReadQAReportFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if report.ResultsBreak != nil || report.RowsBlockedPrePullRequest != 1 || report.RowsBlockedNetworkDenied != 1 || report.RowsBlockedNetworkProhibited != 1 || report.RowsSkipped != 1 ||
		!reflect.DeepEqual(report.DeclaredRows, []QADeclaredRow{{ID: "Q05", Criterion: "the probe reads the live service"}}) {
		t.Fatalf("derived report = %+v", report)
	}
	for _, ending := range []string{"# Notes", "## Notes"} {
		t.Run(ending, func(t *testing.T) {
			t.Parallel()
			results := ReadQAResults("", []byte("## Results\n"+first+"\n"+ending+"\n| orphan | ignored |\n"))
			if results.Break != nil || len(results.StatusTables()) != 1 {
				t.Fatalf("section bounds = %+v", results)
			}
		})
	}
	t.Run("carried exemptions use each table's own Status column", func(t *testing.T) {
		t.Parallel()
		const carried = "carried (established by: earlier.md; head: abc)"
		content := "## Results\n| # | Status |\n| - | --- |\n| C01 | " + carried + " |\n\n" +
			"### Outside evidence\n| # | Provenance | Status |\n| - | --- | --- |\n| C02 | measured | " + carried + " |\n\n" +
			"| # | Observation |\n| - | --- |\n| C03 | " + carried + " |\n\n" +
			"## Evidence\n### C01\ncarried first\n### C02\ncarried second\n### C03\nmeasured third\n"
		exempt := borrowedEvidenceCarriedLines([]byte(content))
		want := map[int]bool{4: true, 9: true, 16: true, 17: true, 18: true, 19: true}
		if !reflect.DeepEqual(exempt, want) {
			t.Fatalf("carried evidence lines = %v, want %v", exempt, want)
		}
	})
}

func TestQAResultsBreakStopsEveryCountingReader(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	path := filepath.Join(dir, "qa/qa-report-2026-11-10.md")
	writeFile(t, filepath.Join(dir, "_prd.md"), pantheonQADeclaration)
	writeFile(t, path, pantheonQAResults)
	report, err := ReadQAReportFile(path)
	if err != nil || report.ResultsBreak == nil {
		t.Fatalf("read = %+v, %v; want a readable verdict and a break", report, err)
	}
	if report.EnvironmentRowsNeedingOverride() != 0 || report.RowsBlockedPrePullRequest != 0 || len(report.DeclaredRows) != 0 {
		t.Fatalf("broken report counted rows: %+v", report)
	}
	assertBreak := func(t *testing.T, err error) {
		t.Helper()
		var brk QAResultsBreak
		if !errors.As(err, &brk) || brk != *report.ResultsBreak {
			t.Fatalf("error = %v, want break %+v", err, report.ResultsBreak)
		}
	}
	for _, verdict := range []string{VerdictPass, VerdictFail, VerdictPartial, VerdictPending, ""} {
		t.Run(verdict, func(t *testing.T) {
			t.Parallel()
			candidate := report
			candidate.Verdict, candidate.Hollow, candidate.RowsBlockedFinding = verdict, true, 1
			candidate.RowsSkipped = 1
			assertBreak(t, QAReportEligibility(dir, candidate))
		})
	}
	actions, err := archiveUnprovenActions(dir, report)
	if actions != nil {
		t.Fatalf("archive actions = %q, want nil", actions)
	}
	assertBreak(t, err)
	options := BorrowedEvidenceOptions{RunBranchPrefix: "roundfix/run-"}
	items, err := FindBorrowedEvidence(t.Context(), path, options)
	if items != nil {
		t.Fatalf("borrowed evidence = %+v, want nil", items)
	}
	assertBreak(t, err)
	assertBreak(t, RequireOwnRunEvidence(t.Context(), path, options))
	brokenCarried := strings.Replace(pantheonQAResults, "| Q02 | pass |", "| Q02 | carried (established by: earlier.md; head: abc) |", 1) + "\n### Q02\n\nBorrowed evidence.\n"
	if exempt := borrowedEvidenceCarriedLines([]byte(brokenCarried)); len(exempt) != 0 {
		t.Fatalf("orphan row exempted evidence lines: %v", exempt)
	}
}
