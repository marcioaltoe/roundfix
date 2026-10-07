package cli

// IN: archive CLI streams, exit code and filesystem in temporary repositories.
// OUT: pin matching details, covered by internal/speccheck.
// Invariant: archive moves nothing on a code pin and accepts Markdown references.

import (
	"bytes"
	"context"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"roundfix/internal/spec"
)

func archivePinWorkspace(t *testing.T) (string, string) {
	t.Helper()
	_, root := newImplementWorkspace(t, nil)
	dir := filepath.Join(root, "docs/specs", implementTestSlug)
	if err := os.Remove(filepath.Join(dir, "_tasks.md")); err != nil {
		t.Fatal(err)
	}
	mustWrite(t, filepath.Join(dir, spec.SupersessionFilename), "---\nsuperseded_by: 0002-delivered-widget\ndate: 2026-10-03\nreason: delivered elsewhere\n---\n\nDelivered elsewhere.\n")
	commitArchiveFixture(t)
	return root, dir
}

func TestArchiveRefusesASpecAFileStillPins(t *testing.T) {
	root, dir := archivePinWorkspace(t)
	before := snapshotDirectoryFiles(t, dir)
	mustWrite(t, filepath.Join(root, "pin_test.go"), "package pin\nvar path = \"docs/specs/"+implementTestSlug+"/_techspec.md\"\n")
	var stdout, stderr bytes.Buffer
	code := runCLIContext(t, context.Background(), []string{"archive", implementTestSlug}, &stdout, &stderr)
	reason := "Spec \"" + implementTestSlug + "\" cannot archive while another file names its active directory docs/specs/" + implementTestSlug + "/: pin_test.go:2"
	if code != exitPreflight || stdout.Len() != 0 || !strings.Contains(stderr.String(), reason) {
		t.Fatalf("exit=%d stdout=%q stderr=%q", code, stdout.String(), stderr.String())
	}
	if after := snapshotDirectoryFiles(t, dir); !reflect.DeepEqual(after, before) {
		t.Fatalf("refusal changed Spec: %+v", after)
	}
	assertPathMissing(t, archiveTestRepositoryPath(root, spec.ArchiveKindSpec, implementTestSlug))
	if err := os.Remove(filepath.Join(root, "pin_test.go")); err != nil {
		t.Fatal(err)
	}
	stdout.Reset()
	stderr.Reset()
	code = runCLIContext(t, context.Background(), []string{"archive", implementTestSlug}, &stdout, &stderr)
	if code != exitOK {
		t.Fatalf("pin removed: exit=%d stderr=%q", code, stderr.String())
	}
	assertPathMissing(t, dir)
	if after := archivedSourceFiles(t, root, archiveTestRepositoryPath(root, spec.ArchiveKindSpec, implementTestSlug)+".md"); !reflect.DeepEqual(after, before) {
		t.Fatal("archive changed superseded files")
	}
}

func TestArchiveIgnoresMarkdownThatNamesTheSpec(t *testing.T) {
	root, dir := archivePinWorkspace(t)
	mustWrite(t, filepath.Join(root, "notes.md"), "docs/specs/"+implementTestSlug+"/_techspec.md\n")
	var stdout, stderr bytes.Buffer
	code := runCLIContext(t, context.Background(), []string{"archive", implementTestSlug}, &stdout, &stderr)
	if code != exitOK || stderr.Len() != 0 {
		t.Fatalf("exit=%d stderr=%q", code, stderr.String())
	}
	assertPathMissing(t, dir)
}
