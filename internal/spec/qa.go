package spec

import (
	"errors"
	"fmt"
	"io"
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

// The terminal row a gate writes when it refuses at a precondition check. The
// gate stopped before it built its matrix, so it has no measured requirement
// to report; row 0 records the refusal itself, which keeps the Results table
// non-empty and names the stop as a refusal rather than as a measurement.
const (
	QAPreconditionRowID         = "0"
	QAPreconditionRowStatus     = "blocked"
	QAPreconditionRowProvenance = "precondition"
)

// The recorded stand-ins for a refusal whose cause the caller could not name.
// A refusal that cannot be written is the deadlock this shape exists to end,
// so an unnamed cause degrades to a placeholder a reader can see rather than
// to no report at all.
const (
	QAPreconditionCheckUnnamed     = "unnamed precondition check"
	QAPreconditionReasonUnrecorded = "reason not recorded"
)

// ErrNoQAReport reports that the Spec has no qa/qa-report-*.md yet.
var ErrNoQAReport = errors.New("no QA Report found")

// QAReport is the validated verdict metadata callers need from a QA Report.
// Each blocked-row cause remains independent; absent counts are zero.
type QAReport struct {
	Verdict                 string
	RowsBlockedEnvironment  int
	RowsBlockedFinding      int
	RowsBlockedDeclared     int
	RowsBlockedPrecondition int
	// Precondition is the refusal the report records, zero when the gate
	// reached its matrix. It is optional metadata: a report that names no
	// precondition is one that refused for none, not a report that is missing
	// something, so its absence never makes a verdict unreadable.
	Precondition PreconditionRefusal
}

// PreconditionRefusal is one QA gate stop at a precondition check: the name of
// the check that refused and the reason it gave. Both reach the report so a
// reader knows what was checked and why the gate never built its matrix, and
// both come back out of it, so the refusal a reader recovers is the one the
// gate recorded rather than a summary of it.
type PreconditionRefusal struct {
	CheckName string
	Reason    string
}

// WritePreconditionRefusalReport writes the whole QA Report of a gate that
// refused at a precondition check: verdict fail, one blocked-precondition row
// in the Results table, and the check and its reason in the frontmatter. The
// caller owns where the report lands and under which name; this contract owns
// its shape, so a refusal cannot be recorded as an empty matrix again.
func WritePreconditionRefusalReport(writer io.Writer, refusal PreconditionRefusal) error {
	check := qaRefusalValue(refusal.CheckName, QAPreconditionCheckUnnamed)
	reason := qaRefusalValue(refusal.Reason, QAPreconditionReasonUnrecorded)

	var report strings.Builder
	report.WriteString("---\n")
	fmt.Fprintf(&report, "verdict: %s\n", VerdictFail)
	// One refusal is one row, and every other cause stays explicitly zero:
	// a closed report carries all its typed counts even when they are empty.
	report.WriteString("rows_blocked_precondition: 1\n")
	report.WriteString("rows_blocked_environment: 0\n")
	report.WriteString("rows_blocked_finding: 0\n")
	report.WriteString("rows_blocked_declared: 0\n")
	fmt.Fprintf(&report, "precondition_check: %s\n", strconv.Quote(check))
	fmt.Fprintf(&report, "precondition_reason: %s\n", strconv.Quote(reason))
	report.WriteString("---\n\n# QA Report\n\n")

	report.WriteString("## Results\n\n")
	report.WriteString("| # | Status | Provenance |\n| - | --- | --- |\n")
	fmt.Fprintf(&report, "| %s | %s | %s |\n\n", QAPreconditionRowID, QAPreconditionRowStatus, QAPreconditionRowProvenance)

	// A list, never a second table: a markdown table outside the Results
	// section is read as more result rows, so prose that justifies the row
	// would become rows of its own.
	report.WriteString("## Precondition refusal\n\n")
	fmt.Fprintf(&report, "- check: %s\n", check)
	fmt.Fprintf(&report, "- reason: %s\n\n", reason)
	report.WriteString("The gate stopped at this check before it built its QA matrix, so no requirement was measured and the row above records the refusal itself.\n")

	if _, err := io.WriteString(writer, report.String()); err != nil {
		return fmt.Errorf("write precondition refusal QA Report: %w", err)
	}
	return nil
}

// qaRefusalValue names one refusal value, falling back when the caller named
// nothing.
func qaRefusalValue(value, fallback string) string {
	if line := qaRefusalLine(value); line != "" {
		return line
	}
	return fallback
}

// qaRefusalLine collapses a refusal value onto one line, so that writing it
// cannot break the frontmatter it goes into and reading it back yields the
// same one line the writer would have produced for it.
func qaRefusalLine(value string) string {
	return strings.Join(strings.Fields(value), " ")
}

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
	return NewestQAReportFromPaths(reports)
}

