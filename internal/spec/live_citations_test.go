package spec

// Suite: Archive Relink planning and exact byte rewrites.
// Invariant: only Live Document citations move; promoted bytes keep their source.
// Boundary IN: Markdown bytes and disposable Git repositories.
// Boundary OUT: command streams are covered in internal/cli.

import (
	"bytes"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"roundfix/internal/gittest"
)

func TestIsLiveDocumentClassifiesLiveFamilies(t *testing.T) {
	t.Parallel()
	for _, tt := range []struct {
		path, root, archive string
		live                bool
	}{
		{"docs/specs/secondary/_prd.md", "docs/specs", "docs/history/specs", true},
		{"docs/backlog/2026-01-02-follow-up.md", "docs/specs", "docs/history/specs", true},
		{"docs/findings/2026-01-02-gap.md", "docs/specs", "docs/history/specs", true},
		{"docs/references/source.md", "docs/specs", "docs/history/specs", true},
		{"docs/handoffs/notes.md", "docs/specs", "docs/history/specs", false},
		{"README.md", "docs/specs", "docs/history/specs", false},
		{"docs/user-guide/archive.md", "docs/specs", "docs/history/specs", false},
		{"docs/history/specs/owner/_prd.md", "docs", "docs/history/specs", false},
		{"plans/_archived/owner.md", "plans", "plans/_archived", false},
		{"plans/secondary/_prd.md", "plans", "plans/_archived", true},
		{"secondary/_prd.md", ".", "_archived", true},
		{"_archived/owner.md", ".", "_archived", false},
		{"docs/specs-near/file.md", "docs/specs", "docs/history/specs", false},
		{"docs/specs/source.json", "docs/specs", "docs/history/specs", false},
	} {
		t.Run(tt.path, func(t *testing.T) {
			if got := IsLiveDocument(tt.path, tt.root, tt.archive); got != tt.live {
				t.Fatalf("live=%v want=%v", got, tt.live)
			}
		})
	}
}

func TestRelinkContentRewritesOnlyDestinationsAndPaths(t *testing.T) {
	t.Parallel()
	const folder = "docs/specs/owner"
	moves := ArchiveCitationMoves(folder, "docs/history/specs/owner.md", []string{folder + "/_prd.md", folder + "/references/source.md"}, map[string]string{folder + "/references/source.md": "docs/references/source.md"})
	for _, tt := range []struct {
		name, before, after string
		count               int
	}{
		{"inline", `[source](../owner/references/source.md "Title")`, `[source](../../references/source.md "Title")`, 1},
		{"angle", `[source](<../owner/references/source.md>)`, `[source](<../../references/source.md>)`, 1},
		{"definition", "[source]: ../owner/references/source.md 'Title'\r\n", "[source]: ../../references/source.md 'Title'\r\n", 1},
		{"fragment and query", `[PRD](../owner/_prd.md?view=1#goals)`, `[PRD](../../history/specs/owner.md?view=1#goals)`, 1},
		{"folder slash", `[owner](../owner/)`, `[owner](../../history/specs/owner.md)`, 1},
		{"prose", "See docs/specs/owner/references/source.md, then docs/specs/owner/_prd.md.", "See docs/references/source.md, then docs/history/specs/owner.md.", 2},
		{"code span", "`docs/specs/owner/references/source.md`", "`docs/references/source.md`", 1},
		{"fence", "```sh\ncat docs/specs/owner/references/source.md\n```\n", "```sh\ncat docs/references/source.md\n```\n", 1},
		{"near prefix", "docs/specs/owner/references/source.md-extra", "docs/specs/owner/references/source.md-extra", 0},
		{"absolute URL", `[web](https://example.com/docs/specs/owner/_prd.md)`, `[web](https://example.com/docs/specs/owner/_prd.md)`, 0},
		{"unrelated", `[other](../../references/unrelated.md)`, `[other](../../references/unrelated.md)`, 0},
		{"label remains", `[docs/specs/owner/_prd.md](../owner/_prd.md)`, `[docs/history/specs/owner.md](../../history/specs/owner.md)`, 2},
	} {
		t.Run(tt.name, func(t *testing.T) {
			got, count := RelinkContent([]byte(tt.before), "docs/specs/secondary/_prd.md", moves)
			if string(got) != tt.after || count != tt.count {
				t.Fatalf("got=%q count=%d want=%q count=%d", got, count, tt.after, tt.count)
			}
		})
	}
	t.Run("escaped destination", func(t *testing.T) {
		moves := map[string]string{folder + "/references/a b.md": "docs/references/a b.md"}
		got, count := RelinkContent([]byte(`[source](../owner/references/a%20b.md#part)`), "docs/specs/secondary/_prd.md", moves)
		if string(got) != `[source](../../references/a%20b.md#part)` || count != 1 {
			t.Fatalf("got=%q count=%d", got, count)
		}
	})
}

