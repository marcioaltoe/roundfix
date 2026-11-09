package cli

import (
	"bytes"
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"roundfix/internal/gittest"
	"roundfix/internal/spec"
)

// Suite: QA Report declared-acceptance command
// Invariant: only a selected report accepted by the shared eligibility decision exits zero.
// Boundary IN: command arguments, report parsing, eligibility, streams, and exit status.
// Boundary OUT: newest-report selection and derived-command rendering, owned by internal/spec/task_test.go.

func TestRunQAReportAcceptCommandAcceptsAQualifyingPartial(t *testing.T) {
	t.Parallel()
	specDir := t.TempDir()
	writeQAReportTestFile(t, filepath.Join(specDir, "_prd.md"), `# Test Spec

## Unreachable Acceptance

- criterion: the unavailable journey
  reason: the gate cannot create a pull request
  satisfied-by: task_01
`)
	reportPath := filepath.Join(specDir, "qa", "qa-report-2026-09-21.md")
	writeQAReportTestFile(t, reportPath, "---\nverdict: partial\nrows_blocked_declared: 1\n---\n")

	var stdout, stderr bytes.Buffer
	code := runCLI(t, []string{"qa-report", "accept", reportPath}, &stdout, &stderr)

	if code != exitOK {
		t.Fatalf("qa-report accept exit = %d, want %d; stderr=%q", code, exitOK, stderr.String())
	}
	if stdout.Len() != 0 || stderr.Len() != 0 {
		t.Fatalf("qa-report accept output = stdout %q, stderr %q; want both empty", stdout.String(), stderr.String())
	}
}

func TestQAReportAcceptRefusesAPendingReport(t *testing.T) {
	t.Parallel()
	specDir := t.TempDir()
	reportPath := filepath.Join(specDir, "qa", "qa-report-2026-09-25.md")
	writeQAReportTestFile(t, reportPath, "---\nverdict: pending\n---\n\n# QA Report\n")
	var stdout, stderr bytes.Buffer

	code := runCLI(t, []string{"qa-report", "accept", reportPath}, &stdout, &stderr)

	if code != exitRunFailed || stdout.Len() != 0 {
		t.Fatalf("qa-report accept = exit %d stdout %q, want failed exit and empty stdout", code, stdout.String())
	}
	if !strings.Contains(stderr.String(), `verdict is "pending"`) {
		t.Fatalf("qa-report accept stderr = %q, want pending-verdict refusal", stderr.String())
	}
}

func TestQAReportAcceptRefusesAHollowPass(t *testing.T) {
	t.Parallel()
	specDir := t.TempDir()
	reportPath := filepath.Join(specDir, "qa", "qa-report-2026-09-25.md")
	writeQAReportTestFile(t, reportPath, "---\nverdict: pass\n---\n\n# QA Report\n\n## Results\n\n| # | Status |\n| - | --- |\n")
	var stdout, stderr bytes.Buffer

	code := runCLI(t, []string{"qa-report", "accept", reportPath}, &stdout, &stderr)

	if code != exitRunFailed || stdout.Len() != 0 {
		t.Fatalf("qa-report accept = exit %d stdout %q, want failed exit and empty stdout", code, stdout.String())
	}
	if !strings.Contains(stderr.String(), "records no QA row") {
		t.Fatalf("qa-report accept stderr = %q, want hollow-report refusal", stderr.String())
	}
}

func TestQAReportAcceptRefusesAnEmptyFrontMatter(t *testing.T) {
	t.Parallel()
	specDir := t.TempDir()
	reportPath := filepath.Join(specDir, "qa", "qa-report-2026-09-25.md")
	content := "---\n---\nverdict: pass\n---\n\n# QA Report\n"
	writeQAReportTestFile(t, reportPath, content)
	var stdout, stderr bytes.Buffer

	code := runCLI(t, []string{"qa-report", "accept", reportPath}, &stdout, &stderr)

	if code == exitOK || stdout.Len() != 0 {
		t.Fatalf("qa-report accept = exit %d stdout %q, want non-zero exit and empty stdout", code, stdout.String())
	}
	if !strings.Contains(stderr.String(), "QA Report front matter is empty") {
		t.Fatalf("qa-report accept stderr = %q, want empty-front-matter refusal", stderr.String())
	}
	if got, err := os.ReadFile(reportPath); err != nil || string(got) != content {
		t.Fatalf("qa-report accept changed or removed report: content=%q err=%v", got, err)
	}
}

