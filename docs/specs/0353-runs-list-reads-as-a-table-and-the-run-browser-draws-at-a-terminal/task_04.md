---
task: task_04
spec: 0353-runs-list-reads-as-a-table-and-the-run-browser-draws-at-a-terminal
status: completed
type: qa
complexity: medium
---

# Task 04: Run the final QA gate

## Overview

This is the authored terminal gate for this Spec. It declares what the matrix
covers and settles the Spec on evidence. The Spec hands the terminal file to
the Run Browser and both cockpits, makes a terminal program refuse any other
writer, prints one closing line when the Run Browser ends without a
selection, names the repository in the Run Browser, gives `runs list` a table
at a terminal with `--format <table|plain>` and `--wide` while keeping the
stable line for everything else, lets `attach` and `runs show` take the short
Run id, and records the decision in ADR-0363 and the Roundfix skill and user
guides.

Every row runs on this repository's own files, temporary repositories,
temporary Run Databases and temporary copies of the audited head. No row
opens an interactive terminal, starts an Agent session or an ACP adapter,
reads or writes the real `~/.roundfix`, reads another repository, reaches the
network, sends a request to a provider, the judge or a reviewer, or reads a
credential. Every `go test` it runs carries an explicit `-timeout`. Every row
that runs `bin/roundfix` sets `HOME` and `TMPDIR` to fresh temporary
directories, because the CLI loads `$HOME/.roundfix/config.yml`, runs with
every Judge Key Variable and `TYPESAFE_API_KEY` unset, and removes the
directories afterwards. The QA never writes into the checkout while a test
suite runs, and records environment variable names only, never values.

## Requirements

1. MUST build the binary with `make build` and run the repository
   Verification, recording its result as a gate fact.
2. MUST verify Success Metrics 5 and 6 and Invariants 3, 4 and 6 by running
   `go test -count=1 -timeout 20m -v` over `./internal/tui` with
   `TestRunBrowserSessionRefusesAWriterThatIsNotAFile`,
   `TestRunCockpitRefusesAWriterThatIsNotAFile`,
   `TestRunBrowserCancelKeysNameTheirReason`,
   `TestRunRepositoryNameNamesTheRepository` and every `TestRunBrowser` test,
   and over `./internal/cli` with `TestTerminalProgramsReceiveTheTerminalFile`,
   `TestRunBrowserEndWithoutASelectionPrintsOneLine`,
   `TestAttachWithoutRunIDNonInteractiveNamesAllRunsList`,
   `TestRunRunsWithoutSubcommandHonorsInteractivity` and every
   `TestReadOnlyRunCommands` test, recording each pass and each subtest.
3. MUST verify Success Metrics 1 to 3 and Invariants 1 and 2 by running
   `go test -count=1 -timeout 20m -v` over `./internal/cli` with
   `TestRunsListPlainFormatKeepsTheStableLine`,
   `TestRunsListTableShowsAlignedColumns`,
   `TestRunsListTableFitsTheWidthByCuttingTheTarget`,
   `TestRunsListWideShowsEveryColumn`, `TestRunsListFormatTranscripts`,
   `TestRunsHelpNamesTheFormatAndWideFlags` and every `TestRunRunsList` and
   `TestRunsList` test, and over `./internal/tui` with
   `TestShortRunIDIsTheFirstEightSuffixDigits`, recording each pass and each
   subtest; and by running Surface Transcripts 1, 2 and 3 of `_techspec.md`
   through `bin/roundfix` from the checkout under a fresh `HOME` and
   `TMPDIR`, comparing every stdout line, every stderr line and the exit
   code.
4. MUST verify Success Metric 4 and Invariant 5 by running
   `go test -count=1 -timeout 20m -v` over `./internal/store` with
   `TestResolveRunReferenceNamesOneRunByItsShortID`, and over
   `./internal/cli` with `TestAttachAndRunsShowAcceptTheShortID`,
   `TestAnAmbiguousShortIDListsEveryMatch` and every `TestRunsShow` and
   `TestAttach` test, recording each pass and each subtest.
