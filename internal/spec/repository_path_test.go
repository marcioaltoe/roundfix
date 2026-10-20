package spec

// Suite: repository path spellings at archive and cancellation boundaries.
// Invariant: linked spellings preserve record, citation and returned-source bytes.
// Boundary IN: resolved temporary repositories and explicit symbolic links.
// Boundary OUT: repository-relative paths and files written by Archive and Cancel.

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"

	"roundfix/internal/gittest"
)

func repositoryPathFixture(t *testing.T) (string, string) {
	t.Helper()
	base, err := filepath.EvalSymlinks(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	root := filepath.Join(base, "repo")
	gittest.InitRepo(t, root, "-b", "main")
	link := filepath.Join(base, "link")
	if err := os.Symlink(root, link); err != nil {
		t.Fatal(err)
	}
	return root, link
}

func TestArchiveFromALinkedSpellingRewritesAndRelinks(t *testing.T) {
	t.Parallel()
	if !gitLocatorTestHome(t) {
		return
	}
	root, link := repositoryPathFixture(t)
	specs := filepath.Join(root, "docs", "specs")
	writeArchiveOverrideFixture(t, specs, StatusCompleted, StatusFailed, nil,
		"\nSee [follow-up](../../backlog/follow-up.md).\n")
	const live = "docs/backlog/follow-up.md"
	writeFile(t, filepath.Join(root, live), "# Follow-up\n\n[owner](../specs/demo/_prd.md)\n")
	revision := legacyCommit(t, root, "docs: linked archive fixture")
	var record, backlog string
	var relinked []string
	for _, spelling := range []string{"control", "linked Spec Root", "linked repository root"} {
		t.Run(spelling, func(t *testing.T) {
			// Restore only this disposable fixture, keeping the same committed provenance.
			gittest.Run(t, root, "reset", "--hard", revision)
			gittest.Run(t, root, "clean", "-fd")
			req := ArchiveRequest{SpecsRoot: specs, RepositoryRoot: root, BuiltInRoot: true,
				Slug: "demo", SourceRevision: revision, ArchivedAt: time.Date(2026, 10, 20, 0, 0, 0, 0, time.UTC),
				QAOverride: &QAArchiveOverride{Approval: "fixture", Reason: "fixture", Revision: revision}}
			if spelling == "linked Spec Root" {
				req.SpecsRoot = filepath.Join(link, "docs", "specs")
			}
			if spelling == "linked repository root" {
				req.RepositoryRoot = link
			}
			result, err := Archive(req)
			if err != nil {
				t.Fatal(err)
			}
			gotRecord := archiveTestReadFile(t, result.RecordPath)
			gotBacklog := archiveTestReadFile(t, filepath.Join(root, live))
			if !strings.Contains(gotRecord, "[follow-up](../../backlog/follow-up.md)") {
				t.Fatalf("record did not rebase outward citation: %s", gotRecord)
			}
			if gotBacklog != "# Follow-up\n\n[owner](../history/specs/demo.md)\n" || !reflect.DeepEqual(result.Relinked, []string{live}) {
				t.Fatalf("backlog=%q relinked=%v", gotBacklog, result.Relinked)
			}
			if spelling == "control" {
				record, backlog, relinked = gotRecord, gotBacklog, result.Relinked
				return
			}
			if gotRecord != record || gotBacklog != backlog || !reflect.DeepEqual(result.Relinked, relinked) {
				t.Fatalf("linked spelling changed archive bytes or relinked paths: %s", spelling)
			}
		})
	}
}

func TestCancelFromALinkedSpellingRelinksAndReturnsSources(t *testing.T) {
	t.Parallel()
	if !gitLocatorTestHome(t) {
		return
	}
	root, link := repositoryPathFixture(t)
	specs := filepath.Join(root, "docs", "specs")
	dir := filepath.Join(specs, "demo")
	writeFile(t, filepath.Join(dir, "_prd.md"), prdFixture("active")+"\nOriginal intent.\n")
	cancelAdopt(t, dir, "backlog", "entry.md", cancelBacklog)
	const live = "docs/backlog/follow-up.md"
	writeFile(t, filepath.Join(root, live), "# Follow-up\n\n[owner](../specs/demo/_prd.md)\n")
	revision := legacyCommit(t, root, "docs: linked cancellation fixture")
	var record, backlog, source string
	var returned []ReturnedSource
	for _, spelling := range []string{"control", "linked Spec Root", "linked repository root"} {
		t.Run(spelling, func(t *testing.T) {
			gittest.Run(t, root, "reset", "--hard", revision)
			gittest.Run(t, root, "clean", "-fd")
			req := CancelRequest{SpecsRoot: specs, RepositoryRoot: root, BuiltInRoot: true,
				Slug: "demo", Reason: "no consumer", SourceRevision: revision,
				CancelledAt: time.Date(2026, 10, 20, 0, 0, 0, 0, time.UTC)}
			if spelling == "linked Spec Root" {
				req.SpecsRoot = filepath.Join(link, "docs", "specs")
			}
			if spelling == "linked repository root" {
				req.RepositoryRoot = link
			}
			result, err := Cancel(req)
			if err != nil {
				t.Fatal(err)
			}
			gotRecord := archiveTestReadFile(t, result.RecordPath)
			gotBacklog := archiveTestReadFile(t, filepath.Join(root, live))
			gotSource := archiveTestReadFile(t, filepath.Join(root, "docs/backlog/entry.md"))
			wantSource := strings.Replace(cancelBacklog, "status: promoted", "status: open", 1)
			wantSource = strings.Replace(wantSource, "spec: demo", "spec: null", 1)
			if gotSource != wantSource || !reflect.DeepEqual(result.Returned, []ReturnedSource{{"docs/specs/demo/references/entry.md", "docs/backlog/entry.md", "open"}}) {
				t.Fatalf("returned source=%q paths=%+v", gotSource, result.Returned)
			}
			for _, item := range result.Returned {
				assertRepositoryRelative(t, item.From)
				assertRepositoryRelative(t, item.To)
			}
			if gotBacklog != "# Follow-up\n\n[owner](../history/specs/demo.md)\n" {
				t.Fatalf("live citation=%q", gotBacklog)
			}
			if spelling == "control" {
				record, backlog, source, returned = gotRecord, gotBacklog, gotSource, result.Returned
				return
			}
			if gotRecord != record || gotBacklog != backlog || gotSource != source || !reflect.DeepEqual(result.Returned, returned) {
				t.Fatalf("linked spelling changed cancellation bytes or returned paths: %s", spelling)
			}
		})
	}
}

func assertRepositoryRelative(t *testing.T, path string) {
	t.Helper()
	if path == ".." || strings.HasPrefix(path, "../") || filepath.IsAbs(path) || strings.Contains(path, "\\") {
		t.Fatalf("unsafe repository-relative result %q", path)
	}
}

func TestRepositoryPathComparesResolvedSpellings(t *testing.T) {
	t.Parallel()
	if !gitLocatorTestHome(t) {
		return
	}
	root, link := repositoryPathFixture(t)
	for _, tc := range []struct{ name, root, path, want string }{
		{"one spelling without existing paths", filepath.Join(root, "absent"), filepath.Join(root, "absent", "file.md"), "file.md"},
		{"root resolved path linked", root, filepath.Join(link, "docs"), "docs"},
		{"root linked path resolved", link, filepath.Join(root, "docs"), "docs"},
		{"nonexistent record below linked parent", root, filepath.Join(link, "docs", "history", "specs", "demo.md"), "docs/history/specs/demo.md"},
		{"root itself", root, root, "."},
		{"same root different spellings", root, link, "."},
		{"sibling sharing prefix", root, root + "2", ""},
		{"outside", link, filepath.Join(filepath.Dir(root), "outside"), ""},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got, err := RepositoryPath(tc.root, tc.path)
			if tc.want == "" {
				want := fmt.Sprintf("path %q is outside the repository %q", tc.path, tc.root)
				if err == nil || err.Error() != want || got != "" {
					t.Fatalf("got %q, %v; want %s", got, err, want)
				}
				return
			}
			if err != nil || got != tc.want {
				t.Fatalf("got %q, %v; want %q", got, err, tc.want)
			}
			assertRepositoryRelative(t, got)
		})
	}
}
