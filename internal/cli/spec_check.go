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
	"roundfix/internal/spec"
	"roundfix/internal/specaudit"
	"roundfix/internal/speccheck"
)

const specUsage = `Usage:
  roundfix spec check [<slug> ...] [--format <text|json>] [--strict]
  roundfix spec audit <slug> [--format <text|json>]

Commands:
  check  Check Spec artifact consistency
  audit  Audit Spec delivery and surviving Git state
`

const specCheckUsage = `Usage:
  roundfix spec check [<slug> ...] [--format <text|json>] [--strict]

Checks the declarations and citations in the requested Specs without changing
their artifacts. With no slug, checks every active Spec in the Spec Root.

Options:
  --format  Output format: text or json (default: text)
  --strict  Promote gap findings to errors

Exit codes:
  0  no errors (clean or gaps only)
  1  at least one error
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
	slugs        []string
	outputFormat string
	strict       bool
}

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
	if err := validateSpecAuditSlug(resolvedSpecsRoot.Path, req.slug); err != nil {
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

	results := make([]speccheck.Result, 0, len(slugs))
	for _, slug := range slugs {
		if err := validateSpecCheckSlug(resolvedSpecsRoot.Path, slug); err != nil {
			printSpecCheckFailure(err, stderr)
			return exitPreflight
		}
		result, err := speccheck.Check(resolvedSpecsRoot.Path, loaded.GitRoot, slug)
		if err != nil {
			printSpecCheckFailure(err, stderr)
			return exitPreflight
		}
		if req.strict {
			promoteSpecCheckGaps(&result)
		}
		results = append(results, result)
	}

	hasError := false
	for _, result := range results {
		for _, finding := range result.Findings {
			if finding.Severity == speccheck.SeverityError {
				hasError = true
				break
			}
		}
		switch req.outputFormat {
		case specCheckFormatText:
			fmt.Fprint(stdout, speccheck.RenderText(result))
		case specCheckFormatJSON:
			data, err := speccheck.RenderJSON(result)
			if err != nil {
				fmt.Fprintf(stderr, "%s: spec check failed: %v\n", app.Name, err)
				return exitRunFailed
			}
			fmt.Fprintln(stdout, string(data))
		}
	}
	if hasError {
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

func validateSpecAuditSlug(specsRoot, slug string) error {
	if slug == "" || filepath.Base(slug) != slug || slug == "." {
		return validationError{message: fmt.Sprintf("invalid Spec slug %q", slug)}
	}
	for _, path := range []string{
		filepath.Join(specsRoot, slug),
		filepath.Join(specsRoot, "_archived", slug),
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

func promoteSpecCheckGaps(result *speccheck.Result) {
	for index := range result.Findings {
		if result.Findings[index].Severity == speccheck.SeverityGap {
			result.Findings[index].Severity = speccheck.SeverityError
		}
	}
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
