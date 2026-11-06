package baseline

import (
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func TestArchiveRecordClausesAreAppended(t *testing.T) {
	t.Parallel()
	const specSentence = "A runtime whose archive leaves an Archive Record, such as `roundfix archive <slug>`, writes `<archive-root>/<slug>.md` and removes the Spec folder in the same change; the folder stays in Git at the record's `source_revision`, so copy every file a later reader needs to `docs/references/`, an accepted ADR, the glossary or an agent guide before or at archive, or send its knowledge in an ai-memory message to the project that owns the Secondbrain for its ingestion."
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
					if moduleID == "spec-workflow" && strings.Contains(clause.Guidance, "Secondbrain inbox") {
						t.Errorf("%s still names the retired Secondbrain inbox", clause.ID)
					}
					if moduleID == "context-workflow" && strings.Contains(clause.Guidance, "docs/history/findings/") {
						t.Errorf("%s still names a retired Findings history family", clause.ID)
					}
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
		if strings.Contains(string(data), "Secondbrain inbox") {
			t.Errorf("generated guide %s still names the retired Secondbrain inbox", path)
		}
		if !strings.Contains(string(data), specSentence) {
			t.Errorf("generated guide %s is missing archive record clause", path)
		}
	}
}

func TestTheArchiveClauseIsRetainedFromTheSourceBaseline(t *testing.T) {
	t.Parallel()
	catalog := mustEmbeddedCatalog(t)
	profile, err := ResolveProfile("", "standard-typescript-monorepo", catalog)
	if err != nil {
		t.Fatal(err)
	}
	active, artifacts, err := resolveManagedArtifacts(catalog, profile, standardTypeScriptDecisions("make verify"), false)
	if err != nil {
		t.Fatal(err)
	}
	managed := make([]ManifestArtifact, len(artifacts))
	for i, a := range artifacts {
		managed[i] = ManifestArtifact{ID: a.ID, Path: a.Path, Kind: a.Kind, Module: a.Module, Template: a.Template, Version: a.Version, Digest: a.Digest}
	}
	source, err := catalog.SourceBaseline("baseline.standard-typescript-monorepo-0.0.1")
	if err != nil {
		t.Fatal(err)
	}
	evidence, delta := classifySourceClauseTransition(source, managed, catalog, active)
	const clauseID = "clause.spec.keep-artifacts-in-spec-folder"
	if got := delta.Dispositions[clauseID]; got != ClauseRetained {
		t.Errorf("%s disposition = %s, want retained", clauseID, got)
	}
	count := 0
	for _, entry := range evidence {
		if entry.FromClause != clauseID {
			continue
		}
		count++
		if entry.Enforcement != "mandatory" || entry.Disposition != string(ClauseRetained) || !reflect.DeepEqual(entry.Targets, []string{clauseID}) {
			t.Errorf("%s retention evidence = %+v, want mandatory retention under its own identity", clauseID, entry)
		}
	}
	if count != 1 {
		t.Errorf("%s has %d evidence entries, want one", clauseID, count)
	}
}
