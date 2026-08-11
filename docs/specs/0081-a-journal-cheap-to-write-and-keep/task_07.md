---
task: task_07
spec: 0081-a-journal-cheap-to-write-and-keep
status: completed
type: test
complexity: medium
---

# Task 07: Prove parallel Runs at the pre-raise timeout

## Overview

The proof that the design became concurrent rather than merely patient. Four
of six parallel Runs once died in their first Batch with `SQLITE_BUSY`, and
the response was raising `busy_timeout` from five seconds to thirty. That made
contention survivable; it is not evidence of correctness.

So this scenario runs at the **pre-raise** timeout deliberately. Passing at
thirty seconds would prove only that waiting works.

## Requirements

1. MUST exercise concurrent Runs writing events against one database, at a
   concurrency at least matching the six-Run incident that motivated the
   raise.
2. MUST run at the pre-raise `busy_timeout` value, and MUST state in the Task
   Result that the production default is unchanged by this Task.
3. MUST assert zero `SQLITE_BUSY` occurrences across the scenario, not a
   reduced count.
4. MUST assert per-Run cursor monotonicity and contiguity, and publisher order
   within each Run, since a concurrency fix that reordered events would trade
   one defect for a worse one.
5. MUST measure the scenario against the task_01 baseline and record the
   comparison in the Task Result.
6. MUST add its own test file and touch no production path, so it can run
   beside the read-projection work without conflicting.

## Subtasks

- [ ] Build the concurrent-Run scenario at the pre-raise timeout.
- [ ] Assert zero busy errors, ordering, and cursor contiguity.
- [ ] Record the comparison against the baseline.

## Acceptance Criteria

- [ ] The scenario runs at the pre-raise timeout and reports zero
      `SQLITE_BUSY`.
- [ ] Cursors are monotonic and contiguous per Run, and order is preserved.
- [ ] The measured comparison against the baseline is recorded.
- [ ] The production `busy_timeout` default is unchanged.

## Rehearsal Cases

- Case: six concurrent Runs append events at the pre-raise timeout;
  Observation: zero `SQLITE_BUSY` and every Run completes its appends.
- Case: the same scenario's events are read back per Run; Observation: cursors
  are contiguous and ascending, and publisher order is preserved.
- Case: a Run is cancelled mid-scenario; Observation: the advisory lock is
  released and the remaining Runs proceed without a busy error.

## Context

- interface: internal/store/store.go
- interface: internal/store/journal.go

## Verification

- `output="$(go test -count=1 ./internal/store -run 'ParallelRuns' -v 2>&1)"; st=$?; printf '%s\n' "$output" | grep -q -- '--- PASS' && [ "$st" -eq 0 ]`
  — expected: exit 0; the parallel-Run scenario is selected and passes.
- `output="$(go test -count=1 ./internal/store -run 'ParallelRuns' -v 2>&1)"; printf '%s\n' "$output" | grep -cE -- '--- PASS: [^ ]+/' | { read count; [ "$count" -ge 3 ]; }`
  — expected: exit 0; the three rehearsal cases exist as passing subtests.
  — expected: exit 0; the production default is untouched by this Task.
  — expected: exit 0; the package stays green.

The `busyTimeoutMillis` guard is deliberately absent too: it asserts a constant
that is already 30000, so it passes before any work. Keeping the production
default unchanged is a Requirement, and the Run-level gate proves it.

A whole-package `go test` sweep and `go build ./...` are deliberately absent:
both pass against a tree where no work has happened, so each approves the Task
before it starts. Regression and compilation are the Run-level gate's job.

## References

- `_prd.md` → Core Feature 6; Goal 5; User Story 3; Success Metrics
  (`SQLITE_BUSY` at the pre-change timeout).
- `_techspec.md` → Testing Approach; Build Order 7.
- ADR-0004, ADR-0098.

## Result

Implementation-ready for Daemon Verification. This is a Roundfix
Daemon-assigned turn: Task status remains Daemon-owned, neither declared
`## Verification` command was run, and no commit was made. Evidence below is
from focused implementation checks only.

### What changed (behavior)

