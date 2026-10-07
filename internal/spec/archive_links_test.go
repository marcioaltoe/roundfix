package spec

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func archiveLinksFixture(t *testing.T, builtIn bool, body string) (ArchiveRequest, string) {
	t.Helper()
	root := filepath.Join(t.TempDir(), "docs", "specs")
	dir := writeArchiveOverrideFixture(t, root, StatusCompleted, StatusCompleted, map[string]string{"qa-report-2026-10-04.md": "---\nverdict: pass\nrows_total: 1\nrows_passed: 1\nrows_blocked: 0\n---\n"}, body)
	return ArchiveRequest{SpecsRoot: root, BuiltInRoot: builtIn, Slug: "demo"}, dir
}

func TestArchiveRewritesRelativeLinksThatLeaveTheSpec(t *testing.T) {
	req, dir := archiveLinksFixture(t, true, "\n[ADR](../../adr/decision.md#choice)\n![image](<../../adr/a%20b.png?raw=1#view> \"title\")\n[ref]: ../../adr/decision.md 'title'\n")
	writeFile(t, filepath.Join(req.SpecsRoot, "..", "adr", "decision.md"), "decision")
	writeFile(t, filepath.Join(req.SpecsRoot, "..", "adr", "a b.png"), "image")
	writeFile(t, filepath.Join(dir, "nested", "notes.md"), "[nested](../../../adr/decision.md)\n")
	result, err := legacyLinkArchive(t, req)
	if err != nil {
		t.Fatal(err)
	}
	if result.RewrittenLinks != 4 {
		t.Fatalf("rewrote %d, want 4", result.RewrittenLinks)
	}
	prd := archiveTestReadFile(t, filepath.Join(result.ArchivedDir, "_prd.md"))
	for _, want := range []string{"../../../adr/decision.md#choice", "<../../../adr/a%20b.png?raw=1#view> \"title\"", "[ref]: ../../../adr/decision.md 'title'"} {
		if !strings.Contains(prd, want) {
			t.Fatalf("missing %q in %s", want, prd)
		}
	}
	if got := archiveTestReadFile(t, filepath.Join(result.ArchivedDir, "nested", "notes.md")); got != "[nested](../../../../adr/decision.md)\n" {
		t.Fatal(got)
	}
}

func TestArchiveKeepsLinksInsideTheSpecAndNonRelativeLinks(t *testing.T) {
	req, dir := archiveLinksFixture(t, true, "")
	body := "[inside](missing.md) [fragment](#a) [url](https://example.com/a) [absolute](/missing)\n`[code](../../missing)`\n```md\n[x](../../missing)\n```\n~~~\n[x](../../missing)\n~~~\n<a href=\"../../missing\">HTML</a>\n"
	writeFile(t, filepath.Join(dir, "notes.md"), body)
	writeFile(t, filepath.Join(dir, "notes.txt"), "[text](../../missing)")
	if err := os.Symlink("missing.md", filepath.Join(dir, "symlink.md")); err != nil {
		t.Fatal(err)
	}
	result, err := legacyLinkArchive(t, req)
	if err != nil {
		t.Fatal(err)
	}
	if result.RewrittenLinks != 0 {
		t.Fatal(result.RewrittenLinks)
	}
	if got := archiveTestReadFile(t, filepath.Join(result.ArchivedDir, "notes.md")); got != body {
		t.Fatal(got)
	}
	if got := archiveTestReadFile(t, filepath.Join(result.ArchivedDir, "notes.txt")); got != "[text](../../missing)" {
		t.Fatal(got)
	}
}

