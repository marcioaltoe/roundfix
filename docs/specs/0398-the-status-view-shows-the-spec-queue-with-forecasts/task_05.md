---
task: task_05
spec: 0398-the-status-view-shows-the-spec-queue-with-forecasts
status: completed
type: qa
complexity: high
---

# Task 05: Run the final QA gate

## Overview

This is the authored terminal gate for this Spec. It declares what the
matrix covers and settles the Spec on evidence. The Spec adds the **Spec
Queue Table** to the status view of the Follow-up Command: one row for each
Delivery Queue item, planned Spec, Spec being authored and open Backlog
Entry, each start and end an actual time, a **Queue Forecast**, a named
condition with two durations or the reason no forecast is shown, with the
**Forecast Basis** stated under the rows. It adds a pure package that
computes the forecasts and replays the 109 delivered Specs of the backtest,
a reader of phase entries, the delivery history read from the Runs and one
`git log`, the rows, the basis and three sources in the JSON report, and
the text of the Follow-up guide, of the Roundfix skill's `follow-up` and
`modes` references, of the glossary and of the notices in ADR-0353 and
ADR-0407.

Every row runs on this repository's own files, temporary homes, temporary
repositories and temporary copies of the audited head or of the Delivery
Base. No row starts an Agent harness, an Agent session or an ACP adapter,
starts a Run, a Daemon or a Delivery Owner, calls ai-memory, reads or
writes the real `~/.roundfix`, the real Run Database, the real Phase Record
or the maintainer's checkout as a fixture, runs the real `gh`, reaches the
network, reads another repository that the row did not create, sends a
request to a provider, the judge or a reviewer, or asserts a wall-clock
duration. Every `go test` it runs carries an explicit `-timeout`.

Every row that runs the product runs `bin/roundfix`, the binary
`make build` made from the audited head, and never a `roundfix` found on
`PATH`; requirement 8 alone also runs a binary the gate builds from the
Delivery Base. Every such row takes `m="$(go env GOMODCACHE)"` and
`c="$(go env GOCACHE)"` first, then sets `HOME` and `TMPDIR` to fresh
temporary directories outside the checkout with `GOMODCACHE="$m"` and
`GOCACHE="$c"`, sets `TZ=UTC`, puts a fake `gh` script that prints `[]`
first on `PATH`, unsets `GH_TOKEN`, `GITHUB_TOKEN`, `TYPESAFE_API_KEY`,
`CLAUDE_CODE_SESSION_ID` and every Judge Key Variable for that command,
reads standard input from `/dev/null`, and removes the temporary
directories afterwards. The QA never writes into the checkout while a test
suite runs. The gate keeps its evidence under names the repository does not
ignore, `.txt` for command output, and links only files it kept.

The `go test` package and test runs of this gate do not run `bin/roundfix`
and keep the session's own `HOME` and `TMPDIR`: their tests make their own
temporary homes, and some, such as `TestGuidanceCompositionJourney/rust-cli`,
call an installed formatter that needs the session's toolchain
configuration under `HOME`. Only rows that run a built binary use fresh
temporary homes.

## Requirements

1. MUST build the binary with `make build` and run the repository
   Verification, recording its result as a gate fact.
2. MUST verify Success Metric 2 and the method by running
   `go test -count=1 -timeout 15m -v` with the six `internal/queueforecast`
   tests that task_01 names in its Verification, recording every test and
   subtest by name and, from the replay, the counts of each kind; then the
   whole `./internal/queueforecast` package and
   `GOOS=windows go build -buildvcs=false ./internal/queueforecast` and the
   same for `linux`, recording each pass.
3. MUST verify Success Metric 7 and Success Metric 8 by running, with
   `-v`, the `internal/phaserecord` test and the seven `internal/cli` tests
   that task_02 names in its Verification, recording every test and subtest
   by name; then the whole `./internal/phaserecord` package and
   `GOOS=windows go build -buildvcs=false ./internal/phaserecord`,
   recording each pass.
4. MUST verify Success Metric 1, Success Metric 3, Success Metric 4,
   Success Metric 5, Success Metric 6, Success Metric 9, Success Metric 10
   and Success Metric 11 by running, with `-v`, the nine `internal/cli`
   tests that task_03 names in its first two Verification commands, the
   test of the history read beside the forge call alone under
   `-timeout 3m`, recording every test and subtest by name; then
   `-run '^TestFollowUp' ./internal/cli`, recording the pass.
