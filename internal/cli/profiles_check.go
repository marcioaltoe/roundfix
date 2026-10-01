package cli

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"strings"

	roundconfig "roundfix/internal/config"
)

const profilesCheckSchema = "roundfix/profiles-check/v1"

type profilesCheckProfile struct {
	Preferred roundconfig.AgentSelection   `json:"preferred"`
	Fallbacks []roundconfig.AgentSelection `json:"fallbacks"`
}

type profilesCheckCategory struct {
	Category    roundconfig.WorkCategory         `json:"category"`
	Status      roundconfig.RecommendationStatus `json:"status"`
	Source      roundconfig.ProfileSource        `json:"source"`
	Configured  profilesCheckProfile             `json:"configured"`
	Recommended profilesCheckProfile             `json:"recommended"`
	Deviation   *roundconfig.ProfileDeviation    `json:"deviation,omitempty"`
}

type profilesCheckResponse struct {
	Schema     string                  `json:"schema"`
	Snapshot   string                  `json:"snapshot"`
	Current    int                     `json:"current"`
	Differ     int                     `json:"differ"`
	Pinned     int                     `json:"pinned"`
	Categories []profilesCheckCategory `json:"categories"`
}

func runProfilesCheckCommand(args []string, stdout, stderr io.Writer, environment commandEnvironment) int {
	if commandWantsHelp(args) {
		fmt.Fprint(stdout, commandUsage("profiles check"))
		return exitOK
	}
	fs := flag.NewFlagSet("profiles check", flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	jsonOutput := fs.Bool("json", false, "Print roundfix/profiles-check/v1 JSON")
	if err := fs.Parse(args); err != nil {
		printProfilesFailure(validationError{message: err.Error()}, stderr)
		return exitPreflight
	}
	if remaining := fs.Args(); len(remaining) > 0 {
		printProfilesFailure(validationError{message: fmt.Sprintf("unexpected argument %q", remaining[0])}, stderr)
		return exitPreflight
	}
	loaded, err := loadCommandConfig(environment, stderr)
	if err != nil {
		printProfilesFailure(err, stderr)
		return exitPreflight
	}
	check, err := roundconfig.CheckRecommendations(loaded.Config)
	if err != nil {
		printProfilesFailure(err, stderr)
		return exitPreflight
	}
	if *jsonOutput {
		encoder := json.NewEncoder(stdout)
		encoder.SetIndent("", "  ")
		if err := encoder.Encode(buildProfilesCheckResponse(check)); err != nil {
			printProfilesFailure(err, stderr)
			return exitRunFailed
		}
	} else {
		printRecommendationCheckText(check, stdout)
	}
	return exitOK
}

func recommendationCounts(check roundconfig.RecommendationCheck) (current, differ, pinned int) {
	for _, row := range check.Categories {
		switch row.Status {
		case roundconfig.RecommendationCurrent:
			current++
		case roundconfig.RecommendationDiffers:
			differ++
		case roundconfig.RecommendationPinned:
			pinned++
		}
	}
	return
}

func buildProfilesCheckResponse(check roundconfig.RecommendationCheck) profilesCheckResponse {
	current, differ, pinned := recommendationCounts(check)
	response := profilesCheckResponse{Schema: profilesCheckSchema, Snapshot: check.Snapshot, Current: current, Differ: differ, Pinned: pinned, Categories: make([]profilesCheckCategory, 0, len(check.Categories))}
	for _, row := range check.Categories {
		response.Categories = append(response.Categories, profilesCheckCategory{
			Category: row.Category, Status: row.Status, Source: row.Source,
			Configured:  profilesCheckProfile{Preferred: row.Configured.Preferred, Fallbacks: row.Configured.Fallbacks},
			Recommended: profilesCheckProfile{Preferred: row.Recommended.Preferred, Fallbacks: row.Recommended.Fallbacks},
			Deviation:   row.Deviation,
		})
	}
	return response
}

// printRecommendationCheckText is shared by the check command and notices.
func printRecommendationCheckText(check roundconfig.RecommendationCheck, stdout io.Writer) {
	for _, row := range check.Categories {
		switch row.Status {
		case roundconfig.RecommendationPinned:
			fmt.Fprintf(stdout, "recommendations: %s pinned against snapshot %s: %s\n", row.Category, check.Snapshot, row.Deviation.Reason)
		case roundconfig.RecommendationDiffers:
			fmt.Fprintf(stdout, "recommendations: %s differs from snapshot %s", row.Category, check.Snapshot)
			if row.Deviation != nil {
				fmt.Fprintf(stdout, " (its deviation was declared against %s)", row.Deviation.From)
			}
			fmt.Fprintln(stdout)
			fmt.Fprintf(stdout, "  configured:  %s (%s)\n", formatRecommendationProfile(row.Configured), row.Source)
			fmt.Fprintf(stdout, "  recommended: %s\n", formatRecommendationProfile(row.Recommended))
		}
	}
	current, differ, pinned := recommendationCounts(check)
	fmt.Fprintf(stdout, "recommendations: snapshot %s; %d current, %d differ, %d pinned\n", check.Snapshot, current, differ, pinned)
	if differ > 0 {
		fmt.Fprintln(stdout, "recommendations: adopt with `roundfix profiles check --apply --scope user|project`, or declare a deviation")
	}
}

func formatRecommendationProfile(profile roundconfig.AgentSelectionProfile) string {
	selections := make([]string, 0, len(profile.Fallbacks)+1)
	selections = append(selections, formatProfileSelection(profile.Preferred))
	for _, fallback := range profile.Fallbacks {
		selections = append(selections, formatProfileSelection(fallback))
	}
	return strings.Join(selections, ", then ")
}
