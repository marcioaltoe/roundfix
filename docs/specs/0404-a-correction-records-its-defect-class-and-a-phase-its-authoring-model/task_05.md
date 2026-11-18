---
task: task_05
spec: 0404-a-correction-records-its-defect-class-and-a-phase-its-authoring-model
status: completed
type: qa
complexity: high
---

# Task 05: Run the final QA gate

## Overview

This is the authored terminal gate for this Spec. It declares what the
matrix covers and settles the Spec on evidence. A Correction is declared
with a Defect Class or the answer `not-a-correction` and a reason, with a
fixed finding or on its own, into one record per Spec that every clone of
the repository reaches; a phase entry of an Authoring Phase records its
model, effort, runner and Repository Identity; the Delivery Queue reads the
Correction Signs at its merge and parks an undeclared one only where
`delivery.require_defect_class` is on; and `roundfix corrections report`
counts Corrections per Spec and plan defects per authoring model beside its
sample, withholding a comparison below 20 Specs a group.

Every row runs on this repository's own files, temporary repositories,
temporary homes and temporary copies of the audited head. No row starts an
Agent harness, an Agent session or an ACP adapter, starts a Run, a Daemon or
a Delivery Owner, calls ai-memory, reads or writes the real `~/.roundfix`,
runs the real `gh`, reaches the network, sends a request to a provider, the
judge or a reviewer, reads another repository, or asserts a wall-clock
duration. Every `go test` it runs carries an explicit `-timeout`.

Every row that runs the product runs `bin/roundfix`, the binary `make build`
made from the audited head, and never a `roundfix` found on `PATH`;
requirement 9 alone also runs a binary the gate builds from the Delivery
Base. Every such row takes `m="$(go env GOMODCACHE)"` and
`c="$(go env GOCACHE)"` first, then sets `HOME` and `TMPDIR` to fresh
temporary directories outside the checkout with `GOMODCACHE="$m"` and
`GOCACHE="$c"`, unsets `GH_TOKEN`, `GITHUB_TOKEN`, `TYPESAFE_API_KEY` and
every Judge Key Variable for that command, sets `CLAUDE_CODE_SESSION_ID` and
`CLAUDE_EFFORT` only where a row names them, reads standard input from
`/dev/null`, and removes the temporary directories afterwards.

A row that runs a whole package with `go test` under a temporary `HOME`, the
whole `./internal/cli` package above all, keeps the host's Rust toolchain:
before it replaces `HOME` it takes `r="${RUSTUP_HOME:-$HOME/.rustup}"` and
`k="${CARGO_HOME:-$HOME/.cargo}"` and passes `RUSTUP_HOME="$r"` and
`CARGO_HOME="$k"` to the command, read only, because the Baseline tests of
that package resolve the Rust toolchain through the home directory and fail
on a home that lacks it, as they do on the Delivery Base. The gate never
writes under either directory. When the host has no Rust toolchain the row
records that fact and runs the same command on the Delivery Base, and a
failure that occurs on both is recorded as environment, not as a finding.
The QA never writes into the checkout while a test suite runs. The gate
keeps its evidence under names the repository does not ignore, `.txt` for
command output, and links only files it kept.

## Requirements

1. MUST build the binary with `make build` and run the repository
   Verification, recording its result as a gate fact.
2. MUST run, with `-v`, the tests task_01, task_02, task_03 and task_04 name
   in their Verification, recording every test and subtest by name, and the
   whole packages `./internal/correction`, `./internal/phaserecord`,
   `./internal/delivery`, `./internal/config`, `./internal/runcause`,
   `./internal/baseline`, `./internal/cli` and `./skills`, and
   `go test -tags docscontract ./internal/docscontract`, under the Rust rule
   of the Overview.
3. MUST reproduce, through `bin/roundfix`, each in a temporary `HOME`,
   Surface Transcript 1, Surface Transcript 2, Surface Transcript 3,
   Surface Transcript 4, Surface Transcript 5, Surface Transcript 6,
   Surface Transcript 7 and Surface Transcript 8 of `_techspec.md`:
   Surface Transcripts 2 and 5 in a temporary Git repository with one
   commit on `main` and the folder `docs/specs/0412-example/`, and the
   others in a temporary Git working tree; Surface Transcript 7 with
   `CLAUDE_CODE_SESSION_ID=3f9a1c20` and `CLAUDE_EFFORT=high`; Surface
   Transcript 8 after a fixture review record that holds `F1`. Record standard output, standard error and the exit code of
   each, compared line by line with the transcript, where `<first day>`,
   `<last day>` and `<home>` stand for the values of the run; and record,
   for Surface Transcript 7, the last line of the day file under that `HOME`, whose
   keys are the nine of a phase entry followed by `repository_id`, `model`,
   `model_source`, `effort`, `effort_source` and `runner`.
