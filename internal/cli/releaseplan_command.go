package cli

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"strings"

	"roundfix/internal/app"
	"roundfix/internal/preflight"
	"roundfix/internal/releaseplan"
)

const (
	releasePlanFormatText = "text"
	releasePlanFormatJSON = "json"
)

type releasePlanCommandRequest struct {
	from         string
	to           string
	resetTo      string
	impact       releaseplan.Impact
	reason       string
	outputFormat string
}

var (
	releasePlanCommandGitRunner preflight.GitRunner = preflight.ExecGitRunner{}
	releasePlanCommandGHRunner  preflight.GHRunner  = preflight.ExecGHRunner{}
)

func runReleaseCommand(ctx context.Context, args []string, stdout, stderr io.Writer, environment commandEnvironment) int {
	if len(args) == 0 {
		fmt.Fprint(stdout, commandUsage("release"))
		return exitOK
	}
	switch args[0] {
	case "-h", "--help", "help":
		fmt.Fprint(stdout, commandUsage("release"))
		return exitOK
	case "plan":
		return runReleasePlanCommand(ctx, args[1:], stdout, stderr, environment)
	default:
		fmt.Fprintf(stderr, "%s: unknown release command %q; run '%s release --help' for usage\n", app.Name, args[0], app.Name)
		return exitPreflight
	}
}

func runReleasePlanCommand(ctx context.Context, args []string, stdout, stderr io.Writer, environment commandEnvironment) int {
	if commandWantsHelp(args) {
		fmt.Fprint(stdout, commandUsage("release plan"))
		return exitOK
	}

	req, err := parseReleasePlanCommand(args)
	if err != nil {
		printReleasePlanFailure(err, stderr)
		return exitPreflight
	}

	if req.resetTo != "" {
		return runReleaseResetPlanCommand(ctx, req, stdout, stderr, environment)
	}

	workDir, err := environment.resolveWorkDir("resolve release plan working directory")
	if err != nil {
		printReleasePlanFailure(err, stderr)
		return exitPreflight
	}
	source := newReleasePlanGitSource(workDir, commandDependenciesForContext(ctx).releasePlanGitRunner)
	plan, err := releaseplan.Build(ctx, releaseplan.Request{
		From:         req.from,
		To:           req.to,
		ManualImpact: req.impact,
		ManualReason: req.reason,
	}, source)
	if err != nil {
		printReleasePlanFailure(err, stderr)
		return exitPreflight
	}

	switch req.outputFormat {
	case releasePlanFormatText:
		printReleasePlanText(plan, stdout)
	case releasePlanFormatJSON:
		if err := printReleasePlanJSON(plan, stdout); err != nil {
			printReleasePlanFailure(err, stderr)
			return exitRunFailed
		}
	default:
		printReleasePlanFailure(validationError{message: fmt.Sprintf("unsupported --format %q; use text or json", req.outputFormat)}, stderr)
		return exitPreflight
	}
	return releasePlanExitCode(plan.State)
}

func runReleaseResetPlanCommand(ctx context.Context, req releasePlanCommandRequest, stdout, stderr io.Writer, environment commandEnvironment) int {
	workDir, err := environment.resolveWorkDir("resolve release plan working directory")
	if err != nil {
		printReleasePlanFailure(err, stderr)
		return exitPreflight
	}
	source := newReleasePlanResetSource(workDir, commandDependenciesForContext(ctx).releasePlanGitRunner, commandDependenciesForContext(ctx).releasePlanGHRunner)
	target, err := source.ResolveTarget(ctx)
	if err != nil {
		printReleasePlanFailure(err, stderr)
		return exitPreflight
	}
	plan, err := releaseplan.BuildReset(ctx, releaseplan.ResetRequest{
		TargetVersion: req.resetTo,
		Target:        target,
	}, source)
	if err != nil {
		printReleasePlanFailure(err, stderr)
		return exitPreflight
	}

	switch req.outputFormat {
	case releasePlanFormatText:
		printReleaseResetPlanText(plan, stdout)
	case releasePlanFormatJSON:
		if err := printReleaseResetPlanJSON(plan, stdout); err != nil {
			printReleasePlanFailure(err, stderr)
			return exitRunFailed
		}
	default:
		printReleasePlanFailure(validationError{message: fmt.Sprintf("unsupported --format %q; use text or json", req.outputFormat)}, stderr)
		return exitPreflight
	}
	return releasePlanExitCode(plan.State)
}

