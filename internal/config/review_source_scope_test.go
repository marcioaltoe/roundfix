// Suite: generated Review Source configuration.
// Invariant: generated config labels legacy PR feedback without selecting a pre-PR reviewer.
// Boundary IN: generated User and Project Config text and Project Config loading.
// Boundary OUT: operational review commands and legacy request coherence.
package config

import (
	"path/filepath"
	"strings"
	"testing"
)

func TestDefaultConfigScopesReviewSourceToPullRequestFeedback(t *testing.T) {
	t.Parallel()

	for _, generated := range []struct {
		name    string
		content string
	}{
		{name: "User Config", content: DefaultConfigYAML()},
		{name: "Project Config", content: DefaultProjectConfigYAML()},
	} {
		generated := generated
		t.Run(generated.name, func(t *testing.T) {
			t.Parallel()
			for _, want := range []string{
				"read only by fetch, watch and resolve",
				"request_review: false",
			} {
				if !strings.Contains(generated.content, want) {
					t.Fatalf("generated %s is missing %q:\n%s", generated.name, want, generated.content)
				}
			}
		})
	}
}

func TestDefaultConfigSelectsNoPrePRReviewProvider(t *testing.T) {
	t.Parallel()

	userConfig := DefaultConfigYAML()
	projectConfig := DefaultProjectConfigYAML()
	for name, content := range map[string]string{
		"User Config":    userConfig,
		"Project Config": projectConfig,
	} {
		if strings.Contains(content, "pre_pr_review") {
			t.Fatalf("generated %s selects a pre-PR review provider:\n%s", name, content)
		}
	}

	homeDir := t.TempDir()
	workDir := t.TempDir()
	mustMkdir(t, filepath.Join(workDir, ".git"))
	mustWrite(t, filepath.Join(workDir, projectConfigName), projectConfig)

	loaded, err := Load(LoadOptions{HomeDir: homeDir, WorkDir: workDir})
	if err != nil {
		t.Fatalf("load generated Project Config: %v", err)
	}
	if got := loaded.Config.PrePRReview.Provider; got != "codex" {
		t.Fatalf("PrePRReview.Provider = %q, want %q", got, "codex")
	}
	if got := loaded.Config.PrePRReview.Source; got != "default" {
		t.Fatalf("PrePRReview.Source = %q, want %q", got, "default")
	}
}
