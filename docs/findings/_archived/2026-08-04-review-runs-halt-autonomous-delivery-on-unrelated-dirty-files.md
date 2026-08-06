---
status: done
created_at: 2026-08-04
updated_at: 2026-08-06
absorbed_by: 2026-08-06-rollup-review-and-delivery-convergence.md
---

# 2026-08-04 — Review Runs halt autonomous delivery on unrelated dirty files

## What was observed

An autonomous delivery of four fluxus Specs stopped dead at the review step of
the first one, for a reason that had nothing to do with the work.

`roundfix watch --source coderabbit --pr 78 --until-clean` refused Preflight:

```
review Run Preflight Validation refused because tracked files are dirty in
checkout /Users/marcio/dev/fluxus.
- M .agents/skills/archive-spec/SKILL.md
- M .agents/skills/qa-gate/SKILL.md
  … 8 files, 502 insertions
Next action: stash or commit tracked changes before running Roundfix
```

Those eight files were modified **before the session began**, are unrelated to
PR #78, and are Skill contracts under `.agents/` — protected tooling the
Supervisor has no bounded-file authorization to commit, stash, or discard on
the user's behalf. Every available action belonged to the human.

The Spec itself was finished: QA `verdict: pass`, 113 rows passing, archived,
pushed, PR open with 5/5 checks green. The loop had nothing left to decide and
still could not proceed.

The asymmetry is the point: `implement` explicitly tolerates this exact state.
Its preflight prints a note — *"working tree has N uncommitted change(s);
implement will run in a Run Worktree"* — and proceeds, because spec Runs are
worktree-isolated. Review Runs execute in the user's checkout and therefore
inherit whatever is in it.

## Root cause

Review Runs need the checkout because a review fix must be a delta over the PR
head branch that Final Push updates. That requirement is real, but it makes
the *cleanliness of unrelated files* a precondition of delivery. Any long-lived
uncommitted work in the repository — an in-progress Skill edit, a scratch
config, a half-finished refactor — becomes a hard stop for an autonomous loop
that is otherwise complete.

For the stated goal of "the human writes the PRD and the agent does the rest",
this is the failure mode with the worst ratio: zero decisions required, full
stop, and the human must context-switch back into a repository they had
deliberately handed over.

## What would settle it

Give review Runs the same isolation spec Runs already have. A review Run can
create a worktree on the PR head branch, apply fixes there, and push from it —
the delta-over-head property is preserved, and the user's checkout is never
read or required to be clean.

If checkout execution must stay, two cheaper mitigations:

- Scope the cleanliness requirement to paths the Run will touch. Files outside
  the PR's diff and outside the review's fix surface cannot conflict with a
  Batch commit that stages only paths changed since its snapshot.
- Offer an explicit, opt-in auto-stash bounded to the offending paths, restored
  on terminal outcome, so the operator can authorize it once instead of being
  asked to hand-manage it mid-loop.

The general principle worth stating: an autonomous loop should never be
blocked by repository state it is not authorized to change. Either it is
isolated from that state, or the requirement is narrowed to what it owns.

## Related

[[2026-08-04-what-still-needs-a-supervisor-between-a-prd-and-a-merge]] enumerates
the structural stops between a PRD and a merge. This is one more, and the only
one that requires no decision at all — which is what makes it the cheapest to
remove. [[2026-07-28-failed-qa-runs-strand-branches-that-block-review-runs]]
records a different route to the same halt.

## Spec pointer

None yet.