5. MUST verify the declared break from the repository history: between the
   Delivery Base (the merge base of the audited head with `main`) and the
   audited head, the diff of `internal/cli/read_only_run_commands_test.go`
   changes only `assertRunReaderBrowser`, and no existing top-level test
   function of `internal/cli/cli_test.go`, `internal/cli/runs_show_test.go`
   or `internal/tui/runbrowser_test.go` changes its body.
6. MUST run each sabotage `_techspec.md` → Testing Approach lists on its own
   temporary copy of the audited head made with `git archive`, with
   `-timeout` on every `go test`, and record that the named test fails, with
   the case whose input makes the mutation observable.
7. MUST execute the tests task_01, task_02 and task_03 name in their
   Verification, and record that `GOOS=windows go build -buildvcs=false ./internal/cli`
   and `GOOS=linux go build -buildvcs=false ./internal/cli` succeed.
8. MUST search `docs/user-guide/`, `docs/adr/` and both copies of the
   Roundfix skill, and record that: both copies of the `runs` reference carry
   the three paragraphs task_01, task_02 and task_03 prescribe and are
   byte-identical; the `runs` and `attach` guides carry the sentences those
   Tasks prescribe; `bin/roundfix runs --help`, `bin/roundfix attach --help`
   and `bin/roundfix runs show --help` print the help lines task_02 and
   task_03 prescribe; and ADR-0363 exists with status `accepted`.
