---
task: task_02
spec: 0005-tui-cockpit
status: pending
type: frontend
complexity: medium
---

# Task 02: Build the two-pane layout and the phase row

## Overview

Replace the current cockpit arrangement with the mocked base structure: Work
Queue left, session timeline right (the widest surface), topped by a phase
row with terminal-text markers — with review and spec variants. Verifiable
through render tests for both Run kinds against the mockup structure.

## Requirements

1. MUST render the two persistent surfaces per the visual contract
   (`design/roundfix-01.png`): queue left, timeline right and dominant, no
   persistent third pane.
2. MUST render the phase row with text markers `[done] [run] [wait] [locked]`:
   review Runs `FETCH > TRIAGE > AGENT > VERIFY > PUSH` (push `[locked]`
   until no Unresolved Review Issues remain); spec Runs `AGENT > VERIFY >
   COMMIT > QA` per current Task position, `QA` `[locked]` until every Task
   is completed and omitted when the Run has no QA opt-in — derived from Run
   state and journal events, no new state.
3. MUST keep both Run kinds flowing through the same layout path (the
   shipped WorkItem mapping), retiring the minimal spec-Run pane.
4. MUST use monospace text, terminal borders, and compact spacing only — no
   icons, logos, or decorative elements.

## Subtasks

- [ ] Base two-pane layout renderer
- [ ] Phase derivation for both Run kinds
- [ ] Phase row renderer with the four markers
- [ ] Render tests for review and spec fixtures

## Acceptance Criteria

- [ ] Review fixture renders queue+timeline+phase row matching the mockup's
      structure (asserted structurally: pane titles, phase sequence, marker
      states — not pixel art).
- [ ] Spec fixture renders the spec phase sequence with QA locked/omitted
      cases covered.
- [ ] The timeline pane is the wider surface at every tested size.
- [ ] Full suite passes; task_01 snapshots updated deliberately (this task
      is the visual change).

## Verification

- `rtk go test ./internal/tui/` — expected: all tests pass.
- `rtk go test ./...` — expected: full suite passes.

## References

`_prd.md` → User Stories 1, 2; Core Features 1, 2. `_techspec.md` → Phase
derivation, Interfaces, Build Order 2. `design/ui-redesign-plan.md` →
Required Changes 1; `design/roundfix-01.png`.
