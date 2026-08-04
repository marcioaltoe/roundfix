# QA-07 — Unmatched declaration refusal

Status: pass

The QA report claimed two declared-blocked rows while the PRD declared one.
The built Archive Command exited 2 and named `rows_blocked_declared is 2`, one
Spec declaration, and shortfall one. A fresh read retained the active Spec and
found no archived destination.
