package specaudit

import (
	"os"
	"path/filepath"
	"reflect"
	"roundfix/internal/spec"
	"testing"
)

func TestAuditClaimsTheArchiveRecord(t *testing.T) {
	for _, delivered := range []bool{false, true} {
		name := "held by branch"
		if delivered {
			name = "delivered"
		}
		t.Run(name, func(t *testing.T) {
			f := newAuditFixture(t)
			active := "docs/specs/" + auditFixtureSlug
			writeAuditFixtureFile(t, filepath.Join(f.repoDir, active, "_prd.md"), "# PRD\n")
			writeAuditFixtureFile(t, filepath.Join(f.repoDir, active, "task_01.md"), "# Task\n")
			writeAuditFixtureFile(t, filepath.Join(f.repoDir, "docs/specs/0001-carrier/_prd.md"), "# Carrier\n")
			f.git("add", ".")
			f.git("commit", "-m", "docs: task artifacts", "-m", "Roundfix-Spec: "+auditFixtureSlug+"\nRoundfix-Task: task_01")
			f.git("switch", "-c", "docs/archive-fixture")
			revision := f.git("rev-parse", "HEAD")
			recordPath := "docs/history/specs/" + auditFixtureSlug + ".md"
			record, err := spec.RenderArchiveRecord(spec.ArchiveRecord{Spec: auditFixtureSlug, Title: "Fixture", Created: "2026-10-06", Archived: "2026-10-06", Disposition: spec.ArchivePass, Source: active, SourceRevision: revision, Outcome: "Archived."})
			if err != nil {
				t.Fatal(err)
			}
			writeAuditFixtureFile(t, filepath.Join(f.repoDir, recordPath), string(record))
			if err := os.RemoveAll(filepath.Join(f.repoDir, active)); err != nil {
				t.Fatal(err)
			}
			f.git("add", "-A")
			f.git("commit", "-m", "docs: archive fixture")
			f.git("switch", "main")
			if delivered {
				f.git("merge", "--ff-only", "docs/archive-fixture")
			}
			before := snapshotGitState(t, f.repoDir)
			result := f.audit()
			var want []Undelivered
			if !delivered {
				want = []Undelivered{{Artifact: recordPath, HeldBy: "docs/archive-fixture"}}
			}
			if len(result.Undelivered) != len(want) || (len(want) > 0 && !reflect.DeepEqual(result.Undelivered, want)) {
				t.Fatalf("undelivered = %+v, want %+v", result.Undelivered, want)
			}
			if !reflect.DeepEqual(before, snapshotGitState(t, f.repoDir)) {
				t.Fatal("audit mutated Git")
			}
		})
	}
}
