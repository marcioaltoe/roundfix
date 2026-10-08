package cli

// Suite: history sanitize refused units.
// Invariant: refusal reasons are observable and refused bytes survive selection and apply.
// Boundary IN: command dispatch and synthetic legacy folders in disposable Git repositories.
// Boundary OUT: stdout, stderr, exit codes and repository files.

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"roundfix/internal/spec"
)

func historyLegacyFolder(t *testing.T, h historyFixture, slug string) string {
	t.Helper()
	folder := filepath.Join(h.repo, "docs/history/specs", slug)
	historyWrite(t, filepath.Join(folder, "_prd.md"), "---\nspec: "+slug+"\ncreated: 2026-09-01\n---\n\n# "+slug+"\n\nPreserved outcome.\n")
	return folder
}

func TestHistorySanitizeConvertsLegacyUnprovenMaps(t *testing.T) {
	h := newHistoryFixture(t, false)
	folder := historyLegacyFolder(t, h, "aaa")
	historyWrite(t, filepath.Join(folder, "_prd.md"), "---\nspec: aaa\ncreated: 2026-09-01\nunproven:\n  - row: 03\n    goal: G2\n    claim: Catalog verification runs alone\n    satisfied-by: [task_04, task_05]\n  - row: 04\n    claim: Release is verified\n    reason: No standalone evidence\n---\n\n# Demo\n\nPreserved outcome.\n")
	historyCommit(t, h.repo, "Legacy unproven maps")
	historyTag(t, h)
	before := historySnapshot(t, h.repo)
	code, out, stderr := historyRun(t, h, "sanitize")
	if code != 0 || stderr != "" || strings.Contains(out, "refused") || !strings.Contains(out, "folder docs/history/specs/aaa:") {
		t.Fatalf("%d %s %s", code, out, stderr)
	}
	historyUnchanged(t, h, before, "")
	code, out, stderr = historyRun(t, h, "sanitize", "--apply", "--batch", "1")
	if code != 0 || stderr != "" || !strings.Contains(out, "wrote 1 Archive Record(s)") {
		t.Fatalf("%d %s %s", code, out, stderr)
	}
	assertPathMissing(t, folder)
	archived, err := spec.ReadArchivedSpec(filepath.Dir(folder), "aaa")
	want := []string{"row 03: Catalog verification runs alone (goal: G2; satisfied-by: task_04, task_05)", "row 04: Release is verified (reason: No standalone evidence)"}
	if err != nil || !reflect.DeepEqual(archived.Record.Unproven, want) {
		t.Fatalf("archive=%+v err=%v", archived, err)
	}
}

func TestHistorySanitizePrintsEachRefusalOnOneLine(t *testing.T) {
	h := newHistoryFixture(t, false)
	folder := historyLegacyFolder(t, h, "aaa")
	historyWrite(t, filepath.Join(folder, "_prd.md"), "---\nspec: aaa\nunproven: {claim: malformed}\n---\n\n# Demo\n")
	historyCommit(t, h.repo, "Malformed unproven")
	historyTag(t, h)
	before := historySnapshot(t, h.repo)
	_, cause := spec.BuildArchiveRecord(spec.ArchiveRecordInput{SpecDir: folder, Legacy: true})
	if cause == nil || !strings.Contains(cause.Error(), "\n") {
		t.Fatalf("fixture needs a multiline error: %v", cause)
	}
	wantRefusal := "refused docs/history/specs/aaa: build legacy record: " + strings.Join(strings.Fields(cause.Error()), " ")
	for _, args := range [][]string{{"sanitize"}, {"sanitize", "--apply", "--batch", "1"}} {
		code, out, stderr := historyRun(t, h, args...)
		wantCode := 0
		if len(args) > 1 {
			wantCode = 2
		}
		if code != wantCode || (wantCode == 0 && stderr != "") || strings.Count(out, "refused ") != 1 || !strings.Contains(out, wantRefusal+"\n") || !strings.Contains(wantRefusal, "cannot unmarshal !!map into") {
			t.Fatalf("%d %s %s; want %q", code, out, stderr, wantRefusal)
		}
		historyUnchanged(t, h, before, "")
	}
}

