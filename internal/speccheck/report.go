// Package speccheck compares declarations and citations inside one Spec's
// artifacts. It reports contradictions without mutating the Spec or opening a
// network connection.
package speccheck

import (
	"encoding/json"
	"fmt"
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
