---
task: task_05
spec: 0246-a-sanitize-that-reads-older-folders-and-names-its-refusals
status: pending
type: qa
complexity: high
---

# Task 05: Run the final QA gate

## Overview

The authored terminal gate for this Spec. It declares what the matrix covers
and settles the Spec on evidence. This Spec changes three things:

- The History Sanitize Command reads a Legacy Archive Folder leniently.
- It records a failed QA without an override as `failed-qa`.
- It lists every Refused Unit while a batch converts the next units it can
  convert.

Every behavior row runs against the built tree with temporary repositories, a
temporary home and synthetic legacy folders. No row reaches a provider, starts
a real Agent Session, reads a credential, writes the real `~/.roundfix` or
writes the adopter repository. No row sends adopter content to Jev.

## Requirements

1. MUST run the repository Verification and record its result as a gate fact.
2. MUST verify API Contract 1 and Success Metric 4 by executing task_02's five
   spec tests. It MUST also confirm from the diff that `spec.Load`,
   `ReadCauseGraph`, `ReadCauseGraphAt` and the Archive Command never set or
   read the lenient path.
3. MUST verify API Contract 2 and Success Metric 2 by executing task_03's five
   spec tests. It MUST confirm from the diff that each reader listed in
   `_techspec.md` → "Readers that branch on disposition" behaves as that list
   states, and that only `ArchivedTaskCompleted` gained a disposition branch.
4. MUST verify API Contract 3, Success Metric 1 and Success Metric 3 by
   executing task_04's CLI tests. It MUST also run Surface Transcript 1,
   Surface Transcript 2 and Surface Transcript 3 through the built binary, each
   in its own temporary repository built as the transcript describes, and
   record the three outputs and exit codes.
5. MUST re-measure the adopter's failure shape with the built binary. It builds
   synthetic folders of the four shapes and a clean folder in a temporary
   repository with an annotated `history-full` tag. It records the plan and
   one `--apply --batch 5`, and reads each record back through a freshly built
   binary's `roundfix history sanitize` and through `ParseArchiveRecord` in a
   test. It records that the same fixture refuses at `405274c6`, or cites the
   TechSpec's Current behavior section when that binary is unavailable.
6. MUST verify Success Metric 5. Every existing `./internal/spec` and
   `./internal/cli` history, archive and archive record test passes, and only
   `TestHistorySanitizePreflightsWholeBatch` changed its expectation.
7. MUST record, as evidence this Spec did not author, each source of
   `_prd.md` → Acceptance evidence: the adopter's triaged report in the
   Secondbrain, the adopter's four folders, read only and never copied, and
   the two published sources. For each source it reaches it MUST record what
   it read and whether it still supports the design. For each source it cannot
   reach it MUST record the row as blocked with the reason.
8. MUST verify that the history command reference and the Roundfix Skill's
   archive reference describe the behavior, that each skill mirror equals its
   canonical file, and that the raised owned-skill version is recorded.
9. MUST perform the glossary check of `docs/agents/domain.md`. It records that
   `CONTEXT.md` defines **Refused Unit** and **Lenient Legacy Reading** and
   names `failed-qa` in **Archive Record**, and whether any other term the
   Spec introduced needs one.
10. MUST verify from the repository history:
    - that each Task's changed files stay within its declarations or its
      `## Recorded paths`;
    - that every Governed Path changed is bounded in `_authorization.md`;
    - that `Makefile`, `go.mod`, `.roundfixrc.yml`, the CI workflows and the
      Baseline modules did not change.

    This row reads Task commits, so it declares a `commit_range` input.
11. MUST record the non-waivable Pull Request row as
    `blocked (environment: no open Pull Request)`, naming the Pull Request row
    in its provenance. It supports the row with the pre-PR equivalent of each
    control:
    - approval: the maintainer's recorded delivery authority;
    - checks and status: the Daemon's repository Verification at the audited
      head;
    - unresolved threads: none, because no Pull Request exists;
    - Merge-Ready acceptance: none yet, because the pre-PR review follows
      archive;
    - review-artifact ancestry: the audited head named as the claimed
      candidate.
