---
task: task_05
spec: 0227-reconcile-releases-the-runs-of-merged-specs
status: pending
type: qa
complexity: medium
---

# Task 05: Run the final QA gate

## Overview

The authored terminal gate for this Spec. It declares what the matrix covers
and settles the Spec on evidence. This Spec lets merge evidence release a
terminal Run of a merged Spec whose files diverged later, lets the
post-merge cleanup prove a squash merge onto a moved default branch, lets
`roundfix reconcile` list and release the item branches of merged Specs, and
describes all three in the Roundfix Skill and the `reconcile` and `deliver`
guides. Every behavior row is exercised against the built tree with
disposable repositories and a disposable Run Database. No command reaches
GitHub, a provider or the network, except the reads of published pages the
outside-evidence row names.

## Requirements

1. MUST run the repository Verification and record its result as a gate fact.
2. MUST verify API Contract 1 and Success Metrics 1 and 2 by executing
   task_01's tests against the built tree, and MUST additionally rebuild the
   PRD's measured fixture with the built `roundfix` binary in a disposable
   repository and Roundfix Home (a Run that inherited an operator commit from
   its item branch, no merge record, a later default-branch commit editing one
   of its files and renaming another), recording that a read-only
   `roundfix reconcile` now reports `superseded` where the binary of the
   starting tree reported "1 Run-only file, 1 differing shared file", and that
   a dirty path outside the declared scope still reports `dirty`.
3. MUST verify API Contract 2 and Success Metric 3 by executing task_02's
   tests, and MUST additionally record, read-only against this repository's
   own history, that `git merge-tree --write-tree` of the first parent of
   merge commits `f0353479` (0224) and `69ace45d` (0226) with candidate heads
   `88e2fe3c` and `e62588bd` equals each merge commit's tree, and that neither
   candidate is an ancestor of, nor has the tree of, its merge commit.
4. MUST verify API Contract 3, Surface Transcript 1 and Success Metric 4 by
   executing task_03's tests and by reproducing Surface Transcript 1 through
   the built binary in a disposable repository holding an item branch of a
   merged Spec, comparing stdout, stderr and exit code under the transcript
   conventions, then running `--apply` there and recording that the branch is
   gone while a live item's branch stays.
5. MUST verify Success Metric 5: every existing test of `internal/worktree`
   and of the reconcile and release suites of `internal/cli` passes, and the
   only changed expectation is the renamed
   `TestArchivedSpecWithMergeEvidenceSupersedesOtherCommitPaths`, read from
   the diff; `internal/cli/cli_test.go` is unchanged.
6. MUST verify that the Roundfix Skill's `reconcile` reference and the
   `reconcile` and `deliver` guides describe the three rules and no longer
   carry the sentences task_04 removes, that each mirror equals its canonical
   file, that the version was raised and recorded, and that the
   `### QA settlement` section of every skill is unchanged.
7. MUST record, as evidence this Spec did not author, each source of
   `_prd.md` → Acceptance evidence: the operator's intervention log entry
   144; the merged 0223, 0224 and 0226 items in the live Run Database, read
   only with `sqlite3 -readonly`; the delivery commits of 0181, 0184, 0200,
   0204, 0205, 0207 and 0217 on this repository's default branch; the Fiscus
   Backlog Entry in the Secondbrain mirror; Git's `git merge-tree`
   documentation; and the published squash-merge cleanup page. For each
   source it reaches it MUST record what it read and whether it still
   supports the design; for each it cannot reach it MUST record the row as
   blocked with the reason.
8. MUST verify that this Spec's own artifacts satisfy the promise rule, and
   perform the glossary check of `docs/agents/domain.md`: record whether
   **Run Worktree Reconciliation** and **Reconcile Command** still describe
   the behavior, and whether merge evidence, the delivery commit or the item
   branch needs a glossary term, as candidates for the maintainer.
9. MUST verify from the repository history that each Task's changed files stay
   within its declarations or its `## Recorded paths`, that every Governed
   Path changed is bounded in `_authorization.md`, and that `Makefile`,
   `.roundfixrc.yml`, `go.mod`, `internal/cli/cli_test.go` and the CI
   workflows did not change. This row reads Task commits, so it declares a
   `commit_range` input.
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
12. MUST NOT start a Delivery Queue, run `roundfix deliver` or
    `roundfix reconcile --apply` outside a disposable repository and Roundfix
    Home, merge, or change anything on GitHub, and MUST NOT write to the live
    Run Database under `~/.roundfix`.

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

- `newest="$(find 'docs/specs/0227-reconcile-releases-the-runs-of-merged-specs/qa' -type f -name "qa-report-*.md" -print 2>/dev/null | awk "{ report=\$0; name=\$0; parts=split(name, path, \"/\"); name=path[parts]; name=substr(name, 11, length(name)-13); date=substr(name, 1, 10); suffix=substr(name, 11); shape=(length(date) == 10 && substr(date, 5, 1) == \"-\" && substr(date, 8, 1) == \"-\"); for (i=1; i <= 10 && shape; i++) { if (i != 5 && i != 8 && index(\"0123456789\", substr(date, i, 1)) == 0) shape=0 } year=substr(date, 1, 4)+0; month=substr(date, 6, 2)+0; day=substr(date, 9, 2)+0; leap=(year%400 == 0 || (year%4 == 0 && year%100 != 0)); days=(month == 2 ? 28+leap : ((month == 4 || month == 6 || month == 9 || month == 11) ? 30 : 31)); dated=(shape && month >= 1 && month <= 12 && day >= 1 && day <= days); if (!dated) date=\"\"; if (suffix == \"\") { sequenced=1; sequence=-1 } else { sequenced=(length(suffix) > 1 && substr(suffix, 1, 1) == \"-\"); for (i=2; i <= length(suffix) && sequenced; i++) { if (index(\"0123456789\", substr(suffix, i, 1)) == 0) sequenced=0 } sequence=(sequenced ? substr(suffix, 2)+0 : 0) } if (dated && sequenced) printf \"%s\\t%s\\t%s\\t%s\\t%s\\n\", dated, date, sequenced, sequence, report }" | sort -k1,1n -k2,2 -k3,3n -k4,4n -k5,5 | tail -1 | cut -f5-)"; test -n "$newest" || exit 1; awk "BEGIN { whitespace=\" \t\r\n\f\v\" } NR == 1 && \$0 == \"---\" { frontmatter=1; next } frontmatter && \$0 == \"---\" { closed=1; exit } frontmatter && index(\$0, \"verdict:\") == 1 { verdict=substr(\$0, 9); while (length(verdict) > 0 && index(whitespace, substr(verdict, 1, 1)) > 0) verdict=substr(verdict, 2); while (length(verdict) > 0 && index(whitespace, substr(verdict, length(verdict), 1)) > 0) verdict=substr(verdict, 1, length(verdict)-1); verdicts++ } END { exit(closed && verdicts == 1 && length(verdict) > 0 ? 0 : 1) }" "$newest"`

## References

`_prd.md` → Goals; Core Features 1-4; Success Metrics 1-5; Acceptance
evidence; `_techspec.md` → Testing Approach 1-4; API Contract 1; API
Contract 2; API Contract 3; Surface Transcript 1; Integration Points;
ADR-0080; ADR-0088; ADR-0091; ADR-0104; ADR-0155; ADR-0156; ADR-0232.
