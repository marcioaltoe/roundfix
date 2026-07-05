# Roundfix TUI Redesign Implementation Plan

## Status

Draft plan. This document captures the agreed terminal UI direction before
implementation. It should guide the code changes in `internal/tui/` without
changing the product workflow.

## Visual Contract

The Live Run View should become a two-surface cockpit:

- `REVIEW QUEUE` on the left.
- `SESSION.TIMELINE` on the right.
- `REVIEW.ISSUE` shown only as a modal overlay when the user asks for detail.

The main screen should keep the live stream visible at all times. Issue detail
is contextual inspection, not a persistent third pane.

Reference states:

- Normal state: issue queue plus expanded session timeline.
- Detail state: same screen dimmed behind a centered terminal modal.

The UI must remain a CLI/TUI. Use monospace text, terminal borders, compact
spacing, and text status markers such as `[run]`, `[done]`, `[wait]`, and
`[locked]`. Do not add provider logos, pictographic icons, browser-like
controls, or decorative dashboard elements.

## Product Fit

This redesign preserves the current Roundfix contract:

- `watch` and `resolve` still drive the Live Run View.
- The daemon still owns the loop.
- The child agent still owns code changes only.
- The cockpit still reads the Run Event Journal, not the live sink.
- Raw agent payloads remain stored as producer JSON; rendering improvements
  happen at read time.

The goal is clearer terminal rendering, not a new workflow system.

## Current Implementation Baseline

Relevant current code:

- `internal/tui/cockpit.go`
  - Owns Bubble Tea model state, focus, key handling, issue selection, detail
    state, layout, footer, progress bar, and issue pane rendering.
- `internal/tui/agent_live.go`
  - Renders agent event text and shared styles.
- `internal/tui/timeline.go`
  - Converts Run Events into timeline text.
- `internal/tui/viewport.go`
  - Owns replay, follow mode, paging, and scrollback behavior.
- `internal/agent/event.go`
  - Converts raw ACP updates back into stream updates for rendering.
- `internal/runevent/event.go`
  - Defines daemon and agent event kinds.
- `docs/adr/0008-run-event-payload-stores-raw-producer-json.md`
  - Requires raw agent payload preservation.
- `docs/adr/0009-cockpit-reads-the-journal-never-the-sink.md`
  - Requires journal-backed cockpit rendering.

The current cockpit already has the useful primitives:

- issue selection
- batch mapping through `BatchSizes`
- issue status labels
- detail loading from markdown artifacts
- detail scroll state
- viewport follow and scrollback
- terminal and attach mode key differences

The redesign should mostly reshape rendering and key semantics around those
existing primitives.

## Required Changes

### 1. Replace The Main Layout With Two Persistent Panes

Current behavior replaces the right pane with issue detail when `model.detail`
is open. Change the base layout so the normal view always renders:

- left issue queue
- right session timeline

The right pane should stay wide enough to be the primary operational surface.
The issue queue should be narrow but information-rich.

Suggested sizing:

- sidebar: 30-34 percent of the inner width
- timeline: remaining width
- minimum sidebar width: keep around the current lower bound
- small terminal fallback: preserve readability before visual fidelity

### 2. Redesign The Header And Phase Row

The header should communicate target and run state without looking like a web
dashboard.

Target shape:

```text
ROUNDFIX // LIVE RUN VIEW // PR #4                  RUN 12d39851  [RESOLVING WITH AGENT]

FETCH [done]  >  TRIAGE [done]  >  AGENT [run]  >  VERIFY [wait]  >  PUSH [locked]
```

This can replace or sit where the current `RUN.PROGRESS` bar lives. If both are
kept temporarily, the phase row should win as the final design because it better
matches the loop shape users need to understand.

The phase row should be derived from run state, issue progress, and daemon
events where available. It should degrade gracefully when data is unknown.

### 3. Strengthen The Review Queue

The queue should keep the selected issue obvious and show batch grouping
without requiring the detail modal.

Target structure:

```text
REVIEW QUEUE (4)

BATCH 001/002                         00:38
> [run] MEDIUM                         #1
  Use panic-safe logger initialization
  packages/backend/src/logging/logger.ts:42

  [done] LOW                           #2
  Avoid hardcoded token in tests
  packages/backend/src/__tests__/auth.test.ts:17

BATCH 002/002
  [wait] HIGH                          #1
  Unvalidated user input in RPC handler
  packages/backend/src/rpc/handlers/user.ts:88

4 issues total  *  1 resolved  *  3 unresolved
```

Implementation notes:

- Keep `batchSeparator` as the source of batch labels, but make it part of the
  richer issue-list rendering.
- Use existing issue metadata: title, severity, status, file, line, and source
  count if available.
- Use terminal text markers only. Avoid provider logos and decorative icons.
- Keep the selected row visually strong with border/accent treatment.

### 4. Make Session Timeline The Primary Surface

The timeline should remain journal-backed and follow-mode aware, but rendering
should group events into scan-friendly blocks.

Target groups:

- `BATCH 001/002 executing 00:38`
- `PLAN`
- `[TOOL] read_file * completed`
- `[TOOL] edit_file * in_progress`
- `THINK checking error paths`
- `[TOOL] go_test * queued`
- `SESSION RUNNING`
- `BATCH 002/002 waiting`

Implementation notes:

- Keep raw event storage untouched.
- Add a reader-side grouping layer over the existing Run Event stream.
- Prefer grouping by event `Batch`, then by event kind.
- Preserve unknown-event skipping behavior.
- Preserve viewport follow, scrollback, page up/down, and `End` to follow.
- Use daemon summaries for daemon events, and reconstructed agent text for
  agent events.

