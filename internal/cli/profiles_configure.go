package cli

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	roundconfig "roundfix/internal/config"
)

const profilesConfigureSchema = "roundfix/profiles-configure/v1"

type profilesConfigureRequest struct {
	scope  string
	file   string
	dryRun bool
	json   bool
	yes    bool
}

type profilesConfigureResponse struct {
	Schema   string                     `json:"schema"`
	Changed  bool                       `json:"changed"`
	Scope    string                     `json:"scope"`
	Path     string                     `json:"path"`
	Profiles []profilesConfigureProfile `json:"profiles"`
	Error    string                     `json:"error,omitempty"`
}

type profilesConfigureProfile struct {
	Category  roundconfig.WorkCategory     `json:"category"`
	Preferred roundconfig.AgentSelection   `json:"preferred"`
	Fallbacks []roundconfig.AgentSelection `json:"fallbacks"`
}

var profilesConfigureInput = func() io.Reader { return os.Stdin }
var confirmProfilesConfigure = defaultConfirmProfilesConfigure

func runProfilesConfigureCommand(ctx context.Context, args []string, stdout, stderr io.Writer) int {
	if commandWantsHelp(args) {
		fmt.Fprint(stdout, commandUsage("profiles configure"))
		return exitOK
	}
	req, err := parseProfilesConfigureCommand(args)
	if err != nil {
		return printProfilesConfigureError(req, roundconfig.ProfileConfigResult{}, err, stdout, stderr)
	}

	profiles, err := profilesForConfigureRequest(ctx, req, stderr)
	if err != nil {
		return printProfilesConfigureError(req, roundconfig.ProfileConfigResult{Scope: req.scope}, err, stdout, stderr)
	}
	result, err := roundconfig.WriteProfilesConfig(ctx, roundconfig.ProfileConfigOptions{
		Scope:    req.scope,
		Profiles: profiles,
		DryRun:   true,
	})
	if err != nil {
		return printProfilesConfigureError(req, roundconfig.ProfileConfigResult{Scope: req.scope}, err, stdout, stderr)
	}
	if req.dryRun {
		if !req.json {
			fmt.Fprint(stdout, profilesConfigurePreview(result))
		}
		if err := printProfilesConfigureSuccess(req, result, stdout); err != nil {
			return printProfilesConfigureOutputError(err, stderr)
		}
		return exitOK
	}

	preview := profilesConfigurePreview(result)
	if !req.yes {
		confirmed, err := confirmProfilesConfigure(ctx, stderr, preview)
		if err != nil {
			return printProfilesConfigureError(req, result, err, stdout, stderr)
		}
		if !confirmed {
			result.Changed = false
			if req.json {
				if err := printProfilesConfigureSuccess(req, result, stdout); err != nil {
					return printProfilesConfigureOutputError(err, stderr)
				}
			} else {
				fmt.Fprintf(stdout, "Profile configuration unchanged: confirmation declined for %s\n", result.Path)
			}
			return exitOK
		}
	}

	result, err = roundconfig.WriteProfilesConfig(ctx, roundconfig.ProfileConfigOptions{
		Scope:    req.scope,
		Profiles: profiles,
	})
	if err != nil {
		return printProfilesConfigureError(req, result, err, stdout, stderr)
	}
	if err := printProfilesConfigureSuccess(req, result, stdout); err != nil {
		return printProfilesConfigureOutputError(err, stderr)
	}
	return exitOK
}

func parseProfilesConfigureCommand(args []string) (profilesConfigureRequest, error) {
	req := profilesConfigureRequest{}
	fs := flagSet("profiles configure")
	fs.StringVar(&req.scope, "scope", "", "Config scope: user or project")
	fs.StringVar(&req.file, "file", "", "Strict profile fragment YAML")
	fs.BoolVar(&req.dryRun, "dry-run", false, "Validate and render without writing")
	fs.BoolVar(&req.json, "json", false, "Print roundfix/profiles-configure/v1 JSON")
	fs.BoolVar(&req.yes, "yes", false, "Write without confirmation after validation")
	if err := fs.Parse(args); err != nil {
		return req, validationError{message: err.Error()}
	}
	if remaining := fs.Args(); len(remaining) > 0 {
		return req, validationError{message: fmt.Sprintf("unexpected argument %q", remaining[0])}
	}
	req.scope = strings.TrimSpace(req.scope)
	req.file = strings.TrimSpace(req.file)
	if req.scope == "" {
		return req, validationError{message: "--scope is required; supported values: user, project"}
	}
	if req.scope != roundconfig.InitScopeUser && req.scope != roundconfig.InitScopeProject {
		return req, validationError{message: fmt.Sprintf("unsupported profiles scope %q; supported values: user, project", req.scope)}
	}
	return req, nil
}

