---
task: task_05
spec: 0005-tui-cockpit
status: completed
type: frontend
complexity: high
---

# Task 05: Open Work Item detail as a centered modal

## Overview

Detail becomes contextual inspection: Enter opens the selected Work Item's
detail as a centered terminal modal over the lightly dimmed cockpit, `D`
toggles it, `Esc` closes back into byte-identical context. Review Runs show
the Review Issue artifact; spec Runs show the task file body — closing
finding 5. Verifiable through synchronous key-routing and render tests.

## Requirements

1. MUST render the modal per the visual contract (`design/roundfix-02.png`):
   centered box with title, status/severity/location line, source reference,
   scrollable body, and a scroll-position footer (`Line a-b of N ·
   PgUp/PgDn page`).
2. MUST route keys: `Enter` opens for the selected Work Item, `D` toggles,
   `Esc` closes, `j/k`/arrows scroll, `PgUp/PgDn` page; while open, queue
   navigation is suspended and Follow Mode keeps advancing underneath;
   closing restores the exact prior render (byte-compared in tests).
3. MUST load detail per Run kind: the Review Issue artifact through the
   existing detail loading, or the task file body read-only with the pane's
   mid-write tolerance (parse failure keeps last good content with a stale
   marker line).
4. MUST dim the background lightly (border/emphasis over heavy shading) and
   keep the modal readable at small widths, degrading to full-surface detail
   below the minimum.
5. MUST keep attach and terminal modes' existing key differences.

## Subtasks

- [x] Modal renderer over the dimmed base
- [x] Key routing: open, toggle, close, scroll, focus suspension
- [x] Detail sources for both Run kinds with mid-write tolerance
- [x] Byte-identical close-context tests and small-width fallback

## Acceptance Criteria

- [x] Key tests: Enter opens, D toggles, Esc closes, scroll and paging move
      the body, queue keys are inert while open.
- [x] Render before opening equals render after closing, byte-compared, for
      both Run kinds.
- [x] Task-file detail shows the body read-only and survives a corrupted
      mid-write reload with the stale marker.
- [x] Full suite passes.

## Verification

- `rtk go test ./internal/tui/` — expected: all tests pass.
- `rtk go test ./...` — expected: full suite passes.

## References

`_prd.md` → User Story 3; Core Feature 5; Decisions (Enter/D/Esc).
`_techspec.md` → Detail sources, Build Order 5, Risks (dimming).
`design/ui-redesign-plan.md` → Required Changes (modal);
`design/roundfix-02.png`. Dogfood finding 5.

## Result

- Modal renderer evidence: `TestCockpitEnterOpensIssueDetailModalAndEscRestoresReviewContext`
  asserts the Work Queue and timeline remain visible behind `REVIEW.ISSUE`
  modal content, with a line-position footer.
- Key-routing evidence: the same test asserts Enter opens, Esc closes,
  queue selection stays fixed while arrow keys scroll the modal, and PgDn
  pages the body; `TestCockpitDetailTogglesWithDAndKeepsTimelineFollowing`
  asserts `D` opens/closes and the timeline follows while the modal is open.
- Spec detail evidence: `TestCockpitSpecTaskDetailModalRestoresContextAndSurvivesStaleReload`
  asserts `SPEC.TASK` renders the task file read-only, keeps the last readable
  body through a corrupted mid-write reload, and shows the stale marker.
- Context/fallback evidence: review and spec tests byte-compare render before
  opening against render after closing; `TestCockpitDetailUsesFullSurfaceFallbackWhenTerminalIsTooShort`
  covers the small-height full-surface fallback; `TestCockpitDetailKeepsAttachDetachKey`
  covers attach-mode `q` while a modal is open.
- Snapshot evidence: the detail-open cockpit golden snapshots were updated
  through `ROUNDFIX_UPDATE_COCKPIT_SNAPSHOTS=1 go test`.
- Verification: `rtk go test ./internal/tui/` passed with 78 tests.
- Verification: `rtk go test ./...` passed with 558 tests across 16 packages.
