---
status: done
created_at: 2026-07-28
updated_at: 2026-07-28
---

# Failed QA Runs strand Run Branches that block the next review Run, and nothing can classify them as superseded (2026-07-28)

Every QA gate attempt creates a Run Branch and commits its QA Report there. When
the gate fails, the outcome is Unresolved, so nothing integrates and the branch
is kept. After a few attempts the repository accumulates Run Branches that each
hold one commit, and Branch Integrity Preflight then refuses to start a review
Run until a human resolves every one of them.

The Reconcile Command cannot clear them: they hold commits absent from the
target branch, so it classifies them `unintegrated` and preserves them. That is
the correct default — but the operator is left doing the judgement by hand, and
the judgement is not "is this integrated" but "is this superseded".

## Reproduction, 2026-07-27

Spec 0038 needed three QA attempts before passing. Starting `roundfix watch` on
the resulting Pull Request then failed Preflight:

```text
Branch Integrity Preflight refused pending Run Branch work for PR Head Branch
"ma/terminal-run-worktree-reconciliation".
- branch=roundfix/run-run_20260727T194155Z_… ahead_commits=1
- branch=roundfix/run-run_20260727T200825Z_… ahead_commits=1
- branch=roundfix/run-run_20260727T201759Z_… ahead_commits=1
Next action: inspect each pending Run Worktree, then run the listed integration
command from the repository root when it is safe.
```

`roundfix reconcile` reported each as `unintegrated` / `preserve`. Its printed
next action — `git merge --ff-only <branch>` — would have been **actively
wrong**: each branch carried a *failing* QA Report at the same dated path the
*passing* report now occupies, so integrating would have replaced the passing
report with a failing one.

The correct action was to discard all three, which required inspecting each
commit by hand to confirm it held nothing but a superseded report, then
`git worktree remove --force` plus `git branch -D` — exactly the manual work the
Reconcile Command exists to eliminate.

## The gap

Roundfix models two states for kept work: integrated or not. Real recovery has a
third — **superseded**: committed, never integrated, and now obsolete because a
later attempt produced the authoritative version at the same path. Nothing
distinguishes it, so the safe default (`preserve`) is applied to work that
should be discarded, and the suggested remedy (`merge --ff-only`) is offered for
work that must not be merged.

The volume scales with QA difficulty: the harder a Spec is to pass, the more
stranded branches it leaves, and the more likely the next review Run is blocked
by them.

## Suggested resolution

1. Give the Reconcile Command a `superseded` classification for a terminal Run
   whose only unintegrated commits touch paths a later Run already wrote on the
   target branch, and let `--apply` release those with the reason recorded.
2. Do not offer `git merge --ff-only` as the next action for a branch whose
   integration would overwrite a newer artifact on the target. Detect the
   conflict-by-supersession and say so.
3. Consider not creating a Run Branch commit for a failing QA gate at all, or
   writing the failing report directly to the target branch. The report is
   evidence, not Task output; stranding it behind an unintegrated branch is what
   creates the debris.
4. Have Branch Integrity Preflight name which pending branches are safely
   discardable, so the operator is not left inspecting each one to answer a
   question Roundfix can already compute.

## Suggested acceptance checks

- Three consecutive failing QA gates followed by a pass leave no branch that
  blocks a review Run.
- `reconcile --apply` releases superseded QA-report branches and preserves
  genuinely unintegrated work.
- No next action suggests an integration that would overwrite a newer artifact.

## What worked — keep

- Reconcile refusing to delete unintegrated work is correct and should not be
  weakened to fix this. The answer is a new classification, not a laxer default.
- Branch Integrity Preflight blocking the review Run was right: it prevented a
  Run from starting on top of unresolved state.

## Addendum — 2026-07-28 — Routed to Spec 0053

The `superseded` classification, its `--apply` release, the
supersession-aware next actions, and the Preflight guard against silently
fast-forwarding a QA-report-only branch are owned by
[Spec 0053 — QA gate reachability and verdict semantics](../specs/0053-qa-gate-reachability-and-verdict-semantics/_prd.md).
