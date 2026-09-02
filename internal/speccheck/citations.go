package speccheck

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"gopkg.in/yaml.v3"

	"roundfix/internal/spec"
)

const (
	// CodeADRUnlisted identifies a Spec citation absent from its Active ADR obligations row.
	CodeADRUnlisted = "SC-ADR-UNLISTED"
	// CodeADRRelated identifies an unlisted accepted ADR one citation edge from a listed ADR.
	CodeADRRelated = "SC-ADR-RELATED"
	// CodeCitationUnsupported identifies a claim whose cited ADR does not carry the claimed subject.
	CodeCitationUnsupported = "SC-CITATION-UNSUPPORTED"
	// CodeCoverageUnmapped identifies a PRD unit absent from the TechSpec Coverage Map.
	CodeCoverageUnmapped = "SC-COVERAGE-UNMAPPED"
	// CodeCoverageUntasked identifies a PRD unit absent from every Task References section.
	CodeCoverageUntasked = "SC-COVERAGE-UNTASKED"
	// CodeReferenceUnresolved identifies a declared repository path that does not resolve.
	CodeReferenceUnresolved = "SC-REF-UNRESOLVED"
	// CodeLoopOrderDivergent identifies disagreeing declared Spec loop orders.
	CodeLoopOrderDivergent = "SC-LOOP-ORDER-DIVERGENT"
	// CodeFindingLifecycle identifies an active Finding without an allowed written lifecycle status.
	CodeFindingLifecycle = "SC-FINDING-LIFECYCLE"
	// CodeRollupMember identifies a declared Rollup member that resolves to no active or archived Finding.
	CodeRollupMember = "SC-ROLLUP-MEMBER"
	// CodeArchiveLicense identifies an archived Finding without a resolvable written absorbed_by pointer.
	CodeArchiveLicense = "SC-ARCHIVE-LICENSE"
)

var (
	citationCoverageDetectorCodes = []string{
		CodeADRUnlisted,
		CodeADRRelated,
		CodeCitationUnsupported,
		CodeCoverageUnmapped,
		CodeCoverageUntasked,
		CodeReferenceUnresolved,
		CodeVerifyWorkIndependent,
		CodeVerifyInvertedExit,
		CodeVerifyNonHermetic,
		CodeRequirementContradictory,
		CodeRehearsalUndeclared,
		CodeWaveCollision,
	}
	adrFilenamePattern     = regexp.MustCompile(`^([0-9]{4})-.*\.md$`)
	findingFilenamePattern = regexp.MustCompile(`^[0-9]{4}-[0-9]{2}-[0-9]{2}-.+\.md$`)
	adrCitationPattern     = regexp.MustCompile(`\bADR-([0-9]{4})\b`)
	// A bare decision number is zero-padded, which is what separates it from a
	// four-digit year in the same prose. Accepting every four-digit token made
	// "applicable in 2026" cite ADR-2026 and suppress SC-ADR-UNLISTED, turning a
	// checker into one that reports less than it appears to. Revisit if decision
	// numbering ever reaches 1000.
	decisionNumberPattern = regexp.MustCompile(`\b(0[0-9]{3})\b`)
	adrAttributionPattern = regexp.MustCompile(`(?i)\bADR-([0-9]{4})\s+(?:already\s+)?(?:makes?|establish(?:es|ed)?|requires?|keeps?|has|places?|puts?|says?)\s+`)
	citationWordPattern   = regexp.MustCompile(`[a-z0-9]+`)
	inactiveStatusPattern = regexp.MustCompile(`(?im)^\s*(?:\*\*)?status(?:\*\*)?:\s*(?:proposed|rejected|deprecated|superseded)\b`)
	featureRefPattern     = regexp.MustCompile(`(?i)\bCore Features?\s+`)
	storyRefPattern       = regexp.MustCompile(`(?i)\b(?:User )?(?:Story|Stories)\s+`)
	numberedItemPattern   = regexp.MustCompile(`^\s*([0-9]+)\.\s+`)
)

const (
	loopOrderShippedClausePath   = "internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/autonomous-work.md"
	loopOrderRepositoryGuidePath = "docs/agents/autonomous-work.md"
	loopOrderBaselineModulePath  = "internal/baseline/assets/modules/autonomous-work.json"
	loopOrderClauseID            = "clause.autonomous.loop-01-qa-once"
)

type loopOrderSource struct {
	name       string
	path       string
	marker     string
	lineNeedle string
	module     bool
}

type loopOrderDeclaration struct {
	source    loopOrderSource
	order     string
	canonical string
	line      int
	found     bool
}

var loopOrderSources = []loopOrderSource{
	{
		name:       "shipped clause",
		path:       loopOrderShippedClausePath,
		marker:     "Follow one order per Spec:",
		lineNeedle: "Follow one order per Spec:",
	},
	{
		name: "repository guide",
		path: loopOrderRepositoryGuidePath,
		// After greenfield adoption the guide is the rendered catalog clause, so
		// it carries the shipped wording rather than a hand-written paraphrase.
		marker:     "Follow one order per Spec:",
		lineNeedle: "Follow one order per Spec:",
	},
	{
		name:       "Baseline module asset",
		path:       loopOrderBaselineModulePath,
		marker:     "Follow one order per Spec:",
		lineNeedle: `"guidance": "Follow one order per Spec:`,
		module:     true,
	},
}

type findingFrontmatterValue struct {
	value string
	line  int
}

type findingFrontmatter struct {
	status     findingFrontmatterValue
	hasStatus  bool
	kind       findingFrontmatterValue
	hasKind    bool
	members    []findingFrontmatterValue
	absorbedBy findingFrontmatterValue
	hasLicense bool
}

type findingDocument struct {
	name        string
	displayPath string
	frontmatter findingFrontmatter
}

var validFindingLifecycle = map[string]bool{
	"pending":  true,
	"partial":  true,
	"deferred": true,
	"done":     true,
}

func detectFindingsConsistency(result *Result, repoRoot string) error {
	const findingsPath = "docs/findings"
	active, present, err := readFindingDocuments(repoRoot, findingsPath)
	if err != nil {
		return err
	}
	archivePath := spec.ArchiveDir(spec.ArchiveKindFinding)
	archived, archivePresent, err := readFindingDocuments(repoRoot, archivePath)
	if err != nil {
		return err
	}
	if !present && !archivePresent {
		addSkip(result, CodeFindingLifecycle, findingsPath)
		addSkip(result, CodeRollupMember, findingsPath)
		addSkip(result, CodeArchiveLicense, archivePath)
		return nil
	}

	if !present {
		addSkip(result, CodeFindingLifecycle, findingsPath)
	} else if len(active) == 0 {
		addSkip(result, CodeFindingLifecycle, findingsPath+" Finding")
	} else {
		detectFindingLifecycle(result, active)
	}

	activeNames := findingDocumentNames(active)
	archivedNames := findingDocumentNames(archived)
	rollups := findingRollups(active)
	if len(rollups) == 0 {
		addSkip(result, CodeRollupMember, findingsPath+" rollup")
	} else {
		detectRollupMembers(result, rollups, activeNames, archivedNames)
	}

	if !archivePresent || len(archived) == 0 {
		addSkip(result, CodeArchiveLicense, archivePath)
		return nil
	}
	activeSpecs, err := repositoryDirectoryNames(filepath.Join(filepath.Clean(repoRoot), "docs", "specs"), true)
	if err != nil {
		return err
	}
	archivedSpecs, err := repositoryDirectoryNames(
		filepath.Join(filepath.Clean(repoRoot), filepath.FromSlash(spec.ArchiveDir(spec.ArchiveKindSpec))),
		false,
	)
	if err != nil {
		return err
	}
	detectArchiveLicenses(result, archived, findingDocumentNames(rollups), activeSpecs, archivedSpecs)
	return nil
}

