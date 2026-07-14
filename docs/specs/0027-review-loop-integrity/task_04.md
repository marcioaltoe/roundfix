---
task: task_04
spec: 0027-review-loop-integrity
status: pending
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

- [ ] Implement the enumeration helper returning the pending-work report structure
- [ ] Implement fast-forward-only integration with a descriptive refusal error
- [ ] Cover enumeration in tests: no run branches, run branch with zero ahead commits, ahead + ff-able, ahead + diverged, branch with and without a worktree
- [ ] Cover integration in tests: successful fast-forward, refusal on divergence

## Acceptance Criteria

- [ ] Enumeration returns exactly the pending Run Branches for a prepared temp repository, with correct ahead counts and ff flags
- [ ] Fast-forward integration moves the base branch and refuses divergent branches without mutating anything
- [ ] The full test suite passes, including with the race detector if goroutines are involved (none expected)

## Context

- interface: `internal/worktree/worktree.go`
- interface: `internal/gitcmd/gitcmd.go`

## Verification

- `rg -q "PendingRunWork" internal/worktree` — expected: exit 0 (report type exists)
- `go test ./internal/worktree/...` — expected: all tests pass
- `go build ./...` — expected: clean build

## References

`_prd.md` → Goal 1, User Story 2, Core Feature 2; `_techspec.md` → Build Order 4 (helper half), Interfaces (PendingRunWork), Integration Points (Git); ADR-0024, ADR-0042.
