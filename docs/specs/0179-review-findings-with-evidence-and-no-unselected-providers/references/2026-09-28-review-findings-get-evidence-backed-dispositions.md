---
type: feat
status: promoted
created: 2026-09-28
spec: 0179-review-findings-with-evidence-and-no-unselected-providers
reason: null
---

# Pre-PR review findings get evidence-backed dispositions, and a changed candidate gets a fresh review

## Opportunity

Spec 0126 Core Feature 5: today a finding is resolved by the operator (a corrective Task, a TechSpec edit, or a recorded dismissal in the PR body) with no durable per-finding record, and nothing forces a fresh review when the candidate changes after review.

## Value

Every finding's disposition (fixed by commit, dismissed with evidence) is recorded with the candidate, and publication consumes review evidence only for the head it reviewed.

## Shape

Record dispositions in the review artifact; require fresh enabled review on a changed head; decide the archive-first late-correction policy (corrective-Spec authority) before implementing it.

Evidence: carried from Spec 0126 when that portfolio Spec was retired on 2026-09-28 (its delivered features shipped in the Specs its supersession names).
