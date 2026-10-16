---
type: feat
status: promoted
created: 2026-10-16
spec: 0274-authoring-guards-hold-the-plan-the-glossary-and-the-prd
reason: null
---

# The glossary is not held to the rule it applies to Specs

## Opportunity

ADR-0244 holds every Spec to the glossary: a bolded phrase of two to five
capitalized words that `CONTEXT.md` does not define is a gap, and the QA gate
and `roundfix archive` refuse the Spec. Nothing holds `CONTEXT.md` to itself.
A Repository Contract Test can: every capitalized phrase of two to five words
in a definition or `_Avoid_` line is a defined term, a plural of one, or on a
short allowlist; every ADR link resolves to an accepted ADR; closed lists that
mirror exported constants (Task Type, Agent Work Category, Cause Class,
Selection Encoding, terminal Run states, Mechanical Refusal Codes) equal them;
an `_Avoid_` line is a plain comma-separated list; and no sentence longer than
sixty characters appears verbatim in two entries.

## Value

The congruence audit of 2026-10-14 found, all mechanically detectable: 16
phrases the glossary used without defining, closed lists that had drifted for
months (Task Type without `qa` since ADR-0091, the ACP Runtime list without
Cursor since ADR-0217, the Backlog Entry lifecycle with three of nine
statuses), one rule stated verbatim in four entries, and 54 `_Avoid_` items
naming another defined term. The branch fixed them by hand (pull request
#545), so the test starts green and fails only on regression.

## Shape

An ordinary test under `internal/docscontract`, mirroring the Spec-side
detector of `internal/speccheck`, with the allowlist and the constant mapping
in one reviewable place. Non-binding.

Triaged on 2026-10-16 from the Secondbrain Inbox Entry
`inbox/roundfix/2026-10-14-hold-the-glossary-to-the-rule-adr-0244-applies-to-specs.md`.
