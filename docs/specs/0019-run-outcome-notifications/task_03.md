---
task: task_03
spec: 0019-run-outcome-notifications
status: pending
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

- [ ] Notifier construction from loaded config at command start
- [ ] Completion-boundary calls for resolve, watch, and implement with the
      payload derived from the completed Run
- [ ] Failure warning on stderr and the Run Event
- [ ] CLI tests: one notification per terminal outcome per command, none for
      fetch, failure path leaves report and exit code unchanged, disabled
      path is silent

## Acceptance Criteria

- [ ] A captured fake notifier records exactly one Outcome per completed
      resolve, watch, and implement Run, with the correct state, kind, and
      target — and records nothing for fetch.
- [ ] A failing notifier produces the stderr warning and the Run Event while
      the command's stdout report and exit code are byte-identical to the
      non-failing case.
- [ ] With notifications disabled, no notifier is invoked and existing CLI
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
