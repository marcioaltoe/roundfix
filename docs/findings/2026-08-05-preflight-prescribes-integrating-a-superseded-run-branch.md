# 2026-08-05 — Preflight prescribes integrating a superseded Run Branch

status: pending

## What was observed

A Spec 0015 Run failed on four Tasks, leaving a Run Branch with one commit —
task_03's work. The Spec was relaunched, the same Task was redone, and the
successful Run integrated its own version onto the Spec branch.

Starting the review then refused:

```
Branch Integrity Preflight refused pending Run Branch work for PR Head Branch
"ma/0015-catalog-sync-operability".
- branch=roundfix/run-run_20260804T181023Z_f9fef8003e217a60 ahead_commits=1
  integration_command="git merge --ff-only roundfix/run-..."
Next action: inspect each pending Run Worktree, then run the listed integration
command from the repository root when it is safe.
```

Running that command would have been wrong. Both commits carry the same
subject — `feat: carry a quarantine for either subject, with reason and
attempts` — and the surviving one is 272 lines ahead in the repository
implementation, having passed through three later Tasks and a QA gate. The
prescribed fast-forward would have re-applied a superseded artifact over newer
work.

`roundfix reconcile` could not release it either: the stranded Run Worktree has
changes, so it classifies `dirty` and preserves, and `--apply` acts only on
`safe` or `superseded`. The only supported way forward was
`--skip-branch-integrity`, which published its audit comment and let the review
proceed.

## Root cause

Preflight detects pending Run Branch work correctly and then prescribes a
single remedy without testing whether the pending work is still wanted. The
`superseded` classification already exists in Reconcile's vocabulary — it is
defined for terminal Runs whose commits are QA-report-only and whose target
carries a newer report — but nothing applies that reasoning to a failed Run
whose Tasks were redone.

The result is a prescription that is safe in the common case and destructive in
the case that produced it. An operator following the printed next action
faithfully would have damaged the branch, and the only thing that prevented it
here was the skill's own warning to verify a stranded branch before discarding
it — guidance that lives in prose, not in the command.

`dirty` compounds it: a failed Run almost always leaves a dirty Worktree, so
the branch it stranded is exactly the one Reconcile refuses to release.

## What would settle it

Extend supersession to redone Tasks. When a pending Run Branch's commits
correspond to Tasks that a later Run settled `completed` on the target branch,
classify it `superseded` and let `reconcile --apply` release it, recording the
superseding commits as evidence — the same shape already used for QA reports.

Failing that, two cheaper fixes:

- **Do not print a bare integration command for work that may be superseded.**
  Name the risk in the refusal: the branch may contain a Task the target has
  already redone, and integrating it would overwrite newer work.
- **Let `--apply` release a `dirty` Worktree whose branch is proven
  superseded.** Dirtiness is evidence about the Worktree, not about whether the
  commits are still wanted.

## Related

[[2026-08-04-review-runs-halt-autonomous-delivery-on-unrelated-dirty-files]]
and [[2026-07-28-failed-qa-runs-strand-branches-that-block-review-runs]] are
the same preflight refusing delivery for reasons the loop cannot resolve on its
own.

## Spec pointer

None yet.
