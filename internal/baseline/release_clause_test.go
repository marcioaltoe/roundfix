package baseline

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestTheReleaseClauseNamesTheSkillsAndGuidesCheck(t *testing.T) {
	t.Parallel()
	const sentence = "Before the release Pull Request, confirm that the skills and guides the repository ships describe the behavior being released; the repository's release runbook owns that check."
	catalog := mustEmbeddedCatalog(t)
	module, ok := catalog.Module("core")
	if !ok {
		t.Fatal("embedded catalog has no core module")
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
		t.Fatalf("decode core module: %v", err)
	}
	var guidance string
	for _, rule := range document.Rules {
		for _, clause := range rule.Clauses {
			if clause.ID == "clause.core.plan-the-release-first" {
				guidance = clause.Guidance
			}
		}
	}
	if !strings.Contains(guidance, sentence) {
		t.Fatalf("embedded release clause is missing %q", sentence)
	}

	golden, ok := catalog.Asset("formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/agent-instructions.md")
	if !ok {
		t.Fatal("embedded formatter golden is missing")
	}
	if !strings.Contains(string(golden.Data), sentence) {
		t.Fatalf("formatter golden is missing %q", sentence)
	}
}
