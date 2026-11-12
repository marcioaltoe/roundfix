---
status: accepted
created_at: 2026-11-05T12:00:00Z
updated_at: 2026-11-11T18:00:00Z
deprecated_at: null
superseded_by: null
---

# A conflict-free merge of the default branch is a review delta, and a segment is measured from its Effective Base

The pre-PR review measures what the candidate itself changed. A merge of the
default branch into the candidate keeps its Reviewer Lineage when the merge is
conflict-free, and every diff the review measures inside the candidate starts
from the Effective Base of the commit it ends at, so the default branch's own
work is never reviewed, counted or shown as reverted.

## Context

On 2026-11-05 Fiscus reported two costs of the same mechanism. Recovering from
`review-findings`, the operator merged the default branch into the item before
the retry; the merge moved the merge base, so the next review opened a new
Reviewer Lineage at round 1, reread the whole candidate and raised new
findings, which happened to Specs 0066 and 0074 on one day. Spec 0070's review
was blocked by size: segmentation by first-parent commits measured an early
item commit, cut from an old default branch, against the new merge base, and
that range showed about 1.5 to 2 MB of the default branch's work as reverted.
A Delivery Retry accepts only a head that descends from the candidate, so the
history could not be redone; the workaround reapplied the change as linear
commits and closed with a tree-identical merge. ADR-0169 already says the
candidate's own merge of the default branch moves the merge base, and
ADR-0197 starts a new lineage whenever it moves. The maintainer answered on
2026-11-05:

> Um merge do main sem conflito, em que o único pai novo é o main, conta como
> delta e não reinicia o review na rodada 1. Quando o review for reiniciar, o
> retry avisa antes.

> Cada trecho do review passa a ser medido contra o merge-base com o main
> atual. Um commit antigo que só parece reverter o main deixa de bloquear o
> review por tamanho.

## Decision

**The Effective Base of a commit** is its merge base with the review's merge
base, which is also its merge base with the current default branch: the
commit of the default branch it was cut from or last merged. A diff from a
commit `S` to a later commit `C` of the candidate starts from:

- `S` itself, when both share one Effective Base;
- the Effective Base of `C`, when `S` is on the default branch and has no
  change of its own;
- otherwise the tree `git merge-tree --write-tree S <Effective Base of C>`
  writes, `S`'s own change carried onto `C`'s Effective Base. When that merge
  conflicts, the diff starts from `S`, as before.

**Review Segments are measured from it.** Every diff the segment plan
measures, the own diff of a commit and the diff of a growing segment, starts
as above. A segment whose diff starts elsewhere than its base records
`measuredFrom` and its prompt says so in one line. A commit cut from an old
default branch therefore shows only its own change, and a conflict-free merge
of the default branch adds nothing. A single commit, a fifth segment, or a
Claude prompt that still exceeds its bound blocks as before.

**A conflict-free merge of the default branch keeps the lineage.** When the
merge base moved since the recorded review, the next review continues that
Reviewer Lineage when all of these hold: the head descends from the reviewed
head, the new merge base descends from the recorded one, every merge on the
head's first-parent path since the reviewed head brings in only commits the
new merge base contains, and the reviewed head carries onto the new merge base
without conflict. The round-2 delta, and a Correction Review's delta, then
start from that carried tree, recorded as `lineage.deltaBase`, and the prompt
says so in one sentence. Any other move of the merge base, a rebase among
them, starts round 1 as before.

**Amended by [ADR-0391](0391-a-conflicting-merge-of-the-default-branch-is-a-review-delta-from-its-conflicted-carry.md) (2026-11-11).** A carry that conflicts no longer measures from the earlier commit or restarts the lineage: the diff starts from its Conflicted Carry, the tree `git merge-tree --write-tree` writes with Git's conflict markers, and a merge of the default branch keeps the lineage with or without a conflict. The other lineage conditions, the Effective Base and the `Review:` line still hold.

**The queue says so before a restart.** When the item's current head will
open a new Reviewer Lineage at round 1 although a review of an earlier head is
recorded, `roundfix deliver status` prints, for a parked or reviewing item, and
`roundfix deliver retry` prints, when the item returns to `reviewing`, one
`Review:` line naming the head and the reason, before the review runs.

## Considered Options

- **Count only a merge whose new head's tree equals a fresh merge of the
  reviewed head and the default branch.** It misses the common case, in which
  the fix commits precede the merge, and the carried tree already gives the
  same answer for a pure merge: an empty delta.
- **Show the merge's `--remerge-diff`.** It works per merge commit only, needs
  Git 2.36, and leaves the commits before and after the merge to be combined
  by hand. ADR-0391 takes this option's idea, a diff from Git's re-merge with
  conflict markers, over the whole range from the reviewed head.
- **Raise the review bound or drop the oversized commit.** ADR-0282 measured
  why the bound exists, and dropping a commit leaves its change unreviewed.
- **Let the operator rewrite the history.** A Delivery Retry accepts only a
  head that descends from the candidate, and rewriting discards the head the
  review and the checks were recorded against.
- **Keep the lineage across a merge that also brings in another branch.** That
  branch's work was never reviewed; the maintainer's answer limits the delta
  to a merge whose only new parent is the default branch.

## Consequences

- An item that merges the default branch after `review-findings` gets round 2
  or its Correction Review, not a new full review; the ceiling of ADR-0197
  still bounds it.
- A segmented review of a candidate that merged the default branch no longer
  blocks on a range that only reverts it, and its segments' bytes drop by the
  default branch's work.
- A conflict in the carry keeps today's behavior: the segment is measured from
  its base, and the lineage restarts with a reason that names the conflict.
- `deliver status` and `deliver retry` may print one more line per item; their
  other lines and exit codes are unchanged.

This decision refines
[ADR-0169](0169-the-pre-pr-review-diffs-the-candidate-from-its-merge-base.md),
whose merge base still moves with the candidate's own merge. It amends
[ADR-0197](0197-a-pre-pr-reviewer-lineage-spans-at-most-two-rounds.md),
[ADR-0282](0282-the-history-cleanup-names-a-superseded-adr-in-its-successor-and-a-large-candidate-is-reviewed-in-segments.md),
[ADR-0286](0286-a-fixed-round-2-finding-gets-one-correction-review-of-its-fix.md)
and
[ADR-0352](0352-the-pre-pr-review-omits-declared-paths-and-a-large-claude-candidate-is-reviewed-in-segments.md).
