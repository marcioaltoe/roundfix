package spec

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"gopkg.in/yaml.v3"
)

// QA Report verdict values written by the qa-gate workflow.
const (
	VerdictPass    = "pass"
	VerdictFail    = "fail"
	VerdictPartial = "partial"
)

// ErrNoQAReport reports that the Spec has no qa/qa-report-*.md yet.
var ErrNoQAReport = errors.New("no QA Report found")

// QAVerdict reads the verdict from the newest qa/qa-report-*.md frontmatter
// in the Spec directory. Missing reports surface as ErrNoQAReport; reports
// whose verdict cannot be read surface as QAReportError.
func QAVerdict(specDir string) (string, error) {
	pattern := filepath.Join(specDir, "qa", "qa-report-*.md")
	reports, err := filepath.Glob(pattern)
	if err != nil {
		return "", fmt.Errorf("find QA Reports matching %q: %w", pattern, err)
	}
	if len(reports) == 0 {
		return "", fmt.Errorf("%w for Spec directory %q; run the qa-gate workflow to produce one", ErrNoQAReport, specDir)
	}
	// Report names embed the date as YYYY-MM-DD, so the lexicographically
	// greatest path is the newest report.
	sort.Strings(reports)
	newest := reports[len(reports)-1]

	content, err := os.ReadFile(newest)
	if err != nil {
		return "", QAReportError{Path: newest, Err: err}
	}
	frontmatterBytes, _, err := splitFrontmatter(content)
	if err != nil {
		return "", QAReportError{Path: newest, Err: err}
	}
	var report struct {
		Verdict string `yaml:"verdict"`
	}
	if err := yaml.Unmarshal(frontmatterBytes, &report); err != nil {
		return "", QAReportError{Path: newest, Err: err}
	}
	switch report.Verdict {
	case VerdictPass, VerdictFail, VerdictPartial:
		return report.Verdict, nil
	case "":
		return "", QAReportError{Path: newest, Err: errors.New("frontmatter has no verdict field")}
	default:
		return "", QAReportError{Path: newest, Err: fmt.Errorf("unsupported verdict %q", report.Verdict)}
	}
}
