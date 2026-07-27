---
task: task_06
spec: 0039-review-source-evidence-and-detached-outcomes
status: pending
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

- [ ] Add notification receipt contracts.
- [ ] Populate additive command environment values.
- [ ] Journal one durable receipt event per attempt.
- [ ] Preserve outcome on notification failure.
- [ ] Add the five-line Detached startup report.
- [ ] Cover sent, skipped, failed, native, and command routes.

## Acceptance Criteria

- [ ] Sent, disabled/unavailable, and failed routes produce sent, skipped, and
      failed receipts respectively.
- [ ] Every attempt appends exactly one route/status/completion-time event.
- [ ] Existing environment variables remain byte-compatible and additive
      values carry terminal context.
- [ ] Notification failure changes neither stored outcome nor top-level exit
      code.
- [ ] Detached startup contains all five documented lines with runnable
      commands.
- [ ] Supervisor monitoring uses the stable outcome filter and no Console Log
      parsing.
- [ ] Native text stays bounded and contains the next action when non-Clean.

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
