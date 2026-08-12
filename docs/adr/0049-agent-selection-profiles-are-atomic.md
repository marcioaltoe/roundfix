---
status: accepted
created_at: 2026-07-17T00:37:58Z
updated_at: 2026-08-11T21:35:56Z
deprecated_at: null
superseded_by: null
---

# Agent Selection Profiles are atomic

Roundfix resolves Agent work through explicit Agent Selection Profiles for `general`, every required Task Type profile, `qa`, and `review`; optional Task Type profiles inherit `general` only when absent. User Config overrides built-ins and Project Config overrides User Config, but each present profile replaces the complete lower-precedence profile and MUST contain one Preferred Selection plus a non-empty Fallback Chain, so no fallback is inherited invisibly from another configuration layer. This refines ADR-0037 and ADR-0040: exact runtime, official model identifier, and reasoning effort remain Roundfix-owned, and every configured tuple must be proven through the installed ACP adapter rather than accepted from a static compatibility assumption.
