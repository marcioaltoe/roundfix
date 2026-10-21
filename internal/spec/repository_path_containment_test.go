package spec

// Suite: resolved containment at the archive and cancellation write boundaries.
// Invariant: escaping links are refused without writes; inside links preserve bytes.
// Boundary IN: resolved temporary repositories, isolated Git config and explicit links.
// Boundary OUT: refusal text, repository status, Spec bytes and destination bytes.

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

func containmentLink(t *testing.T, target, link string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(link), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(target, link); err != nil {
		t.Fatal(err)
	}
}

func containmentFixture(t *testing.T) (ArchiveRequest, CancelRequest, string) {
	t.Helper()
	root, _ := repositoryPathFixture(t)
	specs := filepath.Join(root, "docs", "specs")
	dir := writeArchiveOverrideFixture(t, specs, StatusCompleted, StatusFailed, nil, "\nOriginal intent.\n")
	cancelAdopt(t, dir, "backlog", "entry.md", cancelBacklog)
	writeFile(t, filepath.Join(dir, "references", "notes.md"), "Shared reference bytes.\n")
	writeFile(t, filepath.Join(root, "docs", "specs", "secondary", "_prd.md"), "# Secondary\n\n[source](../demo/references/notes.md)\n[owner](../demo/_prd.md)\n")
	date := time.Date(2026, 10, 20, 0, 0, 0, 0, time.UTC)
	return ArchiveRequest{SpecsRoot: specs, RepositoryRoot: root, BuiltInRoot: true, Slug: "demo", ArchivedAt: date,
			QAOverride: &QAArchiveOverride{Approval: "fixture", Reason: "fixture"}},
		CancelRequest{SpecsRoot: specs, RepositoryRoot: root, BuiltInRoot: true, Slug: "demo", Reason: "no consumer", CancelledAt: date}, dir
}

func containmentRevision(t *testing.T, archive *ArchiveRequest, cancel *CancelRequest) {
	t.Helper()
	revision := legacyCommit(t, archive.RepositoryRoot, "docs: containment fixture")
	archive.SourceRevision = revision
	archive.QAOverride.Revision = revision
	cancel.SourceRevision = revision
}

func containmentNoWrite(t *testing.T, root, dir, outside string, run func() error, want string) {
	t.Helper()
	before := cancelTree(t, root)
	specBefore := cancelTree(t, dir)
	outsideBefore := cancelTree(t, outside)
	status := gittest.Run(t, root, "status", "--porcelain", "--untracked-files=all")
	if err := run(); err == nil || err.Error() != want {
		t.Fatalf("got %v; want %s", err, want)
	}
	if !reflect.DeepEqual(cancelTree(t, root), before) {
		t.Fatal("refusal changed repository tree")
	}
	if !reflect.DeepEqual(cancelTree(t, dir), specBefore) {
		t.Fatal("refusal changed Spec bytes")
	}
	if !reflect.DeepEqual(cancelTree(t, outside), outsideBefore) {
		t.Fatal("refusal wrote outside repository")
	}
	if got := gittest.Run(t, root, "status", "--porcelain", "--untracked-files=all"); got != status {
		t.Fatalf("status changed: %q -> %q", status, got)
	}
}

func TestRepositoryPathRefusesALinkThatLeavesTheRepository(t *testing.T) {
	t.Parallel()
	if !gitLocatorTestHome(t) {
		return
	}
	for _, kind := range []string{"existing directory", "missing file", "linked checkout"} {
		t.Run(kind, func(t *testing.T) {
			archive, cancel, dir := containmentFixture(t)
			root := archive.RepositoryRoot
			outside := filepath.Join(filepath.Dir(root), "outside")
			if err := os.MkdirAll(filepath.Join(outside, "specs"), 0o755); err != nil {
				t.Fatal(err)
			}
			containmentLink(t, outside, filepath.Join(root, "docs", "history"))
			containmentRevision(t, &archive, &cancel)
			spelling := root
			if kind == "linked checkout" {
				spelling = filepath.Join(filepath.Dir(root), "link")
			}
			path := filepath.Join(spelling, "docs", "history", "specs")
			if kind == "missing file" {
				path = filepath.Join(path, "demo.md")
			}
			containmentNoWrite(t, root, dir, outside, func() error {
				rel, err := RepositoryPath(spelling, path)
				if rel != "" {
					t.Errorf("refused path returned %q", rel)
				}
				return err
			}, fmt.Sprintf("path %q is outside the repository %q", path, spelling))
		})
	}
}

func TestRepositoryPathKeepsLinksThatStayInside(t *testing.T) {
	t.Parallel()
	if !gitLocatorTestHome(t) {
		return
	}
	root, link := repositoryPathFixture(t)
	target := filepath.Join(root, "docs", "hist-real")
	if err := os.MkdirAll(target, 0o755); err != nil {
		t.Fatal(err)
	}
	containmentLink(t, target, filepath.Join(root, "docs", "history"))
	const relative = "docs/history/specs/demo.md"
	for _, tc := range []struct{ name, root, path string }{
		{"inside parent", root, filepath.Join(root, filepath.FromSlash(relative))},
		{"both through checkout link", link, filepath.Join(link, filepath.FromSlash(relative))},
		{"linked root resolved path", link, filepath.Join(root, filepath.FromSlash(relative))},
		{"resolved root linked path", root, filepath.Join(link, filepath.FromSlash(relative))},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got, err := RepositoryPath(tc.root, tc.path)
			want := relative
			if tc.name == "linked root resolved path" || tc.name == "resolved root linked path" {
				want = "docs/hist-real/specs/demo.md"
			}
			if err != nil || got != want {
				t.Fatalf("got %q, %v; want %q", got, err, want)
			}
		})
	}
}

