// Boundary: disposable Git repositories, harmless shell commands.
package cli

import (
	"path/filepath"
	"testing"
)

func TestTheRepositoryGateCommandReceivesTheCoreBudget(t *testing.T) {
	t.Parallel()
	for _, row := range []struct {
		name, gomaxprocs, want string
	}{
		{"half divided by two", "", "3|3"},
		{"explicit GOMAXPROCS", "5", "5|5"},
	} {
		t.Run(row.name, func(t *testing.T) {
			t.Parallel()
			f := newDeliveryItemFixture(t)
			withVerificationCoreBudgetInputs(t, row.gomaxprocs)
			f.workflow.loaded.Config.Runs.MaxVerifying = 2
			f.workflow.loaded.Config.Defaults.Verification = `printf '%s|%s' "$GOMAXPROCS" "$ROUNDFIX_VERIFICATION_CORES" > core-budget.txt`
			result, err := f.workflow.Gate(commandContextForTest(t, t.Context()), f.repo, "widget", f.head)
			if err != nil || !result.Passed {
				t.Fatalf("Gate = %+v, err = %v", result, err)
			}
			if got := mustRead(t, filepath.Join(f.repo, "core-budget.txt")); got != row.want {
				t.Fatalf("gate environment = %q, want %q", got, row.want)
			}
		})
	}
}
