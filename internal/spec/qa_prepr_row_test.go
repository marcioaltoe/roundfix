// Suite: pre-PR Pull Request QA report eligibility
// Invariant: only an exact pre-PR Pull Request Results row excuses one environment-blocked row from a declared partial.
// Boundary IN: QA report parsing and QAReportEligibility.
// Boundary OUT: mechanical QA count validation in internal/speccheck.
package spec

import (
	"fmt"
	"path/filepath"
	"strings"
	"testing"
)

func TestAPartialQualifiesWhenItsOnlyEnvironmentRowIsThePrePullRequestRow(t *testing.T) {
	t.Parallel()

	specDir := t.TempDir()
	writePrePullRequestDeclaration(t, specDir)
	report := readPrePullRequestReport(t, specDir, VerdictPartial, 1, 1, fmt.Sprintf(`
## Results

| Requirement | STATUS | Evidence | provenance |
| --- | --- | --- | --- |
| Pull Request journey | %s | no Pull Request exists yet | CLI; %s, repository |
| Declared unreachable | blocked (declared) | declaration | PRD Goal 2 |
`, QANoOpenPullRequestStatus, QAPullRequestRowSource))

	if report.RowsBlockedPrePullRequest != 1 {
		t.Fatalf("RowsBlockedPrePullRequest = %d, want 1", report.RowsBlockedPrePullRequest)
	}
	if err := QAReportEligibility(specDir, report); err != nil {
		t.Fatalf("QAReportEligibility: %v", err)
	}
}

func TestAPartialWithAnotherEnvironmentRowStillRefuses(t *testing.T) {
	t.Parallel()

	specDir := t.TempDir()
	writePrePullRequestDeclaration(t, specDir)
	report := readPrePullRequestReport(t, specDir, VerdictPartial, 2, 1, fmt.Sprintf(`
## Results

| Requirement | Status | Evidence | Provenance |
| --- | --- | --- | --- |
| Pull Request journey | %s | no Pull Request exists yet | %s |
| Release journey | blocked (environment: credentials unavailable) | no release credential | Release row |
| Declared unreachable | blocked (declared) | declaration | PRD Goal 2 |
`, QANoOpenPullRequestStatus, QAPullRequestRowSource))

	if report.RowsBlockedPrePullRequest != 1 {
		t.Fatalf("RowsBlockedPrePullRequest = %d, want 1", report.RowsBlockedPrePullRequest)
	}
	err := QAReportEligibility(specDir, report)
	const want = "rows_blocked_environment is 2, 1 outside the pre-PR Pull Request row; expected 0 outside it"
	if err == nil || err.Error() != want {
		t.Fatalf("error = %v, want %q", err, want)
	}
}

func TestANoPullRequestStatusWithoutThePullRequestSourceStillRefuses(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		status     string
		provenance string
		body       string
	}{
		{
			name:       "another provenance item",
			status:     QANoOpenPullRequestStatus,
			provenance: "CLI row; repository row",
		},
		{
			name:       "provenance merely contains the source words",
			status:     QANoOpenPullRequestStatus,
			provenance: "not a " + QAPullRequestRowSource,
		},
		{
			name:       "status differs in case",
			status:     strings.ToUpper(QANoOpenPullRequestStatus),
			provenance: QAPullRequestRowSource,
		},
		{
			name: "matching table outside Results",
			body: fmt.Sprintf(`
## Evidence

| Status | Provenance |
| --- | --- |
| %s | %s |

## Results

| Status | Provenance |
| --- | --- |
| blocked (environment: no browser) | Frontend row |
`, QANoOpenPullRequestStatus, QAPullRequestRowSource),
		},
		{
			name: "matching table inside a fence",
			body: fmt.Sprintf(`
## Results

~~~markdown
| Status | Provenance |
| --- | --- |
| %s | %s |
~~~

| Status | Provenance |
| --- | --- |
| blocked (environment: no browser) | Frontend row |
`, QANoOpenPullRequestStatus, QAPullRequestRowSource),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			specDir := t.TempDir()
			writePrePullRequestDeclaration(t, specDir)
			body := tt.body
			if body == "" {
				body = fmt.Sprintf(`
## Results

| Requirement | Status | Provenance |
| --- | --- | --- |
| Pull Request journey | %s | %s |
| Declared unreachable | blocked (declared) | PRD Goal 2 |
`, tt.status, tt.provenance)
			}
			report := readPrePullRequestReport(t, specDir, VerdictPartial, 1, 1, body)

			if report.RowsBlockedPrePullRequest != 0 {
				t.Fatalf("RowsBlockedPrePullRequest = %d, want 0", report.RowsBlockedPrePullRequest)
			}
			err := QAReportEligibility(specDir, report)
			const want = "rows_blocked_environment is 1; expected 0"
			if err == nil || err.Error() != want {
				t.Fatalf("error = %v, want %q", err, want)
			}
		})
	}
}

