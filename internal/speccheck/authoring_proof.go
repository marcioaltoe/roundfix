package speccheck

import (
	"fmt"
	"regexp"
	"strings"
)

const (
	// CodeAbsenceUnproven identifies a TechSpec without its absence proof.
	CodeAbsenceUnproven = "SC-ABSENCE-UNPROVEN"
	// CodeADRAlternativesUnrelated identifies a listed ADR without its alternatives line.
	CodeADRAlternativesUnrelated = "SC-ADR-ALTERNATIVES-UNRELATED"
	// CodeCoverageUnread identifies a Coverage Map without a declared cross-read.
	CodeCoverageUnread = "SC-COVERAGE-UNREAD"
)

var (
	techSpecProofCodes    = []string{CodeAbsenceUnproven, CodeADRAlternativesUnrelated, CodeCoverageUnread}
	absenceProofPattern   = regexp.MustCompile("command:\\s*`([^`]+)`\\s*;\\s*result:\\s*(\\S[\\s\\S]*)")
	adrAlternativePattern = regexp.MustCompile(`^\s*-\s+ADR-([0-9]{4}):\s*(\S.*)$`)
)

// detectTechSpecProof reads declarations only; it never executes an authored command
// or exempts an active Spec by its commit date.
func detectTechSpecProof(result *Result, prdContent, techSpecContent []byte, prdDisplay, techSpecDisplay string) {
	prdContent = proofOutsideFences(prdContent)
	techSpecContent = proofOutsideFences(techSpecContent)
	detectAbsenceEvidence(result, techSpecContent, techSpecDisplay)
	detectADRAlternatives(result, prdContent, techSpecContent, prdDisplay, techSpecDisplay)
	detectCoverageCrossRead(result, techSpecContent, techSpecDisplay)
}

// Keep line positions while removing examples, including headings inside fences.
func proofOutsideFences(content []byte) []byte {
	lines := strings.Split(string(content), "\n")
	fence := ""
	for index, line := range lines {
		if fence != "" {
			if closesTranscriptFence(line, fence) {
				fence = ""
			}
			lines[index] = ""
			continue
		}
		if marker, _ := transcriptFence(line); marker != "" {
			fence = marker
			lines[index] = ""
		}
	}
	return []byte(strings.Join(lines, "\n"))
}

func detectAbsenceEvidence(result *Result, content []byte, display string) {
	const section = "Absence Evidence"
	const fix = "Add or repair ## Absence Evidence with None. and a reason, or numbered items with command: in backticks and result: text, following the write-techspec rule \"Absence evidence\"."
	lines := markdownSectionLines(content, section)
	for len(lines) > 0 && strings.TrimSpace(lines[0]) == "" {
		lines = lines[1:]
	}
	if len(lines) > 0 && strings.HasPrefix(strings.TrimSpace(lines[0]), "None.") {
		end := 0
		for end < len(lines) && strings.TrimSpace(lines[end]) != "" {
			end++
		}
		if validNoneParagraph(lines[:end]) {
			return
		}
	}

	items := 0
	for index := 0; index < len(lines); index++ {
		match := numberedItemPattern.FindStringSubmatch(lines[index])
		if len(match) != 2 {
			continue
		}
		items++
		itemLine := sectionLineContaining(content, section, lines[index])
		item := []string{lines[index]}
		for index+1 < len(lines) && strings.TrimSpace(lines[index+1]) != "" && !markdownBlockLine(strings.TrimSpace(lines[index+1])) {
			index++
			item = append(item, lines[index])
		}
		proof := absenceProofPattern.FindStringSubmatch(strings.Join(item, "\n"))
		if len(proof) == 3 && strings.TrimSpace(proof[1]) != "" && strings.TrimSpace(proof[2]) != "" {
			continue
		}
		result.Findings = append(result.Findings, Finding{
			Code: CodeAbsenceUnproven, Severity: SeverityGap,
			Summary: fmt.Sprintf("%s Absence Evidence item %s has no command or no result", display, match[1]),
			Where:   []Location{{Path: display, Line: itemLine}}, Fix: fix,
		})
	}
	if items == 0 {
		result.Findings = append(result.Findings, Finding{
			Code: CodeAbsenceUnproven, Severity: SeverityGap,
			Summary: display + " has no Absence Evidence declaration",
			Where:   []Location{{Path: display, Line: sectionLineContaining(content, section, "None.")}}, Fix: fix,
		})
	}
}

