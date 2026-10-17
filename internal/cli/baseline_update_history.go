package cli

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"roundfix/internal/config"
	"roundfix/internal/spec"
)

type baselineUpdateHistory struct {
	Status         string                          `json:"status"`
	Message        string                          `json:"message,omitempty"`
	Revision       string                          `json:"revision,omitempty"`
	Tag            *baselineUpdateHistoryTag       `json:"tag,omitempty"`
	Units          []baselineUpdateHistoryUnit     `json:"units"`
	Refused        []baselineUpdateHistoryRefusal  `json:"refused"`
	Applied        *baselineUpdateHistoryApplied   `json:"applied,omitempty"`
	Citations      []baselineUpdateHistoryCitation `json:"citations"`
	StaleCitations []baselineUpdateHistoryCitation `json:"staleCitations"`
}
type baselineUpdateHistoryCitation struct {
	Path       string `json:"path"`
	Line       int    `json:"line"`
	Target     string `json:"target"`
	Live       bool   `json:"live"`
	Offer      string `json:"offer,omitempty"`
	PromotedTo string `json:"promotedTo,omitempty"`
	Rewrite    string `json:"rewrite,omitempty"`
}

type baselineUpdateHistoryTag struct {
	Name        string `json:"name"`
	Action      string `json:"action"`
	Commit      string `json:"commit"`
	PushCommand string `json:"pushCommand,omitempty"`
}
type baselineUpdateHistoryUnit struct {
	Unit        string   `json:"unit"`
	Action      string   `json:"action"`
	Paths       []string `json:"paths"`
	Record      string   `json:"record,omitempty"`
	Disposition string   `json:"disposition,omitempty"`
	Files       int      `json:"files"`
	Bytes       int64    `json:"bytes"`
	Links       int      `json:"links,omitempty"`
}
type baselineUpdateHistoryRefusal struct {
	Unit   string `json:"unit"`
	Reason string `json:"reason"`
}
type baselineUpdateHistoryApplied struct {
	Units              int   `json:"units"`
	Records            int   `json:"records"`
	Reduced            int   `json:"reduced"`
	Removed            int   `json:"removed"`
	Bytes              int64 `json:"bytes"`
	RewrittenCitations int   `json:"rewrittenCitations"`
	RewrittenFiles     int   `json:"rewrittenFiles"`
}
type baselineUpdateHistoryPlan struct {
	Report         baselineUpdateHistory
	selected       []historyUnit
	rewrite        historyCitationRewritePlan
	root           string
	historyRoot    string
	archiveRoot    string
	specRootRel    string
	archiveRootRel string
}

func newBaselineUpdateHistory(status string) baselineUpdateHistory {
	return baselineUpdateHistory{Status: status, Units: []baselineUpdateHistoryUnit{}, Refused: []baselineUpdateHistoryRefusal{}, Citations: []baselineUpdateHistoryCitation{}, StaleCitations: []baselineUpdateHistoryCitation{}}
}

