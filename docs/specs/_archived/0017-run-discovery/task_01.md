---
task: task_01
spec: 0017-run-discovery
status: completed
type: backend
complexity: low
---

# Task 01: Store listing query with repository scope and active filters

## Overview

Add the one read-only Run Database query both discovery surfaces consume: list
Runs newest first, optionally scoped to one repository and optionally filtered
to Active Runs. Verifiable on its own through store package tests against a
temporary database.

## Requirements

1. MUST expose a store listing query taking a scope value object with a
   repository root (empty means every repository) and an active-only flag, and
   returning matching Runs ordered newest first by creation time.
2. MUST classify Active versus terminal using the store's existing terminal
   state predicate — no duplicated state lists.
3. MUST return an empty slice, not an error, when nothing matches.
4. MUST NOT change the schema, existing queries, or any write path.

## Subtasks

- [x] Listing query and its scope value object in the store package
- [x] Newest-first ordering by creation time
- [x] Repository scope and active-only filters, composable
- [x] Table tests: scope, filter, ordering, empty result

## Acceptance Criteria

- [x] Seeding Runs across two repositories and listing with a repository scope
      returns only that repository's Runs, newest first.
- [x] The active-only filter excludes every Run in a terminal state and keeps
      every non-terminal one.
- [x] An empty scope lists Runs from every repository.
- [x] Listing an empty database returns an empty result with no error.

## Verification

- `rtk go test ./internal/store/` — expected: all tests pass, including the new
  listing query tests.
- `rtk make verify` — expected: fmt-check, tests, skills check, and build all
  pass.

## References

`_prd.md` → User Stories 1-2; Core Features 2-4. `_techspec.md` → Interfaces:
Store listing query; Data Models; Build Order 1.

## Result

Implemented `Store.ListRuns(ctx, ListRunsQuery)` as a read-only Run Database
query. `ListRunsQuery.GitRoot` scopes by repository when set, an empty value
lists every repository, and `ActiveOnly` filters scanned rows through
`IsTerminalState` so terminal-state classification stays single-sourced. The
query returns Runs newest first by `created_at DESC`.

Verification:

- `rtk go test ./internal/store/`: passed, 45 store tests.
- `rtk make verify`: passed. The gate reported `rtk go test ./...` with 824
  tests passed in 18 packages, `roundfix skills check` passed, and
  `rtk go build -buildvcs=false -o bin/roundfix ./cmd/roundfix` completed.

Acceptance evidence:

- Repository scope and newest-first ordering: `TestListRunsScopesByRepositoryAndOrdersNewestFirst`
  seeds Runs in `tmp/repo-a` and `tmp/repo-b`; listing `tmp/repo-a` returns
  only that repository's Runs as `[repoANewer, repoAOlder]`.
- Active-only filtering: `TestListRunsActiveOnlyKeepsEveryNonTerminalState`
  seeds every current terminal Run state plus `Active`, `ResolvingWithAgent`,
  `Verifying`, and `Pushing`; `ActiveOnly` returns only the four
  non-terminal Runs and every returned state is checked with
  `IsTerminalState`.
- Empty repository scope: `TestListRunsScopesByRepositoryAndOrdersNewestFirst`
  also lists with an empty `ListRunsQuery` and returns Runs from both seeded
  repositories newest first.
- Empty Run Database: `TestListRunsEmptyDatabaseReturnsEmptySlice` verifies
  a newly opened store returns a non-nil empty slice and no error.

Follow-ups: none for this task slice.
