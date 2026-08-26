package speccheck

import (
	"bufio"
	"bytes"
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"gopkg.in/yaml.v3"

	"roundfix/internal/baseline"
	"roundfix/internal/spec"
	"roundfix/internal/suiteguardcontract"
	"roundfix/internal/worktree"
)

const (
	DetectorMechanicalAuthPaths       = "authorization bounded paths"
	DetectorMechanicalConsequentOrder = "consequent-fix commit order"
	DetectorMechanicalReportShape     = "QA Report structure"
	DetectorMechanicalEvidencePath    = "QA evidence paths"

	CodeMechanicalAuthPaths       = "QA-AUTH-PATHS"
	CodeMechanicalConsequentOrder = "QA-CONSEQUENT-ORDER"
	CodeMechanicalReportShape     = "QA-REPORT-SHAPE"
	CodeMechanicalEvidencePath    = "QA-EVIDENCE-PATH"

	blockedCauseLiteral = " — waits on "
)

// MechanicalRequest names the written declarations and observable repository
// facts available to one pre-QA mechanical pass. Empty or missing artifact
// inputs are recorded as MechanicalSkips.
type MechanicalRequest struct {
	RepoRoot          string
	AuthorizationPath string
	TaskCommits       []MechanicalTaskCommit
	ConsequentFixes   []ConsequentFixDeclaration
	ReportPath        string
	TaskRepairPaths   []string
	AssignedRepairs   []AssignedRepair
	Precondition      GatePreconditionResult
}

// AssignedRepair is one exact replacement the gate's Task requires. Path must
// also appear in TaskRepairPaths; Before and After make the operation and its
// verification deterministic instead of asking the mechanical stage to
// interpret prose.
type AssignedRepair struct {
	ID     string
	Path   string
	Before string
	After  string
}

// PerformedRepair records an assigned repair whose write was independently
// read back and verified. WriteMechanicalResult makes this audit trail visible
// in the seeded QA Report.
type PerformedRepair struct {
	ID   string
	Path string
}

// RepairFailure records assigned work the gate could not safely perform or
// verify. It is deliberately not a MechanicalFinding: findings remain
// observations for work the Task did not assign.
type RepairFailure struct {
	ID     string
	Path   string
	Detail string
}

// GatePreconditionCheck names the check a QA gate runs before it builds its
// matrix: the Spec's own strict Spec Consistency Check. One name covers every
// refusal that check produces, because the finding codes recorded beside it
// are what identify the detector that refused.
const GatePreconditionCheck = "spec check --strict"

// GatePreconditionResult separates contradictions that stop a QA gate from
// written inputs the gate itself is assigned to repair. Inputs remain visible
// so the gate can act on them; they are not silently discarded to get green.
type GatePreconditionResult struct {
	Findings []Finding
	Inputs   []Finding
	Blocking bool
}

// GatePrecondition classifies one ordinary Spec Consistency Check result for
// the Spec's own QA gate. Only an undocumented emitted token selected by that
// same Spec's complete Vocabulary Contract becomes gate input. All other
// findings keep their ordinary blocking behavior, and the source Result is not
// mutated.
func GatePrecondition(checked Result) GatePreconditionResult {
	precondition := GatePreconditionResult{
		Findings: []Finding{},
		Inputs:   []Finding{},
	}
	for _, finding := range checked.Findings {
		if finding.Code == CodeVocabularyUndocumented &&
			finding.declaredVocabularySpec != "" &&
			finding.declaredVocabularySpec == checked.Slug {
			precondition.Inputs = append(precondition.Inputs, finding)
			continue
		}
		precondition.Findings = append(precondition.Findings, finding)
	}
	precondition.Blocking = len(precondition.Findings) > 0
	return precondition
}

// PreconditionRefusal derives the audit record of a QA gate stop: the check
// that refused and, from the same strict Spec Consistency Check result the
// gate already classified, every code and sentence behind the refusal. The
// result of that check is read where it is produced rather than re-parsed from
// its rendered text, so a reworded line cannot change what the report records.
//
// A gate whose precondition did not block refuses nothing, which the second
// return separates from a refusal whose cause carries no name. That one is
// still recorded: a refusal the gate cannot write is the deadlock the refusal
// report exists to end.
func PreconditionRefusal(precondition GatePreconditionResult) (spec.PreconditionRefusal, bool) {
	if !precondition.Blocking {
		return spec.PreconditionRefusal{}, false
	}
	reasons := make([]string, 0, len(precondition.Findings))
	recorded := make(map[string]bool, len(precondition.Findings))
	for _, finding := range precondition.Findings {
		// Every distinct cause is kept, because the reason is the only record
		// of why the gate stopped; a repeat of one adds nothing to read.
		reason := preconditionRefusalReason(finding)
		if reason == "" || recorded[reason] {
			continue
		}
		recorded[reason] = true
		reasons = append(reasons, reason)
	}
	return spec.PreconditionRefusal{
		CheckName: GatePreconditionCheck,
		Reason:    strings.Join(reasons, "; "),
	}, true
}

// preconditionRefusalReason renders one refusing finding as its durable code
// followed by the sentence that explains it. The code leads because it is the
// name a later reader and detector share; the sentence beside it may be
// reworded, the code may not. A finding that names only one of the two is
// recorded by that one alone rather than by an empty label.
func preconditionRefusalReason(finding Finding) string {
	code := collapseRefusalText(finding.Code)
	summary := collapseRefusalText(finding.Summary)
	switch {
	case code == "":
		return summary
	case summary == "":
		return code
	default:
		return code + ": " + summary
	}
}

// collapseRefusalText folds one authored value onto a single line so a joined
// reason stays one frontmatter value in the report it is written into.
func collapseRefusalText(value string) string {
	return strings.Join(strings.Fields(value), " ")
}

// MechanicalTaskCommit associates the Daemon-owned Task commit with its Task
// declaration. TaskFile is the one always-allowed progress artifact.
type MechanicalTaskCommit struct {
	TaskID   string
	SHA      string
	TaskFile string
}

// ConsequentFixDeclaration is the written relationship between a corrective
// commit and the commit that caused it.
type ConsequentFixDeclaration struct {
	File        string
	Line        int
	RowHint     string
	CauseCommit string
	FixCommit   string
}

// MechanicalAuthorization reads the Tooling authority declaration from one
// PRD and returns the cited authorization plus its exact bounded paths. An
// absent declaration or artifact is represented by empty values so the
// mechanical detector can record the corresponding presence-aware skip.
func MechanicalAuthorization(repoRoot, prdPath string) (string, []string, error) {
	artifact, present, err := readConstraintArtifact(repoRoot, prdPath)
	if err != nil {
		return "", nil, err
	}
	if !present {
		return "", nil, nil
	}
	row, present := artifact.rows[strings.ToLower(constraintTooling)]
	if !present || strings.TrimSpace(row.AuthorizationPath) == "" {
		return "", nil, nil
	}
	authorizationPath := row.AuthorizationPath
	resolved, ok := resolveRepositoryPath(repoRoot, authorizationPath)
	if !ok {
		return authorizationPath, nil, nil
	}
	content, err := os.ReadFile(resolved)
	if errors.Is(err, os.ErrNotExist) {
		return authorizationPath, nil, nil
	}
	if err != nil {
		return "", nil, fmt.Errorf("read mechanical authorization %q: %w", resolved, err)
	}
	declared := parseMechanicalAuthorizationPaths(content)
	paths := make([]string, 0, len(declared))
	for path := range declared {
		paths = append(paths, path)
	}
	sort.Strings(paths)
	return authorizationPath, paths, nil
}

type mechanicalReportRow struct {
	id     string
	status string
	// provenance is the Results table's Provenance cell, empty when the table
	// declares no such column. It names where a row came from, which is what
	// separates a gate's precondition refusal from a measurement that blocked.
	provenance string
	evidence   string
	inputs     []EvidenceInput
	line       int
}

type mechanicalEvidenceRecord struct {
	head      string
	snapshots []EvidenceSnapshot
}

type mechanicalReport struct {
	path                    string
	rows                    []mechanicalReportRow
	rowsBlockedEnvironment  int
	rowsBlockedFinding      int
	rowsBlockedDeclared     int
	rowsBlockedPrecondition int
	countLines              map[string]int
	evidenceSnapshots       map[string]mechanicalEvidenceRecord
	evidenceSnapshotsErr    error
	parseError              error
}

