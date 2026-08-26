---
status: accepted # proposed | accepted | rejected | deprecated | superseded
created_at: 2026-08-26T00:00:00Z
updated_at: 2026-08-26T00:00:00Z
deprecated_at: null # null or YYYY-MM-DDTHH:MM:SSZ
superseded_by: null # null or ADR-NNNN
---

# Setup is declarative: manifest identity plus declared decision effects

The setup workflow records selected profiles, modules, decisions, and managed artifacts in a versioned manifest and wraps generated Markdown in stable ownership markers; updates resolve those identifiers instead of inferring intent from prose, preserve repository-authored content outside managed boundaries, and ask only when a decision identifier is missing or incompatible. The asset catalog declares each decision's effects — modules activated, artifacts included or excluded, templates selected, dependent decisions, bound values — and preview, audit, and apply all resolve the same declarative Decision Plan. Imperative per-decision branching and profile-per-combination variants were rejected because the three paths would drift or the profile catalog would multiply.

Consolidates ADR-0046 and ADR-0047 (2026-08-26); both are archived under docs/history/adr/.
