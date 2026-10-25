package cli

import (
	"bytes"
	"fmt"
	"strings"
	"testing"
)

// Suite: Baseline public documentation contract
// Invariant: every public Baseline recipe matches the shipped parser and strict document schemas.
// Boundary IN: CLI help, public guides, README, and canonical public skill recipes.
// Boundary OUT: real repository adoption journeys and Python-runtime removal.

func TestBaselineSkillsReconcileExampleParses(t *testing.T) {
	t.Parallel()
	err := parsePublishedBaselineExample([]string{
		"skills", "reconcile",
		"--profile", "go-cli-tui",
		"--source", "example/skills",
		"--revision",
	}, "")
	if err == nil || !strings.Contains(err.Error(), "invalid baseline skills reconcile arguments") {
		t.Fatalf("malformed reconcile example error = %v", err)
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
	case len(args) >= 2 && args[0] == "skills" && args[1] == "reconcile":
		_, err := parseBaselineSkillsReconcileCommand(args[2:])
		return err
	case len(args) >= 2 && args[0] == "assets" && args[1] == "sync":
		_, err := parseBaselineAssetsSyncCommand(args[2:])
		return err
	default:
		_, err := parseBaselineHumanCommand(args)
		return err
	}
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
