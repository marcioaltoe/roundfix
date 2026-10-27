---
task: task_03
spec: 0329-the-skills-lock-matches-what-is-installed-and-a-gate-less-run-settles-green
status: pending
type: backend
complexity: high
---

# Task 03: Every Task settles on the repository Verification, with or without a QA gate

## Overview

The Daemon appends the configured repository Verification to the settlement
attempt of every non-QA Task, whether or not the Task Graph has a QA gate
Task, when `verification.repository_at_settlement` is enabled. The Spec
Consistency and authorization Settlement Checks keep their gate scope. This
answers the adopted Backlog Entry "A Run without a QA gate ends Clean with
the repository red" (Quaestor, 2026-10-26, repeated on v0.125.0 on
2026-10-27) for the Task that breaks the repository, and implements
ADR-0342's first decision.

## Requirements

1. MUST change the condition in `executeTaskWorker`
   (`internal/daemon/task_engine.go`) to `enteredOnRedRepository ||
   plan.RepositoryVerificationAtSettlement`, implementing `_techspec.md` →
   API Contract 4; an empty command still appends nothing, and
   `settlementChecks` keeps `qaTask != nil` for the Spec Consistency
   baseline and the in-process checks (Invariants 3-5).
2. MUST replace `TestUngatedGraphIgnoresEnabledRepositorySettlementVerification`
   in `internal/daemon/settlement_repository_verification_test.go` with
   `TestUngatedGraphRunsRepositoryVerificationAtSettlement`, asserting the
   commands `declared check`, `repository check`. This is the declared break
   of the old contract; every other test in the file keeps its assertions.
3. MUST add, in the same file,
   `TestUngatedGraphRepositoryFailureAtSettlementReturnsFeedback` (a
   repository failure on the first attempt reaches the same Agent Session
   as Verification Feedback naming the command, and a final failure settles
   the Task `failed` with a reason naming it),
   `TestUngatedGraphRepositoryAtSettlementOffRunsOnlyDeclaredCommands` and
   `TestUngatedGraphRunsNoInProcessSettlementCheck` (with a repository
   command configured, a recording `SettlementChecker` is never called and
   the commands are `declared check`, `repository check`).
4. MUST add `TestImplementUngatedGraphRunsRepositoryVerificationAtSettlement`
   to `internal/cli/implement_settlement_checks_test.go`: a graph without a
   QA gate Task run through `roundfix implement` records the Task's declared
   command and then the repository command at settlement. The gated tests of
   that file keep their assertions.
5. MUST keep `TestGatelessGraphNeverCallsTheSettlementChecker` in
   `internal/daemon/settlement_checks_test.go` passing unchanged.
6. MUST migrate the existing implement tests that the change breaks, without
   weakening any assertion. Measured on 2026-10-27 by removing
   `plan.settlementChecks &&` from the condition: 21 tests of
   `./internal/cli` failed, because their graph has no QA gate Task and
   their configuration leaves the built-in repository Verification
   `make verify`, which now runs at settlement in a temporary repository
   without a Makefile. One of them,
   `TestRunImplementWorktreeIsolationExcludesConcurrentUserCommit`, hangs
   instead of failing: its `onTask` commits user work on every Agent turn,
   the repair turn the red `make verify` causes finds nothing to commit, and
   `t.Fatalf` runs off the test goroutine. Each of these tests states its
   repository Verification explicitly in its own configuration: a command
   that passes in its temporary repository, or
   `verification.repository_at_settlement: false` where the test counts
   Verification calls or commits for another purpose, with a one-line
   comment naming the choice. Where a test's count now includes the
   settlement command, the new count is the declared contract. The tests:
   `TestAgentSelectionProfilesMacro`,
   `TestBudgetClockCrossesAfterTheSettledTask`,
   `TestRunImplementBootstrapRunsBeforeAgentWorkAndVerification`,
   `TestRunImplementBootstrapsEachConcurrentTaskWorktreeBeforeAgentWork`,
   `TestRunImplementDetachPrintsReportAndCompletesRun`,
   `TestRunImplementDetachSurvivesCallerProcessGroupKill`,
   `TestRunImplementExecutesSpecEndToEnd`,
   `TestRunImplementOverlapEndsIntegrationPendingAndPrintedCommandWorks`,
   `TestRunImplementPreflightClosesTerminalRunSessionsOnly`,
   `TestRunImplementPreflightPrunesRetainedRunStorage`,
   `TestRunImplementPreflightReapsEmptyTerminalRunAndTaskWorktrees`,
   `TestRunImplementPreflightRetentionPruneFailureIsNonFatal`,
   `TestRunImplementPreflightRetentionZeroSkipsPrune`,
   `TestRunImplementPreflightTerminalReachableChangedBranch`,
   `TestRunImplementPreflightTerminalUniqueChangedBranch`,
   `TestRunImplementRealWorktreeFastForwardsAndCleansPreservingNonOverlappingUserDirt`,
   `TestRunImplementTemporaryVerificationFlowPreservesDeterministicRepair`,
   `TestRunImplementTemporaryVerificationFlowRepeatedTemporaryPreservesTaskWorktree`,
   `TestRunImplementTemporaryVerificationFlowRetriesOnceWithoutAgentRepair`,
   `TestRunImplementVerificationCapacityAndDaemonStatusIntegratedFlow`,
   `TestRunImplementWorktreeIsolationExcludesConcurrentUserCommit`.
   The full `./internal/cli` and `./internal/daemon` suites pass afterwards,
   with `-timeout` on every run.

