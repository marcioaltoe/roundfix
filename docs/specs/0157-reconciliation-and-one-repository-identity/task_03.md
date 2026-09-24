---
task: task_03
spec: 0157-reconciliation-and-one-repository-identity
status: pending
type: backend
complexity: high
---

# Task 03: One identity for every worktree

## Overview

`repoID` hashes the checkout path and the Run store keys Runs by the Git root, so each linked worktree is a separate repository: separate artifact directories, and Runs invisible from the main checkout.

## Requirements

1. MUST derive the repository identity from the main worktree's root, resolved through the common Git directory.
2. MUST keep the main checkout's identity equal to its current value.
3. MUST make Runs recorded from a linked worktree visible from the main checkout.
4. MUST keep artifact directories created under a worktree-derived identity readable and listed with the repository's Runs.

## Subtasks

- [ ] Implement the requirements above.
- [ ] Add a test for each acceptance criterion.

## Acceptance Criteria

- [ ] A linked worktree and the main checkout report the same identity.
- [ ] The main checkout's identity is unchanged.
- [ ] A Run recorded from a linked worktree is listed from the main checkout.
- [ ] An earlier worktree-derived artifact directory remains readable.

## Context

- instruction: `.agents/skills/implement-task/SKILL.md`
- interface: `internal/config/config.go`
- interface: `internal/store/store.go`

## Verification

- `out="$(go test -count=1 -v -run "^TestRepositoryIdentityIsSharedByLinkedWorktrees$" ./internal/config 2>&1)" || { printf "%s\n" "$out"; exit 1; }; printf "%s\n" "$out" | grep -q -- "--- PASS: TestRepositoryIdentityIsSharedByLinkedWorktrees"` — expected: exit 0; before this Task the case does not exist, so the command fails.
- `out="$(go test -count=1 -v -run "^TestRunsFromALinkedWorktreeAreListedFromTheMainCheckout$" ./internal/store 2>&1)" || { printf "%s\n" "$out"; exit 1; }; printf "%s\n" "$out" | grep -q -- "--- PASS: TestRunsFromALinkedWorktreeAreListedFromTheMainCheckout"` — expected: exit 0; before this Task the case does not exist, so the command fails.

## References

- [_techspec.md](_techspec.md)
