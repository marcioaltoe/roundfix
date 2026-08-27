# Run Database lifecycle

Roundfix keeps a compact Run index and bounds only the Run Event Journal.
Every Roundfix-owned durable SQLite table has one lifecycle owner and one
retention rule below. SQLite-owned tables whose names start with `sqlite_` are
engine metadata and follow SQLite's lifecycle rather than a Roundfix retention
decision.

<!-- durable-table-lifecycle:begin -->
| Table | Lifecycle owner | Retention rule |
| --- | --- | --- |
| `runs` | Run lifecycle | Preserve as the compact Run index. Journal Retention never deletes these rows. Any future bound requires measured justification and an explicit policy change. |
| `active_run_locks` | Active Run lifecycle | Keep while the owning Run is Active and release when that Run reaches a terminal outcome. Journal Retention never deletes these locks. |
| `interactive_defaults` | Interactive Input | Keep one current value per key, replacing it when Interactive Input records a newer value. No age-based retention applies. |
| `run_events` | Run Event Journal | Journal Retention may delete events only for terminal Runs older than its configured window. Active Run events are never eligible, and a zero window keeps everything. |
| `run_agent_selections` | Agent Selection lifecycle | Keep as evidence with the owning Run. Delete only with that Run through the existing foreign-key lifecycle, unless a future explicit evidence-retention rule records measured justification. Journal Retention never deletes these rows. |
| `run_windows` | Run Window lifecycle | Keep one current window per repository. Replace it only through an explicit forced set and delete it only through an explicit clear. No age-based retention applies. |
<!-- durable-table-lifecycle:end -->

The policy keeps the existing deletion boundary unchanged: the GC Command can
delete eligible `run_events` and matching Artifact Directory content, but it
does not delete `runs`, `active_run_locks`, `interactive_defaults`, or
`run_agent_selections` or `run_windows` rows. Active Run locks leave the table
only through the Run lifecycle's terminal transition, not because they aged
past a retention window.

The current measurements support that boundary. The storage investigation
found 279 `runs` rows occupying 118,784 bytes, so the Run index did not justify
a bound; the earlier unbounded Run Event Journal had reached about 220 MB and
is governed by [ADR-0033](../adr/0033-the-run-event-journal-is-pruned-by-retention.md).
The same investigation found no Agent Selection rows, so it supplied no
measured justification for an independent evidence-retention window. See the
[storage lifecycle finding](../findings/_archived/2026-07-17-global-run-storage-sanitation-and-compaction.md#3-storage-lifecycle-policy-does-not-cover-every-durable-run-record).

The `internal/store` lifecycle policy test compares this table with every
Roundfix-owned durable table in a migrated Run Database. A schema change that
adds or removes a table without updating this policy fails that check.
