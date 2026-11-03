---
task: task_04
spec: 0343-the-follow-up-names-what-waits-and-where-each-item-went
status: completed
type: qa
complexity: medium
---

# Task 04: Run the final QA gate

## Overview

This is the authored terminal gate for this Spec. It declares what the matrix
covers and settles the Spec on evidence. The Spec makes the follow-up's
`repository-inbox` source and `roundfix deliver plan` skip hidden inbox
entries, reads a Waiting Declaration into a `waiting` class that
`roundfix follow-up` holds back by default, counts per source, never sends to
the judge and lists with `--all`, reports a malformed declaration on its
item, and teaches the follow-up mode to pass `--all`, write the declaration
and name each listed item's destination.

Every row runs on this repository's own files, temporary repositories and
temporary copies of the audited head. No row starts an Agent session or an
ACP adapter, reads or writes the real `~/.roundfix`, reads another
repository, reaches the network, sends a request to a provider, the judge or
a reviewer, or reads a credential. Every `go test` it runs carries an
explicit `-timeout`. Every row that runs `bin/roundfix` sets `HOME` and
`TMPDIR` to fresh temporary directories, because the CLI loads
`$HOME/.roundfix/config.yml`, runs with every Judge Key Variable and
`TYPESAFE_API_KEY` unset, and removes the directories afterwards. A fixture
repository a row creates has no remote, so the `pull-requests` source never
reaches a forge. The QA never writes into the checkout while a test suite
runs.

## Requirements

1. MUST build the binary with `make build` and run the repository
   Verification, recording its result as a gate fact.
2. MUST verify Success Metrics 1 to 6 and Invariants 1 to 6 by running
   `go test -count=1 -timeout 20m -v` with
   `TestDeliverPlanIgnoresHiddenInboxEntries`,
   `TestFollowUpIgnoresHiddenInboxEntries`,
   `TestFollowUpKeepsItsBytesWithoutAWaitingDeclaration`,
   `TestFollowUpHoldsBackWaitingItems`, `TestFollowUpAllListsWaitingItemsLast`,
   `TestFollowUpReportsAnInvalidWaitingDeclaration`,
   `TestFollowUpHoldsBackWaitingItemsTranscript`,
   `TestFollowUpAllListsWaitingItemsTranscript`,
   `TestFollowUpStatusRefusesAllTranscript`,
   `TestFollowUpListHelpMatchesTranscript` and `TestFollowUpHelpTranscript`
   over `./internal/cli`, and
   `TestTheFollowUpGuidesNameTheJSONFieldsAndMessageTypes` over
   `./internal/cli` with `-tags docscontract`, recording each pass.
3. MUST reproduce Surface Transcript 1, Surface Transcript 2, Surface
   Transcript 3 and Surface Transcript 4 of `_techspec.md` through the built
   binary, in a temporary repository the row creates, commits and removes
   with exactly the fixture the Surface Transcripts section describes, and
   compare standard output, standard error and exit code with each block
   under its controlled variations.
4. MUST reproduce the Fluxus report through the built binary in a temporary
   repository the row creates and removes, holding `docs/_inbox/.gitkeep`,
   `docs/_inbox/.obsidian/app.json`, a visible note, a pending Finding with a
   valid Waiting Declaration and an open Backlog Entry: record that
   `roundfix follow-up --format json` lists no hidden entry, does not list
   the Finding, and gives the `findings` source `waiting` 1; that
   `--all --format json` lists the Finding last with `priority_name`
   `waiting`, its `waiting_on` and `since`; and that `roundfix deliver plan`
   prints one inbox row.
