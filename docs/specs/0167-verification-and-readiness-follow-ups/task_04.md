---
task: task_04
spec: 0167-verification-and-readiness-follow-ups
status: pending
type: backend
complexity: medium
---

# Task 04: A quiet project init and a budget test without a race

## Overview

Carried from Spec 0165: `roundfix init` writes the default template, which contains `runs.max_active`, to Project Config, so every later command warns that the key is ignored there; and `TestRunImplementBudgetExceededPreservesRunWorktreeAndBranch` needs a Task to settle inside a real 500 ms budget and flakes without load.

## Requirements

1. MUST leave `runs.max_active` out of the template written to Project Config, keeping it in the User Config template.
2. MUST let a project-scope `init` followed by any command print no ignored-key warning.
3. MUST inject the Implement budget clock through the CLI command dependencies, and make the budget test cross the budget only after the fake committer records the settled Task, keeping the test's assertions.

## Subtasks

- [ ] Implement the requirements above.
- [ ] Add a test for each acceptance criterion.

## Acceptance Criteria

- [ ] Project-scope `init` then `Load` warns nothing.
- [ ] The budget test no longer depends on wall-clock timing.

## Context

- instruction: `.agents/skills/implement-task/SKILL.md`
- interface: `internal/config/config.go`
- interface: `internal/cli/implement.go`

## Verification

- `out="$(go test -count=1 -v -run "^(TestProjectInitWritesNoRunCeiling|TestProjectInitThenLoadWarnsNothing|TestBudgetClockCrossesAfterTheSettledTask)$" ./internal/config ./internal/cli 2>&1)" || { printf "%s\\n" "$out"; exit 1; }; for name in TestProjectInitWritesNoRunCeiling TestProjectInitThenLoadWarnsNothing TestBudgetClockCrossesAfterTheSettledTask; do printf "%s\\n" "$out" | grep -q -- "--- PASS: $name" || exit 1; done` — expected: exit 0; before this Task none of the named cases exists, so the command fails.

## References

- [_techspec.md](_techspec.md) — Build Order
