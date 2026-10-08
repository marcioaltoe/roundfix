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
