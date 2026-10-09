// Suite: report shape findings
// Invariant: import validation uses the same terminal row checks as the mechanical stage.
// Boundary IN: report loader, shape and evidence path detectors
// Boundary OUT: Git import (daemon/qa_prior_pass_test.go)
package speccheck

import (
	"os"
	"path/filepath"
	"testing"
)

func shapeReportForTest(t *testing.T, status string) (string, string) {
	t.Helper()
	root := t.TempDir()
	path := "qa-report-2026-10-01.md"
	content := "---\nverdict: pass\nrows_blocked_environment: 0\nrows_blocked_finding: 0\nrows_blocked_declared: 0\n---\n\n## Results\n\n| # | Status | Provenance | Evidence |\n| --- | --- | --- | --- |\n| 1 | " + status + " | criterion | observed |\n"
	if err := os.WriteFile(filepath.Join(root, path), []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	return root, path
}

func TestReportShapeFindingsNamesAPendingRow(t *testing.T) {
	t.Parallel()
	root, path := shapeReportForTest(t, "pending")
	findings, err := ReportShapeFindings(root, path)
	if err != nil || len(findings) != 1 || findings[0].Code != CodeMechanicalReportShape {
		t.Fatalf("findings = %+v %v", findings, err)
	}
}

func TestReportShapeFindingsAcceptsAClosedReport(t *testing.T) {
	t.Parallel()
	root, path := shapeReportForTest(t, "pass")
	findings, err := ReportShapeFindings(root, path)
	if err != nil || len(findings) != 0 {
		t.Fatalf("findings = %+v %v", findings, err)
	}
}
