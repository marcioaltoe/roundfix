package cli

// Suite: Baseline update names every planned history path.
// Invariant: preview lists exactly the removals and record repairs apply performs.
// Boundary IN: public command arguments and an adopted temporary Git repository.
// Boundary OUT: JSON, text, approval digest and repository files.

import (
	"bytes"
	"encoding/json"
	"path/filepath"
	"reflect"
	"sort"
	"strings"
	"testing"
)

func TestBaselineUpdateNamesEveryFileItWillRemove(t *testing.T) {
	t.Parallel()
	h, record, _ := recordLinksFixture(t, true, true)
	folder := "docs/history/specs/0004-legacy"
	historyLegacyFolder(t, h, "0004-legacy")
	historyWrite(t, filepath.Join(h.repo, folder, "references/lesson.md"), "Preserved lesson.\n")
	for _, path := range []string{"docs/adr/0001-retired.md", "docs/backlog/2026-10-17-done.md", "_archived/findings/old.md"} {
		historyWrite(t, filepath.Join(h.repo, path), "---\nstatus: deprecated\n---\n\n# Retired\n")
	}
	h.revision = historyCommit(t, h.repo, "docs: all pending history shapes")
	historyTag(t, h)
	before, refs, index := historySnapshot(t, h.repo), baselineHistoryGit(t, h, "for-each-ref"), baselineHistoryIndex(t, h)
	preview, out, stderr, code := baselineHistoryRun(t, h, "json")
	if code != 3 || stderr != "" || preview.History.Status != "pending" || len(preview.History.Refused) != 0 || preview.History.Tag.Action != "present" {
		t.Fatalf("preview=%d %s %s", code, out, stderr)
	}
	var wire struct {
		History struct {
			Units []struct {
				Unit   string   `json:"unit"`
				Action string   `json:"action"`
				Paths  []string `json:"paths"`
			} `json:"units"`
		} `json:"history"`
	}
	if err := json.Unmarshal([]byte(out), &wire); err != nil {
		t.Fatal(err)
	}
	want := map[string][]string{
		folder:         {folder + "/_prd.md", folder + "/references/lesson.md"},
		"findings":     {"_archived/findings/old.md", "docs/history/findings/reduced.md"},
		"backlog":      {"docs/backlog/2026-10-17-done.md"},
		"adr":          {"docs/adr/0001-retired.md"},
		"superseded":   {"docs/history/specs/0002-unused.md"},
		"record-links": {record},
	}
	if len(wire.History.Units) != len(want) {
		t.Fatalf("units=%+v, want %v", wire.History.Units, want)
	}
	_, text, stderr, code := baselineHistoryRun(t, h, "text")
	if code != 3 || stderr != "" || !strings.Contains(text, "Plan Digest: "+preview.PlanDigest+"\n") {
		t.Fatalf("text preview=%d %s %s", code, text, stderr)
	}
	removals, seen := map[string]bool{}, map[string]bool{}
	var listed []string
	for _, unit := range wire.History.Units {
		if !reflect.DeepEqual(unit.Paths, want[unit.Unit]) || !sort.StringsAreSorted(unit.Paths) {
			t.Fatalf("unit %s paths=%v, want %v", unit.Unit, unit.Paths, want[unit.Unit])
		}
		delete(want, unit.Unit)
		// Paths must immediately follow their own unit's summary, in JSON order.
		label := "- " + unit.Unit + ":"
		if unit.Action == "convert" {
			label = "- folder " + unit.Unit + ":"
		}
		start := strings.Index(text, label)
		if start < 0 {
			t.Fatalf("missing unit %s in %s", label, text)
		}
		_, tail, _ := strings.Cut(text[start:], "\n")
		var lines strings.Builder
		for _, path := range unit.Paths {
			if seen[path] {
				t.Fatalf("path listed more than once: %s", path)
			}
			seen[path] = true
			listed = append(listed, path)
			lines.WriteString("  - " + path + "\n")
			if unit.Action != "rewrite" {
				removals[path] = true
			}
		}
		if !strings.HasPrefix(tail, lines.String()) {
			t.Fatalf("unit %s missing ordered path lines %q in %s", unit.Unit, lines.String(), text)
		}
	}
	var textPaths []string
	for line := range strings.SplitSeq(text, "\n") {
		if path, ok := strings.CutPrefix(line, "  - "); ok {
			textPaths = append(textPaths, path)
		}
	}
	if len(want) != 0 || !reflect.DeepEqual(textPaths, listed) {
		t.Fatalf("unlisted units=%v text paths=%v JSON paths=%v", want, textPaths, listed)
	}
	baselineHistoryUnchanged(t, h, before, refs, index)
	result, out, stderr, code := baselineHistoryRun(t, h, "json", "--yes")
	if code != 0 || stderr != "" || result.History.Applied == nil || result.History.Applied.Removed != len(removals) || result.ApprovedPlanDigest != preview.PlanDigest {
		t.Fatalf("apply=%d %s %s", code, out, stderr)
	}
	after := historySnapshot(t, h.repo)
	for path := range before {
		_, exists := after[path]
		if exists == removals[path] {
			t.Errorf("path %s exists=%v, planned removal=%v", path, exists, removals[path])
		}
	}
	if bytes.Equal(before[record], after[record]) || strings.Contains(string(after[record]), "](file.md)") {
		t.Fatalf("record-links did not repair %s: %s", record, after[record])
	}
	if _, exists := after[folder+".md"]; !exists {
		t.Fatal("conversion did not write its Archive Record")
	}
}