5. MUST compare, for Success Metric 5 and Invariant 1, the built binary with
   a binary built from the Delivery Base (the merge base of the audited head
   with `main`, unpacked with `git archive` into a temporary directory the
   row removes): on one temporary repository holding an open Backlog Entry, a
   pending Finding, a visible note and a messages file, none declaring a wait
   and no hidden inbox entry, run `roundfix follow-up` and `roundfix
   follow-up --format json` with each binary. Before comparing, run the
   built binary's two commands twice and confirm each output equals itself
   (the report prints the repository root but no clock, PID or duration).
   Then record that the text and JSON of the two binaries are byte-identical,
   and that `roundfix follow-up --status` and `--status --format json` are
   byte-identical too.
6. MUST run each sabotage `_techspec.md` → Testing Approach lists on its own
   temporary copy of the audited head made with `git archive`, with
   `-timeout` on every `go test`, and record that the named test fails, with
   the case whose input makes the mutation observable.
7. MUST execute the tests task_01 to task_03 name in their Verification.
8. MUST search `docs/user-guide/`, `docs/agents/`, `CONTEXT.md` and every
   copy of the Roundfix skill, and record that: the follow-up guide and both
   copies of the follow-up reference carry the synopsis with `[--all]`, the
   Waiting Declaration, the `waiting` class, the judge exclusion, the four
   `waiting_error` texts and the new JSON keys; the follow-up and deliver
   guides and both copies of each reference state the hidden-entry rule;
   both copies of the `modes` reference carry the `all` word, the `--all`
   call, the declaration rule and the destination rule; and no live document
   outside `docs/history/` still says that `repository-inbox` reads every
   file in `docs/_inbox/` without the exception.
9. MUST record, as outside-evidence rows, three sources this Spec did not
   write, each with its date and none blocked:
   - the Fluxus report quoted in
     `references/2026-11-02-the-follow-up-names-what-waits-and-where-each-item-went.md`
     (2026-11-02): compare offline its shape (a `.gitkeep` note, items that
     wait on others with no recorded reason) with requirement 4;
   - the Getting Things Done passage `_prd.md` → Acceptance evidence quotes
     ("If delegating, track the waiting for reminder"): compare it with the
     list holding Waiting Items apart from ranked items;
   - the passage on the Waiting For list `_prd.md` → Acceptance evidence
     quotes ("include the date that each item is requested"): compare it with
     `since` being required.
10. MUST record that the PRD declares no Unreachable Acceptance and count no
    row in `rows_blocked_declared`.
11. MUST perform the glossary check (`docs/agents/domain.md`): record that the
    PRD and TechSpec declare `adds: **Waiting Declaration**`, `adds: **Waiting
    Item**` and the three `changes:` entries, that `CONTEXT.md` defines the two
    new terms and carries the three revisions, and that the Task commit range
    changes no other glossary entry.
