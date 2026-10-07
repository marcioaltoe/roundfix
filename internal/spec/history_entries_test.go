package spec

import (
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"roundfix/internal/gittest"
)

func historyEntryWrite(t *testing.T, root, path, content string) {
	t.Helper()
	absolute := filepath.Join(root, path)
	if err := os.MkdirAll(filepath.Dir(absolute), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(absolute, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}

func TestReduceHistoryEntryKeepsFrontMatterTitleAndFirstParagraph(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct{ name, front, body, paragraph string }{
		{"finding", "---\nstatus: done\nabsorbed_by: 0001-owner\nextra: 'keep  spacing' # comment\n---\n", "# Finding\n\nFirst  paragraph\nwith\tspacing.\n\n## Detail\nDiscarded.\n", "First paragraph with spacing."},
		{"backlog", "---\r\ntype: fix\r\nstatus: closed\r\nreason: no work remains\r\nspec: null\r\n---\r\n", "# Backlog\r\n\r\n## Symptom\r\n\r\nThe   first\r\nparagraph.\r\n\r\n## Expected\r\nDiscarded.\r\n", "The first paragraph."},
		{"hash paragraph", "---\nstatus: done\n---\n", "# Finding\n\n#hashtag is paragraph text.\n\nOther text.\n", "#hashtag is paragraph text."},
	} {
		t.Run(tc.name, func(t *testing.T) {
			path := "docs/history/" + tc.name + "/entry.md"
			revision := strings.Repeat("a", 40)
			got, err := ReduceHistoryEntry([]byte(tc.front+"\n"+tc.body), revision, path)
			if err != nil {
				t.Fatal(err)
			}
			title := "# Finding"
			if tc.name == "backlog" {
				title = "# Backlog"
			}
			want := tc.front + "\n" + title + "\n\n" + tc.paragraph + "\n\nFull text in Git at `" + revision + "`: `" + path + "`.\n"
			if string(got) != want {
				t.Fatalf("reduced bytes = %q, want %q", got, want)
			}
			if !IsReducedHistoryEntry(got) {
				t.Fatal("reduction not recognized")
			}
		})
	}
}

func TestReduceHistoryEntryRefusesWhatItCannotReduce(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct{ name, content string }{
		{"no front matter", "# Title\n\nParagraph.\n"},
		{"unclosed front matter", "---\nstatus: done\n# Title\n"},
		{"no title", "---\nstatus: done\n---\n## Only a section\n"},
		{"already reduced", "---\nstatus: done\n---\n# Title\n\nFull text in Git at `" + strings.Repeat("b", 64) + "`: `docs/history/findings/a.md`.\n\n"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if _, err := ReduceHistoryEntry([]byte(tc.content), strings.Repeat("a", 40), "entry.md"); err == nil {
				t.Fatal("expected refusal")
			}
		})
	}
	for _, tc := range []struct {
		name, last string
		want       bool
	}{
		{"sha1", "Full text in Git at `" + strings.Repeat("a", 40) + "`: `entry.md`.", true},
		{"sha256", "Full text in Git at `" + strings.Repeat("B", 64) + "`: `entry.md`.", true},
		{"short revision", "Full text in Git at `abc`: `entry.md`.", false},
		{"non hex", "Full text in Git at `" + strings.Repeat("z", 40) + "`: `entry.md`.", false},
		{"later content", "Full text in Git at `" + strings.Repeat("a", 40) + "`: `entry.md`.\nLater.", false},
		{"trailing whitespace", "Full text in Git at `" + strings.Repeat("a", 40) + "`: `entry.md`. ", false},
		{"leading whitespace", " Full text in Git at `" + strings.Repeat("a", 40) + "`: `entry.md`.", false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if got := IsReducedHistoryEntry([]byte(tc.last + "\n\n")); got != tc.want {
				t.Fatalf("detector = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestReducedHistoryEntryIsNotPendingAgain(t *testing.T) {
	t.Parallel()
	root := t.TempDir()
	path := ArchiveDir(ArchiveKindFinding) + "/a.md"
	historyEntryWrite(t, root, path, "---\nstatus: done\n---\n# Title\n\nSummary.\n\nDetails.\n")
	plans, err := PlanHistoryKinds(root, strings.Repeat("a", 40))
	if err != nil {
		t.Fatal(err)
	}
	if len(plans) != 1 {
		t.Fatalf("plans = %+v", plans)
	}
	if err := ApplyHistoryKind(root, plans[0]); err != nil {
		t.Fatal(err)
	}
	plans, err = PlanHistoryKinds(root, strings.Repeat("b", 40))
	if err != nil || len(plans) != 0 {
		t.Fatalf("second plan = %+v, %v", plans, err)
	}
}

func TestPlanHistoryKindsOrdersAndMeasures(t *testing.T) {
	t.Parallel()
	root := t.TempDir()
	full := "---\nstatus: done\n---\n# Title\n\nSummary.\n\nLong discarded details.\n"
	revision := strings.Repeat("a", 40)
	files := []string{"docs/history/findings/z.md", "docs/history/findings/a.md", "docs/history/backlog/b.md", "docs/history/reviews/session/a.json", "docs/history/handoffs/h.md"}
	for _, path := range files {
		historyEntryWrite(t, root, path, full)
	}
	adr := "# ADR\n\nPreserved decision.\n"
	historyEntryWrite(t, root, "docs/history/adr/a.md", adr)
	plans, err := PlanHistoryKinds(root, revision)
	if err != nil {
		t.Fatal(err)
	}
	wantFiles := [][]string{{files[1], files[0]}, {files[2]}, {files[3]}, {files[4]}}
	if len(plans) != 4 {
		t.Fatalf("plans = %+v", plans)
	}
	for i, kind := range []ArchiveKind{ArchiveKindFinding, ArchiveKindBacklog, ArchiveKindReview, ArchiveKindHandoff} {
		plan := plans[i]
		if plan.Kind != kind || !reflect.DeepEqual(plan.Files, wantFiles[i]) || plan.BytesBefore != int64(len(full)*len(wantFiles[i])) {
			t.Fatalf("plan %d = %+v", i, plan)
		}
		wantAfter := int64(0)
		if i < 2 {
			if plan.Action != "reduce" {
				t.Fatalf("action = %q", plan.Action)
			}
			for _, path := range wantFiles[i] {
				wantAfter += int64(len("---\nstatus: done\n---\n\n# Title\n\nSummary.\n\nFull text in Git at `" + revision + "`: `" + path + "`.\n"))
			}
		} else if plan.Action != "remove" {
			t.Fatalf("action = %q", plan.Action)
		}
		if plan.BytesAfter != wantAfter {
			t.Fatalf("after = %d, want %d", plan.BytesAfter, wantAfter)
		}
		if err := ApplyHistoryKind(root, plan); err != nil {
			t.Fatal(err)
		}
	}
	got, err := os.ReadFile(filepath.Join(root, "docs/history/adr/a.md"))
	if err != nil || string(got) != adr {
		t.Fatalf("ADR = %q, %v", got, err)
	}
}

func TestApplyHistoryKindRemovesReviewsAndHandoffsOnly(t *testing.T) {
	t.Parallel()
	root := t.TempDir()
	for _, path := range []string{"docs/history/reviews/a/nested/report.md", "docs/history/reviews/b/log.json", "docs/history/handoffs/session.md"} {
		historyEntryWrite(t, root, path, "retired")
	}
	for _, path := range []string{"docs/history/adr/a.md", "docs/references/a.md", "docs/history/specs/record.md"} {
		historyEntryWrite(t, root, path, "preserve")
	}
	plans, err := PlanHistoryKinds(root, strings.Repeat("a", 40))
	if err != nil {
		t.Fatal(err)
	}
	// A file arriving after planning is not part of the removal.
	historyEntryWrite(t, root, "docs/history/reviews/new.md", "preserve")
	for _, plan := range plans {
		if err := ApplyHistoryKind(root, plan); err != nil {
			t.Fatal(err)
		}
	}
	for _, path := range []string{"docs/history/reviews/a", "docs/history/reviews/b", "docs/history/handoffs"} {
		if _, err := os.Stat(filepath.Join(root, path)); !os.IsNotExist(err) {
			t.Fatalf("%s remains: %v", path, err)
		}
	}
	for _, path := range []string{"docs/history/adr/a.md", "docs/references/a.md", "docs/history/specs/record.md", "docs/history/reviews/new.md"} {
		got, err := os.ReadFile(filepath.Join(root, path))
		if err != nil || string(got) != "preserve" {
			t.Fatalf("%s = %q, %v", path, got, err)
		}
	}
	if err := ApplyHistoryKind(root, HistoryKindPlan{Kind: ArchiveKindADR, Action: "remove", Files: []string{"docs/history/adr/a.md"}}); err == nil {
		t.Fatal("ADR removal accepted")
	}
	if err := ApplyHistoryKind(root, HistoryKindPlan{Kind: ArchiveKindReview, Action: "remove", Files: []string{"docs/references/a.md"}}); err == nil {
		t.Fatal("out-of-kind removal accepted")
	}
}

func TestHistoryCitationsNameMarkdownThatCitesARemovedPath(t *testing.T) {
	t.Parallel()
	root := t.TempDir()
	removed := "docs/history/reviews/session"
	historyEntryWrite(t, root, "README.md", "[report](docs/history/reviews/session/report.md)\nSee `docs/history/reviews/session/log.json`.\nNo: docs/history/reviews/session-other/report.md\n")
	historyEntryWrite(t, root, "docs/history/findings/ignored.md", removed+"/report.md\n")
	historyEntryWrite(t, root, "docs/guide.md", "[report](history/reviews/session/report.md#result)\n")
	historyEntryWrite(t, root, "notes.txt", removed+"/report.md\n")
	gittest.InitRepo(t, root)
	gittest.Run(t, root, "add", ".")
	historyEntryWrite(t, root, "untracked.md", removed+"/report.md\n")
	got, err := HistoryCitations(t.Context(), root, []string{removed})
	if err != nil {
		t.Fatal(err)
	}
	want := []HistoryCitation{{Path: "README.md", Line: 1, Target: removed}, {Path: "README.md", Line: 2, Target: removed}, {Path: "docs/guide.md", Line: 1, Target: removed}}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("citations = %+v, want %+v", got, want)
	}
}
