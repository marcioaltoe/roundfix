---
status: accepted
created_at: 2026-09-30T00:00:00Z
updated_at: 2026-09-30T00:00:00Z
deprecated_at: null
superseded_by: null
---

# A Spec names its prerequisites, and the owner waits for them

On 2026-09-30 the queue owner started Spec 0195 while Specs 0192 and 0193 were
parked. 0195's Verification named tests those two Specs create, so its task_04
failed twice, and the later rebase cost a corrective Task. Queue order was the
only ordering the queue knew, and a parked item never held back the items
behind it.

A Spec now names the Specs it needs in the `requires` list of its Task Graph
manifest, beside the Task dependencies the same file already owns. A
prerequisite is met when its archived Spec folder is on the refreshed default
branch, which is where a squash merge leaves it. The owner never starts an
item while a prerequisite is unmet. While the prerequisite is still ahead in
the same queue the item waits, `queued`. When the prerequisite is parked, or
is in no queue at all, the item parks with a blocker that names it.

The owner itself returns such an item to `queued` once every prerequisite is
met. This is not a Delivery Retry: the item never started, so no retry is
counted, nothing is carried forward and no Run evidence is read. A Delivery
Retry of the item is still accepted and has the same effect.

## Consequences

`deliver start` refuses a queue in which a `requires` entry names no Spec in
the repository, names the Spec itself, or closes a cycle among the queued
Specs, because such a queue could never finish. A binary older than this
decision ignores the field and starts the item in queue order, as before. The
owner releases a dependency park only while it is running; when every item is
parked it exits, and the retry of the prerequisite that restarts it also
releases the dependent item once the prerequisite merges.
