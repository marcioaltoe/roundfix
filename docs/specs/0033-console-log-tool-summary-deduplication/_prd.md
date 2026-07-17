---
spec: 0033-console-log-tool-summary-deduplication
status: active
created: 2026-07-16
surfaces: [backend, cli]
---

# Console Log tool-summary deduplication

A Detached Run can write the same compact read or edit summary multiple times when one ACP Runtime repeats identical content across a tool call's lifecycle events. The Run Event Journal is correctly lossless, but the caller-facing Console Log becomes noisy and can make one edit look like many edits. Console rendering must suppress only proven same-tool duplicates while preserving the journal and every distinct operation.

## Goals

- Render identical compact summaries from one tool call once in the Console Log.
- Preserve repeated summaries from distinct tool calls, even when their text is identical.
- Keep every Run Event, raw payload, cursor, and replay record unchanged.
- Preserve current output for non-tool events and for tool lifecycle updates whose visible content changes.

## User Stories

1. As a developer following a Detached Run, I want one visible summary for one repeated tool lifecycle update, so that the Console Log reflects the number of actual operations.
2. As a developer auditing a Run, I want the Run Event Journal to retain every producer event, so that display cleanup never removes forensic evidence.
3. As a developer whose Agent edits the same file in separate tool calls, I want each operation shown, so that deduplication does not hide real repeated work.
4. As a maintainer using non-interactive command output, I want the same bounded Console Log behavior as detached execution, so that the shared plain-text path stays predictable.

## Core Features

1. Plain-text Agent console rendering remembers the last rendered tool summary for each non-empty tool call identifier while that call is active.
2. A tool event is suppressed only when its rendered bytes exactly match the previously rendered bytes for the same tool call identifier.
3. A changed summary for the same tool call is rendered. An identical summary from a different tool call is rendered.
4. Terminal tool states release their remembered entry after the terminal event is processed, keeping renderer state bounded to active calls.
5. Events without a tool call identifier, non-tool events, undecodable future events, and existing Agent-console suppression behavior retain their current contracts.
6. Deduplication occurs after the Run Event Journal fanout boundary. Journal writes, raw payloads, event kinds, summaries, tool metadata, cursors, Attach, and the Live Run View remain unchanged.
7. The renderer never compares or emits raw tool bodies to make the decision; it uses only the existing bounded console text.

## User Experience

Following a Detached Run still shows compact lines such as one file read or edit summary as soon as the first lifecycle event provides it. If the ACP Runtime repeats the same summary when the tool state changes, the Console Log does not print a second identical line. If the visible summary changes, the new line appears. Separate tool calls remain separate even when they operate on the same path with the same line counts.

## Non-Goals / Out of Scope

- Removing, rewriting, or coalescing Run Events in the Run Event Journal.
- Deduplicating the Live Run View or Attach replay.
- Deduplicating by text across different tool call identifiers.
- Hiding repeated Agent messages, thoughts, plans, status events, or raw output.
- Buffering tool output until completion or delaying the first visible summary.
- Changing compact read/edit wording, raw-payload retention, or Agent log policy.

## Success Metrics

- A started and updated event with the same tool call identifier and identical compact edit text produce one Console Log line and two byte-identical journal payloads.
- Two tool call identifiers with identical compact edit text produce two Console Log lines.
- Two different visible summaries from one tool call both appear in order.
- Tool terminal states remove active renderer state without changing their journal records.
- Existing non-tool console, Agent suppression, Attach, and Live Run View tests pass without changed expectations.
- A dogfood replay of the reported duplicate lifecycle pair prints the edit summary once while retaining both journal cursors.

## Decisions

- Deduplication is keyed by tool call identifier and exact rendered bytes; text-only or path-only deduplication is forbidden.
- The first summary is rendered immediately and later exact duplicates are suppressed; terminal buffering is unnecessary.
- ADR-0008 remains authoritative for lossless raw payloads, ADR-0009 keeps the Live Run View on journal replay, and ADR-0030 keeps the Console Log as the unconditional detached record.

## Open Questions

None.
