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

func TestIsReducedHistoryEntryRecognizesLegacyProvenance(t *testing.T) {
	t.Parallel()
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

func TestPlanHistoryKindsOrdersAndMeasures(t *testing.T) {
	t.Parallel()
	root := t.TempDir()
	full := "---\nstatus: done\n---\n# Title\n\nSummary.\n\nLong discarded details.\n"
	files := []string{"docs/history/findings/z.md", "docs/history/findings/a.md", "docs/history/backlog/b.md", "docs/history/reviews/session/a.json", "docs/history/handoffs/h.md"}
	for _, path := range files {
		historyEntryWrite(t, root, path, full)
	}
	adr := full
	historyEntryWrite(t, root, "docs/history/adr/a.md", adr)
	plans, err := PlanHistoryKinds(t.Context(), root)
	if err != nil {
		t.Fatal(err)
	}
	wantFiles := [][]string{{files[1], files[0]}, {files[2]}, {"docs/history/adr/a.md"}, {files[3]}, {files[4]}}
	if len(plans) != 5 {
		t.Fatalf("plans = %+v", plans)
	}
	for i, kind := range []ArchiveKind{ArchiveKindFinding, ArchiveKindBacklog, ArchiveKindADR, ArchiveKindReview, ArchiveKindHandoff} {
		plan := plans[i]
		if plan.Kind != kind || !reflect.DeepEqual(plan.Files, wantFiles[i]) || plan.BytesBefore != int64(len(full)*len(wantFiles[i])) {
			t.Fatalf("plan %d = %+v", i, plan)
		}
		if plan.Action != "remove" {
			t.Fatalf("action = %q", plan.Action)
		}
		if err := ApplyHistoryKind(root, plan); err != nil {
			t.Fatal(err)
		}
	}
	if _, err := os.Stat(filepath.Join(root, "docs/history/adr/a.md")); !os.IsNotExist(err) {
		t.Fatalf("ADR remains: %v", err)
	}
}

func TestApplyHistoryKindRemovesEveryRetiredKind(t *testing.T) {
	t.Parallel()
	root := t.TempDir()
	for _, path := range []string{"docs/history/findings/a.md", "docs/history/backlog/a.md", "docs/history/adr/a.md", "docs/history/reviews/a/nested/report.md", "docs/history/reviews/b/log.json", "docs/history/handoffs/session.md"} {
		historyEntryWrite(t, root, path, "retired")
	}
	for _, path := range []string{"docs/references/a.md", "docs/history/specs/record.md"} {
		historyEntryWrite(t, root, path, "preserve")
	}
	plans, err := PlanHistoryKinds(t.Context(), root)
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
	for _, path := range []string{"docs/history/findings", "docs/history/backlog", "docs/history/adr", "docs/history/reviews/a", "docs/history/reviews/b", "docs/history/handoffs"} {
		if _, err := os.Stat(filepath.Join(root, path)); !os.IsNotExist(err) {
			t.Fatalf("%s remains: %v", path, err)
		}
	}
	for _, path := range []string{"docs/references/a.md", "docs/history/specs/record.md", "docs/history/reviews/new.md"} {
		got, err := os.ReadFile(filepath.Join(root, path))
		if err != nil || string(got) != "preserve" {
			t.Fatalf("%s = %q, %v", path, got, err)
		}
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
	historyEntryWrite(t, root, "docs/history/specs/0001-record.md", removed+"/report.md\n")
	historyEntryWrite(t, root, "docs/guide.md", "[report](history/reviews/session/report.md#result)\n")
	historyEntryWrite(t, root, "notes.txt", removed+"/report.md\n")
	gittest.InitRepo(t, root)
	gittest.Run(t, root, "add", ".")
	historyEntryWrite(t, root, "untracked.md", removed+"/report.md\n")
	for _, tc := range []struct {
		name   string
		target string
	}{
		{name: "canonical folder", target: removed},
		{name: "folder with trailing slash", target: removed + "/"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got, err := HistoryCitations(t.Context(), root, []string{tc.target})
			if err != nil {
				t.Fatal(err)
			}
			want := []HistoryCitation{{Path: "README.md", Line: 1, Target: tc.target}, {Path: "README.md", Line: 2, Target: tc.target}, {Path: "docs/guide.md", Line: 1, Target: tc.target}, {Path: "docs/history/specs/0001-record.md", Line: 1, Target: tc.target}}
			if !reflect.DeepEqual(got, want) {
				t.Fatalf("citations = %+v, want %+v", got, want)
			}
		})
	}
}
