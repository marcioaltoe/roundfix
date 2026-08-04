---
spec: qa-case
date: 2026-08-04
status: closed
verdict: fail
rows_blocked_environment: 0
rows_blocked_finding: 0
rows_blocked_declared: 0
surfaces: [cli]
---

# QA report — Wrongly declared fixture

## Results

| Criterion | Status |
| --- | --- |
| Archive Command help is reachable | pass — `roundfix archive --help` exited 0 |

## Findings

F-001 is a wrongly-declared-row finding: the declared row was reachable and
was run normally, so the declaration is not accepted.
