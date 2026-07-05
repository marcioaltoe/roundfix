---
task: task_03
spec: 0009-parallel-scheduling
status: completed
type: backend
complexity: high
---

# Task 03: Ready-set scheduler in TaskCycle

## Overview

The core: `TaskCycle` becomes a ready-set scheduler running up to
`worktree.concurrency` Tasks simultaneously, each worker owning its whole
per-Task pipeline inside a Task Worktree with a per-Task Agent Session,
results integrating through the serialized queue, and concurrency 1
reproducing today's behavior exactly. Verifiable through instrumented
fake-runner tests under `-race`.

## Requirements

1. MUST implement the owner-loop scheduler: continuously compute the ready
   set (needs completed via prior Runs or this Run's settlements), start
   Tasks up to the cap, receive settlements over a channel, run the
   integration queue serially in completion order, update the status map,
   recompute; dependents of failed Tasks leave the ready set permanently
   and are reported skipped, exactly as today.
2. MUST give each worker goroutine ownership of its full pipeline: Task
   Worktree creation → Agent prompt through a per-Task Agent Session
   (`roundfix-<run-id>-<task_id>`, cwd = the Task Worktree, closed at
   settlement) → reload → verbatim Verification inside the Task Worktree →
   status settlement (ADR-0014 unchanged) → commit in the Task Worktree.
   Integration conflicts settle the Task failed with its worktree kept
   (ADR-0026).
3. MUST assign Batch ordinals at Task start under the scheduler lock
   (unique, monotonic); all journal events carry them as today; agent logs
   remain per-Batch under the Artifact Directory.
4. MUST honor Stop Requests and context cancellation at scheduling
   boundaries: nothing new starts, running workers finish their settlement
   (cooperative agent cancel on hard cancellation), then the Run ends
   through the existing Stopped path.
5. MUST preserve sequential semantics at concurrency 1 through the same
   code path (single worker): the existing TaskCycle test expectations pass
   without behavioral edits.
6. MUST keep the QA step exactly as shipped: after the last settlement,
   alone, in the Run Worktree, on the integrated Run Branch tip. Every
   goroutine has an owner, cancellation, and a shutdown path; the package
   gate includes `-race`.

## Subtasks

- [x] Owner loop with ready-set recomputation and status map
- [x] Worker pipeline with per-Task sessions and worktrees
- [x] Serialized integration wiring with conflict → failed settlement
- [x] Ordinals, journaling, and log paths under concurrency
- [x] Stop/cancellation at boundaries with worker drain
- [x] Sequential-parity and `-race` gates

## Acceptance Criteria

- [x] A 4-Task independent Wave at concurrency 2 shows observed execution
      overlap (instrumented fakes), 4 commits on the Run Branch, correct
      counts, and deterministic journal replay.
- [x] Dependency gating holds under concurrency: a dependent never starts
      before its needs settle; failed-dependency chains end skipped.
- [x] An induced integration conflict settles that Task failed (worktree
      kept, reason journaled) while independents complete.
- [x] A Stop Request mid-Wave lets running Tasks settle and commit, starts
      nothing new, ends Stopped.
- [x] Concurrency 1 passes the existing TaskCycle expectations unchanged;
      `rtk go test -race ./internal/daemon/` is clean.

## Verification

- `rtk go test ./internal/daemon/` — expected: all tests pass.
- `rtk go test -race ./internal/daemon/` — expected: no races.
- `rtk go test ./...` — expected: full suite passes.

## References

`_prd.md` → User Stories 1, 3, 5; Core Features 1–3. `_techspec.md` →
Scheduler, Build Order 3, Risks. ADR-0010, ADR-0014, ADR-0018 (scoped
refinement), ADR-0025, ADR-0026.

## Result

- Implemented the TaskCycle owner-loop scheduler in `internal/daemon`, with
  ready-set recomputation, a concurrency cap, worker settlement collection,
  serialized Task Worktree integration, failed-dependency skipping, and
  stop-boundary draining.
- Added per-Task worker execution in Task Worktrees for parallel-ready graphs:
  Task Worktree creation, per-Task Agent Session naming
  `roundfix-<run-id>-<task_id>`, verification in the Task Worktree, ADR-0014
  settlement, Task Worktree commit, serialized integration, conflict failure
  handling, and success cleanup.
- Wired `worktree.concurrency` and copy-list provisioning from the CLI
  implement path into `TaskPlan`, while preserving concurrency-1 sequential
  behavior through the scheduler path.
- Added scheduler matrix tests:
  `TestTaskCycleSchedulesIndependentWaveWithConcurrencyCap`,
  `TestTaskCycleGatesDependenciesAndSkipsFailedDependencyChainsUnderConcurrency`,
  `TestTaskCycleIntegrationConflictSettlesTaskFailedAndKeepsTaskWorktree`, and
  `TestTaskCycleStopRequestMidWaveDrainsRunningTasksAndStartsNothingNew`.

Acceptance evidence:

- 4-Task independent Wave at concurrency 2: covered by
  `TestTaskCycleSchedulesIndependentWaveWithConcurrencyCap`, which observes
  worker overlap at cap 2, then verifies four integrations, four commits,
  completed counts, and deterministic ordinal journal events.
- Dependency gating and failed-dependency skips: covered by
  `TestTaskCycleGatesDependenciesAndSkipsFailedDependencyChainsUnderConcurrency`,
  which asserts a dependent does not start before its need integrates and that
  failed-dependency chains settle skipped.
- Integration conflict behavior: covered by
  `TestTaskCycleIntegrationConflictSettlesTaskFailedAndKeepsTaskWorktree`,
  which returns a conflict result, settles the Task failed, journals the
  reason, keeps the Task Worktree, and allows independent Tasks to complete.
- Stop Request mid-Wave: covered by
  `TestTaskCycleStopRequestMidWaveDrainsRunningTasksAndStartsNothingNew`,
  which lets started Tasks settle and integrate, prevents new starts, and ends
  through the stopped path.
- Concurrency 1 parity and race cleanliness: existing TaskCycle expectations
  still pass, `TestRunImplementUsesOneAgentSessionPerRunAndCloses` pins
  concurrency 1 for the ADR-0018 path, and
  `rtk go test -race ./internal/daemon/` passed.

Verification evidence:

- `rtk go test ./internal/daemon/` - passed, 57 tests.
- `rtk go test -race ./internal/daemon/` - passed, 57 tests, no races.
- `rtk go test ./...` - passed, 702 tests across 17 packages.
- `rtk make verify` - passed, including full Go tests, Roundfix skill check,
  and build.
