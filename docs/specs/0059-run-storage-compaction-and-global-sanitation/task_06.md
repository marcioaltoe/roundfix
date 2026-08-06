---
task: task_06
spec: 0059-run-storage-compaction-and-global-sanitation
status: pending
type: qa
complexity: medium
---

# Task 06: Run the final QA gate

## Overview

The authored terminal gate. The Daemon executes `qa-gate` as this node once
task_05 settles `completed`.

This Spec deletes bytes irreversibly, so this gate's most valuable rows are the
ones proving it refuses when it should.

## Requirements

1. MUST run only after task_05 settles `completed`.
2. MUST observe compaction refusing with an Active Run present, and confirm the
   database file size is unchanged after the refusal.
3. MUST observe that the preview's reclaimable bytes match what compaction
   actually reclaimed, within the declared tolerance, rather than accepting the
   Task's claim.
4. MUST observe sanitation preserving Review Artifacts and at least one
   ambiguous root, on the filesystem rather than in a log line.
5. MUST observe that a second sanitation run reclaims zero.
6. MUST confirm the storage report mutates nothing, compared before and after.
7. MUST confirm no acceptance in this Spec asserts a recorded size, since that
   is the failure mode the graph was authored to avoid.
8. MUST classify any finding by user impact and record typed blocked-row counts.

## Acceptance Criteria

- [ ] The gate runs only after task_05 settles `completed`.
- [ ] A dated QA report lands under `qa/` with a machine-readable verdict and
      typed blocked-row counts.
- [ ] The report records the Active-Run refusal observed end to end, including
      the unchanged file size.
- [ ] The report records preview-equals-result observed independently.

## Verification

- `ls docs/specs/0059-run-storage-compaction-and-global-sanitation/qa/ | grep -q "qa-report-"`
  — expected: exit 0; a dated QA report exists.
- `grep -l "verdict:" docs/specs/0059-run-storage-compaction-and-global-sanitation/qa/qa-report-*.md | grep -q qa-report`
  — expected: exit 0; the report carries a machine-readable verdict.

## References

- `_prd.md` → Goals; Success Metrics.
- `_techspec.md` → Testing Approach; Measurement, not literals.
- ADR-0033, ADR-0052, ADR-0080, ADR-0091.
