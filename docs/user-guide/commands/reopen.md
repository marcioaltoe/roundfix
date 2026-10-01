### reopen

```bash
roundfix reopen --spec <slug>
```

Returns a completed terminal QA Task to `pending` when one or more of its
dependencies are no longer completed. This is the supported way to reopen a
settled QA gate; do not edit the QA Task file by hand. Reopen retains the prior
QA Report and the Task's prior Result, and records the stale dependency IDs
that invalidated the report. It creates no Run, writes no Run Event Journal
entry, and never commits or pushes.

Options:

- `--spec <slug>` — Spec slug under the configured Spec Root.

The command refuses before mutation when the terminal QA Task is not settled
(not `completed`) or is not stale (every dependency is still `completed`).

Exit codes:

- `0` — stale QA gate reopened.
- `1` — reopen write failed.
- `2` — Preflight Validation failed, including a missing/invalid `--spec`, an
  unsettled QA gate, or a non-stale QA gate.

