---
task: task_05
spec: 0356-a-main-merge-is-a-review-delta-and-a-segment-is-measured-from-its-effective-base
status: pending
type: qa
complexity: medium
---

# Task 05: Run the final QA gate

## Overview

This is the authored terminal gate for this Spec. It declares what the matrix
covers and settles the Spec on evidence. The Spec measures every diff the
pre-PR review takes inside the candidate from its Effective Base, keeps the
Reviewer Lineage across a conflict-free merge of the default branch, names
every restart, prints that restart from `deliver status` and `deliver
retry` (ADR-0366), and states the rules in the guides, the Roundfix skill,
the glossary and four amended ADRs.

Every row runs on this repository's own files, temporary homes, temporary
repositories, the fake reviewers of `internal/cli` and temporary copies of
the audited head. No row starts an Agent session or an ACP adapter, reads or
writes the real `~/.roundfix` or the real home directory, reads another
repository, reaches the network, sends a request to a provider, the judge or a
reviewer, or reads a credential. Every `go test` it runs carries an explicit
`-timeout`. Every row that runs `bin/roundfix` sets `HOME` and `TMPDIR`
to fresh temporary directories, with `GOMODCACHE` and `GOCACHE` kept at
their `go env` values, and removes them afterwards. The QA never writes into
the checkout while a test suite runs. The gate keeps its evidence under names
the repository does not ignore, `.txt` for command output, and links only
files it kept.

## Requirements

1. MUST build the binary with `make build` and run the repository
   Verification, recording its result as a gate fact.
2. MUST verify Success Metrics 1 to 5 and Invariants 1 to 5 by running
   `go test -count=1 -timeout 20m -v` with
   `TestReviewMeasuredFromCarriesOnlyTheCandidatesOwnChange`,
   `TestReviewSegmentsMeasureEachCommitFromItsEffectiveBase`,
   `TestReviewLineageContinuesAcrossAConflictFreeMergeOfTheDefaultBranch`,
   `TestCorrectionReviewAfterAMergeOfTheDefaultBranchReadsOnlyTheFix`,
   `TestReviewLineageRestartNamesItsReason`,
   `TestReviewLineageDecidesEachRound`,
   `TestDeliverStatusNamesAReviewThatOpensANewLineage` and
   `TestDeliverRetryNamesAReviewThatOpensANewLineage` over `./internal/cli`,
   recording every row and subtest by name; then the whole `./internal/cli`
   and `./skills` packages,
   `go test -count=1 -timeout 20m -run '^TestCoverageEquivalence$' ./internal/spec`
   and `-tags docscontract ./internal/docscontract`, recording each pass.
3. MUST record Surface Transcripts 1 and 2 of `_techspec.md` as reproduced
   by `TestDeliverStatusNamesAReviewThatOpensANewLineage` (`rewritten
   head`) and `TestDeliverRetryNamesAReviewThatOpensANewLineage`, which
   drive `deliver status` and `deliver retry` through the command entry
   point with a review record a fake reviewer wrote; a run of the built
   binary would need a real reviewer to write that record, which this gate
   never calls.
4. MUST reproduce the two Fiscus shapes of `_prd.md` → Validation on a
   temporary copy of the Delivery Base made with `git archive` and on a
   temporary copy of the audited head, each with a temporary test file that
   drives `review --base main` through `runCLI` with the fake reviewers, and
   record: on the Delivery Base, the point-2 shape opens round 1 with the full
   candidate diff and the point-4 shape exits `2` with `review diff too
   large`; on the audited head, the first is round 2 without the default
   branch's change and the second exits `0` with its segments. Remove both
   copies afterwards.
