package speccheck

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"roundfix/internal/spec"
)

const (
	// CodeRequirementContradictory identifies declared MUST and MUST NOT clauses over one named subject.
	CodeRequirementContradictory = "SC-REQUIREMENT-CONTRADICTORY"
	// CodeRehearsalUndeclared identifies a gate rehearsal with no complete case declarations.
	CodeRehearsalUndeclared = "SC-REHEARSAL-UNDECLARED"
)

// Stage names the authoring moment a caller is validating.
type Stage string

const (
	// StageAll runs the full detector sweep used by Check.
	StageAll Stage = ""
	// StagePRD runs detectors decidable once the PRD exists.
	StagePRD Stage = "prd"
	// StageTechSpec runs detectors decidable once the PRD and TechSpec exist.
	StageTechSpec Stage = "techspec"
	// StageTasks runs every detector because the complete Spec exists.
	StageTasks Stage = "tasks"
)

type stagedDetector struct {
	code  string
	stage Stage
}

var stagedDetectors = []stagedDetector{
	{code: CodeLoopOrderDivergent, stage: StagePRD},
	{code: CodeFindingLifecycle, stage: StagePRD},
	{code: CodeRollupMember, stage: StagePRD},
	{code: CodeArchiveLicense, stage: StagePRD},
	{code: CodeBacklogUnmoved, stage: StagePRD},
	{code: CodeConstraintMissing, stage: StagePRD},
	{code: CodeConstraintUnreasoned, stage: StagePRD},
	{code: CodeConstraintSource, stage: StagePRD},
	{code: CodeToolingUnauthorized, stage: StagePRD},
	{code: CodeToolingUnbounded, stage: StagePRD},
	{code: CodeCoverageUnmapped, stage: StageTechSpec},
	{code: CodeVocabularyUndocumented, stage: StageTechSpec},
	{code: CodeADRUnlisted, stage: StageTasks},
	{code: CodeADRRelated, stage: StageTasks},
	{code: CodeCoverageUntasked, stage: StageTasks},
	{code: CodeReferenceUnresolved, stage: StageTasks},
	{code: CodeVerifyWorkIndependent, stage: StageTasks},
	{code: CodeRequirementContradictory, stage: StageTasks},
	{code: CodeRehearsalUndeclared, stage: StageTasks},
}

// CheckStage runs the detectors whose inputs exist by stage. StageAll keeps
// the existing Check sweep unchanged, and StageTasks runs that same full sweep.
func CheckStage(specsRoot, repoRoot, slug string, stage Stage) (Result, error) {
	switch stage {
	case StageAll, StageTasks:
		return Check(specsRoot, repoRoot, slug)
	case StagePRD, StageTechSpec:
		return checkAuthoringStage(specsRoot, repoRoot, slug, stage)
	default:
		return Result{Slug: slug, Findings: []Finding{}, Skipped: []SkippedDetector{}},
			fmt.Errorf("unknown Spec authoring stage %q; accepted values: prd, techspec, tasks", stage)
	}
}

func checkAuthoringStage(specsRoot, repoRoot, slug string, stage Stage) (Result, error) {
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

	prdPath := filepath.Join(specDir, "_prd.md")
	prd, present, err := readConstraintArtifact(repoRoot, prdPath)
	if err != nil {
		return result, err
	}
	if !present {
		for _, code := range detectorCodes {
			addSkip(&result, code, artifactDisplayPath(repoRoot, prdPath))
		}
		if stage == StageTechSpec {
			addSkip(&result, CodeCoverageUnmapped, artifactDisplayPath(repoRoot, prdPath))
		}
		addStageSkips(&result, stage)
		return result, nil
	}

	artifacts := []constraintArtifact{prd}
	techSpecPath := filepath.Join(specDir, "_techspec.md")
	techSpecPresent := false
	if stage == StageTechSpec {
		techSpec, found, err := readConstraintArtifact(repoRoot, techSpecPath)
		if err != nil {
			return result, err
		}
		techSpecPresent = found
		if found {
			artifacts = append(artifacts, techSpec)
		} else {
			for _, code := range detectorCodes {
				addSkip(&result, code, artifactDisplayPath(repoRoot, techSpecPath))
			}
		}
		if err := detectVocabularyContract(&result, repoRoot, techSpecPath, found); err != nil {
			return result, err
		}
	}

	for artifactIndex := range artifacts {
		detectConstraintRows(&result, repoRoot, slug, artifacts, artifactIndex)
	}
	if stage == StageTechSpec {
		if err := detectTechSpecCoverage(&result, repoRoot, prdPath, techSpecPath, techSpecPresent); err != nil {
			return result, err
		}
	}
	addStageSkips(&result, stage)
	return result, nil
}

