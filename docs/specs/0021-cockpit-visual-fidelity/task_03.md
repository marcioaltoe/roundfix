---
task: task_03
spec: 0021-cockpit-visual-fidelity
status: completed
type: frontend
complexity: high
---

# Task 03: Timeline fidelity: groups, gutter, bounded summaries

## Overview

Bring the Session Timeline to the approved design: one group per Batch and
event kind with ▼/▶ markers — executing Batch expanded, settled Batches
collapsed to summary rows — an aligned timestamp gutter, kind-colored labels
(PLAN, TOOL, THINK, SESSION, Daemon milestones), exactly one bounded summary
line per event, and the `Live · detail hidden/open` indicator in the pane
header. Demoable by attaching to any Run with Agent events.

## Requirements

1. MUST group timeline rows by Batch and event kind with a ▼ marker on
   expanded groups and ▶ on collapsed ones; the executing Batch renders
   expanded and settled Batches collapse to their summary rows
   automatically — no new keybindings.
2. MUST render an aligned timestamp gutter and kind labels colored through
   the tokens.
3. MUST render every event as exactly one bounded summary line via the
   shared summary helper — raw payloads (tool JSON, markdown bodies) never
   render inline; the Detail Modal keeps full content.
4. MUST show the `Live · detail hidden` / `Live · detail open` indicator in
   the timeline pane header, following the modal state.
5. MUST preserve scrolling, follow semantics, and every existing key; the
   no-color render keeps all distinctions via markers.

## Subtasks

- [x] Group markers and automatic expand/collapse by Batch status
- [x] Timestamp gutter alignment and kind-colored labels
- [x] Summary-only event rows through the shared helper
- [x] `Live · detail` indicator wired to modal state
- [x] Tests: expanded/collapsed groups, gutter alignment, a multi-line tool
      payload rendering one bounded line, indicator states, no-color twins

## Acceptance Criteria

- [x] A settled Batch renders collapsed with ▶ and its summary row; the
      executing Batch renders expanded with ▼ and its events.
- [x] A Run Event whose summary carries a multi-line tool payload renders as
      one line bounded to the pane width, and the full payload remains
      reachable in the Detail Modal.
- [x] Timestamps align in one gutter column across kinds.
- [x] The pane header indicator flips between `detail hidden` and
      `detail open` with the modal.

## Verification

- `rtk go test ./internal/tui/` — expected: all tests pass, including the new
  timeline fidelity tests.
- `rtk make verify` — expected: fmt-check, tests, skills check, and build all
  pass.

## References

`_prd.md` → User Story 3; Core Feature 4. `_techspec.md` → Interfaces:
EventSummary; Build Order 3; Risks (summary robustness). Design refs:
`../_archived/0005-tui-cockpit/design/roundfix-01.png`.

## Result

Status: completed — 2026-07-07.

### What changed

- Batch groups now carry state-driven markers: every Batch group header
  renders `▼ BATCH nnn/nnn [state] [elapsed]`, and a Batch whose journaled
  daemon.batch state word is `completed`, `failed`, or `stopped` collapses
  to that single summary row (`▶ …`) — its event rows disappear from the
  flow while scrollback and the journal keep them. Collapse derives purely
  from Batch state; no keybinding was added.
- Structured events (agent plan/thought/tool/status and every daemon
  milestone kind) render as exactly one row: a 9-column aligned timestamp
  gutter (`HH:MM:SS ` or blanks when the event has no timestamp), the kind
  label (daemon labels prepended; agent summaries self-label), and the
  event's bounded summary line through the shared `EventSummary` helper.
  Raw payload lines (tool commands, diffs, plan entries, markdown bodies)
  never render inline. Events journaled without a summary fall back to
  their reconstructed console text, bounded by the same helper, so older
  journals stay viewable; the shipped skip policy (unknown kinds, session
  lifecycle statuses) is untouched.
- Chunked `agent.message`/`agent.raw` streams keep the shipped coalescing
  contract (0005-pinned: fragments join into whole console lines) — chunk
  fragments cannot meaningfully summarize per event, and the coalescing
  tests continue to pass unchanged.
- The timeline pane header is now pinned (it no longer scrolls away under
  load) and carries the `Live · detail hidden` / `Live · detail open`
  indicator, following the Detail Modal state, muted, right-aligned next to
  the `SESSION.TIMELINE` section label.
- Pane rows are colored through the tokens: batch header label cyan, state
  word in its state token (started/executing/running amber, waiting amber,
  completed green, failed/stopped red), elapsed stamps and the gutter
  muted, `PLAN`/`[TOOL]`/`SESSION` and daemon milestone labels cyan, THINK
  rows muted, streaming text plain. With identity tokens the rows render
  byte-identical marker-only text.
- Scrolling, follow semantics, and every key are untouched — the whole
  viewport suite (follow, freeze, paging, window slides) passes unchanged.

### Commands run

- `rtk go test ./internal/tui/` — 151 passed, 0 failed.
- `rtk make verify` — fmt-check, full test suite (936 passed in 19
  packages), `roundfix skills check`, and build all passed.

### Evidence per acceptance criterion

1. Collapse/expand: `TestTimelineSettledBatchCollapsesAndExecutingBatchExpands`
   pins a settled Batch to exactly `▶ BATCH 001/002 completed 00:38` with
   its event rows absent, and the executing Batch to `▼ BATCH 002/002
   started 00:10` followed by its event rows.
2. Bounded tool payload: `TestTimelineToolPayloadRendersOneBoundedLineInThePane`
   feeds a tool event whose summary and payload carry `$ git apply` and a
   diff body; the full cockpit render shows only
   `12:00:50 [TOOL] apply_patch · completed`, none of the payload lines,
   and every rendered row bounded to the pane width. Full content stays
   reachable in the Detail Modal, which keeps rendering complete bodies
   (`TestCockpitEnterOpensIssueDetailModalAndEscRestoresReviewContext`
   scrolls through the full multi-line artifact body, passing unchanged).
3. Gutter alignment: `TestTimelineGutterAlignsAcrossKinds` pins PLAN, TOOL,
   and TASK rows byte-exactly with one 9-column gutter — timestamped and
   timestamp-less events share the same column.
4. Indicator: `TestTimelinePaneHeaderIndicatorFollowsModalState` flips
   `Live · detail hidden` → `Live · detail open` on Enter and back on Esc.
   No-color twins: `TestTimelineRowsStyledThroughTokensAndNoColorTwin`
   byte-pins the styled rows via token-constructed expectations and asserts
   the identity render is ANSI-free, marker-only, and unchanged.

### Notes and follow-ups

- Techspec's "Daemon milestones per state color" is applied where state is
  knowable from the row: the Batch header state word carries its state
  token; per-event daemon labels render as cyan section labels (their
  summaries carry the state in text). Flagging in case the design review
  wants per-label state mapping instead.
- The old-contract assertions that pinned raw payload lines inline
  (`TestViewportGroupsReviewTimelineByBatchAndKind`, the two Attach replay
  tests) were updated to the new summary-first contract — the degradation
  the PRD names — and now assert the payload text does NOT render.
- Empty-state copy ("No Run Events yet...") kept verbatim; explanatory
  empty states are task_04 scope.
