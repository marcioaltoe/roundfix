---
spec: 0009-parallel-scheduling
status: active
created: 2026-07-05
surfaces: [cli, infra]
---

# Parallel Scheduling

Spec Runs execute one Task at a time even when the Task Graph declares whole
Waves of independent work — the dogfood cycles spent hours in sequential
waits that the graph itself said were unnecessary. The 0008 worktree
isolation was built precisely to host concurrency; this Spec adds the layer
on top: a ready-set scheduler that runs independent Tasks simultaneously,
each in its own Task Worktree, integrating settled work back onto the Run
Branch through a serialized queue. It also delivers the worktree location
configuration and the empty-debris cleanup the last rounds asked for.

## Goals

- Independent Tasks of a Wave execute concurrently up to a configured cap,
  cutting wall-clock for wide graphs without changing any per-Task
  ownership, verification, or commit semantics. See ADR-0025.
- Settled Tasks land on the Run Branch deterministically through a
  serialized integration queue; conflicts between supposedly independent
  Tasks surface as failures, never silent merges. See ADR-0026.
- The worktree parent location becomes configurable (repo > global >
  builtin default), with the repo-slug and unique-run segments always
  appended and never configurable.
- Terminal Runs that produced no commits stop leaving worktrees and Run
  Branches behind.
- Sequential behavior remains available and byte-compatible at
  concurrency 1.

## User Stories

1. As a developer implementing a Spec with a wide Wave, I want independent
   Tasks running concurrently, so that the Run's wall-clock tracks the
   graph's depth instead of its size.
2. As a developer reading the Run report, I want per-Task lines in graph
   order and one commit per Task on the Run Branch regardless of completion
   order, so that concurrency never makes outcomes harder to read.
3. As a developer whose "independent" Tasks actually collided, I want the
   conflicting Task settled failed with its Task Worktree kept and the
   conflict named, so that a graph defect surfaces instead of a mystery
   merge.
4. As a developer watching the cockpit, I want every executing Task visible
   simultaneously with its own state, so that the Live Run View reflects
   real concurrency.
5. As a developer with resource limits, I want a concurrency config with a
   conservative default and a sequential opt-out, so that heavy Agent
   sessions never overwhelm my machine.
6. As a developer, I want the worktree parent directory configurable per
   repo or globally, so that worktrees live where my disk layout wants them
   — while the per-repo slug and unique run id segments stay fixed to
   prevent collisions.
7. As a developer with force-stopped or crashed Runs, I want provably empty
   kept worktrees and Run Branches reaped automatically, so that debris
   stops accumulating.

## Core Features

1. **Ready-set scheduler.** The Implement Command computes the current Wave
   continuously (dependencies satisfied by prior Runs or settlements in this
   Run) and executes up to `worktree.concurrency` Tasks at once (default 2;
   1 = today's sequential behavior exactly). Failure policy is unchanged per
   Task: a failed Task blocks only its dependents; independents continue.
   Stop Requests are honored at scheduling boundaries — running Tasks settle
   first, nothing new starts. See ADR-0025.
2. **Task Worktrees.** Each concurrently executed Task runs in a Task
   Worktree created from the Run Branch tip at its start, with its own
   Agent prompt, verbatim Verification, and settlement inside it; a failed
   Task's worktree is kept as its inspection and settle surface, a
   successful one is removed after integration.
3. **Serialized integration queue.** Settled Tasks integrate onto the Run
   Branch in completion order — fast-forward when possible, cherry-pick
   otherwise; a conflict settles that Task failed with the conflict named.
   The Run Branch → user branch integration is untouched (ADR-0024). The QA
   step still runs last, alone, in the Run Worktree. See ADR-0026.
4. **Worktree location config.** `worktree.location` sets the parent
   directory with repo > global > builtin hierarchy (builtin:
   `~/.roundfix/worktrees/`); the final path is always
   `<location>/<repo-slug>/<run-id>/` — the readable repo slug and the
   unique run id are appended unconditionally and are not configurable.
5. **Empty-debris cleanup.** The preflight sweep and `stop --force` also
   reap kept worktrees and Run Branches of terminal Runs whose branch
   carries no commits beyond its base — provably nothing to lose.
6. **Concurrency-correct surfaces.** Journal events carry per-Task Batch
   ordinals assigned at start; the Work Queue shows every executing Task;
   stdout task lines stay in graph order; Attach replays interleaved
   history deterministically.

## User Experience

Same commands and flags. New config keys (`worktree.concurrency`,
`worktree.location`); the Run header names the effective concurrency; the
cockpit shows multiple `Executing` rows during wide Waves; reports read
exactly as today. At concurrency 1 nothing observable changes.

## Non-Goals / Out of Scope

- Review-path (resolve/watch) parallelism — Batches declare no independence;
  a future spec with conflict-aware planning owns that.
- Retry budgets, escalation, or scheduling priorities (work-plan item 7).
- Cross-Run parallelism (the one-Active-Run-per-target lock stands).
- Merge-conflict resolution — conflicts settle failed by design.
- Distributed or remote execution.

## Success Metrics

- A Spec with a 4-Task independent Wave at concurrency 2 completes in
  roughly half the sequential wall-clock in the live dogfood, with one
  commit per Task on the Run Branch and a byte-identical report shape.
- The next Spec's own implement Run (0010) executes its independent Wave
  concurrently — the standing "next spec tests the previous" convention.
- An induced collision fixture settles the conflicting Task failed, keeps
  its Task Worktree, and completes independents.
- Concurrency 1 reproduces today's behavior against the existing test
  suite.
- After a force-stop with zero settled work, the sweep leaves no worktree
  and no Run Branch behind.

## Decisions

- Spec Runs only; review path deferred (Marcio, 2026-07-05). See ADR-0025.
- `worktree.concurrency`, default 2, repo > global > builtin hierarchy;
  the legacy unused `resolve.concurrent` key is removed from generated
  config (breaking only for a key that never did anything — validation now
  rejects it with a pointer to the new key).
- Serialized cherry-pick queue, conflict → failed. See ADR-0026.
- `worktree.location` parent-only config; slug and run id segments fixed
  (Marcio, 2026-07-05). The repo segment becomes the readable slug (was a
  hash in 0008).
- Glossary gains **Wave** and **Task Worktree**.

## Open Questions

None.
