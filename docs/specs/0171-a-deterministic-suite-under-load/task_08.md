---
task: task_08
spec: 0171-a-deterministic-suite-under-load
status: pending
type: backend
complexity: low
---

# Task 08: The bootstrap timeout test observes the start it classifies

## Overview

Corrective Task from the third QA gate of 2026-09-25, whose repository Verification precondition failed under load on `TestRunBootstrapReturnsBootstrapErrorOnTimeout` in `internal/worktree/worktree_test.go`: `expected post-start classification, got "before-start"`. The test gives `runBootstrap` a 10 ms timeout for `sleep 1`; on a loaded machine the shell has not started within 10 ms, so the timeout fires before start and the classification is legitimately `before-start`. It passes 5/5 in isolation.

## Requirements

1. MUST make the test's post-start timeout case deterministic: the bootstrap command must be observed to have started (for example it creates a marker file the test waits for through `internal/testwait`, or a timeout long enough relative to an observed start) before the timeout can fire, while the command still outlives the timeout; the test keeps asserting `BootstrapFailureAfterStart` and the exact error text for the timeout it sets.
2. MUST keep a separate named test for the before-start classification if one exists, and MUST NOT change `runBootstrap` or any production classification.
3. MUST keep the test's wall time under 5 s.

## Subtasks

- [ ] Implement the requirements above.
- [ ] Prove the fix under the suite's load shape.

## Acceptance Criteria

- [ ] The post-start timeout test passes 20 of 20 iterations at `-cpu 1,4` while another package's tests run beside it.

## Context

- instruction: `.agents/skills/implement-task/SKILL.md`
- interface: `internal/worktree/worktree_test.go`

## Verification

- `out="$(go test -count=20 -cpu 1,4 -v -run "^TestRunBootstrapReturnsBootstrapErrorOnTimeout$" ./internal/worktree 2>&1)" || { printf "%s\\n" "$out"; exit 1; }; printf "%s\\n" "$out" | grep -q -- "--- PASS: TestRunBootstrapReturnsBootstrapErrorOnTimeout" && ! grep -q "Timeout: 10 \* time.Millisecond" internal/worktree/worktree_test.go` — expected: exit 0; before this Task the test still sets a 10 ms timeout, so the command fails.

## References

- [_techspec.md](_techspec.md) — Build Order
