---
task: task_01
spec: 0158-daemon-verification-and-access-readiness
status: pending
type: backend
complexity: medium
---

# Task 01: Collect every independent Verification failure

## Overview

A Task's Verification attempt returns at its first failed command, so the one repair turn ADR-0038 allows sees a single failure even when later commands were already failing for their own reasons.

## Requirements

1. MUST accept Task frontmatter `verification: independent`, record it on `spec.Task`, and refuse any other value by name.
2. MUST, for a Task so declared, run every Verification command even after one fails and publish each failure as today.
3. MUST hand every failed command, with its diagnostics, to the single repair turn, and run every command again on the retry.
4. MUST keep today's behaviour for a Task without the declaration: commands run in order and the attempt stops at the first failure.
5. MUST keep the retry ceiling and the handling of a temporary failure unchanged.

## Subtasks

- [ ] Implement the requirements above.
- [ ] Add a test for each acceptance criterion.

## Acceptance Criteria

- [ ] A declared Task with two failing commands of three hands both failures to its repair turn, and the third command ran.
- [ ] An undeclared Task stops at its first failure and later commands never run.
- [ ] An unknown `verification` value is refused by name.

## Context

- instruction: `.agents/skills/implement-task/SKILL.md`
- interface: `internal/spec/task.go`
- interface: `internal/daemon/engine.go`
- interface: `internal/daemon/task_engine.go`

## Verification

- `out="$(go test -count=1 -v -run "^(TestTaskParsesIndependentVerification|TestIndependentVerificationHandsEveryFailureToRepair|TestUndeclaredVerificationStopsAtFirstFailure)$" ./internal/spec ./internal/daemon 2>&1)" || { printf "%s\n" "$out"; exit 1; }; for name in TestTaskParsesIndependentVerification TestIndependentVerificationHandsEveryFailureToRepair TestUndeclaredVerificationStopsAtFirstFailure; do printf "%s\n" "$out" | grep -q -- "--- PASS: $name" || exit 1; done` — expected: exit 0; before this Task none of the three cases exists, so the command fails.

## References

- [_techspec.md](_techspec.md) — Independent Verification
