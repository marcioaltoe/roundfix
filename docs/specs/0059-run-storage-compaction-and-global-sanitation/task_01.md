---
task: task_01
spec: 0059-run-storage-compaction-and-global-sanitation
status: pending
type: backend
complexity: medium
---

# Task 01: Report where the bytes are

## Overview

Policy decisions about storage are being made without measurement. This slice
adds the read-only report: bytes and row counts grouped by repository, state,
table, and Artifact Root.

It leads the graph because it mutates nothing and because it defines how a byte
total is counted. The three Tasks that follow assert against that vocabulary
rather than each inventing their own.

## Requirements

1. MUST add a read-only storage report grouped by repository, state, table, and
   Artifact Root, exposing measured bytes and row counts.
2. MUST require no flags and be safe to run anywhere, including outside a Git
   repository.
3. MUST reconcile its grouped totals with the measured database file size and
   the measured artifact totals, within a tolerance the Task declares and
   justifies.
4. MUST mutate nothing: no row deleted, no file removed, no lock taken beyond
   reading.
5. MUST report a repository whose recorded root no longer exists as such,
   rather than omitting it or failing.
6. MUST NOT assert any recorded size in its tests. Assert the reconciliation
   relation, which holds at any size.

## Subtasks

- [ ] Add the measurement and its grouping.
- [ ] Add the command surface with no flags.
- [ ] Assert reconciliation within the declared tolerance.

## Acceptance Criteria

- [ ] The report groups by repository, state, table, and Artifact Root.
- [ ] Grouped totals reconcile with the measured file size and artifact totals
      within the declared tolerance, asserted as a relation.
- [ ] The declared tolerance is stated with the reason for its size.
- [ ] Running the report changes no byte, asserted before and after.
- [ ] A recorded root that no longer exists is reported, not omitted.
- [ ] The command runs outside a Git repository.

## Context

- interface: `internal/store/journal.go`
- interface: `internal/cli/gc.go`

## Verification

- `go build -buildvcs=false ./...` — expected: exit 0.
- `output="$(go test ./internal/store ./internal/cli -count=1 -run 'StorageReport|Reconcil' -v 2>&1)"; status=$?; printf '%s\n' "$output"; [ "$status" -eq 0 ] || exit "$status"; printf '%s\n' "$output" | grep -q -- "--- PASS"`
  — expected: exit 0; the report tests ran and passed, with the exit status
  preserved rather than hidden by the pipe.
- `go test ./internal/store ./internal/cli -count=1` — expected: exit 0.
- `go test -parallel 16 ./...` — expected: exit 0.

## References

- `_prd.md` → Core Feature 4; User Story 4; Success Metric 3.
- `_techspec.md` → Measurement, not literals; Build Order 1.
