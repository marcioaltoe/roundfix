---
task: task_04
spec: 0009-parallel-scheduling
status: completed
type: frontend
complexity: medium
---

# Task 04: Concurrency-correct surfaces

## Overview

Make the visible surfaces truthful under parallelism: the Work Queue shows
every executing Task simultaneously (driven by journal events, not Task
Worktree file polling), the Run header names the effective concurrency, the
stdout report keeps graph order regardless of completion order, and Attach
replays interleaved history deterministically. Verifiable through
synchronous cockpit tests over interleaved journal fixtures.

## Requirements

1. MUST derive executing/settled state for concurrent Tasks from
   `daemon.task` journal events; Task Worktree files are never polled by
   the cockpit; post-integration file reads from the Run Worktree behave
   as shipped (mid-write tolerance unchanged).
2. MUST render multiple simultaneous `Executing` rows in the Work Queue
   with stable ordering (graph order), and keep the totals footer correct
   while Tasks settle out of order.
3. MUST add `Concurrency: N` to the spec-Run header (live and plain
   renderers), always shown for spec Runs with the effective resolved
   value; review-Run headers are untouched.
4. MUST keep stdout task lines in graph order under shuffled completion
   (fixture with reversed completion order) and Attach replay
   deterministic over interleaved events for both live-written and
   replayed journals.
5. MUST leave review-Run rendering byte-stable.

## Subtasks

- [x] Journal-driven Work Queue state for concurrent Tasks
- [x] Multiple-executing rendering with stable order and totals
- [x] Header concurrency line in both renderers
- [x] Graph-order report and attach determinism fixtures

## Acceptance Criteria

- [x] An interleaved journal fixture (two Tasks starting before either
      settles) renders two `Executing` rows in graph order; settlements
      out of order update rows and totals correctly.
- [x] Reversed-completion fixture: stdout lines in graph order, byte-exact.
- [x] Attach over the same fixture renders identically to the live pass.
- [x] Review snapshots unchanged; full suite passes.

## Verification

- `rtk go test ./internal/tui/ ./internal/cli/` — expected: all tests pass.
- `rtk go test ./...` — expected: full suite passes.

## References

`_prd.md` → User Stories 2, 4; Core Feature 6. `_techspec.md` →
Concurrency-correct surfaces, Build Order 4. ADR-0009.

## Result

- Interleaved Task journal state: `TestCockpitSpecRunDerivesConcurrentTaskStateFromJournal` covers two `daemon.task` starts before either settles, graph-order Work Queue rows, out-of-order settlement, and totals updates.
- Reversed completion stdout ordering: `TestRenderImplementTaskLinesKeepsGraphOrderWhenCompletionReversed` asserts byte-exact graph-order task lines after statuses are settled in reverse order.
- Attach determinism: `TestCockpitSpecRunInterleavedTaskReplayMatchesLivePolling` verifies live polling and Attach-style replay render identically over the same interleaved event stream.
- Review stability and full verification: `TestCockpitRenderSnapshots` kept review snapshots unchanged; `rtk go test ./internal/tui/ ./internal/cli/`, `rtk go test ./...`, and `rtk make verify` all passed.
