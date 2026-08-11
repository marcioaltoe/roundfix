---
status: accepted
created_at: 2026-07-24T21:27:41Z
updated_at: 2026-08-11T21:35:56Z
deprecated_at: null
superseded_by: null
---

# ADR-0065: Release Plan exposes a read-only reset mode

The Release Plan Command accepts `--reset-to v0.0.1` as a mutually exclusive read-only mode that inventories prior tags and GitHub Releases, produces a digest-bound reset plan, and requires approval without exposing any deletion path. A separate post-QA action needs explicit authority to remove that inventory, preserving ADR-0048's analysis-versus-mutation boundary while keeping the exceptional reset auditable through the public CLI.
