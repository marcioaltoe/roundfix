// IN: archive exit status, streams and filesystem in disposable Git repositories.
// OUT: Daemon settlement and live provider calls.
// Invariant: a Glossary Gap prevents every archive, including QA overrides.
package cli

import (
	"bytes"
	"context"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"roundfix/internal/spec"
)

func glossaryArchiveWorkspace(t *testing.T) (string, string) {
	t.Helper()
	_, root := newImplementWorkspace(t, []implementSeed{{id: "task_01", title: "Write the glossary", status: string(spec.StatusCompleted)}})
	writeArchiveQAReport(t, root, spec.VerdictPass, "rows_blocked_environment: 3")
	dir := filepath.Join(root, "docs/specs", implementTestSlug)
	prd := mustRead(t, filepath.Join(dir, "_prd.md"))
	mustWrite(t, filepath.Join(dir, "_prd.md"), prd+"\n## Glossary\n\nNone.\n\n## Story\n\n**Missing Term**\n")
	mustWrite(t, filepath.Join(root, "CONTEXT.md"), "**Existing Term**: definition\n")
	return root, dir
}
func testArchiveGlossaryRefusal(t *testing.T, override bool) {
	t.Helper()
	root, dir := glossaryArchiveWorkspace(t)
	if override {
		writeArchiveQAReport(t, root, spec.VerdictFail)
	}
	before := snapshotDirectoryFiles(t, root)
	args := []string{"archive", implementTestSlug}
	if override {
		args = append(args, "--qa-override", "--approval", "maintainer test approval", "--reason", "waive QA")
	}
	var stdout, stderr bytes.Buffer
	code := runCLIContext(t, context.Background(), args, &stdout, &stderr)
	want := "Spec \"" + implementTestSlug + "\" cannot archive with a Glossary Gap: SC-GLOSSARY-UNDECLARED: "
	if code != 2 || stdout.Len() != 0 || !strings.Contains(stderr.String(), want) || !strings.Contains(stderr.String(), "Missing Term") {
		t.Fatalf("exit=%d stdout=%q stderr=%q", code, stdout.String(), stderr.String())
	}
	assertPathExists(t, dir)
	assertPathMissing(t, archiveTestRepositoryPath(root, spec.ArchiveKindSpec, implementTestSlug))
	if after := snapshotDirectoryFiles(t, root); !reflect.DeepEqual(after, before) {
		t.Fatal("refusal changed repository files")
	}
}
func TestArchiveRefusesASpecWithAGlossaryGap(t *testing.T) { testArchiveGlossaryRefusal(t, false) }
func TestArchiveWithQAOverrideStillRefusesAGlossaryGap(t *testing.T) {
	testArchiveGlossaryRefusal(t, true)
}
func TestArchiveAcceptsASpecWhoseGlossaryIsCurrent(t *testing.T) {
	root, dir := glossaryArchiveWorkspace(t)
	mustWrite(t, filepath.Join(root, "CONTEXT.md"), "**Missing Term**: definition\n")
	var stdout, stderr bytes.Buffer
	commitArchiveFixture(t)
	code := runCLIContext(t, context.Background(), []string{"archive", implementTestSlug}, &stdout, &stderr)
	if code != 0 || stderr.Len() != 0 || !strings.Contains(stdout.String(), "archived "+implementTestSlug) {
		t.Fatalf("exit=%d stdout=%q stderr=%q", code, stdout.String(), stderr.String())
	}
	assertPathMissing(t, dir)
	assertPathExists(t, archiveTestRepositoryPath(root, spec.ArchiveKindSpec, implementTestSlug)+".md")
}