5. MUST verify Success Metric 12 by running, with `-v`, the `skills` tests
   that task_03 and task_04 name in their Verification and
   `-tags docscontract -run '^TestTheFollowUpGuidesNameTheJSONFieldsAndMessageTypes$' ./internal/cli`,
   recording every test and subtest by name; then the whole `./skills`
   package and `-tags docscontract ./internal/docscontract`, recording each
   pass.
6. MUST reproduce Surface Transcript 1 and Surface Transcript 2 of
   `_techspec.md` through `bin/roundfix follow-up --status`, each with the
   fixture the TechSpec describes: a temporary home with no Run Database
   and no Phase Record and, as the working directory, a temporary Git
   repository on the branch `main` with one commit and no remote, which for
   transcript 1 holds the PRD of `0002-export-the-ledger`, the `_idea.md`
   of `0003-a-report-of-late-items` and the waiting Backlog Entry. Record
   standard output, standard error and the exit code of each, compared line
   by line with the transcript, where `<namespace>`, `<root>`, `<reason>`
   and `<state>` are whatever the command prints. For the fixture of
   transcript 1 also record the JSON report: `queue` holds three rows of
   the kinds `planned`, `authoring` and `backlog`; no cell holds a time or
   a number of seconds; `forecast_basis` holds `delivered_specs` 0 and
   `minimum` 11; and `sources` names `runs`, `delivery-queue`,
   `pull-requests`, `git`, `specs`, `release-plan`, `backlog`,
   `phase-record` and `delivery-history` in that order.
7. MUST record, through `bin/roundfix` in a fresh copy of the fixture of
   transcript 1 for each case:
   - a recorded phase: after `bin/roundfix phase start idea --spec 0004-a-queue-of-reports`
     in that repository and home, the status view prints a row for `0004`
     after the row of `0003` and before the Backlog Entry, with the status
     `authoring: idea`, the start `after its plan is merged` and a note
     that reads `idea since ` followed by a month, day, hour and minute,
     and the source line `phase-record: checked 1`;
   - nothing written: a listing of every file under the repository and the
     home with its size and SHA-256 is the same before and after one text
     run and one JSON run;
   - the refusals: `--status --all` and `--status --messages m.json` exit
     `2` with an empty standard output and the diagnostics the Follow-up
     guide documents;
   - an unreadable Backlog Entry: with a second file under `docs/backlog/`
     whose frontmatter never closes, the report exits `0`, prints
     `backlog: unavailable:` with a reason, keeps the rows of `0002` and
     `0003` and prints no Backlog Entry row;
   - an unreadable Phase Record: with a regular file where the directory
     `phase-record` belongs under the home, the report exits `0`, prints
     `phase-record: unavailable:` with a reason and still prints the row of
     `0003` as `authoring`.
8. MUST compare the audited head with the Delivery Base (the merge base of
   the audited head with the default branch) through the product. Build the
   Delivery Base's binary in a temporary copy made with `git archive` and
   committed in a fresh Git repository, outside the checkout. In the
   fixtures of transcripts 1 and 2, with a fresh home for each run, record
   that every line the Delivery Base's `follow-up --status` prints is
   printed by the audited head's in the same order, and that every key of
   the Delivery Base's JSON report is in the audited head's with the same
   value, except `sources`, which holds the same six entries first. Record
   that `follow-up --help` prints the same text from both.
