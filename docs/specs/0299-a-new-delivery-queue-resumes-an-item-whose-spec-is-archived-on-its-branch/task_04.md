---
task: task_04
spec: 0299-a-new-delivery-queue-resumes-an-item-whose-spec-is-archived-on-its-branch
status: pending
type: qa
complexity: high
---

# Task 04: Run the final QA gate

## Overview

This is the authored terminal gate for this Spec. It declares what the matrix
covers and settles the Spec on evidence. The Spec makes a new Delivery Queue
resume, at review, an item whose continued item branch already carries an
exact archive of its Spec, and park with its cause an item branch whose
archive cannot be proven.

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
2. MUST verify Success Metrics 1, 2 and 5:
   - Run `TestDeliverStartResumesAnArchivedItemBranchAtReview` and
     `TestDeliverStartRevalidatesAContinuedBranchWithoutAnArchive` in
     `./internal/cli`, and record for each subtest the final stage, the Pull
     Request number, the number of Implement calls, the reviewed heads, the
     number of Pull Requests created and whether the item head moved
     (Invariants 1-3).
   - Run `TestANewQueueResumesAProvenArchivedItemAtReview` and
     `TestANewQueueRevalidatesAnItemWithoutAnArchive` in `./internal/delivery`,
     and record the stages each passed through.
3. MUST verify Success Metrics 3 and 4:
   - Run `TestDeliverStartParksAnUnprovenArchivedItemBranch`,
     `TestANewQueueParksAnUnprovenArchivedItem`,
     `TestANewQueueParksAnArchivedItemProofError`,
     `TestARetryReturnsAnUnprovenArchivedItemToQueued` and
     `TestParkClassesAndNextActions`, and record each blocker, the `Park:`
     line `roundfix deliver status` printed, and the stage and retry count
     after the retry (Invariants 4 and 6).
   - On temporary copies of the audited head made with `git archive`, apply
     one sabotage per copy, run the named test, and record that it fails on
     the named case:
     - make the `queued` branch of `advanceItem` ignore a proven archive and
       call `Revalidate`: `TestANewQueueResumesAProvenArchivedItemAtReview`
       and `TestDeliverStartResumesAnArchivedItemBranchAtReview/another_active_Spec`
       fail, the latter parking `delivery-error`;
     - make a proven archive enter `gating` instead of `reviewing`:
       `TestANewQueueResumesAProvenArchivedItemAtReview` and
       `TestDeliverStartResumesAnArchivedItemBranchAtReview/another_active_Spec`
       fail, because no review sees the head;
     - make `ProveArchivedItem` accept a removing commit without
       `archiveCommitIsExact`:
       `TestDeliverStartParksAnUnprovenArchivedItemBranch/inexact_archive_commit`
       fails, because the item reaches `merged`;
     - make `archivePaths` resolve the root with `roundconfig.ResolveSpecsRoot`
       again: `TestDeliverStartResumesAnArchivedItemBranchAtReview/only_active_Spec`
       fails, because the worktree has no `docs/specs`;
     - make `Engine.Retry` treat `archived-item-unproven` like any other
       park: `TestARetryReturnsAnUnprovenArchivedItemToQueued` fails.
4. MUST execute the tests task_01, task_02 and task_03 name in their
   Verification.
