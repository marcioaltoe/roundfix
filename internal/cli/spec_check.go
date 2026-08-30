package cli

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"roundfix/internal/app"
	roundconfig "roundfix/internal/config"
	"roundfix/internal/daemon"
	"roundfix/internal/spec"
	"roundfix/internal/specaudit"
	"roundfix/internal/speccheck"
)

const specUsage = `Usage:
  roundfix spec check [<slug> ...] [--stage <prd|techspec|tasks>] [--format <text|json>] [--strict] [--run-verification]
  roundfix spec audit <slug> [--format <text|json>]

Commands:
  check  Check Spec artifact consistency
  audit  Audit Spec delivery and surviving Git state
`

const specCheckUsage = `Usage:
  roundfix spec check [<slug> ...] [--stage <prd|techspec|tasks>] [--format <text|json>] [--strict] [--run-verification]

Checks the declarations and citations in the requested Specs without changing
their artifacts. With no slug, checks every active Spec in the Spec Root.

Options:
  --stage   Limit checks to the prd, techspec, or tasks authoring stage
  --format  Output format: text or json (default: text)
  --strict  Promote gap findings to errors; a full sweep reports gate repair inputs separately
  --run-verification
            Execute authored Verification commands in a disposable checkout at HEAD

Exit codes:
  0  no errors (clean or gaps only)
  1  at least one error, vacuous command, or unknown command verdict
  2  usage error or unreadable Spec Root
`

const specAuditUsage = `Usage:
  roundfix spec audit <slug> [--format <text|json>]

Audits one active or archived Spec without changing Git state, the Run
Database, or Spec artifacts. Reports every surviving branch and worktree with
its classification evidence, and gives residue an exact reclaim command.

Options:
  --format  Output format: text or json (default: text)

Exit codes:
  0  no residue or undelivered work
  1  residue or undelivered work found, or the audit could not run
  2  usage error or unknown Spec slug
`

const (
	specCheckFormatText    = "text"
	specCheckFormatJSON    = "json"
	specAuditSchemaVersion = "roundfix-specaudit/v1"
	specAuditDocumentType  = "spec.audit"
)

type specCheckRequest struct {
	slugs           []string
	stage           speccheck.Stage
	outputFormat    string
	strict          bool
	runVerification bool
}

type specCheckVerificationReport struct {
	Executed bool                                 `json:"executed"`
	Tree     string                               `json:"tree,omitempty"`
	Commands []specCheckVerificationCommandReport `json:"commands"`
}

type specCheckVerificationCommandReport struct {
	Task    string `json:"task"`
	Command string `json:"command"`
	Verdict string `json:"verdict"`
	Cause   string `json:"cause,omitempty"`
}

type specCheckDocument struct {
	Schema       string                      `json:"schema"`
	Slug         string                      `json:"slug"`
	Findings     []speccheck.Finding         `json:"findings"`
	RepairInputs []specCheckRepairInput      `json:"repairInputs"`
	Skipped      []speccheck.SkippedDetector `json:"skipped"`
	Verification specCheckVerificationReport `json:"verification"`
}

type specCheckOutcome struct {
	result       speccheck.Result
	repairInputs []specCheckRepairInput
}

type specCheckRepairInput struct {
	Code    string               `json:"code"`
	Summary string               `json:"summary"`
	Where   []speccheck.Location `json:"where"`
	Fix     string               `json:"fix"`
}

const (
	specCheckVerificationTreeHEAD       = "HEAD"
	specCheckVerificationVerdictVacuous = "vacuous"
	specCheckVerificationVerdictHonest  = "honest"
	specCheckVerificationVerdictUnknown = "unknown"
)

type specAuditRequest struct {
	slug         string
	outputFormat string
}

type specAuditDocument struct {
	Schema        string                  `json:"schema"`
	SchemaVersion string                  `json:"schemaVersion"`
	Type          string                  `json:"type"`
	OK            bool                    `json:"ok"`
	Slug          string                  `json:"slug"`
	Survivors     []specaudit.Survivor    `json:"survivors"`
	Undelivered   []specaudit.Undelivered `json:"undelivered"`
}

