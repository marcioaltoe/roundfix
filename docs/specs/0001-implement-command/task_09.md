---
task: task_09
spec: 0001-implement-command
status: completed
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

- [x] Work-item view model and Run-Kind-aware pane population
- [x] Task status refresh from task files with mid-write tolerance
- [x] Timeline rendering for the task and qa event Kinds
- [x] Plain-text Live Run View parity for non-TTY output
- [x] Attach coverage over a finished spec Run

## Acceptance Criteria

- [x] Driving the cockpit model synchronously over a spec Run's journal shows Tasks in graph order with statuses updating as `daemon.task` events arrive.
- [x] A corrupted task file mid-refresh leaves the previously shown status in place; no error surfaces in the view.
- [x] Review-Run rendering output is unchanged under the existing tests.
- [x] `attach` on a completed spec Run replays Task lines and the outcome without touching Run state.
- [x] The full existing suite passes unchanged.

## Verification

- `rtk go test ./internal/tui/ ./internal/cli/` — expected: all tests pass.
- `rtk go test ./...` — expected: full suite passes.

## References

`_prd.md` → User Story 6; Core Feature 11; User Experience. `_techspec.md` → System Architecture (tui), Build Order 10, Risks (mid-write reads). ADR-0009.

## Result

Spec Runs are now first-class in the Live Run View and Attach, through a small `WorkItem` view model both Run Kinds map into, keyed on `LiveRunView.RunKind`:

- **Work-item view model**: `internal/tui` gained `WorkItem` (Name, Title, Status) plus `TaskWorkItems`. The cockpit sidebar renders both Run Kinds through a shared `workItemBlock` (marker + name, then the status label); Review Issues map in via `issueWorkItems`, Tasks via `taskWorkItems`. Review rendering is byte-identical — zero review tests changed and all pass unchanged.
- **Cockpit Task pane**: implement Runs render `SPEC.TASKS` with Tasks in Task Graph order. Statuses refresh on every poll tick by re-reading task files through the Run row's git root and Spec slug (`spec.ReloadTask`); the cockpit reads the journal plus task files only (ADR-0009). Labels mirror the review pane: Completed/Failed verbatim, Executing for the Task in flight (first unsettled, or file status `in_progress`), Waiting behind it, Paused after the Run ends. The progress bar counts `N of M Task(s) completed`.
- **Mid-write tolerance**: a task-file parse failure during refresh keeps the last good status; the view never errors (techspec risk: mid-write reads).
- **Attach parity**: `attachRunView` populates the Task pane for implement Runs from the Run row's GitRoot + SpecSlug via `spec.Load`; load failures degrade to an empty pane, mirroring `attachIssues`. Replay-then-follow is unchanged and read-only.
- **Timeline**: `daemon.task` and `daemon.qa` render from their bounded summaries (both are daemon kinds), e.g. `Task task_01 started as Batch 001: <title>`, `Task task_01 settled completed.`, `QA verdict pass for Spec <slug>.` — now pinned by tests in the cockpit, the `RunTimeline`, and the CLI attach replay.
- **Plain-text parity**: for spec Runs `RenderLiveRunView` prints `Spec: <slug>` in the Target block (replacing the PR/Source lines) and a `Tasks` pane with one `task_NN <status> — <title>` line per Task, mirroring the implement stdout contract shape.

Commands run (all pass):

- `rtk go test ./internal/tui/ ./internal/cli/` — 164 tests pass.
- `rtk go test ./...` — 406 tests pass in 16 packages.
- `make verify` — fmt-check, full suite, `roundfix skills check`, build: all green.

Evidence per acceptance criterion:

1. `TestCockpitSpecRunShowsTasksInGraphOrderAndRefreshesStatuses` drives `Update(cockpitTickMsg{})` synchronously: Tasks render in graph order, and after rewriting a task file and appending a `daemon.task` journal event the pane flips to Completed and the timeline shows the settlement line.
2. `TestCockpitSpecRunKeepsLastGoodStatusOnMidWriteTaskFile` truncates the file mid-refresh: the previously shown Completed status and counts stay; no error in the view.
3. Review tests untouched (zero edits) and passing: cockpit and `RenderLiveRunView` review output unchanged.
4. `TestAttachReplaysCompletedSpecRunReadOnly` attaches to a finished implement Run: Task lines, `daemon.task`/`daemon.qa` replay, and the outcome line appear; Run count and terminal state untouched; Agent probe forbidden by the fake.
5. Full suite passes unchanged (`rtk go test ./...`, `make verify`).

Follow-up noted: the Enter detail pane stays Review-Issue-only; a read-only Task detail (task file body) is a candidate for a later slice.
