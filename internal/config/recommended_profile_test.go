package config_test

import (
	"reflect"
	"testing"

	"roundfix/internal/agent"
	config "roundfix/internal/config"
)

func TestRecommendedProfileCoversEveryWorkCategory(t *testing.T) {
	implementation := config.AgentSelectionProfile{Preferred: config.AgentSelection{Runtime: "codex", Model: "gpt-6.1-sol", ReasoningEffort: "high"}, Fallbacks: []config.AgentSelection{{Runtime: "claude", Model: "opus", ReasoningEffort: "high"}}}
	bounded := config.AgentSelection{Runtime: "codex", Model: "gpt-5.6-luna", ReasoningEffort: "max"}
	want := map[config.WorkCategory]config.AgentSelectionProfile{
		config.CategoryGeneral: implementation, config.CategoryBackend: implementation,
		config.CategoryData: implementation, config.CategoryInfra: implementation,
		config.CategoryTest: implementation, config.CategoryQA: implementation,
		config.CategoryFrontend: {Preferred: config.AgentSelection{Runtime: "claude", Model: "opus", ReasoningEffort: "high"}, Fallbacks: []config.AgentSelection{{Runtime: "codex", Model: "gpt-6.1-sol", ReasoningEffort: "xhigh"}}},
		config.CategoryReview:   {Preferred: bounded, Fallbacks: []config.AgentSelection{{Runtime: "codex", Model: "gpt-6.1-sol", ReasoningEffort: "high"}}},
		config.CategoryDocs:     {Preferred: bounded, Fallbacks: []config.AgentSelection{{Runtime: "claude", Model: "sonnet", ReasoningEffort: "high"}}},
		config.CategoryChore:    {Preferred: bounded, Fallbacks: []config.AgentSelection{{Runtime: "claude", Model: "sonnet", ReasoningEffort: "high"}}},
	}
	for _, category := range config.AllWorkCategories() {
		profile, ok := config.RecommendedProfile(category)
		if !ok || !reflect.DeepEqual(profile, want[category]) {
			t.Fatalf("%s: got %+v (%v), want %+v", category, profile, ok, want[category])
		}
		profile.Fallbacks[0].Model = "changed"
		fresh, _ := config.RecommendedProfile(category)
		if !reflect.DeepEqual(fresh, want[category]) {
			t.Fatalf("%s: returned profile aliases table", category)
		}
	}
	if _, ok := config.RecommendedProfile("unknown"); ok {
		t.Fatal("unknown category has a profile")
	}
}

func TestRecommendedProfileFallbackChangesRuntimeExceptReview(t *testing.T) {
	for _, category := range config.AllWorkCategories() {
		profile, ok := config.RecommendedProfile(category)
		if !ok || len(profile.Fallbacks) == 0 {
			t.Fatalf("%s: missing profile or fallback", category)
		}
		if category == config.CategoryReview {
			for _, selection := range append([]config.AgentSelection{profile.Preferred}, profile.Fallbacks...) {
				if selection.Runtime != "codex" {
					t.Fatalf("review must stay on Codex: %+v", selection)
				}
			}
		} else if profile.Preferred.Runtime == profile.Fallbacks[0].Runtime {
			t.Fatalf("%s: first fallback must change runtime", category)
		}
	}
}

func TestRecommendedModelsAreInTheModelCatalog(t *testing.T) {
	for _, category := range config.AllWorkCategories() {
		profile, ok := config.RecommendedProfile(category)
		if !ok {
			t.Fatalf("%s: missing profile", category)
		}
		for _, selection := range append([]config.AgentSelection{profile.Preferred}, profile.Fallbacks...) {
			found := false
			for _, model := range agent.ModelCatalog(selection.Runtime) {
				if model.Value == selection.Model {
					found = true
				}
			}
			if !found {
				t.Errorf("%s: recommended selection absent from catalog: %+v", category, selection)
			}
		}
	}
}

func TestModelRecommendationsListTheRecommendedProfileInOrder(t *testing.T) {
	for _, category := range config.AllWorkCategories() {
		profile, _ := config.RecommendedProfile(category)
		selections := append([]config.AgentSelection{profile.Preferred}, profile.Fallbacks...)
		rows, ok := config.ModelRecommendations(category)
		if !ok || len(rows) != len(selections) {
			t.Fatalf("%s: rows %+v", category, rows)
		}
		for i, row := range rows {
			role := config.RecommendationFallback
			if i == 0 {
				role = config.RecommendationPreferred
			}
			if row.Category != category || row.Rank != i+1 || row.Role != role || row.Selection != selections[i] || row.SourceAsOf != config.ModelRecommendationSnapshotVersion || row.Rationale == "" {
				t.Fatalf("%s row %d: %+v", category, i, row)
			}
		}
		rows[0].Rationale = "changed"
		fresh, _ := config.ModelRecommendations(category)
		if fresh[0].Rationale == "changed" {
			t.Fatal("rows alias recommendation table")
		}
	}
	if _, ok := config.ModelRecommendations("unknown"); ok {
		t.Fatal("unknown category has rows")
	}
}
