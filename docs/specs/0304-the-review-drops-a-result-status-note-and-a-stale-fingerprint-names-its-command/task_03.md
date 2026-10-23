---
task: task_03
spec: 0304-the-review-drops-a-result-status-note-and-a-stale-fingerprint-names-its-command
status: pending
type: qa
complexity: high
---

# Task 03: Run the final QA gate

## Overview

This is the authored terminal gate for this Spec. It declares what the matrix
covers and settles the Spec on evidence. The Spec rewords Delivery Convention
C2 as version 4, makes the pre-PR review dismiss a finding that names a
Result Status Note against a `completed` non-QA Task by validation, without
the sealed validator, keeps every other finding on the validator's
fail-closed path, and makes a stale Behavior Surface fingerprint problem name
the command that re-records it.

Every row runs on temporary repositories under the Run's temporary directory,
with the fake reviewer and fake sealed validator of the `internal/cli` tests.
No row reads or writes the real `~/.roundfix`, runs `roundfix review` against
a real provider, sends a request to a provider, the judge or a reviewer,
reaches the network, reads a credential, or reads another repository. The
gate never runs `roundfix spec judge` and never runs an archive plan without
`--offline`.

## Requirements

1. MUST build the binary with `make build` and run the repository
   Verification, recording its result as a gate fact.
2. MUST verify Success Metrics 1, 2 and 3:
   - Run `TestReviewDismissesAResultStatusNoteOfACompletedTaskWithoutTheValidator`
     in `./internal/cli` and record, per subtest (`0297 shape`, `0290 shape`,
     `front matter anchor`), the exit code, outcome, `validator` value, the
     number of sealed requests, the rule, the reason and the stderr line.
   - Run `TestReviewKeepsAContradictionThatIsNotAResultStatusNote` and
     record, per subtest, the number of questions, the first eligible rule,
     the kept verdict and the exit code.
   - Run `TestReviewValidatorFailureLeavesEveryAskedFindingStanding`,
     `TestReviewValidatorCannotDismissOutsideTheConventionRegion`,
     `TestReviewValidatorNeverAsksNoFailureForAFindingWithAFailureClause`,
     `TestReviewConventionRegionsAdmitAndRefuseAtTheirEdges`,
     `TestReviewConventionC2IncludesEverySettlementSection` and
     `TestReviewDisposeRefusesAConventionDismissal`, and record that the
     validator's fail-closed handling, the convention regions and `review
     dispose` are unchanged.
   - On temporary copies of the audited head made with `git archive`, apply
     one sabotage per copy to `reviewResultStatusNote` or
     `validateReviewConventions` and record that the named case fails:
     - make the front matter check accept any `status:` value: the `failed
       Task` subtest fails;
     - make it ignore `type: qa`: the `QA Task` subtest fails;
     - scan every section after the front matter instead of `## Result`
       only: the `note outside the Result` subtest fails;
     - drop the requirement of a Result Status Note: the `Result without a
       note` and `note outside the Result` subtests fail;
     - drop the check of the finding's text: the `finding names another
       failure` subtest fails;
     - drop the requirement that the conventions hold `C2`: the `anchor
       outside the settlement region` subtest fails;
     - return exit 1 when no finding is left to ask, as before: every subtest
       of `TestReviewDismissesAResultStatusNoteOfACompletedTaskWithoutTheValidator`
       fails.
3. MUST verify Success Metric 4: run `TestDeliveryConventionC2ExemptsTheQATask`
   and `TestReviewValidatorDismissesADaemonSettlementAsC2` in
   `./internal/cli` and `TestTheDeliveryConventionsAreDocumentedAsTheReviewerPromptsThem`
   and `TestTheDeliveryConventionContractReportsAStaleConvention` with
   `-tags docscontract` in `./internal/docscontract`; record the version and
   C2 text the prompt carries. On a temporary copy, restore the version 3 C2
   bullet in `skills/roundfix/references/review.md` only and record that
   `TestTheDeliveryConventionsAreDocumentedAsTheReviewerPromptsThem` fails
   naming that file.
4. MUST verify Success Metric 5: run `TestAStaleBehaviorSurfaceIsReported`,
   `TestTheStaleFingerprintProblemNamesTheProjectRegeneration` and
   `TestTheSkillCoverageMapIsCurrent` with `-tags docscontract` and record
   the stale, missing and extra problem texts. On temporary copies, one per
   sabotage, record that the named case fails:
   - drop the hint from the stale problem:
     `TestAStaleBehaviorSurfaceIsReported/stale` fails;
   - add the hint to the missing problem too:
     `TestAStaleBehaviorSurfaceIsReported/missing` fails;
   - append ` -count=1` to `behaviorSurfaceRecordCommand`:
     `TestTheStaleFingerprintProblemNamesTheProjectRegeneration` fails.
5. MUST execute the tests task_01 and task_02 name in their Verification.
6. MUST record, as evidence this Spec did not author, each source of
   `_prd.md` → Acceptance evidence:
   - read the `source_revision` of the Archive Records
     `docs/history/specs/0297-a-corrective-task-after-a-completed-gate-passes-the-wave-collision-check.md`
     and `docs/history/specs/0290-the-pull-request-record-commit-fits-commitlint.md`
     in the audited checkout; with `git show <source_revision>:<path>`, read
     that Spec's `task_02.md` (0297) and `task_01.md` (0290) and record their
     `status:` line and the Result lines the PRD cites; then record that,
     with whitespace runs collapsed to one space, the 0297 paragraph and the
     0290 sentence are substrings of the `0297 shape` and `0290 shape`
     fixtures in `internal/cli/review_result_status_note_test.go`. When a
     revision's objects are absent from the checkout, record that source as
     `blocked (environment: source revision absent from the checkout)`,
     naming the revision;
   - because this Spec's authorization forbids network access, record the
     Anthropic cookbook "Building Evals" as `blocked (environment: network
     prohibited by Spec authorization)`, quoting the passage the PRD cites.
     Such a row never decides a qualifying `partial`.
7. MUST perform the glossary and documentation checks:
   - record that `CONTEXT.md` defines **Result Status Note**, that
     **Delivery Convention** gained the dismissal sentence, that both link
     ADR-0317, that the glossary contract tests pass, and whether any other
     term the Spec introduced needs an entry (`docs/agents/domain.md`);
   - record the version, the C2 bullet, the dismissal paragraph and the
     validator sentence in `docs/user-guide/commands/review.md` and both
     copies of the Roundfix review reference, and that
     `TestTheSkillCoverageMapIsCurrent` passes;
   - search the guides, the agent guides and the skill copies for wording
     that says the Daemon writes a Task's `## Result`, that
     `roundfix/delivery-conventions/v3` is current, or that only the sealed
     validator can dismiss an anchored finding, and record what it finds.
