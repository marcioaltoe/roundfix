---
task: task_05
spec: 0306-the-follow-up-command-lists-pending-work-in-priority-order
status: completed
type: qa
complexity: high
---

# Task 05: Run the final QA gate

## Overview

This is the authored terminal gate for this Spec. It declares what the matrix
covers and settles the Spec on evidence. The Spec adds the read-only
`roundfix follow-up` command, which lists a repository's Follow-up Items in
Follow-up Priority order beside its work in progress, says when there is
nothing to do, and with `--status` shows the work in progress, the pending
Specs, the next release and what waits on the maintainer.

Every row runs on temporary repositories, homes, Secondbrain Roots and Git
configuration under the Run's temporary directory. No row reads or writes the
real `~/.roundfix` or the real Secondbrain, reaches the network, runs the real
`gh`, sends a request to a provider, the judge or a reviewer, reads a
credential, or reads another repository. The gate never runs `roundfix spec
judge`, never runs `roundfix review` against a real provider, and never runs
an archive plan without `--offline`.

## Requirements

1. MUST build the binary with `make build` and run the repository
   Verification, recording its result as a gate fact.
2. MUST reproduce `_techspec.md` → Surface Transcripts 1, 2 and 3 through
   `bin/roundfix` from an empty temporary directory, comparing standard
   output, standard error and the exit code line by line, and record each
   result.
3. MUST run the built product end to end once in a temporary Git repository
   without remotes, with `HOME` set to a temporary home whose User Config
   names a temporary Secondbrain Root, every Judge Key Variable and
   `TYPESAFE_API_KEY` unset, and a first `PATH` entry holding a temporary
   `gh` script that prints `gh unavailable in QA` to standard error and exits
   1. Seed two pending Inbox Entries, one triaged entry, one open Backlog
   Entry, one pending Finding and one active Spec, and record from
   `bin/roundfix follow-up --format json` and the text view: the five items,
   their order by type, the judge line naming the missing key, the
   `pull-requests` source `unavailable` with `gh unavailable in QA`, and exit
   0. Then remove every seed, commit, and record that the text view prints
   `Nothing to do: no pending item and nothing in progress.` and JSON `state`
   `nothing_to_do`, and that `git status --porcelain` and the Secondbrain
   Root are unchanged by both runs.
4. MUST verify Success Metrics 1-5 and 7 by running, in `./internal/cli`,
   `TestFollowUpListsEveryPendingItemOnce`,
   `TestFollowUpOrdersItemsByTheFourAnswers`,
   `TestFollowUpOrdersByTypeWithoutTheJudge`, `TestFollowUpSaysNothingToDo`,
   `TestFollowUpReportsAMissingSecondbrainRoot`,
   `TestFollowUpNamespaceDefaultsToTheRepositoryName`,
   `TestFollowUpReadsNoUnavailableSourceAsAFailure`,
   `TestFollowUpWritesNothing` and `TestFollowUpRefusesAnInvalidMessagesFile`;
   in `./internal/judge`, the seven `TestFollowUp*` tests task_02 names; and
   in `./internal/config`, the four Secondbrain tests task_03 names. Record
   for each the cases it ran.
5. MUST verify Success Metric 6 and the status view's boundary by running
   `TestFollowUpStatusListsPendingSpecsInDeliveryOrder`,
   `TestFollowUpStatusReportsTheNextRelease`,
   `TestFollowUpStatusListsBlockersAndTheMaintainersQuestion`,
   `TestFollowUpStatusNeverAsksTheJudge` and `TestFollowUpStatusWritesNothing`,
   and record the order, the release lines and the maintainer entries each
   asserted.