func planBaselineUpdateHistory(ctx context.Context, repo string, environment commandEnvironment) baselineUpdateHistoryPlan {
	p := baselineUpdateHistoryPlan{Report: newBaselineUpdateHistory("current")}
	blocked := func(err error) baselineUpdateHistoryPlan {
		p.Report.Status = "blocked"
		p.Report.Message = historyRefusalLine(err)
		p.Report.Units = []baselineUpdateHistoryUnit{}
		p.selected = nil
		p.Report.Citations = []baselineUpdateHistoryCitation{}
		if p.Report.Tag != nil && p.Report.Tag.Action == "create" {
			p.Report.Tag = nil
		}
		return p
	}
	root, err := gitOutput(ctx, repo, "rev-parse", "--show-toplevel")
	if err != nil {
		return blocked(err)
	}
	p.root = root
	environment.workDir, environment.workDirErr = root, nil
	loaded, err := loadCommandConfig(environment, io.Discard)
	if err != nil {
		return blocked(err)
	}
	specs, err := config.ResolveSpecsRoot(loaded, root)
	if err != nil {
		return blocked(err)
	}
	if specs.External {
		return blocked(fmt.Errorf("baseline update history refuses an external Spec Root"))
	}
	archive := spec.ArchiveSpecRoot(specs.Path, specs.BuiltInRoot)
	specRel, err := filepathRelSlash(root, specs.Path)
	if err != nil {
		return blocked(err)
	}
	archiveRel, err := filepathRelSlash(root, archive)
	if err != nil {
		return blocked(err)
	}
	p.historyRoot = filepath.ToSlash(filepath.Dir(archiveRel))
	p.archiveRoot, p.specRootRel, p.archiveRootRel = archive, specRel, archiveRel
	units, err := historyInventory(ctx, root, archive)
	if err != nil {
		return blocked(err)
	}
	revision, err := gitOutput(ctx, root, "rev-parse", "HEAD")
	if err != nil {
		return blocked(err)
	}
	p.Report.Revision = revision
	tagRefs, err := gitOutput(ctx, root, "for-each-ref", "--format=%(refname) %(objecttype)", "refs/tags/history-full")
	if err != nil {
		return blocked(err)
	}
	typ := ""
	for _, ref := range strings.Split(tagRefs, "\n") {
		if kind, ok := strings.CutPrefix(ref, "refs/tags/history-full "); ok {
			typ = kind
			break
		}
	}
	var tagPaths map[string]bool
	if typ != "" {
		if typ != "tag" {
			return blocked(fmt.Errorf("history-full is not annotated"))
		}
		if _, err = gitOutput(ctx, root, "merge-base", "--is-ancestor", "refs/tags/history-full^{commit}", revision); err != nil {
			return blocked(fmt.Errorf("history-full is not an ancestor of HEAD: %w", err))
		}
		commit, err := gitOutput(ctx, root, "rev-parse", "refs/tags/history-full^{commit}")
		if err != nil {
			return blocked(err)
		}
		p.Report.Tag = &baselineUpdateHistoryTag{Name: "history-full", Action: "present", Commit: commit}
		tagPaths, err = baselineHistoryTree(ctx, root, "refs/tags/history-full")
		if err != nil {
			return blocked(err)
		}
	}
	headPaths, err := baselineHistoryTree(ctx, root, revision)
	if err != nil {
		return blocked(err)
	}
	for _, u := range units {
		if u.name == "superseded" && !u.folder {
			var removed []string
			for _, selected := range p.selected {
				removed = append(removed, selected.files...)
			}
			k, e := spec.PlanSupersededRecords(ctx, root, archiveRel, removed)
			u.kind, u.files, u.refusal = &k, k.Files, e
			if len(k.Files) == 0 && e == nil {
				continue
			}
		}
		if u.links != nil {
			u, err = planRecordLinksUnit(ctx, root, archiveRel, p.selected)
			if err != nil {
				return blocked(err)
			}
			if len(u.links.Files) == 0 {
				continue
			}
		}
		u, err = planHistoryUnit(ctx, root, specRel, archive, archiveRel, revision, u)
		if err != nil {
			return blocked(err)
		}
		name := u.name
		if u.folder {
			name = archiveRel + "/" + u.name
		}
		if u.refusal == nil {
			paths := append([]string(nil), u.files...)
			if u.links != nil {
				for path := range u.links.Files {
					paths = append(paths, path)
				}
			}
			if u.folder {
				// Includes untracked files inside the folder.
				paths = append(paths, name)
			}
			status, err := gitOutput(ctx, root, append([]string{"--no-optional-locks", "status", "--porcelain=v1", "--untracked-files=all", "--"}, paths...)...)
			if err != nil {
				return blocked(err)
			}
			dirty := status != ""
			for _, path := range paths {
				if u.folder && path == name {
					continue
				}
				if !headPaths[path] {
					dirty = true
				}
			}
			if dirty {
				path := name
				if u.links != nil {
					path = archiveRel
				} else if !u.folder {
					path = spec.ArchiveDir(spec.ArchiveKind(u.name))
				}
				u.refusal = fmt.Errorf("has uncommitted changes under %s; commit or restore them and rerun", path)
			}
		}
		if u.refusal == nil && tagPaths != nil {
			for _, path := range u.files {
				if !tagPaths[path] {
					locator, err := spec.ResolveGitLocator(ctx, root, path, "history-full", "")
					if err != nil {
						return blocked(err)
					}
					if locator.Commit == "" {
						u.refusal = fmt.Errorf("history-full does not hold %s and no commit the default branch reaches holds it", path)
						break
					}
				}
			}
		}
		if u.refusal != nil {
			p.Report.Refused = append(p.Report.Refused, baselineUpdateHistoryRefusal{Unit: name, Reason: historyRefusalLine(u.refusal)})
			continue
		}
		report := baselineUpdateHistoryUnit{Unit: name, Files: len(u.files), Paths: append([]string{}, u.files...)}
		if u.links != nil {
			for path := range u.links.Files {
				report.Paths = append(report.Paths, path)
			}
		}
		sort.Strings(report.Paths)
		if c := u.conversion; c != nil {
			report.Action = "convert"
			report.Record = c.RecordPath
			report.Disposition = string(c.Record.Disposition)
			report.Bytes = c.Bytes
		} else if u.links != nil {
			report.Action = "rewrite"
			report.Files = len(u.links.Files)
			report.Links = u.links.Links
		} else {
			report.Action = string(u.kind.Action)
			report.Bytes = u.kind.BytesBefore
		}
		p.selected = append(p.selected, u)
		p.Report.Units = append(p.Report.Units, report)
	}
	if len(p.selected) > 0 {
		p.Report.Status = "pending"
		if p.Report.Tag == nil {
			p.Report.Tag = &baselineUpdateHistoryTag{Name: "history-full", Action: "create", Commit: revision, PushCommand: "git push origin history-full"}
		}
	}
	if err := p.listCitations(ctx); err != nil {
		return blocked(err)
	}
	if tagPaths != nil {
		stale, err := spec.StaleHistoryCitations(ctx, root, p.historyRoot, tagPaths)
		if err != nil {
			// Stale reporting is advisory; even its failure leaves the plan intact.
			p.Report.Message = historyRefusalLine(err)
		} else {
			for _, c := range stale {
				p.Report.StaleCitations = append(p.Report.StaleCitations, p.citation(c))
			}
		}
	}
	return p
}

