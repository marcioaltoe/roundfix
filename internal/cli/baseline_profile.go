package cli

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"roundfix/internal/app"
	"roundfix/internal/baseline"
)

const baselineResultSchema = "roundfix/baseline-result/v1"

type baselineProfileResult struct {
	SchemaVersion string                     `json:"schemaVersion"`
	Operation     string                     `json:"operation"`
	State         string                     `json:"state"`
	Profile       *baseline.ResolvedProfile  `json:"profile,omitempty"`
	Profiles      []baseline.ResolvedProfile `json:"profiles,omitempty"`
	Path          string                     `json:"path,omitempty"`
	FromProfile   string                     `json:"fromProfile,omitempty"`
	Category      string                     `json:"category,omitempty"`
	Message       string                     `json:"message,omitempty"`
	NextAction    string                     `json:"nextAction,omitempty"`
}

type baselinePlanPreflightResult struct {
	SchemaVersion string                      `json:"schemaVersion"`
	Operation     string                      `json:"operation"`
	State         string                      `json:"state"`
	Category      string                      `json:"category"`
	Repository    baseline.RepositoryIdentity `json:"repository"`
	Snapshot      baseline.RepositorySnapshot `json:"snapshot"`
	Message       string                      `json:"message"`
	NextAction    string                      `json:"nextAction"`
}

func runBaselineCommand(ctx context.Context, args []string, stdout, stderr io.Writer) int {
	if len(args) == 0 || len(args) == 1 && commandWantsHelp(args) {
		fmt.Fprint(stdout, commandUsage("baseline"))
		return exitOK
	}
	switch args[0] {
	case "plan":
		return runBaselinePlanCommand(ctx, args[1:], stdout, stderr)
	case "profile":
		return runBaselineProfileCommand(args[1:], stdout, stderr)
	default:
		printBaselineProfileFailure(
			"baseline",
			validationError{message: fmt.Sprintf("unknown baseline command %q", args[0])},
			false,
			stdout,
			stderr,
		)
		return exitPreflight
	}
}

func runBaselinePlanCommand(ctx context.Context, args []string, stdout, stderr io.Writer) int {
	if commandWantsHelp(args) {
		fmt.Fprint(stdout, commandUsage("baseline plan"))
		return exitOK
	}
	req, err := parseBaselinePlanCommand(args)
	jsonOutput := req.format == "json" || baselinePlanJSONRequested(args)
	if err != nil {
		printBaselinePlanFailure(err, jsonOutput, baselinePlanPreflightResult{}, stdout, stderr)
		return exitPreflight
	}
	inspection, err := baseline.InspectRepository(ctx, req.repo, nil)
	if err != nil {
		printBaselinePlanFailure(err, jsonOutput, baselinePlanPreflightResult{}, stdout, stderr)
		return exitPreflight
	}
	result := baselinePlanPreflightResult{
		SchemaVersion: baselineResultSchema,
		Operation:     "plan",
		Repository:    inspection.Identity,
		Snapshot:      inspection.Snapshot,
	}
	if len(inspection.Snapshot.Blocking) != 0 {
		result.State = "failed"
		result.Category = "preflight"
		result.Message = "repository preflight found unsafe bounded carriers"
		result.NextAction = "repair each blocking carrier and rerun roundfix baseline plan"
		for _, finding := range inspection.Snapshot.Blocking {
			fmt.Fprintf(stderr, "%s: %s: %s: %s\n", app.Name, finding.Code, finding.Path, finding.Message)
		}
		if jsonOutput {
			if err := encodeBaselinePlanResult(stdout, result); err != nil {
				fmt.Fprintf(stderr, "%s: baseline plan output failed: %v\n", app.Name, err)
				return exitRunFailed
			}
		} else {
			printBaselinePlanText(result, stdout)
		}
		return exitPreflight
	}

	result.State = "action_required"
	result.Category = "decision"
	result.Message = "repository preflight passed; instruction preservation requires a decision"
	result.NextAction = "choose greenfield or preservation mode before completing the Baseline Plan"
	if jsonOutput {
		if err := encodeBaselinePlanResult(stdout, result); err != nil {
			fmt.Fprintf(stderr, "%s: baseline plan output failed: %v\n", app.Name, err)
			return exitRunFailed
		}
	} else {
		printBaselinePlanText(result, stdout)
	}
	return exitUnverified
}

