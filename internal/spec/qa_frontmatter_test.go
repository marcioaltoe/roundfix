package spec

import (
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"roundfix/internal/gittest"
)

// Suite: QA Report front matter reader parity
// Invariant: the Go QA Report reader accepts exactly the front matter shapes accepted by the derived shell verification.
// Boundary IN: report selection, front matter bounds, verdict-line shape, and shell command execution.
// Boundary OUT: verdict eligibility after parsing, owned by qa_test.go and archive_test.go.

func TestReadQAReportRefusesAnEmptyFrontMatter(t *testing.T) {
	t.Parallel()
	assertQAReportFrontmatterError(t, "---\n---\nverdict: pass\n---\n", "QA Report front matter is empty")
}

func TestReadQAReportRefusesAMissingOpeningLine(t *testing.T) {
	t.Parallel()
	assertQAReportFrontmatterError(t, "verdict: pass\n---\n", `QA Report front matter must open with a "---" first line`)
}

func TestReadQAReportRefusesAnUnclosedFrontMatter(t *testing.T) {
	t.Parallel()
	assertQAReportFrontmatterError(t, "---\nverdict: pass\n", `QA Report front matter has no closing "---" line`)
}

func TestReadQAReportRefusesADuplicatedVerdictLine(t *testing.T) {
	t.Parallel()
	assertQAReportFrontmatterError(t, "---\nverdict: pass\nverdict: fail\n---\n", `QA Report front matter has 2 "verdict:" lines; expected exactly 1`)
}

func TestReadQAReportReadsAWellFormedFrontMatter(t *testing.T) {
	t.Parallel()
	const body = "\n# QA Report\n\n## Results\n\n| # | Status | Evidence |\n| - | --- | --- |\n| R01 | pass | observed |\n"
	content := []byte("---\nverdict: pass\nrows_blocked_environment: 0\n---\n" + body)
	path := filepath.Join(t.TempDir(), "qa-report-2026-09-28.md")
	if err := os.WriteFile(path, content, 0o644); err != nil {
		t.Fatalf("write QA Report: %v", err)
	}

	report, err := ReadQAReportFile(path)
	if err != nil {
		t.Fatalf("ReadQAReportFile: %v", err)
	}
	if report.Verdict != VerdictPass {
		t.Fatalf("Verdict = %q, want %q", report.Verdict, VerdictPass)
	}
	_, gotBody, err := splitQAReportFrontmatter(content)
	if err != nil {
		t.Fatalf("splitQAReportFrontmatter: %v", err)
	}
	if string(gotBody) != body {
		t.Fatalf("body = %q, want %q", gotBody, body)
	}
}

func TestQAReportReaderAgreesWithTheDerivedVerification(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		report string
	}{
		{name: "well-formed pass", report: qaFrontmatterReadablePass()},
		{name: "two consecutive opening lines", report: "---\n---\nverdict: pass\n---\n\n# QA Report\n"},
		{name: "two verdict lines", report: "---\nverdict: pass\nverdict: fail\n---\n\n# QA Report\n"},
		{name: "no closing line", report: "---\nverdict: pass\n\n# QA Report\n"},
		{name: "empty verdict value", report: "---\nverdict: \n---\n\n# QA Report\n"},
		{name: "blank line before opening", report: "\n---\nverdict: pass\n---\n\n# QA Report\n"},
		{name: "four-dash closing line", report: "---\nverdict: pass\n----\n\n# QA Report\n"},
		{name: "indented verdict line", report: "---\n verdict: pass\n---\n\n# QA Report\n"},
		{name: "CRLF line endings", report: strings.ReplaceAll(qaFrontmatterReadablePass(), "\n", "\r\n")},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assertQAReportReadersAgree(t, "qa-report-2026-09-28.md", []byte(test.report))
		})
	}
}

