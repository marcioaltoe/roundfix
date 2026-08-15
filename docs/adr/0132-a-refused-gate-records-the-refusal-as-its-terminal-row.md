---
status: accepted # proposed | accepted | rejected | deprecated | superseded
created_at: 2026-08-15T14:30:00Z
updated_at: 2026-08-15T14:30:00Z
deprecated_at: null # null or YYYY-MM-DDTHH:MM:SSZ
superseded_by: null # null or ADR-NNNN
---

# A refused gate records the refusal as its terminal row

A gate that stops at a precondition writes a report whose Results table has no
rows, which its own contract calls malformed, and the next run of the same Spec
then refuses on that report with a fix it cannot perform — it never built a matrix
to materialize rows from. The writer already knows it has zero rows, so it emits
one terminal row naming the refusal that produced it: what stopped the gate, and
that nothing else was executed. That is more useful than an empty table, it is
what actually happened, and it keeps the evidence trail complete where suppressing
the report entirely would leave a later reader unable to tell a refused run from a
run that never happened. The alternative was writing no report at all; it is
smaller in the writer and larger everywhere downstream.
