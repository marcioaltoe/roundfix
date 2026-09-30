---
task: task_02
spec: 0190-a-task-settles-on-the-facts-its-gate-will-check
status: pending
type: backend
complexity: high
---

# Task 02: A Verification attempt audits the commit the Daemon is about to create

## Overview

The QA gate audits each Task commit's Governed Paths against the Spec's grant
(`detectMechanicalAuthPaths` in `internal/speccheck/mechanical.go`), hours
after the Task that made the commit has settled. This Task gives a Verification
attempt in-process checks, and adds the first one: the gate's own audit,
applied to the commit the Daemon is about to create. A finding returns to the
same Agent Session as Verification Feedback. The gate's audit is unchanged and
shares one per-commit function with the new check, so the two cannot disagree.

## Requirements

1. MUST make `runVerificationAttempt` run `req.Checks` after the commands it
   reached, whatever their result, and before it publishes the attempt's
   verdict. It MUST skip them only when the attempt ends without a verdict: a
   Stop Request, an infrastructure error or an unobserved command.
2. MUST record a check that returns a failure, or an error, in
   `CommandFailures` as a `VerificationCommandError` whose `Command` is the
   check's label and whose `OutputPath` is its diagnostics file. The attempt's
   verdict MUST be `failed`. An error MUST never pass. Diagnostics go to the
   path the TechSpec names under "Checks inside a Verification attempt".
3. MUST publish API Contract 2's `verification` events for each check, with
   the label as `command`.
4. MUST add the `SettlementChecker` Engine dependency with the two methods
   the TechSpec shows, defaulted the way `MechanicalStage` is. The default's
   `RefusingFindings` MUST run `speccheck.Check`, `PromoteGaps` and
   `GatePrecondition`, the calls `qaGatePrecondition` makes, through one
   helper both share. This Task adds no caller of `RefusingFindings`; task_03
   does. `verifyTask` MUST fill `Checks` only when the plan's
   `settlementChecks` is true.
5. MUST make `detectMechanicalAuthPaths` evaluate each Task commit through one
   per-commit function that takes the commit's parent revision and its changed
   paths. `speccheck.AuditProspectiveTaskCommit` MUST call that same function
   with `ProspectiveTaskCommit.Parent` and `.Changed`, and add no skip record.
6. MUST keep every existing finding's text for an existing commit. A finding
   for a prospective commit uses the wording the TechSpec gives under "The
   prospective Task commit audit".
7. MUST resolve the authorization record and the delivery target through one
   helper shared by `qaMechanicalRequest` and the settlement check.
8. MUST take the prospective commit's paths from `prepareTaskCommit`'s
   stageable list and its parent from `HEAD` of the Task's work directory. The
   audit runs only when a staged path other than the Task file is a Governed
   Path; otherwise the check passes without reading a grant.
9. MUST settle a Task whose in-process check fails on the final attempt as
   `failed`, with API Contract 3's reason, and create no Task commit.
10. MUST NOT edit `internal/speccheck/mechanical_test.go`, rename or remove a
    top-level test, move the `spec.RequireGovernedOperation` call, or change
    what `runQAGate` checks.
11. MUST put the new tests in the two test files this Task creates. The
    `speccheck` tests run over temporary Git repositories. The engine tests use
    a fake `SettlementChecker` whose `RefusingFindings` returns no finding
    unless a test sets one, and two tests run the default checker over a
    temporary Git repository and a temporary Spec tree.

## Subtasks

- [ ] Run in-process checks inside a Verification attempt and publish them.
- [ ] Share one per-commit audit between the gate and a prospective commit.
- [ ] Add the authorization check and the failure reason.
- [ ] Add one test for each acceptance criterion, each negative case separate.

## Acceptance Criteria

- [ ] A prospective commit that changes a Governed Path outside the grant is
      refused, and the finding names the path.
- [ ] A prospective commit that changes only bounded paths and sanctioned
      regeneration outputs is accepted.
- [ ] A prospective commit that changes the authorization record is refused as
      self-approval.
- [ ] For the same parent and the same changed paths, the prospective audit
      and the gate's audit of the created commit report the same findings.
- [ ] An audit finding returns Verification Feedback, and a repaired tree
      settles `completed`.
