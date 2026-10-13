package spec

// Suite: failed final Spec folder removal.
// Invariant: typed removal failures keep the filled result and every write.
// Boundary: real filesystem permissions in committed temporary repositories.

import (
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func TestArchiveReturnsAFolderRemovalError(t *testing.T) {
	t.Parallel()
	req, dir := recordFixture(t)
	writeFile(t, filepath.Join(dir, "references", "note.md"), "Promoted note.\n")
	req.Promote = []string{"references/note.md"}
	writeFile(t, filepath.Join(req.RepositoryRoot, "docs/backlog/followup.md"), "[Spec](../specs/demo/_prd.md)\n")
	req.SourceRevision = archiveFixtureRevision(t, req.SpecsRoot)
	before := archiveLinksTree(t, dir)
	makeSpecReferencesReadOnly(t, dir)

	result, err := Archive(req)
	assertFolderRemovalError(t, err, dir)
	var size int64
	for _, content := range before {
		size += int64(len(content))
	}
	if result.SourceDir != dir || result.RecordPath != ArchiveRecordPath(ArchiveSpecRoot(req.SpecsRoot, true), req.Slug) || result.ArchivedDir != strings.TrimSuffix(result.RecordPath, ".md") || result.ArchivedOn == "" || result.RemovedFiles != len(before) || result.RemovedBytes != size || !reflect.DeepEqual(result.Promoted, []string{"docs/references/note.md"}) || !reflect.DeepEqual(result.Relinked, []string{"docs/backlog/followup.md"}) {
		t.Fatalf("unfilled removal result: %+v", result)
	}
	data, err := os.ReadFile(result.RecordPath)
	if err != nil {
		t.Fatal(err)
	}
	record, err := ParseArchiveRecord(data)
	if err != nil {
		t.Fatal(err)
	}
	if record.Disposition != ArchivePass || record.SourceRevision != req.SourceRevision || !reflect.DeepEqual(record.Promoted, result.Promoted) || !reflect.DeepEqual(record.Relinked, result.Relinked) {
		t.Fatalf("record=%+v result=%+v", record, result)
	}
	assertRemovalWrite(t, req.RepositoryRoot, "docs/references/note.md", "Promoted note.\n")
	assertRemovalWrite(t, req.RepositoryRoot, "docs/backlog/followup.md", "[Spec](../history/specs/demo.md)\n")
}

func TestCancelReturnsAFolderRemovalError(t *testing.T) {
	t.Parallel()
	req, dir := cancelFixture(t, 8)
	cancelAdopt(t, dir, "finding", "finding.md", cancelFinding)
	writeFile(t, filepath.Join(dir, "references", "note.md"), "Promoted note.\n")
	writeFile(t, filepath.Join(req.RepositoryRoot, "docs/backlog/followup.md"), "[Note](../specs/demo/references/note.md)\n[Spec](../specs/demo/_prd.md)\n")
	req.SourceRevision = archiveFixtureRevision(t, req.SpecsRoot)
	before := archiveLinksTree(t, dir)
	makeSpecReferencesReadOnly(t, dir)

	result, err := Cancel(req)
	assertFolderRemovalError(t, err, dir)
	var size int64
	for _, content := range before {
		size += int64(len(content))
	}
	if result.RecordPath != ArchiveRecordPath(ArchiveSpecRoot(req.SpecsRoot, true), req.Slug) || result.RemovedFiles != len(before) || result.RemovedBytes != size || len(result.RemovedPaths) != len(before) || !reflect.DeepEqual(result.Completed, []string{"task_01", "task_02", "task_03", "task_04", "task_05"}) || !reflect.DeepEqual(result.Pending, []string{"task_06", "task_07", "task_08"}) || len(result.Returned) != 2 {
		t.Fatalf("unfilled removal result: %+v", result)
	}
	for _, path := range result.RemovedPaths {
		if _, ok := before[filepath.Join(req.RepositoryRoot, filepath.FromSlash(path))]; !ok {
			t.Fatalf("unexpected removed path %q", path)
		}
	}
	record := cancelReadRecord(t, result)
	if record.Disposition != ArchiveCancelled || record.SourceRevision != req.SourceRevision || !reflect.DeepEqual(record.CompletedTasks, result.Completed) || !reflect.DeepEqual(record.PendingTasks, result.Pending) || !reflect.DeepEqual(record.ReturnedSources, []string{"docs/backlog/entry.md", "docs/findings/finding.md"}) {
		t.Fatalf("record=%+v result=%+v", record, result)
	}
	for _, source := range result.Returned {
		data, err := os.ReadFile(filepath.Join(req.RepositoryRoot, filepath.FromSlash(source.To)))
		if err != nil || !strings.Contains(string(data), "status: "+source.Status) {
			t.Fatalf("returned source=%+v data=%q err=%v", source, data, err)
		}
	}
	assertRemovalWrite(t, req.RepositoryRoot, "docs/references/note.md", "Promoted note.\n")
	assertRemovalWrite(t, req.RepositoryRoot, "docs/backlog/followup.md", "[Note](../references/note.md)\n[Spec](../history/specs/demo.md)\n")
}

func makeSpecReferencesReadOnly(t *testing.T, dir string) {
	t.Helper()
	path := filepath.Join(dir, "references")
	t.Cleanup(func() {
		if err := os.Chmod(path, 0o755); err != nil && !errors.Is(err, os.ErrNotExist) {
			t.Error(err)
		}
	})
	if err := os.Chmod(path, 0o555); err != nil {
		t.Fatal(err)
	}
}

func assertFolderRemovalError(t *testing.T, err error, folder string) {
	t.Helper()
	var removal *FolderRemovalError
	var cause *os.PathError
	if !errors.As(err, &removal) || removal.Folder != folder || !filepath.IsAbs(removal.Folder) || !errors.Is(err, os.ErrPermission) || !errors.As(err, &cause) || !errors.Is(removal.Err, cause) {
		t.Fatalf("expected typed removal error for %q wrapping permission failure: %v", folder, err)
	}
	if !strings.Contains(err.Error(), folder) || !strings.Contains(err.Error(), cause.Error()) {
		t.Fatalf("removal error lacks folder or cause: %v", err)
	}
}

func assertRemovalWrite(t *testing.T, root, path, want string) {
	t.Helper()
	data, err := os.ReadFile(filepath.Join(root, filepath.FromSlash(path)))
	if err != nil || string(data) != want {
		t.Fatalf("write %s = %q, %v; want %q", path, data, err, want)
	}
}
