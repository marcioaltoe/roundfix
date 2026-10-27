---
task: task_05
spec: 0330-settle-follows-the-branch-plan-and-declared-rows-match-their-criteria
status: pending
type: qa
complexity: high
---

# Task 05: Run the final QA gate

## Overview

This is the authored terminal gate for this Spec. It declares what the matrix
covers and settles the Spec on evidence. The Spec makes QA acceptance match
each declared-blocked row to a declaration naming its criterion, with one
declaration naming several criteria (ADR-0343, making ADR-0240 exact); makes
`roundfix settle` refuse inside a Run or Task Worktree, follow the Spec plan
its Run's branch carries, and refuse the untracked executables a Task commit
refuses (ADR-0340, ADR-0157).

Every row runs on temporary directories, temporary copies of the audited head
made with `git archive`, compiled test binaries, the built `bin/roundfix` and
this repository's own documents and Git history. No row reads or writes
`~/.roundfix`, reaches the network, sends a request to a provider, the judge
or a reviewer, reads a credential, or reads another repository. Every
`go test` carries an explicit `-timeout`, and every row that runs
`bin/roundfix` sets `HOME` to a fresh temporary directory, because the CLI
loads `$HOME/.roundfix/config.yml`; remove each temporary HOME afterwards.
Evidence is written into the checkout only after the suites end.

## Requirements

1. MUST build the binary with `make build` and run the repository
   Verification, recording its result as a gate fact.
2. MUST execute the commands task_01, task_02, task_03 and task_04 name in
   their Verification and record each exit status.
