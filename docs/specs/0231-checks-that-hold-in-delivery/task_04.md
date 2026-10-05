---
task: task_04
spec: 0231-checks-that-hold-in-delivery
status: pending
type: qa
complexity: high
---

# Task 04: Run the final QA gate

## Overview

The authored terminal gate for this Spec. It declares what the matrix covers
and settles the Spec on evidence. This Spec makes the darwin detach fixture
reading count an exiting member as ended, makes CI's pull request job test
the head merged with the default branch tip it fetches and record it as a
`tested-base` annotation, and makes the Delivery Queue re-run a failed check
that tested an older default branch instead of parking it. Every behavior row
is exercised against the built tree with disposable repositories and scripted
GitHub responses. No command reaches GitHub.

## Requirements

1. MUST run the repository Verification and record its result as a gate fact.
2. MUST verify Success Metric 1 by executing task_01's tests: the
   exiting-member test observes the window and reads the group as ended, the
   unreaped-member test still reads a running member as live, and the detach
   death, survivor and both group tests pass ten consecutive iterations each
   with no skip; and MUST record from the diff that no production code and no
   Linux or other Unix reading changed.
3. MUST verify Success Metric 2 and API Contract 1 by running the workflow's
   `Merge the current base branch` step, read from the workflow file, in a
   disposable clone whose default branch moved after the head branched: the
   merge's parents are the head and the moved tip, `VERIFY_BASE` is that tip,
   exactly one `tested-base` notice names it, and a conflicting head fails the
   step; and MUST record that the `push` path and `Verify docs` are
   unchanged.
4. MUST verify Success Metric 3 and API Contracts 2 to 4 by executing task_03's
   tests, including the replay of Pull Request #391, the real-git staleness
   decision and the stale log line.
5. MUST verify Success Metric 4 and API Contract 5: every existing check
   re-run test passes, a failure without a `tested-base` annotation is
   classified as before, and the existing blockers, Park Classes, log lines
   and Warning keep their text.
6. MUST verify that `docs/user-guide/commands/deliver.md` states the
   annotation, the stale rule and the log line of the Vocabulary Contract.
7. MUST record, as evidence this Spec did not author, each source of
   `_prd.md` → Acceptance evidence: GitHub's re-run documentation, the CI
   runs of Pull Request #391 and their checkout lines, GitHub's workflow
   commands and check-run annotation documentation, the `actions/checkout`
   README, XNU's exit, signal and sysctl sources, and the Daemon's
   Verification log of Run `run_20261005T101255Z_f599a5cc8ea7f4c2`. For each
   source it reaches it MUST record what it read and whether it still
   supports the design; for each it cannot reach it MUST record the row as
   blocked with the reason. Reading public pages or this repository's run
   history with `gh` read commands is allowed; nothing on GitHub may change.
8. MUST verify that this Spec's own artifacts satisfy the promise rule, and
   perform the glossary check of `docs/agents/domain.md`: record whether
   **Delivery Queue**, **Park Class** and **Delivery Retry** still describe
   the behavior, and whether "tested base" needs a glossary term.
9. MUST verify from the repository history that each Task's changed files stay
   within its declarations or its `## Recorded paths`, that every Governed
   Path changed is bounded in `_authorization.md`, that task_02's commit
   changes `.github/workflows/ci-verify.yml` alone, and that `Makefile`,
   `go.mod` and every other workflow did not change. This row reads Task
   commits, so it declares a `commit_range` input.
10. MUST record the non-waivable Pull Request row as
    `blocked (environment: no open Pull Request)`, naming the Pull Request row
    in its provenance, and support it with the pre-PR equivalent of each
    control:
    - approval: the maintainer's recorded delivery authority;
    - checks and status: the Daemon's repository Verification at the audited
      head; the new workflow step first runs on this Spec's own Pull Request;
    - unresolved threads: none, because no Pull Request exists;
    - Merge-Ready acceptance: none yet, because the pre-PR review follows
      archive;
    - review-artifact ancestry: the audited head named as the claimed
      candidate.
