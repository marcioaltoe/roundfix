package cli

import (
	"context"
	"flag"
	"fmt"
	"io"
	"path/filepath"
	"strings"

	"roundfix/internal/app"
	roundconfig "roundfix/internal/config"
	"roundfix/internal/spec"
	"roundfix/internal/speccheck"
)

var archiveUsage = `Usage:
  roundfix archive <slug>
  roundfix archive <slug> --qa-override --approval <source> --reason <text>

Archives a Spec after verifying either every Task is completed and the newest
QA Report has verdict: pass (or a partial verdict covered only by declared
Unreachable Acceptance), or a recorded supersession exists for a Spec without
a Task Graph. Writes one Archive Record and removes the committed Spec folder.
The removed bytes stay in Git at the record's source_revision. The destination
is docs/history/specs/<slug>.md for the built-in root, otherwise
<spec-root>/_archived/<slug>.md.
archive creates no Run and never pushes.

The QA override requires an approval source and reason, still requires every
non-QA Task completed, preserves the QA Task and Reports unchanged, and records
the observed QA outcome and archived HEAD. It is refused when QA already
qualifies for normal archive.

Options:
  --qa-override  Archive despite failed, missing, or unreadable QA evidence
  --approval     Source of the maintainer authority for the override
  --reason       Reason the unmet QA prerequisite is being waived

Exit codes:
  0  archived
  2  Preflight Validation failed
`

type archiveCommandRequest struct {
	slug       string
	qaOverride bool
	approval   string
	reason     string
}

