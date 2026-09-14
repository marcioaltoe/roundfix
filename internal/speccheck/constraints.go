package speccheck

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"roundfix/internal/spec"
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
	// CodeToolingUnbounded identifies a declared tooling mutation without bounded files.
	CodeToolingUnbounded = "SC-TOOLING-UNBOUNDED"
	// CodeToolingUntyped identifies an authorization record that states its grant
	// only in prose, so no checker can enumerate what it permits.
	CodeToolingUntyped = "SC-TOOLING-UNTYPED"
	// CodeToolingUnapproved identifies a claimed authorization whose cited
	// record is not an operative grant.
	CodeToolingUnapproved = "SC-TOOLING-UNAPPROVED"
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
		CodeToolingUntyped,
		CodeToolingUnapproved,
	}
	sourcePathPattern   = regexp.MustCompile("(?is)\\bSource:\\s*`([^`]+)`")
	backtickPattern     = regexp.MustCompile("`([^`]+)`")
	markdownLinkPattern = regexp.MustCompile(`\[[^\]]*\]\(([^\s)]+)\)`)
)

type applicability string

const (
	applicable    applicability = "applicable"
	notApplicable applicability = "not applicable"
)

type constraintRow struct {
	Label         string
	Applicability applicability
	Reason        string
	SourcePath    string
	Raw           string
	Authorization authorizationReferenceSelection
	BoundedFiles  bool
	Line          int
}

type authorizationReferenceKind uint8

const (
	authorizationReferenceUnspecified authorizationReferenceKind = iota
	authorizationReferenceOperative
	authorizationReferenceProposed
)

type authorizationReference struct {
	Path         string
	ResolvedPath string
	ReadRoot     string
	ReadPath     string
	Kind         authorizationReferenceKind
}

