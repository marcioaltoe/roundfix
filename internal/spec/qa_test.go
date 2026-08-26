package spec

import (
	"errors"
	"fmt"
	"path/filepath"
	"runtime"
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
	t.Parallel()
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
	t.Parallel()
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
			name:             "declared-blocked pass is unreadable",
			verdict:          VerdictPass,
			extraFrontmatter: []string{"rows_blocked_declared: 3"},
			wantError:        "rows_blocked_declared must be zero when verdict is \"pass\"",
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

func TestReadQAReportBlockedCounts(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name                   string
		fixture                string
		wantVerdict            string
		wantBlockedEnvironment int
		wantBlockedFinding     int
		wantBlockedDeclared    int
	}{
		{
			name:        "absent declared count defaults to zero",
			fixture:     "absent",
			wantVerdict: VerdictPass,
		},
		{
			name:        "explicit zero declared count remains zero",
			fixture:     "zero",
			wantVerdict: VerdictPass,
		},
		{
			name:                   "positive declared count is independent",
			fixture:                "positive",
			wantVerdict:            VerdictPartial,
			wantBlockedEnvironment: 1,
			wantBlockedFinding:     2,
			wantBlockedDeclared:    3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			report, err := ReadQAReport(filepath.Join("testdata", "qa-blocked-counts", tt.fixture))
			if err != nil {
				t.Fatalf("ReadQAReport: %v", err)
			}
			if report.Verdict != tt.wantVerdict {
				t.Errorf("Verdict = %q, want %q", report.Verdict, tt.wantVerdict)
			}
			if report.RowsBlockedEnvironment != tt.wantBlockedEnvironment {
				t.Errorf("RowsBlockedEnvironment = %d, want %d", report.RowsBlockedEnvironment, tt.wantBlockedEnvironment)
			}
			if report.RowsBlockedFinding != tt.wantBlockedFinding {
				t.Errorf("RowsBlockedFinding = %d, want %d", report.RowsBlockedFinding, tt.wantBlockedFinding)
			}
			if report.RowsBlockedDeclared != tt.wantBlockedDeclared {
				t.Errorf("RowsBlockedDeclared = %d, want %d", report.RowsBlockedDeclared, tt.wantBlockedDeclared)
			}
		})
	}
}

func TestReadQAReportRejectsInvalidDeclaredCount(t *testing.T) {
	t.Parallel()
	for _, fixture := range []string{"negative", "non-integer"} {
		t.Run(fixture, func(t *testing.T) {
			_, err := ReadQAReport(filepath.Join("testdata", "qa-blocked-counts", fixture))
			var reportErr QAReportError
			if !errors.As(err, &reportErr) {
				t.Fatalf("error = %v, want QAReportError", err)
			}
			if !strings.Contains(err.Error(), "rows_blocked_declared") {
				t.Errorf("error %q does not name rows_blocked_declared", err)
			}
		})
	}
}

func TestArchivedQAReportCorpusRemainsReadable(t *testing.T) {
	t.Parallel()
	_, testFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("runtime.Caller could not locate the repository")
	}
	pattern := archiveTestRepositoryPath(filepath.Join(filepath.Dir(testFile), "..", ".."), ArchiveKindSpec, "*", "qa", "qa-report-*.md")
	reports, err := filepath.Glob(pattern)
	if err != nil {
		t.Fatalf("find archived QA Reports: %v", err)
	}
	if len(reports) == 0 {
		t.Fatal("archived QA Report corpus is empty")
	}
	seen := make(map[string]struct{})
	for _, reportPath := range reports {
		specDir := filepath.Dir(filepath.Dir(reportPath))
		if _, ok := seen[specDir]; ok {
			continue
		}
		seen[specDir] = struct{}{}
		testSpecDir := specDir
		t.Run(filepath.Base(testSpecDir), func(t *testing.T) {
			if _, err := ReadQAReport(testSpecDir); err != nil {
				t.Fatalf("ReadQAReport(%q): %v", testSpecDir, err)
			}
		})
	}
}

func TestQAVerdictSelectsTheNewestReport(t *testing.T) {
	t.Parallel()
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
	t.Parallel()
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
	t.Parallel()
	specDir := t.TempDir()
	writeFile(t, filepath.Join(specDir, "qa", "notes.md"), "# not a report\n")

	if _, err := NewestQAReport(specDir); !errors.Is(err, ErrNoQAReport) {
		t.Fatalf("error = %v, want ErrNoQAReport", err)
	}
}

