---
task: task_04
spec: 0007-command-lifecycle
status: pending
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

- [ ] `--force` flag composing with all selectors
- [ ] Cancel invocation derivation and best-effort execution
- [ ] Immediate completion and lock release with store assertions
- [ ] Help text and path-distinction reports

## Acceptance Criteria

- [ ] Force-stopping an Active implement Run issues the exact cancel
      invocation (captured via the rig), completes the Run Stopped, and
      releases the target and worktree locks — provable by creating a new
      Run immediately after.
- [ ] Cancel failure paths (missing acpx, unknown runtime) still complete
      the Run, with the stderr note asserted.
- [ ] Graceful-then-force sequence behaves: request recorded, force
      completes immediately.
- [ ] `stop --help` documents the split; full suite passes.

## Verification

- `rtk go test ./internal/cli/` — expected: all tests pass.
- `rtk go run ./cmd/roundfix stop --help` — expected: `--force` documented,
  exit 0.
- `rtk go test ./...` — expected: full suite passes.

## References

`_prd.md` → User Story 5; Core Feature 4. `_techspec.md` → Stop semantics,
Build Order 4. ADR-0022. Round-1 dogfood finding 24.