type authorizationReferenceSelection struct {
	Candidates []authorizationReference
	Reference  authorizationReference
	Selected   bool
	Ambiguous  bool
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
	if err := detectLoopOrderConsistency(&result, repoRoot); err != nil {
		return result, err
	}
	if err := detectFindingsConsistency(&result, repoRoot); err != nil {
		return result, err
	}
	if err := detectBacklogPromotion(&result, repoRoot); err != nil {
		return result, fmt.Errorf("detect backlog promotion: %w", err)
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

	executableTaskGraph := regularRepositoryFile(repoRoot, artifactDisplayPath(repoRoot, filepath.Join(specDir, "_tasks.md")))
	for artifactIndex := range artifacts {
		detectConstraintRows(&result, repoRoot, slug, artifacts, artifactIndex, executableTaskGraph && artifactIndex == 0)
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
	displayPath := artifactDisplayPath(repoRoot, path)
	rows, sectionLine := parseProjectConstraints(content)
	if row, ok := rows[strings.ToLower(constraintTooling)]; ok {
		candidates := authorizationReferences(row.Raw, row.SourcePath, repoRoot, path)
		reference, selected, ambiguous := selectAuthorizationReference(row, candidates)
		row.Authorization = authorizationReferenceSelection{
			Candidates: candidates,
			Reference:  reference,
			Selected:   selected,
			Ambiguous:  ambiguous,
		}
		rows[strings.ToLower(constraintTooling)] = row
	}
	return constraintArtifact{
		displayPath: displayPath,
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
	row := constraintRow{Label: label, Raw: raw, Line: line}
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
	row.BoundedFiles = recordsBoundedFiles(raw)
	return row, true
}

func authorizationReferences(raw, sourcePath, repoRoot, artifactPath string) []authorizationReference {
	type rawReference struct {
		target   string
		start    int
		markdown bool
	}

	var rawReferences []rawReference
	for _, match := range backtickPattern.FindAllStringSubmatchIndex(raw, -1) {
		rawReferences = append(rawReferences, rawReference{
			target: raw[match[2]:match[3]],
			start:  match[0],
		})
	}
	for _, match := range markdownLinkPattern.FindAllStringSubmatchIndex(raw, -1) {
		rawReferences = append(rawReferences, rawReference{
			target:   raw[match[2]:match[3]],
			start:    match[0],
			markdown: true,
		})
	}
	sort.Slice(rawReferences, func(first, second int) bool {
		return rawReferences[first].start < rawReferences[second].start
	})

	seen := make(map[string]bool)
	var references []authorizationReference
	for _, rawReference := range rawReferences {
		reference, ok := authorizationReferencePath(rawReference.target, repoRoot, artifactPath, rawReference.markdown)
		if !ok || reference.Path == sourcePath || seen[reference.Path] {
			continue
		}
		lowerPath := strings.ToLower(reference.Path)
		if !strings.Contains(lowerPath, "authoriz") || !strings.HasSuffix(lowerPath, ".md") {
			continue
		}
		seen[reference.Path] = true
		reference.Kind = authorizationReferenceKindAt(raw, rawReference.start)
		references = append(references, reference)
	}
	return references
}

func authorizationReferencePath(target, repoRoot, artifactPath string, markdown bool) (authorizationReference, bool) {
	target = strings.Trim(strings.TrimSpace(target), "<>")
	if target == "" || strings.Contains(target, `\`) || strings.ContainsAny(target, "?#") || path.IsAbs(target) {
		return authorizationReference{}, false
	}
	clean := path.Clean(target)
	if clean == "." {
		return authorizationReference{}, false
	}
	if !markdown {
		if clean == ".." || strings.HasPrefix(clean, "../") {
			return authorizationReference{}, false
		}
		root, err := filepath.Abs(filepath.Clean(repoRoot))
		if err != nil {
			return authorizationReference{}, false
		}
		return authorizationReference{
			Path:         clean,
			ResolvedPath: filepath.Join(root, filepath.FromSlash(clean)),
			ReadRoot:     root,
			ReadPath:     clean,
		}, true
	}

	artifactPath = filepath.Clean(artifactPath)
	specsRoot := filepath.Dir(filepath.Dir(artifactPath))
	resolvedPath := filepath.Clean(filepath.Join(filepath.Dir(artifactPath), filepath.FromSlash(clean)))
	if !constraintPathWithinRoot(resolvedPath, specsRoot) {
		return authorizationReference{}, false
	}

	root := specsRoot
	if _, ok := constraintRelativePath(repoRoot, resolvedPath); ok {
		root = repoRoot
	}
	displayPath := artifactDisplayPath(repoRoot, resolvedPath)
	root, err := filepath.Abs(filepath.Clean(root))
	if err != nil {
		return authorizationReference{}, false
	}
	resolvedPath, err = filepath.Abs(resolvedPath)
	if err != nil {
		return authorizationReference{}, false
	}
	relativePath, ok := constraintRelativePath(root, resolvedPath)
	if !ok {
		return authorizationReference{}, false
	}
	return authorizationReference{
		Path:         displayPath,
		ResolvedPath: resolvedPath,
		ReadRoot:     root,
		ReadPath:     relativePath,
	}, true
}

func constraintPathWithinRoot(candidate, root string) bool {
	candidate, err := filepath.Abs(candidate)
	if err != nil {
		return false
	}
	root, err = filepath.Abs(root)
	if err != nil {
		return false
	}
	_, ok := constraintRelativePath(canonicalConstraintPath(root), canonicalConstraintPath(candidate))
	return ok
}

func constraintRelativePath(root, candidate string) (string, bool) {
	relative, err := filepath.Rel(filepath.Clean(root), filepath.Clean(candidate))
	if err != nil || relative == "." || relative == ".." || filepath.IsAbs(relative) || strings.HasPrefix(relative, ".."+string(filepath.Separator)) {
		return "", false
	}
	return filepath.ToSlash(relative), true
}

func canonicalConstraintPath(value string) string {
	original := filepath.Clean(value)
	candidate := original
	var tail []string
	for {
		resolved, err := filepath.EvalSymlinks(candidate)
		if err == nil {
			for index := len(tail) - 1; index >= 0; index-- {
				resolved = filepath.Join(resolved, tail[index])
			}
			return filepath.Clean(resolved)
		}
		parent := filepath.Dir(candidate)
		if parent == candidate {
			return original
		}
		tail = append(tail, filepath.Base(candidate))
		candidate = parent
	}
}

func authorizationReferenceKindAt(raw string, start int) authorizationReferenceKind {
	clauseStart := 0
	for _, separator := range []string{";", ". ", "\n"} {
		if index := strings.LastIndex(raw[:start], separator); index >= clauseStart {
			clauseStart = index + len(separator)
		}
	}
	context := strings.ToLower(raw[clauseStart:start])
	if strings.Contains(context, "proposed") {
		return authorizationReferenceProposed
	}
	if claimsOperativeAuthorization(context) {
		return authorizationReferenceOperative
	}
	return authorizationReferenceUnspecified
}

func recordsBoundedFiles(raw string) bool {
	lower := strings.ToLower(raw)
	boundedPrefix := "bounded files:"
	index := strings.Index(lower, boundedPrefix)
	if proposedIndex := strings.Index(lower, "bounded proposed files:"); proposedIndex >= 0 && (index < 0 || proposedIndex < index) {
		boundedPrefix = "bounded proposed files:"
		index = proposedIndex
	}
	if index >= 0 {
		value := raw[index+len(boundedPrefix):]
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

func detectConstraintRows(result *Result, repoRoot, slug string, artifacts []constraintArtifact, artifactIndex int, requireImplement bool) {
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
			detectToolingRow(result, repoRoot, slug, artifact, row, requireImplement)
		}
	}
}

func detectToolingRow(result *Result, repoRoot, slug string, artifact constraintArtifact, row constraintRow, requireImplement bool) {
	rowLocation := Location{Path: artifact.displayPath, Line: row.Line}
	recordLocation := sourceLocation(row)
	selection := row.Authorization
	if selection.Ambiguous {
		locations := []Location{rowLocation}
		for _, reference := range selection.Candidates {
			locations = append(locations, Location{Path: reference.Path, Line: 1})
		}
		result.Findings = append(result.Findings, Finding{
			Code:     CodeToolingUnapproved,
			Severity: SeverityError,
			Summary:  artifact.displayPath + " claims express maintainer authorization, but the claim does not identify exactly one authorization record",
			Where:    locations,
			Fix:      "Cite exactly one operative authorization record for the claim in " + artifact.displayPath + ", and describe every proposal separately.",
		})
	}
	if selection.Selected {
		reference := selection.Reference
		recordLocation = Location{Path: reference.Path, Line: 1}
		content, err := os.ReadFile(reference.ResolvedPath)
		switch {
		case errors.Is(err, os.ErrNotExist):
			addSkip(result, CodeToolingUnauthorized, reference.Path)
		case err != nil:
			result.Findings = append(result.Findings, Finding{
				Code:     CodeToolingUnauthorized,
				Severity: SeverityError,
				Summary:  artifact.displayPath + " cites unreadable " + reference.Path + " for Spec " + slug,
				Where:    []Location{rowLocation, recordLocation},
				Fix:      "Make " + reference.Path + " readable and name Spec " + slug + " in its authorization scope.",
			})
		default:
			detectAuthorizationResolution(result, slug, artifact, row, reference, content, requireImplement)
		}
	}

	if declaresProtectedToolingMutation(row, len(selection.Candidates) != 0) && !row.BoundedFiles {
		result.Findings = append(result.Findings, Finding{
			Code:     CodeToolingUnbounded,
			Severity: SeverityError,
			Summary:  artifact.displayPath + " declares a protected tooling mutation without bounded files alongside " + recordLocation.Path,
			Where:    []Location{rowLocation, recordLocation},
			Fix:      "Add an exact bounded files list to the Tooling authority row in " + artifact.displayPath + ".",
		})
	}
}

func selectAuthorizationReference(row constraintRow, references []authorizationReference) (authorizationReference, bool, bool) {
	if !claimsOperativeAuthorization(row.Reason) {
		if len(references) == 1 {
			return references[0], true, false
		}
		return authorizationReference{}, false, false
	}

	var operative []authorizationReference
	for _, reference := range references {
		if reference.Kind == authorizationReferenceOperative {
			operative = append(operative, reference)
		}
	}
	if len(operative) == 1 {
		return operative[0], true, false
	}
	if len(operative) == 0 && len(references) == 1 && references[0].Kind != authorizationReferenceProposed {
		return references[0], true, false
	}
	return authorizationReference{}, false, true
}

func claimsOperativeAuthorization(value string) bool {
	lower := strings.ToLower(strings.Join(strings.Fields(value), " "))
	return strings.Contains(lower, "express maintainer authorization") ||
		strings.Contains(lower, "protected tooling mutation authorized") ||
		strings.Contains(lower, "authorization is recorded") ||
		strings.Contains(lower, "authorization recorded") ||
		strings.Contains(lower, "authorized at")
}

func declaresProtectedToolingMutation(row constraintRow, citesAuthorization bool) bool {
	if citesAuthorization {
		return true
	}
	reason := strings.ToLower(strings.Join(strings.Fields(row.Reason), " "))
	if strings.Contains(reason, "no protected tooling mutation") {
		return false
	}
	return strings.Contains(reason, "protected tooling mutation") ||
		strings.Contains(reason, "express maintainer authorization")
}

func detectAuthorizationResolution(
	result *Result,
	slug string,
	artifact constraintArtifact,
	row constraintRow,
	reference authorizationReference,
	content []byte,
	requireImplement bool,
) {
	role := authorizationRole(reference.Path)
	resolution := spec.ReadAuthorization(context.Background(), spec.AuthorizationReadRequest{
		RepoRoot:   reference.ReadRoot,
		RecordPath: reference.ReadPath,
		Role:       role,
		AskingSpec: slug,
	})
	if resolution.Outcome == spec.AuthorizationGranted {
		if requireImplement && !resolution.Permits(spec.AuthorizationOperationImplement) {
			rowLocation := Location{Path: artifact.displayPath, Line: row.Line}
			recordLocation := Location{Path: reference.Path, Line: 1}
			result.Findings = append(result.Findings, Finding{
				Code:     CodeToolingUnapproved,
				Severity: SeverityError,
				Summary:  fmt.Sprintf("%s carries an executable Task Graph, but authorization record %s does not permit operation %q", artifact.displayPath, reference.Path, spec.AuthorizationOperationImplement),
				Where:    []Location{rowLocation, recordLocation},
				Fix:      fmt.Sprintf("Add operation %q to %s before treating the Spec as dispatchable.", spec.AuthorizationOperationImplement, reference.Path),
			})
		}
		return
	}

	rowLocation := Location{Path: artifact.displayPath, Line: row.Line}
	recordLocation := Location{Path: reference.Path, Line: 1}
	if role == spec.AuthorizationRoleLegacy &&
		resolution.Outcome == spec.AuthorizationRefused &&
		resolution.Reason.Code == spec.AuthorizationReasonGranted &&
		resolution.Record.GrantedAt.IsZero() {
		if authorizationNamesSpec(content, slug) {
			return
		}
		appendToolingUnauthorized(result, artifact.displayPath, slug, reference.Path, rowLocation, recordLocation)
		return
	}
	if resolution.Outcome == spec.AuthorizationUnresolved {
		result.Findings = append(result.Findings, Finding{
			Code:     CodeToolingUnauthorized,
			Severity: SeverityError,
			Summary:  artifact.displayPath + " cites unresolved " + reference.Path + " for Spec " + slug + ": " + resolution.Reason.Detail,
			Where:    []Location{rowLocation, recordLocation},
			Fix:      "Make " + reference.Path + " readable from the repository and keep its authorization scope repository-relative.",
		})
		return
	}

	switch resolution.Reason.Code {
	case spec.AuthorizationReasonConsuming:
		if len(resolution.Record.Consuming) != 0 {
			appendToolingUnauthorized(result, artifact.displayPath, slug, reference.Path, rowLocation, recordLocation)
			return
		}
		appendToolingUntyped(result, reference.Path, resolution.Reason)
	case spec.AuthorizationReasonMalformedRecord, spec.AuthorizationReasonAction:
		appendToolingUntyped(result, reference.Path, resolution.Reason)
	case spec.AuthorizationReasonStatus:
		if resolution.Record.Status == "" {
			appendToolingUntyped(result, reference.Path, resolution.Reason)
			return
		}
		appendToolingUnapproved(result, artifact.displayPath, reference.Path, rowLocation, recordLocation, resolution.Reason, claimsOperativeAuthorization(row.Reason))
	case spec.AuthorizationReasonPaths:
		if len(resolution.Record.Paths) == 0 {
			appendToolingUntyped(result, reference.Path, resolution.Reason)
			return
		}
		appendToolingUnapproved(result, artifact.displayPath, reference.Path, rowLocation, recordLocation, resolution.Reason, claimsOperativeAuthorization(row.Reason))
	default:
		appendToolingUnapproved(result, artifact.displayPath, reference.Path, rowLocation, recordLocation, resolution.Reason, claimsOperativeAuthorization(row.Reason))
	}
}

func authorizationRole(recordPath string) spec.AuthorizationRole {
	const legacyDirectory = "docs/workflow/authorizations"
	if recordPath == legacyDirectory || strings.HasPrefix(recordPath, legacyDirectory+"/") {
		return spec.AuthorizationRoleLegacy
	}
	return spec.AuthorizationRoleSpec
}

func appendToolingUnauthorized(result *Result, artifactPath, slug, recordPath string, rowLocation, recordLocation Location) {
	result.Findings = append(result.Findings, Finding{
		Code:     CodeToolingUnauthorized,
		Severity: SeverityError,
		Summary:  artifactPath + " cites " + recordPath + ", but that record does not name Spec " + slug,
		Where:    []Location{rowLocation, recordLocation},
		Fix:      "Add Spec " + slug + " to " + recordPath + " or cite the authorization record that already names it.",
	})
}

func appendToolingUntyped(result *Result, recordPath string, reason spec.AuthorizationReason) {
	result.Findings = append(result.Findings, Finding{
		Code:     CodeToolingUntyped,
		Severity: SeverityError,
		Summary:  recordPath + " does not carry an enumerable typed grant: " + reason.Detail,
		Where:    []Location{{Path: recordPath, Line: 1}},
		Fix:      "Open " + recordPath + " with typed authorization frontmatter carrying status, granted, action, consuming, and paths.",
	})
}

func appendToolingUnapproved(
	result *Result,
	artifactPath, recordPath string,
	rowLocation, recordLocation Location,
	reason spec.AuthorizationReason,
	claimsGrant bool,
) {
	if !claimsGrant {
		return
	}
	result.Findings = append(result.Findings, Finding{
		Code:     CodeToolingUnapproved,
		Severity: SeverityError,
		Summary:  artifactPath + " claims express maintainer authorization from " + recordPath + ", but field " + reason.Field + " withholds the grant: " + reason.Detail,
		Where:    []Location{rowLocation, recordLocation},
		Fix:      "Make field " + reason.Field + " operative in " + recordPath + " before claiming authorization, or declare the mutation proposed.",
	})
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
