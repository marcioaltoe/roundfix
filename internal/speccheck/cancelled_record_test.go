package speccheck_test

// Suite: cancelled Archive Record and returned Backlog Entry consistency.
// Boundary IN: public Spec check in a disposable repository and home.
// Boundary OUT: cancellation writes and delivery eligibility.

import (
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"roundfix/internal/gittest"
	"roundfix/internal/spec"
	"roundfix/internal/speccheck"
)

func TestSpecCheckReadsACancelledRecordAndItsReturnedEntry(t *testing.T) {
	// Sequential: isolates process-wide HOME and XDG_CONFIG_HOME for Spec check and Git.
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("XDG_CONFIG_HOME", home)
	root := writeBacklogCarrier(t)
	gittest.InitRepo(t, root, "-b", "main")
	const slug = "0002-cancelled"
	const entry = "docs/backlog/2026-10-12-returned.md"
	recordPath := "docs/history/specs/" + slug + ".md"
	content, err := spec.RenderArchiveRecord(spec.ArchiveRecord{
		Spec: slug, Title: "Cancelled", Created: "2026-10-10", Archived: "2026-10-12",
		Source: "docs/specs/" + slug, SourceRevision: strings.Repeat("a", 40), Disposition: spec.ArchiveCancelled,
		CancelReason: "no consumer", CompletedTasks: []string{"task_01"}, PendingTasks: []string{"task_02"},
		ReturnedSources: []string{entry}, Outcome: "Original intent.",
	})
	if err != nil {
		t.Fatal(err)
	}
	writeFindingsArtifact(t, root, recordPath, string(content))
	writeFindingsArtifact(t, root, entry, "---\ntype: feat\nstatus: open\ncreated: 2026-10-12\nspec: null\n---\n\n# Returned intent\n")
	// A retired Finding still resolves its absorbed_by license through this record.
	writeFindingsArtifact(t, root, "docs/history/findings/2026-10-12-absorbed.md", "---\nstatus: done\nabsorbed_by: "+slug+"\n---\n\n# Retired evidence\n")
	beforeEntry, err := os.ReadFile(filepath.Join(root, entry))
	if err != nil {
		t.Fatal(err)
	}
	result := checkBacklogCarrier(t, root)
	for _, finding := range result.Findings {
		if finding.Code == speccheck.CodeBacklogUnmoved || finding.Code == speccheck.CodeArchiveLicense {
			t.Fatalf("cancelled record or returned entry has a finding: %+v", finding)
		}
		for _, where := range finding.Where {
			if where.Path == entry || where.Path == recordPath {
				t.Fatalf("cancelled record or returned entry has a finding: %+v", finding)
			}
		}
	}
	for name, before := range map[string][]byte{recordPath: content, entry: beforeEntry} {
		after, err := os.ReadFile(filepath.Join(root, name))
		if err != nil || !reflect.DeepEqual(before, after) {
			t.Fatalf("Spec check changed %s: err=%v", name, err)
		}
	}
}
