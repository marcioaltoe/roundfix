---
task: task_04
spec: 0081-a-journal-cheap-to-write-and-keep
status: pending
type: data
complexity: high
---

# Task 04: Amortize appends across a batch

## Overview

Today every agent stdout line is its own immediate transaction, existence
probe, insert, commit, and fsync — 1.6 million of them across fourteen days,
42,000 for a single Run. `AppendRunEvents` already exists and is called only
with one-element slices.

Durability moves from per-event to per-batch and no further: `synchronous`
stays at its default, so a crash can lose at most the batch in flight, and the
events a post-mortem needs most never sit in one.

## Requirements

1. MUST make the batch Store-scoped: every sink the Store hands out shares one
   writer, so boundaries are global to the process rather than per sink.
2. MUST close a batch at its event count, at its linger deadline, before an
   immediate-durability event, and on explicit flush for error paths, Agent
   teardown, terminal settlement, and process shutdown.
3. MUST keep the limits internal constants — not configuration, not CLI
   surface.
4. MUST preserve publisher order and allocate a contiguous cursor range per
   Run inside the write transaction, removing events from the pending batch
   only after a successful commit.
5. MUST preserve the whole pending batch on begin, insert, or commit failure
   and return through the existing critical-sink error path, so a failed
   append still fails the Run exactly as today.
6. MUST reconcile an ambiguous commit by reading the assigned cursor range: an
   exact field-and-payload match settles the batch, no rows permits one retry
   holding the same cursors, and a partial or different match fails as
   corruption. A retry always reuses the original cursor range, since freshly
   allocated cursors would duplicate raw events after an ambiguous success.
7. MUST route the terminal outcome through the direct path that bypasses the
   closed writer, and MUST reject any publish after a successful close.
8. MUST NOT rewrite, prune, or compress any payload — ADR-0008 binds, and this
   Task changes only when bytes reach the database, never what they contain.

## Subtasks

- [ ] Make the writer Store-scoped with its boundaries.
- [ ] Implement ordering, contiguous cursors, and failure preservation.
- [ ] Implement ambiguous-commit reconciliation and the terminal path.

## Acceptance Criteria

- [ ] A batch closes on count, on linger, before an immediate event, and on
      explicit flush.
- [ ] Publisher order and per-Run cursor contiguity hold under concurrency.
- [ ] A failed commit preserves the batch and fails the Run as before.
- [ ] Each of the three ambiguous-commit outcomes behaves as specified.
- [ ] No payload byte differs from what the producer emitted.

## Rehearsal Cases

- Case: a batch reaches its event count; Observation: it commits and the
  pending batch empties only after commit.
- Case: a quiet publisher reaches its linger deadline; Observation: the batch
  commits without waiting for more events.
- Case: an immediate-durability event arrives; Observation: the pending batch
  flushes before it, and the event is durable on its own.
- Case: commit fails outright; Observation: the batch survives intact and the
  error reaches the critical-sink path.
- Case: commit returns ambiguous and the rows are present and identical;
  Observation: the batch settles without a retry.
- Case: commit returns ambiguous and no rows exist; Observation: exactly one
  retry with the same cursors.
- Case: commit returns ambiguous and rows partially match; Observation: it
  fails as corruption rather than retrying.
- Case: publish after a successful close; Observation: rejected.

## Context

- interface: internal/store/journal.go
- interface: internal/runevent/event.go
- interface: internal/daemon/engine.go

## Verification

- `go build -buildvcs=false ./...` — expected: exit 0.
- `output="$(go test -count=1 ./internal/store -run 'Batch|Journal|AmbiguousCommit' -v 2>&1)"; st=$?; printf '%s\n' "$output" | grep -q -- '--- PASS' && [ "$st" -eq 0 ]`
  — expected: exit 0; the batching tests are selected and pass.
- `output="$(go test -count=1 ./internal/store -run 'Batch|AmbiguousCommit' -v 2>&1)"; printf '%s\n' "$output" | grep -cE -- '--- PASS: [^ ]+/' | { read count; [ "$count" -ge 8 ]; }`
  — expected: exit 0; at least the eight rehearsal cases exist as passing
  subtests, so the boundary suite cannot collapse to a happy path.
- `go test -count=1 ./internal/store/... ./internal/runevent/... ./internal/daemon/... ./internal/tui/...`
  — expected: exit 0; every package that reads or writes events stays green.

## References

- `_prd.md` → Core Feature 3; Goal 2; User Stories 2, 5, 6.
- `_techspec.md` → Implementation Design (batching boundary); Build Order 4.
- ADR-0098, ADR-0008, ADR-0030.
