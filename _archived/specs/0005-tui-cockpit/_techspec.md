---
spec: 0005-tui-cockpit
prd: _prd.md
created: 2026-07-05
---

# TUI Cockpit Redesign — Technical Spec

## Executive Summary

A rendering-and-keys reshape of `internal/tui` over primitives that already
exist: the cockpit model's selection/focus/detail state, the `WorkItem` view
model both Run kinds map into, the journal-backed timeline, and Follow Mode.
The accepted trade-off is a **refactor-before-feature step**: the cockpit's
monolithic render path is decomposed into separable renderers (layout, phase
row, work queue, timeline, footer, modal) before any visual change, so every
subsequent task is a small diff against a stable seam instead of edits inside
one large function. No architecture changes: journal-only reads (ADR-0009),
raw payloads untouched (ADR-0008), Bubble Tea v2 synchronous-update testing
throughout.

## System Architecture

All work lands in `internal/tui`; nothing outside it changes behavior.

- `cockpit.go` — model state stays (selection, focus, detail, batch mapping);
  rendering decomposes into per-surface helpers; key routing gains modal
  states.
- `timeline.go` / `agent_live.go` — render-time grouping by Batch and event
  kind; summaries and styles only, no event mutation.
- `viewport.go` — Follow Mode and scrollback as-is; the modal adds one more
  consumer of the existing scroll primitives.
- New small files per surface renderer (phase row, work queue, modal) —
  cohesive, package-internal, no exported API changes beyond what tests need.
- Data inputs are unchanged: Run row (kind, spec slug, PR identity), journal
  events, Review Issue artifacts via existing parsing, task files via
  `spec.ReloadTask` — the mid-write tolerance from 0001 stays.

## Implementation Design

### Interfaces

```go
// internal/tui — package-internal seams the tasks build against
type phase struct{ Name, Marker string } // Marker: done|run|wait|locked
func runPhases(kind store.Kind, s cockpitState) []phase

func renderPhaseRow(width int, phases []phase) string
func renderWorkQueue(width, height int, q workQueueState) string
func renderTimeline(width, height int, t timelineState) string   // exists, regrouped
func renderModal(width, height int, d detailState) string        // centered overlay
func renderFooter(width int, mode footerMode) string             // normal|modal|attach|terminal
```

Modal composition draws the base surfaces, applies the dim style, and
overlays the centered box — one pure function of the model, no new state
machines beyond `detailOpen` already present.

### Phase derivation

- Review Runs: `FETCH`,`TRIAGE`,`AGENT`,`VERIFY`,`PUSH` from Run state and
  cycle position (push `[locked]` until no Unresolved Review Issues remain,
  matching Final Push gating).
- Spec Runs: `AGENT`,`VERIFY`,`COMMIT` per current Task plus `QA` (`[locked]`
  until all Tasks completed; omitted when the Run has no QA opt-in — read
  from the Run's journal events, not new state).

### Detail sources

Review: existing artifact detail loading. Spec: task file body via the spec
package, read-only, same mid-write tolerance as the pane (parse failure keeps
the last good detail with a stale marker line).

### API Contracts

None outward. Key map (documented in footer, asserted in tests): `Tab`,
arrows/`j/k`, `PgUp/PgDn`, `Enter` open, `D` toggle, `Esc` close, `End`
follow, `Ctrl-C` stop — attach and terminal modes keep their existing
differences. Non-TTY output paths are untouched by construction (separate
renderer).

## Coverage Map

- Story 1 → work queue renderer + batch separators (mockup 1)
- Story 2 → `runPhases` spec variant + WorkItem queue for Tasks
- Story 3 → modal renderer + key routing + detail sources (finding 5)
- Story 4 → attach wiring over the same renderers
- Story 5 → viewport/Follow untouched, regrouping tests guard it
- Story 6 → responsive fallback in the layout renderer

## Integration Points

None external. Consumes existing store/journal/spec readers only.

## Testing Approach

The plan's test list is adopted wholesale (see
`design/ui-redesign-plan.md`, "Test Plan"), executed as synchronous
`model.Update` tests and render-string assertions — no terminal emulation:
status labels for every Work Item state, batch separators and elapsed
placement, selected-row rendering, modal open/toggle/close/scroll/focus,
footer per state, timeline grouping per event kind, follow behavior after
grouping, small-terminal fallback, and byte-equal before/after context on
modal close. Step 1 captures current-behavior snapshot assertions first so
the refactor is provably rendering-neutral before any visual change lands.

## Build Order

1. **Renderer decomposition, rendering-neutral** — extract per-surface
   helpers behind snapshot assertions that prove unchanged output.
2. **Two-pane base layout + phase row** (depends on: 1).
3. **Work Queue upgrade** — batch separators, elapsed, markers, severity
   when present, totals footer (depends on: 2).
4. **Timeline grouping** by Batch and event kind with Follow Mode guarded
   (depends on: 2).
5. **Modal detail** — renderer, Enter/D/Esc routing, both detail sources,
   dimming (depends on: 2).
6. **Footer states + responsive fallback** (depends on: 3, 4, 5).
7. **Attach parity pass + docs/skill sync** — attach tests over both Run
   kinds; update the Roundfix skill's Live Run View wording; `make
   skills-sync` (depends on: 6).

## Risks & Considerations

- The decomposition step is the risk sink: it must be provably neutral
  (snapshot-first) or every later diff becomes unreviewable.
- Severity display depends on what Review Issue artifacts actually carry —
  render blank when absent rather than inventing a level (verify the field
  during task 3).
- Small-terminal fallback must never panic on zero/negative widths — table
  tests over degenerate sizes.
- Dimming plus modal must stay readable on light terminals — prefer border
  emphasis over heavy shading (Lip Gloss styles only).

## Decisions

- Refactor-before-feature with snapshot neutrality as step 1.
- Enter opens, D toggles, Esc closes (PRD decision, plan's open question
  settled).
- Phase row replaces the old progress header; totals move to the queue
  footer.
- Render-time grouping only; no daemon event enrichment (ADR-0008).
- The `design/` plan and mockups are the binding visual contract.
