# Focused regression evidence

Build: `ef6eb44ad8951112b1c3641bb7fd21793b440f95`

The full static gate passed first. These narrower reruns name the exact
behavior that supports Task-level criteria not safely created through live
destructive failure injection.

## Store and reconciliation

The focused Store run passed:

- `TestCompleteRunWinnerAndIdenticalReplay`
- `TestTerminalOutcomeConflictPreservesWinner`
- `TestTerminalOutcomeRejectsIntermediateStateUpdate`
- `TestTerminalOutcomeEveryStoredTerminalStateIsImmutable`, including all ten
  stored terminal states
- `TestCompleteRunConcurrentTerminalOutcomesHaveOneWinner`
- `TestReconcileIntegrationPendingRecordsEvidence`
- `TestReconcileIntegrationRejectsIncompleteEvidence`
- `TestReconcileIntegrationRejectsStaleTargetBranch`
- `TestReconcileIntegrationRejectsEveryOtherSourceOutcome`
- `TestReconcileIntegrationRollsBackWhenJournalFails`
- `TestAgentSelectionActiveScopesReturnsLatestLifecycleInStableOrder`
- `TestStoppedRunReleasesActiveLock`

Command:

```text
rtk env GOCACHE=/private/tmp/roundfix-qa-0037-focused-gocache go test -v
./internal/store -run 'Test(...)' -count=1
```

Result: exit `0`.

## Force Stop, publication, and diagnostics

The focused CLI run passed:

- `TestCompletionWinnerOwnerVersusForceStopPublishesOneTerminalOutcome`
- `TestRunOutcomeNotificationFailureWarnsAndJournalsWithoutChangingReportOrExit`
- `TestRunWatchStopRequestBeforeAgentMarksStopped`
- `TestRunForceStopOwnerPermissionAndDeadlineFailuresRetainActiveLock`
- `TestRunForceStopPrimaryFailurePrecedesSecondaryCleanupWarnings`
- `TestRunForceStopOwnerPIDReuseFailsClosed`
- `TestRunForceStopStoppedRunIsIdempotentWithoutOwnerOrSessionActions`

The registered-session run passed:

- `TestRunStopForceAgentSessionCleanupSkipsRunWithoutActiveLifecycle`
- `TestRunStopForceRegisteredAgentSessionCleanupTargetsActiveScopesInOrder`
- `TestRunStopForceRegisteredAgentSessionAbsenceIsIdempotent`
- `TestRunStopForceAgentSessionCleanupFailureRemainsVisibleWithoutClosedLifecycle`

The real helper-process run passed:

- `TestRunForceStopOwnerProcessIntegrationProvesExitBeforeStoreCompletion`
- `TestOwnerProcessControllerGracefulExitProof`
- `TestOwnerProcessControllerForceKillExitProof`
- `TestOwnerProcessControllerRejectsUnprovenCurrentProcess`

## Stop Request boundaries

The focused Watch run passed:

- `TestRunStopRequestDuringStatusWaitStopsAtNextPoll`
- `TestRunStopRequestDuringQuietPeriodStopsBeforeFetch`
- `TestRunStopRequestDuringTransientRetryStopsBeforeNextCheck`
- `TestRunStopRequestDuringMergeReadyWaitStopsBeforeNextCheck`
- `TestRunStopRequestSourceFailureIncludesRunAndOperation`
- `TestRunWithoutStopRequestKeepsRunBudgetBehavior`

The in-flight Work Item contract passed:

- `TestRunImplementStopRequestEndsStoppedWithInterruptMapping`
- `TestRunImplementDatabaseStopRequestAfterTaskCommitEndsStoppedAndReleasesLock`
- `TestTaskCycleStopRequestMidWaveDrainsRunningTasksAndStartsNothingNew`
- `TestTaskCycleStopRequestAfterTaskSettlementHaltsBeforeNextTask`

## Integration recovery and supported platforms

Both Settle recovery regressions passed:

- `TestRunSettleRetargetsKeptRunWorktreeAndCleansUpAfterIntegration`
- `TestRunSettleRetargetsKeptTaskWorktreeAndCleansUpAfterIntegration`

The built CLI listed existing Integration Pending Runs and replayed the
terminal Event Stream for
`run_20260715T125004Z_bceaf37742fa163d`, independently confirming the public
state remains readable. `roundfix settle --help` states that Settle writes no
Run Event Journal entry and never pushes.

The supported-package cross-build passed:

```text
rtk env GOCACHE=/private/tmp/roundfix-qa-0037-focused-gocache \
  GOOS=windows GOARCH=amd64 \
  go build -buildvcs=false ./internal/store ./internal/cli
```

The race run passed for Store, CLI, and Watch:

```text
rtk env GOCACHE=/private/tmp/roundfix-qa-0037-focused-gocache \
  go test -race ./internal/store ./internal/cli ./internal/watch \
  -run 'Test.*(CompleteRun|TerminalOutcome|ForceStop|OwnerProcess|StopRequest|CompletionWinner)' \
  -count=1
```

Verdict: pass.