func (p baselineUpdateHistoryPlan) citation(c spec.HistoryCitation) baselineUpdateHistoryCitation {
	return baselineUpdateHistoryCitation{Path: c.Path, Line: c.Line, Target: c.Target, Rewrite: c.Rewrite, Live: spec.IsLiveDocument(c.Path, p.specRootRel, p.archiveRootRel)}
}

func (p *baselineUpdateHistoryPlan) listCitations(ctx context.Context) error {
	promotable := map[string]string{}
	for _, u := range p.selected {
		if u.folder {
			for _, file := range u.files {
				if !baselineHistoryCoreArtifact(strings.TrimPrefix(file, u.conversion.Folder+"/")) {
					promotable[file] = file
				}
			}
		}
	}
	rewrite, err := planHistoryRewrite(ctx, p.root, p.specRootRel, p.archiveRootRel, p.selected)
	if err != nil {
		return err
	}
	p.rewrite = rewrite
	p.Report.Citations = []baselineUpdateHistoryCitation{}
	citations := rewrite.Citations

	for _, c := range citations {
		report := p.citation(c)
		if report.Live {
			report.Offer = promotable[c.Target]
		}
		p.Report.Citations = append(p.Report.Citations, report)
	}
	return nil
}

// This is the core-artifact refusal set of PlanLegacyConversion, folder-relative.
func baselineHistoryCoreArtifact(path string) bool {
	base := filepath.Base(path)
	return (filepath.Dir(path) == "." && strings.HasPrefix(base, "_") && strings.HasSuffix(base, ".md")) || (strings.HasPrefix(base, "task_") && strings.HasSuffix(base, ".md")) || (strings.HasPrefix(base, "qa-report-") && strings.HasSuffix(base, ".md"))
}

