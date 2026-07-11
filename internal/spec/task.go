package spec

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

type taskFrontmatter struct {
	Task       string `yaml:"task"`
	Spec       string `yaml:"spec"`
	Status     string `yaml:"status"`
	Type       string `yaml:"type"`
	Complexity string `yaml:"complexity"`
}

type taskDocument struct {
	Frontmatter  taskFrontmatter
	Title        string
	Context      []TaskContextRef
	Verification []string
}

const maxTaskContextRefs = 50

// ReloadTask re-reads a Task's file from the Spec Root — typically after an
// Agent has modified it — refreshing Status, Title, Type, and Verification in
// place. The Task Graph fields (ID, File, Needs) belong to the manifest and
// are left alone.
func ReloadTask(specsRoot string, task *Task) error {
	path := filepath.Join(specsRoot, task.File)
	content, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("read Task %q file %q: %w", task.ID, path, err)
	}
	document, err := parseTaskDocument(content)
	if err != nil {
		return TaskFileError{TaskID: task.ID, Path: path, Err: err}
	}
	if len(document.Verification) == 0 {
		return MissingVerificationError{TaskID: task.ID, Path: path}
	}
	task.Title = document.Title
	task.Status = Status(document.Frontmatter.Status)
	task.Type = document.Frontmatter.Type
	task.Context = append([]TaskContextRef(nil), document.Context...)
	task.Verification = document.Verification
	return nil
}

// SetStatus rewrites only the status frontmatter value of a task file,
// preserving every other byte of the file exactly.
func SetStatus(taskPath string, status Status) error {
	if !AllowedStatus(status) {
		return fmt.Errorf("Task status %q is not allowed", status)
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

func parseTaskDocument(content []byte) (taskDocument, error) {
	frontmatterBytes, body, err := splitFrontmatter(content)
	if err != nil {
		return taskDocument{}, err
	}
	var frontmatter taskFrontmatter
	if err := yaml.Unmarshal(frontmatterBytes, &frontmatter); err != nil {
		return taskDocument{}, fmt.Errorf("parse frontmatter: %w", err)
	}
	if !AllowedStatus(Status(frontmatter.Status)) {
		return taskDocument{}, fmt.Errorf("unsupported status %q", frontmatter.Status)
	}
	contextRefs, err := parseTaskContextRefs(body)
	if err != nil {
		return taskDocument{}, err
	}
	return taskDocument{
		Frontmatter:  frontmatter,
		Title:        parseTaskTitle(body),
		Context:      contextRefs,
		Verification: parseVerificationCommands(body),
	}, nil
}

// parseTaskTitle extracts the title from the first level-one heading,
// dropping the "Task NN:" prefix the task template mandates.
func parseTaskTitle(body []byte) string {
	for _, line := range strings.Split(string(body), "\n") {
		line = strings.TrimSpace(line)
		if !strings.HasPrefix(line, "# ") {
			continue
		}
		title := strings.TrimSpace(strings.TrimPrefix(line, "# "))
		if strings.HasPrefix(title, "Task ") {
			if marker := strings.Index(title, ":"); marker >= 0 {
				return strings.TrimSpace(title[marker+1:])
			}
		}
		return title
	}
	return ""
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
// optional ## Context section. Duplicate entries are ignored without changing
// the order of first occurrence.
func parseTaskContextRefs(body []byte) ([]TaskContextRef, error) {
	var refs []TaskContextRef
	seen := map[string]bool{}
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
		key := string(ref.Kind) + "\x00" + ref.Path
		if seen[key] {
			continue
		}
		seen[key] = true
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
		return TaskContextRef{}, TaskContextError{Kind: entry, Reason: `expected "instruction" or "interface" label`}
	}
	kind := ContextKind(strings.TrimSpace(label))
	switch kind {
	case ContextKindInstruction, ContextKindInterface:
	default:
		return TaskContextRef{}, TaskContextError{Kind: strings.TrimSpace(label), Reason: `expected "instruction" or "interface" label`}
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
