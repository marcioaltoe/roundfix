---
task: task_05
spec: 0228-queue-items-that-stay-current-with-main
status: pending
type: qa
complexity: high
---

# Task 05: Run the final QA gate

## Overview

The authored terminal gate for this Spec. It declares what the matrix covers
and settles the Spec on evidence. This Spec makes the owned-skill record
command raise a colliding version, lets the derived merge resolve conflicts
confined to declared version lines, returns a review-only correction after
archive to review round 2, and describes the three rules in the skills, the
guides and the repository's version rule. Every behavior row is exercised
against the built tree with disposable repositories, disposable artifact
directories and a disposable Run Database. No command reaches GitHub, a
provider or the network, except the reads of published pages the
outside-evidence row names.

## Requirements

1. MUST run the repository Verification and record its result as a gate fact.
2. MUST verify API Contract 1 and Success Metric 1 by executing task_01's
   tests against the built tree, and MUST additionally, in a disposable copy
   of the built tree, raise the content of one owned skill under its current
   version, run the record command, and record that both version fields of
   the canonical skill and its mirror name one patch above the highest
   recorded version, that the record gained exactly that entry, and that no
   recorded digest changed.
3. MUST verify API Contract 2 and Success Metric 2 by executing task_02's
   tests against the built tree, and MUST additionally reproduce the
   2026-10-04 collision end to end: in a disposable clone of the built tree
   with a local bare repository as its remote, commit on the default branch
   and on an item branch two different Roundfix Skill edits each recorded
   with the record command, then run the workflow's conflict resolution on
   the item through a probe that writes nothing to the repository (for
   example a `go test -overlay` test). Record that the merge commits with
   the `Roundfix-Delivery: derived-merge` trailer, that the merged skill's
   version is one patch above the default branch's, and that
   `TestEveryOwnedSkillVersionIsRecorded` passes in the merged tree.
4. MUST verify API Contract 3 and Success Metric 3 by executing task_03's
   tests against the built tree, recording that a correction under the
   archived Spec with every finding disposed resumes at `reviewing` with the
   head as its newest candidate, and that each refusal names its reason and
   leaves the item unchanged.
5. MUST verify Success Metric 4: every test of
   `skills/owned_skill_versions_test.go`,
   `internal/cli/deliver_conflict_test.go`,
   `internal/delivery/corrective_spec_test.go`,
   `internal/delivery/archived_head_retry_test.go` and
   `internal/cli/review_lineage_test.go` passes, and the only changed
   expectations are the two record-mode cases of
   `TestRecordingNeverReplacesARecordedVersion` and the repository
   declaration of `TestThisRepositoryDeclaresItsToolsAndDerivedPaths`, read
   from the diff.
6. MUST verify that the `deliver` and `review` guides and references, the
   configuration guide, the `implement-task` skill and the repository's
   version rule carry the phrases task_02 and task_04 require and no longer
   carry the passages task_04 removes, that each mirror equals its canonical
   file, that both raised skill versions are recorded, and that the
   `### QA settlement` section of every skill is unchanged.
7. MUST record, as evidence this Spec did not author, each source of
   `_prd.md` → Acceptance evidence: the operator's intervention log entries
   137, 138, 150, 154 and 157 to 159; the 0225 queue item in the live Run
   Database, read only with `sqlite3 -readonly`; Git's `git merge`
   documentation; the Changesets decisions; and Gerrit's review label
   documentation. For each source it reaches it MUST record what it read and
   whether it still supports the design; for each it cannot reach it MUST
   record the row as blocked with the reason.
8. MUST verify that this Spec's own artifacts satisfy the promise rule, and
   perform the glossary check of `docs/agents/domain.md`: record whether
   **Delivery Retry**, **Reviewer Lineage** and **Park Class** still describe
   the behavior, and whether a line-scoped derived path needs a glossary
   term.
9. MUST verify from the repository history that each Task's changed files
   stay within its declarations or its `## Recorded paths`, that every
   Governed Path changed is bounded in `_authorization.md`, and that
   `Makefile`, `go.mod` and the CI workflows did not change. This row reads
   Task commits, so it declares a `commit_range` input.
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
12. MUST NOT start a Delivery Queue, run `roundfix deliver`, merge, re-run or
    otherwise change anything on GitHub, and MUST NOT write to the live Run
    Database or artifact directory under `~/.roundfix`.

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

- `newest="$(find 'docs/specs/0228-queue-items-that-stay-current-with-main/qa' -type f -name "qa-report-*.md" -print 2>/dev/null | awk "{ report=\$0; name=\$0; parts=split(name, path, \"/\"); name=path[parts]; name=substr(name, 11, length(name)-13); date=substr(name, 1, 10); suffix=substr(name, 11); shape=(length(date) == 10 && substr(date, 5, 1) == \"-\" && substr(date, 8, 1) == \"-\"); for (i=1; i <= 10 && shape; i++) { if (i != 5 && i != 8 && index(\"0123456789\", substr(date, i, 1)) == 0) shape=0 } year=substr(date, 1, 4)+0; month=substr(date, 6, 2)+0; day=substr(date, 9, 2)+0; leap=(year%400 == 0 || (year%4 == 0 && year%100 != 0)); days=(month == 2 ? 28+leap : ((month == 4 || month == 6 || month == 9 || month == 11) ? 30 : 31)); dated=(shape && month >= 1 && month <= 12 && day >= 1 && day <= days); if (!dated) date=\"\"; if (suffix == \"\") { sequenced=1; sequence=-1 } else { sequenced=(length(suffix) > 1 && substr(suffix, 1, 1) == \"-\"); for (i=2; i <= length(suffix) && sequenced; i++) { if (index(\"0123456789\", substr(suffix, i, 1)) == 0) sequenced=0 } sequence=(sequenced ? substr(suffix, 2)+0 : 0) } if (dated && sequenced) printf \"%s\\t%s\\t%s\\t%s\\t%s\\n\", dated, date, sequenced, sequence, report }" | sort -k1,1n -k2,2 -k3,3n -k4,4n -k5,5 | tail -1 | cut -f5-)"; test -n "$newest" || exit 1; awk "BEGIN { whitespace=\" \t\r\n\f\v\" } NR == 1 && \$0 == \"---\" { frontmatter=1; next } frontmatter && \$0 == \"---\" { closed=1; exit } frontmatter && index(\$0, \"verdict:\") == 1 { verdict=substr(\$0, 9); while (length(verdict) > 0 && index(whitespace, substr(verdict, 1, 1)) > 0) verdict=substr(verdict, 2); while (length(verdict) > 0 && index(whitespace, substr(verdict, length(verdict), 1)) > 0) verdict=substr(verdict, 1, length(verdict)-1); verdicts++ } END { exit(closed && verdicts == 1 && length(verdict) > 0 ? 0 : 1) }" "$newest"`

## References

`_prd.md` → Goals; Core Features 1-4; Success Metric 1; Success Metric 2;
Success Metric 3; Success Metric 4; Acceptance evidence; `_techspec.md` →
Testing Approach 1-3; API Contract 1; API Contract 2; API Contract 3; Surface
Transcripts; Integration Points; ADR-0080; ADR-0088; ADR-0091; ADR-0104;
ADR-0155; ADR-0167; ADR-0233.
