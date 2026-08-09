// Suite: Task authoring coherence and stage-scoped Spec checks
// Invariant: coherence findings retain their declared meaning, and a stage runs exactly the detectors it can decide.
// Boundary IN: parsed Task declarations and the public coherence and stage-scope APIs
// Boundary OUT: Markdown parsing, CLI rendering, and Daemon Verification execution
package speccheck_test

import (
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"roundfix/internal/spec"
	"roundfix/internal/speccheck"
)

func TestStageScopeRunsOnlyDetectorsTheStageCanDecide(t *testing.T) {
	t.Parallel()

	repoRoot, specsRoot, slug := writeStageScopeFixture(t)
	prdResult, err := speccheck.CheckStage(specsRoot, repoRoot, slug, speccheck.StagePRD)
	if err != nil {
		t.Fatalf("CheckStage(StagePRD): %v", err)
	}
	if len(findingsWithCode(prdResult, speccheck.CodeConstraintMissing)) == 0 {
		t.Fatalf("StagePRD findings = %#v, want PRD constraint detector findings", prdResult.Findings)
	}
	for _, code := range []string{
		speccheck.CodeCoverageUnmapped,
		speccheck.CodeCoverageUntasked,
		speccheck.CodeVerifyWorkIndependent,
	} {
		if findings := findingsWithCode(prdResult, code); len(findings) != 0 {
			t.Errorf("StagePRD %s findings = %#v, want detector not run", code, findings)
		}
	}

	techSpecResult, err := speccheck.CheckStage(fixtureSpecRoot, "testdata/repo", "vocabulary-missing", speccheck.StageTechSpec)
	if err != nil {
		t.Fatalf("CheckStage(StageTechSpec): %v", err)
	}
	if len(findingsWithCode(techSpecResult, speccheck.CodeVocabularyUndocumented)) != 1 {
		t.Fatalf("StageTechSpec %s findings = %#v, want the TechSpec detector finding", speccheck.CodeVocabularyUndocumented, techSpecResult.Findings)
	}
	for _, code := range []string{
		speccheck.CodeCoverageUntasked,
		speccheck.CodeVerifyWorkIndependent,
		speccheck.CodeRequirementContradictory,
	} {
		if findings := findingsWithCode(techSpecResult, code); len(findings) != 0 {
			t.Errorf("StageTechSpec %s findings = %#v, want detector not run", code, findings)
		}
	}

	tasksResult, err := speccheck.CheckStage(specsRoot, repoRoot, slug, speccheck.StageTasks)
	if err != nil {
		t.Fatalf("CheckStage(StageTasks): %v", err)
	}
	fullResult, err := speccheck.Check(specsRoot, repoRoot, slug)
	if err != nil {
		t.Fatalf("Check(): %v", err)
	}
	if !reflect.DeepEqual(tasksResult, fullResult) {
		t.Errorf("CheckStage(StageTasks) = %#v, want full detector sweep %#v", tasksResult, fullResult)
	}
}

func TestStageScopeRunsTheCitationDetectorInAuthoringStages(t *testing.T) {
	t.Parallel()

	repoRoot := citationCharacterizationFixtureRoot(t)
	specsRoot := filepath.Join(repoRoot, "docs", "specs")
	tests := []struct {
		name  string
		stage speccheck.Stage
	}{
		{name: "PRD", stage: speccheck.StagePRD},
		{name: "TechSpec", stage: speccheck.StageTechSpec},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			result, err := speccheck.CheckStage(specsRoot, repoRoot, citationCharacterizationSlug, tt.stage)
			if err != nil {
				t.Fatalf("CheckStage(%s): %v", tt.stage, err)
			}
			findings := findingsWithCode(result, speccheck.CodeCitationUnsupported)
			if len(findings) != 1 {
				t.Fatalf("CheckStage(%s) %s findings = %#v, want one", tt.stage, speccheck.CodeCitationUnsupported, findings)
			}
			if findings[0].Severity != speccheck.SeverityError {
				t.Errorf("CheckStage(%s) severity = %q, want %q", tt.stage, findings[0].Severity, speccheck.SeverityError)
			}
			for _, want := range []string{
				"ADR-0083 makes `make verify` the only authoritative gate",
				"Adopted sources move to their owning Spec",
			} {
				if !strings.Contains(findings[0].Summary, want) {
					t.Errorf("CheckStage(%s) summary = %q, want %q", tt.stage, findings[0].Summary, want)
				}
			}
		})
	}
}

