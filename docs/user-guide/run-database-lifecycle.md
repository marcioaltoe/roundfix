# Run Database lifecycle

Run Retention is the outer bound for terminal Runs in the Run Database, and
Journal Retention is the inner one for the Run Event Journal and artifact
directory. Every Roundfix-owned durable SQLite table has one lifecycle owner
and one retention rule below. SQLite-owned tables whose names start with
`sqlite_` are engine metadata and follow SQLite's lifecycle rather than a
Roundfix retention decision.

<!-- durable-table-lifecycle:begin -->
| Table | Lifecycle owner | Retention rule |
| --- | --- | --- |
| `runs` | Run lifecycle | Run Retention removes a terminal Run past its configured window (ADR-0255). Journal Retention never deletes these rows. |
| `active_run_locks` | Active Run lifecycle | Keep while the owning Run is Active and release at its terminal outcome. Run Retention removes any remaining lock with its terminal Run past the window (ADR-0255). Journal Retention never deletes these locks. |
| `interactive_defaults` | Interactive Input | Keep one current value per key, replacing it when Interactive Input records a newer value. No age-based retention applies. |
| `run_events` | Run Event Journal | Journal Retention may delete events only for terminal Runs older than its configured window. Active Run events are never eligible, and a zero window keeps everything. Run Retention removes remaining events with their terminal Run past the window (ADR-0255). |
| `run_token_usage` | Run lifecycle | Run Retention removes usage with its terminal Run past the window (ADR-0255). Journal Retention never deletes these rows. |
| `run_agent_selections` | Agent Selection lifecycle | Run Retention removes selections with their terminal Run past the window (ADR-0255). Journal Retention never deletes these rows. |
| `run_windows` | Run Window lifecycle | Keep one current window per repository. Replace it only through an explicit forced set and delete it only through an explicit clear. No age-based retention applies. |
| `delivery_queues` | Delivery Queue lifecycle | Keep one current queue per repository until an explicit Delivery Queue operation removes it. No age-based retention applies. |
| `delivery_queue_runs` | Delivery Queue lifecycle | Keep every Run link, including retries, with the owning queue and delete only through its foreign-key lifecycle. Journal Retention never deletes these rows. |
| `delivery_queue_items` | Delivery Queue lifecycle | Keep ordered items with their owning Delivery Queue and delete them only through that queue's foreign-key lifecycle. |
| `delivery_action_intents` | Delivery Queue lifecycle | Keep every external-action intent with its owning Delivery Queue item so an intent without a receipt remains observable. |
| `delivery_action_receipts` | Delivery Queue lifecycle | Keep each receipt with its intent and delete it only through that intent's foreign-key lifecycle. |
| `run_retention_sweeps` | Run Retention lifecycle | Keep one row recording the last completed sweep and its window, replacing it after each completed sweep (ADR-0255). |
<!-- durable-table-lifecycle:end -->

Run Retention never touches `interactive_defaults`, `run_windows`, or Delivery
Queue tables (`delivery_queues`, `delivery_queue_runs`, `delivery_queue_items`,
`delivery_action_intents`, and `delivery_action_receipts`). The GC Command
keeps Active Runs and active-run locks untouched. Journal Retention can prune
eligible `run_events` and matching Artifact Directory content before Run
Retention removes the whole terminal Run.

The 2026-10-08 measurement recorded by [ADR-0255](../adr/0255-run-retention-removes-terminal-runs-whole-and-compacts-incrementally.md)
found a 1.1 GB Run Database, with 1.08 GB in `run_events` across 832k events
written in the preceding two weeks, about 77 MB a day across 313 Runs. That
growth is the reason the outer Run Retention bound now covers the Run row and
its dependent records as a whole.

The `internal/store` lifecycle policy test compares this table with every
Roundfix-owned durable table in a migrated Run Database. A schema change that
adds or removes a table without updating this policy fails that check.
