---
spec: 0020-run-browser
prd: _prd.md
created: 2026-07-07
---

# Run Browser — Technical Spec

## Executive Summary

The Run Browser is a new Bubble Tea model over the existing `Store.ListRuns`
query, wired as the interactive face of two entry points (`roundfix runs`,
no-arg `attach`) while the plain-text `runs list` remains the agent contract.
The primary trade-off is accepting one more TUI surface to maintain in
exchange for killing the unbounded numbered prompt: the browser reuses the
cockpit's conventions and the store's one listing query, so the new code is a
list model plus row formatting, not a second data path. Because v0.2.0 has
not shipped, the unreleased `--active` flag is superseded by `--state` and
the default filter flips to Active Runs with no public breakage.

## System Architecture

Three modules extend, none added:

- `internal/store` — `ListRunsQuery` gains `Limit` (SQL `LIMIT`, `0` =
  unbounded) and a `States` filter generalizing `ActiveOnly` (active,
  terminal, or all).
- `internal/tui` — new `RunBrowser` Bubble Tea model (v2 API): rows from
  `store.Run` values, newest first, Active-only default with an `a` toggle,
  `Enter` selects, `q`/`Esc` cancels. Pure model logic; no terminal
  emulation in tests.
- `internal/cli` — `runs` without a subcommand and `attach` without a run id
  open the browser in a TTY and loop: browser → existing attach cockpit →
  back to a refreshed browser until cancel. `runs list` grows the new
  columns, `--state`, `--limit`, and hidden-count notes. Non-interactive
  paths keep exit `2` with updated wording.

Flow: entry point → `ListRuns` (state filter) → browser model → selected Run
→ `runAttachCockpit` (unchanged, read-only) → return to browser (fresh
query) → cancel exits `0`.

## Implementation Design

### Interfaces

Store query extension:

```go
type ListRunsQuery struct {
    GitRoot string
    States  RunStateFilter // StatesActive (default) | StatesTerminal | StatesAll
    Limit   int            // 0 = unbounded
}
```

Browser model (in `internal/tui`):

```go
// RunBrowser is the interactive Run discovery list. It never mutates Runs.
type RunBrowser struct { /* rows, cursor, showAll, size */ }

func NewRunBrowser(repo string, active []store.Run, all []store.Run) RunBrowser
func (b RunBrowser) Update(msg tea.Msg) (RunBrowser, tea.Cmd)
func (b RunBrowser) View() string

// BrowserOutcome reports how the browser closed.
type BrowserOutcome struct {
    RunID    string // empty on cancel
    Cancelled bool
}
```

Row formatting shared by text and TUI:

```go
// FormatRunRow renders one Run's discovery fields. Relative time is
// TUI-only; the text surface always uses absolute UTC.
func FormatRunRow(run store.Run, now time.Time, relative bool, withRepo bool) []string
```

### Data Models

No schema change. Every displayed field exists on the Run row: `Agent`,
`CreatedAt` (start), `CompletedAt` (duration; elapsed against now for Active
Runs), `LocalBranch`, `GitRoot`, plus the existing id/state/kind/target.
Short run id in the TUI is the timestamp-less suffix; the text surface keeps
full ids untruncated (agent contract).

### API Contracts

`roundfix runs list` (text, agents):

- Columns: `<run-id>  <state>  <kind>  <target>  <agent>  <started-utc>
  <duration>  <branch>`, plus the repository as the final column with
  `--all`. Full ids, absolute `2026-07-07T14:05:53Z` timestamps, durations
  like `42m` / `1h12m`, `running <elapsed>` for Active Runs.
- `--state <active|terminal|all>` — default `active`. The unreleased
  `--active` flag is removed.
- `--limit N` — default `20`, `0` unbounded, applied after the state filter.
- When the state filter or bound hides Runs, exactly one stderr note names
  the hidden count and the widening flag, shaped like
  `(23 terminal Run(s) hidden; use --state all)` or
  `(15 older Run(s) hidden; use --limit 0)`.
