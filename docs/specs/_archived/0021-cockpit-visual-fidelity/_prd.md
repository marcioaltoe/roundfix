---
spec: 0021-cockpit-visual-fidelity
status: archived
created: 2026-07-07
surfaces: [cli, docs]
archived: "2026-07-07"
source_slug: 0021-cockpit-visual-fidelity
---


# Cockpit Visual Fidelity

Spec 0005 shipped the cockpit's information architecture — two panes, Phase
Row, Work Queue, timeline grouping, Detail Modal — and its QA passed against
that written contract. The visual layer of the approved design
(`docs/specs/_archived/0005-tui-cockpit/design/roundfix-01.png` and
`roundfix-02.png`, with the style rules in `design/ui-redesign-plan.md`)
never became verifiable acceptance criteria, so it silently degraded: the
shipped cockpit is near-monochrome, the timeline renders raw Agent payloads
instead of summary rows, selection is a bare marker instead of a highlighted
card, and empty states explain nothing. Cockpit Visual Fidelity closes that
gap: the cockpit renders to the approved design, and this time the styled
render is the acceptance criterion.

## Goals

- The cockpit visually matches the approved design references: color-coded
  state at a glance, card-style Work Items with a highlighted selection,
  grouped timeline with aligned gutters and summary-first rows, and a styled
  Detail Modal.
- The style system is one shared token set that future TUI surfaces (the Run
  Browser) consume, so the visual language cannot fork per view.
- Color degrades honestly: `NO_COLOR` / `ROUNDFIX_COLOR=never` keeps every
  state readable through the existing text markers.

## User Stories

1. As a user watching a Run, I want phase, batch, and issue state to be
   readable by color and marker at a glance — green done, amber running or
   waiting, red locked or failed, cyan section labels — so that I stop
   parsing prose to know where the Run is.
2. As a user navigating the Work Queue, I want the selected Work Item to be
   an unmistakably highlighted card with severity color, ordinal, and muted
   location, so that selection and priority are visible without reading.
3. As a user reading the timeline, I want events grouped with a ▼/▶ marker,
   an aligned timestamp gutter, kind-labeled rows (PLAN, TOOL, THINK,
   SESSION), and one bounded summary line per event — raw payloads stay
   behind the Detail Modal — so that the timeline reads like the design, not
   like a log dump.
4. As a user attaching to a Run with no Agent events (a Fetch Run, an early
   Run), I want the empty panes to say what they are waiting for or why they
   are empty, so that an empty screen is information instead of confusion.
5. As a user on a no-color terminal, I want every state distinction to
   survive through text markers alone, so that the design degrades without
   losing meaning.

## Core Features

1. A single style token set implementing the approved rules: cyan for active
   borders and section labels, green for completed, amber for running,
   waiting, and pending, red only for locked, failed, or blocking, muted
   gray for timestamps, file paths, and background text. Tokens honor the
   existing `ROUNDFIX_COLOR` / `NO_COLOR` contract.
2. Header and Phase Row per the design: the title bar with Run id and state
   chip on the right, phase markers colored by state.
3. Work Queue per the design: Batch separators with an elapsed-time stamp,
   Work Item cards showing marker, colored severity, ordinal, title, and
   muted location, the selected card highlighted with an accent border, and
   the totals footer.
4. Timeline per the design: one group per Batch and event kind with ▼
   (expanded) and ▶ (collapsed) markers — the executing Batch expanded,
   settled Batches collapsed to summary rows — an aligned timestamp gutter,
   kind-colored labels, and exactly one bounded summary line per event. Full
   payloads remain in the Detail Modal. The pane header carries the
   `Live · detail hidden/open` indicator.
5. Detail Modal per the design: accent border, title block, sectioned body,
   and a position footer.
6. Explanatory empty states for both panes, naming the Run kind and what
   would populate them.
7. Every visual behavior above is pinned by tests that assert the styled
   render (with a forced color profile) and the no-color fallback.

## User Experience

- No layout re-architecture: the 0005 pane structure, keys, and modal
  semantics stay; this spec changes how they render.
- Collapse state is automatic (executing expanded, settled collapsed); no
  new keybindings.
- Terminal-native throughout: monospace, box borders, text markers — no
  icons, gradients, or decorative elements.

## Non-Goals / Out of Scope

- The Run Browser (spec 0020) — it consumes the tokens, it is not built
  here.
- Changes to Run Event storage, event kinds, or the agent contract.
- New interactive controls (manual group collapsing, mouse support).
- Web anything.

## Success Metrics

- Side by side with `roundfix-01.png` and `roundfix-02.png`, the cockpit
  shows the same visual hierarchy: colored phase row, carded queue with
  highlighted selection, grouped gutter-aligned timeline, styled modal.
- A raw Agent payload never renders inline in the timeline — only its
  summary line.

## Decisions

- The style rules from `ui-redesign-plan.md` §7 are the token contract; the
  PNGs are the fidelity reference.
- Styled-render tests are the enforcement mechanism this time — structure
  tests alone allowed the 0005 degradation.
- Group collapsing is state-driven, not interactive (KISS; the design shows
  states, not a control).

## Open Questions

None.