4. MUST verify `_prd.md` → Success Metrics 1 to 4 through `bin/roundfix` in
   one temporary repository and `HOME`: a declaration with each of the five
   classes and with `not-a-correction`; each refusal of Success Metric 1,
   with the file's bytes unchanged after each; one commit named by two
   declarations and two commits named by one, counted as two and one in the
   report's text and JSON; a fixed finding with a class and without one,
   whose ledger line, but for its recorded time, and standard output equal
   those the Delivery Base's binary writes and prints for the same fixture; a dismissal with a class
   refused; a Correction declared again and withdrawn; and a declaration
   made from a linked worktree and one from a second clone, both in the
   one file the first checkout's report reads.
5. MUST verify Success Metrics 5, 6, 7 and 15 by running, with `-v`,
   `TestCorrectionSignsOfAnItemBranchLikePullRequest892`,
   `TestTheMergeStepParksAnUndeclaredCorrectionSign`,
   `TestADefectClassParkRetriesToTheStageItsHeadSupports`,
   `TestDeliverStatusShowsUndeclaredCorrectionSigns`,
   `TestTheQueueOwnerRecordsTheSignsOfAnItemBeforeItsMerge` and
   `TestRequireDefectClassIsAProjectConfigBoolean`, recording each subtest,
   and through `bin/roundfix`: a Project Config with
   `delivery.require_defect_class: maybe` refused by name, the key in a
   temporary User Config ignored with its one warning line, and the Park
   line's next action read from `ClassifyPark`'s output naming the five
   classes, `not-a-correction`, `roundfix corrections declare` and
   `roundfix deliver retry`.
6. MUST verify Success Metrics 8 and 9 through `bin/roundfix phase start`
   in a temporary repository and `HOME`: each of the four Authoring Phases
   with a session, `--model claude-fable-5-2` and `CLAUDE_EFFORT=high`; with
   no `--model`; with no session; with a model of 300 characters (one
   warning, exit 0); with `--effort` alone; and `release` with `--model`
   (one warning, nine keys). Record every result line, warning, exit code
   and written line, and that an entry written by the Delivery Base's binary
   reads back as `unknown` in the head's report. Then record, in a second
   temporary clone whose `origin` is an address with a user name and a
   password, that its phase entries carry the first clone's
   `repository_id`, that no file under the `HOME` holds the user name or the
   password, and that a third repository with another root commit and the
   same Spec slug gets another `repository_id`.
7. MUST verify Success Metrics 10 to 14 and 17 by running, with `-v`, the
   report tests task_04 names, recording each, and through
   `bin/roundfix corrections report` and `--format json` over the fixture
   of requirement 4: that no output holds `%` or one of the words better,
   worse, best and worst (invariant 8); that a `--spec` view prints a reason
   that reads as an instruction on one line inside double quotes; that
   every count in the JSON is an integer; and, comparing a listing of every
   file under the repository and the `HOME` with sizes and modification
   times before and after the run, that nothing changed but the Run
   Database's `-wal` and `-shm` files (invariant 9). The outside-evidence
   row of `_prd.md` → Success Metric 17 quotes, from `_prd.md` → Acceptance
   evidence, the CONSORT checklist item 16 passage "State the results in
   absolute numbers when feasible (eg, 10/20, not 50%)." and the NCHS
   passages on suppression below a denominator and the footnote that gives
   the reason, and records that the report prints every rate as two counts
   and prints the reason beside each withheld comparison.
8. MUST verify invariant 10 and Success Metric 13: for the same fixture,
   run `bin/roundfix runs causes` and `--format json` and the Delivery
   Base's binary the same way, and record that the text, the JSON and the
   signature digest are equal.
9. MUST compare the audited head with the Delivery Base through the
   product. Build the Delivery Base's binary in a temporary copy made with
   `git archive` and committed in a fresh Git repository, outside the
   checkout. With the key absent, record that both binaries print the same
   standard output, standard error and exit code for
   `roundfix review dispose <id> --fixed-by <commit>` without a declaration
   on the same fixture, for `roundfix phase start release` with a session,
   and for `roundfix deliver status` on the same fixture queue; and that the
   written `release` entry has the same keys in the same order.
