---
status: accepted
created_at: 2026-09-29T00:00:00Z
updated_at: 2026-09-29T00:00:00Z
deprecated_at: null
superseded_by: null
---

# A related-ADR gap opens only for ADRs that predate the Spec

`SC-ADR-RELATED` reports an accepted ADR that cites an ADR a Spec lists, when the
Spec does not list the citing ADR itself (ADR-0093). The candidate is a
suggestion to the Spec's author. A new ADR that cites a widely listed
predecessor therefore opened a gap in every other active Spec. Where gaps are
promoted, each gap refused: the strict check, the Delivery Plan, queue
revalidation and the QA gate. PRs #268, #272 and #274 each had to edit other
Specs' Project Constraints before they could merge. None of those Specs could
have accounted for a decision that did not yet exist when it was written.

The check now has a horizon. For a Spec whose `_prd.md` is committed, an ADR is
related only when the commit that added the ADR is an ancestor of the commit
that added the PRD. An ADR added later, or not committed at all, is outside that
Spec's horizon. The full check applies when the PRD is uncommitted, when the ADR
and the PRD sit in different repositories, and when history cannot be read. That
covers the Spec being written and a shallow history, and in both cases the
check fails closed.

## Consequences

The horizon narrows only `SC-ADR-RELATED`. `SC-ADR-UNLISTED` and
`SC-CITATION-UNSUPPORTED` still read everything a Spec cites, whenever the cited
ADR arrived (ADR-0116). The check still reads citations and never infers
relevance (ADR-0093), and it still skips what is absent (ADR-0094). A PRD
revised after a later ADR landed keeps the horizon of its first commit.