func TestArchiveRefusesALinkThatWouldStayBroken(t *testing.T) {
	req, dir := archiveLinksFixture(t, true, "\n[broken](../../adr/missing.md)\n")
	writeFile(t, filepath.Join(dir, "nested", "notes.md"), "[other](../../../findings/missing.md)\n")
	writeFile(t, filepath.Join(req.SpecsRoot, "..", "adr", "decision.md"), "target")
	writeFile(t, filepath.Join(dir, "rewritable.md"), "[valid](../../adr/decision.md)\n")
	before := archiveLinksTree(t, dir)
	_, err := legacyLinkArchive(t, req)
	if err == nil || !strings.Contains(err.Error(), `Spec "demo" has relative links that leave the Spec and do not resolve: _prd.md:`) || !strings.Contains(err.Error(), `"../../adr/missing.md"`) {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := archiveLinksTree(t, dir); !reflect.DeepEqual(got, before) {
		t.Fatal("Spec tree changed")
	}
	if !strings.HasSuffix(err.Error(), "; fix or remove each link, then retry the archive") {
		t.Fatal(err)
	}
	if !strings.Contains(err.Error(), `nested/notes.md:1 "../../../findings/missing.md"`) {
		t.Fatalf("missing second broken link: %v", err)
	}
	if _, err := os.Stat(filepath.Join(ArchiveSpecRoot(req.SpecsRoot, true), req.Slug)); !os.IsNotExist(err) {
		t.Fatal("archive exists")
	}
}

func TestArchiveRewritesLinksUnderAConfiguredSpecRoot(t *testing.T) {
	req, _ := archiveLinksFixture(t, false, "\n[target](../target.md)\n")
	writeFile(t, filepath.Join(req.SpecsRoot, "target.md"), "target")
	result, err := legacyLinkArchive(t, req)
	if err != nil {
		t.Fatal(err)
	}
	if result.RewrittenLinks != 1 || !strings.Contains(archiveTestReadFile(t, filepath.Join(result.ArchivedDir, "_prd.md")), "[target](../../target.md)") {
		t.Fatal(result)
	}
}

func TestArchiveRewritesLinksInASupersededSpec(t *testing.T) {
	req, dir := archiveLinksFixture(t, true, "\n[target](../../adr/decision.md)\n")
	if err := os.Remove(filepath.Join(dir, "_tasks.md")); err != nil {
		t.Fatal(err)
	}
	writeFile(t, filepath.Join(dir, SupersessionFilename), "---\nsuperseded_by: successor\ndate: 2026-10-04\nreason: Delivered elsewhere\n---\n\nDelivered elsewhere.\n")
	prd := "---\nspec: demo\nstatus: active\nsuperseded_by: successor\nsuperseded: 2026-10-04\nsupersession_reason: Delivered elsewhere\n---\n\n[target](../../adr/decision.md)\n"
	writeFile(t, filepath.Join(dir, "_prd.md"), prd)
	writeFile(t, filepath.Join(req.SpecsRoot, "..", "adr", "decision.md"), "target")
	result, err := legacyLinkArchive(t, req)
	if err != nil {
		t.Fatal(err)
	}
	if got := archiveTestReadFile(t, filepath.Join(result.ArchivedDir, "_prd.md")); got != strings.Replace(prd, "../../adr/", "../../../adr/", 1) {
		t.Fatal(got)
	}
}

func TestArchiveKeepsALinkWhoseTargetWasAlreadyArchived(t *testing.T) {
	req, dir := archiveLinksFixture(t, true, "")
	body := "[archived](../../adr/decision.md)\n"
	writeFile(t, filepath.Join(dir, "notes.md"), body)
	writeFile(t, filepath.Join(req.SpecsRoot, "..", "history", "adr", "decision.md"), "target")
	result, err := legacyLinkArchive(t, req)
	if err != nil {
		t.Fatal(err)
	}
	if result.RewrittenLinks != 0 || archiveTestReadFile(t, filepath.Join(result.ArchivedDir, "notes.md")) != body {
		t.Fatal(result)
	}
}

func TestArchiveLinksMatch(t *testing.T) {
	activeDir := "/repo/docs/specs/demo"
	archivedDir := "/repo/docs/history/specs/demo"
	active := []byte("[ADR](../../adr/a.md?raw=1#choice)\n")
	tests := []struct {
		name, archived string
		want           bool
	}{
		{"rewrite", "[ADR](../../../adr/a.md?raw=1#choice)\n", true},
		{"unchanged", string(active), true},
		{"other target", "[ADR](../../../adr/b.md?raw=1#choice)\n", false},
		{"fragment", "[ADR](../../../adr/a.md?raw=1#other)\n", false},
		{"query", "[ADR](../../../adr/a.md?raw=2#choice)\n", false},
		{"body", "[Decision](../../../adr/a.md?raw=1#choice)\n", false},
		{"extra byte", "[ADR](../../../adr/a.md?raw=1#choice)\n!", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ArchiveLinksMatch(active, []byte(tt.archived), activeDir, archivedDir, activeDir); got != tt.want {
				t.Fatalf("match=%v, want %v", got, tt.want)
			}
		})
	}
	if ArchiveLinksMatch([]byte("[local](notes.md)"), []byte("[local](../../../specs/demo/notes.md)"), activeDir, archivedDir, activeDir) {
		t.Fatal("accepted internal rewrite")
	}
}

