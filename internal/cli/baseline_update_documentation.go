package cli

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"roundfix/internal/baseline"
	"roundfix/internal/config"
	"roundfix/internal/spec"
	"roundfix/internal/speccheck"
)

type documentationADR struct {
	Path       string
	Retirement spec.Retirement
	Header     spec.ADRHeader
	Notices    []string
	Content    []byte
}

type documentationFacts struct {
	Modules               []string
	SpecRoot, HistoryRoot string
	Export                baseline.ExportListFacts
	ADRs                  []documentationADR
}

// gatherDocumentationFacts is the shared, read-only repository snapshot for the
// update and Doctor. It reads repository configuration, never a user's home.
func gatherDocumentationFacts(root string) (documentationFacts, error) {
	f := documentationFacts{ADRs: []documentationADR{}}
	anchored, err := os.OpenRoot(root)
	if err != nil {
		return f, fmt.Errorf("open documentation root: %w", err)
	}
	defer anchored.Close()
	manifest, err := anchored.ReadFile("docs/agents/setup-context.json")
	if errors.Is(err, os.ErrNotExist) {
		return f, baseline.ErrNoManifest
	}
	if err != nil {
		return f, fmt.Errorf("read Setup Manifest: %w", err)
	}
	var setup baseline.SetupManifest
	if err := json.Unmarshal(manifest, &setup); err != nil {
		return f, fmt.Errorf("read Setup Manifest: %w", err)
	}
	f.Modules = setup.Modules
	project, err := anchored.ReadFile(".roundfixrc.yml")
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return f, fmt.Errorf("read Project Config: %w", err)
	}
	cfg, err := config.ResolveConfigProposal(nil, project)
	if err != nil {
		return f, fmt.Errorf("read documentation configuration: %w", err)
	}
	specs, err := config.ResolveSpecsRootOrAbsent(config.Loaded{Config: cfg, GitRoot: root}, root)
	if err != nil {
		return f, err
	}
	local := func(path string) string {
		rel, err := filepath.Rel(root, path)
		if err != nil || rel == ".." || strings.HasPrefix(rel, ".."+string(filepath.Separator)) {
			return ""
		}
		return filepath.ToSlash(rel) + "/"
	}
	f.SpecRoot = local(specs.Path)
	archive := spec.ArchiveSpecRoot(specs.Path, specs.BuiltInRoot)
	history := archive
	if specs.BuiltInRoot {
		history = filepath.Dir(archive)
	}
	f.HistoryRoot = local(history)
	f.Export.ProcessPaths = []string{"AGENTS.md", "CLAUDE.md", "docs/agents/", "docs/backlog/", "docs/findings/", "docs/handoffs/"}
	for _, path := range []string{f.SpecRoot, f.HistoryRoot} {
		if path != "" {
			f.Export.ProcessPaths = append(f.Export.ProcessPaths, path)
		}
	}
	f.Export.GlossaryFiles, err = speccheck.GlossaryFiles(root)
	if err != nil {
		return f, err
	}
	docs, err := anchored.Open("docs")
	if err == nil {
		entries, readErr := docs.ReadDir(-1)
		closeErr := docs.Close()
		if readErr != nil {
			return f, readErr
		}
		if closeErr != nil {
			return f, closeErr
		}
		for _, entry := range entries {
			path := "docs/" + entry.Name()
			if entry.IsDir() {
				path += "/"
			}
			f.Export.DocsEntries = append(f.Export.DocsEntries, path)
		}
		slices.Sort(f.Export.DocsEntries)
	} else if !errors.Is(err, os.ErrNotExist) {
		return f, fmt.Errorf("read docs entries: %w", err)
	}
	err = fs.WalkDir(anchored.FS(), "docs/adr", func(path string, entry fs.DirEntry, walkErr error) error {
		if errors.Is(walkErr, os.ErrNotExist) && path == "docs/adr" {
			return nil
		}
		if walkErr != nil {
			return fmt.Errorf("read ADR directory %s: %w", path, walkErr)
		}
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".md") {
			return nil
		}
		content, err := anchored.ReadFile(path)
		if err != nil {
			return fmt.Errorf("read ADR %s: %w", path, err)
		}
		adr := documentationADR{Path: path, Retirement: spec.ClassifyADR(content), Header: spec.ReadADRHeader(content), Notices: spec.ADRAmendmentNotices(content), Content: content}
		f.ADRs = append(f.ADRs, adr)
		if !adr.Retirement.Retired && adr.Header.Knowledge != nil && *adr.Header.Knowledge {
			f.Export.KnowledgeADRs = append(f.Export.KnowledgeADRs, path)
		}
		return nil
	})
	if err != nil {
		return f, err
	}

	return f, nil
}

