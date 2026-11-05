---
task: task_02
spec: 0353-runs-list-reads-as-a-table-and-the-run-browser-draws-at-a-terminal
status: pending
type: backend
complexity: high
---

# Task 02: runs list prints a table at a terminal and keeps its stable line everywhere else

## Overview

`roundfix runs list` prints only the Agents' stable line, about three hundred
characters per Run (`_techspec.md` → Current behavior 6). After this Task it
gains `--format <table|plain>` and `--wide`: at an interactive terminal, with
`--format table` or with `--wide`, it prints a header and aligned columns
with the short id, the repository name when the listing spans repositories,
the local start time and a bare duration, cutting only the target to fit the
terminal width; piped, or with `--format plain`, it prints today's bytes.
This answers "A lista de runs está confusa, talves com informação em excesso"
of the report reported by the maintainer on 2026-11-05
([references](references/2026-11-05-macbook-runs-list-is-confusing-and-attach-shows-nothing.md)).

## Requirements

1. MUST, before changing code, run
   `go test -count=1 -timeout 20m -v -run '^(TestRunRunsList.*|TestRunsListOutsideARepositoryTranscript|TestRunRunsListOutsideRepositoryListsEveryRepository|TestRunsListPrintsNoWarningForARunWhoseCheckoutWasDeleted)$' ./internal/cli`
   and record in the Result that it passes; these pin the stable line.
2. MUST add `--format` and `--wide` to `parseRunsListOptions` with the
   selection, the two usage errors and their exact messages of
   `_techspec.md` → API Contract 6.
3. MUST create `internal/tui/runtable.go` with `ShortRunID`,
   `RunTableOptions` and `RenderRunTable` exactly as API Contract 7 and API
   Contract 8 state, using `RunRepositoryName` (task_01), `truncateDisplay`,
   `padRightDisplay`, `displayWidth` and `formatRunElapsed`.
4. MUST wire the table in `runRunsListCommand` (API Contract 8): add the
   `runsListLocation func() *time.Location` field to `commandDependencies`
   (default `time.Local`), compute the width from `COLUMNS` (40 or more),
   then `terminalColumns` of the unwrapped stdout, then `0`, and print
   `No Runs found.` in both formats when there is no Run Database. Create
   `internal/cli/terminal_columns_unix.go` (`//go:build unix`,
   `unix.IoctlGetWinsize` with `unix.TIOCGWINSZ`) and
   `internal/cli/terminal_columns_other.go` (`//go:build !unix`, returning
   `0, false`). The plain path keeps `printRunsList` unchanged, and every
   stderr note and warning is written as today in both formats.
5. MUST change the `runs` help text in `commandUsage`: the usage line becomes
   `  roundfix runs list [--all] [--state <active|terminal|all>] [--limit N] [--format <table|plain>] [--wide]`;
   the `list` command line becomes
   `  list  Print a table at an interactive terminal, and the stable run id, state,`
   / `        kind, target, agent, start time (UTC), duration and local branch`
   / `        line otherwise or with --format plain`; the option
   `  --format causes: text (default) or json` becomes
   `  --format list: table or plain (default: table at an interactive terminal,`
   / `           plain otherwise); causes: text (default) or json`, followed by
   `  --wide   list: table with the full Run id, kind, branch and worktree`.
