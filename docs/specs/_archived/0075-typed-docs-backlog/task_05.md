---
task: task_05
spec: 0075-typed-docs-backlog
status: completed
type: qa
complexity: medium
---

# Task 05: Run the final QA gate

## Overview

The authored terminal gate. The Daemon executes `qa-gate` as this node once both
non-QA leaves settle, producing the dated report, verdict, and typed
blocked-row counts per ADR-0080.

This Spec ships no binary behaviour, so the gate's evidence is documentary: what
a fresh apply produces, and whether the contract a reader receives is usable.

## Requirements

1. MUST run only after task_03 and task_04 settle `completed`.
2. MUST verify a fresh baseline apply produces a guide documenting
   `docs/backlog/` and its contract, rather than reading the checked-in copy.
3. MUST write one backlog entry of each of the four types from the template
   alone, to prove the templates are usable without reading the Spec.
4. MUST verify no finding was migrated and the findings contract is unchanged.
5. MUST verify no binary behaviour changed: no command reads or validates the
   backlog.
6. MUST classify any finding by user impact and record typed blocked-row counts.

## Acceptance Criteria

- [ ] The gate runs only after both leaves settle `completed`.
- [ ] A dated QA report lands under `qa/` with a machine-readable verdict and
      typed blocked-row counts.
- [ ] The report records a fresh apply producing the contract.
- [ ] The report records one entry written from each template.

## Verification

- `ls docs/specs/0075-typed-docs-backlog/qa/ | grep -q "qa-report-"` — expected:
  exit 0; a dated QA report exists.
- `grep -l "verdict:" docs/specs/0075-typed-docs-backlog/qa/qa-report-*.md | grep -q qa-report`
  — expected: exit 0; the report carries a machine-readable verdict.

## References

- `_prd.md` → Goals; Success Metrics.
- `_techspec.md` → Testing Approach.
- ADR-0080, ADR-0091, ADR-0092.
