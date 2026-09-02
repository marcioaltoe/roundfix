// Suite: Task authoring coherence and stage-scoped Spec checks
// Invariant: coherence findings retain their declared meaning, and a stage runs exactly the detectors it can decide.
// Boundary IN: parsed Task declarations and the public coherence and stage-scope APIs
// Boundary OUT: Markdown parsing, CLI rendering, and Daemon Verification execution
package speccheck_test

import (
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
		speccheck.CodeVerifyInvertedExit,
		speccheck.CodeVerifyNonHermetic,
		speccheck.CodeWaveCollision,
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
		speccheck.CodeVerifyInvertedExit,
		speccheck.CodeVerifyNonHermetic,
		speccheck.CodeRequirementContradictory,
		speccheck.CodeWaveCollision,
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
		speccheck.CodeVerifyInvertedExit,
		speccheck.CodeVerifyNonHermetic,
		speccheck.CodeRequirementContradictory,
		speccheck.CodeRehearsalUndeclared,
		speccheck.CodeWaveCollision,
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

func TestWaveCollisionIsReportedAtAuthoring(t *testing.T) {
	t.Parallel()

	const slug = "wave-collision"
	tests := []struct {
		name          string
		manifestNeeds string
		writeGraph    bool
		wantFinding   bool
	}{
		{
			name:        "colliding graph",
			writeGraph:  true,
			wantFinding: true,
		},
		{
			name:          "serialized graph",
			manifestNeeds: "task_01",
			writeGraph:    true,
		},
		{
			name: "Spec with no graph yet",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			repoRoot := t.TempDir()
			specsRoot := filepath.Join(repoRoot, "docs", "specs")
			writeCitationFixtureFile(t, repoRoot, "internal/shared/first.go", "package shared\n")
			writeCitationFixtureFile(t, repoRoot, "internal/shared/second.go", "package shared\n")
			writeCitationFixtureFile(t, repoRoot, "docs/specs/"+slug+"/_prd.md", `---
spec: wave-collision
status: active
created: 2026-09-02
surfaces: [backend]
---

# Wave collision fixture
`)
			if tt.writeGraph {
				needs := "[]"
				if tt.manifestNeeds != "" {
					needs = "[" + tt.manifestNeeds + "]"
				}
				writeCitationFixtureFile(t, repoRoot, "docs/specs/"+slug+"/_tasks.md", `---
schema: spec-tasks/v1
spec: wave-collision
qa: declined
qa_reason: no behavioral surface in this fixture
graph:
  nodes:
    - id: task_01
      file: task_01.md
      needs: []
    - id: task_02
      file: task_02.md
      needs: `+needs+`
---
`)
				for _, taskID := range []string{"task_01", "task_02"} {
					writeCitationFixtureFile(t, repoRoot, "docs/specs/"+slug+"/"+taskID+".md", `---
status: pending
type: backend
---

# Task: `+taskID+`

## Verification

- `+"`"+`test -f internal/shared/first.go && test -f internal/shared/second.go`+"`"+`
`)
				}
			}

			result, err := speccheck.CheckStage(specsRoot, repoRoot, slug, speccheck.StageTasks)
			if err != nil {
				t.Fatalf("CheckStage(StageTasks): %v", err)
			}
			findings := findingsWithCode(result, speccheck.CodeWaveCollision)
			if tt.wantFinding {
				if len(findings) != 1 {
					t.Fatalf("%s findings = %#v, want one", speccheck.CodeWaveCollision, findings)
				}
				finding := findings[0]
				if finding.Severity != speccheck.SeverityError {
					t.Errorf("severity = %q, want %q", finding.Severity, speccheck.SeverityError)
				}
				for _, want := range []string{
					"task_01",
					"task_02",
					"internal/shared/first.go",
					"internal/shared/second.go",
					string(spec.TouchFromVerification),
				} {
					if !strings.Contains(finding.Summary, want) {
						t.Errorf("summary = %q, want %q", finding.Summary, want)
					}
				}
				if !strings.Contains(finding.Fix, "needs") {
					t.Errorf("fix = %q, want needs-edge remedy", finding.Fix)
				}
				return
			}
			if len(findings) != 0 {
				t.Fatalf("%s findings = %#v, want none", speccheck.CodeWaveCollision, findings)
			}
			if !tt.writeGraph && !hasSkip(result, speccheck.CodeWaveCollision, "_tasks.md") {
				t.Fatalf("Skipped = %#v, want %s missing _tasks.md", result.Skipped, speccheck.CodeWaveCollision)
			}
		})
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

func TestAuthoredQAVerificationIsRefusedByTaskName(t *testing.T) {
	t.Parallel()

	const slug = "authored-qa-verification"
	repoRoot := t.TempDir()
	specsRoot := filepath.Join(repoRoot, "docs", "specs")
	writeCitationFixtureFile(t, repoRoot, "docs/specs/"+slug+"/_prd.md", `---
spec: authored-qa-verification
status: active
created: 2026-08-31
surfaces: [backend]
---

# Authored QA Verification fixture
`)
	writeCitationFixtureFile(t, repoRoot, "docs/specs/"+slug+"/_tasks.md", `---
schema: spec-tasks/v1
spec: authored-qa-verification
qa: task_01
graph:
  nodes:
    - id: task_01
      file: task_01.md
      needs: []
---
`)
	writeCitationFixtureFile(t, repoRoot, "docs/specs/"+slug+"/task_01.md", `---
task: task_01
spec: authored-qa-verification
status: pending
type: qa
---

# Task 01: QA gate

## Verification

- `+"`grep -q \"^verdict: partial$\" docs/specs/authored-qa-verification/qa/qa-report-old.md`"+`
`)

	result, err := speccheck.CheckStage(specsRoot, repoRoot, slug, speccheck.StageTasks)
	if err != nil {
		t.Fatalf("CheckStage(StageTasks): %v", err)
	}
	findings := findingsWithCode(result, speccheck.CodeQAVerificationAuthored)
	if len(findings) != 1 {
		t.Fatalf("%s findings = %#v, want one", speccheck.CodeQAVerificationAuthored, findings)
	}
	for _, want := range []string{"task_01.md", "qa Task task_01", "authored Verification"} {
		if !strings.Contains(findings[0].Summary, want) {
			t.Errorf("finding summary = %q, want %q", findings[0].Summary, want)
		}
	}
}

func TestDerivedQAVerificationIsAcceptedByTaskStageChecker(t *testing.T) {
	t.Parallel()

	const slug = "derived-qa-verification"
	derived := spec.DerivedQAVerification(slug)
	if len(derived) != 1 {
		t.Fatalf("DerivedQAVerification() = %q, want one command", derived)
	}

	// Exercise the hermeticity classifier without the qa exemption first. This
	// keeps the generated command valid even when it is reused or inspected in
	// isolation.
	if findings := speccheck.NonHermeticVerification(spec.Task{
		File:         slug + "/task_01.md",
		Verification: derived,
	}); len(findings) != 0 {
		t.Fatalf("derived QA Verification hermeticity findings = %#v, want none", findings)
	}

	repoRoot := t.TempDir()
	specsRoot := filepath.Join(repoRoot, "docs", "specs")
	writeCitationFixtureFile(t, repoRoot, "docs/specs/"+slug+"/_prd.md", `---
spec: derived-qa-verification
status: active
created: 2026-08-31
surfaces: [backend]
---

# Derived QA Verification fixture
`)
	writeCitationFixtureFile(t, repoRoot, "docs/specs/"+slug+"/_tasks.md", `---
schema: spec-tasks/v1
spec: derived-qa-verification
qa: task_01
graph:
  nodes:
    - id: task_01
      file: task_01.md
      needs: []
---
`)
	writeCitationFixtureFile(t, repoRoot, "docs/specs/"+slug+"/task_01.md", `---
task: task_01
spec: derived-qa-verification
status: pending
type: qa
---

# Task 01: QA gate

## Verification

- `+"`"+derived[0]+"`"+`
`)

	result, err := speccheck.CheckStage(specsRoot, repoRoot, slug, speccheck.StageTasks)
	if err != nil {
		t.Fatalf("CheckStage(StageTasks): %v", err)
	}
	for _, code := range []string{
		speccheck.CodeQAVerificationAuthored,
		speccheck.CodeVerifyNonHermetic,
	} {
		if findings := findingsWithCode(result, code); len(findings) != 0 {
			t.Errorf("StageTasks %s findings = %#v, want none", code, findings)
		}
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