func (p *baselineUpdateHistoryPlan) promote(ctx context.Context, paths []string) error {
	if len(paths) == 0 {
		return nil
	}
	promotions := map[string][]string{}
	destinations := map[string]bool{}
	for _, path := range paths {
		if filepath.IsAbs(path) || filepath.ToSlash(filepath.Clean(path)) != path || strings.Contains(path, "\\") {
			return fmt.Errorf("promotion must be a clean repository-relative file: %q", path)
		}
		for _, refused := range p.Report.Refused {
			if strings.HasPrefix(path, refused.Unit+"/") {
				return fmt.Errorf("promotion %q is in refused unit %s: %s", path, refused.Unit, refused.Reason)
			}
		}
		found := false
		for _, u := range p.selected {
			if !u.folder || !strings.HasPrefix(path, u.conversion.Folder+"/") {
				continue
			}
			relative := strings.TrimPrefix(path, u.conversion.Folder+"/")
			if baselineHistoryCoreArtifact(relative) {
				return fmt.Errorf("cannot promote core artifact %q", path)
			}
			for _, file := range u.files {
				if path == file {
					found = true
					break
				}
			}
			if !found {
				break
			}
			destination := "docs/references/" + filepath.Base(path)
			if destinations[destination] {
				return fmt.Errorf("promotion %q has duplicate destination %s", path, destination)
			}
			if _, err := os.Lstat(filepath.Join(p.root, destination)); !os.IsNotExist(err) {
				if err != nil {
					return fmt.Errorf("promotion %q: inspect destination %s: %w", path, destination, err)
				}
				return fmt.Errorf("promotion %q: destination %s already exists", path, destination)
			}
			destinations[destination] = true
			promotions[u.name] = append(promotions[u.name], path)
			break
		}
		if !found {
			return fmt.Errorf("promotion %q is outside selected legacy folder files", path)
		}
	}
	for i := range p.selected {
		u := &p.selected[i]
		if len(promotions[u.name]) == 0 {
			continue
		}
		c := u.conversion
		// Selected files are committed at Revision. Their copies become references
		// read by every later Spec; PlanLegacyConversion never overwrites a destination.
		converted, err := spec.PlanLegacyConversion(spec.LegacyConversionRequest{RepositoryRoot: p.root, ArchiveRoot: p.archiveRoot, Slug: u.name, SourceRevision: p.Report.Revision, Created: c.Record.Created, Delivery: spec.LegacyDelivery{Commit: c.Record.DeliveryCommit, PullRequest: c.Record.PullRequest, Date: c.Record.Archived}, Promote: promotions[u.name]})
		if err != nil {
			return fmt.Errorf("promotion %q: %w", strings.Join(promotions[u.name], ", "), err)
		}
		u.conversion = &converted
	}
	if err := p.listCitations(ctx); err != nil {
		return err
	}
	promoted := map[string]string{}
	for _, paths := range promotions {
		for _, path := range paths {
			promoted[path] = "docs/references/" + filepath.Base(path)
		}
	}
	for i := range p.Report.Citations {
		c := &p.Report.Citations[i]
		if destination := promoted[c.Target]; destination != "" {
			c.Offer, c.PromotedTo = "", destination
		}
	}
	return nil
}

func baselineHistoryTree(ctx context.Context, root, revision string) (map[string]bool, error) {
	tree, err := gitOutput(ctx, root, "ls-tree", "-r", "--name-only", "-z", revision)
	if err != nil {
		return nil, err
	}
	paths := map[string]bool{}
	for _, p := range strings.Split(tree, "\x00") {
		paths[p] = true
	}
	return paths, nil
}

