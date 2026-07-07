---
task: task_02
spec: 0020-run-browser
status: completed
type: backend
complexity: medium
---

# Task 02: runs list enrichment: columns, state and limit flags, notes

## Overview

Bring the deterministic text surface to the new contract: enriched columns,
an Active-by-default state filter with `--state`, a bounded default with
`--limit`, hidden-count notes, and the self-explaining attach wording for
unknown ids and non-interactive contexts. Demoable by running `runs list`
against seeded Runs.

## Requirements

1. MUST print one Run per line, newest first, with the stable column order:
   run id, state, kind, target, agent, absolute UTC start time, duration
   (`running <elapsed>` for Active Runs), and local branch — plus the
   repository as the final column with `--all`. Run ids are never truncated.
2. MUST default the state filter to Active and support
   `--state <active|terminal|all>`; the unreleased `--active` flag is
   removed.
3. MUST default the bound to the 20 newest matching Runs and support
   `--limit N` with `0` unbounded.
4. MUST print exactly one trailing stderr note when Runs are hidden by the
   state filter or the bound, naming the hidden count and the widening flag
   (shapes: `(N terminal Run(s) hidden; use --state all)`,
   `(N older Run(s) hidden; use --limit 0)`).
5. MUST keep the empty result a single stdout line with exit `0`, stdout
   report-only, and every existing exit-code contract.
6. MUST update `roundfix attach <unknown>` to fail with the error naming
   picker numbers as unstable, and the non-interactive no-run-id failures of
   `attach` (and `runs` without a subcommand) to name `runs list`.

## Subtasks

- [x] Column rendering with the shared row formatter (absolute time, duration)
- [x] `--state` flag replacing `--active`; default active
- [x] `--limit` flag with default 20
- [x] Hidden-count stderr notes
- [x] attach unknown-id and non-interactive wording; `runs` non-TTY wording
- [x] CLI tests: byte-pinned columns, each flag, notes, empty, wording

## Acceptance Criteria

- [x] With seeded Runs of both kinds, `runs list` prints only Active Runs
      with the eight columns and pins the byte shape in a CLI test.
- [x] `--state all` includes terminal Runs; `--state terminal` excludes
      Active ones; hiding produces exactly one matching stderr note.
- [x] 25 seeded matching Runs print 20 lines plus the older-hidden note;
      `--limit 0` prints all 25 without a note.
- [x] `attach 41` (unknown id) exits `2` with the picker-number error;
      no-run-id non-interactive attach and bare `runs` in non-TTY exit `2`
      naming `runs list`.

## Verification

- `rtk go test ./internal/cli/` — expected: all tests pass, including the new
  listing and wording tests.
- `rtk make verify` — expected: fmt-check, tests, skills check, and build all
  pass.

## References

`_prd.md` → User Stories 4-5; Core Features 4-6. `_techspec.md` → API
Contracts: runs list, attach; Build Order 2; Risks (pre-release flag
supersession).

## Result

Text surface brought to the new contract:

- **Shared row formatter** — new `FormatRunRow(run, now, relative, withRepo)`
  in `internal/tui/runrow.go`: run id (never truncated), state, kind, target,
  agent, start time, duration, local branch, plus repository with `withRepo`.
  Absolute UTC RFC 3339 start for the text surface; `relative` renders
  `<elapsed> ago` for the TUI (Task 03 consumer). Durations truncate to
  `53s` / `42m` / `1h12m`; non-terminal Runs render `running <elapsed>`
  against the injected clock; missing fields render `-` so the column count
  stays stable.
- **`runs list`** (`internal/cli/runs.go`) — `--state <active|terminal|all>`
  (default `active`) replaces the unreleased `--active` (now a flag-parse
  error); `--limit N` (default `20`, `0` unbounded) applies after the state
  filter, newest first. One unbounded `StatesAll` query feeds both the
  visible rows and the exact hidden counts (classification through
  `store.IsTerminalState` — no duplicated state lists; no second data path).
  Empty result stays the single `No Runs found.` stdout line with exit `0`;
  stdout carries only the report, notes go to stderr.
- **Hidden-count notes** — exactly one trailing stderr note:
  `(N terminal Run(s) hidden; use --state all)` under the active filter,
  `(N active Run(s) hidden; use --state all)` under the terminal filter,
  `(N older Run(s) hidden; use --limit 0)` when the bound truncates.
  Interpretation settled in code comment: when both the bound and the state
  filter hide Runs, the bound's note wins because it truncates Runs the
  caller asked for.
