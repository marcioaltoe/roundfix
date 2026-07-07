---
task: task_04
spec: 0021-cockpit-visual-fidelity
status: pending
type: frontend
complexity: medium
---

# Task 04: Detail Modal fidelity and pane empty states

## Overview

Bring the Detail Modal to the approved design (`roundfix-02.png`) — accent
border, title block with severity/status/location line, sectioned body,
position footer — and replace the bare empty panes with explanatory states
naming the Run kind and what would populate them. Demoable by opening a
Review Issue detail and by attaching to a Fetch Run.

## Requirements

1. MUST render the modal per the design: accent border, header naming the
   Work Item and the close/scroll keys, a title block with severity, status,
   and location styled through the tokens, sectioned body, and a
   `Line A-B of N` position footer.
2. MUST keep modal semantics unchanged: open/close keys, scrolling, dimmed
   cockpit behind, closing back into unchanged context.
3. MUST render explanatory empty states in both panes — naming the Run kind
   and what populates the pane (for example, a Fetch Run writes artifacts
   and starts no Agent) — instead of the bare placeholders.
4. MUST keep the no-color render fully readable via markers and text.

## Subtasks

- [ ] Modal frame, title block, sections, and position footer through tokens
- [ ] Empty states per pane naming Run kind and expectation
- [ ] Tests: modal render pinned (styled + no-color), empty states per Run
      kind, unchanged open/close/scroll behavior

## Acceptance Criteria

- [ ] The open modal renders the accent frame, styled title block, and
      position footer, pinned under a forced profile with a no-color twin.
- [ ] Attaching to a Fetch Run renders both panes with explanatory copy
      naming the Fetch Run behavior instead of bare placeholders.
- [ ] Existing modal behavior tests pass unchanged.

## Verification

- `rtk go test ./internal/tui/` — expected: all tests pass, including the new
  modal and empty-state tests.
- `rtk make verify` — expected: fmt-check, tests, skills check, and build all
  pass.

## References

`_prd.md` → User Story 4; Core Features 5-6. `_techspec.md` → Coverage Map;
Build Order 4. Design refs:
`../_archived/0005-tui-cockpit/design/roundfix-02.png`.
