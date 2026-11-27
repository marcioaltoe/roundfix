package baseline

import (
	"context"
	"fmt"
	"slices"
	"strings"
	"testing"
)

func TestEveryBuiltInProfileRendersTheLanguageRule(t *testing.T) {
	t.Parallel()
	render := func(t *testing.T, catalog *Catalog, profile string) PlanDocument {
		t.Helper()
		repository := newProjectDecisionPlanRepository(t)
		writeProfileAlignmentFile(t, repository, "go.mod", "module example.invalid/language\n\ngo 1.26\n")
		resolved, err := ResolveProfile(repository, profile, catalog)
		if err != nil {
			t.Fatal(err)
		}
		decisions := slices.DeleteFunc(standardTypeScriptDecisions("make verify"), func(d DecisionValue) bool { return !slices.Contains(resolved.Decisions, d.ID) })
		outcome, err := buildPlanWithCatalog(context.Background(), PlanRequest{
			Repository: repository, ProfileID: profile,
			Decisions:    decisions,
			Preservation: RootPreservationRequest{Mode: PreservationModeGreenfield},
		}, catalog)
		if err != nil || outcome.Plan == nil {
			t.Fatalf("render %s: %v %+v", profile, err, outcome.Result)
		}
		return *outcome.Plan
	}
	for _, profile := range []string{"go-cli-tui", "go-cli-typescript-monorepo", "python", "rust-cli", "standard-typescript-monorepo"} {
		t.Run(profile, func(t *testing.T) {
			if err := languageRuleRenderingError(t, render(t, mustEmbeddedCatalog(t), profile)); err != nil {
				t.Fatal(err)
			}
		})
	}
	t.Run("a_catalog_without_a_clause_fails", func(t *testing.T) {
		catalog := cloneCatalogForRetentionDrift(t, mustEmbeddedCatalog(t))
		removeCatalogClause(t, catalog, "core", "clause.core.code-names-in-english")
		err := languageRuleRenderingError(t, render(t, catalog, "go-cli-tui"))
		if err == nil || !strings.Contains(err.Error(), "Write code and infrastructure names in English") {
			t.Fatalf("missing-clause assertion = %v", err)
		}
	})
}

func languageRuleRenderingError(t *testing.T, plan PlanDocument) error {
	t.Helper()
	guide := strings.Join(strings.Fields(string(planPostimage(t, plan, "docs/agents/agent-instructions.md").Content)), " ")
	for _, want := range []string{
		"- **mandatory**: Write code and infrastructure names in English: variables, functions, classes, files and folders, code comments, database tables, columns and indexes, buckets, prefixes, environment variables, queues, API routes and technical logs.",
		"- **mandatory**: A business term with no faithful English translation keeps its original form inside an English name: an official acronym as it is (`icms_amount`), another term without accents (`apuracao_summary`). Brazilian fiscal and legal terms are the usual case.",
		"- **mandatory**: A new name follows this rule. An existing name changes only when its code is touched, a large rename takes a Spec of its own, and nothing is renamed in bulk without a request. The Pre-PR Review reports a new name that breaks the rule, and the QA gate reports text a user reads outside the Reader Language.",
		"- **mandatory**: Write what the product's users read in the Reader Language named above: screens, user messages, e-mails, PDFs, user documentation, and the names, descriptions and parameters of MCP tools. Specs, ADRs, agent guides, commits and pull requests keep the language they already have.",
	} {
		if !strings.Contains(guide, want) {
			return fmt.Errorf("missing language rule: %s", want)
		}
	}
	if !strings.Contains(guide, "The Reader Language is English.") {
		return fmt.Errorf("missing language rule: The Reader Language is English.")
	}
	domain := strings.Join(strings.Fields(string(planPostimage(t, plan, "docs/agents/domain.md").Content)), " ")
	if !strings.Contains(domain, "- **mandatory**: Use the repository's canonical domain terms in code names, tests, user-facing copy, Specs, and delivery notes. In a code name, a term takes its Name in Code: the term itself, or the name its Glossary Definition records on a line `_In code_:` when the term was translated to English or kept without accents. A Spec that introduces or changes such a term writes that line. Call out a missing term instead of inventing a competing synonym.") {
		return fmt.Errorf("missing canonical-terms clause: %s", domain)
	}
	return nil
}