func (p baselineUpdateHistoryPlan) digest(baselinePlanDigest string) (string, error) {
	if len(p.selected) == 0 {
		return baselinePlanDigest, nil
	}
	type digestUnit struct {
		Unit         string            `json:"unit"`
		Action       string            `json:"action"`
		Paths        []string          `json:"paths"`
		Record       string            `json:"record,omitempty"`
		RecordDigest string            `json:"recordDigest,omitempty"`
		Links        int               `json:"links,omitempty"`
		Rewrites     map[string]string `json:"rewrites,omitempty"`
		Promoted     []string          `json:"promoted,omitempty"`
	}
	type digestCitation struct {
		Path    string `json:"path"`
		Line    int    `json:"line"`
		Target  string `json:"target"`
		Rewrite string `json:"rewrite,omitempty"`
	}
	payload := struct {
		Revision  string           `json:"revision"`
		Action    string           `json:"action"`
		Commit    string           `json:"commit"`
		Units     []digestUnit     `json:"units"`
		Citations []digestCitation `json:"citations,omitempty"`
	}{Revision: p.Report.Revision, Action: p.Report.Tag.Action, Commit: p.Report.Tag.Commit, Units: []digestUnit{}}
	for i, u := range p.selected {
		d := digestUnit{Unit: p.Report.Units[i].Unit, Action: p.Report.Units[i].Action, Paths: append([]string(nil), u.files...)}
		if u.links != nil {
			d.Links = u.links.Links
			d.Rewrites = map[string]string{}
			for path, content := range u.links.Files {
				d.Paths = append(d.Paths, path)
				sum := sha256.Sum256(content)
				d.Rewrites[path] = hex.EncodeToString(sum[:])
			}
		}
		sort.Strings(d.Paths)
		if c := u.conversion; c != nil {
			d.Record = c.RecordPath
			d.Promoted = c.Promoted
			sum := sha256.Sum256(c.Rendered)
			d.RecordDigest = hex.EncodeToString(sum[:])
		}
		payload.Units = append(payload.Units, d)
	}
	for _, c := range p.Report.Citations {
		payload.Citations = append(payload.Citations, digestCitation{Path: c.Path, Line: c.Line, Target: c.Target, Rewrite: c.Rewrite})
	}
	data, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("encode history plan digest: %w", err)
	}
	sum := sha256.Sum256(append([]byte("roundfix/baseline-update-history/v1\x00"+baselinePlanDigest+"\x00"), data...))
	return "sha256:" + hex.EncodeToString(sum[:]), nil
}

func applyBaselineUpdateHistory(ctx context.Context, repo string, p baselineUpdateHistoryPlan) (baselineUpdateHistory, error) {
	report := p.Report
	if len(p.selected) == 0 {
		return report, nil
	}
	repo = p.root
	if report.Tag.Action == "create" {
		if _, err := gitOutput(ctx, repo, "tag", "-a", "history-full", "-m", "docs/history before roundfix baseline update sanitized it (ADR-0254)", report.Revision); err != nil {
			return report, fmt.Errorf("create history-full: %w", err)
		}
		typ, err := gitOutput(ctx, repo, "cat-file", "-t", "refs/tags/history-full")
		if err != nil {
			return report, fmt.Errorf("check history-full: %w", err)
		}
		if typ != "tag" {
			return report, fmt.Errorf("history-full is not annotated after creation")
		}
		tag := *report.Tag
		tag.Action = "created"
		report.Tag = &tag
	}
	report.Applied = &baselineUpdateHistoryApplied{}
	for i, u := range p.selected {
		err := ctx.Err()
		if err == nil {
			if u.conversion != nil {
				err = spec.ApplyLegacyConversion(repo, *u.conversion)
			} else if u.links == nil {
				err = spec.ApplyHistoryKind(repo, *u.kind)
			}
		}
		if err != nil {
			return report, fmt.Errorf("sanitize %s: %w", report.Units[i].Unit, err)
		}
		report.Applied.Units++
		if c := u.conversion; c != nil {
			report.Applied.Records++
			report.Applied.Removed += len(c.Files)
			report.Applied.Bytes += c.Bytes
		} else if u.links == nil {
			report.Applied.Removed += len(u.files)
			report.Applied.Bytes += u.kind.BytesBefore
		}
	}
	var err error
	report.Applied.RewrittenCitations, report.Applied.RewrittenFiles, err = p.rewrite.Apply(ctx, repo)
	if err != nil {
		return report, fmt.Errorf("rewrite history citations: %w", err)
	}
	report.Status = "applied"
	return report, nil
}

