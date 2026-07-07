---
task: task_01
spec: 0020-run-browser
status: completed
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

- [x] State filter type and query wiring, call-site migration
- [x] Limit applied after filtering, newest first
- [x] Table tests: each filter value, limit, filter+limit, empty results

## Acceptance Criteria

- [x] Seeding a mix of Active and terminal Runs, the active filter returns
      only non-terminal Runs, the terminal filter only terminal ones, and the
      all filter every Run — newest first in each case.
- [x] A limit of N returns the N newest matching Runs; `0` returns all.
- [x] An unset filter behaves as active.

## Verification

- `rtk go test ./internal/store/` — expected: all tests pass, including the
  new filter and limit tests.
- `rtk make verify` — expected: fmt-check, tests, skills check, and build all
  pass.

## References

`_prd.md` → Core Features 1, 4; Decisions (default Active). `_techspec.md` →
Interfaces: ListRunsQuery; Build Order 1.

## Result

Replaced `ListRunsQuery.ActiveOnly` with `States RunStateFilter`
(`StatesActive` zero value/default, `StatesTerminal`, `StatesAll`) and added
`Limit int` (`0` = unbounded) in `internal/store/store.go`. Classification
goes through the existing `IsTerminalState` predicate via the unexported
`RunStateFilter.matches` — no duplicated state lists. The limit is applied in
the scan loop after the state filter on the `created_at DESC, id DESC`
ordering, so it keeps the N newest matches. No schema or write-path changes.

Call sites migrated with behavior preserved:

- `internal/cli/runs.go` — `runs list` uses `StatesAll` by default and
  `StatesActive` under `--active`.
- `internal/cli/attach.go` — the attach picker passes `StatesAll` explicitly.

### Acceptance criteria evidence

All covered by `TestListRunsStateFilterAndLimit` in
`internal/store/store_test.go`, which seeds all 13 Run states (4 non-terminal,
9 terminal) in one repository plus a lone Active Run in a second repository:

- Active/terminal/all filters, newest first: subtests "active filter keeps
  only non-terminal Runs", "terminal filter keeps only terminal Runs", "all
  filter with zero limit keeps every Run" assert exact newest-first ID order.
- Limit of N returns the N newest matches, `0` returns all: subtests "limit
  bounds the newest Runs" (N=3), "limit applies after the state filter"
  (terminal, N=2), "limit above the match count keeps every match"; the
  zero-limit queries in the filter subtests return every match.
- Unset filter behaves as active: subtest "unset filter defaults to active"
  queries with only `GitRoot` set and gets exactly the non-terminal Runs.
- Empty results: subtest "filter matching nothing returns an empty slice"
  (terminal filter over a repo with only an Active Run) returns a non-nil
  empty slice.

### Verification evidence

- `rtk go test ./internal/store/ -run TestListRuns -v` — 11 passed (both
  pre-existing ListRuns tests plus the new table test and its 8 subtests).
- `rtk go test ./internal/store/ ./internal/cli/` — 380 passed;
  `rtk gofmt -l` clean on both packages.
- `rtk make verify` — exit 0: fmt-check, `go test ./...` (947 passed in 19
  packages), `roundfix skills check` passed, build succeeded.

### Follow-up notes

- The Run Browser entry point (Task with the TUI listing) can now query
  `ListRunsQuery{States: ..., Limit: ...}` directly; no store follow-ups.
