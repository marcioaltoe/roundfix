---
status: accepted
created_at: 2026-09-28T00:00:00Z
updated_at: 2026-09-28T00:00:00Z
deprecated_at: null
superseded_by: null
---

# A blocking review after archive parks publication for a corrective Spec

The maintainer confirmed on 2026-09-28 that an archived Spec is never edited
to absorb a finding. Publication parks as
`corrective-spec-required: <slug>`, and the correction is a new corrective
Spec with its own `_authorization.md` and QA gate.

No Run budget, corrective-Task ceiling, or queue grant authorizes that Spec.
Roundfix never authors or starts it. How the corrective work reaches the
parked candidate is a recorded limit for a future Spec.

This decision refines [ADR-0153](0153-pre-pr-review-is-an-explicit-provider-policy.md),
which requires an enabled review result before publication, and
[ADR-0120](0120-retired-documentation-lives-under-one-history-root.md), which
keeps retired Specs under the history root. It supersedes neither decision.