8. MUST verify from the repository history:
   - each Task's changed files stay within its declarations or its `##
     Recorded paths`;
   - every Governed Path changed is bounded in `_authorization.md`;
   - these are byte-identical to the Delivery Base: the Makefile, the CI
     workflows, `go.mod`, `go.sum`, `.roundfixrc.yml`, every Baseline module
     under `internal/baseline/assets/`, `docs/agents/setup-context.json`,
     every file under `docs/history/`, `internal/agent/spec_prompt.go`, every
     file under `.agents/skills/implement-task/` and `skills/implement-task/`,
     `internal/cli/review.go`, `internal/cli/review_validation.go`,
     `docs/references/skill-coverage.json`, and every ADR other than
     ADR-0317;
   - the Roundfix skill's version was raised and recorded;
   - no `### QA settlement` section changed.

   This row reads Task commits, so it declares a `commit_range` input.
9. MUST record the non-waivable Pull Request row as
   `blocked (environment: no open Pull Request)`, naming the Pull Request row
   in its provenance. It supports the row with the pre-PR equivalent of each
   control:
   - approval from the recorded delivery authority;
   - checks from the Daemon's repository Verification at the audited head;
   - no unresolved threads, because no Pull Request exists;
   - no Merge-Ready acceptance yet;
   - the audited head named as the claimed candidate.
