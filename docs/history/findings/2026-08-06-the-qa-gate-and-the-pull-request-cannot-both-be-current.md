---
status: done
created_at: 2026-08-06
updated_at: 2026-08-06
absorbed_by: 2026-08-06-rollup-qa-gates-and-verification-evidence.md
---

# 2026-08-06 — The QA gate and the pull request cannot both be current

## What was observed

An authored QA gate that inspects the pull request can never pass on its first run after any
corrective Task, because the artifact it verifies and the artifact the pull request shows are
different by construction.

The gate runs inside the Run Worktree, on the Run Branch. The pull request points at the Spec
target branch. The Run only integrates into the target branch **after** it settles — and a Run
whose gate failed settles `Unresolved`, which keeps the Run Branch unintegrated. So:

1. the gate verifies build `X` on the Run Branch and approves every application row;
2. it reads the pull request, finds head `X-1`, and fails its PR row;
3. the Run ends `Unresolved`, leaving `X` unintegrated;
4. the Supervisor integrates `X` by hand, pushes, and reopens the gate;
5. the gate runs again on `X+1` and repeats the same mismatch if anything moved.

Measured over one night driving three Specs on the `vortex` repository:

| Spec | Gate runs | Failing for a code defect | Failing for this cycle |
| --- | --- | --- | --- |
| 0016 | 5 | 1 | 3 |
| 0019 | 3 | 1 | 1 |
| 0017 | 3 | 1 | 2 |

**Eleven gate runs, three of which found a real defect.** Six of the remaining eight were this
ordering, at roughly forty minutes each.

## The stated order already says otherwise

Spec `0065-loop-order-and-verification-honesty` shipped this as its first Core Feature:

> The loop's order is stated once and consistently — implement the graph, open the Pull Request,
> drive its review to Clean, then request QA.

That order is right, and it is not what happens. An authored terminal `qa` Task is part of the
graph, so **the Daemon executes it the moment the last non-QA Task settles** — which is before the
Supervisor has opened the pull request, let alone driven its review to Clean. The gate cannot
wait for a pull request that the loop opens after the graph closes.

So the conflict is not between the gate and the Supervisor's discipline. It is between two
shipped contracts: the stated loop order puts QA after review, and the authored-gate mechanism
fires QA at graph close. Whichever the Supervisor obeys, the other is violated on the first pass.

## Root cause

Two contracts are individually sound and jointly circular:

- a Task Graph gate is terminal and runs before integration, because integration is what a
  successful Run does at the end;
- a gate that observes pull request review surfaces must read the pull request, and the pull
  request cannot show a commit the Run has not integrated.

`docs/agents/autonomous-work.md` already tells the Supervisor to "open the Pull Request and drive
its review to Clean before requesting QA" for Specs whose acceptance observes review surfaces.
That instruction assumes the gate is requested once, at the end. It does not survive a corrective
Task: the moment a finding produces one, the graph reopens, the gate becomes terminal again, and
the cycle restarts.

## What would settle it

Any one of these breaks the loop:

1. **Integrate before the gate, not after.** If the Daemon fast-forwards completed Task work onto
   the target branch before running an authored gate, the gate reads a pull request that contains
   what it is about to verify.
2. **Let the gate observe the Run Branch as the candidate.** If the PR row compares the pull
   request against the Run Branch and accepts "the PR will contain this once integrated", the row
   can pass on evidence rather than on timing.
3. **Make the PR row non-terminal.** If pull request review state is reported rather than gating,
   the gate stops failing for a fact the gate itself cannot change.
4. **Let the graph declare the gate as review-gated.** If `_tasks.md` could mark its `qa` Task as
   waiting for a review-clean pull request instead of for the last leaf, the authored mechanism
   would express the order Spec 0065 already states, and the Daemon would withhold the gate until
   the Supervisor reports the pull request clean.

The third is the smallest change and the least satisfying: it removes a real check. The first is
the honest one, and it also fixes the related surprise that a failed gate strands completed Task
work on an unintegrated Run Branch.

## Related

- `docs/findings/2026-08-05-review-issues-have-no-identity-across-rounds.md` — the other half of
  the same night. Together they accounted for most of the elapsed time: one hid findings that were
  still open, this one re-ran the gate against a stale pull request.
- `docs/findings/2026-08-05-preflight-prescribes-integrating-a-superseded-run-branch.md` — the
  manual integration this cycle forces is exactly where that trap is waiting. On Spec 0019 the
  superseded Run Branch held a QA report with the **same filename** as the passing one
  (`qa-report-2026-08-05.md`, because the name derives from the date), so following the
  Preflight's suggested `git merge --ff-only` would have overwritten a `pass` with a `fail`.

## Spec pointer

None yet.