func TestHistorySanitizePlanListsEveryRefusedUnit(t *testing.T) {
	h := newHistoryFixture(t, false)
	for _, slug := range []string{"aaa", "bbb", "ccc", "ddd"} {
		folder := historyLegacyFolder(t, h, slug)
		if slug == "aaa" || slug == "ccc" {
			historyWrite(t, filepath.Join(folder, "_prd.md"), "invalid PRD")
		}
	}
	historyCommit(t, h.repo, "Mixed folders")
	before := historySnapshot(t, h.repo)
	code, out, err := historyRun(t, h, "sanitize")
	first, _, _ := strings.Cut(out, "\n")
	if code != 0 || err != "" || !strings.HasSuffix(first, "; 2 unit(s) refused") || strings.Count(out, "\nrefused ") != 2 {
		t.Fatalf("%d %s %s", code, out, err)
	}
	for _, slug := range []string{"aaa", "ccc"} {
		if !strings.Contains(out, "refused docs/history/specs/"+slug+": build legacy record: ") {
			t.Fatalf("missing reason: %s", out)
		}
	}
	if strings.Index(out, "refused docs/history/specs/aaa:") > strings.Index(out, "refused docs/history/specs/ccc:") {
		t.Fatal("refusals out of order")
	}
	historyUnchanged(t, h, before, "")
}

func TestHistorySanitizeApplySkipsARefusedUnitAndFillsTheBatch(t *testing.T) {
	h := newHistoryFixture(t, false)
	for _, slug := range []string{"aaa", "bbb", "ccc"} {
		historyLegacyFolder(t, h, slug)
	}
	historyWrite(t, filepath.Join(h.repo, "docs/history/specs/aaa/_prd.md"), "invalid PRD")
	historyCommit(t, h.repo, "Mixed batch")
	historyTag(t, h)
	before := snapshotDirectoryFiles(t, filepath.Join(h.repo, "docs/history/specs/aaa"))
	code, out, err := historyRun(t, h, "sanitize", "--apply", "--batch", "2")
	if code != 0 || err != "" || !strings.HasPrefix(out, "refused docs/history/specs/aaa: ") || !strings.Contains(out, "applied 2 unit(s): wrote 2 Archive Record(s)") || !strings.HasSuffix(out, "; 1 unit(s) refused; 1 unit(s) remain\n") {
		t.Fatalf("%d %s %s", code, out, err)
	}
	for _, slug := range []string{"bbb", "ccc"} {
		assertPathMissing(t, filepath.Join(h.repo, "docs/history/specs", slug))
		data, e := os.ReadFile(filepath.Join(h.repo, "docs/history/specs", slug+".md"))
		if e != nil {
			t.Fatal(e)
		}
		if _, e := spec.ParseArchiveRecord(data); e != nil {
			t.Fatal(e)
		}
	}
	if !reflect.DeepEqual(before, snapshotDirectoryFiles(t, filepath.Join(h.repo, "docs/history/specs/aaa"))) {
		t.Fatal("refused bytes changed")
	}
}

func TestHistorySanitizeApplyRefusesWhenEveryExaminedUnitIsRefused(t *testing.T) {
	h := newHistoryFixture(t, false)
	for _, slug := range []string{"aaa", "bbb"} {
		folder := historyLegacyFolder(t, h, slug)
		historyWrite(t, filepath.Join(folder, "_prd.md"), "invalid PRD")
	}
	historyCommit(t, h.repo, "All refused")
	historyTag(t, h)
	before := historySnapshot(t, h.repo)
	code, out, err := historyRun(t, h, "sanitize", "--apply", "--batch", "1")
	if code != 2 || strings.Count(out, "refused docs/history/specs/") != 2 || strings.Contains(out, "applied") || !strings.Contains(err, "history sanitize --apply found no convertible unit; 2 unit(s) refused") {
		t.Fatalf("%d %s %s", code, out, err)
	}
	historyUnchanged(t, h, before, "")
}

func TestHistorySanitizeRefusesAPromotionInARefusedUnit(t *testing.T) {
	h := newHistoryFixture(t, true)
	historyWrite(t, filepath.Join(h.repo, "docs/history/specs/aaa/_prd.md"), "invalid PRD")
	historyCommit(t, h.repo, "Refused promotion folder")
	historyTag(t, h)
	before := historySnapshot(t, h.repo)
	p := "docs/history/specs/aaa/references/lesson.md"
	code, out, err := historyRun(t, h, "sanitize", "--apply", "--batch", "1", "--promote", p)
	if code != 2 || out != "" || !strings.Contains(err, fmt.Sprintf("promotion %q is in refused unit docs/history/specs/aaa: build legacy record: ", p)) {
		t.Fatalf("%d %s %s", code, out, err)
	}
	historyUnchanged(t, h, before, "")
}

