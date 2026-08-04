package speccheck

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

const (
	// CodeConstraintMissing identifies an absent required Project Constraint row.
	CodeConstraintMissing = "SC-CONSTRAINT-MISSING"
	// CodeConstraintUnreasoned identifies applicability stated without a reason.
	CodeConstraintUnreasoned = "SC-CONSTRAINT-UNREASONED"
	// CodeConstraintSource identifies a missing cited source path.
	CodeConstraintSource = "SC-CONSTRAINT-SOURCE"
	// CodeToolingUnauthorized identifies a cited authorization that omits the Spec.
	CodeToolingUnauthorized = "SC-TOOLING-UNAUTHORIZED"
	// CodeToolingUnbounded identifies applicable tooling authority without bounded files.
	CodeToolingUnbounded = "SC-TOOLING-UNBOUNDED"
)

const (
	constraintIdentifier = "Identifier strategy"
	constraintAuthHTTP   = "Authentication and HTTP"
	constraintActiveADR  = "Active ADR obligations"
	constraintTooling    = "Tooling authority"
)

var (
	requiredConstraints = []string{
		constraintIdentifier,
		constraintAuthHTTP,
		constraintActiveADR,
		constraintTooling,
	}
	detectorCodes = []string{
		CodeConstraintMissing,
		CodeConstraintUnreasoned,
		CodeConstraintSource,
		CodeToolingUnauthorized,
		CodeToolingUnbounded,
	}
	sourcePathPattern = regexp.MustCompile("(?is)\\bSource:\\s*`([^`]+)`")
	backtickPattern   = regexp.MustCompile("`([^`]+)`")
)

type applicability string

const (
	applicable    applicability = "applicable"
	notApplicable applicability = "not applicable"
)

type constraintRow struct {
	Label             string
	Applicability     applicability
	Reason            string
	SourcePath        string
	AuthorizationPath string
	BoundedFiles      bool
	Line              int
}

type constraintArtifact struct {
	displayPath string
	sectionLine int
	rows        map[string]constraintRow
}

// Check reads one Spec folder and the repository facts cited by its Project
// Constraints. Missing input artifacts skip their detectors; other read
// failures are returned with operation context.
func Check(specsRoot, repoRoot, slug string) (Result, error) {
	result := Result{
		Slug:     slug,
		Findings: []Finding{},
		Skipped:  []SkippedDetector{},
	}

	if strings.TrimSpace(slug) == "" || filepath.Base(slug) != slug || slug == "." {
		return result, fmt.Errorf("invalid Spec slug %q", slug)
	}
	specDir := filepath.Join(filepath.Clean(specsRoot), slug)
	info, err := os.Stat(specDir)
	if err != nil {
		return result, fmt.Errorf("read Spec %q: %w", slug, err)
	}
	if !info.IsDir() {
		return result, fmt.Errorf("read Spec %q: %s is not a directory", slug, specDir)
	}

	prd, present, err := readConstraintArtifact(repoRoot, filepath.Join(specDir, "_prd.md"))
	if err != nil {
		return result, err
	}
	if !present {
		for _, code := range detectorCodes {
			addSkip(&result, code, artifactDisplayPath(repoRoot, filepath.Join(specDir, "_prd.md")))
		}
		for _, code := range citationCoverageDetectorCodes {
			addSkip(&result, code, artifactDisplayPath(repoRoot, filepath.Join(specDir, "_prd.md")))
		}
		return result, nil
	}

	artifacts := []constraintArtifact{prd}
	techSpecPath := filepath.Join(specDir, "_techspec.md")
	techSpec, present, err := readConstraintArtifact(repoRoot, techSpecPath)
	if err != nil {
		return result, err
	}
	if present {
		artifacts = append(artifacts, techSpec)
	} else {
		for _, code := range detectorCodes {
			addSkip(&result, code, artifactDisplayPath(repoRoot, techSpecPath))
		}
	}
	if err := detectVocabularyContract(&result, repoRoot, techSpecPath, present); err != nil {
		return result, err
	}

	for artifactIndex := range artifacts {
		detectConstraintRows(&result, repoRoot, slug, artifacts, artifactIndex)
	}
	if err := detectCitationCoverageAndReferences(&result, specsRoot, repoRoot, slug, specDir, present); err != nil {
		return result, err
	}
	return result, nil
}

