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
	impact       releaseplan.Impact
	reason       string
	outputFormat string
}

func runReleaseCommand(ctx context.Context, args []string, stdout, stderr io.Writer) int {
	if len(args) == 0 {
		fmt.Fprint(stdout, commandUsage("release"))
		return exitOK
	}
	switch args[0] {
	case "-h", "--help", "help":
		fmt.Fprint(stdout, commandUsage("release"))
		return exitOK
	case "plan":
		return runReleasePlanCommand(ctx, args[1:], stdout, stderr)
	default:
		fmt.Fprintf(stderr, "%s: unknown release command %q; run '%s release --help' for usage\n", app.Name, args[0], app.Name)
		return exitPreflight
	}
}

func runReleasePlanCommand(ctx context.Context, args []string, stdout, stderr io.Writer) int {
	if commandWantsHelp(args) {
		fmt.Fprint(stdout, commandUsage("release plan"))
		return exitOK
	}

	req, err := parseReleasePlanCommand(args)
	if err != nil {
		printReleasePlanFailure(err, stderr)
		return exitPreflight
	}

	source := newReleasePlanGitSource("", preflight.ExecGitRunner{})
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

func parseReleasePlanCommand(args []string) (releasePlanCommandRequest, error) {
	req := releasePlanCommandRequest{outputFormat: releasePlanFormatText}
	fs := flag.NewFlagSet("release plan", flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	var impact string
	fs.StringVar(&req.from, "from", "", "Stable release tag to use as the base")
	fs.StringVar(&req.to, "to", "", "Target revision to analyze")
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
	return encoder.Encode(releasePlanJSONFromPlan(plan))
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

func shortReleasePlanSHA(sha string) string {
	if len(sha) <= 12 {
		return sha
	}
	return sha[:12]
}
