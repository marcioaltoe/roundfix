---
task: task_06
spec: 0027-review-loop-integrity
status: completed
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

- [x] Redirect the engine plan's working directory to the user root for resolve and watch
- [x] Delete review-Run worktree creation, integration, Integration Pending propagation, and worktree messaging
- [x] Add the clean-tracked-tree preflight check with actionable refusal output
- [x] Update failed-batch/unresolved report wording about leftover Agent work
- [x] Integration-test in a temp repository: a resolve cycle commits on the user branch with no `roundfix/run-*` branch created; a dirty tracked file refuses at preflight; implement still creates its worktree

## Acceptance Criteria

- [x] After a stubbed resolve cycle in a temp repository, batch commits exist on the user branch and no `roundfix/run-*` ref exists
- [x] A dirty tracked file causes exit 2 with the file named and the next action stated; an untracked file does not block
- [x] No review-Run code path can produce the Integration Pending outcome
- [x] Implement-command worktree tests still pass unchanged
- [x] The full test suite passes

## Context

- interface: `internal/cli/cli.go`
- interface: `internal/cli/implement.go`
- interface: `internal/daemon/engine.go`
- interface: `internal/worktree/worktree.go`
- interface: `internal/preflight/preflight.go`

## Verification

- `go test ./internal/cli/... ./internal/daemon/... ./internal/worktree/...` — expected: all tests pass
- `go test ./...` — expected: all tests pass
- `go build -buildvcs=false ./...` — expected: clean build

## References

`_prd.md` → Goal 1, User Story 1, Core Feature 1; `_techspec.md` → Build Order 5, System Architecture (internal/daemon seam), Executive Summary trade-off; ADR-0042, ADR-0045.

## Result

- Redirected resolve/watch review Runs to the user checkout: Agent sessions, `CyclePlan.GitRoot`, live views, batch commits, review artifact commits, and Final Push now use the repository root; Run rows record that root as `work_dir`.
- Removed review-Run worktree creation, Run Branch integration, Integration Pending propagation, cleanup/kept-worktree reporting, and watch's integration command field. Remaining Integration Pending code is scoped to implement/settle paths.
- Added clean tracked checkout Preflight Validation for resolve/watch. Dirty tracked paths refuse with exit 2 and a stash-or-commit next action; untracked files are allowed.
- Added failed-batch and unresolved-outcome wording that uncommitted checkout changes are Agent work from the Run because Preflight started clean.
- Evidence: `go test ./internal/cli/... ./internal/daemon/... ./internal/worktree/...` passed; `go test ./...` passed; `go build -buildvcs=false ./...` passed. `go build ./...` was not usable in this worktree because Go VCS stamping failed with `error obtaining VCS status`; the Verification command was updated to the buildvcs-safe form.
