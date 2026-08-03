---
task: task_07
spec: 0071-verification-cost
status: completed
type: docs
complexity: low
---

# Task 07: Publish the before-and-after

## Overview

The Spec exists to make verification cheaper, and a claim of "cheaper" is worth
nothing without the two numbers side by side. This Task publishes the
comparison using the baseline's own commands, on the same machine, so the delta
is measured rather than asserted.

## Requirements

1. MUST measure the same figures the baseline records — full suite warm, full
   suite cold, the gate, and per-package times — using the baseline's own
   commands.
2. MUST publish both sets side by side with the delta stated per row.
3. MUST report the parallelism census after, against the counts the baseline
   recorded.
4. MUST state plainly any figure that did not improve, rather than reporting
   only the ones that did.
5. MUST NOT re-derive or edit the recorded baseline; it is the frozen "before".
6. MUST NOT change any code, test, or tooling file.

## Subtasks

- [ ] Re-measure with the baseline's commands.
- [ ] Publish both tables with the delta per row.
- [ ] Report the parallelism census after.
- [ ] State any figure that did not improve.

## Acceptance Criteria

- [ ] The comparison carries the full-suite warm, full-suite cold, gate, and
      per-package figures for both before and after.
- [ ] Every row states its delta.
- [ ] The parallelism census after is reported against the recorded before.
- [ ] Any figure that did not improve is stated as such.
- [ ] The recorded baseline file is byte-identical.
- [ ] `git status --porcelain` shows no path outside this Spec's folder and
      this task file.

## Context

- interface: `docs/specs/0071-verification-cost/baseline/2026-08-02-before.md`

## Verification

- `test -f docs/specs/0071-verification-cost/baseline/2026-08-02-before.md` —
  expected: exit 0; the frozen baseline is present.
- `git diff --name-only HEAD -- docs/specs/0071-verification-cost/baseline | grep -q . && exit 1 || exit 0`
  — expected: exit 0; the baseline was not edited.
- `grep -rqi 'delta' docs/specs/0071-verification-cost` — expected: exit 0; the
  comparison states deltas.
- `git diff --name-only HEAD -- internal/ Makefile .github/ | grep -q . && exit 1 || exit 0`
  — expected: exit 0; this task changed no code or tooling.

## References

- `_prd.md` → Core Features 7; Success Metrics (a published before-and-after
  accompanies the Spec at close).
- `_techspec.md` → Build Order 7; Decisions (the baseline is frozen).