func runSpecCommand(ctx context.Context, args []string, stdout, stderr io.Writer, environment commandEnvironment) int {
	if len(args) == 0 || (len(args) == 1 && commandWantsHelp(args)) {
		fmt.Fprint(stdout, specUsage)
		return exitOK
	}
	switch args[0] {
	case "check":
		return runSpecCheckCommand(ctx, args[1:], stdout, stderr, environment)
	case "audit":
		return runSpecAuditCommand(ctx, args[1:], stdout, stderr, environment)
	default:
		fmt.Fprintf(stderr, "%s: unknown spec command %q\n", app.Name, args[0])
		fmt.Fprintf(stderr, "Run '%s spec --help' for usage.\n", app.Name)
		return exitPreflight
	}
}

func runSpecAuditCommand(ctx context.Context, args []string, stdout, stderr io.Writer, environment commandEnvironment) int {
	if commandWantsHelp(args) {
		fmt.Fprint(stdout, specAuditUsage)
		return exitOK
	}
	req, err := parseSpecAuditCommand(args)
	if err != nil {
		printSpecAuditFailure(err, stderr, true)
		return exitPreflight
	}
	if err := ctx.Err(); err != nil {
		printSpecAuditFailure(err, stderr, false)
		return exitRunFailed
	}
	loaded, err := loadCommandConfig(environment, stderr)
	if err != nil {
		printSpecAuditFailure(err, stderr, true)
		return exitPreflight
	}
	if loaded.GitRoot == "" {
		printSpecAuditFailure(errors.New("spec audit requires a git repository working tree"), stderr, true)
		return exitPreflight
	}
	resolvedSpecsRoot, err := roundconfig.ResolveSpecsRoot(loaded, loaded.GitRoot)
	if err != nil {
		printSpecAuditFailure(err, stderr, true)
		return exitPreflight
	}
	if err := validateSpecAuditSlug(resolvedSpecsRoot.Path, resolvedSpecsRoot.BuiltInRoot, req.slug); err != nil {
		printSpecAuditFailure(err, stderr, true)
		return exitPreflight
	}

	result, err := specaudit.Audit(ctx, loaded.GitRoot, environment.homeDir, req.slug)
	if err != nil {
		printSpecAuditFailure(err, stderr, false)
		return exitRunFailed
	}
	switch req.outputFormat {
	case specCheckFormatText:
		fmt.Fprint(stdout, renderSpecAuditText(result))
	case specCheckFormatJSON:
		data, err := renderSpecAuditJSON(result)
		if err != nil {
			printSpecAuditFailure(err, stderr, false)
			return exitRunFailed
		}
		fmt.Fprintln(stdout, string(data))
	}
	if specAuditNeedsAttention(result) {
		return exitRunFailed
	}
	return exitOK
}

