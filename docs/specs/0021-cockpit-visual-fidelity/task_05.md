---
task: task_05
spec: 0021-cockpit-visual-fidelity
status: completed
type: docs
complexity: low
---

# Task 05: Docs and skill sync for the cockpit visual contract

## Overview

Close the SKILL-matches-CLI gate for the visual contract: update the Live Run
View descriptions where they name concrete cockpit behavior (group markers,
summary rows, indicator, empty states), and re-sync the embedded skill
bundle. Verifiable by the drift check and reading docs against the shipped
render.

## Requirements

1. MUST update the README Live Run View section and the roundfix SKILL.md
   where the cockpit's concrete rendering is described: automatic group
   collapse, summary-only event rows with full content in the Detail Modal,
   the `Live · detail` indicator, and explanatory empty states.
2. MUST note the color behavior: state colors in capable terminals, full
   meaning preserved under `ROUNDFIX_COLOR=never`/`NO_COLOR`.
3. MUST re-sync the embedded bundle so the drift check passes.

## Subtasks

- [x] README Live Run View updates
- [x] roundfix SKILL.md updates and `make skills-sync`
- [x] Drift and skills checks pass

## Acceptance Criteria

- [x] README and SKILL.md describe the shipped rendering truthfully —
      collapse rules, bounded summaries, indicator, empty states, color
      degrade.
- [x] `make skills-sync-check` reports no drift and `roundfix skills check`
      passes.

## Verification

- `rtk make verify` — expected: fmt-check, tests, skills-sync-check, skills
  check, and build all pass.

## References

`_prd.md` → Goals; User Experience. `_techspec.md` → Build Order 5.
CLAUDE.md SKILL.md-matches-CLI HARD RULE.

## Result

Status: completed — 2026-07-07.

### What changed

- README Commands section: the previous README had no Live Run View passage
  at all — cockpit rendering was undocumented. Added one after the `attach`
  description covering the shipped render: `WORK QUEUE` /
  `SESSION.TIMELINE` panes, automatic state-driven Batch collapse
  (`completed`/`failed`/`stopped` fold to one `▶` summary row, others
  expand under `▼`, no key toggles it), one bounded summary row per
  structured event behind the aligned timestamp gutter with raw payloads
  never inline and full content in the Detail Modal, the
  `Live · detail hidden` / `Live · detail open` pane-header indicator,
  per-Run-kind explanatory empty states, and the color contract — cyan
  section labels and active borders, green done, amber running or waiting,
  red locked or failed, muted gray timestamps and paths, with
  `ROUNDFIX_COLOR=never`/`NO_COLOR` keeping every distinction through the
  same layout and text markers.
- `.agents/skills/roundfix/SKILL.md` Live Run View section: added five
  bullets for the same concrete rendering facts (collapse rules, bounded
  summary rows with full content in the Detail Modal, the indicator,
  explanatory empty states, and state colors with the no-color degrade).
- Ran `make skills-sync`; the embedded `skills/roundfix/SKILL.md` now
  matches the canonical source.
- Every documented behavior was transcribed from the shipped code:
  `internal/tui/timeline_group.go` (markers, `timelineBatchSettled`,
  gutter), `internal/tui/styles.go` (`Tokens`, `ResolveTokens`,
  `EventSummary`), `internal/tui/cockpit.go` (`timelinePaneHeaderLine`,
  `workQueueEmptyCopy`, `timelineEmptyCopy`), and
  `internal/cli/style.go` (`ROUNDFIX_COLOR`/`NO_COLOR` resolution).

### Commands run

- `rtk make skills-sync` — exit 0; embedded bundle regenerated.
- `rtk make skills-sync-check` — exit 0; no drift.
- `rtk make verify` — exit 0: fmt-check, `go test ./...` (939 passed in 19
  packages), skills-sync-check, `roundfix skills check` (all 14 bundled
  skills passed), and build.

### Evidence per acceptance criterion

1. Truthful rendering description: README (Commands section, after the
   `attach` paragraph) and both SKILL.md copies now state the collapse rule
   exactly as `timelineBatchSettled` implements it, the summary-row bound as
   `EventSummary` implements it, the indicator strings byte-for-byte as
   `timelinePaneHeaderLine` renders them, the Fetch Run empty-state meaning
   as `workQueueEmptyCopy`/`timelineEmptyCopy` ship it, and the token
   palette and no-color identity-token degrade as `ResolveTokens` defines
   them.
2. Drift and skills checks: `rtk make skills-sync-check` exited 0 after the
   sync, and `roundfix skills check` inside `rtk make verify` reported
   "Roundfix skill check passed" for all 14 bundled skills.

### Notes and follow-ups

- The user-scoped skill copy at `~/.claude/skills/roundfix` predates this
  change; it refreshes through the normal `skills install` flow, outside
  this Task's slice.
