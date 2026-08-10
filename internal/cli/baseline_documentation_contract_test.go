package cli

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// Suite: Baseline public documentation contract
// Invariant: every public Baseline recipe matches the shipped parser and strict document schemas.
// Boundary IN: CLI help, public guides, README, and canonical public skill recipes.
// Boundary OUT: real repository adoption journeys and Python-runtime removal.

func TestBaselineExamplesParse(t *testing.T) {
	t.Parallel()
	root := baselineDocumentationRepoRoot()
	paths := []string{
		filepath.Join(root, "README.md"),
		filepath.Join(root, "docs", "user-guide", "context-driven-development.md"),
		filepath.Join(root, ".agents", "skills", "setup-context-driven", "SKILL.md"),
		filepath.Join(root, ".agents", "skills", "roundfix", "SKILL.md"),
	}

	temporaryRepo := t.TempDir()
	if err := os.Mkdir(filepath.Join(temporaryRepo, ".git"), 0o755); err != nil {
		t.Fatalf("create documentation-test Git marker: %v", err)
	}
	exampleCount := 0
	profileCounter := 0
	for _, path := range paths {
		content := readBaselineDocumentation(t, path)
		for _, command := range baselineBashExamples(content) {
			exampleCount++
			args := strings.Fields(command)
			if len(args) < 2 || args[0] != "roundfix" || args[1] != "baseline" {
				t.Fatalf("%s: invalid Baseline command example %q", path, command)
			}
			if len(args) >= 4 && args[2] == "profile" && args[3] == "init" {
				profileCounter++
				args = normalizeBaselineProfileInitExample(args, profileCounter)
			}
			if err := parsePublishedBaselineExample(args[2:], temporaryRepo); err != nil {
				t.Fatalf("%s: command example %q does not parse: %v", path, command, err)
			}
		}
	}
	if exampleCount < 20 {
		t.Fatalf("parsed %d Baseline examples, want at least 20 public recipes", exampleCount)
	}
}

func parsePublishedBaselineExample(args []string, workDir string) error {
	switch {
	case len(args) == 0:
		_, err := parseBaselineHumanCommand(nil)
		return err
	case args[0] == "update":
		_, err := parseBaselineUpdateCommand(args[1:])
		return err
	case args[0] == "plan":
		_, err := parseBaselinePlanCommand(args[1:])
		return err
	case args[0] == "apply":
		_, err := parseBaselineApplyCommand(args[1:])
		return err
	case len(args) >= 2 && args[0] == "profile" && args[1] == "init":
		var stdout, stderr bytes.Buffer
		if code := runBaselineProfileInitCommand(args[2:], &stdout, &stderr, commandEnvironment{workDir: workDir}); code != exitOK {
			return fmt.Errorf("exit %d: %s", code, stderr.String())
		}
		return nil
	case len(args) >= 2 && args[0] == "profile" && args[1] == "show":
		_, _, err := parseBaselineProfileTargetFormat(args[2:], true)
		return err
	case len(args) >= 2 && args[0] == "profile" && args[1] == "validate":
		_, _, err := parseBaselineProfileTargetFormat(args[2:], false)
		return err
	case len(args) >= 2 && args[0] == "skills" && args[1] == "restore":
		_, err := parseBaselineSkillsRestoreCommand(args[2:])
		return err
	case len(args) >= 2 && args[0] == "assets" && args[1] == "sync":
		_, err := parseBaselineAssetsSyncCommand(args[2:])
		return err
	default:
		_, err := parseBaselineHumanCommand(args)
		return err
	}
}

func normalizeBaselineProfileInitExample(args []string, sequence int) []string {
	normalized := append([]string(nil), args...)
	for index := 0; index+1 < len(normalized); index++ {
		switch normalized[index] {
		case "--id":
			normalized[index+1] = fmt.Sprintf("documentation-profile-%d", sequence)
		case "--from":
			if strings.HasPrefix(normalized[index+1], "<") {
				normalized[index+1] = "go-cli-tui"
			}
		}
	}
	return normalized
}

func baselineBashExamples(content string) []string {
	var examples []string
	inBash := false
	for _, line := range strings.Split(content, "\n") {
		trimmed := strings.TrimSpace(line)
		if trimmed == "```bash" {
			inBash = true
			continue
		}
		if trimmed == "```" {
			inBash = false
			continue
		}
		if inBash && (trimmed == "roundfix baseline" || strings.HasPrefix(trimmed, "roundfix baseline ")) {
			examples = append(examples, trimmed)
		}
	}
	return examples
}

func baselineDocumentationRepoRoot() string {
	root, err := filepath.Abs(filepath.Join("..", ".."))
	if err != nil {
		return filepath.Clean(filepath.Join("..", ".."))
	}
	return root
}

func readBaselineDocumentation(t *testing.T, path string) string {
	t.Helper()
	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	return string(content)
}

func assertBaselineDocumentationContains(t *testing.T, label string, content string, snippets []string) {
	t.Helper()
	for _, snippet := range snippets {
		if !strings.Contains(content, snippet) {
			t.Fatalf("%s missing %q", label, snippet)
		}
	}
}

func betweenBaselineDocumentation(t *testing.T, content, start, end string) string {
	t.Helper()
	startIndex := strings.Index(content, start)
	if startIndex < 0 {
		t.Fatalf("documentation marker %q is missing", start)
	}
	startIndex += len(start)
	endIndex := strings.Index(content[startIndex:], end)
	if endIndex < 0 {
		t.Fatalf("documentation marker %q is missing", end)
	}
	return content[startIndex : startIndex+endIndex]
}

func markdownFenceContent(t *testing.T, content string) string {
	t.Helper()
	const opening = "```markdown"
	start := strings.Index(content, opening)
	if start < 0 {
		t.Fatalf("markdown fence %q is missing", opening)
	}
	start += len(opening)
	content = content[start:]
	end := strings.Index(content, "```")
	if end < 0 {
		t.Fatal("closing markdown fence is missing")
	}
	return strings.TrimSpace(content[:end])
}

// TestProjectConstraintPlanningExamplesParse keeps the parser-level half of
// the project-constraint documentation contract beside the parsers it needs:
// the string contracts on the same two documents live in
// internal/docscontract.
func TestProjectConstraintPlanningExamplesParse(t *testing.T) {
	t.Parallel()
	root := baselineDocumentationRepoRoot()
	guide := readBaselineDocumentation(t, filepath.Join(
		root, "docs", "user-guide", "context-driven-development.md",
	))
	skill := readBaselineDocumentation(t, filepath.Join(
		root, ".agents", "skills", "setup-context-driven", "SKILL.md",
	))
	parsedPlanningExamples := 0
	for _, content := range []string{guide, skill} {
		for _, command := range baselineBashExamples(content) {
			if !strings.Contains(command, "baseline plan") ||
				!strings.Contains(command, "--decision-file") {
				continue
			}
			args := strings.Fields(command)
			if err := parsePublishedBaselineExample(args[2:], ""); err != nil {
				t.Fatalf("project-decision command %q does not parse: %v", command, err)
			}
			parsedPlanningExamples++
		}
	}
	if parsedPlanningExamples < 2 {
		t.Fatalf("parsed project-decision planning examples = %d, want at least 2", parsedPlanningExamples)
	}
}
