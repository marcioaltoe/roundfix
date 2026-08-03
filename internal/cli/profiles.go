package cli

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"strings"

	"roundfix/internal/app"
	roundconfig "roundfix/internal/config"
)

const profilesShowSchema = "roundfix/profiles/v1"

type profilesShowRequest struct {
	category string
	json     bool
}

type profilesShowResponse struct {
	Schema   string                `json:"schema"`
	Profiles []profilesShowProfile `json:"profiles"`
}

type profilesShowProfile struct {
	Category             roundconfig.WorkCategory           `json:"category"`
	Source               roundconfig.ProfileSource          `json:"source"`
	InheritedFrom        roundconfig.WorkCategory           `json:"inherited_from"`
	RecommendationSource roundconfig.WorkCategory           `json:"recommendation_source"`
	Preferred            roundconfig.AgentSelection         `json:"preferred"`
	Fallbacks            []roundconfig.AgentSelection       `json:"fallbacks"`
	Recommendations      []profilesShowRecommendationOutput `json:"recommendations"`
}

type profilesShowRecommendationOutput struct {
	roundconfig.ModelRecommendation
	UnavailableReason string `json:"unavailable_reason,omitempty"`
}

func runProfilesCommand(ctx context.Context, args []string, stdout, stderr io.Writer, environment commandEnvironment) int {
	if len(args) == 0 {
		fmt.Fprint(stdout, commandUsage("profiles"))
		return exitOK
	}
	if len(args) == 1 && commandWantsHelp(args) {
		fmt.Fprint(stdout, commandUsage("profiles"))
		return exitOK
	}
	switch args[0] {
	case "show":
		return runProfilesShowCommand(args[1:], stdout, stderr, environment)
	case "configure":
		return runProfilesConfigureCommand(ctx, args[1:], stdout, stderr, environment)
	case "validate":
		return runProfilesValidateCommand(ctx, args[1:], stdout, stderr, environment)
	default:
		printProfilesFailure(validationError{message: fmt.Sprintf("unknown profiles command %q", args[0])}, stderr)
		return exitPreflight
	}
}

func runProfilesShowCommand(args []string, stdout, stderr io.Writer, environment commandEnvironment) int {
	if commandWantsHelp(args) {
		fmt.Fprint(stdout, commandUsage("profiles show"))
		return exitOK
	}
	req, err := parseProfilesShowCommand(args)
	if err != nil {
		printProfilesFailure(err, stderr)
		return exitPreflight
	}
	categories, err := profilesShowCategories(req)
	if err != nil {
		printProfilesFailure(err, stderr)
		return exitPreflight
	}
	loaded, err := loadCommandConfig(environment, stderr)
	if err != nil {
		printProfilesFailure(err, stderr)
		return exitPreflight
	}
	response, err := buildProfilesShowResponse(loaded.Config, categories)
	if err != nil {
		printProfilesFailure(err, stderr)
		return exitPreflight
	}
	if req.json {
		encoder := json.NewEncoder(stdout)
		encoder.SetIndent("", "  ")
		if err := encoder.Encode(response); err != nil {
			printProfilesFailure(err, stderr)
			return exitRunFailed
		}
		return exitOK
	}
	printProfilesShowText(response, stdout)
	return exitOK
}

