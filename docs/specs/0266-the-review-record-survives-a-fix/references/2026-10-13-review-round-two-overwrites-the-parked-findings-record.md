---
type: fix
status: promoted
created: 2026-10-13
spec: 0266-the-review-record-survives-a-fix
reason: null
---

# A second review round overwrites the parked findings record

## Symptom

Running `roundfix review --base origin/main` on a parked item after a fix
writes its result over the record that held the first round's findings and
dispositions. The correction proof then has nothing to compare against.

## Where

The pre-PR review record is keyed only by checkout, not by head or round, under
`pre-pr-review/<checkout-key>/pre-pr-review.json`.

## Expected

Each round keeps its own record, or the newer round keeps the earlier findings
and their dispositions, so the correction proof can still read them.

## Evidence

Delivery Queue item 0259 on 2026-10-13 (operator log entry 308).
