package spec

// Suite: archive layout characterization through Spec 0094 Task 02
// Invariant: each retired-artifact path, its single owner, the conditional
// Secondbrain clause, and the deliberately re-recorded corpus golden remain
// explicit while the Spec moves each contract in its assigned Task.
// Boundary IN: repository documentation directories and the docscontract corpus golden.
// Boundary OUT: retirement classification, fleet relocation, and documentation
// carriers assigned to other Tasks.

import (
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"
	"testing"
)

type retiredFamilyLocation struct {
	family              string
	kind                ArchiveKind
	directory           string
	presentInRepository bool
}

var archiveLayoutAfterSpec0094Task02 = []retiredFamilyLocation{
	{family: "Specs", kind: ArchiveKindSpec, directory: "docs/history/specs", presentInRepository: true},
	{family: "findings", kind: ArchiveKindFinding, directory: "docs/history/findings", presentInRepository: true},
	{family: "ADRs", kind: ArchiveKindADR, directory: "docs/history/adr", presentInRepository: true},
	{family: "backlog entries", kind: ArchiveKindBacklog, directory: "docs/history/backlog"},
	{family: "Review Artifacts", kind: ArchiveKindReview, directory: "docs/history/reviews"},
}

func TestArchiveLayoutCharacterizationRecordsEveryRetiredFamily(t *testing.T) {
	t.Parallel()

	want := []retiredFamilyLocation{
		{family: "Specs", kind: ArchiveKindSpec, directory: "docs/history/specs", presentInRepository: true},
		{family: "findings", kind: ArchiveKindFinding, directory: "docs/history/findings", presentInRepository: true},
		{family: "ADRs", kind: ArchiveKindADR, directory: "docs/history/adr", presentInRepository: true},
		{family: "backlog entries", kind: ArchiveKindBacklog, directory: "docs/history/backlog"},
		{family: "Review Artifacts", kind: ArchiveKindReview, directory: "docs/history/reviews"},
	}
	if !reflect.DeepEqual(archiveLayoutAfterSpec0094Task02, want) {
		t.Fatalf("archive layout characterization = %#v, want %#v", archiveLayoutAfterSpec0094Task02, want)
	}

	repositoryRoot := archiveLayoutCharacterizationRepositoryRoot(t)
	for _, retiredFamily := range archiveLayoutAfterSpec0094Task02 {
		t.Run(retiredFamily.family, func(t *testing.T) {
			t.Parallel()

			if got := ArchiveDir(retiredFamily.kind); got != retiredFamily.directory {
				t.Fatalf("ArchiveDir(%q) = %q, want %q", retiredFamily.kind, got, retiredFamily.directory)
			}
			if !retiredFamily.presentInRepository {
				return
			}
			info, err := os.Stat(filepath.Join(repositoryRoot, filepath.FromSlash(retiredFamily.directory)))
			if err != nil {
				t.Fatalf("stat current %s directory %q: %v", retiredFamily.family, retiredFamily.directory, err)
			}
			if !info.IsDir() {
				t.Fatalf("current %s location %q is not a directory", retiredFamily.family, retiredFamily.directory)
			}
		})
	}
}

type archiveLayoutCorpusGolden struct {
	Schema string         `json:"schema"`
	Update string         `json:"update"`
	Active map[string]int `json:"active"`
}

func TestArchiveLayoutCharacterizationPinsCorpusGoldenAfterSpec0095(t *testing.T) {
	t.Parallel()

	repositoryRoot := archiveLayoutCharacterizationRepositoryRoot(t)
	path := filepath.Join(repositoryRoot, "internal", "docscontract", "testdata", "corpus-golden.json")
	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read corpus golden %q: %v", path, err)
	}
	var got archiveLayoutCorpusGolden
	decoder := json.NewDecoder(strings.NewReader(string(content)))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&got); err != nil {
		t.Fatalf("parse corpus golden %q: %v", path, err)
	}

	want := archiveLayoutCorpusGolden{
		Schema: "roundfix-speccheck-corpus/v2",
		Update: "Re-recorded because Spec 0095 restored SC-VERIFY-VACUOUS-COMMAND at the tasks stage and added all three of its Verification refusal codes to the characterization corpus. The active corpus reports zero findings for SC-VERIFY-INVERTED-EXIT, SC-VERIFY-NON-HERMETIC, and SC-VERIFY-VACUOUS-COMMAND after their authored commands were accounted. After an intentional detector change, run the focused corpus test, inspect its actual active counts, and update this file in the same change.",
		Active: map[string]int{
			"SC-ADR-RELATED":               0,
			"SC-ADR-UNLISTED":              0,
			"SC-ARCHIVE-LICENSE":           0,
			"SC-CITATION-UNSUPPORTED":      0,
			"SC-CONSTRAINT-MISSING":        0,
			"SC-CONSTRAINT-SOURCE":         0,
			"SC-CONSTRAINT-UNREASONED":     0,
			"SC-COVERAGE-UNMAPPED":         0,
			"SC-COVERAGE-UNTASKED":         0,
			"SC-FINDING-LIFECYCLE":         0,
			"SC-LOOP-ORDER-DIVERGENT":      0,
			"SC-REF-UNRESOLVED":            0,
			"SC-REHEARSAL-UNDECLARED":      0,
			"SC-REQUIREMENT-CONTRADICTORY": 0,
			"SC-ROLLUP-MEMBER":             0,
			"SC-TOOLING-UNAUTHORIZED":      0,
			"SC-TOOLING-UNBOUNDED":         0,
			"SC-VERIFY-INVERTED-EXIT":      0,
			"SC-VERIFY-NON-HERMETIC":       0,
			"SC-VERIFY-VACUOUS-COMMAND":    0,
			"SC-VERIFY-WORK-INDEPENDENT":   0,
			"SC-VOCABULARY-UNDOCUMENTED":   0,
		},
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("corpus golden after Spec 0095 = %#v, want %#v", got, want)
	}
}

func archiveLayoutCharacterizationRepositoryRoot(t *testing.T) string {
	t.Helper()

	_, testFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("runtime.Caller could not locate the repository")
	}
	return filepath.Clean(filepath.Join(filepath.Dir(testFile), "..", ".."))
}
