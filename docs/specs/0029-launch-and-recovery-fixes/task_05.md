---
task: task_05
spec: 0029-launch-and-recovery-fixes
status: completed
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

- [x] Per-candidate status loading in surface resolution with failed-first selection
- [x] `Settle surface:` stderr line on the proceed path
- [x] Per-candidate refusal message when no surface qualifies
- [x] Regression test of the field case: two seeded kept Runs where the stale worktree has the Task `pending` and the checkout has it `failed` — settle selects the checkout and settles
- [x] Tests for the refusal enumeration and for the surface line

## Acceptance Criteria

- [x] With a stale kept worktree (Task `pending`) and the checkout Task `failed`, settle proceeds from the checkout, printing its path in the `Settle surface:` line
- [x] With the Task `failed` only in a kept Run Worktree, settle still selects that worktree as today
- [x] With no surface holding the Task `failed`, the refusal names every candidate path and its status
- [x] The full test suite passes

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

## Result

- Implemented failed-first settle surface resolution across Task Worktree, Run Worktree, and current repository candidates; each candidate loads the target Task status and missing paths are reported explicitly.
- Added `Settle surface: <path>` on stderr for every proceed path before Verification starts, without changing settle stdout lines or integration mechanics.
- Added refusal output that names every candidate surface with its status, or that the path does not exist, while preserving the existing Implement Command guidance for non-failed Task statuses.
- Added regression coverage for the stale kept Run Worktree case where the kept worktree is `pending` and the checkout is `failed`; settle proceeds from the checkout and reports that surface.
- Added coverage that a failed kept Run Worktree still wins when it is the first failed candidate, and that no-failed-surface refusals enumerate all candidate statuses.
- Evidence: `rtk proxy grep -q "Settle surface:" internal/cli/settle.go` passed; `rtk go test ./internal/cli/ -run 'Settle'` passed with 26 settle tests; `rtk go test ./internal/cli/...` passed with 433 CLI tests; `rtk go build -buildvcs=false ./...` passed; `rtk make verify` passed with 1249 tests, skills check, and build.