9. MUST record, as outside-evidence rows, three sources this Spec did not
   write, each with its date and none blocked:
   - the maintainer's report quoted in
     `references/2026-11-05-macbook-runs-list-is-confusing-and-attach-shows-nothing.md`
     (2026-11-05): compare offline its two stable lines and its two silent
     `attach` calls with `TestRunsListTableShowsAlignedColumns`,
     `TestTerminalProgramsReceiveTheTerminalFile` and
     `TestRunBrowserEndWithoutASelectionPrintsOneLine`;
   - the Command Line Interface Guidelines passage `_prd.md` → Acceptance
     evidence quotes (a TTY tells whether a human reads the stream, and
     `--plain` keeps machine-readable output): compare it with the format
     selection and `--format plain`;
   - the Git revision passage `_prd.md` → Acceptance evidence quotes ("a
     leading substring that is unique within the repository"): compare it
     with the short id resolution and its ambiguity refusal.
10. MUST record exactly one row for the PRD's one Unreachable Acceptance
    entry (the Run Browser drawing and closing at the maintainer's MacBook
    terminal) as `blocked (environment: no interactive terminal in the Run)`,
    counting it once in `rows_blocked_declared`, and back it with the rows of
    requirement 2.
11. MUST perform the glossary check (`docs/agents/domain.md`): record that the
    PRD and TechSpec declare `None.` under `## Glossary`, that the Task commit
    range changes no `CONTEXT.md` entry, and that no phrase the Tasks wrote
    into the skill or guides needs a new domain term.
12. MUST verify from the repository history:
    - each Task's changed files stay within its declarations or its `##
      Recorded paths`;
    - every changed Governed Path is listed in `_authorization.md` `paths`;
    - the Task commit range changes none of these: the Makefile, the CI
      workflows, `go.mod`, `go.sum`, `.roundfixrc.yml`, `skills-lock.json`,
      `CHANGELOG.md`, `internal/app/version.go`, `dist/`, `CONTEXT.md`,
      `docs/references/skill-coverage.json`, any existing Archive Record, any
      Baseline asset, `docs/agents/`, `internal/config`, any Roundfix
      reference other than `runs`, any skill other than the Roundfix files
      the Tasks declare;
    - no `### QA settlement` section changed.

    This row reads Task commits, so it declares a `commit_range` input.
13. MUST record the non-waivable Pull Request row as
    `blocked (environment: no open Pull Request)`, naming the Pull Request
    row in its provenance. It supports the row with the pre-PR equivalent of
    each control:
    - approval from the recorded delivery authority;
    - checks from the Daemon's repository Verification at the audited head;
    - no unresolved threads, because no Pull Request exists;
    - no Merge-Ready acceptance yet;
    - the audited head named as the claimed candidate.
14. MUST, for the Outcome for the Archive Record, run only
    `h="$(mktemp -d)"; t="$(mktemp -d)"; HOME="$h" TMPDIR="$t" bin/roundfix archive 0353-runs-list-reads-as-a-table-and-the-run-browser-draws-at-a-terminal --plan --offline; rm -rf "$h" "$t"`
    with every Judge Key Variable and `TYPESAFE_API_KEY` unset, and write the
    report's Outcome as one paragraph that does not end with a colon. It MUST
    NOT do any of these:
    - accept a row whose only evidence is that a file was read;
    - read or write `~/.roundfix`;
    - record any environment value;
    - open an interactive terminal program;
    - send any request to a provider, the judge or a reviewer;
    - write into the checkout while a test suite runs.

## Subtasks

- [ ] Build the binary and run the repository Verification.
- [ ] Run the named Run Browser, table and short id tests and the three transcripts.
- [ ] Check the declared break in the history and run each sabotage on a temporary copy.
- [ ] Record the documents, the outside sources, the unreachable row, the glossary check and the scope audit, and write the QA report.

## Acceptance Criteria

- [ ] The terminal programs get the terminal file, the Run Browser never ends silently, `runs list` keeps its stable line and prints the table, and the short id names one Run, in the tests and the transcripts.
- [ ] Each sabotage makes its named test fail, and the scope audit finds no change outside the declared files.

## Verification

- `newest="$(find 'docs/specs/0353-runs-list-reads-as-a-table-and-the-run-browser-draws-at-a-terminal/qa' -type f -name "qa-report-*.md" -print 2>/dev/null | awk "{ report=\$0; name=\$0; parts=split(name, path, \"/\"); name=path[parts]; name=substr(name, 11, length(name)-13); date=substr(name, 1, 10); suffix=substr(name, 11); shape=(length(date) == 10 && substr(date, 5, 1) == \"-\" && substr(date, 8, 1) == \"-\"); for (i=1; i <= 10 && shape; i++) { if (i != 5 && i != 8 && index(\"0123456789\", substr(date, i, 1)) == 0) shape=0 } year=substr(date, 1, 4)+0; month=substr(date, 6, 2)+0; day=substr(date, 9, 2)+0; leap=(year%400 == 0 || (year%4 == 0 && year%100 != 0)); days=(month == 2 ? 28+leap : ((month == 4 || month == 6 || month == 9 || month == 11) ? 30 : 31)); dated=(shape && month >= 1 && month <= 12 && day >= 1 && day <= days); if (!dated) date=\"\"; if (suffix == \"\") { sequenced=1; sequence=-1 } else { sequenced=(length(suffix) > 1 && substr(suffix, 1, 1) == \"-\"); for (i=2; i <= length(suffix) && sequenced; i++) { if (index(\"0123456789\", substr(suffix, i, 1)) == 0) sequenced=0 } sequence=(sequenced ? substr(suffix, 2)+0 : 0) } if (dated && sequenced) printf \"%s\\t%s\\t%s\\t%s\\t%s\\n\", dated, date, sequenced, sequence, report }" | sort -k1,1n -k2,2 -k3,3n -k4,4n -k5,5 | tail -1 | cut -f5-)"; test -n "$newest" || exit 1; awk "BEGIN { whitespace=\" \t\r\n\f\v\" } NR == 1 && \$0 == \"---\" { frontmatter=1; next } frontmatter && \$0 == \"---\" { closed=1; exit } frontmatter && index(\$0, \"verdict:\") == 1 { verdict=substr(\$0, 9); while (length(verdict) > 0 && index(whitespace, substr(verdict, 1, 1)) > 0) verdict=substr(verdict, 2); while (length(verdict) > 0 && index(whitespace, substr(verdict, length(verdict), 1)) > 0) verdict=substr(verdict, 1, length(verdict)-1); verdicts++ } END { exit(closed && verdicts == 1 && length(verdict) > 0 ? 0 : 1) }" "$newest"`

## References

- `_prd.md` → Goals; User Stories 1-5; Core Features 1-7; Success Metric 1; Success Metric 2; Success Metric 3; Success Metric 4; Success Metric 5; Success Metric 6; Success Metric 7; Acceptance evidence; Unreachable Acceptance; Glossary
- `_techspec.md` → Current behavior; API Contracts 1-10; Surface Transcript 1; Surface Transcript 2; Surface Transcript 3; Invariants 1-6; Declared breaks; Testing Approach; Build Order 4
