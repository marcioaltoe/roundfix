---
task: task_06
spec: 0039-review-source-evidence-and-detached-outcomes
status: completed
type: backend
complexity: high
---

# Task 06: Deliver notification receipts and Detached monitoring

## Overview

Give Detached Run users and Supervisors first-party terminal monitoring while
making every notification attempt durably observable. Notifications receive
the same terminal context as reports, return sent/skipped/failed receipts, and
remain unable to change the Run outcome.

## Requirements

1. MUST return a typed receipt with route, sent/skipped/failed status, and
   completion time for every notification attempt.
2. MUST preserve existing notification environment variables and add the
   documented terminal-context variables.
3. MUST journal exactly one receipt event after each attempt.
4. MUST keep native and command notification failures best-effort and unable
   to change Run state or exit code.
5. MUST print Detached startup as Run ID, Console Log, Attach, Supervisor
   outcome monitor, and Stop lines.
6. MUST use `roundfix events <run-id> --follow --filter outcome` as the
   Supervisor monitor command.
7. MUST bound native notification text while keeping complete structured
   context available to `notify.command`.

## Subtasks

- [x] Add notification receipt contracts.
- [x] Populate additive command environment values.
- [x] Journal one durable receipt event per attempt.
- [x] Preserve outcome on notification failure.
- [x] Add the five-line Detached startup report.
- [x] Cover sent, skipped, failed, native, and command routes.

## Acceptance Criteria

- [x] Sent, disabled/unavailable, and failed routes produce sent, skipped, and
      failed receipts respectively.
- [x] Every attempt appends exactly one route/status/completion-time event.
- [x] Existing environment variables remain byte-compatible and additive
      values carry terminal context.
- [x] Notification failure changes neither stored outcome nor top-level exit
      code.
- [x] Detached startup contains all five documented lines with runnable
      commands.
- [x] Supervisor monitoring uses the stable outcome filter and no Console Log
      parsing.
- [x] Native text stays bounded and contains the next action when non-Clean.

## Context

- instruction: `.agents/skills/agentic-cli-design/SKILL.md`
- instruction: `.agents/skills/coding-guidelines/SKILL.md`
- instruction: `.agents/skills/golang-error-handling/SKILL.md`
- instruction: `.agents/skills/golang-testing/SKILL.md`
- interface: `internal/notify/notify.go`
- interface: `internal/notify/notify_test.go`
- interface: `internal/cli/cli.go`
- interface: `internal/cli/cli_test.go`
- interface: `internal/cli/detach.go`
- interface: `internal/cli/detach_test.go`
- interface: `internal/runevent/event.go`
- interface: `internal/runevent/event_test.go`

## Verification

- `rtk go test ./internal/notify -run 'Test.*(Receipt|Environment|Native|Command)' -count=1`
  — expected: all routes produce bounded context and typed receipts.
- `rtk go test ./internal/cli ./internal/runevent -run 'Test.*(Detached.*Monitor|OutcomeNotification|NotificationReceipt)' -count=1`
  — expected: five-line startup and one durable receipt per attempt pass.
- `rtk go test -race ./internal/notify ./internal/cli -run 'Test.*(Notification|Detached)' -count=1`
  — expected: receipt and Detached reporting are race-free.

## References

- `_prd.md` → Goal 4; User Stories 4–5; Core Features 8–11; User Experience;
  Success Metrics.
- `_techspec.md` → Interfaces: NotificationReceipt; API Contracts: Terminal
  report, stream, and notification; Build Order 6.
- `CONTEXT.md` → Run Outcome Notification, Console Log, Attach, and Run Event
  Stream.

## Result

Notification delivery now returns a typed receipt for every command, native,
disabled, or unavailable attempt. Each receipt carries its route, a
`sent`/`skipped`/`failed` status, and a UTC completion time. The terminal
completion boundary passes the normalized terminal reason, Console Log, Attach
command, Review Issue knowledge, and next action into notifications, then
appends exactly one typed receipt event without changing the completed Run.

The four existing `notify.command` variables retain their names, values, and
order. Five additive variables expose the documented terminal context. Native
notification text is bounded to 256 runes and reserves space for `Next: ...`
on non-Clean outcomes.

Detached startup now prints exactly five lines: Run ID, Console Log, Attach,
the runnable
`roundfix events <run-id> --follow --filter outcome` Supervisor monitor, and
Stop.

### Verification

- `GOCACHE=/tmp/roundfix-task06-gocache rtk go test ./internal/notify -run 'Test.*(Receipt|Environment|Native|Command)' -count=1`
  — passed, 13 tests.
- `GOCACHE=/tmp/roundfix-task06-gocache rtk go test ./internal/cli ./internal/runevent -run 'Test.*(Detached.*Monitor|OutcomeNotification|NotificationReceipt)' -count=1`
  — passed, 14 tests.
- `GOCACHE=/tmp/roundfix-task06-gocache rtk go test -race ./internal/notify ./internal/cli -run 'Test.*(Notification|Detached)' -count=1`
  — passed, 24 tests.
- `GOCACHE=/tmp/roundfix-task06-gocache rtk go test ./internal/cli -run 'TestRunImplementDetach(PrintsReportAndCompletesRun|SurvivesCallerProcessGroupKill)$' -count=1`
  — passed, 2 Detached process tests.
- `GOCACHE=/tmp/roundfix-task06-gocache rtk go test ./internal/notify ./internal/cli ./internal/runevent -skip 'TestBranchIntegrityPreflightMigratesOutdatedRunDatabase|TestRunForceStopOwnerProcessIntegrationProvesExitBeforeStoreCompletion' -count=1`
  — passed, 860 tests.
- `rtk git -c core.fsmonitor=false diff --check` — passed.

### Acceptance evidence

- Command success, native success, disabled configuration, unavailable native
  tooling, command failure, native failure, and timeout cases return the
  required receipt status with a non-zero UTC completion time.
- CLI integration coverage decodes receipt payloads and proves exactly one
  route/status/completion-time event follows each sent, skipped, or failed
  attempt.
- Environment coverage compares the complete ordered environment, preserving
  the four existing entries and adding reason, Console Log, Attach command,
  Review Issue knowledge, and next action.
- A forced notifier failure preserves the Clean stored outcome, stdout report,
  and exit code while emitting one warning and one failed receipt.
- Detached unit and process coverage assert all five lines and the exact stable
  outcome-filter command; no Console Log parsing is used for Supervisor state.
- Native text coverage uses an oversized target, proves the 256-rune bound, and
  retains the non-Clean next action.

### Follow-ups

- The additional broad run confirmed the existing
  `TestBranchIntegrityPreflightMigratesOutdatedRunDatabase` schema-version
  expectation still needs reconciliation with database version 11.
- The managed sandbox still blocks
  `TestRunForceStopOwnerProcessIntegrationProvesExitBeforeStoreCompletion`
  from reading the genuine owner-process identity through `/bin/ps`.
  Both findings are outside Task 06 and were left unchanged.
