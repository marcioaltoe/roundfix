---
task: task_05
spec: 0037-terminal-outcome-integrity
status: completed
type: backend
complexity: high
---

# Task 05: Publish only the winning terminal outcome

## Overview

Carry the completion transition result through all terminal publication paths
so only the winner emits outcome evidence and notification. Preserve primary
failure ordering while retaining registered-session cleanup problems as
explicit secondary diagnostics.

## Requirements

1. MUST publish a terminal outcome Run Event only when completion transitioned
   the Run.
2. MUST attempt Run Outcome Notification only for that same winning transition.
3. MUST make identical replays and conflicting losers observe stored state
   without duplicate publication.
4. MUST preserve one terminal outcome across Resolve, Watch, Implement, Stop,
   and cleanup paths.
5. MUST print and journal the primary failure before any cleanup warning.
6. MUST label cleanup failures as secondary and prevent them from replacing
   the terminal reason or exit code.
7. MUST retain best-effort notification failure without changing the Run
   outcome.

## Subtasks

- [x] Thread the completion transition result through terminal coordinators.
- [x] Gate outcome events and notifications on the winning transition.
- [x] Normalize replay and conflict handling around stored state.
- [x] Order primary failure evidence before cleanup diagnostics.
- [x] Cover terminal paths across operational commands.
- [x] Add deterministic owner-versus-Force-Stop publication coverage.

## Acceptance Criteria

- [x] A completion race stores and publishes exactly one terminal outcome.
- [x] The losing owner emits no second outcome event or notification.
- [x] An identical replay is silent and returns the stored result.
- [x] A conflicting loser cannot change the primary reason, exit code, or
      notification context.
- [x] Cleanup warnings follow the primary failure and are labeled secondary.
- [x] Notification failure remains a warning and leaves persisted outcome
      unchanged.
- [x] Resolve, Watch, Implement, and Stop terminal regressions pass.

## Context

- instruction: `.agents/skills/coding-guidelines/SKILL.md`
- instruction: `.agents/skills/golang-error-handling/SKILL.md`
- instruction: `.agents/skills/golang-testing/SKILL.md`
- instruction: `.agents/skills/testing-boss/SKILL.md`
- interface: `internal/cli/cli.go`
- interface: `internal/cli/cli_test.go`
- interface: `internal/cli/implement.go`
- interface: `internal/cli/implement_test.go`
- interface: `internal/store/journal.go`
- interface: `internal/store/journal_test.go`
- interface: `internal/notify/notify.go`
- interface: `internal/notify/notify_test.go`
- interface: `internal/runevent/event.go`

## Verification

- `rtk go test ./internal/cli ./internal/store ./internal/notify ./internal/runevent -run 'Test.*(CompletionWinner|TerminalOutcome|PrimaryFailure|OutcomeNotification)' -count=1`
  — expected: only the completion winner publishes and primary diagnostics
  precede secondary cleanup warnings.
- `rtk go test -race ./internal/cli ./internal/store -run 'Test.*(CompletionWinner|TerminalOutcome)' -count=1`
  — expected: concurrent terminal publication has one race-free winner.

## References

- `_prd.md` → Goals 1 and 4; User Stories 3–4; Core Features 6–7; User
  Experience; Success Metrics.
- `_techspec.md` → System Architecture: Run Event Journal; Testing Approach;
  Build Order 5.
- `../../adr/0052-run-completion-is-compare-and-set.md` → winner-only
  publication.

## Result

Terminal coordinators now carry `CompleteRunResult.Transitioned` through
Resolve, Watch, Implement, failure, graceful-stop, and Force Stop paths. Only a
winning transition appends `daemon.outcome` and attempts the Run Outcome
Notification. Identical replays return the stored Run without publication, and
intermediate owner-state updates use compare-and-set protection so they cannot
reopen a terminal Run before attempting a conflicting completion.

Force Stop retains Agent Session cleanup failures when owner termination or
completion fails. It prints and journals the primary failure first, then emits
each failed or skipped cleanup diagnostic as a labeled
`Secondary cleanup warning`; those diagnostics do not change the Run state or
command exit code. Notification delivery remains best-effort and records its
own warning after the persisted outcome.

Verification:

- Pre-change
  `rtk env GOCACHE=/private/tmp/roundfix-task05-gocache go test ./internal/cli -run 'Test(CompletionWinner|RunForceStopPrimaryFailure)' -count=1`
  failed because Force Stop emitted no outcome notification or event, the
  losing owner could rewrite Stopped through an intermediate state and then
  store Clean, and cleanup failures disappeared behind the owner failure.
- Pre-change
  `rtk env GOCACHE=/private/tmp/roundfix-task05-gocache go test ./internal/store -run 'TestTerminalOutcomeRejectsIntermediateStateUpdate' -count=1`
  failed because `UpdateRunState` accepted a terminal-to-intermediate rewrite.
- `rtk env GOCACHE=/private/tmp/roundfix-task05-gocache go test ./internal/cli ./internal/store ./internal/notify ./internal/runevent -run 'Test.*(CompletionWinner|TerminalOutcome|PrimaryFailure|OutcomeNotification)' -count=1`
  passed.
- `rtk env GOCACHE=/private/tmp/roundfix-task05-gocache go test -race ./internal/cli ./internal/store -run 'Test.*(CompletionWinner|TerminalOutcome)' -count=1`
  passed.
- `rtk env GOCACHE=/private/tmp/roundfix-task05-gocache go test ./internal/cli ./internal/store ./internal/notify ./internal/runevent -count=1`
  passed.
- `rtk env GOCACHE=/private/tmp/roundfix-task05-gocache make verify` passed:
  2,477 Go tests, the four protected skill tests, the Roundfix skill sync
  check, and the CLI build.
- `rtk git -c core.fsmonitor=false diff --check` passed.

Acceptance evidence:

- `TestCompletionWinnerOwnerVersusForceStopPublishesOneTerminalOutcome`
  deterministically pauses a Resolve owner, lets Force Stop win, then resumes
  the owner. The stored outcome remains Stopped, the owner exits with the Run
  failure code, and the journal and notifier each contain exactly one Stopped
  outcome. Replaying the same Force Stop adds neither an event nor a
  notification.
- `TestTerminalOutcomeRejectsIntermediateStateUpdate` proves a terminal Run
  rejects the losing owner's later intermediate state, preserves its completion
  timestamp, and returns the stored/requested states in
  `TerminalOutcomeConflictError`.
- `TestRunForceStopPrimaryFailurePrecedesSecondaryCleanupWarnings` proves
  printed and journaled owner failure evidence precedes labeled cleanup
  warnings while the Run stays Active and retains its primary exit behavior.
- `TestRunOutcomeNotificationFailureWarnsAndJournalsWithoutChangingReportOrExit`
  proves notifier failure stays a warning and preserves the persisted outcome,
  report, and exit code.
- The focused terminal suite covers the existing Resolve, Watch, Implement,
  Stop, notification, and Run Event regressions; the affected-package run and
  full repository gate passed afterward.

Follow-up: none.