func parseProfilesShowCommand(args []string) (profilesShowRequest, error) {
	req := profilesShowRequest{}
	fs := flag.NewFlagSet("profiles show", flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	fs.StringVar(&req.category, "category", "", "Agent Work Category to show")
	fs.BoolVar(&req.json, "json", false, "Print roundfix/profiles/v1 JSON")
	if err := fs.Parse(args); err != nil {
		return req, validationError{message: err.Error()}
	}
	if remaining := fs.Args(); len(remaining) > 0 {
		return req, validationError{message: fmt.Sprintf("unexpected argument %q", remaining[0])}
	}
	req.category = strings.TrimSpace(req.category)
	return req, nil
}

func profilesShowCategories(req profilesShowRequest) ([]roundconfig.WorkCategory, error) {
	if req.category == "" {
		return roundconfig.AllWorkCategories(), nil
	}
	category, ok := roundconfig.ParseWorkCategory(req.category)
	if !ok {
		return nil, validationError{message: fmt.Sprintf("unknown profile category %q; supported values: %s", req.category, profilesCategoryList())}
	}
	return []roundconfig.WorkCategory{category}, nil
}

func buildProfilesShowResponse(config roundconfig.Config, categories []roundconfig.WorkCategory) (profilesShowResponse, error) {
	return buildProfilesShowResponseWithAvailability(config, categories, nil)
}

func buildProfilesShowResponseWithAvailability(config roundconfig.Config, categories []roundconfig.WorkCategory, unavailable map[roundconfig.AgentSelection]string) (profilesShowResponse, error) {
	response := profilesShowResponse{
		Schema:   profilesShowSchema,
		Profiles: make([]profilesShowProfile, 0, len(categories)),
	}
	for _, category := range categories {
		profile, err := roundconfig.ResolveProfile(config, category, nil)
		if err != nil {
			return profilesShowResponse{}, err
		}
		recommendations, recommendationSource, ok := roundconfig.ModelRecommendations(category)
		if !ok {
			return profilesShowResponse{}, fmt.Errorf("recommendations for Agent Work Category %q are not configured", category)
		}
		outputRecommendations := make([]profilesShowRecommendationOutput, 0, len(recommendations))
		for _, recommendation := range recommendations {
			outputRecommendations = append(outputRecommendations, profilesShowRecommendationOutput{
				ModelRecommendation: recommendation,
				UnavailableReason:   unavailable[recommendation.Selection],
			})
		}
		response.Profiles = append(response.Profiles, profilesShowProfile{
			Category:             category,
			Source:               profile.Source,
			InheritedFrom:        profile.InheritedFrom,
			RecommendationSource: recommendationSource,
			Preferred:            profile.Profile.Preferred,
			Fallbacks:            append([]roundconfig.AgentSelection(nil), profile.Profile.Fallbacks...),
			Recommendations:      outputRecommendations,
		})
	}
	return response, nil
}

func printProfilesShowText(response profilesShowResponse, stdout io.Writer) {
	for index, profile := range response.Profiles {
		if index > 0 {
			fmt.Fprintln(stdout)
		}
		fmt.Fprintf(stdout, "Category: %s\n", profile.Category)
		fmt.Fprintf(stdout, "Profile source: %s\n", profile.Source)
		fmt.Fprintf(stdout, "Profile inherited from: %s\n", categoryOrDash(profile.InheritedFrom))
		fmt.Fprintf(stdout, "Preferred Selection: %s\n", formatProfileSelection(profile.Preferred))
		fmt.Fprintln(stdout, "Fallback Chain:")
		for fallbackIndex, fallback := range profile.Fallbacks {
			fmt.Fprintf(stdout, "  %d. %s\n", fallbackIndex+1, formatProfileSelection(fallback))
		}
		fmt.Fprintf(stdout, "Recommendation source: %s\n", profile.RecommendationSource)
		fmt.Fprintf(stdout, "Recommendations snapshot: %s\n", roundconfig.ModelRecommendationSnapshotVersion)
		fmt.Fprintln(stdout, "Recommendations:")
		for _, recommendation := range profile.Recommendations {
			fmt.Fprintf(stdout, "  %d. %s — %s %s, average cost $%.2f, source %s, category_specific=%t\n",
				recommendation.Rank,
				formatProfileSelection(recommendation.Selection),
				recommendation.Benchmark,
				formatRecommendationPercent(recommendation.ResultPercent),
				recommendation.AverageCostUSD,
				recommendation.SourceAsOf,
				recommendation.CategorySpecific,
			)
			if recommendation.UnavailableReason != "" {
				fmt.Fprintf(stdout, "     unavailable: %s\n", recommendation.UnavailableReason)
			}
			fmt.Fprintf(stdout, "     rationale: %s\n", recommendation.Rationale)
		}
	}
}

func printProfilesFailure(err error, stderr io.Writer) {
	fmt.Fprintf(stderr, "%s: profiles failed: %v\n", app.Name, err)
	fmt.Fprintf(stderr, "Run '%s profiles --help' for usage.\n", app.Name)
}

func categoryOrDash(category roundconfig.WorkCategory) string {
	if category == "" {
		return "—"
	}
	return string(category)
}

func formatProfileSelection(selection roundconfig.AgentSelection) string {
	reasoning := strings.TrimSpace(selection.ReasoningEffort)
	if reasoning == "" {
		reasoning = "model-managed"
	}
	return strings.TrimSpace(selection.Runtime) + " / " + strings.TrimSpace(selection.Model) + " / " + reasoning
}

func formatRecommendationPercent(value float64) string {
	if value == float64(int(value)) {
		return fmt.Sprintf("%d%%", int(value))
	}
	return fmt.Sprintf("%.1f%%", value)
}

func profilesCategoryList() string {
	categories := roundconfig.AllWorkCategories()
	values := make([]string, 0, len(categories))
	for _, category := range categories {
		values = append(values, string(category))
	}
	return strings.Join(values, ", ")
}