type baselinePlanCommandRequest struct {
	repo   string
	format string
}

func parseBaselinePlanCommand(args []string) (baselinePlanCommandRequest, error) {
	req := baselinePlanCommandRequest{repo: ".", format: "text"}
	flags := flag.NewFlagSet("baseline plan", flag.ContinueOnError)
	flags.SetOutput(io.Discard)
	flags.StringVar(&req.repo, "repo", ".", "Git worktree or a path inside it")
	flags.StringVar(&req.format, "format", "text", "Output format: text or json")
	if err := flags.Parse(args); err != nil {
		return baselinePlanCommandRequest{}, validationError{
			message: fmt.Sprintf("invalid baseline plan arguments: %v; run '%s baseline plan --help' for usage", err, app.Name),
		}
	}
	if remaining := flags.Args(); len(remaining) != 0 {
		return baselinePlanCommandRequest{}, validationError{
			message: fmt.Sprintf("unexpected argument %q; run '%s baseline plan --help' for usage", remaining[0], app.Name),
		}
	}
	req.repo = strings.TrimSpace(req.repo)
	req.format = strings.TrimSpace(req.format)
	if req.repo == "" {
		return baselinePlanCommandRequest{}, validationError{message: "--repo cannot be empty"}
	}
	if req.format != "text" && req.format != "json" {
		return baselinePlanCommandRequest{}, validationError{
			message: fmt.Sprintf("unsupported --format %q; use text or json", req.format),
		}
	}
	return req, nil
}

func printBaselinePlanText(result baselinePlanPreflightResult, stdout io.Writer) {
	state := strings.ReplaceAll(result.State, "_", " ")
	fmt.Fprintf(stdout, "Baseline plan preflight: %s\n", state)
	if result.Repository.Digest != "" {
		fmt.Fprintf(stdout, "Repository identity: %s\n", result.Repository.Digest)
		fmt.Fprintf(stdout, "Git object format: %s\n", result.Repository.ObjectFormat)
		fmt.Fprintf(stdout, "Root commits: %d\n", len(result.Repository.RootCommits))
	}
	if result.Snapshot.Digest != "" {
		fmt.Fprintf(stdout, "Bounded snapshot: %s\n", result.Snapshot.Digest)
		fmt.Fprintf(stdout, "Instruction carriers: %d\n", len(result.Snapshot.Carriers))
		fmt.Fprintf(stdout, "Trusted sources: %d\n", len(result.Snapshot.Sources))
		fmt.Fprintf(stdout, "Bounded preimages: %d\n", len(result.Snapshot.Preimages))
	}
	for _, warning := range result.Snapshot.Warnings {
		fmt.Fprintf(stdout, "Warning: %s: %s: %s\n", warning.Code, warning.Path, warning.Message)
	}
	for _, finding := range result.Snapshot.Blocking {
		fmt.Fprintf(stdout, "Blocking: %s: %s: %s\n", finding.Code, finding.Path, finding.Message)
	}
	if result.Message != "" {
		fmt.Fprintf(stdout, "Result: %s\n", result.Message)
	}
	if result.NextAction != "" {
		fmt.Fprintf(stdout, "Next action: %s\n", result.NextAction)
	}
}

