---
status: accepted
created_at: 2026-08-08T00:00:00Z
updated_at: 2026-08-08T00:00:00Z
deprecated_at: null
superseded_by: null
---

# A Spec accepts on evidence it did not author

Spec 0082 passed its QA with zero blocked rows and shipped a command that could
not run on six of the eight repositories it existed to update, because the
rehearsal case and the requirement had one author and one premise: the case
observed that the command blocks, which is what the requirement asked for. A
gate that measures a design against a rubric written by that design's author
confirms rather than tests. Every Spec therefore rests at least one named
acceptance row on evidence originating outside its own artifacts — a real
repository, a measurement, or published literature — and records where that
evidence came from; a row whose external evidence cannot be obtained is recorded
as blocked with its reason rather than dropped. The Secondbrain's
`verificacao-adversarial-e-oraculos-de-agentes` supplies the underlying rule,
that a gate is trustworthy only with evidence it observed the right property and
can fail a known negative, and reports that 46.0% of comparable positive results
in a replayed corpus carried no bug-discriminating information. Separating the
author of acceptance criteria from the author of implementation was considered
and deferred as the structurally more expensive answer to the same problem; this
decision is the cheap one that would have caught the measured failure.
