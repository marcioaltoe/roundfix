package speccheck

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"roundfix/internal/spec"
)

// CodeConsumersUnproven identifies a Task without a proved Consumer Declaration.
const CodeConsumersUnproven = "SC-CONSUMERS-UNPROVEN"

var (
	consumerEntryPattern      = regexp.MustCompile("^-\\s+(.+?)\\s+—\\s+consumers:\\s*(`[^`]+`(?:\\s*,\\s*`[^`]+`)*)\\s*$")
	consumerSuitePattern      = regexp.MustCompile("`([^`]+)`")
	consumerCheckpointPattern = regexp.MustCompile(`^-\s+checkpoint:\s*(task_[0-9]{2,})\s*$`)
)

type consumerEntry struct {
	text string
	line int
}

// detectConsumerDeclarations reads structure only, with no commit horizon and
// without executing any authored command. Graph Tasks are in topological order.
func detectConsumerDeclarations(result *Result, repoRoot string, graph *spec.Graph) error {
	tasks := make(map[string]spec.Task, len(graph.Tasks))
	ancestors := make(map[string]map[string]bool, len(graph.Tasks))
	for _, task := range graph.Tasks {
		tasks[task.ID] = task
		ancestors[task.ID] = make(map[string]bool)
		for _, need := range task.Needs {
			ancestors[task.ID][need] = true
			for ancestor := range ancestors[need] {
				ancestors[task.ID][ancestor] = true
			}
		}
	}
	for _, task := range graph.Tasks {
		if task.Status == spec.StatusCompleted || task.ID == graph.QATaskID || task.Type == spec.TaskTypeQA {
			continue
		}
		path := filepath.Join(filepath.Dir(graph.Spec.Dir), filepath.FromSlash(task.File))
		content, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("read Task %q for Consumer Declaration: %w", path, err)
		}
		content = proofOutsideFences(content)
		display := artifactDisplayPath(repoRoot, path)
		lines := markdownSectionLines(content, "Consumers")
		first := 0
		for first < len(lines) && strings.TrimSpace(lines[first]) == "" {
			first++
		}
		if first == len(lines) {
			appendConsumerGap(result, display, sectionLineContaining(content, "Consumers", "\x00"), "has no Consumer Declaration")
			continue
		}
		if strings.HasPrefix(strings.TrimSpace(lines[first]), "None.") {
			if !validNoneParagraph(lines[first:]) {
				appendConsumerGap(result, display, sectionLineContaining(content, "Consumers", "None."), "has a malformed Consumer Declaration; None. needs a reason in the same paragraph")
			}
			continue
		}

		// Fold indented continuation lines into the bullet while keeping its location.
		entries := make([]consumerEntry, 0)
		headingLine := sectionLineContaining(content, "Consumers", "\x00")
		for index := first; index < len(lines); index++ {
			text := strings.TrimSpace(lines[index])
			if text == "" {
				continue
			}
			entry := consumerEntry{text: text, line: headingLine + index + 1}
			for index+1 < len(lines) {
				next := lines[index+1]
				trimmed := strings.TrimSpace(next)
				if trimmed == "" || markdownBlockLine(trimmed) || !strings.HasPrefix(next, " ") && !strings.HasPrefix(next, "\t") {
					break
				}
				index++
				entry.text += " " + trimmed
			}
			entries = append(entries, entry)
		}

		commands := append([]string(nil), task.Verification...)
		for _, entry := range entries {
			match := consumerCheckpointPattern.FindStringSubmatch(entry.text)
			if len(match) != 2 {
				continue
			}
			checkpoint, present := tasks[match[1]]
			if !present || checkpoint.ID == graph.QATaskID || checkpoint.Type == spec.TaskTypeQA ||
				!ancestors[checkpoint.ID][task.ID] || !graph.QADeclined && !ancestors[graph.QATaskID][checkpoint.ID] {
				appendConsumerGap(result, display, entry.line, "names checkpoint "+match[1]+" outside the path from the Task to the QA gate")
				continue
			}
			commands = append(commands, checkpoint.Verification...)
		}
		for _, entry := range entries {
			if consumerCheckpointPattern.MatchString(entry.text) {
				continue
			}
			match := consumerEntryPattern.FindStringSubmatch(entry.text)
			if len(match) != 3 || strings.TrimSpace(match[1]) == "" {
				appendConsumerGap(result, display, entry.line, "has a malformed Consumer Declaration entry")
				continue
			}
			for _, suiteMatch := range consumerSuitePattern.FindAllStringSubmatch(match[2], -1) {
				suite := suiteMatch[1]
				if strings.TrimSpace(suite) == "" {
					appendConsumerGap(result, display, entry.line, "has a malformed Consumer Declaration entry with a blank suite")
					continue
				}
				proved := false
				for _, command := range commands {
					if strings.Contains(command, suite) {
						proved = true
						break
					}
				}
				if !proved {
					appendConsumerGap(result, display, entry.line, fmt.Sprintf("names consumer suite %q that no Verification command of the Task or its valid checkpoints contains", suite))
				}
			}
		}
	}
	return nil
}

func appendConsumerGap(result *Result, display string, line int, reason string) {
	result.Findings = append(result.Findings, Finding{
		Code: CodeConsumersUnproven, Severity: SeverityGap,
		Summary: display + " " + reason,
		Where:   []Location{{Path: display, Line: line}},
		Fix:     "Add or repair ## Consumers with None. and a reason, or literal suites in the Task's Verification or a valid checkpoint's Verification, following the write-tasks rule \"Consumers of a shared contract\".",
	})
}
