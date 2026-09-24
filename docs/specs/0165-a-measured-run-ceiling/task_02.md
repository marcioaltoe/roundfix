---
task: task_02
spec: 0165-a-measured-run-ceiling
status: pending
type: test
complexity: low
---

# Task 02: A budget test that does not race the header

## Overview

`TestRunImplementBudgetExceededPreservesRunWorktreeAndBranch` configures a 500 ms Run Budget and reads the Run id only from the `Implement Run: <id>` header, which a loaded machine prints after the budget has expired, so the test fails under load, on `main` too.

## Requirements

1. MUST make the budget test identify its Run from the Run Database, or from any stderr line naming it, not from the header alone.
2. MUST add `TestBudgetExceededRunIsFoundWithoutTheHeaderLine`, proving the lookup succeeds for a stderr that lacks the header.
3. MUST NOT change the configured budget or any production code.

## Subtasks

- [ ] Implement the requirements above.
- [ ] Add a test for each acceptance criterion.

## Acceptance Criteria

- [ ] The budget test passes whether or not the header line was printed.

## Context

- instruction: `.agents/skills/implement-task/SKILL.md`
- interface: `internal/cli/implement_test.go`

## Verification

- `out="$(go test -count=1 -v -run "^(TestBudgetExceededRunIsFoundWithoutTheHeaderLine|TestRunImplementBudgetExceededPreservesRunWorktreeAndBranch)$" ./internal/cli 2>&1)" || { printf "%s\\n" "$out"; exit 1; }; for name in TestBudgetExceededRunIsFoundWithoutTheHeaderLine TestRunImplementBudgetExceededPreservesRunWorktreeAndBranch; do printf "%s\\n" "$out" | grep -q -- "--- PASS: $name" || exit 1; done` — expected: exit 0; before this Task the first case does not exist, so the command fails.

## References

- [_techspec.md](_techspec.md) — The budget test
