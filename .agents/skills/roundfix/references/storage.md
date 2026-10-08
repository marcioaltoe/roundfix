## Run storage retention

`store.journal_retention` is a User Config and Project Config key that defaults
to `336h` (14 days). It accepts Go duration strings, and `0` keeps every Run
Event Journal and run artifact directory. Non-zero values make terminal Runs
older than the retention window eligible for pruning. Journal Retention never deletes
Active Runs, `runs` rows, or active-run locks, and it does not remove Review
artifacts under the Spec tree.

Use `roundfix gc [--dry-run]` to inspect or reclaim Run storage. Its Journal Retention block counts only Runs that still hold Run Event
Journal rows or an artifact directory. The dry run prints that reclaimable set, journal rows,
orphaned `runs/<id>` directories, and artifact bytes without changing
anything. A live `roundfix gc` deletes eligible Run Event Journal rows,
removes each pruned Run's `<artifact_dir>/runs/<run-id>` directory, removes
orphaned `runs/<id>` directories under the resolved run artifact root, and
reports Runs, journal rows, and artifact bytes reclaimed on stdout. With
`journal_retention: 0`, its Journal Retention block prints `Journal Retention: 0`
and `No journal pruning performed.`; Run Retention still runs.

Use the three explicit machine-wide storage surfaces separately from that
per-repository retention sweep:

```bash
roundfix gc compact
roundfix gc compact --apply
roundfix gc sanitize
roundfix gc sanitize --apply
roundfix storage report
```

`roundfix gc compact` previews Run Database bytes before, reclaimable, and
projected after without changing the database. Review those measurements, then
use `--apply` to run guarded compaction and print the measured bytes before,
reclaimed, and after. If a live writer advances the database while its exact
snapshot is measured, the bare preview falls back to immutable storage-report
measurements and remains read-only; `--apply` never uses that fallback.
Compaction with `--apply` refuses before mutation when an Active Run exists,
another Run Database connection can be writing, or temporary capacity is
insufficient; each refusal names the Active Run, writer condition, or measured
shortfall that blocked it. `roundfix gc compact --apply` converts a default-mode database to incremental
compaction. A new database already uses incremental mode. Run Retention
automatically returns freed pages in slices; the live `roundfix gc` can also
run guarded full compaction and convert a default-mode database. The budgeted
startup sweep never runs a full compaction.

`roundfix gc sanitize` discovers every recorded Artifact Root from the
machine-wide Run Database, classifies each root from durable and filesystem
evidence, and preserves ambiguous paths. It is a dry-run by default;
`--apply` removes only proven retention-eligible or absent Run artifact
directories. It does not replace or change the per-repository `roundfix gc`
surface.

`roundfix storage report` is read-only, accepts no flags, and needs no Git
repository. It reports measured bytes and row counts by database table,
repository, Run state, and Artifact Root without migrating, opening a writer,
or creating a missing Run Database.

Operational `implement`, `resolve`, and `watch` startup runs the same Journal
Retention prune best-effort when retention is non-zero. It prints the
successful cleanup stderr line only when it reclaimed something, shaped like:

```text
roundfix: pruned Run storage runs=<n> journal_rows=<n> artifact_bytes=<n>
```

Failures print one warning line shaped like:

```text
roundfix: warning: Journal Retention prune failed: <reason>
```

and never block the Run.


## Whole-Run retention

`store.run_retention_days` is User Config only, accepts `7`, `15` or `30`,
and defaults to `30`. A Project Config value is ignored with a warning.
Run Retention removes terminal Runs completed before that window, including
`runs`, Run Event Journal, Agent Selections, token usage, active-run locks
and the Run's artifact directory under a proven root. Review artifacts under
the Spec tree stay outside this sweep. Journal Retention remains the inner
window; `journal_retention: 0` does not disable Run Retention.

Active Runs, Runs named by any Delivery Queue item or Run link, Runs whose
recorded Run Worktree exists, and Runs with an artifact directory under an
unproven root are kept. Removal rechecks terminal state and queue references
inside each transaction. Queue status and reconcile retain the Runs they read.

`implement`, `resolve` and `watch` check after Journal Retention; `deliver start`
checks after recording the Delivery Queue. A sweep is due when none completed,
24 hours have elapsed, or the configured window changed. It has a two-second
budget, checked between Run removals and compaction slices, so a start can
overshoot by one step. A pause records no completion and resumes at the next
start. Quiet sweeps print nothing; other outcomes print one stderr line:

```text
roundfix: Run Retention removed runs=<n> rows=<n> database_bytes_reclaimed=<n>
roundfix: Run Retention paused after its 2s budget; it continues at the next Run start
roundfix: warning: Run Retention failed: <reason>
```

These diagnostics never change the start command's exit code. The removal
line sums removed rows across the five tables and reports logical database
bytes reclaimed, measured from page counts and page size.

`roundfix gc --dry-run` adds a Run Retention section with the window, cutoff,
removable Runs, kept counts by reason, rows per table and estimated reclaimable
database bytes. It writes nothing. `roundfix gc` runs the same sweep without a
budget; its section reports removed Runs and rows, reclaimed artifact bytes,
compaction and measured database bytes before and after. A guarded compaction
refusal appears as `skipped (<reason>)` and still exits `0`.