func TestRunQAReportAcceptCommandFailsClosed(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name      string
		prepare   func(t *testing.T, specDir string) string
		wantError string
	}{
		{
			name: "failed report",
			prepare: func(t *testing.T, specDir string) string {
				path := filepath.Join(specDir, "qa", "qa-report-2026-09-21.md")
				writeQAReportTestFile(t, path, "---\nverdict: fail\n---\n")
				return path
			},
			wantError: `verdict is "fail"`,
		},
		{
			name: "missing report",
			prepare: func(t *testing.T, specDir string) string {
				return filepath.Join(specDir, "qa", "qa-report-2026-09-21.md")
			},
			wantError: "unreadable",
		},
		{
			name: "unparseable report",
			prepare: func(t *testing.T, specDir string) string {
				path := filepath.Join(specDir, "qa", "qa-report-2026-09-21.md")
				writeQAReportTestFile(t, path, "not frontmatter\n")
				return path
			},
			wantError: `front matter must open with a "---" first line`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			specDir := t.TempDir()
			reportPath := tt.prepare(t, specDir)
			var stdout, stderr bytes.Buffer

			code := runCLI(t, []string{"qa-report", "accept", reportPath}, &stdout, &stderr)

			if code != exitRunFailed {
				t.Fatalf("qa-report accept exit = %d, want %d", code, exitRunFailed)
			}
			if stdout.Len() != 0 {
				t.Fatalf("qa-report accept stdout = %q, want empty", stdout.String())
			}
			if !strings.Contains(stderr.String(), tt.wantError) {
				t.Fatalf("qa-report accept stderr = %q, want %q", stderr.String(), tt.wantError)
			}
		})
	}
}

func writeQAReportTestFile(t *testing.T, path, content string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("create QA Report test directory: %v", err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("write QA Report test file: %v", err)
	}
}

const borrowedQAReportContent = "---\nverdict: pass\n---\n\n## Results\n\n| # | Status | Evidence |\n| - | --- | --- |\n| 1 | pass | [run](evidence/2026-11-06-gate/head-exclusive.execution.txt) [replay](evidence/2026-11-06-gate/replay.py.txt) |\n"

const borrowedQAReportRefusal = "QA Report evidence does not belong to this Run: evidence/2026-11-06-gate/head-exclusive.execution.txt (line 9, reused); evidence/2026-11-06-gate/replay.py.txt (line 9, reused)"

func prepareBorrowedQAReport(t *testing.T) (string, string) {
	t.Helper()
	_, root := newImplementWorkspace(t, []implementSeed{implementQAGateSeed(string(spec.StatusFailed))})
	reportRel := "docs/specs/" + implementTestSlug + "/qa/qa-report-2026-11-09.md"
	for _, name := range []string{"head-exclusive.execution.txt", "replay.py.txt"} {
		writeQAReportTestFile(t, filepath.Join(root, filepath.Dir(reportRel), "evidence/2026-11-06-gate", name), "earlier Run evidence\n")
	}
	gittest.Run(t, root, "add", ".")
	gittest.Run(t, root, "commit", "-m", "record earlier evidence")
	gittest.Run(t, root, "commit", "--allow-empty", "-m", "later audited head")
	writeQAReportTestFile(t, filepath.Join(root, reportRel), borrowedQAReportContent)
	return root, reportRel
}

func regenerateBorrowedQAReport(t *testing.T, root, reportRel string) {
	t.Helper()
	for _, name := range []string{"head-exclusive.execution.txt", "replay.py.txt"} {
		writeQAReportTestFile(t, filepath.Join(root, filepath.Dir(reportRel), "evidence/qa-report-2026-11-09", name), "own Run evidence\n")
	}
	writeQAReportTestFile(t, filepath.Join(root, reportRel), strings.ReplaceAll(borrowedQAReportContent, "2026-11-06-gate", "qa-report-2026-11-09"))
}