func readFindingDocuments(repoRoot, relativeDir string) ([]findingDocument, bool, error) {
	directory := filepath.Join(filepath.Clean(repoRoot), filepath.FromSlash(relativeDir))
	entries, err := os.ReadDir(directory)
	if errors.Is(err, os.ErrNotExist) {
		return nil, false, nil
	}
	if err != nil {
		return nil, false, fmt.Errorf("read Findings directory %q: %w", directory, err)
	}

	documents := make([]findingDocument, 0, len(entries))
	for _, entry := range entries {
		if entry.IsDir() || !findingFilenamePattern.MatchString(entry.Name()) {
			continue
		}
		path := filepath.Join(directory, entry.Name())
		content, err := os.ReadFile(path)
		if err != nil {
			return nil, false, fmt.Errorf("read Finding %q: %w", path, err)
		}
		frontmatter, err := parseFindingFrontmatter(content)
		if err != nil {
			return nil, false, fmt.Errorf("parse Finding %q: %w", path, err)
		}
		documents = append(documents, findingDocument{
			name:        entry.Name(),
			displayPath: artifactDisplayPath(repoRoot, path),
			frontmatter: frontmatter,
		})
	}
	return documents, true, nil
}

func parseFindingFrontmatter(content []byte) (findingFrontmatter, error) {
	const opening = "---\n"
	text := string(content)
	if !strings.HasPrefix(text, opening) {
		return findingFrontmatter{}, nil
	}
	rest := text[len(opening):]
	end := strings.Index(rest, "\n---")
	if end < 0 {
		return findingFrontmatter{}, nil
	}

	var document yaml.Node
	if err := yaml.Unmarshal([]byte(rest[:end]), &document); err != nil {
		return findingFrontmatter{}, fmt.Errorf("parse YAML frontmatter: %w", err)
	}
	if len(document.Content) == 0 || len(document.Content[0].Content) == 0 {
		return findingFrontmatter{}, nil
	}
	mapping := document.Content[0]
	if mapping.Kind != yaml.MappingNode {
		return findingFrontmatter{}, errors.New("YAML frontmatter must be a mapping")
	}

	var frontmatter findingFrontmatter
	for index := 0; index+1 < len(mapping.Content); index += 2 {
		key := mapping.Content[index]
		value := mapping.Content[index+1]
		switch key.Value {
		case "status":
			decoded, err := findingScalar(value)
			if err != nil {
				return findingFrontmatter{}, fmt.Errorf("read status: %w", err)
			}
			frontmatter.status = findingFrontmatterValue{value: decoded, line: value.Line + 1}
			frontmatter.hasStatus = true
		case "kind":
			decoded, err := findingScalar(value)
			if err != nil {
				return findingFrontmatter{}, fmt.Errorf("read kind: %w", err)
			}
			frontmatter.kind = findingFrontmatterValue{value: decoded, line: value.Line + 1}
			frontmatter.hasKind = true
		case "members":
			if value.Kind != yaml.SequenceNode {
				return findingFrontmatter{}, errors.New("members must be a YAML sequence")
			}
			for _, member := range value.Content {
				decoded, err := findingScalar(member)
				if err != nil {
					return findingFrontmatter{}, fmt.Errorf("read member: %w", err)
				}
				frontmatter.members = append(frontmatter.members, findingFrontmatterValue{value: decoded, line: member.Line + 1})
			}
		case "absorbed_by":
			decoded, err := findingScalar(value)
			if err != nil {
				return findingFrontmatter{}, fmt.Errorf("read absorbed_by: %w", err)
			}
			frontmatter.absorbedBy = findingFrontmatterValue{value: decoded, line: value.Line + 1}
			frontmatter.hasLicense = true
		}
	}
	return frontmatter, nil
}

func findingScalar(node *yaml.Node) (string, error) {
	if node.Kind != yaml.ScalarNode {
		return "", fmt.Errorf("expected scalar, found YAML kind %d", node.Kind)
	}
	var value string
	if err := node.Decode(&value); err != nil {
		return "", err
	}
	return value, nil
}

func detectFindingLifecycle(result *Result, documents []findingDocument) {
	for _, document := range documents {
		if !document.frontmatter.hasStatus {
			result.Findings = append(result.Findings, Finding{
				Code:     CodeFindingLifecycle,
				Severity: SeverityError,
				Summary:  document.displayPath + " has no lifecycle status",
				Where:    []Location{{Path: document.displayPath, Line: 1}},
				Fix:      "Add status: pending, partial, deferred, or done to the YAML frontmatter in " + document.displayPath + ".",
			})
			continue
		}
		if validFindingLifecycle[document.frontmatter.status.value] {
			continue
		}
		result.Findings = append(result.Findings, Finding{
			Code:     CodeFindingLifecycle,
			Severity: SeverityError,
			Summary:  document.displayPath + " declares unknown lifecycle status " + strconv.Quote(document.frontmatter.status.value),
			Where:    []Location{{Path: document.displayPath, Line: document.frontmatter.status.line}},
			Fix:      "Replace the status in " + document.displayPath + " with exactly one of pending, partial, deferred, or done.",
		})
	}
}

func findingRollups(documents []findingDocument) []findingDocument {
	rollups := make([]findingDocument, 0)
	for _, document := range documents {
		if document.frontmatter.hasKind && document.frontmatter.kind.value == "rollup" {
			rollups = append(rollups, document)
		}
	}
	return rollups
}

func findingDocumentNames(documents []findingDocument) map[string]bool {
	names := make(map[string]bool, len(documents))
	for _, document := range documents {
		names[document.name] = true
	}
	return names
}

func detectRollupMembers(result *Result, rollups []findingDocument, active, archived map[string]bool) {
	for _, rollup := range rollups {
		for _, member := range rollup.frontmatter.members {
			if active[member.value] || archived[member.value] {
				continue
			}
			result.Findings = append(result.Findings, Finding{
				Code:     CodeRollupMember,
				Severity: SeverityError,
				Summary:  rollup.displayPath + " declares unresolved member " + strconv.Quote(member.value),
				Where:    []Location{{Path: rollup.displayPath, Line: member.line}},
				Fix:      "Restore the declared member " + strconv.Quote(member.value) + " under docs/findings/ or " + spec.ArchiveDir(spec.ArchiveKindFinding) + "/, or update members in " + rollup.displayPath + ".",
			})
		}
	}
}