func TestStageScopeDefaultSweepIsUnchanged(t *testing.T) {
	repoRoot := characterizationRepositoryRoot(t)
	activeRoot := filepath.Join(repoRoot, "docs", "specs")
	archivedRoot := materializeArchivedCorpus(t, filepath.Join(activeRoot, "_archived"))
	for _, specsRoot := range []string{activeRoot, archivedRoot} {
		assertDefaultStageSweepUnchanged(t, specsRoot, repoRoot)
	}
}

func assertDefaultStageSweepUnchanged(t *testing.T, specsRoot, repoRoot string) {
	t.Helper()

	entries, err := os.ReadDir(specsRoot)
	if err != nil {
		t.Fatalf("read Spec corpus %q: %v", specsRoot, err)
	}
	for _, entry := range entries {
		if !entry.IsDir() || strings.HasPrefix(entry.Name(), "_") {
			continue
		}
		if _, err := os.Stat(filepath.Join(specsRoot, entry.Name(), "_prd.md")); errors.Is(err, os.ErrNotExist) {
			continue
		} else if err != nil {
			t.Fatalf("inspect Spec %q: %v", entry.Name(), err)
		}
		t.Run(filepath.Base(specsRoot)+"/"+entry.Name(), func(t *testing.T) {
			fullResult, err := speccheck.Check(specsRoot, repoRoot, entry.Name())
			if err != nil {
				t.Fatalf("Check(): %v", err)
			}
			defaultResult, err := speccheck.CheckStage(specsRoot, repoRoot, entry.Name(), speccheck.StageAll)
			if err != nil {
				t.Fatalf("CheckStage(StageAll): %v", err)
			}
			if !reflect.DeepEqual(defaultResult, fullResult) {
				t.Errorf("default stage result = %#v, want unchanged result %#v", defaultResult, fullResult)
			}
		})
	}
}

func TestStageScopeNamesTheDetectorsItSkipped(t *testing.T) {
	t.Parallel()

	repoRoot, specsRoot, slug := writeStageScopeFixture(t)
	result, err := speccheck.CheckStage(specsRoot, repoRoot, slug, speccheck.StagePRD)
	if err != nil {
		t.Fatalf("CheckStage(StagePRD): %v", err)
	}
	for _, code := range []string{
		speccheck.CodeCoverageUnmapped,
		speccheck.CodeVocabularyUndocumented,
		speccheck.CodeADRUnlisted,
		speccheck.CodeADRRelated,
		speccheck.CodeCoverageUntasked,
		speccheck.CodeReferenceUnresolved,
		speccheck.CodeVerifyWorkIndependent,
		speccheck.CodeRequirementContradictory,
		speccheck.CodeRehearsalUndeclared,
	} {
		if !hasSkip(result, code, "stage prd") {
			t.Errorf("StagePRD Skipped = %#v, want named %s detector", result.Skipped, code)
		}
	}
}

func TestStageScopeRejectsUnknownStage(t *testing.T) {
	t.Parallel()

	repoRoot, specsRoot, slug := writeStageScopeFixture(t)
	_, err := speccheck.CheckStage(specsRoot, repoRoot, slug, speccheck.Stage("draft"))
	if err == nil {
		t.Fatal("CheckStage(unknown) error = nil, want accepted stage values")
	}
	for _, want := range []string{"prd", "techspec", "tasks"} {
		if !strings.Contains(err.Error(), want) {
			t.Errorf("CheckStage(unknown) error = %q, want accepted value %q", err, want)
		}
	}
}

func writeStageScopeFixture(t *testing.T) (string, string, string) {
	t.Helper()

	const slug = "stage-scope"
	repoRoot := t.TempDir()
	specsRoot := filepath.Join(repoRoot, "docs", "specs")
	path := filepath.Join(specsRoot, slug, "_prd.md")
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("create stage-scope fixture: %v", err)
	}
	if err := os.WriteFile(path, []byte("# Stage scope\n"), 0o644); err != nil {
		t.Fatalf("write stage-scope fixture: %v", err)
	}
	return repoRoot, specsRoot, slug
}

