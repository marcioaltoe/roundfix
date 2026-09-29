---
status: accepted
created_at: 2026-09-29T00:00:00Z
updated_at: 2026-09-29T00:00:00Z
deprecated_at: null
superseded_by: null
---

# A Task already completed on its target is nothing to carry

Task Carry-Forward proved every settled Task of a Run against its target, even
a Task the target had already completed. A carried Task's own file changes on
the target, so re-inspecting it always found a moved input, and the Run's
whole set was refused. On 2026-09-28 a Delivery Retry re-carried the older of
two Runs of Spec 0175 for exactly this reason, while the newer Run's proved
Tasks stayed behind.

Now a Task whose status is `completed` on the carry-forward target is recorded
as already completed and is not staged, proved or refused. The remaining Tasks
still have to pass every proof as one whole set. The same rule serves the
Reconcile Command's `--carry-forward`, the Implement Command's Preflight and a
Delivery Retry.

## Consequences

This is not carrying part of a refused set, which Spec 0173 excluded. A
completed Task is outside the set: the Daemon wrote its status when it settled
on the target or was carried there
([ADR-0057](0057-daemon-exclusively-owns-implement-task-status.md)), so there
is nothing left to prove or to carry. A refusal of any remaining Task still
refuses every remaining Task.

A Delivery Retry considers every terminal Implement Run of the item's Spec on
the item branch, newest first. The Run recorded on the item is provenance, not
the selector. The newest Run started from an item head that already held the
older carried work, so once it is carried, the older Runs' Tasks read as
already completed. The retry rereads the item's Task Graph after each carried
Run, because a repository fact is never reused across a mutation
([ADR-0090](0090-repository-facts-are-read-in-batches-never-cached-across-mutations.md)).
Each Run's set is carried whole or refused whole. A later Run's refusal leaves
earlier Runs' proved Tasks on the item branch and names them.

This refines
[ADR-0053](0053-terminal-run-worktree-reconciliation-is-proof-based.md) and
[ADR-0158](0158-a-budget-expired-run-settles-budget-exceeded-and-stays-recoverable.md)
without superseding either.
