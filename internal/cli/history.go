package cli

import (
	"context"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"roundfix/internal/config"
	"roundfix/internal/judge"
	"roundfix/internal/spec"
)

const historyUsage = `Usage:
  roundfix history sanitize [--batch <n>] [--advise] [--apply] [--promote <path> ...]

Plans pending legacy folders, findings, backlog, ADRs, reviews, handoffs and superseded records in that
order. Refused Units are listed with their reason, left in place and not counted
toward --batch. Plans write nothing in the repository. Advice is optional and advisory.

Options:
  --batch <n>       Select the next positive number of pending units
  --advise          Add Jev advice for the selected folders (requires --batch)
  --apply           Apply one batch on a clean tree covered by history-full
  --promote <path>  Copy a repository-relative file in the batch to docs/references/

Apply requires an annotated history-full tag at or before HEAD. Every removed
path must be held by that tag or a commit the default branch reaches. This command never commits, tags or pushes.
Exit codes: 0 plan, applied or nothing pending; 1 write failure; 2 preflight refusal.
`

type historyRequest struct {
	batch                   int
	batchSet, advise, apply bool
	promote                 []string
}

type historyUnit struct {
	folder     bool
	name       string
	files      []string // repository-relative paths removed or rewritten
	conversion *spec.LegacyConversion
	kind       *spec.HistoryKindPlan
	refusal    error
}

func parseHistory(args []string) (historyRequest, error) {
	var req historyRequest
	f := flag.NewFlagSet("history sanitize", flag.ContinueOnError)
	f.SetOutput(io.Discard)
	f.IntVar(&req.batch, "batch", 0, "Next units")
	f.BoolVar(&req.advise, "advise", false, "Advice")
	f.BoolVar(&req.apply, "apply", false, "Apply")
	f.Func("promote", "Repository-relative file", func(s string) error { req.promote = append(req.promote, s); return nil })
	if err := f.Parse(args); err != nil {
		return req, err
	}
	f.Visit(func(v *flag.Flag) {
		if v.Name == "batch" {
			req.batchSet = true
		}
	})
	switch {
	case len(f.Args()) > 0:
		return req, fmt.Errorf("unexpected argument %q", f.Args()[0])
	case req.batchSet && req.batch < 1:
		return req, fmt.Errorf("--batch must be a positive integer")
	case req.apply && !req.batchSet:
		return req, fmt.Errorf("--apply requires --batch <n>")
	case req.advise && (!req.batchSet || req.apply):
		return req, fmt.Errorf("--advise requires --batch and cannot be combined with --apply")
	case len(req.promote) > 0 && !req.apply:
		return req, fmt.Errorf("--promote requires --apply")
	}
	return req, nil
}

// Inventory preserves unit errors so planning can continue past a refused unit.
func historyInventory(ctx context.Context, root, archive string) ([]historyUnit, error) {
	slugs, err := spec.LegacyArchiveFolders(archive)
	if err != nil {
		return nil, err
	}
	var units []historyUnit
	var removed []string
	archiveRel, _ := filepathRelSlash(root, archive)
	specRootRel := "docs/specs"
	if archiveRel != spec.ArchiveDir(spec.ArchiveKindSpec) {
		specRootRel = filepath.ToSlash(filepath.Dir(archiveRel))
	}
	for _, slug := range slugs {
		units = append(units, historyUnit{name: slug, folder: true})
		names, err := archivePlanFiles(filepath.Join(archive, slug))
		if err == nil {
			for _, name := range names {
				removed = append(removed, archiveRel+"/"+slug+"/"+name)
			}
		}
	}
	for _, kind := range []spec.ArchiveKind{spec.ArchiveKindFinding, spec.ArchiveKindBacklog, spec.ArchiveKindADR, spec.ArchiveKindReview, spec.ArchiveKindHandoff} {
		plan, err := spec.PlanHistoryKind(ctx, root, kind, specRootRel)
		if len(plan.Files) > 0 || err != nil {
			units = append(units, historyUnit{name: string(kind), files: plan.Files, kind: &plan, refusal: err})
			if err == nil {
				removed = append(removed, plan.Files...)
			}
		}
	}
	plan, err := spec.PlanSupersededRecords(ctx, root, archiveRel, removed)
	if len(plan.Files) > 0 || err != nil {
		units = append(units, historyUnit{name: "superseded", files: plan.Files, kind: &plan, refusal: err})
	}
	return units, nil
}

