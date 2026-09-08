---
status: accepted
created_at: 2026-09-08T17:03:35Z
updated_at: 2026-09-08T17:03:35Z
deprecated_at: null
superseded_by: null
---

# Configured review profiles select the independent reviewer

The maintainer selected Codex as the default independent reviewer and requires an explicit project reviewer to take precedence. Reuse `profiles.review` and its existing built-in, User Config and Project Config provenance, preserving the effective model/effort and declared fallback chain; an automatic launcher must not invent invocation overrides. The legacy `review_source.name` identifies a PR-feedback provider, not the Agent Runtime.
