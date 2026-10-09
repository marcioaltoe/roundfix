package baseline

// Suite: History Relocation citation root containment
// Invariant: a tracked citation read never escapes the repository root.
// Boundary IN: a real Git index and repository paths opened during citation discovery
// Boundary OUT: citation parsing, plan rendering, digest binding, and apply behavior

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func TestRelocationCitationsNeverReadThroughASwappedDirectory(t *testing.T) {
	// Sequential: swaps citationBeforeOpen and citationOpenResult package hooks.
	if runtime.GOOS == "windows" {
		t.Skip("symbolic-link replacement requires Unix permissions")
	}

	repo := newCitationTestRepository(t)
	trackCitationTestFile(t, repo, citationSource, "retired\n")
	const citingPath = "tracked/citation.md"
	trackCitationTestFile(t, repo, citingPath, "in-repository content\n")

	outside := t.TempDir()
	writeCitationTestFile(t, outside, "citation.md", citationSource+"\n")

	previousBeforeOpen := citationBeforeOpen
	previousOpenResult := citationOpenResult
	t.Cleanup(func() {
		citationBeforeOpen = previousBeforeOpen
		citationOpenResult = previousOpenResult
	})

	citationBeforeOpen = func(relative string) {
		if relative != citingPath {
			return
		}
		trackedDirectory := filepath.Join(repo, "tracked")
		if err := os.RemoveAll(trackedDirectory); err != nil {
			t.Fatalf("remove tracked directory: %v", err)
		}
		if err := os.Symlink(outside, trackedDirectory); err != nil {
			t.Fatalf("replace tracked directory with outside link: %v", err)
		}
	}

	var observedOpen bool
	var observedOpenErr error
	citationOpenResult = func(relative string, err error) {
		if relative == citingPath {
			observedOpen = true
			observedOpenErr = err
		}
	}

	findings := citationTestFindings(t, repo, citationTestMove(citationSource), nil)
	if len(findings) != 0 {
		t.Fatalf("relocationCitationFindings() = %#v, want no finding from the outside file", findings)
	}
	if !observedOpen {
		t.Fatalf("open result for %q was not observed", citingPath)
	}
	if observedOpenErr == nil || !strings.Contains(observedOpenErr.Error(), "path escapes from parent") {
		t.Fatalf("open error for %q = %v, want the repository root escape error", citingPath, observedOpenErr)
	}
}

func TestRelocationCitationsReadThroughTheRepositoryRoot(t *testing.T) {
	t.Parallel()
	repo := newCitationTestRepository(t)
	trackCitationTestFile(t, repo, citationSource, "retired\n")
	trackCitationTestFile(t, repo, "guide.md", citationSource+"\n")

	findings := citationTestFindings(t, repo, citationTestMove(citationSource), nil)
	if len(findings) != 1 || findings[0].Path != "guide.md" || !strings.Contains(findings[0].Message, citationSource) {
		t.Fatalf("relocationCitationFindings() = %#v, want the guide's repository-root citation", findings)
	}
}
