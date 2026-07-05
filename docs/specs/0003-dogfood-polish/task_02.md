---
task: task_02
spec: 0003-dogfood-polish
status: completed
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

- [x] Spec-Run header rendering without Budget and Round lines
- [x] Review-path header regression asserts untouched
- [x] Deliberate test updates for the removed lines

## Acceptance Criteria

- [x] An implement Run header renders Target/Run blocks with no `Budget:` and
      no `Round:` line; all remaining lines are unchanged.
- [x] Review-path header tests pass without edits.
- [x] The full suite passes.

## Verification

- `rtk go test ./internal/cli/` — expected: all tests pass.
- `rtk go test ./...` — expected: full suite passes.

## References

`_prd.md` → User Story 2; Core Feature 2. `_techspec.md` → API Contracts,
Build Order 2. Dogfood findings 2 and 3.

## Result

Implemented spec-Run header cleanup for the Live Run View:

- `RenderLiveRunView` now omits the `Round:` and `Budget:` Run lines only when
  rendering an Implement Command spec Run (`RunKind: implement`).
- The remaining spec Run Target and Run lines still render: Spec, Branch,
  Agent, ID, State, Git, Auto-commit, Auto-push, and Last push.
- Review-path header rendering stays unchanged; the existing review Run header
  assertions for `Round: 2 / 6` and `Budget: 38m / 2h` were not edited.

Acceptance evidence:

- `TestRenderLiveRunViewSpecRunRendersTasksAsWorkItems` now sets non-empty
  Round and Budget values and asserts `Round:` and `Budget:` are absent for
  spec Runs while the remaining Target/Run lines are present.
- `TestRenderLiveRunViewGroupsIssuesAndShowsStatusStrips` still asserts the
  review-path `Round:` and `Budget:` lines.

Verification:

- `rtk go test ./internal/tui/` passed: 36 tests.
- `rtk go test ./internal/cli/` passed: 144 tests.
- `rtk go test ./...` passed: 453 tests in 16 packages.
- `rtk make verify` passed: full Go suite, `roundfix skills check`, and build.

Follow-ups: none.
