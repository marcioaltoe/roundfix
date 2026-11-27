package baseline

import (
	"context"
	"reflect"
	"strings"
	"testing"
)

func TestACodeLanguageExceptionSwapsTheCodeClauses(t *testing.T) {
	t.Parallel()
	for _, name := range []string{"a_declared_exception_states_the_language_and_the_source", "a_declaration_without_a_source_is_refused", "no_declaration_states_english"} {
		t.Run(name, func(t *testing.T) {
			repository := newProjectDecisionPlanRepository(t)
			before := snapshotWorkingTree(t, repository)
			decisions := standardTypeScriptDecisions("make verify")
			if name != "no_declaration_states_english" {
				value := map[string]any{"language": "Portuguese"}
				if name != "a_declaration_without_a_source_is_refused" {
					value["source"] = "docs/adr/0141-example.md"
				}
				decisions = append(decisions, DecisionValue{ID: codeLanguageExceptionDecisionID, Value: value})
			}
			outcome, err := BuildPlan(context.Background(), PlanRequest{Repository: repository, ProfileID: "standard-typescript-monorepo", Decisions: decisions, Preservation: RootPreservationRequest{Mode: PreservationModeGreenfield}})
			if name == "a_declaration_without_a_source_is_refused" {
				if err == nil || !strings.Contains(err.Error(), "no source, no exception") {
					t.Fatalf("refusal = %v", err)
				}
				if !reflect.DeepEqual(snapshotWorkingTree(t, repository), before) {
					t.Fatal("refusal wrote files")
				}
				return
			}
			if err != nil || outcome.Plan == nil {
				t.Fatalf("plan: %v %+v", err, outcome)
			}
			guide := string(planPostimage(t, *outcome.Plan, "docs/agents/agent-instructions.md").Content)
			exception := name != "no_declaration_states_english"
			if exception && !strings.Contains(guide, "The Code Language Exception: code and infrastructure names are written in Portuguese, as the existing code is. Source of the owner's decision: `docs/adr/0141-example.md`.") {
				t.Fatal(guide)
			}
			catalog := mustEmbeddedCatalog(t)
			for _, rule := range objectsOrEmpty(catalog.modules["core"]["rules"]) {
				if rule["id"] != "rule.core.code-and-reader-language" {
					continue
				}
				clauses := objectsOrEmpty(rule["clauses"])
				if len(clauses) != 5 {
					t.Fatalf("language rule clauses = %d, want 5", len(clauses))
				}
				for _, clause := range clauses {
					want := !exception
					if clause["id"] == "clause.core.reader-language" {
						want = true
					}
					if clause["id"] == "clause.core.code-names-follow-the-existing-codebase" {
						want = exception
					}
					text, _ := stringValue(clause, "guidance")
					if strings.Contains(strings.Join(strings.Fields(guide), " "), "- **mandatory**: "+text) != want {
						t.Errorf("clause %s presence != %t", clause["id"], want)
					}
				}
			}
			// Exercise punctuation at the renderer boundary as well as the full plan.
			rendered, err := renderProjectDecision(codeLanguageExceptionDecisionID, map[string]any{"language": "Portuguese *legacy*", "source": "decision `one`"}, catalog.decisions[codeLanguageExceptionDecisionID])
			if err != nil || !strings.Contains(rendered, `Portuguese \*legacy\*`) || !strings.Contains(rendered, "`` decision `one` ``") {
				t.Fatalf("escaped rendering: %q %v", rendered, err)
			}
		})
	}
}

func TestACodeLanguageExceptionValueIsRefused(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		name   string
		value  any
		reason string
	}{
		{"not_an_object", "Portuguese", "must be an object with language and source"},
		{"an_unknown_field", map[string]any{"other": true}, "accepts only language and source"},
		{"no_source", map[string]any{"language": "Portuguese"}, "requires a source, the recorded decision of the owner that keeps the codebase's language: no source, no exception"},
		{"a_blank_source", map[string]any{"source": " \t"}, "requires a source, the recorded decision of the owner that keeps the codebase's language: no source, no exception"},
		{"no_language", map[string]any{"source": "decision"}, "requires a language"},
		{"english_in_any_case", map[string]any{"language": "eNgLiSh", "source": "decision"}, "names a language other than English, which is the rule and needs no exception"},
		{"a_multiline_source", map[string]any{"language": "Portuguese", "source": "first\nsecond"}, "structured render code language exception source is non-canonical"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			_, err := normalizeCodeLanguageException(tc.value)
			if err == nil || err.Error() != tc.reason {
				t.Fatalf("normalize = %v, want %s", err, tc.reason)
			}
			err = validateDecisionValue(document{"type": codeLanguageExceptionDecisionType}, tc.value)
			if err == nil || err.Error() != tc.reason {
				t.Fatalf("decision = %v", err)
			}
			repository := newProjectDecisionPlanRepository(t)
			before := snapshotWorkingTree(t, repository)
			decisions := append(standardTypeScriptDecisions("make verify"), DecisionValue{ID: codeLanguageExceptionDecisionID, Value: tc.value})
			outcome, err := BuildPlan(context.Background(), PlanRequest{Repository: repository, ProfileID: "standard-typescript-monorepo", Decisions: decisions, Preservation: RootPreservationRequest{Mode: PreservationModeGreenfield}})
			if err == nil || !strings.HasSuffix(err.Error(), tc.reason) || outcome.Plan != nil {
				t.Fatalf("plan refusal: %v %+v", err, outcome)
			}
			if !reflect.DeepEqual(snapshotWorkingTree(t, repository), before) {
				t.Fatal("invalid declaration wrote files")
			}

		})
	}
}