6. MUST create `internal/cli/runs_list_table_test.go`, each test calling
   `t.Parallel()` first, with `withRunsListLocation(t, loc)` (a helper in the
   same file, as `withRunsListNow` is built) and these tests:
   - `TestRunsListPlainFormatKeepsTheStableLine`: over
     `seedRunsForListColumns`, for each of the flag sets none, `--all`,
     `--state all` and `--all --state all`, the stdout and stderr of
     `runs list <flags> --format plain` equal those of `runs list <flags>`
     byte for byte, and with no flags the stdout equals exactly the line
     `TestRunRunsListPrintsStableColumnsNewestFirst` pins and the stderr is
     `(1 terminal Run(s) hidden; use --state all)\n`.
   - `TestRunsListTableShowsAlignedColumns`: two repositories named `fiscus`
     and `quaestor` under a temporary directory, marked as Git roots the way
     `withCLIWorkspace` marks its repository; Runs at
     `base = 2026-07-06T12:00:00Z`: `implement` Active in `fiscus` with Spec
     `0063-um-gtin-por-mestre-e-a-pendencia-de-identidade` at `base+2m`,
     `implement` Active in `quaestor` with Spec
     `0063-sequencia-pela-emissao-propria-e-pdf-ao-cliente` at `base+4m`,
     and `resolve` Clean in `fiscus` with PR `123` at `base`, completed at
     `base+42m`; now `base+14m`; location `time.FixedZone("-03", -3*60*60)`.
     `runs list --all --format table` prints exactly these lines, with
     `<B>` and `<A>` the `ShortRunID` of the `quaestor` and `fiscus` Active
     Runs:
     `RUN       REPO      STATE   TARGET                                                AGENT  STARTED (-03)  DURATION`,
     `<B>  quaestor  Active  0063-sequencia-pela-emissao-propria-e-pdf-ao-cliente  codex  09:04          10m`,
     `<A>  fiscus    Active  0063-um-gtin-por-mestre-e-a-pendencia-de-identidade   codex  09:02          12m`;
     stderr is `(1 terminal Run(s) hidden; use --state all)\n`. Subtests:
     `one repository` (`--format table` from inside `fiscus`: no `REPO`
     column), `terminal detection` (`ROUNDFIX_TUI=always` without
     `--format` prints the table), `another day` (now
     `2026-07-08T12:00:00Z`: the `fiscus` cell is `Jul 06 09:02`), and
     `owner identity unproven` (a Run seeded as
     `TestRunRunsListRendersOwnerIdentityUnprovenMarker` seeds it adds the
     `NOTE` column with `owner identity unproven`).
   - `TestRunsListTableFitsTheWidthByCuttingTheTarget`: the same Runs with
     `COLUMNS=80` print exactly
     `RUN       REPO      STATE   TARGET                AGENT  STARTED (-03)  DURATION`,
     `<B>  quaestor  Active  0063-sequencia-pela…  codex  09:04          10m`,
     `<A>  fiscus    Active  0063-um-gtin-por-me…  codex  09:02          12m`;
     no line is wider than 80.
   - `TestRunsListWideShowsEveryColumn`: `runs list --all --wide` with
     `COLUMNS=80` prints a header whose column names are, in order, `RUN`,
     `REPO`, `STATE`, `KIND`, `TARGET`, `AGENT`, `STARTED (-03)`,
     `DURATION`, `BRANCH`, `WORKTREE`; each row carries the full Run id, its
     branch and its `GitRoot`; no cell ends in `…`.
   - `TestRunsListFormatTranscripts`: Surface Transcripts 1, 2 and 3 of
     `_techspec.md`, under a fresh `HOME` and working directory, asserting
     every stdout line, every stderr line and the exit code.
   - `TestRunsHelpNamesTheFormatAndWideFlags`: `runs --help` stdout contains
     the three lines of requirement 5 that carry `--format <table|plain>`,
     `--format list:` and `--wide   list:`.
   Add to `internal/tui/runtable_test.go` (`t.Parallel()` first)
   `TestShortRunIDIsTheFirstEightSuffixDigits`:
   `run_20261105T094412Z_93cee2a36ff33269` gives `93cee2a3`, `run_x_abc`
   gives itself, and `legacy` gives itself.
7. MUST keep every existing `TestRunRunsList*`, `TestRunsList*` and
   `TestReadOnlyRunCommands*` test of `./internal/cli` and every
   `./internal/tui` test passing unchanged.
8. MUST add to `.agents/skills/roundfix/references/runs.md` after the
   paragraph that starts "`runs list` exits `0` for matching and empty
   results.", verbatim apart from line wrapping:
   "`runs list` has two formats. Agents and scripts pass `--format plain`,
   which always prints the stable lines above; without `--format` the stable
   lines are also printed whenever standard output is not an interactive
   terminal. At an interactive terminal, or with `--format table`, it prints
   a table for people: a header and aligned columns `RUN` (the eight-digit
   short id), `REPO` (only when the listing spans repositories), `STATE`,
   `TARGET`, `AGENT`, `STARTED (<zone>)` in the machine's local time, and
   `DURATION`; only `TARGET` is cut with `…` to fit the terminal width.
   `--wide` selects the table with the full Run id and adds `KIND`, `BRANCH`
   and `WORKTREE`; `--wide --format plain` and any other `--format` value
   exit `2`. The stderr notes and `No Runs found.` are the same in both
   formats."
   Run `make skills-sync`, then
   `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions`.
9. MUST, in `docs/user-guide/commands/runs.md`, change the usage line
   `roundfix runs list [--state active|terminal|all] [--limit N] [--all]` to
   `roundfix runs list [--state active|terminal|all] [--limit N] [--all] [--format table|plain] [--wide]`
   and add after "Empty results print `No Runs found.` and exit `0`.",
   verbatim apart from line wrapping: "At an interactive terminal, or with
   `--format table`, `runs list` prints a table instead: a header and the
   columns `RUN` (short id), `REPO` (only when the listing spans
   repositories), `STATE`, `TARGET`, `AGENT`, `STARTED (<zone>)` in local
   time and `DURATION`, with only `TARGET` cut to the terminal width.
   `--wide` shows the full Run id and adds `KIND`, `BRANCH` and `WORKTREE`.
   Piped output, and `--format plain`, keep the stable line, so scripts and
   Agents pass `--format plain`." Then, as the LAST step after the last
   guide edit, re-record the Behavior Surface Record with
   `go test -count=1 -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage`.
10. MUST run sabotages 6 to 8 of `_techspec.md` → Testing Approach on a
    temporary copy and record which named test fails, and run
    `GOOS=windows go build -buildvcs=false ./internal/cli` and
    `GOOS=linux go build -buildvcs=false ./internal/cli`.
