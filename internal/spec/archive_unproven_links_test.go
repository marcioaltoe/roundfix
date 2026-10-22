package spec

// Suite: unproven citations crossing the archive boundary.
// Invariant: every entry keeps its target, order and non-link bytes after removal.
// Boundary IN: qualifying partial Specs and legacy folders in temporary Git repositories.
// Boundary OUT: written Archive Records, promoted copies and missing-link scanning.

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"roundfix/internal/gittest"
)

func TestArchiveRebasesTheUnprovenLinksItRecords(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		name, entry                   string
		promote, branchOnly, external bool
	}{
		{name: "removed on main", entry: "[`references/removed.md`](references/removed.md)"},
		{name: "removed only on pull request", entry: "[`references/removed.md`](references/removed.md)", branchOnly: true},
		{name: "promoted", entry: "[`references/promoted.md`](references/promoted.md)", promote: true},
		{name: "plain bytes", entry: "Operator records  the production evidence."},
		{name: "outside folder", entry: "[Repository reference](../../references/outside.md)"},
		{name: "external Spec Root", entry: "[`references/promoted.md`](references/promoted.md)", promote: true, external: true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			root := t.TempDir()
			gittest.InitRepo(t, root, "-b", "main")
			gittest.AppendConfig(t, root, "[user]\n\tname = Archive Test\n\temail = archive@example.com\n[commit]\n\tgpgsign = false\n")
			// An external Spec Root is its own repository; its record is under _archived.
			specsRoot := filepath.Join(root, "docs/specs")
			if tc.external {
				specsRoot = root
			}
			if tc.branchOnly {
				writeFile(t, filepath.Join(root, "README.md"), "Base\n")
				legacyCommit(t, root, "docs: base on main")
				gittest.Run(t, root, "checkout", "-b", "docs/archive-evidence")
			}
			plain := "Keep  this plain entry."
			entries := []string{tc.entry, plain}
			var declarations strings.Builder
			declarations.WriteString("\n## Unreachable Acceptance\n")
			for i, entry := range entries {
				fmt.Fprintf(&declarations, "\n- criterion: production evidence %d\n  reason: needs production access\n  satisfied-by: %s\n", i+1, entry)
			}
			dir := writeArchiveOverrideFixture(t, specsRoot, StatusCompleted, StatusCompleted, map[string]string{
				"qa-report-2026-10-20.md": "---\nverdict: pass\nrows_total: 1\nrows_passed: 1\n---\n\n# Earlier QA\n",
				"qa-report-2026-10-21.md": "---\nverdict: partial\nrows_total: 3\nrows_passed: 1\nrows_blocked: 2\nrows_blocked_declared: 2\nrows_blocked_finding: 0\nrows_blocked_environment: 0\n---\n\n# QA\n\n## Outcome\n\nImplementation verified; production evidence remains declared.\n",
			}, declarations.String())
			writeFile(t, filepath.Join(dir, "references/removed.md"), "Removed evidence\n")
			writeFile(t, filepath.Join(dir, "references/promoted.md"), "Promoted evidence\n")
			writeFile(t, filepath.Join(root, "docs/references/outside.md"), "Repository reference\n")
			revision := legacyCommit(t, root, "docs: qualifying partial evidence")
			req := ArchiveRequest{SpecsRoot: specsRoot, BuiltInRoot: !tc.external, Slug: "demo", RepositoryRoot: root, SourceRevision: revision}
			if tc.promote {
				req.Promote = []string{"references/promoted.md"}
			}
			if tc.branchOnly {
				req.PullRequest = "91"
			}
			result, err := Archive(req)
			if err != nil {
				t.Fatal(err)
			}
			record := readRecordResult(t, result)
			wantEntry := tc.entry
			switch {
			case tc.promote:
				destination := "../../references/promoted.md"
				if tc.external {
					destination = "../docs/references/promoted.md"
				}
				wantEntry = "[`references/promoted.md`](" + destination + ")"
			case strings.Contains(tc.entry, "references/removed.md"):
				locator := "`docs/specs/demo/references/removed.md@" + revision + "`"
				if tc.branchOnly {
					locator = "`docs/specs/demo/references/removed.md` (pull request #91)"
				}
				wantEntry = "`references/removed.md` (" + locator + ")"
			}
			if want := []string{wantEntry, plain}; !reflect.DeepEqual(record.Unproven, want) {
				t.Fatalf("unproven=%q want=%q", record.Unproven, want)
			}
			if record.Disposition != ArchivePartial || record.QAReport != "qa-report-2026-10-21.md" || record.SourceRevision != revision {
				t.Fatalf("record=%+v", record)
			}
			recordRel, err := filepath.Rel(root, result.RecordPath)
			if err != nil {
				t.Fatal(err)
			}
			wantRecord := "docs/history/specs/demo.md"
			if tc.external {
				wantRecord = "_archived/demo.md"
			}
			if filepath.ToSlash(recordRel) != wantRecord {
				t.Fatalf("record path=%s want=%s", recordRel, wantRecord)
			}
			content, err := os.ReadFile(result.RecordPath)
			if err != nil {
				t.Fatal(err)
			}
			if missing := MissingRelativeLinks(root, filepath.ToSlash(recordRel), content); len(missing) != 0 {
				t.Fatalf("broken record links: %+v", missing)
			}
		})
	}
}

