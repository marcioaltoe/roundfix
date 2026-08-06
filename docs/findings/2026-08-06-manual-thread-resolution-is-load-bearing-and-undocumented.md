# 2026-08-06 — Manual thread resolution is load-bearing, and the contract forbids it

status: pending

## What was observed

`2026-08-05-five-frictions-from-a-full-autonomous-spec-night.md` §5 states the constraint plainly:

> the Supervisor is barred from manual thread resolution

Over one night driving three Specs to merge on the `vortex` repository, **the Supervisor resolved
24 review threads by hand across three pull requests**. Without that, none of the three could have
merged. The barred action was the only exit from three distinct deadlocks, and only one of them is
the nitpick case §5 already describes.

| Deadlock | Threads | Why no actor could close it |
| --- | --- | --- |
| Findings on Roundfix's own bookkeeping files | 6 | The finding targets `docs/specs/_reviews/**/issue_*.md`. A Batch is forbidden from editing an unassigned Review Issue file, so the only fix is one no Agent may apply. |
| Threads Roundfix will not re-import | 16 | After one `resolve` settles a Round, later `fetch` returns `Review Issues: none` while GitHub still reports the threads open. The Agent is never handed them. |
| Below-threshold findings under required conversation resolution | 2 | The §5 case: filtered out by design, still counted by branch protection. |

Three `resolve` Runs failed every assigned issue in the first category with the same reason:

```
Batch 001 cannot apply the valid fix because its target is an unassigned Review Issue file.
```

That is a correct refusal against an unsatisfiable assignment. The issue is that nothing downstream
resolves it, so it repeats on every subsequent Round.

## Root cause

Three contracts are each defensible and jointly leave threads that no actor owns:

1. **A Batch may not edit an unassigned Review Issue file.** Correct — it prevents an Agent from
   rewriting another Agent's verdict. But CodeRabbit reviews those files, because they are Markdown
   in the repository, and the resulting findings are valid.
2. **Roundfix does not re-import a thread it considers terminal.** Correct as duplicate suppression.
   But "terminal in the artifacts" and "open at the Review Source" are different facts, and the
   second one is the one that blocks the merge.
3. **The Supervisor does not resolve threads.** Correct as a default — silent resolution is how
   review value erodes. But it is the only remaining actor.

The three together mean: a pull request can reach a state where merging requires an action that
every documented actor is barred from taking.

## What would settle it

The existing §5 proposal — an opt-in courtesy pass that posts the idempotent Outcome Comment and
resolves below-threshold threads — is the right shape. It needs to cover two more cases:

- **Findings on Roundfix-owned artifacts.** The cheapest fix is upstream: recommend excluding
  `docs/specs/**/_reviews/**` from the review source's scope, and say why in the Roundfix docs.
  A repository that reviews the tool's own bookkeeping produces findings nobody may fix. Applied
  on `vortex` as a `.coderabbit.yaml` path filter; it belongs in the documented setup, not in each
  repository's rediscovery.
- **Threads open at the source but terminal in the artifacts.** This is the mirror half of
  `2026-08-05-review-issues-have-no-identity-across-rounds.md`. Keying issues by thread fixes the
  import; until then, a terminal Run that leaves source threads open should say so in its report
  instead of reporting `0 unresolved`.

If manual resolution stays the escape hatch, it should be **documented and bounded** rather than
barred and practiced anyway: name the cases where it is legitimate, require a comment on the pull
request stating which case and why, and keep it outside any assigned Batch. Tonight it was always
accompanied by a comment naming the reason — but that was judgement, not contract.

## Why the framing matters

A rule that is barred in the docs and required in practice is worse than either a permission or a
prohibition. The Supervisor either learns to ignore the contract, or the human wakes up to click
through threads the machine could have closed with a stated reason. Both outcomes are cheaper to
avoid by naming the exception.

## Related

- `2026-08-05-five-frictions-from-a-full-autonomous-spec-night.md` §5 — the nitpick case and the
  courtesy-pass proposal this extends.
- `2026-08-05-review-issues-have-no-identity-across-rounds.md` — the import side of the same
  problem, with its 2026-08-06 addendum measuring the same three Specs.
- `2026-08-06-the-qa-gate-and-the-pull-request-cannot-both-be-current.md` — the ordering that makes
  every one of these repeat once per gate cycle.

## Spec pointer

None yet.
