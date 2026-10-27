---
task: task_03
spec: 0329-the-skills-lock-matches-what-is-installed-and-a-gate-less-run-settles-green
status: completed
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
7. MUST re-record `docs/references/coverage-record.json` after replacing
   `TestUngatedGraphIgnoresEnabledRepositorySettlementVerification`, only with
   `go test ./internal/spec -run '^TestCoverageEquivalence$' -update-coverage-record -count=1`
   (`docs/agents/specific-repository.md`); `TestCoverageEquivalence` runs in
   the repository Verification at settlement and fails while the record still
   lists the replaced test.

## Subtasks

- [ ] Characterize: the new gate-less tests fail on the current tree.
- [ ] Change the settlement condition.
- [ ] Replace the old gate-less test and add the new ones.
- [ ] Migrate the 21 implement fixtures that relied on the gate scope.
- [ ] Re-record the coverage record with its command.

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
- interface: `docs/references/coverage-record.json`

## Consumers

- The settlement command list of `executeTaskWorker` — consumers: `./internal/daemon`, `./internal/cli`

## Verification

- `out="$(go test -count=1 -timeout 20m -v -run '^(TestUngatedGraph.*|TestGatedGraphRunsRepositoryVerificationAfterDeclaredCommands|TestRepositoryVerificationFailure.*|TestTaskDeclaringRepositoryVerificationRunsItOnce|TestRepositoryAtSettlementOffRunsOnlyDeclaredCommands|TestPreWorkProbeNeverRunsTheSettlementRepositoryVerification|TestEmptySettlementRepositoryCommandAppendsNothing|TestGatelessGraphNeverCallsTheSettlementChecker|TestSettlement.*|TestImplementGatedGraphRunsRepositoryVerificationAtSettlement|TestImplementRepositoryAtSettlementOffRunsOnlyDeclaredCommands|TestImplementUngatedGraphRunsRepositoryVerificationAtSettlement|TestAgentSelectionProfilesMacro|TestBudgetClockCrossesAfterTheSettledTask|TestRunImplementBootstrapRunsBeforeAgentWorkAndVerification|TestRunImplementBootstrapsEachConcurrentTaskWorktreeBeforeAgentWork|TestRunImplementDetachPrintsReportAndCompletesRun|TestRunImplementDetachSurvivesCallerProcessGroupKill|TestRunImplementExecutesSpecEndToEnd|TestRunImplementOverlapEndsIntegrationPendingAndPrintedCommandWorks|TestRunImplementPreflightClosesTerminalRunSessionsOnly|TestRunImplementPreflightPrunesRetainedRunStorage|TestRunImplementPreflightReapsEmptyTerminalRunAndTaskWorktrees|TestRunImplementPreflightRetentionPruneFailureIsNonFatal|TestRunImplementPreflightRetentionZeroSkipsPrune|TestRunImplementPreflightTerminalReachableChangedBranch|TestRunImplementPreflightTerminalUniqueChangedBranch|TestRunImplementRealWorktreeFastForwardsAndCleansPreservingNonOverlappingUserDirt|TestRunImplementTemporaryVerificationFlowPreservesDeterministicRepair|TestRunImplementTemporaryVerificationFlowRepeatedTemporaryPreservesTaskWorktree|TestRunImplementTemporaryVerificationFlowRetriesOnceWithoutAgentRepair|TestRunImplementVerificationCapacityAndDaemonStatusIntegratedFlow|TestRunImplementWorktreeIsolationExcludesConcurrentUserCommit)$' ./internal/daemon ./internal/cli 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestUngatedGraphRunsRepositoryVerificationAtSettlement TestUngatedGraphRepositoryFailureAtSettlementReturnsFeedback TestUngatedGraphRepositoryAtSettlementOffRunsOnlyDeclaredCommands TestUngatedGraphRunsNoInProcessSettlementCheck TestImplementUngatedGraphRunsRepositoryVerificationAtSettlement TestGatedGraphRunsRepositoryVerificationAfterDeclaredCommands TestGatelessGraphNeverCallsTheSettlementChecker TestImplementGatedGraphRunsRepositoryVerificationAtSettlement TestAgentSelectionProfilesMacro TestBudgetClockCrossesAfterTheSettledTask TestRunImplementBootstrapRunsBeforeAgentWorkAndVerification TestRunImplementBootstrapsEachConcurrentTaskWorktreeBeforeAgentWork TestRunImplementDetachPrintsReportAndCompletesRun TestRunImplementDetachSurvivesCallerProcessGroupKill TestRunImplementExecutesSpecEndToEnd TestRunImplementOverlapEndsIntegrationPendingAndPrintedCommandWorks TestRunImplementPreflightClosesTerminalRunSessionsOnly TestRunImplementPreflightPrunesRetainedRunStorage TestRunImplementPreflightReapsEmptyTerminalRunAndTaskWorktrees TestRunImplementPreflightRetentionPruneFailureIsNonFatal TestRunImplementPreflightRetentionZeroSkipsPrune TestRunImplementPreflightTerminalReachableChangedBranch TestRunImplementPreflightTerminalUniqueChangedBranch TestRunImplementRealWorktreeFastForwardsAndCleansPreservingNonOverlappingUserDirt TestRunImplementTemporaryVerificationFlowPreservesDeterministicRepair TestRunImplementTemporaryVerificationFlowRepeatedTemporaryPreservesTaskWorktree TestRunImplementTemporaryVerificationFlowRetriesOnceWithoutAgentRepair TestRunImplementVerificationCapacityAndDaemonStatusIntegratedFlow TestRunImplementWorktreeIsolationExcludesConcurrentUserCommit; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done; if grep -q 'func TestUngatedGraphIgnoresEnabledRepositorySettlementVerification' internal/daemon/settlement_repository_verification_test.go; then echo 'the old gate-less test still pins the gate scope' >&2; exit 1; fi` — expected: exit 0. Before this Task the five new tests do not exist and the old test still exists, so the pass check fails. Afterwards the new tests, the existing gated settlement tests, the in-process Settlement Check tests and the 21 migrated implement tests pass, and the old test is gone.