func historyTagCheck(ctx context.Context, root string) error {
	const needs = "history sanitize --apply needs the annotated tag history-full at or before HEAD; tag the last commit before the first batch with git tag -a history-full"
	typ, err := gitOutput(ctx, root, "cat-file", "-t", "refs/tags/history-full")
	if err != nil || typ != "tag" {
		return fmt.Errorf("%s", needs)
	}
	if _, err := gitOutput(ctx, root, "merge-base", "--is-ancestor", "refs/tags/history-full^{commit}", "HEAD"); err != nil {
		return fmt.Errorf("%s", needs)
	}
	return nil
}

func historyTagCoverage(ctx context.Context, root string, units []historyUnit) error {
	tree, err := gitOutput(ctx, root, "ls-tree", "-r", "--name-only", "-z", "refs/tags/history-full")
	if err != nil {
		return err
	}
	paths := map[string]bool{}
	for _, p := range strings.Split(tree, "\x00") {
		paths[p] = true
	}
	for _, u := range units {
		for _, p := range u.files {
			if !paths[p] {
				locator, err := spec.ResolveGitLocator(ctx, root, p, "history-full", "")
				if err != nil {
					return err
				}
				if locator.Commit == "" {
					return fmt.Errorf("history-full does not hold %s and no commit the default branch reaches holds it", p)
				}
			}
		}
	}
	return nil
}

// planHistoryUnit shares conversion planning without selection or command output.
func planHistoryUnit(ctx context.Context, repo, specRel, archive, archiveRel, revision string, u historyUnit) (historyUnit, error) {
	if u.folder {
		names, err := archivePlanFiles(filepath.Join(archive, u.name))
		u.refusal = err
		for _, name := range names {
			u.files = append(u.files, archiveRel+"/"+u.name+"/"+name)
		}
	}
	if u.refusal == nil {
		if u.folder {
			delivery, err := spec.FindLegacyDelivery(ctx, repo, specRel, archiveRel, u.name)
			if err != nil {
				return u, err
			}
			created, err := spec.FindLegacyCreation(ctx, repo, specRel, archiveRel, u.name)
			if err != nil {
				return u, err
			}
			c, err := spec.PlanLegacyConversion(spec.LegacyConversionRequest{RepositoryRoot: repo, ArchiveRoot: archive, Slug: u.name, SourceRevision: revision, Created: created, Delivery: delivery})
			u.refusal = err
			u.conversion = &c
		} else if u.name != "superseded" {
			k, err := spec.PlanHistoryKind(ctx, repo, spec.ArchiveKind(u.name), specRel)
			u.refusal = err
			u.kind = &k
			u.files = k.Files
		}
	}
	return u, nil
}

