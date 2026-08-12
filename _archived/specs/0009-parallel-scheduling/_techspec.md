---
spec: 0009-parallel-scheduling
prd: _prd.md
created: 2026-07-05
---

# Parallel Scheduling — Technical Spec

## Executive Summary

Concurrency lands as a scheduler around unchanged per-Task machinery: the
same prompt, Verification, settlement, and commit path from 0001/0008 runs N
times in parallel, each inside a Task Worktree, with one new serialization
point — the integration queue onto the Run Branch (ADR-0026). The accepted
trade-off is **one more integration boundary inside the Run** (Task
Worktree → Run Branch, cherry-pick based, conflict → failed) in exchange for
wall-clock tracking graph depth instead of graph size. A second consequence
is deliberate: concurrent Tasks cannot share the Run's single Agent Session
(acpx serializes prompts per session), so each concurrent Task gets its own
session scoped to its Task Worktree — a scoped refinement of ADR-0018, with
the per-Run session retained for sequential mode and the QA step.

## System Architecture

- `internal/daemon` — `TaskCycle` gains the scheduler: a worker pool capped
  at `worktree.concurrency`, a guarded status map feeding ready-set
  recomputation on every settlement, Batch ordinals assigned at Task start
  under the scheduler lock, Stop Requests honored at scheduling boundaries
  (running Tasks settle; nothing new starts). Concurrency 1 runs the same
  code with one worker — behavior byte-compatible with today, proven by the
  existing suite.
- `internal/worktree` — Task Worktree lifecycle (`CreateTask`,
  `IntegrateTask`, `CleanupTask`) plus the serialized integration queue; the
  empty-debris reap extends `PruneTerminal`.
- `internal/config` — `worktree.concurrency` (int ≥ 1, default 2) and
  `worktree.location` (parent dir; `~` and absolute; default
  `~/.roundfix/worktrees`), both repo > global > builtin; the dead
  `resolve.concurrent` key is removed and its presence now fails validation
  with a named pointer to `worktree.concurrency`.
- `internal/agent` — nothing structural: session naming gains the per-Task
  variant; prompts/verification/settlement untouched.
- `internal/tui` + attach — concurrent Work Queue state driven by journal
  events (`daemon.task` started/settled) rather than file polling for
  executing Tasks; everything else as shipped.
- `internal/cli` — header line for effective concurrency; settle learns the
  deterministic Task Worktree path; no flag changes.
- Store — no schema change: Task Worktree paths are deterministic
  derivations, not persisted state.

## Implementation Design

### Paths and naming (deterministic, never configurable past the parent)

```
<location>/<repo-slug>/<run-id>/            # Run Worktree (0008, unchanged shape)
<location>/<repo-slug>/<run-id>.<task_id>/  # Task Worktree (sibling, never nested)
roundfix/run-<run-id>                       # Run Branch (unchanged)
roundfix/run-<run-id>-<task_id>             # Task Branch
roundfix-<run-id>-<task_id>                 # per-Task Agent Session name
```

`repo-slug` = sanitized repository directory basename + `-` + 8 hex chars of
the user-root path hash — readable and collision-free (0008 used the bare
hash; this replaces it). Task Worktrees are siblings of the Run Worktree so
no worktree ever nests inside another.

### Interfaces

```go
// internal/worktree
func CreateTask(ctx context.Context, run Ref, taskID string, copyList []string) (TaskRef, error)
    // worktree add -b roundfix/run-<id>-<task> <path> <runBranchTip>

type TaskIntegration struct {
    Mode   string // "ff" | "cherry-pick" | "conflict"
    Reason string // conflict detail when Mode == "conflict"
}
func IntegrateTask(ctx context.Context, run Ref, task TaskRef) (TaskIntegration, error)
    // serialized by the caller; ff when runBranch tip == task base,
    // else cherry-pick the task's commit(s); on conflict: cherry-pick
    // --abort, Run Branch left at its pre-attempt tip, Mode "conflict"
func CleanupTask(ctx context.Context, task TaskRef) error // remove worktree + task branch (success path)
```

### Scheduler (golang-concurrency discipline)

One owner goroutine runs the loop: compute ready set → start Tasks up to the
cap (each Task = one worker goroutine owning its whole per-Task pipeline:
Task Worktree create → prompt/Agent (per-Task session) → reload → verbatim
Verification inside the Task Worktree → settle status in the Task Worktree)
→ workers deliver results over a channel → the owner runs the integration
queue serially in completion order (cherry-pick semantics above; conflict
settles the Task failed and keeps its worktree/branch) → status map updates
→ recompute. Context cancellation drains workers (cooperative agent cancel
per session); every goroutine has an owner and a shutdown path; `-race` on
the package is part of every task gate. After the last settlement the QA
step (when opted in) runs exactly as today — alone, in the Run Worktree, on
the integrated Run Branch tip.