10. MUST, for the Outcome for the Archive Record, run only
    `bin/roundfix archive 0304-the-review-drops-a-result-status-note-and-a-stale-fingerprint-names-its-command --plan --offline`
    with every Judge Key Variable and `TYPESAFE_API_KEY` unset, and write the
    report's `## Outcome` as one paragraph that does not end with a colon. It
    MUST NOT do any of these:
    - accept a row whose only evidence is that a file was read;
    - read or write `~/.roundfix`;
    - record any environment value;
    - send any request to a provider, the judge or a reviewer;
    - write into the checkout while a test suite runs.

## Subtasks

- [ ] Build the binary and run the repository Verification.
- [ ] Run each named test and each sabotage on temporary copies of the audited head.
- [ ] Compare the delivered Task files of Specs 0297 and 0290 with the fixtures, and record the blocked outside source.
- [ ] Record the glossary and documentation checks and the scope audit, and write the QA report.

## Acceptance Criteria

- [ ] The pair is dismissed without the validator in all three shapes, every other case reaches the validator or stands, and each removed guard makes its named case fail.
- [ ] The prompt, the validator input and the documents carry version 4, and the stale fingerprint problem names the `.roundfixrc.yml` command.
- [ ] The scope audit finds no undeclared Governed Path change and no change to the files row 8 names.

## Verification

The following command is rendered unchanged from Roundfix's derived QA contract.
It proves a newest report exists and its verdict is readable; eligibility is
applied in process by whoever settles the Task.

- `newest="$(find 'docs/specs/0304-the-review-drops-a-result-status-note-and-a-stale-fingerprint-names-its-command/qa' -type f -name "qa-report-*.md" -print 2>/dev/null | awk "{ report=\$0; name=\$0; parts=split(name, path, \"/\"); name=path[parts]; name=substr(name, 11, length(name)-13); date=substr(name, 1, 10); suffix=substr(name, 11); shape=(length(date) == 10 && substr(date, 5, 1) == \"-\" && substr(date, 8, 1) == \"-\"); for (i=1; i <= 10 && shape; i++) { if (i != 5 && i != 8 && index(\"0123456789\", substr(date, i, 1)) == 0) shape=0 } year=substr(date, 1, 4)+0; month=substr(date, 6, 2)+0; day=substr(date, 9, 2)+0; leap=(year%400 == 0 || (year%4 == 0 && year%100 != 0)); days=(month == 2 ? 28+leap : ((month == 4 || month == 6 || month == 9 || month == 11) ? 30 : 31)); dated=(shape && month >= 1 && month <= 12 && day >= 1 && day <= days); if (!dated) date=\"\"; if (suffix == \"\") { sequenced=1; sequence=-1 } else { sequenced=(length(suffix) > 1 && substr(suffix, 1, 1) == \"-\"); for (i=2; i <= length(suffix) && sequenced; i++) { if (index(\"0123456789\", substr(suffix, i, 1)) == 0) sequenced=0 } sequence=(sequenced ? substr(suffix, 2)+0 : 0) } if (dated && sequenced) printf \"%s\\t%s\\t%s\\t%s\\t%s\\n\", dated, date, sequenced, sequence, report }" | sort -k1,1n -k2,2 -k3,3n -k4,4n -k5,5 | tail -1 | cut -f5-)"; test -n "$newest" || exit 1; awk "BEGIN { whitespace=\" \t\r\n\f\v\" } NR == 1 && \$0 == \"---\" { frontmatter=1; next } frontmatter && \$0 == \"---\" { closed=1; exit } frontmatter && index(\$0, \"verdict:\") == 1 { verdict=substr(\$0, 9); while (length(verdict) > 0 && index(whitespace, substr(verdict, 1, 1)) > 0) verdict=substr(verdict, 2); while (length(verdict) > 0 && index(whitespace, substr(verdict, length(verdict), 1)) > 0) verdict=substr(verdict, 1, length(verdict)-1); verdicts++ } END { exit(closed && verdicts == 1 && length(verdict) > 0 ? 0 : 1) }" "$newest"`

## References

- `_prd.md` → Goals; Core Features 1-4; Success Metric 1; Success Metric 2; Success Metric 3; Success Metric 4; Success Metric 5; Acceptance evidence; Glossary
- `_techspec.md` → Current behavior; API Contracts 1-7; Invariants 1-6; Testing Approach; Build Order 3
- ADR-0317
