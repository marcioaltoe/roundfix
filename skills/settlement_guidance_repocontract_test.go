// Suite: canonical QA settlement guidance
// Invariant: the QA settlement section is one identical six-row contract in every canonical workflow skill.
// Boundary IN: the three canonical skill files under .agents/skills.
// Boundary OUT: distributed mirror regeneration, which is checked by make skills-sync and its caller.

package skills

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestSettlementGuidanceIsOneTable(t *testing.T) {
	paths := []string{
		filepath.Join("..", ".agents", "skills", "qa-gate", "SKILL.md"),
		filepath.Join("..", ".agents", "skills", "archive-spec", "SKILL.md"),
		filepath.Join("..", ".agents", "skills", "roundfix", "SKILL.md"),
	}
	rows := []string{
		"| `pass` |",
		"| qualifying declared `partial` |",
		"| `environment-blocked` |",
		"| `failed` |",
		"| `missing` |",
		"| `override` |",
	}

	var want string
	for _, path := range paths {
		content, err := os.ReadFile(path)
		if err != nil {
			t.Fatalf("read %s: %v", path, err)
		}
		section, ok := markdownSection(string(content), "### QA settlement")
		if !ok {
			t.Fatalf("%s is missing the QA settlement section", path)
		}
		if want == "" {
			want = section
		} else if section != want {
			t.Fatalf("QA settlement section differs in %s\nwant:\n%s\n got:\n%s", path, want, section)
		}
		for _, row := range rows {
			if !strings.Contains(section, row) {
				t.Errorf("%s is missing QA settlement row %q", path, row)
			}
		}
	}
}

func TestTaskAuthoringGuidanceNamesDeclarations(t *testing.T) {
	paths := []string{
		filepath.Join("..", ".agents", "skills", "write-tasks", "SKILL.md"),
		filepath.Join("..", ".agents", "skills", "write-tasks", "references", "task-template.md"),
	}
	required := []string{
		"verification: independent",
		"precondition_repairs",
		"SC-ORDINAL-CLAIMED",
		"temporal prerequisite",
		"never invents release authority",
		"property-shaped acceptance",
		"test seam",
		"Spec commits narrow",
		"newly required test class",
		"repository gate",
	}

	for _, path := range paths {
		content, err := os.ReadFile(path)
		if err != nil {
			t.Fatalf("read %s: %v", path, err)
		}
		for _, phrase := range required {
			if !strings.Contains(string(content), phrase) {
				t.Errorf("%s is missing authoring guidance %q", path, phrase)
			}
		}
	}
}

func markdownSection(content, heading string) (string, bool) {
	lines := strings.SplitAfter(content, "\n")
	start := -1
	for i, line := range lines {
		if strings.TrimSuffix(line, "\n") == heading {
			start = i
			break
		}
	}
	if start == -1 {
		return "", false
	}

	end := len(lines)
	for i := start + 1; i < len(lines); i++ {
		line := strings.TrimSuffix(lines[i], "\n")
		if strings.HasPrefix(line, "### ") || strings.HasPrefix(line, "## ") {
			end = i
			break
		}
	}
	return strings.Join(lines[start:end], ""), true
}
