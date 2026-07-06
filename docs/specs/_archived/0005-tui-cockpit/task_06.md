---
task: task_06
spec: 0005-tui-cockpit
status: completed
type: frontend
complexity: low
---

# Task 06: Footer states and responsive fallback

## Overview

Finish the chrome: the footer shows the right key hints for normal, modal,
attach, and terminal states, and the layout degrades gracefully — small
terminals collapse to the timeline with a one-line queue summary, degenerate
sizes never panic. Verifiable through footer and sizing table tests.

## Requirements

1. MUST render footer hints per state: normal (`Tab focus · ↑↓ move/scroll ·
   PgUp/PgDn page · Enter issue · D show detail · End follow · Ctrl-C stop`
   — vocabulary adjusted per Run kind for Tasks), modal (`Esc close · j/k
   scroll · PgUp/PgDn page`), attach and terminal variants keeping their
   existing differences.
2. MUST collapse below the minimum two-pane width to a single timeline
   surface with a one-line queue summary and a footer hint to widen;
   the modal falls back per task_05.
3. MUST survive degenerate dimensions (zero/negative width or height) with
   an empty render, never a panic — table-tested.
4. MUST keep width/height distribution stable at common sizes (regression
   asserts at 80×24, 120×40, 200×50).

## Subtasks

- [x] Footer renderer per state and Run kind
- [x] Small-terminal collapse and minimum-size thresholds
- [x] Degenerate-size safety table tests

## Acceptance Criteria

- [x] Footer tests pin each state's hint line for both Run kinds.
- [x] An 80×24 render is the collapsed fallback; 120×40 is two-pane; both
      asserted structurally.
- [x] Zero and negative dimensions return empty strings without panicking.
- [x] Full suite passes.

## Verification

- `rtk go test ./internal/tui/` — expected: all tests pass.
- `rtk go test ./...` — expected: full suite passes.

## References

`_prd.md` → User Story 6; User Experience. `_techspec.md` → Build Order 6,
Risks (degenerate sizes). `design/ui-redesign-plan.md` → Required Changes
(footer, responsive behavior).

## Result

- Footer evidence: `TestCockpitFooterHintsForStatesAndRunKinds` pins normal,
  modal, attach, and terminal footer hints for review and spec Runs, including
  Task vocabulary for spec Runs and the existing `q detach` / `q close`
  variants.
- Responsive evidence: `TestCockpitResponsiveFallbackAndStableSizes` asserts
  80×24 collapses to `SESSION.TIMELINE` plus a `QUEUE.SUMMARY` line and widen
  footer hint, while 120×40 and 200×50 keep the two-pane Work Queue + timeline
  layout with stable widths.
- Degenerate-size evidence: `TestCockpitDegenerateSizesRenderEmptyWithoutPanic`
  covers zero and negative width/height cases and asserts an empty render.
- Snapshot evidence: cockpit snapshots were updated deliberately for the
  task_06 footer copy and still pass with the responsive renderer.
- Verification: `rtk go test ./internal/tui/` passed, including the new footer
  and sizing table tests.
- Verification: `rtk go test ./...` passed with 577 tests across 16 packages.
