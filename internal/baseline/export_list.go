package baseline

import (
	"bytes"
	"fmt"
	"sort"
	"strings"
)

const (
	ExportListPath        = ".secondbrain-export"
	ExportListBeginMarker = "# roundfix:standard-export:begin"
	ExportListEndMarker   = "# roundfix:standard-export:end"
)

// ExportListFacts is the repository snapshot shared by update and Doctor.
type ExportListFacts struct {
	GlossaryFiles []string
	KnowledgeADRs []string
	ProcessPaths  []string
	DocsEntries   []string
}

type ExportLoss struct {
	Path string `json:"path"`
	Kind string `json:"kind"`
}

type ExportListPlan struct {
	Action         string
	Content        []byte
	Standard       []string
	Added, Removed []string
	Converted      []string
	StopsExporting []ExportLoss
	Refusal        string
}

func exportPath(line string) string { return strings.TrimRight(strings.TrimSpace(line), "/") }
func exportWithin(path, parent string) bool {
	return path == parent || strings.HasPrefix(path, parent+"/")
}
func exportProcess(path string, facts ExportListFacts) bool {
	for _, process := range facts.ProcessPaths {
		if exportWithin(exportPath(path), exportPath(process)) {
			return true
		}
	}
	return false
}
func exportContainsProcess(path string, facts ExportListFacts) bool {
	for _, process := range facts.ProcessPaths {
		if exportWithin(exportPath(process), exportPath(path)) {
			return true
		}
	}
	return false
}
func exportEntries(lines [][]byte) []string {
	result := []string{}
	for _, line := range lines {
		s := strings.TrimSpace(string(line))
		if s != "" && !strings.HasPrefix(s, "#") {
			result = append(result, s)
		}
	}
	return result
}
func exportDifference(a, b []string) []string {
	result := []string{}
	for _, item := range a {
		found := false
		for _, other := range b {
			if exportPath(item) == exportPath(other) {
				found = true
				break
			}
		}
		if !found {
			result = append(result, item)
		}
	}
	return result
}

// PlanExportList is pure: it reads neither the filesystem nor configuration.
// Existing bytes outside a valid marker pair are always preserved.
func PlanExportList(current []byte, exists bool, facts ExportListFacts) ExportListPlan {
	p := ExportListPlan{Action: "unchanged", Content: current, Standard: []string{}, Added: []string{}, Removed: []string{}, Converted: []string{}, StopsExporting: []ExportLoss{}}
	candidates := append([]string{"README.md"}, facts.GlossaryFiles...)
	candidates = append(candidates, "docs/user-guide/", "docs/references/", "docs/design/", "docs/lessons/")
	adrs := append([]string{}, facts.KnowledgeADRs...)
	sort.Strings(adrs)
	candidates = append(candidates, adrs...)
	seen := map[string]bool{}
	for _, entry := range candidates {
		if !seen[exportPath(entry)] && !exportProcess(entry, facts) && !exportContainsProcess(entry, facts) {
			p.Standard = append(p.Standard, entry)
			seen[exportPath(entry)] = true
		}
	}
	block := []byte(ExportListBeginMarker + "\n# Kept by the Baseline: roundfix baseline update rewrites the lines between these two markers.\n" + strings.Join(p.Standard, "\n") + "\n" + ExportListEndMarker + "\n")
	if !exists {
		p.Action = "create"
		p.Content = block
		p.Added = append(p.Added, p.Standard...)
		return p
	}
	lines := bytes.SplitAfter(current, []byte("\n"))
	begins, ends := []int{}, []int{}
	for i, line := range lines {
		switch strings.TrimSpace(string(line)) {
		case ExportListBeginMarker:
			begins = append(begins, i)
		case ExportListEndMarker:
			ends = append(ends, i)
		}
	}
	if len(begins)+len(ends) > 0 {
		switch {
		case len(begins) == 0:
			p.Refusal = ExportListBeginMarker + " is missing"
		case len(ends) == 0:
			p.Refusal = ExportListEndMarker + " is missing"
		case len(begins) > 1:
			p.Refusal = fmt.Sprintf("%s appears %d times", ExportListBeginMarker, len(begins))
		case len(ends) > 1:
			p.Refusal = fmt.Sprintf("%s appears %d times", ExportListEndMarker, len(ends))
		case ends[0] < begins[0]:
			p.Refusal = ExportListEndMarker + " comes before " + ExportListBeginMarker
		}
		if p.Refusal != "" {
			p.Action = "refused"
			return p
		}
		start, end := begins[0], ends[0]
		old := exportEntries(lines[start+1 : end])
		p.Added = exportDifference(p.Standard, old)
		p.Removed = exportDifference(old, p.Standard)
		replacement := block
		if !bytes.HasSuffix(lines[end], []byte("\n")) {
			replacement = bytes.TrimSuffix(block, []byte("\n"))
		}
		p.Content = append(bytes.Join(lines[:start], nil), replacement...)
		p.Content = append(p.Content, bytes.Join(lines[end+1:], nil)...)
	} else {
		kept := [][]byte{}
		old := exportEntries(lines)
		broad := []string{}
		losses := map[string]ExportLoss{}
		for _, line := range lines {
			text := strings.TrimSpace(string(line))
			path := exportPath(text)
			remove := false
			if text != "" && !strings.HasPrefix(text, "#") {
				if strings.HasPrefix(text, "!") {
					remove = exportProcess(strings.TrimPrefix(path, "!"), facts)
				} else {
					remove = seen[path] || exportProcess(path, facts) || exportContainsProcess(path, facts)
					if remove && exportProcess(path, facts) {
						losses[text] = ExportLoss{Path: text, Kind: "process"}
					}
					if remove && exportContainsProcess(path, facts) {
						broad = append(broad, path)
					}
				}
			}
			if remove {
				p.Converted = append(p.Converted, text)
			} else {
				kept = append(kept, line)
			}
		}
		p.Added = exportDifference(p.Standard, old)
		own := exportEntries(kept)
		for _, parent := range broad {
			for _, child := range facts.DocsEntries {
				if !exportWithin(exportPath(child), parent) {
					continue
				}
				covered := false
				for _, entry := range append(append([]string{}, p.Standard...), own...) {
					if strings.HasPrefix(entry, "!") {
						if exportWithin(exportPath(child), exportPath(strings.TrimPrefix(entry, "!"))) {
							covered = false
						}
						continue
					}
					if exportWithin(exportPath(child), exportPath(entry)) {
						covered = true
					}
				}
				if !covered {
					kind := "other"
					if exportProcess(child, facts) {
						kind = "process"
					}
					losses[child] = ExportLoss{Path: child, Kind: kind}
				}
			}
		}
		keys := make([]string, 0, len(losses))
		for key := range losses {
			keys = append(keys, key)
		}
		sort.Strings(keys)
		for _, key := range keys {
			p.StopsExporting = append(p.StopsExporting, losses[key])
		}
		p.Content = append(append([]byte{}, block...), '\n')
		p.Content = append(p.Content, bytes.Join(kept, nil)...)
	}
	if !bytes.Equal(current, p.Content) {
		p.Action = "update"
	}
	return p
}
