---
task: task_04
spec: 0307-the-roundfix-skill-modes-and-up-to-three-questions-per-call
status: pending
type: qa
complexity: high
---

# Task 04: Run the final QA gate

## Overview

This is the authored terminal gate for this Spec. It declares what the matrix
covers and settles the Spec on evidence. The Spec rewords the Baseline clause
`clause.core.ask-user-answerable-decisions` to allow at most three independent
questions per call, keeping its identity so an adopter's refresh records it
`retained`, and gives the Roundfix skill its `--follow-up`, `--status` and
`--unattended` Skill Modes in a new topic reference, with OpenCode slash
metadata, answer placement rules and no personal assistant names. The first
word after the skill's name selects the mode with or without a leading `--`
(`/roundfix follow-up`, `$roundfix status`, `@roundfix --unattended`); the
Skill Mode is not the CLI command `roundfix follow-up`, and any other first
word keeps today's behavior.

Every row runs on temporary copies of the audited head and temporary homes
and repositories under the Run's temporary directory. No row reads or writes
the real `~/.roundfix`, reads or writes the Secondbrain, calls ai-memory,
starts an agent harness, reaches the network, sends a request to a provider,
the judge or a reviewer, reads a credential, or reads another repository. The
gate never runs `roundfix spec judge` and never runs an archive plan without
`--offline`. Every row that runs `bin/roundfix` or `go run ./cmd/roundfix`
sets `HOME` and `TMPDIR` to fresh empty temporary directories created for
that row, so no row reads `~/.roundfix` or the machine's shared temporary
directory.

## Requirements

1. MUST build the binary with `make build` and run the repository
   Verification, recording its result as a gate fact.
2. MUST verify Success Metrics 1 and 2:
   - Run `TestTheAskUserClauseAllowsThreeIndependentQuestionsPerCall`,
     `TestTheAskUserClauseRendersInTheAgentInstructionsGuide` and
     `TestAnAdopterRetainsTheAskUserClause` in `./internal/baseline`, and
     record the clause's enforcement, the two sentences found, the guides
     checked and the disposition with its target.
   - Run `TestBaselineClauseForceIsCharacterized`,
     `TestNoTwoBaselineClausesShareText`, `TestCatalogCompatibility` and
     `TestBaselinePlanCharacterization`, and record that they pass.
   - In a temporary copy of the audited head made with `git archive`, run
     `go run -buildvcs=false ./cmd/roundfix baseline update --repo . --no-skills --no-history --yes --format text`
     with `HOME` and `TMPDIR` set to fresh temporary directories, after `git
     init` and one commit of the copy, and record that it reports no file
     change.
   - On temporary copies of the audited head, apply one sabotage per copy and
     record that the named case fails:
     - set the clause's `enforcement` to `mandatory` in
       `internal/baseline/assets/modules/core.json`:
       `TestTheAskUserClauseAllowsThreeIndependentQuestionsPerCall` and
       `TestAnAdopterRetainsTheAskUserClause` fail, the catalog refusing the
       transition target's enforcement;
     - restore, in `docs/agents/agent-instructions.md` only, the clause bullet
       the Delivery Base carries:
       `TestTheAskUserClauseRendersInTheAgentInstructionsGuide` fails in its
       `../../docs/agents/agent-instructions.md` subtest and passes in its
       golden subtest;
     - delete the independence sentence from the Standard TypeScript golden
       `docs/agents/agent-instructions.md` only: the same test fails in its
       golden subtest.
3. MUST verify Success Metric 3:
   - Run `TestTheRoundfixSkillReadsItsModeFromTheRequest`,
     `TestRoundfixSkillIndexNamesEveryReference`,
     `TestRoundfixSkillIndexLinksResolve`,
     `TestRoundfixSkillEntryFileStaysWithinItsBudget` and
     `TestInstallWritesEveryRoundfixReference` in `./skills`, and record the
     metadata value, the invocation forms found, each mode word bare and with
     `--`, the sentence that names the CLI command `roundfix follow-up` as not
     the Skill Mode, the example whose other first word keeps today's
     behavior, the reference count and the entry file's size.
   - On temporary copies, one sabotage per copy, record that the named case
     fails:
     - append a line holding `$ARGUMENTS` to the canonical `SKILL.md` body:
       the `canonical` subtest fails naming the placeholder;
     - delete the `opencode/slash` line from the canonical `SKILL.md`: the
       `canonical` subtest fails naming the metadata;
     - delete `@roundfix unattended` from the Modes section of the embedded
       `skills/roundfix/SKILL.md` only: the `embedded` subtest fails and the
       `canonical` one passes;
     - delete the sentence naming the CLI command `roundfix follow-up` from
       the Modes section of the canonical `SKILL.md`: the `canonical` subtest
       fails naming the missing string;
     - delete "with or without a leading `--`" from the `## How a mode is
       chosen` section of the canonical modes reference: the `canonical`
       subtest fails naming the phrase;
     - delete the `## --status` heading from the canonical modes reference:
       the `canonical` subtest fails naming the section.
