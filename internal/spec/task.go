package spec

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"gopkg.in/yaml.v3"
)

type taskFrontmatter struct {
	Task            string           `yaml:"task"`
	Spec            string           `yaml:"spec"`
	Status          string           `yaml:"status"`
	Type            string           `yaml:"type"`
	Complexity      string           `yaml:"complexity"`
	TaskRepairPaths []string         `yaml:"repair_paths"`
	AssignedRepairs []AssignedRepair `yaml:"assigned_repairs"`
}

type taskDocument struct {
	Frontmatter      taskFrontmatter
	Title            string
	TitleLine        int
	StatusNormalized bool
	Type             TaskType
	Context          []TaskContextRef
	Requirements     []TaskDeclaration
	RehearsalCases   []TaskDeclaration
	Verification     []string
	NegativeControl  []string
	TaskRepairPaths  []string
	AssignedRepairs  []AssignedRepair
}

type taskMarkdownLine struct {
	Text string
	Line int
}

const maxTaskContextRefs = 50

// ReloadTask re-reads a Task's file from the Spec Root — typically after an
// Agent has modified it — refreshing Status, Title, Type, declarations, and
// Verification in place. The Task Graph fields (ID, File, Needs) belong to the
// manifest and are left alone.
func ReloadTask(specsRoot string, task *Task) error {
	path := filepath.Join(specsRoot, task.File)
	content, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("read Task %q file %q: %w", task.ID, path, err)
	}
	document, err := parseTaskDocument(content, path)
	if err != nil {
		return TaskFileError{TaskID: task.ID, Path: path, Err: err}
	}
	if len(document.Verification) == 0 {
		return MissingVerificationError{TaskID: task.ID, Path: path}
	}
	verification, authoredQA := taskVerification(taskSlug(task.File), document.Type, document.Verification)
	task.Title = document.Title
	task.TitleLine = document.TitleLine
	task.Status = Status(document.Frontmatter.Status)
	task.StatusNormalized = document.StatusNormalized
	task.Type = document.Type
	task.Context = append([]TaskContextRef(nil), document.Context...)
	task.Requirements = append([]TaskDeclaration(nil), document.Requirements...)
	task.RehearsalCases = append([]TaskDeclaration(nil), document.RehearsalCases...)
	task.Verification = verification
	task.NegativeControl = append([]string(nil), document.NegativeControl...)
	task.TaskRepairPaths = append([]string(nil), document.TaskRepairPaths...)
	task.AssignedRepairs = append([]AssignedRepair(nil), document.AssignedRepairs...)
	task.qaAuthored = authoredQA
	return nil
}

