package spec

import (
	"errors"
	"fmt"
	"path/filepath"
	"strings"
	"testing"
)

func qaReportFixture(verdict string, extraFrontmatter ...string) string {
	extra := ""
	if len(extraFrontmatter) > 0 {
		extra = strings.Join(extraFrontmatter, "\n") + "\n"
	}
	return fmt.Sprintf(`---
spec: demo
date: 2026-07-04
verdict: %s
%s
surfaces: [cli]
---

# QA Report — Demo
`, verdict, extra)
}

func TestQAVerdictReadsSupportedVerdicts(t *testing.T) {
	for _, verdict := range []string{VerdictPass, VerdictFail, VerdictPartial} {
		t.Run(verdict, func(t *testing.T) {
			specDir := t.TempDir()
			writeFile(t, filepath.Join(specDir, "qa", "qa-report-2026-07-04.md"), qaReportFixture(verdict))

			got, err := QAVerdict(specDir)
			if err != nil {
				t.Fatalf("QAVerdict: %v", err)
			}
			if got != verdict {
				t.Errorf("QAVerdict = %q, want %q", got, verdict)
			}
		})
	}
}

func TestQAVerdictValidatesBlockedCounts(t *testing.T) {
	tests := []struct {
		name             string
		verdict          string
		extraFrontmatter []string
		wantVerdict      string
		wantError        string
	}{
		{
			name:        "absent counts default to zero",
			verdict:     VerdictPass,
			wantVerdict: VerdictPass,
		},
		{
			name:             "environment-blocked pass is readable",
			verdict:          VerdictPass,
			extraFrontmatter: []string{"rows_blocked_environment: 3"},
			wantVerdict:      VerdictPass,
		},
		{
			name:             "finding-blocked pass is unreadable",
			verdict:          VerdictPass,
			extraFrontmatter: []string{"rows_blocked_finding: 1"},
			wantError:        "rows_blocked_finding must be zero when verdict is \"pass\"",
		},
		{
			name:             "finding-blocked partial remains readable",
			verdict:          VerdictPartial,
			extraFrontmatter: []string{"rows_blocked_finding: 1"},
			wantVerdict:      VerdictPartial,
		},
		{
			name:             "finding-blocked fail remains readable",
			verdict:          VerdictFail,
			extraFrontmatter: []string{"rows_blocked_finding: 1"},
			wantVerdict:      VerdictFail,
		},
		{
			name:             "negative environment count is unreadable",
			verdict:          VerdictFail,
			extraFrontmatter: []string{"rows_blocked_environment: -1"},
			wantError:        "rows_blocked_environment must be a non-negative integer",
		},
		{
			name:             "non-integer environment count is unreadable",
			verdict:          VerdictFail,
			extraFrontmatter: []string{"rows_blocked_environment: many"},
			wantError:        "cannot unmarshal",
		},
		{
			name:             "negative finding count is unreadable",
			verdict:          VerdictFail,
			extraFrontmatter: []string{"rows_blocked_finding: -1"},
			wantError:        "rows_blocked_finding must be a non-negative integer",
		},
		{
			name:             "non-integer finding count is unreadable",
			verdict:          VerdictFail,
			extraFrontmatter: []string{"rows_blocked_finding: some"},
			wantError:        "cannot unmarshal",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			specDir := t.TempDir()
			reportPath := filepath.Join(specDir, "qa", "qa-report-2026-07-04.md")
			writeFile(t, reportPath, qaReportFixture(tt.verdict, tt.extraFrontmatter...))

			got, err := QAVerdict(specDir)
			if tt.wantError == "" {
				if err != nil {
					t.Fatalf("QAVerdict: %v", err)
				}
				if got != tt.wantVerdict {
					t.Errorf("QAVerdict = %q, want %q", got, tt.wantVerdict)
				}
				return
			}

			var reportErr QAReportError
			if !errors.As(err, &reportErr) {
				t.Fatalf("error = %v, want QAReportError", err)
			}
			if !strings.Contains(err.Error(), tt.wantError) {
				t.Errorf("error %q does not contain %q", err, tt.wantError)
			}
		})
	}
}

func TestQAVerdictSelectsTheNewestReport(t *testing.T) {
	specDir := t.TempDir()
	writeFile(t, filepath.Join(specDir, "qa", "qa-report-2026-06-28.md"), qaReportFixture(VerdictFail))
	writeFile(t, filepath.Join(specDir, "qa", "qa-report-2026-07-01.md"), qaReportFixture(VerdictPartial))
	writeFile(t, filepath.Join(specDir, "qa", "qa-report-2026-07-04.md"), qaReportFixture(VerdictPass))

	got, err := QAVerdict(specDir)
	if err != nil {
		t.Fatalf("QAVerdict: %v", err)
	}
	if got != VerdictPass {
		t.Errorf("QAVerdict = %q, want the newest report's %q", got, VerdictPass)
	}
}