func readConstraintArtifact(repoRoot, path string) (constraintArtifact, bool, error) {
	content, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return constraintArtifact{}, false, nil
	}
	if err != nil {
		return constraintArtifact{}, false, fmt.Errorf("read Spec artifact %q: %w", path, err)
	}
	rows, sectionLine := parseProjectConstraints(content)
	return constraintArtifact{
		displayPath: artifactDisplayPath(repoRoot, path),
		sectionLine: sectionLine,
		rows:        rows,
	}, true, nil
}

func parseProjectConstraints(content []byte) (map[string]constraintRow, int) {
	lines := strings.Split(string(content), "\n")
	rows := make(map[string]constraintRow)
	sectionLine := 1
	sectionStart := -1
	for index, line := range lines {
		if strings.TrimSpace(line) == "## Project Constraints" {
			sectionLine = index + 1
			sectionStart = index + 1
			break
		}
	}
	if sectionStart == -1 {
		return rows, sectionLine
	}

	var rowLines []string
	rowLine := 0
	flush := func() {
		if len(rowLines) == 0 {
			return
		}
		raw := strings.Join(rowLines, " ")
		row, ok := parseConstraintRow(raw, rowLine)
		if ok {
			rows[strings.ToLower(row.Label)] = row
		}
		rowLines = nil
		rowLine = 0
	}

	for index := sectionStart; index < len(lines); index++ {
		trimmed := strings.TrimSpace(lines[index])
		if strings.HasPrefix(trimmed, "## ") {
			break
		}
		if strings.HasPrefix(trimmed, "- ") {
			flush()
			rowLine = index + 1
			rowLines = append(rowLines, strings.TrimSpace(strings.TrimPrefix(trimmed, "- ")))
			continue
		}
		if len(rowLines) > 0 && trimmed != "" {
			rowLines = append(rowLines, trimmed)
		}
	}
	flush()
	return rows, sectionLine
}

func parseConstraintRow(raw string, line int) (constraintRow, bool) {
	label, declaration, ok := strings.Cut(raw, ":")
	if !ok {
		return constraintRow{}, false
	}
	label = strings.TrimSpace(label)
	declaration = strings.TrimSpace(declaration)
	row := constraintRow{Label: label, Line: line}
	lowerDeclaration := strings.ToLower(declaration)
	switch {
	case strings.HasPrefix(lowerDeclaration, string(notApplicable)):
		row.Applicability = notApplicable
		declaration = declaration[len(notApplicable):]
	case strings.HasPrefix(lowerDeclaration, string(applicable)):
		row.Applicability = applicable
		declaration = declaration[len(applicable):]
	}

	if match := sourcePathPattern.FindStringSubmatch(raw); len(match) == 2 {
		row.SourcePath = filepath.ToSlash(filepath.Clean(strings.TrimSpace(match[1])))
	}
	reason := declaration
	if sourceIndex := strings.Index(strings.ToLower(reason), "source:"); sourceIndex >= 0 {
		reason = reason[:sourceIndex]
	}
	row.Reason = strings.Trim(strings.TrimSpace(reason), "—–-:.; ")
	row.AuthorizationPath = authorizationRecordPath(raw, row.SourcePath)
	row.BoundedFiles = recordsBoundedFiles(raw)
	return row, true
}