5. MUST run each sabotage on its own temporary copy of the audited head made
   with `git archive`, with `-timeout` on every `go test`, and record that
   the named test fails on the named case:
   - `planReviewSegments` measuring from `<start>` again:
     `TestReviewSegmentsMeasureEachCommitFromItsEffectiveBase` fails on
     `codex` and `claude`;
   - `reviewMeasuredFrom` returning `from` when `from` is on the base:
     `TestReviewMeasuredFromCarriesOnlyTheCandidatesOwnChange` fails on
     `start on the base`;
   - the merge-parent check of API Contract 3 skipped:
     `TestReviewLineageRestartNamesItsReason` fails on `merge of another
     branch`;
   - the conflict check skipped: it fails on `conflict`;
   - the delta started from `PreviousHead` when `DeltaBase` is set:
     `TestReviewLineageContinuesAcrossAConflictFreeMergeOfTheDefaultBranch`
     and `TestCorrectionReviewAfterAMergeOfTheDefaultBranchReadsOnlyTheFix`
     fail;
   - `reviewLineageNotice` returning `""` always:
     `TestDeliverStatusNamesAReviewThatOpensANewLineage` fails on
     `rewritten head` and `TestDeliverRetryNamesAReviewThatOpensANewLineage`
     fails.
6. MUST execute the tests tasks 01 to 04 name in their Verification.
7. MUST search both copies of the Roundfix `review` and `deliver`
   references, `docs/user-guide/commands/review.md`,
   `docs/user-guide/commands/deliver.md`, `CONTEXT.md` and ADR-0197,
   ADR-0282, ADR-0286 and ADR-0352, and record that each states the text of
   API Contract 6 for it, that none of the `review` texts still says
   "changed provider or merge base, or a head", and that each of the four
   ADRs carries the amendment naming ADR-0366 with `updated_at` of
   2026-11-05.