func repositoryDirectoryNames(directory string, skipUnderscore bool) (map[string]bool, error) {
	names := make(map[string]bool)
	entries, err := os.ReadDir(directory)
	if errors.Is(err, os.ErrNotExist) {
		return names, nil
	}
	if err != nil {
		return nil, fmt.Errorf("read Spec directory %q: %w", directory, err)
	}
	for _, entry := range entries {
		if !entry.IsDir() || (skipUnderscore && strings.HasPrefix(entry.Name(), "_")) {
			continue
		}
		names[entry.Name()] = true
	}
	return names, nil
}

func detectArchiveLicenses(result *Result, archived []findingDocument, rollups, activeSpecs, archivedSpecs map[string]bool) {
	for _, document := range archived {
		if !document.frontmatter.hasLicense {
			result.Findings = append(result.Findings, Finding{
				Code:     CodeArchiveLicense,
				Severity: SeverityError,
				Summary:  document.displayPath + " has no absorbed_by license",
				Where:    []Location{{Path: document.displayPath, Line: 1}},
				Fix:      "Add absorbed_by to the YAML frontmatter in " + document.displayPath + ", naming an active Rollup basename or an active or archived Spec slug.",
			})
			continue
		}
		license := document.frontmatter.absorbedBy
		if rollups[license.value] || activeSpecs[license.value] || archivedSpecs[license.value] {
			continue
		}
		result.Findings = append(result.Findings, Finding{
			Code:     CodeArchiveLicense,
			Severity: SeverityError,
			Summary:  document.displayPath + " declares unresolved absorbed_by " + strconv.Quote(license.value),
			Where:    []Location{{Path: document.displayPath, Line: license.line}},
			Fix:      "Point absorbed_by in " + document.displayPath + " to an active Rollup basename or an active or archived Spec slug.",
		})
	}
}

func detectLoopOrderConsistency(result *Result, repoRoot string) error {
	declarations := make([]loopOrderDeclaration, 0, len(loopOrderSources))
	missingSource := false
	for _, source := range loopOrderSources {
		declaration, present, err := readLoopOrderDeclaration(repoRoot, source)
		if err != nil {
			return err
		}
		if !present {
			addSkip(result, CodeLoopOrderDivergent, source.path)
			missingSource = true
			continue
		}
		declarations = append(declarations, declaration)
	}
	if missingSource {
		return nil
	}

	reference := declarations[0]
	divergent := !reference.found
	for _, declaration := range declarations[1:] {
		if !declaration.found || declaration.canonical != reference.canonical {
			divergent = true
		}
	}
	if !divergent {
		return nil
	}

	descriptions := make([]string, 0, len(declarations))
	locations := make([]Location, 0, len(declarations))
	for _, declaration := range declarations {
		order := declaration.order
		if !declaration.found {
			order = "<declared order not found>"
		}
		descriptions = append(descriptions, fmt.Sprintf("%s declares %q", declaration.source.name, order))
		locations = append(locations, Location{Path: declaration.source.path, Line: declaration.line})
	}
	result.Findings = append(result.Findings, Finding{
		Code:     CodeLoopOrderDivergent,
		Severity: SeverityError,
		Summary:  "loop order sources disagree: " + strings.Join(descriptions, "; "),
		Where:    locations,
		Fix:      "Make the shipped clause, repository guide, and Baseline module asset declare the same ordered actions.",
	})
	return nil
}

func readLoopOrderDeclaration(repoRoot string, source loopOrderSource) (loopOrderDeclaration, bool, error) {
	path := filepath.Join(filepath.Clean(repoRoot), filepath.FromSlash(source.path))
	content, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return loopOrderDeclaration{}, false, nil
	}
	if err != nil {
		return loopOrderDeclaration{}, false, fmt.Errorf("read loop order source %q: %w", path, err)
	}

	statement := string(content)
	if source.module {
		guidance, err := loopOrderModuleGuidance(content)
		if err != nil {
			return loopOrderDeclaration{}, false, fmt.Errorf("read loop order source %q: %w", path, err)
		}
		statement = guidance
	}
	order, found := declaredLoopOrder(statement, source.marker)
	return loopOrderDeclaration{
		source:    source,
		order:     order,
		canonical: canonicalLoopOrder(order),
		line:      lineContainingText(content, source.lineNeedle),
		found:     found,
	}, true, nil
}

func loopOrderModuleGuidance(content []byte) (string, error) {
	var module struct {
		Rules []struct {
			Clauses []struct {
				ID       string `json:"id"`
				Guidance string `json:"guidance"`
			} `json:"clauses"`
		} `json:"rules"`
	}
	if err := json.Unmarshal(content, &module); err != nil {
		return "", fmt.Errorf("parse Baseline module JSON: %w", err)
	}
	for _, rule := range module.Rules {
		for _, clause := range rule.Clauses {
			if clause.ID == loopOrderClauseID {
				return clause.Guidance, nil
			}
		}
	}
	return "", nil
}

func declaredLoopOrder(statement, marker string) (string, bool) {
	normalized := strings.Join(strings.Fields(statement), " ")
	markerIndex := strings.Index(strings.ToLower(normalized), strings.ToLower(marker))
	if markerIndex < 0 {
		return "", false
	}
	orderStart := markerIndex + len(marker)
	remainder := strings.TrimSpace(normalized[orderStart:])
	orderEnd := strings.Index(remainder, ".")
	if orderEnd < 0 {
		return "", false
	}
	order := strings.TrimSpace(remainder[:orderEnd])
	return order, order != ""
}

func canonicalLoopOrder(order string) string {
	return strings.ToLower(strings.Join(strings.Fields(order), " "))
}

func lineContainingText(content []byte, needle string) int {
	index := strings.Index(string(content), needle)
	if index < 0 {
		return 1
	}
	return strings.Count(string(content[:index]), "\n") + 1
}

type adrRecord struct {
	Number        string
	Title         string
	TitleLine     int
	Text          string
	DisplayPath   string
	Citations     map[string]bool
	CitationLines map[string]int
}

// Claim is one attribution a Spec artifact makes about a decision record.
type Claim struct {
	Artifact string
	Line     int
	Target   string
	Subject  string
	sentence string
}

type resolvedCitationClaim struct {
	claim  Claim
	record adrRecord
}

// CitationClaims parses subject attributions from one PRD or TechSpec. A bare
// ADR token is not a claim because it has no attribution verb or subject.
func CitationClaims(artifact string, content []byte) []Claim {
	var claims []Claim
	var paragraph []string
	paragraphLine := 0
	inFence := false

	flush := func() {
		if len(paragraph) == 0 {
			return
		}
		claims = append(claims, citationClaimsInParagraph(artifact, paragraphLine, strings.Join(paragraph, "\n"))...)
		paragraph = nil
		paragraphLine = 0
	}

	for index, line := range strings.Split(string(content), "\n") {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "```") {
			flush()
			inFence = !inFence
			continue
		}
		if inFence {
			continue
		}
		if trimmed == "" {
			flush()
			continue
		}
		if len(paragraph) > 0 && strings.HasPrefix(trimmed, "- ") {
			flush()
		}
		if len(paragraph) == 0 {
			paragraphLine = index + 1
		}
		paragraph = append(paragraph, line)
	}
	flush()
	return claims
}

