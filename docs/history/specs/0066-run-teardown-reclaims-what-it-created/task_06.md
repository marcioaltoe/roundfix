---
task: task_06
spec: 0066-run-teardown-reclaims-what-it-created
status: completed
type: qa
complexity: high
---

# Task 06: Run the final QA gate

## Overview

The authored terminal gate. The Daemon executes `qa-gate` as this node once
task_05 settles `completed`, producing the dated report, verdict, and typed
blocked-row counts per ADR-0080.

This Spec reclaims things. The gate's most valuable rows are the ones proving
it does not reclaim what it must not.

## Requirements

1. MUST run only after task_05 settles `completed`.
2. MUST observe, on a real fixture, that stopping a Run leaves no descendant
   running — including a child that outlives its immediate parent.
3. MUST observe that an unprovable termination is reported distinctly from a
   proven one, rather than accepting the Task's claim.
4. MUST observe that four failed-cycle Run Branches let `watch` proceed with no
   ref hand-deleted, and that a second reconcile pass after applying is a no-op.
5. MUST observe that an Active Run is untouched by every reclamation path.
6. MUST classify any finding by user impact and record typed blocked-row counts.

## Acceptance Criteria

- [ ] The gate runs only after task_05 settles `completed`.
- [ ] A dated QA report lands under `qa/` with a machine-readable verdict and
      typed blocked-row counts.
- [ ] The report records the outliving-child termination observed directly.
- [ ] The report records the Active Run guard observed independently on every
      reclamation path.

## Verification

- `ls docs/specs/0066-run-teardown-reclaims-what-it-created/qa/ | grep -q "qa-report-"`
  — expected: exit 0; a dated QA report exists.
- `grep -l "verdict:" docs/specs/0066-run-teardown-reclaims-what-it-created/qa/qa-report-*.md | grep -q qa-report`
  — expected: exit 0; the report carries a machine-readable verdict.

## References

- `_prd.md` → Goals; Success Metrics.
- `_techspec.md` → Testing Approach.
- ADR-0044, ADR-0052, ADR-0080, ADR-0091.
