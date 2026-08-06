---
status: done
created_at: 2026-08-03
updated_at: 2026-08-06
absorbed_by: 2026-08-06-rollup-qa-gates-and-verification-evidence.md
---

# 2026-08-03 — A recording-order defect blocked every public QA row

status: pending

## What was observed

Across Spec 0072's four gate executions, the governance preflight (F-001)
classified fifteen of twenty-four matrix rows `blocked (finding: F-001)` —
every public user-journey row, including the ones the Spec existed to prove.
The finding itself narrowed with each remediation until only one clause
remained, and that clause is unsatisfiable:

> "Spec-specific exact authority for Task 04 landed after the Task; the
> later `roundfix` skill authorization was folded into the same commit as
> the protected edit."

The maintainer's *direction* preceded every protected edit — the gate's own
report acknowledges the quoted direction — but the *per-Spec listing in the
grant file* was written at close, after the first gate report named its
absence. Git history is immutable, so no amount of further work can make
the recording precede the edit. From gate #2 onward the Spec could never
reach a `pass`, only a better-documented `fail`.

Meanwhile the substance was green the whole time: `make verify` at 3,103
tests, the PR approved with all review threads resolved, and the authored
gate mechanism itself exercised for real by the very gate runs that were
failing it.

## Root cause

Two compounding rules in the QA gate's governance audit:

1. Any Tooling Authority defect blocks *all* public-flow rows, so a
   recording-order gap converts into a full QA block on rows that have
   nothing to do with tooling.
2. The audit requires the authorization *record* to chronologically precede
   the protected edit, with no recognized remediation path once history
   exists. An honest late record — stating the chronology explicitly, per
   this repository's own "recorded rather than assumed" precedent — still
   fails.

## What would settle it

- Scope the governance block: a Tooling Authority finding blocks the rows
  that depend on tooling trust, not every public row.
- Define the remediation contract for late recording: a grant recorded
  after the fact, quoting the prior direction and bounding the exact paths,
  is a *remediated* finding (non-blocking, noted), not a permanent fail —
  because the alternative spends `qa_override` on Specs whose only defect
  is paperwork order, diluting the mechanism reserved for failed evidence
  (the exact concern Spec 0070 records for unreachable acceptance).

## Consequence for Spec 0072

Closed with `qa_override` under the maintainer's standing direction. The
override rationale: the sole remaining finding is immutable history,
honestly declared in the active artifacts; every functional row that ran
passed; and the user journeys the blocked rows describe were exercised for
real by the four authored-gate executions themselves.

## Spec pointer

None yet. Candidate to fold into Spec 0070's scope (declared-unreachable
acceptance) or a qa-gate skill revision.
