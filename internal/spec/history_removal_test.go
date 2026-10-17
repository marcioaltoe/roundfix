package spec

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"roundfix/internal/gittest"
)

func TestPlanHistoryKindsRemovesEveryRetiredFamily(t *testing.T) {
	t.Parallel()
	if !gitLocatorTestHome(t) {
		return
	}
	root, _ := gitLocatorRepo(t)
	for _, family := range []string{"findings", "backlog", "adr", "reviews", "handoffs"} {
		historyEntryWrite(t, root, "docs/history/"+family+"/entry.md", "legacy bytes\n")
	}
	historyEntryWrite(t, root, "docs/history/findings/reduced.md", "# Reduced\n\nFull text in Git at `"+strings.Repeat("a", 40)+"`: `docs/findings/entry.md`.\n")
	plans, err := PlanHistoryKinds(t.Context(), root)
	if err != nil || len(plans) != 5 {
		t.Fatalf("plans=%+v err=%v", plans, err)
	}
	for _, p := range plans {
		if p.Action != "remove" {
			t.Fatalf("plan=%+v", p)
		}
		for _, path := range p.Files {
			if _, err := os.Stat(filepath.Join(root, path)); err != nil {
				t.Fatal(err)
			}
		}
	}
	if len(plans[0].Files) != 2 {
		t.Fatalf("reduced entry not pending: %+v", plans[0])
	}
}

func TestPlanHistoryKindsRemovesWhatTheOldLayoutLeft(t *testing.T) {
	t.Parallel()
	if !gitLocatorTestHome(t) {
		return
	}
	root, _ := gitLocatorRepo(t)
	files := []string{"docs/findings/_archived/a.md", "_archived/findings/a.md", "_archived/backlog/a.md", "_archived/adr/a.md", "docs/findings/done.md", "docs/backlog/deferred.md", "docs/adr/0001-old.md", "docs/specs/_reviews/one/outcome.md", "docs/specs/reviews/two/outcome.md"}
	for _, path := range files {
		historyEntryWrite(t, root, path, "---\nstatus: deferred\n---\n# Retired\n")
	}
	historyEntryWrite(t, root, "docs/adr/0001-old.md", "---\nstatus: deprecated\n---\n# Old\n")
	for _, path := range files[7:] {
		historyEntryWrite(t, root, path, "---\npull_request_state: closed\nrecorded_at: 2026-10-16\n---\n# Finished\n")
	}
	for _, path := range []string{"docs/findings/pending.md", "docs/findings/unknown.md", "docs/backlog/open.md", "docs/backlog/promoted.md", "docs/adr/0002-active.md", "docs/specs/reviews/live/review.md"} {
		historyEntryWrite(t, root, path, "# Active or undecidable\n")
	}
	plans, err := PlanHistoryKinds(t.Context(), root)
	if err != nil {
		t.Fatal(err)
	}
	got := map[string]bool{}
	for _, p := range plans {
		for _, path := range p.Files {
			got[path] = true
		}
	}
	if len(got) != len(files) {
		t.Fatalf("files=%v", got)
	}
	for _, path := range files {
		if !got[path] {
			t.Fatalf("missing %s", path)
		}
	}
}

func TestPlanHistoryKindsRefusesASupersededADRItsSuccessorDoesNotName(t *testing.T) {
	t.Parallel()
	if !gitLocatorTestHome(t) {
		return
	}
	root, _ := gitLocatorRepo(t)
	old := "docs/adr/0001-old.md"
	historyEntryWrite(t, root, old, "---\nstatus: superseded\nsuperseded_by: ADR-0002\n---\n# Old\n")
	historyEntryWrite(t, root, "docs/adr/0002-new.md", "---\nstatus: accepted\n---\n# New\n")
	if _, err := PlanHistoryKinds(t.Context(), root); err == nil || !strings.Contains(err.Error(), "0001-old.md") {
		t.Fatalf("err=%v", err)
	}
	historyEntryWrite(t, root, "docs/adr/0002-new.md", "---\nstatus: accepted\n---\n# New\n\nConsolidates ADR-0001 (2026-10-16).\n")
	plans, err := PlanHistoryKinds(t.Context(), root)
	if err != nil || len(plans) != 1 || !reflect.DeepEqual(plans[0].Files, []string{old}) {
		t.Fatalf("plans=%+v err=%v", plans, err)
	}
}

func TestApplyHistoryKindKeepsTheActiveDirectories(t *testing.T) {
	t.Parallel()
	if !gitLocatorTestHome(t) {
		return
	}
	root, _ := gitLocatorRepo(t)
	for _, p := range []string{"docs/findings/a.md", "docs/backlog/a.md"} {
		historyEntryWrite(t, root, p, "---\nstatus: done\n---\n# Retired\n")
	}
	historyEntryWrite(t, root, "docs/adr/0001-old.md", "---\nstatus: rejected\n---\n# Old\n")
	historyEntryWrite(t, root, "docs/specs/reviews/one/outcome.md", "---\npull_request_state: closed\nrecorded_at: 2026-10-16\n---\n")
	plans, err := PlanHistoryKinds(t.Context(), root)
	if err != nil {
		t.Fatal(err)
	}
	for _, p := range plans {
		if err := ApplyHistoryKind(root, p); err != nil {
			t.Fatal(err)
		}
	}
	for _, dir := range []string{"docs/adr", "docs/findings", "docs/backlog", "docs/specs"} {
		info, err := os.Stat(filepath.Join(root, dir))
		if err != nil || !info.IsDir() {
			t.Fatalf("directory %s: %v", dir, err)
		}
	}
	again, err := PlanHistoryKinds(t.Context(), root)
	if err != nil || len(again) != 0 {
		t.Fatalf("second plan=%+v err=%v", again, err)
	}
	// Whole-list validation precedes the first deletion, including symlink ancestors.
	first := "docs/history/findings/a.md"
	historyEntryWrite(t, root, first, "keep\n")
	p := HistoryKindPlan{Kind: ArchiveKindFinding, Action: "remove", Files: []string{first, "docs/references/keep.md"}}
	if err := ApplyHistoryKind(root, p); err == nil {
		t.Fatal("unsafe plan accepted")
	}
	if _, err := os.Stat(filepath.Join(root, first)); err != nil {
		t.Fatal("removed before validation")
	}
}