func citationClaimsInParagraph(artifact string, firstLine int, paragraph string) []Claim {
	var claims []Claim
	for _, match := range adrAttributionPattern.FindAllStringSubmatchIndex(paragraph, -1) {
		if citationIsExample(paragraph, match[0]) {
			continue
		}
		subjectStart := match[1]
		subjectEnd := citationSubjectEnd(paragraph, subjectStart)
		subject := strings.TrimSpace(paragraph[subjectStart:subjectEnd])
		if subject == "" {
			continue
		}
		target := "ADR-" + paragraph[match[2]:match[3]]
		claims = append(claims, Claim{
			Artifact: artifact,
			Line:     firstLine + strings.Count(paragraph[:match[0]], "\n"),
			Target:   target,
			Subject:  strings.Join(strings.Fields(subject), " "),
			sentence: citationSentence(paragraph, match[0]),
		})
	}
	return claims
}

func citationIsExample(paragraph string, claimStart int) bool {
	prefix := strings.TrimSpace(paragraph[:claimStart])
	if prefix != "" {
		last := prefix[len(prefix)-1]
		if last == '"' || last == '\'' {
			return true
		}
	}
	lower := strings.ToLower(prefix)
	if len(lower) > 64 {
		lower = lower[len(lower)-64:]
	}
	return strings.Contains(lower, "claiming") || strings.Contains(lower, "example")
}

func citationSubjectEnd(paragraph string, start int) int {
	remainder := paragraph[start:]
	end := len(remainder)
	for _, delimiter := range []string{",", ";", ".", "!", "?", " — ", "\nADR-", " ADR-"} {
		if index := strings.Index(remainder, delimiter); index >= 0 && index < end {
			end = index
		}
	}
	return start + end
}

func citationSentence(paragraph string, claimStart int) string {
	start := 0
	for index := claimStart - 1; index >= 0; index-- {
		if strings.ContainsRune(".!?", rune(paragraph[index])) {
			start = index + 1
			break
		}
	}
	end := len(paragraph)
	for index := claimStart; index < len(paragraph); index++ {
		if strings.ContainsRune(".!?", rune(paragraph[index])) {
			end = index + 1
			break
		}
	}
	return strings.Join(strings.Fields(strings.TrimSpace(paragraph[start:end])), " ")
}

// ResolvedCitationClaimCount reports how many parsed claims name an accepted
// record in the repository corpus. Callers can distinguish no attributions
// from a parser that failed to recognize a known attribution by comparing this
// count with CitationClaims.
func ResolvedCitationClaimCount(repoRoot string, claims []Claim) (int, error) {
	resolved, err := resolveCitationClaims(repoRoot, claims)
	if err != nil {
		return 0, err
	}
	return len(resolved), nil
}

func resolveCitationClaims(repoRoot string, claims []Claim) ([]resolvedCitationClaim, error) {
	corpus, present, err := readADRCorpus(repoRoot)
	if err != nil {
		return nil, err
	}
	if !present {
		return nil, nil
	}
	resolved := make([]resolvedCitationClaim, 0, len(claims))
	for _, claim := range claims {
		record, ok := corpus[strings.TrimPrefix(claim.Target, "ADR-")]
		if !ok {
			continue
		}
		resolved = append(resolved, resolvedCitationClaim{claim: claim, record: record})
	}
	return resolved, nil
}

func detectUnsupportedCitations(result *Result, repoRoot string, claims []Claim) error {
	resolved, err := resolveCitationClaims(repoRoot, claims)
	if err != nil {
		return err
	}
	for _, resolution := range resolved {
		if citationSubjectSupported(resolution.claim.Subject, resolution.record) {
			continue
		}
		claimingSentence := resolution.claim.sentence
		if claimingSentence == "" {
			claimingSentence = resolution.claim.Target + " establishes " + resolution.claim.Subject
		}
		result.Findings = append(result.Findings, Finding{
			Code:     CodeCitationUnsupported,
			Severity: SeverityError,
			Summary: resolution.claim.Artifact + " claims " + strconv.Quote(claimingSentence) +
				", but " + resolution.claim.Target + "'s subject is " + strconv.Quote(resolution.record.Title),
			Where: []Location{
				{Path: resolution.claim.Artifact, Line: resolution.claim.Line},
				{Path: resolution.record.DisplayPath, Line: resolution.record.TitleLine},
			},
			Fix: "Rewrite the claim in " + resolution.claim.Artifact + " to match " + resolution.claim.Target +
				", or cite the decision record that carries the claimed subject.",
		})
	}
	return nil
}

func citationSubjectSupported(subject string, record adrRecord) bool {
	subjectWords := significantCitationWords(subject)
	if len(subjectWords) == 0 {
		return true
	}
	recordWordsInOrder := significantCitationWords(record.Text)
	recordWords := make(map[string]bool)
	for _, word := range recordWordsInOrder {
		recordWords[word] = true
	}
	matched := citationWordOverlap(subjectWords, recordWords)
	if len(subjectWords) <= 2 {
		return matched == len(subjectWords)
	}
	if matched*5 < len(subjectWords)*3 {
		return false
	}
	titleWords := make(map[string]bool)
	for _, word := range significantCitationWords(record.Title) {
		titleWords[word] = true
	}
	return citationWordOverlap(subjectWords, titleWords) >= 2 || citationPhraseMatches(subjectWords, recordWordsInOrder)
}

func citationWordOverlap(words []string, corpus map[string]bool) int {
	matched := 0
	for _, word := range words {
		if corpus[word] {
			matched++
		}
	}
	return matched
}

func citationPhraseMatches(subjectWords, recordWords []string) bool {
	pairs := make(map[string]bool)
	for index := 0; index+1 < len(recordWords); index++ {
		pairs[recordWords[index]+"\x00"+recordWords[index+1]] = true
	}
	for index := 0; index+1 < len(subjectWords); index++ {
		if pairs[subjectWords[index]+"\x00"+subjectWords[index+1]] {
			return true
		}
	}
	return false
}

func significantCitationWords(text string) []string {
	var words []string
	for _, word := range citationWordPattern.FindAllString(strings.ToLower(text), -1) {
		if citationStopWords[word] {
			continue
		}
		words = append(words, citationWordStem(word))
	}
	return words
}

func citationWordStem(word string) string {
	switch {
	case len(word) > 4 && strings.HasSuffix(word, "ies"):
		return strings.TrimSuffix(word, "ies") + "y"
	case len(word) > 5 && strings.HasSuffix(word, "ing"):
		return strings.TrimSuffix(word, "ing")
	case len(word) > 4 && strings.HasSuffix(word, "ed"):
		return strings.TrimSuffix(word, "ed")
	case len(word) > 3 && strings.HasSuffix(word, "s"):
		return strings.TrimSuffix(word, "s")
	default:
		return word
	}
}

var citationStopWords = map[string]bool{
	"a": true, "an": true, "and": true, "as": true, "at": true,
	"be": true, "by": true, "for": true, "from": true, "in": true,
	"is": true, "it": true, "its": true, "of": true, "on": true,
	"or": true, "that": true, "the": true, "this": true, "to": true,
	"with": true,
}

type coverageKind string

const (
	coverageFeature coverageKind = "Core Feature"
	coverageStory   coverageKind = "User Story"
)

type coverageUnit struct {
	Kind   coverageKind
	Number int
	Line   int
}

type referenceSet map[coverageKind]map[int]bool

