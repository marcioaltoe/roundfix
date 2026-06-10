---
title: "Interactive cockpit on attach: focus, keys, issue detail"
type: HITL
category: enhancement
state: ready-for-agent
labels:
  - enhancement
  - ready-for-agent
user_stories:
  - 5
  - 6
  - 7
  - 8
  - 9
  - 10
  - 11
  - 16
blocked_by:
  - 02-timeline-window-follow-state.md
---

# Interactive cockpit on attach: focus, keys, issue detail

## Parent

.scratch/roundfix-interactive-cockpit/01-interactive-live-run-view-prd.md

## What to build

The interactive shell, shipped first on attach where it risks nothing: a Bubble Tea program with input enabled that
hosts the window/Follow Mode engine. Tab switches focus between the Issues pane and the Timeline pane (focused pane
visibly marked); issues navigate with j/k/arrows, Enter opens the read-only Review Issue detail pane (title,
severity, status, file:line, Source Reference, markdown body; missing artifact degrades to "artifact not
available"), Esc closes it; the timeline scrolls with j/k/arrows/PgUp/PgDn/Home/End. Contextual footer shows attach
keys; q and Ctrl-C detach — the Run keeps going, and no stop keybinding exists (attach is non-mutating). Attach in
a TTY opens the cockpit in the alternate screen under the existing TTY gate, including terminal Runs as a read-only
browser exited with q. Non-TTY attach output stays byte-compatible with today.

HITL checkpoint: interaction design review — focus marking, status bar states, detail pane, and footer copy —
before the live commands adopt the cockpit.

## Acceptance criteria

- [ ] Attach in a TTY opens the interactive cockpit; `ROUNDFIX_TUI=never` and non-TTY keep today's text contract
- [ ] Tab focus model with visible focus; timeline scroll keys and issue selection keys work as specified
- [ ] Enter opens the issue detail pane; Esc returns; missing artifact degrades without failing
- [ ] q and Ctrl-C detach; the Run is untouched; no stop keybinding exists in attach
- [ ] Terminal Runs open read-only with full scrollback and exit on q
- [ ] Status bar shows REPLAYING / FOLLOWING / SCROLLED · N below / terminal state
- [ ] Key handling proven by driving the model with synthetic key messages; no real terminal in tests
- [ ] Interaction design review completed

## Blocked by

- 02-timeline-window-follow-state.md
