---
task: task_03
spec: 0021-cockpit-visual-fidelity
status: pending
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

- [ ] Group markers and automatic expand/collapse by Batch status
- [ ] Timestamp gutter alignment and kind-colored labels
- [ ] Summary-only event rows through the shared helper
- [ ] `Live · detail` indicator wired to modal state
- [ ] Tests: expanded/collapsed groups, gutter alignment, a multi-line tool
      payload rendering one bounded line, indicator states, no-color twins

## Acceptance Criteria

- [ ] A settled Batch renders collapsed with ▶ and its summary row; the
      executing Batch renders expanded with ▼ and its events.
- [ ] A Run Event whose summary carries a multi-line tool payload renders as
      one line bounded to the pane width, and the full payload remains
      reachable in the Detail Modal.
- [ ] Timestamps align in one gutter column across kinds.
- [ ] The pane header indicator flips between `detail hidden` and
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
