// Suite: Repository rule placement
// Invariant: Baseline guidance sends repository-authored rules to the canonical specific-repository guide.
// Boundary IN: embedded catalog, formatter goldens and a temporary Baseline plan.
// Boundary OUT: applying a repository refresh and Daemon Verification.

package baseline

import (
	"encoding/json"
	"strings"
	"testing"
)

const (
	keepRootCompactGuidance           = "Keep root agent instructions as short mandatory pointers. Setup owns only marked baseline content; preserve repository-authored bytes outside setup markers and keep project-specific architecture and policy in repository-owned documents. Write each rule or variation of this repository only in `docs/agents/specific-repository.md`, never in a root instruction file or in a Baseline guide, inside or outside its setup markers. Write it there under a heading of its own, as a bullet that states its force, such as `- **mandatory**: <rule>`, with the date of the decision that set it when it has one."
	activateMatchingSkillsGuidance    = "Activate every matching required skill before governed work. When one skill has distinct active-module triggers, retain and follow each trigger. When a skill's default conflicts with a Baseline rule or a Repository-Specific Normative Rule, follow the rule. Read a skill's bundled files from the directory where the skill is installed, never from a fixed absolute path the skill names. Write a stricter trigger, the model or effort a skill runs with, or any other dispatch rule of this repository in `docs/agents/specific-repository.md`, never in this guide."
	docsOneJobPerDirectoryGuidance    = "Give each documentation directory one job: `docs/_inbox/` for raw notes only when a repository has no Secondbrain, `docs/adr/` for decisions, `docs/agents/` for agent guidance, `docs/design/` for design artifacts, `docs/backlog/` for dated, typed intent not yet committed to a Spec, `docs/findings/` for dated investigations, `docs/handoffs/` for session continuity, `docs/references/` for external pointers and durable project reference documents, and `docs/user-guide/` for human documentation. Preserve repository-authored bytes outside setup markers, and add a repository rule only to `docs/agents/specific-repository.md`. Each directory holds live work only. The history root `docs/history/` holds Archive Records, and Git keeps their full text. Spec-contained authorization records live at `<spec-root>/<slug>/_authorization.md`; the preserved legacy location `docs/workflow/authorizations/` remains readable for historical grants. When the runtime's archive leaves an Archive Record, the record names each Spec-owned adopted reference, and the reference leaves the tree with the Spec folder."
	agentInstructionsTemplateSentence = "Repository authors write every project-specific rule or variation, stricter or not, only in `docs/agents/specific-repository.md`, never in this guide or in a root instruction file."
	skillDispatchTemplateSentence     = "This setup-owned guide is generated from the active modules. Activate every matching skill before governed work; write a stricter trigger of this repository in `docs/agents/specific-repository.md`, never in this guide."
	oldDocsLayoutSentence             = "Preserve repository-authored extensions outside setup markers."
)

