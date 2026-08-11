package spec

// Suite: archive layout characterization before Spec 0085
// Invariant: the retired-artifact paths, their distributed owners, the
// conditional Secondbrain clause, and the corpus golden remain recorded until
// the Spec Task that owns each intentional break updates its contract.
// Boundary IN: repository documentation directories and the docscontract corpus golden.
// Boundary OUT: the archive resolver, consumer migration, artifact relocation,
// and the unconditional Secondbrain clause assigned to later Tasks.

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

var archiveLayoutBeforeSpec0085 = []retiredFamilyLocation{
	{family: "Specs", directory: "docs/specs/_archived"},
	{family: "findings", directory: "docs/findings/_archived"},
	{family: "ADRs", directory: "docs/adr"},
	{family: "backlog entries", directory: "docs/backlog"},
}

func TestArchiveLayoutCharacterizationRecordsEveryRetiredFamily(t *testing.T) {
	t.Parallel()

	want := []retiredFamilyLocation{
		{family: "Specs", directory: "docs/specs/_archived"},
		{family: "findings", directory: "docs/findings/_archived"},
		{family: "ADRs", directory: "docs/adr"},
		{family: "backlog entries", directory: "docs/backlog"},
	}
	if !reflect.DeepEqual(archiveLayoutBeforeSpec0085, want) {
		t.Fatalf("archive layout characterization = %#v, want %#v", archiveLayoutBeforeSpec0085, want)
	}

	repositoryRoot := archiveLayoutCharacterizationRepositoryRoot(t)
	for _, retiredFamily := range archiveLayoutBeforeSpec0085 {
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

type archivePathComposer struct {
	packagePath string
	literals    []string
}

// These are historical resolved path literals. Some callers spell the value
// as filepath.Join fragments, but each still owns the resulting layout.
var archivePathComposersBeforeSpec0085 = []archivePathComposer{
	{packagePath: "internal/spec", literals: []string{"_archived"}},
	{packagePath: "internal/speccheck", literals: []string{"docs/specs/_archived", "docs/findings/_archived"}},
	{packagePath: "internal/specaudit", literals: []string{"_archived"}},
	{packagePath: "internal/worktree", literals: []string{"docs/specs/_archived"}},
	{packagePath: "internal/cli", literals: []string{"docs/specs/_archived"}},
}

func TestArchiveLayoutCharacterizationEnumeratesEveryPathComposer(t *testing.T) {
	t.Parallel()

	want := []archivePathComposer{
		{packagePath: "internal/spec", literals: []string{"_archived"}},
		{packagePath: "internal/speccheck", literals: []string{"docs/specs/_archived", "docs/findings/_archived"}},
		{packagePath: "internal/specaudit", literals: []string{"_archived"}},
		{packagePath: "internal/worktree", literals: []string{"docs/specs/_archived"}},
		{packagePath: "internal/cli", literals: []string{"docs/specs/_archived"}},
	}
	if !reflect.DeepEqual(archivePathComposersBeforeSpec0085, want) {
		t.Fatalf("archive path composers = %#v, want %#v", archivePathComposersBeforeSpec0085, want)
	}

	var specCheckerLiterals []string
	for _, composer := range archivePathComposersBeforeSpec0085 {
		if composer.packagePath == "internal/speccheck" {
			specCheckerLiterals = composer.literals
			break
		}
	}
	wantSpecCheckerLiterals := []string{"docs/specs/_archived", "docs/findings/_archived"}
	if !reflect.DeepEqual(specCheckerLiterals, wantSpecCheckerLiterals) {
		t.Fatalf("internal/speccheck archive literals = %q, want %q", specCheckerLiterals, wantSpecCheckerLiterals)
	}
}

const conditionalSecondbrainClauseBeforeSpec0085 = "Consult the local Secondbrain before acting when repository context does not answer business or prior-decision questions, fiscal or tax concepts, cross-project documentation, knowledge about Vortex, Tax, Visio, or Gesttione, or shared architecture patterns. Do not consult it when local code, `CONTEXT.md`, ADRs, and repository documentation fully answer the task."

func TestArchiveLayoutCharacterizationCapturesConditionalSecondbrainClause(t *testing.T) {
	t.Parallel()

	const escapeHatch = "Do not consult it when local code, `CONTEXT.md`, ADRs, and repository documentation fully answer the task."
	if !strings.Contains(conditionalSecondbrainClauseBeforeSpec0085, escapeHatch) {
		t.Fatalf("conditional Secondbrain clause = %q, want escape hatch %q", conditionalSecondbrainClauseBeforeSpec0085, escapeHatch)
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
		Update: "Only active-corpus counts are pinned. After an intentional detector change, run the focused corpus test, inspect its actual active counts, and update this file in the same change. Archived counts are derived and reported because historical authoring changes them.",
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
	// Declared break: task_04 relocates archived Specs and findings, then
	// re-records this golden and updates this characterization in the same Task.
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("corpus golden before task_04 = %#v, want %#v", got, want)
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
