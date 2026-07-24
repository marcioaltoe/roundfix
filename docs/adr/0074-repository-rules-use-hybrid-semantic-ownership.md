---
status: accepted
created_at: 2026-07-24T20:34:11Z
updated_at: 2026-07-24T20:34:11Z
deprecated_at: null
superseded_by: null
---

# ADR-0074: Repository rules use hybrid semantic ownership

Repository rules use a typed project decision when the catalog models their
policy. Other classifiable rules move byte-for-byte into repository-owned
blocks in the active semantic guide; only a rule with no typed or semantic
owner remains in `docs/agents/specific-repository.md`.

The classifier can segment one source entry into byte-exact clauses, but the
deterministic planner accepts a proposal only when the segments cover every
source byte exactly once. Roundfix never rewrites the accepted clause while
moving it.

This decision refines the repository-rule carrier portion of ADR-0067. It does
not weaken ADR-0070: arbitrary nested instruction carriers remain untouched;
only root Readoption evidence and the three recognized repository-rule carriers
participate in automatic redistribution.
