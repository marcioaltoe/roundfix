---
task: task_01
spec: 0020-run-browser
status: pending
type: backend
complexity: low
---

# Task 01: Run state filter and limit on the store listing query

## Overview

Generalize the Run Database listing query: replace the boolean active-only
filter with a three-way state filter (active, terminal, all) and add a result
bound. Verifiable on its own through store package tests.

## Requirements

1. MUST replace `ListRunsQuery.ActiveOnly` with a state filter accepting
   active, terminal, and all, defaulting to active when unset, and migrate
   every call site.
2. MUST add `Limit` to the query (`0` = unbounded), applied after the state
   filter, newest first.
3. MUST classify active versus terminal through the store's existing terminal
   state predicate — no duplicated state lists.
4. MUST NOT change the schema or any write path.

## Subtasks

- [ ] State filter type and query wiring, call-site migration
- [ ] Limit applied after filtering, newest first
- [ ] Table tests: each filter value, limit, filter+limit, empty results

## Acceptance Criteria

- [ ] Seeding a mix of Active and terminal Runs, the active filter returns
      only non-terminal Runs, the terminal filter only terminal ones, and the
      all filter every Run — newest first in each case.
- [ ] A limit of N returns the N newest matching Runs; `0` returns all.
- [ ] An unset filter behaves as active.

## Verification

- `rtk go test ./internal/store/` — expected: all tests pass, including the
  new filter and limit tests.
- `rtk make verify` — expected: fmt-check, tests, skills check, and build all
  pass.

## References

`_prd.md` → Core Features 1, 4; Decisions (default Active). `_techspec.md` →
Interfaces: ListRunsQuery; Build Order 1.