func TestPlanArchiveRelinkPromotesLiveCitedFilesAndTargetsTheRecord(t *testing.T) {
	t.Parallel()
	root := t.TempDir()
	gittest.InitRepo(t, root)
	const folder = "docs/specs/owner"
	files := map[string]string{
		folder + "/_prd.md":              "# Owner\n",
		folder + "/references/source.md": "source\n",
		folder + "/references/manual.md": "manual\n",
		"docs/specs/secondary/_prd.md":   "[source](../owner/references/source.md)\n[PRD](../owner/_prd.md)\n",
		"docs/handoffs/notes.md":         "[source](../specs/owner/references/source.md)\n",
		"docs/history/specs/old.md":      "docs/specs/owner/references/source.md\n",
	}
	for path, content := range files {
		writeFile(t, filepath.Join(root, path), content)
	}
	gittest.Run(t, root, "add", "-A")
	gittest.Run(t, root, "commit", "-m", "docs: citation fixture")
	plan, err := PlanArchiveRelink(t.Context(), root, "docs/specs", "docs/history/specs", folder, "docs/history/specs/owner.md", []string{"references/manual.md"})
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(plan.Promote, []string{"references/manual.md", "references/source.md"}) {
		t.Fatalf("promote=%v", plan.Promote)
	}
	want := "[source](../../references/source.md)\n[PRD](../../history/specs/owner.md)\n"
	if len(plan.Relinked) != 1 || string(plan.Relinked["docs/specs/secondary/_prd.md"]) != want {
		t.Fatalf("relinked=%v", plan.Relinked)
	}
	if !reflect.DeepEqual(plan.Unrelinked, []HistoryCitation{{Path: "docs/handoffs/notes.md", Line: 1, Target: folder + "/references/source.md"}}) {
		t.Fatalf("unrelinked=%v", plan.Unrelinked)
	}
	for path, content := range files {
		got, err := os.ReadFile(filepath.Join(root, path))
		if err != nil || string(got) != content {
			t.Fatalf("plan changed %s: %q %v", path, got, err)
		}
	}
}

func TestPlanArchiveRelinkIsEmptyOutsideAGitWorkTree(t *testing.T) {
	t.Parallel()
	plan, err := PlanArchiveRelink(t.Context(), t.TempDir(), "docs/specs", "docs/history/specs", "docs/specs/owner", "docs/history/specs/owner.md", []string{"references/source.md"})
	if err != nil || len(plan.Promote)+len(plan.Relinked)+len(plan.Unrelinked) != 0 {
		t.Fatalf("plan=%+v err=%v", plan, err)
	}
}

func TestArchiveRecordRendersRelinkedOnlyWhenPresent(t *testing.T) {
	t.Parallel()
	r := ArchiveRecord{Spec: "owner", Title: "Owner", Created: "2026-01-01", Archived: "2026-01-02", Disposition: ArchivePass, Source: "docs/specs/owner", SourceRevision: strings.Repeat("a", 40), QATask: "task_qa", QAReport: "qa-report-2026-01-01.md", QAVerdict: "pass", Outcome: "Archived owner."}
	const legacy = "---\nschema: roundfix/archive-record/v1\nspec: owner\ntitle: Owner\nstatus: archived\ncreated: \"2026-01-01\"\narchived: \"2026-01-02\"\ndisposition: pass\nsource: docs/specs/owner\nsource_revision: aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa\nqa_task: task_qa\nqa_report: qa-report-2026-01-01.md\nqa_verdict: pass\nunproven: []\nadrs: []\nsources: []\nregeneration: []\npromoted: []\n---\n\n# Owner\n\nArchived owner.\n"
	before, err := RenderArchiveRecord(r)
	if err != nil || string(before) != legacy {
		t.Fatalf("legacy bytes changed:\n%s\nerr=%v", before, err)
	}
	r.Relinked = []string{"docs/backlog/follow-up.md", "docs/specs/secondary/_prd.md"}
	content, err := RenderArchiveRecord(r)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Contains(content, []byte("promoted: []\nrelinked:\n")) {
		t.Fatalf("field order: %s", content)
	}
	got, err := ParseArchiveRecord(content)
	if err != nil || !reflect.DeepEqual(got.Relinked, r.Relinked) {
		t.Fatalf("roundtrip=%+v err=%v", got, err)
	}
	parsed, err := ParseArchiveRecord(before)
	if err != nil || len(parsed.Relinked) != 0 {
		t.Fatalf("legacy parse=%+v err=%v", parsed, err)
	}
}

