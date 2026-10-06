---
task: task_05
spec: 0235-one-qa-partial-policy
status: failed
type: qa
complexity: high
---

# Task 05: Run the final QA gate

## Overview

The authored terminal gate for this Spec. It declares what the matrix covers
and settles the Spec on evidence. This Spec makes Daemon settlement,
`roundfix settle`, `roundfix archive` and `roundfix qa-report accept` apply
one QA partial policy. A `partial` whose only unmet rows are pre-PR Pull
Request rows, network-denied outside-evidence rows and covered declared rows
qualifies. The Delivery Queue parks only partials that need the override,
settle names the report it refused, and the skills and guides state the
policy. Every behavior row runs against the built tree in temporary
repositories and homes. No row reaches the network on purpose or reads or
writes the real `~/.roundfix`.

## Requirements

1. MUST run the repository Verification and record its result as a gate fact.
2. MUST verify Success Metric 1 by running the built `roundfix qa-report
   accept` on the newest archived QA Report of every archived Spec under
   `docs/history/specs/` whose verdict is `partial`. Compare each result with
   the measurement in `_prd.md` → Acceptance evidence. 0213 and 0220 must
   now be accepted, 0227 must now be refused only for its non-Pull-Request
   environment row, and every other result must match the measurement. The
   row records the full before and after table.
3. MUST verify Success Metric 2 and API Contract 1 by executing task_01's and
   task_02's policy tests, and by reproducing Surface Transcript 1 and
   Surface Transcript 2 through the built `roundfix qa-report accept` on
   reports written into a disposable Spec directory.
4. MUST verify Success Metric 3 and API Contract 5 by executing
   `TestRunSpecDoesNotParkAQualifyingPartialAsEnvironmentOnly` and
   `TestRunSpecReportsAnEnvironmentOnlyPartialFromTheRunBranch`, and confirm
   from the diff that `qaEnvironmentPartial` reads
   `EnvironmentRowsNeedingOverride` and adds no predicate of its own.
5. MUST verify Success Metric 4 and API Contract 4 by executing
   `TestSettleNamesTheReportItRefused`. It MUST also reproduce Surface
   Transcript 3 through the built `roundfix settle` on a disposable
   repository whose QA Task's newest report is refused.
6. MUST verify API Contract 2, API Contract 3 and Invariants 1, 2, 6 and 10
   from the tests and the diff: the two new messages, the skipped-row
   refusal, the noted Pull Request provenance, the marker's exact form, and
   no caller adding a condition of its own.
7. MUST verify Success Metric 5. Every existing QA, archive, settle, delivery
   and Baseline test passes, and the only changed expectations are those
   task_01, task_02 and task_04 declare.
8. MUST verify that the qa-gate, archive-spec and Roundfix skills, the
   Roundfix archive and settle references, the four user guides and both
   generated Baseline guides state the one policy and the marker, that the
   `### QA settlement` section is identical in the three skills, that the
   Vocabulary Contract's pattern is documented where it says, that each
   skill mirror equals its canonical file, and that every raised owned-skill
   version is recorded.
9. MUST record, as evidence this Spec did not author, each source of
   `_prd.md` → Acceptance evidence. These are the archived QA Reports of
   this repository (Requirement 2), the Fluxus report mirrored in the
   Secondbrain, the Oraculum report in the Secondbrain inbox, and JUnit's
   published assumption contract. For each source it reaches, it MUST
   record what it read and whether it still supports the design. For the
   Fluxus report it MUST record the built binary's verdict on both
   `qa-report-2026-10-05.md` and `qa-report-2026-10-05-01.md`. For a source
   the Run sandbox cannot reach because network access was denied, it MUST
   record the row with this Spec's own marker,
   `blocked (environment: network denied: <host>)`, with an
   `outside-evidence row` item in its provenance. Any other unreachable
   source is recorded as an ordinary environment-blocked row with its
   reason.
10. MUST verify that this Spec's own artifacts satisfy the promise rule, and
    perform the glossary check of `docs/agents/domain.md`. It records
    whether `CONTEXT.md`'s **QA Report** entry, which names only the Pull
    Request row as exempt, needs the network-denied outside-evidence row,
    and whether "network-denied outside-evidence row" needs a glossary term.