## References

- `_prd.md` → Goals; Core Feature 2; User Experience; Success Metric 4; Success Metric 6
- `_techspec.md` → Current behavior 4; Current behavior 6; API Contract 4; Invariants 3-5; Testing Approach; Build Order 3
- ADR-0342; ADR-0182

## Result

Implemented the task_03 slice. `executeTaskWorker` now selects the repository
Verification with `enteredOnRedRepository ||
plan.RepositoryVerificationAtSettlement`; the existing empty-command handling,
Spec Consistency baseline scope and in-process Settlement Check scope remain
unchanged. Replaced the old ungated-contract test, added the four other named
regressions, and configured all 21 listed CLI fixtures explicitly without
changing their existing assertions. The integrated-tree check belongs to
task_04 and is not included here.

Acceptance evidence:

- Declared commands precede repository Verification in an ungated graph:
  `TestUngatedGraphRunsRepositoryVerificationAtSettlement` and
  `TestImplementUngatedGraphRunsRepositoryVerificationAtSettlement` pass.
- Repository failure returns command-naming Verification Feedback in the same
  Agent Session, and persistent failure settles the fixture Task failed with
  a command-naming reason:
  `TestUngatedGraphRepositoryFailureAtSettlementReturnsFeedback/repair` and
  `/final-failure` pass, including status and commit-count assertions.
- Switch-off runs only declared commands:
  `TestUngatedGraphRepositoryAtSettlementOffRunsOnlyDeclaredCommands` passes.
  `TestUngatedGraphRunsNoInProcessSettlementCheck` records both consistency
  and audit calls as zero while observing both commands. The unchanged
  `TestGatelessGraphNeverCallsTheSettlementChecker`, gated repository tests,
  Spec Consistency tests and authorization Settlement Check tests pass.

Checks run in this turn (all Go test runs bounded by an explicit timeout):

- Before the production edit,
  `go test -count=1 -timeout 3m ./internal/daemon ./internal/cli -run '^Test(Implement)?UngatedGraph'`
  exited 1: the new order/checker tests observed only `declared check`,
  feedback observed one Agent request instead of two, persistent repository
  failure incorrectly completed the fixture Task, and the CLI omitted the
  repository command. After the edit the same focused check exited 0.
- `go test -count=1 -timeout 4m -v ./internal/daemon -run '^(TestUngatedGraph.*|TestGatedGraph.*|TestRepositoryVerification.*|TestTaskDeclaringRepositoryVerification.*|TestRepositoryAtSettlementOff.*|TestPreWorkProbeNeverRuns.*|TestEmptySettlementRepository.*|TestGatelessGraphNeverCalls.*|TestSettlement.*)$'`
  exited 0.
- A focused CLI selection of all 21 named migrated tests plus gated/ungated
  settlement and switch-off tests, with `-count=1 -timeout 8m -v`, exited 0.
  Every named migrated test reported PASS. An inspection comparing their
  pre-edit and current condition/assertion lines found them unchanged.
- `go test -count=1 -timeout 20m ./internal/daemon ./internal/cli`
  exited 0: daemon 50.959s, CLI 197.405s. These are the two package suites
  expressly required by Requirement 6, not the declared Verification wrapper.
- `GOOS=linux go build ./internal/cli ./internal/daemon` exited 0.
- `go test -count=1 -timeout 3m ./internal/spec -run '^TestCoverageEquivalence$'`
  exited 0 after regeneration; the replaced test is absent from the record
  and all new test names are present.
- `git diff --check` exited 0.

Coverage regeneration and follow-up:

- Ran only the sanctioned writer:
  `GOFLAGS='-timeout=3m' go test ./internal/spec -run '^TestCoverageEquivalence$' -update-coverage-record -count=1`.
  Its test body passed and wrote `docs/references/coverage-record.json`, but
  the command exited 1 because suiteguard reported
  `modified: docs/references/coverage-record.json`. `internal/spec/main_test.go`
  invokes `suiteguard.Main` without declaring the sanctioned regeneration
  command; suiteguard requires that declaration as well as the approved Spec
  authorization. No guard bypass or hand edit was used. The generated output
  remains in this diff and passes the read-only equivalence check above.
  The recorder also incorporates previously unrecorded tests elsewhere in
  the current repository, as its deterministic full inventory requires.
- Follow-up outside this Task's interfaces: wire the coverage recorder's
  sanctioned-command declaration into the spec package test lifecycle, so
  its first writing invocation can exit successfully under suiteguard.

The incoming `status: in_progress` is preserved. No Task Graph, other Task
file, commit, push or pull request was changed or created. The command under
`## Verification`, selected repository Verification and incremental
Verification were not run; Daemon settlement remains pending.

## Carry-forward provenance

- Source Run: `run_20261027T193605Z_00ac54bea6e5fdc9`
- Source commit: `14013d6ca42b0cc31b0628c1fce4af29dc384751`