The grouping layer should not become a workflow engine. It is a renderer.

### 5. Render Issue Detail As A Modal Overlay

The current `model.detail` state should become a modal state instead of a
right-pane replacement.

Target behavior:

- `Enter` on an issue opens the modal.
- `D` toggles the modal for the selected issue.
- `Esc` closes the modal.
- `j/k`, arrow keys, `PgUp`, and `PgDn` scroll inside the modal while it is
  open.
- `Tab` should not move focus while the modal is open.
- The background remains visible but lower-contrast.

Target modal shape:

```text
+----------------------------------------------------------------+
| REVIEW.ISSUE  #001                     Esc close * j/k scroll  |
|----------------------------------------------------------------|
| Use panic-safe logger initialization                            |
| medium * pending * packages/backend/src/logging/logger.ts:42    |
| source: thread:CRRT_kwDO...                                     |
|                                                                |
| REVIEW COMMENT                                                  |
| "The global logger can panic if initialization fails..."         |
|                                                                |
| ARTIFACT                                                        |
| ## Review Comment                                               |
| ...                                                            |
|                                                                |
| REFERENCES                                                      |
| - packages/backend/src/logging/logger.ts:42                     |
|                                                                |
| Line 1-18 of 47 * PgUp/PgDn page                                |
+----------------------------------------------------------------+
```

Implementation notes:

- Continue loading artifact content through `openDetail`.
- Continue stripping YAML frontmatter through `artifactBodyLines`.
- Reuse `issueDetailView.scroll`.
- Add a `renderModal` or `renderDetailModal` helper.
- Compose the modal over the already-rendered base view.
- If true overlay composition becomes awkward in terminal text, render the base
  view with dimmed styles and place the modal in the same final string using
  Lip Gloss positioning.

### 6. Update Footer Hints

Footer should reflect the modal-first interaction:

Normal state:

```text
Keys: Tab focus * up/down move/scroll * PgUp/PgDn page * Enter issue * D detail modal * End follow * Ctrl-C stop
```

Detail modal open:

```text
Keys: j/k scroll detail * PgUp/PgDn page * Esc close * Ctrl-C stop
```

Attach mode and terminal-run mode should still preserve existing `q detach` and
`q close` semantics.

### 7. Keep Styling Terminal-Native

Style rules:

- Use cyan for active borders and section labels.
- Use green for completed/done.
- Use amber/yellow for running, waiting, and pending.
- Use red only for locked, failed, or blocking states.
- Use muted gray for timestamps, file paths, and background text.
- Avoid gradients, icons, logos, rounded web cards, and decorative elements.

This keeps the UI implementable with Lip Gloss styles and keeps the product
visually honest as a terminal tool.

## Suggested Implementation Order

1. Preserve current tests and snapshot-like expectations before changing
   layout.
2. Refactor rendering helpers in `internal/tui/cockpit.go` so base layout,
   issue queue, timeline, footer, and detail rendering are separable.
3. Implement the two-pane base layout.
4. Upgrade issue queue rows and batch separators.
5. Add the modal detail renderer while keeping `openDetail` and
   `issueDetailView` state.
6. Add `D` key behavior and update modal-specific key routing.
7. Improve timeline grouping by batch and event kind.
8. Update footer text for normal, attach, terminal, and modal states.
9. Tune responsive width/height behavior.
10. Run focused TUI tests, then the broad Go gate.

## Test Plan

Add or update tests around:

- issue status labels for pending, executing, waiting, paused, resolved,
  invalid, duplicated, and failed
- batch separator rendering and elapsed time placement
- selected issue row rendering
- modal open via `Enter`
- modal toggle via `D`
- modal close via `Esc`
- modal scroll handling with `j/k`, arrows, `PgUp`, and `PgDn`
- focus behavior while modal is open
- footer hints in normal, modal, attach, and terminal states
- timeline grouping for agent plan, tool, thought, status, and daemon events
- viewport follow behavior after grouping
- small terminal fallback rendering

Manual smoke after implementation:

```bash
rtk gofmt -w internal/tui/*.go
rtk go test ./...
rtk go run ./cmd/roundfix --help
```

If the implementation touches concurrency, journal polling, or store access,
also run:

```bash
rtk go test -race ./...
```

## Acceptance Criteria

- The normal Live Run View shows only the issue queue and session timeline.
- The selected issue remains visible and easy to scan.
- Batches are visible in the queue and in the timeline.
- The timeline remains the largest and most important surface.
- The issue detail opens as a centered terminal modal.
- Closing the modal returns to the same queue and timeline context.
- The cockpit still reads from the Run Event Journal only.
- Raw agent payload storage is unchanged.
- Unknown or newer event kinds remain skippable.
- Non-TTY output behavior is unchanged.
- The UI uses terminal text markers, not icons or provider logos.
- Existing stop, attach, close, scroll, and follow semantics still work.

## Non-Goals

- No web dashboard layout.
- No persistent third issue-detail pane.
- No provider logos or brand icons.
- No new review workflow.
- No source-thread mutation from the UI.
- No changes to the agent execution contract.
- No changes to Run Event storage format.
- No broad CLI framework migration.

## Open Decisions

- Whether `Enter` opens the modal and `D` toggles it, or whether only `D`
  should own detail behavior.
- Whether the phase row fully replaces `RUN.PROGRESS` or whether a compact
  resolved-count indicator remains in the header.
- Whether timeline grouping should be purely render-time, or whether daemon
  events should add slightly richer summaries for cleaner grouping.
- How much dimming is worth applying behind the modal before it harms terminal
  readability.
