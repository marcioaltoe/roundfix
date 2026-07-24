package cli

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"roundfix/internal/baseline"
)

// Suite: Baseline public documentation contract
// Invariant: every public Baseline recipe matches the shipped parser and strict document schemas.
// Boundary IN: CLI help, public guides, README, and canonical public skill recipes.
// Boundary OUT: real repository adoption journeys and Python-runtime removal.

func TestBaselineDocumentationContract(t *testing.T) {
	helpCases := []struct {
		command  string
		args     []string
		snippets []string
	}{
		{
			command: "baseline",
			args:    []string{"baseline", "--help"},
			snippets: []string{
				"roundfix baseline [--repo <path>] [--format <text|json>]",
				"roundfix baseline plan --profile <id>",
				"roundfix baseline apply --plan <file> --confirm-plan <digest>",
				"roundfix baseline profile init --id <id>",
				"roundfix baseline skills restore --profile <id>",
				"roundfix baseline assets sync --source-dir <path>",
				"explicit confirmation of the displayed Plan Digest",
			},
		},
		{
			command: "baseline plan",
			args:    []string{"baseline", "plan", "--help"},
			snippets: []string{
				"roundfix/baseline-plan/v1",
				"roundfix/baseline-result/v1",
				"--decision-file",
				"never prompts",
				"never",
				"uses the network",
			},
		},
		{
			command: "baseline apply",
			args:    []string{"baseline", "apply", "--help"},
			snippets: []string{
				"--confirm-plan",
				"recoverable",
				"transaction",
				"Repository formatter and Verification commands are reported as",
				"never run",
			},
		},
		{
			command: "baseline profile",
			args:    []string{"baseline", "profile", "--help"},
			snippets: []string{
				"profile init",
				"profile show",
				"profile validate",
			},
		},
		{
			command: "baseline skills restore",
			args:    []string{"baseline", "skills", "restore", "--help"},
			snippets: []string{
				"complete preimage",
				"--source-dir",
				"--confirm-plan",
			},
		},
		{
			command: "baseline assets sync",
			args:    []string{"baseline", "assets", "sync", "--help"},
			snippets: []string{
				"--source-dir",
				"--check",
				"never installs skills",
			},
		},
	}
	for _, test := range helpCases {
		t.Run(test.command, func(t *testing.T) {
			var stdout, stderr bytes.Buffer
			if code := Run(test.args, &stdout, &stderr); code != exitOK {
				t.Fatalf("help exit = %d, want 0 stderr=%q", code, stderr.String())
			}
			if stderr.Len() != 0 {
				t.Fatalf("help stderr = %q, want empty", stderr.String())
			}
			assertBaselineDocumentationContains(t, test.command+" help", stdout.String(), test.snippets)
		})
	}

	root := baselineDocumentationRepoRoot()
	docs := []struct {
		name     string
		path     string
		snippets []string
	}{
		{
			name: "Context-Driven user guide",
			path: filepath.Join(root, "docs", "user-guide", "context-driven-development.md"),
			snippets: []string{
				"### First adoption",
				"Greenfield",
				"Preservation",
				"### Update, profile change, and rejected plans",
				"### Automation",
				"roundfix/baseline-plan/v1",
				"roundfix/baseline-result/v1",
				"### Profiles",
				"### Repository Skill Set restoration",
				"### Canonical asset synchronization",
				"### Recovery and troubleshooting",
				"Interrupted transaction",
				"Incomplete rollback",
				"### Security and execution limits",
				"### Migrate from the script-backed setup skill",
				"profile expectation",
				"repository command",
				"recommendation",
			},
		},
		{
			name: "command reference",
			path: filepath.Join(root, "docs", "user-guide", "commands.md"),
			snippets: []string{
				"### baseline",
				"roundfix baseline plan --profile <id>",
				"roundfix baseline apply --plan <file> --confirm-plan <digest>",
				"stdout",
				"stderr",
				"incomplete-rollback",
			},
		},
		{
			name: "README",
			path: filepath.Join(root, "README.md"),
			snippets: []string{
				"roundfix baseline",
				"complete Baseline adoption",
				"migration",
				"recovery",
			},
		},
		{
			name: "canonical setup skill",
			path: filepath.Join(root, ".agents", "skills", "setup-context-driven", "SKILL.md"),
			snippets: []string{
				"public `roundfix baseline` command family",
				"only runtime authority",
				"roundfix baseline plan",
				"roundfix baseline apply",
				"## Recovery",
			},
		},
		{
			name: "canonical Roundfix skill",
			path: filepath.Join(root, ".agents", "skills", "roundfix", "SKILL.md"),
			snippets: []string{
				"## Context-Driven Baseline",
				"roundfix baseline plan",
				"roundfix baseline apply",
				"no partial plan",
			},
		},
	}
	for _, doc := range docs {
		t.Run(doc.name, func(t *testing.T) {
			content := readBaselineDocumentation(t, doc.path)
			assertBaselineDocumentationContains(t, doc.name, content, doc.snippets)
		})
	}

	guide := readBaselineDocumentation(
		t,
		filepath.Join(root, "docs", "user-guide", "context-driven-development.md"),
	)
	baselineSection := betweenBaselineDocumentation(
		t,
		guide,
		"## Adopt or update the Context-Driven Baseline",
		"## How Roundfix executes it",
	)
	for _, forbidden := range []string{"/Users/", "/home/", `C:\`, "context_setup.py"} {
		if strings.Contains(baselineSection, forbidden) {
			t.Fatalf("Baseline user guide contains environment-specific or retired runtime text %q", forbidden)
		}
	}
}

func TestBaselineExamplesParse(t *testing.T) {
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
	originalWorkingDirectory, err := os.Getwd()
	if err != nil {
		t.Fatalf("read working directory: %v", err)
	}
	if err := os.Chdir(temporaryRepo); err != nil {
		t.Fatalf("enter documentation-test repository: %v", err)
	}
	t.Cleanup(func() {
		if err := os.Chdir(originalWorkingDirectory); err != nil {
			t.Errorf("restore working directory: %v", err)
		}
	})

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
			if err := parsePublishedBaselineExample(args[2:]); err != nil {
				t.Fatalf("%s: command example %q does not parse: %v", path, command, err)
			}
		}
	}
	if exampleCount < 20 {
		t.Fatalf("parsed %d Baseline examples, want at least 20 public recipes", exampleCount)
	}
}

func TestBaselineDecisionExamples(t *testing.T) {
	path := filepath.Join(
		baselineDocumentationRepoRoot(),
		"docs",
		"user-guide",
		"context-driven-development.md",
	)
	content := readBaselineDocumentation(t, path)
	example := betweenBaselineDocumentation(
		t,
		content,
		"<!-- baseline-decision-document:start -->",
		"<!-- baseline-decision-document:end -->",
	)
	example = strings.TrimSpace(example)
	example = strings.TrimPrefix(example, "```json")
	example = strings.TrimSuffix(example, "```")
	example = strings.TrimSpace(example)

	document, err := baseline.ParseDecisionDocument([]byte(example), path)
	if err != nil {
		t.Fatalf("published Decision Document does not pass the strict parser: %v", err)
	}
	if document.SchemaVersion != baseline.DecisionDocumentSchemaVersion ||
		document.Version != baseline.DecisionDocumentVersion ||
		len(document.Decisions) != 9 {
		t.Fatalf("published Decision Document parsed unexpectedly: %#v", document)
	}

	invalid := strings.Replace(example, `"version": "0.0.1",`, `"version": "0.0.1", "unknown": true,`, 1)
	if _, err := baseline.ParseDecisionDocument([]byte(invalid), path+"#invalid"); err == nil {
		t.Fatal("strict parser accepted an unknown field in the published Decision Document shape")
	}
}

func parsePublishedBaselineExample(args []string) error {
	switch {
	case len(args) == 0:
		_, err := parseBaselineHumanCommand(nil)
		return err
	case args[0] == "plan":
		_, err := parseBaselinePlanCommand(args[1:])
		return err
	case args[0] == "apply":
		_, err := parseBaselineApplyCommand(args[1:])
		return err
	case len(args) >= 2 && args[0] == "profile" && args[1] == "init":
		var stdout, stderr bytes.Buffer
		if code := runBaselineProfileInitCommand(args[2:], &stdout, &stderr); code != exitOK {
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
