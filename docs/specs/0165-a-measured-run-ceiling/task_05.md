---
task: task_05
spec: 0165-a-measured-run-ceiling
status: pending
type: test
complexity: low
---

# Task 05: Keep what the budget test proved

## Overview

Corrective Task from the pre-PR review of 2026-09-24. Task 02 went beyond finding the Run without the header: it pre-seeded `task_01` as completed, expected no commit or verification, and let `task_02` be pending or in progress, so the test no longer proves that work settled during the Run survives BudgetExceeded.

## Requirements

1. MUST restore the original fixture and assertions of `TestRunImplementBudgetExceededPreservesRunWorktreeAndBranch`, keeping only the header-independent Run lookup.
2. MUST add `TestBudgetExceededKeepsWorkSettledBeforeTheBudget` proving a Task settled during the Run keeps its commit and status when the budget expires, without depending on a wall-clock race.
3. MUST NOT change production code or the configured budget.

## Subtasks

- [ ] Implement the requirements above.
- [ ] Add a test for each acceptance criterion.

## Acceptance Criteria

- [ ] The budget test again asserts the settled Task's commit and status.

## Context

- instruction: `.agents/skills/implement-task/SKILL.md`
- interface: `internal/cli/implement_test.go`

## Verification

- `out="$(go test -count=1 -v -run "^(TestRunImplementBudgetExceededPreservesRunWorktreeAndBranch|TestBudgetExceededRunIsFoundWithoutTheHeaderLine|TestBudgetExceededKeepsWorkSettledBeforeTheBudget)$" ./internal/cli 2>&1)" || { printf "%s\\n" "$out"; exit 1; }; for name in TestRunImplementBudgetExceededPreservesRunWorktreeAndBranch TestBudgetExceededRunIsFoundWithoutTheHeaderLine TestBudgetExceededKeepsWorkSettledBeforeTheBudget; do printf "%s\\n" "$out" | grep -q -- "--- PASS: $name" || exit 1; done` — expected: exit 0; before this Task the last case does not exist, so the command fails.

## References

- [_techspec.md](_techspec.md) — The budget test