func detectCitationCoverageAndReferences(
	result *Result,
	specsRoot string,
	repoRoot string,
	slug string,
	specDir string,
	techSpecPresent bool,
) error {
	prdPath := filepath.Join(specDir, "_prd.md")
	prdContent, err := os.ReadFile(prdPath)
	if err != nil {
		return fmt.Errorf("read Spec artifact %q: %w", prdPath, err)
	}
	prdDisplayPath := artifactDisplayPath(repoRoot, prdPath)

	if err := detectADRConsistency(result, repoRoot, specDir, prdContent, prdDisplayPath); err != nil {
		return err
	}
	claims := CitationClaims(prdDisplayPath, prdContent)

	units := parsePRDCoverageUnits(prdContent)
	if techSpecPresent {
		techSpecPath := filepath.Join(specDir, "_techspec.md")
		techSpecContent, err := os.ReadFile(techSpecPath)
		if err != nil {
			return fmt.Errorf("read Spec artifact %q: %w", techSpecPath, err)
		}
		claims = append(claims, CitationClaims(artifactDisplayPath(repoRoot, techSpecPath), techSpecContent)...)
		detectCoverageMap(result, units, prdDisplayPath, techSpecContent, artifactDisplayPath(repoRoot, techSpecPath))
	} else {
		addSkip(result, CodeCoverageUnmapped, artifactDisplayPath(repoRoot, filepath.Join(specDir, "_techspec.md")))
	}
	if err := detectUnsupportedCitations(result, repoRoot, claims); err != nil {
		return fmt.Errorf("detect unsupported citations: %w", err)
	}

	graph, graphPresent, err := loadOptionalTaskGraph(specsRoot, slug, specDir)
	if err != nil {
		return err
	}
	if graphPresent {
		if err := detectWaveCollisions(result, repoRoot, graph); err != nil {
			return err
		}
		if err := detectTaskCoverageAndContextReferences(result, repoRoot, specsRoot, graph, units, prdDisplayPath); err != nil {
			return err
		}
	} else {
		manifestDisplayPath := artifactDisplayPath(repoRoot, filepath.Join(specDir, "_tasks.md"))
		addSkip(result, CodeCoverageUntasked, manifestDisplayPath)
		addSkip(result, CodeReferenceUnresolved, manifestDisplayPath)
		addSkip(result, CodeVerifyInvertedExit, manifestDisplayPath)
		addSkip(result, CodeVerifyNonHermetic, manifestDisplayPath)
		addSkip(result, CodeWaveCollision, manifestDisplayPath)
	}

	if err := detectReferenceIndex(result, repoRoot, specDir); err != nil {
		return err
	}
	return nil
}

func detectADRConsistency(result *Result, repoRoot, specDir string, prdContent []byte, prdDisplayPath string) error {
	corpus, present, err := readADRCorpus(repoRoot)
	if err != nil {
		return err
	}
	if !present {
		missing := filepath.ToSlash(filepath.Join("docs", "adr"))
		addSkip(result, CodeADRUnlisted, missing)
		addSkip(result, CodeADRRelated, missing)
		addSkip(result, CodeCitationUnsupported, missing)
		return nil
	}

	rowText, rowLine, rowPresent := projectConstraintEntry(prdContent, constraintActiveADR)
	if !rowPresent {
		missing := prdDisplayPath + " Active ADR obligations row"
		addSkip(result, CodeADRUnlisted, missing)
		addSkip(result, CodeADRRelated, missing)
		return nil
	}
	listed := obligationCitationNumbers([]byte(rowText))
	rowLocation := Location{Path: prdDisplayPath, Line: rowLine}

	specCitations, err := readSpecCitations(repoRoot, specDir)
	if err != nil {
		return err
	}
	for _, number := range sortedCitationNumbers(specCitations) {
		if listed[number] {
			continue
		}
		citation := specCitations[number]
		title := ""
		if record, ok := corpus[number]; ok {
			title = record.Title
		}
		name := adrName(number, title)
		result.Findings = append(result.Findings, Finding{
			Code:     CodeADRUnlisted,
			Severity: SeverityError,
			Summary:  citation.Path + " cites " + name + ", but " + prdDisplayPath + " does not list it under Active ADR obligations",
			Where:    []Location{citation, rowLocation},
			Fix:      "List ADR-" + number + " in the Active ADR obligations row in " + prdDisplayPath + " using the recognised ADR-NNNN form, or remove the stale citation.",
		})
	}

	for _, number := range sortedADRNumbers(corpus) {
		record := corpus[number]
		if listed[number] {
			continue
		}
		relatedNumber, related := firstListedCitation(record, listed)
		if !related {
			continue
		}
		result.Findings = append(result.Findings, Finding{
			Code:     CodeADRRelated,
			Severity: SeverityGap,
			Summary:  adrName(number, record.Title) + " cites listed ADR-" + relatedNumber + " but is not itself listed by the Spec",
			Where: []Location{
				{Path: record.DisplayPath, Line: record.CitationLines[relatedNumber]},
				rowLocation,
			},
			Fix: "Account for ADR-" + number + " in the Active ADR obligations row in " + prdDisplayPath + " by listing it or recording why it does not apply.",
		})
	}
	return nil
}

func readADRCorpus(repoRoot string) (map[string]adrRecord, bool, error) {
	adrDir := filepath.Join(filepath.Clean(repoRoot), "docs", "adr")
	entries, err := os.ReadDir(adrDir)
	if errors.Is(err, os.ErrNotExist) {
		return nil, false, nil
	}
	if err != nil {
		return nil, false, fmt.Errorf("read ADR corpus %q: %w", adrDir, err)
	}
	corpus := make(map[string]adrRecord)
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		match := adrFilenamePattern.FindStringSubmatch(entry.Name())
		if len(match) != 2 {
			continue
		}
		path := filepath.Join(adrDir, entry.Name())
		content, err := os.ReadFile(path)
		if err != nil {
			return nil, false, fmt.Errorf("read ADR %q: %w", path, err)
		}
		status, body, err := adrStatusAndBody(content)
		if err != nil {
			return nil, false, fmt.Errorf("parse ADR %q: %w", path, err)
		}
		if status != "accepted" {
			continue
		}
		citations, lines := citationsWithLines(body)
		lineOffset := strings.Count(string(content[:len(content)-len(body)]), "\n")
		for number := range lines {
			lines[number] += lineOffset
		}
		title := firstHeading(body)
		corpus[match[1]] = adrRecord{
			Number:        match[1],
			Title:         title,
			TitleLine:     lineOffset + lineContainingText(body, "# "+title),
			Text:          string(body),
			DisplayPath:   artifactDisplayPath(repoRoot, path),
			Citations:     citations,
			CitationLines: lines,
		}
	}
	return corpus, true, nil
}