11. MUST NOT change the stable line, its notes, `--all`, `--state`, `--limit`,
    `FormatRunRow`, the Run Browser, `internal/config`, `internal/store`,
    `internal/cli/cli_test.go`, any other skill or reference, the Makefile, CI workflows, `go.mod` or
    `go.sum`.

## Subtasks

- [ ] Characterize the stable line.
- [ ] Parse `--format` and `--wide` and select the format.
- [ ] Render the table with the short id, repository, local time and width budget.
- [ ] Wire the zone seam and the terminal width, with the Unix and other files.
- [ ] Update the help text; add the plain, table, width, wide, transcript and help tests; run the sabotages.
- [ ] Document the formats; sync and record the skill; re-record the Behavior Surface Record last.

## Acceptance Criteria

- [ ] Piped `runs list` and `runs list --format plain` print today's bytes for every flag combination tested.
- [ ] At a terminal, with `--format table` or with `--wide`, `runs list` prints the table of API Contract 7, cutting only `TARGET` to the width.
- [ ] The two usage errors exit `2` with their messages, and the help text names both flags.

## Context

- interface: `internal/cli/runs.go`
- interface: `internal/cli/cli.go`
- creates: `internal/cli/terminal_columns_unix.go`
- creates: `internal/cli/terminal_columns_other.go`
- creates: `internal/cli/runs_list_table_test.go`
- creates: `internal/tui/runtable.go`
- creates: `internal/tui/runtable_test.go`
- interface: `internal/tui/runrow.go`
- interface: `.agents/skills/roundfix/references/runs.md`
- interface: `skills/roundfix/references/runs.md`
- interface: `.agents/skills/roundfix/SKILL.md`
- interface: `skills/roundfix/SKILL.md`
- interface: `skills/testdata/owned-skill-versions.json`
- interface: `docs/user-guide/commands/runs.md`
- interface: `docs/references/behavior-surfaces.json`

## Consumers

- The `runs list` stdout, its flags and its help text — consumers: `./internal/cli`
- `roundtui.ShortRunID` and `RenderRunTable` — consumers: `./internal/tui`, `./internal/cli`

## Verification

- `out="$(go test -count=1 -timeout 20m -v -run '^(TestRunsListPlainFormatKeepsTheStableLine|TestRunsListTableShowsAlignedColumns|TestRunsListTableFitsTheWidthByCuttingTheTarget|TestRunsListWideShowsEveryColumn|TestRunsListFormatTranscripts|TestRunsHelpNamesTheFormatAndWideFlags|TestRunRunsList.*|TestRunsList.*|TestReadOnlyRunCommands.*|TestRunCommandHelp)$' ./internal/cli 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestRunsListPlainFormatKeepsTheStableLine TestRunsListTableShowsAlignedColumns TestRunsListTableFitsTheWidthByCuttingTheTarget TestRunsListWideShowsEveryColumn TestRunsListFormatTranscripts TestRunsHelpNamesTheFormatAndWideFlags; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { echo "missing pass: $name" >&2; exit 1; }; done; out="$(go test -count=1 -timeout 20m -v -run '^(TestShortRunIDIsTheFirstEightSuffixDigits|TestRunBrowser.*|TestFormatRunRowSharedByBothSurfaces)$' ./internal/tui 2>&1)" || { printf '%s\n' "$out"; exit 1; }; printf '%s\n' "$out" | grep -q -- '--- PASS: TestShortRunIDIsTheFirstEightSuffixDigits ' || { echo 'missing pass: TestShortRunIDIsTheFirstEightSuffixDigits' >&2; exit 1; }; GOOS=windows go build -buildvcs=false ./internal/cli && GOOS=linux go build -buildvcs=false ./internal/cli || exit 1; for f in .agents/skills/roundfix/references/runs.md skills/roundfix/references/runs.md docs/user-guide/commands/runs.md; do tr -s '[:space:]' ' ' < "$f" | grep -qF -- '--wide' || { echo "missing --wide in $f" >&2; exit 1; }; done; cmp -s .agents/skills/roundfix/references/runs.md skills/roundfix/references/runs.md || { echo 'runs reference mirror differs' >&2; exit 1; }; go test -count=1 -timeout 20m ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' && go test -count=1 -timeout 20m -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$'` — expected: exit 0. Before this Task none of the seven new tests exists, so the first pass line is missing; afterwards the plain line is unchanged, the table, width budget, wide view and usage errors behave as API Contracts 6 to 8 state, both platform files build, the skill and guide state the formats with a current version record and Behavior Surface Record, and the existing `runs list` tests pass.

## References

- `_prd.md` → Goals; User Story 1; User Story 2; User Story 5; Core Feature 1; Core Feature 2; Core Feature 3; Core Feature 7; Success Metric 1; Success Metric 2; Success Metric 3; Success Metric 7
- `_techspec.md` → Current behavior 6; API Contract 6; API Contract 7; API Contract 8; Surface Transcript 1; Surface Transcript 2; Surface Transcript 3; Invariant 1; Invariant 2; Testing Approach; Build Order 2
