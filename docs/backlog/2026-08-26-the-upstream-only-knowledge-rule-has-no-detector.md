---
type: fix # feat | fix | perf | refactor
status: open # open | promoted | declined
created: 2026-08-26
---

# The upstream-only knowledge rule has no detector

`docs/agents/docs-layout.md` makes it mandatory that the glossary and the agent
guides reference accepted ADRs and never `docs/specs/` or `docs/findings/`
content. Nothing checks it: no docscontract test and no spec check detector
reads guide citations, so a violation survives every gate.

Measured consequence, reported from a fleet repository on 2026-08-26: a session
performing the archival sweep found agent-guide clauses citing findings as
evidence and repointed those citations to `docs/history/findings/` instead of
removing them — preserving the violation in a form the rule's letter no longer
even names, since the rule predates the history root. A prose-only rule is a
partial detector read as a gate, one layer up: the checker family this
repository already measured (see the 0116 PRD's own self-observation).

Fix shape: a docscontract check over `docs/agents/*.md` and `CONTEXT.md` that
refuses references into `docs/specs/`, `docs/findings/`, and their history
roots — with the layout guide itself exempt where it names directories to
define them. The rule's letter should also gain the history roots, since
archived evidence is no more durable a foundation for a guide than active
evidence.
