package skills_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"roundfix/internal/speccheck"
)

func TestTheGlossaryGuidesStateTheDefinitionForms(t *testing.T) {
	t.Parallel()

	root, err := filepath.Abs("..")
	if err != nil {
		t.Fatalf("resolve repository root: %v", err)
	}

	forms := strings.Join(strings.Fields(speccheck.GlossaryDefinitionForms), " ")
	command := speccheck.GlossaryVerificationCommand("<Term>", "CONTEXT.md")
	paths := []string{
		".agents/skills/write-tasks/SKILL.md",
		"skills/write-tasks/SKILL.md",
		".agents/skills/write-prd/references/glossary.md",
		"skills/write-prd/references/glossary.md",
		".agents/skills/roundfix/references/spec.md",
		"skills/roundfix/references/spec.md",
	}

	for _, path := range paths {
		path := path
		t.Run(path, func(t *testing.T) {
			content, err := os.ReadFile(filepath.Join(root, filepath.FromSlash(path)))
			if err != nil {
				t.Fatalf("read %s: %v", path, err)
			}
			normalized := strings.Join(strings.Fields(string(content)), " ")
			if !strings.Contains(normalized, forms) {
				t.Fatalf("%s does not state the accepted Glossary Definition forms", path)
			}
			if strings.Contains(path, "write-tasks/") {
				if !strings.Contains(string(content), command) {
					t.Fatalf("%s does not show the rendered glossary Verification command", path)
				}
				if !strings.Contains(string(content), "reads lines, not whitespace-normalized text") {
					t.Fatalf("%s does not say the Verification reads lines", path)
				}
				if !strings.Contains(string(content), "never more permissive than the check") {
					t.Fatalf("%s does not constrain the Verification to the check", path)
				}
			}
			if strings.Contains(path, "write-prd/") || strings.Contains(path, "roundfix/") {
				if !strings.Contains(normalized, "A heading, an indented or ordered list item and bold text inside a sentence define nothing.") {
					t.Fatalf("%s does not reject non-definition glossary text", path)
				}
			}
		})
	}
}
