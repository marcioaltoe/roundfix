---
task: task_04
spec: 0020-run-browser
status: pending
type: backend
complexity: medium
---

# Task 04: Browser entry points and the attach loop

## Overview

Wire the Run Browser into its two entry points: `roundfix runs` without a
subcommand and `roundfix attach` without a run id open the browser in an
interactive terminal, selection opens the existing read-only Live Run View,
and leaving the view returns to a refreshed browser until cancel. Demoable
end to end at a terminal.

## Requirements

1. MUST open the Run Browser from bare `roundfix runs` and no-run-id
   `roundfix attach` in an interactive terminal, replacing the numbered
   Interactive Input prompt.
2. MUST open the existing attach cockpit for the selected Run with
   byte-for-byte unchanged cockpit behavior, and return to the browser with
   a fresh store query when the cockpit closes.
3. MUST exit `0` with no side effects on browser cancel.
4. MUST keep the non-interactive failures from the text-surface task exactly
   as shipped there (exit `2`, wording naming `runs list`).
5. MUST leave `attach <known-run-id>` byte-for-byte unchanged.

## Subtasks

- [ ] `runs` bare-command dispatch: browser in TTY, exit 2 wording in non-TTY
- [ ] attach no-run-id path: browser in TTY replacing the prompt picker
- [ ] Browser → cockpit → refreshed browser loop
- [ ] CLI tests through the model seam: loop selection, cancel, refresh,
      unchanged explicit-id attach

## Acceptance Criteria

- [ ] Selecting a Run in the browser opens the same Live Run View bytes the
      explicit `attach <run-id>` produces, and closing it returns to the
      browser with current data.
- [ ] Cancelling the browser exits `0` and mutates nothing.
- [ ] Bare `runs` in a non-TTY context exits `2` naming `runs list`.
- [ ] Existing explicit-id attach tests pass unchanged.

## Verification

- `rtk go test ./internal/cli/` — expected: all tests pass, including the new
  entry-point and loop tests.
- `rtk make verify` — expected: fmt-check, tests, skills check, and build all
  pass.

## References

`_prd.md` → User Stories 1, 3; Core Features 1-3, 5. `_techspec.md` → System
Architecture (flow); API Contracts: runs, attach; Build Order 4.
