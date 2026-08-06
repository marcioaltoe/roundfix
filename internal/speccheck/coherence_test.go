// Suite: Task authoring coherence classification
// Invariant: only declared same-subject contradictions and undeclared rehearsal cases produce findings.
// Boundary IN: parsed Task declarations and the public coherence detectors
// Boundary OUT: Markdown parsing, CLI rendering, and Daemon Verification execution
package speccheck_test

import (
	"path/filepath"
	"strings"
	"testing"

	"roundfix/internal/spec"
	"roundfix/internal/speccheck"
)

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
