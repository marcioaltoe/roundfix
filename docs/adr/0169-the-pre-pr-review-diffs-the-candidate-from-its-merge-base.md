---
status: accepted
created_at: 2026-09-29T00:00:00Z
updated_at: 2026-09-29T00:00:00Z
deprecated_at: null
superseded_by: null
---

# The pre-PR review diffs the candidate from its merge base

The pre-PR review compared the tip of the default branch with the candidate
head, so every commit that reached the default branch after the candidate was
cut read as reverted by the candidate. On 2026-09-28 this produced false
findings in four reviews. It also made the delivery engine park an item as
`corrective-spec-required` for a Spec it never touched. Roundfix now resolves
the merge base of the candidate head and the selected base ref once. That
single commit is the diff's before side, the changed-Spec list's before side,
the prompt's `Base commit` and the review record's `baseCommit`. It is the
comparison `git diff A...B` and a GitHub pull request make: what the candidate
introduces.

## Consequences

The record keeps the commit the base ref resolved to as `baseTipCommit`, for
provenance only. A review is reused by repository, merge base, head and
provider. A default branch that moves on therefore no longer invalidates a
review of an unchanged candidate, because its diff has not changed. A
candidate that shares no history with the base ref is refused before any
reviewer call. With several merge bases, Roundfix uses the one
`git merge-base` reports, as `git diff A...B` does. The candidate's own merge
of the default branch still moves the merge base forward, as it does for a
pull request.

This decision refines
[ADR-0153](0153-pre-pr-review-is-an-explicit-provider-policy.md), which
requires a review of the current candidate. It also refines
[ADR-0165](0165-a-blocking-review-after-archive-parks-publication-for-a-corrective-spec.md),
whose corrective-Spec park now fires only for archived Specs the candidate
itself changed. It supersedes neither.
