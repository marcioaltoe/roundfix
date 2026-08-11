// Package speccheck compares declarations and citations inside one Spec's
// artifacts. It reports contradictions without mutating the Spec or opening a
// network connection.
package speccheck

import (
	"encoding/json"
	"fmt"
	"io"
	"strconv"
	"strings"
)

// SchemaVersion identifies the machine-readable Spec Consistency Check
// report contract.
const SchemaVersion = "roundfix-speccheck/v1"

// Severity classifies whether a finding is settled or still needs judgment.
type Severity string

const (
	// SeverityError identifies a contradiction whose two sides are located.
	SeverityError Severity = "error"
	// SeverityGap identifies a candidate the checker cannot settle.
	SeverityGap Severity = "gap"
)

// Location identifies one side of a finding in repository-relative form.
type Location struct {
	Path string `json:"path"`
	Line int    `json:"line"`
}

// Finding is one stable diagnostic with the evidence and repair it needs.
type Finding struct {
	Code     string     `json:"code"`
	Severity Severity   `json:"severity"`
	Summary  string     `json:"summary"`
	Where    []Location `json:"where"`
	Fix      string     `json:"fix"`
}

// SkippedDetector records a detector that could not run because an input
// artifact was absent.
type SkippedDetector struct {
	Code    string `json:"code"`
	Missing string `json:"missing"`
}

// Result is the complete report for one Spec.
type Result struct {
	Slug     string            `json:"slug"`
	Findings []Finding         `json:"findings"`
	Skipped  []SkippedDetector `json:"skipped"`
}

// MechanicalResult is the complete, verdict-free output of the pre-QA
// mechanical stage. Findings are accumulated rather than returned fail-fast.
type MechanicalResult struct {
	Findings []MechanicalFinding
	Carried  []CarriedRow
	Blocked  []BlockedRow
	Skips    []MechanicalSkip
	Blocking bool
}

// MechanicalFinding locates one citation-checkable contradiction and its
// repair. RowHint names the QA matrix row it blocks when that relation is
// written down.
type MechanicalFinding struct {
	Code    string
	File    string
	Line    int
	Detail  string
	Fix     string
	RowHint string
}

// CarriedRow records the report and Git head that established a row. Task 05
// owns deciding whether a row qualifies for this state.
type CarriedRow struct {
	ID              string
	EstablishedBy   string
	EstablishedHead string
	Inputs          []EvidenceInput
}

// BlockedRow ties one report row to the mechanical Finding that stops it.
type BlockedRow struct {
	ID          string
	FindingCode string
	WaitingOn   string
}

// MechanicalSkip records a detector that could not run because its declared
// input artifact was absent.
type MechanicalSkip struct {
	Detector        string
	MissingArtifact string
}

// EvidenceInputKind is the closed vocabulary for a row's evidence source.
type EvidenceInputKind string

const (
	EvidenceRepositoryPath     EvidenceInputKind = "repository_path"
	EvidenceExternalRepository EvidenceInputKind = "external_repository"
	EvidenceLiveService        EvidenceInputKind = "live_service"
	EvidenceElapsedTime        EvidenceInputKind = "elapsed_time"
)

// EvidenceInput names one observation boundary. Only repository paths can be
// eligible for carry-forward, whose decision belongs to Task 05.
type EvidenceInput struct {
	Kind EvidenceInputKind
	Ref  string
}

// EvidenceFile identifies one tracked Git blob by repository-relative path and
// the SHA-256 digest of its bytes.
type EvidenceFile struct {
	Path   string
	SHA256 string
}

// EvidenceSnapshot is the canonical expansion of one repository_path input at
// one Git head. Files must be sorted by path and contain no duplicates.
type EvidenceSnapshot struct {
	Ref   string
	Files []EvidenceFile
}

// ReportRow contains the carry-forward facts for one prior QA Report row.
// AncestryVerified is set only by the repository resolver after Git proves
// EstablishedHead is an ancestor of the current head.
type ReportRow struct {
	ID               string
	Status           string
	EstablishedBy    string
	EstablishedHead  string
	AncestryVerified bool
	Inputs           []EvidenceInput
	EvidencePaths    []string
}

type jsonDocument struct {
	Schema   string            `json:"schema"`
	Slug     string            `json:"slug"`
	Findings []Finding         `json:"findings"`
	Skipped  []SkippedDetector `json:"skipped"`
}

