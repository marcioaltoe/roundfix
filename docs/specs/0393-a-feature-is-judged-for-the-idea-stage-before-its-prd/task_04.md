---
task: task_04
spec: 0393-a-feature-is-judged-for-the-idea-stage-before-its-prd
status: pending
type: qa
complexity: medium
---

# Task 04: Run the final QA gate

## Overview

This is the authored terminal gate for this Spec. It declares what the matrix
covers and settles the Spec on evidence. The Spec adds the stage `idea` to
`roundfix spec judge`, which judges one item that is not yet a Spec or one
request of the session for the Idea stage and never refuses anything; rewords
the three route clauses of the Baseline so a feature is judged before its PRD;
carries the rule and the fielded record into `brainstorming`, `write-idea`,
`write-prd` and the Roundfix skill's follow-up mode; and records the decisions
in ADR-0403 and ADR-0404, the notice in ADR-0201 and the glossary.

Every row runs on this repository's own files, temporary homes, temporary
repositories and temporary copies of the audited head. No row starts an Agent
session or an ACP adapter, calls ai-memory, reads or writes the real
`~/.roundfix` or the real home directory, runs the real `gh`, reaches the
network, reads another repository, sends a request to a provider, the judge
or a reviewer, or reads a credential. Every `go test` it runs carries an
explicit `-timeout`. Every row that runs `bin/roundfix` sets `HOME` and
`TMPDIR` to fresh temporary directories outside the checkout, with
`GOMODCACHE` and `GOCACHE` kept at their `go env` values, puts a fake `gh`
script that prints `[]` first on `PATH`, unsets `GH_TOKEN`, `GITHUB_TOKEN`,
`TYPESAFE_API_KEY` and every Judge Key Variable for that command, and removes
the temporary directories afterwards. The QA never writes into the checkout
while a test suite runs. The gate keeps its evidence under names the
repository does not ignore, `.txt` for command output, and links only files it
kept.

The `go test` package and test runs of this gate do not run `bin/roundfix`
and keep the session's own `HOME` and `TMPDIR`: their tests make their own
temporary homes, and some, such as `TestGuidanceCompositionJourney/rust-cli`,
call an installed formatter that needs the session's toolchain configuration
under `HOME`. Only rows that run `bin/roundfix` use fresh temporary homes.

## Requirements

1. MUST build the binary with `make build` and run the repository
   Verification, recording its result as a gate fact.
2. MUST verify Success Metrics 1 to 4 and 9 and Invariants 1 to 5 and 7 by
   running `go test -count=1 -timeout 20m -v` with the seven
   `internal/judge` tests and the ten `internal/cli` tests that task_01 names
   in its Verification, recording every row and subtest by name; then the
   whole `./internal/judge` package and `-run 'SpecJudge|FollowUp|Help|Usage'
   ./internal/cli`, recording each pass.
3. MUST verify Success Metrics 5 to 7 and Invariant 6 by running, with `-v`,
   `TestTheRouteClausesJudgeAFeatureForTheIdeaStage`,
   `TestTheAuthoringEvidenceClausesCarryTheirText` and
   `TestTheAuthoringEvidenceClausesRenderInTheGuide` over `./internal/baseline`
   and `TestTheSpecReferenceDocumentsTheIdeaStage`,
   `TestTheAuthoringSkillsJudgeAFeatureForTheIdeaStage`,
   `TestTheFollowUpModeRunsTheIdeaStageJudgment` and
   `TestEveryOwnedSkillVersionIsRecorded` over `./skills`, recording every
   subtest by name; then the whole `./internal/baseline` and `./skills`
   packages, `-run 'Baseline|Guidance|SetupManifest' ./internal/cli` and
   `-tags docscontract ./internal/docscontract`, recording each pass.
4. MUST reproduce Surface Transcript 1, Surface Transcript 2, Surface
   Transcript 3, Surface Transcript 4 and Surface Transcript 5 of
   `_techspec.md` through
   `bin/roundfix`, in a temporary Git repository built as the fixture there
   describes (a committed `README.md`, an `origin` remote, the two open
   Backlog Entries, and the request file outside that repository), with the
   working directory in that repository. Record standard output, standard
   error and the exit code of each, compared line by line with the
   transcript; `<reason>` is whatever reason the command gives where the gate
   runs, and the record names it. Also record, with the same fixture, that a
   request file is byte-identical after the command and that the `--format
   json` output of Surface Transcript 5 holds `"kind":"request"` and
   `"verdict":null`.
5. MUST reproduce the managed refresh in a temporary copy of the Delivery Base
   (the merge base of the audited head with the default branch) made with
   `git archive` and committed: record that
   `bin/roundfix baseline update --repo . --no-skills --no-history --yes`
   writes into `docs/agents/spec-routing.md`, exactly once each after
   whitespace normalization, "A feature takes the same route when its Idea
   stage is indicated.", "The Idea stage is indicated when any criterion is
   yes, not indicated when every criterion is no, and undecided otherwise.",
   "which sends the judge only the project's name, the request's first line
   and its first 4,000 characters, after the secret scan.", "The judgment
   informs and never refuses: no command fails for it, and a dismissal needs
   only its recorded reason.", "A fix or a refactor is never judged for the
   Idea stage." and "The Task that implements a Surface Transcript asserts the
   whole transcript in one named test:"; and that a second run of the same
   command changes no file.