// DerivedQAVerification is the Verification Roundfix supplies for a Task of
// type qa. It orders QA Reports by their embedded date and numeric rerun
// sequence, then accepts only the exact passing verdict from the newest one.
// The command is rendered into the Task file so readers can see the contract,
// but changing that rendered command does not change the effective contract.
func DerivedQAVerification(slug string) []string {
	reportDir := filepath.ToSlash(filepath.Join("docs", "specs", slug, "qa"))
	command := `newest="$(find ` + reportDir + ` -type f -name "qa-report-*.md" -print 2>/dev/null | ` +
		`awk "{ report=\$0; name=\$0; parts=split(name, path, \"/\"); name=path[parts]; name=substr(name, 11, length(name)-13); ` +
		`date=substr(name, 1, 10); suffix=substr(name, 11); ` +
		`shape=(length(date) == 10 && substr(date, 5, 1) == \"-\" && substr(date, 8, 1) == \"-\"); ` +
		`for (i=1; i <= 10 && shape; i++) { if (i != 5 && i != 8 && index(\"0123456789\", substr(date, i, 1)) == 0) shape=0 } ` +
		`year=substr(date, 1, 4)+0; month=substr(date, 6, 2)+0; day=substr(date, 9, 2)+0; ` +
		`leap=(year%400 == 0 || (year%4 == 0 && year%100 != 0)); days=(month == 2 ? 28+leap : ((month == 4 || month == 6 || month == 9 || month == 11) ? 30 : 31)); ` +
		`dated=(shape && month >= 1 && month <= 12 && day >= 1 && day <= days); if (!dated) date=\"\"; ` +
		`if (suffix == \"\") { sequenced=1; sequence=-1 } else { sequenced=(length(suffix) > 1 && substr(suffix, 1, 1) == \"-\"); ` +
		`for (i=2; i <= length(suffix) && sequenced; i++) { if (index(\"0123456789\", substr(suffix, i, 1)) == 0) sequenced=0 } sequence=(sequenced ? substr(suffix, 2)+0 : 0) } ` +
		`printf \"%s\\t%s\\t%s\\t%s\\t%s\\n\", dated, date, sequenced, sequence, report }" | ` +
		`sort -k1,1n -k2,2 -k3,3n -k4,4n -k5,5 | tail -1 | cut -f5-)"; ` +
		`test -n "$newest" || exit 1; ` +
		`awk "BEGIN { whitespace=\" \t\r\n\f\v\" } NR == 1 && \$0 == \"---\" { frontmatter=1; next } frontmatter && \$0 == \"---\" { closed=1; exit } ` +
		`frontmatter && index(\$0, \"verdict:\") == 1 { verdict=substr(\$0, 9); ` +
		`while (length(verdict) > 0 && index(whitespace, substr(verdict, 1, 1)) > 0) verdict=substr(verdict, 2); ` +
		`while (length(verdict) > 0 && index(whitespace, substr(verdict, length(verdict), 1)) > 0) verdict=substr(verdict, 1, length(verdict)-1); verdicts++ } ` +
		`END { exit(closed && verdicts == 1 && verdict == \"pass\" ? 0 : 1) }" "$newest"`
	return []string{command}
}

// AuthoredQAVerification reports a qa Task whose rendered command differs
// from the Verification Roundfix derives. Callers refuse this condition rather
// than silently replacing the author's visible text.
func AuthoredQAVerification(task Task) bool {
	return task.Type == TaskTypeQA && task.qaAuthored
}

func taskVerification(slug string, taskType TaskType, rendered []string) ([]string, bool) {
	if taskType != TaskTypeQA {
		return append([]string(nil), rendered...), false
	}
	derived := DerivedQAVerification(slug)
	return derived, !slices.Equal(rendered, derived)
}

func taskSlug(taskFile string) string {
	clean := filepath.Clean(taskFile)
	if separator := strings.IndexRune(clean, filepath.Separator); separator >= 0 {
		return clean[:separator]
	}
	return filepath.Base(filepath.Dir(clean))
}

// SetStatus rewrites only the status frontmatter value of a task file,
// preserving every other byte of the file exactly.
func SetStatus(taskPath string, status Status) error {
	if !AllowedStatus(status) {
		return fmt.Errorf("Task status %q is not allowed (allowed: %s)", status, allowedStatusValues())
	}
	info, err := os.Stat(taskPath)
	if err != nil {
		return fmt.Errorf("stat task file %q: %w", taskPath, err)
	}
	content, err := os.ReadFile(taskPath)
	if err != nil {
		return fmt.Errorf("read task file %q: %w", taskPath, err)
	}
	updated, err := rewriteStatus(content, status)
	if err != nil {
		return fmt.Errorf("rewrite status in task file %q: %w", taskPath, err)
	}
	if err := os.WriteFile(taskPath, updated, info.Mode().Perm()); err != nil {
		return fmt.Errorf("write task file %q: %w", taskPath, err)
	}
	return nil
}