var (
	markdownLinkTargetPattern = regexp.MustCompile(`\[[^\]]*\]\(([^)]+)\)`)
	evidencePathPattern       = regexp.MustCompile(`(?:^|[[:space:],;])((?:qa/)?evidence/[A-Za-z0-9._/-]+)`)
	sha256DigestPattern       = regexp.MustCompile(`^[0-9a-f]{64}$`)
	mechanicalRowHeading      = regexp.MustCompile(`^###\s+([A-Za-z0-9._-]+)(?:\s|$)`)
	carriedRowStatusPattern   = regexp.MustCompile(`^carried \(established by: ([^;]+); head: ([^)]+)\)$`)
)

// Carriable reports whether a passing row has only repository inputs, proven
// ancestry, no intersecting changed path, complete snapshots, and byte-identical
// evidence at the current head. It never computes a QA verdict.
func Carriable(
	prior ReportRow,
	head string,
	changed []string,
	established []EvidenceSnapshot,
	current []EvidenceSnapshot,
) bool {
	if prior.Status != "pass" || strings.TrimSpace(prior.EstablishedBy) == "" ||
		strings.TrimSpace(prior.EstablishedHead) == "" || strings.TrimSpace(head) == "" ||
		!prior.AncestryVerified || len(prior.Inputs) == 0 {
		return false
	}

	seenRefs := make(map[string]bool, len(prior.Inputs))
	for _, input := range prior.Inputs {
		ref := cleanMechanicalPath(input.Ref)
		if input.Kind != EvidenceRepositoryPath || ref == "" || ref != input.Ref || seenRefs[ref] {
			return false
		}
		seenRefs[ref] = true
		matcher := evidencePathMatcher(ref)
		for _, changedPath := range changed {
			clean := cleanMechanicalPath(changedPath)
			if clean == "" || clean != changedPath || matcher(clean) {
				return false
			}
		}
	}

	if len(established) != len(prior.Inputs) || len(current) != len(prior.Inputs) {
		return false
	}
	for index, input := range prior.Inputs {
		if !validEvidenceSnapshot(established[index], input.Ref) ||
			!validEvidenceSnapshot(current[index], input.Ref) ||
			!sameEvidenceFiles(established[index].Files, current[index].Files) {
			return false
		}
	}
	matchers := make([]func(string) bool, len(prior.Inputs))
	for index := range prior.Inputs {
		matchers[index] = evidencePathMatcher(prior.Inputs[index].Ref)
	}
	for _, evidencePath := range prior.EvidencePaths {
		clean := cleanMechanicalPath(evidencePath)
		if clean == "" || clean != evidencePath {
			return false
		}
		covered := false
		for index := range prior.Inputs {
			if matchers[index](clean) && evidenceSnapshotContains(established[index], clean) &&
				evidenceSnapshotContains(current[index], clean) {
				covered = true
				break
			}
		}
		if !covered {
			return false
		}
	}
	return true
}

func validEvidenceSnapshot(snapshot EvidenceSnapshot, ref string) bool {
	if snapshot.Ref != ref || len(snapshot.Files) == 0 {
		return false
	}
	matcher := evidencePathMatcher(ref)
	previous := ""
	for _, file := range snapshot.Files {
		clean := cleanMechanicalPath(file.Path)
		if clean == "" || clean != file.Path || clean <= previous ||
			!matcher(clean) || !sha256DigestPattern.MatchString(file.SHA256) {
			return false
		}
		previous = clean
	}
	return true
}

func sameEvidenceFiles(established, current []EvidenceFile) bool {
	if len(established) != len(current) {
		return false
	}
	for index := range established {
		if established[index] != current[index] {
			return false
		}
	}
	return true
}

func evidenceSnapshotContains(snapshot EvidenceSnapshot, path string) bool {
	for _, file := range snapshot.Files {
		if file.Path == path {
			return true
		}
	}
	return false
}

// evidencePathMatcher compiles one ref into a reusable matcher. A literal ref
// compares by equality; a glob ref compiles its pattern once.
func evidencePathMatcher(ref string) func(string) bool {
	if !strings.ContainsAny(ref, "*?") {
		return func(candidate string) bool { return ref == candidate }
	}
	compiled := regexp.MustCompile(evidenceGlobPattern(ref))
	return compiled.MatchString
}

func evidenceGlobPattern(ref string) string {
	runes := []rune(ref)
	var pattern strings.Builder
	pattern.WriteByte('^')
	for index := 0; index < len(runes); index++ {
		switch runes[index] {
		case '*':
			if index+1 < len(runes) && runes[index+1] == '*' {
				index++
				if index+1 < len(runes) && runes[index+1] == '/' {
					index++
					pattern.WriteString(`(?:.*/)?`)
				} else {
					pattern.WriteString(`.*`)
				}
			} else {
				pattern.WriteString(`[^/]*`)
			}
		case '?':
			pattern.WriteString(`[^/]`)
		default:
			pattern.WriteString(regexp.QuoteMeta(string(runes[index])))
		}
	}
	pattern.WriteByte('$')
	return pattern.String()
}

// RunMechanicalStage performs and verifies explicitly assigned repairs, then
// evaluates every detector and returns all findings in one pass. Repairs may
// write only exact TaskRepairPaths. The stage writes no report and settles no
// Task.
func RunMechanicalStage(ctx context.Context, request MechanicalRequest) (MechanicalResult, error) {
	result := MechanicalResult{
		Findings:       []MechanicalFinding{},
		Performed:      []PerformedRepair{},
		RepairFailures: []RepairFailure{},
		Carried:        []CarriedRow{},
		Blocked:        []BlockedRow{},
		Skips:          []MechanicalSkip{},
	}
	addGatePreconditionFindings(&result, request.Precondition)
	// The refusal travels with the result so the report writer never has to
	// reconstruct why the gate stopped from the findings it happens to carry.
	result.PreconditionRefusal, result.PreconditionRefused = PreconditionRefusal(request.Precondition)
	repoRoot := filepath.Clean(request.RepoRoot)
	if strings.TrimSpace(request.RepoRoot) == "" {
		return result, errors.New("run mechanical stage: repository root is empty")
	}
	if err := performAssignedRepairs(&result, repoRoot, request.TaskRepairPaths, request.AssignedRepairs); err != nil {
		return result, err
	}

	if err := detectMechanicalAuthPaths(ctx, &result, repoRoot, request); err != nil {
		return result, err
	}
	if err := detectMechanicalConsequentOrder(ctx, &result, repoRoot, request.ConsequentFixes); err != nil {
		return result, err
	}
	report, present, err := loadMechanicalReport(repoRoot, request.ReportPath)
	if err != nil {
		return result, err
	}
	if !present {
		missing := request.ReportPath
		if strings.TrimSpace(missing) == "" {
			missing = "QA Report"
		}
		addMechanicalSkip(&result, DetectorMechanicalReportShape, missing)
		addMechanicalSkip(&result, DetectorMechanicalEvidencePath, missing)
	} else {
		detectMechanicalReportShape(&result, report)
		detectMechanicalEvidencePaths(&result, repoRoot, report)
		if err := detectMechanicalEvidenceScratchState(ctx, &result, repoRoot, report); err != nil {
			return result, err
		}
		if report.evidenceSnapshotsErr != nil {
			addMechanicalSkip(&result, DetectorMechanicalReportShape, "evidence_snapshots")
		}
		currentHead, headErr := mechanicalHead(ctx, repoRoot)
		if headErr != nil {
			return result, headErr
		}
		result.Carried, err = resolveCarriedRows(ctx, repoRoot, report, currentHead)
		if err != nil {
			return result, err
		}
	}

	result.Blocking = len(result.Findings) > 0 || len(result.RepairFailures) > 0
	materializeBlockedRows(&result)
	return result, nil
}

func addGatePreconditionFindings(result *MechanicalResult, precondition GatePreconditionResult) {
	for _, finding := range precondition.Findings {
		mechanical := MechanicalFinding{
			Code:   finding.Code,
			Detail: finding.Summary,
			Fix:    finding.Fix,
		}
		if len(finding.Where) > 0 {
			mechanical.File = finding.Where[0].Path
			mechanical.Line = finding.Where[0].Line
		}
		addMechanicalFinding(result, mechanical)
	}
}

type preparedRepairPath struct {
	resolved string
	content  []byte
	mode     os.FileMode
}

