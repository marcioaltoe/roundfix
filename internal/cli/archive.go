package cli

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"roundfix/internal/app"
	roundconfig "roundfix/internal/config"
	"roundfix/internal/judge"
	"roundfix/internal/preflight"
	"roundfix/internal/spec"
	"roundfix/internal/speccheck"
)

var archiveUsage = `Usage:
  roundfix archive <slug>
  roundfix archive <slug> --qa-override --approval <source> --reason <text>
  roundfix archive <slug> --plan [--offline]
  roundfix archive <slug> --pull-request <number>

Archives a Spec after verifying either every Task is completed and the newest
QA Report has verdict: pass (or a partial verdict covered only by declared
Unreachable Acceptance), or a recorded supersession exists for a Spec without
a Task Graph. Writes one Archive Record and removes the committed Spec folder.
The removed bytes stay in Git at the record's source_revision. The destination
is docs/history/specs/<slug>.md for the built-in root, otherwise
<spec-root>/_archived/<slug>.md.
archive creates no Run and never pushes.
Refuses, writing nothing, while the Spec folder has uncommitted changes or Git ignores a file under it.

The QA override requires an approval source and reason, still requires every
non-QA Task completed, preserves the QA Task and Reports unchanged, and records
the observed QA outcome and archived HEAD. It is refused when QA already
qualifies for normal archive.

Options:
  --offline      With --plan, skip Archive Advice and contact no judge
  --qa-override  Archive despite failed, missing, or unreadable QA evidence
  --approval     Source of the maintainer authority for the override
  --reason       Reason the unmet QA prerequisite is being waived
  --pull-request Store the delivery pull request number in the Archive Record

Exit codes:
  0  archived
  1  archive could not finish after writing; stderr says how to finish
  2  Preflight Validation failed
`

