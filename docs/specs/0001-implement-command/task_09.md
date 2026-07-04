---
task: task_09
spec: 0001-implement-command
status: pending
type: frontend
complexity: medium
---

# Task 09: Render Tasks as Work Items in the Live Run View

## Overview

Make spec Runs first-class in the Live Run View and Attach: the work-item pane renders Tasks with their statuses in graph order for `implement` Runs, while review Runs keep their exact rendering. Verifiable by driving the cockpit model synchronously and by attaching to a finished spec Run in CLI tests.

## Requirements

1. MUST introduce a small work-item view model that both Review Issues and Tasks map into, keyed on the Run Kind — review rendering (grouped by Round) stays byte-identical.
2. MUST render Tasks for `implement` Runs in Task Graph order with their current statuses, refreshed by re-reading task files located through the Run row's git root and spec slug.
3. MUST tolerate mid-write task files: a parse failure keeps the last good status and never fails the view.
4. MUST keep Attach working unchanged by run id for spec Runs: replay the Run Event Journal, then follow, with the Task pane populated the same way — read-only, never mutating the Run (ADR-0009 discipline: the view reads the journal, never the live sink).
5. MUST render `daemon.task` and `daemon.qa` events meaningfully in the timeline (Task id, phase, verdict) in both the cockpit and the plain-text renderer.

## Subtasks

- [ ] Work-item view model and Run-Kind-aware pane population
- [ ] Task status refresh from task files with mid-write tolerance
- [ ] Timeline rendering for the task and qa event Kinds
- [ ] Plain-text Live Run View parity for non-TTY output
- [ ] Attach coverage over a finished spec Run

## Acceptance Criteria

- [ ] Driving the cockpit model synchronously over a spec Run's journal shows Tasks in graph order with statuses updating as `daemon.task` events arrive.
- [ ] A corrupted task file mid-refresh leaves the previously shown status in place; no error surfaces in the view.
- [ ] Review-Run rendering output is unchanged under the existing tests.
- [ ] `attach` on a completed spec Run replays Task lines and the outcome without touching Run state.
- [ ] The full existing suite passes unchanged.

## Verification

- `rtk go test ./internal/tui/ ./internal/cli/` — expected: all tests pass.
- `rtk go test ./...` — expected: full suite passes.

## References

`_prd.md` → User Story 6; Core Feature 11; User Experience. `_techspec.md` → System Architecture (tui), Build Order 10, Risks (mid-write reads). ADR-0009.
