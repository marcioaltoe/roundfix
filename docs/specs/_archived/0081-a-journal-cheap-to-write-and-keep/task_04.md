---
task: task_04
spec: 0081-a-journal-cheap-to-write-and-keep
status: completed
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

The declared Verification names `TestBatchClosesOnCountLingerAndImmediate`, which does not exist yet, so it can
fail before the work. Create it to assert that a batch closes on count, on linger, and before an immediate event. A broad pattern over
this package matches cases that already pass and would approve the Task
before it starts.

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

- `output="$(go test -count=1 ./internal/store -run '^TestBatchClosesOnCountLingerAndImmediate$' -v 2>&1)"; st=$?; printf '%s\n' "$output" | grep -q -- '--- PASS' && [ "$st" -eq 0 ]`
  — expected: exit 0; the batching tests are selected and pass.
- `output="$(go test -count=1 ./internal/store -run 'Batch|AmbiguousCommit' -v 2>&1)"; printf '%s\n' "$output" | grep -cE -- '--- PASS: [^ ]+/' | { read count; [ "$count" -ge 8 ]; }`
  — expected: exit 0; at least the eight rehearsal cases exist as passing
  subtests, so the boundary suite cannot collapse to a happy path.
  — expected: exit 0; every package that reads or writes events stays green.

A whole-package `go test` sweep and `go build ./...` are deliberately absent:
both pass against a tree where no work has happened, so each approves the Task
before it starts. Regression and compilation are the Run-level gate's job.

## References

- `_prd.md` → Core Feature 3; Goal 2; User Stories 2, 5, 6.
- `_techspec.md` → Implementation Design (batching boundary); Build Order 4.
- ADR-0098, ADR-0008, ADR-0030.

## Result

Implementation-ready for Daemon Verification. This is a Roundfix Daemon-assigned
turn: Task status is Daemon-owned and left untouched (the work-target lock
proves no live Agent owns it), declared `## Verification` commands are the
Daemon's to run verbatim, and no commit was made. Evidence below is from
focused implementation checks only.

### What changed (behaviour)

- The append path is now Store-scoped and batched. A `journalWriter` lives on
  the writer `Store` (`store.journal`, created in `Open`), and every sink the
  Store hands out shares that one writer — `(*Store).JournalSink()` returns a
  `runevent.Sink` backed by it, so batch boundaries are global to the process
  rather than per sink (requirement 1).
- `journal.writer` holds a pending ordered batch. `publish` appends in call
  order; a batch closes (commits) on count (`journalBatchSize = 128`), on the
  linger deadline (`journalMaxLinger = 100ms`, a timer armed on the batch's
  first event), before an immediate-durability event, and on explicit flush
  (requirement 2). The limits are internal constants, not configuration or CLI
  surface (requirement 3).
- `commitJournalBatch` opens one write transaction through the existing
  `withWriteTx`, groups the ordered batch per Run, and allocates a contiguous
  cursor range per Run (`MAX(cursor)+1 … +n`) inside the transaction. Events
  are removed from the pending batch only after a successful commit (requirement
  4). On any begin, insert, or commit failure the whole pending batch is kept at
  the head of the queue and the error returns through the caller, so a failed
  append still fails the Run exactly as today (requirement 5).
- An ambiguous commit is reconciled by reading the assigned cursor range back
  (requirement 6): an exact field-and-payload match settles the batch (returns
  nil, no retry), no rows permits exactly one retry that reuses the original
  cursor range, and a partial or different match fails as corruption. `withWriteTx`
  now wraps commit failures in a `writeCommitError` (message contract
  unchanged, `commit <operation>: <cause>`) so the writer can distinguish an
  ambiguous commit from a begin/insert failure, while preserving the
  `TestWriteTxIsTheOnlyWriterTransaction` discipline (all `BeginTx` still inside
  the helper). No payload byte is rewritten, pruned, or compressed — each row is
  inserted verbatim from the producer `json.RawMessage` (requirement 8,
  ADR-0008).
- Terminal settlement wires the batching boundary (requirement 7): the batch is
  flushed before the terminal outcome is published, and the terminal outcome
  (`daemon.outcome`) plus the post-terminal notification receipt are routed
  through the **direct immediate path** (`store.AppendRunEvent` / `withWriteTx`),
  never through a `JournalSink` that may have been closed. `CloseJournal`
  flushes the pending batch and marks the writer closed only after that flush
  commits; a failed Close preserves the batch and remains retryable, and every
  later `Publish` after a successful Close is rejected (`errJournalClosed`).