func authorizationRecordPath(raw, sourcePath string) string {
	for _, match := range backtickPattern.FindAllStringSubmatch(raw, -1) {
		path := filepath.ToSlash(filepath.Clean(strings.TrimSpace(match[1])))
		if path == sourcePath {
			continue
		}
		if strings.Contains(strings.ToLower(path), "authoriz") && strings.HasSuffix(strings.ToLower(path), ".md") {
			return path
		}
	}
	return ""
}

func recordsBoundedFiles(raw string) bool {
	lower := strings.ToLower(raw)
	if index := strings.Index(lower, "bounded files:"); index >= 0 {
		value := raw[index+len("bounded files:"):]
		if sourceIndex := strings.Index(strings.ToLower(value), "source:"); sourceIndex >= 0 {
			value = value[:sourceIndex]
		}
		value = strings.Trim(strings.TrimSpace(value), ".; ")
		return value != "" && !strings.EqualFold(value, "none")
	}
	return strings.Contains(lower, "changes to exactly") ||
		strings.Contains(lower, "authorizes exactly") ||
		strings.Contains(lower, "bounded repository-relative")
}

func detectConstraintRows(result *Result, repoRoot, slug string, artifacts []constraintArtifact, artifactIndex int) {
	artifact := artifacts[artifactIndex]
	for _, label := range requiredConstraints {
		row, present := artifact.rows[strings.ToLower(label)]
		if !present {
			secondary := counterpartLocation(artifacts, artifactIndex, label)
			result.Findings = append(result.Findings, Finding{
				Code:     CodeConstraintMissing,
				Severity: SeverityError,
				Summary:  artifact.displayPath + " omits " + label + " required by " + secondary.Path,
				Where: []Location{
					{Path: artifact.displayPath, Line: artifact.sectionLine},
					secondary,
				},
				Fix: "Add the " + label + " row to " + artifact.displayPath + " with applicability, a reason, and an operative Source path.",
			})
			continue
		}

		rowLocation := Location{Path: artifact.displayPath, Line: row.Line}
		secondary := sourceLocation(row)
		if row.Applicability != "" && row.Reason == "" {
			result.Findings = append(result.Findings, Finding{
				Code:     CodeConstraintUnreasoned,
				Severity: SeverityError,
				Summary:  artifact.displayPath + " declares " + label + " " + string(row.Applicability) + " without the reason required by " + secondary.Path,
				Where:    []Location{rowLocation, secondary},
				Fix:      "Write the reason after " + string(row.Applicability) + " in the " + label + " row in " + artifact.displayPath + ".",
			})
		}

		if row.SourcePath == "" || !regularRepositoryFile(repoRoot, row.SourcePath) {
			missingPath := row.SourcePath
			if missingPath == "" {
				missingPath = "docs/agents/"
			}
			result.Findings = append(result.Findings, Finding{
				Code:     CodeConstraintSource,
				Severity: SeverityError,
				Summary:  artifact.displayPath + " cites missing " + missingPath + " for " + label,
				Where: []Location{
					rowLocation,
					{Path: missingPath, Line: 1},
				},
				Fix: "Point the " + label + " row in " + artifact.displayPath + " to an existing operative docs/agents/ source file.",
			})
		}

		if label == constraintTooling && row.Applicability == applicable {
			detectToolingRow(result, repoRoot, slug, artifact, row)
		}
	}
}

