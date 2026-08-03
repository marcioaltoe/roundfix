---
task: task_07
spec: 0072-qa-is-a-task-not-a-flag
status: failed
type: qa
complexity: high
---

# Task 07: Run the final QA gate

## Overview

The authored terminal gate for this Spec — the first one authored under the
contract this Spec built. The Daemon executes `qa-gate` as this node once
every dependency settles `completed`, producing the dated QA report, the
verdict, and the typed blocked-row counts exactly as ADR-0080 defines them.

A first gate attempt (run under the pre-contract flow, before this node
existed) returned `fail` with two blocking findings; both were remediated:
the tooling authorization addendum of 2026-08-03 lists this Spec's bounded
skill files, and the characterization and coverage records were regenerated
deliberately with only digest, identity, and renamed-test fields moving.
The report at `qa/qa-report-2026-08-03.md` is the evidence this re-run
answers.

## Acceptance Criteria

- [ ] The gate runs only after task_04, task_05, and task_06 are settled
      `completed`.
- [ ] A dated QA report lands under `qa/` with a machine-readable verdict
      and typed blocked-row counts.

## Verification

- `ls docs/specs/0072-qa-is-a-task-not-a-flag/qa/ | grep -q "qa-report-"` —
  expected: exit 0; a dated QA report exists.
- `grep -l "verdict:" docs/specs/0072-qa-is-a-task-not-a-flag/qa/qa-report-*.md | grep -q qa-report`
  — expected: exit 0; the report carries a machine-readable verdict.

## References

- `_prd.md` → Core Features 2, 5; Success Metrics 1–3.
- `_techspec.md` → Build Order; ADR-0091.