func printBaselinePlanFailure(
	err error,
	jsonOutput bool,
	base baselinePlanPreflightResult,
	stdout io.Writer,
	stderr io.Writer,
) {
	fmt.Fprintf(stderr, "%s: baseline plan failed: %v\n", app.Name, err)
	fmt.Fprintf(stderr, "Run '%s baseline plan --help' for usage.\n", app.Name)
	if !jsonOutput {
		return
	}
	base.SchemaVersion = baselineResultSchema
	base.Operation = "plan"
	base.State = "failed"
	base.Category = "preflight"
	base.Message = err.Error()
	base.NextAction = "correct the repository or command input and rerun roundfix baseline plan"
	if encodeErr := encodeBaselinePlanResult(stdout, base); encodeErr != nil {
		fmt.Fprintf(stderr, "%s: baseline plan output failed: %v\n", app.Name, encodeErr)
	}
}

func encodeBaselinePlanResult(stdout io.Writer, result baselinePlanPreflightResult) error {
	encoder := json.NewEncoder(stdout)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(result); err != nil {
		return fmt.Errorf("encode Baseline result: %w", err)
	}
	return nil
}

func baselinePlanJSONRequested(args []string) bool {
	for index, arg := range args {
		if arg == "--format" && index+1 < len(args) && strings.TrimSpace(args[index+1]) == "json" {
			return true
		}
		if strings.HasPrefix(arg, "--format=") && strings.TrimSpace(strings.TrimPrefix(arg, "--format=")) == "json" {
			return true
		}
	}
	return false
}

func runBaselineProfileCommand(args []string, stdout, stderr io.Writer) int {
	if len(args) == 0 || len(args) == 1 && commandWantsHelp(args) {
		fmt.Fprint(stdout, commandUsage("baseline profile"))
		return exitOK
	}
	switch args[0] {
	case "init":
		return runBaselineProfileInitCommand(args[1:], stdout, stderr)
	case "show":
		return runBaselineProfileShowCommand(args[1:], stdout, stderr)
	case "validate":
		return runBaselineProfileValidateCommand(args[1:], stdout, stderr)
	default:
		printBaselineProfileFailure(
			"profile",
			validationError{message: fmt.Sprintf("unknown baseline profile command %q", args[0])},
			false,
			stdout,
			stderr,
		)
		return exitPreflight
	}
}