func performAssignedRepairs(result *MechanicalResult, repoRoot string, taskPaths []string, repairs []AssignedRepair) error {
	if len(repairs) == 0 {
		return nil
	}

	allowed := make(map[string]bool, len(taskPaths))
	for _, taskPath := range taskPaths {
		clean := cleanMechanicalPath(taskPath)
		if clean != "" && clean == taskPath {
			allowed[clean] = true
		}
	}

	prepared := make(map[string]preparedRepairPath)
	applied := make([]PerformedRepair, 0, len(repairs))
	seenIDs := make(map[string]bool, len(repairs))
	for _, repair := range repairs {
		id := strings.TrimSpace(repair.ID)
		path := cleanMechanicalPath(repair.Path)
		failure := RepairFailure{ID: id, Path: repair.Path}
		switch {
		case id == "":
			failure.Detail = "assigned repair has no stable identifier and was not performed"
			addRepairFailure(result, failure)
			continue
		case seenIDs[id]:
			failure.Detail = fmt.Sprintf("assigned repair %s repeats an identifier and was not performed", id)
			addRepairFailure(result, failure)
			continue
		}
		seenIDs[id] = true

		if path == "" || path != repair.Path || !allowed[path] {
			failure.Detail = fmt.Sprintf("assigned repair %s names %q outside the Task-named repair paths", id, repair.Path)
			addRepairFailure(result, failure)
			continue
		}
		if repair.Before == "" || repair.Before == repair.After {
			failure.Detail = fmt.Sprintf("assigned repair %s has no deterministic before/after change and was not performed", id)
			addRepairFailure(result, failure)
			continue
		}

		state, loaded := prepared[path]
		if !loaded {
			resolved, ok := resolveRepositoryPath(repoRoot, path)
			if !ok {
				failure.Detail = fmt.Sprintf("assigned repair %s names invalid repository path %q", id, repair.Path)
				addRepairFailure(result, failure)
				continue
			}
			info, err := os.Lstat(resolved)
			if errors.Is(err, os.ErrNotExist) {
				failure.Detail = fmt.Sprintf("assigned repair %s path %s does not exist and was not performed", id, path)
				addRepairFailure(result, failure)
				continue
			}
			if err != nil {
				return fmt.Errorf("inspect assigned repair path %q: %w", resolved, err)
			}
			if !info.Mode().IsRegular() {
				failure.Detail = fmt.Sprintf("assigned repair %s path %s is not a regular file and was not performed", id, path)
				addRepairFailure(result, failure)
				continue
			}
			realRoot, err := filepath.EvalSymlinks(repoRoot)
			if err != nil {
				return fmt.Errorf("resolve assigned repair repository root %q: %w", repoRoot, err)
			}
			realPath, err := filepath.EvalSymlinks(resolved)
			if err != nil {
				return fmt.Errorf("resolve assigned repair path %q: %w", resolved, err)
			}
			if filepath.Clean(realPath) != filepath.Join(filepath.Clean(realRoot), filepath.FromSlash(path)) {
				failure.Detail = fmt.Sprintf("assigned repair %s path %s resolves through a symlink and was not performed", id, path)
				addRepairFailure(result, failure)
				continue
			}
			content, err := os.ReadFile(resolved)
			if err != nil {
				return fmt.Errorf("read assigned repair path %q: %w", resolved, err)
			}
			state = preparedRepairPath{resolved: resolved, content: content, mode: info.Mode().Perm()}
		}

		before := []byte(repair.Before)
		after := []byte(repair.After)
		switch occurrences := bytes.Count(state.content, before); occurrences {
		case 0:
			if bytes.Contains(state.content, after) {
				prepared[path] = state
				continue
			}
			failure.Detail = fmt.Sprintf("assigned repair %s was not performed: %s contains neither its before nor after text", id, path)
			addRepairFailure(result, failure)
		case 1:
			state.content = bytes.Replace(state.content, before, after, 1)
			prepared[path] = state
			applied = append(applied, PerformedRepair{ID: id, Path: path})
		default:
			failure.Detail = fmt.Sprintf("assigned repair %s matches %d places in %s and was not performed", id, occurrences, path)
			addRepairFailure(result, failure)
		}
	}
	if len(result.RepairFailures) > 0 {
		return nil
	}

	written := make(map[string]bool, len(applied))
	for _, repair := range applied {
		if written[repair.Path] {
			continue
		}
		written[repair.Path] = true
		state := prepared[repair.Path]
		if err := os.WriteFile(state.resolved, state.content, state.mode); err != nil {
			return fmt.Errorf("write assigned repair path %q: %w", state.resolved, err)
		}
		verified, err := os.ReadFile(state.resolved)
		if err != nil {
			return fmt.Errorf("verify assigned repair path %q: %w", state.resolved, err)
		}
		if !bytes.Equal(verified, state.content) {
			addRepairFailure(result, RepairFailure{
				ID: repair.ID, Path: repair.Path,
				Detail: fmt.Sprintf("assigned repair %s write to %s could not be verified", repair.ID, repair.Path),
			})
			return nil
		}
	}
	result.Performed = append(result.Performed, applied...)
	return nil
}

func addRepairFailure(result *MechanicalResult, failure RepairFailure) {
	result.RepairFailures = append(result.RepairFailures, failure)
}

func detectMechanicalAuthPaths(ctx context.Context, result *MechanicalResult, repoRoot string, request MechanicalRequest) error {
	missing := request.AuthorizationPath
	if strings.TrimSpace(missing) == "" {
		missing = "tooling authorization"
	}
	path, ok := resolveRepositoryPath(repoRoot, request.AuthorizationPath)
	if !ok {
		addMechanicalSkip(result, DetectorMechanicalAuthPaths, missing)
		return nil
	}
	content, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		addMechanicalSkip(result, DetectorMechanicalAuthPaths, missing)
		return nil
	}
	if err != nil {
		return fmt.Errorf("read mechanical authorization %q: %w", path, err)
	}
	if len(request.TaskCommits) == 0 {
		addMechanicalSkip(result, DetectorMechanicalAuthPaths, "Task commits")
		return nil
	}

	bounded := parseMechanicalAuthorizationPaths(content)
	if len(bounded) == 0 {
		addMechanicalFinding(result, MechanicalFinding{
			Code: CodeMechanicalAuthPaths, File: request.AuthorizationPath, Line: 1,
			Detail: "authorization declares no exact bounded files",
			Fix:    "Declare every authorized repository-relative path in the authorization artifact.",
		})
		return nil
	}
	regenerated, err := parseMechanicalRegenerationOutputs(repoRoot, content)
	if err != nil {
		return fmt.Errorf("resolve sanctioned regeneration outputs from %q: %w", request.AuthorizationPath, err)
	}
	authorizationPath := cleanMechanicalPath(request.AuthorizationPath)

	for _, taskCommit := range request.TaskCommits {
		if strings.TrimSpace(taskCommit.SHA) == "" {
			addMechanicalSkip(result, DetectorMechanicalAuthPaths, "Task commit for "+taskCommit.TaskID)
			continue
		}
		exists, err := mechanicalCommitExists(ctx, repoRoot, taskCommit.SHA)
		if err != nil {
			return err
		}
		if !exists {
			addMechanicalSkip(result, DetectorMechanicalAuthPaths, "Git commit "+taskCommit.SHA)
			continue
		}
		changed, err := mechanicalChangedPaths(ctx, repoRoot, taskCommit.SHA)
		if err != nil {
			return err
		}
		taskFile := cleanMechanicalPath(taskCommit.TaskFile)
		for _, changedPath := range changed {
			if changedPath == taskFile {
				continue
			}
			if authorizationPath != "" && changedPath == authorizationPath {
				addMechanicalFinding(result, MechanicalFinding{
					Code: CodeMechanicalAuthPaths, File: request.AuthorizationPath, Line: 1,
					Detail: fmt.Sprintf(
						"Task %s commit %s changes authorization grant %s in the commit that consumes it",
						taskCommit.TaskID, taskCommit.SHA, changedPath,
					),
					Fix:     "Land the authorization record in its own commit before the Task commit that consumes it.",
					RowHint: taskCommit.TaskID,
				})
				continue
			}
			if !GovernedPath(changedPath) || bounded[changedPath] || regenerated[changedPath] {
				continue
			}
			addMechanicalFinding(result, MechanicalFinding{
				Code: CodeMechanicalAuthPaths, File: request.AuthorizationPath, Line: 1,
				Detail: fmt.Sprintf(
					"Task %s commit %s changes %s outside authorization grant %s's exact bounded files",
					taskCommit.TaskID, taskCommit.SHA, changedPath, request.AuthorizationPath,
				),
				Fix:     "Move the path into an expressly authorized Task or narrow the commit to the written authorization and assigned Task file.",
				RowHint: taskCommit.TaskID,
			})
		}
	}
	return nil
}

