package cli

import (
	"bufio"
	"context"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"roundfix/internal/app"
	"roundfix/internal/baseline"
)

const baselineSetupManifestPath = "docs/agents/setup-context.json"

type baselineHumanRequest struct {
	repo   string
	format string
}

type baselineHumanCommandIO struct {
	input       io.Reader
	interactive bool
}

type baselineHumanPrompt struct {
	reader *bufio.Reader
	writer io.Writer
	step   int
}

type baselineHumanState struct {
	mode             string
	currentProfile   *baseline.ResolvedProfile
	currentDecisions map[string]any
	incompatible     string
}

type baselineHumanActionError struct {
	result baseline.Result
}

func (err *baselineHumanActionError) Error() string {
	return err.result.Message
}

type baselineDecisionDeclaration struct {
	ID      string   `json:"id"`
	Type    string   `json:"type"`
	Values  []string `json:"values"`
	Modes   []string `json:"modes"`
	Summary string   `json:"summary"`
}

func runBaselineHumanCommand(ctx context.Context, args []string, stdout, stderr io.Writer) int {
	return runBaselineHumanCommandWithIO(ctx, args, stdout, stderr, defaultBaselineHumanCommandIO())
}

func runBaselineHumanCommandWithIO(
	ctx context.Context,
	args []string,
	stdout io.Writer,
	stderr io.Writer,
	commandIO baselineHumanCommandIO,
) int {
	request, err := parseBaselineHumanCommand(args)
	jsonOutput := request.format == "json" || baselineHumanJSONRequested(args)
	if err != nil {
		return printBaselineHumanFailure(err, jsonOutput, stdout, stderr)
	}
	if !commandIO.interactive {
		result := baselineHumanActionResult(
			"interactive_input",
			"roundfix baseline requires interactive terminal input",
			"run roundfix baseline plan with explicit --profile and --decision inputs, then apply the emitted plan with roundfix baseline apply",
		)
		return writeBaselineHumanAction(result, jsonOutput, stdout, stderr)
	}
	if commandIO.input == nil {
		return printBaselineHumanFailure(
			errors.New("baseline interactive input reader is required"),
			jsonOutput,
			stdout,
			stderr,
		)
	}

	review := stdout
	if jsonOutput {
		review = stderr
	}
	prompt := &baselineHumanPrompt{
		reader: bufio.NewReader(commandIO.input),
		writer: stderr,
	}
	plan, err := driveHumanBaselinePlan(ctx, request.repo, prompt, review)
	if err != nil {
		var actionErr *baselineHumanActionError
		if errors.As(err, &actionErr) {
			return writeBaselineHumanAction(actionErr.result, jsonOutput, stdout, stderr)
		}
		return printBaselineHumanFailure(err, jsonOutput, stdout, stderr)
	}

	selection, err := prompt.selectOne(
		ctx,
		fmt.Sprintf("Final confirmation for Plan Digest %s", plan.PlanDigest),
		[]string{"Apply this exact Plan Digest", "Decline without writing"},
	)
	if err != nil {
		return printBaselineHumanFailure(err, jsonOutput, stdout, stderr)
	}
	if selection != 0 {
		result := baselineHumanActionResult(
			"approval",
			"Baseline Plan was declined; no repository bytes were written",
			"rerun roundfix baseline to revise or approve a new Plan Digest",
		)
		result.PlanDigest = plan.PlanDigest
		return writeBaselineHumanAction(result, jsonOutput, stdout, stderr)
	}

	result, err := baseline.ApplyPlan(ctx, request.repo, plan, plan.PlanDigest)
	if err != nil {
		return printBaselineApplyFailure(plan, err, jsonOutput, stdout, stderr)
	}
	if err := writeBaselineApplyResult(result, jsonOutput, stdout); err != nil {
		fmt.Fprintf(stderr, "%s: baseline output failed: %v\n", app.Name, err)
		return exitRunFailed
	}
	return exitOK
}

func defaultBaselineHumanCommandIO() baselineHumanCommandIO {
	info, err := os.Stdin.Stat()
	return baselineHumanCommandIO{
		input:       os.Stdin,
		interactive: err == nil && info.Mode()&os.ModeCharDevice != 0,
	}
}