func TestQAReportReaderAgreesWithTheDerivedVerificationOnTheArchive(t *testing.T) {
	t.Parallel()
	_, testFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("runtime.Caller could not locate the repository")
	}
	repoRoot := gittest.PinnedHistory(t, filepath.Join(filepath.Dir(testFile), "..", ".."), ":(glob)"+ArchiveDir(ArchiveKindSpec)+"/*/qa/qa-report-*.md")
	pattern := filepath.Join(repoRoot, filepath.FromSlash("docs/history/specs/*/qa/qa-report-*.md"))
	reportPaths, err := filepath.Glob(pattern)
	if err != nil {
		t.Fatalf("find archived QA Reports: %v", err)
	}
	if len(reportPaths) == 0 {
		t.Fatal("archived QA Report corpus is empty")
	}

	for _, reportPath := range reportPaths {
		relative, err := filepath.Rel(repoRoot, reportPath)
		if err != nil {
			t.Fatalf("make archived QA Report path relative: %v", err)
		}
		t.Run(filepath.ToSlash(relative), func(t *testing.T) {
			content, err := os.ReadFile(reportPath)
			if err != nil {
				t.Fatalf("read archived QA Report: %v", err)
			}
			reportName := filepath.Base(reportPath)
			parsedName := parseQAReportName(reportName)
			if parsedName.date == "" || !parsedName.sequenced {
				// This sweep isolates the front matter readers. A few legacy
				// archive names predate the derived verifier's date/sequence name
				// contract, so give those report bodies one valid sole-report name.
				reportName = "qa-report-2026-09-28.md"
			}
			assertQAReportReadersAgree(t, reportName, content)
		})
	}
}

func assertQAReportFrontmatterError(t *testing.T, content string, want string) {
	t.Helper()
	path := filepath.Join(t.TempDir(), "qa-report-2026-09-28.md")
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("write QA Report: %v", err)
	}
	_, err := ReadQAReportFile(path)
	var reportErr QAReportError
	if !errors.As(err, &reportErr) {
		t.Fatalf("error = %v, want QAReportError", err)
	}
	if reportErr.Err == nil || reportErr.Err.Error() != want {
		t.Fatalf("QAReportError cause = %v, want %q", reportErr.Err, want)
	}
}

func assertQAReportReadersAgree(t *testing.T, reportName string, content []byte) {
	t.Helper()
	const slug = "qa-frontmatter-fixture"
	repoRoot := t.TempDir()
	specDir := filepath.Join(repoRoot, "docs", "specs", slug)
	reportPath := filepath.Join(specDir, "qa", reportName)
	if err := os.MkdirAll(filepath.Dir(reportPath), 0o755); err != nil {
		t.Fatalf("create QA Report directory: %v", err)
	}
	if err := os.WriteFile(reportPath, content, 0o644); err != nil {
		t.Fatalf("copy QA Report: %v", err)
	}

	_, readErr := ReadQAReport(specDir)
	goReadable := readErr == nil
	command := exec.Command("sh", "-c", DerivedQAVerification(slug)[0])
	command.Dir = repoRoot
	output, verificationErr := command.CombinedOutput()
	derivedReadable := verificationErr == nil
	var exitErr *exec.ExitError
	if verificationErr != nil && !errors.As(verificationErr, &exitErr) {
		t.Fatalf("run derived QA Verification: %v", verificationErr)
	}
	if goReadable != derivedReadable {
		t.Fatalf("reader agreement for %q: Go readable=%t (error=%v), derived readable=%t (error=%v, output=%q)", reportName, goReadable, readErr, derivedReadable, verificationErr, output)
	}
}

func qaFrontmatterReadablePass() string {
	return "---\nverdict: pass\nrows_blocked_environment: 0\nrows_blocked_finding: 0\nrows_blocked_declared: 0\n---\n\n# QA Report\n\n## Results\n\n| # | Status | Evidence |\n| - | --- | --- |\n| R01 | pass | observed |\n"
}
