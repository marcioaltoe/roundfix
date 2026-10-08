---
status: accepted
created_at: 2026-10-08T00:00:00Z
updated_at: 2026-10-08T00:00:00Z
deprecated_at: null
superseded_by: null
---

# Run Retention removes terminal Runs whole and compacts incrementally

On 2026-10-08 the machine-wide Run Database was 1.1 GB. `run_events` held
1.08 GB for 832k events written between 2026-09-24 and 2026-10-07, about
77 MB a day, across 313 Runs. Journal Retention (ADR-0033) prunes only the
journal and artifact directory of a terminal Run and never deletes its `runs`
row, so the row index and its dependent rows grow without bound, and nothing
returns freed pages to the filesystem unless someone runs
`roundfix gc compact --apply`. The maintainer chose "Automática + gc
(Recommended)": keep only the data of Runs still running, and of any Run a
Delivery Queue still references, plus the last N days, with N one of 30, 15
or 7 and 30 by default, applied automatically at most once a day and through
`roundfix gc`, then compact the database.

## Decision

- **Window.** User Config key `store.run_retention_days` accepts exactly `7`,
  `15` or `30` and defaults to `30`. Any other value is a configuration
  error that names the three values. There is no value that turns it off.
  The Run Database is machine-wide, so the key is User Config only: a Project
  Config value is ignored with the warning every User-Config-only key prints.
- **Journal Retention stays.** `store.journal_retention` keeps its meaning,
  its scopes and its `336h` default. It is the inner window that prunes the
  bulky journal early; Run Retention is the outer window that removes the Run
  itself. A Journal Retention longer than the Run Retention has no effect past
  the Run Retention, because the Run and its journal leave together. The key
  is neither deprecated nor mapped.
- **What is removed.** A terminal Run whose `completed_at` is older than the
  window is removed whole in one write transaction per Run: its `runs` row,
  and through the existing `ON DELETE CASCADE` foreign keys its `run_events`,
  `run_agent_selections`, `run_token_usage` and `active_run_locks` rows. Its
  `runs/<run-id>` directory is removed under the Run's Artifact Root when that
  root is proven the way `roundfix gc sanitize` proves one. No other table
  references a Run by foreign key: `run_windows`, `interactive_defaults` and
  the Delivery Queue tables are never touched.
- **What is kept.** A Run is never removed while it is not terminal, while
  any Delivery Queue item or Delivery Queue Run link names it, while its
  recorded Run Worktree still exists on disk (reconcile and Task Carry-Forward
  read that Run), or while its artifact directory exists under a root that
  cannot be proven. Each removal transaction re-reads the state and the queue
  references before it deletes, so a Run that changed after the scan is kept.
- **Once a day.** A one-row table `run_retention_sweeps` in the Run Database
  records when the last sweep completed and with which window. The automatic
  sweep is due when no sweep has completed, when the last one completed more
  than 24 hours ago, or when the configured window differs from the recorded
  one. Two processes that start at once may both sweep; every step is
  idempotent, so the overlap costs time and never correctness.
- **Never blocking a start for long.** The automatic sweep runs at the start
  of every Run (`implement`, `resolve`, `watch`) and of every Delivery Queue
  (`deliver start`), after the existing Journal Retention prune, under a
  two-second budget. It takes the machine-wide write lock only per Run
  removal and per compaction slice, so a running Daemon waits at most one such
  step. When the budget runs out it stops between steps, records no
  completion and continues at the next start. A failure prints one warning and
  never blocks the Run. `roundfix gc` runs the same sweep without a budget.
- **Compaction.** A new Run Database is created with
  `auto_vacuum = INCREMENTAL`. In that mode the sweep returns freed pages with
  `PRAGMA incremental_vacuum` in slices of 2,560 pages, each its own write
  transaction. An existing database in the default mode stays in it until a
  guarded full compaction converts it: `roundfix gc compact --apply` and the
  live `roundfix gc` now set `auto_vacuum = INCREMENTAL` before their
  `VACUUM`, and the preview measures the same layout. The automatic sweep
  never runs a full `VACUUM`. When the live `roundfix gc` cannot compact, it
  reports why and still exits 0.
- **Surfaces.** `roundfix gc --dry-run` reports the window, its cutoff, the
  Runs it would remove, the Runs it keeps by reason, the rows per table and an
  estimate of the bytes. The live report adds the measured file bytes before
  and after. `roundfix runs show` and `roundfix events` keep refusing an
  unknown Run, and when the Run ID's embedded creation time predates the
  cutoff they add that Run Retention may have removed it. `roundfix runs list`
  shows no notice, as ADR-0172 keeps for its hot path. `deliver status` and
  `roundfix reconcile` are unchanged, because the Runs they read are kept.

## Measurements

Measured on 2026-10-08 on a `.backup` copy of the live database, through the
same `modernc.org/sqlite` driver, on the maintainer's machine:

- a 30-day window removes 12 Runs and 0 events, a 15-day window 68 Runs and
  0 events, because Journal Retention already pruned their journals; a 7-day
  window removes 200 Runs and 439,456 events;
- the 7-day removal took 4.3 s in one transaction per Run, the slowest Run
  0.7 s; a full `VACUUM` afterwards took 2.3 s and left 391 MB;
- converting the whole database to incremental mode took 4.9 to 5.9 s;
- one heavy day of removal (102 Runs, 107,746 events) took 1.1 s on a
  converted database, and incremental compaction 0.8 s in 22 slices, the
  slowest 50 ms; the same day without conversion needed a 5.0 s full `VACUUM`.

## Consequences

- This amends ADR-0033 and ADR-0171: a `runs` row is still never deleted by
  Journal Retention, but Run Retention deletes it once the Run is past the
  outer window. ADR-0172's Doctor check is unchanged.
- At the 30-day default the file size is still bounded by Journal Retention;
  on this machine a 30-day window reclaims almost no bytes. A 7-day window, or
  a shorter Journal Retention, is what shrinks it.
- A Run removed by Run Retention is gone: `runs list`, `runs show`, `events`
  and token reports no longer see it. A Run Branch whose Run was removed is no
  longer listed by reconcile; it holds no bytes and Git keeps it.
- Schema version 23 adds `run_retention_sweeps`. An older binary refuses the
  migrated database as it refuses every newer schema.
