---
task: task_01
spec: 0021-cockpit-visual-fidelity
status: completed
type: frontend
complexity: medium
---

# Task 01: Style tokens and color-mode resolution

## Overview

Create the single style token set the whole cockpit renders through —
implementing the approved style rules (cyan section labels and active
borders, green done, amber running/waiting/pending, red locked/failed/
blocking, muted gray timestamps and paths) — with color-mode resolution that
degrades to identity styles. Verifiable in isolation through token tests with
a forced color profile.

## Requirements

1. MUST define the token set in one TUI package file as the only style
   source, and fold the existing ad-hoc styles into it.
2. MUST resolve tokens from the effective color mode, reusing the existing
   `ROUNDFIX_COLOR`/`NO_COLOR` resolution — never reimplementing it.
3. MUST render identity (unstyled) tokens when color is disabled, keeping
   every state distinction carried by the existing text markers.
4. MUST expose the bounded event-summary helper: one line per Run Event —
   first line of its summary, truncated to width — never raw payload.
5. MUST NOT change any layout, key handling, or CLI contract.

## Subtasks

- [x] Token set implementing the approved palette rules
- [x] Color-mode resolution to styled or identity tokens
- [x] Fold existing TUI styles into the tokens
- [x] Bounded event-summary helper
- [x] Tests: each token under a forced color profile; identity under
      no-color; summary boundedness on multi-line payloads

## Acceptance Criteria

- [x] Under a forced color profile, each token renders its documented color
      family; with color disabled, output equals the unstyled text.
- [x] A multi-line, oversized event summary renders exactly one line bounded
      to the given width.
- [x] No existing TUI test changes behavior except through the new tokens.

## Verification

- `rtk go test ./internal/tui/` — expected: all tests pass, including the new
  token and summary tests.
- `rtk make verify` — expected: fmt-check, tests, skills check, and build all
  pass.

## References

`_prd.md` → Core Features 1, 7; User Story 5. `_techspec.md` → Interfaces:
Tokens, ResolveTokens, EventSummary; Build Order 1.
`../_archived/0005-tui-cockpit/design/ui-redesign-plan.md` → §7 style rules.

## Result

Implemented `internal/tui/styles.go` as the single style source: the
`Tokens` struct (SectionLabel/ActiveBorder cyan 39, Done green 78,
Running/Waiting/Pending amber 214, Blocked/Failed/Locked red 203, Muted gray
244, Selection cyan card border), `ResolveTokens(colorEnabled bool)` mapping
the effective color mode to styled or identity tokens (border tokens keep
frame structure without color, since borders are layout), and
`EventSummary(event, width)` rendering the first line of the event's
summary truncated to width, never the payload. The ad-hoc style var block
moved out of `agent_live.go` into `styles.go`; the legacy names that map to
tokens (`styleAccent`, `styleMuted`, `styleError`, `styleActiveBorder`) now
alias the token set with identical palette values, so no render changed.
Callers pass the color mode resolved by the existing `ROUNDFIX_COLOR` /
`NO_COLOR` logic in `internal/cli/style.go`; nothing re-reads env vars.
Tests in `internal/tui/styles_test.go` pin the deterministic Lip Gloss v2
ANSI output (the forced color profile).

### Acceptance criteria evidence

1. **Token colors + no-color identity** —
   `TestResolveTokensStyledRendersDocumentedColors` asserts each token
   renders its documented 256-color sequence (`\x1b[38;5;<index>m`);
   `TestResolveTokensIdentityWithoutColor` asserts identity renders equal
   the unstyled text; `TestResolveTokensBorderTokensKeepStructure` asserts
   border tokens carry cyan when styled, zero ANSI when not, with identical
   structure. All pass in `rtk go test ./internal/tui/`.
2. **Bounded one-line summary** — `TestEventSummaryRendersOneBoundedLine`
   covers a 50-char 3-line summary at width 20 (renders exactly
   `aaa…` 20 cells, one line), CRLF endings, empty-summary-with-payload
   (renders empty, never payload), and zero width. Passes.
3. **No existing test behavior change** — no existing test file was
   modified; `rtk go test ./internal/tui/` reports 133 passed and the full
   gate `rtk make verify` reports 918 passed in 19 packages, skills check
   passed, build passed.

### Verification evidence

- `rtk go test ./internal/tui/` → 133 passed (includes new token and
  summary tests).
- `rtk make verify` → fmt-check clean, 918 tests passed in 19 packages,
  `roundfix skills check` passed, build succeeded.

### Follow-up notes (out of this Task's slice)

- Renderer adoption of the semantic tokens per surface (phase row done →
  green, running → amber, etc.) and the cockpit-side color-mode wiring land
  with tasks 02–04, which restyle each pane through the tokens.
- The remaining chrome styles in `styles.go` (`styleBright`, `styleTool`,
  bar/footer/inactive-frame) await those per-surface tasks.
