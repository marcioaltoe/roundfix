---
type: fix
status: promoted
created: 2026-10-04
spec: 0228-queue-items-that-stay-current-with-main
---

# A documentation fix after archive forces a hand-opened Pull Request

## Problem

On 2026-10-04 the pre-PR review of 0225 raised two findings on the archived
Spec's own records: a Task Result that named the superseded skill version,
and an override reason the reviewer misread. The operator fixed the first
in one docs-only commit and dismissed the second with evidence. `deliver
retry` then refused with `corrective-spec-required`: "archived Specs … were
reviewed at parked candidate head …, but the item head is …; author a
corrective Spec". The operator ran the round-2 review by hand (`reviewed`)
and opened PR #382 by hand (intervention log entry 158). Spec 0219's
descendant rule covers a correction after an archived *candidate*, but not
one made in reply to the review of that candidate.

## Expected

A correction that descends from the reviewed candidate, and whose findings
are all disposed (fixed by a descendant commit or dismissed with evidence),
returns the item to review for its second round, as the manual route did,
without a corrective Spec. A change outside the archived Spec's records, or
an undisposed finding, still requires one.

## Notes

Non-binding. ADR-0165 parks a blocking review after archive and ADR-0197
bounds the review to two rounds; the second round is the natural place for
this.
