---
task: task_03
spec: 0297-a-corrective-task-after-a-completed-gate-passes-the-wave-collision-check
status: completed
type: qa
complexity: high
---

# Task 03: Run the final QA gate

## Overview

This is the authored terminal gate for this Spec. It declares what the matrix
covers and settles the Spec on evidence. The Spec makes every
`SC-WAVE-COLLISION` message, the Spec Consistency Check finding with the
refusal reason a QA Report records and the Daemon's pre-dispatch refusal, name
the one `needs` edge that orders the pair, and pins oraculum's corrective
shape so the gate precondition, the strict check, the settlement checks and
the pre-dispatch refusal keep agreeing.

Every row runs on temporary repositories, worktrees and Git configuration
under the Run's temporary directory or on the audited checkout's own history.
No row reads or writes the real `~/.roundfix`, reaches the network, sends a
request to a provider, the judge or a reviewer, reads a credential, or reads
another repository. The gate never runs `roundfix spec judge`, never runs
`roundfix review` against a real provider, never runs `roundfix deliver`, and
never runs an archive plan without `--offline`.

## Requirements

1. MUST build the binary with `make build` and run the repository
   Verification, recording its result as a gate fact.
2. MUST verify Success Metrics 1, 2 and 3:
   - Run `TestQAGateJudgesACorrectiveTaskOnTheCurrentGraph` in
     `./internal/daemon`, `TestWaveCollisionFindingNamesTheResolvingEdge` in
     `./internal/speccheck`, and
     `TestWaveCollisionNamesTheEdgeThatOrdersThePair` and
     `TestWaveCollisionSharedPathListSortsPathsWithTheirSources` in
     `./internal/spec`, and record for each case the reason or edge it
     asserted (Invariants 1-4).
   - In a temporary repository built in the shape the Daemon test builds (a
     plan commit, a Run Branch `roundfix/run-run_A` whose `task_02`
     settlement commit carries `Roundfix-Spec` and `Roundfix-Task` trailers
     and writes `pkg/versions.test.ts` undeclared, a fast-forward, an archive
     commit and its `git revert`, and `task_05` declaring that file under
     `interface:` with `needs: [task_01]`), run `bin/roundfix spec check
     <fixture slug> --format json` and record the `SC-WAVE-COLLISION`
     summary and fix; then commit `needs: [task_02]`, run it again and
     record that no `SC-WAVE-COLLISION` remains.
3. MUST verify Success Metric 4 and Invariant 5:
   - Run `TestTaskCycleRefusesAWaveCollision`,
     `TestTaskCycleDispatchesSerializedCollisionGraph`,
     `TestWaveCollisionIsReportedAtAuthoring`,
     `TestWaveCollisionCheckStillReportsASharedInterfacePath`,
     `TestWaveCollisionRefusesTheSameCommandFile`,
     `TestCollisionsExcludesTransitivelyOrderedTasks` and
     `TestCollisionsReadsSettlementsOnRunBranchesNotOnlyHead`, and record
     that each passes with its test file byte-identical to the Delivery
     Base.
   - On temporary copies of the audited head made with `git archive`, apply
     one sabotage per copy and record that the named test fails:
     - make `ResolvingEdge` name the edge from `First` to `Second`:
       `TestWaveCollisionNamesTheEdgeThatOrdersThePair`,
       `TestWaveCollisionFindingNamesTheResolvingEdge` and
       `TestQAGateJudgesACorrectiveTaskOnTheCurrentGraph` fail, because
       each asserts the exact edge for a pair whose order is known;
     - remove `; ` and `ResolvingEdge` from the summary
       `detectWaveCollisions` writes:
       `TestWaveCollisionFindingNamesTheResolvingEdge` and
       `TestQAGateJudgesACorrectiveTaskOnTheCurrentGraph` fail, because both
       compare the whole reason;
     - make `Collisions` ignore the prior-Run settlement paths:
       `TestQAGateJudgesACorrectiveTaskOnTheCurrentGraph` fails, because its
       first case's `task_02` touches the shared file only through its
       settlement commit;
     - make `Collisions` return no collision:
       `TestWaveCollisionNamesTheEdgeThatOrdersThePair`,
       `TestWaveCollisionFindingNamesTheResolvingEdge` and
       `TestQAGateJudgesACorrectiveTaskOnTheCurrentGraph` fail.
4. MUST execute the tests task_01 and task_02 name in their Verification.
5. MUST record, as evidence this Spec did not author, each source of
   `_prd.md` → Acceptance evidence:
   - in the audited checkout, read Spec 0296's Archive Record under
     `docs/history/specs/` and record its `source_revision`,
     `3dc0b8dd17eeb4cebc861e71ef4cb16d8a079438`, and `qa_verdict`; run
     `git show` of that revision's `_tasks.md`, `task_02.md` and `task_04.md`
     for Spec 0296 and record `task_04`'s `needs` and that both Task files
     name `internal/cli/deliver_queue_subject_test.go`; add a temporary
     detached worktree at that revision, run `bin/roundfix spec check
     0296-the-queues-archive-subject-and-pull-request-title-fit-100-characters
     --format json` there, record the number of `SC-WAVE-COLLISION` findings,
     and remove the worktree;
   - because this Spec's authorization forbids network access, record the
     GCC "Guidelines for Diagnostics" page and the Barik et al. study as
     `blocked (environment: network prohibited by Spec authorization)`,
     quoting the passages the PRD cites. Such a row never decides a
     qualifying `partial`.
