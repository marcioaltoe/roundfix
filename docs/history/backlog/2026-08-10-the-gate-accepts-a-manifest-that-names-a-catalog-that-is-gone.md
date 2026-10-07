---
type: fix # feat | fix | perf | refactor
status: deferred
created: 2026-08-10
spec: null # Spec slug when status: promoted
reason: null # required when status: declined
---

# The gate accepts a manifest that names a catalog that is gone

`docs/agents/setup-context.json` records the `catalogDigest` that produced the current managed bytes. Nothing verifies that the recorded digest matches the catalog the repository actually ships, so the two drift silently.

Full text in Git at `a30dc847a037fe584812b4471aa7da1701a9ec59`: `docs/history/backlog/2026-08-10-the-gate-accepts-a-manifest-that-names-a-catalog-that-is-gone.md`.