func runHistoryCommand(ctx context.Context, args []string, stdout, stderr io.Writer, env commandEnvironment) int {
	fail := func(err error) int { printPreflightFailure("history sanitize", err, stderr); return exitPreflight }
	if len(args) == 0 || commandWantsHelp(args) {
		fmt.Fprint(stdout, historyUsage)
		return exitOK
	}
	if args[0] != "sanitize" {
		return fail(fmt.Errorf("unknown history subcommand %q", args[0]))
	}
	req, err := parseHistory(args[1:])
	if err != nil {
		return fail(err)
	}
	loaded, err := loadCommandConfig(env, stderr)
	if err != nil {
		return fail(err)
	}
	if loaded.GitRoot == "" {
		return fail(fmt.Errorf("history sanitize requires a git repository"))
	}
	root, err := config.ResolveSpecsRoot(loaded, loaded.GitRoot)
	if err != nil {
		return fail(err)
	}
	if root.External {
		return fail(fmt.Errorf("history sanitize refuses an external Spec Root"))
	}
	repo := loaded.GitRoot
	if req.apply {
		status, err := gitOutput(ctx, repo, "status", "--porcelain", "--untracked-files=all")
		if err != nil {
			return fail(err)
		}
		if status != "" {
			return fail(fmt.Errorf("history sanitize --apply requires a clean working tree, including untracked files"))
		}
	}
	revision, err := gitOutput(ctx, repo, "rev-parse", "HEAD")
	if err != nil {
		return fail(err)
	}
	archive := spec.ArchiveSpecRoot(root.Path, root.BuiltInRoot)
	units, err := historyInventory(ctx, repo, archive)
	if err != nil {
		return fail(err)
	}
	total := len(units)
	if len(units) == 0 {
		if len(req.promote) > 0 {
			return fail(fmt.Errorf("promotion is outside the batch"))
		}
		fmt.Fprintln(stdout, "history sanitize plan: 0 unit(s) pending; nothing pending")
		return exitOK
	}
	if req.apply {
		if err := historyTagCheck(ctx, repo); err != nil {
			return fail(err)
		}
	}
	specRel, _ := filepathRelSlash(repo, root.Path)
	archiveRel, _ := filepathRelSlash(repo, archive)
	var selected, refused []historyUnit
	for _, u := range units {
		if req.batchSet && len(selected) == req.batch {
			break
		}
		if u.name == "superseded" && !u.folder {
			var removed []string
			for _, selectedUnit := range selected {
				removed = append(removed, selectedUnit.files...)
			}
			k, e := spec.PlanSupersededRecords(ctx, repo, archiveRel, removed)
			u.kind, u.files, u.refusal = &k, k.Files, e
			if len(k.Files) == 0 && e == nil {
				total--
				continue
			}
		}
		u, err = planHistoryUnit(ctx, repo, specRel, archive, archiveRel, revision, u)
		if err != nil {
			return fail(err)
		}
		if u.refusal != nil {
			if u.folder {
				u.name = archiveRel + "/" + u.name
			}
			refused = append(refused, u)
			continue
		}
		selected = append(selected, u)
	}
	promotions := map[string][]string{}
	destinations := map[string]bool{}
	for _, p := range req.promote {
		if filepath.IsAbs(p) || filepath.ToSlash(filepath.Clean(p)) != p || strings.Contains(p, "\\") {
			return fail(fmt.Errorf("promotion must be repository-relative: %q", p))
		}
		for _, u := range refused {
			if u.folder && strings.HasPrefix(p, u.name+"/") {
				return fail(fmt.Errorf("promotion %q is in refused unit %s: %v", p, u.name, u.refusal))
			}
		}
		found := false
		for _, u := range selected {
			if u.folder && strings.HasPrefix(p, u.conversion.Folder+"/") {
				dest := filepath.Base(p)
				if destinations[dest] {
					return fail(fmt.Errorf("duplicate promotion destination %q", dest))
				}
				destinations[dest] = true
				promotions[u.name] = append(promotions[u.name], p)
				found = true
				break
			}
		}
		if !found {
			return fail(fmt.Errorf("promotion %q is outside the batch", p))
		}
	}
	for i := range selected {
		u := &selected[i]
		if !u.folder || len(promotions[u.name]) == 0 {
			continue
		}
		c := u.conversion
		converted, err := spec.PlanLegacyConversion(spec.LegacyConversionRequest{RepositoryRoot: repo, ArchiveRoot: archive, Slug: u.name, SourceRevision: revision, Created: c.Record.Created, Delivery: spec.LegacyDelivery{Commit: c.Record.DeliveryCommit, PullRequest: c.Record.PullRequest, Date: c.Record.Archived}, Promote: promotions[u.name]})
		if err != nil {
			return fail(err)
		}
		u.conversion = &converted
	}
	if req.apply {
		if err := historyTagCoverage(ctx, repo, selected); err != nil {
			return fail(err)
		}
		printHistoryRefused(refused, stdout)
		if len(selected) == 0 && len(refused) > 0 {
			return fail(fmt.Errorf("history sanitize --apply found no convertible unit; %d unit(s) refused", len(refused)))
		}
	}
	if req.apply {
		rewrite, err := planHistoryRewrite(ctx, repo, specRel, archiveRel, selected)
		if err != nil {
			return fail(err)
		}
		return applyHistoryUnits(ctx, repo, revision, selected, rewrite, len(refused), total-len(selected), stdout, stderr)
	}
	return printHistoryPlan(ctx, req, loaded, env, selected, refused, stdout, stderr)
}

