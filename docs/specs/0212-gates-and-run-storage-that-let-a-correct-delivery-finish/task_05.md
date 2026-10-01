---
task: task_05
spec: 0212-gates-and-run-storage-that-let-a-correct-delivery-finish
status: pending
type: qa
complexity: medium
---

# Task 05: Run the final QA gate

## Overview

The authored terminal gate for this Spec. It declares what the matrix covers
and settles the Spec on evidence. This Spec changes the pre-PR review's diff,
bound and record, the failed-pass QA import, the Task Context kinds with a
settlement check, and reconcile's classification of a merged Spec's Runs,
and describes them in the Roundfix, qa-gate and write-tasks skills and two
command guides. Every behavior row is exercised through the built binary or
by executing the named tests against the built tree, with disposable
repositories, a disposable Roundfix Home and fake runtimes. No command
reaches a provider, GitHub or the network, except the reads of published or
recorded sources the outside-evidence row names.

## Requirements

1. MUST run the repository Verification and record its result as a gate fact.
2. MUST verify, by executing task_01's tests against the built tree, API
   Contracts 1 and 2: the omitted classes and their reasons, the QA Report
   kept, the malformed-lock block, the bound block with no provider call,
   the stderr tail and the older record; and MUST recompute, outside Go,
   that `git diff --no-ext-diff --no-textconv --no-color 64aff3f7 954ad599`
   minus the paths task_01 omits is no larger than the bound, using the
   commit pairs of `_techspec.md` → Measured outside evidence.
3. MUST verify, by executing task_02's tests against the built tree, API
   Contract 3, and that a failed pass holding a Go file under `qa/evidence/`
   no longer reaches the next pass's repository gate.
4. MUST verify, by executing task_03's tests against the built tree, API
   Contract 4, and that `roundfix spec check --strict` of a disposable Spec
   whose Task declares `- deletes: <path>` for a removed path reports no
   `SC-REF-UNRESOLVED`.
5. MUST verify, by executing task_04's tests against the built tree, API
   Contract 5, and through the built binary that `roundfix reconcile` and
   `roundfix reconcile --apply` in a disposable repository and Roundfix Home
   classify and remove a merged Spec's Run with declared leftovers, and keep
   a Run with an undeclared leftover.
6. MUST verify that the Roundfix Skill's `review` and `reconcile` references,
   the `review` and `reconcile` command guides, the qa-gate skill's
   `### Runnable evidence` heading and the write-tasks skill and template
   describe the behavior, that each mirror equals its canonical file, that
   each raised version is recorded, and that the `### QA settlement` section
   of every skill is unchanged.
7. MUST record, as evidence this Spec did not author, each source of
   `_prd.md` → Acceptance evidence: the two review records and their commit
   pairs, re-measured with `git diff`; Spec 0203's archived evidence files;
   Spec 0200's archived task_04; the Vortex note read through the
   Secondbrain; and the read-only count and size of this repository's Run
   Worktrees of merged Specs, measured without opening the Run Database for
   writing. For each source it reaches it MUST record what it read and
   whether it still supports the design; for each it cannot reach it MUST
   record the row as blocked with the reason.
8. MUST verify that this Spec's own artifacts satisfy the promise rule, and
   perform the glossary check of `docs/agents/domain.md`: record whether
   **Recorded Path**, **Reconcile Command** and **Evidence Snapshot** still
   describe the behavior, and whether `deletes` or the review's omitted
   paths need a term.
9. MUST verify from the repository history that each Task's changed files stay
   within its declarations or its `## Recorded paths`, that every Governed
   Path changed is bounded in `_authorization.md`, and that `Makefile`,
   `.roundfixrc.yml`, `go.mod`, `internal/speccheck/coherence.go` and the CI
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
12. MUST NOT run `roundfix reconcile --apply` or `roundfix review` against
    this repository's live Runs or Artifact Directory, and MUST NOT write to
    the live Run Database under `~/.roundfix`.

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

- `newest="$(find 'docs/specs/0212-gates-and-run-storage-that-let-a-correct-delivery-finish/qa' -type f -name "qa-report-*.md" -print 2>/dev/null | awk "{ report=\$0; name=\$0; parts=split(name, path, \"/\"); name=path[parts]; name=substr(name, 11, length(name)-13); date=substr(name, 1, 10); suffix=substr(name, 11); shape=(length(date) == 10 && substr(date, 5, 1) == \"-\" && substr(date, 8, 1) == \"-\"); for (i=1; i <= 10 && shape; i++) { if (i != 5 && i != 8 && index(\"0123456789\", substr(date, i, 1)) == 0) shape=0 } year=substr(date, 1, 4)+0; month=substr(date, 6, 2)+0; day=substr(date, 9, 2)+0; leap=(year%400 == 0 || (year%4 == 0 && year%100 != 0)); days=(month == 2 ? 28+leap : ((month == 4 || month == 6 || month == 9 || month == 11) ? 30 : 31)); dated=(shape && month >= 1 && month <= 12 && day >= 1 && day <= days); if (!dated) date=\"\"; if (suffix == \"\") { sequenced=1; sequence=-1 } else { sequenced=(length(suffix) > 1 && substr(suffix, 1, 1) == \"-\"); for (i=2; i <= length(suffix) && sequenced; i++) { if (index(\"0123456789\", substr(suffix, i, 1)) == 0) sequenced=0 } sequence=(sequenced ? substr(suffix, 2)+0 : 0) } if (dated && sequenced) printf \"%s\\t%s\\t%s\\t%s\\t%s\\n\", dated, date, sequenced, sequence, report }" | sort -k1,1n -k2,2 -k3,3n -k4,4n -k5,5 | tail -1 | cut -f5-)"; test -n "$newest" || exit 1; awk "BEGIN { whitespace=\" \t\r\n\f\v\" } NR == 1 && \$0 == \"---\" { frontmatter=1; next } frontmatter && \$0 == \"---\" { closed=1; exit } frontmatter && index(\$0, \"verdict:\") == 1 { verdict=substr(\$0, 9); while (length(verdict) > 0 && index(whitespace, substr(verdict, 1, 1)) > 0) verdict=substr(verdict, 2); while (length(verdict) > 0 && index(whitespace, substr(verdict, length(verdict), 1)) > 0) verdict=substr(verdict, 1, length(verdict)-1); verdicts++ } END { exit(closed && verdicts == 1 && length(verdict) > 0 ? 0 : 1) }" "$newest"`

## References

`_prd.md` → Goals; User Stories 1-5; Core Features 1-8; Success Metrics 1-5;
Acceptance evidence; `_techspec.md` → Testing Approach 1-4; API Contract 1;
API Contract 2; API Contract 3; API Contract 4; API Contract 5; Measured
outside evidence; Integration Points; ADR-0080; ADR-0088; ADR-0091; ADR-0104;
ADR-0155; ADR-0167; ADR-0212.
