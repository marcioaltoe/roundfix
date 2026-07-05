---
task: task_01
spec: 0005-tui-cockpit
status: pending
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

- [ ] Snapshot assertions over the shipped cockpit states
- [ ] Per-surface helper extraction
- [ ] Neutrality proof: snapshots green post-refactor

## Acceptance Criteria

- [ ] Snapshot tests exist for the listed states/sizes and pass before and
      after the refactor with identical expected strings.
- [ ] No exported API changes; no behavior or key-routing changes.
- [ ] Full suite passes.

## Verification

- `rtk go test ./internal/tui/` — expected: all tests pass, including the new
  snapshots.
- `rtk go test ./...` — expected: full suite passes.

## References

`_prd.md` → Goals; Decisions (refactor-before-feature). `_techspec.md` →
Executive Summary, Interfaces, Build Order 1, Risks (neutrality).
`design/ui-redesign-plan.md` → Suggested Implementation Order 1–2.