func parseBaselineHumanCommand(args []string) (baselineHumanRequest, error) {
	request := baselineHumanRequest{repo: ".", format: "text"}
	flags := flag.NewFlagSet("baseline", flag.ContinueOnError)
	flags.SetOutput(io.Discard)
	flags.StringVar(&request.repo, "repo", ".", "Git worktree or a path inside it")
	flags.StringVar(&request.format, "format", "text", "Output format: text or json")
	if err := flags.Parse(args); err != nil {
		return baselineHumanRequest{}, validationError{
			message: fmt.Sprintf("invalid baseline arguments: %v; run '%s baseline --help' for usage", err, app.Name),
		}
	}
	if remaining := flags.Args(); len(remaining) != 0 {
		return baselineHumanRequest{}, validationError{
			message: fmt.Sprintf("unexpected argument %q; run '%s baseline --help' for usage", remaining[0], app.Name),
		}
	}
	request.repo = strings.TrimSpace(request.repo)
	request.format = strings.TrimSpace(request.format)
	if request.repo == "" {
		return request, validationError{message: "--repo cannot be empty"}
	}
	if request.format != "text" && request.format != "json" {
		return request, validationError{
			message: fmt.Sprintf("unsupported --format %q; use text or json", request.format),
		}
	}
	return request, nil
}

func baselineHumanJSONRequested(args []string) bool {
	for index, arg := range args {
		if arg == "--format=json" {
			return true
		}
		if arg == "--format" && index+1 < len(args) && strings.TrimSpace(args[index+1]) == "json" {
			return true
		}
	}
	return false
}

func driveHumanBaselinePlan(
	ctx context.Context,
	repository string,
	prompt *baselineHumanPrompt,
	review io.Writer,
) (baseline.PlanDocument, error) {
	if prompt == nil {
		return baseline.PlanDocument{}, errors.New("drive human Baseline workflow: prompt adapter is required")
	}
	catalog, err := baseline.LoadEmbeddedCatalog()
	if err != nil {
		return baseline.PlanDocument{}, fmt.Errorf("load Baseline catalog: %w", err)
	}
	inspection, err := baseline.InspectRepository(ctx, repository, nil)
	if err != nil {
		return baseline.PlanDocument{}, fmt.Errorf("inspect Baseline repository: %w", err)
	}
	if len(inspection.Snapshot.Blocking) != 0 {
		return baseline.PlanDocument{}, &baselineHumanActionError{result: baseline.Result{
			SchemaVersion:      baseline.ResultSchemaVersion,
			Operation:          "baseline",
			State:              "action_required",
			Category:           "preflight",
			Message:            "repository preflight found unsafe bounded carriers",
			NextAction:         "repair every blocking carrier and rerun roundfix baseline",
			VerifiedPostimages: []baseline.Postimage{},
			Warnings:           append(inspection.Snapshot.Warnings, inspection.Snapshot.Blocking...),
			Recommendations:    []string{},
		}}
	}

	state, err := inspectBaselineHumanState(inspection.Root, catalog)
	if err != nil {
		return baseline.PlanDocument{}, err
	}
	renderBaselineHumanState(review, state)

	preservationMode, err := promptPreservationMode(ctx, prompt)
	if err != nil {
		return baseline.PlanDocument{}, err
	}
	profile, err := promptBaselineProfile(ctx, prompt, review, inspection.Root, catalog, state)
	if err != nil {
		return baseline.PlanDocument{}, err
	}
	decisions, err := promptBaselineDecisions(ctx, prompt, catalog, profile, state)
	if err != nil {
		return baseline.PlanDocument{}, err
	}
	preservation, err := promptBaselineClassification(
		ctx,
		prompt,
		review,
		inspection,
		preservationMode,
	)
	if err != nil {
		return baseline.PlanDocument{}, err
	}

	outcome, err := baseline.BuildPlan(ctx, baseline.PlanRequest{
		Repository:   inspection.Root,
		ProfileID:    profile.ID,
		Decisions:    decisions,
		Preservation: preservation,
	})
	if err != nil {
		return baseline.PlanDocument{}, fmt.Errorf("build human Baseline Plan: %w", err)
	}
	if outcome.Plan == nil {
		outcome.Result.Operation = "baseline"
		outcome.Result.NextAction = humanBaselineNextAction(outcome.Result)
		return baseline.PlanDocument{}, &baselineHumanActionError{result: outcome.Result}
	}
	printConsolidatedBaselineReview(*outcome.Plan, review)
	return *outcome.Plan, nil
}

