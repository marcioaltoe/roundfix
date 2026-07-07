---
spec: 0020-run-browser
status: archived
created: 2026-07-07
surfaces: [cli, docs]
archived: "2026-07-07"
source_slug: 0020-run-browser
---


# Run Browser

Run discovery shipped as a flat, unbounded listing: the Attach picker dumps
every Run the Run Database has ever recorded — 43 lines and growing — as a
numbered prompt, each row showing only id, state, kind, and target. A user
cannot tell which repository a Run belongs to, which Agent ran it, when it
started, or how long it took; the numbers are ephemeral picker positions that
look like stable arguments (`roundfix attach 41` fails confusingly); and the
history drowns the one question that matters at a terminal: what is running
now. The Run Browser replaces the prompt with a navigable TUI over the same
data, defaults every surface to Active Runs, and enriches each row with the
context a human needs to pick.

## Goals

- Opening Run discovery at an interactive terminal answers "what is running
  now" first: Active Runs by default, full history one toggle away.
- Every row carries enough context to pick without cross-referencing: agent,
  start time, duration, and local branch, alongside id, state, kind, and
  target.
- The deterministic `runs list` text contract stays the agent surface, gains
  the same enrichment and containment rules, and never becomes interactive.
- The number-as-argument trap is closed with an error that explains itself.

## User Stories

1. As a user at a terminal, I want `roundfix runs` or `roundfix attach` to
   open a navigable list of Runs — Active first and by default only Active —
   so that I find the Run I care about without scrolling a 43-line prompt.
2. As a user picking a Run, I want each row to show the Agent, when the Run
   started, how long it ran, and the local branch, so that I can tell two
   Runs on the same Spec or PR apart.
3. As a user browsing, I want to select a Run and land in the existing
   read-only Live Run View, and leave it back into the browser, so that
   inspecting several Runs in a row is one session, not repeated commands.
4. As an agent, I want `runs list` to stay deterministic plain text with a
   bounded default, so that discovery output is parseable and does not grow
   without limit.
5. As a user who types `roundfix attach 41`, I want the failure to tell me
   that picker numbers are not Run ids and what to do instead, so that the
   mistake costs one read instead of a debugging detour.

## Core Features

1. `roundfix runs` (no subcommand) and `roundfix attach` (no run id) at an
   interactive terminal open the Run Browser: a Bubble Tea list of the
   repository's Runs, newest first, showing only Active Runs by default with
   a visible toggle to include terminal Runs.
2. Each browser row shows: short run id, state, kind, target, Agent, relative
   start time, duration (elapsed for Active Runs), and local branch. The
   selected row can be attached with Enter, opening the existing Live Run
   View read-only; leaving the view returns to the browser.
3. The browser is read-only run discovery: no stop, gc, or mutation actions.
   Cancelling exits with no side effects.
4. `runs list` stays non-interactive plain text and gains the same row fields
   (with a stable absolute timestamp instead of a relative one), a default
   state filter of Active Runs, a `--state <active|terminal|all>` filter
   (replacing the unreleased `--active`), and a default bound of the 20
   newest matching Runs with `--limit N` (`0` = unbounded). When Runs are
   hidden by the state filter or the bound, one trailing stderr note names
   the count and the flag that widens the view.
5. In a non-interactive context, `roundfix runs` without a subcommand and
   `roundfix attach` without a run id fail with exit `2` naming `runs list`
   and the run-id requirement — unchanged semantics, updated wording.
6. `roundfix attach <value>` with a value that is not a known run id fails
   with an error stating that picker numbers are not stable Run ids and
   pointing at `roundfix attach` (picker) and `roundfix runs list`.

## User Experience

- Browser layout follows the existing cockpit conventions: header naming the
  repository and filter state, one row per Run, footer keys (`↑↓` move,
  `Enter` attach, `a` toggle all/active, `q` quit). Small terminals degrade
  to fewer columns before breaking layout.
- Active Runs are visually distinct from terminal Runs when the all filter is
  on.
- The empty state names the filter: no Active Runs prints an invitation to
  toggle history instead of a bare empty screen.

## Non-Goals / Out of Scope

- Mutating actions from the browser (stop, gc, settle).
- Cross-repository browsing in the TUI — the browser serves the current
  repository; `runs list --all` keeps covering the rest.
- JSON output for `runs list`.
- Live auto-refresh of the browser list; it reflects the Run Database at open
  and on filter toggle.

## Success Metrics

- Finding and attaching to the Active Run in a repository with 40+ historical
  Runs takes one command and one keypress, with no scrolling past history.
- An agent's `runs list` call returns a bounded, stable-format report whose
  hidden-count note makes truncation explicit.

## Decisions

- One browser, two entry points: `roundfix runs` and no-arg `attach` open the
  same TUI; `runs list` stays the deterministic agent surface.
- Default state filter is Active everywhere; history is opt-in (`a` toggle in
  the TUI, `--state` in text). Decided pre-release, so the unreleased 0017
  `--active` flag is superseded rather than kept as an alias.
- Rows show agent, relative start, duration, and local branch (user
  selection); the text surface uses absolute UTC timestamps for stability.
- Picker numbers stay picker-only; `attach <number>` gets a self-explaining
  error instead of index support, because positions change between
  invocations.

## Open Questions

None.
