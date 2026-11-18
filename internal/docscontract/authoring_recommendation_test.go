//go:build docscontract

//verify:always

package docscontract

import (
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"roundfix/internal/config"
)

func TestTheRecommendedAuthoringDeclarationMatchesThisRepository(t *testing.T) {
	t.Parallel()

	root := baselineDocumentationRepoRoot()
	t.Run("repository_config", func(t *testing.T) {
		loaded, err := config.Load(config.LoadOptions{
			HomeDir:                  t.TempDir(),
			WorkDir:                  root,
			Stderr:                   io.Discard,
			WorktreeLocationOptional: true,
		})
		if err != nil {
			t.Fatalf("load repository config: %v", err)
		}

		want := config.RecommendedAuthoringDeclaration()
		got := loaded.Config.Skills.Authoring.Claude
		if got == nil {
			t.Fatal("skills.authoring.claude is missing")
		}
		var researchGot *config.AuthoringEntry
		if got.Research != nil {
			researchGot = &got.Research.AuthoringEntry
		}
		entries := []struct {
			key  string
			got  *config.AuthoringEntry
			want *config.AuthoringEntry
		}{
			{key: "idea", got: got.Idea, want: want.Idea},
			{key: "prd", got: got.PRD, want: want.PRD},
			{key: "techspec", got: got.TechSpec, want: want.TechSpec},
			{key: "tasks", got: got.Tasks, want: want.Tasks},
			{key: "research", got: researchGot, want: &want.Research.AuthoringEntry},
		}
		for _, entry := range entries {
			if entry.got == nil {
				t.Errorf("skills.authoring.claude.%s is missing", entry.key)
				continue
			}
			if entry.got.Model != entry.want.Model {
				t.Errorf("skills.authoring.claude.%s.model = %q, want %q", entry.key, entry.got.Model, entry.want.Model)
			}
			if entry.got.ReasoningEffort != entry.want.ReasoningEffort {
				t.Errorf("skills.authoring.claude.%s.reasoning_effort = %q, want %q", entry.key, entry.got.ReasoningEffort, entry.want.ReasoningEffort)
			}
		}
	})

	t.Run("configuration_guide", func(t *testing.T) {
		path := filepath.Join(root, "docs", "user-guide", "configuration.md")
		content, err := os.ReadFile(path)
		if err != nil {
			t.Fatalf("read %s: %v", path, err)
		}
		got := firstYAMLBlockUnderHeading(string(content), "## Authoring models")
		if got == "" {
			t.Fatalf("%s has no fenced yaml block under %q", path, "## Authoring models")
		}
		if want := config.RecommendedAuthoringYAML(); got != want {
			t.Fatalf("recommended declaration in %s = %q, want %q", path, got, want)
		}
	})
}

func firstYAMLBlockUnderHeading(content, heading string) string {
	start := strings.Index(content, heading)
	if start < 0 {
		return ""
	}
	section := content[start+len(heading):]
	if next := strings.Index(section, "\n## "); next >= 0 {
		section = section[:next]
	}
	lines := strings.Split(section, "\n")
	for i, line := range lines {
		if strings.TrimSpace(line) != "```yaml" {
			continue
		}
		for end := i + 1; end < len(lines); end++ {
			if strings.TrimSpace(lines[end]) == "```" {
				return strings.Join(lines[i+1:end], "\n") + "\n"
			}
		}
		return ""
	}
	return ""
}
