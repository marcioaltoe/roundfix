package config

import (
	"reflect"
	"testing"
)

func TestCheckRecommendationsReportsCurrentDiffersAndPinned(t *testing.T) {
	t.Parallel()
	cfg := Config{Profiles: builtinProfiles()}
	backend := cfg.Profiles[CategoryBackend]
	backend.Profile.Preferred.Model = "chosen-model"
	backend.Source = ProfileSourceProject
	cfg.Profiles[CategoryBackend] = backend
	review := cfg.Profiles[CategoryReview]
	review.Profile.Preferred.Model = "pinned-model"
	review.Deviation = &ProfileDeviation{From: ModelRecommendationSnapshotVersion, Reason: "quota"}
	cfg.Profiles[CategoryReview] = review
	general := cfg.Profiles[CategoryGeneral]
	general.Deviation = &ProfileDeviation{From: ModelRecommendationSnapshotVersion, Reason: "already matches"}
	cfg.Profiles[CategoryGeneral] = general
	check, err := CheckRecommendations(cfg)
	if err != nil {
		t.Fatal(err)
	}
	if check.Snapshot != ModelRecommendationSnapshotVersion {
		t.Fatalf("snapshot = %q", check.Snapshot)
	}
	want := []RecommendationStatus{RecommendationCurrent, RecommendationDiffers, RecommendationCurrent, RecommendationCurrent, RecommendationPinned}
	for i, row := range check.Categories {
		if row.Status != want[i] {
			t.Errorf("%s status = %s, want %s", row.Category, row.Status, want[i])
		}
		profile, _ := RecommendedProfile(row.Category)
		if !reflect.DeepEqual(row.Recommended, profile) || !reflect.DeepEqual(row.Configured, cfg.Profiles[row.Category].Profile) || row.Source != cfg.Profiles[row.Category].Source || !reflect.DeepEqual(row.Deviation, cfg.Profiles[row.Category].Deviation) {
			t.Errorf("incorrect row: %+v", row)
		}
	}
	if len(check.Categories) != len(want) {
		t.Fatalf("categories = %d", len(check.Categories))
	}
	// Returned profiles and deviations must not alias the caller's configuration.
	check.Categories[0].Configured.Fallbacks[0].Model = "mutated"
	check.Categories[0].Deviation.Reason = "mutated"
	if cfg.Profiles[CategoryGeneral].Deviation.Reason != "already matches" || cfg.Profiles[CategoryGeneral].Profile.Fallbacks[0].Model == "mutated" {
		t.Fatal("check aliases input")
	}
}

func TestCheckRecommendationsEndsADeviationOfAnotherSnapshot(t *testing.T) {
	t.Parallel()
	cfg := Config{Profiles: builtinProfiles()}
	entry := cfg.Profiles[CategoryBackend]
	entry.Profile.Preferred.Model = "chosen-model"
	entry.Deviation = &ProfileDeviation{From: "2000-01-01", Reason: "older choice"}
	cfg.Profiles[CategoryBackend] = entry
	check, err := CheckRecommendations(cfg)
	if err != nil {
		t.Fatal(err)
	}
	row := check.Categories[1]
	if row.Status != RecommendationDiffers || row.Deviation == nil || *row.Deviation != *entry.Deviation {
		t.Fatalf("expired deviation = %+v", row)
	}
}

func TestCheckRecommendationsReportsBuiltinsAsCurrent(t *testing.T) {
	t.Parallel()
	check, err := CheckRecommendations(Config{})
	if err != nil {
		t.Fatal(err)
	}
	if len(check.Categories) != len(RequiredWorkCategories()) {
		t.Fatalf("categories = %+v", check.Categories)
	}
	for i, row := range check.Categories {
		if row.Status != RecommendationCurrent || row.Category != RequiredWorkCategories()[i] || row.Source != ProfileSourceBuiltIn {
			t.Errorf("built-in row = %+v", row)
		}
	}
}

func TestCheckRecommendationsLeavesOutUndefinedOptionalCategories(t *testing.T) {
	t.Parallel()
	cfg := Config{Profiles: builtinProfiles()}
	profile, _ := RecommendedProfile(CategoryDocs)
	cfg.Profiles[CategoryDocs] = ProfileEntry{Profile: profile, Source: ProfileSourceUser}
	check, err := CheckRecommendations(cfg)
	if err != nil {
		t.Fatal(err)
	}
	want := []WorkCategory{CategoryGeneral, CategoryBackend, CategoryFrontend, CategoryDocs, CategoryQA, CategoryReview}
	got := make([]WorkCategory, 0, len(check.Categories))
	for _, row := range check.Categories {
		got = append(got, row.Category)
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("categories = %v, want %v", got, want)
	}
}

func TestCheckRecommendationsComparesTheWholeOrderedFallbackChain(t *testing.T) {
	t.Parallel()
	for _, name := range []string{"preferred", "fallback model", "fallback runtime", "fallback effort", "extra fallback", "fallback order"} {
		t.Run(name, func(t *testing.T) {
			cfg := Config{Profiles: builtinProfiles()}
			entry := cfg.Profiles[CategoryBackend]
			switch name {
			case "preferred":
				entry.Profile.Preferred.ReasoningEffort = "xhigh"
			case "fallback model":
				entry.Profile.Fallbacks[0].Model = "other"
			case "fallback runtime":
				entry.Profile.Fallbacks[0].Runtime = "codex"
			case "fallback effort":
				entry.Profile.Fallbacks[0].ReasoningEffort = "max"
			case "extra fallback":
				entry.Profile.Fallbacks = append(entry.Profile.Fallbacks, AgentSelection{Runtime: "codex", Model: "backup", ReasoningEffort: "high"})
			case "fallback order":
				entry.Profile.Preferred, entry.Profile.Fallbacks[0] = entry.Profile.Fallbacks[0], entry.Profile.Preferred
			}
			cfg.Profiles[CategoryBackend] = entry
			check, err := CheckRecommendations(cfg)
			if err != nil {
				t.Fatal(err)
			}
			if check.Categories[1].Status != RecommendationDiffers {
				t.Fatalf("%s not detected", name)
			}
		})
	}
}

func TestCheckRecommendationsRejectsAnUnresolvableRequiredCategory(t *testing.T) {
	t.Parallel()
	if _, err := CheckRecommendations(Config{Profiles: Profiles{}}); err == nil {
		t.Fatal("missing required profiles accepted")
	}
}
