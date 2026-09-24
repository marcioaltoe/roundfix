---
task: task_05
spec: 0157-reconciliation-and-one-repository-identity
status: pending
type: backend
complexity: high
---

# Task 05: Keep each Run's checkout, share only the identity

## Overview

Corrective Task from the pre-PR review of 2026-09-24. `createRun` now stores every Run's `git_root` as the main checkout, but callers treat `run.GitRoot` as the checkout the Run operates in: `reviewRunTargetGuard` then watches the main checkout and ends a linked-worktree Run as `CheckoutMoved`, `stop` looks for the review session in the wrong directory, and `ActiveRunInGitRoot` now locks every worktree of the repository at once, which forbids the parallel delivery this Spec exists for. `RepositoryRoots` also fails every call when one `.git/worktrees/<id>` entry has no readable `gitdir`, which happens while `git worktree add` runs.

## Requirements

1. MUST store the checkout a Run was created in as its `git_root`, unchanged.
2. MUST keep `ActiveRunInGitRoot` scoped to the exact checkout, so Active Runs in two worktrees of one repository do not block each other.
3. MUST keep listing a linked worktree's Runs from the main checkout and sharing the artifact directory identity.
4. MUST make `RepositoryRoot` resolve the main worktree without enumerating linked worktrees.
5. MUST make `RepositoryRoots` skip a worktree entry whose `gitdir` is missing or unreadable instead of failing.

## Subtasks

- [ ] Implement the requirements above.
- [ ] Add a test for each acceptance criterion.

## Acceptance Criteria

- [ ] A Run created in a linked worktree keeps that worktree as `git_root` and is still listed from the main checkout.
- [ ] Active Runs in two worktrees of one repository coexist.
- [ ] An unreadable worktree entry neither fails identity resolution nor listing.

## Context

- instruction: `.agents/skills/implement-task/SKILL.md`
- interface: `internal/store/store.go`
- interface: `internal/config/config.go`

## Verification

- `out="$(go test -count=1 -v -run "^(TestCreateRunKeepsTheCheckoutGitRoot|TestActiveRunLockStaysPerCheckout|TestRepositoryRootsSkipsAnUnreadableWorktreeEntry|TestRepositoryRootDoesNotEnumerateWorktrees|TestRunsFromALinkedWorktreeAreListedFromTheMainCheckout|TestRepositoryIdentityIsSharedByLinkedWorktrees)$" ./internal/config ./internal/store 2>&1)" || { printf "%s\\n" "$out"; exit 1; }; for name in TestCreateRunKeepsTheCheckoutGitRoot TestActiveRunLockStaysPerCheckout TestRepositoryRootsSkipsAnUnreadableWorktreeEntry TestRepositoryRootDoesNotEnumerateWorktrees TestRunsFromALinkedWorktreeAreListedFromTheMainCheckout TestRepositoryIdentityIsSharedByLinkedWorktrees; do printf "%s\\n" "$out" | grep -q -- "--- PASS: $name" || exit 1; done` — expected: exit 0; before this Task four of the named cases do not exist, so the command fails.

## References

- [_techspec.md](_techspec.md) — Repository identity
