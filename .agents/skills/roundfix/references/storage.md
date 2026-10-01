## Run storage retention

`store.journal_retention` is a User Config and Project Config key that defaults
to `336h` (14 days). It accepts Go duration strings, and `0` keeps every Run
Event Journal and run artifact directory. Non-zero values make terminal Runs
older than the retention window eligible for pruning. Retention never deletes
Active Runs, `runs` rows, or active-run locks, and it does not remove Review
artifacts under the Spec tree.

Use `roundfix gc [--dry-run]` to inspect or reclaim Run storage. Dry-run and
live reports count only Runs that still hold Run Event Journal rows or an
artifact directory. The dry run prints that reclaimable set, journal rows,
orphaned `runs/<id>` directories, and artifact bytes without changing
anything. A live `roundfix gc` deletes eligible Run Event Journal rows,
removes each pruned Run's `<artifact_dir>/runs/<run-id>` directory, removes
orphaned `runs/<id>` directories under the resolved run artifact root, and
reports Runs, journal rows, and artifact bytes reclaimed on stdout. With
`journal_retention: 0`, it prints `GC skipped` and performs no pruning.

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
shortfall that blocked it. Compaction is always explicit. Neither a manual
`roundfix gc` retention sweep nor the operational startup sweep compacts the
Run Database automatically.

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

