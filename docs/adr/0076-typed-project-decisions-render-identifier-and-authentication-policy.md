---
status: accepted
created_at: 2026-07-24T20:34:11Z
updated_at: 2026-07-24T20:34:11Z
deprecated_at: null
superseded_by: null
---

# ADR-0076: Typed project decisions render identifier and authentication policy

The Baseline catalog represents identifier strategy as a discriminated object
and suggests `{"kind":"uuid-v7"}` for new project-owned Internal Identifiers.
It represents a selected Better Auth capability through an `auth.provider`
decision whose typed route exception is deterministically merged into the
repository-owned `http.contract`.

Human callers must confirm or change each unresolved suggestion. Automation
must supply the same structured values explicitly; repository evidence can
suggest applicability but cannot authorize either policy.

The duplicated persisted projection in `auth.provider` and `http.contract`
costs validation complexity. Deterministic conflict detection is required so
the two values can never render contradictory backend guidance.