func runSpecCheckCommand(ctx context.Context, args []string, stdout, stderr io.Writer, environment commandEnvironment) int {
	if commandWantsHelp(args) {
		fmt.Fprint(stdout, specCheckUsage)
		return exitOK
	}
	req, err := parseSpecCheckCommand(args)
	if err != nil {
		printSpecCheckFailure(err, stderr)
		return exitPreflight
	}
	if err := ctx.Err(); err != nil {
		printSpecCheckFailure(err, stderr)
		return exitPreflight
	}
	loaded, err := loadCommandConfig(environment, stderr)
	if err != nil {
		printSpecCheckFailure(err, stderr)
		return exitPreflight
	}
	if loaded.GitRoot == "" {
		printSpecCheckFailure(errors.New("spec check requires a git repository working tree"), stderr)
		return exitPreflight
	}
	resolvedSpecsRoot, err := roundconfig.ResolveSpecsRoot(loaded, loaded.GitRoot)
	if err != nil {
		printSpecCheckFailure(err, stderr)
		return exitPreflight
	}

	slugs := req.slugs
	if len(slugs) == 0 {
		active, skipped, err := spec.ListActiveDetailed(resolvedSpecsRoot.Path)
		if err != nil {
			printSpecCheckFailure(err, stderr)
			return exitPreflight
		}
		printSkippedSpecDiagnostics(stderr, skipped)
		slugs = make([]string, 0, len(active))
		for _, activeSpec := range active {
			slugs = append(slugs, activeSpec.Slug)
		}
	}

	results := make([]specCheckOutcome, 0, len(slugs))
	for _, slug := range slugs {
		if err := validateSpecCheckSlug(resolvedSpecsRoot.Path, slug); err != nil {
			printSpecCheckFailure(err, stderr)
			return exitPreflight
		}
		var result speccheck.Result
		var err error
		if req.stage == speccheck.StageAll {
			result, err = speccheck.Check(resolvedSpecsRoot.Path, loaded.GitRoot, slug)
		} else {
			result, err = speccheck.CheckStage(resolvedSpecsRoot.Path, loaded.GitRoot, slug, req.stage)
		}
		if err != nil {
			printSpecCheckFailure(err, stderr)
			return exitPreflight
		}
		if req.strict {
			speccheck.PromoteGaps(&result)
		}
		results = append(results, classifySpecCheckBoundary(
			result,
			req.strict && req.stage == speccheck.StageAll,
		))
	}

	verificationReports := make([]specCheckVerificationReport, len(slugs))
	for index := range verificationReports {
		verificationReports[index].Commands = []specCheckVerificationCommandReport{}
	}
	if req.runVerification {
		verificationReports, err = probeSpecVerifications(ctx, loaded.GitRoot, resolvedSpecsRoot.Path, slugs)
		if err != nil {
			printSpecCheckFailure(err, stderr)
			return exitRunFailed
		}
	}

	hasError := false
	hasVerificationRefusal := false
	for index, outcome := range results {
		for _, finding := range outcome.result.Findings {
			if finding.Severity == speccheck.SeverityError {
				hasError = true
				break
			}
		}
		switch req.outputFormat {
		case specCheckFormatText:
			fmt.Fprint(stdout, renderSpecCheckText(outcome, verificationReports[index]))
		case specCheckFormatJSON:
			data, err := renderSpecCheckJSON(outcome, verificationReports[index])
			if err != nil {
				fmt.Fprintf(stderr, "%s: spec check failed: %v\n", app.Name, err)
				return exitRunFailed
			}
			fmt.Fprintln(stdout, string(data))
		}
		for _, command := range verificationReports[index].Commands {
			if command.Verdict == specCheckVerificationVerdictVacuous || command.Verdict == specCheckVerificationVerdictUnknown {
				hasVerificationRefusal = true
			}
		}
	}
	if hasError || hasVerificationRefusal {
		return exitRunFailed
	}
	return exitOK
}

func parseSpecCheckCommand(args []string) (specCheckRequest, error) {
	req := specCheckRequest{outputFormat: specCheckFormatText}
	seenSlugs := make(map[string]struct{})
	for index := 0; index < len(args); index++ {
		arg := args[index]
		switch {
		case arg == "--strict":
			req.strict = true
		case arg == "--run-verification":
			req.runVerification = true
		case arg == "--stage":
			index++
			if index >= len(args) {
				return specCheckRequest{}, validationError{message: "--stage requires prd, techspec, or tasks"}
			}
			stage, err := parseSpecCheckStage(args[index])
			if err != nil {
				return specCheckRequest{}, err
			}
			req.stage = stage
		case strings.HasPrefix(arg, "--stage="):
			stage, err := parseSpecCheckStage(strings.TrimPrefix(arg, "--stage="))
			if err != nil {
				return specCheckRequest{}, err
			}
			req.stage = stage
		case arg == "--format":
			index++
			if index >= len(args) {
				return specCheckRequest{}, validationError{message: "--format requires text or json"}
			}
			req.outputFormat = strings.TrimSpace(args[index])
		case strings.HasPrefix(arg, "--format="):
			req.outputFormat = strings.TrimSpace(strings.TrimPrefix(arg, "--format="))
		case strings.HasPrefix(arg, "-"):
			return specCheckRequest{}, validationError{message: fmt.Sprintf("unknown flag %q", arg)}
		default:
			slug := strings.TrimSpace(arg)
			if slug == "" {
				return specCheckRequest{}, validationError{message: "Spec slug must not be empty"}
			}
			if _, seen := seenSlugs[slug]; seen {
				continue
			}
			seenSlugs[slug] = struct{}{}
			req.slugs = append(req.slugs, slug)
		}
	}
	if req.outputFormat != specCheckFormatText && req.outputFormat != specCheckFormatJSON {
		return specCheckRequest{}, validationError{message: fmt.Sprintf("unsupported --format %q; use text or json", req.outputFormat)}
	}
	return req, nil
}