## Subtasks

- [ ] Characterize: the new gate-less tests fail on the current tree.
- [ ] Change the settlement condition.
- [ ] Replace the old gate-less test and add the new ones.
- [ ] Migrate the 21 implement fixtures that relied on the gate scope.

## Acceptance Criteria

- [ ] A Task of a graph without a QA gate settles on its declared commands and then the repository Verification.
- [ ] A repository failure returns as Verification Feedback and fails the Task on the final attempt.
- [ ] The switch turns it off; the in-process checks never run in such a graph; gated graphs are unchanged.

## Context

- instruction: `docs/adr/0342-every-task-settles-on-the-repository-verification-and-a-run-without-a-qa-gate-checks-its-integrated-tree.md`
- instruction: `docs/adr/0182-a-task-settles-on-the-facts-its-gate-will-check.md`
- instruction: `docs/user-guide/commands/implement.md`
- instruction: `internal/daemon/settlement_checks_test.go`
- instruction: `internal/daemon/settlement_spec_consistency_test.go`
- interface: `internal/daemon/task_engine.go`
- interface: `internal/daemon/settlement_repository_verification_test.go`
- interface: `internal/cli/implement_settlement_checks_test.go`
- interface: `internal/cli/implement_test.go`
- interface: `internal/cli/implement_unix_test.go`

## Consumers

- The settlement command list of `executeTaskWorker` — consumers: `./internal/daemon`, `./internal/cli`

## Verification

