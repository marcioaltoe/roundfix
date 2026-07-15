---
task: task_05
spec: 0029-launch-and-recovery-fixes
status: pending
type: backend
complexity: medium
---

# Task 05: Prefer failed-Task surfaces in settle and report the chosen surface

## Overview

Fix the settle surface bug reproduced during the 0027 recovery: resolution picked the latest kept Run's stale worktree — where the target Task had never run and was still `pending` — and refused, while the authoritative checkout had the Task `failed`. Surface resolution now selects the first candidate whose task file has status `failed` (Task Worktree → Run Worktree → current repository), always names the chosen surface, and enumerates every candidate's status when none qualifies.

## Requirements

1. MUST load the target Task's status per candidate surface and select the first, in the existing order (Task Worktree, then the kept Run's Run Worktree, then the current repository), whose task file has status `failed`.
2. MUST print one stderr line `Settle surface: <path>` before running Verification, on every settle that proceeds.
3. MUST, when no candidate surface has the Task `failed`, refuse through the existing preflight-failure convention naming each candidate path with the status found there (or that the path does not exist), keeping the existing guidance about which command owns non-failed statuses.
4. MUST leave the settle stdout contract (`verify …`, `commit <path>`, `settled …`) and all commit/integration mechanics unchanged.

## Subtasks

- [ ] Per-candidate status loading in surface resolution with failed-first selection
- [ ] `Settle surface:` stderr line on the proceed path
- [ ] Per-candidate refusal message when no surface qualifies
- [ ] Regression test of the field case: two seeded kept Runs where the stale worktree has the Task `pending` and the checkout has it `failed` — settle selects the checkout and settles
- [ ] Tests for the refusal enumeration and for the surface line

## Acceptance Criteria

- [ ] With a stale kept worktree (Task `pending`) and the checkout Task `failed`, settle proceeds from the checkout, printing its path in the `Settle surface:` line
- [ ] With the Task `failed` only in a kept Run Worktree, settle still selects that worktree as today
- [ ] With no surface holding the Task `failed`, the refusal names every candidate path and its status
- [ ] The full test suite passes

## Context

- interface: `internal/cli/settle.go`
- interface: `internal/spec/task.go`

## Verification

- `grep -q "Settle surface:" internal/cli/settle.go` — expected: exit 0 (surface reporting exists)
- `go test ./internal/cli/ -run 'Settle'` — expected: settle tests pass, including the stale-surface regression
- `go test ./internal/cli/...` — expected: all tests pass
- `go build -buildvcs=false ./...` — expected: clean build

## References

`_prd.md` → Goal 4, Core Feature 4, Problem 3; `_techspec.md` → Build Order 5, Interfaces (resolveSettleSurface), API Contracts (settle refusal and surface line), Risks (settle ordering visibility).
