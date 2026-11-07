---
status: accepted # proposed | accepted | rejected | deprecated | superseded
created_at: 2026-08-26T00:00:00Z
updated_at: 2026-11-06T23:30:00Z
deprecated_at: null # null or YYYY-MM-DDTHH:MM:SSZ
superseded_by: null # null or ADR-NNNN
---

# Release publication is all-or-nothing, with a bounded token fallback

A read-only publication preflight evaluates the launcher and every platform package as one release set after Verification and before cross-compilation, and no package publishes unless every coordinate is eligible for the exact target version — a partial release is irreversible under npm policy. Because npm proves OIDC identity only at publish time and per-package, the workflow may retry a coordinate failing OIDC auth with the existing NPM_TOKEN during a bounded window, recording each fallback use; token and fallback are removed together after one clean all-OIDC release. The fallback defines no independent rule — it exists so the identity gap cannot produce exactly the partial release the invariant forbids.

**Amended by [ADR-0374](0374-a-rerun-completes-a-release-set-from-the-same-content.md) (2026-11-06).** A coordinate already published at the target version no longer stops publication when its published tarball holds the same content as the tag run's build: the preflight reports it `published`, the tag run compares the decompressed tar streams after the build, and the publish step skips it, so a rerun completes a Release Set that stopped part way. A coordinate published with different content, or whose version was unpublished, still stops publication before any package publishes, and the bounded token fallback is unchanged.

Consolidates ADR-0082 and ADR-0084 (2026-08-26); both are retired, and Git keeps their text.