func parseMechanicalAuthorizationPaths(content []byte) map[string]bool {
	paths := make(map[string]bool)
	text := strings.ReplaceAll(string(content), "\r\n", "\n")
	if strings.HasPrefix(text, "---\n") {
		if body, _, found := strings.Cut(text[len("---\n"):], "\n---"); found {
			var frontmatter struct {
				Paths []string `yaml:"paths"`
			}
			if yaml.Unmarshal([]byte(body), &frontmatter) == nil {
				for _, declared := range frontmatter.Paths {
					if clean := cleanMechanicalPath(declared); clean != "" {
						paths[clean] = true
					}
				}
			}
		}
	}

	inBoundedFiles := false
	for _, line := range strings.Split(text, "\n") {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "## ") {
			inBoundedFiles = strings.EqualFold(strings.TrimSpace(strings.TrimPrefix(trimmed, "## ")), "Bounded files")
			continue
		}
		if !inBoundedFiles {
			continue
		}
		for _, match := range backtickPattern.FindAllStringSubmatch(line, -1) {
			if clean := cleanMechanicalPath(match[1]); clean != "" {
				paths[clean] = true
			}
		}
	}
	return paths
}

func parseMechanicalRegenerationOutputs(repoRoot string, content []byte) (map[string]bool, error) {
	outputs := make(map[string]bool)
	for _, declaration := range suiteguardcontract.ParseSanctionedRegenerations(content) {
		for _, output := range declaration.Outputs {
			outputs[output] = true
		}
		resolved, err := baseline.OutputsFor(repoRoot, declaration.Command)
		if err != nil {
			return nil, fmt.Errorf("resolve outputs for command %q: %w", declaration.Command, err)
		}
		for _, output := range resolved {
			outputs[output] = true
		}
	}
	return outputs, nil
}

func cleanMechanicalPath(value string) string {
	value = strings.TrimSpace(value)
	if value == "" || filepath.IsAbs(value) || strings.Contains(value, `\`) {
		return ""
	}
	clean := filepath.ToSlash(filepath.Clean(filepath.FromSlash(value)))
	if clean == "." || clean == ".." || strings.HasPrefix(clean, "../") {
		return ""
	}
	return clean
}

func detectMechanicalConsequentOrder(ctx context.Context, result *MechanicalResult, repoRoot string, declarations []ConsequentFixDeclaration) error {
	if len(declarations) == 0 {
		addMechanicalSkip(result, DetectorMechanicalConsequentOrder, "consequent-fix declaration")
		return nil
	}
	for _, declaration := range declarations {
		file := declaration.File
		if file == "" {
			file = "Task Graph"
		}
		line := declaration.Line
		if line < 1 {
			line = 1
		}
		if declaration.CauseCommit == "" || declaration.FixCommit == "" {
			addMechanicalSkip(result, DetectorMechanicalConsequentOrder, "consequent-fix commit reference in "+file)
			continue
		}
		if declaration.CauseCommit == declaration.FixCommit {
			addMechanicalFinding(result, MechanicalFinding{
				Code: CodeMechanicalConsequentOrder, File: file, Line: line,
				Detail:  "consequent fix is folded into the commit that caused it: " + declaration.CauseCommit,
				Fix:     "Land the consequent fix in its own commit after the authorized cause commit.",
				RowHint: declaration.RowHint,
			})
			continue
		}
		missingCommit := ""
		for _, sha := range []string{declaration.CauseCommit, declaration.FixCommit} {
			exists, err := mechanicalCommitExists(ctx, repoRoot, sha)
			if err != nil {
				return err
			}
			if !exists {
				missingCommit = sha
				break
			}
		}
		if missingCommit != "" {
			addMechanicalSkip(result, DetectorMechanicalConsequentOrder, "Git commit "+missingCommit)
			continue
		}
		ordered, err := mechanicalIsAncestor(ctx, repoRoot, declaration.CauseCommit, declaration.FixCommit)
		if err != nil {
			return err
		}
		if !ordered {
			addMechanicalFinding(result, MechanicalFinding{
				Code: CodeMechanicalConsequentOrder, File: file, Line: line,
				Detail:  fmt.Sprintf("consequent fix commit %s is not ordered after cause commit %s", declaration.FixCommit, declaration.CauseCommit),
				Fix:     "Order the consequent fix after the commit that caused it, preserving chronological ancestry.",
				RowHint: declaration.RowHint,
			})
		}
	}
	return nil
}

func loadMechanicalReport(repoRoot, reportPath string) (mechanicalReport, bool, error) {
	resolved, ok := resolveRepositoryPath(repoRoot, reportPath)
	if !ok {
		return mechanicalReport{}, false, nil
	}
	content, err := os.ReadFile(resolved)
	if errors.Is(err, os.ErrNotExist) {
		return mechanicalReport{}, false, nil
	}
	if err != nil {
		return mechanicalReport{}, false, fmt.Errorf("read mechanical QA Report %q: %w", resolved, err)
	}
	report := parseMechanicalReport(reportPath, content)
	return report, true, nil
}

func parseMechanicalReport(path string, content []byte) mechanicalReport {
	report := mechanicalReport{
		path:              path,
		countLines:        make(map[string]int),
		evidenceSnapshots: make(map[string]mechanicalEvidenceRecord),
	}
	text := strings.ReplaceAll(string(content), "\r\n", "\n")
	if !strings.HasPrefix(text, "---\n") {
		report.parseError = errors.New("missing YAML frontmatter")
	} else if body, _, found := strings.Cut(text[len("---\n"):], "\n---"); !found {
		report.parseError = errors.New("unterminated YAML frontmatter")
	} else {
		var document yaml.Node
		if err := yaml.Unmarshal([]byte(body), &document); err != nil {
			report.parseError = fmt.Errorf("parse YAML frontmatter: %w", err)
		} else {
			report.evidenceSnapshots, report.evidenceSnapshotsErr = mechanicalEvidenceSnapshots(document)
			counts, lines, err := mechanicalBlockedCounts(document)
			if err != nil {
				report.parseError = err
			} else {
				report.rowsBlockedEnvironment = counts["rows_blocked_environment"]
				report.rowsBlockedFinding = counts["rows_blocked_finding"]
				report.rowsBlockedDeclared = counts["rows_blocked_declared"]
				report.rowsBlockedPrecondition = counts["rows_blocked_precondition"]
				report.countLines = lines
			}
		}
	}

	allLines := strings.Split(text, "\n")
	inResults := false
	statusColumn, evidenceColumn, provenanceColumn := -1, -1, -1
	for index, line := range allLines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "## ") {
			if inResults {
				break
			}
			inResults = strings.EqualFold(strings.TrimSpace(strings.TrimPrefix(trimmed, "## ")), "Results")
			continue
		}
		if !inResults {
			continue
		}
		cells := markdownCells(line)
		if len(cells) == 0 {
			continue
		}
		if statusColumn < 0 {
			for cellIndex, cell := range cells {
				switch strings.ToLower(strings.TrimSpace(cell)) {
				case "status":
					statusColumn = cellIndex
				case "evidence":
					evidenceColumn = cellIndex
				case "provenance":
					provenanceColumn = cellIndex
				}
			}
			continue
		}
		if markdownSeparator(cells) || statusColumn >= len(cells) {
			continue
		}
		row := mechanicalReportRow{id: strings.TrimSpace(cells[0]), status: strings.TrimSpace(cells[statusColumn]), line: index + 1}
		if evidenceColumn >= 0 && evidenceColumn < len(cells) {
			row.evidence = strings.TrimSpace(cells[evidenceColumn])
		}
		if provenanceColumn >= 0 && provenanceColumn < len(cells) {
			row.provenance = strings.TrimSpace(cells[provenanceColumn])
		}
		report.rows = append(report.rows, row)
	}
	parseMechanicalRowInputs(allLines, report.rows)
	return report
}

func mechanicalEvidenceSnapshots(document yaml.Node) (map[string]mechanicalEvidenceRecord, error) {
	records := make(map[string]mechanicalEvidenceRecord)
	if len(document.Content) == 0 || document.Content[0].Kind != yaml.MappingNode {
		return records, nil
	}
	mapping := document.Content[0]
	for index := 0; index+1 < len(mapping.Content); index += 2 {
		if mapping.Content[index].Value != "evidence_snapshots" {
			continue
		}
		var raw map[string]struct {
			Head   string `yaml:"head"`
			Inputs []struct {
				Ref   string `yaml:"ref"`
				Files []struct {
					Path   string `yaml:"path"`
					SHA256 string `yaml:"sha256"`
				} `yaml:"files"`
			} `yaml:"inputs"`
		}
		if err := mapping.Content[index+1].Decode(&raw); err != nil {
			return nil, fmt.Errorf("decode evidence_snapshots: %w", err)
		}
		for rowID, snapshot := range raw {
			record := mechanicalEvidenceRecord{head: strings.TrimSpace(snapshot.Head)}
			for _, input := range snapshot.Inputs {
				files := make([]EvidenceFile, 0, len(input.Files))
				for _, file := range input.Files {
					files = append(files, EvidenceFile{Path: strings.TrimSpace(file.Path), SHA256: strings.TrimSpace(file.SHA256)})
				}
				record.snapshots = append(record.snapshots, EvidenceSnapshot{Ref: strings.TrimSpace(input.Ref), Files: files})
			}
			records[strings.TrimSpace(rowID)] = record
		}
		return records, nil
	}
	return records, nil
}

func parseMechanicalRowInputs(lines []string, rows []mechanicalReportRow) {
	rowIndexes := make(map[string]int, len(rows))
	for index := range rows {
		rowIndexes[rows[index].id] = index
	}
	currentRow := -1
	seen := make(map[int]bool)
	invalid := make(map[int]bool)
	for index := 0; index < len(lines); index++ {
		trimmed := strings.TrimSpace(lines[index])
		if match := mechanicalRowHeading.FindStringSubmatch(trimmed); len(match) == 2 {
			currentRow = -1
			if rowIndex, ok := rowIndexes[match[1]]; ok {
				currentRow = rowIndex
			}
			continue
		}
		if currentRow < 0 || trimmed != "```yaml" {
			continue
		}
		start := index + 1
		for index++; index < len(lines) && strings.TrimSpace(lines[index]) != "```"; index++ {
		}
		if index >= len(lines) {
			invalid[currentRow] = true
			break
		}
		block := strings.Join(lines[start:index], "\n")
		var declaration struct {
			Inputs *[]struct {
				Kind EvidenceInputKind `yaml:"kind"`
				Ref  string            `yaml:"ref"`
			} `yaml:"inputs"`
		}
		if err := yaml.Unmarshal([]byte(block), &declaration); err != nil {
			if strings.HasPrefix(strings.TrimSpace(block), "inputs:") {
				invalid[currentRow] = true
			}
			continue
		}
		if declaration.Inputs == nil {
			continue
		}
		if seen[currentRow] {
			invalid[currentRow] = true
			continue
		}
		seen[currentRow] = true
		inputs := make([]EvidenceInput, 0, len(*declaration.Inputs))
		for _, input := range *declaration.Inputs {
			inputs = append(inputs, EvidenceInput{Kind: input.Kind, Ref: strings.TrimSpace(input.Ref)})
		}
		rows[currentRow].inputs = inputs
	}
	for rowIndex := range invalid {
		rows[rowIndex].inputs = nil
	}
}

