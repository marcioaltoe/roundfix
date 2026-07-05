---
task: task_01
spec: 0005-tui-cockpit
status: completed
type: frontend
complexity: medium
---

# Task 01: Decompose the cockpit renderers, rendering-neutral

## Overview

Extract the cockpit's monolithic render path into separable per-surface
helpers — layout, header/phase area, work queue, timeline, footer, detail —
with zero visual change, proven by snapshot assertions captured before the
refactor. Every later task in this Spec diffs against these seams.
Verifiable by the snapshots passing before and after.

## Requirements

1. MUST capture current-behavior snapshot assertions first: full cockpit
   render strings for a review Run fixture (normal, detail-open, attach, and
   terminal modes) at two terminal sizes, plus the spec-Run pane as shipped.
2. MUST extract per-surface render helpers (base layout, header area, work
   queue, timeline, footer, detail) as package-internal functions of model
   state, leaving model state, key routing, and update logic untouched.
3. MUST keep the captured snapshots passing unchanged after the extraction —
   the task's whole point is provable neutrality.
4. MUST follow Bubble Tea v2 discipline: pure render functions, synchronous
   `model.Update` tests, no terminal emulation.

## Subtasks

- [x] Snapshot assertions over the shipped cockpit states
- [x] Per-surface helper extraction
- [x] Neutrality proof: snapshots green post-refactor

## Acceptance Criteria

- [x] Snapshot tests exist for the listed states/sizes and pass before and
      after the refactor with identical expected strings.
- [x] No exported API changes; no behavior or key-routing changes.
- [x] Full suite passes.

## Verification

- `rtk go test ./internal/tui/` — expected: all tests pass, including the new
  snapshots.
- `rtk go test ./...` — expected: full suite passes.

## References

`_prd.md` → Goals; Decisions (refactor-before-feature). `_techspec.md` →
Executive Summary, Interfaces, Build Order 1, Risks (neutrality).
`design/ui-redesign-plan.md` → Suggested Implementation Order 1–2.

## Result

Implemented the rendering-neutral cockpit decomposition behind snapshot
coverage. `TestCockpitRenderSnapshots` now compares the full visible cockpit
render for review normal, detail-open, attach, and terminal states at 88x24
and 120x40, plus the spec-Run pane at both sizes.

Evidence:

- Snapshot capture before refactor: `ROUNDFIX_UPDATE_COCKPIT_SNAPSHOTS=1 rtk go test ./internal/tui/ -run TestCockpitRenderSnapshots` passed with 11 snapshot cases.
- Pre-refactor snapshot proof: `rtk go test ./internal/tui/ -run TestCockpitRenderSnapshots` passed with 11 snapshot cases.
- Post-refactor neutrality proof: `rtk go test ./internal/tui/ -run TestCockpitRenderSnapshots` passed with the same golden strings.
- Focused verification: `rtk go test ./internal/tui/` passed, 52 tests.
- Full verification: `rtk go test ./...` passed, 532 tests across 16 packages.
- Repository gate: `rtk make verify` passed (`go test ./...`, `roundfix skills check`, and build).

Acceptance criteria evidence:

- Snapshot tests cover the listed states and sizes, including the shipped
  spec-Run pane, and passed before and after helper extraction without
  changing the golden files.
- The production diff adds only package-internal render helpers and a
  `cockpitLayout` struct in `internal/tui`; exported APIs, model state,
  key routing, and update logic were not changed.
- The full suite and repository verification gate passed.

Follow-ups: none.
