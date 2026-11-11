---
status: accepted
created_at: 2026-11-11T18:00:00Z
updated_at: 2026-11-11T18:00:00Z
deprecated_at: null
superseded_by: null
---

# A conflicting merge of the default branch is a review delta from its Conflicted Carry

When the candidate's own change does not carry onto a newer default branch
without conflict, the pre-PR review measures the candidate from the tree Git
writes for that carry, with conflict markers in the conflicted files, instead
of from the candidate commit. A merge of the default branch keeps the
Reviewer Lineage whether or not it conflicts, and every review prompt asks the
reviewer to read the whole file before a finding says that something outside
the diff is missing.

## Context

On 2026-11-11 Conexus reported two false findings from round 1 of the pre-PR
review on Roundfix 0.138.0 (provider `codex`): "the added `try` block has no
`catch` or `finally`" and "`StatementBalanceMovements` is referenced without
being imported", each about a file whose `catch` or import lay outside the
hunk and whose head type-checked. Both came right after the operator merged
`main` into the item branch to clear a `pull-request-conflict` park, as the
park asks.

A reproduction with temporary repositories proved the cause. ADR-0366 keeps
the lineage only when the reviewed head carries onto the new merge base
without conflict, and otherwise measures from the earlier commit. A
`pull-request-conflict` park means exactly that the reviewed head conflicts
with the default branch, so every recovery from it restarted the lineage at
round 1 (`the merge base moved from <prior base> to <base> and the reviewed
head <prior head> conflicts with it`) and sent the whole candidate again. In a
segmented review the same fallback measured a segment that spans the merge
from an item commit, and that segment carried the default branch's newly
landed work as the candidate's. The prompt let the reviewer open files but
never asked it to confirm an absence it inferred from a hunk.

The maintainer approved this Spec on 2026-11-11 as "review pré-PR revisa só o
delta e com contexto depois de merge da main". The answer of 2026-11-05 that
ADR-0366 rests on limited the delta to "um merge do main sem conflito"; this
answer covers the merge the conflict park asks for.

## Decision

**A Conflicted Carry is a measurement start.** When `git merge-tree
--write-tree <commit> <Effective Base>` exits 1, Roundfix reads, with
`merge.conflictStyle=merge` so a user's setting cannot change the tree, the
tree it writes and the paths it lists as conflicted. Every diff the review
measures inside the candidate, a round-2 delta, a Correction Review delta and
a Review Segment, starts from that tree. In conflicted files the diff shows
the conflict markers removed and the candidate's resolution; elsewhere it
shows only the candidate's own change. Any other exit status still blocks the
review.

**A conflicting merge keeps the lineage.** The lineage rules of ADR-0366 keep
every other condition: the head descends from the reviewed head, the new
merge base descends from the recorded one, and every merge on the head's
first-parent path brings in only commits the new merge base contains. The
restart reason for a conflicting carry is retired; seven reasons remain. The
two-round ceiling of ADR-0197 is unchanged, so a conflicting merge after a
clean round 2 closes at the ceiling without a review, as any later change
does.

**The prompt says so.** After the line that names the delta's tree, and after
a segment's `This segment's diff starts from` line, the prompt adds
`Carrying <commit> onto the default branch conflicted; where tree <tree>
holds Git's conflict markers, the diff shows how the candidate resolved the
conflict.` and, when Git lists paths, `Conflicted paths: <path>, <path>.`

**An absence is confirmed in the file.** Every review prompt adds, after its
read-only line: `A diff hunk shows only the lines near a change. Before a
finding states that an import, declaration, symbol or closing part of a block
is missing or undefined, read the whole file in the repository, which holds
the head commit, and raise the finding only if the file confirms it.`

## Considered Options

- **Keep restarting on a conflict (ADR-0366).** It sends the whole candidate
  again after every conflict park and resets the two-round ceiling, the cost
  Conexus measured.
- **Measure from the merge commit's own parents.** It sees only the merge
  commit and loses the commits before and after it, which ADR-0366 already
  rejected for `--remerge-diff`. The Conflicted Carry applies the same idea,
  a diff from Git's re-merge with conflict markers to the recorded result,
  over the whole range from the reviewed head.
- **Send the whole of every touched file.** It spends the 917504-byte bound
  of ADR-0282 on unchanged lines, while the read-only session can already
  open any file.
- **Run the repository's typecheck before the review and give its result to
  the reviewer.** The review runs no repository command and serves every
  stack; the Repository Verification already builds the head before the Pull
  Request, and the false findings were claims the file at the head refuted.
- **Keep the lineage but show the conflicted files whole.** The conflict
  markers already show both sides and the resolution in the diff the review
  bounds.

## Consequences

- An item that merges a conflicting default branch, the recovery of every
  `pull-request-conflict` park, gets round 2 or its Correction Review of the
  delta, not a new full review.
- A segmented review no longer shows the default branch's work in a segment
  that spans a conflicting merge.
- `deliver status` and `deliver retry` print no `Review:` line for a
  conflicting merge; the `review` guide lists seven restart reasons.
- Every review prompt grows by one sentence, and a conflicting delta by one
  or two lines.
- A conflict resolution made after a clean round 2 is not reviewed; the
  ceiling of ADR-0197 already bounds every change after round 2.

This decision amends
[ADR-0366](0366-a-conflict-free-merge-of-the-default-branch-is-a-review-delta-and-a-segment-is-measured-from-its-effective-base.md)
and
[ADR-0197](0197-a-pre-pr-reviewer-lineage-spans-at-most-two-rounds.md), and
keeps
[ADR-0282](0282-the-history-cleanup-names-a-superseded-adr-in-its-successor-and-a-large-candidate-is-reviewed-in-segments.md)
and
[ADR-0286](0286-a-fixed-round-2-finding-gets-one-correction-review-of-its-fix.md).
