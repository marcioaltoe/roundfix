---
title: "Interactive cockpit on attach: focus, keys, issue detail"
type: HITL
category: enhancement
state: completed
labels:
  - enhancement
  - completed
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

- [x] Attach in a TTY opens the interactive cockpit; `ROUNDFIX_TUI=never` and non-TTY keep today's text contract
- [x] Tab focus model with visible focus; timeline scroll keys and issue selection keys work as specified
- [x] Enter opens the issue detail pane; Esc returns; missing artifact degrades without failing
- [x] q and Ctrl-C detach; the Run is untouched; no stop keybinding exists in attach
- [x] Terminal Runs open read-only with full scrollback and exit on q
- [x] Status bar shows REPLAYING / FOLLOWING / SCROLLED · N below / terminal state
- [x] Key handling proven by driving the model with synthetic key messages; no real terminal in tests
- [x] Interaction design review completed

## Blocked by

- 02-timeline-window-follow-state.md

## Comments

**2026-06-10 (agent):** Built the interactive cockpit (`internal/tui/cockpit.go`): a Bubble Tea model hosting the
`TimelineViewport` over a `CockpitSource` (TimelineSource + DataVersion + Run lookup; `*store.Store` satisfies it,
journal-only per ADR 0009). Mode-aware key ownership is plumbed now: attach → q/Ctrl-C detach (`tea.Quit`), no stop
key; owning → q ignored, Ctrl-C invokes the `OnStop` callback (wired in issue 04), contextual footer differs.
Tab switches focus with the focused panel border highlighted; timeline keys (↑↓/j/k, PgUp/PgDn, Home, End/G) drive
the viewport; issues keys (j/k bounded selection with sidebar windowing, Enter) open the read-only detail pane —
structured header (title, severity·status·file:line, Source Reference) plus the artifact's markdown body with the
YAML frontmatter stripped; Esc closes; missing artifacts degrade to "artifact not available". Follow polling is a
250ms `tea.Tick` gated by data_version (idle ticks read zero event rows — proven); a terminal transition flips the
status to `<STATE> · READ-ONLY` and stops ticking entirely. `runAttachCockpit` in `internal/cli/attach.go` opens
the alt-screen cockpit when `liveTUIEnabled(stdout)` — including terminal Runs as read-only browsers — and prints
the detach/terminal closing line after exit; non-TTY and `ROUNDFIX_TUI=never` keep the existing text contract
byte-identical (all prior attach/follow tests pass unchanged). Preview-driven refinement: scrolling when everything
fits on screen no longer freezes Follow Mode. Tests drive the model with synthetic `tea.KeyPressMsg` values
(validated against the same `Key.String()` matching used in production) — focus/selection/clamping, detail
open/close with a real persisted Round artifact, degraded missing-artifact pane, detach keys quitting with the Run
untouched, status-bar narration across FOLLOWING → SCROLLED·N → End-resume, idle-tick row-read suppression,
terminal read-only with no ticking, and owning-mode key differences. HITL interaction review: three rendered
states presented and **approved as-is**. Verification: `rtk go vet ./...` clean, `rtk go test -race` tui+cli 112
passed, full `rtk go test ./...` 226 passed in 15 packages, `rtk go run ./cmd/roundfix --help` green.
