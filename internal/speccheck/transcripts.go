package speccheck

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"roundfix/internal/spec"
)

// Transcript is one declared Surface Transcript. Malformed is empty for a valid block.
type Transcript struct {
	Number    int
	Title     string
	Line      int
	Command   string
	Stdout    []string
	Stderr    []string
	ExitCode  int
	Malformed string
}

var transcriptItemPattern = regexp.MustCompile(`^\s*([0-9]+)\.\s+Surface Transcript:\s*(.*)$`)
var transcriptRefPattern = regexp.MustCompile(`(?i)\bSurface Transcripts?\s+`)

// SurfaceTranscripts reads declarations and blocks without executing any command.
func SurfaceTranscripts(content []byte) []Transcript {
	_, transcripts := parseTranscriptSection(content)
	return transcripts
}

// transcriptFence recognizes a Markdown fence and its exact info string.
func transcriptFence(line string) (string, string) {
	trimmed := strings.TrimSpace(line)
	if len(trimmed) < 3 || trimmed[0] != '`' && trimmed[0] != '~' {
		return "", ""
	}
	n := 0
	for n < len(trimmed) && trimmed[n] == trimmed[0] {
		n++
	}
	if n < 3 {
		return "", ""
	}
	return trimmed[:n], strings.TrimSpace(trimmed[n:])
}
func closesTranscriptFence(line, fence string) bool {
	marker, info := transcriptFence(line)
	return marker != "" && marker[0] == fence[0] && len(marker) >= len(fence) && info == ""
}

func parseTranscriptSection(content []byte) (promiseDeclaration, []Transcript) {
	lines := strings.Split(string(content), "\n")
	declaration := promiseDeclaration{state: promiseNoDeclaration, line: 1}
	level := 0
	fence := ""
	blockTranscript := false
	blockIndent := ""
	var section []string
	var transcripts []Transcript
	var blocks [][]string
	var block []string
	finish := func() {
		if len(transcripts) == 0 {
			return
		}
		item := &transcripts[len(transcripts)-1]
		switch len(blocks) {
		case 0:
			item.Malformed = "without a transcript block"
		case 1:
			parseTranscriptBlock(item, blocks[0])
		default:
			item.Malformed = "with more than one transcript block"
		}
		blocks = nil
	}
	for index, line := range lines {
		if fence != "" {
			if closesTranscriptFence(line, fence) {
				if blockTranscript {
					blocks = append(blocks, block)
				}
				fence = ""
				blockTranscript = false
				block = nil
			} else if blockTranscript {
				block = append(block, strings.TrimPrefix(line, blockIndent))
			}
			if level != 0 {
				section = append(section, line)
			}
			continue
		}
		marker, info := transcriptFence(line)
		if marker != "" {
			fence = marker
			blockTranscript = level != 0 && len(transcripts) > 0 && info == "transcript"
			blockIndent = line[:len(line)-len(strings.TrimLeft(line, " \t"))]
			if level != 0 {
				section = append(section, line)
			}
			continue
		}
		heading, title := markdownHeading(line)
		if level == 0 {
			if (heading == 2 || heading == 3) && title == "Surface Transcripts" {
				level = heading
				declaration.line = index + 1
			}
			continue
		}
		if heading > 0 && heading <= level {
			break
		}
		section = append(section, line)
		match := transcriptItemPattern.FindStringSubmatch(line)
		if len(match) == 0 {
			continue
		}
		number, err := strconv.Atoi(match[1])
		if err != nil {
			continue
		}
		finish()
		transcripts = append(transcripts, Transcript{Number: number, Title: match[2], Line: index + 1})
	}
	// An unclosed fence cannot satisfy the complete four-line block contract.
	if blockTranscript {
		blocks = append(blocks, nil)
	}
	finish()
	if len(transcripts) > 0 {
		declaration.state = promiseDeclaredUnits
		for _, item := range transcripts {
			declaration.units = append(declaration.units, coverageUnit{Kind: coverageTranscript, Number: item.Number, Line: item.Line})
		}
	} else {
		for index, line := range section {
			trimmed := strings.TrimSpace(line)
			if trimmed == "" {
				continue
			}
			if strings.HasPrefix(trimmed, "None.") && validNoneParagraph(section[index:]) {
				declaration.state = promiseExplicitNone
			}
			break
		}
	}
	return declaration, transcripts
}

