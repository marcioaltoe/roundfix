---
task: task_02
spec: 0155-a-verify-that-runs-what-changed
status: pending
type: test
complexity: medium
---

# Task 02: A partition that covers every test once

## Overview

A selector that drops a test from both sets makes verification faster by making
it blind. This Task adds the contract that makes that impossible to do quietly.

## Requirements

1. MUST prove every package listed by `go list ./...` is selected by exactly one
   set.
2. MUST prove every test function in `internal/cli` is selected by exactly one of
   the two `internal/cli` invocations.
3. MUST include negative controls showing the contract fails when an item is
   dropped from both sets and when it is added to both.
4. MUST read the set definitions from `internal/verifyselect`, the same source
   the Makefile uses.

## Subtasks

- [ ] Add the package partition check.
- [ ] Add the `internal/cli` test partition check.
- [ ] Add the negative controls.

## Acceptance Criteria

- [ ] The contract passes on the tree as delivered.
- [ ] A package removed from both sets fails the contract, naming it.
- [ ] An `internal/cli` test selected by both invocations fails the contract,
      naming it.

## Context

- instruction: `.agents/skills/implement-task/SKILL.md`
- interface: `Makefile`

## Verification

- `out="$(go test -count=1 -v -run "^TestPartitionCoversEveryTestExactlyOnce" ./internal/verifyselect 2>&1)" || { printf "%s\n" "$out"; exit 1; }; printf "%s\n" "$out" | grep -q -- "--- PASS: TestPartitionCoversEveryTestExactlyOnce"` — expected: exit 0; before this Task the case does not exist, so the command fails.

## References

- [_techspec.md](_techspec.md) — The two sets
