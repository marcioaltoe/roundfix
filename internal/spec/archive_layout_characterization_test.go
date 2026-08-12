package spec

// Suite: archive layout characterization through Spec 0085
// Invariant: each retired-artifact path, its single owner, the conditional
// Secondbrain clause, and the deliberately re-recorded corpus golden remain
// explicit while the Spec moves each contract in its assigned Task.
// Boundary IN: repository documentation directories and the docscontract corpus golden.
// Boundary OUT: retired ADR relocation and the unconditional Secondbrain clause
// assigned to later Tasks.

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
	family    string
	directory string
}

var archiveLayoutAfterTask04 = []retiredFamilyLocation{
	{family: "Specs", directory: "_archived/specs"},
	{family: "findings", directory: "_archived/findings"},
	{family: "ADRs", directory: "docs/adr"},
	{family: "backlog entries", directory: "docs/backlog"},
}

func TestArchiveLayoutCharacterizationRecordsEveryRetiredFamily(t *testing.T) {
	t.Parallel()

	want := []retiredFamilyLocation{
		{family: "Specs", directory: "_archived/specs"},
		{family: "findings", directory: "_archived/findings"},
		{family: "ADRs", directory: "docs/adr"},
		{family: "backlog entries", directory: "docs/backlog"},
	}
	if !reflect.DeepEqual(archiveLayoutAfterTask04, want) {
		t.Fatalf("archive layout characterization = %#v, want %#v", archiveLayoutAfterTask04, want)
	}

	repositoryRoot := archiveLayoutCharacterizationRepositoryRoot(t)
	for _, retiredFamily := range archiveLayoutAfterTask04 {
		t.Run(retiredFamily.family, func(t *testing.T) {
			t.Parallel()

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

func TestArchiveLayoutCharacterizationPinsCorpusGoldenBeforeRelocation(t *testing.T) {
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
		Update: "Re-recorded because Spec 0085 Task 04 moved archived Specs and findings from their active-document trees to _archived/specs and _archived/findings. Active-corpus counts remain unchanged because retired artifacts are excluded. After an intentional detector change, run the focused corpus test, inspect its actual active counts, and update this file in the same change.",
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
			"SC-VERIFY-WORK-INDEPENDENT":   0,
			"SC-VOCABULARY-UNDOCUMENTED":   0,
		},
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("corpus golden after task_04 = %#v, want %#v", got, want)
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
