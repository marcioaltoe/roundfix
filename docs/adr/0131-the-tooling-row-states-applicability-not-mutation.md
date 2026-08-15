---
status: accepted # proposed | accepted | rejected | deprecated | superseded
created_at: 2026-08-15T08:11:00Z
updated_at: 2026-08-15T08:11:00Z
deprecated_at: null # null or YYYY-MM-DDTHH:MM:SSZ
superseded_by: null # null or ADR-NNNN
---

# The tooling row states applicability, not mutation

The PRD template tells an author to record the Tooling authority row as
`applicable — no protected tooling mutation proposed or authorized`, and the
checker refuses exactly that, because it reads `applicable` as "this Spec mutates
tooling" and then demands bounded files. Both readings are defensible and only
one can be mechanical. The row states whether the constraint governs the Spec,
which it always does — the rule binds every Spec, including the ones that change
no tooling — and bounded files become required only when the row's own reason
declares a proposed or authorized mutation. The checker changes rather than the
template, because the template's reading is the correct one and because changing
the template would itself be a protected tooling mutation needing a grant to fix
a defect in the thing that hands grants out.
