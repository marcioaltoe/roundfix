package cli

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"strings"

	"roundfix/internal/baseline"
)

type baselineStartersCommandRequest struct {
	repo         string
	profile      string
	confirmation string
	format       string
}

func parseBaselineStartersCommand(args []string) (baselineStartersCommandRequest, error) {
	request := baselineStartersCommandRequest{repo: ".", format: "text"}
	flags := flag.NewFlagSet("baseline starters", flag.ContinueOnError)
	flags.SetOutput(io.Discard)
	flags.StringVar(&request.repo, "repo", ".", "Repository path")
	flags.StringVar(&request.profile, "profile", "", "Baseline Profile (required)")
	flags.StringVar(&request.confirmation, "confirm-plan", "", "Exact approved Plan Digest")
	flags.StringVar(&request.format, "format", "text", "Output format: text or json")
	if err := flags.Parse(args); err != nil {
		return request, fmt.Errorf("invalid baseline starters arguments: %w; run 'roundfix baseline starters --help' for usage", err)
	}
	if len(flags.Args()) != 0 {
		return request, fmt.Errorf("unexpected argument %q; run 'roundfix baseline starters --help' for usage", flags.Args()[0])
	}
	request.repo = strings.TrimSpace(request.repo)
	request.profile = strings.TrimSpace(request.profile)
	request.confirmation = strings.TrimSpace(request.confirmation)
	request.format = strings.TrimSpace(request.format)
	if request.profile == "" {
		return request, errors.New("--profile is required")
	}
	if request.repo == "" {
		return request, errors.New("--repo cannot be empty")
	}
	if request.format != "text" && request.format != "json" {
		return request, fmt.Errorf("unsupported --format %q; use text or json", request.format)
	}
	return request, nil
}

func runBaselineStartersCommand(ctx context.Context, args []string, stdout, stderr io.Writer) int {
	if commandWantsHelp(args) {
		fmt.Fprint(stdout, commandUsage("baseline starters"))
		return exitOK
	}
	request, err := parseBaselineStartersCommand(args)
	if err != nil {
		printBaselineUsageError("baseline starters", err.Error(), stderr)
		return exitPreflight
	}
	catalog, err := baseline.LoadEmbeddedCatalog()
	if err != nil {
		fmt.Fprintf(stderr, "roundfix baseline starters: %v\n", err)
		return exitRunFailed
	}
	payload, err := baseline.Starters(ctx, baseline.StartersRequest{Repository: request.repo, ProfileID: request.profile, ConfirmPlan: request.confirmation}, catalog)
	if err != nil {
		if errors.Is(err, baseline.ErrStartersStale) {
			fmt.Fprintf(stderr, "roundfix baseline starters: %s\n", baseline.ErrStartersStale)
			return exitUnverified
		}
		var profileErr *baseline.ProfileResolutionError
		if errors.As(err, &profileErr) {
			printBaselineUsageError("baseline starters", err.Error(), stderr)
			return exitPreflight
		}
		fmt.Fprintf(stderr, "roundfix baseline starters: %v\n", err)
		return exitRunFailed
	}
	if err := writeBaselineStartersPayload(payload, request.format, stdout); err != nil {
		fmt.Fprintf(stderr, "roundfix baseline starters: output failed: %v\n", err)
		return exitRunFailed
	}
	if payload.State == "preview" {
		return exitUnverified
	}
	return exitOK
}

func writeBaselineStartersPayload(payload baseline.StartersPayload, format string, stdout io.Writer) error {
	if format == "json" {
		return json.NewEncoder(stdout).Encode(payload)
	}
	if _, err := fmt.Fprintf(stdout, "Baseline starters: %s\n", payload.Profile); err != nil {
		return err
	}
	for _, starter := range payload.Starters {
		if _, err := fmt.Fprintf(stdout, "- %s %s\n", starter.Action, starter.Path); err != nil {
			return err
		}
	}
	if _, err := fmt.Fprintf(stdout, "State: %s\nPlan Digest: %s\n", payload.State, payload.PlanDigest); err != nil {
		return err
	}
	if payload.State == "preview" {
		_, err := fmt.Fprintf(stdout, "Next: roundfix baseline starters --profile %s --confirm-plan %s\n", payload.Profile, payload.PlanDigest)
		return err
	}
	return nil
}
