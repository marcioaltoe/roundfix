---
type: fix
status: promoted
created: 2026-09-28
spec: 0182-delivery-that-reviews-and-retries-from-where-the-item-stands
reason: null
---

# The delivery pre-PR review diffs the candidate against the current main, not its merge base

## Symptom

Spec 0175's delivery review reported two false findings: that the candidate deletes active Specs 0177–0180 and reactivates retired Specs 0126/0127. Both came from commits that reached `main` (#265, #266) after the item branch was created; the review compared the current default branch with the candidate head, so main's newer commits read as reverted by the candidate.

## Where

The pre-PR reviewer invoked by `internal/delivery/engine.go` `reviewCandidate` (and `roundfix review --base`, which defaults to the main branch): the diff base.

## Expected

The review diffs from the merge base of the candidate and the default branch (three-dot semantics), so only the candidate's own changes are reviewed.

## Evidence

`pre-pr-review.json` for head `74a51f23`, 2026-09-28.