// NewestQAReportFromPaths returns the newest path from a non-empty set of QA
// Report paths using the same filename recency contract as NewestQAReport.
// Callers that obtain paths from a Git tree use this helper so filesystem and
// historical-tree selection cannot drift.
func NewestQAReportFromPaths(reports []string) (string, error) {
	if len(reports) == 0 {
		return "", ErrNoQAReport
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
		// The day's first report sorts below every numeric sibling. It used to
		// be sequence 1, which `-01` also parses to, and the path-order
		// tie-break then picked the unsuffixed report over a newer one. -1
		// keeps that impossible for any suffix the contract can produce,
		// including `-00`.
		parsed.sequenced, parsed.sequence = true, -1
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

// QAVerdict reads and validates the verdict metadata from the newest
// qa/qa-report-*.md frontmatter in the Spec directory. Missing reports surface
// as ErrNoQAReport; unreadable or inconsistent metadata surfaces as
// QAReportError.
func QAVerdict(specDir string) (string, error) {
	report, err := ReadQAReport(specDir)
	if err != nil {
		return "", err
	}
	return report.Verdict, nil
}

// ReadQAReport reads and validates the verdict and typed blocked-row counts
// from the newest qa/qa-report-*.md in the Spec directory.
func ReadQAReport(specDir string) (QAReport, error) {
	newest, err := NewestQAReport(specDir)
	if err != nil {
		return QAReport{}, err
	}
	return readQAReport(newest)
}

func readQAReport(path string) (QAReport, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return QAReport{}, QAReportError{Path: path, Err: err}
	}
	frontmatterBytes, _, err := splitFrontmatter(content)
	if err != nil {
		return QAReport{}, QAReportError{Path: path, Err: err}
	}
	var frontmatter struct {
		Verdict                 string    `yaml:"verdict"`
		RowsBlockedEnvironment  yaml.Node `yaml:"rows_blocked_environment"`
		RowsBlockedFinding      yaml.Node `yaml:"rows_blocked_finding"`
		RowsBlockedDeclared     yaml.Node `yaml:"rows_blocked_declared"`
		RowsBlockedPrecondition yaml.Node `yaml:"rows_blocked_precondition"`
		PreconditionCheck       string    `yaml:"precondition_check"`
		PreconditionReason      string    `yaml:"precondition_reason"`
	}
	if err := yaml.Unmarshal(frontmatterBytes, &frontmatter); err != nil {
		return QAReport{}, QAReportError{Path: path, Err: err}
	}
	rowsBlockedEnvironment, err := qaBlockedCount(frontmatter.RowsBlockedEnvironment, "rows_blocked_environment")
	if err != nil {
		return QAReport{}, QAReportError{Path: path, Err: err}
	}
	rowsBlockedFinding, err := qaBlockedCount(frontmatter.RowsBlockedFinding, "rows_blocked_finding")
	if err != nil {
		return QAReport{}, QAReportError{Path: path, Err: err}
	}
	rowsBlockedDeclared, err := qaBlockedCount(frontmatter.RowsBlockedDeclared, "rows_blocked_declared")
	if err != nil {
		return QAReport{}, QAReportError{Path: path, Err: err}
	}
	rowsBlockedPrecondition, err := qaBlockedCount(frontmatter.RowsBlockedPrecondition, "rows_blocked_precondition")
	if err != nil {
		return QAReport{}, QAReportError{Path: path, Err: err}
	}
	// The refusal is read exactly as the report recorded it: an absent check or
	// reason stays empty rather than becoming the writer's placeholder, so a
	// reader can tell a refusal whose cause was unnamed from one that was never
	// written down.
	report := QAReport{
		Verdict:                 frontmatter.Verdict,
		RowsBlockedEnvironment:  rowsBlockedEnvironment,
		RowsBlockedFinding:      rowsBlockedFinding,
		RowsBlockedDeclared:     rowsBlockedDeclared,
		RowsBlockedPrecondition: rowsBlockedPrecondition,
		Precondition: PreconditionRefusal{
			CheckName: qaRefusalLine(frontmatter.PreconditionCheck),
			Reason:    qaRefusalLine(frontmatter.PreconditionReason),
		},
	}
	switch report.Verdict {
	case VerdictPass:
		if report.RowsBlockedFinding > 0 {
			return QAReport{}, QAReportError{
				Path: path,
				Err:  fmt.Errorf("rows_blocked_finding must be zero when verdict is %q", VerdictPass),
			}
		}
		if report.RowsBlockedDeclared > 0 {
			return QAReport{}, QAReportError{
				Path: path,
				Err:  fmt.Errorf("rows_blocked_declared must be zero when verdict is %q", VerdictPass),
			}
		}
		// A gate that refused at a precondition measured no requirement, so a
		// pass it claims describes nothing that ran.
		if report.RowsBlockedPrecondition > 0 {
			return QAReport{}, QAReportError{
				Path: path,
				Err:  fmt.Errorf("rows_blocked_precondition must be zero when verdict is %q", VerdictPass),
			}
		}
		return report, nil
	case VerdictFail, VerdictPartial:
		return report, nil
	case "":
		return QAReport{}, QAReportError{Path: path, Err: errors.New("frontmatter has no verdict field")}
	default:
		return QAReport{}, QAReportError{Path: path, Err: fmt.Errorf("unsupported verdict %q", report.Verdict)}
	}
}

// qaBlockedCount reads one optional typed blocked-cause count. An absent field
// is zero; a present field must be a non-negative YAML integer scalar.
func qaBlockedCount(node yaml.Node, field string) (int, error) {
	if node.Kind == 0 {
		return 0, nil
	}
	if node.Kind != yaml.ScalarNode || node.Tag != "!!int" {
		return 0, fmt.Errorf("cannot unmarshal %s as a non-negative integer", field)
	}
	var count int
	if err := node.Decode(&count); err != nil {
		return 0, fmt.Errorf("cannot unmarshal %s as a non-negative integer: %w", field, err)
	}
	if count < 0 {
		return 0, fmt.Errorf("%s must be a non-negative integer", field)
	}
	return count, nil
}