func TestCoherenceFindingsRetainDeclarationSourceLines(t *testing.T) {
	t.Parallel()

	graph, err := spec.Load(filepath.Join("testdata", "repo", "docs", "specs"), replay0060Task03)
	if err != nil {
		t.Fatalf("Load(%s): %v", replay0060Task03, err)
	}
	if len(graph.Tasks) != 1 {
		t.Fatalf("Tasks = %d, want 1", len(graph.Tasks))
	}
	task := graph.Tasks[0]

	contradiction, ok := speccheck.ContradictoryRequirements(task)
	if !ok {
		t.Fatal("ContradictoryRequirements() found no contradiction")
	}
	wantContradictionLines := []int{15, 13}
	for index, want := range wantContradictionLines {
		if got := contradiction.Where[index].Line; got != want {
			t.Errorf("contradiction location %d line = %d, want %d", index, got, want)
		}
	}

	rehearsal, ok := speccheck.UndeclaredRehearsal(task)
	if !ok {
		t.Fatal("UndeclaredRehearsal() found no missing declaration")
	}
	if got, want := rehearsal.Where[0].Line, 9; got != want {
		t.Errorf("rehearsal location line = %d, want %d", got, want)
	}
}

func TestContradictoryRequirementsRefusesSameNamedSubject(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		requirements []spec.TaskDeclaration
		wantFinding  bool
	}{
		{
			name: "same subject is required and forbidden",
			requirements: []spec.TaskDeclaration{
				{Text: "MUST keep the named gate enabled.", Line: 11},
				{Text: "MUST NOT keep the named gate enabled.", Line: 12},
			},
			wantFinding: true,
		},
		{
			name: "subject cannot be identified",
			requirements: []spec.TaskDeclaration{
				{Text: "MUST publish the report.", Line: 11},
				{Text: "MUST NOT delete the scratch directory.", Line: 12},
			},
		},
		{
			name: "one sided requirement",
			requirements: []spec.TaskDeclaration{
				{Text: "MUST keep the named gate enabled.", Line: 11},
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			finding, got := speccheck.ContradictoryRequirements(spec.Task{
				File:         "fixture/task_01.md",
				Requirements: tt.requirements,
			})
			if got != tt.wantFinding {
				t.Fatalf("ContradictoryRequirements() finding = %v, want %v; finding: %#v", got, tt.wantFinding, finding)
			}
			if !got {
				return
			}
			if finding.Code != speccheck.CodeRequirementContradictory || finding.Severity != speccheck.SeverityError {
				t.Errorf("finding identity = %s/%s, want %s/%s", finding.Code, finding.Severity, speccheck.CodeRequirementContradictory, speccheck.SeverityError)
			}
			if len(finding.Where) != 2 {
				t.Errorf("finding locations = %#v, want both declaring requirements", finding.Where)
			}
		})
	}
}

func TestUndeclaredRehearsalRequiresCasesAndObservations(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		title          string
		titleLine      int
		rehearsalCases []spec.TaskDeclaration
		wantFinding    bool
	}{
		{
			name:        "rehearsal has no cases",
			title:       "Rehearse the lifecycle and prove the gates",
			titleLine:   9,
			wantFinding: true,
		},
		{
			name:        "stated purpose proves a gate fires",
			title:       "Prove the archive gate fires",
			titleLine:   9,
			wantFinding: true,
		},
		{
			name:  "rehearsal declares cases and observations",
			title: "Rehearse the lifecycle and prove the gates",
			rehearsalCases: []spec.TaskDeclaration{
				{Text: "Case: stale link; Observation: archive reports the link.", Line: 20},
				{Text: "Case: prose mention; Observation: archive remains accepted.", Line: 21},
			},
		},
		{
			name:  "case without observation is incomplete",
			title: "Rehearse the lifecycle and prove the gates",
			rehearsalCases: []spec.TaskDeclaration{
				{Text: "Case: stale link", Line: 20},
			},
			wantFinding: true,
		},
		{
			name:  "ordinary task needs no rehearsal declaration",
			title: "Implement the gate detector",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			finding, got := speccheck.UndeclaredRehearsal(spec.Task{
				File:           "fixture/task_01.md",
				Title:          tt.title,
				TitleLine:      tt.titleLine,
				RehearsalCases: tt.rehearsalCases,
			})
			if got != tt.wantFinding {
				t.Fatalf("UndeclaredRehearsal() finding = %v, want %v; finding: %#v", got, tt.wantFinding, finding)
			}
			if !got {
				return
			}
			if finding.Code != speccheck.CodeRehearsalUndeclared || finding.Severity != speccheck.SeverityError {
				t.Errorf("finding identity = %s/%s, want %s/%s", finding.Code, finding.Severity, speccheck.CodeRehearsalUndeclared, speccheck.SeverityError)
			}
			if !strings.Contains(finding.Fix, "Rehearsal Cases") {
				t.Errorf("finding fix = %q, want the authored section name", finding.Fix)
			}
		})
	}
}