func TestPlanArchiveRelinkRecognizesScannerFormsAndExcludesArchiveRoots(t *testing.T) {
	t.Parallel()
	root := t.TempDir()
	gittest.InitRepo(t, root)
	files := map[string]string{
		"plans/owner/_prd.md":           "# Owner\n",
		"plans/owner/references/a b.md": "source\n",
		"plans/secondary/_prd.md":       "[source](<../owner/references/a%20b.md>)\n[source]: ../owner/references/a%20b.md#part\n",
		"plans/_archived/old.md":        "[source](../owner/references/a%20b.md)\n",
		"docs/history/backlog/old.md":   "plans/owner/references/a b.md\n",
		"README.md":                     "[owner](plans/owner/)\n",
	}
	for path, content := range files {
		writeFile(t, filepath.Join(root, path), content)
	}
	gittest.Run(t, root, "add", "-A")
	gittest.Run(t, root, "commit", "-m", "docs: scanner fixture")
	plan, err := PlanArchiveRelink(t.Context(), root, "plans", "plans/_archived", "plans/owner", "plans/_archived/owner.md", nil)
	if err != nil {
		t.Fatal(err)
	}
	want := "[source](<../../docs/references/a%20b.md>)\n[source]: ../../docs/references/a%20b.md#part\n"
	if !reflect.DeepEqual(plan.Promote, []string{"references/a b.md"}) || len(plan.Relinked) != 1 || string(plan.Relinked["plans/secondary/_prd.md"]) != want {
		t.Fatalf("plan=%+v", plan)
	}
	if !reflect.DeepEqual(plan.Unrelinked, []HistoryCitation{{Path: "README.md", Line: 1, Target: "plans/owner"}}) {
		t.Fatalf("unrelinked=%+v", plan.Unrelinked)
	}
}

func TestArchiveRelinkRestoresFilesAfterALaterWriteFailure(t *testing.T) {
	t.Parallel()
	req, dir := recordFixture(t)
	const first = "docs/backlog/2026-01-01-first.md"
	const second = "docs/specs/secondary/_prd.md"
	for _, path := range []string{first, second} {
		writeFile(t, filepath.Join(req.RepositoryRoot, path), "Read `docs/specs/demo/references/adopted.md`.\n")
	}
	req.SourceRevision = archiveFixtureRevision(t, req.SpecsRoot)
	path := filepath.Join(req.RepositoryRoot, second)
	if err := os.Chmod(path, 0o444); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if err := os.Chmod(path, 0o644); err != nil {
			t.Error(err)
		}
	})
	before := archiveLinksTree(t, req.RepositoryRoot)
	_, err := Archive(req)
	if err == nil || !strings.Contains(err.Error(), "write archive relink") {
		t.Fatalf("expected a later relink write failure, got %v", err)
	}
	after := archiveLinksTree(t, req.RepositoryRoot)
	if !reflect.DeepEqual(before, after) {
		for path, content := range after {
			if before[path] != content {
				t.Errorf("changed or created: %s", path)
			}
		}
		for path := range before {
			if _, ok := after[path]; !ok {
				t.Errorf("removed: %s", path)
			}
		}
		t.Fatal("rollback failed to restore files or remove created files")
	}
	if _, err := os.Stat(dir); err != nil {
		t.Fatalf("Spec removed on write failure: %v", err)
	}
}

func TestPlanArchiveRelinkRefusesGitFailureInAWorkTree(t *testing.T) {
	t.Parallel()
	root := t.TempDir()
	// A broken worktree pointer must refuse, rather than fall back to non-Git behavior.
	writeFile(t, filepath.Join(root, ".git"), "gitdir: missing\n")
	plan, err := PlanArchiveRelink(t.Context(), root, "docs/specs", "docs/history/specs", "docs/specs/owner", "docs/history/specs/owner.md", nil)
	if err == nil || len(plan.Promote)+len(plan.Relinked)+len(plan.Unrelinked) != 0 {
		t.Fatalf("plan=%+v err=%v", plan, err)
	}
}
