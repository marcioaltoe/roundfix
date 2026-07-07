---
task: task_04
spec: 0021-cockpit-visual-fidelity
status: completed
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

- [x] Modal frame, title block, sections, and position footer through tokens
- [x] Empty states per pane naming Run kind and expectation
- [x] Tests: modal render pinned (styled + no-color), empty states per Run
      kind, unchanged open/close/scroll behavior

## Acceptance Criteria

- [x] The open modal renders the accent frame, styled title block, and
      position footer, pinned under a forced profile with a no-color twin.
- [x] Attaching to a Fetch Run renders both panes with explanatory copy
      naming the Fetch Run behavior instead of bare placeholders.
- [x] Existing modal behavior tests pass unchanged.

## Verification

- `rtk go test ./internal/tui/` — expected: all tests pass, including the new
  modal and empty-state tests.
- `rtk make verify` — expected: fmt-check, tests, skills check, and build all
  pass.

## References

`_prd.md` → User Story 4; Core Features 5-6. `_techspec.md` → Coverage Map;
Build Order 4. Design refs:
`../_archived/0005-tui-cockpit/design/roundfix-02.png`.

## Result

Status: completed — 2026-07-07.

### What changed

- The Detail Modal renders through the tokens end to end: the frame uses
  the `ActiveBorder` accent token (cyan in color mode, plain structural
  border without color), the title line splits into a `SectionLabel` Work
  Item label and muted `Esc close · j/k scroll` hints, the rule and source
  lines render muted, severity and status render in the artifact's state
  token (pending amber, resolved/completed green, failed red, in_progress
  amber, invalid/duplicated/skipped muted) with the location muted, markdown
  headings in the body render as section labels (sectioned body), and the
  `Line A-B of N · PgUp/PgDn page` footer renders muted. The dimmed cockpit
  behind the modal dims through the muted token. Every text byte is
  unchanged, so all shipped modal behavior tests (open/close via
  Enter/D/Esc, scrolling, context restore, full-surface fallback, stale
  degrade) pass without modification and the goldens did not move.
- Both panes replace their bare placeholders with explanatory empty states
  naming the Run kind and what populates the pane: a Fetch Run explains "A
  Fetch Run writes Review artifacts to disk and starts no Agent" in both
  panes, review Runs explain that Review Issues queue once the Run records
  them and that Agent/Daemon activity streams into the timeline, implement
  Runs explain the Spec's Task Graph and Task execution stream. Copy renders
  muted and truncated to the pane width.
- `attach` now forwards the Run kind for every Run (previously only
  implement Runs), so attaching to a Fetch Run can render its explanatory
  states; `specRunView` still keys strictly on the implement kind, so
  review rendering is unchanged.

### Commands run

- `rtk go test ./internal/tui/` — 154 passed, 0 failed.
- `rtk make verify` — fmt-check, full test suite (939 passed in 19
  packages), `roundfix skills check`, and build all passed.

### Evidence per acceptance criterion

1. Modal pinned: `TestCockpitDetailModalStyledThroughTokensWithNoColorTwin`
   byte-pins the title line, rule, subject, severity/status/location line,
   source line, section heading, and position footer via token-constructed
   expectations; asserts the frame carries the cyan border sequence; and
   pins the no-color twin as ANSI-free, byte-identical to the styled
   render's stripped text, with every marker preserved.
2. Fetch Run empty states: `TestCockpitFetchRunAttachRendersExplanatoryEmptyStates`
   attaches to a Fetched Run and asserts both panes carry the Fetch Run
   copy while the generic placeholders are gone;
   `TestCockpitEmptyStatesNameReviewAndSpecRunExpectations` covers the
   review and implement Run variants.
3. Existing modal behavior: the whole shipped modal suite
   (`TestCockpitEnterOpensIssueDetailModalAndEscRestoresReviewContext`,
   `TestCockpitDetailTogglesWithDAndKeepsTimelineFollowing`,
   `TestCockpitSpecTaskDetailModalRestoresContextAndSurvivesStaleReload`,
   `TestCockpitDetailUsesFullSurfaceFallbackWhenTerminalIsTooShort`,
   `TestCockpitDetailKeepsAttachDetachKey`) passes unchanged in the 154.

### Notes and follow-ups

- Pane focus is still signaled only by the border color (`panel()` chrome):
  in no-color mode the focused pane border is indistinguishable from the
  unfocused one. Pre-existing gap outside this task's modal/empty-state
  slice — flagging for a follow-up (a focus marker would fix it).
- The inactive pane border (gray 238) and progress-bar chrome remain global
  styles pending no per-surface owner; the modal itself is fully
  token-pure.
