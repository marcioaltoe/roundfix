---
task: task_01
spec: 0353-runs-list-reads-as-a-table-and-the-run-browser-draws-at-a-terminal
status: completed
type: frontend
complexity: medium
---

# Task 01: The Run Browser and both cockpits draw on the terminal, and the Run Browser never ends silently

## Overview

Bare `roundfix attach` and bare `roundfix runs` open a Run Browser that draws
nothing, because the command hands Bubble Tea its environment wrapper of
stdout instead of the terminal file (`_techspec.md` → Current behavior 1–3);
the attach cockpit and the owning cockpit of a foreground Run get the same
wrapper. After this Task the three terminal programs receive the terminal
file, a terminal program refuses any writer that is not a file with an error
naming it, a Run Browser closed without a selection prints one stderr line
naming the key, and the Run Browser's `REPO` column names the repository.
This answers the "roundfix attach não mostra nada" of the report reported by
the maintainer on 2026-11-05
([references](references/2026-11-05-macbook-runs-list-is-confusing-and-attach-shows-nothing.md)).

## Requirements

1. MUST, before changing code, run
   `go test -count=1 -timeout 20m -v -run '^(TestReadOnlyRunCommandsKeepTheirOutputInsideARepository|TestAttachRunBrowserCancelExitsZeroWithoutAttaching|TestAttachWithoutRunIDNonInteractiveNamesAllRunsList|TestRunRunsWithoutSubcommandHonorsInteractivity)$' ./internal/cli`
   and record in the Result that it passes, with the browser's stderr equal
   to the warning alone.
2. MUST add `terminalOutput` to `internal/cli` (API Contract 1) and pass
   `terminalOutput(stdout)` to `roundtui.RunBrowserSession` in
   `defaultRunBrowserSession`, `terminalOutput(stdout)` to
   `roundtui.RunCockpit` in `runAttachCockpit`, and `terminalOutput(stderr)`
   to `roundtui.RunCockpit` in `startRunUI`; every color and TUI test keeps
   reading the wrapper.
3. MUST make `RunBrowserSession` and `RunCockpit` refuse a writer that is not
   an `*os.File` before building a program or a model, with exactly
   `run browser session: output %T is not a terminal file` and
   `run cockpit: output %T is not a terminal file` (API Contract 2).
4. MUST add `Reason` to `roundtui.BrowserOutcome`, set to `q`, `esc` or
   `ctrl+c` by the cancel keys and to `interrupted` when the program ended
   without a key outcome (API Contract 3).
5. MUST make `runRunBrowserLoop` and the empty pass of `runRunsBrowserCommand`
   write, after every cancelled outcome, exactly one stderr line
   `roundfix: Run Browser closed (<reason>); no Run selected.`, or
   `roundfix: Run Browser closed; no Run selected.` when the reason is empty,
   and return exit `0` with nothing on stdout (API Contract 4).
6. MUST add `roundtui.RunRepositoryName` (API Contract 5) in
   `internal/tui/runrow.go` and use it for the Run Browser's repository
   column in place of `browserRepoName`.
7. MUST create `internal/tui/terminal_output_test.go`, each test calling
   `t.Parallel()` first, with:
   - `TestRunBrowserSessionRefusesAWriterThatIsNotAFile`: a `*bytes.Buffer`
     and a struct embedding `io.Writer` each give the API Contract 2 error
     naming their type;
   - `TestRunCockpitRefusesAWriterThatIsNotAFile`: the same for `RunCockpit`;
   - `TestRunBrowserCancelKeysNameTheirReason`: `q`, `esc` and `ctrl+c` give
     `Cancelled` with that `Reason`, and `enter` on a row gives no reason;
   - `TestRunRepositoryNameNamesTheRepository`: `RepositoryRoot`
     `/repos/fiscus` with `GitRoot`
     `/worktrees/fiscus-2f7556b8/roundfix-deliver-0063-x-411942b2` gives
     `fiscus`; only `GitRoot` `/repos/quaestor` gives `quaestor`; neither
     gives `-`; and a Run Browser over that first Run renders `fiscus` and
     not `roundfix-deliver`.
