// Suite: one QA partial eligibility policy.
// Invariant: only exempt environment rows and covered declarations qualify.
// Boundary IN: ReadQAReportFile and QAReportEligibility on temporary Specs.
// Boundary OUT: callers and mechanical frontmatter count validation.
package spec

import (
	"fmt"
	"path/filepath"
	"testing"
)

func TestAPartialWhoseOnlyUnmetRowsArePullRequestRowsQualifies(t *testing.T) {
	t.Parallel()
	for _, name := range []string{"Spec 0213", "Spec 0220", "Oraculum"} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			dir := t.TempDir()
			report := readPartialPolicyReport(t, dir, 1, 0, 0, partialPolicyTable(partialPolicyPRRow()))
			assertPartialPolicyEligibility(t, dir, report, "")
		})
	}
}

func TestAPullRequestRowWithANoteInItsProvenanceIsRecognized(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	report := readPartialPolicyReport(t, dir, 1, 0, 0, partialPolicyTable(fmt.Sprintf("| %s | %s (Requirement 10 constrains execution) |\n", QANoOpenPullRequestStatus, QAPullRequestRowSource)))
	if report.RowsBlockedPrePullRequest != 1 {
		t.Fatalf("PR count = %d, want 1", report.RowsBlockedPrePullRequest)
	}
	assertPartialPolicyEligibility(t, dir, report, "")
}

func TestANetworkDeniedOutsideEvidenceRowNeverDecidesAQualifyingPartial(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name                  string
		pr, network, declared int
	}{
		{"network only", 0, 1, 0}, {"PR and network", 1, 1, 0},
		{"Fluxus PR and two declarations", 1, 0, 2},
		{"all three kinds", 1, 1, 2}, {"network and declared", 0, 1, 2},
		{"declared only", 0, 0, 2},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			dir := t.TempDir()
			rows := ""
			if tt.pr > 0 {
				rows += partialPolicyPRRow()
			}
			if tt.network > 0 {
				rows += partialPolicyNetworkRow()
			}
			if tt.declared > 0 {
				writeFile(t, filepath.Join(dir, "_prd.md"), "# Test Spec\n\n## Unreachable Acceptance\n\n- criterion: first declared acceptance\n  reason: unavailable environment\n  satisfied-by: task_01\n- criterion: second declared acceptance\n  reason: unavailable environment\n  satisfied-by: task_01\n")
				rows += "| blocked (declared) | Requirement 1 |\n| blocked (declared) | Requirement 2 |\n"
			}
			report := readPartialPolicyReport(t, dir, tt.pr+tt.network, 0, tt.declared, partialPolicyTable(rows))
			if report.RowsBlockedNetworkDenied != tt.network {
				t.Fatalf("network count = %d, want %d", report.RowsBlockedNetworkDenied, tt.network)
			}
			assertPartialPolicyEligibility(t, dir, report, "")
		})
	}
}

func TestANetworkDeniedRowWithoutOutsideEvidenceProvenanceStillRefuses(t *testing.T) {
	t.Parallel()
	tests := []struct{ name, status, provenance string }{
		{"Requirement 8", QANetworkDeniedStatusPrefix + "docs.example.com)", "Requirement 8"},
		{"empty host", QANetworkDeniedStatusPrefix + ")", QAOutsideEvidenceRowSource},
		{"blank host", QANetworkDeniedStatusPrefix + "  )", QAOutsideEvidenceRowSource},
		{"missing closing parenthesis", QANetworkDeniedStatusPrefix + "docs.example.com", QAOutsideEvidenceRowSource},
		{"wrong source boundary", QANetworkDeniedStatusPrefix + "docs.example.com)", QAOutsideEvidenceRowSource + "s"},
		{"embedded source", QANetworkDeniedStatusPrefix + "docs.example.com)", "not an " + QAOutsideEvidenceRowSource},
		{"case sensitive source", QANetworkDeniedStatusPrefix + "docs.example.com)", "Outside-evidence row"},
		{"case sensitive status", "Blocked (environment: network denied: docs.example.com)", QAOutsideEvidenceRowSource},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			dir := t.TempDir()
			report := readPartialPolicyReport(t, dir, 1, 0, 0, partialPolicyTable(fmt.Sprintf("| %s | %s |\n", tt.status, tt.provenance)))
			if report.RowsBlockedNetworkDenied != 0 {
				t.Fatalf("network count = %d, want 0", report.RowsBlockedNetworkDenied)
			}
			assertPartialPolicyEligibility(t, dir, report, "rows_blocked_environment is 1; expected 0")
		})
	}
}

