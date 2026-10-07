### reopen

```bash
roundfix reopen --spec <slug>
```

Returns a completed terminal QA Task to `pending` when one or more of its
dependencies are no longer completed, or when the gate has a **Late
Dependency**. Reopen proves a Late Dependency from the Task Graph manifest at
the commit that added the newest QA Report: a Task in the current dependency
closure that the recorded closure lacked reopens the gate whatever its status.
The invalidation record names it under `Dependencies added after the QA Report`.
This is the supported way to reopen a settled QA gate; do not edit the
QA Task file by hand. Reopen retains the prior QA Report and the Task's prior
Result. Without the commit that added the newest report, reopen refuses as
before. It creates no Run, writes no Run Event Journal entry, and never commits
or pushes.

Options:

- `--spec <slug>` — Spec slug under the configured Spec Root.

The command refuses before mutation when the terminal QA Task is not settled
(not `completed`), is neither stale nor above a proven Late Dependency, or when
the newest QA Report has no commit from which the prior Task Graph manifest can
be proven.

Exit codes:

- `0` — a settled QA gate that was stale or above a proven Late Dependency
  was reopened.
- `1` — reopen write failed.
- `2` — Preflight Validation failed, including a missing/invalid `--spec`, an
  unsettled QA gate, or a QA gate that is neither stale nor above a proven Late
  Dependency.