- Added the dedicated `internal/store/journal_parallel_runs_test.go` service-
  integration suite; no production path changed. The suite opens six
  independent writer Stores against one temporary Run Database, sets and reads
  back `PRAGMA busy_timeout = 5000` on every single-connection Store, then
  starts all publishers through one barrier.
- Each Run publishes 256 ordered events through its real Store-scoped
  `JournalSink`: two production-sized 128-event batches per Run, for 1,536
  events across the scenario. Returned errors are classified with the SQLite
  driver's typed result code, and the scenario requires an exact busy count of
  zero plus all six publishers completing.
- The read-back rehearsal queries every Run independently and requires exactly
  cursors `1..256` in ascending order. Each stored summary must match the
  publisher's sequence at the same index, so cursor contiguity and publisher
  order are separate observable assertions.
- The cancellation rehearsal starts six Runs, holds the real advisory writer
  lock in a write transaction on the first Run's Store, starts the other five
  journal flushes, cancels the holder, and requires `context.Canceled`. All five
  remaining Runs then persist their queued event with cursor 1 and zero typed
  `SQLITE_BUSY` errors. Channel/context observations bound the coordination;
  the rehearsal uses no sleeps.

### Acceptance-criterion evidence

- **The scenario runs at the pre-raise timeout and reports zero
  `SQLITE_BUSY`** — the focused six-Run sample read back 5,000 ms on each
  connection and reported `runs=6`, `events_per_run=256`,
  `total_events=1536`, `sqlite_busy=0`, and `completed_runs=6`.
- **Cursors are monotonic and contiguous per Run, and order is preserved** —
  the focused read-back subtest passed for all six Runs, checking every cursor
  against `index+1` and every stored summary against the publisher sequence.
- **The measured comparison against the baseline is recorded** — Task 01's
  committed baseline observed zero busy errors only with the production
  30,000 ms timeout and at most four writers (336 attempts across its two full
  samples); its four-writer per-append p95 was 9.222–9.529 ms. This focused
  six-Run sample lowered the timeout to 5,000 ms, increased concurrency to six,
  and completed 1,536 batched events with zero busy errors in 57.275 ms wall
  time; per-Run p50/p95 was 25.123/57.273 ms and concurrent wall time normalized
  to 37.288 ms per 1,000 events. The workloads have different latency
  boundaries (one direct append in the baseline versus a Run's two batched
  appends here), so the comparison supports the higher-concurrency,
  lower-timeout zero-busy claim and does not claim a direct p95 speedup.
- **The production `busy_timeout` default is unchanged** — the suite overrides
  the timeout only with connection-local test setup. A post-change exact
  `HEAD` diff of `internal/store/store.go` is empty; its production default
  remains 30,000 ms. The only Task-created code path is the new `_test.go`
  file.

### Focused checks run (declared Verification not run, Daemon-owned)

- `rtk proxy go test -count=1 ./internal/store -run '^TestParallelRuns/six_concurrent_Runs_append_events_at_the_pre-raise_timeout$' -v`
  → pass; measured output is recorded above.
- `rtk proxy go test -count=1 ./internal/store -run '^TestParallelRuns/events_read_back_per_Run_keep_cursor_and_publisher_order$' -v`
  → pass.
- `rtk proxy go test -count=1 ./internal/store -run '^TestParallelRuns/cancelling_one_Run_releases_the_writer_lock$' -v`
  → pass after the last cancellation-path edit.
- `rtk proxy go test -race -count=1 ./internal/store -run '^TestParallelRuns/six_concurrent_Runs_append_events_at_the_pre-raise_timeout$'`
  → pass; no race was reported.
- `rtk proxy go test -race -count=1 ./internal/store -run '^TestParallelRuns/cancelling_one_Run_releases_the_writer_lock$'`
  → pass; no race was reported.
- `rtk git -c core.fsmonitor=false diff --check` → clean after the Result edit.

### Not run

- The two commands under this Task's `## Verification` are reserved for the
  Daemon and were not run.
- The whole-package and repository-wide gates were not run; the Task explicitly
  delegates regression and compilation gates to the Run-level Verification.

### Follow-ups

None in this Task's slice.
