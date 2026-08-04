# QA-13 — Typed blocked counts

Status: pass

Three independent public fixtures kept the causes separate:

- finding: `rows_blocked_finding is 1; expected 0`;
- environment: `rows_blocked_environment is 1; expected 0`;
- declared shortfall: `rows_blocked_declared is 2`, one declaration, shortfall
  one.

The declared-only accepted fixture used environment 0, finding 0, declared 1.
This closed report has no blocked row of any cause, so all three exact counts
are zero; its one documentation contradiction is a failed row, not a blocked
row, and is not folded into `rows_blocked_finding`.
