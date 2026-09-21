package cli

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// Suite: QA Report declared-acceptance command
// Invariant: only a selected report accepted by the shared eligibility decision exits zero.
// Boundary IN: command arguments, report parsing, eligibility, streams, and exit status.
// Boundary OUT: newest-report selection and derived-command rendering, owned by internal/spec/task_test.go.

func TestRunQAReportAcceptCommandAcceptsAQualifyingPartial(t *testing.T) {
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

func TestRunQAReportAcceptCommandFailsClosed(t *testing.T) {
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
			wantError: "frontmatter",
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