- **Wording** — `attach <unknown>` fails exit `2` with the techspec string:
  `Run "41" does not exist; picker numbers are not stable Run ids — pass a
  run id or run 'roundfix attach' to pick interactively`. Non-interactive
  no-run-id attach: `missing run id in non-interactive mode; pass a run id
  or run 'roundfix runs list' to discover Runs`. Bare `runs` non-interactive
  (new seam `runsInteractiveInputAvailable`): exit `2`,
  `runs requires a subcommand in non-interactive mode; use 'roundfix runs
  list'`; interactive bare `runs` keeps printing usage until Task 04 wires
  the browser.
- **Help text** (`internal/cli/cli.go`) — root usage line and the `runs`
  usage block document `--state`, `--limit`, the enriched columns, and the
  note behavior.

### Acceptance criteria evidence

All in `internal/cli/cli_test.go` (buffer-captured `Run(...)`, pinned clock
via `withRunsListNow`, pinned `completed_at` via `setListedRunCompletedAt`):

- Eight columns, Active default, byte-pinned:
  `TestRunRunsListPrintsStableColumnsNewestFirst` seeds resolve+implement+
  other-repo Runs and pins the full byte shape
  (`<id>  Active  implement  spec:0017-run-discovery  codex
  2026-07-06T12:02:00Z  running 12m  ma/spec-run`) plus the terminal-hidden
  note.
- State filter and notes: `TestRunRunsListStateFlagFiltersAndNotes` pins
  `--state all` (both rows, no note), `--state terminal` (terminal row only,
  `(1 active Run(s) hidden; use --state all)`), and `--all --state all`
  (repository column) byte-for-byte; stderr compared exactly, so exactly one
  note.
- Bound: `TestRunRunsListLimitBoundsNewestMatches` seeds 25 Active Runs —
  default prints 20 newest plus `(5 older Run(s) hidden; use --limit 0)`;
  `--limit 0` prints all 25 with empty stderr; `--limit 3` prints 3 plus the
  22-hidden note.
- Wording and exits: `TestAttachUnknownRunFailsBeforeTUIStart` runs
  `attach 41` → exit `2` with the picker-number error;
  `TestAttachWithoutRunIDNonInteractiveNamesRunsList` (no TTY and
  `--no-input`) → exit `2` naming `runs list` and the run-id requirement;
  `TestRunRunsWithoutSubcommandHonorsInteractivity` → non-interactive bare
  `runs` exits `2` naming `runs list`, interactive prints usage exit `0`.
  Also `TestRunRunsListAllRowsHiddenKeepsSingleEmptyLine` (empty-after-filter
  stays `No Runs found.` + note, exit `0`) and
  `TestRunRunsListUsageErrors` (`--state bogus`, `--limit -1`, removed
  `--active` → exit `2`).

### Verification evidence

- `rtk go test ./internal/cli/ -run 'TestRunRunsList|TestRunRunsWithoutSubcommand|TestAttachUnknownRunFailsBeforeTUIStart|TestAttachWithoutRunIDNonInteractiveNamesRunsList|TestRunCommandHelp' -v`
  — 33 passed.
- `rtk go test ./internal/cli/ ./internal/tui/ ./internal/store/` — 548
  passed; `rtk gofmt -l` clean.
- `rtk make verify` — exit 0: fmt-check, `go test ./...` (961 passed in 19
  packages), `roundfix skills check` passed, build succeeded.

### Follow-up notes

- Task 03: `FormatRunRow` ships with the `relative` mode rendering
  `<elapsed> ago`; the browser can shorten run ids on its side (text surface
  keeps them full per the agent contract).
- Task 04: interactive bare `runs` still prints usage
  (`TestRunRunsWithoutSubcommandHonorsInteractivity/interactive prints
  usage` pins the interim behavior) and the numbered attach picker
  (`renderAttachRunPicker`, `runListState`, `runListTarget`) is untouched —
  both are replaced by the Run Browser wiring.
- Task 05: README/usage-guide/SKILL.md sync for the new flags and wording
  (`roundfix skills check` passes today; `runs list` flag docs live in the
  usage guide and skill).