func TestAPartialWithAnotherEnvironmentRowBesideANetworkDeniedRowNamesBoth(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	report := readPartialPolicyReport(t, dir, 2, 0, 0, partialPolicyTable(partialPolicyNetworkRow()+"| blocked (environment: credentials unavailable) | Requirement 8 |\n"))
	assertPartialPolicyEligibility(t, dir, report, "rows_blocked_environment is 2, 1 outside the pre-PR Pull Request row and network-denied outside-evidence rows; expected 0 outside them")
}

func TestAPartialWithASkippedRowNeverQualifies(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name                           string
		environment, finding, declared int
		rows, want                     string
	}{
		{"PR and skipped", 1, 0, 0, partialPolicyPRRow() + "| skipped | Requirement 9 |\n", "partial records 1 skipped row(s); a qualifying partial records none"},
		{"skipped alone", 0, 0, 0, "| SKIPPED | Requirement 9 |\n", "partial records 1 skipped row(s); a qualifying partial records none"},
		{"skipped precedes declarations", 1, 0, 1, partialPolicyPRRow() + "| skipped | Requirement 9 |\n| blocked (declared) | Requirement 1 |\n", "partial records 1 skipped row(s); a qualifying partial records none"},
		{"finding precedes environment and skipped", 2, 1, 0, partialPolicyPRRow() + "| blocked (environment: credentials unavailable) | Requirement 8 |\n| skipped | Requirement 9 |\n| blocked (finding) | Requirement 10 |\n", "rows_blocked_finding is 1; expected 0"},
		{"environment precedes skipped", 1, 0, 0, "| blocked (environment: credentials unavailable) | Requirement 8 |\n| skipped | Requirement 9 |\n", "rows_blocked_environment is 1; expected 0"},
		{"uncovered declaration", 1, 0, 2, partialPolicyPRRow() + "| blocked (declared) | Requirement 1 |\n| blocked (declared) | Requirement 2 |\n", "rows_blocked_declared is 2, but Spec declares 1 unreachable acceptance; shortfall is 1"},
		{"no unmet row", 0, 0, 0, "| pass | Requirement 1 |\n", `newest QA Report verdict is "partial"; expected "pass"`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			dir := t.TempDir()
			if tt.name == "uncovered declaration" {
				writePrePullRequestDeclaration(t, dir)
			}
			report := readPartialPolicyReport(t, dir, tt.environment, tt.finding, tt.declared, partialPolicyTable(tt.rows))
			assertPartialPolicyEligibility(t, dir, report, tt.want)
		})
	}
}

