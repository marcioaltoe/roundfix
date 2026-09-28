---
type: feat
status: promoted
created: 2026-09-28
spec: 0180-a-prepared-queue-that-revalidates-before-each-spec
reason: null
---

# Queue-wide limits and a single pending question

## Opportunity

Spec 0127 Core Features 4 and 5: approval and limits accompany each Spec, but queue-wide time, spending, concurrency and correction limits are not decided or enforced, and nothing guarantees only one maintainer question is pending (silence or a recommended option must never count as an answer).

## Value

An unattended queue stops at an explicit limit instead of an inferred one, and the maintainer faces one decision at a time.

## Shape

Declare queue limits in the queue record; stop dependent actions on a missing or changed grant; one durable pending-question record.

Evidence: carried from Spec 0127 when that portfolio Spec was retired on 2026-09-28 (its delivered features shipped in the Specs its supersession names).