func TestQAReportAcceptRefusesBorrowedEvidence(t *testing.T) {
	t.Parallel()
	root, reportRel := prepareBorrowedQAReport(t)
	var stdout, stderr bytes.Buffer
	args := []string{"qa-report", "accept", filepath.Join(root, reportRel)}
	code := runCLIContext(t, t.Context(), args, &stdout, &stderr)
	want := "roundfix: QA Report eligibility refused: " + borrowedQAReportRefusal + "\n"
	if code != exitRunFailed || stdout.Len() != 0 || stderr.String() != want {
		t.Fatalf("accept = %d, stdout %q, stderr %q; want refusal %q", code, stdout.String(), stderr.String(), want)
	}
	regenerateBorrowedQAReport(t, root, reportRel)
	stdout.Reset()
	stderr.Reset()
	if code := runCLIContext(t, t.Context(), args, &stdout, &stderr); code != exitOK || stdout.Len() != 0 || stderr.Len() != 0 {
		t.Fatalf("regenerated accept = %d, stdout %q, stderr %q", code, stdout.String(), stderr.String())
	}
}

func TestSettleRefusesAQAReportWithBorrowedEvidence(t *testing.T) {
	t.Parallel()
	root, reportRel := prepareBorrowedQAReport(t)
	taskPath := implementTaskPath(root, "task_qa")
	before := mustRead(t, taskPath)
	var stdout, stderr bytes.Buffer
	args := []string{"settle", "--spec", implementTestSlug, "--task", "task_qa"}
	code := runCLIContext(t, t.Context(), args, &stdout, &stderr)
	want := "Settle surface: " + root + "\nroundfix: settle QA Report is ineligible: " + borrowedQAReportRefusal + " (report " + reportRel + ")\n"
	if code != exitRunFailed || stderr.String() != want {
		t.Fatalf("settle = %d, stdout %q, stderr %q; want refusal %q", code, stdout.String(), stderr.String(), want)
	}
	if got := mustRead(t, taskPath); got != before {
		t.Fatalf("refusal changed Task: %s", got)
	}
	regenerateBorrowedQAReport(t, root, reportRel)
	stdout.Reset()
	stderr.Reset()
	if code := runCLIContext(t, t.Context(), args, &stdout, &stderr); code != exitOK || stderr.String() != "Settle surface: "+root+"\n" {
		t.Fatalf("regenerated settle = %d, stdout %q, stderr %q", code, stdout.String(), stderr.String())
	}
	if got := mustRead(t, taskPath); !strings.Contains(got, "status: completed") {
		t.Fatalf("regenerated fixture did not settle: %s", got)
	}
}

func TestQAReportAcceptRefusesAnUnresolvedEvidenceLink(t *testing.T) {
	t.Parallel()
	_, root := newImplementWorkspace(t, []implementSeed{implementQAGateSeed(string(spec.StatusFailed))})
	writeQAReportTestFile(t, filepath.Join(root, ".gitignore"), "*.log\n")
	gittest.Run(t, root, "add", ".gitignore")
	gittest.Run(t, root, "commit", "-m", "ignore logs")
	reportPath := filepath.Join(root, "docs/specs", implementTestSlug, "qa/qa-report-2026-11-01.md")
	content := "---\nverdict: pass\n---\n\n## Results\n\n| # | Status | Evidence |\n| - | --- | --- |\n| 1 | pass | [consumer](evidence/consumer-tests.log) [journeys](evidence/journeys.txt) |\n"
	writeQAReportTestFile(t, reportPath, content)
	writeQAReportTestFile(t, filepath.Join(filepath.Dir(reportPath), "evidence/consumer-tests.log"), "consumer evidence\n")
	var stdout, stderr bytes.Buffer
	code := runCLIContext(t, t.Context(), []string{"qa-report", "accept", reportPath}, &stdout, &stderr)
	want := "roundfix: QA Report eligibility refused: QA Report evidence links do not resolve in the committed tree: evidence/consumer-tests.log (line 9, ignored); evidence/journeys.txt (line 9, missing)\n"
	if code != exitRunFailed || stdout.Len() != 0 || stderr.String() != want {
		t.Fatalf("accept = %d, stdout %q, stderr %q; want exit 1, empty stdout, stderr %q", code, stdout.String(), stderr.String(), want)
	}
	writeQAReportTestFile(t, filepath.Join(filepath.Dir(reportPath), "evidence/consumer-tests.txt"), "consumer evidence\n")
	writeQAReportTestFile(t, filepath.Join(filepath.Dir(reportPath), "evidence/journeys.txt"), "journey evidence\n")
	writeQAReportTestFile(t, reportPath, strings.ReplaceAll(content, "consumer-tests.log", "consumer-tests.txt"))
	stdout.Reset()
	stderr.Reset()
	code = runCLIContext(t, t.Context(), []string{"qa-report", "accept", reportPath}, &stdout, &stderr)
	if code != exitOK || stdout.Len() != 0 || stderr.Len() != 0 {
		t.Fatalf("repaired accept = %d, stdout %q, stderr %q", code, stdout.String(), stderr.String())
	}
}

