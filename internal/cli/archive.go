package cli

import (
	"context"
	"fmt"
	"io"
	"path/filepath"
	"strings"

	"roundfix/internal/app"
	roundconfig "roundfix/internal/config"
	"roundfix/internal/spec"
)

const archiveUsage = `Usage:
  roundfix archive <slug>

Archives a completed Spec after verifying every Task is completed and the
newest QA Report has verdict: pass or a partial verdict whose blocked rows are
covered only by declared Unreachable Acceptance. Stamps archive metadata and
moves docs/specs/<slug>/ to
docs/specs/_archived/<slug>/. archive creates no Run and never pushes.

Exit codes:
  0  archived
  2  Preflight Validation failed
`

func runArchiveCommand(ctx context.Context, args []string, stdout, stderr io.Writer, environment commandEnvironment) int {
	if commandWantsHelp(args) {
		fmt.Fprint(stdout, archiveUsage)
		return exitOK
	}
	slug, err := parseArchiveCommand(args)
	if err != nil {
		printPreflightFailure("archive", err, stderr)
		return exitPreflight
	}
	loaded, err := loadCommandConfig(environment, stderr)
	if err != nil {
		printPreflightFailure("archive", err, stderr)
		return exitPreflight
	}
	if loaded.GitRoot == "" {
		printPreflightFailure("archive", fmt.Errorf("archive requires a git repository working tree"), stderr)
		return exitPreflight
	}
	if err := ctx.Err(); err != nil {
		printPreflightFailure("archive", err, stderr)
		return exitPreflight
	}
	resolvedSpecsRoot, err := roundconfig.ResolveSpecsRoot(loaded, loaded.GitRoot)
	if err != nil {
		printPreflightFailure("archive", err, stderr)
		return exitPreflight
	}
	result, err := spec.Archive(spec.ArchiveRequest{
		SpecsRoot: resolvedSpecsRoot.Path,
		Slug:      slug,
	})
	if err != nil {
		printPreflightFailure("archive", err, stderr)
		return exitPreflight
	}
	rel, err := filepathRelSlash(loaded.GitRoot, result.ArchivedDir)
	if err != nil {
		fmt.Fprintf(stderr, "%s: archive completed but could not format path: %v\n", app.Name, err)
		return exitRunFailed
	}
	fmt.Fprintf(stdout, "archived %s -> %s\n", slug, rel)
	return exitOK
}

func filepathRelSlash(base string, target string) (string, error) {
	rel, err := filepath.Rel(base, target)
	if err != nil {
		return "", err
	}
	if rel == ".." || strings.HasPrefix(rel, ".."+string(filepath.Separator)) || filepath.IsAbs(rel) {
		return filepath.ToSlash(target), nil
	}
	return filepath.ToSlash(rel), nil
}

func parseArchiveCommand(args []string) (string, error) {
	if len(args) == 0 {
		return "", validationError{message: "missing required Spec slug; pass roundfix archive <slug>"}
	}
	if len(args) > 1 {
		return "", validationError{message: fmt.Sprintf("unexpected argument %q", args[1])}
	}
	slug := strings.TrimSpace(args[0])
	if slug == "" {
		return "", validationError{message: "missing required Spec slug; pass roundfix archive <slug>"}
	}
	return slug, nil
}
