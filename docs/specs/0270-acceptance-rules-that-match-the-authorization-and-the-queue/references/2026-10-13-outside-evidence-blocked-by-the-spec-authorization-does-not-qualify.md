---
type: fix
status: promoted
created: 2026-10-13
spec: 0270-acceptance-rules-that-match-the-authorization-and-the-queue
reason: null
---

# Outside evidence blocked by the Spec's own authorization does not qualify

## Symptom

The QA of 0258 found no defect, but recorded four outside-evidence rows
(original publications) as `blocked (environment: network prohibited by Spec
authorization)`. The qualifying-partial policy only accepts `blocked
(environment: network denied: <host>)`, so the item parked as
`qa-environment-partial` and needed a QA override to archive.

## Where

The QA skill's wording for an outside-evidence row, and the declared-acceptance
eligibility policy shared by settlement and archive.

## Expected

When the Spec's authorization forbids the network, the QA records the row in
the qualifying form, or the policy accepts the authorization-prohibited form
for an outside-evidence row, so a partial with no finding archives without an
override.

## Evidence

Delivery Queue item 0258 on 2026-10-13 (operator log entry 311).
