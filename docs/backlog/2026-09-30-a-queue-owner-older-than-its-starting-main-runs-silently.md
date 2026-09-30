---
type: fix
status: open
created: 2026-09-30
spec: null
reason: null
---

# A queue owner older than its items' starting main runs without warning

## Symptom

The Wave 5 queue owner was built at `c98b1641`. Spec 0181 then merged, and items started from a main that included it. The owner kept the defects 0181 fixed for the whole queue: the unfixed ADR cascade, scope refusals for undeclared paths and the unclassifiable review verdict. That caused four parks. Nothing in `deliver start` or `deliver status` said that the owner binary was older than the main its items started from.

## Where

The Delivery Queue owner: `deliver start` and `deliver status`, and the item start in the delivery engine. The binary's build commit is available through `internal/app`, and the Auditor Staleness evidence already compares a build commit with a base.

## Expected

When an item starts from a main whose history holds commits that change Roundfix's own source after the owner's build commit, and the repository is Roundfix itself (its module path), `deliver status` records one warning naming the owner build and the starting main. The warning only reports; it never stops the queue.

## Evidence

`docs/findings/2026-09-29-the-first-full-queue-trial-needed-manual-recovery.md`, Wave 5 addendum.
