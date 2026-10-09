package baseline

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestArchiveRecordClausesAreAppended(t *testing.T) {
	t.Parallel()
	const specSentence = "A runtime whose archive leaves an Archive Record, such as `roundfix archive <slug>`, writes `<archive-root>/<slug>.md` and removes the Spec folder in the same change; the folder stays in Git at the record's `source_revision`, so copy every file a later reader needs to `docs/references/`, an accepted ADR, the glossary, an agent guide or the Secondbrain inbox before or at archive."
	const contextSentence = "When the runtime's archive leaves an Archive Record, the record names each Spec-owned adopted reference, and the reference leaves the tree with the Spec folder."
	catalog, err := LoadEmbeddedCatalog()
	if err != nil {
		t.Fatal(err)
	}
	for moduleID, entry := range map[string]struct{ clauseID, sentence string }{"spec-workflow": {"clause.spec.keep-artifacts-in-spec-folder", specSentence}, "context-workflow": {"clause.context.docs-one-job-per-directory", contextSentence}} {
		module, ok := catalog.Module(moduleID)
		if !ok {
			t.Fatalf("missing module %s", moduleID)
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
				if clause.ID == entry.clauseID {
					found = strings.Contains(clause.Guidance, entry.sentence)
				}
			}
		}
		if !found {
			t.Errorf("embedded %s is missing appended sentence", entry.clauseID)
		}
	}
	for _, path := range []string{"docs/agents/docs-layout.md", "formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/docs-layout.md"} {
		var data []byte
		if path == "docs/agents/docs-layout.md" {
			var err error
			data, err = os.ReadFile(filepath.Join("..", "..", path))
			if err != nil {
				t.Fatal(err)
			}
		} else {
			asset, ok := catalog.Asset(path)
			if !ok {
				t.Fatalf("missing generated guide %s", path)
			}
			data = asset.Data
		}
		if !strings.Contains(string(data), specSentence) {
			t.Errorf("generated guide %s is missing archive record clause", path)
		}
	}
}