9. MUST run each sabotage below on its own temporary copy of the audited
   head made with `git archive` and committed in a fresh Git repository
   (the suite guard refuses a tree outside Git), with `-timeout` on every
   `go test`, and record that the named test fails on the named subtest.
   No sabotage regenerates or synchronizes anything:
   1. The limit is taken at position `floor(85 × n ÷ 100)`:
      `TestThePercentileIsTheNearestRank/eleven_values` and
      `TestTheReplayReproducesTheBacktest/at_its_start`.
   2. `MinimumHistory` is 10:
      `TestThePercentileIsTheNearestRank/ten_values_give_no_pair` and
      `TestTheReplayReproducesTheBacktest/at_its_start`.
   3. The sum of two sets returns its first set:
      `TestASetOfSumsIsCutTo512Values/pairwise_sum` and
      `TestTheReplayReproducesTheBacktest/one_behind`.
   4. The cut keeps the first 512 values:
      `TestASetOfSumsIsCutTo512Values/513_values`.
   5. `MaximumBehind` is 4:
      `TestForecastGivesEachRowItsCells/fourth_behind`.
   6. A row with a wait of its own is counted as a row behind:
      `TestForecastGivesEachRowItsCells/own_wait_is_not_counted`.
   7. `Phases` also returns prompt entries:
      `TestPhasesReadsBackThePhaseEntriesOfTheWindow/prompt_entries_are_skipped`.
   8. The history takes a Spec's latest Implement Run as its start:
      `TestTheDeliveryHistoryIsReadFromRunsAndOneGitLog/start_is_the_first_implement_run`.
   9. A cancelled Archive Record gives an end:
      `TestTheDeliveryHistoryIsReadFromRunsAndOneGitLog/a_cancelled_record_is_not_delivered`.
   10. A phase entry of any remote counts:
       `TestASpecIsBeingAuthored/another_remote`.
   11. The note is not cut:
       `TestTheNoteIsOneLineOfRepositoryText/a_long_waiting_declaration_is_cut`.
   12. The table cuts a limit to the minute as it cuts a likely value:
       `TestTheSpecQueueTableRendersTimesAndDurations/the_recorded_moment`.
   13. A Backlog Entry row makes the status `state` `work_found`:
       `TestTheStatusViewShowsTheSpecQueueTable/state_is_unchanged`.
   14. The sentence that begins "A row's `note` and `title` are repository
       text" leaves the `--status` section of the canonical `modes`
       reference only:
       `TestTheStatusModeRelaysTheSpecQueueTable/canonical/a_note_is_data`.
10. MUST execute the tests tasks 01, 02, 03 and 04 name in their
    Verification.
11. MUST search and record:
    - the Follow-up guide carries, in its section `## Spec queue`, the
      three paragraphs of task_02 and the seven of task_03, the five
      sentences of JSON fields and the status example with `queue` and
      `forecast_basis`, and the sentences task_03 replaces;
    - both copies of the Roundfix `follow-up` and `modes` references carry
      the text tasks 03 and 04 prescribe, the `--status` section of the
      `modes` reference keeps the paragraph it had on the Delivery Base, and
      the entry file of the skill differs from the Delivery Base only in
      its two version fields;
    - `CONTEXT.md` defines **Spec Queue Table**, **Queue Forecast** and
      **Forecast Basis** with the text task_04 prescribes and carries the
      revised **Follow-up Command**; ADR-0353 and ADR-0407 each carry one
      notice that names ADR-0410; and ADR-0409 and ADR-0410 are `accepted`;
    - no non-test Go file of `internal/queueforecast` imports `os` or
      `io/fs`, calls `time.Now`, `time.Since` or `time.Until`, or imports a
      package outside the standard library;
    - `internal/cli/follow_up_status.go`, `internal/cli/follow_up_history.go`
      and `internal/cli/follow_up_queue.go` call no `time.Now`, and
      `internal/phaserecord/read.go` and the two reader files create, write
      and remove no file;
    - the command dependencies gained the one field `followUpNow` from this
      Spec;
    - the replay fixture holds 109 objects with only the keys `spec`,
      `plan`, `start` and `end`.
