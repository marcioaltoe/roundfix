---
task: task_04
spec: 0007-command-lifecycle
status: completed
type: backend
complexity: medium
---

# Task 04: Force stop with cooperative Agent Session cancel

## Overview

Complete ADR-0022's immediate half: `roundfix stop --force` sends the
cooperative cancel through the target Run's Agent Session, then completes
the Run in the database and releases its locks — the recovery path for dead
or runaway processes. Verifiable through invocation-capture tests over the
fake acpx rig and store assertions.

## Requirements

1. MUST add `--force` to the stop command, composing with every existing
   selector: best-effort `acpx <agent> cancel -s roundfix-<run-id>` for the
   target Run (agent derived from the Run's stored runtime identity; the
   session name derives from the run id by construction), then immediate
   database completion as Stopped and lock release — today's force
   behavior plus the cancel.
2. MUST keep the cancel best-effort: a failing or impossible cancel (no
   runtime recorded, acpx absent, session gone) is reported on stderr and
   never blocks the completion.
3. MUST make the graceful default (task_01) and `--force` mutually
   coherent: `--force` works on Runs with or without a pending Stop
   Request; the stop report distinguishes the two paths.
4. MUST document in help text when to use `--force` (dead process, stuck
   Run) versus the graceful default.

## Subtasks

- [x] `--force` flag composing with all selectors
- [x] Cancel invocation derivation and best-effort execution
- [x] Immediate completion and lock release with store assertions
- [x] Help text and path-distinction reports

## Acceptance Criteria

- [x] Force-stopping an Active implement Run issues the exact cancel
      invocation (captured via the rig), completes the Run Stopped, and
      releases the target and worktree locks — provable by creating a new
      Run immediately after.
- [x] Cancel failure paths (missing acpx, unknown runtime) still complete
      the Run, with the stderr note asserted.
- [x] Graceful-then-force sequence behaves: request recorded, force
      completes immediately.
- [x] `stop --help` documents the split; full suite passes.

## Verification

- `rtk go test ./internal/cli/` — expected: all tests pass.
- `rtk go run ./cmd/roundfix stop --help` — expected: `--force` documented,
  exit 0.
- `rtk go test ./...` — expected: full suite passes.

## References

`_prd.md` → User Story 5; Core Feature 4. `_techspec.md` → Stop semantics,
Build Order 4. ADR-0022. Round-1 dogfood finding 24.

## Result

- Added force-stop Agent Session cancel through the existing acpx runner path:
  `CancelSession` builds `--cwd <workdir> <agent> cancel -s roundfix-<run-id>`
  and `roundfix stop --force` calls it best-effort before completing the Run
  Stopped.
- Persisted the selected Agent identity on newly created resolve, watch, and
  implement Runs so a later `stop --force` process can derive the cancel
  runtime. Existing rows keep an empty value and report a non-blocking cancel
  skipped note.
- Added force-stop reports that distinguish immediate force completion from
  the graceful Stop Request path, and updated `stop --help` plus the Roundfix
  skill docs to describe when `--force` is appropriate.
- Acceptance evidence: `TestACPXCancelSessionInvokesSessionCancel` captures
  the exact acpx cancel argv, and
  `TestRunStopForceCancelsImplementRunAndReleasesLocks` proves an Active
  implement Run is completed Stopped and the Spec lock is released by creating
  a new Run immediately after.
- Acceptance evidence:
  `TestRunStopForceReportsCancelFailuresButCompletes` covers missing acpx,
  unknown runtime, and no recorded runtime; all complete Stopped and assert the
  stderr note.
- Acceptance evidence: `TestRunStopGracefulThenForceCompletesImmediately`
  records a graceful Stop Request first, then force-stops the same Run
  immediately and proves the lock release.
- Acceptance evidence: `TestRunStopHelpListsSpecSelector` asserts the help
  text names graceful stop and `--force`; `rtk go run ./cmd/roundfix stop
  --help` exited 0 and printed the documented split.
- Verification: `rtk go test ./internal/cli/` passed
  (`223 passed in 1 packages`).
- Verification: `rtk go run ./cmd/roundfix stop --help` passed and printed
  `--force` guidance.
- Verification: `rtk go test ./...` passed (`647 passed in 16 packages`).
- Verification: `rtk make verify` passed, including full tests, skill check,
  and build.
