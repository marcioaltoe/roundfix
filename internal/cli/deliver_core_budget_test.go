// Boundary: disposable Git repositories, harmless shell commands and Docker stand-ins.
package cli

import (
	"path/filepath"
	"slices"
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

func TestDeliverLanePassesItsVerificationCoreBudget(t *testing.T) {
	t.Parallel()
	for _, row := range []struct {
		name, gomaxprocs, want string
	}{
		{"half divided by two", "", "3"},
		{"explicit GOMAXPROCS", "5", "5"},
	} {
		t.Run(row.name, func(t *testing.T) {
			t.Parallel()
			f := newDeliverLaneFixture(t)
			writeUserConfig(t, f.workflow.loaded.HomeDir, "runs:\n  max_verifying: 2\n")
			withVerificationCoreBudgetInputs(t, row.gomaxprocs)
			deps := commandEnvironmentForTest(t).dependencies
			// Keep the fixture's tool environment and the hermetic budget inputs.
			f.environment = append(f.environment, "GOMAXPROCS="+row.gomaxprocs)
			deps.environ = func() []string { return slices.Clone(f.environment) }
			ctx := contextWithCommandDependencies(t.Context(), deps)
			stdout, stderr, code := invokeDeliverLaneContext(t, ctx, f, "--all")
			if code != exitOK || stderr != "" {
				t.Fatalf("lane exit = %d; stdout = %s; stderr = %s", code, stdout, stderr)
			}
			starts := f.starts(t)
			if len(starts) != 4 {
				t.Fatalf("starts = %+v", starts)
			}
			args := starts[3].Args
			want := []string{"-e", "GOSUMDB=off", "-e", "GOMAXPROCS=" + row.want, "-e", "ROUNDFIX_VERIFICATION_CORES=" + row.want, "-v"}
			at := slices.Index(args, "GOSUMDB=off") - 1
			if at < 0 || at+len(want) > len(args) || !slices.Equal(args[at:at+len(want)], want) {
				t.Fatalf("container environment arguments = %q, want sequence %q", args, want)
			}
		})
	}
}
