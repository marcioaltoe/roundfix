package speccheck

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"gopkg.in/yaml.v3"
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
	id       string
	status   string
	evidence string
	line     int
}

type mechanicalReport struct {
	path                   string
	rows                   []mechanicalReportRow
	rowsBlockedEnvironment int
	rowsBlockedFinding     int
	rowsBlockedDeclared    int
	countLines             map[string]int
	parseError             error
}

var (
	markdownLinkTargetPattern = regexp.MustCompile(`\[[^\]]*\]\(([^)]+)\)`)
	evidencePathPattern       = regexp.MustCompile(`(?:^|[[:space:],;])((?:qa/)?evidence/[A-Za-z0-9._/-]+)`)
)

// RunMechanicalStage evaluates every detector and returns all findings in one
// pass. It reads Git and files only; it writes no report and settles no Task.
func RunMechanicalStage(ctx context.Context, request MechanicalRequest) (MechanicalResult, error) {
	result := MechanicalResult{
		Findings: []MechanicalFinding{},
		Carried:  []CarriedRow{},
		Blocked:  []BlockedRow{},
		Skips:    []MechanicalSkip{},
	}
	repoRoot := filepath.Clean(request.RepoRoot)
	if strings.TrimSpace(request.RepoRoot) == "" {
		return result, errors.New("run mechanical stage: repository root is empty")
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
	}

	result.Blocking = len(result.Findings) > 0
	materializeBlockedRows(&result)
	return result, nil
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

	allowed := parseMechanicalAuthorizationPaths(content)
	if len(allowed) == 0 {
		addMechanicalFinding(result, MechanicalFinding{
			Code: CodeMechanicalAuthPaths, File: request.AuthorizationPath, Line: 1,
			Detail: "authorization declares no exact bounded files",
			Fix:    "Declare every authorized repository-relative path in the authorization artifact.",
		})
		return nil
	}

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
			if allowed[changedPath] || changedPath == taskFile {
				continue
			}
			addMechanicalFinding(result, MechanicalFinding{
				Code: CodeMechanicalAuthPaths, File: request.AuthorizationPath, Line: 1,
				Detail:  fmt.Sprintf("Task %s commit %s changes %s outside the exact bounded files", taskCommit.TaskID, taskCommit.SHA, changedPath),
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
	report := mechanicalReport{path: path, countLines: make(map[string]int)}
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
			counts, lines, err := mechanicalBlockedCounts(document)
			if err != nil {
				report.parseError = err
			} else {
				report.rowsBlockedEnvironment = counts["rows_blocked_environment"]
				report.rowsBlockedFinding = counts["rows_blocked_finding"]
				report.rowsBlockedDeclared = counts["rows_blocked_declared"]
				report.countLines = lines
			}
		}
	}

	allLines := strings.Split(text, "\n")
	inResults := false
	statusColumn, evidenceColumn := -1, -1
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
		report.rows = append(report.rows, row)
	}
	return report
}

func mechanicalBlockedCounts(document yaml.Node) (map[string]int, map[string]int, error) {
	counts := map[string]int{
		"rows_blocked_environment": 0,
		"rows_blocked_finding":     0,
		"rows_blocked_declared":    0,
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
		case typedBlockedStatus(lower, "finding") && strings.Contains(lower, " — waits on "):
			actual["rows_blocked_finding"]++
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
			addMechanicalFinding(result, MechanicalFinding{
				Code: CodeMechanicalReportShape, File: report.path, Line: line,
				Detail: fmt.Sprintf("%s is %d but the Results table contains %d matching rows", field, declared[field], actual[field]),
				Fix:    "Set " + field + " to the exact number of matching typed blocked rows.",
			})
		}
	}
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
	output, err := command.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("read changed paths for Git commit %q: %w: %s", sha, err, strings.TrimSpace(string(output)))
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
		if finding.RowHint == "" {
			continue
		}
		key := finding.RowHint + "\x00" + finding.Code
		if seen[key] {
			continue
		}
		seen[key] = true
		result.Blocked = append(result.Blocked, BlockedRow{
			ID:          finding.RowHint,
			FindingCode: finding.Code,
			WaitingOn:   finding.Detail,
		})
	}
}
