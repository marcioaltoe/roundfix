// Suite: Baseline release-gate ACP journeys
// Invariant: classification always accepts one sealed validated proposal or returns a complete deterministic manual proposal.
// Boundary IN: preferred selection, byte-identical fallback, invalid semantic output, and manual classification.
// Boundary OUT: CLI prompting, repository mutation, final plan approval, and model quality.

package baselineacp

import "testing"

func TestBaselineMacroJourneysACPProposals(t *testing.T) {
	tests := []struct {
		name string
		run  func(*testing.T)
	}{
		{name: "preferred proposal", run: TestPreferredFallbackValidPreferredPreventsFallback},
		{name: "fallback receives identical snapshot", run: TestPreferredFallbackInvalidPreferredUsesIdenticalSnapshotBytes},
		{name: "manual classification fallback", run: TestManualClassificationFallbackWhenSelectionsUnavailableOrInvalid},
		{name: "tool timeout and oversized output fail closed", run: TestPreferredFallbackRejectsToolTimeoutAndOversizedOutput},
	}
	for _, test := range tests {
		t.Run(test.name, test.run)
	}
}
