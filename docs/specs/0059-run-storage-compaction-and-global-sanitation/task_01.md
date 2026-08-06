---
task: task_01
spec: 0059-run-storage-compaction-and-global-sanitation
status: completed
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

## Result

Implemented the read-only `roundfix storage report` surface. The store now
measures SQLite table allocation and row counts, free pages, repository and
state Run rows plus Run Artifact bytes, and each unique recorded Artifact
Root. The command reads only Roundfix Home, so it does not require repository
discovery or Project Config. A dedicated immutable SQLite reader prevents the
otherwise read-triggered WAL and shared-memory sidecars.

Focused checks run after the final implementation edits:

- `rtk zsh -c 'GOCACHE="$PWD/.gocache" rtk go test ./internal/store -count=1 -run "^TestStorageReportReconcilesMeasuredTotalsWithoutMutation$"'`
  — passed (`1` test).
- `rtk zsh -c 'GOCACHE="$PWD/.gocache" rtk go test ./internal/cli -count=1 -run "^TestRunStorageReportOutsideGitRepository$"'`
  — passed (`1` test).
- `rtk zsh -c 'GOCACHE="$PWD/.gocache" rtk go test ./internal/cli -count=1 -run "^TestRunStorageReportRejectsFlags$"'`
  — passed (`1` test).
- `rtk zsh -c 'GOCACHE="$PWD/.gocache" rtk go test ./internal/cli -count=1 -run "^TestRunStorageReportWithoutDatabaseCreatesNothing$"'`
  — passed (`1` test).
- `rtk zsh -c 'GOCACHE="$PWD/.gocache" rtk go test ./internal/cli -count=1 -run "^TestRunGCDryRunListsEligibleRunsAndChangesNothing$"'`
  — passed (`1` existing GC regression test).
- `rtk zsh -c 'GOCACHE="$PWD/.gocache" rtk go test ./internal/store ./internal/cli -count=1 -run "^$"'`
  — exited `0`; both packages and their tests compiled without selecting a
  test.
- `rtk zsh -c 'GOCACHE="$PWD/.gocache" rtk go vet ./internal/store ./internal/cli'`
  — exited `0`.
- `rtk git diff --check` — exited `0`.

Acceptance evidence:

- Repository, state, table, and Artifact Root grouping: the store-focused test
  reconciles every repository/state/root Run-row sum to the measured `runs`
  table row count, reconciles repository and state Run Artifact byte sums, and
  exercises table and unique Artifact Root groups. The CLI-focused test checks
  that all four groups are emitted.
- Database and artifact reconciliation: the store-focused test sums table page
  bytes plus measured free-page bytes and compares that relation with the
  measured database file; it independently sums unique Artifact Root bytes and
  compares them with the measured artifact total. No recorded byte value is
  asserted.
- Declared tolerance: `StorageReportReconciliationToleranceBytes` is the
  database's measured SQLite page size. The report states that one page is
  allowed because page accounting and the independent file stat can land on
  adjacent page boundaries without taking a writer lock; the test asserts the
  tolerance-to-page-size relation rather than a literal.
- No mutation: the store-focused test hashes the complete Run Database and
  Artifact Root trees before and after reporting and proves the path sets and
  bytes are identical. The outside-Git CLI test also compares the database file
  before and after, and the empty-home CLI test proves the command does not
  create a database.
- Missing recorded roots: both store and CLI tests create a Run whose recorded
  repository and Artifact Root no longer exist and assert that both paths are
  reported with `missing` status.
- Outside Git: the CLI test runs with a work directory containing no `.git`
  entry and exits `0`; the empty-home companion does the same without creating
  state.

The commands under `## Verification` were not run; the Daemon owns those
commands and Task settlement.
