package speccheck_test

// Suite: archived Task Context through the public Spec check.
// Boundary IN: temporary active Specs and main-only clones with a bare forge.
// Boundary OUT: network and archive writes.

import (
	"path/filepath"
	"strings"
	"testing"

	"roundfix/internal/gittest"
	"roundfix/internal/spec"
	"roundfix/internal/speccheck"
)

func TestArchivedTaskContextRecoversASquashDroppedRevision(t *testing.T) {
	t.Parallel()
	for _, specRoot := range []string{"docs/specs", "planning"} {
		t.Run(specRoot, func(t *testing.T) {
			for _, number := range []string{"7", "", "unreachable"} {
				t.Run(map[string]string{"7": "record", "": "adding subject", "unreachable": "unreachable forge"}[number], func(t *testing.T) {
					const slug = "0002-other"
					source := specRoot + "/" + slug
					archiveRoot := "docs/history/specs"
					if specRoot != "docs/specs" {
						archiveRoot = specRoot + "/_archived"
					}
					root := pathPinRepository(t, specRoot, "pending", source+"/_techspec.md")
					gittest.Run(t, root, "checkout", "-b", "main")
					gittest.Run(t, root, "add", ".")
					gittest.Run(t, root, "commit", "-m", "docs: active Spec")
					gittest.Run(t, root, "checkout", "-b", "feat/source")
					writeFindingsArtifact(t, root, source+"/_techspec.md", "# Historical context\n")
					gittest.Run(t, root, "add", ".")
					gittest.Run(t, root, "commit", "-m", "docs: source")
					revision := strings.TrimSpace(gittest.Run(t, root, "rev-parse", "HEAD"))
					gittest.Run(t, root, "rm", "-r", source)
					r := spec.ArchiveRecord{Spec: slug, Title: "Archived context", Source: source, SourceRevision: revision, Disposition: spec.ArchivePass, PullRequest: number}
					if number == "unreachable" {
						r.PullRequest = "7"
					}
					content, err := spec.RenderArchiveRecord(r)
					if err != nil {
						t.Fatal(err)
					}
					writeFindingsArtifact(t, root, archiveRoot+"/"+slug+".md", string(content))
					gittest.Run(t, root, "add", ".")
					gittest.Run(t, root, "commit", "-m", "docs: archive")
					gittest.Run(t, root, "checkout", "main")
					gittest.Run(t, root, "merge", "--squash", "feat/source")
					gittest.Run(t, root, "commit", "-m", "feat: delivery (#7)")
					clone, _ := gittest.PullRequestClone(t, root, "feat/source", revision)
					if number == "unreachable" {
						gittest.Run(t, clone, "remote", "set-url", "origin", filepath.Join(t.TempDir(), "absent.git"))
					}
					result, err := speccheck.Check(filepath.Join(clone, specRoot), clone, pathPinSlug)
					if err != nil {
						t.Fatal(err)
					}
					if number == "unreachable" {
						if findings := findingsWithCode(result, speccheck.CodeReferenceUnresolved); len(findings) != 1 {
							t.Fatalf("unavailable context findings=%+v", findings)
						}
						return
					}
					if findings := findingsWithCode(result, speccheck.CodeReferenceUnresolved); len(findings) != 0 {
						t.Fatalf("context unresolved: %+v", findings)
					}
					gittest.Run(t, clone, "cat-file", "-e", revision+"^{commit}")
				})
			}
		})
	}
}
