---
status: accepted
created_at: 2026-07-24T21:27:41Z
updated_at: 2026-08-11T21:35:56Z
deprecated_at: null
superseded_by: null
---

# ADR-0066: Context-Driven Baseline execution belongs to the CLI

The public Roundfix CLI is the sole runtime authority for Baseline Profile
authoring, audit, Decision Plan, Change Plan, and confirmation-gated apply.
The `setup-context-driven` skill becomes an instructional layer over that
public API, and the Python runtime is removed after behavioral parity, so
humans and automations can operate the Context-Driven Baseline without Codex.
This supersedes Spec 0045's choice to keep the canonical implementation inside
the skill.
