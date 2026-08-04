---
task: task_03
spec: 0076-force-stop-exit-proof
status: pending
type: qa
complexity: medium
---

# Task 03: Run the final QA gate

## Overview

The authored terminal gate. The Daemon executes `qa-gate` as this node once
task_02 settles `completed`, producing the dated report, verdict, and typed
blocked-row counts per ADR-0080.

This Spec's acceptance is unusual in that its subject is a proof. The gate
should judge whether the repaired test can now fail for the right reason, not
merely whether it passes — a passing test is what this Spec exists to distrust.

## Requirements

1. MUST run only after task_02 settles `completed`.
2. MUST exercise the helper directly, as a user of the test suite would: run it
   alone, observe liveness, and confirm no runtime fatal error.
3. MUST confirm the repeated parent/helper run is clean, since that probe is
   what exposed the original defect.
4. MUST verify the premature-exit regression can fail, rather than accepting
   the Task's recorded claim that it did.
5. MUST classify any finding by user impact and record typed blocked-row counts.

## Acceptance Criteria

- [ ] The gate runs only after task_02 settles `completed`.
- [ ] A dated QA report lands under `qa/` with a machine-readable verdict and
      typed blocked-row counts.
- [ ] The report records the helper's observed liveness and the absence of a
      runtime fatal error.
- [ ] The report records an independently observed failure of the
      premature-exit regression.

## Verification

- `ls docs/specs/0076-force-stop-exit-proof/qa/ | grep -q "qa-report-"` —
  expected: exit 0; a dated QA report exists.
- `grep -l "verdict:" docs/specs/0076-force-stop-exit-proof/qa/qa-report-*.md | grep -q qa-report`
  — expected: exit 0; the report carries a machine-readable verdict.

## References

- `_prd.md` → Goals; Success Metrics.
- `_techspec.md` → Testing Approach.
- ADR-0080, ADR-0089, ADR-0091.
