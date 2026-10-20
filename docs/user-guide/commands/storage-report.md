### storage report

```bash
roundfix storage report
```

Measures bytes and row counts by repository, state, table, and Artifact Root.
It reads the machine-wide Run Database and recorded Artifact Roots from
Roundfix Home. The report is read-only, accepts no flags, and needs no Git
repository. It never migrates the database, locks for writes, or changes any
byte.

The report measures one Storage Snapshot, which reads the write-ahead log
when it holds frames. `Snapshot bytes` is the snapshot's total page size;
`Table bytes`, `Free bytes`, `Pointer-map bytes` and `Lock-byte page bytes`
group every page, so the grouped bytes equal the snapshot bytes exactly.
`File bytes` measures the main file separately; `WAL bytes` and
`Shared-memory bytes` measure its sidecar files. The reconciliation tolerance
is 0 bytes. No database or WAL byte changes, and no sidecar file is created.

If a live writer advances the database during `gc compact`'s preview, the
preview falls back to the storage report's Storage Snapshot. Its file bytes
and free bytes supply the fallback measurements; `gc compact --apply` keeps
its guarded compaction behavior.
