//go:build docscontract

package docscontract

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"roundfix/internal/config"
)

func TestModelSelectionReferenceStatesTheShippedSnapshot(t *testing.T) {
	root := baselineDocumentationRepoRoot()
	path := filepath.Join(root, "docs", "references", "model-selection.md")
	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	reference := string(content)

	wantUpdated := "Updated " + config.ModelRecommendationSnapshotVersion + "."
	if !strings.Contains(reference, wantUpdated) {
		t.Fatalf("reference missing snapshot marker %q", wantUpdated)
	}

	for _, category := range config.AllWorkCategories() {
		category := category
		t.Run(string(category), func(t *testing.T) {
			profile, ok := config.RecommendedProfile(category)
			if !ok {
				t.Fatalf("no Recommended Profile for %q", category)
			}

			selections := append([]config.AgentSelection{profile.Preferred}, profile.Fallbacks...)
			for _, selection := range selections {
				want := strings.Join([]string{selection.Runtime, selection.Model, selection.ReasoningEffort}, " / ")
				if !strings.Contains(reference, want) {
					t.Errorf("reference missing Recommended Profile selection %q", want)
				}
			}
		})
	}
}
