---
task: task_04
spec: 0329-the-skills-lock-matches-what-is-installed-and-a-gate-less-run-settles-green
status: completed
type: backend
complexity: high
---

# Task 04: A Run without a QA gate checks its integrated tree before it ends Clean

## Overview

Settlement Checks see one Task's tree, so two Tasks of one Wave can each
pass and integrate into a red tree. In a graph with a QA gate Task the
gate's precondition finds it; in a graph without one, nothing does. After
this Task, such a Run whose Task Capacity is above one and that completed
two or more Tasks runs the repository Verification once on its integrated
Run Worktree, the Integrated Repository Verification, and ends `Unresolved`
with a named reason when it fails. This completes the answer to the adopted
Backlog Entry "A Run without a QA gate ends Clean with the repository red"
(Quaestor, 2026-10-26 and 2026-10-27) and implements ADR-0342's second
decision.

## Requirements

1. MUST implement `_techspec.md` → API Contract 5 in `TaskCycle`
   (`internal/daemon/task_engine.go`): after `runTaskScheduler` returns
   without error, run the repository Verification once in `plan.WorkDir`
   when the plan has no QA gate Task, the switch is on, the trimmed command
   is not empty, `plan.Concurrency > 1`, `result.Completed >= 2`, and no
   Task failed or was skipped. Use the Verifier under the Verification
   Capacity gate, the diagnostics file
   `<artifact-root>/runs/<run-id>/verification/integrated-repository.log`,
   the single exclusive retry after exit `75`, and a `daemon_verification`
   event with the summary `Integrated Repository Verification passed:
   <command>.` or `Integrated Repository Verification failed: <command>.`.
   A Stop Request or an exhausted Run Budget during the check ends the
   cycle as either does today.
2. MUST add `IntegratedVerificationFailure` to `TaskCycleResult`, set to
   `Integrated Repository Verification failed: <command>; diagnostics:
   <path>` on a failure. The check never changes a Task status and never
   sends an Agent turn (Invariant 6).
3. MUST implement API Contract 6 in `runImplementCommand`
   (`internal/cli/implement.go`): a non-empty failure makes the outcome
   `Unresolved` before the Run Branch is integrated with the user's branch,
   prints the reason on stderr before `Implement Run <id> reached
   Unresolved.`, carries it as the terminal completion reason, keeps the Run
   Worktree and exits `1`.
4. MUST add `internal/daemon/integrated_repository_verification_test.go`
   with `TestIntegratedRepositoryVerificationRunsOnceAfterAParallelUngatedGraph`
   (Task Capacity 2, two completed Tasks: exactly one extra repository
   command, run in the Run Worktree after both Tasks settle, its event and
   diagnostics path), `TestIntegratedRepositoryVerificationFailureLeavesTaskStatuses`
   (failure sets the result field, both Tasks stay `completed`, no Agent
   request after the Tasks), and `TestIntegratedRepositoryVerificationDoesNotRun`
   with subtests `capacity-one`, `one-completed`, `failed-task`,
   `switch-off`, `empty-command` and `gated-graph`. A recording verifier
   tells the integrated call from settlement calls by its diagnostics path.
5. MUST add `TestImplementUngatedRunEndsUnresolvedWhenTheIntegratedTreeIsRed`
   and `TestImplementUngatedRunEndsCleanWhenTheIntegratedTreeIsGreen` to
   `internal/cli/implement_settlement_checks_test.go`: two independent Tasks
   without a QA gate, `worktree.concurrency: 2` and a repository command;
   the fake verifier fails only the integrated call in the first test. The
   red test asserts exit `1`, the stderr line, the Run record `Unresolved`
   and the kept Run Worktree; the green test asserts exit `0` and Clean.
6. MUST keep the 21 implement tests task_03 migrated passing. A gate-less
   fixture with two or more completed Tasks at the default Task Capacity 2
   now also runs the Integrated Repository Verification; where a test
   counts Verification calls, its count includes that call as the declared
   contract, and no other assertion changes. The full `./internal/cli` and
   `./internal/daemon` suites pass afterwards, with `-timeout` on every run.

## Subtasks

- [ ] Characterize: the new tests fail on the current tree.
- [ ] Add the check to `TaskCycle` and the result field.
- [ ] Make the implement outcome read it.
- [ ] Add the daemon and implement tests, and keep the migrated implement tests passing.

## Acceptance Criteria

- [ ] A parallel gate-less Run with two completed Tasks runs the repository Verification once on the integrated tree.
- [ ] A failure ends the Run `Unresolved` with the command and its diagnostics file named, and changes no Task status.
- [ ] The check never runs at Task Capacity one, with one completed Task, after a failed or skipped Task, with the switch off, without a command, or in a gated graph.