4. MUST verify Success Metrics 4 and 5:
   - Run `TestTheFollowUpModeRunsItsStepsInOrder`,
     `TestTheUnattendedModeAsksWithinTheBaselineLimit` and
     `TestTheRoundfixSkillNamesNoPersonalAssistant` in `./skills`, and record
     the eight step headings, the item that first names each checked term,
     and the three pinned sentences.
   - On temporary copies, one sabotage per copy, record that the named case
     fails:
     - swap the third and fourth `--follow-up` items of the canonical modes
       reference: `TestTheFollowUpModeRunsItsStepsInOrder` fails in its
       `canonical` subtest, `--unattended` now first appearing in item 4;
     - remove `memory_message_send` from item 8 of the canonical modes
       reference: `TestTheFollowUpModeRunsItsStepsInOrder` fails in its
       `canonical` subtest, the term found in no item 8;
     - append the word `resolved_to` to item 4 of the canonical modes
       reference: the same subtest fails naming the retired word;
     - change "at most three questions per call" to "at most four questions
       per call" in the canonical modes reference only:
       `TestTheUnattendedModeAsksWithinTheBaselineLimit` fails in its
       `canonical` subtest;
     - change the same phrase in `internal/baseline/assets/modules/core.json`
       only: the same test fails before its subtests, naming the core module;
     - delete "and never an authorization" from the canonical modes
       reference: the same test fails in its `canonical` subtest;
     - add the word `GT3` to the embedded modes reference only:
       `TestTheRoundfixSkillNamesNoPersonalAssistant` fails in its `embedded`
       subtest.
   - Run `make skills-sync-check` and
     `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$'`, and
     record the Roundfix skill's version at the Delivery Base and at the
     audited head.
5. MUST execute the tests task_01, task_02 and task_03 name in their
   Verification.
6. MUST read the modes reference and the Modes section and record, for each
   of `_prd.md` → Core Features 3, 4 and 5, the sentence that delivers it:
   the read-only `--status` report and what it never does; each `--follow-up`
   step; the `--unattended` authorizations list, option rules, chat fallback
   and answer placement; each optional source's detection and the missing
   source report, and the sentences that make the first word select the
   mode with or without `--`, separate it from the CLI command `roundfix
   follow-up` and keep today's behavior for any other first word. Record
   every command and flag the reference names and that `bin/roundfix --help`
   and that command's own `--help`, each run with `HOME` and `TMPDIR` set to
   fresh temporary directories, list each one. A skill invocation form
   (`/roundfix status`, `$roundfix --status`, `@roundfix unattended` and the
   like) is a Skill Mode, not a command, and is not checked against the
   help; any other named command or flag the binary lacks is a finding.
7. MUST record, as evidence this Spec did not author, each source of
   `_prd.md` → Acceptance evidence:
   - read
     `internal/baseline/assets/source-baselines/baseline.standard-typescript-monorepo-0.0.1/manifest.json`
     in the audited checkout and record the entry for
     `clause.core.ask-user-answerable-decisions` with its `enforcement`, that
     the file is byte-identical to the Delivery Base, and that
     `TestAnAdopterRetainsTheAskUserClause`, whose adopter carries that
     Source Baseline, records the clause `retained`;
   - because this Spec's authorization forbids network access, record the
     Codex pull request "Feat: request user input tool" and the Claude Code
     skills documentation as `blocked (environment: network prohibited by
     Spec authorization)`, quoting the passages the PRD cites. Such a row
     never decides a qualifying `partial`.
8. MUST record the live harness criterion as exactly one row, `blocked
   (declared unreachable)`, against the PRD's single Unreachable Acceptance
   entry, and count it once in `rows_blocked_declared`: no row may start an
   agent harness or ask a person.
9. MUST perform the glossary and documentation checks:
   - record that `CONTEXT.md` defines **Skill Mode** with its ADR-0320 link,
     that the glossary contract tests pass, and whether any other term the
     Spec introduced needs an entry (`docs/agents/domain.md`);
   - search the repository outside `docs/history/` for "exactly one
     question", "one question per call", "single question" and "Never bundle
     questions", and record each match; matches in the owned interview skills
     are expected by ADR-0320, and any other is a finding.
10. MUST verify from the repository history:
    - each Task's changed files stay within its declarations or its `##
      Recorded paths`;
    - every Governed Path changed is bounded in `_authorization.md`;
    - these are byte-identical to the Delivery Base: the Makefile, the CI
      workflows, `go.mod`, `go.sum`, `.roundfixrc.yml`, every Baseline module
      other than `internal/baseline/assets/modules/core.json`, everything
      under `internal/baseline/assets/source-baselines/` and
      `internal/baseline/assets/retention/`, every file under `docs/history/`
      and `docs/user-guide/`, every owned skill other than `roundfix`,
      `.agents/skills/roundfix/agents/`, `docs/references/behavior-surfaces.json`,
      `docs/references/skill-coverage.json`, and every Go file that is not a
      `_test.go` file;
    - in `core.json`, only the clause's guidance, its rule's version and the
      module's version changed;
    - the Roundfix skill's version was raised and recorded;
    - no `### QA settlement` section changed.

    This row reads Task commits, so it declares a `commit_range` input.
