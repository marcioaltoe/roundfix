---
status: accepted
created_at: 2026-07-24T21:27:41Z
updated_at: 2026-08-11T21:35:56Z
deprecated_at: null
superseded_by: null
---

# ADR-0067: Custom Baseline Profiles are repository-owned

Roundfix ships built-in Baseline Profiles, while every custom Baseline Profile
lives in and is versioned with the repository that uses it. User-scoped custom
profiles and cross-repository precedence are excluded so audit, automation, and
review resolve the same profile from repository evidence. Custom profiles live
at `.roundfix/baseline/profiles/<id>.json` and may compose only the modules,
decisions, Repository Capabilities, and templates in the CLI's embedded
catalog. Custom executable code, templates, and modules are excluded;
repository-authored rules remain in
`docs/agents/specific-repository.md`. Baseline creates and links that carrier
only when it has non-empty Repository-Specific Normative Rules to preserve.