func detectADRAlternatives(result *Result, prdContent, techSpecContent []byte, prdDisplay, techSpecDisplay string) {
	listed := make(map[string]Location)
	for _, artifact := range []struct {
		content []byte
		display string
	}{{prdContent, prdDisplay}, {techSpecContent, techSpecDisplay}} {
		row, line, _ := projectConstraintEntry(artifact.content, "Active ADR obligations")
		for number := range obligationCitationNumbers([]byte(row)) {
			if _, present := listed[number]; !present {
				listed[number] = Location{Path: artifact.display, Line: line}
			}
		}
	}
	related := make(map[string]bool)
	for _, line := range markdownSectionLines(techSpecContent, "ADR Alternatives") {
		if match := adrAlternativePattern.FindStringSubmatch(line); len(match) == 3 {
			related[match[1]] = true
		}
	}
	for _, number := range sortedCitationNumbers(listed) {
		if related[number] {
			continue
		}
		location := listed[number]
		result.Findings = append(result.Findings, Finding{
			Code: CodeADRAlternativesUnrelated, Severity: SeverityGap,
			Summary: location.Path + " lists ADR-" + number + " under Active ADR obligations, but " + techSpecDisplay + " has no ADR Alternatives line for it",
			Where:   []Location{location, {Path: techSpecDisplay, Line: sectionLineContaining(techSpecContent, "ADR Alternatives", "ADR-"+number+":")}},
			Fix:     "Add - ADR-" + number + ": <relation> under ## ADR Alternatives, following the write-techspec rule \"Rejected alternatives\".",
		})
	}
}

func detectCoverageCrossRead(result *Result, content []byte, display string) {
	const section = "Coverage Map"
	const prefix = "Deviations from the PRD:"
	const fix = "Add or repair the Deviations from the PRD: line in ## Coverage Map with None. or bullets naming PRD units, following the write-techspec rule \"Coverage cross-read\"."
	lines := markdownSectionLines(content, section)
	declaration := -1
	for index, line := range lines {
		if strings.HasPrefix(strings.TrimSpace(line), prefix) {
			declaration = index
			break
		}
	}
	line := sectionLineContaining(content, section, prefix)
	// None. ends the declaration, but a deviation bullet after it must still name a PRD unit.
	none := declaration >= 0 && strings.TrimSpace(strings.TrimPrefix(strings.TrimSpace(lines[declaration]), prefix)) == "None."
	bullets := 0
	if declaration >= 0 {
		for index := declaration + 1; index < len(lines); index++ {
			trimmed := strings.TrimSpace(lines[index])
			if !strings.HasPrefix(trimmed, "- ") && !strings.HasPrefix(trimmed, "* ") && !strings.HasPrefix(trimmed, "+ ") {
				continue
			}
			bullets++
			references := newReferenceSet()
			addDeclaredReferences(references, trimmed)
			if len(references[coverageStory])+len(references[coverageFeature])+len(references[coverageMetric]) > 0 {
				continue
			}
			result.Findings = append(result.Findings, Finding{
				Code: CodeCoverageUnread, Severity: SeverityGap,
				Summary: display + " Coverage Cross-Read deviation names no PRD unit",
				Where:   []Location{{Path: display, Line: line + index - declaration}}, Fix: fix,
			})
		}
	}
	if bullets == 0 && !none {
		summary := display + " Coverage Map has no Deviations from the PRD: line"
		if declaration >= 0 {
			summary = display + " Deviations from the PRD: line has neither None. nor a following bullet"
		}
		result.Findings = append(result.Findings, Finding{
			Code: CodeCoverageUnread, Severity: SeverityGap, Summary: summary,
			Where: []Location{{Path: display, Line: line}}, Fix: fix,
		})
	}
}
