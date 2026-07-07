---
task: task_03
spec: 0020-run-browser
status: completed
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

- [x] Browser model: rows, cursor, filter state, outcome
- [x] Shared row formatter (relative/absolute time, duration, repo column)
- [x] Toggle, navigation, selection, and cancel key handling
- [x] Filter-aware empty state and small-width degradation
- [x] Model tests: navigation, toggle, selection, cancel, empty, width

## Acceptance Criteria

- [x] Driving the model with key messages selects a Run by `Enter` and
      reports its run id; `q` reports cancel with no selection.
- [x] The default view contains only Active Runs; after `a`, terminal Runs
      appear and Active rows are distinguishable.
- [x] The empty Active view renders the toggle invitation.
- [x] A narrow width renders without overflow, dropping columns in a
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

## Result

New `internal/tui/runbrowser.go`: the Run Browser as a pure Bubble Tea v2
value model (`charm.land/bubbletea/v2`, `tea.KeyPressMsg`/`tea.Key`), per the
techspec interface — `NewRunBrowser(repo, active, all)`,
`Update(tea.Msg) (RunBrowser, tea.Cmd)`, `View() string`, `Outcome()
BrowserOutcome`.

- **Rows** — `FormatRunRow(run, now, relative=true, withRepo=false)` (the
  Task 02 shared formatter; requirement 5 satisfied by consumption, so the
  surfaces cannot drift), with the id column swapped to `browserShortRunID`
  (the timestamp-less suffix of `run_<ts>_<suffix>`; other shapes render in
  full). Columns align via per-column max widths; `now` is captured at
  construction (no live auto-refresh, per PRD Non-Goals) and injectable in
  package tests.
- **Filter and header** — Active-only default backed by the pre-queried
  `active` slice; `a` swaps to `all` and back, resetting the cursor. Header:
  `Run Browser — <repo> — ACTIVE|ALL`. Footer: `↑↓ move  enter attach
  a all/active  q quit`.
- **Keys** — `↑`/`k`, `↓`/`j` move with edge clamping; `Enter` records
  `BrowserOutcome{RunID}` and returns `tea.Quit` (ignored on an empty list);
  `q`/`Esc`/`Ctrl-C` record `BrowserOutcome{Cancelled: true}` and quit. The
  model holds only slices and ints — no store handle, so it cannot attach,
  mutate, or stop anything.
- **Empty states** — Active: `No active Runs — press a to show all Runs.`
  (techspec copy); all: `No Runs found.` (matches the text surface).
- **Degradation** — `browserColumnDropOrder` documents the small-width
  order: branch → agent → start → target → kind drop until rows fit; short
  id, state, and duration always survive; every line is width-truncated as
  the last resort. Height windows the row list around the cursor
  (`visibleWindow`), so small terminals never overflow.
- **Distinct Active rows** — under the all filter, non-terminal rows render
  through `cockpitTokens.Running` (amber); the `running <elapsed>` duration
  keeps the distinction text-only readable when color is off (ResolveTokens
  identity contract).

### Acceptance criteria evidence

All in `internal/tui/runbrowser_test.go`, driving `Update` synchronously
with `tea.KeyPressMsg` (same synthetic-key seam as the cockpit tests; no
terminal emulation):

- Selection and cancel: `TestRunBrowserEnterReportsSelectedRunID` (`↓` +
  `Enter` → `tea.QuitMsg` and the selected full run id);
  `TestRunBrowserCancelKeysReportCancel` table over `q`/`esc`/`ctrl+c`
  (Cancelled true, empty RunID).
- Active default and distinct rows:
  `TestRunBrowserDefaultShowsActiveRunsNewestFirst` (only Active Runs, short
  ids, relative start, `running 12m`, header `ACTIVE`, cursor on newest);
  `TestRunBrowserToggleShowsAllAndDistinguishesActive` (after `a`: header
  `ALL`, terminal Run with `42m` visible, Active row line carries ANSI
  styling while the terminal row is byte-identical to its stripped form;
  second `a` returns to Active-only).
- Empty invitation: `TestRunBrowserEmptyStatesNameTheFilter` (exact toggle
  invitation, Enter inert on empty, `No Runs found.` under the all filter).
- Narrow width: `TestRunBrowserNarrowWidthDropsColumnsInOrder` at widths
  120/90/60 asserts every line fits the width and columns disappear in the
  documented order (120 keeps all; 90 drops branch, agent, start; 60 keeps
  id, state, duration only). `TestRunBrowserSmallHeightKeepsCursorVisible`
  windows three rows into a height-6 terminal with the cursor visible.
- Formatter sharing: `TestFormatRunRowSharedByBothSurfaces` pins absolute
  UTC vs `12m ago` from the same call, identical durations, untruncated id
  in absolute mode, and the repository column under `withRepo`.

### Verification evidence

- `rtk go test ./internal/tui/ -run 'TestRunBrowser|TestFormatRunRow' -v` —
  15 passed.
- `rtk go test ./internal/tui/` — 169 passed; `rtk gofmt -l internal/tui`
  clean.
- `rtk make verify` — exit 0: fmt-check, `go test ./...` (976 passed in 19
  packages), `roundfix skills check` passed, build succeeded.

### Follow-up notes

- Task 04 wires the entry points: wrap `RunBrowser` in a `tea.Program`
  (the model exposes no `Init`; the wrapper owns the program loop), re-query
  and rebuild the model on return from the cockpit, and act on
  `Outcome()`. `NewRunBrowser` takes pre-queried active/all slices — the
  entry point supplies both (the techspec's browser `Limit` note applies
  there).
- The browser styles ride the package-level `cockpitTokens` (styled
  default), matching the current renderers; per-surface color-mode wiring
  stays with the cockpit's token-adoption plan noted in `styles.go`.