type baselineDocumentationExport struct {
	Path           string                `json:"path"`
	Action         string                `json:"action"`
	Standard       []string              `json:"standard"`
	Added          []string              `json:"added"`
	Removed        []string              `json:"removed"`
	Converted      []string              `json:"converted"`
	StopsExporting []baseline.ExportLoss `json:"stopsExporting"`
	Refusal        string                `json:"refusal"`
}
type baselineDocumentationADRs struct {
	Total    int `json:"total"`
	Unmarked int `json:"unmarked"`
}
type baselineDocumentationHeader struct {
	Path      string   `json:"path"`
	AmendedBy []string `json:"amendedBy"`
}
type baselineUpdateDocumentation struct {
	Status     string                        `json:"status"`
	ExportList baselineDocumentationExport   `json:"exportList"`
	ADRs       baselineDocumentationADRs     `json:"adrs"`
	ADRHeaders []baselineDocumentationHeader `json:"adrHeaders"`
}
type baselineDocumentationPlan struct {
	Report      baselineUpdateDocumentation
	Export      baseline.ExportListPlan
	FileChanges []baseline.FileChange
}

func newBaselineUpdateDocumentation() baselineUpdateDocumentation {
	return baselineUpdateDocumentation{Status: "current", ExportList: baselineDocumentationExport{Path: baseline.ExportListPath, Action: "none", Standard: []string{}, Added: []string{}, Removed: []string{}, Converted: []string{}, StopsExporting: []baseline.ExportLoss{}}, ADRHeaders: []baselineDocumentationHeader{}}
}

func documentationExportReport(p baseline.ExportListPlan) baselineDocumentationExport {
	return baselineDocumentationExport{Path: baseline.ExportListPath, Action: p.Action, Standard: p.Standard, Added: p.Added, Removed: p.Removed, Converted: p.Converted, StopsExporting: p.StopsExporting, Refusal: p.Refusal}
}

func documentationContentDigest(content []byte) string {
	sum := sha256.Sum256(content)
	return "sha256:" + hex.EncodeToString(sum[:])
}

func documentationNotices(adr documentationADR) []string {
	// Paths carry the record's own numeric identifier; never add it from a notice.
	own, _, _ := strings.Cut(filepath.Base(adr.Path), "-")
	result := []string{}
	for _, id := range adr.Notices {
		if strings.TrimLeft(strings.TrimPrefix(id, "ADR-"), "0") != strings.TrimLeft(own, "0") {
			result = append(result, id)
		}
	}
	return result
}

func planBaselineUpdateDocumentation(root string) (baselineDocumentationPlan, error) {
	p := baselineDocumentationPlan{Report: newBaselineUpdateDocumentation(), FileChanges: []baseline.FileChange{}}
	facts, err := gatherDocumentationFacts(root)
	if err != nil {
		return p, err
	}
	if slices.Contains(facts.Modules, "secondbrain") {
		anchored, err := os.OpenRoot(root)
		if err != nil {
			return p, err
		}
		current, readErr := anchored.ReadFile(baseline.ExportListPath)
		closeErr := anchored.Close()
		if readErr != nil && !errors.Is(readErr, os.ErrNotExist) {
			return p, fmt.Errorf("read export list: %w", readErr)
		}
		if closeErr != nil {
			return p, closeErr
		}
		p.Export = baseline.PlanExportList(current, readErr == nil, facts.Export)
		p.Report.ExportList = documentationExportReport(p.Export)
		if p.Export.Action == "create" || p.Export.Action == "update" {
			p.FileChanges = append(p.FileChanges, baseline.FileChange{Path: baseline.ExportListPath, Action: p.Export.Action, BeforeIdentity: documentationContentDigest(current), AfterIdentity: documentationContentDigest(p.Export.Content), ManagedEntries: []string{"export.standard-list"}})
		}
	}
	for _, adr := range facts.ADRs {
		if adr.Retirement.Retired {
			continue
		}
		p.Report.ADRs.Total++
		if adr.Header.Knowledge == nil {
			p.Report.ADRs.Unmarked++
		}
		notices := documentationNotices(adr)
		after := spec.WithADRAmendedBy(adr.Content, notices)
		if bytes.Equal(after, adr.Content) {
			continue
		}
		p.Report.ADRHeaders = append(p.Report.ADRHeaders, baselineDocumentationHeader{Path: adr.Path, AmendedBy: spec.ReadADRHeader(after).AmendedBy})
		p.FileChanges = append(p.FileChanges, baseline.FileChange{Path: adr.Path, Action: "update", BeforeIdentity: documentationContentDigest(adr.Content), AfterIdentity: documentationContentDigest(after), ManagedEntries: []string{"adr.amendment-list"}})
	}
	if len(p.FileChanges) > 0 {
		p.Report.Status = "pending"
	}
	return p, nil
}

