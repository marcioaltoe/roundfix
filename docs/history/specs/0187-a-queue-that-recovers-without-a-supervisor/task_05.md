---
task: task_05
spec: 0187-a-queue-that-recovers-without-a-supervisor
status: completed
type: qa
complexity: medium
---

# Task 05: Run the final QA gate

## Overview

The authored terminal gate for this Spec. It declares what the matrix covers
and settles the Spec on evidence. This Spec changes production code in
`internal/cli`, `internal/delivery` and `internal/speccheck`, and the guides
that describe delivery cleanup, retry, owner warnings and the authorization
audit. Every behavior row is exercised through the built binary or by
executing the named tests against the built tree. Each binary command runs in a
disposable repository with a disposable Roundfix Home. No command reaches
GitHub, a provider or a live remote.

## Requirements

1. MUST run the repository Verification and record its result as a gate fact.
2. MUST verify, by executing task_01's tests against the built tree, that an
   item merged on a bare `origin` but not fetched locally is released without
   a warning, that an unreachable `origin` reports `refresh default branch`,
   and that a repository without `origin` resolves as before.
3. MUST verify, by executing task_02's tests against the built tree, that a
   refusal caused only by an amending non-Task commit names it and prints the
   five ordered commands, and that every other refusal keeps today's text.
4. MUST verify, by executing task_03's tests against the built tree:
   - a source change on the starting main after the owner's build commit
     yields `owner-older-than-main`;
   - a docs-only main, a current owner and an absent build commit yield no
     warning;
   - the item still reaches `running`.
5. MUST verify, by executing task_04's tests against the built tree, that the
   grant a Task ran under authorizes its commit when the delivery target
   holds it, and that a record main never held, a self-approving commit and a
   fork-point-covered commit behave as before.
6. MUST record, as evidence this Spec did not author, the authorization
   audit table of the QA reports on the Run branches
   `roundfix/run-run_20260929T200238Z_ca4a5877c9dc9d5c` and
   `roundfix/run-run_20260929T201521Z_3d1a54f26096ba81`. Read it with
   `git show <branch>:docs/specs/0181-gates-that-refuse-only-what-someone-can-act-on/qa/qa-report-2026-09-29.md`.
   It shows task_07 read at revision `6784210b`. The gate MUST also record
   worktrunk issue #3519
   (<https://github.com/max-sixty/worktrunk/issues/3519>) as an independent
   report of the stale-local-ref class. When a branch or the page is
   unavailable, it MUST record the row as blocked with its reason.
7. MUST verify that `docs/user-guide/commands.md`, the Roundfix skill and its
   mirror carry `refreshes the default branch from the delivery remote`,
   `amended by`, `owner-older-than-main` and `the grant the Task ran under`,
   and that `make skills-sync-check` exits `0`.
8. MUST verify that this Spec's own artifacts satisfy the promise rule, and
   check whether it introduced, changed or retired a glossary term that
   `CONTEXT.md` does not carry.
9. MUST verify from the repository history that each Task's changed files stay
   within its declarations or its `## Recorded paths`, and that every Governed
   Path is bounded in `_authorization.md`.
10. MUST record the non-waivable Pull Request row as
    `blocked (environment: no open Pull Request)`, naming the Pull Request row
    in its provenance, and support it with the pre-PR equivalent of each
    control:
    - approval: the maintainer's recorded delivery authority;
    - checks and status: the Daemon's repository Verification at the audited
      head;
    - unresolved threads: none, because no Pull Request exists;
    - Merge-Ready acceptance: none yet, because the pre-PR review follows
      archive;
    - review-artifact ancestry: the audited head named as the claimed
      candidate.
11. MUST NOT accept a row whose only evidence is that a file was read.
12. MUST NOT write to the live Run Database under `~/.roundfix`.

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

- `newest="$(find 'docs/specs/0187-a-queue-that-recovers-without-a-supervisor/qa' -type f -name "qa-report-*.md" -print 2>/dev/null | awk "{ report=\$0; name=\$0; parts=split(name, path, \"/\"); name=path[parts]; name=substr(name, 11, length(name)-13); date=substr(name, 1, 10); suffix=substr(name, 11); shape=(length(date) == 10 && substr(date, 5, 1) == \"-\" && substr(date, 8, 1) == \"-\"); for (i=1; i <= 10 && shape; i++) { if (i != 5 && i != 8 && index(\"0123456789\", substr(date, i, 1)) == 0) shape=0 } year=substr(date, 1, 4)+0; month=substr(date, 6, 2)+0; day=substr(date, 9, 2)+0; leap=(year%400 == 0 || (year%4 == 0 && year%100 != 0)); days=(month == 2 ? 28+leap : ((month == 4 || month == 6 || month == 9 || month == 11) ? 30 : 31)); dated=(shape && month >= 1 && month <= 12 && day >= 1 && day <= days); if (!dated) date=\"\"; if (suffix == \"\") { sequenced=1; sequence=-1 } else { sequenced=(length(suffix) > 1 && substr(suffix, 1, 1) == \"-\"); for (i=2; i <= length(suffix) && sequenced; i++) { if (index(\"0123456789\", substr(suffix, i, 1)) == 0) sequenced=0 } sequence=(sequenced ? substr(suffix, 2)+0 : 0) } if (dated && sequenced) printf \"%s\\t%s\\t%s\\t%s\\t%s\\n\", dated, date, sequenced, sequence, report }" | sort -k1,1n -k2,2 -k3,3n -k4,4n -k5,5 | tail -1 | cut -f5-)"; test -n "$newest" || exit 1; awk "BEGIN { whitespace=\" \t\r\n\f\v\" } NR == 1 && \$0 == \"---\" { frontmatter=1; next } frontmatter && \$0 == \"---\" { closed=1; exit } frontmatter && index(\$0, \"verdict:\") == 1 { verdict=substr(\$0, 9); while (length(verdict) > 0 && index(whitespace, substr(verdict, 1, 1)) > 0) verdict=substr(verdict, 2); while (length(verdict) > 0 && index(whitespace, substr(verdict, length(verdict), 1)) > 0) verdict=substr(verdict, 1, length(verdict)-1); verdicts++ } END { exit(closed && verdicts == 1 && length(verdict) > 0 ? 0 : 1) }" "$newest"`

## References

`_prd.md` → Goals 1-4; Core Features 1-4; Success Metrics 1-4; Acceptance
evidence; `_techspec.md` → Testing Approach 1-5; API Contracts 1-3; ADR-0053;
ADR-0080; ADR-0081; ADR-0088; ADR-0090; ADR-0091; ADR-0093; ADR-0094;
ADR-0096; ADR-0104; ADR-0117; ADR-0149; ADR-0155; ADR-0156; ADR-0158;
ADR-0160; ADR-0161; ADR-0166; ADR-0167; ADR-0169; ADR-0170; ADR-0176;
ADR-0178.

## Result
