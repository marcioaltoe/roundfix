---
task: task_07
spec: 0081-a-journal-cheap-to-write-and-keep
status: pending
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
- `grep -q 'busyTimeoutMillis = 30000' internal/store/store.go`
  — expected: exit 0; the production default is untouched by this Task.
- `go test -count=1 ./internal/store/...`
  — expected: exit 0; the package stays green.

## References

- `_prd.md` → Core Feature 6; Goal 5; User Story 3; Success Metrics
  (`SQLITE_BUSY` at the pre-change timeout).
- `_techspec.md` → Testing Approach; Build Order 7.
- ADR-0004, ADR-0098.
