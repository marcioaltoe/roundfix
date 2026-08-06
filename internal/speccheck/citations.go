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
	// CodeCoverageUnmapped identifies a PRD unit absent from the TechSpec Coverage Map.
	CodeCoverageUnmapped = "SC-COVERAGE-UNMAPPED"
	// CodeCoverageUntasked identifies a PRD unit absent from every Task References section.
	CodeCoverageUntasked = "SC-COVERAGE-UNTASKED"
	// CodeReferenceUnresolved identifies a declared repository path that does not resolve.
	CodeReferenceUnresolved = "SC-REF-UNRESOLVED"
	// CodeLoopOrderDivergent identifies disagreeing declared Spec loop orders.
	CodeLoopOrderDivergent = "SC-LOOP-ORDER-DIVERGENT"
)

var (
	citationCoverageDetectorCodes = []string{
		CodeADRUnlisted,
		CodeADRRelated,
		CodeCoverageUnmapped,
		CodeCoverageUntasked,
		CodeReferenceUnresolved,
		CodeVerifyWorkIndependent,
		CodeRequirementContradictory,
		CodeRehearsalUndeclared,
	}
	adrFilenamePattern    = regexp.MustCompile(`^([0-9]{4})-.*\.md$`)
	adrCitationPattern    = regexp.MustCompile(`\bADR-([0-9]{4})\b`)
	legacyInactivePattern = regexp.MustCompile(`(?im)^\s*(?:\*\*)?status(?:\*\*)?:\s*(?:proposed|rejected|deprecated|superseded)\b`)
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
		name:       "repository guide",
		path:       loopOrderRepositoryGuidePath,
		marker:     "follow one order:",
		lineNeedle: "order: implement the graph",
	},
	{
		name:       "Baseline module asset",
		path:       loopOrderBaselineModulePath,
		marker:     "Follow one order per Spec:",
		lineNeedle: `"guidance": "Follow one order per Spec:`,
		module:     true,
	},
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
	DisplayPath   string
	Citations     map[string]bool
	CitationLines map[string]int
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

	units := parsePRDCoverageUnits(prdContent)
	if techSpecPresent {
		techSpecPath := filepath.Join(specDir, "_techspec.md")
		techSpecContent, err := os.ReadFile(techSpecPath)
		if err != nil {
			return fmt.Errorf("read Spec artifact %q: %w", techSpecPath, err)
		}
		detectCoverageMap(result, units, prdDisplayPath, techSpecContent, artifactDisplayPath(repoRoot, techSpecPath))
	} else {
		addSkip(result, CodeCoverageUnmapped, artifactDisplayPath(repoRoot, filepath.Join(specDir, "_techspec.md")))
	}

	graph, graphPresent, err := loadOptionalTaskGraph(specsRoot, slug, specDir)
	if err != nil {
		return err
	}
	if graphPresent {
		if err := detectTaskCoverageAndContextReferences(result, repoRoot, specsRoot, graph, units, prdDisplayPath); err != nil {
			return err
		}
	} else {
		manifestDisplayPath := artifactDisplayPath(repoRoot, filepath.Join(specDir, "_tasks.md"))
		addSkip(result, CodeCoverageUntasked, manifestDisplayPath)
		addSkip(result, CodeReferenceUnresolved, manifestDisplayPath)
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
		return nil
	}

	rowText, rowLine, rowPresent := projectConstraintEntry(prdContent, constraintActiveADR)
	if !rowPresent {
		missing := prdDisplayPath + " Active ADR obligations row"
		addSkip(result, CodeADRUnlisted, missing)
		addSkip(result, CodeADRRelated, missing)
		return nil
	}
	listed := citationNumbers([]byte(rowText))
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
			Fix:      "List ADR-" + number + " in the Active ADR obligations row in " + prdDisplayPath + " or remove the stale citation.",
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
		corpus[match[1]] = adrRecord{
			Number:        match[1],
			Title:         firstHeading(body),
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
		if legacyInactivePattern.Match(content) {
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

func citationNumbers(content []byte) map[string]bool {
	numbers := make(map[string]bool)
	for _, match := range adrCitationPattern.FindAllSubmatch(content, -1) {
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
		if finding, ok := WorkIndependentVerification(task); ok {
			finding.Where[0] = Location{
				Path: artifactDisplayPath(repoRoot, taskPath),
				Line: sectionLineContaining(content, "Verification", task.Verification[0]),
			}
			finding.Summary = finding.Where[0].Path + strings.TrimPrefix(finding.Summary, task.File)
			result.Findings = append(result.Findings, finding)
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
