---
task: task_04
spec: 0165-a-measured-run-ceiling
status: pending
type: backend
complexity: medium
---

# Task 04: A ceiling that frees itself and cannot be raced

## Overview

Corrective Task from the pre-PR review of 2026-09-24. The ceiling counts Active Runs whose owner process is dead, forever, and its refusal names no way out; Project Config can raise or disable the machine-wide ceiling; and the count is read outside the transaction that creates the Run, so two concurrent `implement`s in different repositories can both pass it.

## Requirements

1. MUST reclaim each holder whose owner process is dead, the way the per-checkout path already does, before comparing the count.
2. MUST make the refusal's next action name `roundfix stop <id>` for each holder and the `runs.max_active` key.
3. MUST ignore, with a warning, `runs.max_active` from Project Config; only User Config sets it; state this in `docs/user-guide/configuration.md`.
4. MUST count Active Implement Runs inside the same write transaction that creates the Run and refuse there with a typed ceiling error, keeping the preflight check for the early message.

## Subtasks

- [ ] Implement the requirements above.
- [ ] Add a test for each acceptance criterion.

## Acceptance Criteria

- [ ] A holder with a dead owner is reclaimed and no longer blocks.
- [ ] The refusal names `roundfix stop` and `runs.max_active`.
- [ ] A Project Config value does not raise the User Config ceiling.
- [ ] Two creations racing at the ceiling cannot both succeed.

## Context

- instruction: `.agents/skills/implement-task/SKILL.md`
- interface: `internal/cli/implement.go`
- interface: `internal/store/store.go`
- interface: `internal/config/config.go`

## Verification

- `out="$(go test -count=1 -v -run "^(TestRunCeilingReclaimsADeadHolder|TestRunCeilingRefusalNamesTheWayOut|TestProjectConfigCannotRaiseTheRunCeiling|TestRunCeilingIsCountedInsideTheCreateTransaction)$" ./internal/config ./internal/store ./internal/cli 2>&1)" || { printf "%s\\n" "$out"; exit 1; }; for name in TestRunCeilingReclaimsADeadHolder TestRunCeilingRefusalNamesTheWayOut TestProjectConfigCannotRaiseTheRunCeiling TestRunCeilingIsCountedInsideTheCreateTransaction; do printf "%s\\n" "$out" | grep -q -- "--- PASS: $name" || exit 1; done` — expected: exit 0; before this Task none of the named cases exists, so the command fails.

## References

- [_techspec.md](_techspec.md) — The ceiling
