package spec

// Suite: tracked citation readers and working-tree entries.
// Invariant: missing entries are skipped; every other read error is retained.
// Boundary IN: committed sources changed only on disk in isolated Git repositories.
// Boundary OUT: source paths, surviving citations and errors naming blocked paths.

import (
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"roundfix/internal/gittest"
)

const workingTreeCitationTarget = "docs/history/findings/x.md"

func trackedSourceReaderFixture(t *testing.T) string {
	t.Helper()
	root := t.TempDir()
	gittest.InitRepo(t, root, "-b", "main")
	for _, path := range []string{"README.md", "docs/gone.md", "pkg/gone.test.ts", "notes/a.md"} {
		writeStaleCitationSource(t, root, path, "`"+workingTreeCitationTarget+"`\n")
	}
	legacyCommit(t, root, "docs: tracked citation sources")
	return root
}

func trackedSourceReaderPaths(t *testing.T, root, reader string) ([]string, error) {
	t.Helper()
	if reader == "citation-sources" {
		return CitationSources(t.Context(), root, false)
	}
	var citations []HistoryCitation
	var err error
	if reader == "history-citations" {
		citations, err = HistoryCitations(t.Context(), root, []string{workingTreeCitationTarget})
	} else {
		citations, err = StaleHistoryCitations(t.Context(), root, "docs/history", map[string]bool{workingTreeCitationTarget: true})
	}
	var paths []string
	for _, citation := range citations {
		if citation.Target != workingTreeCitationTarget || citation.Line != 1 {
			t.Fatalf("unexpected citation: %+v", citation)
		}
		paths = append(paths, citation.Path)
	}
	return paths, err
}

func TestTrackedSourceReadersSkipAPathTheWorkingTreeLacks(t *testing.T) {
	t.Parallel()
	for _, reader := range []string{"citation-sources", "history-citations", "stale-history-citations"} {
		t.Run(reader, func(t *testing.T) {
			t.Parallel()
			if !citationSourcesTestHome(t) {
				return
			}
			root := trackedSourceReaderFixture(t)
			for _, path := range []string{"docs/gone.md", "pkg/gone.test.ts"} {
				if err := os.Remove(filepath.Join(root, path)); err != nil {
					t.Fatal(err)
				}
			}
			got, err := trackedSourceReaderPaths(t, root, reader)
			want := []string{"README.md", "notes/a.md"}
			if err != nil || !reflect.DeepEqual(got, want) {
				t.Fatalf("paths=%v err=%v want=%v", got, err, want)
			}
		})
	}
}

func TestTrackedSourceReadersStillRefuseAPathBehindAFile(t *testing.T) {
	t.Parallel()
	for _, reader := range []string{"citation-sources", "history-citations", "stale-history-citations"} {
		t.Run(reader, func(t *testing.T) {
			t.Parallel()
			if !citationSourcesTestHome(t) {
				return
			}
			root := trackedSourceReaderFixture(t)
			notes := filepath.Join(root, "notes")
			if err := os.RemoveAll(notes); err != nil {
				t.Fatal(err)
			}
			if err := os.WriteFile(notes, []byte("blocking file\n"), 0644); err != nil {
				t.Fatal(err)
			}
			_, err := trackedSourceReaderPaths(t, root, reader)
			if err == nil || !strings.Contains(err.Error(), "notes/a.md") {
				t.Fatalf("error=%v; want error naming notes/a.md", err)
			}
		})
	}
}