func historyRemovalRecord(t *testing.T, root, slug string, granted bool) string {
	t.Helper()
	r := ArchiveRecord{Spec: slug, Title: slug, Source: "docs/specs/" + slug, SourceRevision: strings.Repeat("a", 40), Disposition: ArchiveSuperseded, SupersededBy: "0099-successor", Outcome: "Retired."}
	if granted {
		r.Regeneration = []ArchiveRegeneration{{Command: "generate", Outputs: []string{"out.json"}}}
	}
	content, err := RenderArchiveRecord(r)
	if err != nil {
		t.Fatal(err)
	}
	path := "docs/history/specs/" + slug + ".md"
	historyEntryWrite(t, root, path, string(content))
	return path
}

func TestPlanSupersededRecordsSkipsCitedAndGrantedRecords(t *testing.T) {
	t.Parallel()
	if !gitLocatorTestHome(t) {
		return
	}
	root, _ := gitLocatorRepo(t)
	uncited := historyRemovalRecord(t, root, "0001-uncited", false)
	for _, slug := range []string{"0002-path", "0003-slug", "0004-link", "0005-granted", "0006-removed-source", "0007-skill", "0008-fixture", "0009-locator", "0010-unreachable", "0011-pull-request"} {
		historyRemovalRecord(t, root, slug, slug == "0005-granted")
	}
	historyEntryWrite(t, root, "unusual/document.md", "spec: 0003-slug\n[record](../docs/history/specs/0004-link.md)\n`docs/history/specs/0002-path.md`\n`docs/history/specs/0009-locator.md@"+strings.Repeat("b", 40)+"`\n")
	removed := "docs/history/findings/source.md"
	historyEntryWrite(t, root, "other/locators.md", "`docs/history/specs/0010-unreachable.md` (not reachable from any branch or tag)\n`docs/history/specs/0011-pull-request.md` (pull request #42)\n")
	historyEntryWrite(t, root, removed, "absorbed_by: 0006-removed-source\n")
	historyEntryWrite(t, root, "skills/example/SKILL.md", "0007-skill\n")
	historyEntryWrite(t, root, "nested/testdata/fixture.md", "0008-fixture\n")
	legacyCommit(t, root, "docs: superseded records")
	p, err := PlanSupersededRecords(t.Context(), root, "docs/history/specs", []string{removed})
	want := []string{uncited, "docs/history/specs/0006-removed-source.md", "docs/history/specs/0007-skill.md", "docs/history/specs/0008-fixture.md", "docs/history/specs/0009-locator.md", "docs/history/specs/0010-unreachable.md", "docs/history/specs/0011-pull-request.md"}
	if err != nil || !reflect.DeepEqual(p.Files, want) {
		t.Fatalf("plan=%+v err=%v want=%v", p, err, want)
	}
}

func TestRemovalRewriteLocatesAReducedEntryAtItsProvenance(t *testing.T) {
	t.Parallel()
	if !gitLocatorTestHome(t) {
		return
	}
	root, _ := gitLocatorRepo(t)
	active := "docs/findings/original.md"
	historyEntryWrite(t, root, active, "# Original full text\n")
	provenance := legacyCommit(t, root, "docs: original finding")
	gittest.Run(t, root, "rm", active)
	retired := "docs/history/findings/reduced.md"
	historyEntryWrite(t, root, retired, fmt.Sprintf("# Summary\n\nFull text in Git at `%s`: `%s`.\n", provenance, active))
	guide := "docs/user-guide/a.md"
	historyEntryWrite(t, root, guide, "# Guide\n\n[finding](../history/findings/reduced.md)\n")
	historyEntryWrite(t, root, "README.md", retired+"\n")
	legacyCommit(t, root, "docs: retired finding")
	plan, err := PlanHistoryRewrite(t.Context(), root, "docs/specs", "docs/history/specs", []string{retired}, nil)
	if err != nil {
		t.Fatal(err)
	}
	want := "# Guide\n\nfinding (`" + active + "@" + provenance + "`)\n"
	if string(plan.Files[guide]) != want || len(plan.Files) != 1 || plan.Count != 1 {
		t.Fatalf("plan=%+v want=%s", plan, want)
	}
	count, files, err := plan.Apply(t.Context(), root)
	if err != nil || count != 1 || files != 1 {
		t.Fatalf("apply=%d %d %v", count, files, err)
	}
	got, err := os.ReadFile(filepath.Join(root, guide))
	if err != nil || string(got) != want {
		t.Fatalf("guide=%s err=%v", got, err)
	}
}
