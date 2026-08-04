---
task: task_07
spec: 0068-spec-close-audit
status: pending
type: qa
complexity: high
---

# Task 07: Run the final QA gate

## Overview

The authored terminal gate. The Daemon executes `qa-gate` as this node once
both non-QA leaves settle, producing the dated report, verdict, and typed
blocked-row counts per ADR-0080.

## Requirements

1. MUST run only after task_05 and task_06 settle `completed`.
2. MUST exercise the audit as an operator does — against a real Spec in this
   repository and against fixtures for each survivor kind — rather than by
   reading tests.
3. MUST verify the audit reclaims nothing: Git state before and after must be
   identical.
4. MUST verify an Active Run is never classified as residue, independently of
   the Task's own claim.
5. MUST verify the deleted-target Run resolves by content rather than `unknown`.
6. MUST classify any finding by user impact and record typed blocked-row counts.

## Acceptance Criteria

- [ ] The gate runs only after both leaves settle `completed`.
- [ ] A dated QA report lands under `qa/` with a machine-readable verdict and
      typed blocked-row counts.
- [ ] The report records Git state unchanged across an audit run.
- [ ] The report records the Active Run guard observed independently.

## Verification

- `ls docs/specs/0068-spec-close-audit/qa/ | grep -q "qa-report-"` — expected:
  exit 0; a dated QA report exists.
- `grep -l "verdict:" docs/specs/0068-spec-close-audit/qa/qa-report-*.md | grep -q qa-report`
  — expected: exit 0; the report carries a machine-readable verdict.

## References

- `_prd.md` → Goals; Success Metrics.
- `_techspec.md` → Testing Approach.
- ADR-0052, ADR-0080, ADR-0091.