func detectToolingRow(result *Result, repoRoot, slug string, artifact constraintArtifact, row constraintRow) {
	rowLocation := Location{Path: artifact.displayPath, Line: row.Line}
	recordLocation := sourceLocation(row)
	if row.AuthorizationPath != "" {
		recordLocation = Location{Path: row.AuthorizationPath, Line: 1}
		recordPath, ok := resolveRepositoryPath(repoRoot, row.AuthorizationPath)
		if !ok {
			addSkip(result, CodeToolingUnauthorized, row.AuthorizationPath)
		} else {
			content, err := os.ReadFile(recordPath)
			switch {
			case errors.Is(err, os.ErrNotExist):
				addSkip(result, CodeToolingUnauthorized, row.AuthorizationPath)
			case err != nil:
				result.Findings = append(result.Findings, Finding{
					Code:     CodeToolingUnauthorized,
					Severity: SeverityError,
					Summary:  artifact.displayPath + " cites unreadable " + row.AuthorizationPath + " for Spec " + slug,
					Where:    []Location{rowLocation, recordLocation},
					Fix:      "Make " + row.AuthorizationPath + " readable and name Spec " + slug + " in its authorization scope.",
				})
			case !authorizationNamesSpec(content, slug):
				result.Findings = append(result.Findings, Finding{
					Code:     CodeToolingUnauthorized,
					Severity: SeverityError,
					Summary:  artifact.displayPath + " cites " + row.AuthorizationPath + ", but that record does not name Spec " + slug,
					Where:    []Location{rowLocation, recordLocation},
					Fix:      "Add Spec " + slug + " to " + row.AuthorizationPath + " or cite the authorization record that already names it.",
				})
			}
		}
	}

	if !row.BoundedFiles {
		result.Findings = append(result.Findings, Finding{
			Code:     CodeToolingUnbounded,
			Severity: SeverityError,
			Summary:  artifact.displayPath + " makes Tooling authority applicable without bounded files alongside " + recordLocation.Path,
			Where:    []Location{rowLocation, recordLocation},
			Fix:      "Add an exact bounded files list to the Tooling authority row in " + artifact.displayPath + ".",
		})
	}
}

func authorizationNamesSpec(content []byte, slug string) bool {
	text := strings.ToLower(string(content))
	if strings.Contains(text, strings.ToLower(slug)) {
		return true
	}
	number, _, hasSuffix := strings.Cut(slug, "-")
	if !hasSuffix || number == "" {
		return false
	}
	for _, char := range number {
		if char < '0' || char > '9' {
			return false
		}
	}
	pattern := regexp.MustCompile(`(^|[^0-9])` + regexp.QuoteMeta(number) + `([^0-9]|$)`)
	return pattern.MatchString(text)
}

func counterpartLocation(artifacts []constraintArtifact, artifactIndex int, label string) Location {
	for index, artifact := range artifacts {
		if index == artifactIndex {
			continue
		}
		if row, ok := artifact.rows[strings.ToLower(label)]; ok {
			return Location{Path: artifact.displayPath, Line: row.Line}
		}
	}
	return Location{Path: "docs/agents/spec-routing.md", Line: 1}
}

func sourceLocation(row constraintRow) Location {
	if row.SourcePath != "" {
		return Location{Path: row.SourcePath, Line: 1}
	}
	return Location{Path: "docs/agents/spec-routing.md", Line: 1}
}

func addSkip(result *Result, code, missing string) {
	for _, skipped := range result.Skipped {
		if skipped.Code == code && skipped.Missing == missing {
			return
		}
	}
	result.Skipped = append(result.Skipped, SkippedDetector{Code: code, Missing: missing})
}

func regularRepositoryFile(repoRoot, relative string) bool {
	path, ok := resolveRepositoryPath(repoRoot, relative)
	if !ok {
		return false
	}
	info, err := os.Stat(path)
	return err == nil && info.Mode().IsRegular()
}

func resolveRepositoryPath(repoRoot, relative string) (string, bool) {
	if filepath.IsAbs(relative) {
		return "", false
	}
	clean := filepath.Clean(filepath.FromSlash(relative))
	if clean == "." || clean == ".." || strings.HasPrefix(clean, ".."+string(filepath.Separator)) {
		return "", false
	}
	return filepath.Join(filepath.Clean(repoRoot), clean), true
}

func artifactDisplayPath(repoRoot, path string) string {
	relative, err := filepath.Rel(filepath.Clean(repoRoot), filepath.Clean(path))
	if err != nil {
		return filepath.ToSlash(filepath.Clean(path))
	}
	return filepath.ToSlash(relative)
}
