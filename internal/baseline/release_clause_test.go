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

func TestTheReleaseClauseLetsACeilingAnswerTheApproval(t *testing.T) {
	t.Parallel()
	catalog := mustEmbeddedCatalog(t)
	for _, tt := range []struct {
		module, clause, guide, sentence string
	}{
		{
			module:   "core",
			clause:   "clause.core.plan-the-release-first",
			guide:    "agent-instructions.md",
			sentence: "A Release Ceiling the maintainer recorded in the `_authorization.md` of every Spec the release delivers answers a minor or major plan at or below that ceiling, never a breaking or manually classified one.",
		},
		{
			module:   "spec-workflow",
			clause:   "clause.spec.project-constraints-02-tooling-authorization",
			guide:    "spec-routing.md",
			sentence: "the impact ceiling a `release` grant answers without a second approval (`release_ceiling`: `patch`, `minor` or `major`) with the version-only paths that release may change (`release_paths`),",
		},
	} {
		t.Run(tt.module, func(t *testing.T) {
			t.Parallel()
			module, ok := catalog.Module(tt.module)
			if !ok {
				t.Fatalf("embedded catalog has no %s module", tt.module)
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
				t.Fatalf("decode %s module: %v", tt.module, err)
			}
			var guidance string
			for _, rule := range document.Rules {
				for _, clause := range rule.Clauses {
					if clause.ID == tt.clause {
						guidance = clause.Guidance
					}
				}
			}
			if !strings.Contains(guidance, tt.sentence) {
				t.Errorf("embedded clause %s is missing %q", tt.clause, tt.sentence)
			}

			path := "formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/" + tt.guide
			golden, ok := catalog.Asset(path)
			if !ok {
				t.Fatalf("embedded formatter golden %s is missing", path)
			}
			if !strings.Contains(strings.Join(strings.Fields(string(golden.Data)), " "), tt.sentence) {
				t.Errorf("formatter golden %s is missing %q", path, tt.sentence)
			}
		})
	}
}
