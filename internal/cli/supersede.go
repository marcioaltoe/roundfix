package cli

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"roundfix/internal/app"
	roundconfig "roundfix/internal/config"
	"roundfix/internal/spec"
)

const supersedeUsage = `Usage:
  roundfix supersede --spec <slug> --by <slug> --reason <text>

Records that another active or archived Spec delivered this Spec's content.
Writes only _supersession.md in the superseded Spec. supersede creates no Run,
writes no Run Event Journal entry, and never commits or pushes.

Options:
  --spec    Superseded Spec slug under the configured Spec Root
  --by      Active or archived Spec slug that delivered the content
  --reason  Explanation recorded in the supersession amendment

Exit codes:
  0  supersession recorded
  1  supersession write failed
  2  Preflight Validation failed
`

type supersedeRequest struct {
	specSlug string
	bySlug   string
	reason   string
}

func runSupersedeCommand(ctx context.Context, args []string, stdout, stderr io.Writer, environment commandEnvironment) int {
	if commandWantsHelp(args) {
		fmt.Fprint(stdout, supersedeUsage)
		return exitOK
	}
	req, err := parseSupersedeCommand(args)
	if err != nil {
		printPreflightFailure("supersede", err, stderr)
		return exitPreflight
	}
	resolvedRoot, err := preflightSupersede(ctx, req, stderr, environment)
	if err != nil {
		printPreflightFailure("supersede", err, stderr)
		return exitPreflight
	}
	_, err = spec.WriteSupersession(filepath.Join(resolvedRoot.Path, req.specSlug), req.bySlug, req.reason, time.Now())
	if errors.Is(err, spec.ErrSupersessionExists) {
		printPreflightFailure("supersede", fmt.Errorf("superseded Spec %q already carries a supersession", req.specSlug), stderr)
		return exitPreflight
	}
	if err != nil {
		fmt.Fprintf(stderr, "%s: supersede failed: %v\n", app.Name, err)
		return exitRunFailed
	}
	fmt.Fprintf(stdout, "superseded %s by %s\n", req.specSlug, req.bySlug)
	return exitOK
}

func parseSupersedeCommand(args []string) (supersedeRequest, error) {
	var req supersedeRequest
	flags := flag.NewFlagSet("supersede", flag.ContinueOnError)
	flags.SetOutput(io.Discard)
	flags.StringVar(&req.specSlug, "spec", "", "Superseded Spec slug")
	flags.StringVar(&req.bySlug, "by", "", "Superseding Spec slug")
	flags.StringVar(&req.reason, "reason", "", "Supersession explanation")
	if err := flags.Parse(args); err != nil {
		return req, validationError{message: err.Error()}
	}
	if remaining := flags.Args(); len(remaining) > 0 {
		return req, validationError{message: fmt.Sprintf("unexpected argument %q", remaining[0])}
	}
	req.specSlug = strings.TrimSpace(req.specSlug)
	req.bySlug = strings.TrimSpace(req.bySlug)
	req.reason = strings.TrimSpace(req.reason)
	if req.specSlug == "" {
		return req, validationError{message: "missing required --spec; pass --spec <slug>"}
	}
	if err := validateSupersedeSlug("--spec", req.specSlug); err != nil {
		return req, err
	}
	if req.bySlug == "" {
		return req, validationError{message: "missing required --by; pass --by <slug>"}
	}
	if err := validateSupersedeSlug("--by", req.bySlug); err != nil {
		return req, err
	}
	if req.reason == "" {
		return req, validationError{message: "missing required --reason; pass --reason <text>"}
	}
	return req, nil
}

