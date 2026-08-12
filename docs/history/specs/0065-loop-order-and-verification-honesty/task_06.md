---
task: task_06
spec: 0065-loop-order-and-verification-honesty
status: completed
type: qa
complexity: medium
---

# Task 06: Run the final QA gate

## Overview

The authored terminal gate. The Daemon executes `qa-gate` as this node once
task_05 settles `completed`.

This Spec exists because a Task settled `completed` having done nothing, so
this gate's most valuable row is the one proving that Task can no longer be
authored.

## Requirements

1. MUST run only after task_05 settles `completed`.
2. MUST observe that Spec 0060's `task_03`, replayed as written, is refused for
   both its contradictory requirements and its work-independent Verification,
   rather than accepting the Tasks' claims.
3. MUST observe that a legitimate Verification containing a repository-wide
   gate still passes, since over-reach is this Spec's stated risk.
4. MUST observe the loop order stated identically in all three sources, and a
   seeded divergence in each one detected.
5. MUST confirm every active and archived Spec still checks as it does today,
   and that the corpus budget test passes.
6. MUST confirm no archived Spec artifact was rewritten and no recorded
   evidence retroactively invalidated, which the PRD declares out of scope.
7. MUST classify any finding by user impact and record typed blocked-row
   counts.

## Acceptance Criteria

- [ ] The gate runs only after task_05 settles `completed`.
- [ ] A dated QA report lands under `qa/` with a machine-readable verdict and
      typed blocked-row counts.
- [ ] The report records the Spec 0060 replay observed end to end.
- [ ] The report records the false-positive check observed independently.

## Verification

- `ls docs/specs/0065-loop-order-and-verification-honesty/qa/ | grep -q "qa-report-"`
  — expected: exit 0; a dated QA report exists.
- `grep -l "verdict:" docs/specs/0065-loop-order-and-verification-honesty/qa/qa-report-*.md | grep -q qa-report`
  — expected: exit 0; the report carries a machine-readable verdict.

## References

- `_prd.md` → Goals; Success Metrics.
- `_techspec.md` → Testing Approach.
- ADR-0080, ADR-0091, ADR-0093.
