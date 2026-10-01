### events

```bash
roundfix events <run-id> [--follow] [--filter task-status,batch,verification,outcome,agent-selection]
```

Writes only `roundfix-events/v1` JSONL records to stdout; diagnostics go to
stderr. The public Supervisor categories, in journal cursor order, are
`task-status`, `batch`, `verification`, `outcome`, and `agent-selection`;
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