func inspectBaselineHumanState(root string, catalog *baseline.Catalog) (baselineHumanState, error) {
	state := baselineHumanState{
		mode:             "adoption",
		currentDecisions: map[string]any{},
	}
	data, err := os.ReadFile(filepath.Join(root, filepath.FromSlash(baselineSetupManifestPath)))
	if errors.Is(err, fs.ErrNotExist) {
		return state, nil
	}
	if err != nil {
		return state, fmt.Errorf("read current Setup Manifest: %w", err)
	}
	var manifest baseline.SetupManifest
	if err := json.Unmarshal(data, &manifest); err != nil {
		state.incompatible = "the existing Setup Manifest is invalid JSON"
		return state, nil
	}
	if manifest.SchemaVersion != baseline.ManifestSchema ||
		manifest.Version != baseline.ManifestVersion ||
		strings.TrimSpace(manifest.Profile) == "" {
		state.incompatible = "the existing Setup Manifest is incompatible with the current Baseline"
		return state, nil
	}
	profile, err := baseline.ResolveProfile(root, manifest.Profile, catalog)
	if err != nil || profile.Digest != manifest.ProfileDigest {
		state.incompatible = "the existing Setup Manifest references an unavailable or changed Baseline Profile"
		return state, nil
	}
	state.mode = "update"
	state.currentProfile = &profile
	for id, decision := range manifest.Decisions {
		state.currentDecisions[id] = decision.Value
	}
	return state, nil
}

func renderBaselineHumanState(output io.Writer, state baselineHumanState) {
	if state.mode == "update" {
		fmt.Fprintf(output, "Baseline workflow: update\nCurrent Baseline Profile: %s\n", state.currentProfile.ID)
		return
	}
	fmt.Fprintln(output, "Baseline workflow: adoption")
	if state.incompatible != "" {
		fmt.Fprintf(output, "Existing state: incompatible — %s\n", state.incompatible)
	} else {
		fmt.Fprintln(output, "Existing state: unconfigured")
	}
}

func promptPreservationMode(
	ctx context.Context,
	prompt *baselineHumanPrompt,
) (baseline.PreservationMode, error) {
	selected, err := prompt.selectOne(ctx, "Instruction preservation", []string{
		"Greenfield: back up root instruction carriers without importing their rules",
		"Preservation: back up root instruction carriers and review every proposed classification",
	})
	if err != nil {
		return "", err
	}
	if selected == 0 {
		return baseline.PreservationModeGreenfield, nil
	}
	return baseline.PreservationModePreservation, nil
}

func promptBaselineProfile(
	ctx context.Context,
	prompt *baselineHumanPrompt,
	review io.Writer,
	root string,
	catalog *baseline.Catalog,
	state baselineHumanState,
) (baseline.ResolvedProfile, error) {
	if state.currentProfile != nil {
		selected, err := prompt.selectOne(ctx, "Baseline Profile", []string{
			"Keep current profile " + state.currentProfile.ID,
			"Change Baseline Profile",
		})
		if err != nil {
			return baseline.ResolvedProfile{}, err
		}
		if selected == 0 {
			return *state.currentProfile, nil
		}
		fmt.Fprintln(review, "Profile change requested; the next Plan Digest will bind the replacement profile.")
	}

	profiles, err := availableBaselineProfiles(root, catalog)
	if err != nil {
		return baseline.ResolvedProfile{}, err
	}
	options := make([]string, len(profiles))
	for index, profile := range profiles {
		options[index] = fmt.Sprintf("%s (%s)", profile.ID, profile.Source)
	}
	selected, err := prompt.selectOne(ctx, "Select exactly one Baseline Profile", options)
	if err != nil {
		return baseline.ResolvedProfile{}, err
	}
	return profiles[selected], nil
}

func availableBaselineProfiles(root string, catalog *baseline.Catalog) ([]baseline.ResolvedProfile, error) {
	ids := catalog.ProfileIDs()
	profiles := make([]baseline.ResolvedProfile, 0, len(ids))
	for _, id := range ids {
		profile, err := baseline.ResolveProfile(root, id, catalog)
		if err != nil {
			return nil, fmt.Errorf("resolve built-in Baseline Profile %q: %w", id, err)
		}
		profiles = append(profiles, profile)
	}
	custom, err := baseline.DiscoverRepositoryProfiles(root, catalog)
	if err != nil {
		return nil, fmt.Errorf("discover repository-owned Baseline Profiles: %w", err)
	}
	return append(profiles, custom...), nil
}

