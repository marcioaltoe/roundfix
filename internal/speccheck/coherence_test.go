// Suite: Task authoring coherence classification
// Invariant: only declared same-subject contradictions and undeclared rehearsal cases produce findings.
// Boundary IN: parsed Task declarations and the public coherence detectors
// Boundary OUT: Markdown parsing, CLI rendering, and Daemon Verification execution
package speccheck_test

import (
	"strings"
	"testing"

	"roundfix/internal/spec"
	"roundfix/internal/speccheck"
)

func TestContradictoryRequirementsRefusesSameNamedSubject(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		requirements []string
		wantFinding  bool
	}{
		{
			name: "same subject is required and forbidden",
			requirements: []string{
				"MUST keep the named gate enabled.",
				"MUST NOT keep the named gate enabled.",
			},
			wantFinding: true,
		},
		{
			name: "subject cannot be identified",
			requirements: []string{
				"MUST publish the report.",
				"MUST NOT delete the scratch directory.",
			},
		},
		{
			name: "one sided requirement",
			requirements: []string{
				"MUST keep the named gate enabled.",
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
		rehearsalCases []string
		wantFinding    bool
	}{
		{
			name:        "rehearsal has no cases",
			title:       "Rehearse the lifecycle and prove the gates",
			wantFinding: true,
		},
		{
			name:        "stated purpose proves a gate fires",
			title:       "Prove the archive gate fires",
			wantFinding: true,
		},
		{
			name:  "rehearsal declares cases and observations",
			title: "Rehearse the lifecycle and prove the gates",
			rehearsalCases: []string{
				"Case: stale link; Observation: archive reports the link.",
				"Case: prose mention; Observation: archive remains accepted.",
			},
		},
		{
			name:  "case without observation is incomplete",
			title: "Rehearse the lifecycle and prove the gates",
			rehearsalCases: []string{
				"Case: stale link",
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