func adrStatusAndBody(content []byte) (string, []byte, error) {
	text := string(content)
	const opening = "---\n"
	if !strings.HasPrefix(text, opening) {
		if inactiveStatusPattern.Match(content) {
			return "inactive", content, nil
		}
		return "accepted", content, nil
	}
	rest := text[len(opening):]
	end := strings.Index(rest, "\n---")
	if end < 0 {
		return "", nil, errors.New("missing YAML frontmatter closing marker")
	}
	var frontmatter struct {
		Status string `yaml:"status"`
	}
	if err := yaml.Unmarshal([]byte(rest[:end]), &frontmatter); err != nil {
		return "", nil, fmt.Errorf("parse frontmatter: %w", err)
	}
	bodyStart := end + len("\n---")
	if strings.HasPrefix(rest[bodyStart:], "\n") {
		bodyStart++
	}
	return strings.TrimSpace(frontmatter.Status), []byte(rest[bodyStart:]), nil
}

func firstHeading(content []byte) string {
	for _, line := range strings.Split(string(content), "\n") {
		if strings.HasPrefix(strings.TrimSpace(line), "# ") {
			return strings.TrimSpace(strings.TrimPrefix(strings.TrimSpace(line), "# "))
		}
	}
	return ""
}

func readSpecCitations(repoRoot, specDir string) (map[string]Location, error) {
	citations := make(map[string]Location)
	err := filepath.WalkDir(specDir, func(path string, entry fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if entry.IsDir() || filepath.Ext(entry.Name()) != ".md" {
			return nil
		}
		content, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		displayPath := artifactDisplayPath(repoRoot, path)
		for index, line := range strings.Split(string(content), "\n") {
			for _, match := range adrCitationPattern.FindAllStringSubmatch(line, -1) {
				if _, exists := citations[match[1]]; !exists {
					citations[match[1]] = Location{Path: displayPath, Line: index + 1}
				}
			}
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("read Spec citations in %q: %w", specDir, err)
	}
	return citations, nil
}

// obligationCitationNumbers reads decision numbers only after the caller has
// isolated the Active ADR obligations row. That context makes a bare four-digit
// number unambiguous and leaves list punctuation, including English "and" and
// Portuguese "e", outside the citation token itself.
func obligationCitationNumbers(content []byte) map[string]bool {
	numbers := make(map[string]bool)
	for _, match := range decisionNumberPattern.FindAllSubmatch(content, -1) {
		numbers[string(match[1])] = true
	}
	return numbers
}

func citationsWithLines(content []byte) (map[string]bool, map[string]int) {
	citations := make(map[string]bool)
	lines := make(map[string]int)
	for index, line := range strings.Split(string(content), "\n") {
		for _, match := range adrCitationPattern.FindAllStringSubmatch(line, -1) {
			citations[match[1]] = true
			if lines[match[1]] == 0 {
				lines[match[1]] = index + 1
			}
		}
	}
	return citations, lines
}

func firstListedCitation(record adrRecord, listed map[string]bool) (string, bool) {
	var related []string
	for number := range record.Citations {
		if listed[number] {
			related = append(related, number)
		}
	}
	if len(related) == 0 {
		return "", false
	}
	sort.Strings(related)
	return related[0], true
}

func sortedADRNumbers(corpus map[string]adrRecord) []string {
	numbers := make([]string, 0, len(corpus))
	for number := range corpus {
		numbers = append(numbers, number)
	}
	sort.Strings(numbers)
	return numbers
}

func sortedCitationNumbers(citations map[string]Location) []string {
	numbers := make([]string, 0, len(citations))
	for number := range citations {
		numbers = append(numbers, number)
	}
	sort.Strings(numbers)
	return numbers
}

func adrName(number, title string) string {
	if title == "" {
		return "ADR-" + number
	}
	return "ADR-" + number + " (" + title + ")"
}

func projectConstraintEntry(content []byte, label string) (string, int, bool) {
	lines := strings.Split(string(content), "\n")
	sectionStart := -1
	for index, line := range lines {
		if strings.TrimSpace(line) == "## Project Constraints" {
			sectionStart = index + 1
			break
		}
	}
	if sectionStart == -1 {
		return "", 0, false
	}
	for index := sectionStart; index < len(lines); index++ {
		trimmed := strings.TrimSpace(lines[index])
		if strings.HasPrefix(trimmed, "## ") {
			break
		}
		prefix := "- " + label + ":"
		if !strings.HasPrefix(strings.ToLower(trimmed), strings.ToLower(prefix)) {
			continue
		}
		entry := []string{strings.TrimSpace(strings.TrimPrefix(trimmed, "- "))}
		for next := index + 1; next < len(lines); next++ {
			continuation := strings.TrimSpace(lines[next])
			if continuation == "" || strings.HasPrefix(continuation, "- ") || strings.HasPrefix(continuation, "## ") {
				break
			}
			entry = append(entry, continuation)
		}
		return strings.Join(entry, " "), index + 1, true
	}
	return "", 0, false
}

func parsePRDCoverageUnits(content []byte) []coverageUnit {
	var units []coverageUnit
	units = append(units, parseNumberedSection(content, "User Stories", coverageStory)...)
	units = append(units, parseNumberedSection(content, "Core Features", coverageFeature)...)
	return units
}

func parseNumberedSection(content []byte, heading string, kind coverageKind) []coverageUnit {
	lines := strings.Split(string(content), "\n")
	inSection := false
	var units []coverageUnit
	for index, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "## ") {
			if inSection {
				break
			}
			inSection = strings.TrimSpace(strings.TrimPrefix(trimmed, "## ")) == heading
			continue
		}
		if !inSection {
			continue
		}
		match := numberedItemPattern.FindStringSubmatch(line)
		if len(match) != 2 {
			continue
		}
		number, err := strconv.Atoi(match[1])
		if err == nil {
			units = append(units, coverageUnit{Kind: kind, Number: number, Line: index + 1})
		}
	}
	return units
}

func detectCoverageMap(result *Result, units []coverageUnit, prdDisplayPath string, content []byte, techSpecDisplayPath string) {
	references, sectionLine, present := parseCoverageMapReferences(content)
	if !present {
		sectionLine = 1
	}
	for _, unit := range units {
		if references[unit.Kind][unit.Number] {
			continue
		}
		name := coverageUnitName(unit)
		result.Findings = append(result.Findings, Finding{
			Code:     CodeCoverageUnmapped,
			Severity: SeverityError,
			Summary:  prdDisplayPath + " declares " + name + ", but " + techSpecDisplayPath + " has no Coverage Map entry for it",
			Where: []Location{
				{Path: prdDisplayPath, Line: unit.Line},
				{Path: techSpecDisplayPath, Line: sectionLine},
			},
			Fix: "Add " + name + " to the Coverage Map in " + techSpecDisplayPath + ".",
		})
	}
}

func parseCoverageMapReferences(content []byte) (referenceSet, int, bool) {
	references := newReferenceSet()
	lines := strings.Split(string(content), "\n")
	inSection := false
	sectionLine := 0
	for index, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "## ") {
			if inSection {
				break
			}
			if strings.TrimSpace(strings.TrimPrefix(trimmed, "## ")) == "Coverage Map" {
				inSection = true
				sectionLine = index + 1
			}
			continue
		}
		if inSection {
			addDeclaredReferences(references, line)
		}
	}
	return references, sectionLine, inSection || sectionLine > 0
}

