package config

import (
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"gopkg.in/yaml.v3"
)

func TestBuiltinProfilesAreTheRecommendedProfile(t *testing.T) {
	loaded, err := Load(LoadOptions{HomeDir: t.TempDir(), WorkDir: t.TempDir()})
	if err != nil {
		t.Fatal(err)
	}
	for _, category := range RequiredWorkCategories() {
		t.Run(string(category), func(t *testing.T) {
			want, ok := RecommendedProfile(category)
			if !ok {
				t.Fatal("missing Recommended Profile")
			}
			got, err := ResolveProfile(loaded.Config, category, nil)
			if err != nil {
				t.Fatal(err)
			}
			if got.Source != ProfileSourceBuiltIn || got.InheritedFrom != "" || !reflect.DeepEqual(got.Profile, want) {
				t.Fatalf("resolved = %#v, want built-in %#v", got, want)
			}
		})
	}
}

func TestGeneratedConfigRendersTheBuiltinProfiles(t *testing.T) {
	for _, tc := range []struct{ name, content string }{
		{"user", DefaultConfigYAML()}, {"project", DefaultProjectConfigYAML()},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var doc yaml.Node
			if err := yaml.Unmarshal([]byte(tc.content), &doc); err != nil {
				t.Fatal(err)
			}
			root := documentRootMapping(&doc)
			var profiles *yaml.Node
			for i := 0; i < len(root.Content); i += 2 {
				if root.Content[i].Value == "profiles" {
					profiles = root.Content[i+1]
				}
			}
			if profiles == nil {
				t.Fatal("missing profiles")
			}
			var overlay profilesOverlay
			if err := profiles.Decode(&overlay); err != nil {
				t.Fatal(err)
			}
			if len(overlay.entries) != len(RequiredWorkCategories()) {
				t.Fatalf("profiles = %#v", overlay.entries)
			}
			for i, category := range RequiredWorkCategories() {
				if profiles.Content[2*i].Value != string(category) {
					t.Fatalf("category order = %#v", profiles.Content)
				}
				if !reflect.DeepEqual(overlay.entries[category].Profile, Builtin().Profiles[category].Profile) {
					t.Fatalf("%s did not round-trip: %#v", category, overlay.entries[category])
				}
			}
			// Exercise the normal scope loader too, including validation.
			home, work := t.TempDir(), t.TempDir()
			path := filepath.Join(work, projectConfigName)
			if tc.name == "user" {
				path = filepath.Join(home, userConfigRelPath)
			}
			if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
				t.Fatal(err)
			}
			if err := os.WriteFile(path, []byte(tc.content), 0600); err != nil {
				t.Fatal(err)
			}
			loaded, err := Load(LoadOptions{HomeDir: home, WorkDir: work})
			if err != nil {
				t.Fatal(err)
			}
			for _, category := range RequiredWorkCategories() {
				if !reflect.DeepEqual(loaded.Config.Profiles[category].Profile, Builtin().Profiles[category].Profile) {
					t.Fatalf("%s loader changed profile", category)
				}
			}
		})
	}
}

func TestLegacyCodexDefaultIsTheGeneralCodexSelection(t *testing.T) {
	home, work := t.TempDir(), t.TempDir()
	if err := os.WriteFile(filepath.Join(work, projectConfigName), []byte("defaults:\n  agent: codex\n"), 0600); err != nil {
		t.Fatal(err)
	}
	loaded, err := Load(LoadOptions{HomeDir: home, WorkDir: work})
	if err != nil {
		t.Fatal(err)
	}
	want, ok := RecommendedProfile(CategoryGeneral)
	if !ok || want.Preferred.Runtime != "codex" {
		t.Fatal("general must prefer Codex")
	}
	got, err := ResolveProfile(loaded.Config, CategoryGeneral, nil)
	if err != nil {
		t.Fatal(err)
	}
	if got.Profile.Preferred != want.Preferred {
		t.Fatalf("legacy = %#v, want %#v", got.Profile.Preferred, want.Preferred)
	}
	if loaded.Config.Runtimes.Codex.Model != want.Preferred.Model || loaded.Config.Runtimes.Codex.ReasoningEffort != want.Preferred.ReasoningEffort {
		t.Fatal("legacy runtime differs from general")
	}
	if strings.TrimSpace(loaded.Config.Runtimes.Claude.ReasoningEffort) != "" || loaded.Config.Runtimes.Claude.Model != defaultClaudeModel {
		t.Fatal("Claude legacy default changed")
	}
}

func TestBuiltinProfilesDefineNoOptionalCategory(t *testing.T) {
	cfg := Builtin()
	if len(cfg.Profiles) != len(RequiredWorkCategories()) {
		t.Fatalf("profiles = %#v", cfg.Profiles)
	}
	want, _ := RecommendedProfile(CategoryGeneral)
	for _, category := range OptionalWorkCategories() {
		if _, defined := cfg.Profiles[category]; defined {
			t.Fatalf("%s is built in", category)
		}
		got, err := ResolveProfile(cfg, category, nil)
		if err != nil {
			t.Fatal(err)
		}
		if got.Source != ProfileSourceBuiltIn || got.InheritedFrom != CategoryGeneral || !reflect.DeepEqual(got.Profile, want) {
			t.Fatalf("%s inheritance = %#v", category, got)
		}
	}
}