func (p baselineDocumentationPlan) digest(previous string) string {
	if p.Report.Status != "pending" {
		return previous
	}
	// The payload is made only of strings and slices, and cannot fail encoding.
	payload, _ := json.Marshal(struct {
		Digest        string                        `json:"digest"`
		ExportAction  string                        `json:"exportAction"`
		ExportContent string                        `json:"exportContent"`
		Refusal       string                        `json:"refusal"`
		Headers       []baselineDocumentationHeader `json:"headers"`
	}{previous, p.Report.ExportList.Action, documentationContentDigest(p.Export.Content), p.Report.ExportList.Refusal, p.Report.ADRHeaders})
	return documentationContentDigest(payload)
}

func applyBaselineUpdateDocumentation(ctx context.Context, root string, p baselineDocumentationPlan) (baselineUpdateDocumentation, error) {
	report := p.Report
	if report.Status != "pending" {
		return report, nil
	}
	anchored, err := os.OpenRoot(root)
	if err != nil {
		return report, err
	}
	defer anchored.Close()
	return p.apply(ctx, root, anchored, anchored.WriteFile)
}

func (p baselineDocumentationPlan) apply(ctx context.Context, root string, anchored *os.Root, writeFile func(string, []byte, os.FileMode) error) (baselineUpdateDocumentation, error) {
	report := p.Report
	for _, header := range report.ADRHeaders {
		if err := ctx.Err(); err != nil {
			return report, err
		}
		content, err := anchored.ReadFile(header.Path)
		if errors.Is(err, os.ErrNotExist) {
			continue
		}
		if err != nil {
			return report, fmt.Errorf("read ADR %s: %w", header.Path, err)
		}
		after := spec.WithADRAmendedBy(content, header.AmendedBy)
		if bytes.Equal(content, after) {
			continue
		}
		if err := writeFile(header.Path, after, 0644); err != nil {
			return report, fmt.Errorf("write ADR %s: %w", header.Path, err)
		}
	}
	// History may have retired an ADR or changed its successor. Resolve the list
	// from the tree after those writes, never from the preview's stale bytes.
	if report.ExportList.Action != "none" && report.ExportList.Action != "refused" {
		if err := ctx.Err(); err != nil {
			return report, err
		}
		facts, err := gatherDocumentationFacts(root)
		if err != nil {
			return report, err
		}
		current, err := anchored.ReadFile(baseline.ExportListPath)
		if err != nil && !errors.Is(err, os.ErrNotExist) {
			return report, fmt.Errorf("read export list: %w", err)
		}
		refreshed := baseline.PlanExportList(current, err == nil, facts.Export)
		if refreshed.Action == "create" || refreshed.Action == "update" {
			if err := writeFile(baseline.ExportListPath, refreshed.Content, 0644); err != nil {
				return report, fmt.Errorf("write export list: %w", err)
			}
		}
	}
	report.Status = "applied"
	return report, nil
}

func printBaselineUpdateDocumentation(w io.Writer, d baselineUpdateDocumentation) {
	e := d.ExportList
	if d.Status != "pending" && d.Status != "applied" && e.Action != "refused" && d.ADRs.Unmarked == 0 {
		return
	}
	fmt.Fprintf(w, "Documentation: %s\n", d.Status)
	if e.Action == "create" || e.Action == "update" {
		fmt.Fprintf(w, "Export list: %s .secondbrain-export (%d line(s) in the Standard Export List)\n", e.Action, len(e.Standard))
	}
	for _, line := range e.Added {
		fmt.Fprintf(w, "- adds %s\n", line)
	}
	for _, line := range append(append([]string{}, e.Removed...), e.Converted...) {
		fmt.Fprintf(w, "- removes %s\n", line)
	}
	for _, kind := range []string{"process", "other"} {
		for _, loss := range e.StopsExporting {
			if loss.Kind != kind {
				continue
			}
			detail := "process"
			if kind == "other" {
				detail = "not in the standard list"
			}
			fmt.Fprintf(w, "- stops exporting %s (%s)\n", loss.Path, detail)
		}
	}

	if e.Action == "refused" {
		fmt.Fprintf(w, "Export list refused: %s; nothing is written to .secondbrain-export\n", e.Refusal)
	}
	if d.ADRs.Unmarked > 0 {
		fmt.Fprintf(w, "ADRs without a Knowledge Mark: %d of %d; they are not exported\n", d.ADRs.Unmarked, d.ADRs.Total)
	}
	if len(d.ADRHeaders) > 0 {
		fmt.Fprintf(w, "ADR headers: %d\n", len(d.ADRHeaders))
		for _, header := range d.ADRHeaders {
			fmt.Fprintf(w, "- %s gains amended_by: [%s]\n", header.Path, strings.Join(header.AmendedBy, ", "))
		}
	}
}
