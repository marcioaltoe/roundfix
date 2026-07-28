---
status: pending
created_at: 2026-07-28
updated_at: 2026-07-28
---

# QA gate — acceptance rows that need a live Pull Request are structurally unreachable, so their Spec can never be archived (2026-07-28)

The QA Agent runs in a Run Worktree checked out on `roundfix/run-<run-id>`. A
Pull Request belongs to the user's feature branch, never to that Run Branch, so
`gh pr view` from the QA surface always reports no Pull Request — no matter how
many are open. Combined with the rule that Agents never push, any acceptance row
whose journey requires a live Pull Request or reaches Final Push cannot be
executed by the gate.

Those rows are recorded `blocked`, which caps the verdict at `partial`. Since
`roundfix archive` requires `verdict: pass`, a Spec containing even one such row
can never be archived through the normal loop.

## Evidence

Spec `0039-review-source-evidence-and-detached-outcomes`, 2026-07-28. The gate
reported no product finding and every other row executed, yet returned
`partial`:

- `QA-10` — the successful approval journey continues through Final Push.
- `QA-11` — the public inheritance journey requires Roundfix to create its
  separate Daemon review-artifact commit and continue to Final Push.
- `QA-19` — blocked only for its non-skip Clean, CleanUnverified-after-push, and
  separate artifact-commit branches.

The gate was rerun *after* Pull Request #40 was opened on the feature branch, and
the same rows stayed blocked with the same reason. Confirmed directly:

```console
$ git rev-parse --abbrev-ref HEAD          # inside the QA Run Worktree
roundfix/run-run_20260728T041231Z_48d72c1b142ea37b
$ gh pr view --json number
no pull requests found for branch "roundfix/run-run_20260728T041231Z_48d72c1b142ea37b"
```

This is not a missing Pull Request. It is the Run Worktree isolation that every
spec Run depends on, meeting a QA matrix that reasons about the user's branch.

## Why it matters beyond one Spec

The affected rows are exactly the ones a review-oriented Spec exists to prove.
Spec 0039 owns Review Source Evidence and Merge-Ready acceptance; its most
valuable journeys are the ones the gate cannot reach. The stricter the Spec's
subject matter, the more of its matrix is unreachable — the gate is weakest
precisely where the risk is highest.

It also creates a false signal. `partial` reads like incomplete diligence when
the real cause is an environment the Agent cannot have. A reader comparing two
`partial` reports cannot tell "we could not verify this" from "we chose not to".

## Suggested resolution

1. Let the QA surface resolve the Pull Request from the Run's recorded target
   branch rather than the Run Worktree's checked-out branch. Roundfix already
   records the target branch on the Run; the QA Agent should be told the Pull
   Request number and repository rather than being left to infer them from a
   checkout that structurally cannot know.
2. Introduce a read-only Pull Request journey mode so approval, Merge-Ready
   acceptance, and artifact-inheritance rows can be observed against a real Pull
   Request without granting the Agent push authority — the constraint that
   currently blocks them is push, not observation.
3. Distinguish `blocked-by-environment` from `blocked-by-finding` in the report
   and in the verdict rule. A matrix whose only blocked rows are environmental
   should be able to reach `pass` when a supervisor records the equivalent
   evidence, instead of silently capping the Spec at `partial` forever.
4. Until one of the above lands, document how a supervisor closes such rows —
   the journeys *are* exercisable outside the gate by running `roundfix watch`
   against the real Pull Request, which performs the approval journey, creates
   the separate Daemon review-artifact commit, and reaches Final Push.

## Suggested acceptance checks

- A Spec whose matrix includes Pull Request journeys reaches `pass` when those
  journeys succeed against a real Pull Request.
- The QA Agent resolves the Pull Request without depending on the Run Worktree's
  branch name.
- A report distinguishes rows blocked by environment from rows blocked by a
  finding, and the verdict rule treats them differently.

## What worked — keep

- The gate refused to credit unexecuted journeys from Task evidence or Result
  prose. Capping the verdict was the honest response to rows it could not run;
  the defect is that it can never run them, not that it declined to guess.

## Addendum — 2026-07-28 — Suggestion 1 implemented; 2 and 3 remain

Suggestion 1 merged to `main` through Pull Request #40 (squash commit
`ed4abec`): `BuildQAPrompt` now renders the Run Worktree branch, the Spec
target branch, and the user checkout as labeled facts, and the rerun gate
resolved Pull Request #40 against the target branch and verified the build
match. Suggestion 4 was also exercised for real: the `watch` Run on Pull
Request #40 performed the approval journey, created the separate Daemon
review-artifact commit (`d240e99`), and reached Final Push. Suggestions 2 and
3 — read-only Pull Request journeys and the blocked-by-environment verdict
distinction — remain unimplemented, so Spec 0039 stays capped at `partial`
until one of them lands or a supervisor records the equivalent evidence.