5. MUST record, as evidence this Spec did not author, each source of
   `_prd.md` → Acceptance evidence:
   - run `git log -1 --format=%s --grep='(#619)'` in the audited checkout and
     record the subject of Spec 0292's squash commit; run `git show
     --name-status --format= <that commit>` and record that it adds
     `docs/history/specs/0292-runs-causes-reads-the-diagnostic-where-the-run-wrote-it.md`
     and deletes the Spec folder's `_prd.md`, and read `pull_request` from
     that record at the audited head (`"619"`), so the Pull Request merged by
     hand carried both the archive and the Pull Request Record Commit, the
     shape of `TestDeliverStartResumesAnArchivedItemBranchAtReview`;
   - because this Spec's authorization forbids network access, record
     GitHub's protected-branch documentation as `blocked (environment:
     network prohibited by Spec authorization)`, quoting the passage the PRD
     cites. Such a row never decides a qualifying `partial`.
6. MUST record Success Metric 6 as exactly one row, `blocked (declared
   unreachable)`, against the PRD's single Unreachable Acceptance entry, and
   count it once in `rows_blocked_declared`: no row may reach a forge.
7. MUST perform the glossary and documentation checks:
   - record that `CONTEXT.md` defines **Archived Item Resume**, links
     ADR-0312, and that **Delivery Revalidation** says it does not run for
     it; record that the glossary contract tests pass and whether any other
     term the Spec introduced needs an entry (`docs/agents/domain.md`);
   - record the paragraph and sentences `docs/user-guide/commands/deliver.md`
     and both copies of the Roundfix `deliver` reference gained, and that
     `TestTheSkillCoverageMapIsCurrent` passes;
   - search the guides, the agent guides and the skill copies for wording
     that says every continued item is revalidated before its first Run, and
     record what it finds.
8. MUST verify from the repository history:
   - each Task's changed files stay within its declarations or its `##
     Recorded paths`;
   - every Governed Path changed is bounded in `_authorization.md`;
   - these are byte-identical to the Delivery Base: the Makefile, the CI
     workflows, `go.mod`, `go.sum`, `.roundfixrc.yml`, every Baseline module
     under `internal/baseline/assets/`, `docs/agents/setup-context.json`,
     every file under `docs/history/`, every file under `internal/spec/` and
     `internal/store/`, every `_test.go` file under `internal/cli` that
     exists at the Delivery Base, `internal/cli/deliver_revalidate.go`,
     `internal/cli/archive_reader.go` and `docs/references/skill-coverage.json`;
   - in `internal/delivery/engine.go`, only `advanceItem`, `Retry`, the
     blocker constants, the new types and the dependency wiring changed;
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
    `bin/roundfix archive 0299-a-new-delivery-queue-resumes-an-item-whose-spec-is-archived-on-its-branch --plan --offline`
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
- [ ] Read Spec 0292's squash commit and record from history.
- [ ] Record the outside evidence, the declared row, the glossary and documentation checks and the scope audit, and write the QA report.

## Acceptance Criteria

- [ ] A continued branch with a proven archive reaches `merged` through `roundfix deliver start` with one review and no Implement Run, with and without another active Spec, and each removed guard makes its named test fail.
- [ ] Each unproven branch parks `archived-item-unproven` with its cause and the `finding` Park Class, and its retry returns it to `queued`.
- [ ] The scope audit finds no undeclared Governed Path change and no change outside the declared functions and files.

## Verification

The following command is rendered unchanged from Roundfix's derived QA contract.
It proves a newest report exists and its verdict is readable; eligibility is
applied in process by whoever settles the Task.

- `newest="$(find 'docs/specs/0299-a-new-delivery-queue-resumes-an-item-whose-spec-is-archived-on-its-branch/qa' -type f -name "qa-report-*.md" -print 2>/dev/null | awk "{ report=\$0; name=\$0; parts=split(name, path, \"/\"); name=path[parts]; name=substr(name, 11, length(name)-13); date=substr(name, 1, 10); suffix=substr(name, 11); shape=(length(date) == 10 && substr(date, 5, 1) == \"-\" && substr(date, 8, 1) == \"-\"); for (i=1; i <= 10 && shape; i++) { if (i != 5 && i != 8 && index(\"0123456789\", substr(date, i, 1)) == 0) shape=0 } year=substr(date, 1, 4)+0; month=substr(date, 6, 2)+0; day=substr(date, 9, 2)+0; leap=(year%400 == 0 || (year%4 == 0 && year%100 != 0)); days=(month == 2 ? 28+leap : ((month == 4 || month == 6 || month == 9 || month == 11) ? 30 : 31)); dated=(shape && month >= 1 && month <= 12 && day >= 1 && day <= days); if (!dated) date=\"\"; if (suffix == \"\") { sequenced=1; sequence=-1 } else { sequenced=(length(suffix) > 1 && substr(suffix, 1, 1) == \"-\"); for (i=2; i <= length(suffix) && sequenced; i++) { if (index(\"0123456789\", substr(suffix, i, 1)) == 0) sequenced=0 } sequence=(sequenced ? substr(suffix, 2)+0 : 0) } if (dated && sequenced) printf \"%s\\t%s\\t%s\\t%s\\t%s\\n\", dated, date, sequenced, sequence, report }" | sort -k1,1n -k2,2 -k3,3n -k4,4n -k5,5 | tail -1 | cut -f5-)"; test -n "$newest" || exit 1; awk "BEGIN { whitespace=\" \t\r\n\f\v\" } NR == 1 && \$0 == \"---\" { frontmatter=1; next } frontmatter && \$0 == \"---\" { closed=1; exit } frontmatter && index(\$0, \"verdict:\") == 1 { verdict=substr(\$0, 9); while (length(verdict) > 0 && index(whitespace, substr(verdict, 1, 1)) > 0) verdict=substr(verdict, 2); while (length(verdict) > 0 && index(whitespace, substr(verdict, length(verdict), 1)) > 0) verdict=substr(verdict, 1, length(verdict)-1); verdicts++ } END { exit(closed && verdicts == 1 && length(verdict) > 0 ? 0 : 1) }" "$newest"`

## References

- `_prd.md` → Goals; Core Features 1-4; Success Metric 1; Success Metric 2; Success Metric 3; Success Metric 4; Success Metric 5; Success Metric 6; Acceptance evidence; Unreachable Acceptance; Glossary
- `_techspec.md` → Current behavior; API Contracts 1-7; Invariants 1-6; Testing Approach; Build Order 4
- ADR-0223
- ADR-0229
- ADR-0312