func runBaselineProfileInitCommand(args []string, stdout, stderr io.Writer) int {
	if commandWantsHelp(args) {
		fmt.Fprint(stdout, commandUsage("baseline profile init"))
		return exitOK
	}
	fs := flag.NewFlagSet("baseline profile init", flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	id := fs.String("id", "", "Repository-owned Baseline Profile ID")
	from := fs.String("from", "go-cli-tui", "Built-in Baseline Profile source")
	if err := fs.Parse(args); err != nil {
		printBaselineProfileFailure("profile.init", validationError{message: err.Error()}, false, stdout, stderr)
		return exitPreflight
	}
	if remaining := fs.Args(); len(remaining) != 0 {
		printBaselineProfileFailure("profile.init", validationError{message: fmt.Sprintf("unexpected argument %q", remaining[0])}, false, stdout, stderr)
		return exitPreflight
	}
	if strings.TrimSpace(*id) == "" {
		printBaselineProfileFailure("profile.init", validationError{message: "--id is required"}, false, stdout, stderr)
		return exitPreflight
	}
	repoRoot, err := findBaselineRepositoryRoot()
	if err != nil {
		printBaselineProfileFailure("profile.init", err, false, stdout, stderr)
		return exitPreflight
	}
	if repoRoot == "" {
		printBaselineProfileFailure("profile.init", errors.New("repository-owned Baseline Profile init requires a Git repository"), false, stdout, stderr)
		return exitPreflight
	}
	catalog, err := baseline.LoadEmbeddedCatalog()
	if err != nil {
		printBaselineProfileFailure("profile.init", err, false, stdout, stderr)
		return exitRunFailed
	}
	result, err := baseline.InitCustomProfile(repoRoot, *id, *from, catalog)
	if err != nil {
		printBaselineProfileFailure("profile.init", err, false, stdout, stderr)
		if isBaselineProfileExecutionError(err) {
			return exitRunFailed
		}
		return exitPreflight
	}
	relativePath := baselineProfileOutputPath(repoRoot, result.Path)
	fmt.Fprintf(stdout, "Created repository-owned Baseline Profile %s\n", result.Profile.ID)
	fmt.Fprintf(stdout, "Path: %s\n", relativePath)
	fmt.Fprintf(stdout, "Built-in source: %s\n", result.FromProfile)
	fmt.Fprintf(stdout, "Digest: %s\n", result.Profile.Digest)
	return exitOK
}

func runBaselineProfileShowCommand(args []string, stdout, stderr io.Writer) int {
	if commandWantsHelp(args) {
		fmt.Fprint(stdout, commandUsage("baseline profile show"))
		return exitOK
	}
	target, format, err := parseBaselineProfileTargetFormat(args, true)
	jsonOutput := format == "json" || baselineProfileJSONRequested(args)
	if err != nil {
		printBaselineProfileFailure("profile.show", err, jsonOutput, stdout, stderr)
		return exitPreflight
	}
	catalog, err := baseline.LoadEmbeddedCatalog()
	if err != nil {
		printBaselineProfileFailure("profile.show", err, jsonOutput, stdout, stderr)
		return exitRunFailed
	}
	repoRoot, err := findBaselineRepositoryRoot()
	if err != nil {
		printBaselineProfileFailure("profile.show", err, jsonOutput, stdout, stderr)
		return exitPreflight
	}
	profile, err := resolveBaselineProfileTarget(repoRoot, target, catalog)
	if err != nil {
		printBaselineProfileFailure("profile.show", err, jsonOutput, stdout, stderr)
		return exitPreflight
	}
	profile.Path = baselineProfileOutputPath(repoRoot, profile.Path)
	if jsonOutput {
		return encodeBaselineProfileResult(stdout, stderr, baselineProfileResult{
			SchemaVersion: baselineResultSchema,
			Operation:     "profile.show",
			State:         "valid",
			Profile:       &profile,
		})
	}
	printBaselineProfileText(profile, stdout)
	return exitOK
}

func runBaselineProfileValidateCommand(args []string, stdout, stderr io.Writer) int {
	if commandWantsHelp(args) {
		fmt.Fprint(stdout, commandUsage("baseline profile validate"))
		return exitOK
	}
	target, format, err := parseBaselineProfileTargetFormat(args, false)
	jsonOutput := format == "json" || baselineProfileJSONRequested(args)
	if err != nil {
		printBaselineProfileFailure("profile.validate", err, jsonOutput, stdout, stderr)
		return exitPreflight
	}
	catalog, err := baseline.LoadEmbeddedCatalog()
	if err != nil {
		printBaselineProfileFailure("profile.validate", err, jsonOutput, stdout, stderr)
		return exitRunFailed
	}
	repoRoot, err := findBaselineRepositoryRoot()
	if err != nil {
		printBaselineProfileFailure("profile.validate", err, jsonOutput, stdout, stderr)
		return exitPreflight
	}
	var profiles []baseline.ResolvedProfile
	if target == "" {
		if repoRoot == "" {
			printBaselineProfileFailure("profile.validate", errors.New("repository profile discovery requires a Git repository"), jsonOutput, stdout, stderr)
			return exitPreflight
		}
		profiles, err = baseline.DiscoverRepositoryProfiles(repoRoot, catalog)
	} else {
		var profile baseline.ResolvedProfile
		profile, err = resolveBaselineProfileTarget(repoRoot, target, catalog)
		if err == nil {
			profiles = []baseline.ResolvedProfile{profile}
		}
	}
	if err != nil {
		printBaselineProfileFailure("profile.validate", err, jsonOutput, stdout, stderr)
		return exitPreflight
	}
	for index := range profiles {
		profiles[index].Path = baselineProfileOutputPath(repoRoot, profiles[index].Path)
	}
	if jsonOutput {
		if profiles == nil {
			profiles = []baseline.ResolvedProfile{}
		}
		return encodeBaselineProfileResult(stdout, stderr, baselineProfileResult{
			SchemaVersion: baselineResultSchema,
			Operation:     "profile.validate",
			State:         "valid",
			Profiles:      profiles,
		})
	}
	if len(profiles) == 0 {
		fmt.Fprintln(stdout, "No repository-owned Baseline Profiles found.")
		return exitOK
	}
	for _, profile := range profiles {
		fmt.Fprintf(stdout, "%s: valid (%s) %s\n", profile.ID, profile.Source, profile.Digest)
	}
	return exitOK
}

func parseBaselineProfileTargetFormat(args []string, requireTarget bool) (string, string, error) {
	format := "text"
	target := ""
	formatSeen := false
	for index := 0; index < len(args); index++ {
		arg := args[index]
		switch {
		case arg == "--format":
			if formatSeen {
				return "", "", validationError{message: "--format may be specified only once"}
			}
			index++
			if index >= len(args) {
				return "", "", validationError{message: "flag needs an argument: --format"}
			}
			format = strings.TrimSpace(args[index])
			formatSeen = true
		case strings.HasPrefix(arg, "--format="):
			if formatSeen {
				return "", "", validationError{message: "--format may be specified only once"}
			}
			format = strings.TrimSpace(strings.TrimPrefix(arg, "--format="))
			formatSeen = true
		case strings.HasPrefix(arg, "-"):
			return "", "", validationError{message: fmt.Sprintf("unknown flag %q", arg)}
		default:
			if target != "" {
				return "", "", validationError{message: fmt.Sprintf("unexpected argument %q", arg)}
			}
			target = strings.TrimSpace(arg)
		}
	}
	if format != "text" && format != "json" {
		return "", "", validationError{message: fmt.Sprintf("unsupported format %q; supported values: text, json", format)}
	}
	if requireTarget && target == "" {
		return "", "", validationError{message: "Baseline Profile ID is required"}
	}
	return target, format, nil
}

func resolveBaselineProfileTarget(repoRoot, target string, catalog *baseline.Catalog) (baseline.ResolvedProfile, error) {
	if baselineProfileTargetIsPath(target) {
		if repoRoot == "" {
			return baseline.ResolvedProfile{}, errors.New("explicit Baseline Profile paths require a Git repository")
		}
		return baseline.LoadCustomProfilePath(repoRoot, target, catalog)
	}
	if _, builtIn := catalog.Profile(target); builtIn {
		return baseline.ResolveProfile("", target, catalog)
	}
	if repoRoot == "" {
		return baseline.ResolvedProfile{}, fmt.Errorf("repository-owned Baseline Profile %q requires a Git repository", target)
	}
	return baseline.ResolveProfile(repoRoot, target, catalog)
}

func baselineProfileTargetIsPath(target string) bool {
	return filepath.Ext(target) == ".json" ||
		strings.ContainsRune(target, filepath.Separator) ||
		strings.Contains(target, "/") ||
		strings.Contains(target, `\`)
}

func printBaselineProfileText(profile baseline.ResolvedProfile, stdout io.Writer) {
	fmt.Fprintf(stdout, "Baseline Profile: %s\n", profile.ID)
	fmt.Fprintf(stdout, "Source: %s\n", profile.Source)
	if profile.Path != "" {
		fmt.Fprintf(stdout, "Path: %s\n", profile.Path)
	}
	fmt.Fprintf(stdout, "Catalog schema: %s\n", profile.CatalogSchema)
	fmt.Fprintf(stdout, "Digest: %s\n", profile.Digest)
	printBaselineProfileList(stdout, "Modules", profile.Modules)
	printBaselineProfileList(stdout, "Decisions", profile.Decisions)
	printBaselineProfileList(stdout, "Repository Capabilities", profile.Capabilities)
	printBaselineProfileList(stdout, "Templates", profile.Templates)
	if len(profile.Values) == 0 {
		fmt.Fprintln(stdout, "Values: none")
	} else {
		fmt.Fprintln(stdout, "Values:")
		keys := make([]string, 0, len(profile.Values))
		for key := range profile.Values {
			keys = append(keys, key)
		}
		sort.Strings(keys)
		for _, key := range keys {
			value, _ := json.Marshal(profile.Values[key])
			fmt.Fprintf(stdout, "  %s=%s\n", key, value)
		}
	}
}

func printBaselineProfileList(stdout io.Writer, label string, values []string) {
	if len(values) == 0 {
		fmt.Fprintf(stdout, "%s: none\n", label)
		return
	}
	fmt.Fprintf(stdout, "%s:\n", label)
	for _, value := range values {
		fmt.Fprintf(stdout, "  - %s\n", value)
	}
}

func printBaselineProfileFailure(operation string, err error, jsonOutput bool, stdout, stderr io.Writer) {
	fmt.Fprintf(stderr, "%s: baseline %s failed: %v\n", app.Name, operation, err)
	fmt.Fprintf(stderr, "Run '%s baseline profile --help' for usage.\n", app.Name)
	if !jsonOutput {
		return
	}
	_ = json.NewEncoder(stdout).Encode(baselineProfileResult{
		SchemaVersion: baselineResultSchema,
		Operation:     operation,
		State:         "failed",
		Category:      "preflight",
		Message:       err.Error(),
		NextAction:    "correct the profile input and rerun the command",
	})
}

func encodeBaselineProfileResult(stdout, stderr io.Writer, result baselineProfileResult) int {
	encoder := json.NewEncoder(stdout)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(result); err != nil {
		fmt.Fprintf(stderr, "%s: baseline %s output failed: %v\n", app.Name, result.Operation, err)
		return exitRunFailed
	}
	return exitOK
}

func baselineProfileJSONRequested(args []string) bool {
	for index, arg := range args {
		if arg == "--format" && index+1 < len(args) && strings.TrimSpace(args[index+1]) == "json" {
			return true
		}
		if strings.TrimSpace(strings.TrimPrefix(arg, "--format=")) == "json" && strings.HasPrefix(arg, "--format=") {
			return true
		}
	}
	return false
}

func baselineProfileOutputPath(repoRoot, path string) string {
	if path == "" {
		return ""
	}
	if repoRoot == "" {
		return filepath.ToSlash(path)
	}
	relative, err := filepath.Rel(repoRoot, path)
	if err != nil || relative == ".." || strings.HasPrefix(relative, ".."+string(filepath.Separator)) {
		return filepath.ToSlash(path)
	}
	return filepath.ToSlash(relative)
}

func findBaselineRepositoryRoot() (string, error) {
	workDir, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("resolve current directory: %w", err)
	}
	workDir, err = filepath.Abs(workDir)
	if err != nil {
		return "", fmt.Errorf("resolve current directory: %w", err)
	}
	for {
		gitPath := filepath.Join(workDir, ".git")
		if info, err := os.Lstat(gitPath); err == nil {
			if info.IsDir() || info.Mode().IsRegular() {
				return workDir, nil
			}
			return "", fmt.Errorf("unsafe Git marker %q", gitPath)
		} else if !errors.Is(err, fs.ErrNotExist) {
			return "", fmt.Errorf("inspect Git marker %q: %w", gitPath, err)
		}
		parent := filepath.Dir(workDir)
		if parent == workDir {
			return "", nil
		}
		workDir = parent
	}
}

func isBaselineProfileExecutionError(err error) bool {
	return errors.Is(err, fs.ErrPermission) || errors.Is(err, fs.ErrExist)
}
