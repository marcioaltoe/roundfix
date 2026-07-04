---
title: "PRD: interactive Live Run View"
type: PRD
category: enhancement
state: ready-for-agent
labels:
  - enhancement
  - ready-for-agent
blocked_by: []
---

# PRD: interactive Live Run View

## Problem Statement

The Live Run View renders but does not respond. The live cockpit runs Bubble Tea with input disabled, so no key is processed; the timeline shows only the tail of a 300-line ring, so anything that scrolled past is unreachable even though the full history sits in the Run Event Journal; and the sidebar lists Review Issues with no way to select or inspect one. Attach in a TTY prints a static frame and streams text — it never opens the cockpit at all.

The user expects an operational console: scroll back through everything a Run did, jump between Review Issues and read them, and always know whether the view is replaying backlog, following live events, or frozen where they scrolled to.

The codebase also carries a latent hazard this feature must resolve: in the live cockpit Ctrl-C is a Stop Request, while in attach Ctrl-C is detach. One key, opposite consequences.

## Solution

Make the Live Run View one interactive cockpit shared by attach (in a TTY) and the live resolve and watch commands. The cockpit consumes Run Events exclusively from the Run Event Journal — cursor-paged replay plus data-version follow polling — even in the process that is producing the events, per ADR 0009. Live and replayed rendering become identical by construction, and scrollback is the same mechanism as the tail: one bounded sliding window of events paged by cursor in both directions.

The cockpit gains a focus model (Issues pane and Timeline pane), line-based scrolling backed by event paging, a Follow Mode state machine with an explicit backlog/following/scrolled indicator, and a read-only Review Issue detail pane. Non-TTY output is untouched on every command.

## User Stories

1. As a developer, I want to scroll back through the entire Run timeline, so that diagnosis is not limited to the last screenful.
2. As a developer, I want the timeline to follow new events automatically when I am at the bottom, so that the cockpit behaves like a live console.
3. As a developer, I want the view to freeze while I am scrolled up and show how many new events arrived below, so that reading history and watching progress do not fight each other.
4. As a developer, I want one key to jump back to the tail and resume following, so that catching up is instant.
5. As a developer, I want to always see whether the cockpit is replaying backlog, following live, scrolled, or showing a terminal Run, so that I trust what I am looking at.
6. As a developer, I want to move a selection through the Review Issues sidebar, so that long issue lists are navigable.
7. As a developer, I want to open the selected Review Issue and read its artifact (title, severity, status, file:line, Source Reference, body), so that triage context is one keypress away.
8. As a developer, I want attach in a terminal to open the interactive cockpit, including for terminal Runs, so that I can browse a finished Run's history instead of receiving a single static dump.
9. As a developer, I want attach to keep its plain-text behavior when output is not a terminal, so that pipes and CI contracts do not regress.
10. As a developer, I want detach in attach to be a harmless keypress (q or Ctrl-C), so that closing the cockpit never stops a Run.
11. As a developer, I want Stop Request to stay exactly where it is today — Ctrl-C in the owning resolve/watch process, or the stop command — and to not exist as an attach keybinding, so that no accidental keypress can terminate a Run.
12. As a developer, I want the live resolve and watch cockpit to render from the journal, so that what I see live is exactly what attach and replay will show.
13. As a developer, I want one cockpit for an entire Watch Run, so that the view stops being torn down and rebuilt per Batch.
14. As a developer, I want scrolling to never slow down or affect the Run, so that inspection is always safe.
15. As a developer, I want console memory to stay bounded no matter how long the Run is, so that long Watch Runs remain stable.
16. As a test author, I want the window, follow states, and key handling provable with fake event sources and no real terminal, so that the cockpit's behavior is tested where it lives.

## Implementation Decisions

