---
task: task_03
spec: 0081-a-journal-cheap-to-write-and-keep
status: completed
type: data
complexity: medium
---

# Task 03: Stop scanning the journal at every Run start

## Overview

The retention sweep runs at the start of every `implement`, `watch`, and
`resolve`. Its eligibility query joins the entire event table and groups by
Run with no time predicate in SQL — the cutoff is applied in Go, after the
scan — and it runs inside the write transaction. So every Run start takes the
machine-wide write lock and holds it across a full scan of the largest table,
even when nothing is eligible, which the measurement showed is the steady
state: 302 eligible Runs, zero prunable rows.

The event count that scan produces is used for reporting only.

## Requirements

1. MUST express the retention cutoff as a SQL predicate so eligibility work is
   bounded by the candidate set rather than by the table.
2. MUST move the eligibility scan out of the write transaction; only the prune
   itself needs the writer, and only when rows are actually eligible.
3. MUST either derive the reported event count cheaply from the candidate set
   or drop it from the hot path, and MUST say which in the Task Result.
4. MUST preserve the retention policy's meaning exactly: age-based,
   terminal-only, `runs` rows and Active Run locks never touched, and zero
   meaning keep everything.
5. MUST keep `roundfix gc` output shapes and exit behaviour unchanged for
   every case that has rows, and MUST state any output change for the
   no-eligible-rows case.
6. MUST NOT change the schema or vacuum anything — compaction stays the
   separate, explicitly fenced command it already is.

The declared Verification names `TestRetentionScanIsBoundedByCandidates`, which does not exist yet, so it can
fail before the work. Create it to assert that eligibility work is bounded by the candidate set rather than the event table, with no aggregate over that table inside a write transaction. A broad pattern over
this package matches cases that already pass and would approve the Task
before it starts.

## Subtasks

- [ ] Move the cutoff into SQL and the scan out of the write transaction.
- [ ] Resolve the reported event count cheaply or drop it.
- [ ] Prove policy meaning and `gc` behaviour unchanged.

## Acceptance Criteria

- [ ] Eligibility work is bounded by the candidate set, not the event table.
- [ ] No aggregate over the event table runs inside a write transaction.
- [ ] The retention policy's meaning is unchanged for every case.
- [ ] `gc` still prunes exactly the Runs it pruned before, and never touches
      `runs` rows or locks.

## Context

- interface: internal/store/journal.go
- interface: internal/cli/gc.go

## Verification

- `output="$(go test -count=1 ./internal/store -run '^TestRetentionScanIsBoundedByCandidates$' -v 2>&1)"; st=$?; printf '%s\n' "$output" | grep -q -- '--- PASS' && [ "$st" -eq 0 ]`
  — expected: exit 0; the retention tests are selected and pass.
- `output="$(go test -count=1 ./internal/store -run 'RetentionScanOutsideWriteTransaction' -v 2>&1)"; st=$?; printf '%s\n' "$output" | grep -q -- '--- PASS' && [ "$st" -eq 0 ]`
  — expected: exit 0; a named test proves the scan left the write transaction,
  rather than the property being asserted in prose.
  — expected: exit 0; the store and the `gc` command stay green.

A whole-package `go test` sweep and `go build ./...` are deliberately absent:
both pass against a tree where no work has happened, so each approves the Task
before it starts. Regression and compilation are the Run-level gate's job.

## References

- `_prd.md` → Core Feature 2; Goal 1; User Story 1.
- `_techspec.md` → System Architecture (the retention path); Build Order 3.
- ADR-0033.
- `references/2026-08-06-event-journal-payload-economics.md`.

## Result

Implementation-ready for Daemon Verification. No Task status was edited; no
commit was made.

### What changed (behavior)

- `terminalRunPruneCandidates` now selects only `runs` rows — the cutoff is a
  SQL predicate (`julianday(r.completed_at) <= julianday(?)`) that bounds
  eligibility to the candidate set, and the old `LEFT JOIN run_events` with
  `COUNT(e.run_id)` aggregate is gone from the eligibility query entirely. The
  authoritative exact boundary check (`parseTime(completed_at).Before(cutoff)`)
  still runs in Go over now-bounded candidates, so the retention policy's
  meaning is preserved exactly.
- `PruneTerminalRuns` runs the eligibility scan on the read connection first
  and only takes the machine-wide write lock — and only when rows are actually
  eligible. A no-op prune no longer opens a write transaction.
- The reported event count is **derived cheaply from the candidate set**
  (requirement 3's option a), via a bounded `COUNT(*) ... WHERE run_id IN (<candidates>)`
  query that runs on the read connection outside any write transaction. The
  Run-start sweep's reported `journal_rows` already came from the actual DELETE
  `RowsAffected`, which is unchanged and never touched a full-table aggregate.
- Schema, `run_events`/`runs` rows, Active Run locks, and VACUUM are untouched.

### Acceptance-criterion evidence

- **Eligibility bounded by candidate set, not the event table** — new
  `TestRetentionScanIsBoundedByCandidates` inspects the eligibility query via
  the package AST and asserts it contains no `run_events` reference and that
  the cutoff is a SQL predicate (on `completed_at` via `julianday`).
- **No aggregate over the event table inside a write transaction** — the
  eligibility query has no `run_events` aggregate; the only `run_events`
  aggregate (the candidate-set event count) lives in
  `countPruneCandidateEvents`, which runs on the read connection. New
  `TestRetentionScanOutsideWriteTransaction` holds the machine-wide advisory
  write lock and asserts a no-op `PruneTerminalRuns` completes immediately
  instead of blocking on a write transaction.
- **Policy meaning unchanged** — existing `TestPruneTerminalRunsDeletesOnlyEligibleJournalRows`,
  `TestPruneTerminalRunsNoOpsWhenCutoffSelectsNothing`,
  `TestRetentionPreservesRunLifecycleRecords`, and `TestPruneTerminalRunsNeverVacuumsAllocatedPages`
  pass focused (excluded from this run only the two declared Verification
  commands that the Daemon owns).
- **`gc` behavior unchanged** — `TestRunGCDryRunListsEligibleRunsAndChangesNothing`,
  `TestRunGCPrunesEligibleJournalsArtifactsAndOrphans`, `TestRunGCSkipsWhenJournalRetentionIsZero`,
  and the implement preflight sweep tests pass. Output shapes are unchanged for
  every case that has rows; for the no-eligible-rows case both the dry-run
  `Journal rows eligible: 0` and the sweep's skip-printing behavior are
  byte-identical to before, so **no output change** applies there.

### Focused checks run (declared Verification not run, Daemon-owned)

- `go build -buildvcs=false ./...` → build OK.
- `go vet -buildvcs=false ./internal/store` → OK.
- `gofmt -l internal/store/journal.go internal/store/journal_test.go` → clean.
- Per-package compile of `./internal/store ./internal/cli` test binaries → OK.
- Existing store retention/lifecycle suites and cli `gc`/preflight suites pass
  as listed above; `TestJournalMeasurementHarness` passes (self-seeded Active
  Run remains ineligible).

### Follow-ups

None in scope. The candidate-set event count is currently a bounded per-candidate
`IN` query; it exists only to keep the `gc` dry-run "Journal rows eligible" line
byte-identical, and could be dropped entirely if that report line were ever
retired.
