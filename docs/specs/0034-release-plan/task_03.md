---
task: task_03
spec: 0034-release-plan
status: pending
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

- [ ] Adapt the existing Git runner to the Release Plan source interface.
- [ ] Resolve default and explicit release-range endpoints.
- [ ] Load complete commit messages and changed paths deterministically.
- [ ] Detect dirty, missing, malformed, empty, and reversed inputs.
- [ ] Snapshot repository state around read-only operations.
- [ ] Add real temporary-repository integration fixtures and tests.

## Acceptance Criteria

- [ ] Omitted endpoints resolve the latest reachable stable tag through committed `HEAD`.
- [ ] Explicit valid endpoints return the expected ordered commits and paths.
- [ ] Dirty tracked and untracked paths block resolution with an actionable diagnostic payload.
- [ ] Invalid tags, revisions, targets, and ranges fail before any classification result is emitted.
- [ ] Repository bytes, refs, tags, remotes, and configuration are identical before and after every successful or failing source operation.

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