func promptBaselineDecisions(
	ctx context.Context,
	prompt *baselineHumanPrompt,
	catalog *baseline.Catalog,
	profile baseline.ResolvedProfile,
	state baselineHumanState,
) ([]baseline.DecisionValue, error) {
	answers := make([]baseline.DecisionValue, 0, len(profile.Decisions))
	answered := make(map[string]struct{}, len(profile.Decisions))
	for _, id := range profile.Decisions {
		if _, fixed := profile.Values[id]; fixed {
			continue
		}
		value, err := promptBaselineDecision(ctx, prompt, catalog, id, state.currentDecisions)
		if err != nil {
			return nil, err
		}
		answers = append(answers, baseline.DecisionValue{ID: id, Value: value})
		answered[id] = struct{}{}
	}
	for {
		normalized, missing, err := baseline.ResolveDecisionInput(profile, answers, catalog)
		if err != nil {
			return nil, err
		}
		if len(missing) == 0 {
			return normalized, nil
		}
		for _, id := range missing {
			if _, exists := answered[id]; exists {
				return nil, fmt.Errorf("resolve Baseline decisions: %q remains missing after an answer", id)
			}
			value, err := promptBaselineDecision(ctx, prompt, catalog, id, state.currentDecisions)
			if err != nil {
				return nil, err
			}
			answers = append(answers, baseline.DecisionValue{ID: id, Value: value})
			answered[id] = struct{}{}
		}
	}
}

func promptBaselineDecision(
	ctx context.Context,
	prompt *baselineHumanPrompt,
	catalog *baseline.Catalog,
	id string,
	current map[string]any,
) (any, error) {
	declaration, err := baselineDecisionDefinition(catalog, id)
	if err != nil {
		return nil, err
	}
	if value, ok := current[id]; ok {
		encoded, _ := json.Marshal(value)
		selected, err := prompt.selectOne(ctx, declaration.Summary, []string{
			fmt.Sprintf("Keep %s=%s", id, encoded),
			"Change " + id,
		})
		if err != nil {
			return nil, err
		}
		if selected == 0 {
			return value, nil
		}
	}

	switch declaration.Type {
	case "boolean":
		selected, err := prompt.selectOne(ctx, declaration.Summary, []string{"Yes", "No"})
		if err != nil {
			return nil, err
		}
		return selected == 0, nil
	case "enum":
		selected, err := prompt.selectOne(ctx, declaration.Summary, declaration.Values)
		if err != nil {
			return nil, err
		}
		return declaration.Values[selected], nil
	case "http-contract":
		selected, err := prompt.selectOne(ctx, declaration.Summary, declaration.Modes)
		if err != nil {
			return nil, err
		}
		return map[string]any{"mode": declaration.Modes[selected]}, nil
	case "string":
		return prompt.readNonEmpty(ctx, declaration.Summary+" ("+id+")")
	default:
		return nil, fmt.Errorf("prompt Baseline decision %q: unsupported type %q", id, declaration.Type)
	}
}

func baselineDecisionDefinition(
	catalog *baseline.Catalog,
	id string,
) (baselineDecisionDeclaration, error) {
	entry, ok := catalog.Decision(id)
	if !ok {
		return baselineDecisionDeclaration{}, fmt.Errorf("unknown Baseline decision %q", id)
	}
	var declaration baselineDecisionDeclaration
	if err := json.Unmarshal(entry.Data, &declaration); err != nil {
		return baselineDecisionDeclaration{}, fmt.Errorf("decode Baseline decision %q: %w", id, err)
	}
	return declaration, nil
}

