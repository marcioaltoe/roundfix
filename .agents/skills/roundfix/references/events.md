## Supervisor Run Event Stream

Use `roundfix events <run-id>` when a Supervisor, script, or CI process needs a
machine-readable Run projection. It replays one explicit Run from the Run
Database as `roundfix-events/v1` JSONL:

```bash
roundfix events <run-id>
roundfix events <run-id> --follow
roundfix events <run-id> --filter verification,outcome
```

The command never picks the newest Run implicitly. Discover Run IDs with
`roundfix runs list`, the Run Browser, or the Detached Run stdout report.
stdout contains JSONL records only. Diagnostics and validation errors go to
stderr. Missing Run ID, unknown Run ID, unknown filter category, and empty
filter exit `2`; store and output write errors exit `1`. When projection fails,
the command skips a record it cannot project, writes one warning to stderr with
its cursor, event kind, and projection error, then continues replay or follow;
stdout remains JSONL only. SIGINT or SIGTERM during `--follow` exits `130`
without a stdout trailer. With `--follow`, replay drains first and live follow
starts without duplicating the boundary event; terminal Runs replay and exit
immediately with `0`.

Default replay emits these public categories in journal cursor order:
`task-status`, `batch`, `verification`, `outcome`, and `agent-selection`.
`--filter` accepts a comma-separated subset of only those category names. Internal Run Event kinds,
raw Agent payloads, command strings, and diagnostic paths are not filters.
Internal Run Event kinds and raw Agent payloads are not projected.

Stable fields:

| category | fields |
| --- | --- |
| `task-status` | `schema`, `run_id`, `category`, `time`, `cursor`, `batch`, `work_item`, `phase`, `status`, `summary` |
| `batch` | `schema`, `run_id`, `category`, `time`, `cursor`, `batch`, `phase`, `summary` |
| `verification` | `schema`, `run_id`, `category`, `time`, `cursor`, `batch`, `work_item`, `attempt`, `phase`, `verdict`, `summary` |
| `outcome` | `schema`, `run_id`, `category`, `time`, `cursor`, `outcome`, `summary`; optional terminal `reason`, `next_action`, `review_issues_known`, `console_log`, `attach_command`, `evidence_kind`, `evidence_head_sha`, and `verified_head_sha` |

An Unobserved Verification adds classification `verification_unknown` with
`command`, `reason`, and `diagnostic_path` on both its `failed` and `verdict`
records. `reason` carries the runner cause or `reason unavailable`, and
`diagnostic_path` carries the retained path or `unavailable`. These fields let
a Supervisor distinguish "we did not find out" from a command verdict.

A Vacuous Verification adds classification `verification_vacuous` and the
`commands` field containing the commands that passed against the unchanged
tree.

Copy-paste examples:

```json
{"schema":"roundfix-events/v1","run_id":"run_20260710T120000Z_demo","category":"batch","time":"2026-07-10T12:00:00Z","cursor":1,"batch":1,"phase":"started","summary":"batch started"}
{"schema":"roundfix-events/v1","run_id":"run_20260710T120000Z_demo","category":"verification","time":"2026-07-10T12:00:01Z","cursor":2,"batch":1,"attempt":1,"work_item":"task_01","phase":"verdict","verdict":"failed","summary":"verification attempt 1 verdict failed"}
{"schema":"roundfix-events/v1","run_id":"run_20260710T120000Z_demo","category":"outcome","time":"2026-07-10T12:00:02Z","cursor":3,"outcome":"Unresolved","summary":"outcome Unresolved"}
```

Use `events` for unattended monitoring. Use `attach` for the human Live Run
View. Do not grep the Detached Run Console Log for state; it is a compact text
record, not a stable state API.


### Token usage

`roundfix events <run-id> --filter usage` selects the `usage` category of the
Run Event Stream (`roundfix-events/v1`). It is included in the default filter.
Each record names the Work Item scope, runtime, model, reasoning effort and
`token_basis` (`turn`, `request-sum` or `unreported`), with tokens and reported
cost when present. For example, its summary can say
`task_01 used 5639755 tokens (request-sum)`. An unreported prompt has no token
field. Retention can remove these events; durable usage totals remain
available through `roundfix runs show <run-id>`.
