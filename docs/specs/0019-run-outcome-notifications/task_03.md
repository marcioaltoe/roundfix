---
task: task_03
spec: 0019-run-outcome-notifications
status: completed
type: backend
complexity: medium
---

# Task 03: Terminal-outcome wiring with best-effort warning reporting

## Overview

Fire one notification when `resolve`, `watch`, or `implement` completes a Run
with a terminal outcome — including Detached Runs, which notify from the
detached process — and report notifier failures as a stderr warning plus a
Daemon-source Run Event without touching the Run's report, outcome, or exit
code. Demoable end to end with a fake notifier in CLI tests.

## Requirements

1. MUST invoke the notifier exactly once after a successful Run completion at
   the terminal-outcome boundaries of `resolve`, `watch`, and `implement`,
   with the Run id, terminal state, kind, and target. Fetch completion,
   settle, and archive MUST NOT notify.
2. MUST notify from the process that owns the Run completion, so Detached
   Runs notify from the detached child.
3. MUST report a notifier failure as one stderr line shaped like
   `roundfix: outcome notification failed: <reason>` and one Daemon-source
   Run Event, leaving the report bytes, outcome, and exit code unchanged.
4. MUST keep `notify.enabled: false` byte-for-byte identical to today's
   behavior.
5. MUST inject the notifier through a seam CLI tests can capture — no desktop
   side effects in tests.

## Subtasks

- [x] Notifier construction from loaded config at command start
- [x] Completion-boundary calls for resolve, watch, and implement with the
      payload derived from the completed Run
- [x] Failure warning on stderr and the Run Event
- [x] CLI tests: one notification per terminal outcome per command, none for
      fetch, failure path leaves report and exit code unchanged, disabled
      path is silent

## Acceptance Criteria

- [x] A captured fake notifier records exactly one Outcome per completed
      resolve, watch, and implement Run, with the correct state, kind, and
      target — and records nothing for fetch.
- [x] A failing notifier produces the stderr warning and the Run Event while
      the command's stdout report and exit code are byte-identical to the
      non-failing case.
- [x] With notifications disabled, no notifier is invoked and existing CLI
      tests pass unchanged.

## Verification

- `rtk go test ./internal/cli/` — expected: all tests pass, including the new
  notification wiring tests.
- `rtk make verify` — expected: fmt-check, tests, skills check, and build all
  pass.

## References

`_prd.md` → User Story 1; Core Features 1, 5-6. `_techspec.md` → System
Architecture (flow); Coverage Map; Build Order 3; Risks (post-completion
firing).

## Result

Implemented task_03.

- Acceptance 1: `TestRunOutcomeNotificationsCaptureTerminalResolveWatchAndImplement` records one fake-notifier `Outcome` for completed `resolve`, `watch`, and `implement` Runs with the completed Run id, state, kind, and `pr:<number>`/`spec:<slug>` target. `TestRunOutcomeNotificationsSkipFetch` verifies `fetch` records no notification.
- Acceptance 2: `TestRunOutcomeNotificationFailureWarnsAndJournalsWithoutChangingReportOrExit` verifies the exact stderr warning shape, a Daemon-source Run Event, and byte-identical stdout plus unchanged exit code versus a non-failing notifier.
- Acceptance 3: `TestRunOutcomeNotificationsDisabledSkipsNotifier` sets `notify.enabled: false`, fails if the notifier factory is constructed, and verifies the normal stdout stays unchanged and stderr has no notification warning. The full CLI package also passed unchanged.
- Verification: `rtk go test ./internal/cli/` passed with 320 tests in 1 package.
- Verification: `rtk make verify` exited 0; it ran `rtk go test ./...` with 878 tests in 19 packages, `roundfix skills check`, and `go build`.
