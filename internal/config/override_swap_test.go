package config

import (
	"reflect"
	"strings"
	"testing"
)

func TestAnOverrideEqualToAFallbackSwapsWithThePreferred(t *testing.T) {
	t.Parallel()
	preferred := AgentSelection{Runtime: "codex", Model: "preferred", ReasoningEffort: "high"}
	fallbacks := []AgentSelection{
		{Runtime: "claude", Model: "first", ReasoningEffort: "high"},
		{Runtime: "codex", Model: "second", ReasoningEffort: "max"},
		{Runtime: "claude", Model: "third", ReasoningEffort: ""},
	}
	for index, override := range fallbacks {
		t.Run(override.Model, func(t *testing.T) {
			profile := AgentSelectionProfile{Preferred: preferred, Fallbacks: append([]AgentSelection(nil), fallbacks...)}
			config := Config{Profiles: Profiles{CategoryGeneral: {Profile: profile, Source: ProfileSourceProject}}}
			resolved, err := ResolveProfile(config, CategoryData, &override)
			if err != nil {
				t.Fatal(err)
			}
			wantFallbacks := append([]AgentSelection(nil), fallbacks...)
			wantFallbacks[index] = preferred
			want := ResolvedProfile{Category: CategoryData, Source: ProfileSourceInvocation, InheritedFrom: CategoryGeneral,
				Profile: AgentSelectionProfile{Preferred: override, Fallbacks: wantFallbacks}}
			if !reflect.DeepEqual(resolved, want) {
				t.Fatalf("resolved = %+v, want %+v", resolved, want)
			}
			if !reflect.DeepEqual(config.Profiles[CategoryGeneral].Profile, AgentSelectionProfile{Preferred: preferred, Fallbacks: fallbacks}) {
				t.Fatal("override mutated the configured profile")
			}
		})
	}
	t.Run("configured duplicate remains refused", func(t *testing.T) {
		config := Config{Profiles: Profiles{CategoryGeneral: {Profile: AgentSelectionProfile{
			Preferred: preferred, Fallbacks: []AgentSelection{preferred, fallbacks[0]},
		}}}}
		for _, override := range []AgentSelection{preferred, fallbacks[0], {Runtime: "codex", Model: "distinct", ReasoningEffort: "high"}} {
			if _, err := ResolveProfile(config, CategoryGeneral, &override); err == nil || !strings.Contains(err.Error(), "duplicate Agent Selection") {
				t.Fatalf("override %+v: error = %v, want duplicate refusal", override, err)
			}
		}
	})
}

func TestAnOverrideDistinctFromTheChainKeepsTheFallbacks(t *testing.T) {
	t.Parallel()
	profile := AgentSelectionProfile{
		Preferred: AgentSelection{Runtime: "codex", Model: "preferred", ReasoningEffort: "high"},
		Fallbacks: []AgentSelection{
			{Runtime: "claude", Model: "first", ReasoningEffort: "high"},
			{Runtime: "codex", Model: "second", ReasoningEffort: "max"},
		},
	}
	config := Config{Profiles: Profiles{CategoryBackend: {Profile: profile, Source: ProfileSourceUser}}}
	original := AgentSelectionProfile{Preferred: profile.Preferred, Fallbacks: append([]AgentSelection(nil), profile.Fallbacks...)}
	override := AgentSelection{Runtime: "claude", Model: "distinct", ReasoningEffort: ""}
	resolved, err := ResolveProfile(config, CategoryBackend, &override)
	if err != nil {
		t.Fatal(err)
	}
	want := ResolvedProfile{Category: CategoryBackend, Source: ProfileSourceInvocation,
		Profile: AgentSelectionProfile{Preferred: override, Fallbacks: profile.Fallbacks}}
	if !reflect.DeepEqual(resolved, want) {
		t.Fatalf("resolved = %+v, want %+v", resolved, want)
	}
	if !reflect.DeepEqual(config.Profiles[CategoryBackend].Profile, original) {
		t.Fatal("override mutated the configured profile")
	}
}