func resolveCarriedRows(ctx context.Context, repoRoot string, priorReport mechanicalReport, currentHead string) ([]CarriedRow, error) {
	carried := make([]CarriedRow, 0)
	for _, priorRow := range priorReport.rows {
		establishingReport := priorReport
		establishingRow := priorRow
		establishedBy := priorReport.path
		requiredHead := ""

		if reportPath, head, ok := carriedRowCitation(priorRow.status); ok {
			if reportPath == priorReport.path {
				continue
			}
			loaded, present, err := loadMechanicalReport(repoRoot, reportPath)
			if err != nil {
				return nil, err
			}
			if !present {
				continue
			}
			row, found := mechanicalReportRowByID(loaded.rows, priorRow.id)
			if !found {
				continue
			}
			establishingReport = loaded
			establishingRow = row
			establishedBy = reportPath
			requiredHead = head
		}

		if establishingRow.status != "pass" || len(establishingRow.inputs) == 0 {
			continue
		}
		if !mechanicalSnapshotKeysNamePassingRows(establishingReport) {
			continue
		}
		record, ok := establishingReport.evidenceSnapshots[priorRow.id]
		if !ok || strings.TrimSpace(record.head) == "" || requiredHead != "" && requiredHead != record.head {
			continue
		}
		exists, err := mechanicalCommitExists(ctx, repoRoot, record.head)
		if err != nil {
			return nil, err
		}
		if !exists {
			continue
		}
		ancestor, err := mechanicalIsAncestor(ctx, repoRoot, record.head, currentHead)
		if err != nil {
			return nil, err
		}
		if !ancestor {
			continue
		}
		changed, err := worktree.PriorChangedFiles(ctx, repoRoot, record.head)
		if err != nil {
			return nil, fmt.Errorf("resolve carry-forward changed paths from %q: %w", record.head, err)
		}
		currentSnapshots, resolved, err := buildEvidenceSnapshots(ctx, repoRoot, currentHead, establishingRow.inputs)
		if err != nil {
			return nil, err
		}
		if !resolved {
			continue
		}
		row := ReportRow{
			ID:               priorRow.id,
			Status:           establishingRow.status,
			EstablishedBy:    establishedBy,
			EstablishedHead:  record.head,
			AncestryVerified: true,
			Inputs:           append([]EvidenceInput(nil), establishingRow.inputs...),
		}
		evidencePaths, resolved := mechanicalEvidenceRepositoryPaths(repoRoot, establishingReport.path, establishingRow.evidence)
		if !resolved {
			continue
		}
		row.EvidencePaths = evidencePaths
		if !Carriable(row, currentHead, changed, record.snapshots, currentSnapshots) {
			continue
		}
		carried = append(carried, CarriedRow{
			ID:              priorRow.id,
			EstablishedBy:   establishedBy,
			EstablishedHead: record.head,
			Inputs:          append([]EvidenceInput(nil), establishingRow.inputs...),
		})
	}
	return carried, nil
}

