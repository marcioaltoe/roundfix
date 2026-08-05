---
task: task_05
spec: 0077-a-green-check-is-not-a-review
status: failed
type: qa
complexity: medium
---

# Task 05: Run the final QA gate

## Overview

The authored terminal gate. The Daemon executes `qa-gate` as this node once
task_04 settles `completed`.

This Spec exists because a gate approved unreviewed code, so this gate's most
valuable rows are the ones proving the new gate cannot.

## Requirements

1. MUST run only after task_04 settles `completed`.
2. MUST observe, against the Pull Request #107 payload, that the evidence
   resolves `skipped` and the watch declines to merge — the exact case that
   merged 125 files unreviewed.
3. MUST observe that an unrecognised green check resolves `pending` and does not
   merge, rather than accepting the Task's claim.
4. MUST confirm no recorded payload that verifies today stopped verifying.
5. MUST confirm no automatic retry was introduced, since it is explicitly out of
   scope.
6. MUST classify any finding by user impact and record typed blocked-row counts.

## Acceptance Criteria

- [ ] The gate runs only after task_04 settles `completed`.
- [ ] A dated QA report lands under `qa/` with a machine-readable verdict and
      typed blocked-row counts.
- [ ] The report records the #107 replay observed end to end.
- [ ] The report records the unrecognised-signal case observed independently.

## Verification

- `ls docs/specs/0077-a-green-check-is-not-a-review/qa/ | grep -q "qa-report-"`
  — expected: exit 0; a dated QA report exists.
- `grep -l "verdict:" docs/specs/0077-a-green-check-is-not-a-review/qa/qa-report-*.md | grep -q qa-report`
  — expected: exit 0; the report carries a machine-readable verdict.

## References

- `_prd.md` → Goals; Success Metrics.
- `_techspec.md` → Testing Approach.
- ADR-0054, ADR-0080, ADR-0091.