6. MUST, on temporary copies of the audited head made with `git archive`,
   apply one sabotage per copy, run the named test, and record that it fails
   on the named case:
   - make the yes verdict need a Noul strictly above `yes_at_least`:
     `TestFollowUpVerdictsUseTheThresholds` fails on its exact `0.7` answer;
   - let an `unsure` blocks answer count as yes in `followUpOrder`:
     `TestFollowUpOrdersItemsByTheFourAnswers` fails on the item whose blocks
     answer is `0.5` and fix answer `0.9`;
   - remove the secret scan before sending:
     `TestFollowUpNeverSendsASecretShapedItem` fails, because a request
     carries the secret-shaped item;
   - cut the item text by bytes instead of runes:
     `TestFollowUpCutsTextOnARuneBoundary` fails on the text of 3,999 ASCII
     characters followed by `ção`;
   - let Project Config set `secondbrain.root`:
     `TestProjectConfigCannotSetTheSecondbrainRoot` fails;
   - read the inbox namespace recursively, `_triaged/` included:
     `TestFollowUpListsEveryPendingItemOnce` fails on its triaged entry;
   - take an Inbox Entry's date from its file name instead of `created_at`:
     `TestFollowUpListsEveryPendingItemOnce` fails, because its two entries'
     `created_at` order is the reverse of their file-name order;
   - ignore the work in progress in the stop rule:
     `TestFollowUpSaysNothingToDo` fails on its case with one active Run;
   - order `pending_specs` by slug only:
     `TestFollowUpStatusListsPendingSpecsInDeliveryOrder` fails, because
     `0009-second` is queued before `0003-first`;
   - drop the `ErrInvalidReleaseRange` mapping:
     `TestFollowUpStatusReportsTheNextRelease` fails on its subtest at the
     release tag;
   - call the judge from the status view:
     `TestFollowUpStatusNeverAsksTheJudge` fails.
7. MUST execute the tests task_01, task_02, task_03 and task_04 name in
   their Verification.
8. MUST record, as evidence this Spec did not author, each source of
   `_prd.md` → Acceptance evidence:
   - quote, from `docs/agents/secondbrain.md` at the audited head, the two
     passages the PRD cites, and record that the inbox reader matches them
     (positional status, `_triaged/` excluded, oldest first by `created_at`);
   - because this Spec's authorization forbids network access, record
     TypeSafe's Noul and primitives pages and Rex Black's triage framework as
     `blocked (environment: network prohibited by Spec authorization)`,
     quoting the passages the PRD cites. Such rows never decide a qualifying
     `partial`.
9. MUST record Success Metric 8 as exactly one row, `blocked (declared
   unreachable)`, against the PRD's single Unreachable Acceptance entry, and
   count it once in `rows_blocked_declared`.
10. MUST perform the glossary and documentation checks:
    - record that `CONTEXT.md` defines **Follow-up Command**, **Follow-up
      Item**, **Follow-up Priority** and **Secondbrain Root**, each linking
      ADR-0319 where its definition names it; record that the glossary
      contract tests pass and whether any other term the Spec introduced
      needs an entry (`docs/agents/domain.md`);
    - record that `docs/user-guide/commands/follow-up.md`, both copies of
      the `follow-up` reference and the configuration guide describe the
      sources, the order, the stop rule, both views and both keys as the
      built product behaves in Requirement 3, and that
      `TestTheSkillCoverageMapIsCurrent` passes;
    - record that the reference index names the `follow-up` family once and
      that no reference describes the skill modes Spec 0307 owns.
11. MUST verify from the repository history:
    - each Task's changed files stay within its declarations or its
      `## Recorded paths`;
    - every Governed Path changed is bounded in `_authorization.md`;
    - these are byte-identical to the Delivery Base: the Makefile, the CI
      workflows, `go.mod`, `go.sum`, `.roundfixrc.yml`, every file under
      `docs/agents/`, `internal/baseline/`, `internal/store/`,
      `internal/delivery/` and `internal/releaseplan/`, every file under
      `docs/history/`, `internal/cli/releaseplan_command.go`,
      `internal/cli/deliver.go`, `internal/cli/runs.go` and
      `internal/cli/spec_judge.go`;
    - in `internal/judge/questions.json` only the `follow-up-priority`
      judgment was added;
    - the Roundfix skill version was raised and recorded;
    - no `### QA settlement` section changed.

    This row reads Task commits, so it declares a `commit_range` input.
12. MUST record the non-waivable Pull Request row as
    `blocked (environment: no open Pull Request)`, naming the Pull Request
    row in its provenance. It supports the row with the pre-PR equivalent of
    each control:
    - approval from the recorded delivery authority;
    - checks from the Daemon's repository Verification at the audited head;
    - no unresolved threads, because no Pull Request exists;
    - no Merge-Ready acceptance yet;
    - the audited head named as the claimed candidate.
