# QA-13 — Typed blocked counts

Status: pass

Fresh public fixtures kept the three causes separate:

- finding: `rows_blocked_finding is 1; expected 0`;
- environment: `rows_blocked_environment is 1; expected 0`; and
- declared: `rows_blocked_declared is 2`, one declaration, shortfall one.

The accepted declared-only fixture used environment 0, finding 0, and declared
1. This rerun has no blocked row of any cause, so its three exact frontmatter
counts close at zero.