func applyHistoryUnits(ctx context.Context, repo, revision string, units []historyUnit, rewrite spec.HistoryRewritePlan, refused, remaining int, stdout, stderr io.Writer) int {
	var records, removed, promoted int
	var bytes int64
	for _, u := range units {
		err := ctx.Err()
		if err == nil {
			if c := u.conversion; c != nil {
				err = spec.ApplyLegacyConversion(repo, *c)
				records++
				removed += len(c.Files)
				bytes += c.Bytes
				promoted += len(c.Promoted)
			} else {
				err = spec.ApplyHistoryKind(repo, *u.kind)
				removed += len(u.files)
				bytes += u.kind.BytesBefore
			}
		}
		if err != nil {
			fmt.Fprintf(stderr, "history sanitize write failed for %s: %v; restore with git restore and git clean\n", u.name, err)
			return exitRunFailed
		}
	}
	rewritten, rewrittenFiles, err := rewrite.Apply(ctx, repo)
	if err != nil {
		fmt.Fprintf(stderr, "history sanitize rewrite failed: %v; restore with git restore and git clean\n", err)
		return exitRunFailed
	}
	fmt.Fprintf(stdout, "history sanitize applied %d unit(s): wrote %d Archive Record(s), removed %d file(s) (%d bytes) kept in Git at %.12s and tag history-full; rewrote %d citation(s) in %d file(s); promoted %d file(s) to docs/references/%s; %d unit(s) remain\n", len(units), records, removed, bytes, revision, rewritten, rewrittenFiles, promoted, historyRefusalSuffix(refused), remaining)
	return exitOK
}

