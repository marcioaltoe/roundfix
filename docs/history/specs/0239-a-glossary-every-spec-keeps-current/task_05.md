---
task: task_05
spec: 0239-a-glossary-every-spec-keeps-current
status: completed
type: qa
complexity: high
---

# Task 05: Run the final QA gate

## Overview

The authored terminal gate for this Spec. It declares what the matrix covers
and settles the Spec on evidence. This Spec makes the Spec Consistency Check
and the Archive Command hold a Glossary Declaration to the glossary, rewords
the Baseline's domain clauses, makes the authoring, QA and Roundfix skills
follow the rule, and adds the dropped terms to `CONTEXT.md`. Every behavior row
runs against the built tree in temporary repositories. No row reaches a
provider, calls the Jev judge, reads a credential, writes the real
`~/.roundfix`, or writes to an adopter repository.

## Requirements

1. MUST run the repository Verification and `make verify-docs`, and record
   each result as a gate fact.
2. MUST verify Success Metric 1, Success Metric 2 and Success Metric 3 by
   executing task_01's glossary and archive tests and the corpus tests, and by
   running the built binary on a temporary repository: `roundfix spec check`
   on a Spec with an undeclared bold term, with a planned term, and with a
   completed Task whose term is absent; and `roundfix archive` on the last,
   with and without `--qa-override`, recording each exit code and the
   Glossary Gap reason (API Contract 1, API Contract 2, API Contract 3,
   API Contract 4, API Contract 5).
3. MUST verify that this Spec passes its own detector: the built
   `roundfix spec check 0239-a-glossary-every-spec-keeps-current --strict`
   reports no glossary finding, and every term of its Glossary Declaration is
   defined in `CONTEXT.md`.
4. MUST verify Success Metric 4 by executing task_02's clause tests and a
   Managed Refresh dry run that reports the guides current (API Contract 6).
5. MUST verify Success Metric 5: for each catch-up entry, read the ADR
   `_techspec.md` → Catch-up glossary entries names for it and record whether
   the definition agrees, and confirm that no `GLOSSARY.md` exists (API
   Contract 8).
6. MUST verify Success Metric 6: the authoring skills name domain-modeling and
   the declaration on the autonomous route and keep grilling and
   grill-with-docs as the interactive entry, write-tasks gives each declared
   term a glossary requirement, the QA gate records declared terms, the
   Roundfix Skill's spec and archive references describe the codes and the
   refusal, each mirror equals its canonical file, and every raised version is
   recorded (API Contract 7). The Vocabulary Contract's pattern MUST be
   documented where it says.
7. MUST record, as evidence this Spec did not author, each source of
   `_prd.md` → Acceptance evidence. It MUST re-run the adopter measurement
   read-only (last change of each adopter's glossary, defined terms, Specs
   archived after it, from local Git history without fetching), never writing
   to an adopter or sending its content anywhere; and MUST reach the two
   published sources and the upstream changelog. For each source it reaches it
   MUST record what it read and whether it still supports the design; a source
   the Run sandbox cannot reach is recorded as blocked with its reason, and a
   published source refused only because the sandbox denied network access is
   recorded as `blocked (environment: network denied: <host>)`.
8. MUST perform the glossary check of `docs/agents/domain.md` for this Spec:
   record, for each term `_prd.md` → Glossary adds or changes, whether
   `CONTEXT.md` defines it after the work, and whether the Spec introduced any
   other term.
9. MUST verify from the repository history that each Task's changed files stay
   within its declarations or its `## Recorded paths`, that every Governed Path
   changed is bounded in `_authorization.md`, and that `Makefile`, `go.mod`,
   the CI workflows, the Source Baseline corpus and the upstream
   `domain-modeling`, `grilling` and `grill-with-docs` skills did not change.
   This row reads Task commits, so it declares a `commit_range` input.
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
    MUST NOT start a real review, call the Jev judge or send any request to a
    provider.

