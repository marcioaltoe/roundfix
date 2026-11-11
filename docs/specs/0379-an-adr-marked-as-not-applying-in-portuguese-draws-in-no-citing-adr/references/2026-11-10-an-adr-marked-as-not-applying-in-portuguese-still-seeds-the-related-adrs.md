---
type: feat
status: promoted
created: 2026-11-10
spec: 0379-an-adr-marked-as-not-applying-in-portuguese-draws-in-no-citing-adr
reason: null
---

# An ADR marked as not applying in Portuguese still seeds the related ADRs

## Opportunity

An adopter that writes its Specs in Portuguese could mark an ADR as not applying in its own words and have the strict check read it. ADR-0377 reads one fixed leading token, `not applicable`, in the TechSpec's `## ADR Alternatives` line, through `notApplicableRelationPattern` in `internal/speccheck/authoring_proof.go`.

## Value

Fiscus reported on 2026-11-10 that the 76 ADR lines behind the cascade Spec 0367 fixed use "não se aplica". The strict check of its Spec 0081 reports no finding on Roundfix 0.137.6, so nothing fails today. Those lines are still read as applied ADRs, so the amendment of ADR-0377 does not shorten the related-ADR closure for a Portuguese Spec until its author rewrites each line with the English token. The hypothesis is that reading the adopter's own marker saves that rewrite and keeps the closure small for every Portuguese Spec.

## Shape

Non-binding. A small closed list of leading tokens that includes the Portuguese form; or a Project Config key that names the adopter's token; or no code change and a line in the `write-techspec` skill that prescribes the English token whatever the Spec's language. ADR-0377 rejected a marker phrase in the PRD row because "a phrase written in one Spec's language binds no other", so a Spec that takes this up amends that ADR.

Reported by fiscus in an ai-memory message whose id begins `01a1c323`, 2026-11-10, in answer to the operator's request to run the strict check on its Spec 0081. Fiscus did not ask for the change.