func promptBaselineClassification(
	ctx context.Context,
	prompt *baselineHumanPrompt,
	review io.Writer,
	inspection baseline.RepositoryInspection,
	mode baseline.PreservationMode,
) (baseline.RootPreservationRequest, error) {
	request := baseline.RootPreservationRequest{Mode: mode}
	if mode != baseline.PreservationModePreservation {
		return request, nil
	}
	plan, err := baseline.PlanRootPreservation(inspection, request)
	if err != nil {
		return baseline.RootPreservationRequest{}, fmt.Errorf("prepare consolidated classification review: %w", err)
	}
	if plan.State == baseline.PreservationStateBlocked {
		return baseline.RootPreservationRequest{}, &baselineHumanActionError{
			result: baselineHumanActionResult("preflight", plan.NextAction, plan.NextAction),
		}
	}
	if plan.DecisionSkeleton == nil {
		return request, nil
	}
	document := plan.DecisionSkeleton.Document
	fmt.Fprintln(review, "\nConsolidated editable classification review:")
	for index, disposition := range document.Readoption.Dispositions {
		entry := plan.SourceBaseline.Entries[index]
		fmt.Fprintf(
			review,
			"%d. %s bytes %d-%d: %s -> %s\n",
			index+1,
			entry.Path,
			entry.Start,
			entry.End,
			disposition.Classification,
			disposition.Disposition,
		)
	}
	selected, err := prompt.selectOne(ctx, "Classification review", []string{
		"Accept the complete proposal",
		"Edit classifications and dispositions",
	})
	if err != nil {
		return baseline.RootPreservationRequest{}, err
	}
	if selected == 1 {
		for index := range document.Readoption.Dispositions {
			disposition := &document.Readoption.Dispositions[index]
			entry := plan.SourceBaseline.Entries[index]
			classification, err := prompt.selectOne(ctx, "Classification for "+entry.Path, []string{
				"Normative Clause",
				"Operational Contract",
				"Recommendation",
				"Non-governed evidence",
			})
			if err != nil {
				return baseline.RootPreservationRequest{}, err
			}
			switch classification {
			case 0:
				disposition.Classification = "normative-clause"
			case 1:
				disposition.Classification = "operational-contract"
			case 2:
				disposition.Classification = "recommendation"
			default:
				disposition.Classification = "non-governed"
			}
			keep := 1
			canPreserve := entry.Kind != "managed-block" &&
				len(strings.TrimSpace(string(entry.SourceBytes))) != 0
			if disposition.Classification != "non-governed" && canPreserve {
				keep, err = prompt.selectOne(ctx, "Disposition for "+entry.Path, []string{
					"Preserve in Repository-Specific Normative Rules",
					"Reject with an individual reason",
				})
				if err != nil {
					return baseline.RootPreservationRequest{}, err
				}
			}
			if disposition.Classification == "non-governed" || keep == 1 {
				reason, err := prompt.readNonEmpty(ctx, "Reason for rejecting "+entry.Path)
				if err != nil {
					return baseline.RootPreservationRequest{}, err
				}
				disposition.Disposition = "rejected"
				disposition.Destination = nil
				disposition.Reason = reason
				continue
			}
			sum := sha256.Sum256(entry.SourceBytes)
			disposition.Disposition = "repository-rules"
			disposition.Destination = &baseline.ReadoptionDestination{
				DocumentType:  "repository-rules",
				Path:          "docs/agents/repository-rules.md",
				Digest:        hex.EncodeToString(sum[:]),
				ProposedBytes: base64.StdEncoding.EncodeToString(entry.SourceBytes),
			}
			disposition.Reason = "Preserve this source entry as a Repository-Specific Normative Rule after human review."
		}
	}
	request.Decisions = &document
	return request, nil
}

func printConsolidatedBaselineReview(plan baseline.PlanDocument, output io.Writer) {
	fmt.Fprintln(output, "\nConsolidated Change Plan review")
	fmt.Fprintln(output, "File changes:")
	for index, change := range plan.FileChanges {
		fmt.Fprintf(output, "%d. %s %s\n", index+1, change.Action, change.Path)
		fmt.Fprintf(output, "   Before: %s\n", change.BeforeIdentity)
		fmt.Fprintf(output, "   After: %s\n", change.AfterIdentity)
		fmt.Fprintf(output, "   Managed entries: %s\n", strings.Join(change.ManagedEntries, ", "))
	}
	fmt.Fprintln(output, "Complete managed-entry ledger:")
	for _, entry := range plan.ManagedEntries {
		fmt.Fprintf(output, "%d. %s %s %s\n", entry.Ordinal, entry.Action, entry.Path, entry.ID)
	}
	fmt.Fprintln(output, "Complete Upgrade Retention Contract ledger:")
	if len(plan.Retention) == 0 {
		fmt.Fprintln(output, "0 entries")
	}
	for index, retention := range plan.Retention {
		fmt.Fprintf(
			output,
			"%d. %s -> %s (%s)\n",
			index+1,
			retention.FromClause,
			retention.Disposition,
			strings.Join(retention.Targets, ", "),
		)
	}
	for _, warning := range plan.Warnings {
		fmt.Fprintf(output, "Warning: %s: %s: %s\n", warning.Code, warning.Path, warning.Message)
	}
	fmt.Fprintf(output, "Plan Digest: %s\n", plan.PlanDigest)
}

