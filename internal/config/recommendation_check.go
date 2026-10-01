package config

import (
	"fmt"
	"slices"
)

type RecommendationStatus string

const (
	RecommendationCurrent RecommendationStatus = "current"
	RecommendationDiffers RecommendationStatus = "differs"
	RecommendationPinned  RecommendationStatus = "pinned"
)

type CategoryRecommendation struct {
	Category    WorkCategory
	Status      RecommendationStatus
	Source      ProfileSource
	Configured  AgentSelectionProfile
	Recommended AgentSelectionProfile
	Deviation   *ProfileDeviation
}

type RecommendationCheck struct {
	Snapshot   string
	Categories []CategoryRecommendation
}

// CheckRecommendations compares configured categories with the shipped snapshot.
// It reads only config and the recommendation, without IO or Agent Sessions.
func CheckRecommendations(config Config) (RecommendationCheck, error) {
	check := RecommendationCheck{Snapshot: ModelRecommendationSnapshotVersion}
	for _, category := range ConfiguredWorkCategories(config) {
		resolved, err := ResolveProfile(config, category, nil)
		if err != nil {
			return RecommendationCheck{}, err
		}
		recommended, ok := RecommendedProfile(category)
		if !ok {
			return RecommendationCheck{}, fmt.Errorf("recommendations for Agent Work Category %q are not configured", category)
		}
		status := RecommendationDiffers
		if resolved.Profile.Preferred == recommended.Preferred && slices.Equal(resolved.Profile.Fallbacks, recommended.Fallbacks) {
			status = RecommendationCurrent
		} else if resolved.Deviation != nil && resolved.Deviation.From == ModelRecommendationSnapshotVersion {
			status = RecommendationPinned
		}
		check.Categories = append(check.Categories, CategoryRecommendation{
			Category: category, Status: status, Source: resolved.Source,
			Configured: resolved.Profile, Recommended: recommended, Deviation: resolved.Deviation,
		})
	}
	return check, nil
}