type archiveCommandRequest struct {
	slug        string
	qaOverride  bool
	approval    string
	reason      string
	plan        bool
	offline     bool
	promote     []string
	pullRequest string
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
	if req.plan {
		if req.qaOverride || req.approval != "" || req.reason != "" || len(req.promote) > 0 {
			printPreflightFailure("archive", validationError{message: "--plan cannot be combined with --promote, --qa-override, --approval or --reason"}, stderr)
			return exitPreflight
		}
		return runArchivePlan(ctx, req, loaded, stdout, stderr, environment)
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
	revisionRoot := loaded.GitRoot
	if resolvedSpecsRoot.External {
		revisionRoot = resolvedSpecsRoot.Path
	}
	if req.pullRequest != "" {
		recordOnly, err := recordArchiveCommandPullRequest(ctx, req, resolvedSpecsRoot, revisionRoot, loaded.GitRoot, stdout)
		if err != nil {
			printPreflightFailure("archive", err, stderr)
			return exitPreflight
		}
		if recordOnly {
			return exitOK
		}
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
	status, err := readSpecFolderStatus(ctx, repositoryRoot, specDir)
	if err != nil {
		printPreflightFailure("archive", err, stderr)
		return exitPreflight
	}
	statusFolder := specDir
	if len(status.Changes) == 0 && len(status.Ignored) > 0 {
		// Git reports the physical repository root; the configured root may
		// contain a symlink. Keep the existing uncommitted refusal unchanged.
		statusFolder, err = filepath.EvalSymlinks(specDir)
		if err != nil {
			printPreflightFailure("archive", fmt.Errorf("resolve ignored Spec folder: %w", err), stderr)
			return exitPreflight
		}
	}
	folderRel, err := filepathRelSlash(repositoryRoot, statusFolder)
	if err != nil {
		printPreflightFailure("archive", err, stderr)
		return exitPreflight
	}
	if err := refuseSpecFolderStatus("archive", req.slug, folderRel, status); err != nil {
		printPreflightFailure("archive", err, stderr)
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
		PullRequest:    req.pullRequest,
		RepositoryRoot: repositoryRoot,
		Promote:        req.promote,
	})
	var removal *spec.FolderRemovalError
	if err != nil && !errors.As(err, &removal) {
		printPreflightFailure("archive", err, stderr)
		return exitPreflight
	}
	if req.pullRequest != "" {
		if _, err := spec.RecordArchivePullRequest(result.RecordPath, req.pullRequest); err != nil {
			fmt.Fprintf(stderr, "%s: archive completed but could not record pull request: %v\n", app.Name, err)
			return exitRunFailed
		}
	}
	rel, err := filepathRelSlash(loaded.GitRoot, result.RecordPath)
	if err != nil {
		fmt.Fprintf(stderr, "%s: archive completed but could not format path: %v\n", app.Name, err)
		return exitRunFailed
	}
	if removal != nil {
		printFolderRemovalRecovery("archive", rel, folderRel, removal, stderr)
		return exitRunFailed
	}
	suffix := fmt.Sprintf("; removed %d file(s) (%d bytes) kept in Git at %.12s", result.RemovedFiles, result.RemovedBytes, revision)
	if len(result.Promoted) > 0 {
		suffix += fmt.Sprintf("; promoted %d file(s) to docs/references/", len(result.Promoted))
	}
	if len(result.Relinked) > 0 {
		suffix += fmt.Sprintf("; relinked %d file(s)", len(result.Relinked))
	}
	if req.pullRequest != "" {
		suffix += "; pull request " + req.pullRequest
	}
	warnArchiveSourceRevision(ctx, revisionRoot, revision, req.pullRequest, stderr)
	for _, citation := range result.Unrelinked {
		fmt.Fprintf(stderr, "%s: warning: archive leaves %s:%d naming removed %s unrelinked\n", app.Name, citation.Path, citation.Line, citation.Target)
	}
	if result.QAOverride {
		fmt.Fprintf(stdout, "archived %s with QA override -> %s%s\n", req.slug, rel, suffix)
	} else {
		fmt.Fprintf(stdout, "archived %s -> %s%s\n", req.slug, rel, suffix)
	}
	return exitOK
}

// archiveDefaultRef uses origin/HEAD when present and the conventional current
// branch name otherwise, matching Preflight Validation's default detection.
func archiveDefaultRef(ctx context.Context, repo string) (ref, display string) {
	branch, err := gitOutput(ctx, repo, "symbolic-ref", "--short", "HEAD")
	if err != nil {
		branch = ""
	}
	detection := preflight.DetectDefaultBranch(ctx, repo, branch, nil)
	switch detection.Source {
	case preflight.DefaultBranchFromOriginHead:
		return "refs/remotes/origin/" + detection.Name, "origin/" + detection.Name
	case preflight.DefaultBranchFromNameMatch:
		return "refs/heads/" + detection.Name, detection.Name
	default:
		return "", ""
	}
}

// recordArchiveCommandPullRequest decides the record-only form before any
// active-Spec checks. A default-branch record is never rewritten.
func recordArchiveCommandPullRequest(ctx context.Context, req archiveCommandRequest, root roundconfig.SpecsRoot, repo, outputRoot string, stdout io.Writer) (bool, error) {
	ref, display := archiveDefaultRef(ctx, repo)
	if ref == "" {
		return false, errors.New("the default branch is unknown; cannot record a pull request")
	}
	if _, err := gitOutput(ctx, repo, "rev-parse", "--verify", ref+"^{commit}"); err != nil {
		return false, fmt.Errorf("resolve default branch %s: %w", display, err)
	}
	specDir := filepath.Join(root.Path, req.slug)
	if _, err := os.Lstat(specDir); err == nil {
		return false, nil
	} else if !errors.Is(err, os.ErrNotExist) {
		return false, fmt.Errorf("stat active Spec: %w", err)
	}
	recordPath := spec.ArchiveRecordPath(spec.ArchiveSpecRoot(root.Path, root.BuiltInRoot), req.slug)
	if _, err := os.Lstat(recordPath); errors.Is(err, os.ErrNotExist) {
		return false, nil
	} else if err != nil {
		return false, fmt.Errorf("stat Archive Record: %w", err)
	}
	if len(req.promote) > 0 || req.qaOverride || req.approval != "" || req.reason != "" {
		return true, errors.New("record-only --pull-request cannot be combined with --promote, --qa-override, --approval or --reason")
	}
	repositoryRoot, err := gitOutput(ctx, repo, "rev-parse", "--show-toplevel")
	if err != nil {
		return true, err
	}
	recordRel, err := filepathRelSlash(repositoryRoot, recordPath)
	if err != nil {
		return true, err
	}
	// ls-tree distinguishes an absent path from an unreadable default ref.
	entry, err := gitOutput(ctx, repositoryRoot, "ls-tree", ref, "--", recordRel)
	if err != nil {
		return true, fmt.Errorf("read default-branch Archive Record: %w", err)
	}
	if entry != "" {
		return true, fmt.Errorf("Archive Record %s is already on %s; an archived record stays byte-identical", recordRel, display)
	}
	archived, err := spec.ReadArchivedSpec(spec.ArchiveSpecRoot(root.Path, root.BuiltInRoot), req.slug)
	if err != nil {
		return true, fmt.Errorf("read Archive Record: %w", err)
	}
	if archived.Form != spec.ArchivedRecord {
		return true, errors.New("--pull-request requires an Archive Record")
	}
	rel, err := filepathRelSlash(outputRoot, recordPath)
	if err != nil {
		return true, err
	}
	changed, err := spec.RecordArchivePullRequest(recordPath, req.pullRequest)
	if err != nil {
		return true, err
	}
	if changed {
		fmt.Fprintf(stdout, "recorded pull request %s in %s\n", req.pullRequest, rel)
	} else {
		fmt.Fprintf(stdout, "%s already records pull request %s\n", rel, req.pullRequest)
	}
	return true, nil
}

func warnArchiveSourceRevision(ctx context.Context, repo, revision, number string, stderr io.Writer) {
	ref, display := archiveDefaultRef(ctx, repo)
	if ref == "" {
		fmt.Fprintf(stderr, "%s: warning: the default branch is unknown, so source_revision %.12s was not checked\n", app.Name, revision)
		return
	}
	if _, err := gitOutput(ctx, repo, "merge-base", "--is-ancestor", revision, ref); err != nil {
		var exitErr *exec.ExitError
		if !errors.As(err, &exitErr) || exitErr.ExitCode() != 1 {
			fmt.Fprintf(stderr, "%s: warning: source_revision %.12s could not be checked against %s: %v\n", app.Name, revision, display, err)
			return
		}
		if number == "" {
			fmt.Fprintf(stderr, "%s: warning: source_revision %.12s is not an ancestor of %s; a squash merge keeps it only in the pull request's head ref, so record the pull request with --pull-request once it opens\n", app.Name, revision, display)
		} else {
			fmt.Fprintf(stderr, "%s: warning: source_revision %.12s is not an ancestor of %s; readers recover it from refs/pull/%s/head\n", app.Name, revision, display, number)
		}
	}
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
	flags.BoolVar(&req.plan, "plan", false, "Plan the archive without changing the repository")
	flags.BoolVar(&req.offline, "offline", false, "With --plan, skip Archive Advice and contact no judge")
	flags.StringVar(&req.approval, "approval", "", "Override approval source")
	flags.StringVar(&req.reason, "reason", "", "Override reason")
	flags.StringVar(&req.pullRequest, "pull-request", "", "Store the delivery pull request number in the Archive Record")
	flags.Func("promote", "Promote a Spec-relative file to docs/references/ (repeatable)", func(value string) error { req.promote = append(req.promote, value); return nil })
	if err := flags.Parse(args[1:]); err != nil {
		return req, validationError{message: err.Error()}
	}
	if remaining := flags.Args(); len(remaining) > 0 {
		return req, validationError{message: fmt.Sprintf("unexpected argument %q", remaining[0])}
	}
	if req.offline && !req.plan {
		return req, validationError{message: "--offline requires --plan"}
	}
	req.approval = strings.TrimSpace(req.approval)
	req.reason = strings.TrimSpace(req.reason)
	if filepath.Base(req.slug) != req.slug || req.slug == "." || req.slug == ".." {
		return req, validationError{message: fmt.Sprintf("unsafe Spec slug %q", req.slug)}
	}
	var pullRequestSet bool
	var approvalSet bool
	var reasonSet bool
	flags.Visit(func(parsed *flag.Flag) {
		switch parsed.Name {
		case "pull-request":
			pullRequestSet = true
		case "approval":
			approvalSet = true
		case "reason":
			reasonSet = true
		}
	})
	if pullRequestSet && !regexp.MustCompile(`^[1-9][0-9]{0,9}$`).MatchString(req.pullRequest) {
		return req, validationError{message: "--pull-request requires a positive integer matching ^[1-9][0-9]{0,9}$"}
	}
	if req.plan {
		if pullRequestSet {
			return req, validationError{message: "--pull-request cannot be combined with --plan"}
		}
		if len(req.promote) > 0 {
			return req, validationError{message: "--plan cannot be combined with --promote"}
		}
		if approvalSet || reasonSet {
			return req, validationError{message: "--approval and --reason cannot be used with --plan"}
		}
		return req, nil
	}
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

func runArchivePlan(ctx context.Context, req archiveCommandRequest, loaded roundconfig.Loaded, stdout, stderr io.Writer, environment commandEnvironment) int {
	root, err := roundconfig.ResolveSpecsRoot(loaded, loaded.GitRoot)
	if err != nil {
		printPreflightFailure("archive", err, stderr)
		return exitPreflight
	}
	specDir := filepath.Join(root.Path, req.slug)
	info, err := os.Stat(specDir)
	if err != nil || !info.IsDir() {
		printPreflightFailure("archive", fmt.Errorf("unknown active Spec slug %q", req.slug), stderr)
		return exitPreflight
	}
	files, err := archivePlanFiles(specDir)
	if err != nil {
		printPreflightFailure("archive", err, stderr)
		return exitPreflight
	}
	var core, evidence, candidates []string
	var evidenceBytes, totalBytes int64
	for _, name := range files {
		path := filepath.Join(specDir, name)
		stat, statErr := os.Stat(path)
		if statErr != nil {
			printPreflightFailure("archive", statErr, stderr)
			return exitPreflight
		}
		totalBytes += stat.Size()
		switch {
		case judgeArchiveCore(name):
			core = append(core, name)
		case strings.HasPrefix(filepath.ToSlash(name), "qa/evidence/"):
			evidence = append(evidence, name)
			evidenceBytes += stat.Size()
		default:
			candidates = append(candidates, name)
		}
	}
	var report judge.ArchiveAdviceReport
	if !req.offline {
		q, err := judge.Load()
		if err != nil {
			printPreflightFailure("archive", err, stderr)
			return exitPreflight
		}
		keys := map[string]string{}
		allowed := map[string]bool{}
		for _, key := range q.KeyVariables() {
			allowed[key] = true
		}
		for _, entry := range environment.environ {
			name, value, _ := strings.Cut(entry, "=")
			if allowed[name] {
				keys[name] = value
			}
		}
		report, err = judge.AdviseArchive(ctx, q.WithMonthlyCeiling(loaded.Config.Jev.MonthlyCeilingUSD), judge.ArchiveAdviceRequest{RepoRoot: loaded.GitRoot, SpecDir: specDir, Spec: req.slug, Files: candidates, Keys: keys, HomeDir: loaded.HomeDir, Transport: environment.dependencies.judgeTransport, Now: environment.dependencies.judgeNow})
		if err != nil {
			printPreflightFailure("archive", err, stderr)
			return exitPreflight
		}
	}
	relRecord := filepath.ToSlash(filepath.Join(spec.ArchiveSpecRoot(root.Path, root.BuiltInRoot), req.slug+".md"))
	relRecord, _ = filepathRelSlash(loaded.GitRoot, relRecord)
	repositoryRoot := loaded.GitRoot
	if root.External {
		repositoryRoot, err = gitOutput(ctx, root.Path, "rev-parse", "--show-toplevel")
		if err != nil {
			printPreflightFailure("archive", err, stderr)
			return exitPreflight
		}
	}
	specRootRel, err := filepathRelSlash(repositoryRoot, root.Path)
	if err != nil {
		printPreflightFailure("archive", err, stderr)
		return exitPreflight
	}
	archiveRootRel, err := filepathRelSlash(repositoryRoot, spec.ArchiveSpecRoot(root.Path, root.BuiltInRoot))
	if err != nil {
		printPreflightFailure("archive", err, stderr)
		return exitPreflight
	}
	folderRel, err := filepathRelSlash(repositoryRoot, specDir)
	if err != nil {
		printPreflightFailure("archive", err, stderr)
		return exitPreflight
	}
	recordRel, err := filepathRelSlash(repositoryRoot, spec.ArchiveRecordPath(spec.ArchiveSpecRoot(root.Path, root.BuiltInRoot), req.slug))
	if err != nil {
		printPreflightFailure("archive", err, stderr)
		return exitPreflight
	}
	relinkPlan, err := spec.PlanArchiveRelink(ctx, repositoryRoot, specRootRel, archiveRootRel, folderRel, recordRel, nil)
	if err != nil {
		printPreflightFailure("archive", err, stderr)
		return exitPreflight
	}
	fmt.Fprintf(stdout, "archive plan for %s: removes %d file(s) (%d bytes) and writes %s\n", req.slug, len(files), totalBytes, relRecord)
	fmt.Fprintf(stdout, "core %d file(s): %s\n", len(core), strings.Join(core, ", "))
	fmt.Fprintf(stdout, "evidence %d file(s) (%d bytes) under qa/evidence/\n", len(evidence), evidenceBytes)
	if req.offline {
		fmt.Fprintln(stdout, "judge: skipped (offline)")
		for _, name := range candidates {
			stat, err := os.Stat(filepath.Join(specDir, filepath.FromSlash(name)))
			if err != nil {
				printPreflightFailure("archive", err, stderr)
				return exitPreflight
			}
			fmt.Fprintf(stdout, "candidate %s %d bytes: no advice (offline)\n", name, stat.Size())
		}
	}
	for _, advice := range report.Advice {
		stat, _ := os.Stat(filepath.Join(specDir, filepath.FromSlash(advice.File)))
		size := int64(0)
		if stat != nil {
			size = stat.Size()
		}
		detail := "no advice"
		if advice.Choice != nil {
			detail = *advice.Choice
			if advice.Probabilities != nil {
				if probability, ok := advice.Probabilities[*advice.Choice]; ok {
					detail += fmt.Sprintf(" (p=%.2f)", probability)
				}
			}
		} else if advice.Reason != nil {
			detail += " (" + *advice.Reason + ")"
		}
		fmt.Fprintf(stdout, "candidate %s %d bytes: %s\n", advice.File, size, detail)
	}
	promotions := map[string]string{}
	for _, path := range relinkPlan.Promote {
		destination := "docs/references/" + filepath.Base(path)
		promotions[folderRel+"/"+path] = destination
		fmt.Fprintf(stdout, "promotes %s to %s\n", path, destination)
	}
	folderFiles := make([]string, 0, len(files))
	for _, path := range files {
		folderFiles = append(folderFiles, folderRel+"/"+path)
	}
	moves := spec.ArchiveCitationMoves(folderRel, recordRel, folderFiles, promotions)
	var relinked []string
	for path := range relinkPlan.Relinked {
		relinked = append(relinked, path)
	}
	sort.Strings(relinked)
	for _, path := range relinked {
		content, err := gitOutput(ctx, repositoryRoot, "show", "HEAD:"+path)
		if err != nil {
			printPreflightFailure("archive", err, stderr)
			return exitPreflight
		}
		_, count := spec.RelinkContent([]byte(content), path, moves)
		fmt.Fprintf(stdout, "relinks %s: %d citation(s)\n", path, count)
	}
	for _, citation := range relinkPlan.Unrelinked {
		fmt.Fprintf(stdout, "leaves %s:%d naming %s unrelinked\n", citation.Path, citation.Line, citation.Target)
	}
	fmt.Fprintf(stdout, "promote with: roundfix archive %s --promote <path>\n", req.slug)
	return exitOK
}

func archivePlanFiles(dir string) ([]string, error) {
	var files []string
	err := filepath.WalkDir(dir, func(path string, entry os.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if entry.IsDir() {
			return nil
		}
		if entry.Type()&os.ModeSymlink != 0 {
			return fmt.Errorf("Spec contains symlink %q", path)
		}
		rel, err := filepath.Rel(dir, path)
		if err != nil {
			return err
		}
		files = append(files, filepath.ToSlash(rel))
		return nil
	})
	sort.Strings(files)
	return files, err
}
func judgeArchiveCore(name string) bool {
	base := filepath.Base(name)
	return name == "_prd.md" || name == "_techspec.md" || name == "_tasks.md" || name == "_authorization.md" || (strings.HasPrefix(base, "task_") && strings.HasSuffix(base, ".md")) || (strings.HasPrefix(base, "qa-report-") && strings.HasSuffix(base, ".md")) || name == "references/_index.md"
}
