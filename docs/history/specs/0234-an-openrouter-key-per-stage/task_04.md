---
task: task_04
spec: 0234-an-openrouter-key-per-stage
status: completed
type: qa
complexity: high
---

# Task 04: Run the final QA gate

## Overview

The authored terminal gate for this Spec. It declares what the matrix covers
and settles the Spec on evidence. This Spec makes the Jev judge read
`ROUNDFIX_OPENROUTER_JUDGE_API_KEY` and implementation on an open model read
`ROUNDFIX_OPENROUTER_IMPLEMENT_API_KEY` before the shared
`ROUNDFIX_OPENROUTER_API_KEY`, records the variable each used, and lists each
stage's variables in the Doctor Command. Every behavior row runs against the
built tree with fake transports, fake environments and a temporary Roundfix
home. No row reaches OpenRouter or TypeSafe, sets a real key, or reads or
writes the real `~/.roundfix`.

## Requirements

1. MUST run the repository Verification and record its result as a gate fact.
2. MUST verify Success Metric 1 and Success Metric 2 by executing task_01's
   judge and `spec judge` tests, and verify Surface Transcript 1 and Surface
   Transcript 2 through the built `roundfix spec judge` in a disposable
   repository with an active Spec, a temporary home and an environment that
   sets none of the three judge variables: stdout, stderr and exit match.
3. MUST verify Success Metric 3 by executing task_03's three tests, and
   confirm from the diff that Spec 0233's helper calls the implementation
   stage's `Select` and that no other reader names the shared key for that
   stage.
4. MUST verify Success Metric 4 and API Contract 3 by executing task_02's
   doctor tests, and by running the built `roundfix doctor` with fake
   sentinel values in the judge and implementation variables: the
   `environment:` line lists the five entries in order and the output
   contains no sentinel.
5. MUST verify API Contract 1, API Contract 2 and API Contract 4 from the
   tests and the diff: the summary and JSON name the variable, every Judge
   Log line carries `key_variable`, each implementation spend record carries
   it, and the generic `OPENROUTER_API_KEY` is read nowhere.
6. MUST verify Success Metric 5: every existing judge, doctor and light-tier
   test passes, and the only changed expectations are those task_01,
   task_02 and task_03 declare.
7. MUST verify that `docs/user-guide/commands/spec.md`,
   `docs/user-guide/commands/doctor.md`, `docs/user-guide/configuration.md`,
   the Roundfix Skill's spec and runtime references, and the write-prd and
   write-techspec skills name the stage keys, the order and the fallback,
   that the Vocabulary Contract's patterns are documented where it says,
   that each skill mirror equals its canonical file, and that every raised
   owned-skill version is recorded.
8. MUST record, as evidence this Spec did not author, each source of
   `_prd.md` → Acceptance evidence: OpenRouter's Activity guide, its API key
   reference and its limits reference, and the Measurement section of the
   Backlog Entry "A judge-assigned model tier per Task" of 2026-10-05. For
   each source it reaches it MUST record what it read and whether it still
   supports the design; for each it cannot reach it MUST record the row as
   blocked with the reason.
9. MUST verify that this Spec's own artifacts satisfy the promise rule, and
   perform the glossary check of `docs/agents/domain.md`: record whether
   **Judge Log** still describes the record now that it names a key
   variable, and whether "stage key" needs a glossary term.
10. MUST verify from the repository history that each Task's changed files
    stay within its declarations or its `## Recorded paths`, that every
    Governed Path changed is bounded in `_authorization.md`, and that
    `Makefile`, `go.mod` and the CI workflows did not change. This row reads
    Task commits, so it declares a `commit_range` input.
11. MUST record the non-waivable Pull Request row as
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
12. MUST NOT accept a row whose only evidence is that a file was read, and
    MUST NOT set a real OpenRouter or TypeSafe key or send any request to
    either service.

## Subtasks

- [ ] Build the binary and run the repository Verification.
- [ ] Exercise the judge, doctor and implementation rows.
- [ ] Reproduce both Surface Transcripts through the built binary.
- [ ] Record the outside evidence, glossary check and scope audit.
- [ ] Write the QA report with the Pull Request row blocked.

## Acceptance Criteria

- [ ] Every behavior row passes on recorded evidence.
- [ ] Both Surface Transcripts match through the built binary.
- [ ] The outside-evidence row records what OpenRouter's documentation says
      or why it was unreachable.
- [ ] The scope audit finds no undeclared Governed Path change.

## Verification

The following command is rendered unchanged from Roundfix's derived QA contract.
It proves a newest report exists and its verdict is readable; eligibility is
applied in process by whoever settles the Task.

- `newest="$(find 'docs/specs/0234-an-openrouter-key-per-stage/qa' -type f -name "qa-report-*.md" -print 2>/dev/null | awk "{ report=\$0; name=\$0; parts=split(name, path, \"/\"); name=path[parts]; name=substr(name, 11, length(name)-13); date=substr(name, 1, 10); suffix=substr(name, 11); shape=(length(date) == 10 && substr(date, 5, 1) == \"-\" && substr(date, 8, 1) == \"-\"); for (i=1; i <= 10 && shape; i++) { if (i != 5 && i != 8 && index(\"0123456789\", substr(date, i, 1)) == 0) shape=0 } year=substr(date, 1, 4)+0; month=substr(date, 6, 2)+0; day=substr(date, 9, 2)+0; leap=(year%400 == 0 || (year%4 == 0 && year%100 != 0)); days=(month == 2 ? 28+leap : ((month == 4 || month == 6 || month == 9 || month == 11) ? 30 : 31)); dated=(shape && month >= 1 && month <= 12 && day >= 1 && day <= days); if (!dated) date=\"\"; if (suffix == \"\") { sequenced=1; sequence=-1 } else { sequenced=(length(suffix) > 1 && substr(suffix, 1, 1) == \"-\"); for (i=2; i <= length(suffix) && sequenced; i++) { if (index(\"0123456789\", substr(suffix, i, 1)) == 0) sequenced=0 } sequence=(sequenced ? substr(suffix, 2)+0 : 0) } if (dated && sequenced) printf \"%s\\t%s\\t%s\\t%s\\t%s\\n\", dated, date, sequenced, sequence, report }" | sort -k1,1n -k2,2 -k3,3n -k4,4n -k5,5 | tail -1 | cut -f5-)"; test -n "$newest" || exit 1; awk "BEGIN { whitespace=\" \t\r\n\f\v\" } NR == 1 && \$0 == \"---\" { frontmatter=1; next } frontmatter && \$0 == \"---\" { closed=1; exit } frontmatter && index(\$0, \"verdict:\") == 1 { verdict=substr(\$0, 9); while (length(verdict) > 0 && index(whitespace, substr(verdict, 1, 1)) > 0) verdict=substr(verdict, 2); while (length(verdict) > 0 && index(whitespace, substr(verdict, length(verdict), 1)) > 0) verdict=substr(verdict, 1, length(verdict)-1); verdicts++ } END { exit(closed && verdicts == 1 && length(verdict) > 0 ? 0 : 1) }" "$newest"`

## References

- `_prd.md` → Goals; User Stories 1-4; Core Features 1-6; Success Metric 1; Success Metric 2; Success Metric 3; Success Metric 4; Success Metric 5; Acceptance evidence
- `_techspec.md` → API Contract 1; API Contract 2; API Contract 3; API Contract 4; Surface Transcript 1; Surface Transcript 2; Vocabulary Contract; Build Order 4
- ADR-0239; ADR-0080; ADR-0091; ADR-0104
