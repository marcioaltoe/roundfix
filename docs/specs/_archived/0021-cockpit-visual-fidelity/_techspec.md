---
spec: 0021-cockpit-visual-fidelity
prd: _prd.md
created: 2026-07-07
---

# Cockpit Visual Fidelity — Technical Spec

## Executive Summary

The cockpit's renderers already exist and are rendering-neutral (0005
task_01); this spec threads one style token set through them and tightens
event rendering to summary-first. The primary trade-off is test strategy:
styled output must become assertable, so TUI tests force a fixed Lip Gloss
color profile and pin ANSI-bearing renders for key states — heavier fixtures
than plain-text assertions, accepted because unpinned styling is exactly how
the approved design degraded after 0005. No layout, key, event-storage, or
CLI contract changes.

## System Architecture

One module changes:

- `internal/tui` — a new `styles.go` (token set + color-mode resolution)
  consumed by the existing cockpit renderers (`cockpit.go`, `agent_live.go`
  styles fold into the tokens). Event summary extraction becomes a bounded
  render-time helper. No new packages; the Run Browser (0020) will import
  the same tokens later.

`internal/cli` is untouched except where it already forwards the color mode;
the existing `ROUNDFIX_COLOR` / `NO_COLOR` resolution is reused, not
reimplemented.

## Implementation Design

### Interfaces

```go
// Tokens is the cockpit style contract (ui-redesign-plan §7).
type Tokens struct {
    SectionLabel, ActiveBorder          lipgloss.Style // cyan
    Done                                lipgloss.Style // green
    Running, Waiting, Pending           lipgloss.Style // amber
    Blocked, Failed, Locked             lipgloss.Style // red
    Muted                               lipgloss.Style // gray: timestamps, paths
    Selection                           lipgloss.Style // accent card border
}

// ResolveTokens maps the effective color mode to styled or no-op tokens.
func ResolveTokens(colorEnabled bool) Tokens

// EventSummary returns one bounded line for a Run Event: first line of the
// summary, truncated to width, never raw payload.
func EventSummary(event runevent.RunEvent, width int) string
```

### Data Models

None changed. Batch elapsed stamps derive from existing Run Event
timestamps (first to last event of the Batch, or now for the executing
Batch). Group collapse state derives from Batch status already present in
the timeline model.

### API Contracts

- Rendering only. Stdout/stderr contracts of every command are unchanged;
  the cockpit is stderr/TTY surface.
- `ROUNDFIX_COLOR=never` / `NO_COLOR`: tokens become identity styles; all
  information carried by the existing text markers (`[done]`, `[run]`,
  `[wait]`, `[locked]`, `▼/▶`, severity words).
- Kind labels and their colors: PLAN cyan, TOOL blue/cyan family, THINK
  muted, SESSION cyan, Daemon milestones per state color. Exact palette
  indices live in `styles.go` as the single source.

## Coverage Map

- Story 1 (state by color) → `Tokens`, phase row + header renderers
- Story 2 (carded queue, selection) → Work Queue renderer + `Selection`
  token, Batch elapsed stamps
- Story 3 (grouped summary timeline) → timeline renderer, `EventSummary`,
  collapse rules, `Live · detail` indicator
- Story 4 (empty states) → pane empty-state renderers
- Story 5 (no-color degrade) → `ResolveTokens(false)` path + fallback tests
- Core Feature 7 (pinned styled render) → forced-profile TUI tests

## Integration Points

None external. Lip Gloss v2 already in the tree.

## Testing Approach

Existing seam — synchronous `model.Update`/`View` — extended with a forced
color profile so `View()` output carries deterministic ANSI:

- Token tests: each token renders the documented color; no-color mode
  renders identity.
- Renderer tests per surface: phase row states, queue card + selected card
  (ANSI-pinned), batch stamp, timeline group expanded/collapsed, gutter
  alignment, summary truncation (a multi-line payload renders exactly one
  bounded line), modal frame, empty states, `Live · detail` indicator.
- Fallback twins: the same fixtures with color disabled assert the
  marker-only render.

## Build Order

1. Style tokens and color-mode resolution (`styles.go`), folding the
   existing `agent_live.go` styles in, with token tests.
2. Header, Phase Row, and Work Queue fidelity: chips, cards, severity
   colors, selection border, Batch elapsed stamps, totals footer
   (depends on: 1).
3. Timeline fidelity: group markers and collapse rules, timestamp gutter,
   kind labels, `EventSummary` boundedness, `Live · detail` indicator
   (depends on: 1).
4. Detail Modal fidelity and pane empty states (depends on: 1).
5. Docs and skill sync: README/usage/SKILL touch-ups where the cockpit is
   described, `make skills-sync` (depends on: 2, 3, 4).

## Risks & Considerations

- ANSI-pinned tests are brittle against palette tweaks; mitigation: pin
  through the token names (assert the render contains the token's rendering
  of the word) rather than hard-coded escape literals scattered per test.
- Summary extraction must never panic on malformed payloads — first line of
  `Summary`, already bounded by `runevent.BoundSummary`.
- The executing-expanded/settled-collapsed rule changes how much history is
  visible at once; the Detail Modal and scrollback keep full access.

## Decisions

- One `styles.go` token file as the single style source; 0020's Run Browser
  imports it (sequenced after this spec).
- Tests force a fixed color profile — the enforcement the 0005 contract
  lacked.
- No interactive collapsing; state-driven only (PRD Decisions).
