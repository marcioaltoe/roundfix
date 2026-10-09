package baseline

import (
	"encoding/json"
	"strings"
	"testing"
)

const (
	qaOverrideCommand             = "roundfix archive <slug> --qa-override --approval <source> --reason <text>"
	qaOverrideHandEditProhibition = "Never hand-edit the override stamp"
	qaOverrideUnsupportedRuntime  = "A runtime command without override support must be reported as unsupported rather than given an invented flag."
)

func TestSpecDocsLayoutClauseRoutesTheOverrideThroughTheCommand(t *testing.T) {
	t.Parallel()
	guidance := specDocsLayoutClauseGuidance(t)

	for _, required := range []string{
		qaOverrideCommand,
		qaOverrideHandEditProhibition,
		qaOverrideUnsupportedRuntime,
	} {
		if !strings.Contains(guidance, required) {
			t.Errorf("embedded docs-layout clause is missing %q", required)
		}
	}
}

func TestSpecDocsLayoutClauseNoLongerAsksForAHandStamp(t *testing.T) {
	t.Parallel()
	const obsoleteHandStamp = "preserve any supplied reason"

	if guidance := specDocsLayoutClauseGuidance(t); strings.Contains(guidance, obsoleteHandStamp) {
		t.Errorf("embedded docs-layout clause still contains %q", obsoleteHandStamp)
	}
	if golden := formatterGoldenDocsLayout(t); strings.Contains(golden, obsoleteHandStamp) {
		t.Errorf("formatter golden docs-layout guide still contains %q", obsoleteHandStamp)
	}
}

func TestFormatterGoldenDocsLayoutCarriesTheOverrideCommand(t *testing.T) {
	t.Parallel()
	golden := formatterGoldenDocsLayout(t)

	for _, required := range []string{
		qaOverrideCommand,
		qaOverrideHandEditProhibition,
	} {
		if !strings.Contains(golden, required) {
			t.Errorf("formatter golden docs-layout guide is missing %q", required)
		}
	}
}

func specDocsLayoutClauseGuidance(t *testing.T) string {
	t.Helper()

	catalog, err := LoadEmbeddedCatalog()
	if err != nil {
		t.Fatalf("load embedded catalog: %v", err)
	}
	module, ok := catalog.Module("spec-workflow")
	if !ok {
		t.Fatal("embedded catalog has no spec-workflow module")
	}
	var document struct {
		Rules []struct {
			ID      string `json:"id"`
			Clauses []struct {
				ID       string `json:"id"`
				Guidance string `json:"guidance"`
			} `json:"clauses"`
		} `json:"rules"`
	}
	if err := json.Unmarshal(module.Data, &document); err != nil {
		t.Fatalf("decode spec-workflow module: %v", err)
	}
	for _, rule := range document.Rules {
		if rule.ID != "rule.spec.docs-layout" {
			continue
		}
		for _, clause := range rule.Clauses {
			if clause.ID == "clause.spec.keep-artifacts-in-spec-folder" {
				return clause.Guidance
			}
		}
	}
	t.Fatal("spec-workflow module has no docs-layout artifact clause")
	return ""
}

func formatterGoldenDocsLayout(t *testing.T) string {
	t.Helper()

	catalog, err := LoadEmbeddedCatalog()
	if err != nil {
		t.Fatalf("load embedded catalog: %v", err)
	}
	const assetPath = "formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/docs-layout.md"
	asset, ok := catalog.Asset(assetPath)
	if !ok {
		t.Fatalf("embedded catalog has no asset %q", assetPath)
	}
	return string(asset.Data)
}
