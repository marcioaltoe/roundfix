---
task: task_01
spec: 0007-command-lifecycle
status: pending
type: backend
complexity: high
---

# Task 01: Route Stop Requests through the Run Database

## Overview

Implement ADR-0022's graceful half: `roundfix stop` on a live Active Run
records a Stop Request in the Run Database (schema v5), and both engines
honor it at their next settlement boundary — the in-flight Work Item still
verifies, settles, and commits, then the Run ends Stopped. Verifiable
through store migration tests and engine boundary tests.

## Requirements

1. MUST migrate the Run Database to schema v5: `runs` gains nullable
   `stop_requested_at`, existing rows and locks survive untouched,
   `user_version` becomes 5; a fresh database creates v5 directly.
2. MUST add store methods to record a Stop Request for an Active Run and to
   read the flag; requesting a stop on a terminal Run is a named error.
3. MUST make `roundfix stop` (all selectors) default to recording the Stop
   Request when the target Run is Active, reporting `Stop Request recorded;
   the Run stops after the current Work Item settles.` and naming `--force`
   for dead processes; today's immediate completion moves behind `--force`
   (delivered in task_04 — this task keeps a bare immediate fallback flag
   wiring compilable but the graceful path is the default).
4. MUST make both engines (resolve cycle and TaskCycle) check the flag
   after each Work Item settlement and before starting the QA step; when
   set: journal the stop through existing event kinds, end the Run Stopped
   through the existing paths, with locks released by completion.
5. MUST keep Ctrl-C in-terminal behavior and every exit-code mapping
   unchanged; a gracefully stopped Run reports the standard Stopped
   outcome.

## Subtasks

- [ ] Schema v5 migration with populated v4 fixture test
- [ ] Store request/read methods with terminal-Run refusal
- [ ] Stop command graceful default and report
- [ ] Engine boundary checks in both cycles, QA step included
- [ ] Stopped-outcome and lock-release tests

## Acceptance Criteria

- [ ] Migration test: populated v4 fixture (runs in several states + one
      active lock) opens as v5 with all rows intact.
- [ ] Engine tests: a Stop Request set mid-Task lets that Task verify,
      settle, and commit before the Run ends Stopped; a Request set before
      the QA step skips QA entirely; the resolve cycle mirrors both.
- [ ] Stop command tests: Active Run → request recorded + report; terminal
      Run → the named refusal; selectors and exit codes unchanged.
- [ ] `rtk go test -race ./internal/daemon/` clean (new polled read).
- [ ] Full suite passes.

## Verification

- `rtk go test ./internal/store/ ./internal/daemon/ ./internal/cli/` —
  expected: all tests pass.
- `rtk go test -race ./internal/daemon/` — expected: no races.
- `rtk go test ./...` — expected: full suite passes.

## References

`_prd.md` → User Story 4; Core Feature 3. `_techspec.md` → Stop semantics,
Data Models, Build Order 1. ADR-0004, ADR-0022. Round-1 dogfood finding 24.