8. MUST record, as outside-evidence rows, two sources this Spec did not
   write, each with its date and none blocked:
   - the Fiscus report quoted in
     `references/2026-11-05-fiscus-the-delivery-queue-takes-days-and-a-main-merge-restarts-the-review.md`
     (2026-11-05): compare offline its points 2 and 4 with the results of
     requirement 4;
   - Git's own merge, through the `git` binary: record the `conflicting
     carry` row of `TestReviewMeasuredFromCarriesOnlyTheCandidatesOwnChange`
     and the exit status of `git merge-tree --write-tree` for that pair in a
     temporary repository, and compare them with the git-merge-tree(1)
     passage `_prd.md` → Acceptance evidence quotes.
9. MUST record that the PRD declares no Unreachable Acceptance and count no
   row in `rows_blocked_declared`.
10. MUST perform the glossary check (`docs/agents/domain.md`): record that
    the PRD and TechSpec declare `adds: **Effective Base**`, `changes:
    **Reviewer Lineage**` and `changes: **Review Segment**`, that
    `CONTEXT.md` defines each as API Contract 6 states, and that the
    glossary test of `./internal/docscontract` passes.
11. MUST verify from the repository history:
    - each Task's changed files stay within its declarations or its `##
      Recorded paths`;
    - every changed Governed Path is listed in `_authorization.md` `paths`;
    - the Task commit range changes none of these: the Makefile, the CI
      workflows, `go.mod`, `go.sum`, `.roundfixrc.yml`, `skills-lock.json`,
      `CHANGELOG.md`, `internal/app/version.go`,
      `dist/npm/roundfix/package.json`,
      `docs/references/skill-coverage.json`,
      `docs/references/coverage-record.json`, anything under
      `internal/baseline/`, `internal/delivery/`, `internal/store/`,
      `internal/config/` or `docs/agents/`, any ADR other than the four
      task_04 amends, any skill other than the Roundfix skill files task_04
      declares;
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
    `h="$(mktemp -d)"; t="$(mktemp -d)"; HOME="$h" TMPDIR="$t" GOMODCACHE="$(go env GOMODCACHE)" GOCACHE="$(go env GOCACHE)" bin/roundfix archive 0356-a-main-merge-is-a-review-delta-and-a-segment-is-measured-from-its-effective-base --plan --offline; rm -rf "$h" "$t"`
    with every Judge Key Variable and `TYPESAFE_API_KEY` unset, and write the
    report's Outcome as one paragraph that does not end with a colon. It MUST
    NOT do any of these:
    - accept a row whose only evidence is that a file was read;
    - read or write `~/.roundfix` or the real home directory;
    - record any environment value;
    - send any request to a provider, the judge or a reviewer;
    - write into the checkout while a test suite runs;
    - link a file it did not keep, or keep evidence under a name the
      repository ignores.

## Subtasks

- [ ] Build the binary and run the repository Verification.
- [ ] Run the named tests and reproduce the two Fiscus shapes on the Delivery Base and the audited head.
- [ ] Run each sabotage on a temporary copy.
- [ ] Record the wording search, the outside sources, the glossary check and the scope audit, and write the QA report.

## Acceptance Criteria

- [ ] A candidate that merged the default branch is segmented and reviewed on its own change only, and a conflict-free merge keeps its Reviewer Lineage.
- [ ] Every restart names its reason, and `deliver status` and `deliver retry` print it before the review runs.
- [ ] Each sabotage makes its named test fail, and the scope audit finds no change outside the declared files.

## Verification

- `newest="$(find 'docs/specs/0356-a-main-merge-is-a-review-delta-and-a-segment-is-measured-from-its-effective-base/qa' -type f -name "qa-report-*.md" -print 2>/dev/null | awk "{ report=\$0; name=\$0; parts=split(name, path, \"/\"); name=path[parts]; name=substr(name, 11, length(name)-13); date=substr(name, 1, 10); suffix=substr(name, 11); shape=(length(date) == 10 && substr(date, 5, 1) == \"-\" && substr(date, 8, 1) == \"-\"); for (i=1; i <= 10 && shape; i++) { if (i != 5 && i != 8 && index(\"0123456789\", substr(date, i, 1)) == 0) shape=0 } year=substr(date, 1, 4)+0; month=substr(date, 6, 2)+0; day=substr(date, 9, 2)+0; leap=(year%400 == 0 || (year%4 == 0 && year%100 != 0)); days=(month == 2 ? 28+leap : ((month == 4 || month == 6 || month == 9 || month == 11) ? 30 : 31)); dated=(shape && month >= 1 && month <= 12 && day >= 1 && day <= days); if (!dated) date=\"\"; if (suffix == \"\") { sequenced=1; sequence=-1 } else { sequenced=(length(suffix) > 1 && substr(suffix, 1, 1) == \"-\"); for (i=2; i <= length(suffix) && sequenced; i++) { if (index(\"0123456789\", substr(suffix, i, 1)) == 0) sequenced=0 } sequence=(sequenced ? substr(suffix, 2)+0 : 0) } if (dated && sequenced) printf \"%s\\t%s\\t%s\\t%s\\t%s\\n\", dated, date, sequenced, sequence, report }" | sort -k1,1n -k2,2 -k3,3n -k4,4n -k5,5 | tail -1 | cut -f5-)"; test -n "$newest" || exit 1; awk "BEGIN { whitespace=\" \t\r\n\f\v\" } NR == 1 && \$0 == \"---\" { frontmatter=1; next } frontmatter && \$0 == \"---\" { closed=1; exit } frontmatter && index(\$0, \"verdict:\") == 1 { verdict=substr(\$0, 9); while (length(verdict) > 0 && index(whitespace, substr(verdict, 1, 1)) > 0) verdict=substr(verdict, 2); while (length(verdict) > 0 && index(whitespace, substr(verdict, length(verdict), 1)) > 0) verdict=substr(verdict, 1, length(verdict)-1); verdicts++ } END { exit(closed && verdicts == 1 && length(verdict) > 0 ? 0 : 1) }" "$newest"`

## References

- `_prd.md` → Goals; Core Features 1-5; Success Metric 1; Success Metric 2; Success Metric 3; Success Metric 4; Success Metric 5; Acceptance evidence; Unreachable Acceptance; Glossary
- `_techspec.md` → Current behavior; API Contracts 1-6; Surface Transcript 1; Surface Transcript 2; Invariants 1-5; Declared breaks; Testing Approach; Build Order 5
