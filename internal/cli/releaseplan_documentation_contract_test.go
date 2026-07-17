package cli

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// Suite: Release Plan documentation contract
// Invariant: public release guidance starts with the same read-only Release Plan Command and approval boundary.
// Boundary IN: root help, command help, durable docs, and owned Roundfix skill copies.
// Boundary OUT: publication workflow execution, GitHub, npm, and release assets.

func TestReleasePlanDocumentationContract(t *testing.T) {
	var rootStdout, rootStderr bytes.Buffer
	if code := Run([]string{"--help"}, &rootStdout, &rootStderr); code != exitOK {
		t.Fatalf("root help exit = %d, want 0 stderr=%q", code, rootStderr.String())
	}
	if rootStderr.String() != "" {
		t.Fatalf("root help stderr = %q, want empty", rootStderr.String())
	}
	assertReleasePlanDocumentationContains(t, "root help", rootStdout.String(), []string{
		"roundfix release plan [--from <tag>] [--to <revision>] [--format <text|json>]",
		"release    Plan the next release version without mutating repository or release state",
	})

	var commandStdout, commandStderr bytes.Buffer
	if code := Run([]string{"release", "plan", "--help"}, &commandStdout, &commandStderr); code != exitOK {
		t.Fatalf("release plan help exit = %d, want 0 stderr=%q", code, commandStderr.String())
	}
	if commandStderr.String() != "" {
		t.Fatalf("release plan help stderr = %q, want empty", commandStderr.String())
	}
	assertReleasePlanDocumentationContains(t, "release plan help", commandStdout.String(), []string{
		"roundfix release plan [--from <tag>] [--to <revision>] [--impact <none|patch|minor|major> --reason <text>] [--format <text|json>]",
		"ready",
		"approval_required",
		"manual_classification_required",
		"no_release",
		"--from",
		"--to",
		"--impact",
		"--reason",
		"--format",
		"creates no Run",
		"never mutates",
	})

	repoRoot := releasePlanDocumentationRepoRoot()
	docs := []struct {
		name     string
		path     string
		snippets []string
	}{
		{
			name: "release runbook",
			path: filepath.Join(repoRoot, "docs", "user-guide", "release-runbook.md"),
			snippets: []string{
				"roundfix release plan",
				"before changelog edits",
				"version-file edits",
				"GitHub Release creation",
				"generic release request authorizes",
				"patch proposal",
				"minor, major, or version-zero breaking",
				"--impact <none|patch|minor|major> --reason <text>",
				"does not approve",
				"tag-triggered",
				"publication workflow",
			},
		},
		{
			name: "usage command index",
			path: filepath.Join(repoRoot, "docs", "user-guide", "usage.md"),
			snippets: []string{
				"Release planning before publication",
				"roundfix release plan",
				"generic release",
				"authorizes only a conclusive patch plan",
				"version-zero",
				"breaking plans require explicit human approval",
				"--impact <none|patch|minor|major> --reason <text>",
				"does not approve",
				"| `release plan` | Classify committed release changes without mutating release state |",
			},
		},
		{
			name: "root Agent pointer",
			path: filepath.Join(repoRoot, "AGENTS.md"),
			snippets: []string{
				"HARD RULE — release planning",
				"roundfix release plan",
				"before changelog, version, tag, push, package, asset",
				"conclusive patch plan",
				"minor, major, version-zero breaking",
				"manual",
				"classification outcomes require the decisions",
			},
		},
		{
			name: "canonical Roundfix skill",
			path: filepath.Join(repoRoot, ".agents", "skills", "roundfix", "SKILL.md"),
			snippets: []string{
				"## Release planning",
				"roundfix release plan",
				"before changelog edits",
				"version-file edits",
				"GitHub Release creation",
				"generic release request authorizes only a conclusive patch plan",
				"minor,",
				"major, or version-zero breaking",
				"--impact <none|patch|minor|major> --reason <text>",
				"does not approve",
			},
		},
		{
			name: "embedded Roundfix skill",
			path: filepath.Join(repoRoot, "skills", "roundfix", "SKILL.md"),
			snippets: []string{
				"## Release planning",
				"roundfix release plan",
				"generic release request authorizes only a conclusive patch plan",
				"--impact <none|patch|minor|major> --reason <text>",
			},
		},
	}
	for _, doc := range docs {
		t.Run(doc.name, func(t *testing.T) {
			content := readReleasePlanDocumentation(t, doc.path)
			assertReleasePlanDocumentationContains(t, doc.name, content, doc.snippets)
		})
	}

	canonical := readReleasePlanDocumentation(t, filepath.Join(repoRoot, ".agents", "skills", "roundfix", "SKILL.md"))
	embedded := readReleasePlanDocumentation(t, filepath.Join(repoRoot, "skills", "roundfix", "SKILL.md"))
	if canonical != embedded {
		t.Fatal("embedded Roundfix skill differs from canonical .agents/skills/roundfix/SKILL.md; run make skills-sync")
	}
}

func releasePlanDocumentationRepoRoot() string {
	return filepath.Clean(filepath.Join("..", ".."))
}

func readReleasePlanDocumentation(t *testing.T, path string) string {
	t.Helper()
	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	return string(content)
}

func assertReleasePlanDocumentationContains(t *testing.T, label string, content string, snippets []string) {
	t.Helper()
	for _, snippet := range snippets {
		if !strings.Contains(content, snippet) {
			t.Fatalf("%s missing %q", label, snippet)
		}
	}
}
