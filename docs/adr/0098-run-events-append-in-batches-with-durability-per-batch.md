---
status: accepted
created_at: 2026-08-06T21:23:58Z
updated_at: 2026-08-11T21:35:56Z
deprecated_at: null
superseded_by: null
---

# Run Events append in batches with durability per batch

Run Event appends group into one immediate transaction per batch instead of one
transaction per event, and `synchronous` stays at SQLite's `FULL` default, so a
crash can lose at most the batch in flight. Today every agent stdout line is its
own `BEGIN IMMEDIATE`, existence probe, insert, commit, and fsync; the measured
journal holds 1,645,457 events across fourteen days, one Run contributing
42,000. Batching removes the great majority of those fsyncs and lock
acquisitions without weakening the durability guarantee, which matters because
ADR-0030 removed per-Batch agent log files precisely on the grounds that the
journal stores every raw payload durably — the journal is the only durable copy,
not a duplicate. A batch closes when it reaches its configured event count,
when its linger deadline expires, before an immediate-durability event, and on
explicit flush during error, shutdown, Agent teardown, or terminal settlement.
The Store owns the one open batch, so every sink handle shares those boundaries.

Events keep publisher order and receive a contiguous per-Run cursor range in
the write transaction. A failed begin, insert, or commit leaves the batch
intact and returns the error to the Run lifecycle. A commit with an ambiguous
outcome is reconciled by reading that exact cursor range: an exact byte match
settles the batch, no rows permits one retry with the same cursors, and a
partial or different match fails as corruption. The primary key therefore
makes retry idempotent without dropping or duplicating a raw event. A process
crash can lose only the Store's open batch: never more than the configured
count and never more than the configured linger window for a quiet publisher.

The faster `synchronous=NORMAL` was considered and declined for now: under WAL
it would survive process crashes and lose only on host crash or power loss, but
the batch-level guarantee already captures most of the gain, and a durability
relaxation should be bought with a measurement rather than assumed. Retrying a
failed batch with newly allocated cursors was also rejected because an
ambiguous successful commit would duplicate the raw events. ADR-0008 stands
untouched: payloads are still stored as raw producer JSON, write-once and
read-as-blob, and batching changes when bytes are committed, never what they
contain.
