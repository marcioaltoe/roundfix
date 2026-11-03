---
type: feat
status: open
created: 2026-11-03
spec: null
reason: null
---

# The Backlog and Findings contracts name the Waiting Declaration

## Opportunity

Spec 0343 lets an open Backlog Entry or an unresolved Finding declare
`waiting_on` and `since`, which `roundfix follow-up` reads to hold the item
back as waiting. The Baseline's copyable Backlog Operational Contract and
Findings Operational Contract in `docs/agents/docs-layout.md` do not name the
two fields, so an agent in an adopter that writes a Finding outside the
Roundfix skill's follow-up mode never learns that a wait can be declared.

## Value

Every adopter's guide would show the fields where it already shows `status`
and `created`, and the declaration would stop depending on the follow-up mode
being the writer. ADR-0353 deferred this to keep Spec 0343 off the Baseline.

## Shape

Add the two optional keys, commented, to both copyable frontmatter blocks of
the `context-workflow` module, with one sentence on when to write them; the
module version, the derived Baseline files and the Source Baseline retention
of the two clauses follow the Baseline module rules. Non-binding.