### Concurrency-correct surfaces

Batch ordinals are start-ordered and unique under the scheduler lock; every
journal event carries them as today, so Attach replay is deterministic by
cursor (unchanged property). The Work Queue derives executing/settled state
for concurrent Tasks from `daemon.task` events (files in Task Worktrees are
not polled); after integration, file reads from the Run Worktree behave as
shipped. stdout task lines keep graph order regardless of completion order
(the report already renders from final statuses). The Run header adds
`Concurrency: N` for spec Runs.

### Settle and debris

A failed Task keeps its Task Worktree; `roundfix settle --spec --task`
resolves the deterministic path first (falling back to the Run Worktree for
sequential-mode failures), verifies there, settles, and hands the commit to
the same integration queue mechanics (a settle-time conflict reports and
leaves the worktree kept). `PruneTerminal` (preflight sweep and
`stop --force`) additionally reaps worktrees and branches — Run or Task —
belonging to terminal Runs whose branch has no commits beyond its base
(`git merge-base --is-ancestor` equality check): provably empty debris.

## Coverage Map

- Story 1 → scheduler + Task Worktrees (ADR-0025)
- Story 2 → integration queue ordering + graph-order report (ADR-0026)
- Story 3 → conflict → failed with kept worktree (ADR-0026)
- Story 4 → journal-driven Work Queue concurrency
- Story 5 → `worktree.concurrency` config, sequential parity at 1
- Story 6 → `worktree.location` hierarchy with fixed slug/run segments
- Story 7 → empty-debris reap (round-3 finding 1)

## Integration Points

git only (worktree add/remove, cherry-pick, merge-base), through the
existing wrappers; acpx sessions per Task through the existing invocation
builders. No new external systems.

## Testing Approach

Hermetic temp repos throughout. `internal/worktree`: Task Worktree matrix —
ff first-finisher, cherry-pick second-finisher, induced conflict (two Tasks
editing one file) with Run Branch tip unmoved and worktree kept, cleanup,
empty-debris reap for Run and Task branches. `internal/daemon`: scheduler
tests with scripted fake runners — wide-Wave concurrency (observed overlap
via instrumented fakes), dependency gating under concurrency, failure
isolation, Stop Request at boundaries, ordinal uniqueness, and a
sequential-parity run at concurrency 1 against existing expectations;
`-race` mandatory. Config: hierarchy resolution, validation, the
`resolve.concurrent` rejection message. CLI/TUI: concurrent Work Queue
rendering from journal fixtures, header line, report order with shuffled
completion. The live dogfood metric: spec 0010's own Run executing its
independent Wave concurrently.

## Build Order

1. Config: `worktree.concurrency`, `worktree.location`, slug path scheme,
   `resolve.concurrent` removal (no deps)
2. Task Worktree lifecycle and serialized integration queue in
   `internal/worktree`, plus the empty-debris reap (no deps)
3. Scheduler in `TaskCycle`: worker pool, per-Task sessions, ordinals, stop
   boundaries, sequential parity (depends on: 1, 2)
4. Concurrency surfaces: journal-driven Work Queue, header, report order,
   attach determinism (depends on: 3)
5. Settle over Task Worktrees and debris-reap wiring in preflight and
   `stop --force` (depends on: 2, 3)
6. Docs and skill sync (depends on: 3, 4, 5)

## Risks & Considerations

- Concurrent `make verify` instances contend for CPU; default 2 is the
  deliberate ceiling and the docs say so. Go caches are concurrency-safe.
- Two Agent sessions double codex load; the per-Task session is also the
  natural acpx scope (sessions key on cwd) — no shared-session queueing.
- Cherry-pick preserves the Task commit message/trailers; tests byte-compare
  them post-integration.
- The scheduler must never deadlock on a failed dependency chain: dependents
  of failed Tasks leave the ready set permanently and are reported skipped,
  as today.
- Location config moves worktrees; the debris reap and settle must resolve
  paths through the same derivation helper — one source of truth.

## Decisions

- Spec Runs only; concurrency cap `worktree.concurrency` default 2; parity
  at 1. See ADR-0025.
- Serialized cherry-pick integration; conflict → failed. See ADR-0026.
- Per-Task Agent Sessions for concurrent Tasks — a scoped refinement of
  ADR-0018 (the per-Run session remains for sequential mode and QA);
  recorded here rather than a new ADR since ADR-0025 already states it.
- Deterministic sibling paths; repo segment becomes slug+hash; no store
  schema change.
- `resolve.concurrent` removed with a pointing validation error.
