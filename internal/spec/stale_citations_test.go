package spec

// Suite: Stale History Citations.
// Invariant: only absent files held by the tag are reported, outside the History Root.
// Boundary IN: tracked Markdown and a tag-file set in disposable Git repositories.
// Boundary OUT: baseline state and rendering are tested in internal/cli.

import (
	"context"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"roundfix/internal/gittest"
)

func TestStaleHistoryCitationsReportsFilesTheTagHolds(t *testing.T) {
	t.Parallel()
	root := t.TempDir()
	gittest.InitRepo(t, root, "-b", "main")
	writeStaleCitationSource(t, root, "docs/findings/gap.md", "# Gap\n\n[lesson](../history/specs/legacy/references/lesson.md#source)\nSee `docs/history/specs/legacy/_prd.md`.\n")
	gittest.Run(t, root, "add", ".")
	held := map[string]bool{"docs/history/specs/legacy/references/lesson.md": true, "docs/history/specs/legacy/_prd.md": true}
	got, err := StaleHistoryCitations(context.Background(), root, "docs/history", held)
	want := []HistoryCitation{{Path: "docs/findings/gap.md", Line: 3, Target: "docs/history/specs/legacy/references/lesson.md"}, {Path: "docs/findings/gap.md", Line: 4, Target: "docs/history/specs/legacy/_prd.md"}}
	if err != nil || !reflect.DeepEqual(got, want) {
		t.Fatalf("citations=%+v err=%v want=%+v", got, err, want)
	}
}

func TestStaleHistoryCitationsIgnoresDirectoriesAndPresentFiles(t *testing.T) {
	t.Parallel()
	root := t.TempDir()
	gittest.InitRepo(t, root, "-b", "main")
	writeStaleCitationSource(t, root, "docs/findings/gap.md", "[folder](../history/specs/legacy/)\n[still here](../history/specs/legacy/present.md)\n`docs/history/specs/legacy/`\n`docs/history/specs/legacy/present.md`\n`docs/history/specs/legacy/missing.md-extra`\n`docs/history/specs/legacy/not-held.md`\n")
	writeStaleCitationSource(t, root, "docs/history/specs/legacy/present.md", "[missing](missing.md)\n")
	writeStaleCitationSource(t, root, "untracked.md", "docs/history/specs/legacy/missing.md\n")
	gittest.Run(t, root, "add", "docs")
	held := map[string]bool{"docs/history/specs/legacy/present.md": true, "docs/history/specs/legacy/missing.md": true}
	got, err := StaleHistoryCitations(context.Background(), root, "docs/history", held)
	if err != nil || len(got) != 0 {
		t.Fatalf("citations=%+v err=%v", got, err)
	}
}

func writeStaleCitationSource(t *testing.T, root, path, content string) {
	t.Helper()
	absolute := filepath.Join(root, path)
	if err := os.MkdirAll(filepath.Dir(absolute), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(absolute, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
}
