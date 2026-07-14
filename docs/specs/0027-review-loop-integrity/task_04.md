---
task: task_04
spec: 0027-review-loop-integrity
status: completed
type: backend
complexity: medium
---

# Task 04: Enumerate pending Run Branch work and fast-forward-integrate it

## Overview

The Branch Integrity Preflight needs a deterministic answer to "what Roundfix Run Branch work is pending against this branch, and can it be integrated safely?". This task builds that capability in the worktree package as pure, read-mostly git operations: enumeration of `roundfix/run-*` branches based on a given branch (with ahead-commit counts, associated worktree paths, and fast-forward eligibility) and a fast-forward-only integration operation.

## Requirements

1. MUST enumerate every Roundfix Run Branch whose merge base with a given branch is that branch's tip ancestor — reporting branch name, associated worktree path when one exists, number of commits ahead of the base branch, and whether the base branch can fast-forward to the Run Branch tip.
2. MUST provide a fast-forward-only integration operation that advances the base branch to a Run Branch tip and refuses with a descriptive error when fast-forward is impossible; it never creates merge commits and never resolves conflicts.
3. MUST use porcelain git commands only, per the existing integration decision (ADR-0024), reusing the package's existing branch-enumeration and worktree-list parsing patterns.
4. MUST treat a Run Branch with zero ahead commits as not pending.

## Subtasks

- [x] Implement the enumeration helper returning the pending-work report structure
- [x] Implement fast-forward-only integration with a descriptive refusal error
- [x] Cover enumeration in tests: no run branches, run branch with zero ahead commits, ahead + ff-able, ahead + diverged, branch with and without a worktree
- [x] Cover integration in tests: successful fast-forward, refusal on divergence

## Acceptance Criteria

- [x] Enumeration returns exactly the pending Run Branches for a prepared temp repository, with correct ahead counts and ff flags
- [x] Fast-forward integration moves the base branch and refuses divergent branches without mutating anything
- [x] The full test suite passes, including with the race detector if goroutines are involved (none expected)

## Context

- interface: `internal/worktree/worktree.go`
- interface: `internal/gitcmd/gitcmd.go`

## Verification

- `grep -R -q "PendingRunWork" internal/worktree` — expected: exit 0 (report type exists; `rg` is unavailable in this environment)
- `go test ./internal/worktree/...` — expected: all tests pass
- `go build -buildvcs=false ./...` — expected: clean build

## References

`_prd.md` → Goal 1, User Story 2, Core Feature 2; `_techspec.md` → Build Order 4 (helper half), Interfaces (PendingRunWork), Integration Points (Git); ADR-0024, ADR-0042.

## Result

- Added `PendingRunWork`, `ListPendingRunWork`, and `IntegratePendingRunWork` in `internal/worktree`, using the existing git runner plus `worktree list --porcelain`, `for-each-ref`, ancestry checks, and fast-forward porcelain.
- Pre-change signal: `rtk grep -R -n "PendingRunWork" internal/worktree` exited 1 before implementation, proving the report type/helper surface was absent.
- Acceptance 1: `TestListPendingRunWorkReportsAheadRunBranches` creates a temp repository covering no Run Branches, zero-ahead, fast-forwardable ahead work with a worktree path, and diverged ahead work without a worktree; `rtk go test ./internal/worktree/...` passed with 24 tests.
- Acceptance 2: `TestIntegratePendingRunWorkFastForwardsBaseBranch` verifies the base branch moves to the Run Branch tip; `TestIntegratePendingRunWorkRefusesDivergedBranch` verifies divergence returns a fast-forward refusal and leaves both refs unchanged.
- Acceptance 3: `rtk make verify` passed: 1160 Go tests, `roundfix skills check`, and `go build -buildvcs=false -o bin/roundfix ./cmd/roundfix`. No goroutines were added, so race-specific verification was not applicable.
- Verification run: `rtk grep -R -q "PendingRunWork" internal/worktree`, `rtk go test ./internal/worktree/...`, and `rtk go build -buildvcs=false ./...` all exited 0.
