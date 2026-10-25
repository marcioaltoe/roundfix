// Suite: setup-context-driven skill invocation contract
// Invariant: an Agent may start the setup skill while Plan Digest approval remains required
// Boundary: person-only skills keep their invocation lock

package skills

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestAnAgentMayStartTheSetupSkill(t *testing.T) {
	t.Parallel()

	for _, test := range []struct {
		name string
		path string
	}{
		{name: "canonical", path: filepath.Join("..", ".agents", "skills", "setup-context-driven", "SKILL.md")},
		{name: "shipped", path: filepath.Join("setup-context-driven", "SKILL.md")},
	} {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			data, err := os.ReadFile(test.path)
			if err != nil {
				t.Fatal(err)
			}
			frontmatter, body := setupSkillFrontmatterAndBody(string(data))
			if strings.Contains(frontmatter, "disable-model-invocation:") {
				t.Fatalf("%s still disables model invocation", test.path)
			}
			if strings.Index(body, "## Who starts this skill") >= strings.Index(body, "## Boundaries") {
				t.Fatalf("%s places the startup section after Boundaries", test.path)
			}
			normalized := strings.Join(strings.Fields(body), " ")
			for _, want := range []string{
				"Starting it authorizes no write.",
				"Ask the maintainer to approve the exact current `planDigest`",
			} {
				if !strings.Contains(normalized, want) {
					t.Errorf("%s is missing %q", test.path, want)
				}
			}
		})
	}

	t.Run("person-only-skills-keep-the-lock", func(t *testing.T) {
		t.Parallel()
		path := filepath.Join("..", ".agents", "skills", "implement-spec", "SKILL.md")
		data, err := os.ReadFile(path)
		if err != nil {
			t.Fatal(err)
		}
		frontmatter, _ := setupSkillFrontmatterAndBody(string(data))
		if !strings.Contains(frontmatter, "disable-model-invocation: true") {
			t.Fatalf("%s no longer keeps the person-only invocation lock", path)
		}
	})
}

func setupSkillFrontmatterAndBody(content string) (string, string) {
	parts := strings.SplitN(content, "---", 3)
	if len(parts) != 3 {
		return content, ""
	}
	return parts[1], parts[2]
}
