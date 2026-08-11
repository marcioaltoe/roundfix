---
status: accepted
created_at: 2026-07-24T21:27:41Z
updated_at: 2026-08-11T21:35:56Z
deprecated_at: null
superseded_by: null
---

# ADR-0063: Repositories own their HTTP contract

The Standard TypeScript Monorepo Profile requires each repository to persist an HTTP Contract Decision instead of imposing universal REST or POST-only semantics. Setup reuses an existing supported contract or asks for the mode and its typed exceptions, because framework capability cannot safely infer application API policy.