func TestTheBaselineSendsRepositoryRulesToTheSpecificRepositoryGuide(t *testing.T) {
	t.Parallel()

	catalog := mustEmbeddedCatalog(t)
	tests := []struct {
		name string
		run  func(*testing.T)
	}{
		{
			name: "keep_root_compact",
			run: func(t *testing.T) {
				assertEmbeddedClause(t, catalog, "core", "rule.core.compact-root", "clause.core.keep-root-compact", keepRootCompactGuidance)
			},
		},
		{
			name: "activate_matching_skills",
			run: func(t *testing.T) {
				assertEmbeddedClause(t, catalog, "core", "rule.core.skill-dispatch", "clause.core.activate-matching-skills", activateMatchingSkillsGuidance)
			},
		},
		{
			name: "docs_one_job_per_directory",
			run: func(t *testing.T) {
				assertEmbeddedClause(t, catalog, "context-workflow", "rule.context.docs-layout", "clause.context.docs-one-job-per-directory", docsOneJobPerDirectoryGuidance)
			},
		},
		{
			name: "agent_instructions_guide",
			run: func(t *testing.T) {
				assertGuideGolden(t, catalog, "agent-instructions", keepRootCompactGuidance, agentInstructionsTemplateSentence)
			},
		},
		{
			name: "skill_dispatch_guide",
			run: func(t *testing.T) {
				assertGuideGolden(t, catalog, "skill-dispatch", activateMatchingSkillsGuidance, skillDispatchTemplateSentence)
			},
		},
		{
			name: "docs_layout_guide",
			run: func(t *testing.T) {
				assertGuideGolden(t, catalog, "docs-layout", docsOneJobPerDirectoryGuidance, "")
			},
		},
		{
			name: "root_pointer",
			run: func(t *testing.T) {
				repo := newPlanRepository(t)
				const rules = "# Repository rules\n\nKeep the repository boundary explicit.\n"
				writeInspectionFile(t, repo, legacyRepositoryPath, rules)
				commitInspectionRepository(t, repo, "seed legacy repository rules")
				plan := buildPlanWithRepositoryExtension(t, repo)
				agents := planPostimage(t, plan, "AGENTS.md")
				want := "- Read `docs/agents/specific-repository.md` for project-specific hard rules, and write every rule or variation of this repository there, never in this file or in a Baseline guide."
				normalized := strings.Join(strings.Fields(string(agents.Content)), " ")
				if count := strings.Count(normalized, want); count != 1 {
					t.Fatalf("AGENTS.md pointer count = %d, want 1", count)
				}
			},
		},
		{
			name: "starter",
			run: func(t *testing.T) {
				// The expected bytes come only from the embedded starter
				// template, never from a fixture this test writes.
				want := "# Repository-Specific Normative Rules\n\nWrite each rule or variation of this repository here, never in a root instruction file or in a Baseline guide. Give each rule a heading of its own and state it as a bullet with its force, such as `- **mandatory**: run the schema check before changing a migration (2026-11-16)`, with the date of the decision that set it when it has one.\n"
				starter, ok := catalog.TemplateContent("template.extension.repository-rules")
				if !ok {
					t.Fatal("missing template.extension.repository-rules")
				}
				if starter.Path != "templates/extensions/specific-repository.md" {
					t.Fatalf("starter template path = %q, want templates/extensions/specific-repository.md", starter.Path)
				}
				if string(starter.Data) != want {
					t.Fatalf("specific-repository.md starter = %q, want %q", starter.Data, want)
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, tt.run)
	}
}

func assertEmbeddedClause(t *testing.T, catalog *Catalog, moduleID, ruleID, clauseID, guidance string) {
	t.Helper()

	asset, ok := catalog.Module(moduleID)
	if !ok {
		t.Fatalf("missing module %s", moduleID)
	}
	var module document
	if err := json.Unmarshal(asset.Data, &module); err != nil {
		t.Fatalf("decode module %s: %v", moduleID, err)
	}
	foundRule := false
	foundClause := false
	for _, rule := range objectsOrEmpty(module["rules"]) {
		if rule["id"] != ruleID {
			continue
		}
		foundRule = true
		for _, clause := range objectsOrEmpty(rule["clauses"]) {
			if clause["id"] != clauseID {
				continue
			}
			foundClause = true
			if clause["enforcement"] != "mandatory" {
				t.Errorf("%s enforcement = %v, want mandatory", clauseID, clause["enforcement"])
			}
			if clause["guidance"] != guidance {
				t.Errorf("%s guidance = %v, want %q", clauseID, clause["guidance"], guidance)
			}
			if _, exists := clause["replaces"]; exists {
				t.Errorf("%s unexpectedly carries replaces", clauseID)
			}
			if clauseID == "clause.context.docs-one-job-per-directory" && strings.Contains(clause["guidance"].(string), oldDocsLayoutSentence) {
				t.Errorf("%s still carries the old placement sentence", clauseID)
			}
		}
	}
	if !foundRule {
		t.Fatalf("missing rule %s in module %s", ruleID, moduleID)
	}
	if !foundClause {
		t.Fatalf("missing clause %s in rule %s", clauseID, ruleID)
	}
}

func assertGuideGolden(t *testing.T, catalog *Catalog, guide, guidance, templateSentence string) {
	t.Helper()

	asset, ok := catalog.Asset("formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/" + guide + ".md")
	if !ok {
		t.Fatalf("missing Standard TypeScript %s golden", guide)
	}
	wantLine := "- **mandatory**: " + guidance
	lineCount := 0
	for _, line := range strings.Split(string(asset.Data), "\n") {
		if line == wantLine {
			lineCount++
		}
	}
	if lineCount != 1 {
		t.Fatalf("%s golden forced clause count = %d, want 1", guide, lineCount)
	}
	if templateSentence != "" {
		text := strings.Join(strings.Fields(string(asset.Data)), " ")
		want := strings.Join(strings.Fields(templateSentence), " ")
		if count := strings.Count(text, want); count != 1 {
			t.Fatalf("%s golden template sentence count = %d, want 1", guide, count)
		}
	}
}