10. MUST run each sabotage below on its own temporary copy of the audited
    head made with `git archive` and committed in a fresh Git repository,
    with `-timeout` on every `go test`, and record that the named test
    fails. No sabotage regenerates or synchronizes anything:
    1. `correction.Validate` accepts a `plan` declaration without an
       artifact: `TestValidateNamesTheAcceptedValues`.
    2. `correction.Undeclared` counts a sign covered by a
       `not-a-correction` answer as undeclared:
       `TestUndeclaredCountsEverySignNoStandingDeclarationNames`.
    3. The merge step parks with the key off:
       `TestTheMergeStepParksAnUndeclaredCorrectionSign`.
    4. The sign reader drops the first-parent restriction and reads every
       commit of `D..H`: `TestCorrectionSignsOfAnItemBranchLikePullRequest892`.
    5. `roundfix phase start` exits 1 for a model of 300 characters:
       `TestPhaseStartRecordsTheModelAndEffortOfEachAuthoringPhase`.
    6. The report prints a rate as a percentage:
       `TestCorrectionsReportPrintsNoPercentAndNoRanking`.
    7. The comparison accepts a group of 19 Specs:
       `TestCorrectionsReportComparesOnlyAboveTwentySpecs`.
    8. `roundfix review dispose --dismiss` accepts `--class`:
       `TestReviewDisposeSurfaceTranscript8`.
11. MUST perform the glossary check (`docs/agents/domain.md`): record that
    the PRD declares `adds:` for **Correction**, **Defect Class**,
    **Correction Sign** and **Corrections Report** and `changes:` for
    **Phase Record** and **Review Finding Disposition**, and the TechSpec
    `adds:` for **Repository Identity**; that each added term has one
    Glossary Definition line and each changed one was revised; that the
    **Defect Class** entry lists the five classes and names
    `not-a-correction` as an answer that is not a Defect Class; that no
    other entry changed; and that `TestGlossaryPhrasesAreDefinedTerms`,
    `TestGlossaryLinksNameAcceptedADRs` and the closed-list test of
    `internal/docscontract/user_guide_contract_test.go` pass.
12. MUST verify from the repository history:
    - each Task's changed files stay within its declarations or its
      `## Recorded paths`;
    - every changed Governed Path is listed in `_authorization.md` `paths`;
    - the Task commit range changes none of these: the Makefile, the CI
      workflows, `go.mod`, `go.sum`, `skills-lock.json`, `CHANGELOG.md`,
      `internal/app/version.go`, `dist/npm/roundfix/package.json`,
      `docs/references/coverage-record.json`, `.roundfixrc.yml`,
      `internal/cli/cli_test.go`, `internal/docscontract/publicdocs_test.go`,
      `internal/docscontract/testdata/corpus-golden.json`,
      `internal/baseline/repository_test.go`, `internal/cli/runs_causes.go`,
      `.agents/skills/roundfix/references/runs.md`, any file under
      `docs/agents/` or `.claude/agents/`, any Baseline module, any ADR
      other than ADR-0407, and any owned skill other than `roundfix`,
      `write-idea`, `write-prd`, `write-techspec` and `write-tasks`;
    - no `### QA settlement` section and no `## Declared model` section of
      an authoring skill changed.

    This row reads Task commits, so it declares a `commit_range` input.
13. MUST record the non-waivable Pull Request row as
    `blocked (environment: no open Pull Request)`, naming the Pull Request
    row in its provenance. It supports the row with the pre-PR equivalent
    of each control:
    - approval from the recorded delivery authority;
    - checks from the Daemon's repository Verification at the audited head;
    - no unresolved threads, because no Pull Request exists;
    - no Merge-Ready acceptance yet;
    - the audited head named as the claimed candidate.
14. MUST record the Unreachable Acceptance of `_prd.md` as not reachable in
    this gate, naming its `satisfied-by`, and require nothing of it.
15. MUST, for the Outcome for the Archive Record, run only
    `m="$(go env GOMODCACHE)"; c="$(go env GOCACHE)"; h="$(mktemp -d)"; t="$(mktemp -d)"; HOME="$h" TMPDIR="$t" GOMODCACHE="$m" GOCACHE="$c" bin/roundfix archive 0404-a-correction-records-its-defect-class-and-a-phase-its-authoring-model --plan --offline; rm -rf "$h" "$t"`
    with every Judge Key Variable and `TYPESAFE_API_KEY` unset, and write
    the report's Outcome as one paragraph that does not end with a colon.
    It MUST require only the invariants `_techspec.md` numbers, and nothing
    `_techspec.md` → Interfaces lists as not an invariant: Roundfix Home
    unchanged by a delivery with the key off, the fixing commit of a fixed
    finding as a second sign, the help of `review dispose`, `phase`, `gc`
    and the root command unchanged, an Authoring Phase entry byte-identical
    to the Delivery Base, or the same phase entry shown from another clone.
    It MUST NOT do any of these:
    - accept a row whose only evidence is that a file was read;
    - run a `roundfix` found on `PATH`;
    - read or write `~/.roundfix`, or write under the host's Rust
      toolchain directories;
    - run the real `gh`, call ai-memory or reach the network;
    - start an Agent harness, an Agent session, an ACP adapter, a Run, a
      Daemon or a Delivery Owner;
    - send any message, or any request to a provider, the judge or a
      reviewer;
    - assert a wall-clock duration;
    - write into the checkout while a test suite runs;
    - read a ref, tag or history of this repository that a fresh checkout
      or the Linux Lane's copy lacks;
    - link a file it did not keep, or keep evidence under a name the
      repository ignores.

