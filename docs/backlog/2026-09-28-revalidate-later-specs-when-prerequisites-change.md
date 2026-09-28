---
type: feat
status: open
created: 2026-09-28
spec: null
reason: null
---

# Later queued Specs are revalidated when a prerequisite changes their assumptions

## Opportunity

Spec 0127 Core Feature 6: after an earlier queue item merges, later items run against the new main without checking that their authored premises (touched files, signatures, governed paths) still hold — Onda 2 hit this when Spec 0175 was authored before 0173 changed the delivery code it edits.

## Value

A contradiction, a new authority requirement or a new irreversible action becomes a durable blocker instead of an improvised decision in a Run.

## Shape

Re-run `spec check` and the governed-path audit for each later item at its start; park with a named blocker on failure.

Evidence: carried from Spec 0127 when that portfolio Spec was retired on 2026-09-28 (its delivered features shipped in the Specs its supersession names).
