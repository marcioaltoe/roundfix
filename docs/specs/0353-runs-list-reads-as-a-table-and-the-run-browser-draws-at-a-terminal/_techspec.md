---
spec: 0353-runs-list-reads-as-a-table-and-the-run-browser-draws-at-a-terminal
prd: _prd.md
created: 2026-11-05
---

# runs list reads as a table, and the Run Browser draws, at a terminal — Technical Spec

## Executive Summary

Three slices. First, the three terminal programs (the Run Browser, the attach
cockpit and the owning cockpit) receive the terminal file instead of the
command's environment wrapper, a terminal program refuses any writer that is
not a file, the Run Browser reports which key closed it, bare `attach` and
bare `runs` print one closing line, and the Run Browser's `REPO` column names
the repository. Second, `runs list` gains `--format <table|plain>` and
`--wide`, and prints a table rendered in `internal/tui` when its standard
output is a terminal. Third, the store resolves a short Run id, `attach` and
`runs show` use it, and ADR-0363 records the dual surface and the short id.
The trade-off this design accepts is a second text shape for one command: an
Agent whose shell is a pseudo-terminal now gets the table unless it passes
`--format plain`, which the skill tells it to do, in exchange for a listing a
person can read without a flag.

## Project Constraints

- Identifier strategy: applicable — no identifier is minted; the short id is
  the first eight random hexadecimal digits of an existing Run id, resolved
  back to that id by one store query and never stored. Source:
  `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — no request path, endpoint or
  credential changes; the commands read the Run Database and the terminal.
  Source: `docs/agents/go.md`, `docs/agents/cli.md`.
- Active ADR obligations: applicable — ADR-0357: "Outside a Git repository `runs list` behaves as `runs list --all`, with the
  repository column and exit `0`"; the table follows the same rule for its
  `REPO` column, and `attach` and `runs show` stay Read-Only Run Commands.
  ADR-0362 (the Git-root fix, delivered first) changes only which `.git`
  entry marks a Git root; the listing, the table and the Run Browser carry no
  obligation from it. ADR-0363, written by task_03, records the decision below. Source:
  `docs/agents/domain.md`, `docs/agents/cli.md`.
- Tooling authority: applicable — express maintainer authorization: the cycle
  AZ answers of 2026-11-05 in the authorization record; bounded files are the
  Roundfix skill files the Tasks change. No Makefile, workflow, `go.mod`,
  `go.sum`, lint, formatter or Baseline module change; `golang.org/x/sys` is
  already a direct requirement. Source: `docs/agents/agent-instructions.md`,
  `docs/agents/spec-routing.md`; Spec-contained authorization record:
  `docs/specs/0353-runs-list-reads-as-a-table-and-the-run-browser-draws-at-a-terminal/_authorization.md`;
  bounded files: `.agents/skills/roundfix/SKILL.md`, `.agents/skills/roundfix/references/runs.md`, `skills/roundfix/SKILL.md`.

## Current behavior

Read at `38cc5854`.

1. `Run` (`internal/cli/cli.go`) wraps `stdout` and `stderr` with
   `commandWriter`, a `commandEnvironmentWriter` struct embedding the writer.
   `defaultRunBrowserSession` (`internal/cli/runbrowser.go`),
   `runAttachCockpit` (`internal/cli/attach.go`) and `startRunUI`
   (`internal/cli/runui.go`) pass that wrapper to `roundtui.RunBrowserSession`
   and `roundtui.RunCockpit`, which give it to `tea.NewProgram` through
   `tea.WithOutput`. Bubble Tea measures the screen only through a writer
   that is a terminal file, so with the wrapper the program has no size and
   draws nothing; `colorEnabled` and `liveTUIEnabled` already unwrap with
   `environmentForWriter`.
2. Probe (pseudo-terminal sized 120 by 40, roundfix built at `38cc5854`):
   bare `runs` wrote only the alternate-screen, keyboard-mode and clear
   sequences for 2.5 seconds; a `SIGWINCH` after a resize drew nothing; `q`
   and `Esc` restored the screen and exited `0` with no text. A program that
   called `RunBrowserSession` with `os.Stdout` drew the header, the empty
   state and the footer; with `struct{ io.Writer }{os.Stdout}` it drew
   nothing. A lone escape byte followed 200 ms later by `[?2026;2$y` closed
   the browser as `Esc` and left the rest on the shell. The same prototype
   binary with `defaultRunBrowserSession` passing the unwrapped writer drew
   the Run Browser in the same pseudo-terminal and, after `q`, printed the
   closing line.
3. `runRunBrowserLoop` returns `exitOK` and prints nothing when the outcome is
   `Cancelled`; `runRunsBrowserCommand` ignores the outcome of its empty pass.
   `RunBrowserSession` turns an end without a key into `Cancelled`.
4. Bare `attach` and bare `runs` refuse with exit `2` when standard input is
   not a character device (`defaultAttachInteractiveInputAvailable`) or
   `liveTUIEnabled(stdout)` is false;
   `TestAttachWithoutRunIDNonInteractiveNamesAllRunsList` and
   `TestRunRunsWithoutSubcommandHonorsInteractivity` pin it.
5. The Run Browser's repository column is `browserRepoName(run.GitRoot)`;
   for a queued Run `GitRoot` is the Run Worktree, so it shows
   `roundfix-deliver-…`. `store.Run.RepositoryRoot` holds the repository's
   main checkout.
6. `runs list` (`internal/cli/runs.go`) has `--all`, `--state`, `--limit`;
   `printRunsList` joins `roundtui.FormatRunRow` fields with two spaces and
   appends `owner_identity_unproven=true`; no header. The notes and
   `No Runs found.` are fixed strings; `TestRunRunsListPrintsStableColumnsNewestFirst`
   and its neighbours pin the bytes.
7. `store.(*Store).Run` selects `WHERE id = ?`; `attach`, `runs show`,
   `events`, `stop` and `reconcile` look Runs up only by the exact id.
   `newRunID` writes `run_<YYYYMMDDTHHMMSSZ>_<16 lowercase hex>`.

## System Architecture

- `internal/cli`: a `terminalOutput` helper that unwraps the command writer
  for the three terminal programs; the closing line in the Run Browser loop;
  `runs list` format selection, a width budget and a local-zone seam; Run id
  resolution in `attach` and `runs show`.
- `internal/tui`: the writer guard in `RunBrowserSession` and `RunCockpit`;
  `BrowserOutcome.Reason`; `RunRepositoryName`, shared by the Run Browser and
  the table; `ShortRunID` and `RenderRunTable`.
- `internal/store`: `ResolveRunReference` and `AmbiguousRunReferenceError`.
- Documents: the Roundfix skill's `runs` reference, the `runs` and `attach`
  user guides, the help text and ADR-0363.

## Implementation Design

### Interfaces

```go
// internal/cli
func terminalOutput(writer io.Writer) io.Writer // environmentForWriter's writer
func browserClosedLine(outcome roundtui.BrowserOutcome) string
func terminalColumns(file *os.File) (int, bool) // unix: TIOCGWINSZ; other: 0, false

// internal/tui
type BrowserOutcome struct { RunID string; Cancelled bool; Reason string }
func RunRepositoryName(run store.Run) string
func ShortRunID(id string) string
type RunTableOptions struct {
	Now      time.Time
	Location *time.Location
	Repo     bool // the listing spans repositories
	Wide     bool
	Width    int  // 0: no budget
}
func RenderRunTable(runs []store.Run, options RunTableOptions) string

// internal/store
type AmbiguousRunReferenceError struct { Reference string; RunIDs []string }
func (store *Store) ResolveRunReference(ctx context.Context, reference string) (Run, bool, error)
```

### Data Models

None. No schema, column or record changes; the short id is computed.

### API Contracts

1. API Contract: the terminal output. `terminalOutput(w)` returns the writer
   `environmentForWriter(w)` returns. `defaultRunBrowserSession`,
   `runAttachCockpit` and `startRunUI` pass `terminalOutput(<writer>)` to the
   terminal program; colors and the TUI test keep using the wrapper.
2. API Contract: the writer guard. `RunBrowserSession` returns, before
   building a program, `fmt.Errorf("run browser session: output %T is not a terminal file", output)`
   when `output` is not an `*os.File`; `RunCockpit` returns
   `fmt.Errorf("run cockpit: output %T is not a terminal file", output)`
   before building its model. A regular file or pipe passes the guard.
3. API Contract: the closing reason. `handleKey` sets `Reason` to the key
   string (`q`, `esc`, `ctrl+c`) on cancel; `RunBrowserSession` sets
   `Cancelled: true, Reason: "interrupted"` when the program ended without a
   key outcome.
4. API Contract: the closing line. After every `Cancelled` outcome,
   `runRunBrowserLoop` and the empty pass of `runRunsBrowserCommand` write
   `browserClosedLine(outcome)` and a newline to stderr, then return `exitOK`.
   The line is `roundfix: Run Browser closed (<reason>); no Run selected.`,
   or `roundfix: Run Browser closed; no Run selected.` when `Reason` is
   empty. Stdout gets nothing.
5. API Contract: the repository name. `RunRepositoryName(run)` is
   `filepath.Base` of the trimmed `RepositoryRoot`, else of the trimmed
   `GitRoot`, else `-`. The Run Browser's `REPO` column uses it.
6. API Contract: `runs list` flags. `--format` accepts `table` or `plain`,
   default empty; any other value is
   `validationError{"unknown --format \"<v>\"; use table or plain"}`;
   `--wide` with `--format plain` is
   `validationError{"--wide applies only to --format table"}`. Both print
   through `printRunsListFailure` and exit `2`. The table is selected when
   `--format table`, or when `--format` is empty and `--wide` is set or
   `liveTUIEnabled(stdout)` is true. Usage line:
   `roundfix runs list [--all] [--state <active|terminal|all>] [--limit N] [--format <table|plain>] [--wide]`.
7. API Contract: the table. `RenderRunTable` returns `No Runs found.\n` for no
   Runs. Otherwise a header line and one line per Run, cells left-aligned to
   each column's widest display width, joined by two spaces, trailing spaces
   trimmed, each line ending in `\n`. Columns: `RUN` (`ShortRunID`; full id
   when `Wide`), `REPO` (when `Repo`; `RunRepositoryName`), `STATE`, `KIND`
   (when `Wide`), `TARGET` (the Spec slug for `implement`, `pr:<n>` for the
   others, `-` when empty), `AGENT` (`-` when empty), `STARTED (<zone>)` with
   `<zone>` = `Now.In(Location).Format("MST")` and cells
   `CreatedAt.In(Location)` as `15:04` on `Now`'s local date and
   `Jan 02 15:04` otherwise, `DURATION` (the elapsed or completion span in
   `formatRunElapsed` form without `running `, `-` when a terminal Run has no
   completion), `BRANCH` and `WORKTREE` (`GitRoot`) when `Wide`, and `NOTE`
   (`owner identity unproven` or empty) only when a listed Run has an
   unproven owner identity. When `Width > 0`, not `Wide`, and a line exceeds
   `Width`, every `TARGET` cell is cut with `truncateDisplay` to
   `max(12, <TARGET width> - <excess>)`.
8. API Contract: `runs list` table wiring. `ShortRunID(id)` is the first eight
   characters of the text after the last `_` of a `run_`-prefixed id with at
   least eight such characters, else the id. `runRunsListCommand` renders the
   table with `Now = runsListNow()`, `Location = runsListLocation()` (a new
   `commandDependencies` field, default `time.Local`), `Repo = opts.all`, and
   `Width` = `COLUMNS` when it parses to 40 or more, else
   `terminalColumns` of the unwrapped stdout when it is an `*os.File`, else
   `0`. With no Run Database it prints `No Runs found.` in both formats. The
   stderr notes and warnings are unchanged.
9. API Contract: short id resolution. `ResolveRunReference` trims the
   reference; an empty one is the error `Run ID is required`. It returns the
   Run whose id equals it when one exists. Otherwise, when the reference
   matches `^[0-9a-f]{8,16}$`, it selects the Runs whose id matches
   `run\_%\_<reference>%` (`LIKE … ESCAPE '\'`) ordered by id descending, which is newest first because the id starts with its creation time: one gives that
   Run and `true`; two or more give `*AmbiguousRunReferenceError`, whose
   `Error()` is
   `Run id prefix "<reference>" matches <n> Runs: <id>, <id>; pass more digits or the full Run id`;
   none gives `false`. Any other reference gives `false`.
10. API Contract: callers. `runAttachCommand` (explicit id) and
    `runRunsShowCommand` call `ResolveRunReference` instead of `Run`; an
    ambiguity prints through each command's failure shape
    (`roundfix attach failed: <error>`;
    `roundfix: runs show failed: <error>` and its usage line) and exits `2`;
    not found keeps each message. `runs show` reads token usage with the
    resolved `run.ID`. `events`, `stop`, `reconcile` and the Run Browser
    keep `Run`.

### Surface Transcripts

1. Surface Transcript: an unknown format is a usage error.

   ```transcript
   $ roundfix runs list --format yaml
   stdout:
   stderr:
   roundfix: runs list failed: unknown --format "yaml"; use table or plain
   Run 'roundfix runs list --help' for usage.
   exit: 2
   ```

2. Surface Transcript: `--wide` needs the table.

   ```transcript
   $ roundfix runs list --wide --format plain
   stdout:
   stderr:
   roundfix: runs list failed: --wide applies only to --format table
   Run 'roundfix runs list --help' for usage.
   exit: 2
   ```

3. Surface Transcript: an empty table, with no Run Database under a fresh
   `HOME`.

   ```transcript
   $ roundfix runs list --format table
   stdout:
   No Runs found.
   stderr:
   exit: 0
   ```

## Invariants

1. Without `--format` and `--wide`, a `runs list` whose stdout is not an
   interactive terminal prints the bytes `38cc5854` prints; `--format plain`
   prints the same bytes.
2. Every stderr note, warning and exit code of `runs list` is the same in
   both formats.
3. A terminal program never receives the command environment wrapper, and
   refuses any non-file writer with an error naming its type.
4. An end of the Run Browser without a selection prints exactly one line on
   stderr and nothing on stdout, and exits `0`.
5. A short id resolves to at most one Run; an exact id is matched before any
   prefix; a reference of fewer than eight digits never matches a prefix.
6. Bare `attach` and bare `runs` with a non-terminal standard input keep
   exiting `2` with today's messages.

## Declared breaks

- `assertRunReaderBrowser` (`internal/cli/read_only_run_commands_test.go`)
  expects stderr equal to the warning; after task_01 it expects the warning
  followed by `roundfix: Run Browser closed; no Run selected.\n`, because its
  scripted outcome has no reason. Its callers
  `TestReadOnlyRunCommandsKeepTheirOutputInsideARepository` (four browser
  subtests) and `TestReadOnlyRunCommandsWarnOnAStrayGitDirectory` (two
  browser subtests per marker) fail until it changes; a prototype on a
  scratch copy of `38cc5854` that only unwrapped the writer and printed that
  line failed exactly those subtests among every `TestReadOnlyRunCommands`,
  `TestAttach`, `TestRunRuns` and `TestBrowserAttachCockpit` test.

No other test is expected to change. Tests that assert empty stdout after a
cancelled browser keep passing, since the line goes to stderr; the help test's
`runs list` usage prefix still matches.

## ADR-0363

task_03 creates `docs/adr/0363-runs-list-prints-a-table-at-a-terminal-and-a-short-run-id-names-a-run.md`
with exactly this content, apart from line wrapping:

```markdown
---
status: accepted
created_at: 2026-11-05T12:00:00Z
updated_at: 2026-11-05T12:00:00Z
deprecated_at: null
superseded_by: null
---

# runs list prints a table at a terminal, and a short Run id names a Run

`roundfix runs list` prints an aligned table with a header when its standard
output is an interactive terminal, and keeps its stable line for every other
reader. The table's eight-digit short id is accepted by `attach` and
`runs show` as a unique prefix of a Run id's random suffix.

## Context

`runs list` was written for Agents: one line per Run with the full Run id,
the UTC start, the branch and, with `--all` or outside a repository, the Run
Worktree path, and no header. On 2026-11-05 the maintainer ran it on his
MacBook outside a repository and wrote: "A lista de runs está confusa, talves
com informação em excesso". The same morning he chose:

> No terminal interativo, runs list mostra cabeçalho, colunas alinhadas, id
> curto, nome do repositório, Spec, hora em BRT e duração. Branch e worktree só
> com --wide. Em pipe ou com --format plain, continua a linha estável de hoje
> para os agentes.

Outside a Git repository `runs list` lists every repository
([ADR-0357](0357-a-read-only-run-command-needs-no-repository-and-only-warns-on-worktree-location.md)).

## Decision

`--format plain` prints the stable line, byte for byte. `--format table`
prints the table. Without `--format`, the table is printed when `--wide` is
given or standard output is an interactive terminal by the test the Live Run
View uses, and the stable line otherwise. `--wide --format plain` is a usage
error.

The table has a header and the columns `RUN`, `REPO`, `STATE`, `TARGET`,
`AGENT`, `STARTED (<zone>)` and `DURATION`. `REPO` is the repository's base
name and appears only when the listing spans repositories. `STARTED` is the
machine's local time, its zone named once in the header. `--wide` shows the
full Run id and adds `KIND`, `BRANCH` and `WORKTREE`. Only `TARGET` is cut,
with `…`, to fit the terminal width. The stderr notes and `No Runs found.`
are unchanged.

The short id is the first eight hexadecimal digits of a Run id's random
suffix. `attach` and `runs show` match the exact id first; a reference of
eight to sixteen lowercase hexadecimal digits names the one Run whose suffix
starts with it, and two or more matches are a usage error listing the full
ids. Every other command keeps requiring the exact id.

## Considered Options

- One format for everyone: rejected, because Agents parse the stable line and
  a person reads the terminal, and neither shape serves the other.
- The table only behind a flag: rejected, because the person at the terminal
  is the one who would not know the flag.
- The full Run id in the table: rejected, because it is the widest column and
  the maintainer chose the short id; accepting the prefix keeps it usable.
- The prefix in every command that takes a Run id: rejected for now, because
  `stop` and `reconcile` act on the Run and `events` serves automation that
  holds the full id.
- A hard-coded America/Sao_Paulo zone: rejected, because the zone is the
  machine's, and the maintainer's machines already run in it.

## Consequences

A person at a terminal reads a short table. An Agent whose shell is a
pseudo-terminal would get the table too, so the Roundfix skill tells Agents
to pass `--format plain`. A short id can become ambiguous as Runs accumulate;
the refusal names every match.
```

## Coverage Map

- Goal 1, Story 1, Story 2, Core Features 1–2, Success Metrics 1–2 →
  `runs list` format selection, `RenderRunTable`, `runsListLocation`,
  `terminalColumns` (task_02).
- Story 5, Core Feature 3, Success Metric 3 → the wide table and the usage
  errors (task_02).
- Goal 2, Story 3, Core Feature 4, Success Metric 4 → `ResolveRunReference`,
  `attach` and `runs show` (task_03).
- Goal 3, Story 4, Core Features 5–6, Success Metrics 5–6 → `terminalOutput`,
  the writer guard, `BrowserOutcome.Reason`, `browserClosedLine`,
  `RunRepositoryName` (task_01).
- Goal 4, Core Feature 7, Success Metric 7 → the paragraphs each Task writes
  into the Roundfix `runs` reference and its mirror and into the `runs` and
  `attach` guides, which name `--format plain` for Agents, the table, `--wide`,
  the short id and the closing line (task_01, task_02, task_03); the `runs`,
  `attach` and `runs show` help text (task_02, task_03); and the ADR-0363
  record (task_03).

Deviations from the PRD: None.

## Absence Evidence

1. `terminalOutput` and the writer guard — command: `grep -rn 'terminalOutput\|is not a terminal file' internal`; result: no match.
2. The closing line — command: `grep -rn 'no Run selected\|Run Browser closed' internal`; result: only the `BrowserOutcome` doc comment in `internal/tui/runbrowser.go`.
3. `RenderRunTable`, `ShortRunID`, `RunRepositoryName` — command: `grep -rn 'RenderRunTable\|ShortRunID\|RunRepositoryName' internal`; result: only the unexported `browserShortRunID` of `internal/tui/runbrowser.go`, which keeps the sixteen-digit suffix.
4. The `runs list` format flags and zone seam — command: `grep -n 'format' internal/cli/runs.go; grep -rn 'runsListLocation' internal`; result: no match.
5. `terminalColumns` — command: `grep -rn 'IoctlGetWinsize\|TIOCGWINSZ' internal`; result: no match.
6. `ResolveRunReference` — command: `grep -rn 'ResolveRunReference\|AmbiguousRunReference\|LIKE .run' internal/store`; result: no match; `selectRun` uses `WHERE id = ?` only.

## ADR Alternatives

- ADR-0357: it rejected keeping `runs list` refusing outside a repository and
  making every command lenient; this design keeps both rejections, keeps the
  repository column rule for the table, and changes no command's leniency.
- ADR-0362: it rejected asking Git with `git rev-parse`, refusing only a
  home-directory `.git`, and the other `gitdir:` variants; this design finds
  the repository through the same config walk and takes none of them.
- ADR-0363: written here; it takes none of its own rejected options (one
  format, a flag-only table, the full id, the prefix everywhere, a fixed
  zone).

## Integration Points

Bubble Tea v2 programs keep their options; only the writer they get changes.
`golang.org/x/sys/unix` (already required) gives the terminal width on Unix.
No network, no provider, no other repository.

## Testing Approach

Model and command tests through existing seams: `withRunBrowserSession`,
`withRunsListNow`, a new `withRunsListLocation`, `setCommandEnvForTest`
(`ROUNDFIX_TUI`, `COLUMNS`), `seedRunsForList`, and a raw-database helper that
renames a seeded Run's id (opened with `_pragma=foreign_keys(0)`, as
`setListedRunCreatedAt` opens it) to build colliding prefixes. No test opens a
terminal. Repository directories are made the way the package's tests mark a
Git root today (`withCLIWorkspace`), so Spec 0352's change to that rule
reaches them through the shared helper.

