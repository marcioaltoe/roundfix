---
task: task_04
spec: 0069-review-run-targets-its-pull-request
status: failed
type: qa
complexity: medium
---

# Task 04: Run the final QA gate

## Overview

The authored terminal gate. The Daemon executes `qa-gate` as this node once
task_03 settles `completed`.

This Spec exists because two legitimate security findings failed for an
environmental reason, so this gate's most valuable rows are the ones proving
that an environmental stop no longer reads as a defect.

## Requirements

1. MUST run only after task_03 settles `completed`.
2. MUST observe that a Run started on a branch other than the Pull Request's
   head refuses at Preflight with exit `2`, creating no Run and writing
   nothing.
3. MUST observe that a checkout moved mid-Run reaches the interruption outcome
   with its Review Issues left unsettled, rather than accepting the Task's
   claim.
4. MUST confirm every review artifact commit in the exercised Runs landed on
   the Pull Request's head branch, read from Git rather than from a log.
5. MUST confirm a Run whose checkout matches behaves as before.
6. MUST confirm Roundfix moved no working tree during any observed Run.
7. MUST classify any finding by user impact and record typed blocked-row
   counts.

## Acceptance Criteria

- [ ] The gate runs only after task_03 settles `completed`.
- [ ] A dated QA report lands under `qa/` with a machine-readable verdict and
      typed blocked-row counts.
- [ ] The report records the Preflight refusal observed end to end.
- [ ] The report records the mid-Run interruption observed independently.

## Verification

- `ls docs/specs/0069-review-run-targets-its-pull-request/qa/ | grep -q "qa-report-"`
  — expected: exit 0; a dated QA report exists.
- `grep -l "verdict:" docs/specs/0069-review-run-targets-its-pull-request/qa/qa-report-*.md | grep -q qa-report`
  — expected: exit 0; the report carries a machine-readable verdict.

## References

- `_prd.md` → Goals; Success Metrics.
- `_techspec.md` → Testing Approach.
- ADR-0052, ADR-0080, ADR-0091.