func parseTaskDocument(content []byte, taskPath string) (taskDocument, error) {
	frontmatterBytes, body, err := splitFrontmatter(content)
	if err != nil {
		return taskDocument{}, err
	}
	var frontmatter taskFrontmatter
	if err := yaml.Unmarshal(frontmatterBytes, &frontmatter); err != nil {
		return taskDocument{}, fmt.Errorf("parse frontmatter: %w", err)
	}
	rawStatus := frontmatter.Status
	normalizedStatus := NormalizeStatus(rawStatus)
	if !AllowedStatus(Status(normalizedStatus)) {
		return taskDocument{}, fmt.Errorf("unsupported status %q (allowed: %s)", frontmatter.Status, allowedStatusValues())
	}
	frontmatter.Status = normalizedStatus
	taskType, err := ParseTaskType(taskPath, frontmatter.Type)
	if err != nil {
		return taskDocument{}, err
	}
	contextRefs, err := parseTaskContextRefs(body)
	if err != nil {
		return taskDocument{}, err
	}
	bodyLineOffset := bytes.Count(content[:len(content)-len(body)], []byte{'\n'})
	title, titleLine := parseTaskTitle(body, bodyLineOffset)
	return taskDocument{
		Frontmatter:      frontmatter,
		Title:            title,
		TitleLine:        titleLine,
		StatusNormalized: strings.TrimSpace(rawStatus) != normalizedStatus,
		Type:             taskType,
		Context:          contextRefs,
		Requirements:     parseTaskRequirements(body, bodyLineOffset),
		RehearsalCases:   parseTaskSectionBullets(body, "Rehearsal Cases", bodyLineOffset),
		Verification:     parseVerificationCommands(body),
		NegativeControl:  parseNegativeControlDeclarations(body),
		TaskRepairPaths:  append([]string(nil), frontmatter.TaskRepairPaths...),
		AssignedRepairs:  append([]AssignedRepair(nil), frontmatter.AssignedRepairs...),
	}, nil
}

// parseTaskTitle extracts the title from the first level-one heading,
// dropping the "Task NN:" prefix the task template mandates.
func parseTaskTitle(body []byte, lineOffset int) (string, int) {
	for index, line := range strings.Split(string(body), "\n") {
		line = strings.TrimSpace(line)
		if !strings.HasPrefix(line, "# ") {
			continue
		}
		title := strings.TrimSpace(strings.TrimPrefix(line, "# "))
		if strings.HasPrefix(title, "Task ") {
			if marker := strings.Index(title, ":"); marker >= 0 {
				return strings.TrimSpace(title[marker+1:]), lineOffset + index + 1
			}
		}
		return title, lineOffset + index + 1
	}
	return "", 1
}

// parseTaskRequirements extracts numbered requirements and joins their
// continuation lines without interpreting the requirement prose.
func parseTaskRequirements(body []byte, lineOffset int) []TaskDeclaration {
	var requirements []TaskDeclaration
	for _, sourceLine := range taskSectionLines(body, "Requirements", lineOffset) {
		trimmed := strings.TrimSpace(sourceLine.Text)
		if trimmed == "" {
			continue
		}
		if item, ok := numberedTaskItem(trimmed); ok {
			requirements = append(requirements, TaskDeclaration{Text: item, Line: sourceLine.Line})
			continue
		}
		if len(requirements) > 0 {
			requirements[len(requirements)-1].Text += " " + trimmed
		}
	}
	return requirements
}

func parseTaskSectionBullets(body []byte, heading string, lineOffset int) []TaskDeclaration {
	var entries []TaskDeclaration
	for _, sourceLine := range taskSectionLines(body, heading, lineOffset) {
		trimmed := strings.TrimSpace(sourceLine.Text)
		if strings.HasPrefix(trimmed, "- ") {
			entries = append(entries, TaskDeclaration{
				Text: strings.TrimSpace(strings.TrimPrefix(trimmed, "- ")),
				Line: sourceLine.Line,
			})
		}
	}
	return entries
}

