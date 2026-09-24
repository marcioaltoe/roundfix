---
task: task_07
spec: 0155-a-verify-that-runs-what-changed
status: pending
type: test
complexity: medium
---

# Task 07: The partition follows what the recipes run

## Overview

Corrective Task from the pre-PR review of 2026-09-24. The `internal/cli` half of the partition contract sorts tests with its own copy of the pattern, so it cannot fail when the Makefile's `-run` and `-skip` recipes drift apart or one of them is dropped.

## Requirements

1. MUST derive each set's `internal/cli` membership from what the `verify-changed` recipes actually run, reading their arguments from `make -n` or listing with `go test -list` under the same `-run` and `-skip` patterns.
2. MUST fail when the two lists do not partition every top-level test of the package.
3. MUST carry a negative control that changes a recipe pattern and observes the failure.

## Subtasks

- [ ] Implement the requirements above.
- [ ] Add a test for each acceptance criterion.

## Acceptance Criteria

- [ ] The contract passes on the tree and fails when a recipe pattern drifts.

## Context

- instruction: `.agents/skills/implement-task/SKILL.md`
- interface: `internal/verifyselect/verifyselect_test.go`

## Verification

- `out="$(go test -count=1 -v -run "^(TestPartitionFollowsTheMakefileRecipes)$" ./internal/verifyselect 2>&1)" || { printf "%s\\n" "$out"; exit 1; }; for name in TestPartitionFollowsTheMakefileRecipes; do printf "%s\\n" "$out" | grep -q -- "--- PASS: $name" || exit 1; done` — expected: exit 0; before this Task the case does not exist, so the command fails.

## References

- [_techspec.md](_techspec.md) — Build Order