// RenderText renders one result with every diagnostic location and fix.
func RenderText(result Result) string {
	var report strings.Builder
	report.WriteString("Spec ")
	report.WriteString(result.Slug)
	report.WriteByte('\n')
	if len(result.Findings) == 0 {
		report.WriteString("No findings.\n")
	}
	for _, finding := range result.Findings {
		report.WriteByte('[')
		report.WriteString(string(finding.Severity))
		report.WriteString("] ")
		report.WriteString(finding.Code)
		report.WriteString(": ")
		report.WriteString(finding.Summary)
		report.WriteByte('\n')
		for _, location := range finding.Where {
			report.WriteString("  at ")
			report.WriteString(location.Path)
			report.WriteByte(':')
			report.WriteString(strconv.Itoa(location.Line))
			report.WriteByte('\n')
		}
		report.WriteString("  fix: ")
		report.WriteString(finding.Fix)
		report.WriteByte('\n')
	}
	if len(result.Skipped) > 0 {
		report.WriteString("Skipped:\n")
		for _, skipped := range result.Skipped {
			report.WriteString("  ")
			report.WriteString(skipped.Code)
			report.WriteString(": missing ")
			report.WriteString(skipped.Missing)
			report.WriteByte('\n')
		}
	}
	return report.String()
}

// RenderJSON renders one compact roundfix-speccheck/v1 object.
func RenderJSON(result Result) ([]byte, error) {
	findings := result.Findings
	if findings == nil {
		findings = []Finding{}
	}
	skipped := result.Skipped
	if skipped == nil {
		skipped = []SkippedDetector{}
	}
	document := jsonDocument{
		Schema:   SchemaVersion,
		Slug:     result.Slug,
		Findings: findings,
		Skipped:  skipped,
	}
	data, err := json.Marshal(document)
	if err != nil {
		return nil, fmt.Errorf("render Spec Consistency Check JSON: %w", err)
	}
	return data, nil
}

// WriteMechanicalResult writes the mechanical sections and row seeds that a
// Daemon-owned report lifecycle can place in a new QA Report. It does not
// choose a report path, edit an existing report, or compute a verdict.
func WriteMechanicalResult(writer io.Writer, result MechanicalResult) error {
	var report strings.Builder
	report.WriteString("## Mechanical findings\n\n")
	if len(result.Findings) == 0 {
		report.WriteString("None.\n\n")
	}
	for _, finding := range result.Findings {
		report.WriteString("### ")
		report.WriteString(markdownText(finding.Code))
		report.WriteString("\n\n- location: `")
		report.WriteString(markdownCode(finding.File))
		report.WriteByte(':')
		report.WriteString(strconv.Itoa(finding.Line))
		report.WriteString("`\n- detail: ")
		report.WriteString(markdownText(finding.Detail))
		report.WriteString("\n- fix: ")
		report.WriteString(markdownText(finding.Fix))
		if finding.RowHint != "" {
			report.WriteString("\n- blocked row: `")
			report.WriteString(markdownCode(finding.RowHint))
			report.WriteByte('`')
		}
		report.WriteString("\n\n")
	}

	report.WriteString("## Mechanical rows\n\n")
	report.WriteString("| # | Status | Provenance |\n| - | --- | --- |\n")
	for _, row := range result.Carried {
		status := "carried (established by: " + row.EstablishedBy + "; head: " + row.EstablishedHead + ")"
		report.WriteString("| ")
		report.WriteString(markdownCell(row.ID))
		report.WriteString(" | ")
		report.WriteString(markdownCell(status))
		report.WriteString(" | report and head retained |\n")
	}
	for _, row := range result.Blocked {
		status := "blocked (finding: " + row.FindingCode + " — waits on " + row.WaitingOn + ")"
		report.WriteString("| ")
		report.WriteString(markdownCell(row.ID))
		report.WriteString(" | ")
		report.WriteString(markdownCell(status))
		report.WriteString(" | mechanical finding |\n")
	}

	report.WriteString("\n## Mechanical skips\n\n")
	if len(result.Skips) == 0 {
		report.WriteString("None.\n")
	} else {
		report.WriteString("| Detector | Missing artifact |\n| --- | --- |\n")
		for _, skip := range result.Skips {
			report.WriteString("| ")
			report.WriteString(markdownCell(skip.Detector))
			report.WriteString(" | ")
			report.WriteString(markdownCell(skip.MissingArtifact))
			report.WriteString(" |\n")
		}
	}

	if _, err := io.WriteString(writer, report.String()); err != nil {
		return fmt.Errorf("write mechanical QA Report sections: %w", err)
	}
	return nil
}

func markdownText(value string) string {
	return strings.ReplaceAll(strings.TrimSpace(value), "\n", " ")
}

func markdownCode(value string) string {
	return strings.ReplaceAll(markdownText(value), "`", "'")
}

func markdownCell(value string) string {
	return strings.ReplaceAll(markdownText(value), "|", "\\|")
}
