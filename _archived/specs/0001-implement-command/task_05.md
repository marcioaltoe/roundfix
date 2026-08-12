---
task: task_05
spec: 0001-implement-command
status: completed
type: backend
complexity: high
---

# Task 05: Execute the Task Graph through the Daemon engine

## Overview

Add `TaskCycle` to the Daemon engine as a sibling of the resolve cycle: it walks the Task Graph in topological order, runs one Agent per Task as a Batch of one, runs the Task's Verification verbatim, settles status, commits per Task, and applies the failure policy. This is the core of the feature and is verifiable on its own by driving the cycle with fake collaborators over a real temporary git repository and spec directory.

## Requirements

1. MUST execute non-completed Tasks in topological order, one at a time; a Task runs only when every `needs` entry is `completed` in a live status map seeded from task files and updated as the Run progresses — otherwise it is skipped and left `pending`.
2. MUST run each Task as a Batch of one through the existing Agent runner with the task prompt; the Batch number is the 1-based execution ordinal and Run Events carry the Task id in the Work Item field.
3. MUST take the before-snapshot at Task start so pre-existing worktree changes never enter the Task commit.
4. MUST re-read the task file after the Agent finishes, then run every Verification command sequentially through the existing verifier; all must pass.
5. MUST settle status per ADR-0014: passing verification settles `completed` (writing the status if the Agent forgot); an Agent error, an Agent-reported `failed`, or a verification failure settles `failed` with the reason journaled — and `completed` is never settled without passing verification.
6. MUST commit code changes plus the updated task file on success, with the message `<type>: <task title>` mapping `docs→docs`, `test→test`, `chore→chore`, everything else `→feat`, plus `Roundfix-Spec` and `Roundfix-Task` trailers; a failed Task produces no commit and its worktree changes are preserved.
7. MUST continue with independent Tasks after a failure; only Stop Requests and infrastructure errors halt the cycle (generalizing ADR-0010).
8. MUST add the `daemon.task` and `daemon.qa` Run Event Kinds and emit `daemon.task` phase events plus the existing verification, commit, and outcome events; never wire the pusher or the Review Source resolver.

## Subtasks

- [x] `TaskPlan`/`TaskCycleResult` types and the cycle skeleton with needs-gating
- [x] Per-Task Agent invocation with Batch-of-one identity and journaling
- [x] Verbatim per-Task verification through the existing verifier
- [x] Status settling for completed, forgotten, and failed outcomes
- [x] Task commit with type mapping and trailers
- [x] Failure policy: dependents stay pending, independents continue, stop/infra halts
- [x] New Run Event Kinds and their emission

## Acceptance Criteria

- [x] Happy path: a two-Task chain produces two commits, each containing only that Task's changes plus its task file, with correct messages and trailers for at least the `docs` and default `feat` mappings.
- [x] A failed middle Task produces no commit, its dependents end the cycle still `pending`, an independent Task still completes, and the result counts completed/failed/skipped correctly.
- [x] An Agent that leaves status untouched is settled by the Daemon: `completed` on passing verification, `failed` on failing verification.
- [x] A Task marked `completed` by the Agent whose verification fails is settled `failed`.
- [x] Pre-existing dirty files from before a Task are absent from that Task's commit.
- [x] A canceled context halts the cycle with no further Tasks executed.
- [x] The full existing review-path suite passes unchanged.

## Verification

- `rtk go test ./internal/daemon/ ./internal/runevent/` — expected: all tests pass.
- `rtk go test -race ./internal/daemon/` — expected: no races.
- `rtk go test ./...` — expected: full suite passes.

## References

`_prd.md` → User Stories 1, 2; Core Features 3, 5, 6, 7. `_techspec.md` → Interfaces (TaskPlan/TaskCycle), Data Models (Run Events), Build Order 3, 6. ADR-0010, ADR-0013, ADR-0014.

## Result

`Engine.TaskCycle(ctx, TaskPlan) (TaskCycleResult, error)` now executes a Spec's Task Graph as a sibling of the resolve cycle (`internal/daemon/task_engine.go`). Non-completed Tasks run sequentially in the given topological order behind a live needs gate seeded from task-file statuses (prior-Run `completed` counts; skipped Tasks stay `pending` on disk and journal a `daemon.task` skip event with the unmet needs). Stale `in_progress` and `failed` Tasks re-run fresh. Each executed Task: before-snapshot → `ResolvingWithAgent` → fresh task-file read → `BuildTaskPrompt` → `Runner.Run` as a Batch of one (Batch number = 1-based execution ordinal, Task id in the event Work Item field, log path `<WorkDir>/.roundfix/runs/<run-id>/agent/batch-NNN.log`) → `ReloadTask` → `Verifying` → every Verification command verbatim through the existing `Verifier` in WorkDir, fail-fast, with `defaults.verification` never appended → settle per ADR-0014 (`completed` only after passing verification, writing the status when the Agent forgot; Agent error, Agent-set `failed`, or verification failure settles `failed` with the reason journaled) → on success, one commit of the snapshot diff plus the task file with `<type>: <title>` (docs/test/chore pass through, everything else feat) and `Roundfix-Spec`/`Roundfix-Task` trailers after a blank line. Failed Tasks create no commit, preserve worktree changes, and the cycle continues with independent Tasks; only Stop Requests and infrastructure errors halt (ADR-0010 semantics via the existing `isStop`). `internal/runevent` gains `KindDaemonTask` (`daemon.task`) and `KindDaemonQA` (`daemon.qa`), both covered by `IsDaemonKind`; the cycle ends with a `daemon.outcome` event carrying the completed/failed/skipped counts. Pusher and Review Source resolver are never invoked. The QA step is a marked seam in `TaskCycle` (comment between the Task walk and the outcome event) for task_07; `TaskCycleResult.QAVerdict` stays empty.

Design note for task_06: `TaskPlan` matches the techspec exactly (no Artifact Directory field), so the Agent log path derives from the documented `.roundfix` default under `WorkDir`; a configured `defaults.artifact_dir` does not redirect spec-Run logs.

Verification evidence (all fresh, after the last edit):

- `rtk go test ./internal/daemon/ ./internal/runevent/` — 42 tests passed.
- `rtk go test -race ./internal/daemon/` — 29 tests passed, no races.
- `rtk go test ./...` — 351 tests passed in 16 packages (review-path suite unchanged).
- `make verify` — fmt-check, tests, `roundfix skills check`, and build all passed.

Evidence per acceptance criterion (`internal/daemon/task_engine_test.go`): happy path with real git repo, per-Task commit contents, docs/feat messages and trailers — `TestTaskCycleRealRepoCommitsPerTaskExcludingPreexistingDirt` and `TestTaskCycleExecutesAgentVerifySettleCommitContract`; failed middle Task, pending dependents, continuing independents, correct counts — `TestTaskCycleFailedTaskSkipsDependentsAndContinuesIndependents`; forgotten-status settling both ways — `TestTaskCycleSettlesForgottenAgentStatus`; Agent-claimed `completed` with failing verification settled `failed` — asserted on `task_01` in `TestTaskCycleFailedTaskSkipsDependentsAndContinuesIndependents`; pre-existing dirty files excluded — real-repo test plus `TestTaskCycleCommitStagesSnapshotDiffPlusTaskFile`; canceled context halts — `TestTaskCycleStopBeforeTaskPublishesStopAndDoesNothing` and `TestTaskCycleStopDuringAgentPreservesTaskAndHalts`; review-path suite unchanged — full `go test ./...` green with zero edits to existing tests.
