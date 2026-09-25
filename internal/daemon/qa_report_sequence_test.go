package daemon

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"roundfix/internal/spec"
	"roundfix/internal/speccheck"
)

// Suite: mechanical QA Report allocation
// Invariant: the path written for a date is the path NewestQAReport selects.
// Boundary IN: report-directory inspection, exclusive creation, and report-name ordering.
// Boundary OUT: QA Report contents and verdict settlement, owned by qa_hollow_report_test.go.

func TestMechanicalQAReportAllocatesAboveTheHighestSequence(t *testing.T) {
	t.Parallel()
	engine, plan, date := mechanicalQAReportSequenceFixture(t)
	writeExistingQAReportForSequenceTest(t, plan.Spec.Dir, fmt.Sprintf("qa-report-%s.md", date))
	writeExistingQAReportForSequenceTest(t, plan.Spec.Dir, fmt.Sprintf("qa-report-%s-02.md", date))

	reportPath, err := engine.writeMechanicalQAReport(context.Background(), plan, speccheck.MechanicalResult{Blocking: true})
	if err != nil {
		t.Fatalf("writeMechanicalQAReport: %v", err)
	}
	assertMechanicalQAReportPathForSequenceTest(t, plan, reportPath, fmt.Sprintf("qa-report-%s-03.md", date))
}

func TestMechanicalQAReportStartsUnsuffixedOnAFreshDate(t *testing.T) {
	t.Parallel()
	engine, plan, date := mechanicalQAReportSequenceFixture(t)

	reportPath, err := engine.writeMechanicalQAReport(context.Background(), plan, speccheck.MechanicalResult{Blocking: true})
	if err != nil {
		t.Fatalf("writeMechanicalQAReport: %v", err)
	}
	assertMechanicalQAReportPathForSequenceTest(t, plan, reportPath, fmt.Sprintf("qa-report-%s.md", date))
}

func TestMechanicalQAReportFollowsTheUnsuffixedReportWithSequenceOne(t *testing.T) {
	t.Parallel()
	engine, plan, date := mechanicalQAReportSequenceFixture(t)
	writeExistingQAReportForSequenceTest(t, plan.Spec.Dir, fmt.Sprintf("qa-report-%s.md", date))

	reportPath, err := engine.writeMechanicalQAReport(context.Background(), plan, speccheck.MechanicalResult{Blocking: true})
	if err != nil {
		t.Fatalf("writeMechanicalQAReport: %v", err)
	}
	assertMechanicalQAReportPathForSequenceTest(t, plan, reportPath, fmt.Sprintf("qa-report-%s-01.md", date))
}

func TestMechanicalQAReportRefusesBehindALaterDatedReport(t *testing.T) {
	t.Parallel()
	engine, plan, _ := mechanicalQAReportSequenceFixture(t)
	laterName := fmt.Sprintf("qa-report-%s.md", taskCycleNowForTest().AddDate(0, 0, 1).Format("2006-01-02"))
	laterPath := writeExistingQAReportForSequenceTest(t, plan.Spec.Dir, laterName)

	reportPath, err := engine.writeMechanicalQAReport(context.Background(), plan, speccheck.MechanicalResult{Blocking: true})
	if err == nil {
		t.Fatalf("writeMechanicalQAReport returned path %q, want refusal behind %s", reportPath, laterName)
	}
	if reportPath != "" {
		t.Errorf("writeMechanicalQAReport path = %q, want empty path on refusal", reportPath)
	}
	if !strings.Contains(err.Error(), laterName) {
		t.Errorf("writeMechanicalQAReport error = %q, want it to name %q", err, laterName)
	}
	entries, readErr := os.ReadDir(filepath.Join(plan.Spec.Dir, "qa"))
	if readErr != nil {
		t.Fatalf("read QA Report directory after refusal: %v", readErr)
	}
	if len(entries) != 1 || entries[0].Name() != laterName {
		t.Errorf("QA Report directory after refusal = %#v, want only %q", entries, laterName)
	}
	newest, newestErr := spec.NewestQAReport(plan.Spec.Dir)
	if newestErr != nil {
		t.Fatalf("NewestQAReport: %v", newestErr)
	}
	if newest != laterPath {
		t.Errorf("NewestQAReport = %q, want unchanged later report %q", newest, laterPath)
	}
}

func mechanicalQAReportSequenceFixture(t *testing.T) (*Engine, TaskPlan, string) {
	t.Helper()
	now := taskCycleNowForTest()
	repoRoot := t.TempDir()
	specDir := filepath.Join(repoRoot, "docs", "specs", taskCycleSlug)
	engine := &Engine{deps: Dependencies{Now: func() time.Time { return now }}}
	plan := TaskPlan{WorkDir: repoRoot, Spec: spec.Spec{Slug: taskCycleSlug, Dir: specDir}}
	return engine, plan, now.Format("2006-01-02")
}

func writeExistingQAReportForSequenceTest(t *testing.T, specDir, name string) string {
	t.Helper()
	path := filepath.Join(specDir, "qa", name)
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("create QA Report directory: %v", err)
	}
	if err := os.WriteFile(path, []byte("---\nverdict: pass\n---\n"), 0o644); err != nil {
		t.Fatalf("write existing QA Report %q: %v", path, err)
	}
	return path
}

func assertMechanicalQAReportPathForSequenceTest(t *testing.T, plan TaskPlan, reportPath, wantName string) {
	t.Helper()
	wantPath := filepath.Join(plan.Spec.Dir, "qa", wantName)
	if got := filepath.Join(plan.WorkDir, filepath.FromSlash(reportPath)); got != wantPath {
		t.Errorf("writeMechanicalQAReport path = %q, want %q", got, wantPath)
	}
	newest, err := spec.NewestQAReport(plan.Spec.Dir)
	if err != nil {
		t.Fatalf("NewestQAReport: %v", err)
	}
	if newest != wantPath {
		t.Errorf("NewestQAReport = %q, want allocated path %q", newest, wantPath)
	}
}