func mechanicalEvidenceRepositoryPaths(repoRoot, reportPath, evidence string) ([]string, bool) {
	reportDirectory := filepath.Dir(filepath.Join(repoRoot, filepath.FromSlash(reportPath)))
	seen := make(map[string]bool)
	var paths []string
	for _, target := range mechanicalEvidenceTargets(evidence) {
		target = strings.TrimSpace(strings.Split(target, "#")[0])
		if target == "" || strings.Contains(target, "://") || strings.HasPrefix(target, "mailto:") {
			continue
		}
		if filepath.IsAbs(target) || strings.Contains(target, `\`) {
			return nil, false
		}
		base := reportDirectory
		switch {
		case strings.HasPrefix(target, "qa/"):
			base = filepath.Dir(reportDirectory)
		case strings.HasPrefix(target, "docs/"):
			base = repoRoot
		}
		resolved := filepath.Clean(filepath.Join(base, filepath.FromSlash(target)))
		relative, err := filepath.Rel(repoRoot, resolved)
		if err != nil {
			return nil, false
		}
		clean := cleanMechanicalPath(filepath.ToSlash(relative))
		if clean == "" || seen[clean] {
			if clean == "" {
				return nil, false
			}
			continue
		}
		seen[clean] = true
		paths = append(paths, clean)
	}
	sort.Strings(paths)
	return paths, true
}

func mechanicalSnapshotKeysNamePassingRows(report mechanicalReport) bool {
	for rowID := range report.evidenceSnapshots {
		row, found := mechanicalReportRowByID(report.rows, rowID)
		if !found || row.status != "pass" {
			return false
		}
	}
	return true
}

func carriedRowCitation(status string) (string, string, bool) {
	match := carriedRowStatusPattern.FindStringSubmatch(strings.TrimSpace(status))
	if len(match) != 3 {
		return "", "", false
	}
	reportPath := cleanMechanicalPath(match[1])
	head := strings.TrimSpace(match[2])
	if reportPath == "" || reportPath != strings.TrimSpace(match[1]) || head == "" {
		return "", "", false
	}
	return reportPath, head, true
}

func mechanicalReportRowByID(rows []mechanicalReportRow, id string) (mechanicalReportRow, bool) {
	for _, row := range rows {
		if row.id == id {
			return row, true
		}
	}
	return mechanicalReportRow{}, false
}

func mechanicalHead(ctx context.Context, repoRoot string) (string, error) {
	command := exec.CommandContext(ctx, "git", "-C", repoRoot, "rev-parse", "--verify", "HEAD^{commit}")
	output, err := command.Output()
	if err != nil {
		return "", fmt.Errorf("resolve mechanical-stage Git head: %w", err)
	}
	head := strings.TrimSpace(string(output))
	if head == "" {
		return "", errors.New("resolve mechanical-stage Git head: empty commit")
	}
	return head, nil
}

func buildEvidenceSnapshots(ctx context.Context, repoRoot, head string, inputs []EvidenceInput) ([]EvidenceSnapshot, bool, error) {
	paths, err := mechanicalBlobPaths(ctx, repoRoot, head)
	if err != nil {
		return nil, false, err
	}
	snapshots := make([]EvidenceSnapshot, 0, len(inputs))
	seenRefs := make(map[string]bool, len(inputs))
	ordered := make([]mechanicalSnapshotInput, 0, len(inputs))
	for _, input := range inputs {
		ref := cleanMechanicalPath(input.Ref)
		if input.Kind != EvidenceRepositoryPath || ref == "" || ref != input.Ref || seenRefs[ref] {
			return nil, false, nil
		}
		seenRefs[ref] = true
		matcher := evidencePathMatcher(ref)
		var matches []string
		for _, candidate := range paths {
			if matcher(candidate) {
				matches = append(matches, candidate)
			}
		}
		if len(matches) == 0 {
			return nil, false, nil
		}
		ordered = append(ordered, mechanicalSnapshotInput{ref: ref, matches: matches})
	}

	digests, err := mechanicalBlobDigests(ctx, repoRoot, head, ordered)
	if err != nil {
		if ctxErr := ctx.Err(); ctxErr != nil {
			return nil, false, ctxErr
		}
		return nil, false, err
	}
	for _, input := range ordered {
		files := make([]EvidenceFile, 0, len(input.matches))
		for _, path := range input.matches {
			digest := digests[path]
			if digest == "" {
				return nil, false, fmt.Errorf("read blob at Git head %q for path %q: no content", head, path)
			}
			files = append(files, EvidenceFile{Path: path, SHA256: digest})
		}
		snapshots = append(snapshots, EvidenceSnapshot{Ref: input.ref, Files: files})
	}
	return snapshots, true, nil
}

// mechanicalSnapshotInput groups one evidence ref with the tracked paths it
// matched, so blob contents can be read in one batched process per head.
type mechanicalSnapshotInput struct {
	ref     string
	matches []string
}

// mechanicalBlobDigests reads the blob content for every matched path in one
// git cat-file --batch process, writing each object spec to stdin and parsing
// the length-prefixed records from stdout in match order. It returns a SHA256
// digest per repository-relative path.
func mechanicalBlobDigests(ctx context.Context, repoRoot, head string, inputs []mechanicalSnapshotInput) (map[string]string, error) {
	ordered := make([]string, 0)
	saw := make(map[string]bool)
	for _, input := range inputs {
		for _, path := range input.matches {
			if !saw[path] {
				saw[path] = true
				ordered = append(ordered, path)
			}
		}
	}

	command := exec.CommandContext(ctx, "git", "-C", repoRoot, "cat-file", "--batch")
	stdin, err := command.StdinPipe()
	if err != nil {
		return nil, fmt.Errorf("open cat-file stdin: %w", err)
	}
	stdout, err := command.StdoutPipe()
	if err != nil {
		return nil, fmt.Errorf("open cat-file stdout: %w", err)
	}
	if err := command.Start(); err != nil {
		return nil, fmt.Errorf("start git cat-file --batch: %w", err)
	}
	writeErr := make(chan error, 1)
	go func() {
		writer := bufio.NewWriter(stdin)
		for _, path := range ordered {
			if _, err := writer.WriteString(head + ":" + path + "\n"); err != nil {
				writeErr <- err
				return
			}
		}
		if err := writer.Flush(); err != nil {
			writeErr <- err
			return
		}
		if err := stdin.Close(); err != nil {
			writeErr <- err
			return
		}
		writeErr <- nil
	}()

	digests := make(map[string]string, len(ordered))
	reader := bufio.NewReader(stdout)
	var readErr error
	for _, path := range ordered {
		digest, err := readCatFileRecord(reader)
		if err != nil {
			readErr = fmt.Errorf("read blob at Git head %q for path %q: %w", head, path, err)
			break
		}
		digests[path] = digest
	}
	if readErr != nil {
		// Reap the child and stop feeding it before surfacing the failure, so an
		// error path never leaks the batched process or its writer goroutine.
		if command.Process != nil {
			_ = command.Process.Kill()
		}
		_ = command.Wait()
		return nil, readErr
	}
	if err := command.Wait(); err != nil {
		if werr := <-writeErr; werr != nil {
			return nil, fmt.Errorf("write cat-file specs: %w", werr)
		}
		return nil, fmt.Errorf("git cat-file --batch: %w", err)
	}
	if werr := <-writeErr; werr != nil {
		return nil, fmt.Errorf("write cat-file specs: %w", werr)
	}
	return digests, nil
}

func readCatFileRecord(reader *bufio.Reader) (string, error) {
	header, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	header = strings.TrimSuffix(header, "\n")
	fields := strings.Split(header, " ")
	if len(fields) != 3 {
		return "", fmt.Errorf("unexpected record header %q", header)
	}
	if fields[1] != "blob" {
		return "", fmt.Errorf("record %q is not a blob", header)
	}
	size, err := strconv.Atoi(fields[2])
	if err != nil || size < 0 {
		return "", fmt.Errorf("invalid blob size %q", fields[2])
	}
	content := make([]byte, size)
	if _, err := io.ReadFull(reader, content); err != nil {
		return "", err
	}
	if _, err := reader.ReadByte(); err != nil {
		return "", err
	}
	digest := sha256.Sum256(content)
	return fmt.Sprintf("%x", digest), nil
}

func mechanicalBlobPaths(ctx context.Context, repoRoot, head string) ([]string, error) {
	command := exec.CommandContext(ctx, "git", "-C", repoRoot, "ls-tree", "-r", "-z", "--full-tree", head)
	output, err := command.Output()
	if err != nil {
		return nil, fmt.Errorf("list tracked blobs at Git head %q: %w", head, err)
	}
	var paths []string
	for _, record := range bytes.Split(output, []byte{0}) {
		parts := bytes.SplitN(record, []byte{'\t'}, 2)
		if len(parts) != 2 || !bytes.Contains(parts[0], []byte(" blob ")) {
			continue
		}
		path := string(parts[1])
		clean := cleanMechanicalPath(path)
		if clean == "" || clean != path {
			return nil, fmt.Errorf("list tracked blobs at Git head %q: invalid path %q", head, path)
		}
		paths = append(paths, clean)
	}
	sort.Strings(paths)
	return paths, nil
}

func mechanicalBlockedCounts(document yaml.Node) (map[string]int, map[string]int, error) {
	counts := map[string]int{
		"rows_blocked_environment":  0,
		"rows_blocked_finding":      0,
		"rows_blocked_declared":     0,
		"rows_blocked_precondition": 0,
	}
	lines := make(map[string]int)
	if len(document.Content) == 0 || document.Content[0].Kind != yaml.MappingNode {
		return nil, nil, errors.New("YAML frontmatter must be a mapping")
	}
	mapping := document.Content[0]
	for index := 0; index+1 < len(mapping.Content); index += 2 {
		key, value := mapping.Content[index], mapping.Content[index+1]
		if _, tracked := counts[key.Value]; !tracked {
			continue
		}
		if value.Kind != yaml.ScalarNode || value.Tag != "!!int" {
			return nil, nil, fmt.Errorf("%s must be a non-negative integer", key.Value)
		}
		count, err := strconv.Atoi(value.Value)
		if err != nil || count < 0 {
			return nil, nil, fmt.Errorf("%s must be a non-negative integer", key.Value)
		}
		counts[key.Value] = count
		lines[key.Value] = value.Line + 1
	}
	return counts, lines, nil
}

func detectMechanicalReportShape(result *MechanicalResult, report mechanicalReport) {
	if report.parseError != nil {
		addMechanicalFinding(result, MechanicalFinding{
			Code: CodeMechanicalReportShape, File: report.path, Line: 1,
			Detail: "report structure cannot be parsed: " + report.parseError.Error(),
			Fix:    "Repair the QA Report frontmatter and keep all three typed blocked-cause counts as non-negative integers.",
		})
	}
	actual := map[string]int{
		"rows_blocked_environment":  0,
		"rows_blocked_finding":      0,
		"rows_blocked_declared":     0,
		"rows_blocked_precondition": 0,
	}
	unparsed := map[string]int{
		"rows_blocked_environment": 0,
		"rows_blocked_finding":     0,
		"rows_blocked_declared":    0,
	}
	if len(report.rows) == 0 {
		addMechanicalFinding(result, MechanicalFinding{
			Code: CodeMechanicalReportShape, File: report.path, Line: 1,
			Detail: "Results table has no report rows",
			Fix:    "Materialize every planned QA row with one terminal status.",
		})
	}
	for _, row := range report.rows {
		status := strings.TrimSpace(row.status)
		lower := strings.ToLower(status)
		switch {
		// The row a gate writes when a precondition stopped it before it built
		// a matrix: the bare blocked status no typed cause covers, with
		// provenance precondition. Both cells are read, because the status
		// alone says only that nothing was measured and the provenance alone
		// would let any status claim to be a refusal.
		case strings.EqualFold(row.provenance, spec.QAPreconditionRowProvenance):
			if lower != spec.QAPreconditionRowStatus {
				addMechanicalFinding(result, MechanicalFinding{
					Code: CodeMechanicalReportShape, File: report.path, Line: row.line,
					Detail: "row " + row.id + " has provenance " + spec.QAPreconditionRowProvenance +
						" with status " + strconv.Quote(status),
					Fix: "Record the precondition refusal as status " + spec.QAPreconditionRowStatus +
						" beside provenance " + spec.QAPreconditionRowProvenance + ".",
					RowHint: row.id,
				})
				continue
			}
			actual["rows_blocked_precondition"]++
		case lower == "pass", lower == "fail", lower == "skipped", strings.HasPrefix(lower, "carried (") && strings.HasSuffix(lower, ")"):
		case lower == "pending" || lower == "":
			addMechanicalFinding(result, MechanicalFinding{
				Code: CodeMechanicalReportShape, File: report.path, Line: row.line,
				Detail:  "row " + row.id + " remains pending instead of carrying a terminal status",
				Fix:     "Set row " + row.id + " to pass, fail, a typed blocked status, carried, or skipped.",
				RowHint: row.id,
			})
		case typedBlockedStatus(lower, "environment"):
			actual["rows_blocked_environment"]++
		case typedBlockedStatus(lower, "finding") && strings.Contains(lower, blockedCauseLiteral):
			actual["rows_blocked_finding"]++
		case typedBlockedStatus(lower, "finding"):
			unparsed["rows_blocked_finding"]++
			addMechanicalFinding(result, MechanicalFinding{
				Code: CodeMechanicalReportShape, File: report.path, Line: row.line,
				Detail:  "row " + row.id + " has a finding-typed blocked status without required literal \" — waits on \"",
				Fix:     "Add the required literal \" — waits on \" after the finding type.",
				RowHint: row.id,
			})
		case typedBlockedStatus(lower, "declared"):
			actual["rows_blocked_declared"]++
		case strings.HasPrefix(lower, "blocked"):
			addMechanicalFinding(result, MechanicalFinding{
				Code: CodeMechanicalReportShape, File: report.path, Line: row.line,
				Detail:  "row " + row.id + " has a blocked cause outside environment, finding, or declared",
				Fix:     "Use blocked (environment: ...), blocked (finding: ...), or blocked (declared: ...).",
				RowHint: row.id,
			})
		default:
			addMechanicalFinding(result, MechanicalFinding{
				Code: CodeMechanicalReportShape, File: report.path, Line: row.line,
				Detail:  "row " + row.id + " has non-terminal status " + strconv.Quote(status),
				Fix:     "Replace the status with a terminal QA row status.",
				RowHint: row.id,
			})
		}
	}

	declared := map[string]int{
		"rows_blocked_environment": report.rowsBlockedEnvironment,
		"rows_blocked_finding":     report.rowsBlockedFinding,
		"rows_blocked_declared":    report.rowsBlockedDeclared,
	}
	if report.parseError != nil {
		// The frontmatter never parsed, so the declared counts and their line
		// numbers are unknown. The parse finding above already names the
		// repair; absence claims here would misdirect it.
		return
	}
	for _, field := range []string{"rows_blocked_environment", "rows_blocked_finding", "rows_blocked_declared"} {
		line := report.countLines[field]
		if line == 0 {
			line = 1
		}
		if _, present := report.countLines[field]; !present {
			addMechanicalFinding(result, MechanicalFinding{
				Code: CodeMechanicalReportShape, File: report.path, Line: line,
				Detail: field + " is absent from report frontmatter",
				Fix:    "Record all three typed blocked-cause counts in every closed report.",
			})
			continue
		}
		if declared[field] != actual[field] {
			if declared[field] == actual[field]+unparsed[field] {
				continue
			}
			addMechanicalCountMismatch(result, report, field, line, declared[field], actual[field])
		}
	}
	detectPreconditionCount(result, report, actual["rows_blocked_precondition"])
}

// detectPreconditionCount reconciles rows_blocked_precondition with the refusal
// rows the Results table actually carries. The field is required only of a
// report that records a refusal: every report closed before the refusal row
// existed declares three counts, and demanding a fourth from all of them would
// raise a fresh blocker on gates that never refused — the deadlock this row
// exists to end.
func detectPreconditionCount(result *MechanicalResult, report mechanicalReport, rows int) {
	const field = "rows_blocked_precondition"
	line, present := report.countLines[field]
	switch {
	case present:
		if report.rowsBlockedPrecondition != rows {
			addMechanicalCountMismatch(result, report, field, line, report.rowsBlockedPrecondition, rows)
		}
	case rows > 0:
		addMechanicalFinding(result, MechanicalFinding{
			Code: CodeMechanicalReportShape, File: report.path, Line: 1,
			Detail: field + " is absent from report frontmatter",
			Fix:    "Record " + field + " beside the precondition refusal row the gate wrote.",
		})
	}
}

func addMechanicalCountMismatch(result *MechanicalResult, report mechanicalReport, field string, line, declared, actual int) {
	if line < 1 {
		line = 1
	}
	addMechanicalFinding(result, MechanicalFinding{
		Code: CodeMechanicalReportShape, File: report.path, Line: line,
		Detail: fmt.Sprintf("%s is %d but the Results table contains %d matching rows", field, declared, actual),
		Fix:    "Set " + field + " to the exact number of matching typed blocked rows.",
	})
}

func detectMechanicalEvidencePaths(result *MechanicalResult, repoRoot string, report mechanicalReport) {
	reportDirectory := filepath.Dir(filepath.Join(repoRoot, filepath.FromSlash(report.path)))
	for _, row := range report.rows {
		for _, target := range mechanicalEvidenceTargets(row.evidence) {
			target = strings.TrimSpace(strings.Split(target, "#")[0])
			if target == "" || strings.Contains(target, "://") || strings.HasPrefix(target, "mailto:") {
				continue
			}
			if filepath.IsAbs(target) || strings.Contains(target, `\`) {
				addUnresolvedMechanicalEvidence(result, report, row, target)
				continue
			}
			base := reportDirectory
			switch {
			case strings.HasPrefix(target, "qa/"):
				base = filepath.Dir(reportDirectory)
			case strings.HasPrefix(target, "docs/"):
				base = repoRoot
			}
			resolved := filepath.Clean(filepath.Join(base, filepath.FromSlash(target)))
			relative, relErr := filepath.Rel(repoRoot, resolved)
			if relErr != nil || relative == ".." || strings.HasPrefix(relative, ".."+string(filepath.Separator)) {
				addUnresolvedMechanicalEvidence(result, report, row, target)
				continue
			}
			info, err := os.Stat(resolved)
			if err == nil && !info.IsDir() {
				continue
			}
			addUnresolvedMechanicalEvidence(result, report, row, target)
		}
	}
}

type mechanicalEvidenceEntry struct {
	mode   string
	object string
	path   string
}

func detectMechanicalEvidenceScratchState(ctx context.Context, result *MechanicalResult, repoRoot string, report mechanicalReport) error {
	evidenceRoot := cleanMechanicalPath(filepath.ToSlash(filepath.Join(filepath.Dir(filepath.FromSlash(report.path)), "evidence")))
	entries, err := mechanicalTrackedEvidence(ctx, repoRoot, evidenceRoot)
	if err != nil {
		return err
	}
	if len(entries) == 0 {
		resolved, ok := resolveRepositoryPath(repoRoot, evidenceRoot)
		if !ok {
			addMechanicalSkip(result, DetectorMechanicalEvidencePath, evidenceRoot)
			return nil
		}
		info, statErr := os.Stat(resolved)
		if errors.Is(statErr, os.ErrNotExist) || statErr == nil && !info.IsDir() {
			addMechanicalSkip(result, DetectorMechanicalEvidencePath, evidenceRoot)
			return nil
		}
		if statErr != nil {
			return fmt.Errorf("inspect mechanical evidence directory %q: %w", resolved, statErr)
		}
		return nil
	}

	for _, entry := range entries {
		switch entry.mode {
		case "160000":
			addMechanicalEvidenceScratchFinding(result, entry.path, "gitlink")
		case "100644", "100755":
			binary, err := mechanicalBlobIsBinary(ctx, repoRoot, entry.object)
			if err != nil {
				return fmt.Errorf("inspect mechanical evidence %q: %w", entry.path, err)
			}
			if binary {
				addMechanicalEvidenceScratchFinding(result, entry.path, "built binary")
			}
		}
	}
	return nil
}

func mechanicalTrackedEvidence(ctx context.Context, repoRoot, evidenceRoot string) ([]mechanicalEvidenceEntry, error) {
	command := exec.CommandContext(ctx, "git", "-C", repoRoot, "ls-files", "--stage", "-z", "--", evidenceRoot)
	output, err := command.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("list tracked mechanical evidence under %q: %w: %s", evidenceRoot, err, strings.TrimSpace(string(output)))
	}
	var entries []mechanicalEvidenceEntry
	for _, raw := range bytes.Split(output, []byte{0}) {
		if len(raw) == 0 {
			continue
		}
		header, rawPath, ok := bytes.Cut(raw, []byte{'\t'})
		fields := bytes.Fields(header)
		if !ok || len(fields) != 3 {
			return nil, fmt.Errorf("parse tracked mechanical evidence entry %q", raw)
		}
		path := string(rawPath)
		clean := cleanMechanicalPath(path)
		if clean == "" || clean != path || (clean != evidenceRoot && !strings.HasPrefix(clean, evidenceRoot+"/")) {
			return nil, fmt.Errorf("tracked mechanical evidence path escapes %q: %q", evidenceRoot, path)
		}
		if string(fields[2]) != "0" {
			continue
		}
		entries = append(entries, mechanicalEvidenceEntry{
			mode:   string(fields[0]),
			object: string(fields[1]),
			path:   path,
		})
	}
	return entries, nil
}

func mechanicalBlobIsBinary(ctx context.Context, repoRoot, object string) (bool, error) {
	const probeSize = 8000
	command := exec.CommandContext(ctx, "git", "-C", repoRoot, "cat-file", "blob", object)
	var stderr bytes.Buffer
	command.Stderr = &stderr
	stdout, err := command.StdoutPipe()
	if err != nil {
		return false, fmt.Errorf("open git blob %s: %w", object, err)
	}
	if err := command.Start(); err != nil {
		return false, fmt.Errorf("start git blob %s: %w", object, err)
	}
	probe := make([]byte, probeSize)
	n, readErr := io.ReadFull(stdout, probe)
	if readErr != nil && !errors.Is(readErr, io.EOF) && !errors.Is(readErr, io.ErrUnexpectedEOF) {
		_ = command.Wait()
		return false, fmt.Errorf("read git blob %s: %w", object, readErr)
	}
	if _, err := io.Copy(io.Discard, stdout); err != nil {
		_ = command.Wait()
		return false, fmt.Errorf("drain git blob %s: %w", object, err)
	}
	if err := command.Wait(); err != nil {
		return false, fmt.Errorf("git cat-file blob %s: %w: %s", object, err, strings.TrimSpace(stderr.String()))
	}
	return bytes.IndexByte(probe[:n], 0) >= 0, nil
}

func addMechanicalEvidenceScratchFinding(result *MechanicalResult, path, kind string) {
	addMechanicalFinding(result, MechanicalFinding{
		Code:   CodeMechanicalEvidencePath,
		File:   path,
		Line:   1,
		Detail: "tracked evidence path " + path + " is a " + kind,
		Fix:    "Remove the gate's scratch state from the Spec evidence directory and keep only reader-openable artifacts.",
	})
}

func mechanicalEvidenceTargets(value string) []string {
	seen := make(map[string]bool)
	var targets []string
	add := func(target string) {
		target = strings.TrimSpace(target)
		if target == "" || seen[target] {
			return
		}
		seen[target] = true
		targets = append(targets, target)
	}
	for _, match := range markdownLinkTargetPattern.FindAllStringSubmatch(value, -1) {
		add(match[1])
	}
	for _, match := range evidencePathPattern.FindAllStringSubmatch(value, -1) {
		add(match[1])
	}
	for _, match := range backtickPattern.FindAllStringSubmatch(value, -1) {
		if evidencePathPattern.MatchString(match[1]) {
			add(match[1])
		}
	}
	return targets
}

func typedBlockedStatus(status, cause string) bool {
	prefix := "blocked (" + cause + ":"
	return strings.HasPrefix(status, prefix) && strings.HasSuffix(status, ")") && strings.TrimSpace(status[len(prefix):len(status)-1]) != ""
}

func addUnresolvedMechanicalEvidence(result *MechanicalResult, report mechanicalReport, row mechanicalReportRow, target string) {
	addMechanicalFinding(result, MechanicalFinding{
		Code: CodeMechanicalEvidencePath, File: report.path, Line: row.line,
		Detail:  "row " + row.id + " cites unresolved evidence path " + target,
		Fix:     "Create the cited evidence artifact or update the row to its resolving report-relative path.",
		RowHint: row.id,
	})
}

func mechanicalCommitExists(ctx context.Context, repoRoot, sha string) (bool, error) {
	command := exec.CommandContext(ctx, "git", "-C", repoRoot, "cat-file", "-e", sha+"^{commit}")
	output, err := command.CombinedOutput()
	if err == nil {
		return true, nil
	}
	var exitErr *exec.ExitError
	if errors.As(err, &exitErr) {
		return false, nil
	}
	return false, fmt.Errorf("inspect Git commit %q: %w: %s", sha, err, strings.TrimSpace(string(output)))
}

func mechanicalChangedPaths(ctx context.Context, repoRoot, sha string) ([]string, error) {
	command := exec.CommandContext(ctx, "git", "-C", repoRoot, "diff-tree", "--no-commit-id", "--name-only", "-r", "--root", sha)
	output, err := command.Output()
	if err != nil {
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) && len(bytes.TrimSpace(exitErr.Stderr)) > 0 {
			return nil, fmt.Errorf("read changed paths for Git commit %q: %w: %s", sha, err, bytes.TrimSpace(exitErr.Stderr))
		}
		return nil, fmt.Errorf("read changed paths for Git commit %q: %w", sha, err)
	}
	var paths []string
	for _, line := range strings.Split(string(output), "\n") {
		if clean := cleanMechanicalPath(line); clean != "" {
			paths = append(paths, clean)
		}
	}
	sort.Strings(paths)
	return paths, nil
}

func mechanicalIsAncestor(ctx context.Context, repoRoot, ancestor, descendant string) (bool, error) {
	command := exec.CommandContext(ctx, "git", "-C", repoRoot, "merge-base", "--is-ancestor", ancestor, descendant)
	output, err := command.CombinedOutput()
	if err == nil {
		return true, nil
	}
	var exitErr *exec.ExitError
	if errors.As(err, &exitErr) && exitErr.ExitCode() == 1 {
		return false, nil
	}
	return false, fmt.Errorf("compare Git commit order %q..%q: %w: %s", ancestor, descendant, err, strings.TrimSpace(string(output)))
}

func addMechanicalFinding(result *MechanicalResult, finding MechanicalFinding) {
	if finding.Line < 1 {
		finding.Line = 1
	}
	result.Findings = append(result.Findings, finding)
}

func addMechanicalSkip(result *MechanicalResult, detector, missing string) {
	result.Skips = append(result.Skips, MechanicalSkip{Detector: detector, MissingArtifact: missing})
}

func materializeBlockedRows(result *MechanicalResult) {
	seen := make(map[string]bool)
	for _, finding := range result.Findings {
		rowID := strings.TrimSpace(finding.RowHint)
		if rowID == "" {
			rowID = finding.Code
		}
		key := rowID + "\x00" + finding.Code
		if seen[key] {
			continue
		}
		seen[key] = true
		result.Blocked = append(result.Blocked, BlockedRow{
			ID:          rowID,
			FindingCode: finding.Code,
			WaitingOn:   finding.Detail,
		})
	}
}
