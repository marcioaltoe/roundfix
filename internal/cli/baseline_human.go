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
	"roundfix/internal/baselineacp"
)

const baselineSetupManifestPath = "docs/agents/setup-context.json"

type baselineHumanRequest struct {
	repo   string
	format string
}

type baselineHumanCommandIO struct {
	input            io.Reader
	interactive      bool
	revisionAnalyzer baselineRevisionAnalyzer
}

type baselineRevisionAnalyzer interface {
	Revise(context.Context, baseline.RevisionSnapshot) (baseline.RevisionProposal, error)
}

type baselineHumanPrompt struct {
	reader *bufio.Reader
	writer io.Writer
	step   int
}

type baselineHumanOutput struct {
	writer io.Writer
	err    error
}

func (output *baselineHumanOutput) Write(data []byte) (int, error) {
	if output.err != nil {
		return 0, output.err
	}
	written, err := output.writer.Write(data)
	if err != nil {
		output.err = err
	}
	return written, err
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
	ID         string   `json:"id"`
	Type       string   `json:"type"`
	Values     []string `json:"values"`
	Modes      []string `json:"modes"`
	Summary    string   `json:"summary"`
	Default    any      `json:"default"`
	Suggestion any      `json:"suggestion"`
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

	reviewDestination := stdout
	if jsonOutput {
		reviewDestination = stderr
	}
	review := &baselineHumanOutput{writer: reviewDestination}
	promptOutput := &baselineHumanOutput{writer: stderr}
	prompt := &baselineHumanPrompt{
		reader: bufio.NewReader(commandIO.input),
		writer: promptOutput,
	}
	var planRequest baseline.PlanRequest
	plan, err := driveHumanBaselinePlanWithRequest(ctx, request.repo, prompt, review, &planRequest)
	if err != nil {
		var actionErr *baselineHumanActionError
		if errors.As(err, &actionErr) {
			return writeBaselineHumanAction(actionErr.result, jsonOutput, stdout, stderr)
		}
		return printBaselineHumanFailure(err, jsonOutput, stdout, stderr)
	}
	if err := baselineHumanOutputFailure(review, promptOutput); err != nil {
		return printBaselineHumanFailure(
			fmt.Errorf("write Baseline review: %w", err),
			jsonOutput,
			stdout,
			stderr,
		)
	}

	originalPlan := plan
	for {
		selection, selectionErr := prompt.selectOne(
			ctx,
			fmt.Sprintf("Final confirmation for Plan Digest %s", plan.PlanDigest),
			[]string{
				"Apply this exact Plan Digest",
				"Decline without writing",
				"Reject and revise one decision area",
			},
		)
		if selectionErr != nil {
			return printBaselineHumanFailure(selectionErr, jsonOutput, stdout, stderr)
		}
		if selection == 0 {
			break
		}
		if selection == 1 {
			result := baselineHumanActionResult(
				"approval",
				"Baseline Plan was declined; no repository bytes were written",
				"rerun roundfix baseline to revise or approve a new Plan Digest",
			)
			result.PlanDigest = plan.PlanDigest
			return writeBaselineHumanAction(result, jsonOutput, stdout, stderr)
		}
		revisedPlan, revisedRequest, revisionErr := reviseHumanBaselinePlan(
			ctx,
			prompt,
			review,
			originalPlan,
			plan,
			planRequest,
			commandIO.revisionAnalyzer,
		)
		if revisionErr != nil {
			return printBaselineHumanFailure(revisionErr, jsonOutput, stdout, stderr)
		}
		plan = revisedPlan
		planRequest = revisedRequest
	}
	if err := baselineHumanOutputFailure(review, promptOutput); err != nil {
		return printBaselineHumanFailure(
			fmt.Errorf("write Baseline confirmation: %w", err),
			jsonOutput,
			stdout,
			stderr,
		)
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

func baselineHumanOutputFailure(outputs ...*baselineHumanOutput) error {
	for _, output := range outputs {
		if output != nil && output.err != nil {
			return output.err
		}
	}
	return nil
}

func defaultBaselineHumanCommandIO() baselineHumanCommandIO {
	info, err := os.Stdin.Stat()
	return baselineHumanCommandIO{
		input:            os.Stdin,
		interactive:      err == nil && info.Mode()&os.ModeCharDevice != 0,
		revisionAnalyzer: baselineacp.NewDefaultAnalyzer(),
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
	return driveHumanBaselinePlanWithRequest(ctx, repository, prompt, review, nil)
}

func driveHumanBaselinePlanWithRequest(
	ctx context.Context,
	repository string,
	prompt *baselineHumanPrompt,
	review io.Writer,
	captured *baseline.PlanRequest,
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

	preservationMode, err := promptPreservationMode(
		ctx,
		prompt,
		baselineHasRootInstructions(inspection.Snapshot),
	)
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
	profile, decisions, profileDraft, err := promptBaselineProfileAlignment(
		ctx,
		prompt,
		review,
		inspection.Root,
		catalog,
		state,
		profile,
		decisions,
	)
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

	request := baseline.PlanRequest{
		Repository:   inspection.Root,
		Decisions:    decisions,
		Preservation: preservation,
	}
	if profileDraft == nil {
		request.ProfileID = profile.ID
	} else {
		request.ProfileDraft = profileDraft
	}
	outcome, err := baseline.BuildPlan(ctx, request)
	if err != nil {
		return baseline.PlanDocument{}, fmt.Errorf("build human Baseline Plan: %w", err)
	}
	if outcome.Plan == nil {
		outcome.Result.Operation = "baseline"
		outcome.Result.NextAction = humanBaselineNextAction(outcome.Result)
		return baseline.PlanDocument{}, &baselineHumanActionError{result: outcome.Result}
	}
	printConsolidatedBaselineReview(*outcome.Plan, review)
	if captured != nil {
		*captured = request
	}
	return *outcome.Plan, nil
}

func reviseHumanBaselinePlan(
	ctx context.Context,
	prompt *baselineHumanPrompt,
	review io.Writer,
	original baseline.PlanDocument,
	current baseline.PlanDocument,
	request baseline.PlanRequest,
	analyzer baselineRevisionAnalyzer,
) (baseline.PlanDocument, baseline.PlanRequest, error) {
	selected, err := prompt.selectOne(ctx, "Decision area to revisit", []string{
		"Baseline Profile",
		"Repository-Specific Normative Rules",
		"Repository divergences and decisions",
		"Projected files",
	})
	if err != nil {
		return baseline.PlanDocument{}, request, err
	}
	areas := []baseline.RevisionArea{
		baseline.RevisionAreaProfile,
		baseline.RevisionAreaRepositoryRules,
		baseline.RevisionAreaDivergences,
		baseline.RevisionAreaFiles,
	}
	area := areas[selected]
	inputMode, err := prompt.selectOne(ctx, "Revision input", []string{
		"Make structured changes manually",
		"Translate a free-form suggestion within Baseline scope",
	})
	if err != nil {
		return baseline.PlanDocument{}, request, err
	}

	revised := request
	semanticAccepted := false
	if inputMode == 1 {
		suggestion, readErr := prompt.readNonEmpty(ctx, "Scoped Baseline revision suggestion")
		if readErr != nil {
			return baseline.PlanDocument{}, request, readErr
		}
		if analyzer == nil {
			fmt.Fprintln(review, "Semantic revision is unavailable; continue with the structured manual revision.")
		} else {
			snapshot, snapshotErr := baseline.NewRevisionSnapshot(current, area, suggestion)
			if snapshotErr != nil {
				return baseline.PlanDocument{}, request, snapshotErr
			}
			proposal, proposalErr := analyzer.Revise(ctx, snapshot)
			if proposalErr != nil {
				if ctx.Err() != nil {
					return baseline.PlanDocument{}, request, proposalErr
				}
				fmt.Fprintf(review, "Semantic revision was discarded: %v\n", proposalErr)
				fmt.Fprintln(review, "Next action: make the correction through the structured manual revision.")
			} else if proposal.Manual {
				fmt.Fprintln(review, "Semantic revision is unavailable; continue with the structured manual revision.")
			} else {
				decisions, decisionErr := baseline.DecisionsFromRevisionProposal(snapshot, proposal)
				if decisionErr != nil {
					fmt.Fprintf(review, "Semantic revision was discarded: %v\n", decisionErr)
					fmt.Fprintln(review, "Next action: make the correction through the structured manual revision.")
				} else {
					catalog, loadErr := baseline.LoadEmbeddedCatalog()
					if loadErr != nil {
						return baseline.PlanDocument{}, request, loadErr
					}
					profile := current.Profile
					normalized, missing, normalizeErr := baseline.ResolveDecisionInput(profile, decisions, catalog)
					if normalizeErr != nil || len(missing) != 0 {
						fmt.Fprintln(review, "Semantic revision was discarded because it did not produce complete valid Baseline decisions.")
						fmt.Fprintln(review, "Next action: make the correction through the structured manual revision.")
					} else {
						revised.Decisions = normalized
						semanticAccepted = true
					}
				}
			}
		}
	}
	if !semanticAccepted {
		revised, err = promptStructuredBaselineRevision(ctx, prompt, review, current, request, area)
		if err != nil {
			return baseline.PlanDocument{}, request, err
		}
	}
	outcome, err := baseline.RecalculatePlan(ctx, original, revised)
	if err != nil {
		return baseline.PlanDocument{}, request, err
	}
	if outcome.Plan == nil {
		outcome.Result.Operation = "baseline"
		outcome.Result.NextAction = humanBaselineNextAction(outcome.Result)
		return baseline.PlanDocument{}, request, &baselineHumanActionError{result: outcome.Result}
	}
	if outcome.Plan.PlanDigest == current.PlanDigest {
		return baseline.PlanDocument{}, request, errors.New(
			"recalculate Baseline Plan: revision made no decision change; choose a different structured value",
		)
	}
	fmt.Fprintln(review, "\nRejected Plan revision accepted; review the newly computed complete Plan.")
	printConsolidatedBaselineReview(*outcome.Plan, review)
	return *outcome.Plan, revised, nil
}

func promptStructuredBaselineRevision(
	ctx context.Context,
	prompt *baselineHumanPrompt,
	review io.Writer,
	current baseline.PlanDocument,
	request baseline.PlanRequest,
	area baseline.RevisionArea,
) (baseline.PlanRequest, error) {
	catalog, err := baseline.LoadEmbeddedCatalog()
	if err != nil {
		return request, err
	}
	currentValues := make(map[string]any, len(current.Decisions))
	for _, decision := range current.Decisions {
		currentValues[decision.ID] = decision.Value
	}
	switch area {
	case baseline.RevisionAreaProfile:
		profile := current.Profile
		revisionState := baselineHumanState{
			mode:             "update",
			currentProfile:   &profile,
			currentDecisions: currentValues,
		}
		selected, selectErr := promptBaselineProfile(ctx, prompt, review, request.Repository, catalog, revisionState)
		if selectErr != nil {
			return request, selectErr
		}
		decisions, decisionErr := promptBaselineDecisions(ctx, prompt, catalog, selected, revisionState)
		if decisionErr != nil {
			return request, decisionErr
		}
		if request.ProfileDraft == nil || selected.ID != current.Profile.ID {
			request.ProfileID = selected.ID
			request.ProfileDraft = nil
		}
		request.Decisions = decisions
	case baseline.RevisionAreaRepositoryRules:
		inspection, inspectErr := baseline.InspectRepository(ctx, request.Repository, nil)
		if inspectErr != nil {
			return request, inspectErr
		}
		mode, modeErr := promptPreservationMode(
			ctx,
			prompt,
			baselineHasRootInstructions(inspection.Snapshot),
		)
		if modeErr != nil {
			return request, modeErr
		}
		preservation, preservationErr := promptBaselineClassification(ctx, prompt, review, inspection, mode)
		if preservationErr != nil {
			return request, preservationErr
		}
		request.Preservation = preservation
	case baseline.RevisionAreaFiles:
		fmt.Fprintln(review, "Projected files are derived; select the owning Baseline decision to change:")
		for index, change := range current.FileChanges {
			fmt.Fprintf(review, "%d. %s %s\n", index+1, change.Action, change.Path)
		}
		fallthrough
	case baseline.RevisionAreaDivergences:
		decisions, decisionErr := promptOneBaselineDecisionRevision(ctx, prompt, catalog, current.Decisions, currentValues)
		if decisionErr != nil {
			return request, decisionErr
		}
		request.Decisions = decisions
	default:
		return request, fmt.Errorf("structured Baseline revision: unsupported area %q", area)
	}
	return request, nil
}

func promptOneBaselineDecisionRevision(
	ctx context.Context,
	prompt *baselineHumanPrompt,
	catalog *baseline.Catalog,
	decisions []baseline.DecisionValue,
	current map[string]any,
) ([]baseline.DecisionValue, error) {
	options := make([]string, len(decisions))
	for index, decision := range decisions {
		options[index] = decision.ID
	}
	selected, err := prompt.selectOne(ctx, "Baseline decision to change", options)
	if err != nil {
		return nil, err
	}
	value, err := promptBaselineDecision(ctx, prompt, catalog, decisions[selected].ID, current)
	if err != nil {
		return nil, err
	}
	revised := append([]baseline.DecisionValue(nil), decisions...)
	revised[selected].Value = value
	return revised, nil
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
	for id, decision := range manifest.Decisions {
		state.currentDecisions[id] = decision.Value
	}
	profile, err := baseline.ResolveProfile(root, manifest.Profile, catalog)
	if err != nil {
		state.incompatible = "the existing Setup Manifest references an unavailable or changed Baseline Profile"
		return state, nil
	}
	state.currentProfile = &profile
	if profile.Digest != manifest.ProfileDigest {
		state.incompatible = "the existing Setup Manifest references an unavailable or changed Baseline Profile"
		return state, nil
	}
	stored := make([]baseline.DecisionValue, 0, len(manifest.Decisions))
	for id, decision := range manifest.Decisions {
		if _, fixed := profile.Values[id]; fixed {
			continue
		}
		stored = append(stored, baseline.DecisionValue{ID: id, Value: decision.Value})
	}
	normalized, missing, err := baseline.ResolveDecisionInput(profile, stored, catalog)
	if err != nil || len(missing) != 0 {
		state.incompatible = "the existing Setup Manifest contains incompatible project decisions"
		delete(state.currentDecisions, "auth.provider")
		delete(state.currentDecisions, "http.contract")
		return state, nil
	}
	state.currentDecisions = make(map[string]any, len(normalized))
	for _, decision := range normalized {
		state.currentDecisions[decision.ID] = decision.Value
	}
	state.mode = "update"
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
	hasInstructions bool,
) (baseline.PreservationMode, error) {
	defaultIndex := 0
	if hasInstructions {
		defaultIndex = 1
	}
	selected, err := prompt.selectOneDefault(ctx, "Instruction preservation", []string{
		"Greenfield: back up root instruction carriers without importing their rules",
		"Preservation: back up root instruction carriers and review every proposed classification",
	}, defaultIndex)
	if err != nil {
		return "", err
	}
	if selected == 0 {
		return baseline.PreservationModeGreenfield, nil
	}
	return baseline.PreservationModePreservation, nil
}

func baselineHasRootInstructions(snapshot baseline.RepositorySnapshot) bool {
	for _, carrier := range snapshot.Carriers {
		if carrier.Scope == "root" {
			return true
		}
	}
	return false
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
		reuseLabel := "Keep current profile " + state.currentProfile.ID
		if state.mode != "update" {
			reuseLabel = "Reuse existing profile " + state.currentProfile.ID
		}
		selected, err := prompt.selectOneDefault(ctx, "Baseline Profile", []string{
			reuseLabel,
			"Change Baseline Profile",
		}, 0)
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

func promptBaselineProfileAlignment(
	ctx context.Context,
	prompt *baselineHumanPrompt,
	review io.Writer,
	root string,
	catalog *baseline.Catalog,
	state baselineHumanState,
	profile baseline.ResolvedProfile,
	decisions []baseline.DecisionValue,
) (baseline.ResolvedProfile, []baseline.DecisionValue, *baseline.ProfileDraftInput, error) {
	sourceProfileID := profile.ID
	var draft *baseline.ProfileDraftInput
	for {
		alignment, err := baseline.ResolveProfileAlignment(
			ctx,
			root,
			baseline.ProfileAlignmentRequest{
				ProfileID:            profile.ID,
				Decisions:            decisionsSelectedByProfile(decisions, profile),
				Profile:              profilePointer(profile),
				RemediationProfileID: sourceProfileID,
			},
			catalog,
		)
		if err != nil {
			return baseline.ResolvedProfile{}, nil, nil, err
		}
		renderBaselineProfileAlignment(review, alignment)
		if alignment.Ready {
			return profile, decisions, draft, nil
		}

		profileCapabilities := make(map[string]struct{}, len(profile.Capabilities))
		for _, capabilityID := range profile.Capabilities {
			profileCapabilities[capabilityID] = struct{}{}
		}
		var adaptable []baseline.ProfileDivergence
		var nonRemovable []baseline.ProfileDivergence
		for _, divergence := range alignment.Divergences {
			if !divergence.Blocking {
				continue
			}
			if _, profileSpecific := profileCapabilities[divergence.ID]; profileSpecific {
				adaptable = append(adaptable, divergence)
			} else {
				nonRemovable = append(nonRemovable, divergence)
			}
		}
		if len(nonRemovable) != 0 {
			nextActions := uniqueProfileDivergenceActions(nonRemovable)
			return baseline.ResolvedProfile{}, nil, nil, &baselineHumanActionError{
				result: baselineHumanActionResult(
					"decision",
					"Profile alignment has non-removable required divergences; instruction classification did not start",
					strings.Join(nextActions, " "),
				),
			}
		}
		if len(adaptable) == 0 {
			return baseline.ResolvedProfile{}, nil, nil, errors.New(
				"Profile alignment is not ready but has no actionable blocking divergence",
			)
		}

		selected, err := prompt.selectOne(ctx, "Profile divergence resolution", []string{
			"Change Baseline Profile",
			"Create a reviewed repository-owned Profile adaptation",
			"Decline without writing",
		})
		if err != nil {
			return baseline.ResolvedProfile{}, nil, nil, err
		}
		switch selected {
		case 0:
			changed, selectErr := promptBaselineProfile(
				ctx,
				prompt,
				review,
				root,
				catalog,
				baselineHumanState{},
			)
			if selectErr != nil {
				return baseline.ResolvedProfile{}, nil, nil, selectErr
			}
			changedDecisions, decisionErr := promptBaselineDecisions(
				ctx,
				prompt,
				catalog,
				changed,
				state,
			)
			if decisionErr != nil {
				return baseline.ResolvedProfile{}, nil, nil, decisionErr
			}
			profile = changed
			decisions = changedDecisions
			sourceProfileID = changed.ID
			draft = nil
		case 1:
			adapted, adaptedDecisions, input, adaptationErr := promptBaselineProfileAdaptation(
				ctx,
				prompt,
				review,
				root,
				catalog,
				sourceProfileID,
				profile,
				decisions,
				adaptable,
			)
			if adaptationErr != nil {
				return baseline.ResolvedProfile{}, nil, nil, adaptationErr
			}
			profile = adapted
			decisions = adaptedDecisions
			draft = &input
		default:
			return baseline.ResolvedProfile{}, nil, nil, &baselineHumanActionError{
				result: baselineHumanActionResult(
					"decision",
					"Profile divergence resolution was declined; instruction classification did not start and no repository bytes were written",
					"rerun roundfix baseline to change the Profile or review a repository-owned adaptation",
				),
			}
		}
	}
}

func promptBaselineProfileAdaptation(
	ctx context.Context,
	prompt *baselineHumanPrompt,
	review io.Writer,
	root string,
	catalog *baseline.Catalog,
	sourceProfileID string,
	source baseline.ResolvedProfile,
	decisions []baseline.DecisionValue,
	blocking []baseline.ProfileDivergence,
) (baseline.ResolvedProfile, []baseline.DecisionValue, baseline.ProfileDraftInput, error) {
	removedModules := suggestedProfileModuleRemovals(source, blocking)
	removedCapabilities := make([]string, 0, len(blocking))
	for _, divergence := range blocking {
		removedCapabilities = append(removedCapabilities, divergence.ID)
	}
	for {
		fmt.Fprintln(review, "\nRepository-owned Profile adaptation proposal")
		printProfileRemovalReview(review, "Modules", removedModules)
		printProfileRemovalReview(review, "Capabilities", removedCapabilities)
		selected, err := prompt.selectOne(ctx, "Profile adaptation review", []string{
			"Accept every listed removal and re-audit",
			"Review every module and capability selection",
			"Decline without writing",
		})
		if err != nil {
			return baseline.ResolvedProfile{}, nil, baseline.ProfileDraftInput{}, err
		}
		if selected == 2 {
			return baseline.ResolvedProfile{}, nil, baseline.ProfileDraftInput{},
				&baselineHumanActionError{result: baselineHumanActionResult(
					"decision",
					"Profile adaptation was declined; no repository bytes were written",
					"rerun roundfix baseline to review a different Profile adaptation",
				)}
		}
		if selected == 1 {
			removedModules, err = promptProfileRemovals(
				ctx,
				prompt,
				"module",
				source.Modules,
				removedModules,
			)
			if err != nil {
				return baseline.ResolvedProfile{}, nil, baseline.ProfileDraftInput{}, err
			}
			removedCapabilities, err = promptProfileRemovals(
				ctx,
				prompt,
				"capability",
				source.Capabilities,
				removedCapabilities,
			)
			if err != nil {
				return baseline.ResolvedProfile{}, nil, baseline.ProfileDraftInput{}, err
			}
		}

		id, err := prompt.readNonEmpty(ctx, "Repository-owned Baseline Profile ID")
		if err != nil {
			return baseline.ResolvedProfile{}, nil, baseline.ProfileDraftInput{}, err
		}
		input, err := baseline.NewProfileAdaptationDraft(
			sourceProfileID,
			id,
			removedModules,
			removedCapabilities,
			catalog,
		)
		if err != nil {
			fmt.Fprintf(review, "Profile adaptation is invalid: %v\n", err)
			fmt.Fprintln(review, "Review the removals and provide a valid repository-owned Profile ID.")
			continue
		}
		resolved, canonical, err := baseline.ResolveProfileDraft(root, input, catalog)
		if err != nil {
			fmt.Fprintf(review, "Profile adaptation is invalid: %v\n", err)
			fmt.Fprintln(review, "Review the removals and provide a valid repository-owned Profile ID.")
			continue
		}
		input.Document = canonical
		return resolved, decisionsSelectedByProfile(decisions, resolved), input, nil
	}
}

func suggestedProfileModuleRemovals(
	profile baseline.ResolvedProfile,
	blocking []baseline.ProfileDivergence,
) []string {
	suggestions := map[string]string{
		"capability.stack.bun":          "bun",
		"capability.stack.turborepo":    "monorepo",
		"capability.stack.typescript":   "typescript",
		"capability.workspace.backend":  "backend",
		"capability.workspace.frontend": "frontend",
	}
	selected := make(map[string]struct{}, len(profile.Modules))
	for _, moduleID := range profile.Modules {
		selected[moduleID] = struct{}{}
	}
	removed := make(map[string]struct{})
	for _, divergence := range blocking {
		moduleID := suggestions[divergence.ID]
		if _, exists := selected[moduleID]; exists {
			removed[moduleID] = struct{}{}
		}
	}
	if _, selectedAutonomousWork := selected["autonomous-work"]; selectedAutonomousWork {
		removed["autonomous-work"] = struct{}{}
	}
	result := make([]string, 0, len(removed))
	for _, moduleID := range profile.Modules {
		if _, remove := removed[moduleID]; remove {
			result = append(result, moduleID)
		}
	}
	return result
}

func promptProfileRemovals(
	ctx context.Context,
	prompt *baselineHumanPrompt,
	kind string,
	values []string,
	current []string,
) ([]string, error) {
	removed := make(map[string]struct{}, len(current))
	for _, value := range current {
		removed[value] = struct{}{}
	}
	result := make([]string, 0)
	for _, value := range values {
		defaultIndex := 0
		if _, remove := removed[value]; remove {
			defaultIndex = 1
		}
		selected, err := prompt.selectOneDefault(
			ctx,
			fmt.Sprintf("Profile %s %s", kind, value),
			[]string{"Keep", "Remove"},
			defaultIndex,
		)
		if err != nil {
			return nil, err
		}
		if selected == 1 {
			result = append(result, value)
		}
	}
	return result, nil
}

func decisionsSelectedByProfile(
	decisions []baseline.DecisionValue,
	profile baseline.ResolvedProfile,
) []baseline.DecisionValue {
	selected := make(map[string]struct{}, len(profile.Decisions))
	for _, decisionID := range profile.Decisions {
		selected[decisionID] = struct{}{}
	}
	result := make([]baseline.DecisionValue, 0, len(decisions))
	for _, decision := range decisions {
		if _, retained := selected[decision.ID]; retained {
			result = append(result, decision)
		}
	}
	return result
}

func renderBaselineProfileAlignment(output io.Writer, alignment baseline.ProfileAlignment) {
	fmt.Fprintf(output, "\nBaseline Profile alignment: %s\n", alignment.State)
	fmt.Fprintf(output, "Profile: %s\n", alignment.Profile.ID)
	if len(alignment.Divergences) == 0 {
		fmt.Fprintln(output, "Divergences: none")
		return
	}
	fmt.Fprintln(output, "Divergences:")
	for _, divergence := range alignment.Divergences {
		severity := "advisory"
		if divergence.Blocking {
			severity = "blocking"
		}
		fmt.Fprintf(
			output,
			"- %s %s (%s): %s\n",
			severity,
			divergence.ID,
			divergence.Code,
			divergence.Message,
		)
		if divergence.NextAction != "" {
			fmt.Fprintf(output, "  Next action: %s\n", divergence.NextAction)
		}
	}
}

func printProfileRemovalReview(output io.Writer, label string, removed []string) {
	fmt.Fprintf(output, "%s removed (%d):\n", label, len(removed))
	if len(removed) == 0 {
		fmt.Fprintln(output, "- none")
		return
	}
	for _, id := range removed {
		fmt.Fprintf(output, "- %s\n", id)
	}
}

func uniqueProfileDivergenceActions(divergences []baseline.ProfileDivergence) []string {
	seen := make(map[string]struct{})
	actions := make([]string, 0, len(divergences))
	for _, divergence := range divergences {
		action := strings.TrimSpace(divergence.NextAction)
		if action == "" {
			continue
		}
		if _, duplicate := seen[action]; duplicate {
			continue
		}
		seen[action] = struct{}{}
		actions = append(actions, action)
	}
	if len(actions) == 0 {
		return []string{"resolve every non-removable required divergence and rerun roundfix baseline"}
	}
	return actions
}

func profilePointer(profile baseline.ResolvedProfile) *baseline.ResolvedProfile {
	return &profile
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
		if baseline.ValidateDecisionValue(catalog, id, value) == nil {
			encoded, _ := json.Marshal(value)
			selected, err := prompt.selectOneDefault(ctx, declaration.Summary, []string{
				fmt.Sprintf("Keep %s=%s", id, encoded),
				"Change " + id,
			}, 0)
			if err != nil {
				return nil, err
			}
			if selected == 0 {
				return value, nil
			}
		}
	}

	switch declaration.Type {
	case "identifier-strategy":
		suggestion, ok := declaration.Suggestion.(map[string]any)
		if !ok {
			return nil, fmt.Errorf("prompt Baseline decision %q: suggestion is invalid", id)
		}
		encoded, err := json.Marshal(suggestion)
		if err != nil {
			return nil, fmt.Errorf("prompt Baseline decision %q: encode suggestion: %w", id, err)
		}
		selected, err := prompt.selectOneDefault(
			ctx,
			declaration.Summary,
			[]string{
				fmt.Sprintf("Keep suggested UUID version 7: %s=%s", id, encoded),
				"Use a repository-defined identifier strategy",
			},
			0,
		)
		if err != nil {
			return nil, err
		}
		if selected == 0 {
			return suggestion, nil
		}
		guidance, err := prompt.readNonEmpty(
			ctx,
			"Operative rule for new project-owned Internal Identifiers",
		)
		if err != nil {
			return nil, err
		}
		return map[string]any{
			"kind":     "repository-defined",
			"guidance": guidance,
		}, nil
	case "auth-provider":
		suggestion, err := baselineAuthProviderSuggestion(declaration, current, catalog)
		if err != nil {
			return nil, fmt.Errorf("prompt Baseline decision %q: %w", id, err)
		}
		exception, ok := suggestion["routeException"].(map[string]any)
		if !ok {
			return nil, fmt.Errorf("prompt Baseline decision %q: route suggestion is invalid", id)
		}
		encoded, err := json.Marshal(suggestion)
		if err != nil {
			return nil, fmt.Errorf("prompt Baseline decision %q: encode suggestion: %w", id, err)
		}
		selected, err := prompt.selectOneDefault(
			ctx,
			declaration.Summary,
			[]string{
				fmt.Sprintf("Keep suggested Better Auth provider: %s=%s", id, encoded),
				"Change the Better Auth route exception",
			},
			0,
		)
		if err != nil {
			return nil, err
		}
		if selected == 0 {
			return suggestion, nil
		}
		scope, _ := exception["scope"].(string)
		scope, err = prompt.readNonEmptyDefault(ctx, "Better Auth route scope", scope)
		if err != nil {
			return nil, err
		}
		methodSelection, err := prompt.selectOneDefault(
			ctx,
			"Better Auth route methods",
			[]string{"GET and POST", "GET only", "POST only"},
			0,
		)
		if err != nil {
			return nil, err
		}
		methods := []any{"GET", "POST"}
		if methodSelection == 1 {
			methods = []any{"GET"}
		}
		if methodSelection == 2 {
			methods = []any{"POST"}
		}
		reason, _ := exception["reason"].(string)
		reason, err = prompt.readNonEmptyDefault(ctx, "Better Auth provider-protocol reason", reason)
		if err != nil {
			return nil, err
		}
		return map[string]any{
			"kind": "better-auth",
			"routeException": map[string]any{
				"scope": scope, "methods": methods,
				"owner": "Better Auth", "reason": reason,
			},
		}, nil
	case "boolean":
		defaultIndex := -1
		if value, ok := declaration.Default.(bool); ok {
			defaultIndex = 1
			if value {
				defaultIndex = 0
			}
		}
		selected, err := prompt.selectOneDefault(
			ctx,
			declaration.Summary,
			[]string{"Yes", "No"},
			defaultIndex,
		)
		if err != nil {
			return nil, err
		}
		return selected == 0, nil
	case "enum":
		selected, err := prompt.selectOneDefault(
			ctx,
			declaration.Summary,
			declaration.Values,
			baselineDecisionDefaultIndex(declaration.Values, declaration.Default),
		)
		if err != nil {
			return nil, err
		}
		return declaration.Values[selected], nil
	case "http-contract":
		defaultMode := ""
		if value, ok := declaration.Default.(map[string]any); ok {
			defaultMode, _ = value["mode"].(string)
		}
		selected, err := prompt.selectOneDefault(
			ctx,
			declaration.Summary,
			declaration.Modes,
			baselineDecisionDefaultIndex(declaration.Modes, defaultMode),
		)
		if err != nil {
			return nil, err
		}
		return map[string]any{"mode": declaration.Modes[selected]}, nil
	case "string":
		if value, ok := declaration.Default.(string); ok && strings.TrimSpace(value) != "" {
			return prompt.readNonEmptyDefault(ctx, declaration.Summary+" ("+id+")", value)
		}
		return prompt.readNonEmpty(ctx, declaration.Summary+" ("+id+")")
	default:
		return nil, fmt.Errorf("prompt Baseline decision %q: unsupported type %q", id, declaration.Type)
	}
}

func baselineAuthProviderSuggestion(
	declaration baselineDecisionDeclaration,
	current map[string]any,
	catalog *baseline.Catalog,
) (map[string]any, error) {
	raw, ok := declaration.Suggestion.(map[string]any)
	if !ok {
		return nil, errors.New("suggestion is invalid")
	}
	suggestion, err := cloneBaselineDecisionObject(raw)
	if err != nil {
		return nil, fmt.Errorf("clone suggestion: %w", err)
	}

	httpValue, ok := current["http.contract"]
	if !ok {
		return suggestion, nil
	}
	contract, ok := httpValue.(map[string]any)
	if !ok {
		return suggestion, nil
	}
	exceptions, ok := contract["exceptions"].([]any)
	if !ok {
		return suggestion, nil
	}
	profile := baseline.ResolvedProfile{
		ID:        "human-auth-provider-suggestion",
		Decisions: []string{"auth.provider", "http.contract"},
	}
	for _, rawException := range exceptions {
		exception, ok := rawException.(map[string]any)
		if !ok {
			continue
		}
		reason, ok := exception["reason"].(string)
		if !ok {
			continue
		}
		candidate, err := cloneBaselineDecisionObject(suggestion)
		if err != nil {
			return nil, fmt.Errorf("clone suggestion candidate: %w", err)
		}
		routeException, ok := candidate["routeException"].(map[string]any)
		if !ok {
			return nil, errors.New("route suggestion is invalid")
		}
		routeException["reason"] = reason
		_, missing, err := baseline.ResolveDecisionInput(
			profile,
			[]baseline.DecisionValue{
				{ID: "auth.provider", Value: candidate},
				{ID: "http.contract", Value: httpValue},
			},
			catalog,
		)
		if err == nil && len(missing) == 0 {
			return candidate, nil
		}
	}
	return suggestion, nil
}

func cloneBaselineDecisionObject(value map[string]any) (map[string]any, error) {
	encoded, err := json.Marshal(value)
	if err != nil {
		return nil, err
	}
	var cloned map[string]any
	if err := json.Unmarshal(encoded, &cloned); err != nil {
		return nil, err
	}
	return cloned, nil
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

func baselineDecisionDefaultIndex(options []string, value any) int {
	text, ok := value.(string)
	if !ok {
		return -1
	}
	for index, option := range options {
		if option == text {
			return index
		}
	}
	return -1
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
				Path:          "docs/agents/specific-repository.md",
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
	return prompt.selectOneDefault(ctx, label, options, -1)
}

func (prompt *baselineHumanPrompt) selectOneDefault(
	ctx context.Context,
	label string,
	options []string,
	defaultIndex int,
) (int, error) {
	if len(options) == 0 {
		return 0, fmt.Errorf("prompt %q has no options", label)
	}
	if defaultIndex >= len(options) {
		return 0, fmt.Errorf("prompt %q default index %d is outside its options", label, defaultIndex)
	}
	prompt.step++
	fmt.Fprintf(prompt.writer, "\nPrompt %d: %s\n", prompt.step, label)
	for index, option := range options {
		suffix := ""
		if index == defaultIndex {
			suffix = " (default)"
		}
		fmt.Fprintf(prompt.writer, "  %d. %s%s\n", index+1, option, suffix)
	}
	for {
		line, err := prompt.readLine(ctx)
		if err != nil {
			return 0, err
		}
		answer := strings.TrimSpace(line)
		if answer == "" && defaultIndex >= 0 {
			return defaultIndex, nil
		}
		choice, err := strconv.Atoi(answer)
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

func (prompt *baselineHumanPrompt) readNonEmptyDefault(
	ctx context.Context,
	label string,
	defaultValue string,
) (string, error) {
	if strings.TrimSpace(defaultValue) == "" {
		return "", fmt.Errorf("prompt %q has an empty default value", label)
	}
	prompt.step++
	fmt.Fprintf(prompt.writer, "\nPrompt %d: %s\n> [default: %s] ", prompt.step, label, defaultValue)
	line, err := prompt.readLine(ctx)
	if err != nil {
		return "", err
	}
	if value := strings.TrimSpace(line); value != "" {
		return value, nil
	}
	return defaultValue, nil
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