Sabotages, each on a temporary copy, with the test and the input that sees it:

1. Pass `stdout` instead of `terminalOutput(stdout)` in
   `defaultRunBrowserSession`: `TestTerminalProgramsReceiveTheTerminalFile/run_browser`
   fails, its error naming `cli.commandEnvironmentWriter`.
2. The same in `runAttachCockpit`: `/attach_cockpit` fails.
3. The same in `startRunUI`: `/owning_cockpit` fails.
4. Drop the closing line in `runRunBrowserLoop`:
   `TestRunBrowserEndWithoutASelectionPrintsOneLine/attach` fails.
5. `RunRepositoryName` reads `GitRoot` first:
   `TestRunRepositoryNameNamesTheRepository` fails on the Run whose
   `RepositoryRoot` is `/repos/fiscus` and `GitRoot` a worktree path.
6. Auto mode always renders the table:
   `TestRunsListPlainFormatKeepsTheStableLine` fails.
7. Cut `AGENT` instead of `TARGET`, or never cut:
   `TestRunsListTableFitsTheWidthByCuttingTheTarget` (`COLUMNS=80`) fails.
8. Render `STARTED` in UTC: `TestRunsListTableShowsAlignedColumns` fails on
   the `-03` location (`09:02` expected, `12:02` rendered).