func validateSupersedeSlug(flagName string, slug string) error {
	if slug == "." || slug == ".." || filepath.IsAbs(slug) || filepath.Base(slug) != slug || strings.ContainsAny(slug, `/\`) {
		return validationError{message: fmt.Sprintf("invalid Spec slug %q; %s must be one Spec directory name", slug, flagName)}
	}
	return nil
}

func preflightSupersede(ctx context.Context, req supersedeRequest, stderr io.Writer, environment commandEnvironment) (roundconfig.SpecsRoot, error) {
	loaded, err := loadCommandConfig(environment, stderr)
	if err != nil {
		return roundconfig.SpecsRoot{}, err
	}
	if loaded.GitRoot == "" {
		return roundconfig.SpecsRoot{}, validationError{message: "supersede requires a git repository working tree"}
	}
	if err := ctx.Err(); err != nil {
		return roundconfig.SpecsRoot{}, err
	}
	resolvedRoot, err := roundconfig.ResolveSpecsRoot(loaded, loaded.GitRoot)
	if err != nil {
		return roundconfig.SpecsRoot{}, err
	}
	supersededDir := filepath.Join(resolvedRoot.Path, req.specSlug)
	known, err := knownSpecDirectory(supersededDir)
	if err != nil {
		return roundconfig.SpecsRoot{}, err
	}
	if !known {
		return roundconfig.SpecsRoot{}, validationError{message: fmt.Sprintf("superseded Spec %q is unknown", req.specSlug)}
	}
	if req.specSlug == req.bySlug {
		return roundconfig.SpecsRoot{}, validationError{message: fmt.Sprintf("Spec %q cannot supersede itself", req.specSlug)}
	}
	knownDeliverer, err := knownDelivererSpec(resolvedRoot, req.bySlug)
	if err != nil {
		return roundconfig.SpecsRoot{}, err
	}
	if !knownDeliverer {
		return roundconfig.SpecsRoot{}, validationError{message: fmt.Sprintf("superseding Spec %q is neither active nor archived", req.bySlug)}
	}
	if _, err := os.Stat(filepath.Join(supersededDir, spec.SupersessionFilename)); err == nil {
		return roundconfig.SpecsRoot{}, validationError{message: fmt.Sprintf("superseded Spec %q already carries a supersession", req.specSlug)}
	} else if !errors.Is(err, os.ErrNotExist) {
		return roundconfig.SpecsRoot{}, fmt.Errorf("stat supersession record for Spec %q: %w", req.specSlug, err)
	}
	return resolvedRoot, nil
}

func knownDelivererSpec(root roundconfig.SpecsRoot, slug string) (bool, error) {
	activeDir := filepath.Join(root.Path, slug)
	status, err := spec.ReadPRDStatus(activeDir)
	if err == nil {
		if status != "active" {
			return false, validationError{message: fmt.Sprintf("superseding Spec %q is not active: _prd.md frontmatter status is %q; expected %q", slug, status, "active")}
		}
		return true, nil
	}
	if !errors.Is(err, os.ErrNotExist) {
		return false, validationError{message: fmt.Sprintf("superseding Spec %q has malformed _prd.md in the active Spec Root: %v", slug, err)}
	}

	archivedDir := filepath.Join(spec.ArchiveSpecRoot(root.Path, root.BuiltInRoot), slug)
	status, err = spec.ReadPRDStatus(archivedDir)
	if errors.Is(err, os.ErrNotExist) {
		return false, nil
	}
	if err != nil {
		return false, validationError{message: fmt.Sprintf("superseding Spec %q has malformed _prd.md in the archive: %v", slug, err)}
	}
	if status == "archived" {
		return true, nil
	}
	if _, err := spec.ReadSupersession(archivedDir); err == nil {
		return true, nil
	} else if !errors.Is(err, spec.ErrNoSupersession) {
		return false, validationError{message: fmt.Sprintf("superseding Spec %q has malformed supersession record in the archive: %v", slug, err)}
	}
	return false, validationError{message: fmt.Sprintf("superseding Spec %q is not archived: archived _prd.md frontmatter status is %q; expected %q", slug, status, "archived")}
}

func knownSpecDirectory(directory string) (bool, error) {
	info, err := os.Stat(filepath.Join(directory, "_prd.md"))
	if errors.Is(err, os.ErrNotExist) {
		return false, nil
	}
	if err != nil {
		return false, fmt.Errorf("inspect Spec directory %q: %w", directory, err)
	}
	return !info.IsDir(), nil
}
