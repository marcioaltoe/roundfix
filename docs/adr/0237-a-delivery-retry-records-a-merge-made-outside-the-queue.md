---
status: accepted
created_at: 2026-10-05T00:00:00Z
updated_at: 2026-10-05T00:00:00Z
deprecated_at: null
superseded_by: null
---

# A Delivery Retry records a merge made outside the queue

On 2026-10-05 the operator fixed CI on Spec 0231's item branch and merged its
Pull Request #404 by hand. The Delivery Queue still showed the item parked
`checks-failed`, and no command could answer its Pending Question: the owner
skips parked items, `deliver status` reads only the Run Database, and a
Delivery Retry asks for the item branch and then applies the archived-head
rules of ADR-0223 and ADR-0229, which can send the item back to review and
publication for work that is already merged. Specs 0225 and 0229, whose Pull
Requests the operator opened and merged by hand, stayed parked the same way.

A Delivery Retry of a parked item now first asks whether the item is already
merged, before the Delivery Queue Limits, the item workspace and every other
retry rule. The item is merged when its recorded Pull Request is merged from
the item branch, or, for any item without that answer, when the local default
branch holds the Spec's archived `_prd.md` through a delivery commit that is
not reachable from the item head, the merge evidence of ADR-0232. The retry
then records the item `merged` with that merge commit and the merged head as
its newest candidate, without changing its retry count, and hands the queue
to its owner, whose next pass runs the normal post-merge cleanup and releases
the items that wait for it. A recorded Pull Request that was closed without
merging and has no merge evidence keeps the item parked, and the retry is
refused with that reason.

## Considered Options

- Let `deliver status` ask GitHub: rejected, because status is a read-only,
  offline account of the Run Database and would start changing state.
- Let the owner probe every parked item on each pass: rejected, because a
  park is a Pending Question that only an explicit Delivery Retry or a new
  queue answers, and the probe would add a GitHub read per parked item to
  every pass of every queue.
- Mark the item merged only from the recorded Pull Request: rejected, because
  Specs 0225 and 0229 were merged through Pull Requests the queue never
  recorded.

## Consequences

A retry of an item whose recorded Pull Request cannot be read is refused,
because the retry would otherwise republish work that may already be merged.
Without a recorded Pull Request the retry reads only local Git, so the local
default branch must already hold the delivery, as it must for ADR-0232. A
merge recorded this way is not refused by a retry limit, a queue deadline or
a Queue Token Ceiling, because it resumes no work; this refines the retry
refusal of ADR-0199 for that one case, and ADR-0223 and ADR-0229 keep their
rules for every item that is not merged. A cleanup that cannot prove
the merged head is reported as today's cleanup warning.
