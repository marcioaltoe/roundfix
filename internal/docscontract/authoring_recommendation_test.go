//go:build docscontract

//verify:always

package docscontract

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"slices"
	"strings"
	"testing"

	"gopkg.in/yaml.v3"

	"roundfix/internal/config"
	"roundfix/internal/gittest"
)

func TestTheRecommendedAuthoringDeclarationMatchesThisRepository(t *testing.T) {
	t.Parallel()

	root := baselineDocumentationRepoRoot()

	load := func(t *testing.T, root string) *config.AuthoringHarness {
		t.Helper()
		loaded, err := config.Load(config.LoadOptions{
			HomeDir: t.TempDir(), WorkDir: root, Stderr: io.Discard, WorktreeLocationOptional: true,
		})
		if err != nil {
			t.Fatalf("load repository config: %v", err)
		}
		return loaded.Config.Skills.Authoring.Claude
	}
	t.Run("repository_config", func(t *testing.T) {
		if err := authoringRecommendationMismatch(load(t, root)); err != nil {
			t.Fatal(err)
		}
	})
	fixture := func(t *testing.T, phase string, fallback *config.AuthoringFallback) string {
		t.Helper()
		content, err := os.ReadFile(filepath.Join(root, ".roundfixrc.yml"))
		if err != nil {
			t.Fatal(err)
		}
		var document yaml.Node
		if err := yaml.Unmarshal(content, &document); err != nil {
			t.Fatal(err)
		}
		node := document.Content[0]
		for _, key := range []string{"skills", "authoring", "claude", phase} {
			var child *yaml.Node
			for i := 0; i < len(node.Content); i += 2 {
				if node.Content[i].Value == key {
					child = node.Content[i+1]
					break
				}
			}
			if child == nil {
				t.Fatalf("fixture lacks %s", key)
			}
			node = child
		}
		// Drop any existing fallback so the fixture remains useful after adoption.
		for i := 0; i < len(node.Content); i += 2 {
			if node.Content[i].Value == "fallback" {
				node.Content = append(node.Content[:i], node.Content[i+2:]...)
				break
			}
		}
		if fallback != nil {
			var value yaml.Node
			if err := value.Encode(fallback); err != nil {
				t.Fatal(err)
			}
			node.Content = append(node.Content, &yaml.Node{Kind: yaml.ScalarNode, Tag: "!!str", Value: "fallback"}, &value)
		}
		content, err = yaml.Marshal(&document)
		if err != nil {
			t.Fatal(err)
		}
		target := t.TempDir()
		gittest.MarkRepository(t, target)
		if err := os.WriteFile(filepath.Join(target, ".roundfixrc.yml"), content, 0o600); err != nil {
			t.Fatal(err)
		}
		return target
	}
	phases := []struct {
		key   string
		entry *config.AuthoringEntry
	}{
		{"idea", config.RecommendedAuthoringDeclaration().Idea},
		{"prd", config.RecommendedAuthoringDeclaration().PRD},
		{"techspec", config.RecommendedAuthoringDeclaration().TechSpec},
		{"tasks", config.RecommendedAuthoringDeclaration().Tasks},
	}
	t.Run("a_declared_fallback_that_differs_fails", func(t *testing.T) {
		for _, phase := range phases {
			for _, field := range []string{"model", "reasoning_effort"} {
				t.Run(phase.key+"/"+field, func(t *testing.T) {
					fallback := *phase.entry.Fallback
					if field == "model" {
						fallback.Model = "another-model"
					} else {
						fallback.ReasoningEffort = "low"
					}
					err := authoringRecommendationMismatch(load(t, fixture(t, phase.key, &fallback)))
					if err == nil || !strings.Contains(err.Error(), "skills.authoring.claude."+phase.key+".fallback") {
						t.Fatalf("contract error=%v; want the differing fallback", err)
					}
				})
			}
		}
	})
	t.Run("an_undeclared_fallback_passes", func(t *testing.T) {
		for _, phase := range phases {
			t.Run(phase.key, func(t *testing.T) {
				if err := authoringRecommendationMismatch(load(t, fixture(t, phase.key, nil))); err != nil {
					t.Fatal(err)
				}
			})
		}
	})
	t.Run("a_matching_declared_fallback_passes", func(t *testing.T) {
		for _, phase := range phases {
			t.Run(phase.key, func(t *testing.T) {
				if err := authoringRecommendationMismatch(load(t, fixture(t, phase.key, phase.entry.Fallback))); err != nil {
					t.Fatal(err)
				}
			})
		}
	})

	t.Run("a_research_entry_without_a_recommended_tool_fails", func(t *testing.T) {
		for _, tc := range []struct {
			name  string
			tools []string
		}{
			{"nil", nil},
			{"empty", []string{}},
			{"search_only", []string{"WebSearch"}},
			{"fetch_and_mcp", []string{"WebFetch", "mcp__kb__read"}},
		} {
			t.Run(tc.name, func(t *testing.T) {
				declaration := config.RecommendedAuthoringDeclaration()
				declaration.Research.Tools = tc.tools
				err := authoringRecommendationMismatch(&declaration)
				if err == nil || !strings.Contains(err.Error(), "skills.authoring.claude.research.tools") {
					t.Fatalf("contract error=%v; want the missing recommended research tool", err)
				}
			})
		}
		declaration := config.RecommendedAuthoringDeclaration()
		declaration.Research.Tools = append(declaration.Research.Tools, "mcp__kb__read")
		if err := authoringRecommendationMismatch(&declaration); err != nil {
			t.Fatalf("additional research tool: %v", err)
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

func authoringRecommendationMismatch(got *config.AuthoringHarness) error {
	if got == nil {
		return fmt.Errorf("skills.authoring.claude is missing")
	}
	want := config.RecommendedAuthoringDeclaration()
	var researchGot *config.AuthoringEntry
	if got.Research != nil {
		researchGot = &got.Research.AuthoringEntry
	}
	entries := []struct {
		key       string
		got, want *config.AuthoringEntry
	}{
		{"idea", got.Idea, want.Idea}, {"prd", got.PRD, want.PRD},
		{"techspec", got.TechSpec, want.TechSpec}, {"tasks", got.Tasks, want.Tasks},
		{"research", researchGot, &want.Research.AuthoringEntry},
	}
	for _, entry := range entries {
		path := "skills.authoring.claude." + entry.key
		if entry.got == nil {
			return fmt.Errorf("%s is missing", path)
		}
		if entry.got.Model != entry.want.Model {
			return fmt.Errorf("%s.model = %q, want %q", path, entry.got.Model, entry.want.Model)
		}
		if entry.got.ReasoningEffort != entry.want.ReasoningEffort {
			return fmt.Errorf("%s.reasoning_effort = %q, want %q", path, entry.got.ReasoningEffort, entry.want.ReasoningEffort)
		}
		if entry.got.Fallback != nil && !reflect.DeepEqual(entry.got.Fallback, entry.want.Fallback) {
			return fmt.Errorf("%s.fallback = %+v, want %+v", path, entry.got.Fallback, entry.want.Fallback)
		}
	}
	for _, tool := range want.Research.Tools {
		if !slices.Contains(got.Research.Tools, tool) {
			return fmt.Errorf("skills.authoring.claude.research.tools lacks recommended tool %q", tool)
		}
	}
	return nil
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