8. MUST create `internal/cli/run_browser_terminal_test.go`, each test calling
   `t.Parallel()` first, with:
   - `TestTerminalProgramsReceiveTheTerminalFile`, subtests `run browser`
     (`defaultRunBrowserSession` with `commandWriter(&bytes.Buffer{}, …)`
     returns an error naming `*bytes.Buffer` and not
     `commandEnvironmentWriter`), `attach cockpit` (`runAttachCockpit` over a
     seeded Run with the same wrapped writer exits `1` and its stderr names
     `*bytes.Buffer`), and `owning cockpit` (`startRunUI` with a wrapped
     `*bytes.Buffer` whose environment sets `ROUNDFIX_TUI=always` delivers on
     its cockpit channel an error naming `*bytes.Buffer`; close the UI
     afterwards);
   - `TestRunBrowserEndWithoutASelectionPrintsOneLine`, subtests `attach`,
     `runs` and `runs without a Run Database` (scripted outcome
     `Cancelled: true, Reason: "q"`; stderr exactly
     `roundfix: Run Browser closed (q); no Run selected.\n`, stdout empty,
     exit `0`) and `no reason` (scripted `Cancelled: true`; stderr exactly
     `roundfix: Run Browser closed; no Run selected.\n`).
9. MUST update `assertRunReaderBrowser`
   (`internal/cli/read_only_run_commands_test.go`) as `_techspec.md` →
   Declared breaks states, and keep every other test of `./internal/tui` and
   every `TestAttach*`, `TestRunRuns*`, `TestReadOnlyRunCommands*` and
   `TestBrowserAttachCockpit*` test of `./internal/cli` passing unchanged.
10. MUST add to `.agents/skills/roundfix/references/runs.md`, after the
    paragraph that ends "lists every repository outside one.", verbatim apart
    from line wrapping:
    "When the Run Browser closes without a selection, bare `runs` and
    `attach` print one line on stderr and exit `0`:
    `roundfix: Run Browser closed (<key>); no Run selected.`, where `<key>`
    is `q`, `esc` or `ctrl+c`, or `interrupted` when a signal ended it. The
    Run Browser and the Live Run View draw on the terminal itself; given any
    other output they refuse with `output <type> is not a terminal file`.
    The Run Browser's repository column names the repository, not the Run
    Worktree."
    Run `make skills-sync`, then
    `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions`.
11. MUST add to `docs/user-guide/commands/runs.md` after "In a
    non-interactive context it exits `2` naming `runs list`." and to
    `docs/user-guide/commands/attach.md` after the sentence that ends "in
    non-interactive mode it exits `2` naming `roundfix runs list`.", verbatim
    apart from line wrapping: "When the Run Browser closes without a
    selection it prints `roundfix: Run Browser closed (<key>); no Run
    selected.` on stderr and exits `0`." Then, as the LAST step after the
    last guide edit, re-record the Behavior Surface Record with
    `go test -count=1 -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage`.
12. MUST run sabotages 1 to 5 of `_techspec.md` → Testing Approach on a
    temporary copy and record which named test fails.
13. MUST NOT change the Run Browser's keys, layout, short id or column drop
    order, the non-interactive refusals of bare `attach` and bare `runs`,
    `runs list`, `internal/config`, `internal/store`,
    `internal/cli/cli_test.go`, any other skill or reference, the Makefile, CI workflows, `go.mod` or `go.sum`.

## Subtasks

- [ ] Characterize the cancel path and the wrapped writer.
- [ ] Hand the terminal file to the three terminal programs and guard the writer.
- [ ] Name the closing key and print the closing line.
- [ ] Name the repository in the Run Browser.
- [ ] Add the model and command tests; run the sabotages on a temporary copy.
- [ ] Document the closing line; sync and record the skill; re-record the Behavior Surface Record last.

## Acceptance Criteria

- [ ] The Run Browser, the attach cockpit and the owning cockpit receive the terminal file, and a terminal program given a non-file writer refuses naming its type.
- [ ] A Run Browser closed without a selection prints exactly one stderr line naming its key, stdout stays empty and the exit is `0`.
- [ ] The Run Browser's repository column names the repository for a queued Run.

## Context

- interface: `internal/tui/runbrowser.go`
- interface: `internal/tui/cockpit.go`
- interface: `internal/tui/runrow.go`
- interface: `internal/tui/runbrowser_test.go`
- creates: `internal/tui/terminal_output_test.go`
- interface: `internal/cli/runbrowser.go`
- interface: `internal/cli/attach.go`
- interface: `internal/cli/runui.go`
- creates: `internal/cli/tui_output.go`
- creates: `internal/cli/run_browser_terminal_test.go`
- interface: `internal/cli/read_only_run_commands_test.go`
- interface: `.agents/skills/roundfix/references/runs.md`
- interface: `skills/roundfix/references/runs.md`
- interface: `.agents/skills/roundfix/SKILL.md`
- interface: `skills/roundfix/SKILL.md`
- interface: `skills/testdata/owned-skill-versions.json`
- interface: `docs/user-guide/commands/runs.md`
- interface: `docs/user-guide/commands/attach.md`
- interface: `docs/references/behavior-surfaces.json`

