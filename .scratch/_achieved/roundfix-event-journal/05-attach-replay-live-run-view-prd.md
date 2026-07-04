---
title: "PRD: attach and replay Live Run View"
type: PRD
category: enhancement
state: ready-for-agent
labels:
  - enhancement
  - ready-for-agent
blocked_by:
  - 01-run-event-seam-prd.md
  - 02-run-event-journal-prd.md
  - 03-journaled-agent-stream-prd.md
---

# PRD: attach and replay Live Run View

## Problem Statement

Roundfix has a Live Run View for active command execution, but it is not yet a durable attach surface. A user who closes the terminal, opens another shell, or wants to inspect a completed Run cannot start from the Run Database cursor and replay the structured event stream.

The user wants a cockpit that behaves like an operational console: attach to an Active Run, replay historical events, then follow new events without losing context or duplicating output.

## Solution

Add an attach/replay experience for the Live Run View. The UI should load the Run snapshot, replay Run Events from the journal, render them through the same timeline renderer used for live Agent output, and then follow newly appended events while the Run remains active.

The same renderer should support active follow mode and completed-run inspection. Detaching the UI must not stop the Run. Stopping the Run must remain an explicit Stop Request.

## User Stories

1. As a developer, I want to attach to an Active Run by Run ID, so that I can resume watching a Run from another terminal.
2. As a developer, I want the attach view to replay prior events, so that I do not lose context when reconnecting.
3. As a developer, I want replay and live events to use the same visual language, so that the cockpit is consistent.
4. As a developer, I want the UI to avoid duplicate replay after reconnect, so that event timelines stay readable.
5. As a developer, I want the UI to continue following new events after replay, so that it behaves like a live console.
6. As a developer, I want to inspect a completed Run, so that I can diagnose failures after the fact.
7. As a developer, I want detach to leave the Run active, so that closing the UI does not cancel work.
8. As a developer, I want Stop Request to remain explicit, so that accidental keypresses do not terminate a Run.
9. As a developer, I want the left pane to show Run metadata and Review Issues while replaying events, so that the console is anchored to the current scope.
10. As a developer, I want the right pane to show structured Agent events, so that tool calls and outputs remain readable.
11. As a developer, I want the UI to show when it is replaying backlog versus following live events, so that I understand what I am seeing.
12. As a developer, I want attach to fail clearly when the Run ID does not exist, so that I know whether to list or start a Run.
13. As a developer, I want attach to work when no Agent is currently running, so that completed Runs are still inspectable.
14. As a developer, I want bounded console memory during attach, so that large journals do not exhaust memory.
15. As a daemon, I want attach clients to use cursor reads, so that slow UIs do not block event producers.
16. As a test author, I want attach behavior to be testable with a fake event source, so that terminal rendering is not required for every test.

## Implementation Decisions

- Add a CLI attach surface for Run inspection and Active Run follow mode.
- Reuse the existing Live Run View layout rather than creating a separate screen.
- Use the Run Event Journal as the source of console history. Do not parse Agent log files for replay.
- Implement replay as cursor-based reads with a limit. The UI advances its cursor only after accepting events.
- Implement follow mode through polling the Run Event Journal first, using the journal's change-detection signal (SQLite data version) so idle polls do not read rows. The attach reader uses a read-only database connection with short autocommit reads and never holds a long-lived read transaction. A later optimization may replace polling with daemon push notifications.
- Treat attach as non-mutating. It must not create Runs, fetch Review Source issues, start Agents, commit, push, or resolve Review Source threads.
- Keep Stop Request separate from detach. Detach exits the UI; Stop Request transitions the Run toward stopped behavior.
- Render completed Runs as read-only. If the Run is terminal, the UI should replay history and show terminal state without waiting forever.
- Render missing Run, missing journal, and database errors as CLI errors before starting the TUI.
- Render live and replayed Run Events through one timeline renderer that consumes Run Events only. Live delivery and journal replay are two adapters feeding the same renderer.
- Decode Agent payloads from the raw ACP session update JSON per ADR 0008 and skip unknown event kinds instead of failing, so journals written by newer ACP Runtime versions stay viewable.
- Keep console memory bounded inside the renderer with a ring buffer, and coalesce high-frequency message chunks before rendering so chatty streams do not thrash the terminal.
- Keep the first attach implementation local-first and Run Database-backed. Do not require a daemon socket to attach.

## Testing Decisions

- Good tests should prove attach reads Run metadata, replays events after cursor zero, follows later events, and exits cleanly for terminal Runs.
- CLI tests should cover missing Run ID, unknown Run ID, non-mutating attach, and terminal Run replay.
- TUI tests should cover rendering replayed tool events, state strips, and follow-mode status text.
- Store-backed tests should prove cursor progression avoids duplicate events.
- Cancellation tests should prove detach exits the UI without completing or stopping the Run.
- Polling tests should use injected clocks or fake event sources; no sleeps.
- Race tests are required if follow mode uses goroutines or channels.

## Out of Scope

- Remote daemon sockets.
- Multi-user attach.
- Browser UI.
- Full Run list UI.
- Event filtering by issue or tool call.
- Editing Review Issues from the attach UI.
- Starting new Resolve or Watch Runs from attach mode.

## Further Notes

This PRD intentionally starts with a local Run Database polling model. That is simpler than a daemon transport and still stronger than a purely in-memory active stream because it survives terminal disconnects and process restarts.
