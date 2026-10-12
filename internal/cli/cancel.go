package cli

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"roundfix/internal/app"
	roundconfig "roundfix/internal/config"
	"roundfix/internal/spec"
	"roundfix/internal/speccheck"
	"roundfix/internal/store"
)

var cancelUsage = `Usage:
  roundfix cancel <slug> --reason <text>

Cancels an active Spec that will not be delivered. Writes one Archive Record
with disposition cancelled, the reason, the date and the Tasks completed and
pending, returns each adopted Backlog Entry to docs/backlog/ (status: open)
and each adopted Finding to docs/findings/ (status: pending), and removes the
committed Spec folder. The removed bytes stay in Git at the record's
source_revision. cancel creates no Run, never changes a Task status, and
never commits or pushes.

Refuses while an Active Run or an unfinished Delivery Queue item holds the
Spec, while an active Spec requires it, and for a Spec that carries a
supersession.

Options:
  --reason  Why the Spec is cancelled, recorded in the Archive Record

Exit codes:
  0  cancelled
  1  cancel write failed
  2  Preflight Validation failed
`

func runCancelCommand(ctx context.Context, args []string, stdout, stderr io.Writer, environment commandEnvironment) int {
	if commandWantsHelp(args) {
		fmt.Fprint(stdout, cancelUsage)
		return exitOK
	}
	refuse := func(err error) int {
		printPreflightFailure("cancel", err, stderr)
		return exitPreflight
	}
	slug, reason, err := parseCancelCommand(args)
	if err != nil {
		return refuse(err)
	}
	loaded, err := loadCommandConfig(environment, stderr)
	if err != nil {
		return refuse(err)
	}
	if loaded.GitRoot == "" {
		return refuse(errors.New("cancel requires a git repository working tree"))
	}
	root, err := roundconfig.ResolveSpecsRoot(loaded, loaded.GitRoot)
	if err != nil {
		return refuse(err)
	}
	folder := filepath.Join(root.Path, slug)
	info, err := os.Lstat(folder)
	if errors.Is(err, os.ErrNotExist) || (err == nil && !info.IsDir()) {
		return refuse(fmt.Errorf("unknown active Spec slug %q", slug))
	}
	if err != nil {
		return refuse(fmt.Errorf("inspect active Spec: %w", err))
	}
	pins, err := speccheck.ActiveSpecPathPins(loaded.GitRoot, root.Path, slug)
	if err != nil {
		return refuse(err)
	}
	if len(pins) > 0 {
		var locations []string
		for _, pin := range pins {
			locations = append(locations, fmt.Sprintf("%s:%d", pin.Path, pin.Line))
		}
		rel, err := filepathRelSlash(loaded.GitRoot, folder)
		if err != nil {
			return refuse(err)
		}
		return refuse(fmt.Errorf("Spec %q cannot cancel while another file names its active directory %s/: %s", slug, rel, strings.Join(locations, ", ")))
	}
	revisionRoot := loaded.GitRoot
	if root.External {
		revisionRoot = root.Path
	}
	repositoryRoot, err := gitOutput(ctx, revisionRoot, "rev-parse", "--show-toplevel")
	if err != nil {
		return refuse(err)
	}
	changes, err := gitOutput(ctx, revisionRoot, "status", "--porcelain", "--untracked-files=all", "--", folder)
	if err != nil {
		return refuse(err)
	}
	if strings.TrimSpace(changes) != "" {
		rel, err := filepathRelSlash(repositoryRoot, folder)
		if err != nil {
			return refuse(err)
		}
		return refuse(fmt.Errorf("Spec %q has changes not committed at HEAD under %s; commit them before cancel so the Archive Record's source_revision holds the Spec", slug, rel))
	}
	parked, err := checkCancelWork(ctx, loaded.HomeDir, loaded.GitRoot, slug, stderr)
	if err != nil {
		return refuse(err)
	}
	if err := checkCancelRequires(root.Path, slug); err != nil {
		return refuse(err)
	}
	revision, err := gitOutput(ctx, revisionRoot, "rev-parse", "HEAD")
	if err != nil {
		return refuse(fmt.Errorf("resolve cancellation revision from HEAD: %w", err))
	}
	// HistoryCitations reads tracked working-tree files, including the Spec's
	// own files. Collect advice against the same relink plan Cancel will apply.
	citations, err := cancelCitations(ctx, repositoryRoot, root, folder)
	if err != nil {
		return refuse(err)
	}
	if err := ctx.Err(); err != nil {
		return refuse(err)
	}
	folderRel, err := filepathRelSlash(repositoryRoot, folder)
	if err != nil {
		return refuse(err)
	}
	recordRel, err := filepathRelSlash(loaded.GitRoot, spec.ArchiveRecordPath(spec.ArchiveSpecRoot(root.Path, root.BuiltInRoot), slug))
	if err != nil {
		return refuse(err)
	}
	result, err := spec.Cancel(spec.CancelRequest{
		SpecsRoot: root.Path, BuiltInRoot: root.BuiltInRoot, Slug: slug,
		Reason: reason, SourceRevision: revision, RepositoryRoot: repositoryRoot,
	})
	if err != nil {
		if cancelWriteError(err) {
			fmt.Fprintf(stderr, "%s: cancel write failed: %v\n", app.Name, err)
			return exitRunFailed
		}
		return refuse(err)
	}
	fmt.Fprintf(stdout, "cancelled %s -> %s; %d completed Task(s), %d pending; removed %d file(s) (%d bytes) kept in Git at %.12s\n", slug, recordRel, len(result.Completed), len(result.Pending), result.RemovedFiles, result.RemovedBytes, revision)
	for _, source := range result.Returned {
		fmt.Fprintf(stdout, "returned %s as %s\n", source.To, source.Status)
	}
	if parked {
		fmt.Fprintf(stderr, "%s: warning: Delivery Queue item %s stays parked; a Delivery Retry of a cancelled Spec is refused, so close its branch and pull request by hand\n", app.Name, slug)
	}
	removed := make(map[string]bool, len(result.RemovedPaths)+1)
	for _, path := range result.RemovedPaths {
		removed[path] = true
	}
	removed[folderRel+"/"] = true
	warned := make(map[string]bool)
	for _, citation := range citations {
		location := fmt.Sprintf("%s:%d", citation.Path, citation.Line)
		if !removed[citation.Target] || warned[location] {
			continue
		}
		fmt.Fprintf(stderr, "%s: warning: %s cites %s, which this cancellation removed\n", app.Name, location, citation.Target)
		warned[location] = true
	}
	return exitOK
}

