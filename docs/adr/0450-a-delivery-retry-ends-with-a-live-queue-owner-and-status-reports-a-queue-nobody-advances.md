---
status: accepted
knowledge: true
created_at: 2026-11-24T00:00:00Z
updated_at: 2026-11-24T00:00:00Z
deprecated_at: null
superseded_by: null
---

# A Delivery Retry ends with a live Queue Owner, and status reports a queue nobody advances

## Context

On 2026-11-23 quaestor reported that a Delivery Retry left an item at
`running` for six hours with no owner and no Active Run. Status printed its
row, limits and usage with no warning. The preceding Implement Run ended
Failed after only its QA gate failed; its completed implementation Tasks were
not carried forward. The item's recorded Run also lagged behind its executor.
The Praetor reported that status exposed neither the Run nor the Tasks in
progress, so its periodic sweep falsely reported no progress on another
project's running item.

## Decision

This decision amends [ADR-0279](0279-status-reports-a-queue-owner-proven-gone-and-a-dismissal-survives-a-fix-that-leaves-its-anchor.md)
and [ADR-0395](0395-a-queue-item-stops-only-for-its-own-work.md) in part.

The Implement executor returns the Run its child created on every exit,
including Failed. The Queue Owner records that identifier on the item before
handling the result, with a context its own stop cannot cancel. A store error
is joined with the executor error; the ensuing park keeps its blocker and
Park Class.

A QA-Only Failed Run is a terminal Failed Implement Run whose Run Branch
exists and holds a loadable committed Task Graph with the QA gate's Task not
completed and every other Task completed. Delivery Retry accepts it under the
existing Task Carry-Forward proofs, removing the QA gate's evidence before
inspection and refusing the whole set on a proof failure. A missing branch or
unloadable graph is skipped. The Reconcile Command and Implement Preflight
keep their accepted outcomes unchanged.

After storing the re-entry stage, Delivery Retry hands the item to a proven
live recorded owner or starts one and reads the queue again. A recorded owner
proven gone is reclaimed; an uncheckable live identity fails with its existing
diagnostic. Success requires a proven live owner unless every item is already
merged or parked. Success prints the merged-outside notice, carry-forward
lines, `Retried`, review notice and buffered owner report in their existing
order, and exits `0`.

Failure prints no standard output and exits `2`. When the start failed, no
owner is recorded, or its recorded owner is proven gone, it releases a record
proven gone and restores the parked item in one guarded transaction. The guard
requires no recorded owner and the item's re-entry stage with an empty
blocker; it restores stage, blocker, retry count and candidate commits, keeping
the Run, Item Branch, worktree and Run links. A failed guard changes nothing.
An uncheckable live identity changes nothing more, because that owner counts
as present. The failure report names the reason, stored item stage and blocker
(`-` when empty), any carried Runs, and `roundfix deliver retry <slug>` when
reparked or `roundfix deliver status` otherwise. Reasons distinguish a failed
start, an absent PID, a different identity, no owner after start and an identity
that could not be checked.

When an item is neither `merged` nor `parked` and no owner is recorded, status
prints `Owner: Delivery Queue has no recorded owner and <n> item(s) can still advance; run roundfix deliver resume`, counting those items. Existing
proven-gone owner lines and uncheckable-owner diagnostics keep their bytes;
a proven live owner or a queue with no advanceable item adds no owner line.
Status stays read-only and reclaims nothing.

For each item at `running`, in queue order after `Turn:` and before `Park:`,
status prints one `Run:` line independently of owner proof. Its Active Run is
the newest non-terminal Implement Run for the repository, Spec and Item Branch,
read with `ListRuns`. Tasks in progress come from `RunEventsOfKinds`: a Task's
latest `started` event has no later `settled` or `skipped`, ordered by their
starts. The forms are `Run: <slug> no Active Run`,
`Run: <slug> <run-id>; Tasks in progress: <task>, <task>` and
`Run: <slug> <run-id>; no Task in progress`. The QA gate has no start event and
uses the last form. Database read errors fail status. Existing lines keep
their bytes and relative order, rows retain four tab-separated fields, and
help, schema, stages, commands, flags and exit codes stay unchanged.

## Considered Options

- Record the Run at its start. Rejected: the child creates it and the owner
  reads its output only after it exits. Earlier recording would need polling
  beside the child; status instead reads the Active Run directly.
- Restart a stopped Queue Owner automatically. Rejected: recovery remains an
  explicit operator action through resume or retry. Status reports the missing
  owner without mutating the queue.