6. MUST run each sabotage of `_techspec.md` → Testing Approach on its own
   temporary copy of the audited head made with `git archive`, with `-timeout`
   on every `go test`, and record that the named test fails: sabotages 1 to
   11 as listed there. Sabotage 9 edits only the module file and sabotage 10
   only the canonical `write-prd` skill file; neither regenerates anything,
   so the named test fails on the edited copy alone.
7. MUST execute the tests tasks 01, 02 and 03 name in their Verification.
8. MUST search and record:
   - `internal/judge/questions.json` holds the `idea-stage` entry with exactly
     the three criteria `shape_open`, `unverified_assumption` and
     `market_question`, `yes_at_least` 0.7, `no_at_most` 0.3 and the four types
     that are not judged, and no Go file holds a second copy of a criterion's
     text outside tests;
   - no Go file outside `internal/cli/spec_judge.go` and
     `internal/cli/spec_judge_idea.go` and their tests calls
     `JudgeIdeaStage`, and no command reads a verdict, a record or
     `_idea.md`;
   - both copies of `brainstorming`, `write-idea`, `write-prd`, their two
     templates and the Roundfix `modes` and `spec` references carry API
     Contracts 6 and 7, and `.agents/skills` and `skills` hold the same bytes
     for each changed skill file;
   - this repository's `docs/agents/spec-routing.md`, the formatter golden of
     the Standard TypeScript profile and the spec-workflow module carry API
     Contract 5; the three clause ids, their `mandatory` force and their
     absence of `replaces` are unchanged; and every file under
     `internal/baseline/assets/source-baselines/` is byte-identical to the
     Delivery Base;
   - `CONTEXT.md` carries API Contract 8, ADR-0201 carries API Contract 9, and
     ADR-0403 and ADR-0404 are `accepted`;
   - `docs/user-guide/commands/spec.md` documents the stage, both inputs, what
     is sent, the three criteria, both thresholds and the exit codes.
9. MUST record, as outside-evidence rows, the published sources this Spec did
   not write, each with its date and none blocked, through the passages
   `_prd.md` → Acceptance evidence quotes, looking nothing up:
   - Kohavi and colleagues, KDD 2013 (read 2026-11-14): compare "Only one
     third of the ideas tested at Microsoft improved the metric(s) they were
     designed to improve." with the rule of requirement 5, which asks whether
     to build a feature before its PRD;
   - Herbei and Wegkamp, 2005 (read 2026-11-14): compare "The option R
     expresses doubt and is to be used for few observations that are hard to
     classify in an automatic way." with the undecided verdict of Surface
     Transcripts and the two thresholds of requirement 8;
   - TypeSafe documentation, "Noul" (read 2026-11-14): compare "A good
     starting point is 0.7 for yes and 0.3 for no." with the thresholds of
     requirement 8.
   It MUST also record that
   `docs/references/idea-stage-judgment-measurement.md` exists and states a
   cost of US$0.000054 per request, for Success Metric 8.
10. MUST record the two entries of `_prd.md` → Unreachable Acceptance, one row
    each, as `blocked` with `declared: ` followed by the entry's criterion
    text: the live judgment with a real key after the release, and the
    measures after two months of use.
11. MUST perform the glossary check (`docs/agents/domain.md`): record that the
    PRD and TechSpec declare `adds:` for **Idea Stage Judgment** and
    `changes:` for **Advisory Judgment**, that `CONTEXT.md` defines both with
    the text of API Contract 8, and that
    `TestGlossaryPhrasesAreDefinedTerms` passes.
