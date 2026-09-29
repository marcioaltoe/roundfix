---
task: task_07
spec: 0179-review-findings-with-evidence-and-no-unselected-providers
status: pending
type: backend
complexity: low
---

# Task 07: A user-supplied revision is never read as a Git option

## Overview

Corrective Task from the second pre-PR review of 2026-09-29. `roundfix review dispose --fixed-by <commit>` in `internal/cli/review.go` passes the operator's value straight to `git rev-parse`; a value such as `--output=<path>` is parsed as a Git option instead of a revision and can write to an arbitrary file.

## Requirements

1. MUST pass `--end-of-options` before every operator-supplied revision given to Git by `review dispose` (and any other `roundfix review` path that forwards an operator revision), so the value is always read as a revision.
2. MUST refuse a `--fixed-by` value that does not resolve to a commit with the existing refusal and exit `2`, including a value that starts with `-`, and write nothing.
3. MUST keep every existing test and message unchanged, and put the new tests in a file of their own.

## Subtasks

- [ ] Implement the requirements above.
- [ ] Add a named test for each acceptance criterion.

## Acceptance Criteria

- [ ] `--fixed-by --output=<path>` exits `2`, creates no file at `<path>`, and appends nothing to the ledger.
- [ ] A valid descendant commit given to `--fixed-by` is still recorded.

## Context

- interface: `internal/cli/review.go`
- creates: `internal/cli/review_dispose_revision_test.go`

## Verification

- `out="$(go test -count=1 -v -run "^(TestDisposeFixedByOptionLikeValueIsRefusedAndWritesNothing|TestDisposeFixedByValidCommitIsStillRecorded)$" ./internal/cli 2>&1)" || { printf "%s\\n" "$out"; exit 1; }; for name in TestDisposeFixedByOptionLikeValueIsRefusedAndWritesNothing TestDisposeFixedByValidCommitIsStillRecorded; do printf "%s\\n" "$out" | grep -q -- "--- PASS: $name" || exit 1; done` — expected: exit 0; before this Task the two new tests do not exist, so the command fails.

## References

- [_techspec.md](_techspec.md) — Build Order