6. MUST record Success Metric 5 as exactly one row, `blocked (declared
   unreachable)`, against the PRD's single Unreachable Acceptance entry, and
   count it once in `rows_blocked_declared`: no row may read another
   repository.
7. MUST perform the glossary and documentation checks:
   - record that the Spec declares no glossary term, that the glossary
     contract tests pass, and whether any term the Spec introduced needs an
     entry (`docs/agents/domain.md`);
   - record the row `docs/user-guide/commands/spec.md` gained, the bullet
     both copies of the Roundfix `spec` reference gained, and that
     `TestTheSkillCoverageMapIsCurrent` passes;
   - search the guides, the agent guides and the skill copies for wording
     about Wave collisions or the `needs` remedy that contradicts the new
     finding, and record what it finds.
8. MUST verify from the repository history:
   - each Task's changed files stay within its declarations or its `##
     Recorded paths`;
   - every Governed Path changed is bounded in `_authorization.md`;
   - these are byte-identical to the Delivery Base: the Makefile, the CI
     workflows, `go.mod`, `go.sum`, `.roundfixrc.yml`, every Baseline module
     under `internal/baseline/assets/`, `docs/agents/setup-context.json`,
     `CONTEXT.md`, every file under `docs/history/` and `internal/delivery/`,
     every `_test.go` file under `internal/spec`, `internal/speccheck` and
     `internal/daemon` that exists at the Delivery Base, and
     `docs/references/skill-coverage.json`;
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
    `bin/roundfix archive 0297-a-corrective-task-after-a-completed-gate-passes-the-wave-collision-check --plan --offline`
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
- [ ] Reproduce the corrective shape through `roundfix spec check` and read Spec 0296's corrective graph from history.
- [ ] Record the outside evidence, the declared row, the glossary and documentation checks and the scope audit, and write the QA report.

## Acceptance Criteria

- [ ] Every Wave collision message names the resolving edge, and each removed guard makes its named test fail.
- [ ] In the corrective shape the gate, the strict check, the settlement check and the pre-dispatch refusal agree before and after the edge.
- [ ] The scope audit finds no undeclared Governed Path change and no change to an existing test.

## Verification


The following command is rendered unchanged from Roundfix's derived QA contract.
It proves a newest report exists and its verdict is readable; eligibility is
applied in process by whoever settles the Task.

- `newest="$(find 'docs/specs/0297-a-corrective-task-after-a-completed-gate-passes-the-wave-collision-check/qa' -type f -name "qa-report-*.md" -print 2>/dev/null | awk "{ report=\$0; name=\$0; parts=split(name, path, \"/\"); name=path[parts]; name=substr(name, 11, length(name)-13); date=substr(name, 1, 10); suffix=substr(name, 11); shape=(length(date) == 10 && substr(date, 5, 1) == \"-\" && substr(date, 8, 1) == \"-\"); for (i=1; i <= 10 && shape; i++) { if (i != 5 && i != 8 && index(\"0123456789\", substr(date, i, 1)) == 0) shape=0 } year=substr(date, 1, 4)+0; month=substr(date, 6, 2)+0; day=substr(date, 9, 2)+0; leap=(year%400 == 0 || (year%4 == 0 && year%100 != 0)); days=(month == 2 ? 28+leap : ((month == 4 || month == 6 || month == 9 || month == 11) ? 30 : 31)); dated=(shape && month >= 1 && month <= 12 && day >= 1 && day <= days); if (!dated) date=\"\"; if (suffix == \"\") { sequenced=1; sequence=-1 } else { sequenced=(length(suffix) > 1 && substr(suffix, 1, 1) == \"-\"); for (i=2; i <= length(suffix) && sequenced; i++) { if (index(\"0123456789\", substr(suffix, i, 1)) == 0) sequenced=0 } sequence=(sequenced ? substr(suffix, 2)+0 : 0) } if (dated && sequenced) printf \"%s\\t%s\\t%s\\t%s\\t%s\\n\", dated, date, sequenced, sequence, report }" | sort -k1,1n -k2,2 -k3,3n -k4,4n -k5,5 | tail -1 | cut -f5-)"; test -n "$newest" || exit 1; awk "BEGIN { whitespace=\" \t\r\n\f\v\" } NR == 1 && \$0 == \"---\" { frontmatter=1; next } frontmatter && \$0 == \"---\" { closed=1; exit } frontmatter && index(\$0, \"verdict:\") == 1 { verdict=substr(\$0, 9); while (length(verdict) > 0 && index(whitespace, substr(verdict, 1, 1)) > 0) verdict=substr(verdict, 2); while (length(verdict) > 0 && index(whitespace, substr(verdict, length(verdict), 1)) > 0) verdict=substr(verdict, 1, length(verdict)-1); verdicts++ } END { exit(closed && verdicts == 1 && length(verdict) > 0 ? 0 : 1) }" "$newest"`

## References

- `_prd.md` → Goals; Core Features 1-4; Success Metric 1; Success Metric 2; Success Metric 3; Success Metric 4; Success Metric 5; Acceptance evidence; Unreachable Acceptance
- `_techspec.md` → Current behavior; API Contracts 1-6; Invariants 1-5; Testing Approach; Build Order 3
- ADR-0117
- ADR-0182
- ADR-0096