func TestArchiveRestoresRewrittenBytesWhenRenameFails(t *testing.T) {
	req, dir := archiveLinksFixture(t, false, "\n[target](../target.md)\n")
	writeFile(t, filepath.Join(req.SpecsRoot, "target.md"), "target")
	notes := "[target](../target.md)\n"
	writeFile(t, filepath.Join(dir, "notes.md"), notes)
	before := archiveTestReadFile(t, filepath.Join(dir, "_prd.md"))
	// Moving the configured archive directory into itself fails at rename,
	// after the link writes, without permissions or timing assumptions.
	req.Slug = "_archived"
	moved := filepath.Join(req.SpecsRoot, req.Slug)
	if err := os.Rename(dir, moved); err != nil {
		t.Fatal(err)
	}
	_, err := legacyLinkArchive(t, req)
	if err == nil || !strings.Contains(err.Error(), "move Spec") {
		t.Fatalf("expected rename failure: %v", err)
	}
	if got := archiveTestReadFile(t, filepath.Join(moved, "notes.md")); got != notes {
		t.Fatal("notes not restored")
	}
	if got := archiveTestReadFile(t, filepath.Join(moved, "_prd.md")); got != before {
		t.Fatal("PRD not restored")
	}
}

func archiveLinksTree(t *testing.T, dir string) map[string]string {
	t.Helper()
	files := make(map[string]string)
	err := filepath.WalkDir(dir, func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.Type()&os.ModeSymlink != 0 {
			target, err := os.Readlink(path)
			if err != nil {
				return err
			}
			files[path] = "symlink:" + target
		}
		if entry.Type().IsRegular() {
			files[path] = archiveTestReadFile(t, path)
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	return files
}

func TestArchiveLinksMarkdownForms(t *testing.T) {
	req, dir := archiveLinksFixture(t, true, "")
	target := filepath.Join(req.SpecsRoot, "..", "adr", "a(b).md")
	writeFile(t, target, "target")
	before := "[![badge](../../adr/a(b).md)](../../adr/a(b).md)\n[malformed](../../missing\n[nested [label]](../../adr/a(b).md 'title')\n![image](<../../adr/a(b).md>)\n[reference]: <../../adr/a(b).md> (title)\n`` [code](../../missing) ` ``\n\\[escaped](../../missing)\n"
	writeFile(t, filepath.Join(dir, "notes.md"), before)
	result, err := legacyLinkArchive(t, req)
	if err != nil {
		t.Fatal(err)
	}
	if result.RewrittenLinks != 5 {
		t.Fatalf("rewrote %d, want 5", result.RewrittenLinks)
	}
	after := archiveTestReadFile(t, filepath.Join(result.ArchivedDir, "notes.md"))
	if !ArchiveLinksMatch([]byte(before), []byte(after), dir, result.ArchivedDir, dir) {
		t.Fatal("Archive output does not match")
	}
}

// Legacy folders retain ADR-0230's link pass; new cuts never call it.
func legacyLinkArchive(t *testing.T, req ArchiveRequest) (ArchiveResult, error) {
	t.Helper()
	source := filepath.Join(req.SpecsRoot, req.Slug)
	destination := filepath.Join(ArchiveSpecRoot(req.SpecsRoot, req.BuiltInRoot), req.Slug)
	rewrites, count, err := prepareArchiveLinks(source, destination, req.Slug)
	if err != nil {
		return ArchiveResult{}, err
	}
	if err := os.MkdirAll(filepath.Dir(destination), 0755); err != nil {
		return ArchiveResult{}, err
	}
	for _, rewrite := range rewrites {
		if err := os.WriteFile(rewrite.path, rewriteArchiveLinks(rewrite.original, rewrite.destinations), rewrite.mode); err != nil {
			return ArchiveResult{}, err
		}
	}
	if err := os.Rename(source, destination); err != nil {
		return ArchiveResult{}, errors.Join(fmt.Errorf("move Spec: %w", err), restoreArchiveLinks(rewrites))
	}
	return ArchiveResult{SourceDir: source, ArchivedDir: destination, RewrittenLinks: count}, nil
}