func parseCancelCommand(args []string) (string, string, error) {
	if len(args) == 0 || strings.TrimSpace(args[0]) == "" || strings.HasPrefix(args[0], "-") {
		return "", "", validationError{message: "missing required Spec slug; pass roundfix cancel <slug> --reason <text>"}
	}
	slug := strings.TrimSpace(args[0])
	flags := flag.NewFlagSet("cancel", flag.ContinueOnError)
	flags.SetOutput(io.Discard)
	reason := flags.String("reason", "", "Why the Spec is cancelled")
	parseErr := flags.Parse(args[1:])
	normalized := strings.Join(strings.Fields(*reason), " ")
	if normalized == "" {
		return "", "", validationError{message: "missing required --reason; pass --reason <text>"}
	}
	if parseErr != nil {
		return "", "", validationError{message: parseErr.Error()}
	}
	if remaining := flags.Args(); len(remaining) > 0 {
		return "", "", validationError{message: fmt.Sprintf("unexpected argument %q", remaining[0])}
	}
	if filepath.Base(slug) != slug || slug == "." || slug == ".." || strings.Contains(slug, "\\") {
		return "", "", validationError{message: fmt.Sprintf("unsafe Spec slug %q", slug)}
	}
	return slug, normalized, nil
}

func checkCancelWork(ctx context.Context, home, repo, slug string, stderr io.Writer) (bool, error) {
	if _, err := os.Stat(store.DatabasePath(home)); errors.Is(err, os.ErrNotExist) {
		return false, nil
	} else if err != nil {
		return false, fmt.Errorf("inspect Run Database before cancel: %w", err)
	}
	db, err := store.Open(ctx, home)
	if err != nil {
		return false, fmt.Errorf("open Run Database before cancel: %w", err)
	}
	defer func() { _ = db.Close() }()
	active, found, err := db.ActiveSpecRun(ctx, repo, slug)
	if err != nil {
		return false, fmt.Errorf("check Active Run before cancel: %w", err)
	}
	if found {
		_, reclaimed, err := reclaimOrphanedActiveRun(ctx, db, active, stderr)
		if err != nil {
			return false, err
		}
		if !reclaimed {
			return false, fmt.Errorf("Active Run %s holds Spec %q; stop it with: roundfix stop %s", active.ID, slug, active.ID)
		}
	}
	queue, _, err := db.DeliveryQueue(ctx, repo)
	if err != nil {
		return false, fmt.Errorf("read Delivery Queue before cancel: %w", err)
	}
	for _, item := range queue.Items {
		if item.SpecSlug != slug {
			continue
		}
		switch item.Stage {
		case store.DeliveryStageMerged:
			return false, nil
		case store.DeliveryStageParked:
			return true, nil
		default:
			return false, fmt.Errorf("Spec %q cannot be cancelled while its Delivery Queue item is at stage %s; stop the queue with: roundfix deliver stop", slug, item.Stage)
		}
	}
	return false, nil
}