func TestQAVerdictPrefersTheSameDateRerun(t *testing.T) {
	t.Parallel()
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
	t.Parallel()
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
	t.Parallel()
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

// qaResultsRows returns the data rows of the report's Results table: every
// table line under the `## Results` heading except the header and its
// separator, up to the next heading. The refusal contract is about how many
// rows a refused gate materializes, so the count has to be read from the table
// rather than from the frontmatter that claims it.
func qaResultsRows(t *testing.T, report string) []string {
	t.Helper()
	var rows []string
	inResults := false
	for _, line := range strings.Split(report, "\n") {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "## ") {
			if inResults {
				break
			}
			inResults = strings.TrimPrefix(trimmed, "## ") == "Results"
			continue
		}
		if !inResults || !strings.HasPrefix(trimmed, "|") {
			continue
		}
		if strings.Contains(trimmed, "Status") || strings.Trim(trimmed, "| -") == "" {
			continue
		}
		rows = append(rows, trimmed)
	}
	return rows
}

func TestWritePreconditionRefusalReportWritesOneTerminalRow(t *testing.T) {
	t.Parallel()
	var report strings.Builder
	if err := WritePreconditionRefusalReport(&report, PreconditionRefusal{
		CheckName: "strict",
		Reason:    "SC-VOCABULARY-UNDOCUMENTED",
	}); err != nil {
		t.Fatalf("WritePreconditionRefusalReport: %v", err)
	}

	rows := qaResultsRows(t, report.String())
	if len(rows) != 1 {
		t.Fatalf("Results rows = %q, want exactly one terminal row", rows)
	}
	wantRow := fmt.Sprintf("| %s | %s | %s |", QAPreconditionRowID, QAPreconditionRowStatus, QAPreconditionRowProvenance)
	if rows[0] != wantRow {
		t.Errorf("terminal row = %q, want %q", rows[0], wantRow)
	}
	for _, want := range []string{
		"verdict: " + VerdictFail + "\n",
		"rows_blocked_precondition: 1\n",
		"rows_blocked_environment: 0\n",
		"rows_blocked_finding: 0\n",
		"rows_blocked_declared: 0\n",
		`precondition_check: "strict"` + "\n",
		`precondition_reason: "SC-VOCABULARY-UNDOCUMENTED"` + "\n",
	} {
		if !strings.Contains(report.String(), want) {
			t.Errorf("report does not record %q:\n%s", want, report.String())
		}
	}
}

func TestWritePreconditionRefusalReportIsReadableAsARefusal(t *testing.T) {
	t.Parallel()
	var content strings.Builder
	if err := WritePreconditionRefusalReport(&content, PreconditionRefusal{
		CheckName: "strict",
		Reason:    "SC-REQUIREMENT-CONTRADICTORY",
	}); err != nil {
		t.Fatalf("WritePreconditionRefusalReport: %v", err)
	}
	specDir := t.TempDir()
	writeFile(t, filepath.Join(specDir, "qa", "qa-report-2026-08-25.md"), content.String())

	report, err := ReadQAReport(specDir)
	if err != nil {
		t.Fatalf("ReadQAReport of a refusal the gate itself wrote: %v", err)
	}
	if report.Verdict != VerdictFail {
		t.Errorf("Verdict = %q, want %q", report.Verdict, VerdictFail)
	}
	if report.RowsBlockedPrecondition != 1 {
		t.Errorf("RowsBlockedPrecondition = %d, want 1", report.RowsBlockedPrecondition)
	}
	if report.RowsBlockedEnvironment != 0 || report.RowsBlockedFinding != 0 || report.RowsBlockedDeclared != 0 {
		t.Errorf("other blocked causes = %d/%d/%d, want zero; a refusal measured nothing else",
			report.RowsBlockedEnvironment, report.RowsBlockedFinding, report.RowsBlockedDeclared)
	}
}