func profilesForConfigureRequest(ctx context.Context, req profilesConfigureRequest, stderr io.Writer) (roundconfig.Profiles, error) {
	if req.file != "" {
		content, err := os.ReadFile(req.file)
		if err != nil {
			return nil, fmt.Errorf("read profiles file %q: %w", req.file, err)
		}
		profiles, err := roundconfig.ParseProfilesFragment(content)
		if err != nil {
			return nil, fmt.Errorf("read profiles file %q: %w", req.file, err)
		}
		return profiles, nil
	}
	return collectProfilesConfigureInput(ctx, profilesConfigureInput(), stderr)
}

func collectProfilesConfigureInput(ctx context.Context, input io.Reader, output io.Writer) (roundconfig.Profiles, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	reader := bufio.NewReader(input)
	fmt.Fprintln(output, "Roundfix Profiles Configure")
	fmt.Fprintln(output, "Create one complete Agent Selection Profile.")
	category, err := readProfilesConfigureCategory(ctx, reader, output)
	if err != nil {
		return nil, err
	}
	printProfilesConfigureRecommendations(output, category)

	preferred, err := readProfilesConfigureSelection(ctx, reader, output, "Preferred Selection")
	if err != nil {
		return nil, err
	}
	fallbacks := []roundconfig.AgentSelection{}
	for {
		fmt.Fprintf(output, "Fallback %d runtime (blank to finish): ", len(fallbacks)+1)
		runtime, err := readProfilesLine(ctx, reader)
		if err != nil {
			return nil, err
		}
		if strings.TrimSpace(runtime) == "" {
			if len(fallbacks) == 0 {
				return nil, validationError{message: "profiles." + string(category) + ".fallbacks must include at least one complete Agent Selection before confirmation"}
			}
			break
		}
		fallback, err := readProfilesConfigureSelectionAfterRuntime(ctx, reader, output, "Fallback Selection", runtime)
		if err != nil {
			return nil, err
		}
		fallbacks = append(fallbacks, fallback)
	}
	return roundconfig.NormalizeProfilesFragment(roundconfig.Profiles{
		category: {Profile: roundconfig.AgentSelectionProfile{Preferred: preferred, Fallbacks: fallbacks}},
	})
}

func readProfilesConfigureCategory(ctx context.Context, reader *bufio.Reader, output io.Writer) (roundconfig.WorkCategory, error) {
	fmt.Fprintf(output, "Category (%s): ", profilesCategoryList())
	value, err := readProfilesLine(ctx, reader)
	if err != nil {
		return "", err
	}
	category, ok := roundconfig.ParseWorkCategory(value)
	if !ok {
		return "", validationError{message: fmt.Sprintf("unknown profile category %q; supported values: %s", strings.TrimSpace(value), profilesCategoryList())}
	}
	return category, nil
}

func readProfilesConfigureSelection(ctx context.Context, reader *bufio.Reader, output io.Writer, label string) (roundconfig.AgentSelection, error) {
	fmt.Fprintf(output, "%s runtime: ", label)
	runtime, err := readProfilesLine(ctx, reader)
	if err != nil {
		return roundconfig.AgentSelection{}, err
	}
	return readProfilesConfigureSelectionAfterRuntime(ctx, reader, output, label, runtime)
}

func readProfilesConfigureSelectionAfterRuntime(ctx context.Context, reader *bufio.Reader, output io.Writer, label string, runtime string) (roundconfig.AgentSelection, error) {
	fmt.Fprintf(output, "%s model: ", label)
	model, err := readProfilesLine(ctx, reader)
	if err != nil {
		return roundconfig.AgentSelection{}, err
	}
	fmt.Fprintf(output, "%s reasoning_effort (blank means model-managed): ", label)
	reasoning, err := readProfilesLine(ctx, reader)
	if err != nil {
		return roundconfig.AgentSelection{}, err
	}
	profiles, err := roundconfig.NormalizeProfilesFragment(roundconfig.Profiles{
		roundconfig.CategoryGeneral: {
			Profile: roundconfig.AgentSelectionProfile{
				Preferred: roundconfig.AgentSelection{
					Runtime:         runtime,
					Model:           model,
					ReasoningEffort: reasoning,
				},
				Fallbacks: []roundconfig.AgentSelection{{Runtime: "codex", Model: "temporary-validation-fallback", ReasoningEffort: ""}},
			},
		},
	})
	if err != nil {
		return roundconfig.AgentSelection{}, err
	}
	return profiles[roundconfig.CategoryGeneral].Profile.Preferred, nil
}

