---
type: feat
status: promoted
created: 2026-09-28
spec: 0179-review-findings-with-evidence-and-no-unselected-providers
reason: null
---

# A blocking review after archive parks publication and asks for a new Spec

## Opportunity

Spec 0127 Core Feature 9: when a pre-PR review blocks after the Spec is archived on the branch, the correction path is ad hoc.

## Value

The archived Spec stays intact, publication parks, and the correction is authored as a new Spec with its own authority and gate; the candidate is reviewed again after it archives.

## Shape

A `review-blocked` park reason that names the needed corrective Spec; no cycle budget silently authorizes it.

Evidence: carried from Spec 0127 when that portfolio Spec was retired on 2026-09-28 (its delivered features shipped in the Specs its supersession names).