func TestHistorySanitizeRefusesAKindUnitAndPlansTheRest(t *testing.T) {
	h := newHistoryFixture(t, true)
	historyWrite(t, filepath.Join(h.repo, "docs/history/findings/entry.md"), "# Retired\n\nMissing front matter.\n")
	historyCommit(t, h.repo, "Malformed Finding")
	before := historySnapshot(t, h.repo)
	code, out, err := historyRun(t, h, "sanitize")
	if code != 0 || err != "" || !strings.Contains(out, "refused findings: plan history kind findings: reduce history entry") || !strings.Contains(out, "missing front matter") || !strings.Contains(out, "backlog: reduces") || !strings.Contains(out, "reviews: removes") || strings.Count(out, "\nfolder ") != 3 {
		t.Fatalf("%d %s %s", code, out, err)
	}
	historyUnchanged(t, h, before, "")
}

func TestHistorySanitizeConvertsTheFourLegacyShapes(t *testing.T) {
	h := newHistoryFixture(t, false)
	// The synthetic manifest uses the public persisted schema, not adopter files.
	const graph = `---
schema: spec-tasks/v1
qa: task_qa
graph:
  nodes:
    - id: task_01
      file: task_01.md
    # - id: task_02
    #   file: task_02.md
    - id: task_qa
      file: task_qa.md
      needs: [task_01]
---

| id | title | type | complexity | needs |
| --- | --- | --- | --- | --- |
`
	rows := []string{
		"| task_01 | Implementation | backend | low | - |\n| task_02 | Commented | backend | low | - |\n",
		"| task_01 | Implementation | backend | low | - |\n| task_03 | Omitted | backend | low | - |\n| task_04 | Omitted | backend | low | - |\n",
		"| task_01 | Implementation | refactor | low | - |\n",
		"| task_01 | Implementation | backend | low | - |\n",
	}
	for i, slug := range []string{"aaa", "bbb", "ccc", "ddd"} {
		folder := historyLegacyFolder(t, h, slug)
		historyWrite(t, filepath.Join(folder, "_tasks.md"), graph+rows[i]+"| task_qa | QA | qa | low | task_01 |\n")
		for _, task := range []string{"task_01", "task_qa"} {
			historyWrite(t, filepath.Join(folder, task+".md"), "---\ntask: "+task+"\nspec: "+slug+"\nstatus: completed\ntype: backend\n---\n\n# Task\n")
		}
		if slug == "ddd" {
			for day := 1; day <= 3; day++ {
				historyWrite(t, filepath.Join(folder, "qa", fmt.Sprintf("qa-report-2026-09-%02d.md", day)), "---\nverdict: fail\n---\n\n# Failed QA\n")
			}
		}
	}
	historyCommit(t, h.repo, "Four legacy shapes")
	historyTag(t, h)
	before := historySnapshot(t, h.repo)
	code, out, err := historyRun(t, h, "sanitize")
	if code != 0 || err != "" || strings.Contains(out, "refused") || !strings.Contains(out, "failed-qa") {
		t.Fatalf("%d %s %s", code, out, err)
	}
	for _, line := range []string{
		"tolerates docs/history/specs/aaa: projection row task_02 names a Task outside the graph",
		"tolerates docs/history/specs/bbb: projection row task_03 names a Task outside the graph",
		"tolerates docs/history/specs/bbb: projection row task_04 names a Task outside the graph",
		`tolerates docs/history/specs/ccc: projection row task_01 has retired Task type "refactor"`,
	} {
		if !strings.Contains(out, line+"\n") {
			t.Fatalf("missing %q: %s", line, out)
		}
	}
	historyUnchanged(t, h, before, "")
	code, out, err = historyRun(t, h, "sanitize", "--apply", "--batch", "4")
	if code != 0 || err != "" || !strings.Contains(out, "applied 4 unit(s): wrote 4 Archive Record(s)") {
		t.Fatalf("%d %s %s", code, out, err)
	}
	for _, slug := range []string{"aaa", "bbb", "ccc", "ddd"} {
		assertPathMissing(t, filepath.Join(h.repo, "docs/history/specs", slug))
		data, e := os.ReadFile(filepath.Join(h.repo, "docs/history/specs", slug+".md"))
		if e != nil {
			t.Fatal(e)
		}
		record, e := spec.ParseArchiveRecord(data)
		if e != nil {
			t.Fatal(e)
		}
		if slug == "ddd" && (record.Disposition != spec.ArchiveFailedQA || record.QAVerdict != spec.VerdictFail || record.QAReport != "qa-report-2026-09-03.md" || record.QAOverride != nil || spec.ArchivedTaskCompleted(record, record.QATask)) {
			t.Fatalf("failed QA record=%+v", record)
		}
	}
}

