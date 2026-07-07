---
task: task_02
spec: 0021-cockpit-visual-fidelity
status: pending
type: frontend
complexity: high
---

# Task 02: Header, Phase Row, and Work Queue fidelity

## Overview

Bring the cockpit's left half to the approved design (`roundfix-01.png`):
the header with the Run id and state chip, the color-coded Phase Row, and a
Work Queue of Work Item cards — colored severity, ordinal, muted location,
an accent-bordered selected card, Batch separators with elapsed stamps, and
the totals footer. Demoable by attaching to any recorded Run.

## Requirements

1. MUST render the header per the design: title segments, Run id and state
   chip on the right, using the section-label and state tokens.
2. MUST color Phase Row markers by state through the tokens while keeping
   the exact `[done]`/`[run]`/`[wait]`/`[locked]` marker text.
3. MUST render each Work Item as a card: state marker, severity in its state
   color, ordinal, title, and location in muted gray; the selected card
   carries the accent selection border and is unmistakable.
4. MUST render Batch separators with an elapsed stamp derived from the
   Batch's Run Event timestamps (running Batches count against now), and
   keep the totals footer.
5. MUST preserve every existing key, layout region, and text marker; the
   no-color render keeps all distinctions via markers.

## Subtasks

- [ ] Header and state chip through tokens
- [ ] Phase Row state coloring
- [ ] Work Item cards with severity color, muted location, selection border
- [ ] Batch separators with elapsed stamps; totals footer
- [ ] Tests: styled renders pinned under a forced profile plus no-color
      twins for header, phases, card, selected card, and batch stamp

## Acceptance Criteria

- [ ] Under a forced color profile, the selected card's render differs from
      an unselected card's only by the documented selection styling, and both
      pin their byte shape in tests.
- [ ] Severity words render in their state colors; locations render muted.
- [ ] A Batch with events at T and T+38s renders a `00:38`-shaped stamp.
- [ ] With color disabled, the same fixtures render marker-only and remain
      distinguishable.

## Verification

- `rtk go test ./internal/tui/` — expected: all tests pass, including the new
  queue fidelity tests.
- `rtk make verify` — expected: fmt-check, tests, skills check, and build all
  pass.

## References

`_prd.md` → User Stories 1-2; Core Features 2-3. `_techspec.md` → Coverage
Map; Build Order 2. Design refs:
`../_archived/0005-tui-cockpit/design/roundfix-01.png`.