12. MUST NOT accept a row whose only evidence is that a file was read. It MUST
    NOT modify the adopter repository or copy its content, and MUST NOT send
    any request to a provider.

## Subtasks

- [ ] Build the binary and run the repository Verification.
- [ ] Exercise the lenient reading, `failed-qa` and Refused Unit rows.
- [ ] Run the three Surface Transcripts and re-measure the adopter's shapes.
- [ ] Record the outside evidence, glossary check and scope audit.
- [ ] Write the QA report with the Pull Request row blocked.

## Acceptance Criteria

- [ ] Every behavior row passes on recorded evidence.
- [ ] The three Surface Transcripts match through the built binary.
- [ ] The four adopter shapes convert in one apply.
- [ ] The scope audit finds no undeclared Governed Path change.

## Verification

The following command is rendered unchanged from Roundfix's derived QA contract.
It proves a newest report exists and its verdict is readable; eligibility is
applied in process by whoever settles the Task.

- `newest="$(find 'docs/specs/0246-a-sanitize-that-reads-older-folders-and-names-its-refusals/qa' -type f -name "qa-report-*.md" -print 2>/dev/null | awk "{ report=\$0; name=\$0; parts=split(name, path, \"/\"); name=path[parts]; name=substr(name, 11, length(name)-13); date=substr(name, 1, 10); suffix=substr(name, 11); shape=(length(date) == 10 && substr(date, 5, 1) == \"-\" && substr(date, 8, 1) == \"-\"); for (i=1; i <= 10 && shape; i++) { if (i != 5 && i != 8 && index(\"0123456789\", substr(date, i, 1)) == 0) shape=0 } year=substr(date, 1, 4)+0; month=substr(date, 6, 2)+0; day=substr(date, 9, 2)+0; leap=(year%400 == 0 || (year%4 == 0 && year%100 != 0)); days=(month == 2 ? 28+leap : ((month == 4 || month == 6 || month == 9 || month == 11) ? 30 : 31)); dated=(shape && month >= 1 && month <= 12 && day >= 1 && day <= days); if (!dated) date=\"\"; if (suffix == \"\") { sequenced=1; sequence=-1 } else { sequenced=(length(suffix) > 1 && substr(suffix, 1, 1) == \"-\"); for (i=2; i <= length(suffix) && sequenced; i++) { if (index(\"0123456789\", substr(suffix, i, 1)) == 0) sequenced=0 } sequence=(sequenced ? substr(suffix, 2)+0 : 0) } if (dated && sequenced) printf \"%s\\t%s\\t%s\\t%s\\t%s\\n\", dated, date, sequenced, sequence, report }" | sort -k1,1n -k2,2 -k3,3n -k4,4n -k5,5 | tail -1 | cut -f5-)"; test -n "$newest" || exit 1; awk "BEGIN { whitespace=\" \t\r\n\f\v\" } NR == 1 && \$0 == \"---\" { frontmatter=1; next } frontmatter && \$0 == \"---\" { closed=1; exit } frontmatter && index(\$0, \"verdict:\") == 1 { verdict=substr(\$0, 9); while (length(verdict) > 0 && index(whitespace, substr(verdict, 1, 1)) > 0) verdict=substr(verdict, 2); while (length(verdict) > 0 && index(whitespace, substr(verdict, length(verdict), 1)) > 0) verdict=substr(verdict, 1, length(verdict)-1); verdicts++ } END { exit(closed && verdicts == 1 && length(verdict) > 0 ? 0 : 1) }" "$newest"`

## References

- `_prd.md` → Goals; Core Features 1-4; Success Metric 1; Success Metric 2; Success Metric 3; Success Metric 4; Success Metric 5; Acceptance evidence
- `_techspec.md` → API Contract 1; API Contract 2; API Contract 3; Surface Transcripts; Build Order 5
- ADR-0251; ADR-0248; ADR-0154; ADR-0184
