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
not a duplicate.

The faster `synchronous=NORMAL` was considered and declined for now: under WAL
it would survive process crashes and lose only on host crash or power loss, but
the batch-level guarantee already captures most of the gain, and a durability
relaxation should be bought with a measurement rather than assumed. ADR-0008
stands untouched: payloads are still stored as raw producer JSON, write-once and
read-as-blob, and batching changes when bytes are committed, never what they
contain.
