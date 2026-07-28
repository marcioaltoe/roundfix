package spec

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

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

// NewestQAReport returns the path of the newest qa/qa-report-*.md in the
// Spec directory. Report names are qa-report-YYYY-MM-DD.md for the first run
// of a date and qa-report-YYYY-MM-DD-NN.md for same-date reruns, so recency
// is the embedded date first and then the run sequence compared as a number:
// an unsuffixed report counts as sequence 1, which makes -02 newer than it
// and -10 newer than -02. Raw byte order would invert that pair, because "-"
// sorts before ".". A name whose suffix is not a number (the qa-gate's
// scope-or-build variant) carries no sequence, so it loses to every
// sequenced report of the same date but still beats every earlier date. A
// name with no valid date loses to every dated report and wins only when no
// report parses at all. Equal keys fall back to path order, so the result is
// always deterministic. Missing reports surface as ErrNoQAReport.
func NewestQAReport(specDir string) (string, error) {
	pattern := filepath.Join(specDir, "qa", "qa-report-*.md")
	reports, err := filepath.Glob(pattern)
	if err != nil {
		return "", fmt.Errorf("find QA Reports matching %q: %w", pattern, err)
	}
	if len(reports) == 0 {
		return "", fmt.Errorf("%w for Spec directory %q; run the qa-gate workflow to produce one", ErrNoQAReport, specDir)
	}
	newest := parseQAReportName(reports[0])
	for _, report := range reports[1:] {
		if candidate := parseQAReportName(report); newest.olderThan(candidate) {
			newest = candidate
		}
	}
	return newest.path, nil
}

// qaReportName is the recency key of one qa/qa-report-*.md path.
type qaReportName struct {
	path string
	// date is the embedded YYYY-MM-DD, empty when the name carries no valid
	// date. The fixed-width form makes byte order date order.
	date string
	// sequenced reports the name carries a run sequence: absent (the first
	// run of the date) or an all-digit -NN suffix.
	sequenced bool
	// sequence is the run number within the date, 1 for an unsuffixed name.
	sequence int
}

// qaReportNamePrefix and qaReportNameSuffix bracket the date and the optional
// run suffix in a report file name.
const (
	qaReportNamePrefix = "qa-report-"
	qaReportNameSuffix = ".md"
	qaReportDateLayout = "2006-01-02"
)

// parseQAReportName derives the recency key of a report path. Names the
// qa-gate never writes stay parseable enough to order: an unrecognized run
// suffix drops the sequence, and an unrecognized date drops the date.
func parseQAReportName(path string) qaReportName {
	parsed := qaReportName{path: path}
	rest := filepath.Base(path)
	rest = strings.TrimPrefix(rest, qaReportNamePrefix)
	rest = strings.TrimSuffix(rest, qaReportNameSuffix)
	if len(rest) < len(qaReportDateLayout) {
		return parsed
	}
	date := rest[:len(qaReportDateLayout)]
	if _, err := time.Parse(qaReportDateLayout, date); err != nil {
		return parsed
	}
	parsed.date = date

	switch suffix := rest[len(qaReportDateLayout):]; {
	case suffix == "":
		parsed.sequenced, parsed.sequence = true, 1
	case strings.HasPrefix(suffix, "-") && isDigits(suffix[1:]):
		sequence, err := strconv.Atoi(suffix[1:])
		if err != nil {
			return parsed
		}
		parsed.sequenced, parsed.sequence = true, sequence
	}
	return parsed
}

// olderThan reports whether name is less recent than other, ordering by date,
// then by presence of a run sequence, then by that sequence as a number, and
// finally by path so equal keys stay deterministic.
func (name qaReportName) olderThan(other qaReportName) bool {
	if name.date != other.date {
		return name.date < other.date
	}
	if name.sequenced != other.sequenced {
		return !name.sequenced
	}
	if name.sequenced && name.sequence != other.sequence {
		return name.sequence < other.sequence
	}
	return name.path < other.path
}

// isDigits reports whether value is a non-empty run of ASCII digits, the only
// form a QA Report run sequence takes.
func isDigits(value string) bool {
	if value == "" {
		return false
	}
	for _, char := range value {
		if char < '0' || char > '9' {
			return false
		}
	}
	return true
}

// QAVerdict reads the verdict from the newest qa/qa-report-*.md frontmatter
// in the Spec directory. Missing reports surface as ErrNoQAReport; reports
// whose verdict cannot be read surface as QAReportError.
func QAVerdict(specDir string) (string, error) {
	newest, err := NewestQAReport(specDir)
	if err != nil {
		return "", err
	}

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