## Consumers

- `roundtui.BrowserOutcome`, `RunBrowserSession`, `RunCockpit` and `RunRepositoryName` — consumers: `./internal/tui`, `./internal/cli`

## Verification

- `out="$(go test -count=1 -timeout 20m -v -run '^(TestRunBrowserSessionRefusesAWriterThatIsNotAFile|TestRunCockpitRefusesAWriterThatIsNotAFile|TestRunBrowserCancelKeysNameTheirReason|TestRunRepositoryNameNamesTheRepository|TestRunBrowser.*|TestFormatRunRowSharedByBothSurfaces)$' ./internal/tui 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestRunBrowserSessionRefusesAWriterThatIsNotAFile TestRunCockpitRefusesAWriterThatIsNotAFile TestRunBrowserCancelKeysNameTheirReason TestRunRepositoryNameNamesTheRepository; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { echo "missing pass: $name" >&2; exit 1; }; done; out="$(go test -count=1 -timeout 20m -v -run '^(TestTerminalProgramsReceiveTheTerminalFile|TestRunBrowserEndWithoutASelectionPrintsOneLine|TestReadOnlyRunCommands.*|TestAttach.*|TestRunRuns.*|TestBrowserAttachCockpit.*)$' ./internal/cli 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestTerminalProgramsReceiveTheTerminalFile TestRunBrowserEndWithoutASelectionPrintsOneLine; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { echo "missing pass: $name" >&2; exit 1; }; done; for f in .agents/skills/roundfix/references/runs.md skills/roundfix/references/runs.md docs/user-guide/commands/runs.md docs/user-guide/commands/attach.md; do tr -s '[:space:]' ' ' < "$f" | grep -qF -- 'no Run selected.' || { echo "missing closing line in $f" >&2; exit 1; }; done; cmp -s .agents/skills/roundfix/references/runs.md skills/roundfix/references/runs.md || { echo 'runs reference mirror differs' >&2; exit 1; }; go test -count=1 -timeout 20m ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' && go test -count=1 -timeout 20m -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$'` — expected: exit 0. Before this Task none of the six new tests exists, so the first pass line is missing; afterwards the terminal programs get the terminal file, the closing line prints, the Run Browser names the repository, the skill and guides state it with a current version record and Behavior Surface Record, and the existing browser, attach and runs tests pass.

## References

- `_prd.md` → Goals; User Story 4; Core Feature 5; Core Feature 6; Success Metric 5; Success Metric 6
- `_techspec.md` → Current behavior 1–5; API Contract 1; API Contract 2; API Contract 3; API Contract 4; API Contract 5; Invariant 3; Invariant 4; Invariant 6; Declared breaks; Testing Approach; Build Order 1

## Result

Skills loaded: implement-task, coding-guidelines, bubbletea, golang-error-handling, golang-testing, tech-writer

### What changed

- `internal/cli/tui_output.go` adds `terminalOutput`, which strips the
  command environment wrapper. `defaultRunBrowserSession` and
  `runAttachCockpit` pass `terminalOutput(stdout)`, and `startRunUI` passes
  `terminalOutput(stderr)`. `colorEnabled` and `liveTUIEnabled` still read
  the wrapper.
- `RunBrowserSession` and `RunCockpit` check the writer before they build a
  program or model. A writer that is not an `*os.File` is refused with
  `run browser session: output %T is not a terminal file` or
  `run cockpit: output %T is not a terminal file`.
- `BrowserOutcome.Reason` is set to the cancel key (`q`, `esc`,
  `ctrl+c`). It is `interrupted` when the program ended without a key
  outcome.
- `browserClosedLine` builds the closing line. `runRunBrowserLoop` and the
  empty pass of `runRunsBrowserCommand` print it on stderr after every
  cancelled outcome and exit `0`, with nothing on stdout.
- `roundtui.RunRepositoryName` (`internal/tui/runrow.go`) uses the base name
  of `RepositoryRoot`, then of `GitRoot`, and `-` when both are empty. It
  replaces `browserRepoName`, which is deleted.
