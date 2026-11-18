package cli

import (
	"fmt"
	"io"

	roundconfig "roundfix/internal/config"
)

type updateCheckRecommendations struct {
	Snapshot   string                          `json:"snapshot"`
	Categories []updateCheckRecommendationItem `json:"categories"`
}

type updateCheckRecommendationItem struct {
	Category    roundconfig.WorkCategory      `json:"category"`
	State       string                        `json:"state"`
	Source      roundconfig.ProfileSource     `json:"source"`
	Configured  profilesCheckProfile          `json:"configured"`
	Recommended profilesCheckProfile          `json:"recommended"`
	Deviation   *roundconfig.ProfileDeviation `json:"deviation,omitempty"`
	Repin       string                        `json:"repin,omitempty"`
	Adopt       string                        `json:"adopt"`
}

func updateCheckProfileRecommendations(loaded roundconfig.Loaded) *updateCheckRecommendations {
	if loaded.RepositoryVersionErr != nil {
		return nil
	}
	check, err := roundconfig.CheckRecommendations(loaded.Config)
	if err != nil {
		// The optional comparison must not change the Update Check's report contract.
		return nil
	}
	report := &updateCheckRecommendations{Snapshot: check.Snapshot}
	for _, row := range check.Categories {
		if row.Status != roundconfig.RecommendationDiffers {
			continue
		}
		item := updateCheckRecommendationItem{
			Category: row.Category, State: "differs", Source: row.Source,
			Configured:  profilesCheckProfile{Preferred: row.Configured.Preferred, Fallbacks: row.Configured.Fallbacks},
			Recommended: profilesCheckProfile{Preferred: row.Recommended.Preferred, Fallbacks: row.Recommended.Fallbacks},
			Adopt:       "roundfix profiles check --apply --scope " + string(row.Source),
		}
		if row.Deviation != nil {
			item.State = "pinned_to_older_snapshot"
			item.Deviation = row.Deviation
			item.Repin = fmt.Sprintf("profiles.%s.deviation.from: %s", row.Category, check.Snapshot)
		}
		report.Categories = append(report.Categories, item)
	}
	if len(report.Categories) == 0 {
		return nil
	}
	return report
}

func printUpdateCheckRecommendations(stdout io.Writer, report *updateCheckRecommendations) {
	if report == nil {
		return
	}
	for _, item := range report.Categories {
		recommended := formatRecommendationProfile(roundconfig.AgentSelectionProfile{Preferred: item.Recommended.Preferred, Fallbacks: item.Recommended.Fallbacks})
		if item.State == "pinned_to_older_snapshot" {
			label := "Project"
			if item.Source == roundconfig.ProfileSourceUser {
				label = "User"
			}
			fmt.Fprintf(stdout, "Profiles: %s (%s) is pinned to an older snapshot (%s: %s); recommended %s as of %s; re-pin with %s in %s Config, or adopt with %s.\n", item.Category, item.Source, item.Deviation.From, updateCheckLine(item.Deviation.Reason, 300), recommended, report.Snapshot, item.Repin, label, item.Adopt)
			continue
		}
		configured := formatRecommendationProfile(roundconfig.AgentSelectionProfile{Preferred: item.Configured.Preferred, Fallbacks: item.Configured.Fallbacks})
		fmt.Fprintf(stdout, "Profiles: %s (%s) differs from the Recommended Profile of %s: configured %s; recommended %s; adopt with %s.\n", item.Category, item.Source, report.Snapshot, configured, recommended, item.Adopt)
	}
}
