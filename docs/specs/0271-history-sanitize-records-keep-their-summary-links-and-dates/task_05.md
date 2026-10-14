---
task: task_05
spec: 0271-history-sanitize-records-keep-their-summary-links-and-dates
status: pending
type: qa
complexity: medium
---

# Task 05: Run the final QA gate

## Overview

The authored terminal gate for this Spec. It declares what the matrix covers
and settles the Spec on evidence. This Spec makes the History Sanitize Command
and the history section of `roundfix baseline update` write Archive Records
from Legacy Archive Folders that keep their summary, state their override
gaps, list their adopted sources, carry a creation date and rebase the links
of the outcome and of each promoted copy, while every record already written
and the Archive Command stay as they are. No row reads or writes the real
`~/.roundfix`, reaches a provider, reads a credential or reads another
repository.

## Requirements

1. MUST build the binary with `make build` and run the repository
   Verification, recording its result as a gate fact.
2. MUST verify Success Metric 1 through the built `bin/roundfix` with `HOME`
   set to a temporary home, `NODE_OPTIONS` unset and a temporary Git
   repository holding an annotated `history-full` tag and synthetic Legacy
   Archive Folders, one per defect as `_techspec.md` → Current behavior
   describes them: a two-line block quote summary, a title followed by a
   section, `qa_override: true` without approval or reason, a summary link to
   a promoted reference, a promoted copy linking `../_prd.md` and an ADR
   outside the folder, an adopted Backlog Entry and an adopted Finding without
   `_index.md`, and a PRD without a date beside a dated TechSpec plus one
   beside an undated TechSpec. It MUST run `roundfix history sanitize` (plan),
   then `--apply --batch <n>` with the promotions, and record for each record
   and copy that the defect is gone (API Contracts 1-5), that every rewritten
   destination resolves from its file to an existing path, and that each
   record converted without a promotion has the byte size the plan printed.
   The plan takes no `--promote`, so a promoted record's size is not compared
   with it; Requirement 3 compares the promoted bytes of the two commands.
   Path text naming a folder file that is neither a core artifact nor
   promoted stays as the Archive Relink's path pass leaves it (ADR-0277:
   "No other byte changes"); it still names that file at the record's
   `source_revision`, and the report records it as observed, not as a
   finding.
3. MUST repeat Requirement 2's apply through `roundfix baseline update --yes
   --promote <path>` on an identical fixture adopted by the Baseline, and
   record that it writes the same record and copy bytes.
4. MUST verify Success Metric 2 and Invariants 1-3 against the delivery base,
   the merge base of the audited head with `main`: build a second binary from
   a temporary worktree at that commit, convert a folder without any of the
   six shapes and a promoted file without relative links with both binaries
   on identical fixtures, and record that the records and copies are
   byte-identical.
5. MUST verify Success Metric 3 and API Contract 6: record that
   `git diff --name-only --diff-filter=MD <delivery base> <audited head> -- docs/history/specs`
   is empty, so every Archive Record present at the delivery base is
   byte-identical at the audited head, and run the Archive Command's tests
   (`go test -count=1 -run 'Archive' ./internal/spec ./internal/cli`),
   recording that they pass with no test file of the Archive Command changed.
6. MUST execute the tests task_01, task_02 and task_03 name in their
   Verification.
7. MUST record, as evidence this Spec did not author, each source of
   `_prd.md` → Acceptance evidence: the CommonMark Spec 0.31.2 section 5.1
   passage on the block quote marker and the RFC 3986 section 4.2 passage on
   relative references, quoting each and stating whether it still supports
   the design; for a source it cannot reach it records
   `blocked (environment: network denied: <host>)`.
8. MUST perform the glossary check of `docs/agents/domain.md`: record that
   `CONTEXT.md` defines **Legacy Link Rebase** and that the revised
   **Lenient Legacy Reading** names the new readings, and whether any other
   term the Spec introduced needs one.
9. MUST verify from the repository history that each Task's changed files stay
   within its declarations or its `## Recorded paths`; that every Governed Path
   changed is bounded in `_authorization.md`; that the Makefile, the CI
   workflows, `go.mod`, `go.sum`, every Baseline module and
   `internal/cli/review_conventions.go` did not change; that the Roundfix
   skill version was raised and recorded; and that no `### QA settlement`
   section changed. This row reads Task commits, so it declares a
   `commit_range` input.
