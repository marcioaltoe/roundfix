package speccheck_test

// IN: Spec Check and Git/filesystem path scanning in disposable repositories.
// OUT: archive command and settlement policy, covered by their package suites.
// Invariant: active directory dependencies refuse, while archive-resolved Context remains valid.

import (
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"roundfix/internal/gittest"
	"roundfix/internal/speccheck"
)

const pathPinSlug = "0001-path-pin"

func pathPinRepository(t *testing.T, specRoot string, status string, contextPath string) string {
	t.Helper()
	root := t.TempDir()
	gittest.InitRepo(t, root)
	writeCitationFixtureFile(t, root, specRoot+"/"+pathPinSlug+"/_prd.md", "---\nstatus: active\n---\n\n# Path pin\n")
	writeCitationFixtureFile(t, root, specRoot+"/"+pathPinSlug+"/_tasks.md", "---\nschema: spec-tasks/v1\nspec: "+pathPinSlug+"\ngraph:\n  nodes:\n    - id: task_01\n      file: task_01.md\n      needs: []\n---\n")
	task := "---\ntask: task_01\nspec: " + pathPinSlug + "\nstatus: " + status + "\ntype: backend\ncomplexity: low\n---\n\n# Task\n"
	if contextPath != "" {
		task += "\n## Context\n\n- interface: `" + contextPath + "`\n"
	}
	task += "\n## Verification\n\n- `go test ./...`\n"
	writeCitationFixtureFile(t, root, specRoot+"/"+pathPinSlug+"/task_01.md", task)
	return root
}

func TestActiveSpecPathPinIsAnError(t *testing.T) {
	for _, status := range []string{"pending", "in_progress", "completed"} {
		t.Run(status, func(t *testing.T) {
			root := pathPinRepository(t, "docs/specs", status, "")
			writeCitationFixtureFile(t, root, "pin_test.go", "package pin\nvar file = \"docs/specs/"+pathPinSlug+"/_techspec.md\"\n")
			gittest.Run(t, root, "add", "pin_test.go")
			result, err := speccheck.Check(filepath.Join(root, "docs/specs"), root, pathPinSlug)
			if err != nil {
				t.Fatal(err)
			}
			findings := findingsWithCode(result, speccheck.CodeSpecPathPinned)
			if len(findings) != 1 || findings[0].Severity != speccheck.SeverityError || !reflect.DeepEqual(findings[0].Where, []speccheck.Location{{Path: "pin_test.go", Line: 2}}) {
				t.Fatalf("pin findings = %+v", findings)
			}
			if findings[0].Fix != "Read a fixture or an exported constant instead of the Spec's file; a Spec archives and may be deleted." {
				t.Fatalf("fix = %q", findings[0].Fix)
			}
		})
	}
}

func TestActiveSpecPathPinSkipsMarkdownSpecRootsAndIgnoredFiles(t *testing.T) {
	root := pathPinRepository(t, "docs/specs", "pending", "")
	pin := "docs/specs/" + pathPinSlug + "/_techspec.md\n"
	for _, path := range []string{"notes.md", "notes.MD", "docs/specs/other/fixture.go", "docs/history/specs/other/fixture.go", "docs/history/fixture.go", "ignored.go"} {
		writeCitationFixtureFile(t, root, path, pin)
	}
	writeCitationFixtureFile(t, root, ".gitignore", "ignored.go\n")
	writeCitationFixtureFile(t, root, "similar.go", "docs/specs/"+pathPinSlug+"extra/file\ndocs/specs/"+pathPinSlug+"_suffix/file\n")
	writeCitationFixtureFile(t, root, "binary.dat", "\x00"+pin)
	pins, err := speccheck.ActiveSpecPathPins(root, filepath.Join(root, "docs/specs"), pathPinSlug)
	if err != nil || len(pins) != 0 {
		t.Fatalf("pins = %+v, %v", pins, err)
	}
	writeCitationFixtureFile(t, root, "untracked.go", pin+pin)
	pins, err = speccheck.ActiveSpecPathPins(root, filepath.Join(root, "docs/specs"), pathPinSlug)
	if err != nil || !reflect.DeepEqual(pins, []speccheck.SpecPathPin{{Path: "untracked.go", Line: 1}, {Path: "untracked.go", Line: 2}}) {
		t.Fatalf("untracked pins = %+v, %v", pins, err)
	}
}

func TestActiveSpecPathPinsOutsideWorkTree(t *testing.T) {
	root := t.TempDir()
	pin := "planning/" + pathPinSlug + "/file\n"
	for _, path := range []string{"notes.md", "planning/_archived/other/fixture.go", "docs/history/fixture.go", ".git/fixture.go"} {
		writeCitationFixtureFile(t, root, path, pin)
	}
	writeCitationFixtureFile(t, root, "fixture.go", pin)
	pins, err := speccheck.ActiveSpecPathPins(root, filepath.Join(root, "planning"), pathPinSlug)
	if err != nil || !reflect.DeepEqual(pins, []speccheck.SpecPathPin{{Path: "fixture.go", Line: 1}}) {
		t.Fatalf("walk pins = %+v, %v", pins, err)
	}
	pins, err = speccheck.ActiveSpecPathPins(root, t.TempDir(), pathPinSlug)
	if err != nil || len(pins) != 0 {
		t.Fatalf("external pins = %+v, %v", pins, err)
	}
}

func TestTaskContextResolvesThroughTheArchive(t *testing.T) {
	for _, specRoot := range []string{"docs/specs", "planning"} {
		t.Run(specRoot, func(t *testing.T) {
			ref := specRoot + "/0002-other/_techspec.md"
			root := pathPinRepository(t, specRoot, "pending", ref)
			archive := "docs/history/specs"
			if specRoot != "docs/specs" {
				archive = specRoot + "/_archived"
			}
			writeCitationFixtureFile(t, root, archive+"/0002-other/_techspec.md", "# Archived\n")
			result, err := speccheck.Check(filepath.Join(root, specRoot), root, pathPinSlug)
			if err != nil {
				t.Fatal(err)
			}
			if findings := findingsWithCode(result, speccheck.CodeReferenceUnresolved); len(findings) != 0 {
				t.Fatalf("unresolved = %+v", findings)
			}
			// The active form also resolves without an archive; no file is rewritten.
			writeCitationFixtureFile(t, root, ref, "# Active\n")
			if err := os.Remove(filepath.Join(root, archive, "0002-other/_techspec.md")); err != nil {
				t.Fatal(err)
			}
			result, err = speccheck.Check(filepath.Join(root, specRoot), root, pathPinSlug)
			if err != nil || len(findingsWithCode(result, speccheck.CodeReferenceUnresolved)) != 0 {
				t.Fatalf("active resolution = %+v, %v", result.Findings, err)
			}
		})
	}
}

func TestTaskContextMissingInBothPlacesIsUnresolved(t *testing.T) {
	root := pathPinRepository(t, "docs/specs", "pending", "docs/specs/0002-other/missing.md")
	writeCitationFixtureFile(t, root, "docs/history/specs/0002-other/unrelated.md", "# Unrelated\n")
	result, err := speccheck.Check(filepath.Join(root, "docs/specs"), root, pathPinSlug)
	if err != nil {
		t.Fatal(err)
	}
	findings := findingsWithCode(result, speccheck.CodeReferenceUnresolved)
	if len(findings) != 1 || !strings.Contains(findings[0].Summary, "0002-other/missing.md") {
		t.Fatalf("unresolved = %+v", findings)
	}
}
