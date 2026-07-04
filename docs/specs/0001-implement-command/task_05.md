---
task: task_05
spec: 0001-implement-command
status: pending
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

- [ ] `TaskPlan`/`TaskCycleResult` types and the cycle skeleton with needs-gating
- [ ] Per-Task Agent invocation with Batch-of-one identity and journaling
- [ ] Verbatim per-Task verification through the existing verifier
- [ ] Status settling for completed, forgotten, and failed outcomes
- [ ] Task commit with type mapping and trailers
- [ ] Failure policy: dependents stay pending, independents continue, stop/infra halts
- [ ] New Run Event Kinds and their emission

## Acceptance Criteria

- [ ] Happy path: a two-Task chain produces two commits, each containing only that Task's changes plus its task file, with correct messages and trailers for at least the `docs` and default `feat` mappings.
- [ ] A failed middle Task produces no commit, its dependents end the cycle still `pending`, an independent Task still completes, and the result counts completed/failed/skipped correctly.
- [ ] An Agent that leaves status untouched is settled by the Daemon: `completed` on passing verification, `failed` on failing verification.
- [ ] A Task marked `completed` by the Agent whose verification fails is settled `failed`.
- [ ] Pre-existing dirty files from before a Task are absent from that Task's commit.
- [ ] A canceled context halts the cycle with no further Tasks executed.
- [ ] The full existing review-path suite passes unchanged.

## Verification

- `rtk go test ./internal/daemon/ ./internal/runevent/` — expected: all tests pass.
- `rtk go test -race ./internal/daemon/` — expected: no races.
- `rtk go test ./...` — expected: full suite passes.

## References

`_prd.md` → User Stories 1, 2; Core Features 3, 5, 6, 7. `_techspec.md` → Interfaces (TaskPlan/TaskCycle), Data Models (Run Events), Build Order 3, 6. ADR-0010, ADR-0013, ADR-0014.
