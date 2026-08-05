---
task: task_06
spec: 0078-roundfix-asks-for-the-review
status: completed
type: qa
complexity: medium
---

# Task 06: Run the final QA gate

## Overview

The authored terminal gate. The Daemon executes `qa-gate` as this node once
task_05 settles `completed`.

This Spec exists because a correctness fix left the loop unable to advance, so
this gate's most valuable rows are the ones proving the loop advances now and
that it consumes exactly what the configuration says it will.

## Requirements

1. MUST run only after task_05 settles `completed`.
2. MUST observe that a Round which pushes publishes exactly one review request,
   including the Round whose Final Push is followed by the artifact-only docs
   commit.
3. MUST observe that the same head asked twice publishes once, rather than
   accepting the Task's claim.
4. MUST observe all four rows of the Preflight coherence table, including both
   refusals and their exit codes.
5. MUST confirm `fetch` publishes no request under any configuration.
6. MUST confirm no automatic retry, backoff, or capacity wait was introduced,
   since all are explicitly out of scope.
7. MUST confirm the repository's own committed configuration pair is coherent,
   because an incoherent pair strands every later Run.
8. MUST classify any finding by user impact and record typed blocked-row counts.

## Acceptance Criteria

- [ ] The gate runs only after task_05 settles `completed`.
- [ ] A dated QA report lands under `qa/` with a machine-readable verdict and
      typed blocked-row counts.
- [ ] The report records the one-request-per-Round observation end to end.
- [ ] The report records both Preflight refusals observed independently.

## Verification

- `ls docs/specs/_archived/0078-roundfix-asks-for-the-review/qa/ | grep -q "qa-report-"`
  — expected: exit 0; a dated QA report exists.
- `grep -l "verdict:" docs/specs/_archived/0078-roundfix-asks-for-the-review/qa/qa-report-*.md | grep -q qa-report`
  — expected: exit 0; the report carries a machine-readable verdict.

## References

- `_prd.md` → Goals; Success Metrics.
- `_techspec.md` → Testing Approach.
- ADR-0054, ADR-0080, ADR-0091.
