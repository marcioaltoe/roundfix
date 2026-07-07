---
task: task_02
spec: 0021-cockpit-visual-fidelity
status: completed
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

- [x] Header and state chip through tokens
- [x] Phase Row state coloring
- [x] Work Item cards with severity color, muted location, selection border
- [x] Batch separators with elapsed stamps; totals footer
- [x] Tests: styled renders pinned under a forced profile plus no-color
      twins for header, phases, card, selected card, and batch stamp

## Acceptance Criteria

- [x] Under a forced color profile, the selected card's render differs from
      an unselected card's only by the documented selection styling, and both
      pin their byte shape in tests.
- [x] Severity words render in their state colors; locations render muted.
- [x] A Batch with events at T and T+38s renders a `00:38`-shaped stamp.
- [x] With color disabled, the same fixtures render marker-only and remain
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

## Result

Status: completed — 2026-07-07.

### What changed

- The cockpit now resolves its token set once per view: `CockpitConfig`
  gained `ColorEnabled`, the model carries `Tokens` from
  `ResolveTokens(cfg.ColorEnabled)`, and both CLI entry points (`runui.go`
  owning cockpit on stderr, `attach.go` attach cockpit on stdout) forward
  the existing `ROUNDFIX_COLOR` / `NO_COLOR` resolution from
  `internal/cli/style.go`. Nothing re-reads env vars in the TUI.
- Header: `ROUNDFIX` renders through `SectionLabel`, the remaining title
  segments and the Run id through `Muted`, and the state chip through a
  run-state token — green for Clean/Fetched/IntegrationPending, red for
  every other terminal outcome, amber while the Run moves. Text bytes are
  unchanged. Call sites no longer add `Bold` on top of tokens, so identity
  tokens render truly marker-only.
- Phase Row: `[done]`/`[run]`/`[wait]`/`[locked]` marker text is untouched
  and colored via `Done`/`Running`/`Waiting`/`Locked`; phase names render
  plain, separators muted.
- Work Item cards: marker plus severity render in the item's state token
  (Executing amber, Resolved/Completed green, Failed red, Paused red,
  Invalid/Duplicated/Skipped muted, Waiting amber); ordinal and location
  render muted; title renders plain. The selected card is wrapped in the
  `Selection` accent border (structural in no-color mode) and keeps the `>`
  marker.
- Batch separators now derive their elapsed stamp from the Batch's Run
  Event timestamps (first-to-last event; the executing Batch of a live Run
  counts against `now`), replacing the old poll-detection wall clock
  (`batchStartedAt` removed). The model folds journal events into per-Batch
  spans through an incremental cursor scan on each poll. Labels render
  `SectionLabel`, stamps `Muted`. The totals footer stays, rendered muted.
- All existing keys, layout regions, and text markers preserved; golden
  snapshots regenerated to carry the selected-card border, batch stamps,
  and the timestamped fixture's timeline grouping.

### Commands run

- `rtk go test ./internal/tui/` — 146 passed, 0 failed.
- `rtk make verify` — fmt-check, full test suite (931 passed in 19
  packages), `roundfix skills check`, and build all passed.

### Evidence per acceptance criterion

1. Selection styling pinned: `TestCockpitWorkItemCardPinsSelectionStyling`
   byte-pins the unselected and selected card renders against expectations
   constructed from the tokens themselves; the selected expectation is the
   same content lines (with the `>` marker) wrapped in `Tokens.Selection` —
   proving the render differs only by the documented selection styling.
2. Severity/state colors and muted locations: the same test pins
   `Tokens.Running`-rendered `  [run] MAJOR` and `Tokens.Muted`-rendered
   location/ordinal segments; header/phase colors pinned by
   `TestCockpitHeaderRendersSectionLabelAndStateChipTokens` and
   `TestCockpitPhaseRowColorsMarkersThroughTokens`.
3. `00:38` stamp: `TestCockpitBatchSeparatorStampDerivesFromRunEventTimestamps`
   feeds Batch 1 events at T and T+38s and pins the separator's
   `BATCH 001/002 … 00:38` shape; a Batch without timestamped events shows
   no stamp; `TestCockpitWorkQueueBatchSeparatorShowsOrdinalAndElapsed`
   covers the executing Batch counting against now (01:23).
4. No-color twins: `TestCockpitHeaderNoColorRendersMarkerOnlyText`,
   the identity half of `TestCockpitPhaseRowColorsMarkersThroughTokens`,
   `TestCockpitWorkItemCardNoColorKeepsMarkerDistinctions`, and the plain
   half of the batch stamp test assert zero ANSI, marker parity with the
   styled render, and that the selected card stays distinguishable via the
   `>` marker plus structural border.

### Follow-ups

- At the 88x24 minimum two-pane size the taller selected card can push the
  last Work Item below the fold, leaving its Batch separator visible as a
  "more below" hint (pre-existing row-budget behavior, now easier to hit).
  If tighter packing is wanted, that is a layout tweak for a later slice.
- `SESSION.TIMELINE` labels, collapsed-pane summary, and Detail Modal still
  render through the legacy global styles by design — task_03/task_04
  scope.