12. MUST verify from the repository history:
    - each Task's changed files stay within its declarations or its `##
      Recorded paths`;
    - every changed Governed Path is listed in `_authorization.md` `paths`;
    - the Task commit range changes none of these: the Makefile, the CI
      workflows, `go.mod`, `go.sum`, `.roundfixrc.yml`, `skills-lock.json`,
      `CHANGELOG.md`, `internal/app/version.go`,
      `dist/npm/roundfix/package.json`,
      `docs/references/coverage-record.json`, `internal/cli/cli_test.go`,
      `internal/cli/follow_up.go`, `internal/cli/follow_up_inventory.go`, any
      file under `internal/speccheck/` or `internal/spec/`, any Vendored Skill
      under `.agents/skills`, any file under
      `internal/baseline/assets/modules/` other than `spec-workflow.json`, any
      file under `internal/baseline/assets/source-baselines/`, any file under
      `docs/user-guide/` other than `commands/spec.md`, any ADR other than
      ADR-0201, and any skill other than `brainstorming`, `write-idea`,
      `write-prd` and `roundfix`;
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
    `h="$(mktemp -d)"; t="$(mktemp -d)"; HOME="$h" TMPDIR="$t" GOMODCACHE="$(go env GOMODCACHE)" GOCACHE="$(go env GOCACHE)" bin/roundfix archive 0393-a-feature-is-judged-for-the-idea-stage-before-its-prd --plan --offline; rm -rf "$h" "$t"`
    with every Judge Key Variable and `TYPESAFE_API_KEY` unset, and write the
    report's Outcome as one paragraph that does not end with a colon. It MUST
    NOT do any of these:
    - accept a row whose only evidence is that a file was read;
    - read or write `~/.roundfix` or the real home directory;
    - run the real `gh`, call ai-memory or reach the network;
    - send any message, or any request to a provider, the judge or a
      reviewer, a request of Surface Transcript 5 included;
    - record any environment value;
    - write into the checkout while a test suite runs;
    - link a file it did not keep, or keep evidence under a name the
      repository ignores.

## Subtasks

- [ ] Build the binary and run the repository Verification.
- [ ] Run the named tests and reproduce the five Surface Transcripts and the managed refresh.
- [ ] Run each sabotage on a temporary copy.
- [ ] Record the searches, the outside sources, the two declared rows, the glossary check and the scope audit, and write the QA report.

## Acceptance Criteria

- [ ] The command judges one inventoried item or the start of one request file, shows every answer with its probability, sends nothing it must not, and exits zero for every report; the five Surface Transcripts reproduce through the built binary.
- [ ] An adopter's managed refresh carries the rule in the three route clauses, with their ids, force and kept guarantees, and both copies of the four skills, the templates, the glossary and ADR-0201 agree with ADR-0403 and ADR-0404.
- [ ] Each sabotage makes its named test fail, and the scope audit finds no change outside the declared files.

## Verification

- `newest="$(find 'docs/specs/0393-a-feature-is-judged-for-the-idea-stage-before-its-prd/qa' -type f -name "qa-report-*.md" -print 2>/dev/null | awk "{ report=\$0; name=\$0; parts=split(name, path, \"/\"); name=path[parts]; name=substr(name, 11, length(name)-13); date=substr(name, 1, 10); suffix=substr(name, 11); shape=(length(date) == 10 && substr(date, 5, 1) == \"-\" && substr(date, 8, 1) == \"-\"); for (i=1; i <= 10 && shape; i++) { if (i != 5 && i != 8 && index(\"0123456789\", substr(date, i, 1)) == 0) shape=0 } year=substr(date, 1, 4)+0; month=substr(date, 6, 2)+0; day=substr(date, 9, 2)+0; leap=(year%400 == 0 || (year%4 == 0 && year%100 != 0)); days=(month == 2 ? 28+leap : ((month == 4 || month == 6 || month == 9 || month == 11) ? 30 : 31)); dated=(shape && month >= 1 && month <= 12 && day >= 1 && day <= days); if (!dated) date=\"\"; if (suffix == \"\") { sequenced=1; sequence=-1 } else { sequenced=(length(suffix) > 1 && substr(suffix, 1, 1) == \"-\"); for (i=2; i <= length(suffix) && sequenced; i++) { if (index(\"0123456789\", substr(suffix, i, 1)) == 0) sequenced=0 } sequence=(sequenced ? substr(suffix, 2)+0 : 0) } if (dated && sequenced) printf \"%s\\t%s\\t%s\\t%s\\t%s\\n\", dated, date, sequenced, sequence, report }" | sort -k1,1n -k2,2 -k3,3n -k4,4n -k5,5 | tail -1 | cut -f5-)"; test -n "$newest" || exit 1; awk "BEGIN { whitespace=\" \t\r\n\f\v\" } NR == 1 && \$0 == \"---\" { frontmatter=1; next } frontmatter && \$0 == \"---\" { closed=1; exit } frontmatter && index(\$0, \"verdict:\") == 1 { verdict=substr(\$0, 9); while (length(verdict) > 0 && index(whitespace, substr(verdict, 1, 1)) > 0) verdict=substr(verdict, 2); while (length(verdict) > 0 && index(whitespace, substr(verdict, length(verdict), 1)) > 0) verdict=substr(verdict, 1, length(verdict)-1); verdicts++ } END { exit(closed && verdicts == 1 && length(verdict) > 0 ? 0 : 1) }" "$newest"`

## References

- `_prd.md` → Goals; User Stories 1-7; Core Features 1-10; Success Metric 1; Success Metric 2; Success Metric 3; Success Metric 4; Success Metric 5; Success Metric 6; Success Metric 7; Success Metric 8; Success Metric 9; Acceptance evidence; Unreachable Acceptance; Glossary
- `_techspec.md` → API Contracts 1-9; Surface Transcript 1; Surface Transcript 2; Surface Transcript 3; Surface Transcript 4; Surface Transcript 5; Invariants 1-7; Declared breaks; Testing Approach; Build Order 4