func runArchiveCommand(ctx context.Context, args []string, stdout, stderr io.Writer, environment commandEnvironment) int {
	if commandWantsHelp(args) {
		fmt.Fprint(stdout, archiveUsage)
		return exitOK
	}
	req, err := parseArchiveCommand(args)
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
	pins, err := speccheck.ActiveSpecPathPins(loaded.GitRoot, resolvedSpecsRoot.Path, req.slug)
	if err != nil {
		printPreflightFailure("archive", err, stderr)
		return exitPreflight
	}
	if len(pins) > 0 {
		locations := make([]string, 0, len(pins))
		for _, pin := range pins {
			locations = append(locations, fmt.Sprintf("%s:%d", pin.Path, pin.Line))
		}
		rel, err := filepathRelSlash(loaded.GitRoot, filepath.Join(resolvedSpecsRoot.Path, req.slug))
		if err != nil {
			printPreflightFailure("archive", err, stderr)
			return exitPreflight
		}
		printPreflightFailure("archive", fmt.Errorf("Spec %q cannot archive while another file names its active directory %s/: %s", req.slug, rel, strings.Join(locations, ", ")), stderr)
		return exitPreflight
	}
	findings, err := speccheck.GlossaryFindings(resolvedSpecsRoot.Path, loaded.GitRoot, req.slug)
	if err != nil {
		printPreflightFailure("archive", err, stderr)
		return exitPreflight
	}
	if len(findings) > 0 {
		reasons := make([]string, 0, len(findings))
		for _, finding := range findings {
			reasons = append(reasons, finding.Code+": "+finding.Summary)
		}
		printPreflightFailure("archive", fmt.Errorf("Spec %q cannot archive with a Glossary Gap: %s", req.slug, strings.Join(reasons, "; ")), stderr)
		return exitPreflight
	}
	revisionRoot := loaded.GitRoot
	if resolvedSpecsRoot.External {
		revisionRoot = resolvedSpecsRoot.Path
	}
	repositoryRoot, err := gitOutput(ctx, revisionRoot, "rev-parse", "--show-toplevel")
	if err != nil {
		printPreflightFailure("archive", err, stderr)
		return exitPreflight
	}
	revision, err := gitOutput(ctx, revisionRoot, "rev-parse", "HEAD")
	if err != nil {
		printPreflightFailure("archive", fmt.Errorf("resolve archive revision from HEAD: %w", err), stderr)
		return exitPreflight
	}
	specDir := filepath.Join(resolvedSpecsRoot.Path, req.slug)
	changes, err := gitOutput(ctx, revisionRoot, "status", "--porcelain", "--untracked-files=all", "--", specDir)
	if err != nil {
		printPreflightFailure("archive", err, stderr)
		return exitPreflight
	}
	if strings.TrimSpace(changes) != "" {
		rel, _ := filepathRelSlash(repositoryRoot, specDir)
		printPreflightFailure("archive", fmt.Errorf("Spec %q has changes not committed at HEAD under %s; commit them before archive so the Archive Record's source_revision holds the Spec", req.slug, rel), stderr)
		return exitPreflight
	}
	var qaOverride *spec.QAArchiveOverride
	if req.qaOverride {
		qaOverride = &spec.QAArchiveOverride{Approval: req.approval, Reason: req.reason, Revision: revision}
	}
	result, err := spec.Archive(spec.ArchiveRequest{
		SpecsRoot:      resolvedSpecsRoot.Path,
		BuiltInRoot:    resolvedSpecsRoot.BuiltInRoot,
		Slug:           req.slug,
		QAOverride:     qaOverride,
		SourceRevision: revision,
		RepositoryRoot: repositoryRoot,
	})
	if err != nil {
		printPreflightFailure("archive", err, stderr)
		return exitPreflight
	}
	rel, err := filepathRelSlash(loaded.GitRoot, result.RecordPath)
	if err != nil {
		fmt.Fprintf(stderr, "%s: archive completed but could not format path: %v\n", app.Name, err)
		return exitRunFailed
	}
	suffix := fmt.Sprintf("; removed %d file(s) (%d bytes) kept in Git at %.12s", result.RemovedFiles, result.RemovedBytes, revision)
	if result.QAOverride {
		fmt.Fprintf(stdout, "archived %s with QA override -> %s%s\n", req.slug, rel, suffix)
	} else {
		fmt.Fprintf(stdout, "archived %s -> %s%s\n", req.slug, rel, suffix)
	}
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

func parseArchiveCommand(args []string) (archiveCommandRequest, error) {
	var req archiveCommandRequest
	if len(args) == 0 {
		return req, validationError{message: "missing required Spec slug; pass roundfix archive <slug>"}
	}
	req.slug = strings.TrimSpace(args[0])
	if req.slug == "" {
		return req, validationError{message: "missing required Spec slug; pass roundfix archive <slug>"}
	}
	flags := flag.NewFlagSet("archive", flag.ContinueOnError)
	flags.SetOutput(io.Discard)
	flags.BoolVar(&req.qaOverride, "qa-override", false, "Archive with a QA override")
	flags.StringVar(&req.approval, "approval", "", "Override approval source")
	flags.StringVar(&req.reason, "reason", "", "Override reason")
	if err := flags.Parse(args[1:]); err != nil {
		return req, validationError{message: err.Error()}
	}
	if remaining := flags.Args(); len(remaining) > 0 {
		return req, validationError{message: fmt.Sprintf("unexpected argument %q", remaining[0])}
	}
	req.approval = strings.TrimSpace(req.approval)
	req.reason = strings.TrimSpace(req.reason)
	var approvalSet bool
	var reasonSet bool
	flags.Visit(func(parsed *flag.Flag) {
		switch parsed.Name {
		case "approval":
			approvalSet = true
		case "reason":
			reasonSet = true
		}
	})
	if !req.qaOverride {
		if approvalSet || reasonSet {
			return req, validationError{message: "--approval and --reason require --qa-override"}
		}
		return req, nil
	}
	if req.approval == "" && req.reason == "" {
		return req, validationError{message: "--qa-override requires --approval <source> and --reason <text>"}
	}
	if req.approval == "" {
		return req, validationError{message: "--qa-override requires --approval <source>"}
	}
	if req.reason == "" {
		return req, validationError{message: "--qa-override requires --reason <text>"}
	}
	return req, nil
}
