---
status: accepted
created_at: 2026-10-04T00:00:00Z
updated_at: 2026-10-06T00:00:00Z
deprecated_at: null
superseded_by: null
---

# An operator archive resumes any park from the Run start

On 2026-10-04 the Delivery Queue parked Spec 0220 `run-unresolved` with no
candidate head. Its newest QA Report was `partial` with no finding-blocked
row; its two environment-blocked rows waited for an open Pull Request. The
operator archived it with the QA Archive Override, as the standing
authorization allows, and `roundfix deliver retry` refused with
`candidate head is missing`. Only a `qa-environment-partial` park could use
the Implement start head of its Run in place of a candidate, so the operator
opened the Pull Request and ran the review by hand. Spec 0211 had ended the
same way on 2026-10-02, before that exception existed.

Two rules settle it:

- A Delivery Retry of an archived item whose archive records
  `qa_override: true` and that has no candidate head accepts the item head
  when Git proves it descends from the Implement start head of the item's
  recorded Run, whatever the park. It records that head as the candidate and
  resumes at `reviewing`. An archive without the override still needs a
  recorded candidate, and `corrective-spec-required` and
  `pull-request-conflict` keep their own rules.
- An unresolved Run whose newest QA Report is `partial`, has no
  finding-blocked row and has at least one environment-blocked row parks
  `qa-environment-partial`, including when every environment-blocked row is
  one that waits for an open Pull Request. Such a partial has no finding to
  fix: a `run-unresolved` retry would run the gate again on the same
  environment and park again.

ADR-0167 is unchanged: the pre-PR Pull Request row still never decides a
qualifying declared partial, and the Archive Command still refuses a partial
that is not qualifying unless the operator records the override. ADR-0223's
descendant rule is unchanged for an item that has a candidate.

The alternative of letting the queue archive a partial blocked only by Pull
Request rows was rejected: the gate wrote `partial` because it judged that
those rows decide it, and only the operator may override that verdict.

## Consequences

The operator's recovery for an environment-only partial is the same whatever
rows block it, and it no longer ends in a hand-made Pull Request. The Run
start head is a weak anchor, since every commit of the item branch descends
from it, so the override record is what authorizes the retry; review,
repository gate, delivery authorization and required checks still run on the
recorded head.

**Qualifying partials do not park (2026-10-06).** Superseded in part by ADR-0240: a partial whose only environment-blocked rows are pre-PR Pull Request rows or network-denied outside-evidence rows qualifies, settles `completed`, and never parks `qa-environment-partial`; only an environment-blocked row outside those two kinds parks it. The operator archive's retry rule stands.
