package cli

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"path/filepath"
	"sort"
	"strings"

	"roundfix/internal/config"
	"roundfix/internal/spec"
)

type baselineUpdateHistory struct {
	Status   string                         `json:"status"`
	Message  string                         `json:"message,omitempty"`
	Revision string                         `json:"revision,omitempty"`
	Tag      *baselineUpdateHistoryTag      `json:"tag,omitempty"`
	Units    []baselineUpdateHistoryUnit    `json:"units"`
	Refused  []baselineUpdateHistoryRefusal `json:"refused"`
	Applied  *baselineUpdateHistoryApplied  `json:"applied,omitempty"`
}
type baselineUpdateHistoryTag struct {
	Name        string `json:"name"`
	Action      string `json:"action"`
	Commit      string `json:"commit"`
	PushCommand string `json:"pushCommand,omitempty"`
}
type baselineUpdateHistoryUnit struct {
	Unit        string `json:"unit"`
	Action      string `json:"action"`
	Record      string `json:"record,omitempty"`
	Disposition string `json:"disposition,omitempty"`
	Files       int    `json:"files"`
	Bytes       int64  `json:"bytes"`
	BytesAfter  int64  `json:"bytesAfter,omitempty"`
}
type baselineUpdateHistoryRefusal struct {
	Unit   string `json:"unit"`
	Reason string `json:"reason"`
}
type baselineUpdateHistoryApplied struct {
	Units   int   `json:"units"`
	Records int   `json:"records"`
	Reduced int   `json:"reduced"`
	Removed int   `json:"removed"`
	Bytes   int64 `json:"bytes"`
}
type baselineUpdateHistoryPlan struct {
	Report      baselineUpdateHistory
	selected    []historyUnit
	root        string
	historyRoot string
}

func newBaselineUpdateHistory(status string) baselineUpdateHistory {
	return baselineUpdateHistory{Status: status, Units: []baselineUpdateHistoryUnit{}, Refused: []baselineUpdateHistoryRefusal{}}
}

func planBaselineUpdateHistory(ctx context.Context, repo string, environment commandEnvironment) baselineUpdateHistoryPlan {
	p := baselineUpdateHistoryPlan{Report: newBaselineUpdateHistory("current")}
	blocked := func(err error) baselineUpdateHistoryPlan {
		p.Report.Status = "blocked"
		p.Report.Message = historyRefusalLine(err)
		p.Report.Units = []baselineUpdateHistoryUnit{}
		p.selected = nil
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
	units, err := historyInventory(root, archive)
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
			if u.folder {
				// Includes untracked files inside the folder.
				paths = append(paths, name)
			}
			status, err := gitOutput(ctx, root, append([]string{"--no-optional-locks", "status", "--porcelain=v1", "--untracked-files=all", "--"}, paths...)...)
			if err != nil {
				return blocked(err)
			}
			dirty := status != ""
			for _, path := range u.files {
				if !headPaths[path] {
					dirty = true
				}
			}
			if dirty {
				path := name
				if !u.folder {
					path = spec.ArchiveDir(spec.ArchiveKind(u.name))
				}
				u.refusal = fmt.Errorf("has uncommitted changes under %s; commit or restore them and rerun", path)
			}
		}
		if u.refusal == nil && tagPaths != nil {
			for _, path := range u.files {
				if !tagPaths[path] {
					u.refusal = fmt.Errorf("history-full does not hold %s", path)
					break
				}
			}
		}
		if u.refusal != nil {
			p.Report.Refused = append(p.Report.Refused, baselineUpdateHistoryRefusal{Unit: name, Reason: historyRefusalLine(u.refusal)})
			continue
		}
		report := baselineUpdateHistoryUnit{Unit: name, Files: len(u.files)}
		if c := u.conversion; c != nil {
			report.Action = "convert"
			report.Record = c.RecordPath
			report.Disposition = string(c.Record.Disposition)
			report.Bytes = c.Bytes
		} else {
			report.Action = string(u.kind.Action)
			report.Bytes = u.kind.BytesBefore
			if u.kind.Action == "reduce" {
				report.BytesAfter = u.kind.BytesAfter
			}
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
	return p
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
		Unit         string   `json:"unit"`
		Action       string   `json:"action"`
		Paths        []string `json:"paths"`
		Record       string   `json:"record,omitempty"`
		RecordDigest string   `json:"recordDigest,omitempty"`
		BytesAfter   *int64   `json:"bytesAfter,omitempty"`
	}
	payload := struct {
		Revision string       `json:"revision"`
		Action   string       `json:"action"`
		Commit   string       `json:"commit"`
		Units    []digestUnit `json:"units"`
	}{Revision: p.Report.Revision, Action: p.Report.Tag.Action, Commit: p.Report.Tag.Commit, Units: []digestUnit{}}
	for i, u := range p.selected {
		d := digestUnit{Unit: p.Report.Units[i].Unit, Action: p.Report.Units[i].Action, Paths: append([]string(nil), u.files...)}
		sort.Strings(d.Paths)
		if c := u.conversion; c != nil {
			d.Record = c.RecordPath
			sum := sha256.Sum256(c.Rendered)
			d.RecordDigest = hex.EncodeToString(sum[:])
		} else {
			d.BytesAfter = &u.kind.BytesAfter
		}
		payload.Units = append(payload.Units, d)
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
			} else {
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
		} else if u.kind.Action == "reduce" {
			report.Applied.Reduced += len(u.files)
		} else {
			report.Applied.Removed += len(u.files)
			report.Applied.Bytes += u.kind.BytesBefore
		}
	}
	report.Status = "applied"
	return report, nil
}

func printBaselineUpdateHistory(w io.Writer, h baselineUpdateHistory) {
	if h.Status != "pending" && h.Status != "applied" && h.Status != "blocked" && len(h.Refused) == 0 {
		return
	}
	fmt.Fprintf(w, "History: %s\n", h.Status)
	if h.Status == "blocked" {
		fmt.Fprintf(w, "History result: %s\n", h.Message)
	}
	if h.Status == "pending" {
		var files int
		var bytes int64
		for _, u := range h.Units {
			files += u.Files
			bytes += u.Bytes
		}
		fmt.Fprintf(w, "History units: %d (%d file(s), %d bytes leave docs/history)\n", len(h.Units), files, bytes)
		for _, u := range h.Units {
			switch u.Action {
			case "convert":
				fmt.Fprintf(w, "- folder %s: removes %d file(s) (%d bytes) and writes %s (%s)\n", u.Unit, u.Files, u.Bytes, u.Record, u.Disposition)
			case "reduce":
				fmt.Fprintf(w, "- %s: reduces %d file(s) from %d to %d bytes\n", u.Unit, u.Files, u.Bytes, u.BytesAfter)
			case "remove":
				fmt.Fprintf(w, "- %s: removes %d file(s) (%d bytes)\n", u.Unit, u.Files, u.Bytes)
			}
		}
	}
	if h.Status == "applied" {
		a := h.Applied
		fmt.Fprintf(w, "History applied: %d unit(s): wrote %d Archive Record(s), reduced %d file(s), removed %d file(s) (%d bytes) kept in Git at %.12s and tag history-full\n", a.Units, a.Records, a.Reduced, a.Removed, a.Bytes, h.Revision)
	}
	if len(h.Refused) > 0 {
		fmt.Fprintf(w, "History refused: %d\n", len(h.Refused))
		for _, r := range h.Refused {
			fmt.Fprintf(w, "- refused %s: %s\n", r.Unit, r.Reason)
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
