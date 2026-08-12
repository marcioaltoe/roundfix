---
status: accepted # proposed | accepted | rejected | deprecated | superseded
created_at: 2026-08-12T12:59:56Z
updated_at: 2026-08-12T12:59:56Z
deprecated_at: null # null or YYYY-MM-DDTHH:MM:SSZ
superseded_by: null # null or ADR-NNNN
---

# A proposed decision is pending, not retired

The documentation layout treats `proposed`, `rejected`, `deprecated`, and
`superseded` ADRs alike as inactive, which is correct for the question it answers
— whether a decision is in force — and wrong as a retirement rule, because a
proposed ADR is a live proposal awaiting a decision rather than history. Only
`rejected`, `deprecated`, and `superseded` retire an ADR to the archive root; a
`proposed` ADR stays in the active directory, and so does a legacy ADR whose body
does not mark it inactive. Reading one status field for two different questions
is what would otherwise file a pending proposal as history, where nobody looking
for an open decision would find it.
