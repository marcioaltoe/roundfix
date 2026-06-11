---
title: "PRD: durable Run Event Journal"
type: PRD
category: enhancement
state: ready-for-agent
labels:
  - enhancement
  - ready-for-agent
blocked_by:
  - 01-run-event-seam-prd.md
---

# PRD: durable Run Event Journal

## Problem Statement

Roundfix can show Agent and verification output while a Run is active, but the Live Run View still depends on the currently attached terminal stream. If the terminal closes, the TUI exits, or a developer wants to inspect a still-running or recently completed Run from another shell, Roundfix has no durable event cursor to replay what happened.

This makes long Resolve Runs and Watch Runs harder to trust. The user needs the Run Database to become the durable source of truth for Run events, not just Run lifecycle metadata.

## Solution

Add a durable Run Event Journal to the Run Database. Every meaningful runtime event should be appended with a monotonic cursor that can be queried later. The journal should preserve structured ACP stream updates, Run state messages, daemon decisions, verification milestones, Batch boundaries, Stop Request events, and push/source-resolution events.

The user should be able to rely on the Run Database to answer: what happened, in what order, for which Run, Batch, issue, tool call, and ACP Runtime.

## User Stories

1. As a developer, I want Roundfix to persist Agent output as Run Events, so that I can inspect a Run after the terminal closes.
2. As a developer, I want every Run Event to have a cursor, so that a UI can replay events without duplicates.
3. As a developer, I want Run Events to preserve ACP tool call identity, so that I can correlate a tool start, update, and result.
4. As a developer, I want Run Events to preserve Batch numbers, so that I can understand which Batch produced each event.
5. As a developer, I want Run Events to preserve the event source, so that I can distinguish Agent, daemon, verification, git, Review Source, and TUI events.
6. As a developer, I want Run Events to preserve structured payloads, so that the TUI can render tool input, output, diffs, terminal blocks, and plain text differently.
7. As a developer, I want Run Events to be queryable after completion, so that failure analysis does not depend on log files alone.
8. As a developer, I want the journal to be stored in the central Run Database, so that it follows the existing Roundfix Home model.
9. As a developer, I want the journal to avoid unbounded in-memory buffers, so that long Watch Runs remain stable.
10. As a developer, I want journal writes to be ordered per Run, so that replay is deterministic.
11. As a developer, I want journal reads to support a cursor and limit, so that a UI can page through large histories.
12. As a developer, I want journal rows to include timestamps, so that I can compare runtime delays and daemon decisions.
13. As a developer, I want the journal to survive process restart, so that future daemon attach flows can reconnect safely.
14. As a daemon, I want to append lifecycle events through a narrow API, so that command logic does not write database rows directly.
15. As a TUI, I want a single event query API, so that replay behavior is independent of how events were produced.
16. As a test author, I want deterministic event insertion, so that cursor ordering and replay can be proven without sleeps.

## Implementation Decisions

- Add a Run Event Journal to the existing Run Database instead of adding a new file format or external broker.
- Store one row per event with a per-Run cursor that is monotonic and never reused.
- Use the existing Run identity as the parent relation. Events must not exist without a valid Run.
- Preserve enough normalized fields for efficient rendering and filtering: event kind, event source, Batch number, Review Issue reference when known, ACP tool ID, ACP tool state, text summary, timestamp, and structured payload JSON.
- Store the raw producer JSON as the durable payload per ADR 0008. For Agent events that is the raw ACP session update exactly as received. The schema must not force every future ACP content type into a dedicated column, and readers must skip unknown event kinds instead of failing.
- The journal consumes Run Events through the sink interface defined by the Run Event seam PRD and registers as a critical sink: an append failure after Run start fails the Run.
- Add store-level APIs for appending one event, appending a small ordered batch of events when useful, and listing events after a cursor with a caller-supplied limit.
- Enforce single-writer connection discipline: one writer connection with immediate transactions, WAL journal mode enabled at database creation before any reader connects, and busy timeouts applied on every connection.
- Allocate the per-Run cursor inside the insert transaction from the current per-Run maximum, which is race-free under the single-writer rule.
- Store payloads as plain TEXT JSON, not a binary JSON encoding, so payloads round-trip byte-exact and stay inspectable with standard tooling.
- Expose a cheap change-detection signal for pollers based on the SQLite data version, so follow mode in the attach PRD can detect new events without reading rows.
- Treat the cursor as an opaque integer for callers. Callers should only compare it by ordering and use it as the next replay position.
- Keep journal insertion inside database transactions when the event must be committed atomically with a Run state transition.
- Do not replace existing Agent log files in this PRD. Logs remain useful for raw diagnosis; the journal becomes the structured product surface.
- Do not introduce network services, sockets, or daemon IPC in this PRD. This PRD only creates durable state and store APIs.
- Do not prune journal entries in the first implementation. Retention can be added later after the attach/replay behavior is stable.

## Testing Decisions

- Good tests should prove externally visible store behavior: migration, append, cursor ordering, pagination, payload round-trip, and invalid Run handling.
- Store tests should create a temporary Roundfix Home and real SQLite Run Database.
- Migration tests should prove existing Run Database creation still works and the new journal schema is present.
- Append tests should prove cursors increase monotonically per Run.
- Replay tests should prove listing after cursor returns only newer events and respects limits.
- Payload tests should prove structured ACP-like JSON survives round-trip without losing block kind or tool metadata.
- Negative tests should prove appending to a missing Run fails with a clear error.
- Reader tests should prove a second read-only connection can page events while the writer connection appends.
- Concurrency-sensitive tests should avoid sleeps and use explicit goroutine completion. If concurrent append support is added, race tests must cover it.

## Out of Scope

- Live TUI attach or replay commands.
- ACP Runtime execution changes.
- Daemon subscriber fanout.
- Event retention, pruning, compression, or export.
- Cross-machine synchronization.
- Review Source thread resolution changes.

## Further Notes

This PRD is the durable foundation for the remaining event-stream work. It lands after the Run Event seam PRD, which defines the Run Event type and sink interface this journal persists. Later slices should consume the Run Event Journal instead of inventing their own persistence.
