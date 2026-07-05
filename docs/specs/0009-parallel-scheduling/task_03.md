---
task: task_03
spec: 0009-parallel-scheduling
status: pending
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

- [ ] Owner loop with ready-set recomputation and status map
- [ ] Worker pipeline with per-Task sessions and worktrees
- [ ] Serialized integration wiring with conflict → failed settlement
- [ ] Ordinals, journaling, and log paths under concurrency
- [ ] Stop/cancellation at boundaries with worker drain
- [ ] Sequential-parity and `-race` gates

## Acceptance Criteria

- [ ] A 4-Task independent Wave at concurrency 2 shows observed execution
      overlap (instrumented fakes), 4 commits on the Run Branch, correct
      counts, and deterministic journal replay.
- [ ] Dependency gating holds under concurrency: a dependent never starts
      before its needs settle; failed-dependency chains end skipped.
- [ ] An induced integration conflict settles that Task failed (worktree
      kept, reason journaled) while independents complete.
- [ ] A Stop Request mid-Wave lets running Tasks settle and commit, starts
      nothing new, ends Stopped.
- [ ] Concurrency 1 passes the existing TaskCycle expectations unchanged;
      `rtk go test -race ./internal/daemon/` is clean.

## Verification

- `rtk go test ./internal/daemon/` — expected: all tests pass.
- `rtk go test -race ./internal/daemon/` — expected: no races.
- `rtk go test ./...` — expected: full suite passes.

## References

`_prd.md` → User Stories 1, 3, 5; Core Features 1–3. `_techspec.md` →
Scheduler, Build Order 3, Risks. ADR-0010, ADR-0014, ADR-0018 (scoped
refinement), ADR-0025, ADR-0026.
