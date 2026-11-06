// Suite: Baseline Profile decision gaps
// Invariant: only unreachable required bindings on unconditional profile templates are gaps.
// Boundary IN: catalog decisions, module artifacts, profile reachability, diagnostics and plan refusal.
// Boundary OUT: CLI streams and update orchestration.

package baseline

import (
	"encoding/json"
	"os"
	"slices"
	"testing"
)

func TestNoBuiltInProfileHasADecisionGap(t *testing.T) {
	t.Parallel()
	catalog, err := LoadEmbeddedCatalog()
	if err != nil {
		t.Fatal(err)
	}
	for id := range catalog.profiles {
		t.Run(id, func(t *testing.T) {
			t.Parallel()
			profile, err := resolveBuiltInProfile(id, catalog)
			if err != nil {
				t.Fatal(err)
			}
			if gaps := ProfileDecisionGaps(profile, catalog); len(gaps) != 0 {
				t.Fatalf("built-in gaps = %+v", gaps)
			}
		})
	}
}

func TestAProfileWithoutTheIncrementalDecisionHasOneGap(t *testing.T) {
	t.Parallel()
	catalog, err := LoadEmbeddedCatalog()
	if err != nil {
		t.Fatal(err)
	}
	for id := range catalog.profiles {
		t.Run(id, func(t *testing.T) {
			t.Parallel()
			profile, err := resolveBuiltInProfile(id, catalog)
			if err != nil {
				t.Fatal(err)
			}
			profile.Decisions = slices.DeleteFunc(slices.Clone(profile.Decisions), func(id string) bool { return id == "verification.incremental" })
			want := []ProfileDecisionGap{{Decision: "verification.incremental", Template: "template.guide.agent-instructions", Artifact: "guide.agent-instructions", Suggested: "make verify-incremental"}}
			if gaps := ProfileDecisionGaps(profile, catalog); !slices.Equal(gaps, want) {
				t.Fatalf("gaps = %+v, want %+v", gaps, want)
			}
		})
	}
}

func TestProfileDecisionGapTextNamesTheLineAndTheFlag(t *testing.T) {
	t.Parallel()
	profile := ResolvedProfile{ID: "sqllib", Path: ".roundfix/baseline/profiles/sqllib.json"}
	incremental := ProfileDecisionGap{Decision: "verification.incremental", Template: "template.guide.agent-instructions", Suggested: "make verify-incremental"}
	gate := ProfileDecisionGap{Decision: "verification.gate", Template: "template.guide.agent-instructions", Suggested: "make verify"}
	tests := []struct {
		name            string
		gaps            []ProfileDecisionGap
		message, repair string
	}{
		{"one", []ProfileDecisionGap{incremental},
			"profile sqllib does not list decision verification.incremental, which template template.guide.agent-instructions requires",
			`add "verification.incremental" to "decisions" in .roundfix/baseline/profiles/sqllib.json, then rerun with --decision 'verification.incremental=make verify-incremental'`},
		{"two sorted", []ProfileDecisionGap{gate, incremental},
			"profile sqllib does not list decision verification.gate, which template template.guide.agent-instructions requires; profile sqllib does not list decision verification.incremental, which template template.guide.agent-instructions requires",
			`add "verification.gate", "verification.incremental" to "decisions" in .roundfix/baseline/profiles/sqllib.json, then rerun with --decision 'verification.gate=make verify' --decision 'verification.incremental=make verify-incremental'`},
		{"empty suggestion", []ProfileDecisionGap{{Decision: "branch.prefix", Template: "template.guide.agent-instructions"}},
			"profile sqllib does not list decision branch.prefix, which template template.guide.agent-instructions requires",
			`add "branch.prefix" to "decisions" in .roundfix/baseline/profiles/sqllib.json, then rerun with --decision 'branch.prefix=<value>'`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := ProfileDecisionGapMessage(profile, tt.gaps); got != tt.message {
				t.Fatalf("message = %q, want %q", got, tt.message)
			}
			if got := ProfileDecisionGapRepair(profile, tt.gaps, "rerun with"); got != tt.repair {
				t.Fatalf("repair = %q, want %q", got, tt.repair)
			}
		})
	}
	// A decision can bind more than one template; its repair is still listed once.
	if got := ProfileDecisionGapRepair(profile, []ProfileDecisionGap{incremental, gate, incremental}, "rerun with"); got != tests[1].repair {
		t.Fatalf("deduplicated repair = %q", got)
	}
}

func TestBuildPlanRefusesAProfileDecisionGapBeforeRendering(t *testing.T) {
	t.Parallel()
	catalog, err := LoadEmbeddedCatalog()
	if err != nil {
		t.Fatal(err)
	}
	repository := newPlanRepository(t)
	created, err := InitCustomProfile(repository, "sqllib", "go-cli-tui", catalog)
	if err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(created.Path)
	if err != nil {
		t.Fatal(err)
	}
	var profile customProfileDocument
	if err := json.Unmarshal(data, &profile); err != nil {
		t.Fatal(err)
	}
	profile.Decisions = slices.DeleteFunc(profile.Decisions, func(id string) bool { return id == "verification.incremental" })
	data, err = json.Marshal(profile)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(created.Path, data, 0o644); err != nil {
		t.Fatal(err)
	}
	outcome, err := buildPlanWithCatalog(t.Context(), PlanRequest{Repository: repository, ProfileID: "sqllib", Preservation: RootPreservationRequest{Mode: PreservationModeGreenfield}}, catalog)
	if err != nil {
		t.Fatal(err)
	}
	const message = "profile sqllib does not list decision verification.incremental, which template template.guide.agent-instructions requires"
	const repair = `add "verification.incremental" to "decisions" in .roundfix/baseline/profiles/sqllib.json, then rerun with --decision 'verification.incremental=make verify-incremental'`
	if outcome.Plan != nil || outcome.Result.State != "action_required" || outcome.Result.Category != "decision" || outcome.Result.Message != message || outcome.Result.NextAction != repair {
		t.Fatalf("gap plan outcome = %+v", outcome)
	}
}

