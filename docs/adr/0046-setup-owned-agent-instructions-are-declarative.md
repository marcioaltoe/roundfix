---
status: accepted
created_at: 2026-07-17T00:37:58Z
updated_at: 2026-08-11T21:35:56Z
deprecated_at: null
superseded_by: null
---

# ADR-0046: Setup-owned agent instructions are declarative

The setup workflow records selected profiles, modules, decisions, and managed artifacts in a versioned manifest and wraps generated Markdown in stable ownership markers. Updates resolve these identifiers instead of inferring intent from prose, preserve repository-authored content outside managed boundaries, and request confirmation only when a decision identifier is missing or incompatible. This makes audit and safe correction deterministic, portable, and idempotent across template revisions.
