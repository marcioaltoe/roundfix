---
task: task_02
spec: 0081-a-journal-cheap-to-write-and-keep
status: pending
type: data
complexity: high
---

# Task 02: Give the writer one transaction discipline

## Overview

Every write transaction enters through one Store helper, and Roundfix
processes serialize before SQLite ever sees a writer. Today each operational
command opens its own writer connection with `_txlock=immediate`, so six
parallel Runs on one machine become six processes racing for the same lock —
which is how four of six Runs died in their first Batch with `SQLITE_BUSY` and
why the timeout was raised from five seconds to thirty.

The raise made contention survivable. It did not make the design concurrent.
This Task does.

## Requirements

1. MUST route every writer transaction — cursor allocation, inserts, state
   changes, and commit — through one Store helper, so no path opens an ad-hoc
   immediate transaction.
2. MUST serialize Roundfix processes with a machine-wide advisory lock before
   SQLite sees a writer, and MUST release it on every exit path including
   error and cancellation.
3. MUST keep one writer Store per operational process with its single
   connection, leaving read-only Store values separate and unaffected.
4. MUST return context cancellation and advisory-lock failures to the caller
   rather than swallowing them.
5. MUST keep `busy_timeout` in place as defence against an unknown
   non-Roundfix writer, while no Roundfix path relies on it for concurrency.
6. MUST keep cursor allocation inside the locked transaction so concurrent
   Runs preserve input order and monotonic cursors.
7. MUST NOT change the schema, the stream contract, or any command's output.

## Subtasks

- [ ] Introduce the write-transaction helper and route every writer through it.
- [ ] Add the machine-wide advisory lock with its release paths.
- [ ] Prove ordering and monotonic cursors under concurrency.

## Acceptance Criteria

- [ ] No writer path opens a transaction outside the helper.
- [ ] The advisory lock is released on success, error, and cancellation.
- [ ] Concurrent writers produce monotonic, contiguous cursors per Run.
- [ ] Read-only Store behaviour is unchanged.

## Context

- interface: internal/store/store.go
- interface: internal/store/journal.go

## Verification

- `go build -buildvcs=false ./...` — expected: exit 0.
- `output="$(go test -count=1 ./internal/store -run 'WriteTx|AdvisoryLock|Concurrent' -v 2>&1)"; st=$?; printf '%s\n' "$output" | grep -q -- '--- PASS' && [ "$st" -eq 0 ]`
  — expected: exit 0; the discipline tests are selected and pass.
- `output="$(grep -c 'BeginTx' internal/store/*.go)"; printf '%s' "$output" | grep -q . && grep -q 'withWriteTx' internal/store/store.go`
  — expected: exit 0; the helper exists on the real surface.
- `go test -count=1 ./internal/store/... ./internal/cli/...`
  — expected: exit 0; the store and the commands that open it stay green.

## References

- `_prd.md` → Goal 5; User Story 3; Core Feature 6.
- `_techspec.md` → Implementation Design (`withWriteTx`); Build Order 2.
- ADR-0004.
- `docs/findings/_archived/` → the 2026-08-06 run-lifecycle rollup, whose
  members record the six-parallel-Runs incident.