12. MUST record, as outside-evidence rows, the published sources this Spec
    did not write, each with its date (read 2026-11-16) and none blocked,
    through the passages `_prd.md` → Acceptance evidence quotes, looking
    nothing up:
    - The Kanban Guide, May 2025: compare "a forecast of how long it should
      take a work item to flow from started to finished" and "should be
      based on historical cycle time" with the end of the Spec in progress
      in `TestTheStatusViewForecastsFromTheDeliveryHistory/in_progress`,
      which is its start plus a time taken from the fixture's own delivered
      Specs;
    - the same guide: compare "The SLE itself has two parts: a period of
      elapsed time and a probability associated with that period" with the
      basis line of requirement 6 and of the recorded moment, which names
      the median and the 85th percentile;
    - ProKanban, *The Benefits of Defining Your Service Level Expectation*:
      compare "The key is to maintain communicating it as a range and level
      of confidence so as not to create confusion toward a guarantee." with
      the two times of a forecast cell and with the guide's sentence "a
      forecast is a range with its confidence, never a commitment";
    - Nick Brown, ASOS Tech Blog, *The Full Monte*: compare the 85th
      percentile right in 361 of 400 forecasts and the 50th in 270 of 400
      with the counts requirement 2 records for the ends forecast at the
      start;
    - Jenkins, issue JENKINS-21099: compare "just give no estimate at all"
      with the cells `not enough history (0 of 11)` of Surface Transcript 1;
    - Mergify, *Merge Queue* API: compare "`null` when it cannot be
      estimated" with the `null` in `likely` and `limit` of each cell that
      shows a condition in requirement 6.
13. MUST record the four criteria of the one entry of `_prd.md` →
    Unreachable Acceptance, one row each, as `blocked` with `declared: `
    followed by that criterion's text, never a paraphrase.
14. MUST perform the glossary check (`docs/agents/domain.md`): record that
    the PRD declares `adds:` for **Spec Queue Table**, **Queue Forecast**
    and **Forecast Basis** and `changes:` for **Follow-up Command**, that
    task_04 wrote those four entries and no other, and that
    `TestGlossaryPhrasesAreDefinedTerms` passes.
15. MUST verify from the repository history:
    - each Task's changed files stay within its declarations or its
      `## Recorded paths`;
    - every changed Governed Path is listed in `_authorization.md` `paths`;
    - the Task commit range changes none of these: the Makefile, the CI
      workflows, `go.mod`, `go.sum`, `skills-lock.json`, `CHANGELOG.md`,
      `internal/app/version.go`, `dist/npm/roundfix/package.json`,
      `docs/references/coverage-record.json`, `internal/cli/cli_test.go`,
      `.roundfixrc.yml`, any Vendored Skill under `.agents/skills`, any
      file under `internal/baseline/`, any file under `docs/agents/`, any
      guide under `docs/user-guide/` other than `commands/follow-up.md`,
      any ADR other than ADR-0353 and ADR-0407, any owned skill other than
      `roundfix`, any Roundfix reference other than `follow-up` and
      `modes`, and any file of `internal/phaserecord` other than `read.go`
      and its test;
    - no `### QA settlement` section changed.

    This row reads Task commits, so it declares a `commit_range` input.
16. MUST record the non-waivable Pull Request row as
    `blocked (environment: no open Pull Request)`, naming the Pull Request
    row in its provenance. It supports the row with the pre-PR equivalent of
    each control:
    - approval from the recorded delivery authority;
    - checks from the Daemon's repository Verification at the audited head;
    - no unresolved threads, because no Pull Request exists;
    - no Merge-Ready acceptance yet;
    - the audited head named as the claimed candidate.
17. MUST, for the Outcome for the Archive Record, run only
    `m="$(go env GOMODCACHE)"; c="$(go env GOCACHE)"; h="$(mktemp -d)"; t="$(mktemp -d)"; HOME="$h" TMPDIR="$t" GOMODCACHE="$m" GOCACHE="$c" bin/roundfix archive 0398-the-status-view-shows-the-spec-queue-with-forecasts --plan --offline; rm -rf "$h" "$t"`
    with every Judge Key Variable and `TYPESAFE_API_KEY` unset, and write the
    report's Outcome as one paragraph that does not end with a colon. It MUST
    NOT do any of these:
    - accept a row whose only evidence is that a file was read;
    - run a `roundfix` found on `PATH`;
    - read or write `~/.roundfix`, the real Run Database, the real Phase
      Record or the real home directory;
    - run the real `gh`, call ai-memory or reach the network;
    - start an Agent harness, an Agent session, an ACP adapter, a Run, a
      Daemon or a Delivery Owner;
    - send any message, or any request to a provider, the judge or a
      reviewer;
    - assert a wall-clock duration, or compare a time a fixture seeded with
      the real clock;
    - write into the checkout while a test suite runs;
    - link a file it did not keep, or keep evidence under a name the
      repository ignores.

## Subtasks

