---
task: task_06
spec: 0073-skill-versions-decoupled-from-the-binary
status: pending
type: qa
complexity: medium
---

# Task 06: Run the final QA gate

## Overview

The authored terminal gate. The Daemon executes `qa-gate` as this node once
task_05 settles `completed`.

This Spec removes a gate rather than adding one, so this gate's most valuable
rows are the ones proving nothing Roundfix genuinely owns lost its protection.

## Requirements

1. MUST run only after task_05 settles `completed`.
2. MUST observe, end to end, that editing an owned skill leaves `make verify`
   green with no regeneration step — the Spec's reason for existing — rather
   than accepting the Task's claim.
3. MUST observe all three readiness states produced independently, including
   an unreachable source reported distinctly from a missing skill.
4. MUST observe that a third-party skill without a version passes and is never
   compared against an owned minimum.
5. MUST confirm every digest protecting an artifact Roundfix generates is still
   enforced, since this Spec removes one class of digest and must not remove
   the others.
6. MUST confirm a Baseline applied before this Spec still validates and that
   archived Spec artifacts are byte-identical.
7. MUST confirm no acceptance in this Spec asserts a recorded digest or a
   recorded version.
8. MUST classify any finding by user impact and record typed blocked-row counts.

## Acceptance Criteria

- [ ] The gate runs only after task_05 settles `completed`.
- [ ] A dated QA report lands under `qa/` with a machine-readable verdict and
      typed blocked-row counts.
- [ ] The report records the owned-skill edit observed end to end.
- [ ] The report records the third-party boundary observed independently.

## Verification

- `ls docs/specs/0073-skill-versions-decoupled-from-the-binary/qa/ | grep -q "qa-report-"`
  — expected: exit 0; a dated QA report exists.
- `grep -l "verdict:" docs/specs/0073-skill-versions-decoupled-from-the-binary/qa/qa-report-*.md | grep -q qa-report`
  — expected: exit 0; the report carries a machine-readable verdict.

## References

- `_prd.md` → Goals; Success Metrics.
- `_techspec.md` → Testing Approach.
- ADR-0080, ADR-0081, ADR-0085, ADR-0091.