func loadOptionalTaskGraph(specsRoot, slug, specDir string) (*spec.Graph, bool, error) {
	manifestPath := filepath.Join(specDir, "_tasks.md")
	_, err := os.Stat(manifestPath)
	if errors.Is(err, os.ErrNotExist) {
		return nil, false, nil
	}
	if err != nil {
		return nil, false, fmt.Errorf("read Task Graph %q: %w", manifestPath, err)
	}
	graph, err := spec.Load(specsRoot, slug)
	if err != nil {
		return nil, true, fmt.Errorf("load Task Graph for Spec %q: %w", slug, err)
	}
	return graph, true, nil
}

func detectTaskCoverageAndContextReferences(
	result *Result,
	repoRoot string,
	specsRoot string,
	graph *spec.Graph,
	units []coverageUnit,
	prdDisplayPath string,
) error {
	references := newReferenceSet()
	for _, task := range graph.Tasks {
		taskPath := filepath.Join(filepath.Clean(specsRoot), task.File)
		content, err := os.ReadFile(taskPath)
		if err != nil {
			return fmt.Errorf("read Task %q: %w", taskPath, err)
		}
		for _, line := range markdownSectionLines(content, "References") {
			addDeclaredReferences(references, line)
		}
		detectTaskContextReferences(result, repoRoot, taskPath, content, task.Context)
		// Completed Tasks carry historical evidence whose authoring contract is
		// not retroactively changed by a newly shipped detector.
		if task.Status == spec.StatusCompleted {
			continue
		}
		if finding, ok := AuthoredQAVerification(task); ok {
			finding.Where[0] = Location{
				Path: artifactDisplayPath(repoRoot, taskPath),
				Line: sectionLineContaining(content, "Verification", "`"),
			}
			finding.Summary = finding.Where[0].Path + strings.TrimPrefix(finding.Summary, task.File)
			result.Findings = append(result.Findings, finding)
		}
		if finding, ok := WorkIndependentVerification(task); ok {
			finding.Where[0] = Location{
				Path: artifactDisplayPath(repoRoot, taskPath),
				Line: sectionLineContaining(content, "Verification", task.Verification[0]),
			}
			finding.Summary = finding.Where[0].Path + strings.TrimPrefix(finding.Summary, task.File)
			result.Findings = append(result.Findings, finding)
		}
		createdPaths := make(map[string]bool)
		for _, command := range task.Verification {
			commandFindings := InvertedExitVerification(spec.Task{
				File:         task.File,
				Verification: []string{command},
			})
			for _, finding := range commandFindings {
				finding.Where[0] = Location{
					Path: artifactDisplayPath(repoRoot, taskPath),
					Line: sectionLineContaining(content, "Verification", command),
				}
				finding.Summary = finding.Where[0].Path + strings.TrimPrefix(finding.Summary, task.File)
				result.Findings = append(result.Findings, finding)
			}
			if task.Type != spec.TaskTypeQA {
				form, matched := nonHermeticVerificationCommand(command, createdPaths)
				if matched {
					finding := nonHermeticFinding(task.File, command, form)
					finding.Where[0] = Location{
						Path: artifactDisplayPath(repoRoot, taskPath),
						Line: sectionLineContaining(content, "Verification", command),
					}
					finding.Summary = finding.Where[0].Path + strings.TrimPrefix(finding.Summary, task.File)
					result.Findings = append(result.Findings, finding)
				}
			}
		}
		for _, command := range VacuousVerificationCommands(task) {
			result.Findings = append(result.Findings, Finding{
				Code:     CodeVerifyVacuousCommand,
				Severity: SeverityError,
				Summary: artifactDisplayPath(repoRoot, taskPath) +
					" declares a Verification command that already passes against the unchanged tree: " + command,
				Where: []Location{{
					Path: artifactDisplayPath(repoRoot, taskPath),
					Line: sectionLineContaining(content, "Verification", command),
				}},
				Fix: "Replace it with a command that names this Task's own effect, or drop it: a boundary is proved by auditing the Task commit, not by a command that runs before the commit exists.",
			})
		}
		if finding, ok := ContradictoryRequirements(task); ok {
			for index := range finding.Where {
				finding.Where[index].Path = artifactDisplayPath(repoRoot, taskPath)
			}
			finding.Summary = finding.Where[0].Path + strings.TrimPrefix(finding.Summary, task.File)
			result.Findings = append(result.Findings, finding)
		}
		if finding, ok := UndeclaredRehearsal(task); ok {
			finding.Where[0].Path = artifactDisplayPath(repoRoot, taskPath)
			finding.Summary = finding.Where[0].Path + strings.TrimPrefix(finding.Summary, task.File)
			result.Findings = append(result.Findings, finding)
		}
	}

	manifestDisplayPath := artifactDisplayPath(repoRoot, filepath.Join(graph.Spec.Dir, "_tasks.md"))
	for _, unit := range units {
		if references[unit.Kind][unit.Number] {
			continue
		}
		name := coverageUnitName(unit)
		result.Findings = append(result.Findings, Finding{
			Code:     CodeCoverageUntasked,
			Severity: SeverityError,
			Summary:  prdDisplayPath + " declares " + name + ", but no Task in " + manifestDisplayPath + " references it",
			Where: []Location{
				{Path: prdDisplayPath, Line: unit.Line},
				{Path: manifestDisplayPath, Line: 1},
			},
			Fix: "Reference " + name + " from the Task that implements or verifies it.",
		})
	}
	return nil
}

func detectTaskContextReferences(result *Result, repoRoot, taskPath string, content []byte, refs []spec.TaskContextRef) {
	taskDisplayPath := artifactDisplayPath(repoRoot, taskPath)
	for _, ref := range refs {
		if ref.Kind == spec.ContextKindCreates {
			continue
		}
		if repositoryPathExists(repoRoot, ref.Path) {
			continue
		}
		line := sectionLineContaining(content, "Context", ref.Path)
		result.Findings = append(result.Findings, Finding{
			Code:     CodeReferenceUnresolved,
			Severity: SeverityError,
			Summary:  taskDisplayPath + " declares unresolved Task Context path " + ref.Path,
			Where: []Location{
				{Path: taskDisplayPath, Line: line},
				{Path: ref.Path, Line: 1},
			},
			Fix: "Create " + ref.Path + " or update the Task Context entry in " + taskDisplayPath + ".",
		})
	}
}

func detectReferenceIndex(result *Result, repoRoot, specDir string) error {
	indexPath := filepath.Join(specDir, "references", "_index.md")
	content, err := os.ReadFile(indexPath)
	if errors.Is(err, os.ErrNotExist) {
		addSkip(result, CodeReferenceUnresolved, artifactDisplayPath(repoRoot, indexPath))
		return nil
	}
	if err != nil {
		return fmt.Errorf("read reference index %q: %w", indexPath, err)
	}
	indexDisplayPath := artifactDisplayPath(repoRoot, indexPath)
	for _, entry := range parseReferenceIndexEntries(content) {
		resolvedPath, ok := resolveIndexedReference(filepath.Dir(indexPath), entry.Path)
		if ok {
			_, err = os.Stat(resolvedPath)
			ok = err == nil
			if err != nil && !errors.Is(err, os.ErrNotExist) {
				return fmt.Errorf("resolve indexed reference %q: %w", resolvedPath, err)
			}
		}
		if ok {
			continue
		}
		missingDisplayPath := entry.Path
		if resolvedPath != "" {
			missingDisplayPath = artifactDisplayPath(repoRoot, resolvedPath)
		}
		result.Findings = append(result.Findings, Finding{
			Code:     CodeReferenceUnresolved,
			Severity: SeverityError,
			Summary:  indexDisplayPath + " declares unresolved reference path " + entry.Path,
			Where: []Location{
				{Path: indexDisplayPath, Line: entry.Line},
				{Path: missingDisplayPath, Line: 1},
			},
			Fix: "Restore " + entry.Path + " beside the reference index or update the declaring row.",
		})
	}
	return nil
}