func printHistoryPlan(ctx context.Context, req historyRequest, loaded config.Loaded, env commandEnvironment, units, refused []historyUnit, stdout, stderr io.Writer) int {
	var files int
	var bytes int64
	var removed []string
	for _, u := range units {
		files += len(u.files)
		removed = append(removed, u.files...)
		if u.conversion != nil {
			bytes += u.conversion.Bytes
		} else {
			bytes += u.kind.BytesBefore
		}
	}
	archiveRel := spec.ArchiveDir(spec.ArchiveKindSpec)
	specs, err := config.ResolveSpecsRoot(loaded, loaded.GitRoot)
	if err != nil {
		printPreflightFailure("history sanitize", err, stderr)
		return exitPreflight
	}
	specRel, _ := filepathRelSlash(loaded.GitRoot, specs.Path)
	archiveRel, _ = filepathRelSlash(loaded.GitRoot, spec.ArchiveSpecRoot(specs.Path, specs.BuiltInRoot))
	rewrite, err := planHistoryRewrite(ctx, loaded.GitRoot, specRel, archiveRel, units)
	if err != nil {
		printPreflightFailure("history sanitize", err, stderr)
		return exitPreflight
	}
	citations := rewrite.Citations
	fmt.Fprintf(stdout, "history sanitize plan: %d unit(s) pending; %d file(s) (%d bytes) leave docs/history%s\n", len(units), files, bytes, historyRefusalSuffix(len(refused)))
	for _, u := range units {
		if c := u.conversion; c != nil {
			delivery := "unknown"
			if c.Record.DeliveryCommit != "" {
				delivery = fmt.Sprintf("%.12s", c.Record.DeliveryCommit)
				if c.Record.PullRequest != "" {
					delivery += " #" + c.Record.PullRequest
				}
			}
			fmt.Fprintf(stdout, "folder %s: removes %d file(s) (%d bytes) and writes %s (%d bytes, %s, delivery %s)\n", c.Folder, len(c.Files), c.Bytes, c.RecordPath, len(c.Rendered), c.Record.Disposition, delivery)
			for _, tolerance := range c.Tolerated {
				fmt.Fprintf(stdout, "tolerates %s: %s\n", c.Folder, tolerance)
			}
			var candidates []string
			for _, p := range c.Files {
				if !judgeArchiveCore(p) && !strings.HasPrefix(p, "qa/evidence/") {
					candidates = append(candidates, p)
				}
			}
			advice := map[string]string{}
			if req.advise {
				q, err := judge.Load()
				if err != nil {
					for _, p := range candidates {
						advice[p] = "no advice (" + err.Error() + ")"
					}
				} else {
					keys := map[string]string{}
					for _, key := range q.KeyVariables() {
						for _, entry := range env.environ {
							n, v, _ := strings.Cut(entry, "=")
							if n == key {
								keys[n] = v
							}
						}
					}
					report, err := judge.AdviseArchive(ctx, q.WithMonthlyCeiling(loaded.Config.Jev.MonthlyCeilingUSD), judge.ArchiveAdviceRequest{RepoRoot: loaded.GitRoot, SpecDir: filepath.Join(loaded.GitRoot, c.Folder), Spec: c.Slug, Files: candidates, Keys: keys, HomeDir: loaded.HomeDir, Transport: env.dependencies.judgeTransport, Now: env.dependencies.judgeNow})
					if err != nil {
						for _, p := range candidates {
							advice[p] = "no advice (" + err.Error() + ")"
						}
					} else {
						for _, a := range report.Advice {
							detail := "no advice"
							if a.Choice != nil {
								detail = *a.Choice
								if probability, ok := a.Probabilities[*a.Choice]; ok {
									detail += fmt.Sprintf(" (p=%.2f)", probability)
								}
							} else if a.Reason != nil {
								detail += " (" + *a.Reason + ")"
							}
							advice[a.File] = detail
						}
					}
				}
			}
			for _, p := range candidates {
				info, err := os.Stat(filepath.Join(loaded.GitRoot, c.Folder, p))
				if err != nil {
					printPreflightFailure("history sanitize", err, stderr)
					return exitPreflight
				}
				suffix := ""
				if req.advise {
					suffix = ": " + advice[p]
				}
				fmt.Fprintf(stdout, "candidate %s/%s %d bytes%s\n", c.Folder, p, info.Size(), suffix)
			}
		} else {
			k := u.kind
			fmt.Fprintf(stdout, "%s: removes %d file(s) (%d bytes)\n", u.name, len(k.Files), k.BytesBefore)

		}
	}
	printHistoryRefused(refused, stdout)
	for _, c := range citations {
		fmt.Fprintf(stdout, "cites %s:%d names %s", c.Path, c.Line, c.Target)
		if c.Rewrite != "" {
			fmt.Fprintf(stdout, "; rewrites to %s", c.Rewrite)
		}
		fmt.Fprintln(stdout)
	}
	fmt.Fprintln(stdout, "apply with: roundfix history sanitize --apply --batch <n> (needs the annotated tag history-full at or before HEAD)")
	return exitOK
}

func historyRefusalSuffix(count int) string {
	if count == 0 {
		return ""
	}
	return fmt.Sprintf("; %d unit(s) refused", count)
}

func printHistoryRefused(units []historyUnit, stdout io.Writer) {
	for _, u := range units {
		fmt.Fprintf(stdout, "refused %s: %s\n", u.name, historyRefusalLine(u.refusal))
	}
}

func historyRefusalLine(err error) string {
	return strings.Join(strings.Fields(err.Error()), " ")
}

// Both commands share the same removal rewrite, frozen before applying units.
func planHistoryRewrite(ctx context.Context, repo, specRel, archiveRel string, units []historyUnit) (spec.HistoryRewritePlan, error) {
	var removed []string
	moves := map[string]string{}
	for _, u := range units {
		removed = append(removed, u.files...)
		if c := u.conversion; c != nil {
			promotions := map[string]string{}
			for _, destination := range c.Promoted {
				for _, file := range u.files {
					if filepath.Base(file) == filepath.Base(destination) {
						promotions[file] = destination
					}
				}
			}
			for source, destination := range spec.ArchiveCitationMoves(c.Folder, c.RecordPath, u.files, promotions) {
				moves[source] = destination
			}
		}
	}
	return spec.PlanHistoryRewrite(ctx, repo, specRel, archiveRel, removed, moves)
}