- [ ] A check that fails on the final attempt settles the Task `failed` with
      `Settlement check failed: settlement check: authorization:` in its
      reason, and no Task commit exists.
- [ ] A checker error fails the check.
- [ ] A check runs after a failed command, so one Feedback turn carries both.
- [ ] Each check publishes `started` and then `command-passed` or `failed`.
- [ ] A graph without a QA gate Task never calls the checker.
- [ ] The default checker refuses a Governed Path outside the grant in a
      temporary Git repository.
- [ ] The default checker, over a temporary Spec tree with a refusing finding,
      returns the findings the gate's precondition reports for that tree.

## Context

- interface: `internal/daemon/engine.go`
- interface: `internal/daemon/task_engine.go`
- interface: `internal/speccheck/mechanical.go`
- creates: `internal/daemon/settlement_checks.go`
- creates: `internal/daemon/settlement_checks_test.go`
- creates: `internal/speccheck/prospective_commit_audit_test.go`
- instruction: `docs/adr/0182-a-task-settles-on-the-facts-its-gate-will-check.md`
- instruction: `docs/adr/0130-the-audit-judges-governed-paths-and-history-keeps-the-set-honest.md`
- instruction: `docs/adr/0178-a-task-commit-is-authorized-by-the-grant-it-ran-under.md`

## Verification

- `out="$(go test -count=1 -v -run "^(TestProspectiveAuditRefusesAGovernedPathOutsideTheGrant|TestProspectiveAuditAcceptsBoundedAndSanctionedPaths|TestProspectiveAuditRefusesSelfApproval|TestProspectiveAuditMatchesTheGateAuditOfTheCreatedCommit)$" ./internal/speccheck 2>&1)" || { printf "%s\\n" "$out"; exit 1; }; for name in TestProspectiveAuditRefusesAGovernedPathOutsideTheGrant TestProspectiveAuditAcceptsBoundedAndSanctionedPaths TestProspectiveAuditRefusesSelfApproval TestProspectiveAuditMatchesTheGateAuditOfTheCreatedCommit; do printf "%s\\n" "$out" | grep -q -- "--- PASS: $name" || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done` — expected: exit 0; before this Task none of the four named tests exists, so the command fails.
- `out="$(go test -count=1 -v -run "^(TestSettlementAuditFindingReturnsFeedbackAndRepairSettles|TestSettlementCheckFailureOnFinalAttemptFailsWithNamedReason|TestSettlementCheckerErrorFailsTheCheck|TestSettlementChecksRunAfterAFailedCommand|TestSettlementChecksPublishVerificationEvents|TestGatelessGraphNeverCallsTheSettlementChecker|TestDefaultSettlementCheckerAuditsAProspectiveCommit|TestDefaultSettlementCheckerReturnsTheGatePreconditionFindings)$" ./internal/daemon 2>&1)" || { printf "%s\\n" "$out"; exit 1; }; for name in TestSettlementAuditFindingReturnsFeedbackAndRepairSettles TestSettlementCheckFailureOnFinalAttemptFailsWithNamedReason TestSettlementCheckerErrorFailsTheCheck TestSettlementChecksRunAfterAFailedCommand TestSettlementChecksPublishVerificationEvents TestGatelessGraphNeverCallsTheSettlementChecker TestDefaultSettlementCheckerAuditsAProspectiveCommit TestDefaultSettlementCheckerReturnsTheGatePreconditionFindings; do printf "%s\\n" "$out" | grep -q -- "--- PASS: $name" || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done` — expected: exit 0; before this Task none of the eight named tests exists, so the command fails.

## References

- [_prd.md](_prd.md) — Goal 1; Goal 2; User Story 2; Core Feature 4; Core Feature 5; Core Feature 6; Success Metric 3
- [_techspec.md](_techspec.md) — Interfaces; Checks inside a Verification attempt; The prospective Task commit audit; API Contract 2; API Contract 3; Testing Approach 2; Testing Approach 3; Build Order 2
- ADR-0182; ADR-0038; ADR-0096; ADR-0111; ADR-0117; ADR-0130; ADR-0135; ADR-0166; ADR-0178; ADR-0179

## Result
