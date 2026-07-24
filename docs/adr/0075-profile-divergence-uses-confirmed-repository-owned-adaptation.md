---
status: accepted
created_at: 2026-07-24T20:34:11Z
updated_at: 2026-07-24T20:34:11Z
deprecated_at: null
superseded_by: null
---

# ADR-0075: Profile divergence uses confirmed repository-owned adaptation

A missing required profile-specific capability cannot be waived while keeping
the built-in Profile identity. The human Baseline flow can instead propose a
repository-owned Profile adaptation, validate its modules and capabilities
against the embedded catalog, re-run alignment, and include the Profile file in
the final digest-bound Change Plan.

Automation supplies the same strict custom Profile document explicitly.
Universal required capabilities are not removable through adaptation; their
absence remains blocking and names the supported remediation operation.

The trade-off is a larger planning state and one additional public automation
input. It preserves the meaning of `required` and the no-mutation-before-plan
contract from ADR-0068.
