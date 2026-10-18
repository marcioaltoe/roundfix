package speccheck

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"roundfix/internal/spec"
)

// CodeVerifyPathUndeclared identifies a Verification path the plan does not provide.
const CodeVerifyPathUndeclared = "SC-VERIFY-PATH-UNDECLARED"

var verificationPathPattern = regexp.MustCompile(`^[A-Za-z0-9._@+-]+(?:/[A-Za-z0-9._@+-]+)+/?$`)

func detectVerificationPaths(result *Result, repoRoot, specDir string, graph *spec.Graph) error {
	if graph == nil {
		addSkip(result, CodeVerifyPathUndeclared, artifactDisplayPath(repoRoot, filepath.Join(specDir, "_tasks.md")))
		return nil
	}
	tasks := make(map[string]spec.Task, len(graph.Tasks))
	for _, task := range graph.Tasks {
		tasks[task.ID] = task
	}
	for _, task := range graph.Tasks {
		if task.Status != spec.StatusPending || task.ID == graph.QATaskID || task.Type == spec.TaskTypeQA {
			continue
		}
		taskPath := filepath.Join(filepath.Dir(specDir), filepath.FromSlash(task.File))
		content, err := os.ReadFile(taskPath)
		if err != nil {
			return fmt.Errorf("read Verification Task %q: %w", taskPath, err)
		}
		var provided []string
		for _, ref := range task.Context {
			provided = append(provided, ref.Path)
		}
		provided = append(provided, spec.RecordedTaskPaths(content)...)
		visited := make(map[string]bool)
		var addCreates func(string)
		addCreates = func(id string) {
			if visited[id] {
				return
			}
			visited[id] = true
			dependency := tasks[id]
			for _, ref := range dependency.Context {
				if ref.Kind == spec.ContextKindCreates {
					provided = append(provided, ref.Path)
				}
			}
			for _, need := range dependency.Needs {
				addCreates(need)
			}
		}
		for _, need := range task.Needs {
			addCreates(need)
		}
		for _, command := range task.Verification {
			for _, word := range verificationPathWords(command) {
				candidate := strings.TrimPrefix(word, "./")
				if !verificationPathPattern.MatchString(candidate) || strings.Contains(candidate, "...") {
					continue
				}
				absolute, ok := resolveRepositoryPath(repoRoot, candidate)
				if !ok {
					continue
				}
				first, _, _ := strings.Cut(candidate, "/")
				if first == "." || first == ".." {
					continue
				}
				info, err := os.Stat(filepath.Join(repoRoot, first))
				if errors.Is(err, os.ErrNotExist) {
					continue
				}
				if err != nil {
					return fmt.Errorf("stat Verification directory %q: %w", first, err)
				}
				if !info.IsDir() {
					continue
				}
				// The Spec Root can be configured outside docs/specs.
				relative, err := filepath.Rel(filepath.Dir(specDir), absolute)
				if err == nil && relative != ".." && !strings.HasPrefix(relative, ".."+string(filepath.Separator)) && !filepath.IsAbs(relative) {
					continue
				}
				candidate = filepath.ToSlash(filepath.Clean(candidate))
				supplied := false
				for _, declared := range provided {
					declared = filepath.ToSlash(filepath.Clean(declared))
					if candidate == declared || strings.HasPrefix(declared, candidate+"/") {
						supplied = true
						break
					}
				}
				if supplied {
					continue
				}
				if _, err := os.Lstat(absolute); err == nil {
					continue
				} else if !errors.Is(err, os.ErrNotExist) {
					return fmt.Errorf("lstat Verification path %q: %w", absolute, err)
				}
				result.Findings = append(result.Findings, Finding{
					Code: CodeVerifyPathUndeclared, Severity: SeverityError,
					Summary: fmt.Sprintf("Task %s Verification names %s, which the Task does not declare, no Task it depends on creates, and the checked tree does not hold", task.ID, candidate),
					Where:   []Location{{Path: artifactDisplayPath(repoRoot, taskPath), Line: sectionLineContaining(content, "Verification", command)}},
					Fix:     fmt.Sprintf("Declare %s in this Task, create it in a Task this Task needs, or correct the command, following the write-tasks rule \"Verification paths\".", candidate),
				})
			}
		}
	}
	return nil
}

// Read operands only; the shell is never invoked. Redirection targets form
// their own segment so printed messages cannot hide a redirected file.
func verificationPathWords(command string) []string {
	words := shellWords(command)
	var paths []string
	start := 0
	for index, word := range words {
		switch word.text {
		case ";", "&&", "||", "|", "(", ")", "<", ">", ">>":
			paths = append(paths, verificationPathOperands(words[start:index])...)
			start = index + 1
		}
	}
	return append(paths, verificationPathOperands(words[start:])...)
}

func verificationPathOperands(words []shellWord) []string {
	if len(words) == 0 {
		return nil
	}
	index := 0
	// These wrappers occur in repository-authored Verification.
	for index < len(words) && (words[index].text == "!" || words[index].text == "rtk" || words[index].text == "proxy" || words[index].text == "{" || words[index].text == "}") {
		index++
	}
	if index == len(words) {
		return nil
	}
	command := filepath.Base(words[index].text)
	switch command {
	case "printf", "echo":
		return nil
	case "test", "[":
		if index+1 < len(words) && words[index+1].text == "!" {
			return nil
		}
	}
	excluded := make(map[int]bool)
	if command == "grep" || command == "rg" {
		explicit := false
		options := true
		firstOperand := -1
		for i := index + 1; i < len(words); i++ {
			word := words[i].text
			if options && word == "--" {
				options = false
				continue
			}
			if options && (word == "-e" || word == "-f" || word == "--regexp" || word == "--file") {
				explicit = true
				if i+1 < len(words) {
					i++
					excluded[i] = true
				}
				continue
			}
			if options && (strings.HasPrefix(word, "--regexp=") || strings.HasPrefix(word, "--file=")) {
				explicit = true
				continue
			}
			if options && strings.HasPrefix(word, "-") && word != "-" {
				// Short options can combine, for example -qFe PATTERN.
				if flag := strings.IndexAny(word[1:], "ef"); !strings.HasPrefix(word, "--") && flag >= 0 {
					explicit = true
					if flag+2 == len(word) {
						if i+1 < len(words) {
							i++
							excluded[i] = true
						}
					}
				} else if grepOptionTakesValue(word) && i+1 < len(words) {
					i++
				}
				continue
			}
			if firstOperand < 0 {
				firstOperand = i
			}
		}
		if !explicit && firstOperand >= 0 {
			excluded[firstOperand] = true
		}
	}
	var operands []string
	for i, word := range words {
		if !excluded[i] {
			operands = append(operands, word.text)
		}
	}
	return operands
}