## Subtasks

- [ ] Build the binary and run the repository Verification.
- [ ] Run the named tests of the four Tasks and their packages, keeping the host's Rust toolchain.
- [ ] Reproduce the eight Surface Transcripts and the comparison with the Delivery Base.
- [ ] Verify the declarations, the signs, the park, the phase entries and the report through the product.
- [ ] Run each sabotage on a temporary copy.
- [ ] Record the glossary check, the scope audit and the Pull Request row, and write the QA report.

## Acceptance Criteria

- [ ] Every Correction is recorded with its class or `not-a-correction` and its reason, a refusal writes nothing, and a declaration from any clone reaches one record.
- [ ] With the key off a delivery, a fixed disposition and `roundfix deliver status` behave as on the Delivery Base; with it on an undeclared sign parks the item and a declaration with a retry merges it.
- [ ] Each Authoring Phase entry records its model, effort, runner and Repository Identity, and no other phase's entry changes.
- [ ] The report shows counts over counts, withholds a comparison below 20 Specs a group, prints no percentage or ranking, and changes no file.
- [ ] Each sabotage makes its named test fail.

## Verification

- `newest="$(find 'docs/specs/0404-a-correction-records-its-defect-class-and-a-phase-its-authoring-model/qa' -type f -name "qa-report-*.md" -print 2>/dev/null | awk "{ report=\$0; name=\$0; parts=split(name, path, \"/\"); name=path[parts]; name=substr(name, 11, length(name)-13); date=substr(name, 1, 10); suffix=substr(name, 11); shape=(length(date) == 10 && substr(date, 5, 1) == \"-\" && substr(date, 8, 1) == \"-\"); for (i=1; i <= 10 && shape; i++) { if (i != 5 && i != 8 && index(\"0123456789\", substr(date, i, 1)) == 0) shape=0 } year=substr(date, 1, 4)+0; month=substr(date, 6, 2)+0; day=substr(date, 9, 2)+0; leap=(year%400 == 0 || (year%4 == 0 && year%100 != 0)); days=(month == 2 ? 28+leap : ((month == 4 || month == 6 || month == 9 || month == 11) ? 30 : 31)); dated=(shape && month >= 1 && month <= 12 && day >= 1 && day <= days); if (!dated) date=\"\"; if (suffix == \"\") { sequenced=1; sequence=-1 } else { sequenced=(length(suffix) > 1 && substr(suffix, 1, 1) == \"-\"); for (i=2; i <= length(suffix) && sequenced; i++) { if (index(\"0123456789\", substr(suffix, i, 1)) == 0) sequenced=0 } sequence=(sequenced ? substr(suffix, 2)+0 : 0) } if (dated && sequenced) printf \"%s\\t%s\\t%s\\t%s\\t%s\\n\", dated, date, sequenced, sequence, report }" | sort -k1,1n -k2,2 -k3,3n -k4,4n -k5,5 | tail -1 | cut -f5-)"; test -n "$newest" || exit 1; awk "BEGIN { whitespace=\" \t\r\n\f\v\" } NR == 1 && \$0 == \"---\" { frontmatter=1; next } frontmatter && \$0 == \"---\" { closed=1; exit } frontmatter && index(\$0, \"verdict:\") == 1 { verdict=substr(\$0, 9); while (length(verdict) > 0 && index(whitespace, substr(verdict, 1, 1)) > 0) verdict=substr(verdict, 2); while (length(verdict) > 0 && index(whitespace, substr(verdict, length(verdict), 1)) > 0) verdict=substr(verdict, 1, length(verdict)-1); verdicts++ } END { exit(closed && verdicts == 1 && length(verdict) > 0 ? 0 : 1) }" "$newest"`

## References

- `_prd.md` → Goals; User Stories 1–12; Core Features 1–22; Success Metrics 1–17; Acceptance evidence; Unreachable Acceptance; Glossary
- `_techspec.md` → Interfaces (invariants 1–11); API Contracts 1–12; Surface Transcripts 1–8; Testing Approach; Build Order 5; Glossary
- [`references/2026-11-17-fable-and-opus-authoring-study.md`](references/2026-11-17-fable-and-opus-authoring-study.md), the study of 2026-11-17
- ADR-0417, ADR-0418, ADR-0419, ADR-0420, ADR-0421, ADR-0422, ADR-0423, ADR-0424, ADR-0425, ADR-0426, ADR-0407