func TestLegacyConversionRebasesTheUnprovenLinksItRecords(t *testing.T) {
	t.Parallel()
	req := legacyFixture(t, "status: archived\narchived: 2026-10-21\nunproven:\n  - '[`references/removed.md`](references/removed.md)'\n  - 'Keep  this plain entry.'\n", "")
	folder := filepath.Join(req.ArchiveRoot, req.Slug)
	writeFile(t, filepath.Join(folder, "references/removed.md"), "Legacy evidence\n")
	req.SourceRevision = legacyCommit(t, req.RepositoryRoot, "docs: legacy unproven evidence")
	want := []string{"`references/removed.md` (`docs/history/specs/demo/references/removed.md@" + req.SourceRevision + "`)", "Keep  this plain entry."}
	before := legacyTree(t, req)
	conversion := legacyPlan(t, req)
	legacyAssertWhole(t, req, before)
	if !reflect.DeepEqual(conversion.Record.Unproven, want) {
		t.Fatalf("planned unproven=%q want=%q", conversion.Record.Unproven, want)
	}
	if err := ApplyLegacyConversion(req.RepositoryRoot, conversion); err != nil {
		t.Fatal(err)
	}
	archived, err := ReadArchivedSpec(req.ArchiveRoot, req.Slug)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(archived.Record.Unproven, want) {
		t.Fatalf("written unproven=%q want=%q", archived.Record.Unproven, want)
	}
	content, err := os.ReadFile(filepath.Join(req.ArchiveRoot, req.Slug+".md"))
	if err != nil {
		t.Fatal(err)
	}
	if missing := MissingRelativeLinks(req.RepositoryRoot, "docs/history/specs/demo.md", content); len(missing) != 0 {
		t.Fatalf("broken legacy record links: %+v", missing)
	}
}

func TestArchiveUnprovenRebaseKeepsTheCallerSlice(t *testing.T) {
	t.Parallel()
	for _, mode := range []string{"nil rebase", "rewritten copy", "second entry error"} {
		t.Run(mode, func(t *testing.T) {
			t.Parallel()
			_, dir := archiveLinksFixture(t, true, "")
			entries := []string{"[promoted](references/promoted.md)", "[removed](references/removed.md)", "Keep  plain bytes."}
			before := append([]string(nil), entries...)
			input := ArchiveRecordInput{SpecDir: dir, Slug: "demo", Unproven: entries}
			cause := errors.New("locator unavailable")
			if mode != "nil rebase" {
				input.Rebase = &LegacyLinkRebase{
					Folder: "docs/specs/demo", RecordPath: "docs/history/specs/demo.md",
					Promoted: map[string]string{"docs/specs/demo/references/promoted.md": "docs/references/promoted.md"},
					Locate: func(path string) (GitLocator, bool, error) {
						if mode == "second entry error" {
							return GitLocator{}, false, cause
						}
						return GitLocator{Path: path, Commit: strings.Repeat("a", 40)}, true, nil
					},
				}
			}
			record, err := BuildArchiveRecord(input)
			if !reflect.DeepEqual(entries, before) {
				t.Fatalf("caller slice changed: %q want=%q", entries, before)
			}
			if mode == "second entry error" {
				if !errors.Is(err, cause) || err.Error() != "rebase archive unproven entry 2: locator unavailable" {
					t.Fatalf("error=%v", err)
				}
				return
			}
			if err != nil {
				t.Fatal(err)
			}
			want := before
			if mode == "rewritten copy" {
				want = []string{"[promoted](../../references/promoted.md)", "removed (`docs/specs/demo/references/removed.md@" + strings.Repeat("a", 40) + "`)", "Keep  plain bytes."}
			}
			if !reflect.DeepEqual(record.Unproven, want) {
				t.Fatalf("unproven=%q want=%q", record.Unproven, want)
			}
			if mode == "rewritten copy" {
				record.Unproven[2] = "Changed returned record"
				if !reflect.DeepEqual(entries, before) {
					t.Fatal("record aliases caller slice")
				}
			}
		})
	}
}
