---
task: task_03
spec: 0290-the-pull-request-record-commit-fits-commitlint
status: pending
type: qa
complexity: high
---

# Task 03: Run the final QA gate

## Overview

This is the authored terminal gate for this Spec. It declares what the matrix
covers and settles the Spec on evidence. The Spec makes the Delivery Queue's
Pull Request Record Commit name the Spec by its number, so its subject fits
100 characters for every slug, keeps adopting a record commit made by hand
under any subject, and makes a commit the queue makes that a hook refuses park
with a reason naming the hook, the subject and the hand-commit recovery.

Every row runs on temporary repositories, homes and Git configuration under
the Run's temporary directory. No row reads or writes the real `~/.roundfix`,
reaches the network, sends a request to a provider, the judge or a reviewer,
reads a credential, or reads another repository. The gate never runs
`roundfix spec judge`, never runs `roundfix review` against a real provider,
never runs `roundfix deliver` against a real forge, and never runs an archive
plan without `--offline`.

## Requirements

1. MUST build the binary with `make build` and run the repository
   Verification, recording its result as a gate fact.
2. MUST verify Success Metrics 1, 2, 3 and 5:
   - Run `TestPullRequestRecordSubjectNamesTheSpecNumber`,
     `TestDeliveryWorkflowRecordCommitPassesAHeaderLimitHook` and
     `TestDeliveryWorkflowAdoptsARecordCommitWhateverItsSubject`, and record
     each subtest's subject and its length (Invariants 1 and 2).
   - Run the record subtests of
     `TestDeliveryQueueCommitHookRefusalNamesTheHookAndTheRecovery` and
     record that each adopts the hand commit made after the refusal.
3. MUST verify Success Metric 4:
   - Run `TestDeliveryQueueCommitHookRefusalNamesTheHookAndTheRecovery` and
     `TestDeliveryQueueCommitWithoutAHookKeepsTheGitError`, and record, for
     each subtest, the hook label, the refused subject, the exit status and
     the recovery sentence in the error, and the Park Class (Invariants 3-5).
   - Run `TestClassifyCommitHookRefusal`,
     `TestClassifyCommitHookRefusalReadsTheRepositoryWhenOutputNamesNothing`,
     `TestGitCommitterClassifiesHookRefusalAndLeavesWorkStaged`,
     `TestGitCommitterKeepsPlainGitFailureUnclassified` and
     `TestTaskCycleHookRefusalKeepsTaskCompletedAndNamesRecovery` in
     `./internal/daemon`, and record that the Daemon's Task commit refusal is
     unchanged (Invariant 6).
   - On temporary copies of the audited head made with `git archive`, apply
     one sabotage per copy and record that the named test fails:
     - make `pullRequestRecordSubject` return the slug form again:
       `TestPullRequestRecordSubjectNamesTheSpecNumber` and
       `TestDeliveryWorkflowRecordCommitPassesAHeaderLimitHook` fail;
     - make the queue commit helper return the runner's error unclassified:
       `TestDeliveryQueueCommitHookRefusalNamesTheHookAndTheRecovery` fails;
     - make the helper classify on the runner's whole error text, arguments
       included: the subtest whose slug contains `pre-commit` fails;
     - make `pullRequestRecordCommitIsExact` also require the queue's own
       subject: `TestDeliveryWorkflowAdoptsARecordCommitWhateverItsSubject`
       fails.
4. MUST execute the tests task_01 and task_02 name in their Verification.
5. MUST record, as evidence this Spec did not author, each source of
   `_prd.md` → Acceptance evidence:
   - quote from the adopted Inbox Entry under `references/` the park reason
     it recorded, the 107-character subject and the hand-commit subject; then,
     in a temporary repository with an executable `commit-msg` hook that
     exits 1 when the message's first line exceeds 100 characters, record
     that Git refuses a commit with that 107-character subject, leaves the
     change staged, and commits it with the 57-character subject;
   - run `git log -1 --format=%s 9e92ac39` in the audited checkout and record
     the subject and its length, 117 characters;
   - because this Spec's authorization forbids network access, record Git's
     hook documentation, the commitlint conventional configuration and the
     Stack Overflow answer as `blocked (environment: network prohibited by
     Spec authorization)`, quoting the passages the PRD cites. Such a row
     never decides a qualifying `partial`.
6. MUST record Success Metric 6 as exactly one row, `blocked (declared
   unreachable)`, against the PRD's single Unreachable Acceptance entry, and
   count it once in `rows_blocked_declared`: no row may read another
   repository or the real `~/.roundfix`.
7. MUST perform the glossary and documentation checks:
   - record that `CONTEXT.md` revises **Hook Refusal** to cover a Delivery
     Queue commit and links ADR-0302, that the glossary contract tests pass,
     and whether any other term the Spec introduced needs an entry
     (`docs/agents/domain.md`);
   - record the two paragraphs `docs/user-guide/commands/deliver.md` and both
     copies of the Roundfix `deliver` reference gained, and that
     `TestTheSkillCoverageMapIsCurrent` passes;
   - search the guides, the agent guides and the skill copies for wording
     that gives the record subject with the slug or says a refused queue
     commit parks with a Git error only, and record what it finds.
