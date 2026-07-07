---
task: task_01
spec: 0021-cockpit-visual-fidelity
status: pending
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

- [ ] Token set implementing the approved palette rules
- [ ] Color-mode resolution to styled or identity tokens
- [ ] Fold existing TUI styles into the tokens
- [ ] Bounded event-summary helper
- [ ] Tests: each token under a forced color profile; identity under
      no-color; summary boundedness on multi-line payloads

## Acceptance Criteria

- [ ] Under a forced color profile, each token renders its documented color
      family; with color disabled, output equals the unstyled text.
- [ ] A multi-line, oversized event summary renders exactly one line bounded
      to the given width.
- [ ] No existing TUI test changes behavior except through the new tokens.

## Verification

- `rtk go test ./internal/tui/` — expected: all tests pass, including the new
  token and summary tests.
- `rtk make verify` — expected: fmt-check, tests, skills check, and build all
  pass.

## References

`_prd.md` → Core Features 1, 7; User Story 5. `_techspec.md` → Interfaces:
Tokens, ResolveTokens, EventSummary; Build Order 1.
`../_archived/0005-tui-cockpit/design/ui-redesign-plan.md` → §7 style rules.