func (prompt *baselineHumanPrompt) selectOne(
	ctx context.Context,
	label string,
	options []string,
) (int, error) {
	if len(options) == 0 {
		return 0, fmt.Errorf("prompt %q has no options", label)
	}
	prompt.step++
	fmt.Fprintf(prompt.writer, "\nPrompt %d: %s\n", prompt.step, label)
	for index, option := range options {
		fmt.Fprintf(prompt.writer, "  %d. %s\n", index+1, option)
	}
	for {
		line, err := prompt.readLine(ctx)
		if err != nil {
			return 0, err
		}
		choice, err := strconv.Atoi(strings.TrimSpace(line))
		if err == nil && choice >= 1 && choice <= len(options) {
			return choice - 1, nil
		}
		fmt.Fprintf(prompt.writer, "Enter a number from 1 to %d: ", len(options))
	}
}

func (prompt *baselineHumanPrompt) readNonEmpty(ctx context.Context, label string) (string, error) {
	prompt.step++
	fmt.Fprintf(prompt.writer, "\nPrompt %d: %s\n> ", prompt.step, label)
	for {
		line, err := prompt.readLine(ctx)
		if err != nil {
			return "", err
		}
		if value := strings.TrimSpace(line); value != "" {
			return value, nil
		}
		fmt.Fprint(prompt.writer, "Enter a non-empty value: ")
	}
}

func (prompt *baselineHumanPrompt) readLine(ctx context.Context) (string, error) {
	if err := ctx.Err(); err != nil {
		return "", err
	}
	line, err := prompt.reader.ReadString('\n')
	if err != nil && !errors.Is(err, io.EOF) {
		return "", fmt.Errorf("read Baseline prompt: %w", err)
	}
	if errors.Is(err, io.EOF) && line == "" {
		return "", errors.New("interactive Baseline input ended before the workflow was complete")
	}
	if err := ctx.Err(); err != nil {
		return "", err
	}
	return line, nil
}

func baselineHumanActionResult(category, message, nextAction string) baseline.Result {
	return baseline.Result{
		SchemaVersion:      baseline.ResultSchemaVersion,
		Operation:          "baseline",
		State:              "action_required",
		Category:           category,
		Message:            message,
		NextAction:         nextAction,
		VerifiedPostimages: []baseline.Postimage{},
		Warnings:           []baseline.Finding{},
		Recommendations:    []string{},
	}
}

func humanBaselineNextAction(result baseline.Result) string {
	switch result.Category {
	case "classification":
		return "review the reported classification state, then rerun roundfix baseline"
	case "preflight":
		return "repair the reported repository state, then rerun roundfix baseline"
	default:
		return "revise the selected profile or repository decisions, then rerun roundfix baseline"
	}
}

func writeBaselineHumanAction(
	result baseline.Result,
	jsonOutput bool,
	stdout io.Writer,
	stderr io.Writer,
) int {
	if err := writeBaselinePlanResult(result, jsonOutput, stdout); err != nil {
		fmt.Fprintf(stderr, "%s: baseline output failed: %v\n", app.Name, err)
		return exitRunFailed
	}
	if result.Category == "preflight" {
		for _, finding := range result.Warnings {
			fmt.Fprintf(stderr, "%s: %s: %s: %s\n", app.Name, finding.Code, finding.Path, finding.Message)
		}
		return exitPreflight
	}
	return exitUnverified
}

func printBaselineHumanFailure(
	err error,
	jsonOutput bool,
	stdout io.Writer,
	stderr io.Writer,
) int {
	fmt.Fprintf(stderr, "%s: baseline failed: %v\n", app.Name, err)
	result := baselineHumanActionResult(
		"execution",
		err.Error(),
		"correct the reported failure and rerun roundfix baseline",
	)
	exit := exitRunFailed
	var validation validationError
	switch {
	case errors.As(err, &validation):
		result.Category = "usage"
		result.NextAction = "correct the command input and rerun roundfix baseline"
		exit = exitPreflight
	case errors.Is(err, context.Canceled), errors.Is(err, context.DeadlineExceeded):
		result.Category = "canceled"
		result.NextAction = "rerun roundfix baseline when the operation can complete"
		exit = exitSIGINT
	}
	if jsonOutput {
		if writeErr := writeBaselinePlanResult(result, true, stdout); writeErr != nil {
			fmt.Fprintf(stderr, "%s: baseline output failed: %v\n", app.Name, writeErr)
			return exitRunFailed
		}
	}
	return exit
}
