---
task: task_03
spec: 0020-run-browser
status: pending
type: frontend
complexity: high
---

# Task 03: RunBrowser TUI model with row formatting

## Overview

Build the Run Browser as a pure Bubble Tea v2 model in the TUI package: a
navigable, read-only list of Runs with the enriched row fields, Active-only
default, an all/active toggle, selection, and cancel — plus the shared row
formatter both surfaces use. Verifiable in isolation by driving `Update`
synchronously.

## Requirements

1. MUST render one row per Run, newest first: short run id, state, kind,
   target, agent, relative start time, duration (elapsed for Active Runs),
   and local branch. Active rows are visually distinct when the all filter is
   on.
2. MUST default to Active Runs only, with the `a` key toggling between
   Active and all; the header names the repository and the current filter.
3. MUST support `↑`/`↓` movement, `Enter` reporting the selected run id,
   and `q`/`Esc`/`Ctrl-C` reporting cancel — the model itself never attaches,
   mutates, or stops anything.
4. MUST render a filter-aware empty state
   (`No active Runs — press a to show all Runs.`) and degrade to fewer
   columns at small widths before breaking layout.
5. MUST expose a shared row formatter usable by the text surface (absolute
   UTC time) and the browser (relative time), so the two surfaces cannot
   drift.
6. MUST use Bubble Tea v2 module paths and API, with model tests driving
   `Update` synchronously — no terminal emulation.

## Subtasks

- [ ] Browser model: rows, cursor, filter state, outcome
- [ ] Shared row formatter (relative/absolute time, duration, repo column)
- [ ] Toggle, navigation, selection, and cancel key handling
- [ ] Filter-aware empty state and small-width degradation
- [ ] Model tests: navigation, toggle, selection, cancel, empty, width

## Acceptance Criteria

- [ ] Driving the model with key messages selects a Run by `Enter` and
      reports its run id; `q` reports cancel with no selection.
- [ ] The default view contains only Active Runs; after `a`, terminal Runs
      appear and Active rows are distinguishable.
- [ ] The empty Active view renders the toggle invitation.
- [ ] A narrow width renders without overflow, dropping columns in a
      documented order.

## Verification

- `rtk go test ./internal/tui/` — expected: all tests pass, including the new
  browser model tests.
- `rtk make verify` — expected: fmt-check, tests, skills check, and build all
  pass.

## References

`_prd.md` → User Stories 1-2; Core Features 1-3; User Experience.
`_techspec.md` → Interfaces: RunBrowser, FormatRunRow; Build Order 3.
CONTEXT.md → Run Browser.