- `out="$(go test -count=1 -timeout 20m -v -run '^(TestUngatedGraph.*|TestGatedGraphRunsRepositoryVerificationAfterDeclaredCommands|TestRepositoryVerificationFailure.*|TestTaskDeclaringRepositoryVerificationRunsItOnce|TestRepositoryAtSettlementOffRunsOnlyDeclaredCommands|TestPreWorkProbeNeverRunsTheSettlementRepositoryVerification|TestEmptySettlementRepositoryCommandAppendsNothing|TestGatelessGraphNeverCallsTheSettlementChecker|TestSettlement.*|TestImplementGatedGraphRunsRepositoryVerificationAtSettlement|TestImplementRepositoryAtSettlementOffRunsOnlyDeclaredCommands|TestImplementUngatedGraphRunsRepositoryVerificationAtSettlement|TestAgentSelectionProfilesMacro|TestBudgetClockCrossesAfterTheSettledTask|TestRunImplementBootstrapRunsBeforeAgentWorkAndVerification|TestRunImplementBootstrapsEachConcurrentTaskWorktreeBeforeAgentWork|TestRunImplementDetachPrintsReportAndCompletesRun|TestRunImplementDetachSurvivesCallerProcessGroupKill|TestRunImplementExecutesSpecEndToEnd|TestRunImplementOverlapEndsIntegrationPendingAndPrintedCommandWorks|TestRunImplementPreflightClosesTerminalRunSessionsOnly|TestRunImplementPreflightPrunesRetainedRunStorage|TestRunImplementPreflightReapsEmptyTerminalRunAndTaskWorktrees|TestRunImplementPreflightRetentionPruneFailureIsNonFatal|TestRunImplementPreflightRetentionZeroSkipsPrune|TestRunImplementPreflightTerminalReachableChangedBranch|TestRunImplementPreflightTerminalUniqueChangedBranch|TestRunImplementRealWorktreeFastForwardsAndCleansPreservingNonOverlappingUserDirt|TestRunImplementTemporaryVerificationFlowPreservesDeterministicRepair|TestRunImplementTemporaryVerificationFlowRepeatedTemporaryPreservesTaskWorktree|TestRunImplementTemporaryVerificationFlowRetriesOnceWithoutAgentRepair|TestRunImplementVerificationCapacityAndDaemonStatusIntegratedFlow|TestRunImplementWorktreeIsolationExcludesConcurrentUserCommit)$' ./internal/daemon ./internal/cli 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestUngatedGraphRunsRepositoryVerificationAtSettlement TestUngatedGraphRepositoryFailureAtSettlementReturnsFeedback TestUngatedGraphRepositoryAtSettlementOffRunsOnlyDeclaredCommands TestUngatedGraphRunsNoInProcessSettlementCheck TestImplementUngatedGraphRunsRepositoryVerificationAtSettlement TestGatedGraphRunsRepositoryVerificationAfterDeclaredCommands TestGatelessGraphNeverCallsTheSettlementChecker TestImplementGatedGraphRunsRepositoryVerificationAtSettlement TestAgentSelectionProfilesMacro TestBudgetClockCrossesAfterTheSettledTask TestRunImplementBootstrapRunsBeforeAgentWorkAndVerification TestRunImplementBootstrapsEachConcurrentTaskWorktreeBeforeAgentWork TestRunImplementDetachPrintsReportAndCompletesRun TestRunImplementDetachSurvivesCallerProcessGroupKill TestRunImplementExecutesSpecEndToEnd TestRunImplementOverlapEndsIntegrationPendingAndPrintedCommandWorks TestRunImplementPreflightClosesTerminalRunSessionsOnly TestRunImplementPreflightPrunesRetainedRunStorage TestRunImplementPreflightReapsEmptyTerminalRunAndTaskWorktrees TestRunImplementPreflightRetentionPruneFailureIsNonFatal TestRunImplementPreflightRetentionZeroSkipsPrune TestRunImplementPreflightTerminalReachableChangedBranch TestRunImplementPreflightTerminalUniqueChangedBranch TestRunImplementRealWorktreeFastForwardsAndCleansPreservingNonOverlappingUserDirt TestRunImplementTemporaryVerificationFlowPreservesDeterministicRepair TestRunImplementTemporaryVerificationFlowRepeatedTemporaryPreservesTaskWorktree TestRunImplementTemporaryVerificationFlowRetriesOnceWithoutAgentRepair TestRunImplementVerificationCapacityAndDaemonStatusIntegratedFlow TestRunImplementWorktreeIsolationExcludesConcurrentUserCommit; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done; if grep -q 'func TestUngatedGraphIgnoresEnabledRepositorySettlementVerification' internal/daemon/settlement_repository_verification_test.go; then echo 'the old gate-less test still pins the gate scope' >&2; exit 1; fi` — expected: exit 0. Before this Task the five new tests do not exist and the old test still exists, so the pass check fails. Afterwards the new tests, the existing gated settlement tests, the in-process Settlement Check tests and the 21 migrated implement tests pass, and the old test is gone.

## References

- `_prd.md` → Goals; Core Feature 2; User Experience; Success Metric 4; Success Metric 6
- `_techspec.md` → Current behavior 4; Current behavior 6; API Contract 4; Invariants 3-5; Testing Approach; Build Order 3
- ADR-0342; ADR-0182