- Empty result stays one stdout line, exit `0`. Exit codes unchanged.

`roundfix runs` (no subcommand):

- TTY: opens the Run Browser. Non-TTY: exit `2`,
  `runs requires a subcommand in non-interactive mode; use 'roundfix runs
  list'`.

`roundfix attach`:

- No run id + TTY: opens the Run Browser (replaces the numbered prompt).
- No run id + non-TTY: exit `2`, wording updated to name both `runs list`
  and the run-id requirement.
- Unknown run id (for example a picker number): exit `2`,
  `Run "41" does not exist; picker numbers are not stable Run ids — pass a
  run id or run 'roundfix attach' to pick interactively`.
- `attach <known-run-id>`: byte-for-byte unchanged.

Run Browser (TUI):

- Header: repository name and filter state (`ACTIVE` / `ALL`).
- Keys: `↑↓` move, `Enter` attach (opens the existing Live Run View
  read-only; leaving it returns to a refreshed browser), `a` toggle
  active/all, `q`/`Esc`/`Ctrl-C` quit with no side effects, exit `0`.
- Empty Active view names the toggle:
  `No active Runs — press a to show all Runs.`

## Coverage Map

- Story 1 (active-first navigable list) → `RunBrowser`, default
  `StatesActive`, entry-point wiring
- Story 2 (row context) → `FormatRunRow`, Run row fields
- Story 3 (attach and return) → browser→cockpit loop in `internal/cli`
- Story 4 (bounded deterministic text) → `--state`/`--limit`, hidden-count
  notes, `ListRunsQuery.Limit`
- Story 5 (number trap) → attach unknown-id error wording
- Core Feature 3 (read-only) → browser model has no mutation paths

## Integration Points

None external. Bubble Tea v2 / Lip Gloss v2 are already dependencies
(cockpit).

## Testing Approach

Existing seams:

- `internal/store` table tests for `States` and `Limit` combinations.
- `internal/tui` model tests driving `Update` synchronously with
  `tea.KeyPressMsg`: navigation, toggle, selection outcome, cancel outcome,
  empty states, small-size degradation.
- `internal/cli` buffer-captured tests: `runs list` new columns and notes
  byte-pinned, `--state`/`--limit` behavior, non-TTY `runs` exit `2`,
  attach unknown-id wording, no-arg attach non-TTY wording. The browser
  loop is exercised through the model seam, not terminal emulation.

## Build Order

1. Store: `RunStateFilter` and `Limit` on `ListRunsQuery`, migrating the
   `ActiveOnly` call sites, with store tests.
2. Text surface: `runs list` columns, `--state`, `--limit`, hidden-count
   notes, and the attach/`runs` non-interactive and unknown-id wording, with
   CLI tests (depends on: 1).
3. `RunBrowser` model and row formatting in `internal/tui`, with model tests
   (depends on: 1).
4. Entry-point wiring: `roundfix runs` and no-arg `attach` open the browser
   in a TTY, browser→cockpit→browser loop, with CLI tests (depends on:
   2, 3).
5. Docs and skill sync: README, usage guide, roundfix SKILL.md, CONTEXT.md
   Run Browser term, `make skills-sync` (depends on: 4).

## Risks & Considerations

- The `--active` → `--state` supersession is safe only while unreleased;
  this spec must ship in v0.2.0 with 0017 or the flag becomes public API.
- The browser refresh-on-return re-queries the store; with thousands of Runs
  the `Limit` also bounds the browser's initial page — the toggle loads all,
  which is acceptable at local scale.
- Duration for Active Runs uses wall clock against `CreatedAt`; clock skew
  is cosmetic only.

## Decisions

- One shared listing query and one shared row formatter feed both surfaces —
  no second data path.
- `--state` replaces `--active` (pre-release supersession, PRD Decisions).
- Browser refreshes by re-querying on toggle and on return from the cockpit;
  no live auto-refresh (PRD Non-Goals).