10. MUST record the non-waivable Pull Request row as
    `blocked (environment: no open Pull Request)`, naming the Pull Request row
    in its provenance, and support it with the pre-PR equivalent of each
    control: approval from the recorded delivery authority, checks from the
    Daemon's repository Verification at the audited head, no unresolved
    threads because no Pull Request exists, no Merge-Ready acceptance yet, and
    the audited head named as the claimed candidate.
11. MUST NOT accept a row whose only evidence is that a file was read. It MUST
    NOT read or write `~/.roundfix`, record any environment value, send any
    request to a provider, or start a workflow.

## Subtasks

- [ ] Build the binary and run the repository Verification.
- [ ] Convert the six synthetic shapes through both commands.
- [ ] Compare unchanged shapes and existing records with the delivery base.
- [ ] Record the outside evidence, the glossary check and the scope audit, and write the QA report.

## Acceptance Criteria

- [ ] Each of the six defects is absent from what both commands write.
- [ ] Shapes the Spec does not touch convert byte-identically to the delivery base.
- [ ] No Archive Record present at the delivery base changed.
- [ ] The scope audit finds no undeclared Governed Path change.

## Verification

The following command is rendered unchanged from Roundfix's derived QA contract.
It proves a newest report exists and its verdict is readable; eligibility is
applied in process by whoever settles the Task.

- `newest="$(find 'docs/specs/0271-history-sanitize-records-keep-their-summary-links-and-dates/qa' -type f -name "qa-report-*.md" -print 2>/dev/null | awk "{ report=\$0; name=\$0; parts=split(name, path, \"/\"); name=path[parts]; name=substr(name, 11, length(name)-13); date=substr(name, 1, 10); suffix=substr(name, 11); shape=(length(date) == 10 && substr(date, 5, 1) == \"-\" && substr(date, 8, 1) == \"-\"); for (i=1; i <= 10 && shape; i++) { if (i != 5 && i != 8 && index(\"0123456789\", substr(date, i, 1)) == 0) shape=0 } year=substr(date, 1, 4)+0; month=substr(date, 6, 2)+0; day=substr(date, 9, 2)+0; leap=(year%400 == 0 || (year%4 == 0 && year%100 != 0)); days=(month == 2 ? 28+leap : ((month == 4 || month == 6 || month == 9 || month == 11) ? 30 : 31)); dated=(shape && month >= 1 && month <= 12 && day >= 1 && day <= days); if (!dated) date=\"\"; if (suffix == \"\") { sequenced=1; sequence=-1 } else { sequenced=(length(suffix) > 1 && substr(suffix, 1, 1) == \"-\"); for (i=2; i <= length(suffix) && sequenced; i++) { if (index(\"0123456789\", substr(suffix, i, 1)) == 0) sequenced=0 } sequence=(sequenced ? substr(suffix, 2)+0 : 0) } if (dated && sequenced) printf \"%s\\t%s\\t%s\\t%s\\t%s\\n\", dated, date, sequenced, sequence, report }" | sort -k1,1n -k2,2 -k3,3n -k4,4n -k5,5 | tail -1 | cut -f5-)"; test -n "$newest" || exit 1; awk "BEGIN { whitespace=\" \t\r\n\f\v\" } NR == 1 && \$0 == \"---\" { frontmatter=1; next } frontmatter && \$0 == \"---\" { closed=1; exit } frontmatter && index(\$0, \"verdict:\") == 1 { verdict=substr(\$0, 9); while (length(verdict) > 0 && index(whitespace, substr(verdict, 1, 1)) > 0) verdict=substr(verdict, 2); while (length(verdict) > 0 && index(whitespace, substr(verdict, length(verdict), 1)) > 0) verdict=substr(verdict, 1, length(verdict)-1); verdicts++ } END { exit(closed && verdicts == 1 && length(verdict) > 0 ? 0 : 1) }" "$newest"`

## References

- `_prd.md` → Goals; User Stories 1-3; Core Features 1-4; Success Metric 1; Success Metric 2; Success Metric 3; Acceptance evidence
- `_techspec.md` → API Contract 1; API Contract 2; API Contract 3; API Contract 4; API Contract 5; API Contract 6; Invariants 1-6; Current behavior; Build Order 5
- ADR-0277; ADR-0251; ADR-0260; ADR-0261