11. MUST record the non-waivable Pull Request row as
    `blocked (environment: no open Pull Request)`, naming the Pull Request row
    in its provenance. It supports the row with the pre-PR equivalent of each
    control:
    - approval from the recorded delivery authority;
    - checks from the Daemon's repository Verification at the audited head;
    - no unresolved threads, because no Pull Request exists;
    - no Merge-Ready acceptance yet;
    - the audited head named as the claimed candidate.
12. MUST, for the Outcome for the Archive Record, run only
    `bin/roundfix archive 0307-the-roundfix-skill-modes-and-up-to-three-questions-per-call --plan --offline`
    with every Judge Key Variable and `TYPESAFE_API_KEY` unset and with `HOME`
    and `TMPDIR` each set to a fresh empty temporary directory created for
    this row (the CLI
    loads User Config from `$HOME/.roundfix/config.yml`, so the default home
    would read `~/.roundfix`), and write the
    report's `## Outcome` as one paragraph that does not end with a colon. It
    MUST NOT do any of these:
    - accept a row whose only evidence is that a file was read;
    - read or write `~/.roundfix`;
    - record any environment value;
    - send any request to a provider, the judge or a reviewer;
    - write into the checkout while a test suite runs.

## Subtasks

- [ ] Build the binary and run the repository Verification.
- [ ] Run each named test, the refresh in a temporary copy, and each sabotage on temporary copies of the audited head.
- [ ] Read the modes against the Core Features and the binary's help, and record the outside sources and the declared row.
- [ ] Record the glossary and wording checks and the scope audit, and write the QA report.

## Acceptance Criteria

- [ ] The clause allows three independent questions per call in the catalog and both guides, an adopter keeps it `retained`, and each sabotage makes its named case fail.
- [ ] The skill reads its mode in the three invocation forms, bare and with `--`, its modes reference delivers every Core Feature with commands the binary has, and each sabotage makes its named case fail.
- [ ] The scope audit finds no undeclared Governed Path change and no change to the files row 10 names.

## Verification

The following command is rendered unchanged from Roundfix's derived QA contract.
It proves a newest report exists and its verdict is readable; eligibility is
applied in process by whoever settles the Task.

- `newest="$(find 'docs/specs/0307-the-roundfix-skill-modes-and-up-to-three-questions-per-call/qa' -type f -name "qa-report-*.md" -print 2>/dev/null | awk "{ report=\$0; name=\$0; parts=split(name, path, \"/\"); name=path[parts]; name=substr(name, 11, length(name)-13); date=substr(name, 1, 10); suffix=substr(name, 11); shape=(length(date) == 10 && substr(date, 5, 1) == \"-\" && substr(date, 8, 1) == \"-\"); for (i=1; i <= 10 && shape; i++) { if (i != 5 && i != 8 && index(\"0123456789\", substr(date, i, 1)) == 0) shape=0 } year=substr(date, 1, 4)+0; month=substr(date, 6, 2)+0; day=substr(date, 9, 2)+0; leap=(year%400 == 0 || (year%4 == 0 && year%100 != 0)); days=(month == 2 ? 28+leap : ((month == 4 || month == 6 || month == 9 || month == 11) ? 30 : 31)); dated=(shape && month >= 1 && month <= 12 && day >= 1 && day <= days); if (!dated) date=\"\"; if (suffix == \"\") { sequenced=1; sequence=-1 } else { sequenced=(length(suffix) > 1 && substr(suffix, 1, 1) == \"-\"); for (i=2; i <= length(suffix) && sequenced; i++) { if (index(\"0123456789\", substr(suffix, i, 1)) == 0) sequenced=0 } sequence=(sequenced ? substr(suffix, 2)+0 : 0) } if (dated && sequenced) printf \"%s\\t%s\\t%s\\t%s\\t%s\\n\", dated, date, sequenced, sequence, report }" | sort -k1,1n -k2,2 -k3,3n -k4,4n -k5,5 | tail -1 | cut -f5-)"; test -n "$newest" || exit 1; awk "BEGIN { whitespace=\" \t\r\n\f\v\" } NR == 1 && \$0 == \"---\" { frontmatter=1; next } frontmatter && \$0 == \"---\" { closed=1; exit } frontmatter && index(\$0, \"verdict:\") == 1 { verdict=substr(\$0, 9); while (length(verdict) > 0 && index(whitespace, substr(verdict, 1, 1)) > 0) verdict=substr(verdict, 2); while (length(verdict) > 0 && index(whitespace, substr(verdict, length(verdict), 1)) > 0) verdict=substr(verdict, 1, length(verdict)-1); verdicts++ } END { exit(closed && verdicts == 1 && length(verdict) > 0 ? 0 : 1) }" "$newest"`

## References

- `_prd.md` → Goals; Core Features 1-6; Success Metric 1; Success Metric 2; Success Metric 3; Success Metric 4; Success Metric 5; Acceptance evidence; Unreachable Acceptance; Glossary
- `_techspec.md` → Current behavior; API Contracts 1-7; Invariants 1-7; Testing Approach; Build Order 4
- ADR-0320
- ADR-0330