func detectTechSpecCoverage(result *Result, repoRoot, prdPath, techSpecPath string, techSpecPresent bool) error {
	prdContent, err := os.ReadFile(prdPath)
	if err != nil {
		return fmt.Errorf("read Spec artifact %q: %w", prdPath, err)
	}
	if !techSpecPresent {
		addSkip(result, CodeCoverageUnmapped, artifactDisplayPath(repoRoot, techSpecPath))
		return nil
	}
	techSpecContent, err := os.ReadFile(techSpecPath)
	if err != nil {
		return fmt.Errorf("read Spec artifact %q: %w", techSpecPath, err)
	}
	detectCoverageMap(
		result,
		parsePRDCoverageUnits(prdContent),
		artifactDisplayPath(repoRoot, prdPath),
		techSpecContent,
		artifactDisplayPath(repoRoot, techSpecPath),
	)
	return nil
}

func addStageSkips(result *Result, stage Stage) {
	for _, detector := range stagedDetectors {
		if stageIncludes(stage, detector.stage) {
			continue
		}
		addSkip(result, detector.code, "stage "+string(stage))
	}
}

func stageIncludes(scope, detectorStage Stage) bool {
	rank := func(stage Stage) int {
		switch stage {
		case StagePRD:
			return 1
		case StageTechSpec:
			return 2
		case StageTasks:
			return 3
		default:
			return 0
		}
	}
	return rank(detectorStage) <= rank(scope)
}

var (
	modalMentionPattern     = regexp.MustCompile(`(?i)\bMUST\s+and\s+MUST\s+NOT\b`)
	requirementModalPattern = regexp.MustCompile(`(?i)(?:^|[,;]\s+|\band\s+)MUST(\s+NOT)?\s+`)
	requirementWordPattern  = regexp.MustCompile(`[a-z0-9]+`)
	rehearsalWordPattern    = regexp.MustCompile(`(?i)\brehears[a-z]*\b`)
	proveWordPattern        = regexp.MustCompile(`(?i)\bprov(?:e|es|ed|ing)\b`)
	gateWordPattern         = regexp.MustCompile(`(?i)\bgates?\b`)
)

type declaredRequirementClause struct {
	Line      int
	Forbidden bool
	Words     []string
}

// ContradictoryRequirements reports one requirement forbidding a state that
// another requirement needs, decided from declared MUST/MUST NOT pairs over
// the same named subject.
func ContradictoryRequirements(task spec.Task) (Finding, bool) {
	clauses := declaredRequirementClauses(task.Requirements)
	for _, needed := range clauses {
		if needed.Forbidden {
			continue
		}
		for _, forbidden := range clauses {
			if !forbidden.Forbidden {
				continue
			}
			subject, ok := commonNamedSubject(needed.Words, forbidden.Words)
			if !ok {
				continue
			}
			return Finding{
				Code:     CodeRequirementContradictory,
				Severity: SeverityError,
				Summary:  task.File + " declares MUST and MUST NOT requirements over named subject " + subject,
				Where: []Location{
					{Path: task.File, Line: needed.Line},
					{Path: task.File, Line: forbidden.Line},
				},
				Fix: "Rewrite one requirement so the declared state of " + subject + " is consistent.",
			}, true
		}
	}
	return Finding{}, false
}

// UndeclaredRehearsal reports a Task whose title declares a gate rehearsal but
// whose Rehearsal Cases section has no complete case and observation entries.
func UndeclaredRehearsal(task spec.Task) (Finding, bool) {
	if !declaresGateRehearsal(task.Title) {
		return Finding{}, false
	}
	for _, declaration := range task.RehearsalCases {
		if !completeRehearsalCase(declaration.Text) {
			return undeclaredRehearsalFinding(task, declaration.Line), true
		}
	}
	if len(task.RehearsalCases) == 0 {
		return undeclaredRehearsalFinding(task, task.TitleLine), true
	}
	return Finding{}, false
}