func readProfilesLine(ctx context.Context, reader *bufio.Reader) (string, error) {
	line, err := reader.ReadString('\n')
	if err != nil && !errors.Is(err, io.EOF) {
		return "", fmt.Errorf("read profiles configure input: %w", err)
	}
	if err := ctx.Err(); err != nil {
		return "", err
	}
	return strings.TrimSpace(line), nil
}

func printProfilesConfigureRecommendations(output io.Writer, category roundconfig.WorkCategory) {
	recommendations, source, ok := roundconfig.ModelRecommendations(category)
	if !ok {
		return
	}
	fmt.Fprintf(output, "Recommendations for %s (source: %s, advisory only):\n", category, source)
	for _, recommendation := range recommendations {
		fmt.Fprintf(output, "  %d. %s\n", recommendation.Rank, formatProfileSelection(recommendation.Selection))
	}
}

func defaultConfirmProfilesConfigure(ctx context.Context, stderr io.Writer, preview string) (bool, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	if err := ctx.Err(); err != nil {
		return false, err
	}
	if _, err := fmt.Fprint(stderr, preview); err != nil {
		return false, fmt.Errorf("write profiles configure preview: %w", err)
	}
	if _, err := fmt.Fprint(stderr, "Write this Agent Selection Profile config? [y/N]: "); err != nil {
		return false, fmt.Errorf("write profiles configure confirmation prompt: %w", err)
	}
	line, err := bufio.NewReader(os.Stdin).ReadString('\n')
	if err != nil && !errors.Is(err, io.EOF) {
		return false, fmt.Errorf("read profiles configure confirmation prompt: %w", err)
	}
	if err := ctx.Err(); err != nil {
		return false, err
	}
	return normalizeSetupYesNo(line)
}

func profilesConfigurePreview(result roundconfig.ProfileConfigResult) string {
	var builder strings.Builder
	fmt.Fprintf(&builder, "Profile Configure Preview\n")
	fmt.Fprintf(&builder, "Scope: %s\n", result.Scope)
	fmt.Fprintf(&builder, "Path: %s\n", result.Path)
	for _, profile := range profilesConfigureProfiles(result.Profiles) {
		fmt.Fprintf(&builder, "Category: %s\n", profile.Category)
		fmt.Fprintf(&builder, "Preferred Selection: %s\n", formatProfileSelection(profile.Preferred))
		fmt.Fprintf(&builder, "Fallback Chain:\n")
		for index, fallback := range profile.Fallbacks {
			fmt.Fprintf(&builder, "  %d. %s\n", index+1, formatProfileSelection(fallback))
		}
	}
	return builder.String()
}

func printProfilesConfigureSuccess(req profilesConfigureRequest, result roundconfig.ProfileConfigResult, stdout io.Writer) error {
	if req.json {
		return json.NewEncoder(stdout).Encode(profilesConfigureResponseForResult(result, ""))
	}
	if req.dryRun {
		fmt.Fprintf(stdout, "Profile configuration dry run: %s\n", result.Path)
		return nil
	}
	if result.Changed {
		fmt.Fprintf(stdout, "Profile configuration written: %s\n", result.Path)
		return nil
	}
	fmt.Fprintf(stdout, "Profile configuration unchanged: %s\n", result.Path)
	return nil
}

func printProfilesConfigureOutputError(err error, stderr io.Writer) int {
	printProfilesFailure(fmt.Errorf("encode profiles configure JSON: %w", err), stderr)
	return exitRunFailed
}

func printProfilesConfigureError(req profilesConfigureRequest, result roundconfig.ProfileConfigResult, err error, stdout, stderr io.Writer) int {
	if req.json {
		response := profilesConfigureResponseForResult(result, err.Error())
		response.Changed = false
		_ = json.NewEncoder(stdout).Encode(response)
	}
	printProfilesFailure(err, stderr)
	return exitPreflight
}

func profilesConfigureResponseForResult(result roundconfig.ProfileConfigResult, errorMessage string) profilesConfigureResponse {
	return profilesConfigureResponse{
		Schema:   profilesConfigureSchema,
		Changed:  result.Changed,
		Scope:    result.Scope,
		Path:     result.Path,
		Profiles: profilesConfigureProfiles(result.Profiles),
		Error:    errorMessage,
	}
}

func profilesConfigureProfiles(profiles roundconfig.Profiles) []profilesConfigureProfile {
	output := []profilesConfigureProfile{}
	for _, category := range roundconfig.AllWorkCategories() {
		entry, ok := profiles[category]
		if !ok {
			continue
		}
		output = append(output, profilesConfigureProfile{
			Category:  category,
			Preferred: entry.Profile.Preferred,
			Fallbacks: append([]roundconfig.AgentSelection(nil), entry.Profile.Fallbacks...),
		})
	}
	return output
}

func flagSet(name string) *flag.FlagSet {
	fs := flag.NewFlagSet(name, flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	return fs
}
