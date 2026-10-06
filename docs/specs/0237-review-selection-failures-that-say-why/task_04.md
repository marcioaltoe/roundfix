---
task: task_04
spec: 0237-review-selection-failures-that-say-why
status: pending
type: qa
complexity: high
---

# Task 04: Run the final QA gate

## Overview

The authored terminal gate for this Spec. It declares what the matrix covers
and settles the Spec on evidence. This Spec makes the ACPX Runner keep the
protocol step and the adapter's message of a prompt process that failed
without a result, and makes `roundfix review` name them and retry a failure
before the prompt once on the same selection. Every behavior row runs against
the built tree with fake runners, a fake acpx or a fake ACP adapter, and a
temporary home. No row reaches a provider, reads a credential, or reads or
writes the real `~/.roundfix` or `~/.acpx`.

## Requirements

1. MUST run the repository Verification and record its result as a gate fact.
2. MUST verify Success Metric 1 and Success Metric 2 by executing task_02's
   five runner tests, and confirm from the diff that `Reason`, `Err`, the
   exit-code mapping and the daemon packages did not change (API Contract 1,
   API Contract 2).
3. MUST verify Success Metric 3 and API Contract 3 and API Contract 4 by
   executing task_03's seven retry tests, and confirm from the diff that a
   `RunPrepared` failure is retried only when `PromptSent` is false and that
   no path records `none`.
4. MUST verify Success Metric 4: every existing review and runner test passes,
   and the only changed expectations are the three fallback tests task_03
   names.
5. MUST re-measure the outside evidence against the installed acpx, if
   present, with a fake ACP adapter the gate writes in a temporary directory
   and a temporary `HOME`: one `acpx prompt --format json` per failing step
   (`initialize`, `session/set_model`, `session/prompt`) and one adapter
   exit during `session/prompt`, recording stdout, stderr and exit, and feed
   each recorded stdout to the built runner to show the step and message the
   review would report. When acpx or Node is absent, record the row as
   blocked with that reason and rest it on `_prd.md` → Acceptance evidence.
6. MUST record, as evidence this Spec did not author, each source of
   `_prd.md` → Acceptance evidence: the Pantheon report in the Secondbrain,
   the acpx measurement above, and the Agent Client Protocol overview. For
   each source it reaches it MUST record what it read and whether it still
   supports the design; for each it cannot reach it MUST record the row as
   blocked with the reason.
7. MUST verify that `docs/user-guide/commands/review.md` and the Roundfix
   Skill's `review` reference describe the step, the message, the retry and
   the three record fields, that the Vocabulary Contract's patterns are
   documented where it says, that each skill mirror equals its canonical file,
   and that the raised owned-skill version is recorded.
8. MUST verify that this Spec's own artifacts satisfy the promise rule, and
   perform the glossary check of `docs/agents/domain.md`: record whether
   **Agent Selection** and **Fallback Selection** still describe the review's
   behavior with a same-selection retry, and whether "protocol step" needs a
   glossary term.
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
11. MUST NOT accept a row whose only evidence is that a file was read, and
    MUST NOT start a real review or send any request to a provider.

## Subtasks

- [ ] Build the binary and run the repository Verification.
- [ ] Exercise the runner and review rows.
- [ ] Re-measure acpx with a fake adapter or record the row blocked.
- [ ] Record the outside evidence, glossary check and scope audit.
- [ ] Write the QA report with the Pull Request row blocked.

## Acceptance Criteria

- [ ] Every behavior row passes on recorded evidence.
- [ ] The acpx re-measurement shows each failing step and message, or the row
      is blocked with its reason.
- [ ] The outside-evidence row records what each source says or why it was
      unreachable.
- [ ] The scope audit finds no undeclared Governed Path change.

## Verification

The following command is rendered unchanged from Roundfix's derived QA contract.
It proves a newest report exists and its verdict is readable; eligibility is
applied in process by whoever settles the Task.

- `newest="$(find 'docs/specs/0237-review-selection-failures-that-say-why/qa' -type f -name "qa-report-*.md" -print 2>/dev/null | awk "{ report=\$0; name=\$0; parts=split(name, path, \"/\"); name=path[parts]; name=substr(name, 11, length(name)-13); date=substr(name, 1, 10); suffix=substr(name, 11); shape=(length(date) == 10 && substr(date, 5, 1) == \"-\" && substr(date, 8, 1) == \"-\"); for (i=1; i <= 10 && shape; i++) { if (i != 5 && i != 8 && index(\"0123456789\", substr(date, i, 1)) == 0) shape=0 } year=substr(date, 1, 4)+0; month=substr(date, 6, 2)+0; day=substr(date, 9, 2)+0; leap=(year%400 == 0 || (year%4 == 0 && year%100 != 0)); days=(month == 2 ? 28+leap : ((month == 4 || month == 6 || month == 9 || month == 11) ? 30 : 31)); dated=(shape && month >= 1 && month <= 12 && day >= 1 && day <= days); if (!dated) date=\"\"; if (suffix == \"\") { sequenced=1; sequence=-1 } else { sequenced=(length(suffix) > 1 && substr(suffix, 1, 1) == \"-\"); for (i=2; i <= length(suffix) && sequenced; i++) { if (index(\"0123456789\", substr(suffix, i, 1)) == 0) sequenced=0 } sequence=(sequenced ? substr(suffix, 2)+0 : 0) } if (dated && sequenced) printf \"%s\\t%s\\t%s\\t%s\\t%s\\n\", dated, date, sequenced, sequence, report }" | sort -k1,1n -k2,2 -k3,3n -k4,4n -k5,5 | tail -1 | cut -f5-)"; test -n "$newest" || exit 1; awk "BEGIN { whitespace=\" \t\r\n\f\v\" } NR == 1 && \$0 == \"---\" { frontmatter=1; next } frontmatter && \$0 == \"---\" { closed=1; exit } frontmatter && index(\$0, \"verdict:\") == 1 { verdict=substr(\$0, 9); while (length(verdict) > 0 && index(whitespace, substr(verdict, 1, 1)) > 0) verdict=substr(verdict, 2); while (length(verdict) > 0 && index(whitespace, substr(verdict, length(verdict), 1)) > 0) verdict=substr(verdict, 1, length(verdict)-1); verdicts++ } END { exit(closed && verdicts == 1 && length(verdict) > 0 ? 0 : 1) }" "$newest"`

## References

- `_prd.md` → Goals; User Stories 1-3; Core Features 1-4; Success Metric 1; Success Metric 2; Success Metric 3; Success Metric 4; Acceptance evidence
- `_techspec.md` → API Contract 1; API Contract 2; API Contract 3; API Contract 4; Vocabulary Contract; Build Order 4
- ADR-0242; ADR-0080; ADR-0091; ADR-0104
