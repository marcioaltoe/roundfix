---
status: accepted
created_at: 2026-09-30T00:00:00Z
updated_at: 2026-09-30T00:00:00Z
deprecated_at: null
superseded_by: null
---

# A Task commit is authorized by the grant it ran under

The QA gate's mechanical authorization audit read each Task commit's grant at
`merge-base <delivery target> <commit>^1`, which is the point where the item
branch forked from the default branch. That is a topological proxy for "the
grant that existed before the commit". It misses the grant a Task actually ran
under, when that grant reached the item branch after the fork.

On 2026-09-29 Spec 0181's grant was widened on main (#277) and the same record
was cherry-picked onto the item branch before task_07 ran. The audit still read
the grant at the fork, `6784210b`, and refused the task's commit. Merging main
into the item afterwards did not change that commit's fork point. Only
rebasing the item onto main did.

The audit now also accepts the grant recorded in the Task commit's parent,
when the delivery target's current tip holds a byte-identical authorization
record at the same path. That grant existed before the commit, because the
commit's parent carries it. It is not a local invention, and it has not been
narrowed or revoked since, because the default branch holds exactly that
record now. An older identical record in history is never enough. The fork-point read stays the first choice.

## Consequences

Spec 0119's rule that a later amendment never authorizes an earlier change
still holds, because the grant must already be in the commit's parent. The
self-approval check is unchanged: a Task commit that changes the authorization
record is refused. A record that exists only on the item branch, or that
differs from every version on the default branch, authorizes nothing new. The
audit reports which revision authorized the commit in its `Revision` column.