13. MUST, for the Outcome for the Archive Record, run only
    `bin/roundfix archive 0306-the-follow-up-command-lists-pending-work-in-priority-order --plan --offline`
    with every Judge Key Variable and `TYPESAFE_API_KEY` unset and with `HOME`
    set to a fresh empty temporary directory created for this row (the CLI
    loads User Config from `$HOME/.roundfix/config.yml`, so the default home
    would read `~/.roundfix`), and write the
    report's `## Outcome` as one paragraph that does not end with a colon. It
    MUST NOT do any of these:
    - accept a row whose only evidence is that a file was read;
    - read or write `~/.roundfix` or the real Secondbrain;
    - record any environment value;
    - send any request to a provider, the judge or a reviewer;
    - write into the checkout while a test suite runs.

## Subtasks

- [ ] Build the binary, run the repository Verification and reproduce the three Surface Transcripts.
- [ ] Run the built product end to end in temporary repositories and homes.
- [ ] Run each named test and each sabotage on temporary copies of the audited head.
- [ ] Record the outside evidence, the declared row, the glossary and documentation checks and the scope audit, and write the QA report.

## Acceptance Criteria

- [ ] The built product lists the seeded items, orders them by type without a key, reports the missing sources and nothing to do, and writes nothing.
- [ ] Each sabotage makes its named test fail on the named case.
- [ ] The scope audit finds no undeclared Governed Path change and no change outside the declared files.

## Verification

The following command is rendered unchanged from Roundfix's derived QA contract.
It proves a newest report exists and its verdict is readable; eligibility is
applied in process by whoever settles the Task.

- `newest="$(find 'docs/specs/0306-the-follow-up-command-lists-pending-work-in-priority-order/qa' -type f -name "qa-report-*.md" -print 2>/dev/null | awk "{ report=\$0; name=\$0; parts=split(name, path, \"/\"); name=path[parts]; name=substr(name, 11, length(name)-13); date=substr(name, 1, 10); suffix=substr(name, 11); shape=(length(date) == 10 && substr(date, 5, 1) == \"-\" && substr(date, 8, 1) == \"-\"); for (i=1; i <= 10 && shape; i++) { if (i != 5 && i != 8 && index(\"0123456789\", substr(date, i, 1)) == 0) shape=0 } year=substr(date, 1, 4)+0; month=substr(date, 6, 2)+0; day=substr(date, 9, 2)+0; leap=(year%400 == 0 || (year%4 == 0 && year%100 != 0)); days=(month == 2 ? 28+leap : ((month == 4 || month == 6 || month == 9 || month == 11) ? 30 : 31)); dated=(shape && month >= 1 && month <= 12 && day >= 1 && day <= days); if (!dated) date=\"\"; if (suffix == \"\") { sequenced=1; sequence=-1 } else { sequenced=(length(suffix) > 1 && substr(suffix, 1, 1) == \"-\"); for (i=2; i <= length(suffix) && sequenced; i++) { if (index(\"0123456789\", substr(suffix, i, 1)) == 0) sequenced=0 } sequence=(sequenced ? substr(suffix, 2)+0 : 0) } if (dated && sequenced) printf \"%s\\t%s\\t%s\\t%s\\t%s\\n\", dated, date, sequenced, sequence, report }" | sort -k1,1n -k2,2 -k3,3n -k4,4n -k5,5 | tail -1 | cut -f5-)"; test -n "$newest" || exit 1; awk "BEGIN { whitespace=\" \t\r\n\f\v\" } NR == 1 && \$0 == \"---\" { frontmatter=1; next } frontmatter && \$0 == \"---\" { closed=1; exit } frontmatter && index(\$0, \"verdict:\") == 1 { verdict=substr(\$0, 9); while (length(verdict) > 0 && index(whitespace, substr(verdict, 1, 1)) > 0) verdict=substr(verdict, 2); while (length(verdict) > 0 && index(whitespace, substr(verdict, length(verdict), 1)) > 0) verdict=substr(verdict, 1, length(verdict)-1); verdicts++ } END { exit(closed && verdicts == 1 && length(verdict) > 0 ? 0 : 1) }" "$newest"`

## References

- `_prd.md` → Goals; User Stories 1-4; Core Features 1-8; Success Metric 1; Success Metric 2; Success Metric 3; Success Metric 4; Success Metric 5; Success Metric 6; Success Metric 7; Success Metric 8; Acceptance evidence; Unreachable Acceptance; Glossary
- `_techspec.md` → API Contracts 1-9; Surface Transcripts 1-3; Invariants 1-7; Testing Approach; Build Order 5
- ADR-0095
- ADR-0143
- ADR-0209
- ADR-0256
- ADR-0319
