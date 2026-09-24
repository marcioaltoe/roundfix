---
task: task_05
spec: 0157-reconciliation-and-one-repository-identity
status: completed
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

## Result

Implemented checkout-local Run storage and locking while preserving shared
repository identity for artifact resolution and repository-scoped listing.
`RepositoryRoot` now resolves the main worktree directly from the current
checkout's Git metadata, and `RepositoryRoots` ignores an individual linked
worktree entry when its `gitdir` cannot be read.

Focused-check evidence:

- A pre-change focused run of
  `go test -count=1 -run '^(TestCreateRunKeepsTheCheckoutGitRoot|TestActiveRunLockStaysPerCheckout)$' ./internal/store`
  failed because Run creation rewrote the linked checkout to the main root and
  the second worktree could not acquire its lock. The same command exited 0
  after the implementation. This covers checkout-preserving `git_root`, exact
  `ActiveRunInGitRoot` and `ActiveSpecRun` lookup, and coexisting Active Runs.
- A pre-change focused run of
  `go test -count=1 -run '^(TestRepositoryRootDoesNotEnumerateWorktrees|TestRepositoryRootsSkipsAnUnreadableWorktreeEntry)$' ./internal/config`
  failed on both the non-directory worktree registry and the unreadable entry.
  The same command exited 0 after the implementation. This covers direct main
  worktree resolution and tolerant linked-worktree enumeration.
- `go test -count=1 -run '^TestRunsFromALinkedWorktreeAreListedFromTheMainCheckout$' ./internal/store`
  exited 0 with an unreadable sibling entry present, proving the linked Run is
  still listed from the main checkout while retaining its linked checkout in
  `git_root`.
- `go test -count=1 -run '^TestRepositoryIdentityIsSharedByLinkedWorktrees$' ./internal/config`
  exited 0, proving the main and linked checkouts still resolve the same
  Artifact Directory identity.
- `go test -count=1 ./internal/config` and
  `go test -count=1 ./internal/store` both exited 0 after the final code and
  test edits.

The Daemon-owned `## Verification` command was not run.

## Carry-forward provenance

- Source Run: `run_20260924T173241Z_9612fd1f6966c280`
- Source commit: `4fbd57febf82df70d3874e98d9b85164cf4c0b10`
