---
task: task_03
spec: 0017-run-discovery
status: pending
type: backend
complexity: medium
---

# Task 03: Attach picker: no-argument Interactive Input and no-input error

## Overview

Make `roundfix attach` without a run id usable: in an interactive terminal it
opens Interactive Input listing the repository's Runs and attaches to the
selection; in a non-interactive context it fails with an error naming the
listing command. Demoable by running `attach` with no arguments against seeded
Runs.

## Requirements

1. MUST open Interactive Input when `attach` receives no run id in an
   interactive terminal: a numbered list of the repository's Runs, newest
   first with Active Runs first, each entry showing state, kind, and target,
   accepting a number or a run id.
2. MUST attach to the selected Run through the existing Attach path with
   byte-for-byte unchanged behavior; `attach <run-id>` is untouched.
3. MUST list terminal Runs too — Attach replays them from the Run Event
   Journal.
4. MUST exit `0` with no side effects when the picker is cancelled.
5. MUST fail with exit `2` in a non-interactive context (no TTY or
   `--no-input`), with an error naming the listing command as the discovery
   alternative.
6. MUST reuse the listing query from the store — no second query shape.

## Subtasks

- [ ] Split attach validation: missing run id becomes picker (interactive) or
      actionable error (non-interactive)
- [ ] Picker rendering and selection following the Interactive Input pattern
- [ ] Wire selection into the existing attach path
- [ ] CLI tests: scripted picker input (number, run id, cancel), no-input
      error, unchanged explicit-id behavior

## Acceptance Criteria

- [ ] With seeded Runs, `attach` with no arguments in an interactive context
      lists them (Active first, newest first) and attaches to the selection by
      number and by run id.
- [ ] Cancelling the picker exits `0` without attaching.
- [ ] `attach` with no arguments and no TTY (or `--no-input`) exits `2` and the
      error names the listing command.
- [ ] `attach <run-id>` output is unchanged for an existing Run.

## Verification

- `rtk go test ./internal/cli/` — expected: all tests pass, including the new
  picker and no-input tests.
- `rtk make verify` — expected: fmt-check, tests, skills check, and build all
  pass.

## References

`_prd.md` → User Story 3; Core Features 5-6; User Experience. `_techspec.md` →
Interfaces: pickAttachRun; API Contracts: attach without a run id; Build
Order 3.
