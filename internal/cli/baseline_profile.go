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

func runBaselineCommand(ctx context.Context, args []string, stdout, stderr io.Writer) int {
	if len(args) == 1 && commandWantsHelp(args) {
		fmt.Fprint(stdout, commandUsage("baseline"))
		return exitOK
	}
	if len(args) == 0 || strings.HasPrefix(args[0], "-") {
		return runBaselineHumanCommand(ctx, args, stdout, stderr)
	}
	switch args[0] {
	case "plan":
		return runBaselinePlanCommand(ctx, args[1:], stdout, stderr)
	case "apply":
		return runBaselineApplyCommand(ctx, args[1:], stdout, stderr)
	case "profile":
		return runBaselineProfileCommand(args[1:], stdout, stderr)
	case "skills":
		if len(args) == 1 || args[1] != "restore" {
			printBaselineProfileFailure(
				"baseline skills",
				validationError{message: "baseline skills requires the restore command"},
				false,
				stdout,
				stderr,
			)
			return exitPreflight
		}
		return runBaselineSkillsRestoreCommand(ctx, args[2:], stdout, stderr)
	case "assets":
		if len(args) == 1 || args[1] != "sync" {
			printBaselineProfileFailure(
				"baseline assets",
				validationError{message: "baseline assets requires the sync command"},
				false,
				stdout,
				stderr,
			)
			return exitPreflight
		}
		return runBaselineAssetsSyncCommand(ctx, args[2:], stdout, stderr)
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
		printBaselinePlanFailure(err, jsonOutput, stdout, stderr)
		return exitPreflight
	}
	decisionInput, preservation, err := loadBaselinePlanDecisions(req)
	if err != nil {
		printBaselinePlanFailure(err, jsonOutput, stdout, stderr)
		return exitPreflight
	}
	outcome, err := baseline.BuildPlan(ctx, baseline.PlanRequest{
		Repository:   req.repo,
		ProfileID:    req.profile,
		Decisions:    decisionInput,
		Preservation: preservation,
	})
	if err != nil {
		printBaselinePlanFailure(err, jsonOutput, stdout, stderr)
		return exitPreflight
	}
	if outcome.Plan == nil {
		if err := writeBaselinePlanResult(outcome.Result, jsonOutput, stdout); err != nil {
			fmt.Fprintf(stderr, "%s: baseline plan output failed: %v\n", app.Name, err)
			return exitRunFailed
		}
		if outcome.Result.Category == "preflight" {
			for _, finding := range outcome.Result.Warnings {
				fmt.Fprintf(stderr, "%s: %s: %s: %s\n", app.Name, finding.Code, finding.Path, finding.Message)
			}
			return exitPreflight
		}
		return exitUnverified
	}
	if jsonOutput {
		data, err := baseline.MarshalPlanDocument(*outcome.Plan)
		if err != nil {
			fmt.Fprintf(stderr, "%s: baseline plan output failed: %v\n", app.Name, err)
			return exitRunFailed
		}
		if _, err := stdout.Write(data); err != nil {
			fmt.Fprintf(stderr, "%s: baseline plan output failed: %v\n", app.Name, err)
			return exitRunFailed
		}
		return exitOK
	}
	printBaselinePlanText(*outcome.Plan, stdout)
	return exitOK
}

type baselinePlanCommandRequest struct {
	repo          string
	profile       string
	format        string
	decisions     stringListFlag
	decisionFiles stringListFlag
}