func parseReleasePlanCommand(args []string) (releasePlanCommandRequest, error) {
	req := releasePlanCommandRequest{outputFormat: releasePlanFormatText}
	fs := flag.NewFlagSet("release plan", flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	var impact string
	fs.StringVar(&req.from, "from", "", "Stable release tag to use as the base")
	fs.StringVar(&req.to, "to", "", "Target revision to analyze")
	fs.StringVar(&req.resetTo, "reset-to", "", "Stable version for a read-only release history reset plan")
	fs.StringVar(&impact, "impact", "", "Manual impact for ambiguous changes")
	fs.StringVar(&req.reason, "reason", "", "Manual classification reason")
	fs.StringVar(&req.outputFormat, "format", releasePlanFormatText, "Output format: text or json")
	if err := fs.Parse(args); err != nil {
		return releasePlanCommandRequest{}, validationError{message: fmt.Sprintf("invalid release plan arguments: %v; run '%s release plan --help' for usage", err, app.Name)}
	}
	if remaining := fs.Args(); len(remaining) > 0 {
		return releasePlanCommandRequest{}, validationError{message: fmt.Sprintf("unexpected argument %q; run '%s release plan --help' for usage", remaining[0], app.Name)}
	}
	req.from = strings.TrimSpace(req.from)
	req.to = strings.TrimSpace(req.to)
	req.resetTo = strings.TrimSpace(req.resetTo)
	req.outputFormat = strings.TrimSpace(req.outputFormat)
	req.reason = strings.TrimSpace(req.reason)
	if req.outputFormat == "" {
		req.outputFormat = releasePlanFormatText
	}
	if req.outputFormat != releasePlanFormatText && req.outputFormat != releasePlanFormatJSON {
		return releasePlanCommandRequest{}, validationError{message: fmt.Sprintf("unsupported --format %q; use text or json", req.outputFormat)}
	}
	if trimmedImpact := strings.TrimSpace(impact); trimmedImpact != "" {
		req.impact = releaseplan.Impact(trimmedImpact)
	}
	visited := map[string]bool{}
	fs.Visit(func(flag *flag.Flag) {
		visited[flag.Name] = true
	})
	if visited["reset-to"] {
		for _, incompatible := range []string{"from", "to", "impact", "reason"} {
			if visited[incompatible] {
				return releasePlanCommandRequest{}, validationError{
					message: fmt.Sprintf("--reset-to cannot be combined with --%s; choose range planning or reset planning", incompatible),
				}
			}
		}
		if _, err := releaseplan.ParseStableVersion(req.resetTo); err != nil {
			return releasePlanCommandRequest{}, err
		}
	}
	return req, nil
}

func releasePlanExitCode(state releaseplan.State) int {
	switch state {
	case releaseplan.StateReady, releaseplan.StateNoRelease:
		return exitOK
	case releaseplan.StateApprovalRequired, releaseplan.StateManualClassificationRequired:
		return exitUnverified
	default:
		return exitPreflight
	}
}

func printReleasePlanFailure(err error, stderr io.Writer) {
	fmt.Fprintf(stderr, "%s: release plan failed: %v\n", app.Name, err)
}

func printReleasePlanText(plan releaseplan.Plan, stdout io.Writer) {
	fmt.Fprintf(stdout, "Decision: %s\n", plan.State)
	fmt.Fprintf(stdout, "Base: %s (%s)\n", plan.Base.Tag, shortReleasePlanSHA(plan.Base.CommitSHA))
	fmt.Fprintf(stdout, "Target: %s (%s)\n", plan.Target.Name, shortReleasePlanSHA(plan.Target.CommitSHA))
	fmt.Fprintf(stdout, "Impact: %s\n", plan.Classification.Impact)
	fmt.Fprintf(stdout, "Breaking: %t\n", plan.Classification.Breaking)
	if plan.ProposedVersion == "" {
		fmt.Fprintln(stdout, "Proposed version: none")
	} else {
		fmt.Fprintf(stdout, "Proposed version: %s\n", plan.ProposedVersion)
	}
	if plan.Classification.ManualReason != "" {
		fmt.Fprintf(stdout, "Manual reason: %s\n", plan.Classification.ManualReason)
	}

	switch plan.State {
	case releaseplan.StateReady:
		fmt.Fprintln(stdout, "Approval required: no")
		fmt.Fprintf(stdout, "Next action: release may proceed for %s after independent release verification.\n", plan.ProposedVersion)
	case releaseplan.StateNoRelease:
		fmt.Fprintln(stdout, "Approval required: no")
		fmt.Fprintln(stdout, "Next action: no release is required for the committed range.")
	case releaseplan.StateApprovalRequired:
		fmt.Fprintln(stdout, "Approval required: yes")
		fmt.Fprintf(stdout, "Approval question: %s\n", plan.Approval.Question)
		fmt.Fprintln(stdout, "Next action: answer the approval question before any release mutation.")
	case releaseplan.StateManualClassificationRequired:
		fmt.Fprintln(stdout, "Approval required: no")
		fmt.Fprintf(stdout, "Next action: rerun roundfix release plan --from %s --to %s --impact <none|patch|minor|major> --reason <text>\n", plan.Base.Tag, plan.Target.Name)
	}

	changes := releasePlanTextChanges(plan)
	if len(changes) == 0 {
		return
	}
	if plan.State == releaseplan.StateManualClassificationRequired {
		fmt.Fprintln(stdout, "Blocking commits:")
	} else {
		fmt.Fprintln(stdout, "Determining commits:")
	}
	for _, change := range changes {
		fmt.Fprintf(stdout, "- %s %s [%s]\n", shortReleasePlanSHA(change.CommitSHA), change.Subject, releasePlanChangeLabel(change))
	}
}

func releasePlanTextChanges(plan releaseplan.Plan) []releaseplan.ChangeEvidence {
	if plan.State == releaseplan.StateManualClassificationRequired {
		return releasePlanBlockingChanges(plan)
	}
	if plan.State == releaseplan.StateNoRelease && plan.Classification.ManualReason == "" {
		return append([]releaseplan.ChangeEvidence(nil), plan.Changes...)
	}

	changes := make([]releaseplan.ChangeEvidence, 0, len(plan.Changes))
	seen := map[string]bool{}
	for _, change := range plan.Changes {
		include := false
		switch {
		case plan.Classification.ManualReason != "" && change.AutomaticImpact == releaseplan.ImpactNone && change.CrossesMaintenanceOnlyBoundary:
			include = true
		case plan.Classification.Breaking && change.Breaking:
			include = true
		case change.AutomaticImpact != releaseplan.ImpactNone && change.AutomaticImpact == plan.Classification.Impact:
			include = true
		}
		if include && !seen[change.CommitSHA] {
			seen[change.CommitSHA] = true
			changes = append(changes, change)
		}
	}
	return changes
}

func releasePlanBlockingChanges(plan releaseplan.Plan) []releaseplan.ChangeEvidence {
	blocking := map[string]bool{}
	for _, sha := range plan.Classification.BlockingCommits {
		blocking[sha] = true
	}
	changes := make([]releaseplan.ChangeEvidence, 0, len(blocking))
	for _, change := range plan.Changes {
		if blocking[change.CommitSHA] {
			changes = append(changes, change)
		}
	}
	return changes
}

func releasePlanChangeLabel(change releaseplan.ChangeEvidence) string {
	labels := []string{fmt.Sprintf("impact=%s", change.AutomaticImpact)}
	if change.ConventionalType != "" {
		labels = append(labels, "type="+change.ConventionalType)
	}
	if change.Breaking {
		labels = append(labels, "breaking=true")
	}
	if change.AutomaticImpact == releaseplan.ImpactNone && change.CrossesMaintenanceOnlyBoundary {
		labels = append(labels, "manual")
	}
	return strings.Join(labels, ", ")
}

func printReleasePlanJSON(plan releaseplan.Plan, stdout io.Writer) error {
	encoder := json.NewEncoder(stdout)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(releasePlanJSONFromPlan(plan)); err != nil {
		return fmt.Errorf("encode release plan JSON: %w", err)
	}
	return nil
}

func printReleaseResetPlanText(plan releaseplan.ResetPlan, stdout io.Writer) {
	fmt.Fprintf(stdout, "Decision: %s\n", plan.State)
	fmt.Fprintf(stdout, "Reset target: %s\n", plan.TargetVersion)
	fmt.Fprintf(stdout, "Target: %s (%s)\n", plan.TargetRevision, plan.TargetCommit)
	fmt.Fprintf(stdout, "Plan digest: %s\n", plan.PlanDigest)
	fmt.Fprintln(stdout, "Approval required: yes")
	fmt.Fprintf(stdout, "Approval question: %s\n", plan.Approval.Question)
	fmt.Fprintln(stdout, "Stable tags:")
	for _, tag := range plan.Tags {
		location := string(tag.Source)
		if tag.Remote != "" {
			location += ":" + tag.Remote
		}
		fmt.Fprintf(
			stdout,
			"- %s %s identity=%s target=%s location=%s\n",
			tag.Name,
			tag.Ref,
			tag.ImmutableID,
			tag.TargetCommit,
			location,
		)
	}
	fmt.Fprintln(stdout, "GitHub Releases:")
	for _, release := range plan.Releases {
		target := release.TargetCommit
		if target == "" {
			target = "unavailable"
		}
		fmt.Fprintf(
			stdout,
			"- %s name=%q identity=%s node=%s target=%s targetCommitish=%s\n",
			release.TagName,
			release.Name,
			release.ImmutableID,
			release.NodeID,
			target,
			release.TargetCommitish,
		)
	}
	fmt.Fprintln(stdout, "Next action: review this complete inventory; any tag or GitHub Release deletion requires separate explicit post-QA authority.")
}

func printReleaseResetPlanJSON(plan releaseplan.ResetPlan, stdout io.Writer) error {
	encoder := json.NewEncoder(stdout)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(releaseResetPlanJSONFromPlan(plan)); err != nil {
		return fmt.Errorf("encode release reset plan JSON: %w", err)
	}
	return nil
}

type releasePlanJSON struct {
	SchemaVersion   string                        `json:"schemaVersion"`
	State           releaseplan.State             `json:"state"`
	Base            releasePlanVersionRefJSON     `json:"base"`
	Target          releasePlanRevisionRefJSON    `json:"target"`
	Classification  releasePlanClassificationJSON `json:"classification"`
	ProposedVersion string                        `json:"proposedVersion,omitempty"`
	Approval        releasePlanApprovalJSON       `json:"approval"`
	Changes         []releasePlanChangeJSON       `json:"changes"`
}

type releasePlanVersionRefJSON struct {
	Tag       string `json:"tag"`
	Version   string `json:"version"`
	CommitSHA string `json:"commitSHA"`
}

type releasePlanRevisionRefJSON struct {
	Name      string `json:"name"`
	CommitSHA string `json:"commitSHA"`
}

type releasePlanClassificationJSON struct {
	Source                       releaseplan.ClassificationSource `json:"source"`
	Impact                       releaseplan.Impact               `json:"impact"`
	Breaking                     bool                             `json:"breaking"`
	ManualReason                 string                           `json:"manualReason,omitempty"`
	ManualClassificationRequired bool                             `json:"manualClassificationRequired"`
	BlockingCommits              []string                         `json:"blockingCommits,omitempty"`
}

type releasePlanApprovalJSON struct {
	Required        bool                      `json:"required"`
	Increment       releaseplan.IncrementKind `json:"increment"`
	ProposedVersion string                    `json:"proposedVersion,omitempty"`
	Question        string                    `json:"question,omitempty"`
}

type releasePlanChangeJSON struct {
	CommitSHA                      string             `json:"commitSHA"`
	Subject                        string             `json:"subject"`
	ConventionalType               string             `json:"conventionalType,omitempty"`
	Breaking                       bool               `json:"breaking"`
	AutomaticImpact                releaseplan.Impact `json:"automaticImpact"`
	CrossesMaintenanceOnlyBoundary bool               `json:"crossesMaintenanceOnlyBoundary"`
}

type releaseResetPlanJSON struct {
	SchemaVersion string                     `json:"schemaVersion"`
	State         releaseplan.State          `json:"state"`
	TargetVersion string                     `json:"targetVersion"`
	Target        releasePlanRevisionRefJSON `json:"target"`
	PlanDigest    string                     `json:"planDigest"`
	Approval      releasePlanApprovalJSON    `json:"approval"`
	Tags          []releaseResetTagJSON      `json:"tags"`
	Releases      []releaseResetReleaseJSON  `json:"releases"`
}

type releaseResetTagJSON struct {
	Name         string                `json:"name"`
	Source       releaseplan.TagSource `json:"source"`
	Remote       string                `json:"remote,omitempty"`
	Ref          string                `json:"ref"`
	ImmutableID  string                `json:"immutableID"`
	TargetCommit string                `json:"targetCommit"`
}

type releaseResetReleaseJSON struct {
	ID              int64  `json:"id"`
	NodeID          string `json:"nodeID"`
	Name            string `json:"name"`
	TagName         string `json:"tagName"`
	TargetCommitish string `json:"targetCommitish,omitempty"`
	TargetCommit    string `json:"targetCommit,omitempty"`
	ImmutableID     string `json:"immutableID"`
}

func releasePlanJSONFromPlan(plan releaseplan.Plan) releasePlanJSON {
	changes := make([]releasePlanChangeJSON, 0, len(plan.Changes))
	for _, change := range plan.Changes {
		changes = append(changes, releasePlanChangeJSON{
			CommitSHA:                      change.CommitSHA,
			Subject:                        change.Subject,
			ConventionalType:               change.ConventionalType,
			Breaking:                       change.Breaking,
			AutomaticImpact:                change.AutomaticImpact,
			CrossesMaintenanceOnlyBoundary: change.CrossesMaintenanceOnlyBoundary,
		})
	}
	return releasePlanJSON{
		SchemaVersion: releaseplan.SchemaVersion,
		State:         plan.State,
		Base: releasePlanVersionRefJSON{
			Tag:       plan.Base.Tag,
			Version:   plan.Base.Version.String(),
			CommitSHA: plan.Base.CommitSHA,
		},
		Target: releasePlanRevisionRefJSON{
			Name:      plan.Target.Name,
			CommitSHA: plan.Target.CommitSHA,
		},
		Classification: releasePlanClassificationJSON{
			Source:                       plan.Classification.Source,
			Impact:                       plan.Classification.Impact,
			Breaking:                     plan.Classification.Breaking,
			ManualReason:                 plan.Classification.ManualReason,
			ManualClassificationRequired: plan.Classification.ManualClassificationRequired,
			BlockingCommits:              append([]string(nil), plan.Classification.BlockingCommits...),
		},
		ProposedVersion: plan.ProposedVersion,
		Approval: releasePlanApprovalJSON{
			Required:        plan.Approval.Required,
			Increment:       plan.Approval.Increment,
			ProposedVersion: plan.Approval.ProposedVersion,
			Question:        plan.Approval.Question,
		},
		Changes: changes,
	}
}

func releaseResetPlanJSONFromPlan(plan releaseplan.ResetPlan) releaseResetPlanJSON {
	tags := make([]releaseResetTagJSON, 0, len(plan.Tags))
	for _, tag := range plan.Tags {
		tags = append(tags, releaseResetTagJSON{
			Name:         tag.Name,
			Source:       tag.Source,
			Remote:       tag.Remote,
			Ref:          tag.Ref,
			ImmutableID:  tag.ImmutableID,
			TargetCommit: tag.TargetCommit,
		})
	}
	releases := make([]releaseResetReleaseJSON, 0, len(plan.Releases))
	for _, release := range plan.Releases {
		releases = append(releases, releaseResetReleaseJSON{
			ID:              release.ID,
			NodeID:          release.NodeID,
			Name:            release.Name,
			TagName:         release.TagName,
			TargetCommitish: release.TargetCommitish,
			TargetCommit:    release.TargetCommit,
			ImmutableID:     release.ImmutableID,
		})
	}
	return releaseResetPlanJSON{
		SchemaVersion: plan.SchemaVersion,
		State:         plan.State,
		TargetVersion: plan.TargetVersion,
		Target: releasePlanRevisionRefJSON{
			Name:      plan.TargetRevision,
			CommitSHA: plan.TargetCommit,
		},
		PlanDigest: plan.PlanDigest,
		Approval: releasePlanApprovalJSON{
			Required:        plan.Approval.Required,
			Increment:       plan.Approval.Increment,
			ProposedVersion: plan.Approval.ProposedVersion,
			Question:        plan.Approval.Question,
		},
		Tags:     tags,
		Releases: releases,
	}
}

func shortReleasePlanSHA(sha string) string {
	if len(sha) <= 12 {
		return sha
	}
	return sha[:12]
}