func TestWritePreconditionRefusalReportKeepsTheRefusalOnOneLine(t *testing.T) {
	t.Parallel()
	var content strings.Builder
	if err := WritePreconditionRefusalReport(&content, PreconditionRefusal{
		CheckName: "spec check --strict",
		Reason:    "SC-VOCABULARY-UNDOCUMENTED:\nterm | \"Run Ledger\" is undocumented",
	}); err != nil {
		t.Fatalf("WritePreconditionRefusalReport: %v", err)
	}
	report := content.String()
	if rows := qaResultsRows(t, report); len(rows) != 1 {
		t.Fatalf("Results rows = %q, want exactly one terminal row", rows)
	}
	wantReason := `precondition_reason: "SC-VOCABULARY-UNDOCUMENTED: term | \"Run Ledger\" is undocumented"`
	if !strings.Contains(report, wantReason+"\n") {
		t.Errorf("report does not record the collapsed reason %q:\n%s", wantReason, report)
	}

	specDir := t.TempDir()
	writeFile(t, filepath.Join(specDir, "qa", "qa-report-2026-08-25.md"), report)
	if _, err := ReadQAReport(specDir); err != nil {
		t.Fatalf("ReadQAReport of a multi-line refusal reason: %v", err)
	}
}

func TestWritePreconditionRefusalReportRecordsAnUnnamedRefusal(t *testing.T) {
	t.Parallel()
	var content strings.Builder
	if err := WritePreconditionRefusalReport(&content, PreconditionRefusal{CheckName: "  ", Reason: ""}); err != nil {
		t.Fatalf("WritePreconditionRefusalReport: %v", err)
	}
	report := content.String()
	if rows := qaResultsRows(t, report); len(rows) != 1 {
		t.Fatalf("Results rows = %q, want the refusal recorded even when it is unnamed", rows)
	}
	for _, want := range []string{
		`precondition_check: "` + QAPreconditionCheckUnnamed + `"`,
		`precondition_reason: "` + QAPreconditionReasonUnrecorded + `"`,
	} {
		if !strings.Contains(report, want+"\n") {
			t.Errorf("report does not record %q:\n%s", want, report)
		}
	}
}

func TestReadQAReportRejectsAPreconditionBlockedPass(t *testing.T) {
	t.Parallel()
	specDir := t.TempDir()
	writeFile(t, filepath.Join(specDir, "qa", "qa-report-2026-08-25.md"),
		qaReportFixture(VerdictPass, "rows_blocked_precondition: 1"))

	_, err := ReadQAReport(specDir)
	var reportErr QAReportError
	if !errors.As(err, &reportErr) {
		t.Fatalf("error = %v, want QAReportError; a gate that never ran cannot pass", err)
	}
	if !strings.Contains(err.Error(), "rows_blocked_precondition must be zero") {
		t.Errorf("error %q does not name the precondition count", err)
	}
}

func TestReadQAReportRejectsANegativePreconditionCount(t *testing.T) {
	t.Parallel()
	specDir := t.TempDir()
	writeFile(t, filepath.Join(specDir, "qa", "qa-report-2026-08-25.md"),
		qaReportFixture(VerdictFail, "rows_blocked_precondition: -1"))

	_, err := ReadQAReport(specDir)
	var reportErr QAReportError
	if !errors.As(err, &reportErr) {
		t.Fatalf("error = %v, want QAReportError", err)
	}
	if !strings.Contains(err.Error(), "rows_blocked_precondition must be a non-negative integer") {
		t.Errorf("error %q does not name the precondition count", err)
	}
}

func TestReadQAReportRecordsThePreconditionRefusal(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		written PreconditionRefusal
		want    PreconditionRefusal
	}{
		{
			name: "a named refusal reads back as it was written",
			written: PreconditionRefusal{
				CheckName: "spec check --strict",
				Reason:    `SC-VOCABULARY-UNDOCUMENTED: term "Run Ledger" is undocumented`,
			},
			want: PreconditionRefusal{
				CheckName: "spec check --strict",
				Reason:    `SC-VOCABULARY-UNDOCUMENTED: term "Run Ledger" is undocumented`,
			},
		},
		{
			name: "a refusal spread over lines reads back on one",
			written: PreconditionRefusal{
				CheckName: "spec check\t--strict",
				Reason:    "SC-REQUIREMENT-CONTRADICTORY:\n_prd.md requires one report per run and forbids writing one",
			},
			want: PreconditionRefusal{
				CheckName: "spec check --strict",
				Reason:    "SC-REQUIREMENT-CONTRADICTORY: _prd.md requires one report per run and forbids writing one",
			},
		},
		{
			name:    "an unnamed refusal reads back as the placeholder it recorded",
			written: PreconditionRefusal{},
			want: PreconditionRefusal{
				CheckName: QAPreconditionCheckUnnamed,
				Reason:    QAPreconditionReasonUnrecorded,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var content strings.Builder
			if err := WritePreconditionRefusalReport(&content, tt.written); err != nil {
				t.Fatalf("WritePreconditionRefusalReport: %v", err)
			}
			specDir := t.TempDir()
			writeFile(t, filepath.Join(specDir, "qa", "qa-report-2026-08-25.md"), content.String())

			report, err := ReadQAReport(specDir)
			if err != nil {
				t.Fatalf("ReadQAReport of a refusal the gate itself wrote: %v", err)
			}
			if report.Precondition != tt.want {
				t.Errorf("Precondition = %#v, want %#v", report.Precondition, tt.want)
			}
		})
	}
}

