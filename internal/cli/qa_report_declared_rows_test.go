package cli

import (
	"bytes"
	"fmt"
	"path/filepath"
	"strings"
	"testing"
)

func TestQAReportAcceptMatchesEachDeclaredRowToItsCriterion(t *testing.T) {
	t.Parallel()
	for _, tt := range []struct {
		name, id     string
		metric, code int
	}{
		{"Quaestor shape", "Q11", 11, 0},
		{"undeclared row", "Q12", 12, 1},
	} {
		t.Run(tt.name, func(t *testing.T) {
			dir := t.TempDir()
			writeQAReportTestFile(t, filepath.Join(dir, "_prd.md"), "# PRD\n\n## Unreachable Acceptance\n\n- criterion: Success Metric 10, after release\n  criterion: Success Metric 11, after release\n  reason: future evidence\n  satisfied-by: operator\n")
			path := filepath.Join(dir, "qa", "qa-report-2026-10-27.md")
			writeQAReportTestFile(t, path, fmt.Sprintf("---\nverdict: partial\nrows_blocked_declared: 2\n---\n\n## Results\n\n| # | Status |\n| - | --- |\n| Q10 | blocked (declared: Success Metric 10) |\n| %s | blocked (declared: Success Metric %d) |\n", tt.id, tt.metric))
			var stdout, stderr bytes.Buffer
			code := runCLI(t, []string{"qa-report", "accept", path}, &stdout, &stderr)
			if code != tt.code || stdout.Len() != 0 {
				t.Fatalf("exit=%d stdout=%q stderr=%q; want exit %d and no stdout", code, stdout.String(), stderr.String(), tt.code)
			}
			if tt.code == 0 && stderr.Len() != 0 {
				t.Fatalf("accepted report stderr=%q", stderr.String())
			}
			if tt.code != 0 && !strings.Contains(stderr.String(), "declared-blocked row Q12") {
				t.Fatalf("refusal lacks unmatched row: %q", stderr.String())
			}
		})
	}
}
