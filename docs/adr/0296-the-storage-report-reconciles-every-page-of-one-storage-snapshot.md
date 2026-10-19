---
status: accepted
created_at: 2026-10-19T20:00:00Z
updated_at: 2026-10-19T20:00:00Z
deprecated_at: null
superseded_by: null
---

# The storage report reconciles every page of one Storage Snapshot

On 2026-10-18 `roundfix storage report` failed while the Delivery Queue ran
with `file bytes=1113714688 grouped bytes=1112350720 difference=1363968
tolerance=4096`. The report adds the `dbstat` table pages and the free pages
and compares the sum with the Run Database file size within one page. SQLite's
documentation says "The DBSTAT table only reports on the content of btrees
within the database file. Freelist pages, pointer-map pages, and the lock page
are omitted from the analysis." The Run Database uses incremental auto-vacuum
(ADR-0255), so it holds one pointer-map page for every 820 pages of 4096 bytes,
and a file past 1 GiB holds the lock-byte page. At 271,903 pages that is 332
pointer-map pages and one lock-byte page: 333 pages of 4096 bytes, exactly the
reported difference. On 2026-10-19 the report failed the same way on an idle
store whose write-ahead log held 0 bytes (337 pages). Separately, the report
opens the database with SQLite's `immutable` parameter, which skips locking and
the write-ahead log; SQLite warns that a file that "does in fact change can
result in incorrect query results and/or SQLITE_CORRUPT errors". With a live
WAL holding 2,000 rows not yet checkpointed, the immutable reader measured
3,013 pages while a WAL reader of the same files measured 5,020. The maintainer decided on 2026-10-19
that the report accounts for the live WAL or reads a consistent snapshot, and
that widening the tolerance is not the fix.

## Decision

**One Storage Snapshot measures every page group.** The storage report opens
the Run Database the way `roundfix gc --dry-run` already does: as an immutable
file when the `-wal` file is absent or holds no frame (32 bytes or fewer), and
otherwise as a read-only WAL reader, so committed frames stay visible and no
sidecar file is created. Every measurement runs inside one read transaction on
that connection.

**Every page is grouped.** The grouped bytes are the table pages that `dbstat`
reports, the free pages, the pointer-map pages and the lock-byte page. The
pointer-map and lock-byte counts come from the database file format: with
incremental or full auto-vacuum the first pointer-map page is page 2 and one
follows every `U/5 + 1` pages, where `U` is the usable page size, moved one
page later where it would fall on the lock-byte page; the lock-byte page is the
page holding byte offset 1,073,741,824 when the database is that large. The
grouped bytes must equal the snapshot's page count times its page size exactly,
so the reconciliation tolerance is 0.

**The files are reported apart.** The report prints the main file's bytes, the
`-wal` and `-shm` files' bytes, and the snapshot's bytes. The main file differs
from the snapshot whenever the write-ahead log holds frames not yet
checkpointed, so the main file is reported, never reconciled.

## Considered Options

- Raising the tolerance to the observed difference was rejected by the
  maintainer, and it would grow with the database.
- Counting the pages `dbstat` does not report as page count minus table and
  free pages was rejected: the reconciliation would then hold by construction
  and could no longer catch a grouping defect.
- Always opening a WAL reader was rejected: on an idle store it creates the
  `-wal` and `-shm` files that the report promises not to create.
- Checkpointing before measuring was rejected: the report is read-only.

## Consequences

- The report works on a store of any size and while a queue writes.
- The text report gains the pointer-map, lock-byte page, WAL, shared-memory
  and snapshot lines, and its tolerance line reads 0 bytes.
- `roundfix gc compact`'s live-writer fallback reads the same snapshot.
- This refines ADR-0255, whose incremental auto-vacuum created the pointer-map
  pages, and ADR-0172, whose Doctor storage check stays a cheap read without
  `dbstat`.