## Context

- instruction: `docs/adr/0342-every-task-settles-on-the-repository-verification-and-a-run-without-a-qa-gate-checks-its-integrated-tree.md`
- instruction: `docs/user-guide/commands/implement.md`
- instruction: `.agents/skills/roundfix/references/implement.md`
- instruction: `internal/daemon/settlement_repository_verification_test.go`
- interface: `internal/daemon/task_engine.go`
- creates: `internal/daemon/integrated_repository_verification_test.go`
- interface: `internal/cli/implement.go`
- interface: `internal/cli/implement_settlement_checks_test.go`
- interface: `internal/cli/implement_test.go`
- interface: `internal/cli/implement_unix_test.go`

## Consumers

- `TaskCycleResult.IntegratedVerificationFailure` — consumers: `./internal/cli`
- The Task cycle's terminal step — consumers: `./internal/daemon`

## Verification

- `out="$(go test -count=1 -timeout 20m -v -run '^(TestIntegratedRepositoryVerification.*|TestUngatedGraph.*|TestGatedGraphRunsRepositoryVerificationAfterDeclaredCommands|TestImplementUngatedRun.*|TestImplementUngatedGraphRunsRepositoryVerificationAtSettlement|TestImplementGatedGraphRunsRepositoryVerificationAtSettlement|TestAgentSelectionProfilesMacro|TestBudgetClockCrossesAfterTheSettledTask|TestRunImplementBootstrapRunsBeforeAgentWorkAndVerification|TestRunImplementBootstrapsEachConcurrentTaskWorktreeBeforeAgentWork|TestRunImplementDetachPrintsReportAndCompletesRun|TestRunImplementDetachSurvivesCallerProcessGroupKill|TestRunImplementExecutesSpecEndToEnd|TestRunImplementOverlapEndsIntegrationPendingAndPrintedCommandWorks|TestRunImplementPreflightClosesTerminalRunSessionsOnly|TestRunImplementPreflightPrunesRetainedRunStorage|TestRunImplementPreflightReapsEmptyTerminalRunAndTaskWorktrees|TestRunImplementPreflightRetentionPruneFailureIsNonFatal|TestRunImplementPreflightRetentionZeroSkipsPrune|TestRunImplementPreflightTerminalReachableChangedBranch|TestRunImplementPreflightTerminalUniqueChangedBranch|TestRunImplementRealWorktreeFastForwardsAndCleansPreservingNonOverlappingUserDirt|TestRunImplementTemporaryVerificationFlowPreservesDeterministicRepair|TestRunImplementTemporaryVerificationFlowRepeatedTemporaryPreservesTaskWorktree|TestRunImplementTemporaryVerificationFlowRetriesOnceWithoutAgentRepair|TestRunImplementVerificationCapacityAndDaemonStatusIntegratedFlow|TestRunImplementWorktreeIsolationExcludesConcurrentUserCommit)$' ./internal/daemon ./internal/cli 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestIntegratedRepositoryVerificationRunsOnceAfterAParallelUngatedGraph TestIntegratedRepositoryVerificationFailureLeavesTaskStatuses TestIntegratedRepositoryVerificationDoesNotRun TestIntegratedRepositoryVerificationDoesNotRun/capacity-one TestIntegratedRepositoryVerificationDoesNotRun/gated-graph TestImplementUngatedRunEndsUnresolvedWhenTheIntegratedTreeIsRed TestImplementUngatedRunEndsCleanWhenTheIntegratedTreeIsGreen TestUngatedGraphRunsRepositoryVerificationAtSettlement TestImplementGatedGraphRunsRepositoryVerificationAtSettlement TestAgentSelectionProfilesMacro TestBudgetClockCrossesAfterTheSettledTask TestRunImplementBootstrapRunsBeforeAgentWorkAndVerification TestRunImplementBootstrapsEachConcurrentTaskWorktreeBeforeAgentWork TestRunImplementDetachPrintsReportAndCompletesRun TestRunImplementDetachSurvivesCallerProcessGroupKill TestRunImplementExecutesSpecEndToEnd TestRunImplementOverlapEndsIntegrationPendingAndPrintedCommandWorks TestRunImplementPreflightClosesTerminalRunSessionsOnly TestRunImplementPreflightPrunesRetainedRunStorage TestRunImplementPreflightReapsEmptyTerminalRunAndTaskWorktrees TestRunImplementPreflightRetentionPruneFailureIsNonFatal TestRunImplementPreflightRetentionZeroSkipsPrune TestRunImplementPreflightTerminalReachableChangedBranch TestRunImplementPreflightTerminalUniqueChangedBranch TestRunImplementRealWorktreeFastForwardsAndCleansPreservingNonOverlappingUserDirt TestRunImplementTemporaryVerificationFlowPreservesDeterministicRepair TestRunImplementTemporaryVerificationFlowRepeatedTemporaryPreservesTaskWorktree TestRunImplementTemporaryVerificationFlowRetriesOnceWithoutAgentRepair TestRunImplementVerificationCapacityAndDaemonStatusIntegratedFlow TestRunImplementWorktreeIsolationExcludesConcurrentUserCommit; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done` — expected: exit 0. Before this Task the integrated-check tests do not exist, so the first pass check fails. Afterwards they pass with task_03's gate-less settlement tests, the gated implement test and the 21 implement tests task_03 migrated.

