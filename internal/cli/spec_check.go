package cli

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"roundfix/internal/app"
	roundconfig "roundfix/internal/config"
	"roundfix/internal/spec"
	"roundfix/internal/speccheck"
)

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

const (
	specCheckFormatText = "text"
	specCheckFormatJSON = "json"
)

type specCheckRequest struct {
	slugs        []string
	outputFormat string
	strict       bool
}

func runSpecCommand(ctx context.Context, args []string, stdout, stderr io.Writer, environment commandEnvironment) int {
	if len(args) == 0 || commandWantsHelp(args) {
		fmt.Fprint(stdout, specCheckUsage)
		return exitOK
	}
	if args[0] != "check" {
		fmt.Fprintf(stderr, "%s: unknown spec command %q\n", app.Name, args[0])
		fmt.Fprintf(stderr, "Run '%s spec check --help' for usage.\n", app.Name)
		return exitPreflight
	}
	return runSpecCheckCommand(ctx, args[1:], stdout, stderr, environment)
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
