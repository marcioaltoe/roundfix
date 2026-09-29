---
status: accepted
created_at: 2026-09-29T00:00:00Z
updated_at: 2026-09-29T00:00:00Z
deprecated_at: null
superseded_by: null
---

# A retention prune reports only what it reclaimed

Journal Retention prunes the Run Event Journal rows and the artifact directory
of a terminal Run that completed before the cutoff
([ADR-0033](0033-the-run-event-journal-is-pruned-by-retention.md)). It never
deletes the `runs` row, so a Run stays eligible for every later prune. The
prune and its report did not distinguish a Run that still held something from
one emptied by an earlier prune. On 2026-09-29 `roundfix gc` reported
`Runs pruned: 13` with no journal rows, no artifact bytes and no orphan
directories removed, and a `gc --dry-run` right after listed the same 13 Runs
as eligible. The best-effort prune that every operational command runs had
printed `pruned Run storage runs=12 journal_rows=0 artifact_bytes=0` to
stderr in 133 Run console logs. Each of those prunes also took the machine-wide
write lock to delete nothing.

Now a Run counts as reclaimed, and as eligible in a dry run, only when it still
holds Run Event Journal rows or an artifact directory at the cutoff. The prune
counts each eligible Run's events on the read connection first. It takes the
write lock only when some eligible Run still has events, and it deletes only
those Runs' events. A second prune after a complete one reclaims nothing and
reports zero, and the operational sweep prints nothing then.

## Consequences

Eligibility is unchanged: a terminal Run that completed before the cutoff is
still the only Run whose journal and artifacts may be pruned, and `runs` rows
and Active Run locks are still never touched. Only the count and the lock use
change. An artifact directory that outlived its journal rows, for example after
a failed removal, is still found and removed, because the directory alone makes
its Run reclaimable.

This refines ADR-0033 without superseding it.
