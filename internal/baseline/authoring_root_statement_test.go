// Suite: Authoring Phase root statement
// Invariant: the Spec workflow root block states the Claude Code model rule and its harness scope once.
// Boundary IN: embedded template, temporary Baseline plan, and Standard TypeScript formatter golden.
// Boundary OUT: applying a repository refresh and user-level instruction files.

package baseline

import (
	"strings"
	"testing"
)

func TestTheSpecWorkflowRootBlockStatesHowAnAuthoringPhaseKeepsItsModel(t *testing.T) {
	t.Parallel()

	const firstSentence = "In Claude Code, a skill's `model` and `effort` header applies only when a person types the skill (`/skill`); a skill an Agent activates runs on the session's model (anthropics/claude-code#98898, measured 2026-11-14)."
	const lastSentence = "This bullet does not apply to Codex or OpenCode, whose models come from the Agent Selection Profiles."
	catalog := mustEmbeddedCatalog(t)
	template, ok := catalog.TemplateContent("template.root.spec-workflow")
	if !ok {
		t.Fatal("missing Spec workflow root template")
	}
	assertSentences := func(t *testing.T, source, content string) {
		t.Helper()
		normalized := strings.Join(strings.Fields(content), " ")
		for _, sentence := range []string{firstSentence, lastSentence} {
			if count := strings.Count(normalized, sentence); count != 1 {
				t.Errorf("%s sentence count = %d, want 1: %q", source, count, sentence)
			}
		}
	}

	t.Run("plan_postimage", func(t *testing.T) {
		plan := buildTestPlan(t, newPlanRepository(t))
		content := string(planPostimage(t, plan, "AGENTS.md").Content)
		assertSentences(t, "AGENTS.md postimage", content)
		_, block, ok := strings.Cut(content, "<!-- setup-context-driven:begin id=root.spec-workflow ")
		if !ok {
			t.Fatal("AGENTS.md postimage lacks the Spec workflow begin marker")
		}
		block, _, ok = strings.Cut(block, "<!-- setup-context-driven:end id=root.spec-workflow -->")
		if !ok {
			t.Fatal("AGENTS.md postimage lacks the Spec workflow end marker")
		}
		assertSentences(t, "root.spec-workflow block", block)
	})

	t.Run("formatter_golden", func(t *testing.T) {
		golden, ok := catalog.Asset("formatter-fixtures/standard-typescript-monorepo/golden/AGENTS.md")
		if !ok {
			t.Fatal("missing Standard TypeScript AGENTS.md formatter golden")
		}
		assertSentences(t, "Spec workflow template", string(template.Data))
		assertSentences(t, "Standard TypeScript AGENTS.md formatter golden", string(golden.Data))
	})

	t.Run("scoped_to_claude_code", func(t *testing.T) {
		assertSentences(t, "Spec workflow template", string(template.Data))
		for _, bullet := range strings.Split(string(template.Data), "\n- ")[1:] {
			normalized := strings.Join(strings.Fields(bullet), " ")
			if strings.Contains(normalized, lastSentence) {
				if !strings.HasPrefix(normalized, firstSentence) {
					t.Errorf("authoring bullet must open with %q, got %q", firstSentence, normalized)
				}
				return
			}
		}
		t.Fatal("missing authoring bullet scoped to Claude Code")
	})
}
