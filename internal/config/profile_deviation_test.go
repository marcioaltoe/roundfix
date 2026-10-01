package config

import (
	"context"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

const deviationProfile = `profiles:
  general:
    preferred: {runtime: codex, model: chosen-model, reasoning_effort: high}
    fallbacks:
      - {runtime: codex, model: backup-model, reasoning_effort: ""}
`
const validDeviation = `    deviation:
      from: 2026-09-30
      reason: "  Keep the validated model  "
`

func deviationWorkspace(t *testing.T) (string, string) {
	t.Helper()
	home, repo := t.TempDir(), t.TempDir()
	mustMkdir(t, filepath.Join(repo, ".git"))
	mustMkdir(t, filepath.Join(home, ".roundfix"))
	return home, repo
}

func TestProfileDeviationLoadsFromUserAndProjectConfig(t *testing.T) {
	t.Parallel()
	for _, scope := range []string{"user", "project"} {
		t.Run(scope, func(t *testing.T) {
			home, repo := deviationWorkspace(t)
			path := filepath.Join(home, ".roundfix", "config.yml")
			source := ProfileSourceUser
			if scope == "project" {
				path = filepath.Join(repo, ".roundfixrc.yml")
				source = ProfileSourceProject
			}
			mustWrite(t, path, deviationProfile+validDeviation)
			loaded, err := Load(LoadOptions{HomeDir: home, WorkDir: repo})
			if err != nil {
				t.Fatal(err)
			}
			want := ProfileDeviation{From: "2026-09-30", Reason: "Keep the validated model"}
			for _, category := range []WorkCategory{CategoryGeneral, CategoryData} {
				resolved, err := ResolveProfile(loaded.Config, category, nil)
				if err != nil {
					t.Fatal(err)
				}
				if resolved.Source != source || resolved.Deviation == nil || *resolved.Deviation != want {
					t.Fatalf("resolved = %+v", resolved)
				}
				resolved.Deviation.Reason = "mutated copy"
				if *loaded.Config.Profiles[CategoryGeneral].Deviation != want {
					t.Fatal("resolution aliases configured deviation")
				}
			}
			override := AgentSelection{Runtime: "codex", Model: "override", ReasoningEffort: ""}
			resolved, err := ResolveProfile(loaded.Config, CategoryGeneral, &override)
			if err != nil || resolved.Deviation != nil {
				t.Fatalf("invocation deviation = %+v, err = %v", resolved.Deviation, err)
			}
			if scope == "user" {
				mustWrite(t, filepath.Join(repo, ".roundfixrc.yml"), deviationProfile)
				replaced, err := Load(LoadOptions{HomeDir: home, WorkDir: repo})
				if err != nil {
					t.Fatal(err)
				}
				if replaced.Config.Profiles[CategoryGeneral].Deviation != nil {
					t.Fatal("project replacement retained user deviation")
				}
			}
		})
	}
}

func TestProfileDeviationRejectsMalformedValues(t *testing.T) {
	t.Parallel()
	cases := []struct{ name, value, path string }{
		{"unknown deviation key", "{from: 2026-09-30, reason: keep, extra: true}", ".extra"},
		{"missing from", "{reason: keep}", ".from"},
		{"missing reason", "{from: 2026-09-30}", ".reason"},
		{"empty reason", "{from: 2026-09-30, reason: '  '}", ".reason"},
		{"not a calendar date", "{from: 2026-02-30, reason: keep}", ".from"},
		{"not a mapping", "keep", ""},
		{"duplicate key", "{from: 2026-09-30, reason: keep, reason: again}", ".reason"},
		{"non string reason", "{from: 2026-09-30, reason: 42}", ".reason"},
		{"non scalar from", "{from: [], reason: keep}", ".from"},
		{"date time", "{from: 2026-09-30T00:00:00Z, reason: keep}", ".from"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			content := deviationProfile + "    deviation: " + tc.value + "\n"
			home, repo := deviationWorkspace(t)
			mustWrite(t, filepath.Join(repo, ".roundfixrc.yml"), content)
			_, err := Load(LoadOptions{HomeDir: home, WorkDir: repo})
			if err == nil || !strings.Contains(err.Error(), "profiles.general.deviation"+tc.path) {
				t.Fatalf("expected path refusal, got %v", err)
			}
			if _, err := ParseProfilesFragment([]byte(content)); err == nil || !strings.Contains(err.Error(), "profiles.general.deviation"+tc.path) {
				t.Fatalf("fragment refusal = %v", err)
			}
		})
	}
}

