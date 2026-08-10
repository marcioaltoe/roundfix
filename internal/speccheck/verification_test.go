// Suite: work-independent Task Verification classification
// Invariant: Tasks whose declared commands are all repository gates or clean-tree checks are classified independently of prose or status.
// Boundary IN: parsed Task Verification commands and the public detector
// Boundary OUT: Markdown parsing, CLI rendering, and Daemon Verification execution
package speccheck_test

import (
	"strings"
	"testing"

	"roundfix/internal/spec"
	"roundfix/internal/speccheck"
)

func TestWorkIndependentVerificationRefusesOnlyWorkIndependentCommands(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		status       spec.Status
		verification []string
		wantFinding  bool
	}{
		{
			name:         "repository gate and clean tree",
			status:       spec.StatusPending,
			verification: []string{"make verify", "git status --porcelain"},
			wantFinding:  true,
		},
		{
			name:         "repository wide Go gates and diff cleanliness",
			status:       spec.StatusInProgress,
			verification: []string{"go build -buildvcs=false ./...", "go test -parallel 16 ./...", "git diff --check"},
			wantFinding:  true,
		},
		{
			name:         "repository gate plus focused effect assertion",
			status:       spec.StatusPending,
			verification: []string{"make verify", "go test ./internal/speccheck -run '^TestTaskEffect$'"},
		},
		{
			name:         "repository wide command with effect selector",
			status:       spec.StatusPending,
			verification: []string{"go test ./... -run '^TestTaskEffect$'"},
		},
		{
			name:         "only effect assertions",
			status:       spec.StatusPending,
			verification: []string{"grep -q 'expected effect' internal/speccheck/verification.go", "go test ./internal/speccheck"},
		},
		{
			name:         "status does not change the declared command property",
			status:       spec.StatusCompleted,
			verification: []string{"make verify", "git status --porcelain"},
			wantFinding:  true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			finding, got := speccheck.WorkIndependentVerification(spec.Task{
				File:         "fixture/task_01.md",
				Status:       tt.status,
				Verification: tt.verification,
			})
			if got != tt.wantFinding {
				t.Fatalf("WorkIndependentVerification() finding = %v, want %v; finding: %#v", got, tt.wantFinding, finding)
			}
			if !got {
				return
			}
			if finding.Code != speccheck.CodeVerifyWorkIndependent || finding.Severity != speccheck.SeverityError {
				t.Errorf("finding identity = %s/%s, want %s/%s", finding.Code, finding.Severity, speccheck.CodeVerifyWorkIndependent, speccheck.SeverityError)
			}
			if len(finding.Where) != 1 || finding.Where[0].Path != "fixture/task_01.md" || finding.Where[0].Line < 1 {
				t.Errorf("finding locations = %#v, want the declaring Task", finding.Where)
			}
			if !strings.Contains(finding.Fix, "Task's own effect") {
				t.Errorf("finding fix = %q, want an effect-asserting repair", finding.Fix)
			}
		})
	}
}

func TestVacuousVerificationCommandIsCaughtBesideHonestSiblings(t *testing.T) {
	t.Parallel()

	for _, testCase := range []struct {
		name    string
		command string
		vacuous bool
	}{
		{
			name:    "an assertion over which paths changed",
			command: "git diff --name-only HEAD | grep -v -E '^(kept)$' | grep -q . && exit 1 || exit 0",
			vacuous: true,
		},
		{
			name:    "a porcelain cleanliness check",
			command: "git status --porcelain | grep -q . && exit 1 || exit 0",
			vacuous: true,
		},
		{
			// Verified on a clean tree on 2026-08-10: git writes nothing, grep
			// matches nothing, and the command exits 1. It fails on an
			// unchanged tree rather than passing, so it is not vacuous.
			name:    "a nonempty-diff assertion fails on an unchanged tree",
			command: "git diff --name-status HEAD | grep -q .",
			vacuous: false,
		},
		{
			name:    "the same paths with a predicate that passes when empty",
			command: "git diff --name-status HEAD | grep -q . && exit 1 || exit 0",
			vacuous: true,
		},
		{
			name:    "a two-snapshot comparison",
			command: `s1="$(git status --porcelain)"; make regen >/dev/null 2>&1; s2="$(git status --porcelain)"; [ "$s1" = "$s2" ]`,
			vacuous: true,
		},
		{
			name:    "a named test that does not exist yet",
			command: "go test ./internal/agent -run '^TestX$' -count=1 -v 2>&1 | tee /dev/stderr | grep -q '^--- PASS: TestX'",
			vacuous: false,
		},
		{
			name:    "a grep over declared source",
			command: "grep -c 'Declared break' internal/agent/selection_test.go",
			vacuous: false,
		},
	} {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			task := spec.Task{Verification: []string{testCase.command}}
			vacuous := speccheck.VacuousVerificationCommands(task)
			if (len(vacuous) != 0) != testCase.vacuous {
				t.Fatalf("VacuousVerificationCommands(%q) = %#v, want vacuous %v", testCase.command, vacuous, testCase.vacuous)
			}
		})
	}
}

// TestOneHonestCommandDoesNotAbsolveAVacuousSibling pins the unit of judgement:
// the Daemon's pre-work probe refuses a Task for any single vacuous command, so
// the static check must not stay silent because a sibling is honest.
func TestOneHonestCommandDoesNotAbsolveAVacuousSibling(t *testing.T) {
	t.Parallel()

	task := spec.Task{Verification: []string{
		"go test ./internal/agent -run '^TestY$' -count=1 -v 2>&1 | grep -q '^--- PASS: TestY'",
		"test -z \"$(git diff --name-only HEAD)\"",
	}}
	vacuous := speccheck.VacuousVerificationCommands(task)
	if len(vacuous) != 1 || !strings.Contains(vacuous[0], "--name-only") {
		t.Fatalf("VacuousVerificationCommands() = %#v, want only the working-tree assertion", vacuous)
	}
	if _, reported := speccheck.WorkIndependentVerification(task); reported {
		t.Fatal("WorkIndependentVerification judges the whole Task and must stay silent here")
	}
}