12. MUST verify from the repository history:
    - each Task's changed files stay within its declarations or its `##
      Recorded paths`;
    - every changed Governed Path is listed in `_authorization.md` `paths`;
    - the Task commit range changes none of these: the Makefile, the CI
      workflows, `go.mod`, `go.sum`, `.roundfixrc.yml`, `skills-lock.json`,
      `CHANGELOG.md`, `internal/app/version.go`, `dist/`,
      `docs/references/skill-coverage.json`, any existing Archive Record, any
      Baseline asset, `docs/agents/`, `internal/judge`, `internal/speccheck`,
      `internal/cli/follow_up_status.go`, any skill other than the Roundfix
      files the Tasks declare;
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
    `h="$(mktemp -d)"; t="$(mktemp -d)"; HOME="$h" TMPDIR="$t" bin/roundfix archive 0343-the-follow-up-names-what-waits-and-where-each-item-went --plan --offline; rm -rf "$h" "$t"`
    with every Judge Key Variable and `TYPESAFE_API_KEY` unset, and write the
    report's ` Outcome` as one paragraph that does not end with a colon. It
    MUST NOT do any of these:
    - accept a row whose only evidence is that a file was read;
    - read or write `~/.roundfix`;
    - record any environment value;
    - send any request to a provider, the judge or a reviewer;
    - write into the checkout while a test suite runs.

## Subtasks

- [ ] Build the binary and run the repository Verification.
- [ ] Run the named tests and reproduce the transcripts and the Fluxus report through the binary.
- [ ] Compare the bytes with the Delivery Base binary and run each sabotage on a temporary copy.
- [ ] Record the wording search, the outside sources, the glossary check and the scope audit, and write the QA report.

## Acceptance Criteria

- [ ] Hidden inbox entries are listed by neither reader, in the tests and through the binary.
- [ ] A Waiting Item is held back with its count, listed last with `--all`, never judged, and a malformed declaration is reported, in the tests and through the binary.
- [ ] A repository without a declaration keeps its bytes against the Delivery Base, each sabotage makes its named test fail, and the scope audit finds no change outside the declared files.

## Verification

- `newest="$(find 'docs/specs/0343-the-follow-up-names-what-waits-and-where-each-item-went/qa' -type f -name "qa-report-*.md" -print 2>/dev/null | awk "{ report=\$0; name=\$0; parts=split(name, path, \"/\"); name=path[parts]; name=substr(name, 11, length(name)-13); date=substr(name, 1, 10); suffix=substr(name, 11); shape=(length(date) == 10 && substr(date, 5, 1) == \"-\" && substr(date, 8, 1) == \"-\"); for (i=1; i <= 10 && shape; i++) { if (i != 5 && i != 8 && index(\"0123456789\", substr(date, i, 1)) == 0) shape=0 } year=substr(date, 1, 4)+0; month=substr(date, 6, 2)+0; day=substr(date, 9, 2)+0; leap=(year%400 == 0 || (year%4 == 0 && year%100 != 0)); days=(month == 2 ? 28+leap : ((month == 4 || month == 6 || month == 9 || month == 11) ? 30 : 31)); dated=(shape && month >= 1 && month <= 12 && day >= 1 && day <= days); if (!dated) date=\"\"; if (suffix == \"\") { sequenced=1; sequence=-1 } else { sequenced=(length(suffix) > 1 && substr(suffix, 1, 1) == \"-\"); for (i=2; i <= length(suffix) && sequenced; i++) { if (index(\"0123456789\", substr(suffix, i, 1)) == 0) sequenced=0 } sequence=(sequenced ? substr(suffix, 2)+0 : 0) } if (dated && sequenced) printf \"%s\\t%s\\t%s\\t%s\\t%s\\n\", dated, date, sequenced, sequence, report }" | sort -k1,1n -k2,2 -k3,3n -k4,4n -k5,5 | tail -1 | cut -f5-)"; test -n "$newest" || exit 1; awk "BEGIN { whitespace=\" \t\r\n\f\v\" } NR == 1 && \$0 == \"---\" { frontmatter=1; next } frontmatter && \$0 == \"---\" { closed=1; exit } frontmatter && index(\$0, \"verdict:\") == 1 { verdict=substr(\$0, 9); while (length(verdict) > 0 && index(whitespace, substr(verdict, 1, 1)) > 0) verdict=substr(verdict, 2); while (length(verdict) > 0 && index(whitespace, substr(verdict, length(verdict), 1)) > 0) verdict=substr(verdict, 1, length(verdict)-1); verdicts++ } END { exit(closed && verdicts == 1 && length(verdict) > 0 ? 0 : 1) }" "$newest"`

## References

- `_prd.md` → Goals; Core Features 1-7; Success Metric 1; Success Metric 2; Success Metric 3; Success Metric 4; Success Metric 5; Success Metric 6; Success Metric 7; Success Metric 8; Acceptance evidence; Unreachable Acceptance; Glossary
- `_techspec.md` → Current behavior; API Contracts 1-9; Surface Transcript 1; Surface Transcript 2; Surface Transcript 3; Surface Transcript 4; Invariants 1-6; Declared breaks; Testing Approach; Build Order 4
