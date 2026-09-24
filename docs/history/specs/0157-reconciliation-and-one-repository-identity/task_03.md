---
task: task_03
spec: 0157-reconciliation-and-one-repository-identity
status: completed
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

## Result

### Implementation

- Default Artifact Directory identity now resolves a linked worktree through
  its common Git directory and hashes the main worktree root. A main checkout
  continues to hash its existing root.
- Run creation stores that same main-worktree root. Repository-scoped Run
  lookups include the main root and every path registered in the common Git
  directory, which keeps earlier checkout-derived rows discoverable without
  rewriting their recorded Artifact Directory.
- Active Spec Run, Run Window, retained Spec Run, and repository-scoped active
  Run lookups use the same repository-root resolution.

### Focused checks

- Before the production change,
  `rtk env GOCACHE=/tmp/roundfix-task03-gocache go test ./internal/config ./internal/store`
  exited 1: the linked worktree resolved a different Artifact Directory, both
  main-checkout Run listings were empty, and the earlier artifact was absent
  from the listing.
- After the production change, the same focused command exited 0 for both
  packages.
- The first `rtk make verify-incremental` run reached the full suite; Task 03's
  packages passed, but two force-stop integration tests could not read the
  sandbox process table. The permission-enabled rerun exited 0 after
  formatting, vet, all Go package tests, skill checks, and the build.
- The Task's authored `## Verification` commands were not run; the Daemon owns
  those checks and Task settlement.

### Acceptance-criterion evidence

- **A linked worktree and the main checkout report the same identity:**
  `TestRepositoryIdentityIsSharedByLinkedWorktrees` creates a real Git main
  checkout and linked worktree, then observes the same default Artifact
  Directory from both.
- **The main checkout's identity is unchanged:** the same config test compares
  the resolved main-checkout directory with the pre-change hash of the main
  root.
- **A Run recorded from a linked worktree is listed from the main checkout:**
  `TestRunsFromALinkedWorktreeAreListedFromTheMainCheckout` records through the
  linked path, lists through the main path, and observes the Run under the main
  repository root.
- **An earlier worktree-derived Artifact Directory remains readable:**
  `TestEarlierWorktreeDerivedArtifactDirectoryRemainsReadable` seeds an earlier
  linked-path Run and its path-hashed artifact file, lists it from the main
  checkout, and reads the file through the Run's preserved Artifact Directory.
