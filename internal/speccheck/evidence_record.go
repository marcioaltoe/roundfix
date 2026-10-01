package speccheck

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"gopkg.in/yaml.v3"

	"roundfix/internal/spec"
)

// EvidenceRecord identifies the rows snapshotted at the audited head, in Results order.
type EvidenceRecord struct {
	Head string
	Rows []string
}

// EvidenceReportFileError distinguishes report IO failures from Git read errors.
type EvidenceReportFileError struct{ Err error }

func (err *EvidenceReportFileError) Error() string { return err.Err.Error() }
func (err *EvidenceReportFileError) Unwrap() error { return err.Err }

// AlwaysObserved names rows whose observations must be repeated on every pass.
func AlwaysObserved(provenance string, inputs []EvidenceInput) (string, bool) {
	for _, item := range strings.FieldsFunc(provenance, func(r rune) bool { return r == ';' || r == ',' }) {
		switch strings.TrimSpace(item) {
		case "repository Verification":
			return CarryReasonRepositoryVerification, true
		case spec.QAPullRequestRowSource:
			return CarryReasonPullRequestRow, true
		}
	}
	for _, input := range inputs {
		if input.Kind == EvidenceCommitRange {
			return CarryReasonCommitRangeInput, true
		}
	}
	return "", false
}

var evidenceSnapshotKey = regexp.MustCompile(`^(?:evidence_snapshots|"evidence_snapshots"|'evidence_snapshots')\s*:`)

// RecordEvidenceSnapshots replaces Agent-written snapshots with Git blob evidence.
// A Git read error still writes the report with the untrusted key removed.
func RecordEvidenceSnapshots(ctx context.Context, repoRoot, reportPath, head string) (EvidenceRecord, error) {
	record := EvidenceRecord{Head: head}
	if !filepath.IsAbs(reportPath) {
		reportPath = filepath.Join(repoRoot, reportPath)
	}
	content, err := os.ReadFile(reportPath)
	if err != nil {
		return record, &EvidenceReportFileError{Err: fmt.Errorf("read evidence report: %w", err)}
	}
	stripped, insertion, newline := stripEvidenceSnapshots(content)
	relative, err := filepath.Rel(repoRoot, reportPath)
	if err != nil {
		return record, fmt.Errorf("resolve evidence report path: %w", err)
	}
	report := parseMechanicalReport(filepath.ToSlash(relative), stripped)
	mapping := &yaml.Node{Kind: yaml.MappingNode}
	var readErr error
	if report.parseError == nil && insertion >= 0 {
		for _, row := range report.rows {
			if row.status != "pass" || len(row.inputs) == 0 {
				continue
			}
			if _, observed := AlwaysObserved(row.provenance, row.inputs); observed {
				continue
			}
			repositoryOnly := true
			for _, input := range row.inputs {
				if input.Kind != EvidenceRepositoryPath {
					repositoryOnly = false
				}
			}
			if !repositoryOnly {
				continue
			}
			snapshots, resolved, err := buildEvidenceSnapshots(ctx, repoRoot, head, row.inputs)
			if err != nil {
				readErr = err
				break
			}
			if !resolved {
				continue
			}
			paths, resolved := mechanicalEvidenceRepositoryPaths(repoRoot, report.path, row.evidence)
			if !resolved {
				continue
			}
			covered := true
			for _, path := range paths {
				found := false
				for _, snapshot := range snapshots {
					if evidenceSnapshotContains(snapshot, path) {
						found = true
					}
				}
				if !found {
					covered = false
				}
			}
			if !covered {
				continue
			}
			// Keep declaration and Results ordering while letting YAML quote scalars.
			inputs := &yaml.Node{Kind: yaml.SequenceNode}
			for _, snapshot := range snapshots {
				files := &yaml.Node{Kind: yaml.SequenceNode}
				for _, file := range snapshot.Files {
					files.Content = append(files.Content, evidenceMapping("path", evidenceScalar(file.Path), evidenceScalar("sha256"), evidenceScalar(file.SHA256)))
				}
				inputs.Content = append(inputs.Content, evidenceMapping("ref", evidenceScalar(snapshot.Ref), evidenceScalar("files"), files))
			}
			mapping.Content = append(mapping.Content, &yaml.Node{Kind: yaml.ScalarNode, Tag: "!!str", Style: yaml.DoubleQuotedStyle, Value: row.id}, evidenceMapping("head", evidenceScalar(head), evidenceScalar("inputs"), inputs))
			record.Rows = append(record.Rows, row.id)
		}
	}
	if readErr != nil {
		record.Rows = nil
	}
	if len(record.Rows) > 0 {
		var block bytes.Buffer
		encoder := yaml.NewEncoder(&block)
		encoder.SetIndent(2)
		if err := encoder.Encode(evidenceMapping("evidence_snapshots", mapping)); err != nil {
			return record, fmt.Errorf("encode evidence snapshots: %w", err)
		}
		rendered := bytes.ReplaceAll(block.Bytes(), []byte("\n"), newline)
		stripped = append(append(append([]byte(nil), stripped[:insertion]...), rendered...), stripped[insertion:]...)
	}
	if err := os.WriteFile(reportPath, stripped, 0o644); err != nil {
		return record, &EvidenceReportFileError{Err: fmt.Errorf("write evidence report: %w", err)}
	}
	if ctx.Err() != nil {
		return record, ctx.Err()
	}
	return record, readErr
}

func evidenceScalar(value string) *yaml.Node {
	return &yaml.Node{Kind: yaml.ScalarNode, Tag: "!!str", Value: value}
}
func evidenceMapping(key string, value *yaml.Node, rest ...*yaml.Node) *yaml.Node {
	return &yaml.Node{Kind: yaml.MappingNode, Content: append([]*yaml.Node{evidenceScalar(key), value}, rest...)}
}

// stripEvidenceSnapshots edits whole lines only inside frontmatter.
func stripEvidenceSnapshots(content []byte) ([]byte, int, []byte) {
	lines := bytes.SplitAfter(content, []byte("\n"))
	newline := []byte("\n")
	if len(lines) == 0 || strings.TrimRight(string(lines[0]), "\r\n") != "---" {
		return content, -1, newline
	}
	if bytes.HasSuffix(lines[0], []byte("\r\n")) {
		newline = []byte("\r\n")
	}
	var out bytes.Buffer
	out.Write(lines[0])
	skipping := false
	for i := 1; i < len(lines); i++ {
		line := lines[i]
		text := strings.TrimRight(string(line), "\r\n")
		if text == "---" {
			insertion := out.Len()
			for _, rest := range lines[i:] {
				out.Write(rest)
			}
			return out.Bytes(), insertion, newline
		}
		if evidenceSnapshotKey.MatchString(text) {
			skipping = true
			continue
		}
		if skipping && (strings.HasPrefix(text, " ") || strings.HasPrefix(text, "\t")) {
			continue
		}
		if text != "" {
			skipping = false
		}
		out.Write(line)
	}
	return out.Bytes(), -1, newline
}