func TestPreconditionRefusalRoundTripsThroughTheQAReport(t *testing.T) {
	t.Parallel()
	refusal := PreconditionRefusal{
		CheckName: "spec check --strict",
		Reason:    `SC-VOCABULARY-UNDOCUMENTED: term "Run Ledger" is undocumented; SC-COVERAGE-UNMAPPED: Core Feature 2 has no TechSpec section`,
	}
	var written strings.Builder
	if err := WritePreconditionRefusalReport(&written, refusal); err != nil {
		t.Fatalf("WritePreconditionRefusalReport: %v", err)
	}
	specDir := t.TempDir()
	writeFile(t, filepath.Join(specDir, "qa", "qa-report-2026-08-25.md"), written.String())

	report, err := ReadQAReport(specDir)
	if err != nil {
		t.Fatalf("ReadQAReport: %v", err)
	}
	if report.Precondition != refusal {
		t.Fatalf("Precondition = %#v, want the refusal that was written, %#v", report.Precondition, refusal)
	}
	// The metadata a reader hands back has to be enough to write the same
	// refusal again; anything the read drops would be evidence the next writer
	// cannot recover.
	var rewritten strings.Builder
	if err := WritePreconditionRefusalReport(&rewritten, report.Precondition); err != nil {
		t.Fatalf("WritePreconditionRefusalReport from the report that was read: %v", err)
	}
	if rewritten.String() != written.String() {
		t.Errorf("rewritten report differs from the one read:\n--- read ---\n%s\n--- rewritten ---\n%s", written.String(), rewritten.String())
	}
}

func TestReadQAReportLeavesThePreconditionUnrecordedWhenNoneRefused(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name             string
		verdict          string
		extraFrontmatter []string
		want             PreconditionRefusal
	}{
		{
			name:    "a pass that names no precondition stays readable",
			verdict: VerdictPass,
			want:    PreconditionRefusal{},
		},
		{
			name:             "a check without a reason records the check alone",
			verdict:          VerdictFail,
			extraFrontmatter: []string{`precondition_check: "spec check --strict"`},
			want:             PreconditionRefusal{CheckName: "spec check --strict"},
		},
		{
			name:             "a reason without a check records the reason alone",
			verdict:          VerdictFail,
			extraFrontmatter: []string{`precondition_reason: "SC-REQUIREMENT-CONTRADICTORY"`},
			want:             PreconditionRefusal{Reason: "SC-REQUIREMENT-CONTRADICTORY"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			specDir := t.TempDir()
			writeFile(t, filepath.Join(specDir, "qa", "qa-report-2026-08-25.md"),
				qaReportFixture(tt.verdict, tt.extraFrontmatter...))

			report, err := ReadQAReport(specDir)
			if err != nil {
				t.Fatalf("ReadQAReport: %v", err)
			}
			if report.Verdict != tt.verdict {
				t.Errorf("Verdict = %q, want %q", report.Verdict, tt.verdict)
			}
			if report.Precondition != tt.want {
				t.Errorf("Precondition = %#v, want %#v", report.Precondition, tt.want)
			}
		})
	}
}

func TestReadQAReportRejectsANonScalarPrecondition(t *testing.T) {
	t.Parallel()
	for _, field := range []string{
		"precondition_check: [strict, vocabulary]",
		"precondition_reason:\n  code: SC-REQUIREMENT-CONTRADICTORY",
	} {
		t.Run(strings.SplitN(field, ":", 2)[0], func(t *testing.T) {
			specDir := t.TempDir()
			writeFile(t, filepath.Join(specDir, "qa", "qa-report-2026-08-25.md"),
				qaReportFixture(VerdictFail, field))

			_, err := ReadQAReport(specDir)
			var reportErr QAReportError
			if !errors.As(err, &reportErr) {
				t.Fatalf("error = %v, want QAReportError; a refusal that is not one recorded line is not readable", err)
			}
			if !strings.Contains(err.Error(), "cannot unmarshal") {
				t.Errorf("error %q does not name the unreadable value", err)
			}
		})
	}
}