- [ ] Build the binary and run the repository Verification.
- [ ] Run the named tests of the four Tasks and the packages.
- [ ] Reproduce the two Surface Transcripts and the product rows, and compare with the Delivery Base.
- [ ] Run each sabotage on a temporary copy.
- [ ] Record the searches, the outside sources, the four declared rows, the glossary check and the scope audit, and write the QA report.

## Acceptance Criteria

- [ ] The status view prints the Spec Queue Table after the pending Specs with a row for each Delivery Queue item, planned Spec, Spec being authored and open Backlog Entry, and every line and JSON field of the Delivery Base is still there; both Surface Transcripts reproduce through the built binary.
- [ ] The replay reproduces the counts the method was accepted on; a Spec in progress and the three rows behind it show pairs of times, every other wait a named condition with two durations, and a small history no time.
- [ ] The view writes nothing, asks no judge, continues without an unreadable source and refuses `--all` and `--messages` as before.
- [ ] The guide, both copies of the two references, the glossary and the two amended records agree with ADR-0409 and ADR-0410.
- [ ] Each sabotage makes its named test fail, and the scope audit finds no change outside the declared files.

## Verification

- `newest="$(find 'docs/specs/0398-the-status-view-shows-the-spec-queue-with-forecasts/qa' -type f -name "qa-report-*.md" -print 2>/dev/null | awk "{ report=\$0; name=\$0; parts=split(name, path, \"/\"); name=path[parts]; name=substr(name, 11, length(name)-13); date=substr(name, 1, 10); suffix=substr(name, 11); shape=(length(date) == 10 && substr(date, 5, 1) == \"-\" && substr(date, 8, 1) == \"-\"); for (i=1; i <= 10 && shape; i++) { if (i != 5 && i != 8 && index(\"0123456789\", substr(date, i, 1)) == 0) shape=0 } year=substr(date, 1, 4)+0; month=substr(date, 6, 2)+0; day=substr(date, 9, 2)+0; leap=(year%400 == 0 || (year%4 == 0 && year%100 != 0)); days=(month == 2 ? 28+leap : ((month == 4 || month == 6 || month == 9 || month == 11) ? 30 : 31)); dated=(shape && month >= 1 && month <= 12 && day >= 1 && day <= days); if (!dated) date=\"\"; if (suffix == \"\") { sequenced=1; sequence=-1 } else { sequenced=(length(suffix) > 1 && substr(suffix, 1, 1) == \"-\"); for (i=2; i <= length(suffix) && sequenced; i++) { if (index(\"0123456789\", substr(suffix, i, 1)) == 0) sequenced=0 } sequence=(sequenced ? substr(suffix, 2)+0 : 0) } if (dated && sequenced) printf \"%s\\t%s\\t%s\\t%s\\t%s\\n\", dated, date, sequenced, sequence, report }" | sort -k1,1n -k2,2 -k3,3n -k4,4n -k5,5 | tail -1 | cut -f5-)"; test -n "$newest" || exit 1; awk "BEGIN { whitespace=\" \t\r\n\f\v\" } NR == 1 && \$0 == \"---\" { frontmatter=1; next } frontmatter && \$0 == \"---\" { closed=1; exit } frontmatter && index(\$0, \"verdict:\") == 1 { verdict=substr(\$0, 9); while (length(verdict) > 0 && index(whitespace, substr(verdict, 1, 1)) > 0) verdict=substr(verdict, 2); while (length(verdict) > 0 && index(whitespace, substr(verdict, length(verdict), 1)) > 0) verdict=substr(verdict, 1, length(verdict)-1); verdicts++ } END { exit(closed && verdicts == 1 && length(verdict) > 0 ? 0 : 1) }" "$newest"`

## References

- `_prd.md` → Goals; User Stories 1-9; Core Features 1-17; Success Metric 1; Success Metric 2; Success Metric 3; Success Metric 4; Success Metric 5; Success Metric 6; Success Metric 7; Success Metric 8; Success Metric 9; Success Metric 10; Success Metric 11; Success Metric 12; Acceptance evidence; Unreachable Acceptance; Glossary
- `_techspec.md` → API Contracts 1-12; Surface Transcript 1; Surface Transcript 2; Testing Approach; Build Order 5
- ADR-0409, ADR-0410, ADR-0353, ADR-0407