11. MUST verify from the repository history that each Task's changed files
    stay within its declarations or its `## Recorded paths`, that every
    Governed Path changed is bounded in `_authorization.md`, that no derived
    Baseline file was edited outside `make baseline-digests` and the Managed
    Refresh, and that `Makefile`, `go.mod` and the CI workflows did not
    change. This row reads Task commits, so it declares a `commit_range`
    input.
12. MUST record the non-waivable Pull Request row as
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
13. MUST NOT accept a row whose only evidence is that a file was read.

## Subtasks

- [ ] Build the binary and run the repository Verification.
- [ ] Replay the archived partial reports and the Fluxus reports through the built binary.
- [ ] Reproduce the three Surface Transcripts.
- [ ] Record the outside evidence, glossary check and scope audit.
- [ ] Write the QA report with the Pull Request row blocked.

## Acceptance Criteria

- [ ] Every behavior row passes on recorded evidence.
- [ ] The archived-report replay changes exactly the verdicts the PRD
      predicts.
- [ ] The three Surface Transcripts match through the built binary.
- [ ] The scope audit finds no undeclared Governed Path change.

## Verification

The following command is rendered unchanged from Roundfix's derived QA contract.
It proves a newest report exists and its verdict is readable; eligibility is
applied in process by whoever settles the Task.

- `newest="$(find 'docs/specs/0235-one-qa-partial-policy/qa' -type f -name "qa-report-*.md" -print 2>/dev/null | awk "{ report=\$0; name=\$0; parts=split(name, path, \"/\"); name=path[parts]; name=substr(name, 11, length(name)-13); date=substr(name, 1, 10); suffix=substr(name, 11); shape=(length(date) == 10 && substr(date, 5, 1) == \"-\" && substr(date, 8, 1) == \"-\"); for (i=1; i <= 10 && shape; i++) { if (i != 5 && i != 8 && index(\"0123456789\", substr(date, i, 1)) == 0) shape=0 } year=substr(date, 1, 4)+0; month=substr(date, 6, 2)+0; day=substr(date, 9, 2)+0; leap=(year%400 == 0 || (year%4 == 0 && year%100 != 0)); days=(month == 2 ? 28+leap : ((month == 4 || month == 6 || month == 9 || month == 11) ? 30 : 31)); dated=(shape && month >= 1 && month <= 12 && day >= 1 && day <= days); if (!dated) date=\"\"; if (suffix == \"\") { sequenced=1; sequence=-1 } else { sequenced=(length(suffix) > 1 && substr(suffix, 1, 1) == \"-\"); for (i=2; i <= length(suffix) && sequenced; i++) { if (index(\"0123456789\", substr(suffix, i, 1)) == 0) sequenced=0 } sequence=(sequenced ? substr(suffix, 2)+0 : 0) } if (dated && sequenced) printf \"%s\\t%s\\t%s\\t%s\\t%s\\n\", dated, date, sequenced, sequence, report }" | sort -k1,1n -k2,2 -k3,3n -k4,4n -k5,5 | tail -1 | cut -f5-)"; test -n "$newest" || exit 1; awk "BEGIN { whitespace=\" \t\r\n\f\v\" } NR == 1 && \$0 == \"---\" { frontmatter=1; next } frontmatter && \$0 == \"---\" { closed=1; exit } frontmatter && index(\$0, \"verdict:\") == 1 { verdict=substr(\$0, 9); while (length(verdict) > 0 && index(whitespace, substr(verdict, 1, 1)) > 0) verdict=substr(verdict, 2); while (length(verdict) > 0 && index(whitespace, substr(verdict, length(verdict), 1)) > 0) verdict=substr(verdict, 1, length(verdict)-1); verdicts++ } END { exit(closed && verdicts == 1 && length(verdict) > 0 ? 0 : 1) }" "$newest"`

## References

- `_prd.md` → Goals; Core Features 1-6; Success Metric 1; Success Metric 2; Success Metric 3; Success Metric 4; Success Metric 5; Acceptance evidence
- `_techspec.md` → API Contract 1; API Contract 2; API Contract 3; API Contract 4; API Contract 5; Surface Transcript 1; Surface Transcript 2; Surface Transcript 3; Vocabulary Contract; Build Order 5
- ADR-0240; ADR-0080; ADR-0091; ADR-0104