func TestProfileWithoutDeviationLoadsUnchanged(t *testing.T) {
	t.Parallel()
	home, repo := deviationWorkspace(t)
	mustWrite(t, filepath.Join(repo, ".roundfixrc.yml"), deviationProfile)
	loaded, err := Load(LoadOptions{HomeDir: home, WorkDir: repo})
	if err != nil {
		t.Fatal(err)
	}
	fragment, err := ParseProfilesFragment([]byte(deviationProfile))
	if err != nil {
		t.Fatal(err)
	}
	resolved, err := ResolveProfile(loaded.Config, CategoryGeneral, nil)
	if err != nil {
		t.Fatal(err)
	}
	if resolved.Deviation != nil || !reflect.DeepEqual(resolved.Profile, fragment[CategoryGeneral].Profile) {
		t.Fatalf("profile changed: %+v", resolved)
	}
	for _, content := range []string{"", "defaults:\n  agent: codex\n"} {
		mustWrite(t, filepath.Join(repo, ".roundfixrc.yml"), content)
		loaded, err := Load(LoadOptions{HomeDir: home, WorkDir: repo})
		if err != nil {
			t.Fatal(err)
		}
		for category := range loaded.Config.Profiles {
			resolved, err := ResolveProfile(loaded.Config, category, nil)
			if err != nil || resolved.Deviation != nil {
				t.Fatalf("built-in/legacy deviation: %+v, %v", resolved, err)
			}
		}
	}
}

func TestProfilesFragmentPersistsTheDeviation(t *testing.T) {
	t.Parallel()
	home, repo := deviationWorkspace(t)
	fragment, err := ParseProfilesFragment([]byte(deviationProfile + validDeviation))
	if err != nil {
		t.Fatal(err)
	}
	proposal, err := PrepareProfilesConfig(context.Background(), ProfileConfigOptions{Scope: "project", HomeDir: home, WorkDir: repo, Profiles: fragment})
	if err != nil {
		t.Fatal(err)
	}
	want := *fragment[CategoryGeneral].Deviation
	result := proposal.Result()
	if result.Changes.Changes[0].Deviation == nil || *result.Changes.Changes[0].Deviation != want {
		t.Fatal("change set lost deviation")
	}
	fragment[CategoryGeneral].Deviation.Reason = "mutated input"
	result.Profiles[CategoryGeneral].Deviation.Reason = "mutated result"
	result.Changes.Changes[0].Deviation.Reason = "mutated changes"
	result, err = PersistProfilesConfig(context.Background(), proposal)
	if err != nil {
		t.Fatal(err)
	}
	if !result.Changed || *result.Profiles[CategoryGeneral].Deviation != want || *result.Changes.Changes[0].Deviation != want {
		t.Fatal("proposal aliases deviations")
	}
	loaded, err := Load(LoadOptions{HomeDir: home, WorkDir: repo})
	if err != nil {
		t.Fatal(err)
	}
	got := loaded.Config.Profiles[CategoryGeneral].Deviation
	if got == nil || *got != want {
		t.Fatalf("written deviation = %+v, want %+v", got, want)
	}
}

func TestReplacingAProfileDropsItsDeviation(t *testing.T) {
	t.Parallel()
	home, repo := deviationWorkspace(t)
	mustWrite(t, filepath.Join(repo, ".roundfixrc.yml"), deviationProfile+validDeviation)
	fragment, err := ParseProfilesFragment([]byte(deviationProfile))
	if err != nil {
		t.Fatal(err)
	}
	_, err = WriteProfilesConfig(context.Background(), ProfileConfigOptions{Scope: "project", HomeDir: home, WorkDir: repo, Profiles: fragment})
	if err != nil {
		t.Fatal(err)
	}
	loaded, err := Load(LoadOptions{HomeDir: home, WorkDir: repo})
	if err != nil {
		t.Fatal(err)
	}
	if loaded.Config.Profiles[CategoryGeneral].Deviation != nil {
		t.Fatal("replacement retained deviation")
	}
}
