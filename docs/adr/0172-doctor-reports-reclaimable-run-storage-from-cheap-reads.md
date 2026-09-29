---
status: accepted
created_at: 2026-09-29T00:00:00Z
updated_at: 2026-09-29T00:00:00Z
deprecated_at: null
superseded_by: null
---

# Doctor reports reclaimable Run storage from cheap reads

Run storage grows with every Run, and only the GC Command says how much of it
can be reclaimed, and only when someone runs it. A capture of 2026-09-17
recorded 4.4 GB of Run storage with 619 MB reclaimable after one day, and
nobody had been told. The
full measurement, `roundfix storage`, walks every recorded Artifact Root and
scans the database page by page, so it cannot run on a routine path.

The Doctor Command now carries a `storage` check. It opens the Run Database
through the read-only reader and reads only cheap facts:

- the free-page count and page size, which SQLite keeps in the database header;
- the terminal Runs past the Journal Retention cutoff that still hold Run Event
  Journal rows, counted per candidate;
- whether each such Run's artifact directory exists under the repository's
  Artifact Root, checked with one `lstat` and never walked.

When a Run past the cutoff still holds something, or free space reaches the
threshold, the line reports `found` and names `roundfix gc` or
`roundfix gc compact`. Otherwise it reports `ok`. A database it cannot read
reports `partial`.

## Consequences

The check never reports `failed`, so it never changes Doctor's exit code:
reclaimable storage is a notice, not a readiness defect. It never deletes,
compacts or writes anything, and it never runs `dbstat` or measures directory
bytes. The count uses the same reclaimable predicate as `roundfix gc --dry-run`
([ADR-0171](0171-a-retention-prune-reports-only-what-it-reclaimed.md)), so
the two never disagree on which Runs remain.

`runs list` does not carry the notice. It is the Agents' hot path, and it
already prints at most one note.