func taskSectionLines(body []byte, heading string, lineOffset int) []taskMarkdownLine {
	var lines []taskMarkdownLine
	inSection := false
	for index, line := range strings.Split(string(body), "\n") {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "## ") {
			if inSection {
				break
			}
			inSection = strings.TrimSpace(strings.TrimPrefix(trimmed, "## ")) == heading
			continue
		}
		if inSection && isTaskMarkdownHeading(trimmed) {
			break
		}
		if inSection {
			lines = append(lines, taskMarkdownLine{Text: line, Line: lineOffset + index + 1})
		}
	}
	return lines
}

func isTaskMarkdownHeading(line string) bool {
	level := 0
	for level < len(line) && level < 6 && line[level] == '#' {
		level++
	}
	return level > 0 && (level == len(line) || line[level] == ' ' || line[level] == '\t')
}

func numberedTaskItem(line string) (string, bool) {
	dot := strings.IndexByte(line, '.')
	if dot < 1 || dot+1 >= len(line) || (line[dot+1] != ' ' && line[dot+1] != '\t') {
		return "", false
	}
	for _, char := range line[:dot] {
		if char < '0' || char > '9' {
			return "", false
		}
	}
	item := strings.TrimSpace(line[dot+1:])
	return item, item != ""
}

// parseVerificationCommands extracts commands verbatim from the backticked
// bullet entries of every ## Verification section, in section order. Bullets
// without a backticked span carry no command and are skipped.
func parseVerificationCommands(body []byte) []string {
	var commands []string
	inSection := false
	for _, line := range strings.Split(string(body), "\n") {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "# ") {
			inSection = false
			continue
		}
		if strings.HasPrefix(trimmed, "## ") {
			inSection = strings.TrimSpace(strings.TrimPrefix(trimmed, "## ")) == "Verification"
			continue
		}
		if !inSection || !strings.HasPrefix(trimmed, "- ") {
			continue
		}
		if command, ok := firstBacktickSpan(trimmed); ok && command != "" {
			commands = append(commands, command)
		}
	}
	return commands
}

// parseNegativeControlDeclarations extracts declarations verbatim from the
// backticked bullet entries of every ## Negative Control section, in section
// order. It carries declarations only; executing them belongs to no parser.
func parseNegativeControlDeclarations(body []byte) []string {
	var declarations []string
	inSection := false
	for _, line := range strings.Split(string(body), "\n") {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "# ") {
			inSection = false
			continue
		}
		if strings.HasPrefix(trimmed, "## ") {
			inSection = strings.TrimSpace(strings.TrimPrefix(trimmed, "## ")) == "Negative Control"
			continue
		}
		if !inSection || !strings.HasPrefix(trimmed, "- ") {
			continue
		}
		if declaration, ok := firstBacktickSpan(trimmed); ok && declaration != "" {
			declarations = append(declarations, declaration)
		}
	}
	return declarations
}

func firstBacktickSpan(line string) (string, bool) {
	start := strings.IndexByte(line, '`')
	if start < 0 {
		return "", false
	}
	length := strings.IndexByte(line[start+1:], '`')
	if length < 0 {
		return "", false
	}
	return line[start+1 : start+1+length], true
}

// parseTaskContextRefs extracts labeled repository-relative paths from the
// optional ## Context section. Duplicate input entries are ignored without
// changing the order of first occurrence; a declared output must be unique.
func parseTaskContextRefs(body []byte) ([]TaskContextRef, error) {
	var refs []TaskContextRef
	seen := map[string]bool{}
	pathKinds := map[string]ContextKind{}
	inSection := false
	for _, line := range strings.Split(string(body), "\n") {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "# ") {
			inSection = false
			continue
		}
		if strings.HasPrefix(trimmed, "## ") {
			inSection = strings.TrimSpace(strings.TrimPrefix(trimmed, "## ")) == "Context"
			continue
		}
		if !inSection || trimmed == "" {
			continue
		}
		if !strings.HasPrefix(trimmed, "- ") {
			continue
		}
		ref, err := parseTaskContextLine(trimmed)
		if err != nil {
			return nil, err
		}
		if previousKind, ok := pathKinds[ref.Path]; ok &&
			(ref.Kind == ContextKindCreates || previousKind == ContextKindCreates) {
			return nil, TaskContextError{
				Kind:   string(ref.Kind),
				Path:   ref.Path,
				Reason: "path must be unique within Task Context",
			}
		}
		key := string(ref.Kind) + "\x00" + ref.Path
		if seen[key] {
			continue
		}
		seen[key] = true
		pathKinds[ref.Path] = ref.Kind
		refs = append(refs, ref)
		if len(refs) > maxTaskContextRefs {
			return nil, TaskContextError{Reason: fmt.Sprintf("declares more than %d unique entries", maxTaskContextRefs)}
		}
	}
	return refs, nil
}