func parseTranscriptBlock(item *Transcript, lines []string) {
	if len(lines) == 0 || !strings.HasPrefix(lines[0], "$ ") || strings.TrimSpace(strings.TrimPrefix(lines[0], "$ ")) == "" {
		item.Malformed = "without a command line"
		return
	}
	item.Command = strings.TrimPrefix(lines[0], "$ ")
	if len(lines) < 2 || lines[1] != "stdout:" {
		item.Malformed = "without a stdout line"
		return
	}
	stderr := -1
	for index := 2; index < len(lines); index++ {
		if lines[index] == "stderr:" {
			stderr = index
			break
		}
	}
	if stderr == -1 {
		item.Malformed = "without a stderr line after its stdout line"
		return
	}
	item.Stdout = append([]string(nil), lines[2:stderr]...)
	if len(lines) <= stderr+1 || !strings.HasPrefix(lines[len(lines)-1], "exit: ") {
		item.Malformed = "without an exit line"
		return
	}
	code, err := strconv.Atoi(strings.TrimPrefix(lines[len(lines)-1], "exit: "))
	if err != nil || code < 0 || code > 255 {
		item.Malformed = "with an exit code outside 0-255"
		return
	}
	item.Stderr = append([]string(nil), lines[stderr+1:len(lines)-1]...)
	item.ExitCode = code
}

func detectTranscripts(result *Result, horizon contractHorizon, content []byte, artifact string) []coverageUnit {
	declaration, transcripts := parseTranscriptSection(content)
	if horizon.held {
		detectPromiseDeclaration(result, declaration, artifact, "Surface Transcripts", CodeTranscriptUndeclared)
	} else {
		addSkip(result, CodeTranscriptUndeclared, horizon.missing)
	}
	for _, item := range transcripts {
		if item.Malformed == "" {
			continue
		}
		result.Findings = append(result.Findings, Finding{Code: CodeTranscriptMalformed, Severity: SeverityError,
			Summary: fmt.Sprintf("%s declares Surface Transcript %d %s", artifact, item.Number, item.Malformed),
			Where:   []Location{{Path: artifact, Line: item.Line}},
			Fix:     "Write the block as a command line starting with \"$ \", then \"stdout:\", \"stderr:\" and \"exit: <code>\", in that order."})
	}
	return coverageUnitsDeclaredIn(declaration.units, artifact)
}

func detectTranscriptGate(result *Result, repoRoot string, graph *spec.Graph, units []coverageUnit) {
	if graph.QADeclined || graph.QATaskID == "" {
		return
	}
	for _, task := range graph.Tasks {
		if task.ID != graph.QATaskID || task.Status == spec.StatusCompleted {
			continue
		}
		references := newReferenceSet()
		for _, requirement := range task.Requirements {
			addDeclaredReferences(references, requirement.Text)
		}
		path := artifactDisplayPath(repoRoot, filepath.Join(filepath.Dir(graph.Spec.Dir), task.File))
		line := 1
		if len(task.Requirements) > 0 {
			line = task.Requirements[0].Line
		}
		for _, unit := range units {
			if unit.Kind != coverageTranscript || references[coverageTranscript][unit.Number] {
				continue
			}
			name := coverageUnitName(unit)
			result.Findings = append(result.Findings, Finding{Code: CodeTranscriptUngated, Severity: SeverityError,
				Summary: unit.DeclaredIn + " declares " + name + ", but QA Task " + path + " names it in no Requirement",
				Where:   []Location{{Path: unit.DeclaredIn, Line: unit.Line}, {Path: path, Line: line}},
				Fix:     "Name " + name + " in a Requirement of " + path + " so the gate reproduces it."})
		}
	}
}
