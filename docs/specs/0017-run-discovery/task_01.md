---
task: task_01
spec: 0017-run-discovery
status: pending
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

- [ ] Listing query and its scope value object in the store package
- [ ] Newest-first ordering by creation time
- [ ] Repository scope and active-only filters, composable
- [ ] Table tests: scope, filter, ordering, empty result

## Acceptance Criteria

- [ ] Seeding Runs across two repositories and listing with a repository scope
      returns only that repository's Runs, newest first.
- [ ] The active-only filter excludes every Run in a terminal state and keeps
      every non-terminal one.
- [ ] An empty scope lists Runs from every repository.
- [ ] Listing an empty database returns an empty result with no error.

## Verification

- `rtk go test ./internal/store/` — expected: all tests pass, including the new
  listing query tests.
- `rtk make verify` — expected: fmt-check, tests, skills check, and build all
  pass.

## References

`_prd.md` → User Stories 1-2; Core Features 2-4. `_techspec.md` → Interfaces:
Store listing query; Data Models; Build Order 1.
