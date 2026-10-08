### events

```bash
roundfix events <run-id> [--follow] [--filter task-status,batch,verification,outcome,agent-selection,usage]
```

Writes only `roundfix-events/v1` JSONL records to stdout; diagnostics go to
stderr. The public Supervisor categories, in journal cursor order, are
`task-status`, `batch`, `verification`, `outcome`, `agent-selection`, and `usage`;
`--filter` accepts a comma-separated subset of those names only. Missing or
unknown Run IDs and invalid filters exit `2`; stream/store errors exit `1`;
interrupting `--follow` exits `130`. A terminal Run replays and exits `0`. Use
`events` for automation, `attach` for the human view, and the Detached Run
Console Log as a compact text record — not a state API.

If a journal entry cannot be projected,
the command skips a record it cannot project, writes one warning to stderr with
its cursor, event kind, and
projection error, then continues replay or follow. stdout remains JSONL only;
store and output write errors still exit `1`.

For a Detached Run's stable terminal subscription, use:

```bash
roundfix events <run-id> --follow --filter outcome
```

To inspect Task and Verification capacity flow, use the profile-led Run ID from
`implement --detach`:

```bash
roundfix events <run-id> --follow \
  --filter task-status,verification,outcome > run-events.jsonl \
  2> run-events.diagnostics
```

The Task-status records identify Agent work and Daemon settlement.
`daemon.verification` payloads use the canonical phases `waiting`, `started`,
`command-passed`, `failed`, and `verdict`. They carry the Task, numbered
attempt, shared or exclusive mode, retry identity, and capacity where
applicable. A Temporary Verification Failure adds classification `temporary`,
reason `temporary_verification_failure`, retained `diagnostic_path`, and
whether the exclusive retry remains available. Requested JSONL remains on
stdout; follow progress and operational diagnostics remain on stderr.

An Unobserved Verification adds classification `verification_unknown` with
`command`, `reason`, and `diagnostic_path` on both its `failed` and `verdict`
records. `reason` carries the runner cause or `reason unavailable`, and
`diagnostic_path` carries the retained path or `unavailable`.

A Vacuous Verification adds classification `verification_vacuous` and the
`commands` field containing the commands that passed against the unchanged
tree.

The outcome record carries the terminal state plus bounded reason and next
action when non-Clean. When available, it also carries Review Issue knowledge,
Console Log, Attach command, accepted Evidence kind and head, and the verified
parent head used by artifact-only inheritance.


Usage records are on by default. Select only prompt usage with `--filter usage`.
Each `daemon.token_usage` event projects to category `usage`, with `scope_kind`,
`scope_id`, `runtime`, `model`, `reasoning_effort`, and `token_basis`. A reported
prompt carries `tokens`; an unreported prompt carries `token_basis: unreported`
and omits `tokens`. Optional `input_tokens`, `output_tokens`,
`cached_read_tokens`, `cached_write_tokens`, `thought_tokens`, `cost_amount`,
and `cost_currency` appear only when recorded. Reported zero values are kept.
Tokens describe the prompt; cost is the adapter's cumulative Agent Session
reading. The schema stays `roundfix-events/v1`.

Surface Transcript 7 shows one prompt counted with basis `request-sum`:

```transcript
$ roundfix events run_20261001T120000Z_0123456789abcdef --filter usage
stdout:
{"schema":"roundfix-events/v1","run_id":"run_20261001T120000Z_0123456789abcdef","category":"usage","time":"2026-10-01T12:04:00Z","cursor":41,"work_item":"task_01","summary":"task_01 used 5639755 tokens (request-sum)","scope_kind":"task","scope_id":"task_01","runtime":"codex","model":"gpt-6.1-sol","reasoning_effort":"high","token_basis":"request-sum","tokens":5639755}
stderr:
exit: 0
```

An unreported prompt's summary is `<scope_id> reported no usage`. A failed or
stopped prompt still records whatever usage it returned. If persistence fails,
the Run's progress output carries `roundfix: warning: token usage not recorded
for <scope_kind> <scope_id>: <error>`; that warning does not change the prompt or
Run outcome.

An unknown Run ID exits `2`. When its `run_<YYYYMMDD>T<HHMMSS>Z_<hex>`
creation time predates the configured Run Retention cutoff, the diagnostic
adds `; Run Retention may have removed it, because it removes terminal Runs
that completed more than <N> days ago`, using `store.run_retention_days` from
User Config (30 by default). Other IDs, including `run_missing`, keep their
existing refusal. The hint describes a possible removal, not proof of one.