11. MUST NOT accept a row whose only evidence is that a file was read.
12. MUST NOT start a Delivery Queue, merge, re-run, reopen or otherwise change
    anything on GitHub, and MUST NOT read or write the real `~/.roundfix`
    except the read-only Run log named in Requirement 7.

## Subtasks

- [ ] Build the matrix from the Requirements above and the sources no
      declaration waives.
- [ ] Execute each row against the built tree and record its evidence.
- [ ] Write the dated QA Report with its verdict.

## Acceptance Criteria

- [ ] The QA Report records a verdict and names, in each row's provenance, the
      sources that row covers.
- [ ] The repository Verification result is recorded as a fact, not re-derived.

## Context

- instruction: `.agents/skills/qa-gate/SKILL.md`

## Verification

The following command is rendered unchanged from Roundfix's derived QA contract.
It proves a newest report exists and its verdict is readable; eligibility is
applied in process by whoever settles the Task.

- `newest="$(find 'docs/specs/0231-checks-that-hold-in-delivery/qa' -type f -name "qa-report-*.md" -print 2>/dev/null | awk "{ report=\$0; name=\$0; parts=split(name, path, \"/\"); name=path[parts]; name=substr(name, 11, length(name)-13); date=substr(name, 1, 10); suffix=substr(name, 11); shape=(length(date) == 10 && substr(date, 5, 1) == \"-\" && substr(date, 8, 1) == \"-\"); for (i=1; i <= 10 && shape; i++) { if (i != 5 && i != 8 && index(\"0123456789\", substr(date, i, 1)) == 0) shape=0 } year=substr(date, 1, 4)+0; month=substr(date, 6, 2)+0; day=substr(date, 9, 2)+0; leap=(year%400 == 0 || (year%4 == 0 && year%100 != 0)); days=(month == 2 ? 28+leap : ((month == 4 || month == 6 || month == 9 || month == 11) ? 30 : 31)); dated=(shape && month >= 1 && month <= 12 && day >= 1 && day <= days); if (!dated) date=\"\"; if (suffix == \"\") { sequenced=1; sequence=-1 } else { sequenced=(length(suffix) > 1 && substr(suffix, 1, 1) == \"-\"); for (i=2; i <= length(suffix) && sequenced; i++) { if (index(\"0123456789\", substr(suffix, i, 1)) == 0) sequenced=0 } sequence=(sequenced ? substr(suffix, 2)+0 : 0) } if (dated && sequenced) printf \"%s\\t%s\\t%s\\t%s\\t%s\\n\", dated, date, sequenced, sequence, report }" | sort -k1,1n -k2,2 -k3,3n -k4,4n -k5,5 | tail -1 | cut -f5-)"; test -n "$newest" || exit 1; awk "BEGIN { whitespace=\" \t\r\n\f\v\" } NR == 1 && \$0 == \"---\" { frontmatter=1; next } frontmatter && \$0 == \"---\" { closed=1; exit } frontmatter && index(\$0, \"verdict:\") == 1 { verdict=substr(\$0, 9); while (length(verdict) > 0 && index(whitespace, substr(verdict, 1, 1)) > 0) verdict=substr(verdict, 2); while (length(verdict) > 0 && index(whitespace, substr(verdict, length(verdict), 1)) > 0) verdict=substr(verdict, 1, length(verdict)-1); verdicts++ } END { exit(closed && verdicts == 1 && length(verdict) > 0 ? 0 : 1) }" "$newest"`

## References

`_prd.md` → Goals; Core Features 1-3; Success Metrics 1-4; Acceptance
evidence; `_techspec.md` → Interfaces; API Contract 1; API Contract 2; API
Contract 3; API Contract 4; API Contract 5; Vocabulary Contract; Testing
Approach; Build Order 4; ADR-0236; ADR-0213; ADR-0080; ADR-0091.
