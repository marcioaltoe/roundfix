// Suite: Baseline release-gate engine journeys
// Invariant: the assembled Baseline engine preserves every safety property that cannot be induced through a production CLI flag.
// Boundary IN: root preservation, unsafe carriers, stale and cross-clone plans, transaction rollback/recovery, reapply, and finding corrections.
// Boundary OUT: real process dispatch, human prompts, sealed ACP supervision, external formatter execution, and repository Verification.

package baseline

import "testing"

func TestBaselineMacroJourneysEngineSafety(t *testing.T) {
	tests := []struct {
		name string
		run  func(*testing.T)
	}{
		{name: "greenfield", run: TestGreenfieldPlanBacksUpWithoutImport},
		{name: "preservation", run: TestPreservationPlanAcceptsCompleteDecisionDocument},
		{name: "profile selection", run: TestProfileAlignmentResolvesExactlyOneProfile},
		{name: "stale plan", run: TestApplyStalePreimage},
		{name: "cross-clone apply", run: TestApplyCrossClone},
		{name: "unsafe carrier", run: TestInstructionAliasUnsafeTargetsBlock},
		{name: "rollback", run: TestTransactionRollback},
		{name: "recovery", run: TestTransactionRecovery},
		{name: "empty reapply", run: TestEmptyReapply},
		{name: "rejected-plan revision", run: TestRejectedPlanRevision},
		{name: "renewed approval", run: TestRevisionRequiresNewApproval},
	}
	for _, test := range tests {
		t.Run(test.name, test.run)
	}
}

func TestBaselineFindingRegressions(t *testing.T) {
	tests := []struct {
		name string
		run  func(*testing.T)
	}{
		{name: "safe alias evidence", run: TestInstructionAliasRetainsOneSourceEvidence},
		{name: "consolidated decision schema", run: TestDecisionDocumentSkeletonPassesStrictParser},
		{name: "HTTP candidate facts", run: TestHTTPRouteCandidatesContainFactsWithoutNormativeClause},
		{name: "PostgreSQL diagnostics", run: TestPostgreSQLEvidenceSeparatesImplementationAndContract},
		{name: "file-level projection", run: TestFileChangesProjectionRejectsMismatch},
		{name: "repository-executable commands", run: TestExecutableVerificationCommandRequiresLocalDeclaration},
		{name: "audit executes nothing", run: TestCapabilityAuditNoExecution},
	}
	for _, test := range tests {
		t.Run(test.name, test.run)
	}
}