func TestHistorySanitizeInventoryRefusalsAndBatchBoundary(t *testing.T) {
	for _, location := range []string{"docs/history/specs/aaa/link", "docs/history/findings/link"} {
		t.Run(location, func(t *testing.T) {
			h := newHistoryFixture(t, false)
			historyLegacyFolder(t, h, "aaa")
			historyLegacyFolder(t, h, "bbb")
			link := filepath.Join(h.repo, location)
			if err := os.MkdirAll(filepath.Dir(link), 0755); err != nil {
				t.Fatal(err)
			}
			if err := os.Symlink("../../missing", link); err != nil {
				t.Fatal(err)
			}
			historyCommit(t, h.repo, "Inventory refusal")
			code, out, err := historyRun(t, h, "sanitize")
			if code != 0 || err != "" || !strings.Contains(out, "Spec contains symlink") || !strings.Contains(out, "; 1 unit(s) refused") || !strings.Contains(out, "folder docs/history/specs/bbb:") {
				t.Fatalf("%d %s %s", code, out, err)
			}
			code, out, err = historyRun(t, h, "sanitize", "--batch", "1")
			if code != 0 || err != "" {
				t.Fatalf("%d %s %s", code, out, err)
			}
			if strings.HasPrefix(location, "docs/history/findings/") && strings.Contains(out, "refused") {
				t.Fatalf("examined kind after batch: %s", out)
			}
		})
	}
}

func TestHistorySanitizeTagCoverageSelectsOnlyConvertibleUnits(t *testing.T) {
	h := newHistoryFixture(t, false)
	historyLegacyFolder(t, h, "bbb")
	historyCommit(t, h.repo, "Covered folder")
	historyTag(t, h)
	folder := historyLegacyFolder(t, h, "aaa")
	historyWrite(t, filepath.Join(folder, "_prd.md"), "invalid PRD")
	historyCommit(t, h.repo, "Uncovered refused folder")
	before := snapshotDirectoryFiles(t, folder)
	code, out, err := historyRun(t, h, "sanitize", "--apply", "--batch", "1")
	if code != 0 || err != "" || !strings.Contains(out, "refused docs/history/specs/aaa:") || !strings.Contains(out, "wrote 1 Archive Record(s)") {
		t.Fatalf("%d %s %s", code, out, err)
	}
	if !reflect.DeepEqual(before, snapshotDirectoryFiles(t, folder)) {
		t.Fatal("uncovered refused folder changed")
	}
}

func TestHistorySanitizePromotionSeparatesFolderFromKind(t *testing.T) {
	h := newHistoryFixture(t, false)
	folder := historyLegacyFolder(t, h, "findings")
	historyWrite(t, filepath.Join(folder, "references/lesson.md"), "Reusable knowledge.\n")
	historyWrite(t, filepath.Join(h.repo, "docs/history/findings/entry.md"), "---\nstatus: closed\n---\n\n# Retired\n\nFirst paragraph.\n\n## Details\n\nMore.\n")
	historyCommit(t, h.repo, "Same folder and kind name")
	historyTag(t, h)
	code, out, err := historyRun(t, h, "sanitize", "--apply", "--batch", "2", "--promote", "docs/history/specs/findings/references/lesson.md")
	if code != 0 || err != "" || !strings.Contains(out, "wrote 1 Archive Record(s), reduced 1 file(s)") || !strings.Contains(out, "promoted 1 file(s)") {
		t.Fatalf("%d %s %s", code, out, err)
	}
	data, e := os.ReadFile(filepath.Join(h.repo, "docs/references/lesson.md"))
	if e != nil || string(data) != "Reusable knowledge.\n" {
		t.Fatalf("promotion=%s err=%v", data, e)
	}
}