func printBaselineUpdateHistory(w io.Writer, h baselineUpdateHistory) {
	if h.Status != "pending" && h.Status != "applied" && h.Status != "blocked" && len(h.Refused) == 0 && len(h.Citations) == 0 && len(h.StaleCitations) == 0 && h.Message == "" {
		return
	}
	fmt.Fprintf(w, "History: %s\n", h.Status)
	if h.Message != "" {
		fmt.Fprintf(w, "History result: %s\n", h.Message)
	}
	if h.Status == "pending" {
		var files int
		var bytes int64
		for _, u := range h.Units {
			if u.Action == "rewrite" {
				continue
			}
			files += u.Files
			bytes += u.Bytes
		}
		fmt.Fprintf(w, "History units: %d (%d file(s), %d bytes leave docs/history)\n", len(h.Units), files, bytes)
		for _, u := range h.Units {
			switch u.Action {
			case "convert":
				fmt.Fprintf(w, "- folder %s: removes %d file(s) (%d bytes) and writes %s (%s)\n", u.Unit, u.Files, u.Bytes, u.Record, u.Disposition)
			case "rewrite":
				fmt.Fprintf(w, "- record-links: rewrites %d link(s) in %d Archive Record(s)\n", u.Links, u.Files)
			case "remove":
				fmt.Fprintf(w, "- %s: removes %d file(s) (%d bytes)\n", u.Unit, u.Files, u.Bytes)
			}
			for _, path := range u.Paths {
				fmt.Fprintf(w, "  - %s\n", path)
			}
		}
	}
	if h.Status == "applied" {
		a := h.Applied
		fmt.Fprintf(w, "History applied: %d unit(s): wrote %d Archive Record(s), removed %d file(s) (%d bytes) kept in Git at %.12s and tag history-full; rewrote %d citation(s) in %d file(s)\n", a.Units, a.Records, a.Removed, a.Bytes, h.Revision, a.RewrittenCitations, a.RewrittenFiles)
	}
	if len(h.Refused) > 0 {
		fmt.Fprintf(w, "History refused: %d\n", len(h.Refused))
		for _, r := range h.Refused {
			fmt.Fprintf(w, "- refused %s: %s\n", r.Unit, r.Reason)
		}
	}
	if len(h.Citations) > 0 {
		fmt.Fprintf(w, "History citations: %d\n", len(h.Citations))
		for _, c := range h.Citations {
			fmt.Fprintf(w, "- cites %s:%d names %s", c.Path, c.Line, c.Target)
			switch {
			case c.PromotedTo != "":
				fmt.Fprintf(w, " (promoted to %s)", c.PromotedTo)
			case c.Offer != "":
				fmt.Fprintf(w, " (live; promote with --promote %s)", c.Offer)
			case c.Live:
				fmt.Fprint(w, " (live)")
			}
			if c.Rewrite != "" {
				fmt.Fprintf(w, "; rewrites to %s", c.Rewrite)
			}
			fmt.Fprintln(w)
		}
	}
	if len(h.StaleCitations) > 0 {
		fmt.Fprintf(w, "History stale citations: %d\n", len(h.StaleCitations))
		for _, c := range h.StaleCitations {
			fmt.Fprintf(w, "- cites %s:%d names %s; history-full holds it: git show history-full:%s\n", c.Path, c.Line, c.Target, c.Target)
		}
	}
	if h.Tag != nil {
		switch h.Tag.Action {
		case "create":
			fmt.Fprintf(w, "History tag: creates annotated history-full at %.12s; push it with %s\n", h.Tag.Commit, h.Tag.PushCommand)
		case "created":
			fmt.Fprintf(w, "History tag: created annotated history-full at %.12s; push it with %s\n", h.Tag.Commit, h.Tag.PushCommand)
		case "present":
			fmt.Fprintf(w, "History tag: history-full at %.12s holds every planned path\n", h.Tag.Commit)
		}
	}
}