9. Accept seven-digit prefixes: `TestResolveRunReferenceNamesOneRunByItsShortID/too_short`
   fails.
10. Return the newest of several matches:
    `TestResolveRunReferenceNamesOneRunByItsShortID/ambiguous` fails.

## Build Order

1. The Run Browser and cockpits draw on the terminal and the Run Browser
   never ends silently; its `REPO` names the repository; the `runs` reference
   and the `runs` and `attach` guides say so (task_01).
2. `runs list` table, `--format`, `--wide`, width and zone; the help text,
   reference and guide (task_02; depends on: 1, for `RunRepositoryName` and
   the shared documents).
3. Short id resolution in the store, `attach` and `runs show`; ADR-0363 and
   its documents (task_03; depends on: 2, for `ShortRunID`, the help text and
   the shared documents).
4. QA gate (depends on: 3).

Spec 0352 is delivered first and may change the Roundfix `SKILL.md` pair,
the version record, `references/runs.md`, the `runs` guide and the Behavior
Surface Record; each Task here rebases onto it and re-runs the record
commands, which recompute versions and fingerprints.

## Risks & Considerations

- An Agent in a pseudo-terminal gets the table; mitigated by the skill text
  and `--format plain`. A script piping stdout is unaffected.
- The MacBook diagnostics may show a second cause of the instant close, such
  as a terminal reply split across reads and read as `Esc`; the closing line
  will name `esc`, and a follow-up Spec treats it. The Secondbrain records
  that the maintainer works through cmux, mosh and tmux, where such a split
  is likelier.
- A short id can collide; the refusal lists every match.
- `terminalColumns` is Unix-only; elsewhere the width falls back to
  `COLUMNS` or no budget. `GOOS=windows go build` proves the stub.

## Glossary

None.

## Decisions

- The writer guard lives in the terminal programs, so any future caller that
  passes a wrapper fails loudly instead of drawing nothing.
- The closing line goes to stderr, because it is a diagnostic and stdout of
  bare `attach` stays empty on cancel.
- The table renderer lives in `internal/tui`, next to the display-width
  helpers and `FormatRunRow`; `runs list` only selects and wires it.
- `STARTED` uses `HH:MM` for today and `Mon DD HH:MM` otherwise; the zone is
  named once in the header. See ADR-0363.
- The short id is accepted only by `attach` and `runs show`. See ADR-0363.
