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
