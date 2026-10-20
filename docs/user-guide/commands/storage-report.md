### storage report

```bash
roundfix storage report
```

Measures bytes and row counts by repository, state, table, and Artifact Root.
It reads the machine-wide Run Database, the recorded repositories' Git
metadata, and their resolved Artifact Roots. The report is read-only, accepts
no flags, and needs no Git repository. It never migrates the database, locks for writes, or changes any
byte.

A Run that recorded no Artifact Root is measured under its repository's default root, `~/.roundfix/artifacts/<repository id>`.
A relative Artifact Root is measured under the Run's Git root.
An absolute Artifact Root is measured as recorded. For older Runs, the report
uses the checkout's default root when it holds the Run's directory and the
repository's default does not, or when the repository default cannot be resolved.
A missing root stays listed as `missing`; the report creates no directory.

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