func parseSpecCheckStage(value string) (speccheck.Stage, error) {
	stage := speccheck.Stage(strings.TrimSpace(value))
	switch stage {
	case speccheck.StagePRD, speccheck.StageTechSpec, speccheck.StageTasks:
		return stage, nil
	default:
		return speccheck.StageAll, validationError{
			message: fmt.Sprintf("unsupported --stage %q; use prd, techspec, or tasks", value),
		}
	}
}

func probeSpecVerifications(
	ctx context.Context,
	repoRoot string,
	specsRoot string,
	slugs []string,
) (reports []specCheckVerificationReport, returnErr error) {
	checkoutDir, cleanup, err := speccheck.DisposableCheckout(ctx, repoRoot)
	if err != nil {
		return nil, fmt.Errorf("prepare Verification tree at HEAD: %w", err)
	}
	defer func() {
		returnErr = errors.Join(returnErr, cleanup())
	}()

	diagnosticDir, err := os.MkdirTemp("", "roundfix-spec-check-verification-")
	if err != nil {
		return nil, fmt.Errorf("prepare Verification diagnostics: %w", err)
	}
	defer func() {
		returnErr = errors.Join(returnErr, os.RemoveAll(diagnosticDir))
	}()

	reports = make([]specCheckVerificationReport, len(slugs))
	verifier := daemon.ExecVerifier{}
	for specIndex, slug := range slugs {
		report := specCheckVerificationReport{
			Executed: true,
			Tree:     specCheckVerificationTreeHEAD,
			Commands: []specCheckVerificationCommandReport{},
		}
		manifestPath := filepath.Join(specsRoot, slug, "_tasks.md")
		if _, err := os.Stat(manifestPath); errors.Is(err, os.ErrNotExist) {
			reports[specIndex] = report
			continue
		} else if err != nil {
			return nil, fmt.Errorf("inspect Task Graph for Spec %q: %w", slug, err)
		}

		graph, err := spec.Load(specsRoot, slug)
		if err != nil {
			return nil, fmt.Errorf("load Verification commands for Spec %q: %w", slug, err)
		}
		for taskIndex, task := range graph.Tasks {
			verdicts, err := daemon.ProbeCommands(ctx, verifier, checkoutDir, task.Verification, func(commandIndex int) string {
				return filepath.Join(
					diagnosticDir,
					fmt.Sprintf("spec-%03d-task-%03d-command-%03d.log", specIndex+1, taskIndex+1, commandIndex+1),
				)
			})
			if err != nil {
				return nil, fmt.Errorf("probe Spec %q Task %s Verification: %w", slug, task.ID, err)
			}
			for _, verdict := range verdicts {
				report.Commands = append(report.Commands, specCheckVerificationCommand(task.ID, verdict))
			}
		}
		reports[specIndex] = report
	}
	return reports, nil
}

func specCheckVerificationCommand(taskID string, verdict daemon.CommandVerdict) specCheckVerificationCommandReport {
	report := specCheckVerificationCommandReport{
		Task:    taskID,
		Command: verdict.Command,
		Verdict: specCheckVerificationVerdictHonest,
	}
	switch {
	case verdict.Unknown:
		report.Verdict = specCheckVerificationVerdictUnknown
		if verdict.Cause != nil {
			report.Cause = verdict.Cause.Error()
			var unknownErr *daemon.VerificationUnknownError
			if errors.As(verdict.Cause, &unknownErr) && unknownErr.Err != nil {
				report.Cause = unknownErr.Err.Error()
			}
		}
	case verdict.Vacuous:
		report.Verdict = specCheckVerificationVerdictVacuous
	}
	return report
}

