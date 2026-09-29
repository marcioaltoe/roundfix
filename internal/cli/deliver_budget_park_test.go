// Suite: Delivery Run result mapping.
// Invariant: only a supported Run created by this Implement invocation can explain exit 1.
// Boundary IN: the pure Delivery Run decision from process result and Run Database observations.
// Boundary OUT: process spawning, Git, and Run Database reads stay outside this suite.
package cli

import (
	"strings"
	"testing"

	"roundfix/internal/delivery"
	"roundfix/internal/store"
)

func TestDeliveryRunResultMapsABudgetExceededRun(t *testing.T) {
	t.Parallel()
	after := store.Run{ID: "run-budget", State: store.StateBudgetExceeded}

	got, err := deliveryRunResult(
		roundfixCommandResult{exitCode: exitRunFailed, stderr: "  Run Budget exceeded  \n"},
		"",
		nil,
		&after,
	)
	if err != nil {
		t.Fatalf("map budget-ended Run: %v", err)
	}
	if got.Outcome != delivery.RunOutcomeBudgetExceeded || got.RunID != after.ID || got.Reason != "Run Budget exceeded" {
		t.Fatalf("budget-ended result = %+v", got)
	}
}

func TestDeliveryRunResultMapsAnUnresolvedRun(t *testing.T) {
	t.Parallel()
	after := store.Run{ID: "run-unresolved", State: store.StateUnresolved}

	got, err := deliveryRunResult(
		roundfixCommandResult{exitCode: exitRunFailed, stderr: "  unresolved work  "},
		"",
		nil,
		&after,
	)
	if err != nil {
		t.Fatalf("map unresolved Run: %v", err)
	}
	if got.Outcome != delivery.RunOutcomeUnresolved || got.RunID != after.ID || got.Reason != "unresolved work" {
		t.Fatalf("unresolved result = %+v", got)
	}
}

func TestDeliveryRunResultMapsACleanExit(t *testing.T) {
	t.Parallel()
	after := store.Run{ID: "run-clean", State: store.StateClean}

	got, err := deliveryRunResult(
		roundfixCommandResult{exitCode: exitOK},
		"  candidate-head  \n",
		nil,
		&after,
	)
	if err != nil {
		t.Fatalf("map clean exit: %v", err)
	}
	if got.Outcome != delivery.RunOutcomeClean || got.RunID != after.ID || len(got.CandidateCommits) != 1 || got.CandidateCommits[0] != "candidate-head" {
		t.Fatalf("clean result = %+v", got)
	}
}

func TestDeliveryRunResultRefusesAFailedRun(t *testing.T) {
	t.Parallel()
	after := store.Run{ID: "run-failed", State: store.StateFailed}

	_, err := deliveryRunResult(
		roundfixCommandResult{exitCode: exitRunFailed, stderr: "Run failed"},
		"",
		nil,
		&after,
	)
	if err == nil || !strings.Contains(err.Error(), "roundfix implement failed with exit code 1") {
		t.Fatalf("failed Run error = %v, want Implement executor failure", err)
	}
}

func TestDeliveryRunResultRefusesAnExitWithoutANewRun(t *testing.T) {
	t.Parallel()

	_, err := deliveryRunResult(
		roundfixCommandResult{exitCode: exitRunFailed, stderr: "executor stopped without a Run"},
		"",
		nil,
		nil,
	)
	if err == nil || !strings.Contains(err.Error(), "roundfix implement failed with exit code 1") {
		t.Fatalf("runless exit error = %v, want Implement executor failure", err)
	}
}

func TestDeliveryRunResultIgnoresARunThatExistedBeforeTheInvocation(t *testing.T) {
	t.Parallel()
	before := store.Run{ID: "run-existing", State: store.StateUnresolved}
	after := before

	_, err := deliveryRunResult(
		roundfixCommandResult{exitCode: exitRunFailed, stderr: "executor stopped"},
		"",
		&before,
		&after,
	)
	if err == nil || !strings.Contains(err.Error(), "roundfix implement failed with exit code 1") {
		t.Fatalf("stale Run error = %v, want Implement executor failure", err)
	}
}