- `Store.Close` now closes the Store-scoped journal writer (flush + reject)
  alongside the lock file and database.
- Agent-selection transition events
  (`daemon.agent_selection_*`) are treated as immediately durable, because an
  ACP adapter prepares the next (fallback) session only after reading the
  durable fallback notification back from the journal — confirmed by the macro
  ACPX shim contract.

### CLI wiring

- All 12 `store.JournalSink{Store: runStore}` construction sites were replaced
  with `runStore.JournalSink()` (runui.go + cli.go), so every command shares the
  Store-scoped writer.
- `publishRunOutcome` uses `runStore.AppendRunEvent` directly (bypasses the
  closed writer).
- `journalOutcomeNotificationReceipt` uses `runStore.AppendRunEvent` directly
  (post-terminal receipt never enters the closed batch).
- `publishTerminalCompletionWithContext` flushes the journal before publishing
  the terminal outcome so the batch precedes the outcome in cursor order.
- `runUI.Close` flushes the shared journal writer at the Agent-teardown
  boundary, so events published through the sink are durable before the command
  closes its Store.
- `CompleteRun`'s signature is **unchanged**; the terminal outcome is routed
  through the direct immediate path rather than by folding the outcome event
  into `CompleteRun`, satisfying requirement 7 without a wide call-site ripple
  (74 callers remain untouched).

### Acceptance-criterion evidence (focused checks)

- **A batch closes on count, on linger, before an immediate event, and on
  explicit flush** — `TestBatchClosesOnCountLingerAndImmediate` proves all four
  boundaries (count, linger, immediate pre-flush, explicit flush).
- **Publisher order and per-Run cursor contiguity hold under concurrency** —
  `TestBatchPublishPreservesOrderAndContiguousCursors` (single-threaded cursor
  contiguity) and the `concurrent publishers` subtest of
  `TestBatchBeginInsertCommitFailurePreservesBatch` (4 goroutines × 25 events →
  contiguous cursors 1..100) pass, also under `-race`.
- **A failed commit preserves the batch and fails the Run as before** — the
  `insert failure preserves the batch and fails the Run` subtest proves a
  missing-Run flush preserves the whole pending batch and keeps failing until
  drained.
- **Each of the three ambiguous-commit outcomes behaves as specified** —
  `TestBatchAmbiguousCommit` covers exact-match settle (no retry), no-rows
  one-same-cursor retry, partial-match corruption, and different-payload
  corruption.
- **No payload byte differs from what the producer emitted** —
  `TestBatchPreservesPayloadBytes` round-trips a raw payload unchanged; the
  payload-preservation check in `TestAgentConsoleDisplaySinkUsesStatefulSinkForNonTTYAndDetachedLogWriter`
  also passes.
- **Store-scoped shared writer** — the `multiple sinks share one store-scoped
  writer` subtest proves two `JournalSink()` handles publish into one pending
  batch.

### Focused checks run (declared Verification not yet run, Daemon-owned)

- `go test -count=1 ./internal/store -run '^TestBatchClosesOnCountLingerAndImmediate$' -v`
  → `--- PASS` (the declared Verification's first selector).
- `go test -count=1 ./internal/store -run 'Batch|AmbiguousCommit' -v` → 11
  nested `--- PASS` subtests (the declared Verification's second selector, ≥8).
- `go test -count=1 ./internal/store/` → 229 passed.
- `go test -count=1 ./...` → 3918 passed across 27 packages.
- `go test -race -count=1 ./internal/store/` → 229 passed.
- `go vet -buildvcs=false ./internal/store/ ./internal/cli/` → clean.
- `gofmt -l` on all changed files → clean.
- `go build -buildvcs=false ./...` → clean.

### Follow-ups (not this task's slice)

- The batched writer and the `CompleteRun` terminal outcome remain separate
  paths (outcome via direct `AppendRunEvent`). If a future task folds the
  outcome into `CompleteRun`'s transaction for stricter atomicity, that is a
  deliberate follow-up, not a gap against this Task's requirement 7.
- Non-terminal one-shot records that still go through `JournalSink()` (e.g.
  cleanup-warning events) now land at the next flush/linger rather than on the
  immediate call; their callers already tolerate best-effort journaling.
