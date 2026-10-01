---
task: task_03
spec: 0210-evidence-snapshots-that-stay-small
status: pending
type: qa
complexity: medium
---

# Task 03: Run the final QA gate

## Overview

The authored terminal gate for this Spec. It declares what the matrix covers
and settles the Spec on evidence. This Spec changes the Evidence Snapshot
recorder, its reader and the carry proof in `internal/speccheck`, four test
assertions in `internal/speccheck` and `internal/daemon`, and the
Context-Driven Development guide. Every behavior row is exercised by executing
the named tests against the built tree. No command reaches GitHub, a provider
or a live remote.

## Requirements

1. MUST run the repository Verification and record its result as a gate fact.
2. MUST verify, by executing task_01's new tests against the built tree:
   - a row whose glob matches 2,500 files records one line for that input,
     and the snapshot block's line count equals one plus three per row plus
     one per declared input;
   - the recorded digest equals an independently built `h1` summary digest;
   - a new-form record carries at an unchanged head, and a changed, added or
     removed file under an input re-runs the row with
     `input moved: <ref>`, naming the ref and not the files;
   - `Carriable` accepts a recorded pair against current files and refuses a
     differing pair.
3. MUST verify, by executing the governed `TestCarriable` and
   `TestMechanicalStageCarriable…` tests, every `TestRowCarry…` test and
   task_01's `TestCarryReadsAPerFileSnapshotAsItsDigest`, that a report in
   Spec 0202's per-file form is still read: a valid list carries, a list
   missing a file re-runs as moved, and a malformed list re-runs as
   `no evidence snapshot` with no error and no finding. It MUST confirm from
   the Git history that `internal/speccheck/mechanical_test.go`,
   `internal/speccheck/qa_row_carry_test.go` and
   `internal/daemon/qa_prior_pass_test.go` are unchanged by this Spec.
4. MUST verify, by executing `TestTwoGatePassesCarryAnUnmovedRowIntoTheCommittedReport`,
   `TestTwoGatePassesReRunARowWhoseInputMoved` and
   `TestQAGateCarriesARowFromAnUnintegratedFailedPass` against the built
   tree, that the prior pass import of ADR-0194 keeps a failed pass recorded
   in the new form and in the per-file form, and that the next pass carries
   its unmoved row. No public command runs the QA stage without a live ACP
   Runtime, so this row runs through the Daemon's own test seam.
5. MUST record, as evidence this Spec did not author:
   - Spec 0203's QA Report before commit `a1fc8402`, read with
     `git show a1fc8402^:docs/specs/0203-a-reviewer-that-validates-its-findings-and-remembers-its-rounds/qa/qa-report-2026-10-01.md`.
     Recount its lines, its `sha256:` lines, its recorded rows and their
     declared inputs, and compute the line count the new form gives the same
     rows. The PRD states 157,212 lines, 78,254 `sha256:` lines, 11 rows of
     7 inputs and 111 lines;
   - the documentation of `golang.org/x/mod/sumdb/dirhash.Hash1`
     (<https://pkg.go.dev/golang.org/x/mod/sumdb/dirhash#Hash1>), confirming
     the summary line format ADR-0210 adopts.

   When a source is unavailable, it MUST record that row as blocked with its
   reason.
6. MUST verify that the guide contains `one line per declared input` and
   `input moved: <refs>`, no longer contains `input moved: <paths>`, and
   still names every word of Spec 0202's Vocabulary Contract.
7. MUST verify that this Spec's own artifacts satisfy the promise rule, and
   check whether it introduced, changed or retired a glossary term.
   **Evidence Snapshot**'s definition in `CONTEXT.md` names no per-file
   shape; record whether it still describes the record.
8. MUST verify from the repository history that each Task's changed files stay
   within its declarations or its `## Recorded paths`, and that no Governed
   Path changed, as `_authorization.md` records `paths: []`. This row reads
   Task commits, so it declares a `commit_range` input.
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
11. MUST NOT write to the live Run Database under `~/.roundfix`.

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

- `newest="$(find 'docs/specs/0210-evidence-snapshots-that-stay-small/qa' -type f -name "qa-report-*.md" -print 2>/dev/null | awk "{ report=\$0; name=\$0; parts=split(name, path, \"/\"); name=path[parts]; name=substr(name, 11, length(name)-13); date=substr(name, 1, 10); suffix=substr(name, 11); shape=(length(date) == 10 && substr(date, 5, 1) == \"-\" && substr(date, 8, 1) == \"-\"); for (i=1; i <= 10 && shape; i++) { if (i != 5 && i != 8 && index(\"0123456789\", substr(date, i, 1)) == 0) shape=0 } year=substr(date, 1, 4)+0; month=substr(date, 6, 2)+0; day=substr(date, 9, 2)+0; leap=(year%400 == 0 || (year%4 == 0 && year%100 != 0)); days=(month == 2 ? 28+leap : ((month == 4 || month == 6 || month == 9 || month == 11) ? 30 : 31)); dated=(shape && month >= 1 && month <= 12 && day >= 1 && day <= days); if (!dated) date=\"\"; if (suffix == \"\") { sequenced=1; sequence=-1 } else { sequenced=(length(suffix) > 1 && substr(suffix, 1, 1) == \"-\"); for (i=2; i <= length(suffix) && sequenced; i++) { if (index(\"0123456789\", substr(suffix, i, 1)) == 0) sequenced=0 } sequence=(sequenced ? substr(suffix, 2)+0 : 0) } if (dated && sequenced) printf \"%s\\t%s\\t%s\\t%s\\t%s\\n\", dated, date, sequenced, sequence, report }" | sort -k1,1n -k2,2 -k3,3n -k4,4n -k5,5 | tail -1 | cut -f5-)"; test -n "$newest" || exit 1; awk "BEGIN { whitespace=\" \t\r\n\f\v\" } NR == 1 && \$0 == \"---\" { frontmatter=1; next } frontmatter && \$0 == \"---\" { closed=1; exit } frontmatter && index(\$0, \"verdict:\") == 1 { verdict=substr(\$0, 9); while (length(verdict) > 0 && index(whitespace, substr(verdict, 1, 1)) > 0) verdict=substr(verdict, 2); while (length(verdict) > 0 && index(whitespace, substr(verdict, length(verdict), 1)) > 0) verdict=substr(verdict, 1, length(verdict)-1); verdicts++ } END { exit(closed && verdicts == 1 && length(verdict) > 0 ? 0 : 1) }" "$newest"`

## References

`_prd.md` → Goals 1-3; Core Features 1-5; Success Metrics 1-4; Acceptance
evidence; `_techspec.md` → Testing Approach 1-3; API Contract 1; API Contract
2; Integration Points; ADR-0080; ADR-0088; ADR-0091; ADR-0097; ADR-0104;
ADR-0155; ADR-0167; ADR-0194; ADR-0195; ADR-0210.
