---
task: task_02
spec: 0190-a-task-settles-on-the-facts-its-gate-will-check
status: completed
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

Implemented the Task 02 slice. Verification attempts now run in-process checks
after the commands they reach, including failed commands, and before publishing
the attempt verdict. Check failures and checker errors retain diagnostics as
`VerificationCommandError` entries, return through the existing same-session
Verification Feedback repair, and name the Settlement Check in a final failure.
Stop Requests, infrastructure errors and unobserved commands bypass the checks.
Temporary retries rerun them.

The authorization check uses `prepareTaskCommit` with the pre-Agent snapshot
and takes the prospective parent from the Task work directory's `HEAD`. It
reads authorization only when a stageable path other than the Task file is a
Governed Path. The gate and settlement share authorization/delivery-target
resolution and one parent-and-paths per-commit audit. Existing commit findings
retain their wording; prospective findings use the specified wording and
produce no skip records.

Added the defaulted `SettlementChecker` dependency. Its `RefusingFindings`
method and `qaGatePrecondition` share the same `Check`, `PromoteGaps` and
`GatePrecondition` helper. The production baseline and consistency-check caller
remain task_03's work. The existing governed-operation calls stay after the
final attempt, and `runQAGate` keeps its checks.

### Focused-check evidence

- Before implementation, the prospective audit API, SettlementChecker seam,
  in-process attempt checks and the two new test files were absent. The only
  pre-existing dirty path was this Task file's Daemon-owned status change.
- `GOCACHE=/tmp/roundfix-task02-gocache rtk proxy go test -count=1 ./internal/daemon ./internal/speccheck -run 'TestSettlement|TestGatelessGraph|TestDefaultSettlement|TestProspective|TestMechanical|TestAudit'`
  exited 0 after the final code change: daemon and speccheck passed. This
  focused selection also exercises existing mechanical and authorization
  audit regressions without editing `mechanical_test.go`.
- `GOCACHE=/tmp/roundfix-task02-gocache rtk make verify-incremental` exited 0
  after the final code change, with the permission needed by CLI process-tree
  integration tests. Formatting, vet, package tests, skill checks and build
  passed. The initial sandboxed run encountered process-table permission
  failures and repository-boundary guard failures because fixtures were edited
  while it ran; subsequent runs used a stable tree. The initial default Go
  cache access failure was resolved with the task-scoped cache above.
- `rtk proxy git -c core.fsmonitor=false diff --check` exited 0.

Every test below passed in the final focused selection and was included in the
incremental package tests.

| Acceptance criterion | Implementation and focused evidence |
| --- | --- |
| Governed Path outside the grant is refused and named | Shared per-commit bounded-path audit; `TestProspectiveAuditRefusesAGovernedPathOutsideTheGrant` asserts the exact prospective finding sentence. |
| Bounded paths and sanctioned outputs are accepted | The shared audit keeps the grant and regeneration rules; `TestProspectiveAuditAcceptsBoundedAndSanctionedPaths`. |
| Authorization self-approval is refused | The shared self-approval detector sees the prospective changed paths; `TestProspectiveAuditRefusesSelfApproval`. |
| Prospective and created-commit audits agree | Both use the same parent-and-paths function; `TestProspectiveAuditMatchesTheGateAuditOfTheCreatedCommit` compares all finding fields with the specified wording substitution for accepted, escaped and self-approval commits. `TestProspectiveAuditMatchesTheGateWithoutAnAuthorizationReference` preserves the gate's existing behavior when the request has no reference. |
| Finding returns Feedback and repair settles completed | The existing repair loop receives the check's diagnostic; `TestSettlementAuditFindingReturnsFeedbackAndRepairSettles` verifies the repaired tree, unchanged Agent Session, Feedback path, successful settlement and one commit. |
| Final-attempt check failure settles failed with the named reason and no Task commit | `TestSettlementCheckFailureOnFinalAttemptFailsWithNamedReason` verifies the reason prefix, persisted diagnostics, failed Task status, zero committer calls and unchanged Git HEAD. |
| Checker error fails | Error text becomes check diagnostics; `TestSettlementCheckerErrorFailsTheCheck` verifies Feedback, failed settlement and no commit. |
| Check runs after a failed command | The reached command failure stops later commands but proceeds to checks; `TestSettlementChecksRunAfterAFailedCommand` verifies both failures in one Feedback turn. |
| Check publishes started then passed or failed | Existing verification event publication carries the check label; `TestSettlementChecksPublishVerificationEvents` verifies started/failed on attempt 1 and started/command-passed after repair. |
| Gateless graph never calls the checker | Checks are built only for `settlementChecks`; `TestGatelessGraphNeverCallsTheSettlementChecker` verifies no audit or refusing-findings calls. |
| Default checker refuses an outside-grant Governed Path | `TestDefaultSettlementCheckerAuditsAProspectiveCommit` uses a temporary Git repository and committed grant. |
| Default checker returns the gate's refusing findings | `TestDefaultSettlementCheckerReturnsTheGatePreconditionFindings` creates a refusing temporary Spec tree and compares the findings with `qaGatePrecondition`. |

Additional focused tests cover skipping checks on Stop Requests,
infrastructure errors and unobserved commands; avoiding grant reads for
ordinary paths; excluding pre-existing Governed Paths from the prospective
commit; and rerunning checks on a Temporary Verification Failure retry.
Diagnostics use
`runs/<run-id>/verification/batch-NNN-attempt-M-settlement-authorization.log`
under the Artifact Directory.

The declared `## Verification` commands were not run. Task status remains
Daemon-owned. No other Task file or Task Graph manifest was edited, and no
commit, push or Pull Request was created. The Daemon owns the next declared
Verification and settlement.

## Carry-forward provenance

- Source Run: `run_20260930T230750Z_197aa9aff15490a9`
- Source commit: `02d361d8d8dd19eda538c06411d1a6a74d03b18a`
