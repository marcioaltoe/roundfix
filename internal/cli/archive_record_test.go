package cli

// Suite: Archive Command retirement and Git provenance.
// Invariant: the confirmation accounts for bytes recoverable at source_revision.
// Boundary IN: real command in committed temporary repositories.
// Boundary OUT: record construction lives in internal/spec.

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"path/filepath"
	"strings"
	"testing"

	roundconfig "roundfix/internal/config"
	"roundfix/internal/gittest"
	"roundfix/internal/spec"
)

func initGitRepository(t *testing.T, root string) {
	t.Helper()
	gittest.InitRepo(t, root)
	gittest.Run(t, root, "config", "user.name", "Archive Test")
	gittest.Run(t, root, "config", "user.email", "archive@example.com")
	gittest.Run(t, root, "add", "-A")
	gittest.Run(t, root, "commit", "-m", "docs: archive fixture")
}
func commitArchiveFixture(t *testing.T) {
	t.Helper()
	if !commandEnvironmentOverridesForTest(t).workDirSet {
		t.Fatal("archive fixture requires an explicit disposable work directory")
	}
	loaded, err := loadCommandConfig(commandEnvironmentForTest(t), io.Discard)
	if err != nil {
		t.Fatal(err)
	}
	root, err := roundconfig.ResolveSpecsRoot(loaded, loaded.GitRoot)
	if err != nil {
		t.Fatal(err)
	}
	repo := loaded.GitRoot
	if root.External {
		repo = root.Path
	}
	if strings.TrimSpace(gittest.Run(t, repo, "status", "--porcelain")) == "" {
		return
	}
	gittest.Run(t, repo, "add", "-A")
	gittest.Run(t, repo, "commit", "-m", "docs: committed archive fixture")
}
func archiveRecordAt(t *testing.T, path string) spec.ArchiveRecord {
	t.Helper()
	record, err := spec.ParseArchiveRecord([]byte(mustRead(t, path)))
	if err != nil {
		t.Fatal(err)
	}
	return record
}
func archivedSourceFiles(t *testing.T, repo, recordPath string) map[string][]byte {
	t.Helper()
	r := archiveRecordAt(t, recordPath)
	if !strings.HasPrefix(recordPath, repo+string(filepath.Separator)) {
		repo = filepath.Dir(filepath.Dir(recordPath))
	}
	listing := gittest.Run(t, repo, "ls-tree", "-r", "--name-only", r.SourceRevision, "--", r.Source)
	files := map[string][]byte{}
	for _, path := range strings.Split(strings.TrimSpace(listing), "\n") {
		if path != "" {
			files[strings.TrimPrefix(path, r.Source+"/")] = []byte(gittest.Run(t, repo, "show", r.SourceRevision+":"+path))
		}
	}
	return files
}
func archiveConfirmation(t *testing.T, repo, path string, override bool) string {
	t.Helper()
	r := archiveRecordAt(t, path)
	files := archivedSourceFiles(t, repo, path)
	size := 0
	for _, text := range files {
		size += len(text)
	}
	disposition := ""
	if override {
		disposition = " with QA override"
	}
	rel, err := filepathRelSlash(repo, path)
	if err != nil {
		t.Fatal(err)
	}
	return fmt.Sprintf("archived %s%s -> %s; removed %d file(s) (%d bytes) kept in Git at %.12s\n", r.Spec, disposition, rel, len(files), size, r.SourceRevision)
}
func TestArchiveCommandLeavesTheArchiveRecord(t *testing.T) {
	t.Parallel()
	for _, override := range []bool{false, true} {
		t.Run(fmt.Sprint(override), func(t *testing.T) {
			_, repo := newImplementWorkspace(t, []implementSeed{{id: "task_01", status: string(spec.StatusCompleted)}})
			verdict := spec.VerdictPass
			if override {
				verdict = spec.VerdictFail
			}
			writeArchiveQAReport(t, repo, verdict)
			commitArchiveFixture(t)
			args := []string{"archive", implementTestSlug}
			if override {
				args = append(args, "--qa-override", "--approval", "maintainer", "--reason", strings.Repeat("r", 300))
			}
			var out, errout bytes.Buffer
			code := runCLIContext(t, t.Context(), args, &out, &errout)
			path := archiveTestRepositoryPath(repo, spec.ArchiveKindSpec, implementTestSlug) + ".md"
			if code != exitOK {
				t.Fatalf("exit=%d stderr=%s", code, &errout)
			}
			if out.String() != archiveConfirmation(t, repo, path, override) || errout.Len() != 0 {
				t.Fatalf("stdout=%q stderr=%q", &out, &errout)
			}
			assertPathMissing(t, filepath.Join(repo, "docs/specs", implementTestSlug))
			if len(mustRead(t, path)) > spec.ArchiveRecordTargetBytes {
				t.Fatal("oversized record")
			}
		})
	}
}
func TestArchiveCommandRefusesUncommittedSpecChanges(t *testing.T) {
	t.Parallel()
	for _, kind := range []string{"modified", "staged", "untracked", "deleted"} {
		t.Run(kind, func(t *testing.T) {
			_, repo := newImplementWorkspace(t, []implementSeed{{id: "task_01", status: string(spec.StatusCompleted)}})
			writeArchiveQAReport(t, repo, spec.VerdictPass)
			mustWrite(t, filepath.Join(repo, "docs/specs", implementTestSlug, "evidence.txt"), "evidence")
			commitArchiveFixture(t)
			dir := filepath.Join(repo, "docs/specs", implementTestSlug)
			path := filepath.Join(dir, "_prd.md")
			switch kind {
			case "untracked":
				mustWrite(t, filepath.Join(dir, "new.md"), "new")
			case "deleted":
				gittest.Run(t, repo, "rm", filepath.Join("docs/specs", implementTestSlug, "evidence.txt"))
			default:
				mustWrite(t, path, mustRead(t, path)+"\nchanged\n")
				if kind == "staged" {
					gittest.Run(t, repo, "add", "-A")
				}
			}
			before := snapshotDirectoryFiles(t, dir)
			var out, errout bytes.Buffer
			code := runCLIContext(t, context.Background(), []string{"archive", implementTestSlug}, &out, &errout)
			reason := fmt.Sprintf("Spec %q has changes not committed at HEAD under docs/specs/%s; commit them before archive so the Archive Record's source_revision holds the Spec", implementTestSlug, implementTestSlug)
			if code != exitPreflight || out.Len() != 0 || !strings.Contains(errout.String(), reason) {
				t.Fatalf("exit=%d stdout=%s stderr=%s", code, &out, &errout)
			}
			after := snapshotDirectoryFiles(t, dir)
			if fmt.Sprint(before) != fmt.Sprint(after) {
				t.Fatal("refusal changed files")
			}
			assertPathMissing(t, archiveTestRepositoryPath(repo, spec.ArchiveKindSpec, implementTestSlug)+".md")
		})
	}
}
func TestArchiveCommandRemovedBytesStayInGit(t *testing.T) {
	t.Parallel()
	_, repo := newImplementWorkspace(t, []implementSeed{{id: "task_01", status: string(spec.StatusCompleted)}})
	writeArchiveQAReport(t, repo, spec.VerdictPass)
	dir := filepath.Join(repo, "docs/specs", implementTestSlug)
	mustMkdir(t, filepath.Join(dir, "qa/evidence"))
	mustWrite(t, filepath.Join(dir, "qa/evidence", "binary.dat"), "\x00\x01\xff\n")
	before := snapshotDirectoryFiles(t, dir)
	commitArchiveFixture(t)
	var out, errout bytes.Buffer
	if code := runCLIContext(t, t.Context(), []string{"archive", implementTestSlug}, &out, &errout); code != exitOK {
		t.Fatalf("exit=%d stderr=%s", code, &errout)
	}
	path := archiveTestRepositoryPath(repo, spec.ArchiveKindSpec, implementTestSlug) + ".md"
	after := archivedSourceFiles(t, repo, path)
	if len(before) != len(after) {
		t.Fatalf("files %d != %d", len(before), len(after))
	}
	for name, text := range before {
		if !bytes.Equal(after[name], text) {
			t.Fatalf("removed bytes differ for %s", name)
		}
	}
	if out.String() != archiveConfirmation(t, repo, path, false) {
		t.Fatalf("confirmation=%s", &out)
	}
}
func TestSupersedeAcceptsAnArchiveRecord(t *testing.T) {
	// Sequential: calls the already parallel TestSupersedeAcceptsASupersessionArchivedDeliverer on the same testing.T.
	TestSupersedeAcceptsASupersessionArchivedDeliverer(t)
}