func TestProfileDecisionGapReachabilityAndTemplateSelection(t *testing.T) {
	t.Parallel()
	// A small in-memory catalog isolates the graph rule from embedded asset changes.
	binding := map[string]any{"artifact": "guide", "template": "template.guide", "token": "value"}
	required := map[string]any{"effects": []any{map[string]any{"renderBindings": []any{binding, binding}}}, "default": false}
	tests := []struct {
		name      string
		decisions map[string]document
		selected  []string
		wantGap   bool
	}{
		{"unreachable binding deduplicates", map[string]document{"required": required}, nil, true},
		{"optional binding", map[string]document{"required": {"optional": true, "effects": required["effects"]}}, nil, false},
		{"conditional transitive cycle", map[string]document{
			"start":    {"effects": []any{map[string]any{"when": map[string]any{"equals": false}, "requireDecisions": []any{"middle"}}}},
			"middle":   {"effects": []any{map[string]any{"activateModules": []any{"extra"}, "requireDecisions": []any{"start"}}}},
			"required": required,
		}, []string{"start"}, false},
		{"alternative template", map[string]document{
			"required": required,
			"selector": {"effects": []any{map[string]any{"selectTemplates": []any{map[string]any{"artifact": "guide", "template": "template.other"}}}}},
		}, nil, false},
		{"same template remains required", map[string]document{
			"required": required,
			"selector": {"effects": []any{map[string]any{"selectTemplates": []any{map[string]any{"artifact": "guide", "template": "template.guide"}}}}},
		}, nil, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			catalog := &Catalog{decisions: tt.decisions, modules: map[string]document{
				"core":  {"supportingGuides": []any{map[string]any{"id": "guide", "template": "template.guide"}}},
				"extra": {"requiredDecisions": []any{"required"}},
			}}
			profile := ResolvedProfile{Modules: []string{"core"}, Decisions: tt.selected}
			gaps := ProfileDecisionGaps(profile, catalog)
			if tt.wantGap {
				want := []ProfileDecisionGap{{Decision: "required", Template: "template.guide", Artifact: "guide", Suggested: "false"}}
				if !slices.Equal(gaps, want) {
					t.Fatalf("gaps = %+v, want %+v", gaps, want)
				}
			} else if len(gaps) != 0 {
				t.Fatalf("unexpected gaps = %+v", gaps)
			}
			profile.Modules = nil
			if gaps := ProfileDecisionGaps(profile, catalog); len(gaps) != 0 {
				t.Fatalf("binding on unselected module = %+v", gaps)
			}
		})
	}
}

func TestProfileDecisionGapsOrderBindingsAndSuggestions(t *testing.T) {
	t.Parallel()
	catalog := &Catalog{
		modules: map[string]document{"core": {"rootBlocks": []any{
			map[string]any{"id": "root.z", "template": "template.z"},
			map[string]any{"id": "root.a", "template": "template.a"},
		}}},
		decisions: map[string]document{
			"a": {"suggestion": true, "default": false, "effects": []any{map[string]any{"renderBindings": []any{
				map[string]any{"artifact": "root.z", "template": "template.z", "token": "a"},
				map[string]any{"artifact": "root.a", "template": "template.a", "token": "a"},
			}}}},
			"z": {"effects": []any{map[string]any{"renderBindings": []any{
				map[string]any{"artifact": "root.a", "template": "template.a", "token": "z"},
			}}}},
		},
	}
	want := []ProfileDecisionGap{
		{Decision: "a", Template: "template.a", Artifact: "root.a", Suggested: "true"},
		{Decision: "a", Template: "template.z", Artifact: "root.z", Suggested: "true"},
		{Decision: "z", Template: "template.a", Artifact: "root.a"},
	}
	if got := ProfileDecisionGaps(ResolvedProfile{Modules: []string{"core"}}, catalog); !slices.Equal(got, want) {
		t.Fatalf("ordered gaps = %+v, want %+v", got, want)
	}
}

func TestProfileDecisionGapAcceptsExistingRenderValues(t *testing.T) {
	t.Parallel()
	catalog, err := LoadEmbeddedCatalog()
	if err != nil {
		t.Fatal(err)
	}
	tests := []struct {
		name      string
		gap       ProfileDecisionGap
		decisions []DecisionValue
		want      bool
	}{
		{"missing incremental", ProfileDecisionGap{Decision: "verification.incremental", Template: "template.guide.agent-instructions", Artifact: "guide.agent-instructions"}, nil, false},
		{"normalized incremental", ProfileDecisionGap{Decision: "verification.incremental", Template: "template.guide.agent-instructions", Artifact: "guide.agent-instructions"}, []DecisionValue{{ID: "verification.incremental", Value: "make verify-incremental"}}, true},
		{"backend provider default", ProfileDecisionGap{Decision: "auth.provider", Template: "template.guide.backend", Artifact: "guide.backend"}, nil, true},
		{"backend HTTP default", ProfileDecisionGap{Decision: "http.contract", Template: "template.guide.backend", Artifact: "guide.backend"}, nil, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := profileDecisionGapHasRenderValue(tt.gap, tt.decisions, catalog); got != tt.want {
				t.Fatalf("render value available = %v, want %v", got, tt.want)
			}
		})
	}
}