func TestNewestQAReportOrdersByDateThenRunSequence(t *testing.T) {
	tests := []struct {
		name    string
		reports []string
		want    string
	}{
		{
			name:    "single unsuffixed report",
			reports: []string{"qa-report-2026-07-28.md"},
			want:    "qa-report-2026-07-28.md",
		},
		{
			name:    "same-date rerun beats the first run of the date",
			reports: []string{"qa-report-2026-07-28.md", "qa-report-2026-07-28-02.md"},
			want:    "qa-report-2026-07-28-02.md",
		},
		{
			// The first suffix the naming contract produces. It once parsed to
			// the same sequence as the unsuffixed report, so a path-order
			// tie-break returned the stale one and a Spec could archive on a
			// superseded verdict.
			name:    "the first numeric rerun beats the unsuffixed report",
			reports: []string{"qa-report-2026-07-28.md", "qa-report-2026-07-28-01.md"},
			want:    "qa-report-2026-07-28-01.md",
		},
		{
			name:    "the first numeric rerun wins regardless of input order",
			reports: []string{"qa-report-2026-07-28-01.md", "qa-report-2026-07-28.md"},
			want:    "qa-report-2026-07-28-01.md",
		},
		{
			name:    "a zero suffix still outranks the unsuffixed report",
			reports: []string{"qa-report-2026-07-28.md", "qa-report-2026-07-28-00.md"},
			want:    "qa-report-2026-07-28-00.md",
		},
		{
			name:    "run sequence compares as a number",
			reports: []string{"qa-report-2026-07-28.md", "qa-report-2026-07-28-02.md", "qa-report-2026-07-28-10.md"},
			want:    "qa-report-2026-07-28-10.md",
		},
		{
			name:    "a later date beats an earlier date's highest sequence",
			reports: []string{"qa-report-2026-07-28-10.md", "qa-report-2026-07-29.md"},
			want:    "qa-report-2026-07-29.md",
		},
		{
			name:    "an unsequenced suffix loses to the same date's sequenced reports",
			reports: []string{"qa-report-2026-07-28.md", "qa-report-2026-07-28-ffd6852.md"},
			want:    "qa-report-2026-07-28.md",
		},
		{
			name:    "an unsequenced suffix still beats an earlier date",
			reports: []string{"qa-report-2026-07-27-09.md", "qa-report-2026-07-28-ffd6852.md"},
			want:    "qa-report-2026-07-28-ffd6852.md",
		},
		{
			name:    "an undated name loses to every dated report",
			reports: []string{"qa-report-2026-07-01.md", "qa-report-draft.md", "qa-report-2026-13-40.md"},
			want:    "qa-report-2026-07-01.md",
		},
		{
			name:    "only undated names fall back to path order",
			reports: []string{"qa-report-draft.md", "qa-report-alpha.md"},
			want:    "qa-report-draft.md",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			specDir := t.TempDir()
			for _, report := range tt.reports {
				writeFile(t, filepath.Join(specDir, "qa", report), qaReportFixture(VerdictPass))
			}

			got, err := NewestQAReport(specDir)
			if err != nil {
				t.Fatalf("NewestQAReport: %v", err)
			}
			if want := filepath.Join(specDir, "qa", tt.want); got != want {
				t.Errorf("NewestQAReport = %q, want %q", got, want)
			}
		})
	}
}

func TestNewestQAReportReportsMissingReports(t *testing.T) {
	specDir := t.TempDir()
	writeFile(t, filepath.Join(specDir, "qa", "notes.md"), "# not a report\n")

	if _, err := NewestQAReport(specDir); !errors.Is(err, ErrNoQAReport) {
		t.Fatalf("error = %v, want ErrNoQAReport", err)
	}
}

func TestQAVerdictPrefersTheSameDateRerun(t *testing.T) {
	specDir := t.TempDir()
	writeFile(t, filepath.Join(specDir, "qa", "qa-report-2026-07-28.md"), qaReportFixture(VerdictFail))
	writeFile(t, filepath.Join(specDir, "qa", "qa-report-2026-07-28-02.md"), qaReportFixture(VerdictPass))

	got, err := QAVerdict(specDir)
	if err != nil {
		t.Fatalf("QAVerdict: %v", err)
	}
	if got != VerdictPass {
		t.Errorf("QAVerdict = %q, want the same-date rerun's %q", got, VerdictPass)
	}
}

func TestQAVerdictReportsMissingReports(t *testing.T) {
	tests := []struct {
		name  string
		setup func(t *testing.T, specDir string)
	}{
		{
			name:  "no qa directory",
			setup: func(t *testing.T, specDir string) {},
		},
		{
			name: "qa directory without reports",
			setup: func(t *testing.T, specDir string) {
				writeFile(t, filepath.Join(specDir, "qa", "notes.md"), "# not a report\n")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			specDir := t.TempDir()
			tt.setup(t, specDir)

			_, err := QAVerdict(specDir)
			if !errors.Is(err, ErrNoQAReport) {
				t.Fatalf("error = %v, want ErrNoQAReport", err)
			}
		})
	}
}

func TestQAVerdictReportsUnreadableReports(t *testing.T) {
	tests := []struct {
		name    string
		content string
		wantMsg string
	}{
		{
			name:    "no frontmatter",
			content: "# QA Report — Demo\n",
			wantMsg: "frontmatter",
		},
		{
			name: "no verdict field",
			content: `---
spec: demo
date: 2026-07-04
---

# QA Report — Demo
`,
			wantMsg: "no verdict field",
		},
		{
			name:    "unsupported verdict value",
			content: qaReportFixture("maybe"),
			wantMsg: `"maybe"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			specDir := t.TempDir()
			writeFile(t, filepath.Join(specDir, "qa", "qa-report-2026-07-04.md"), tt.content)

			_, err := QAVerdict(specDir)
			var reportErr QAReportError
			if !errors.As(err, &reportErr) {
				t.Fatalf("error = %v, want QAReportError", err)
			}
			if !strings.Contains(err.Error(), tt.wantMsg) {
				t.Errorf("error %q does not contain %q", err, tt.wantMsg)
			}
			if !strings.Contains(err.Error(), "qa-report-2026-07-04.md") {
				t.Errorf("error %q does not name the report path", err)
			}
		})
	}
}
