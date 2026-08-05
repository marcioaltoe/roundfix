---
task: task_05
spec: 0067-derived-artifact-regeneration-boundary
status: completed
type: qa
complexity: medium
---

# Task 05: Run the final QA gate

## Overview

The authored terminal gate. The Daemon executes `qa-gate` as this node once
task_07 settles `completed`, producing the dated report, verdict, and typed
blocked-row counts per ADR-0080.

## Requirements

1. MUST run only after task_07 settles `completed`.
2. MUST reproduce the 2026-08-01 failure directly: edit an owned skill, run the
   sanctioned command once, and observe `make verify` green.
3. MUST verify the exhaustiveness test fails on an artifact added without a
   record, rather than accepting the Task's claim.
4. MUST verify no digest value moved as a result of this Spec.
5. MUST classify any finding by user impact and record typed blocked-row counts.

## Acceptance Criteria

- [ ] The gate runs only after task_07 settles `completed`.
- [ ] A dated QA report lands under `qa/` with a machine-readable verdict and
      typed blocked-row counts.
- [ ] The report records the owned-skill-edit journey observed end to end.
- [ ] The report records that the derived tree's digests are unchanged.

## Verification

- `ls docs/specs/0067-derived-artifact-regeneration-boundary/qa/ | grep -q "qa-report-"`
  — expected: exit 0; a dated QA report exists.
- `grep -l "verdict:" docs/specs/0067-derived-artifact-regeneration-boundary/qa/qa-report-*.md | grep -q qa-report`
  — expected: exit 0; the report carries a machine-readable verdict.

## References

- `_prd.md` → Goals; Success Metrics.
- `_techspec.md` → Testing Approach.
- ADR-0080, ADR-0081, ADR-0085, ADR-0091.
