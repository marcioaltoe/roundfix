---
task: task_02
spec: 0003-dogfood-polish
status: pending
type: backend
complexity: low
---

# Task 02: Make the implement Run header truthful

## Overview

The Implement Command's Run header prints two review-path facts that do not
apply to spec Runs: a Run Budget that is never enforced and a `Round: -`
placeholder. Remove both lines for spec Runs only. Verifiable through the
existing header rendering tests.

## Requirements

1. MUST omit the `Budget:` line from the implement Run header — the Implement
   Command deliberately does not enforce the Run Budget (0001 techspec
   decision).
2. MUST omit the `Round:` line from the implement Run header — Round is a
   review-path concept.
3. MUST keep the review-path headers (fetch, resolve, watch) byte-identical.
4. MUST update only the header assertions that cover the two removed lines.

## Subtasks

- [ ] Spec-Run header rendering without Budget and Round lines
- [ ] Review-path header regression asserts untouched
- [ ] Deliberate test updates for the removed lines

## Acceptance Criteria

- [ ] An implement Run header renders Target/Run blocks with no `Budget:` and
      no `Round:` line; all remaining lines are unchanged.
- [ ] Review-path header tests pass without edits.
- [ ] The full suite passes.

## Verification

- `rtk go test ./internal/cli/` — expected: all tests pass.
- `rtk go test ./...` — expected: full suite passes.

## References

`_prd.md` → User Story 2; Core Feature 2. `_techspec.md` → API Contracts,
Build Order 2. Dogfood findings 2 and 3.
