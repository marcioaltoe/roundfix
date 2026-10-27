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
snapshot is measured, the bare preview falls back to the storage report's Storage Snapshot and remains read-only; `--apply` never uses that fallback.
Compaction with `--apply` refuses before mutation when an Active Run exists,
another Run Database connection can be writing, or temporary capacity is
insufficient; each refusal names the Active Run, writer condition, or measured
shortfall that blocked it. `roundfix gc compact --apply` converts a default-mode database to incremental
compaction. A new database already uses incremental mode. Run Retention
automatically returns freed pages in slices; the live `roundfix gc` can also
run guarded full compaction and convert a default-mode database. The budgeted
startup sweep never runs a full compaction.

`roundfix storage report`, `roundfix gc sanitize` and Run Retention resolve a
recorded Artifact Directory with one resolver. `roundfix gc sanitize` groups
Runs by the Artifact Root each wrote to and classifies each root with its
existing proofs. A Run whose Artifact Root cannot be resolved makes its group
`unsafe`, so the group is preserved. It is a dry-run by default; `--apply`
removes only proven retention-eligible or absent Run artifact directories and
rechecks the Run Database immediately before removing a directory proved
absent. Each root reports the recorded values that resolved to it with
`Recorded as:`, followed by its per-root candidate, preserved and reclaimable
lines. It does not replace or change the per-repository `roundfix gc` surface.

`roundfix storage report` is read-only, accepts no flags, and needs no Git
repository. It reports measured bytes and row counts by database table,
repository, Run state, and Artifact Root without migrating, opening a writer,
or creating a missing Run Database.

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

## Worktree Build Caches

`worktree.caches` declares rebuildable directories in User Config or Project
Config, with an empty default. After its Run Retention section,
`roundfix gc --dry-run` prints the declared entries, reclaimable directory
count, bytes, kept count and paths, without writing:

```text
  Worktree Build Caches declared: .gocache
  Worktree Build Caches reclaimable: 1
  Worktree Build Cache bytes reclaimable: 11
  Worktree Build Caches kept: 1
  Reclaimable Worktree Build Caches:
    <worktree>/.gocache
```

The live report uses `Worktree Build Caches removed`,
`Worktree Build Cache bytes reclaimed` and `Removed Worktree Build Caches`.
It counts only directories it actually removed. GC rereads ownership state
immediately before each removal and removes only the declared directory.
A removal failure is diagnosed on stderr; other removals proceed, and the
command exits `1` with the report of successful removals.

A terminal Run's recorded Run Worktree and retained sibling Task Worktrees,
and a Delivery Queue item worktree at stage `merged`, are eligible regardless
of retention age. Active Runs and items in any other stage keep their caches.
GC skips worktrees outside `worktree.location`, declared paths that are
symlinks, and paths that resolve outside their worktree. Source files and Git
metadata stay intact. Shared Build Caches are declared by
`worktree.shared_caches`; GC trims them above their declared limit by removing
only entries no build used for the declared idle time, oldest first.
With `worktree.caches: []`, GC omits the Worktree Build Cache section. `gc compact` and
`gc sanitize` retain their separate scopes.

## Shared Build Caches

`worktree.shared_caches` declares build caches shared by multiple worktrees.
Each entry names its path template, kind, maximum size and idle window. GC
trims a declared Shared Build Cache above its limit by removing only cache
entries no build used for the declared idle time, oldest first.

### Automatic trims

Roundfix also trims every declared Shared Build Cache by itself, with the
rules above: at the start of every Implement Run, after the Run is recorded
and before its Run Worktree exists, and before each repository gate of the
Delivery Queue. One budget of thirty seconds covers all declared caches at
each trigger; when it runs out, the trim keeps what it removed and the next
trigger continues. A refusal, an error or the end of the budget never fails,
parks or delays the Run or the gate beyond the budget. Each Implement Run
records one `daemon.status` Run Event per declared cache whose payload holds
`event: shared_cache_trim`, `trigger: run_start`, `path`, `outcome`
(`trimmed`, `unchanged`, `kept`, `failed` or `stopped`), `size_bytes`,
`limit_bytes`, `removed_entries`, `removed_bytes` and, when there is one,
`cause`. A trim that removed entries, failed or stopped also writes a line,
on standard error at Run start and in the Queue Owner's log at the gate:

```text
roundfix: [prefix]Shared Build Cache <path>: removed <n> idle entries (<bytes> bytes) P
roundfix: [prefix]Shared Build Cache <path>: trim stopped after its <budget> budget P; removed <n> idle entries (<bytes> bytes); the next trim continues
roundfix: warning: [prefix]Shared Build Cache <path>: trim failed P: <cause>; C
```

Without a declared cache nothing runs and no line or event is written.
`roundfix gc` still trims on demand (ADR-0334).

## Item Binaries

After the Worktree Build Cache section, GC scans
`<artifact dir>/delivery/*/item-binary` for the repository it runs in. An Item
Binary directory is reclaimable when that repository's Delivery Queue has no
item for its slug or holds it at stage `merged`. Every other stage, including
`parked`, keeps its copies.

When at least one Item Binary directory exists, `roundfix gc --dry-run`
prints these lines and the sorted reclaimable paths without writing anything:

```text
  Item Binaries reclaimable: <n>
  Item Binary bytes reclaimable: <bytes>
  Item Binaries kept: <n>
  Reclaimable Item Binaries:
    <path>
```

The live report uses `Item Binaries removed`, `Item Binary bytes reclaimed`
and `Removed Item Binaries`, counting only directories it removed. GC rereads
the queue immediately before each removal. It keeps symlinked Item Binary
directories and paths that resolve outside `<artifact dir>/delivery`.
The item's other artifacts, including `item-binary-build.log` and
`repository-gate.log`, stay. A removal failure is diagnosed on stderr;
other removals proceed, and GC exits `1` with the report of successful removals.
With no Item Binary directory, GC output stays byte-identical.

## Retention policy

ADR-0290's retention policy is documented in
`docs/user-guide/commands/gc.md`. Run journals and artifacts follow Journal
Retention; terminal Run records follow Run Retention. Worktree Build Caches
follow ended ownership, and Item Binaries stay while their queue item exists
and is not merged. Run Worktrees and Run Branches leave only when reconciliation
proves them `safe` or `superseded`, never by age. Pre-PR Review Records,
dispositions, delivery logs and archive evidence are retained.

Preview without changing storage with `roundfix gc --dry-run`,
`roundfix gc compact` and `roundfix reconcile`.
