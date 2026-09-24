---
task: task_01
spec: 0165-a-measured-run-ceiling
status: pending
type: backend
complexity: medium
---

# Task 01: A machine-wide Active Run ceiling

## Overview

Nothing bounds how many Implement Runs one machine starts; five at once on 2026-09-24 overloaded the machine enough to fail an unrelated QA gate.

## Requirements

1. MUST add `runs.max_active` to User Config: integer, default 3, `0` disables, a negative value refused by validation, documented in `docs/user-guide/configuration.md`.
2. MUST count Active Runs of kind implement across every repository in the Run Database.
3. MUST make `implement` refuse with exit 2 before creating a Run when that count has reached the ceiling, listing each holding Run's id, repository and Spec.
4. MUST leave the per-checkout Active Run check unchanged.

## Subtasks

- [ ] Implement the requirements above.
- [ ] Add a test for each acceptance criterion.

## Acceptance Criteria

- [ ] At `max_active: 2` with Active Runs in two repositories, a third `implement` refuses naming both and creates nothing.
- [ ] `0` disables the bound; a negative value is refused.

## Context

- instruction: `.agents/skills/implement-task/SKILL.md`
- interface: `internal/config/config.go`
- interface: `internal/store/store.go`
- interface: `internal/cli/implement.go`

## Verification

- `out="$(go test -count=1 -v -run "^(TestConfigRunsMaxActiveDefaultsToThree|TestConfigRefusesANegativeRunsMaxActive|TestActiveImplementRunsAreCountedAcrossRepositories|TestImplementRefusesAtTheActiveRunCeiling|TestImplementRunCeilingZeroDisables)$" ./internal/config ./internal/store ./internal/cli 2>&1)" || { printf "%s\\n" "$out"; exit 1; }; for name in TestConfigRunsMaxActiveDefaultsToThree TestConfigRefusesANegativeRunsMaxActive TestActiveImplementRunsAreCountedAcrossRepositories TestImplementRefusesAtTheActiveRunCeiling TestImplementRunCeilingZeroDisables; do printf "%s\\n" "$out" | grep -q -- "--- PASS: $name" || exit 1; done` — expected: exit 0; before this Task none of the named cases exists, so the command fails.

## References

- [_techspec.md](_techspec.md) — The ceiling
