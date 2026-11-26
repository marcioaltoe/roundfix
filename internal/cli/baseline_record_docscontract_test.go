//go:build docscontract

// Suite: repository record contracts.
// Invariant: records and repository rules remain current.
// Boundary IN: live repository records and source declarations.
// Boundary OUT: product execution and publication.
//
//verify:always
package cli

import (
	"bytes"
	"os"
	"path/filepath"
	"reflect"
	"roundfix/internal/baseline"
	"strings"
	"testing"
)

func TestBaselineUpdateSurfaceTranscriptThisRepositoryIsCurrent(t *testing.T) {
	t.Parallel()
	root, err := filepath.Abs(filepath.Join("..", ".."))
	if err != nil {
		t.Fatal(err)
	}
	setCommandEnvironmentForTest(t, t.TempDir(), root)
	var stdout, stderr bytes.Buffer
	code := runCLI(t, []string{"baseline", "update", "--repo", root, "--no-skills", "--no-history", "--format", "text"}, &stdout, &stderr)
	if code != exitOK || stderr.Len() != 0 {
		t.Fatalf("exit=%d stderr=%s stdout=%s", code, stderr.String(), stdout.String())
	}
	catalog, err := baseline.LoadEmbeddedCatalog()
	if err != nil {
		t.Fatal(err)
	}
	want := []string{"Baseline update: current", "Result: the repository already matches the current Baseline catalog", "Prior catalog: " + catalog.Digest(), "Current catalog: " + catalog.Digest(), "File changes: 0", "Misplaced repository rules: 0", "Retention evidence: 0", "Declared decisions: 0"}
	lines := strings.Split(strings.TrimSuffix(stdout.String(), "\n"), "\n")
	if len(lines) <= len(want) || !reflect.DeepEqual(lines[:len(want)], want) {
		t.Fatalf("transcript first eight lines=%q, want %q", lines, want)
	}
}

func TestBaselineUpdateSurfaceTranscriptKeepsTheClaudeInstructionLink(t *testing.T) {
	t.Parallel()
	root, err := filepath.Abs(filepath.Join("..", ".."))
	if err != nil {
		t.Fatal(err)
	}
	setCommandEnvironmentForTest(t, t.TempDir(), root)
	assertLink := func() {
		t.Helper()
		target, err := os.Readlink(filepath.Join(root, "CLAUDE.md"))
		if err != nil || target != "AGENTS.md" {
			t.Fatalf("CLAUDE.md target=%q error=%v", target, err)
		}
	}
	assertLink()
	var stdout, stderr bytes.Buffer
	code := runCLI(t, []string{"baseline", "update", "--repo", root, "--no-skills", "--no-history", "--format", "text"}, &stdout, &stderr)
	if code != exitOK || stderr.Len() != 0 {
		t.Fatalf("exit=%d stderr=%s stdout=%s", code, stderr.String(), stdout.String())
	}
	catalog, err := baseline.LoadEmbeddedCatalog()
	if err != nil {
		t.Fatal(err)
	}
	want := []string{"Baseline update: current", "Result: the repository already matches the current Baseline catalog", "Prior catalog: " + catalog.Digest(), "Current catalog: " + catalog.Digest(), "File changes: 0", "Misplaced repository rules: 0", "Retention evidence: 0", "Declared decisions: 0"}
	lines := strings.Split(strings.TrimSuffix(stdout.String(), "\n"), "\n")
	if len(lines) <= len(want) || !reflect.DeepEqual(lines[:len(want)], want) {
		t.Fatalf("transcript=%q, want prefix=%q followed by remaining output", lines, want)
	}
	for _, line := range lines {
		if strings.Contains(line, "baseline.root.claude-link") || strings.Contains(line, "CLAUDE.md") {
			t.Fatalf("unexpected link plan or warning: %s", line)
		}
	}
	assertLink()
}