type indexedReference struct {
	Path string
	Line int
}

func parseReferenceIndexEntries(content []byte) []indexedReference {
	var entries []indexedReference
	headerFound := false
	for index, line := range strings.Split(string(content), "\n") {
		cells := markdownCells(line)
		if !headerFound {
			if len(cells) == 5 && strings.EqualFold(cells[0], "source") && strings.EqualFold(cells[4], "path") {
				headerFound = true
			}
			continue
		}
		if len(cells) == 0 {
			break
		}
		if markdownSeparator(cells) {
			continue
		}
		if len(cells) >= 5 && strings.TrimSpace(cells[4]) != "" {
			entries = append(entries, indexedReference{Path: strings.TrimSpace(cells[4]), Line: index + 1})
		}
	}
	return entries
}

func resolveIndexedReference(indexDir, relative string) (string, bool) {
	if filepath.IsAbs(relative) || strings.Contains(relative, `\`) {
		return "", false
	}
	clean := filepath.Clean(filepath.FromSlash(relative))
	if clean == "." || clean == ".." || strings.HasPrefix(clean, ".."+string(filepath.Separator)) {
		return "", false
	}
	return filepath.Join(indexDir, clean), true
}

// markdownHeading reports the depth and text of one ATX heading line, and depth
// zero for a line that is not a heading. Depth separates a document's sections
// from the subsections written inside one.
func markdownHeading(line string) (int, string) {
	trimmed := strings.TrimSpace(line)
	depth := 0
	for depth < len(trimmed) && trimmed[depth] == '#' {
		depth++
	}
	if depth == 0 || depth >= len(trimmed) || trimmed[depth] != ' ' {
		return 0, ""
	}
	return depth, strings.TrimSpace(trimmed[depth+1:])
}

func markdownCells(line string) []string {
	trimmed := strings.TrimSpace(line)
	if !strings.HasPrefix(trimmed, "|") || !strings.HasSuffix(trimmed, "|") {
		return nil
	}
	trimmed = strings.TrimPrefix(strings.TrimSuffix(trimmed, "|"), "|")
	parts := strings.Split(trimmed, "|")
	cells := make([]string, 0, len(parts))
	for _, part := range parts {
		cells = append(cells, strings.TrimSpace(part))
	}
	return cells
}

func markdownSeparator(cells []string) bool {
	if len(cells) == 0 {
		return false
	}
	for _, cell := range cells {
		trimmed := strings.TrimSpace(cell)
		if trimmed == "" {
			return false
		}
		for _, char := range trimmed {
			if char != '-' && char != ':' {
				return false
			}
		}
	}
	return true
}

func repositoryPathExists(repoRoot, relative string) bool {
	path, ok := resolveRepositoryPath(repoRoot, relative)
	if !ok {
		return false
	}
	_, err := os.Stat(path)
	return err == nil
}

func markdownSectionLines(content []byte, heading string) []string {
	var section []string
	inSection := false
	for _, line := range strings.Split(string(content), "\n") {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "## ") {
			if inSection {
				break
			}
			inSection = strings.TrimSpace(strings.TrimPrefix(trimmed, "## ")) == heading
			continue
		}
		if inSection {
			section = append(section, line)
		}
	}
	return section
}

func sectionLineContaining(content []byte, heading, needle string) int {
	lines := strings.Split(string(content), "\n")
	inSection := false
	sectionLine := 1
	for index, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "## ") {
			if inSection {
				break
			}
			if strings.TrimSpace(strings.TrimPrefix(trimmed, "## ")) == heading {
				inSection = true
				sectionLine = index + 1
			}
			continue
		}
		if inSection && strings.Contains(line, needle) {
			return index + 1
		}
	}
	return sectionLine
}

func newReferenceSet() referenceSet {
	return referenceSet{
		coverageFeature: {},
		coverageStory:   {},
	}
}

func addDeclaredReferences(references referenceSet, line string) {
	for _, number := range referenceNumbers(line, featureRefPattern) {
		references[coverageFeature][number] = true
	}
	for _, number := range referenceNumbers(line, storyRefPattern) {
		references[coverageStory][number] = true
	}
}

func referenceNumbers(text string, pattern *regexp.Regexp) []int {
	var numbers []int
	for _, match := range pattern.FindAllStringIndex(text, -1) {
		numbers = append(numbers, parseNumberSequence(text[match[1]:])...)
	}
	return numbers
}

func parseNumberSequence(text string) []int {
	cursor := skipSpaces(text, 0)
	first, next, ok := parsePositiveInt(text, cursor)
	if !ok {
		return nil
	}
	numbers := []int{first}
	cursor = next
	for {
		cursor = skipSpaces(text, cursor)
		if dashLength := rangeDashLength(text[cursor:]); dashLength > 0 {
			cursor = skipSpaces(text, cursor+dashLength)
			end, after, parsed := parsePositiveInt(text, cursor)
			if !parsed {
				return numbers
			}
			start := numbers[len(numbers)-1]
			if end >= start {
				for number := start + 1; number <= end; number++ {
					numbers = append(numbers, number)
				}
			}
			cursor = after
			continue
		}

		separatorLength := listSeparatorLength(text[cursor:])
		if separatorLength == 0 {
			return numbers
		}
		cursor = skipSpaces(text, cursor+separatorLength)
		if strings.HasPrefix(strings.ToLower(text[cursor:]), "and ") {
			cursor = skipSpaces(text, cursor+len("and"))
		}
		number, after, parsed := parsePositiveInt(text, cursor)
		if !parsed {
			return numbers
		}
		numbers = append(numbers, number)
		cursor = after
	}
}

func skipSpaces(text string, cursor int) int {
	for cursor < len(text) && (text[cursor] == ' ' || text[cursor] == '\t') {
		cursor++
	}
	return cursor
}

func parsePositiveInt(text string, cursor int) (int, int, bool) {
	start := cursor
	for cursor < len(text) && text[cursor] >= '0' && text[cursor] <= '9' {
		cursor++
	}
	if cursor == start {
		return 0, start, false
	}
	number, err := strconv.Atoi(text[start:cursor])
	return number, cursor, err == nil
}

func rangeDashLength(text string) int {
	for _, dash := range []string{"-", "–", "—"} {
		if strings.HasPrefix(text, dash) {
			return len(dash)
		}
	}
	return 0
}

func listSeparatorLength(text string) int {
	if strings.HasPrefix(text, ",") || strings.HasPrefix(text, "&") {
		return 1
	}
	lower := strings.ToLower(text)
	if strings.HasPrefix(lower, "and ") {
		return len("and")
	}
	return 0
}

func coverageUnitName(unit coverageUnit) string {
	return string(unit.Kind) + " " + strconv.Itoa(unit.Number)
}
