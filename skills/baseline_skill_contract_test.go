package skills

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// Suite: Baseline skill distribution contract
// Invariant: canonical and shipped Baseline guidance is identical and invokes only the public CLI.
// Boundary IN: setup-context-driven and Roundfix owned skill guidance.
// Boundary OUT: executable setup-runtime removal and Baseline CLI behavior.

func TestBaselineSkillContract(t *testing.T) {
	repoRoot := filepath.Clean(filepath.Join(".."))
	skillNames := []string{"setup-context-driven", "roundfix"}
	embeddedFiles, err := Files()
	if err != nil {
		t.Fatalf("read embedded skills: %v", err)
	}
	embeddedByPath := make(map[string][]byte, len(embeddedFiles))
	for _, file := range embeddedFiles {
		embeddedByPath[file.Path] = file.Data
	}

	for _, name := range skillNames {
		t.Run(name, func(t *testing.T) {
			canonicalPath := filepath.Join(repoRoot, ".agents", "skills", name, "SKILL.md")
			distributedPath := filepath.Join(repoRoot, "skills", name, "SKILL.md")
			canonical := readBaselineSkillContractFile(t, canonicalPath)
			distributed := readBaselineSkillContractFile(t, distributedPath)
			if !bytes.Equal(canonical, distributed) {
				t.Fatalf("%s canonical and distributed guidance differ; run make skills-sync", name)
			}
			embedded := embeddedByPath[filepath.ToSlash(filepath.Join(name, "SKILL.md"))]
			if !bytes.Equal(distributed, embedded) {
				t.Fatalf("%s distributed and embedded guidance differ", name)
			}
			body := baselineSkillBody(string(canonical))
			if strings.Contains(strings.ToLower(body), "python") ||
				strings.Contains(body, "context_setup.py") ||
				strings.Contains(body, "scripts/context") {
				t.Fatalf("%s skill body invokes a retired independent setup runtime", name)
			}
			for _, forbidden := range []string{"/Users/", "/home/", `C:\`} {
				if strings.Contains(body, forbidden) {
					t.Fatalf("%s skill body contains environment-specific path %q", name, forbidden)
				}
			}
		})
	}

	setup := string(readBaselineSkillContractFile(
		t,
		filepath.Join(repoRoot, ".agents", "skills", "setup-context-driven", "SKILL.md"),
	))
	for _, required := range []string{
		"only runtime authority",
		"roundfix baseline --repo . --format text",
		"roundfix baseline plan",
		"roundfix baseline apply",
		"roundfix/baseline-plan/v1",
		"roundfix/baseline-result/v1",
		"Repository Skill Set restoration",
		"Canonical asset synchronization",
		"behavioral fallback",
	} {
		if !strings.Contains(setup, required) {
			t.Fatalf("setup-context-driven skill missing %q", required)
		}
	}
}

func readBaselineSkillContractFile(t *testing.T, path string) []byte {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	return data
}

func baselineSkillBody(content string) string {
	const delimiter = "---"
	parts := strings.SplitN(content, delimiter, 3)
	if len(parts) != 3 {
		return content
	}
	return parts[2]
}