func checkCancelRequires(root, slug string) error {
	entries, err := os.ReadDir(root)
	if err != nil {
		return fmt.Errorf("read active Specs before cancel: %w", err)
	}
	for _, entry := range entries {
		if !entry.IsDir() || entry.Name() == slug || strings.HasPrefix(entry.Name(), "_") {
			continue
		}
		manifest, err := os.ReadFile(filepath.Join(root, entry.Name(), "_tasks.md"))
		if errors.Is(err, os.ErrNotExist) {
			continue
		}
		if err != nil {
			return fmt.Errorf("read active Spec %s Task Graph: %w", entry.Name(), err)
		}
		requires, err := spec.ParseRequiredSpecs(manifest, entry.Name())
		if err != nil {
			return fmt.Errorf("read active Spec %s prerequisites: %w", entry.Name(), err)
		}
		for _, required := range requires {
			if required == slug {
				return fmt.Errorf("Spec %q cannot be cancelled while active Spec %s requires it; remove it from that Spec's requires first", slug, entry.Name())
			}
		}
	}
	return nil
}

func cancelCitations(ctx context.Context, repo string, root roundconfig.SpecsRoot, folder string) ([]spec.HistoryCitation, error) {
	folderRel, err := filepathRelSlash(repo, folder)
	if err != nil {
		return nil, err
	}
	removed := []string{folderRel + "/"}
	err = filepath.WalkDir(folder, func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.IsDir() {
			return nil
		}
		rel, err := filepathRelSlash(repo, path)
		if err != nil {
			return err
		}
		removed = append(removed, rel)
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("inspect cancellation citation targets: %w", err)
	}
	citations, err := spec.HistoryCitations(ctx, repo, removed)
	if err != nil {
		return nil, err
	}
	specRootRel, err := filepathRelSlash(repo, root.Path)
	if err != nil {
		return nil, err
	}
	archiveRootRel, err := filepathRelSlash(repo, spec.ArchiveSpecRoot(root.Path, root.BuiltInRoot))
	if err != nil {
		return nil, err
	}
	recordRel, err := filepathRelSlash(repo, spec.ArchiveRecordPath(spec.ArchiveSpecRoot(root.Path, root.BuiltInRoot), filepath.Base(folder)))
	if err != nil {
		return nil, err
	}
	plan, err := spec.PlanArchiveRelink(ctx, repo, specRootRel, archiveRootRel, folderRel, recordRel, nil)
	if err != nil {
		return nil, fmt.Errorf("plan cancellation citation advice: %w", err)
	}
	// A relink can change a line while leaving a code example on that same line
	// untouched. Use the shared plan to identify the citations that will remain.
	remaining := make(map[string]bool)
	for _, citation := range plan.Unrelinked {
		remaining[fmt.Sprintf("%s:%d", citation.Path, citation.Line)] = true
	}
	var live []spec.HistoryCitation
	for _, citation := range citations {
		if remaining[fmt.Sprintf("%s:%d", citation.Path, citation.Line)] {
			live = append(live, citation)
		}
	}
	return live, nil
}

// Cancel returns ordinary wrapped errors. Its mutation phase has these explicit
// operation labels; all earlier planning and validation failures are preflight.
func cancelWriteError(err error) bool {
	for _, operation := range []string{"create cancellation destination directory:", "create cancellation destination ", "write cancellation destination ", "write cancellation relink ", "remove cancelled Spec;"} {
		if strings.HasPrefix(err.Error(), operation) {
			return true
		}
	}
	return false
}
