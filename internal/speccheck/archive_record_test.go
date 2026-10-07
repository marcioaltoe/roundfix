package speccheck_test

import (
	"os"
	"path/filepath"
	"roundfix/internal/gittest"
	"roundfix/internal/spec"
	"roundfix/internal/speccheck"
	"strings"
	"testing"
)

func writeCheckArchiveRecord(t *testing.T, root, archiveRoot, source, slug, revision string) {
	t.Helper()
	content, err := spec.RenderArchiveRecord(spec.ArchiveRecord{
		Spec: slug, Title: "Archived fixture", Created: "2026-10-06", Archived: "2026-10-06",
		Disposition: spec.ArchivePass, Source: source, SourceRevision: revision,
		QATask: "task_qa", QAReport: "qa-report-2026-10-06.md", QAVerdict: "pass", Outcome: "Preserved fixture.",
	})
	if err != nil {
		t.Fatal(err)
	}
	writeFindingsArtifact(t, root, archiveRoot+"/"+slug+".md", string(content))
}

func TestArchivedSlugsIncludeArchiveRecords(t *testing.T) {
	root := writeBacklogCarrier(t)
	const slug = "0002-record-only"
	recordPath := "docs/history/specs/" + slug + ".md"
	writeCheckArchiveRecord(t, root, "docs/history/specs", "docs/specs/"+slug, slug, strings.Repeat("a", 40))
	writeFindingsArtifact(t, root, "docs/history/findings/2026-10-06-absorbed.md", "---\nstatus: done\nabsorbed_by: "+slug+"\n---\n# Absorbed\n")
	writeFindingsArtifact(t, root, "docs/backlog/2026-10-06-promoted.md", "---\nstatus: promoted\nspec: "+slug+"\n---\n# Promoted\n")
	result := checkBacklogCarrier(t, root)
	if findings := findingsWithCode(result, speccheck.CodeArchiveLicense); len(findings) != 0 {
		t.Fatalf("record-only finding license unresolved: %+v", findings)
	}
	findings := findingsWithCode(result, speccheck.CodeBacklogUnmoved)
	if len(findings) != 1 || strings.Contains(findings[0].Summary, "unresolvable") || !strings.Contains(findings[0].Fix, recordPath) || strings.Contains(findings[0].Fix, slug+"/references/") {
		t.Fatalf("promotion findings: %+v", findings)
	}
}

func TestTaskContextResolvesThroughTheArchiveRecord(t *testing.T) {
	for _, specRoot := range []string{"docs/specs", "planning"} {
		t.Run(specRoot, func(t *testing.T) {
			for _, spelling := range []string{"active", "archived"} {
				t.Run(spelling, func(t *testing.T) {
					const slug = "0002-other"
					source := specRoot + "/" + slug
					archiveRoot := "docs/history/specs"
					if specRoot != "docs/specs" {
						archiveRoot = specRoot + "/_archived"
					}
					ref := source + "/_techspec.md"
					if spelling == "archived" {
						ref = archiveRoot + "/" + slug + "/_techspec.md"
					}
					root := pathPinRepository(t, specRoot, "pending", ref)
					gittest.InitRepo(t, root, "-b", "main")
					writeFindingsArtifact(t, root, source+"/_techspec.md", "# Historical context\n")
					gittest.Run(t, root, "add", ".")
					gittest.Run(t, root, "commit", "-m", "docs: historical context")
					revision := strings.TrimSpace(gittest.Run(t, root, "rev-parse", "HEAD"))
					if err := os.RemoveAll(filepath.Join(root, source)); err != nil {
						t.Fatal(err)
					}
					writeCheckArchiveRecord(t, root, archiveRoot, source, slug, revision)
					check := func() []speccheck.Finding {
						result, err := speccheck.Check(filepath.Join(root, specRoot), root, pathPinSlug)
						if err != nil {
							t.Fatal(err)
						}
						return findingsWithCode(result, speccheck.CodeReferenceUnresolved)
					}
					if findings := check(); len(findings) != 0 {
						t.Fatalf("historical context unresolved: %+v", findings)
					}
					writeCheckArchiveRecord(t, root, archiveRoot, source, slug, strings.Repeat("a", 40))
					if findings := check(); len(findings) != 1 {
						t.Fatalf("unavailable source resolved: %+v", findings)
					}
					writeCheckArchiveRecord(t, root, archiveRoot, source, slug, revision)
					// A known commit does not prove a file absent from that commit.
					if err := os.Remove(filepath.Join(root, archiveRoot, slug+".md")); err != nil {
						t.Fatal(err)
					}
					if findings := check(); len(findings) != 1 {
						t.Fatalf("missing record resolved: %+v", findings)
					}
				})
			}
		})
	}
}
