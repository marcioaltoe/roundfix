---
type: fix # feat | fix | perf | refactor
status: open # open | promoted | declined
created: 2026-08-12
spec: null # Spec slug when status: promoted
reason: null # required when status: declined
---

# A hook failure kills a Run whose work was already verified

## Symptom

The Daemon runs the authoritative Verification and then commits. When the
repository's `pre-commit` hook refuses, `lint-staged` reverts and the Run ends
`Failed`. The repair loop covers a Verification failure, not a hook failure.

Three Runs died this way in one Spec, every time with the work **correct and
already verified**, left staged in the Task Worktree:

- a route function of 82 lines against the repository's limit of 80
- a generated OpenAPI client file of 2,462 lines against a limit of 500
- `Array#sort()` where the rule asks for `toSorted()`

All three findings were legitimate. The cost was not the check — it was **where**
it ran. Two gates where one cannot be satisfied nor recovered is a design
conflict, and the invariant that resolves it is written nowhere: a commit hook
may never be stricter than the authoritative Verification.

A second symptom rides along. `settle`, the command designed for recovery,
refuses work that is `completed` but uncommitted:

```text
Task task_13 has no failed settle surface; candidates: <task worktree>: status completed;
<run worktree>: status pending; <repo>: status pending; Task task_13 status is completed;
nothing to do
```

The Daemon marked it `completed`, the commit died in the hook, and recovery
refused because the status was not `failed`. It happened twice, and both times
`git diff HEAD` was extracted from the worktree and applied by hand. Whether
Spec 0092 already closes this half needs checking before any Task runs.

## Where

`internal/daemon` — the commit that follows the authoritative Verification —
and `internal/cli/settle.go`, whose recovery contract assumes lost work is
always `failed`.

## Expected

A hook refusal is either impossible by construction — the Daemon is the
verification authority and commits accordingly — or it is a repairable class
that spends a repair round, or at minimum a detected case that names the
recovery instead of ending `Failed` in silence. Recovery covers the state
"verified, settled, not committed".

## Evidence

`docs/findings/2026-08-12-three-consecutive-specs-measure-the-loop.md` records
the adjacent friction of the same session. This entry is minted from the Inbox
Entry `inbox/roundfix/2026-08-08-run-morre-em-falha-de-hook-e-settle-nao-recupera.md`
in the Secondbrain, observed in `fluxus` Spec 0021 across 2026-08-07/08.