3. MUST verify Success Metrics 1 to 3 and Invariants 1 and 2: run the five
   tests task_01 adds with `-count=3 -timeout 15m -v` and record each passes;
   then, through the built product with a fresh `HOME`, write in a temporary
   directory a Spec folder whose `_prd.md` has one Unreachable Acceptance
   declaration with the two criteria `Success Metric 10, the first` and
   `Success Metric 11, the second`, and a `qa/qa-report-2026-10-28.md` with
   `verdict: partial`, `rows_blocked_declared: 2` and Results rows `Q10`
   `blocked (declared: Success Metric 10)` and `Q11`
   `blocked (declared: Success Metric 11)`; record that
   `bin/roundfix qa-report accept <report>` exits `0`, and that after adding
   a row `Q12` `blocked (declared: Success Metric 12)` and setting
   `rows_blocked_declared: 3` it exits `1` with stderr naming
   `shortfall is 2` and `declared-blocked row Q12`. In a second Spec folder
   whose `_prd.md` has the one declaration `no Run delivered after this Spec
   fails the test` and whose report has `rows_blocked_declared: 1` and the
   row `Q09` `blocked (declared: no later Run fails the test)`, record that
   `bin/roundfix qa-report accept` exits `0`, as it does at the Delivery
   Base (today's count, kept).
4. MUST verify Success Metrics 4 to 6 and Invariants 3 and 4: run the tests
   task_02, task_03 and task_04 add, and
   `TestSettleNamesNoRemovalWhenIntegrationIsPending`, with
   `-count=3 -timeout 15m -v`, and record each passes; record from
   `bin/roundfix settle --help` (fresh `HOME`) that the usage carries the
   sentences of API Contracts 9 and 10.
5. MUST run each sabotage on its own temporary copy of the audited head made
   with `git archive`, and record that the named test fails, with the case
   whose input makes the mutation observable:
   - compare the listed declared rows by count of criteria instead of by
     match: `TestQAReportEligibilityMatchesEachDeclaredRow` fails (input:
     rows `Q10` and `Q12` against two declared criteria, where a count
     passes);
   - drop today's count when the report lists declared rows, so an unmatched
     row refuses: `TestQAReportEligibilityKeepsTodaysCountForAParaphrasedRow`
     fails (input: the paraphrased row `Q09`, one declaration,
     `rows_blocked_declared: 1`);
   - let a second `criterion:` line replace the first:
     `TestUnreachableDeclarationNamesEveryCriterion` fails (input: two
     `criterion:` lines);
   - match a bare prefix without the boundary characters:
     `TestCriterionNamedMatchesTextOrLeadingLabel` fails (input:
     `Success Metric 1` against `Success Metric 10, …`);
   - drop the Run Branch check in settle preflight:
     `TestSettleRefusesInsideTheRunWorktree` fails (input: settle run with
     the kept Run Worktree as working directory);
   - drop the fast-forward: `TestSettleAdvancesTheRunWorktreeToTheBranchPlan`
     fails (input: a branch that amended the PRD after the Run started);
   - advance whenever the branch is ahead, Spec folder or not:
     `TestSettleNamesNoRemovalWhenIntegrationIsPending` fails (input: a
     branch commit outside the Spec folder);
   - merge instead of refusing a diverged drift:
     `TestSettleRefusesASpecChangedOnADivergedBranch` fails (input: a Run
     Branch commit plus a branch PRD amendment);
   - drop the executable check in settle preflight:
     `TestSettleRefusesAnUnlistedExecutableInANewDirectory` fails (input: an
     untracked `newdir/helper.sh` with mode `0755`);
   - pass no authorization to `UnauthorizedExecutables`:
     `TestSettleCommitsAListedExecutableWithItsMode` fails (input: a granted
     record listing `scripts/launch.sh`).
6. MUST verify Success Metric 7: record that each file task_01, task_02,
   task_03 and task_04 name in their phrase checks carries its rule, that
   the `qa-gate` skill has the `### Declared rows name their criterion`
   section outside `### QA settlement`, and that `TestSettlementGuidanceIsOneTable`
   passes.
7. MUST record, as evidence this Spec did not author, each source of
   `_prd.md` → Acceptance evidence:
   - because this Spec's authorization forbids network access, record the
     `git-merge` documentation as `blocked (environment: network prohibited
     by Spec authorization)`, naming the outside-evidence row in its
     provenance and quoting the two passages the PRD cites; such a row never
     decides a qualifying `partial`;
   - in a temporary repository, create a linked worktree on a second branch,
     stage a new file and edit a tracked file there, commit a change to a
     different file on the first branch and record that
     `git merge --ff-only <first branch>` in the worktree exits `0` keeping
     both local changes; then commit a change to the edited file on the first
     branch and record that the next fast-forward exits `1` with
     `Your local changes to the following files would be overwritten by
     merge`, Git's own behaviour on the QA host;
   - extract with `git show` into a temporary directory the `_prd.md` and the
     QA Report of Spec 0326 at `1398a64d3626e4d168f17a0234b749b2625d03e2`
     (`qa-report-2026-10-27-01.md`), Spec 0327 at
     `520f7ceffb7c8d9a4688daa692e51ada8244b55b`, Spec 0328 at
     `781f3fb0f26de048c7dab1ff83cb9fefab12667b` (both
     `qa-report-2026-10-27.md`) and Spec 0324 at
     `ba32f39db32e218cc973bf9e2ec777ae545a4cce` (`qa-report-2026-10-26.md`),
     the `source_revision` of each Archive Record, each into
     `<dir>/<slug>/` with its report under `qa/`; record that
     `bin/roundfix qa-report accept` (fresh `HOME`) exits `0` for all four:
     the first three by match, Spec 0324, whose row paraphrases its
     criterion, by today's count; at authoring the v0.125.0 binary also
     accepted all four, record that as context only;
   - quote Quaestor's report from the adopted Backlog Entry under
     `references/`.
8. MUST record the PRD's single Unreachable Acceptance entry as exactly one
   row whose status is `blocked (declared: <criterion>)` with the entry's
   full criterion text, counted as one in `rows_blocked_declared`.
9. MUST perform the glossary check: record that the revised
   **Unreachable Acceptance**, **Settle Command**, **Run Worktree** and
   **Authorized Executable** entries of `CONTEXT.md` carry the sentences API
   Contracts 8 to 11 give, and that no bolded phrase of the Spec's artifacts
   lacks a glossary entry or a declaration (`docs/agents/domain.md`).
10. MUST verify from the Task commits of this Run (a `commit_range` input):
    - each Task's changed files stay within its declarations or its
      `## Recorded paths`;
    - every changed Governed Path is listed in `_authorization.md` `paths`;
    - no Task commit changes the Makefile, `.github/workflows/`, `go.mod`,
      `go.sum`, `.roundfixrc.yml`, `internal/baseline/`, `docs/agents/`,
      `docs/adr/`, `docs/history/` or the text of any `### QA settlement`
      section.
11. MUST record the non-waivable Pull Request row as
    `blocked (environment: no open Pull Request)`, naming the Pull Request
    row in its provenance, supported by the pre-PR equivalent of each
    control: approval from the recorded delivery authority; checks from the
    Daemon's repository Verification at the audited head; no unresolved
    threads, because no Pull Request exists; no Merge-Ready acceptance yet;
    the audited head named as the claimed candidate.
12. MUST, for the Outcome for the Archive Record, run only
    `HOME="$(mktemp -d)" bin/roundfix archive 0330-settle-follows-the-branch-plan-and-declared-rows-match-their-criteria --plan --offline`
    with every Judge Key Variable and `TYPESAFE_API_KEY` unset, remove that
    temporary HOME afterwards, and write the report's `## Outcome` as one
    paragraph that does not end with a colon.
    It MUST NOT do any of these:
    - accept a row whose only evidence is that a file was read;
    - read or write `~/.roundfix`;
    - record any environment value;
    - send any request to a provider, the judge or a reviewer;
    - write into the checkout while a test suite runs.

## Subtasks

- [ ] Build the binary, run the repository Verification and every Task Verification.
- [ ] Run the named tests, the built-product acceptance rows and each sabotage on temporary copies.
- [ ] Record the outside sources, the declared row, the glossary check and the scope audit, and write the QA report.

## Acceptance Criteria

- [ ] Every named test passes three times, and the built product accepts one declaration with two criteria and refuses an undeclared row.
- [ ] Each sabotage makes its named test fail.
- [ ] The scope audit finds no change outside the declared files and the granted Governed Paths.

## Verification

The following command is rendered unchanged from Roundfix's derived QA contract.
It proves a newest report exists and its verdict is readable; eligibility is
applied in process by whoever settles the Task.

- `newest="$(find 'docs/specs/0330-settle-follows-the-branch-plan-and-declared-rows-match-their-criteria/qa' -type f -name "qa-report-*.md" -print 2>/dev/null | awk "{ report=\$0; name=\$0; parts=split(name, path, \"/\"); name=path[parts]; name=substr(name, 11, length(name)-13); date=substr(name, 1, 10); suffix=substr(name, 11); shape=(length(date) == 10 && substr(date, 5, 1) == \"-\" && substr(date, 8, 1) == \"-\"); for (i=1; i <= 10 && shape; i++) { if (i != 5 && i != 8 && index(\"0123456789\", substr(date, i, 1)) == 0) shape=0 } year=substr(date, 1, 4)+0; month=substr(date, 6, 2)+0; day=substr(date, 9, 2)+0; leap=(year%400 == 0 || (year%4 == 0 && year%100 != 0)); days=(month == 2 ? 28+leap : ((month == 4 || month == 6 || month == 9 || month == 11) ? 30 : 31)); dated=(shape && month >= 1 && month <= 12 && day >= 1 && day <= days); if (!dated) date=\"\"; if (suffix == \"\") { sequenced=1; sequence=-1 } else { sequenced=(length(suffix) > 1 && substr(suffix, 1, 1) == \"-\"); for (i=2; i <= length(suffix) && sequenced; i++) { if (index(\"0123456789\", substr(suffix, i, 1)) == 0) sequenced=0 } sequence=(sequenced ? substr(suffix, 2)+0 : 0) } if (dated && sequenced) printf \"%s\\t%s\\t%s\\t%s\\t%s\\n\", dated, date, sequenced, sequence, report }" | sort -k1,1n -k2,2 -k3,3n -k4,4n -k5,5 | tail -1 | cut -f5-)"; test -n "$newest" || exit 1; awk "BEGIN { whitespace=\" \t\r\n\f\v\" } NR == 1 && \$0 == \"---\" { frontmatter=1; next } frontmatter && \$0 == \"---\" { closed=1; exit } frontmatter && index(\$0, \"verdict:\") == 1 { verdict=substr(\$0, 9); while (length(verdict) > 0 && index(whitespace, substr(verdict, 1, 1)) > 0) verdict=substr(verdict, 2); while (length(verdict) > 0 && index(whitespace, substr(verdict, length(verdict), 1)) > 0) verdict=substr(verdict, 1, length(verdict)-1); verdicts++ } END { exit(closed && verdicts == 1 && length(verdict) > 0 ? 0 : 1) }" "$newest"`

## References

- `_prd.md` → Goals; Core Features 1-6; Success Metric 1; Success Metric 2; Success Metric 3; Success Metric 4; Success Metric 5; Success Metric 6; Success Metric 7; Acceptance evidence; Unreachable Acceptance; Glossary
- `_techspec.md` → Current behavior; API Contracts 1-11; Declared breaks; Invariants 1-4; Testing Approach; Build Order 5
