---
task: task_06
spec: 0027-review-loop-integrity
status: pending
type: backend
complexity: high
---

# Task 06: Execute review Runs in the user's checkout

## Overview

With the Branch Integrity Preflight guarding the branch, remove worktree isolation from review Runs: resolve and watch point the Daemon engine at the user's checkout, create no Run Worktree and no Run Branch, and drop the integration step — batch commits land directly on the PR Head Branch. Integration Pending ceases to exist as a review-Run outcome. Spec Runs (implement) keep their worktree contract unchanged.

## Requirements

1. MUST pass the user's repository root as the engine's working directory for resolve and watch, so batch snapshots, Verification, and commits all operate on the checked-out PR Head Branch.
2. MUST NOT create a Run Worktree or a Run Branch for review Runs; the Run row records the user root as its working directory.
3. MUST remove the review-Run integration step, the Integration Pending outcome path, and the kept/cleaned worktree messaging for resolve and watch, leaving the implement command's worktree contract untouched.
4. MUST extend review-Run Preflight Validation to require a clean tracked working tree, refusing with each dirty path and the stash-or-commit next action; untracked files remain allowed (ADR-0045).
5. MUST state in the failed-batch and unresolved-outcome reporting that uncommitted changes in the checkout are Agent work from this Run, since the tree was clean at start.
6. MUST leave fetch behavior unchanged (it already runs in the checkout).

## Subtasks

- [ ] Redirect the engine plan's working directory to the user root for resolve and watch
- [ ] Delete review-Run worktree creation, integration, Integration Pending propagation, and worktree messaging
- [ ] Add the clean-tracked-tree preflight check with actionable refusal output
- [ ] Update failed-batch/unresolved report wording about leftover Agent work
- [ ] Integration-test in a temp repository: a resolve cycle commits on the user branch with no `roundfix/run-*` branch created; a dirty tracked file refuses at preflight; implement still creates its worktree

## Acceptance Criteria

- [ ] After a stubbed resolve cycle in a temp repository, batch commits exist on the user branch and no `roundfix/run-*` ref exists
- [ ] A dirty tracked file causes exit 2 with the file named and the next action stated; an untracked file does not block
- [ ] No review-Run code path can produce the Integration Pending outcome
- [ ] Implement-command worktree tests still pass unchanged
- [ ] The full test suite passes

## Context

- interface: `internal/cli/cli.go`
- interface: `internal/cli/implement.go`
- interface: `internal/daemon/engine.go`
- interface: `internal/worktree/worktree.go`
- interface: `internal/preflight/preflight.go`

## Verification

- `go test ./internal/cli/... ./internal/daemon/... ./internal/worktree/...` — expected: all tests pass
- `go test ./...` — expected: all tests pass
- `go build ./...` — expected: clean build

## References

`_prd.md` → Goal 1, User Story 1, Core Feature 1; `_techspec.md` → Build Order 5, System Architecture (internal/daemon seam), Executive Summary trade-off; ADR-0042, ADR-0045.