## References

- `_prd.md` → Goals; Core Feature 3; User Experience; Success Metric 5; Success Metric 6
- `_techspec.md` → Current behavior 5; Interfaces; API Contract 5; API Contract 6; Invariants 3-6; Testing Approach; Build Order 4
- ADR-0342

## Result

Implemented the integrated Run-tree check after the Task scheduler returns.
The check requires every API Contract 5 condition, uses the Run Worktree and
Verification Capacity gate, writes `verification/integrated-repository.log`,
and retries exit `75` once with exclusive capacity. Its terminal
`daemon_verification` event names the command and verdict. A red tree sets
`TaskCycleResult.IntegratedVerificationFailure` without a Task repair turn
or a Task-status change. Stop Requests and Run Budget expiry retain their
existing terminal paths.

The implement command reads that failure before integrating the Run Branch.
It keeps the Run Worktree, records `Unresolved`, publishes the failure as the
terminal completion reason, prints it before the terminal stderr line, and
returns exit `1`.

Acceptance evidence:

- Parallel ungated graph: `TestIntegratedRepositoryVerificationRunsOnceAfterAParallelUngatedGraph`
  observes both settlement events before the integrated verifier call, exactly
  one extra repository command, the Run Worktree, diagnostics file, and success
  event. The new CLI green test exercises two independent Tasks with Task
  Capacity two and real disposable Git worktrees, then asserts Clean and exit zero.
- Red integrated tree: `TestIntegratedRepositoryVerificationFailureLeavesTaskStatuses`
  asserts the exact failure reason, both completed statuses on disk, the failure
  event, and only two Agent requests. Its Tasks are dependency-ordered to avoid
  the existing scheduler fake's whole-tree copy overwriting the first Task's
  status. The CLI red test separately exercises independent parallel Tasks with
  real disposable Git worktrees: both remain completed, the user's branch stays
  unchanged, the Run Worktree stays present, exit is one, and stderr and the
  terminal outcome event carry the reason.
- Exclusions: `TestIntegratedRepositoryVerificationDoesNotRun` covers
  `capacity-one`, `one-completed`, `failed-task`, `skipped-task`, `switch-off`,
  `empty-command`, and `gated-graph`. Additional tests cover an exclusive retry
  that succeeds, an exclusive retry that again exits `75` without a third call,
  a durable Stop Request during verification, and Run Budget cancellation during
  verification.

Focused checks:

- Before implementation, `go test -timeout 2m ./internal/daemon -run '^TestIntegratedRepositoryVerification' -count=1`
  exited one because the new result field did not exist.
- `go test -timeout 3m ./internal/daemon -run '^(TestIntegratedRepositoryVerification.*|TestUngatedGraph.*|TestRepository.*Settlement.*|TestGatedGraphRunsRepositoryVerificationAfterDeclaredCommands)$' -count=1`
  exited zero on the final code (`ok roundfix/internal/daemon`, 1.945s).
- `go test -timeout 10m ./internal/cli -run '^(TestImplement.*|TestRunImplement.*|TestAgentSelectionProfilesMacro|TestBudgetClockCrossesAfterTheSettledTask)$' -count=1`
  exited zero on the final code (`ok roundfix/internal/cli`, 27.374s). This
  selection includes the migrated implement fixtures; their existing assertions
  needed no changes. An earlier run's tests passed but its suite guard rejected
  a concurrent source edit; the recorded successful run used a stable tree.
- `git diff --check` exited zero.

The declared Verification command, selected repository Verification, and full
package suites were not run in this child turn; the Daemon owns those gates.
The pre-existing `status: in_progress` change remains untouched. No other Task,
Task Graph, commit, push, or Pull Request was changed or created.

## Carry-forward provenance

- Source Run: `run_20261027T193605Z_00ac54bea6e5fdc9`
- Source commit: `5ee76ace2e541a6cd5d2d8fe92484e6462c79fd3`
