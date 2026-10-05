---
status: accepted
created_at: 2026-10-04T00:00:00Z
updated_at: 2026-10-04T00:00:00Z
deprecated_at: null
superseded_by: null
---

# Merge evidence releases the Runs and item branches of a merged Spec

ADR-0161 and ADR-0212 release a terminal Run of a merged Spec only when every
commit its Run Branch holds and the merged head lacks is represented there,
and they represent a commit that is neither a Task commit nor Spec-directory
work by comparing its changed files with the merged head. After a squash merge
that comparison measures how far the default branch moved since, not whether
work was lost: on 2026-10-04 reconcile kept six Runs of five Specs archived and
merged on main ("21 differing shared files", "1 Run-only file"), never looked
at two item branches the Delivery Queue had recreated, and the post-merge
cleanup refused every merge whose default branch had moved during the item.

A Spec now carries merge evidence when the default-branch head holds its
archived `_prd.md` and the commit that added that file, its delivery commit,
is not reachable from the Run or item branch being proven. With merge
evidence, a terminal Run of that Spec is `superseded` when each of its Task
commits names the Spec and its Task is completed in the archived Spec; every
other commit, QA Report commits included, is superseded by the delivery
without a content comparison. A `roundfix/deliver-<slug>-<16 hex>` branch that
is not the branch of a live Delivery Queue item is released under the same
proof, with its worktree only when that worktree is clean. Without merge
evidence the content proof of ADR-0161 and ADR-0212 stays exactly as it is.

The post-merge cleanup accepts a merge commit whose tree equals the tree Git
computes by merging the candidate head into the merge commit's first parent
(`git merge-tree --write-tree`), so a merge that holds the candidate's changes
plus later default-branch commits is proven. Any other difference, or a Git
without that mode, still refuses with today's text.

## Consequences

Uncommitted work stays protected: a dirty Run Worktree of a merged Spec is
released only under ADR-0212's rule, when every dirty path lies in the Spec's
directories or in a path its Tasks declared or recorded; any other dirty path
keeps the Run `dirty`, and a dirty item worktree keeps its item branch. A
committed change that never reached the delivery is released with the Run, and
its commits remain reachable by the Run head the Run Database records with the
reconciliation, and by the head the report prints for an item branch, until
Git collects them. Releasing remains the explicit act of ADR-0053 and
ADR-0161: `roundfix reconcile --apply` or the delivery owner after a merge, and
`--discard-superseded` keeps its Branch Disposition (ADR-0115). This refines
ADR-0161 and ADR-0212 for merged Specs without superseding them.
