---
task: task_03
spec: 0017-run-discovery
status: completed
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

- [x] Split attach validation: missing run id becomes picker (interactive) or
      actionable error (non-interactive)
- [x] Picker rendering and selection following the Interactive Input pattern
- [x] Wire selection into the existing attach path
- [x] CLI tests: scripted picker input (number, run id, cancel), no-input
      error, unchanged explicit-id behavior

## Acceptance Criteria

- [x] With seeded Runs, `attach` with no arguments in an interactive context
      lists them (Active first, newest first) and attaches to the selection by
      number and by run id.
- [x] Cancelling the picker exits `0` without attaching.
- [x] `attach` with no arguments and no TTY (or `--no-input`) exits `2` and the
      error names the listing command.
- [x] `attach <run-id>` output is unchanged for an existing Run.

## Verification

- `rtk go test ./internal/cli/` — expected: all tests pass, including the new
  picker and no-input tests.
- `rtk make verify` — expected: fmt-check, tests, skills check, and build all
  pass.

## References

`_prd.md` → User Story 3; Core Features 5-6; User Experience. `_techspec.md` →
Interfaces: pickAttachRun; API Contracts: attach without a run id; Build
Order 3.

## Result

Implemented no-argument `roundfix attach` as an Interactive Input picker over
the store `ListRuns` query. The picker scopes to the current repository,
orders Active Runs before terminal Runs while preserving newest-first order
within each group, lists terminal Runs too, accepts either a number or a Run
id, and then passes the selected Run id into the existing Attach path. The
explicit `attach <run-id>` path is unchanged.

Non-interactive no-argument attach now exits `2` before opening the Run
Database and names `roundfix runs list` as the discovery alternative.
Cancelling the picker exits `0` without replaying or following any Run.

Acceptance evidence:

- `TestAttachPickerSelectsRunByNumberActiveFirstNewestFirst` seeds Active and
  terminal Runs across two repositories, verifies the picker lists only the
  current repository with Active Runs first and newest-first ordering, selects
  by number, and observes Attach replay for the selected terminal Run.
- `TestAttachPickerSelectsRunByIDWithExplicitAttachOutput` selects by Run id
  and compares the picker attach stdout byte-for-byte with explicit
  `attach <run-id>` stdout.
- `TestAttachPickerCancelExitsZeroWithoutAttaching` cancels with a blank
  picker input, verifies exit `0`, no attach stdout, unchanged Run count, and
  unchanged terminal Run state.
- `TestAttachWithoutRunIDNonInteractiveNamesRunsList` covers both no TTY and
  `--no-input` paths, expecting exit `2`, no stdout, and stderr naming
  `roundfix runs list`.

Verification:

- `rtk go test ./internal/cli/` passed: 309 CLI tests.
- `rtk make verify` passed: `rtk go test ./...` reported 838 tests across 18
  packages, `roundfix skills check` passed, and `go build` completed.