func renderSpecCheckText(outcome specCheckOutcome, verification specCheckVerificationReport) string {
	var report strings.Builder
	executed, unexecuted := 0, 0
	for _, command := range verification.Commands {
		if command.Verdict == specCheckVerificationVerdictUnknown {
			unexecuted++
			continue
		}
		executed++
	}
	report.WriteString(speccheck.RenderText(outcome.result, speccheck.VerificationCoverage{
		Ran:        verification.Executed,
		Executed:   executed,
		Unexecuted: unexecuted,
	}))
	if len(outcome.repairInputs) > 0 {
		report.WriteString("Repair inputs:\n")
		for _, input := range outcome.repairInputs {
			fmt.Fprintf(&report, "[repair input] %s: %s\n", input.Code, input.Summary)
			for _, location := range input.Where {
				fmt.Fprintf(&report, "  at %s:%d\n", location.Path, location.Line)
			}
			fmt.Fprintf(&report, "  fix: %s\n", input.Fix)
		}
	}
	if !verification.Executed {
		return report.String()
	}
	fmt.Fprintf(&report, "Verification tree: %s\n", verification.Tree)
	if len(verification.Commands) == 0 {
		report.WriteString("No authored Verification commands.\n")
		return report.String()
	}
	for _, command := range verification.Commands {
		fmt.Fprintf(&report, "- %s: %s — %q", command.Task, command.Verdict, command.Command)
		switch command.Verdict {
		case specCheckVerificationVerdictVacuous:
			report.WriteString(" (exited zero before work)")
		case specCheckVerificationVerdictHonest:
			report.WriteString(" (exited non-zero before work)")
		case specCheckVerificationVerdictUnknown:
			if command.Cause != "" {
				fmt.Fprintf(&report, " (%s)", command.Cause)
			}
		}
		report.WriteByte('\n')
	}
	return report.String()
}

func renderSpecCheckJSON(outcome specCheckOutcome, verification specCheckVerificationReport) ([]byte, error) {
	findings := outcome.result.Findings
	if findings == nil {
		findings = []speccheck.Finding{}
	}
	repairInputs := outcome.repairInputs
	if repairInputs == nil {
		repairInputs = []specCheckRepairInput{}
	}
	skipped := outcome.result.Skipped
	if skipped == nil {
		skipped = []speccheck.SkippedDetector{}
	}
	if verification.Commands == nil {
		verification.Commands = []specCheckVerificationCommandReport{}
	}
	data, err := json.Marshal(specCheckDocument{
		Schema:       speccheck.SchemaVersion,
		Slug:         outcome.result.Slug,
		Findings:     findings,
		RepairInputs: repairInputs,
		Skipped:      skipped,
		Verification: verification,
	})
	if err != nil {
		return nil, fmt.Errorf("render Spec Consistency Check JSON: %w", err)
	}
	return data, nil
}

func classifySpecCheckBoundary(result speccheck.Result, gateSweep bool) specCheckOutcome {
	outcome := specCheckOutcome{
		result:       result,
		repairInputs: []specCheckRepairInput{},
	}
	if !gateSweep {
		return outcome
	}

	precondition := speccheck.GatePrecondition(result)
	outcome.result.Findings = precondition.Findings
	for _, input := range precondition.Inputs {
		outcome.repairInputs = append(outcome.repairInputs, specCheckRepairInput{
			Code:    input.Code,
			Summary: input.Summary,
			Where:   input.Where,
			Fix:     input.Fix,
		})
	}
	return outcome
}

func parseSpecAuditCommand(args []string) (specAuditRequest, error) {
	req := specAuditRequest{outputFormat: specCheckFormatText}
	for index := 0; index < len(args); index++ {
		arg := args[index]
		switch {
		case arg == "--format":
			index++
			if index >= len(args) {
				return specAuditRequest{}, validationError{message: "--format requires text or json"}
			}
			req.outputFormat = strings.TrimSpace(args[index])
		case strings.HasPrefix(arg, "--format="):
			req.outputFormat = strings.TrimSpace(strings.TrimPrefix(arg, "--format="))
		case strings.HasPrefix(arg, "-"):
			return specAuditRequest{}, validationError{message: fmt.Sprintf("unknown flag %q", arg)}
		case req.slug != "":
			return specAuditRequest{}, validationError{message: fmt.Sprintf("unexpected argument %q", arg)}
		default:
			req.slug = strings.TrimSpace(arg)
		}
	}
	if req.slug == "" {
		return specAuditRequest{}, validationError{message: "spec audit requires one Spec slug"}
	}
	if req.outputFormat != specCheckFormatText && req.outputFormat != specCheckFormatJSON {
		return specAuditRequest{}, validationError{
			message: fmt.Sprintf("unsupported --format %q; use text or json", req.outputFormat),
		}
	}
	return req, nil
}