func TestAClauseGateOnPresence(t *testing.T) {
	t.Parallel()
	catalog := mustEmbeddedCatalog(t)
	for _, tc := range []struct {
		name     string
		gate     map[string]any
		recorded map[string]any
		applies  bool
		invalid  bool
	}{
		{"present_true_applies_while_recorded", map[string]any{"decision": codeLanguageExceptionDecisionID, "present": true}, map[string]any{codeLanguageExceptionDecisionID: map[string]any{"language": "Portuguese", "source": "decision"}}, true, false},
		{"present_false_applies_while_absent", map[string]any{"decision": codeLanguageExceptionDecisionID, "present": false}, nil, true, false},
		{"a_default_is_not_presence", map[string]any{"decision": "branch.prefix", "present": true}, nil, false, false},
		{"equals_is_unchanged", map[string]any{"decision": "frontend.layout", "equals": "systems"}, nil, true, false},
		{"both_forms_are_refused", map[string]any{"decision": "frontend.layout", "equals": "systems", "present": true}, nil, false, true},
		{"a_non_boolean_present_is_refused", map[string]any{"decision": codeLanguageExceptionDecisionID, "present": "true"}, nil, false, true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			clause := document{"appliesWhen": tc.gate}
			loader := &catalogLoader{}
			loader.validateClauseGate(catalog, "test", clause)
			if (len(loader.diagnostics) > 0) != tc.invalid {
				t.Fatalf("diagnostics = %+v", loader.diagnostics)
			}
			if !tc.invalid && clauseApplies(catalog, clause, tc.recorded) != tc.applies {
				t.Fatal("gate result mismatch")
			}
			if tc.name == "present_true_applies_while_recorded" || tc.name == "present_false_applies_while_absent" {
				opposite := map[string]any{codeLanguageExceptionDecisionID: map[string]any{"language": "Portuguese", "source": "decision"}}
				clauseID := "clause.core.code-names-follow-the-existing-codebase"
				if tc.recorded != nil {
					opposite = nil
					clauseID = "clause.core.code-names-in-english"
				}
				if clauseApplies(catalog, clause, opposite) {
					t.Fatal("presence gate applies in the opposite state")
				}
				source := SourceBaseline{Entries: []SourceBaselineEntry{{ID: clauseID, Kind: "normative-clause", Enforcement: "mandatory", Carrier: "docs/agents/agent-instructions.md"}}}
				evidence, delta := classifySourceClauseTransition(source, []ManifestArtifact{{Kind: "guide", Path: "docs/agents/agent-instructions.md"}}, catalog, []string{"core"}, tc.recorded)
				if delta.Dispositions[clauseID] != ClauseReasonedRejection || len(evidence) != 1 || !strings.Contains(evidence[0].Reason, "presence") {
					t.Fatalf("presence transition: %+v %+v", delta, evidence)
				}
			}
		})
	}
}

func TestTheCatalogRefusesADefaultForACodeLanguageException(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct{ name, field string }{{"a_default", "default"}, {"a_suggestion", "suggestion"}} {
		t.Run(tc.name, func(t *testing.T) {
			catalog := cloneCatalogForRetentionDrift(t, mustEmbeddedCatalog(t))
			decision := catalog.decisions[codeLanguageExceptionDecisionID]
			decision[tc.field] = map[string]any{"language": "Portuguese", "source": "decision"}
			loader := &catalogLoader{}
			loader.validateOptionalDecision(codeLanguageExceptionDecisionID, decision)
			if len(loader.diagnostics) != 1 || loader.diagnostics[0].Code != "catalog.decision.optional.invalid" {
				t.Fatalf("diagnostics = %+v", loader.diagnostics)
			}
		})
	}
}

// snapshotWorkingTree records every file of the repository outside .git. A plan
// may run git status, which can rewrite .git/index when a file's timestamp is
// racily clean; that refresh is no write of the repository's files.
func snapshotWorkingTree(t *testing.T, repository string) map[string]manifestInputFileSnapshot {
	t.Helper()
	snapshot := snapshotManifestInputRepository(t, repository)
	for path := range snapshot {
		if strings.HasPrefix(path, ".git/") {
			delete(snapshot, path)
		}
	}
	return snapshot
}
