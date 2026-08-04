---
task: task_06
spec: 0070-declared-unreachable-acceptance
status: pending
type: qa
complexity: high
---

# Task 06: Run the final QA gate

## Overview

The authored terminal gate. The Daemon executes `qa-gate` as this node once
task_05 settles `completed`, producing the dated report, verdict, and typed
blocked-row counts per ADR-0080.

This gate is unusual: it runs under the contract it just changed. Core Feature
7 — a static-gate failure blocking only the rows it implicates — has no code to
test it, because its correctness lives in this skill's own prompt. This gate is
where that behaviour is observed, and observing it is part of the acceptance.

## Requirements

1. MUST run only after task_05 settles `completed`.
2. MUST exercise the archive boundary as an operator does — a declared case, a
   finding-blocked case, an environment-blocked case, and an unmatched
   declaration — rather than by reading tests.
3. MUST record the three typed blocked-row counts separately, and MUST NOT fold
   one cause into another.
4. MUST observe and record its own scoping behaviour: if any check fails, the
   report states which rows that failure implicates and shows the rows it did
   not block still reporting their own results.
5. MUST verify that `qa_override` still archives a genuinely failed Spec.
6. MUST classify any finding by user impact.

## Acceptance Criteria

- [ ] The gate runs only after task_05 settles `completed`.
- [ ] A dated QA report lands under `qa/` with a machine-readable verdict and
      all three typed blocked-row counts.
- [ ] The report records the observed archive outcome for the declared,
      finding-blocked, environment-blocked, and unmatched cases.
- [ ] The report states, in its own words, which rows any failed check
      implicates — the evidence for Core Feature 7.
- [ ] The report records that `qa_override` remains able to archive a failed
      Spec.

## Verification

- `ls docs/specs/0070-declared-unreachable-acceptance/qa/ | grep -q "qa-report-"`
  — expected: exit 0; a dated QA report exists.
- `grep -l "verdict:" docs/specs/0070-declared-unreachable-acceptance/qa/qa-report-*.md | grep -q qa-report`
  — expected: exit 0; the report carries a machine-readable verdict.
- `grep -q "rows_blocked_declared" docs/specs/0070-declared-unreachable-acceptance/qa/qa-report-*.md`
  — expected: exit 0; the third typed count is present in the report.

## References

- `_prd.md` → Goals; Success Metrics; Core Feature 7.
- `_techspec.md` → Testing Approach; Risks & Considerations.
- ADR-0080, ADR-0091.
