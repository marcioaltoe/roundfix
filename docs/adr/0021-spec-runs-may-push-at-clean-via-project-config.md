---
status: accepted
created_at: 2026-07-05T22:17:04Z
updated_at: 2026-08-11T21:35:56Z
deprecated_at: null
superseded_by: null
---

# Spec Runs may push at Clean via Project Config

A repository can opt a spec Run into pushing its branch when — and only when — the Run ends Clean, through a Project Config key that defaults to off; every other outcome never pushes, and opening pull requests stays permanently outside Roundfix's scope. This supersedes the "never push" half of ADR-0013 (the commit-per-Task ownership half stands): the original veto protected early dogfooding, while a per-repository opt-in matches the review path's existing configurable Final Push and keeps the decision where it belongs — in the repo's own config.
