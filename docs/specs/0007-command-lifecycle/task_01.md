---
task: task_01
spec: 0007-command-lifecycle
status: completed
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

- [x] Schema v5 migration with populated v4 fixture test
- [x] Store request/read methods with terminal-Run refusal
- [x] Stop command graceful default and report
- [x] Engine boundary checks in both cycles, QA step included
- [x] Stopped-outcome and lock-release tests

## Acceptance Criteria

- [x] Migration test: populated v4 fixture (runs in several states + one
      active lock) opens as v5 with all rows intact.
- [x] Engine tests: a Stop Request set mid-Task lets that Task verify,
      settle, and commit before the Run ends Stopped; a Request set before
      the QA step skips QA entirely; the resolve cycle mirrors both.
- [x] Stop command tests: Active Run → request recorded + report; terminal
      Run → the named refusal; selectors and exit codes unchanged.
- [x] `rtk go test -race ./internal/daemon/` clean (new polled read).
- [x] Full suite passes.

## Verification

- `rtk go test ./internal/store/ ./internal/daemon/ ./internal/cli/` —
  expected: all tests pass.
- `rtk go test -race ./internal/daemon/` — expected: no races.
- `rtk go test ./...` — expected: full suite passes.

## References

`_prd.md` → User Story 4; Core Feature 3. `_techspec.md` → Stop semantics,
Data Models, Build Order 1. ADR-0004, ADR-0022. Round-1 dogfood finding 24.

## Result

Implemented schema v5 with nullable `runs.stop_requested_at`, v4→v5 and
fresh-create paths, plus `RequestStop` and `StopRequested`. Terminal Runs now
refuse Stop Requests with `ErrTerminalRunStopRequest`.

`roundfix stop` now records a Stop Request by default for every selector and
reports `Stop Request recorded; the Run stops after the current Work Item
settles.` The report names `--force` for dead or runaway Runs; `--force` keeps
the immediate fallback path wired for task_04.

Both engines poll the Run Database Stop Request flag before new work and after
settlement boundaries. TaskCycle also polls before QA, so a Stop Request before
QA skips the QA step. The stop is journaled with the existing daemon status
event and surfaces through the existing Stopped completion path.

Evidence:

- Migration/store: `TestOpenMigratesV4RunDatabasePreservingRunsLocksAndAddingStopRequests`,
  `TestRequestStopRecordsStopRequestForActiveRun`, and
  `TestRequestStopRejectsTerminalRunWithNamedError` passed in
  `rtk go test ./internal/store/ ./internal/daemon/ ./internal/cli/`.
- Engine boundaries: `TestTaskCycleStopRequestAfterTaskSettlementHaltsBeforeNextTask`,
  `TestTaskCycleStopRequestBeforeQAStepSkipsQA`,
  `TestResolveCycleStopRequestBeforeBatchPublishesStopAndDoesNothing`, and
  `TestResolveCycleStopRequestAfterBatchSettlementHaltsBeforeNextBatch` passed
  in the focused package test.
- Stop command: `TestRunStopByRunIDRecordsStopRequest`,
  `TestRunStopByPullRequestRecordsStopRequestForMatchingActiveRun`,
  `TestRunStopBySpecRecordsStopRequestForMatchingActiveRun`, and
  `TestRunStopRejectsAlreadyTerminalRun` passed in the focused package test.
- Stopped outcome and lock release:
  `TestRunImplementDatabaseStopRequestAfterTaskCommitEndsStoppedAndReleasesLock`
  passed in the focused package test.
- `rtk go test ./internal/store/ ./internal/daemon/ ./internal/cli/`: passed,
  288 tests across 3 packages.
- `rtk go test -race ./internal/daemon/`: passed, 53 tests in 1 package.
- `rtk go test ./...`: passed, 613 tests across 16 packages.
- `rtk make verify`: passed; it ran `rtk go test ./...`,
  `rtk go run -buildvcs=false ./cmd/roundfix skills check`, and
  `rtk go build -buildvcs=false -o bin/roundfix ./cmd/roundfix`.

Follow-up: task_04 still owns the cooperative Agent Session cancel behavior
behind `roundfix stop --force`.
