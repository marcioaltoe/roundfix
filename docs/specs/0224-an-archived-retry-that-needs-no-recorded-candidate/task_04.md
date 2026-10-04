---
task: task_04
spec: 0224-an-archived-retry-that-needs-no-recorded-candidate
status: pending
type: qa
complexity: medium
---

# Task 04: Run the final QA gate

## Overview

The authored terminal gate for this Spec. It declares what the matrix covers
and settles the Spec on evidence. This Spec lets a Delivery Retry resume an
item the operator archived with the QA Archive Override from the Implement
start head of its Run whatever its park, parks a zero-finding QA partial with
any environment-blocked row as `qa-environment-partial`, and describes both in
the Roundfix Skill and the `deliver` guide. Every behavior row is exercised by
executing the named tests against the built tree, with disposable
repositories, real archives and a disposable Run Database. No command reaches
GitHub, a provider or the network, except the reads of published pages the
outside-evidence row names.

## Requirements

1. MUST run the repository Verification and record its result as a gate fact.
2. MUST verify API Contract 1 and Success Metric 1 by executing task_01's
   tests against the built tree, and record that
   `TestArchivedRetryOfARunUnresolvedItemWithoutACandidateReturnsToReview`
   resumes a `run-unresolved` item at `reviewing` with the archived head as
   its only candidate after a real `roundfix archive --qa-override`, and that
   both refusals of API Contract 1 keep their text and leave the item
   unchanged, and by executing task_05's test, that an override archive
   without a candidate and without History refuses with the archived-head
   text naming an empty candidate and leaves the item unchanged.
3. MUST verify API Contract 2 and Success Metric 2 by executing task_02's
   tests against the built tree, and MUST additionally feed the archived
   0220 QA Report (`docs/history/specs/0220-tests-and-pins-that-hold-in-every-environment/qa/qa-report-2026-10-03.md`),
   copied into a disposable repository's Run Branch, through the classification
   with a probe that writes nothing to the repository (for example a
   `go test -overlay` test), recording that it now sets the environment
   partial and that the Archive Command's eligibility for the same report is
   unchanged.
4. MUST verify Success Metric 3: every existing test of
   `internal/delivery/operator_archive_retry_test.go`,
   `internal/cli/deliver_archived_retry_test.go` and
   `internal/cli/deliver_operator_archive_test.go` passes, and the only
   changed expectation is the `pre-PR only` case, read from the diff.
5. MUST verify that the Roundfix Skill's `deliver` reference and the
   `deliver` guide describe both rules and no longer carry the sentences
   task_03 removes, that each mirror equals its canonical file, that the
   version was raised and recorded, and that the `### QA settlement` section
   of every skill is unchanged.
6. MUST record, as evidence this Spec did not author, each source of
   `_prd.md` → Acceptance evidence: the operator's intervention log entries
   96 and 140 to 142; the 0220 queue item in the live Run Database, read only
   with `sqlite3 -readonly`; the archived 0220 QA Report; Git's
   `git merge-base` documentation; and Temporal's recovery guide. For each
   source it reaches it MUST record what it read and whether it still
   supports the design; for each it cannot reach it MUST record the row as
   blocked with the reason.
7. MUST verify that this Spec's own artifacts satisfy the promise rule, and
   perform the glossary check of `docs/agents/domain.md`: record whether
   **Delivery Retry**, **QA Archive Override** and **Park Class** still
   describe the behavior, and whether the Implement start head needs a
   glossary term.
8. MUST verify from the repository history that each Task's changed files stay
   within its declarations or its `## Recorded paths`, that every Governed
   Path changed is bounded in `_authorization.md`, and that `Makefile`,
   `.roundfixrc.yml`, `go.mod` and the CI workflows did not change. This row
   reads Task commits, so it declares a `commit_range` input.
9. MUST record the non-waivable Pull Request row as
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
10. MUST NOT accept a row whose only evidence is that a file was read.
11. MUST NOT start a Delivery Queue, run `roundfix deliver`, merge, re-run or
    otherwise change anything on GitHub, and MUST NOT write to the live Run
    Database under `~/.roundfix`.

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

- `newest="$(find 'docs/specs/0224-an-archived-retry-that-needs-no-recorded-candidate/qa' -type f -name "qa-report-*.md" -print 2>/dev/null | awk "{ report=\$0; name=\$0; parts=split(name, path, \"/\"); name=path[parts]; name=substr(name, 11, length(name)-13); date=substr(name, 1, 10); suffix=substr(name, 11); shape=(length(date) == 10 && substr(date, 5, 1) == \"-\" && substr(date, 8, 1) == \"-\"); for (i=1; i <= 10 && shape; i++) { if (i != 5 && i != 8 && index(\"0123456789\", substr(date, i, 1)) == 0) shape=0 } year=substr(date, 1, 4)+0; month=substr(date, 6, 2)+0; day=substr(date, 9, 2)+0; leap=(year%400 == 0 || (year%4 == 0 && year%100 != 0)); days=(month == 2 ? 28+leap : ((month == 4 || month == 6 || month == 9 || month == 11) ? 30 : 31)); dated=(shape && month >= 1 && month <= 12 && day >= 1 && day <= days); if (!dated) date=\"\"; if (suffix == \"\") { sequenced=1; sequence=-1 } else { sequenced=(length(suffix) > 1 && substr(suffix, 1, 1) == \"-\"); for (i=2; i <= length(suffix) && sequenced; i++) { if (index(\"0123456789\", substr(suffix, i, 1)) == 0) sequenced=0 } sequence=(sequenced ? substr(suffix, 2)+0 : 0) } if (dated && sequenced) printf \"%s\\t%s\\t%s\\t%s\\t%s\\n\", dated, date, sequenced, sequence, report }" | sort -k1,1n -k2,2 -k3,3n -k4,4n -k5,5 | tail -1 | cut -f5-)"; test -n "$newest" || exit 1; awk "BEGIN { whitespace=\" \t\r\n\f\v\" } NR == 1 && \$0 == \"---\" { frontmatter=1; next } frontmatter && \$0 == \"---\" { closed=1; exit } frontmatter && index(\$0, \"verdict:\") == 1 { verdict=substr(\$0, 9); while (length(verdict) > 0 && index(whitespace, substr(verdict, 1, 1)) > 0) verdict=substr(verdict, 2); while (length(verdict) > 0 && index(whitespace, substr(verdict, length(verdict), 1)) > 0) verdict=substr(verdict, 1, length(verdict)-1); verdicts++ } END { exit(closed && verdicts == 1 && length(verdict) > 0 ? 0 : 1) }" "$newest"`

## References

`_prd.md` → Goals; Core Features 1-3; Success Metrics 1-3; Acceptance
evidence; `_techspec.md` → Testing Approach 1-3; API Contract 1; API
Contract 2; Surface Transcripts; Integration Points; ADR-0080; ADR-0088;
ADR-0091; ADR-0104; ADR-0155; ADR-0167; ADR-0229.
