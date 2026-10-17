package baseline

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestHistorySanitizeClauseIsAppended(t *testing.T) {
	t.Parallel()
	const sentence = "The runtime's `roundfix history sanitize` removes what an earlier layout left."

	catalog, err := LoadEmbeddedCatalog()
	if err != nil {
		t.Fatal(err)
	}
	module, ok := catalog.Module("context-workflow")
	if !ok {
		t.Fatal("missing context-workflow module")
	}
	var document struct {
		Rules []struct {
			Clauses []struct {
				ID       string `json:"id"`
				Guidance string `json:"guidance"`
			} `json:"clauses"`
		} `json:"rules"`
	}
	if err := json.Unmarshal(module.Data, &document); err != nil {
		t.Fatal(err)
	}
	found := false
	for _, rule := range document.Rules {
		for _, clause := range rule.Clauses {
			if clause.ID == "clause.context.findings-09-archive" {
				found = true
				if !strings.Contains(clause.Guidance, sentence) {
					t.Fatal("embedded clause is missing the history sanitize sentence")
				}
			}
		}
	}
	if !found {
		t.Fatal("embedded clause is missing")
	}

	paths := []string{
		"docs/agents/docs-layout.md",
		"formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/docs-layout.md",
	}
	for _, path := range paths {
		var data []byte
		if path == "docs/agents/docs-layout.md" {
			data, err = os.ReadFile(filepath.Join("..", "..", path))
		} else {
			asset, found := catalog.Asset(path)
			if !found {
				t.Fatalf("missing generated guide %s", path)
			}
			data = asset.Data
		}
		if err != nil {
			t.Fatalf("read %s: %v", path, err)
		}
		if !strings.Contains(string(data), sentence) {
			t.Errorf("generated guide %s is missing the history sanitize sentence", path)
		}
	}
}
