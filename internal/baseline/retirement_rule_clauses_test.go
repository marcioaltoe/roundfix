package baseline

import (
	"encoding/json"
	"strings"
	"testing"

	"roundfix/internal/spec"
)

func embeddedClauseByID(t *testing.T, id string) characterizedBaselineClause {
	t.Helper()
	for _, clause := range embeddedBaselineClauses(t) {
		if clause.ID == id {
			return clause
		}
	}
	t.Fatalf("missing clause %s", id)
	return characterizedBaselineClause{}
}

func TestTheRetirementRuleIsStatedOnce(t *testing.T) {
	t.Parallel()
	catalog := mustEmbeddedCatalog(t)
	asset, ok := catalog.Module("context-workflow")
	if !ok {
		t.Fatal("missing context-workflow module")
	}
	var module document
	if err := json.Unmarshal(asset.Data, &module); err != nil {
		t.Fatal(err)
	}
	for _, rule := range objectsOrEmpty(module["rules"]) {
		for _, clause := range objectsOrEmpty(rule["clauses"]) {
			id, _ := clause["id"].(string)
			guidance, _ := clause["guidance"].(string)
			if id == "clause.context.findings-09-archive" {
				continue
			}
			for _, forbidden := range []string{
				"docs/history/findings/", "docs/history/backlog/", "docs/history/adr/",
				"reduced entry", "stay whole", "capture door",
			} {
				if strings.Contains(guidance, forbidden) {
					t.Errorf("%s names retired history guidance %q", id, forbidden)
				}
			}
		}
	}
}

func TestTheDocsLayoutRegionStaysLean(t *testing.T) {
	t.Parallel()
	catalog := mustEmbeddedCatalog(t)
	asset, ok := catalog.Asset("formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/docs-layout.md")
	if !ok {
		t.Fatal("missing docs-layout formatter golden")
	}
	begin := "<!-- setup-context-driven:begin id=guide.docs-layout"
	end := "<!-- setup-context-driven:end id=guide.docs-layout -->"
	text := string(asset.Data)
	start := strings.Index(text, begin)
	finish := strings.Index(text, end)
	if start < 0 || finish < start {
		t.Fatalf("docs-layout managed region markers missing")
	}
	if size := len(text[start : finish+len(end)]); size > 11000 {
		t.Fatalf("guide.docs-layout region = %d bytes, want at most 11000", size)
	}
}

func TestTheRetirementClauseCarriesTheSentinel(t *testing.T) {
	t.Parallel()
	clause := embeddedClauseByID(t, "clause.context.findings-09-archive")
	if !strings.HasPrefix(clause.Guidance, spec.RetirementRuleSentinel) || !strings.Contains(clause.Guidance, spec.RetirementRuleSentinel) {
		t.Fatalf("retirement clause does not open with %q: %q", spec.RetirementRuleSentinel, clause.Guidance)
	}
}

func TestTheRetirementClauseReplacesTheClausesItAbsorbs(t *testing.T) {
	t.Parallel()
	clause := embeddedClauseByID(t, "clause.context.findings-09-archive")
	var raw struct {
		Replaces []string `json:"replaces"`
	}
	data, err := json.Marshal(clause)
	if err != nil {
		t.Fatal(err)
	}
	if err := json.Unmarshal(data, &raw); err != nil {
		t.Fatal(err)
	}
	want := []string{"clause.context.findings-10-live-work-health", "clause.context.findings-11-rollup-closure"}
	if len(raw.Replaces) != len(want) || raw.Replaces[0] != want[0] || raw.Replaces[1] != want[1] {
		t.Fatalf("replaces = %v, want %v", raw.Replaces, want)
	}
}

func TestTheCaptureTriggerAsksForAConditionalLesson(t *testing.T) {
	t.Parallel()
	catalog := mustEmbeddedCatalog(t)
	asset, ok := catalog.Module("core")
	if !ok {
		t.Fatal("missing core module")
	}
	var module document
	if err := json.Unmarshal(asset.Data, &module); err != nil {
		t.Fatal(err)
	}
	for _, rule := range objectsOrEmpty(module["rules"]) {
		for _, clause := range objectsOrEmpty(rule["clauses"]) {
			if clause["id"] != "clause.core.message-the-owning-project" {
				continue
			}
			guidance := clause["guidance"].(string)
			for _, phrase := range []string{
				"When a Finding closes",
				"ADR, Archive Record or reference",
				"Only then send one message",
				"<repository>@<commit>:<path>",
			} {
				if !strings.Contains(guidance, phrase) {
					t.Errorf("capture trigger lacks %q", phrase)
				}
			}
			return
		}
	}
	t.Fatal("capture trigger is missing")
}