func TestEnvironmentRowsNeedingOverrideCountsOnlyNonExemptRows(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name                          string
		environment                   int
		body                          string
		pr, network, skipped, outside int
	}{
		{"no environment rows", 0, partialPolicyTable("| pass | Requirement 1 |\n"), 0, 0, 0, 0},
		{"PR only", 1, partialPolicyTable(partialPolicyPRRow()), 1, 0, 0, 0},
		{"network only", 1, partialPolicyTable(partialPolicyNetworkRow()), 0, 1, 0, 0},
		{"both exemptions and another environment", 3, partialPolicyTable(partialPolicyPRRow() + partialPolicyNetworkRow() + "| blocked (environment: no browser) | Requirement 8 |\n"), 1, 1, 0, 1},
		{"exemptions exceed frontmatter", 1, partialPolicyTable(partialPolicyPRRow() + partialPolicyNetworkRow()), 1, 1, 0, 0},
		{"zero frontmatter clamps exemptions", 0, partialPolicyTable(partialPolicyPRRow() + partialPolicyNetworkRow()), 1, 1, 0, 0},
		{"no provenance", 2, "## Results\n\n| Status | Evidence |\n| --- | --- |\n| " + QANoOpenPullRequestStatus + " | none |\n| " + QANetworkDeniedStatusPrefix + "docs.example.com) | none |\n| SkIpPeD | none |\n", 0, 0, 1, 2},
		{"no status", 1, "## Results\n\n| Evidence | Provenance |\n| --- | --- |\n| " + QANetworkDeniedStatusPrefix + "docs.example.com) | " + QAOutsideEvidenceRowSource + " |\n| skipped | Requirement 9 |\n", 0, 0, 0, 1},
		{"outside Results", 1, "## Evidence\n\n" + partialPolicyTable(partialPolicyPRRow() + partialPolicyNetworkRow() + "| skipped | Requirement 9 |\n")[len("## Results\n\n"):] + "\n## Results\n\n| Status | Provenance |\n| --- | --- |\n| pass | Requirement 1 |\n", 0, 0, 0, 1},
		{"fenced tables", 1, "## Results\n\n~~~markdown\n" + partialPolicyTable(partialPolicyPRRow()+partialPolicyNetworkRow()+"| skipped | Requirement 9 |\n") + "~~~\n\n```markdown\n" + partialPolicyTable(partialPolicyPRRow()+partialPolicyNetworkRow()+"| skipped | Requirement 9 |\n") + "```\n\n| Status | Provenance |\n| --- | --- |\n| pass | Requirement 1 |\n", 0, 0, 0, 1},
		{"multiple tables and deeper headings", 2, partialPolicyTable(partialPolicyPRRow()) + "\n### Outside evidence\n\n| Provenance | STATUS |\n| --- | --- |\n| " + QAOutsideEvidenceRowSource + ": published guide | " + QANetworkDeniedStatusPrefix + "docs.example.com) |\n\n## Notes\n\n| Status | Provenance |\n| --- | --- |\n| skipped | Requirement 9 |\n", 1, 1, 0, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			report := readPartialPolicyReport(t, t.TempDir(), tt.environment, 0, 0, tt.body)
			if report.RowsBlockedPrePullRequest != tt.pr || report.RowsBlockedNetworkDenied != tt.network || report.RowsSkipped != tt.skipped {
				t.Fatalf("derived counts PR/network/skipped = %d/%d/%d, want %d/%d/%d", report.RowsBlockedPrePullRequest, report.RowsBlockedNetworkDenied, report.RowsSkipped, tt.pr, tt.network, tt.skipped)
			}
			if got := report.EnvironmentRowsNeedingOverride(); got != tt.outside {
				t.Fatalf("override count = %d, want %d", got, tt.outside)
			}
		})
	}
}

func partialPolicyNetworkRow() string {
	return fmt.Sprintf("| %sdocs.example.com) | Requirement 7; %s (published guide) |\n", QANetworkDeniedStatusPrefix, QAOutsideEvidenceRowSource)
}

func partialPolicyPRRow() string {
	return fmt.Sprintf("| %s | %s |\n", QANoOpenPullRequestStatus, QAPullRequestRowSource)
}

func partialPolicyTable(rows string) string {
	return "## Results\n\n| Status | Provenance |\n| --- | --- |\n" + rows
}

func readPartialPolicyReport(t *testing.T, dir string, environment, finding, declared int, body string) QAReport {
	t.Helper()
	path := filepath.Join(dir, "qa", "qa-report-2026-10-06.md")
	writeFile(t, path, fmt.Sprintf("---\nverdict: partial\nrows_blocked_environment: %d\nrows_blocked_finding: %d\nrows_blocked_declared: %d\n---\n\n%s", environment, finding, declared, body))
	report, err := ReadQAReportFile(path)
	if err != nil {
		t.Fatalf("ReadQAReportFile: %v", err)
	}
	return report
}

func assertPartialPolicyEligibility(t *testing.T, dir string, report QAReport, want string) {
	t.Helper()
	err := QAReportEligibility(dir, report)
	if want == "" {
		if err != nil {
			t.Fatalf("QAReportEligibility: %v", err)
		}
	} else if err == nil || err.Error() != want {
		t.Fatalf("error = %v, want %q", err, want)
	}
}