func declaredRequirementClauses(requirements []spec.TaskDeclaration) []declaredRequirementClause {
	var clauses []declaredRequirementClause
	for _, requirement := range requirements {
		text := modalMentionPattern.ReplaceAllString(requirement.Text, "declared modal pair")
		matches := requirementModalPattern.FindAllStringSubmatchIndex(text, -1)
		for matchIndex, match := range matches {
			end := len(text)
			if matchIndex+1 < len(matches) {
				end = matches[matchIndex+1][0]
			}
			clauseText := strings.Trim(strings.TrimSpace(text[match[1]:end]), " .,:;—–-")
			if clauseText == "" {
				continue
			}
			clauses = append(clauses, declaredRequirementClause{
				Line:      requirement.Line,
				Forbidden: match[2] >= 0,
				Words:     significantRequirementWords(clauseText),
			})
		}
	}
	return clauses
}

func declaresGateRehearsal(title string) bool {
	if !gateWordPattern.MatchString(title) {
		return false
	}
	return rehearsalWordPattern.MatchString(title) || proveWordPattern.MatchString(title)
}

func significantRequirementWords(clause string) []string {
	var words []string
	for _, word := range requirementWordPattern.FindAllString(strings.ToLower(clause), -1) {
		word = requirementWordStem(word)
		if word == "" || requirementStopWords[word] {
			continue
		}
		words = append(words, word)
	}
	return words
}

func requirementWordStem(word string) string {
	switch {
	case strings.HasPrefix(word, "committ") || word == "commits":
		return "commit"
	case len(word) > 4 && strings.HasSuffix(word, "ies"):
		return strings.TrimSuffix(word, "ies") + "y"
	case len(word) > 4 && strings.HasSuffix(word, "s"):
		return strings.TrimSuffix(word, "s")
	default:
		return word
	}
}

func commonNamedSubject(needed, forbidden []string) (string, bool) {
	forbiddenWords := make(map[string]bool, len(forbidden))
	for _, word := range forbidden {
		forbiddenWords[word] = true
	}
	neededWords := make(map[string]bool, len(needed))
	for _, word := range needed {
		neededWords[word] = true
	}
	union := make(map[string]bool, len(neededWords)+len(forbiddenWords))
	for word := range neededWords {
		union[word] = true
	}
	for word := range forbiddenWords {
		union[word] = true
	}

	shared := 0
	subject := ""
	for _, word := range needed {
		if forbiddenWords[word] && subject == "" {
			subject = word
		}
	}
	for word := range neededWords {
		if forbiddenWords[word] {
			shared++
		}
	}
	if shared == 1 && subject == "commit" {
		return subject, true
	}
	if shared < 2 || shared*2 < len(union) {
		return "", false
	}
	return subject, true
}

func completeRehearsalCase(entry string) bool {
	lower := strings.ToLower(strings.TrimSpace(entry))
	if !strings.HasPrefix(lower, "case:") {
		return false
	}
	separator := strings.Index(lower, "; observation:")
	if separator < 0 {
		return false
	}
	caseText := strings.TrimSpace(entry[len("case:"):separator])
	observationText := strings.TrimSpace(entry[separator+len("; observation:"):])
	return caseText != "" && observationText != ""
}

func undeclaredRehearsalFinding(task spec.Task, line int) Finding {
	if line < 1 {
		line = 1
	}
	return Finding{
		Code:     CodeRehearsalUndeclared,
		Severity: SeverityError,
		Summary:  task.File + " declares a gate rehearsal but has no complete Rehearsal Cases declaration",
		Where:    []Location{{Path: task.File, Line: line}},
		Fix:      "Add a ## Rehearsal Cases section with one `- Case: <case>; Observation: <observation>` entry for each case.",
	}
}

var requirementStopWords = map[string]bool{
	"a": true, "an": true, "and": true, "another": true, "as": true,
	"at": true, "be": true, "by": true, "declared": true, "each": true,
	"for": true, "from": true, "in": true, "into": true, "is": true,
	"it": true, "its": true, "no": true, "of": true, "on": true,
	"named": true, "one": true, "or": true, "over": true, "same": true, "that": true,
	"the": true, "their": true, "them": true, "this": true, "to": true,
	"when": true, "whose": true, "with": true,
	"add": true, "change": true, "create": true, "decide": true,
	"declare": true, "define": true, "edit": true, "implement": true,
	"keep": true, "leave": true, "prove": true, "report": true,
	"requirement": true, "reuse": true, "run": true, "task": true,
	"update": true, "use": true,
}
