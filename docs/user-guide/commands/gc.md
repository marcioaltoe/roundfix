### gc

```bash
roundfix gc --dry-run
roundfix gc
```

Non-interactive storage reclamation. Resolves `store.journal_retention`,
prunes eligible terminal Runs' Run Event Journal rows and
`<artifact-dir>/runs/<run-id>` directories, removes orphaned `runs/<id>`
directories, and reports Runs, journal rows, and artifact bytes reclaimed.
Journal Retention reports count only Runs that still hold Run Event Journal rows
or an artifact directory. `--dry-run` lists that reclaimable set without
deleting. Journal Retention keeps Run rows and active-run locks. Review artifacts
under the Spec Root stay intact.

## Run Retention

After the Journal Retention report, GC applies `store.run_retention_days` from
User Config: `30` by default, with only `7`, `15` or `30` accepted. A Project
Config value is ignored with a warning. Terminal Runs whose completion time is
before the cutoff leave the database whole, including their journal, Agent
Selections, token usage and locks. GC removes each Run's artifact directory
under its proven Artifact Root before removing the database rows.

Active Runs stay. The report counts terminal Runs kept past the cutoff by
reason: queue-referenced, worktree present, and artifact root unproven. GC keeps
a Run when any Delivery Queue item or Run link names it, its recorded Run
Worktree exists, or its artifact directory exists under a root whose physical
path inside Roundfix Home cannot be proven. Multiple reasons can apply to one
kept Run.

`roundfix gc --dry-run` opens only the reader and changes nothing. Its Run
Retention section prints the window, cutoff, removable Runs, kept reasons,
rows per table and `Database bytes reclaimable (estimated)`, followed by the
removable Run IDs. The bytes are an estimate from event payloads and summaries;
they are not a promise of the space compaction will return.

The live report prints removed Runs and rows, Run artifact bytes reclaimed,
and database bytes before and after measured as page count times page size.
Compaction reports one of:

- `incremental (<n> pages)` when it returns free pages in slices;
- `full` when it rebuilds an existing non-incremental database;
- `full (converted to incremental)` when it converts a default-mode database;
- `not needed` when there are no free pages in incremental mode;
- `skipped (<reason>)` when the guarded full compaction refuses, including
  when an Active Run is present. This refusal still exits 0.

With `store.journal_retention: 0`, GC prints its normal `GC dry-run` or
`GC complete` header, `Journal Retention: 0` and `No journal pruning performed.`
and still applies Run Retention. Repeating a removal reports only what that
invocation reclaimed.

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
symlinks, and paths that resolve outside their worktree. Source files, Git
metadata, the user's checkout cache and the shared Go build cache stay intact.
With `worktree.caches: []`, GC omits the Worktree Build Cache section. `gc compact` and
`gc sanitize` retain their separate scopes.

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

[ADR-0290](../../adr/0290-merged-items-give-back-their-item-binaries-and-run-worktrees-leave-only-by-proof.md)
defines what Roundfix Home keeps and who reclaims it:

| Area | Kept while | Reclaimed by |
| --- | --- | --- |
| Run Event Journal and run artifact directory | the Run is Active or younger than Journal Retention | the GC Command and Run Retention (ADR-0033, ADR-0255) |
| Run row and its records | the Run is not terminal, is queue-referenced, keeps its Run Worktree, or is younger than Run Retention | Run Retention (ADR-0255) |
| Run Worktree and Run Branch | reconciliation does not prove them `safe` or `superseded` | `roundfix reconcile --apply` and the Delivery Queue's release after a merge, never by age (ADR-0053, ADR-0212, ADR-0289) |
| Worktree Build Cache | its Run is not terminal or its item is not `merged` | the GC Command (ADR-0266) |
| Item Binary | its Delivery Queue item exists and is not `merged` | the Delivery Queue's release after a merge, and the GC Command |
| Pre-PR Review Records, dispositions and delivery logs | always | nothing |
| Archive evidence | always | nothing: it lives in Git |

Preview without changing storage with `roundfix gc --dry-run`,
`roundfix gc compact` and `roundfix reconcile`. Run Worktrees leave only by
proof; retention age alone never releases them.
