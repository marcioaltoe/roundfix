---
task: task_01
spec: 0250-a-run-database-that-keeps-only-recent-runs
status: pending
type: data
complexity: high
---

# Task 01: The Run Database removes a terminal Run whole and compacts incrementally

## Overview

Adds the store half of Run Retention to `internal/store`: schema version 23
with the one-row sweep record, incremental auto-vacuum for a new Run Database,
the candidate scan, the per-Run removal that re-checks the Run inside its
write transaction, the incremental compaction slice, and the conversion of a
default-mode database by the guarded full compaction. The Run Database
lifecycle policy gains the new table and names Run Retention as the age bound
of the Run rows. It is verifiable on its own through fixture databases in
temporary homes.

## Requirements

1. MUST raise the schema to version 23 with the `run_retention_sweeps` table of
   `_techspec.md` → Data Models, created on a fresh database and by a v22 to
   v23 migration, and MUST make a database this binary creates report
   `PRAGMA auto_vacuum` 2. The TechSpec's Current behavior section records
   that the pragma must reach a new file before `journal_mode = WAL`, or be
   followed by a `VACUUM`. An existing database keeps its mode.
2. MUST implement `RunRetentionCandidates`, `RemoveRetainedRun`,
   `LastRunRetentionSweep`, `RecordRunRetentionSweep` and
   `CompactIncrementally` with the signatures of `_techspec.md` → Interfaces,
   as `_techspec.md` → API Contract 2 and API Contract 3 state:
   - the scan selects terminal Runs completed before the cutoff, never an
     Active Run, sets `QueueReferenced` from `delivery_queue_items.run_id` and
     `delivery_queue_runs.run_id`, fills the row counts of `runs`,
     `run_events`, `run_agent_selections`, `run_token_usage` and
     `active_run_locks`, and sets `EstimatedBytes` to the summed byte length
     of the events' `payload` and `summary`; it runs on a reader;
   - the removal is one `withWriteTx` transaction that re-reads the Run,
     returns `RunRetentionKeptError` with a reason when it is no longer
     terminal, no longer before the cutoff or queue-referenced, deletes the
     `runs` row and returns the rows the cascade removed; an absent Run returns
     zero counts and no error;
   - the compaction runs `PRAGMA incremental_vacuum(<maxPages>)` in one write
     transaction when the database is in incremental mode, returns the pages
     released and the free pages left, and does nothing and reports
     `Incremental: false` otherwise.
3. MUST make `PreviewCompaction` and `Compact` set `auto_vacuum = INCREMENTAL`
   before `VACUUM INTO` and `VACUUM`, so a default-mode database is converted
   by `Compact` and its result still matches the preview within the existing
   one-page tolerance. Their refusals and errors are unchanged.
4. MUST NOT change `PruneTerminalRuns`, `TerminalRunPruneCandidates`, the
   Journal Retention semantics, or any table other than the new one
   (`_techspec.md` → Invariants 6 and 10).
5. MUST update the table between the `durable-table-lifecycle` markers of
   `docs/user-guide/run-database-lifecycle.md`: add the `run_retention_sweeps`
   row, owned by the Run Retention lifecycle, keeping one row that records the
   last completed sweep; and change the retention rule of `runs`,
   `run_events`, `run_agent_selections`, `run_token_usage` and
   `active_run_locks` so each says that Run Retention removes it with its
   terminal Run past the window (ADR-0255). The lifecycle policy test must
   pass with the new table.