8. MUST verify from the repository history:
   - each Task's changed files stay within its declarations or its `##
     Recorded paths`;
   - every Governed Path changed is bounded in `_authorization.md`;
   - these are byte-identical to the Delivery Base: the Makefile, the CI
     workflows, `go.mod`, `go.sum`, `.roundfixrc.yml`, every Baseline module
     under `internal/baseline/assets/`, `docs/agents/setup-context.json`,
     every file under `docs/history/`, every file under `internal/delivery/`
     and `internal/preflight/`, every file under `internal/daemon/` other than
     `commit_hook.go`, and `docs/references/skill-coverage.json`;
   - every `_test.go` file under `internal/cli` that exists at the Delivery
     Base is byte-identical to it, except
     `internal/cli/deliver_record_pull_request_test.go`, whose diff from the
     Delivery Base changes only the expected subject of
     `TestDeliveryWorkflowRecordsThePullRequestInTheArchiveRecord`;
   - the Roundfix skill version was raised and recorded;
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
    `bin/roundfix archive 0290-the-pull-request-record-commit-fits-commitlint --plan --offline`
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
- [ ] Reproduce the Inbox Entry's refusal and the hand commit with real Git, and read the squash subject from history.
- [ ] Record the outside evidence, the declared row, the glossary and documentation checks and the scope audit, and write the QA report.

## Acceptance Criteria

- [ ] The record subject fits 100 characters behind a header-limit hook, a hand commit under any subject is adopted, and each removed guard makes its named test fail.
- [ ] A refused queue commit names the hook, the subject and the recovery, with class `environment`, and a failure without a hook keeps today's error.
- [ ] The scope audit finds no undeclared Governed Path change and no change to the Delivery Engine or the Daemon's Task commit refusal.

## Verification


The following command is rendered unchanged from Roundfix's derived QA contract.
It proves a newest report exists and its verdict is readable; eligibility is
applied in process by whoever settles the Task.

- `newest="$(find 'docs/specs/0290-the-pull-request-record-commit-fits-commitlint/qa' -type f -name "qa-report-*.md" -print 2>/dev/null | awk "{ report=\$0; name=\$0; parts=split(name, path, \"/\"); name=path[parts]; name=substr(name, 11, length(name)-13); date=substr(name, 1, 10); suffix=substr(name, 11); shape=(length(date) == 10 && substr(date, 5, 1) == \"-\" && substr(date, 8, 1) == \"-\"); for (i=1; i <= 10 && shape; i++) { if (i != 5 && i != 8 && index(\"0123456789\", substr(date, i, 1)) == 0) shape=0 } year=substr(date, 1, 4)+0; month=substr(date, 6, 2)+0; day=substr(date, 9, 2)+0; leap=(year%400 == 0 || (year%4 == 0 && year%100 != 0)); days=(month == 2 ? 28+leap : ((month == 4 || month == 6 || month == 9 || month == 11) ? 30 : 31)); dated=(shape && month >= 1 && month <= 12 && day >= 1 && day <= days); if (!dated) date=\"\"; if (suffix == \"\") { sequenced=1; sequence=-1 } else { sequenced=(length(suffix) > 1 && substr(suffix, 1, 1) == \"-\"); for (i=2; i <= length(suffix) && sequenced; i++) { if (index(\"0123456789\", substr(suffix, i, 1)) == 0) sequenced=0 } sequence=(sequenced ? substr(suffix, 2)+0 : 0) } if (dated && sequenced) printf \"%s\\t%s\\t%s\\t%s\\t%s\\n\", dated, date, sequenced, sequence, report }" | sort -k1,1n -k2,2 -k3,3n -k4,4n -k5,5 | tail -1 | cut -f5-)"; test -n "$newest" || exit 1; awk "BEGIN { whitespace=\" \t\r\n\f\v\" } NR == 1 && \$0 == \"---\" { frontmatter=1; next } frontmatter && \$0 == \"---\" { closed=1; exit } frontmatter && index(\$0, \"verdict:\") == 1 { verdict=substr(\$0, 9); while (length(verdict) > 0 && index(whitespace, substr(verdict, 1, 1)) > 0) verdict=substr(verdict, 2); while (length(verdict) > 0 && index(whitespace, substr(verdict, length(verdict), 1)) > 0) verdict=substr(verdict, 1, length(verdict)-1); verdicts++ } END { exit(closed && verdicts == 1 && length(verdict) > 0 ? 0 : 1) }" "$newest"`

## References

- `_prd.md` → Goals; Core Features 1-4; Success Metric 1; Success Metric 2; Success Metric 3; Success Metric 4; Success Metric 5; Success Metric 6; Acceptance evidence; Unreachable Acceptance; Glossary
- `_techspec.md` → Current behavior; API Contracts 1-7; Invariants 1-6; Testing Approach; Build Order 3
- ADR-0302
