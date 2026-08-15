---
status: accepted # proposed | accepted | rejected | deprecated | superseded
created_at: 2026-08-15T14:31:00Z
updated_at: 2026-08-15T14:31:00Z
deprecated_at: null # null or YYYY-MM-DDTHH:MM:SSZ
superseded_by: null # null or ADR-NNNN
---

# The gate performs the repairs its Task names, and its vocabulary check is one of them

The authoring contract gives the QA gate the glossary update for a term the Spec
coined, and the gate's own static precondition refuses on that term being
undocumented — so the repair is assigned to the actor forbidden from reaching it,
and two Specs stalled there on consecutive days. The vocabulary precondition
therefore does not refuse a term declared by the Spec under gate; the gate
documents it and proves it documented. More generally, a repair a gate's Task file
names is work the gate performs and then verifies, not a finding it reports: a
gate given two named repairs on 2026-08-15 wrote both as findings and failed,
which leaves contract-assigned work to whoever reads the report. Reporting stays
correct for everything the Task did not assign.
