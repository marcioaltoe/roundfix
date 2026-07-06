---
spec: 0014-run-store-retention
status: active
created: 2026-07-06
surfaces: [cli, infra, docs]
---

# Run Store Retention

The Run Database's Run Event Journal is append-only and never pruned, so
`run_events` grows without bound — a dogfood database reached ~220 MB, almost
entirely journal rows carrying every agent payload of every Run, forever. The
redundant per-Batch agent log files are already going opt-in (spec 0011), but
the journal itself and its on-disk artifact directories keep accumulating. This
Spec bounds the Run store: terminal Runs older than a configured retention
window have their journal and artifact directory pruned — on demand through a
GC Command and best-effort during the preflight sweep — while the load-bearing
Run state (rows and active-run locks) is never touched.

## Goals

- The Run Event Journal stops growing without bound: terminal Runs past a
  configured retention window become eligible for pruning. See ADR-0033.
- A developer can reclaim Run storage on demand and see what was freed.
- Pruning is safe by construction: Active Runs, `runs` rows, and active-run
  locks are never removed by retention — only journal events and artifact
  directories of terminal Runs.
- Storage self-heals: the operational preflight sweep prunes eligible Runs
  best-effort, so the database does not require manual maintenance.

## User Stories

1. As a developer dogfooding Roundfix, I want the Run Event Journal to stop
   accumulating unbounded rows, so that my Run Database does not grow to
   hundreds of megabytes of stale history.
2. As a developer, I want a command that reclaims Run storage and reports what
   it freed, so that I can clean up on demand and trust what happened.
3. As a developer, I want a dry-run that shows what pruning would remove without
   removing it, so that I can review before reclaiming.
4. As a developer, I want retention to never touch Active Runs or the state that
   powers stop/attach/detach/one-active-run, so that cleanup can never break a
   live or recoverable Run. See ADR-0033.
5. As a developer who never runs cleanup, I want the preflight sweep to prune
   eligible terminal Runs automatically, so that storage self-heals without a
   manual step.

## Core Features

1. **Journal Retention window.** A config key sets the age after which a
   terminal Run's Run Event Journal and artifact directory are eligible for
   pruning; `0` disables pruning (keep everything). Active Runs are never
   eligible. See ADR-0033.
2. **GC Command.** `roundfix gc` prunes eligible terminal Runs' journal events
   and artifact directories, removes orphaned run artifact directories (no
   matching Run row), and reports the rows and bytes freed. A `--dry-run` lists
   what would be pruned without changing anything.
3. **Retention-safe scope.** Pruning removes only Run Event Journal rows and
   on-disk artifact directories of terminal Runs past the window; it never
   deletes `runs` rows, active-run locks, or anything belonging to an Active
   Run.
4. **Self-healing prune.** The operational preflight sweep runs a best-effort
   retention prune (when the window is non-zero), so storage is bounded without
   requiring the GC Command; failures are warnings, never fatal.

## User Experience

`roundfix gc` prints a deterministic report of what it freed (Runs pruned,
journal rows removed, artifact bytes reclaimed) and `--dry-run` prints the same
shape prefixed as a preview. The preflight sweep gains one optional summary line
when it prunes. No other output changes; retention defaults to a bounded window
so a fresh install is already self-limiting, and setting it to `0` restores
unlimited history.

## Non-Goals / Out of Scope

- Deleting or archiving `runs` rows themselves — retention prunes journals and
  artifacts, not the Run index (the compact row stays as durable history).
- Compacting or vacuuming the SQLite file itself beyond row deletion (a
  follow-up if disk reclamation needs it).
- Changing the Run Event Journal format, cursor semantics, or replay (ADR-0008).
- Retention for review artifacts under the spec tree (spec 0011 owns their
  placement; they are versioned by the repository owner, not GC'd here).
- Making the journal itself opt-in per Run (retention bounds it; opt-in is a
  separate future decision).

## Success Metrics

- After `roundfix gc` on a database with terminal Runs older than the window,
  their journal rows and artifact directories are gone and the report names the
  freed counts; Active Runs and `runs` rows remain.
- `roundfix gc --dry-run` changes nothing and lists the same set it would prune.
- With a non-zero window, running an operational command prunes eligible
  terminal Runs during preflight without a manual step.
- A retention of `0` prunes nothing, on demand or in the sweep.

## Decisions

- Retention bounds the Run Event Journal and artifact directories of terminal
  Runs by an age window; `runs` rows and locks are never pruned. See ADR-0033.
- Cleanup is both on demand (GC Command, with `--dry-run`) and best-effort in
  the preflight sweep; the glossary gains GC Command and Journal Retention.
- Default retention is a bounded window (not unlimited and not aggressive), with
  `0` meaning keep everything — the exact default is set in the TechSpec.

## Open Questions

None.