- `assertRunReaderBrowser` now expects the warning followed by
  `roundfix: Run Browser closed; no Run selected.\n`, as listed in Declared
  breaks.
- The closing-line paragraph is in the `runs` reference, and `make
  skills-sync` copied it to the mirror. The skill version went from 0.1.185
  to 0.1.186 through `-record-skill-versions`. The sentence is in the `runs`
  and `attach` guides. The last step re-recorded the Behavior Surface Record
  with `-record-skill-coverage`.

### Evidence

- Req 1 (before any code change): the four named tests passed.
  `TestReadOnlyRunCommandsKeepTheirOutputInsideARepository` passed all 10
  subtests, including the 4 browser subtests. Those subtests asserted that
  the browser's stderr was exactly the warning.
- New model tests: `go test -count=1 -run '^(TestRunBrowserSessionRefusesAWriterThatIsNotAFile|TestRunCockpitRefusesAWriterThatIsNotAFile|TestRunBrowserCancelKeysNameTheirReason|TestRunRepositoryNameNamesTheRepository)$' -v ./internal/tui` → all four PASS. `go test -count=1 ./internal/tui` → ok.
- New command tests: `go test -count=1 -v -run '^(TestTerminalProgramsReceiveTheTerminalFile|TestRunBrowserEndWithoutASelectionPrintsOneLine)$' ./internal/cli` → PASS for `run_browser`, `attach_cockpit`, `owning_cockpit`, `attach`, `runs`, `runs_without_a_Run_Database` and `no_reason`.
- Existing tests: `go test -count=1 -timeout 20m -run '^(TestReadOnlyRunCommands.*|TestAttach.*|TestRunRuns.*|TestBrowserAttachCockpit.*)$' ./internal/cli` → ok.
- Other checks: `gofmt -l` clean; `go vet ./internal/cli ./internal/tui` clean; `GOOS=windows go build ./...` ok; `cmp` shows the runs reference and its mirror are identical; `TestTheSkillCoverageMapIsCurrent` ok after recording.
- Sabotages (Req 12) were run on a temporary copy under `/tmp`. The copy was
  made a git repository because the suite guard runs `git ls-files`. Each
  sabotage was reverted and the copy deleted afterwards.
  1. `defaultRunBrowserSession` passes `stdout` →
     `TestTerminalProgramsReceiveTheTerminalFile/run_browser` FAIL with
     `output cli.commandEnvironmentWriter is not a terminal file`.
  2. `runAttachCockpit` passes `stdout` → `/attach_cockpit` FAIL with
     `roundfix attach failed: run cockpit: output cli.commandEnvironmentWriter …`.
  3. `startRunUI` passes `stderr` → `/owning_cockpit` FAIL with
     `run cockpit: output cli.commandEnvironmentWriter …`.
  4. The closing line is removed from `runRunBrowserLoop` →
     `TestRunBrowserEndWithoutASelectionPrintsOneLine/attach`, `/runs` and
     `/no_reason` FAIL with stderr `""`.
  5. `RunRepositoryName` reads `GitRoot` first →
     `TestRunRepositoryNameNamesTheRepository/repository_root_over_the_Run_Worktree`
     FAIL with `"roundfix-deliver-0063-x-411942b2"; want "fiscus"`, and
     `/the_Run_Browser_names_the_repository` also FAIL.

### Acceptance criteria

- Terminal file and writer guard: shown by
  `TestTerminalProgramsReceiveTheTerminalFile`, the two `…RefusesAWriterThatIsNotAFile` tests and sabotages 1–3.
- One closing line, empty stdout, exit `0`: shown by
  `TestRunBrowserEndWithoutASelectionPrintsOneLine`, the updated
  `assertRunReaderBrowser`, `TestRunBrowserCancelKeysNameTheirReason` and
  sabotage 4.
- The repository column names the repository for a queued Run: shown by
  `TestRunRepositoryNameNamesTheRepository` and sabotage 5.

### Notes

- The owning-cockpit subtest waits for the cockpit inside `ui.waitOnce`
  itself, so the deferred `Close` does not wait a second time on the drained
  channel.
- Nothing changed in the keys, layout, short id, column drop order,
  non-interactive refusals, `runs list`, `internal/config`, `internal/store`,
  `cli_test.go`, the Makefile, CI, `go.mod` or `go.sum`.