func parseBaselinePlanCommand(args []string) (baselinePlanCommandRequest, error) {
	req := baselinePlanCommandRequest{repo: ".", format: "text"}
	flags := flag.NewFlagSet("baseline plan", flag.ContinueOnError)
	flags.SetOutput(io.Discard)
	flags.StringVar(&req.repo, "repo", ".", "Git worktree or a path inside it")
	flags.StringVar(&req.profile, "profile", "", "Built-in or repository-owned Baseline Profile")
	flags.Var(&req.decisions, "decision", "Baseline decision as id=value (repeatable)")
	flags.Var(&req.decisionFiles, "decision-file", "Strict Decision Document path (repeatable)")
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
	req.profile = strings.TrimSpace(req.profile)
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

type stringListFlag []string

func (values *stringListFlag) String() string {
	return strings.Join(*values, ",")
}

func (values *stringListFlag) Set(value string) error {
	*values = append(*values, value)
	return nil
}

func loadBaselinePlanDecisions(
	request baselinePlanCommandRequest,
) ([]baseline.DecisionValue, baseline.RootPreservationRequest, error) {
	var decisions []baseline.DecisionValue
	preservation := baseline.RootPreservationRequest{}
	for _, raw := range request.decisions {
		id, value, ok := strings.Cut(raw, "=")
		id = strings.TrimSpace(id)
		if !ok || id == "" {
			return nil, baseline.RootPreservationRequest{}, validationError{
				message: fmt.Sprintf("invalid --decision %q; expected id=value", raw),
			}
		}
		if id == "preservation.mode" {
			mode := baseline.PreservationMode(strings.TrimSpace(value))
			if mode != baseline.PreservationModeGreenfield && mode != baseline.PreservationModePreservation {
				return nil, baseline.RootPreservationRequest{}, validationError{
					message: "preservation.mode must be greenfield or preservation",
				}
			}
			if preservation.Mode != "" {
				return nil, baseline.RootPreservationRequest{}, validationError{
					message: "preservation.mode may be specified only once",
				}
			}
			preservation.Mode = mode
			continue
		}
		decisions = append(decisions, baseline.DecisionValue{
			ID:    id,
			Value: parseBaselineInlineDecision(value),
		})
	}
	for _, decisionPath := range request.decisionFiles {
		data, err := os.ReadFile(decisionPath)
		if err != nil {
			return nil, baseline.RootPreservationRequest{}, fmt.Errorf(
				"read decision file %q: %w", decisionPath, err,
			)
		}
		document, err := baseline.ParseDecisionDocument(data, decisionPath)
		if err != nil {
			return nil, baseline.RootPreservationRequest{}, err
		}
		for _, decision := range document.Decisions {
			if decision.ID != "preservation.mode" {
				decisions = append(decisions, decision)
				continue
			}
			mode, ok := decision.Value.(string)
			if !ok {
				return nil, baseline.RootPreservationRequest{}, validationError{
					message: "preservation.mode in a decision file must be a string",
				}
			}
			normalized := baseline.PreservationMode(strings.TrimSpace(mode))
			if normalized != baseline.PreservationModeGreenfield &&
				normalized != baseline.PreservationModePreservation {
				return nil, baseline.RootPreservationRequest{}, validationError{
					message: "preservation.mode must be greenfield or preservation",
				}
			}
			if preservation.Mode != "" {
				return nil, baseline.RootPreservationRequest{}, validationError{
					message: "preservation.mode may be specified only once",
				}
			}
			preservation.Mode = normalized
		}
		if document.Readoption != nil {
			if preservation.Decisions != nil {
				return nil, baseline.RootPreservationRequest{}, validationError{
					message: "readoption decisions may appear in only one decision file",
				}
			}
			copy := document
			preservation.Decisions = &copy
		}
	}
	return decisions, preservation, nil
}

func parseBaselineInlineDecision(value string) any {
	trimmed := strings.TrimSpace(value)
	var parsed any
	if json.Unmarshal([]byte(trimmed), &parsed) == nil {
		switch parsed.(type) {
		case bool, map[string]any, []any, float64:
			return parsed
		}
	}
	return trimmed
}

func printBaselinePlanText(plan baseline.PlanDocument, stdout io.Writer) {
	fmt.Fprintln(stdout, "Baseline Plan: ready")
	fmt.Fprintf(stdout, "Plan Digest: %s\n", plan.PlanDigest)
	fmt.Fprintf(stdout, "Repository identity: %s\n", plan.Repository.Digest)
	fmt.Fprintf(stdout, "Catalog identity: %s\n", plan.Catalog.Digest)
	fmt.Fprintf(stdout, "Baseline Profile: %s\n", plan.Profile.ID)
	fmt.Fprintf(stdout, "Managed entries: %d\n", len(plan.ManagedEntries))
	fmt.Fprintf(stdout, "File changes: %d\n", len(plan.FileChanges))
	for _, change := range plan.FileChanges {
		fmt.Fprintf(stdout, "- %s %s (%d managed entries)\n",
			change.Action, change.Path, len(change.ManagedEntries))
	}
	for _, warning := range plan.Warnings {
		fmt.Fprintf(stdout, "Warning: %s: %s: %s\n", warning.Code, warning.Path, warning.Message)
	}
}

func printBaselinePlanFailure(
	err error,
	jsonOutput bool,
	stdout io.Writer,
	stderr io.Writer,
) {
	fmt.Fprintf(stderr, "%s: baseline plan failed: %v\n", app.Name, err)
	fmt.Fprintf(stderr, "Run '%s baseline plan --help' for usage.\n", app.Name)
	if !jsonOutput {
		return
	}
	result := baseline.Result{
		SchemaVersion:      baseline.ResultSchemaVersion,
		Operation:          "plan",
		State:              "failed",
		Category:           "preflight",
		Message:            err.Error(),
		NextAction:         "correct the repository or command input and rerun roundfix baseline plan",
		VerifiedPostimages: []baseline.Postimage{},
		Warnings:           []baseline.Finding{},
		Recommendations:    []string{},
	}
	if encodeErr := writeBaselinePlanResult(result, true, stdout); encodeErr != nil {
		fmt.Fprintf(stderr, "%s: baseline plan output failed: %v\n", app.Name, encodeErr)
	}
}

func writeBaselinePlanResult(result baseline.Result, jsonOutput bool, stdout io.Writer) error {
	if jsonOutput {
		data, err := baseline.MarshalResult(result)
		if err != nil {
			return err
		}
		_, err = stdout.Write(data)
		return err
	}
	fmt.Fprintf(stdout, "Baseline plan: %s\n", strings.ReplaceAll(result.State, "_", " "))
	fmt.Fprintf(stdout, "Category: %s\n", result.Category)
	fmt.Fprintf(stdout, "Result: %s\n", result.Message)
	fmt.Fprintf(stdout, "Next action: %s\n", result.NextAction)
	for _, warning := range result.Warnings {
		fmt.Fprintf(stdout, "Warning: %s: %s: %s\n", warning.Code, warning.Path, warning.Message)
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
