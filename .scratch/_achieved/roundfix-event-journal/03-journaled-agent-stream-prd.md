---
title: "PRD: journaled ACP Agent stream"
type: PRD
category: enhancement
state: ready-for-agent
labels:
  - enhancement
  - ready-for-agent
blocked_by:
  - 01-run-event-seam-prd.md
  - 02-run-event-journal-prd.md
---

# PRD: journaled ACP Agent stream

## Problem Statement

Roundfix now receives structured ACP stream updates, but those updates are still primarily delivered to the active writer and the Agent log. If the terminal disconnects or the TUI stops, the structured stream is not durably available as a product event stream.

The user needs the Agent console to behave like a durable Run timeline: visible live when attached, replayable after disconnect, and precise enough to diagnose tool calls, verification output, and destructive-action prevention.

## Solution

Register the Run Event Journal as a critical sink behind the existing Run Event seam, so every Agent Run Event is appended to the journal while the existing output adapters keep working. The Agent runner already publishes Run Events through the sink seam and stays unaware of terminal UI and persistence details; this slice adds the journal adapter, not a new boundary.

This should make Resolve Runs durable without requiring the attach command yet.

## User Stories

1. As a developer, I want ACP Runtime output to be journaled, so that I can inspect Agent behavior after the process exits.
2. As a developer, I want tool start events to be persisted, so that I can see what the Agent attempted.
3. As a developer, I want tool update events to be persisted, so that I can see completed, failed, and in-progress tool outcomes.
4. As a developer, I want tool input and output blocks to be persisted, so that debugging does not require raw protocol logs.
5. As a developer, I want diff and terminal blocks to be persisted, so that the Live Run View can render them distinctly.
6. As a developer, I want Agent message chunks to be persisted, so that final summaries are replayable.
7. As a developer, I want Agent thought chunks to remain supported in the stream model, so that runtimes that expose them do not break the journal.
8. As a developer, I want session status events persisted, so that I can tell whether an ACP session completed, failed, or stopped.
9. As a developer, I want Batch metadata on stream events, so that multi-Batch Resolve Runs are understandable.
10. As a developer, I want Agent log files to keep working, so that existing raw diagnosis workflows do not regress.
11. As a developer, I want journal failures to fail the Run clearly, so that Roundfix does not silently lose the source of truth.
12. As a developer, I want bounded text summaries, so that huge tool output does not make list views unusable.
13. As a daemon, I want one stream sink boundary, so that future subscribers can attach without modifying ACP runner code.
14. As a TUI, I want stream events to arrive in the same structure used by replay, so that live and replay rendering stay consistent.
15. As a test author, I want a fake journal sink, so that Agent stream behavior can be tested without a terminal.
16. As a test author, I want the current no-live output path to keep working, so that CLI output contracts do not regress.

## Implementation Decisions

- Keep ACP Runtime execution under the Agent boundary and reuse the Run Event seam introduced by the Run Event seam PRD. No new publication boundary is created here.
- Persist Run Events exactly as published by the seam: normalized fields for filtering plus the raw ACP session update payload per ADR 0008. Never parse formatted console text.
- Persist both the bounded text summary and the raw payload. The summary exists for quick list rendering; the payload exists for faithful TUI rendering.
- Run ID and Batch number arrive already stamped by the Run Event seam. This slice must not restamp or infer them.
- Keep the existing Agent log output path as a parallel artifact. The journal should not depend on the log file parser.
- Ensure stream update formatting remains backwards-compatible for non-live output where practical.
- Treat journal append failure as a Run failure once the Run has started, because losing the durable event stream violates this feature's product contract.
- Keep live subscriber fanout out of this PRD. The sink can be shaped to support fanout later, but should only persist and write current output in this slice.
- Preserve Stop Request behavior: when a context is canceled, the Agent should stop gracefully and the stop event should be journaled before later daemon actions are skipped.
- Do not add a dependency for this slice. The existing SQLite store and structured stream types are enough.

## Testing Decisions

- Good tests should prove that structured ACP stream updates become Run Events with the expected source, kind, Batch number, tool metadata, summary, and payload.
- Agent tests should cover tool start, tool update, message chunk, plan, status, and raw writer fallback.
- CLI or daemon-level tests should prove a Resolve Run with a fake Agent persists journal events while preserving existing console/log behavior.
- Stop tests should prove canceled Agent execution records a stopped event before returning.
- Negative tests should prove journal append errors surface as Run failures and do not get swallowed.
- Tests should use fake Agent runners and temporary Run Databases rather than launching real ACP runtimes.
- Existing stream formatting tests should remain at the TUI seam, not by asserting private implementation branches.

## Out of Scope

- New attach commands.
- Live subscriber fanout.
- Watch Run event streaming.
- Event pruning or compression.
- Changing ACP Runtime selection or model behavior.
- Replacing the Agent log file.

## Further Notes

This PRD should be implemented after the Run Event Journal exists. It is the first user-visible durability improvement because Resolve Runs are where the most valuable ACP output is produced.
