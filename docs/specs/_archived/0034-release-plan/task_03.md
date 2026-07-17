---
task: task_03
spec: 0034-release-plan
status: completed
type: backend
complexity: medium
---

# Task 03: Resolve committed release ranges through local Git

## Overview

Supply the classifier with a reproducible committed range from the local repository through the existing context-aware Git boundary. Temporary-repository integration tests make tag selection, revision validation, changed paths, cleanliness, and read-only behavior observable without a network or Roundfix Run.

## Requirements

1. MUST implement the TechSpec `GitSource` contract through the existing context-aware Git runner rather than shelling out inside the domain package.
2. MUST resolve explicit `--from` and `--to` identities and default the target to committed `HEAD`.
3. MUST select the latest stable semantic-version tag reachable from the target when the base is omitted.
4. MUST return normalized commits with full messages and changed paths for the complete non-reversed range.
5. MUST reject missing tags, malformed bases, unresolved revisions, reversed or empty ranges, and non-commit targets with actionable errors.
6. MUST reject tracked or untracked working-tree changes and name the paths plus the commit, stash, or remove next action.
7. MUST honor context cancellation, perform no network access, and leave files, refs, tags, remotes, and configuration unchanged.

## Subtasks

- [x] Adapt the existing Git runner to the Release Plan source interface.
- [x] Resolve default and explicit release-range endpoints.
- [x] Load complete commit messages and changed paths deterministically.
- [x] Detect dirty, missing, malformed, empty, and reversed inputs.
- [x] Snapshot repository state around read-only operations.
- [x] Add real temporary-repository integration fixtures and tests.

## Acceptance Criteria

- [x] Omitted endpoints resolve the latest reachable stable tag through committed `HEAD`.
- [x] Explicit valid endpoints return the expected ordered commits and paths.
- [x] Dirty tracked and untracked paths block resolution with an actionable diagnostic payload.
- [x] Invalid tags, revisions, targets, and ranges fail before any classification result is emitted.
- [x] Repository bytes, refs, tags, remotes, and configuration are identical before and after every successful or failing source operation.

## Context

- interface: `internal/preflight/preflight.go`
- interface: `internal/cli/cli.go`

## Verification

- `go test ./internal/cli -run 'TestReleasePlanGitSource' -count=1` — expected: real temporary-repository range, dirty-tree, invalid-input, cancellation, and mutation-audit cases pass.
- `go test ./internal/releaseplan ./internal/cli -run 'Test(ReleasePlanGitSource|Classify|ParseStableVersion|CalculateProposal)' -count=1` — expected: the Git adapter and domain contract pass together.

## References

- `_prd.md` → Goal 1; User Stories 1 and 3; Core Features 1, 4, and 7; User Experience.
- `_techspec.md` → Interfaces; API Contracts: range defaults and invalid inputs; Testing Approach; Build Order 3.
- ADR-0048 → Release planning is read-only and confirmation-gated.

## Result

Implemented the Release Plan `GitSource` contract and a CLI adapter backed by the existing context-aware `preflight.GitRunner`. The adapter resolves clean local repositories only, defaults `--to` to committed `HEAD`, selects the highest reachable stable release tag when `--from` is omitted, validates explicit stable base tags and commit targets, loads ordered commit messages and changed paths, and returns typed actionable errors for dirty worktrees and invalid ranges.

Verification:

- `go test ./internal/cli -run 'TestReleasePlanGitSource' -count=1` — passed after implementation: 13 temporary-repository GitSource cases.
- `go test ./internal/releaseplan ./internal/cli -run 'Test(ReleasePlanGitSource|Classify|ParseStableVersion|CalculateProposal)' -count=1` — passed after implementation: 57 combined adapter and domain cases.
- `make verify` — passed after implementation: Go tests, Python skill tests, `roundfix skills check`, and build all completed.

Acceptance evidence:

- Omitted endpoints select the latest reachable stable tag through `HEAD`: `TestReleasePlanGitSourceDefaultRangeResolvesLatestReachableStableTag`.
- Explicit endpoints return ordered commits and deterministic changed paths: `TestReleasePlanGitSourceExplicitEndpointsReturnOrderedCommitsAndPaths`.
- Dirty tracked and untracked paths return `ErrDirtyWorktree`, path payloads, and commit/stash/remove guidance: `TestReleasePlanGitSourceDirtyWorktreeReportsPathsAndPreservesRepo`.
- Malformed bases, pre-release bases, missing tags, unresolved targets, non-commit targets, empty ranges, and reversed ranges fail with typed errors and no partial range: `TestReleasePlanGitSourceInvalidInputsFailReadOnly`.
- Successful and failing operations preserve worktree bytes, refs, tags, remotes, local config, and status snapshots: all `TestReleasePlanGitSource*` tests use the mutation audit helper.

Follow-ups:

- Build orchestration, CLI parsing, output rendering, and exit-code mapping remain in later tasks.
