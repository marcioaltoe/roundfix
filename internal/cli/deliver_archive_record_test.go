package cli

// Suite: Delivery Queue exact retirement proof.
// Invariant: only the record, complete deletions and matching promotions may change.
// Boundary IN: real Git commits and the real Archive Command.
// Boundary OUT: review and publication remain the Delivery Engine's concern.

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	roundconfig "roundfix/internal/config"
	"roundfix/internal/gittest"
	"roundfix/internal/preflight"
	"roundfix/internal/spec"
)

func retirementWorkflow(repo, home string) *commandDeliveryWorkflow {
	return &commandDeliveryWorkflow{loaded: roundconfig.Loaded{GitRoot: repo, HomeDir: home, Config: roundconfig.Config{Specs: roundconfig.Specs{Root: "docs/specs"}}}, git: preflight.ExecGitRunner{}}
}
func retirementFixture(t *testing.T, promote bool) (*commandDeliveryWorkflow, string, string, string, string) {
	t.Helper()
	home, repo := newImplementWorkspace(t, []implementSeed{{id: "task_01", status: string(spec.StatusCompleted)}})
	writeArchiveQAReport(t, repo, spec.VerdictPass)
	dir := filepath.Join(repo, "docs/specs", implementTestSlug)
	mustMkdir(t, filepath.Join(dir, "references"))
	mustWrite(t, filepath.Join(dir, "references", "knowledge.md"), "Reusable knowledge\n")
	commitArchiveFixture(t)
	parent := strings.TrimSpace(gittest.Run(t, repo, "rev-parse", "HEAD"))
	workflow := retirementWorkflow(repo, home)
	source, destination, err := workflow.archivePaths(repo, implementTestSlug)
	if err != nil {
		t.Fatal(err)
	}
	req := spec.ArchiveRequest{SpecsRoot: filepath.Join(repo, "docs/specs"), BuiltInRoot: true, Slug: implementTestSlug, SourceRevision: parent, RepositoryRoot: repo}
	if promote {
		req.Promote = []string{"references/knowledge.md"}
	}
	if _, err := spec.Archive(req); err != nil {
		t.Fatal(err)
	}
	return workflow, repo, parent, source, destination
}
func TestDeliveryAcceptsAnExactRetirement(t *testing.T) {
	for _, promote := range []bool{false, true} {
		t.Run(map[bool]string{false: "record", true: "promotion"}[promote], func(t *testing.T) {
			w, repo, parent, source, destination := retirementFixture(t, promote)
			exact, err := w.archiveDiffIsExact(t.Context(), repo, source, destination)
			if err != nil || !exact {
				t.Fatalf("working retirement=%v err=%v", exact, err)
			}
			gittest.Run(t, repo, "add", "-A")
			gittest.Run(t, repo, "commit", "-m", "docs: archive record")
			head := strings.TrimSpace(gittest.Run(t, repo, "rev-parse", "HEAD"))
			exact, err = w.archiveCommitIsExact(t.Context(), repo, parent, head, source, destination)
			if err != nil || !exact {
				t.Fatalf("committed retirement=%v err=%v", exact, err)
			}
			reconciled, err := w.reconcileArchiveCommit(t.Context(), repo, implementTestSlug, parent, head)
			if err != nil || !reconciled.ExactSpecMove {
				t.Fatalf("reconcile=%+v err=%v", reconciled, err)
			}
		})
	}
}
func TestDeliveryRefusesAnInexactRetirement(t *testing.T) {
	for _, kind := range []string{"keeps file", "extra path", "other revision", "different promotion", "other title", "other created"} {
		t.Run(kind, func(t *testing.T) {
			w, repo, parent, source, destination := retirementFixture(t, true)
			recordPath := filepath.Join(repo, destination+".md")
			switch kind {
			case "keeps file":
				mustMkdir(t, filepath.Join(repo, source))
				mustWrite(t, filepath.Join(repo, source, "_prd.md"), gittest.Run(t, repo, "show", parent+":"+source+"/_prd.md"))
			case "extra path":
				mustWrite(t, filepath.Join(repo, "unrelated.md"), "Extra\n")
			case "different promotion":
				mustWrite(t, filepath.Join(repo, "docs/references/knowledge.md"), "Different\n")
			default:
				r := archiveRecordAt(t, recordPath)
				switch kind {
				case "other revision":
					r.SourceRevision = strings.Repeat("b", 40)
				case "other title":
					r.Title = "Other"
				case "other created":
					r.Created = "2026-10-01"
				}
				content, err := spec.RenderArchiveRecord(r)
				if err != nil {
					t.Fatal(err)
				}
				mustWrite(t, recordPath, string(content))
			}
			exact, err := w.archiveDiffIsExact(t.Context(), repo, source, destination)
			if err != nil || exact {
				t.Fatalf("working proof=%v err=%v", exact, err)
			}
			gittest.Run(t, repo, "add", "-A")
			gittest.Run(t, repo, "commit", "-m", "docs: inexact archive")
			head := strings.TrimSpace(gittest.Run(t, repo, "rev-parse", "HEAD"))
			exact, err = w.archiveCommitIsExact(t.Context(), repo, parent, head, source, destination)
			if err != nil || exact {
				t.Fatalf("commit proof=%v err=%v", exact, err)
			}
		})
	}
}
func TestDeliveryKeepsTheLegacyExactMove(t *testing.T) {
	repo, parent, head := commitLinkRewritingArchive(t, []string{"_prd.md", "task_01.md", "qa/links.md"}, nil)
	w := retirementWorkflow(repo, "")
	source, destination, err := w.archivePaths(repo, implementTestSlug)
	if err != nil {
		t.Fatal(err)
	}
	exact, err := w.archiveCommitIsExact(t.Context(), repo, parent, head, source, destination)
	if err != nil || !exact {
		t.Fatalf("legacy proof=%v err=%v", exact, err)
	}
}
func TestDeliveryArchiveStageCommitsTheArchiveRecord(t *testing.T) {
	home, repo := newImplementWorkspace(t, []implementSeed{{id: "task_01", status: string(spec.StatusCompleted)}})
	writeArchiveQAReport(t, repo, spec.VerdictPass)
	commitArchiveFixture(t)
	parent := strings.TrimSpace(gittest.Run(t, repo, "rev-parse", "HEAD"))
	gittest.PersistIdentity(t, repo)
	t.Setenv(cliTestHelperEnv, "1")
	w := retirementWorkflow(repo, home)
	var log bytes.Buffer
	w.log = &log
	result, err := w.Archive(t.Context(), repo, implementTestSlug, parent)
	if err != nil || !result.ExactSpecMove || result.Parent != parent || result.Head == parent {
		t.Fatalf("archive stage=%+v err=%v", result, err)
	}
	recordPath := archiveTestRepositoryPath(repo, spec.ArchiveKindSpec, implementTestSlug) + ".md"
	r := archiveRecordAt(t, recordPath)
	if r.SourceRevision != parent {
		t.Fatalf("source revision=%s", r.SourceRevision)
	}
	assertPathMissing(t, filepath.Join(repo, r.Source))
	if status := strings.TrimSpace(gittest.Run(t, repo, "status", "--porcelain")); status != "" {
		t.Fatalf("dirty archive: %s", status)
	}
	again, err := w.Archive(t.Context(), repo, implementTestSlug, result.Head)
	if err != nil || !again.AlreadyArchived {
		t.Fatalf("already archived=%+v err=%v", again, err)
	}
	if _, err := os.Stat(recordPath); err != nil {
		t.Fatal(err)
	}
}
