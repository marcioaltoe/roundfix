// Suite: Cross-project continuity
// Invariant: ai-memory replaces every inbox clause with the same force and guides safe delivery.
// Boundary IN: embedded catalog, Source Baseline transitions, and Plan postimages.
// Boundary OUT: applying a Plan, live ai-memory, and the Secondbrain filesystem.

package baseline

import (
	"reflect"
	"strings"
	"testing"
)

func TestEveryInboxClauseAnAdopterHoldsIsReplaced(t *testing.T) {
	t.Parallel()
	catalog := mustEmbeddedCatalog(t)
	profile, err := ResolveProfile("", "standard-typescript-monorepo", catalog)
	if err != nil {
		t.Fatal(err)
	}
	decisions := standardTypeScriptDecisions("make verify")
	for i := range decisions {
		if decisions[i].ID == "secondbrain.enabled" {
			decisions[i].Value = true
		}
	}
	active, artifacts, err := resolveManagedArtifacts(catalog, profile, decisions, false)
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
	for id, disposition := range delta.Dispositions {
		if disposition == ClauseUnaccounted {
			t.Errorf("unaccounted clause %s", id)
		}
	}
	replacements := map[string]string{
		"inbox-auto-capture":     "clause.core.ai-memory-is-optional",
		"inbox-write-permission": "clause.core.ai-memory-is-optional",
		"capture-trigger":        "clause.core.message-the-owning-project",
		"capture-self-contained": "clause.core.messages-are-self-contained",
		"inbox-entry-contract":   "clause.core.messages-are-self-contained",
		"inbox-empty-state":      "clause.core.consume-messages-after-recording",
		"inbox-triage-order":     "clause.core.consume-messages-after-recording",
		"research-capture":       "clause.secondbrain.research-to-the-knowledge-owner",
	}
	for old, successor := range replacements {
		id := "clause.secondbrain." + old
		if got := delta.Dispositions[id]; got != ClauseReplaced {
			t.Errorf("%s disposition = %s, want replaced", id, got)
		}
		count := 0
		for _, entry := range evidence {
			if entry.FromClause != id {
				continue
			}
			count++
			if entry.Enforcement != "mandatory" || entry.Disposition != string(ClauseReplaced) || !reflect.DeepEqual(entry.Targets, []string{successor}) {
				t.Errorf("%s replacement evidence = %+v, want mandatory replacement by %s", id, entry, successor)
			}
		}
		if count != 1 {
			t.Errorf("%s has %d evidence entries, want one", id, count)
		}
	}
	for _, id := range []string{"clause.secondbrain.prohibit-writes", "clause.secondbrain.baseline-owned-guidance", "clause.context.inbox-01-triage", "clause.context.inbox-02-fleet-flow", "clause.context.inbox-03-extend-before-minting", "clause.spec.sources-01-group-by-shared-context"} {
		if got := delta.Dispositions[id]; got != ClauseRetained {
			t.Errorf("%s disposition = %s, want retained", id, got)
		}
	}
}

func TestTheGuidesSendCrossProjectWorkThroughAiMemory(t *testing.T) {
	t.Parallel()
	decisions := standardTypeScriptDecisions("make verify")
	for i := range decisions {
		if decisions[i].ID == "secondbrain.enabled" {
			decisions[i].Value = true
		}
	}
	plan := buildProjectDecisionPlan(t, newProjectDecisionPlanRepository(t), decisions)
	instructions := strings.Join(strings.Fields(string(planPostimage(t, plan, "docs/agents/agent-instructions.md").Content)), " ")
	for _, phrase := range []string{
		"ai-memory is optional and detected at run time",
		"tell the user that the ai-memory MCP server or that project is missing",
		"put the handoff or the message text in the response for the user to deliver",
		"sent to that project's inbox with `memory_message_send`",
		"Consume a message with `memory_message_pop` only after its content is recorded in this repository",
	} {
		if !strings.Contains(instructions, phrase) {
			t.Errorf("agent instructions lack %q", phrase)
		}
	}
	dispatch := string(planPostimage(t, plan, "docs/agents/skill-dispatch.md").Content)
	for _, skill := range []string{"ai-memory-handoff", "ai-memory-messaging"} {
		if !strings.Contains(dispatch, "`"+skill+"`") {
			t.Errorf("skill dispatch lacks %s", skill)
		}
	}
	for _, guide := range []string{"agent-instructions.md", "skill-dispatch.md", "secondbrain.md", "docs-layout.md", "spec-routing.md"} {
		text := string(planPostimage(t, plan, "docs/agents/"+guide).Content)
		for _, forbidden := range []string{"Inbox Entr", "inbox/<destination>", "_triaged", "`handoff`"} {
			if strings.Contains(text, forbidden) {
				t.Errorf("%s still names %q", guide, forbidden)
			}
		}
	}
}