## Subtasks

- [ ] Build the binary and run the repository Verification and the docs gate.
- [ ] Exercise the detector, the archive refusal and the Baseline rows.
- [ ] Check the catch-up entries against their ADRs and the skills.
- [ ] Record the outside evidence, glossary check and scope audit.
- [ ] Write the QA report with the Pull Request row blocked.

## Acceptance Criteria

- [ ] Every behavior row passes on recorded evidence.
- [ ] This Spec passes its own glossary detector.
- [ ] The outside-evidence row records what each source says or why it was
      unreachable.
- [ ] The scope audit finds no undeclared Governed Path change.

## Verification

The following command is rendered unchanged from Roundfix's derived QA contract.
It proves a newest report exists and its verdict is readable; eligibility is
applied in process by whoever settles the Task.

- `newest="$(find 'docs/specs/0239-a-glossary-every-spec-keeps-current/qa' -type f -name "qa-report-*.md" -print 2>/dev/null | awk "{ report=\$0; name=\$0; parts=split(name, path, \"/\"); name=path[parts]; name=substr(name, 11, length(name)-13); date=substr(name, 1, 10); suffix=substr(name, 11); shape=(length(date) == 10 && substr(date, 5, 1) == \"-\" && substr(date, 8, 1) == \"-\"); for (i=1; i <= 10 && shape; i++) { if (i != 5 && i != 8 && index(\"0123456789\", substr(date, i, 1)) == 0) shape=0 } year=substr(date, 1, 4)+0; month=substr(date, 6, 2)+0; day=substr(date, 9, 2)+0; leap=(year%400 == 0 || (year%4 == 0 && year%100 != 0)); days=(month == 2 ? 28+leap : ((month == 4 || month == 6 || month == 9 || month == 11) ? 30 : 31)); dated=(shape && month >= 1 && month <= 12 && day >= 1 && day <= days); if (!dated) date=\"\"; if (suffix == \"\") { sequenced=1; sequence=-1 } else { sequenced=(length(suffix) > 1 && substr(suffix, 1, 1) == \"-\"); for (i=2; i <= length(suffix) && sequenced; i++) { if (index(\"0123456789\", substr(suffix, i, 1)) == 0) sequenced=0 } sequence=(sequenced ? substr(suffix, 2)+0 : 0) } if (dated && sequenced) printf \"%s\\t%s\\t%s\\t%s\\t%s\\n\", dated, date, sequenced, sequence, report }" | sort -k1,1n -k2,2 -k3,3n -k4,4n -k5,5 | tail -1 | cut -f5-)"; test -n "$newest" || exit 1; awk "BEGIN { whitespace=\" \t\r\n\f\v\" } NR == 1 && \$0 == \"---\" { frontmatter=1; next } frontmatter && \$0 == \"---\" { closed=1; exit } frontmatter && index(\$0, \"verdict:\") == 1 { verdict=substr(\$0, 9); while (length(verdict) > 0 && index(whitespace, substr(verdict, 1, 1)) > 0) verdict=substr(verdict, 2); while (length(verdict) > 0 && index(whitespace, substr(verdict, length(verdict), 1)) > 0) verdict=substr(verdict, 1, length(verdict)-1); verdicts++ } END { exit(closed && verdicts == 1 && length(verdict) > 0 ? 0 : 1) }" "$newest"`

## References

- `_prd.md` → Goals; User Stories 1-4; Core Features 1-7; Success Metric 1; Success Metric 2; Success Metric 3; Success Metric 4; Success Metric 5; Success Metric 6; Acceptance evidence; Glossary
- `_techspec.md` → API Contract 1; API Contract 2; API Contract 3; API Contract 4; API Contract 5; API Contract 6; API Contract 7; API Contract 8; Catch-up glossary entries; Vocabulary Contract; Build Order 5
- ADR-0244; ADR-0080; ADR-0091; ADR-0104; ADR-0240