func parseTaskContextLine(line string) (TaskContextRef, error) {
	entry := strings.TrimSpace(strings.TrimPrefix(line, "- "))
	label, rest, ok := strings.Cut(entry, ":")
	if !ok {
		return TaskContextRef{}, TaskContextError{Kind: entry, Reason: `expected "instruction", "interface", or "creates" label`}
	}
	kind := ContextKind(strings.TrimSpace(label))
	switch kind {
	case ContextKindInstruction, ContextKindInterface, ContextKindCreates:
	default:
		return TaskContextRef{}, TaskContextError{Kind: strings.TrimSpace(label), Reason: `expected "instruction", "interface", or "creates" label`}
	}
	path := strings.TrimSpace(rest)
	if span, ok := firstBacktickSpan(path); ok {
		path = span
	}
	clean, err := cleanTaskContextPath(path)
	if err != nil {
		return TaskContextRef{}, TaskContextError{Kind: string(kind), Path: path, Reason: err.Error()}
	}
	return TaskContextRef{Kind: kind, Path: clean}, nil
}

func cleanTaskContextPath(path string) (string, error) {
	path = strings.TrimSpace(path)
	if path == "" {
		return "", errors.New("path is required")
	}
	if filepath.IsAbs(path) {
		return "", errors.New("path must be repository-relative")
	}
	if strings.Contains(path, `\`) {
		return "", errors.New("path must use slash separators")
	}
	clean := filepath.ToSlash(filepath.Clean(path))
	if clean == "." || clean == ".." || strings.HasPrefix(clean, "../") {
		return "", errors.New("path must stay inside the repository")
	}
	if clean != path {
		return "", errors.New("path must be clean")
	}
	return clean, nil
}

// rewriteStatus replaces the value of the first status field in the YAML
// frontmatter, byte-preserving everything around it, including any trailing
// comment on the status line.
func rewriteStatus(content []byte, status Status) ([]byte, error) {
	text := string(content)
	const opening = "---\n"
	if !strings.HasPrefix(text, opening) {
		return nil, errors.New("missing YAML frontmatter opening marker")
	}
	end := strings.Index(text[len(opening):], "\n---")
	if end < 0 {
		return nil, errors.New("missing YAML frontmatter closing marker")
	}
	frontmatter := text[len(opening) : len(opening)+end]
	lines := strings.Split(frontmatter, "\n")
	for index, line := range lines {
		if !strings.HasPrefix(line, "status:") {
			continue
		}
		lines[index] = rewriteStatusLine(line, status)
		return []byte(text[:len(opening)] + strings.Join(lines, "\n") + text[len(opening)+end:]), nil
	}
	return nil, errors.New("frontmatter has no status field")
}

func rewriteStatusLine(line string, status Status) string {
	rest := line[len("status:"):]
	valueStart := 0
	for valueStart < len(rest) && (rest[valueStart] == ' ' || rest[valueStart] == '\t') {
		valueStart++
	}
	remainder := rest[valueStart:]
	valueEnd := len(remainder)
	if marker := strings.Index(remainder, "#"); marker >= 0 {
		valueEnd = marker
	}
	valueLength := len(strings.TrimRight(remainder[:valueEnd], " \t"))
	return "status:" + rest[:valueStart] + string(status) + remainder[valueLength:]
}