6. MUST add these tests, each on a fixture Run Database in a `t.TempDir()`
   home:
   - `TestRunRetentionCandidatesCountEveryDependentRow` in
     `internal/store/run_retention_test.go`: an old terminal Run with events,
     an Agent Selection and a token usage row, an Active Run, a recent Run and
     a queue-referenced old Run give exactly the old and queue-referenced
     candidates, with the queue flag, exact row counts and a byte estimate
     equal to the summed lengths;
   - `TestRemoveRetainedRunDeletesTheRunWhole`: after removal no row of that
     Run is left in any of the five tables, every other Run keeps every row,
     and `run_windows`, `interactive_defaults` and the Delivery Queue tables
     are unchanged;
   - `TestRemoveRetainedRunRefusesARunThatChanged`: a Run linked to a
     Delivery Queue after the scan, or moved inside the window, is refused
     with `RunRetentionKeptError` and keeps every row;
   - `TestRunRetentionSweepRecordRoundTrips`: no record, then a recorded
     sweep read back with its time and window, then a replacement;
   - `TestNewRunDatabaseUsesIncrementalAutoVacuum`: a database created by
     `Open` reports mode 2;
   - `TestCompactIncrementallyShrinksTheFile`: after a removal on a fresh
     database, repeated slices bring the free pages to 0 and the page count
     below its value before the removal; a default-mode database reports
     `Incremental: false` and is unchanged;
   - `TestCompactConvertsADefaultModeDatabase`: a database created in mode 0
     is in mode 2 after `PreviewCompaction` and `Compact`, the preview stays
     read-only, and the result matches the preview within one page;
   - `TestMigrateAddsRunRetentionSweeps` in `internal/store/migrate_test.go`:
     a v22 database migrates to v23 with the new table and every existing
     row intact.

## Subtasks

- [ ] Add schema 23, the migration and incremental mode for a new database.
- [ ] Add the candidate scan, the removal and the sweep record.
- [ ] Add the incremental compaction slice and the conversion in `Compact`.
- [ ] Update the lifecycle policy table.
- [ ] Write the eight tests and run the whole store package.

## Acceptance Criteria

- [ ] A terminal Run past the cutoff is removed with every dependent row, and nothing else changes.
- [ ] A queue-referenced, Active or changed Run is never removed.
- [ ] A new database compacts incrementally, and `Compact` converts an old one.
- [ ] The lifecycle policy names every table, the new one included.

## Context

- interface: `internal/store/store.go`
- interface: `internal/store/journal.go`
- interface: `internal/store/journal_test.go`
- interface: `internal/store/migrate_test.go`
- interface: `docs/user-guide/run-database-lifecycle.md`
- creates: `internal/store/run_retention.go`
- creates: `internal/store/run_retention_test.go`
- instruction: `docs/adr/0255-run-retention-removes-terminal-runs-whole-and-compacts-incrementally.md`
- instruction: `docs/adr/0033-the-run-event-journal-is-pruned-by-retention.md`
- instruction: `docs/adr/0171-a-retention-prune-reports-only-what-it-reclaimed.md`

## Verification

- `out="$(go test -count=1 -v -run '^(TestRunRetentionCandidatesCountEveryDependentRow|TestRemoveRetainedRunDeletesTheRunWhole|TestRemoveRetainedRunRefusesARunThatChanged|TestRunRetentionSweepRecordRoundTrips|TestNewRunDatabaseUsesIncrementalAutoVacuum|TestCompactIncrementallyShrinksTheFile|TestCompactConvertsADefaultModeDatabase|TestMigrateAddsRunRetentionSweeps)$' ./internal/store 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for t in TestRunRetentionCandidatesCountEveryDependentRow TestRemoveRetainedRunDeletesTheRunWhole TestRemoveRetainedRunRefusesARunThatChanged TestRunRetentionSweepRecordRoundTrips TestNewRunDatabaseUsesIncrementalAutoVacuum TestCompactIncrementallyShrinksTheFile TestCompactConvertsADefaultModeDatabase TestMigrateAddsRunRetentionSweeps; do printf '%s\n' "$out" | grep -q -- "--- PASS: $t " || { printf 'missing passing test %s\n' "$t" >&2; exit 1; }; done; p="$(tr -s '[:space:]' ' ' < docs/user-guide/run-database-lifecycle.md)"; printf '%s\n' "$p" | grep -q -- '[|] .run_retention_sweeps. [|]' || { printf 'missing the run_retention_sweeps row in the lifecycle policy\n' >&2; exit 1; }; go test -count=1 ./internal/store` — expected: exit 0. Before this Task the eight tests do not exist and the policy has no `run_retention_sweeps` row, so the first check fails. After it, the eight tests pass, the policy names the new table, and every `./internal/store` test passes, the lifecycle policy test included.

## References

- `_prd.md` → Core Feature 1; Core Feature 2; Core Feature 4; Success Metric 1; Success Metric 4
- `_techspec.md` → Data Models; API Contract 2; API Contract 3; Interfaces; Invariants 1, 2, 5, 6, 10; Current behavior; Build Order 1
- ADR-0255; ADR-0033; ADR-0171