func TestArchiveRefusesALinkThatLeavesTheRepositoryBeforeWriting(t *testing.T) {
	t.Parallel()
	if !gitLocatorTestHome(t) {
		return
	}
	for _, family := range []string{"history", "references"} {
		t.Run(family, func(t *testing.T) {
			archive, cancel, dir := containmentFixture(t)
			root := archive.RepositoryRoot
			outside := filepath.Join(filepath.Dir(root), "outside")
			if err := os.MkdirAll(outside, 0o755); err != nil {
				t.Fatal(err)
			}
			containmentLink(t, outside, filepath.Join(root, "docs", family))
			containmentRevision(t, &archive, &cancel)
			destination := filepath.Join(root, "docs", "history", "specs")
			operation := "Archive Root for archive relink"
			if family == "references" {
				archive.Promote = []string{"references/notes.md"}
				destination = filepath.Join(root, "docs", "references", "notes.md")
				operation = "promotion destination"
			}
			containmentNoWrite(t, root, dir, outside, func() error { _, err := Archive(archive); return err },
				fmt.Sprintf("resolve %s: path %q is outside the repository %q", operation, destination, root))
		})
	}
}

func TestCancelRefusesALinkThatLeavesTheRepositoryBeforeWriting(t *testing.T) {
	t.Parallel()
	if !gitLocatorTestHome(t) {
		return
	}
	for _, family := range []string{"backlog", "findings", "references", "history"} {
		t.Run(family, func(t *testing.T) {
			archive, cancel, dir := containmentFixture(t)
			root := archive.RepositoryRoot
			if family == "findings" {
				cancelAdopt(t, dir, "finding", "finding.md", cancelFinding)
			}
			outside := filepath.Join(filepath.Dir(root), "outside")
			if err := os.MkdirAll(outside, 0o755); err != nil {
				t.Fatal(err)
			}
			containmentLink(t, outside, filepath.Join(root, "docs", family))
			containmentRevision(t, &archive, &cancel)
			basename := map[string]string{"backlog": "entry.md", "findings": "finding.md", "references": "notes.md", "history": "specs"}[family]
			destination := filepath.Join(root, "docs", family, basename)
			containmentNoWrite(t, root, dir, outside, func() error { _, err := Cancel(cancel); return err },
				fmt.Sprintf("cancellation path %q is outside the repository", destination))
		})
	}
}

func TestArchiveAndCancelKeepLinksThatStayInside(t *testing.T) {
	t.Parallel()
	if !gitLocatorTestHome(t) {
		return
	}
	for _, command := range []string{"archive", "cancel"} {
		t.Run(command, func(t *testing.T) {
			archive, cancel, _ := containmentFixture(t)
			root := archive.RepositoryRoot
			containmentRevision(t, &archive, &cancel)
			revision := archive.SourceRevision
			var controlRecord, controlLive, controlSource string
			for _, spelling := range []string{"control", "inside history link", "linked checkout"} {
				check := func(t *testing.T) {
					gittest.Run(t, root, "reset", "--hard", revision)
					gittest.Run(t, root, "clean", "-fd")
					ar, cr := archive, cancel
					if spelling == "inside history link" {
						target := filepath.Join(root, "docs", "hist-real")
						if err := os.MkdirAll(target, 0o755); err != nil {
							t.Fatal(err)
						}
						containmentLink(t, target, filepath.Join(root, "docs", "history"))
					}
					if spelling == "linked checkout" {
						link := filepath.Join(filepath.Dir(root), "link")
						ar.RepositoryRoot, cr.RepositoryRoot = link, link
						ar.SpecsRoot, cr.SpecsRoot = filepath.Join(link, "docs", "specs"), filepath.Join(link, "docs", "specs")
					}
					var recordPath string
					if command == "archive" {
						result, err := Archive(ar)
						if err != nil {
							t.Fatal(err)
						}
						recordPath = result.RecordPath
					} else {
						result, err := Cancel(cr)
						if err != nil {
							t.Fatal(err)
						}
						recordPath = result.RecordPath
					}
					record := archiveTestReadFile(t, recordPath)
					live := archiveTestReadFile(t, filepath.Join(root, "docs", "specs", "secondary", "_prd.md"))
					source := archiveTestReadFile(t, filepath.Join(root, "docs", "references", "notes.md"))
					if command == "cancel" {
						source += archiveTestReadFile(t, filepath.Join(root, "docs", "backlog", "entry.md"))
					}
					if !strings.Contains(live, "../../references/notes.md") || !strings.Contains(live, "../../history/specs/demo.md") {
						t.Fatalf("citations were not relinked: %s", live)
					}
					if spelling == "control" {
						controlRecord, controlLive, controlSource = record, live, source
						return
					}
					if record != controlRecord || live != controlLive || source != controlSource {
						t.Fatal("linked spelling changed record, relinked or copied bytes")
					}
				}
				if spelling == "control" {
					check(t)
				} else {
					t.Run(spelling, check)
				}
			}
		})
	}
}
