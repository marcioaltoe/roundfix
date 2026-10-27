// Suite: authorization subjects in the shared QA partial policy.
// Invariant: bounded subjects with outside-evidence provenance count once.
// Boundary: temporary reports and Specs; no network or archived-file reads.
package spec

import (
	"fmt"
	"strings"
	"testing"
)

func TestANetworkProhibitedRowNamesWhatTheAuthorizationForbids(t *testing.T) {
	t.Parallel()
	const refused = "rows_blocked_environment is 1; expected 0"
	status := func(subject string) string {
		return "blocked (environment: " + subject + " prohibited by Spec authorization)"
	}
	tests := []struct {
		name, status, provenance, extra, reason string
		environment, prohibited                 int
	}{
		{"accepts Spec 0309 R07 another repository", status("another repository"), QAOutsideEvidenceRowSource, "", "", 1, 1},
		{"accepts Spec 0317 QA-08 upstream checkout and network", status("upstream checkout and network"), QAOutsideEvidenceRowSource, "", "", 1, 1},
		{"accepts the canonical network form", QANetworkProhibitedStatus, QAOutsideEvidenceRowSource, "", "", 1, 1},
		{"accepts a subject of 80 multibyte characters", status("  " + strings.Repeat("é", QAProhibitedSubjectMaxRunes) + "  "), QAOutsideEvidenceRowSource, "", "", 1, 1},
		{"refuses an empty subject", status(""), QAOutsideEvidenceRowSource, "", refused, 1, 0},
		{"refuses a blank subject", status("   "), QAOutsideEvidenceRowSource, "", refused, 1, 0},
		{"refuses a subject of 81 characters", status(strings.Repeat("a", QAProhibitedSubjectMaxRunes+1)), QAOutsideEvidenceRowSource, "", refused, 1, 0},
		{"refuses a subject with a parenthesis", status("another (upstream) repository"), QAOutsideEvidenceRowSource, "", refused, 1, 0},
		{"refuses a changed suffix case", "blocked (environment: network prohibited by spec authorization)", QAOutsideEvidenceRowSource, "", refused, 1, 0},
		{"refuses a row without outside evidence", status("another repository"), "Requirement 8", "", refused, 1, 0},
		{"a row in both forms counts once", status("network denied: docs.example.com"), QAOutsideEvidenceRowSource, "| blocked (environment: credentials unavailable) | Requirement 8 |\n", "rows_blocked_environment is 2, 1 outside the pre-PR Pull Request row and the network-denied and network-prohibited outside-evidence rows; expected 0 outside them", 2, 1},
		{"a skipped row still refuses", status("another repository"), QAOutsideEvidenceRowSource, "| skipped | Requirement 9 |\n", "partial records 1 skipped row(s); a qualifying partial records none", 1, 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			dir := t.TempDir()
			rows := fmt.Sprintf("| %s | %s |\n", tt.status, tt.provenance) + tt.extra
			report := readPartialPolicyReport(t, dir, tt.environment, 0, 0, partialPolicyTable(rows))
			if report.RowsBlockedNetworkProhibited != tt.prohibited || report.RowsBlockedNetworkDenied != 0 {
				t.Fatalf("prohibited/denied counts = %d/%d, want %d/0", report.RowsBlockedNetworkProhibited, report.RowsBlockedNetworkDenied, tt.prohibited)
			}
			assertPartialPolicyEligibility(t, dir, report, tt.reason)
		})
	}
}

func TestTheArchivedAuthorizationRowsQualify(t *testing.T) {
	t.Parallel()
	// Results tables copied from the archived reports at their recorded revisions.
	// Keep these fixtures local: archived Specs may be removed from the repository.
	tests := []struct {
		name                    string
		environment, prohibited int
		body                    string
	}{
		// 078869bf
		{"Spec 0258", 5, 4, `## Results

| # | Status | Provenance |
| - | --- | --- |
| 01 | pass | task_05 Requirement 1; repository Verification |
| 10 | pass | task_05 Requirement 10 |
| 11 | pass | outside-evidence row (adopted Backlog Entry) |
| 12 | blocked (environment: network prohibited by Spec authorization) | outside-evidence row (Robinson 2006) |
| 13 | blocked (environment: network prohibited by Spec authorization) | outside-evidence row (Juergens 2011) |
| 14 | blocked (environment: network prohibited by Spec authorization) | outside-evidence row (Zimmermann 2022) |
| 15 | blocked (environment: network prohibited by Spec authorization) | outside-evidence row (Gotel and Finkelstein 1994) |
| 18 | blocked (declared: fewer QA failures in Specs authored after release) | PRD Unreachable Acceptance: fewer QA failures in Specs authored after release |
| 19 | blocked (environment: no open Pull Request) | Pull Request row |
`},
		// 06a9b7ce
		{"Spec 0309", 4, 3, `## Results

| # | Status | Provenance |
| - | --- | --- |
| R01 | pass | repository Verification |
| R02 | pass | Requirement 2 |
| R03 | pass | Requirement 3 |
| R04 | pass | Requirement 2; Requirement 3 |
| R05 | pass | Requirement 3 |
| R06 | blocked (environment: network prohibited by Spec authorization) | outside-evidence row |
| R07 | blocked (environment: another repository prohibited by Spec authorization) | outside-evidence row |
| R08 | blocked (environment: network prohibited by Spec authorization) | outside-evidence row |
| R09 | blocked (declared unreachable) | PRD Unreachable Acceptance: a real Python repository adopts the python profile with roundfix baseline update, restores the 52 skills and runs its own make verify under the new guide |
| R11 | pass | Requirement 9 |
| R12 | blocked (environment: no open Pull Request) | Pull Request row |
`},
		// ef1a4bab
		{"Spec 0317", 3, 2, `## Results

| # | Status | Provenance |
| - | --- | --- |
| QA-01 | pass | Requirement 1; repository Verification |
| QA-02 | pass | Requirement 2 |
| QA-03 | pass | Requirement 3 |
| QA-04 | pass | Requirement 4 |
| QA-05 | pass | Requirement 5 |
| QA-06 | pass | Requirement 6 |
| QA-07 | pass | Requirement 7 |
| QA-08 | blocked (environment: upstream checkout and network prohibited by Spec authorization) | outside-evidence row |
| QA-09 | blocked (environment: network prohibited by Spec authorization) | outside-evidence row |
| QA-10 | blocked (declared unreachable) | PRD Unreachable Acceptance: Success Metric 6 |
| QA-11 | pass | Requirement 10 |
| QA-12 | pass | Requirement 11 |
| QA-13 | blocked (environment: no open Pull Request) | Pull Request row |
| QA-14 | pass | Requirement 13 |
`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			dir := t.TempDir()
			writePrePullRequestDeclaration(t, dir)
			report := readPartialPolicyReport(t, dir, tt.environment, 0, 1, tt.body)
			if report.RowsBlockedNetworkProhibited != tt.prohibited || report.RowsBlockedPrePullRequest != 1 || report.RowsBlockedNetworkDenied != 0 || report.EnvironmentRowsNeedingOverride() != 0 {
				t.Fatalf("archived counts = %+v; want %d prohibited, one PR, no denied or override rows", report, tt.prohibited)
			}
			assertPartialPolicyEligibility(t, dir, report, "")
		})
	}
}
