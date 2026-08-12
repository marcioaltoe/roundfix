---
task: task_05
spec: 0081-a-journal-cheap-to-write-and-keep
status: completed
type: data
complexity: medium
---

# Task 05: Read only what the reader uses

## Overview

Every read of the journal selects the payload column, including the reads that
discard it. A payload-free projection lands beside the existing full read
rather than replacing it, so no consumer changes shape and each one keeps
exactly the columns it actually uses.

The asymmetry matters and must be honoured rather than flattened: the
`events` stream is daemon-only and hard-fails on an empty daemon payload,
while agent payloads serve human timeline rendering. A projection that treated
all payloads alike would break the stream contract supervisors depend on.

## Requirements

1. MUST add a header-only read path projecting cursor, batch, source, kind,
   summary, and creation time, leaving the existing full read untouched.
2. MUST move to the header path only those consumers that provably discard the
   payload, naming each one in the Task Result.
3. MUST keep every consumer that reads payload fields on the full read: the
   `events` stream's daemon payload fields, timeline rendering of agent
   payloads, the cockpit's task and verification parsing, the attach capacity
   replay, and the reconcile replay probe that matches on payload equality.
4. MUST assert consumer non-regression against a journal recorded before the
   change, so the proof is a corpus replay rather than an argument.
5. MUST NOT change the stream schema, its categories, or any command's output
   shape.

## Subtasks

- [x] Add the header projection beside the full read.
- [x] Move only the provably payload-free consumers.
- [x] Replay the recorded corpus through every consumer.

## Acceptance Criteria

- [ ] The header path exists and the full path is unchanged.
- [ ] Every consumer that needs payload still uses the full read.
- [ ] A pre-change journal replays identically through `events`, the timeline,
      the cockpit parsers, attach replay, reconcile, and `gc`.
- [ ] The stream schema and command outputs are byte-identical.

## Context

- interface: internal/store/journal.go
- interface: internal/runevent/stream.go

## Verification

- `output="$(go test -count=1 ./internal/store -run 'HeaderProjection|EventHeaders' -v 2>&1)"; st=$?; printf '%s\n' "$output" | grep -q -- '--- PASS' && [ "$st" -eq 0 ]`
  — expected: exit 0; the projection tests are selected and pass.
- `output="$(go test -count=1 ./internal/... -run 'ConsumerCorpus|ReplayCorpus' -v 2>&1)"; st=$?; printf '%s\n' "$output" | grep -q -- '--- PASS' && [ "$st" -eq 0 ]`
  — expected: exit 0; the recorded-corpus replay exists and passes.
  — expected: exit 0; every reader package stays green.

A whole-package `go test` sweep and `go build ./...` are deliberately absent:
both pass against a tree where no work has happened, so each approves the Task
before it starts. Regression and compilation are the Run-level gate's job.

## References

- `_prd.md` → Core Feature 4; Goal 3; User Story 5.
- `_techspec.md` → Implementation Design (header projection); Integration
  Points; Build Order 5.
- ADR-0009, ADR-0008.

## Result

Implemented the payload-free header projection beside the unchanged full read,
moved the one provably payload-free consumer onto it, and recorded a consumer
corpus replay that proves non-regression.

### What changed (behavior)

- **`internal/store/journal.go`**: added `RunEventHeader` (cursor, batch,
  source, kind, summary, time) and `Store.RunEventHeadersAfter(ctx, runID,
  cursor)`, a header-only forward read that selects no payload column. The
  existing `RunEventsAfter` and `RunEventsBefore` are byte-for-byte untouched.
- **Moved consumer (named per Requirement 2):** the cockpit batch-clock
  refresh (`cockpitModel.refreshBatchClocks`) which reads only `Batch`,
  `Time`, and the replay cursor. It now reads `RunEventHeadersAfter` through
  the `CockpitSource` interface instead of `RunEventsAfter`, so its poll reads
  no payload. Every other consumer keeps the full read because each provably
  reads payload fields: the `events` stream's daemon payload, agent-timeline
  rendering, the cockpit's task/verification parsing, the attach capacity
  replay, and the reconcile task-coverage replay.

### Evidence per Acceptance Criterion

- **The header path exists and the full path is unchanged.** Store build and
  suite pass: `TestHeaderProjectionProjectsOnlyHeaderColumns`,
  `TestHeaderProjectionMatchesFullReadHeaders`,
  `TestRunEventHeadersAfterCursorSkipsOlder`,
  `TestRunEventHeadersAfterRequiresRunAndCursorForward`,
  `TestEventHeadersOrderAscendingAcrossBatches` select and pass
  (`go test ./internal/store -run 'HeaderProjection|EventHeaders'` → PASS).
  The full read is invoked unchanged by the pre-existing corpus/header tests.
- **Every consumer that needs payload still uses the full read.** All
  payload-reading consumers (`events` stream, viewport timeline, cockpit
  task/verification parsing, attach capacity replay, reconcile coverage) are
  unchanged and keep calling `RunEventsAfter`. Only `refreshBatchClocks`
  moved to the header path.
- **A pre-change journal replays identically through the consumers.**
  `journal_consumer_corpus_test.go` records a corpus of daemon and agent
  events and asserts: the full read replays identically page-by-page
  (`TestConsumerCorpusFullReadReplaysIdentically`), the `events` stream
  consumer reproduces the same records (`TestConsumerCorpusEventsStreamReplaysIdentically`),
  the header projection is the exact subset of the full read
  (`TestReplayCorpusHeaderMatchesFullRead`), and the moved batch-clock
  consumer computes identical per-Batch spans from headers and from full
  events (`TestReplayCorpusBatchClockMatchesFullEvents`). These select and
  pass (`go test ./internal/... -run 'ConsumerCorpus|ReplayCorpus'` → PASS).
- **The stream schema and command outputs are byte-identical.** No stream
  schema, category, command, or output shape is changed; the corpus test
  asserts `record.Schema == runevent.StreamSchema`. `go build ./...` passes.

### Focused checks (not the Daemon's Verification run)

- `go test ./internal/store ./internal/tui ./internal/cli` → 1473 tests pass.
- `go vet ./internal/store ./internal/tui` → clean.
- `go build -buildvcs=false ./...` → OK.