func validateSpecCheckSlug(specsRoot, slug string) error {
	if slug == "" || filepath.Base(slug) != slug || slug == "." {
		return validationError{message: fmt.Sprintf("invalid Spec slug %q", slug)}
	}
	path := filepath.Join(specsRoot, slug)
	info, err := os.Stat(path)
	if errors.Is(err, os.ErrNotExist) {
		return validationError{message: fmt.Sprintf("unknown Spec slug %q", slug)}
	}
	if err != nil {
		return fmt.Errorf("read Spec %q: %w", slug, err)
	}
	if !info.IsDir() {
		return validationError{message: fmt.Sprintf("Spec slug %q does not name a directory", slug)}
	}
	return nil
}

func validateSpecAuditSlug(specsRoot string, builtInRoot bool, slug string) error {
	if slug == "" || filepath.Base(slug) != slug || slug == "." {
		return validationError{message: fmt.Sprintf("invalid Spec slug %q", slug)}
	}
	for _, path := range []string{
		filepath.Join(specsRoot, slug),
		filepath.Join(spec.ArchiveSpecRoot(specsRoot, builtInRoot), slug),
	} {
		info, err := os.Stat(path)
		if err == nil {
			if !info.IsDir() {
				return validationError{message: fmt.Sprintf("Spec slug %q does not name a directory", slug)}
			}
			return nil
		}
		if !errors.Is(err, os.ErrNotExist) {
			return fmt.Errorf("read Spec %q: %w", slug, err)
		}
	}
	return validationError{message: fmt.Sprintf("unknown Spec slug %q", slug)}
}

func renderSpecAuditText(result specaudit.Result) string {
	var report strings.Builder
	fmt.Fprintf(&report, "Spec audit %s\n", result.Slug)
	if len(result.Survivors) == 0 && len(result.Undelivered) == 0 {
		report.WriteString("No residue or undelivered work.\n")
		return report.String()
	}
	if len(result.Survivors) > 0 {
		report.WriteString("Survivors:\n")
		for _, survivor := range result.Survivors {
			survivorType := "branch"
			if survivor.IsWorktree {
				survivorType = "worktree"
			}
			fmt.Fprintf(&report, "- %s %s %s\n", survivor.Kind, survivorType, survivor.Name)
			fmt.Fprintf(&report, "  evidence: %s\n", survivor.Evidence)
			if survivor.Kind == specaudit.KindResidue {
				fmt.Fprintf(&report, "  reclaim: %s\n", survivor.Reclaim)
			}
		}
	}
	if len(result.Undelivered) > 0 {
		report.WriteString("Undelivered artifacts:\n")
		for _, artifact := range result.Undelivered {
			fmt.Fprintf(&report, "- %s\n", artifact.Artifact)
			if artifact.HeldBy == "" {
				report.WriteString("  held by: no surviving branch\n")
				continue
			}
			fmt.Fprintf(&report, "  held by: %s\n", artifact.HeldBy)
		}
	}
	return report.String()
}

func renderSpecAuditJSON(result specaudit.Result) ([]byte, error) {
	document := specAuditDocument{
		Schema:        specAuditSchemaVersion,
		SchemaVersion: specAuditSchemaVersion,
		Type:          specAuditDocumentType,
		OK:            !specAuditNeedsAttention(result),
		Slug:          result.Slug,
		Survivors:     result.Survivors,
		Undelivered:   result.Undelivered,
	}
	data, err := json.Marshal(document)
	if err != nil {
		return nil, fmt.Errorf("render Spec audit JSON: %w", err)
	}
	return data, nil
}

func specAuditNeedsAttention(result specaudit.Result) bool {
	if len(result.Undelivered) > 0 {
		return true
	}
	for _, survivor := range result.Survivors {
		if survivor.Kind == specaudit.KindResidue {
			return true
		}
	}
	return false
}

func printSpecCheckFailure(err error, stderr io.Writer) {
	fmt.Fprintf(stderr, "%s: spec check failed: %v\n", app.Name, err)
	fmt.Fprintf(stderr, "Run '%s spec check --help' for usage.\n", app.Name)
}

func printSpecAuditFailure(err error, stderr io.Writer, includeUsage bool) {
	fmt.Fprintf(stderr, "%s: spec audit failed: %v\n", app.Name, err)
	if includeUsage {
		fmt.Fprintf(stderr, "Run '%s spec audit --help' for usage.\n", app.Name)
	}
}
