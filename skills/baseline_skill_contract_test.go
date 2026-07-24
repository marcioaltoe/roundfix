package skills

import (
	"bytes"
	"os"
	"path/filepath"
	"reflect"
	"sort"
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

func TestNoPythonBaselineRuntime(t *testing.T) {
	repoRoot := filepath.Clean(filepath.Join(".."))
	for _, root := range []string{
		filepath.Join(repoRoot, ".agents", "skills", "setup-context-driven"),
		filepath.Join(repoRoot, "skills", "setup-context-driven"),
	} {
		var files []string
		if err := filepath.WalkDir(root, func(filePath string, entry os.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if entry.IsDir() {
				return nil
			}
			relative, err := filepath.Rel(root, filePath)
			if err != nil {
				return err
			}
			files = append(files, filepath.ToSlash(relative))
			return nil
		}); err != nil {
			t.Fatalf("walk setup skill %s: %v", root, err)
		}
		sort.Strings(files)
		if !reflect.DeepEqual(files, []string{"SKILL.md"}) {
			t.Fatalf("setup skill %s ships non-guidance files: %v", root, files)
		}
	}

	makefile := string(readBaselineSkillContractFile(t, filepath.Join(repoRoot, "Makefile")))
	for _, forbidden := range []string{
		"setup-context-check",
		"PYTHONDONTWRITEBYTECODE",
		"python3",
	} {
		if strings.Contains(makefile, forbidden) {
			t.Fatalf("post-cutover Makefile invokes retired runtime marker %q", forbidden)
		}
	}

	for _, relative := range []string{
		"README.md",
		filepath.Join("docs", "user-guide", "commands.md"),
		filepath.Join("docs", "user-guide", "context-driven-development.md"),
		filepath.Join(".agents", "skills", "setup-context-driven", "SKILL.md"),
		filepath.Join(".agents", "skills", "roundfix", "SKILL.md"),
	} {
		content := string(readBaselineSkillContractFile(t, filepath.Join(repoRoot, relative)))
		for _, forbidden := range []string{
			"context_" + "setup.py",
			"context_" + "baseline.py",
			"python3",
			"Python fallback",
		} {
			if strings.Contains(content, forbidden) {
				t.Fatalf("%s references retired Baseline runtime marker %q", relative, forbidden)
			}
		}
	}
}

func TestThinSetupSkill(t *testing.T) {
	files, err := Files()
	if err != nil {
		t.Fatalf("read embedded skills: %v", err)
	}
	var setupFiles []string
	for _, file := range files {
		if file.Skill == "setup-context-driven" {
			setupFiles = append(setupFiles, file.Path)
		}
	}
	sort.Strings(setupFiles)
	if !reflect.DeepEqual(setupFiles, []string{"setup-context-driven/SKILL.md"}) {
		t.Fatalf("embedded setup skill files = %v", setupFiles)
	}
	for _, diagnostic := range Check() {
		if strings.HasPrefix(diagnostic.Path, "setup-context-driven/") {
			t.Fatalf("thin setup skill diagnostic: %s: %s", diagnostic.Path, diagnostic.Message)
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
