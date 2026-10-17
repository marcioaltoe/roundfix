// Suite: Baseline release-gate engine journeys
// Invariant: the assembled Baseline engine preserves every safety property that cannot be induced through a production CLI flag.
// Boundary IN: root preservation, unsafe carriers, stale and cross-clone plans, transaction rollback/recovery, reapply, and finding corrections.
// Boundary OUT: real process dispatch, human prompts, sealed ACP supervision, external formatter execution, and repository Verification.

package baseline

import "testing"

func TestProjectDecisionJourneyEngine(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		run  func(*testing.T)
	}{
		{name: "strict decision objects", run: TestProjectDecisionValidation},
		{name: "missing decision produces no partial Plan", run: TestPlanDocumentMissingDecisionsReturnsResultWithoutPartialPlan},
		{name: "identifier strategy is required", run: TestIdentifierStrategyDecision},
		{name: "Better Auth is capability-bound", run: TestAuthProviderDecision},
		{name: "derived HTTP conflict stops planning", run: TestHTTPContractConflict},
		{name: "equivalent normalized decisions", run: TestProfileAlignmentEquivalentNormalizedDecisions},
		{name: "render apply and empty reapply", run: TestProjectDecisionRendering},
	}
	for _, test := range tests {
		t.Run(test.name, test.run)
	}
}

func TestToolingAuthorizationJourneyCoreClause(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		run  func(*testing.T)
	}{
		{name: "renders for every maintained Profile", run: TestToolingAuthorityClause},
		{name: "cannot be disabled", run: TestToolingAuthorityCannotBeDisabled},
		{name: "remains source-accounted", run: TestToolingAuthorityAccounting},
	}
	for _, test := range tests {
		t.Run(test.name, test.run)
	}
}

func TestBaselineFindingRegressions(t *testing.T) {
	t.Parallel()

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