- One interactive cockpit component is shared by attach (TTY) and the live resolve/watch commands. Non-TTY output paths are untouched everywhere.
- The cockpit reads the Run Event Journal only — never the live sink — per ADR 0009: cursor-paged replay, then follow polling on the data-version signal over a read-only connection. The TUI's event-buffer sink delivery is retired; the best-effort sink remains for non-TTY text output.
- The store gains a backward read: the N events immediately before a cursor, returned ascending, with the same read-only short-autocommit discipline as the forward read.
- The timeline holds a bounded sliding window of roughly 500 events (internal constant, not configuration), rendered to lines. Scrolling moves the line viewport; reaching an edge pages the window by events and evicts from the other end. The whole journal is reachable by scrolling; memory stays bounded by the window.
- Follow Mode state machine: REPLAYING while loading backlog; FOLLOWING when the viewport is at the bottom (tail auto-advances); SCROLLED when the user scrolled up — the viewport freezes and the status bar counts new events below with a hint to press End; terminal Runs show the terminal state, read-only, with full scrollback. End/G jumps to the tail and resumes FOLLOWING; manually scrolling to the bottom also resumes. Scrolling never affects the Run: the reader simply stops paging forward.
- Focus model: Tab switches between the Issues pane and the Timeline pane; the focused pane is visibly marked. Issues pane: j/k or arrows move the selection, Enter opens the Review Issue detail pane, Esc closes it. Timeline pane: j/k/arrows scroll lines, PgUp/PgDn page, Home goes to the window top, End/G to the tail.
- The Review Issue detail pane renders the issue artifact read-only (title, severity, status, file:line, Source Reference, markdown body) in place of the timeline. Missing or cleaned artifacts degrade to "artifact not available" without failing. No editing, ever, from the cockpit.
- Timeline filtering by issue stays out: Agent-source events carry no per-issue reference, so filtering would render an empty timeline. The future path is per-issue stamping at event production, never parsing rendered text.
- Key ownership is mode-dependent and shown in a contextual footer. Owning cockpit (resolve/watch): Ctrl-C is Stop Request, unchanged; there is no detach — q does nothing. Attach cockpit: q and Ctrl-C are detach; Stop Request does not exist as a key, because attach is non-mutating by recorded contract — stopping is the explicit stop command.
- Attach in a TTY opens the cockpit in the alternate screen under the existing TTY gate (ROUNDFIX_TUI and terminal detection), including for terminal Runs, which open as a read-only browser exited with q. This consciously amends the earlier "terminal Run exits cleanly" criterion: in a TTY, exiting cleanly now means exiting on q. ROUNDFIX_TUI=never forces today's plain-text attach.
- The live watch command runs one cockpit for the whole Run instead of constructing a view per cycle, which the journal-only model makes natural.

## Testing Decisions

- Good tests prove viewport behavior through the window seam with fake event sources: what is visible, in which follow state, after which sequence of events and keys — never private helper calls.
- Backward-paging store tests prove the window walks the full journal both ways without duplicates or gaps at page boundaries, against a real temporary Run Database.
- Follow-state tests prove FOLLOWING advances the tail, SCROLLED freezes it and counts arrivals, End resumes, and terminal Runs never enter FOLLOWING.
- Key-handling tests drive the Bubble Tea model with synthetic key messages — no real terminal — covering focus switching, issue selection, detail open/close, and detach keys per mode.
- Attach contract tests prove TTY mode opens the cockpit while non-TTY output stays byte-compatible with today's static-plus-streaming text.
- Live-mode tests prove the resolve/watch cockpit renders journal content while the engine writes, that Ctrl-C still produces a Stop Request, and that retiring the sink delivery does not change non-TTY output.
- Detail-pane tests cover a present artifact and a missing one.
- Race tests are required: follow polling, the owning process writing while its own cockpit reads, and detach during active polling.

## Out of Scope

- Timeline filtering by Review Issue or tool call (requires per-issue event stamping at production).
- Editing Review Issues from the cockpit.
- Stop Request as an attach keybinding, or detach in the owning cockpit (run-continues-headless).
- Mouse support, search within the timeline, and configurable keybindings.
- Rich diff rendering, grouped tool entries, and durations (reader-side fidelity work deferred from the event-journal sequence).
- A Run list / picker UI; attach still takes a Run ID.
- Remote daemon sockets or multi-client coordination.

## Further Notes

This PRD builds directly on the event-journal sequence: the journal with per-Run cursors (PRD 02), journaled Agent and daemon events (PRDs 03 and 06), and the attach replay/follow foundation (PRD 05). ADR 0009 records the journal-only consumption decision; CONTEXT.md defines Follow Mode. The grilled decisions here are closed — implement, don't redesign.
