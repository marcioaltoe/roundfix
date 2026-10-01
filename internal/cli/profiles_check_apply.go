package cli

import (
	"context"
	"fmt"
	"io"
	"strings"

	roundconfig "roundfix/internal/config"
)

func runProfilesCheckApplyCommand(ctx context.Context, args []string, stdout, stderr io.Writer, environment commandEnvironment) int {
	if commandWantsHelp(args) {
		fmt.Fprint(stdout, commandUsage("profiles check"))
		return exitOK
	}
	req := profilesConfigureRequest{}
	fs := flagSet("profiles check")
	apply := fs.Bool("apply", false, "Adopt differing Recommended Profiles")
	fs.StringVar(&req.scope, "scope", "", "Config scope: user or project")
	fs.BoolVar(&req.dryRun, "dry-run", false, "Prove and preview without writing")
	fs.BoolVar(&req.yes, "yes", false, "Write without confirmation after proof")
	fs.BoolVar(&req.json, "json", false, "Print roundfix/profiles-configure/v1 JSON")
	if err := fs.Parse(args); err != nil {
		printProfilesFailure(validationError{message: err.Error()}, stderr)
		return exitPreflight
	}
	req.scope = strings.TrimSpace(req.scope)
	if !*apply || (req.scope != roundconfig.InitScopeUser && req.scope != roundconfig.InitScopeProject) || len(fs.Args()) > 0 {
		printProfilesFailure(validationError{message: "--apply requires --scope user|project and no positional arguments"}, stderr)
		return exitPreflight
	}
	loaded, err := loadCommandConfig(environment, stderr)
	if err != nil {
		return printProfilesConfigureError(req, roundconfig.ProfileConfigResult{Scope: req.scope}, err, stdout, stderr)
	}
	check, err := roundconfig.CheckRecommendations(loaded.Config)
	if err != nil {
		return printProfilesConfigureError(req, roundconfig.ProfileConfigResult{Scope: req.scope}, err, stdout, stderr)
	}
	profiles := roundconfig.Profiles{}
	for _, row := range check.Categories {
		if row.Status != roundconfig.RecommendationDiffers {
			continue
		}
		if req.scope == roundconfig.InitScopeUser && row.Source == roundconfig.ProfileSourceProject {
			if _, err := fmt.Fprintf(stderr, "recommendations: %s is defined by Project Config; use --scope project to adopt it\n", row.Category); err != nil {
				return printProfilesConfigureOutputError(err, stderr)
			}
			continue
		}
		profiles[row.Category] = roundconfig.ProfileEntry{Profile: row.Recommended}
	}
	if len(profiles) == 0 {
		if req.json {
			req.dryRun = false
			if err := printProfilesConfigureSuccess(req, roundconfig.ProfileConfigResult{Scope: req.scope, Profiles: profiles}, stdout); err != nil {
				return printProfilesConfigureOutputError(err, stderr)
			}
		} else if _, err := fmt.Fprintln(stdout, "Profile configuration unchanged: nothing to adopt"); err != nil {
			return printProfilesConfigureOutputError(err, stderr)
		}
		return exitOK
	}
	loadOptions, err := environment.loadOptions(stderr)
	if err != nil {
		return printProfilesConfigureError(req, roundconfig.ProfileConfigResult{Scope: req.scope}, err, stdout, stderr)
	}
	return writeProfilesConfiguration(ctx, req, profiles, loadOptions, stdout, stderr, environment)
}
