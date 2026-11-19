---
status: accepted
created_at: 2026-11-18T21:00:00Z
updated_at: 2026-11-18T21:00:00Z
deprecated_at: null
superseded_by: null
---

# A delivery Pull Request behind its base branch is brought up to date by a merge of the default branch

When GitHub reports a delivery Pull Request's merge state as `BEHIND`, the
Queue Owner stops the Check Wait at once and merges the refreshed default
branch into the item branch, the same merge it already makes for a conflict
confined to derived paths. The item then goes back through the repository
gate, the push and the Head Check Set, never through the pre-PR review, the QA
gate or the Run. One owner makes at most three such Base Updates for one item;
when the Pull Request is behind a fourth time, the item parks
`pull-request-behind`.

## Context

On 2026-11-16 quaestor reported (ai-memory message `01a1e522`) that its
`main` now requires branches to be up to date before merging. A delivery Pull
Request that falls behind `main` is refused with "base branch policy prohibits
the merge"; the queue returns to `checking` once for that head, and the second
refusal parks `delivery-error`. The recovery runs the review, the gate and CI
again. The report asked the queue to bring a branch that fell behind up to
date by itself.

A reproduction on `1a8e08b4` confirmed the sequence, at the Delivery Engine
and through the command engine with a strict local forge. The checking stage
waits only while the merge state is `BLOCKED` or `UNKNOWN`, so `BEHIND` counts
as settled and the item moves to `merging`. The only merge of the default
branch the queue makes runs when GitHub reports `mergeable: CONFLICTING`
([ADR-0192](0192-a-conflict-confined-to-declared-derived-paths-is-resolved-by-regeneration.md)).
Called on a branch that is only behind, that merge commits a clean merge, and
it finds nothing to merge when the branch already holds the fetched default
branch.

GitHub documents `BEHIND` as "The head ref is out of date." and tells the
author to merge or rebase the base branch into the branch when branch
protection requires it to be up to date. A published measurement on a real
pull request (multiplex-ai/muggle-ai-works #265) found that `BEHIND` surfaces
only when the base requires up-to-date branches, and that `DIRTY` and
`BLOCKED`, which a pending or failing required check causes, take precedence
over it.

## Decision

**The trigger.** At `checking`, after the `CONFLICTING` test, a merge state of
`BEHIND` stops the Check Wait. A merge state of `BEHIND` never moves an item
to `merging`, as `BLOCKED` and `UNKNOWN` never do.

**The update.** The owner merges the refreshed default branch into the item
branch in the item worktree, through the merge ADR-0192 defines: never a
rebase, never a forced push. A merge without conflict is committed with the
same subject and the trailer `Roundfix-Delivery: derived-merge`. The new head
becomes the item's newest candidate, and the item returns to `gating`, so the
repository gate runs on the merged tree, the head is pushed, and the Head
Check Set of the new head is awaited before the merge.

**What does not run again.** The pre-PR review, the QA gate and the Run do not
run again. The merge adds only commits already on the default branch, and the
item's own change against its merge base does not move, as ADR-0192 holds for
the derived merge and
[ADR-0366](0366-a-conflict-free-merge-of-the-default-branch-is-a-review-delta-and-a-segment-is-measured-from-its-effective-base.md)
holds for the review lineage. `review-stale` still parks only a Pull Request
head the queue did not record as a candidate.

**A conflict.** A conflict confined to declared derived paths is resolved by
regeneration, as in ADR-0192. Any other conflict aborts the merge and parks
`pull-request-conflict: <path>, …`, whose Delivery Retry resumes at
`reviewing` because a hand resolution touches source.

**Nothing to merge.** When the fetched default branch is already in the item
branch, GitHub's merge state lags; the Check Wait goes on polling, and its
timeout parks `checks-timeout: … ; merge state BEHIND`.

**The bound.** One Queue Owner makes at most three Base Updates for one item,
counted in its memory. When the Pull Request is behind again after the third,
the item parks `pull-request-behind: pull request #<n> is still behind its
base branch after 3 Base Updates`, of Park Class `environment`. The count
starts over when the item parks; a Delivery Retry of that park resumes at
`checking`, where the next `BEHIND` updates the branch again. An owner without
a merge of the default branch parks `pull-request-behind: pull request #<n> is
behind its base branch` at the first `BEHIND`.

**The policy refusal.** A merge refused by the branch policy still returns to
`checking` once for that head, where a `BEHIND` merge state now updates the
branch. A second refusal at the same head still parks `delivery-error`.

## Considered Options

- GitHub's update-branch endpoint (`PUT /repos/{owner}/{repo}/pulls/{number}/update-branch`
  with `expected_head_sha`): rejected. The forge makes the merge commit, so the
  item branch in the owner's worktree falls behind its own remote and the
  owner would have to adopt a head it did not make; derived paths cannot be
  regenerated; and the repository gate never sees the merged tree before CI.
  The local merge keeps the item branch the owner's, and a squash merge folds
  it away like any other commit of the branch.
- A rebase onto the default branch: rejected for the reasons ADR-0192 gives.
  It rewrites a pushed branch, needs a forced push and discards the head the
  review and the checks were recorded against.
- Detecting a stale branch from the compare API's `behind_by`: rejected. It
  would also update branches whose base does not require it, paying a gate
  and a CI run for nothing. `BEHIND` is the base branch's own requirement.
  The price is that a required check pending or failing on the behind head
  shows `BLOCKED` first, so the update usually comes after one CI round on a
  head that could not merge.
- Running the pre-PR review again after the update: rejected, because the
  reviewed change of the item does not move.
- Skipping the repository gate and relying on CI: rejected. The gate is the
  queue's own Verification of the tree it pushes, whatever CI an adopter runs,
  and ADR-0192 relies on it as the net after a merge.
- Updating without a bound: rejected. A base branch that moves during every
  gate and CI round would keep an item cycling without end.
- Persisting the count in the Run Database: rejected. The bound only has to
  stop one owner's loop; a schema change buys nothing a retry does not already
  give, and the policy-refusal memory of the same stage already lives in the
  owner.

## Consequences

An adopter whose base branch requires up-to-date branches gets its delivery
merged without a park and without a second review. Each move of the base
branch costs one repository gate and one CI round on the merged head. A
repository whose base does not require up-to-date branches never sees
`BEHIND` and keeps its behavior.

The Base Update commit carries the derived-merge trailer and merges the
default branch, so the Correction Sign check of
[ADR-0422](0422-a-commit-without-a-mark-of-a-run-or-of-the-queue-on-the-item-branch-is-a-correction-sign.md)
never counts it, and a commit hook that refuses it parks `delivery-error` with
the hand merge that recovers it, as
[ADR-0302](0302-a-queue-commit-names-the-spec-by-number-and-a-hook-refusal-parks-with-the-hand-commit.md)
decides for the derived merge. The checking stage of
[ADR-0369](0369-the-checking-stage-waits-for-the-head-check-set-and-parks-on-any-failed-check.md)
keeps its Head Check Set and its parks; it gains one merge state that never
merges. The Check Wait of
[ADR-0311](0311-the-delivery-queues-check-wait-is-a-user-config-duration-sixty-minutes-by-default.md)
keeps its key, default and bounds.