func TestAResultsTableWithoutProvenanceGainsNoPrePullRequestRow(t *testing.T) {
	t.Parallel()

	specDir := t.TempDir()
	writePrePullRequestDeclaration(t, specDir)
	report := readPrePullRequestReport(t, specDir, VerdictPartial, 1, 1, fmt.Sprintf(`
## Results

| Requirement | Status | Evidence |
| --- | --- | --- |
| Pull Request journey | %s | %s |
| Declared unreachable | blocked (declared) | PRD Goal 2 |
`, QANoOpenPullRequestStatus, QAPullRequestRowSource))

	if report.RowsBlockedPrePullRequest != 0 {
		t.Fatalf("RowsBlockedPrePullRequest = %d, want 0", report.RowsBlockedPrePullRequest)
	}
	err := QAReportEligibility(specDir, report)
	const want = "rows_blocked_environment is 1; expected 0"
	if err == nil || err.Error() != want {
		t.Fatalf("error = %v, want %q", err, want)
	}
}

func TestAPassIsUnchangedByThePrePullRequestRow(t *testing.T) {
	t.Parallel()

	t.Run("pass stays eligible", func(t *testing.T) {
		t.Parallel()

		report := QAReport{
			Verdict:                   VerdictPass,
			RowsBlockedEnvironment:    1,
			RowsBlockedPrePullRequest: 1,
		}
		if err := QAReportEligibility(t.TempDir(), report); err != nil {
			t.Fatalf("QAReportEligibility: %v", err)
		}
	})

	t.Run("partial without a declared row qualifies", func(t *testing.T) {
		t.Parallel()

		report := QAReport{
			Verdict:                   VerdictPartial,
			RowsBlockedEnvironment:    1,
			RowsBlockedPrePullRequest: 1,
		}
		if err := QAReportEligibility(t.TempDir(), report); err != nil {
			t.Fatalf("QAReportEligibility: %v", err)
		}
	})
}

func readPrePullRequestReport(t *testing.T, specDir, verdict string, environment, declared int, body string) QAReport {
	t.Helper()

	path := filepath.Join(specDir, "qa", "qa-report-2026-09-29.md")
	writeFile(t, path, fmt.Sprintf(`---
verdict: %s
rows_blocked_environment: %d
rows_blocked_finding: 0
rows_blocked_declared: %d
---

# QA Report
%s`, verdict, environment, declared, body))
	report, err := ReadQAReportFile(path)
	if err != nil {
		t.Fatalf("ReadQAReportFile: %v", err)
	}
	return report
}

func writePrePullRequestDeclaration(t *testing.T, specDir string) {
	t.Helper()

	writeFile(t, filepath.Join(specDir, "_prd.md"), `# Test Spec

## Unreachable Acceptance

- criterion: declared acceptance
  reason: unavailable environment
  satisfied-by: task_01
`)
}
