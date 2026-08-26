---
status: accepted # proposed | accepted | rejected | deprecated | superseded
created_at: 2026-08-26T00:00:00Z
updated_at: 2026-08-26T00:00:00Z
deprecated_at: null # null or YYYY-MM-DDTHH:MM:SSZ
superseded_by: null # null or ADR-NNNN
---

# Release publication is all-or-nothing, with a bounded token fallback

A read-only publication preflight evaluates the launcher and every platform package as one release set after Verification and before cross-compilation, and no package publishes unless every coordinate is eligible for the exact target version — a partial release is irreversible under npm policy. Because npm proves OIDC identity only at publish time and per-package, the workflow may retry a coordinate failing OIDC auth with the existing NPM_TOKEN during a bounded window, recording each fallback use; token and fallback are removed together after one clean all-OIDC release. The fallback defines no independent rule — it exists so the identity gap cannot produce exactly the partial release the invariant forbids.

Consolidates ADR-0082 and ADR-0084 (2026-08-26); both are archived under docs/history/adr/.
