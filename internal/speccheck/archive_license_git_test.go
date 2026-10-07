package speccheck_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"roundfix/internal/gittest"
	"roundfix/internal/speccheck"
)

func TestArchiveLicenseResolvesThroughTheRecordOrGit(t *testing.T) {
	root := writeBacklogCarrier(t)
	gittest.InitRepo(t, root, "-b", "main")
	writeFindingsArtifact(t, root, "docs/history/specs/0003-legacy/_prd.md", "# Legacy Spec\n")
	gittest.Run(t, root, "add", ".")
	gittest.Run(t, root, "commit", "-m", "docs: legacy archive")
	revision := strings.TrimSpace(gittest.Run(t, root, "rev-parse", "HEAD"))
	gittest.Run(t, root, "rm", "-r", "docs/history/specs/0003-legacy")
	gittest.Run(t, root, "commit", "-m", "docs: remove legacy folder")
	writeCheckArchiveRecord(t, root, "docs/history/specs", "docs/specs/0002-record-only", "0002-record-only", revision)
	for _, slug := range []string{"0002-record-only", "0003-legacy", "0004-never-archived"} {
		writeFindingsArtifact(t, root, "docs/history/findings/2026-10-06-"+slug+".md", "---\nstatus: done\nabsorbed_by: "+slug+"\n---\n# Absorbed\n")
	}
	findings := findingsWithCode(checkBacklogCarrier(t, root), speccheck.CodeArchiveLicense)
	if len(findings) != 1 || !strings.Contains(findings[0].Summary, "0004-never-archived") {
		t.Fatalf("archive licenses: %+v, want only 0004-never-archived unresolved", findings)
	}
	assertFilesystemFallback := func(t *testing.T) {
		t.Helper()
		findings := findingsWithCode(checkBacklogCarrier(t, root), speccheck.CodeArchiveLicense)
		if len(findings) != 2 || !strings.Contains(findings[0].Summary, "0003-legacy") || !strings.Contains(findings[1].Summary, "0004-never-archived") {
			t.Fatalf("filesystem-only licenses: %+v, want legacy and never-archived unresolved", findings)
		}
	}
	t.Run("Git fails without HEAD", func(t *testing.T) {
		gittest.Run(t, root, "update-ref", "-d", "refs/heads/main")
		assertFilesystemFallback(t)
	})
	t.Run("outside Git", func(t *testing.T) {
		if err := os.RemoveAll(filepath.Join(root, ".git")); err != nil {
			t.Fatal(err)
		}
		assertFilesystemFallback(t)
	})
}