func TestSettleRefusesAQAReportWithAnIgnoredEvidenceLink(t *testing.T) {
	t.Parallel()
	_, root := newImplementWorkspace(t, []implementSeed{implementQAGateSeed(string(spec.StatusFailed))})
	writeQAReportTestFile(t, filepath.Join(root, ".gitignore"), "*.log\n")
	gittest.Run(t, root, "add", ".gitignore")
	gittest.Run(t, root, "commit", "-m", "ignore logs")
	reportRel := "docs/specs/" + implementTestSlug + "/qa/qa-report-2026-11-01.md"
	reportPath := filepath.Join(root, reportRel)
	writeQAReportTestFile(t, reportPath, "---\nverdict: pass\n---\n\n## Results\n\n| # | Status | Evidence |\n| - | --- | --- |\n| 1 | pass | [consumer](evidence/consumer-tests.log) |\n")
	writeQAReportTestFile(t, filepath.Join(filepath.Dir(reportPath), "evidence/consumer-tests.log"), "consumer evidence\n")
	taskPath := implementTaskPath(root, "task_qa")
	before := mustRead(t, taskPath)
	var stdout, stderr bytes.Buffer
	code := runCLIContext(t, context.Background(), []string{"settle", "--spec", implementTestSlug, "--task", "task_qa"}, &stdout, &stderr)
	want := "Settle surface: " + root + "\nroundfix: settle QA Report is ineligible: QA Report evidence links do not resolve in the committed tree: evidence/consumer-tests.log (line 9, ignored) (report " + reportRel + ")\n"
	if code != exitRunFailed || stderr.String() != want {
		t.Fatalf("settle = %d, stdout %q, stderr %q; want exit 1, stderr %q", code, stdout.String(), stderr.String(), want)
	}
	if got := mustRead(t, taskPath); got != before {
		t.Fatalf("refusal changed Task file: %s", got)
	}
	logPath := filepath.Join(filepath.Dir(reportPath), "evidence/consumer-tests.log")
	txtPath := filepath.Join(filepath.Dir(reportPath), "evidence/consumer-tests.txt")
	if err := os.Rename(logPath, txtPath); err != nil {
		t.Fatal(err)
	}
	writeQAReportTestFile(t, reportPath, strings.ReplaceAll(mustRead(t, reportPath), "consumer-tests.log", "consumer-tests.txt"))
	stdout.Reset()
	stderr.Reset()
	code = runCLIContext(t, t.Context(), []string{"settle", "--spec", implementTestSlug, "--task", "task_qa"}, &stdout, &stderr)
	if code != exitOK || stderr.String() != "Settle surface: "+root+"\n" {
		t.Fatalf("repaired settle = %d, stdout %q, stderr %q", code, stdout.String(), stderr.String())
	}
	if got := mustRead(t, taskPath); !strings.Contains(got, "status: completed") {
		t.Fatalf("repaired fixture Task did not settle: %s", got)
	}
	if got := strings.TrimSpace(gittest.Run(t, root, "ls-files", "--", filepath.ToSlash(filepath.Join(filepath.Dir(reportRel), "evidence/consumer-tests.txt")))); got == "" {
		t.Fatal("repaired settlement omitted the saved evidence")
	}
}
